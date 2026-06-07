package cli

import (
	"errors"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// mapComposeRuntimeSentinel klassifiziert den Compose-Daemon- und
// Runtime-Sentinel-Cluster für die JSON-Mapper-Tabelle von
// `u-boot up` und `u-boot down` (slice-v1-cli-json-dry-run-up-down
// T5 / T0-(e) R2-LOW-2). Beide Subcommands teilen sich diese drei
// Pfade (Daemon-unreachable / Compose-Runtime-Failure / Stabili-
// zations-Timeout) — Helper-Heim verhindert Duplikation in den
// per-Subcommand-Mappers.
//
// Switch-Order-Disziplin (T0-(e) R3-HIGH-1): `mapUpErrorToDiagnostic`
// und `mapDownErrorToDiagnostic` rufen diesen Helper NACH dem
// FS-Sentinel-Check und VOR den fachlichen Subcommand-spezifischen
// Sentinels — Multi-`%w`-Wraps mit FS+Docker fallen auf die FS-
// Klasse (LH-NFA-REL-003 / Exit 14), nicht auf Docker (Exit 11).
//
// Returnt `(code, true)` wenn einer der zwei Sentinels matched, sonst
// `("", false)` — Caller fällt auf den nächsten Switch-Case durch.
//
// LH-Code für beide Pfade ist `LH-NFA-REL-003` (T0-(f) Konsolidierung
// mit Doku-/Test-Pin-Pflicht für die `(code, exitCode)`-Tupel-
// Disambiguation in `cli-json-output.md` §6.7). Exit-Code wird vom
// separaten [ExitCode]-Helper aus dem Sentinel abgeleitet
// (driven.ErrDockerUnavailable → 11, driven.ErrComposeRuntime → 12) —
// dieser Helper liefert nur den LH-Code für `diagnosticItem.Code`.
func mapComposeRuntimeSentinel(err error) (code string, matched bool) {
	switch {
	case errors.Is(err, driven.ErrDockerUnavailable):
		return "LH-NFA-REL-003", true
	case errors.Is(err, driven.ErrComposeRuntime):
		return "LH-NFA-REL-003", true
	}
	return "", false
}
