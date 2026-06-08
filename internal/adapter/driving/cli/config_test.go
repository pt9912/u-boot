package cli_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestErrDryRunNotApplicable_Identity pins the slice-v1-cli-json-
// dry-run-config T2 CLI sentinel (T0-(g) Option (i.a)): the
// read-only config forms (`config`, `config get`) reject
// `--dry-run`/`--diff` Envelope-konform via this sentinel. T5
// registers the synthetic flags, wires the reject path, and adds
// the isUsageError branch (Exit 2); T2 only defines the sentinel.
func TestErrDryRunNotApplicable_Identity(t *testing.T) {
	t.Parallel()
	if cli.ErrDryRunNotApplicable == nil {
		t.Fatal("ErrDryRunNotApplicable is nil")
	}
	want := "--dry-run/--diff is only valid for `config set`"
	if got := cli.ErrDryRunNotApplicable.Error(); got != want {
		t.Errorf("ErrDryRunNotApplicable.Error() = %q, want %q", got, want)
	}
}

// TestErrDryRunNotApplicable_SurvivesWrap pins the wrap contract:
// T5's runConfigShow/runConfigGet reject path wraps the sentinel
// (Pattern-Erbe logs ErrFollowJSONNotSupported) and isUsageError
// branches via errors.Is to assign Exit 2.
func TestErrDryRunNotApplicable_SurvivesWrap(t *testing.T) {
	t.Parallel()
	wrapped := fmt.Errorf("config get: %w", cli.ErrDryRunNotApplicable)
	if !errors.Is(wrapped, cli.ErrDryRunNotApplicable) {
		t.Errorf("errors.Is(wrapped, ErrDryRunNotApplicable) = false, want true")
	}
}

// TestExitCode_ConfigValidationSentinels pins that every config
// fachlich-validation sentinel maps to Exit-Code 10. The T0-(m)
// split (T3) introduced ErrConfigWriteRejected +
// ErrConfigPostPatchSanityFailed; ExitCode is a classifier
// INDEPENDENT of the T5 mapConfigErrorToDiagnostic mapper, so a
// forgotten entry silently drops the new classes to Exit 1
// (Independent-Review finding R-IR-1). This table guards every
// config Exit-10 sentinel against that regression, wrapped to
// mimic the application layer's fmt.Errorf("%w: …", sentinel).
func TestExitCode_ConfigValidationSentinels(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		err  error
	}{
		{"ErrConfigPathUnknown", driving.ErrConfigPathUnknown},
		{"ErrConfigValueInvalid", driving.ErrConfigValueInvalid},
		{"ErrConfigWriteRejected", driving.ErrConfigWriteRejected},
		{"ErrConfigPostPatchSanityFailed", driving.ErrConfigPostPatchSanityFailed},
		{"ErrConfigSchemaInvalid", driving.ErrConfigSchemaInvalid},
		{"ErrConfigValueNotSet", driving.ErrConfigValueNotSet},
	}
	for _, c := range cases {
		wrapped := fmt.Errorf("config service: %w: detail", c.err)
		if got := cli.ExitCode(wrapped); got != 10 {
			t.Errorf("ExitCode(wrap of %s) = %d, want 10 (LH-FA-CLI-006 validation)", c.name, got)
		}
	}
}
