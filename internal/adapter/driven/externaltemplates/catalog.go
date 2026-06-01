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
	"bytes"
	"context"
	"embed"
	"fmt"
	iofs "io/fs"
	"path"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// templatesFS holds every built-in template's metadata file plus
// (eventually) every datafile. The embed pattern is `templates/*`
// recursive so slice-v1-template-init's per-template file-tree
// can use the same FS without a second embed declaration.
//
//go:embed templates/*
var templatesFS embed.FS

// catalogRoot is the directory inside [templatesFS] that holds the
// per-template subdirectories. Constant — the embed pattern above
// matches the same prefix; same constant is used by the test seam
// fixture FS so the production and test paths stay parallel.
const catalogRoot = "templates"

// metadataFile is the per-template metadata filename. ADR-0009
// §Entscheidung pins it; any directory inside [catalogRoot] without
// this file is treated as an invalid template (caught at List time
// so a packaging mistake fails fast in CI).
const metadataFile = "template.yaml"

// supportedAPIVersion is the only `apiVersion` value the adapter
// accepts. ADR-0009 §Entscheidung pins the v1 schema; future
// schema evolution must land in a new slice that bumps this
// constant (and either rejects or migrates v1-shaped files). The
// gate runs at load time so a packaging mistake — or a future
// template author copying a v2 example — fails fast with
// `domain.ErrInvalidTemplate` rather than silently rendering
// under wrong assumptions (review-followup N4).
const supportedAPIVersion = "github.com/pt9912/u-boot/template/v1"

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
//   - A directory missing its [metadataFile], or with a malformed
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
		meta, err := readTemplate(c.fs, path.Join(catalogRoot, ent.Name()))
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

// readTemplate parses a single `template.yaml` into a domain
// [TemplateMetadata] and validates it. The intermediate
// [rawTemplateYAML] struct keeps yaml.v3 tag knowledge inside the
// adapter; the domain type stays free of YAML-library imports per
// the depguard `domain-isoliert` rule.
//
// Strict decoding (review-followup N1): the decoder is run with
// KnownFields(true) so a typo like `requiredTool:` (singular) or
// `addOns:` (instead of `supportedAddOns:`) fails at load time
// instead of silently dropping the author's intent. Same protection
// covers stray vendor extensions ahead of any schema-v2 work.
//
// apiVersion gate (review-followup N4): rawTemplateYAML.APIVersion
// must equal [supportedAPIVersion]; a template carrying a future
// version (e.g. v2) is rejected with a `domain.ErrInvalidTemplate`
// wrap so the message is uniform with the other validation paths.
func readTemplate(fs iofs.FS, dir string) (domain.TemplateMetadata, error) {
	metaPath := path.Join(dir, metadataFile)
	data, err := iofs.ReadFile(fs, metaPath)
	if err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("read %s: %w", metaPath, err)
	}
	var raw rawTemplateYAML
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("parse %s: %w", metaPath, err)
	}
	if raw.APIVersion != supportedAPIVersion {
		return domain.TemplateMetadata{}, fmt.Errorf("%s: %w: apiVersion %q is not supported (want %q)",
			metaPath, domain.ErrInvalidTemplate, raw.APIVersion, supportedAPIVersion)
	}
	meta := raw.toDomain()
	if err := meta.Validate(); err != nil {
		return domain.TemplateMetadata{}, fmt.Errorf("%s: %w", metaPath, err)
	}
	return meta, nil
}

// rawTemplateYAML mirrors the on-disk schema per ADR-0009
// §Entscheidung. Private so the adapter is the only place that
// knows the `apiVersion`-tagged YAML shape; the application layer
// consumes only the curated domain projection.
//
// `apiVersion` is parsed but not validated at the listing stage —
// slice-v1-template-init can introduce a version-rejection sentinel
// when it needs to gate render behaviour on the schema version.
type rawTemplateYAML struct {
	APIVersion      string           `yaml:"apiVersion"`
	Name            string           `yaml:"name"`
	Description     string           `yaml:"description"`
	Version         string           `yaml:"version"`
	SupportedAddOns []string         `yaml:"supportedAddOns"`
	GeneratedFiles  []string         `yaml:"generatedFiles"`
	RequiredTools   []string         `yaml:"requiredTools"`
	Variables       []rawTemplateVar `yaml:"variables"`
}

// rawTemplateVar is the on-disk variable schema.
type rawTemplateVar struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
}

// toDomain copies the raw YAML projection into the domain value.
// Flat copy by design — see ADR-0009 §Entscheidung; any future
// schema evolution (vX → vX+1) lives in the adapter so the domain
// shape stays stable.
func (r rawTemplateYAML) toDomain() domain.TemplateMetadata {
	out := domain.TemplateMetadata{
		Name:            r.Name,
		Description:     r.Description,
		Version:         r.Version,
		SupportedAddOns: r.SupportedAddOns,
		GeneratedFiles:  r.GeneratedFiles,
		RequiredTools:   r.RequiredTools,
	}
	if len(r.Variables) > 0 {
		out.Variables = make([]domain.TemplateVariable, len(r.Variables))
		for i, v := range r.Variables {
			out.Variables[i] = domain.TemplateVariable{
				Name:        v.Name,
				Description: v.Description,
				Default:     v.Default,
				Required:    v.Required,
			}
		}
	}
	return out
}
