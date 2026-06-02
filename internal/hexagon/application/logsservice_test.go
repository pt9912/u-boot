package application_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/application"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// logsFixture bundles the LogsService under test with the two
// fake adapters it uses (fs + engine). The fs starts seeded with
// `u-boot.yaml` and `compose.yaml` so most tests can override one
// or the other; tests that exercise the missing-file paths call
// `removeFile` after the fixture builds.
type logsFixture struct {
	svc    *application.LogsService
	fs     *fakeFS
	engine *fakeDockerEngine
}

func newLogsFixture(t *testing.T) *logsFixture {
	t.Helper()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	if err := fs.WriteFile("/proj/compose.yaml",
		[]byte("services:\n  postgres:\n    image: postgres:16-alpine\n"), 0o644); err != nil {
		t.Fatalf("seed compose.yaml: %v", err)
	}
	engine := newFakeDockerEngine()
	svc := application.NewLogsService(fs, engine, nil)
	return &logsFixture{svc: svc, fs: fs, engine: engine}
}

func TestLogsService_BaseDirEmpty_NonSentinelError(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	_, err := fix.svc.Logs(context.Background(), driving.LogsRequest{})
	if err == nil {
		t.Fatalf("expected error for empty BaseDir, got nil")
	}
	// No sentinel — the CLI would map this to the default Exit-1
	// path (programmer error, not user error).
	if errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("empty BaseDir leaked into ErrProjectNotInitialized: %v", err)
	}
}

func TestLogsService_NoUbootYaml_ReturnsErrProjectNotInitialized(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	engine := newFakeDockerEngine()
	svc := application.NewLogsService(fs, engine, nil)
	_, err := svc.Logs(context.Background(), driving.LogsRequest{BaseDir: "/proj"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrProjectNotInitialized) {
		t.Errorf("err = %v, want wrap of ErrProjectNotInitialized", err)
	}
	if fix := engine.logsCallCount; fix != 0 {
		t.Errorf("ComposeLogs called %d times despite missing u-boot.yaml, want 0", fix)
	}
}

func TestLogsService_NoComposeYaml_ReturnsErrComposeFileMissing(t *testing.T) {
	t.Parallel()
	fs := newFakeFS()
	if err := fs.WriteFile("/proj/u-boot.yaml",
		[]byte("schemaVersion: 1\nproject:\n  name: demo\n"), 0o644); err != nil {
		t.Fatalf("seed u-boot.yaml: %v", err)
	}
	engine := newFakeDockerEngine()
	svc := application.NewLogsService(fs, engine, nil)
	_, err := svc.Logs(context.Background(), driving.LogsRequest{BaseDir: "/proj"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driving.ErrComposeFileMissing) {
		t.Errorf("err = %v, want wrap of ErrComposeFileMissing", err)
	}
	if fix := engine.logsCallCount; fix != 0 {
		t.Errorf("ComposeLogs called %d times despite missing compose.yaml, want 0", fix)
	}
}

// TestLogsService_HappyPath_TailNormalised pins the T0-(c)-contract:
// empty Tail in the Request → `"all"` in the adapter options. The
// fake stores `logsOptions`; we assert on the post-call state.
func TestLogsService_HappyPath_TailNormalised(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	fix.engine.scriptLogs(nil)
	var sink bytes.Buffer
	_, err := fix.svc.Logs(context.Background(), driving.LogsRequest{
		BaseDir:    "/proj",
		Service:    "", // empty → no Services slice (Compose-Default)
		Follow:     false,
		Tail:       "",
		OutputSink: &sink,
	})
	if err != nil {
		t.Fatalf("Logs: %v", err)
	}
	if fix.engine.logsOptions.Tail != "all" {
		t.Errorf("Tail = %q, want %q (T0-(c) normalisation)", fix.engine.logsOptions.Tail, "all")
	}
	if len(fix.engine.logsOptions.Services) != 0 {
		t.Errorf("Services = %v, want empty (T0-(a) Compose-Default)", fix.engine.logsOptions.Services)
	}
	if fix.engine.logsOptions.Sink == nil {
		t.Errorf("Sink is nil; want OutputSink to be forwarded")
	}
	if fix.engine.logsOptions.Follow {
		t.Errorf("Follow = true, want false")
	}
}

// TestLogsService_HappyPath_ServiceFilter pins the single-service
// path: Request.Service = "postgres" → ComposeLogsOptions.Services
// = ["postgres"].
func TestLogsService_HappyPath_ServiceFilter(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	fix.engine.scriptLogs(nil)
	if _, err := fix.svc.Logs(context.Background(), driving.LogsRequest{
		BaseDir: "/proj",
		Service: "postgres",
		Follow:  true,
		Tail:    "100",
	}); err != nil {
		t.Fatalf("Logs: %v", err)
	}
	got := fix.engine.logsOptions
	if len(got.Services) != 1 || got.Services[0] != "postgres" {
		t.Errorf("Services = %v, want [\"postgres\"]", got.Services)
	}
	if !got.Follow {
		t.Errorf("Follow = false, want true")
	}
	if got.Tail != "100" {
		t.Errorf("Tail = %q, want %q", got.Tail, "100")
	}
}

// TestLogsService_SIGINT_ContextCanceled_ReturnsNil pins the
// slice-v1-logs §SIGINT-Vertrag Schicht 2: when the adapter
// returns a `context.Canceled`-wrapped error, the use case
// short-circuits to (LogsResponse{}, nil) so the CLI exits 0.
func TestLogsService_SIGINT_ContextCanceled_ReturnsNil(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	// Adapter signals Ctrl-C via ctx.Err()-unverdeckt-Return.
	fix.engine.scriptLogs(context.Canceled)
	resp, err := fix.svc.Logs(context.Background(), driving.LogsRequest{
		BaseDir: "/proj",
		Follow:  true,
	})
	if err != nil {
		t.Errorf("err = %v, want nil (SIGINT-Vertrag Schicht 2)", err)
	}
	if resp != (driving.LogsResponse{}) {
		t.Errorf("resp = %+v, want zero value", resp)
	}
}

// TestLogsService_SIGINT_DeadlineExceeded_ReturnsNil pins the
// symmetric path for `context.DeadlineExceeded` — same short-
// circuit so test-timeouts and RPC deadlines don't surface as
// errors.
func TestLogsService_SIGINT_DeadlineExceeded_ReturnsNil(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	fix.engine.scriptLogs(context.DeadlineExceeded)
	if _, err := fix.svc.Logs(context.Background(), driving.LogsRequest{
		BaseDir: "/proj",
	}); err != nil {
		t.Errorf("err = %v, want nil (DeadlineExceeded same as Canceled)", err)
	}
}

// TestLogsService_ComposeRuntimeError_PropagatesSentinel pins the
// 12-exit-code path: a runtime failure from Compose (unknown
// service, exit ≠ 0) propagates with `driven.ErrComposeRuntime`
// intact — the CLI relies on `errors.Is` to map it to Exit-12.
func TestLogsService_ComposeRuntimeError_PropagatesSentinel(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	composeErr := fmt.Errorf("docker compose logs failed: %w", driven.ErrComposeRuntime)
	fix.engine.scriptLogs(composeErr)
	_, err := fix.svc.Logs(context.Background(), driving.LogsRequest{
		BaseDir: "/proj",
		Service: "psotgres", // typo: at runtime Compose says unknown
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("err = %v, want wrap of ErrComposeRuntime (Exit-12)", err)
	}
}

// TestLogsService_DockerUnavailable_PropagatesSentinel pins the
// 11-exit-code path: a docker-environment failure from the
// adapter propagates with `driven.ErrDockerUnavailable` intact.
func TestLogsService_DockerUnavailable_PropagatesSentinel(t *testing.T) {
	t.Parallel()
	fix := newLogsFixture(t)
	composeErr := fmt.Errorf("docker daemon unreachable: %w", driven.ErrDockerUnavailable)
	fix.engine.scriptLogs(composeErr)
	_, err := fix.svc.Logs(context.Background(), driving.LogsRequest{BaseDir: "/proj"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("err = %v, want wrap of ErrDockerUnavailable (Exit-11)", err)
	}
}
