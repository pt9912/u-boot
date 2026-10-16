package cli

import "github.com/pt9912/u-boot/internal/hexagon/port/driving"

// visibleWarnings returns the warnings the human renderer writes.
// Without `quiet` that is every warning; under `--quiet` only the
// entries flagged [driving.WarningEntry.SecurityRelaxation]
// (LH-FA-CLI-005 exception 1: the disclosure of a sandbox security
// relaxation per LH-FA-DEV-006/-007/-008) or
// [driving.WarningEntry.DegradationFallback] (exception 2: the
// disclosure of a degradation fallback per LH-FA-DEV-007); every
// other notice is dropped. JSON mode does not filter: `--quiet
// --json` keeps every entry in `diagnostics[]`.
func visibleWarnings(ws []driving.WarningEntry, quiet bool) []driving.WarningEntry {
	if !quiet {
		return ws
	}
	var out []driving.WarningEntry
	for _, w := range ws {
		if w.SecurityRelaxation || w.DegradationFallback {
			out = append(out, w)
		}
	}
	return out
}
