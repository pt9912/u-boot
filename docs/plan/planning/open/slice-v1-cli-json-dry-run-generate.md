# Slice V1: `generate --json` / `--dry-run` / `--diff` — Vier-Artefakt-Surface

> **Status:** T0-Discovery + R1 adressiert, `open/`. Vierter Folge-Slice (4/9) des
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
> devcontainer `MkdirAll`), die **Devcontainer-Atomicity-Pflicht**
> (zwei Files atomar oder gar nicht, bestehender Code-Kommentar
> `generate.go:623-625`) und der **`--allow-external-feature-sources`-
> Side-Effect** (mutiert `u-boot.yaml` als zusätzlichen Schreib-
> Vorgang außerhalb des Artefakt-Targets).
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
- ✅ **Vier-Artefakt-Symmetrie**: identische Envelope-Shape unabhängig
  vom Artefakt, mit `subcommand="<artifact>"` (oder `command="generate"`
  + Artefakt im `data.artifact` — Sub-Decision (m)).
- ✅ **Action-Klassifikation** im Voll-Schema-Envelope: `changes[]`
  trägt die `GenerateAction`-String-Form (`created` / `updated-block`
  / `no-op` / `repaired-manual`) im `kind`-Feld (oder eigenem
  `action`-Feld — Sub-Decision (f)).
- ✅ **NoOp-Pin**: `--dry-run --json` bei bereits idempotenter
  Datei liefert `plannedFiles: []` (kein WriteFile-Capture) plus
  `status: ok` und eine `action: no-op`-Marker — Sub-Decision (f).
- ✅ **UpdatedBlock-Hunks**: `--diff --json` bei `UpdatedBlock`
  rendert Hunks **nur** für den managed-block-Bereich (Sub-Decision
  (g): block-only vs. full-file-LCS).
- ✅ **RepairedManual-Diff**: changelog-only Sonderfall mit Single-
  Line-Insert (`## [Unreleased]`-Header) — Diff-Pin testet, dass
  Hunks korrekt rendern (Sub-Decision (h)).
- ✅ **Devcontainer-Pre-Write-Validation-Pin**: Phase 1
  (`planDevcontainerFiles`) ist atomar — wenn auch nur ein File
  als present-no-block / malformed klassifiziert wird, returnt
  der Use-Case `ErrGenerateManualConflict` (Exit 10) **ohne ein
  einziges WriteFile**. Acceptance-Pin: `--dry-run --json` mit
  einem manuell editierten `.devcontainer/devcontainer.json`
  ohne Marker → `plannedFiles: []`, `diagnostics: [LH-FA-CLI-006
  oder LH-FA-DEV-001]`, kein FS-Touch.
- ✅ **Devcontainer-Phase-2-Half-Write-Carveout** (T0-(i)):
  Mid-Write zweiter File in Phase 2 (`executeDevcontainerPlans`)
  hinterlässt **halbgeschriebenen Zustand** auf Disk — File 1
  ist committed, File 2 fehlt/teilweise. Carveout-Eintrag in
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts mit Re-Trigger auf einen Devcontainer-Rollback-aware-
  Slice (offen, V2-Scope; Cluster-Out-of-Scope per T0-(b) Variante 3).
  Envelope-Form: `plannedFiles[]` zeigt File-1-Capture, `diagnostics[].file`
  markiert File 2 als Failure-Position, `exitCode: 14`.
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
- ✅ **Allowlist-Erweiterung**: `"u-boot generate"` in
  `jsonallowlist.go` (mit Sub-Decision: einer pro Artefakt
  `"u-boot generate changelog"` etc. oder nur die parent-Form? —
  Sub-Decision (l)).
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
- **T0-(e)** **Switch-Order-Pflicht** im neuen
  `mapGenerateErrorToDiagnostic`: ErrGenerateFileSystem FIRST,
  weil Multi-`%w` sonst Exit-14 auf Exit-10 downgraded.
- **T0-(f)** **NoOp-Envelope-Form**: `plannedFiles: []` plus
  Action-Marker (`changes: []` mit `action: no-op`? Oder
  `data.action: "no-op"` als Voll-Schema-Top-Level-Feld?).
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
  (b) Phase-2-Mid-Write-Failure → File 1 committed, File 2 fehlt;
      `plannedFiles[]` zeigt File 1, `diagnostics[].file` markiert
      File 2, Exit 14 (`LH-NFA-REL-003`).
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
- **T0-(m)** **Envelope-Shape**: `command="generate"` mit
  `subcommand="<artifact>"` (analog `template list`-Form), oder
  `command="generate"` ohne subcommand und Artefakt im `data`-
  Block?
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
| T2 | Port-Types: `GenerateRequest.PreviewMode`, `GenerateResponse.PlannedFiles`/`Changes`-Felder, Action-Marker-Feld (Sub-Decision (f)); `ErrGenerateFileSystem` ist schon da. | ~60 | T0 |
| T3 | Application-Layer: `GenerateService.fsFactory` + `generateMu sync.Mutex` + `NewGenerateServiceWithFactory` + `Generate()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; Multi-`%w`-Wrap an den 8 FS-Wrap-Stellen. | ~200 | T2 |
| T4 | Composition-Root-Wiring `generateFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runGenerate` ruft generische Helper mit `command="generate"`, `mapErr=mapGenerateErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; `mapGenerateErrorToDiagnostic` neu. | ~180 | T1 + T2 (T4 für Run-time-Smoke, Code-parallelisierbar) |
| T6 | Acceptance-Tests: 4 Artefakte × 4 Flag-Kombos + NoOp/UpdatedBlock/RepairedManual/Devcontainer-Atomicity/Allow-External-Side-Effect-Pins. | ~500 | T5 |
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
  [`slice-v1-cli-cleanup-add-preview-mode-alias`](slice-v1-cli-cleanup-add-preview-mode-alias.md)
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
