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

- Ohne `service`-Argument streamt es Logs aller in
  `cfg.Services` mit `enabled: true` registrierten Services
  (Compose-Default ist „alle"; der Use-Case nutzt
  `activeServiceNames(cfg)` als Filter falls nötig — siehe
  T0-Decision unten).
- Mit `service`-Argument streamt nur diesen einen.
- `--follow` blockiert bis Ctrl-C; SIGINT beendet sauber mit
  Exit-Code 0 (analog `tail -f`-Konvention).
- `--tail <n>` mit `n ≥ 0` zeigt nur die letzten n Zeilen pro
  Service.

Compose-CLI-Failures (Service unbekannt, Project nicht gestartet,
Docker nicht erreichbar) klassifizieren wie M6 `up`/`down`:
Exit-Codes 10 (User-Fehler) bzw. 14 (technisch).

## Akzeptanzkriterien

- ✅ **Driven-Port-Erweiterung:**
  `driven.DockerEngine.ComposeLogs(ctx, dir, opts)` ergänzt;
  `opts` trägt Sink (`io.Writer` für stdout-Streaming, analog
  `ProgressSink` aus `ComposeUpOptions`), `Services []string`
  (leer = alle), `Follow bool`, `Tail string` (Compose-Konvention:
  „all" oder Dezimalzahl-String). Adapter shellt zu
  `docker compose -f <dir>/compose.yaml logs ...` aus. Context-
  Cancellation propagiert via `exec.CommandContext` — `--follow`
  bleibt unterbrechbar.
- ✅ **Driving-Port:** Neuer `LogsUseCase` mit
  `LogsRequest{BaseDir, Service, Follow, Tail}` und
  `LogsResponse` (vermutlich leer — Output ist Stream auf Writer).
  Application-Service `LogsService` orchestriert Project-State-
  Check (u-boot.yaml + compose.yaml vorhanden, analog
  `UpService` §M6-T1) plus Service-Name-Validation (siehe T0
  unten).
- ✅ **CLI:** `u-boot logs [service]` mit optionalem Positional-
  Arg (Cobra-`MaximumNArgs(1)`). Flags `--follow` / `--tail
  <n>`. Service-Name-Validation via `domain.NewServiceName`
  (Exit-Code 10 bei Format-Fehler, mappt durch
  `isServiceValidationError` analog `add`/`remove`).
- ✅ **Streaming-Disziplin:** Adapter line-buffert auf den Sink,
  damit `--follow` real-time ankommt (Compose-Default kann
  block-buffern bei pipe-stdout). Tradeoff dokumentiert; ggf.
  via `docker compose logs --no-log-prefix` / `--timestamps`
  Subentscheidung in T0.
- ✅ **SIGINT-Vertrag:** Ctrl-C im `--follow`-Pfad beendet mit
  Exit-Code 0 (tail-konform); innerhalb der Use-Case-Schicht
  als `context.Canceled` erkannt und nicht als Fehler propagiert.
- ✅ **Tests:**
  - Application-Unit-Tests (fakeDockerEngine, mock-ComposeLogs)
    pinnen: Request→Adapter-Call-Mapping, Service-Validation-
    Fehlerpfad, Project-State-Check (kein u-boot.yaml → Exit-
    Code 10), SIGINT-Pass-Through.
  - Adapter-Unit-Test für die Cobra-CLI-Wiring (analog
    `cli_test.go`-Pattern).
  - **Docker-tag E2E** in `internal/e2e/` (analog
    `up_acceptance_docker_test.go`): postgres hochfahren, `u-boot
    logs postgres --tail 5` zeigt mindestens eine Log-Zeile;
    `u-boot logs --follow` startet und wird via Test-Timeout
    beendet.
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
| T0  | **Discovery / Design.** Vier Sub-Decisions: (a) leerer Service-Filter → Compose-Default vs. `activeServiceNames`-Filter (M6-Konvention; macht keinen Unterschied für `logs`, aber Konsistenz); (b) Service-Name-Validation-Tiefe (`domain.NewServiceName` + Katalog-Membership-Check wie `add`, oder nur Regex + Pass-Through zu Compose?); (c) `--tail`-Default ohne Flag (Compose-Default ist „all"; übernehmen oder `tail=100`?); (d) Output-Format-Sub-Entscheidung (`--no-log-prefix` per Default, `--timestamps` opt-in?). Ergebnis als §T0-Outcomes im Plan analog `slice-v1-devcontainer-features`. | — (Plan-Arbeit) |
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
