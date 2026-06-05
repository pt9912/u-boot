package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/diff"
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

// previewModeFromFlags maps the --dry-run / --diff flag combination
// to the [driving.AddPreviewMode] enum per slice-v1-cli-json-dry-run-
// add T0-(b) Wahrheitstabelle:
//
//	--dry-run | --diff | mode              | production write?
//	-----------+--------+-------------------+------------------
//	no        | no    | PreviewNone       | yes (Normal-Mode)
//	yes       | no    | PreviewDryRun     | no  (Plan only)
//	no        | yes   | PreviewAndApply   | yes (Plan + Write)
//	yes       | yes   | PreviewDryRun     | no  (Diff preview)
//
// The --dry-run-wins rule on the (yes, yes) cell matches LH-FA-CLI-
// 007 (dry-run is a hard "no write") combined with LH-FA-CLI-008
// (--diff alone is preview-and-apply). With both flags set the user
// asked for a diff preview WITHOUT writing — that's PreviewDryRun.
func previewModeFromFlags(dryRun, diffFlag bool) driving.AddPreviewMode {
	if dryRun {
		return driving.PreviewDryRun
	}
	if diffFlag {
		return driving.PreviewAndApply
	}
	return driving.PreviewNone
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
	if flags.Yes && flags.NoInteractive {
		// Review #5: --json contract requires every error to produce a
		// JSON envelope. The early-return path used to bypass JSON-mode
		// and ship the raw error to stderr.
		return reportAddError(out, ErrConflictingModeFlags, driving.AddServiceResponse{}, flags)
	}

	svcName, err := domain.NewServiceName(args[0])
	if err != nil {
		return reportAddError(out, err, driving.AddServiceResponse{}, flags)
	}

	cwd, err := getwd()
	if err != nil {
		// Review #6: getwd failure used to skip the JSON envelope.
		return reportAddError(out, fmt.Errorf("determine working directory: %w", err), driving.AddServiceResponse{}, flags)
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
		return reportAddError(out, addErr, resp, flags)
	}

	if flags.JSON {
		return writeAddJSON(out, resp, flags.DryRun, flags.Diff)
	}

	if flags.Diff {
		if err := writeAddDiff(out, resp); err != nil {
			return err
		}
	}
	return printAddSummary(out, resp, flags.DryRun)
}

// reportAddError is the single error-emission gate of runAdd: in
// human mode it surfaces the raw error to Cobra (which lets main.go
// render it to stderr and compute the exit code); in JSON mode it
// writes the envelope to stdout AND still returns the original error
// so cli.ExitCode picks up the right exit code (review #2 — without
// the propagation the shell would see 0 while the envelope claims
// e.g. 14).
func reportAddError(out io.Writer, addErr error, resp driving.AddServiceResponse, flags addFlags) error {
	if !flags.JSON {
		return addErr
	}
	if err := writeAddErrorEnvelope(out, addErr, resp, flags.DryRun, flags.Diff); err != nil {
		return err
	}
	return addErr
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
	pfs, chs := mapResponseToWire(resp, diffFlag)
	env := newFullEnvelope("add", "", dryRun, diffFlag, pfs, chs, nil, 0)
	return writeEnvelope(out, env)
}

// writeAddErrorEnvelope renders the JSON envelope on the error path.
// The use case returns a non-empty Response (PlannedFiles populated
// up to the failure point — T0-(b) Mid-Write-Failure / T0-(j)
// Round-4 H2), so the envelope still ships a voll-schema view when
// --dry-run / --diff was requested. Validation errors (Pre-Write,
// e.g. invalid service name) ship the minimal envelope shape since
// no plan was made.
//
// dryRun and diffFlag are forwarded from the user's flag state
// (review #4): the previous form hardcoded `dryRun=false, diff=true`
// which misrepresented `--dry-run --diff --json` invocations as
// preview-and-apply, contradicting Spec §326 voll-schema fields.
// The envelope now always reports the actual user-requested mode.
func writeAddErrorEnvelope(out io.Writer, addErr error, resp driving.AddServiceResponse, dryRun, diffFlag bool) error {
	diag := mapErrorToDiagnostic(addErr)
	exitCode := ExitCode(addErr)
	// Annotate the diagnostic with the failure path when the
	// application layer carries one (Mid-Write-Failure surfaces the
	// failing path via the resp.PlannedFiles tail entry); not all
	// error classes know a path, in which case `file` stays empty.
	if path := lastPlannedPath(resp); path != "" {
		diag.File = path
	}
	// Voll-schema applies whenever the recorder captured anything OR
	// the user explicitly asked for it via --dry-run/--diff. Without
	// a recorder capture and without a preview flag the envelope
	// shape is the minimal contract (Spec §1841).
	wantsFullSchema := len(resp.PlannedFiles) > 0 || dryRun || diffFlag
	if !wantsFullSchema {
		env := newMinimalEnvelope("add", "", []diagnosticItem{diag}, exitCode)
		return writeEnvelope(out, env)
	}
	pfs, chs := mapResponseToWire(resp, diffFlag)
	env := newFullEnvelope("add", "", dryRun, diffFlag, pfs, chs, []diagnosticItem{diag}, exitCode)
	return writeEnvelope(out, env)
}

// lastPlannedPath returns the path of the last PlannedFile in the
// response — convenient for Mid-Write-Failure-Diagnostics where the
// recorder's tail entry is the failing path.
func lastPlannedPath(resp driving.AddServiceResponse) string {
	if len(resp.PlannedFiles) == 0 {
		return ""
	}
	return resp.PlannedFiles[len(resp.PlannedFiles)-1].Path
}

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

// mapResponseToWire converts the driving-layer recorder capture into
// the CLI wire-types. The Hunks-field is populated only when
// withHunks is true (--diff / preview-and-apply); `changes[].count`
// always follows the T0-(g) semantics regardless of flag state —
// for modify-actions that means we compute hunks even without --diff
// just to sum their NewLines, since the alternative (CountLines on
// the whole new file) overstated the count by orders of magnitude
// for any add-on-into-existing-file case (review #1).
func mapResponseToWire(resp driving.AddServiceResponse, withHunks bool) ([]plannedFile, []changeEntry) {
	pfs := make([]plannedFile, 0, len(resp.PlannedFiles))
	chs := make([]changeEntry, 0, len(resp.PlannedFiles))
	for _, pf := range resp.PlannedFiles {
		wirePF := plannedFile{Path: pf.Path, Action: pf.Action}
		count, hunks := computeChangeCountAndHunks(pf)
		if withHunks && len(hunks) > 0 {
			wirePF.Hunks = toCLIHunks(hunks)
		}
		pfs = append(pfs, wirePF)
		chs = append(chs, changeEntry{Path: pf.Path, Count: count})
	}
	return pfs, chs
}

// computeChangeCountAndHunks applies the T0-(g) `changes[].count`
// semantics AND returns the hunks (or nil for binary/no-change paths)
// so the caller can re-use them for the wire-Hunks field when --diff
// is set. The double-return keeps the diff invocation single per
// PlannedFile regardless of flag combination — the previous form
// computed hunks twice (once in computeChangeCount via CountLines as
// a wrong fallback, once again in mapResponseToWire via CountFromHunks)
// and the modify-no-diff path returned the whole-new-file line count,
// violating the T0-(g) contract.
//
// Action-rules:
//   - "create": count = CountLines(NewContent), hunks computed for
//     full-file insertion shape.
//   - "modify": count = sum(hunk.NewLines) over computed hunks.
//   - "delete": count = 0 (review #8: even for binary deletes, where
//     a naive CountBytesDiff would return len(OldContent)).
//   - binary content (non-delete): count = CountBytesDiff, hunks=nil
//     so wirePF.Hunks remains omitted (T0-(l) Spec-konformes Fallback).
func computeChangeCountAndHunks(pf driving.PlannedFile) (int, []driving.Hunk) {
	// Delete always returns 0 (T0-(g)) — short-circuit BEFORE the
	// binary-check so the CountBytesDiff trap doesn't fire for
	// binary deletes (review #8).
	if pf.Action == "delete" {
		return 0, nil
	}
	if diff.IsBinary(pf.OldContent, pf.NewContent) {
		return diff.CountBytesDiff(pf.OldContent, pf.NewContent), nil
	}
	hunks := diff.Compute(pf.OldContent, pf.NewContent)
	switch pf.Action {
	case "create":
		return diff.CountLines(pf.NewContent), hunks
	case "modify":
		// CountAdditions (not CountFromHunks): Spec §477 example
		// pins `count: 6` for the 6-line postgres-block append, NOT
		// 6+context. Slice T0-(g)'s original sum(hunk.NewLines) form
		// included context lines, drift against §477 (review-round-7
		// finding B).
		return diff.CountAdditions(hunks), hunks
	default:
		// Unknown action — keep parity with the create branch as the
		// safe fallback; the spec restricts action to {create, modify,
		// delete} (Spec §354) so this branch is unreachable today.
		return diff.CountLines(pf.NewContent), hunks
	}
}

// toCLIHunks copies driving.Hunk values into the CLI hunk wire-type.
// The field-level JSON tags are identical (T0-(l)); the copy is a
// schicht-separation guarantee, not a re-shape.
func toCLIHunks(src []driving.Hunk) []hunk {
	if len(src) == 0 {
		return nil
	}
	out := make([]hunk, len(src))
	for i, h := range src {
		out[i] = hunk{
			OldStart: h.OldStart,
			OldLines: h.OldLines,
			NewStart: h.NewStart,
			NewLines: h.NewLines,
			Content:  h.Content,
		}
	}
	return out
}

// mapErrorToDiagnostic maps an add-path error to a diagnosticItem
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
func mapErrorToDiagnostic(err error) diagnosticItem {
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
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-005", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// writeAddDiff renders the human-mode unified-diff string for each
// planned file (LH-FA-CLI-008). One file header per planned file
// followed by the hunks; binary files render only a header note.
// A blank line between file blocks keeps multi-file diffs visually
// separated; content-identical modifies render a "(no changes)"
// hint so the user does not interpret the empty body as a missed
// diff (review #15).
//
// Returns the first write error (broken pipe via `… | head -1`)
// instead of silently swallowing it (review #3) — the previous form
// dropped errors and the CLI exited 0 even after truncated output.
func writeAddDiff(out io.Writer, resp driving.AddServiceResponse) error {
	for i, pf := range resp.PlannedFiles {
		if i > 0 {
			if _, err := fmt.Fprintln(out); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(out, "--- %s (%s)\n", pf.Path, pf.Action); err != nil {
			return err
		}
		if diff.IsBinary(pf.OldContent, pf.NewContent) {
			if _, err := fmt.Fprintln(out, "(binary content — diff suppressed)"); err != nil {
				return err
			}
			continue
		}
		hunks := diff.Compute(pf.OldContent, pf.NewContent)
		if len(hunks) == 0 {
			if _, err := fmt.Fprintln(out, "(no changes)"); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprint(out, diff.Render(hunks)); err != nil {
			return err
		}
	}
	return nil
}

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
