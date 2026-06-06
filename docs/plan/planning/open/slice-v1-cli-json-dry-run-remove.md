# Slice V1: `remove --json` / `--dry-run` / `--diff` — Add-Inverse mit Purge-Gate

> **Status:** T0-Discovery + R1/R2 adressiert, `open/`. Fünfter Folge-Slice (5/9) des
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
  `InitProjectService.fsFactory`.
- **T0-(d)** `ErrRemoveFileSystem`-Sentinel-Einführung +
  Wrap-Audit (R2-MED-F3-Kalibrierung): heute Single-`%w` an
  **10 FS-Stellen** in `removeservice.go` (NICHT ~6 wie initialer
  Stub):
  - **Write/Remove**: Z. 235 (WriteFile), Z. 241 (RemoveAll
    extraFiles) — 2 Stellen.
  - **Read/Exists/Stat**: Z. 272 (Exists compose/env), Z. 282
    (ReadFile compose/env), Z. 286 (Lstat compose/env), Z. 321
    (ReadFile u-boot.yaml), Z. 325 (Lstat u-boot.yaml), Z. 358
    (Exists extraFiles) — 6 Stellen.
  - **YAML-Codec**: Z. 330 (yaml.PatchScalar) — 1 Stelle.
  - **Managedblock-Scanner**: Z. 307 (scan compose/env for block)
    — 1 Stelle.
  - **AUSGESCHLOSSEN**: Z. 304 (managedblock-malformed wrappt
    `ErrServiceInconsistent`, KEIN FS-Wrap) — bleibt fachlich-
    Klasse.

  T3 migriert alle 10 auf Multi-`%w` analog init's
  `initproject.go:925/967/1015/1117/1143`-Pattern. T6 ergänzt
  Read-Pfad-FS-Failure-Pin (mindestens einer) damit die Switch-
  Order-Garantie nicht löchrig wird.
- **T0-(e)** **Switch-Order-Pflicht** im neuen
  `mapRemoveErrorToDiagnostic`. Diagnostic-Code-Tabelle (T6-Pin-
  Pflicht pro Zeile):

  | Sentinel | LH-Code | Exit |
  | --- | --- | --- |
  | `ErrRemoveFileSystem` | `LH-NFA-REL-003` | 14 |
  | `ErrServiceUnsupported` | `LH-FA-ADD-002` | 10 |
  | `ErrServiceUnregistered` | `LH-FA-ADD-007` | 10 |
  | `ErrServiceInconsistent` | `LH-FA-ADD-005` | 10 |
  | `ErrProjectNotInitialized` | `LH-FA-ADD-001` | 10 |
  | `ErrConfirmationRequired` | `LH-FA-INIT-005` | 10 |
  | `ErrConfirmerUnavailable` (NEU, T2; wrappt heutigen string-Error in `removeservice.go:171`) | `LH-FA-CLI-005A` | 10 |
  | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 |
  | `ErrConflictingModeFlags` (`--yes ⊕ --no-interactive`) | `LH-FA-CLI-005A` | 2 |
  | Default (unknown) | `LH-FA-CLI-006` | 1 |

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

  **Volume-Presence-Pflicht** (R2-MED-F5-Fix): WARN-Diagnostic
  wird NUR emittiert wenn der Catalog-Entry tatsächlich Volumes
  deklariert (`catalogueFor(svc).Volumes != nil`, analog
  `planExtraFileDeletes` Z. 349). Bei einem zukünftigen Catalog-
  Entry ohne Volumes (z. B. Config-only Add-on) wäre eine WARN
  semantisch falsch ("Purge deferred" obwohl es nichts zu purgen
  gibt). T6-Pin: `TestRemove_PurgeOnVolumelessService_NoWarn` —
  Mock-Catalog mit `Volumes: nil`, `--purge --json` →
  `diagnostics: []`, kein WARN-Eintrag, `data.volumesPurged:
  false`. Pattern-Vorlauf für Keycloak/OTel-Catalog-Erweiterungen.
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
  (a) `--purge --diff --json` ohne `--yes` → defensiveNoopConfirmer
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
- **T0-(j)** **Confirmer-Swap-Pattern (NEU, R1-HIGH-2-Fix)**:
  remove etabliert das Pattern; es ist NICHT geerbt von init
  (init swappt nur ProgressPort, nicht Confirmer). Form:
  `req.SilenceConfirmer = flags.JSON`. Bei `--purge --json` ohne
  `--yes`: ConfirmerPort wird auf einen `defensiveNoopConfirmer`
  umgeswapt der `false, nil` returnt — `runPurgeGate`
  (removeservice.go:173-176) wandelt das in
  `ErrConfirmationRequired`. **Semantik-Klarstellung**: das ist
  KEIN Silencing (keine UX-Information-Verlust-Symmetrie zu
  noopProgress), sondern eine **bewusste Behaviour-Change** im
  JSON-Mode — User muss explizit `--yes` setzen um im
  JSON-Mode zu purgen. Pattern-Erbe-Disziplin T0-(a) Spalte
  führt das als remove-spezifisch.

  **`--purge --yes --json`-Pfad** (R1-MED-5-Fix): bei
  `req.Yes==true` skipped runPurgeGate (Z. 162-164) ohne
  Confirmer-Call → Execute läuft durch → `VolumesPurged: false`
  (v0.3.0 deferred). Plan-Vertrag: trotzdem WARN-Diagnostic
  emittiert (T0-(g)), aber `exitCode: 0` UND `status: warn`
  (Spec §447-Kopplung) — Warn-only verschiebt nicht den
  Exit-Code. T6-Pin: `TestRemove_PurgeYesJSON_WarnOnly` mit
  `status: warn`, exit 0, `data.volumesPurged: false`.
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
  + Adversarial).

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | ------ | --------------- | --- |
| T0 | Discovery + Sub-Decisions (a)-(o) klären; Review-Runden | — (Plan) | — |
| T1 | Refactor-Tranche (wenn überhaupt nötig — generate-Pattern ist etabliert; möglicher T1-Scope: noopConfirmer-Helper im application package, analog noopProgress) | ~30-60 | T0 |
| T2 | Port-Types: `RemoveServiceRequest.PreviewMode` + `SilenceConfirmer`-Feld, `RemoveServiceResponse.PlannedFiles`/`Changes`-Felder, **zwei neue Sentinels**: `ErrRemoveFileSystem` (FS-Klasse, T0-(d)) UND `ErrConfirmerUnavailable` (Confirmer-I/O-Error-Klasse, R2-HIGH-F1-Fix für T0-(e)-Tabelle). | ~90 | T0 |
| T3 | Application-Layer: `RemoveServiceService.fsFactory` + `removeMu sync.Mutex` + `NewRemoveServiceServiceWithFactory` + `Remove()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; Multi-`%w`-Wrap an den ~6 FS-Wrap-Stellen; Confirmer-Swap auf noopConfirmer im JSON-Mode. | ~200 | T2 |
| T4 | Composition-Root-Wiring `removeFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runRemove` ruft generische Helper mit `command="remove"`, `mapErr=mapRemoveErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; `mapRemoveErrorToDiagnostic` neu; `data`-Struct (`removeEnvelopeData`); WARNING-Migration in `diagnostics[]` (`level: "warn"`); **Pre-UC-Sentinel-Kanal** für `domain.ErrInvalidServiceName` und `ErrConflictingModeFlags`: müssen via `reportError`-Helper emittiert werden, NICHT durch Cobra-Default-Print (R1-LOW-7-Fix), damit der JSON-stdout-Cleanliness-Pin aus init T0-(o) hält. **Human-Mode-Diff-Renderer** (R2-LOW-F6-Fix): bei `--purge --diff` ohne `--json` bleibt die deferred-Volumes-Prosa auf `errOut`, NICHT im Diff-Body. T6-Pin: `TestRemove_PurgeHumanDiff_StderrSeparation` mit getrennten Buffer-Assertions. | ~250 | T1 + T2 |
| T6 | Acceptance-Tests: ~20-25 Tests (drei JSON-Modi + NoOp Single+Repeat + Mid-Write-Failure + ConfirmationRequired-Pfade × 3 Varianten + Service-Sentinels × 4 + WARNING-Migration-Pin + `--purge`-on/off × Dry-Run-Kombos (T0-(h)) + `--purge --yes --json` WarnOnly-Pin (T0-(j) R1-MED-5) + `ErrConflictingModeFlags`-Pin). R1-MED-6-Kalibrierung: ~600-700 LOC realistisch (Confirmer-Pattern-Neumuster zieht Test-Surface). | ~650 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe). | ~80 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.6/§7, roadmap, slice nach done/ mit DoD-Hash-Tabelle. | — (Doku) | T7 |

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
