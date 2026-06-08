package cli

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// TestMapConfigErrorToDiagnostic_AllRows covers every switch-case in
// mapConfigErrorToDiagnostic (T0-(f) Mapper-Tabelle) so each
// config-sentinel maps to its spec-konforme LH-Kennung. White-box
// (package cli) because the mapper is unexported.
//
// LH-FA-CONF-005 is deliberately multi-use (Path-Unknown /
// Write-Rejected / Value-Not-Set, T0-(m)/R3-MED-2) — three rows
// share the code; consumers disambiguate by message, not code.
func TestMapConfigErrorToDiagnostic_AllRows(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		err    error
		wantLH string
	}{
		{"ErrConfigFileSystem", driving.ErrConfigFileSystem, "LH-NFA-REL-003"},
		{"ErrConfigSchemaInvalid", driving.ErrConfigSchemaInvalid, "LH-FA-CONF-002"},
		{"ErrConfigPostPatchSanityFailed", driving.ErrConfigPostPatchSanityFailed, "LH-FA-CONF-002"},
		{"ErrConfigPathUnknown", driving.ErrConfigPathUnknown, "LH-FA-CONF-005"},
		{"ErrConfigWriteRejected", driving.ErrConfigWriteRejected, "LH-FA-CONF-005"},
		{"ErrConfigValueInvalid", driving.ErrConfigValueInvalid, "LH-FA-CONF-001"},
		{"ErrConfigValueNotSet", driving.ErrConfigValueNotSet, "LH-FA-CONF-005"},
		{"ErrProjectNotInitialized", driving.ErrProjectNotInitialized, "LH-FA-INIT-001"},
		{"ErrDryRunNotApplicable", ErrDryRunNotApplicable, "LH-FA-CLI-006"},
		{"unknown → default", errors.New("boom"), "LH-FA-CLI-006"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			diag := mapConfigErrorToDiagnostic(tc.err)
			if diag.Code != tc.wantLH {
				t.Errorf("code: want %q, got %q", tc.wantLH, diag.Code)
			}
			if diag.Level != "error" {
				t.Errorf("level: want error, got %q", diag.Level)
			}
			if diag.Message == "" {
				t.Error("message must not be empty")
			}
		})
	}
}

// TestMapConfigErrorToDiagnostic_SwitchOrderFSFirst_ByDesign pins the
// T0-(f) FS-first switch-order: a synthetic multi-`%w` chain with
// BOTH ErrConfigFileSystem AND ErrConfigValueInvalid resolves to
// LH-NFA-REL-003 (FS-class) in the diagnostic code — not LH-FA-CONF-001.
//
// `_ByDesign`: [ExitCode] is an INDEPENDENT classifier that checks
// isValidationError BEFORE isFilesystemError, so the same synthetic
// chain yields exit 10 (validation-sub-class). The (code, exitCode)
// tuple is the disambiguation contract (cli-json-output.md §6.7
// pattern): the FS code signals the class, the exit differentiates
// the sub-sentinel source. No real code path chains both today —
// this is a defense-only robustness pin.
func TestMapConfigErrorToDiagnostic_SwitchOrderFSFirst_ByDesign(t *testing.T) {
	t.Parallel()
	chain := fmt.Errorf("config service: synthetic: %w: %w",
		driving.ErrConfigFileSystem, driving.ErrConfigValueInvalid)

	if diag := mapConfigErrorToDiagnostic(chain); diag.Code != "LH-NFA-REL-003" {
		t.Errorf("switch-order violation: want LH-NFA-REL-003 (FS-first), got %q", diag.Code)
	}
	if !errors.Is(chain, driving.ErrConfigFileSystem) || !errors.Is(chain, driving.ErrConfigValueInvalid) {
		t.Fatal("multi-wrap must match both sentinels via errors.Is")
	}
	// By-design (code, exitCode) split: ExitCode's isValidationError
	// runs before isFilesystemError → 10, even though the code is FS.
	if got := ExitCode(chain); got != 10 {
		t.Errorf("ExitCode by-design: want 10 (validation-sub-class via independent classifier), got %d", got)
	}
}

// TestConfigFileSystem_PureChainExit14 pins the real (non-synthetic)
// FS path: a pure ErrConfigFileSystem wrap (FS sentinel + raw error,
// no validation sentinel) maps to LH-NFA-REL-003 AND exit 14 — the
// production shape of the five Multi-`%w` FS wrap-sites in
// application/config.go.
func TestConfigFileSystem_PureChainExit14(t *testing.T) {
	t.Parallel()
	chain := fmt.Errorf("config service: read %q: %w: %w",
		"u-boot.yaml", driving.ErrConfigFileSystem, errors.New("permission denied"))
	if got := ExitCode(chain); got != 14 {
		t.Errorf("pure FS chain: want exit 14, got %d", got)
	}
	if diag := mapConfigErrorToDiagnostic(chain); diag.Code != "LH-NFA-REL-003" {
		t.Errorf("pure FS chain: want LH-NFA-REL-003, got %q", diag.Code)
	}
}
