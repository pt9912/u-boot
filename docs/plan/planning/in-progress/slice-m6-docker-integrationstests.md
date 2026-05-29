# Slice M6: Docker-Integrationstests (Build-Tag-Pfad)

> **Status:** In progress (Sub-Tranchen)
> **DoD:** Sub-T1 ✅ `ab9ff4a` (+ review `d5ac2c3`) / Sub-T2 ✅ `bcd8486` / Sub-T3 ⏳ / Sub-T4 ⏳ / Stabilisierung (3× grün, `continue-on-error: false`) ⏳

## Sub-Tranchen-Schnitt

Der Slice wird in vier Sub-Tranchen gegliedert, damit jede Schicht
einzeln review- und merge-bar bleibt. Jeder Sub-Hash steht in der
DoD-Zeile; die Reihenfolge ist hart (Sub-T4 setzt T1–T3 voraus).

- **Sub-T1 — Adapter-Verhaltens-Pins.** Die zwei Pins, die direkt
  am `internal/adapter/driven/docker/`-Adapter ansetzen:
  `engine_progressstream_docker_test.go` (LH-NFA-PERF-002 via
  `io.Pipe`-Event-Ordering) + `engine_psjsonschema_docker_test.go`
  (LH-FA-DIAG-002 Compose-JSON-Schema-Snapshot). Beide nutzen ein
  inline-deklariertes minimales `compose.yaml`-Fixture. ~250 LoC.
- **Sub-T2 — Application-Verhaltens-Pins.** UpService-spezifische
  Pins, die über das Adapter-Verhalten hinaus die polling-Loop-
  Klassifikation in der Application-Schicht prüfen:
  `internal/hexagon/application/upservice_healthcheck_docker_test.go`
  (LH-FA-UP-001 §966) +
  `internal/hexagon/application/upservice_portprobe_docker_test.go`
  (LH-FA-UP-001 §968). Beide instanziieren `application.NewUpService`
  mit den realen Adaptern (DockerEngine + NetProbe + Clock + FS +
  YAML). ~200 LoC.
- **Sub-T3 — End-to-end-Verhaltens-Pins.** Voller Stack via
  `internal/e2e/` (neues Package):
  `postgres_acceptance_docker_test.go` (LH-AK-002) +
  `down_volumes_docker_test.go` (LH-FA-UP-004 §1015). Diese Tests
  wickeln den kompletten `init → add postgres → up → down`-Flow
  gegen die echte Engine ab. ~250 LoC.
- **Sub-T4 — CI-Wiring + Doku.** GitHub-Actions-Workflow
  (`integration-docker`-Job mit `continue-on-error: true` als
  Stabilisierungs-Maßnahme) + `docs/user/quality.md` §2 Tests-
  Update + DoD-Eintrag in dieser Datei. Carveout-Eintrag bleibt
  offen bis 3× grün + erster Lauf ohne `continue-on-error`. ~100 LoC.

Nach Sub-T4 wechselt der Slice-Status auf
**"In progress — Stabilisierung pending"**. Die finale
Carveout-Aufhebung folgt mit einer separaten PR, die das
Stabilization-Evidence-Block-Audit erfüllt (siehe
"Akzeptanzkriterien" unten); diese PR mergt
`continue-on-error: false` und entfernt den `carveouts.md`-Eintrag.

## Auslöser

`spec/architecture.md` §5 beschreibt eine Build-Tag-Konvention für
Adapter-Integrationstests gegen die echte Docker-Engine:

```
//go:build docker
```

mit `go test -tags docker ./...`. Stand M6-T2 (`84a676c`) existiert
der `internal/adapter/driven/docker/`-Adapter (Engine + Probe), aber
es fehlen weiterhin:

- ein CI-Stage / Make-Target, der die getaggten Tests ausführt;
- echte Compose-gegen-Daemon-Tests jenseits des
  `engine_docker_test.go`-Skeletons aus T2 (das nur die
  Build-Tag-Verkabelung absichert);
- Verhaltens-Pins gegen reales `docker compose`, insbesondere für
  `LH-NFA-PERF-002` (Pull/Create/Start/Healthcheck-Phasen-Stream
  in stderr — siehe Akzeptanzkriterien unten).

Die Build-Tag-Konvention ist bisher nur dokumentiert
(`LH-FA-PROJDOCS-005`), aber jetzt teilweise im Code verankert; der
Carveout-Status bleibt aktiv, bis CI den getaggten Pfad ausführt
und die Verhaltens-Pins existieren.

## Aufhebungsbedingung

Stand M6-T2 (`84a676c`) ist der Adapter implementiert und das
Build-Tag-Skeleton (`engine_docker_test.go`) liegt vor. Was zur
Aufhebung des Carveouts fehlt:

1. **Verhaltens-Pins gegen echte Compose-Engine** — über das
   Smoke-Skeleton hinaus (siehe Akzeptanzkriterien-Tabelle unten);
2. **`make test-docker`**-Target, der die getaggten Tests ausführt;
3. **CI-Job**, der das Docker-Socket mountet und `make test-docker`
   aufruft (ergänzt `make ci`, nicht `make gates`).

## Akzeptanzkriterien

### Strukturelle Bedingungen

- **Netzwerk-Namespace-Voraussetzung für Tests mit TCP-Probe-
  Assertions**: alle Pins, die `net.DialTCP` von der Test-Seite
  gegen einen Compose-veröffentlichten Port machen — heute Sub-T2
  §968 (`upservice_portprobe_docker_test.go`) plus Sub-T3 LH-AK-002
  (`postgres_acceptance_docker_test.go`) — verlangen, dass die
  Test-Prozess-Netzwerk-Namespace und die Docker-Daemon-Netzwerk-
  Namespace identisch sind. Konkret: entweder das Test-Binary
  läuft direkt auf dem Host (mit lokal installiertem `docker` und
  `docker compose`), oder es läuft in einem Container mit
  `--network=host`. Die einfache "Container mit gemountetem
  `/var/run/docker.sock`"-Variante reicht **nicht**: Compose
  veröffentlicht die Ports auf dem **Host**-Loopback, das
  Test-Binary sieht aber den **Container**-Loopback — der Probe
  schlägt mit `connection refused` fehl, obwohl der Service
  korrekt läuft. **Sub-T4-Makefile-Verkabelung muss diese
  Anforderung erfüllen** (Empfehlung: `docker run --network=host`
  oder Host-natives `go test -tags docker`); andernfalls würde
  der `integration-docker`-CI-Job §968 und LH-AK-002 falsch-rot
  liefern.
- `internal/adapter/driven/docker/engine_docker_test.go` (existiert
  seit T2) wird um zusätzliche Verhaltens-Tests aus der Tabelle
  unten ergänzt; optional weitere `*_docker_test.go`-Dateien.
- `make test-docker` führt `go test -tags docker ./...` aus (nicht
  nur `./internal/adapter/driven/docker/...`) — entweder direkt
  auf dem Host (mit gemountetem `/var/run/docker.sock`) oder in
  einer Docker-in-Docker-Variante. **Begründung für den
  Repo-weiten Scope**: die Verhaltens-Pins unten liegen über zwei
  Schichten verteilt — adapter-lokale Pins
  (`internal/adapter/driven/docker/`) und application-/end-to-end-
  Pins (`internal/hexagon/application/` bzw. ein dediziertes
  `internal/e2e/`-Verzeichnis). Ein auf das Adapter-Package
  begrenzter Scope würde die LH-AK-002-, LH-FA-UP-001- und
  LH-FA-UP-004-Pins formal als „nicht im getaggten Pfad" gelten
  lassen — Carveout-Aufhebung wäre nominell erfüllt, faktisch
  nicht. Die Test-Datei-Verortung jedes Pins ist in der Tabelle
  unten als eigene Spalte verbindlich; Reviewer prüfen vor dem
  Mergen, dass jede Datei tatsächlich existiert und unter
  `//go:build docker` läuft.
- `.github/workflows/ci.yml` bekommt einen neuen Job
  `integration-docker`, oder ein eigener Workflow
  `.github/workflows/integration.yml`. **`continue-on-error` ist
  zulässig, aber nur als zeitlich befristete
  Stabilisierungs-Maßnahme mit hartem Exit-Kriterium**: spätestens
  nach **drei aufeinanderfolgenden grünen Läufen** des neuen Jobs
  (auf `main`, ohne Re-Runs) wird `continue-on-error: false`
  gesetzt. Der Eintrag in `carveouts.md` darf erst entfernt
  werden, **nachdem** der Job ohne `continue-on-error` mindestens
  einmal grün durchgelaufen ist — sonst wäre die Carveout-
  Aufhebung selbst eine Hintertür (formell „erledigt", praktisch
  nicht-blockierend). Die drei-grünen-Läufe-Heuristik ist eine
  Carveout-Disziplin-Regel für diesen Slice; sie generalisiert
  nicht auf andere Slices, weil andere Adapter-Integrations-Pfade
  andere Failure-Modes haben.
- **Auditierbarkeits-Pflicht für die drei-grünen-Läufe-Regel**
  (mechanischer Beleg statt Vertrauenswort): die PR, die
  `continue-on-error: false` setzt, **muss** in ihrer Beschreibung
  einen Block der folgenden Form enthalten — sonst ist die PR
  nicht mergebar:

  ```
  ## Stabilization Evidence (slice-m6-docker-integrationstests)
  Three consecutive green `integration-docker` runs on `main`,
  each a first attempt (no re-runs, no workflow_dispatch):
  1. <commit-sha>  <workflow-run-url>  run_attempt=1  conclusion=success  event=push
  2. <commit-sha>  <workflow-run-url>  run_attempt=1  conclusion=success  event=push
  3. <commit-sha>  <workflow-run-url>  run_attempt=1  conclusion=success  event=push
  ```

  Reviewer-Pflicht (mechanisch, nicht „Vertrauen"):
  1. Jede URL öffnen und auf der GitHub-Run-Seite verifizieren,
     dass **`run_attempt`** = `1` ist — sichtbar im URL-Pfad
     (`/attempts/1`) oder am UI-Element „Run #N · Attempt #1".
     Re-Runs zeigen `/attempts/2+` und sind disqualifiziert,
     selbst wenn der Trigger ursprünglich `push` war.
  2. **`event`** = `push` bestätigen (nicht
     `workflow_dispatch`/`schedule`/`repository_dispatch`).
  3. **Commit-SHA** auf `main` verifizieren (nicht auf einem
     temporär gepushten Branch, der später force-gelöscht wurde
     — GitHub würde den Run trotzdem zeigen).
  4. Optional, aber empfohlen: per `gh api repos/{repo}/actions/runs/{run-id}`
     die drei Felder `run_attempt`, `event` und `head_branch`
     gegenprüfen; ein Mismatch zwischen UI und API-Antwort wäre
     ein Manipulations-Signal.

  Die PR, die später den `carveouts.md`-Eintrag entfernt, **muss**
  zusätzlich auf die URL des **ersten** Laufs **ohne**
  `continue-on-error` verweisen — auch mit `run_attempt=1`-Pin —
  damit die Aufhebung selbst an einem konkreten Commit-SHA hängt,
  nicht an einer Vertrauens-Aussage. Mit diesem Audit-Pfad bleibt
  die Regel ohne neue Repo-Files / Scripting mechanisch
  nachprüfbar; Reviewer und künftige Auditoren haben die
  Belegkette in der PR-History selbst.
- `docs/user/quality.md` §2 Tests wird um den Docker-Pfad ergänzt.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf
  den Aufhebungs-Commit als gelöst markieren.

### Verhaltens-Pins (Spec-IDs, die der getaggte Pfad einlösen muss)

Diese Pins decken Verhalten ab, das der **Unit-Pfad in T2 bewusst
nicht** testet (die T2-Substitute `/bin/echo` und `/bin/false`
faken nur Exit-Codes, kein echtes Compose-Stderr-Volumen oder
-Format). Damit klare Trennung: Unit deckt Klassifikation
ab, getaggter Integration-Pfad deckt End-to-end-Verhalten ab.

| Spec-ID                | Verhalten, das pin-getestet wird                                                                                                                                          | Test-Datei-Verortung (verbindlich)                                                       | Vorgeschlagener Test-Pfad                                                                                                            |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **LH-AK-002**          | Voller PostgreSQL-Flow: `u-boot init && u-boot add postgres && u-boot up` ⇒ Healthcheck `healthy` in ≤60 s, Port 5432 erreichbar.                                          | `internal/e2e/postgres_acceptance_docker_test.go` (neues Verzeichnis, package `e2e_test`) | End-to-end-Test mit echtem temp-Verzeichnis + InitProjectService + AddServiceService + UpService gegen echte Engine.                 |
| **LH-FA-UP-001 §966**  | Service mit Healthcheck stabilisiert erst auf `healthy`, nicht auf `running` allein.                                                                                       | `internal/hexagon/application/upservice_healthcheck_docker_test.go`                       | Compose-Fixture mit langsamem `pg_isready`-Healthcheck; assert dass `up` nicht früh returnt.                                          |
| **LH-FA-UP-001 §968**  | TCP-Port-Probe gegen `localhost` greift sobald Healthcheck `healthy` ist.                                                                                                  | `internal/hexagon/application/upservice_portprobe_docker_test.go`                         | Postgres-Fixture; assert `net.DialTCP("127.0.0.1:5432")` wird vom UpService durchgeführt und `nil`-error returnt.                     |
| **LH-NFA-PERF-002**    | **Compose-Stderr-Phasen-Stream** (`Pulling…`, `Creating…`, `Starting…`, `Healthchecking…`) reicht **live** an `opts.ProgressSink` durch — kein Buffer, kein Verlust bei fehlschlagendem `up`. | `internal/adapter/driven/docker/engine_progressstream_docker_test.go`                     | Pin via **`io.Pipe`-Event-Ordering** statt absoluter Timing-Schwellen. Setup: `r, w := io.Pipe()`; `opts.ProgressSink = w`; eine Test-Goroutine liest aus `r` und sammelt `(chunk, recvAt time.Time)`-Events bis EOF. Test-Assertion ist rein ereignis-relational, **keine Wall-Clock-Zahlen**: (a) mindestens ein nicht-leerer Chunk empfangen; (b) `events[0].recvAt < composeUpReturnedAt` — das **erste** Chunk-Event muss **vor** dem `ComposeUp`-Return-Zeitpunkt liegen (`composeUpReturnedAt := time.Now()` direkt nach dem `ComposeUp`-Aufruf in der Test-Goroutine). Damit ist „live" als **happens-before**-Relation operational definiert: ein nachträglicher Buffer-Flush würde alle Events erst **nach** dem Return zustellen, wodurch (b) reißt — unabhängig davon, ob der CI-Runner schnell (Pull dauert 2 s) oder langsam (Pull dauert 30 s) ist; relative Reihenfolge bleibt invariant. Fixture: Service mit `image: postgres:16-alpine` (ergo echter Pull mit garantierter `Pulling…`-stderr) + forcibly-fehlschlagender Healthcheck-Definition, damit `up` mit Code 12 endet — der Failure-Pfad ist der härtere Test (Erfolgs-Pfad könnte alle Phasen erfolgreich buffern, der Fail-Pfad nicht). |
| **LH-FA-UP-004 §1015** | `compose down -v` löscht das postgres-data-Volume wirklich (nicht nur den Container).                                                                                      | `internal/e2e/down_volumes_docker_test.go`                                                | Fixture: schreibe Test-Daten ins Postgres-Volume, `down --volumes`, assert dass Volume auf Docker-Ebene weg ist.                      |
| **LH-FA-DIAG-002**     | `docker compose ps --format json` (v2.20+) liefert die im T2-Parser angenommenen Feldnamen (`Service`, `State`, `Health`, `Publishers`).                                  | `internal/adapter/driven/docker/engine_psjsonschema_docker_test.go`                       | Snapshot-Test gegen reale Compose-Ausgabe; bricht laut bei Compose-Schema-Drift.                                                     |

**Audit-Check für die Aufhebungs-PR**: Reviewer prüft, dass jede der
sechs Test-Datei-Pfade aus der Spalte „Test-Datei-Verortung"
tatsächlich existiert, das `//go:build docker`-Tag trägt und im
neuen Job `integration-docker` mitläuft (sichtbar an der Test-
Output-Sektion des CI-Logs). Eine Datei fehlt = Pin nicht
eingelöst = Carveout-Aufhebung blockiert, unabhängig vom
`continue-on-error`-Status oder der „3 grüne Läufe"-Heuristik.

Die Verhaltens-Pins sind die eigentliche Aufhebungs-Bedingung —
ohne sie wäre der getaggte Pfad nur ein Smoke-Test ohne Spec-
Verankerung. Die Tabelle ist die kanonische Quelle dafür, welche
Spec-IDs **nicht** im T2-Unit-Pfad eingelöst werden.

## Out of Scope

- Andere Build-Tags (`//go:build keycloak`, `//go:build otel`) —
  separate Slices pro Adapter.
- Kubernetes-Smoke — u-boot orchestriert Compose-Stacks, nicht
  Kubernetes; ein Cluster-Smoke-Pfad ist nicht im Roadmap-Bereich.
- **Shell-Script-Mock-basierte Stderr-Pins.** Ein
  per-args-branching-Script über `t.TempDir()` als Engine-Binary
  wäre die einzige Möglichkeit, LH-NFA-PERF-002 ohne reale
  Docker-Engine zu pinnen — aber das würde nur die Sink-
  **Verkabelung** prüfen, nicht den realen Compose-Stderr-Output.
  Echte Verhaltensvalidierung gehört in den getaggten
  Integration-Pfad; ein Shell-Mock wäre danach redundant. T2-
  Review (`84a676c`) hat die Option explizit erwogen und
  zugunsten dieser Klärung verworfen.

## Bezug

- Auslösende Spec: `spec/architecture.md` §5 Build-Tag-Konvention.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Docker-Integrationstests fehlen.
- Hängt von: M6 `u-boot up`/`down` — erst dort entsteht
  `internal/adapter/driven/docker/`, gegen das die getaggten Tests
  überhaupt laufen können. **Adapter-Stand**: `internal/adapter/
  driven/docker/engine.go` existiert seit M6-T2 (`84a676c`); das
  Skeleton `engine_docker_test.go` mit `//go:build docker` ebenfalls.
- Vorgelagerte Reviews: M6-T2 Review identifizierte
  LH-NFA-PERF-002-Pin als Mittel-Lücke; verortet hier statt im
  T2-Unit-Pfad (siehe Out-of-Scope-Begründung oben).
- Phase: M6 (zusammen mit dem Docker-Adapter), Aufhebung
  spätestens vor M6 Done.
