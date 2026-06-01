package driven

import (
	"context"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// TemplateCatalog enumerates the external project templates u-boot
// can render via `u-boot init --template <name>`. The interface is
// read-only — slice-v1-template-list (this slice) needs only the
// listing path; the future slice-v1-template-init Adds a separate
// per-template file-tree accessor.
//
// Adapters back the catalog from any source — the production
// adapter (`internal/adapter/driven/externaltemplates`) embeds a
// curated set via `embed.FS` per ADR-0009 §Entscheidung; a future
// slice-later-local-templates adapter walks the filesystem for
// `--template ./pfad`.
//
// Contract:
//
//   - The returned slice is sorted by [domain.TemplateMetadata.Name]
//     so callers can render deterministic output without re-sorting.
//   - Each element passes [domain.TemplateMetadata.Validate] — the
//     adapter validates at load time so the application service
//     does not need to re-check.
//   - Errors propagate raw; the application layer wraps them with
//     a driving-port sentinel before they leave the hexagon.
//
// ctx is honored so a slow future adapter (network-backed catalog,
// large local tree) can be cancelled; the production embed.FS
// adapter ignores it.
type TemplateCatalog interface {
	List(ctx context.Context) ([]domain.TemplateMetadata, error)
}
