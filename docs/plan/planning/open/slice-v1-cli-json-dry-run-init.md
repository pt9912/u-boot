# Slice V1: `init --json` / `--dry-run` / `--diff` — modifying-Surface erbt von Add

> **Status:** geplant für v0.4.0 — dritter Folge-Slice des Cluster-
> Slice [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 3/9). Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
> 1:1 für die Carrier-Types, den `RecordingFileSystem`-driven-
> Adapter, den Pure-Go Diff-Renderer, das `previewModeFromFlags`-
> Mapping und die Error-Envelope-Pipeline; init-spezifisch sind
> die Mutations-Matrix (MkdirAll + WriteFile direkt plus
> CopyExclusive/Mkdir/MkdirAll/Copy/RemoveAll indirekt via
> `BackupPath`), drei zusätzliche Sentinels
> (`LH-FA-INIT-001`..`-006`), die Soft-Existing-Detection
> (`LH-FA-INIT-004`) und der Template-Modus
> (`init --template <name>`, Catalog-Read im Dry-Run).
>
> Init ist der erste Folge-Slice, der vom Add-Pattern erbt — die
> Erbschafts-Disziplin (was 1:1, was init-spezifisch) ist Sub-
> Decision in §T0-Discovery. Slice liegt in `open/`.

## Auslöser

Cluster-Slice `slice-v1-cli-json-dry-run` §T0-Outcomes (a)+(b)+(e)
machen jeden modifying-Subcommand für `--json`/`--dry-run`/`--diff`
verbindlich (`LH-NFA-USE-004` §1813, `LH-FA-CLI-007` §326,
`LH-FA-CLI-008` §451-489). `u-boot init` ist nach `add` der zweite
modifying-Subcommand und der **wichtigste Onboarding-Use-Case**
(Cluster-§T0-Discovery Z. 320 nennt `init --dry-run --diff --json`
explizit als Beispiel-Hauptanwendung): ein neuer Nutzer will sehen,
was `u-boot init <project>` an Files/Dirs anlegt, bevor er das auf
einer existierenden Codebase ausführt.

Spec-Bezug (geerbt von add-Slice):

- `LH-FA-CLI-007` (Dry-Run, Voll-Schema §326)
- `LH-FA-CLI-008` (Diff, §451-489)
- `LH-NFA-USE-004` (Minimalkontrakt §1841)

Init-spezifische Sentinels und Spec-Stellen:

- `LH-FA-INIT-001`..`-007` (Projekt-Skeleton, Verzeichnisstruktur,
  Soft-Existing-Detection, Backup-Pfad-Failures, Service-Name-
  Validation)
- `LH-NFA-REL-003` (FS-Failure-Klasse, geerbt für Mid-Write-Failure
  analog `add`)

Heute-Stand-Pre-Scan
(`internal/hexagon/application/initproject.go`, 1079 LOC):

| Phase | Methode | Pfade (typisch) | Indirekt? |
| --- | --- | --- | --- |
| Skeleton-Dirs (Z. 776) | `MkdirAll` | `.devcontainer/`, `docker/`, `scripts/`, `docs/`, evtl. `.github/workflows/` | direkt |
| Skeleton-Files (Z. 818, 866, 968) | `WriteFile` | `u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`, `.env.example`, `.gitignore`, `Makefile`, devcontainer-/template-Files | direkt |
| Backup (Z. 978, [`backup.go`](../../../../internal/hexagon/application/backup.go)) | `WriteFileExclusive`, `CopyExclusive`, `Mkdir`, `MkdirAll`, `Copy`, `RemoveAll` | `<file>.bak.<n>` plus Backup-Verzeichnis | indirekt via `BackupPath` |

Damit nutzt init **alle 8** `driven.FileSystem`-Mutations-Methoden
in der Praxis — der `RecordingFileSystem` aus add-T1-B deckt das
bereits ab (Drift-Schutz war Cluster-Pflicht). Kein neuer
driven-Adapter nötig.

## Aufhebungsbedingung

Acht Flag-Kombinationen für `u-boot init <project>` liefern spec-
konforme Outputs (geerbt von add, ein-zu-eins-Symmetrie):

```bash
u-boot init myproj                          # human, schreibt
u-boot init myproj --dry-run                # human Vorschau, kein Write
u-boot init myproj --diff                   # human Unified-Diff + Write
u-boot init myproj --dry-run --diff         # human Unified-Diff, kein Write
u-boot init myproj --json                   # Minimalkontrakt-Envelope, schreibt
u-boot init myproj --dry-run --json         # Voll-Schema-Envelope, kein Write
u-boot init myproj --diff --json            # Voll-Schema-Envelope, Hunks, Write
u-boot init myproj --dry-run --diff --json  # Voll-Schema, Hunks, kein Write
```

`make test` + `make lint` + `make docs-check` grün.
`jsonAllowlist` migriert: `"u-boot init": true`, Reject-Pin-Test
schrumpft (10 → 9 Reject-Cases).

Konkrete Pin-Form für `init myproj --dry-run --json` (frisch-CWD,
`basic`-Template):

```json
{
  "status": "ok",
  "command": "init",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    {"path": "myproj", "action": "create"},
    {"path": "myproj/.devcontainer", "action": "create"},
    {"path": "myproj/docker", "action": "create"},
    {"path": "myproj/scripts", "action": "create"},
    {"path": "myproj/docs", "action": "create"},
    {"path": "myproj/u-boot.yaml", "action": "create"},
    {"path": "myproj/compose.yaml", "action": "create"},
    {"path": "myproj/.env.example", "action": "create"},
    {"path": "myproj/README.md", "action": "create"},
    {"path": "myproj/CHANGELOG.md", "action": "create"},
    {"path": "myproj/.gitignore", "action": "create"}
  ],
  "changes": [/* identische 11 Einträge mit count = CountLines(NewContent) */],
  "diagnostics": [],
  "exitCode": 0
}
```

Exakte Reihenfolge und vollständige Pfad-Liste = Sub-Decision T0-(g)
(orientiert sich an heutigem `initproject.go`-Aufruf-Pattern;
Mutations-Matrix für jeden Pfad in T0-(c) finalisieren).

Negative-Pin: bei `--dry-run` null Production-FS-Mutationen, gleicher
Spy-Mechanismus wie in add T5.

Soft-Existing-Detection-Pin (`LH-FA-INIT-004`):
`u-boot init myproj --dry-run --json` auf eine **existierende**
Projekt-CWD ohne `--backup`/`--force` liefert:

```json
{
  "status": "error",
  "command": "init",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [],
  "changes": [],
  "diagnostics": [
    {"level": "error", "code": "LH-FA-INIT-004", "message": "project 'myproj' already exists; use --backup or --force"}
  ],
  "exitCode": 10
}
```

## Akzeptanzkriterien

- ✅ **Pattern-Erbe von add 1:1** (T0-(a)): `RecordingFileSystem`,
  `driving.PlannedFile/ChangeEntry/Hunk`, `diff.Compute`/`Render`/
  `CountAdditions`/`CountLines`/`CountBytesDiff`/`IsBinary`,
  `jsontestutil.AssertFullEnvelope` mit `checkHunks`, `cliJSONEnvelope`
  + `newMinimalEnvelope`/`newFullEnvelope`/`writeEnvelope` —
  alle ohne Änderung wiederverwendet.
- ✅ **`previewModeFromFlags` extrahiert** (T0-(b)): aus
  `internal/adapter/driving/cli/add.go` Z. 114 in ein neues
  `cli`-Paket-internes Helper-File (Vorschlag `previewmode.go`),
  damit init+remove+generate+config-set ihn ohne Copy-Paste
  konsumieren. Behält die T0-(b)-Wahrheitstabelle aus add (Dry-Run
  wins bei (yes,yes)).
- ✅ **`driving.PreviewMode` shared** (T0-(c)): umbenennen
  `driving.AddPreviewMode` → `driving.PreviewMode` (Konstanten
  `PreviewNone`/`PreviewDryRun`/`PreviewAndApply` analog).
  `driving.AddServiceRequest.PreviewMode` bleibt, neue
  `driving.InitProjectRequest.PreviewMode` ergänzt. `fsFactory`-
  Signatur bekommt `driving.PreviewMode` statt
  `driving.AddPreviewMode` — gemeinsame Closure-Form für alle 5
  modifying-Services.
- ✅ **`InitProjectService.fsFactory`** (T0-(d)): analog zu
  `AddServiceService.fsFactory` ergänzen plus
  `NewInitProjectServiceWithFactory`-Konstruktor. `Add()`-
  Wrapper-Pattern aus add 1:1 spiegeln (Mutex-serialisiert,
  s.fs-Swap, defer-Restore, `mapCaptureToPlannedFiles(captured,
  req.BaseDir)`).
- ✅ **CLI-RunE für `u-boot init`**: drei Flag-Pfade analog add,
  Error-Envelope-Gate via gemeinsamem Helper (T0-(e)).
- ✅ **`mapErrorToDiagnostic` für init** (T0-(f)): per-subcommand
  Switch analog add, init-spezifische LH-Codes:
  - `LH-FA-INIT-001`: `ErrInvalidProjectName` (domain)
  - `LH-FA-INIT-002`: `ErrProjectExists`/`ErrFileExists` (soft-
    existing)
  - `LH-FA-INIT-003`: noch nicht modelliert (Spec-Sub-Decision —
    siehe T0-(f))
  - `LH-FA-INIT-004`: `ErrConfirmationRequired`,
    `ErrForceRequiresBackup`, `ErrBackupUnsupportedKind`
  - `LH-FA-INIT-005`: `ErrBackupSuffixExhausted`,
    `ErrBackupSourceMissing`
  - `LH-FA-INIT-006`: `domain.ErrInvalidServiceName` (geteilt mit
    add, falls --template ein Service-Name-Validation auslöst)
  - `LH-NFA-REL-003`: neuer `driving.ErrInitFileSystem`-Sentinel
    für FS-Failures (analog `ErrAddFileSystem`)
- ✅ **Composition-Root-Wiring** (T0-(g)): `initFSFactory` analog
  `addFSFactory` in `cmd/uboot/main.go`; gleiches Closure-Pattern.
  Pflicht-Erweiterung: `initSvc` migriert auf
  `NewInitProjectServiceWithFactory`. App-Struct + `cli.New(...)`
  bleiben unverändert.
- ✅ **Mutations-Matrix `cli-json-output.md` §7**: init-Zeile
  ergänzt (alle 8 Methoden möglich via BackupPath; siehe Pre-Scan
  oben).
- ✅ **count-Semantik für BackupPath-CopyExclusive** (T0-(h)):
  Backup-Files sind content-identisch zu Original. T0-(h) finalisiert
  die `changes[].count`-Form: 0 (identische Bytes) oder
  `CountLines(NewContent)` (gleicher Inhalt, eigene Datei). Pin-Test
  pro Variante.
- ✅ **Template-Modus** (T0-(i)): `--template basic` lädt im
  Dry-Run-Pfad nur die Catalog-Reads (kein WriteFile), die Templates
  landen als geplante PlannedFiles aus dem Recorder. Sub-Decision:
  Catalog-Read-Failure-Pfad (`ErrTemplateNotFound`/`ErrTemplateRender`/
  `ErrTemplateCatalog` → bestehende LH-Codes, kein neuer Sentinel).
- ✅ **Test-Pflichten**: Acceptance-Tests für alle 8 Flag-
  Kombinationen + Soft-Existing-Pin + Backup-Pfad-Pin + Template-
  Modus-Pin + Mid-Write-Failure-Scenario + Null-FS-Mutationen-Spy
  (auf Recorder-Ebene wiederverwertet aus add T1-B).
- ✅ **`docs/user/cli-json-output.md`**: §6 Migrations-Tabelle
  init→done, §6.4 neue init-Sektion (analog §6.3 add), §7
  Mutations-Matrix-init-Zeile.
- ✅ **CHANGELOG**-Eintrag (Pattern aus add-Slice).

## T0-Discovery (vor `next/`-Übergang)

Sub-Decisions, die dieser Slice klären muss, bevor er in `next/`
wandert. Bewusst kondensiert vs. add-Slice — add hat 12 Sub-
Decisions, weil es Pattern-Vorbild war; init kann das meiste
referenzieren.

### T0-(a) Pattern-Erbe-Disziplin: was wird 1:1 übernommen?

Init darf nur dort von add abweichen, wo init-spezifische
Verhalten dokumentiert sind (Mutations-Matrix, Backup-Indirektion,
Soft-Existing-Detection, Template-Mode, init-LH-Codes). **Vorschlag
(T0-Festlegung)**: harte Erbe-Liste in T0-Outcomes; jeder add-
Helper, der nicht 1:1 wiederverwendet werden kann, braucht eine
Begründung im Outcome.

### T0-(b) `previewModeFromFlags` extrahieren

`add.go` hat heute `previewModeFromFlags(dryRun, diff)` als private
Funktion (slice-v1-cli-json-dry-run-add Z. 114). Init braucht
dieselbe Wahrheitstabelle. **Vorschlag (T0-Festlegung)**:
Extraktion in ein neues `cli`-Paket-internes File
(`internal/adapter/driving/cli/previewmode.go`) als Package-Helper.
Erste Refactor-Tranche des Slices (T1), bevor init-RunE darauf
zugreift.

### T0-(c) `driving.PreviewMode` umbenennen

Heute heißt der Enum `driving.AddPreviewMode` mit Add-Prefix —
historisch korrekt, weil er für add eingeführt wurde. Für die
4 folgenden modifying-Slices ist das schief: jeder bekommt seinen
eigenen Mode-Type oder importiert `AddPreviewMode` unter falscher
Bedeutung.

**Drei Optionen:**

1. **Umbenennen** zu `driving.PreviewMode` (Konstanten unverändert).
   `AddPreviewMode` als type-Alias erhalten für Backward-Compat —
   add-Code ändert sich nicht.
2. **Eigener Enum pro Service** (`InitPreviewMode`,
   `GeneratePreviewMode`, …). Drift-Risiko bei Konstanten-Werten.
3. **Lassen**, init importiert `AddPreviewMode` direkt.
   Semantisch-Drift (init-Code referenziert „Add"-Type).

**Vorschlag (T0-Festlegung)**: **Option 1** — Rename + type-Alias.
Slice-T1-Tranche macht den Rename plus `type AddPreviewMode =
PreviewMode`-Alias als Carveout, der nach Cluster-T_close entfernt
wird.

### T0-(d) `InitProjectService.fsFactory` Konstruktor-Form

Add hat einen zweiten Konstruktor (`NewAddServiceServiceWithFactory`)
und behält den Legacy-Konstruktor für Backward-Compat
(`addservice.go:203`). Tests, die den legacy-Konstruktor benutzen,
zeigen `PlannedFiles: nil` (Recorder nil).

**Vorschlag (T0-Festlegung)**: gleiche Form für init — neuer
`NewInitProjectServiceWithFactory(fsFactory, …)`-Konstruktor neben
dem heutigen `NewInitProjectService(fs, …)`. Legacy bleibt für
existierende Tests funktional.

### T0-(e) Error-Envelope-Helper gemeinsam machen?

`add.go` hat `reportAddError` und `writeAddErrorEnvelope` —
beide könnten zu `reportError(out, err, resp, flags, cmd)` und
`writeErrorEnvelope(out, err, resp, cmd, dryRun, diff)`
generalisiert werden. Init duplizieren wäre 2× Copy-Paste, viele
Wartungs-Stellen.

**Drei Optionen:**

1. **Pro Subcommand eigene Funktion** (Status quo bei add). Init
   bekommt `reportInitError`/`writeInitErrorEnvelope`. N×Duplikation.
2. **Gemeinsame `report{Error,Envelope}` mit `command`-Parameter**.
   Init-spezifisches Verhalten (Mid-Failure-Voll-Schema-Switch,
   wenn `len(resp.PlannedFiles) > 0`) lebt in der gemeinsamen
   Funktion; Caller liefern command-String.
3. **Helper-Struct `envelopeWriter{command, mapErr}`** mit Methode
   `report(out, err, resp, flags)`. Mehr Ceremony, aber sauberer
   bei zukünftigen Erweiterungen.

**Vorschlag (T0-Festlegung)**: **Option 2** — pragmatischer Mittel-
weg. Init-T1 macht den Refactor; add migriert in derselben Tranche
auf den gemeinsamen Helper. Acceptance-Tests aus add bleiben grün
(reine Refactor-Tranche).

### T0-(f) Diagnostic-Code-Quelle für init

Plan-Bezug: `LH-FA-INIT-001`..`-007`. Sub-Decision: gibt es einen
init-spezifischen FS-Failure-Sentinel (analog `ErrAddFileSystem` →
`LH-NFA-REL-003`)?

**Vorschlag (T0-Festlegung)**: ja — neuer
`driving.ErrInitFileSystem`-Sentinel in `port/driving/initproject.go`,
gemappt auf `LH-NFA-REL-003`/Exit-Code 14 in
`cli.isFilesystemError`. Wrap-Stellen in `initproject.go`
(`WriteFile`-Sites Z. 818/866/968 und backup-relevante Pfade)
ergänzen. Analog zum add-Pattern.

**Code-Map** (T0-Outcomes-Tabelle finalisiert):

| Sentinel | LH-Code | Exit-Code |
| --- | --- | --- |
| `domain.ErrInvalidProjectName` | `LH-FA-INIT-001` | 10 |
| `driving.ErrProjectExists`, `driving.ErrFileExists` | `LH-FA-INIT-004` | 10 |
| `driving.ErrConfirmationRequired`, `driving.ErrForceRequiresBackup`, `driving.ErrBackupUnsupportedKind` | `LH-FA-INIT-004` | 10 |
| `driving.ErrBackupSuffixExhausted`, `driving.ErrBackupSourceMissing` | `LH-FA-INIT-005` | 14 |
| `domain.ErrInvalidServiceName` (geteilt mit add) | `LH-FA-INIT-006` | 10 |
| **`driving.ErrInitFileSystem` (neu)** | **`LH-NFA-REL-003`** | **14** |

### T0-(g) `plannedFiles[]`-Reihenfolge + Catalog-Read-Phase

Init ruft die FS-Operationen in einer deterministischen Reihenfolge
(siehe `initproject.go`): erst `MkdirAll` für alle Skeleton-Dirs,
dann `WriteFile` für alle Skeleton-Files. Der Recorder capturet
in Aufrufreihenfolge.

**Sub-Decision**: `plannedFiles[]`-Reihenfolge im Voll-Schema
folgt der Aufruf-Reihenfolge des Use-Cases (deterministisch durch
Code-Pfad) — kein Re-Sort, kein Stable-Sort. Pinnt in Acceptance-
Tests die exakte Liste für den `basic`-Template-Default.

**Catalog-Read-Phase**: `--template basic` liest die Catalog-Files
über `templateCatalogAdapter` (externaltemplates). Diese **Reads**
landen NICHT im Recorder (Reads passieren am underlying-FS, nicht
am Recorder). Der Recorder sieht nur die ResultWrites. Das ist OK
für V1.

### T0-(h) `count`-Semantik bei BackupPath-Indirektion

Backup-Operationen kopieren Original → `.bak.N`. Inhalt ist
identisch. T0-(g) aus add-Slice sagt: `create` =
`CountLines(NewContent)`. Bei einem Backup-Copy ist NewContent =
Original-Body, also `count` = Lines des Original-Files.

**Vorschlag (T0-Festlegung)**: gleiche `CountAdditions`/`CountLines`/
`CountBytesDiff`-Form wie add — keine Backup-spezifische Sub-
Logik. Backup-File-Eintrag im `plannedFiles[]` hat
`action: "create"`, count = Lines(Original). User sieht:
„`.env.example.bak.1`, action create, count 4" — semantisch
korrekt.

### T0-(i) Template-Mode-Failure-Pfade

`--template <name>` mit unbekanntem Name → `ErrTemplateNotFound`
(`LH-FA-CLI-006` heute, bleibt). `--template`-Render-Failure →
`ErrTemplateRender` (`LH-NFA-REL-003`? oder `LH-FA-CLI-006`?).

**Vorschlag (T0-Festlegung)**: Template-Failures bleiben mit ihrer
heutigen LH-Klassifikation (keine Erweiterung in diesem Slice —
Out-of-Scope, weil Template-Catalog ein eigener Komplex ist).

### T0-(j) `init --json` (Minimalkontrakt) ohne `--dry-run`/`--diff`

Analog add T0-(k): Spec-streng Minimal — kein
`plannedFiles[]`/`changes[]`/`dryRun`/`diff` im Output. Hint im
Doku-Block: „use `--dry-run --json` to preview".

## T0-Outcomes

Verbindliche Festzurrung wandert nach `next/`-Übergang in diesen
Block — analog add-Slice. Erwartete Form: Tabelle mit den 9 Sub-
Decisions plus Implementations-Pflicht-Spalte (T1-T6).

## Tranchen (vorgeschlagen)

LOC-Schätzung deutlich schlanker als add (das alles neu erfunden
hat); init erbt das Pattern.

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | Discovery + Sub-Decisions klären; Pattern-Erbe-Liste pinnen | — (Plan) |
| T1 | **Refactor-Tranche**: `previewModeFromFlags` extrahieren (T0-(b)), `driving.PreviewMode` umbenennen + Alias (T0-(c)), `reportError`/`writeErrorEnvelope` generalisieren (T0-(e)). Add-Tests bleiben grün ohne Test-Edits. | ~100 |
| T2 | **Port-Types + Sentinel**: `driving.InitProjectRequest.PreviewMode`-Field, `driving.ErrInitFileSystem`-Sentinel, `cli.isFilesystemError`-Erweiterung. | ~40 |
| T3 | **Application-Layer**: `InitProjectService.fsFactory`-Feld + `NewInitProjectServiceWithFactory`-Konstruktor + `Init()`-Wrapper mit Mutex/Swap analog `AddServiceService.Add()`; `mapCaptureToPlannedFiles(captured, req.BaseDir)`. FS-Write-Wraps mit `%w: ErrInitFileSystem`. | ~150 |
| T4 | **Composition-Root-Wiring** in `cmd/uboot/main.go`: `initFSFactory`-Closure analog `addFSFactory`. | ~30 |
| T5 | **CLI-RunE**: `init`-RunE auf den gemeinsamen Error-Envelope-Helper aus T1, drei JSON-Pfade analog add, Allowlist-Migration (`"u-boot init": true`), Reject-Pin-Test 10→9. | ~140 |
| T6 | **Acceptance-Tests**: 8 Flag-Kombinationen × frisch-CWD + Soft-Existing-Pin + Backup-Pfad-Pin + Template-Modus-Pin + Mid-Write-Failure + 3-Flag-Combo. | ~280 |
| T7 | **Closure**: CHANGELOG, `cli-json-output.md` §6+§6.4+§7-Update (init-Zeile in Mutations-Matrix), roadmap-Update (3/9 done), Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) |

LOC-Bilanz: ~740 LOC — Pattern-Erbe spart ~50 % gegenüber add (war
~1380 LOC inkl. Renderer-Erfindung).

## Out of Scope

- **Backup-Konsistenz-Re-Validation** (Read-after-Write): falls
  init in Zukunft ein Cleanup-Detect-Schritt hinzufügt, der nach
  einem Backup-Copy den Ziel-Inhalt liest, müsste der Recorder
  einen Overlay-Map-Cache ergänzen. Heute nicht der Fall.
- **Template-Catalog-Erweiterung**: neue Templates landen in einem
  eigenen Slice; init-Slice liefert nur die JSON-Migration des
  bestehenden `basic`-Templates.
- **Generisches `mapErrorToDiagnostic`-Registry**: Altitude-Reviewer-
  Vorschlag aus add R6 #I1 (sentinel→LH-Code-Registry). Cluster-
  T_close-Aufgabe, nicht init-Aufgabe.
- **Generischer `previewFSFactory`-Konstruktor** in
  `cmd/uboot/main.go`: Altitude-Reviewer-Vorschlag aus add R6 #I3.
  Folgt erst, wenn 3+ Subcommands ihre eigenen Factories haben und
  das Drift-Risiko sichtbar wird.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  §T0-Outcomes — Vorgaben für den Folge-Slice-Block.
- Pattern-Vorbild:
  [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
  — T0-T6 + Review-Rounds 6-8 voll abgeschlossen. Erbschafts-
  Disziplin in §T0-(a) dieses Slices.
- Spec: `LH-FA-CLI-007/008`, `LH-NFA-USE-004`,
  `LH-FA-INIT-001..007`, `LH-NFA-REL-003`
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- Code-Anker heute:
  [`initproject.go`](../../../../internal/hexagon/application/initproject.go)
  (1079 LOC, Skeleton-Dirs + WriteFiles + BackupPath-Indirektion),
  [`backup.go`](../../../../internal/hexagon/application/backup.go)
  (213 LOC, `BackupPath`-Helper),
  [`cli/init.go`](../../../../internal/adapter/driving/cli/init.go)
  (236 LOC, RunE-Erweiterungs-Ziel),
  [`cmd/uboot/main.go`](../../../../cmd/uboot/main.go) (`initSvc`-
  Konstruktor-Migration auf `WithFactory`-Form).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
