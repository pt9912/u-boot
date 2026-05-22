// Package progress is the driven-adapter implementations of
// [driven.ProgressPort]. Today only a text formatter for stdout
// exists; a JSON variant for `u-boot init --json` (LH-NFA-USE-004,
// V1) plugs in here later.
//
// Layer rule: adapters may import the domain and their driven-port
// interface, plus external libraries; they may not import application
// or other adapter packages (LH-FA-ARCH-003, depguard-enforced).
package progress

import (
	"fmt"
	"io"

	"github.com/pt9912/u-boot/internal/hexagon/port/driven"
)

// TextWriter implements [driven.ProgressPort] by formatting events
// as human-friendly lines on an io.Writer (typically os.Stdout in
// production, a *bytes.Buffer in tests). Format follows the M3-T4b
// convention so existing user-doc and screenshots stay valid.
type TextWriter struct {
	out io.Writer
}

// Static check: TextWriter satisfies the ProgressPort interface.
// Lives in text.go (not in a `_test.go` file) so a mismatch breaks
// the package build, not only the test build.
var _ driven.ProgressPort = (*TextWriter)(nil)

// NewText returns a TextWriter that writes events to out. Passing
// nil is a programmer error — wire io.Discard explicitly when no
// output is wanted.
func NewText(out io.Writer) *TextWriter {
	return &TextWriter{out: out}
}

// AffectedFiles renders the LH-FA-INIT-005 §609 summary as a
// header line plus one line per row. The "(with backup)" marker
// follows the canonical em-dash + parenthetical format.
func (t *TextWriter) AffectedFiles(baseDir string, rows []driven.AffectedFile) {
	fmt.Fprintf(t.out, "Affected files in %s:\n", baseDir)
	for _, r := range rows {
		marker := ""
		if r.Backup {
			marker = " (with backup)"
		}
		fmt.Fprintf(t.out, "  - %s — %s%s\n", r.Path, actionLabel(r.Action), marker)
	}
}

// actionLabel maps the driven-port action enum to a human label.
// Centralised here (not as a Stringer on the enum) so presentation
// stays out of the port and a future JSON adapter can use its own
// vocabulary (e.g. snake_case keys).
func actionLabel(a driven.AffectedAction) string {
	switch a {
	case driven.AffectedReplaceBlock:
		return "replace managed block"
	case driven.AffectedOverwriteFull:
		return "full overwrite"
	}
	return fmt.Sprintf("action(%d)", int(a))
}
