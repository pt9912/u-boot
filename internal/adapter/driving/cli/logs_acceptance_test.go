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

// newAppWithLogsStub wires the logs use-case stub plus a
// deterministic getwd so tests do not depend on the runner's CWD.
// Mirrors `newAppWithUpStub`/`newAppWithDownStub`.
func newAppWithLogsStub(stub driving.LogsUseCase) *cli.App {
	return newAppWithLogs(stub, cli.WithGetwd(func() (string, error) { return "/tmp/u-boot-logs-test/demo", nil }))
}

// unmarshalLogsEnv parses the JSON envelope for further structural
// pins beyond the AssertMinimalEnvelope shape.
func unmarshalLogsEnv(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v\nraw=%s", err, raw)
	}
	return env
}

// ----------------------------------------------------------------------
// T0-(a) Verbatim + T0-(i) Validation-Reihenfolge
// ----------------------------------------------------------------------

// TestLogsJSON_FollowJSONReject_MapsCLI006Exit2 pins T0-(a) Option (A)
// verbatim: `--follow --json` MUST be rejected with
// ErrFollowJSONNotSupported / Exit 2 BEFORE the use case runs.
func TestLogsJSON_FollowJSONReject_MapsCLI006Exit2(t *testing.T) {
	stub := &fakeLogsUseCase{}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs", "--follow"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrFollowJSONNotSupported) {
		t.Fatalf("expected ErrFollowJSONNotSupported, got %v", err)
	}
	if cli.ExitCode(err) != 2 {
		t.Errorf("exit code: want 2, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	if stub.called {
		t.Errorf("use case called despite --follow --json reject")
	}
}

// TestLogsJSON_ValidationOrder_FollowJSONBeatsInvalidTail pins
// T0-(i) Pre-UC-Validation-Reihenfolge: `--follow --json --tail=-1`
// MUST emit ErrFollowJSONNotSupported FIRST (Row 7), NOT
// ErrInvalidLogsTail (Row 8). Step 1 in runLogs runs before
// Step 2.
func TestLogsJSON_ValidationOrder_FollowJSONBeatsInvalidTail(t *testing.T) {
	stub := &fakeLogsUseCase{}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs", "--follow", "--tail=-1"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrFollowJSONNotSupported) {
		t.Fatalf("validation-order broken: want ErrFollowJSONNotSupported (T0-(i) Step 1), got %v", err)
	}
	if errors.Is(err, cli.ErrInvalidLogsTail) {
		t.Errorf("validation-order broken: ErrInvalidLogsTail leaked through Step 1 reject")
	}
}

// ----------------------------------------------------------------------
// Mapper-Coverage (Rows 1-8 + Default)
// ----------------------------------------------------------------------

// TestLogsJSON_ErrLogsFileSystem_MapsRel003Exit14 pins Row 1
// (FS-first Switch-Order, T0-(e) R3-HIGH-1 Defense).
func TestLogsJSON_ErrLogsFileSystem_MapsRel003Exit14(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: Exists(%q): %w: %w",
			"/tmp/u-boot-logs-test/demo/compose.yaml", driving.ErrLogsFileSystem, errors.New("permission denied")),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrLogsFileSystem) {
		t.Fatalf("expected ErrLogsFileSystem, got %v", err)
	}
	if cli.ExitCode(err) != 14 {
		t.Errorf("exit code: want 14, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(14),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestLogsJSON_DockerUnavailable_MapsRel003Exit11 pins Row 2 via
// shared helper mapComposeRuntimeSentinel.
func TestLogsJSON_DockerUnavailable_MapsRel003Exit11(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: %w", driven.ErrDockerUnavailable),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Fatalf("expected ErrDockerUnavailable, got %v", err)
	}
	if cli.ExitCode(err) != 11 {
		t.Errorf("exit code: want 11, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(11),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestLogsJSON_ComposeRuntime_MapsRel003Exit12 pins Row 3 via
// shared helper.
func TestLogsJSON_ComposeRuntime_MapsRel003Exit12(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: %w", driven.ErrComposeRuntime),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Fatalf("expected ErrComposeRuntime, got %v", err)
	}
	if cli.ExitCode(err) != 12 {
		t.Errorf("exit code: want 12, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(12),
		jsontestutil.WithExpectedCodes("LH-NFA-REL-003"),
	)
}

// TestLogsJSON_ComposeFileMissing_MapsUp001Exit10 pins Row 4
// Cluster-Konsens (T0-(g) — same Sentinel → same LH-Code).
func TestLogsJSON_ComposeFileMissing_MapsUp001Exit10(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: %q absent: %w",
			"/tmp/u-boot-logs-test/demo/compose.yaml", driving.ErrComposeFileMissing),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrComposeFileMissing) {
		t.Fatalf("expected ErrComposeFileMissing, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-UP-001"),
	)
}

// TestLogsJSON_ProjectNotInitialized_MapsInit001Exit10 pins Row 5
// Cross-Slice-Klassen-Pin: Environment-Operation (logs/up/down/
// generate) → LH-FA-INIT-001, NOT LH-FA-ADD-001 (add/remove).
func TestLogsJSON_ProjectNotInitialized_MapsInit001Exit10(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: %q absent: %w",
			"/tmp/u-boot-logs-test/demo/u-boot.yaml", driving.ErrProjectNotInitialized),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Fatalf("expected ErrProjectNotInitialized, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-001"),
	)
}

// TestLogsJSON_InvalidServiceName_MapsInit006Exit10 pins Row 6
// (domain-level format validation).
func TestLogsJSON_InvalidServiceName_MapsInit006Exit10(t *testing.T) {
	stub := &fakeLogsUseCase{}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs", "Postgres"}, &stdout, &stderr)
	if !errors.Is(err, domain.ErrInvalidServiceName) {
		t.Fatalf("expected ErrInvalidServiceName, got %v", err)
	}
	if cli.ExitCode(err) != 10 {
		t.Errorf("exit code: want 10, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(10),
		jsontestutil.WithExpectedCodes("LH-FA-INIT-006"),
	)
	if stub.called {
		t.Errorf("use case called despite invalid service name")
	}
}

// TestLogsJSON_InvalidTail_MapsCLI006Exit2 pins Row 8.
func TestLogsJSON_InvalidTail_MapsCLI006Exit2(t *testing.T) {
	stub := &fakeLogsUseCase{}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs", "--tail=-1"}, &stdout, &stderr)
	if !errors.Is(err, cli.ErrInvalidLogsTail) {
		t.Fatalf("expected ErrInvalidLogsTail, got %v", err)
	}
	if cli.ExitCode(err) != 2 {
		t.Errorf("exit code: want 2, got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(2),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
	if stub.called {
		t.Errorf("use case called despite invalid tail")
	}
}

// TestLogsJSON_UnknownError_MapsCLI006Exit1 pins the Default-Branch
// (Row 9 — unknown error falls through with LH-FA-CLI-006/Exit 1).
func TestLogsJSON_UnknownError_MapsCLI006Exit1(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: errors.New("logs service: synthetic unknown error"),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error")
	}
	if cli.ExitCode(err) != 1 {
		t.Errorf("exit code: want 1 (default), got %d", cli.ExitCode(err))
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(1),
		jsontestutil.WithExpectedCodes("LH-FA-CLI-006"),
	)
}

// ----------------------------------------------------------------------
// Acceptance-Pins
// ----------------------------------------------------------------------

// TestLogsJSON_BoundedTail_EmitsLinesEnvelope pins the happy path:
// `--tail=N --json` buffers the Compose-Stream and emits a
// Minimal+Data envelope with `data.lines []string`.
func TestLogsJSON_BoundedTail_EmitsLinesEnvelope(t *testing.T) {
	stub := &fakeLogsUseCase{
		onLogs: func(req driving.LogsRequest) error {
			_, err := req.OutputSink.Write([]byte("postgres  | 2026-01-01 LOG: ready\nredis     | 2026-01-01 LOG: ready\n"))
			return err
		},
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "logs", "--tail=2"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v (stderr=%s)", err, stderr.String())
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(0),
	)
	env := unmarshalLogsEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	if data == nil {
		t.Fatalf("expected data carrier, got nil")
	}
	lines, _ := data["lines"].([]any)
	if len(lines) != 2 {
		t.Errorf("data.lines length: want 2, got %d (lines=%v)", len(lines), lines)
	}
	if first, _ := lines[0].(string); !strings.Contains(first, "postgres") {
		t.Errorf("data.lines[0] should contain `postgres`, got %q", first)
	}
}

// TestLogsJSON_EmptyServiceSet_LinesIsEmptyArrayNotNull pins the
// Empty-Array contract (T0-(j) Cluster-Pattern from up-down R5-LOW-3):
// when the use case writes nothing to OutputSink, `data.lines` MUST
// serialize as `[]`, not `null`.
func TestLogsJSON_EmptyServiceSet_LinesIsEmptyArrayNotNull(t *testing.T) {
	stub := &fakeLogsUseCase{
		onLogs: func(_ driving.LogsRequest) error { return nil },
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	raw := stdout.Bytes()
	if bytes.Contains(raw, []byte(`"lines":null`)) {
		t.Errorf("lines field MUST serialize as [], got null; raw=%s", raw)
	}
	if !bytes.Contains(raw, []byte(`"lines":[]`)) {
		t.Errorf("Empty-Array-Pin: lines MUST serialize as []; raw=%s", raw)
	}
}

// TestLogsJSON_QuietJSON_StillEmitsEnvelope pins Cluster-T0-(a)
// doctor-Pattern (--quiet --json semantically identical to --json).
func TestLogsJSON_QuietJSON_StillEmitsEnvelope(t *testing.T) {
	stub := &fakeLogsUseCase{
		onLogs: func(req driving.LogsRequest) error {
			_, err := req.OutputSink.Write([]byte("postgres  | LOG: ready\n"))
			return err
		},
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--quiet", "--json", "logs"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if stdout.Len() == 0 {
		t.Fatalf("--quiet --json MUST emit envelope on stdout (Cluster-T0-(a)); got empty")
	}
	jsontestutil.AssertMinimalEnvelope(t, stdout.Bytes(),
		jsontestutil.WithCommand("logs"),
		jsontestutil.WithExitCode(0),
	)
}

// TestLogsJSON_FSError_SanitizesPathInDiagnosticMessage pins the
// Path-Leak-Defense via shared cli/sanitize.go-Helper from up-down T5.
//
//nolint:dupl // Sanitizer-Pin-Pattern bewusst symmetrisch zu up/down (per-Subcommand-Pfad-Leak-Defense).
func TestLogsJSON_FSError_SanitizesPathInDiagnosticMessage(t *testing.T) {
	stub := &fakeLogsUseCase{
		err: fmt.Errorf("logs service: Exists(%q): %w: %w",
			"/tmp/u-boot-logs-test/demo/u-boot.yaml",
			driving.ErrLogsFileSystem, errors.New("permission denied")),
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr); err == nil {
		t.Fatal("expected error")
	}
	env := unmarshalLogsEnv(t, stdout.Bytes())
	diags, _ := env["diagnostics"].([]any)
	if len(diags) != 1 {
		t.Fatalf("want 1 diagnostic, got %d", len(diags))
	}
	diag, _ := diags[0].(map[string]any)
	msg, _ := diag["message"].(string)
	if !strings.Contains(msg, "u-boot.yaml") {
		t.Errorf("sanitized message MUST contain project-relative path; got %q", msg)
	}
	if strings.Contains(msg, "/tmp/u-boot-logs-test/demo") {
		t.Errorf("path-leak: absolute BaseDir MUST NOT appear; got %q", msg)
	}
}

// TestLogsJSON_TrailingNewline_StrippedFromLastLine pins
// splitLogLines: the buffer typically ends with `\n` from the last
// log line, and the trailing empty token MUST be stripped so the
// last data.lines entry is the real last line, not `""`.
func TestLogsJSON_TrailingNewline_StrippedFromLastLine(t *testing.T) {
	stub := &fakeLogsUseCase{
		onLogs: func(req driving.LogsRequest) error {
			_, err := req.OutputSink.Write([]byte("postgres  | A\nredis     | B\n"))
			return err
		},
	}
	app := newAppWithLogsStub(stub)
	var stdout, stderr bytes.Buffer
	if err := app.Execute(context.Background(), []string{"--json", "logs"}, &stdout, &stderr); err != nil {
		t.Fatalf("execute: %v", err)
	}
	env := unmarshalLogsEnv(t, stdout.Bytes())
	data, _ := env["data"].(map[string]any)
	lines, _ := data["lines"].([]any)
	if len(lines) != 2 {
		t.Fatalf("trailing-newline-strip broken: want 2 lines, got %d (lines=%v)", len(lines), lines)
	}
	last, _ := lines[1].(string)
	if last == "" {
		t.Errorf("last line is empty — trailing-newline produced a phantom empty token")
	}
	if !strings.Contains(last, "redis") {
		t.Errorf("last line should be the `redis`-line, got %q", last)
	}
}
