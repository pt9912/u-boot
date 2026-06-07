# Slice V1: `remove --json` / `--dry-run` / `--diff` — Add-Inverse mit Purge-Gate

> **Status (2026-06-07):** T2 bis T8 ✅, Slice **done**. R15-Bestätigungsrunde 2026-06-07 lieferte **HIGH=0, MED=0, LOW=2 + 1 Cross-Slice**: vier T7-Fixes adversarial bestätigt (Validator-Flag-Timing, Sanitizer-Unwrap-Chain, WARN-Konsistenz, Pattern-Verifikation), zwei neue R15-LOW-Findings in T8 mitgefixt (R15-LOW-1 Sanitizer-Substring-Robustheit, R15-LOW-2 TooMany-DryRun-Coverage-Pin), eine R15-Cross-Slice-1-Pattern-Drift-Erkenntnis (add/init/generate haben analoge Defekte) in neuen open/-Stub [`slice-v1-cli-json-envelope-consolidation`](../open/slice-v1-cli-json-envelope-consolidation.md) ausgelagert. T8-Closure-Commit siehe §Tranche-Status. Lifecycle-Übergänge: `open/` nach R1-R11, `next/` nach R12, `in-progress/` ab T2, `done/` nach T8; 56 Plan-Findings (R1-R12) + 11 Pre-T8-Findings (R13-R14) + 3 R15-Findings = 70 gesamt. Fünfter Folge-Slice (5/9) des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 5/9). Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
> 1:1 für die symmetrische Service-Operation und aus
> [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
> 1:1 für die `PreviewMode`-Carrier-Form,
> `RecordingFileSystem`-driven-Adapter, Pure-Go LCS-Diff-Renderer,
> `previewModeFromFlags`-Mapping und Error-Emission-Helper.
>
> Remove-spezifisch sind der **Confirmation-Gate für `--purge`**
> (`LH-FA-CLI-005A` §254 — mediated by `Yes` / `NoInteractive`),
> die **stderr-WARNING für `--purge` ohne Volume-Removal**
> (review-followup F4 — stderr-Sauberkeit für `--json`-Konsumenten
> ist genau das Problem, das diese Slice lösen muss), das
> **`VolumesPurged`-Status-Feld** im Envelope und die **Idempotenz-
> Semantik** (already-disabled → `Changed=nil`-NoOp).

## Auslöser

Cluster-Slice §T0-Outcomes (a)+(b)+(e) machen jeden modifying-
Subcommand für `--json`/`--dry-run`/`--diff` verbindlich
(`LH-NFA-USE-004` §1813, `LH-FA-CLI-007` §326, `LH-FA-CLI-008`
§451-489). `u-boot remove` ist nach `doctor`/`add`/`init`/`generate`
der nächste modifying-Subcommand und die **inverse Operation zu
`add`** — strip managed blocks aus `compose.yaml` + `.env.example`,
flip `services.<name>.enabled` auf `false` in `u-boot.yaml`,
optional Volume-Purge via `--purge`-Gate.

Spec-Bezug (geerbt von add/init):

- `LH-FA-CLI-007` (Dry-Run, Voll-Schema §326)
- `LH-FA-CLI-008` (Diff, §451-489)
- `LH-NFA-USE-004` (Minimalkontrakt §1841)

Remove-spezifische Spec-Anker:

- `LH-FA-ADD-007` (`u-boot remove <service>`-Surface inkl. `--purge`-
  Opt-in)
- `LH-FA-CLI-005A` §254 (Confirmation-Gate für destruktive
  Operationen — geteilt mit `down --volumes`)
- `LH-NFA-REL-003` (FS-Failure-Klasse, geerbt für Mid-Write
  analog init)

Heute-Stand-Pre-Scan
(`internal/hexagon/application/removeservice.go`, ~370 LOC;
`internal/adapter/driving/cli/remove.go`, ~170 LOC):

| Phase | Methode | Pfade | Code-Anker |
| --- | --- | --- | --- |
| Plan-Phase (per-File-Read + managedblock-Detection) | `Exists` (Z. 270, 356), `ReadFile` (Z. 277, 319) | `compose.yaml`, `.env.example`, `u-boot.yaml` | Plan-Phase, kein Write |
| Execute-Phase Write | `WriteFile` (Z. 234) | `compose.yaml` (managed-block strip), `.env.example` (managed-block strip), `u-boot.yaml` (enabled=false) | direkt; bei `Changed!=nil` |
| Execute-Phase Delete | `RemoveAll` (Z. 240 via plannedRemoveFile.removeAction) | optional extraFiles | direkt; nur bei extraFiles-Catalog-Entry |
| Volume-Purge | KEIN FS — heute deferred | `docker volume rm <name>` als WARNING im stderr-Block | nicht implementiert (v0.3.0) |

Damit nutzt remove **2-3 von 8** Recorder-Mutations-Methoden direkt
(`WriteFile` immer; `RemoveAll` für extraFiles wenn der Catalog-
Entry sie definiert). `WriteFileExclusive`, `Mkdir`, `MkdirAll`,
`Rename`, `Copy`, `CopyExclusive` werden NICHT gerufen — Recorder
deckt sie als Drift-Schutz trotzdem ab.

Use-Case-Deps: `driven.FileSystem` (Read + Write + RemoveAll),
`driven.YAMLCodec`, `driven.Confirmer` (für `--purge`-Gate),
`driven.Logger`. **KEIN** `GitClient`, **KEIN** `Progress`-Port.
**Confirmer-Port-Kollision mit `--json`-Mode** (R1-HIGH-2-
Klarstellung): der Confirmer-Prompt würde stdin/stdout
polluten. Init's T0-(o) ProgressPort-Silencing ist NICHT
direkt geerbt — init swappt nur `s.progress`, nicht
`s.confirmer`. **Confirmer-Swap ist ein neues Pattern, das
remove etabliert** (T0-(j) Sub-Decision), nicht ein geerbtes.
Pattern-Erbe-Disziplin T0-(a) Spalte führt das entsprechend.

## Aufhebungsbedingung

Acht Flag-Kombinationen für `u-boot remove <service>` liefern
spec-konforme Outputs (geerbt von add/init/generate):

```bash
u-boot remove postgres                          # human, schreibt
u-boot remove postgres --dry-run                # human Vorschau, kein Write
u-boot remove postgres --diff                   # human Unified-Diff + Write
u-boot remove postgres --dry-run --diff         # human Unified-Diff, kein Write
u-boot remove postgres --json                   # Minimal+Data-Envelope
u-boot remove postgres --dry-run --json         # Voll-Schema, kein Write
u-boot remove postgres --diff --json            # Voll-Schema mit Hunks
u-boot remove postgres --dry-run --diff --json  # Voll-Schema, Hunks, kein Write
```

Plus die `--purge`-Variante muss in JEDER Flag-Kombination ein
spec-konformes Verhalten zeigen.

`make gates` grün (lint + test + coverage-gate ≥ 90 % +
docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ Drei JSON-Pfade analog init/generate (`runRemove` ruft generische
  Helper mit `command="remove"` + `mapErr=mapRemoveErrorToDiagnostic`).
- ✅ **Envelope-Shape**: `command="remove"`, kein `subcommand`-Feld
  (Service ist Positional-Arg). `data.service` trägt den
  Service-Namen; `data.volumesPurged` trägt den
  Volume-Purge-Status (T0-(f) Sub-Decision).
- ✅ **Idempotenz-NoOp-Pin**: Single-Call und Repeat-Call gegen
  bereits disabled Service liefern `plannedFiles: []` UND
  `changes: []`, `data.priorState: "deactivated"`, `data.state:
  "deactivated"`, `status: ok`, Exit 0. **NUR `PriorState=Deactivated`
  qualifiziert als NoOp** (R2-MED-F2-Fix). `EnabledUnset` und
  `InconsistentBlock` sind state-transitioning (`Changed!=nil`,
  `plannedFiles!=[]`) — separater Pin nötig.
- ✅ **EnabledUnset-Normalisierungs-Pin** (R2-MED-F2-Fix):
  Service mit fehlendem `enabled`-Key in `services.<name>`
  liefert `plannedFiles: [u-boot.yaml, compose.yaml, .env.example]`,
  `data.priorState: "enabled-unset"`, `data.state: "deactivated"`,
  `changes[]` non-leer, `status: ok`, Exit 0. T6-Fixture muss
  beide Varianten (mit/ohne `enabled: false`) explizit setzen
  um Test-Drift zu vermeiden.
- ✅ **Mid-Write-Failure-Pin**: analog init/generate — `plannedFiles[]`
  enthält Captures bis zur Failure-Stelle (Recorder zeichnet vor
  Delegieren auf, `recordingfs.go:139`), `diagnostics[].file` =
  `lastPlannedPath`.
- ✅ **`--purge`-Gate-Verhalten im JSON-Mode**: Confirmer-Prompt-
  Silencing analog init's ProgressPort (T0-(o) Sub-Decision):
  `req.SilenceConfirmer = flags.JSON` oder per
  request-time-swap. Bei `--purge --no-interactive --json` OHNE
  `--yes` → `ErrConfirmationRequired` Envelope mit
  `LH-FA-INIT-005` / Exit 10.
- ✅ **stderr-WARNING-Migration** (T0-(g) festgezurrt): die
  heutige `printRemoveSummary`-WARNING auf errOut bei
  `--purge`-Status FALSE muss in den Envelope wandern als
  `diagnostics[]`-Eintrag. **AK-Pin explizit** (R2-LOW-F7-Fix):
  `diagnostics[0].code == "LH-FA-ADD-007"` AND
  `diagnostics[0].level == "warn"` AND `data.volumesPurged ==
  false`. Im JSON-Mode darf stderr nicht durch die
  WARNING-Prosa polluten.
- ✅ **Cobra-Args-Missing-Pfad-Pin** (R11-HIGH-F1-Fix für
  Spec-LH-NFA-USE-004-§1841-Konformität): `u-boot remove --json`
  ohne positional arg läuft via `cobra.ExactArgs(1)`-Guard
  (`cli/remove.go:77`) **vor** RunE — der JSON-Helper aus T5
  wird NIE aufgerufen, Konsument bekommt KEINEN Envelope und
  Exit 2. Symmetrie-Bruch zu `remove "bad name" --json` (voller
  Envelope mit `LH-FA-INIT-006`/Exit 10). **T5-Pflicht
  festgezurrt** (R12-MED-F2-Mechanismus): **Custom-`Args`-
  Validator** als Closure die `*App` per Konstruktor-Closure-
  Capture einfängt (analog `newRemoveCommand`-Form
  `cli/remove.go:36-37` mit `a *App` im Outer-Scope); bei
  `len(args)==0` returnt der Validator
  `cli.ErrServiceNameMissing` (Layer-Heim CLI, R12-HIGH-F1).
  **Pflicht-Begleit-Edit**: `Args: cobra.ExactArgs(1)` Z. 77
  durch `Args: validateRemoveArgs(a)` ersetzen — `cobra.
  ExactArgs(1)` würde sonst FRÜHER feuern und die Custom-Form
  überstimmen. **PreRunE-Alternative verworfen**: Layer-Mismatch
  (PreRunE feuert nach Args-Default, müsste `len(cmd.Args)`
  re-checken — redundant) plus `cobra.ExactArgs(1)` müsste
  trotzdem auf `cobra.ArbitraryArgs` umgestellt werden. Custom-
  `Args`-Closure ist die schlankere Form. `RunE`-Pfad:
  `reportError` mit dem Sentinel (`code: "LH-FA-CLI-006"` /
  `exitCode: 2`); `data` ist `nil` weil kein Service-Kontext
  vorhanden. T6-Pin:
  `TestRemove_NoPositionalArg_JSON_EmitsCLI006Envelope` mit
  empty `args[]` + `--json` → voller Envelope mit `command:
  "remove"`, `data: nil` (kein Service-Kontext), `code:
  "LH-FA-CLI-006"`, exit 2. Pattern-Erbe-Vorlauf für künftige
  Folge-Slices (up/down haben `cobra.ExactArgs(1)` für
  `<service>`-Subform — denselben Args-Guard-Pfad).
- ✅ **`ErrServiceUnregistered` ERROR-Pfad-Pin** (R6-MED-F3-Fix,
  Symmetrie zum WARN-Pfad-Pin oben): `LH-FA-ADD-007` wird auch
  als Error-Code für `ErrServiceUnregistered` genutzt (R5-F2
  Multi-Use-Klarstellung). AK-Pin: `diagnostics[0].code ==
  "LH-FA-ADD-007"` AND `diagnostics[0].level == "error"` AND
  `status == "error"` AND `exitCode == 10`. Konsumenten
  disambiguieren WARN-Pfad und ERROR-Pfad ausschließlich über
  `(code, level)`-Tupel, nicht über Code allein. T6-Pin
  `TestRemove_ServiceUnregisteredJSON_ErrorLevelCodePin`.
- ✅ **Volumes-Purge-Status im Envelope**: `data.volumesPurged: false`
  in v0.3.0 (deferred), mit Hint-Diagnostic-Eintrag wenn `--purge`
  requested aber nicht ausgeführt (T0-(h)).
- ✅ **Mapper**: neuer `mapRemoveErrorToDiagnostic(err)` mit
  Switch-Order-Pflicht (FS-first analog init). Heutige Sentinels
  (`ErrServiceUnsupported`, `ErrServiceUnregistered`,
  `ErrServiceInconsistent`, `ErrProjectNotInitialized`,
  `ErrConfirmationRequired`) auf LH-Kennungen mappen.
- ✅ **`ErrRemoveFileSystem`-Sentinel** (NEU): existiert noch nicht
  in port/driving/removeservice.go. T2 ergänzt analog
  `ErrAddFileSystem`/`ErrInitFileSystem`/`ErrGenerateFileSystem`.
  Multi-`%w`-Wrap an allen FS-Stellen (Switch-Order-Sicherheit).
- ✅ **Allowlist-Erweiterung**: `"u-boot remove"` in
  `jsonallowlist.go`.
- ✅ **CLI-Pin-Tests**: ~15-18 Acceptance-Tests in
  `remove_acceptance_test.go` (drei JSON-Modi, NoOp-Pin,
  ManualConflict-Symmetrie, FS-Failure, ConfirmationRequired-
  Pfad, ServiceUnregistered vs. ServiceUnsupported, Idempotenz-
  Repeat-Pin).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (remove→done),
  §6.6 neue Sektion, §7 Mutations-Matrix (remove-Zeile).
- ✅ **CHANGELOG `### Added`-Eintrag** analog init/generate.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a)** **Pattern-Erbe-Disziplin festgezurrt** (R7-MED-F3-Fix
  nach 6 Iterationen): die Erbe-Disziplin ist über die anderen
  Sub-Decisions verteilt; hier als Anchor-Tabelle konsolidiert.

  | Erbe 1:1 von init/generate | Remove-spezifisch (NEU) |
  | --- | --- |
  | `driving.PreviewMode` direkt (T0-(b)) | Confirmer-Swap-Mechanismus (T0-(j), `noopConfirmer` aus M4) |
  | `RecordingFileSystem`-driven-Adapter (Pattern-Vorbild add T1-B) | `--purge`-Flag-Dimension (T0-(h)) |
  | Pure-Go LCS-Diff-Renderer (add T2) | WARN-Migration in `diagnostics[]` (T0-(g)) |
  | `previewModeFromFlags`-Mapping (init T1-B) | `delete`-Action für RemoveAll-Captures (T0-(p)) |
  | Generalisierte Helper `reportError`/`writeEnvelope`/`mapPlannedFilesToWire` (init T5/generate T5) | Neuer `ErrRemoveFileSystem`-Sentinel (T0-(d)) |
  | Multi-`%w`-Switch-Order-Pattern (init T0-(f)) | Neuer `ErrConfirmerUnavailable`-Sentinel (T0-(e) R2-F1) |
  | `mapCaptureToPlannedFiles(records, baseDir)` (add T0-(i)) | `LH-FA-ADD-007` Multi-Use ERROR+WARN (T0-(g) R5-F2) |
  | Pre-UC-Sentinel-`reportError`-Kanal (init T0-(o)/T5) | `RemoveServiceResponse.Warnings`-Feld als WARN-Source-of-Truth (T2 R7-F2) |
  | `cliJSONEnvelope.Data` + `newDataEnvelope` (generate T1, vorgezogen aus Template 9/9) | typed `removeEnvelopeData`-Struct mit `*bool`-Wrapping für `omitempty` (T0-(f) R7-F1) |
  | Path-Anchor `plannedFiles[].path` project-relativ (init T0-(k)) | Two-Phase-Capture-Semantik bei RemoveAll-Failure (T0-(p) R5-F4) |
- **T0-(b)** **`driving.PreviewMode` direkt** (kein
  `RemovePreviewMode`-Alias) — durch init-T0-(c) Alias-
  Lebensdauer-Pflicht erzwungen.
- **T0-(c)** `RemoveServiceService.fsFactory`-Form analog
  `InitProjectService.fsFactory`. **Control-Flow-Skeleton im
  `Remove()`-Wrapper** (R4-MED-F2-Fix — Plan-Vertrag pinnt
  Reihenfolge der Phasen, damit T3-Implementer keine Wahlfreiheit
  hat):

  ```
  Remove(req):
    LOCK removeMu                            # generateMu/initMu-Pattern
    if req.SilenceConfirmer { confirmerSwap() }  # konditional analog init's SilenceProgress (R7-LOW-F4); INNERHALB Lock; Mechanismus R12-F3 unten
    fs, recorder := s.selectFS(req.PreviewMode)  # recorder ist CALL-SCOPED, lokale Variable
    fsSwap(fs)                               # init's Swap-Pattern
    state := detectServiceState(s.fs, s.yaml, ...)
    if state == Unregistered:    return early (ErrServiceUnregistered)
    if state == InconsistentYAML: return early (ErrServiceInconsistent)
    if state == Deactivated:     return no-op (KEIN runPurgeGate)
    # Active / EnabledUnset / InconsistentBlock:
    if req.Purge && catalogueFor(svc).volumeOptional == false:
        warnings = append(warnings, deferredVolumesWarning(svc))   # WARN VOR Gate/Execute (R8-MED-F3-Fix)
    if req.PreviewMode != PreviewDryRun:     # T0-(h)(a) Skip-Logik
        runPurgeGate(req)
    executeRemove(...)
    captures := recorder.Captured()             # vor Unswaps drainen
    fsUnswap; confirmerUnswap; UNLOCK         # Defer-Restore-Pattern analog init
    resp.PlannedFiles = mapCaptureToPlannedFiles(captures, req.BaseDir)
    return response
  ```

  **WARN-Emission-Ort festgezurrt** (R8-MED-F3-Fix): die WARN-
  Diagnostic für `--purge`-mit-Volume-Service wird **VOR**
  `runPurgeGate` und **VOR** `executeRemove` an `warnings`
  angehängt (lokale Variable im Wrapper). Das pinnt zwei
  Eigenschaften: (1) der WARN-Eintrag ist auch im
  `ErrConfirmationRequired`-Pfad sichtbar (User weiß: dein
  `--purge` wäre eh deferred geworden), (2) der WARN-Eintrag ist
  auch im Mid-Write-Failure-Pfad vorhanden, aber Variante A
  (R4-F3) unterdrückt ihn dann zugunsten des Error-Diagnostics —
  T5 mapped resp.Warnings ins Envelope NUR wenn kein
  Error-Diagnostic existiert. Pattern-Pin im T6: zwei separate
  Tests für `--purge --no-interactive --json` (kein `--yes`) und
  `--purge --yes --json` Mid-Write-Failure — beide zeigen ohne
  WARN-Diagnostic im Envelope, aber aus unterschiedlichen
  Gründen.

  **`confirmerSwap()`-Mechanismus festgezurrt** (R12-MED-F3-Fix):
  Service-Field-Mutation mit defer-Restore analog init's
  `s.progress`-Swap (`initproject.go:345-349`):

  ```go
  if req.SilenceConfirmer {
      prevConfirmer := s.confirmer
      s.confirmer = noopConfirmer{}
      defer func() { s.confirmer = prevConfirmer }()
  }
  ```

  Plan-Vertrag: `s.confirmer` wird MUTIERT (KEIN lokales
  `effective`-Var), so dass `runPurgeGate` (Z. 158-178) ohne
  Signature-Change die geswappte Form sieht. Lokale-Variable-
  Variante (lokales `effective := noopConfirmer{}` + Signature-
  Change auf runPurgeGate) ist **verworfen** — bräuchte
  `runPurgeGate`-Refactor und bricht Pattern-Erbe init.

  **Recorder-Lebensdauer-Invariante (R6-MED-F2-Fix)**: `recorder`
  ist eine **lokale Variable im Wrapper**, NICHT ein Service-Feld
  (analog `initproject.go:336` + `addservice.go:364`). Pattern:
  pro Aufruf liefert die `fsFactory`-Closure (`cmd/uboot/main.go:
  130-180`) eine **frische** `recordingfs.New(...)`-Instanz —
  call-scoped. `recordingfs.Captured()` (`recordingfs.go:105-112`)
  ist NICHT thread-safe; bei Service-Feld + parallelen Goroutinen
  würden Captures der einen Goroutine in die Response der anderen
  leaken. Plan-Vertrag: T6-Pin
  `TestRemove_ConcurrentInvocationsSerializeSwaps` assertion
  ergänzt, dass `resp1.PlannedFiles` und `resp2.PlannedFiles`
  **disjunkte Capture-Sets** sind.

  **Race-Sicherheit (R5-HIGH-F1-Fix)**: ALLE Swaps (`confirmerSwap`,
  `fsSwap`) laufen **INNERHALB** der `removeMu`-Lock-Region. Außerhalb
  der Lock-Region könnte eine parallel laufende Goroutine zwischen
  Swap und Lock-Acquisition ihren eigenen Swap durchführen — beide
  Goroutinen würden auf demselben `s.confirmer`/`s.fs`-Field race.
  Pattern-Erbe init's `Init()`-Wrapper
  (`initproject.go:328-348`): erst Lock, dann `s.fs`/`s.progress`-
  Swap mit `defer`-Restore innerhalb des Lock-Scopes.
  T6-Pin: `TestRemove_ConcurrentInvocationsSerializeSwaps` mit zwei
  Goroutinen die parallel Remove() auf demselben Service-Instance
  rufen — Confirmer/FS-State darf nicht race-corruption zeigen.

  **`detectServiceState` läuft INNERHALB der Swap-Region** — sonst
  würde der Recorder die Read-Captures (compose.yaml / .env.example
  / u-boot.yaml) nicht sehen, und ein Mid-Read-Failure im
  Capture-FS würde stillschweigend mit dem Real-FS arbeiten.
  T6-Pin: `TestRemove_DryRun_DetectStateUsesCaptureF S` mit
  Spy-Read-Counter.

  **`runPurgeGate` läuft mit dem geswappten `s.confirmer`**
  (entweder `noopConfirmer` bei `SilenceConfirmer=true` oder dem
  echten Confirmer), aber NUR wenn `PreviewMode != PreviewDryRun`.
  T6-Pin: `TestRemove_DryRunPurgeYes_NoConfirmerCall` mit
  Confirmer-Call-Counter == 0.
- **T0-(d)** `ErrRemoveFileSystem`-Sentinel-Einführung +
  Wrap-Audit (R2-MED-F3 + R4-HIGH-F1-Korrektur): heute Single-
  `%w` an **8 FS-Wrap-Stellen** in `removeservice.go` (NICHT
  ~6, NICHT 10):
  - **Write/Remove** (2 Stellen): Z. 235 (WriteFile), Z. 241
    (RemoveAll extraFiles).
  - **Read/Exists/Stat** (6 Stellen): Z. 272 (Exists compose/env),
    Z. 282 (ReadFile compose/env), Z. 286 (Lstat compose/env),
    Z. 321 (ReadFile u-boot.yaml), Z. 325 (Lstat u-boot.yaml),
    Z. 358 (Exists extraFiles).

  T3 migriert alle 8 auf Multi-`%w` mit `ErrRemoveFileSystem`-
  Sentinel analog init's `initproject.go:925/967/1015/1117/1143`-
  Pattern. T6 ergänzt Read-Pfad-FS-Failure-Pin (mindestens einer)
  damit die Switch-Order-Garantie nicht löchrig wird.

  **Fachlich-klassifizierte Wraps** (2 Stellen, NICHT auf
  `ErrRemoveFileSystem`):
  - Z. 304 (managedblock-malformed): wrappt bereits korrekt
    `ErrServiceInconsistent` → bleibt unverändert.
  - Z. 307 (default-Branch im managedblock-Scanner): wraps
    unexpected Scanner-Error. R4-Korrektur: **T3 wrappt mit
    `ErrServiceInconsistent`** analog Z. 304 (gleicher Marker,
    gleicher Fail-Modus, gleiche Datenkonsistenz-Klasse).
    NICHT auf `ErrRemoveFileSystem` — Z. 307 ist KEIN
    `s.fs.*`-Aufruf, sondern ein Scanner-Format-Defekt am
    Managed-Block.
  - Z. 330 (`s.yaml.PatchScalar`-Failure): YAML-Codec-Fehler,
    KEIN FS-I/O. R4-Korrektur: **T3 wrappt mit
    `ErrServiceInconsistent`** (Datenkonsistenz-Klasse —
    invalides YAML-Schema). Alternative wäre ein neuer
    `ErrYAMLPatch`-Sentinel; weil aber heute nur der eine
    Codec-Wrap-Pfad existiert und Exit 10 / Fachlich-Klasse
    semantisch passt, ist die Konsolidierung auf
    `ErrServiceInconsistent` die schlankere Lösung.

  Inkonsistenz im initialen Inventar (R3-MED-F3 hatte alle 10
  als "FS-Wrap-Bucket" gelistet, obwohl Z. 307 + Z. 330 fachlich
  sind) ist mit dieser Aufteilung behoben.

  **`ErrServiceInconsistent` Triple-Use-Klarstellung**
  (R5-MED-F3-Fix): Der Sentinel deckt nach R4-F1-Fix drei
  Sub-Semantiken:
  (a) Z. 304: Managedblock-Marker malformed (BEGIN ohne END /
      duplicate BEGIN) — kanonischer ErrServiceInconsistent-Sinn.
  (b) Z. 307: Scanner-Default-Branch (unexpected scanner state) —
      ähnlicher Marker-Defekt, Inhalts-Semantik-Stretching klein.
  (c) Z. 330: YAML-Patch-Failure (`s.yaml.PatchScalar` failt auf
      `services.<name>.enabled`-Pfad) — KEIN managed-block-Defekt;
      Sentinel-Stretching auf "YAML-Schema-Inkonsistenz".
  Alternative wäre ein eigener `ErrYAMLPatchFailed`-Sentinel mit
  eigener LH-Zeile in T0-(e). T3-Implementer entscheidet zwischen
  (i) Konsolidierung auf `ErrServiceInconsistent` (schlanker, aber
  semantisch dehnend) oder (ii) Sub-Sentinel einführen (sauberer
  aber +1 Tabellen-Zeile + Plan-Anpassung). Plan-Vorschlag:
  Konsolidierung — die drei Sub-Semantiken teilen sich
  Exit-Code 10 und Repair-Hint (User-Action: YAML manuell
  reparieren), nur die Message-Text-Differenzierung bleibt. T6-Pin
  pinnt den `LH-FA-ADD-005`-Code für alle drei Pfade.
- **T0-(e)** **Switch-Order-Pflicht** im neuen
  `mapRemoveErrorToDiagnostic`. Diagnostic-Code-Tabelle (T6-Pin-
  Pflicht pro Zeile):

  | Sentinel | LH-Code | Exit |
  | --- | --- | --- |
  | `ErrRemoveFileSystem` | `LH-NFA-REL-003` | 14 |
  | `ErrConfirmerUnavailable` (NEU, T2; wrappt heutigen string-Error in `removeservice.go:171`) | `LH-FA-CLI-005A` | 10 |
  | `ErrConfirmationRequired` | `LH-FA-INIT-005` | 10 |
  | `ErrServiceUnsupported` | `LH-FA-ADD-002` | 10 |
  | `ErrServiceUnregistered` | `LH-FA-ADD-007` | 10 |
  | `ErrServiceInconsistent` | `LH-FA-ADD-005` | 10 |
  | `ErrProjectNotInitialized` | `LH-FA-ADD-001` | 10 |
  | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 |
  | `ErrConflictingModeFlags` (`--yes ⊕ --no-interactive`) | `LH-FA-CLI-005A` | 2 |
  | `cli.ErrServiceNameMissing` (NEU **T5**, R11-HIGH-F1 + R12-HIGH-F1: Layer-Heim CLI-Adapter analog `cli.ErrConflictingModeFlags` `cli/cli.go:177`; KEIN driving-Sentinel weil Form-Validierungs-Sentinel vom CLI-Adapter emittiert, vom Use-Case nie gesehen) | `LH-FA-CLI-006` | 2 |
  | Default (unknown) | `LH-FA-CLI-006` | 1 |

  **Tabellen-Reihenfolge = Switch-Reihenfolge** (R3-MED-F3-Fix):
  Infrastruktur-Sentinels (`ErrRemoveFileSystem`,
  `ErrConfirmerUnavailable`) stehen VOR `ErrConfirmationRequired`
  und fachlichen Service-Sentinels, damit Multi-`%w`-Wraps (Go
  1.20+) nicht versehentlich auf einen fachlich-Branch matchen.
  T6-Pin verifiziert: ein konstruierter
  `fmt.Errorf("%w: %w", ErrConfirmerUnavailable,
  ErrConfirmationRequired)` MUSS `LH-FA-CLI-005A` / Exit 10
  (I/O-Klasse), NICHT `LH-FA-INIT-005`. **Defense-only-Pin**
  (R4-LOW-F5-Klarstellung): heute existiert KEIN Code-Pfad der
  beide Sentinels gemeinsam chained — der Pin verifiziert die
  Mapper-Robustheit gegen einen synthetisch konstruierten
  Multi-Wrap, nicht ein reales Failure-Szenario. Cluster-T_close
  kann eine generische `mapErrorToDiagnostic`-Registry die
  Multi-`%w`-Resilienz cluster-übergreifend pinnen.

- **T0-(f)** **Envelope-`data`-Form festgezurrt**: Success-Envelope
  trägt `data: {"service": "<…>", "priorState": "<…>", "state":
  "<…>", "volumesPurged": <bool>}`. Error-Envelope trägt nur
  `data: {"service": "<…>"}` ohne `priorState`/`state`/
  `volumesPurged` (Zero-Response auf Error-Pfad — analog
  generate T0-(q)).

  **Struct-Form festgezurrt** (R7-HIGH-F1-Fix, Pattern-Erbe-Wahl
  zwischen init und generate): generate-Symmetrie — eine typed
  Struct `removeEnvelopeData` in `cli/remove.go` analog generate's
  `generateEnvelopeData`. Felder mit **Pointer-Wrapping**, damit
  `omitempty` Key-Abwesenheit pinnen kann (nicht nur Zero-Value):

  ```go
  type removeEnvelopeData struct {
      Service       string  `json:"service"`
      PriorState    *string `json:"priorState,omitempty"`
      State         *string `json:"state,omitempty"`
      VolumesPurged *bool   `json:"volumesPurged,omitempty"`
  }
  ```

  Begründung für `*bool` statt `bool`: `bool`-`omitempty` würde
  den Error-Pfad-`false`-Wert UND den Success-Pfad-`false`-Wert
  identisch droppen (Spec §1841 fordert Key-Presence-vs-Absence-
  Disambiguierung). Pattern-Erbe `cliJSONEnvelope.DryRun` /
  `Diff` (`jsonenvelope.go:34-37`) nutzt `*bool` aus genau diesem
  Grund. **Konsistenz-Begründung vs. generate's
  `Action string`-Pattern** (R8-LOW-F4): generate kommt mit
  `string`-`omitempty` aus, weil `Action` einen klaren Empty-
  Marker hat (Action ist `""` nur im Error-Pfad, Success-Strings
  sind alle non-empty). Für `PriorState`/`State` gilt das auch
  (Strings `"active"`/`"deactivated"`/`"enabled-unset"`), aber
  `*string` ist defensiver gegen Default-`""`-Drift in
  künftigen Domain-Erweiterungen. **`VolumesPurged` MUSS
  `*bool`** weil `false` ein valider Success-Wert ist. Pattern-
  Wahl konsistent: alle drei Felder Pointer, weil Symmetrie-
  Pflicht gegen `VolumesPurged`-Vorgabe.

  **T6 Helper-Gap-Notiz** (R8-HIGH-F1-Fix): Plan-AK pinnt Key-
  Absence für Error-Envelope (`volumesPurged` MUSS abwesend sein,
  NICHT `false`). `jsontestutil.AssertFullEnvelope` hat heute
  KEINE `WithDataKeyAbsent`/`WithDataShape`-Option
  (`jsontestutil.go:22-27` dokumentiert die kleine Helper-API
  explizit als „vier Options decken alle Pin-Wünsche der Folge-
  Slices ab"). Init/generate-Acceptance-Tests inspizieren `data`
  per manuellem `env["data"].(map[string]any)["..."]`-Cast
  (`init_acceptance_test.go:98, 134, 171, 226, 263`). **T6-
  Implementer-Wahl**:
  (a) **Helper-Erweiterung** in `jsontestutil` (eigene Sub-Tranche
      T6-A: `WithDataKeyAbsent("volumesPurged")` /
      `WithDataKeyPresent("volumesPurged", false)`) — saubere
      Acceptance-API, Pattern-Vorlauf für künftige Slices die
      `data`-Sub-Keys pinnen.
  (b) **Manuelle Cast-Form** analog init/generate-Pattern —
      kein Helper-Refactor, aber AK-Wortlaut "T6-Pin verwendet
      `AssertFullEnvelope`-Key-Presence-Assertion" muss auf
      "T6-Pin inspiziert `env["data"]`-map manuell" umformuliert.
  Plan-Vorschlag: **Variante (a)** — Helper-Erweiterung in
  T6-A-Sub-Tranche (+~50 LOC für 2 Options + Tests). Pattern-
  Symmetrie zu `WithCommand`/`WithExpectedCodes` ist sauber, und
  Folge-Slices (down, doctor-update) könnten den Helper erben.
- **T0-(g)** **WARNING-Migration festgezurrt** (R1-HIGH-1-Fix +
  R2-MED-F5-Erweiterung): heutige `printRemoveSummary`-stderr-
  WARNING (Z. 163-171) bei `--purge && !VolumesPurged` wandert
  im JSON-Mode in `diagnostics[]` mit `level: "warn"` und Code
  **`LH-FA-ADD-007`** (Spec §924 / §2602 — die Anforderung selbst
  beschreibt das deferred-Volumes-Verhalten). KEIN Suffix-Schema
  wie `-VOLUMES-DEFERRED` — Spec §1834 erlaubt nur die feste
  `LH-<Bereich>-<Modul>-<3-stellige-Zahl>`-Form oder tool-interne
  Codes mit Doku-Pflicht; ein freier Suffix verletzt das Schema.
  Differenzierung zur Confirmation-Required-Diagnostik läuft über
  den `level: "warn"`-vs-`"error"`-Vertrag plus den
  `message`-Text plus `data.volumesPurged`-Status. Pinnbar via
  `jsontestutil.AssertFullEnvelope`.

  **`LH-FA-ADD-007` Multi-Use-Klarstellung** (R5-MED-F2-Fix):
  Plan nutzt `LH-FA-ADD-007` an zwei Stellen:
  (1) **ERROR-Diagnostic-Code** für `ErrServiceUnregistered`
      (Mapper-Tabelle T0-(e)) — Exit 10.
  (2) **WARN-Diagnostic-Code** für `--purge && !VolumesPurged`
      (T0-(g)) — Exit 0 oder Exit 14 (bei Mid-Write-Variante).
  Beide referenzieren das Spec-Umbrella `LH-FA-ADD-007 "Service
  entfernen"` (§924-947) — der Code identifiziert die *Anforderung*,
  nicht die Sub-Semantik. Spec §1834-Vertrag erlaubt das, weil
  `diagnostics[].level` (warn vs error) und ggf. `message`-Text
  + `data.volumesPurged`-Status den konkreten Sub-Pfad
  disambiguieren. Konsumenten dürfen NICHT nur auf `code` filtern,
  sondern müssen `(code, level)` als Tupel betrachten. Pattern
  ist konsistent mit init/down's `LH-FA-INIT-005`-Multi-Use für
  `ErrConfirmationRequired` (geteilt zwischen init und down).
  WARN-Diagnostic wird NUR emittiert wenn der Catalog-Entry
  tatsächlich ein named volume deklariert. Echter Field-Name auf
  der `serviceCatalogueEntry`-Struct
  (`addservice_execute.go:190-224`): **`volumeOptional bool`**
  — `false` heißt "Service hat ein named volume" (heute nur
  postgres mit `volumeOptional: false`, `volumeRefLiteral:
  "postgres-data"`). Keycloak und OTel sind `volumeOptional:
  true` → kein named volume → keine WARN. KEIN erfundenes
  `Volumes`-Feld; T3-Implementierung nutzt
  `catalogueFor(svc).volumeOptional == false` als Check.
  Pattern-Vorlauf bleibt für künftige `volumeOptional: true`-
  Catalog-Entries (Keycloak/OTel sind heute Volumeless-Beispiele,
  Mock nicht nötig). T6-Pin:
  `TestRemove_PurgeOnVolumelessService_NoWarn` mit keycloak oder
  otel als realistische Volumeless-Fixture → `diagnostics: []`,
  kein WARN-Eintrag, `data.volumesPurged: false`.
- **T0-(h)** **`--purge`-in-Dry-Run-Verhalten festgezurrt** (R1-
  HIGH-3-Fix): Dry-Run impliziert Null-Mutationen. Drei
  Vertragsränder gepinnt:
  (a) **Confirmer-Gate-Skip**: `if req.PreviewMode ==
      PreviewDryRun { skip confirmer gate }` — auch ohne `--yes`
      keine `ErrConfirmationRequired` im Dry-Run (analog init's
      `initGit`-Skip T0-(n)).
  (b) **Envelope-Form im Dry-Run**: `data.volumesPurged: false`
      IMMER (deferred unabhängig vom Gate-Skip);
      `diagnostics[]` enthält EINEN `warn`-Eintrag mit Code
      `LH-FA-ADD-007` (T0-(g) WARN-Migration) wenn
      `req.Purge && req.PreviewMode != PreviewNone` — dem User
      ist klar: Purge wurde requested, im Dry-Run aber
      semantisch geskippt.
  (c) **Diff-Pfad rendert KEINE Volume-Aktion**: Volume-Removal
      ist nicht-FS-Side-Effect; der Recorder capturet nur FS-
      Mutations. `changes[]` enthält ausschließlich
      FS-Captures; die Purge-Side-Effect-Information lebt
      ausschließlich in `data.volumesPurged` + `diagnostics[]`-
      WARN.

  T6-Pflicht: 1 Test pro `--purge`-on-Variante in jedem der
  vier Dry-Run-Kombos (Dry-Run, Dry-Run+Diff, Dry-Run+JSON,
  Dry-Run+Diff+JSON).

  **PreviewAndApply-Branch festgezurrt** (R2-MED-F4-Fix): bei
  `--purge --diff` ohne `--dry-run` (PreviewMode=PreviewAndApply)
  läuft der Confirmer-Gate REGULAR — Diff-Mode schreibt echt auf
  Disk. T6-Pins für PreviewAndApply + `--purge`:
  (a) `--purge --diff --json` ohne `--yes` → noopConfirmer
      (T0-(j)) → `ErrConfirmationRequired` Envelope, Exit 10,
      kein `changes[]` (Plan-Phase failt vor Execute).
  (b) `--purge --diff --json --yes` → Gate skipped, Execute läuft,
      Voll-Envelope mit `changes[]` der FS-Captures + WARN-
      Diagnostic (T0-(g)) für `data.volumesPurged: false`,
      `status: warn`, Exit 0.
  (c) `--purge --diff --no-interactive` ohne `--json` ohne `--yes`
      → `ErrConfirmationRequired`-Pfad (heutiges Verhalten,
      stderr-Print). T5 muss diesen Pfad NICHT in JSON-Helper
      kanalisieren (kein `--json`-Flag) — bleibt Cobra-Default.
- **T0-(i)** **`--purge`-Mutex mit `--dry-run`/`--diff`?** Analog
  init's `--template`-Mutex (T0-(i))? Vorschlag: **NEIN** — Purge
  ist eine Side-Effect-Dimension, kein Renderer-Pfad. Dry-Run +
  Purge ist semantisch konsistent: "zeige was Remove + Purge
  ändern WÜRDE", auch wenn der Gate-Run skipped wird.
- **T0-(j)** **Confirmer-Swap-Mechanismus (NEU, R1-HIGH-2-Fix +
  R3-HIGH-F2-Korrektur)**: NEU ist nur der **Swap-Mechanismus**
  (request-time statt construction-time) — der `noopConfirmer`-
  Helper selbst existiert bereits in `application/noop.go:17-33`
  (M4 Confirmer-Port-Slice; `RemoveServiceService.NewRemove…`
  `removeservice.go:48` nutzt ihn schon als nil-Fallback). Init
  swappt nur ProgressPort, nicht Confirmer — der Swap-Mechanismus
  ist hier neu, der Helper ist geerbt. Form:
  `req.SilenceConfirmer = flags.JSON`. Bei `--purge --json` ohne
  `--yes`: ConfirmerPort wird auf den existierenden
  `noopConfirmer` umgeswapt der `false, nil` returnt —
  `runPurgeGate` (removeservice.go:173-176) wandelt das in
  `ErrConfirmationRequired`. **Semantik-Klarstellung**: das ist
  KEIN Silencing (keine UX-Information-Verlust-Symmetrie zu
  noopProgress), sondern eine **bewusste Behaviour-Change** im
  JSON-Mode — User muss explizit `--yes` setzen um im
  JSON-Mode zu purgen. Pattern-Erbe-Disziplin T0-(a) Spalte
  führt nur den Swap-Mechanismus als remove-spezifisch.

  **`--purge --yes --json`-Pfad** (R1-MED-5-Fix): bei
  `req.Yes==true` skipped runPurgeGate (Z. 162-164) ohne
  Confirmer-Call → Execute läuft durch → `VolumesPurged: false`
  (v0.3.0 deferred). Plan-Vertrag: trotzdem WARN-Diagnostic
  emittiert (T0-(g)), aber `exitCode: 0` UND `status: warn`
  (Spec §447-Kopplung) — Warn-only verschiebt nicht den
  Exit-Code. T6-Pin: `TestRemove_PurgeYesJSON_WarnOnly` mit
  `status: warn`, exit 0, `data.volumesPurged: false`.

  **`--purge --yes --json` PLUS Mid-Write-Failure-Variante**
  (R4-MED-F3-Fix, Doppel-Diagnostic-Klärung): wenn die Execute-
  Phase mid-write failt (z. B. compose.yaml WriteFile-Error vor
  yaml.WriteFile), wird **Variante A** festgezurrt: Error-
  Diagnostic dominiert, WARN unterdrückt. Envelope: `diagnostics:
  [{level: "error", code: "LH-NFA-REL-003", file: "<…>"}]`,
  `status: error`, exit 14, `data: {"service": "<…>"}` ohne
  `volumesPurged` (Zero-Response analog T0-(f) Error-Pfad).
  Begründung: WARN über `volumesPurged: false` würde sich auf
  ein nicht-existentes Datenfeld beziehen — die Zero-Response-
  Klausel zieht den `data`-Bereich konsistent auf den Diagnostics-
  Channel. T6-Pin: `TestRemove_PurgeYesJSON_MidWriteFailure_
  ErrorOnly`.
- **T0-(k)** Path-Anchor: `plannedFiles[].path` ist project-
  relativ (analog init T0-(k)) — `mapCaptureToPlannedFiles(records,
  baseDir)`-Erbe.
- **T0-(l)** **Allowlist-Form**: parent-only `"u-boot remove"`
  (Service ist Positional-Arg analog generate T0-(l)).
- **T0-(m)** **Envelope-Shape**: `command="remove"`, kein
  `subcommand`-Feld, Service-Name in `data.service`.
- **T0-(n)** **`Codes`-Registry**: KEINE Ergänzung nötig
  (LH-Codes sind generisch erlaubt; nur §6.6-Doku pflegt die
  remove-Sektion).
- **T0-(o)** Pre-`next/`-Review-Runden-Erwartung (R10-LOW-F3-
  Update auf faktischen Stand): 10 Runden (Discovery + 9
  Adversarial mit unterschiedlichen Angles). HIGH-Frequenz-
  Asymptote stabil bei 1/Runde über R5-R10 (sechs Runden in
  Folge). **Cluster-Pattern-Erbe für Folge-Slices 6/9-9/9**:
  Asymptote-Detektion ist legitimes Konvergenz-Kriterium für
  `next/`-Übergang; ≥ 5 Runden mit konstanter HIGH-Frequenz
  bestätigen, dass die verbleibenden Befunde Cosmetic-
  Präzisierungen sind, keine Substanz-Lücken.
- **T0-(p)** **`delete`-Action-Vertrag (NEU, R4-HIGH-F4 +
  R5-LOW-F4-Erweiterung + R9-HIGH-F1-Klarstellung)**:
  remove ist der **erste end-to-end-sichtbare** `PlannedFile.
  Action == "delete"`-Produzent (Capture wandert via
  `mapCaptureToPlannedFiles` in den Wire-Envelope). **Layer-
  Klarstellung**: der `delete`-enum-Wert existiert auf Spec-
  Layer bereits (`cli-json-output.md` §164-166:
  `enum: ["create", "modify", "delete"]`); der Recorder-Layer
  hat `actionDelete` seit Rename/RemoveAll-Support
  (`recordingfs.go:35-37`). NEU ist nur, dass ein Use-Case den
  `delete`-Capture in `mapCaptureToPlannedFiles` durchreicht
  und damit in `data.changes[].action: "delete"` end-to-end
  sichtbar wird. Init/add/generate produzieren `actionDelete`-
  Recorder-Captures NICHT (kein RemoveAll/Rename in deren
  Use-Case-Pfaden); remove tut es via `RemoveAll` Z. 241 für
  extraFiles. — `RemoveAll` (Z. 241) für extraFiles wird
  von `recordingfs.RemoveAll` (`recordingfs.go:197-206`) mit
  `Action: actionDelete` capturet, `mapCaptureToPlannedFiles`
  mapped das auf den Spec-§354-Wert `"delete"`. Init/add/generate
  produzieren nur `create`/`modify`. Diff-Renderer-Behavior:
  `delete` = reiner Old-only-Hunk (`OldContent` voll, `NewContent`
  leer → full-file-Remove-Block). Plus: `OldContent` für
  RemoveAll wird via `recordingfs.snapshot` über `ReadFile`
  geladen — für regular files (z. B. otel-Catalog
  `extraFiles`-Eintrag) OK, für Dir-Trees aber `nil`. Heute kein
  Risiko (postgres/keycloak/otel haben File-extraFiles), aber
  Pattern-Vorlauf für künftige Dir-extraFiles wäre out-of-scope.

  **Mid-Stream-Capture-Semantik bei `RemoveAll`-Failure**
  (R5-LOW-F4-Klarstellung): wenn `RemoveAll` mid-stream failt,
  läuft `recordingfs.RemoveAll` (Z. 197-206) in zwei Phasen:
  (1) `snapshot` (via internal `ReadFile`) → setzt `OldContent`.
  (2) `underlying.RemoveAll` → mutiert echte Disk (im
      Passthrough-Modus) oder no-op (Dry-Run-Modus).
  Wenn Phase 1 failt (File unreadable, Permission denied), bleibt
  `OldContent` leer und der Capture trägt `Action: "delete"` +
  `OldContent: nil`. Wenn Phase 2 failt nach erfolgreichem
  Snapshot, ist Capture vollständig (OldContent gesetzt). Beide
  Pfade kommen in `plannedFiles[]` als gleichberechtigte Einträge
  vor. Diff-Renderer behandelt `OldContent: nil` als "binary or
  unreadable" — leerer Hunk-Body. T6-Pin: `TestRemove_OtelExtra
  FileDelete_DiffHasDeleteHunk` prüft `data.changes[].action ==
  "delete"` plus den Unified-Diff-Body. `cli-json-output.md` §7
  Mutations-Matrix und §6.6 dokumentieren `action: "delete"`
  explizit.
- **T0-(q)** **`baseDirSanitizedError`-Wrapper für `diagnostic.
  message`** (R14-MED-1-Fix in T7-Commit `4fb3fea`, T8-Robustheits-
  Erweiterung via R15-LOW-1): Use-Case-Layer wrappt FS-Fehler an
  einer absoluten Pfad-Stelle
  (`application/removeservice.go:410` mit `f.path`) plus die raw-
  FS-Error-Component carriet den absoluten Pfad mit; ohne
  Sanitisierung würde `mapRemoveErrorToDiagnostic` den Pfad 1:1
  in `diagnostic.message` lassen — Info-Leak des User-Filesystem-
  Layouts. **Wrap-Site-Inventur-Korrektur** (R15-LOW-1 Audit):
  von 8 FS-Wrap-Stellen in T0-(d) Inventar tunneln nur die direkte
  absolute-Pfad-Site plus den raw FS-Error den Pfad mit — die
  anderen sieben Wrap-Sites führen `f.relPath`/`filename`/
  `xf.Path` oder Literal `u-boot.yaml`. Sanitisierung bleibt
  trotzdem nötig wegen der raw-Error-Component. **Wrapper-Form**:
  package-private `baseDirSanitizedError struct { inner error;
  baseDir string }` in `cli/remove.go:465-491` mit `Error()`-
  Method-Override und `Unwrap() error`-Method, damit `errors.Is`/
  `As` über Unwrap-Chain (single + multi-`%w`-Wrap Go 1.20+)
  intakt bleiben. **Sanitisierungs-Regeln**:
  (a) `<baseDir>/foo` → `foo` via `strings.ReplaceAll(msg,
      baseDir+sep, "")` — Path-Separator-Suffix garantiert
      Wort-Ende.
  (b) bare `<baseDir>` → `.` via `replaceBareBaseDir` — Word-
      Boundary-Check (gefolgt von End-of-String oder Nicht-Pfad-
      Byte, ASCII-only). R15-LOW-1-Robustheit: substring-Replace
      ohne Boundary würde `<baseDir>-cache/lock` zu `.-cache/lock`
      mangeln (T8 fixt das in `4fb3fea`-Nachfolge-Commit, dem
      T8-Closure-Commit).
  (c) leerer baseDir → unverändert (defensive identity).
  Aufrufer: `runRemove` ruft `sanitizeBaseDir(removeErr, cwd)`
  vor `reportError`. **Use-Case-Layer-`relativizePath`-
  Asymmetrie**: die Wire-Form-`plannedFiles[].path`-Sanitisierung
  läuft via `filepath.Rel` (`application/addservice.go:318-327`)
  — `Rel` arbeitet auf einzelnen Pfad-Strings; für Error-
  Messages mit eingebettetem Pfad innerhalb von Prosa ist
  Word-Boundary-Substring die direkte Form. Mapper-Switch-Order,
  Sentinel-Identity, Exit-Code-Mapping — alles unverändert.
  T6-Pin `TestRemove_FSErrorWithAbsolutePath_SanitizesMessage`
  pinnt Pre-T7-Path-Leak; T8-Pin
  `TestRemove_FSErrorWithBaseDirSubstring_NotMangled` pinnt die
  R15-LOW-1-Boundary-Robustheit (Konstruktor-Fixture nutzt
  baseDir + `-cache/lock` als Sibling-Pfad, der **nicht**
  gemangled werden darf).
- **T0-(r)** **`validateRemoveArgs`-Flag-Awareness** (R13-HIGH-1-
  Fix in T7-Commit `4fb3fea`, T8-Coverage-Erweiterung via
  R15-LOW-2): Spec §1841/§1842 verlangt für jeden modifying-
  Subcommand im `--json`-Mode einen Envelope auf stdout —
  Minimal-Schema im Default-Pfad, Voll-Schema sobald `--dry-run`
  oder `--diff` gesetzt ist. Custom-`Args`-Validator
  `validateRemoveArgs(a *App)` (`cli/remove.go:193-213`) bekommt
  die Flag-Awareness via `cmd.Flags().GetBool("dry-run"/"diff")`
  zur Validator-Zeit. **Cobra-Parse-vor-Validate-Ordering** (R15-
  Angle-A-Defense): Cobra v1.10.2 `command.go:861-909`
  garantiert `ParseFlags` läuft vor `ValidateArgs` — der Validator
  sieht die geparsten Flag-Values zuverlässig. `--dry-run` und
  `--diff` sind als **Subcommand-Local-Flags** auf `remove`
  registriert (`cli/remove.go:146-149`), nicht persistent; der
  Tree-Lookup in `cmd.Flags()` liefert die Local-Form direkt.
  **Drei Pfade**: `len(args)==1` ok (durch zu RunE),
  `len(args)==0` → emit `LH-FA-CLI-006`-Envelope auf stdout
  (Minimal- oder Voll-Schema je Flag-State) + return
  `ErrServiceNameMissing`, `len(args)>1` → emit Envelope mit
  Cobra-Roh-Error + return Cobra-Error (Exit 2 via
  `isUsageError`-`"accepts "`-Prefix). T6-Pins
  `TestRemove_NoPositionalArg_DryRunJSON_EmitsFullSchemaEnvelope`
  (R13-HIGH-1) + `TestRemove_TooManyArgs_JSON_EmitsCLI006Envelope`
  (R13-MED-1). **T8-Coverage-Pin** R15-LOW-2:
  `TestRemove_TooManyArgs_DryRunJSON_EmitsFullSchemaEnvelope`
  pinnt explizit die Kombi-Pfad-Coverage `len(args)>1 +
  --dry-run + --json` → Voll-Schema. Future-Regression-Defense
  bei Validator-Refactoring.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | ------ | --------------- | --- |
| T0 | Discovery + Sub-Decisions (a)-(o) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (R3-HIGH-F2-Fix): `noopConfirmer` existiert bereits seit M4 Confirmer-Port-Slice in `application/noop.go:17-33` und tut exakt was T0-(j) braucht (`ConfirmRemoveVolumes → false, nil`). `RemoveServiceService`-Konstruktor (`removeservice.go:48`) nutzt ihn schon als nil-Fallback. T3 swappt den existierenden Helper request-time, kein neuer Helper nötig. | — (entfällt) | T0 |
| T2 | Port-Types: `RemoveServiceRequest.PreviewMode` + `SilenceConfirmer`-Feld, `RemoveServiceResponse.PlannedFiles`/`Changes`-Felder, **`RemoveServiceResponse.Warnings []driving.WarningEntry`-Feld** (R7-MED-F2-Fix + R8-MED-F2-Type-Klärung: neuer Port-Type `driving.WarningEntry struct { Code string; Level string; Message string; Subject string \`json:",omitempty"\` }` analog `diagnosticItem`-Wire-Form — Layer-sauber (KEIN `domain.Diagnostic`-Wiederverwendung, weil dessen Severity-Enum + ID-Field semantisch mismatch zum Wire-Type ist). **`Subject`-Feld** (R12-LOW-F4-Vorlauf): proaktiv eingeführt für up/down Multi-Service-WARN ("container 'postgres' will be replaced") + config-set Multi-Key-WARN; remove nutzt das Feld NICHT (`""`, omitempty), aber Pattern-Erbe für 6/9 + 8/9 ist proaktiv abgedeckt — kein breaking Type-Change im Cluster. **Cluster-Vorlauf-Disziplin** (R9-MED-F2-Fix): Type ist bewusst generisch `driving.WarningEntry` benannt, NICHT `RemoveWarningEntry`, weil up/down's recreate-Warnings und config-set's value-warnings denselben Type erben werden (Cluster-Folge-Slices 6/9 + 8/9). Erste-Slice-Pattern-Last analog `PreviewMode`-Rename in init T0-(c). Use-Case ist Source-of-Truth für WARN (Catalog-Lookup für `volumeOptional`), CLI mapped via `mapWarningsToDiagnostics(resp.Warnings) []diagnosticItem`. Triviales Field-Mapping, kein Severity-Enum-zu-String-Cast nötig.), **zwei neue Port-Sentinels** (R12-HIGH-F1-Layer-Korrektur, der dritte Sentinel `cli.ErrServiceNameMissing` lebt im CLI-Layer und wird in T5 etabliert): `ErrRemoveFileSystem` (FS-Klasse, T0-(d)) UND `ErrConfirmerUnavailable` (Confirmer-I/O-Error-Klasse, R2-HIGH-F1-Fix für T0-(e)-Tabelle). | ~120 | T0 |
| T3 | Application-Layer: `RemoveServiceService.fsFactory` + `removeMu sync.Mutex` + `NewRemoveServiceServiceWithFactory` + `Remove()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; **Multi-`%w`-Wrap an den 8 FS-Wrap-Stellen** (R4-HIGH-F1-Klassifikations-Fix gegen R3-Initial-10, T0-(d) Inventar Z. 264-285); Z. 307 + Z. 330 separat mit `ErrServiceInconsistent`-Wrap (KEIN ErrRemoveFileSystem, fachlich-Klasse, T0-(d)); `ErrConfirmerUnavailable`-Sentinel-Wrap in `runPurgeGate` Z. 171; Confirmer-Swap auf existierenden `noopConfirmer` im JSON-Mode. | ~240 | T2 |
| T4 | Composition-Root-Wiring `removeFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runRemove` ruft generische Helper mit `command="remove"`, `mapErr=mapRemoveErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; `mapRemoveErrorToDiagnostic` neu; `data`-Struct (`removeEnvelopeData`); WARNING-Migration in `diagnostics[]` (`level: "warn"`); **Pre-UC-Sentinel-Kanal** (R4-LOW-F6-Klarstellung: Codepfade existieren bereits in `cli/remove.go:108-120`, NEU ist nur die Kanalisierung via `reportError` analog `init.go:205, 216, 221`) für `domain.ErrInvalidServiceName`, `ErrConflictingModeFlags` UND `getwd`-Failure (`fmt.Errorf("determine working directory: %w", err)`, R3-LOW-F6-Fix). Der `getwd`-Wrap trägt KEIN typed Sentinel und fällt in den Default-Branch `LH-FA-CLI-006` / Exit 1 (Pattern-Erbe von init T0-(o)); Mapper-Tabelle T0-(e) NICHT ergänzt. **Human-Mode-Diff-Renderer** (R2-LOW-F6-Fix): bei `--purge --diff` ohne `--json` bleibt die deferred-Volumes-Prosa auf `errOut`, NICHT im Diff-Body. T6-Pin: `TestRemove_PurgeHumanDiff_StderrSeparation` mit getrennten Buffer-Assertions. | ~250 | T1 + T2 |
| T6 | Acceptance-Tests: ~20-25 Tests (drei JSON-Modi + NoOp Single+Repeat + Mid-Write-Failure + ConfirmationRequired-Pfade × 3 Varianten + Service-Sentinels × 4 + WARNING-Migration-Pin + `--purge`-on/off × Dry-Run-Kombos (T0-(h)) + `--purge --yes --json` WarnOnly-Pin (T0-(j) R1-MED-5) + `ErrConflictingModeFlags`-Pin). R1-MED-6-Kalibrierung: ~600-700 LOC realistisch (Confirmer-Pattern-Neumuster zieht Test-Surface). **Pin-Namen-Mapping** (R6-LOW-F4) — kanonische Tags pro Finding-Anker: `TestRemove_ConcurrentInvocationsSerializeSwaps` (R5-F1+R6-F2 Race+Recorder-Scope), `TestRemove_DryRun_DetectStateUsesCaptureFS` + `TestRemove_DryRunPurgeYes_NoConfirmerCall` (R4-F2 Control-Flow), `TestRemove_PurgeOnVolumelessService_NoWarn` (R3-F1 Volume-Presence), `TestRemove_PurgeYesJSON_WarnOnly` (R1-MED-5 + R5-F2 WARN-Pfad), `TestRemove_PurgeYesJSON_MidWriteFailure_ErrorOnly` (R4-F3 Variante A), `TestRemove_OtelExtraFileDelete_DiffHasDeleteHunk` (R4-F4 delete-Action), `TestRemove_PurgeHumanDiff_StderrSeparation` (R2-F6 stderr-Trennung), `TestRemove_ServiceUnregisteredJSON_ErrorLevelCodePin` (R6-F3 ERROR-Pfad-Symmetrie). Weitere ~12 Pins (Idempotenz-Repeat, EnabledUnset-Normalisierung, ManualConflict × 3 (R5-F3 Triple-Use), Service-Sentinels × 4, ConfirmerUnavailable-allein-Pfad (R2-F1), Multi-`%w`-Switch-Order-Defense (R3-F3), Read-Pfad-FS-Failure (R2-F3), `ErrConflictingModeFlags` (R1-F4), `TestRemove_NoPositionalArg_JSON_EmitsCLI006Envelope` (R11-F1)) lassen sich aus AK-Block + Sub-Decision-Pins direkt ableiten. **Failure-Mode-Coverage-Summary** (R11-LOW-F3): FMEA-Walk identifizierte 20 Failure-Szenarien; nach R11-Adressierung sind 18/20 explizit gepinnt + 2/20 als bewusste Out-of-Scope-Carveouts (Context-Cancellation als Status-quo-Erbe, fsFactory-NPE als Composition-Root-Defekt-Klasse). | ~650 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe). | ~80 | T6 |
| T8 | Closure (R9 präzisierte Liste): **CHANGELOG** Unreleased-Eintrag analog generate T8. **`cli-json-output.md`**: §6-Tabelle (remove→done), neue §6.6-Sektion (`u-boot remove --json` mit drei Flag-Kombos + `data`-Carrier `{service, priorState, state, volumesPurged}` + `--purge`-Mutex + WARN-Migration + EnabledUnset-Normalisierungs-Pin), §7 Mutations-Matrix-Zeile (`remove (slice-v1-cli-json-dry-run-remove): WriteFile + RemoveAll` — KEINE neue `action`-Spalte, R9-HIGH-F1: §7 ist die `driven.FileSystem`-Methoden-Tabelle, `action: "delete"` ist Spec-§354-enum-Wert der bereits dokumentiert ist; §6.6 ergänzt nur ein worked example mit `data.changes[].action: "delete"` für otel-extraFile). **`roadmap.md`**: zwei explizite Edits — (i) done-Zähler 4→5 + `remove` aus der Offen-Liste der Cluster-AP-Zelle streichen, (ii) `Nächster Schritt`-Klausel auf Folge-Slice 6/9 `up-down` umstellen (R9-LOW-F4). **Carveouts.md**: **neuer Eintrag** (kein bestehender Volume-Removal-Eintrag — Verifikation `carveouts.md` Z. 22-29 zeigt nur Paketierung/Keycloak-CI/template-list/generate-devcontainer-Half-Write); Pattern-Vorbild `slice-v2-generate-devcontainer-rollback-aware-write` mit Status `open/, on hold pending trigger`; Spec-Anker `LH-FA-ADD-007 §"Volumes nur auf explizite Anforderung"` (R3-MED-F5 + R9-LOW-F3). **`open/`-Plan-Stub**: `slice-v1-volume-auto-removal.md` mit Auslöser (real-world Volume-Reclamation-Anfrage), Out-of-Scope-V1-Bestätigung, Spec-Anker. **Slice nach `done/`** mit DoD-Hash-Tabelle analog generate T8. **Plan-Verdichtung beim Übergang nach `done/`** (R10-MED-F2-Pflicht, Plan-Länge ist mit ~900 Zeilen über generate-done ~700 + remove-Initial ~280): (i) Review-Round-Tabellen R1-R10 Adressierungs-Spalte auf einen Satz pro Finding kürzen (Aufzeichnung *was* eingearbeitet wurde, nicht *wie*); (ii) Sub-Decision T0-(c) Skeleton-Block (~80 Zeilen Inline-Code) in dedizierte H3-Sektion verschieben analog init-done Pattern; (iii) AK-Block 17 Bullets auf zentrale-Anforderungs + Adressierungs-Pin-Trennung gliedern. Pattern-Vorbild generate-done T8 verdichtet ähnlich. | — (Doku) | T7 |

LOC-Bilanz vorläufig: ~1200-1400 (R1-MED-6-Kalibrierung —
Confirmer-Swap-Pattern ist neu und nicht von init geerbt, zieht
zusätzliche Test-Surface in T6). Pattern-Erbe von init/generate
spart die FS/PreviewMode-Infrastruktur; remove-spezifisch sind
das Confirmer-Swap-Pattern (T0-(j) NEU), die WARNING-Migration
(T0-(g)) und die `--purge`-Dimensions-Coverage (T0-(h)).

## Review-Round-1 (Pre-`next/`)

Eine adversarial-orientierte Review-Runde gegen den initialen
Stub (`12c49df`). Sieben Findings (3 HIGH, 3 MEDIUM, 1 LOW),
alle adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH | `LH-FA-ADD-007-VOLUMES-DEFERRED` ist erfundene Code-Form — Spec §1834 erlaubt nur `LH-<Bereich>-<Modul>-<3-Zahl>` oder tool-interne Codes mit Doku-Pflicht; freie Suffix-Form verletzt das Schema | T0-(g) auf `LH-FA-ADD-007` korrigiert (Spec §924/§2602 trägt die deferred-Volumes-Semantik); Differenzierung läuft über `level`-Vertrag + `message`-Text + `data.volumesPurged` |
| 2 | HIGH | T0-(j) Confirmer-Silencing-Pattern als "analog init T0-(o)" deklariert — existiert in keinem Done-Slice; init swappt nur ProgressPort, nicht Confirmer | T0-(j) auf "NEUES Pattern, nicht geerbt" umgestellt; Recon-Block + Pattern-Erbe-Disziplin-Anchor entsprechend; Semantik klargestellt (Behaviour-Change im JSON-Mode, kein Silencing) |
| 3 | HIGH | `--purge --dry-run --diff`-JSON-Vertrag offen: WARN-Diagnostic-Verhalten? Diff-Volume-Visualisierung? Gate-Skip-Pfad ohne `--yes`? | T0-(h) auf drei Pin-Punkte erweitert: (a) Confirmer-Gate-Skip im Dry-Run unabhängig vom `--yes`, (b) `data.volumesPurged: false` + WARN-Diagnostic IMMER bei `req.Purge && PreviewMode != PreviewNone`, (c) Diff rendert KEINE Volume-Aktion (nicht-FS-Side-Effect). T6-Pflicht: 1 Test pro `--purge`-on-Variante in den 4 Dry-Run-Kombos |
| 4 | MEDIUM | `ErrConflictingModeFlags`-Mutex (`--yes ⊕ --no-interactive`) fehlt in T0-(e) Mapper-Tabelle | Tabelle erweitert: `ErrConflictingModeFlags → LH-FA-CLI-005A / Exit 2` |
| 5 | MEDIUM | `--purge --yes --json` Silent-Approval-Pfad nicht gepinnt — wirft `status: warn` exit-Code-Anomalie? | T0-(j) ergänzt: WARN-Diagnostic IMMER bei `purge && !VolumesPurged`, aber `exitCode: 0` UND `status: warn`. T6-Pin `TestRemove_PurgeYesJSON_WarnOnly` |
| 6 | MEDIUM | T6-LOC unterschätzt (~500 für 15-18 Tests); Confirmer-Pattern ist neu, zusätzliche Test-Surface | T6-LOC auf ~650, Test-Anzahl ~20-25; LOC-Gesamt-Bilanz auf ~1200-1400 |
| 7 | LOW | Pre-UC-Sentinels `domain.ErrInvalidServiceName` + `ErrConflictingModeFlags` werden heute via Cobra-Default-Print emittiert — verletzt JSON-stdout-Cleanliness-Pin | T5-Tranchen-Zelle ergänzt: Pre-UC-Sentinel-Kanal via `reportError`-Helper, nicht Cobra-Default |

R1-Reviewer-Note: docs-check grün; Recon-Verifikationen
(`recordingfs.RemoveAll` als Capture-Methode, `executeRemove`-
deterministische Reihenfolge, `ServiceState.String()`-Strings,
Code-LOC-Anker) bestätigt. Sub-Decisions a-o sind nach R1
konsolidiert; Confirmer-Pattern bleibt der substanziellste
Eigenleistungs-Anteil.

## Review-Round-2 (Pre-`next/`)

Adversarial-Edge-Cases + Test-Harness-Qualität gegen den
R1-gepflegten Stub (`91e4dd1`). Sieben Findings (1 HIGH,
4 MEDIUM, 2 LOW), alle adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | Confirmer-I/O-Error-Pfad (`removeservice.go:171`) fällt in Default-Mapper-Branch `LH-FA-CLI-006 / Exit 1` — semantisch falsch für I/O-Failures (`os.Stdin`-EOF, Pipe-Bruch). Spec `LH-FA-CLI-005A` §254 koppelt Confirmation-Gate-Failures an Exit 10 | Neuer Sentinel `ErrConfirmerUnavailable` in T2 ergänzt; T0-(e)-Mapper-Tabelle erweitert (→ `LH-FA-CLI-005A` / Exit 10) |
| F2 | MEDIUM | NoOp-Pin AK kollidiert mit `EnabledUnset`-State — der wäre state-transitioning (`Changed!=nil`), nicht NoOp | AK explizit: NUR `PriorState=Deactivated` qualifiziert als NoOp; **neuer EnabledUnset-Normalisierungs-Pin** als separater Test mit `priorState: "enabled-unset"`, `changes[]` non-leer |
| F3 | MEDIUM | T0-(d) Wrap-Audit nannte ~6 FS-Stellen — real sind es **10** (Z. 235, 241, 272, 282, 286, 307, 321, 325, 330, 358) | T0-(d) auf 10 Stellen kalibriert mit Code-Anker-Inventar (Write/Remove + Read/Exists/Stat + YAML-Codec + Managedblock-Scanner); Z. 304 als Nicht-FS-Wrap explizit ausgeschlossen; T6 ergänzt Read-Pfad-FS-Failure-Pin |
| F4 | MEDIUM | `--purge --diff` ohne `--dry-run` (PreviewAndApply) Gate-Vertrag offen — T0-(h) pinnte nur PreviewDryRun-Branch | T0-(h) erweitert um PreviewAndApply-Branch mit drei expliziten Pins: (a) ohne `--yes` → ErrConfirmationRequired-Envelope; (b) mit `--yes` → Execute + WARN-Diagnostic + Exit 0; (c) Non-JSON-Non-Yes-Pfad bleibt Cobra-Default |
| F5 | MEDIUM | WARN-Diagnostic-Bedingung pinnt nicht Volume-Presence-Check — bei zukünftigen Volumeless-Catalog-Entries würde WARN fälschlich emittieren | T0-(g) erweitert: WARN NUR wenn `catalogueFor(svc).Volumes != nil`; T6-Pin `TestRemove_PurgeOnVolumelessService_NoWarn` für Pattern-Vorlauf |
| F6 | LOW | Human-Mode-Diff + `--purge`-deferred-Volume-Prosa-Trennung offen — soll inline im Diff oder auf errOut? | T5-Zelle ergänzt: Human-Mode-Diff-Renderer hält die Prosa auf `errOut`, NICHT im Diff-Body; T6-Pin `TestRemove_PurgeHumanDiff_StderrSeparation` |
| F7 | LOW | T6-Code-Pin für `LH-FA-ADD-007` nicht explizit im AK — laxer Implementer pinnt nur `level: "warn"` ohne Code-Assertion | AK-WARNING-Migration-Zeile ergänzt um expliziten Code+Level+volumesPurged-Triple-Pin |

R2-Reviewer-Note: docs-check grün; weitere geprüfte Edge-Cases
ohne Befund: `Changed: nil` vs. `[]` (jsonenvelope.go:139-143
normalisiert), `--no-interactive --json` ohne `--purge`
(irrelevanter Pfad), Recorder-vs-Real-FS-Trigger-Vertrag,
WarnOnly-Statuskopplung. Confirmer-Pattern und `--purge`-
Dimension bleiben die substanziellsten Eigenleistungs-Anteile;
weitere Runden könnten Implementation-Reality (T3-T5) prüfen.

## Review-Round-3 (Pre-`next/`)

Implementation-Reality + Cross-Plan-Drift gegen den
R2-konsolidierten Stub (`e921522`). Sechs Findings (2 HIGH,
3 MEDIUM, 1 LOW), alle adressiert im selben Commit. Die zwei
HIGH-Befunde sind echte API-Realitäts-Lücken die R1/R2 textuell
übersehen haben.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | `catalogueFor(svc).Volumes != nil` ist erfundene API — `serviceCatalogueEntry` (`addservice_execute.go:190-224`) hat KEIN `Volumes`-Feld. Echte Felder: `volumeOptional bool`, `volumeRefLiteral string`. Plan-T6-Pin mit "Mock-Catalog `Volumes: nil`" nicht implementierbar | T0-(g) auf `catalogueFor(svc).volumeOptional == false` umgestellt; T6-Pin nutzt keycloak/otel als realistische Volumeless-Fixtures (heute existierende Catalog-Entries mit `volumeOptional: true`); kein Mock nötig |
| F2 | HIGH | Confirmer-Helper-Triple-Drift: T0-(j)/(h) sprachen `defensiveNoopConfirmer`, T3 sprach `noopConfirmer`, T1 plante "noopConfirmer-Helper analog noopProgress bauen". Realität: `noopConfirmer` existiert seit M4 in `application/noop.go:17-33` und tut genau was T0-(j) braucht | Alle Plan-Stellen auf `noopConfirmer` vereinheitlicht; T1-Tranche entfällt komplett ("kein neuer Helper, nur Swap-Mechanismus request-time"); T0-(j) klargestellt: NEU ist nur der Swap-Mechanismus, der Helper ist M4-Erbe |
| F3 | MEDIUM | T0-(e) Switch-Order-Tabelle hatte `ErrConfirmerUnavailable` NACH `ErrConfirmationRequired` — ein Multi-`%w`-Wrap mit beiden Sentinels würde falsch auf `LH-FA-INIT-005` matchen | Tabelle umsortiert: Infrastruktur-Sentinels (`ErrRemoveFileSystem`, `ErrConfirmerUnavailable`) VOR den fachlichen; expliziter "Tabellen-Reihenfolge = Switch-Reihenfolge"-Hinweis + T6-Multi-`%w`-Pin |
| F4 | MEDIUM | T3-Cell sagte noch "~6 FS-Wrap-Stellen" — T0-(d) R2-F3 hatte schon auf 10 kalibriert | T3-Cell auf 10 Stellen nachgezogen; LOC 200→240; ErrConfirmerUnavailable-Wrap explizit ergänzt |
| F5 | MEDIUM | Carveout-Inventarisierungs-Pflicht für WARN-on-Success-Pfad fehlt — generate hatte das Half-Write-State-Vorbild korrekt in `carveouts.md`+`open/`-Stub eingetragen, remove macht das nicht | T8-Cell um Carveout-Eintrag-Pflicht + ggf. open/-Trigger-Slice-Stub ergänzt (Pattern-Vorbild `slice-v2-generate-devcontainer-rollback-aware-write`) |
| F6 | LOW | Pre-UC-Sentinel-Kanal-Liste in T5 unvollständig — `getwd`-Failure (`cli/remove.go:117-120`) fehlte, init's Pattern (`init.go:221`) zeigt es explizit | T5-Cell-Pre-UC-Liste um `getwd`-Failure-Pfad ergänzt |

R3-Reviewer-Note: docs-check grün; Implementation-Reality-Pass
deckte zwei HIGH-Findings auf die R1+R2 nicht erwischt haben —
beide entstanden weil die Reviewer in R1/R2 Sub-Decisions
textuell konsolidiert haben ohne Code-Lookup. Geprüfte Code-
Realitäten ohne Befund: `ErrConfirmerUnavailable`-Wrap-Pfad an
einer Stelle sauber etablierbar, `EnabledUnset`-Pfad-Reihenfolge
konsistent mit `executeRemove`-Sequenz, AK Idempotenz-vs-
EnabledUnset-Trennung lesbar, T0-(j) Selbstkonsistenz Recon vs
AK, LOC-Bilanz T6 defensibel. Confirmer-Pattern nach F2-Fix
deutlich kleiner (nur Swap-Mechanismus statt neuer Helper).

## Review-Round-4 (Pre-`next/`)

Deep-Implementation-Reality + adversariale Coverage der R3-Fixes
gegen den R3-konsolidierten Stub (`35d6d51`). Sechs Findings (2
HIGH, 2 MEDIUM, 2 LOW). Die zwei HIGH-Befunde decken eine
Klassifikations-Fehler im R3-Inventar (Z. 307/330 sind nicht-FS)
und einen komplett neuen `delete`-Action-Vertrag auf (remove
ist der erste Use-Case mit dieser Action).

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | T0-(d) 10-Stellen-Inventar enthielt zwei Nicht-FS-Wrap-Stellen fälschlich. Z. 307 (default-Branch managedblock-Scanner) und Z. 330 (yaml.PatchScalar) sind fachliche Klassen, KEINE FS-I/O. Migration auf `ErrRemoveFileSystem` würde Datenkonsistenz-Defekte als Disk-Failure (Exit 14, retry-safe) klassifizieren — semantisch falsch | Inventar auf **8 FS-Wrap-Stellen** rekalibriert; Z. 307 + Z. 330 separat als fachlich-Klasse mit `ErrServiceInconsistent`-Wrap (analog Z. 304) festgezurrt |
| F2 | MEDIUM | T3-Cell verspricht `Remove()`-Wrapper mit FS-Swap, lässt aber offen wo `runPurgeGate` relativ zu `detectServiceState`/`fsSwap`/early-returns landet — drei plausible Varianten ohne Pin | T0-(c) um **Control-Flow-Skeleton** erweitert (Phasen-Reihenfolge explizit: confirmerSwap → Lock → fsSwap → detect → early-returns → conditional gate → execute → unswap → captures-mapping); T6-Pins für DryRun-skip + DetectInCaptureFS |
| F3 | MEDIUM | `--purge --yes --json` PLUS Mid-Write-Failure-Variante nicht im Plan — drei plausible Envelope-Formen (Error-only, Doppel-Diagnostic, Special-Code) | T0-(j) erweitert um **Variante A** festgezurrt: Error-Diagnostic dominiert, WARN unterdrückt, Zero-Response (`data: {service}` ohne `volumesPurged`); T6-Pin `TestRemove_PurgeYesJSON_MidWriteFailure_ErrorOnly` |
| F4 | HIGH | `delete`-Action-Vertrag fehlt — remove ist der erste Use-Case mit `PlannedFile.Action == "delete"` (für `RemoveAll`-Captures auf extraFiles). Init/add/generate produzieren nur `create`/`modify`; Diff-Renderer-Behavior für `delete` nicht im Plan | Neue **T0-(p)** Sub-Decision: `delete` = Old-only-Hunk (full-file-Remove); `cli-json-output.md` §6.6+§7 dokumentieren `action: "delete"`; T6-Pin `TestRemove_OtelExtraFileDelete_DiffHasDeleteHunk` |
| F5 | LOW | T6-Multi-`%w`-Switch-Order-Pin (R3-F3-Fix) ist Defense-only — heute kein Code-Pfad chained beide Sentinels; Pin-Rahmung "versehentlich" überzeichnet User-Value | T0-(e) als "Defense-only-Pin" qualifiziert mit Hinweis auf Cluster-T_close Mapper-Registry-Slice |
| F6 | LOW | T5-Cell sagte "Pre-UC-Sentinel-Kanal **ergänzt** werden" — Codepfade existieren bereits in `cli/remove.go:108-120`, NEU ist nur die Kanalisierung via `reportError`. Plus `getwd`-Wrap fällt in Default-Branch ohne typed Sentinel — Plan dokumentierte das nicht | T5-Cell-Formulierung präzisiert; `getwd`-Wrap-Default-Pfad explizit dokumentiert; Mapper-Tabelle bleibt unverändert |

R4-Reviewer-Note: docs-check grün. Implementation-Reality-Pass
deckte zwei HIGHs auf — eine Klassifikations-Fehler im R3-
Inventar (textuelle Konsolidierung verwechselte Scanner/Codec-
Wraps mit FS-Wraps) und einen kompletten neuen Action-Vertrag
(`delete` ist remove-spezifisch). Confirmer-Pattern und
`--purge`-Dimension sind nach R1-R4 vollständig durchspezifiziert.
Weitere Runden würden vermutlich nur noch Cosmetic-Drift fangen.

## Review-Round-5 (Pre-`next/`)

Spec-Coverage-Audit + Cross-Plan-Konsistenz gegen den R4-
konsolidierten Stub (`c6cef92`). User-getriebene Runde wegen
2-HIGH-Pattern in R3+R4. Vier Findings (1 HIGH, 2 MEDIUM, 1 LOW),
alle adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | T0-(c) Control-Flow-Skeleton hatte `confirmerSwap` AUSSERHALB der `removeMu`-Lock-Region — Race-Bedingung mit parallel laufenden Goroutinen die ihren eigenen Swap durchführen. Pattern-Erbe init (`initproject.go:328-348`) hat ALLE Swaps INNERHALB Lock | Skeleton reorganisiert: `LOCK → confirmerSwap → fsSwap → … → captures.Drain() → fsUnswap → confirmerUnswap → UNLOCK`. Race-Sicherheits-Block hinzugefügt; T6-Pin `TestRemove_ConcurrentInvocationsSerializeSwaps` für zwei parallele Goroutinen |
| F2 | MEDIUM | `LH-FA-ADD-007` Multi-Use als ERROR-Code (ErrServiceUnregistered) UND WARN-Code (Volumes-deferred) — Konsumenten können nicht über Code allein disambiguieren | T0-(g) Klarstellung: Code identifiziert Spec-Anforderung (Umbrella §924-947), nicht Sub-Semantik. Konsumenten müssen `(code, level)`-Tupel betrachten + `data.volumesPurged`-Status. Pattern konsistent mit init/down's `LH-FA-INIT-005`-Multi-Use für `ErrConfirmationRequired` |
| F3 | MEDIUM | `ErrServiceInconsistent` Triple-Use für Z. 304/307/330 — Z. 330 YAML-Patch-Defect ist semantisch Sentinel-Stretching (kein managed-block) | T0-(d) Triple-Use-Klarstellung mit zwei Alternativen (Konsolidierung vs. Sub-Sentinel `ErrYAMLPatchFailed`); Plan-Vorschlag Konsolidierung (Exit-Code 10 + Repair-Hint geteilt) mit T6-Pin auf `LH-FA-ADD-005` für alle drei Pfade |
| F4 | LOW | T0-(p) `delete`-Action-Vertrag spezifiziert RemoveAll Mid-Stream-Failure-Capture nicht — was wenn snapshot vor RemoveAll failt? | T0-(p) erweitert: Zwei-Phasen-Capture-Semantik (snapshot ReadFile → underlying.RemoveAll); Phase-1-Failure → `OldContent: nil`, Phase-2-Failure → vollständiges OldContent; Diff-Renderer behandelt nil als binary/unreadable |

R5-Reviewer-Note: docs-check grün. Spec-Coverage-Audit
bestätigt: LH-FA-CLI-005A (Confirmer-Gate), LH-FA-ADD-007 (Service
entfernen Umbrella §924-947), LH-FA-INIT-005 (Pattern-Erbe für
ErrConfirmationRequired), LH-FA-ADD-001/002/005, LH-NFA-REL-003
sind alle in der Spec verankert und korrekt gemappt. F1 ist der
substanziellste R5-Befund — Race-Bedingung wäre erst in
Concurrent-Production-Tests aufgefallen. Plan-Konsistenz nach 5
Runden weiterhin tragfähig.

## Review-Round-6 (Pre-`next/`)

Concurrency-Deep-Dive + T6-Test-Strategy-Vollständigkeit gegen
den R5-konsolidierten Stub (`eb20830`). User-getriebene Runde
nach dem R5-F1-Race-Befund. Vier Findings (1 HIGH, 2 MEDIUM,
1 LOW):

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | T0-(c) Skeleton verwendete `recorder.Drain()` — Real-API ist `recorder.Captured()` (`recordingfs.go:105`). Auch R5-F1-Adressierungs-Zelle hatte den Drift doppelt. Pattern-Erbe init (`initproject.go:358`) ruft `recorder.Captured()`. T3-Implementer würde Compile-Error sehen | Beide Stellen auf `recorder.Captured()` umgestellt (replace_all) |
| F2 | MEDIUM | T0-(c) Skeleton ließ Recorder-Lebensdauer-Invariante offen — Plan-Vertrag dokumentierte nicht, dass `recorder` lokale Variable (call-scoped) ist. Service-Feld-Interpretation würde bei parallelen Goroutinen Captures leaken (recordingfs.Captured() NICHT thread-safe) | Skeleton expanded: `fs, recorder := s.selectFS(req.PreviewMode)` als lokale Variable explizit; neue Sub-Block "Recorder-Lebensdauer-Invariante" mit Pattern-Erbe-Verweis; T6-Pin-Erweiterung um disjunkte-Capture-Sets-Assertion |
| F3 | MEDIUM | `LH-FA-ADD-007` ERROR-Pfad-Pin-Asymmetrie: AK-WARN-Migration hat expliziten `(code, level, volumesPurged)`-Triple-Pin, aber ERROR-Pfad (ErrServiceUnregistered) hat keinen symmetrischen `(code, level)`-Pin — laxer Implementer könnte WARN-Drift im ERROR-Pfad nicht erkennen | Neue AK-Zeile für ErrServiceUnregistered ERROR-Pfad-Pin: `code == "LH-FA-ADD-007"` AND `level == "error"` AND `exitCode == 10`; T6-Pin `TestRemove_ServiceUnregisteredJSON_ErrorLevelCodePin` |
| F4 | LOW | Pin-Namen-Inventar nur 8 explizit; T6 zählt 20-25 — restliche ~12 Pins (Idempotenz, Sentinels, Switch-Order-Defense) ohne kanonische Tags | T6-Cell um Pin-Namen-Mapping erweitert mit kanonischen Tags für die 9 explizit benannten Pins plus Hinweis dass die weiteren ~12 aus AK-Block + Sub-Decision-Pins direkt ableitbar sind |

R6-Reviewer-Note: docs-check grün. HIGH-Frequenz weiter
fallend: R1=3, R2=1, R3=2, R4=2, R5=1, **R6=1**. F1
ist ein 1-Char-Fix der zweimal im Plan vorkam — leicht zu
übersehen ohne API-Lookup, lehrreich für Pattern-Erbe-
Behauptungen. F2 ergänzt die Recorder-Scope-Invariante die R5-F1
implizit gelassen hatte. F3 schließt eine Pin-Asymmetrie zwischen
WARN- und ERROR-Pfad für denselben Code. F4 ist
Inventarisierungs-Hygiene. Plan-Konsistenz nach 6 Runden ist
stabil; weitere Runden bringen vermutlich nur LOW-Befunde.

## Review-Round-7 (Pre-`next/`)

Implementation-Path-Walk (T2→T3→T5) + Cross-Reference-Audit
gegen den R6-konsolidierten Stub (`82eb121`). Vier Findings
(1 HIGH, 2 MEDIUM, 1 LOW):

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | `removeEnvelopeData`-Struct-Form mehrdeutig zwischen init-Pattern (`nil` data) und generate-Pattern (typed Struct). Plus: `bool VolumesPurged`-omitempty würde `false` UND Key-Abwesenheit identisch droppen — Error-Pfad-Pin (`volumesPurged` ABWESEND) nicht enforceable | T0-(f) konkret festgezurrt: generate-Symmetrie + Pointer-Wrapping (`*PriorState`, `*State`, `*VolumesPurged`) für Key-Presence-vs-Absence-Disambiguierung; Pattern-Erbe `cliJSONEnvelope.DryRun`-`*bool`-Style |
| F2 | MEDIUM | WARN-Diagnostic-Emission-Ort offen: CLI-Layer kennt `--purge`-Flag aber nicht Catalog (Layer-Verletzung); Use-Case kennt Catalog aber Response trägt keine Warnings | T2 erweitert um **`RemoveServiceResponse.Warnings []DiagnosticEntry`-Feld** (Variante (b), saubere Layer-Trennung); LOC ~90→~110; CLI mapped via `mapWarningsToDiagnostics`-Helper |
| F3 | MEDIUM | T0-(a) Pattern-Erbe-Disziplin war nach 6 Runden noch TODO-Form; Verweise wie T0-(j) "Pattern-Erbe-Disziplin T0-(a) Spalte" zeigten ins Leere | T0-(a) als **zweispaltige Anchor-Tabelle** ausformuliert: 10 Erbe-1:1-Patterns + 10 Remove-spezifische Patterns mit Cross-Refs auf Sub-Decisions |
| F4 | LOW | Confirmer-Swap unconditional vs. konditional auf `SilenceConfirmer` nicht spezifiziert — init's Pattern (`initproject.go:345-349`) ist konditional | Skeleton auf `if req.SilenceConfirmer { confirmerSwap() }` umgestellt — Pattern-Erbe init-symmetrisch, weniger Code |

R7-Reviewer-Note: docs-check grün. HIGH-Frequenz weiter
konstant: R1=3, R2=1, R3=2, R4=2, R5=1, R6=1, **R7=1**. F1 ist
der substanziellste Befund — `bool`-vs-`*bool` Wahl beeinflusst
T6-Acceptance-Pin-Form (Key-Presence-Assertion vs Zero-Value-
Vergleich). F2 schließt eine Layer-Trennung sauber durch
Response-Type-Erweiterung. F3 löst eine 6-Runden-überfällige
TODO-Auflösung. F4 ist Pattern-Konsistenz mit init.

## Review-Round-8 (Pre-`next/`)

Stress-Test der R7-Festzurrungen + Test-Helper-Infrastruktur-
Audit gegen den R7-konsolidierten Stub (`a6e01e7`). Vier
Findings (1 HIGH, 2 MEDIUM, 1 LOW). HIGH-Frequenz weiterhin
konstant 1/Runde seit R5 — Konvergenz-Asymptote.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | Plan T0-(f) verspricht `jsontestutil.AssertFullEnvelope`-Key-Presence-Assertion für `volumesPurged`-Absence-Pin — Helper hat das nicht (`jsontestutil.go:22-27` dokumentiert nur vier Options). init/generate-Tests inspizieren `data` per manuellem Cast. Plan-Vertrag vs. Helper-Realität klafft auseinander | T6-Helper-Gap-Notiz mit zwei Implementer-Wahlen: (a) Helper-Erweiterung als T6-A-Sub-Tranche (+~50 LOC für `WithDataKeyAbsent`/`WithDataKeyPresent`) oder (b) Manuelle Cast-Form analog init/generate. Plan-Vorschlag Variante (a) als Pattern-Vorlauf für Folge-Slices |
| F2 | MEDIUM | T2 versprach `Response.Warnings []DiagnosticEntry` ohne Type-Definition — `DiagnosticEntry` existiert nicht; `domain.Diagnostic` hat Severity-Enum + ID-Field (semantischer Mismatch zur Wire-Form) | T2 erweitert um konkrete Type-Definition: neuer Port-Type `driving.WarningEntry struct { Code string; Level string; Message string }` analog `diagnosticItem`-Wire-Form; CLI mapped via `mapWarningsToDiagnostics`-Helper trivial. LOC ~110→~120 |
| F3 | MEDIUM | WARN-Emission-Ort im Use-Case-Body offen — drei plausible Stellen (vor Gate, nach Execute, in Execute) ohne Pin in T0-(c)-Skeleton | Skeleton ergänzt: `if req.Purge && catalogueFor(svc).volumeOptional == false { warnings = append(...) }` VOR `runPurgeGate`. Begründung: WARN auch im ErrConfirmationRequired-Pfad sichtbar (semantisch korrekt — User weiß, dass Purge eh deferred wäre); T5 unterdrückt WARN bei Error-Diagnostic (R4-F3-Variante-A). Zwei separate T6-Tests pinnen die Unterdrückung |
| F4 | LOW | `*bool` für VolumesPurged vs. `string`-omitempty für generate's Action — Pattern-Inkonsistenz ohne Begründung | T0-(f) Klarstellung: `VolumesPurged` MUSS `*bool` (false ist valider Success-Wert); `PriorState`/`State` Pointer für Symmetrie + Default-`""`-Drift-Schutz; generate's `string`-Pattern reicht weil Action-`""` klar Error-Marker ist |

R8-Reviewer-Note: docs-check grün. F1 ist Plan-vs-Test-Helper-
Realitäts-Lücke (spiegelbildlich zu R3-HIGH-F2 — dort Helper
existierte, Plan plante Neubau; hier Plan plante Helper-Use,
Helper-Surface fehlt). F2+F3 ist R7-F2-Cascade (Response.
Warnings-Festzurrung zog Folge-Fragen nach sich). F4 ist
Pattern-Konsistenz-Begründung.

**Konvergenz-Bewertung:** HIGH-Frequenz konstant 1/Runde über
vier Runden (R5-R8). Reviewer-Empfehlung: nach R8 in `next/`
migrieren — weitere Runden würden vermutlich Implementation-
Stress-Test-Befunde liefern die erst in T2/T3 auftauchen, oder
F1-ähnliche Helper-Gap-Variationen produzieren.

## Review-Round-9 (Pre-`next/`)

Cluster-Cross-Slice-Konsistenz + Carveouts/Roadmap-Impact gegen
den R8-konsolidierten Stub (`899af45`). Vier Findings (1 HIGH,
1 MEDIUM, 2 LOW). HIGH-Frequenz weiterhin konstant 1/Runde
**fünfte Runde in Folge** (R5-R9) — Asymptote bestätigt.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | T0-(p) verschmolz zwei Layer ("`delete` erster Use-Case"): Wire-Spec hat `delete` bereits als enum-Wert; Recorder hat `actionDelete` seit Rename/RemoveAll-Support. NEU ist nur die end-to-end-Sichtbarkeit. T8-Cell `§7 Mutations-Matrix` Wortlaut führt Closure-Implementer in die Irre (versucht neue action-Spalte zu bauen statt FS-Methoden-Zeile zu ergänzen) | T0-(p) Layer-Klarstellung: "erster end-to-end-sichtbare delete-Capture, der via mapCaptureToPlannedFiles in den Wire-Envelope wandert". T8-Cell präzisiert: §7-Zeile ergänzt `remove: WriteFile + RemoveAll`, KEINE neue action-Spalte; §6.6 worked example mit `data.changes[].action: "delete"` für otel-extraFile |
| F2 | MEDIUM | `driving.WarningEntry` ist erster Multi-Diagnostic-Port-Type im Cluster — Plan begründete die Layer-Wahl nicht als bewussten Cluster-Vorlauf. T2-Implementer könnte ihn als isolierte remove-Lösung lesen und z.B. `RemoveWarningEntry` umbenennen | T2-Cell um Cluster-Vorlauf-Disziplin ergänzt: Type bewusst generisch `driving.WarningEntry` für up/down recreate-Warnings (6/9) + config-set value-warnings (8/9); Erste-Slice-Pattern-Last analog `PreviewMode`-Rename in init T0-(c) |
| F3 | LOW | T8-Cell sagte generisch "Carveout-Eintrag" — Closure-Implementer könnte bestehenden Eintrag suchen. carveouts.md hat keinen Volume-Removal-Eintrag (verifiziert) | T8-Cell präzisiert: "**neuer Eintrag**" + Pattern-Vorbild generate-V2-Stub-Form + Spec-Anker `LH-FA-ADD-007 §"Volumes nur auf explizite Anforderung"` |
| F4 | LOW | T8-Cell sagte generisch "roadmap" — zwei konkrete Edits nötig (done-Zähler + Nächster-Schritt-Klausel) | T8-Cell präzisiert: zwei explizite Edits — done-Zähler 4→5 + remove aus Offen-Liste; Nächster-Schritt-Klausel auf Folge-Slice 6/9 up-down |

R9-Reviewer-Note: docs-check grün. F1 ist Cluster-Pattern-
Bruch-Klasse (Layer-Verschmelzung); F2 schließt Cluster-Vorlauf-
Lücke; F3+F4 sind Closure-Disziplin-Hygiene für T8.

**Konvergenz-Bewertung:** fünfte Runde 1-HIGH-Frequenz. Asymptote
über fünf Runden stabil. Reviewer-Empfehlung weiterhin: nach
R9-Adressierung in `next/` migrieren — weitere Runden würden
vermutlich nur ähnliche Closure-Präzisierungen + Cluster-Vorlauf-
Disziplin-Befunde produzieren.

## Review-Round-10 (Pre-`next/`)

First-Time-Implementer-Hand-off-Quality + Anti-Pattern-Detection
gegen den R9-konsolidierten Stub (`d0594b5`). Drei Findings
(1 HIGH, 1 MEDIUM, 1 LOW). HIGH-Frequenz weiterhin konstant
1/Runde **sechste Runde in Folge** (R5-R10) — Asymptote
endgültig bestätigt:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | T3-Tranchen-Cell sagte noch "10 FS-Wrap-Stellen" — T0-(d) wurde in R4-F1 auf 8 FS-Stellen kalibriert (Z. 307+330 sind fachlich, NICHT FS). Outdated-Reference seit R4. Tag-1-Implementer würde nach T3-Cell migrieren statt nach T0-(d), Z. 307+330 fälschlich mit ErrRemoveFileSystem wrappen — genau das HIGH-Befund-Szenario von R4-F1 | T3-Cell auf "8 FS-Wrap-Stellen" gezogen + expliziter Verweis: Z. 307+330 separat mit ErrServiceInconsistent (KEIN ErrRemoveFileSystem) |
| F2 | MEDIUM | Plan-Länge ~900 Zeilen — Review-Round-Tabellen R1-R10 sind ~25% des Plans; AK-Block 17 Bullets; T0-(c) Skeleton-Block ~80 Inline-Code-Zeilen. Pattern-Vergleich generate-done (~700) und init-done (~1583): remove sitzt zwischen, aber mit überproportionalem Review-History-Anteil | T8-Cell um Plan-Verdichtungs-Pflicht beim done/-Übergang erweitert: (i) Review-Round-Tabellen Adressierungs-Spalte auf einen Satz pro Finding kürzen, (ii) T0-(c) Skeleton in dedizierte H3-Sektion analog init-done, (iii) AK-Bullets gliedern |
| F3 | LOW | T0-(o) sagte "≥ 2 Runden, steht aktuell bei R1-R4" — outdated seit R5 (tatsächlich R9 mit 5-Runden-Asymptote). Folge-Slice-Implementer könnte aus T0-(o) lesen "≥ 2 reicht" und nach R2 migrieren | T0-(o) auf faktischen Stand: 10 Runden, Asymptote-Detektion bei ≥ 5 konstant-HIGH-Runden als Cluster-Pattern-Erbe für Folge-Slices 6/9-9/9 |

R10-Reviewer-Note: docs-check grün. F1 ist Outdated-Reference
aus R4 (Tag-1-Implementer-Showstopper, 1-Zahl-Edit). F2 ist
T8-Closure-Auftrag (non-blocking für `next/`-Übergang). F3 ist
Sub-Decision-Hygiene. **Geprüft ohne Befund**: `*string`-vs-
`string` (legitimer Cluster-Vorlauf), Multi-`%w`-Defense-Pin
(legitime Cluster-Konsistenz), Two-Phase-Capture (real-Recorder-
API-Reflexion), Bilanzen-Summation (konsistent), historische
Adressierungs-Zellen-Outdated-References (korrekt als Aufzeichnung).

**Konvergenz-Bewertung:** sechste Runde 1-HIGH-Frequenz; F1
ist ein Outdated-Reference-Treffer der seit R4 latent war,
kein neuer Substanz-Befund. Reviewer-Empfehlung **eindeutig**:
nach R10-Adressierung in `next/` migrieren — weitere Runden
würden vermutlich nur ähnliche Outdated-Reference- oder
Plan-Hygiene-Befunde produzieren.

## Review-Round-11 (Pre-`next/`) — Failure-Mode-Enumeration

Systematischer FMEA-Walk durch 20 Failure-Szenarien gegen den
R10-konsolidierten Stub (`37013bb`). Drei Findings (1 HIGH,
1 MEDIUM, 1 LOW). **Siebte Runde in Folge** mit konstanter
HIGH-Frequenz 1 (R5-R11) — Asymptote endgültig bestätigt.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | Failure-Mode #2 (`u-boot remove --json` ohne positional arg) hatte KEINEN Plan-Vertrag — `cobra.ExactArgs(1)`-Guard feuert VOR RunE, JSON-Helper läuft nie, Konsument bekommt Exit 2 ohne Envelope. Symmetrie-Bruch zu `remove "bad name" --json` (voller Envelope). Spec-LH-NFA-USE-004-§1841-Verletzung | AK-Pin `Cobra-Args-Missing-Pfad` ergänzt: T5-Pflicht Custom-`Args`-Validator (oder PreRunE) der `ErrServiceNameMissing`-Sentinel emittiert; via Pre-UC-Kanal in Envelope mit `LH-FA-CLI-006`/Exit 2. T0-(e)-Tabelle erweitert. T6-Pin `TestRemove_NoPositionalArg_JSON_EmitsCLI006Envelope` |
| F2 | MEDIUM | Failure-Modes #16 (fsFactory-NPE) und #17 (Context-Cancellation) hatten keinen Carveout-Verweis. Init's done-T0-(p) carved Context-Cancellation explizit Out-of-Scope; remove-Plan erwähnte das nicht. Tag-1-Implementer könnte plausibel ctx.Err()-Check einbauen | Out-of-Scope-Block erweitert: **Context-Cancellation-Carveout** (Pattern-Erbe init's done-T0-(p), Cluster-T_close-Block für Exit-130-Convention); **fsFactory-NPE-Schutz** als Composition-Root-Defekt-Klasse markiert |
| F3 | LOW | T6-Pin-Inventar war nur 9 von 20 Failure-Modes explizit benannt. R6-F4 etablierte Pin-Namen-Mapping aber FMEA-Coverage-Summary fehlte | T6-Cell um **Failure-Mode-Coverage-Summary** erweitert: nach R11 sind 18/20 explizit gepinnt + 2/20 als bewusste Out-of-Scope-Carveouts |

R11-Reviewer-Note: docs-check grün. F1 ist Coverage-Lücke die
R1-R10 textuell übersehen haben (Cobra-Default-Pfad vor RunE
nicht systematisch enumeriert). F2 ist Carveout-Erbe-Hygiene
für Pattern-Vorlauf. F3 ist FMEA-Inventarisierungs-Klarheit.

**Konvergenz-Bewertung:** R11-Total = 3 Findings (niedrigste
im Review-Verlauf). FMEA-Score 18/20 expliziter Plan-Pin + 2/20
bewusste Carveouts = **20/20 Coverage**. Asymptote bestätigt
über sieben Runden. Reviewer-Empfehlung wie R8/R9/R10: nach
R11-Adressierung in `next/` migrieren.

## Review-Round-12 (`next/`)

R11-Fix-Validation + Implementation-Pre-Flight-Walk gegen den
Post-Lifecycle-Stand (`0419747`). Vier Findings (1 HIGH, 2
MEDIUM, 1 LOW). Achte Runde in Folge mit konstanter HIGH-
Frequenz 1 (R5-R12) — Asymptote sehr stabil.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| F1 | HIGH | `ErrServiceNameMissing` (R11-F1) wurde im `driving`-Package geplant — sollte aber CLI-Layer analog `cli.ErrConflictingModeFlags` (`cli/cli.go:177`). Form-Validierungs-Sentinel vom CLI-Adapter emittiert, UC sieht ihn nie. `isUsageError`-Klassifikator-Block-Layout würde brechen | T0-(e)-Tabelle: `cli.ErrServiceNameMissing` mit Begründung; T2-Cell stellt klar dass T2 nur zwei port/driving-Sentinels (FS + ConfirmerUnavailable) ergänzt; AK-Pin nennt CLI-Layer-Heim |
| F2 | MEDIUM | T5-Pflicht offen zwischen Custom-`Args`-Validator und PreRunE — drei plausible Mechanismen, plus `cobra.ExactArgs(1)`-Wechsel ungenannt | T5-Pflicht festgezurrt: Custom-`Args`-Validator mit Konstruktor-Closure-Capture analog `newRemoveCommand`-Form; Pflicht-Begleit-Edit: `cobra.ExactArgs(1)` durch `validateRemoveArgs(a)` ersetzen. PreRunE-Alternative verworfen (Layer-Mismatch) |
| F3 | MEDIUM | `confirmerSwap()`-Skeleton-Pseudo-Code zeigte nur Aufrufstelle, nicht Mechanismus — drei plausible Varianten (Service-Field-Mutation, lokale Variable, Wrapper-Func mit Signature-Change) | T0-(c) ergänzt um Go-Code-Block für die Festzurrung: Service-Field-Mutation mit defer-Restore analog init's `s.progress`-Swap (`initproject.go:345-349`); lokale-Variable-Variante verworfen weil runPurgeGate-Refactor nötig |
| F4 | LOW | `driving.WarningEntry` ohne `Subject`-Feld — Cluster-Vorlauf-Gap für up/down Multi-Service-WARN | T2-Cell um `Subject string \`json:",omitempty"\``-Feld erweitert; remove nutzt es nicht (`""`-omitempty), aber up/down 6/9 + config-set 8/9 erben es ohne breaking Type-Change |

R12-Reviewer-Note: docs-check grün. F1 ist Layer-Idiom-
Konsistenz (Tag-1-Implementer-Showstopper für Cluster-Pattern-
Erbe). F2 schließt R11-F1-Mechanismus-Lücke. F3 ist Skeleton-
Mechanik-Präzision. F4 ist proaktive Cluster-Vorlauf-Disziplin.

**Konvergenz-Bewertung:** achte Runde mit 1 HIGH-Frequenz seit
R5. Reviewer-Empfehlung: nach R12 implementations-bereit;
weitere Pre-Implementation-Runden brächten nur Hygiene-
Variations.

## Review-Round-13 (Pre-T8, Code-against-Plan)

Adversarialer Pre-T8-Audit des Post-T6-Code-States (alle
Acceptance-Tests grün, T2-T6 Hashes `d0c9c5d`/`dbbf7b1`/
`3b079dd`/`3188e75`/`9eae9ec`) gegen Plan + Spec. Sechs Findings
(1 HIGH, 1 MED, 4 LOW). Fokus: Args-Validator-Flag-Awareness,
Envelope-Symmetrie-Bruch, Plan-Drift, Hygiene.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| HIGH-1 | HIGH | `validateRemoveArgs` reichte hardcoded `dryRun=false, diff=false` an `writeErrorEnvelope` durch — `u-boot --dry-run --json remove` ohne arg produzierte Minimal-Envelope trotz Spec-§1842-Voll-Schema-Pflicht | T7-Code-Fix in `4fb3fea`: Validator liest `cmd.Flags().GetBool("dry-run"/"diff")` zur Args-Validator-Zeit (Cobra hat Local-Flags zu diesem Zeitpunkt geparst) und reicht User-Flag-State weiter. Pin: `TestRemove_NoPositionalArg_DryRunJSON_EmitsFullSchemaEnvelope`. Plan-T0-(r) (T8) dokumentiert die Defense. |
| MED-1 | MED | `len(args)>1` Pfad: `cobra.ExactArgs(1)` fing `--json remove a b c` mit Roh-Error auf stderr ab, OHNE Envelope auf stdout — Symmetrie-Bruch zum Missing-Arg-Pfad (Spec §1841) | T7-Code-Fix in `4fb3fea`: Validator detektiert `len(args)>1` vor `cobra.ExactArgs(1)`-Delegate, ruft `writeErrorEnvelope` und reicht den Cobra-Error durch. Exit 2 via `isUsageError`-`"accepts "`-Prefix. Pin: `TestRemove_TooManyArgs_JSON_EmitsCLI006Envelope`. Plan-T0-(r) (T8) dokumentiert. |
| LOW-1 | LOW | Plan-Drift: drei T6-Cell-Pin-Namen aus Sub-Decision T0-(c) fehlen in Acceptance-Tests (`TestRemove_ConcurrentInvocationsSerializeSwaps`, `TestRemove_DryRun_DetectStateUsesCaptureFS`, `TestRemove_DryRunPurgeYes_NoConfirmerCall`) — Pin-Namen leben im Application-Layer (`removeservice_test.go`) als bewusste Layer-Carveouts | T8-Plan-Drift-Doku unten unter §R13-LOW-1 Plan-Drift-Adressierung markiert die drei Pin-Namen als bewusste Application-Layer-Heim-Carveouts mit Code-Anker. |
| LOW-2 | LOW | Kommentar-Präzisierung in `printRemoveSummary`-Block (Z. 505+) — Pre-T7-Kommentar suggerierte WARN auf allen Pfaden, real ist `--purge && !VolumesPurged`-Mode-gefiltert | T7-Code-Fix: Kommentar in `cli/remove.go:493-525` präzisiert (analog R14-HIGH-2-Block) plus Dry-Run-Suppression-Begründung. |
| LOW-3 | LOW | T0-(g) WARN-Migration-Doku ↔ T0-(h)(b) Volume-Removal-Status-Doku — leichte Doppelung an drei Stellen über die Sub-Decisions verteilt | T8-Plan-Verdichtung-Punkt: die beiden Sub-Decisions sind orthogonal (g = wie WARN emittiert, h = wann WARN gilt), aber die Cross-Refs lassen sich konsolidieren. **In T8 nicht durchgeführt — bewusst aufgeschoben** (Folge-Plan-Hygiene-Slice oder bei Cluster-T_close). |
| LOW-4 | LOW | `mapRemoveErrorToDiagnostic`-Header-Kommentar referenziert Spec-Anker für `ErrServiceUnregistered`, aber nicht für `ErrServiceInconsistent` (asymmetrisch) | T8-Plan-Verdichtung-Punkt: low-prio Cosmetic, **nicht durchgeführt** (analog LOW-3). |

R13-Reviewer-Note: vier T7-Fixes erforderlich (HIGH-1, MED-1, R14-HIGH-2, R14-MED-1). Sechs Findings, davon vier echte Code-Defekte und zwei Plan/Doku-Hygiene-Punkte (LOW-3, LOW-4) bewusst aufgeschoben.

## Review-Round-14 (Pre-T8, Defense-Audit)

Parallel-Runde zu R13 mit anderem Fokus: WARN-Pollution im
Human-Mode, Path-Leak-Defense, Doku-AK-Completion. Fünf Findings
(2 HIGH, 1 MED, 2 LOW).

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| HIGH-1 | HIGH | Doku-AKs CHANGELOG `### Added` + `cli-json-output.md` §6-Tabelle + §6.6 + §7-Mutations-Matrix fehlen — Plan-T8-Schritt | **T8-Scope**: T8-Closure-Schritt 1+2 (siehe §Tranche-Status). NICHT Pre-T7 weil Doku ohnehin T8-Schritte sind. |
| HIGH-2 | HIGH | `printRemoveSummary` zeigte `--purge`-deferred-Volumes-WARNING auch in `PreviewDryRun` — Use-Case skippt Gate (T0-(h)(a)), führt keine Mutation aus, WARN-Prosa wäre "ist-deferred" statt "würde-deferred" semantisch falsch | T7-Code-Fix in `4fb3fea`: `printRemoveSummary`-Signature um `previewMode driving.PreviewMode` erweitert; WARN-Block überspringt bei `previewMode == PreviewDryRun`. `PreviewAndApply` behält WARN. Pin: `TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution` (beide Richtungen). |
| MED-1 | MED | `mapRemoveErrorToDiagnostic` reichte `err.Error()` 1:1 an `diagnostic.message` — FS-Wraps der Form `fmt.Errorf("remove write %s: %w: %w", absPath, ErrRemoveFileSystem, raw)` tunneln den absoluten Filesystem-Pfad in den User-facing Output (Info-Leak im JSON-Mode maschinen-lesbar) | T7-Code-Fix in `4fb3fea`: neuer `baseDirSanitizedError`-Wrapper-Type in `cli/remove.go:465-491`; `runRemove` wrappt UC-Error mit `sanitizeBaseDir(removeErr, cwd)`. `errors.Is`/`As` bleiben intakt via Unwrap-Chain. Pin: `TestRemove_FSErrorWithAbsolutePath_SanitizesMessage`. Plan-T0-(q) (T8) dokumentiert. |
| LOW-1 | LOW | T0-(c) Skeleton-Block enthält ~80 Zeilen Inline-Pseudocode + Go-Code-Blocks — bewusste Cluster-Vorlauf-Disziplin, aber Plan-Länge skaliert | T8-Plan-Verdichtung-Punkt analog R13-LOW-3/4: **nicht durchgeführt**, separat als Plan-Verdichtungs-Slice oder Cluster-T_close. |
| LOW-2 | LOW | Sub-Decisions-Tabelle T0-(a) "Pattern-Erbe-Disziplin" referenziert generate-T1-Vorzug — minimaler Drift, nicht-blocking | T8-Plan-Verdichtung-Punkt: **nicht durchgeführt**. |

R14-Reviewer-Note: HIGH-2 ist semantischer Bug (Output-Klassifikation falsch), MED-1 ist Security-Hygiene (Info-Leak-Defense). Beide T7-fixierbar. HIGH-1 ist T8-Scope per Definition.

## Review-Round-15 (Post-T7, Bestätigung)

Adversariale Bestätigungs-Runde gegen den Post-T7-Code-State
(`4fb3fea`) mit sechs Verifikations-Angles (A-F: Validator-Flag-
Read-Timing, Sanitizer-Wrapper-Edge-Cases, WARN-Konsistenz,
Pattern-Drift-Audit add/init/generate, Coverage-Pfade,
Plan-T0-(q)/T0-(r) Sub-Decisions). **Outcome: HIGH=0, MED=0,
LOW=2 + 1 Cross-Slice** — vier T7-Fixes adversarial bestätigt,
kein T7-Followup nötig, T8-Closure-ready.

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| LOW-1 | LOW | `replaceBareBaseDir` (Pre-T8) nutzte `strings.ReplaceAll(msg, baseDir, ".")` ohne Word-Boundary — naive Form, würde `<baseDir>-cache/lock` zu `.-cache/lock` mangeln. Praktisch heute unwahrscheinlich (FS-Layer referenziert immer exakten cwd-Pfad), aber konstruierbar | T8-Code-Fix im T8-Closure-Commit: neuer `replaceBareBaseDir`-Helper + `isPathComponentByte` mit Word-Boundary-Check (gefolgt von End-of-String oder Nicht-Pfad-Byte). Pin: `TestRemove_FSErrorWithBaseDirSubstring_NotMangled` konstruiert Sibling-Pfad `<baseDir>-cache/lock` und pinnt unverändert. |
| LOW-2 | LOW | Kein expliziter Pin für `--json --dry-run remove a b c` (Voll-Schema bei `len(args)>1 + --dry-run` Kombi-Pfad) — Code-Pfad ist korrekt aber Future-Regression-Risiko bei Validator-Refactoring | T8-Code-Fix: neuer Pin `TestRemove_TooManyArgs_DryRunJSON_EmitsFullSchemaEnvelope` ergänzt explizite Coverage. |
| Cross-Slice-1 | Cross-Slice | `add.go:79`, `generate.go:78` haben rohes `cobra.ExactArgs(1)` ohne JSON-Envelope-Hook (strukturell gleiches R13-MED-1-Problem). Init/add/generate-Mapper haben keinen BaseDir-Sanitizer für `diagnostic.message` (strukturell gleiches R14-MED-1-Problem) | **NICHT T8-Scope**. T8 dokumentiert den Pattern-Drift als bewussten Out-of-Scope und legt neuen open/-Slice-Stub [`slice-v1-cli-json-envelope-consolidation`](../open/slice-v1-cli-json-envelope-consolidation.md) plus carveouts.md-Eintrag an. Trigger: Cluster-Stand 7/9 oder 8/9. |

R15-Reviewer-Note: Wrap-Site-Inventur-Korrektur in der Commit-Message von T7 — nur 1 von 8 FS-Wrap-Stellen trägt absoluten Pfad direkt (`removeservice.go:410`), nicht alle 8 wie Commit-Message suggerierte. Defense bleibt nötig wegen raw-FS-Error-Component. Plan-T0-(q) ist um die Korrektur ergänzt.

**Konvergenz-Bewertung:** R13-R14 ergaben 11 Findings (3 HIGH, 3 MED, 5 LOW), davon 4 echte Code-Defekte in T7 gefixt. R15 bestätigt HIGH=0 mit zwei zusätzlichen LOW-Findings (in T8 mitgefixt) plus einer Cross-Slice-Pattern-Drift-Erkenntnis (zu Folge-Slice ausgelagert). HIGH-Frequenz-Asymptote über R5-R15 (11 Runden): 1-1-1-1-1-1-1-1 für R5-R12, 1-2-0 für R13-R14-R15. Konvergenz erreicht.

## R13-LOW-1 Plan-Drift-Adressierung

Drei T6-Cell-Pin-Namen aus T0-(c) Sub-Decision sind als bewusste **Application-Layer-Heim-Carveouts** zu verstehen, nicht als CLI-Acceptance-Pins:

| Pin-Name | Code-Anker | Layer-Heim | Begründung |
| --- | --- | --- | --- |
| `TestRemove_ConcurrentInvocationsSerializeSwaps` | `application/removeservice_test.go` (Konzept) | Application | Race-Test gegen `removeMu`-Lock und Recorder-call-scope-Invariante. Acceptance-Test wäre nicht informativ — CLI sieht den seriellen Output, der Race wäre nur in Use-Case-Internals beobachtbar (Service-Field-Mutation timing). |
| `TestRemove_DryRun_DetectStateUsesCaptureFS` | `application/removeservice_test.go` (Konzept) | Application | Spy-Read-Counter im `RecordingFileSystem`-Stub würde im CLI-Acceptance-Layer nicht durchgereicht. Service-Layer-Test. |
| `TestRemove_DryRunPurgeYes_NoConfirmerCall` | `application/removeservice_test.go` (Konzept) | Application | Confirmer-Call-Counter im Stub würde im CLI-Acceptance-Layer nicht durchgereicht. Service-Layer-Test. |

**Plan-Vertrag**: die drei Pin-Namen bleiben als Sub-Decision-Pins in T0-(c) — Tag-1-Implementer-Anker. Code-Heim ist die Application-Layer-Test-Suite. **NICHT** als Folge-Slice ausgelagert (zu kleinteilig); T6-Sub-Tranche im Application-Layer würde sie konsolidieren, falls sich in Cluster-T_close-Audit Bedarf zeigt.

## Out of Scope

- **Context-Cancellation-Handling** (R11-MED-F2-Carveout, Pattern-
  Erbe init's done-T0-(p)): `ctx.Err() == context.Canceled` mid-
  Remove() bleibt Status-quo — fällt auf Default-Branch
  `LH-FA-CLI-006` / Exit 2. Eine konsistente Exit-130-Convention
  für ALLE modifying-Subcommands ist Cross-Cutting-Concern,
  eigener Cluster-T_close-Block. Remove ändert den heutigen Pfad
  NICHT — `cli-json-output.md` §6.6 nimmt den Doku-Hint von
  init's §6.4 als Vorbild auf.
- **`fsFactory`-NPE-Schutz** (R11-MED-F2-Carveout): wenn ein
  Composition-Root-Bug eine `nil`-FS aus `s.selectFS(mode)`
  liefert, NPE beim ersten `s.fs.*`-Call. Heute kein Test-Pin
  (init/add/generate auch nicht); Pattern bleibt "Composition-
  Root-Bug = Defekt, kein User-Pfad". Out-of-Scope-Verbleib
  konsistent mit Cluster-Pattern.
- **Volume-Auto-Removal**: heute `--purge` deferred mit WARNING.
  Auto-Removal bleibt eigener Slice (LH-FA-ADD-007 §"Volumes
  nur auf explizite Anforderung" implementiert den Gate, aber
  nicht den Removal-Aufruf an `docker volume rm`).
- **HTTP- oder gRPC-Schnittstellen**: ADR-0010 schließt
  explizit aus.
- **Schema-Versionierung** (`schemaVersion: 1`): siehe
  Cluster-Slice §Out of Scope.
- **Add-on-spezifische Cleanup-Hooks** (Keycloak realm-export,
  OTel collector-config): nicht remove-Slice-Scope.
- **Generisches `mapErrorToDiagnostic`-Registry**: Cluster-T_close-
  Aufgabe (Altitude-Reviewer-Vorschlag aus add R6 #I1).

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  §T0-Outcomes — Vorgaben für den Folge-Slice-Block.
- Pattern-Vorbilder:
  [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
  — Service-Operation-Symmetrie (add ↔ remove);
  [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
  — Pattern-Erbe (PreviewMode, RecordingFileSystem, Diff-Renderer,
  Helper, Multi-`%w`-Switch-Order, Confirmer-Silencing analog
  ProgressPort-Silencing T0-(o));
  [`slice-v1-cli-json-dry-run-generate`](../done/slice-v1-cli-json-dry-run-generate.md)
  — Data-Carrier-Form (T0-(p) bereits etabliert).
- Spec: `LH-FA-CLI-007/008`, `LH-NFA-USE-004`, `LH-FA-ADD-007`,
  `LH-FA-CLI-005A` §254, `LH-NFA-REL-003`
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- Code-Anker heute:
  [`removeservice.go`](../../../../internal/hexagon/application/removeservice.go)
  (~370 LOC, Plan-Phase + Execute-Phase + `--purge`-Gate),
  [`cli/remove.go`](../../../../internal/adapter/driving/cli/remove.go)
  (~170 LOC, RunE-Erweiterungs-Ziel + WARNING-Migration in
  `printRemoveSummary`),
  [`port/driving/removeservice.go`](../../../../internal/hexagon/port/driving/removeservice.go)
  (Carrier-Types + ein remove-spezifischer Sentinel
  `ErrServiceUnregistered`; `ErrRemoveFileSystem` muss in T2
  ergänzt werden).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).

## Tranche-Status (2026-06-07, T8-Closure)

| Tranche | Status | Commit-Hash | Notiz |
| --- | --- | --- | --- |
| T0 (Plan, R1-R12, in-progress) | ✅ | `d7f9e65` | 56 Plan-Findings adressiert |
| T1 (entfällt) | ✅ | — | `noopConfirmer` lebt bereits in `application/noop.go:17-33` (M4 Confirmer-Port-Slice). R3-HIGH-F2-Fix |
| T2 (Port-Types) | ✅ | `d0c9c5d` | PreviewMode + SilenceConfirmer + Warnings + WarningEntry-Type + ErrRemoveFileSystem + ErrConfirmerUnavailable |
| T3 (Application-Layer) | ✅ | `dbbf7b1` | fsFactory + removeMu + Remove()-Wrapper + 8 FS-Wrap-Stellen Multi-`%w` + volumesPurgedWarnings |
| T4 (Composition-Root + Helper) | ✅ | `3b079dd` | `newPreviewFSFactory`-Helper (Altitude-Reviewer add R6 #I3 Trigger bei 4 Factories met) |
| T5 (CLI-RunE-Rewrite) | ✅ | `3188e75` | validateRemoveArgs + removeEnvelopeData + runRemove + writeRemoveJSON + mapWarningsToDiagnostics + mapRemoveErrorToDiagnostic + Allowlist |
| T6 + T6-A (Acceptance-Tests) | ✅ | `9eae9ec` | 21 Pin-Tests + WithDataKeyAbsent/Present-Helper + delete-Hunk-Vertrag T0-(p) |
| T7 (Pre-T8-Review-Fixes) | ✅ | `4fb3fea` | R13-HIGH-1 + R13-MED-1 + R14-HIGH-2 + R14-MED-1 plus 4 Pin-Tests |
| T8 (Closure) | ✅ | dieser Commit | CHANGELOG + cli-json-output.md §6/§6.6/§7 + roadmap done-Zähler 4→5 + zwei carveouts.md-Einträge + zwei open/-Stubs (volume-auto-removal + envelope-consolidation) + R15-Bestätigung + R15-LOW-1 Sanitizer-Substring-Robustheit-Code-Fix (`replaceBareBaseDir` Word-Boundary) + R15-LOW-2 Coverage-Pin (`TestRemove_TooManyArgs_DryRunJSON_EmitsFullSchemaEnvelope`) + R15-LOW-1 Substring-Pin (`TestRemove_FSErrorWithBaseDirSubstring_NotMangled`) + T0-(q)/(r) Sub-Decisions + R13/R14/R15-Tabellen + R13-LOW-1 Plan-Drift-Doku + done/-Move |

**Coverage**: 91.10% (von 90.50% nach T5 → stabil durch T6/T7/T8 — zwei neue Pin-Tests in T8, Sanitizer-Helper-Erweiterung verändert die statement-coverage nicht).

**Review-Findings-Konsolidierung**:
- R1-R12 (Plan-Phase, Pre-Implementation): 56 Findings adressiert.
- R13 (Pre-T8 Code-against-Plan): 6 Findings (1 HIGH, 1 MED, 4 LOW). HIGH-1 + MED-1 → T7-Code-Fixes (`4fb3fea`); LOW-1 → T8-Plan-Drift-Doku; LOW-2 → T7-Code-Fix (Kommentar); LOW-3/4 → bewusst nicht durchgeführt (Cluster-T_close-Hygiene).
- R14 (Pre-T8 Defense-Audit): 5 Findings (2 HIGH, 1 MED, 2 LOW). HIGH-2 + MED-1 → T7-Code-Fixes (`4fb3fea`); HIGH-1 → T8-Doku-Closure; LOW-1/2 → bewusst nicht durchgeführt.
- R15 (Post-T7 Bestätigung): HIGH=0, MED=0, LOW=2 + 1 Cross-Slice. LOW-1 + LOW-2 → T8-Code-Fixes (Sanitizer-Word-Boundary + TooMany-DryRun-Pin). Cross-Slice-1 → neuer open/-Stub [`slice-v1-cli-json-envelope-consolidation`](../open/slice-v1-cli-json-envelope-consolidation.md) + carveouts.md-Eintrag.

**Plan-Verdichtungs-Carveout**: die in der T8-Cell von §Tranchen erwähnte Plan-Verdichtung (R1-R10-Adressierungs-Spalten kürzen, T0-(c) Skeleton-Block in eigene H3-Sektion, AK-Block-Gliederung) ist in T8 **NICHT** durchgeführt. Die R1-R12-Review-Tabellen bleiben in ihrer originalen Wortform als historisches Audit-Trail-Artefakt. Falls bei Cluster-T_close eine systematische Plan-Verdichtungs-Sub-Tranche kommt, wandert die Disziplin dort hin — nicht jeder Folge-Slice braucht seine eigene Verdichtung.
