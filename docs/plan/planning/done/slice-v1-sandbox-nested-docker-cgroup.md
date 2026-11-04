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

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer, Degradationstabelle), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Dateien). Entscheidung: [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md).

**Berührte Spec-Stellen:** [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Statusdateien der Startscripte), [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Grenzen; Wirkung auf den verschachtelten Docker), [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) (neu, Träger [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)); kein Lastenheft-CR ([`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) Entscheidung 5).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Sandbox-Devcontainer mit `nestedRuntime: docker` baut beim Start seine cgroup so um, dass der verschachtelte Docker Grenzen (`--memory`, `--cpus`) anwenden kann, bevor `dockerd` läuft — die Schritte (1) bis (4) des Befunds, idempotent, nur bei cgroup v2 und nur, wenn `memory` nicht schon an die Kinder weitergegeben wird. Ein Fehlschlag ist best-effort: der Umbau bricht ab, eine Zeile auf stderr nennt den Schritt, und `dockerd` startet wie ohne Umbau; Statusdatei und Exit-Code ändern sich nicht, die Policy `onUnavailable` greift nicht ([`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) Entscheidung 5, [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Keine Änderung an `nestedRuntime: podman` und am Standardprofil ohne nested Docker:** der Befund betrifft den Daemon, den `docker-init.sh` startet.
- **Keine Host-Einstellung und kein `--cgroupns=host`:** das hebt die Isolation auf und verschiebt die Frage nur.
- **Keine Zusage für andere Hosts oder Engines:** gemessen ist ein Container, nicht die Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Linux `amd64`, lokale Engine). Die Messungen 1, 2, 4 und 5 bleiben bei [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md).
- **Kein Umbau des Werteformats** der Grenzen.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [x] **Entscheidung.** [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) (`Accepted`, Index nachgezogen) wählt den Startort (`postStartCommand`, nicht `ENTRYPOINT`; gelesen an der Dev Container CLI 0.89.0); kein Lastenheft-CR (Entscheidung 5, Lastenheft unverändert). *Geändert auf Entscheidung des Projektinhabers 2026-11-04: Das Messprotokoll unter `docs/messungen/slice-v1-sandbox-nested-docker-cgroup/` und die Messung der Risiken 2 bis 4 (Startort-Wirkung nach Neuanlage, `docker exec` und neue Terminals, Prozesse der Erweiterung) entfallen für diesen Slice und wandern an [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md); die lesende Quellenlage zu Risiko 1 steht in der ADR.*
- [x] **Erzeugung.** Das Template startet den Umbau vor `dockerd` (Ort laut ADR); vier Unit-Tests für Ort, Bedingung und Schritte, Fehlschlag und „nur mit nested Docker“ (mit und ohne nested Docker) samt Zeile in `harness/replay.md`; ein Docker-Test (Build-Tag `docker`) belegt `memory.max` und `cpu.max` im verschachtelten Container eines Wegwerf-Containers (ein Lauf in der verschachtelten Engine dieses Entwicklungs-Devcontainers); das eigene Projekt ist regeneriert. *Geändert auf Entscheidung des Projektinhabers 2026-11-04: der Satz „die zwei roten Tests laufen in einem frisch angelegten Container grün“ entfällt hier und wandert an den Folge-Slice. Dry-Run: kein eigener Fall für `docker-init.sh`; die allgemeinen Dry-Run-Tests des Generators (`generate_rollback_test.go`, `generate_envfile_test.go`, `dryrun_log_test.go`) laufen im grünen `make test`, pinnen aber nicht diese Datei (Verification V-6, Entscheidung des Planners).*
- [x] **Doku.** Handbuch (`docs/user/`) und `CHANGELOG.md` nennen nur das Gemessene; Nutzer-Doku ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); `docs/maintainer/devcontainer-arbeiten.md` nennt den Umbau.
- [x] `make gates` grün im Devcontainer-Container (Verifier-Lauf auf `e6083d7d`, Exit 0). *(nicht gezählt — Gate-Lauf)*
- [x] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge, Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Fundstellen am 2026-11-04 gelesen; die Tabelle ist nach der letzten Fix-Runde (HEAD `e6083d7d`) gegen `git diff --stat 60e75605..HEAD` ohne den Planungs-Baum abgeglichen (11 Dateien, jede mit Zeile) und ohne den Commit `6ea73133` (gehört zu [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers); seine Änderungen an `CHANGELOG.md`, `docs/user/benutzerhandbuch.md` und `spec/spezifikation.md` stehen nicht in dieser Tabelle) ([`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl` | update | startet `dockerd` per `setsid nohup` nach der Bereitschaftsprüfung; der Umbau steht davor im selben Zweig, best-effort ([`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md), [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag)) |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | unverändert (kein `ENTRYPOINT`, [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)) | Block `{{- if .NestedDocker}}` (Zeile 53) kopiert `docker-init.sh` nach `/usr/local/bin/u-boot-docker-init` (Zeile 65) |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` | unverändert (kein `overrideCommand`, [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)) | `postStartCommand` kommt aus `data.PostStart` (`devcontainer_sandbox.go`) |
| `internal/hexagon/application/devcontainer_sandbox.go` | unverändert (kein neuer Pfad, [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)) | `devcontainerTemplateData`: Zweig `nestedDocker(dc)` setzt `PostStart`, `Mounts`, `RunArgs` |
| `internal/hexagon/application/generate.go` | prüfen | `planDevcontainerFiles` (Zeile 1071): `docker-init.sh` ist bereits im Plan |
| `internal/hexagon/application/generate_sandbox_test.go` | update | vier neue Tests der Beleg-Tabelle von [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) (`…_CgroupPrepBeforeDaemon`, `…_CgroupPrepConditionsAndSteps`, `…_CgroupPrepBestEffort`, `TestGenerateDevcontainer_CgroupPrepOnlyWithNestedDocker`), Hilfen `generatedDockerInit`, `funcBody`; `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` um das Fragment der stderr-Zeile erweitert (Zählungen unverändert); `TestGenerateDevcontainer_SandboxDocker`, `…_ScriptConflict`, `…_StaleBlockKeepsUserContent` unverändert, sie pinnen den Wortlaut nicht (gelesen: Fragmente `wait_seconds=…`, `dockerd --host=…`, `exit 11`) |
| `internal/e2e/sandbox_devcontainer_docker_test.go` | update | `TestE2E_SandboxDevcontainer_NestedDocker_CgroupLimitsReachNestedContainer` (Build-Tag `docker`); Fix-Runde: Vorbedingung vor dem Skript (`memory` schon in `subtree_control` → Skip), Fehlschlagzeile des ersten Laufs ist ein Fehler, `memory` in `subtree_control` nach dem ersten Lauf |
| `internal/e2e/resources_docker_test.go` | unverändert | `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` prüft nur `HostConfig` des äußeren Containers (Kopfkommentar „Not checked: the effect on the nested Docker“) und belegt den Umbau nicht ([`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) führt ihn nicht in der Beleg-Tabelle); Kopfkommentar (Zeilen 3 bis 10) nennt den Ausgangszustand der verschachtelten Engine |
| `.devcontainer/docker-init.sh` dieses Repos | update | eigenes Erzeugnis von `u-boot generate devcontainer` (zweiter Lauf: `no-op`); `.devcontainer/Dockerfile`, `devcontainer.json` und `.gitignore` unverändert (der Lauf änderte sie nicht) |
| `harness/replay.md` | update | Zeile `generate devcontainer` mit `nestedRuntime: docker` um den cgroup-Umbau und den Test erweitert |
| `spec/spezifikation.md` ([`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) neu, [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) Verweis, Historie), `docs/plan/adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md` samt Index | neu / update | Entscheidung des Architekten; `spec/lastenheft.md` unverändert (kein CR, [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) Entscheidung 5) |
| `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/devcontainer-arbeiten.md` | update | Doku; READMEs unverändert (kein Satz dort betrifft den Umbau) |
| `docs/messungen/slice-v1-sandbox-nested-docker-cgroup/` | entfällt | kein Messprotokoll in diesem Slice (Entscheidung des Projektinhabers 2026-11-04); das Protokoll entsteht im Folge-Slice [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) |
| `docs/reviews/2026-11-04-slice-v1-sandbox-nested-docker-cgroup.md`, `docs/reviews/verification-2026-11-04-slice-v1-sandbox-nested-docker-cgroup.md`, Planungs-Baum (diese Datei, Folge-Slice, `open/README.md`, `in-progress/roadmap.md`) und Register (`observations/BEO-DOK/…`) | neu / update | Closure: Review- und Verification-Report, Folge-Slice mit Index-Zeile, Roadmap-Ruhe-Marker, Belege im Register |

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber priorisiert (2026-11-04); kein anderer Slice in `in-progress/` (WIP-Limit 1). Die Neuanlage des Containers mit `make dev-up RECREATE=1` ist nicht mehr Start-Bedingung dieses Slice; sie gehört zum Start des Folge-Slice.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Greift weder `ENTRYPOINT` noch `postStartCommand` vor `dockerd`, oder verhindert der Umbau neue Prozesse in der Wurzel (Risiko 3), gilt die Alternative: Grenzen im verschachtelten Docker bleiben nicht zugesagt, Messung auf einem Runner bei [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md); die Prüfung des Risikos 3 liegt beim Folge-Slice.
- `in-progress` → `next` (zu groß): Messung und Entscheidung als eigener Slice vor dem Code.

## 5. Closure-Trigger

DoD vollständig; Review und Verifikation durch andere Rollen; Closure-Notiz; `git mv` nach `done/`. *Geändert auf Entscheidung des Projektinhabers 2026-11-04: die Bedingungen „Neuanlage des Devcontainers belegt“ und „`make test-docker` ohne die zwei cgroup-Fehlschläge“ entfallen hier und sind Closure-Trigger von [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md).*

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

1. **Der `ENTRYPOINT` des Images läuft nicht.** PID 1 dieses Containers ist `/bin/sh -c "echo Container started …"`; das deutet auf `--entrypoint /bin/sh` der Dev-Containers-Erweiterung (Eigenschaft `overrideCommand`). — **Ausgang:** entfallen, lesend belegt (2026-11-04): die Dev Container CLI 0.89.0 übergibt `--entrypoint /bin/sh` und lässt `ENTRYPOINT`/`Cmd` des Images nur bei `overrideCommand: false` stehen (Quelle: [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)); Startort ist `postStartCommand`. Einschränkung: die Erweiterung des Editors selbst ist nicht gelesen.
2. **Der Zustand überlebt keinen Neustart / die Neuanlage setzt ihn nicht zurück.** Gemessen ist nur der laufende Container. — **Ausgang:** weiter offen → [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) (Messung nach `make dev-up RECREATE=1`, nach dem Release `v0.23.0`; Entscheidung des Projektinhabers 2026-11-04).
3. **Neue Prozesse in der Wurzel nach dem Umbau** (`docker exec`, neue Terminals der Erweiterung): ein Domain-cgroup mit gesetztem `subtree_control` nimmt nach Kenntnis des Kernels keine Prozesse mehr auf; der Docker-Test prüft ein `docker exec` in einer Engine, nicht die Erweiterung. — **Ausgang:** weiter offen → [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md).
4. **Der Umbau verschiebt Prozesse, die die Erweiterung braucht** (z. B. Server-Prozesse). — **Ausgang:** weiter offen → [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) (Anhängen nach dem Umbau).
5. **Die Unit-Tests pinnen Text und Reihenfolge der erzeugten Datei, nicht das Verhalten des Skripts** (Review F-2, Verification V-4): „sonst geschieht nichts, ohne Ausgabe“ und „Fehlschlag startet `dockerd` trotzdem“ sind nur als Skript-Form gepinnt; der Docker-Test deckt den Normalzustand, nicht `domain threaded`. — **Ausgang:** weiter offen → Beobachtung [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Beleg `evidence/slice-v1-sandbox-nested-docker-cgroup.md`).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Folge ADR → Spezifikation → Template → Tests → Docker-Test hielt: [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) und [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) stehen vor dem Code, alle sieben Belege der Beleg-Tabelle haben einen existierenden Test (Verification). `make gates` ist grün (Verifier-Lauf auf `e6083d7d`, Exit 0, Coverage 92,10 %). Der Docker-Test `TestE2E_SandboxDevcontainer_NestedDocker_CgroupLimitsReachNestedContainer` ist PASS (nicht SKIP); die Mutationsprobe des Verifiers (Aufruf `cgroup_prepare` entfernt) färbt ihn rot — er unterscheidet den Umbau.
- **Grenze des Belegs (ehrlich):** Der PASS ist ein Lauf in der verschachtelten Engine dieses Entwicklungs-Devcontainers und unterscheidet den Umbau nur im Normalzustand (Domain-Wurzel mit Prozessen, `memory` nicht delegiert; Vorbedingung, sonst Skip). Den Zustand `domain threaded`, den die Messung der Sitzung zeigte, stellt er nicht her. Es gibt **keine Messung nach Neuanlage**; der CI-Lauf des Jobs `integration-docker` auf dem Runner steht aus. Die zwei `TestE2E_Resources_*` sind in diesem Container bis zur Neuanlage nicht grün und belegen den Umbau ohnehin nicht (sie prüfen `HostConfig` des äußeren Containers).
- **Was ging anders als geplant:** (1) Auf Entscheidung des Projektinhabers 2026-11-04 (Wortlaut: „dazu muss ich diese session beenden -- deshalb machen wir das erst nach dem Release“) sind die DoD-Zeilen *Messung und Entscheidung* und *Erzeugung* und §5 angepasst: Messung der Risiken 2 bis 4, das Messprotokoll und der Lauf der zwei Resources-Tests nach Neuanlage wandern an [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md), nach dem Release `v0.23.0`; die Anpassung steht auch an den geänderten DoD-Zeilen und in §5. (2) Der Review fand vier MEDIUM-Befunde (F-1 Beleg-Zeile mit Resources-Test, F-2 Bedingung nur als Text gepinnt, F-3 Spec enger als Code, F-4 Handbuch-Widerspruch), alle in der Fix-Runde behoben bis auf die Unit-Test-Grenze (Risiko 5). (3) ADR 0045 und [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) sind nach `Accepted` in place berichtigt worden, weil nichts veröffentlicht war (Verification V-3; nach dem Push nur noch neue ADR mit `Supersedes`). (4) Verification V-1: das Plan-Ziel nannte `onUnavailable` und eine Statusdatei, ADR und Code sagen das Gegenteil; Ziel und Zustandsmarker sind berichtigt. (5) Verification V-6: kein eigener Dry-Run-Fall für `docker-init.sh`; die DoD bindet an die allgemeinen Dry-Run-Tests (siehe DoD *Erzeugung*).
- **Steering-Loop-Eintrag:** Zwei Klassen gezählt, nicht verkörpert: Zusage weiter als die Messung (Beleg-Tabelle nannte einen Test ohne Bezug, Handbuch „kann Grenzen anwenden“ gegen „nicht zugesagt“) und Plan-Aussage über Code (Plan-Ziel gegen ADR und Code, Zustandsmarker).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-v1-sandbox-nested-docker-cgroup.md` in [`BEO-DOK/zusage-weiter-als-messung/`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (F-1, F-4), in [`BEO-DOK/plan-aussage-ueber-code/`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (V-1) und in [`BEO-DOK/spec-bedingung-ohne-code-entsprechung/`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (F-2, V-4) ergänzt; der Zähler folgt aus den Dateien, kein Ausgang zugewiesen.
- **Folge-Slices:** [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) ist eine Datei in `open/`.
- **Risiken aus §6:** je ein Ausgang, siehe §6 (1 entfallen, 2 bis 4 weiter offen im Folge-Slice, 5 weiter offen im Register).
- **Drei Paarungen:** Anker — kein Herkunfts-Anker gesetzt (nichts verkörpert); Folge-Slice — die Datei in `open/` liegt vor, die Index-Zeile steht in `open/README.md`; Register — die drei Evidence-Dateien liegen vor.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `docs/plan/`, `docs/messungen/` (DOK), `internal/hexagon/application` (APP) und `internal/e2e` (E2E); alle sind in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) deklariert, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register gesichtet (Closure 2026-11-04); Treffer: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (offen über der Schwelle, kein Ausgang zugewiesen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10); der Zähler ist je die Zahl der Evidence-Dateien. Dieser Slice trägt je einen weiteren Beleg bei (§7).

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
