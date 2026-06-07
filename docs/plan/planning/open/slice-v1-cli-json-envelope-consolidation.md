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
  Voll-Schema (Spec §1842-Verletzung, `LH-FA-CLI-007`-Vertrag;
  analog R13-HIGH-1 für remove).
- `u-boot --diff --json add` → analog, Minimal-Envelope statt
  Voll-Schema mit Hunks (`LH-FA-CLI-008` §451-489-Verletzung).
  Defense-Symmetrie zum `--dry-run`-Pfad: Validator MUSS auch
  `--diff` lesen und Voll-Schema-Wahl pinnen.

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

- **Cluster-Stand 8/9 erreicht** (config done): erst dann tragen
  ALLE modifying-Subcommands (add/init/generate/remove/up-down/
  logs/config) JSON-Envelopes, und die Helper-Extraktion lohnt
  sich gegen das vollständige Konsumenten-Set. **Cluster-
  Reihenfolge** (`slice-v1-cli-json-dry-run.md` §Per-Command-
  Folge-Slices): 6 up-down, 7 logs, 8 config, 9 template
  (template ist read-only Array-Output, kein Args-Validator-
  Pattern-Drift). Up-down (6/9) + logs (7/9) **erben** den
  Pattern-Drift, weil sie nach diesem Stub landen (Konsolidierung
  ist post-hoc, keine Pre-Pattern-Vorlage). Der Slice ist deshalb
  **retroaktive Konsolidierung** für alle dann existierenden
  Subcommands — nicht ein "von Anfang an erben"-Vorlauf.
- **Pre-Cluster-T_close-Hygiene-Pass**: alternativ bündelt
  Cluster-T_close (`slice-v1-cli-json-dry-run` Closure-Slice)
  die Konsolidierung mit der Allowlist-/PersistentPreRunE-
  Entfernung in einem Refactoring-Schritt. Dann entfällt dieser
  Stub und wird in Cluster-T_close absorbiert.
- **Real-World-Beschwerde** über fehlende Envelope-Symmetrie
  (z. B. CI-Konsument bricht auf `u-boot add --json` ohne arg
  ab weil kein JSON kommt) — feuert auch vor 8/9.
- **Security-Audit-Befund** zum Path-Leak in `diagnostic.message`
  (Info-Disclosure-CVE-Klasse) — feuert auch vor 8/9.

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

- **Pre-Pattern-Vorlage für Folge-Slices 6/7/8**: up-down (6/9),
  logs (7/9) und config (8/9) landen **vor** diesem
  Konsolidierungs-Slice (Cluster-Reihenfolge ist fix, kein
  Re-Sequencing). Sie kopieren das Pattern-Drift-Vorbild aus
  add/init/generate weiter — diese Slice ist deshalb
  retroaktive Helper-Extraktion über das vollständige
  Subcommand-Set, NICHT ein Pre-Pattern-Vorlauf den 6/7/8 von
  Anfang an erben.
- **`template` (9/9)**: read-only Array-Output ohne
  modifying-Args-Validator. Kein Pattern-Drift, deshalb nicht
  Teil dieses Slice.
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
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag bei
  `--dry-run --json` (Flag-Awareness-Pflicht: Voll-Schema bei
  `--dry-run`).
- `LH-FA-CLI-008` §451-489 — Voll-Schema-Vertrag bei
  `--diff --json` (Hunks-Pflicht plus Pre-Service-Validation-
  Symmetrie: NoPositionalArg- und TooManyArgs-Pfad mit
  `--diff` MUSS Voll-Schema-Envelope tragen, sonst Spec-
  Verletzung analog R13-HIGH-1 für `--dry-run`).
- `LH-FA-CLI-006` — usage-Error-Klasse für Form-Validierung
  (gemeinsamer Code-Pfad für NoPositionalArg + TooManyArgs).
