# Changelog

All notable changes to **u-boot** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Internal `u-boot generate changelog` ([`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001-generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005-idempotenz), [`LH-AK-007`](spec/lastenheft.md#lh-ak-007-changelog-generator))
maintains a Keep-a-Changelog-formatted changelog for user projects;
this file is the same format applied to u-boot itself.

## [Unreleased]

### Added

- `feat(template): u-boot init --template ./pfad` — lokale
  User-Templates ([`LH-FA-TPL-003`](spec/lastenheft.md#lh-fa-tpl-003-eigene-templates), [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Entscheidung „Lokale
  User-Templates"). `--template` löst jetzt neben Katalog-Namen
  (`basic`) auch Dateisystem-Pfade auf (`./mein-tpl`, `/abs/tpl`,
  `~/tpl`). Die Klassifikation ist eine reine, plattformunabhängige
  `domain.ClassifyTemplateRef`-Regel (kein FS-Stat); ein Composite-
  Resolver delegiert an den eingebauten `embed.FS`-Katalog oder den
  neuen Filesystem-Resolver (`localtemplates`). Gleiche `template.yaml`-
  Validierung (apiVersion-Gate + Metadaten-Minimum) wie der Katalog
  über das geteilte `templateyaml`-Paket. Fehlerklassen: fehlender
  Pfad / kein Verzeichnis / fehlende `template.yaml` → Exit 10
  (`ErrTemplateNotFound`); malformed `template.yaml` → Exit 10
  (`ErrTemplateInvalid`); Symlink im Template-Baum → Exit 10
  (Pfad-Safety-Reject, kein Teil-Output); Render-Fehler → Exit 14.
  In einem lokalen `template.yaml` deklarierte `variables:` werden
  beim Render aktuell **ignoriert** (keine Substitution, kein
  Prompt); Variable-Auflösung (`--var key=value`) ist out-of-scope
  und folgt in einem eigenen Slice.

## [0.4.0] - 2026-06-08

Fourth release. Completes the V1 machine-readable-CLI milestone: the
full [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md) cluster delivers `--json` /
`--dry-run` / `--diff` for **all ten** spec-enum subcommands
([`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) Minimalkontrakt + [`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) Voll-Schema),
followed by the `jsonArgsValidator`-consolidation for add/init/
generate. Plus the `u-boot logs` subcommand, devcontainer-features
support with a drift-doctor check, and the transitional reject-gate
teardown at cluster T_close. [ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md)-Re-Eval-Trigger-2 (JSON-CLI
as the canonical machine interface) is satisfied. Details below.

### Added

- `feat(cli): u-boot config/config get/config set --json
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001-projektkonfiguration)..[`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005-konfiguration-anzeigen-und-ändern)) —
  achter Folge-Slice (8/9) des Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). **Erster Read-only+Modifying-
  Hybrid**: drei Sub-Formen teilen `command: "config"` mit
  Pflicht-`subcommand` (`"show"`/`"get"`/`"set"`, §322). Bare +
  `get` sind read-only (nur `--json`, Minimal+Data-Carrier
  `{body}` bzw. `{path, value}`); `set` ist modifying mit
  `--dry-run`/`--diff` (Voll-Schema, `plannedFiles[]` =
  genau eine Zeile `u-boot.yaml`, `hunks[]` aus patched-vs-
  current Bytes). **NoOp**: `oldValue == newValue` → kein
  WriteFile, `plannedFiles: []` + `data.noOp: true` + leeres
  `diagnostics: []` (kein `level: "info"`, Spec §2.1).
  **`--dry-run`/`--diff`-Reject auf den Read-only-Formen** via
  neuem `cli.ErrDryRunNotApplicable` → Exit 2 (T0-(g),
  Pattern-Erbe logs' `ErrFollowJSONNotSupported`). **Drei-Klassen-
  Sentinel-Split** (T0-(m)): `ErrConfigValueInvalid` aufgeteilt in
  `driving.ErrConfigWriteRejected` (non-writable Pfad, Hint
  `u-boot add <svc>`) + `driving.ErrConfigPostPatchSanityFailed`
  (Post-Patch-Roundtrip) — beide Exit 10, damit JSON-Konsumenten
  per `code` disambiguieren statt per Message-Substring.
  **Mapper-Tabelle mit 10 Rows** (T0-(f), FS-first): Row 1
  `ErrConfigFileSystem`→[`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/Exit 14, Rows 2-3
  Schema/Post-Patch-Sanity→[`LH-FA-CONF-002`](spec/lastenheft.md#lh-fa-conf-002-inhalt-der-konfiguration), Rows 4-5-7
  [`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005-konfiguration-anzeigen-und-ändern) (Multi-Use: Path-Unknown/Write-Rejected/
  Value-Not-Set — Disambiguation per Message-Prefix), Row 6
  `ErrConfigValueInvalid`→[`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001-projektkonfiguration), Row 8
  `ErrProjectNotInitialized`→[`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren) (Environment-
  Operation Pattern-Erbe up/down/generate/logs), Row 9 Reject
  →Exit 2, Row 10 Default. **PreviewMode-Cluster**:
  `ConfigService.fsFactory` + `selectFS` + `NewConfigServiceWith
  Factory` + `cmd/uboot/main.go`-Wiring (fünfter Preview-Factory);
  `config set` routet seinen WriteFile über die mode-spezifische
  FS und surface't `ConfigSetResponse.PlannedFiles` aus dem
  Recorder für den `--diff`-Renderer. **SilenceLogger** silenced
  die fünf `s.logger.*`-Sites im JSON-Mode; **Orphan-Feature-WARN**
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)/`level: "warn"`) migriert in `diagnostics[]`
  (Dual-Emission). **Subcommand-bewusste `reportErrorSub`/
  `writeErrorEnvelopeSub`** (additiv; Single-Form-Caller
  unverändert) erfüllen die §322-Pflicht auch auf dem Error-Pfad.
  Allowlist-Reject-Liste schrumpft von 4 auf 1 (nur noch
  `template (bare)`). Drei Review-Runden (zwei HIGH + ein MED, alle
  gefixt). Doku in [`docs/user/cli-json-output.md §6.9`](docs/user/cli-json-output.md).
- `feat(cli): u-boot logs --json ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run) /
  [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)) — siebter Folge-Slice
  (7/9) des Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md).
  **Read-only-Klasse** auf lokalem FS (analog up-down):
  weder `--dry-run` noch `--diff` — nur `--json` mit
  typisiertem Data-Carrier. **`logsStatusData.lines []string`
  ohne omitempty** (Empty-Array-Pin: `[]`, NICHT `null`,
  Pattern-Erbe up-down's `services []serviceStatus` ohne
  omitempty). **T0-(a) Single-Envelope + `--follow --json`
  Reject** (Option (A)): Spec-§1841-Konsens (Single-Envelope
  pro CLI-Call) wird honoriert; `--follow --json` ist
  inkompatibel und wird in `runLogs` Stage-1 (vor UC-Call)
  mit neuem `ErrFollowJSONNotSupported` → [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/
  Exit 2 rejected — bounded `--tail=N`-Pfad ist die
  einzige Akquisitions-Form. **T0-(i) Validation-Order**:
  `--follow --json` schlägt `--tail=-1`; CLI-Stage-1
  reihenfolge pinnt Reject-Sentinel VOR Tail-Validation
  (`TestLogsJSON_ValidationOrder_FollowJSONBeatsInvalidTail`
  Pin). **Neuer FS-Sentinel** `driving.ErrLogsFileSystem`
  mit Read-spezifischer Message-Form (`"logs: filesystem
  read failed"`, Pattern-Erbe up-down). T3 wrapt zwei
  FS-Read-Stellen (`logsservice.go:117/137` —
  `checkProjectInitialized` + `checkComposeFile`) auf
  Multi-`%w` (`fmt.Errorf("logs service: Exists(%q): %w: %w",
  path, driving.ErrLogsFileSystem, err)`). **Mapper-Tabelle
  mit 9 Rows** (T0-(f)): Row 1 FS-Sentinel-first (LH-NFA-REL-
  003/Exit 14), Rows 2-3 shared `mapComposeRuntimeSentinel`-
  Helper (Docker/ComposeRuntime → [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)), Row 4
  shared `ErrComposeFileMissing` ([`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001-umgebung-starten)/Exit 10),
  Row 5 cross-cutting `ErrProjectNotInitialized` (LH-FA-INIT-
  001/Exit 10, Pattern-Erbe generate als Environment-
  Operation), Row 6 domain-level `ErrInvalidServiceName`
  ([`LH-FA-INIT-006`](spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung)/Exit 10), Row 7 logs-only
  `ErrFollowJSONNotSupported` ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/Exit 2 — T0-(a)
  Reject-Pfad), Row 8 logs-only `ErrInvalidLogsTail`
  ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/Exit 2), Row 9 Default [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/Exit 1.
  **Switch-Order-FS-first-Defense** (Pattern-Erbe up-down
  R2-HIGH-2 / R3-HIGH-1) via `TestLogsJSON_MultiWrap_
  FSAndDocker_SwitchOrderFSFirst_ByDesign`-Pin — synthetische
  FS+Docker-Chain → diagnostics[0].code = [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)
  (FS-Klasse, Mapper-FS-first), `exitCode = 11` (Docker-
  Sub-Klasse, ExitCode-Helper-Driven-first); **`(code,
  exitCode)`-Tupel-Disambiguation** per `cli-json-output.md
  §6.7` ist der Vertrag (Pattern-Erbe up-down T8).
  **`baseDirSanitizedError`-Wiederverwendung**: `runLogs`
  wrappt UC-Errors mit `sanitizeBaseDir(err, cwd)` vor
  `reportError` (Path-Leak-Defense, Pattern-Erbe up-down
  T5). **`runLogs(ctx, stdout, errOut io.Writer, args,
  flags, uc, getwd)`-Signatur** Cluster-konsistent mit
  up/down/remove (errOut für strukturierte Pfade reserviert,
  heute `_ = errOut`-Stub). **`logsFlags{Follow, Tail,
  Service, JSON, Quiet}`** mit neuen `JSON`/`Quiet`-Boolean-
  Feldern und `IsValid()`-Builder (Pattern-Erbe up's
  `upFlags`). **Cluster-Allowlist** erweitert
  (`jsonallowlist.go`): `"u-boot logs": true`; Reject-Liste
  von 5 → 4 (`config bare/get/set`, `template bare` bleiben).
  **`cli.isFilesystemError`** erweitert um
  `driving.ErrLogsFileSystem` → Exit 14;
  **`cli.isUsageError`** erweitert um
  `cli.ErrFollowJSONNotSupported` → Exit 2 (`cli.go`).
  **15 Acceptance-Pins** (`logs_acceptance_test.go`):
  T0-(a)/T0-(i)/T0-(j)(ii) verbatim, 9 Mapper-Coverage-
  Pins, Empty-Array-Pin, Trailing-Newline-Strip-Pin, Path-
  Leak-Sanitizer-Pin, FS+Docker-Switch-Order-Defense-Pin.
  Pre-T6-Review: HIGH=0, MED=4, LOW=6 (T7 fixte MED-1 +
  LOW-5; MED-3 → §6.8-Doku-Pflicht hier eingelöst).
  Pre-T8-Bestätigungsrunde: HIGH=0, MED=2, LOW=2 (MED-1
  Mapper-Kommentar-Drift + MED-2 Defense-Pin + LOW-1 Plan-
  Drift in `ba7d06f` gefixt; LOW-2 CRLF-Lücke in §6.8 als
  bekannte Limitation dokumentiert). **Vier neue
  open/-Stubs** (T6 R2-LOW): [`slice-v1-logs-format-flags`](docs/plan/planning/open/slice-v1-logs-format-flags.md),
  [`slice-v1-logs-multi-service-filter`](docs/plan/planning/open/slice-v1-logs-multi-service-filter.md),
  [`slice-v1-logs-time-range-filter`](docs/plan/planning/open/slice-v1-logs-time-range-filter.md).
  `[ba7d06f, b502cd5, 343e622, 69cfc0d,
  c21ba28, 0fe74e4]`.

- `feat(cli): u-boot up --json / u-boot down --json
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001-umgebung-starten)/[`LH-FA-UP-003`](spec/lastenheft.md#lh-fa-up-003-startstatus-anzeigen)/[`LH-FA-UP-004`](spec/lastenheft.md#lh-fa-up-004-umgebung-stoppen) /
  [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung)) — sechster Folge-Slice (6/9) des Cluster-
  Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). **Read-only-Klasse** auf
  lokalem FS: weder `--dry-run` noch `--diff` (Cluster-Slice
  Z. 464-467) — nur `--json` mit typisierten Data-Carriern.
  `u-boot up` und `u-boot down` sind im selben Slice gebündelt
  weil beide den Compose-Status lesen und das Confirmer-Swap-
  Pattern teilen. **`upStatusData.services[]`** trägt
  `serviceStatus{name, state, port, healthcheck}` (plain Go-
  Strings mit `omitempty` für port/healthcheck — keine Three-
  State-Disambiguation nötig; [`LH-FA-UP-003`](spec/lastenheft.md#lh-fa-up-003-startstatus-anzeigen) Mindestangaben);
  `data.timeoutFireAndForget *bool omitempty` als Marker nur im
  `--timeout=0`-Pfad (Pattern-Erbe remove's `*bool`-Key-
  Absence-Disambiguation). **`downStatusData.removedVolumes
  bool`** ohne omitempty (`false` ist legitimer Success-Wert
  "nichts entfernt"). **`UpRequest.SilenceProgress bool` +
  `DownRequest.SilenceConfirmer bool`** symmetrisch zum
  remove-Pattern; CLI setzt sie auf `flags.JSON`. **Application-
  Layer-ProgressSink-Branch** (`UpService.Up`): `effective :=
  req.ProgressSink; if req.SilenceProgress { effective =
  io.Discard }` — Compose-Phase-Stream wird im JSON-Mode
  unterdrückt, nil-Default bleibt im DockerEngine-Adapter
  (`progressSinkOrDiscard`-Pattern). **Application-Layer-
  Confirmer-Branch** (`DownService.runConfirmationGate` Row 4):
  Request-time Gate-Branch ohne Field-Mutation (`confirmer :=
  s.confirmer; if req.SilenceConfirmer { confirmer =
  noopConfirmer{} }`) — kein neuer `downMu`-Mutex nötig, race-
  frei by construction. **Refuse-by-Default-Semantik** im
  JSON-Mode: bei `--volumes --json` OHNE `--yes` returnt
  `noopConfirmer.ConfirmRemoveVolumes` `(false, nil)` → fällt
  durch in `ErrConfirmationRequired`/Exit 10 (Symmetrie zum
  `--no-interactive`-Pfad). JSON-Konsumenten MUSSEN `--yes`
  explizit setzen für destructive `--volumes`. **Zwei neue
  FS-Sentinels** `driving.ErrUpFileSystem` /
  `driving.ErrDownFileSystem` mit Read-spezifischer Message-Form
  (`"<cmd>: filesystem read failed"`, NICHT `"mutation
  failed"` weil up/down read-only auf lokalem FS). T3 migriert
  fünf FS-Read-Wraps (`upservice.go:105/138/148` +
  `downservice.go:81/97`) auf Multi-`%w`. **Mapper-Tabelle mit
  verbindlicher Switch-Order** (T0-(e) R3-HIGH-1): zehn Rows
  mit Mapper-Heim-Spalte (`mapUp`/`mapDown`/`helper`/`beide`).
  Row 1 FS-Sentinel-first ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/Exit 14), Rows 2-3
  shared `mapComposeRuntimeSentinel(err)`-Helper in neuem File
  `cli/composesentinel.go` ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) für `driven.ErrDocker
  Unavailable` und `driven.ErrComposeRuntime`), Row 4 up-only
  `ErrStabilizationTimeout`/[`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001-umgebung-starten)/Exit 12, Row 5 down-
  only `ErrConfirmationRequired`/[`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005-überschreibschutz)/Exit 10
  (geteilt mit init/remove), Row 6 shared `ErrComposeFileMissing`
  /[`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001-umgebung-starten), Row 7 cross-cutting `ErrProjectNotInitialized`/
  **[`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)** (Pattern-Erbe generate als Environment-
  Operation, NICHT add/remove [`LH-FA-ADD-001`](spec/lastenheft.md#lh-fa-add-001-add-on-befehl)), Row 8 up-only
  `ErrInvalidTimeout`/[`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes), Row 9 down-only
  `ErrConflictingModeFlags`/[`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung), Row 10 Default
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes). **`(code, exitCode)`-Tupel-Disambiguation**
  (cli-json-output.md §6.7-Doku-Pin): bei Multi-`%w`-Wraps
  liefern Mapper (FS-first) und ExitCode-Helper (Driven-first
  per `cli.go:285-313`) zwei getrennte Klassifikationen —
  Konsument disambiguiert über das Tupel, NICHT über `code`
  allein. Pattern-Erbe remove's [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007-service-entfernen) Multi-Use.
  **`baseDirSanitizedError`-Helper-Extraktion**: aus
  `cli/remove.go:465-538` nach neuem File `cli/sanitize.go`
  (package-intern; `remove.go:299` nutzt unverändert weiter).
  `runUp`/`runDown` wrappen UC-Errors mit `sanitizeBaseDir(err,
  cwd)` vor `reportError` — 11 FS-Read- und Compose-Runtime-
  Wraps in upservice/downservice tunneln keinen absoluten
  Filesystem-Pfad mehr in `diagnostic.message`. **Allowlist-
  Migration**: `u-boot up` + `u-boot down` als zwei separate
  Einträge in `jsonAllowlist()` (Cobra-CommandPath-Form);
  Reject-Liste schrumpft von 7 auf 5 (logs, config bare/get/set,
  template bare). Tranchen-Reihenfolge: T1 + T4 entfallen
  (noopConfirmer/io.Discard existieren als Helper; Composition-
  Root braucht kein Wiring-Update bei Bool-Field-Pattern); T2
  Port-Types (`e966a83`), T3 Application-Layer (`86fb5b2`), T5
  CLI-RunE + Mapper + Sanitizer-Helper-Extraktion (`a5aaf9c`),
  T6 28 Acceptance-Tests (`2473988`), T7 Pre-T8-Adressierung
  von 11 Findings (`31f7238`), T8 Closure. **Acceptance-
  Coverage**: 28 CLI-Tests (18 up + 10 down) plus 2 Application-
  Layer-Tests in `downservice_test.go` für die `noopConfirmer`-
  Branch-Defense (`removeVolumesCalls == 0` + Contrast-Pin).
  **Out-of-Scope-V1 Carveouts mit open/-Stubs**: Recreate-
  Detection ([`slice-v1-recreate-detection`](docs/plan/planning/open/slice-v1-recreate-detection.md)), Volume-Named-Liste
  ([`slice-v1-down-volumes-named-list`](docs/plan/planning/open/slice-v1-down-volumes-named-list.md)), Partial-Snapshot bei
  Mid-ComposeUp-Failure
  ([`slice-v1-up-partial-snapshot-on-failure`](docs/plan/planning/open/slice-v1-up-partial-snapshot-on-failure.md)),
  strukturierte Multi-Port-Liste
  ([`slice-v1-multi-port-services`](docs/plan/planning/open/slice-v1-multi-port-services.md)).
  Coverage-Gate 91 %.
- `feat(cli): u-boot remove --json / --dry-run / --diff
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007-service-entfernen) /
  [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung)) — fünfter Folge-Slice (5/9) des Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). `u-boot remove` ist die **inverse
  Operation zu `add`** (strip managed-block aus `compose.yaml` +
  `.env.example`, flip `services.<name>.enabled` auf `false` in
  `u-boot.yaml`, optional Volume-Purge via `--purge`-Gate). Acht
  Flag-Kombinationen plus `--purge`-Dimension: drei JSON-Modi
  (Minimal+Data, Voll-Schema Dry-Run, Voll-Schema Preview-and-
  Apply) plus Human-Mode mit/ohne `--diff`/`--purge`. Pattern-
  Erbe von add/init/generate 1:1: `driving.PreviewMode` direkt
  (kein Service-Prefix-Alias), `RecordingFileSystem`-driven-
  Adapter, Pure-Go LCS-Diff-Renderer, `previewModeFromFlags`-
  Mapping, generalisierte Error-Emission-Helper. Remove-spezifisch:
  **Confirmer-Gate für `--purge`** ([`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung) §254) mit
  `noopConfirmer`-Swap im JSON-Mode (T0-(j) — neues Pattern,
  nicht aus init-Progress-Swap geerbt: Service-Field-Mutation
  mit defer-Restore innerhalb `removeMu`-Lock-Scope; lokale-
  Variable-Variante verworfen weil `runPurgeGate`-Signature-
  Refactor); **`data.volumesPurged`** als `*bool` im Envelope
  (`false` deferred-Status in v0.3.0 ist valider Success-Wert,
  Plain-`bool`+omitempty würde Error-Pfad-Zero und Success-Pfad-
  `false` identisch droppen); **WARN-Migration**: heutige
  `printRemoveSummary`-stderr-WARNING bei `--purge && !VolumesPurged`
  wandert im JSON-Mode in `diagnostics[]` mit
  code: "[`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007-service-entfernen)", `level: "warn"` (Multi-Use des Codes
  für ERROR `ErrServiceUnregistered` UND WARN deferred-Volumes —
  Konsumenten disambiguieren über `(code, level)`-Tupel);
  **Custom-`Args`-Validator** `validateRemoveArgs(a *App)` als
  Cobra-PositionalArgs-Closure mit `*App`-Capture (R11/R12/R13-Pin),
  ersetzt `cobra.ExactArgs(1)` — emittiert den
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)-Envelope auf stdout BEVOR der Sentinel zu Cobra
  zurückgeht (Spec §1841-Symmetrie); **Voll-Schema bei `--dry-run`/
  `--diff`** auch im NoPositionalArg- und TooManyArgs-Pfad
  (Spec §1842, R13-HIGH-1 Validator-Flag-Awareness via
  `cmd.Flags().GetBool`); **`baseDirSanitizedError`-Wrapper** für
  `diagnostic.message`: FS-Wraps der Form
  `fmt.Errorf("... %s: %w: %w", absPath, ErrRemoveFileSystem, raw)`
  tunneln den absoluten Filesystem-Pfad in den User-facing Output —
  Sanitizer ersetzt `<baseDir>/foo` durch `foo` und bare `<baseDir>`
  durch `.`, an Word-Boundaries (`replaceBareBaseDir` ist robust
  gegen Substring-Kollisionen wie `<baseDir>-cache/lock`,
  R14-MED-1 + R15-LOW-1); `errors.Is`/`As` bleiben intakt via
  Unwrap-Chain. **Dry-Run-WARN-Suppression** in
  `printRemoveSummary`: Use-Case skippt `runPurgeGate` in
  `PreviewDryRun` (T0-(h)(a)) und führt keine Mutation aus — die
  WARN-Prosa wäre semantisch falsch ("ist-deferred" statt
  "würde-deferred"); Fix unterdrückt WARN-Block bei
  `previewMode == PreviewDryRun`, `PreviewAndApply` behält die
  WARN. Remove-spezifische LH-Code-Mapper-Tabelle in
  `mapRemoveErrorToDiagnostic` (FS-first Switch-Order):
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/14 für `ErrRemoveFileSystem`, [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung)/10
  für `ErrConfirmerUnavailable` (neuer Sentinel) UND
  `ErrConflictingModeFlags`/2, [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005-überschreibschutz)/10 für
  `ErrConfirmationRequired`, `LH-FA-ADD-{001,002,005,007}`/10 für
  fachliche Service-Sentinels, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/2 für
  `ErrServiceNameMissing` und Cobra-too-many-args. **Idempotenz-
  NoOp-Semantik**: nur `PriorState=Deactivated` qualifiziert
  (Single+Repeat-Call → `plannedFiles: []`, `changes: []`,
  `data.priorState=data.state="deactivated"`, Exit 0);
  `EnabledUnset` und `InconsistentBlock` sind state-transitioning
  (`Changed!=nil`, Voll-`plannedFiles[]`). **`delete`-Action im
  Recorder**: `RemoveAll`-Captures für `extraFiles` werden als
  `plannedFiles[].action: "delete"` gewired — `--diff` rendert
  Old-Content + leeren New-Content. Tranchen-Reihenfolge (T1
  entfällt — `noopConfirmer` lebt bereits in M4 Confirmer-Slice):
  T2 Port-Types (`d0c9c5d`), T3 Application-Layer mit 8 FS-Wrap-
  Stellen Multi-`%w` (`dbbf7b1`), T4 `newPreviewFSFactory`-Helper
  (`3b079dd`), T5 CLI-RunE-Rewrite (`3188e75`), T6+T6-A
  Acceptance-Tests inkl. `WithDataKeyAbsent`/`WithDataKeyPresent`-
  Helper-Erweiterung in `jsontestutil` (`9eae9ec`), T7 Pre-T8-
  Review-Fixes R13-HIGH-1 + R13-MED-1 + R14-HIGH-2 + R14-MED-1
  (`4fb3fea`), T8 Closure mit R15-LOW-1 Sanitizer-Substring-
  Robustheit + R15-LOW-2 Coverage-Pin. Acceptance-Coverage:
  23 Pin-Tests in `remove_acceptance_test.go` (drei JSON-Modi,
  Idempotenz-Repeat, EnabledUnset-Normalisierung, Mid-Write-
  Failure mit `plannedFiles[]`-Capture, ConfirmationRequired-
  Pfad, ServiceUnregistered ERROR vs. deferred-Volumes WARN,
  Concurrent-Invocations-Mutex, Dry-Run-WARN-Suppression mit
  PreviewAndApply-Kontrast, Sanitizer-Path-Leak-Defense plus
  Substring-Kollisions-Robustheit, TooManyArgs-Voll-Schema bei
  `--dry-run`, `delete`-Action-Hunk-Vertrag). Coverage-Gate
  91.10 %.
- `feat(cli): u-boot generate --json / --dry-run / --diff
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001-generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005-idempotenz) /
  [`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)/[`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)) — vierter Folge-Slice (4/9) des Cluster-
  Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md). `u-boot generate` ist nach
  doctor/add/init der nächste modifying-Subcommand und der
  **erste**, der mehrere Artefakte (changelog/readme/env-example/
  devcontainer) über einen einzigen Subcommand bedient. Vier
  Flag-Kombinationen mit derselben Symmetrie wie add/init:
  `--json` (Minimal+Data-Envelope), `--dry-run --json`
  (Voll-Schema mit `plannedFiles[]`/`changes[]`, kein FS-Write),
  `--diff --json` (Voll-Schema mit Hunks, Preview-and-Apply),
  `--dry-run --diff --json` (Vorschau plus Hunks, kein Write).
  Human-Mode `generate --diff` rendert Unified-Diff am stdout.
  **Action-Klassifikation via `data.action`** (`created` /
  `updated-block` / `no-op` / `repaired-manual`):
  UpdatedBlock und RepairedManual sind FS-semantisch identisch
  (`plannedFiles[i].action: "modify"`), `data.action` ist der
  einzige Discriminator. **Multi-Artefakt-Envelope-Form**:
  `command="generate"`, kein `subcommand`-Feld (Cobra-Positional-
  Arg-Semantik), `data.artifact` trägt das Artefakt. Pattern-
  Erbe init→generate 1:1: `PreviewMode`-Carrier (direkt, kein
  Service-Prefix-Alias), `RecordingFileSystem`-driven-Adapter,
  Pure-Go LCS-Diff-Renderer, `previewModeFromFlags`-Mapping,
  generalisierte Error-Emission-Helper. Generate-spezifische
  Erweiterungen: **2 von 8 Recorder-Mutations-Methoden** im
  Capture-Set (`WriteFile`, `MkdirAll`); **kein** `GitClient`
  und **kein** `ProgressPort` (schmaler als init);
  **`generateMu sync.Mutex`** auf `GenerateService`;
  **per-Artefakt LH-Code-Tabelle** im neuen
  `mapGenerateErrorToDiagnostic(err, artifact)` (changelog→
  [`LH-FA-GEN-002`](spec/lastenheft.md#lh-fa-gen-002-changelog-erzeugen), readme→[`LH-FA-GEN-003`](spec/lastenheft.md#lh-fa-gen-003-readme-erzeugen), env-example→
  [`LH-FA-GEN-004`](spec/lastenheft.md#lh-fa-gen-004-beispiel-env-erzeugen), devcontainer→[`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen));
  **`ErrConfigValueInvalid`-Sentinel-Wrap** auf
  `validateAllowExternalFeatureSourcesEntries`/
  `applyAllowExternalFeatureSources` für den
  [`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)-URL-Reject-Pfad (Spec §720 fordert exakt
  Exit 10 für ungültige `--allow-external-feature-sources`-URLs;
  ohne den Sentinel-Wrap wäre der Pfad auf Default
  [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)/Exit 1 gefallen); **Multi-`%w`-Wrap auf
  ~17 FS-Wrap-Stellen** in `application/generate.go` (Switch-
  Order-Sicherheit analog init T6; ohne Multi-`%w` würde ein
  Multi-Wrap mit FS-Sentinel + fachlich-Sentinel auf Exit 10
  downgraden). **`cliJSONEnvelope.Data`-Feld + `newDataEnvelope`-
  Konstruktor** wurden aus dem Template-Slice 9/9 in generate
  vorgezogen (Generate ist der erste Multi-Artefakt-Konsument
  mit `data`-Bedarf); Template-Slice 9/9 erbt das Feld nur noch.
  **`writeErrorEnvelope`/`reportError` um `data any`-Trailing-
  Param erweitert** (init/add reichen `nil` durch — nicht-
  brechende Erweiterung). Acceptance-Coverage: 15 Tests in
  `generate_acceptance_test.go` (drei JSON-Modi, 4 ManualConflict-
  Codes als Sub-Tests, URL-Reject-[`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003-devcontainer-features), ArtifactUnknown-
  Exit-2, ProjectNotInitialized-Exit-10, FS-Failure-Exit-14,
  Allow-External-Mutex, NoOp-Empty-Arrays, UpdatedBlock-vs-
  RepairedManual-Action-Discriminator, Human-Mode-Summary +
  Diff-Rendering). Devcontainer-Phase-1-Atomicity + Phase-2-
  Half-Write-Carveout (V2-Open-Slice
  [`slice-v2-generate-devcontainer-rollback-aware-write`](docs/plan/planning/open/slice-v2-generate-devcontainer-rollback-aware-write.md)) +
  Repeat-Idempotency leben in den Application-Layer-Tests.
  Coverage-Gate ≥ 91 %.

### Changed

- **BREAKING** `u-boot template list --json` Ausgabe-Format
  ([slice-v1-cli-json-dry-run-template](docs/plan/planning/done/slice-v1-cli-json-dry-run-template.md), neunter und letzter
  Folge-Slice des Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md)):
  der bisherige rohe, pretty-indented `[]templateJSON`-Array-
  Output wurde auf den [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe)-Minimalkontrakt-Envelope
  migriert — `{"status":"ok","command":"template","subcommand":
  "list","diagnostics":[],"exitCode":0,"data":[…]}` (single-line,
  compact). Die Template-Liste lebt jetzt im `data`-Feld;
  Konsumenten, die das Top-Level-Array lasen, müssen auf `.data`
  umstellen. Damit ist das letzte nicht-spec-konforme `--json`-
  Surface geschlossen (alle neun Folge-Slices des Clusters tragen
  jetzt den Minimalkontrakt; der bewusste Doctor-Slice-Carveout ist
  aufgelöst). bare `u-boot template --json` bleibt Exit-2-Reject
  (§1838: `subcommand` verpflichtend für `command="template"`;
  Help-Parent ohne eigenes Datum). Ein Katalog-IO-Fehler mappt auf
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)/Exit 14. Doku in
  [`docs/user/cli-json-output.md §6.2`](docs/user/cli-json-output.md).
- **Cluster-T_close** (Abschluss des [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md)-
  Clusters): die transitionale `--json`-Reject-Mechanik (Allowlist-
  Map + `applyJSONRejectGate` am Root-`PersistentPreRunE` +
  `ErrJSONNotImplemented`) wurde **entfernt** — alle Spec-Enum-Forms
  sind migriert, es gibt nichts mehr zu rejecten. Der bare-
  `u-boot template --json`-Reject ist von gate-getragen
  (`ErrJSONNotImplemented`) auf **RunE-getragen**
  (`ErrTemplateSubcommandRequired`, Exit 2, envelope-LOS) umgestellt,
  damit er den Gate-Abbau überlebt ohne Hilfetext zu leaken. Keine
  Verhaltensänderung für Konsumenten (Exit 2 bleibt); rein interner
  Mechanik-Abbau (netto −90 LOC). Doku in
  [`docs/user/cli-json-output.md §6.1`](docs/user/cli-json-output.md).

### Fixed

- **`add`/`init`/`generate` `--json` Envelope-Symmetrie + Path-Leak-
  Defense** ([slice-v1-cli-json-envelope-consolidation](docs/plan/planning/done/slice-v1-cli-json-envelope-consolidation.md), R15-Cross-
  Slice-1): ein Wrong-Arg-Aufruf unter `--json` (`u-boot --json add`
  ohne Service, `u-boot --json add a b` mit zu vielen) emittierte
  bisher nur eine nackte Cobra-stderr-Meldung und **kein** JSON auf
  stdout (Spec §1841-Verletzung); `--dry-run`/`--diff` wählten zudem
  nicht das Voll-Schema (§1842). Jetzt tragen alle drei Commands den
  Args-Envelope (Exit 2) über den geteilten `jsonArgsValidator` —
  konsolidiert mit `config`/`remove`, die das Muster schon hatten.
  Zusätzlich: absolute Filesystem-Pfade aus Use-Case-Fehler-Wraps
  werden nun via `sanitizeBaseDir` aus `diagnostic.message`
  entfernt (Path-Leak/Info-Disclosure-Defense, symmetrisch zu
  `config`/`remove`). Kein Verhaltenswechsel für korrekte Aufrufe.
- `fix(cli): mapAddErrorToDiagnostic Backup-Sentinels auf
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) (add↔init Diagnostic-Code-Harmonisierung) —
  `mapAddErrorToDiagnostic` mappte `ErrBackupSuffixExhausted`/
  `ErrBackupSourceMissing` auf [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005-überschreibschutz) (Validation-Klasse,
  Exit-10-Suggestion), während `isFilesystemError` sie ohnehin
  zu Exit 14 routet — Envelope-Code und Exit-Klasse waren also
  desynchron. `mapInitErrorToDiagnostic` macht es bereits korrekt
  ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) + Exit 14). Cleanup zieht add auf dieselbe
  Klassifikation nach. Die Branch ist im add-Pfad defensiv (heute
  ruft kein add-Use-Case `runBackup`/`BackupPath`), bleibt aber
  für zukünftige Catalog-Erweiterungen erhalten. Cross-Slice-
  Drift-Finding aus [`slice-v1-cli-json-dry-run-init`](docs/plan/planning/done/slice-v1-cli-json-dry-run-init.md) Review-Round-9
  (`d7f9e65`); Plan-Pfad
  [`slice-v1-cli-cleanup-add-backup-error-class`](docs/plan/planning/done/slice-v1-cli-cleanup-add-backup-error-class.md).
  Coverage-Gate unverändert grün (91.10 %).

### Added

- `feat(cli): u-boot init --json / --dry-run / --diff
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe)) — zweiter modifying-Sub-
  command mit JSON-Envelope-Migration und wichtigster Onboarding-
  Use-Case (Cluster-Slice [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice
  3/9). Vier neue Flag-Kombinationen mit derselben Symmetrie wie
  `add`: `--json` (Minimalkontrakt-Envelope), `--dry-run --json`
  (Voll-Schema mit `plannedFiles[]`/`changes[]`, kein FS-Write),
  `--diff --json` (Voll-Schema mit `plannedFiles[].hunks[]`,
  Preview-and-Apply), `--dry-run --diff --json` (Vorschau plus
  Hunks, kein Write). Human-Mode `init --diff` rendert Unified-
  Diff am stdout. Pattern-Erbe von add 1:1: gemeinsame
  `PreviewMode`-Carrier-Types (`AddPreviewMode` als Alias auf
  kanonischen `driving.PreviewMode`), `RecordingFileSystem`-
  driven-Adapter, Pure-Go LCS-Diff-Renderer,
  `previewModeFromFlags`-Mapping, generalisierte Error-Emission-
  Helper (`reportError`/`writeErrorEnvelope`/`writeDiff`/
  `lastPlannedPath`) in `cli/erroremission.go`. Init-spezifische
  Erweiterungen: **sechs der acht** `driven.FileSystem`-
  Mutations-Methoden im Capture-Set (`MkdirAll`, `WriteFile`
  direkt plus `CopyExclusive`/`Mkdir`/`MkdirAll`/`Copy`/
  `RemoveAll` indirekt via `BackupPath`); **`initGit`-Skip im
  Dry-Run** (Composition-Root-`initFSFactory`-Closure liefert
  Recorder; `Init()` skippt den separaten `driven.GitClient`-Port
  bei `PreviewMode == PreviewDryRun`, weil git am Recorder vorbei
  auf die echte Disk schreibt); **ProgressPort-Silencing im
  JSON-Mode** (`req.SilenceProgress = flags.JSON` swappt
  `s.progress` auf Noop, weil add keinen stdout-bound Port hat
  und init der erste mit Progress-Events während des Use-Case-
  Laufs ist); **Service-Race-Safe via `sync.Mutex`** auf
  `InitProjectService.initMu` (analog `AddServiceService`).
  Sieben init-spezifische LH-Codes als Spec-Anker
  ([`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung)),
  drei mit dedizierten Sentinels in
  `mapInitErrorToDiagnostic` ([`LH-FA-INIT-004`](spec/lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) Marker-Kollision,
  [`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005-überschreibschutz) Force/Backup-Usage, [`LH-FA-INIT-006`](spec/lastenheft.md#lh-fa-init-006-projektnamen-validierung) Name-
  Validierung); **neuer `driving.ErrInitFileSystem`-Sentinel**
  als Multi-`%w`-Wrap auf fünf FS-Wrap-Stellen (vier direkte in
  `initproject.go` — `MkdirAll` + drei `WriteFile`-Actions — plus
  ein Caller-side Wrap auf das `runBackup`-Ergebnis, der typed
  Backup-Sentinels durchreicht und nur rohe FS-Errors wrappt;
  R3-Wrap-Strategie konsolidiert die ursprünglich pro
  `backup.go`-Site geplanten Wraps auf den Aufruferand).
  Switch-Order-Pflicht: FS-first, weil Multi-`%w` sonst Exit-14
  auf Exit-10 downgraded. Template-Modus (`init --template <name>`) ist in
  V1 **mutex** zu `--dry-run`/`--diff` —
  `driving.ErrTemplateConflictsWithFlag` rejected die Kombination
  (Exit-Code 2, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)); Template-Preview wandert in
  eigenen Folge-Slice. Context-Cancellation bleibt
  Status-quo-Carveout (Cross-Cutting Exit-130-Convention für alle
  modifying-Subcommands ist eigener Block). Acceptance-Coverage:
  ~17 Flag-Matrix-Tests inkl. Soft-Existing × `--devcontainer`,
  Planning-Phase-Force-Failure (Exit 10), Mid-Write-Failure mit
  Switch-Order-Pin (Exit 14), Concurrent-Init-Mutex-Pin auf zwei
  Goroutinen, Path-Anchor-Pin für positional `<name>` mit
  trailing-slash/dot-slash/abs-path, initGit-Skip-Pin (`.git/`
  fehlt + Spy-Counter 0), JSON-stdout-Cleanliness-Pin
  (`json.Decode → io.EOF`). Coverage-Gate ≥ 91 %.

- `feat(cli): u-boot add --json / --dry-run / --diff
  ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe) / [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe)) — erster modifying-Sub-
  command mit JSON-Envelope-Migration (Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice 2/9). Vier neue
  Flag-Kombinationen: `--json` (Minimalkontrakt-Envelope ohne
  Plan), `--dry-run --json` (Voll-Schema mit `plannedFiles[]`/
  `changes[]`, kein FS-Write), `--diff --json` (Voll-Schema
  Preview-and-Apply mit `plannedFiles[].hunks[]`),
  `--dry-run --diff --json` (Vorschau plus Hunks, kein Write).
  Human-Mode-`--diff` rendert Unified-Diff-String an stdout
  (`+`/`-`/space-Prefix plus `@@`-Header). Cluster-Infrastruktur
  als Pattern-Vorbild: neuer `RecordingFileSystem`-driven-Adapter
  (`internal/adapter/driven/recordingfs/`) implementiert alle 8
  Mutations-Methoden mit Passthrough-Schalter, modelliert den
  impliziten `MkdirAll`-Effekt auf Parent-Dirs; Pure-Go LCS-Diff-
  Renderer (`internal/adapter/driving/cli/diff/`); Composition-
  Root-`fsFactory(driving.AddPreviewMode)`-Closure in
  `cmd/uboot/main.go`. Diagnostic-Codes sind LH-Kennungen
  (`LH-FA-ADD-{001,002,005,006}`/`LH-FA-INIT-{004,006}`/
  [`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)). Mid-Write-Failure-UX: Voll-Schema-Envelope
  zeigt `plannedFiles[]` bis zur Failure-Stelle, `diagnostics[].file`
  markiert die Failure-Position, `exitCode: 14` ([`LH-NFA-REL-003`](spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)).
  `changes[].count`-Semantik: für `create` total Lines der neuen
  Datei, für `modify` `+`-Lines aus den Hunks (Spec §477 exakt),
  für `delete` `0`, für binary `CountBytesDiff`. Service-Race-
  Safe via `sync.Mutex` auf `AddServiceService`; recursive
  Dep-Installs (für künftige Catalogue-Erweiterungen) erben den
  Outer-`PreviewMode`. AssertFullEnvelope erweitert um
  `checkHunks`-Helper mit Field-Name-Drift-Pin
  (`offset` statt `oldStart` failt sofort).

- `feat(cli): u-boot doctor --json plus Root-PersistentFlag
  --json ([`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) / [`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)) — Pattern-Vorbild für
  die maschinen-lesbare CLI-Surface (Cluster-Slice
  [`slice-v1-cli-json-dry-run`](docs/plan/planning/done/slice-v1-cli-json-dry-run.md), Folge-Slice 1/9). **Doctor**
  emittiert mit `--json` einen Spec-§1841-Minimalkontrakt-
  Envelope (`status`/`command`/`diagnostics`/`exitCode`);
  All-OK-Fall ergibt `diagnostics: []` gemäß Lastenheft-Beispiel
  §1846-1852. `SeverityOK` und `SeverityInfo` werden gefiltert
  (Spec §1834 erlaubt `level` nur `warn|error`). `--quiet --json`
  ist semantisch identisch zu reinem `--json`; `--strict --json`
  upgraded Warn auf Exit-Code 11 ohne `status`-Drift (Spec §1837
  koppelt `status` an höchsten `level`, nicht an `--strict`).
  Broken-Pipe-Resistenz: fachlicher Exit-Code 11 hat Vorrang vor
  Write-Fehlern. **Root-PersistentFlag `--json`** für alle 10
  Spec-Enum-Subcommands; nicht-migrierte Forms rejecten mit
  `ErrJSONNotImplemented` (Exit-Code 2, [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes)-Klasse)
  plus Folge-Slice-Verweis. Zentrale Allowlist in
  `cli/root.go` via `PersistentPreRunE`; `--help` als
  read-only Escape-Hatch durchgelassen. **`u-boot template list
  --json`** wandert vom lokalen Flag aufs Root-Flag (beide
  Schreibweisen `template list --json` und `--json template list`
  identisches Output); Envelope-Migration folgt mit
  [`slice-v1-cli-json-dry-run-template`](docs/plan/planning/done/slice-v1-cli-json-dry-run-template.md), dokumentiert im
  `carveouts.md`-Temporär-Eintrag. **Common-Envelope
  `cliJSONEnvelope`** (`internal/adapter/driving/cli/jsonenvelope.go`)
  trägt Minimalkontrakt- und Voll-Schema-Felder über zwei
  Konstruktoren (`newMinimalEnvelope`/`newFullEnvelope`); Voll-
  Schema-Felder via Pointer-Wrapping (`*bool`/`*[]T`) als
  Anti-Drift gegen `omitempty`-Semantik-Refactor. **Schema-
  Helper-Sub-Package** `internal/adapter/driving/cli/jsontestutil/`
  mit zwei Modi `AssertMinimalEnvelope`/`AssertFullEnvelope`
  (Options-Pattern, kein neuer Dep) plus `DefaultAllowedCodes`-
  Registry für die 13 Doctor-Check-Codes. **Drei aktive Drift-
  Gates** schützen die Code-Registry: (1) `application.DoctorCheckIDs()`
  ↔ Map-Vollständigkeit, (2) Map ↔ Markdown-Roundtrip-Parser
  auf `docs/user/cli-json-output.md` §5.1 mit HTML-Marker-
  Sektion-Begrenzung, (3) Helper-Reject im Acceptance-Pfad
  für undokumentierte Codes. **Schema-Vertrag-Doku**
  ([`docs/user/cli-json-output.md`](docs/user/cli-json-output.md))
  zitiert Minimalkontrakt und Voll-Schema verbatim, dokumentiert
  Code-Registry und Per-Command-Migrations-Reihenfolge.
- `feat(logs): u-boot logs [service] [--follow] [--tail <n>]
  ([`LH-FA-UP-005`](spec/lastenheft.md#lh-fa-up-005-logs-anzeigen)) — neuer Subcommand streamt Compose-Logs als
  V1-Erweiterung der `up`/`down`-Familie. Ohne Service-Argument
  läuft `docker compose logs` über alle Services aus
  `compose.yaml` (T0-(a) Compose-Facade-Semantik, kein
  u-boot.yaml-Filter); mit Service-Argument nur diesen einen
  (Format-Validation via `domain.NewServiceName`, Existenz-Check
  delegiert an Compose). `--follow` blockt bis Ctrl-C und
  beendet via SIGINT-Vertrag-Schicht-1+2+3 (Adapter gibt
  `ctx.Err()` unverdeckt zurück, Use-Case übersetzt zu
  `(LogsResponse{}, nil)`, CLI exit-code 0). `--tail <n>` mit
  Stage-1-Validation auf nicht-negative Ganzzahlen (Default
  leer → Use-Case-Normalisierung zu Compose-`"all"`). Exit-
  Code-Mapping analog `up`/`down`: 10 (User/Project-State),
  11 (`ErrDockerUnavailable`), 12 (`ErrComposeRuntime`),
  14 (FS), 2 (CLI-Usage). Docker-tag E2E-Tests gegen echten
  postgres-Stack pinnen `--tail`-Buffer-Content und
  `--follow`-SIGINT-Vertrag.
- `feat(devcontainer): Devcontainer-Features-Allowlist und Katalog
  ([`LH-FA-DEV-003`](spec/lastenheft.md#lh-fa-dev-003-devcontainer-features)) — 8 Built-in-Features (`git`, `docker-cli`,
  `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`)
  plus External-Source-Allowlist via
  `devcontainer.featureSources.allow`. CLI:
  `u-boot config set devcontainer.features.<name>.{enabled,source,version}`
  plus `--allow-external-feature-sources <url>[,<url>...]` auf
  den drei Spec-§714-717-Pfaden (`init --devcontainer`,
  `generate devcontainer`, `config set devcontainer.featureSources.allow`).
  Doctor-Check `devcontainer.features.allowlist` (Error bei
  Allowlist-Violation, Warn bei Orphan-Activation oder fehlendem
  `enabled:`-Key). User-Doku in
  [`docs/user/devcontainer-features.md`](docs/user/devcontainer-features.md).
- `feat(doctor): devcontainer.features.drift Check` — über-Spec
  Drift-Erkennung zwischen `u-boot.yaml`'s Features-Map und den
  Keys im gerenderten `.devcontainer/devcontainer.json`. Drei
  Warn-Cases: aktiviertes Feature fehlt im JSON (Case 1, inkl.
  Datei-fehlt-Disziplin), deaktiviertes Feature noch im JSON
  (Case 2a), JSON-Key ohne cfg-Pendant (Case 2b, Hand-Edit-Hint).
  Repair-Hint `u-boot generate devcontainer`. Doctor-Total
  steigt 12→13.

## [0.3.0] - 2026-06-01

Third release. Completes the V1 „Add-on Catalogue Expansion"
milestone (5/5): the catalogue now ships three integrated
service add-ons — Postgres (since MVP), Keycloak ([`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003-keycloak-hinzufügen)
/ [`LH-AK-003`](spec/lastenheft.md#lh-ak-003-keycloak-flow)) and OpenTelemetry ([`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004-opentelemetry-hinzufügen) / [`LH-AK-004`](spec/lastenheft.md#lh-ak-004-opentelemetry-flow)) —
plus the matching `u-boot remove <service>` mirror, the
[`LH-FA-ADD-006`](spec/lastenheft.md#lh-fa-add-006-add-on-abhängigkeiten) `--with-deps` dependency-resolution mechanism, and
a doku-only audit closure for three V1 spec-IDs. Architectural
side-effect: the per-service catalogue pattern grew from a flat
`(compose, env, volume)`-tuple in M5 to a declarative entry with
`requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional` /
`healthcheckOptional` / `extraFiles` — any new add-on plugs in
by adding one catalogue row and three templates.

### Verified

- **Three V1 spec-IDs audit-closed** by
  [`slice-v1-audit-done`](docs/plan/planning/done/slice-v1-audit-done.md)
  — Doku-only verification that the existing code/doc state
  already satisfies the requirements:
  [`LH-FA-BUILD-006`](spec/lastenheft.md#lh-fa-build-006-aggregator-targets) (Aggregator-Targets `gates`/`ci`/`fullbuild` in
  the Makefile),
  [`LH-NFA-MAINT-004`](spec/lastenheft.md#lh-nfa-maint-004-dokumentierte-schnittstellen) (Add-on and template interfaces documented via
  [ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md)/-0009 + driving/driven port doc-comments + slice docs),
  [`LH-NFA-PORT-003`](spec/lastenheft.md#lh-nfa-port-003-containerfreundlichkeit) (u-boot itself runs in container / devcontainer:
  GHCR distroless image + container-aware `doctor` since v0.2.0 +
  six-platform binary distribution + `init --devcontainer`-generated
  files).

### Added

- **`u-boot add otel`** — [`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004-opentelemetry-hinzufügen) / [`LH-AK-004`](spec/lastenheft.md#lh-ak-004-opentelemetry-flow). Third and
  final add-on of the v0.3.0 milestone catalogue. Compose-Service
  mit Image-Pin `otel/opentelemetry-collector:0.108.0` (Stable),
  Port-Mappings `4317:4317` (OTLP/gRPC) + `4318:4318` (OTLP/HTTP),
  `command: --config=/etc/otel-collector-config.yaml`, Bind-Mount
  der gerenderten Config-Datei. Kein Healthcheck im Mindest-
  Setup ([`LH-AK-004`](spec/lastenheft.md#lh-ak-004-opentelemetry-flow) §2374 toleriert `running` ODER `healthy`).
  Mindest-Collector-Config in `otel-collector-config.yaml`:
  Receivers `otlp/grpc+http`, Processors `batch`, Exporters
  `debug` (stdout), Pipelines `logs`/`metrics`/`traces` — alle
  drei Signal-Typen aus [`LH-FA-ADD-004`](spec/lastenheft.md#lh-fa-add-004-opentelemetry-hinzufügen) §880.
  
  Internal: Catalogue-Pattern wächst um drei Felder pro Service —
  `extraFiles []extraFileEntry` für whole-file artefacts abseits
  von compose+env+volume (für OTel die Collector-Config-Datei),
  plus `envOptional` (implizit via leerem `envTmpl`) und
  `healthcheckOptional` für Services, die das Standard-Pattern
  legitim nicht brauchen. `executeAdd` schreibt extraFiles als
  vierten Slot nach yaml/compose/env; `executeRemove` löscht sie
  symmetrisch. `serviceComplete` skipt healthcheck-presence für
  `healthcheckOptional`; explicit `healthcheck.disable: true`
  bleibt hart abgelehnt. Acceptance-Helper-Reuse aus
  [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md) T3 (`acceptance_helpers.go`) — OTel-E2E
  bleibt ~30 Zeilen. **Makefile-Patch**: `test-docker`-Target
  mountet jetzt `/tmp` host-shared, damit Compose-Bind-Mount-
  Pfade vom Daemon (Host) aufgelöst werden können — sonst sieht
  der Daemon nur den Container-Pfad `t.TempDir()` nicht und
  erstellt einen leeren Verzeichnis-Mount, der den Collector
  beim Config-Read crasht. See
  [`slice-v1-otel`](docs/plan/planning/done/slice-v1-otel.md).
- **`u-boot add keycloak`** — [`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003-keycloak-hinzufügen) / [`LH-AK-003`](spec/lastenheft.md#lh-ak-003-keycloak-flow). Second
  add-on in the catalogue after Postgres. Compose-Service mit
  Image-Pin `quay.io/keycloak/keycloak:26.0` (LTS), Port-Mapping
  `8080:8080`, `command: start-dev` für [`LH-AK-003`](spec/lastenheft.md#lh-ak-003-keycloak-flow)-Boot, Healthcheck
  via `/dev/tcp/localhost/9000` (bash-builtin, kein curl im
  Image) gegen `/health/ready`. Admin-Credentials via Placeholder-
  Env-Block (`KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` +
  `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`).
  **Persistenz: flüchtige H2-In-Container-Datenbank** — kein
  Volume, nach `docker compose down` weg; [`LH-AK-003`](spec/lastenheft.md#lh-ak-003-keycloak-flow) verlangt nur
  Endpoint-200/302. Persistente externe Postgres-Anbindung
  ([`LH-FA-ADD-003`](spec/lastenheft.md#lh-fa-add-003-keycloak-hinzufügen) §857) bleibt als eigener Folge-Slice
  (`slice-v1-keycloak-external-postgres`, Trigger: Nutzer-Bedarf).
  Internal refactor: `renderPostgresTemplates` → generischer
  `renderServiceTemplates(svc)` über neue Service-Catalogue-
  Tabelle; `hasRequiredEnvKeys` / `contentScanState` /
  `inspectVolumeArtefact` / `patchTargetsFor` werden per-Service
  über `requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional`
  parametrisiert, damit Keycloak's volume-loser Pfad nicht in
  den Postgres-Repair-Loop läuft. Test-Helper-Extraktion:
  `internal/e2e/acceptance_helpers.go` teilt die init+add+up-
  Pipeline mit dem [`LH-AK-002`](spec/lastenheft.md#lh-ak-002-postgresql-flow)-Postgres-Test (Boot-Zeit-Carveout
  für Keycloak: 4 min UpService-Timeout vs. 90 s Postgres). See
  [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md).
- **`u-boot add <service> --with-deps`** — [`LH-FA-ADD-006`](spec/lastenheft.md#lh-fa-add-006-add-on-abhängigkeiten) add-on
  dependency mechanism. New domain type `AddOnDependency` (path-
  conditional service dependency declaration) + per-service
  catalogue side-table `dependenciesFor(svc)` (Postgres has none
  today; first non-nil row lands with [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md)). When
  the requested add-on declares a dep that is not yet registered
  in `u-boot.yaml`, the four-mode dispatch decides what happens:
  `--with-deps` auto-installs the chain (recursive `Add` calls,
  flag inherited so transitive deps follow); `--yes` has the same
  effect; `--no-interactive` (without `--yes`/`--with-deps`)
  fails fast with the new `ErrDependenciesRequired` sentinel
  (exit 10); default-interactive prompts via the new
  `Confirmer.ConfirmAddDependency(ctx, svc, missing)` driven-port
  method (mirror of `ConfirmRemoveVolumes` from M6). Postgres-
  only flows are unchanged — the no-deps short-circuit keeps the
  load+resolve cost out of the MVP catalogue path. Breaking
  refactor in the application layer: `NewAddServiceService`
  now takes a `Confirmer` between `yaml` and `logger`; all eight
  callsites updated in lock-step. See
  [`slice-v1-addons-deps`](docs/plan/planning/done/slice-v1-addons-deps.md).
- **`u-boot remove <service> [--purge]`** — first slice of the
  v0.3.0 milestone ("Add-on Catalogue Expansion"). Mirror of
  `u-boot add`: detects the [`LH-FA-ADD-005`](spec/lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern) service state, strips
  the `service.<name>` managed block from `compose.yaml` and
  `.env.example`, then sets `services.<name>.enabled: false` in
  `u-boot.yaml`. Idempotent: removing an already-disabled service
  is a no-op with a clear message. Inconsistent project state
  (orphan block, missing entry) surfaces as `ErrServiceInconsistent`
  with a manual-cleanup hint. New driving sentinel
  `ErrServiceUnregistered` (exit 10) distinguishes "service was
  never added" from "service name not in the catalogue"
  (`ErrServiceUnsupported`). [`LH-FA-ADD-007`](spec/lastenheft.md#lh-fa-add-007-service-entfernen) §"Volumes nur auf
  explizite Anforderung": `--purge` opts in destructively and
  triggers the [`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung) §254 confirmation gate (mirror of
  `u-boot down --volumes`); auto-removal of volumes is deferred
  to a follow-up slice — v0.3.0's `--purge` summary points at
  `docker volume rm <name>` for the manual cleanup. Internal:
  `detectServiceState` extracted from the M5 add path to a
  package-level function so both add and remove share it without
  duplication. See
  [`slice-v1-add-remove`](docs/plan/planning/done/slice-v1-add-remove.md).

## [0.2.0] - 2026-06-01

Second release. Adds the first two V1 template features
(`template list` + `init --template`), a cross-platform binary
distribution (six platforms as GitHub-Release assets), and a
container-aware `doctor` that no longer mis-reports a healthy host
as 4 errors when run from inside the distroless image. v0.1.1 was
originally planned as a patch-only tag for the doctor fix but is
skipped in favour of this minor bump — three features landed before
the tag-push and strict SemVer wants a MINOR bump for them.

### Added

- **`u-boot template list [--json]`** — first V1 template
  subcommand ([`LH-FA-TPL-004`](spec/lastenheft.md#lh-fa-tpl-004-templates-auflisten)). Enumerates the built-in project-
  template catalog with name, description, and version in a
  tabwriter-aligned table; `--json` emits a structured array
  with the full [`LH-FA-TPL-002`](spec/lastenheft.md#lh-fa-tpl-002-template-metadaten) metadata surface (`supportedAddOns`,
  `generatedFiles`, `requiredTools`, `variables`). Bootstrap
  built-in: `basic` (one template; further built-ins follow on
  demand per [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Folgepunkte 4). Fully hexagonal:
  `domain.TemplateMetadata` + `Validate()` (kebab-case-name
  regex, `ErrInvalidTemplate` sentinel), driven port
  `TemplateCatalog`, embed.FS-backed `externaltemplates` adapter,
  application `TemplateListService` (multi-`%w` so the original
  `domain.ErrInvalidTemplate` chain survives), CLI
  `template list` rendering. Adapter directory consolidated to
  `internal/adapter/driven/externaltemplates/` (no hyphen) for
  consistency with the existing `driven/`-adapter naming; [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md)
  §Entscheidung updated to match. See
  [`slice-v1-template-list`](docs/plan/planning/done/slice-v1-template-list.md).
- **`u-boot init <name> --template <name>`** — second V1 template
  feature, the render path of [`LH-FA-TPL-001`](spec/lastenheft.md#lh-fa-tpl-001-projektvorlagen) / [`LH-FA-TPL-002`](spec/lastenheft.md#lh-fa-tpl-002-template-metadaten). The
  init service delegates file rendering to the new
  `TemplateInitService` when `--template` is set; project structure
  directories and `git init` stay with the InitProjectService so
  the user-observable flow is one command. Byte-identity
  guarantee: `u-boot init demo --template basic` produces a
  project byte-identical to `u-boot init demo` for the six default
  files (`u-boot.yaml`, `compose.yaml`, `README.md`,
  `CHANGELOG.md`, `.env.example`, `.gitignore`) — pinned by an
  E2E `diff -r` test against the production catalog. Render engine:
  Go `text/template` for `*.tmpl` files, 1:1 copy for non-`.tmpl`
  files (per [ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md) §Entscheidung); `template.yaml` metadata is
  skipped. Two-phase render-then-write: a render error in any file
  short-circuits before the first disk write, so a buggy template
  no longer leaves a half-populated project. New
  `domain.TemplatePath` validator rejects `..` segments, absolute
  paths, Windows drive letters, backslashes, NUL bytes, and empty
  strings ([`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes) exit 10 via `ErrInvalidTemplatePath`).
  Mutex with `--devcontainer`/`--force`/`--backup`: surfaces as
  `ErrTemplateConflictsWithFlag` (exit 2) — v1 is fresh-init-only.
  Soft-existing-detection is skipped on the template path because
  `--template` resolves the "is this an existing project?"
  ambiguity by intent; the hard-existing check
  (`u-boot.yaml` present → `ErrProjectExists`) remains the
  safety net. Variable resolution + `--var key=value` deferred to
  a future slice (basic has no variables). See
  [`slice-v1-template-init`](docs/plan/planning/done/slice-v1-template-init.md).
- **Cross-platform binary distribution** for six platforms
  (Linux/macOS/Windows × amd64/arm64). `make build-binaries`
  cross-compiles every supported `GOOS`/`GOARCH` combination via
  the pinned `golang:$(GO_VERSION)` builder image (CGO disabled,
  `-ldflags "-s -w -X main.version=$(VERSION)"`, output to
  `bin/u-boot-<os>-<arch>[.exe]`). `.github/workflows/publish.yml`
  builds the same set after the GHCR push on every `v*` tag and
  attaches them as GitHub-Release assets via `gh release upload`.
  See
  [`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md)
  — [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md) §Folgepunkte 1 trigger pulled forward by the
  doctor-container-awareness feedback.
- Quickstart in `README.md` / `README.de.md` gets a host-native
  install block (`curl -sSL … | chmod +x` for Linux/macOS,
  `Invoke-WebRequest` for Windows) as the primary recommended path;
  the GHCR `docker run …` block is demoted to "alternative for
  container/CI workflows".
- `internal/hexagon/port/driven.RuntimeEnvironment` port plus
  `internal/adapter/driven/runtime.FileEnv` adapter: best-effort
  container detection via `/.dockerenv` (Docker Engine / Desktop)
  and `/run/.containerenv` (Podman / CRI-O / buildah). Drives the
  doctor-container-awareness change below.

### Changed

- `u-boot doctor` now skips the four host-prerequisite checks
  (`git.installed`, `docker.installed`, `docker.reachable`,
  `docker.compose.installed`) when running inside a container, with
  a `SeverityInfo` diagnostic and a hint that points at
  [`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md)
  for the rationale. Effect: `docker run --rm
  ghcr.io/pt9912/u-boot:0.2.0 doctor` no longer mis-reports a
  healthy host as 4 errors; exit code on an otherwise-clean project
  goes from 11 to 0. This addresses real-world feedback from the
  first `v0.1.0` GHCR pull (2026-05-31) where the distroless image's
  lack of bundled `docker` / `git` binaries surfaced as false-
  positive errors.
- Host installations are unaffected — `runtime.FileEnv` returns
  `false` outside containers, so the existing
  [`LH-FA-DIAG-002`](spec/lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen)-classified errors / warnings remain.

### Notes

`releases/latest/download/u-boot-<os>-<arch>[.exe]` resolves to the
highest stable tag — since `v0.1.0` predates binary assets, the
`latest`-shortcut starts working with `v0.2.0` (or any later tag).

## [0.1.0] - 2026-05-31

First public release. Closes the MVP scope from
[`spec/lastenheft.md`](spec/lastenheft.md) MVP-priority IDs: all
`LH-AK-*`, `LH-FA-*` and `LH-SA-*` items are delivered (audit trail
in [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md)
§MVP-Bilanz).

### Added — Subcommands

- `u-boot init [name] [--devcontainer]` — generate project skeleton
  (`u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`,
  `.env.example`, `.gitignore`, `docker/`, `scripts/`, `docs/`)
  and run `git init` ([`LH-FA-INIT-001`](spec/lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](spec/lastenheft.md#lh-fa-init-007-git-repository-initialisierung), [`LH-AK-001`](spec/lastenheft.md#lh-ak-001-minimaler-init-flow), [`LH-AK-005`](spec/lastenheft.md#lh-ak-005-devcontainer-flow)).
  Mode flags `--yes` / `--no-interactive` / `--assume-existing`
  ([`LH-FA-CLI-005A`](spec/lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung)); re-init with `--force` / `--backup` for
  managed-block edits vs. full overwrite ([`LH-FA-INIT-005`](spec/lastenheft.md#lh-fa-init-005-überschreibschutz)).
- `u-boot doctor [--strict]` — 11 diagnostic checks against the
  local environment + project, severity-classified
  (ok / warn / error), repair-hint output, exit-code 11 on errors
  (or warns under `--strict`) ([`LH-FA-DIAG-001`](spec/lastenheft.md#lh-fa-diag-001-doctor-befehl)..[`LH-FA-DIAG-004`](spec/lastenheft.md#lh-fa-diag-004-reparaturhinweise)).
- `u-boot add <service>` — idempotent state-machine for service
  add-ons; today's catalogue: `postgres` only ([`LH-FA-ADD-001`](spec/lastenheft.md#lh-fa-add-001-add-on-befehl)/[`LH-FA-ADD-002`](spec/lastenheft.md#lh-fa-add-002-postgresql-hinzufügen)/[`LH-FA-ADD-005`](spec/lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern),
  [`LH-AK-002`](spec/lastenheft.md#lh-ak-002-postgresql-flow), [`LH-AK-006`](spec/lastenheft.md#lh-ak-006-idempotenz)). Keycloak ([`LH-AK-003`](spec/lastenheft.md#lh-ak-003-keycloak-flow)) and
  OpenTelemetry ([`LH-AK-004`](spec/lastenheft.md#lh-ak-004-opentelemetry-flow)) follow in V1.
- `u-boot up [--timeout <sec>]` and `u-boot down [--volumes]` —
  Compose wrapper with healthcheck polling and TCP port probes
  ([`LH-FA-UP-001`](spec/lastenheft.md#lh-fa-up-001-umgebung-starten)..[`LH-FA-UP-004`](spec/lastenheft.md#lh-fa-up-004-umgebung-stoppen)). `--timeout 0` is fire-and-forget.
- `u-boot generate <changelog|readme|env-example|devcontainer>` —
  idempotent block-replace via the `U-BOOT MANAGED BLOCK: init`
  marker; user content outside the managed region is preserved
  byte-identically. `changelog` carries the [`LH-AK-007`](spec/lastenheft.md#lh-ak-007-changelog-generator) pin
  (no destructive edits to existing entries). Exit codes
  `0` / `2` / `10` / `14` per [`LH-FA-CLI-006`](spec/lastenheft.md#lh-fa-cli-006-exit-codes) ([`LH-FA-GEN-001`](spec/lastenheft.md#lh-fa-gen-001-generate-befehl)..[`LH-FA-GEN-005`](spec/lastenheft.md#lh-fa-gen-005-idempotenz),
  [`LH-FA-DEV-001`](spec/lastenheft.md#lh-fa-dev-001-devcontainer-erzeugen)/[`LH-FA-DEV-004`](spec/lastenheft.md#lh-fa-dev-004-benutzerrechte)/[`LH-FA-DEV-005`](spec/lastenheft.md#lh-fa-dev-005-ports)).
- `u-boot config [get <path> | set <path> <value>]` — whitelist-
  scoped reads/writes with two-stage schema validation (struct
  unmarshal + per-path domain re-validation) before any
  `WriteFile`. `services.<svc>.enabled` is get-only; toggling
  happens through `add` / `remove` to keep the add-on state
  machine atomic ([`LH-FA-CONF-001`](spec/lastenheft.md#lh-fa-conf-001-projektkonfiguration)..[`LH-FA-CONF-005`](spec/lastenheft.md#lh-fa-conf-005-konfiguration-anzeigen-und-ändern)).

### Added — CI & release infrastructure

- GitHub Actions CI workflow `.github/workflows/ci.yml` with three
  PR-blocking jobs ([`LH-QA-003`](spec/lastenheft.md#lh-qa-003-ci-fähigkeit-github-actions)): `gates (lint + test +
  coverage-gate)`, `security-gates (govulncheck)`,
  `image-scan (trivy HIGH+CRITICAL)`. All actions SHA-pinned;
  Docker-only runner ([`LH-FA-BUILD-007`](spec/lastenheft.md#lh-fa-build-007-docker-only-workflow)); per-job minimal
  permissions.
- GitHub Actions release workflow
  `.github/workflows/publish.yml` triggered on `v*` tags. Strict
  SemVer-2.0 validation (rejects leading-zero numeric prereleases
  and build-metadata `+...` tags), GHCR image push to
  `ghcr.io/pt9912/u-boot:<version>` (plus `:latest` for stable
  tags), OCI label verification, and live `--version` smoke test
  against the tag-derived `VERSION`.
- Local outer/inner-loop parity: `make image-scan` reproduces the
  `image-scan` CI job using the same Trivy version
  (`TRIVY_VERSION ?= 0.70.0`) the action installs.
- Multi-stage distroless runtime image (`gcr.io/distroless/static-debian12:nonroot`)
  built via `make build`; CGO-disabled static binary; version
  injected at build time as `-X main.version=<UBOOT_VERSION>`
  and as the `org.opencontainers.image.version` label.

### Added — Architecture & documentation

- Hexagonal architecture ([`LH-FA-ARCH-001`](spec/lastenheft.md#lh-fa-arch-001-hexagonales-pattern)..[`LH-FA-ARCH-003`](spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement), [ADR-0002](docs/plan/adr/0002-hexagonale-architektur.md)):
  `internal/hexagon/{domain,application,port/{driving,driven}}`
  + `internal/adapter/{driving,driven}`. `depguard` enforces
  layer rules in CI.
- 10 ADRs cover language (Go), architecture (hexagonal), lint
  profile (SOLID-near), CI system, CLI framework (Cobra), revive
  custom rules, distribution path (GHCR), plugin system (static —
  no plugins), template format (YAML + Go `text/template`), and
  the HTTP adapter (not built; CLI-only).
- User-facing setup docs:
  [`docs/user/quality.md`](docs/user/quality.md) (quality-gates
  overview) and
  [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
  (one-time GitHub UI activation of required status checks).
- German `spec/lastenheft.md` (~3000 lines, 14 sections + 4 open
  points all decided) is the single source of truth; English
  `README.md` / German `README.de.md` are equivalent.

### Known limitations and deliberate carve-outs

- **Add-on catalogue is intentionally small:** only `postgres`
  ships in v0.1.0. Keycloak and OpenTelemetry are V1.
- **Templates implementation is V1.** Format is decided
  ([ADR-0009](docs/plan/adr/0009-template-format-yaml-files.md): YAML + `text/template`); the three implementation
  slices ([`slice-v1-template-list`](docs/plan/planning/done/slice-v1-template-list.md), [`slice-v1-template-init`](docs/plan/planning/done/slice-v1-template-init.md),
  [`slice-later-local-templates`](docs/plan/planning/done/slice-later-local-templates.md)) follow on demand.
- **JSON / machine-readable output is V1.** `--json` and
  `--dry-run` flags ([`LH-FA-CLI-007`](spec/lastenheft.md#lh-fa-cli-007-dry-run)/[`LH-FA-CLI-008`](spec/lastenheft.md#lh-fa-cli-008-diff-ausgabe), [`LH-NFA-USE-004`](spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe)) are
  not yet shipped; [ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md) (no HTTP adapter) explicitly relies
  on this V1 track landing.
- **Distribution is GHCR-only.** Binary, Homebrew, Debian/RPM
  paths are deferred with explicit trigger slices in [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md).
  `npm` / `pip` are rejected (ecosystem mismatch).
- **No plugin loader.** Add-on system stays statically compiled
  into u-boot ([ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md)). Four re-evaluation triggers documented
  in [ADR-0008](docs/plan/adr/0008-plugin-system-statisch.md) §Folgepunkte.
- **CLI-only.** No HTTP / daemon adapter ([ADR-0010](docs/plan/adr/0010-kein-http-driving-adapter.md)); programmatic
  consumers use subprocess + `--json` once V1 lands.
- **Inner-loop is Docker-only** ([`LH-FA-BUILD-007`](spec/lastenheft.md#lh-fa-build-007-docker-only-workflow)). GNU `make`
  remains the single non-Docker host dependency (permanent
  carve-out to [`LH-NFA-PORT-002`](spec/lastenheft.md#lh-nfa-port-002-keine-unnötigen-systemabhängigkeiten)).

### Setup — required one-time GitHub UI activation

Before merging external PRs against `main`, activate the three
required status checks in GitHub UI per
[`docs/user/branch-protection.md`](docs/user/branch-protection.md):
the exact match strings are the workflow `name:` fields
(`gates (lint + test + coverage-gate)`,
`security-gates (govulncheck)`,
`image-scan (trivy HIGH+CRITICAL)`), not the shorter `jobs.<key>`
identifiers.

[Unreleased]: https://github.com/pt9912/u-boot/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/pt9912/u-boot/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/pt9912/u-boot/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pt9912/u-boot/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pt9912/u-boot/releases/tag/v0.1.0
