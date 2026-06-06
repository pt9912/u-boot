# Slice V1: `remove --json` / `--dry-run` / `--diff` — Add-Inverse mit Purge-Gate

> **Status:** T0-Discovery, `open/`. Fünfter Folge-Slice (5/9) des
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
**Confirmer-Port-Kollision mit `--json`-Mode**: der Confirmer
prompt würde stdin/stdout polluten — analog zu init's
ProgressPort-Silencing T0-(o) braucht remove eine
Silence-Variante (Sub-Decision T0-(o)).

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
  "deactivated"`, `status: ok`, Exit 0.
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
- ✅ **stderr-WARNING-Migration**: die heutige `printRemoveSummary`-
  WARNING auf errOut bei `--purge`-Status FALSE muss in den
  Envelope wandern (`diagnostics[]` mit `level: "warn"`-Eintrag
  oder `data.warnings[]`-Array — Sub-Decision T0-(g)). Im
  JSON-Mode darf stderr nicht durch die WARNING-Prosa polluten.
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
  Wrap-Audit: heute Single-`%w` an den ~6 FS-Stellen
  (WriteFile/RemoveAll/Exists/ReadFile/Stat). T3 erweitert
  auf Multi-`%w` analog init/generate.
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
  | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 |
  | Default (unknown) | `LH-FA-CLI-006` | 1 |

- **T0-(f)** **Envelope-`data`-Form festgezurrt**: Success-Envelope
  trägt `data: {"service": "<…>", "priorState": "<…>", "state":
  "<…>", "volumesPurged": <bool>}`. Error-Envelope trägt nur
  `data: {"service": "<…>"}` ohne `priorState`/`state`/
  `volumesPurged` (Zero-Response auf Error-Pfad — analog
  generate T0-(q)).
- **T0-(g)** **WARNING-Migration**: heutige `printRemoveSummary`-
  stderr-WARNING (Z. 163-171) bei `--purge && !VolumesPurged` muss
  im JSON-Mode in den Envelope. Drei Optionen:
  (a) `diagnostics[]` mit `level: "warn"` und eigenem LH-Code
  (z. B. `LH-FA-ADD-007-VOLUMES-DEFERRED`).
  (b) `data.warnings[]`-Array als Free-Form-Liste.
  (c) `data.volumesPurged: false` allein + Konsument leitet
  WARNING ab.
  Vorschlag: (a) — passt zum Diagnostics-Schema-Vertrag (Spec
  §1834 erlaubt `warn`-Level), pinnbar via `jsontestutil.
  AssertFullEnvelope`.
- **T0-(h)** **`--purge`-in-Dry-Run-Verhalten**: Dry-Run impliziert
  Null-Mutationen. Sollte der Confirmer-Gate trotzdem laufen?
  Vorschlag: **nein** — analog init's `initGit`-Skip im Dry-Run
  (T0-(n)). Implementierung: `if req.PreviewMode == PreviewDryRun
  { skip confirmer gate, set VolumesPurged: false }`.
- **T0-(i)** **`--purge`-Mutex mit `--dry-run`/`--diff`?** Analog
  init's `--template`-Mutex (T0-(i))? Vorschlag: **NEIN** — Purge
  ist eine Side-Effect-Dimension, kein Renderer-Pfad. Dry-Run +
  Purge ist semantisch konsistent: "zeige was Remove + Purge
  ändern WÜRDE", auch wenn der Gate-Run skipped wird.
- **T0-(j)** **Confirmer-Silencing-Form**: analog init T0-(o):
  `req.SilenceConfirmer = flags.JSON`. Bei `--purge --json` ohne
  `--yes`: ConfirmerPort wird auf einen `noopConfirmer` umgeswapt
  der **defensiv refused** (gibt `false` zurück, was zu
  `ErrConfirmationRequired` führt). User muss explizit `--yes`
  setzen um im JSON-Mode zu purgen.
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
| T2 | Port-Types: `RemoveServiceRequest.PreviewMode` + `SilenceConfirmer`-Feld, `RemoveServiceResponse.PlannedFiles`/`Changes`-Felder, neuer `ErrRemoveFileSystem`-Sentinel. | ~70 | T0 |
| T3 | Application-Layer: `RemoveServiceService.fsFactory` + `removeMu sync.Mutex` + `NewRemoveServiceServiceWithFactory` + `Remove()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; Multi-`%w`-Wrap an den ~6 FS-Wrap-Stellen; Confirmer-Swap auf noopConfirmer im JSON-Mode. | ~200 | T2 |
| T4 | Composition-Root-Wiring `removeFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runRemove` ruft generische Helper mit `command="remove"`, `mapErr=mapRemoveErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; `mapRemoveErrorToDiagnostic` neu; `data`-Struct (`removeEnvelopeData`); WARNING-Migration in `diagnostics[]` (`level: "warn"`). | ~220 | T1 + T2 |
| T6 | Acceptance-Tests: ~15-18 Tests (drei JSON-Modi + NoOp Single+Repeat + Mid-Write-Failure + ConfirmationRequired-Pfade + Service-Sentinels + WARNING-Migration-Pin). | ~500 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe). | ~80 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.6/§7, roadmap, slice nach done/ mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: ~1100-1200 — schmaler als init (~1480), in
der Größenordnung von generate (~1150). Pattern-Erbe von init/
generate spart die Infrastruktur; remove-spezifisch sind nur die
Confirmer-Silencing- und WARNING-Migration-Sub-Decisions.

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
