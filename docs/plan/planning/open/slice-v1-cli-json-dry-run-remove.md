# Slice V1: `remove --json` / `--dry-run` / `--diff` — Add-Inverse mit Purge-Gate

> **Status:** T0-Discovery + R1/R2/R3/R4/R5 adressiert, `open/`. Fünfter Folge-Slice (5/9) des
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

- **T0-(a)** Pattern-Erbe-Disziplin (was 1:1 von init/generate,
  was remove-spezifisch).
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
    confirmerSwap(req.SilenceConfirmer)     # NEUER Swap-Mechanismus T0-(j), INNERHALB Lock
    fsSwap(req.PreviewMode != PreviewNone)   # init's Swap-Pattern
    state := detectServiceState(s.fs, s.yaml, ...)
    if state == Unregistered:    return early (ErrServiceUnregistered)
    if state == InconsistentYAML: return early (ErrServiceInconsistent)
    if state == Deactivated:     return no-op (KEIN runPurgeGate)
    # Active / EnabledUnset / InconsistentBlock:
    if req.PreviewMode != PreviewDryRun:     # T0-(h)(a) Skip-Logik
        runPurgeGate(req)
    executeRemove(...)
    captures := recorder.Drain()             # vor Unswaps drainen
    fsUnswap; confirmerUnswap; UNLOCK         # Defer-Restore-Pattern analog init
    resp.PlannedFiles = mapCaptureToPlannedFiles(captures, req.BaseDir)
    return response
  ```

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
- **T0-(o)** Pre-`next/`-Review-Runden-Erwartung: ≥ 2 (Discovery
  + Adversarial). Steht aktuell bei R1-R4.
- **T0-(p)** **`delete`-Action-Vertrag (NEU, R4-HIGH-F4 +
  R5-LOW-F4-Erweiterung)**:
  remove ist der **erste** Use-Case, der `PlannedFile.Action ==
  "delete"` produziert — `RemoveAll` (Z. 241) für extraFiles wird
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

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | ------ | --------------- | --- |
| T0 | Discovery + Sub-Decisions (a)-(o) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (R3-HIGH-F2-Fix): `noopConfirmer` existiert bereits seit M4 Confirmer-Port-Slice in `application/noop.go:17-33` und tut exakt was T0-(j) braucht (`ConfirmRemoveVolumes → false, nil`). `RemoveServiceService`-Konstruktor (`removeservice.go:48`) nutzt ihn schon als nil-Fallback. T3 swappt den existierenden Helper request-time, kein neuer Helper nötig. | — (entfällt) | T0 |
| T2 | Port-Types: `RemoveServiceRequest.PreviewMode` + `SilenceConfirmer`-Feld, `RemoveServiceResponse.PlannedFiles`/`Changes`-Felder, **zwei neue Sentinels**: `ErrRemoveFileSystem` (FS-Klasse, T0-(d)) UND `ErrConfirmerUnavailable` (Confirmer-I/O-Error-Klasse, R2-HIGH-F1-Fix für T0-(e)-Tabelle). | ~90 | T0 |
| T3 | Application-Layer: `RemoveServiceService.fsFactory` + `removeMu sync.Mutex` + `NewRemoveServiceServiceWithFactory` + `Remove()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; **Multi-`%w`-Wrap an den 10 FS-Wrap-Stellen** (R3-MED-F4-Kalibrierung, T0-(d) Inventar); `ErrConfirmerUnavailable`-Sentinel-Wrap in `runPurgeGate` Z. 171; Confirmer-Swap auf existierenden `noopConfirmer` im JSON-Mode. | ~240 | T2 |
| T4 | Composition-Root-Wiring `removeFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runRemove` ruft generische Helper mit `command="remove"`, `mapErr=mapRemoveErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; `mapRemoveErrorToDiagnostic` neu; `data`-Struct (`removeEnvelopeData`); WARNING-Migration in `diagnostics[]` (`level: "warn"`); **Pre-UC-Sentinel-Kanal** (R4-LOW-F6-Klarstellung: Codepfade existieren bereits in `cli/remove.go:108-120`, NEU ist nur die Kanalisierung via `reportError` analog `init.go:205, 216, 221`) für `domain.ErrInvalidServiceName`, `ErrConflictingModeFlags` UND `getwd`-Failure (`fmt.Errorf("determine working directory: %w", err)`, R3-LOW-F6-Fix). Der `getwd`-Wrap trägt KEIN typed Sentinel und fällt in den Default-Branch `LH-FA-CLI-006` / Exit 1 (Pattern-Erbe von init T0-(o)); Mapper-Tabelle T0-(e) NICHT ergänzt. **Human-Mode-Diff-Renderer** (R2-LOW-F6-Fix): bei `--purge --diff` ohne `--json` bleibt die deferred-Volumes-Prosa auf `errOut`, NICHT im Diff-Body. T6-Pin: `TestRemove_PurgeHumanDiff_StderrSeparation` mit getrennten Buffer-Assertions. | ~250 | T1 + T2 |
| T6 | Acceptance-Tests: ~20-25 Tests (drei JSON-Modi + NoOp Single+Repeat + Mid-Write-Failure + ConfirmationRequired-Pfade × 3 Varianten + Service-Sentinels × 4 + WARNING-Migration-Pin + `--purge`-on/off × Dry-Run-Kombos (T0-(h)) + `--purge --yes --json` WarnOnly-Pin (T0-(j) R1-MED-5) + `ErrConflictingModeFlags`-Pin). R1-MED-6-Kalibrierung: ~600-700 LOC realistisch (Confirmer-Pattern-Neumuster zieht Test-Surface). | ~650 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe). | ~80 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.6/§7, roadmap, slice nach done/ mit DoD-Hash-Tabelle; **Carveout-Eintrag in `carveouts.md`** für deferred-Volume-Auto-Removal + WARN-on-Success-Semantik mit Trigger auf einen Volume-Auto-Removal-Folge-Slice (R3-MED-F5-Fix, Pattern-Vorbild `slice-v2-generate-devcontainer-rollback-aware-write`); ggf. **`open/`-Plan-Stub für den Trigger-Slice** anlegen analog generate's V2-Rollback-Stub. | — (Doku) | T7 |

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

## Out of Scope

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
