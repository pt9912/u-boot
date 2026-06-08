# Slice V1: CLI-JSON-Envelope-Pattern-Konsolidierung add/init/generate

> **Status:** `open/` — **Trigger gefeuert (2026-06-08): Cluster
> `slice-v1-cli-json-dry-run` vollständig done (9/9 + T_close).**
> T0-Discovery gefahren (post-cluster Pre-Scan + Sub-Decisions
> geschärft); R-Runden + Lifecycle `open/`→`next/` **noch offen**.
> Konsolidierungs-Slice für den R15-Cross-Slice-1-Pattern-Drift aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> §Review-Round-15. Carveout-Plan-Anker ([[feedback_carveouts_need_plans]]);
> verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.
>
> **Standalone, NICHT in T_close absorbiert**: der Stub nannte als
> Option, die Konsolidierung in den Cluster-T_close zu bündeln
> (§Trigger Punkt 2). Das wurde **nicht** getan — der T_close war
> rein Mechanik-Abbau (Allowlist/Gate). Diese Konsolidierung ist
> damit ein eigenständiger post-cluster-Hygiene-Slice. **Verifiziert
> noch akkurat** (2026-06-08): `add.go:79` + `generate.go:78` rohes
> `cobra.ExactArgs(1)`; `init.go:138` `MaximumNArgs(1)`; add/init/
> generate nutzen **kein** `sanitizeBaseDir` (grep == 0).

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

**Subset already covered durch up-down T5** (Update 2026-06-07
nach up-down-Stub R3-MED-4 + R5-MED-3 Wortlaut-Präzisierung +
R6-MED-2 Inventur-Korrektur + up-down T8-Closure):
**HELPER-EXTRAKTION VOLLZOGEN** in up-down T5 — `cli/sanitize.go`
existiert seit Commit `a5aaf9c` als File im bestehenden
`package cli`. Sub-Decision 2 unten (File-vs-Sub-Package-Heim)
ist **OBSOLET** weil festgelegt.

up/down's **FS-Read-Wraps UND Compose-Runtime-Wraps** — insgesamt
**12 Stellen** (7 in `upservice.go:80,105,108,138,141,146,148` +
5 in `downservice.go:69,81,84,97,100`) — beide Klassen tragen
absolute Pfade via `req.BaseDir` oder `filepath.Join(...)`.
Function-Anker (robust gegen Zeilen-Drift): in `upservice`
betroffen `readUBootYAML` + `readComposeFile` + `Up`-Wrapper
selbst; in `downservice` analog `readUBootYAML` +
`readComposeFile` + `Down`-Wrapper. Alle sind durch den
up-down-Slice (Folge-Slice 6/9) selbst sanitized.
Pre-R5-Wortlaut "FS-Read-Wraps" war zu eng — Compose-Runtime-
Wraps (`ComposeUp on %q`, `ComposeDown on %q`) tragen
ebenfalls Pfad-Leak. Dieser Konsolidierungs-Slice braucht
beide Klassen NICHT mehr im Scope. Liste oben bleibt für
add/init/generate die Refactor-Ziele.

remove löste das in T7/T8 via `baseDirSanitizedError`-Wrapper
(`cli/remove.go:465-491`) + `replaceBareBaseDir`-Word-Boundary-
Helper (R15-LOW-1 robust gegen Substring-Kollisionen wie
`<baseDir>-cache/...`). **Helper-Heim festgelegt durch up-down
T5** (Update 2026-06-07 nach up-down-Stub R3-MED-3): up-down's
T5 extrahiert den Wrapper aus `cli/remove.go:465-540` in einen
neuen File **`cli/sanitize.go`** (im bestehenden `package cli`,
keine Sub-Package-Form). Wenn dieser Konsolidierungs-Slice
zieht, ist `cli/sanitize.go` schon das Heim — Sub-Decision 2
unten (Sub-Package-vs-File) ist obsolet, übernimm File-Heim.

## Heute-Stand-Pre-Scan (T0-Discovery 2026-06-08, post-cluster)

Code-Realität verifiziert (nach Cluster-Abschluss):

| Aspekt | add | init | generate | (Vorbild) remove | (Vorbild) config |
| --- | --- | --- | --- | --- | --- |
| `Args:` | `cobra.ExactArgs(1)` (`add.go:79`) | `cobra.MaximumNArgs(1)` (`init.go:138`) | `cobra.ExactArgs(1)` (`generate.go:78`) | `validateRemoveArgs(a)` | `configArgsValidator(a, sub, base)` |
| Envelope bei Args-Fehler | ❌ kein | ❌ (nur `>1` offen) | ❌ kein | ✅ | ✅ |
| `sanitizeBaseDir` an reportError | ❌ (grep 0) | ❌ | ❌ | ✅ `remove.go:297` | ✅ `config.go:264/300/355` |

**Zwei Vorbilder existieren jetzt** (beide nach dem Stub-Schreiben
gelandet) — und das **config-Muster ist bereits generisch**:
`configArgsValidator(a, subcommand, base)` (`config.go:224`) ruft
`base(cmd, args)` und bei Fehler im `--json`-Modus
`writeErrorEnvelopeSub(out, err, …, "config", subcommand,
mapConfigErrorToDiagnostic, nil)`. Der einzige config-Spezifikum-
Rest: hartkodiertes `"config"` + `mapConfigErrorToDiagnostic` + die
`subcommand == "set"`-Flag-Guard. **Generalisierung ist trivial** —
`command` + `mapErr` als Params, und die Flag-Reads unbedingt
versuchen (`cmd.Flags().GetBool` liefert `false` wenn das Flag nicht
registriert ist, also kein Guard nötig).

**Greedy-Sanitize ist das etablierte Muster** (NICHT selektiv, wie
der §Lösungs-Skizze-Punkt-2 oben vorschlug): config + remove wrappen
**jeden** UC-Error mit `sanitizeBaseDir(err, cwd)` vor `reportError`
(`config.go:264/300/355`, `remove.go:297`) — der Wrapper schreibt nur
um, wenn `baseDir` wirklich im String vorkommt, kein Performance-
Penalty. Stub-Empfehlung „selektiv" ist damit **überholt** → greedy
übernehmen (Defense-in-Depth, Symmetrie mit config/remove).

## Sub-Decisions (T0-geschärft — zum Review)

- **SD-A — Shared-Helper für Args-Validator.** Extrahiere
  `jsonArgsValidator(a *App, command, subcommand string, base
  cobra.PositionalArgs, mapErr func(error) diagnosticItem)
  cobra.PositionalArgs` (generalisierte `configArgsValidator`-Form,
  Flag-Reads unbedingt). **Reichweite (Sub-Decision zum Review):**
  (a) **Voll-Konsolidierung** — `configArgsValidator` +
  `validateRemoveArgs` refactoren auf den Shared-Helper (ihre
  Wrapper werden 1-Zeilen-Delegationen / entfallen), add/init/
  generate adoptieren ihn. Entfernt die Duplizierung vollständig
  (= Carveout-Geist), berührt aber done+getesteten config/remove-
  Code (Risiko durch starke Test-Suiten beherrscht). (b)
  **Minimal** — nur add/init/generate bekommen den Helper,
  config/remove bleiben wie sie sind (zwei Near-Duplikate bleiben).
  **Plan-Empfehlung: (a)** — die Konsolidierung ist genau der
  Carveout-Zweck; ein Near-Duplikat-Rest würde den Drift nur
  verschieben.
- **SD-B — Sanitize-Granularität → greedy** (überholt §Lösungs-
  Skizze-2): add/init/generate wrappen ihre UC-Errors an JEDER
  `reportError`-Site mit `sanitizeBaseDir(err, cwd)`, symmetrisch zu
  config/remove. `cli/sanitize.go` ist das Heim (up-down T5,
  `a5aaf9c`).
- **SD-C — init `MaximumNArgs(1)`-Sonderfall.** init nimmt 0-1
  Args (0 → Default-Projektname). Der `>1`-Pfad ist heute offen
  (kein Envelope). Lösung: `base = cobra.MaximumNArgs(1)` durch den
  Shared-Helper → der `>1`-Fehler trägt dann den Envelope. 0-Arg-
  Default-Logik bleibt in der RunE unberührt.
- **SD-D — Migrations-Reihenfolge → eine Tranche pro Schicht**
  (Helper → Adoptionen → Sanitize), nicht per-Command-Sub-Tranchen
  (R15-Empfehlung): der Helper-Extract berührt ohnehin alle Files;
  per-Command-Splitting würde Helper-Drift einführen.

## Aufhebungsbedingung

- `add`/`init`/`generate` tragen einen Args-Validator mit JSON-
  Envelope-Hook: `--json <cmd>` mit falscher Arg-Zahl
  (NoPositionalArg + TooManyArg) emittiert den Envelope auf stdout
  (§1841), und `--dry-run`/`--diff --json` wählt das Voll-Schema
  (§1842 / `LH-FA-CLI-007`/`-008`).
- Kein absoluter Filesystem-Pfad leakt mehr in `diagnostic.message`
  von add/init/generate (greedy `sanitizeBaseDir`).
- Bei SD-A (a): genau **ein** Shared-`jsonArgsValidator`; keine
  Near-Duplikate mehr (`configArgsValidator`/`validateRemoveArgs`
  delegieren oder sind weg).
- `make gates` grün (inkl. Pin-Tests pro Command für beide Pfade).
- Carveout-Eintrag (carveouts.md Z. 30) entfernt.

## Tranchen (vorläufig)

| Tranche | Inhalt | LOC |
| --- | --- | --- |
| T0 | Discovery (dieser Pre-Scan) + Sub-Decisions geschärft | — |
| T1 | Shared-Helper `jsonArgsValidator` in `cli/` + (SD-A (a)) `configArgsValidator`/`validateRemoveArgs` darauf refactoren; bestehende config/remove-Pins müssen grün bleiben | ~40 |
| T2 | add/init/generate adoptieren `jsonArgsValidator` (init via `MaximumNArgs(1)`-base, SD-C) + Pin-Tests pro Command (NoArg/TooMany × Minimal/Voll-Schema) | ~120 |
| T3 | greedy `sanitizeBaseDir` an allen add/init/generate-`reportError`-Sites + Path-Leak-Pins (mapper-message trägt keinen abs. Pfad) | ~60 |
| T4 | Closure: carveouts.md Z. 30 entfernen, ggf. CHANGELOG (Bug-Fix: §1841-Envelope-Symmetrie + Path-Leak-Defense), `done/`-Move + DoD | — |

LOC-Schätzung **~220** (mittlerer Slice). Read-bestätigt: kein neuer
Sentinel (Args-Fehler = `LH-FA-CLI-006`/Exit 2, schon vorhanden),
kein Port-Contract-Change.

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
2. **Helper-Heim für `baseDirSanitizedError`** — **festgelegt
   durch up-down T5 (2026-06-07, R3-MED-3 Festzurrung)**:
   Helper-Heim ist `cli/sanitize.go` (File im bestehenden
   `package cli`), extrahiert aus `cli/remove.go:465-540`.
   Sub-Decision (File-vs-Sub-Package) obsolet. Verbleibende
   Sub-Decision für DIESEN Slice: **Aufruf-Granularität** —
   greedy (wrap ALLE Errors) vs. selektiv (nur FS-Errors, weil
   Pre-Service-Sentinels nichts mit Pfad tunneln — overhead
   vs. defense-in-depth). Plan-Empfehlung: **selektiv**
   (Defense-in-Depth via Switch-Order, kein Performance-Penalty
   für CLI-Form-Validierungen die eh kein FS gesehen haben).
   Add/init/generate wrappen ihre UC-Errors analog `runRemove
   :299` mit `sanitizeBaseDir(err, cwd)`.
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
- **Confirmer-Swap-Pattern für `--purge`-Gate** (Update
  2026-06-07 nach up-down-Stub R4-MED-3): nach up-down-T2 ist
  `SilenceConfirmer`-Bool-Pattern in BEIDEN
  `RemoveServiceRequest` UND `DownRequest` etabliert
  (identische Form). `down --volumes` nutzt jetzt dasselbe
  Pattern wie `remove --purge`. Konsolidierungs-Wert für
  zukünftige Slices: das `SilenceConfirmer`-Bool plus Request-
  time Gate-Branch ist ein Pattern, kein remove-Spezifikum.
  Falls künftige Slices weitere Confirmer-Pfade einführen
  (z. B. `config set` destructive-Reset oder ein
  hypothetisches `prune`-Subcommand), erben sie die etablierte
  Bool-Form direkt. **Aktive Konsolidierungs-Pflicht** falls
  ein dritter Confirmer-Subcommand landet (Trigger-Schwelle):
  Helper-Heim für die `noopConfirmer{}`-Branch-Logic in einen
  geteilten `cli`-Sub-Helper (analog `mapComposeRuntimeSentinel`-
  Pattern für Mappers).
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
