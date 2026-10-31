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

**Verantwortlich:** pt9912 (gesetzt beim Start, 2026-10-31)

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

- [x] **Modus verkörpert.** [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) nennen beide Modi; `repo.mk` schaltet über `BUILD_MODE` (Standard `docker`); `make modus-pruefung` belegt je Modus, welche Targets `docker build` rufen.
- [x] **Devcontainer trägt die Toolchain.** `.devcontainer/Dockerfile` installiert Go und `golangci-lint` in den Versionen von `repo.mk` mit SHA-256-Prüfung und setzt `BUILD_MODE=local` und `GOCACHE`; `make gates` im Devcontainer ist grün; die gemessenen Zeiten vor und nach der Änderung stehen in §7.
- [x] **Regeln und Doku nachgezogen.** [`AGENTS.md`](../../../../AGENTS.md) §3.1, [`docs/maintainer/quality.md`](../../../maintainer/quality.md) und [`harness/README.md`](../../../../harness/README.md) nennen den Nachweis „`make gates` im Modus der Umgebung“.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register ([`observations/`](../observations/)) fortgeschrieben, oder in §7 „keine Beobachtung angefallen“.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Schritt M0 (vor jedem Code): Messlauf.** Klärt, ob `go test` mit Coverage-Profil Ergebnisse cached, und misst `make gates` im Docker-Modus als Ausgangswert (Wanduhr, der Ablauf ist parallel). Das Ergebnis legt fest, ob Test und Coverage im lokalen Modus ein Lauf werden ([ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 7).

**Ergebnis M0** (gemessen am 2026-10-31, Linux `amd64`, ein Lauf je Zelle, [Messprotokoll](../../../messungen/slice-v1-devcontainer-bau-lokal/README.md)): `go test` cached auch mit `-coverpkg` und `-coverprofile` (24 von 24 Paketen aus dem Cache, 0,8 bis 1,2 s; vollständig kalt war nur der erste Lauf mit 11,9 s, die beiden Läufe mit Coverage-Profil, die im selben Cache-Volume nach diesem ersten Lauf folgten, hatten schon einen warmen Build-Cache: 5,5 s und 4,6 s); das Profil des gecachten Laufs hat dieselbe Größe und dieselbe Summe (92,1 %) wie das des ersten Coverage-Laufs, sein Hash ist verschieden. Nach der Änderung einer Testdatei lief 1 von 24 Paketen neu. Das Lint-Image hat `bash` und Go `1.27.1` (nicht `1.27.3`). `make gates` im Docker-Modus: 89,9 s und 89,3 s. **Folge:** Test und Coverage bleiben im lokalen Modus zwei Targets, beide aus dem Cache; [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleibt unverändert, die Folge-ADR des Architect-Reports entfällt; `scripts/bau.sh` darf für `lint` `bash` voraussetzen.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/README.md` | unverändert | Gelesen: Zeile 65 trägt die Indexzeile für [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (Commit `b371e262`); der Architect-Report F-6 berichtigt die frühere Planzeile „update“ |
| `spec/spezifikation.md` | update | [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Zeilen 368 bis 374) um den lokalen Modus erweitern; [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (Zeilen 1146 bis 1168): Spalte *Zweck* nennt `docker build --target …` (Docker-Modus) und den lokalen Befehl; Zeile 99 (CI-Runner, „Keine Host-Go-Toolchain“) bleibt, weil die CI im Docker-Modus läuft. [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) sagt zusätzlich: Der lokale Modus gilt für alle Make-Targets mit Docker-Pfad, nicht nur für die in [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) geführten (`test-docker`, `govulncheck` und `build-binaries` stehen dort nicht); Release-Artefakte (Binaries, Pakete, Images) entstehen in der Release-Pipeline im Docker-Modus (`publish.yml` ruft `make build-binaries` ohne `BUILD_MODE`); `packages` baut seine Binaries immer im Docker-Modus, ein lokales `build-binaries` ist Entwicklungs-Komfort; ein Teillauf (`PKG`, `RUN`) ist kein Nachweis. [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) (Zeilen 58 bis 71) und [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleiben unverändert, gelesen; werden Test und Coverage im lokalen Modus ein Lauf, schärft eine Folge-ADR [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) (Architect-Report F-2) |
| `scripts/bau.sh` | neu | Stage-Befehle `compile`, `lint`, `test`, `coverage`: eine Quelle für `Dockerfile` und lokalen Modus; setzt `CGO_ENABLED=0` bei `compile`, `test` und `coverage` und sonst nichts; `GOFLAGS` kommt aus der Umgebung (Docker: `ENV` der Stage `deps`, `Dockerfile` Zeile 51; lokal: `ENV` im `.devcontainer/Dockerfile`). `lint` setzt nichts: Die Stage baut auf `golangci/golangci-lint`, nicht auf `deps`, und kennt weder `GOFLAGS` noch `CGO_ENABLED`; so bleibt der Docker-Modus unverändert (Architect-Report F-3). Das Lint-Image hat `bash` und Go `1.27.1` (M0, [Messprotokoll](../../../messungen/slice-v1-devcontainer-bau-lokal/README.md), Abschnitt C) |
| `scripts/coverage-gate.sh` | unverändert | Kopf gelesen (Aufruf `coverage-gate.sh <coverage-func.txt> <threshold>`, Bootstrap-Modus über `COVERAGE_BOOTSTRAP`); der Aufruf aus `scripts/bau.sh` ist derselbe wie aus `Dockerfile` Zeile 146 |
| `Dockerfile` | update | Die `RUN`-Zeilen der Stages `compile` (Zeile 68), `lint` (Zeile 82), `test` (Zeile 88) und `coverage` (Zeile 136) rufen `scripts/bau.sh`; die Stage-Struktur nach [`SPEC-025`](../../../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) bleibt |
| `repo.mk` | update | `BUILD_MODE ?= docker`; lokale Rezepte für `deps`, `compile`, `lint`, `test`, `coverage-gate`, `test-docker`, `govulncheck`, `build-binaries`; `PKG` und `RUN` für `test` (nur von der Kommandozeile, über `BAU_PKG` und `BAU_RUN` in der Umgebung an das Skript; im Docker-Modus ein Fehler); `modus-pruefung`; `--no-cache-filter` (Zeilen 77 bis 79) bleibt für den Docker-Modus. `packages` baut seine Binaries immer im Docker-Modus (Rezept mit `$(MAKE) BUILD_MODE=docker build-binaries`), weil es sonst von `build-binaries` im Modus der Umgebung abhinge (Review F-1). Eine Ablehnung von `PKG` und `RUN` beim Einlesen über `MAKECMDGOALS` für `gates`, `record-gates`, `ci` und `fullbuild`. Kopplungskommentar bei `GO_VERSION` und `GOLANGCI_LINT_VERSION` (Zeilen 43 und 44): „Hebung betrifft drei Stellen: `Dockerfile`, diese Datei, `.devcontainer/Dockerfile`“ (Architect-Report F-4, F-5) |
| `.devcontainer/Dockerfile` | update | Hinter dem verwalteten Block: Go und `golangci-lint` als Release-Archive mit SHA-256 (je Architektur), `ENV BUILD_MODE=local`, `ENV GOCACHE` im Workspace-Volume; die Zeilen überstehen `u-boot generate devcontainer` (gemessen am 2026-10-31). Der verwaltete Block endet mit `USER vscode`; die Installation braucht davor `USER root` und danach wieder `USER vscode`. Der Kopplungshinweis „Hebung betrifft drei Stellen“ steht auch hier (Architect-Report F-4) |
| `.gitignore`, `.dockerignore` | update | Der Go-Cache liegt unter `.cache/` im Workspace (`GOCACHE=/workspaces/u-boot/.cache/go-build`); der Eintrag `.cache/` gehört in beide Dateien: nicht in Git, nicht in den Docker-Kontext (kein Eintrag dafür heute, Architect-Report F-7) |
| `.d-check.yml` | update | `scan.ignore` nimmt `.cache/**` auf (Go-Cache des Devcontainers im Workspace, fremde Markdown-Dateien der Abhängigkeiten; sonst prüft das Doku-Gate den Cache mit: gemessen 923 statt 904 Dateien). Kein Gate, keine Schwelle und kein Modul ändert sich ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| `tools/messung/devcontainer-bau-lokal/`, `docs/messungen/slice-v1-devcontainer-bau-lokal/`, `docs/messungen/README.md` | neu / update | Skripte und Protokoll von M0 und der Messung im Devcontainer-Image (Abschnitt F); Indexzeile |
| `spec/architecture.md` | update | Zeile 714 beschreibt `make test` als Docker-Stage und erhält den Zusatz „im Docker-Modus“; Zeile 716 (ein getrenntes Make-Target in der CI) gelesen, unverändert (Review F-6); Schichten und Importregeln bleiben unberührt, weil der Slice keinen Go-Quelltext ändert (Architect-Report F-7) |
| `scripts/build-packages.sh` | unverändert | Gelesen: Kopf und Aufruf (Zeilen 3 bis 4, 29): `nfpm` läuft per `docker run`; der Slice ändert das nicht (Architect-Report F-5) |
| `tools/mk/devcontainer.mk` | unverändert | Gelesen: die Targets `dev-up`, `dev-down`, `dev-shell` rufen nur `u-boot devcontainer`; sie kennen keinen Modus |
| `AGENTS.md` | update | §3.1: der Nachweis ist `make gates` im Modus der Umgebung; die Begründung (Reproduzierbarkeit, Supply-Chain) bleibt; Bezug auf [ADR-0037](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) als Plan-Anker nach §3.6 |
| `docs/maintainer/quality.md` | update | Zeilen 20 bis 36 (Lint im Docker-Modus und im lokalen Modus) und Zeilen 316 bis 318 (CI: „Keine Host-Toolchain“ gilt im Docker-Modus, die Workflows setzen `BUILD_MODE` nicht); Zeile 72 (`make gates`) gelesen, unverändert (Review F-6) |
| `harness/README.md` | update | Sensors-Tabelle: Zeile für `make modus-pruefung` (kein Gate); Satz zum Modus unter *Sensors* |
| `harness/verification.md` | unverändert | Gelesen (Zeilen 75 bis 87): die Sensor-Tabelle nennt Make-Targets ohne Modus |
| `.github/workflows/ci.yml`, `.github/workflows/integration.yml`, `.github/workflows/publish.yml` | unverändert | Gelesen: sie rufen `make gates`, `make govulncheck`, `make image-scan` und `make test-docker` (ci.yml und integration.yml) und `make build-binaries VERSION="$VERSION"` (publish.yml, Zeile 177), alle ohne `BUILD_MODE`, also im Docker-Modus |
| `docs/plan/planning/open/slice-v1-devcontainer-instanz-label.md`, `docs/plan/planning/done/slice-v1-konfigdatei-name-als-konstante.md`, `docs/plan/planning/open/slice-v1-devcontainer-claude-oauth-token.md`, `docs/plan/planning/open/README.md` | neu / update (fremde Vorgänge) | Entwürfe von Folgevorgängen aus dem Gespräch zu diesem Slice und die Indexzeilen dazu; sie liegen in der Range dieses Slice, sind aber **nicht** Teil seiner DoD |
| `CHANGELOG.md` | unverändert | Kein öffentlicher Vertrag berührt: Produktverhalten und Release-Artefakte ändern sich nicht (Kopf gelesen, Abschnitt `Unreleased` leer) |
| `docs/reviews/2026-10-31-slice-v1-devcontainer-bau-lokal*.md` | neu | Die vier Reports der Kette: Architect, Review, Nachlauf-Review, Verification |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/README.md` | update | Closure: Ruhe-Marker „Nichts in Arbeit.“ statt der Arbeitszeile (d-check `planning`), Datum „Letzte Änderung“; Zeile des Slice in `open/README.md` auf „geschlossen“, Verweis-Nachzug des Moves durch `make slice-mv` |
| `docs/plan/planning/observations/` | neu / update | Closure: drei neue Beobachtungen (je ein Verzeichnis mit `observation.md`, `state.md` und einer Evidence-Datei) und sechs weitere Evidence-Dateien in bestehenden Beobachtungen, mit Zeilen in deren `state.md` (Aufzählung in §7) |

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

- **Die Modi laufen auseinander** (ein Lauf lokal grün, in der CI rot), weil `scripts/bau.sh` nur Teile der Stages teilt — **Ausgang:** entfallen: beide Modi führen die Stage-Befehle aus `scripts/bau.sh` aus, `make gates` ist in beiden grün (§7); der CI-Lauf im Docker-Modus steht noch aus.
- **Die Pins stehen an drei Stellen** (`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`) ohne Abgleich, bis d-check das Modul liefert — **Ausgang:** weiter offen: → [`BEO-WZG/pin-an-drei-stellen-ohne-abgleich`](../observations/BEO-WZG/pin-an-drei-stellen-ohne-abgleich/observation.md) im Register, mit Verweis auf Issue 14 in `pt9912/d-check`.
- **Ein Teillauf zählt als Nachweis** (`make test PKG=… RUN=…`): `repo.mk` lehnt `PKG` und `RUN` beim Einlesen ab, wenn `gates`, `record-gates`, `ci` oder `fullbuild` unter den Zielen stehen; mit anderen Zielen (etwa `make lint test PKG=…`) bleibt ein Teillauf erlaubt, ohne Nachweis zu sein; `make modus-pruefung` prüft die Ablehnung — **Ausgang:** entfallen: die Ablehnung ist gepinnt (Gegenproben, sechs Mutationen im Nachlauf-Review rot); ein Teillauf mit anderen Zielen bleibt erlaubt, ohne Nachweis zu sein.
- **`govulncheck` im lokalen Modus lädt sein Werkzeug über das Netz** (`go install …@$(GOVULNCHECK_VERSION)`, Modul-Proxy statt SHA-256-Archiv); im Docker-Modus ist die Lage gleich (Version gepinnt, kein Hash), also kein neuer Befund und ein akzeptiertes Negativ (Architect-Report F-4) — **Ausgang:** entfallen: im Docker-Modus ist die Lage gleich, kein neues Risiko dieses Slice; lokal Exit 0 (Messprotokoll G).
- **`test-docker` im verschachtelten Daemon ist nicht gleichwertig zum Docker-Modus am Host:** 9 Tests schlagen fehl (vier davon, weil das Startscript bei eingeschränkten User-Namespaces nicht mit Exit 11 endet; bei fünf ist die Ursache nicht geklärt), im Docker-Modus am selben Host 0 Fehlschläge und 4 `SKIP` ([Messprotokoll](../../../messungen/slice-v1-devcontainer-bau-lokal/README.md), Abschnitt G). [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) sagt deshalb, dass `test-docker` im Devcontainer nicht als Nachweis zugesagt ist. — **Ausgang:** weiter offen: → [`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`](../observations/BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig/observation.md) im Register; die fünf ungeklärten Ursachen sind ihr Gegenstand.
- **Test-Cache verdeckt einen umgebungsabhängigen Test** im lokalen Modus — **Ausgang:** entfallen: der Docker-Modus (Host und CI) läuft weiter ohne Cache (`--no-cache-filter` unverändert) und fängt es.
- **`node:2.1.1` mit `claude-code:1.0.5`** ist nicht gemessen (das Handbuch hat `claude-code` mit `node:1` gemessen); der Devcontainer-Build kann scheitern — **Ausgang:** entfallen: `make dev-up` baut mit allen Features, Exit 0, 334 s (Messprotokoll H).
- **Der Hook-Guard blockt `go` in der Befehlsposition**, auch im Devcontainer; ein Agent, der `go test` direkt ruft, wird gestoppt — **Ausgang:** weiter offen: → [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md); der Guard ist unverändert, ein Lauf eines Agenten im Devcontainer ist nicht gemessen.
- **`GOCACHE` im Workspace-Volume** wächst ohne Grenze — **Ausgang:** entfallen: `.cache/` ist aus Git, Docker-Kontext und Doku-Scan ausgenommen und ohne Folge löschbar; das Wachstum selbst ist nicht gemessen.

## 7. Closure-Notiz

Alle Zeiten: Linux `amd64`, lokale Engine, **ein Lauf je Zelle**, keine Streuung bestimmt ([Messprotokoll](../../../messungen/slice-v1-devcontainer-bau-lokal/README.md)).

- **Was hat funktioniert:**
  - **Docker-Modus unverändert gleichwertig.** Review (Frage 1) und Verifier (Abgleich mit [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 4) lesen die Stages `compile`, `lint`, `test`, `coverage` mit `scripts/bau.sh` gleichwertig zu vorher; `make gates` im Docker-Modus: 89,9 s und 89,3 s vorher, 91,1 s nach dem Umbau des `Dockerfile`, Coverage 92,10 % bei Schwelle 90 %.
  - **Lokaler Modus, Zeiten vor und nach** (`make gates`): Docker-Modus 89,9 s / 89,3 s / 91,1 s; lokal im Image des Devcontainers 42,7 s kalt und 26,4 s warm (Stand vor der Fix-Runde), 43,4 s und 26,0 s auf dem Endstand; im gestarteten Devcontainer 26,5 s und 27,7 s mit warmem Cache (der erste Lauf mit leerem Cache dauerte 48,6 s, seine Rohdaten sind ersetzt). Einzeln warm: `lint` 0,6 bis 0,7 s, `test` 0,4 s, `coverage-gate` 1,2 bis 1,3 s gegen 20 bis 31 s, 18 bis 20 s und 19 bis 20 s im Docker-Modus. Die 26 s warm bestehen zum größten Teil aus den Gates ohne Go (`ablage-pruefung`, `docs-check`, `baseline-verify`), nicht aus `lint`, `test` und `coverage-gate` (zusammen etwa 2 s).
  - **M0 entschied vor dem Code:** `go test` cached auch mit `-coverpkg` und `-coverprofile` (24 von 24 Paketen aus dem Cache), nach der Änderung einer Testdatei läuft 1 von 24 Paketen neu. Test und Coverage bleiben zwei Targets, [`SPEC-032`](../../../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) blieb unverändert, die vom Architect erwogene Folge-ADR entfiel.
  - **Der gestartete Devcontainer trägt es** (Messung H, Verifier F-1): `make dev-up` mit allen Features (`node:2.1.1`, `claude-code:1.0.5`, `homebrew:2.0.5`, `ollama`, `gh`) und dem neuen `.devcontainer/Dockerfile` baut in 334 s, `go` 1.27.3 und `golangci-lint` 2.14.0 laufen, `GOCACHE` und `GOMODCACHE` im Workspace-Volume sind für `vscode` schreibbar.
  - **Der Wächter `make modus-pruefung` ist scharf:** sechs Mutationen an `repo.mk` färben ihn im Nachlauf-Review rot, der Verifier bestätigte die Bedingungen von Hand mit `make -n`; die vier SHA-256-Werte stimmen mit go.dev und dem Release von `golangci-lint` überein (Verifier).
- **Was ging anders als geplant:**
  - **Eine ADR-Annahme hielt der Messung nicht stand:** [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 3 führt `test-docker` als im lokalen Modus laufend; Messung G zeigt im verschachtelten Daemon 9 Fehlschläge gegen 0 im Docker-Modus am selben Host (Verifier F-2). Die Entscheidung ist deshalb nicht verletzt (der Nachweis ist `make gates`, und `test-docker` gehört nicht dazu); [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) führt `test-docker` im Devcontainer als nicht zugesagt. Eine Folge-ADR war nach dem Verifier nicht nötig; würde `test-docker` dort Nachweis, wäre sie es.
  - **Der erste Review war merge-blockierend** (2 MEDIUM: `packages` baute im lokalen Modus seine Binaries lokal; Bedingungen in [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) ohne pinnenden Test). Behoben in der Fix-Runde (`packages` ruft `$(MAKE) BUILD_MODE=docker build-binaries`, `modus-pruefung` erweitert); der Nachlauf-Review fand 0 HIGH und 0 MEDIUM, 6 LOW, 5 INFO, der Verifier 0 HIGH, 1 MEDIUM, 4 LOW, 2 INFO.
  - **Der Verifier schränkte ein:** „`make gates` im Devcontainer“ war am Image ohne die Features gemessen. Die Messung H im echten Devcontainer behob das und fand den Fehler im Wächter, der `BUILD_MODE=local` aus der Umgebung nicht neutralisierte (behoben mit `env -u BUILD_MODE`).
  - **Der Entwurf von ADR und Plan entstand im Implementer-Kontext** ohne Vorlage und ohne Architect; der Architect lief nach (Beobachtung unten).
  - **Noch nicht erfüllt ist der Teil des Closure-Triggers (§5), der den Pull Request und die CI betrifft:** Ein Pull Request ist nicht gemergt; die CI (Docker-Modus) lief für diesen Stand nicht, der Verifier fand für `98699f3e` keinen Lauf. Der Beleg für den Docker-Modus ist der lokale Lauf von `make gates`.
- **Steering-Loop-Eintrag:** Sensor `make modus-pruefung` neu (Fitness Function von [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md)): pinnt die Bedingungen der zwei Modi aus [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) mit `make -n`, mit Gegenproben und Mutationsbelegen. Lehre: Ein Wächter über eine Umgebung wird in der ausgelieferten Umgebung gefahren, nicht nur an einer Variante von ihr; erst der gestartete Devcontainer zeigte, dass der Wächter `BUILD_MODE` aus der Umgebung voraussetzte. Gezählt, nicht verkörpert: Der Sensor entstand aus der ADR, nicht aus einem Auslöser von drei; die Lehre ist `BEO-DOK/messvariante-belegt-nicht-die-auslieferung` (9×, Ausgang beim Projektinhaber). Benannte Spec-Lücke ist keine entstanden: [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) trägt die Nicht-Zusage zu `test-docker`.
- **Beobachtungs-Register (`../observations/`):**
  - Neu angelegt, je mit Beleg `evidence/slice-v1-devcontainer-bau-lokal.md`, Zähler 1×: [`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`](../observations/BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig/observation.md) (Kandidat a, Risiko in §6); [`BEO-WZG/pin-an-drei-stellen-ohne-abgleich`](../observations/BEO-WZG/pin-an-drei-stellen-ohne-abgleich/observation.md) (Kandidat c, Risiko in §6; verwandt mit `BEO-WZG/kopplung-nicht-nachgezogen`, die eingetretene Abweichung ist dort der Gegenstand, hier nur der fehlende Sensor); [`BEO-DOK/implementer-schreibt-adr-und-plan-ohne-vorlage-und-architect`](../observations/BEO-DOK/implementer-schreibt-adr-und-plan-ohne-vorlage-und-architect/observation.md) (Kandidat f; weder `dod-wortlaut-vom-implementer-geaendert` noch `implementer-entscheidet-adr-trigger` passt wörtlich, beide sind in der Observation als verwandt benannt und zählen getrennt).
  - Ergänzt: `evidence/slice-v1-devcontainer-bau-lokal.md` in `BEO-DOK/messvariante-belegt-nicht-die-auslieferung` (Kandidat b: Zähler steht bei 9×, über der Schwelle, Ausgang beim Projektinhaber, unverändert; der Fall ist der Trigger-Fall „Variante gemessen, Auslieferung gelesen“); in `BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle` (Kandidat d: 2×, die Passung ist in der Evidence als eingeschränkt benannt, weil der Fremdzustand dort ein erzeugter Eintrag war, hier die Gleichwertigkeit einer fremden Umgebung; unter der Schwelle); in `BEO-WZG/host-werkzeug-ausserhalb-docker` (Kandidat e, Selbstauskunft des Implementers, an keiner Datei nachprüfbar: Zähler steht bei 17×, über der Schwelle, Ausgang beim Projektinhaber, unverändert; die Stand-Zeile nennt noch elf Belege, der Zähler folgt aus den Dateien).
  - Rückfälle nach der Verkörperung aus den Finding-Klassen der Reports, je eine Evidence-Datei: `BEO-DOK/plan-aussage-ueber-code` (61×), `BEO-DOK/spec-bedingung-ohne-code-entsprechung` (26×), `BEO-CLI/kommentar-beschreibt-abwesendes` (10×); jeweils keine weitere Schärfung entschieden.
  - Über der Schwelle von drei stehen damit fünf der neun berührten Beobachtungen (`messvariante-belegt-nicht-die-auslieferung`, `host-werkzeug-ausserhalb-docker`, `plan-aussage-ueber-code`, `spec-bedingung-ohne-code-entsprechung`, `kommentar-beschreibt-abwesendes`). Für keine wird in dieser Closure ein Ausgang zugewiesen: Bei den ersten beiden steht die Entscheidung des Projektinhabers aus; die Regeln zu den anderen drei bestehen in `AGENTS.md` §3.10 und §3.9 und wurden vor diesem Slice entschieden.
- **Folge-Slices:** [`slice-v1-devcontainer-claude-oauth-token`](../open/slice-v1-devcontainer-claude-oauth-token.md), [`slice-v1-devcontainer-instanz-label`](../open/slice-v1-devcontainer-instanz-label.md) und [`slice-v1-konfigdatei-name-als-konstante`](../done/slice-v1-konfigdatei-name-als-konstante.md) sind Dateien in `open/`; sie entstanden im Gespräch zu diesem Slice und sind nicht Teil seiner DoD. Verwandt, nicht übernommen: [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md) (Feature-Optionen in `u-boot.yaml`, §1).
- **Risiken aus §6:** je Risiko genau ein Ausgang, in der Reihenfolge von §6 (Ausgänge auch dort):
  - Modi laufen auseinander: **entfallen**, soweit gemessen: beide Modi führen die Stage-Befehle aus `scripts/bau.sh`, `make gates` ist im Docker-Modus (91,1 s, 92,10 %) und im lokalen Modus grün; die CI im Docker-Modus steht noch aus (siehe oben).
  - Pins an drei Stellen: **weiter offen** → `BEO-WZG/pin-an-drei-stellen-ohne-abgleich`.
  - Teillauf zählt als Nachweis: **entfallen**: `repo.mk` lehnt `PKG` und `RUN` beim Einlesen ab, wenn `gates`, `record-gates`, `ci` oder `fullbuild` unter den Zielen stehen, `modus-pruefung` pinnt es mit Gegenproben (sechs Mutationen rot); mit anderen Zielen bleibt ein Teillauf erlaubt, ohne Nachweis zu sein.
  - `govulncheck` lokal lädt sein Werkzeug über das Netz: **entfallen** als Risiko dieses Slice: im Docker-Modus ist die Lage gleich (Version gepinnt, kein Hash), also kein neuer Befund (Architect F-4); lokal Exit 0, „No vulnerabilities found.“ (Messung G).
  - `test-docker` im verschachtelten Daemon: **weiter offen** → `BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig` (fünf ungeklärte Ursachen).
  - Test-Cache verdeckt einen umgebungsabhängigen Test: **entfallen**: Im Docker-Modus (Host und CI) läuft die Suite weiter ohne Cache (`--no-cache-filter` unverändert, Review Frage 4); ein lokal verdeckter Fehler fällt dort auf.
  - `node:2.1.1` mit `claude-code:1.0.5`: **entfallen**: `make dev-up` baut mit allen Features, Exit 0, 334 s (Messung H).
  - Hook-Guard blockt `go` in der Befehlsposition: **weiter offen** → `BEO-WZG/host-werkzeug-ausserhalb-docker` (der Guard ist unverändert; ein Lauf eines Agenten im Devcontainer ist nicht gemessen, die Entscheidung zum Guard steht dort aus).
  - `GOCACHE` im Workspace-Volume wächst ohne Grenze: **entfallen** als Risiko des Slice, mit Einschränkung: `.cache/` ist aus Git, Docker-Kontext und Doku-Scan ausgenommen und ohne Folge löschbar (ein Lauf mit leerem Cache kostete 48,6 s); das Wachstum selbst ist nicht gemessen, eine Größe liegt nicht vor.
- **Drei Paarungen:** Anker: kein Eintrag mit dem Feld für einen Herkunfts-Anker in diesem Slice, nichts zu paaren. Folge-Slice: die drei Folge-Slices sind Dateien in `open/`; jedes Risiko mit Ausgang „weiter offen“ hat sein Verzeichnis im Register, keines mit Ausgang „eingetreten“. Register: `make docs-check` nach dem Move ohne Befund (924 Dateien, 0 Befunde); jeder Verweis des Slice auf eine Beobachtung löst auf.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `tools/`, `scripts/` (WZG), `spec/`, `harness/`, `docs/plan/` (DOK). Beide sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md) als Sub-Areas geführt; eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist durchgegangen. Treffer: [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (Schwelle erreicht, kein Ausgang zugewiesen, Entscheidung des Projektinhabers zum Stolperdraht steht aus). Der Slice berührt den Gegenstand: Im Devcontainer ist `go` über `make` erlaubt, der Guard bleibt unverändert, und ein direkter Aufruf bleibt blockiert (§6).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
