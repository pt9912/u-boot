# Slice V1: `up --json` / `down --json` — read-only Compose-Status-Envelope

> **Status:** `open/`. Sechster Folge-Slice (6/9) des Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 6/9). **Read-only-Klasse**: weder `--dry-run`
> noch `--diff` (Cluster-Slice Z. 464-467), nur `--json` mit
> typisiertem Data-Carrier. **Bündelt up+down** in einem Slice
> (Cluster-T0-(e) Z. 369-372): beide read-only-JSON, gemeinsamer
> Compose-Status-Reader, Confirmer-Pattern bei `down --volumes`
> (Spec §1015) 1:1 erbbar aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> T0-(j).
>
> Erbt das `Data any`-Wire-Field aus
> [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
> T0-(c)/(d): explizit *"für `slice-v1-cli-json-dry-run-up-down`"*
> vorgesehen mit Vorbild
> `type upStatusData struct { Services []serviceStatus
> ` + "`json:\"services\"`" + ` }`. Erbt
> `driving.WarningEntry`-Type (inkl. `Subject`-Feld) aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> T2 (R7-MED-F2 + R9-MED-F2 + R12-LOW-F4): bewusst generisch für
> up/down *"recreate-Warnings"* (Multi-Service-WARN "container
> 'postgres' will be replaced") plus config-set (8/9) value-warnings.

## Auslöser

Cluster-Slice §T0-Outcomes (a) macht jeden read-only-Subcommand
für `--json` verbindlich (`LH-NFA-USE-004` §1813). `up` und `down`
sind nach `doctor`/`add`/`init`/`generate`/`remove` die nächsten
und einzig verbleibenden modifying-CLI-Subcommands mit
Compose-Side-Effects: sie schreiben **nichts** auf das lokale
Filesystem (nur `ReadFile(compose.yaml)`), aber sie ändern den
Docker-Daemon-State (Container starten/stoppen, optional Volumes
entfernen).

Spec-Bezug:

- `LH-FA-UP-001` (Umgebung starten, §955-§978) — `u-boot up` mit
  `--timeout`-Stabilisierung
- `LH-FA-UP-002` (Docker Compose verwenden, §980-§986)
- `LH-FA-UP-003` (Startstatus anzeigen, §988-§1000) — Service-
  Name, Containerstatus, Port, Healthcheck-Status als
  Mindestangaben → direkter Carrier-Field-Anker für
  `upStatusData`
- `LH-FA-UP-004` (Umgebung stoppen, §1003-§1019) — `u-boot down`
  mit `--volumes`-destructive-Opt-in
- `LH-FA-CLI-005A` §234/§246/§254 — Confirmation-Gate für
  `down --volumes` (gemeinsam mit `remove --purge`, `init
  --force`, `config set` destructive paths)
- `LH-NFA-USE-004` §1813 / §1841 — Minimalkontrakt-Pflicht
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag (NICHT für
  up/down weil keine `--dry-run`-Variante, aber das
  `Data any`-Field gehört zur `cliJSONEnvelope`-Struktur)

Heute-Stand-Pre-Scan
(`internal/adapter/driving/cli/up.go`, 110 LOC;
`internal/adapter/driving/cli/down.go`, 117 LOC;
`internal/hexagon/application/upservice.go`, 521 LOC;
`internal/hexagon/application/downservice.go`, 131 LOC):

| Aspekt | up | down |
| --- | --- | --- |
| Positional-Args | `cobra.NoArgs` (`up.go:63`) | `cobra.NoArgs` (`down.go:73`) |
| Lokale Flags | `--timeout` (default 60, 0=fire-and-forget) | `--volumes` (destructive opt-in) |
| Persistente Flags read-through | `--quiet` | `--yes`, `--no-interactive`, `--quiet` |
| FS-Mutation | KEINE | KEINE |
| FS-Read | `ReadFile(compose.yaml)` einmal in `readComposeFile` | `ReadFile(compose.yaml)` analog |
| Docker-Mutation | `engine.ComposeUp` + `engine.ComposePs`-Polling | `engine.ComposeDown` (+ optional `ComposeRm --volumes`) |
| WARN-Emission heute | `renderUpDiagnostics(stdout, resp.Result.Diagnostics, quiet)` — auf **stdout** (M6 §T6) | `renderDownSuccess` — minimal output, keine Diagnostics |
| Confirmer-Pfad | KEIN | `down --volumes` ruft Confirmer analog `remove --purge` (Spec §1015 + §254) |
| ProgressSink | `stderr` für Compose-pull/create/start/healthcheck (LH-NFA-PERF-002) | `stderr` für compose-down-phases (LH-NFA-PERF-002) |
| Exit-Codes | 10 (no-yaml), 11 (Docker unreachable), 12 (Compose-Runtime) | 10 (no-yaml, confirmation refused), 11, 12 |
| Allowlist heute | `jsonallowlist.go:74-75` Reject mit Follow-up `up-down` | analog |

Use-Case-Deps `UpService`: `driven.FileSystem`, `driven.YAMLCodec`,
`driven.DockerEngine`, `driven.NetProbe`, `driven.Clock`,
`driven.Logger`. Use-Case-Deps `DownService`: `driven.FileSystem`,
`driven.YAMLCodec`, `driven.DockerEngine`, `driven.Confirmer`,
`driven.Logger` (+ optional). **KEIN** RecordingFileSystem-Bedarf
(read-only FS), **KEIN** PreviewMode-Bedarf, **KEIN** Diff-
Renderer-Bedarf.

## Aufhebungsbedingung

Zwei Flag-Kombinationen pro Subcommand liefern spec-konforme
Outputs:

```bash
u-boot up                                # Human-Mode (heutiges Verhalten)
u-boot up --json                         # Minimal+Data-Envelope
u-boot up --json --timeout=0             # analog fire-and-forget
u-boot up --json --quiet                 # semantisch identisch zu --json (Cluster-T0-(a))

u-boot down                              # Human-Mode (heutiges Verhalten)
u-boot down --json                       # Minimal+Data-Envelope, KEIN Confirmer-Prompt
u-boot down --volumes --yes --json       # Confirmer geswapped, Volumes removed, Envelope
u-boot down --volumes --no-interactive --json   # ErrConfirmationRequired-Envelope, exit 10
u-boot down --volumes --json             # default interactive prompt — analog remove?
                                          # (T0-Sub-Decision: Prompt im JSON-Mode unterdrücken oder erlauben?)
```

`make gates` grün (lint + test + coverage-gate ≥ 90 % + docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ **`--json`-Allowlist-Migration**: `"u-boot up"` und
  `"u-boot down"` aus dem Reject-Block in `jsonallowlist.go:74-75`
  raus, beide Subcommands tragen Envelope.
- ✅ **Envelope-Shape**: `command="up"` bzw. `command="down"`,
  KEIN `subcommand`-Feld (beide sind Top-Level-Subcommands ohne
  Sub-Form). KEIN `dryRun`/`diff`/`plannedFiles`/`changes`/`hunks`-
  Feld (read-only-Klasse). Pflicht-Felder pro Spec §1841:
  `status`/`command`/`diagnostics`/`exitCode` plus typed `data`-
  Carrier.
- ✅ **`upStatusData`-Carrier-Form** (doctor T0-(c)/(d) Vorlauf):
  pro Service `{name, state, ports, healthcheck}` als
  `serviceStatus`-Sub-Struct mit `json:"…"`-Tags. **Pointer-
  Wrapping** auf optionalen Feldern (`Healthcheck *string` weil
  nicht alle Services healthchecks haben).
- ✅ **`downStatusData`-Carrier-Form** (T0-Sub-Decision):
  `{removedVolumes []string, stoppedContainers []string}` oder
  schmaler `{removedVolumes []string}` — abhängig davon, was die
  Use-Case-Response heute trägt. Bei `down` ohne `--volumes`
  ist `removedVolumes: []`.
- ✅ **Idempotenz-Pin**: `down` gegen bereits-gestoppte Umgebung
  liefert `removedVolumes: []`, `status: ok`, Exit 0 (analog
  remove NoOp-Semantik).
- ✅ **`--quiet --json` semantisch identisch zu `--json`**
  (Cluster-T0-(a) Pattern aus doctor T6-Pin): `--quiet` darf den
  JSON-Output NICHT unterdrücken — JSON ist die Maschinen-
  Schnittstelle, `--quiet` ist die Human-Schnittstelle-Unterdrückung
  und kollidiert semantisch.
- ✅ **ProgressSink-Silencing im JSON-Mode**: Compose-Phase-
  Streaming auf stderr (LH-NFA-PERF-002 "pull/create/start/
  healthcheck phases stream to stderr live") MUSS in `--json`
  unterdrückt werden, sonst polluten Live-Phasen-Logs den stderr
  für JSON-Konsumenten. Pattern-Erbe von init T0-(o)
  `ProgressPort`-Silencing — aber up/down nutzen `ProgressSink
  io.Writer` (kein Port), also via Request-Field
  `req.ProgressSink = io.Discard` oder Service-Field-Swap.
  T0-Sub-Decision: welche der zwei Formen.
- ✅ **`down --volumes`-Confirmer-Pattern-Erbe** (T0-(j) aus
  remove): `req.SilenceConfirmer = flags.JSON` triggert
  Service-Field-Mutation auf `noopConfirmer`-Instanz (analog
  remove `RemoveServiceService.Remove()`-Wrapper). Bei
  `--volumes --no-interactive --json` OHNE `--yes` →
  `ErrConfirmationRequired`-Envelope mit `LH-FA-INIT-005` /
  Exit 10 (geteilter Code mit init/remove).
- ✅ **WARN-Migration in `diagnostics[]`**: heutige
  `renderUpDiagnostics`-Calls auf stdout (M6 §T6) wandern im
  JSON-Mode in `diagnostics[]` mit `level: "warn"` und passenden
  LH-Codes. Multi-Service-WARNs (z. B. *"container 'postgres'
  will be replaced"*) nutzen das proaktiv eingeführte
  `Subject`-Feld auf `driving.WarningEntry` (R12-LOW-F4-Vorlauf
  aus remove T2).
- ✅ **Mapper-Tabelle** (analog `mapRemoveErrorToDiagnostic`):
  per-Subcommand-Mapper `mapUpErrorToDiagnostic` und
  `mapDownErrorToDiagnostic` ODER gebündelter
  `mapComposeErrorToDiagnostic(err, command string)` (T0-Sub-
  Decision). Sentinels: `ErrDockerUnavailable` →
  `LH-FA-CLI-006`?/Exit 11, `ErrComposeRuntime` →
  `LH-FA-CLI-006`?/Exit 12, `ErrConfirmationRequired` →
  `LH-FA-INIT-005`/Exit 10, `ErrConflictingModeFlags` →
  `LH-FA-CLI-005A`/Exit 2, `ErrInvalidTimeout` →
  `LH-FA-CLI-006`/Exit 2, `ErrProjectNotInitialized` →
  `LH-FA-UP-001`/Exit 10. T0-Sub-Decision: welcher LH-Code für
  Docker-/Compose-Runtime-Klasse — heute gibt's keinen
  `LH-NFA-REL`-Eintrag für Docker-Daemon-Verfügbarkeit.
- ✅ **Mid-Operation-Failure-UX** (analog Mid-Write aber für
  Docker): wenn `ComposeUp` mid-stream failt (Container A
  started, Container B failed), liefert der Envelope den
  Status-Snapshot bis zur Failure-Stelle plus `diagnostics[]`-
  Eintrag mit Failure-Position und Exit-Code 12.
- ✅ **CLI-Pin-Tests**: ~10-14 Acceptance-Tests in
  `up_acceptance_test.go` + `down_acceptance_test.go` (oder
  einer gebündelter Form `updown_acceptance_test.go`).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (up-down→done),
  neue §6.7-Sektion mit Pattern-Vorgabe aus doctor T0-(d), §7
  Mutations-Matrix-Zeilen ("`up`: nur ReadFile", "`down`: nur
  ReadFile" — beide read-only).
- ✅ **CHANGELOG `### Added`-Eintrag** analog
  remove/generate/init.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a) Bündelung up+down in einem Slice — wie viel Code-
  Sharing?** Cluster-T0-(e) Z. 369-372 sagt "denselben Compose-
  Status-Reader brauchen". Konkret: gibt's ein gemeinsames
  `composeStatusEnvelope`-Helper-Pattern, oder bleibt jeder
  Subcommand selbst-tragend mit kopierter Envelope-Logik?
  Sub-Decision-Optionen:
  (a) Beide Subcommands rufen einen geteilten
      `writeComposeStatusJSON(out, data, warnings)`-Helper,
      lebt in `cli/composestatus.go` (neu).
  (b) Jeder Subcommand hat eigenes `writeUpJSON` /
      `writeDownJSON`, kein geteilter Helper — Pattern wie
      remove `writeRemoveJSON`. Mehr Code-Duplikation, weniger
      Abstraktions-Last.
  Plan-Empfehlung: (b) — die Carrier-Types unterscheiden sich
  (Services-Array vs. RemovedVolumes-Array), das geteilte
  Pattern wäre nur `newDataEnvelope`-Call + Allowlist-Eintrag
  (beides existierende Helper). Helper-Extraction-Druck reift
  erst nach 8/9 (consolidation-Slice).
- **T0-(b) `--quiet --json` Pattern festzurren** (Cluster-T0-(a)
  doctor T6 Vorbild): `--quiet` UNterdrückt im Human-Mode
  status-table + diagnostics-Section auf stdout; im JSON-Mode
  MUSS der Envelope erscheinen. Code-Pfad-Pin: `if flags.JSON
  { … return writeUpJSON(...) } if flags.Quiet { return nil }
  … renderUpStatus(...)`. T6-Pin verifiziert beide Reihenfolgen
  `--quiet --json` und `--json --quiet`.
- **T0-(c) ProgressSink-Silencing-Form**: drei Optionen:
  (a) `req.ProgressSink = io.Discard` aus dem CLI bei
      `flags.JSON` — request-time Substitution, Use-Case sieht
      das Discard-Writer.
  (b) Service-Field-Swap analog init T0-(o) ProgressPort —
      braucht Service-Mutation und defer-Restore innerhalb der
      Lock-Region (gibt's heute keinen `upMu` / `downMu`).
  (c) CLI-Layer-Wrap: ein `discardOnJSONWriter`-Decorator,
      injiziert in den ProgressSink-Field.
  Plan-Empfehlung: (a) — `io.Discard` ist die schlankste Form,
  Use-Case-Signatur unverändert, kein Mutex-Erfordernis.
- **T0-(d) `down --volumes` Confirmer-Swap-Erbe**: 1:1 aus
  remove T0-(j). `req.SilenceConfirmer = flags.JSON` (oder
  request-time-swap, analog T0-(c) ProgressSink). Sub-
  Entscheidung: gleicher Form-Erbe wie T0-(c), Konsistenz-Pin.
- **T0-(e) Mapper-Tabelle Layer-Heim**: zwei separate Mapper
  `mapUpErrorToDiagnostic`/`mapDownErrorToDiagnostic` ODER ein
  gemeinsamer `mapComposeErrorToDiagnostic(err, command
  string)`. Heutige Sentinels überlappen stark
  (`ErrDockerUnavailable`, `ErrComposeRuntime`, `ErrProject
  NotInitialized` sind in beiden), aber `down`-spezifisch sind
  `ErrConfirmationRequired`/`ErrConflictingModeFlags`/
  `volumes`-Sentinels und `up`-spezifisch `ErrInvalidTimeout`.
  Plan-Empfehlung: separate Mapper analog
  `mapRemoveErrorToDiagnostic`/`mapAddErrorToDiagnostic`
  (Pattern-Erbe), aber geteilter Helper für die
  Docker-/Compose-Runtime-Sentinels als interner
  `mapComposeRuntimeSentinel(err)`-Helper.
- **T0-(f) LH-Code-Klassifikation für Docker-Runtime-Klasse**:
  `ErrDockerUnavailable` → Exit 11 ist gesetzt, aber welcher
  LH-Code? Spec hat `LH-NFA-REL-003` für FS-Failure und
  `LH-FA-CLI-006` als Default. Für Docker-Daemon-
  Unverfügbarkeit gibt's keinen dedizierten LH-Code.
  Sub-Decision: (i) neuer `LH-NFA-REL-005`-Code (Plan-
  Annahme) ODER (ii) Konsolidierung auf existierenden
  `LH-NFA-REL-003` mit Sub-Semantik-Dehnung (analog remove's
  Triple-Use von `LH-FA-ADD-005`). Plan-Vorschlag: (ii)
  Konsolidierung — der Exit-Code 11 differenziert bereits;
  Spec-Eintrag bleibt schmal.
- **T0-(g) `upStatusData`-Field-Granularität**: doctor T0-(c)
  zitiert das Vorbild als `Services []serviceStatus`, aber
  `serviceStatus`-Felder sind nicht festgenagelt. Spec
  LH-FA-UP-003 fordert: Name, Containerstatus, Port,
  Healthcheck (optional). Sub-Decision: (i) genau diese vier
  Felder ODER (ii) erweitert um zusätzliche Container-Metadaten
  (Image-Digest, Uptime, Restart-Count). Plan-Empfehlung: (i)
  Spec-konform-minimal, Erweiterung als eigener
  Sub-Slice falls Real-World-Druck.
- **T0-(h) `downStatusData`-Field-Definition**: was trägt der
  Carrier? Heute liefert `DownService.Down`
  `DownResponse{RemovedVolumes []string}` — das ist die
  einzige nicht-triviale Datenstruktur. Sub-Decision: nur
  `{removedVolumes []string}` ODER ergänzt um
  `{stoppedContainers []string, removedVolumes []string}` für
  Konsumenten-Klassifikation. Plan-Empfehlung: minimal nur
  `removedVolumes` (das ist das einzige destructive-Sub-Result),
  Container-Liste ist Compose-Implicit.
- **T0-(i) Mid-`ComposeUp`-Failure-Capture-Vertrag** (analog
  init/add Mid-Write-Failure-UX): wenn `ComposeUp` mid-stream
  failt (Container A started, B failed), trägt der Envelope
  `data.services[]` Snapshot bis zur Failure-Stelle plus
  `diagnostics[]` mit Failure-Position. Sub-Decision: trägt
  `serviceStatus` ein `state: "failed"`-Wert ODER ist Failure
  nur in `diagnostics[]` markiert (Service fehlt im `services[]`-
  Array)? Plan-Empfehlung: `state: "failed"` analog
  `domain.ContainerState`-Enum.
- **T0-(j) `--timeout=0` Fire-and-Forget im JSON-Mode**: heute
  bedeutet `--timeout=0` "no polling, no probes, status table
  omitted, info diagnostic shown" (`up.go:53-54`). Im JSON-Mode
  bleibt das info-diagnostic erhalten (`level: "info"` ist
  Spec §1834 verboten — siehe doctor-Pattern aus §97 doctor-
  Slice). Sub-Decision: `level: "warn"` UpRound oder Field-
  Drop. Plan-Empfehlung: Field-Drop weil "fire-and-forget" kein
  WARN-Niveau erreicht; ggf. eigener `data.timeoutFireAndForget:
  true`-Marker im Carrier.
- **T0-(k) Recreate-Warnings-Semantik** (R12-LOW-F4 aus remove
  T2 setzt den Type-Vorlauf): wann emittiert `up` eine WARN
  *"container 'postgres' will be replaced"*? Heute existiert
  in `upservice.go` keine recreate-Detection. Sub-Decision:
  (i) Recreate-Detection als V1-Scope ODER (ii) als Carveout
  für Folge-Slice (analog Volume-Auto-Removal aus remove).
  Plan-Empfehlung: (ii) — Recreate-Detection braucht
  Compose-Plan-Pre-Walk (`docker compose config`-Parse +
  Container-Hash-Vergleich), nicht-trivial. WARN-Carrier-Type
  ist proaktiv da, aber konkrete Detection kommt später.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(k) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog remove T1): `noopConfirmer` lebt bereits in `application/noop.go:17-33`, `io.Discard` ist Go-stdlib — beide Helper für ProgressSink-Silencing und Confirmer-Swap existieren | — (entfällt) | T0 |
| T2 | Port-Types: `UpRequest.SilenceProgress`/`SilenceDiagnostics` (oder request-time ProgressSink-Discard); `DownRequest.SilenceConfirmer`-Feld analog `RemoveServiceRequest`. `UpResponse.Warnings []driving.WarningEntry` (Type schon da aus remove T2). Sentinels: ggf. neuer `LH-FA-UP-001`-Sub-Sentinel für ProjectNotInitialized falls heute Roh-Error | ~80 | T0 |
| T3 | Application-Layer: ProgressSink-Discard-Wiring im JSON-Mode (oder Service-Field-Swap), `DownService.Down()`-Wrapper mit Confirmer-Swap (analog `RemoveServiceService.Remove()`-Wrapper) | ~120 | T2 |
| T4 | Composition-Root: heute existiert KEINE up/down-fsFactory (kein Recorder) — T4 prüft ob `cmd/uboot/main.go` Wiring-Updates braucht (vermutlich nur Confirmer-Bezug für down) | ~20 | T3 |
| T5 | CLI-RunE: `runUp`/`runDown` ruft generische Helper mit `command="up"`/`"down"`, `mapErr=mapUpErrorToDiagnostic`/`mapDownErrorToDiagnostic`; Envelope-Pfade; Allowlist-Migration; Mapper neu; `data`-Structs (`upStatusData`/`downStatusData`); WARN-Migration | ~200 | T2 |
| T6 | Acceptance-Tests: ~10-14 Tests (Envelope-Pin both Subcommands, Idempotenz-Pin für down, `--quiet --json`-Pin, ProgressSink-Silencing-Pin, Confirmer-Swap-Pin für down --volumes, ConflictingModeFlags-Pin, Service-Sentinels-Pins) | ~400-500 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe) | ~50 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.7/§7 (zwei §7-Zeilen "nur ReadFile"), roadmap done-Zähler 5→6, ggf. carveouts.md (Recreate-Warnings-Carveout aus T0-(k)), Slice nach `done/` mit DoD-Hash-Tabelle | — (Doku) | T7 |

LOC-Bilanz vorläufig: ~800-1000 (Schätzung nach unten korrigiert
gegenüber remove ~1200-1400 wegen Read-Only-Klasse). Pattern-Erbe
von remove (Confirmer-Swap), init (ProgressPort-Silencing-Vorbild)
und doctor (typed Data-Carrier).

## Out of Scope

- **Recreate-Detection** (T0-(k) Sub-Decision): Compose-Plan-Pre-
  Walk + Container-Hash-Vergleich für *"container 'postgres' will
  be replaced"*-WARN ist V1-Out-of-Scope. WARN-Carrier-Type ist
  via `driving.WarningEntry`-Vorlauf (remove T2 R12-LOW-F4)
  proaktiv vorhanden, konkrete Detection wandert in Folge-Slice
  (Trigger: User-Feedback über fehlende Replace-Warnings oder
  Cluster-T_close-Audit).
- **`up --service <name>`-Selective-Form**: heute liefert `up`
  alle Compose-Services (`cobra.NoArgs`); ein zukünftiger
  Sub-Form für Single-Service-Start wäre eigener Slice mit Args-
  Validator-Pattern (würde dann auch das envelope-consolidation-
  Pattern erben).
- **Pattern-Drift-Konsolidierung** mit add/init/generate
  (R15-Cross-Slice-1 aus remove): up/down landen VOR der
  Konsolidierung
  ([`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md))
  und übernehmen den heutigen Drift-Status mit (kein Args-
  Validator weil `cobra.NoArgs`, kein BaseDirSanitizedError-
  Bedarf für Compose-Runtime-Errors).
- **Docker-Daemon-Version-Reporting im Envelope** (z. B.
  `data.dockerVersion`): nicht in Spec gefordert, eigener Slice
  falls Konsumenten-Bedarf.

## Bezug

- Cluster:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  (Folge-Slice 6/9).
- Pattern-Vorbilder:
  [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  (Data-Carrier `upStatusData`-Vorbild + Read-Only-Envelope-Form),
  [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
  (ProgressPort-Silencing-Vorbild im JSON-Mode),
  [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
  (Confirmer-Swap-Pattern T0-(j) + `driving.WarningEntry`-Type).
- Code-Anker:
  [`cli/up.go`](../../../../internal/adapter/driving/cli/up.go),
  [`cli/down.go`](../../../../internal/adapter/driving/cli/down.go),
  [`application/upservice.go`](../../../../internal/hexagon/application/upservice.go),
  [`application/downservice.go`](../../../../internal/hexagon/application/downservice.go),
  [`cli/jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  Z. 29/74-75.
- Folge-Slices:
  [`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md)
  (R15-Cross-Slice-1 — retroaktive Helper-Extraction, Trigger
  Cluster-Stand 8/9).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
