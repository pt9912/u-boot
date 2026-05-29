# Slice M6: `u-boot up` / `u-boot down`-Flow

> **Status:** In progress
> **DoD:** T1 ✅ `9f8badd` / T2 ✅ `84a676c` / T3 ⏳ / T4 ⏳ / T5 ⏳ / T6 ⏳ / T7 ⏳

## Auslöser

Nach M3 (`u-boot init`), M4 (`u-boot doctor`) und M5 (`u-boot add postgres`)
fehlen die letzten beiden MVP-Subkommandos, um die Compose-Umgebung
tatsächlich zu betreiben: `u-boot up` und `u-boot down`. Erst damit
schließt sich der MVP-Acceptance-Flow `LH-AK-002` (`u-boot init &&
u-boot add postgres && u-boot up`).

Spec-Pflicht für M6 (alle MVP-Priorität, `spec/lastenheft.md` §4.6):

- **`LH-FA-UP-001`** Umgebung starten. Standard-Wartezeit 60 s auf
  Stabilisierung; `--timeout <sek>` überschreibt; `--timeout=0`
  bedeutet fire-and-forget (kein Warten). Negative Werte ⇒ Exit-Code 2.
  Stabilisierungsdefinition:
  - Healthcheck-Service: Zielzustand `healthy`.
  - Service ohne Healthcheck: Zielzustand `running` ausreicht.
  - Ports: TCP-Probe auf `localhost`; nicht-TCP / nicht eindeutig
    probebar ⇒ `warn`-Diagnose statt Abbruch.
- **`LH-FA-UP-002`** Intern Docker Compose verwenden.
- **`LH-FA-UP-003`** Startstatus anzeigen — Mindestangaben:
  Service-Name, Containerstatus, Port, Healthcheck-Status (falls
  vorhanden).
- **`LH-FA-UP-004`** Umgebung stoppen; `--volumes` für vollständiges
  Aufräumen (Container + Volumes).

Plus aus angrenzenden Spec-Punkten:

- **`LH-FA-CLI-005A`** Destruktive Bestätigung für `down --volumes`
  (`--yes` oder interaktive Bestätigung; `--no-interactive` ohne
  `--yes` ⇒ Exit-Code 10).
- **`LH-FA-CLI-006`** Exit-Code-Mapping (vollständige
  Sentinel→Code-Tabelle in T7; hier nur die Übersicht):
  - `2` — CLI-Validierung (`--timeout=-1`, gleichzeitiges `--yes`
    + `--no-interactive`, unbekannte Flags).
  - `10` — fachlicher Validierungsfehler (kein `u-boot.yaml`, kein
    `compose.yaml`, fehlende Bestätigung bei `down --volumes`).
  - `11` — Umgebungsfehler. Wird **vor** dem eigentlichen
    Compose-Call durch Adapter-Pre-Probes erkannt (T2:
    `docker version`-Roundtrip + `docker compose version`-
    Roundtrip), nicht durch Stderr-Parse des `up`/`down`-Calls.
    Sentinel: `ErrDockerUnavailable`.
  - `12` — Ausführungsfehler **nach** bestandenen Pre-Probes:
    Compose-Start-Failure (`ErrComposeRuntime`),
    Stabilisierungs-Timeout (`ErrStabilizationTimeout`).
- **`LH-NFA-PERF-002`** Fortschritt einzelner Services (Pull/Create/
  Start/Healthcheck) sichtbar darstellen.

Out of Scope (V1):

- **`LH-FA-UP-005`** `u-boot logs` / `u-boot logs <service>` — eigener
  V1-Slice (eigene Flags `--follow`, `--tail`).
- **`LH-NFA-USE-004`** `--json`-Output für `up`/`down` — bewusst V1,
  analog M4-doctor-Entscheidung. Text-Output zuerst stabilisieren.

## Vorbereitende Slices (Status)

- [`slice-m4-soft-existing-detection`](../done/slice-m4-soft-existing-detection.md)
  — `Confirmer`-Port liegt vor; `down --volumes` reuse't ihn.
- [`slice-m4-logging-port`](../done/slice-m4-logging-port.md) —
  `Logger`-Port; `up`/`down` emittieren Debug-/Info-Events pro Phase.
- [`slice-m6-docker-integrationstests`](../open/slice-m6-docker-integrationstests.md)
  — `//go:build docker`-CI-Pfad; M6 ist der erste Adapter, für den
  reale Compose-Calls existieren — dieser Carveout-Slice wird mit T7
  von M6 fachlich freigeschaltet, bleibt aber als separater Slice
  bestehen (Aufweichungs-Auflöser, kein Tranche-Folie).

## Architektur-Punkte

- **Neuer Driven-Port `DockerEngine`** (state-mutierend), getrennt
  vom existierenden read-only `DockerProbe` (siehe Kommentar in
  `internal/hexagon/port/driven/docker.go:23`). Methoden
  (verbindliche Signaturen — T2-Detail-Sektion ist die kanonische
  Quelle; Port-Typen sind **bewusst nicht** Domain-Typen, weil
  die Application zwischen Port- und Domain-Ebene mappt):
  `ComposeUp(ctx, dir, opts ComposeUpOptions) (ComposeUpResult, error)`
  (umfasst Pull + Create + Start),
  `ComposeDown(ctx, dir, opts ComposeDownOptions) error`,
  `ComposePs(ctx, dir) ([]ComposeService, error)`.
  Adapter in `internal/adapter/driven/docker/engine.go` mit
  `os/exec docker compose`.
  Eventuell `ComposePull` separat herausziehen, wenn Fortschritts-
  Streaming sonst zu eng wird; T2-Entscheidung.

  **Typ-Schichtung (verbindlich, sonst entsteht Drift zwischen
  Architektur-Übersicht und Tranchen-Detail):**
  - Driven-Port-Typen in `internal/hexagon/port/driven/`:
    `ComposeUpOptions`, `ComposeDownOptions`, `ComposeUpResult`,
    `ComposeService` — rohe Compose-Beobachtungen, kein Domain-
    Wissen.
  - Domain-Typen in `internal/hexagon/domain/serviceup.go`:
    `UpResult`, `ServiceStatus`, `containerState`,
    `StabilizationOutcome` — u-boot-eigene Klassifikation.
  - Driving-Port-Typen in `internal/hexagon/port/driving/`:
    `UpRequest`, `UpResponse`, `DownRequest`, `DownResponse` —
    Application-Schnittstelle.
  - Mapping `ComposeService → ServiceStatus` und
    `ComposeUpResult → UpResult` lebt **ausschließlich** in der
    Application (`upservice.go`); Adapter darf keine Domain-Typen
    importieren (depguard-Regel `adapter-no-domain`, falls noch
    nicht aktiv, in T2 ergänzen).
- **Neuer Driven-Port `NetProbe`** für TCP-Reachability:
  `DialTCP(ctx, host string, port int, timeout time.Duration) error`.
  Adapter in `internal/adapter/driven/net/probe.go` mit
  `net.DialTimeout`. Begründung für eigenen Port: hält den
  Application-Service frei von `net`-Importen (depguard-Regel
  `application-no-net`/`application-no-stdlib-io` — siehe carveouts);
  Fake im Test injizierbar ohne realen TCP-Versuch.
- **Streaming-Output für `LH-NFA-PERF-002`:** `docker compose up`
  emittiert die Phasen (Pull/Create/Start/Healthcheck) bereits auf
  stderr. Adapter sollte die Streams durchreichen, ohne sie zu
  puffern; Detail-Entscheidung (Channel vs. Writer-Injection) in T2.
- **Confirmer-Reuse:** `down --volumes` ruft den bereits in M4
  etablierten `Confirmer`-Port (`Confirm(prompt, defaultYes) (bool,
  error)`); kein neuer Port nötig.

## State-Modell (`u-boot up`)

Pro Service normalisiert `up` zuerst den rohen Compose-State-String
in eine eigene `containerState`-Enum (in `domain/serviceup.go`), um
String-Vergleiche im Polling-Pfad auszuschließen — Compose ändert
Casing und Werte zwischen Releases (`Running` vs. `running`,
`restarting` vs. nicht dokumentiert), Freitext-Vergleich wäre
fragil.

**Compose-State → `containerState`-Mapping** (Quelle: `docker
compose ps --format json`, Feld `State`; case-insensitive
gematched):

| Compose-Raw-State                              | `containerState`   | Polling-Konsequenz                                  |
| ---------------------------------------------- | ------------------ | --------------------------------------------------- |
| `running`                                      | `stateRunning`     | Healthcheck/Port-Probe entscheidet (siehe unten)    |
| `restarting`                                   | `stateRestarting`  | siehe Restart-Loop-Regel unten                      |
| `created`, `starting`, `paused`                | `stateStarting`    | Polling läuft weiter (Service hat noch nicht voll gestartet) |
| `exited`, `dead`, `removing`, `removed`        | `stateDead`        | sofort `Failed` (kein Polling-Retry) — diese fünf bilden die **Dead-Allowlist**: explizit dokumentierte Compose-States, deren Bedeutung „Container existiert nicht / ist tot" eindeutig ist |
| alles andere (unbekannter Compose-State)       | `stateUnknown`     | **Polling läuft weiter** (degradiert zu `RunningOnly`) + persistente Diagnose-warn auf den Raw-String — kein sofortiges `Failed` |

**Begründung für die Soft-Behandlung von `stateUnknown`:**
Compose erweitert sein State-Vokabular zwischen Versionen
(historisch z. B. neue Werte für Pulling/Health-Übergänge); eine
hart-on-`Failed`-Logik würde u-boot bei jeder Compose-Erweiterung
spontan brechen, obwohl der Service möglicherweise korrekt läuft.
`stateUnknown` ⇒ Service zählt als `RunningOnly` (Polling läuft
weiter, Healthcheck/Port-Probe entscheiden letztendlich); bleibt
der Service nach `req.Timeout` weiter ohne Stabilisierung, fällt
das in den normalen `ErrStabilizationTimeout`-Pfad (Code 12).
Damit ist die einzige harte `Failed`-Klassifikation an die
**explizite** Dead-Allowlist (`exited`/`dead`/`removing`/`removed`)
gebunden — eine geschlossene, dokumentierte Liste, die u-boot
selbst pflegt. Trade-off bewusst: marginal längere Stabilisierungs-
zeit für Service, die sich in einem zukünftigen Compose-State
„fest setzen" — akzeptabel gegen das größere Risiko
falsch-positiver Failures bei Compose-Updates.

Pro `stateUnknown`-Beobachtung wird **genau ein**
Diagnose-Eintrag pro Service ausgegeben (Idempotenz via
`unknownStateReported`-Set in der Polling-Datenstruktur, damit
nicht jede der 120 Polls denselben Eintrag duplikat erzeugt).
Diagnose-ID: `up.state.<service>.unknown`; Message enthält den
Raw-String für User-Debugging.

**Restart-Loop-Detection (Statt Freitext „Restart-Loop"):**
Compose `ps --format json` liefert pro Container ein `Health`-
Feld und einen Service-State; einen direkten `RestartCount` gibt
es im strukturierten Output nicht. T4 zählt deshalb **selbst**
in der Polling-Loop-Datenstruktur, wie oft ein Service in
aufeinanderfolgenden Iterationen `stateRestarting` zeigt:
`restartObservations[serviceName]++` bei jedem `stateRestarting`-
Tick, Reset auf 0 sobald der Service einen anderen State zeigt.
Schwelle `restartLoopThreshold = 3` (Konstante). Erreicht ein
Service drei aufeinanderfolgende Restart-Observations ⇒ `Failed`.
Begründung: ein einzelner `restarting`-Tick kann ein normaler
Compose-internen Restart sein; eine echte Loop manifestiert sich
über mehrere Iterationen. MVP-pragmatische Heuristik;
Tuning-Möglichkeit per Flag bleibt V1-Folge-Slice.

**Klassifikation pro Service (auf normalisiertem State):**

| `containerState`                      | zusätzliche Bedingung                                                         | Klassifikation |
| ------------------------------------- | ----------------------------------------------------------------------------- | -------------- |
| `stateRunning`                        | **kein Healthcheck definiert UND kein TCP-Port deklariert** ⇒ Service zählt **allein auf Basis von `running`** als stabilisiert (LH-FA-UP-001 §967: „Für Dienste ohne Healthcheck ist `running` als Zielzustand ausreichend"). | `Stabilized` (`ok`) |
| `stateRunning`                        | Healthcheck definiert UND Healthcheck `healthy`; deklarierte Ports werden zusätzlich geprüft, aber nur als `warn`-Diagnose, nicht als Stabilisierungs-Veto (Healthcheck dominiert die Klassifikation) | `Stabilized` (`ok`) |
| `stateRunning`                        | kein Healthcheck definiert UND TCP-Port deklariert UND Port erreichbar         | `Stabilized` (`ok`) |
| `stateRunning`                        | Healthcheck definiert, aber noch nicht `healthy` (z. B. `starting`)            | `RunningOnly` (`warn`, Polling läuft) |
| `stateRunning`                        | kein Healthcheck definiert UND TCP-Port deklariert UND Port noch nicht erreichbar | `RunningOnly` (`warn`, Polling läuft) |
| `stateRunning`                        | kein Healthcheck definiert UND nur nicht-probebare Ports (UDP, Range, Long-Syntax `udp`) | `Stabilized` (`ok`) — Healthcheck-Slot ist leer, Port-Probe ist nicht möglich; per §969 ist das ein `warn`-Diagnose-Pfad, kein Stabilisierungs-Blocker |
| `stateStarting`                       | beliebig                                                                      | `RunningOnly` (Polling läuft) |
| `stateRestarting`                     | `restartObservations < restartLoopThreshold`                                  | `RunningOnly` (Polling läuft) |
| `stateRestarting`                     | `restartObservations >= restartLoopThreshold`                                 | `Failed` (sofort) |
| `stateDead`                           | beliebig                                                                      | `Failed` (sofort) |
| `stateUnknown`                        | beliebig                                                                      | `RunningOnly` (Polling läuft + persistente Diagnose-warn pro Service-Erstbeobachtung; bei Timeout normaler `ErrStabilizationTimeout`-Pfad) |

Globaler Abbruch wenn:

- ein einzelner Service in `Failed` läuft (Exit-Code 12, Compose-Start-
  Fehler), oder
- nach `--timeout` Sekunden mindestens ein Service noch nicht in
  `Stabilized` ist (Exit-Code 12, Stabilisierungs-Timeout).

**`--timeout=0` (fire-and-forget — verbindliche Sonderfall-Spez,
mit T4 Schritt 5 abgestimmt):** kein Polling, **kein**
`ComposePs`-Roundtrip, **keine** Port-/Healthcheck-Probes. `up`
returnt direkt nach erfolgreichem `ComposeUp` mit
`UpResult{Services: nil, Stabilized: false, Diagnostics:
[{ID: "up.fire-and-forget", Severity: SeverityInfo}]}`. Die
LH-FA-UP-003-Statusanzeige entfällt in diesem Modus per Spec
(LH-FA-UP-001 §970 ist die explizite Ausnahme zu §988); CLI
rendert statt der Service-Tabelle nur die Info-Diagnose-Zeile.
Dieser Pfad ist an **drei** Stellen test-gepinnt, damit
zukünftige Implementierungen nicht versehentlich auf einen
opportunistischen `ComposePs`-Call zurückfallen:
1. **T4-Application-Pin** (siehe Schritt 5): Fake-Engine mit
   `ComposePs`-Panic; Aufruf darf den Panic nicht triggern.
2. **T6-Output-Pin** (siehe T6-Tests): golden-Fixture
   `up_timeout_0.txt` enthält **keine** Service-Tabelle, nur
   die Info-Zeile aus `up.fire-and-forget`.
3. **T7-CLI-Pin** (siehe T7-Tests): `up --timeout=0` ⇒
   Fake-Engine-Assertion `PsCalls == 0`, Exit-Code 0.

Nicht-TCP / nicht eindeutig probebare Ports (UDP, Range-Syntax,
fehlende Host-IP) erzeugen einen `warn`-Diagnose-Eintrag, blockieren
aber den Stabilisierungsentscheid nicht; der Service zählt
basierend allein auf Healthcheck/Running-Status als stabilisiert.

## Tranchen-Schnitt

Jede Tranche eigener Commit, grün durch alle Gates (lint + test +
coverage-gate + docs-check), DoD-Line in dieser Datei trägt den
Commit-Hash (Konvention `[[feedback-done-slice-dod-hash]]`).

1. **T1 — Domain-Types + Driving-Ports + Sentinels.**
   - `internal/hexagon/domain/serviceup.go`: `ServiceStatus`-Struct
     (`{Name, ContainerStatus, Port, Healthcheck}`), `UpResult`
     (`Services []ServiceStatus`, `Stabilized bool`,
     **`Diagnostics []domain.Diagnostic`**),
     `StabilizationOutcome`-Enum (`Stabilized`/`RunningOnly`/`Failed`).
     - **`Diagnostics`-Vertrag:** Reuse des `domain.Diagnostic`-
       Typs aus M4 (`internal/hexagon/domain/diagnostic.go`;
       `{ID, Severity, Message, Hint}`). T1 friert den Vertrag
       ein, T4 füllt ihn: alle nicht-fatalen Beobachtungen aus
       dem Up-Lauf (insbesondere nicht-TCP-/Range-/unparsbare
       Ports aus `parseComposePort`-Tabelle) landen hier als
       `SeverityWarn`-Einträge, **nicht** als separater
       Stderr-Print. CLI-Output-Layer (T6) entscheidet auf
       Basis dieses Feldes, ob nach der Status-Tabelle eine
       zusätzliche Warn-Sektion gerendert wird. Damit ist
       klar, dass Diagnostics **Teil des Domain-Contracts** sind
       und nicht CLI-seitig nebenher entstehen — sonst wären
       sie aus Application-Tests heraus nicht assertierbar.
     - **`SeverityInfo`-Erweiterung in T1:** M4 hat heute nur
       `SeverityOK`/`SeverityWarn`/`SeverityError`. T1 prüft, ob
       `SeverityInfo` (informativer Hinweis, nicht-bewertend)
       schon existiert; falls nein, ergänzt T1 diesen Wert. Damit
       lassen sich Modus-Hinweise wie der `up.fire-and-forget`-
       Eintrag (Schritt 5 unten) ohne falsch-warn-Färbung
       transportieren. M4-`doctor` ist nicht-betroffen, weil es
       die neue Severity nicht emittiert; bestehende Severity-
       Tabellen in `domain.Diagnostic`-Tests werden um den Wert
       erweitert.
     - **ID-Konvention:** `up.port.<service>.<index>` für
       Port-Parse-Warnungen (Index zählt die `ports:`-Array-
       Position des Services); `up.fire-and-forget` für den
       `--timeout=0`-Modus-Hinweis; weitere Präfixe
       (`up.health.…`, `up.dns.…`) bleiben für künftige
       Tranchen reserviert.
   - `internal/hexagon/port/driving/up.go`: `UpRequest`
     (`BaseDir`, `Timeout time.Duration`, `ProgressSink io.Writer`),
     `UpResponse` (`Result domain.UpResult`), `UpUseCase`-Interface
     mit `Up(ctx, req)`.
     - **Timeout-Semantik (verbindlich, mit T7 abgestimmt):**
       `Timeout` ist eine `time.Duration`, **kein** roher
       Sekunden-Wert. `Timeout == 0` ⇒ fire-and-forget (kein
       Polling). `Timeout > 0` ⇒ Polling-Loop-Maximaldauer.
       `Timeout < 0` ⇒ Validation-Fehler, non-sentinel-error (CLI
       mappt auf Code 2).
     - **Konvertierung CLI → Domain (T7-Pflicht):** das CLI-Flag
       `--timeout <int>` (Sekunden) muss in der Subcommand-
       `RunE`-Closure explizit als
       `time.Duration(secs) * time.Second` konvertiert werden,
       **nicht** als `time.Duration(secs)` (das wäre Nanosekunden
       → 60 ns Timeout statt 60 s, und der Polling-Loop würde
       sofort ablaufen). Test in `cli_test.go`: ein Lauf mit
       `--timeout=60` und einer Mock-`Clock`, die in der ersten
       Iteration 30 s simuliert ⇒ erwartet **kein**
       `ErrStabilizationTimeout`. Pin der Konvertierung.
   - `internal/hexagon/port/driving/down.go`: `DownRequest`
     (`BaseDir`, `RemoveVolumes bool`, `AssumeYes bool`,
     `NonInteractive bool`, `ProgressSink io.Writer`),
     `DownResponse` (`RemovedVolumes bool`), `DownUseCase`-
     Interface mit `Down(ctx, req)`.
     - **Bewusst keine Stopped/Removed-Counters:** `docker
       compose down` liefert keinen strukturierten Counter-Output
       (stderr ist menschenlesbarer Phasen-Stream). Statt eine
       „unknown = -1"-Konvention zu pflegen, die ohnehin kein
       Caller auswerten kann, hält der MVP-Vertrag nur den
       binären Effekt (`RemovedVolumes` spiegelt den Input wider,
       sobald Compose erfolgreich returnt). CLI rendert eine
       Erfolgsmeldung der Form „environment stopped" bzw.
       „environment stopped, volumes removed". Falls strukturierte
       Counter später gebraucht werden (z. B. für `--json`),
       wäre das eigener V1-Slice mit `ComposePs`-Diff vor/nach
       `Down`.
     - **`NonInteractive` ist eigenes Request-Feld**, **nicht**
       aus `AssumeYes` ableitbar (siehe T5-Algorithmus für die
       Wahrheitstabelle). Es spiegelt das persistente CLI-Flag
       `--no-interactive` (LH-FA-CLI-005A) und macht den
       Confirmation-Pfad deterministisch ohne dass der Confirmer-
       Adapter den Interaktivitätsmodus selbst kennen muss.
   - **Sentinel-Schichtung (zentral, damit `errors.Is` über die
     Schichten hinweg trägt):** Sentinels liegen in zwei
     verschiedenen Packages, je nachdem wer sie emittiert. Das
     ist die einzige Schichtung, die `errors.Is`-Durchleitung
     **ohne** dass die Application über Driven-Port-Sentinels
     schweigen müsste — und ohne die wäre Code 11 in der Praxis
     unzugänglich (ein „erst-wrappen-dann-CLI-mappt"-Ansatz
     verbirgt das Engine-Sentinel im Wrap, sodass `errors.Is`
     in T7-CLI nichts mehr findet):

     **Driven-Port-Sentinels** in
     `internal/hexagon/port/driven/docker_engine.go` — vom Adapter
     emittiert, von der Application **als Wert durchgereicht**
     (via `fmt.Errorf("up service: ComposeUp: %w", err)`,
     **nicht** durch erneutes Wrappen unter einem
     Application-Sentinel):
     - `driven.ErrDockerUnavailable` (Docker-Daemon nicht erreichbar
       oder Compose-Plugin fehlt) → CLI-Code **11**. Adapter
       erkennt das deterministisch über die zwei Pre-Probes (siehe
       T2: `docker version`-Roundtrip + `docker compose version`-
       Roundtrip) plus `exec.LookPath`-Fehler beim Start des
       `docker`-Binaries; **nicht** über Stderr-Parse des
       eigentlichen `up`/`down`-Calls.
     - `driven.ErrComposeRuntime` (Compose-Stderr-Fehler **nach**
       bestandenen Pre-Probes) → CLI-Code **12**.

     **Application-/Driving-Sentinels** in
     `internal/hexagon/port/driving/` — vom Application-Service
     direkt emittiert:
     - `ErrProjectNotInitialized` (Reuse aus M5; kein `u-boot.yaml`)
       → Code **10**.
     - `ErrComposeFileMissing` (kein `compose.yaml` — fachlich vs.
       M3 fresh-init) → Code **10**.
     - `ErrConfirmationRequired` (`down --volumes` ohne `--yes` im
       `--no-interactive`-Modus) → Code **10**.
     - `ErrStabilizationTimeout` (Polling-Timeout) → Code **12**.

     **`errors.Is`-Durchleitungs-Vertrag (verbindlich für T4/T5):**
     1. Application-Code wrappt Engine-Fehler **niemals** in einem
        eigenen Sentinel. Konkret: kein `fmt.Errorf("...: %w",
        ErrComposeRuntime)` und kein `fmt.Errorf("...: %w",
        ErrDockerUnavailable)` in der Application; das würde die
        ursprüngliche Engine-Fehler-Identität überschreiben.
     2. Application wrappt nur **kontextuell** und nur unter dem
        **Original-Fehler** (`fmt.Errorf("up service: ComposeUp
        on %q: %w", baseDir, err)`), sodass
        `errors.Is(err, driven.ErrDockerUnavailable)` **und**
        `errors.Is(err, driven.ErrComposeRuntime)` aus T7-CLI
        nach dem Wrap weiter true liefern.
     3. T7-CLI-Mapping prüft die Sentinels in dieser festen
        Reihenfolge (am-spezifischsten zuerst): zuerst
        `driven.ErrDockerUnavailable` (Code 11), dann
        `driven.ErrComposeRuntime` (Code 12), dann
        `ErrStabilizationTimeout` (Code 12), dann die
        10er-Application-Sentinels, sonst Default Code 1.
     4. Pin-Test in `cli_test.go`: Fake-Engine returnt
        `driven.ErrDockerUnavailable` direkt aus
        `ComposeUp`/`ComposeDown`/`ComposePs`; CLI muss in allen
        drei Fällen Code 11 liefern, **nicht** Code 12. Zweites
        Pin-Test: Fake-Engine returnt `driven.ErrComposeRuntime`;
        CLI liefert Code 12. Drittes Pin-Test: Fake-Engine returnt
        einen generischen `errors.New("some compose stderr")`
        (ohne Sentinel); Application wickelt das in
        `driven.ErrComposeRuntime` am **Adapter-Rand** (T2:
        unbekannte Compose-Fehler werden vom Adapter selbst noch
        in `ErrComposeRuntime` gehüllt, bevor sie die Application
        sehen) — somit kommt am CLI immer ein Sentinel an, und
        Code 12 ist der Default-Pfad.
   - Tests: Tabellen-Tests für die Outcome-Enum, Sentinel-Mapping,
     Timeout-Sentinel-Validation.
   - DoD T1: Domain + Ports kompilieren, Sentinels einzeln in einem
     `cli_exit_mapping_test.go`-Stil getestet.

2. **T2 — `DockerEngine`-Driven-Port + Adapter + Fake.**
   - `internal/hexagon/port/driven/docker_engine.go` (neue Datei,
     getrennt von `docker.go`, um den DockerProbe-Kommentar nicht
     aufzublähen): `DockerEngine`-Interface mit
     `ComposeUp(ctx, dir string, opts ComposeUpOptions) (ComposeUpResult, error)`,
     `ComposeDown(ctx, dir string, opts ComposeDownOptions) error`,
     `ComposePs(ctx, dir string) ([]ComposeService, error)`.
     - `ComposeUpOptions`: `Detach bool` (true für M6),
       `ProgressSink io.Writer` (für Phasen-Output zum CLI; nil =
       discard).
     - `ComposeDownOptions`: `RemoveVolumes bool`,
       `ProgressSink io.Writer`.
     - `ComposeService`: `{Name, State, Health, Ports []string}`.
     - `ComposeUpResult`: `{Services []ComposeService}` (Snapshot
       direkt nach `up`-Return, vor Healthcheck-Polling).
   - `internal/adapter/driven/docker/engine.go`: Adapter mit
     `os/exec docker compose -f <dir>/compose.yaml up -d` /
     `down [-v]` / `ps --format json`. JSON-Format ist seit Compose
     v2.20 stabil (das ist gerade die doctor-Mindestversion —
     `LH-FA-DIAG-002`); kein Fallback nötig.
   - **Env-Error-Klassifikation (Adapter-Vertrag):** der Adapter
     übernimmt die Unterscheidung Env-Fehler vs. Runtime-Fehler
     **nicht** über Stderr-Parse (zu fragil, Compose-Strings
     wechseln zwischen Releases), sondern über zwei robustere
     Signale:
     1. **Binary-Lookup:** `exec.LookPath("docker")` vor jedem Call;
        Fehler ⇒ `ErrDockerUnavailable` (Mapping „Docker-Binary
        fehlt" → Code 11).
     2. **`exec.ExitError`-Pre-Probe:** vor `ComposeUp`/`ComposeDown`
        ruft der Adapter einmal `docker version --format '{{.Server.Version}}'`
        (daemon-roundtrip) und einmal `docker compose version
        --short` (plugin-roundtrip). Beide Probes liefern bereits
        in `DockerProbe` die identische Klassifikation; T2 zieht
        die Probe-Funktionen in `internal/adapter/driven/docker/`
        package-private heraus und wiederverwendet sie. Fehler in
        einem der Pre-Probes ⇒ `ErrDockerUnavailable` mit
        eingebettetem Probe-Detail. Nur wenn beide Pre-Probes grün
        sind, geht der eigentliche `up`/`down`-Aufruf weiter; ein
        Fehler dort ⇒ generischer `ErrComposeRuntime` (Code 12).
     Diese Schichtung garantiert ein deterministisches Code-11-vs-
     12-Verhalten und macht keinen Compose-Stderr-String load-
     bearing. Trade-off: zwei zusätzliche Roundtrips pro `up`-
     Aufruf (je ~50 ms). Akzeptabel, weil `up` ohnehin eine
     User-initiierte interaktive Operation ist; ein zusätzliches
     `--skip-preflight`-Flag wäre Premature Optimization.
   - **Stderr-Forwarding:** `cmd.Stderr = opts.ProgressSink`
     (default: `io.Discard`). Kein Buffer, kein Parse — Compose ist
     selbst für das Format zuständig. CLI verbindet später
     `os.Stderr`.
   - `internal/adapter/driven/docker/engine_test.go`: Build-Tagged
     `//go:build docker` (siehe `spec/architecture.md` §5 +
     `slice-m6-docker-integrationstests`). M6 selbst legt nur die
     Test-Datei + die `make test-docker`-Target-Erweiterung; der
     CI-Job für den Docker-Pfad bleibt im separaten Carveout-Slice.
   - `internal/hexagon/application/fakes_test.go` erweitert um
     `fakeDockerEngine` mit Skript-Setup `(Up/Down/Ps)→(result, err)`.
   - DoD T2: Adapter implementiert; `make verify-depguard` grün
     (Adapter darf `os/exec`, Application nicht); Build-Tag-Test
     existiert, schlägt lokal ohne Docker fehl (erwartet:
     `// +build docker` skip).

3. **T3 — `NetProbe`-Driven-Port + Adapter + Fake.**
   - `internal/hexagon/port/driven/netprobe.go`:
     `NetProbe`-Interface mit
     `DialTCP(ctx context.Context, host string, port int, timeout time.Duration) error`.
     Nil-error = erreichbar; non-nil error = nicht erreichbar
     (timeout / refused / unresolved).
   - `internal/adapter/driven/net/probe.go`: Adapter mit
     `net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)`.
     Bei `ctx.Err() != nil`: priorisiert ctx-Fehler vor net-Fehler.
   - Depguard: neue Regel `application-no-net` (Application darf
     `net` nicht importieren) wird in T3 aktiviert; das ist analog
     zur existierenden `application-no-yaml`/`application-no-os-exec`-
     Familie und kein eigener Carveout-Slice.
   - `internal/adapter/driven/net/probe_test.go`: Unit-Tests mit
     lokalem `net.Listen` + zufälligem Port (offen / refused) sowie
     `127.0.0.1:1` (refused-Path).
   - `fakeNetProbe` im Application-Test mit Map
     `{host:port → error}`.
   - DoD T3: Adapter + Fake + Depguard-Regel grün;
     `make verify-depguard` zeigt die neue Regel.

4. **T4 — `UpService` + Stabilisierungs-Polling.**
   - `internal/hexagon/application/upservice.go`:
     `UpService` mit DI `fs driven.FileSystem`,
     `yaml driven.YAMLCodec`, `engine driven.DockerEngine`,
     `net driven.NetProbe`, `logger driven.Logger`,
     `clock driven.Clock` (neuer kleiner Port, nur
     `Now() time.Time` + `Sleep(d time.Duration)`; sonst werden
     Polling-Tests flaky).
   - Falls noch kein `Clock`-Port existiert: T4 legt ihn an
     (analog zu `Confirmer`/`Logger` aus M4 — kleiner Helper-Port,
     der Tests deterministisch macht). Wenn schon vorhanden: reuse.
   - `Up(ctx, req)`-Algorithmus:
     1. `BaseDir != ""`, `Timeout >= 0` validieren; sonst
        non-sentinel-Error (CLI mappt auf Code 2).
     2. `u-boot.yaml`-Existenz + Parse ⇒ sonst
        `ErrProjectNotInitialized`.
     3. `compose.yaml`-Existenz prüfen ⇒ sonst
        `ErrComposeFileMissing`.
     4. `engine.ComposeUp(ctx, baseDir, {Detach: true,
        ProgressSink: req.ProgressSink})` ausführen. Fehler **nicht**
        in ein Application-Sentinel umhüllen — Adapter hat den
        Engine-Fehler bereits in `driven.ErrDockerUnavailable`
        oder `driven.ErrComposeRuntime` klassifiziert (T2-Vertrag).
        Application macht nur kontextuellen Wrap:
        `return UpResponse{}, fmt.Errorf("up service: ComposeUp on %q: %w", baseDir, err)`
        — damit bleibt `errors.Is(err, driven.ErrDockerUnavailable)`
        in T7-CLI true und Code 11 ist erreichbar.
     5. Wenn `req.Timeout == 0`: **sofort returnen ohne weiteren
        Engine-Call** —
        `return UpResponse{Result: domain.UpResult{Services: nil,
        Stabilized: false, Diagnostics: []domain.Diagnostic{{
        ID: "up.fire-and-forget", Severity: SeverityInfo,
        Message: "started with --timeout=0; status check skipped",
        Hint: "run `u-boot doctor` or `docker compose ps` for
        live status",
        }}}}, nil`.
        - **Begründung (Spec-Treue):** `LH-FA-UP-001` §970 sagt
          „Mit `--timeout=0` wird auf das Warten verzichtet; `up`
          beendet nach Initiierung der Compose-Aktionen." Ein
          zusätzlicher `engine.ComposePs`-Call wäre ein
          blockierender Daemon-Roundtrip nach dem `ComposeUp`-
          Return — semantisch kein Warten auf Stabilisierung,
          aber praktisch ein Sync-Punkt, der mit einem
          `driven.ErrDockerUnavailable` (Code 11) fehlschlagen
          könnte, obwohl der eigentliche `up` schon angekommen
          ist. Der User hat explizit fire-and-forget gewählt;
          jedes zusätzliche Engine-Risiko verwässert die Garantie.
        - **Konflikt mit `LH-FA-UP-003` (Status anzeigen) bewusst
          akzeptiert:** §988 fordert Status nach dem Start, aber
          nur als Konsequenz aus dem Default-Stabilisierungs-Pfad.
          Im `--timeout=0`-Modus ist Status per Definition nicht
          verfügbar (kein Polling) — die Spec-Forderung kollidiert
          nicht logisch, weil §970 die explizite Ausnahme ist. Der
          `up.fire-and-forget`-Diagnose-Eintrag (`SeverityInfo`,
          neue Severity-Variante in M4-`domain.Diagnostic` —
          falls noch nicht vorhanden, in T1 ergänzen) macht den
          Bewusstheits-Pfad sichtbar und CLI-Output-Layer (T6)
          rendert ihn als Info-Zeile statt einer leeren
          Status-Tabelle.
        - **Pin-Test in T4:** Fake-Engine mit Skript
          `ComposeUp → nil` und `ComposePs` so präpariert, dass
          ein Aufruf einen `panic("ComposePs must not be called
          when Timeout=0")` triggert. `Up` mit `Timeout=0` darf
          den Panic nicht auslösen. Damit ist die Garantie
          test-gepinnt, nicht nur dokumentiert.
     6. Sonst Polling-Loop:
        - `clock.Now()` als Startzeitpunkt; Loop-Intervall 500 ms
          (Konstante `pollInterval`).
        - In jeder Iteration `engine.ComposePs` rufen, Healthcheck/
          Running-Status klassifizieren; für jeden Service mit
          deklariertem TCP-Port `net.DialTCP` ausführen
          (sequentiell, je Service `dialTimeout=300ms` —
          Konstante).
        - **`ComposePs`-Fehler-Politik (verbindlich):** der
          Daemon kann mitten im Polling-Loop wegbrechen
          (User dreht Docker Desktop aus, Socket-Permission-
          Wechsel etc.) — der Adapter klassifiziert auch hier
          deterministisch in `driven.ErrDockerUnavailable` /
          `driven.ErrComposeRuntime` (selbe Pre-Probe-Heuristik
          wie T2-`ComposeUp`). Application reicht das Sentinel
          1:1 durch kontextuellen Wrap weiter und **bricht den
          Polling-Loop sofort ab**:
          `return UpResponse{}, fmt.Errorf("up service: ComposePs at t=%s: %w", elapsed, err)`.
          Damit bleibt sowohl der Env-Fall (Code 11) als auch der
          Runtime-Fall (Code 12) aus dem Loop heraus erreichbar.
          Wichtige Konsequenz: ein einzelner `ComposePs`-Fehler
          beendet `up`; es gibt **kein** Soft-Retry („vielleicht
          nächste Iteration"), weil ein fehlerhafter Daemon-Roundtrip
          kein Wackelkontakt ist, sondern ein Diagnose-Signal.
          Pin-Test: Fake-Engine returnt in der dritten
          `ComposePs`-Iteration `driven.ErrDockerUnavailable` ⇒
          `Up` returnt einen Fehler, `errors.Is(err,
          driven.ErrDockerUnavailable) == true`, kein weiterer
          Polling-Step erfolgt (Fake-Engine `PsCalls == 3`).
        - Erste vollständige Stabilisierung ⇒ Loop verlassen,
          `UpResult{Stabilized: true, …}` returnen.
        - `Failed`-Status eines Service ⇒ Polling-Abbruch mit
          fachlichem Wrap: `return UpResponse{},
          fmt.Errorf("up service: %q reached state %q: %w",
          name, state, driven.ErrComposeRuntime)`. Bewusst
          **mit** `%w` auf das Driven-Sentinel, weil
          `Failed`-State eine Compose-Ausführungsbeobachtung ist
          und denselben Code-12-Pfad bedienen soll wie ein
          stderr-getriebener `ComposeUp`-Fehler. Das ist die
          **einzige** Stelle, an der die Application ein
          Driven-Sentinel **erzeugt** statt durchzuleiten;
          gerechtfertigt, weil hier kein Engine-Call den
          Fehler trägt — der `Failed`-State kommt aus den
          ComposePs-**Daten**, nicht aus einem Engine-Error.
        - `clock.Now()-start > req.Timeout` ⇒
          `ErrStabilizationTimeout` mit Liste der noch nicht
          stabilisierten Services (Code 12, eigener
          Application-Sentinel — kein Driven-Wrap, weil der
          Timeout u-boot-eigene Polling-Logik ist, nicht Compose).
        - `ctx.Err()` ⇒ wrappen + returnen (Cancel-Pfad,
          non-sentinel).
   - **Service-Definition + TCP-Port-Extraktion (robust):** der
     Polling-Pfad braucht pro Service die deklarierten Ports und
     die Existenz eines Healthchecks. Compose erlaubt vier
     beobachtbare Port-Syntaxen, die T4 alle akzeptieren muss
     ohne den Unmarshal zum Scheitern zu bringen — gescheitertes
     Unmarshal eines Pflichtfelds würde `up` mit Code 12 abbrechen,
     obwohl die Spec einen `warn`-Pfad fordert. Die Parsing-Logik
     lebt als isolierter Helper `parseComposePort(raw any)
     (portProbeTarget, error)` in `upservice.go`, getestet pro
     Syntax:

     | Compose-Syntax (`ports:`-Element)          | Beispiel                                             | Klassifikation                                                |
     | ------------------------------------------ | ---------------------------------------------------- | ------------------------------------------------------------- |
     | nackte Integer-Kurzform                    | `5432`                                               | TCP `localhost:5432`                                          |
     | String-Kurzform `host:container`           | `"5432:5432"`                                        | TCP `localhost:5432` (Host-Port = links)                      |
     | String-Kurzform mit Host-Bind              | `"127.0.0.1:5432:5432"`                              | TCP `127.0.0.1:5432`                                          |
     | String-Kurzform mit Protokoll              | `"5432:5432/udp"`                                    | `warn`-Diagnose, kein Probe (nicht-TCP)                       |
     | String-Kurzform mit Range                  | `"5000-5010:5000-5010"`                              | `warn`-Diagnose, kein Probe (nicht eindeutig)                 |
     | Long-Syntax-Mapping                        | `{target: 5432, published: 5432, protocol: tcp}`     | TCP `localhost:5432`; `protocol: udp` ⇒ `warn`                |
     | Long-Syntax mit `host_ip`                  | `{target: 5432, published: 5432, host_ip: 127.0.0.1}`| TCP `127.0.0.1:5432`                                          |
     | unbekannte / fehlerhafte Form              | `[1, 2, 3]`, leerer String, nicht-numerischer Host   | `warn`-Diagnose, **kein Fail** (graceful, Service stabilisiert auf Healthcheck/Running) |

     Verbindlich: `ports`-Wert wird als `[]any` mit
     `YAMLCodec.Unmarshal` gelesen (heterogene Slice-Elemente:
     Strings, Integer, Mapping). Jedes Element läuft durch
     `parseComposePort`. Parse-Fehler ist **nie** ein Service-
     Fehler — er erzeugt einen `Diagnostics`-Eintrag in
     `UpResult` (Severity `warn`, ID `up.port.<service>.<index>`)
     und der Service stabilisiert allein auf Basis seines
     Healthcheck-/Running-Status (LH-FA-UP-001 §969). Test-
     Fixture: jeder Tabellen-Zeile genau ein Unit-Test.

     Healthcheck-Existenz: `services.<name>.healthcheck` als
     Mapping vorhanden ⇒ Healthcheck-Pfad. Compose erlaubt auch
     `healthcheck: { disable: true }` ⇒ zählt als **nicht
     vorhanden** (running ausreichend). Pin-Test für diesen
     Edge-Case.
   - Tests (alle in `_test`-Package, mit `fakeDockerEngine`,
     `fakeNetProbe`, `fakeClock`):
     - Fire-and-forget (`Timeout=0`) ⇒ kein Polling, kein
       NetProbe-Call.
     - Single-Healthcheck-Service stabilisiert in 2 Polls.
     - Single-Running-Service (kein Healthcheck) **und kein Port
       deklariert** ⇒ stabilisiert sobald `stateRunning` in der
       ersten Polling-Iteration (allein-`running`-Pfad,
       LH-FA-UP-001 §967). Pin verhindert die Drift „warten, bis
       irgendwas anderes auch noch passiert".
     - Single-Running-Service (kein Healthcheck) **mit
       deklariertem TCP-Port** ⇒ stabilisiert sobald `running` +
       Port erreichbar.
     - Timeout: Service bleibt unhealthy ⇒ `ErrStabilizationTimeout`,
       Liste der Pending-Services in error-Detail.
     - Failed-Service ⇒ `ErrComposeRuntime` mit Service-Name,
       Polling-Abbruch ohne weitere Iteration.
     - Nicht-TCP-Port-Service ⇒ Diagnostics-Eintrag (`warn`), aber
       Stabilisierung auf Basis Healthcheck/Running.
     - **State-Normalisierung (Pin):** Fake-Engine returnt
       `ComposeService.State` als `"Running"` (Großbuchstabe),
       `"RUNNING"`, `"running"` ⇒ alle drei werden auf
       `stateRunning` normalisiert, Stabilisierung verhält sich
       identisch. Selbe Tabellen-Tests für `restarting`/
       `RESTARTING`, `exited`/`Exited`.
     - **Restart-Loop-Detection (Pin):** Fake-Engine returnt
       in den ersten zwei Polling-Iterationen
       `stateRestarting`, danach `stateRunning` + healthy ⇒
       Service stabilisiert (`restartObservations` resetet auf
       0, Loop erreicht Threshold nie). Zweiter Sub-Test:
       Fake returnt `stateRestarting` in drei aufeinander-
       folgenden Iterationen ⇒ Service als `Failed` klassifiziert
       (Service-Name + State im Error-Detail), Polling-Abbruch
       in der dritten Iteration.
     - **Unbekannter State-String (Soft-Pfad, Pin):** Fake
       returnt `State: "frobnicating"` in Iteration 1 und 2,
       dann `stateRunning` + healthy in Iteration 3 ⇒ Service
       stabilisiert (kein `Failed`), `UpResult.Diagnostics`
       enthält **genau einen** Eintrag mit ID
       `up.state.<service>.unknown` und `"frobnicating"` im
       `Message`-Feld (Idempotenz-Test: nicht zwei Einträge
       trotz zwei Beobachtungen). Zweiter Sub-Test: Fake
       returnt `"frobnicating"` über den vollen Timeout ⇒
       `ErrStabilizationTimeout` (Code 12, normaler Pfad), nicht
       `ErrComposeRuntime` — damit ist der Compose-Erweiterungs-
       Schutz test-gepinnt.
     - `ErrProjectNotInitialized` / `ErrComposeFileMissing`-Pfade.
     - Ctx-Cancel mid-poll ⇒ ctx.Err() durchgereicht.
     - **Engine-Sentinel-Durchleitung (drei Pins):**
       (a) Fake-Engine returnt `driven.ErrDockerUnavailable`
       aus `ComposeUp` ⇒ `Up` returnt einen Fehler mit
       `errors.Is(err, driven.ErrDockerUnavailable) == true`,
       keine Polling-Iteration.
       (b) Fake-Engine returnt `driven.ErrComposeRuntime` aus
       `ComposeUp` ⇒ analog mit `ErrComposeRuntime`.
       (c) Fake-Engine returnt in der dritten `ComposePs`-
       Iteration `driven.ErrDockerUnavailable` ⇒ `Up` returnt
       einen Fehler mit `errors.Is(err,
       driven.ErrDockerUnavailable) == true`, Fake-Engine
       `PsCalls == 3` (kein Soft-Retry).
   - DoD T4: alle obigen Tests grün; `make gates` grün; keine
     time.Sleep-basierten Tests (Clock-Port).

5. **T5 — `DownService` + `--volumes`-Confirmer-Pfad.**
   - `internal/hexagon/application/downservice.go`: `DownService`
     mit DI `fs`, `engine`, `confirmer driven.Confirmer`, `logger`.
   - `Down(ctx, req)`-Algorithmus:
     1. `BaseDir != ""` validieren.
     2. `u-boot.yaml` + `compose.yaml`-Existenz wie in `Up`.
     3. Wahrheitstabelle für den `--volumes`-Bestätigungspfad
        (verbindlicher Vertrag; T7-CLI mappt die drei
        DownRequest-Flags direkt aus `--volumes`/`--yes`/
        `--no-interactive`):

        | `RemoveVolumes` | `AssumeYes` | `NonInteractive` | Verhalten                                          |
        | --------------- | ----------- | ---------------- | -------------------------------------------------- |
        | false           | *           | *                | direkt zu Schritt 4, kein Confirmer-Call           |
        | true            | true        | *                | direkt zu Schritt 4, kein Confirmer-Call           |
        | true            | false       | true             | **sofort `ErrConfirmationRequired`** (Code 10), kein Confirmer-Call, kein Engine-Call |
        | true            | false       | false            | `confirmer.Confirm("Remove all volumes? Data will be lost.", false)`; `(true, nil)` ⇒ Schritt 4; `(false, nil)` ⇒ `ErrConfirmationRequired`; `error` ⇒ wrappen |

        Das `NonInteractive`-Flag macht den `--no-interactive`-
        Fail-Fast-Pfad **vor** dem Confirmer-Call deterministisch:
        die Application braucht keinen Confirmer-Adapter, der
        seinen Interaktivitätsmodus selbst kennt, und der Spec-
        Pfad LH-FA-CLI-005A („nicht-interaktiv ohne `--yes` ⇒
        Code 10") ist eindeutig an dieser Tabelle ablesbar.
     4. `engine.ComposeDown(ctx, baseDir, {RemoveVolumes:
        req.RemoveVolumes, ProgressSink: req.ProgressSink})`.
        Engine-Pre-Probes (T2) klassifizieren Env- vs.
        Runtime-Fehler bereits am Adapter-Rand; Application macht
        **nur** kontextuellen Wrap unter dem Original-Fehler:
        `return DownResponse{}, fmt.Errorf("down service: ComposeDown on %q: %w", baseDir, err)`.
        `errors.Is(err, driven.ErrDockerUnavailable)` (Code 11)
        und `errors.Is(err, driven.ErrComposeRuntime)` (Code 12)
        bleiben nach diesem Wrap aus T7-CLI sichtbar.
     5. `DownResponse{RemovedVolumes: req.RemoveVolumes}`
        zurückgeben. Keine Counter (siehe T1-Begründung); CLI
        rendert auf Basis dieses Booleans die Erfolgsmeldung.
   - Tests (eine Zeile pro Tabellenzeile + Engine-Fehlerpfade):
     - `RemoveVolumes=false`: kein Confirmer-Call, Engine-Call
       mit `RemoveVolumes=false`.
     - `RemoveVolumes=true, AssumeYes=true`: kein Confirmer-Call,
       Engine-Call mit `RemoveVolumes=true` (zwei Sub-Tests für
       `NonInteractive=true`/`false`; beide identisch).
     - `RemoveVolumes=true, AssumeYes=false, NonInteractive=true`:
       `ErrConfirmationRequired`, **kein** Confirmer-Call,
       **kein** Engine-Call. Assertion am Fake-Engine: `Calls ==
       0`.
     - `RemoveVolumes=true, AssumeYes=false, NonInteractive=false`,
       Confirmer returnt `(true, nil)` ⇒ Engine-Call mit
       `RemoveVolumes=true`.
     - Selbe Konstellation, Confirmer returnt `(false, nil)` ⇒
       `ErrConfirmationRequired`, kein Engine-Call.
     - Selbe Konstellation, Confirmer returnt error ⇒ wrappped
       error (kein Sentinel).
     - `ErrProjectNotInitialized`/`ErrComposeFileMissing`-Pfade.
     - **Engine-Sentinel-Durchleitung (Pin):** Fake-Engine
       returnt `driven.ErrDockerUnavailable` aus `ComposeDown`
       ⇒ `Down` returnt einen Fehler mit
       `errors.Is(err, driven.ErrDockerUnavailable) == true`.
       Zweiter Test mit `driven.ErrComposeRuntime` ⇒
       `errors.Is(err, driven.ErrComposeRuntime) == true`.
       Beide Sentinels überleben den Application-Wrap.
   - DoD T5: alle Tests grün.

6. **T6 — Output-Layer (LH-FA-UP-003 + LH-NFA-PERF-002).**
   - `internal/adapter/driving/cli/statusview.go` (neue Datei oder
     in `up.go` lokal — T6-Entscheidung): Tabellen-Renderer für
     `UpResult.Services` mit den vier Pflicht-Spalten Name /
     Container-Status / Port / Healthcheck. Sortierung: alphabetisch
     nach Service-Name; reproduzierbarer Output für Tests.
   - Healthcheck-Spalte: `-` wenn nicht definiert, sonst
     `healthy`/`unhealthy`/`starting` (Compose-Strings).
   - Port-Spalte: kommagetrennt `"5432:5432"`, mehrere Mappings
     je Zeile zusammengezogen.
   - **Progress-Forwarding:** CLI verkabelt `os.Stderr` als
     `ProgressSink` an `engine.ComposeUp`/`ComposeDown`, sodass
     Compose-Phasen (`Pulling…`, `Creating…`, `Starting…`,
     `Healthchecking…`) live durchlaufen.
   - **`--quiet` (LH-FA-CLI-005, persistent Root-Flag, von M4-T7
     verkabelt) — Geltungsbereich pro Subcommand explizit
     getrennt:**
     - **`up --quiet`**: unterdrückt die finale Status-Tabelle
       (LH-FA-UP-003) und die Diagnose-Warn-Sektion am Ende; der
       Erfolgs-Exit (Code 0) bleibt ohne stdout-Print stehen.
     - **`down --quiet`**: unterdrückt die einzeilige
       Erfolgsmeldung (`"environment stopped"` /
       `"environment stopped, volumes removed"`); `down`
       produziert keine Status-Tabelle, daher gibt es hier
       nichts Tabellarisches zu unterdrücken — die Klarstellung
       ist nur, dass `--quiet` **trotzdem** den Erfolgs-Text
       wegnimmt, damit `down --quiet` in CI-Skripten einen
       leeren stdout-Output liefert. Asymmetrie zu `up` ist
       beabsichtigt und an die jeweilige Subcommand-Ausgabe
       gebunden, nicht an einen gemeinsamen „finale Tabelle"-
       Begriff.
     - **Progress-Stream bleibt in beiden Subcommands**: der
       `engine.ComposeUp`/`ComposeDown`-stderr-Stream wird
       **nicht** durch `--quiet` stummgeschaltet, weil
       LH-NFA-PERF-002 die Sichtbarkeit der Pull/Create/Start/
       Healthcheck-Phasen explizit fordert. Trade-off bewusst:
       `--quiet` reduziert das u-boot-eigene UI, nicht den
       durchgereichten Compose-Output; ein zukünftiges
       `--silent` (LH-OPEN-Folgepunkt, nicht in M6) wäre der
       Schalter, der auch Compose stummschaltet.
     - **Tests** in `cli_test.go` separat pro Subcommand:
       (a) `up --quiet` ⇒ kein Tabellen-Output, kein Diagnose-
       Block, `engine.ComposeUp` mit nicht-nil `ProgressSink`
       aufgerufen.
       (b) `down --quiet` ⇒ keine Erfolgsmeldung,
       `engine.ComposeDown` mit nicht-nil `ProgressSink`
       aufgerufen.
   - Tests in `cli_test.go`: Tabellen-Output mit golden-File-
     Fixtures pro Szenario (0 Services, 1 Service ohne
     Healthcheck, 1 Service mit Healthcheck, 2 Services
     gemischt).
   - **Fire-and-forget-Output-Pin (verbindlich):** zusätzliches
     golden-Fixture `up_timeout_0.txt`, dessen Inhalt
     **garantiert keine** Service-Tabellen-Spaltenüberschriften
     enthält (Assertion via `strings.Contains` auf
     `"SERVICE"`/`"NAME"`-Header == false), sondern ausschließlich
     die Info-Zeile aus dem `up.fire-and-forget`-Diagnostic-
     Eintrag. Wenn ein zukünftiger Renderer-Refactor versehentlich
     auch leere Tabellen-Header druckt, schlägt dieser Test fehl,
     bevor die Spec-Verletzung in Production landet.
   - DoD T6: golden-Tests grün; Output stabil bei wiederholtem
     Lauf.

7. **T7 — CLI-Subcommands `up`/`down` + e2e-Acceptance.**
   - `internal/adapter/driving/cli/up.go`: Cobra-Subcommand `up`
     mit Flag `--timeout <int>` (default 60). Negative Werte ⇒
     Cobra-Validation-Error (Code 2). DI über das `cli.Build`-
     Konstruktor-Muster (siehe `cli.go` / Add-Subcommand-Muster
     aus M5-T7).
   - `internal/adapter/driving/cli/down.go`: Cobra-Subcommand
     `down` mit Flags `--volumes` (bool, default false), `--yes`
     (bool, default false). `--no-interactive` ist persistentes
     Root-Flag (M4-T7) und wird in der `RunE`-Closure aus
     `cmd.Flags().GetBool("no-interactive")` gelesen und 1:1 als
     `DownRequest.NonInteractive` durchgereicht; analog
     `--yes` → `AssumeYes`. `--volumes` → `RemoveVolumes`.
   - **Zwei verschiedene Exklusivitäts-/Konflikt-Pfade — bewusst
     getrennt halten** (Spec hat hier zwei nicht-überlappende
     Regeln, die beim ersten Lesen wie ein Widerspruch wirken):

     | Regel                                  | Spec-Stelle                             | Wann triggert                                                                   | Exit-Code | Wo implementiert                                            |
     | -------------------------------------- | --------------------------------------- | ------------------------------------------------------------------------------- | --------- | ----------------------------------------------------------- |
     | **Flag-Exklusivität**                  | LH-FA-CLI-005A §235 / §258              | `--yes` UND `--no-interactive` werden **beide** gesetzt — unabhängig von der Aktion | **2**     | Root-Resolver (M4-T7, kein M6-Touch)                        |
     | **Destruktiver Bestätigungs-Abbruch**  | LH-FA-CLI-005A §254                     | `--volumes` UND `--no-interactive` UND **nicht** `--yes` — explizit destruktiv  | **10**    | T5-Application (`ErrConfirmationRequired`)                  |

     Lese-Reihenfolge in `down`: Cobra/Root-Resolver prüft §235
     **zuerst** (Exit 2, fail-fast vor jedem Service-Aufruf);
     erst wenn beide Flags zusammen valide sind (höchstens eines
     gesetzt) erreicht der Aufruf den `DownService`, der dann
     §254 über die T5-Wahrheitstabelle umsetzt (Exit 10).
     Tests pinnen beide Pfade getrennt:
     (a) `down --volumes --yes --no-interactive` ⇒ Code 2,
     kein `DownService`-Call, kein `Confirmer`-Call;
     (b) `down --volumes --no-interactive` (ohne `--yes`) ⇒
     Code 10, `DownService`-Call ja, `Confirmer`-Call nein.
     Damit hat das Team eine eindeutige Lesart: Exklusivität ist
     CLI-Validation, Bestätigungsabbruch ist fachlich.
   - `cli.go`: beide Subcommands registrieren; Exit-Code-Mapping
     erweitern. Pin als Tabelle (Test in `cli_test.go` pinned das
     1:1):

     | Sentinel (Package)                  | Exit-Code | Begründung                                        |
     | ----------------------------------- | --------- | ------------------------------------------------- |
     | `driving.ErrProjectNotInitialized`  | 10        | fachliche Validierung (kein u-boot.yaml)          |
     | `driving.ErrComposeFileMissing`     | 10        | fachliche Validierung (kein compose.yaml)         |
     | `driving.ErrConfirmationRequired`   | 10        | LH-FA-CLI-005A (`down --volumes` ohne `--yes`)    |
     | `driven.ErrDockerUnavailable`       | **11**    | Umgebung (Docker/Compose-Plugin nicht verfügbar) |
     | `driving.ErrStabilizationTimeout`   | 12        | Ausführung (Polling-Timeout)                      |
     | `driven.ErrComposeRuntime`          | 12        | Ausführung (Compose-Stderr-Fehler)                |
     | (CLI-Validation in Cobra)           | 2         | `--timeout=-1`, unbekannte Flags                  |

     **Mapping-Reihenfolge in der `RunE`-Closure (verbindlich,
     siehe T1-Durchleitungs-Vertrag Schritt 3):** die
     Driven-Sentinels werden **zuerst** abgefragt
     (`errors.Is(err, driven.ErrDockerUnavailable)` →
     `errors.Is(err, driven.ErrComposeRuntime)`), erst danach
     die Driving-/Application-Sentinels. Falsche Reihenfolge
     (z. B. zuerst auf `driving.ErrStabilizationTimeout`) ist
     kein Funktionsfehler (Sentinels überschneiden sich nicht),
     aber Lese-Reihenfolge soll der Schicht-Hierarchie folgen.
   - `cli_test.go`: end-to-end-Tests gegen `fakeDockerEngine` +
     `fakeNetProbe` + `fakeClock` für:
     - `up` Happy Path (postgres-Service stabilisiert).
     - `up --timeout=0` (kein Polling): Fake-Engine-Assertion
       `PsCalls == 0` (kein einziger `ComposePs`-Roundtrip),
       Exit-Code 0, golden-Fixture `up_timeout_0.txt` matched.
     - `up --timeout=1` (Timeout-Pfad mit Pending-Service).
     - `up --timeout=-1` ⇒ Code 2.
     - `down` ohne `--volumes`.
     - `down --volumes --yes` ⇒ keine Confirmation.
     - `down --volumes` ohne `--yes`, Confirmer returnt false ⇒
       Code 10.
     - `up` in Verzeichnis ohne `u-boot.yaml` ⇒ Code 10.
   - **`LH-AK-002`-Acceptance** (Integration über die existierenden
     M5-Postgres-Templates): Test seed'et ein
     u-boot-Projekt mit aktivem `postgres`-Service (analog
     M5-T4c-e2e), führt `up` gegen `fakeDockerEngine` aus, das
     einen `healthy`-State nach 2 Polls liefert, und asserted die
     vier LH-AK-002-Erwartungen (Service in compose.yaml,
     `.env.example`-Block, Healthcheck erreicht, Port erreichbar).
     Der echte Compose-Lauf (LH-AK-002 mit echter Docker-Engine)
     bleibt im separaten [`slice-m6-docker-integrationstests`](../open/slice-m6-docker-integrationstests.md).
   - DoD T7: `make gates` grün; LH-AK-002-Acceptance-Test grün;
     Roadmap auf M6 = Done; M5-add-postgres-DoD-Link auf
     LH-AK-002-Erfüllung erweitern (separater Doku-Commit).

## Akzeptanzkriterien (Slice-Level)

- `LH-FA-UP-001..004` abgehakt.
- `LH-AK-002` (PostgreSQL-Flow) als CLI-e2e-Test grün
  (Docker-Engine-Fake; reale Docker-Acceptance über den
  Carveout-Slice).
- `make gates` grün (lint + test + coverage-gate + docs-check).
- Exit-Codes 2/10/11/12 für die in §Auslöser tabellierten Fehlerpfade
  pinned.
- Keine neuen temporären Carveouts ohne Slice-Plan
  (`LH-FA-PROJDOCS-005`); falls T-Granular doch einen Carveout
  öffnet, parallel Slice in `open/` + Eintrag in `carveouts.md`.

## Out of Scope

- **`LH-FA-UP-005` Logs**: eigener V1-Slice
  (`slice-v1-up-logs.md` o. ä., wird angelegt, sobald V1-Stoßrichtung
  klar ist; Carveout entsteht erst bei Bedarf).
- **`--json`-Output für up/down**: V1, analog M4-doctor-Vertagung.
- **Auto-Recovery / Retry**: scheitert ein Service, bricht `up` ab;
  keine eigene Retry-Logik in M6 (User entscheidet via `up`-Re-Run).
- **Multi-Compose-File-Support** (`-f file1 -f file2`): M6 hält
  am Single-`compose.yaml` aus M5-Scaffold fest.
- **`docker compose pull` als eigener Subcommand**: in M6 ist Pull
  Teil von `up`; ein expliziter `u-boot pull` ist nicht in der Spec
  und käme als V1-Folge.
- **Healthcheck-Definition durch u-boot**: Healthchecks kommen
  ausschließlich aus den Add-on-Templates (M5 hat `pg_isready` für
  postgres geliefert); `up` interpretiert nur, was Compose meldet.

## Risiken & offene Punkte (zur Klärung in T-Implementation)

- **`docker compose ps --format json` Schema-Stabilität:** seit
  Compose v2.20 (doctor-Mindestversion) als stabil dokumentiert,
  aber Felder können sich erweitern. T2 friert ein minimales
  `ComposeService`-Mapping ein und ignoriert unbekannte Felder
  (`json.Decoder` mit `DisallowUnknownFields = false`).
- **`net.DialTimeout`-Granularität:** 300 ms Probe-Timeout vs.
  500 ms Poll-Intervall — Polling-Loop muss serialisiert sein,
  sonst überlappen sich Probes. T4-Implementierung: pro Iteration
  alle Probes sequentiell, danach `clock.Sleep(pollInterval)`.
- **Compose-Stderr-Volumen bei großen Pulls:** kein Buffer, direktes
  Forwarding an `os.Stderr` — funktioniert für interaktive CLI;
  Logging-Adapter (slog) sieht den Stream nicht. Akzeptiert: das
  ist Compose-eigener Output und gehört nicht ins strukturierte
  Logging.
- **`down` ohne vorheriges `up`:** Compose returnt typischerweise
  einen No-op; T5 prüft das **nicht** vorher (kein
  `ComposePs`-Pre-Check), sondern verlässt sich auf
  Compose-Idempotenz. Falls Tests einen Pre-Check erzwingen
  wollen, ist das ein T5-Tranche-Detail.

## Bezug

- Auslösende Spec: `LH-FA-UP-001..004` (`spec/lastenheft.md` §4.6),
  `LH-FA-CLI-005A`, `LH-FA-CLI-006`, `LH-NFA-PERF-002`, `LH-AK-002`.
- Vorgänger: [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  (liefert die Postgres-Templates und den `services`-Schema-Block,
  ohne den `up` keinen Service zu starten hätte).
- Nachfolger: [`slice-m6-docker-integrationstests`](../open/slice-m6-docker-integrationstests.md)
  (fachlich freigeschaltet durch T7); danach M7 (`u-boot generate`)
  oder MVP-Closure (`LH-FA-DEV-001..005` Devcontainer-Mindestumfang
  + `LH-AK-001/005..007`).
