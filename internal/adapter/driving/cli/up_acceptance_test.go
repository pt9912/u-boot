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
	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// upUseCaseStub returns a fixed response/err and ignores the request.
// `lastReq` captures the Request the CLI assembled so individual
// tests can pin the SilenceProgress / Timeout / BaseDir flag
// propagation (slice-v1-cli-json-dry-run-up-down T6).
type upUseCaseStub struct {
	resp    driving.UpResponse
	err     error
	called  bool
	lastReq driving.UpRequest
}

func (s *upUseCaseStub) Up(_ context.Context, req driving.UpRequest) (driving.UpResponse, error) {
	s.called = true
	s.lastReq = req
	return s.resp, s.err
}

// newAppWithUpStub wires the up use-case stub plus a deterministic
// getwd so tests do not depend on the runner's CWD. Mirrors
// `newAppWithAddStub` from add_json_test.go.
func newAppWithUpStub(stub driving.UpUseCase) *cli.App {
	return newAppWithUp(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-up-test/demo", nil }))
}

// unmarshalUpEnv parses the JSON envelope for further structural
// pins beyond the AssertMinimalEnvelope shape.
func unmarshalUpEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}

// ----------------------------------------------------------------------
// Happy-Path / Carrier
// ----------------------------------------------------------------------

// TestUpJSON_HappyPath_EmitsMinimalDataEnvelope pins the T0-(g) data
// carrier shape: name/state/port/healthcheck per service.
func TestUpJSON_HappyPath_EmitsMinimalDataEnvelope(t *testing.T) {
	stub := &upUseCaseStub{
		resp: driving.UpResponse{
			Result: domain.UpResult{
				Stabilized: true,
				Services: []domain.ServiceStatus{
					{Name: "postgres", ContainerStatus: domain.StateRunning, Port: "5432:5432", Healthcheck: "healthy"},
				},
			},
		},
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalUpEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if data == nil {
		t.Fatalf("expected data carrier, got nil")
	}
	services, _ := data["services"].([]any)
	if len(services) != 1 {
		t.Fatalf("services length: want 1, got %d", len(services))
	}
	svc, _ := services[0].(map[string]any)
	if svc["name"] != "postgres" {
		t.Errorf("name: want postgres, got %v", svc["name"])
	}
	if svc["port"] != "5432:5432" {
		t.Errorf("port: want 5432:5432, got %v", svc["port"])
	}
	if svc["healthcheck"] != "healthy" {
		t.Errorf("healthcheck: want healthy, got %v", svc["healthcheck"])
	}
}

// TestUpJSON_FireAndForget_HasTimeoutFireAndForgetMarker pins T0-(j):
// `--timeout=0` produces services: [] (NOT null) plus the optional
// marker `data.timeoutFireAndForget: true`.
func TestUpJSON_FireAndForget_HasTimeoutFireAndForgetMarker(t *testing.T) {
	stub := &upUseCaseStub{resp: driving.UpResponse{Result: domain.UpResult{Stabilized: false, Services: nil}}}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "up", "--timeout=0"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalUpEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if data == nil {
		t.Fatalf("expected data carrier on fire-and-forget, got nil")
	}
	marker, present := data["timeoutFireAndForget"]
	if !present {
		t.Fatalf("data.timeoutFireAndForget MUST be present in fire-and-forget mode")
	}
	if got, _ := marker.(bool); !got {
		t.Errorf("data.timeoutFireAndForget: want true, got %v", marker)
	}
	// Empty-Array-Pin (R5-LOW-3): JSON-Layer MUSS [] serialisieren,
	// nicht null. RawMessage-Check.
	pinServicesIsEmptyArrayNotNull(t, stdout.Bytes())
}

// TestUpJSON_HappyPath_TimeoutFireAndForgetMarkerAbsent pins that the
// marker MUST be absent in non-fire-and-forget mode (Key-Absence-
// Disambiguation, R4 marker discipline).
func TestUpJSON_HappyPath_TimeoutFireAndForgetMarkerAbsent(t *testing.T) {
	stub := &upUseCaseStub{
		resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{}}},
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	env := unmarshalUpEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if _, present := data["timeoutFireAndForget"]; present {
		t.Errorf("data.timeoutFireAndForget MUST be absent outside fire-and-forget; got %v", data["timeoutFireAndForget"])
	}
}

// pinServicesIsEmptyArrayNotNull verifies that the `services` field
// serializes as `[]`, not `null` — Empty-Array-Pin (T0-(j) R5-LOW-3
// + R6-LOW-2). RawMessage-Check: byte sequence `"services":[]` MUST
// appear; `"services":null` is a regression.
func pinServicesIsEmptyArrayNotNull(t *testing.T, raw []byte) {
	t.Helper()
	if bytes.Contains(raw, []byte(`"services":null`)) {
		t.Errorf("R5-LOW-3 regression: services field MUST serialize as [], got `null`; raw=%s", raw)
	}
	if !bytes.Contains(raw, []byte(`"services":[]`)) {
		t.Errorf("Empty-Array-Pin: services MUST serialize as []; raw=%s", raw)
	}
}

// TestUpJSON_AllStable_DiagnosticsIsEmptyArrayNotNull pins T0-(j) for
// the diagnostics field — empty MUST serialize as [], not null.
func TestUpJSON_AllStable_DiagnosticsIsEmptyArrayNotNull(t *testing.T) {
	stub := &upUseCaseStub{
		resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{}}},
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	raw := stdout.Bytes()
	if bytes.Contains(raw, []byte(`"diagnostics":null`)) {
		t.Errorf("diagnostics MUST serialize as [], got `null`; raw=%s", raw)
	}
	if !bytes.Contains(raw, []byte(`"diagnostics":[]`)) {
		t.Errorf("Empty-Array-Pin: diagnostics MUST serialize as []; raw=%s", raw)
	}
}

// ----------------------------------------------------------------------
// Flag-Propagation
// ----------------------------------------------------------------------

// TestUpJSON_QuietJSON_StillEmitsEnvelope pins T0-(b): `--quiet --json`
// is semantically identical to `--json` — quiet does NOT suppress the
// envelope (Cluster-T0-(a) doctor-Pattern).
func TestUpJSON_QuietJSON_StillEmitsEnvelope(t *testing.T) {
	stub := &upUseCaseStub{
		resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{}}},
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--quiet", "--json", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatalf("--quiet --json MUST emit envelope on stdout (T0-(b)); got empty")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(0),
	)
}

// TestUpJSON_JSONQuiet_ReversedFlagOrder pins symmetry: `--json --quiet`
// produces the same envelope as `--quiet --json`.
func TestUpJSON_JSONQuiet_ReversedFlagOrder(t *testing.T) {
	stub := &upUseCaseStub{
		resp: driving.UpResponse{Result: domain.UpResult{Stabilized: true, Services: []domain.ServiceStatus{}}},
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "--quiet", "up"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(0),
	)
}

// TestUpJSON_SilenceProgress_TrueWhenJSON pins T0-(c) Form (d): the
// CLI sets req.SilenceProgress = flags.JSON, so the use case sees
// `true` in --json mode and `false` otherwise.
func TestUpJSON_SilenceProgress_TrueWhenJSON(t *testing.T) {
	stub := &upUseCaseStub{}
	app := newAppWithUpStub(stub)
	if err := app.Execute(context.Background(), []string{"--json", "up"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute --json: %v", err)
	}
	if !stub.lastReq.SilenceProgress {
		t.Errorf("--json MUST set req.SilenceProgress=true; got false")
	}
	// Contrast: without --json, SilenceProgress MUST be false.
	stub2 := &upUseCaseStub{}
	app2 := newAppWithUpStub(stub2)
	if err := app2.Execute(context.Background(), []string{"up"}, &bytes.Buffer{}, &bytes.Buffer{}); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if stub2.lastReq.SilenceProgress {
		t.Errorf("without --json, req.SilenceProgress MUST be false; got true")
	}
}

// ----------------------------------------------------------------------
// Error-Pfade
// ----------------------------------------------------------------------

// TestUpJSON_InvalidTimeout_EmitsCLI006Envelope pins the Pre-UC-
// Validation-Pfad: ErrInvalidTimeout fällt durch reportError mit
// data=nil interface (R2-MED-4).
func TestUpJSON_InvalidTimeout_EmitsCLI006Envelope(t *testing.T) {
	stub := &upUseCaseStub{}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up", "--timeout=-1"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrInvalidTimeout) {
		t.Fatalf("expected ErrInvalidTimeout, got %v", err)
	}
	if cli.ExitCode(err) != 2 {
		t.Errorf("exit code: want 2, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	if stub.called {
		t.Errorf("use case called despite Pre-UC-Validation failure")
	}
}

// TestUpJSON_ProjectNotInitialized_LHFAINIT001_Exit10 pins the Cross-
// Slice-Klassen-Pin (R4-MED-2): up als Environment-Operation mappt
// ErrProjectNotInitialized auf LH-FA-INIT-001 (Pattern-Erbe generate),
// NICHT auf LH-FA-ADD-001 (add/remove als Service-Operations).
func TestUpJSON_ProjectNotInitialized_LHFAINIT001_Exit10(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: %q absent: %w", "/tmp/u-boot-up-test/demo/u-boot.yaml", driving.ErrProjectNotInitialized),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("expected ErrProjectNotInitialized, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-001"),
	)
}

// TestUpJSON_ErrUpFileSystem_LHNFAREL003_Exit14 pins the FS-Sentinel-
// Pfad (Row 1 der Mapper-Tabelle).
func TestUpJSON_ErrUpFileSystem_LHNFAREL003_Exit14(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: ReadFile(%q): %w: %w",
			"/tmp/u-boot-up-test/demo/compose.yaml", driving.ErrUpFileSystem, errors.New("disk read error")),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrUpFileSystem) {
		t.Fatalf("expected ErrUpFileSystem, got %v", err)
	}
	if cli.ExitCode(err) != 14 {
		t.Errorf("exit code: want 14, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestUpJSON_ComposeRuntime_LHNFAREL003_Exit12 pins Row 3 der Mapper-
// Tabelle (T7-MED-2): mapComposeRuntimeSentinel matched
// driven.ErrComposeRuntime auf LH-NFA-REL-003, Exit-Code 12 vom
// ExitCode-Helper.
func TestUpJSON_ComposeRuntime_LHNFAREL003_Exit12(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: ComposeUp on %q: %w",
			"/tmp/u-boot-up-test/demo", driven.ErrComposeRuntime),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Fatalf("expected ErrComposeRuntime, got %v", err)
	}
	if cli.ExitCode(err) != 12 {
		t.Errorf("exit code: want 12, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(12),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestUpJSON_StabilizationTimeout_LHFAUP001_Exit12 pins Row 4 (T7-MED-3):
// up-spezifischer Runtime-Sentinel ErrStabilizationTimeout → LH-FA-UP-001,
// Exit-Code 12 via ExitCode-Helper (cli.go:296).
func TestUpJSON_StabilizationTimeout_LHFAUP001_Exit12(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: stabilization timeout after 60s: %w", driving.ErrStabilizationTimeout),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrStabilizationTimeout) {
		t.Fatalf("expected ErrStabilizationTimeout, got %v", err)
	}
	if cli.ExitCode(err) != 12 {
		t.Errorf("exit code: want 12, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(12),
		jsontestutil.WithExpectedCodes("LH-FA-UP-001"),
	)
}

// TestUpJSON_ComposeFileMissing_LHFAUP001_Exit10 pins Row 6 der
// Mapper-Tabelle (T7-LOW-3): geteilte fachliche Validierung (auch in
// mapDown). Einen Pin in up genügt; down nutzt den selben Code.
func TestUpJSON_ComposeFileMissing_LHFAUP001_Exit10(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: %q absent: %w",
			"/tmp/u-boot-up-test/demo/compose.yaml", driving.ErrComposeFileMissing),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrComposeFileMissing) {
		t.Fatalf("expected ErrComposeFileMissing, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-UP-001"),
	)
}

// TestUpJSON_UnknownError_DefaultsToLHFACLI006_Exit1 pins Row 10 der
// Mapper-Tabelle (T7-LOW-4): ein ungemappter Error fällt auf den
// Default-Branch LH-FA-CLI-006 / Exit 1.
func TestUpJSON_UnknownError_DefaultsToLHFACLI006_Exit1(t *testing.T) {
	stub := &upUseCaseStub{
		err: errors.New("up service: synthetic unknown error"),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error")
	}
	if cli.ExitCode(err) != 1 {
		t.Errorf("exit code: want 1 (default), got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(1),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
}

// TestUpJSON_MultiWrap_FSAndStabilizationTimeout_FSFirst_ByDesign is
// T7-LOW-5 / T7-MED-1 by-design verification: zusätzlicher Defense-
// Pin für ein FS+Driving-Sentinel-Pair. Konstruierter Multi-`%w`-Wrap
// mit ErrUpFileSystem + ErrStabilizationTimeout demonstriert die
// Zwei-Pfad-Disambiguation:
//   - Mapper FS-first → code = LH-NFA-REL-003 (FS-Klasse via Row 1)
//   - ExitCode-Helper checkt driving.ErrStabilizationTimeout VOR
//     isFilesystemError (cli.go:296) → exit = 12 (StabilizationTimeout
//     Sub-Klasse).
//
// Beide Pfade unabhängig; Konsumenten disambiguieren über
// (code, exitCode)-Tupel per T8 §6.7-Doku-Pin. Selbe Logik wie beim
// FS+Docker-Wrap (TestUpJSON_MultiWrap_FSAndDocker).
func TestUpJSON_MultiWrap_FSAndStabilizationTimeout_FSFirst_ByDesign(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: synthetic chain: %w: %w",
			driving.ErrUpFileSystem, driving.ErrStabilizationTimeout),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrUpFileSystem) {
		t.Errorf("errors.Is(ErrUpFileSystem) MUST match; got %v", err)
	}
	if !errors.Is(err, driving.ErrStabilizationTimeout) {
		t.Errorf("errors.Is(ErrStabilizationTimeout) MUST match; got %v", err)
	}
	// ExitCode-Helper StabilizationTimeout-first (cli.go:296) → 12.
	if cli.ExitCode(err) != 12 {
		t.Errorf("ExitCode: want 12 (StabilizationTimeout-Sub-Klasse via cli.go:296), got %d", cli.ExitCode(err))
	}
	// Mapper FS-first → diagnostics[0].code = LH-NFA-REL-003.
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(12),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestUpJSON_DockerUnavailable_LHNFAREL003_Exit11 pins the shared
// helper (Row 2 der Mapper-Tabelle): mapComposeRuntimeSentinel
// matched ErrDockerUnavailable auf LH-NFA-REL-003, Exit-Code 11 vom
// ExitCode-Helper.
func TestUpJSON_DockerUnavailable_LHNFAREL003_Exit11(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: ComposeUp on %q: %w",
			"/tmp/u-boot-up-test/demo", driven.ErrDockerUnavailable),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Fatalf("expected ErrDockerUnavailable, got %v", err)
	}
	if cli.ExitCode(err) != 11 {
		t.Errorf("exit code: want 11, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExitCode(11),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestUpJSON_MultiWrap_FSAndDocker_SwitchOrderFSFirst is the R2-HIGH-2
// + R3-HIGH-1 Defense-only-Pin: a synthetic multi-`%w` chain with
// BOTH ErrUpFileSystem AND ErrDockerUnavailable MUST resolve to
// LH-NFA-REL-003/Exit 14 (FS class), NOT Exit 11 (Docker class).
//
// Heute existiert kein realer Code-Pfad der beide Sentinels chained
// (readComposeFile failed VOR ComposeUp). Der Pin verifiziert die
// Mapper-/ExitCode-Robustheit gegen einen synthetisch konstruierten
// Multi-Wrap.
func TestUpJSON_MultiWrap_FSAndDocker_SwitchOrderFSFirst(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: synthetic chain: %w: %w",
			driving.ErrUpFileSystem, driven.ErrDockerUnavailable),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr)
	// errors.Is matched beide Sentinels via Go 1.20+ multi-`%w`-Chain.
	if !errors.Is(err, driving.ErrUpFileSystem) {
		t.Errorf("errors.Is(ErrUpFileSystem) MUST match in multi-wrap; got %v", err)
	}
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("errors.Is(ErrDockerUnavailable) MUST match in multi-wrap; got %v", err)
	}
	// T7-MED-1 by-design: Mapper ist FS-first → diagnostics[0].code
	// = LH-NFA-REL-003 (FS-Klasse). ExitCode-Helper (cli.go:290)
	// checked driven-Sentinels first → exitCode = 11
	// (Docker-Sub-Klasse). Beide Pfade getrennt; (code, exitCode)-
	// Tupel-Disambiguation per §6.7 ist der Vertrag — der FS-Code
	// signalisiert die Klasse, der Exit differenziert die Sub-
	// Sentinel-Quelle. Pin verifiziert das Zwei-Pfad-Verhalten.
	if code := cli.ExitCode(err); code != 11 {
		t.Errorf("ExitCode: want 11 (Docker-Sub-Klasse via cli.go:290 Driven-first), got %d", code)
	}
	// Mapper-Tabelle ist FS-first: diagnostics[0].code MUSS
	// LH-NFA-REL-003 (FS-Klasse) sein.
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("up"),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestUpJSON_FSError_SanitizesPathInDiagnosticMessage pins R2-MED-5:
// runUp wrappt UC-Errors mit sanitizeBaseDir vor reportError. Der
// absolute Pfad `/tmp/u-boot-up-test/demo/compose.yaml` wird
// project-relative `compose.yaml`.
//
//nolint:dupl // Sanitizer-Pin-Pattern bewusst symmetrisch zu down (per-Subcommand-Pfad-Leak-Defense).
func TestUpJSON_FSError_SanitizesPathInDiagnosticMessage(t *testing.T) {
	stub := &upUseCaseStub{
		err: fmt.Errorf("up service: ReadFile(%q): %w: %w",
			"/tmp/u-boot-up-test/demo/compose.yaml",
			driving.ErrUpFileSystem, errors.New("disk read error")),
	}
	app := newAppWithUpStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "up"}, &stdout, &stderr); err == nil {
		t.Fatal("expected error")
	}
	env := unmarshalUpEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	if !strings.Contains(msg, "compose.yaml") {
		t.Errorf("sanitized message MUST contain project-relative path; got: %q", msg)
	}
	if strings.Contains(msg, "/tmp/u-boot-up-test/demo") {
		t.Errorf("R2-MED-5 path-leak: absolute BaseDir MUST NOT appear in diagnostic.message; got: %q", msg)
	}
}
