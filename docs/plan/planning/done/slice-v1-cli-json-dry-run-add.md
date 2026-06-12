# Slice V1: `add --json` / `--dry-run` / `--diff` — Pattern-Vorbild für modifying-Surface

> **Status:** geplant für v0.4.0 — zweiter Folge-Slice des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 2). Etabliert die schwerere Cluster-Infrastruktur,
> die der Doctor-Slice (Platz 1) bewusst ausgelassen hat:
> **`RecordingFileSystem`-driven-Adapter** (Cluster-T0-(b)
> Variante 2 mit Passthrough-Modus), **Diff-Renderer** (Cluster-
> T0-(d) zweigleisig: Unified-Human + strukturierte Hunks im
> JSON), **Composition-Root-Wiring** der driving-Port-
> Instanzen via `fsSelector`-Closure (T0-(e) Option 4) und
> **`AssertFullEnvelope`-Erstnutzung** (im Doctor-Slice T2 nur
> als Stub angelegt). Stub liegt in `next/` (nach fünf Review-
> Runden für die Iteration übernommen); Review-Findings
> aus drei Runden adressiert: Runde 1 H1-L4 (3× HIGH, 6× MEDIUM,
> 4× LOW), Runde 2 H1-M2 (3× HIGH, 2× MEDIUM — fsSelector-Drei-
> Modi, Exit-Code 14 statt 11, Spec-konformer Binary-Fallback,
> Recorder-Capture-Realität bei Mid-Failure, count-Semantik-
> Konsistenz), Runde 3 H1-M1 (2× HIGH, 1× MEDIUM —
> T3-Tranche-Drift gegen T0-(e)-Enum, `driven.RecorderPort`-
> Interface als sauberer Datenpfad, Trailing-Newline-robuste
> count-Formel), Runde 4 H1-M1 (2× HIGH, 1× MEDIUM —
> `FileMutationRecord` mit OldContent/NewContent als Diff-Basis,
> `ErrAddFileSystem`-Sentinel + non-empty Response on Error-Pfad
> + `cli.ExitCode`-Erweiterung, AK auf `fsFactory`-Form
> synchronisiert), Runde 5 H1-M2 (1× HIGH, 2× MEDIUM —
> `driving.PlannedFile` mit `NewContent`/`OldContent` (`json:"-"`)
> als Content-Träger zwischen Layern, `ErrAddFileSystem`-LH-Code
> auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) korrigiert, T0-(e)-Skizze auf `fsFactory`-
> Tuple-Vertrag synchronisiert). **Slice ✅ done** — Cluster-
> Folge-Slice 2/9 (Add) komplett abgeschlossen, alle Tranchen
> und zwei Code-Review-Runden (R6: 15 Findings, R7: 3 Findings)
> adressiert. DoD-Tabelle siehe unten. Slice-Datei wandert nach
> `done/`. Nächster Cluster-Schritt: Folge-Slice 3/9 (init).
>
> **DoD-Tranchen-Hashes** (alle T0-T6 + R6/R7-Findings):
>
> | Tranche / Round | Inhalt | Commit |
> | --- | --- | --- |
> | T0 | T0-Outcomes festgezurrt (12 Sub-Decisions) | `424b3ec` |
> | T1-A | Port-Types (Carrier + Sentinel + RecorderPort) | `fbcbce8` |
> | T1-B | recordingfs-Adapter mit Passthrough-Schalter | `e9ed7d8` |
> | T1-C | AddServiceService.fsFactory + Recorder-Capture + ErrAddFileSystem-Wrap | `195f146` |
> | T1-D | Composition-Root-Wiring (fsFactory-Closure in main.go) | `505f974` |
> | T2 | Pure-Go Diff-Renderer + checkHunks-Helper | `551da9f` |
> | T4 | u-boot add JSON-RunE-Pfad (3 Modi + Allowlist-Migration) | `93babcb` |
> | T5 | Acceptance-Pins Variante A/B + Idempotent + Diff-Struktur | `d300bfb` |
> | R6-Application | s.fs-Race + dep-PreviewMode + Response.ServiceName | `9ac4fd2` |
> | R6-CLI | JSON-Envelope-Vollständigkeit + Renderer-Verträge (Findings #1-#6, #8, #15) | `55ac6ed` |
> | R6-Defense | Diagnostic-Order + splitLines-Compliance (#11, #14) | `c666a4b` |
> | R6-Tests | WithGetwd + 3-Flag-Combo + Exit-Code-Pins (#12, #13, #2) | `a22db99` |
> | R7 | Impliziter MkdirAll im Recorder + count exakt Spec §477 | `4507a53` |
> | T6 | Closure (DoD-Hashes, user-docs, roadmap, CHANGELOG, done/-Move) | dieser Commit |

## Auslöser

Cluster-Slice [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md) §T0-Outcomes (b)+(d)+(e)
schreiben für diesen Slice drei Cluster-übergreifende Etablierungen
fest:

- **T0-(b) `RecordingFileSystem`-Architektur:** der driven-Adapter
  implementiert
  [`driven.FileSystem`](../../../../internal/hexagon/port/driven/filesystem.go)
  und capturet alle **8** Mutations-Methoden (`WriteFile`,
  `WriteFileExclusive`, `Mkdir`, `MkdirAll`, `Rename`, `RemoveAll`,
  `Copy`, `CopyExclusive`) — auch die heute aus `add` nicht
  aufgerufenen, als Drift-Schutz für zukünftige Use-Cases. Passthrough-
  Schalter trennt Dry-Run-Modus (`Passthrough=false`, kein Production-
  FS-Aufruf) von Preview-and-Apply (`Passthrough=true`, capturen plus
  durchreichen).
- **T0-(d) Diff-Renderer:** Human-Mode rendert Unified-Diff als
  String an stdout; JSON-Mode hängt strukturierte Hunks per
  `plannedFiles[].hunks`-Array an jeden Eintrag (Spec-§326-Voll-
  Schema). Beide Renderer teilen denselben LCS-Hunk-Datentyp.
- **T0-(e) Reihenfolge Platz 2:** `add` ist Pattern-Vorbild für
  modifying-Surface; die nachfolgenden vier modifying-Slices
  (`init`, `generate`, `remove`, `config set`) erben das Wiring
  als geschlossenen Outcome-Block.

Spec-Bezug:

- **[`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)** (Dry-Run, V1, Voll-Schema): für dateiver-
  ändernde Befehle muss `--dry-run` die geplanten Änderungen
  ohne FS-Schreiben zeigen. Bei `--dry-run --json` greift das
  Pflicht-Schema mit `dryRun`/`diff`/`plannedFiles`/`changes`
  als Pflichtfelder (§326,
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) §302-447).
- **[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe)** (Diff, V1): `--diff` zeigt Unterschiede
  zwischen aktuellem und geplantem Zustand. Bei `--diff --json`
  gilt das Voll-Schema mit `diff: true`. Kombinierbar mit
  `--dry-run` (Vorschau ohne Schreiben) oder ohne (Vorschau plus
  Schreiben, „Preview-and-Apply"). Spec §451-489.
- **[`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)** (Maschinen-lesbar, V1): wenn `--json` ohne
  `--dry-run`/`--diff` aufgerufen wird, gilt der Minimalkontrakt
  (§1841) — auch `add` muss diese Form tragen, weil `--json` für
  alle 10 Spec-Enum-Subcommands Pflicht ist.

Heute-Stand-Pre-Scan (Cluster-T0-(b) §483-498-Matrix bestätigt;
Review-Finding H1 adressiert — Schleifen-Anzahl und impliziter
MkdirAll explizit dokumentiert):

- [`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go)
  ruft **direkt** `WriteFile` an **zwei Schleifen-Sites**:
  - **Z. 664** (`for _, w := range []*fileWrite{ep.UBootYAML,
    ep.Compose, ep.EnvExample}`): bis zu **3 Slots** pro Add-Run
    (`u-boot.yaml`, `compose.yaml`, `.env.example`).
  - **Z. 674** (`for _, w := range plan.ExtraFiles`): N
    Extra-Files je nach Catalog-Eintrag (z. B. OTel-Service
    bringt `otel-collector-config.yaml`); Anzahl wächst mit
    jedem zukünftigen Add-on und kann **Sub-Dir-Pfade**
    enthalten (heute flach, künftig denkbar `otel/collector.yaml`).
- **Impliziter `MkdirAll`-Effekt** des
  [`driven.FileSystem.WriteFile`-Vertrags](../../../../internal/hexagon/port/driven/filesystem.go)
  (§25-30: „creating parent directories with mode 0o755 as needed");
  Production-Adapter
  [`driven/fs/fs.go`](../../../../internal/adapter/driven/fs/) Z. 57-61
  ruft `os.MkdirAll(filepath.Dir(path), 0o755)` vor `os.WriteFile`.
  Der `RecordingFileSystem` muss diesen impliziten Parent-Dir-
  Anlage-Effekt im Capture **modellieren** — Sub-Decision T0-(b)
  finalisiert die Form (eigene capturete `MkdirAll`-Plan-Einträge
  vor jedem `WriteFile` mit Sub-Dir-Pfad, oder bewusste
  YAGNI-Auslassung mit Doku-Pin).
- **`BackupPath`-Indirektion**: keine. `add` ist im Pre-Scan
  schlanker als `init` (das sowohl direkte `MkdirAll`/`WriteFile`
  als auch `BackupPath`-CopyExclusive/Mkdir/MkdirAll/Copy/
  RemoveAll-Indirektion hat).
- **Der `RecordingFileSystem` deckt trotzdem alle 8 Mutations-
  Methoden ab** — der Drift-Schutz ist Cluster-Pflicht, auch
  wenn `add` heute nur `WriteFile` (+ implizit `MkdirAll`) ruft.

Vorgänger-Slice (Doctor-Platz-1) hat etabliert:

- `cliJSONEnvelope` mit Pointer-Wrapping auf `*bool`/`*[]T`
  (Anti-Drift M1) — Voll-Schema-Felder erscheinen im modifying-Pfad
  auch bei `false`/`[]` ([`jsonenvelope.go`](../../../../internal/adapter/driving/cli/jsonenvelope.go)).
- `newFullEnvelope`-Konstruktor bereits da, **noch nicht** verwendet
  — Erstnutzung in diesem Slice.
- `jsontestutil.AssertFullEnvelope` als Stub angelegt; voll
  funktional, prüft [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) §326 Required-Set, `action`-
  Enum, `count ≥ 0`.
- Root-PersistentFlag `--json` + Reject-Allowlist + Code-Registry-
  Disziplin. `u-boot add` ist heute Reject-Eintrag; dieser Slice
  migriert ihn auf die Migrate-Liste.

## Aufhebungsbedingung

Alle vier Flag-Kombinationen für `u-boot add <service>` liefern
spec-konforme Outputs:

```bash
u-boot add postgres                      # human, schreibt
u-boot add postgres --dry-run            # human Vorschau, kein Write
u-boot add postgres --diff               # human Unified-Diff + Write (Preview-and-Apply)
u-boot add postgres --dry-run --diff     # human Unified-Diff, kein Write
u-boot add postgres --json               # Minimalkontrakt-Envelope, schreibt
u-boot add postgres --dry-run --json     # Voll-Schema-Envelope, kein Write
u-boot add postgres --diff --json        # Voll-Schema-Envelope, Hunks, Write
u-boot add postgres --dry-run --diff --json  # Voll-Schema, Hunks, kein Write
```

`make test` + `make lint` + `make docs-check` grün. **Reject-Eintrag
für `u-boot add` in
`jsonallowlist.go`
entfernt**; statt dessen Allowlist-Eintrag `"u-boot add": true`.

Konkrete Pin-Form für `add --dry-run --json` (Spec §326-Voll-Schema;
Review-Finding H2 adressiert: `add` mutiert immer **mindestens 3
Files**, Action-Werte projekt-state-abhängig). Zwei Pin-Varianten:

**Variante A — frisch-init Projekt** (`compose.yaml` existiert nicht
oder ist leer):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    {"path": "u-boot.yaml", "action": "modify"},
    {"path": "compose.yaml", "action": "create"},
    {"path": ".env.example", "action": "create"}
  ],
  "changes": [
    {"path": "u-boot.yaml", "count": 2},
    {"path": "compose.yaml", "count": 12},
    {"path": ".env.example", "count": 4}
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

**Variante B — existierendes Setup** (`compose.yaml` hat bereits
andere Services):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    {"path": "u-boot.yaml", "action": "modify"},
    {"path": "compose.yaml", "action": "modify"},
    {"path": ".env.example", "action": "modify"}
  ],
  "changes": [
    {"path": "u-boot.yaml", "count": 2},
    {"path": "compose.yaml", "count": 6},
    {"path": ".env.example", "count": 2}
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

Reihenfolge der `plannedFiles[]` spiegelt
[`addservice.go:140-143`](../../../../internal/hexagon/port/driving/addservice.go)
(`Changed`-Reihenfolge: `u-boot.yaml` → `compose.yaml` →
`.env.example`). Wenn der Catalog-Eintrag zusätzliche `ExtraFiles`
trägt (z. B. OTel `otel-collector-config.yaml`), erscheinen sie
nach `.env.example` in der `ExtraFiles`-Reihenfolge der
Catalog-Definition. `count`-Semantik: gemäß T0-(g) (Spec §430+§477
konsistent — siehe T0-Discovery).

Negative-Pin (Cluster-T0-(b) Pflicht): bei `--dry-run` darf der
RecordingFileSystem **null** Production-FS-Mutations-Aufrufe
durchreichen. Acceptance-Test instrumentiert die Production-FS
mit einem Spy auf alle 8 Mutations-Methoden und assertet `Calls
== 0`.

## Akzeptanzkriterien

- ✅ **`RecordingFileSystem`-driven-Adapter** (Cluster T0-(b)
  Variante 2): neuer Sub-Package
  `internal/adapter/driven/recordingfs/` (T0-(a) dieses Slices
  finalisiert die Lokation). Implementiert
  [`driven.FileSystem`](../../../../internal/hexagon/port/driven/filesystem.go);
  delegiert die **4** Read-Methoden (`Exists`, `ReadFile`,
  `ReadDir`, `Lstat`) an die underlying Production-`fs.FileSystem`.
  Für die **8** Mutations-Methoden gilt der Passthrough-Schalter
  (`Passthrough=false`/`true` per Konstruktor-Option, exakter
  Field-Name = Sub-Decision T0-(b)). **Alle 8** Methoden müssen
  capturet werden, auch die heute aus `add` ungenutzten — der
  Recorder ist der Drift-Anker für zukünftige Folge-Slices.
  **Impliziter `MkdirAll`-Effekt** (Review-Finding H1) ist gemäß
  T0-(b)-Outcome modelliert.
- ✅ **Recorder-Carrier-Typ über die Schicht-Grenze** (Review-
  Finding H3): gemäß T0-(i) liegt der Carrier-Type-Definition in
  `internal/hexagon/port/driving/addservice.go` als Public-Types
  `PlannedFile`/`ChangeEntry`/`Hunk`; `AddServiceResponse`
  bekommt zwei neue Felder. `make lint` depguard grün; weder CLI-
  noch driven-Adapter importiert den jeweils anderen.
- ✅ **Composition-Root-Wiring mit `fsFactory`-Closure**
  (T0-(e) Option 4 + T0-(i) `RecorderPort`; Review-Round-4-
  Finding M1 adressiert — vorheriger Wortlaut „zwei App-Struct-
  Felder bzw. zwei driving-Port-Instanzen" war Stub-Annahme aus
  Round 1, durch T0-(e) verworfen):
  [`cmd/uboot/main.go`](../../../../cmd/uboot/main.go) konstruiert
  einen `fsFactory func(driving.AddPreviewMode) (driven.FileSystem,
  driven.RecorderPort)`-Closure und injiziert ihn in den
  `AddServiceService`-Konstruktor. App-Struct in
  [`cli/cli.go`](../../../../internal/adapter/driving/cli/cli.go)
  und `cli.New(...)`-Signatur **bleiben unverändert**; CLI-RunE
  setzt `req.PreviewMode` gemäß T0-(b)-Wahrheitstabelle. Der
  CLI-Adapter importiert weder `driven.FileSystem` noch
  `recordingfs` direkt (Hard-Rule [`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)/[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement),
  geprüft via `make lint` depguard).
- ✅ **Diff-Renderer zweigleisig** (Cluster T0-(d)): Pure-Go
  LCS-Hunk-Algorithmus (~150-200 LOC) oder
  `github.com/pmezard/go-difflib` (0-Dep, MIT) — Dep-Sub-Decision
  T0-(d) dieses Slices. Beide Renderer arbeiten auf gemeinsamem
  Hunk-Datentyp:
  - **Human-Mode** (`--diff` ohne `--json`): Unified-Diff-String
    mit `+`/`-`-Prefix und `@@ -oldStart,oldLines +newStart,newLines @@`-Header.
  - **JSON-Mode** (`--diff --json`): strukturierte Hunks per
    `plannedFiles[].hunks: [{oldStart, oldLines, newStart, newLines, content}]`
    (Field-Name-Pin T0-(c)).
- ✅ **`u-boot add` JSON-Pfad** in
  [`cli/add.go`](../../../../internal/adapter/driving/cli/add.go):
  drei Code-Pfade je nach Flag-Kombination:
  - `--json` ohne Voll-Schema-Flags → `newMinimalEnvelope`
    (`command: "add"`, Diagnostic-Mapping aus AddServiceResponse).
  - `--dry-run --json` und/oder `--diff --json` → `newFullEnvelope`
    mit `dryRun`/`diff` korrekt gesetzt, `plannedFiles[]` aus dem
    Recorder, `changes[]` per Diff-Hunk-Counter, optional
    `hunks` bei `--diff`.
  - Human-Mode unverändert (existierende Plaintext-Logik bleibt).
- ✅ **`AssertFullEnvelope`-Erstnutzung + Hunks-Helper-
  Erweiterung** (Review-Finding M5 adressiert): Acceptance-Tests
  rufen
  [`jsontestutil.AssertFullEnvelope`](../../../../internal/adapter/driving/cli/jsontestutil/jsontestutil.go)
  mit `WithCommand("add")` plus `WithExpectedCodes(...)` und
  pinnen den Voll-Schema-Required-Set. **Erste Verwendung** des
  Voll-Helpers (Doctor-Slice trug ihn nur als Stub mit Tests).
  `checkPlannedFiles` wird in T2 um `checkHunks` erweitert
  (Hunks-Struktur-Validierung gemäß T0-(l)-Schema-Pin: Pflicht-
  Felder, Zahl-Ranges, Field-Names) — positive Pin (drei valide
  Hunks) und negative Pin (falscher Field-Name `offset` statt
  `oldStart`).
- ✅ **Negative-Pin "null FS-Mutationen im Dry-Run"** (Cluster
  T0-(b) §256-272 Mutations-Matrix-Pflicht): pro `--dry-run`-Pfad
  ein Acceptance-Test, der die Production-FS mit einem Counting-
  Spy umhüllt und assertet, dass nach dem Run alle 8 Mutations-
  Call-Counter exakt 0 sind.
- ✅ **Read-after-Write-Audit** (Cluster T0-(b) §529-547 Pflicht):
  pro Use-Case-Pfad prüfen, ob ein `WriteFile(p, …)` gefolgt von
  `Exists(p)`/`ReadFile(p)`/`Lstat(p)` auf demselben Pfad
  auftritt. Falls ja: **Overlay-Map** in den Recorder ergänzen
  (~30 LOC); falls nein: Audit-Ergebnis im T0-Outcomes
  dokumentieren. Stichprobe T0 zeigt: `add` ist Read-then-Write
  (catalog → service-Files), kein Write-then-Read → kein Overlay
  nötig (Verifikation in T0).
- ✅ **Diff-Renderer-`changes[].count`-Semantik**: gleiches Hunk-
  Datum für Human- und JSON-Modus; `changes[i].count` gemäß T0-(g)
  **`newLines`-Form** (Spec §430+§477 konsistent): bei
  `action: "create"` = total lines der neuen Datei; bei
  `action: "modify"` = Summe der neuen Zeilen über alle Hunks;
  bei `action: "delete"` = `0`. (Vorheriger AK-Wortlaut
  „`oldLines + newLines`" war Stub-Annahme; Review-Round-2-
  Finding M2 korrigiert.)
- ✅ **Allowlist-Migration**: `u-boot add` raus aus dem
  Reject-Pfad in `jsonAllowlist`, rein in den Migrate-Pfad.
  Bestehender Pin-Test
  `TestRootJSON_RejectsAllNonMigratedForms` schrumpft entsprechend
  (10 statt 11 Reject-Cases).
- ✅ **Diagnostic-Codes per LH-Kennung** (Review-Finding M1
  adressiert, T0-(j)): `add`-Diagnostics mappen die sieben
  Sentinels auf `LH-FA-ADD-{001,002,005,006}`/`LH-FA-INIT-{004,
  005,006}` (Tabelle in T0-(j)). **Keine** Erfindung tool-interner
  `add.*`-Codes; `jsontestutil.codeAllowed` lässt LH-Codes via
  `strings.HasPrefix("LH-")` ohne Registry-Pflege durch. Keine
  `DefaultAllowedCodes`-Erweiterung nötig; keine Markdown-
  Sektion-Edit.
- ✅ **Schema-Vertrag-Doku**: `docs/user/cli-json-output.md` §6.1
  (Migrations-Reihenfolge-Tabelle) auf Status "T0 in Arbeit" für
  Platz 2 nachgezogen. Bei T_close auf "done" + Commit-Hash.
- ✅ **Architektur-Grenzen sauber**: `make lint` (depguard) grün;
  CLI-Adapter importiert **kein** `recordingfs` (Wiring via
  driving-Port-Instanzen in `cmd/uboot/main.go`), `recordingfs`
  importiert **kein** `application`/`driving` (driven-Layer-
  Disziplin).

## T0-Discovery (vor `next/`-Übergang)

Sub-Decisions, die dieser Folge-Slice klären muss, bevor er in
`next/` wandert.

### T0-(a) `RecordingFileSystem`-Lokation

Vorschlag: neuer driven-Adapter-Sub-Package
`internal/adapter/driven/recordingfs/` (analog `fs/`, `git/`,
`docker/` etc.). Sub-Decision: Public-API-Form — wrapping
constructor (`NewRecordingFS(underlying driven.FileSystem,
opts ...Option)`) oder zwei separate Konstruktoren
(`NewDryRun(underlying)` / `NewPassthrough(underlying)`)?
Erste Variante (Constructor + Options) ist Repo-Konvention
(`cli.New(...)`, `progress.NewBar(...)`), Vorschlag also Option 1.

### T0-(b) Passthrough-Schalter-Form

`Passthrough bool` als Konstruktor-Option (Default `false` =
Dry-Run-Modus) oder zwei separate Adapter-Konstruktoren? Plus:
exakte Aufruf-Reihenfolge bei `Passthrough=true` für
Mutation-Failure-Vertrag.

Vorschlag (T0-Festlegung):

- **Konstruktor-Option** `WithPassthrough(bool)`.
- **Aufruf-Reihenfolge im Passthrough=true-Pfad:** (1) Plan-
  Eintrag capturen (in `plannedFiles[]` ergänzen), (2)
  Production-Mutation ausführen, (3) bei Mutation-Fehler den
  Plan-Eintrag **bestehen lassen** und ein `diagnostics[]`-Item
  mit `level: "error"` und dem Fehler-Code ergänzen. Begründung:
  der User soll im JSON sehen, was geplant war, **und** wo es
  gescheitert ist; ein Roll-back des Plan-Eintrags würde die
  Drift-Info verschlucken.

**Drei Failure-Scenarios explizit gepinnt** (Review-Findings L2 +
Round 2 H2 + M1 adressiert — Mid-Failure-UX, Exit-Code-Klasse
und Recorder-Capture-Realität):

| Scenario | `plannedFiles[]` | `diagnostics[]` | `status` | `exitCode` |
| --- | --- | --- | --- | --- |
| **Success-Sequenz** (alle 3 Files OK) | alle 3 Files mit korrekter Action | `[]` | `"ok"` | `0` |
| **Mid-Write-Failure** (File 1 OK, File 2 failt) | nur die bis zum Failure tatsächlich gecaptureten Aufrufe (File 1 + File 2) | 1× `level:"error"` mit [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)-konformem FS-Code, `file:` für File 2 | `"error"` | `14` |
| **Pre-Write-Validation-Failure** (z. B. ungültiger Service-Name vor erstem Write) | `[]` (keine Files geplant) | 1× `level:"error"` mit [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung) | `"error"` | `10` |

**Round-2-H2-Korrektur (Exit-Code-Klasse)**: FS-Write-Failure
(Mid-Write-Failure-Scenario) klassifiziert nach [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
als **technischer Persistenz-/Dateisystem-Fehler** → Exit-Code
**14**, **nicht** 11. Code 11 ist für fachliche
Umgebungs-/Prüfungsfehler (`ErrDoctorFailures`,
`ErrDockerUnavailable`). Vorheriger Tabellen-Wert „11" war
Stub-Annahme.

**Round-2-M1-Korrektur (Recorder-Capture-Realität)**: Im Preview-
and-Apply-Modus (`--diff` ohne `--dry-run`) läuft der Recorder
**production-parallel** —
[`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go)
returnt beim zweiten WriteFile-Fehler aus der Schleife und ruft
File 3 nie auf. Der Recorder sieht damit nur die Aufrufe bis zur
Failure-Stelle. Vorheriger Tabellen-Wortlaut „alle 3 Files (auch
ungeschriebene)" war Wunschdenken; der gewählte Capture-
Mechanismus liefert das nicht ohne Pre-Plan-Extractor oder
Wrapper-Use-Case (beides Out-of-Scope für V1).

**Dry-Run-Modus** ist davon **unberührt**: bei `--dry-run` (mit
oder ohne `--diff`) läuft der Recorder mit `Passthrough=false`
und capturet **alle** geplanten Mutations ohne Production-Call.
Hier sieht der User die vollständige Liste.

UX-Hinweis Mid-Write-Failure: das `diagnostics[].file`-Optional-
Feld (Spec §382) markiert die Failure-Stelle. Doku-Hint in
`cli-json-output.md` §6.1 (Add-Sektion): „On mid-write failure
during `--diff` (preview-and-apply), `plannedFiles[]` contains
only the files captured before and at the failure point. The
diagnostic's `file` field identifies the failure point; later
files that the use-case would have written are **not** listed."
Roll-back-aware Capture (alle bereits geschriebenen Files
reverten) ist **Out-of-Scope** für V1 (würde Cluster-T0-(b)
Variante 3 ChangeSet-Pattern erfordern, das explizit verworfen
wurde).

### T0-(c) Diff-Hunk-Field-Name im JSON

Cluster-T0-(d)-Vorschlag: `plannedFiles[].hunks` als Top-Level-
Array pro Datei. Alternative: `plannedFiles[].diff: { hunks: [...] }`
als Sub-Objekt (klarer Namespace für zukünftige Diff-Metadata
wie Encoding-Hints).

Vorschlag (T0-Festlegung): `plannedFiles[].hunks` direkt
(flacher, Spec verlangt kein Sub-Objekt; jsonenvelope.go würde
einen neuen optional-Field-Tag bekommen, Pin-Test ergänzt).

### T0-(d) Diff-Library: Pure-Go intern vs. `pmezard/go-difflib`

`go.mod` ist auf 4 Deps disziplinär minimal. Cluster-Plan
§610-614 nennt `pmezard/go-difflib` als 0-Dep MIT-Lib als Option.
Pure-Go-LCS+Unified-Diff intern wäre ~150-200 LOC.

Vorschlag (T0-Festlegung): **Pure-Go intern**. Begründung:
LCS+Unified-Diff sind klassisches Algorithm-Material, gut
testbar, kein neuer Dep nötig; die Dep-Disziplin (`go.mod`-
4-Dep-Linie) ist langfristiger Wert als die ~150 LOC Ersparnis.
Falls die Pure-Go-Implementation in T2 unerwartete Komplexität
zeigt (z. B. Binary-Diff-Sonderfälle), Sub-Decision-Revert auf
`go-difflib` möglich, aber nicht Default.

### T0-(e) Composition-Root-Doppel-Wiring-Form

Cluster-T0-(b) §441-456 sagt: zwei driving-Port-Instanzen pro
modifying Subcommand (Normal-Mode + Preview-Mode). Sub-Decision:
**App-Struct-Erweiterung-Form**.

Optionen:

1. **Zwei separate App-Struct-Felder**: `addServiceUseCase`
   (Normal) und `addServicePreviewUseCase` (Preview); CLI-RunE
   selektiert. App-Struct wächst um 2 Felder × 5 modifying
   Use-Cases = **10 zusätzliche Felder** im Cluster-Endzustand.
   `cli.New(...)` Signatur wächst von 11 auf 16+ Parameter.
2. **Ein Selector-Field**: `addServiceUseCase func(preview bool)
   driving.AddServiceUseCase`; CLI-RunE ruft `app.addServiceUseCase(preview)`.
3. **Wrapping in App**: ein UseCase-Field plus eine Methode
   `(a *App) addServiceFor(preview bool)`, die intern zwischen
   zwei vorkonstruierten Instanzen auswählt.
4. **`PreviewMode bool` im Request-Type** (Review-Finding M3
   adressiert): `driving.AddServiceRequest{..., PreviewMode bool}`;
   Composition-Root wiret den Use-Case **einmal** mit einem
   FS-Selektor (`fsSelector(preview bool) driven.FileSystem`),
   der intern zwischen Production-FS und RecordingFileSystem
   wählt. CLI-RunE setzt `Request.PreviewMode = a.dryRun || a.diff`.
   **App-Struct unverändert**; `cli.New(...)`-Signatur unverändert.

Vorschlag (T0-Festlegung): **Option 4**. Begründung: die ersten
drei Optionen erzwingen Field-/Symbol-Duplikation in
[`cli.go`](../../../../internal/adapter/driving/cli/cli.go) und
in den 8 Test-Helpers (`newApp`/`newAppWithDoctor`/`newAppWithAdd`/…)
mal 5 modifying Use-Cases, was eine Wartungs-Last über die
gesamte Cluster-Serie hinweg trägt. Option 4 verschiebt die
Wahl in den Application-Layer (`Request`-Feld); die
Composition-Root in `cmd/uboot/main.go` konstruiert den Use-Case
einmal mit einem **Constructor-Closure** als FS-Selector.

**Drei Modi statt zwei** (Review-Round-2-Finding H1 adressiert):
der Selector muss `--dry-run` von `--diff`-ohne-`--dry-run`
unterscheiden — sonst kollabiert Spec §465-468 Preview-and-Apply
zu einem No-Write-Preview. `Request.PreviewMode` ist deshalb
**kein** Boolean, sondern eine Enum:

```go
// internal/hexagon/port/driving/addservice.go (Skizze)
type AddPreviewMode int
const (
    PreviewNone        AddPreviewMode = iota // Normal-Mode: Production-FS direkt
    PreviewDryRun                            // --dry-run: kein Write, vollständiger Plan
    PreviewAndApply                          // --diff ohne --dry-run: Plan capturen + schreiben
)
type AddServiceRequest struct {
    // ...
    PreviewMode AddPreviewMode
}

// cmd/uboot/main.go (Skizze, vollständig in T0-(i) — siehe dort
// für den fsFactory-Vertrag mit Tuple-Return inkl. RecorderPort.
// Diese Skizze hier ist nur die Mode-Switch-Logik; vorheriger
// Stub-Wortlaut mit `selector` und `NewAddServiceServiceWithSelector`
// war Drift gegen T0-(i)/AK/T3 — Review-Round-5-Finding M2
// adressiert: fsFactory ist der einzige Vertrag, der OldContent/
// NewContent-Capture korrekt zurückreicht).
fsFactory := func(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort) {
    switch mode {
    case driving.PreviewDryRun:
        rec := recordingfs.New(prodFS, recordingfs.WithPassthrough(false))
        return rec, rec
    case driving.PreviewAndApply:
        rec := recordingfs.New(prodFS, recordingfs.WithPassthrough(true))
        return rec, rec
    default:
        return prodFS, nil
    }
}
addService := application.NewAddServiceService(fsFactory, ...)
```

CLI-RunE-Mapping: `Request.PreviewMode = previewModeFromFlags(a.dryRun, a.diff)`
mit der Wahrheitstabelle aus T0-(b):

| `--dry-run` | `--diff` | `PreviewMode` | Production-Write? |
| --- | --- | --- | --- |
| nein | nein | `PreviewNone` | ja (Normal-Mode) |
| ja | nein | `PreviewDryRun` | nein (Plan only) |
| nein | ja | `PreviewAndApply` | ja (Plan + Write) |
| ja | ja | `PreviewDryRun` | nein (Diff-Vorschau, kein Write) |

Use-Case-Methode bleibt `Add(ctx, req)`-Signatur
(`req.PreviewMode` wird intern an `selector(mode)`
weitergereicht). depguard bleibt sauber, da weder CLI noch
Use-Case `recordingfs` importieren — nur `cmd/uboot/main.go`
(Wiring-Layer, exempt von `adapter-driving-no-driven` per
`.golangci.yml`-Allowlist).

Verworfene Form-Variante: Option 1 (App-Struct-Doppel-Felder) — der
ursprüngliche Stub-Vorschlag. Trade-off neu bewertet: Wartungs-Last
über 5 Folge-Slices und 8 Test-Helpers wiegt schwerer als die
zusätzliche Sub-Decision für den Application-Layer-Request-Type.

### T0-(f) Mutations-Matrix-Pre-Scan-Dokumentation

Cluster-T0-(b) §256-276 ist Mutations-Matrix-Pflicht. Pre-Scan
für `add` zeigt heute: **nur** `WriteFile` (direkt, 2 Stellen).
Sub-Decision: wo wandert die vollständige Pro-Subcommand-Matrix?

Vorschlag (T0-Festlegung): in den `cli-json-output.md` Doku-Block
§7 (neu) als Tabelle plus Verweis im Recorder-Doc-Comment. Die
Matrix ist der Drift-Anker — wenn ein Folge-Slice einen neuen
Use-Case anlegt, der eine andere Mutations-Methode ruft, muss er
die Matrix erweitern (Pflicht-Eintrag im jeweiligen Slice-AK).

### T0-(g) `changes[i].count`-Semantik

Spec §365-371 fordert `count ≥ 0` als Integer pro `plannedFiles[]`-
Eintrag; der Wortlaut ist semantisch offen. Spec liefert **zwei**
Beispiele für dieselbe `add postgres`-Operation:

- §430-435 (`--dry-run --json`, `action: "create"`): `count: 12`
- §477-482 (`--diff --json` ohne `--dry-run`, `action: "modify"`):
  `count: 6`

Beide Beispiele sind nur konsistent erklärbar, wenn `count =
newLines` interpretiert wird (Review-Finding M2 adressiert —
der vorherige Stub-Vorschlag „`oldLines + newLines`" passt
**nicht** zu beiden Spec-Beispielen, sondern nur zum einen
oder anderen):

- **Bei `action: "create"`**: ganze Datei ist neu, `newLines` =
  totalLines der neuen Datei (Postgres-compose-Service-Block ~12
  Zeilen → Spec §430 `count: 12` ✅).
- **Bei `action: "modify"`**: nur die echten `+`-Zeilen aus den
  Hunks (Implementation `diff.CountAdditions`). Die ursprünglich
  diskutierte „`sum-over-hunks(hunk.newLines)`"-Form schloss auch
  die Context-Lines mit ein und drifted gegen Spec §477 — Round-7
  Finding B-Korrektur (Postgres-Block fügt 6 Zeilen zu existierender
  compose.yaml hinzu → Spec §477 `count: 6` ✅).

**Entscheidung (T0-Festlegung)**: `count = newLines`
(hinzugefügte/totale Zeilen in der neuen Datei-Version). Konkrete
Formel mit **Trailing-Newline-Robustheit** (Review-Round-3-Finding
MEDIUM 1 adressiert — die naive `strings.Split(content, "\n")`-
Form zählt `"a\n"` als 2 Elemente statt 1):

- **`action: "create"`**:
  `count = bytes.Count(newContent, []byte("\n"))` falls
  `newContent` mit `"\n"` endet; sonst `+1` für die letzte unter-
  minated Zeile. Go-Idiom:
  ```go
  func countLines(content []byte) int {
      n := bytes.Count(content, []byte("\n"))
      if len(content) > 0 && !bytes.HasSuffix(content, []byte("\n")) {
          n++
      }
      return n
  }
  ```
  Generierte YAML/`.env.example`-Templates enden konventionell mit
  einem trailing newline; `countLines("a\nb\n") == 2`,
  `countLines("a\nb") == 2`, `countLines("") == 0`. Spec §430
  `count: 12` für eine 12-Zeilen-Postgres-Block-Datei → der
  Renderer liefert das korrekt unabhängig von der Trailing-Newline-
  Konvention.
- **`action: "modify"`**: `sum(+-Zeilen über alle Hunks)` —
  Implementation `diff.CountAdditions(hunks)`. **Review-Round-7
  Finding B-Korrektur**: die ursprüngliche Stub-Form
  `sum(hunk.newLines)` zählte auch die Kontextzeilen, was die
  Spec §477 Beispiel-Zahl `count: 6` für das 6-Zeilen-Postgres-
  Block-Append systematisch überschritt (Beispiel: existierende
  4-Zeilen-redis-compose + 6-Zeilen-postgres-Append → ein Hunk
  mit hunk.NewLines=10, davon 6 echte `+`-Zeilen und 4 Context-
  Zeilen; nur die 6 echten Additions sind Spec-§477-konform).
  `diff.CountFromHunks` (Roh-NewLines-Sum) bleibt als exported
  Helper erhalten für Diff-Rendering-Anwendungen, die die volle
  hunk-side-line-count brauchen.
- **`action: "delete"`**: `0` (keine neuen Zeilen) — Implementation
  short-circuit VOR dem `IsBinary`-Check, damit auch Binary-Deletes
  korrekt `0` zurückgeben (Review-Round-6 Finding #8).

Pin-Tests in T5 müssen vier Trailing-Newline-Edge-Cases pinnen:
(1) `"a\n"` → 1, (2) `"a"` → 1, (3) `""` → 0, (4) `"a\nb\n"` → 2.

Verworfene Optionen:

- **Lines-Sum** (`oldLines + newLines`): scheitert an Spec-Beispiel
  §477 (modify-Postgres-Block hätte oldLines=0, newLines=6 →
  Sum=6, passt zufällig; aber §430 Create-Postgres hätte oldLines=0,
  newLines=12 → Sum=12, passt — Sum-Variante würde sich nur bei
  modify-Fällen mit oldLines>0 von newLines unterscheiden, ist
  also nicht verifizierbar gegen die Spec-Beispiele). Verworfen
  zugunsten der eindeutigeren `newLines`-Form.
- **`max(oldLines, newLines)`**: passt zu beiden Spec-Beispielen
  ebenfalls (max(0,12)=12; max(0,6)=6) — aber semantisch
  unklarer („betroffene Zeilen"). `newLines` ist UX-näher
  („gewachsene Zeilen").
- **Anzahl Hunks**: `len(hunks)` pro Datei. Passt zu **keinem**
  Spec-Beispiel (Postgres-Add wäre wohl 1 Hunk → count=1, nicht
  6 oder 12).

Pin-Test in T5 verifiziert die `newLines`-Form gegen einen
deterministischen Test-Fixture (`testdata/add-postgres-compose-
yaml-fresh.yaml`-Setup → `count` muss 12 sein; `testdata/add-
postgres-compose-yaml-existing.yaml`-Setup → `count` muss 6 sein).

### T0-(h) Pre-Scan-Read-after-Write-Stichprobe für `add`

Cluster-T0-(b) §529-537 fordert pro Use-Case eine Read-after-
Write-Re-Validation. Pre-Scan-Befund aus
[`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go):
`add` liest catalog (`ReadFile` über Template-Adapter) und
schreibt service-Files (`WriteFile`) — **Read-then-Write**-Muster,
kein Write-then-Read. Sub-Decision: dokumentiert in T0-Outcomes
„`add` braucht kein Overlay" oder Re-Validation in T2 mit
ausführlicher Tabellen-Auflistung pro `addservice_*.go`-Datei?

Vorschlag (T0-Festlegung): Stichprobe ist verbindlich
(Cluster-T0-(b)-Lieferpflicht), wandert in T0-Outcomes als
Pro-File-Tabelle. Wenn ein Write-then-Read gefunden wird, T2
ergänzt die Overlay-Map; sonst T2-LOC bleibt schlanker.

### T0-(i) Recorder-Carrier-Typ über die Schicht-Grenze

**(Review-Finding H3 adressiert)**

`plannedFile`/`changeEntry` sind heute CLI-Adapter-private Typen
([`jsonenvelope.go:61-71`](../../../../internal/adapter/driving/cli/jsonenvelope.go)).
Der `RecordingFileSystem` lebt im **driven-Adapter** (T0-(a)
Vorschlag `internal/adapter/driven/recordingfs/`). depguard-Regeln
[`.golangci.yml:241-257`](../../../../.golangci.yml)
verbieten beide Direct-Import-Richtungen
(`adapter-driving-no-driven` und `adapter-driven-no-driving`).
Carrier-Typ-Sub-Decision ist deshalb **Pflicht-T0-Schritt**, nicht
T2-Implementierungs-Detail.

Drei Optionen:

1. **`driving.AddServiceResponse`-Erweiterung**: bestehender
   Response-Type bekommt zwei neue Felder `PlannedFiles
   []driving.PlannedFile` und `Changes []driving.ChangeEntry`
   (neutrale Domain-Wire-Types im `port/driving`-Sub-Package).
   Use-Case befüllt sie aus dem Recorder; CLI-Adapter mapped sie
   1:1 in `cliJSONEnvelope.PlannedFiles`/`Changes`. depguard-
   konform: CLI importiert `port/driving` (heute schon), Use-Case
   importiert `port/driven` (heute schon), `recordingfs` gibt
   konkrete Capture-Datenstruktur an Use-Case, Use-Case mapped
   auf `driving.PlannedFile`.
2. **Domain-Wire-Types in `internal/hexagon/domain/`**: neuer
   Sub-Package `domain/fsplan/` mit `PlannedFile`/`ChangeEntry`/
   `Hunk`. Beide Adapter-Schichten dürfen das Domain importieren
   (depguard erlaubt das explizit). Mehr Ceremony als Option 1.
3. **Decorator-Pattern im CLI-Layer**: ein `cli`-interner
   Wrapper um `driving.AddServiceUseCase` kapselt den Recorder.
   Verletzt das depguard-Verbot dennoch nicht, weil der CLI-Wrapper
   den `RecordingFileSystem` über die Composition-Root als
   `driven.FileSystem`-Interface erhält (kein Concrete-Import).
   Aber: Use-Case sieht weder Plan noch Hunks; das widerspricht
   T0-(b) §441-456 ("Use-Case-Code bleibt unverändert").
   Verworfen.

**Entscheidung (T0-Festlegung): Option 1**. Begründung:
- Kleinster Eingriff (zwei neue Felder im bestehenden Response-Type
  statt neues Sub-Package).
- Hält die Schicht-Disziplin sauber (CLI-Adapter mappt
  driving-Layer-Types auf eigene Wire-Types — heute schon das
  Pattern für `templateJSON`).
- Carrier-Type-Definition liegt in `internal/hexagon/port/driving/addservice.go`
  als neue Public-Types:
  ```go
  type PlannedFile struct {
      Path       string `json:"path"`
      Action     string `json:"action"`
      // NewContent/OldContent: CLI-Renderer-internal — JSON-Tag
      // "-" hält die Rohinhalte aus dem Wire-Output (sonst
      // Base64-Drift und Spec-§326-Verletzung). Diff-Renderer
      // im CLI-Adapter (T2) konsumiert sie für LCS-Hunks und
      // changes[].count (Review-Round-5-Finding H1 adressiert).
      NewContent []byte `json:"-"`
      OldContent []byte `json:"-"`
  }
  type ChangeEntry struct {
      Path  string `json:"path"`
      Count int    `json:"count"`
  }
  type Hunk struct {
      OldStart int    `json:"oldStart"`
      OldLines int    `json:"oldLines"`
      NewStart int    `json:"newStart"`
      NewLines int    `json:"newLines"`
      Content  string `json:"content"`
  }
  ```
  Application-Layer mappt `recorder.Captured()` 1:1 in
  `[]PlannedFile` (Field-Rename `Content → NewContent`/`OldContent`,
  keine Content-Verlust). Beim JSON-Marshal werden die zwei
  `json:"-"`-Felder weggelassen — Wire-Form trägt nur die
  Spec-§326-konformen `path`/`action`/`hunks` (plus optional
  weitere Voll-Schema-Felder).

**Datenpfad Recorder → Response — `driven.RecorderPort`-Interface**
(Review-Round-3-Finding H2 adressiert):

`fsFactory` aus T0-(e) gibt **zwei** Interface-Werte zurück, nicht
nur einen FS:

```go
// internal/hexagon/port/driven/recordingport.go (NEU)
type RecorderPort interface {
    // Captured liefert die seit Konstruktor-Zeit gecaptureten
    // Mutations-Aufrufe in Aufruf-Reihenfolge. Production-FS-
    // Adapter implementieren das nicht; nur RecordingFileSystem
    // erfüllt das Interface.
    Captured() []FileMutationRecord
}

// FileMutationRecord trägt sowohl die Mutations-Klassifikation
// (Path/Action) als auch die rohen Inhalts-Snapshots, die der
// CLI-Adapter für Diff-Rendering und Lines-Count braucht
// (Review-Round-4-Finding H1 adressiert).
//
// NewContent: der bei WriteFile/WriteFileExclusive übergebene
// Body — KANN auch bei delete-Actions leer sein (= nil).
// OldContent: vom Recorder VOR der Mutation per ReadFile
// (delegiert an Production-FS bzw. Production-Snapshot)
// erfasst. Bei create-Action (Datei existierte vorher nicht)
// = nil; bei modify/delete = der bisherige Datei-Inhalt.
//
// Beide Felder werden vom CLI-Adapter (T2 Diff-Renderer) für
// LCS-Hunk-Generierung und Lines-Count konsumiert. Use-Case
// selbst sieht sie nicht — er ruft nur Captured() ab und reicht
// die Records (ohne weitere Verarbeitung) in
// AddServiceResponse.PlannedFiles weiter; CLI-Adapter macht das
// Rendering. (Mapping-Aufgabe T4.)
type FileMutationRecord struct {
    Path       string
    Action     string // "create" | "modify" | "delete"
    NewContent []byte
    OldContent []byte
}
```

**Capture-Mechanik im `recordingfs`-Adapter** (T1-Aufgabe): bei
jedem `WriteFile(path, body, mode)`-Call ruft der Recorder VOR
dem Capture einmal `underlying.ReadFile(path)` (für OldContent;
bei Not-Found = `Action: "create"`, OldContent = nil; sonst
`Action: "modify"`, OldContent = read-Result). NewContent =
übergebener `body`. Bei `RemoveAll(path)`: OldContent = letzter
ReadFile-Stand (für leere Verzeichnisse = nil), NewContent = nil,
`Action: "delete"`. Bei `Copy(src, dst, ...)`: NewContent =
`ReadFile(src)`-Snapshot, OldContent = `ReadFile(dst)` falls
existent, sonst nil — Action je nach dst-Existenz `"create"`
oder `"modify"`.

```go
// internal/hexagon/application/addservice.go (Erweiterung)
type AddServiceService struct {
    fsFactory func(driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort)
    // ... bestehende Felder ...
}

func (s *AddServiceService) Add(ctx context.Context, req AddServiceRequest) (AddServiceResponse, error) {
    fs, recorder := s.fsFactory(req.PreviewMode)
    s.fs = fs  // Use-Case-internal swap
    // ... bestehende Add-Logik unverändert ...
    resp := buildResponse(...)
    if recorder != nil {
        resp.PlannedFiles = mapCapture(recorder.Captured())
        // Diff-Hunks und Changes-Count kommen vom Diff-Renderer
        // (T2) im CLI-Adapter NACH Use-Case-Return; Use-Case
        // selbst füllt nur PlannedFiles aus dem Recorder.
    }
    return resp, nil
}
```

```go
// cmd/uboot/main.go (Composition-Root)
prodFS := fs.New(...)
factory := func(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort) {
    switch mode {
    case driving.PreviewDryRun:
        rec := recordingfs.New(prodFS, recordingfs.WithPassthrough(false))
        return rec, rec // recordingfs.RecordingFileSystem implementiert beide
    case driving.PreviewAndApply:
        rec := recordingfs.New(prodFS, recordingfs.WithPassthrough(true))
        return rec, rec
    default:
        return prodFS, nil // PreviewNone: kein Recorder
    }
}
addService := application.NewAddServiceService(factory, ...)
```

**Schicht-Disziplin sauber**: Application sieht nur Interfaces
(`driven.FileSystem` + `driven.RecorderPort`); CLI sieht nur
driving-Types (`AddServiceResponse.PlannedFiles`); `recordingfs`
(driven-Adapter) ist die **einzige** Implementation, die beide
Ports erfüllt; nur `cmd/uboot/main.go` (Wiring-Layer, depguard-
exempt) kennt den konkreten Typ. depguard `adapter-no-application`
+ `application-no-adapter` bleiben grün.

**Aufgaben-Verteilung der Tranchen**:

- T1-A (Port-Types): Carrier-Types `PlannedFile`/`ChangeEntry`/
  `Hunk` in `port/driving/addservice.go`, `ErrAddFileSystem`-
  Sentinel, `RecorderPort`-Interface + `FileMutationRecord` in
  `port/driven/recordingport.go`.
- T1-B (`recordingfs`): `RecordingFileSystem`-Struct implementiert
  `driven.FileSystem` + `driven.RecorderPort`; `Captured()`-Methode
  returnt internen Mutations-Log.
- T1-C (Application): `AddServiceService.fsFactory`-Feld +
  Verkabelung; `AddServiceResponse.PlannedFiles` aus
  `recorder.Captured()` befüllt; FS-Write-Errors gewrappt mit
  `ErrAddFileSystem`.
- T1-D (`fsFactory` in `cmd/uboot/main.go`): das oben skizzierte
  Wiring — ursprünglich als T3 geplant, zur T1-Klammer
  hochgezogen, weil Wiring direkt am Carrier-/Recorder-Vertrag hängt.
- T2 (Diff-Renderer im CLI-Adapter): pro `PlannedFile.Path` Diff
  rendern, `Hunks[]` und `count` befüllen — getrennt von
  Use-Case-Return-Pfad, weil Diff-Rendering CLI-Concern ist.
- T4 (CLI-RunE): konsumiert `AddServiceResponse.PlannedFiles`,
  ruft Diff-Renderer, baut `cliJSONEnvelope`.

### T0-(j) Diagnostic-Code-Quelle für `add`

**(Review-Finding M1 adressiert)**

`internal/hexagon/port/driving/addservice.go:106-138` definiert
**vier** Spec-konforme Sentinels:

- `ErrServiceUnsupported` ([`LH-FA-ADD-002`](../../../../spec/lastenheft.md#lh-fa-add-002--postgresql-hinzufügen) — unbekannter Service)
- `ErrServiceInconsistent` ([`LH-FA-ADD-005`](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) — Catalog-State-Mismatch)
- `ErrDependenciesRequired` ([`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) — fehlende Add-On-Deps)
- `ErrProjectNotInitialized` ([`LH-FA-ADD-001`](../../../../spec/lastenheft.md#lh-fa-add-001--add-on-befehl) — kein u-boot.yaml)

Plus drei Application-Sentinels, die heute beim Add aufkommen
können (siehe [`cli.go`](../../../../internal/adapter/driving/cli/cli.go)
`ExitCode`-Mapping):

- `ErrInvalidServiceName` ([`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung) — Service-Name-
  Validation, geteilt mit init)
- `ErrFileExists`/`ErrProjectExists` ([`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen) —
  Marker-Kollision, fachlich auch bei add möglich)
- `ErrBackupSuffixExhausted` / `ErrBackupSourceMissing`
  ([`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) — Backup-Pfad-Failures)

**Entscheidung (T0-Festlegung): LH-Codes**. Begründung:

- Spec §445 erlaubt **LH-Kennung der verursachenden Anforderung**
  als kanonische Code-Form ("z. B. [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)").
- `jsontestutil.codeAllowed`
  ([`jsontestutil.go:280-291`](../../../../internal/adapter/driving/cli/jsontestutil/jsontestutil.go))
  akzeptiert LH-Codes (`strings.HasPrefix("LH-")`) **ohne**
  Registry-Pflege — keine Doppel-Doku-Last in
  `DefaultAllowedCodes` plus Markdown-Sektion.
- Keine Erfindung tool-interner `add.*`-Codes nötig; die Drift-
  Risiken aus Doctor-Slice T0-(h) Option 3 entfallen.

Mapping-Tabelle (im T0-Outcomes finalisiert, fester Bestandteil
der Mutations-Matrix-Doku aus T0-(f)):

| Sentinel | LH-Code | Level | Exit-Code |
| --- | --- | --- | --- |
| `ErrProjectNotInitialized` | [`LH-FA-ADD-001`](../../../../spec/lastenheft.md#lh-fa-add-001--add-on-befehl) | error | 10 |
| `ErrServiceUnsupported` | [`LH-FA-ADD-002`](../../../../spec/lastenheft.md#lh-fa-add-002--postgresql-hinzufügen) | error | 10 |
| `ErrServiceInconsistent` | [`LH-FA-ADD-005`](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) | error | 10 |
| `ErrDependenciesRequired` | [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) | error | 10 |
| `ErrInvalidServiceName` | [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung) | error | 10 |
| `ErrFileExists` | [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen) | error | 10 |
| `ErrBackupSuffixExhausted` | [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) | error | 14 |
| **`ErrAddFileSystem` (neu)** | **[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)** | **error** | **14** |

**Neuer Sentinel `ErrAddFileSystem`** (Review-Round-4-Finding H2
adressiert; Round-5-Finding M1 LH-Code-Korrektur): heutige
Add-Writes wrappen rohe FS-Fehler ohne spec-konformen Sentinel;
`cli.ExitCode` `isFilesystemError`-Helper mappt nur die bekannten
Sentinels auf 14, rohe `os.WriteFile`-Errors landen im `default:
return 1`-Zweig. Der Mid-Write-Failure-Scenario aus T0-(b)
(`exitCode: 14`) ist damit ohne Sentinel nicht spec-konform
erreichbar.

**LH-Code-Wahl** (Round-5-M1-Korrektur): [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)
([`spec/lastenheft.md`](../../../../spec/lastenheft.md) §1875-1879
„Abbruch bei kritischen Fehlern: bei kritischen Fehlern muss das
Produkt abbrechen und eine klare Fehlermeldung ausgeben") matched
die FS-Write-Failure-Semantik exakt. Der vorherige Round-4-Vorschlag
[`LH-FA-ADD-002`](../../../../spec/lastenheft.md#lh-fa-add-002--postgresql-hinzufügen) war Drift — der Code beschreibt unbekannte
Services, nicht FS-/Persistenzfehler. T1 ergänzt:

```go
// internal/hexagon/port/driving/addservice.go
var ErrAddFileSystem = errors.New("add: filesystem mutation failed")
```

Application-Layer wrapt rohe FS-Errors: `fmt.Errorf("add: write
%s: %w: %w", path, ErrAddFileSystem, rawErr)` (Go 1.20+ Multi-
`%w`). `cli.ExitCode` `isFilesystemError`-Helper bekommt
`ErrAddFileSystem` ergänzt — `errors.Is(err, ErrAddFileSystem)
→ 14`.

**Non-empty Response on Error-Pfad** (Review-Round-4-Finding H2):
heutige `runExecutePlan`-Form returnt `(AddServiceResponse{},
wrappedErr)` bei FS-Failure. Für Spec-§326-Voll-Schema-Output mit
`plannedFiles[]` aus dem Recorder MUSS der Add-Use-Case bei
Mid-Write-Failure eine **non-empty Response** zurückgeben:

```go
// internal/hexagon/application/addservice.go (T1-Skizze)
func (s *AddServiceService) Add(ctx context.Context, req AddServiceRequest) (AddServiceResponse, error) {
    fs, recorder := s.fsFactory(req.PreviewMode)
    s.fs = fs
    // ... bestehende Add-Logik mit FS-Writes ...
    resp := AddServiceResponse{ /* bestehende Felder */ }
    if recorder != nil {
        resp.PlannedFiles = mapCapture(recorder.Captured())
    }
    if writeErr != nil {
        // Non-empty Response: User soll plannedFiles[] (Calls
        // bis Failure) UND den fachlichen Sentinel sehen.
        return resp, fmt.Errorf("add: write %s: %w: %w", failedPath, ErrAddFileSystem, writeErr)
    }
    return resp, nil
}
```

CLI-RunE: bei Error-Return baut der RunE-Pfad den Voll-Schema-
Envelope aus der Response (`PlannedFiles` ist befüllt) und ergänzt
`diagnostics[]`-Eintrag mit `level: "error"`, code: "[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)",
`file: failedPath`, `message`. Exit-Code via `cli.ExitCode(err)
→ 14`. Pin-Test in T5 verifiziert den kompletten Pfad:
FS-Fake-Failure bei zweitem WriteFile → JSON mit File 1 + File 2,
diagnostics[0].code: "[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern)", `diagnostics[0].file:
"compose.yaml"`, `status: "error"`, `exitCode: 14`.

**Success-Pfade** (`AddServiceResponse.Changed != nil`):
`status: "ok"`, `diagnostics: []`. Idempotent-no-op-Form
(`Changed == nil`): `status: "ok"`, `diagnostics: []`,
`plannedFiles: []`, `changes: []` (keine Mutation geplant). Keine
`Info`-Diagnostics zur State-Action (User hat
`AddServiceResponse.State`-Field zur Aufklärung; tool-internes
Echo wäre Drift-Pfad ohne UX-Wert).

**`DefaultAllowedCodes`-Erweiterung**: KEINE nötig (alle Codes
sind LH-Prefix, gehen über `codeAllowed`'s `strings.HasPrefix`-
Pfad ohne Registry-Eintrag).

### T0-(k) `add --json` (Minimalkontrakt) ohne `--dry-run`/`--diff`: Output-Form

**(Review-Finding M4 adressiert)**

Spec §1841 ist bindend: `--json` **ohne** `--dry-run`/`--diff`
trägt **nur** den Minimalkontrakt
(`status`/`command`/`subcommand?`/`diagnostics`/`exitCode`). Spec
§1842 sagt: das Voll-Schema (`plannedFiles`/`changes`/`dryRun`/
`diff`) gilt **nur** für `--dry-run --json` und `--diff --json`.

Heißt für `u-boot add postgres --json` (ohne Dry-Run/Diff): die
Add-Operation schreibt das FS **um**, das JSON-Output trägt aber
**keine** Plan- oder Change-Information. UX-Spannung (User
skripten `--json` für Automation und wollen wissen, welche Files
verändert wurden), aber Spec verbietet die Felder im
Minimalkontrakt.

Drei Optionen:

1. **Spec-streng Minimal**: `add postgres --json` ohne Dry-Run/Diff
   gibt **nur** den Minimalkontrakt aus (`status: "ok"`,
   `diagnostics: []`, `exitCode: 0`). User-Hint im
   `cli-json-output.md`-Doku-Block: „use `--dry-run --json` to
   preview, `--diff --json` to preview-and-apply with FS-Plan".
2. **Tool-internes `mutated`-Feld**: Spec §1839 erlaubt
   zusätzliche Felder (sind nicht aktiv verboten). Envelope-
   Erweiterung um `mutated: ["u-boot.yaml", "compose.yaml",
   ".env.example"]`-Array (nur Pfade, keine Action/Count — das
   wäre Voll-Schema-Aufweichung). `AssertMinimalEnvelope`
   müsste angepasst werden, um `mutated` zuzulassen (Helper-
   Erweiterung).
3. **Bare-`--json` weiter rejecten**: `u-boot add --json` ohne
   Dry-Run/Diff bleibt in der Reject-Liste, bis eine bessere
   UX-Sub-Decision die richtige Form gefunden hat. Pin-Test
   schrumpft Reject-Liste **nicht** auf 10, sondern bleibt bei
   11 für die `add`-Form (mit Aufhebung nur für `--dry-run`/
   `--diff`-Aufrufe).

**Entscheidung (T0-Festlegung): Option 1 (Spec-streng Minimal)**.
Begründung:

- Spec-Konformität ist die V1-Hard-Rule. Spec §1841 sagt klar:
  Voll-Schema-Felder sind im normalen `--json` **nicht** zulässig.
  `AssertMinimalEnvelope` rejected sie aktiv.
- Tool-internes `mutated`-Feld würde Doctor-Slice's
  Single-Source-of-Truth-Disziplin brechen
  (`checkNoFullFields` müsste eine Whitelist neuer Felder
  bekommen — Drift-Anfangspunkt).
- Reject (Option 3) wäre eine User-feindliche Krücke; spec-
  konformes Minimal ist konsistent und der Hint im Doku-Block
  ist leicht zu pflegen.

Doku-Hint in `cli-json-output.md` §6.1 (Add-Sektion-Erweiterung
in T6): „For a list of files that would change, use `--dry-run
--json` (preview) or `--diff --json` (preview-and-apply)."

### T0-(l) Hunks-Schema-Pin + Binary-Content-Detection

**(Review-Findings M5 + L4 adressiert)**

`plannedFile`-Struct
([`jsonenvelope.go:61-64`](../../../../internal/adapter/driving/cli/jsonenvelope.go))
hat heute nur `Path` und `Action`. Add-Slice erweitert um
`Hunks []hunk`. `AssertFullEnvelope`/`checkPlannedFiles`
([`jsontestutil.go:306-329`](../../../../internal/adapter/driving/cli/jsontestutil/jsontestutil.go))
validiert heute nur `path`+`action` — die `hunks`-Struktur
bleibt ungeprüft.

**Hunks-Schema-Pin** (T0-Festlegung, finalisiert in T1):

```go
type hunk struct {
    OldStart int    `json:"oldStart"` // 1-basiert, ≥ 1
    OldLines int    `json:"oldLines"` // ≥ 0
    NewStart int    `json:"newStart"` // 1-basiert, ≥ 1
    NewLines int    `json:"newLines"` // ≥ 0
    Content  string `json:"content"`  // Multi-Line, mit +/-/space-Prefix
}
```

Field-Name-Wahl (Cluster-T0-(c)-Festlegung war `plannedFiles[].hunks`):
**flach unter `plannedFiles[]`**, kein `plannedFiles[].diff`-Sub-
Objekt (Cluster-Plan §552-595). Spec verlangt das nicht; flache
Form spart eine Schicht.

**`AssertFullEnvelope`-Erweiterung** (T2-Pflicht):
`checkPlannedFiles` bekommt einen zusätzlichen Pfad, der
`plannedFiles[i].hunks` (optional) prüft: wenn vorhanden, dann
Array von Objekten mit Pflichtfeldern `oldStart`/`oldLines`/
`newStart`/`newLines`/`content`, Zahlen ≥ 0, `start`-Werte ≥ 1
bei Lines > 0 (sonst irrelevant). Negative-Pin: ungültiges
Field-Name (`offset` statt `oldStart`) bricht den Test.

**Binary-Content-Detection** (Review-Findings L4 + Round 2 H3
adressiert — Spec §354 erlaubt nur `create|modify|delete`):
`add`-Templates sind **heute** reine Text-Files
([`addservice_execute.go:252-261`](../../../../internal/hexagon/application/addservice_execute.go)
Catalog-Map; `embed.FS`-Templates mit `.compose.tmpl`/`.env.tmpl`-
Suffix). Der Diff-Renderer in T2 muss trotzdem einen UTF-8-
Validity-Check vor LCS-Diff-Rendering durchführen.

**Bei Binary-Content** (`!utf8.Valid(newContent)` oder gleicher
Check auf alte Datei) gilt **Spec-konformes Fallback**:

- `plannedFiles[].action` bleibt im Spec-Enum
  (`create`/`modify`/`delete`) — **kein** „binary"-Wert
  (Round-2-Finding H3-Korrektur: der ursprüngliche Stub-Vorschlag
  `action: "binary"` war Spec-widrig und hätte
  `AssertFullEnvelope` und JSON-Konsumenten korrekt zum Scheitern
  gebracht).
- `plannedFiles[].hunks` wird **nicht gesetzt** (Field via
  `omitempty` weggelassen — Diff-Rendering einfach übersprungen).
- `changes[].count` zählt Byte-Längen-Differenz statt Zeilen
  (`abs(len(newContent) - len(oldContent))` bei modify;
  `len(newContent)` bei create; `0` bei delete) — Spec §365-371
  lässt `count` als Integer ≥ 0 offen, Byte-Diff ist eine valide
  Form.
- Optional: `diagnostics[]`-Eintrag mit `level: "warn"` und Code
  [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Binary-Diff-Skip-Hinweis) als User-Hint —
  Sub-Decision T0-(l) Variante: notwendig oder out-of-scope?
  Vorschlag: ohne diagnostics-Eintrag, weil das Information-Level
  ist (Spec §1834 verbietet `level: "info"`), und Warn würde den
  `status` auf `warn` upgraden, was bei einer reinen Binary-
  Detektion semantisch nicht passt.

**Drift-Anker für Folge-Slices**: zukünftige Add-on-Catalog-Erweiterungen
müssen entweder Text-Templates anbieten oder die Binary-Detection
nachziehen (`embed.FS` erlaubt grundsätzlich Binary). T0-Outcomes-
Doku ergänzt ein Pflicht-Pin: jeder Folge-Slice, der `embed.FS`-
Templates hinzufügt, prüft UTF-8-Validity beim Render-Test.

## T0-Outcomes

Verbindliche Festzurrung der 12 Sub-Decisions aus §T0-Discovery
nach fünf Review-Runden (24 Findings adressiert). Tabellen-Form
kondensiert; Begründungen verbleiben in den jeweiligen §T0-
Discovery-Sub-Sektionen (verlinkt). Implementations-Pflicht-Spalte
nennt die Tranche, die das Outcome materialisiert.

| Sub-Decision | Outcome | Implementations-Pflicht |
| --- | --- | --- |
| [T0-(a)](#t0-a-recordingfilesystem-lokation) RecordingFileSystem-Lokation | Sub-Package `internal/adapter/driven/recordingfs/` analog `driven/fs/`/`driven/git/`-Repo-Konvention | T1 |
| [T0-(b)](#t0-b-passthrough-schalter-form) Passthrough + Failure-Scenarios | Konstruktor-Option `WithPassthrough(bool)`; Aufruf-Reihenfolge im Passthrough=true-Pfad: (1) capturen → (2) Production-Mutation → (3) bei Fehler Plan-Eintrag bestehen lassen + `diagnostics[]`-Item mit [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern). **Drei Failure-Scenarios** spec-konform gepinnt: Success (`status:"ok"`/`exitCode:0`), Mid-Write-Failure (Capture bis Failure-Stelle, `status:"error"`/`exitCode:14`), Pre-Write-Validation-Failure (`plannedFiles:[]`, [`LH-FA-INIT-006`](../../../../spec/lastenheft.md#lh-fa-init-006--projektnamen-validierung)/`exitCode:10`). Dry-Run-Modus liefert vollständige Liste, Preview-and-Apply nur Calls bis Failure. | T1 (Recorder-Mechanik) + T4 (CLI-RunE-Error-Pfad) + T5 (Acceptance-Pins) |
| [T0-(c)](#t0-c-diff-hunk-field-name-im-json) Diff-Hunk-Field-Name | `plannedFiles[].hunks` als flaches Top-Level-Array pro Datei (kein `plannedFiles[].diff`-Sub-Objekt) | T2 (Renderer) + T4 (Envelope-Befüllung) |
| [T0-(d)](#t0-d-diff-library-pure-go-intern-vs-pmezardgo-difflib) Diff-Library | **Pure-Go intern** (~150-200 LOC LCS+Unified-Diff); kein neuer Dep (`go.mod`-4-Dep-Disziplin); Sub-Decision-Revert auf `pmezard/go-difflib` möglich, falls Pure-Go-Implementation in T2 unerwartete Komplexität zeigt | T2 |
| [T0-(e)](#t0-e-composition-root-doppel-wiring-form) Composition-Root-Wiring | **Option 4**: `Request.PreviewMode` als Enum `AddPreviewMode {PreviewNone, PreviewDryRun, PreviewAndApply}`; Composition-Root-Closure `fsFactory(mode) (driven.FileSystem, driven.RecorderPort)`-Tuple; App-Struct + `cli.New(...)`-Signatur unverändert | T1-D (ursprünglich T3) |
| [T0-(f)](#t0-f-mutations-matrix-pre-scan-dokumentation) Mutations-Matrix | Pro-Subcommand-Matrix wandert in `docs/user/cli-json-output.md` §7 (neu); Recorder-Doc-Comment verweist auf §7; jeder Folge-Slice ergänzt seine Use-Case-Mutations-Spalte | T6 (Doc-Update) + Recorder-Doc-Comment in T1 |
| [T0-(g)](#t0-g-changesicount-semantik) `changes[].count`-Semantik | **`newLines`-Form** (Spec §430+§477 konsistent): `create` = total lines neue Datei (trailing-newline-robust via `bytes.Count + HasSuffix`); `modify` = `sum(+-Zeilen über alle Hunks)` via `diff.CountAdditions` (Round-7-B-Korrektur: Stub-Form `sum(hunk.newLines)` zählte Context mit, drifted gegen Spec §477 `count: 6`); `delete` = `0` (short-circuit vor IsBinary, Round-6 #8). Vier Edge-Case-Pin-Tests in T5 (`"a\n"→1`, `"a"→1`, `""→0`, `"a\nb\n"→2`). | T2 (Counter im Renderer) + T5 (Pin-Tests) |
| [T0-(h)](#t0-h-pre-scan-read-after-write-stichprobe-für-add) Read-after-Write-Stichprobe | `add` ist Read-then-Write (catalog → service-Files), **kein** Write-then-Read → kein Overlay-Map-Fallback nötig. T0-Outcomes-Tabelle pro `addservice_*.go`-Datei mit Pre-Scan-Ergebnis in T1-Doc-Comment. | T1 (Doc-Comment) |
| [T0-(i)](#t0-i-recorder-carrier-typ-über-die-schicht-grenze) Recorder-Carrier-Typ | **Option 1**: Carrier-Types in `internal/hexagon/port/driving/addservice.go` als Public-Types (`PlannedFile{Path, Action, NewContent json:"-", OldContent json:"-"}`, `ChangeEntry`, `Hunk`). `AddServiceResponse` bekommt `PlannedFiles []PlannedFile` + `Changes []ChangeEntry`. **`driven.RecorderPort`-Interface** (neu, `port/driven/recordingport.go`) mit `Captured() []FileMutationRecord`-Methode (`FileMutationRecord` trägt `Path`/`Action`/`NewContent`/`OldContent`). `recordingfs.RecordingFileSystem` implementiert beide Ports. depguard `adapter-no-application`/`application-no-adapter` bleiben grün. | T1 (Carrier-Types + RecorderPort + Recorder-Implementation + Capture-Mechanik) + T4 (Mapping recorder.Captured → Response in Use-Case) |
| [T0-(j)](#t0-j-diagnostic-code-quelle-für-add) Diagnostic-Code-Quelle | **LH-Kennungen** (keine erfundenen `add.*`-Codes). Mapping-Tabelle pinnt acht Sentinels auf `LH-FA-ADD-{001,002,005,006}`/`LH-FA-INIT-{004,005,006}`/[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern). Success-Pfad: `status:"ok"`/`diagnostics:[]`. Idempotent-no-op: leere `plannedFiles[]`/`changes[]`. **Neuer Sentinel `ErrAddFileSystem`** (in `port/driving/addservice.go`, gemappt auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern), Exit-Code 14). `cli.ExitCode.isFilesystemError` ergänzt um `ErrAddFileSystem`. **Non-empty Response on Error-Pfad**: Use-Case returnt `(Response{PlannedFiles: ...}, wrappedErr)` bei FS-Failure. Keine `DefaultAllowedCodes`-Erweiterung (LH-Prefix-Pfad). | T1 (Sentinel + Wrap) + `cli.ExitCode`-Erweiterung in T4 + Mapping in T4 |
| [T0-(k)](#t0-k-add---json-minimalkontrakt-ohne---dry-run--diff-output-form) `add --json` (Minimal) Output | **Spec-streng Minimalkontrakt**: ohne `--dry-run`/`--diff` keine FS-Plan-Information im JSON; nur `status`/`command`/`diagnostics`/`exitCode`. Doku-Hint in `cli-json-output.md` §6.1 (Add-Sektion): „use `--dry-run --json` to preview, `--diff --json` to preview-and-apply with FS-Plan". | T4 (Code-Pfad-Verzweigung) + T6 (Doku-Hint) |
| [T0-(l)](#t0-l-hunks-schema-pin--binary-content-detection) Hunks-Schema + Binary-Detection | Hunks-Schema: `{OldStart, OldLines, NewStart, NewLines int, Content string}` mit `json:"oldStart"`/etc. Field-Name-Tags und Range-Constraints (`Start ≥ 1` bei `Lines > 0`, sonst egal). **`AssertFullEnvelope`-Erweiterung** in `jsontestutil.checkPlannedFiles` mit `checkHunks`-Helper (positive + negative Pin gegen Field-Name-Drift `offset` statt `oldStart`). **Binary-Content-Detection** (`!utf8.Valid(...)`): `action` bleibt Spec-Enum, `hunks` weggelassen via `omitempty`, `count` = Byte-Diff statt Lines, **kein** Diagnostic-Eintrag. | T1 (Hunk-Type-Definition) + T2 (UTF-8-Check + checkHunks-Helper) |

**Status nach Outcomes**: T0 ✅ festgezurrt. Sub-Decisions sind
verbindlich; jede T1-T6-Tranche referenziert die zugehörigen
Outcomes via Sub-Decision-Anker. Abweichung von einem Outcome
während T1-T6 erfordert Plan-Update (separater Commit), nicht
nur Code-Anpassung — der Plan ist die Schicht-Spec, T1-T6 sind
deren Realisierung.

**`next/`-Übergang erfolgt** (siehe Status-Header oben); nächster
Plan-Schritt ist `in-progress/`-Übergang plus T1-Start
(`recordingfs`-Adapter + `RecorderPort` + Carrier-Types in
`port/driving/addservice.go` + `ErrAddFileSystem`-Sentinel +
Mutations-Matrix-Doc-Comment).

## Tranchen (vorgeschlagen)

T1 wurde während der Realisierung in vier Sub-Tranchen geteilt
(A/B/C/D). Die ursprünglich als T3 gelistete Composition-Root-
Wiring-Tranche wurde zu **T1-D** umetikettiert, weil das Wiring
direkt den Carrier-Type-Vertrag und `RecorderPort` aus T1-A/T1-B
verdrahtet — Trennung wäre ohne reale Konsumenten künstlich.
T2/T4/T5/T6 behalten ihre Nummern (kein Renumbering, um Cross-
Refs aus Reviews stabil zu halten).

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + Sub-Decisions** aus §T0-Discovery klären (zwölf Sub-Decisions: (a) Recorder-Lokation, (b) Passthrough-Schalter + 3 Failure-Scenarios, (c) Hunk-Field-Name, (d) Diff-Library, (e) Composition-Root-Wiring-Form, (f) Mutations-Matrix-Dokumentations-Ort, (g) `count`-Semantik, (h) Read-after-Write-Stichprobe, (i) Recorder-Carrier-Typ über Schicht-Grenze, (j) Diagnostic-Code-Quelle, (k) Minimal-Output für `add --json` ohne Dry-Run/Diff, (l) Hunks-Schema-Pin + Binary-Detection). Entscheidungen mit Begründung in einem `T0-Outcomes`-Block dokumentieren. | — (Plan-Arbeit) |
| T1-A ✅ `fbcbce8` | **Port-Types** (T0-(i) + T0-(j)): Carrier-Types `PlannedFile`/`ChangeEntry`/`Hunk` plus `AddPreviewMode`-Enum + `PreviewMode`-Request-Feld + `ErrAddFileSystem`-Sentinel in `port/driving/addservice.go`. `RecorderPort`-Interface + `FileMutationRecord` in `port/driven/recordingport.go`. depguard-Konformität geprüft. | ~110 |
| T1-B ✅ `e9ed7d8` | **`recordingfs`-driven-Adapter** (T0-(a) + T0-(b)): `RecordingFileSystem`-Struct implementiert `driven.FileSystem` + `driven.RecorderPort`; Konstruktor mit `WithPassthrough`-Option + 4 Read-Delegationen + 8 Mutations-Methoden (alle 8, auch ungenutzte) + impliziter `MkdirAll`-Capture-Modell + Pre-Write-ReadFile-Capture für `OldContent`. Unit-Tests pro Methode (Dry-Run + Passthrough + Failure-Pfad). | ~210 |
| T1-C ✅ `195f146` | **Application-Layer-Anbindung** (T0-(b) Mid-Failure + T0-(j) Error-Wrap): `AddServiceService.fsFactory`-Feld + `NewAddServiceServiceWithFactory`-Konstruktor + `Add()`-Wrapper mit `selectFS(mode)` + `s.fs`-Swap + `recorder.Captured()` → `resp.PlannedFiles`-Mapping auch auf Error-Pfad. `addservice_execute.go`-WriteFile-Stellen wrappen FS-Errors mit `%w: ErrAddFileSystem`. Drei Factory-Tests (DryRun-Capture, Legacy-Backward-Compat, Write-Failure-Wrap). | ~140 |
| T1-D | **Composition-Root-Wiring** (T0-(e) Option 4, ursprünglich T3): `cmd/uboot/main.go` konstruiert einen `fsFactory(mode driving.AddPreviewMode) (driven.FileSystem, driven.RecorderPort)`-Closure (drei Switch-Cases: `PreviewNone` → `prodFS, nil`; `PreviewDryRun` → `recordingfs.New(prodFS, WithPassthrough(false))`; `PreviewAndApply` → `recordingfs.New(prodFS, WithPassthrough(true))`). `addSvc` migriert auf `NewAddServiceServiceWithFactory`. App-Struct + `cli.New(...)`-Signatur **unverändert**; CLI-RunE-Mapping (`previewModeFromFlags`) bleibt T4-Aufgabe. | ~50 |
| T2 | **Diff-Renderer** + **`AssertFullEnvelope`-Hunks-Helper-Erweiterung** (Review-Finding M5). Pure-Go LCS-Hunk + Unified-String-Renderer + UTF-8-Validity-Check vor LCS (T0-(l), Binary-Detection-Fallback). Hunk-Datentyp gemeinsam für beide Modi. **`checkHunks`-Helper** im `jsontestutil`-Package (Struktur-Pflicht-Felder, Field-Names, Zahl-Ranges). **Golden-File-Tests** gegen Spec-Beispiel-Fixtures (`testdata/add-postgres-compose-fresh.golden` und `-existing.golden`, Review-Finding L3). Unit-Tests gegen klassische LCS-Edge-Cases (leere Inputs, identische Inputs, einseitiger Append, Mitten-Modify, Binary-Detection). | ~310 |
| ~~T3~~ | **In T1-D aufgegangen** — Composition-Root-Wiring war zu eng an T1-A/T1-B (Carrier-Types + `RecorderPort`) gekoppelt, um als separater Schritt sinnvolle Test-Inhalte zu liefern. Pin-Tests für die vier Flag-Kombinationen aus dem ursprünglichen T3-Wortlaut wandern nach T5. | — |
| T4 | **`u-boot add` JSON-RunE-Pfad**: drei Code-Pfade je nach Flag-Kombination — (a) `--json` ohne Dry-Run/Diff → `newMinimalEnvelope` (T0-(k), Spec-streng Minimal); (b) `--dry-run --json` → `newFullEnvelope` mit Recorder-Capture, `dryRun=true`/`diff=false`; (c) `--diff --json` (mit oder ohne Dry-Run) → `newFullEnvelope` mit Hunks (gerendert aus `FileMutationRecord.OldContent`/`NewContent`), `diff=true`. **Error-Response-Pfad** (T0-(j) Round-4 H2): bei `ErrAddFileSystem`-Return baut der RunE den Voll-Schema-Envelope aus der **non-empty Response** plus `diagnostics[]`-Eintrag mit `file:` und LH-Code; Exit-Code 14 via bestehendem `cli.ExitCode`-Pfad (mit ergänztem `isFilesystemError`-Sentinel-Check). Diagnostic-Code-Mapping aus T0-(j). Allowlist-Migration: `u-boot add` raus aus Reject, rein in Migrate. Reject-Pin-Test schrumpft (11 → 10). | ~230 |
| T5 | **Acceptance-Tests** für alle vier Flag-Kombinationen + drei Failure-Scenarios (T0-(b)) + Negative-Pin (Null-FS-Mutationen im Dry-Run) + Diff-Output-Pin (Unified-Struktur stimmt) + Pin-Test für `count`-Semantik gegen Test-Fixture (T0-(g) `newLines`-Pin). Erstnutzung von `jsontestutil.AssertFullEnvelope` mit `checkHunks`-Erweiterung aus T2. Two-Pin-Form (Variante A frisch-init / Variante B existing). | ~360 |
| T6 | **Closure.** CHANGELOG `## [Unreleased]` Added-Eintrag, roadmap.md Cluster-Slice-Zelle aktualisiert (Add done, nächster Schritt `init`), cli-json-output.md §6.1-Tabelle auf done plus Add-Sektion-Erweiterung mit Minimal-vs-Voll-Output-Hinweis (T0-(k)) und Mid-Failure-UX-Hint (T0-(b)). Slice-File `in-progress/` → `done/` mit DoD-Hash-Tranchen-Tabelle. `make docs-check` grün. | — (Doku) |

LOC-Schätzung Folge-Slice: ~1380 LOC nach vier Review-Runden
(Carrier-Types in port/driving, MkdirAll-Capture-Modell, Hunks-
Helper-Erweiterung, Golden-File-Tests, drei Failure-Scenarios,
`driven.RecorderPort`-Interface, Pre-Write-ReadFile-Capture für
OldContent/NewContent, `ErrAddFileSystem`-Sentinel + non-empty
Response on Error-Pfad + cli.ExitCode-Erweiterung) — von
ursprünglich ~1120; weiterhin **deutlich** über der vom Cluster-
Slice gesetzten 200..600-Bandbreite. Begründung: dieser Slice
etabliert die schwerere Cluster-Infrastruktur (Recorder + Carrier-
Types + Diff-Renderer + Composition-Root-Selector-Wiring + Hunks-
Helper), die die vier nachfolgenden modifying-Slices (`init`,
`generate`, `remove`, `config set`) als geschlossenen Outcome-Block
erben. LOC-Bandbreite ist deshalb für diesen Pattern-Vorbild-Slice
**keine** Hard-Rule (analog Doctor-Slice ~630).

**Aufteilungs-Erwägung** (Review-Finding L1 adressiert): ein
separater `slice-v1-cli-json-dry-run-recordingfs`-Sub-Slice (T1
+ Composition-Root-Skeleton, ~420 LOC) wäre architekturell
denkbar. **Verworfen**, weil:

- Der Carrier-Typ-Schicht-Grenze-Vertrag (T0-(i)) verbindet
  Recorder und CLI-Adapter so eng, dass ein dazwischenliegender
  Sub-Slice nur einen "ungenutzten Recorder ohne Konsumenten"
  liefern würde — Test-Inhalte wären synthetisch.
- Composition-Root-Eingriff (`fsSelector`-Closure aus T0-(e)
  Option 4) ist ein Single-Point-of-Change in `cmd/uboot/main.go`;
  ein Split würde denselben Punkt zweimal modifizieren.
- Pattern-Vorbild-Last: die vier Folge-Slices erben das
  geschlossene `add`-Pattern; ein Split würde die Vorbild-Form
  verteilen.

## Out of Scope

- **Envelope-Migration für `template list`**: bleibt
  [`slice-v1-cli-json-dry-run-template`](slice-v1-cli-json-dry-run-template.md)
  ([`open/slice-v1-cli-json-dry-run-template.md`](slice-v1-cli-json-dry-run-template.md))
  vorbehalten (Platz 9).
- **`cliJSONEnvelope.Data`-Feld**: im Doctor-Slice entfernt (Review
  H1); kommt mit `template`-Slice wieder, weil `template list` ein
  Subcommand-spezifisches Payload braucht. `add` benötigt es nicht
  — `plannedFiles[]`/`changes[]` decken die Add-Spezifik ab.
- **Read-after-Write-Overlay**: nur falls T0-(h)-Stichprobe einen
  Write-then-Read-Use-Case findet (heute nicht erwartet). Sonst
  Overlay-LOC nicht in T2.
- **Per-Use-Case-`add`-Refactor**: wenn die Mutations-Matrix für
  `add` über die Zeit wächst (z. B. ein zukünftiger Add-on bringt
  `Copy`-Aufrufe mit), bleibt der Recorder unverändert (alle 8
  Methoden bereits da); die Matrix-Doku wird im jeweiligen Folge-
  Slice ergänzt.
- **Binary-File-Diff**: Unified-Diff ist Text-orientiert. `add`-
  Templates sind heute reine Text-Files (`compose.yaml`,
  `.env.example`). Falls Binary-Templates dazukommen, eigenes
  Slice mit `--diff`-Verhalten-Sub-Decision (Hex-Dump? Hash-
  Vergleich? Skip?).
- **`u-boot add --diff` ohne `--dry-run` ohne `--json`** Render-
  Form: Cluster T0-(d) sagt "Unified-Diff-String an stdout"; die
  exakte Trennzeichen-Form (Color, Pager-Integration) ist
  YAGNI für V1 — Pure-Plaintext mit `+`/`-`-Prefix reicht.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
  §T0-Outcomes (b)+(d)+(e) sind die Vorgaben dieses Slices.
- Vorgänger-Slice:
  [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  T2-T5 etabliert die Infrastruktur (Envelope, Helper, Allowlist,
  Drift-Gates), die dieser Slice **nicht erneut** baut.
- Spec: [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Voll-Schema §326), [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) (Diff
  §451-489), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Minimalkontrakt §1841)
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- ADR: [`ADR-0010`](../../adr/0010-kein-http-driving-adapter.md)
  §Folgepunkte Re-Eval-Trigger 2.
- Code-Anker heute:
  [`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go)
  (WriteFile-Stellen Z. 664+674),
  [`cli/add.go`](../../../../internal/adapter/driving/cli/add.go)
  (RunE-Erweiterungs-Ziel),
  [`driven/fs/`](../../../../internal/adapter/driven/fs/)
  (Production-FS-Vorbild für `recordingfs/`),
  [`cmd/uboot/main.go`](../../../../cmd/uboot/main.go)
  (Composition-Root-Doppel-Wiring-Ziel),
  [`cli/jsonenvelope.go`](../../../../internal/adapter/driving/cli/jsonenvelope.go)
  (`*[]plannedFile`/`*[]changeEntry` bereits da),
  [`cli/jsontestutil/`](../../../../internal/adapter/driving/cli/jsontestutil/)
  (`AssertFullEnvelope`-Erstnutzung),
  `cli/jsonallowlist.go`
  (Allowlist-Migration).
- Vorbild-Slice für T0-Outcomes-Layout:
  [`done/slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  §T0-Outcomes (acht Sub-Decisions).
- Phase: V1 (Teil des V1-pünktlichen
  [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)-Clusters).
