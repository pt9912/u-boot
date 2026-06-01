package application_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// LH-FA-ADD-006 four-mode dispatch (slice-v1-addons-deps T3).
//
// Strategy: drive the unexported [AddServiceService.handle-
// MissingDependencies] directly via the [HandleMissing-
// DependenciesForTest] export seam with a synthetic missing list
// containing an unsupported service name ("ghost-service"). The
// recursive sub-Add fails fast on the catalogue check
// (driving.ErrServiceUnsupported) without touching the disk —
// observing that error in the chain proves the dispatch took the
// autoInstall arm, while its absence (paired with
// driving.ErrDependenciesRequired) proves the fail-fast arm.

func mustNewServiceName(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}

func newDispatchService(t *testing.T, conf *fakeConfirmer) *application.AddServiceService {
	t.Helper()
	return application.NewAddServiceService(newFakeFS(), &fakeYAML{}, conf, nil)
}

func TestHandleMissingDeps_WithDeps_AutoInstallsWithoutPrompt(t *testing.T) {
	t.Parallel()
	conf := &fakeConfirmer{}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "postgres"),
		WithDeps:    true,
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want wrap of ErrServiceUnsupported (proves recursive Add fired)", err)
	}
	if errors.Is(err, driving.ErrDependenciesRequired) {
		t.Errorf("err = %v, must NOT carry ErrDependenciesRequired in WithDeps mode", err)
	}
	if len(conf.addDepCalls) != 0 {
		t.Errorf("ConfirmAddDependency called %d times, want 0 (WithDeps skips prompt)", len(conf.addDepCalls))
	}
}

func TestHandleMissingDeps_Yes_AutoInstallsWithoutPrompt(t *testing.T) {
	t.Parallel()
	conf := &fakeConfirmer{}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "postgres"),
		Yes:         true,
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want wrap of ErrServiceUnsupported", err)
	}
	if len(conf.addDepCalls) != 0 {
		t.Errorf("ConfirmAddDependency called %d times, want 0 (--yes skips prompt)", len(conf.addDepCalls))
	}
}

func TestHandleMissingDeps_NoInteractive_FailFastWithoutPrompt(t *testing.T) {
	t.Parallel()
	conf := &fakeConfirmer{}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:       "/proj",
		ServiceName:   mustNewServiceName(t, "postgres"),
		NoInteractive: true,
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrDependenciesRequired) {
		t.Fatalf("err = %v, want wrap of ErrDependenciesRequired", err)
	}
	if errors.Is(err, driving.ErrServiceUnsupported) {
		t.Errorf("err = %v, must NOT trigger recursive Add in --no-interactive fail-fast", err)
	}
	if len(conf.addDepCalls) != 0 {
		t.Errorf("ConfirmAddDependency called %d times, want 0 (--no-interactive skips prompt)", len(conf.addDepCalls))
	}
	if !strings.Contains(err.Error(), "--with-deps") {
		t.Errorf("err = %v, missing --with-deps hint in message", err)
	}
	if !strings.Contains(err.Error(), "ghost-service") {
		t.Errorf("err = %v, missing service name in message", err)
	}
}

func TestHandleMissingDeps_Default_PromptYes_PromotesToAutoInstall(t *testing.T) {
	t.Parallel()
	conf := &fakeConfirmer{addDepAnswer: true}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "postgres"),
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want wrap of ErrServiceUnsupported (recursive Add after prompt-yes)", err)
	}
	if len(conf.addDepCalls) != 1 {
		t.Fatalf("ConfirmAddDependency called %d times, want exactly 1", len(conf.addDepCalls))
	}
	call := conf.addDepCalls[0]
	if call.Service != "postgres" {
		t.Errorf("ConfirmAddDependency.Service = %q, want %q", call.Service, "postgres")
	}
	if len(call.Missing) != 1 || call.Missing[0] != "ghost-service" {
		t.Errorf("ConfirmAddDependency.Missing = %v, want [ghost-service]", call.Missing)
	}
}

func TestHandleMissingDeps_Default_PromptNo_FailsFast(t *testing.T) {
	t.Parallel()
	conf := &fakeConfirmer{addDepAnswer: false}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "postgres"),
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrDependenciesRequired) {
		t.Fatalf("err = %v, want wrap of ErrDependenciesRequired", err)
	}
	if errors.Is(err, driving.ErrServiceUnsupported) {
		t.Errorf("err = %v, recursive Add must NOT fire after user declined", err)
	}
	if len(conf.addDepCalls) != 1 {
		t.Errorf("ConfirmAddDependency called %d times, want 1", len(conf.addDepCalls))
	}
}

func TestHandleMissingDeps_Default_PromptErrors_SurfacesError(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("confirmer-io-failed")
	conf := &fakeConfirmer{addDepErr: sentinel}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:     "/proj",
		ServiceName: mustNewServiceName(t, "postgres"),
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want wrap of confirmer sentinel", err)
	}
	if errors.Is(err, driving.ErrDependenciesRequired) {
		t.Errorf("err = %v, I/O failure must NOT collapse to ErrDependenciesRequired", err)
	}
}

func TestHandleMissingDeps_YesBeatsNoInteractive(t *testing.T) {
	t.Parallel()
	// --yes + --no-interactive → autoInstall arm. CLI mutual-exclusion
	// check (ErrConflictingModeFlags) prevents this combo from reaching
	// the service, but the dispatch's flag truth-table still resolves
	// it cleanly to the safer "yes pre-confirms" semantics.
	conf := &fakeConfirmer{}
	svc := newDispatchService(t, conf)
	missing := []domain.ServiceName{mustNewServiceName(t, "ghost-service")}
	req := driving.AddServiceRequest{
		BaseDir:       "/proj",
		ServiceName:   mustNewServiceName(t, "postgres"),
		Yes:           true,
		NoInteractive: true,
	}

	err := svc.HandleMissingDependenciesForTest(context.Background(), req, missing)

	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("err = %v, want wrap of ErrServiceUnsupported (--yes wins over --no-interactive)", err)
	}
	if errors.Is(err, driving.ErrDependenciesRequired) {
		t.Errorf("err = %v, must NOT fail-fast when --yes is set", err)
	}
	if len(conf.addDepCalls) != 0 {
		t.Errorf("ConfirmAddDependency called %d times, want 0", len(conf.addDepCalls))
	}
}
