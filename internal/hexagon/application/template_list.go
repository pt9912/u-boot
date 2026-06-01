package application

import (
	"context"
	"fmt"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TemplateListService implements [driving.TemplateListUseCase] for
// `u-boot template list` (LH-FA-TPL-004). The service is a thin
// pass-through over the [driven.TemplateCatalog] port — sorting and
// validation already happen in the adapter, so the application
// layer only handles error-wrapping and the noop-logger fallback.
//
// Stateless and concurrent-safe: every call delegates to the
// catalog with the caller's context.
type TemplateListService struct {
	catalog driven.TemplateCatalog
	logger  driven.Logger
}

// Static check: TemplateListService satisfies the driving port.
var _ driving.TemplateListUseCase = (*TemplateListService)(nil)

// NewTemplateListService constructs the service with the driven
// adapters injected by the wiring layer. catalog is mandatory;
// logger accepts nil and is routed to [noopLogger] (matching the
// other services' nil-tolerant pattern).
func NewTemplateListService(catalog driven.TemplateCatalog, logger driven.Logger) *TemplateListService {
	if logger == nil {
		logger = noopLogger{}
	}
	return &TemplateListService{catalog: catalog, logger: logger}
}

// List delegates to the catalog and wraps any failure with
// [driving.ErrTemplateCatalog]. The wrap uses Go 1.20+ multi-`%w`
// so callers can still `errors.Is(err, domain.ErrInvalidTemplate)`
// for catalog-validation diagnostics in tests or future surface
// areas that want to distinguish parse from validation failures.
func (s *TemplateListService) List(ctx context.Context, _ driving.TemplateListRequest) (driving.TemplateListResponse, error) {
	metas, err := s.catalog.List(ctx)
	if err != nil {
		s.logger.Debug("template list: catalog error", "err", err.Error())
		return driving.TemplateListResponse{}, fmt.Errorf("%w: %w", driving.ErrTemplateCatalog, err)
	}
	s.logger.Debug("template list: catalog ok", "count", len(metas))
	return driving.TemplateListResponse{Templates: metas}, nil
}
