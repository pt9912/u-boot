# Slice V1: `init --json` / `--dry-run` / `--diff` — modifying-Surface erbt von Add

> **Status:** ✅ **done** — dritter Folge-Slice (3/9) des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 3/9). Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-add`](slice-v1-cli-json-dry-run-add.md)
> 1:1 für die Carrier-Types, den `RecordingFileSystem`-driven-
> Adapter, den Pure-Go Diff-Renderer, das `previewModeFromFlags`-
> Mapping und die Error-Envelope-Pipeline; init-spezifisch sind
> die Mutations-Matrix (`MkdirAll` + `WriteFile` direkt plus
> `CopyExclusive`/`Mkdir`/`MkdirAll`/`Copy`/`RemoveAll` indirekt
> via `BackupPath` — sechs der acht `driven.FileSystem`-Mutations-
> Methoden; `WriteFileExclusive` und `Rename` werden NICHT aus
> init-Pfaden gerufen, Recorder deckt sie als Drift-Schutz
> trotzdem ab), sieben init-spezifische LH-Codes als Spec-Anker
> ([`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung)) — davon drei mit dedizierten
> Sentinels in der `mapErrorToDiagnostic`-Map (INIT-004 für
> Marker-Kollision; INIT-005 für `--force`/`--backup`-Usage-
> Failures; INIT-006 für Name-Validierung), die anderen vier
> rein als Phasen-Anker. **[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) für Backup-FS-
> Failures** (Suffix-Exhaustion, Source-Missing) und Mid-Write-
> FS-Failures. Template-Modus (`init --template <name>`) ist in
> V1 mutex zu `--dry-run`/`--diff` (siehe T0-(i) Out-of-Scope-
> Carveout).
>
> Init war der erste Folge-Slice, der vom Add-Pattern erbt — die
> Erbschafts-Disziplin (was 1:1, was init-spezifisch) ist in
> §T0-Discovery festgezurrt. Slice wanderte über drei Pre-`next/`-
> Review-Runden (R1/R2/R3 ≈ 93 Findings insgesamt; 17 Sub-
> Decisions (a)-(q) verbindlich festgezurrt inkl. der vier R3-
> Adversarial-Funde T0-(n)/(o)/(p)/(q) initGit-Skip /
> ProgressPort-Silencing / Context-Cancellation / Planning-Phase-
> Failures) und Review-Round-9 nach T6 (6 Findings, 4 R-Commits
> R1-R4, 1 Folge-Slice ausgelagert, 1 trivialer Dead-Code-
> Cleanup). Cluster-Stand nach init-Closure: **3/9 done**.
>
> **DoD-Tranchen-Hashes** (alle T0-T8 + R-Runden):
>
> | Tranche / Round | Inhalt | Commit |
> | --- | --- | --- |
> | T0 — Stub | `open/`-Stub Cluster-Folge-Slice 3/9 | `c79c4d2` |
> | T0 — R2-Findings | Pre-`next/` Review-Round-2 (31 Findings, 4 Angles) | `e45d30f` |
> | T0 — `open/→next/` | Lifecycle-Übergang | `bac9463` |
> | T0 — `next/→in-progress/` | Lifecycle-Übergang | `6ceb2a9` |
> | T0 — T1-E Re-Sequenz | Helper-Generalisierung nach T5 verschoben (Variante B) | `8c933d7` |
> | T1-A | `driving.AddPreviewMode → PreviewMode` + Type-Alias | `ad56550` |
> | T1-B | `previewModeFromFlags` nach `cli/previewmode.go` + Test-File-Move | `8ea5359` |
> | T1-C | `recordImplicitMkdir`-Dedup via `knownDirs` (R2-Finding B-4) | `8ba0250` |
> | T1-D | `mapResponseToWire → mapPlannedFilesToWire` nach `wireshapes.go` | `b058fb9` |
> | T1-E | Helper-Generalisierung (revertiert, wandert nach T5) | `8b858eb` / Revert `94dd78a` |
> | T1-F | Interne Refs auf kanonischen `driving.PreviewMode` | `9ed3e34` |
> | T2 | Port-Types `InitProject.PreviewMode`/`SilenceProgress`/`PlannedFiles` + `ErrInitFileSystem`-Sentinel | `a883870` |
> | T3 | Application-Layer: `fsFactory` + `initMu` + `initGit`-Skip + ProgressPort-Swap + 9 FS-Wrap-Stellen | `22d8402` |
> | T4 | Composition-Root-Wiring `initFSFactory`-Closure in `cmd/uboot/main.go` | `ab67de3` |
> | T5 | Helper-Generalisierung (`reportError`/`writeErrorEnvelope`/`writeDiff`/`lastPlannedPath`) + init-RunE für JSON/Dry-Run/Diff | `0689b32` |
> | T6 | Acceptance-Pins — Flag-Matrix + Error-Scenarios + T0-Pflicht-Pins | `80d5624` |
> | T6+ | Coverage-Bump auf 91.00 % (+1.0 % Sicherheitsabstand) | `bab6b13` |
> | R1 | `mapInitErrorToDiagnostic` erfasst `ErrInvalidFeatureSource` (+ R5 Dead-Code-Cleanup) | `6e5ad01` |
> | R2 | `initFromTemplate` Defense-in-Depth-Guard für `PreviewMode` | `e897fa7` |
> | R3 | `runBackup` Wrap-Strategie pinnen — raw FS vs typed Sentinel | `e10b57d` |
> | R4 | T0-(k) Path-Anchor Acceptance-Pin für positional `<name>` | `ee30c3c` |
> | T7 — Doku-Closure | Review-Round-9-Tabelle + Folge-Slice-Stub [`slice-v1-cli-cleanup-add-backup-error-class`](slice-v1-cli-cleanup-add-backup-error-class.md) | `d7f9e65` |
> | T8 — Closure | CHANGELOG, `cli-json-output.md` §6/§6.4/§7, `cli.go`-Godoc-Backup-Sentinels-Korrektur, roadmap-Update, `open/[slice-v1-cli-cleanup-add-preview-mode-alias](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md)`-Stub, Slice in `done/` | dieser Commit |

## Auslöser

Cluster-Slice [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md) §T0-Outcomes (a)+(b)+(e)
machen jeden modifying-Subcommand für `--json`/`--dry-run`/`--diff`
verbindlich ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) §1813, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run) §326,
[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) §451-489). `u-boot init` ist nach `add` der zweite
modifying-Subcommand und der **wichtigste Onboarding-Use-Case**
(Cluster-§T0-Discovery Z. 320 nennt `init --dry-run --diff --json`
explizit als Beispiel-Hauptanwendung): ein neuer Nutzer will sehen,
was `u-boot init <project>` an Files/Dirs anlegt, bevor er das auf
einer existierenden Codebase ausführt.

Spec-Bezug (geerbt von add-Slice):

- [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run) (Dry-Run, Voll-Schema §326)
- [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) (Diff, §451-489)
- [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) (Minimalkontrakt §1841)

Init-spezifische Sentinels und Spec-Stellen:

- [`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung) (Projekt-Skeleton, Verzeichnisstruktur,
  Soft-Existing-Detection, Backup-Pfad-Failures, Service-Name-
  Validation)
- [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) (FS-Failure-Klasse, geerbt für Mid-Write-Failure
  analog `add`)

Heute-Stand-Pre-Scan
(`internal/hexagon/application/initproject.go`, 1079 LOC):

| Phase | Methode | Pfade (Default, ohne `--template`/`--devcontainer`) | Code-Anker |
| --- | --- | --- | --- |
| Skeleton-Dirs (`writeDirectories` Z. 768 → `projectStructureDirs` Z. 30) | `MkdirAll` (Call Z. 776) | `docker/`, `scripts/`, `docs/` (immer); `.devcontainer/` (nur bei `--devcontainer`) | direkt |
| Skeleton-Files (`executeTemplatedFiles` → `fileTemplates()` in [`templates.go`](../../../../internal/hexagon/application/templates.go) Z. 73-81) | `WriteFile` (Calls Z. 818 actionWrite, Z. 866 actionReplaceBlock, Z. 968 actionOverwriteFull) | `README.md`, `CHANGELOG.md`, `compose.yaml`, `.env.example`, `.gitignore` (in dieser Aufruf-Reihenfolge aus `fileTemplates()`); devcontainer-Files (nur bei `--devcontainer`) | direkt |
| u-boot.yaml (`Init()` Z. 302 ruft `executeUBootYAML` Z. 1037 → `executeFile` Z. 814 → `WriteFile` Z. 818/866/968) | `WriteFile` | `u-boot.yaml` (ZULETZT — nach Dirs und Skeleton-Files; [`LH-FA-INIT-002`](../../../../spec/lastenheft.md#lh-fa-init-002-projektname) anchor) | direkt |
| Backup (Aufrufer: `initproject.go` Z. 978 `runBackup` → [`backup.go`](../../../../internal/hexagon/application/backup.go) `BackupPath` Z. 57) | `RemoveAll` (Z. 88), `CopyExclusive` (Z. 139), `Mkdir` (Z. 149), `MkdirAll` (Z. 198), `Copy` (Z. 209) | `<file>.bak.<n>` plus Backup-Verzeichnis | indirekt via `BackupPath` |

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
`InitProjectService.Init()` (Z. 245-318), Execute-Phase Z.
289-309: (1) `writeDirectories` → `docker`, `scripts`, `docs`;
(2) `executeTemplatedFiles` → `README.md`, `CHANGELOG.md`,
`compose.yaml`, `.env.example`, `.gitignore` (Reihenfolge aus
`fileTemplates()` in `templates.go` Z. 73-81); (3)
`executeUBootYAML` Z. 1037 → `u-boot.yaml` (ZULETZT). Bei
`--devcontainer` wird `.devcontainer/` als vierter Dir-Eintrag
plus devcontainer-Files ergänzt — siehe T0-(m) Flag-Matrix-
Coverage.

Negative-Pin: bei `--dry-run` null Production-FS-Mutationen, gleicher
Spy-Mechanismus wie in add T5 (Recorder schickt nichts an die
underlying-FS bei `WithPassthrough(false)`).

Soft-Existing-Detection-Pin ([`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen)):
`u-boot init myproj --dry-run --json` auf eine **existierende**
Projekt-CWD ohne `--backup`/`--force`/`--no-interactive` liefert
einen Error-Envelope (drei Disambiguatoren, nicht zwei; siehe
`checkSoftExisting` in initproject.go Z. 478-508). Error-Message
folgt `softExistingAbort` (Z. 531-534) Format exakt:
`"%w: %d structure elements detected (%s) via %s; add --backup or --force to re-init"`
(wraps `driving.ErrProjectExists` als ersten `%w`); der CLI-Adapter
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

**Template-Reject-Pin** ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes), T0-(i) Out-of-Scope-
Carveout für V1): `u-boot init myproj --template basic --dry-run
--json` rejects am CLI-RunE-Level (T5 mutex-check vor uc.Init-
Call):

```json
{
  "status": "error",
  "command": "init",
  "diagnostics": [
    {"level": "error", "code": "LH-FA-CLI-006", "message": "--template is mutually exclusive with --dry-run/--diff (V1 carveout — see slice-v1-cli-json-dry-run-template-preview)"}
  ],
  "exitCode": 2
}
```

Minimal-Envelope (kein plannedFiles/changes), weil die Validation
VOR jedem Recorder-Setup fired.

**Planning-Phase-Force-Failure-Pin** ([`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz), T0-(q)):
`u-boot init myproj --force --dry-run --json` auf CWD mit
existierender `.gitignore` (kein managed-Block, kein `--backup`)
failed im planFile bevor irgendein Write das Capture erreicht.
`plannedFiles[]` bleibt leer:

```json
{
  "status": "error",
  "command": "init",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [],
  "changes": [],
  "diagnostics": [
    {"level": "error", "code": "LH-FA-INIT-005", "message": "<ErrForceRequiresBackup message>"}
  ],
  "exitCode": 10
}
```

**Mid-Write-Failure-Pin** ([`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern), T0-(f) Switch-Order
+ T0-(k) writeInitDiff-Verträge):
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
  - [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen): `driving.ErrProjectExists`,
    `driving.ErrFileExists` (Marker-Kollision, „Bestehendes
    Projekt erkennen" §567)
  - [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz): `driving.ErrConfirmationRequired`,
    `driving.ErrForceRequiresBackup`,
    `driving.ErrBackupUnsupportedKind` (Überschreibschutz §595-619
    Usage-Klasse → Exit 10)
  - [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern): `driving.ErrBackupSuffixExhausted`,
    `driving.ErrBackupSourceMissing` (FS-Klasse: Suffix-Exhaustion
    und Source-Missing sind technische Filesystem-Failures, kein
    User-Action; Exit 14 via `isFilesystemError` — präzisere
    Klassifikation als Spec §605/§619 dem User die richtige Klasse
    signalisiert)
  - [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung): `domain.ErrInvalidProjectName` UND
    `domain.ErrInvalidServiceName` (Name-Validierung §625;
    Konvention aus `add.go:410` weitergeführt — Carveout-Pin
    siehe T0-(f) Footnote, dass §625 strikt nur „Projektname"
    benennt)
  - [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes): `driving.ErrTemplateConflictsWithFlag` (Usage-
    Error, Exit-Code 2 via `isUsageError`)
  - [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern): neuer `driving.ErrInitFileSystem`-Sentinel
    für FS-Failures (analog `ErrAddFileSystem` → Exit-Code 14
    via `isFilesystemError`)
  - [`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)/[`LH-FA-INIT-002`](../../../../spec/lastenheft.md#lh-fa-init-002-projektname)/[`LH-FA-INIT-003`](../../../../spec/lastenheft.md#lh-fa-init-003-projektstruktur-erzeugen)/[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung) sind heute ohne dedizierten
    Sentinel — Spec-Anker für Use-Case-Phasen, kein Error-Pfad.
- ✅ **Composition-Root-Wiring** (T4-Tranche, kein eigener T0-
  Sub-Decision-Slot weil Pattern-Erbe von add-T1-D 1:1):
  `initFSFactory` analog
  `addFSFactory` in `cmd/uboot/main.go`; gleiches Closure-Pattern.
  Pflicht-Erweiterung: `initSvc` migriert auf
  `NewInitProjectServiceWithFactory`. App-Struct + `cli.New(...)`
  bleiben unverändert.
- ✅ **Mutations-Matrix `cli-json-output.md` §7**: init-Zeile
  ergänzt — 6 von 8 Methoden via init-Pfaden (direkt: `WriteFile`,
  `MkdirAll`; indirekt via `BackupPath`: `RemoveAll`,
  `CopyExclusive`, `Mkdir`, `MkdirAll`, `Copy`);
  `WriteFileExclusive` + `Rename` ungenutzt aber Recorder-
  abgedeckt — siehe Pre-Scan-Tabelle oben.
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
Funktion ([slice-v1-cli-json-dry-run-add](slice-v1-cli-json-dry-run-add.md) Z. 114). Init braucht
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
Slice-T1-Tranche macht den Rename plus expliziter Type-Alias-
Syntax `type AddPreviewMode = PreviewMode` (**Gleichheits-Zeichen
ist Pflicht** — `type AddPreviewMode PreviewMode` wäre ein NEUER
Defined Type, der die Factory-Signatur `func(driving.AddPreviewMode)
(driven.FileSystem, driven.RecorderPort)` NICHT mehr assignment-
kompatibel zu `func(driving.PreviewMode) ...` macht und damit
addservice_factory_test.go bricht). Die Konstanten-Deklarationen
in `port/driving/addservice.go` Z. 150
(`PreviewNone AddPreviewMode = iota`) werden auf den kanonischen
Type umgestellt (`PreviewNone PreviewMode = iota`). Pin-Test
in T1: `var _ driving.AddPreviewMode = driving.PreviewMode(0)`
plus Factory-Signature-Identity-Check.

**Carveout-Plan-Pflicht** (MEMORY.md
[[feedback_carveouts_need_plans]]): die Alias-Lebensdauer „bis
Cluster-T_close" braucht einen eigenen Slice-Plan-Stub im
`open/`-Verzeichnis
([`slice-v1-cli-cleanup-add-preview-mode-alias`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md),
T8 dieses Slices legt ihn an). Ohne Plan wäre der
Carveout ein loser Hänger ohne Cleanup-Owner. Alternative:
**Alias als permanente Backward-Compat-Garantie** deklarieren
und das Cluster-T_close-Removal-Versprechen ganz fallen lassen.

**Vorschlag (T0-Festlegung)**: Cleanup-Plan-Stub-Variante —
init-T8 erzeugt
[`open/slice-v1-cli-cleanup-add-preview-mode-alias.md`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md)
mit Auslöser („Carveout aus init-Slice T0-(c)"),
einer AK („Alias-Decl raus, Verifikation via
addservice_factory_test.go") und LOC-Schätzung (~10 LOC, ein
git rm + ein paar Test-Aliases). **Alias-Lebensdauer-Pflicht**:
`AddPreviewMode` ist die EINZIGE Service-Prefix-Alias. Folge-Slices (remove/generate/
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
weg, ergänzt um drei konkret gepinnte Sub-Entscheidungen aus
Review-Round-2 Findings:

- **Helper-Signatur (Decomposed-Slices statt Response-Pointer)**:
  `reportError(out, err, plannedFiles, changes, dryRun, diff,
  command, mapErr)` und
  `writeErrorEnvelope(out, err, plannedFiles, changes, dryRun,
  diff, command, mapErr)` — nehmen `[]driving.PlannedFile` und
  `[]driving.ChangeEntry` als separate Parameter STATT eines
  `resp driving.AddServiceResponse`-Pointers. Begründung: init's
  `driving.InitProjectResponse` hat heute KEINE
  `PlannedFiles`/`Changes`-Felder; T2 ergänzt sie zwar, aber die
  decomposed-Form lässt T1 und T2 trotzdem parallel laufen und
  vermeidet die response-shape-Kopplung. CLI-Caller extrahiert
  `resp.PlannedFiles`/`resp.Changes` selbst beim Aufruf.

- **`mapErr`-Source-Pflicht**: jeder Subcommand-RunE definiert
  `mapErr := mapXxxErrorToDiagnostic` als erste Zeile im
  Funktions-Body und reicht den Function-Value an reportError
  weiter. Keine App-Struct-Erweiterung. **Add-Migration mit
  Rename**: T1 benennt heutiges `add.go::mapErrorToDiagnostic`
  in `mapAddErrorToDiagnostic` um (eine Decl + eine Call-Site
  bei add.go:246, ~5 LOC); init definiert `mapInitErrorToDiagnostic`
  parallel. Symmetrie über Subcommand-Prefix vermeidet Cross-Package-
  Name-Konflikte und macht das `mapErr := mapXxxErrorToDiagnostic`-
  Pattern in jedem RunE konkret nachvollziehbar.

- **Alternative-Form `errorEnvelopeReq`-Struct erwogen, abgelehnt**:
  Round-3 Reviewer vorschlug ein Config-Struct
  (`type errorEnvelopeReq struct { Out, Err, Planned, Changes,
  DryRun, Diff, Command, MapErr }`) gegen die 8-Param-Signatur.
  T0-Entscheidung: bei der decomposed-Form bleiben — Add-Pattern
  hat keine Config-Structs für andere Helper, einführen würde
  inkonsistent gegen das etablierte Pattern. Falls Cluster-T_close
  einen einheitlichen Helper-Stil etabliert, wandert die
  Migration in einen eigenen Cleanup-Slice.

- **`mapResponseToWire`-Migration-Pflicht**:
  `add.go::mapResponseToWire(resp, withHunks)` und das interne
  `computeChangeCountAndHunks(pf)` sind heute add-private.
  Decomposed-Helper-Signatur erfordert auch deren Migration:
  T1 benennt um auf `mapPlannedFilesToWire(pfs, withHunks)`
  und verschiebt nach `cli/wireshapes.go` (oder direkt in den
  `previewmode.go`-Block als allgemeiner cli-Helper).
  `computeChangeCountAndHunks` zieht mit. Zwei add-Call-Sites
  (add.go writeAddJSON Z. 227, writeAddErrorEnvelope Z. 264)
  aktualisieren. LOC-Delta ~15.

- **writeDiff-Header-Format (Option a, command-agnostisch)**:
  `writeDiff(out, plannedFiles)` emittiert `--- <path> (<action>)`-
  Header identisch für alle Subcommands. Per-command-Header-
  Overrides sind Out-of-Scope V1. Falls ein zukünftiger Subcommand
  (z.B. generate) eine andere Header-Form will, kann er writeDiff
  per Helper-Override umgehen.

Init-T1 macht den Refactor; add migriert in derselben Tranche
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

Plan-Bezug: [`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung). Sub-Decision: gibt es einen
init-spezifischen FS-Failure-Sentinel (analog `ErrAddFileSystem` →
[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern))?

**Vorschlag (T0-Festlegung)**: ja — neuer
`driving.ErrInitFileSystem`-Sentinel in `port/driving/initproject.go`,
gemappt auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/Exit-Code 14 in
`cli.isFilesystemError`. Wrap-Stellen in `initproject.go`
(`WriteFile`-Sites Z. 818/866/968 und backup-relevante Pfade)
ergänzen. Analog zum add-Pattern.

**Code-Map** (T0-Outcomes-Tabelle finalisiert; korrigiert
gegenüber dem ursprünglichen Stub gemäß
`cli.go:217-224`-Konvention und Lastenheft §-Nummern):

| Sentinel | Quell-Datei | LH-Code | Exit-Code |
| --- | --- | --- | --- |
| `domain.ErrInvalidProjectName` | `domain/projectname.go` | [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung) | 10 |
| `driving.ErrProjectExists`, `driving.ErrFileExists` | `port/driving/initproject.go` | [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) | 10 |
| `driving.ErrConfirmationRequired` (shared) | `port/driving/down.go` | [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) | 10 |
| `driving.ErrForceRequiresBackup`, `driving.ErrBackupUnsupportedKind` | `port/driving/initproject.go` | [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) | 10 |
| `driving.ErrBackupSuffixExhausted`, `driving.ErrBackupSourceMissing` | `port/driving/initproject.go` | [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) | 14 (heute schon via `cli.go::isFilesystemError` Z. 369-370 — siehe Doku-Korrektur unten) |
| `domain.ErrInvalidServiceName` (geteilt mit add) | `domain/servicename.go` | [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung) | 10 |
| `driving.ErrTemplateConflictsWithFlag` | `port/driving/initproject.go` | [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) | 2 |
| **`driving.ErrInitFileSystem` (neu)** | `port/driving/initproject.go` (T2) | **[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)** | **14** |

**Footnote — INIT-006-Carveout**: Spec §625 nennt [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung)
strikt „Projektnamen-Validierung". Die etablierte Codebase-
Konvention (`cli.go:217-220`, `add.go:410`) erweitert das auf
Service-Name-Validierung. Init übernimmt diese Konvention; ein
dedizierter LH-Code für Service-Name-Validation bleibt
Cluster-T_close-Sub-Decision.

**Footnote — Backup-Sentinel-Doku-Korrektur (T8-Pflicht)**:
heutiges `cli.go::ExitCode`-Godoc Z. 241-244 labelt
`ErrBackupSuffixExhausted` + `ErrBackupSourceMissing` als
[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)-Klasse, obwohl `isFilesystemError` (Z. 369-370)
sie schon auf Exit 14 routet. Spec §595-619 (INIT-005 „Über-
schreibschutz") spricht NICHT von „filesystem-failure-class";
die Slice-Engineering-Entscheidung shiftet die LH-Code-Klassifi-
kation auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) (technical-FS-failure), um Envelope-
Code und Exit-Code-Klasse zu synchronisieren. T8-Doku-Edit:
`cli.go` Z. 241-244 Godoc-Comment auf neue Klassifikation
nachziehen.

**Switch-Order-Pflicht (Add R6 #11 erblich)**:
`mapErrorToDiagnostic` für init MUSS `ErrInitFileSystem` als
ERSTEN `errors.Is`-case prüfen. Multi-`%w`-Wraps (Go 1.20+) machen
`errors.Is(err, sentinel)` für BEIDE gewrappte Sentinels in der
gleichen Chain true; **T3 wrappt** FS-Errors als
`fmt.Errorf("%s: %w: %w", path, ErrInitFileSystem, rawErr)` —
heutiger Stand ist Single-`%w` (`initproject.go` Z. 819/867/969:
`fmt.Errorf("write %s: %w", plan.Template.Path, err)`), T3
erweitert auf Multi-`%w` analog `addservice_execute.go`. Ohne
FS-first-Order würde ein künftiger fachlicher Sentinel im Multi-
Wrap die FS-Klassifikation ([`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) / Exit-Code 14) auf
einen Exit-Code-10-Fachfehler downgraden.

**Switch-Order verbindlich** (T6-Pin verifiziert die Reihenfolge):

```go
switch {
case errors.Is(err, driving.ErrInitFileSystem):  // 1. FS-first (LH-NFA-REL-003 / 14)
case errors.Is(err, driving.ErrBackupSuffixExhausted), errors.Is(err, driving.ErrBackupSourceMissing):
                                                  // 2. FS-Klasse (LH-NFA-REL-003 / 14)
case errors.Is(err, driving.ErrTemplateConflictsWithFlag):
                                                  // 3. Usage (LH-FA-CLI-006 / 2)
case errors.Is(err, driving.ErrConfirmationRequired),
     errors.Is(err, driving.ErrForceRequiresBackup),
     errors.Is(err, driving.ErrBackupUnsupportedKind):
                                                  // 4. INIT-005 Usage-Klasse (10)
case errors.Is(err, driving.ErrProjectExists), errors.Is(err, driving.ErrFileExists):
                                                  // 5. INIT-004 Marker-Kollision (10)
case errors.Is(err, domain.ErrInvalidProjectName),
     errors.Is(err, domain.ErrInvalidServiceName):
                                                  // 6. INIT-006 Name-Validation (10)
default:
    // 7. Default → LH-FA-CLI-006 als Envelope-Code, Exit-Code
    //    via cli.ExitCode(err) (NICHT automatisch 2 — isUsageError
    //    matched nur ErrTemplateConflictsWithFlag + Cobra-Usage-
    //    Errors). Unbekannte Sentinels landen damit als
    //    LH-FA-CLI-006 / Exit 1, was die korrekte Fallback-
    //    Klassifikation ist.
}
```

T6-Pin-Test mit künstlich konstruiertem
`fmt.Errorf("%w: %w", ErrInitFileSystem, ErrProjectExists)` MUSS
[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) / Exit-14 erzeugen — NICHT [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) / Exit-10.

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
  [`LH-FA-INIT-003`](../../../../spec/lastenheft.md#lh-fa-init-003-projektstruktur-erzeugen)-Indikator-Pfade via `Exists`.
- `planFile`/`fileHasManagedBlock` lesen Templates und existierende
  Files.

Der `RecordingFileSystem` delegiert Reads grundsätzlich an den
underlying-FS und captured sie nicht (siehe `recordingfs.go`
Read-Methods Z. 103-122 — Exists/ReadFile/ReadDir/Lstat — plus
Test `TestRecordingFS_ReadsAlwaysDelegate`). T6-Spy-Check muss
zwischen **Writes** (verboten im Dry-Run, Counter MUSS 0 sein)
und **Reads** (erlaubt, Counter irrelevant) unterscheiden. Pin-
Doku in T6-Test-Kommentar explizit ausweisen.

**Soft-Existing × `--devcontainer`-Kollision**: `.devcontainer/
devcontainer.json` ist der 6. softIndicator (initproject.go
Z. 446). `init myproj --devcontainer --dry-run --json` auf
existierender CWD mit dem File trifft `checkSoftExisting` BEVOR
die `--devcontainer`-aware planning läuft. T0-Outcome: das
`--devcontainer`-Flag ändert die Soft-Detection NICHT;
nur `--force`/`--backup`/`--no-interactive` (T0-Disambigua-
toren) hebeln Detection auf. T6 ergänzt einen Pin-Test für
diese Interaktion (Existing-Project-Fixture mit `.devcontainer/
devcontainer.json` + `--devcontainer --dry-run --json` →
ErrProjectExists wie ohne `--devcontainer`).

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
Form wie add-modify (T0-(g) gilt unverändert).

**Content-identical Edge-Case (Idempotenz-Signalisierung)**:
content-identischer Block-Replace ergibt `CountAdditions = 0` und
`diff.Compute(...)==nil` (keine Hunks). Sub-Decision:
PlannedFile-Eintrag bleibt **sichtbar** im Envelope mit
`{action: "modify", count: 0, hunks: omitted}` — NICHT
suppressed. Begründung: Suppression würde den Mid-Write-Failure-
Trace verkleinern und die Recorder-Capture-Liste lückenhaft
machen; sichtbar-mit-count-0 ist UX-transparent
(`writeInitDiff` rendert dazu das `(no changes)`-Hint aus
T0-(k)). T6 pinnt die Form via Re-init-Pin-Test gegen ein
existierendes Projekt mit identischem managed-Block.

**`actionOverwriteFull`-Content-identical-Subform** (
`--backup`-Pfad mit byte-identischem Body, z. B. Re-init nach
manuellem Restore): Recorder capturet (a) den Backup-PlannedFile
als `{action: "create", count: Lines(Original)}` — Backup ist
immer real, identisch oder nicht; (b) den File-PlannedFile als
`{action: "modify", count: 0, hunks: omitted}`; (c)
`writeInitDiff` rendert das `(no changes)`-Hint analog zum
actionReplaceBlock-identical-Fall. T6-Pin via Re-init mit hand-
restoriertem identischen `.gitignore` o.ä.

### T0-(i) Template-Mode-Preview als V1-Out-of-Scope-Carveout

**Round-2 Finding B-4 strukturelle Lücke**: `--template <name>`
ruft im heutigen Wiring `InitProjectService.Init()` →
`initFromTemplate` Z. 409 → `s.templateInit.Init(ctx, ...)` auf
einer **separaten** `TemplateInitService`-Instanz (siehe
`cmd/uboot/main.go` Z. 120). Die separate Instanz hält ihren
eigenen `fsAdapter` und ist **NICHT** an die per-request fsFactory
des InitProjectService gebunden. Konsequenz: bei
`init --template basic --dry-run --json` würde der TemplateInit
direkt auf die Production-FS schreiben — der Recorder sieht
nichts, der Dry-Run schreibt trotzdem.

**Drei Lösungs-Optionen**:

1. **Composition-Root-Refactor**: TemplateInitService bekommt
   auch eine fsFactory; main.go-Wiring shared eine Factory zwischen
   beiden Services oder lokal pro Init-Request synchronisiert.
   Großer Side-Quest, berührt M3-Slice-Code (templateInit).
2. **TemplateInitRequest.PreviewMode**: TemplateInit-API um
   PreviewMode-Override erweitern; InitProjectService ruft mit
   eigenem PreviewMode weiter. Mittlerer Impact, ändert die
   TemplateInit-Port-Signatur.
3. **V1-Out-of-Scope-Carveout**: `init --template <name>` lehnt
   `--dry-run`/`--diff` ab; CLI emittiert eine
   `ErrTemplateConflictsWithFlag`-Diagnostic ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes),
   Exit 2). Folge-Slice
   `slice-v1-cli-json-dry-run-template-preview` (neu in open/
   anzulegen) löst die Composition-Root-Refactor sauber als
   eigene Tranche.

**Vorschlag (T0-Festlegung)**: **Option 3** — Out-of-Scope für
diesen Slice. **NEUER CLI-Level-Mutex-Check** in `init.go` RunE
(T5): `if flags.Template != "" && (flags.DryRun || flags.Diff)
{ return ErrTemplateConflictsWithFlag }` — gibt den existierenden
Sentinel mit einem zweckmäßigen Message-Wrap zurück. Die
`initproject.go` Z. 360-367-Raises decken nur
`--template + --devcontainer/--force/--backup` ab und sind
unverändert; die NEUE `--template + --dry-run|--diff`-Mutex ist
ein CLI-Layer-Check, weil die Use-Case-Request heute keine
PreviewMode/DryRun-Felder hat (T2 fügt sie hinzu, aber die
Mutex bleibt CLI-seitig für klare Fehler-Lokalisation). Pin-
Form siehe §Aufhebungsbedingung Template-Reject-Pin oben. Doku-
Hint in `cli-json-output.md` §6.4: „`--template` ist in V1 mutex
zu `--dry-run`/`--diff`; siehe Folge-Slice
`slice-v1-cli-json-dry-run-template-preview` (geplant) für die
Composition-Root-Refactor-Variante."

Template-Failures (`ErrTemplateNotFound`/`ErrTemplateRender`/
`ErrTemplateCatalog`) bleiben mit ihrer heutigen LH-Klassifikation
([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes), Exit 2 für Conflicts; [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) Exit 14
für Catalog/Render via existierender `isFilesystemError`).

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
fallback-Form mit deterministischem cwd-basename. CLI-Adapter
(`init.go` Z. 117 nutzt `cobra.MaximumNArgs(1)` — positional ist
OPTIONAL, fallback ist reachable). Test-Setup: `WithGetwd` stellt
einen festen Pfad (`/tmp/test-deterministic-projname`) ein;
`resolveProjectName` (initproject.go Z. 542) leitet daraus
`filepath.Base(req.BaseDir)` ab. T6-Pin verifiziert (a) in
`u-boot.yaml.NewContent` den String `project:\n  name:
test-deterministic-projname` UND (b) dass plannedFiles[].Path
KEINEN Pfad-Prefix `test-deterministic-projname/` trägt (=
Path-Anchor-Konsistenz mit T0-(k) Option 1). Doku-Hint in
`cli-json-output.md` §6.4.

### T0-(m) Flag-Matrix-Coverage im Aufhebungsbedingung-Pin

Init hat 5 verhaltens-modifizierende Flags neben `--dry-run`/
`--diff`/`--json`: `--devcontainer`, `--template <name>`,
`--no-git`, `--force`, `--backup`, plus `--allow-external-feature-
sources`. Jede ändert den plannedFiles[]-Shape unterschiedlich.

**Flag-Matrix-Pin-Plan** (T0-Festlegung):

| Flag-Set | Default-Pin | Pflicht in T6? |
| --- | --- | --- |
| Default (kein Flag) | §Aufhebungsbedingung-JSON oben | ja (4 JSON-Kombos × default = 4 Tests) |
| `--devcontainer` | + `.devcontainer/` Dir + 2 Files (devcontainer.json, Dockerfile) | ja (~2 Tests, default + diff-Pfad) |
| `--template basic` | **Out-of-Scope V1** — rejects mit `ErrTemplateConflictsWithFlag` (siehe T0-(i)) | ja (1 negativer Test) |
| `--no-git` | im Recorder-Output identisch zu Default (initGit läuft POST-Capture); im Apply-Pfad suppresst `git init`-Side-Effect, unsichtbar im Envelope | ja (1 Doku-Hint-Test) |
| `--force` (Re-init) | `action: "modify"` für `compose.yaml`/`u-boot.yaml` + `actionReplaceBlock` Pfad | ja (~2 Tests, T0-(h) Pin) |
| `--backup` (Re-init) | + Backup-Files als `action: "create"` (Original-Lines als count) | ja (~2 Tests, Backup-Pfad-Pin) |
| `--allow-external-feature-sources X` | u-boot.yaml content-Variation; plannedFiles[]-Liste identisch zu Default | optional (1 Test als Acceptance-Cluster) |

**T6-Math erklärt**: ~4 (Default × 4 JSON-Kombos) +
~2 (--devcontainer) + 1 (--template-Reject) + 1 (--no-git-Doku-
Pin) + 2 (--force) + 2 (--backup) + 1 (Soft-Existing-Pin) =
**~13 Top-Level-Tests** (statt naiv 6 × 4 = 24). Begründung:
nicht jede Flag-Variation braucht die volle 4-Kombo-Matrix —
Variationen wie `--no-git` (Recorder-irrelevant) und
`--template` (Out-of-Scope-Reject) brauchen nur einen Pin;
`--devcontainer`/`--force`/`--backup` brauchen Default- plus
Diff-Pfad. Plus Mid-Write-Failure-Pin (2 Positionen) +
Concurrent-Init-Mutex-Pin + 3-Flag-Combo +
Path-Anchor-Pin = **~17 Tests total** im T6-Block.

**recordImplicitMkdir-Duplikations-Hazard für `--devcontainer`**:
Round-2 Finding B-4 zeigt: ein expliziter `MkdirAll('.devcontainer')`
wird vom Recorder als `actionCreate` capturet, aber das Dir
existiert danach NICHT auf der underlying-FS (Passthrough=false
im Dry-Run). Wenn dann `WriteFile('.devcontainer/devcontainer.json')`
folgt, triggert `recordImplicitMkdir` einen ZWEITEN synthetischen
`.devcontainer/`-Record. T1-Refactor MUSS `recordImplicitMkdir`
um einen Dedup-Check ergänzen: skip wenn das Dir bereits in
`r.records` als `actionCreate` mit demselben Path steht (oder
einen kleinen In-Memory-Dir-Existence-Overlay-Set führen, der
auf MkdirAll/Mkdir/recordImplicitMkdir-Calls updated wird).
Bonus: dieser Fix gilt auch für add — Add-Code-Mutation in T1
zieht das gleichzeitig mit. **Dedup-Datenstruktur** (T1-Pin):
`map[string]bool` als private `r.knownDirs`-Field der
`RecordingFileSystem`-Struct, initialisiert in `New()` und
gefüllt aus `recordDir` (Mkdir/MkdirAll) und
`recordImplicitMkdir` selbst. `RemoveAll` löscht den Eintrag;
`Rename` schiebt um. T1 ergänzt `TestRecordingFS_ImplicitMkdir
Deduplicates` als Pin.

### T0-(n) `initGit`-Skip im Dry-Run (per-request-FS-Swap reicht nicht aus)

**Round-3 Adversarial-Finding C-1**: `InitProjectService.Init()`
ruft `s.initGit(ctx, ...)` (initproject.go Z. 311-315), das den
separaten `driven.GitClient`-Port nutzt — NICHT den per-request-
swappable `s.fs`. Konsequenz: `u-boot init myproj --dry-run
--json` shells out `git init` auf die echte Festplatte und legt
`.git/` an, obwohl `--dry-run` Null-FS-Mutationen verspricht.
Der Negative-Pin aus §Aufhebungsbedingung (Spy auf fs.*-Calls)
sieht das nicht — git-Operationen laufen am Recorder vorbei.

**Drei Optionen**:

1. **Direct-Skip in Init()**: `if req.PreviewMode == PreviewDryRun
   { skip s.initGit }`. Einfach; PreviewAndApply (`--diff` ohne
   `--dry-run`) lässt `git init` laufen.
2. **noopGit-Adapter im fsFactory**: Composition-Root liefert
   pro PreviewMode auch einen Git-Adapter (production oder noop).
   Symmetrisch zu fsFactory, aber Factory-Signatur wird 3-Tuple.
3. **gitFactory analog fsFactory**: separate Factory für git.
   Cleaner aber 2 Factories statt einer.

**Vorschlag (T0-Festlegung)**: **Option 1** — Direct-Skip mit
Wahrheitstabelle `PreviewDryRun → skip; PreviewAndApply → run;
PreviewNone → run`. Kein Architektur-Side-Effect, minimaler
Code-Eingriff in T3. T6 ergänzt einen `--dry-run --json`-Test
in non-git CWD und verifiziert via `os.Stat(".git/")` dass KEIN
.git/-Dir entstanden ist plus via Spy auf `s.git` dass
`Init`-Calls-Counter == 0 ist.

### T0-(o) JSON-Mode-ProgressPort-Silencing

**Round-3 Adversarial-Finding C-2**: init schreibt im
`emitSummary`-Pfad (initproject.go Z. 283)
`progress.AffectedFiles`-Events DURING der Use-Case-Call auf
stdout — über den `driven.Progress`-Port, der im Composition-Root
mit `progress.NewText(stdout)` verdrahtet ist (main.go Z. 64).
In `--json`-Mode landet das auf stdout VOR dem JSON-Envelope und
JSON-Parser-Konsumenten brechen (zwei JSON-Objects, oder Text-
Prefix vor JSON). add hat KEINEN Progress-Port — init ist der
erste modifying-Service mit stdout-bound Port aus dem Use-Case.

**Vorschlag (T0-Festlegung)**: **Composition-Root-Swap analog
fsFactory**: das CLI-RunE detected `flags.JSON` und konstruiert
einen `progress.Noop`-Adapter. Aktivierung über
`req.SilenceProgress bool`-Field auf `InitProjectRequest` (T2)
oder über einen Setter (`s.SetProgress(noop)`) im Init()-
Wrapper. Pin-Form (T6): `init myproj --json` und `init myproj
--dry-run --json` produzieren stdout das exakt EINE JSON-Object
enthält (`json.Decode → io.EOF` nach einem Decode). Doku-Hint
in `cli-json-output.md` §7: „Modifying-Subcommands mit stdout-
bound ProgressPorts MÜSSEN den Port in JSON-Mode silencen — der
Recorder schützt nur die FS-Layer, nicht stdout."

### T0-(p) Context-Cancellation-Handling (Status-quo-Carveout)

**Round-3 Adversarial-Finding C-5**: Init und Add haben heute
keine `ctx.Err()`/`select-on-Done`-Checks. Eine
`context.Canceled` mid-init produziert keinen Spec-konformen
Envelope — sie fällt durch zur default-Klausel
([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) / Exit 2).

**Vorschlag (T0-Festlegung)**: **Status-quo Out-of-Scope** für
V1. Context-Cancellation ist Cross-Cutting-Concern für ALLE
modifying-Subcommands; ein konsistenter Exit-Code-130-Convention-
Slice wäre eigener Cluster-T_close-Block. Init-Slice ändert
heutigen Pfad NICHT — Doku-Hint in `cli-json-output.md` §7:
„Ctrl-C / context-cancellation während eines modifying-Sub-
commands fällt heute auf [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) / Exit 2; eine
Interrupt-aware Exit-130-Convention bleibt eigener Folge-Slice."

### T0-(q) Planning-Phase-Failures (Force/Validation vor Recorder)

**Round-3 Adversarial-Finding C-3**: Init's `planTemplatedFiles`
(Z. 273) durchläuft `planFile` (Z. 642+) PRO Skeleton-File und
kann ErrForceRequiresBackup / ErrBackupUnsupportedKind / Template-
Read-Failures werfen BEVOR irgendein Recorder-Capture stattfindet.
Mid-Write-Failure-Pin aus §Aufhebungsbedingung deckt nur
Execute-Phase-Failures ab. Planning-Phase-Failures haben
`plannedFiles[] == []` aber `status: error` mit fachlichem Code.

**Vorschlag (T0-Festlegung)**: explizite Pin-Form in
§Aufhebungsbedingung (Planning-Phase-Force-Failure-Pin oben
ergänzt). T6-Acceptance ergänzt einen dedizierten Test für
`init --force --dry-run --json` auf CWD mit unmanaged
`.gitignore` — Pin: `plannedFiles: [], diagnostics:
[[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)], exitCode: 10`. Die Unterscheidung zu Mid-
Write-Failure (`exitCode: 14`) ist load-bearing — Planning-
Errors sind User-Action-Klasse (Exit 10), Write-Errors sind
FS-Klasse (Exit 14).

## T0-Outcomes

Verbindliche Festzurrung wandert nach `next/`-Übergang in diesen
Block — analog add-Slice. Erwartete Form: Tabelle mit den **17
Sub-Decisions (a)-(q)** plus Implementations-Pflicht-Spalte
(T1-T7; T8 ist reine Doku-Schließung). Nach Review-Round-3 vier
neue Sub-Decisions ergänzt: T0-(n) initGit-Skip, T0-(o)
ProgressPort-Silencing, T0-(p) Context-Cancellation-Carveout,
T0-(q) Planning-Phase-Failures.

## Tranchen (vorgeschlagen)

LOC-Schätzung schlanker als add (das alles neu erfunden hat);
init erbt das Pattern, aber die Mutations-Matrix ist breiter und
die Acceptance-Test-Matrix größer.

| T | Inhalt | LOC | Voraussetzung |
| - | ------ | --- | --- |
| T0 | Discovery + 13 Sub-Decisions klären; Pattern-Erbe-Tabelle pinnen | — (Plan) | — |
| T1 | **Refactor-Tranche**: (a) `previewModeFromFlags` extrahieren nach `cli/previewmode.go` (T0-(b); folgt etabliertem Pattern aus `jsonenvelope.go`/`statusview.go`/`jsonallowlist.go`) + Test-File-Move `add_internal_test.go → previewmode_internal_test.go`; (b) `driving.PreviewMode` umbenennen + `type AddPreviewMode = PreviewMode`-Alias (T0-(c)); (c) `mapResponseToWire`/`computeChangeCountAndHunks` → `mapPlannedFilesToWire` umbenennen + nach `cli/wireshapes.go` verschieben (T0-(e) partial — die wireshape-Funktionen haben heute keine Caller-Variation-Frage, weil sie nur Carrier-Mapping ohne `command`-Param machen); (d) Add-Godoc-Wahrheitstabellen in 5 Files (`port/driving/addservice.go`, `recordingfs.go`, `recordingport.go`, `cmd/uboot/main.go`, `add.go`) mitgezogen; (e) **recordImplicitMkdir-Dedup-Fix** mit `r.knownDirs map[string]bool` (Round-2 Finding B-4 — gilt auch für add). Add-Tests bleiben grün durch Type-Alias-`=`-Syntax. LOC-Sub-Breakdown: previewmode.go ~30; Rename+Alias ~15; mapResponseToWire-Migration ~15; Godoc ~50; recordImplicitMkdir-Dedup ~15 = ~125. | ~125 | T0 |
| T2 | **Port-Types + Sentinel**: `driving.InitProjectRequest.PreviewMode`-Field + `driving.InitProjectRequest.SilenceProgress`-Field (T0-(o)), `driving.InitProjectResponse.PlannedFiles`/`Changes`-Felder (analog `AddServiceResponse`), `driving.ErrInitFileSystem`-Sentinel, `cli.isFilesystemError`-Erweiterung. Unit-Test für Sentinel-Identity + ExitCode-Routing. | ~70 | T0 |
| T3 | **Application-Layer**: `InitProjectService.fsFactory`-Feld + `initMu sync.Mutex` + `NewInitProjectServiceWithFactory`-Konstruktor + `Init()`-Wrapper mit Mutex/Swap analog `AddServiceService.Add()`; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; **initGit-Skip-Logic** (T0-(n): `if req.PreviewMode == PreviewDryRun { skip s.initGit }`); **ProgressPort-Swap** (T0-(o): `if req.SilenceProgress { s.progress = noopProgress }` via Setter oder request-time-swap analog s.fs). **Neun FS-Wrap-Stellen** (initproject.go Z. 776 MkdirAll-Loop / Z. 818 WriteFile actionWrite / Z. 866 WriteFile actionReplaceBlock (re-init only) / Z. 968 WriteFile actionOverwriteFull (re-init only); backup.go Z. 88 RemoveAll / Z. 139 CopyExclusive / Z. 149 Mkdir / Z. 198 MkdirAll / Z. 209 Copy) mit `%w: ErrInitFileSystem` umhüllt — Wrap-Form ist Multi-`%w` analog addservice_execute.go (heute Single-`%w`, T3 erweitert). Code-Sites ≠ Runtime-Call-Count, weil Z. 776 in einer Loop läuft. Factory-Tests analog `addservice_factory_test.go` (~200 LOC Test-Datei). | ~320 | T2 |
| T4 | **Composition-Root-Wiring** in `cmd/uboot/main.go`: `initFSFactory`-Closure analog `addFSFactory`. | ~30 | T3 |
| T5 | **CLI-RunE + Helper-Generalisierung**: zwei Sub-Schritte zusammen, weil init's RunE der zweite Caller ist und damit den Helper-Refactor erst real motiviert (`unparam`-Linter-Friendliness statt premature abstraction in T1): (a) **Helper-Generalisierung** `reportAddError`/`writeAddErrorEnvelope`/`writeAddDiff`/`lastPlannedPath` aus add.go extrahieren nach `cli/erroremission.go` als `reportError`/`writeErrorEnvelope`/`writeDiff`/`lastPlannedPath` mit decomposed-Slices-Signatur (T0-(e)); 4 Add-Call-Sites in runAdd migrieren; `mapErrorToDiagnostic → mapAddErrorToDiagnostic` Rename. (b) **init-RunE**: ruft die generischen Helper mit `command="init"` + `mapErr=mapInitErrorToDiagnostic`; **NEUER CLI-Mutex-Check** `--template + --dry-run|--diff → ErrTemplateConflictsWithFlag` (T0-(i)); drei JSON-Pfade analog add; `req.SilenceProgress = flags.JSON` setzen (T0-(o)); Allowlist-Migration (`"u-boot init": true`); Reject-Pin-Test `TestRootJSON_RejectsAllNonMigratedForms` in `internal/adapter/driving/cli/jsonallowlist_test.go` (T0-Outcome verifiziert pre-T5-Count durch lokales `make test`; post-T5 = pre-T5 − 1). | ~280 (Helper-Generalisierung ~120 + init-RunE ~160) | T1 + T2 (T4 für Run-time-Smoke aber Code-parallelisierbar) |
| T6 | **Acceptance-Tests**: ~13 Flag-Matrix-Tests (T0-(m)); plus Soft-Existing-Pin (3 Disambiguatoren) + Soft-Existing × `--devcontainer` (T0-(g)); Planning-Phase-Force-Failure-Pin (T0-(q), exitCode 10); Mid-Write-Failure-Pin (zwei Positionen, T0-(f) Switch-Order-Pin mit Multi-`%w`-Konstrukt, exitCode 14); Template-Reject-Pin (T0-(i), exitCode 2); 3-Flag-Combo `--dry-run --diff --json`; Concurrent-Init-Mutex-Pin (zwei Goroutinen auf ein InitProjectService-Instance, unterschiedliche TempDirs); Path-Anchor-Pin (`PlannedFile.Path` ist project-relativ); **initGit-Skip-Pin** (T0-(n): `--dry-run --json` in non-git CWD → kein .git/-Dir + Spy-Counter 0); **JSON-stdout-Cleanliness-Pin** (T0-(o): `json.Decode → io.EOF`). Test-Fixture-Helper `initFixture(t, opts)` für TempDir + ExistingProject-Setup (shared, ~50 LOC) — per-Test-Body ~25 LOC. ~17 Tests + Mid-Failure-Helper-Cluster + Helper-File = ~600 LOC realistisch. | ~600 | T5 |
| T7 | **Review-Fix-Rounds** (~1-2 Runden bei Pattern-Erbe; add hatte R6/R7/R8): Diff aus Reviewer-Findings konsolidieren, Fixes als eigene Sub-Commits, DoD-Hash-Tabelle ergänzen. | ~80 | T6 |
| T8 | **Closure**: CHANGELOG-Eintrag, `cli-json-output.md` §6-Tabelle (init→done) + §6.1-Reject-Liste (init raus) + §6.4 neue init-Sektion (inkl. Context-Cancellation-Carveout T0-(p) und ProgressPort-Silencing-Hint T0-(o)) + §7 Mutations-Matrix (init-Zeile); `cli.go` Z. 241-244 Godoc-Korrektur (Backup-Sentinels auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) nachziehen, T0-(f) Footnote); roadmap-Update (3/9 done); **[`slice-v1-cli-cleanup-add-preview-mode-alias`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md) als open/-Stub anlegen** (Carveout-Plan-Pflicht T0-(c)); Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz: ~1480 LOC (unchanged trotz T1-E-Verschiebung — die
Helper-Generalisierung wandert nur aus T1 nach T5, Total bleibt
gleich). Pattern-Erbe deckt nur noch ~7 % gegenüber add (~1380),
weil init's Adversarial-Findings (T0-(n)/(o)/(p)/(q)) plus
init-spezifische Eigenheiten (initGit-Port, ProgressPort,
Planning-Phase-Failures, Soft-Existing-Detection) den Pattern-
Spar-Effekt weitgehend aufzehren.

**Helper-Refactor-Sequenzierungs-Entscheidung (post-T1-E-Revert)**:
Variante B aus T1-E-Diskussion: Helper-Generalisierung wandert
nach T5 (statt T1-E), weil ein zweiter Caller (init's RunE)
zeitnah folgt — keine `unparam`-Lint-Suppression nötig, kein
premature-abstraction-Vorwurf. Add bleibt während T2-T4 auf
seinen add-spezifischen Helpern (`reportAddError` etc.); T5
extrahiert + migriert beide gleichzeitig. T1-Restumfang reduziert
auf ~125 LOC (vorher ~250); T5-Umfang auf ~280 LOC (vorher ~160).

**Reihenfolge-Pflicht**: T1 und T2 sind code-parallel
(unterschiedliche Files); T3 wartet auf T2 (braucht
InitProjectRequest.PreviewMode + SilenceProgress +
InitProjectResponse-Felder + ErrInitFileSystem); T4 wartet auf
T3; T5 wartet auf T1 + T2 (braucht gemeinsamen Helper +
Sentinel-Mapping; kann parallel zu T4 entwickelt werden, weil
T5's Interface bereits via T3 fixiert ist); T6 wartet auf T5;
T7 wartet auf T6; T8 wartet auf T7.

**DoD-Hash-Snapshot-Policy** (MEMORY.md feedback
[[feedback_done_slice_dod_hash]]): die DoD-Hash-Tabelle nutzt
**Commit-Hashes pro Tranche** (nicht File-Content-Hashes —
entspricht etablierter Praxis in
[`done/slice-v1-cli-json-dry-run-add.md`](slice-v1-cli-json-dry-run-add.md)
und anderen done-Slices). T1-Commits mutieren bis
zu **8 Files** aus dem add-Slice's done-Snapshot
(`port/driving/addservice.go`, `application/addservice.go`,
`recordingfs.go`, `recordingport.go`, `cmd/uboot/main.go`,
`adapter/driving/cli/add.go`, `adapter/driving/cli/add_internal_test.go`,
plus `adapter/driving/cli/recordingfs.go` falls Dedup-Test-
Datei mutiert). Diese Mutationen sind **additiv** auf der add-
Slice's DoD-Hash-Tabelle:

- Init-Slice's T8-Tranchen-Tabelle bekommt für T1 einen
  Commit-Hash plus eine Footnote mit der vollständigen
  Liste der mutierten add-Slice-Files.
- Add-Slice's done/-Datei bleibt unverändert (accepted Slices
  werden nicht umgeschrieben — AGENTS.md §Slice-Disziplin).
- Forward-Pointer: ein neuer Status-Header-Eintrag in init-
  Slice's done-Datei: „T1 migriert add-Slice-Files, post-T1-
  Revisionen siehe T1-Commit-Hash".

## Review-Round-9 (T7)

Eine Review-Runde via Code-Reviewer-Agent gegen den Diff-Range
`ad56550..bab6b13` (T1-A bis T6 + Coverage-Bump). Ergebnis:
sechs Findings, von denen vier echte Bugs im init-Slice waren,
einer eine Cross-Slice-Divergenz aus add (Folge-Slice ausgelagert)
und einer trivialer Dead-Code-Cleanup.

| # | Sev  | Finding                                                  | Adressierung                                                       |
| - | ---  | -------------------------------------------------------- | ------------------------------------------------------------------ |
| 1 | med  | `mapInitErrorToDiagnostic` fehlt `ErrInvalidFeatureSource`-Case → Code/Exit-Klassen-Drift ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) + Exit 10) | R1: Case auf [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features); Test im `AllCases`-Table — `6e5ad01` |
| 2 | med  | `initFromTemplate` ohne `PreviewMode`-Guard im Application-Layer (CLI fängt es, UC asymmetrisch) | R2: `PreviewMode != PreviewNone → ErrTemplateConflictsWithFlag` am UC-Eintritt; Acceptance-Pin — `e897fa7` |
| 3 | low  | `runBackup` Wrap-Strategie (raw FS ↔ typed Sentinel) ohne direkten Application-Test | R3: Zwei Tests via `RunBackupForTest`-Bridge — `e10b57d`            |
| 4 | low  | T0-(k) Path-Anchor für positional `<name>` ungetestet (trailing-slash, dot-slash, abs-path) | R4: Vier-Cases Acceptance-Table mit `validatingInitUseCaseStub` — `ee30c3c` |
| 5 | low  | Add↔Init divergieren bei `ErrBackupSuffixExhausted`-Code (Add: [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) + Exit 14 → inkonsistent; Init: [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) + Exit 14 → konsistent) | Folge-Slice: [`slice-v1-cli-cleanup-add-backup-error-class`](slice-v1-cli-cleanup-add-backup-error-class.md) |
| 6 | info | Init's mapErr-Switch hat `ErrInvalidServiceName`-Case — dead-code (Init hat keinen Service-Arg) | R1: Case entfernt im selben Commit — `6e5ad01`                     |

T7-LOC-Bilanz: 4 R-Commits (~90 LOC) + 1 Folge-Slice-Stub
(~70 LOC Plan-Markdown). Coverage-Gate bleibt grün (91.10% nach R3,
unchanged nach R4).

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
  [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
  §T0-Outcomes — Vorgaben für den Folge-Slice-Block.
- Pattern-Vorbild:
  [`slice-v1-cli-json-dry-run-add`](slice-v1-cli-json-dry-run-add.md)
  — T0-T6 + Review-Rounds 6-8 voll abgeschlossen. Erbschafts-
  Disziplin in §T0-(a) dieses Slices.
- Spec: [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe),
  [`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung), [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)
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
