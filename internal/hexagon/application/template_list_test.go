package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// fakeCatalog implements driven.TemplateCatalog for the
// TemplateListService unit tests. Returns the configured metas /
// err verbatim and records the last-seen context so the
// context-propagation pin can assert on it.
type fakeCatalog struct {
	metas      []domain.TemplateMetadata
	err        error
	lastCtx    context.Context //nolint:containedctx // test-only seam for context-propagation pin
	callCount  int
}

func (f *fakeCatalog) List(ctx context.Context) ([]domain.TemplateMetadata, error) {
	f.callCount++
	f.lastCtx = ctx
	return f.metas, f.err
}

func TestTemplateListService_List_DelegatesToCatalog(t *testing.T) {
	t.Parallel()
	want := []domain.TemplateMetadata{
		{Name: "alpha", Description: "first", Version: "1.0.0"},
		{Name: "basic", Description: "skeleton", Version: "0.1.0"},
	}
	cat := &fakeCatalog{metas: want}
	svc := application.NewTemplateListService(cat, nil)

	resp, err := svc.List(context.Background(), driving.TemplateListRequest{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if cat.callCount != 1 {
		t.Errorf("catalog.callCount = %d, want 1", cat.callCount)
	}
	if len(resp.Templates) != 2 {
		t.Fatalf("len(Templates) = %d, want 2", len(resp.Templates))
	}
	for i, m := range want {
		if resp.Templates[i].Name != m.Name {
			t.Errorf("Templates[%d].Name = %q, want %q (order must match catalog)", i, resp.Templates[i].Name, m.Name)
		}
	}
}

func TestTemplateListService_List_EmptyCatalogIsNotError(t *testing.T) {
	t.Parallel()
	svc := application.NewTemplateListService(&fakeCatalog{metas: nil}, nil)

	resp, err := svc.List(context.Background(), driving.TemplateListRequest{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(resp.Templates) != 0 {
		t.Errorf("len(Templates) = %d, want 0 (empty catalog → empty response, no error)", len(resp.Templates))
	}
}

func TestTemplateListService_List_CatalogErrorWrapsSentinel(t *testing.T) {
	t.Parallel()
	bareErr := errors.New("catalog boom")
	svc := application.NewTemplateListService(&fakeCatalog{err: bareErr}, nil)

	resp, err := svc.List(context.Background(), driving.TemplateListRequest{})
	if err == nil {
		t.Fatal("List: want error, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateCatalog) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateCatalog", err)
	}
	if !errors.Is(err, bareErr) {
		t.Errorf("err = %v, want wrap of original cause %v (multi-%%w chain)", err, bareErr)
	}
	if len(resp.Templates) != 0 {
		t.Errorf("Templates = %v, want zero value on error", resp.Templates)
	}
}

func TestTemplateListService_List_ErrInvalidTemplateChainPreserved(t *testing.T) {
	t.Parallel()
	// Simulate the catalog returning an ErrInvalidTemplate-wrapped
	// error — slice-v1-template-list T1 adapter does this when a
	// fixture catalog has a malformed template.yaml. The service
	// wraps with ErrTemplateCatalog but must preserve the domain
	// sentinel via multi-%w so debug-build callers (or future tests
	// wiring a fixture catalog) can still classify.
	validateErr := fmt.Errorf("templates/broken/template.yaml: %w", domain.ErrInvalidTemplate)
	svc := application.NewTemplateListService(&fakeCatalog{err: validateErr}, nil)

	_, err := svc.List(context.Background(), driving.TemplateListRequest{})
	if err == nil {
		t.Fatal("want error, got nil")
	}
	if !errors.Is(err, driving.ErrTemplateCatalog) {
		t.Errorf("err = %v, want wrap of driving.ErrTemplateCatalog", err)
	}
	if !errors.Is(err, domain.ErrInvalidTemplate) {
		t.Errorf("err = %v, want preserve of domain.ErrInvalidTemplate through multi-%%w", err)
	}
}

func TestTemplateListService_List_PropagatesContext(t *testing.T) {
	t.Parallel()
	cat := &fakeCatalog{metas: nil}
	svc := application.NewTemplateListService(cat, nil)

	type ctxKey int
	const sentinel ctxKey = 42
	ctx := context.WithValue(context.Background(), sentinel, "marker")

	if _, err := svc.List(ctx, driving.TemplateListRequest{}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if cat.lastCtx == nil {
		t.Fatal("catalog.lastCtx = nil, want propagated context")
	}
	if got := cat.lastCtx.Value(sentinel); got != "marker" {
		t.Errorf("catalog ctx value = %v, want %q (context must propagate verbatim)", got, "marker")
	}
}

func TestTemplateListService_New_NilLoggerIsAccepted(t *testing.T) {
	t.Parallel()
	// Pattern-pin: matches ConfigService / DoctorService / etc. The
	// nil-tolerant logger fallback removes nil-guard noise from the
	// hot path in List.
	svc := application.NewTemplateListService(&fakeCatalog{}, nil)
	if svc == nil {
		t.Fatal("New returned nil")
	}
	if _, err := svc.List(context.Background(), driving.TemplateListRequest{}); err != nil {
		t.Fatalf("List with nil logger: %v", err)
	}
}
