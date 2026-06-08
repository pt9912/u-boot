# Slice V1: `config --json` / `config get --json` / `config set --json` — drei Sub-Forms unter einem Folge-Slice

> **Status:** `done/` — **alle Tranchen T0–T7 + drei Review-Runden
> abgeschlossen, T8-Closure (2026-06-08)**, `make gates` grün
> (Coverage 91.30 %).
>
> **DoD-Tranchen-Hashes** (T0–T8 + Review-Runden):
>
> | Tranche / Round | Inhalt | Commit |
> | --- | --- | --- |
> | T0 — Stub | Discovery-Stub für 8/9 config in `open/` | `da3c7f0` |
> | T0 — R1 | adversarialer Stub-Review (4 HIGH + 6 MED + 4 LOW) | `13b76d3` |
> | T0 — R2 | adversarialer Stub-Review (3 HIGH + 6 MED + 4 LOW) | `0990a9d` |
> | T0 — R3 | adversarialer Stub-Review (0 HIGH + 4 MED + 4 LOW; Asymptote) + Lifecycle `open/`→`next/` | `e813b05` |
> | T2 | Port-Felder (`PreviewMode`/`SilenceLogger`/`Warnings`) + 2 Sentinels (`ErrConfigWriteRejected`/`ErrConfigPostPatchSanityFailed`) + `cli.ErrDryRunNotApplicable` + `configSetFlags.JSON`/`Quiet` + Pin-Tests + Lifecycle `next/`→`in-progress/` | `acbd97d` |
> | T3 | Application-Layer: Multi-`%w`-FS-Wraps + Sentinel-Split-Wiring + SilenceLogger-Branch + Orphan-WARN→`Warnings` Dual-Emission | `e29fd1e` |
> | T4 | PreviewMode-Cluster: `fsFactory` + `selectFS` + Write-Routing + `NewConfigServiceWithFactory` + `cmd/uboot/main.go`-Wiring (fünfter Preview-Factory) | `a2aafc8` |
> | T4 — Review R-T4-1 (HIGH) | Selbst-Review: `ConfigSetResponse.PlannedFiles`-Feld + Recorder-Surfacing (sonst `--diff` ohne Byte-Quelle) | `1ad1344` |
> | Review R-IR-1 (HIGH) | unabhängiger Agent (Port/App): zwei Split-Sentinels fehlten in `isConfigValidationError` → Exit-10→Exit-1-Regression gefixt + ExitCode-Pin | `ebd5119` |
> | T5 | CLI-RunE-Neufassung: 3 Data-Carrier + Subcommand-Pflicht (inkl. Error-Pfad via `reportErrorSub`) + Allowlist 3 Forms + Mapper Switch-Order T0-(f) + `configArgsValidator` + Voll-Schema/Dry-Run/Diff + Reject + WARN→diagnostics + Acceptance-Tests | `11aea03` |
> | T6 | Acceptance-Vervollständigung: white-box Mapper-Rows + Switch-Order-`_ByDesign` + Cobra-unknown-sub + Help-Edge-Case + CONF-005-Disambiguation + Sanitizer + Subcommand-Pflicht + `--quiet --json` + Mid-Stage-Shapes | `4cb7e90` |
> | T7 — Review R-CLI-1 (MED) | unabhängiger Agent (CLI-Layer): Args-Validator-Voll-Schema-Leak auf Read-only-Forms gefixt + 2 Regression-Pins + Diff-Hunks-Assertion verschärft | `5cca4d5` |
> | T8 — Closure | CHANGELOG + `cli-json-output.md` §6.9/§6-Tabelle/§7 + roadmap done-Zähler 7→8 + 4 carveouts-Einträge + `done/`-Move | `724a283` |
>
> T2–T7 geliefert:
> Port-Felder + zwei Sentinels
> `ErrConfigWriteRejected`/`ErrConfigPostPatchSanityFailed`
> (T0-(m)-Split) + `cli.ErrDryRunNotApplicable` +
> `configSetFlags.JSON`/`Quiet` (T2); Multi-`%w`-FS-Wraps +
> Sentinel-Split-Wiring + SilenceLogger-Branch + Orphan-
> WARN→`Warnings` Dual-Emission (T3); `fsFactory`-Feld + nil-safe
> `selectFS` + Write-Routing über `selectFS(req.PreviewMode)` +
> `NewConfigServiceWithFactory` + `cmd/uboot/main.go`-Wiring (T4);
> vollständige `config.go`-CLI-Neufassung — 3 Data-Carrier,
> Subcommand-Pflicht (auch Error-Pfad via subcommand-bewusste
> `reportErrorSub`/`writeErrorEnvelopeSub`), Allowlist 3 Forms
> (Reject 4→1), `mapConfigErrorToDiagnostic` (Switch-Order T0-(f)),
> `configArgsValidator`, Voll-Schema/Dry-Run/Diff, bare/get-Reject
> via `ErrDryRunNotApplicable`+`isUsageError`, WARN→diagnostics,
> Acceptance-Tests (T5); T6 ergänzte Mapper-Rows + Switch-Order-
> `_ByDesign` (white-box) + `--quiet --json`-Pins + Cobra-unknown-
> sub-Pin + Help-Edge-Case + CONF-005-Disambiguation + Sanitizer-
> Worst-Case + Subcommand-Pflicht-Pin + Mid-Stage-Shapes (T6); drei
> Review-Runden (R-T4-1 + R-IR-1 + R-CLI-1, alle gefixt) (T7).
> **Review-Runde vor T5 (2026-06-08, zwei Stufen)**: (1) Selbst-
> Review fand **R-T4-1 (HIGH)** — `ConfigSetResponse.PlannedFiles
> []driving.PlannedFile`-Feld neu + Recorder-Surfacing
> (`mapCaptureToPlannedFiles`), weil `config set --diff` die
> patched/current Bytes für den geteilten Diff-Renderer braucht
> (T4-Recorder-Verzicht war falsch). (2) **Unabhängiger** Reviewer-
> Agent fand **R-IR-1 (HIGH)** — die T3-Sentinel-Split-Sentinels
> `ErrConfigWriteRejected`/`ErrConfigPostPatchSanityFailed` fehlten
> in `cli.go isConfigValidationError` (ein von der T5-Mapper-Tabelle
> **unabhängiger**, bereits live wirksamer ExitCode-Klassifikator) →
> heutige Exit-10→Exit-1-Regression auf dem Plain-CLI-Pfad. Gefixt +
> `TestExitCode_ConfigValidationSentinels`-Tabellen-Pin gegen künftige
> Split-Regressions. Übrige Punkte (Multi-`%w`, Read/Write-Trennung,
> SilenceLogger, kein Mutex nötig) clear.
> Die zuvor nach T5 verschobenen Flag-Structs
> (`configGetFlags`/`configShowFlags` + `DryRun`/`Diff` +
> Flag-Registrierung) sind in T5 gelandet. **T8-Closure
> abgeschlossen** (`724a283`): CHANGELOG, `cli-json-output.md`
> §6.9 + §6-Tabelle (Reject 4→1) + §7-Mutations-Zeile, roadmap
> done-Zähler 7→8, carveouts.md für die vier Folge-Stubs, Slice
> nach `done/` mit DoD-Hash-Tabelle. Slice damit vollständig done;
> Cluster-Stand 8/9 (nur noch template 9/9 offen).
> R1+R2+R3-
> Adressierung gefahren: R1=4+6+4, R2=3+6+4, R3=0+4+4 — Asymptote
> erreicht; T0-Sub-Decisions (a)-(p) komplett mit R3-festgezurrten
> Beschlüssen für (a)/(b)/(g)/(o)/(p); vier Folge-Slice-Stubs in
> `open/` gespawned; LOC-Bilanz ~1500-1900.
> Achter Folge-Slice (8/9) des Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 8/9). **Drei Sub-Forms unter einem Slice**
> (analog up-down-Bündelung): `u-boot config` (bare), `u-boot
> config get <path>`, `u-boot config set <path> <value>`. Bare +
> Get sind **Read-only-Klasse**; Set ist **Modifying-Klasse** —
> trägt zusätzlich `--dry-run` + `--diff` (Cluster-Plan Z. 91-100).
> Der Slice ist damit der **erste Folge-Slice der Read-only- und
> Modifying-Pattern in einer einzigen Stub-Lieferung bündelt**.
>
> Erbt Modifying-Klassen-Disziplin aus
> [`slice-v1-cli-json-dry-run-remove`](slice-v1-cli-json-dry-run-remove.md)
> (RecordingFileSystem + Diff-Renderer + Custom-Args-Validator +
> Sanitizer für Pre-UC-Validation). Erbt Read-only-Klassen-
> Disziplin aus
> [`slice-v1-cli-json-dry-run-up-down`](slice-v1-cli-json-dry-run-up-down.md)
> (FS-Sentinel-Pattern + Mapper-Switch-Order +
> `cli/sanitize.go`-Helper) und
> [`slice-v1-cli-json-dry-run-logs`](slice-v1-cli-json-dry-run-logs.md)
> (Single-Envelope für Read-only + `(code, exitCode)`-Tupel-
> Disambiguation Pattern). Erbt `subcommand`-Pflichtfeld-Form aus
> [`slice-v1-cli-json-dry-run-template`](slice-v1-cli-json-dry-run-template.md)-
> Stub (`command="config"` + `subcommand="list"|"get"|"set"` —
> T0-(b) Sub-Decision).

## Auslöser

Cluster-Slice §T0-Outcomes (a) macht `--json` für jeden
Spec-Enum-Subcommand verbindlich (`LH-NFA-USE-004` §1813);
Cluster-Slice Z. 102-113 **§config-Cluster-Pflicht** (Review-
Finding MEDIUM) zwingt explizit alle drei Sub-Forms in **einen**
Slice: gemeinsamer `ConfigUseCase`, gemeinsame Sentinels
(`ErrConfigPathUnknown`, `ErrConfigValueInvalid`,
`ErrConfigSchemaInvalid`, `ErrConfigFileSystem`,
`ErrConfigValueNotSet`), gemeinsamer Mapper.

`config` ist der einzige Folge-Slice der **read-only und
modifying gleichzeitig** trägt; `config set` ist die **letzte
nicht-migrierte Modifying-Surface** (add/init/generate/remove
sind alle in `done/`). Damit ist der Slice der **definitive
Test-Träger für die Modifying-Pattern-Wiederverwendung** —
keine neue Pattern-Erfindung, ausschließlich Wiederverwendung.

Spec-Bezug:

- `LH-FA-CONF-001..005` — Config-Subcommand-Vertrag.
- `LH-NFA-USE-004` §1813 / §1841 — Minimalkontrakt-Pflicht für
  alle drei Sub-Forms.
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag (für `config
  set` mit `--dry-run`/`--diff`; bare/get nutzen Minimalkontrakt).
- `LH-FA-CLI-008` §451-489 — Diff-Vertrag (für `config set
  --diff`).

Heute-Stand-Pre-Scan
([`internal/adapter/driving/cli/config.go`](../../../../internal/adapter/driving/cli/config.go),
224 LOC;
[`internal/hexagon/application/config.go`](../../../../internal/hexagon/application/config.go),
845 LOC;
[`internal/hexagon/port/driving/config.go`](../../../../internal/hexagon/port/driving/config.go),
187 LOC):

| Aspekt | bare (`config`) | `config get` | `config set` |
| --- | --- | --- | --- |
| Positional-Args | `NoArgs` | `ExactArgs(1)` (`<path>`) | `ExactArgs(2)` (`<path> <value>`) |
| Lokale Flags | — | — | `--allow-external-feature-sources` (StringSlice, LH-FA-DEV-003) |
| FS-Mutation | — | — | `WriteFile(u-boot.yaml, patched)` (`config.go:194`) |
| FS-Read | `ReadFile(u-boot.yaml)` (Show) | `ReadFile(u-boot.yaml)` + parse | `ReadFile(u-boot.yaml)` + parse |
| UC-Method | `ConfigUseCase.Show` | `ConfigUseCase.Get` | `ConfigUseCase.Set` |
| Output heute | `out.Write(resp.Body)` byte-identisch | `fmt.Fprintln(out, resp.Value)` bare-Scalar | `printConfigSetSummary` Two-Shape (`OldValue → NewValue` / `already X`) |
| Sentinels | `ErrProjectNotInitialized` (10), `ErrConfigFileSystem` (14) | `ErrConfigPathUnknown` (10), `ErrConfigValueNotSet` (10), `ErrConfigSchemaInvalid` (10), plus bare-Sentinels | alle Get-Sentinels plus `ErrConfigValueInvalid` (10), `ErrConfigSchemaInvalid` (10) |
| Allowlist heute | Reject (Z. 29 `jsonallowlist.go`) | Reject | Reject |

Use-Case-Deps `ConfigService`: `driven.FileSystem`,
`driven.YAMLCodec`, `driven.Logger`. KEIN `Confirmer`, KEIN
`DockerEngine`, KEIN `Progress`-bound-state. **Fünf Sentinels
existieren bereits typed** (`ErrConfigPathUnknown`,
`ErrConfigValueInvalid`, `ErrConfigSchemaInvalid`,
`ErrConfigFileSystem`, `ErrConfigValueNotSet`); `ErrConfigFileSystem`
existiert bereits mit Read-Message-Form (`"config: filesystem
error"`), **wird aber heute single-`%w` + `%v`-tail gewrapped**
(`fmt.Errorf("%w: read %q: %v", driving.ErrConfigFileSystem,
path, err)` an allen 5 Sites). T3 migriert die 5 Sites auf
**echtes Multi-`%w`** (`%w: ...: %w` — Pattern-Erbe up-down
T3) damit die Switch-Order-Defense-Tests (Mapper FS-first vs.
ExitCode-Helper Driven-first) gegen synthetische Multi-Wrap-
Chains greifen. (R1-HIGH-1-Adressierung.)

Bemerkenswert: anders als up/down/logs hat `config` **schon
einen FS-Sentinel** (`ErrConfigFileSystem`, Z. 141 `port/
driving/config.go`). Das bedeutet **T2 ist substanziell
kleiner als bei up-down/logs** — kein neuer Sentinel nötig,
nur Switch-Order-Disziplin im neuen Mapper. Außerdem ist
`driving.ErrConfigFileSystem` **bereits in `cli.go:405`
`isFilesystemError`** registriert (Pre-Cluster-Slice) — keine
T5-Co-Migration nötig (R1-LOW-1-Adressierung).

## Aufhebungsbedingung

Sechs Flag-Kombinationen für drei Sub-Forms liefern spec-
konforme Outputs:

```bash
u-boot config --json                                     # bare → Minimal+Data{body}
u-boot config get project.name --json                    # Get → Minimal+Data{path, value}
u-boot config set project.name x --json                  # Set → Voll-Schema (RecordingFS) mit dataEnvelope
u-boot config set project.name x --dry-run --json        # Set Dry-Run → plannedFiles[] ohne WriteFile-Call
u-boot config set project.name x --diff                  # Set Diff Human-Mode (unified)
u-boot config set project.name x --dry-run --diff --json # Set Dry-Run + strukturierte Hunks
u-boot config --quiet --json                             # bare + --quiet → identisch zu --json (doctor-Pattern)
u-boot config get project.name --quiet --json            # get + --quiet → identisch zu --json
u-boot config set project.name x --quiet --json          # set + --quiet → identisch zu --json
```

(R1-MED-2-Adressierung: drei `--quiet --json`-Pins ergänzt.)

`make gates` grün (lint + test + coverage-gate ≥ 90 % +
docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ **`--json`-Allowlist-Migration**: `"u-boot config": true`,
  `"u-boot config get": true`, `"u-boot config set": true` in
  `jsonAllowlist()`; Reject-Liste schrumpft **von 4 auf 1**
  (heute `jsonallowlist_test.go:27-32`: 4 Reject-Forms `config
  (bare)`, `config get`, `config set`, `template (bare)` —
  nach Slice bleibt nur `template (bare)`; `template list
  --json` ist bereits M3-migriert, R2-LOW-4-Wortlaut-
  Schärfung).
- ✅ **Envelope-Shape (Minimalkontrakt)** für **bare + get**:
  `command="config"`, **`subcommand` Pflicht** (T0-(b) Sub-
  Decision: Wert für bare festzurren — Kandidaten `"show"`,
  `"list"`, `""`). `data`-Carrier:
  - bare: `configShowData{body string}` (heutiges
    `ConfigShowResponse.Body []byte` → `string` für JSON-Safety).
  - get: `configGetData{path string, value string}` (zwei
    String-Felder ohne `omitempty`).
- ✅ **Envelope-Shape (Voll-Schema)** für **set** mit `--dry-
  run`/`--diff`: `cliJSONEnvelope` mit `Subcommand="set"`,
  `DryRun` flag, `PlannedFiles[]` (immer **eine** Zeile —
  `u-boot.yaml` — falls Set kein NoOp; `[]` falls NoOp oder
  Validation-Failure), `Changes[].count` (1 File modify oder
  0 NoOp), `Hunks[]` für `--diff`-Mode mit dem patched-vs-
  current u-boot.yaml-Diff. Konsumenten erkennen Set-NoOp am
  leeren `plannedFiles[]` + `data.noOp == true` + leerem
  `diagnostics: []` (T0-(d) Sub-Decision; R2-HIGH-1: kein
  `level: "info"`-Diagnostic — Spec-§2.1-Bruch).
- ✅ **`config set` Two-Shape Summary verschwindet**: heutige
  `printConfigSetSummary` (`OldValue → NewValue` /`already X`)
  wird im JSON-Mode nicht emittiert — `data`-Carrier trägt
  die Info strukturiert (`configSetData{path, oldValue,
  newValue, noOp bool, appendedSources []string omitempty}`).
- ✅ **`config set --allow-external-feature-sources` im JSON-
  Mode**: heute via `cli/config.go:107-108` als StringSlice-
  Flag. Pre-UC-Validation-Pfad (Path-Kind-Mismatch, Z. 182-
  187) ergänzt um `reportError` analog up-down/logs-Stub.
  **Daten-Mapping**: Hybrid-Form per T0-(c) (R1-MED-1) —
  `oldValue`/`newValue` als CSV-Strings (Status quo);
  zusätzliches `appendedSources []string omitempty` (nur bei
  `path.Kind == ConfigDevcontainerFeatureSourcesAllow` gesetzt)
  damit Konsument weiß, was der Flag beigetragen hat.
- ✅ **`config show` Body als JSON-String** (R1-Lücke):
  heutiges `ConfigShowResponse.Body []byte` byte-identisch
  zur Disk-Datei. Im JSON-Mode wird `data.body` ein
  Go-`string` → trägt UTF-8-Escape-Sequenzen wenn YAML
  CR/Tab/Non-Printables enthält. Kein semantischer Bruch
  (`json.Unmarshal` resynthetisiert die exakten Bytes),
  aber Doku-Pin in §6.9 nötig.
- ✅ **`config set` Custom-Args-Validator** (T0-(l), R1-HIGH-4):
  `validateConfigSetArgs` ersetzt `cobra.ExactArgs(2)` und
  emittiert Envelope bei NoPositionalArg + TooManyArgs
  (analog `validateRemoveArgs`); `validateConfigGetArgs`
  analog für `ExactArgs(1)`. Beide via `isUsageError` →
  Exit 2.
- ✅ **`--quiet --json` semantisch identisch zu `--json`**
  (Cluster-T0-(a) doctor-Pattern). Bare-Show emittiert dann
  `data.body` im JSON statt raw-write.
- ✅ **`--dry-run` für bare/get rejected**: Cluster-Plan Z. 91-100
  sagt explizit "nur modifying tragen Dry-Run". Bare/get sind
  Read-only → `--dry-run --json` muss mit Exit 2 (ggf. neuer
  `ErrDryRunNotApplicable`-Sentinel — Sub-Decision T0-(g))
  rejected werden. Analog logs-`--follow --json`-Reject (T0-(a)).
- ✅ **`--diff` für bare/get rejected**: dito.
- ✅ **`subcommand`-Pflichtfeld-Validierung** (`LH-FA-CLI-007`
  §322): jeder **RunE-emittierte** Envelope mit
  `command="config"` MUSS `subcommand` setzen; T6-Pin gegen
  Empty-Subcommand-Drift. **Cobra-Help-Edge-Case ausgenommen**
  (R1-MED-6): `u-boot config --help --json` läuft durch
  Help-Escape-Hatch in `applyJSONRejectGate`
  (`jsonallowlist.go:112`) und emittiert KEINEN Envelope
  (Cobra rendert Help auf stdout).
- ✅ **Mapper-Tabelle** (`mapConfigErrorToDiagnostic`) analog
  up/down/logs-Pattern mit Switch-Order FS-first → existing
  `ErrConfigFileSystem`, dann Get-/Set-spezifische Sentinels
  (siehe T0-(f)).
- ✅ **Path-Leak-Defense**: `runConfigShow`/`runConfigGet`/
  `runConfigSet` wrappen UC-Errors mit `sanitizeBaseDir(err,
  cwd)` vor `reportError` analog up-down-T5.
- ✅ **RecordingFileSystem-Wiederverwendung** für `config set
  --dry-run`: Pattern-Erbe add T5 1:1. `fsFactory(mode)`
  liefert `RecordingFileSystem` für Dry-Run, Passthrough
  sonst.
- ✅ **CLI-Pin-Tests**: ~18-22 Acceptance-Tests in
  `config_acceptance_test.go` (bare-Envelope + Get-Envelope +
  Set-Envelope-Voll-Schema + Set-Dry-Run + Set-Diff + Set-
  NoOp-Pin + Mapper-Rows + Pre-UC-Validation + Sanitizer +
  Subcommand-Pflicht-Pin + `--dry-run`-Reject auf bare/get).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (config→done),
  neue §6.9-Sektion mit drei Sub-Form-Envelopes + Set-Voll-
  Schema-Beispiel + Subcommand-Pflicht-Doku + Reject-Block
  für `--dry-run`/`--diff` auf Read-only-Forms; §7
  Mutations-Matrix-Zeile aktualisieren (`config set`
  bereits drin als `WriteFile`).
- ✅ **CHANGELOG `### Added`**-Eintrag analog up-down/logs.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a) Slice-Bündelung: drei Sub-Forms in einem Slice?**
  Cluster-Plan Z. 102-113 hat das bereits implizit festgezurrt
  ("alle drei brauchen `--json`"), aber expliziter Festzurrung
  hier:
  (i) **drei Sub-Forms gebündelt** (analog up-down T0-(e) Z.
      369-372): ein gemeinsamer Stub, eine T0-Review-Runde,
      drei Mapper-Zellen, ein T6-Test-File. **Vorteil**:
      gemeinsamer `ConfigUseCase`, gemeinsame Sentinels,
      gemeinsamer Mapper.
  (ii) drei separate Mini-Slices (8a/8b/8c). **Nachteil**: 3x
       Stub-Overhead, 3x Review-Runden, redundante Pattern-
       Erbe-Sektionen.
  Plan-Empfehlung: **(i)** Bündelung. Pattern-Erbe up-down
  bewiesen tragfähig. **R3-festgezurrt: (i) Bündelung**.

- **T0-(b) Bare-Subcommand-Wert**: was steht in
  `envelope.subcommand` für **bare** `u-boot config`?
  (i) `"show"` — semantisch ehrlich (heute `runConfigShow`).
  (ii) `"list"` — Cluster-Plan-Vorschlag Z. 111.
  (iii) `""` (Leer-String) — Cluster-Plan-Vorschlag Z. 112-113
        falls Spec leeren Subcommand erlaubt.
  Plan-Empfehlung: **(i) `"show"`** weil Code-Heim `runConfigShow`
  ist und der Wert `subcommand`-feld in Konsumenten-Filtern
  unmittelbar mit der Code-Realität abgleichbar bleibt. `"list"`
  ist Drift gegen die Code-Benennung; `""` ist Spec-§322-
  fragwürdig.
  **R3-festgezurrt: (i) `"show"`**. Cluster-Plan Z. 111-113
  schlug `"list"`/`""` als Kandidaten vor; R3-Cross-Reference
  schlägt Cluster-Plan-Vorschlag mit Code-Heim-Begründung
  (`runConfigShow` ist der real existierende RunE). Cluster-
  Plan-Klausel war Vorschlag, kein Beschluss (R3-LOW-3-
  Adressierung).

- **T0-(c) DTO-Carrier-Layout**: drei Sub-Forms tragen drei
  Carrier-Types:
  - bare: `configShowData{body string}` — `body` ohne
    `omitempty` (auch leeres `u-boot.yaml` ist legitim, `""`
    statt `null` per Empty-Pin).
  - get: `configGetData{path string, value string}` — beide
    ohne `omitempty`; `value` ist Bare-Scalar-String (`true`/
    `false` für Bool, raw für String).
  - set: `configSetData{path string, oldValue string,
    newValue string, noOp bool}` — `noOp` ohne `omitempty`
    (legitimer Success-False), `oldValue`/`newValue` ohne
    `omitempty` (Empty-String `""` = legitimer initial-unset).

  **`--allow-external-feature-sources`-Repräsentation**
  (R1-MED-1-Adressierung): bei `set
  devcontainer.featureSources.allow <url>
  --allow-external-feature-sources <extra>` baut
  `application/config.go:769,778` `oldValue` und `newValue`
  als `strings.Join(..., ",")` (CSV-String). Sub-Decision:
  (i) Status quo CSV-String in `oldValue`/`newValue` belassen
      — Konsument splittet selbst. Einfach, aber konsumenten-
      unfreundlich.
  (ii) Variante-Type `configSetData` um `oldValues []string`
       + `newValues []string` + `appendedSources []string`
       ergänzen wenn `path.Kind ==
       ConfigDevcontainerFeatureSourcesAllow`. Mehr Code,
       konsumenten-freundlicher.
  (iii) **Hybrid**: `oldValue`/`newValue` bleiben CSV-Strings;
        zusätzlich `appendedSources []string omitempty`
        (NUR bei Allow-Path gesetzt) damit Konsument weiß,
        was der Flag beigetragen hat.
  Plan-Empfehlung: **(iii) Hybrid** — Pattern-Erbe up-down
  `removedVolumes bool` ohne omitempty + zusätzliches Field
  nur bei spezifischen Sub-Decisions. T0-Review präzisiert.

  **`appendedSources`-Quelle** (R2-MED-1-Adressierung): das
  Field wird im **CLI-Layer** (nicht im UC) befüllt aus
  `flags.AllowExternalFeatureSources` **raw** — also der
  User-Input vor Dedupe. KEINE Port-Erweiterung von
  `ConfigSetResponse` nötig (kein
  `ConfigSetResponse.AppendedSources`-Field). Begründung:
  der Konsument soll sehen "was wollte der User anhängen",
  NICHT "was hat sich nach Dedupe wirklich geändert" —
  letzteres ist über `oldValue`/`newValue`-Diff ableitbar.
  Bei Dedupe-Pfad (Allow-Liste enthielt URL schon) zeigt
  `appendedSources: ["url"]` plus `oldValue == newValue` →
  Konsument weiß "Flag war No-op gegen Bestand".

- **T0-(d) `config set` NoOp-Envelope-Form**: heute returnt
  Set bei `OldValue == NewValue` ohne `WriteFile`-Call. Wie
  reagiert der Voll-Schema-Envelope?
  (i) `plannedFiles: []`, `changes: []`, `data.noOp: true`,
      `diagnostics: []`. **Konsument-Disambiguation**: leeres
      `plannedFiles` plus `data.noOp == true` = NoOp.
      **KEIN `level: "info"`-Diagnostic** — Spec-§2.1 (`cli-
      json-output.md` Z. 97-103) verbietet `level: "info"`
      verbatim ("Das Lastenheft beschränkt `diagnostics[].
      level` strikt auf `warn` oder `error`"). Pattern-Erbe
      doctor: All-OK-Fall serialisiert ein leeres
      `diagnostics: []` (R2-HIGH-1-Adressierung — ursprünglich
      vorgeschlagene Info-Diagnostic war Spec-Bruch).
  (ii) `plannedFiles: [{path: "u-boot.yaml", action:
       "modify", changes: [{count: 0}]}]` mit Zero-Count.
       **Drift gegen add-Pattern**: dort wird Zero-Count nie
       emittiert.
  Plan-Empfehlung: **(i)** Empty-PlannedFiles + Info-
  Diagnostic. Pattern-Erbe add: `plannedFiles` listet nur
  echte Mutations.

- **T0-(e) FS-Sentinel-Wiederverwendung**: anders als up-
  down/logs hat `config` **schon einen FS-Sentinel**
  (`ErrConfigFileSystem`, Z. 141 port/driving/config.go).
  Zwei Sub-Decisions:
  - (i) `ErrConfigFileSystem` bereits ausreichend → T2 ist
        rein Co-Migration-Tranche (Wrap-Sites auf Multi-`%w`
        falls heute single).
  - (ii) Separater Read- vs. Write-Sentinel-Split (`ErrConfig
        FileSystemRead` / `ErrConfigFileSystemWrite`) für
        feinere Diagnose-Klassen. **Über-Granularität-Risiko**
        — Pattern-Erbe up-down/logs nutzt einen Sentinel pro
        Subcommand, nicht pro Direction.
  Plan-Empfehlung: **(i)** bestehender Sentinel ausreichend.
  T2 ist Co-Migration + Co-Drift-Check auf alle
  `s.fs.WriteFile`/`s.fs.ReadFile`-Sites in `config.go`.

- **T0-(f) Mapper-Tabelle** (`mapConfigErrorToDiagnostic`)
  Switch-Order — gilt **für alle drei Sub-Forms** (eine
  Mapper-Function, weil Sentinels geteilt sind):

  | # | Sentinel | LH-Code | Exit | Begründung |
  | - | -------- | ------- | ---- | ---------- |
  | 1 | `driving.ErrConfigFileSystem` | `LH-NFA-REL-003` | 14 | FS-first damit Multi-`%w` mit FS+Validation auf FS-Klasse fällt |
  | 2 | `driving.ErrConfigSchemaInvalid` | `LH-FA-CONF-002` | 10 | Schema-Bruch vor Path-Unknown (Schema-Defekt > Form-Defekt) |
  | 3 | `driving.ErrConfigPostPatchSanityFailed` (NEU, T2 — R2-MED-4) | `LH-FA-CONF-002` | 10 | Post-Patch-Roundtrip-Sanity-Fehler (Schema-Drift-Indikator); semantisch nahe `ErrConfigSchemaInvalid` aber separat extrahiert aus heutigem `ErrConfigValueInvalid`-Multi-Use |
  | 4 | `driving.ErrConfigPathUnknown` | `LH-FA-CONF-005` | 10 | Path-Whitelist-Bruch |
  | 5 | `driving.ErrConfigWriteRejected` (NEU, T2 — R2-MED-4) | `LH-FA-CONF-005` | 10 | WriteAllowed-Reject (`services.<svc>.enabled` etc.); separat extrahiert aus heutigem `ErrConfigValueInvalid`-Multi-Use; Hint "u-boot add <svc>" im message |
  | 6 | `driving.ErrConfigValueInvalid` | `LH-FA-CONF-001` | 10 | **Nur Value-Coercion-Bruch** (set; nach T2-Separation der zwei anderen Klassen) |
  | 7 | `driving.ErrConfigValueNotSet` | `LH-FA-CONF-005` | 10 | Optionaler Pfad nicht gesetzt (nur get) |
  | 8 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | Pattern-Erbe up/down/generate/logs (Environment-Operation) |
  | 9 | `cli.ErrDryRunNotApplicable` (NEU, T2) | `LH-FA-CLI-006` | 2 | bare/get rejecten `--dry-run` (T0-(g) Option (i.a)) |
  | 10 | Default (unknown) | `LH-FA-CLI-006` | 1 | Fallback |

  **Cross-Slice-Klassen-Pin**: `ErrProjectNotInitialized`
  mappt hier auf **`LH-FA-INIT-001`** (Environment-Operation
  Pattern-Erbe up/down/generate/logs) — NICHT auf
  `LH-FA-ADD-001` wie bei add/remove. Bewusste Cluster-
  Konvention.

- **T0-(g) `--dry-run`/`--diff` auf bare/get Reject-Sentinel**:
  bare/get sind Read-only → tragen kein `--dry-run`/`--diff`.
  Drei Reject-Optionen:
  (i.a) **Neuer `cli.ErrDryRunNotApplicable`-Sentinel** + Flag
       wird an bare/get Cobra-Cmd **registriert** und im RunE
       rejected (Envelope-konform). Pattern-Erbe logs T0-(a)
       Option (A): `--follow` ist registriert, `--follow --json`
       wird in `runLogs`-Stage-1 rejected.
  (i.b) **Neuer Sentinel** ABER Flag **nicht** an bare/get
       Cobra-Cmd registriert → Cobra emittiert `unknown flag
       --dry-run` mit Roh-stderr-Output (kein Envelope!) →
       Spec-§1841-Bruch.
  (ii) Re-use `cli.ErrJSONNotImplemented` (heute für
       Allowlist-Reject). **Semantischer Drift**: dieser
       Sentinel sagt "Form noch nicht migriert", nicht
       "Form unterstützt das Flag nicht".
  (iii) Cobra-Native: `MarkFlagsMutuallyExclusive` auf
        Cobra-Ebene. **Nachteil**: keine JSON-Envelope-
        Emission im Reject-Pfad (Cobra schreibt direkt
        nach stderr). Spec-§1841-Bruch.
  Plan-Empfehlung: **(i.a)** — Flag registrieren UND im RunE
  rejecten. Pattern-Erbe logs `ErrFollowJSONNotSupported`-
  Pfad ist (i.a)-konform (`--follow` ist registriert,
  Reject im RunE). (R1-MED-3-Adressierung.)

  **R3-festgezurrt: (i.a)**. Begründung: Variante (i.c)
  "PersistentFlag nur am set + Parent-Chain-Lookup via
  `cmd.Flags().Changed("dry-run")`" hängt von Cobra-Parent-
  Chain-Default-Behavior ab, das nicht stabil dokumentiert
  ist — Drift-Risiko gegen Cobra-Upgrades (heute v1.10.2 →
  künftig v2). Die Help-UX-Pollution durch synthetische
  Flags ist akzeptabel: der `--help`-Output zeigt zwei Flags
  mit klarer Reject-Beschreibung ("only valid for `config
  set`") — Pattern-Erbe init's `--no-confirm` ist analog
  semantisch synthetisch für bestimmte Pfade.

- **T0-(h) `subcommand`-Pflicht-Form für `config get` / `config
  set`**: bei Cobra-Compound (`u-boot config get`) trägt
  envelope.subcommand `"get"` bzw. `"set"`. **Quelle**:
  `cmd.Name()` im RunE liefert das Cobra-Sub-Verb (`"get"`,
  `"set"`); für bare `u-boot config` liefert `cmd.Name()`
  `"config"` und der CLI-Layer setzt `subcommand` manuell auf
  den T0-(b)-festgezurrten Wert (`"show"`). Kein Args-Inspect,
  kein `cmd.CommandPath()`-Parse (R1-LOW-3-Adressierung).
  Test-Pin gegen Empty-Subcommand-Drift. **Pattern-Erbe**:
  template-Slice hat dasselbe Problem (`template list`). Sub-
  Decision-Form: geteilter Helper `cobraPathToSubcommand(cmd)`
  in `cli/`-Sub-Package extrahieren ODER inline-Switch im
  jeweiligen RunE.
  Plan-Empfehlung: **inline-Switch** (zwei Stellen reicht
  noch nicht für Helper-Extraktion); falls config + template
  zusammen tragen, Helper in Cluster-T_close-Tranche.

  **Cobra-Help-Edge-Case** (R1-MED-6-Adressierung):
  `u-boot config --help --json` ist KEIN Envelope-Pfad —
  `applyJSONRejectGate` (`jsonallowlist.go:112`) Help-Escape-
  Hatch returnt vor RunE; Cobra rendert Help auf stdout. Die
  `subcommand`-Pflicht gilt ausschließlich für RunE-Pfade.
  AK + T6-Pin entsprechend formulieren.

- **T0-(i) `config set` Pre-UC-Validation-Pfade**: heute
  `runConfigSet:174-187` validiert (a) Path-Parse via
  `domain.NewConfigPath`, (b) AllowExternalFeatureSources-
  Path-Kind-Mismatch. Beide sind Pre-UC-Errors und brauchen
  `reportError`-Wrap analog up-down/logs T5.

- **T0-(j) Echte Spec-Anker für `LH-FA-CONF-*`-Codes statt
  `LH-FA-CLI-006`-Fallback**: heutiger Doc-Block in
  `cli/config.go:14-65` nutzt nur generische Exit-Code-Tabelle
  (Z. 53 "Exit codes (LH-FA-CLI-006)") und nennt **keinen
  per-Sentinel-LH-FA-CONF-Mapping**; lediglich Spec-Bezug für
  Subcommand-Tree (Z. 15) und writable paths (Z. 37) als
  `LH-FA-CONF-001 / §D1`. Heutige Mapping-Form im CLI-Layer
  fällt deshalb auf `LH-FA-CLI-006` für **alle** Validation-
  Errors (Pfad-Unknown, Wert-Invalid, Schema-Invalid). Der neue
  `mapConfigErrorToDiagnostic` muss die echten Spec-Anker
  zuweisen — Pattern-Erbe up-down/logs nutzt echte Spec-Anker
  (`LH-FA-UP-001`, `LH-FA-INIT-001`, `LH-FA-INIT-006`) statt
  generischen Fallback. (R1-HIGH-3-Adressierung: ursprüngliche
  Begründung „heutige Code-Doc-Block-Mapping honorieren" war
  Hineinprojektion — der Doc-Block enthält das nicht.)
  Plan-Empfehlung: **echte Spec-Anker** `LH-FA-CONF-001/002/005`
  per Mapper-Row T0-(f), siehe dort.

- **T0-(k) `config set --diff`-Renderer-Pfad**: Pattern-Erbe
  add/init/generate-Slices nutzen Pure-Go-Diff (Cluster-T0-(d)
  Option (Beides) Z. 582-621). `config set` patcht eine
  einzige u-boot.yaml-Datei mit einem Scalar-Wert — der Diff
  ist eine Single-Hunk-Modifikation auf einer Datei. Sub-
  Decision: ist der existierende Diff-Renderer-Helper direkt
  wiederverwendbar oder braucht config einen eigenen wegen
  YAML-Schema-Quirks?
  Plan-Empfehlung: **direkte Wiederverwendung**. Patched-
  Bytes vs. Current-Bytes durch denselben Pure-Go-Diff-
  Renderer wie add/init/generate.

- **T0-(l) Custom-Args-Validator `validateConfigSetArgs`**
  (R1-HIGH-4-Adressierung): `config set` trägt heute
  `cobra.ExactArgs(2)` (`cli/config.go:102`). Bei `u-boot
  config set <path>` (nur 1 Positional) emittiert Cobra **vor**
  PersistentPreRunE einen `accepts 2 arg(s), received 1`-Error
  — **ohne JSON-Envelope**. Spec-§1841/§1842-Bruch.
  Pattern-Erbe remove (`slice-v1-cli-json-dry-run-remove` T7
  Custom-Args-Validator): `validateRemoveArgs(a *App)` ersetzt
  `cobra.ExactArgs(1)` und emittiert den Envelope **selbst**
  bei NoPositionalArg + TooManyArgs (Exit 2 via
  `isUsageError`).
  Sub-Decision:
  (i) **`validateConfigSetArgs(a *App)`-Closure** analog
      remove: prüft `len(args) == 2`, emittiert bei Mismatch
      JSON-Envelope auf stdout VOR Cobra-Return mit Exit 2.
  (ii) `cobra.ExactArgs(2)` belassen, Cobra-Roh-Output bei
       Args-Mismatch akzeptieren. **Pattern-Erbe-Bruch**
       gegen remove + Cluster-Plan §1841.
  Plan-Empfehlung: **(i)** Custom-Args-Validator. Analog für
  `config get`-`ExactArgs(1)` und für bare `config`:
  **alle drei Sub-Forms** bekommen `validateConfigSetArgs` /
  `validateConfigGetArgs` / `validateConfigShowArgs` aus
  Konsistenz-Disziplin (R2-HIGH-3-Adressierung).

  **`u-boot config foo`-Pfad** (unknown command, R2-HIGH-3):
  bei bare-Cobra-Cmd mit `cobra.NoArgs` plus AddCommand-
  Children (`get`/`set`) emittiert Cobra `Error: unknown
  command "foo" for "u-boot config"` mit Exit 2 — **ohne
  JSON-Envelope**. Symmetrie-Bruch gegen `remove`-Pattern
  (`4fb3fea` R13-MED-1). `validateConfigShowArgs` muss bei
  Positional > 0 (das ist ein nicht-Subcommand-Token) den
  Envelope-konformen Reject emittieren. Pattern: prüft erst
  ob `args[0]` ein registriertes Sub-Verb ist (`get`/`set`);
  falls nein → Envelope-Emit mit
  `unknown sub-command for config`-Diagnostic. Cobra ruft
  `validateConfigShowArgs` nur dann auf, wenn kein Child-
  Match — also ist der Validator-Pfad sauber.

  **Cobra-Mechanik-Anker** (R3-MED-1-Adressierung):
  `spf13/cobra/command.go:Find()` führt `legacyArgs`-Pfad
  aus, der bei `u-boot config foo` (kein Child-Match) auf
  den Parent-Cmd `config` als Target dispatcht plus
  `args=["foo"]`. Cobra v1.10.2 ruft danach `Args(cmd,
  args)` — also genau `validateConfigShowArgs(cmd, ["foo"])`.
  Pattern-Erbe `remove` ist **kein direkter Vorbild**
  (`remove` ist Leaf-Cmd ohne Children); Vorbild ist
  konzeptionell die Custom-Validator-Form selbst. T6-Pin
  `TestConfigUnknownSubcommand_FooEmitsEnvelope` pinnt
  den Pfad und schützt gegen Cobra-Mechanik-Drift bei
  Version-Upgrades.

- **T0-(m) `ErrConfigValueInvalid`-Drei-Klassen-Disambiguation**
  (R2-MED-4-Adressierung): heute mappt `ErrConfigValueInvalid`
  drei semantisch unterschiedliche Pfade auf einen Sentinel:
  (a) **Value-Coercion-Bruch** (`config.go:278,287,296,308,313`):
      `ParseBool`/`NewProjectName`/CSV-Parse-Fehler.
  (b) **WriteAllowed-Reject** (`config.go:251-256`): `services.
      <svc>.enabled` ist nicht writable, Hint "u-boot add <svc>".
  (c) **Post-Patch-Schema-Sanity** (`config.go:376,388`):
      Roundtrip nach `PatchScalar` produziert nicht den
      erwarteten Wert (sehr selten — Schema-Drift-Indikator).
  Konsumenten können die drei Klassen heute nicht auseinander-
  halten (gleicher `code: LH-FA-CONF-001`, gleicher Exit 10).
  Sub-Decision-Optionen:
  (i) **Separate Sentinels** `ErrConfigWriteRejected` (für b)
      und `ErrConfigPostPatchSanityFailed` (für c) im Port-
      Layer. Mapper-Tabelle T0-(f) bekommt drei neue Rows.
      Konsument disambiguiert per `code`. Pattern-Erbe doctor
      (per-Check-Sentinels).
  (ii) **Multi-Use mit Code-Suffix**: `LH-FA-CONF-001-Coerce`
       vs. `LH-FA-CONF-001-WriteReject` vs. `LH-FA-CONF-001-
       PostPatchSanity`. KEIN Spec-Anker — eigener Cluster-
       Konvention-Bruch.
  (iii) **Status-quo**: alle drei auf `LH-FA-CONF-001`
        belassen, Konsument disambiguiert per
        `diagnostic.message`-Substring-Match. Brüchig.
  Plan-Empfehlung: **(i)** separate Sentinels — Pattern-Erbe
  init/add (`ErrInitFileSystem` separat vom `ErrInvalid…`).
  Erfordert Port-Erweiterung in T2 (drei neue Sentinels
  exportiert in `port/driving/config.go`), Mapper-Rows in
  T0-(f) entsprechend erweitert (5 → 7 Rows fachlich plus
  FS/Schema/Default = 9-10 Rows total).

  **`LH-FA-CONF-005`-Multi-Use mit drei Sentinels**
  (R3-MED-2-Adressierung): Mapper-Tabelle T0-(f) belegt
  `LH-FA-CONF-005` jetzt in **drei Rows** (Row 4
  `ErrConfigPathUnknown`, Row 5 `ErrConfigWriteRejected`,
  Row 7 `ErrConfigValueNotSet`). Pattern-Erbe `LH-FA-ADD-007`
  (remove) war doppelt belegt — drei Sentinels auf einem
  LH-Code ist eine Steigerung. Konsumenten können per
  `code: LH-FA-CONF-005` allein nicht disambiguieren —
  müssen auf den Sentinel-Klassen-Hint im Message-Body
  zurückgreifen. T8 dokumentiert das in `cli-json-output.md`
  §6.9 als expliziter Konsumenten-Disambiguation-Block
  analog `(code, exitCode)`-Tupel für up-down/logs (nur
  hier: `(code, message-prefix)` weil keine ExitCode-
  Disambiguation hilft, alle drei sind Exit 10).

- **T0-(n) Logger-Output im JSON-Mode + Orphan-Feature-Warn-
  Migration** (R2-MED-5-Adressierung; R3-MED-4-Erweiterung
  um zwei Debug-Sites): heute **fünf** `s.logger.*`-Sites
  im `ConfigService.Set`-Pfad: drei `Info` plus zwei
  `Debug`. Alle werden durch den `SilenceLogger`-Swap
  betroffen (`logger := s.logger; if req.SilenceLogger {
  logger = noopLogger{} }` an Method-Beginn). Sites:
  `config.go:158` (NoOp-Debug), `:200` (Set-Success-Info),
  `:237` (`maybeWarnOrphanFeatureActivation` — User-Warn
  dass ein Feature aktiviert wurde aber sein Service nicht
  registriert ist), `:782` (Allow-NoOp-Debug), `:813`
  (Allow-Set-Success-Info).
  Im JSON-Mode emittieren sie stderr-Logs parallel zum
  Envelope auf stdout. Pattern-Erbe up-down `SilenceProgress`,
  remove `SilenceConfirmer` haben Bool-Field-Pattern.
  Sub-Decision-Optionen:
  (i) **`ConfigSetRequest.SilenceLogger bool`-Field**: CLI-
      Layer setzt `req.SilenceLogger = flags.JSON`; UC-Body
      swappt `s.logger` lokal auf `noopLogger{}`. Pattern-
      Erbe up-down ProgressSink-Branch.
  (ii) **Status-quo belassen**: stderr ist separate Channel,
       JSON-Konsumenten die stderr capturen müssen Mix
       parsen. Pattern-Erbe-Bruch gegen up-down/remove.
  Plan-Empfehlung: **(i)** Bool-Field-Pattern.
  **Orphan-Feature-Warn-Migration** (separat, eigener Sub-
  Decision-Block):
  (a) **WARN-Migration in `diagnostics[]`**: Pattern-Erbe
      remove `mapWarningsToDiagnostics` — der
      `maybeWarnOrphanFeatureActivation`-Output wird als
      `LH-FA-DEV-003` / `level: "warn"`-Entry in `diagnostics[]`
      des Envelopes emittiert. ConfigSetResponse braucht
      `Warnings []driving.WarningEntry`-Field analog
      `RemoveResponse`. Pattern-konsistent + JSON-Konsument
      sieht User-Warn ohne stderr-Capture-Pflicht.
  (b) **Status-quo stderr-Log**: WARN bleibt unsichtbar im
      JSON-Mode. Konsument-Klage-Risiko.
  Plan-Empfehlung: **(a)** WARN-Migration. Erfordert Port-
  Erweiterung in T2 (`ConfigSetResponse.Warnings`-Field).

- **T0-(o) Mid-Stage-Failure-Snapshot bei `config set --dry-
  run --json`** (R2-Sub-Decision-Lücke G-1): bei Stage-N-
  Failure zwischen Stages 1-5 (vor WriteFile) emittiert der
  Voll-Schema-Envelope:
  (i) **Minimal-Envelope mit Error-Diagnostic** + leeres
      `plannedFiles: []` + kein `data`-Carrier. Konsument
      sieht "Validation failed, kein Plan erzeugt".
  (ii) **Voll-Schema-Envelope** mit `data.path`/`data.oldValue`
       befüllt soweit bekannt, `data.newValue` leer falls
       Stage-1 fehlschlägt; `plannedFiles: []` weil
       Recorder leer. Mehr Datum für Konsument-Debug.
  Plan-Empfehlung: **(i)** Minimal-Form für Validation-
  Failures (Stages 1-4); **(ii)** Voll-Form ab Stage 5
  (NoOp/Schema-Sanity nach Patch). Pattern-Erbe init T0-(l)
  Stage-Map-Form analog. T6-Pin pro Stage.
  **R3-festgezurrt: Plan-Empfehlung** (Mixed-Form je Stage).

- **T0-(p) Hint-String-Sanitization für Error-Messages**
  (R2-Sub-Decision-Lücke D-Hint): `ErrConfigPathUnknown`
  und `ErrConfigSchemaInvalid` tragen ggf. Pfad-/YAML-
  Decoder-Output mit Filename. `sanitizeBaseDir(err, cwd)`
  greift nur auf BaseDir-Pfade — andere Pfade (z. B. yaml.v3-
  Multi-Doc-Include) würden durchgereicht. T5-Pflicht:
  `sanitizeBaseDir`-Wrap **alle** UC-Errors umfassen (heute
  Pattern-Erbe), plus Worst-Case-T6-Pin mit synthetisch
  konstruiertem Error der Filename-Leak enthält. Kein
  separater Sanitizer-Helper nötig — `cli/sanitize.go` ist
  ausreichend wenn `runConfig*`-RunE konsequent `sanitize
  BaseDir(err, cwd)` vor `reportError` aufruft.
  **R3-festgezurrt: Plan-Empfehlung** (alle UC-Errors mit
  `sanitizeBaseDir` wrappen via `cli/sanitize.go`).

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(p) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog up-down/logs T1): `cli/sanitize.go`-Helper, `RecordingFileSystem`-Adapter, Pure-Go-Diff-Renderer existieren bereits aus add/init/generate/remove/up-down T5 | — (entfällt) | T0 |
| T2 ✅ (2026-06-08) | **Geliefert**: **`ConfigSetRequest.PreviewMode driving.PreviewMode`-Field** + **`ConfigSetRequest.SilenceLogger bool`-Field** (Pattern-Erbe `UpRequest.SilenceProgress`; R2-MED-5 T0-(n)). **`ConfigSetResponse.Warnings []driving.WarningEntry`-Field** für Orphan-Feature-WARN-Migration (Pattern-Erbe `RemoveResponse.Warnings`; R2-MED-5 T0-(n)). **Zwei neue Port-Sentinels** (R2-MED-4 T0-(m)): `driving.ErrConfigWriteRejected` + `driving.ErrConfigPostPatchSanityFailed`. `cli.ErrDryRunNotApplicable`-Sentinel (T0-(g) Option (i.a)) im CLI-Layer (R2-LOW-3-Fix). `configSetFlags.JSON`/`Quiet` read-through (in Set-Closure populated, Pattern-Erbe logs T2). Pin-Tests (`port/driving/config_test.go` + `cli/config_test.go`). **KEIN neuer FS-Sentinel** (`ErrConfigFileSystem` existiert). **KEIN `ConfigSetResponse.AppendedSources`-Field** — `appendedSources` lebt CLI-Layer-only (R2-MED-1 T0-(c); R3-LOW-2). **Nach T5 verschoben** (Lint-/Behavior-Grund — Präzedenz logs/up-down T2 = `feat(port)`; black-box `cli_test` kann unexported Structs nicht ohne RunE-Signatur-Refactor referenzieren, + user-sichtbare `--dry-run`/`--diff` ohne Reject-Wiring wäre Behavior-Trap): `configGetFlags{JSON, Quiet}` + `configShowFlags{JSON, Quiet}` + `configSetFlags.DryRun`/`.Diff`-Felder + `--dry-run`/`--diff`-Cobra-Flag-Registrierung. | ~130 (Port + Scaffold + Pins) | T0 |
| T3 ✅ (2026-06-08) | **Geliefert** (`make gates` grün, Coverage 91.20 %): **Multi-`%w`-Wrap-Migration** der 5 FS-Read/Write-Wrap-Sites in `application/config.go` (real: 3× `%w: read %q: %w` + 2× `%w: write %q: %w`) von `%v`-tail auf echtes Multi-`%w` (Pattern-Erbe up-down T3). **`ErrConfigValueInvalid`-Multi-Use-Splitting** (R2-MED-4 T0-(m)): WriteAllowed-Reject (`writeRejectedError`, beide Sites) → `ErrConfigWriteRejected`; Post-Patch-Sanity → `ErrConfigPostPatchSanityFailed`; Value-Coercion (`coerceConfigValue`) + Allowlist-Enforcement (`featureSourceInAllow`-Reject, user-actionable) bleiben `ErrConfigValueInvalid`. **Logger-Branch** im `Set`- + `setFeatureSourcesAllow`-Body (`logger := s.logger; if req.SilenceLogger { logger = noopLogger{} }`, T0-(n)). **Orphan-Feature-WARN-Migration**: `maybeWarnOrphanFeatureActivation` (jetzt freie Funktion mit `logger`-Param) liefert zusätzlich `[]driving.WarningEntry{{Code:"LH-FA-DEV-003", Level:"warn", Subject:<feature>}}` → `ConfigSetResponse.Warnings` (Dual-Emission: stderr-Info + strukturierte WARN). T3-Tests (`config_t3_test.go`: SilenceLogger-Suppression, Dual-Emission, Warnings-survive-suppression, Non-Orphan-nil) + zwei WriteRejected-Test-Updates. **Refinement gegenüber Plan-Zeilennummern**: Pre-Scan nannte nur Z. 376/388 für Post-Patch-Sanity; die strukturell identischen `revalidateFeatureEntry`-Sites (absent/unbound/empty) wurden konsistenz-halber mitmigriert (sonst Mapper-Row-6-„nur Value-Coercion"-Bruch). **Nach T4 verschoben**: `PreviewMode`-Handling (`fsFactory`-Feld + `selectFS` + Set-Write-Routing) — zusammen mit dem `WithFactory`-Konstruktor + Composition-Root, weil die `selectFS`-Nicht-nil-Branch sonst untestbar/uncovered wäre (Coverage-Gate). KEIN ProgressSink-Branch. KEIN Confirmer-Branch. | ~80 (ohne PreviewMode → T4) | T2 |
| T4 ✅ (2026-06-08) | **Geliefert** (`make gates` grün, Coverage 91.20 %): **PreviewMode-Cluster (aus T3 verschoben)** + **Composition-Root-Erweiterung** (R2-HIGH-2): `ConfigService.fsFactory`-Feld + nil-safe `selectFS(mode) (FS, RecorderPort)`-Methode (Pattern-Erbe `AddServiceService`); `Set` + `setFeatureSourcesAllow` routen ihren `WriteFile` über `fs, recorder := s.selectFS(req.PreviewMode)` (Reads bleiben auf Production-`s.fs`) und surfacen `recorder.Captured()` via `mapCaptureToPlannedFiles` → **`ConfigSetResponse.PlannedFiles`** (R-T4-1, s.u.). `NewConfigServiceWithFactory(fsFactory, yaml, logger)`-Konstruktor neu (Bootstrap-FS aus `fsFactory(PreviewNone)`). `cmd/uboot/main.go`: `configFSFactory := newPreviewFSFactory(fsAdapter)` (jetzt **fünf** Factories) + Konstruktor-Wechsel `NewConfigService` → `NewConfigServiceWithFactory`. Damit ist `ConfigService` nicht mehr der einzige Modifying-Service ohne `WithFactory`-Variante. Tests `config_factory_test.go`: DryRun-touch-nichts-in-Production + `PlannedFiles[0].NewContent` befüllt, PreviewNone-persistiert + PlannedFiles nil, Legacy-Konstruktor-ignoriert-PreviewMode + PlannedFiles nil. **T4-Review-Followup R-T4-1 (HIGH)**: ursprünglicher T4-Entwurf verwarf den Recorder („CLI baut plannedFiles statisch") — falsch, weil `config set --diff` die patched+current Bytes (`PlannedFile.NewContent`/`OldContent`) für den geteilten `mapPlannedFilesToWire`/`writeDiff`-Renderer braucht; die existieren nur im Recorder. Fix: `ConfigSetResponse.PlannedFiles []driving.PlannedFile`-Feld neu + Recorder-Surfacing wie add. | ~60 + R-T4-1 | T3 |
| T5 ✅ (2026-06-08) | **Geliefert** (`make gates` grün, Coverage 91.20 %): vollständige `config.go`-Neufassung — drei `runConfig*` auf Cluster-Signatur (flags-Struct + JSON-Pfade); drei Data-Carrier (`configShowData{body}`, `configGetData{path,value}`, `configSetData{path,oldValue,newValue,noOp,appendedSources omitempty}`); Subcommand-Pflicht (`show`/`get`/`set`) auf ALLEN Envelopes inkl. Error-Pfad via neue **subcommand-bewusste `reportErrorSub`/`writeErrorEnvelopeSub`** in `erroremission.go` (additiv, bestehende Single-Form-Caller unverändert; template 9/9 erbt sie). Allowlist-Migration 3 Forms (Reject 4→1). `mapConfigErrorToDiagnostic` Switch-Order T0-(f) FS-first (10 Rows). Konsolidierter Custom-Args-Validator `configArgsValidator(a, sub, base)` (T0-(l): NoArgs für bare, ExactArgs für get/set; Envelope-Reject vor Cobra-Return). `config set` Voll-Schema mit `previewModeFromFlags`→`PreviewMode` + plannedFiles/changes/hunks aus `resp.PlannedFiles` (R-T4-1) via `mapPlannedFilesToWire`/`writeDiff`. bare/get `--dry-run`/`--diff`-Reject via `ErrDryRunNotApplicable` (synthetische Flags an die Struct-Felder gebunden) + `isUsageError`-Branch (Exit 2). `SilenceLogger=flags.JSON`-Wiring im Set-Request. WARN→`diagnostics[]` via `mapWarningsToDiagnostics`. `sanitizeBaseDir` auf allen UC-Errors (T0-(p)). Acceptance-Tests `config_acceptance_test.go` (Show/Get/Set-Data, NoOp, Dry-Run-Voll-Schema, Diff-Hunks, Reject bare/get, Args-Mismatch, WriteRejected→Exit10, WARN→diagnostics, Human-Mode). **Cobra-unknown-subcommand-Mechanik-Pin** (R2/R3) → T6. | ~350-450 | T2, T4 |
| T6 ✅ (2026-06-08) | **Geliefert** (`make gates` grün, Coverage 91.30 %): T5 lieferte schon die Basis-Suite (`config_acceptance_test.go`: Show/Get/Set-Data, NoOp, Dry-Run-Voll-Schema, Diff-Hunks, bare/get-Reject, Args-Mismatch, WriteRejected→Exit10, WARN→diagnostics, Human-Mode). T6 ergänzte: **`--quiet --json`-Pin für alle drei Forms** (R1-MED-2), **white-box Mapper-Rows 1-10** + **FS-first-Switch-Order-`_ByDesign`-Pin** + Pure-FS-Exit-14-Pin (`config_internal_test.go`, R2-MED-4/T0-(f)), **`TestConfigJSON_UnknownSubcommandEmitsEnvelope`** (R3-MED-1: Cobra-`config foo`-Dispatch → Parent-`configArgsValidator` → Envelope; bestätigt die Cobra-Mechanik-Annahme), **`LH-FA-CONF-005`-Disambiguation-Pin** (R3-MED-2: drei Rows gleicher Code, paarweise verschiedene Messages), **Cobra-Help-Edge-Case-Pin** (`--help --json` kein Envelope, R1-MED-6), **Sanitizer-Worst-Case-Pin** (abs-Pfad relativiert, T0-(p)), **Subcommand-Pflicht-Pin** (alle Forms + Error-Pfad, §322), **Mid-Stage-Shape-Pin** (plain=minimal vs --dry-run=voll, T0-(o)). Orphan-WARN/SilenceLogger sind über die Application-Tests (T3 `config_t3_test.go`) + den CLI-WARN→diagnostics-Pin abgedeckt. | T5+T6 ~750-950 | T5 |
| T7 ✅ (2026-06-08) | **Geliefert** (`make gates` grün, Coverage 91.30 %): drei Review-Runden insgesamt — (1) Selbst-Review T2–T4 → **R-T4-1 (HIGH)** PlannedFiles-Surfacing; (2) **unabhängiger** Reviewer-Agent Port/Application → **R-IR-1 (HIGH)** ExitCode-Regression; (3) **unabhängiger** Reviewer-Agent CLI-Layer (T5/T6) → **R-CLI-1 (MED)**: `configArgsValidator` leakte ein Voll-Schema-Error-Envelope auf den Read-only-Forms (`config foo --dry-run --json` / `config get --dry-run --json` ohne Arg), inkonsistent zum RunE-Reject-Pfad. Fix: Read-only-Validatoren hartkodieren `dryRun/diff=false` (nur `set` reicht echte Flags durch) + zwei Regression-Pins (Read-only-Args-Error bleibt Minimal trotz `--dry-run`). LOW-1: Diff-Hunks-Test auf echte Hunk-Inhalt-Assertion (old→new) verschärft. LOW-2 (`Quiet`-Feld zugewiesen-aber-ungelesen) bewusst belassen (no-op-Semantik, pattern-konsistent, im Struct-Doc dokumentiert). Übriger CLI-Layer clear (Subcommand-Pflicht alle Pfade, Mapper↔ExitCode-Kopplung, erroremission-Refactor verhaltensidentisch, NoOp-Envelope, Sanitizer, ctx-Propagation). | ~150-250 | T6 |
| T8 | Closure: CHANGELOG, `cli-json-output.md` **neue §6.9-Sektion** mit drei Sub-Form-Envelopes + Set-Voll-Schema-Beispiel + Subcommand-Pflicht-Doku (RunE-only-Geltungsbereich) + `--dry-run`/`--diff`-Reject-Doku für Read-only-Forms; **`(code, exitCode)`-Tupel-Disambiguation-Block entfällt für config** (R2-MED-6); **NEU: `LH-FA-CONF-005`-Multi-Use-Disambiguation-Block** (R3-MED-2) — `code: LH-FA-CONF-005` mit drei Sentinels (Path-Unknown / Write-Rejected / Value-Not-Set), Konsumenten disambiguieren per Message-Prefix oder Sentinel-Klassen-Hint da alle drei Exit 10 sind. **§6.1 Reject-Liste-Update von 4 auf 1**. **§6 Tabelle Z. 374 Status `offen → done`**. **§7 NEUEINTRAG** `config set` mit `WriteFile`-Spalte ✓. roadmap done-Zähler 7→8. carveouts.md-Einträge **für vier Folge-Slice-Stubs** (R3-MED-3): `slice-v1-config-multi-path-set`, `slice-v1-config-list-subcommand`, `slice-v1-config-multi-path-get`, `slice-v1-config-structured-hint`. Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: **~1500-1900** (T2 ~130 + T3 ~80 + T4
~40 + T5 ~350-450 + T6 ~750-950 + T7 ~150-250 = ~1500-1900
Tranchen-Summe). Deutlich größer als logs ~700-800 weil
drei Sub-Forms gleichzeitig + Set-Modifying-Surface mit
RecordingFileSystem + Diff + Dry-Run + Custom-Args-Validator
+ T4-Composition-Root-Wechsel (R2-HIGH-2) + zwei neue Port-
Sentinels (R2-MED-4) + Logger-Silencing-Pattern + WARN-
Migration (R2-MED-5). Pattern-Erbe von add/init/generate/
remove (Modifying-Pattern: RecordingFS + Diff + fsFactory +
Dry-Run + Custom-Args-Validator + `WithFactory`-Konstruktor)
plus up-down/logs (Read-only-Pattern: FS-Sentinel-Switch-
Order + Sanitizer + Reject-Sentinel für inkompatible Flag-
Kombi). (R1-MED-5 + R2-HIGH-2 + R2-MED-4/5-Adressierung.)

## Out of Scope

- **Multi-Path-Set** (`u-boot config set a.b 1 c.d 2` —
  zwei Pfade in einem Call): heute `ExactArgs(2)`, ein Pfad
  pro Set-Aufruf. Multi-Path würde Transaktions-Semantik
  brauchen (alle oder keine schreiben) und ist Spec-Erweiterung
  außerhalb V1. Folge-Slice falls Real-World-Druck.
- **`config list` als eigener Subcommand** (Listing aller
  Pfade mit Werten): heute nicht in Spec. Bare `u-boot config`
  liefert byte-identisch das gesamte `u-boot.yaml`; ein
  strukturierter Pfad-Wert-Tree wäre eigener Subcommand.
  Folge-Slice falls Real-World-Druck.
- **`config get --json-array`** (mehrere Pfade in einem Call):
  Plan-Vorschlag T0-(b) bleibt bei single-path. Multi-Path-Get
  analog Multi-Path-Set Out-of-Scope.
- **WriteAllowed-Reverse-Mapping als Hint-Field**: heute
  Reject-Message embed `"u-boot add <svc>"` als String. Ein
  strukturiertes `data.hint{action: "add", argument: "<svc>"}`
  wäre konsument-freundlicher. Out-of-Scope V1; Folge-Slice
  falls Real-World-Druck.
- **`subcommand`-Pflicht für alle Forms aufweichen falls Spec
  ändert**: T0-(b)-Entscheidung gilt für aktuellen Spec-§322-
  Stand. Cluster-T_close-Audit darf neu festzurren falls Spec-
  Update.
- **Context-Cancellation mid-`config set`** (R1-MED-4-
  Adressierung, Pattern-Erbe init T0-(p) + remove R11-MED-F2):
  Ctrl-C zwischen Stage 1-5 (Validation) und Stage 6 (WriteFile,
  `application/config.go:194`) bleibt Status-quo Default-Branch
  Exit 1 — kein partial-write-Risk weil `WriteFile` atomar ist
  (oder gar nicht ausgeführt wird). Cross-Cutting-Folge-Slice
  ist Pattern-Erbe init/remove-Carveout, nicht config-spezifisch.
- **`fsFactory`-NPE-Schutz** (R1-MED-4-Adressierung, Pattern-
  Erbe remove R11-MED-F2): Composition-Root-Bug, der `nil`-FS
  aus `selectFS(mode)` liefert, ist Defekt-Klasse — kein User-
  Pfad. Status-quo wie add/init/generate/remove (Composition-
  Root-Tests fangen das via panic in Acceptance-Setup).
- **YAML-Comments-Preservation bei `config set
  devcontainer.featureSources.allow`** (R1-MED-4-Adressierung):
  `setFeatureSourcesAllow` in `application/config.go:733-820`
  macht Marshal-Rewrite (Z. 800), **verliert Comments** für
  den list-path — Spec-§711-721 macht keine Comment-
  Preservation-Aussage für diesen Path. Scalar-Pfade behalten
  Comments via `yaml.v3.PatchScalar` (Z. 166). T8-Doku in §6.9
  vermerkt das als bekannte Limitation; kein dedizierter
  Folge-Slice (Marshal-Rewrite ist Spec-konform).

## Bezug

- Cluster:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  (Folge-Slice 8/9).
- Pattern-Vorbilder:
  - **Modifying-Klasse**:
    [`slice-v1-cli-json-dry-run-add`](slice-v1-cli-json-dry-run-add.md)
    (RecordingFS + Diff + fsFactory),
    [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)
    (ProgressPort-Silencing, hier nicht relevant),
    [`slice-v1-cli-json-dry-run-generate`](slice-v1-cli-json-dry-run-generate.md)
    (data-Envelope-Form, per-Artefakt-Mapper),
    [`slice-v1-cli-json-dry-run-remove`](slice-v1-cli-json-dry-run-remove.md)
    (Custom-Args-Validator + Sanitizer + Pre-UC-`reportError`).
  - **Read-only-Klasse**:
    [`slice-v1-cli-json-dry-run-up-down`](slice-v1-cli-json-dry-run-up-down.md)
    (FS-Sentinel-Pattern + Mapper-Switch-Order + Sanitizer-
    Helper-Quelle + (code, exitCode)-Tupel-Disambiguation),
    [`slice-v1-cli-json-dry-run-logs`](slice-v1-cli-json-dry-run-logs.md)
    (Single-Envelope-Vertrag + Reject-Sentinel-Pattern für
    inkompatible Flag-Kombi).
- Code-Anker:
  [`cli/config.go`](../../../../internal/adapter/driving/cli/config.go),
  [`application/config.go`](../../../../internal/hexagon/application/config.go),
  [`port/driving/config.go`](../../../../internal/hexagon/port/driving/config.go),
  [`cli/jsonallowlist.go``
  Z. 29 (heutige Reject-Liste).
- Folge-Slices: noch keine direkten Forward-Refs aus config
  heraus; T0-Review kann Folge-Stubs spawnen (Kandidaten: Multi-
  Path-Set, structured-hint-Form, `config list`).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
