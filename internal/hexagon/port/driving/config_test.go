package driving_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestConfigT2Sentinels_Identity pins the two sentinels added in
// slice-v1-cli-json-dry-run-config T2 (T0-(m) three-class split of
// the formerly overloaded ErrConfigValueInvalid). Identity = exists,
// non-nil, expected message. A renamed message would silently break
// the §6.9 consumer-disambiguation doc.
func TestConfigT2Sentinels_Identity(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"ErrConfigWriteRejected", driving.ErrConfigWriteRejected, "config: write rejected for non-writable path"},
		{"ErrConfigPostPatchSanityFailed", driving.ErrConfigPostPatchSanityFailed, "config: post-patch sanity check failed"},
	}
	for _, c := range cases {
		if c.err == nil {
			t.Fatalf("%s is nil", c.name)
		}
		if got := c.err.Error(); got != c.want {
			t.Errorf("%s.Error() = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestConfigT2Sentinels_DistinctFromValueInvalid pins that the two
// new sentinels are NOT aliased to ErrConfigValueInvalid (nor to
// each other). The whole point of T0-(m) is that a JSON consumer
// can disambiguate the three classes by `code`; an accidental alias
// would silently re-merge them and retarget the mapper rows
// (T0-(f) rows 3, 5, 6 all carry exit 10 — only the sentinel
// identity keeps them apart).
func TestConfigT2Sentinels_DistinctFromValueInvalid(t *testing.T) {
	t.Parallel()
	// Pointer comparison (==): all operands are pure sentinel
	// values, so errors.Is chain-walking adds no signal and trips
	// staticcheck SA1032 (sentinel-on-left ordering).
	all := []error{
		driving.ErrConfigValueInvalid,
		driving.ErrConfigWriteRejected,
		driving.ErrConfigPostPatchSanityFailed,
		driving.ErrConfigSchemaInvalid,
		driving.ErrConfigPathUnknown,
		driving.ErrConfigValueNotSet,
		driving.ErrConfigFileSystem,
	}
	for i := range all {
		for j := i + 1; j < len(all); j++ {
			if all[i] == all[j] {
				t.Errorf("config sentinel %v aliased to %v", all[i], all[j])
			}
		}
	}
}

// TestConfigT2Sentinels_SurviveContextualWrap pins the wrap
// contract: T3 will wrap these sentinels with fmt.Errorf("%w", …)
// in application/config.go and the CLI mapper branches via
// errors.Is. Verify both survive the expected wrap chain.
func TestConfigT2Sentinels_SurviveContextualWrap(t *testing.T) {
	t.Parallel()
	wr := fmt.Errorf("config service: %w", driving.ErrConfigWriteRejected)
	if !errors.Is(wr, driving.ErrConfigWriteRejected) {
		t.Errorf("errors.Is(wrapped, ErrConfigWriteRejected) = false, want true")
	}
	sp := fmt.Errorf("config service: %w", driving.ErrConfigPostPatchSanityFailed)
	if !errors.Is(sp, driving.ErrConfigPostPatchSanityFailed) {
		t.Errorf("errors.Is(wrapped, ErrConfigPostPatchSanityFailed) = false, want true")
	}
}

// TestConfigSetRequest_T2Fields pins the two T2-added request
// fields' types and zero values. PreviewMode shares the
// cluster-wide modifying-subcommand contract (default PreviewNone =
// today's production write); SilenceLogger defaults false (today's
// stderr logging unchanged). Compile-time + zero-value pin: a
// rename or type change fails the build.
func TestConfigSetRequest_T2Fields(t *testing.T) {
	t.Parallel()
	var req driving.ConfigSetRequest

	if req.PreviewMode != driving.PreviewNone {
		t.Errorf("PreviewMode zero-value: want PreviewNone, got %v", req.PreviewMode)
	}
	req.PreviewMode = driving.PreviewDryRun
	if req.PreviewMode != driving.PreviewDryRun {
		t.Errorf("PreviewMode assignment: want PreviewDryRun, got %v", req.PreviewMode)
	}

	if req.SilenceLogger {
		t.Error("SilenceLogger zero-value: want false (today's stderr logging unchanged)")
	}
	req.SilenceLogger = true
	if !req.SilenceLogger {
		t.Error("SilenceLogger assignment to true did not stick")
	}
}

// TestConfigSetResponse_WarningsField pins the T2-added Warnings
// field (T0-(n) Orphan-Feature-WARN-Migration). Type must be
// []driving.WarningEntry so T5 can map it to the envelope's
// diagnostics[] array exactly like remove/up. Nil on the happy
// path; appendable.
func TestConfigSetResponse_WarningsField(t *testing.T) {
	t.Parallel()
	var resp driving.ConfigSetResponse
	if resp.Warnings != nil {
		t.Errorf("Warnings zero-value: want nil, got %v", resp.Warnings)
	}
	resp.Warnings = append(resp.Warnings, driving.WarningEntry{
		Code:    "LH-FA-DEV-003",
		Level:   "warn",
		Message: "feature activated but service not registered",
	})
	if len(resp.Warnings) != 1 || resp.Warnings[0].Level != "warn" {
		t.Errorf("Warnings append did not stick: %+v", resp.Warnings)
	}
}
