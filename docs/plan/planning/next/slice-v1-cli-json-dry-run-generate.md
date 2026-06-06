# Slice V1: `generate --json` / `--dry-run` / `--diff` — Vier-Artefakt-Surface

> **Status:** T0-Discovery + R1/R2/R3/R4/R5 adressiert, `next/` (Lifecycle-Übergang aus `open/` nach fünf Pre-`next/`-Review-Runden; 21 Findings gesamt: 4 HIGH, 14 MED, 3 LOW). Vierter Folge-Slice (4/9) des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 4/9). Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
> 1:1 für den `PreviewMode`-Carrier (kanonisch — kein neuer Service-
> Prefix-Alias, Alias-Lebensdauer-Pflicht aus init-T0-(c)), den
> `RecordingFileSystem`-driven-Adapter, den Pure-Go LCS-Diff-
> Renderer, den `previewModeFromFlags`-Mapper und die generalisierten
> Error-Emission-Helper (`reportError`/`writeErrorEnvelope`/
> `writeDiff`/`lastPlannedPath`).
>
> Generate-spezifisch sind **vier Artefakt-Handler** (changelog /
> readme / env-example / devcontainer) mit unterschiedlichen
> Action-Semantiken (`Created` / `UpdatedBlock` / `NoOp` /
> `RepairedManual`), eine **schmale FS-Surface** (2 von 8
> Recorder-Mutations-Methoden im Default-Pfad: `WriteFile` + bei
> devcontainer `MkdirAll`), die **Devcontainer-Atomicity-Asymmetrie**
> (Phase 1 `planDevcontainerFiles` ist Pre-Write-Validation-atomar,
> Phase 2 `executeDevcontainerPlans` ist **nicht** Roll-back-atomar
> — bewusster Carveout mit `carveouts.md`-Eintrag, siehe T0-(i))
> und der **`--allow-external-feature-sources`-Side-Effect**
> (mutiert `u-boot.yaml` als zusätzlichen Schreib-Vorgang außerhalb
> des Artefakt-Targets).
>
> **Voraussetzungen aus dem Init-Slice**: das Pattern-Vorbild
> (Closure-Form `fsFactory(driving.PreviewMode)`, Multi-`%w`-Wraps
> mit Switch-Order-Pflicht, Helper-Generalisierung, Path-Anchor-
> Disziplin) ist nach init-T0-T8 voll etabliert. Der
> `AddPreviewMode → PreviewMode`-Alias-Cleanup steht im open/-
> Plan `slice-v1-cli-cleanup-add-preview-mode-alias` und wartet
> explizit auf „mindestens einen weiteren Folge-Slice" — generate
> ist genau dieser Folge-Slice, und MUSS deshalb `driving.PreviewMode`
> direkt referenzieren (kein neuer `GeneratePreviewMode`-Alias).

## Auslöser

Cluster-Slice §T0-Outcomes (a)+(b)+(e) machen jeden modifying-
Subcommand für `--json`/`--dry-run`/`--diff` verbindlich
(`LH-NFA-USE-004` §1813, `LH-FA-CLI-007` §326, `LH-FA-CLI-008`
§451-489). `u-boot generate <artifact>` ist nach `doctor`, `add`
und `init` der vierte Subcommand — und der erste, der
**mehrere Artefakte** über einen einzigen Subcommand bedient. Die
JSON-Surface muss pro Artefakt die spec-konforme Vorschau plus
die Action-Klassifikation (`Created`/`UpdatedBlock`/`NoOp`/
`RepairedManual`) tragen, ohne das Subcommand selbst zu zerlegen.

Spec-Bezug (geerbt von init/add):

- `LH-FA-CLI-007` (Dry-Run, Voll-Schema §326)
- `LH-FA-CLI-008` (Diff, §451-489)
- `LH-NFA-USE-004` (Minimalkontrakt §1841)

Generate-spezifische Spec-Anker:

- `LH-FA-GEN-001` (Subcommand-Surface mit unbekanntem Artefakt
  → Exit 2)
- `LH-FA-GEN-002` (changelog / `LH-AK-007` Keep-a-Changelog)
- `LH-FA-GEN-003` (readme)
- `LH-FA-GEN-004` (env-example)
- `LH-FA-DEV-001` (devcontainer)
- `LH-FA-DEV-003` (`--allow-external-feature-sources` Side-Effect)
- `LH-FA-GEN-005` (Idempotenz / NoOp)
- `LH-NFA-REL-003` (FS-Failure-Klasse, geerbt für Mid-Write
  analog init)

Heute-Stand-Pre-Scan
(`internal/hexagon/application/generate.go`, ~960 LOC;
`internal/adapter/driving/cli/generate.go`, ~160 LOC):

| Phase | Methode | Pfade (Default-Artefakt-Pfade) | Code-Anker |
| --- | --- | --- | --- |
| changelog Create/Update | `WriteFile` (Z. 289/344) | `CHANGELOG.md` | direkt |
| readme/env-example (`generateManagedFile`) | `WriteFile` (Z. 504/563) | `README.md` / `.env.example` | direkt |
| devcontainer Multi-File (`writeDevcontainerPlan`) | `MkdirAll` (Z. 848) + `WriteFile` × 2 (Z. 852/864) | `.devcontainer/devcontainer.json` + `.devcontainer/Dockerfile` | direkt; **Phase-1-Pre-Write-Validation atomar** (Plan-Phase capturet keine Schreiboperation; `ErrGenerateManualConflict` returnt vor jedem Write), **Phase-2-Execute-Phase NICHT atomar** (Mid-Write zweiter File → Half-State auf Disk, T0-(i) Carveout) |
| `--allow-external-feature-sources` u-boot.yaml-Mutation | `WriteFile` (Z. 951) | `u-boot.yaml` | direkt — Side-Effect aus Flag |

Damit nutzt generate **2 von 8** Recorder-Mutations-Methoden
direkt (`WriteFile`, `MkdirAll`) plus `--allow-external-feature-
sources` als Side-Effect auf `u-boot.yaml` (auch `WriteFile`).
`WriteFileExclusive`, `Mkdir`, `Rename`, `RemoveAll`, `Copy`,
`CopyExclusive` werden NICHT gerufen — Recorder deckt sie als
Drift-Schutz trotzdem ab (etabliert in add).

Use-Case-Deps: `driven.FileSystem` (Read + Write), `driven.YAMLCodec`,
`driven.Logger`. **KEIN** `GitClient`, **KEIN** `Progress`-Port —
generate ist deutlich schmaler als init (`initGit`-Skip-Carveout
und `ProgressPort`-Silencing-Carveout entfallen).

## Aufhebungsbedingung

Acht Flag-Kombinationen für `u-boot generate <artifact>` liefern
spec-konforme Outputs (geerbt von add/init, ein-zu-eins-Symmetrie),
mal vier Artefakte:

```bash
u-boot generate <artifact>                       # human, schreibt
u-boot generate <artifact> --dry-run             # human Vorschau, kein Write
u-boot generate <artifact> --diff                # human Unified-Diff + Write
u-boot generate <artifact> --dry-run --diff      # human Unified-Diff, kein Write
u-boot generate <artifact> --json                # Minimalkontrakt-Envelope, schreibt
u-boot generate <artifact> --dry-run --json      # Voll-Schema, kein Write
u-boot generate <artifact> --diff --json         # Voll-Schema mit Hunks, schreibt
u-boot generate <artifact> --dry-run --diff --json
```

Für `<artifact>` ∈ {changelog, readme, env-example, devcontainer}.

`make gates` grün (lint + test + coverage-gate ≥ 90 % + docs-check) —
modifying-CLI-Slice schließt auf Hard-Gates-Form analog init/add,
nicht auf reduziertem `make test + lint + docs-check`.

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ Drei JSON-Pfade analog init (`runGenerate` ruft generische
  Helper mit `command="generate"` + `mapErr=mapGenerateErrorToDiagnostic`).
- ✅ **Vier-Artefakt-Symmetrie** (T0-(m) festgezurrt): identische
  Envelope-Shape unabhängig vom Artefakt. **`command="generate"`,
  kein `subcommand`-Feld** (Cobra-`<artifact>` ist Positional-Arg,
  nicht Subcommand — analog T0-(l)-Allowlist). Artefakt wird in
  `data.artifact: "<changelog|readme|env-example|devcontainer>"`
  geführt. Helper-Signatur (`reportError`/`writeErrorEnvelope`/
  `writeDiff`) bleibt unverändert — generate setzt nur `command`,
  kein `subcommand`.
- ✅ **Action-Klassifikation via `data.action`** (T0-(f)
  festgezurrt): Generate-Action wird im Voll-Schema-Envelope
  als Top-Level-Feld `data.action:
  "<created|updated-block|no-op|repaired-manual>"` getragen
  (schema-konform — `data` ist im Voll-Schema freies Feld).
  Begründung: `plannedFiles[i].action` (`create|modify|delete`)
  unterscheidet **nicht** zwischen `UpdatedBlock` (managed-block-
  rewrite) und `RepairedManual` (Single-Line-Header-Insert) —
  beide sind FS-semantisch `modify` mit nicht-leeren Arrays. Das
  zusätzliche `data.action`-Feld macht die Generate-Semantik
  eindeutig, ohne in das `changes[]`/`plannedFiles[]`-Schema
  einzugreifen (`jsontestutil.AssertFullEnvelope` bleibt
  unverändert anwendbar).
- ✅ **NoOp-Pin**: `--dry-run --json` bei bereits idempotenter
  Datei liefert `plannedFiles: []` UND `changes: []` (beide
  Arrays leer), `data.action: "no-op"`, `status: ok`, Exit 0.
- ✅ **UpdatedBlock-Hunks**: `--diff --json` bei `UpdatedBlock`
  rendert Hunks **nur** für den managed-block-Bereich (Sub-Decision
  (g): block-only vs. full-file-LCS). Envelope-Pin: `data.action:
  "updated-block"`, `plannedFiles[i].action: "modify"`,
  `plannedFiles[i].hunks[]` nicht-leer.
- ✅ **RepairedManual-Diff**: changelog-only Sonderfall mit Single-
  Line-Insert (`## [Unreleased]`-Header) — Diff-Pin testet, dass
  Hunks korrekt rendern (Sub-Decision (h)). Envelope-Pin:
  `data.action: "repaired-manual"`, `plannedFiles[i].action:
  "modify"` (FS-semantisch identisch zu UpdatedBlock — die
  Unterscheidung lebt **nur** in `data.action`).
- ✅ **Action-Diskriminations-Pin**: ein Acceptance-Test für
  changelog setzt zwei Szenarien gegenüber — (a) managed-block-
  ist-stale → `data.action: "updated-block"`, (b) `## [Unreleased]`-
  Header fehlt → `data.action: "repaired-manual"`. Beide
  produzieren identische `plannedFiles[i].action: "modify"`-
  Form; `data.action` ist die einzige Discriminator-Quelle.
- ✅ **Devcontainer-Pre-Write-Validation-Pin**: Phase 1
  (`planDevcontainerFiles`) ist atomar — wenn auch nur ein File
  als present-no-block / malformed klassifiziert wird, returnt
  der Use-Case `ErrGenerateManualConflict` (Exit 10) **ohne ein
  einziges WriteFile**. Acceptance-Pin: `--dry-run --json` mit
  einem manuell editierten `.devcontainer/devcontainer.json`
  ohne Marker → `plannedFiles: []`, **`diagnostics[].code:
  "LH-FA-DEV-001"`** (Devcontainer-Render-Spec, NICHT
  `LH-FA-CLI-006` — der ist Default-Fallback und würde Drift
  signalisieren), `exitCode: 10`, kein FS-Touch.
- ✅ **Devcontainer-Phase-2-Half-Write-Carveout** (T0-(i)):
  Mid-Write zweiter File in Phase 2 (`executeDevcontainerPlans`)
  hinterlässt **halbgeschriebenen Zustand** auf Disk — File 1
  ist committed, File 2 fehlt/teilweise. Carveout-Eintrag in
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts mit Re-Trigger auf einen Devcontainer-Rollback-aware-
  Slice (offen, V2-Scope; Cluster-Out-of-Scope per T0-(b) Variante 3).
  **Envelope-Form (Recorder-Realität, R4-Korrektur)**: der
  `RecordingFileSystem.WriteFile`
  (`recordingfs.go:139`) zeichnet den Aufruf **vor** dem
  Delegieren auf — der fehlgeschlagene File-2-Write steht
  deshalb in `plannedFiles[]` mit drin (planned[] = [File 1,
  File 2]). `lastPlannedPath` (`erroremission.go:73, 86-90`)
  liefert den letzten Eintrag, also File 2, als
  `diagnostics[].file`. Korrekte Pin-Form: `plannedFiles[]`
  enthält **beide** Captures (File 1 mit success-Content, File 2
  mit attempt-Content); `diagnostics[].code: "LH-NFA-REL-003"`,
  `diagnostics[].file: "<File 2-Pfad>"`, `exitCode: 14`.
  Konsistent mit init's Mid-Write-Failure-Pattern (init T6
  pinnt das genauso).
- ✅ **`--allow-external-feature-sources`-Side-Effect-Capture**:
  die Mutation auf `u-boot.yaml` (Z. 951) wird im Recorder erfasst
  und erscheint als zusätzlicher `plannedFiles[]`-Eintrag mit
  eigener Diff-Hunk-Sequenz — Sub-Decision (j).
- ✅ **Neuer `mapGenerateErrorToDiagnostic`** mit Switch-Order-
  Pflicht (ErrGenerateFileSystem FIRST, dann ErrGenerateManualConflict,
  dann ErrArtifactUnknown, dann Default) — Sub-Decision (e).
- ✅ **`ErrGenerateFileSystem`-Vervollständigung**: existiert
  bereits in port/driving/generate.go:154, aber Application-Code
  wrappt FS-Writes heute mit Single-`%w`. T3 erweitert auf Multi-
  `%w` (Switch-Order-Sicherheit) — Sub-Decision (d).
- ✅ **Allowlist-Erweiterung**: `"u-boot generate"` (parent-only)
  in `jsonallowlist.go`. Per-Artefakt-Form ist technisch
  unmöglich, weil `<artifact>` Positional-Arg ist und
  `cmd.CommandPath()` immer `"u-boot generate"` returnt
  (T0-(l) festgezurrt).
- ✅ **CLI-Pin-Tests**: 4 Artefakte × 8 Flag-Kombinationen
  (deckt die Aufhebungsbedingung 1:1 ab — die vier Human-Mode-
  Pfade ohne JSON sind öffentlicher CLI-Vertrag und müssen
  geprüft werden) plus NoOp/UpdatedBlock/RepairedManual/
  Devcontainer-Phase-1-Validation/Devcontainer-Phase-2-Half-Write/
  Allow-External-Side-Effect-Special-Pins. Total: ~36-40 Tests.
  Helper-Pattern analog `initFixture(t, opts)` für TempDir +
  u-boot.yaml-Setup (shared, ~80 LOC).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (generate→done),
  §6.5 neue Sektion, §7 Mutations-Matrix (generate-Zeile).
- ✅ **CHANGELOG `### Added`-Eintrag** analog init mit Pattern-
  Erbe-Notiz + generate-Spezifika.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a)** Pattern-Erbe-Disziplin (was 1:1 von init, was
  generate-spezifisch).
- **T0-(b)** **`driving.PreviewMode` direkt** (kein
  `GeneratePreviewMode`-Alias) — durch init-T0-(c) Alias-
  Lebensdauer-Pflicht erzwungen.
- **T0-(c)** `GenerateService.fsFactory`-Form analog
  `InitProjectService.fsFactory`.
- **T0-(d)** `ErrGenerateFileSystem`-Wrap-Audit: heute Single-`%w`
  in den 7 WriteFile-Stellen + 1 MkdirAll-Stelle. T3 erweitert
  auf Multi-`%w` analog init's Z. 925/967/1015/1117/1143-Stellen.
- **T0-(e)** **Switch-Order-Pflicht und Artefakt-Kontext** im
  neuen `mapGenerateErrorToDiagnostic`:
  (a) **Switch-Order**: ErrGenerateFileSystem FIRST (Multi-`%w`-
  Sicherheit, sonst Exit-14 → Exit-10-Downgrade);
  (b) **Mapper-Signatur**: erweitert auf `(err error, artifact
  domain.Artifact)`, weil `ErrGenerateManualConflict` pro Artefakt
  einen anderen LH-Anker hat (siehe Tabelle unten). Add/Init haben
  Single-Artefakt-Mapper und brauchen keinen Artefakt-Param —
  generate ist hier abweichend.
  (c) **Diagnostic-Code-Tabelle** (T6-Pin-Pflicht pro Zeile):

  | Sentinel | Artefakt | LH-Code | Exit |
  | --- | --- | --- | --- |
  | `ErrGenerateFileSystem` | * (alle) | `LH-NFA-REL-003` | 14 |
  | `ErrGenerateManualConflict` | changelog | `LH-FA-GEN-002` | 10 |
  | `ErrGenerateManualConflict` | readme | `LH-FA-GEN-003` | 10 |
  | `ErrGenerateManualConflict` | env-example | `LH-FA-GEN-004` | 10 |
  | `ErrGenerateManualConflict` | devcontainer | `LH-FA-DEV-001` | 10 |
  | `ErrArtifactUnknown` | * (alle) | `LH-FA-CLI-006` | 2 |
  | Default (unknown) | * (alle) | `LH-FA-CLI-006` | 1 |

  Switch-Order verbindlich: FS-first, dann ManualConflict
  (artefakt-spezifischer Code via `artifact`-Param), dann
  ArtifactUnknown, dann Default.
- **T0-(f)** **Action-Klassifikation via `data.action` festgezurrt**
  (R3-Festzurrung, R2-Variante „kein Marker" verworfen, weil sie
  UpdatedBlock vs. RepairedManual nicht unterscheidbar machte):
  Generate-Action wird im Voll-Schema-Envelope als
  `data.action: "<created|updated-block|no-op|repaired-manual>"`
  getragen. NoOp produziert zusätzlich `plannedFiles: []` UND
  `changes: []` (beide leer; ein `changes`-Eintrag mit „nur
  `action`" wäre schema-illegal — `jsontestutil.AssertFullEnvelope`/
  `checkChanges` Z. 384-400 enforced `path`/`count` Pflicht pro
  `changes[]`-Eintrag, Spec §365/§368). Created/UpdatedBlock/
  RepairedManual produzieren nicht-leere Arrays; `data.action`
  ist der eindeutige Discriminator zwischen UpdatedBlock und
  RepairedManual (beide haben `plannedFiles[i].action: "modify"`,
  FS-semantisch identisch). Human-Mode bleibt unverändert
  (`printGenerateSummary`-Branches).
- **T0-(p)** **`cliJSONEnvelope.Data`-Feld-Migration vorgezogen**
  (R4-Finding 2): heutige `cliJSONEnvelope`
  (`internal/adapter/driving/cli/jsonenvelope.go:38-44`) hat
  **kein** `Data`-Feld — der Doctor-Slice-T0-(c)-Kommentar
  reserviert es für den Template-Slice 9/9
  (`newDataEnvelope`-Konstruktor + Pin-Test). Generate ist 4/9
  und der **erste** Voll-Schema-Konsument mit Data-Feld-Bedarf
  (`data.artifact` aus T0-(m) + `data.action` aus T0-(f)).
  Sub-Decision: Generate-Slice **zieht die Migration vor** —
  neue T1-Sub-Tranche oder T2-Erweiterung ergänzt
  `Data any \`json:"data,omitempty"\`` plus
  `newDataEnvelope(command, subcommand string, data any, diags
  []diagnosticItem, exitCode int)`-Konstruktor mit Marshal-Pin-
  Test (analog Template-Slice-Plan §Akzeptanzkriterien H1-Finding).
  **Signatur mit `subcommand string`-Param** (R5-Finding):
  Generate ruft mit `subcommand=""` (T0-(m) — kein Subcommand);
  Template-Slice 9/9 ruft mit `subcommand="list"`. Der Konstruktor
  serialisiert `subcommand` mit `omitempty`, so dass die Generate-
  Envelope das Feld nicht trägt und die Template-Envelope schon.
  Damit muss Template-Slice den Konstruktor nicht erneut ändern —
  Ownership-Verschiebung ist hier abgeschlossen. T8-Closure
  aktualisiert den Template-Slice-Plan: Data-Feld + Konstruktor
  werden dort nur noch genutzt, nicht eingeführt. Begründung
  insgesamt: ohne diese Vorziehung kann generate die T0-(f)/(m)-
  Verträge nicht umsetzen, und der Konstruktor wäre ohne den
  subcommand-Param für Template unzureichend.
- **T0-(g)** **UpdatedBlock-Diff-Granularität**: managed-block-
  only Hunks vs. full-file-LCS. Block-only ist semantisch sauber
  (User sieht NUR die u-boot-managed Änderung), aber existing
  Renderer macht full-file — Migration nötig?
- **T0-(h)** **RepairedManual-Diff-Pin**: changelog-Header-Insert
  als minimaler Hunk.
- **T0-(i)** **Devcontainer-Atomicity-Klärung**: das Application-
  Code-Kommentar (`generate.go:618-624`) beschreibt **Pre-Write-
  Validation-Atomicity** (Phase 1 ist atomar — kein WriteFile bei
  Validation-Conflict), NICHT Roll-back-Atomicity in Phase 2.
  Sub-Decision pinnt zwei Acceptance-Verhalten:
  (a) Phase-1-Failure (manuell editierter devcontainer-File ohne
      Marker) → `plannedFiles: []`, kein FS-Touch, Exit 10
      (`LH-FA-DEV-001`).
  (b) Phase-2-Mid-Write-Failure → File 1 committed auf Disk,
      File 2 underlying-Write fehlgeschlagen. Envelope-Form
      (R4-Recorder-Realität): `plannedFiles[]` enthält **beide**
      Captures (File 1 mit success-Content, File 2 mit
      attempt-Content — Recorder zeichnet vor Delegieren auf,
      `recordingfs.go:139`), `lastPlannedPath` liefert File 2
      als `diagnostics[].file`, `diagnostics[].code:
      "LH-NFA-REL-003"`, Exit 14. Konsistent mit init T6
      Mid-Write-Pattern.
  Der **Phase-2-Half-State** ist ein bewusster Carveout
  ([[feedback_carveouts_need_plans]]); `carveouts.md`-Eintrag
  mit Re-Trigger auf einen Devcontainer-Rollback-aware-Slice
  (V2-Scope, Cluster-T0-(b) Variante 3 verworfen). Pre-Write-
  Validation-Atomicity bleibt der ehrliche Vertrag.
- **T0-(j)** **`--allow-external-feature-sources`-Side-Effect**:
  Mutation auf `u-boot.yaml` als zusätzlicher `plannedFiles[]`-
  Eintrag im Envelope. Diff zeigt YAML-Hunks. Acceptance-Pin
  testet beide Schreib-Operationen (devcontainer-Files + yaml).
- **T0-(k)** Path-Anchor: `plannedFiles[].path` ist project-
  relativ (analog init T0-(k)) — `mapCaptureToPlannedFiles(records,
  baseDir)`-Erbe.
- **T0-(l)** **Allowlist-Form festgezurrt: parent-only**
  `"u-boot generate"`. Per-Artefakt-Form (`"u-boot generate
  changelog"` etc.) ist **nicht möglich** — `<artifact>` ist
  ein Cobra-Positional-Argument, kein Subcommand, und
  `cmd.CommandPath()` returnt für jedes Artefakt nur
  `"u-boot generate"`. Per-Artefakt-Einträge würden die Reject-
  Gate-Mechanik (`applyJSONRejectGate` in `jsonallowlist.go`)
  niemals matchen — `--json` würde dauerhaft rejected bleiben.
  Konsistenz zu init/add ist nebensächlich; der eigentliche
  Grund ist die CommandPath-Semantik.
- **T0-(m)** **Envelope-Shape festgezurrt**:
  `command="generate"`, **kein `subcommand`-Feld**, Artefakt im
  `data.artifact:
  "<changelog|readme|env-example|devcontainer>"`. Begründung:
  (1) Cobra-`<artifact>` ist Positional-Arg, kein Subcommand
  (analog T0-(l)); ein synthetisches `subcommand`-Feld im
  Envelope würde von der CLI-Layer-Realität abweichen.
  (2) Die generischen Error-Emission-Helper (`reportError`/
  `writeErrorEnvelope` etc.) setzen heute **kein** Subcommand —
  bei Wahl von `subcommand="<artifact>"` wäre eine Helper-
  Signatur-Erweiterung nötig (Drift-Risk gegen init/add).
  Plus: Action-Klassifikation lebt ohnehin in `data.action`
  (T0-(f)), so dass `data` der natürliche Träger für
  Generate-spezifische Felder ist.
- **T0-(n)** **`Codes`-Registry-Ergänzung NICHT nötig**: die
  Codes-Map in `jsontestutil.DefaultAllowedCodes` (cli-json-output.md
  §5) ist für tool-interne **dotted** Codes (`add.*`, `init.*`-
  Style) gedacht — LH-Codes werden vom Acceptance-Helper bereits
  generisch erlaubt (`LH-FA-GEN-001..005`/`LH-FA-DEV-001/003`/
  `LH-NFA-REL-003`/`LH-FA-CLI-006` sind alle LH-präfixiert). T6
  ergänzt keine Registry-Einträge; Doku-Update für §6.5
  (Per-Command-Sektion) genügt.
- **T0-(o)** Pre-`next/`-Review-Runden-Erwartung: init hatte
  3 vor `next/`, add hatte 5. Generate: ≥ 2 (Discovery-Tiefe +
  Adversarial).

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | ------ | --------------- | --- |
| T0 | Discovery + Sub-Decisions (a)-(o) klären; Review-Runden | — (Plan) | — |
| T1 | Refactor-Tranche (wenn überhaupt nötig — generate hat schmalere FS-Surface; ggf. nur ErrGenerateFileSystem-Multi-`%w`-Audit oder gar kein T1) | ~30-80 | T0 |
| T2 | Port-Types: `GenerateRequest.PreviewMode`, `GenerateResponse.PlannedFiles`/`Changes`-Felder. **`data.action`-Klassifikation** liegt im Envelope-Layer (T5), nicht im Port — die existierende `GenerateResponse.Action` (`GenerateAction`-Enum) wird in T5 zum `data.action`-String gerendert; keine neuen Port-Felder dafür (T0-(f) Festzurrung). `ErrGenerateFileSystem` ist schon da. | ~50 | T0 |
| T3 | Application-Layer: `GenerateService.fsFactory` + `generateMu sync.Mutex` + `NewGenerateServiceWithFactory` + `Generate()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; Multi-`%w`-Wrap an den 8 FS-Wrap-Stellen. | ~200 | T2 |
| T4 | Composition-Root-Wiring `generateFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runGenerate` ruft generische Helper mit `command="generate"` (kein subcommand, T0-(m)), `mapErr=mapGenerateErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; **`mapGenerateErrorToDiagnostic(err, artifact)` neu mit Artefakt-Parameter** (T0-(e); per-Artefakt LH-Code für ErrGenerateManualConflict). `data.action` aus `resp.Action.String()` gerendert; `data.artifact` aus `req.Artifact.String()`. Helper-Generalisierung (`reportError`/`writeErrorEnvelope`) bleibt unverändert (Signatur trägt heute kein subcommand, T0-(m)). | ~200 | T1 + T2 (T4 für Run-time-Smoke, Code-parallelisierbar) |
| T6 | Acceptance-Tests: 4 Artefakte × 8 Flag-Kombos (4 Human-Mode + 4 JSON, deckt Aufhebungsbedingung 1:1) + NoOp/UpdatedBlock/RepairedManual/Devcontainer-Phase-1-Validation/Devcontainer-Phase-2-Half-Write/Allow-External-Side-Effect-Pins; Helper `generateFixture(t, opts)` shared (~80 LOC). | ~640 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden). | ~80 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.5/§7, roadmap, slice nach done/ mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: ~1080-1130 — schmaler als init (~1480), weil
keine `initGit`/`Progress`-Carveouts und schmalere FS-Surface,
aber breiter als add (~1380? — add war Pattern-Founder mit viel
Infrastruktur).

## Review-Round-1 (Pre-`next/`)

Eine adversarial-orientierte Review-Runde gegen den initialen
Stub (`fbef9b5`). Fünf Findings (2 HIGH, 2 MEDIUM, 1 LOW), alle
adressiert im selben Commit wie dieser R1-Block:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH | Coverage-Gate fehlt im Slice-Abschluss (Aufhebungsbedingung nannte nur `test + lint + docs-check` — Code-ändernder modifying-CLI-Slice braucht `make gates`-Hard-Form mit coverage-gate ≥ 90 %) | Aufhebungsbedingung auf `make gates` umgestellt (inkl. expliziter Coverage-Gate-Pflicht) |
| 2 | HIGH | Devcontainer-Atomicity-Widerspruch (Stub sagte „zwei Files atomar oder gar nicht", pinnte gleichzeitig Half-Write-State + Roll-back-aware Out-of-Scope) | Atomicity-Klärung in Pre-Scan-Tabelle und T0-(i): Phase 1 (`planDevcontainerFiles`) ist atomar (Pre-Write-Validation), Phase 2 (`executeDevcontainerPlans`) NICHT — Half-State ist bewusster Carveout mit `carveouts.md`-Eintrag und Re-Trigger auf späteren Rollback-Slice |
| 3 | MEDIUM | Testmatrix deckt Aufhebungsbedingung nicht ab (8 Flag-Kombos pro Artefakt vs. AK-Plan 4×4 — Human-Mode-Pfade ohne JSON fehlten) | AK auf 4 Artefakte × 8 Flag-Kombos + Special-Pins erweitert, Total ~36-40 Tests |
| 4 | MEDIUM | Per-Artefakt-Allowlist ist kein gültiger CommandPath (Cobra `cmd.CommandPath()` für `u-boot generate <artifact>` ist nur `"u-boot generate"`, weil `<artifact>` Positional-Arg ist; per-Artefakt-Einträge würden Reject-Gate nie matchen) | T0-(l) auf **parent-only** festgezurrt mit CommandPath-Semantik-Begründung |
| 5 | LOW | Codes-Map-Ergänzung driftet gegen Registry-Konvention (LH-Codes sind im Test-Helper bereits generisch erlaubt; Registry ist für tool-interne dotted Codes) | T0-(n) auf „NICHT nötig" umgestellt — nur §6.5-Per-Command-Doku, keine Registry-Einträge |

R1-Reviewer-Note: Markdown-Link-Sensor (`make docs-check`)
und Spec-Anker-Bezüge sind grün; Pattern-Verweise auf
PreviewMode/Recorder/Diff-Helper/Error-Emission existieren im
aktuellen Code + done-Slices.

## Review-Round-2 (Pre-`next/`)

Zweite adversarial-orientierte Runde gegen den R1-gepflegten Stub
(`f3134bd`), Fokus auf Plan-Drift, Carveout-Inventarisierung und
Schema-Kompatibilität. Fünf Findings (1 HIGH, 3 MEDIUM, 1 LOW),
alle adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH | Phase-2-Half-Write-Carveout war im Slice angekündigt, aber NICHT in `carveouts.md` inventarisiert und ohne Plan-Stub (verletzt `LH-FA-PROJDOCS-005` und MEMORY-Rule [[feedback_carveouts_need_plans]]) | Neuer Eintrag in `carveouts.md` §Temporäre Carveouts mit Re-Trigger; neuer open/-Plan-Stub `slice-v2-generate-devcontainer-rollback-aware-write` (Status `on hold pending trigger`, drei Lösungs-Skizzen) |
| 2 | MEDIUM | Header-Vertrag „zwei Files atomar oder gar nicht" widersprach T0-(i)-Ergebnis (Phase 1 atomar, Phase 2 nicht) | Header-Text auf „Devcontainer-Atomicity-**Asymmetrie**" umgestellt mit expliziter Phase-1/Phase-2-Trennung |
| 3 | MEDIUM | T6-Tranchen-Zelle plante 4×4-Matrix; Aufhebungsbedingung + AK fordern 4×8 | T6-Zelle auf 4×8 + Special-Pins erweitert, LOC-Schätzung ~640 (vorher ~500) |
| 4 | MEDIUM | T0-(f) NoOp-Action-Marker war schema-unklar — `changes: []` mit `action: no-op` wäre schema-illegal (Helper enforced `path`/`count` pro `changes[]`-Eintrag) | T0-(f) festgezurrt: NoOp = `plannedFiles: []` UND `changes: []`, keine Action-Marker-Schema-Erweiterung; Konsumenten leiten NoOp aus Leerheit beider Arrays ab. AK-Block entsprechend nachgezogen |
| 5 | LOW | AK „Allowlist-Erweiterung" referenzierte Sub-Decision (l) als offen, obwohl T0-(l) parent-only bereits festzurrt | AK-Zeile auf festgezurrte parent-only-Form umgestellt, Sub-Decision-Verweis raus |

R2-Reviewer-Note: docs-check grün, Pattern-Verweise weiterhin
tragfähig.

## Review-Round-3 (Pre-`next/`)

Dritte Runde gegen den R2-gepflegten Stub (`cffda2c`), Fokus auf
V2-Rollback-Korrektheit und Plan-interne Konsistenz. Vier
Findings (1 HIGH gegen V2-Stub, 3 MEDIUM gegen V1-Stub), alle
adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH (V2) | V2-Stub bevorzugte Option 3 (per-File Temp+Rename) löste das Multi-File-Half-State-Problem NICHT — wenn Rename 1 succeeds und Rename 2 failt, bleibt File 1 committed. Per-File-Atomicity ≠ Multi-File-Atomicity. | V2-Stub bevorzugte Skizze auf **Option 1 (Snapshot + Rollback-on-Failure)** umgestellt — echte Multi-File-Atomicity mit Best-Effort-Rollback. Option 3 explizit als „verworfen" gelabelt. Failure-Injection-Pin im Trigger-Slice ergänzt: „erste Datei committed, zweite Rename/YAML-Write failt → Restore aktiviert" |
| 2 | MEDIUM | Action-Vertrag widersprüchlich: AK fordert vier Generate-Actions, T0-(f) verbietet Action-Marker, T2 plante zugleich ein „Action-Marker-Feld". `plannedFiles[i].action: "modify"` unterscheidet UpdatedBlock und RepairedManual nicht. | T0-(f) auf **`data.action: "<…>"`** umgestellt (Top-Level-Voll-Schema-Feld in `data`, schema-konform). T2-Tranchen-Zelle: keine Port-Felder, Rendering im T5-Layer. Action-Diskriminations-Acceptance-Pin in den AK-Block (changelog UpdatedBlock vs. RepairedManual mit identischem `plannedFiles[i].action: "modify"`, Discriminator `data.action`) |
| 3 | MEDIUM | Diagnostic-Code zu lose gepinnt (`[LH-FA-CLI-006 oder LH-FA-DEV-001]`); ErrGenerateManualConflict braucht Artefakt-Kontext, geplanter `mapGenerateErrorToDiagnostic(err)` hat keinen | T0-(e) erweitert um **Mapper-Signatur `(err, artifact)`** mit Begründung und um eine **per-Artefakt LH-Code-Tabelle** (changelog→GEN-002, readme→GEN-003, env-example→GEN-004, devcontainer→DEV-001). Acceptance-Pin auf exakten `LH-FA-DEV-001` umgestellt (statt „either/or"). T5-Tranche um Mapper-Signatur-Erweiterung ergänzt |
| 4 | MEDIUM | Envelope-Shape offen (T0-(m)) — `subcommand="<artifact>"` vs. `data.artifact`. Helper-Signatur-Drift-Risiko, weil heutige Helper kein subcommand setzen | T0-(m) festgezurrt auf **`command="generate"`, kein subcommand, `data.artifact: "<…>"`**. Begründungen: (1) Cobra-Positional-Arg-Semantik analog T0-(l), (2) Helper-Signatur bleibt unverändert, (3) `data` ist ohnehin Träger für `data.action`. AK Vier-Artefakt-Symmetrie-Zeile entsprechend |

R3-Reviewer-Note: docs-check grün, V2-Stub-Korrektur ist die
wichtigste Erkenntnis (Stub hätte sonst eine falsche Lösung
empfohlen, die den Carveout nicht schließt).

## Review-Round-4 (Pre-`next/`)

Vierte Runde gegen den R3-gepflegten Stub (`0edcf25`), Fokus
Plan-Drift, Carveout-Inventarisierung-Sichtbarkeit und Recorder-
Realität. Vier Findings (1 HIGH gegen V2-Stub/Roadmap, 3 MEDIUM
gegen V1-Stub), alle adressiert im selben Commit:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH | V2-Slice fehlte als Roadmap-Slice-Zeile (LH-FA-PROJDOCS-005 fordert Doppel-Sichtbarkeit in `carveouts.md` UND `roadmap.md`) | Neuer Roadmap-Eintrag in `roadmap.md` §AP-Tabelle für `slice-v2-generate-devcontainer-rollback-aware-write` mit Status `on hold pending trigger` und Carveout-Anker-Verweis |
| 2 | MEDIUM | `data.action`/`data.artifact`-Vertrag widersprüchlich: AK forderte beide Felder, T0-(f) sagte aber gleichzeitig „nicht über separates Top-Level-Feld geführt" (R2-Drift). Zusätzlich: heutige `cliJSONEnvelope` hat **kein** `Data`-Feld (Doctor-Slice T0-(c) reserviert es für Template-Slice 9/9 mit `newDataEnvelope`-Konstruktor). | T0-(f)-Text konsolidiert auf R3-Festzurrung (`data.action` als eindeutiger Discriminator); R2-Drift-Sätze entfernt. **Neuer T0-(p)**: `cliJSONEnvelope.Data`-Feld-Migration aus Template-Slice 9/9 in den generate-Slice vorgezogen — Generate ist erster Data-Konsument, ergänzt `Data any \`json:"data,omitempty"\`` plus `newDataEnvelope`-Konstruktor in T1/T2-Sub-Tranche; T8-Closure pflegt den Template-Slice-Plan nach |
| 3 | MEDIUM | Mid-Write-Envelope passt nicht zur Recorder-Semantik: Plan sagte `plannedFiles[]` zeigt File 1, `diagnostics[].file` markiert File 2 — aber `RecordingFileSystem.WriteFile` (`recordingfs.go:139`) zeichnet VOR dem Delegieren auf, der fehlgeschlagene Write steht also in `plannedFiles[]` mit drin | Devcontainer-Phase-2-Pin auf Recorder-Realität korrigiert: `plannedFiles[]` enthält BEIDE Captures (File 1 success-Content, File 2 attempt-Content), `lastPlannedPath` liefert File 2 als `diagnostics[].file` — konsistent mit init's Mid-Write-Failure-Pattern (init T6) |
| 4 | MEDIUM | V2-Rollback-Skizze deckt Verzeichnis-/Scratch-State nicht ab: `MkdirAll` erstellt `.devcontainer/` vor Write-Sequenz; `.bak.<n>`-Snapshots wären Scratch-Artefakte. „Disk-Zustand nach Aufruf == vor Aufruf" muss diese explizit fordern. | V2-Stub um **Rollback-Scope-Erklärung** erweitert: drei Side-Effects (Dir-Anlage, YAML-Mutation, Snapshot-Persistierung); Acceptance-Pin auf frischem Projekt ohne `.devcontainer/`-Dir mit `tree`-Vergleich Pre-State == Post-State |

R4-Reviewer-Note: docs-check grün; Pattern-Verweise konsistent
nach Plan-Konsolidierung; Recorder-Code-Realität war der
load-bearing Befund (Plan hatte init's eigenes Pattern nicht
übernommen).

## Review-Round-5 (Pre-`next/`)

Fünfte Runde gegen den R4-gepflegten Stub (`0b3e1ad`), Fokus auf
intra-Plan-Konsistenz nach den Konsolidierungen und Acceptance-
Pin-Vollständigkeit. Drei MEDIUM-Findings, alle adressiert:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | MEDIUM | T0-(i) widersprach der R4-korrigierten Recorder-Realität (sagte noch „plannedFiles[] zeigt File 1") | T0-(i) auf R4-Form gezogen: „plannedFiles[] enthält beide Captures" mit Recorder-Verweis und init-T6-Konsistenz-Notiz |
| 2 | MEDIUM | T0-(p) skizzierte `newDataEnvelope(command, data, diags, exitCode)` ohne `subcommand`-Param — Template-Slice 9/9 braucht `subcommand="list"` und müsste den Konstruktor erneut ändern | Konstruktor-Signatur auf `newDataEnvelope(command, subcommand string, data any, diags, exitCode)` erweitert mit `omitempty`-Serialisierung (Generate ruft mit `subcommand=""`, Template mit `subcommand="list"`); Ownership-Verschiebung damit final |
| 3 | MEDIUM | V2-Acceptance-Pin ließ File 2 failen, das passiert VOR der u-boot.yaml-Mutation — YAML-Rollback-Pfad wurde nie geprüft | V2-Stub auf **zwei separate Failure-Injection-Pins** umgestellt: Pin A (File-2-Failure, Devcontainer-Rollback + Dir-Cleanup) und Pin B (u-boot.yaml-Write-Failure NACH Devcontainer-Success, voller Three-Side-Effect-Rollback) |

R5-Reviewer-Note: docs-check grün; Stub konsolidiert ohne neue
HIGH-Befunde; weitere Runden sind diminishing returns.

## Out of Scope

- **Roll-back-aware Recorder** für Devcontainer-Phase-2-Atomicity
  (Cluster-T0-(b) Variante 3 verworfen, V2-Scope). Phase 2-Mid-
  Write hinterlässt halbgeschriebenes File 1 auf Disk; Carveout-
  Eintrag in `carveouts.md` mit Re-Trigger auf einen späteren
  Devcontainer-Rollback-Slice.
- **HTTP- oder gRPC-Schnittstellen**: ADR-0010 schließt
  explizit aus.
- **Schema-Versionierung** (`schemaVersion: 1`): siehe
  Cluster-Slice §Out of Scope.
- **Neue Artefakte** (z. B. `generate dockerfile` standalone):
  V1-Scope ist die existierenden vier.
- **Generic `mapErrorToDiagnostic`-Registry**: Altitude-Reviewer-
  Vorschlag aus add R6 #I1. Cluster-T_close-Aufgabe.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  §T0-Outcomes — Vorgaben für den Folge-Slice-Block.
- Pattern-Vorbild:
  [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
  — T0-T8 + Review-Round-9 voll abgeschlossen. Erbschafts-
  Disziplin in T0-(a) dieses Slices; Alias-Lebensdauer-Pflicht
  aus init-T0-(c) zwingt direktes `driving.PreviewMode`.
- Add-Slice (sekundär):
  [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
  — Pattern-Founder; relevante Sub-Decisions
  (`CountAdditions`-Semantik §477, `checkHunks`-Helper) bleiben
  geerbt.
- Cleanup-Stub:
  [`slice-v1-cli-cleanup-add-preview-mode-alias`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md)
  — wartet auf „mindestens einen weiteren Folge-Slice"; generate
  IST dieser Folge-Slice, MUSS `driving.PreviewMode` direkt
  nutzen.
- Spec: `LH-FA-CLI-007/008`, `LH-NFA-USE-004`,
  `LH-FA-GEN-001..005`, `LH-FA-DEV-001/003`, `LH-NFA-REL-003`
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- Code-Anker heute:
  [`generate.go`](../../../../internal/hexagon/application/generate.go)
  (~960 LOC, vier Artefakt-Handler + managedblock-Surface),
  [`cli/generate.go`](../../../../internal/adapter/driving/cli/generate.go)
  (~160 LOC, RunE-Erweiterungs-Ziel),
  [`port/driving/generate.go`](../../../../internal/hexagon/port/driving/generate.go)
  (Carrier-Types + drei Sentinels inkl. `ErrGenerateFileSystem`).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
