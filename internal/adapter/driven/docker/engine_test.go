package docker_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os/exec"
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
