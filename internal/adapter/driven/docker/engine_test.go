package docker_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driven/docker"
	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

func TestEngine_SatisfiesDockerEnginePort(t *testing.T) {
	t.Parallel()
	// Why: pin that *docker.Engine continues to implement the
	// driven.DockerEngine interface. A method-signature drift on
	// either side would break the wiring in cmd/uboot/main.go (T7);
	// catching it here makes the failure point obvious.
	var _ driven.DockerEngine = docker.NewEngine()
	var _ driven.DockerEngine = docker.WithEngineBinary("/bin/echo")
}

func TestEngine_ComposeUp_MissingBinary_ReturnsErrDockerUnavailable(t *testing.T) {
	t.Parallel()
	e := docker.WithEngineBinary("/does/not/exist/docker-binary")
	_, err := e.ComposeUp(context.Background(), "/tmp/demo", driven.ComposeUpOptions{Detach: true})
	assertErrIs(t, err, driven.ErrDockerUnavailable)
}

func TestEngine_ComposeDown_MissingBinary_ReturnsErrDockerUnavailable(t *testing.T) {
	t.Parallel()
	e := docker.WithEngineBinary("/does/not/exist/docker-binary")
	err := e.ComposeDown(context.Background(), "/tmp/demo", driven.ComposeDownOptions{RemoveVolumes: true})
	assertErrIs(t, err, driven.ErrDockerUnavailable)
}

func TestEngine_ComposePs_MissingBinary_ReturnsErrDockerUnavailable(t *testing.T) {
	t.Parallel()
	e := docker.WithEngineBinary("/does/not/exist/docker-binary")
	_, err := e.ComposePs(context.Background(), "/tmp/demo")
	assertErrIs(t, err, driven.ErrDockerUnavailable)
}

func TestEngine_MissingBinary_DoesNotReturnErrComposeRuntime(t *testing.T) {
	t.Parallel()
	// Why: pin the M6 slice's §Sentinel-Schichtung 11-vs-12 split.
	// A missing binary is an env failure (code 11), not a runtime
	// failure (code 12). If the preflight ever started classifying
	// LookPath failures as ErrComposeRuntime, the CLI would
	// silently move this case to exit code 12.
	e := docker.WithEngineBinary("/does/not/exist/docker-binary")
	_, err := e.ComposeUp(context.Background(), "/tmp/demo", driven.ComposeUpOptions{})
	if errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("ComposeUp with missing binary returned ErrComposeRuntime; want ErrDockerUnavailable only")
	}
}

// shellBinaryAvailable skips the test when the requested
// always-present POSIX utility is missing. CI runs inside
// `golang:1.26.3` (Debian) and `golangci/golangci-lint:v2.12.2-alpine`,
// both of which ship `/bin/echo`, `/bin/true`, and `/bin/false`.
// The skip keeps the tests portable for developer environments
// that strip them.
func shellBinaryAvailable(t *testing.T, path string) {
	t.Helper()
	if _, err := exec.LookPath(path); err != nil {
		t.Skipf("%s not available: %v", path, err)
	}
}

func TestEngine_Preflight_InfoFails_ReturnsErrDockerUnavailable(t *testing.T) {
	t.Parallel()
	// Why: cover the second preflight branch (probe.Info failure)
	// without a real docker daemon. `/bin/false` returns exit code
	// 1 for every invocation, so LookPath succeeds (binary exists)
	// but the daemon-roundtrip emulation fails. Pins the
	// errors.Is wrap: a non-zero exit on the second probe stays in
	// the ErrDockerUnavailable class, not ErrComposeRuntime.
	shellBinaryAvailable(t, "/bin/false")
	e := docker.WithEngineBinary("/bin/false")
	_, err := e.ComposePs(context.Background(), "/tmp/demo")
	assertErrIs(t, err, driven.ErrDockerUnavailable)
	if errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("/bin/false preflight failure leaked into ErrComposeRuntime: %v", err)
	}
}

func TestEngine_AllProbesPass_ButPsOutputNotJSON_ReturnsErrComposeRuntime(t *testing.T) {
	t.Parallel()
	// Why: cover the post-preflight ComposePs path. `/bin/echo`
	// returns exit code 0 for every invocation (LookPath ok,
	// version ok, compose-version ok → preflight green) and then
	// prints its own args as stdout — which is not valid JSON, so
	// parseComposePsOutput rejects it. End result: ComposePs
	// returns ErrComposeRuntime (NOT ErrDockerUnavailable).
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	_, err := e.ComposePs(context.Background(), "/tmp/demo")
	if err == nil {
		t.Fatal("expected error from non-JSON ps output")
	}
	if errors.Is(err, driven.ErrDockerUnavailable) {
		t.Errorf("preflight should have passed with /bin/echo, but got ErrDockerUnavailable: %v", err)
	}
	assertErrIs(t, err, driven.ErrComposeRuntime)
}

func TestEngine_ComposeDown_AllProbesPass_Succeeds(t *testing.T) {
	t.Parallel()
	// Why: cover the post-preflight ComposeDown success path.
	// `/bin/echo` exits 0 for the down call (no JSON parsing
	// downstream) so the use case returns nil. Pins the "down has
	// no snapshot" contract — no ComposePs follow-up.
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	err := e.ComposeDown(context.Background(), "/tmp/demo", driven.ComposeDownOptions{RemoveVolumes: true})
	if err != nil {
		t.Errorf("ComposeDown with passing preflight + 0-exit binary: %v", err)
	}
}

func TestEngine_ComposeUp_AllProbesPass_NoPostUpPsRoundtrip(t *testing.T) {
	t.Parallel()
	// Why: pin the LH-FA-UP-001 §970 fire-and-forget contract at
	// the adapter level. A successful `compose up` MUST NOT
	// follow up with a `compose ps` roundtrip — if it did, the
	// extra call could surface ErrComposeRuntime after the `up`
	// itself already succeeded, leaking past the §970 guarantee.
	//
	// `/bin/echo` returns 0 for `up`; if the adapter had a
	// follow-up ps it would receive non-JSON stdout and surface
	// ErrComposeRuntime. The pin asserts the call returns nil.
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	_, err := e.ComposeUp(context.Background(), "/tmp/demo", driven.ComposeUpOptions{Detach: true})
	if err != nil {
		t.Errorf("ComposeUp must not perform a follow-up ps roundtrip; got: %v", err)
	}
}

func TestParseComposePsOutput_Empty(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  []byte
	}{
		{"nil", nil},
		{"empty", []byte("")},
		{"only-whitespace", []byte("  \n\t  ")},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := docker.ParseComposePsOutputForTest(tc.raw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != nil {
				t.Errorf("expected nil services for empty input, got %v", got)
			}
		})
	}
}

func TestParseComposePsOutput_NDJSON_SingleService(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"Service":"postgres","Name":"demo-postgres-1","State":"running","Health":"healthy","Publishers":[{"PublishedPort":5432,"TargetPort":5432,"Protocol":"tcp"}]}`)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}
	want := driven.ComposeService{
		Name:   "postgres",
		State:  "running",
		Health: "healthy",
		Ports:  []string{"5432:5432"},
	}
	assertComposeServiceEqual(t, got[0], want)
}

func TestParseComposePsOutput_NDJSON_MultipleServices(t *testing.T) {
	t.Parallel()
	raw := []byte(
		`{"Service":"postgres","Name":"demo-postgres-1","State":"running","Health":"healthy","Publishers":[{"PublishedPort":5432,"TargetPort":5432}]}` + "\n" +
			`{"Service":"redis","Name":"demo-redis-1","State":"running","Health":"","Publishers":[]}` + "\n",
	)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 services, got %d", len(got))
	}
	if got[0].Name != "postgres" || got[1].Name != "redis" {
		t.Errorf("services out of order: %v", got)
	}
	if got[1].Health != "" {
		t.Errorf("empty Health field collapsed: %q", got[1].Health)
	}
	if len(got[1].Ports) != 0 {
		t.Errorf("empty Publishers should produce empty Ports, got %v", got[1].Ports)
	}
}

func TestParseComposePsOutput_NDJSON_SkipsBlankLines(t *testing.T) {
	t.Parallel()
	raw := []byte(
		`{"Service":"postgres","State":"running"}` + "\n" +
			"\n" +
			`{"Service":"redis","State":"running"}` + "\n",
	)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 services after blank-line skip, got %d", len(got))
	}
}

func TestParseComposePsOutput_ArrayForm(t *testing.T) {
	t.Parallel()
	// Compose v2.21+ emits a JSON array instead of NDJSON.
	raw := []byte(`[
		{"Service":"postgres","Name":"demo-postgres-1","State":"running","Health":"healthy","Publishers":[{"PublishedPort":5432,"TargetPort":5432}]},
		{"Service":"redis","Name":"demo-redis-1","State":"running"}
	]`)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 services from array form, got %d", len(got))
	}
}

func TestParseComposePsOutput_MalformedJSON_ReturnsError(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  []byte
	}{
		{"truncated-object", []byte(`{"Service":"postgres"`)},
		{"truncated-array", []byte(`[{"Service":"x"}`)},
		{"not-json", []byte(`this is not json`)},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := docker.ParseComposePsOutputForTest(tc.raw)
			if err == nil {
				t.Fatalf("expected error for malformed input, got nil")
			}
		})
	}
}

func TestParseComposePsOutput_UnknownFieldsAreIgnored(t *testing.T) {
	t.Parallel()
	// Why: pin the "Compose may add fields between versions"
	// robustness contract. The parser must accept unknown JSON keys
	// (json.Unmarshal default behavior) so a Compose upgrade that
	// adds e.g. a `Resources` field does not break u-boot up.
	raw := []byte(`{"Service":"postgres","State":"running","Resources":{"NewField":42},"FutureKey":"futureValue"}`)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error for unknown-field input: %v", err)
	}
	if len(got) != 1 || got[0].Name != "postgres" || got[0].State != "running" {
		t.Errorf("known fields not parsed: %v", got)
	}
}

func TestParseComposePsOutput_MultiplePublishersConcatenate(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"Service":"otel","State":"running","Publishers":[{"PublishedPort":4317,"TargetPort":4317},{"PublishedPort":4318,"TargetPort":4318}]}`)
	got, err := docker.ParseComposePsOutputForTest(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}
	if len(got[0].Ports) != 2 {
		t.Fatalf("expected 2 port mappings, got %v", got[0].Ports)
	}
	if got[0].Ports[0] != "4317:4317" || got[0].Ports[1] != "4318:4318" {
		t.Errorf("port mapping order or format wrong: %v", got[0].Ports)
	}
}

func TestProgressSinkOrDiscard(t *testing.T) {
	t.Parallel()
	// Why: nil sink must default to io.Discard, not panic. Pin the
	// adapter's nil-tolerance convention so a future caller passing
	// a zero-value ComposeUpOptions does not crash the engine.
	if got := docker.ProgressSinkOrDiscardForTest(nil); got != io.Discard {
		t.Errorf("nil sink should default to io.Discard, got %T", got)
	}
	var buf bytes.Buffer
	if got := docker.ProgressSinkOrDiscardForTest(&buf); got != &buf {
		t.Errorf("non-nil sink should pass through, got %T", got)
	}
}

// assertErrIs is the test helper for the sentinel-chain pin used
// across all engine tests.
func assertErrIs(t *testing.T, err error, target error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error wrapping %v, got nil", target)
	}
	if !errors.Is(err, target) {
		t.Fatalf("err = %v; want errors.Is(err, %v)", err, target)
	}
}

// assertComposeServiceEqual is the field-by-field comparator used
// by the NDJSON parser tests. reflect.DeepEqual would also work but
// produces less actionable error messages.
func assertComposeServiceEqual(t *testing.T, got, want driven.ComposeService) {
	t.Helper()
	if got.Name != want.Name {
		t.Errorf("Name = %q, want %q", got.Name, want.Name)
	}
	if got.State != want.State {
		t.Errorf("State = %q, want %q", got.State, want.State)
	}
	if got.Health != want.Health {
		t.Errorf("Health = %q, want %q", got.Health, want.Health)
	}
	if len(got.Ports) != len(want.Ports) {
		t.Errorf("Ports length = %d, want %d (got %v, want %v)", len(got.Ports), len(want.Ports), got.Ports, want.Ports)
		return
	}
	for i := range got.Ports {
		if got.Ports[i] != want.Ports[i] {
			t.Errorf("Ports[%d] = %q, want %q", i, got.Ports[i], want.Ports[i])
		}
	}
}

// TestEngine_ComposeLogs_MissingBinary_ReturnsErrDockerUnavailable
// pins that the slice-v1-logs T1 adapter inherits the 11-vs-12
// Sentinel-Schichtung from M6: a missing docker binary is an
// environment failure (ErrDockerUnavailable / Exit-Code 11), not
// a compose runtime failure (Exit-12).
func TestEngine_ComposeLogs_MissingBinary_ReturnsErrDockerUnavailable(t *testing.T) {
	t.Parallel()
	e := docker.WithEngineBinary("/does/not/exist/docker-binary")
	err := e.ComposeLogs(context.Background(), "/tmp/demo", driven.ComposeLogsOptions{})
	assertErrIs(t, err, driven.ErrDockerUnavailable)
	if errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("missing binary leaked into ErrComposeRuntime: %v", err)
	}
}

// TestEngine_ComposeLogs_AllProbesPass_StreamsToSink pins the
// happy-path streaming contract: /bin/echo passes the preflight
// (exit 0 for every invocation), prints its arguments to stdout,
// and the adapter forwards that stdout to opts.Sink. Pins that
// (a) Sink receives bytes, and (b) the constructed argv contains
// the expected `compose -f <dir>/compose.yaml logs` shape with
// the optional --follow / --tail / services suffixes when set.
func TestEngine_ComposeLogs_AllProbesPass_StreamsToSink(t *testing.T) {
	t.Parallel()
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	var sink bytes.Buffer
	err := e.ComposeLogs(context.Background(), "/tmp/demo", driven.ComposeLogsOptions{
		Services: []string{"postgres"},
		Follow:   true,
		Tail:     "100",
		Sink:     &sink,
	})
	if err != nil {
		t.Fatalf("ComposeLogs with passing preflight: %v", err)
	}
	// /bin/echo prints all of its args to stdout. The argv built
	// inside ComposeLogs starts with `compose -f <dir>/compose.yaml
	// logs` and appends the optional flags + service filter. We
	// don't pin the full string (path-separator portability), but
	// each token must be present.
	got := sink.String()
	for _, want := range []string{"compose", "compose.yaml", "logs", "--follow", "--tail", "100", "postgres"} {
		if !bytes.Contains([]byte(got), []byte(want)) {
			t.Errorf("Sink output missing %q\n  full: %s", want, got)
		}
	}
}

// TestEngine_ComposeLogs_NonZeroExit_WrapsErrComposeRuntime pins
// the post-preflight runtime-error path. `/bin/false` passes the
// preflight (exit 0 on the lookup, exit 1 on the compose-args
// invocation — preflight uses different sub-args). Hmm actually
// /bin/false returns exit 1 for everything including preflight,
// so this scenario requires a different fake. Use a script-style
// approach: not portable; instead use `sh -c 'exit 1'`-ish via
// the existing /bin/false → preflight will fail with
// ErrDockerUnavailable, not what we want.
//
// Practical alternative: pin the ctx.Err()-pass-through path
// (cancel-before-call), which exercises the same return-site
// branch that the slice plan §AK calls out as load-bearing for
// SIGINT semantics.
func TestEngine_ComposeLogs_ContextCanceled_ReturnsCtxErrUnverdeckt(t *testing.T) {
	t.Parallel()
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel BEFORE ComposeLogs runs.
	err := e.ComposeLogs(ctx, "/tmp/demo", driven.ComposeLogsOptions{Sink: io.Discard})
	if err == nil {
		t.Fatalf("expected error from cancelled context, got nil")
	}
	// Slice-v1-logs §AK + Plan-Followup P3: ctx.Err() unverdeckt;
	// MUST NOT be wrapped in ErrComposeRuntime.
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want wrap of context.Canceled", err)
	}
	if errors.Is(err, driven.ErrComposeRuntime) {
		t.Errorf("ctx.Canceled leaked into ErrComposeRuntime — SIGINT-Vertrag verletzt: %v", err)
	}
}

// TestEngine_ComposeLogs_TailEmpty_SkipsFlag_F6 is the Review-
// Followup F6 anti-regression pin: the adapter's `if opts.Tail
// != ""` branch (`engine.go:155`) skips the `--tail`-flag when
// the caller passes an empty string. The LogsService normalises
// empty→"all" before calling, so this branch is never reached
// in production; but a hypothetical direct adapter caller (T4
// E2E, future RPC) must see Compose's own default-all behaviour
// rather than a `--tail <empty>` argv error. Pinned via the
// /bin/echo-arg-inspection pattern.
func TestEngine_ComposeLogs_TailEmpty_SkipsFlag_F6(t *testing.T) {
	t.Parallel()
	shellBinaryAvailable(t, "/bin/echo")
	e := docker.WithEngineBinary("/bin/echo")
	var sink bytes.Buffer
	if err := e.ComposeLogs(context.Background(), "/tmp/demo", driven.ComposeLogsOptions{
		Tail: "", // empty — adapter must skip the flag entirely
		Sink: &sink,
	}); err != nil {
		t.Fatalf("ComposeLogs: %v", err)
	}
	got := sink.String()
	if bytes.Contains([]byte(got), []byte("--tail")) {
		t.Errorf("argv contains `--tail` despite Tail==\"\"; want flag skipped\n  argv echo: %s", got)
	}
}

func TestEngine_ComposeLogs_LineBufferedStdoutAndStderr(t *testing.T) {
	t.Parallel()
	binary := dockerScript(t, `#!/bin/sh
if [ "$1" = "version" ]; then
  echo "24.0.0"
  exit 0
fi
if [ "$1" = "compose" ] && [ "$2" = "version" ]; then
  echo "v2.20.0"
  exit 0
fi
if [ "$1" = "compose" ]; then
  printf 'stdout one\nstdout partial'
  printf 'stderr one\nstderr partial' >&2
  exit 0
fi
exit 1
`)
	e := docker.WithEngineBinary(binary)
	sink := &recordingWriter{}
	if err := e.ComposeLogs(context.Background(), "/tmp/demo", driven.ComposeLogsOptions{
		Tail: "all",
		Sink: sink,
	}); err != nil {
		t.Fatalf("ComposeLogs: %v", err)
	}

	got := sink.chunks()
	sort.Strings(got)
	want := []string{"stderr one\n", "stderr partial", "stdout one\n", "stdout partial"}
	if !equalStringSlices(got, want) {
		t.Errorf("line-buffered chunks = %#v, want %#v", got, want)
	}
}

// TestWrapComposeRunError_F3 is the Review-Followup F3 anti-
// regression pin for the SIGINT-Pass-Through Schicht 1 (post-
// cmd.Run). The wrap-helper was extracted from
// [Engine.ComposeLogs] so it can be unit-tested without a real
// subprocess — the real `exec.CommandContext`-mid-flight-kill
// path is the responsibility of T4's Docker-Tag E2E test.
//
// Four cases pin the contract: (a) nil context + nil run-error
// → nil; (b) nil context + run-error → ErrComposeRuntime; (c)
// canceled context + nil run-error → context.Canceled (the
// "cancel raced with successful exit" edge); (d) canceled
// context + run-error like "signal: killed" → context.Canceled
// unverdeckt, NOT ErrComposeRuntime (the load-bearing mid-flight-
// SIGINT path).
func TestWrapComposeRunError_F3(t *testing.T) {
	t.Parallel()

	t.Run("nil ctx + nil run-error → nil", func(t *testing.T) {
		t.Parallel()
		err := docker.WrapComposeRunErrorForTest(context.Background(), nil, "logs")
		if err != nil {
			t.Errorf("err = %v, want nil", err)
		}
	})

	t.Run("nil ctx + run-error → ErrComposeRuntime", func(t *testing.T) {
		t.Parallel()
		runErr := errors.New("exit status 1")
		err := docker.WrapComposeRunErrorForTest(context.Background(), runErr, "logs")
		if !errors.Is(err, driven.ErrComposeRuntime) {
			t.Errorf("err = %v, want wrap of ErrComposeRuntime", err)
		}
	})

	t.Run("canceled ctx + nil run-error → context.Canceled", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := docker.WrapComposeRunErrorForTest(ctx, nil, "logs")
		if !errors.Is(err, context.Canceled) {
			t.Errorf("err = %v, want context.Canceled", err)
		}
	})

	t.Run("canceled ctx + run-error → context.Canceled UNVERDECKT", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		// Mid-flight SIGINT: exec.CommandContext kills subprocess,
		// cmd.Run() returns "signal: killed". The wrap MUST return
		// ctx.Err() instead of wrapping into ErrComposeRuntime.
		runErr := errors.New("signal: killed")
		err := docker.WrapComposeRunErrorForTest(ctx, runErr, "logs")
		if !errors.Is(err, context.Canceled) {
			t.Errorf("err = %v, want context.Canceled (Schicht-1 post-cmd.Run guard)", err)
		}
		if errors.Is(err, driven.ErrComposeRuntime) {
			t.Errorf("mid-flight cancel leaked into ErrComposeRuntime — Exit-12 statt Exit-0: %v", err)
		}
	})
}

func dockerScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("write fake docker script: %v", err)
	}
	return path
}

type recordingWriter struct {
	chunk []string
}

func (w *recordingWriter) Write(p []byte) (int, error) {
	w.chunk = append(w.chunk, string(p))
	return len(p), nil
}

func (w *recordingWriter) chunks() []string {
	out := make([]string, len(w.chunk))
	copy(out, w.chunk)
	return out
}

func equalStringSlices(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
