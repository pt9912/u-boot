# Slice V1: `logs --json` — Streaming-Output mit Modell-Entscheidung

> **Status:** `open/`. Siebter Folge-Slice (7/9) des Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 7/9). **Read-only-Klasse** mit Streaming-
> Vertrag: weder `--dry-run` noch `--diff` (Cluster-Slice
> Z. 464-467); zentrale Sub-Decision ist das **Output-Modell**
> für `--json` (JSON-Lines vs. Single-Envelope) — Cluster-Plan
> Z. 326-329 hat diese explizit hierher ausgelagert.
>
> Erbt Read-only-Klassen-Disziplin aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> (kein PreviewMode, kein RecordingFileSystem, kein
> `--dry-run`/`--diff`). Erbt `cli/sanitize.go`-Helper
> (Pfad-Leak-Defense) und `cli/composesentinel.go`-Mapper-
> Helper aus up-down T5. Erbt `driving.WarningEntry`-Type aus
> remove T2 falls Logs WARN emittiert (heute keine
> bekannten Pfade).

## Auslöser

Cluster-Slice §T0-Outcomes (a) macht `--json` für jeden
Subcommand verbindlich (`LH-NFA-USE-004` §1813). `u-boot logs`
ist nach `up`/`down` der nächste Read-only-Subcommand und der
**erste Streaming-Subcommand** — alle bisherigen Folge-Slices
(doctor/add/init/generate/remove/up-down) liefern Single-
Envelopes nach abgeschlossener Operation. Logs ist
strukturell anders: ohne `--follow` Bounded-Output mit
`--tail <n>`; mit `--follow` Unbounded-Stream bis SIGINT.

Spec-Bezug:

- `LH-FA-UP-005` (Logs anzeigen) — Streaming-Vertrag + Tail-
  Semantik.
- `LH-NFA-USE-004` §1813 / §1841 — Minimalkontrakt-Pflicht.
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag (NICHT für
  logs weil keine `--dry-run`-Variante; Plan-Anker bleibt für
  `cliJSONEnvelope`-Struktur).

Heute-Stand-Pre-Scan
([`internal/adapter/driving/cli/logs.go`](../../../../internal/adapter/driving/cli/logs.go),
181 LOC;
[`internal/hexagon/application/logsservice.go`](../../../../internal/hexagon/application/logsservice.go),
164 LOC;
[`internal/hexagon/port/driving/logs.go`](../../../../internal/hexagon/port/driving/logs.go),
95 LOC):

| Aspekt | Heute |
| --- | --- |
| Positional-Args | `cobra.MaximumNArgs(1)` (`logs.go:87`) — optionaler Service-Name |
| Lokale Flags | `--follow` (default false), `--tail <n>` (default "" = Compose-Default "all") |
| Persistent Flags read-through | `--quiet` (heute nicht im `runLogs`-Signaturpfad — eigener Read-Check nötig) |
| FS-Mutation | KEINE |
| FS-Read | `u-boot.yaml` + `compose.yaml` Pre-Checks (analog up/down `checkProjectInitialized` + `checkComposeFilePresent`) |
| Docker-Operation | `engine.ComposeLogs(ctx, baseDir, ...)` — Streaming-Adapter |
| Output | direkter Compose-Stream via `req.OutputSink` (`cmd.OutOrStdout()`) — Service-Prefix + Lines |
| Sentinels | `ErrInvalidLogsTail` (Exit 2), `ErrProjectNotInitialized` (10), `ErrComposeFileMissing` (10), `domain.ErrInvalidServiceName` (10), `driven.ErrDockerUnavailable` (11), `driven.ErrComposeRuntime` (12) |
| SIGINT | `(LogsResponse{}, nil)` — Application-Layer fängt `context.Canceled` ab → Exit 0 |
| Allowlist heute | `jsonallowlist.go:74-75` Reject mit Follow-up `logs` |

Use-Case-Deps `LogsService`: `driven.FileSystem` (read-only),
`driven.DockerEngine` (ComposeLogs streaming). KEIN
`driven.YAMLCodec`, KEIN `Confirmer`, KEIN `Logger`-bound-state.

## Aufhebungsbedingung

Vier Flag-Kombinationen für `u-boot logs [service]` liefern
spec-konforme Outputs:

```bash
u-boot logs                                # Human-Mode (heutiges Verhalten)
u-boot logs --json                         # T0-(a) Sub-Decision: Single-Envelope ODER JSON-Lines
u-boot logs postgres --json --tail=100     # bounded, identische Form-Wahl wie obig
u-boot logs --follow                       # Human-Mode bis SIGINT
u-boot logs --follow --json                # T0-(a) Sub-Decision: ist Stream-Output spec-konform?
```

`make gates` grün (lint + test + coverage-gate ≥ 90 % +
docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ **Output-Modell festgezurrt** (T0-(a) Sub-Decision — siehe
  unten): eine der drei Optionen mit T0-Review-Begründung
  gewählt; Pattern-Erbe-Disziplin gegen Cluster-Slice §§1841
  belegt.
- ✅ **`--json`-Allowlist-Migration**: `"u-boot logs": true` in
  `jsonAllowlist()`; Reject-Liste schrumpft auf 4 (config
  bare/get/set, template bare).
- ✅ **Envelope-Shape** (Single-Envelope-Pfad, falls T0-(a)
  Option A): `command="logs"`, kein `subcommand`-Feld; KEIN
  `dryRun`/`diff`/`plannedFiles`/`changes`/`hunks`. Pflicht-
  Felder pro Spec §1841: `status`/`command`/`diagnostics`/
  `exitCode` plus typed `data`-Carrier `logsData{lines
  []string}` ODER strukturiert `[]logLine{service, line}`
  (T0-(c) Sub-Decision).
- ✅ **JSON-Lines-Pfad** (falls T0-(a) Option B): pro Compose-
  Log-Zeile ein NDJSON-Object `{"level": "info", "code":
  "LH-FA-UP-005", "message": "<line>", "service":
  "<prefix>"}` auf stdout. Letzte Zeile ein Final-Envelope
  mit `status`/`exitCode`. **Spec-§1841-Vertrag-Bruch**:
  Konsument bekommt nicht EINEN Envelope sondern N — als
  Streaming-Sub-Pattern dokumentiert in `cli-json-output.md`
  §6.8 mit explizitem Carveout-Vermerk.
- ✅ **`--follow --json` Semantik**: ist Stream-Output
  überhaupt spec-konform? T0-Review prüft drei Alternativen
  (a) Reject mit Exit 2 (b) Stream-Pattern (NDJSON) (c)
  Buffer + Single-Envelope nach SIGINT.
- ✅ **`--quiet --json` semantisch identisch zu `--json`**
  (Cluster-T0-(a) doctor-Pattern). `--quiet` heute im
  `runLogs`-Pfad nicht expliziert — T5 muss das ändern
  (analog up/down).
- ✅ **`--tail` im JSON-Mode**: `--tail=10 --json` bounded =>
  data-Form je T0-(a) Wahl.
- ✅ **SIGINT-Vertrag im JSON-Mode**: heute SIGINT → Exit 0 +
  nichts. Im JSON-Mode (a) Single-Envelope mit `status: ok`
  vor SIGINT-Handler (b) NDJSON-Final-Object mit `status:
  ok` (c) kein Output, nur Exit-Code.
- ✅ **Mapper-Tabelle** (`mapLogsErrorToDiagnostic`) analog up/
  down-Pattern mit Switch-Order FS-first → ggf. neuer
  `ErrLogsFileSystem`-Sentinel, dann `mapComposeRuntime
  Sentinel`-Helper, dann Logs-spezifische Sentinels
  (ErrInvalidLogsTail), dann cross-cutting fachlich
  (ErrComposeFileMissing/ErrProjectNotInitialized),
  Default.
- ✅ **Path-Leak-Defense**: `runLogs` wrappt UC-Errors mit
  `sanitizeBaseDir(err, cwd)` vor `reportError` analog up/
  down-T5.
- ✅ **`baseDirSanitizedError`-Wiederverwendung**: aus
  `cli/sanitize.go` (etabliert in up-down T5).
- ✅ **CLI-Pin-Tests**: ~10-14 Acceptance-Tests in
  `logs_acceptance_test.go` (Envelope-Pin Single-/Stream-
  Form je T0-(a) Wahl, --quiet --json, --follow --json
  Sub-Decision-Pin, --tail-Bounded-Pin, Service-Sentinels
  Mapper-Rows, SIGINT-Vertrag-Pin, Path-Leak-Sanitizer-Pin).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (logs→done),
  neue §6.8-Sektion mit Output-Modell-Form + ggf. Streaming-
  Carveout-Vermerk, §7 Mutations-Matrix-Zeile (logs: nur
  ReadFile).
- ✅ **CHANGELOG `### Added`**-Eintrag analog up-down.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a) Output-Modell festzurren** — **zentrale Sub-Decision**
  (Cluster-Plan Z. 326-329 hierher ausgelagert): drei Optionen:

  | # | Form | Vor | Contra |
  | - | --- | --- | --- |
  | A | **Single-Envelope nach Stream-Ende** (Pattern-konsistent mit Cluster) | Spec-§1841-konform; eine Envelope-Form für ALLE Subcommands; Konsument-Parsing simpel | `--follow` macht Single-Envelope unmöglich (Unbounded-Stream); `--tail` bounded ist OK aber bricht beim Follow-Pfad |
  | B | **JSON-Lines (NDJSON)** ein Object pro Log-Zeile + Final-Envelope | Streaming-tauglich (auch `--follow`); semantisch logs-natural | Spec-§1841-Bruch (N Objects statt 1); Konsument braucht NDJSON-Parser; eigener Doku-Carveout |
  | C | **Hybrid**: bei `--follow` JSON-Lines, sonst Single-Envelope | beide Welten | zwei Vertragsformen unter einem Flag-Suffix; Konsument muss Detection |

  Plan-Empfehlung: **(B) JSON-Lines** mit expliziter
  cli-json-output.md §6.8-Carveout-Dokumentation. Begründung:
  (a) logs ist semantisch ein Streaming-Subcommand —
  Single-Envelope für `--follow` unmöglich. (b) NDJSON ist
  Industrie-Standard für Streaming-JSON (Docker-Compose
  selbst nutzt es). (c) Hybrid-Form (C) bricht den Vertrag
  noch heftiger als (B). (d) Konsumenten die strukturiert
  parsen können auch NDJSON via `json.Decoder`-Loop. (e)
  Final-Envelope am Ende trägt das `status`/`exitCode`-
  Signal das Spec §1841 fordert.

  **Alternative-Pfad** wäre (A) mit `--follow` Reject (Exit 2
  bei `--follow --json`-Kombi). Trifft auf Konsumenten die
  ihren Output-Buffer überlasten würden. Plan-Vorschlag: nur
  falls Real-World-Push-Back gegen NDJSON.

- **T0-(b) NDJSON-Per-Line-Object-Form** (falls T0-(a) Option B):
  drei Felder gegen Spec §1834 (`level` nur `warn|error`)
  prüfen:
  (i) `{"line": "<raw-compose-output>"}` — schmalste Form;
      `level`-Feld weggelassen.
  (ii) `{"level": "info", "code": "LH-FA-UP-005", "message":
       "<line>"}` — diagnosticItem-Form. Aber `level: "info"`
       ist Spec §1834 NICHT erlaubt — Bruch.
  (iii) `{"service": "postgres", "line": "<line>"}` — angereichert
        mit Compose-Service-Prefix als Sub-Feld.
  Plan-Empfehlung: **(iii)** Per-Line-Object mit `service` +
  `line`-Feld. `level` weglassen weil Compose-Logs keine
  strukturierten Severity-Level haben.

- **T0-(c) Final-Envelope-Form** (Stream-Ende-Marker): die letzte
  NDJSON-Zeile MUSS einen vollen Minimalkontrakt-Envelope
  tragen damit Konsument `status`/`exitCode` ausliest. Format-
  Optionen:
  (i) **Diskriminator-Feld** auf Object-Ebene
      (`{"type": "envelope", "status": "ok", ...}`); Per-Line
      hat `{"type": "line", "service": "...", "line": "..."}`.
      Klare Disambiguation, aber alle Objects bekommen ein
      `type`-Feld.
  (ii) **Implicit by Schema**: das letzte Object ist immer
       der Envelope; Konsument parsed sequentiell und der
       letzte ist der Envelope. Brittle wenn der Stream
       abbricht.
  Plan-Empfehlung: **(i)** Diskriminator-Feld `type:
  "line"|"envelope"` für klare Sequential-Parsing-Semantik.

- **T0-(d) `--follow --json` Final-Envelope-Trigger**: bei
  Unbounded-Stream ist Stream-Ende = SIGINT-Cancel. Wird der
  Final-Envelope vor oder nach dem SIGINT-Handler emittiert?
  (i) Application-Layer ruft `OutputSink.Write` für Final-
      Envelope direkt vor `return (LogsResponse{}, nil)`
      bei Cancel.
  (ii) CLI-Layer (`runLogs`) emittiert Final-Envelope nach
       `useCase.Logs(...)`-Return.
  Plan-Empfehlung: **(ii)** CLI-Layer-Emission damit
  Application-Layer-Vertrag (`LogsResponse{}`) unverändert
  bleibt. Use-Case bleibt Format-agnostisch.

- **T0-(e) FS-Sentinel `ErrLogsFileSystem`?**: heute drei
  FS-Read-Wrap-Stellen ohne typed Sentinel:
  `logsservice.go:??` `checkProjectInitialized` + `??`
  `checkComposeFilePresent` (TBD per Code-Recon in T2).
  Pattern-Erbe up-down T2: zwei neue Sentinels für FS-first
  Switch-Order-Defense. Sub-Decision:
  (i) Neuer `driving.ErrLogsFileSystem` mit Read-Message-Form
      `"logs: filesystem read failed"`. Pattern 1:1 zu
      `ErrUpFileSystem`/`ErrDownFileSystem`.
  (ii) Re-use `driving.ErrUpFileSystem` (semantisch shared
       "filesystem read failed" auf compose.yaml/u-boot.yaml).
       Sentinel-Cluster-Konsolidierung.
  Plan-Empfehlung: **(i)** neuer Sentinel. Pattern-Disziplin
  > Konsolidierung — jeder Subcommand-Pfad bekommt seinen
  eigenen FS-Sentinel-Anker.

- **T0-(f) Mapper-Tabelle** (`mapLogsErrorToDiagnostic`)
  Switch-Order:

  | # | Sentinel | LH-Code | Exit | Mapper-Heim | Begründung |
  | - | -------- | ------- | ---- | ----------- | ---------- |
  | 1 | `driving.ErrLogsFileSystem` (NEU, T2) | `LH-NFA-REL-003` | 14 | `mapLogs` | FS-first damit Multi-`%w` mit FS+Docker auf FS-Klasse fällt |
  | 2 | `driven.ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | `helper` | shared via `mapComposeRuntimeSentinel` aus up-down T5 |
  | 3 | `driven.ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | `helper` | dito |
  | 4 | `driving.ErrComposeFileMissing` | `LH-FA-UP-005` | 10 | `mapLogs` | shared mit up/down (auf LH-FA-UP-001) — aber logs nutzt eigenen LH-Anker (Sub-Decision T0-(g)) |
  | 5 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | `mapLogs` | Pattern-Erbe up/down/generate (Environment-Operation) |
  | 6 | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 | `mapLogs` | Pattern-Erbe init |
  | 7 | `cli.ErrInvalidLogsTail` | `LH-FA-CLI-006` | 2 | `mapLogs` | CLI-Form-Validierung |
  | 8 | Default (unknown) | `LH-FA-CLI-006` | 1 | `mapLogs` | Fallback |

- **T0-(g) `ErrComposeFileMissing` LH-Code-Drift**: up/down
  haben das auf `LH-FA-UP-001` gemappt (Logs-Spec ist aber
  `LH-FA-UP-005`). Sub-Decision: logs-spezifischer
  `LH-FA-UP-005`-Anker ODER up/down-Konsens `LH-FA-UP-001`?
  Cluster-Konvention sagt: same Sentinel → same LH-Code (R4-
  MED-2 R7-Pattern). Plan-Empfehlung: **`LH-FA-UP-001`** für
  Cluster-Konsistenz (logs erbt up/down-Mapping).

- **T0-(h) SIGINT-Vertrag im JSON-Mode**: Final-Envelope ja
  oder nein bei SIGINT-Cancel?
  (i) Bei `--follow --json` mit SIGINT: Final-Envelope mit
      `status: ok`, `exitCode: 0`, kein Diagnostic. Klare
      Stream-Ende-Markierung für Konsument.
  (ii) Kein Output, nur Exit 0 (analog heute). Konsument
       muss am Stream-Ende selbst Schluss machen.
  Plan-Empfehlung: **(i)** Final-Envelope für klare Streaming-
  Semantik.

- **T0-(i) Heute-Validation-Pfad-Drift**: `runLogs:118-121`
  ruft `validateLogsTailFlag` VOR `domain.NewServiceName`. Im
  JSON-Mode bedeutet das: ein Args-Error mit invalid Service-
  Name liefert nicht den Service-Name-Validation-Code sondern
  den Tail-Fehler — Drift gegen Plan-Mapper-Tabelle T0-(f).
  Pattern-Erbe up/down: Pre-UC-Validation läuft via
  `reportError`. Sub-Decision-Form: T5 ergänzt `--json`-
  Awareness in `runLogs` mit `reportError`-Aufruf für jeden
  Validation-Branch.

- **T0-(j) `--quiet --json` Pattern**: heute `runLogs` liest
  `--quiet` nicht (nur Compose-Stream auf OutputSink). Im
  JSON-Mode: `--quiet --json` = `--json` (Cluster-T0-(a)).
  T5 muss `flags.Quiet` lesen UND ignorieren wenn JSON-Mode.

- **T0-(k) Compose-Log-Service-Prefix-Parsing**: Compose
  liefert pro Zeile `<service-name> | <line>` (mit Pipe-
  Separator). Im JSON-Lines-Mode trägt das Per-Line-Object
  `{service, line}` — der Service-Prefix wird vom Application-
  oder CLI-Layer abgeschnitten. Sub-Decision: Wo lebt das
  Parsing?
  (i) Application-Layer (Tee-Writer der pro Zeile parsed und
      Per-Line-Object emittiert).
  (ii) CLI-Layer (Wrapping-OutputSink der das Compose-Output
       zerlegt und als NDJSON emittiert).
  Plan-Empfehlung: **(ii)** CLI-Layer — Application-Layer
  bleibt Format-agnostisch.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(k) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog up-down T1): `cli/sanitize.go` + `cli/composesentinel.go`-Helper existieren bereits aus up-down T5 | — (entfällt) | T0 |
| T2 | Port-Types: ggf. neuer `LogsRequest.SilenceProgress bool` (sehr wahrscheinlich nicht — logs ist OutputSink-driven, kein Progress); `driving.ErrLogsFileSystem`-Sentinel (T0-(e) Option (i)); Co-Migration der Port-Sentinel-Kommentare (`logs.go:??` falls heute generische `LH-FA-CLI-006`-Anker) | ~60 | T0 |
| T3 | Application-Layer: Multi-`%w`-Wrap-Migration der FS-Read-Stellen auf `ErrLogsFileSystem`. KEIN ProgressSink-Branch nötig (OutputSink ist Stream-Sink, nicht Phase-Sink). | ~30 | T2 |
| T4 | **Entfällt-Kandidat**: Composition-Root `cmd/uboot/main.go` hat heute schon `NewLogsService` mit allen Deps. Kein Wiring-Update bei Bool-Field-Pattern (falls T2 keinen Bool einführt). | — (entfällt erwartet) | T3 |
| T5 | CLI-RunE: `runLogs` ergänzt um `--json`/`--quiet`-Awareness; Allowlist-Migration; neuer `mapLogsErrorToDiagnostic`; Pre-UC-Validation-Pfade via `reportError`; **NDJSON-OutputSink-Wrapping** (T0-(k) Option (ii)) — wrappt `cmd.OutOrStdout()` in einen `ndjsonOutputSink` der Per-Line-Object emittiert; Final-Envelope-Emission nach `useCase.Logs(...)`-Return (T0-(c) Diskriminator-Feld + T0-(d) CLI-Layer-Emission); Sanitizer-Aufrufe via `cli/sanitize.go`. | ~200 | T2 |
| T6 | Acceptance-Tests: ~12-15 Tests (Single-Tail-Pin, Follow-Stream-Pin, --quiet --json, Mapper-Rows, SIGINT-Vertrag, Final-Envelope-Form, NDJSON-Per-Line-Object-Form, Path-Leak-Sanitizer) | ~450-550 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe) | ~50 | T6 |
| T8 | Closure: CHANGELOG, `cli-json-output.md` §6/§6.8/§7 mit **NDJSON-Streaming-Carveout-Doku** in §6.8 (erster Subcommand der Spec-§1841 als N-Objects-pro-Aufruf interpretiert — explizit dokumentieren), roadmap done-Zähler 6→7, Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: **~790-890** (deutlich kleiner als up-
down ~1035-1135 weil keine zwei Subcommands zu bündeln,
kein zweiter FS-Sentinel, keine zwei Mapper-Files). Pattern-
Erbe von up-down (FS-Sentinel-Pattern + Mapper-Switch-Order +
Sanitizer-Helper-Wiederverwendung + ComposeRuntime-Helper) und
remove (`reportError`-Helper-Form für Pre-UC-Validation-
Pfade).

## Out of Scope

- **`--no-log-prefix` / `--timestamps`** (Spec-Erweiterung für
  Compose-Logs-Format-Flags): bewusste Logs-Slice-Erweiterung
  außerhalb des V1-Scope. Pattern-Vorbild Compose-CLI direkt
  passend; Folge-Slice falls Real-World-Druck.
- **Multi-Service-Filter** (`u-boot logs svc1 svc2`): heute
  Single-Service via `cobra.MaximumNArgs(1)`. Multi-Args-Form
  wäre Spec-Erweiterung (LH-FA-UP-005 spricht Singular); eigener
  Folge-Slice.
- **`--since` / `--until` Time-Range-Filter**: nicht in Spec;
  Folge-Slice falls Real-World-Druck.
- **WARN-Migration**: `driving.WarningEntry`-Type ist aus
  remove T2 verfügbar, aber logs hat heute keine bekannten
  WARN-Pfade (kein Recreate-Detection, kein deferred-
  Volumes-Pattern). Falls künftige Erweiterung WARN braucht
  (z. B. "service has no logs"), wandert das in den
  Recreate-Detection-Folge-Slice oder einen eigenen.
- **JSON-Lines vs. Spec-§1841 Cluster-Audit**: wenn Cluster-
  T_close (nach 9/9) auf strikten Single-Envelope-Vertrag
  besteht, müsste logs auf Option A (Single-Envelope mit
  `--follow`-Reject) migrieren. Heute bewusste Sub-Decision
  in diesem Slice; Re-Eval bei T_close erlaubt.

## Bezug

- Cluster:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  (Folge-Slice 7/9).
- Pattern-Vorbilder:
  [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
  (Read-only-Klassen-Disziplin + FS-Sentinel-Pattern +
  Sanitizer-Helper-Quelle + ComposeRuntime-Helper-Quelle),
  [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
  (`reportError`-Helper-Form für Pre-UC-Validation),
  [`slice-v1-logs`](../done/slice-v1-logs.md) (M6-Logs-Auslieferung
  mit T0-Outcomes — der `--json`-Pfad ist bewusst hierher
  ausgelagert worden).
- Code-Anker:
  [`cli/logs.go`](../../../../internal/adapter/driving/cli/logs.go),
  [`application/logsservice.go`](../../../../internal/hexagon/application/logsservice.go),
  [`port/driving/logs.go`](../../../../internal/hexagon/port/driving/logs.go),
  [`cli/jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  Z. 29/74.
- Folge-Slices: keine direkten Forward-Refs aus logs heraus;
  `slice-v1-recreate-detection` ist up-down-Carveout (nicht
  logs).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
