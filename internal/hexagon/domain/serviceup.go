package domain

import "strings"

// ContainerState is the u-boot-normalized form of `docker compose ps`
// container states. The raw Compose state strings are case-sensitive
// and have grown across Compose versions (e.g. new transitional
// states for pull/health handshakes); a Freitext-Vergleich in the
// polling loop would silently break on every Compose upgrade. The
// enum below pins a fixed vocabulary; [ParseContainerState] is the
// single entry point that maps raw Compose strings to the enum.
//
// Layer rules (LH-FA-ARCH-002, LH-FA-ARCH-003): pure domain type,
// imported by application (M6 [UpService] in T4) for the
// stabilization classifier. Adapters obtain the raw Compose string
// from `docker compose ps --format json` (T2 driven port) and pass
// it to the application; the application calls [ParseContainerState].
type ContainerState int

const (
	// StateUnknown represents any Compose state string that does not
	// match the explicit allowlist below. M6 treats unknown states
	// as `RunningOnly` (poll continues) plus a persistent
	// `up.state.<service>.unknown` warn diagnostic, so a Compose
	// upgrade that introduces a new state value never causes a hard
	// failure. Zero value of the enum on purpose — uninitialized
	// state defaults to the conservative "we don't know" path.
	StateUnknown ContainerState = iota

	// StateStarting covers the three Compose states that signal the
	// container has not yet reached `running` (`created`, `starting`,
	// `paused`). Polling continues until the container transitions or
	// the timeout expires.
	StateStarting

	// StateRunning is the main happy path. The stabilization
	// classifier then looks at the service's healthcheck and TCP
	// port to decide between [OutcomeStabilized] and
	// [OutcomeRunningOnly] (LH-FA-UP-001 §966–§969).
	StateRunning

	// StateRestarting is a transitional state. A single restart tick
	// is normal; a sustained restart loop is detected via
	// [RestartLoopThreshold] consecutive observations.
	StateRestarting

	// StateDead aggregates the explicit "container is gone or
	// failed" Compose states (`exited`, `dead`, `removing`,
	// `removed`). These are the *only* states that map to a hard
	// [OutcomeFailed] without retry, per the M6 slice's explicit
	// dead-allowlist.
	StateDead
)

// String returns the canonical lowercase identifier of the container
// state, used by the [Logger] port and the CLI's diagnostic output.
// The identifiers are stable: CI dashboards and log scrapers may pin
// them.
func (s ContainerState) String() string {
	switch s {
	case StateUnknown:
		return "unknown"
	case StateStarting:
		return "starting"
	case StateRunning:
		return "running"
	case StateRestarting:
		return "restarting"
	case StateDead:
		return "dead"
	default:
		return "unknown"
	}
}

// ParseContainerState maps a raw `docker compose ps` state string to
// the u-boot-normalized [ContainerState]. The match is
// case-insensitive and trims surrounding whitespace so a Compose
// release that flips casing (`"Running"` vs. `"running"`) does not
// break the polling classifier. Any string outside the explicit
// allowlist returns [StateUnknown] — see the constant's doc for the
// soft-fail rationale.
func ParseContainerState(raw string) ContainerState {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "running":
		return StateRunning
	case "restarting":
		return StateRestarting
	case "created", "starting", "paused":
		return StateStarting
	case "exited", "dead", "removing", "removed":
		return StateDead
	default:
		return StateUnknown
	}
}

// RestartLoopThreshold is the number of consecutive
// [StateRestarting] observations the [UpService] polling loop
// tolerates before classifying a service as [OutcomeFailed]. A
// single `restarting` tick is normal Compose behavior; a sustained
// loop manifests across multiple poll iterations.
//
// The threshold is a domain constant (not application-tunable) so
// the test surface in T4 can pin the value deterministically and so
// the LH-FA-UP-001 stabilization semantics stay stable across
// releases. A future per-project override would be a separate slice.
const RestartLoopThreshold = 3

// StabilizationOutcome is the per-service classifier result inside
// the [UpService] polling loop. It is recomputed every iteration
// from [ContainerState] + healthcheck signal + TCP port probe.
type StabilizationOutcome int

const (
	// OutcomeRunningOnly means the service has not yet reached a
	// terminal classification — polling continues until the timeout
	// expires. Zero value of the enum on purpose: an unclassified
	// outcome is treated as "keep polling", which is fail-safe (the
	// loop's timeout will eventually trigger
	// [ErrStabilizationTimeout]).
	OutcomeRunningOnly StabilizationOutcome = iota

	// OutcomeStabilized is the terminal success: the service has
	// reached the Zielzustand per LH-FA-UP-001 (`healthy` for
	// services with a healthcheck, `running` for services without,
	// plus a successful TCP port probe when the service declares
	// a probable port).
	OutcomeStabilized

	// OutcomeFailed is the terminal failure: a service has been
	// observed in [StateDead] or has crossed [RestartLoopThreshold]
	// consecutive [StateRestarting] observations. Polling stops
	// immediately; the use case returns wrapped
	// `driven.ErrComposeRuntime` (M6 slice T4-step 6).
	OutcomeFailed
)

// String returns the canonical lowercase identifier of the outcome.
// Used by the [Logger] port and by tests to assert classifier
// results.
func (o StabilizationOutcome) String() string {
	switch o {
	case OutcomeRunningOnly:
		return "running-only"
	case OutcomeStabilized:
		return "stabilized"
	case OutcomeFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// ServiceStatus is the per-service snapshot the [UpService] returns
// inside [UpResult]. The fields exactly match the four mandatory
// columns of LH-FA-UP-003: Service-Name, Container-Status, Port,
// Healthcheck.
type ServiceStatus struct {
	// Name is the Compose service name (e.g. "postgres"). Stable
	// across iterations.
	Name string

	// ContainerStatus is the normalized container state at the
	// moment the result was assembled — typically the last poll
	// iteration before stabilization or timeout.
	ContainerStatus ContainerState

	// Port is the human-readable port mapping list as the CLI
	// renders it (e.g. "5432:5432" or "5432:5432, 127.0.0.1:9091:9091").
	// Empty when the service declares no ports. The application
	// constructs the string from the per-port probe targets so the
	// CLI does not need to re-render Compose syntax.
	Port string

	// Healthcheck is the Compose-reported health status string
	// ("healthy", "unhealthy", "starting") or empty for services
	// without a healthcheck. Pinning the strings to Compose's
	// vocabulary lets CI dashboards filter on stable values.
	Healthcheck string
}

// UpResult is the aggregate result of one `u-boot up` invocation.
// Returned inside [driving.UpResponse]; the CLI adapter renders
// Services as the LH-FA-UP-003 status table and Diagnostics as the
// warn/info section below it.
//
// `--timeout=0` fire-and-forget pin: in that mode the use case
// returns Services=nil, Stabilized=false, and Diagnostics carrying
// a single `up.fire-and-forget` [SeverityInfo] entry. The CLI then
// renders only the info line, no status table.
type UpResult struct {
	// Services is the per-service snapshot. Order is deterministic
	// (alphabetical by Name) so golden-file tests can assert
	// byte-exact output.
	Services []ServiceStatus

	// Stabilized is true exactly when every entry in Services
	// reached [OutcomeStabilized] within the request's timeout.
	// False when the use case returned early via fire-and-forget or
	// via [ErrStabilizationTimeout] (though the timeout path
	// normally returns a non-nil error before constructing the
	// result).
	Stabilized bool

	// Diagnostics carries non-fatal observations from the polling
	// loop — port-parse warns from
	// [driving.UpRequest]-declared but non-TCP-probable ports
	// (LH-FA-UP-001 §969), unknown-state warns from Compose-states
	// outside [ParseContainerState]'s allowlist, and the
	// `up.fire-and-forget` info entry when Timeout=0. Diagnostics
	// are part of the domain contract (not CLI-side) so the
	// application can assert them in tests.
	Diagnostics []Diagnostic
}
