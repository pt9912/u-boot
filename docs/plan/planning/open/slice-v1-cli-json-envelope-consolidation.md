# Slice V1: CLI-JSON-Envelope-Pattern-Konsolidierung add/init/generate

> **Status:** `open/`, on hold pending trigger. Konsolidierungs-Slice
> für den R15-Cross-Slice-1-Pattern-Drift aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> §Review-Round-15. Carveout-Plan-Anker ([[feedback_carveouts_need_plans]]);
> verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

Die remove-R13/R14-Pre-T8-Review-Runden haben zwei Pattern-Defekte
aufgedeckt, die strukturell **auch** in add/init/generate vorhanden
sind. Pre-T7-Code-State im remove-Slice + R15-Cross-Slice-Audit:

### Defekt 1 — Custom-`Args`-Validator fehlt in add/generate

`cli/add.go:79` und `cli/generate.go:78` nutzen rohes
`cobra.ExactArgs(1)` ohne JSON-Envelope-Hook. Konsequenz:

- `u-boot --json add` ohne positional arg → Cobra emittiert
  `"accepts 1 arg(s), received 0"` auf stderr, **kein** JSON-
  Envelope auf stdout (Spec §1841-Verletzung).
- `u-boot --json add postgres extra` → analog, kein Envelope
  (Spec §1841-Verletzung).
- `u-boot --dry-run --json add` → Minimal-Envelope statt
  Voll-Schema (Spec §1842-Verletzung; analog R13-HIGH-1 für
  remove).

Generate hat dasselbe Pattern an `cli/generate.go:78`. Init
nutzt `cobra.MaximumNArgs(1)` mit interner Logik die zumindest
für `len(args)==0` einen Default einsetzt — anders gelagert,
aber `len(args)>1` ist offen.

remove löste das in T5/T7 via `validateRemoveArgs(a *App)`:
Custom-`cobra.PositionalArgs`-Closure mit `*App`-Capture, die
für `len(args)==0` UND `len(args)>1` das stdout-Envelope **vor**
dem Cobra-Return emittiert. Plus Flag-Awareness via
`cmd.Flags().GetBool("dry-run"/"diff")` für Voll-Schema-Wahl.

### Defekt 2 — `baseDirSanitizedError`-Wrapper fehlt in add/init/generate

`mapAddErrorToDiagnostic` / `mapInitErrorToDiagnostic` /
`mapGenerateErrorToDiagnostic` reichen `err.Error()` 1:1 an
`diagnostic.message`. Use-Case-Wraps der Form

```go
fmt.Errorf("...write %s: %w: %w", absPath, ErrXxxFileSystem, raw)
```

tunneln den absoluten Filesystem-Pfad in den User-facing Output.
Im JSON-Mode ist der Pfad maschinen-lesbar abgreifbar →
Info-Leak des User-Filesystem-Layouts.

Wrap-Site-Inventar (Stichproben aus R15-Audit):

- `application/addservice_execute.go:672, 689`
- `application/initproject.go:925, 967, 1015, 1117, 1143`
- `application/generate.go:210` (+ weitere im
  devcontainer-Phase-2-Block)

remove löste das in T7/T8 via `baseDirSanitizedError`-Wrapper
(`cli/remove.go:465-491`) + `replaceBareBaseDir`-Word-Boundary-
Helper (R15-LOW-1 robust gegen Substring-Kollisionen wie
`<baseDir>-cache/...`).

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Cluster-Stand 7/9 oder 8/9 erreicht**: die Konsolidierung
  lohnt sich erst wenn alle modifying-Subcommands JSON-Envelopes
  tragen — sonst verliert die Helper-Extraktion an Konsumenten.
  Up-down (6/9) + Config (8/9) sind die nächsten Folge-Slices,
  die das Pattern erben — wenn die ohne Konsolidierung landen,
  vergrößert sich der Refactoring-Druck.
- **Real-World-Beschwerde** über fehlende Envelope-Symmetrie
  (z. B. CI-Konsument bricht auf `u-boot add --json` ohne arg
  ab weil kein JSON kommt).
- **Security-Audit-Befund** zum Path-Leak in `diagnostic.message`
  (Info-Disclosure-CVE-Klasse).

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen, vor der eigentlichen Implementierung
zu klären:

1. **Helper-Heim für `validateArgs`-Pattern**: extrahiere
   `validateRemoveArgs` auf eine generische Form
   `validateArgsForJSON(a *App, command string, expected int)`-
   Closure-Factory. Lebt in `cli/` oder neuem `cli/cobraargs/`-
   Sub-Package. Add/init/generate rüsten ihre `Args:`-Felder
   um. Sub-Decision: per-Command-Custom-Mapper (analog
   `mapRemoveErrorToDiagnostic`) als Parameter mitgeben, ODER
   die generische Form ohne Mapper (`LH-FA-CLI-006` ist der
   einzige Code in beiden Pfaden, Pre-Service-Validation-
   Sentinel-Klasse).
2. **Helper-Heim für `baseDirSanitizedError`**: extrahiere den
   Wrapper + `replaceBareBaseDir` + `isPathComponentByte` aus
   `cli/remove.go` in eine eigene `cli/sanitize/` Sub-Package
   oder als `cli/`-Helper. Add/init/generate wrappen ihre
   UC-Errors analog `runRemove:299` mit `sanitizeBaseDir(err,
   cwd)`. Sub-Decision: greedy (wrap ALLE Errors) vs. selektiv
   (nur FS-Errors, weil Pre-Service-Sentinels nichts mit Pfad
   tunneln — overhead vs. defense-in-depth).
3. **Migrations-Reihenfolge**: alle drei (add/init/generate) in
   einer PR, ODER per-Command-Sub-Tranchen mit jeweils eigener
   Pin-Test-Sequenz. R15-Cross-Slice empfiehlt eine PR weil das
   Helper-Extract-Refactoring ohnehin die drei Files berührt
   und einzelne Sub-Tranchen Helper-Drift einführen würden.

## Out of Scope

- **`up`/`down`/`logs`/`config`/`template`**: diese landen erst
  in den Folge-Slices 6-9. Sie erben das konsolidierte Pattern
  von Anfang an (Plan-Vertrag: Folge-Slices 6-9 referenzieren
  diesen Stub als Pattern-Vorbild, falls die Konsolidierung vor
  ihnen landet).
- **Confirmer-Swap-Pattern für `--purge`-Gate**: das ist
  remove-spezifisch (Confirmer ist heute nur im remove-Pfad
  aktiv). Falls künftige Slices Confirmer brauchen (down --volumes
  hat einen, ist aber kein modifying-CLI-Subcommand mit
  `--purge`-Symmetrie), eigener Konsolidierungs-Schritt.
- **Pattern-Inventur für `LH-FA-ADD-007` Multi-Use**: derselbe
  LH-Code für ERROR + WARN ist heute nur in remove relevant.
  Falls ein zukünftiger Slice analogen Multi-Use einführt,
  wandert der Disambiguations-Vertrag in einen separaten
  Helper.

## Spec-Bezug

- `LH-NFA-USE-004` §1841 — Minimalkontrakt-Envelope-Vertrag
  (Symmetrie-Pflicht für alle JSON-Pfade).
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag (Flag-
  Awareness-Pflicht: Voll-Schema bei `--dry-run` ODER `--diff`).
- `LH-FA-CLI-006` — usage-Error-Klasse für Form-Validierung.
