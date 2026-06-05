package cli

import (
	"fmt"
	"io"

	"github.com/pt9912/u-boot/internal/adapter/driving/cli/diff"
	"github.com/pt9912/u-boot/internal/hexagon/port/driving"
)

// reportError is the single error-emission gate for every modifying-
// subcommand RunE (slice-v1-cli-json-dry-run-init T0-(e) / T5-a).
//
// Decomposed-Slices-Signatur (statt service-spezifischer Response-
// Pointer): nimmt `[]driving.PlannedFile` direkt, damit init/remove/
// generate/config-set ohne ihre eigenen Response-Types das Helper
// teilen können. `command` und `mapErr` werden vom Caller-RunE
// gesetzt (`mapErr := mapAddErrorToDiagnostic` für add,
// `mapErr := mapInitErrorToDiagnostic` für init — Symmetrie-
// Pattern aus T0-(e) mapErr-Source-Pflicht).
//
// Verhalten:
//   - Human-Mode (`jsonFlag=false`): returnt den raw error — Cobra
//     propagiert ihn, main.go rendert die stderr-Meldung und
//     cli.ExitCode bestimmt den Shell-Exit.
//   - JSON-Mode (`jsonFlag=true`): schreibt den Envelope auf stdout
//     UND returnt den original err so cli.ExitCode den richtigen
//     Exit-Code bestimmen kann (add review #2 — ohne Propagation
//     wäre die Shell-Exit 0 trotz envelope-claimed 14).
func reportError(
	out io.Writer,
	err error,
	planned []driving.PlannedFile,
	dryRun, diffFlag, jsonFlag bool,
	command string,
	mapErr func(error) diagnosticItem,
) error {
	if !jsonFlag {
		return err
	}
	if envErr := writeErrorEnvelope(out, err, planned, dryRun, diffFlag, command, mapErr); envErr != nil {
		return envErr
	}
	return err
}

// writeErrorEnvelope renders the JSON envelope on the error path.
//
// Voll-Schema-Switch (add review #4): voll-schema applies whenever
// the recorder captured anything (`len(planned) > 0`) OR the user
// explicitly asked for it via `--dry-run`/`--diff`. Without a
// recorder capture AND without a preview flag the envelope shape is
// the minimal contract (Spec §1841).
//
// dryRun/diffFlag werden VOM USER-FLAG-STATE durchgereicht, NICHT
// hardgecodet (add review #4 — frühere Form hatte `false, true`
// fix, was --dry-run --diff --json mit Mid-Failure auf falsche
// dryRun/diff-Werte mappte).
//
// File-Annotation: `diag.File = lastPlannedPath(planned)` für
// Mid-Write-Failure-Diagnostics (add review-round-6 #lastPlannedPath
// erblich).
func writeErrorEnvelope(
	out io.Writer,
	addErr error,
	planned []driving.PlannedFile,
	dryRun, diffFlag bool,
	command string,
	mapErr func(error) diagnosticItem,
) error {
	diag := mapErr(addErr)
	exitCode := ExitCode(addErr)
	if path := lastPlannedPath(planned); path != "" {
		diag.File = path
	}
	wantsFullSchema := len(planned) > 0 || dryRun || diffFlag
	if !wantsFullSchema {
		env := newMinimalEnvelope(command, "", []diagnosticItem{diag}, exitCode)
		return writeEnvelope(out, env)
	}
	pfs, chs := mapPlannedFilesToWire(planned, diffFlag)
	env := newFullEnvelope(command, "", dryRun, diffFlag, pfs, chs, []diagnosticItem{diag}, exitCode)
	return writeEnvelope(out, env)
}

// lastPlannedPath returns the path of the last PlannedFile in the
// list — convenient for Mid-Write-Failure-Diagnostics where the
// recorder's tail entry is the failing path. Returns "" for an
// empty list.
func lastPlannedPath(planned []driving.PlannedFile) string {
	if len(planned) == 0 {
		return ""
	}
	return planned[len(planned)-1].Path
}

// writeDiff renders the human-mode unified-diff string for each
// planned file (LH-FA-CLI-008). One file header per planned file
// followed by the hunks; binary files render only a header note.
// A blank line between file blocks keeps multi-file diffs visually
// separated; content-identical modifies render a "(no changes)"
// hint so the user does not interpret the empty body as a missed
// diff (add review #15).
//
// Returns the first write error (broken pipe via `… | head -1`)
// instead of silently swallowing it (add review #3) — the previous
// form dropped errors and the CLI exited 0 even after truncated
// output.
//
// Command-agnostisch (T0-(k) Option a): identischer
// `--- <path> (<action>)`-Header für alle Subcommands; per-command
// Header-Overrides sind Out-of-Scope V1.
func writeDiff(out io.Writer, planned []driving.PlannedFile) error {
	for i, pf := range planned {
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
