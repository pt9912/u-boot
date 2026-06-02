package driven

import (
	"context"
	"errors"
	"io"
)

// ErrDockerUnavailable signals a Docker-environment-class failure:
// the host `docker` binary is not on PATH, the daemon is not
// reachable, or the `docker compose` plugin is missing. The adapter
// detects this *before* attempting the actual `compose up`/`down`/
// `ps` call via three pre-probes (LookPath, `docker version
// --format '{{.Server.Version}}'`, `docker compose version
// --short`) — not by parsing Compose stderr, because Compose
// error messages drift between versions and a string-match would
// silently break.
//
// CLI mapping: LH-FA-CLI-006 exit code 11 (fachlicher
// Umgebungsfehler). M6 slice §Sentinel-Schichtung pins the
// errors.Is-Durchleitungs-Vertrag: the application layer wraps
// engine errors only with contextual `fmt.Errorf("...: %w", err)`,
// never under a driving-layer sentinel — so `errors.Is(err,
// driven.ErrDockerUnavailable)` continues to hold at the CLI.
var ErrDockerUnavailable = errors.New("docker unavailable")

// ErrComposeRuntime signals a Compose-runtime-class failure that
// surfaced *after* the pre-probes confirmed a working Docker
// environment: `compose up` exited non-zero, `compose down` failed,
// `compose ps` produced unparseable output, or any other Compose-
// side execution error.
//
// CLI mapping: LH-FA-CLI-006 exit code 12 (fachlicher
// Ausführungsfehler). Distinct from
// `driving.ErrStabilizationTimeout` (which also maps to 12): the
// runtime sentinel is "Compose said no", the stabilization sentinel
// is "u-boot gave up waiting" — same exit code, different cause.
var ErrComposeRuntime = errors.New("compose runtime error")

// ComposeUpOptions configures a [DockerEngine.ComposeUp] call.
type ComposeUpOptions struct {
	// Detach mirrors `docker compose up -d`. M6 always sets this
	// true: the polling loop in the [UpService] application
	// service decides when the call returns, not the foreground
	// compose process.
	Detach bool

	// ProgressSink is the writer the adapter forwards Compose
	// stderr to — the live `Pulling…` / `Creating…` / `Starting…`
	// / `Healthchecking…` phase stream demanded by LH-NFA-PERF-002.
	// `nil` is treated as `io.Discard`. The CLI wires this to
	// `os.Stderr`; `--quiet` does **not** silence the stream
	// (see the M6 slice for the rationale).
	ProgressSink io.Writer
}

// ComposeDownOptions configures a [DockerEngine.ComposeDown] call.
type ComposeDownOptions struct {
	// RemoveVolumes mirrors `docker compose down -v`. When true,
	// named volumes are dropped alongside containers — the
	// LH-FA-UP-004 destructive path that needs the LH-FA-CLI-005A
	// confirmation flow before reaching the adapter.
	RemoveVolumes bool

	// ProgressSink behaves the same as in [ComposeUpOptions].
	ProgressSink io.Writer
}

// ComposeLogsOptions configures a [DockerEngine.ComposeLogs] call
// (LH-FA-UP-005). Slice-v1-logs §T0-Outcomes pinned the surface:
// minimal Compose-Facade — only the Spec-mandated `--follow` and
// `--tail` knobs are exposed; no `--no-log-prefix`/`--timestamps`
// (the latter live in a future Folge-Slice with a documented
// trigger).
type ComposeLogsOptions struct {
	// Services is the list of service names to filter on (empty =
	// Compose-Default = all services in compose.yaml). Slice-v1-logs
	// T0-(a) decision: no `activeServiceNames(cfg)`-filter — Compose
	// decides, so manually-added compose.yaml services stay visible.
	// The application service passes either an empty slice (no
	// positional argument) or a single-element slice (`u-boot logs
	// <service>`) here; the adapter just forwards.
	Services []string

	// Follow mirrors `docker compose logs --follow`. When true the
	// adapter blocks until the underlying compose process exits —
	// SIGINT propagation lives in the ctx.Err()-pass-through
	// contract documented on [DockerEngine.ComposeLogs].
	Follow bool

	// Tail is the value passed to `docker compose logs --tail
	// <Tail>`. Compose accepts the literal `"all"` (all available
	// lines) or a decimal integer string (`"0"`, `"100"`, …).
	// Slice-v1-logs T0-(c): the application service normalises an
	// empty incoming string to `"all"` and validates numeric
	// inputs >= 0 in the CLI layer; the adapter trusts the value.
	Tail string

	// Sink is the writer the adapter forwards BOTH of Compose's
	// streams to — stdout (log lines) AND stderr (compose status
	// like `Attaching to …`, service-exit notices). Review-Followup
	// F2: the previous Doc-Kommentar promised "stdout only", but
	// `cmd.Stderr` is also routed here to keep T0-(d) intact (Spec-
	// treu: Compose-Output unverändert). Splitting the streams to
	// distinct writers would either require a `--no-log-prefix`-
	// style discriminator the slice deliberately omits, or a
	// second sink field which would re-introduce the API surface
	// that T0-(d) closed. So: one Sink for both, documented.
	//
	// `nil` is treated as `io.Discard` by the adapter
	// (`progressSinkOrDiscard`). The CLI wires this to
	// `cmd.OutOrStdout()` via [driving.LogsRequest.OutputSink].
	Sink io.Writer
}

// ComposeService is the per-service snapshot returned by `docker
// compose ps --format json`. The fields are the *raw* values from
// Compose; the application service normalizes them
// (state → [domain.ContainerState] via
// `domain.ParseContainerState`, health → display string,
// ports → probe targets).
//
// Layer rules: the adapter exposes raw observations here; the
// classification (running vs. healthy vs. failed) is the
// application layer's responsibility.
type ComposeService struct {
	// Name is the Compose service name (e.g. "postgres"), not the
	// container name. Stable across runs.
	Name string

	// State is the raw Compose state string ("running",
	// "restarting", "exited", …). The application service feeds
	// this to [domain.ParseContainerState] for case-insensitive
	// normalization.
	State string

	// Health is the raw Compose health status ("healthy",
	// "unhealthy", "starting") or empty for services without a
	// healthcheck. Pinning the strings to Compose's vocabulary
	// lets CI dashboards filter on stable values.
	Health string

	// Ports is the list of host-published port mappings in the
	// canonical form "host:container" (e.g. "5432:5432"). Empty
	// for services without exposed ports.
	Ports []string
}

// ComposeUpResult is the return value of [DockerEngine.ComposeUp].
// Carries no fields today — a post-T6 review (M6 closure) dropped
// the original ComposeUp-internal `compose ps` snapshot to honor
// the LH-FA-UP-001 §970 fire-and-forget contract (no `ps`
// roundtrip after a successful `up`). The type stays so future
// metadata (per-call timing, image-pull stats) can land without
// another signature break; for M6 it is intentionally empty.
type ComposeUpResult struct{}

// DockerEngine is the state-mutating Compose port (M6+), separate
// from the read-only [DockerProbe] used by M4 `doctor`. Splitting
// the two keeps each port narrow: probe answers "is Docker
// available?", engine answers "start/stop/inspect a Compose
// project".
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): driven port; the
// production implementation lives in
// `internal/adapter/driven/docker/engine.go` and is the only place
// `os/exec docker compose ...` runs.
//
// Pre-flight contract: every method runs the three-step
// environment pre-probe (LookPath + `docker version` roundtrip +
// `docker compose version`) before the actual Compose call. A
// failure in any pre-probe step returns wrapped
// [ErrDockerUnavailable] (CLI code 11). A failure in the actual
// `compose up`/`down`/`ps` call returns wrapped
// [ErrComposeRuntime] (CLI code 12). This determinism is what
// makes the LH-FA-CLI-006 11-vs-12 distinction reachable from
// `errors.Is` at the CLI.
type DockerEngine interface {
	// ComposeUp shells out to `docker compose -f <dir>/compose.yaml
	// up [-d]` and, on success, returns a snapshot via
	// `docker compose ps`. Methods take a [context.Context] because
	// Compose can block on network (pulls) or hang on a stale
	// socket.
	ComposeUp(ctx context.Context, dir string, opts ComposeUpOptions) (ComposeUpResult, error)

	// ComposeDown shells out to `docker compose -f <dir>/compose.yaml
	// down [-v]`. No snapshot is returned — `down` is the inverse
	// of `up` and the application service does not poll after it.
	ComposeDown(ctx context.Context, dir string, opts ComposeDownOptions) error

	// ComposePs shells out to `docker compose -f <dir>/compose.yaml
	// ps --format json`. Used by the [UpService] polling loop for
	// per-iteration state observation; a failure here aborts the
	// polling loop and is forwarded to the CLI with the same
	// errors.Is identity as a `ComposeUp` failure.
	ComposePs(ctx context.Context, dir string) ([]ComposeService, error)

	// ComposeLogs shells out to `docker compose -f <dir>/compose.yaml
	// logs [--follow] [--tail <value>] [<service>…]` (LH-FA-UP-005).
	// Streams the compose stdout and stderr to `opts.Sink` line-
	// buffered so `--follow` arrives real-time once the subprocess
	// pipe yields complete log records.
	//
	// **SIGINT contract (slice-v1-logs §AK + Plan-Followup P3):**
	// When the call returns and `ctx.Err() != nil`
	// (`context.Canceled` or `context.DeadlineExceeded`), the
	// adapter MUST return `ctx.Err()` UNVERDECKT — it MUST NOT
	// wrap it in [ErrComposeRuntime] or any other sentinel.
	// Reason: Ctrl-C in the `--follow` path maps to CLI Exit-Code 0
	// via the application service's
	// `errors.Is(err, context.Canceled)` short-circuit; wrapping
	// would degrade Exit-0 to Exit-12 (compose runtime). All other
	// non-zero-exit conditions wrap into [ErrComposeRuntime] as
	// usual.
	ComposeLogs(ctx context.Context, dir string, opts ComposeLogsOptions) error
}
