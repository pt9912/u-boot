package cli

// cliJSONEnvelope ist der Wire-Type für `u-boot --json`-Ausgaben
// (slice-v1-cli-json-dry-run-doctor T2). Das Lastenheft trennt
// zwei Kontraktstufen — Minimalkontrakt (LH-NFA-USE-004 §1841) und
// Voll-Schema (LH-FA-CLI-007 §326). Beide werden in **diesem** Typ
// gerendert, der Konstruktor pinnt die Stufe:
//
//   - newMinimalEnvelope: read-only-Aufrufe (doctor, logs, up, …)
//     setzen DryRun/Diff/PlannedFiles/Changes nicht; sie fallen
//     per `omitempty` aus dem JSON.
//   - newFullEnvelope: --dry-run/--diff-Aufrufe (add, init, …)
//     setzen alle vier Voll-Felder explizit.
//
// IMPORTANT: DryRun und Diff sind *bool — NOT plain bool. Spec
// §326 verlangt dryRun/diff im modifying-Pfad auch wenn der Wert
// false ist. Plain `bool` + omitempty würde false aus dem JSON
// werfen und das Spec-Required-Set verletzen. Siehe
// docs/plan/planning/in-progress/slice-v1-cli-json-dry-run-doctor.md
// §T0-(d) (Review-Finding M1).
type cliJSONEnvelope struct {
	Status     string `json:"status"`
	Command    string `json:"command"`
	Subcommand string `json:"subcommand,omitempty"`
	// DryRun/Diff: *bool (see IMPORTANT note above).
	DryRun *bool `json:"dryRun,omitempty"`
	Diff   *bool `json:"diff,omitempty"`
	// PlannedFiles/Changes: *[]T pointer wrapper, same reasoning as
	// DryRun/Diff. omitempty on a plain []T drops len==0 slices,
	// which would hide the spec-required `"plannedFiles":[]`/
	// `"changes":[]` form in the full-envelope path. The pointer
	// is nil for the minimal path (field omitted) and points to a
	// (possibly empty) slice for the full path (field always present).
	PlannedFiles *[]plannedFile   `json:"plannedFiles,omitempty"`
	Changes      *[]changeEntry   `json:"changes,omitempty"`
	Diagnostics  []diagnosticItem `json:"diagnostics"`
	ExitCode     int              `json:"exitCode"`
	// Data: optionales Free-Form-Feld (Spec §1839 erlaubt zusätzliche
	// Felder im Minimal-Mode). slice-v1-cli-json-dry-run-generate T0-(p)
	// hat das ursprünglich für Template-Slice 9/9 reservierte Feld
	// vorgezogen, weil generate als erster Multi-Artefakt-Subcommand
	// ein `data.artifact: "<…>"` braucht (T0-(m)) plus `data.action:
	// "<…>"` als Discriminator zwischen UpdatedBlock und
	// RepairedManual (T0-(f)). Template-Slice 9/9 erbt das Feld nur
	// noch (`data: []templateJSON`). Disziplin „Konstruktor pinnt
	// die Stufe": `newDataEnvelope` ist der einzige Konstruktor, der
	// Data setzt; `omitempty` hält die Minimal/Voll-Envelopes
	// data-frei. Marshal-Pin-Test in jsonenvelope_test.go.
	Data any `json:"data,omitempty"`
}

// diagnosticItem ist ein Eintrag im `diagnostics[]`-Array. Spec
// §1834 lässt `level` nur `warn` oder `error` zu; SeverityOK- und
// SeverityInfo-Items werden beim Mapping aus
// domain.DiagnosticReport übersprungen, nicht als level: "ok"/"info"
// serialisiert.
type diagnosticItem struct {
	Level   string `json:"level"`
	Code    string `json:"code"`
	Message string `json:"message"`
	File    string `json:"file,omitempty"`
}

// plannedFile ist ein Eintrag im Voll-Schema-`plannedFiles[]`. Spec
// §354 erlaubt `action` nur create/modify/delete. Hunks sind optional
// (LH-FA-CLI-008 §477-482 macht sie Pflicht nur im `--diff --json`-
// Pfad; ohne `--diff` bleibt das Feld via omitempty weg —
// slice-v1-cli-json-dry-run-add T0-(l)).
type plannedFile struct {
	Path   string `json:"path"`
	Action string `json:"action"`
	Hunks  []hunk `json:"hunks,omitempty"`
}

// hunk ist die Wire-Form eines Diff-Hunks in
// `plannedFiles[].hunks` (LH-FA-CLI-008 §477-482; slice-v1-cli-json-
// dry-run-add T0-(l)). Coordinates sind 1-basiert wenn die jeweilige
// *Lines-Zahl > 0; pure-add/-delete-Hunks setzen die Off-Seite auf 0
// per Unified-Diff-Konvention. Field-Tag-Drift (`offset` statt
// `oldStart`) wird vom jsontestutil.checkHunks-Helper gepinnt.
type hunk struct {
	OldStart int    `json:"oldStart"`
	OldLines int    `json:"oldLines"`
	NewStart int    `json:"newStart"`
	NewLines int    `json:"newLines"`
	Content  string `json:"content"`
}

// changeEntry ist ein Eintrag im Voll-Schema-`changes[]`. Spec §368
// fordert count ≥ 0.
type changeEntry struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// newMinimalEnvelope baut einen Minimalkontrakt-Envelope (LH-NFA-
// USE-004 §1841). DryRun/Diff/PlannedFiles/Changes/Data bleiben nil
// und fallen per omitempty aus dem JSON. status wird aus diags
// abgeleitet (Spec §447 / §1837): error → "error"; warn ohne
// error → "warn"; sonst "ok".
//
// Für Minimal+Data-Pfade (generate, template list) gibt es den
// dedizierten [newDataEnvelope]-Konstruktor — "Konstruktor pinnt
// die Stufe", siehe Doc oben am Typ.
func newMinimalEnvelope(command, subcommand string, diags []diagnosticItem, exitCode int) cliJSONEnvelope {
	if diags == nil {
		diags = []diagnosticItem{}
	}
	return cliJSONEnvelope{
		Status:      statusFromDiagnostics(diags),
		Command:     command,
		Subcommand:  subcommand,
		Diagnostics: diags,
		ExitCode:    exitCode,
	}
}

// newFullEnvelope baut einen Voll-Schema-Envelope (LH-FA-CLI-007
// §326). Alle vier Voll-Felder werden explizit gesetzt; bei
// `dryRun=false`/`diff=false` erscheint im JSON entsprechend
// `"dryRun":false`/`"diff":false` (Spec-Required-Set).
//
// data trägt subcommand-spezifische Free-Form-Inhalte (slice-v1-
// cli-json-dry-run-generate T0-(p)/(q)); init/add reichen `nil`
// durch, generate reicht `{"artifact":"<…>","action":"<…>"}` durch.
// `data == nil` lässt das Feld via omitempty aus dem JSON fallen.
func newFullEnvelope(
	command, subcommand string,
	dryRun, diff bool,
	planned []plannedFile,
	changes []changeEntry,
	data any,
	diags []diagnosticItem,
	exitCode int,
) cliJSONEnvelope {
	if diags == nil {
		diags = []diagnosticItem{}
	}
	if planned == nil {
		planned = []plannedFile{}
	}
	if changes == nil {
		changes = []changeEntry{}
	}
	return cliJSONEnvelope{
		Status:       statusFromDiagnostics(diags),
		Command:      command,
		Subcommand:   subcommand,
		DryRun:       &dryRun,
		Diff:         &diff,
		PlannedFiles: &planned,
		Changes:      &changes,
		Diagnostics:  diags,
		ExitCode:     exitCode,
		Data:         data,
	}
}

// newDataEnvelope baut einen Minimalkontrakt-Envelope mit einem
// zusätzlichen `data`-Träger für subcommand-spezifische Free-Form-
// Inhalte (slice-v1-cli-json-dry-run-generate T0-(p)). DryRun/Diff/
// PlannedFiles/Changes bleiben nil und fallen per omitempty aus dem
// JSON — `data` ist orthogonal zum Voll-Schema und ergänzt den
// Minimalkontrakt um ein Konsumenten-lesbares Datenfeld.
//
// Aufrufer:
//   - generate (Folge-Slice 4/9): `command="generate"`,
//     `subcommand=""`, `data={"artifact":"<…>","action":"<…>"}`.
//   - template list (Folge-Slice 9/9): `command="template"`,
//     `subcommand="list"`, `data=[]templateJSON`.
//
// `data == nil` lässt das Feld via omitempty aus dem JSON fallen,
// sodass dieser Konstruktor auch als trivialer Wrapper-Pfad für
// Fälle ohne data-Träger nutzbar bleibt (init/add im Error-Pfad
// ohne Artefakt-Kontext reichen `nil` durch).
func newDataEnvelope(
	command, subcommand string,
	data any,
	diags []diagnosticItem,
	exitCode int,
) cliJSONEnvelope {
	if diags == nil {
		diags = []diagnosticItem{}
	}
	return cliJSONEnvelope{
		Status:      statusFromDiagnostics(diags),
		Command:     command,
		Subcommand:  subcommand,
		Diagnostics: diags,
		ExitCode:    exitCode,
		Data:        data,
	}
}

// statusFromDiagnostics implementiert Spec §447 / §1837 — die
// `status`-Kopplung an das höchste `level` in `diagnostics`.
// SeverityOK/SeverityInfo erscheinen nie als diagnosticItem
// (Filter beim Mapping aus domain.DiagnosticReport).
func statusFromDiagnostics(diags []diagnosticItem) string {
	hasWarn := false
	for _, d := range diags {
		switch d.Level {
		case "error":
			return "error"
		case "warn":
			hasWarn = true
		}
	}
	if hasWarn {
		return "warn"
	}
	return "ok"
}
