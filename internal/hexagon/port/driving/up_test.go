package driving_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

func TestUpSentinels_Identity(t *testing.T) {
	t.Parallel()
	// Why: pin that each Up sentinel is a non-nil unique error value
	// so `errors.Is` at the CLI adapter (M6-T7 mapping table) sees
	// the expected identities. A drift (e.g. accidental reassignment
	// in another file) would silently turn an exit-code-10 path into
	// an exit-code-1 path.
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"ErrComposeFileMissing", driving.ErrComposeFileMissing, "compose file missing"},
		{"ErrStabilizationTimeout", driving.ErrStabilizationTimeout, "stabilization timeout"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.err == nil {
				t.Fatalf("%s is nil", tc.name)
			}
			if got := tc.err.Error(); got != tc.want {
				t.Errorf("%s.Error() = %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestUpSentinels_MutuallyDistinct(t *testing.T) {
	t.Parallel()
	// Why: a future refactor that accidentally aliased one sentinel
	// to another would make the exit-code mapping silently merge two
	// classes (e.g. stabilization timeout becomes the same as
	// compose-file missing). Pointer comparison (==) is the precise
	// test here — both sides are pure sentinel values, so errors.Is
	// adds no signal and trips staticcheck SA1032.
	pairs := []struct {
		a, b error
	}{
		{driving.ErrComposeFileMissing, driving.ErrStabilizationTimeout},
		{driving.ErrComposeFileMissing, driving.ErrProjectNotInitialized},
		{driving.ErrStabilizationTimeout, driving.ErrProjectNotInitialized},
	}
	for _, p := range pairs {
		if p.a == p.b {
			t.Errorf("sentinels aliased: %v == %v", p.a, p.b)
		}
	}
}

func TestUpSentinels_SurviveContextualWrap(t *testing.T) {
	t.Parallel()
	// Why: pin the application-level wrap contract from the M6 slice
	// §Sentinel-Schichtung — wrapping a sentinel under fmt.Errorf
	// with %w must leave errors.Is intact. M6 T4 (UpService) wraps
	// engine and validation errors under contextual prefixes; if
	// %w → %v slipped in anywhere, exit-code mapping at the CLI
	// would silently lose the sentinel identity.
	wrapped := fmt.Errorf("up service: ComposeUp on %q: %w", "/tmp/demo", driving.ErrComposeFileMissing)
	if !errors.Is(wrapped, driving.ErrComposeFileMissing) {
		t.Errorf("errors.Is(wrapped, ErrComposeFileMissing) = false, want true")
	}
	doubleWrapped := fmt.Errorf("cli: %w", wrapped)
	if !errors.Is(doubleWrapped, driving.ErrComposeFileMissing) {
		t.Errorf("errors.Is(doubleWrapped, ErrComposeFileMissing) = false, want true")
	}
}

func TestUpRequest_TimeoutSemanticsAreValueBased(t *testing.T) {
	t.Parallel()
	// Why: the M6 slice (and the CLI flag conversion contract) pin
	// UpRequest.Timeout as a time.Duration value — not a raw second
	// integer. Pin the zero-value (fire-and-forget) and the
	// negative-value (validation error) semantics so a future
	// refactor that flipped the type silently breaks the CLI
	// adapter.
	var zero driving.UpRequest
	if zero.Timeout != 0 {
		t.Errorf("zero-value Timeout = %v, want 0 (fire-and-forget)", zero.Timeout)
	}
	// Negative duration is representable; the application layer is
	// responsible for rejecting it before calling Compose. We only
	// pin here that the type can carry it.
	req := driving.UpRequest{Timeout: -1}
	if req.Timeout >= 0 {
		t.Errorf("negative Timeout assignment lost: %v, want < 0", req.Timeout)
	}
}
