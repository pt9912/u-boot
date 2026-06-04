# Slice V1: `add --json` / `--dry-run` / `--diff` — Pattern-Vorbild für modifying-Surface

> **Status:** geplant für v0.4.0 — zweiter Folge-Slice des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 2). Etabliert die schwerere Cluster-Infrastruktur,
> die der Doctor-Slice (Platz 1) bewusst ausgelassen hat:
> **`RecordingFileSystem`-driven-Adapter** (Cluster-T0-(b)
> Variante 2 mit Passthrough-Modus), **Diff-Renderer** (Cluster-
> T0-(d) zweigleisig: Unified-Human + strukturierte Hunks im
> JSON), **Composition-Root-Doppel-Wiring** der driving-Port-
> Instanzen und **`AssertFullEnvelope`-Erstnutzung** (im
> Doctor-Slice T2 nur als Stub angelegt).

## Auslöser

Cluster-Slice `slice-v1-cli-json-dry-run` §T0-Outcomes (b)+(d)+(e)
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

- **`LH-FA-CLI-007`** (Dry-Run, V1, Voll-Schema): für dateiver-
  ändernde Befehle muss `--dry-run` die geplanten Änderungen
  ohne FS-Schreiben zeigen. Bei `--dry-run --json` greift das
  Pflicht-Schema mit `dryRun`/`diff`/`plannedFiles`/`changes`
  als Pflichtfelder (§326,
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) §302-447).
- **`LH-FA-CLI-008`** (Diff, V1): `--diff` zeigt Unterschiede
  zwischen aktuellem und geplantem Zustand. Bei `--diff --json`
  gilt das Voll-Schema mit `diff: true`. Kombinierbar mit
  `--dry-run` (Vorschau ohne Schreiben) oder ohne (Vorschau plus
  Schreiben, „Preview-and-Apply"). Spec §451-489.
- **`LH-NFA-USE-004`** (Maschinen-lesbar, V1): wenn `--json` ohne
  `--dry-run`/`--diff` aufgerufen wird, gilt der Minimalkontrakt
  (§1841) — auch `add` muss diese Form tragen, weil `--json` für
  alle 10 Spec-Enum-Subcommands Pflicht ist.

Heute-Stand-Pre-Scan (Cluster-T0-(b) §483-498-Matrix bestätigt):

- [`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go)
  ruft **direkt** `WriteFile` an zwei Stellen (Z. 664 + Z. 674),
  **indirekt** via `BackupPath` keine — `add` ist im Pre-Scan
  schlanker als `init` (`init` hat sowohl direkte `MkdirAll`/
  `WriteFile` als auch `BackupPath`-CopyExclusive/Mkdir/MkdirAll/
  Copy/RemoveAll-Indirektion).
- Der `RecordingFileSystem` deckt trotzdem alle 8 ab — der Drift-
  Schutz ist Cluster-Pflicht.

Vorgänger-Slice (Doctor-Platz-1) hat etabliert:

- `cliJSONEnvelope` mit Pointer-Wrapping auf `*bool`/`*[]T`
  (Anti-Drift M1) — Voll-Schema-Felder erscheinen im modifying-Pfad
  auch bei `false`/`[]` ([`jsonenvelope.go`](../../../../internal/adapter/driving/cli/jsonenvelope.go)).
- `newFullEnvelope`-Konstruktor bereits da, **noch nicht** verwendet
  — Erstnutzung in diesem Slice.
- `jsontestutil.AssertFullEnvelope` als Stub angelegt; voll
  funktional, prüft `LH-FA-CLI-007` §326 Required-Set, `action`-
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
[`jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
entfernt**; statt dessen Allowlist-Eintrag `"u-boot add": true`.

Konkrete Pin-Form für `add --dry-run --json` (Spec §326-Voll-Schema):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    {"path": "compose.yaml", "action": "modify"},
    {"path": ".env.example", "action": "create"}
  ],
  "changes": [
    {"path": "compose.yaml", "count": 12},
    {"path": ".env.example", "count": 4}
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

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
- ✅ **Composition-Root-Doppel-Wiring**: das App-Struct in
  [`cli/cli.go`](../../../../internal/adapter/driving/cli/cli.go)
  trägt **zwei** Felder pro modifying Use-Case (Normal + Preview)
  oder eine `addServiceUseCase`-Variante mit selector-Funktion
  (Sub-Decision T0-(e)). [`cmd/uboot/main.go`](../../../../cmd/uboot/main.go)
  konstruiert in der Composition-Root beide driving-Port-Instanzen
  und injiziert sie. Der **CLI-Adapter** importiert weder
  `driven.FileSystem` noch `recordingfs` direkt (Hard-Rule
  `LH-FA-ARCH-002`/`-003`, geprüft via `make lint` depguard).
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
- ✅ **`AssertFullEnvelope`-Erstnutzung**: Acceptance-Tests
  rufen
  [`jsontestutil.AssertFullEnvelope`](../../../../internal/adapter/driving/cli/jsontestutil/jsontestutil.go)
  mit `WithCommand("add")` plus `WithExpectedCodes(...)` und
  pinnen den Voll-Schema-Required-Set. **Erste Verwendung** des
  Voll-Helpers (Doctor-Slice trug ihn nur als Stub mit Tests).
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
  Datum für Human- und JSON-Modus; `changes[i].count` zählt
  geänderte Zeilen (`oldLines + newLines` Summe pro Hunk, dann
  pro Datei aggregiert). Sub-Decision T0-(g) finalisiert die
  Counter-Form.
- ✅ **Allowlist-Migration**: `u-boot add` raus aus dem
  Reject-Pfad in `jsonAllowlist`, rein in den Migrate-Pfad.
  Bestehender Pin-Test
  `TestRootJSON_RejectsAllNonMigratedForms` schrumpft entsprechend
  (10 statt 11 Reject-Cases).
- ✅ **Code-Registry-Erweiterung**: falls `add` neue Diagnostic-
  Codes emittiert (z. B. `add.service-conflict`,
  `add.template-render`), landen sie in
  [`jsontestutil.DefaultAllowedCodes`](../../../../internal/adapter/driving/cli/jsontestutil/coderegistry.go)
  **plus** in der Markdown-Sektion-§5 von
  [`docs/user/cli-json-output.md`](../../../user/cli-json-output.md)
  zwischen den `<!-- code-registry:start/end -->`-Markern. Beide
  Drift-Gates aus dem Doctor-Slice T2 erzwingen die Doppel-
  Pflege automatisch.
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

1. **Zwei separate Felder**: `addServiceUseCase` (Normal) und
   `addServicePreviewUseCase` (Preview); CLI-RunE selektiert.
2. **Ein Selector-Field**: `addServiceUseCase func(preview bool)
   driving.AddServiceUseCase`; CLI-RunE ruft `app.addServiceUseCase(preview)`.
3. **Wrapping in App**: ein UseCase-Field plus eine Methode
   `(a *App) addServiceFor(preview bool)`, die intern zwischen
   zwei vorkonstruierten Instanzen auswählt.

Vorschlag (T0-Festlegung): Option 1 (zwei Felder). Einfachster
Composition-Root, kein neues Pattern, leicht testbar. App-Struct
wird größer; aber das ist Pattern-Vorbild-Last (vier modifying
Folge-Slices machen dasselbe).

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
Eintrag. Vager Wortlaut. Sub-Decision: **was wird gezählt**?

Optionen:

1. **Diff-Lines-Sum**: `oldLines + newLines` pro Hunk, dann pro
   Datei aggregiert. Repräsentiert "Größe des Eingriffs".
2. **Geänderte Zeilen netto**: `max(oldLines, newLines)` pro
   Hunk. Repräsentiert "Anzahl gewordener Zeilen".
3. **Anzahl Hunks**: `len(hunks)` pro Datei. Repräsentiert
   "Anzahl unabhängiger Änderungs-Blöcke".

Vorschlag (T0-Festlegung): Option 1 (Lines-Sum). Lastenheft-
Beispiel §430-435 zeigt `{"path": "compose.yaml", "count": 12}`
für eine `add postgres`-modify mit 12 hinzugefügten Zeilen —
das passt zu Option 1 (Δ-Lines).

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

## Tranchen (vorgeschlagen)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + Sub-Decisions** aus §T0-Discovery klären (acht Sub-Decisions, inkl. Mutations-Matrix und Read-after-Write-Stichprobe). Entscheidungen mit Begründung in einem `T0-Outcomes`-Block dokumentieren. | — (Plan-Arbeit) |
| T1 | **`recordingfs`-driven-Adapter** anlegen. `RecordingFileSystem`-Struct + Konstruktor + 4 Read-Delegationen + 8 Mutations-Methoden (alle 8, auch ungenutzte). Unit-Tests pro Methode: Dry-Run-Mode (kein Production-Call), Passthrough-Mode (capture + delegate), Mutation-Failure-Pfad. depguard-Konformität geprüft (driven-Layer-Disziplin). | ~280 |
| T2 | **Diff-Renderer** (Pure-Go LCS-Hunk + Unified-String-Renderer). Hunk-Datentyp gemeinsam für beide Modi. Unit-Tests gegen klassische LCS-Edge-Cases (leere Inputs, identische Inputs, einseitiger Append, Mitten-Modify, …). | ~220 |
| T3 | **Composition-Root-Doppel-Wiring**: App-Struct um zwei `addService*`-Felder (Normal + Preview), `cmd/uboot/main.go` konstruiert beide. `cli.New(...)` Signatur erweitert (mit Test-Helper-Anpassung). Test, der den Doppel-Pfad pinnt: `--dry-run` ruft Preview-Instanz, ohne Flag Normal-Instanz. | ~140 |
| T4 | **`u-boot add` JSON-RunE-Pfad**: drei Code-Pfade je nach Flag-Kombination (Minimal, Voll-Schema mit Dry-Run, Voll-Schema mit Diff). Allowlist-Migration: `u-boot add` raus aus Reject, rein in Migrate. Reject-Pin-Test schrumpft (11 → 10). | ~180 |
| T5 | **Acceptance-Tests** für alle vier Flag-Kombinationen plus Negative-Pin (Null-FS-Mutationen im Dry-Run) plus Diff-Output-Pin (Unified-Struktur stimmt). Erstnutzung von `jsontestutil.AssertFullEnvelope`. Carveouts-Eintrag (falls neue Diagnostics-Codes auftreten und Code-Registry-Edits nötig sind). | ~300 |
| T6 | **Closure.** CHANGELOG `## [Unreleased]` Added-Eintrag, roadmap.md Cluster-Slice-Zelle aktualisiert (Add done, nächster Schritt `init`), Cluster-Slice §Per-Command-Folge-Slices §6.1-Tabelle in cli-json-output.md auf done. Slice-File `in-progress/` → `done/` mit DoD-Hash-Tranchen-Tabelle. `make docs-check` grün. | — (Doku) |

LOC-Schätzung Folge-Slice: ~1120 LOC — **deutlich** über der vom
Cluster-Slice gesetzten 200..600-Bandbreite. Begründung: dieser
Slice etabliert die schwerere Cluster-Infrastruktur (Recorder +
Diff-Renderer + Composition-Root-Doppel-Wiring), die die vier
nachfolgenden modifying-Slices (`init`, `generate`, `remove`,
`config set`) als geschlossenen Outcome-Block erben. LOC-
Bandbreite ist deshalb für diesen Pattern-Vorbild-Slice **keine**
Hard-Rule (analog Doctor-Slice ~630).

## Out of Scope

- **Envelope-Migration für `template list`**: bleibt
  `slice-v1-cli-json-dry-run-template`
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
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  §T0-Outcomes (b)+(d)+(e) sind die Vorgaben dieses Slices.
- Vorgänger-Slice:
  [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  T2-T5 etabliert die Infrastruktur (Envelope, Helper, Allowlist,
  Drift-Gates), die dieser Slice **nicht erneut** baut.
- Spec: `LH-FA-CLI-007` (Voll-Schema §326), `LH-FA-CLI-008` (Diff
  §451-489), `LH-NFA-USE-004` (Minimalkontrakt §1841)
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
  [`cli/jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  (Allowlist-Migration).
- Vorbild-Slice für T0-Outcomes-Layout:
  [`done/slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  §T0-Outcomes (acht Sub-Decisions).
- Phase: V1 (Teil des V1-pünktlichen
  `slice-v1-cli-json-dry-run`-Clusters).
