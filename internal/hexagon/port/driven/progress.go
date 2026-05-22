package driven

// ProgressPort is the application's side-channel for reporting
// affected-paths information to the user during a re-init
// (LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 — "vor dem Schreiben
// muss eine Zusammenfassung der betroffenen Pfade ausgegeben
// werden"). The application emits structured events through this
// port; the adapter (text on stdout today, JSON for `--json` later)
// is responsible for presentation.
//
// Introduced in the M3-T4c-review to close the hexagonal-layer-leak
// finding #8: previously the application wrote pre-formatted text
// to an io.Writer, baking presentation into the use-case package.
//
// All methods are best-effort — implementations must not return
// errors. A failure to emit progress is never reason to abort the
// underlying use case; the worst observable consequence is a user
// missing a one-line heads-up.
type ProgressPort interface {
	// AffectedFiles is called once per Init invocation, before any
	// write side-effect. baseDir is the absolute project root the
	// rows are relative to; rows are sorted by Path so the
	// adapter's output is deterministic. An empty rows slice is
	// never passed (application skips the call instead) so an
	// adapter does not need to special-case "no affected files".
	AffectedFiles(baseDir string, rows []AffectedFile)
}

// AffectedFile is one entry in a [ProgressPort.AffectedFiles]
// report. Fields are intentionally narrow: Path + Action + Backup
// are enough for the LH-FA-INIT-005 §609 summary; presentation
// (labels, indentation, glyphs) belongs in the adapter.
type AffectedFile struct {
	// Path is the file path relative to baseDir (e.g. "compose.yaml").
	Path string
	// Action is the kind of change the use case will perform.
	Action AffectedAction
	// Backup is true when the use case will copy the file to
	// `<path>.bak[.N]` before mutating it (LH-FA-INIT-005 §605/§617).
	Backup bool
}

// AffectedAction enumerates the kinds of mutations a re-init can
// apply to an existing file. The application emits the enum value;
// adapters choose how to render it.
type AffectedAction int

const (
	// AffectedReplaceBlock means only the file's
	// `U-BOOT MANAGED BLOCK: init` region is rewritten — content
	// outside the markers survives unchanged (LH-FA-INIT-005 §613).
	AffectedReplaceBlock AffectedAction = iota
	// AffectedOverwriteFull means the entire file is replaced —
	// LH-FA-INIT-005 §619 requires --backup for this path.
	AffectedOverwriteFull
)
