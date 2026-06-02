package docker

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// Engine is the production [driven.DockerEngine] adapter. Construct
// with [NewEngine] (default `docker` binary on $PATH) or
// [WithEngineBinary] (test injection).
//
// The Engine reuses the existing [Probe] adapter for the pre-flight
// environment checks — both adapters live in the same package, so
// the reuse stays local without exposing Probe internals beyond
// `internal/adapter/driven/docker`.
type Engine struct {
	// binary is the path to the docker CLI; defaults to "docker".
	binary string
	// probe runs the LookPath / `docker version` / `docker compose
	// version` roundtrips that classify env failures before the
	// real Compose call.
	probe *Probe
}

// Static check: Engine satisfies the DockerEngine port.
var _ driven.DockerEngine = (*Engine)(nil)

// NewEngine returns an Engine that shells out to the `docker` binary
// on `$PATH`.
func NewEngine() *Engine {
	return &Engine{binary: "docker", probe: New()}
}

// WithEngineBinary overrides the docker binary path; intended for
// tests (integration tests under build-tag `docker` may point this
// at a container-runtime alias).
func WithEngineBinary(path string) *Engine {
	return &Engine{binary: path, probe: WithBinary(path)}
}

// preflight runs the three-step environment classifier:
//
//  1. exec.LookPath — the docker binary exists on PATH.
//  2. probe.Info — the daemon answers a roundtrip (covers daemon
//     down, socket permission denied).
//  3. probe.ComposeVersion — the compose plugin is installed.
//
// Any failure returns wrapped [driven.ErrDockerUnavailable] with
// the original cause embedded textually (the cause's identity is
// not preserved by errors.Is; the CLI only needs the sentinel).
// Only when all three succeed does the caller proceed to the real
// Compose invocation; any failure there is classified as
// [driven.ErrComposeRuntime] by the caller.
//
// Performance: two extra roundtrips per Compose call (~100 ms total
// in the typical-case). Acceptable trade-off for the deterministic
// LH-FA-CLI-006 11-vs-12 mapping; the M6 slice notes a future
// `--skip-preflight` flag would be premature optimization.
func (e *Engine) preflight(ctx context.Context) error {
	if _, err := exec.LookPath(e.binary); err != nil {
		return fmt.Errorf("docker binary not on PATH (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	if err := e.probe.Info(ctx); err != nil {
		return fmt.Errorf("docker daemon unreachable (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	if _, err := e.probe.ComposeVersion(ctx); err != nil {
		return fmt.Errorf("docker compose plugin not available (%s): %w", err.Error(), driven.ErrDockerUnavailable)
	}
	return nil
}

// ComposeUp implements [driven.DockerEngine].
func (e *Engine) ComposeUp(ctx context.Context, dir string, opts driven.ComposeUpOptions) (driven.ComposeUpResult, error) {
	if err := e.preflight(ctx); err != nil {
		return driven.ComposeUpResult{}, err
	}
	args := []string{"compose", "-f", filepath.Join(dir, "compose.yaml"), "up"}
	if opts.Detach {
		args = append(args, "-d")
	}
	cmd := exec.CommandContext(ctx, e.binary, args...)
	cmd.Stderr = progressSinkOrDiscard(opts.ProgressSink)
	if err := cmd.Run(); err != nil {
		return driven.ComposeUpResult{}, fmt.Errorf("docker compose up failed (%s): %w", err.Error(), driven.ErrComposeRuntime)
	}
	// LH-FA-UP-001 §970 fire-and-forget pin: ComposeUp must NOT
	// follow up with a `compose ps` roundtrip. The original T2
	// design included a post-up snapshot in ComposeUpResult.Services,
	// but a post-T6 review confirmed UpService never reads the field
	// (the polling loop calls ComposePs separately) and the extra
	// roundtrip could fail with ErrComposeRuntime after a successful
	// up — leaking past the §970 fire-and-forget guarantee.
	return driven.ComposeUpResult{}, nil
}

// ComposeDown implements [driven.DockerEngine].
func (e *Engine) ComposeDown(ctx context.Context, dir string, opts driven.ComposeDownOptions) error {
	if err := e.preflight(ctx); err != nil {
		return err
	}
	args := []string{"compose", "-f", filepath.Join(dir, "compose.yaml"), "down"}
	if opts.RemoveVolumes {
		args = append(args, "-v")
	}
	cmd := exec.CommandContext(ctx, e.binary, args...)
	cmd.Stderr = progressSinkOrDiscard(opts.ProgressSink)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose down failed (%s): %w", err.Error(), driven.ErrComposeRuntime)
	}
	return nil
}

// ComposePs implements [driven.DockerEngine].
func (e *Engine) ComposePs(ctx context.Context, dir string) ([]driven.ComposeService, error) {
	if err := e.preflight(ctx); err != nil {
		return nil, err
	}
	return e.composePs(ctx, dir)
}

// ComposeLogs implements [driven.DockerEngine] (LH-FA-UP-005).
// Streams `docker compose logs` stdout and stderr to opts.Sink,
// line-by-line, so `--follow` callers receive complete log records
// as soon as the subprocess pipe yields them.
//
// SIGINT contract (slice-v1-logs §AK + Plan-Followup P3): when
// the underlying `cmd.Run()` returns and `ctx.Err() != nil`, the
// adapter returns `ctx.Err()` unverdeckt — NOT wrapped in
// `driven.ErrComposeRuntime`. Reason: the application layer
// short-circuits on `errors.Is(err, context.Canceled)` and the
// CLI maps that path to Exit-Code 0 (tail-konform). Wrapping
// would degrade Exit-0 to Exit-12 (compose runtime).
func (e *Engine) ComposeLogs(ctx context.Context, dir string, opts driven.ComposeLogsOptions) error {
	// SIGINT-Pass-Through (1/2): early return for an already-
	// cancelled context. The preflight below would otherwise
	// surface ErrDockerUnavailable (the daemon probe fails on
	// a cancelled ctx), masking the user's Ctrl-C with a 12-vs-0
	// exit-code drift.
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if err := e.preflight(ctx); err != nil {
		return err
	}
	args := []string{"compose", "-f", filepath.Join(dir, "compose.yaml"), "logs"}
	if opts.Follow {
		args = append(args, "--follow")
	}
	if opts.Tail != "" {
		args = append(args, "--tail", opts.Tail)
	}
	args = append(args, opts.Services...)
	cmd := exec.CommandContext(ctx, e.binary, args...)
	return wrapComposeRunError(ctx, runLineBuffered(cmd, progressSinkOrDiscard(opts.Sink)), "logs")
}

func runLineBuffered(cmd *exec.Cmd, sink io.Writer) error {
	var mu sync.Mutex
	stdout := newLineBufferingWriter(lockedWriter{sink: sink, mu: &mu})
	stderr := newLineBufferingWriter(lockedWriter{sink: sink, mu: &mu})
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	runErr := cmd.Run()
	stdoutErr := stdout.Flush()
	stderrErr := stderr.Flush()
	if runErr != nil {
		return runErr
	}
	if stdoutErr != nil {
		return stdoutErr
	}
	return stderrErr
}

type lockedWriter struct {
	sink io.Writer
	mu   *sync.Mutex
}

func (w lockedWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.sink.Write(p)
}

type lineBufferingWriter struct {
	dst     io.Writer
	pending []byte
}

func newLineBufferingWriter(dst io.Writer) *lineBufferingWriter {
	return &lineBufferingWriter{dst: dst}
}

func (w *lineBufferingWriter) Write(p []byte) (int, error) {
	written := 0
	for len(p) > 0 {
		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			w.pending = append(w.pending, p...)
			written += len(p)
			return written, nil
		}
		segment := p[:idx+1]
		w.pending = append(w.pending, segment...)
		written += len(segment)
		if _, err := w.dst.Write(w.pending); err != nil {
			return written, err
		}
		w.pending = w.pending[:0]
		p = p[idx+1:]
	}
	return written, nil
}

func (w *lineBufferingWriter) Flush() error {
	if len(w.pending) == 0 {
		return nil
	}
	_, err := w.dst.Write(w.pending)
	w.pending = w.pending[:0]
	return err
}

// wrapComposeRunError implements the SIGINT-Pass-Through Schicht 1
// (post-cmd.Run): if `ctx.Err() != nil` the context cancellation
// is returned unverdeckt — `exec.CommandContext` killed the
// subprocess, so `runErr` would otherwise be `signal: killed` and
// get wrapped into ErrComposeRuntime, flipping Exit-0 (tail-
// konform Ctrl-C) into Exit-12 (compose runtime). Extracted from
// [Engine.ComposeLogs] for Review-Followup F3: the helper is
// unit-testable without a real subprocess.
//
// kind is the compose verb name used in the error message
// (`"logs"`, future `"exec"`, …). Returns nil when runErr is nil
// and ctx is healthy.
func wrapComposeRunError(ctx context.Context, runErr error, kind string) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return ctxErr
	}
	if runErr != nil {
		return fmt.Errorf("docker compose %s failed (%s): %w", kind, runErr.Error(), driven.ErrComposeRuntime)
	}
	return nil
}

// composePs is the preflight-less internal helper; called by both
// the public ComposePs (which preflights first) and ComposeUp
// (which already preflighted at the start).
func (e *Engine) composePs(ctx context.Context, dir string) ([]driven.ComposeService, error) {
	cmd := exec.CommandContext(ctx, e.binary, "compose", "-f", filepath.Join(dir, "compose.yaml"), "ps", "--format", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker compose ps failed (%s): %w", err.Error(), driven.ErrComposeRuntime)
	}
	services, parseErr := parseComposePsOutput(out)
	if parseErr != nil {
		return nil, fmt.Errorf("parse compose ps output (%s): %w", parseErr.Error(), driven.ErrComposeRuntime)
	}
	return services, nil
}

// progressSinkOrDiscard returns the caller-provided writer or
// [io.Discard] when nil. Adapter convention: nil sinks mean "don't
// forward Compose stderr" rather than panicking on a nil writer.
func progressSinkOrDiscard(sink io.Writer) io.Writer {
	if sink == nil {
		return io.Discard
	}
	return sink
}

// composePsLine is the JSON shape of one `docker compose ps
// --format json` entry. Compose v2.20+ originally emitted NDJSON
// (one object per line); v2.21+ may emit a JSON array. Both shapes
// are accepted by [parseComposePsOutput].
//
// Fields are intentionally a subset of what Compose returns —
// unknown fields are ignored by the json decoder so a Compose
// upgrade adding new fields cannot break the parser.
type composePsLine struct {
	// Service is the Compose service name (e.g. "postgres"). The
	// adapter maps this to [driven.ComposeService.Name] (not the
	// container name, which carries the project prefix).
	Service string `json:"Service"`

	// Name is the container name (e.g. "demo-postgres-1");
	// captured for completeness but not surfaced via
	// [driven.ComposeService]. Kept here so a future field
	// surfacing it (e.g. for `--json` output, V1) does not need a
	// re-parse round.
	Name string `json:"Name"`

	// State is the raw Compose container state ("running",
	// "restarting", "exited", …).
	State string `json:"State"`

	// Health is the raw healthcheck status ("healthy",
	// "unhealthy", "starting") or empty.
	Health string `json:"Health"`

	// Publishers lists the host-published ports for the service.
	// Empty for services without exposed ports.
	Publishers []composePsPublisher `json:"Publishers"`
}

// composePsPublisher is the Compose JSON shape for one host->
// container port mapping. PublishedPort=0 means the port is
// declared but not bound to a host port (rare; usually a no-publish
// `expose:` style entry).
type composePsPublisher struct {
	PublishedPort int    `json:"PublishedPort"`
	TargetPort    int    `json:"TargetPort"`
	Protocol      string `json:"Protocol"`
}

// parseComposePsOutput accepts either NDJSON (Compose v2.20) or a
// JSON array (Compose v2.21+) and produces the canonical
// [driven.ComposeService] slice. Empty input → nil slice, no error
// (Compose returns nothing when the project has no services
// running).
func parseComposePsOutput(raw []byte) ([]driven.ComposeService, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return nil, nil
	}
	lines, err := decodeComposePsLines(raw)
	if err != nil {
		return nil, err
	}
	result := make([]driven.ComposeService, 0, len(lines))
	for _, l := range lines {
		svc := driven.ComposeService{
			Name:   l.Service,
			State:  l.State,
			Health: l.Health,
		}
		for _, p := range l.Publishers {
			svc.Ports = append(svc.Ports, fmt.Sprintf("%d:%d", p.PublishedPort, p.TargetPort))
		}
		result = append(result, svc)
	}
	return result, nil
}

// decodeComposePsLines dispatches on the first non-whitespace byte
// of `raw` to pick the right JSON shape decoder: `[` → array form
// (Compose v2.21+), anything else → NDJSON (Compose v2.20). Caller
// guarantees `raw` is already trimmed and non-empty.
func decodeComposePsLines(raw []byte) ([]composePsLine, error) {
	if raw[0] == '[' {
		var lines []composePsLine
		if err := json.Unmarshal(raw, &lines); err != nil {
			return nil, fmt.Errorf("parse JSON array: %w", err)
		}
		return lines, nil
	}
	return decodeNDJSON(raw)
}

// decodeNDJSON scans `raw` line-by-line, json-Unmarshal'ing each
// non-empty line into a [composePsLine]. Empty/whitespace-only
// lines are skipped — Compose's NDJSON output sometimes carries a
// trailing newline, and the project's CI fixtures may stack blank
// lines from heredoc construction.
func decodeNDJSON(raw []byte) ([]composePsLine, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	// Compose ps output per service is small (~500 bytes); the
	// 1 MiB buffer here is generous and covers projects with
	// dozens of services without per-line allocation.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var lines []composePsLine
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var l composePsLine
		if err := json.Unmarshal(line, &l); err != nil {
			return nil, fmt.Errorf("parse NDJSON line: %w", err)
		}
		lines = append(lines, l)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan NDJSON: %w", err)
	}
	return lines, nil
}
