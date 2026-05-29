package domain_test

import (
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestContainerState_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		state domain.ContainerState
		want  string
	}{
		{domain.StateUnknown, "unknown"},
		{domain.StateStarting, "starting"},
		{domain.StateRunning, "running"},
		{domain.StateRestarting, "restarting"},
		{domain.StateDead, "dead"},
		{domain.ContainerState(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.state.String(); got != tc.want {
			t.Errorf("ContainerState(%d).String() = %q, want %q", tc.state, got, tc.want)
		}
	}
}

func TestParseContainerState(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  string
		want domain.ContainerState
	}{
		// Direct matches from the Compose vocabulary.
		{"running-lower", "running", domain.StateRunning},
		{"restarting-lower", "restarting", domain.StateRestarting},
		{"created-lower", "created", domain.StateStarting},
		{"starting-lower", "starting", domain.StateStarting},
		{"paused-lower", "paused", domain.StateStarting},
		{"exited-lower", "exited", domain.StateDead},
		{"dead-lower", "dead", domain.StateDead},
		{"removing-lower", "removing", domain.StateDead},
		{"removed-lower", "removed", domain.StateDead},

		// Case-insensitive: pins the M6 slice's "Compose may flip
		// casing between releases" robustness contract.
		{"running-title", "Running", domain.StateRunning},
		{"running-upper", "RUNNING", domain.StateRunning},
		{"exited-title", "Exited", domain.StateDead},
		{"restarting-upper", "RESTARTING", domain.StateRestarting},

		// Whitespace trim — Compose JSON output is normally clean,
		// but a future format-string adapter or a `ps --format`
		// override could leak whitespace. Trim defensively.
		{"running-trimmed", "  running  ", domain.StateRunning},

		// Soft-unknown: any string outside the allowlist must
		// degrade to StateUnknown so a Compose upgrade with a new
		// state value never causes a hard u-boot up failure.
		{"unknown-empty", "", domain.StateUnknown},
		{"unknown-new-value", "frobnicating", domain.StateUnknown},
		{"unknown-pull", "pulling", domain.StateUnknown},
		{"unknown-typo", "runing", domain.StateUnknown},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := domain.ParseContainerState(tc.raw); got != tc.want {
				t.Errorf("ParseContainerState(%q) = %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestStabilizationOutcome_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		outcome domain.StabilizationOutcome
		want    string
	}{
		{domain.OutcomeRunningOnly, "running-only"},
		{domain.OutcomeStabilized, "stabilized"},
		{domain.OutcomeFailed, "failed"},
		{domain.StabilizationOutcome(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.outcome.String(); got != tc.want {
			t.Errorf("StabilizationOutcome(%d).String() = %q, want %q", tc.outcome, got, tc.want)
		}
	}
}

func TestStabilizationOutcome_ZeroValueIsRunningOnly(t *testing.T) {
	t.Parallel()
	// Why: an uninitialized [domain.StabilizationOutcome] must default
	// to "keep polling", not to a terminal classification. If a
	// future renumbering flipped the zero value to Stabilized or
	// Failed, the application's polling loop would either
	// false-positive on the first tick (Stabilized) or abort
	// immediately (Failed). Pin the zero value here so the regression
	// is visible.
	var zero domain.StabilizationOutcome
	if zero != domain.OutcomeRunningOnly {
		t.Errorf("zero-value StabilizationOutcome = %v, want OutcomeRunningOnly", zero)
	}
}

func TestRestartLoopThreshold_PinnedValue(t *testing.T) {
	t.Parallel()
	// Why: the M6 slice fixes the threshold at 3 consecutive
	// restart observations. T4's polling-loop tests rely on this
	// exact value to construct deterministic fixtures (two restart
	// ticks must recover, three must fail). A drift would silently
	// move the boundary.
	if domain.RestartLoopThreshold != 3 {
		t.Errorf("RestartLoopThreshold = %d, want 3", domain.RestartLoopThreshold)
	}
}

func TestStateUnknown_IsZeroValue(t *testing.T) {
	t.Parallel()
	// Why: the soft-fail path for unknown Compose states relies on
	// StateUnknown being the natural zero value of [ContainerState]
	// so an uninitialized field never accidentally classifies as a
	// known state. M6 slice's dead-allowlist rationale.
	var zero domain.ContainerState
	if zero != domain.StateUnknown {
		t.Errorf("zero-value ContainerState = %v, want StateUnknown", zero)
	}
}
