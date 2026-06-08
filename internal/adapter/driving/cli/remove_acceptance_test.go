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
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// removeUseCaseStub returns a fixed response/err and ignores the
// request. Wraps fakeRemoveServiceUseCase with a shape tailored for
// the JSON-pin tests below; every test fills resp / err directly.
//
// `lastReq` captures the Request the CLI assembled so individual
// tests can pin the SilenceConfirmer / PreviewMode / Purge / Yes
// flag propagation (T0-(j) + T0-(h)).
type removeUseCaseStub struct {
	resp    driving.RemoveServiceResponse
	err     error
	called  bool
	lastReq driving.RemoveServiceRequest
}

func (s *removeUseCaseStub) Remove(_ context.Context, req driving.RemoveServiceRequest) (driving.RemoveServiceResponse, error) {
	s.called = true
	s.lastReq = req
	return s.resp, s.err
}

// newAppWithRemoveStub wires the remove use-case stub plus a
// deterministic getwd so tests do not depend on the runner's CWD.
// Mirrors `newAppWithAddStub` from add_json_test.go.
func newAppWithRemoveStub(stub driving.RemoveServiceUseCase) *cli.App {
	return newAppWithRemove(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-remove-test/demo", nil }))
}

// pgRemoveName is a shorthand for the validated postgres ServiceName
// in remove acceptance tests (separate from add's pgServiceName to
// keep the per-file fixtures self-contained).
func pgRemoveName(t *testing.T) domain.ServiceName {
	t.Helper()
	n, err := domain.NewServiceName("postgres")
	if err != nil {
		t.Fatalf("NewServiceName(postgres): %v", err)
	}
	return n
}

// unmarshalRemoveEnv parses the JSON envelope for further structural
// pins beyond the AssertMinimalEnvelope / AssertFullEnvelope shape.
func unmarshalRemoveEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}

// ----------------------------------------------------------------------
// JSON Path A — Minimal + Data Envelope (--json without --dry-run / --diff)
// ----------------------------------------------------------------------

// TestRemoveJSON_BareUsesDataEnvelope is the T0-(f)/(m) success pin:
// `u-boot --json remove postgres` (no preview-flag) ships the
// Minimal+Data envelope with `data: {service, priorState, state,
// volumesPurged}`. Spec §1841: no plannedFiles/changes/dryRun/diff.
func TestRemoveJSON_BareUsesDataEnvelope(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   pgRemoveName(t),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"compose.yaml", ".env.example", "u-boot.yaml"},
			VolumesPurged: false,
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("service", "postgres"),
		jsontestutil.WithDataKeyPresent("priorState", "active"),
		jsontestutil.WithDataKeyPresent("state", "deactivated"),
		jsontestutil.WithDataKeyPresent("volumesPurged", false),
	)
	// SilenceConfirmer-Pin (T0-(j)): --json automatically sets
	// req.SilenceConfirmer so the JSON envelope stays uncorrupted by
	// interactive prompts. The use-case-Stub records lastReq.
	if !stub.lastReq.SilenceConfirmer {
		t.Errorf("SilenceConfirmer: --json must propagate SilenceConfirmer=true; got false (lastReq=%+v)", stub.lastReq)
	}
}

// TestRemoveJSON_BareIdempotentNoOpKeepsDataEnvelope pins the
// T0-(f)/(m) idempotent-no-op path: PriorState==State==Deactivated
// and Changed=nil produces the same data shape but with state/
// priorState both "deactivated" (Konsumenten leiten NoOp aus dem
// Tupel ab).
func TestRemoveJSON_BareIdempotentNoOpKeepsDataEnvelope(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateDeactivated,
			State:       domain.ServiceStateDeactivated,
			// Changed=nil
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("priorState", "deactivated"),
		jsontestutil.WithDataKeyPresent("state", "deactivated"),
		jsontestutil.WithDataKeyPresent("volumesPurged", false),
	)
}

// TestRemoveJSON_EnabledUnsetNormalizesToDeactivated pins the
// EnabledUnset → Deactivated normalisation path (AK-Block: "Normalised
// X (enabled key was missing)") rendered into the data carrier.
func TestRemoveJSON_EnabledUnsetNormalizesToDeactivated(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateEnabledUnset,
			State:       domain.ServiceStateDeactivated,
			Changed:     []string{"u-boot.yaml"},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("priorState", "enabled-unset"),
		jsontestutil.WithDataKeyPresent("state", "deactivated"),
	)
}

// ----------------------------------------------------------------------
// JSON Path B — Dry-Run Voll-Envelope (--dry-run --json)
// ----------------------------------------------------------------------

// TestRemoveJSON_DryRunUsesFullEnvelope pins the --dry-run --json
// path: voll-schema with dryRun=true, diff=false, plannedFiles from
// the recorder (3 captures), data carrier mit allen vier Feldern.
func TestRemoveJSON_DryRunUsesFullEnvelope(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("services:\n  postgres:\n    enabled: true\n"), NewContent: []byte("services:\n  postgres:\n    enabled: false\n")},
				{Path: "compose.yaml", Action: "modify", OldContent: []byte("services:\n  postgres:\n    image: postgres:16\n"), NewContent: []byte("services: {}\n")},
				{Path: ".env.example", Action: "modify", OldContent: []byte("POSTGRES_DB=app\n"), NewContent: []byte("")},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"remove", "postgres", "--dry-run", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("service", "postgres"),
		jsontestutil.WithDataKeyPresent("priorState", "active"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true, envelope=%s", stdout.String())
	}
	if got, _ := env["diff"].(bool); got {
		t.Errorf("diff: want false, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 3 {
		t.Errorf("plannedFiles: want 3, got %d (envelope=%s)", len(pfs), stdout.String())
	}
	// PreviewMode propagation (T0-(b)): --dry-run wires PreviewDryRun.
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode: want PreviewDryRun, got %v", stub.lastReq.PreviewMode)
	}
}

// ----------------------------------------------------------------------
// JSON Path C — Diff Voll-Envelope (--diff --json) plus delete-Action
// ----------------------------------------------------------------------

// TestRemoveJSON_DiffWithoutDryRunRendersHunks pins the --diff --json
// preview-and-apply path: voll-schema with diff=true, dryRun=false,
// hunks for non-binary files.
func TestRemoveJSON_DiffWithoutDryRunRendersHunks(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "modify",
					OldContent: []byte("services:\n  postgres:\n    image: postgres:16\n"),
					NewContent: []byte("services: {}\n")},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"remove", "postgres", "--diff", "--json"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true, envelope=%s", stdout.String())
	}
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 1 {
		t.Fatalf("plannedFiles: want 1, got %d", len(pfs))
	}
	first, _ := pfs[0].(map[string]any)
	hunks, ok := first["hunks"].([]any)
	if !ok || len(hunks) == 0 {
		t.Errorf("plannedFiles[0].hunks: want non-empty, got %v", first["hunks"])
	}
	if stub.lastReq.PreviewMode != driving.PreviewAndApply {
		t.Errorf("PreviewMode: want PreviewAndApply, got %v", stub.lastReq.PreviewMode)
	}
}

// TestRemove_OtelExtraFileDelete_DiffHasDeleteHunk is the T0-(p)
// R4-HIGH-F4 pin: remove is the first end-to-end-visible `delete`-
// Action producer (RemoveAll on otel-extraFiles). The envelope
// MUST carry `plannedFiles[].action: "delete"` with a hunk that
// renders the old content as a removal block.
func TestRemove_OtelExtraFileDelete_DiffHasDeleteHunk(t *testing.T) {
	otelName, err := domain.NewServiceName("opentelemetry")
	if err != nil {
		t.Fatalf("NewServiceName(opentelemetry): %v", err)
	}
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: otelName,
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "modify",
					OldContent: []byte("services:\n  otel:\n    image: otel/opentelemetry-collector:latest\n"),
					NewContent: []byte("services: {}\n")},
				{Path: "otel-collector-config.yaml", Action: "delete",
					OldContent: []byte("receivers:\n  otlp: {}\nexporters:\n  logging: {}\n"),
					NewContent: []byte{}},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"remove", "opentelemetry", "--diff", "--json"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	pfs, _ := env["plannedFiles"].([]any)
	if len(pfs) != 2 {
		t.Fatalf("plannedFiles: want 2, got %d", len(pfs))
	}
	// Find the delete entry — order is recorder-determined, not
	// alphabetic, so iterate.
	var deleteFound bool
	for _, raw := range pfs {
		pf, _ := raw.(map[string]any)
		if pf["action"] == "delete" {
			deleteFound = true
			path, _ := pf["path"].(string)
			if path != "otel-collector-config.yaml" {
				t.Errorf("delete entry path: want otel-collector-config.yaml, got %q", path)
			}
			hunks, _ := pf["hunks"].([]any)
			if len(hunks) == 0 {
				t.Errorf("delete entry MUST carry a hunk rendering the old content; got empty hunks")
			}
		}
	}
	if !deleteFound {
		t.Errorf("envelope missing plannedFiles[].action: \"delete\" (T0-(p)); got: %s", stdout.String())
	}
}

// ----------------------------------------------------------------------
// WARN-Migration (T0-(g)) + --purge --yes --json WarnOnly (T0-(j))
// ----------------------------------------------------------------------

// TestRemove_PurgeYesJSON_WarnOnly is the T0-(j) R1-MED-5 + R5-F2
// pin: `--purge --yes --json` succeeds (Yes skips the gate without
// confirmer call), VolumesPurged stays false (v0.3.0 deferred), the
// Use-Case emits a WARN-Diagnostic with LH-FA-ADD-007. status=warn,
// exit=0 (warn-only does not shift exit code, Spec §447).
func TestRemove_PurgeYesJSON_WarnOnly(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   pgRemoveName(t),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"u-boot.yaml", "compose.yaml", ".env.example"},
			VolumesPurged: false,
			Warnings: []driving.WarningEntry{
				{Code: "LH-FA-ADD-007", Level: "warn",
					Message: "--purge requested for service \"postgres\" but volume removal is deferred (v0.3.0); the named volumes are still on disk and untouched"},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--yes", "--json", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithExpectedCodes("LH-FA-ADD-007"),
		jsontestutil.WithDataKeyPresent("volumesPurged", false),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if status, _ := env["status"].(string); status != "warn" {
		t.Errorf("status: want \"warn\" (Spec §447 warn-only), got %q", status)
	}
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("diagnostics: want 1 WARN entry, got %d (envelope=%s)", len(diags), stdout.String())
	}
	diag, _ := diags[0].(map[string]any)
	if level, _ := diag["level"].(string); level != "warn" {
		t.Errorf("diagnostics[0].level: want \"warn\", got %q", level)
	}
	if !stub.lastReq.Purge || !stub.lastReq.Yes {
		t.Errorf("Purge/Yes not propagated: %+v", stub.lastReq)
	}
}

// TestRemove_PurgeOnVolumelessService_NoWarn is the T0-(g) R3-F1
// pin: services with volumeOptional=true (keycloak, otel) produce
// NO WARN even with --purge (Use-Case-Logic — leeres Warnings-Feld
// returns; CLI emits diagnostics=[]).
func TestRemove_PurgeOnVolumelessService_NoWarn(t *testing.T) {
	kcName, err := domain.NewServiceName("keycloak")
	if err != nil {
		t.Fatalf("NewServiceName(keycloak): %v", err)
	}
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   kcName,
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"u-boot.yaml", "compose.yaml"},
			VolumesPurged: false,
			Warnings:      nil, // Use-Case returns nil for volumeless
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--yes", "--json", "remove", "keycloak", "--purge"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
		jsontestutil.WithDataKeyPresent("volumesPurged", false),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 0 {
		t.Errorf("volumeless --purge MUST NOT emit WARN; got diagnostics=%v", diags)
	}
	if status, _ := env["status"].(string); status != "ok" {
		t.Errorf("status: want \"ok\" (no diagnostics), got %q", status)
	}
}

// ----------------------------------------------------------------------
// Mid-Write-Failure — Variante A (Error dominates, WARN suppressed)
// ----------------------------------------------------------------------

// TestRemove_PurgeYesJSON_MidWriteFailure_ErrorOnly is the T0-(j)
// R4-MED-F3 Variante-A pin: wenn die Execute-Phase mid-write failt
// (ErrRemoveFileSystem) UND der Use-Case auch ein Warnings-Feld
// gesetzt hatte (WARN-Emission VOR Gate, T0-(c)), MUSS Error
// dominieren: nur ein ERROR-Eintrag im diagnostics-Array, kein
// WARN. data-Form: nur `service`, KEIN `volumesPurged`/`priorState`/
// `state` (Zero-Response-Klausel, T0-(f)).
func TestRemove_PurgeYesJSON_MidWriteFailure_ErrorOnly(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			// Warnings populated — would normally surface, but Variante A
			// suppresses them on the error path.
			Warnings: []driving.WarningEntry{
				{Code: "LH-FA-ADD-007", Level: "warn", Message: "would-be deferred"},
			},
			PlannedFiles: []driving.PlannedFile{
				{Path: "u-boot.yaml", Action: "modify", OldContent: []byte("a\n"), NewContent: []byte("b\n")},
				{Path: "compose.yaml", Action: "modify", OldContent: []byte("svc\n"), NewContent: []byte("")},
			},
		},
		err: fmt.Errorf("remove write compose.yaml: %w: %w",
			driving.ErrRemoveFileSystem, errors.New("disk full")),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--yes", "--json", "remove", "postgres", "--purge", "--diff"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrRemoveFileSystem) {
		t.Fatalf("expected ErrRemoveFileSystem to propagate; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit code: want 14, got %d", code)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
		jsontestutil.WithDataKeyPresent("service", "postgres"),
		jsontestutil.WithDataKeyAbsent("volumesPurged", "priorState", "state"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	// Variante A: exactly ONE entry, and it's an ERROR.
	if len(diags) != 1 {
		t.Fatalf("Variante A: want 1 diagnostic (error only, WARN suppressed); got %d: %s", len(diags), stdout.String())
	}
	diag, _ := diags[0].(map[string]any)
	if level, _ := diag["level"].(string); level != "error" {
		t.Errorf("diagnostics[0].level: want \"error\", got %q", level)
	}
	if status, _ := env["status"].(string); status != "error" {
		t.Errorf("status: want \"error\", got %q", status)
	}
}

// TestRemove_ReadPathFSFailure_BeforeStateDetect is the R2-F3 pin:
// FS-Failure auf dem Read-Pfad (Exists / ReadFile / Lstat) während
// detectServiceState() landet ebenfalls als ErrRemoveFileSystem im
// Envelope — exit 14, LH-NFA-REL-003. Es gibt KEINE PlannedFiles
// (Failure passiert vor dem ersten Write).
func TestRemove_ReadPathFSFailure_BeforeStateDetect(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: read compose.yaml: %w: %w",
			driving.ErrRemoveFileSystem, errors.New("permission denied")),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrRemoveFileSystem) {
		t.Fatalf("expected ErrRemoveFileSystem; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit code: want 14, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
		jsontestutil.WithDataKeyPresent("service", "postgres"),
	)
}

// ----------------------------------------------------------------------
// ConfirmationRequired-Pfade × 3 Varianten
// ----------------------------------------------------------------------

// TestRemoveJSON_ConfirmationRequired_NoInteractive pins the gate-
// refused path in JSON-mode: `--no-interactive --json --purge` without
// `--yes` → ErrConfirmationRequired → LH-FA-INIT-005, exit 10.
func TestRemoveJSON_ConfirmationRequired_NoInteractive(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove service: --purge refused in --no-interactive without --yes: %w",
			driving.ErrConfirmationRequired),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--no-interactive", "--json", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-005"),
		jsontestutil.WithDataKeyPresent("service", "postgres"),
	)
}

// TestRemoveJSON_ConfirmationRequired_NoYesNoFlag_DefaultsToRefuse
// pins the JSON-mode behaviour-change (T0-(j)): even without
// --no-interactive, `--purge --json` ohne `--yes` MUST refuse — the
// CLI sets SilenceConfirmer=true so the use case sees a noopConfirmer
// returning `false, nil` → ErrConfirmationRequired.
func TestRemoveJSON_ConfirmationRequired_NoYesNoFlag_DefaultsToRefuse(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove service: --purge confirmation refused: %w",
			driving.ErrConfirmationRequired),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrConfirmationRequired) {
		t.Fatalf("expected ErrConfirmationRequired (JSON-mode behaviour-change); got %v", err)
	}
	if !stub.lastReq.SilenceConfirmer {
		t.Errorf("SilenceConfirmer: --json must propagate true even without --yes; got %+v", stub.lastReq)
	}
}

// TestRemove_ConfirmerUnavailableJSON_RoutesToCLI005A is the R2-F1
// pin: ErrConfirmerUnavailable (I/O-Failure im Confirmer-Call)
// routes to LH-FA-CLI-005A / exit 10, distinct from
// ErrConfirmationRequired's LH-FA-INIT-005.
func TestRemove_ConfirmerUnavailableJSON_RoutesToCLI005A(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: confirmer: %w: %w",
			driving.ErrConfirmerUnavailable, errors.New("stdin EOF")),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrConfirmerUnavailable) {
		t.Fatalf("expected ErrConfirmerUnavailable; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10 (Confirmation-Gate-Klasse), got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-005A"),
	)
}

// TestRemove_MultiWWrap_SwitchOrderDefense_FSWins is the R3-F3
// defense pin: a synthetically-constructed multi-`%w`-Wrap chaining
// both ErrConfirmerUnavailable AND ErrConfirmationRequired MUST
// route to LH-FA-CLI-005A (Infrastruktur-Klasse), NICHT
// LH-FA-INIT-005 (Confirmation-Gate-Klasse) — Switch-Order T0-(e)
// pins Infrastruktur-First. Heute existiert KEIN realer Code-Pfad
// der beide Sentinels chained — der Pin verifiziert die Mapper-
// Robustheit gegen synthetische Konstrukte.
func TestRemove_MultiWWrap_SwitchOrderDefense_FSWins(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("synthetic chain: %w: %w",
			driving.ErrConfirmerUnavailable, driving.ErrConfirmationRequired),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres", "--purge"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error to propagate")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		// Switch-Order: ErrConfirmerUnavailable check fires FIRST →
		// LH-FA-CLI-005A wins over the LH-FA-INIT-005 inner-sentinel.
		jsontestutil.WithExpectedCodes("LH-FA-CLI-005A"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	if code, _ := diag["code"].(string); code == "LH-FA-INIT-005" {
		t.Errorf("Switch-Order violation: routed to LH-FA-INIT-005 instead of LH-FA-CLI-005A (multi-%%w synthetic chain)")
	}
}

// ----------------------------------------------------------------------
// Service-Sentinels × 4 — ServiceUnsupported / Unregistered / Inconsistent / ProjectNotInitialized
// ----------------------------------------------------------------------

// TestRemoveJSON_ServiceUnsupported_LH002 pins LH-FA-ADD-002 (unknown
// catalogue entry) routing.
func TestRemoveJSON_ServiceUnsupported_LH002(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: service \"redis\" not in catalogue: %w",
			driving.ErrServiceUnsupported),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrServiceUnsupported) {
		t.Fatalf("expected ErrServiceUnsupported; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-ADD-002"),
	)
}

// TestRemove_ServiceUnregisteredJSON_ErrorLevelCodePin is the R6-F3
// + R5-F2 pin: `ErrServiceUnregistered` routes to `LH-FA-ADD-007`
// with `level: "error"` (Symmetrie zur WARN-Pfad-Nutzung desselben
// Codes in TestRemove_PurgeYesJSON_WarnOnly). Konsumenten
// disambiguieren WARN-Pfad und ERROR-Pfad ausschließlich über
// `(code, level)`-Tupel.
func TestRemove_ServiceUnregisteredJSON_ErrorLevelCodePin(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: service \"postgres\" was never added: %w",
			driving.ErrServiceUnregistered),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrServiceUnregistered) {
		t.Fatalf("expected ErrServiceUnregistered; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-ADD-007"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	if level, _ := diag["level"].(string); level != "error" {
		t.Errorf("ERROR-Pfad: diagnostics[0].level: want \"error\", got %q (symmetry to WARN-Pfad pin)", level)
	}
	if status, _ := env["status"].(string); status != "error" {
		t.Errorf("status: want \"error\", got %q", status)
	}
}

// TestRemoveJSON_ServiceInconsistent_LH005 pins LH-FA-ADD-005 routing
// for ErrServiceInconsistent (managed-block orphan / yaml-patch fail).
// Triple-Use-Marker (R5-F3): the same Code covers three sub-paths in
// the application service (managedblock-malformed, scanner default,
// yaml.PatchScalar) — alle drei landen identisch im Envelope.
func TestRemoveJSON_ServiceInconsistent_LH005(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: managed block malformed: %w",
			driving.ErrServiceInconsistent),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrServiceInconsistent) {
		t.Fatalf("expected ErrServiceInconsistent; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-ADD-005"),
	)
}

// TestRemoveJSON_ProjectNotInitialized_LH001 pins LH-FA-ADD-001 for
// the `no u-boot.yaml`-Pfad.
func TestRemoveJSON_ProjectNotInitialized_LH001(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove: project not initialized: %w",
			driving.ErrProjectNotInitialized),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("expected ErrProjectNotInitialized; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-ADD-001"),
	)
}

// ----------------------------------------------------------------------
// CLI-Form-Sentinels — ConflictingModeFlags + NoPositionalArg
// ----------------------------------------------------------------------

// TestRemoveJSON_ConflictingModeFlags_EmitsCLI005A is the R1-F4 pin:
// `--yes --no-interactive --json remove postgres` (mode-flag mutex)
// → LH-FA-CLI-005A / exit 2.
func TestRemoveJSON_ConflictingModeFlags_EmitsCLI005A(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--yes", "--no-interactive", "--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrConflictingModeFlags) {
		t.Fatalf("expected ErrConflictingModeFlags; got %v", err)
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("exit code: want 2, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-005A"),
	)
	if stub.called {
		t.Errorf("use case called despite mode-flag mutex; got %+v", stub.lastReq)
	}
}

// TestRemove_NoPositionalArg_JSON_EmitsCLI006Envelope is the
// R11-HIGH-F1 + R12-MED-F2 Custom-Args-Validator pin: `u-boot --json
// remove` (kein positional arg) emittiert den Minimal-Envelope mit
// LH-FA-CLI-006 / Exit 2 AUF STDOUT (NICHT nur Cobra-Standard-
// Error-Print), weil der Custom-Validator den Envelope VOR der
// Sentinel-Propagation an Cobra emittiert. Symmetrie-Bruch zu
// `remove "bad name" --json` (voller Envelope mit LH-FA-INIT-006/
// Exit 10) ist damit aufgehoben — beide Pre-UC-Validation-Pfade
// emittieren Envelope auf stdout.
func TestRemove_NoPositionalArg_JSON_EmitsCLI006Envelope(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrServiceNameMissing) {
		t.Fatalf("expected ErrServiceNameMissing; got %v", err)
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("exit code: want 2 (LH-FA-CLI-006), got %d", code)
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected JSON envelope on stdout (Spec §1841 missing-arg pin); got empty")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	if stub.called {
		t.Errorf("use case called despite missing positional arg")
	}
	// Service-Kontext fehlt → kein data-Carrier.
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if _, present := env["data"]; present {
		t.Errorf("missing-arg envelope MUST NOT carry `data` (no service context); got data=%v", env["data"])
	}
}

// TestRemoveJSON_InvalidServiceName_EmitsINIT006 is the symmetric
// pin to TestRemove_NoPositionalArg_JSON: `remove "bad name" --json`
// MUSS ebenfalls den vollen Envelope auf stdout emittieren, aber
// mit LH-FA-INIT-006 / Exit 10 (Domain-Validation, not CLI-Form).
// Beide Pre-UC-Validation-Pfade gehen über `reportError`.
func TestRemoveJSON_InvalidServiceName_EmitsINIT006(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "INVALID NAME WITH SPACES"}, &stdout, &stderr)
	if !errors.Is(err, domain.ErrInvalidServiceName) {
		t.Fatalf("expected ErrInvalidServiceName; got %v", err)
	}
	if code := cli.ExitCode(err); code != 10 {
		t.Errorf("exit code: want 10 (Domain-Validation), got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-006"),
	)
}

// ----------------------------------------------------------------------
// Human-Mode tests — stderr-Trennung, --diff WARNING-Routing
// ----------------------------------------------------------------------

// TestRemove_PurgeHumanDiff_StderrSeparation is the R2-F6 pin: bei
// `--purge --diff` ohne `--json` läuft der Diff-Renderer auf stdout
// (unified-diff body), WÄHREND die `--purge`-deferred-Volumes-WARNING
// auf errOut bleibt — der Diff-Body wird NIE durch die WARNING-Prose
// polluted. Stderr keeps its today's prose form for backwards
// compatibility (existing tooling that grep'd `WARNING:` on stderr).
func TestRemove_PurgeHumanDiff_StderrSeparation(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   pgRemoveName(t),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"compose.yaml"},
			VolumesPurged: false,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "modify",
					OldContent: []byte("services:\n  postgres:\n    image: postgres:16\n"),
					NewContent: []byte("services: {}\n")},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--yes", "remove", "postgres", "--purge", "--diff"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	out := stdout.String()
	errs := stderr.String()
	// Diff body MUST be on stdout — sentinel header is enough.
	if !strings.Contains(out, "--- compose.yaml (modify)") {
		t.Errorf("--diff body missing from stdout; got:\n%s", out)
	}
	// WARNING prose MUST be on stderr, NOT stdout.
	if strings.Contains(out, "WARNING:") {
		t.Errorf("WARNING leaked onto stdout (polluted diff body); got:\n%s", out)
	}
	if !strings.Contains(errs, "WARNING: --purge was requested") {
		t.Errorf("WARNING missing from stderr; got:\n%s", errs)
	}
	if !strings.Contains(errs, "docker volume rm") {
		t.Errorf("manual-cleanup hint missing from stderr; got:\n%s", errs)
	}
}

// TestRemove_DryRun_PropagatesPreviewDryRunFlag is the T0-(b) Truth-
// Table-Pin für `--dry-run` ohne `--diff` ohne `--json`: human-mode,
// PreviewMode=PreviewDryRun, kein Production-Write. Die Tatsache, dass
// PreviewDryRun an die Use-Case weitergereicht wird, ist die wichtige
// Form-Drift-Defense — alles andere wäre eine Regression im
// previewModeFromFlags-Helper.
func TestRemove_DryRun_PropagatesPreviewDryRunFlag(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			Changed:     []string{"compose.yaml"},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"remove", "postgres", "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode: want PreviewDryRun, got %v", stub.lastReq.PreviewMode)
	}
	// Human mode, no --json → no JSON envelope on stdout.
	if strings.HasPrefix(strings.TrimSpace(stdout.String()), "{") {
		t.Errorf("non-JSON mode emitted JSON on stdout; got:\n%s", stdout.String())
	}
}

// ----------------------------------------------------------------------
// T7 Pre-T8-Review-Fixes (R13-HIGH-1 + R13-MED-1 + R14-HIGH-2 + R14-MED-1)
// ----------------------------------------------------------------------

// TestRemove_NoPositionalArg_DryRunJSON_EmitsFullSchemaEnvelope
// (R13-HIGH-1): `--dry-run --json remove` ohne positional arg MUSS
// das Voll-Schema-Envelope emittieren (Spec §1842), NICHT das
// Minimal-Schema. Pre-T7 hatte der Args-Validator hardcoded
// `dryRun=false, diff=false` und produzierte einen Minimal-Envelope
// trotz `--dry-run` → Spec §1842 Verletzung. Fix in T7: Validator
// liest `cmd.Flags().GetBool("dry-run"/"diff")` und reicht den
// User-Flag-State an `writeErrorEnvelopeSub` durch.
func TestRemove_NoPositionalArg_DryRunJSON_EmitsFullSchemaEnvelope(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--dry-run", "--json", "remove"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrServiceNameMissing) {
		t.Fatalf("expected ErrServiceNameMissing; got %v", err)
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("exit code: want 2, got %d", code)
	}
	// Voll-Schema MUSS dryRun/diff/plannedFiles/changes haben.
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true (--dry-run flag was set), got false; envelope=%s", stdout.String())
	}
}

// TestRemove_TooManyArgs_JSON_EmitsCLI006Envelope (R13-MED-1):
// `--json remove a b c` (zwei extra positional args) MUSS ebenfalls
// einen Envelope auf stdout produzieren — Symmetrie zum missing-arg-
// Pfad. Pre-T7 ist der `len(args)>1`-Pfad nur durch `cobra.ExactArgs(1)`
// abgefangen worden ohne stdout-Envelope (Spec §1841 Verletzung).
// Fix in T7: Validator emittiert Envelope vor dem Cobra-Error-Return.
func TestRemove_TooManyArgs_JSON_EmitsCLI006Envelope(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres", "extra-arg"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error (too many positional args)")
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("exit code: want 2 (LH-FA-CLI-006 usage class), got %d", code)
	}
	if stdout.Len() == 0 {
		t.Fatalf("expected JSON envelope on stdout (Spec §1841 too-many-args symmetry to missing-arg); got empty")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	if stub.called {
		t.Errorf("use case called despite too-many-args validation failure")
	}
}

// TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution (R14-
// HIGH-2): bei `--purge --dry-run` ohne `--json` skipped die Use-Case
// den Confirmer-Gate (Plan T0-(h)(a)) und führt keine Mutation
// aus → `VolumesPurged: false`, aber semantisch "nichts versucht",
// nicht "deferred work". Pre-T7 zeigte `printRemoveSummary` die
// WARNING-Prosa trotzdem (Logic `--purge && !VolumesPurged`) und
// suggestierte fälschlich, dass Volume-Cleanup deferred wurde. Fix
// in T7: WARNING wird in `PreviewDryRun` unterdrückt; nur das
// Summary-Headline und die "Would change"-Liste landen auf stdout
// (Headline-Wording bleibt heute "Removed service" — separate
// "Would remove"-Wording-Drift wäre eigene Sub-Decision).
func TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   pgRemoveName(t),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"compose.yaml", ".env.example", "u-boot.yaml"},
			VolumesPurged: false,
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--yes", "remove", "postgres", "--purge", "--dry-run"}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	// WARNING-Prosa darf in Dry-Run NICHT auf stderr landen — der
	// Use-Case-Gate ist geskippt und nichts wurde versucht.
	errs := stderr.String()
	if strings.Contains(errs, "WARNING:") {
		t.Errorf("Dry-Run-Pollution: WARNING-Prosa darf in PreviewDryRun NICHT auf stderr landen; got:\n%s", errs)
	}
	if strings.Contains(errs, "docker volume rm") {
		t.Errorf("Dry-Run-Pollution: manual-cleanup-Hint darf in PreviewDryRun NICHT erscheinen; got:\n%s", errs)
	}
	// PreviewMode korrekt durchgereicht.
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode: want PreviewDryRun, got %v", stub.lastReq.PreviewMode)
	}
	// PreviewAndApply-Konstrast-Pin: bei `--purge` OHNE `--dry-run`
	// muss die WARNING wieder erscheinen (sonst hätten wir einen
	// falschen Suppression-Fix).
	var stdout2, stderr2 bytes.Buffer
	stub2 := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName:   pgRemoveName(t),
			PriorState:    domain.ServiceStateActive,
			State:         domain.ServiceStateDeactivated,
			Changed:       []string{"compose.yaml"},
			VolumesPurged: false,
		},
	}
	app2 := newAppWithRemoveStub(stub2)
	if err := app2.Execute(context.Background(),
		[]string{"--yes", "remove", "postgres", "--purge"}, &stdout2, &stderr2); err != nil {
		t.Fatalf("execute (PreviewAndApply contrast): %v", err)
	}
	if !strings.Contains(stderr2.String(), "WARNING:") {
		t.Errorf("PreviewAndApply contrast: WARNING MUST appear when actually-applied; got stderr:\n%s", stderr2.String())
	}
}

// TestRemove_FSErrorWithAbsolutePath_SanitizesMessage (R14-MED-1):
// FS-Wraps der Form `fmt.Errorf("remove write %s: %w: %w", absPath,
// ErrRemoveFileSystem, raw)` enthalten den absoluten Filesystem-
// Pfad. Pre-T7 wäre dieser Pfad 1:1 in `diagnostic.message` gelandet
// (Info-Leak des Filesystem-Layouts an JSON-Konsumenten). Fix in T7:
// `runRemove` wrappt `removeErr` mit `sanitizeBaseDir(removeErr, cwd)`
// vor dem `reportError`-Call. Der Sanitizer ersetzt `<baseDir>/foo`
// durch `foo` und bare `<baseDir>` durch `.`. errors.Is bleibt intakt
// via Unwrap-Chain.
//
// Test-Fixture: `newAppWithRemoveStub` setzt cwd auf
// `/tmp/u-boot-remove-test/demo`. Der konstruierte Fehler enthält
// `/tmp/u-boot-remove-test/demo/compose.yaml` — der Sanitizer MUSS
// das auf `compose.yaml` reduzieren.
func TestRemove_FSErrorWithAbsolutePath_SanitizesMessage(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf("remove write /tmp/u-boot-remove-test/demo/compose.yaml: %w: %w",
			driving.ErrRemoveFileSystem, errors.New("disk full")),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres"}, &stdout, &stderr)
	// errors.Is MUSS trotz Sanitizer-Wrap funktionieren — Unwrap-
	// Chain bleibt intakt.
	if !errors.Is(err, driving.ErrRemoveFileSystem) {
		t.Fatalf("errors.Is(ErrRemoveFileSystem) broken by sanitizer; got %v", err)
	}
	if code := cli.ExitCode(err); code != 14 {
		t.Errorf("exit code: want 14, got %d", code)
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	// Sanitisierte Form: Pfad ist project-relative.
	if !strings.Contains(msg, "compose.yaml") {
		t.Errorf("sanitized message MUST contain project-relative path \"compose.yaml\"; got: %q", msg)
	}
	// Path-Leak-Defense: absoluter BaseDir-Prefix MUSS weg sein.
	if strings.Contains(msg, "/tmp/u-boot-remove-test/demo") {
		t.Errorf("R14-MED-1 path-leak: absolute BaseDir MUST NOT appear in diagnostic.message; got: %q", msg)
	}
}

// TestRemove_FSErrorWithBaseDirSubstring_NotMangled (R15-LOW-1):
// Sanitizer-Robustheit gegen Substring-Kollisionen. Pre-T8 nutzte der
// Sanitizer `strings.ReplaceAll(msg, baseDir, ".")` — naive Form, die
// einen Error-Pfad `<baseDir>-cache/<…>` zu `.-cache/<…>` mangeln
// würde. Fix in T8: `replaceBareBaseDir` ersetzt baseDir nur an
// Word-Boundaries (gefolgt von End-of-String oder Nicht-Pfad-Byte).
//
// Test-Fixture-Trick: `newAppWithRemoveStub` setzt cwd auf
// `/tmp/u-boot-remove-test/demo`. Der konstruierte Fehler enthält
// `/tmp/u-boot-remove-test/demo-cache/lock` — Substring-Prefix von
// cwd, aber NICHT cwd selbst. Robust-Sanitizer MUSS den Pfad
// unangetastet lassen; daneben aber `/tmp/u-boot-remove-test/demo/
// compose.yaml` korrekt zu `compose.yaml` reduzieren.
func TestRemove_FSErrorWithBaseDirSubstring_NotMangled(t *testing.T) {
	stub := &removeUseCaseStub{
		err: fmt.Errorf(
			"remove write /tmp/u-boot-remove-test/demo/compose.yaml: "+
				"unrelated lock held at /tmp/u-boot-remove-test/demo-cache/lock: %w: %w",
			driving.ErrRemoveFileSystem, errors.New("disk full")),
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--json", "remove", "postgres"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrRemoveFileSystem) {
		t.Fatalf("errors.Is(ErrRemoveFileSystem) broken; got %v", err)
	}
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	// Project-relative Pfad MUSS erscheinen (baseDir+sep-Pass).
	if !strings.Contains(msg, "compose.yaml") {
		t.Errorf("sanitized message MUST contain project-relative path %q; got: %q", "compose.yaml", msg)
	}
	// Substring-Kollision MUSS NICHT mangled werden. `-cache/lock`
	// nach dem nicht-cwd-Pfad-Prefix darf NICHT zu `.-cache/lock`
	// werden — der Defense-Check pinnt das ungemangelte Original.
	if strings.Contains(msg, ".-cache/lock") {
		t.Errorf("R15-LOW-1 substring-collision: baseDir-prefixed sibling path MUST NOT be mangled to %q; got: %q", ".-cache/lock", msg)
	}
	if !strings.Contains(msg, "/tmp/u-boot-remove-test/demo-cache/lock") {
		t.Errorf("R15-LOW-1 substring-collision: sibling path MUST stay intact (no boundary-less ReplaceAll); got: %q", msg)
	}
	// Bare baseDir-Form: kein cwd-Pfad als nackter Suffix mehr im
	// Pre-Pass-Output (nur die Substring-Kollisions-Form).
	if strings.Contains(msg, "/tmp/u-boot-remove-test/demo/compose.yaml") {
		t.Errorf("R14-MED-1 regression: cwd-prefixed path MUST be relativized; got: %q", msg)
	}
}

// TestRemove_TooManyArgs_DryRunJSON_EmitsFullSchemaEnvelope (R15-LOW-2):
// Coverage-Pin für die Flag-Awareness-Pfad-Kombination von R13-MED-1
// (too-many-args-Envelope-Symmetrie) UND R13-HIGH-1 (Voll-Schema bei
// `--dry-run`). `--json --dry-run remove a b c` MUSS einen Voll-Schema-
// Envelope auf stdout produzieren — Symmetrie zum NoPositionalArg-
// DryRun-Pin und zum TooManyArgs-Minimal-Pin. Future-Regression-Defense:
// wenn jemand den `len(args)>1`-Pfad ohne Flag-Read refaktoriert, würde
// dieser Pin als Erstes fallen.
func TestRemove_TooManyArgs_DryRunJSON_EmitsFullSchemaEnvelope(t *testing.T) {
	stub := &removeUseCaseStub{}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(),
		[]string{"--dry-run", "--json", "remove", "postgres", "extra-arg"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error (too many positional args)")
	}
	if code := cli.ExitCode(err); code != 2 {
		t.Errorf("exit code: want 2 (LH-FA-CLI-006 usage class), got %d", code)
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true (--dry-run flag set), got false; envelope=%s", stdout.String())
	}
	if stub.called {
		t.Errorf("use case called despite too-many-args validation failure")
	}
}

// TestRemove_DryRunDiffJSONCombo pins the 3-flag combo: `--dry-run
// --diff --json` produces voll-schema with BOTH dryRun=true AND
// diff=true, hunks populated, no production write (previewModeFromFlags
// (true, true) → PreviewDryRun per init T0-(b) Wahrheitstabelle).
func TestRemove_DryRunDiffJSONCombo(t *testing.T) {
	stub := &removeUseCaseStub{
		resp: driving.RemoveServiceResponse{
			ServiceName: pgRemoveName(t),
			PriorState:  domain.ServiceStateActive,
			State:       domain.ServiceStateDeactivated,
			PlannedFiles: []driving.PlannedFile{
				{Path: "compose.yaml", Action: "modify",
					OldContent: []byte("services:\n  postgres: {}\n"),
					NewContent: []byte("services: {}\n")},
			},
		},
	}
	app := newAppWithRemoveStub(stub)

	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(),
		[]string{"remove", "postgres", "--dry-run", "--diff", "--json"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertFullEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("remove"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalRemoveEnv(t, stdout.Bytes())
	if got, _ := env["dryRun"].(bool); !got {
		t.Errorf("dryRun: want true (3-flag combo)")
	}
	if got, _ := env["diff"].(bool); !got {
		t.Errorf("diff: want true (3-flag combo)")
	}
	// PreviewDryRun (--dry-run wins on the truth-table).
	if stub.lastReq.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode: want PreviewDryRun (3-flag combo, --dry-run wins), got %v", stub.lastReq.PreviewMode)
	}
}
