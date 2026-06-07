package driving

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// UpRequest is the input for [UpUseCase.Up]. It is the
// application-layer expression of `u-boot up` per LH-FA-UP-001 /
// LH-FA-UP-002; the CLI adapter (M6-T7) translates the `--timeout`
// flag into [UpRequest.Timeout] via
// `time.Duration(secs) * time.Second` — never as raw
// `time.Duration(secs)` (which would be nanoseconds).
type UpRequest struct {
	// BaseDir is the absolute path of the initialized u-boot project.
	// Mandatory; the CLI adapter defaults it to the current working
	// directory (mirroring `u-boot init` / `u-boot add`).
	BaseDir string

	// Timeout is the maximum wall-clock duration the polling loop
	// waits for every declared service to reach
	// [domain.OutcomeStabilized]. LH-FA-UP-001 §963: the spec default
	// is 60 s; the CLI adapter applies that default before
	// constructing the request, so the application sees a populated
	// value here.
	//
	// Semantics:
	//
	//   - Timeout > 0: regular polling loop, bounded by this
	//     duration; if any service is still not stabilized when the
	//     bound is hit the use case returns [ErrStabilizationTimeout].
	//   - Timeout == 0: fire-and-forget. The use case returns
	//     immediately after a successful `ComposeUp`, with no
	//     `ComposePs` roundtrip and no port/healthcheck probes
	//     (LH-FA-UP-001 §970). The result carries a single
	//     `up.fire-and-forget` [domain.SeverityInfo] diagnostic.
	//   - Timeout < 0: validation error. The use case returns a
	//     non-sentinel error before any Compose call; the CLI maps
	//     it to LH-FA-CLI-006 exit code 2.
	Timeout time.Duration

	// ProgressSink is the writer the application passes to the
	// `DockerEngine.ComposeUp` adapter so the Compose stderr stream
	// (pull/create/start/healthcheck phases per LH-NFA-PERF-002)
	// reaches the user. The CLI adapter wires this to `os.Stderr`;
	// `nil` is treated as `io.Discard`.
	// `--quiet` does **not** silence this stream — see the M6 slice
	// for the rationale.
	ProgressSink io.Writer

	// SilenceProgress switches the application-layer ProgressSink to
	// `io.Discard` for the duration of this request, so the Compose
	// stderr stream does not pollute machine-consumable output
	// (slice-v1-cli-json-dry-run-up-down T0-(c) form (d)). The CLI
	// adapter sets this to `true` when `--json` is active. Pattern is
	// symmetric to [RemoveServiceRequest.SilenceConfirmer]: a boolean
	// request flag that lets the use case branch internally without
	// the CLI having to know about `io.Discard`.
	SilenceProgress bool
}

// UpResponse is the output of [UpUseCase.Up]. The CLI adapter (T6)
// renders [Result.Services] as the LH-FA-UP-003 status table and
// [Result.Diagnostics] as the warn / info section below.
type UpResponse struct {
	// Result carries the per-service snapshot and the
	// Stabilized flag. See [domain.UpResult] for the contract on the
	// `--timeout=0` fire-and-forget shape (Services=nil,
	// Stabilized=false, Diagnostics carrying a single info entry).
	Result domain.UpResult

	// Warnings carries soft-warning diagnostics the use case wants the
	// CLI adapter to surface via the JSON envelope's `diagnostics[]`
	// array (slice-v1-cli-json-dry-run-up-down T2 — type inherited
	// from remove T2 cluster-vorlauf R12-LOW-F4). Empty / nil on the
	// happy path; populated when recreate-detection (T0-(k) follow-up
	// slice) or future read-side WARN paths land.
	Warnings []WarningEntry
}

// All Up sentinels below live in the `driving` package (not in
// `application`) so the CLI adapter can branch on them via
// [errors.Is] without importing `application` — the LH-FA-ARCH-003
// depguard rule forbids that cross-layer import.
//
// Sentinel layering (M6 slice §Sentinel-Schichtung): driving
// sentinels here are the *application-level* failures. The
// *environment* and *runtime* failures from the docker engine ride
// on `driven.ErrDockerUnavailable` (→ exit code 11) and
// `driven.ErrComposeRuntime` (→ exit code 12); the application
// returns them through a contextual `fmt.Errorf("...: %w", err)`
// wrap, never under a driving sentinel — that would shadow the
// original sentinel identity and make exit code 11 unreachable
// from `errors.Is` at the CLI.

// ErrUpFileSystem signals that the up use case hit a raw filesystem
// error during the read-only phase that loads `u-boot.yaml` /
// `compose.yaml` (slice-v1-cli-json-dry-run-up-down T2 / T0-(d)
// inherited from remove's ErrRemoveFileSystem; up/down read-only
// so the message form is "read failed", not "mutation failed").
// T3 wraps the FS-Read Stellen in upservice.go (Z. 105, 138, 148)
// with multi-`%w`-form (Go 1.20+):
//
//	`fmt.Errorf("up service: <action>(%q): %w: %w", path,
//	            ErrUpFileSystem, rawErr)`
//
// Switch-Order in `mapUpErrorToDiagnostic` (T0-(e)) MUST check
// ErrUpFileSystem FIRST so multi-`%w` chains that include both
// ErrUpFileSystem AND a fachlich sentinel route to the FS-class
// (LH-NFA-REL-003 / exit 14), not the fachlich-class. Maps to
// LH-NFA-REL-003 exit code 14 via cli's `isFilesystemError`.
var ErrUpFileSystem = errors.New("up: filesystem read failed")

// ErrComposeFileMissing signals that `BaseDir` contains a
// `u-boot.yaml` (passes the [ErrProjectNotInitialized] gate) but no
// `compose.yaml` — the file the engine would feed to
// `docker compose -f`. Distinct from [ErrProjectNotInitialized]
// because the user message and repair path differ: missing
// `u-boot.yaml` → "run `u-boot init`"; missing `compose.yaml` →
// "compose file was deleted, restore it or re-init". Maps to
// LH-FA-UP-001 exit code 10 (fachliche Validierung).
var ErrComposeFileMissing = errors.New("compose file missing")

// ErrStabilizationTimeout signals that the polling loop reached the
// [UpRequest.Timeout] bound without classifying every service as
// [domain.OutcomeStabilized]. The wrapped error carries the list of
// services still pending so the CLI can render the offending names.
// Maps to LH-FA-UP-001 exit code 12 (fachlicher Ausführungsfehler).
//
// Distinct from `driven.ErrComposeRuntime` (T2) even though both
// map to exit code 12: the timeout is u-boot's own polling-loop
// concern, not a Compose-stderr observation. Keeping the sentinels
// separate lets future tooling distinguish "Compose said no" from
// "u-boot gave up waiting".
var ErrStabilizationTimeout = errors.New("stabilization timeout")

// UpUseCase is the driving-port for `u-boot up`. The CLI adapter
// holds a reference and calls [Up] from the Cobra command handler.
//
// Contract:
//
//   - On success the response carries a populated [domain.UpResult]
//     and the error is nil. For Timeout > 0, Stabilized is true; for
//     Timeout == 0 (fire-and-forget), Stabilized is false and the
//     diagnostics carry the `up.fire-and-forget` info entry.
//   - On a use-case failure the response is the zero value and the
//     error wraps one of the sentinels above, [ErrProjectNotInitialized]
//     (defined in `addservice.go` and shared across use cases), or a
//     `driven.*` sentinel forwarded from the engine.
//   - The returned error carries the Engine sentinel identity intact
//     via `fmt.Errorf("up service: ...: %w", err)`, so
//     `errors.Is(err, driven.ErrDockerUnavailable)` and
//     `errors.Is(err, driven.ErrComposeRuntime)` continue to hold at
//     the CLI level.
type UpUseCase interface {
	Up(ctx context.Context, req UpRequest) (UpResponse, error)
}
