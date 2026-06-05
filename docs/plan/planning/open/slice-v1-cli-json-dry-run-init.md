# Slice V1: `init --json` / `--dry-run` / `--diff` — modifying-Surface erbt von Add

> **Status:** geplant für v0.4.0 — dritter Folge-Slice des Cluster-
> Slice [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 3/9). Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
> 1:1 für die Carrier-Types, den `RecordingFileSystem`-driven-
> Adapter, den Pure-Go Diff-Renderer, das `previewModeFromFlags`-
> Mapping und die Error-Envelope-Pipeline; init-spezifisch sind
> die Mutations-Matrix (`MkdirAll` + `WriteFile` direkt plus
> `CopyExclusive`/`Mkdir`/`MkdirAll`/`Copy`/`RemoveAll` indirekt
> via `BackupPath` — sechs der acht `driven.FileSystem`-Mutations-
> Methoden; `WriteFileExclusive` und `Rename` werden NICHT aus
> init-Pfaden gerufen, Recorder deckt sie als Drift-Schutz
> trotzdem ab), sieben init-spezifische LH-Codes
> (`LH-FA-INIT-001`..`-007`), die Soft-Existing-Detection
> (`LH-FA-INIT-004` für Marker-Kollision; `LH-FA-INIT-005` für
> Backup-/Force-Failures) und der Template-Modus
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

| Phase | Methode | Pfade (Default, ohne `--template`/`--devcontainer`) | Indirekt? |
| --- | --- | --- | --- |
| Skeleton-Dirs (Z. 776, `writeDirectories` → `projectStructureDirs`) | `MkdirAll` | `docker/`, `scripts/`, `docs/` (immer); `.devcontainer/` (nur bei `--devcontainer`) | direkt |
| Skeleton-Files (Z. 818/866/968, `executeTemplatedFiles` → `fileTemplates`) | `WriteFile` | `README.md`, `CHANGELOG.md`, `compose.yaml`, `.env.example`, `.gitignore` (in dieser Aufruf-Reihenfolge); devcontainer-Files (nur bei `--devcontainer`) | direkt |
| u-boot.yaml (Z. 302/865, `executeUBootYAML`) | `WriteFile` | `u-boot.yaml` (ZULETZT — nach Dirs und Skeleton-Files; LH-FA-INIT-002 anchor) | direkt |
| Backup (Z. 978, [`backup.go`](../../../../internal/hexagon/application/backup.go)) | `CopyExclusive` (Z. 139), `Mkdir` (Z. 149), `MkdirAll` (Z. 198), `Copy` (Z. 209), `RemoveAll` (Z. 88) | `<file>.bak.<n>` plus Backup-Verzeichnis | indirekt via `BackupPath` |

Damit nutzt init **sechs der acht** `driven.FileSystem`-Mutations-
Methoden in der Praxis (`WriteFile`, `MkdirAll` direkt;
`CopyExclusive`, `Mkdir`, `MkdirAll`, `Copy`, `RemoveAll` über
Backup). `WriteFileExclusive` und `Rename` werden aus keinem
init-Pfad gerufen (Cluster-§499-502 dokumentiert das). Der
`RecordingFileSystem` aus add-T1-B deckt aber **alle 8** ab
(Drift-Schutz war Cluster-Pflicht) — kein neuer driven-Adapter
nötig.

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

**Wichtig: Path-Anchor-Klärung** (T0-(k)): Heutiges `init <name>`
nutzt `BaseDir = cwd` und `req.Name = <name>`; der Use-Case
schreibt **direkt in cwd** (`filepath.Join(cwd, "u-boot.yaml")`),
**NICHT** in ein `<cwd>/<name>/`-Subdir. Die Pin-Form unten
spiegelt das tatsächliche Verhalten. Falls ein Folge-Slice das
Verhalten auf `cwd/<name>/` ändert, muss die Pin-Form mitwandern
(T0-(k) Sub-Decision dokumentiert diese Entscheidung explizit).

Konkrete Pin-Form für `init myproj --dry-run --json` (frisch-CWD,
**Default-Pfad** ohne `--template`/`--devcontainer`):

```json
{
  "status": "ok",
  "command": "init",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    {"path": "docker", "action": "create"},
    {"path": "scripts", "action": "create"},
    {"path": "docs", "action": "create"},
    {"path": "README.md", "action": "create"},
    {"path": "CHANGELOG.md", "action": "create"},
    {"path": "compose.yaml", "action": "create"},
    {"path": ".env.example", "action": "create"},
    {"path": ".gitignore", "action": "create"},
    {"path": "u-boot.yaml", "action": "create"}
  ],
  "changes": [/* identische 9 Einträge mit count = CountLines(NewContent) */],
  "diagnostics": [],
  "exitCode": 0
}
```

Reihenfolge folgt der tatsächlichen Aufruf-Sequenz in
`InitProjectService.Init()` (Z. 289-309): (1) `writeDirectories`
→ `docker`, `scripts`, `docs`; (2) `executeTemplatedFiles` →
`README.md`, `CHANGELOG.md`, `compose.yaml`, `.env.example`,
`.gitignore` (Reihenfolge aus `fileTemplates()` Z. 566+);
(3) `executeUBootYAML` → `u-boot.yaml` (ZULETZT). Bei
`--devcontainer` wird `.devcontainer/` als vierter Dir-Eintrag
plus devcontainer-Files ergänzt — siehe T0-(m) Flag-Matrix-
Coverage.

Negative-Pin: bei `--dry-run` null Production-FS-Mutationen, gleicher
Spy-Mechanismus wie in add T5 (Recorder schickt nichts an die
underlying-FS bei `WithPassthrough(false)`).

Soft-Existing-Detection-Pin (`LH-FA-INIT-004`):
`u-boot init myproj --dry-run --json` auf eine **existierende**
Projekt-CWD ohne `--backup`/`--force`/`--no-interactive` liefert
einen Error-Envelope (drei Disambiguatoren, nicht zwei; siehe
`checkSoftExisting` in initproject.go Z. 478-507). Error-Message
folgt `softExistingAbort` (Z. 531-534) Format
`"%d structure elements detected (...) via ..."`; der CLI-Adapter
normalisiert ggf. via `mapErrorToDiagnostic` — Sub-Decision T0-(f):

```json
{
  "status": "error",
  "command": "init",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [],
  "changes": [],
  "diagnostics": [
    {"level": "error", "code": "LH-FA-INIT-004", "message": "<softExistingAbort message>"}
  ],
  "exitCode": 10
}
```

**Mid-Write-Failure-Pin** (`LH-NFA-REL-003`, T0-(o)):
`u-boot init myproj --dry-run --diff --json` mit FS-Failure bei
File-Index N im Use-Case-Pfad — Recorder hat die ersten N
Mutationen captured, Use-Case wrapt `os.WriteFile`-Error als
`fmt.Errorf("...: %w", driving.ErrInitFileSystem)`:

```json
{
  "status": "error",
  "command": "init",
  "dryRun": true,
  "diff": true,
  "plannedFiles": [/* die ersten N+1 captured-Records bis Failure-Stelle */],
  "changes": [/* identische Liste mit count = CountAdditions/CountLines pro Action */],
  "diagnostics": [
    {"level": "error", "code": "LH-NFA-REL-003", "file": "<failing path>", "message": "..."}
  ],
  "exitCode": 14
}
```

T6-Acceptance pinnt mindestens zwei Failure-Positionen: früh
(MkdirAll-Fehler nach Index 1) und spät (WriteFile-Fehler bei
u-boot.yaml als letztem File). Roll-back-aware Capture (alle
bereits geschriebenen Files reverten) ist Out-of-Scope (V1, gleiches
Argument wie add-T0-(b)).

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
- ✅ **`mapErrorToDiagnostic` für init** (T0-(f), Order-Pflicht
  siehe dort): per-subcommand Switch analog add, init-spezifische
  LH-Codes (Mapping respektiert `cli.go:217-224`-Konvention für
  bereits etablierte Sentinels — abweichend von der naiven
  „Spec-Code = LH-Sektionsnummer"-Lesart):
  - `LH-FA-INIT-004`: `driving.ErrProjectExists`,
    `driving.ErrFileExists` (Marker-Kollision, „Bestehendes
    Projekt erkennen" §567)
  - `LH-FA-INIT-005`: `driving.ErrConfirmationRequired`,
    `driving.ErrForceRequiresBackup`,
    `driving.ErrBackupUnsupportedKind`,
    `driving.ErrBackupSuffixExhausted`,
    `driving.ErrBackupSourceMissing` (Überschreibschutz §595-619)
  - `LH-FA-INIT-006`: `domain.ErrInvalidProjectName` UND
    `domain.ErrInvalidServiceName` (Name-Validierung §625;
    Konvention aus `add.go:410` weitergeführt — Carveout-Pin
    siehe T0-(f) Footnote, dass §625 strikt nur „Projektname"
    benennt)
  - `LH-FA-CLI-006`: `driving.ErrTemplateConflictsWithFlag` (Usage-
    Error, Exit-Code 2 via `isUsageError`)
  - `LH-NFA-REL-003`: neuer `driving.ErrInitFileSystem`-Sentinel
    für FS-Failures (analog `ErrAddFileSystem` → Exit-Code 14
    via `isFilesystemError`)
  - LH-FA-INIT-001/-002/-003/-007 sind heute ohne dedizierten
    Sentinel — Spec-Anker für Use-Case-Phasen, kein Error-Pfad.
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
wird. **Alias-Lebensdauer-Pflicht**: `AddPreviewMode` ist die
EINZIGE Service-Prefix-Alias. Folge-Slices (remove/generate/
config-set) referenzieren `driving.PreviewMode` direkt — keine
weiteren `XxxPreviewMode`-Aliases. Carveout-Liste in Cluster-
T_close enthält damit nur eine Alias-Zeile (statt 5+).

**Project-relative paths erblich**: Add-Round-8 Finding A
etablierte `mapCaptureToPlannedFiles(records, baseDir)` als
inverse Strippung zu `filepath.Join(baseDir, …)`. Init muss
dieses Mapping 1:1 erben — sonst leakt das Envelope absolute
cwd-Pfade. Init-Pflicht: T3 ruft
`mapCaptureToPlannedFiles(recorder.Captured(), req.BaseDir)`,
T6-Acceptance-Test verifiziert dass kein `PlannedFile.Path` mit
`"/"` oder `req.BaseDir` beginnt. Mit der Path-Anchor-Klärung
aus T0-(k) (heute BaseDir=cwd, kein cwd/name-Subdir) ergibt das
die in §Aufhebungsbedingung pinned bare-basename-Form.

### T0-(d) `InitProjectService.fsFactory` Konstruktor-Form

Add hat einen zweiten Konstruktor (`NewAddServiceServiceWithFactory`)
und behält den Legacy-Konstruktor für Backward-Compat
(`addservice.go:203`). Tests, die den legacy-Konstruktor benutzen,
zeigen `PlannedFiles: nil` (Recorder nil).

**Vorschlag (T0-Festlegung)**: gleiche Form für init — neuer
`NewInitProjectServiceWithFactory(fsFactory, …)`-Konstruktor neben
dem heutigen `NewInitProjectService(fs, …)`. Legacy bleibt für
existierende Tests funktional. **`initMu sync.Mutex`-Pflicht**:
`InitProjectService` bekommt das gleiche Mutex-Pattern wie
`AddServiceService.addMu` (Add-Review #10) — der per-Request
s.fs-Swap zwischen Goroutines würde sonst racen. Lock/defer-
Unlock umschließt `selectFS`/`s.fs`-Swap/`runInit`/`mapCapture`.
T6 ergänzt einen konkurrenten Init-Pin-Test.

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

**Verhaltens-Verträge (1:1 aus add-R6 erblich)**:

- **Exit-Code-Propagation (add #2)**: nach erfolgreichem
  Envelope-Write returnt der gemeinsame `reportError` IMMER den
  original `addErr` — `cli.ExitCode(err)` würde sonst 0
  zurückgeben trotz envelope-claimed 14. T6 pinnt das mit
  `errors.Is(err, ErrInitFileSystem)` + `cli.ExitCode(err) == 14`.
- **Broken-pipe-Propagation (add #3)**: `writeInitDiff` und
  `printInitSummary` returnen `error` (analog `writeAddDiff` /
  `printAddSummary` — beide returnen heute `error`); runInit
  propagiert. Parität mit
  `TestDoctorJSON_BrokenPipePreservesExitCode`.
- **Mid-Failure-File-Annotation (add R6 #lastPlannedPath)**:
  `reportError` setzt `diag.File = lastPlannedPath(resp)` für
  Mid-Write-Failure-Diagnostics. Init erbt die Logik 1:1.
- **`wantsFullSchema`-Switch (add R6 #4)**:
  `len(resp.PlannedFiles) > 0 || dryRun || diffFlag` bestimmt
  voll- vs. minimal-Envelope auf Error-Pfad. Init-T1-Refactor
  zieht das in den gemeinsamen Helper, mit `dryRun`/`diff` als
  Parameter (nicht hardgecodet).

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

**Code-Map** (T0-Outcomes-Tabelle finalisiert; korrigiert
gegenüber dem ursprünglichen Stub gemäß
`cli.go:217-224`-Konvention und Lastenheft §-Nummern):

| Sentinel | Quell-Datei | LH-Code | Exit-Code |
| --- | --- | --- | --- |
| `domain.ErrInvalidProjectName` | `domain/projectname.go` | `LH-FA-INIT-006` | 10 |
| `driving.ErrProjectExists`, `driving.ErrFileExists` | `port/driving/initproject.go` | `LH-FA-INIT-004` | 10 |
| `driving.ErrConfirmationRequired` (shared) | `port/driving/down.go` | `LH-FA-INIT-005` | 10 |
| `driving.ErrForceRequiresBackup`, `driving.ErrBackupUnsupportedKind` | `port/driving/initproject.go` | `LH-FA-INIT-005` | 10 |
| `driving.ErrBackupSuffixExhausted`, `driving.ErrBackupSourceMissing` | `port/driving/initproject.go` | `LH-FA-INIT-005` | 14 |
| `domain.ErrInvalidServiceName` (geteilt mit add) | `domain/servicename.go` | `LH-FA-INIT-006` | 10 |
| `driving.ErrTemplateConflictsWithFlag` | `port/driving/initproject.go` | `LH-FA-CLI-006` | 2 |
| **`driving.ErrInitFileSystem` (neu)** | `port/driving/initproject.go` (T2) | **`LH-NFA-REL-003`** | **14** |

**Footnote — INIT-006-Carveout**: Spec §625 nennt LH-FA-INIT-006
strikt „Projektnamen-Validierung". Die etablierte Codebase-
Konvention (`cli.go:217-220`, `add.go:410`) erweitert das auf
Service-Name-Validierung. Init übernimmt diese Konvention; ein
dedizierter LH-Code für Service-Name-Validation bleibt
Cluster-T_close-Sub-Decision.

**Switch-Order-Pflicht (Add R6 #11 erblich)**:
`mapErrorToDiagnostic` für init MUSS `ErrInitFileSystem` als
ERSTEN `errors.Is`-case prüfen — `executeWriteFiles`/`backup.go`
wrappen FS-Errors als `fmt.Errorf("...: %w: %w", path,
ErrInitFileSystem, rawErr)` (Multi-`%w`). Falls künftiger Code
einen fachlichen Sentinel in die gleiche Chain wrapt, würde
fachlich-zuerst-Order die FS-Klassifikation (`LH-NFA-REL-003` /
Exit-Code 14) auf einen Exit-Code-10-Fachfehler downgraden.

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
am Recorder). Der Recorder sieht nur die Result-Writes. Das ist OK
für V1.

**Reads im Dry-Run sind erlaubt**: zusätzlich zu Catalog-Reads
liest init im Dry-Run-Pfad auch über `s.fs.Exists`/`Lstat`/
`ReadFile`:
- `checkSoftExisting` (initproject.go Z. 478-507) prüft 6
  LH-FA-INIT-003-Indikator-Pfade via `Exists`.
- `planFile`/`fileHasManagedBlock` lesen Templates und existierende
  Files.

Der `RecordingFileSystem` delegiert Reads grundsätzlich an den
underlying-FS und captured sie nicht (siehe `recordingfs.go`
Read-Methods Z. 90-120 plus Test
`TestRecordingFS_ReadsAlwaysDelegate`). T6-Spy-Check muss
zwischen **Writes** (verboten im Dry-Run, Counter MUSS 0 sein)
und **Reads** (erlaubt, Counter irrelevant) unterscheiden. Pin-
Doku in T6-Test-Kommentar explizit ausweisen.

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

**`actionReplaceBlock`-Sonderform (Re-init mit `--force`)**: bei
existierendem Projekt + `--force` (`executeReplaceBlock`-Pfad)
ist der WriteFile-Body nur der managed-block-Bereich, nicht das
ganze File. Recorder sieht: `WriteFile(compose.yaml, blockBody)`
— NewContent ist der Block, OldContent ist der bisherige Block
(VOR `WriteFile` via `s.snapshot`). Action: `modify`. count =
`CountAdditions(diffHunks(OldContent, NewContent))` — gleiche
Form wie add-modify (T0-(g) gilt unverändert). Edge-Case:
Content-identischer Block-Replace ergibt `CountAdditions = 0` und
keine Hunks (Idempotenz-Signalisierung im Envelope).

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

### T0-(k) Path-Anchor bei positional Project-Name + writeInitDiff-Verträge

**Path-Anchor**: heutiges `init <name>` setzt `req.BaseDir = cwd`
und `req.Name = <name>`; die Use-Case-Schichten schreiben direkt
in cwd (`filepath.Join(cwd, "compose.yaml")` etc.), **NICHT** in
`cwd/<name>/`. Drei Optionen:

1. **Status quo**: BaseDir=cwd, Pin-Form trägt bare basenames
   (`u-boot.yaml`, `compose.yaml`, `docker`). Init `<name>` ist
   nur ein Project-Identity-Field, kein Pfad-Prefix.
2. **Verhaltens-Änderung**: `init <name>` legt vorher
   `cwd/<name>/` an, BaseDir wird intern auf `cwd/<name>`
   umgesetzt; Pin-Form trägt `<name>/`-Prefix. Bricht heutige
   CLI-Verträge und integration-Tests (`main_test.go::TestRun_…`
   nutzen Option 1).
3. **Hybrid**: Flag `--in-subdir <name>` opt-in für Variante 2.

**Vorschlag (T0-Festlegung)**: **Option 1** — Status quo halten,
Pin-Form ohne `<name>/`-Prefix (siehe §Aufhebungsbedingung-JSON).
Verhaltens-Änderung wäre ein eigener Slice (Out-of-Scope V1).

**`writeInitDiff`-Verträge (Add R6 #15 erblich)**: vier Pflichten
gegenüber dem extraktiven `writeDiff`-Helper aus T1:
- Blank-Line-Separator zwischen Multi-File-Diffs
- `--- <path> (<action>)` Header pro File
- Binary-Hint `(binary content — diff suppressed)` bei
  `IsBinary`-Match
- `(no changes)`-Hint bei `diff.Compute(...)==nil` (content-
  identischer modify)

### T0-(l) Positional-Arg-Fallback `init` ohne `<name>`

`resolveProjectName` (initproject.go Z. 542) leitet bei
leerem `req.Name` den Project-Name aus `filepath.Base(req.BaseDir)`
ab. Das ist eine UX-relevante CLI-Form: `cd /tmp/foo && u-boot
init --dry-run --json` produziert ein Envelope mit project=`foo`
in `u-boot.yaml`.

**Sub-Decision**: Acceptance-Test ergänzt einen Pin für die
fallback-Form mit deterministischem cwd-basename (testfaker via
WithGetwd). plannedFiles[]-Output ist identisch zur
positional-Form (Path-Anchor T0-(k) Option 1), nur das geschriebene
`u-boot.yaml` enthält den anderen Project-Name. Doku-Hint in
`cli-json-output.md` §6.4.

### T0-(m) Flag-Matrix-Coverage im Aufhebungsbedingung-Pin

Init hat 5 verhaltens-modifizierende Flags neben `--dry-run`/
`--diff`/`--json`: `--devcontainer`, `--template <name>`,
`--no-git`, `--force`, `--backup`, plus `--allow-external-feature-
sources`. Jede ändert den plannedFiles[]-Shape unterschiedlich.

**Flag-Matrix-Pin-Plan** (T0-Festlegung):

| Flag-Set | Default-Pin | Pflicht in T6? |
| --- | --- | --- |
| Default (kein Flag) | §Aufhebungsbedingung-JSON oben | ja |
| `--devcontainer` | + `.devcontainer/` Dir + 2 Files | ja (`+3` Einträge) |
| `--template basic` | identisch zu Default (basic template = default file-set) | ja |
| `--no-git` | identisch zu Default (`--no-git` wirkt POST-write, nicht im Recorder) | ja (Doku-Hint) |
| `--force` (Re-init) | `action: "modify"` für `compose.yaml`/`u-boot.yaml` + `actionReplaceBlock` Pfad | ja (T0-(h)) |
| `--backup` (Re-init) | + Backup-Files als `action: "create"` | ja (Backup-Pfad-Pin) |
| `--allow-external-feature-sources X` | u-boot.yaml content-Variation; plannedFiles[]-Liste identisch zu Default | optional (T6 deckt es als Acceptance-Cluster) |

Pflicht-Pins in T6 sind die ersten sechs (default + 5
Variationen); `--allow-external-feature-sources` ist optional
weil reine Content-Variation des u-boot.yaml-Write.

## T0-Outcomes

Verbindliche Festzurrung wandert nach `next/`-Übergang in diesen
Block — analog add-Slice. Erwartete Form: Tabelle mit den 9 Sub-
Decisions plus Implementations-Pflicht-Spalte (T1-T6).

## Tranchen (vorgeschlagen)

LOC-Schätzung schlanker als add (das alles neu erfunden hat);
init erbt das Pattern, aber die Mutations-Matrix ist breiter und
die Acceptance-Test-Matrix größer.

| T | Inhalt | LOC | Voraussetzung |
| - | ------ | --- | --- |
| T0 | Discovery + Sub-Decisions klären; Pattern-Erbe-Tabelle pinnen | — (Plan) | — |
| T1 | **Refactor-Tranche**: `previewModeFromFlags` extrahieren (T0-(b)), `driving.PreviewMode` umbenennen + `AddPreviewMode`-Alias (T0-(c)), `reportError`/`writeErrorEnvelope`/`writeDiff`-Helper generalisieren (T0-(e)). Add-Tests + Add-Godoc-Wahrheitstabellen (5 Files, ~30 LOC Comment-Edits) mitgezogen. | ~130 | T0 |
| T2 | **Port-Types + Sentinel**: `driving.InitProjectRequest.PreviewMode`-Field, `driving.ErrInitFileSystem`-Sentinel, `cli.isFilesystemError`-Erweiterung. Unit-Test für Sentinel-Identity + ExitCode-Routing. | ~50 | T0 |
| T3 | **Application-Layer**: `InitProjectService.fsFactory`-Feld + `initMu sync.Mutex` + `NewInitProjectServiceWithFactory`-Konstruktor + `Init()`-Wrapper mit Mutex/Swap analog `AddServiceService.Add()`; `mapCaptureToPlannedFiles(captured, req.BaseDir)`. **Neun FS-Wrap-Stellen** (initproject.go Z. 776/818/866/968 plus backup.go Z. 88/139/149/198/209) mit `%w: ErrInitFileSystem` umhüllt. Factory-Tests analog `addservice_factory_test.go` (~200 LOC Test-Datei). | ~280 | T2 |
| T4 | **Composition-Root-Wiring** in `cmd/uboot/main.go`: `initFSFactory`-Closure analog `addFSFactory`. | ~30 | T3 |
| T5 | **CLI-RunE**: `init`-RunE auf den gemeinsamen Error-Envelope-Helper aus T1, drei JSON-Pfade analog add, Allowlist-Migration (`"u-boot init": true`), Reject-Pin-Test `TestRootJSON_RejectsAllNonMigratedForms` (cases-Slice 10→9). | ~140 | T1 + T2 |
| T6 | **Acceptance-Tests** (T0-(m)-Matrix): 6 Pflicht-Flag-Sets (default + `--devcontainer`/`--template basic`/`--no-git`/`--force`/`--backup`) × 4 JSON-Flag-Kombos = ~13 Top-Level-Tests; Soft-Existing-Pin (3 Disambiguatoren); Mid-Write-Failure-Pin (zwei Positionen: früh + spät); 3-Flag-Combo `--dry-run --diff --json`; Concurrent-Init-Mutex-Pin; Path-Anchor-Pin (`PlannedFile.Path` ist project-relativ, kein `/`-Prefix). Table-driven wo möglich, sonst individuell. | ~400 | T5 |
| T7 | **Review-Fix-Rounds** (~1-2 Runden bei Pattern-Erbe; add hatte R6/R7/R8): Diff aus Reviewer-Findings konsolidieren, Fixes als eigene Sub-Commits, DoD-Hash-Tabelle ergänzen. | ~80 | T6 |
| T8 | **Closure**: CHANGELOG-Eintrag, `cli-json-output.md` §6-Tabelle (init→done) + §6.1-Reject-Liste (init raus) + §6.4 neue init-Sektion + §7 Mutations-Matrix (init-Zeile), roadmap-Update (3/9 done), Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz: ~1110 LOC — Pattern-Erbe spart ~20 % gegenüber add
(war ~1380 LOC inkl. Renderer-Erfindung). Init ist umfangreicher
als ursprünglich geschätzt, weil die Mutations-Matrix breiter
(9 statt 2 wrap-sites), die Test-Matrix größer (13 statt 7
Szenarien), und der T1-Refactor sowohl init-Vorbereitung als
auch Add-Code-Migration umfasst.

**Reihenfolge-Pflicht**: T1 und T2 sind parallel; T3 wartet auf
T2 (braucht InitProjectRequest.PreviewMode + ErrInitFileSystem);
T5 wartet auf T1 + T2 (braucht den gemeinsamen Helper + Sentinel-
Mapping); T6 wartet auf T5; T7 wartet auf T6; T8 wartet auf T7.

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
