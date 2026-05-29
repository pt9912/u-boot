package driving_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestDownSentinels_Identity(t *testing.T) {
	t.Parallel()
	if driving.ErrConfirmationRequired == nil {
		t.Fatal("ErrConfirmationRequired is nil")
	}
	if got, want := driving.ErrConfirmationRequired.Error(), "confirmation required"; got != want {
		t.Errorf("ErrConfirmationRequired.Error() = %q, want %q", got, want)
	}
}

func TestDownSentinels_DistinctFromOthers(t *testing.T) {
	t.Parallel()
	// Why: the §254 destructive-abort sentinel must not be
	// confused with any other driving sentinel. A future refactor
	// that aliased it to ErrComposeFileMissing or
	// ErrProjectNotInitialized would silently retarget exit codes.
	// Use pointer comparison (==) here — both sides are pure
	// sentinel values, so `errors.Is` chain-walking adds no signal
	// and trips staticcheck SA1032 (sentinel-on-left argument
	// ordering).
	others := []error{
		driving.ErrComposeFileMissing,
		driving.ErrStabilizationTimeout,
		driving.ErrProjectNotInitialized,
		driving.ErrServiceUnsupported,
		driving.ErrServiceInconsistent,
	}
	for _, other := range others {
		if driving.ErrConfirmationRequired == other {
			t.Errorf("ErrConfirmationRequired aliased to %v", other)
		}
	}
}

func TestDownSentinels_SurviveContextualWrap(t *testing.T) {
	t.Parallel()
	// Why: same wrap-contract pin as up_test.go — application
	// service wraps with fmt.Errorf("%w", ...); CLI mapping uses
	// errors.Is. Verify that ErrConfirmationRequired survives the
	// expected wrap chain.
	wrapped := fmt.Errorf("down service: confirmation refused: %w", driving.ErrConfirmationRequired)
	if !errors.Is(wrapped, driving.ErrConfirmationRequired) {
		t.Errorf("errors.Is(wrapped, ErrConfirmationRequired) = false, want true")
	}
}

func TestDownRequest_FieldsAreOrthogonal(t *testing.T) {
	t.Parallel()
	// Why: the M6 slice's §T5 truth table treats RemoveVolumes,
	// AssumeYes, NonInteractive as three independent inputs. Pin
	// that all eight combinations are representable in the request
	// struct — a future refactor that collapsed AssumeYes and
	// NonInteractive into a single tri-state would silently break
	// the LH-FA-CLI-005A §235 vs. §254 distinction.
	combinations := []driving.DownRequest{
		{RemoveVolumes: false, AssumeYes: false, NonInteractive: false},
		{RemoveVolumes: false, AssumeYes: false, NonInteractive: true},
		{RemoveVolumes: false, AssumeYes: true, NonInteractive: false},
		{RemoveVolumes: false, AssumeYes: true, NonInteractive: true},
		{RemoveVolumes: true, AssumeYes: false, NonInteractive: false},
		{RemoveVolumes: true, AssumeYes: false, NonInteractive: true},
		{RemoveVolumes: true, AssumeYes: true, NonInteractive: false},
		{RemoveVolumes: true, AssumeYes: true, NonInteractive: true},
	}
	// Pin uniqueness: every combination must produce a unique
	// (bool, bool, bool) tuple — i.e. the three fields don't
	// alias each other.
	seen := make(map[[3]bool]struct{}, len(combinations))
	for _, c := range combinations {
		key := [3]bool{c.RemoveVolumes, c.AssumeYes, c.NonInteractive}
		if _, dup := seen[key]; dup {
			t.Errorf("DownRequest combination collapsed: %+v", c)
		}
		seen[key] = struct{}{}
	}
	if len(seen) != 8 {
		t.Errorf("expected 8 distinct combinations, got %d", len(seen))
	}
}

func TestDownResponse_RemovedVolumesEcho(t *testing.T) {
	t.Parallel()
	// Why: the response field naming is part of the contract — the
	// M6 slice §T1 explicitly removed Stop/Remove counters in favor
	// of this single boolean echo. A rename or type-flip would
	// break the CLI's one-line success-message rendering.
	if r := (driving.DownResponse{RemovedVolumes: true}); !r.RemovedVolumes {
		t.Errorf("DownResponse{RemovedVolumes: true}.RemovedVolumes = false, want true")
	}
	if r := (driving.DownResponse{RemovedVolumes: false}); r.RemovedVolumes {
		t.Errorf("DownResponse{RemovedVolumes: false}.RemovedVolumes = true, want false")
	}
}
