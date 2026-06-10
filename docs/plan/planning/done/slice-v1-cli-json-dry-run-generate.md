# Slice V1: `generate --json` / `--dry-run` / `--diff` — Vier-Artefakt-Surface

> **Status:** ✅ **done** — vierter Folge-Slice (4/9) des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 4/9). Lifecycle-Übergänge: `open/` nach R1-R5,
> `next/` nach R6/R7, `in-progress/` ab T1, `done/` mit T8-Closure.
> Findings-Bilanz: 33 gesamt (7 HIGH, 19 MED, 7 LOW); R10 ergab
> 5 (0 HIGH, 2 MED, 3 LOW) gegen die T1-T6-Implementation.
>
> Generate ist der **erste** Subcommand, der mehrere Artefakte
> (changelog/readme/env-example/devcontainer) über einen einzigen
> Subcommand bedient. Pattern-Erbe init→generate 1:1; Generate-
> spezifische Erweiterungen: `command="generate"` ohne
> `subcommand`-Feld (Cobra-Positional-Arg-Semantik),
> `data.artifact`/`data.action`-Carrier-Form, per-Artefakt
> LH-Code-Tabelle in `mapGenerateErrorToDiagnostic(err, artifact)`,
> `ErrConfigValueInvalid`-Sentinel-Wrap für den
> [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)-URL-Reject-Pfad (Spec §720). `cliJSONEnvelope.Data`
> + `newDataEnvelope`-Konstruktor wurden aus dem Template-Slice 9/9
> vorgezogen — Template-Slice erbt das Feld nur noch. Cluster-Stand
> nach Closure: **4/9 done** (doctor, add, init, generate); offene
> 5/9: remove, up-down, logs, config, template.
>
> **DoD-Tranchen-Hashes** (alle T0-T8 + R-Runden):
>
> | Tranche / Round | Inhalt | Commit |
> | --- | --- | --- |
> | T0 — Stub | `open/`-Stub Cluster-Folge-Slice 4/9 | `fbef9b5` |
> | T0 — R1 | Pre-`next/` Review-Round-1 (5 Findings, 2 HIGH / 2 MED / 1 LOW) | `f3134bd` |
> | T0 — R2 | Pre-`next/` Review-Round-2 + Carveout-Inventarisierung | `cffda2c` |
> | T0 — R3 | Pre-`next/` Review-Round-3 (V2-Stub-Korrektur, Action-Vertrag, Diagnostic-Code-Tabelle, Envelope-Shape) | `0edcf25` |
> | T0 — R4 | Pre-`next/` Review-Round-4 (Carveout-Roadmap-Sichtbarkeit, Data-Feld-Migration, Recorder-Realität, V2-Verzeichnis-Cleanup) | `0b3e1ad` |
> | T0 — R5 | Pre-`next/` Review-Round-5 (3 MED Intra-Plan-Drift) | `b14d2e8` |
> | T0 — `next/`-Übergang | Lifecycle aus `open/` | `2e3d577` |
> | T0 — R6 | `next/` Review-Round-6 (Implementation-Reality + Spec-Coverage) | `8d7d847` |
> | T0 — R7 | `next/` Review-Round-7 (Test-Härte + Plan-Drift) | `2a6f4d4` |
> | T0 — `in-progress/`-Übergang | Lifecycle aus `next/` | `17a50f4` |
> | T1 | `cliJSONEnvelope.Data` + `newDataEnvelope` + Helper-`data any`-Param (aus Template-Slice 9/9 vorgezogen) | `bd3de20` |
> | T2 | Port-Types: `GenerateRequest.PreviewMode` + `GenerateResponse.PlannedFiles`/`Changes` | `96edf40` |
> | T3 | Application-Layer: `GenerateService.fsFactory` + `generateMu` + Multi-`%w` (~17 Stellen) + `ErrConfigValueInvalid`-Wrap | `242690f` |
> | T4 | Composition-Root: `generateFSFactory`-Closure in `cmd/uboot/main.go` | `41f0231` |
> | T5 | CLI-RunE: drei JSON-Pfade + Allowlist-Migration + `mapGenerateErrorToDiagnostic(err, artifact)` + `generateEnvelopeData` | `9f8937d` |
> | T6 | 15 Acceptance-Pins — JSON-Pfade + Error-Scenarios + Action-Discriminator | `b0b31e0` |
> | T7 — R10 R1+R2 | Switch-Order auf Plan-T0-(e) + `manualConflictCodeFor` explizit | `031fd79` |
> | T7 — R10 LOW-Bundle | `unparam`-Suppression + Doc-Drifts + T6-Coverage-Text | `1e20c87` |
> | T7 — Doku-Closure | Review-Round-10-Tabelle | `a6a7d2f` |
> | T8 — Closure | CHANGELOG + `cli-json-output.md` §6/§6.5/§7 + Template-Slice-Plan-Update + roadmap-Update (4/9) + done/-Move | dieser Commit |
>
> Konsumiert das Pattern-Vorbild aus
> [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)
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
> Plan [`slice-v1-cli-cleanup-add-preview-mode-alias`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md) und wartet
> explizit auf „mindestens einen weiteren Folge-Slice" — generate
> ist genau dieser Folge-Slice, und MUSS deshalb `driving.PreviewMode`
> direkt referenzieren (kein neuer `GeneratePreviewMode`-Alias).

## Auslöser

Cluster-Slice §T0-Outcomes (a)+(b)+(e) machen jeden modifying-
Subcommand für `--json`/`--dry-run`/`--diff` verbindlich
([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) §1813, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run) §326, [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe)
§451-489). `u-boot generate <artifact>` ist nach `doctor`, `add`
und `init` der vierte Subcommand — und der erste, der
**mehrere Artefakte** über einen einzigen Subcommand bedient. Die
JSON-Surface muss pro Artefakt die spec-konforme Vorschau plus
die Action-Klassifikation (`Created`/`UpdatedBlock`/`NoOp`/
`RepairedManual`) tragen, ohne das Subcommand selbst zu zerlegen.

Spec-Bezug (geerbt von init/add):

- [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run) (Dry-Run, Voll-Schema §326)
- [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) (Diff, §451-489)
- [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) (Minimalkontrakt §1841)

Generate-spezifische Spec-Anker:

- [`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001-generate-befehl) (Subcommand-Surface mit unbekanntem Artefakt
  → Exit 2)
- [`LH-FA-GEN-002`](../../../../spec/lastenheft.md#lh-fa-gen-002-changelog-erzeugen) (changelog / [`LH-AK-007`](../../../../spec/lastenheft.md#lh-ak-007-changelog-generator) Keep-a-Changelog)
- [`LH-FA-GEN-003`](../../../../spec/lastenheft.md#lh-fa-gen-003-readme-erzeugen) (readme)
- [`LH-FA-GEN-004`](../../../../spec/lastenheft.md#lh-fa-gen-004-beispiel-env-erzeugen) (env-example)
- [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen) (devcontainer)
- [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features) (`--allow-external-feature-sources` Side-Effect)
- [`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005-idempotenz) (Idempotenz / NoOp)
- [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) (FS-Failure-Klasse, geerbt für Mid-Write
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

Tabelle listet **Write/Mkdir**-Pfade (Recorder-relevant). Zusätzlich
existieren ~10 **Read/Exists/Marshal**-Stellen mit
`ErrGenerateFileSystem`-Wrap (Z. 138, 283-284, 304-305, 498-499,
523-524, 716-717, 729-731, 784-785, 791-794, 918-921, 947-949 —
R6-Audit). Recorder zeichnet Reads NICHT auf (immutable Operation),
aber T3 muss ALLE FS-Wrap-Stellen auf Multi-`%w` migrieren — sonst
ist die Switch-Order-Garantie löchrig (T0-(d)).

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
  `writeDiff`) wird im T5 um einen `data any`-Param erweitert
  (R7-Korrektur, T0-(q)) — generate reicht
  `data: {"artifact": "<…>"}` durch, init/add reichen `nil`
  durch (nicht-brechende Trailing-Param-Erweiterung). Subcommand-
  Feld bleibt wie bisher unbenutzt — generate setzt nur
  `command`, kein `subcommand`.
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
- ✅ **NoOp-Pin (Single-Call)**: `--dry-run --json` bei bereits
  idempotenter Datei liefert `plannedFiles: []` UND `changes:
  []` (beide Arrays leer), `data.action: "no-op"`, `status: ok`,
  Exit 0.
- ✅ **Repeat-Idempotency-Pin** ([`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005-idempotenz) §1203-1213
  Wiederholungs-Eigenschaft, Port-Vertrag generate.go:171-174
  „calling `Generate` twice with the same request is safe"):
  mindestens für **changelog** (wegen Hash-Heuristik-Fragilität
  in `generate.go:249-265`) und **devcontainer** (wegen
  Two-Phase): zweimaliger Aufruf hintereinander; zweiter Lauf
  liefert `Action: NoOp`, `changes: []`, **0** Recorder-
  Mutations-Records (Spy verifiziert), `data.action: "no-op"`
  im JSON-Envelope. Single-Call-NoOp prüft nur Pre-existing-
  Idempotenz; Repeat-Call prüft echte Wiederholungs-Sicherheit
  über den vollen Use-Case-Pfad.
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
  "[`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)"`** (Devcontainer-Render-Spec, NICHT
  [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) — der ist Default-Fallback und würde Drift
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
  mit attempt-Content); diagnostics[].code: "[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)",
  `diagnostics[].file: "<File 2-Pfad>"`, `exitCode: 14`.
  Konsistent mit init's Mid-Write-Failure-Pattern (init T6
  pinnt das genauso).
  **Error-Envelope `data`-Form** (R6-Festzurrung, T0-(q) Sub-
  Decision): der Error-Envelope trägt `data: {"artifact":
  "<changelog|readme|env-example|devcontainer>"}` (ableitbar aus
  `req.Artifact`), aber **kein `data.action`**-Feld (Use-Case-
  Response auf Error-Pfad ist Zero-`GenerateResponse`, keine
  Action existiert). T5-Pflicht: `writeErrorEnvelope`/
  `reportError` werden für generate auf den neuen
  `newDataEnvelope`-Konstruktor umgestellt und tragen das
  Artefakt-`data` durch — entweder via neuer Signatur (`data
  any`-Param) oder via generate-spezifischem Helper-Wrapper, der
  vor dem Aufruf das Artefakt-`data` setzt. Sub-Decision (e) im
  T5-Entwurf konkretisiert die Variante (Signatur-Erweiterung
  vs. lokaler Wrapper).
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
- **T0-(d)** `ErrGenerateFileSystem`-Wrap-Audit (R6-Korrektur):
  heute Single-`%w` an **~17 FS-Stellen** in `generate.go`,
  nicht 8 — der ursprüngliche Stub zählte nur die Write/Mkdir-
  Pfade. Realer Wrap-Inventar:
  - **Write/Mkdir**: Z. 289/344 (changelog), Z. 504/563
    (managedFile), Z. 852/864 (devcontainer), Z. 951 (yaml-
    allowlist), Z. 848 (MkdirAll devcontainer-Dir) — 8 Stellen.
  - **Read/Exists**: Z. 138 (readProjectConfig), Z. 283-284 +
    Z. 304-305 (changelog Exists/Read), Z. 498-499 +
    Z. 523-524 (managedFile), Z. 716-717 + Z. 729-731
    (collectPorts compose), Z. 784-785 + Z. 791-794
    (devcontainer), Z. 918-921 (yaml-allowlist Read),
    Z. 947-949 (yaml Marshal) — ~10 Stellen.
  - **Wrap-Form heute**: `fmt.Errorf("%w: <op>: %v",
    driving.ErrGenerateFileSystem, err)` — Sentinel-`%w` plus
    raw-`%v` (kein Multi-`%w`), `errors.Is(err, raw)` würde
    NICHT matchen. T3 zieht auf Multi-`%w` nach (analog init's
    `initproject.go:925/967/1015/1117/1143`-Stellen mit zwei
    `%w`-Verbs).

  T3 muss alle ~17 Stellen migrieren, nicht nur die 8 Write/
  Mkdir. T6 erweitert die Switch-Order-Pflicht-Pin-Tests um
  **mindestens einen Read-Pfad-FS-Failure-Test** (z. B. `ReadFile`-
  Permission-Denied bei `applyAllowExternalFeatureSources`),
  sonst ist die FS-first-Garantie löchrig.
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
  | `ErrGenerateFileSystem` | * (alle) | [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) | 14 |
  | `ErrGenerateManualConflict` | changelog | [`LH-FA-GEN-002`](../../../../spec/lastenheft.md#lh-fa-gen-002-changelog-erzeugen) | 10 |
  | `ErrGenerateManualConflict` | readme | [`LH-FA-GEN-003`](../../../../spec/lastenheft.md#lh-fa-gen-003-readme-erzeugen) | 10 |
  | `ErrGenerateManualConflict` | env-example | [`LH-FA-GEN-004`](../../../../spec/lastenheft.md#lh-fa-gen-004-beispiel-env-erzeugen) | 10 |
  | `ErrGenerateManualConflict` | devcontainer | [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen) | 10 |
  | `ErrConfigValueInvalid` (`--allow-external-feature-sources` URL) | devcontainer | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features) | 10 |
  | `ErrArtifactUnknown` | * (alle) | [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) | 2 |
  | Default (unknown) | * (alle) | [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) | 1 |

  Switch-Order verbindlich: FS-first, dann ManualConflict
  (artefakt-spezifischer Code via `artifact`-Param), dann
  ConfigValueInvalid (devcontainer-only, `--allow-external-
  feature-sources`-URL-Reject — Spec §720 fordert exakt
  [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)/Exit 10), dann ArtifactUnknown, dann Default.
  **T3-Pflicht**: `validateAllowExternalFeatureSourcesEntries`
  (`generate.go:893-901`) und `applyAllowExternalFeatureSources`
  (`generate.go:912-934`) wrappen heute mit
  `fmt.Errorf("generate devcontainer: …: %w", err)` ohne typed
  Sentinel — Code-Kommentar in Z. 908-911 verspricht
  `ErrConfigValueInvalid` aber wrappt ihn nicht. T3 zieht den
  Sentinel-Wrap nach (Multi-`%w` analog FS-Wraps).
  **T6-Pflicht**: Reject-Pin-Test für `--allow-external-feature-
  sources <invalid-url> --json` → `code: "[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)"`,
  `exitCode: 10`.
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
      ([`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)).
  (b) Phase-2-Mid-Write-Failure → File 1 committed auf Disk,
      File 2 underlying-Write fehlgeschlagen. Envelope-Form
      (R4-Recorder-Realität): `plannedFiles[]` enthält **beide**
      Captures (File 1 mit success-Content, File 2 mit
      attempt-Content — Recorder zeichnet vor Delegieren auf,
      `recordingfs.go:139`), `lastPlannedPath` liefert File 2
      als `diagnostics[].file`, `diagnostics[].code:
      "[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)"`, Exit 14. Konsistent mit init T6
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
  bei Wahl von `subcommand="<artifact>"` wäre eine zusätzliche
  Signatur-Änderung nötig. (Der `data any`-Param wird ohnehin
  via T0-(q)/R7 nachgezogen — das ist nicht-brechend; eine
  parallele `subcommand`-Erweiterung wäre dagegen semantisch
  schiefe Symmetrie.)
  Plus: Action-Klassifikation lebt ohnehin in `data.action`
  (T0-(f)), so dass `data` der natürliche Träger für
  Generate-spezifische Felder ist.
- **T0-(n)** **`Codes`-Registry-Ergänzung NICHT nötig**: die
  Codes-Map in `jsontestutil.DefaultAllowedCodes` (cli-json-output.md
  §5) ist für tool-interne **dotted** Codes (`add.*`, `init.*`-
  Style) gedacht — LH-Codes werden vom Acceptance-Helper bereits
  generisch erlaubt ([`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001-generate-befehl)..[`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005-idempotenz)/[`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)/[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)/
  [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) sind alle LH-präfixiert). T6
  ergänzt keine Registry-Einträge; Doku-Update für §6.5
  (Per-Command-Sektion) genügt.
- **T0-(o)** Pre-`next/`-Review-Runden-Erwartung: init hatte
  3 vor `next/`, add hatte 5. Generate steht nach R6 bei sechs
  (R1-R6, 26 Findings); Konvergenz erreicht.
- **T0-(q)** **Error-Envelope-`data`-Vertrag festgezurrt** (R6-
  Finding 2): Success-Envelope trägt `data: {"artifact":
  "<…>", "action": "<…>"}` (Top-Level-`data`-Feld, R3-
  Festzurrung). Error-Envelope (Mid-Write-Failure, Validation-
  Conflict, URL-Reject) trägt `data: {"artifact": "<…>"}`
  **ohne `action`** — Use-Case-Response auf Error-Pfad ist
  Zero-`GenerateResponse`, Action existiert nicht. Konsumenten
  können daher Generate-Action im Erfolgsfall lesen
  (`data.action`), im Fehlerfall nur das Artefakt
  (`data.artifact`). Symmetrie zu init/add: dort trägt der
  Error-Envelope kein `data` (init/add haben kein `data`-Feld);
  generate ist hier abweichend, weil `data.artifact` für
  multi-artifact-Kontext load-bearing ist (sonst wüssten
  Konsumenten nicht, welches der vier Artefakte gefailt hat).
  T5-Sub-Decision: entweder `writeErrorEnvelope`-Signatur um
  `data any`-Param erweitern (zieht init/add mit hoch, gleicher
  Constructor-Pfad), oder generate-lokaler Wrapper. Vorschlag:
  Signatur-Erweiterung — niedriges Drift-Risiko, weil init/add
  schlicht `nil` durchreichen.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | ------ | --------------- | --- |
| T0 | Discovery + Sub-Decisions (a)-(o) klären; Review-Runden | — (Plan) | — |
| T1 | Refactor-Tranche (wenn überhaupt nötig — generate hat schmalere FS-Surface; ggf. nur ErrGenerateFileSystem-Multi-`%w`-Audit oder gar kein T1) | ~30-80 | T0 |
| T2 | Port-Types: `GenerateRequest.PreviewMode`, `GenerateResponse.PlannedFiles`/`Changes`-Felder. **`data.action`-Klassifikation** liegt im Envelope-Layer (T5), nicht im Port — die existierende `GenerateResponse.Action` (`GenerateAction`-Enum) wird in T5 zum `data.action`-String gerendert; keine neuen Port-Felder dafür (T0-(f) Festzurrung). `ErrGenerateFileSystem` ist schon da. | ~50 | T0 |
| T3 | Application-Layer: `GenerateService.fsFactory` + `generateMu sync.Mutex` + `NewGenerateServiceWithFactory` + `Generate()`-Wrapper mit FS-Swap; `mapCaptureToPlannedFiles(captured, req.BaseDir)`; Multi-`%w`-Wrap an den **~17 FS-Wrap-Stellen** (T0-(d) Audit, R6-Kalibrierung — Write/Mkdir + Read/Exists/Marshal); `ErrConfigValueInvalid`-Sentinel-Wrap für `--allow-external-feature-sources`-URL-Reject (T0-(e) [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)-Pfad). | ~280 | T2 |
| T4 | Composition-Root-Wiring `generateFSFactory`-Closure in `cmd/uboot/main.go`. | ~30 | T3 |
| T5 | CLI-RunE: `runGenerate` ruft generische Helper mit `command="generate"` (kein subcommand, T0-(m)), `mapErr=mapGenerateErrorToDiagnostic`; drei JSON-Pfade; Allowlist-Migration; **`mapGenerateErrorToDiagnostic(err, artifact)` neu mit Artefakt-Parameter** (T0-(e); per-Artefakt LH-Code für ErrGenerateManualConflict). `data.action` aus `resp.Action.String()` gerendert; `data.artifact` aus `req.Artifact.String()`. **Helper-Signatur-Erweiterung (R7-Korrektur, T0-(q))**: `writeErrorEnvelope`/`reportError` werden um einen `data any`-Param erweitert (Default `nil`), damit Generate `data: {"artifact": "<…>"}` auch im Error-Envelope durchreichen kann (Mid-Write-Failure, ManualConflict, URL-Reject). Init/add-RunE-Stellen reichen `nil` durch — Drift-Risiko gering, weil die Signatur-Änderung nicht-brechend ist (neuer Trailing-Param). | ~240 | T1 + T2 (T4 für Run-time-Smoke, Code-parallelisierbar) |
| T6 | Acceptance-Tests: **repräsentative Pin-Coverage** statt 8×4-Matrix-Durchlauf (R10-LOW-2-Klarstellung) — pro Flag-Modus ein Pin auf einem repräsentativen Artefakt, pro Artefakt mindestens ein Pin (ManualConflict-Sub-Tests decken die 4-Artefakt-Symmetrie). Pflicht-Pins: drei JSON-Modi (Bare/DryRun/Diff/Combo); per-Artefakt ManualConflict-Codes; URL-Reject ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)); ArtifactUnknown/ProjectNotInitialized/FS-Failure-Exit-Codes; Allow-External-Mutex; NoOp-Empty-Arrays; UpdatedBlock-vs-RepairedManual-Action-Diskriminator; Human-Mode-Summary + Diff-Rendering. Devcontainer-Phase-1/Phase-2-Half-Write und Repeat-Idempotency leben im Application-Layer (`generate_test.go`/`generate_features_test.go`) — CLI-Stub-Tests können FS-Phasen nicht simulieren. Total 15 Tests (~530 LOC realistic). | ~530 | T5 |
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
| 1 | HIGH | Phase-2-Half-Write-Carveout war im Slice angekündigt, aber NICHT in `carveouts.md` inventarisiert und ohne Plan-Stub (verletzt [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin) und MEMORY-Rule [[feedback_carveouts_need_plans]]) | Neuer Eintrag in `carveouts.md` §Temporäre Carveouts mit Re-Trigger; neuer open/-Plan-Stub [`slice-v2-generate-devcontainer-rollback-aware-write`](../open/slice-v2-generate-devcontainer-rollback-aware-write.md) (Status `on hold pending trigger`, drei Lösungs-Skizzen) |
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
| 3 | MEDIUM | Diagnostic-Code zu lose gepinnt ([[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) oder [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)]); ErrGenerateManualConflict braucht Artefakt-Kontext, geplanter `mapGenerateErrorToDiagnostic(err)` hat keinen | T0-(e) erweitert um **Mapper-Signatur `(err, artifact)`** mit Begründung und um eine **per-Artefakt LH-Code-Tabelle** (changelog→GEN-002, readme→GEN-003, env-example→GEN-004, devcontainer→DEV-001). Acceptance-Pin auf exakten [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen) umgestellt (statt „either/or"). T5-Tranche um Mapper-Signatur-Erweiterung ergänzt |
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
| 1 | HIGH | V2-Slice fehlte als Roadmap-Slice-Zeile ([`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin) fordert Doppel-Sichtbarkeit in `carveouts.md` UND `roadmap.md`) | Neuer Roadmap-Eintrag in `roadmap.md` §AP-Tabelle für [`slice-v2-generate-devcontainer-rollback-aware-write`](../open/slice-v2-generate-devcontainer-rollback-aware-write.md) mit Status `on hold pending trigger` und Carveout-Anker-Verweis |
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

## Review-Round-6 (`next/`, Implementation-Reality + Spec-Coverage)

Sechste Runde gegen den nach `next/`-Übergang konsolidierten Stub
(`2e3d577`), Angle: **Implementation-Reality-Audit** (jeder
proposierte T1-T8-Inhalt gegen existierende Codebase belegt) +
**Spec-Coverage-Audit** (jeder LH-Anker gegen `lastenheft.md`)
+ **Geerbte-Pattern-Check** (init-Vergleiche gegen done-File).
Fünf Findings (2 HIGH, 2 MEDIUM, 1 LOW gegen V2-Stub):

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)-URL-Reject-Pfad nicht im Diagnostic-Code-Mapping abgebildet — Spec §720 fordert exakt [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)/Exit 10 bei ungültiger `--allow-external-feature-sources`-URL; heutiger Code (generate.go:898-933) wrappt ohne typed Sentinel; T0-(e)-Tabelle hatte keinen Eintrag → Acceptance-Test würde auf Default-Branch ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes)/Exit 1) fallen | T0-(e)-Tabelle um Zeile ErrConfigValueInvalid | devcontainer | [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features) | 10 erweitert; Switch-Order-Block erweitert; T3-Pflicht für `ErrConfigValueInvalid`-Sentinel-Wrap nachgezogen; T6 um Reject-Pin-Test ergänzt; T3-LOC angehoben |
| 2 | HIGH | `data.action` im Error-Envelope undefiniert — AK pinnte `data.action` für Success, Phase-2-Half-Write-AK schwieg; `writeErrorEnvelope` setzt heute `subcommand=""` hardcoded und hat keinen Daten-Slot → Vertragsambiguität | Phase-2-Half-Write-AK um Error-Envelope-`data`-Klärung ergänzt (`data.artifact` JA, `data.action` NEIN — Zero-Response auf Error-Pfad); **neue T0-(q)** Sub-Decision für volle Symmetrie-Klärung; T5-Pflicht: `writeErrorEnvelope`-Signatur um `data any`-Param erweitern (zieht init/add mit hoch — `nil` durchreichen) |
| 3 | MEDIUM | [`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005-idempotenz)-Idempotenz nur Single-Call gepinnt; Spec §1203-1213 + Port-Vertrag generate.go:171-174 fordern Wiederholungs-Eigenschaft | Repeat-Idempotency-Pin in AK ergänzt (mindestens changelog wegen Hash-Heuristik + devcontainer wegen Two-Phase); zweiter Aufruf → `NoOp`, 0 Recorder-Mutations-Records (Spy), `data.action: "no-op"`. T6-Zelle um den Pin erweitert; Test-Total ~40-42 |
| 4 | MEDIUM | T0-(d)-Wrap-Audit zählte 8 Stellen, real ~17 (Read-Pfade fehlten); Wrap-Form ist `%w: …: %v` (Single-`%w`) — kein Multi-`%w`-Pattern | T0-(d) auf ~17 Stellen kalibriert mit Code-Anker-Inventar (Write/Mkdir + Read/Exists/Marshal); Pre-Scan-Tabelle um Read-Pfad-Notiz; T3-LOC von ~200 auf ~280 angehoben; T6 um Read-Pfad-FS-Failure-Pin (mindestens einer) |
| 5 | LOW (V2) | V2-Side-Effect-Liste übersah `collectDevcontainerForwardPorts`-Pre-Read-Sequenz; heutiger Flow ist konsistent (Reads laufen vor MkdirAll), aber V2-Stub nicht zukunftsfest gegen Schema-Erweiterung | V2-Stub um Sequenz-Reihenfolge-Notiz (heutige `generate.go:636-672`-Sequenz) + Trigger-Zukunftsfestigkeit-Hinweis ergänzt |

R6-Reviewer-Note: docs-check grün. Implementation-Reality-Audit
hat die Helper-Pattern (`mapCaptureToPlannedFiles`,
`previewModeFromFlags`, Recorder-vor-Delegieren, init's Multi-`%w`)
alle als real existierend bestätigt — Pattern-Erbe-Behauptungen
tragen. Die Lücken lagen an den **Vertrags-Rändern** (neuer
Spec-Anker nicht in Mapping; Error-Envelope-Symmetrie nicht
durchspezifiziert) und an einer **Audit-Untererfassung** (Read-
Pfade in Wrap-Inventar fehlten). Keine fundamentalen Plan-
Reorganisationen nötig.

## Review-Round-7 (`next/`, Test-Härte + Plan-Drift)

Siebte Runde gegen den R6-gepflegten Stub (`8d7d847`), Fokus
Test-Realismus und Plan-Drift nach den R6-Vertragserweiterungen.
Zwei Findings (1 HIGH gegen V2-Stub, 1 MEDIUM gegen V1-T5):

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | HIGH (V2) | YAML-Rollback-Pin B kann mit no-op-failing Spy bestanden werden, ohne den echten Restore-Pfad zu belegen. Echter `FS.WriteFile` (`fs.go:56-61`) delegiert auf `os.WriteFile` → truncate-overwrite mit anschließend möglichem Failure (Disk-Full nach Truncate, Signal-Interrupt) ist realistisch. Ohne explizite Forderung an den Fake bleibt der YAML-Restore-Pfad ungetestet. | V2-Stub Pin B explizit: Spy muss vor dem Failure eine **partielle/truncated Mutation** auf u-boot.yaml ausführen. Hash-Snapshot pre/post Pflicht (nicht nur tree-Vergleich); Post-Hash IDENTISCH zum Pre-Hash → Rollback hat den truncierten Zustand restauriert, nicht am Original vorbeigelaufen |
| 2 | MEDIUM | T5-Tranchen-Zelle sagte „Helper-Generalisierung bleibt unverändert" — widerspricht T0-(q) und AK-Block, die `writeErrorEnvelope`/`reportError` um `data any`-Param für `data.artifact` im Error-Envelope erweitern fordern. Plan-interne Drift nach R6-Erweiterung. | T5-Zelle nachgezogen: Helper-Signatur-Erweiterung um `data any`-Param explizit (init/add reichen `nil` durch, nicht-brechende Trailing-Param-Erweiterung); LOC ~200 → ~240. AK-Zeile (R3-Festzurrung) und T0-(m)-Drift-Risk-Notiz auf die nicht-brechende Form aktualisiert |

R7-Reviewer-Note: docs-check grün. Keine neuen Vertragslücken;
R7 räumt Plan-interne Drift nach R6 auf und härtet einen
Acceptance-Pin gegen no-op-passing-Tests. Stub konvergiert
weiterhin.

## Review-Round-10 (T7 Pre-T8 Code-Review)

Code-Review nach T1-T6 (Diff-Range `bd3de20..b0b31e0`), Fokus auf
Code-Korrektheit gegen die Plan-Versprechen, Pattern-Erbe-
Konsistenz init→generate und potentielle Inkonsistenzen zwischen
den Tranchen. Fünf Findings (2 MEDIUM, 3 LOW), alle adressiert:

| # | Sev | Finding | Adressierung |
| - | --- | --- | --- |
| 1 | MEDIUM | Switch-Order in `mapGenerateErrorToDiagnostic` war `FS → ConfigValueInvalid → ManualConflict → …` — Plan-T0-(e) pinnt aber `FS → ManualConflict → ConfigValueInvalid → …`. Operational heute kein Bug (Sentinels in getrennten Pfaden), aber Plan-Drift gegen verbindliche Sub-Decision | Switch-Order auf Plan-Konformität gezogen (Commit `031fd79`); Doc-Kommentar dokumentiert die Plan-Reihenfolge mit Begründung |
| 2 | MEDIUM | `manualConflictCodeFor` hatte impliziten Fall-Through auf [`LH-FA-GEN-002`](../../../../spec/lastenheft.md#lh-fa-gen-002-changelog-erzeugen) (changelog) — Plan-T0-(e) Default-Zeile sagt aber [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes). Zero-Value Artifact wirkte zufällig korrekt (zero=changelog), aber Enum-Erweiterung ohne Switch-Update würde Code stillschweigend falsch routen | Alle vier Artefakte explizit als case gelistet (inkl. changelog); Default-Fallback auf [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006-exit-codes) — macht Switch-Lücken nach Enum-Erweiterung sichtbar (Commit `031fd79`) |
| 3 | LOW | `//nolint:unparam` auf `reportError.data` mit "T5 ist der zweite Caller"-Begründung — nach T5-Merge ist generate der reale zweite Caller, unparam akzeptiert die Mischform | Suppression entfernt; Doc auf etablierte Realität umgestellt (Commit `1e20c87`) |
| 4 | LOW | T6-Plan-Text suggerierte strikte 4×8=32-Matrix; T6 implementiert 15 repräsentative Pins. Devcontainer-Phase-1/Phase-2 + Repeat-Idempotency leben im Application-Layer | Plan-T6-Zelle auf "repräsentative Pin-Coverage" präzisiert; Bezug auf Application-Layer-Tests; LOC-Schätzung 680→530 (Commit `1e20c87`) |
| 5 | LOW | Doc-Drift in `GenerateResponse.PlannedFiles`-Kommentar — verwies auf JSON-Tags am Slice statt am PlannedFile-Type; Port serialisiert nichts direkt | Kommentar nachgezogen: Verweis auf `PlannedFile`-Definition + Klarstellung dass CLI eigene Wire-Typen baut (Commit `1e20c87`) |

R10-Reviewer-Note: docs-check + lint + test + coverage grün
(91.10 %). Keine HIGH-Befunde — Pattern-Erbe init→generate ist
strukturell sauber (`Generate()`-Wrapper 1:1 zu `Init()`-Wrapper,
Multi-`%w`-Migration vollständig auf den ~17 FS-Stellen,
`ErrConfigValueInvalid`-Sentinel-Wrap für URL-Reject etabliert,
Allowlist-Migration konsistent). Lücken nur an Switch-Order-
Reihenfolge, Defensiv-Fallbacks und Doku-Drift.

## Out of Scope

- **Roll-back-aware Recorder** für Devcontainer-Phase-2-Atomicity
  (Cluster-T0-(b) Variante 3 verworfen, V2-Scope). Phase 2-Mid-
  Write hinterlässt halbgeschriebenes File 1 auf Disk; Carveout-
  Eintrag in `carveouts.md` mit Re-Trigger auf einen späteren
  Devcontainer-Rollback-Slice.
- **HTTP- oder gRPC-Schnittstellen**: [ADR-0010](../../adr/0010-kein-http-driving-adapter.md) schließt
  explizit aus.
- **Schema-Versionierung** (`schemaVersion: 1`): siehe
  Cluster-Slice §Out of Scope.
- **Neue Artefakte** (z. B. `generate dockerfile` standalone):
  V1-Scope ist die existierenden vier.
- **Generic `mapErrorToDiagnostic`-Registry**: Altitude-Reviewer-
  Vorschlag aus add R6 #I1. Cluster-T_close-Aufgabe.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
  §T0-Outcomes — Vorgaben für den Folge-Slice-Block.
- Pattern-Vorbild:
  [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)
  — T0-T8 + Review-Round-9 voll abgeschlossen. Erbschafts-
  Disziplin in T0-(a) dieses Slices; Alias-Lebensdauer-Pflicht
  aus init-T0-(c) zwingt direktes `driving.PreviewMode`.
- Add-Slice (sekundär):
  [`slice-v1-cli-json-dry-run-add`](slice-v1-cli-json-dry-run-add.md)
  — Pattern-Founder; relevante Sub-Decisions
  (`CountAdditions`-Semantik §477, `checkHunks`-Helper) bleiben
  geerbt.
- Cleanup-Stub:
  [`slice-v1-cli-cleanup-add-preview-mode-alias`](../open/slice-v1-cli-cleanup-add-preview-mode-alias.md)
  — wartet auf „mindestens einen weiteren Folge-Slice"; generate
  IST dieser Folge-Slice, MUSS `driving.PreviewMode` direkt
  nutzen.
- Spec: [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe),
  [`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001-generate-befehl)..[`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005-idempotenz), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)/[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003-devcontainer-features), [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- Code-Anker heute:
  [`generate.go`](../../../../internal/hexagon/application/generate.go)
  (~960 LOC, vier Artefakt-Handler + managedblock-Surface),
  [`cli/generate.go`](../../../../internal/adapter/driving/cli/generate.go)
  (~160 LOC, RunE-Erweiterungs-Ziel),
  [`port/driving/generate.go`](../../../../internal/hexagon/port/driving/generate.go)
  (Carrier-Types + drei Sentinels inkl. `ErrGenerateFileSystem`).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
