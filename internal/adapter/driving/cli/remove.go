package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// removeFlags bundles the per-invocation flag state of
// `u-boot remove <service>`. Purge is the local destructive opt-in
// (LH-FA-ADD-007 §"Volumes nur auf explizite Anforderung"); Yes /
// NoInteractive / JSON read through from the root command's
// PersistentFlags so `u-boot --json remove postgres --purge` and
// `u-boot remove postgres --purge --json` behave identically.
//
// The slice-v1-cli-json-dry-run-remove slice (T5) extends this
// with --dry-run / --diff (LH-FA-CLI-007/008) and the read-through
// of the root --json flag, so [runRemove] sees the full preview-
// mode state in one struct.
type removeFlags struct {
	Purge         bool
	Yes           bool
	NoInteractive bool
	DryRun        bool
	Diff          bool
	JSON          bool
}

// removeEnvelopeData is the typed `data` carrier for the JSON
// envelope of `u-boot remove` (slice-v1-cli-json-dry-run-remove
// T0-(f)/(m)). Pointer-Wrapping pinnt Key-Presence-vs-Absence
// (Spec §1841): Success-Pfad setzt alle vier Felder; Error-Pfad
// trägt nur Service (Zero-Response für PriorState/State/
// VolumesPurged → die *-Felder bleiben nil und fallen via
// omitempty aus dem JSON). Pre-Service-Validation-Pfade
// (NoPositionalArg, ErrConflictingModeFlags) übergeben `data=nil`
// an reportError; kein Service-Kontext existiert dort.
//
// Pointer-Wahl pro Feld:
//   - Service: plain string — auf allen Pfaden gesetzt, sobald
//     [domain.NewServiceName] passiert ist; `""` ist nie ein
//     valider Wert.
//   - PriorState/State: *string — analog VolumesPurged für
//     Symmetrie, plus Defense gegen Default-`""`-Drift in
//     künftigen [domain.ServiceState]-Erweiterungen.
//   - VolumesPurged: *bool — `false` ist ein valider Success-Wert
//     (v0.3.0 deferred Volumes); plain `bool`+omitempty würde den
//     Error-Pfad-Zero und Success-Pfad-false identisch droppen
//     (Pattern analog [cliJSONEnvelope.DryRun]/`Diff` in
//     jsonenvelope.go:34-37).
type removeEnvelopeData struct {
	Service       string  `json:"service"`
	PriorState    *string `json:"priorState,omitempty"`
	State         *string `json:"state,omitempty"`
	VolumesPurged *bool   `json:"volumesPurged,omitempty"`
}

// newRemoveCommand builds the `u-boot remove <service>` Cobra
// subcommand (LH-FA-ADD-007).
//
// Symmetric to `add <service>`: positional service name validated
// via [domain.NewServiceName] before reaching the use case; catalog
// rejection and state-machine mismatches both surface as exit code
// 10 via [ExitCode]. The local `--purge` flag opts into volume
// removal (LH-FA-CLI-005A §254-style confirmation gate); the
// application service handles the gate consistency with `down
// --volumes`.
//
// Custom Args-Validator (slice-v1-cli-json-dry-run-remove T5,
// R11-HIGH-F1 + R12-HIGH-F1 + R12-MED-F2 mechanism): replaces
// the legacy `cobra.ExactArgs(1)` guard so that
// `u-boot remove --json` ohne positional arg den JSON-Envelope
// auf stdout emittiert (Spec §1841) BEVOR Cobra den
// [ErrServiceNameMissing]-Sentinel an Execute() propagiert.
// `cobra.ExactArgs(1)` würde sonst FRÜHER feuern und die
// envelope-emission überstimmen.
func newRemoveCommand(a *App) *cobra.Command {
	flags := &removeFlags{}

	cmd := &cobra.Command{
		Use:   "remove <service>",
		Short: "Remove a service add-on from the u-boot project",
		Long: `Remove an integrated service add-on from the current u-boot project
(LH-FA-ADD-007 — mirror of "u-boot add"). The positional <service>
argument names the catalogue entry; today MVP supports only
postgres. Keycloak / OpenTelemetry will follow in their own V1
slices.

The command runs in an initialized project (u-boot.yaml present per
LH-FA-ADD-001) and is idempotent: removing an already-disabled
service is a no-op with a clear message. The state-machine handles
the inverse of add — strip the service.<name> managed block from
compose.yaml and .env.example, then set services.<name>.enabled to
false in u-boot.yaml.

--purge is the explicit destructive opt-in for volume removal
(LH-FA-ADD-007). The same LH-FA-CLI-005A §254 confirmation gate as
"u-boot down --volumes" applies: --no-interactive without --yes
exits 10 (ErrConfirmationRequired); interactive mode prompts with a
safe default-No. In v0.3.0 the gate's "approved" outcome does NOT
auto-remove volumes yet — the CLI surfaces the deferred work so the
user can clean up manually with "docker volume rm <name>".

Flag combinations (LH-FA-CLI-007/008):
  --dry-run            preview without writing files
  --diff               show unified diff of planned changes
  --dry-run --diff     unified diff preview, no write
  --json               JSON output; pairs with --dry-run / --diff
                       for the LH-FA-CLI-007 §326 voll-schema

Exit codes (LH-FA-CLI-006):
  0   success (state transition OR idempotent no-op)
  2   CLI / flag errors (unknown subcommand, missing positional,
      conflicting mode flags)
  10  fachlich: ErrServiceUnsupported (not in catalogue),
      ErrServiceUnregistered (never added), ErrServiceInconsistent
      (orphan block or missing entry), ErrProjectNotInitialized,
      ErrConfirmationRequired (purge gate refused)
  14  filesystem error (permission, disk full, race)

Examples:
  u-boot remove postgres                    # state-transitioning remove
  u-boot remove postgres                    # idempotent re-run: no-op
  u-boot remove postgres --purge --yes      # opt into volume cleanup
  u-boot remove postgres --dry-run --json   # preview as JSON
  u-boot remove redis                       # exit 10 — not in catalogue`,
		Args: validateRemoveArgs(a),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Yes = a.yes
			flags.NoInteractive = a.noInteractive
			flags.JSON = a.json
			return runRemove(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr(), args, *flags, a.removeServiceUseCase, a.getwd)
		},
	}

	cmd.Flags().BoolVar(&flags.Purge, "purge", false,
		"also request volume removal for the service (LH-FA-ADD-007). Destructive: triggers the LH-FA-CLI-005A §254 confirmation gate (refuses in --no-interactive without --yes). v0.3.0 does NOT auto-remove volumes after approval — the summary points at `docker volume rm` for manual cleanup.")
	cmd.Flags().BoolVar(&flags.DryRun, "dry-run", false,
		"preview the planned changes without writing files (LH-FA-CLI-007)")
	cmd.Flags().BoolVar(&flags.Diff, "diff", false,
		"render a unified diff of the planned changes (LH-FA-CLI-008)")
	return cmd
}

// validateRemoveArgs is the `u-boot remove <service>` positional-args
// validator. Originally a bespoke closure (slice-v1-cli-json-dry-run-
// remove T5/T7: R11/R12/R13 — JSON-envelope-emission BEFORE the Cobra
// return so `--json remove` without an arg emits the envelope per
// Spec §1841/§1842, plus too-many-args symmetry and --dry-run/--diff
// flag-awareness for the Voll-Schema).
//
// At slice-v1-cli-json-envelope-consolidation T1 (SD-A (a)) it became
// a delegation to the shared [jsonArgsValidator]: previewFlags=true
// (remove carries --dry-run/--diff → Voll-Schema on arg errors), and
// the base closure [removeArgsBase] preserves the
// [ErrServiceNameMissing] missing-arg sentinel (Schutzplanke 2)
// instead of cobra's bare "accepts 1 arg(s)". `data` stays nil — no
// service context exists before the positional is parsed.
func validateRemoveArgs(a *App) cobra.PositionalArgs {
	return jsonArgsValidator(a, "remove", "", removeArgsBase, mapRemoveErrorToDiagnostic, true)
}

// removeArgsBase is the per-command base validator: len(args)==0 →
// ErrServiceNameMissing (custom sentinel, Exit 2 via isUsageError),
// len(args)==1 → ok, len(args)>1 → cobra's TooMany error.
func removeArgsBase(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return ErrServiceNameMissing
	}
	return cobra.ExactArgs(1)(cmd, args)
}

// runRemove is split from the Cobra closure for direct unit-testing
// (no Cobra construction needed). Mirrors [runAdd]'s shape; the
// mutual-exclusion check on --yes / --no-interactive lives here for
// the same reason (CLI-level usage error, not a use case concern).
//
// Five output paths (slice-v1-cli-json-dry-run-remove T5):
//
//   - Human (no --json), no preview-flag → [printRemoveSummary]
//     plus the today's `--purge`-deferred-volumes WARNING on errOut.
//   - Human (no --json) with --diff → [writeDiff] for the unified
//     diff string on stdout, then [printRemoveSummary]; --purge
//     WARNING-prosa stays on errOut (T0-(c) R2-LOW-F6 stderr-
//     separation: the diff body never gets prose-polluted).
//   - --json without --dry-run/--diff → [newDataEnvelope] with
//     `command="remove"`, `data` carrying the success quartet
//     `{service, priorState, state, volumesPurged}` (or just
//     `service` on the error path — Zero-Response analog
//     T0-(f)). WARN-Diagnostics from `resp.Warnings` get mapped
//     via [mapWarningsToDiagnostics].
//   - --dry-run --json (with or without --diff) → [newFullEnvelope]
//     with plannedFiles[] from the recorder, dryRun=true, no
//     FS-write, plus the same data carrier + WARN-Diagnostics.
//   - --diff --json without --dry-run → [newFullEnvelope]
//     preview-and-apply, plannedFiles[] + hunks, dryRun=false,
//     diff=true.
//
// Error-Pfad (T0-(j) R4-MED-F3 Variante A): Error-Diagnostic
// dominiert, WARN unterdrückt. reportError reicht `resp.Warnings`
// NICHT durch — der `mapErr`-Pfad in [writeErrorEnvelopeSub] emittiert
// ausschließlich den Error-Diagnostic (`level: "error"`); ein
// WARN würde sich auf nicht-vorhandene `data.volumesPurged`-Daten
// beziehen, weil die Zero-Response-Klausel `priorState`/`state`/
// `volumesPurged` aus dem `data`-Block zieht. T6-Pin:
// `TestRemove_PurgeYesJSON_MidWriteFailure_ErrorOnly`.
//
// `--purge --diff` ohne `--json` (T0-(h) PreviewAndApply): die
// `printRemoveSummary`-WARNING bleibt auf errOut, NICHT im Diff-
// Body (R2-LOW-F6 stderr-Trennung).
func runRemove(
	ctx context.Context,
	out, errOut io.Writer,
	args []string,
	flags removeFlags,
	uc driving.RemoveServiceUseCase,
	getwd func() (string, error),
) error {
	mapErr := mapRemoveErrorToDiagnostic

	if flags.Yes && flags.NoInteractive {
		return reportError(out, ErrConflictingModeFlags, nil, flags.DryRun, flags.Diff, flags.JSON, "remove", mapErr, nil)
	}

	svcName, err := domain.NewServiceName(args[0])
	if err != nil {
		return reportError(out, err, nil, flags.DryRun, flags.Diff, flags.JSON, "remove", mapErr, nil)
	}

	data := removeEnvelopeData{Service: svcName.String()}

	cwd, err := getwd()
	if err != nil {
		return reportError(out, fmt.Errorf("determine working directory: %w", err), nil, flags.DryRun, flags.Diff, flags.JSON, "remove", mapErr, data)
	}

	mode := previewModeFromFlags(flags.DryRun, flags.Diff)
	resp, removeErr := uc.Remove(ctx, driving.RemoveServiceRequest{
		BaseDir:          cwd,
		ServiceName:      svcName,
		Purge:            flags.Purge,
		Yes:              flags.Yes,
		NoInteractive:    flags.NoInteractive,
		PreviewMode:      mode,
		SilenceConfirmer: flags.JSON,
	})

	if removeErr != nil {
		// T7 R14-MED-1 path-leak fix: FS-Wraps in der Use-Case
		// (`fmt.Errorf("remove write %s: %w: %w", absPath, …)`) tragen
		// absolute Pfade die `err.Error()` 1:1 in `diagnostic.message`
		// laufen lässt. Pre-sanitize gegen `cwd` damit der User-facing
		// Output project-relative Pfade zeigt (analog
		// `mapCaptureToPlannedFiles` für `plannedFiles[].path`).
		// errors.Is-Matching bleibt intakt via Unwrap-Chain — der
		// Sanitizer wraps, replaced nicht.
		return reportError(out, sanitizeBaseDir(removeErr, cwd), resp.PlannedFiles, flags.DryRun, flags.Diff, flags.JSON, "remove", mapErr, data)
	}

	if flags.JSON {
		return writeRemoveJSON(out, resp, flags.DryRun, flags.Diff, svcName)
	}

	if flags.Diff {
		if err := writeDiff(out, resp.PlannedFiles); err != nil {
			return err
		}
	}
	printRemoveSummary(out, errOut, resp, flags.Purge, mode)
	return nil
}

// writeRemoveJSON renders the success-path JSON envelope. Two
// shapes (T0-(f)/(m) Festzurrungen):
//
//   - dryRun=false && diff=false → [newDataEnvelope] (Minimal-
//     kontrakt plus `data: {service, priorState, state,
//     volumesPurged}`).
//   - dryRun=true || diff=true   → [newFullEnvelope] mit
//     plannedFiles/changes plus dem gleichen `data`-Träger;
//     optional hunks im --diff-Pfad.
//
// `resp.Warnings` werden via [mapWarningsToDiagnostics] in
// `diagnostics[]` mit `level: "warn"` gemapped — Status-Kopplung
// (Spec §447) macht `status: "warn"` und exit-code bleibt 0
// (Warn-only verschiebt den Exit-Code nicht). T0-(j) R1-MED-5-Pin
// `TestRemove_PurgeYesJSON_WarnOnly`.
func writeRemoveJSON(out io.Writer, resp driving.RemoveServiceResponse, dryRun, diffFlag bool, svcName domain.ServiceName) error {
	priorState := resp.PriorState.String()
	state := resp.State.String()
	volumesPurged := resp.VolumesPurged
	data := removeEnvelopeData{
		Service:       svcName.String(),
		PriorState:    &priorState,
		State:         &state,
		VolumesPurged: &volumesPurged,
	}
	warnDiags := mapWarningsToDiagnostics(resp.Warnings)
	if !dryRun && !diffFlag {
		env := newDataEnvelope("remove", "", data, warnDiags, 0)
		return writeEnvelope(out, env)
	}
	pfs, chs := mapPlannedFilesToWire(resp.PlannedFiles, diffFlag)
	env := newFullEnvelope("remove", "", dryRun, diffFlag, pfs, chs, data, warnDiags, 0)
	return writeEnvelope(out, env)
}

// mapWarningsToDiagnostics maps a `[]driving.WarningEntry` (Use-
// Case-Source-of-Truth) to a `[]diagnosticItem` (CLI-Wire-Form)
// for the JSON envelope's `diagnostics[]` array (slice-v1-cli-
// json-dry-run-remove T0-(g) + T5). Triviales Field-Mapping —
// Code/Level/Message 1:1, `Subject` landet auf `File` weil das
// die heute existierende Diagnostic-Field-Form für freien
// `subject`-Kontext ist (Cluster-Vorlauf-Annahme aus T2 R12-
// LOW-F4; ein dedizierter `subject`-Field-Add wäre Folge-Slice-
// Arbeit, nicht remove-Scope).
//
// Returns nil for empty/nil input — newDataEnvelope/
// newFullEnvelope normalisieren das auf `diagnostics: []` im
// JSON. Source-of-Truth ist der Use-Case: nur Catalog-Lookup für
// `volumeOptional` kennt die Service-Volume-Semantik (T0-(g)
// R3-F1).
func mapWarningsToDiagnostics(ws []driving.WarningEntry) []diagnosticItem {
	if len(ws) == 0 {
		return nil
	}
	out := make([]diagnosticItem, 0, len(ws))
	for _, w := range ws {
		d := diagnosticItem{
			Level:   w.Level,
			Code:    w.Code,
			Message: w.Message,
		}
		if w.Subject != "" {
			d.File = w.Subject
		}
		out = append(out, d)
	}
	return out
}

// mapRemoveErrorToDiagnostic maps a remove-path error to a
// [diagnosticItem] with the spec-konforme LH-Kennung per T0-(e)
// Switch-Order-Pflicht.
//
// Switch-Order verbindlich (slice-Plan T0-(e), R3-MED-F3-Fix):
// Infrastruktur-Sentinels first (FS, ConfirmerUnavailable) →
// Confirmation-Gate → fachliche Service-Sentinels →
// CLI-Form-Sentinels → Default. Infrastruktur-First schützt vor
// Multi-`%w`-Wraps (Go 1.20+) die einen fachlich-Sentinel
// versehentlich tunneln — z. B. ein synthetisch konstruierter
// `fmt.Errorf("%w: %w", ErrConfirmerUnavailable,
// ErrConfirmationRequired)` MUSS `LH-FA-CLI-005A` / Exit 10
// (I/O-Klasse) erzeugen, NICHT `LH-FA-INIT-005`. Defense-only-
// Pin (T0-(e) R4-LOW-F5): heute existiert KEIN Code-Pfad der
// beide Sentinels gemeinsam chained — T6 verifiziert die
// Mapper-Robustheit gegen einen synthetischen Multi-Wrap, nicht
// ein reales Failure-Szenario.
//
// `LH-FA-ADD-007` Multi-Use (T0-(g) R5-F2): derselbe Code wird
// für `ErrServiceUnregistered` (ERROR-Pfad, hier) UND für die
// `--purge && !VolumesPurged`-WARN-Diagnostic (WARN-Pfad, via
// [mapWarningsToDiagnostics] in `writeRemoveJSON`) genutzt.
// Disambiguation läuft über `(code, level)`-Tupel — Konsumenten
// dürfen NICHT nur auf `code` filtern.
//
// `ErrConfirmationRequired` → `LH-FA-INIT-005` (geteilt mit
// init/down) — der Spec-Anker §195 lebt im INIT-005-Block.
// `ErrConflictingModeFlags` → `LH-FA-CLI-005A` (Mutex-Verträge
// für `--yes`/`--no-interactive` §235).
// `ErrServiceNameMissing` → `LH-FA-CLI-006` (Form-Validierung
// vom CLI-Adapter emittiert, Exit 2 via [isUsageError]).
//
// Default fallback maps to LH-FA-CLI-006 / Exit 1 via
// [ExitCode] (NICHT automatisch Exit 2 — [isUsageError] matched
// nur die CLI-Sentinels, nicht den generischen Default-Pfad).
func mapRemoveErrorToDiagnostic(err error) diagnosticItem {
	switch {
	case errors.Is(err, driving.ErrRemoveFileSystem):
		return diagnosticItem{Level: "error", Code: "LH-NFA-REL-003", Message: err.Error()}
	case errors.Is(err, driving.ErrConfirmerUnavailable):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-005A", Message: err.Error()}
	case errors.Is(err, driving.ErrConfirmationRequired):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-005", Message: err.Error()}
	case errors.Is(err, driving.ErrServiceUnsupported):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-002", Message: err.Error()}
	case errors.Is(err, driving.ErrServiceUnregistered):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-007", Message: err.Error()}
	case errors.Is(err, driving.ErrServiceInconsistent):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-005", Message: err.Error()}
	case errors.Is(err, driving.ErrProjectNotInitialized):
		return diagnosticItem{Level: "error", Code: "LH-FA-ADD-001", Message: err.Error()}
	case errors.Is(err, domain.ErrInvalidServiceName):
		return diagnosticItem{Level: "error", Code: "LH-FA-INIT-006", Message: err.Error()}
	case errors.Is(err, ErrConflictingModeFlags):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-005A", Message: err.Error()}
	case errors.Is(err, ErrServiceNameMissing):
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	default:
		return diagnosticItem{Level: "error", Code: "LH-FA-CLI-006", Message: err.Error()}
	}
}

// Sanitizer-Helper (baseDirSanitizedError, replaceBareBaseDir,
// isPathComponentByte, sanitizeBaseDir) leben nach slice-v1-cli-
// json-dry-run-up-down T5 in cli/sanitize.go — package-intern
// wiederverwendet von up/down + remove. Call-Site `runRemove:299`
// nutzt sie unverändert.

// printRemoveSummary writes a short, deterministic summary of the
// remove outcome. Two shapes for the stdout summary:
//
//   - Idempotent no-op (Changed=nil):
//     "Service <name> is already disabled; no changes."
//   - State transition (Changed!=nil):
//     "Removed service <name>." + list of changed paths.
//
// When `--purge` was requested AND the response shows VolumesPurged=
// false (always true in v0.3.0 — actual volume removal is deferred)
// AND we are NOT in dry-run preview, a WARNING block is appended to
// errOut (review-followup F4: stderr keeps stdout clean for future
// --json consumers). The wording (review-followup F5) is explicit
// about what was NOT done so users don't trust the prior NOTE-as-
// aside framing.
//
// Dry-Run-Suppression (slice-v1-cli-json-dry-run-remove T7 R14-
// HIGH-2 Fix): in `PreviewDryRun` läuft die Use-Case ohne Gate-Aufruf
// und ohne tatsächliche Mutation (Plan T0-(h)(a) Gate-Skip). Ohne
// Suppression würde die WARNING-Prosa suggerieren "deferred work" —
// aber im Dry-Run wurde gar nichts versucht. Die Prosa wäre damit
// semantisch falsch ("würde-deferred" wäre korrekt, "ist-deferred"
// ist es nicht). Fix: WARN-Block überspringen wenn
// `previewMode == PreviewDryRun`. PreviewAndApply behält die WARN,
// weil der Gate dann läuft und die Mutation wirklich stattfindet —
// hier ist `VolumesPurged: false` der bewusst-deferred Status.
//
// JSON-mode bypasses this helper entirely (runRemove returns
// directly from `writeRemoveJSON`); the WARNING migrates into the
// envelope's `diagnostics[]` via the application service's
// `resp.Warnings` + [mapWarningsToDiagnostics] (T0-(g) WARN-
// Migration). Human-mode keeps the legacy stderr prose intact —
// existing tooling that grep'd `WARNING:` on stderr keeps working.
func printRemoveSummary(out, errOut io.Writer, resp driving.RemoveServiceResponse, purge bool, previewMode driving.PreviewMode) {
	name := resp.ServiceName.String()

	if len(resp.Changed) == 0 {
		fmt.Fprintf(out, "Service %q is already disabled; no changes.\n", name)
	} else {
		fmt.Fprintf(out, "Removed service %q.\n\nChanged:\n", name)
		for _, p := range resp.Changed {
			fmt.Fprintln(out, "  - "+p)
		}
	}

	if purge && !resp.VolumesPurged && previewMode != driving.PreviewDryRun {
		fmt.Fprintf(errOut,
			"\nWARNING: --purge was requested but volume removal is NOT yet automated in v0.3.0.\n"+
				"         The %s service's named volumes are still on disk and untouched.\n"+
				"         Remove them manually after confirming the data is no longer needed:\n"+
				"           docker volume ls --filter label=com.docker.compose.project=<your-project>\n"+
				"           docker volume rm <name>\n",
			name)
	}
}
