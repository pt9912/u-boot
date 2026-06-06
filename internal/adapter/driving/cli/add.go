package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// addFlags bundles the per-invocation flag state of
// `u-boot add <service>` (--yes / --no-interactive read through
// from the root command, plus the add-specific --with-deps from
// LH-FA-ADD-006). The slice-v1-cli-json-dry-run-add slice extends
// this with --dry-run / --diff (LH-FA-CLI-007/008) and a read-through
// of the root --json flag, so [runAdd] sees the full preview-mode
// state in one struct.
type addFlags struct {
	Yes           bool
	NoInteractive bool
	WithDeps      bool
	DryRun        bool
	Diff          bool
	JSON          bool
}

// newAddCommand builds the `u-boot add <service>` Cobra subcommand.
//
// Positional argument layout: `add <service>` instead of fixed
// sub-subcommands (`add postgres`, `add keycloak`, …) so a typo or
// an unsupported service name reaches [domain.NewServiceName] and the
// in-app catalogue check. Either rejection then maps to exit-code 10
// via [ExitCode]; a fixed sub-subcommand structure would surface
// `unknown command` as exit-code 2, hiding the fachliche distinction
// between "invalid name" and "name not in catalogue".
//
// The persistent --yes / --no-interactive / --json flags are read
// through from the root command; today --yes / --no-interactive are
// no-op for the postgres MVP (no per-service interactive prompts),
// but the mutual-exclusion check still fires here so future add-ons
// inherit the same usage contract.
func newAddCommand(a *App) *cobra.Command {
	flags := &addFlags{}

	cmd := &cobra.Command{
		Use:   "add <service>",
		Short: "Add a service add-on to the u-boot project",
		Long: `Add an integrated service add-on to the current u-boot project. The
positional <service> argument names the catalogue entry; today MVP
supports only postgres (LH-FA-ADD-002). Keycloak (LH-FA-ADD-003) and
OpenTelemetry (LH-FA-ADD-004) join in V1.

The command runs in an initialized project (u-boot.yaml present per
LH-FA-ADD-001) and is idempotent: calling it twice on the same active
service is a no-op. The state-machine (LH-FA-ADD-005) handles
re-activation, block rebuild after a partial cleanup, and abort with
a repair hint when the project state is inconsistent (orphan compose
block or user-managed entry without u-boot marker).

Flag combinations (LH-FA-CLI-007/008):
  --dry-run            preview without writing files
  --diff               show unified diff of planned changes
  --dry-run --diff     unified diff preview, no write
  --json               JSON output; pairs with --dry-run / --diff
                       for the LH-FA-CLI-007 §326 voll-schema

Examples:
  u-boot add postgres                 # first add: register + write
  u-boot add postgres                 # idempotent re-run: no-op
  u-boot add postgres --dry-run       # preview, no write
  u-boot add postgres --diff --json   # voll-schema with hunks
  u-boot add redis                    # exit 10 — not in catalogue
  u-boot add keycloak --with-deps     # auto-install missing deps`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			flags.JSON = a.json
			return runAdd(cmd.Context(), cmd.OutOrStdout(), args, *flags, a.addServiceUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.WithDeps, "with-deps", false,
		"auto-install missing add-on dependencies (LH-FA-ADD-006) without prompting")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"preview the planned changes without writing files (LH-FA-CLI-007)")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"render a unified diff of the planned changes (LH-FA-CLI-008)")

	return cmd
}

// runAdd is split from the Cobra closure for direct unit-testing
// (no Cobra command construction needed). ctx is taken as the first
// parameter explicitly so contextcheck can see the propagation.
//
// The function performs the LH-FA-CLI-005A mutual-exclusion check on
// --yes / --no-interactive (same shape as runInit), validates the
// positional service name via [domain.NewServiceName] (LH-FA-ADD-001
// invalid-name path), then delegates to the AddServiceUseCase. The
// use case owns the LH-FA-ADD-005 dispatch — runAdd only translates
// the response into output. Output shape is mode-dependent:
//
//   - --json without --dry-run/--diff → minimal envelope (Spec §1841
//     / T0-(k)). FS-mutations happen as in normal mode; the envelope
//     carries no plan information.
//   - --dry-run --json (with or without --diff) → voll-schema with
//     plannedFiles[] from the recorder, dryRun=true, no FS-write.
//   - --diff --json without --dry-run → voll-schema preview-and-apply,
//     plannedFiles[] + hunks, dryRun=false, diff=true.
//   - Human-mode --diff (with or without --dry-run) → unified-diff
//     string on stdout via [diff.Render].
//   - Human-mode plain → existing [printAddSummary] behaviour.
func runAdd(
	ctx context.Context,
	out io.Writer,
	args []string,
	flags addFlags,
	uc driving.AddServiceUseCase,
	getwd func() (string, error),
) error {
	// mapErr-Source-Pflicht (slice-v1-cli-json-dry-run-init T0-(e)):
	// jeder Subcommand-RunE definiert seine eigene mapErr-Funktion
	// und reicht sie an reportError weiter. Symmetrie zu init's
	// mapInitErrorToDiagnostic.
	mapErr := mapAddErrorToDiagnostic

	if flags.Yes && flags.NoInteractive {
		return reportError(out, ErrConflictingModeFlags, nil, flags.DryRun, flags.Diff, flags.JSON, "add", mapErr)
	}

	svcName, err := domain.NewServiceName(args[0])
	if err != nil {
		return reportError(out, err, nil, flags.DryRun, flags.Diff, flags.JSON, "add", mapErr)
	}

	cwd, err := getwd()
	if err != nil {
		return reportError(out, fmt.Errorf("determine working directory: %w", err), nil, flags.DryRun, flags.Diff, flags.JSON, "add", mapErr)
	}

	mode := previewModeFromFlags(flags.DryRun, flags.Diff)
	resp, addErr := uc.Add(ctx, driving.AddServiceRequest{
		BaseDir:       cwd,
		ServiceName:   svcName,
		WithDeps:      flags.WithDeps,
		Yes:           flags.Yes,
		NoInteractive: flags.NoInteractive,
		PreviewMode:   mode,
	})

	if addErr != nil {
		return reportError(out, addErr, resp.PlannedFiles, flags.DryRun, flags.Diff, flags.JSON, "add", mapErr)
	}

	if flags.JSON {
		return writeAddJSON(out, resp, flags.DryRun, flags.Diff)
	}

	if flags.Diff {
		if err := writeDiff(out, resp.PlannedFiles); err != nil {
			return err
		}
	}
	return printAddSummary(out, resp, flags.DryRun)
}

// writeAddJSON renders the success-path JSON envelope. Three shapes
// per T0-(k) + T0-(b)/(d):
//
//   - dryRun=false && diff=false → minimal envelope (Spec §1841).
//   - dryRun=true                → voll-schema, plannedFiles from
//     recorder, optional hunks if diff=true.
//   - diff=true                  → voll-schema preview-and-apply,
//     plannedFiles + hunks.
func writeAddJSON(out io.Writer, resp driving.AddServiceResponse, dryRun, diffFlag bool) error {
	if !dryRun && !diffFlag {
		env := newMinimalEnvelope("add", "", nil, 0)
		return writeEnvelope(out, env)
	}
	pfs, chs := mapPlannedFilesToWire(resp.PlannedFiles, diffFlag)
	env := newFullEnvelope("add", "", dryRun, diffFlag, pfs, chs, nil, 0)
	return writeEnvelope(out, env)
}

// writeAddErrorEnvelope/lastPlannedPath wurden in slice-v1-cli-json-
// dry-run-init T5-a nach `cli/erroremission.go` extrahiert als
// `writeErrorEnvelope`/`lastPlannedPath` mit decomposed-Slices-
// Signatur. runAdd ruft jetzt `reportError(out, err, planned, ...,
// "add", mapAddErrorToDiagnostic)`.

// writeEnvelope marshals env and writes it with a trailing newline.
// Centralised so all three add JSON paths share the same I/O shape.
func writeEnvelope(out io.Writer, env cliJSONEnvelope) error {
	raw, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal add envelope: %w", err)
	}
	if _, err := fmt.Fprintln(out, string(raw)); err != nil {
		return fmt.Errorf("write add envelope: %w", err)
	}
	return nil
}

// mapResponseToWire/computeChangeCountAndHunks/toCLIHunks wurden in
// slice-v1-cli-json-dry-run-init T1-D nach
// `internal/adapter/driving/cli/wireshapes.go` extrahiert. add ruft
// jetzt mapPlannedFilesToWire(resp.PlannedFiles, diffFlag).

// mapAddErrorToDiagnostic maps an add-path error to a diagnosticItem
// with the spec-konforme LH-Kennung per T0-(j). Unknown errors fall
// back to a generic LH-FA-CLI-006 wrapper (default error path); the
// invariants Spec §1834 (level ∈ {warn, error}) and §1837 (status
// coupling) carry the rest.
//
// Order matters: addservice_execute.go wraps FS-Write-Failures as
// `fmt.Errorf("write %s: %w: %w", path, ErrAddFileSystem, rawErr)` —
// a multi-%w wrap (Go 1.20+). If a future code path adds a fachlich
// sentinel to the same chain (e.g. atomicity rollback also wrapping
// ErrServiceInconsistent), errors.Is would match BOTH this case and
// ErrAddFileSystem. ErrAddFileSystem checks first (review #11) so
// the user sees the technical-persistence diagnostic and exit-14
// classification rather than a misleading fachlich code (LH-FA-ADD-005
// → exit 10) that the wrap accidentally included.
//
// Renamed from mapErrorToDiagnostic in slice-v1-cli-json-dry-run-
// init T5-a (T0-(e) mapErr-Source-Pflicht) — init parallel definiert
// mapInitErrorToDiagnostic, beide werden per Function-Value an
// reportError gereicht.
func mapAddErrorToDiagnostic(err error) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrAddFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-001", Message: err.Error()}
	case errors.Is(err, driving.ErrServiceUnsupported):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-002", Message: err.Error()}
	case errors.Is(err, driving.ErrServiceInconsistent):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-005", Message: err.Error()}
	case errors.Is(err, driving.ErrDependenciesRequired):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-006", Message: err.Error()}
	case errors.Is(err, domain.ErrInvalidServiceName):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-006", Message: err.Error()}
	case errors.Is(err, driving.ErrFileExists), errors.Is(err, driving.ErrProjectExists):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-004", Message: err.Error()}
	case errors.Is(err, driving.ErrBackupSuffixExhausted), errors.Is(err, driving.ErrBackupSourceMissing):
		// Defensive branch — add ruft heute keine Backup-Logik
		// (kein runBackup/BackupPath im add-Use-Case-Pfad). Branch
		// bleibt für zukünftige Catalog-Erweiterungen, die Backup
		// rufen könnten, und wird auf LH-NFA-REL-003 + Exit 14
		// klassifiziert — analog mapInitErrorToDiagnostic. Vorher
		// fälschlich LH-FA-INIT-005 (Validation-Klasse, würde Exit 10
		// suggerieren); isFilesystemError routet ohnehin zu Exit 14,
		// also war Envelope-Code und Exit-Klasse desynchron
		// (slice-v1-cli-cleanup-add-backup-error-class).
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// writeAddDiff wurde in slice-v1-cli-json-dry-run-init T5-a nach
// `cli/erroremission.go` als generisches `writeDiff(out, planned)`
// extrahiert. runAdd ruft jetzt writeDiff(out, resp.PlannedFiles).

// printAddSummary writes a short, deterministic summary of the add
// outcome. dryRun=true switches the lead-in from "Added" to "Would
// add" so the user sees that no FS-write happened. Three shapes:
//
//   - PriorState=Active and Changed=nil: prints
//     "Service <name> is already active; no changes."
//   - PriorState=Active and Changed!=nil: prints
//     "Repaired service <name> artefacts." plus the changed paths.
//   - Other PriorState transitions: prints
//     "Added service <name>." plus the changed paths.
//
// Order of Changed follows the AddServiceUseCase contract
// (u-boot.yaml → compose.yaml → .env.example), so the user sees a
// stable rollback hint without re-sorting.
func printAddSummary(out io.Writer, resp driving.AddServiceResponse, dryRun bool) error {
	header := addSummaryHeader(resp, dryRun)
	if _, err := fmt.Fprint(out, header); err != nil {
		return err
	}
	if resp.PriorState == resp.State && len(resp.Changed) == 0 {
		// Header already carries the full "already active" line.
		return nil
	}
	for _, p := range resp.Changed {
		if _, err := fmt.Fprintln(out, "  - "+p); err != nil {
			return err
		}
	}
	return nil
}

// addSummaryHeader picks the lead-in text for printAddSummary. The
// three-state switch lives here so the printer stays linear (the
// linter's cognitive-complexity budget caps the printer at ≤ 20).
func addSummaryHeader(resp driving.AddServiceResponse, dryRun bool) string {
	name := resp.ServiceName.String()
	switch {
	case resp.PriorState == resp.State && len(resp.Changed) == 0:
		return fmt.Sprintf("Service %q is already active; no changes.\n", name)
	case resp.PriorState == resp.State && len(resp.Changed) > 0:
		if dryRun {
			return fmt.Sprintf("Would repair service %q artefacts.\n\nWould change:\n", name)
		}
		return fmt.Sprintf("Repaired service %q artefacts.\n\nChanged:\n", name)
	default:
		if dryRun {
			return fmt.Sprintf("Would add service %q.\n\nWould change:\n", name)
		}
		return fmt.Sprintf("Added service %q.\n\nChanged:\n", name)
	}
}
