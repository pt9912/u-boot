// Package externaltemplates is the embed.FS-backed implementation of
// the `port/driven.TemplateCatalog` interface (LH-FA-ARCH-002), per
// ADR-0009 §Entscheidung. The package bundles the built-in external
// project templates that ship with u-boot. `template.yaml` metadata
// is enumerated by `u-boot template list` (slice-v1-template-list)
// and (later) resolved by name from `u-boot init --template <name>`
// (slice-v1-template-init).
//
// Directory layout (mirrors ADR-0009 §Entscheidung):
//
//	externaltemplates/
//	├── catalog.go
//	├── catalog_test.go
//	└── templates/
//	    └── <name>/
//	        ├── template.yaml   # metadata (LH-FA-TPL-002)
//	        └── …                # file-templates (added by slice-v1-template-init)
package externaltemplates

import (
	"context"
	"embed"
	"fmt"
	iofs "io/fs"
	"path"
	"sort"
	"strings"

	"github.com/pt9912/u-boot/internal/adapter/driven/templateyaml"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// templatesFS holds every built-in template's metadata file plus
// every per-template datafile. The pattern uses the `all:` prefix
// so dotfile templates like `.gitignore.tmpl` and `.env.example.tmpl`
// (slice-v1-template-init T3) are embedded — Go's default
// `templates/*` would silently drop them per the `embed` package
// dotfile-exclusion rule.
//
//go:embed all:templates/*
var templatesFS embed.FS

// catalogRoot is the directory inside [templatesFS] that holds the
// per-template subdirectories. Constant — the embed pattern above
// matches the same prefix; same constant is used by the test seam
// fixture FS so the production and test paths stay parallel.
const catalogRoot = "templates"

// Catalog is the production TemplateCatalog adapter. The zero value
// is NOT usable — production callers go through [New] (or [NewWithFS]
// for tests) so the embedded FS handle is set.
type Catalog struct {
	// fs is the source of template directories. Production uses the
	// package-level [templatesFS]; tests inject a [fstest.MapFS]
	// via [NewWithFS] to avoid touching the production set.
	fs iofs.FS
}

// Static check: Catalog satisfies the TemplateCatalog port.
var _ driven.TemplateCatalog = (*Catalog)(nil)

// Static check: the same Catalog satisfies the TemplateFiles port
// (slice-v1-template-init T1). Same struct, two roles — listing
// (List) and per-template file-tree access (Open). The wiring layer
// passes the same instance to both consumer services.
var _ driven.TemplateFiles = (*Catalog)(nil)

// New returns a production Catalog backed by the embedded
// `templates/` bundle. The bundle currently contains the `basic`
// bootstrap template per ADR-0009 §Folgepunkte 4; further built-ins
// land in their own slices on concrete demand.
func New() *Catalog { return &Catalog{fs: templatesFS} }

// NewWithFS is the test seam — wraps an arbitrary iofs.FS so the
// adapter's table-driven tests can supply fixture trees rooted at
// [catalogRoot]. Production code uses [New].
func NewWithFS(fs iofs.FS) *Catalog { return &Catalog{fs: fs} }

// List walks the catalog root, parses every per-template
// `template.yaml`, and returns the metadata list sorted by
// [domain.TemplateMetadata.Name]. Non-directory entries at the
// catalog root (e.g. a stray `README.md`) are skipped silently —
// only directories are treated as templates.
//
// Errors:
//
//   - ctx cancelled before any IO happens → returns ctx.Err() so
//     callers wiring a short deadline observe cancellation even
//     though embed.FS itself runs in microseconds (review-followup
//     N5 — the port contract claims `ctx is honored` and the
//     adapter now does, instead of relying on the documented
//     special case).
//   - The catalog root must exist; a missing root surfaces as a
//     wrapped [iofs.PathError] (production embed guarantees it
//     exists, so a runtime miss is a packaging bug).
//   - A directory missing its [templateyaml.MetadataFile], or with a malformed
//     `template.yaml`, fails the call — partial enumeration would
//     hide build problems. The per-template prefix in the wrap
//     ("template %q: …") lets the caller pinpoint the offender.
func (c *Catalog) List(ctx context.Context) ([]domain.TemplateMetadata, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	entries, err := iofs.ReadDir(c.fs, catalogRoot)
	if err != nil {
		return nil, fmt.Errorf("read catalog root %q: %w", catalogRoot, err)
	}

	out := make([]domain.TemplateMetadata, 0, len(entries))
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		meta, err := templateyaml.Read(c.fs, path.Join(catalogRoot, ent.Name()))
		if err != nil {
			return nil, fmt.Errorf("template %q: %w", ent.Name(), err)
		}
		out = append(out, meta)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Open returns an [iofs.FS] rooted at the per-template directory
// (slice-v1-template-init T1). The application render-loop walks
// the result to discover `.tmpl` files (rendered) and plain files
// (copied 1:1) per ADR-0009 §Entscheidung. Existence is checked via
// `iofs.ReadDir` before `iofs.Sub`; a missing or empty name maps to
// [driven.ErrTemplateNotFound] so the application service can wrap
// it for LH-FA-CLI-006 code-10 mapping.
//
// Returned FS includes `template.yaml`; callers that only want
// renderable files filter it out themselves (the rendering pipeline
// doesn't need the metadata, since [List] already validated it).
func (c *Catalog) Open(ctx context.Context, name string) (iofs.FS, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if name == "" {
		return nil, fmt.Errorf("%w: empty template name", driven.ErrTemplateNotFound)
	}
	// Defense-in-depth (Post-Closure-Review #1): a catalog name is a bare
	// identifier. Reject path-shaped names (`.`/`..`/separators) so a
	// mis-classified ref can never resolve to the catalog ROOT (`.` →
	// path.Join collapses to catalogRoot) or escape it (`..`). The
	// Composite already routes these to the FS resolver via
	// domain.ClassifyTemplateRef; this guard holds even if a future
	// caller bypasses that classification.
	if name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return nil, fmt.Errorf("%w: %q is not a catalog name", driven.ErrTemplateNotFound, name)
	}
	subDir := path.Join(catalogRoot, name)
	entries, err := iofs.ReadDir(c.fs, subDir)
	if err != nil {
		// Both "directory does not exist" and "name points at a
		// non-directory" collapse to template-not-found. The user
		// cares that "the name is invalid", not why iofs failed.
		return nil, fmt.Errorf("%w: %q", driven.ErrTemplateNotFound, name)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("%w: %q (directory is empty)", driven.ErrTemplateNotFound, name)
	}
	sub, err := iofs.Sub(c.fs, subDir)
	if err != nil {
		return nil, fmt.Errorf("template subfs %q: %w", name, err)
	}
	return sub, nil
}
