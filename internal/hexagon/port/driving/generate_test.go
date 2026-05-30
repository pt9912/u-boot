package driving_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestGenerateSentinels_Distinct(t *testing.T) {
	t.Parallel()
	sentinels := []error{
		driving.ErrArtifactUnknown,
		driving.ErrGenerateManualConflict,
		driving.ErrGenerateFileSystem,
	}
	for i, a := range sentinels {
		if a == nil {
			t.Errorf("sentinel #%d is nil", i)
			continue
		}
		for j := i + 1; j < len(sentinels); j++ {
			if a == sentinels[j] {
				t.Errorf("sentinel #%d aliased to #%d", i, j)
			}
		}
	}
	// Cross-check: none of the M7 sentinels collide with the reused
	// M5/M6 ErrProjectNotInitialized. A future refactor that aliased
	// ErrGenerateManualConflict to ErrProjectNotInitialized would
	// silently retarget the exit-code mapping.
	for i, a := range sentinels {
		if a == driving.ErrProjectNotInitialized {
			t.Errorf("M7 sentinel #%d aliased to ErrProjectNotInitialized", i)
		}
	}
}

func TestGenerateSentinels_SurviveContextualWrap(t *testing.T) {
	t.Parallel()
	// Mirror the wrap-contract pin from up_test.go / down_test.go:
	// the application service wraps with fmt.Errorf("%w", ...) and
	// the CLI mapping uses errors.Is. Verify the three M7 sentinels
	// survive the expected wrap chain.
	for _, sentinel := range []error{
		driving.ErrArtifactUnknown,
		driving.ErrGenerateManualConflict,
		driving.ErrGenerateFileSystem,
	} {
		wrapped := fmt.Errorf("generate service: %w", sentinel)
		if !errors.Is(wrapped, sentinel) {
			t.Errorf("errors.Is(wrapped, %v) = false", sentinel)
		}
	}
}

func TestGenerateAction_String(t *testing.T) {
	t.Parallel()
	cases := map[driving.GenerateAction]string{
		driving.GenerateActionCreated:        "created",
		driving.GenerateActionUpdatedBlock:   "updated-block",
		driving.GenerateActionNoOp:           "no-op",
		driving.GenerateActionRepairedManual: "repaired-manual",
		driving.GenerateAction(99):           "unknown",
	}
	for act, want := range cases {
		if got := act.String(); got != want {
			t.Errorf("GenerateAction(%d).String() = %q, want %q", act, got, want)
		}
	}
}
