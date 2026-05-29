package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// doctorFlags bundles the per-invocation flag state of
// `u-boot doctor` that runDoctor actually consumes today: --strict
// (LH-FA-DIAG-003) plus --quiet (LH-FA-CLI-005, filters OK items
// from the rendered report). The persistent --verbose / --debug
// flags exist on the root command per LH-FA-CLI-005 and are
// load-bearing at the LOGGER level since
// `slice-followup-verbosity-wiring` (`buildRootCommand`'s
// PersistentPreRunE flips a shared *slog.LevelVar). The doctor
// renderer itself does not consume them — service-level
// logger.Debug/Info calls are the surface they govern.
type doctorFlags struct {
	Strict bool
	Quiet  bool
}

// idColumnWidth is the padding used for the diagnostic-ID column in
// the rendered report. Sized for the longest ID today
// (`devcontainer.dockerfile.valid`, 29 chars) plus one space of
// margin; bump when a future check introduces a longer ID.
const idColumnWidth = 30

// newDoctorCommand builds the `u-boot doctor` Cobra subcommand
// (LH-FA-DIAG-001).
//
// Local flag:
//
//	--strict        treat any Warn as fail-grade — exit code 11 fires
//	                even without an Error (LH-FA-DIAG-003).
//
// The persistent flags --quiet / --verbose / --debug (LH-FA-CLI-005)
// are read from the App after Cobra parses them.
func newDoctorCommand(a *App) *cobra.Command {
	flags := &doctorFlags{}

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose the local environment for u-boot prerequisites",
		Long: `Run a battery of checks against the current directory and the local
environment, classify each finding as ok / warn / error per
LH-FA-DIAG-003, and exit with code 11 when any error is present
(or any warn, when --strict is set).

Default checks (LH-FA-DIAG-002):
  - write-permissions in the current directory
  - git binary availability
  - docker binary + version (≥ 24.0) + daemon reachability
  - docker compose plugin + version (≥ 2.20)
  - u-boot.yaml syntax + schemaVersion + project.name regex
  - compose.yaml syntax + top-level services-shape
  - .devcontainer/devcontainer.json syntax + minimum compat (when present)
  - .devcontainer/Dockerfile FROM-directive (when present)

Examples:
  u-boot doctor                  # full report on stdout
  u-boot doctor --quiet          # hide ok entries
  u-boot doctor --strict         # warn → exit 11 (CI-strict mode)`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			flags.Quiet = a.quiet
			return runDoctor(cmd.Context(), cmd.OutOrStdout(), *flags, a.doctorUseCase, a.getwd)
		},
	}
	cmd.Flags().BoolVar(&flags.Strict, "strict", false,
		"treat any warning as fail-grade (LH-FA-DIAG-003)")
	return cmd
}

// runDoctor is split from the Cobra closure for testability with
// fake use-case + fake getwd. Context is the first parameter so
// contextcheck can see the propagation (the closure boundary itself
// is `contextcheck`-excluded for this package; see
// `.golangci.yml`).
func runDoctor(
	ctx context.Context,
	out io.Writer,
	flags doctorFlags,
	uc driving.DoctorUseCase,
	getwd func() (string, error),
) error {
	cwd, err := getwd()
	if err != nil {
		return fmt.Errorf("determine working directory: %w", err)
	}
	resp, err := uc.Check(ctx, driving.DoctorRequest{BaseDir: cwd})
	if err != nil {
		return err
	}

	writeDoctorReport(out, cwd, resp.Report, flags.Quiet)

	// Exit-code policy (LH-FA-DIAG-003):
	//   - any Error              → ErrDoctorFailures
	//   - --strict + any Warn    → ErrDoctorFailures
	//   - otherwise              → nil (exit 0)
	if resp.Report.HasErrors() {
		return ErrDoctorFailures
	}
	if flags.Strict && resp.Report.HasWarnings() {
		return ErrDoctorFailures
	}
	return nil
}

// writeDoctorReport renders the diagnostic report on out. Items are
// sorted Errors-first then Warns then OKs (via
// [domain.DiagnosticReport.SortedByIssuesFirst]); the OK band is
// suppressed when `quiet` is true. Hints render only on warn/error
// items (OK items rarely carry a hint and listing one would clutter
// the report).
func writeDoctorReport(out io.Writer, baseDir string, report domain.DiagnosticReport, quiet bool) {
	fmt.Fprintf(out, "Diagnostic report for %s\n", baseDir)
	fmt.Fprintln(out, "──────────────────────────────────────")

	items := report.SortedByIssuesFirst()
	var nOK, nWarn, nErr int
	for _, item := range items {
		switch item.Severity {
		case domain.SeverityOK:
			nOK++
		case domain.SeverityWarn:
			nWarn++
		case domain.SeverityError:
			nErr++
		}
		if quiet && item.Severity == domain.SeverityOK {
			continue
		}
		fmt.Fprintf(out, "%s  %-*s %s\n", severityGlyph(item.Severity), idColumnWidth, item.ID, item.Message)
		if item.Hint != "" && item.Severity != domain.SeverityOK {
			fmt.Fprintf(out, "   → %s\n", item.Hint)
		}
	}

	fmt.Fprintf(out, "\nSummary: %d error, %d warn, %d ok\n", nErr, nWarn, nOK)
}

// severityGlyph returns the Unicode glyph that prefixes each
// diagnostic line in the rendered report. Pure ASCII fallbacks are
// not provided — the project's other adapters (e.g. the init
// summary in [printInitSummary]) also use Unicode arrows / dashes,
// and the test pipeline runs in UTF-8-clean containers.
func severityGlyph(s domain.Severity) string {
	switch s {
	case domain.SeverityError:
		return "✗"
	case domain.SeverityWarn:
		return "⚠"
	case domain.SeverityOK:
		return "✓"
	}
	return "?"
}
