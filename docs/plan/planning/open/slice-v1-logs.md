# Slice V1: `u-boot logs` (`LH-FA-UP-005`)

> **Status:** geplant für v0.4.0 — Spec ✅
> ([`spec/lastenheft.md:1023-1040`](../../../../spec/lastenheft.md)),
> Port-Anker ✅
> ([`internal/hexagon/port/driving/README.md:39`](../../../../internal/hexagon/port/driving/README.md)
> nennt `LogsUseCase` als V1-Erweiterung, driven/README §"Geplante
> Erweiterungen" listet `Logs`/`Exec`-Verb auf `DockerEngine`),
> Implementation ausstehend. T0 (Discovery + Design) festgehalten
> beim Übergang nach `next/`; Tranchen-Schnitt wird dort
> verfeinert.

## Auslöser

`LH-FA-UP-005` ist die einzige fehlende `u-boot up`/`down`-
Familie-Spec-ID in der V1-Phase: M6 hat `up`/`down` ausgeliefert
(`LH-FA-UP-001..004`), `logs` ist die V1-Erweiterung. Heute muss
der User für Logs direkt auf `docker compose -f compose.yaml
logs ...` ausweichen — das funktioniert, umgeht aber die
M6-Konvention, dass der Compose-Adapter alle Compose-Calls
kanalisiert.

Spec-Wortlaut (knapp):

```bash
u-boot logs                  # alle aktiven Services
u-boot logs postgres         # einzelner Service
```

Pflicht-Flags: `--follow` (fortlaufend), `--tail <n>` (letzte n
Zeilen). Keine weiteren AKs im Spec-Text.

Roadmap-Notiz: „Gehört zusammen mit dem Dry-Run-/JSON-Slice"
([`roadmap.md`](../in-progress/roadmap.md) §v0.4.0). Bedeutet:
beide V1-CLI-Erweiterungen werden in v0.4.0 gebündelt; harte
Code-Abhängigkeit besteht nicht — `--json`-Mode kommt im
Folge-Slice `slice-v1-cli-json-dry-run` nachträglich auf `logs`
drauf.

## Aufhebungsbedingung

`u-boot logs [service] [--follow] [--tail <n>]` in einem
initialisierten Projekt:

- Ohne `service`-Argument streamt es Logs der in u-boot.yaml
  registrierten Services. **Welche genau** (alle Services aus
  compose.yaml vs. nur `cfg.Services` mit `enabled: true` via
  `activeServiceNames(cfg)`) ist **T0-Sub-Decision (a)** — der
  Default wird vor T1-Start festgezurrt und im §T0-Outcomes-
  Block ergänzt (Plan-Followup-P4-Konsistenz: dieser Slice
  spricht heute an zwei Stellen widersprüchlich darüber).
- Mit `service`-Argument streamt nur diesen einen.
- `--follow` blockiert bis Ctrl-C; SIGINT beendet sauber mit
  Exit-Code 0 (analog `tail -f`-Konvention).
- `--tail <n>` akzeptiert `n ≥ 0` (Ganzzahl als String) **oder**
  die Compose-Konstante `"all"`. Default ohne Flag: Compose-
  Default (effektiv „all"). Negative oder nicht-numerische Werte
  außer `"all"` → CLI-Usage-Error, Exit-Code 2 (Cobra-/Stage-1-
  Validation in `runLogs`, vor Use-Case-Aufruf — analog der
  bestehenden `--timeout`/`--tail`-Parse-Konventionen aus M6).

Compose-/Docker-Failures klassifizieren strikt analog M6
`up`/`down` (vgl. `internal/adapter/driving/cli/cli.go:237 ff.`
`ExitCode`-Mapping):

- **Exit-Code 11** — Docker-Environment-Fehler (Docker nicht
  erreichbar / nicht installiert); `driven.ErrDockerUnavailable`
  aus dem Compose-Adapter.
- **Exit-Code 12** — Compose-Runtime-Fehler (Compose-Stack nicht
  gestartet, unbekannter Service zur Laufzeit, Compose-Exit ≠ 0);
  `driven.ErrComposeRuntime` aus dem Adapter.
- **Exit-Code 10** — User-Validation: ungültiger Service-Name
  (Format), fehlendes u-boot.yaml/compose.yaml (Project-State-
  Check, via `ErrProjectNotInitialized` /
  `ErrComposeFileMissing`).
- **Exit-Code 14** — technischer Persistenz-/FS-Fehler (z. B.
  Lesefehler auf compose.yaml während des Project-State-Checks);
  selten erreichbar bei `logs`, aber für Symmetrie mit
  `up`/`down` erhalten.
- **Exit-Code 0** — `--follow` durch SIGINT (siehe SIGINT-Vertrag
  unten).

## Akzeptanzkriterien

- ✅ **Driven-Port-Erweiterung:**
  `driven.DockerEngine.ComposeLogs(ctx, dir, opts)` ergänzt;
  `opts` trägt `Sink io.Writer` für stdout-Streaming (analog
  `ProgressSink` aus `ComposeUpOptions`), `Services []string`
  (leer = Default gemäß T0-(a)), `Follow bool`, `Tail string`
  (Compose-Konvention: `"all"` oder Ganzzahl-String).
  Adapter shellt zu `docker compose -f <dir>/compose.yaml
  logs ...` aus. **Adapter-Kontrakt für SIGINT
  (Plan-Followup-P3):** wenn `ctx.Err() != nil` nach
  `cmd.Run()`, gibt der Adapter `ctx.Err()` (also
  `context.Canceled` bzw. `context.DeadlineExceeded`)
  **unverdeckt** zurück — **nicht** in
  `driven.ErrComposeRuntime` wrappen. Sonst wird Ctrl-C zu
  Exit-Code 12 statt 0. `exec.CommandContext` killt den
  Compose-Prozess; der Wrap-Filter sitzt im Adapter direkt am
  `cmd.Run()`-Returnpunkt.
- ✅ **Driving-Port (Plan-Followup-P2):** Neuer `LogsUseCase`
  mit `LogsRequest{BaseDir, Service, Follow, Tail, OutputSink io.Writer}`
  (Sink im Request analog `UpRequest.ProgressSink` /
  `DownRequest.ProgressSink` aus M6 — die CLI gibt `cmd.OutOrStdout()`
  rein, der Use-Case reicht weiter an `ComposeLogsOptions.Sink`).
  `LogsResponse` leer (Output ist Stream, keine strukturierte
  Rückgabe). Application-Service `LogsService` orchestriert
  Project-State-Check (u-boot.yaml + compose.yaml vorhanden,
  analog `UpService` §M6-T1) plus Service-Name-Validation
  (siehe T0-Sub-Decision (b)).
- ✅ **CLI:** `u-boot logs [service]` mit optionalem Positional-
  Arg (Cobra-`MaximumNArgs(1)`). Flags `--follow` / `--tail
  <n>`. Service-Name-Validation via `domain.NewServiceName`
  (Exit-Code 10 bei Format-Fehler, mappt durch
  `isServiceValidationError` analog `add`/`remove`). `--tail`-
  Parse: Strings `"all"`, `"0"`, `"1"`, …, `"100"` akzeptiert;
  negative oder andere non-numerische Werte → Exit-Code 2 via
  CLI-Stage-1-Parse-Fehler (vor Use-Case-Aufruf).
- ✅ **Streaming-Disziplin:** Adapter line-buffert auf den Sink,
  damit `--follow` real-time ankommt (Compose-Default kann
  block-buffern bei pipe-stdout). Tradeoff dokumentiert; ggf.
  via `docker compose logs --no-log-prefix` / `--timestamps`
  Subentscheidung in T0-(d).
- ✅ **SIGINT-Vertrag:** Ctrl-C im `--follow`-Pfad beendet mit
  Exit-Code 0 (tail-konform). Der Vertrag besteht aus drei
  Schichten:
  1. **Adapter:** gibt `ctx.Err()` unverdeckt zurück (siehe
     Driven-Port-Erweiterung oben).
  2. **Use-Case:** prüft `errors.Is(err, context.Canceled)` und
     `errors.Is(err, context.DeadlineExceeded)`; in beiden
     Fällen Rückgabe `(LogsResponse{}, nil)` — kein Fehler.
  3. **CLI:** `cmd.Context()` mit `signal.NotifyContext(ctx,
     os.Interrupt)` wired (analog vermutlich `up --follow`
     falls schon vorhanden — siehe `internal/adapter/driving/cli/up.go`,
     ansonsten neu in `logs.go`).
- ✅ **Tests:**
  - Application-Unit-Tests (fakeDockerEngine, mock-ComposeLogs)
    pinnen: Request→Adapter-Call-Mapping (OutputSink wird
    durchgereicht), Service-Validation-Fehlerpfad,
    Project-State-Check (kein u-boot.yaml → Exit-Code 10),
    SIGINT-Pass-Through (`context.Canceled` aus Adapter →
    Use-Case-`nil`-Return).
  - Adapter-Unit-Test mit fake-cmd-runner: Konstruktion von
    `docker compose logs`-Argumenten je nach Flag-Kombination,
    Sink-Streaming, `ctx.Err()`-Pass-Through am `cmd.Run()`-
    Returnpunkt.
  - CLI-Test analog `cli_test.go`-Pattern: Flag-Parsing,
    `--tail`-Validierungs-Failures (Exit-2-Pin), Mapping aller
    vier Use-Case-Sentinels auf Exit-Codes 10/11/12/14.
  - **Docker-tag E2E** in `internal/e2e/` (analog
    `up_acceptance_docker_test.go`): postgres hochfahren, `u-boot
    logs postgres --tail 5` zeigt mindestens eine Log-Zeile;
    `u-boot logs --follow` startet und wird via Test-Timeout +
    Context-Cancellation beendet (Exit 0).
- ✅ **Spec-Pin:** `internal/hexagon/application/acceptance_test.go`
  oder Docker-e2e-Test deckt `LH-FA-UP-005` ab; Test-Naming
  `TestLHFAUP005_Logs<…>` analog `TestLHFADEV003_*`.
- ✅ **Doku:** README (EN + DE) Quickstart-Block oder Subcommand-
  Tabelle um `u-boot logs` ergänzt; ggf. neue
  `docs/user/logs.md` falls Verhalten Detail-Erklärung verdient
  (vermutlich nicht — `--help`-Output reicht für simple Flags).

## Tranchen (Skizze, wird beim Übergang nach `next/` verfeinert)

| T   | Inhalt (Skizze) | LOC (Schätzung) |
| --- | --------------- | --------------- |
| T0  | **Discovery / Design.** Vier Sub-Decisions (Plan-Followup-P4: (a) ist **blockierend** für die §Aufhebungsbedingung — der Plan-Text widerspricht sich heute zwischen „nur enabled" und „leer = alle"; vor T1-Start eindeutig auflösen): (a) leerer Service-Filter → Compose-Default vs. `activeServiceNames`-Filter (entscheidet, ob deaktivierte/manuell-Compose-Services Logs leaken können — Source-of-Truth-Frage); (b) Service-Name-Validation-Tiefe (`domain.NewServiceName` + Katalog-Membership-Check wie `add`, oder nur Regex + Pass-Through zu Compose?); (c) `--tail`-Default ohne Flag (Compose-Default ist „all"; übernehmen oder `tail=100`?); (d) Output-Format-Sub-Entscheidung (`--no-log-prefix` per Default, `--timestamps` opt-in?). Ergebnis als §T0-Outcomes im Plan analog `slice-v1-devcontainer-features`. | — (Plan-Arbeit) |
| T1  | **Driven-Port + Adapter.** `ComposeLogs`-Methode + `ComposeLogsOptions`-Struct in `port/driven/docker_engine.go`. Adapter-Implementation in `adapter/driven/docker/engine.go` mit `exec.CommandContext` (Context-Cancellation für SIGINT). Unit-Tests für Adapter (mock-out via fake-cmd-runner). | ~120 |
| T2  | **Use-Case.** `LogsRequest`/`LogsResponse`/`LogsUseCase` in `port/driving/logsservice.go`; `LogsService` in `application/logsservice.go`. Project-State-Check (analog `UpService`); Service-Name-Validation gemäß T0-(b). | ~120 |
| T3  | **CLI-Subcommand.** `internal/adapter/driving/cli/logs.go`; Cobra-Command `logs [service]` mit `--follow`/`--tail`-Flags; SIGINT-Handler oder Context-Cancellation via `cmd.Context()`. App-Wiring in `cli.go:New`. | ~80 |
| T4  | **Docker-Tag E2E + Spec-Pin.** `internal/e2e/logs_acceptance_docker_test.go`: postgres up + logs --tail + logs --follow (mit Test-Timeout-Cancellation). Plus Application-Layer-Acceptance-Test `TestLHFAUP005_Logs*`. | ~150 |
| T5  | **Doku + Closure.** README EN+DE Subcommand-Tabelle, CHANGELOG `## [Unreleased]`-Eintrag. Slice `open/` → `done/` mit DoD-Hash-Line. Roadmap-Status. | — (Doku) |

LOC-Summe T1-T4 ≈ **470 LOC** (Schätzung); unter 800-LOC-Carveout-
Schwelle. Re-Check vor T4-Start.

## Out of Scope

- **`--since <timestamp>` / `--until <timestamp>`-Flags:** Spec
  verlangt nur `--follow` und `--tail`. Compose-CLI hat beides,
  aber bis Trigger-Nachfrage YAGNI.
- **Multi-Service-Filter** (`u-boot logs postgres keycloak`):
  Spec-Beispiel hat nur einen Service-Filter, Cobra-`MaximumNArgs(1)`
  wahrt das. Spätere Erweiterung auf `MaximumNArgs(N)` falls
  Trigger.
- **`--json`-Output:** lebt im Folge-Slice
  `slice-v1-cli-json-dry-run` (Roadmap-Notiz); diese Slice
  verzichtet auf Output-Format-Switches, damit der
  Maschinen-Schnittstellen-Slice einen klaren Greenfield-Punkt
  hat.
- **`u-boot exec <service> <cmd>`:** driven-README listet `Exec`
  als V1-Geschwister-Erweiterung, aber Spec-ID ist nicht
  zugeordnet. Eigener Folge-Slice falls Trigger.

## Bezug

- Spec: `LH-FA-UP-005`
  ([`spec/lastenheft.md:1023`](../../../../spec/lastenheft.md)).
- Port-Anker:
  [`internal/hexagon/port/driving/README.md:39`](../../../../internal/hexagon/port/driving/README.md)
  (`LogsUseCase`),
  [`internal/hexagon/port/driven/README.md:60`](../../../../internal/hexagon/port/driven/README.md)
  (`Logs` als Geplante DockerEngine-Erweiterung).
- Vorbild-Slices:
  [`slice-m6-up-down`](../done/slice-m6-up-down.md)
  (Compose-Adapter-Pattern + Project-State-Check + Cobra-
  Subcommand-Wiring),
  [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md)
  T0-Outcomes-Layout als Doku-Template für die T0-Sub-Decisions.
- Code-Anker:
  `internal/adapter/driven/docker/engine.go:81` (ComposeUp-
  Pattern für CommandContext + ProgressSink-Streaming),
  `internal/hexagon/application/downservice.go:50` (Use-Case-
  Skelett mit Project-State-Check),
  `internal/adapter/driving/cli/down.go` (Cobra-Subcommand-
  Pattern mit Context-Propagation).
- Roadmap: [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0 —
  Bündelung mit `slice-v1-cli-json-dry-run`.
- Phase: V1, geplant für v0.4.0-Release.
