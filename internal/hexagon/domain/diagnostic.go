package domain

import "sort"

// Severity classifies the gravity of a [Diagnostic] per
// LH-FA-DIAG-003. The three levels are deliberately small — the
// spec leaves no room for project-specific intermediate severities.
//
// Ordering: SeverityOK < SeverityInfo < SeverityWarn < SeverityError
// so that [DiagnosticReport.MaxSeverity] reduces over a slice with a
// straightforward max operation. SeverityInfo sits between OK and
// Warn — see its constant doc for the rationale.
type Severity int

const (
	// SeverityOK means the diagnostic passed. Reports with only OK
	// diagnostics map to exit code 0.
	SeverityOK Severity = iota
	// SeverityInfo means the diagnostic carries a non-judgmental hint
	// — typically a mode acknowledgement (e.g. M6 `u-boot up
	// --timeout=0` emits an `up.fire-and-forget` info entry to
	// signal that status polling was skipped). Info is strictly
	// between OK and Warn: it must not push the report into a
	// non-zero exit code, but it is also not a clean OK that the
	// CLI's `--quiet` filter should drop alongside OK entries.
	// Added in M6-T1; existing M4 doctor checks never emit Info.
	SeverityInfo
	// SeverityWarn means the check found a concern that does not
	// block proceeding. Without --strict, reports with only OK+Warn
	// still map to exit code 0 (LH-FA-DIAG-003).
	SeverityWarn
	// SeverityError means the check failed in a way that blocks the
	// project from working. Any single Error makes the report
	// fail-grade and forces a non-zero exit code.
	SeverityError
)

// String returns the canonical lowercase label used in CLI text
// output ("ok", "warn", "error"). Adapters that need other casings
// or glyphs are expected to format from the enum value, not from
// this string.
func (s Severity) String() string {
	switch s {
	case SeverityOK:
		return "ok"
	case SeverityInfo:
		return "info"
	case SeverityWarn:
		return "warn"
	case SeverityError:
		return "error"
	default:
		return "unknown"
	}
}

// Diagnostic is one entry in a [DiagnosticReport]. ID is the stable
// machine-readable check identifier (e.g. `"docker.version"`); it
// stays consistent across runs so log scrapers and CI dashboards can
// pin alerts. Message is the human-readable explanation; Hint is the
// optional LH-FA-DIAG-004 repair suggestion.
type Diagnostic struct {
	// ID identifies the check that emitted this diagnostic (e.g.
	// `"docker.available"`, `"yaml.valid"`). Stable across runs and
	// versions; once shipped it lives in CI configurations and
	// must not change without a deprecation cycle.
	ID string

	// Severity is the LH-FA-DIAG-003 classification.
	Severity Severity

	// Message is the human-readable summary of what the check
	// observed. Sentence-case, ends with a period.
	Message string

	// Hint is the LH-FA-DIAG-004 repair suggestion, empty when none
	// is available. Sentence-case, ends with a period. Mostly set
	// for Warn / Error; OK diagnostics rarely carry a hint.
	Hint string
}

// DiagnosticReport is the aggregate of all [Diagnostic] entries a
// single `u-boot doctor` run produced. The slice carries the
// deterministic emit-order of the service; the CLI adapter may
// re-sort for presentation (typically by severity then ID).
type DiagnosticReport struct {
	// Items holds every diagnostic the service emitted, in the
	// order it ran the checks. An empty slice means no checks ran;
	// "all OK" is represented by a non-empty slice of only
	// SeverityOK entries.
	Items []Diagnostic
}

// MaxSeverity returns the highest [Severity] across all Items. An
// empty report returns SeverityOK (the natural identity for max);
// the CLI then maps that to exit code 0.
func (r DiagnosticReport) MaxSeverity() Severity {
	highest := SeverityOK
	for _, item := range r.Items {
		if item.Severity > highest {
			highest = item.Severity
		}
	}
	return highest
}

// HasErrors reports whether any diagnostic in the report has
// SeverityError. Convenience for the CLI's exit-code dispatch
// (LH-FA-DIAG-003: any error → exit ≠ 0).
func (r DiagnosticReport) HasErrors() bool {
	return r.MaxSeverity() == SeverityError
}

// HasWarnings reports whether the highest severity is exactly
// SeverityWarn. False when the report only has OK entries or when
// any Error is present (Error is strictly worse than Warn). Used
// by the CLI's `--strict`-mode exit dispatch (LH-FA-DIAG-003: with
// --strict any Warn also forces exit ≠ 0).
func (r DiagnosticReport) HasWarnings() bool {
	return r.MaxSeverity() == SeverityWarn
}

// SortedByIssuesFirst returns the items reordered so that Errors
// come first, then Warnings, then OK; within a severity bucket
// the original ID-alphabetical order is preserved. Helpful for
// CLI rendering where the operator's eye should land on real
// problems first. Does not mutate the receiver.
func (r DiagnosticReport) SortedByIssuesFirst() []Diagnostic {
	out := append([]Diagnostic(nil), r.Items...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return out[i].Severity > out[j].Severity
		}
		return out[i].ID < out[j].ID
	})
	return out
}
