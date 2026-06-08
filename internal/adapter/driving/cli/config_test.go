package cli_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli"
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
