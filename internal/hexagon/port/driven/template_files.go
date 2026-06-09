package driven

import (
	"context"
	"errors"
	iofs "io/fs"
)

// ErrTemplateNotFound signals that the requested template name does
// not exist in the catalog (typo, mis-cased name, or a name the
// shipped catalog never had). Adapter contract — the application
// service for `u-boot init --template <name>` wraps this with a
// driving-port sentinel that maps to LH-FA-CLI-006 code 10
// (user-actionable: pick a real template or generate one).
//
// Lives in the driven port (not the adapter) so the application
// layer can branch on it via `errors.Is` without importing the
// concrete catalog implementation (depguard `application-no-
// adapter` rule).
var ErrTemplateNotFound = errors.New("template not found")

// ErrTemplateInvalid signals that a template was found but its
// `template.yaml` is malformed, carries an unsupported apiVersion, or
// fails the LH-FA-TPL-002 metadata minimum. Distinct from
// [ErrTemplateNotFound] (the template is present, just not valid).
//
// Produced by the filesystem-backed resolver (slice-later-local-
// templates), which wraps the layer-neutral
// `templateyaml.Read` / [domain.ErrInvalidTemplate] result with this
// sentinel. The embed.FS catalog validates at build time and does not
// emit it on the user path. The application service for
// `u-boot init --template <ref>` wraps it with a driving-port sentinel
// that maps to LH-FA-CLI-006 exit code 10 (user-actionable: fix the
// template metadata) — NOT exit 14 (technical render failure).
//
// Lives in the driven port (not the adapter) so the application layer
// can branch on it via `errors.Is` without importing the concrete
// resolver (depguard `application-no-adapter` rule).
var ErrTemplateInvalid = errors.New("template invalid")

// TemplateFiles exposes the per-template file tree for the render
// path of `u-boot init --template <name>` (slice-v1-template-init).
// Read-only; complements [TemplateCatalog]'s listing role with the
// asset-access role.
//
// Why a separate port: the listing slice's [TemplateCatalog] is
// already in production (with at least one external fake in test
// code); a new method on it would force every implementor to ship
// a body. Two ports lets the production adapter satisfy both with
// a single struct while existing fakes stay untouched. Tests that
// want render-only behaviour wire only [TemplateFiles].
//
// Contract:
//
//   - Open returns an [iofs.FS] rooted at the template's per-template
//     directory. Callers walk it with [iofs.WalkDir] / read with
//     [iofs.ReadFile]. `template.yaml` is included in the returned
//     tree (callers ignore it if they only want renderable files).
//   - Unknown / mistyped name → [ErrTemplateNotFound]. Empty name
//     gets the same sentinel — the catalog has no nameless template.
//   - ctx is honored on the entry path (same convention as
//     [TemplateCatalog.List]).
type TemplateFiles interface {
	Open(ctx context.Context, name string) (iofs.FS, error)
}
