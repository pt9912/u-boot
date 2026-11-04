package cli_test

import (
	"strings"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// LH-FA-DEV-012: the warning of `devcontainer.envfile` is a plain warn item
// (neither SecurityRelaxation nor DegradationFallback): the full report
// names it with its hint, `--quiet` drops it and exits 0.
func TestDoctor_EnvFile_QuietDropsWarning(t *testing.T) {
	item := domain.Diagnostic{
		ID: "devcontainer.envfile", Severity: domain.SeverityWarn,
		Message: "devcontainer.json passes --env-file .devcontainer/devcontainer.env, but the file is missing; starting the container fails.",
		Hint:    "Run `u-boot generate devcontainer` to create .devcontainer/devcontainer.env.",
	}
	full, err := runDoctorWith(t, []domain.Diagnostic{item}, "doctor")
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if !strings.Contains(full, "devcontainer.envfile") || !strings.Contains(full, "u-boot generate devcontainer") {
		t.Errorf("the full report lacks the warning or its hint:\n%s", full)
	}
	quiet, err := runDoctorWith(t, []domain.Diagnostic{item}, "--quiet", "doctor")
	if err != nil || quiet != "" {
		t.Errorf("--quiet: err = %v, stdout = %q, want silence and exit 0", err, quiet)
	}
}
