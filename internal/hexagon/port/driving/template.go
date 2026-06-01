package driving

import (
	"context"
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// TemplateListRequest is the input for [TemplateListUseCase.List].
// Empty today — the listing returns the full catalog. Future slices
// (template-init, local-templates) MAY add filter fields without
// breaking the request shape (additive only).
type TemplateListRequest struct{}

// TemplateListResponse is the output of [TemplateListUseCase.List].
// Templates is the catalog list sorted by Name, as the driven
// adapter guarantees (`port/driven.TemplateCatalog` contract); the
// CLI adapter renders it without re-sorting.
type TemplateListResponse struct {
	Templates []domain.TemplateMetadata
}

// ErrTemplateCatalog wraps adapter-side failures from the
// [driven.TemplateCatalog]: filesystem IO, YAML parse, and metadata
// validation (e.g. a malformed embedded `template.yaml`). All map
// to LH-FA-CLI-006 exit code 14 (technical persistence/data
// failure) — invalid embedded metadata is a packaging bug the CI
// build should have caught, not a user-actionable error.
//
// The wrapped chain preserves the original cause via `%w`, so
// `errors.Is(err, domain.ErrInvalidTemplate)` still works at the
// CLI for the rare case where a debug build wires a fixture catalog
// with broken entries.
var ErrTemplateCatalog = errors.New("template: catalog error")

// TemplateListUseCase implements `u-boot template list`
// (LH-FA-TPL-004). Read-only — never mutates state.
//
// Contract:
//
//   - The returned Templates slice mirrors the catalog order
//     (sorted by Name); the CLI adapter renders it as-is.
//   - On success Templates MAY be empty (a catalog with zero
//     entries — defensively allowed but not expected in
//     production); the CLI adapter renders an empty-state message
//     in that case.
//   - On failure the response is the zero value and the error
//     wraps [ErrTemplateCatalog].
//
// ctx is honored so a slow future adapter (network-backed catalog)
// can be cancelled; the embed.FS-backed adapter shipped with
// slice-v1-template-list ignores it.
type TemplateListUseCase interface {
	List(ctx context.Context, req TemplateListRequest) (TemplateListResponse, error)
}
