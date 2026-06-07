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
| FS-Read | `u-boot.yaml` + `compose.yaml` Pre-Checks (analog up/down `checkProjectInitialized` + `checkComposeFile`) |
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
  down-T5 — der Helper lebt seit up-down T5 in
  `cli/sanitize.go` package-intern (R1-LOW-1 Redundanz-
  Konsolidierung).
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

- **T0-(a) Output-Modell** — **zentrale Sub-Decision, OFFEN**
  (Cluster-Plan Z. 326-329 hierher ausgelagert; R1-HIGH-1
  De-Festlegung): drei Optionen, **keine vorab-präferiert**.
  Entscheidungs-Kriterien in R2/R3-Review-Runden klären.

  | # | Form | Vor | Contra |
  | - | --- | --- | --- |
  | A | **Single-Envelope + `--follow --json` Reject (Exit 2)** | Spec-§1841-konform; cluster-konsistent mit 6/9 done-Slices; eine Envelope-Form Cluster-weit; Konsument-Parsing simpel (`json.Unmarshal`); Pattern-Erbe up-down 1:1; ggf. migrierbar zu (B) falls Real-World-Push-Back; kein T_close-Audit-Pflicht | `--follow --json` Konsumenten verlieren strukturierten Output; bounded-Pfad (`--tail=N --json`) bleibt funktional; alle Log-Zeilen müssen im Use-Case gepuffert werden (Memory-Skalierung bei großen `--tail`-Werten) |
  | B | **JSON-Lines (NDJSON)** ein Object pro Log-Zeile + Final-Envelope | Streaming-tauglich (auch `--follow`); semantisch logs-natural; **Konsument-O(1)-Stream-Processing** (keine Buffer-Wartezeit bis Stream-Ende); **memory-konstant** auch bei Unbounded-Stream; Konsument kann jede Zeile sofort verarbeiten/filtern; Pattern-Erbe Industrie-Standards (Docker `--format json`, `journalctl -o json` opt-in) | Spec-§1841-Bruch (N Objects statt 1); Konsument braucht NDJSON-Parser (`json.Decoder`-Loop); eigener Doku-Carveout in `cli-json-output.md` §6.8; bricht Cluster-Pattern; Cluster-Slice T_close-Tranche-Konsens-Pflicht (R2-LOW-2) |
  | C | **Hybrid**: bei `--follow` NDJSON, sonst Single-Envelope | beide Welten | zwei Vertragsformen unter einem Flag-Suffix; Konsument muss Format-Detection machen; bricht §1841 nur unter `--follow` aber inkonsistent — der heftigere Vertrags-Bruch von (B) wird durch (A)-Form unter `--tail` getarnt; Cluster-Audit-Aufwand vermutlich identisch zu (B) |

  **Entscheidungs-Kriterien für R2/R3:**
  - **Real-World-Konsumenten-Belege** (R2-MED-3
    Operationalisierung): gibt es heute belegte CI-Skripte /
    Konsumenten die `--follow --json` brauchen? **Such-
    Pflicht** in der T0-Discovery-Tranche, konkrete Pfade:
    (i) `.github/workflows/*.yml` — alle CI-Workflows nach
        `u-boot logs`-Aufrufen mit `--json`-Flag durchsuchen.
    (ii) `Makefile` + `scripts/` — Build-/Dev-Skripte nach
         `logs --json`-Pattern grep'en.
    (iii) `docs/user/` + `docs/plan/` — Doku-Quellen nach
          Konsumenten-Use-Cases für strukturiertes
          `logs --json` durchsuchen.
    (iv) externes Issue-Tracking — gibt es offene Issues/
         PRs die `logs --json --follow` einfordern?
    Wenn alle vier Pfade Null-Belege liefern → **Memory
    `diagnose_vor_carveout`-konform: (A) wählen**, nicht
    spekulieren. Belegs-Pflicht ist Pre-`next/`-Trigger für
    T0-(a)-Entscheidung.
  - **Cluster-T_close-Konsequenz**: bei (B) muss T_close-Audit
    den NDJSON-Carveout explizit absegnen (oder logs auf (A)
    migrieren). Bei (A) keine T_close-Sonderbehandlung nötig.
  - **Real-World-Pattern-Vergleich** (korrekt zitiert): Docker-
    Compose-Default liefert plain text; `--format json`
    erst auf Opt-in NDJSON. `kubectl logs` plain text.
    `journalctl -o json` NDJSON aber `-o json` ist explizit.
    **Keiner der Standards mischt Single-Envelope mit
    Streaming unter einem Flag** — sie haben einen
    expliziten Format-Switch. u-boot's `--json` ist heute
    cluster-weit das Format-Signal — kein Sub-Switch
    erforderlich, wenn (A) gewählt wird.
  - **NDJSON-Per-Line-Schema-Pflicht** (falls (B)): wenn
    "Industrie-Standard" das Argument ist, muss das Per-
    Line-Object dem Docker-Compose-`--format json`-Schema
    folgen (R1-HIGH-2). Siehe T0-(b).
  - **Diskriminator-Field-Schema-Erweiterung** (falls (B)):
    `cliJSONEnvelope` (`jsonenvelope.go:21-50`) hat kein
    `type`-Feld. Erweiterung oder logs-spezifischer Wire-
    Type pflichtig (R1-HIGH-3). Siehe T0-(c).

  **R2/R3-Pflicht**: eine der drei Optionen mit dokumentierten
  Konsumenten-Belegen oder Cluster-T_close-Konsens festzurren.
  Vorab-Festlegung ist ausgeschlossen — die Begründung muss
  belastbar sein.

- **T0-(b) NDJSON-Per-Line-Object-Form** (relevant nur falls
  T0-(a) Option B, R1-HIGH-2 Schema-Konsistenz-Korrektur):
  drei Optionen mit Industrie-Standard-Belegen:
  (i) `{"line": "<raw-compose-output>"}` — schmalste Form;
      `level`-Feld weggelassen. Konsument muss Service-Prefix
      selbst parsen.
  (ii) `{"service": "postgres", "line": "<line>"}` — schmaler
       angereichert mit Compose-Service-Prefix als Sub-Feld.
       **R1-Drift**: kein Industrie-Standard hat genau diese
       Form.
  (iii) `{"time": "<ts>", "service": "<name>", "container":
        "<container-id>", "log": "<line>"}` — **Docker-
        Compose-`--format json`-Schema 1:1** (R1-HIGH-2
        Konsistenz: wenn (B) "Industrie-Standard" zitiert,
        muss das Per-Line-Object dem Standard folgen).
        **R2-MED-1 Realitäts-Constraint**: heutige Stream-
        Quelle liefert nur `<service>-<idx>`-Prefix
        (Container-Index, NICHT Container-ID). Echtes
        `container`-Feld bräuchte einen Pre-Walk via
        `compose ps`-Roundtrip (Out-of-Scope für diesen
        Slice). Lösungs-Form: **`container` als
        `omitempty`**-Field — bleibt heute weg, Schema bleibt
        forward-kompatibel falls künftiger Slice den
        Pre-Walk einführt. Field-Drop heute, Field-Add
        morgen.
  Plan-Empfehlung **bedingt auf T0-(a) Option B**: **(iii) mit
  `container omitempty`** — Docker-Compose-Schema forward-
  kompatibel, ohne Pre-Walk-Pflicht. `level` weglassen weil
  Compose-Logs keine strukturierten Severity-Level haben
  (Spec §1834 erlaubt sowieso nur `warn|error`, nicht
  `info`). Implementations-Pflicht: CLI-Layer-Parsing liefert
  drei Felder (`service` aus Prefix, `log` aus Zeile, `time`
  aus `--timestamps`-Form falls jemals aktiviert);
  `container` bleibt heute weg.

- **T0-(c) Final-Envelope-Form** (Stream-Ende-Marker, relevant
  nur falls T0-(a) Option B oder C, R1-HIGH-3 Schema-Bruch-
  Auflösung): die letzte NDJSON-Zeile MUSS einen vollen
  Minimalkontrakt-Envelope tragen damit Konsument
  `status`/`exitCode` ausliest. **R1-HIGH-3-Konflikt**: das
  `cliJSONEnvelope`-Schema (`jsonenvelope.go:21-50`) hat
  KEIN `type`-Feld. Drei Optionen:
  (i) **Eigener `logsLineEnvelope`-Wire-Type**: separater
      Wire-Struct nur in `cli/logs.go`, kein
      `cliJSONEnvelope`-Schema-Refactor. Per-Line-Form ist
      logs-spezifisch (T0-(b)); Final-Envelope nutzt aber
      `cliJSONEnvelope`-Struct ohne Diskriminator-Feld.
      Konsument trennt Per-Line von Final-Envelope am
      Feld-Set (`time`/`service`/`container`/`log` →
      Per-Line; `status`/`command`/`diagnostics`/`exitCode`
      → Envelope). Schmaler Eingriff.
  (ii) **`cliJSONEnvelope.Type *string`-Schema-Erweiterung**:
       `cliJSONEnvelope`-Struct um `Type *string
       json:"type,omitempty"` ergänzen. Existierende done-
       Slices setzen das nicht (Pointer-omitempty → Feld
       fällt weg). Logs setzt `type: "envelope"`; Per-Line
       hat eigenen Wire-Type mit `type: "line"`. Konsistente
       Diskriminator-Form aber Schema-Erweiterung berührt
       6 done-Slices (Tests müssen verifizieren dass Feld
       weiterhin abwesend bleibt).
  (iii) **Diskriminator-loose**: keine `type`-Felder; Per-
        Line-Form (i) ohne Diskriminator. Konsument nutzt
        Sequential-Parse: jedes Object außer dem letzten ist
        Per-Line; das letzte ist Envelope. Brittle bei
        Stream-Abbruch.
  Plan-Empfehlung **bedingt auf T0-(a) Option B**: **(i)
  Eigener `logsLineEnvelope`-Wire-Type** — schmalster
  Eingriff, kein `cliJSONEnvelope`-Schema-Refactor, kein
  Risiko an Pinnt-Tests in `jsonenvelope_test.go`. Konsument
  disambiguiert über das Feld-Set, nicht über expliziten
  Diskriminator.

- **T0-(d) Final-Envelope-Trigger** (relevant falls T0-(a) B/C):
  bei Unbounded-Stream ist Stream-Ende = SIGINT-Cancel ODER
  natural-end (bounded `--tail` ohne `--follow`). **R1-HIGH-4
  Plan-Konflikt-Auflösung**: heute liefert
  `LogsService.Logs` für **beide Fälle** `(LogsResponse{},
  nil)` (`logsservice.go:102-110`). CLI-Layer kann SIGINT
  NICHT von natural-end unterscheiden. Drei Lösungs-Pfade:
  (i) **CLI-Layer prüft `ctx.Err()` post-UC-Return**: nach
      `useCase.Logs(...)` ohne strukturierte Response
      pürft der CLI `ctx.Err() == context.Canceled` und
      setzt Final-Envelope-`status: ok`, `exitCode: 0`
      mit optionaler Info-Diagnostic
      `"stream cancelled by SIGINT"` ODER ohne Diagnostic.
      Schmaler Eingriff, keine Port-Type-Änderung.
  (ii) **`LogsResponse.TerminatedBy string`-Feld**:
       Port-Type-Erweiterung in T2. Application-Layer setzt
       `TerminatedBy: "stream-end"` ODER `"cancel"`. CLI-
       Layer liest das Feld direkt. Klarer Vertrag aber
       Port-Type-Refactor.
  (iii) **Application-Layer emittiert Final-Envelope selbst**:
        widerspricht T0-(d)-Plan-Form (CLI-Layer-Emission)
        und bricht Format-Agnostik der Application-Layer.
        Verworfen.
  Plan-Empfehlung **bedingt auf T0-(a) Option B/C**: **(i)
  `ctx.Err()`-Check post-UC**. Schmalster Eingriff, keine
  Port-Type-Änderung nötig. T0-(d) Plan-Form bleibt erhalten
  (CLI-Layer-Emission). Pattern-Erbe `runLogs`-`runConfirmation
  Gate` Pre-/Post-Check-Idiomen.

  **R2-LOW-3 Race-Toleranz absegnen**: bei `--follow --tail=10`
  mit Stream-Exhaust + SIGINT gleichzeitig (Compose-Subprozess-
  Exit und Signal-Delivery simultan) ist `ctx.Err()` race-abhängig
  — entweder nil (Stream-Ende-vor-Signal) oder `Canceled`
  (Signal-vor-Stream-Ende). **Semantisch identisch**: beide
  Pfade bedeuten "Stream sauber beendet" für den Konsument; der
  Final-Envelope trägt `status: ok`/`exitCode: 0` in beiden
  Fällen. Race ist **bewusst toleriert** weil die User-facing
  Semantik gleich ist — kein Test-Pin nötig, kein
  Determinismus-Carveout.

  Bei Option (A) ist diese Sub-Decision irrelevant —
  Single-Envelope wird nach UC-Return immer geschrieben (kein
  `--follow --json`-Pfad weil Reject).

- **T0-(e) FS-Sentinel `ErrLogsFileSystem`** (R1-MED-1
  Festzurrung, keine "ggf."-Aufweichung): heute **zwei** FS-
  Read-Wrap-Stellen ohne typed Sentinel:
  `logsservice.go:117-127` `checkProjectInitialized` (Z. ~121
  `Exists(%q): %w`) + `:133-143` `checkComposeFile`
  (Z. ~137 `Exists(%q): %w`). Pattern-Erbe up-down T2: neue
  Sentinels für FS-first Switch-Order-Defense sind **Pflicht**,
  nicht optional. Sub-Decision-Pfad:
  (i) **Neuer `driving.ErrLogsFileSystem`** mit Read-Message-
      Form `"logs: filesystem read failed"`. Pattern 1:1 zu
      `ErrUpFileSystem`/`ErrDownFileSystem`.
  (ii) Re-use `driving.ErrUpFileSystem` (semantisch shared
       "filesystem read failed" auf compose.yaml/u-boot.yaml).
       Sentinel-Cluster-Konsolidierung.
  Plan-Empfehlung **festgezurrt**: **(i)** neuer Sentinel.
  Pattern-Disziplin > Konsolidierung — jeder Subcommand-Pfad
  bekommt seinen eigenen FS-Sentinel-Anker (Cluster-
  Konvention aus up-down R3-MED-1). T2-Cell-Wortlaut auf
  "neuer Sentinel" (statt "ggf.").

- **T0-(f) Mapper-Tabelle** (`mapLogsErrorToDiagnostic`)
  Switch-Order:

  | # | Sentinel | LH-Code | Exit | Mapper-Heim | Begründung |
  | - | -------- | ------- | ---- | ----------- | ---------- |
  | 1 | `driving.ErrLogsFileSystem` (NEU, T2) | `LH-NFA-REL-003` | 14 | `mapLogs` | FS-first damit Multi-`%w` mit FS+Docker auf FS-Klasse fällt |
  | 2 | `driven.ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | `helper` | shared via `mapComposeRuntimeSentinel` aus up-down T5 |
  | 3 | `driven.ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | `helper` | dito |
  | 4 | `driving.ErrComposeFileMissing` | `LH-FA-UP-001` | 10 | `mapLogs` | Cluster-Konsens mit up/down (T0-(g) festgezurrt: same Sentinel → same LH-Code) |
  | 5 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | `mapLogs` | Pattern-Erbe up/down/generate (Environment-Operation) |
  | 6 | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 | `mapLogs` | Pattern-Erbe init |
  | 7 | `cli.ErrInvalidLogsTail` | `LH-FA-CLI-006` | 2 | `mapLogs` | CLI-Form-Validierung |
  | 8 | Default (unknown) | `LH-FA-CLI-006` | 1 | `mapLogs` | Fallback |

- **T0-(g) `ErrComposeFileMissing` LH-Code Cluster-Konsens**
  (R1-MED-2 festgezurrt): up/down haben das auf
  `LH-FA-UP-001` gemappt. Cluster-Konvention "same Sentinel →
  same LH-Code" (R4-MED-2 Pattern aus up-down) gilt auch für
  logs — `ErrComposeFileMissing` ist derselbe Port-Sentinel
  egal welcher Subcommand ihn auslöst. Plan-Empfehlung
  **festgezurrt**: **`LH-FA-UP-001`** für Cluster-Konsistenz.
  Mapper-Tabelle Z. 235 entsprechend korrigiert (R1-MED-2
  Tabellen-Drift behoben).

- **T0-(h) SIGINT-Vertrag im JSON-Mode** (R1-HIGH-4 plan-
  intern aufgelöst via T0-(d)(i) `ctx.Err()`-Check):
  Final-Envelope-Emission und SIGINT-Distinction sind durch
  T0-(d) gelöst. Sub-Decision-Form:
  - **Option (A) Single-Envelope**: Single-Envelope nach
    UC-Return mit `status: ok` (bounded `--tail`-Pfad) ODER
    Final-Envelope-mit-Diagnostic falls UC-Error. Kein
    SIGINT-Pfad weil `--follow --json` rejected.
  - **Option (B/C) NDJSON**: CLI-Layer prüft `ctx.Err()` post-
    UC. Bei `context.Canceled` → Final-Envelope `{type:
    "envelope" (falls T0-(c)(ii)) | implicit-by-feldset (falls
    T0-(c)(i)), status: "ok", exitCode: 0}`. Bei natural-end
    identisches Format. Konsument disambiguiert NICHT — beide
    Pfade sind aus Konsument-Sicht "Stream sauber beendet".
  Plan-Empfehlung **bedingt auf T0-(a)**: konsistent mit
  T0-(d)(i) `ctx.Err()`-Check. Application-Layer-Vertrag
  (`LogsResponse{}`, kein TerminatedBy-Feld) bleibt
  unverändert.

- **T0-(i) Heute-Validation-Pfad-Drift**: `runLogs:118-121`
  ruft `validateLogsTailFlag` VOR `domain.NewServiceName`. Im
  JSON-Mode bedeutet das: ein Args-Error mit invalid Service-
  Name liefert nicht den Service-Name-Validation-Code sondern
  den Tail-Fehler — Drift gegen Plan-Mapper-Tabelle T0-(f).
  Pattern-Erbe up/down: Pre-UC-Validation läuft via
  `reportError`. Sub-Decision-Form: T5 ergänzt `--json`-
  Awareness in `runLogs` mit `reportError`-Aufruf für jeden
  Validation-Branch.

- **T0-(j) `--quiet --json` + `runLogs`-Signatur** (R1-MED-4
  Signatur-Form pinnen): heute `runLogs(ctx, out, args,
  flags, uc, getwd)` (`logs.go:111-118`) liest weder
  `a.quiet` noch `a.json` — Pattern-Bruch zu
  doctor/add/init/generate/remove/up-down die alle
  `a.json`/`a.quiet` via App-State oder durchgereichte
  bool-Parameter lesen. T5-Pflicht: **Signatur-Refactor**.
  Drei Optionen:
  (i) **`runLogs(ctx, out, args, flags, uc, getwd, jsonMode
      bool, quietMode bool)`** — zwei bool-Parameter
      durchreichen. Pattern-Erbe up-down `runUp`/`runDown`.
  (ii) **`runLogs(ctx, out, args, flags{Quiet, JSON, …}, uc,
       getwd)`** — `logsFlags`-Struct um `Quiet bool`
       und `JSON bool` erweitern. Closure liest `a.quiet`/
       `a.json` und füllt die Fields VOR `runLogs`-Aufruf
       (Pattern-Erbe add/init/generate/remove).
  (iii) `runLogs(ctx, out, args, flags, uc, getwd, a *App)`
        — App-Struct durchreichen. Bricht Testbarkeit-
        Pattern.
  Plan-Empfehlung **festgezurrt**: **(ii)** `logsFlags`-Struct
  um `Quiet` + `JSON` Fields erweitern, Closure liest
  `a.quiet`/`a.json` (Pattern-Erbe add/init/generate/remove —
  Closure-Idiom Z. 88-95 in `cli/up.go` als Vorbild).
  **`--quiet --json` Semantik** (Cluster-T0-(a) doctor-
  Pattern): im JSON-Mode ist `--quiet` ein No-Op weil JSON
  immer auf stdout muss. Im Human-Mode kann `--quiet`
  weiterhin die heutigen `Compose-Stream`-Output
  unterdrücken (oder bleibt Status-quo "Compose-Stream-Output
  ignoriert `--quiet`" — siehe Out-of-Scope).

- **T0-(k) Compose-Log-Output-Form** (R1-MED-5 Format-
  Korrektur, R1-LOW-4 Belastbarkeit): Compose liefert pro
  Zeile **`<service-name>-<idx>  | <line>`** (Container-Index
  `-1`/`-2`, ZWEI Leerzeichen vor dem Pipe, NICHT `<service>
  | <line>` wie Pre-R1-Form annahm). Plus: Adapter
  (`docker_engine.go:99-113`) mischt **stdout UND stderr**
  in den `OutputSink` — stderr-Lines vom Adapter
  (`"Attaching to postgres"`, `"postgres exited with code 1"`)
  sind als reguläre `line`-Objects ODER `diagnostic`-Items
  zu klassifizieren. Plus: Adapter macht bereits
  Line-Buffering (`engine.go:166-183` `runLineBuffered`),
  also empfängt der Sink Line-für-Line, nicht Byte-für-Byte
  — gut für NDJSON-Wrapping.

  Sub-Decision Form (relevant nur falls T0-(a) Option B):
  Wo lebt das Compose-Prefix-Parsing?
  (i) Application-Layer (Tee-Writer der pro Zeile parsed
      und Per-Line-Object emittiert).
  (ii) CLI-Layer (Wrapping-OutputSink der das Compose-Output
       zerlegt und als NDJSON emittiert).
  Plan-Empfehlung **bedingt auf T0-(a) Option B**: **(ii)**
  CLI-Layer — Application-Layer bleibt Format-agnostisch
  (`LogsRequest.OutputSink io.Writer` ist heute der direkte
  Compose-Stream-Pfad; CLI-Wrapping ändert das nicht).

  **stderr-vs-stdout-Klassifikation**: stderr-Lines wie
  `"Attaching to <svc>"`/`"<svc> exited"` sind Compose-
  Steuersignale, nicht Service-Logs. **R2-MED-4 Adapter-
  Constraint**: heutiger Adapter (`engine.go:163` ruft
  `runLineBuffered(cmd, progressSinkOrDiscard(opts.Sink))`)
  mixt stdout UND stderr in **einen** Sink. Trennung wäre
  Driven-Port-Refactor (Out-of-Scope). Sub-Optionen:
  (a) Alle Lines (stdout+stderr) als `line`-Objects mit
      `service`-Feld aus Prefix; stderr-Lines tragen
      `service: null` (kein Prefix-Match möglich), `line:
      "<raw>"`. Konsument-Filter über `service`-Field-
      Presence.
  (b) **VERWORFEN**: Stderr-Lines als `diagnostic`-Items
      mit `level: "info"` — Spec §1834 erlaubt sowieso nur
      `warn|error`, plus Adapter-Refactor wäre nötig.
  (c) Stderr-Lines unterdrücken — heute kein Konsument-
      Bedarf, aber verliert Compose-Steuersignal-Information.
  Plan-Empfehlung: **(a)** — single-Sink-Constraint des
  heutigen Adapters konsistent unterstützt; Konsument-Filter
  via `service`-Field-Presence.

  Bei T0-(a) Option (A) Single-Envelope: kein Per-Line-
  Parsing nötig — Compose-Stream wird gesammelt und in
  `data.lines []string` (oder ähnliche Form, T2-Sub-
  Decision) gebündelt.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(k) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog up-down T1): `cli/sanitize.go` + `cli/composesentinel.go`-Helper existieren bereits aus up-down T5 | — (entfällt) | T0 |
| T2 | Port-Types: **`driving.ErrLogsFileSystem`-Sentinel** (T0-(e) Option (i) festgezurrt, R1-MED-1); Read-spezifische Message-Form `"logs: filesystem read failed"`; Heim-Position in `port/driving/logs.go` analog `up.go` (vor `ErrComposeFileMissing`-Cluster). Plus: **`logsFlags.JSON bool` + `logsFlags.Quiet bool` Felder** im CLI-Layer-Struct (T0-(j)(ii) festgezurrt, R1-MED-4). KEIN `SilenceProgress`-Field — `LogsRequest` hat `OutputSink io.Writer`, kein ProgressSink. Co-Migration der heutigen Port-Sentinel-Kommentare falls heute generische `LH-FA-CLI-006`-Anker (Code-Recon in T2). T4 entfällt (kein Composition-Root-Wechsel). | ~70 | T0 |
| T3 | Application-Layer: Multi-`%w`-Wrap-Migration der **zwei FS-Read-Stellen** (`logsservice.go:117-127` `checkProjectInitialized` + `:133-143` `checkComposeFile`) auf `ErrLogsFileSystem`. KEIN ProgressSink-Branch nötig (OutputSink ist Stream-Sink, nicht Phase-Sink). KEIN `LogsResponse`-Field-Erweiterung (`TerminatedBy`-Feld verworfen via T0-(d)(i) `ctx.Err()`-Check). | ~30 | T2 |
| T4 | **Entfällt** (analog up-down T4): Composition-Root `cmd/uboot/main.go` hat heute schon `NewLogsService` mit allen Deps. T2 führt nur Port-Sentinel + CLI-Flag-Fields ein — kein Service-Wiring-Wechsel. | — (entfällt) | T3 |
| T5 | CLI-RunE: **`runLogs(ctx, stdout, errOut io.Writer, args, flags, uc, getwd)`-Signatur-Refactor** (R2-HIGH-3 Cluster-Pattern-Konsistenz mit up/down/remove `up.go:133`/`down.go:128`/`remove.go:253`) + `logsFlags.JSON`/`logsFlags.Quiet` Fields durchreichen analog up/down (T0-(j)(ii)). Allowlist-Migration `"u-boot logs": true` in `jsonAllowlist()`. **`isFilesystemError`-Co-Migration** (`cli/cli.go:401-428`, R1-MED-6): `driving.ErrLogsFileSystem` ergänzen damit Exit-Code-Mapping auf 14 fällt. Neuer `mapLogsErrorToDiagnostic` mit Switch-Order T0-(f). Pre-UC-Validation-Pfade via `reportError` für Single-Envelope-Form (Option (A)) ODER **Single-Object-NDJSON-Pre-Stream-Output** im Stream-Pfad für NDJSON (Option (B), R1-MED-3 Aufschlüsselung — `reportError` ist NDJSON-Frame-konform weil `writeEnvelope`/`fmt.Fprintln` einen Newline-terminierten Single-Object schreibt; kein Wrap-Layer nötig). Bei Option (B): **NDJSON-OutputSink-Wrapping** (T0-(k)(ii)) — wrappt `cmd.OutOrStdout()` in `ndjsonOutputSink` der Per-Line-Object emittiert; Final-Envelope-Emission nach `useCase.Logs(...)`-Return mit `ctx.Err()`-Check (T0-(d)(i)); eigener `logsLineEnvelope`-Wire-Type (T0-(c)(i)) plus `writeLogsLineEnvelope`-Constructor analog `writeEnvelope`. Sanitizer-Aufrufe via `cli/sanitize.go`. | (A) ~150-200 / (B) ~250-300 | T2 |
| T6 | Acceptance-Tests: ~14-18 Tests (bounded `--tail`-Pin, `--follow --json`-Pin je T0-(a) Form, `--quiet --json`-Pin, Mapper-Rows 1-8, SIGINT-Vertrag (Option (B) only), Final-Envelope-Form (Option (B) only), NDJSON-Per-Line-Object-Form (Option (B) only), Path-Leak-Sanitizer-Pin, `--follow --json`-Reject (Option (A) only)) | (A) ~400-500 / (B) ~500-600 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe) | ~50 | T6 |
| T8 | Closure: CHANGELOG, `cli-json-output.md` §6/§6.8/§7 (Form je T0-(a)-Wahl: bei (A) §6.8 als reguläre Read-only-Sektion; bei (B) §6.8 mit **NDJSON-Streaming-Carveout-Doku** als erster Subcommand der Spec-§1841 als N-Objects-pro-Aufruf interpretiert — explizit dokumentieren), roadmap done-Zähler 6→7, carveouts.md-Einträge für die drei `open/`-Stubs (`slice-v1-logs-format-flags`, `-multi-service-filter`, `-time-range-filter` — bereits in `open/` angelegt bei R2-Adressierung, R2-HIGH-2 Memory-Disziplin), Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) | T7 |

LOC-Bilanz vorläufig: **~900-1050** (R1-LOW-3 Korrektur — T5 +50-100 für NDJSON-Wrapping/Line-Parsing falls Option B, T6 +50 für NDJSON-Stream-Decoding-Tests). Bei T0-(a) Option (A) wird die Bilanz auf ~780-880 reduziert (kein NDJSON-Wrapping). Deutlich kleiner als up-
down ~1035-1135 weil keine zwei Subcommands zu bündeln,
kein zweiter FS-Sentinel, keine zwei Mapper-Files). Pattern-
Erbe von up-down (FS-Sentinel-Pattern + Mapper-Switch-Order +
Sanitizer-Helper-Wiederverwendung + ComposeRuntime-Helper) und
remove (`reportError`-Helper-Form für Pre-UC-Validation-
Pfade).

## Out of Scope

Memory-Feedback `carveouts_need_plans` (R1-MED-7 + R2-HIGH-2):
**sofort** mit Plan in `open/` versehen. Drei `open/`-Stubs
sind bei R2-Adressierung angelegt (Verzicht auf T8-Verschiebung
um Memory-konform zu bleiben).

- **`--no-log-prefix` / `--timestamps`** (Spec-Erweiterung für
  Compose-Logs-Format-Flags): bewusste Logs-Slice-Erweiterung
  außerhalb des V1-Scope. Pattern-Vorbild Compose-CLI direkt
  passend. Plan-Stub:
  [`slice-v1-logs-format-flags`](slice-v1-logs-format-flags.md)
  (`open/`, Status `on hold pending trigger`).
- **Multi-Service-Filter** (`u-boot logs svc1 svc2`): heute
  Single-Service via `cobra.MaximumNArgs(1)`. Multi-Args-Form
  wäre Spec-Erweiterung (LH-FA-UP-005 spricht Singular).
  Plan-Stub:
  [`slice-v1-logs-multi-service-filter`](slice-v1-logs-multi-service-filter.md)
  (`open/`, Status `on hold pending trigger`).
- **`--since` / `--until` Time-Range-Filter**: nicht in Spec.
  Plan-Stub:
  [`slice-v1-logs-time-range-filter`](slice-v1-logs-time-range-filter.md)
  (`open/`, Status `on hold pending trigger`).
- **WARN-Migration**: `driving.WarningEntry`-Type ist aus
  remove T2 verfügbar, aber logs hat heute keine bekannten
  WARN-Pfade. **KEIN eigener Folge-Slice-Stub** — falls
  künftige Erweiterung WARN braucht (z. B. "service has no
  logs"), wandert das in den `slice-v1-recreate-detection`-
  Folge-Slice (Memory-Wieder-Verknüpfung mit existing
  up-down-Carveout-Stub).
- **JSON-Lines vs. Spec-§1841 Cluster-Audit**: bei T0-(a)
  Option (B) wird Cluster-T_close auf NDJSON-Vertrag-Audit
  pflicht; bei Option (A) entfällt das. **KEIN eigener
  Folge-Slice-Stub** — gehört zu Cluster-Slice T_close-Tranche, nicht
  zu logs-Folge.

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
