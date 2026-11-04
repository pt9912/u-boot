# Slice V1: cgroup-Umbau beim Start des Sandbox-Devcontainers mit nested Docker — die Grenzen für Arbeitsspeicher und CPU erreichen verschachtelte Container, der Umbau läuft wie im dind-Entrypoint vor `dockerd`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung offen (keine Closure-Bedingung außerhalb der DoD).

**Anlass.** In diesem Devcontainer scheitert jeder Start mit `--memory` oder `--cpus` (Exit `125`, `unable to apply cgroup configuration: cannot enter cgroupv2 "/sys/fs/cgroup/docker" with domain controllers -- it is in threaded mode`); zwei Tests von `make test-docker` (`TestE2E_Resources_DefaultProfileLimitsReachTheContainer`, `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer`) sind deshalb rot, und die offenen Messungen von [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md) sind hier nicht fahrbar. Entscheidung des Projektinhabers 2026-11-04: die Hürde im Devcontainer überwinden, nach dem Muster des dind-Entrypoints.

**Befund der Sitzung 2026-11-04** (an diesem Container gelesen und gelaufen, nicht auf anderen Hosts):

- Die Wurzel-cgroup hatte den Typ `domain threaded`, `subtree_control` trug `cpuset cpu pids`, `memory` nicht; `docker/` und `docker/buildkit` waren `threaded`. Der Kernel nennt eine Domäne mit Prozessen und threaded Controllern in `subtree_control` so.
- Die Folge von Schritten *(1) `-cpuset -cpu -pids` aus `cgroup.subtree_control` der Wurzel, (2) Unter-cgroup `init` anlegen und alle Prozesse der Wurzel hineinschreiben, (3) alle Controller aus `cgroup.controllers` in `subtree_control` aktivieren, (4) `dockerd` starten* setzte die Wurzel auf `domain`; `docker run --memory=64m --cpus=0.5` zeigte danach `memory.max = 67108864` und `cpu.max = 50000 100000`. Ohne Schritt 1 (nur 2 und 3) blieb die Wurzel `domain threaded`, die Prozesse ließen sich nicht verschieben.
- Der Zustand lebt in der cgroup des Containers und endet mit ihr; die Änderung war eine Handarbeit mit `sudo` im laufenden Container. Ob ein Neustart ihn zurücksetzt, ist nicht gemessen (§6 Risiko 2).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer, Degradationstabelle), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Dateien). Die ADR-Nummer vergibt der Architekt-Schritt über den ADR-Index.

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Statusdateien der Startscripte), [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Grenzen; Wirkung auf den verschachtelten Docker), je nach Entscheidung ein neuer `SPEC-<NNN>` und ein Lastenheft-CR.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Sandbox-Devcontainer mit `nestedRuntime: docker` baut beim Start seine cgroup so um, dass der verschachtelte Docker Grenzen (`--memory`, `--cpus`) anwenden kann, bevor `dockerd` läuft — die Schritte (1) bis (4) des Befunds, idempotent, nur bei cgroup v2 und nur, wenn `memory` nicht schon an die Kinder weitergegeben wird. Ein Fehlschlag folgt der Policy `onUnavailable` und erscheint in der Statusdatei.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Änderung an `nestedRuntime: podman` und am Standardprofil ohne nested Docker:** der Befund betrifft den Daemon, den `docker-init.sh` startet.
- **Keine Host-Einstellung und kein `--cgroupns=host`:** das hebt die Isolation auf und verschiebt die Frage nur.
- **Keine Zusage für andere Hosts oder Engines:** gemessen ist ein Container, nicht die Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Linux `amd64`, lokale Engine). Die Messungen 1, 2, 4 und 5 bleiben bei [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md).
- **Kein Umbau des Werteformats** der Grenzen.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [ ] **Messung und Entscheidung.** Messprotokoll unter `docs/messungen/slice-v1-sandbox-nested-docker-cgroup/` zu den Risiken 1 bis 3 (Startort, Neustart, `docker exec`); ADR (`Accepted`, Index nachgezogen) zur Wahl des Startorts (`ENTRYPOINT` des Images oder `postStartCommand`); bei Wirkung auf die Zusage ein Lastenheft-CR mit Bestätigung des Projektinhabers vor dem Code.
- [ ] **Erzeugung.** Das Template startet den Umbau vor `dockerd` (Ort laut ADR); Golden-/Replay-Cases für beide Ausgaben (mit und ohne nested Docker, `--dry-run` schreibt nichts); die zwei roten Tests laufen in einem frisch angelegten Container grün; ein Docker-Test (Build-Tag `docker`) belegt `memory.max` und `cpu.max` im verschachtelten Container.
- [ ] **Doku.** Handbuch (`docs/user/`) und `CHANGELOG.md` nennen nur das Gemessene; Nutzer-Doku ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); `docs/maintainer/devcontainer-arbeiten.md` nennt den Umbau.
- [ ] `make gates` grün im Devcontainer-Container. *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge, Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Fundstellen am 2026-11-04 gelesen; der Implementer ergänzt die Tabelle im ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl` | update (Ort nach ADR) | startet `dockerd` per `setsid nohup` nach der Bereitschaftsprüfung; der Umbau steht davor |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update (nur bei `ENTRYPOINT`) | Block `{{- if .NestedDocker}}` (Zeile 53) kopiert `docker-init.sh` nach `/usr/local/bin/u-boot-docker-init` (Zeile 65) |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | update (nur bei `overrideCommand`) | `postStartCommand` kommt aus `data.PostStart` (`devcontainer_sandbox.go`) |
| `internal/hexagon/application/devcontainer_sandbox.go` | update (nur bei neuem Pfad) | `devcontainerTemplateData`: Zweig `nestedDocker(dc)` setzt `PostStart`, `Mounts`, `RunArgs` |
| `internal/hexagon/application/generate.go` | prüfen | `planDevcontainerFiles` (Zeile 1071): `docker-init.sh` ist bereits im Plan |
| Tests: `generate_sandbox_test.go`, Golden-Cases, `internal/e2e/sandbox_devcontainer_docker_test.go`, `internal/e2e/resources_docker_test.go` | update | Golden für beide Ausgaben; Docker-Test mit den Werten im verschachtelten Container |
| `.devcontainer/*` und `.gitignore` dieses Repos | update | eigenes Erzeugnis von `u-boot generate devcontainer` |
| `spec/lastenheft.md`, `spec/spezifikation.md`, ADR samt Index | update / neu | nach Messung und Entscheidung |
| `docs/messungen/slice-v1-sandbox-nested-docker-cgroup/`, `docs/user/*`, `docs/maintainer/devcontainer-arbeiten.md`, `CHANGELOG.md` | neu / update | Protokoll und Doku |

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber priorisiert; ein Container steht für das Anlegen mit `make dev-up RECREATE=1` zur Verfügung (Neuanlage verwirft diese Sitzung); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Greift weder `ENTRYPOINT` noch `postStartCommand` vor `dockerd`, oder verhindert der Umbau neue Prozesse in der Wurzel (Risiko 3), gilt die Alternative: Grenzen im verschachtelten Docker bleiben nicht zugesagt, Messung auf einem Runner bei [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md).
- `in-progress` → `next` (zu groß): Messung und Entscheidung als eigener Slice vor dem Code.

## 5. Closure-Trigger

DoD vollständig; Neuanlage des Devcontainers belegt; `make test-docker` ohne die zwei cgroup-Fehlschläge; Review und Verifikation durch andere Rollen; Closure-Notiz; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

1. **Der `ENTRYPOINT` des Images läuft nicht.** PID 1 dieses Containers ist `/bin/sh -c "echo Container started …"`; das deutet auf `--entrypoint /bin/sh` der Dev-Containers-Erweiterung (Eigenschaft `overrideCommand`). Ein `ENTRYPOINT` nach dind-Muster liefe dann nicht. Gelesen ist nur PID 1, nicht der Aufruf der Erweiterung. — **Ausgang:** offen → Messung im Slice.
2. **Der Zustand überlebt keinen Neustart / die Neuanlage setzt ihn nicht zurück.** Gemessen ist nur der laufende Container. — **Ausgang:** offen → Messung mit `make dev-up RECREATE=1`.
3. **Neue Prozesse in der Wurzel nach dem Umbau** (`docker exec`, neue Terminals der Erweiterung): ein Domain-cgroup mit gesetztem `subtree_control` nimmt nach Kenntnis des Kernels keine Prozesse mehr auf; ungemessen. — **Ausgang:** offen → Messung mit Anhängen der Erweiterung und `docker exec` nach dem Umbau.
4. **Der Umbau verschiebt Prozesse, die die Erweiterung braucht** (z. B. Server-Prozesse). — **Ausgang:** offen → Messung (Anhängen nach dem Umbau).

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `docs/plan/`, `docs/messungen/` (DOK), `internal/hexagon/application` (APP) und `internal/e2e` (E2E); alle sind in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) deklariert, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** nicht durchgegangen; der Implementer sichtet das Register im ersten Lauf, bis dahin steht hier kein Befund.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
