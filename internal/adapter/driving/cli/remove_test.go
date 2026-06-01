package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestRemove_StateTransitionPrintsChangedPaths(t *testing.T) {
	uc := &fakeRemoveServiceUseCase{
		resp: driving.RemoveServiceResponse{
			ServiceName: mustServiceNameCLI(t, "postgres"),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			Changed:     []string{"compose.yaml", ".env.example", "u-boot.yaml"},
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(), []string{"remove", "postgres"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	if !uc.called {
		t.Error("remove use case was not called")
	}
	out := stdout.String()
	if !strings.Contains(out, `Removed service "postgres".`) {
		t.Errorf("stdout missing transition headline; got:\n%s", out)
	}
	for _, p := range []string{"compose.yaml", ".env.example", "u-boot.yaml"} {
		if !strings.Contains(out, "  - "+p) {
			t.Errorf("stdout missing %q in Changed list; got:\n%s", p, out)
		}
	}
}

func TestRemove_IdempotentNoOpPrintsAlreadyDisabled(t *testing.T) {
	uc := &fakeRemoveServiceUseCase{
		resp: driving.RemoveServiceResponse{
			ServiceName: mustServiceNameCLI(t, "postgres"),
			PriorState:  domain.ServiceStateDeactivated,
			State:       domain.ServiceStateDeactivated,
			// Changed=nil
		},
	}
	var stdout, stderr bytes.Buffer
	if err := newAppWithRemove(uc).Execute(context.Background(), []string{"remove", "postgres"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, `Service "postgres" is already disabled; no changes.`) {
		t.Errorf("stdout missing idempotent-no-op message; got:\n%s", out)
	}
	if strings.Contains(out, "Changed:") {
		t.Errorf("stdout shows Changed list on no-op; got:\n%s", out)
	}
}

func TestRemove_PurgeSurfacesManualCleanupHint(t *testing.T) {
	// --purge through the gate → executeRemove ran successfully, but
	// VolumesPurged stays false (T3 defers). CLI summary appends the
	// manual-cleanup NOTE.
	uc := &fakeRemoveServiceUseCase{
		resp: driving.RemoveServiceResponse{
			ServiceName:   mustServiceNameCLI(t, "postgres"),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"compose.yaml", ".env.example", "u-boot.yaml"},
			VolumesPurged: false,
		},
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(), []string{"--yes", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	// The fake captured the request — confirm Purge + Yes propagated.
	if !uc.lastReq.Purge {
		t.Error("Request.Purge = false; --purge flag did not propagate")
	}
	if !uc.lastReq.Yes {
		t.Error("Request.Yes = false; --yes persistent flag did not propagate")
	}
	out := stdout.String()
	if !strings.Contains(out, "--purge was requested") {
		t.Errorf("stdout missing --purge NOTE; got:\n%s", out)
	}
	if !strings.Contains(out, "docker volume rm") {
		t.Errorf("stdout missing manual-cleanup hint; got:\n%s", out)
	}
}

func TestRemove_ConfirmationGateRefusedSurfacesCode10(t *testing.T) {
	// Use case returns ErrConfirmationRequired wrap; ExitCode → 10.
	uc := &fakeRemoveServiceUseCase{
		err: fmt.Errorf("remove service: --purge refused in --no-interactive without --yes: %w",
			driving.ErrConfirmationRequired),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(),
		[]string{"--no-interactive", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute: want error, got nil")
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("ExitCode = %d, want 10", code)
	}
}

func TestRemove_UnregisteredSurfacesCode10(t *testing.T) {
	uc := &fakeRemoveServiceUseCase{
		err: fmt.Errorf("remove: %q was never added: %w", "postgres", driving.ErrServiceUnregistered),
	}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(), []string{"remove", "postgres"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute: want error, got nil")
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("ExitCode = %d, want 10 (ErrServiceUnregistered → validation)", code)
	}
}

func TestRemove_UnknownServiceFailsAtDomainValidation(t *testing.T) {
	// Invalid service name (uppercase letter rejected by
	// domain.NewServiceName) — fails before reaching the use case.
	uc := &fakeRemoveServiceUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(), []string{"remove", "Postgres"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute: want domain.ErrInvalidServiceName, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidServiceName) {
		t.Errorf("err = %v, want wrap of domain.ErrInvalidServiceName", err)
	}
	if uc.called {
		t.Error("use case was called despite domain validation failure")
	}
}

func TestRemove_ConflictingModeFlagsExits2(t *testing.T) {
	uc := &fakeRemoveServiceUseCase{}
	var stdout, stderr bytes.Buffer
	err := newAppWithRemove(uc).Execute(context.Background(),
		[]string{"--yes", "--no-interactive", "remove", "postgres"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute: want ErrConflictingModeFlags, got nil")
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("ExitCode = %d, want 2 (mode-flag mutex is CLI usage error)", code)
	}
}

func TestRemove_HelpListsPurgeFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	uc := &fakeRemoveServiceUseCase{}
	if err := newAppWithRemove(uc).Execute(context.Background(), []string{"remove", "--help"}, &stdout, &stderr); err != nil {
		t.Fatalf("Execute: %v (stderr=%q)", err, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "--purge") {
		t.Errorf("remove --help missing --purge flag mention; got:\n%s", out)
	}
}

func mustServiceNameCLI(t *testing.T, raw string) domain.ServiceName {
	t.Helper()
	name, err := domain.NewServiceName(raw)
	if err != nil {
		t.Fatalf("NewServiceName(%q): %v", raw, err)
	}
	return name
}
