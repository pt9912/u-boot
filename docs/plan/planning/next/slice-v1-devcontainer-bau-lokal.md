# Slice V1: Bauen im Devcontainer — lokaler Modus für die Pflicht-Targets mit gepinnter Toolchain, Docker-Modus bleibt Standard für Host und CI

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist.

**Bezug:** [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) (Host-Abhängigkeiten bleiben Docker und `make`), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (der Devcontainer, in dem gebaut wird), [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Entscheidung; Status `Accepted`, 2026-10-31).

**Berührte Spec-Stellen:** [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (neuer lokaler Modus), [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (Spalte *Zweck* der Pflicht-Targets), [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) (die Stages rufen das gemeinsame Skript), [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) (Coverage-Lauf, falls Test und Coverage ein Lauf werden).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Im Devcontainer baut, lintet und testet `make` mit der dort installierten, gepinnten Toolchain und den Caches des Containers; am Host und in der CI läuft unverändert der Docker-Modus.

**Befund** (gelesen und gemessen am 2026-10-31, Fundstellen in [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) §Kontext; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- `repo.mk` Zeilen 77 bis 79, 100, 103, 142: `--no-cache-filter` lässt `lint`, `test` und `coverage` bei jedem Lauf neu laufen. Das `Dockerfile` hat keinen `--mount=type=cache`; `test` (Zeile 82) und `coverage` (Zeilen 138 bis 143) führen beide die ganze Suite aus.
- Gemessen: `make lint` 20 bis 31 s, `make test` 18 bis 20 s, `make coverage-gate` 19 bis 20 s, ein Lauf je Zelle, mit und ohne Quelländerung gleich.
- Der Architect hat den Plan in frischem Kontext gegen [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md), [ADR-0004](../../adr/0004-ci-system.md), das Referenzmodell und [`AGENTS.md`](../../../../AGENTS.md) geprüft: 0 HIGH, 2 MEDIUM (F-3, F-4), kein Widerspruch, keine Folge-ADR heute ([Report](../../../reviews/2026-10-31-slice-v1-devcontainer-bau-lokal-architect.md)). Seine Befunde sind in §3 und §6 eingearbeitet.
- `.claude/hooks/pretooluse-command-guard.sh` blockt `go` in der Befehlsposition. Der Slice ändert den Guard nicht: Der lokale Modus läuft über `make`, nicht über direkte `go`-Aufrufe des Agenten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die CI und ein gemeinsames Toolchain-Image (Option E der ADR).** Der Projektinhaber hat entschieden, zuerst den Devcontainer-Modus einzuführen; die CI bleibt der Referenzlauf im Docker-Modus. Ein Folge-Slice dazu ist nicht angelegt, weil die Entscheidung nach diesem Slice fällt.
- **Ein Abgleich der Pins über `Dockerfile`, `repo.mk` und `.devcontainer/Dockerfile`.** Der Projektinhaber will kein Skript im Repo; das Gate liefert d-check (Issue 14 in `pt9912/d-check`). Bis dahin ist das Risiko in §6 benannt.
- **Feature-Optionen in `u-boot.yaml` und ein Pin des Basis-Images.** Das ist eine Änderung an `u-boot` selbst und ein anderer Vorgang. Der offene Slice [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md) nennt Feature-Optionen als einen möglichen Weg für ein verwandtes Problem; er nimmt diesen Punkt nicht an.
- **Mutationstests und Änderungen am Hook-Guard.** Mutationstests gibt es nicht und sie sind ein eigener Vorgang (Sensor außerhalb von `make gates`). Der Guard gehört `ai-harness-init` ([`MR-012`](../../../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md)).
- **`build` (Runtime-Image), `run`, `image-scan` und `packages`.** Sie arbeiten am Image selbst; ihr Docker-Pfad bleibt ([ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 3).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Modus verkörpert.** [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) nennen beide Modi; `repo.mk` schaltet über `BUILD_MODE` (Standard `docker`); `make modus-pruefung` belegt je Modus, welche Targets `docker build` rufen.
- [ ] **Devcontainer trägt die Toolchain.** `.devcontainer/Dockerfile` installiert Go und `golangci-lint` in den Versionen von `repo.mk` mit SHA-256-Prüfung und setzt `BUILD_MODE=local` und `GOCACHE`; `make gates` im Devcontainer ist grün; die gemessenen Zeiten vor und nach der Änderung stehen in §7.
- [ ] **Regeln und Doku nachgezogen.** [`AGENTS.md`](../../../../AGENTS.md) §3.1, [`docs/maintainer/quality.md`](../../../maintainer/quality.md) und [`harness/README.md`](../../../../harness/README.md) nennen den Nachweis „`make gates` im Modus der Umgebung“.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register ([`observations/`](../observations/)) fortgeschrieben, oder in §7 „keine Beobachtung angefallen“.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Schritt M0 (vor jedem Code): Messlauf.** Klärt, ob `go test` mit Coverage-Profil Ergebnisse cached, und misst `make gates` im Docker-Modus als Ausgangswert (Wanduhr, der Ablauf ist parallel). Das Ergebnis legt fest, ob Test und Coverage im lokalen Modus ein Lauf werden ([ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 7).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/README.md` | unverändert | Gelesen: Zeile 65 trägt die Indexzeile für [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Commit `b371e262`); der Architect-Report F-6 berichtigt die frühere Planzeile „update“ |
| `spec/spezifikation.md` | update | [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Zeilen 368 bis 374) um den lokalen Modus erweitern; [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (Zeilen 1146 bis 1168): Spalte *Zweck* nennt `docker build --target …` (Docker-Modus) und den lokalen Befehl; Zeile 99 (CI-Runner, „Keine Host-Go-Toolchain“) bleibt, weil die CI im Docker-Modus läuft. [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) sagt zusätzlich: Der lokale Modus gilt für alle Make-Targets mit Docker-Pfad, nicht nur für die in [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) geführten (`test-docker`, `govulncheck` und `build-binaries` stehen dort nicht); Release-Artefakte (`build-binaries`, `packages`) entstehen nur im Docker-Modus; ein Teillauf (`PKG`, `RUN`) ist kein Nachweis. [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) (Zeilen 58 bis 71) und [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleiben unverändert, gelesen; werden Test und Coverage im lokalen Modus ein Lauf, schärft eine Folge-ADR [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) (Architect-Report F-2) |
| `scripts/bau.sh` | neu | Stage-Befehle `compile`, `lint`, `test`, `coverage`: eine Quelle für `Dockerfile` und lokalen Modus; setzt je Stage nur, was die `RUN`-Zeile im `Dockerfile` heute selbst setzt: `CGO_ENABLED=0` bei `compile`, `test` und `coverage` (Zeilen 62, 82, 138); `GOFLAGS` kommt aus der Umgebung (Docker: `ENV` der Stage `deps`, Zeilen 42 bis 46; lokal: `ENV` im `.devcontainer/Dockerfile`). `lint` setzt nichts: Die Stage baut auf `golangci/golangci-lint`, nicht auf `deps`, und ruft `golangci-lint run ./...` ohne diese Variablen (Zeilen 71 bis 76); so bleibt der Docker-Modus unverändert (Architect-Report F-3). M0 prüft vorab die Go-Version im Image `golangci/golangci-lint:v2.14.0` gegen `GO_VERSION` und ob `bash` im Lint-Image vorhanden ist (nicht geprüft; das Skript muss sonst für `lint` POSIX-`sh` sein) |
| `scripts/coverage-gate.sh` | unverändert | Kopf gelesen (Aufruf `coverage-gate.sh <coverage-func.txt> <threshold>`, Bootstrap-Modus über `COVERAGE_BOOTSTRAP`); der Aufruf aus `scripts/bau.sh` ist derselbe wie aus `Dockerfile` Zeile 146 |
| `Dockerfile` | update | Die `RUN`-Zeilen der Stages `compile`, `test`, `coverage` (Zeilen 62, 82, 129 bis 146) rufen `scripts/bau.sh`; Stage `lint` (Zeile 76) ebenso; die Stage-Struktur nach [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) bleibt |
| `repo.mk` | update | `BUILD_MODE ?= docker`; lokale Rezepte für `deps`, `compile`, `lint`, `test`, `coverage-gate`, `test-docker`, `govulncheck`, `build-binaries`; `PKG` und `RUN` für `test`; `modus-pruefung`; `--no-cache-filter` (Zeilen 77 bis 79) bleibt für den Docker-Modus. `packages` hängt von `build-binaries` ab (Zeile 167): Im lokalen Modus bleibt `packages` im Docker-Pfad, `build-binaries` lokal ist nur Entwicklungs-Komfort. Die Rezepte von `gates`-Targets setzen `PKG` und `RUN` leer oder brechen ab, wenn sie gesetzt sind. Kopplungskommentar bei `GO_VERSION` und `GOLANGCI_LINT_VERSION` (Zeilen 43 und 44): „Hebung betrifft drei Stellen: `Dockerfile`, diese Datei, `.devcontainer/Dockerfile`“ (Architect-Report F-4, F-5) |
| `.devcontainer/Dockerfile` | update | Hinter dem verwalteten Block: Go und `golangci-lint` als Release-Archive mit SHA-256 (je Architektur), `ENV BUILD_MODE=local`, `ENV GOCACHE` im Workspace-Volume; die Zeilen überstehen `u-boot generate devcontainer` (gemessen am 2026-10-31). Der verwaltete Block endet mit `USER vscode`; die Installation braucht davor `USER root` und danach wieder `USER vscode`. Der Kopplungshinweis „Hebung betrifft drei Stellen“ steht auch hier (Architect-Report F-4) |
| `.gitignore`, `.dockerignore` | update | Der Go-Cache liegt unter `.cache/` im Workspace (`GOCACHE=/workspaces/u-boot/.cache/go-build`); der Eintrag `.cache/` gehört in beide Dateien: nicht in Git, nicht in den Docker-Kontext (kein Eintrag dafür heute, Architect-Report F-7) |
| `spec/architecture.md` | update | Zeilen 714 und 716 beschreiben `make test` als Docker-Stage und ein getrenntes Make-Target in der CI (gelesen über den Architect-Report); sie erhalten den Zusatz „im Docker-Modus“; Schichten und Importregeln bleiben unberührt, weil der Slice keinen Go-Quelltext ändert (Architect-Report F-7) |
| `scripts/build-packages.sh` | unverändert | Gelesen: Kopf und Aufruf (Zeilen 3 bis 4, 29): `nfpm` läuft per `docker run`; der Slice ändert das nicht (Architect-Report F-5) |
| `tools/mk/devcontainer.mk` | unverändert | Gelesen: die Targets `dev-up`, `dev-down`, `dev-shell` rufen nur `u-boot devcontainer`; sie kennen keinen Modus |
| `AGENTS.md` | update | §3.1: der Nachweis ist `make gates` im Modus der Umgebung; die Begründung (Reproduzierbarkeit, Supply-Chain) bleibt; Bezug auf [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) als Plan-Anker nach §3.6 |
| `docs/maintainer/quality.md` | update | Zeilen 20 bis 31 (Lint über `docker build`), Zeile 72 (`make gates`), Zeilen 300 bis 314 (CI: „Keine Host-Toolchain“ gilt für die CI unverändert) |
| `harness/README.md` | update | Sensors-Tabelle: Zeile für `make modus-pruefung` (kein Gate); Satz zum Modus unter *Sensors* |
| `harness/verification.md` | unverändert | Gelesen (Zeilen 75 bis 87): die Sensor-Tabelle nennt Make-Targets ohne Modus |
| `.github/workflows/ci.yml`, `.github/workflows/integration.yml` | unverändert | Gelesen: sie rufen `make gates`, `make govulncheck`, `make image-scan` und `make test-docker` (ohne `BUILD_MODE`, also Docker) |
| `CHANGELOG.md` | unverändert | Kein öffentlicher Vertrag berührt: Produktverhalten und Release-Artefakte ändern sich nicht (Kopf gelesen, Abschnitt `Unreleased` leer) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): Der Projektinhaber priorisiert den Slice, `Verantwortlich:` ist gesetzt; kein anderer Slice steht in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn M0 ergibt, dass der lokale Modus ohne Umbau des Coverage-Laufs nicht trägt, oder wenn `repo.mk` mehr als die genannten Targets umbauen müsste.
- `in-progress` → `open` (blockiert): wenn `golangci-lint` oder Go im Devcontainer nicht ohne ein Feature mit Lockerungen installierbar sind.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD ist vollständig, `make gates` im Devcontainer ist grün, die Zeiten vor und nach der Änderung stehen in §7, ein Pull Request ist gemergt und die CI (Docker-Modus) ist grün, und die Closure-Notiz mit dem Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Modi laufen auseinander** (ein Lauf lokal grün, in der CI rot), weil `scripts/bau.sh` nur Teile der Stages teilt — **Ausgang:** bei Closure.
- **Die Pins stehen an drei Stellen** (`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`) ohne Abgleich, bis d-check das Modul liefert — **Ausgang:** bei Closure (weiter offen: Beobachtung im Register, Verweis auf Issue 14 in `pt9912/d-check`).
- **Ein Teillauf zählt als Nachweis** (`make test PKG=… RUN=…`): Die Rezepte von `gates`-Targets lehnen `PKG` und `RUN` ab oder setzen sie leer (§3, `repo.mk`); `make modus-pruefung` prüft das — **Ausgang:** bei Closure.
- **`govulncheck` im lokalen Modus lädt sein Werkzeug über das Netz** (`go install …@$(GOVULNCHECK_VERSION)`, Modul-Proxy statt SHA-256-Archiv); im Docker-Modus ist die Lage gleich (Version gepinnt, kein Hash), also kein neuer Befund und ein akzeptiertes Negativ (Architect-Report F-4) — **Ausgang:** bei Closure.
- **Test-Cache verdeckt einen umgebungsabhängigen Test** im lokalen Modus — **Ausgang:** bei Closure (die CI im Docker-Modus läuft ohne Cache und fängt es).
- **`node:2.1.1` mit `claude-code:1.0.5`** ist nicht gemessen (das Handbuch hat `claude-code` mit `node:1` gemessen); der Devcontainer-Build kann scheitern — **Ausgang:** bei Closure.
- **Der Hook-Guard blockt `go` in der Befehlsposition**, auch im Devcontainer; ein Agent, der `go test` direkt ruft, wird gestoppt — **Ausgang:** bei Closure.
- **`GOCACHE` im Workspace-Volume** wächst ohne Grenze — **Ausgang:** bei Closure.

## 7. Closure-Notiz

<!-- Wird vor dem `git mv` nach done/ gefüllt (Planner-Rolle). -->

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `tools/`, `scripts/` (WZG), `spec/`, `harness/`, `docs/plan/` (DOK). Beide sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md) als Sub-Areas geführt; eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist durchgegangen. Treffer: [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Schwelle erreicht, kein Ausgang zugewiesen, Entscheidung des Projektinhabers zum Stolperdraht steht aus). Der Slice berührt den Gegenstand: Im Devcontainer ist `go` über `make` erlaubt, der Guard bleibt unverändert, und ein direkter Aufruf bleibt blockiert (§6).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
