package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/adapter/driving/cli/jsontestutil"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// downUseCaseStub returns a fixed response/err and ignores the
// request. `lastReq` captures the Request the CLI assembled so
// individual tests can pin the SilenceConfirmer / RemoveVolumes /
// AssumeYes flag propagation (slice-v1-cli-json-dry-run-up-down T6).
type downUseCaseStub struct {
	resp    driving.DownResponse
	err     error
	called  bool
	lastReq driving.DownRequest
}

func (s *downUseCaseStub) Down(_ context.Context, req driving.DownRequest) (driving.DownResponse, error) {
	s.called = true
	s.lastReq = req
	return s.resp, s.err
}

// newAppWithDownStub wires the down use-case stub plus a deterministic
// getwd so tests do not depend on the runner's CWD.
func newAppWithDownStub(stub driving.DownUseCase) *cli.App {
	return newAppWithDown(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-down-test/demo", nil }))
}

// unmarshalDownEnv parses the JSON envelope for further structural
// pins beyond the AssertMinimalEnvelope shape.
func unmarshalDownEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}

// ----------------------------------------------------------------------
// Happy-Path / Idempotenz
// ----------------------------------------------------------------------

// TestDownJSON_HappyPath_EmitsMinimalDataEnvelope pins the T0-(h) data
// carrier shape: removedVolumes bool (false on the regular stop path).
func TestDownJSON_HappyPath_EmitsMinimalDataEnvelope(t *testing.T) {
	stub := &downUseCaseStub{resp: driving.DownResponse{RemovedVolumes: false}}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalDownEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if data == nil {
		t.Fatalf("expected data carrier, got nil")
	}
	// R5-MED-1 + R6-MED-1: removedVolumes MUST be present (no
	// omitempty) — false is a legitimate success value.
	rv, present := data["removedVolumes"]
	if !present {
		t.Fatalf("data.removedVolumes MUST be present (no omitempty); got absent")
	}
	if got, _ := rv.(bool); got {
		t.Errorf("data.removedVolumes: want false (no --volumes), got %v", rv)
	}
}

// TestDownJSON_WithVolumesAndYes_RemovedVolumesTrue pins the
// destructive success path: --volumes --yes proceed without prompt,
// resp.RemovedVolumes=true mirrors to data.removedVolumes=true.
func TestDownJSON_WithVolumesAndYes_RemovedVolumesTrue(t *testing.T) {
	stub := &downUseCaseStub{resp: driving.DownResponse{RemovedVolumes: true}}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--yes", "--json", "down", "--volumes"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalDownEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if got, _ := data["removedVolumes"].(bool); !got {
		t.Errorf("data.removedVolumes: want true, got %v", data["removedVolumes"])
	}
	if !stub.lastReq.RemoveVolumes {
		t.Errorf("req.RemoveVolumes: want true, got false")
	}
	if !stub.lastReq.AssumeYes {
		t.Errorf("req.AssumeYes: want true, got false")
	}
}

// ----------------------------------------------------------------------
// Confirmer-Branch
// ----------------------------------------------------------------------

// TestDownJSON_SilenceConfirmer_TrueWhenJSON pins T0-(d) Option (b):
// CLI sets req.SilenceConfirmer = flags.JSON. Contrast pin: without
// --json the field is false.
func TestDownJSON_SilenceConfirmer_TrueWhenJSON(t *testing.T) {
	stub := &downUseCaseStub{}
	app := newAppWithDownStub(stub)
	if err := app.Execute(context.Background(), []string{"--yes", "--json", "down"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute --json: %v", err)
	}
	if !stub.lastReq.SilenceConfirmer {
		t.Errorf("--json MUST set req.SilenceConfirmer=true; got false")
	}
	stub2 := &downUseCaseStub{}
	app2 := newAppWithDownStub(stub2)
	if err := app2.Execute(context.Background(), []string{"down"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if stub2.lastReq.SilenceConfirmer {
		t.Errorf("without --json, req.SilenceConfirmer MUST be false; got true")
	}
}

// ----------------------------------------------------------------------
// Error-Pfade
// ----------------------------------------------------------------------

// TestDownJSON_ConflictingModeFlags_EmitsCLI005AEnvelope pins the
// LH-FA-CLI-005A §235 Pre-UC-Validation: --yes + --no-interactive →
// Exit 2 with code LH-FA-CLI-005A (R2-MED-2 Mode-Mutex-Pattern).
func TestDownJSON_ConflictingModeFlags_EmitsCLI005AEnvelope(t *testing.T) {
	stub := &downUseCaseStub{}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--yes", "--no-interactive", "--json", "down"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Fatalf("expected ErrConflictingModeFlags, got %v", err)
	}
	if cli.ExitCode(err) != 2 {
		t.Errorf("exit code: want 2, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-005A"),
	)
	if stub.called {
		t.Errorf("use case called despite Pre-UC-Validation failure")
	}
}

// TestDownJSON_ConfirmationRequired_LHFAINIT005_Exit10 pins the
// down-only Confirmer-Refuse path (Mapper-Tabelle Row 5). LH-Code
// LH-FA-INIT-005 ist geteilt mit init/remove (R3-HIGH-3 Co-Migration).
func TestDownJSON_ConfirmationRequired_LHFAINIT005_Exit10(t *testing.T) {
	stub := &downUseCaseStub{
		err: fmt.Errorf("down service: --volumes declined by user: %w", driving.ErrConfirmationRequired),
	}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "down", "--volumes"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-005"),
	)
}

// TestDownJSON_ProjectNotInitialized_LHFAINIT001_Exit10 pins the
// Cross-Slice-Klassen-Pin (R4-MED-2): down als Environment-Operation
// mappt ErrProjectNotInitialized auf LH-FA-INIT-001 (identisch zu up),
// NICHT auf LH-FA-ADD-001 (Service-Operations).
func TestDownJSON_ProjectNotInitialized_LHFAINIT001_Exit10(t *testing.T) {
	stub := &downUseCaseStub{
		err: fmt.Errorf("down service: %q absent: %w",
			"/tmp/u-boot-down-test/demo/u-boot.yaml", driving.ErrProjectNotInitialized),
	}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("expected ErrProjectNotInitialized, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-001"),
	)
}

// TestDownJSON_ErrDownFileSystem_LHNFAREL003_Exit14 pins Row 1 der
// Mapper-Tabelle für down.
func TestDownJSON_ErrDownFileSystem_LHNFAREL003_Exit14(t *testing.T) {
	stub := &downUseCaseStub{
		err: fmt.Errorf("down service: Exists(%q): %w: %w",
			"/tmp/u-boot-down-test/demo/u-boot.yaml", driving.ErrDownFileSystem, errors.New("permission denied")),
	}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrDownFileSystem) {
		t.Fatalf("expected ErrDownFileSystem, got %v", err)
	}
	if cli.ExitCode(err) != 14 {
		t.Errorf("exit code: want 14, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestDownJSON_QuietJSON_StillEmitsEnvelope pins T0-(b) für down
// (T7-LOW-1 Symmetrie zu up). --quiet --json MUSS Envelope emittieren.
func TestDownJSON_QuietJSON_StillEmitsEnvelope(t *testing.T) {
	stub := &downUseCaseStub{resp: driving.DownResponse{RemovedVolumes: false}}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--quiet", "--json", "down"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatalf("--quiet --json MUST emit envelope on stdout (T0-(b)); got empty")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(0),
	)
}

// TestDownJSON_DockerUnavailable_LHNFAREL003_Exit11 pins Row 2 der
// Mapper-Tabelle für down (T7-LOW-2 Symmetrie zu up). Selber shared
// helper mapComposeRuntimeSentinel.
func TestDownJSON_DockerUnavailable_LHNFAREL003_Exit11(t *testing.T) {
	stub := &downUseCaseStub{
		err: fmt.Errorf("down service: ComposeDown on %q: %w",
			"/tmp/u-boot-down-test/demo", driven.ErrDockerUnavailable),
	}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr)
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Fatalf("expected ErrDockerUnavailable, got %v", err)
	}
	if cli.ExitCode(err) != 11 {
		t.Errorf("exit code: want 11, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("down"),
		jsontestutil.WithExitCode(11),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestDownJSON_FSError_SanitizesPathInDiagnosticMessage pins R2-MED-5
// für down (analog up).
//
//nolint:dupl // Sanitizer-Pin-Pattern bewusst symmetrisch zu up (per-Subcommand-Pfad-Leak-Defense).
func TestDownJSON_FSError_SanitizesPathInDiagnosticMessage(t *testing.T) {
	stub := &downUseCaseStub{
		err: fmt.Errorf("down service: Exists(%q): %w: %w",
			"/tmp/u-boot-down-test/demo/u-boot.yaml",
			driving.ErrDownFileSystem, errors.New("permission denied")),
	}
	app := newAppWithDownStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "down"}, &stdout, &stderr); err == nil {
		t.Fatal("expected error")
	}
	env := unmarshalDownEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	if !strings.Contains(msg, "u-boot.yaml") {
		t.Errorf("sanitized message MUST contain project-relative path; got: %q", msg)
	}
	if strings.Contains(msg, "/tmp/u-boot-down-test/demo") {
		t.Errorf("R2-MED-5 path-leak: absolute BaseDir MUST NOT appear; got: %q", msg)
	}
}
