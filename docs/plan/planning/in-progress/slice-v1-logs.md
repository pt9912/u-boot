# Slice V1: `u-boot logs` (`LH-FA-UP-005`)

> **Status:** geplant für v0.4.0 — Spec ✅
> ([`spec/lastenheft.md:1023-1040`](../../../../spec/lastenheft.md)),
> Port-Anker ✅
> ([`internal/hexagon/port/driving/README.md:39`](../../../../internal/hexagon/port/driving/README.md)
> nennt `LogsUseCase` als V1-Erweiterung, driven/README §"Geplante
> Erweiterungen" listet `Logs`/`Exec`-Verb auf `DockerEngine`),
> Plan-Followup P1..P5 ✅ (Review-Findings adressiert),
> T0-Discovery ✅ (siehe §T0-Outcomes). In `in-progress/`: T1
> (Driven-Port + Adapter) läuft, T2..T5 ausstehend.

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
([`roadmap.md`](roadmap.md) §v0.4.0). Bedeutet:
beide V1-CLI-Erweiterungen werden in v0.4.0 gebündelt; harte
Code-Abhängigkeit besteht nicht — `--json`-Mode kommt im
Folge-Slice `slice-v1-cli-json-dry-run` nachträglich auf `logs`
drauf.

## Aufhebungsbedingung

`u-boot logs [service] [--follow] [--tail <n>]` in einem
initialisierten Projekt:

- Ohne `service`-Argument streamt es Logs **aller Services aus
  dem Compose-Projekt** (`compose.yaml`) gemäß Compose-Default
  — kein u-boot.yaml-Filter, kein `activeServiceNames(cfg)`-
  Filter. Konkret: `ComposeLogsOptions.Services` bleibt leer,
  Compose entscheidet. Begründung in §T0-Outcomes (a):
  `u-boot logs` ist Compose-Facade, kein State-Machine-Tool;
  manuell ergänzte Compose-Services bleiben sichtbar.
- Mit `service`-Argument streamt nur diesen einen.
- `--follow` blockiert bis Ctrl-C; SIGINT beendet sauber mit
  Exit-Code 0 (analog `tail -f`-Konvention).
- `--tail <n>` akzeptiert ausschließlich Ganzzahlen `n ≥ 0`
  (Compose-Konvention für Numeric-Tail). Default ohne Flag: CLI
  setzt leeren String, Use-Case normalisiert zu `"all"`; Adapter
  übersetzt das in `docker compose logs --tail all`.
  Negative Werte oder nicht-numerische Strings → CLI-Usage-
  Error, **Exit-Code 2** (Cobra-/Stage-1-Validation in `runLogs`,
  vor Use-Case-Aufruf). Der String `"all"` selbst ist **nicht**
  CLI-User-Input — er entsteht nur durch die Use-Case-
  Normalisierung; Tests pinnen alle drei Pfade (kein Flag,
  `--tail 0`, negativ/non-numerisch).

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
  <n>`. Service-Name-Validation **ausschließlich** via
  `domain.NewServiceName` (Format-Regex; Exit-Code 10 bei
  Format-Fehler, mappt durch `isServiceValidationError`
  analog `add`/`remove`) — **keine** `cfg.Services`- oder
  Katalog-Membership-Prüfung (Compose-Delegation, T0-(b)).
  Unbekannter Service zur Laufzeit landet via
  `driven.ErrComposeRuntime` bei **Exit-Code 12**; das ist
  Plan-Akzeptanz (T0-(b)-Folge), nicht ein Fehler — der User
  sieht eine Compose-Runtime-Meldung wie bei `up`/`down`.
  `--tail`-Parse: Ganzzahlen `n ≥ 0` akzeptiert (`"0"`, `"1"`,
  …, `"100"`); negative oder non-numerische Werte → Exit-Code
  2 via CLI-Stage-1-Parse-Fehler (vor Use-Case-Aufruf). Der
  String `"all"` selbst ist intern (Use-Case-Normalisierung),
  nicht User-Input.
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
| T0  | **Discovery / Design.** ✅ Vier Sub-Decisions festgezurrt — siehe [§T0-Outcomes](#t0-outcomes). (a) Compose-Default ohne `activeServiceNames`-Filter; (b) nur Regex-Validation, Compose macht Existenz-Check (unbekannter Service → Exit-12); (c) `--tail`-Default leerer String → Use-Case-Normalisierung zu `"all"`; (d) Spec-treu — nur `--follow`+`--tail`, keine Output-Format-Flags. | — (Plan-Arbeit) |
| T1  | **Driven-Port + Adapter.** ✅ Done. `ComposeLogsOptions{Services, Follow, Tail, Sink}` + `ComposeLogs(ctx, dir, opts)` in `port/driven/docker_engine.go`. Adapter in `adapter/driven/docker/engine.go` mit `exec.CommandContext`; **zweistufiger SIGINT-Pass-Through (P3-Vertrag):** (1) `ctx.Err()`-Pre-Preflight-Check, (2) `ctx.Err()`-Post-`cmd.Run()`-Check. Beide returnen `ctx.Err()` unverdeckt, damit Ctrl-C nicht in ErrComposeRuntime/ErrDockerUnavailable maskiert. 3 Adapter-Tests (missing-binary→ErrDockerUnavailable, happy-path-stream-to-sink, ctx-canceled-pre-call→context.Canceled). `fakeDockerEngine` um `ComposeLogs`-Stub erweitert (T2 ergänzt die scripting-Helper). | ~120 geschätzt / **~102 real** (−15 %, unter Budget) |
| T2  | **Use-Case.** ✅ Done. `LogsRequest{BaseDir, Service, Follow, Tail, OutputSink}` + `LogsResponse{}` + `LogsUseCase` in `port/driving/logs.go` (mit ausführlichen T0-Outcomes-Referenzen im Doc-Kommentar). `LogsService` in `application/logsservice.go`: BaseDir-Check → Project-State-Check (u-boot.yaml + compose.yaml) → Tail-Normalisierung (T0-(c): leer → `"all"`) → ComposeLogs-Aufruf → SIGINT-Pass-Through (`context.Canceled` und `context.DeadlineExceeded` → `(LogsResponse{}, nil)`). Service-Name-Validation NICHT im Use-Case (T0-(b): nur Regex auf CLI-Ebene, Compose macht Existenz-Check). 8 Tests (BaseDir-empty, ohne u-boot.yaml, ohne compose.yaml, Happy-Path-Tail-Normalisierung, Happy-Path-Service-Filter, SIGINT-Canceled, SIGINT-Deadline, ErrComposeRuntime-Propagation, ErrDockerUnavailable-Propagation). | ~120 geschätzt / **~245 real** (+104 %; getrieben durch ausführliche Doc-Kommentare mit T0/P3-Referenzen — analog Parent-Slice T4-Verlauf) |
| T3  | **CLI-Subcommand.** `internal/adapter/driving/cli/logs.go`; Cobra-Command `logs [service]` mit `--follow`/`--tail`-Flags; SIGINT-Handler oder Context-Cancellation via `cmd.Context()`. App-Wiring in `cli.go:New`. | ~80 |
| T4  | **Docker-Tag E2E + Spec-Pin.** `internal/e2e/logs_acceptance_docker_test.go`: postgres up + logs --tail + logs --follow (mit Test-Timeout-Cancellation). Plus Application-Layer-Acceptance-Test `TestLHFAUP005_Logs*`. | ~150 |
| T5  | **Doku + Closure.** README EN+DE Subcommand-Tabelle, CHANGELOG `## [Unreleased]`-Eintrag. Slice `open/` → `done/` mit DoD-Hash-Line. Roadmap-Status. | — (Doku) |

LOC-Summe T1-T4 ≈ **470 LOC** (Schätzung); unter 800-LOC-Carveout-
Schwelle. Re-Check vor T4-Start.

## T0-Outcomes

Vier Sub-Decisions vor T1-Start festgezurrt. Plan-Followup-P4
hatte (a) als blockierend markiert; alle vier sind hier
verbindlich.

### T0-(a) Service-Filter: Compose-Default (alle)

**Entscheidung:** `u-boot logs` ohne Argument lässt
`ComposeLogsOptions.Services` leer; Compose entscheidet, welche
Services geloggt werden (= alle in `compose.yaml`, unabhängig
von `cfg.Services.<name>.enabled`).

**Begründung:** `u-boot logs` ist operativ Compose-Facade, kein
State-Machine-Tool. `up`/`down` lesen u-boot.yaml nur als
Projektmarker und führen Compose gegen `compose.yaml` aus; ein
`activeServiceNames(cfg)`-Filter ausgerechnet beim Inspect-
Befehl würde eine andere Runtime-Sicht einführen. Manuell
ergänzte Compose-Services bleiben für `logs` sichtbar — das
ist konsistent mit der Compose-Delegations-Linie aus M6.

### T0-(b) Service-Name-Validation: nur Regex, Compose macht Existenz

**Entscheidung:** Service-Name-Validation ausschließlich via
`domain.NewServiceName` (Format-Regex). Keine
`cfg.Services`-Membership-Prüfung, keine Katalog-Prüfung.
Unbekannte Services zur Laufzeit landen via
`driven.ErrComposeRuntime` bei Exit-Code 12.

**Begründung:** Konsistent mit T0-(a). Eine cfg- oder Katalog-
Prüfung würde `u-boot logs` enger als `docker compose logs`
machen und manuell ergänzte Services ausschließen — Bruch der
Compose-Facade-Semantik. **Akzeptierte Folge:** Ein Tippfehler
(`u-boot logs psotgres`) wird nicht früh durch eine Validation
gefangen, sondern landet bei Exit-Code 12 mit Compose-Runtime-
Meldung. Das ist Plan-Akzeptanz, kein Bug. Format-Fehler (z. B.
`u-boot logs Postgres` mit Großbuchstaben) bleiben bei
Exit-Code 10 via `domain.ErrInvalidServiceName`.

### T0-(c) `--tail`-Default: Compose-Default `"all"`

**Entscheidung:** CLI defaultet `--tail` auf leeren String;
Use-Case normalisiert leer → `"all"`. Adapter setzt
`--tail all` an die Compose-CLI weiter. Akzeptierter
User-Input-Range: Ganzzahl-Strings `"0"`, `"1"`, …
(keine Obergrenze; Compose mapped sehr große Werte selbst).
Negative oder non-numerische Werte (außer dem internen `"all"`)
→ Exit-Code 2 in der CLI-Stage-1-Validation.

**Begründung:** `u-boot logs` ist Compose-Facade; principle-of-
least-surprise. Defensive Defaults wie `tail=100` (kubectl-
Konvention) wären nur sinnvoll, wenn die Spec Performance gegen
lange Historien priorisierte — tut sie nicht. Performance-
Probleme bei langlaufenden Containern löst der User mit
explizitem `--tail 100`.

### T0-(d) Output-Format: Spec-treu, keine zusätzlichen Flags

**Entscheidung:** CLI exponiert nur die zwei Spec-Pflicht-
Flags `--follow` und `--tail`. Adapter gibt Compose-Output
unverändert an stdout. Kein `--no-log-prefix`, kein
`--timestamps`. Format-Flag-Wünsche sind explizit Out of Scope
mit benanntem Trigger (siehe §Out of Scope).

**Begründung:** Hält den Surface klein und vermeidet, dass der
spätere `slice-v1-cli-json-dry-run`-Slice direkt noch
Format-Sonderfälle mitziehen muss. Spec-treu (LH-FA-UP-005
fordert genau zwei Flags). Compose-Output-Default (mit
Service-Prefix, ohne Timestamps) bleibt erhalten — der User
sieht das gleiche Layout wie bei `docker compose logs`.

---

## Out of Scope

- **`--no-log-prefix` / `--timestamps`-Flags (Output-Format,
  T0-(d)):** Spec-Pflicht sind nur `--follow` und `--tail`;
  dieser Slice exponiert genau diese. Compose-Defaults bleiben
  (Service-Prefix sichtbar, keine Timestamps). Trigger für
  einen Folge-Slice: konkrete Debugging-/Logformat-Nachfrage
  (z. B. „Timestamps für Incident-Forensik fehlen"). Bewusst
  *vor* dem `slice-v1-cli-json-dry-run`-Slice, weil
  `--json` einen eigenen Output-Modus mitbringt — Format-
  Flags hier würden den später kommenden Maschinen-Format-
  Slice unnötig verschachteln.
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
- Roadmap: [`roadmap.md`](roadmap.md) §v0.4.0 —
  Bündelung mit `slice-v1-cli-json-dry-run`.
- Phase: V1, geplant für v0.4.0-Release.
