# Slice V1: `config --json` / `config get --json` / `config set --json` — drei Sub-Forms unter einem Folge-Slice

> **Status:** `open/`. Achter Folge-Slice (8/9) des Cluster-Slice
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
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> (RecordingFileSystem + Diff-Renderer + Custom-Args-Validator +
> Sanitizer für Pre-UC-Validation). Erbt Read-only-Klassen-
> Disziplin aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> (FS-Sentinel-Pattern + Mapper-Switch-Order +
> `cli/sanitize.go`-Helper) und
> [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
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
`DockerEngine`, KEIN `Progress`-bound-state. **Vier Sentinels
existieren bereits typed** (`ErrConfigPathUnknown`,
`ErrConfigValueInvalid`, `ErrConfigSchemaInvalid`,
`ErrConfigFileSystem`, `ErrConfigValueNotSet`); `ErrConfigFileSystem`
ist **bereits Multi-`%w`-fähig** mit Read-Message-Form
(`"config: filesystem error"`) — Pattern-Erbe up-down/logs.

Bemerkenswert: anders als up/down/logs hat `config` **schon
einen FS-Sentinel** (`ErrConfigFileSystem`, Z. 141 `port/
driving/config.go`). Das bedeutet **T2 ist substanziell
kleiner als bei up-down/logs** — kein neuer Sentinel nötig,
nur Switch-Order-Disziplin im neuen Mapper + Co-Migration
der bestehenden Wrap-Sites auf Multi-`%w` falls heute single-`%w`.

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
```

`make gates` grün (lint + test + coverage-gate ≥ 90 % +
docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ **`--json`-Allowlist-Migration**: `"u-boot config": true`,
  `"u-boot config get": true`, `"u-boot config set": true` in
  `jsonAllowlist()`; Reject-Liste schrumpft auf 1 (`template
  bare`, da `template list --json` bereits in M3 migriert ist).
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
  leeren `plannedFiles[]` + `level="info"` Diagnostic
  (T0-(d) Sub-Decision).
- ✅ **`config set` Two-Shape Summary verschwindet**: heutige
  `printConfigSetSummary` (`OldValue → NewValue` /`already X`)
  wird im JSON-Mode nicht emittiert — `data`-Carrier trägt
  die Info strukturiert (`configSetData{path, oldValue,
  newValue, noOp bool}`).
- ✅ **`config set --allow-external-feature-sources` im JSON-
  Mode**: heute via `cli/config.go:107-108` als StringSlice-
  Flag. Pre-UC-Validation-Pfad (Path-Kind-Mismatch, Z. 182-
  187) ergänzt um `reportError` analog up-down/logs-Stub.
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
  §322): jeder Envelope mit `command="config"` MUSS
  `subcommand` setzen; T6-Pin gegen Empty-Subcommand-Drift.
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
  bewiesen tragfähig.

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

- **T0-(d) `config set` NoOp-Envelope-Form**: heute returnt
  Set bei `OldValue == NewValue` ohne `WriteFile`-Call. Wie
  reagiert der Voll-Schema-Envelope?
  (i) `plannedFiles: []`, `changes: []`, `data.noOp: true`,
      `diagnostics[0].level: "info"`, `diagnostics[0].code:
      "LH-FA-CONF-001"`, `diagnostics[0].message: "already X"`.
      **Konsument-Disambiguation**: leeres `plannedFiles`
      plus `info`-Diagnostic = NoOp.
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
  | 3 | `driving.ErrConfigPathUnknown` | `LH-FA-CONF-005` | 10 | Path-Whitelist-Bruch |
  | 4 | `driving.ErrConfigValueInvalid` | `LH-FA-CONF-001` | 10 | Wert-Coercion-Bruch (nur set) |
  | 5 | `driving.ErrConfigValueNotSet` | `LH-FA-CONF-005` | 10 | Optionaler Pfad nicht gesetzt (nur get) |
  | 6 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | Pattern-Erbe up/down/generate/logs (Environment-Operation) |
  | 7 | `cli.ErrDryRunNotApplicable` (NEU, T5, falls T0-(g) (i)) | `LH-FA-CLI-006` | 2 | bare/get rejecten `--dry-run` |
  | 8 | Default (unknown) | `LH-FA-CLI-006` | 1 | Fallback |

  **Cross-Slice-Klassen-Pin**: `ErrProjectNotInitialized`
  mappt hier auf **`LH-FA-INIT-001`** (Environment-Operation
  Pattern-Erbe up/down/generate/logs) — NICHT auf
  `LH-FA-ADD-001` wie bei add/remove. Bewusste Cluster-
  Konvention.

- **T0-(g) `--dry-run`/`--diff` auf bare/get Reject-Sentinel**:
  bare/get sind Read-only → tragen kein `--dry-run`/`--diff`.
  Drei Reject-Optionen:
  (i) **Neuer `cli.ErrDryRunNotApplicable`-Sentinel** im
      CLI-Layer (`cli/config.go`) der vor UC-Aufruf rejected
      mit Exit 2. Klare Disambiguation — Konsument sieht
      "Read-only-Form, falsche Flag-Kombi".
  (ii) Re-use `cli.ErrJSONNotImplemented` (heute für
       Allowlist-Reject). **Semantischer Drift**: dieser
       Sentinel sagt "Form noch nicht migriert", nicht
       "Form unterstützt das Flag nicht".
  (iii) Cobra-Native: `MarkFlagsMutuallyExclusive` auf
        Cobra-Ebene. **Nachteil**: keine JSON-Envelope-
        Emission im Reject-Pfad (Cobra schreibt direkt
        nach stderr). Spec-§1841-Bruch.
  Plan-Empfehlung: **(i)** neuer Sentinel. Pattern-Erbe
  logs `ErrFollowJSONNotSupported` für inkompatible Flag-
  Kombi.

- **T0-(h) `subcommand`-Pflicht-Form für `config get` / `config
  set`**: bei Cobra-Compound (`u-boot config get`) trägt
  envelope.subcommand `"get"` bzw. `"set"`. Test-Pin gegen
  Empty-Subcommand-Drift. **Pattern-Erbe**: template-Slice
  hat dasselbe Problem (`template list`). Sub-Decision-Form:
  geteilter Helper `cobraPathToSubcommand(cmd) string` in
  `cli/`-Sub-Package extrahieren ODER inline-Switch im
  jeweiligen RunE.
  Plan-Empfehlung: **inline-Switch** (zwei Stellen reicht
  noch nicht für Helper-Extraktion); falls config + template
  zusammen tragen, Helper in Cluster-T_close-Tranche.

- **T0-(i) `config set` Pre-UC-Validation-Pfade**: heute
  `runConfigSet:174-187` validiert (a) Path-Parse via
  `domain.NewConfigPath`, (b) AllowExternalFeatureSources-
  Path-Kind-Mismatch. Beide sind Pre-UC-Errors und brauchen
  `reportError`-Wrap analog up-down/logs T5.

- **T0-(j) Cluster-Anker-Drift für `LH-FA-CONF-*`-Codes**:
  heutige Code-Mapping in `config.go`-Doc-Block sagt:
  - `LH-FA-CONF-001` für allgemeine Mapping (Spec-anchor
    Set-Flow).
  - `LH-FA-CONF-002` für Schema-Validation.
  - `LH-FA-CONF-005` für Path-Whitelist.
  T0-Review prüft ob die Mapper-Tabelle T0-(f) diese
  Cluster-Anker honoriert oder LH-FA-CLI-006 als generischer
  Fallback nutzt. Plan-Empfehlung: **echte Spec-Anker**
  (LH-FA-CONF-001/002/005) — analog up-down/logs die
  echte Spec-Anker (LH-FA-UP-001/INIT-001) statt
  LH-FA-CLI-006 nutzen.

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

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(k) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog up-down/logs T1): `cli/sanitize.go`-Helper, `RecordingFileSystem`-Adapter, Pure-Go-Diff-Renderer existieren bereits aus add/init/generate/remove/up-down T5 | — (entfällt) | T0 |
| T2 | Port-Types: `configFlags{JSON, Quiet, DryRun, Diff}` (Set-Form), `configGetFlags{JSON, Quiet}`, `configShowFlags{JSON, Quiet}` — ggf. zu **einem** `configFlags` konsolidiert wenn Bare/Get die Modifying-Felder ignorieren. **`ErrConfigFileSystem`-Co-Migration-Check** (T0-(e) Option (i)): Wrap-Sites auf Multi-`%w` falls heute single. KEIN neuer Sentinel. Plus: ggf. `cli.ErrDryRunNotApplicable`-Sentinel (T0-(g) Option (i)). T4 entfällt (kein Composition-Root-Wechsel). | ~60 | T0 |
| T3 | Application-Layer: Multi-`%w`-Wrap-Migration der bestehenden FS-Read/Write-Stellen auf `ErrConfigFileSystem` falls heute single-`%w`. KEIN ProgressSink-Branch nötig (config emittiert keinen Stream). KEIN Confirmer-Branch nötig (config set nicht destructive). | ~30 | T2 |
| T4 | **Entfällt** (analog up-down/logs T4): Composition-Root `cmd/uboot/main.go` hat heute schon `NewConfigService` mit allen Deps. T2 führt nur Flag-Fields ein. | — (entfällt) | T3 |
| T5 | CLI-RunE: drei `runConfig*`-Refactors auf Cluster-Signatur (ctx, stdout, errOut, args, flags, uc, getwd). Allowlist-Migration 3 Forms. Neuer `mapConfigErrorToDiagnostic` mit Switch-Order T0-(f). Pre-UC-Validation via `reportError`. **`config set` Voll-Schema-Pfad** mit `fsFactory(mode)` für Dry-Run (Pattern-Erbe add T5: `RecordingFileSystem` vs. Passthrough). **`config set --diff`-Pfad** mit Pure-Go-Diff-Renderer auf Patched-Bytes vs. Current-Bytes (Pattern-Erbe add T5). **bare/get `--dry-run`/`--diff`-Reject** via `ErrDryRunNotApplicable` (T0-(g)). Sanitizer-Aufrufe via `cli/sanitize.go`. | ~300-400 | T2 |
| T6 | Acceptance-Tests: **~18-22 Tests**: bare-Envelope-Pin (data.body), Get-Envelope-Pin (data.path/value), Set-Voll-Schema-Pin (Subcommand-Pflicht + plannedFiles + changes), Set-Dry-Run-Pin (kein WriteFile-Call), Set-Diff-Pin (hunks-Form), Set-NoOp-Pin (empty plannedFiles + info-Diagnostic, T0-(d)), Mapper-Rows 1-8, Pre-UC-Validation-Pin (Path-Kind-Mismatch für `--allow-external-feature-sources`), Sanitizer-Pin, Subcommand-Pflicht-Pin (`""`-Reject), `--dry-run`-Reject-Pin auf bare/get (T0-(g)), **FS+Schema-Multi-`%w`-Switch-Order-Defense-Pin** analog up-down/logs `_ByDesign`-Suffix (synthetisch konstruierte Multi-`%w`-Chain). | ~600-800 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe) | ~50 | T6 |
| T8 | Closure: CHANGELOG, `cli-json-output.md` **neue §6.9-Sektion** mit drei Sub-Form-Envelopes + Set-Voll-Schema-Beispiel + Subcommand-Pflicht-Doku + `--dry-run`/`--diff`-Reject-Doku für Read-only-Forms + (code, exitCode)-Tupel-Disambiguation-Block (Pattern-Erbe §6.7/§6.8); §7 keine Änderung (`config set` bereits drin). roadmap done-Zähler 7→8. carveouts.md-Einträge falls T0-Review Folge-Slices spawnt. Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: **~990-1280** (deutlich größer als
logs ~700-800 weil drei Sub-Forms gleichzeitig + Set-Modifying-
Surface mit RecordingFileSystem + Diff + Dry-Run). Pattern-Erbe
von add/init/generate/remove (Modifying-Pattern: RecordingFS +
Diff + fsFactory + Dry-Run) plus up-down/logs (Read-only-Pattern:
FS-Sentinel-Switch-Order + Sanitizer + (code, exitCode)-Tupel).

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

## Bezug

- Cluster:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  (Folge-Slice 8/9).
- Pattern-Vorbilder:
  - **Modifying-Klasse**:
    [`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
    (RecordingFS + Diff + fsFactory),
    [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
    (ProgressPort-Silencing, hier nicht relevant),
    [`slice-v1-cli-json-dry-run-generate`](../done/slice-v1-cli-json-dry-run-generate.md)
    (data-Envelope-Form, per-Artefakt-Mapper),
    [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
    (Custom-Args-Validator + Sanitizer + Pre-UC-`reportError`).
  - **Read-only-Klasse**:
    [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
    (FS-Sentinel-Pattern + Mapper-Switch-Order + Sanitizer-
    Helper-Quelle + (code, exitCode)-Tupel-Disambiguation),
    [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
    (Single-Envelope-Vertrag + Reject-Sentinel-Pattern für
    inkompatible Flag-Kombi).
- Code-Anker:
  [`cli/config.go`](../../../../internal/adapter/driving/cli/config.go),
  [`application/config.go`](../../../../internal/hexagon/application/config.go),
  [`port/driving/config.go`](../../../../internal/hexagon/port/driving/config.go),
  [`cli/jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  Z. 29 (heutige Reject-Liste).
- Folge-Slices: noch keine direkten Forward-Refs aus config
  heraus; T0-Review kann Folge-Stubs spawnen (Kandidaten: Multi-
  Path-Set, structured-hint-Form, `config list`).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
