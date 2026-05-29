package cli

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

// renderUpStatus writes the LH-FA-UP-003 four-column status table
// to `out`. Columns: SERVICE / CONTAINER / PORT / HEALTH. Empty Port
// and Healthcheck render as "-" so columns align visually even for
// services that don't expose either.
//
// Service ordering: the [UpService] application service already
// sorts by Name (deterministic CLI output guaranteed); this
// function does not re-sort.
//
// Fire-and-forget contract: when `services` is empty (e.g.
// `up --timeout=0`) the function emits NO table — not even an
// empty header — so the LH-NFA-USE-004 golden-file pin can assert
// "no SERVICE header in stdout" for that mode.
func renderUpStatus(out io.Writer, services []domain.ServiceStatus) error {
	if len(services) == 0 {
		return nil
	}
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "SERVICE\tCONTAINER\tPORT\tHEALTH"); err != nil {
		return fmt.Errorf("write status header: %w", err)
	}
	for _, s := range services {
		port := s.Port
		if port == "" {
			port = "-"
		}
		health := s.Healthcheck
		if health == "" {
			health = "-"
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.Name, s.ContainerStatus.String(), port, health); err != nil {
			return fmt.Errorf("write status row %q: %w", s.Name, err)
		}
	}
	return w.Flush()
}

// renderUpDiagnostics writes the Diagnostic section beneath the
// status table. Both [domain.SeverityInfo] and [domain.SeverityWarn]
// entries are shown; [domain.SeverityOK] entries are filtered
// (M6 currently doesn't emit any). When `quiet` is true the whole
// section is suppressed — LH-FA-CLI-005 quiet semantics for up.
//
// Layout: one blank line before the section (visual separator from
// the status table), then one line per entry of the form
// `<severity>: <message>`, plus an indented `hint: <hint>` line
// when the diagnostic carries a hint.
func renderUpDiagnostics(out io.Writer, diagnostics []domain.Diagnostic, quiet bool) {
	if quiet {
		return
	}
	emitted := 0
	for _, d := range diagnostics {
		if d.Severity == domain.SeverityOK {
			continue
		}
		if emitted == 0 {
			fmt.Fprintln(out)
		}
		emitted++
		fmt.Fprintf(out, "%s: %s\n", d.Severity.String(), d.Message)
		if d.Hint != "" {
			fmt.Fprintf(out, "  hint: %s\n", d.Hint)
		}
	}
}

// renderDownSuccess writes the one-line down success message. When
// `quiet` is true nothing is written (LH-FA-CLI-005 quiet semantics
// for down — the asymmetric `--quiet` contract from the M6 slice).
// The message form mirrors the slice plan §T6:
//   - "environment stopped"                  (RemovedVolumes=false)
//   - "environment stopped, volumes removed" (RemovedVolumes=true)
func renderDownSuccess(out io.Writer, removedVolumes, quiet bool) {
	if quiet {
		return
	}
	if removedVolumes {
		fmt.Fprintln(out, "environment stopped, volumes removed")
		return
	}
	fmt.Fprintln(out, "environment stopped")
}
