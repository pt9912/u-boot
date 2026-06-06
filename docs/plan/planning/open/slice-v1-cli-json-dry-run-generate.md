# Slice V1: `generate --json` / `--dry-run` / `--diff` — Vier-Artefakt-Surface

> **Status:** T0-Discovery, `open/`. Vierter Folge-Slice (4/9) des
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
| devcontainer Multi-File (`writeDevcontainerPlan`) | `MkdirAll` (Z. 848) + `WriteFile` × 2 (Z. 852/864) | `.devcontainer/devcontainer.json` + `.devcontainer/Dockerfile` | direkt; atomar-pflichtig |
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

`make test` + `make lint` + `make docs-check` grün.

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
- ✅ **Devcontainer-Atomicity-Pin**: Mid-Write-Failure beim
  zweiten der zwei Files → `plannedFiles[]` enthält **nur** die
  Capture bis zur Failure-Stelle, `exitCode: 14`. Atomar-or-nothing
  ist Application-Layer-Vertrag (existierender Code-Kommentar
  `generate.go:623-625`), Sub-Decision (i) entscheidet, ob der
  Recorder Roll-back-aware ist (Cluster-Out-of-Scope: nein).
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
- ✅ **CLI-Pin-Tests**: ~16+ Acceptance-Tests (4 Artefakte × 4
  Flag-Kombos plus NoOp/UpdatedBlock/RepairedManual/Devcontainer-
  Atomicity-Special-Pins).
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
- **T0-(i)** **Devcontainer-Atomicity-Carveout**: Recorder ist
  NICHT Roll-back-aware (Cluster-T0-(b) Variante 3 verworfen).
  Mid-Write zweiter File → halbgeschriebener Zustand auf Disk;
  `plannedFiles[]` zeigt nur File 1 + diag-Marker auf File 2.
  Application-Code-Kommentar `generate.go:623-625` dokumentiert
  die Pflicht — Sub-Decision pinnt das Acceptance-Verhalten.
- **T0-(j)** **`--allow-external-feature-sources`-Side-Effect**:
  Mutation auf `u-boot.yaml` als zusätzlicher `plannedFiles[]`-
  Eintrag im Envelope. Diff zeigt YAML-Hunks. Acceptance-Pin
  testet beide Schreib-Operationen (devcontainer-Files + yaml).
- **T0-(k)** Path-Anchor: `plannedFiles[].path` ist project-
  relativ (analog init T0-(k)) — `mapCaptureToPlannedFiles(records,
  baseDir)`-Erbe.
- **T0-(l)** **Allowlist-Form**: parent `"u-boot generate"` oder
  per-Artefakt? Init/add nutzen parent-only — Konsistenz-
  Empfehlung.
- **T0-(m)** **Envelope-Shape**: `command="generate"` mit
  `subcommand="<artifact>"` (analog `template list`-Form), oder
  `command="generate"` ohne subcommand und Artefakt im `data`-
  Block?
- **T0-(n)** **`Codes`-Map-Ergänzung**: generate-eigene Codes
  (heute `LH-FA-GEN-001..005` im CLI, `LH-FA-DEV-001/003`, plus
  geerbte `LH-NFA-REL-003`/`LH-FA-CLI-006`).
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

## Out of Scope

- **Roll-back-aware Recorder** für Devcontainer-Atomicity:
  Cluster-T0-(b) Variante 3 (ChangeSet-Pattern) ist V1-out-of-
  scope. Mid-Write zweiter devcontainer-File hinterlässt
  halbgeschriebenes File 1 auf Disk; `plannedFiles[]` zeigt
  Position bis zur Failure-Stelle.
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
