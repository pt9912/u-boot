# Verification-Report: [`slice-v1-devcontainer-bau-lokal`](../plan/planning/done/slice-v1-devcontainer-bau-lokal.md) — 2026-10-31

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) und den ausgeführten Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Außer dieser Datei ist nichts geändert, nichts committet, nichts gepusht; `git status` war am Ende leer.

**Gegenstand:** `git diff 9b01cb44..98699f3e` (9 Commits, 27 Dateien, `HEAD` = `98699f3e`, Arbeitsbaum zu Beginn sauber). Eingang (gelesen, nicht als Beweis genommen): Architect-Report, Review-Report und Nachlauf-Report unter `docs/reviews/` mit dem Präfix des Slice-Datums, Messprotokoll [`README.md`](../messungen/slice-v1-devcontainer-bau-lokal/README.md) samt Rohdaten.

**Verdikt:** DoD (Liefer-Punkte 1 bis 3 und `make gates`): **erfüllt**, mit der Einschränkung F-1 (der echte Devcontainer ist nicht gestartet). Die Punkte Review, Closure-Notiz, Register, Risiko-Ausgänge und Paarungen sind Closure-Aufgaben des Planners und stehen offen; nichts im Stand verhindert sie. **Bereit für die Closure: ja**, unter der Bedingung, dass F-1 in der Closure-Notiz entschieden wird (entweder ein Lauf im echten Devcontainer oder ein ausdrücklich so benannter Rest). Keine HIGH-Befunde, ein MEDIUM, vier LOW, zwei INFO.

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-bau-lokal`](../plan/planning/done/slice-v1-devcontainer-bau-lokal.md)
- IDs: [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md), [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) und [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) (beide unverändert), [`LH-NFA-PORT-002`](../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten)
- Artefakte: `repo.mk`, `Dockerfile`, `scripts/bau.sh`, `.devcontainer/Dockerfile`, `spec/spezifikation.md`, `spec/architecture.md`, `AGENTS.md`, `docs/maintainer/quality.md`, `harness/README.md`, `.d-check.yml`, `.gitignore`, `.dockerignore`, `tools/messung/devcontainer-bau-lokal/`, [Messprotokoll-Verzeichnis](../messungen/slice-v1-devcontainer-bau-lokal/)

### 1. DoD-Abgleich

- [x] **Modus verkörpert** — Evidence: [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) hat im Diff (`spec/spezifikation.md`) den Docker-Modus-Punkt, den Punkt „Lokaler Modus“, den Release-Punkt und den Nachweis-Punkt; [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) nennt `BUILD_MODE` in den Variablen und in der Spalte *Zweck* je Target den lokalen Befehl. `repo.mk` Zeilen 77 bis 83: `BUILD_MODE ?= docker`, `$(error …)` für Werte außer `docker` und `local`; Rezepte unter `ifeq ($(BUILD_MODE),local)` (Zeile 130 ff.). Selbst gefahren: `make modus-pruefung` Exit 0, „OK (8 targets in both modes; build, run, image-scan, packages, docs-check in docker; BUILD_MODE, PKG and RUN rules)“. Stichproben mit `make -s -n`: `BUILD_MODE=local test-docker coverage-gate deps lint` liefert `go test -v -tags docker …`, `OUT=coverage … bash scripts/bau.sh coverage`, `go mod download`, `bash scripts/bau.sh lint` ohne `docker build`; `BUILD_MODE=docker test RUN=X` bricht mit „PKG and RUN exist in the local mode only“ ab; `BUILD_MODE=local packages` ruft `make BUILD_MODE=docker build-binaries` und zeigt `docker build … binary-archive`.
- [x] **Devcontainer trägt die Toolchain** — Evidence: `.devcontainer/Dockerfile` Zeilen 29 bis 71 (Diff): `USER root`, Go und `golangci-lint` als Release-Archive je Architektur, `sha256sum -c` vor dem Entpacken, `ENV BUILD_MODE=local`, `GOFLAGS`, `GOTOOLCHAIN=local`, `GOCACHE` und `GOMODCACHE` unter `/workspaces/u-boot/.cache`, `USER vscode`. **Prüfsummen selbst gegen die Quellen geprüft** (curl am 2026-10-31): `go.dev/dl/?mode=json&include=all` für `go1.27.3` liefert `linux-amd64` `65efbc80…379b40` und `linux-arm64` `a06e9930…165f49`; `golangci-lint-2.14.0-checksums.txt` des Releases `v2.14.0` liefert `linux-amd64` `ab90aeb7…1032ab` und `linux-arm64` `ee7ec5f3…299bae`. Alle vier stimmen mit den `ARG`-Zeilen 40, 41, 43, 44 überein. Versionen gleich an drei Stellen: `repo.mk` (`GO_VERSION ?= 1.27.3`, `GOLANGCI_LINT_VERSION ?= v2.14.0`), `Dockerfile` Zeilen 35 und 36, `.devcontainer/Dockerfile` Zeilen 39 und 42. Zeiten vor und nach (Messprotokoll, Abschnitte D und F): Docker-Modus vor der Änderung 89,9 s und 89,3 s; lokal im Image 43,4 s (erster Lauf, d-check-Image wird geladen) und 26,0 s (warm); Docker-Modus nach dem Umbau des `Dockerfile` 91,1 s; einzeln lokal warm `lint` 0,7 s, `test` 0,4 s, `coverage-gate` 1,2 s gegen 20 bis 31 s je Gate im Docker-Modus (ADR-Kontext). „`make gates` im Devcontainer grün“ gilt für das Image aus diesem Dockerfile, nicht für den gestarteten Devcontainer: siehe F-1. Die Zeiten stehen im Messprotokoll; §7 des Slice ist noch leer und wird in der Closure gefüllt.
- [x] **Regeln und Doku nachgezogen** — Evidence: `AGENTS.md` §3.1 (Docker-Modus am Host und in der CI, lokaler Modus im Devcontainer, „Nachweis … ist `make gates` im Modus der Umgebung; ein Teillauf zählt nicht“, Verweis auf [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) als Plan-Anker nach §3.6, Begründung unverändert); `docs/maintainer/quality.md` (Lint in beiden Modi, `test-docker`-Absatz mit Messung, CI „Keine Host-Toolchain“ im Docker-Modus); `harness/README.md` (Sensors-Zeile `make modus-pruefung`, Absatz „Modus“); `spec/architecture.md` Zeile 714 („im Docker-Modus … im lokalen Modus derselbe Befehl aus `scripts/bau.sh`“), Zeile 716 unverändert gelesen. Die Architektur-Sicht trägt weiterhin keine ADR- oder Slice-Kennung.
- [x] `make gates` grün — selbst gefahren (Docker-Modus am Host): Exit 0, `coverage-gate: OK — coverage 92.10% meets threshold 90%`, `ablage-pruefung: OK`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`, `d-check: 907 Datei(en) geprüft, 0 Befund(e)`.
- [ ] Review durchgeführt — Report (und Nachlauf) liegen vor; Verdikt der Reviews: 0 HIGH, 0 MEDIUM im Nachlauf. Das Häkchen setzt der Planner.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (Planner); §7 ist Platzhalter.
- [ ] Beobachtungs-Register fortgeschrieben — offen (Planner). Der Plan kündigt zwei Beobachtungen an (Pin-Abgleich, `test-docker` im verschachtelten Daemon); beide existieren heute nicht (Nachlauf N-6). Der Stand verhindert das nicht.
- [ ] Risiken aus §6 mit Ausgang — offen (Planner); neun Risiken, alle „bei Closure“. Belege für die Ausgänge: siehe F-1 (Risiko `node:2.1.1` mit `claude-code:1.0.5`) und Abschnitt 2.
- [ ] Drei Paarungen — offen (Planner).

### 2. Plan gegen Code

`git diff --stat 9b01cb44..98699f3e` gegen die Plan-Tabelle §3:

| Diff-Datei | Plan-Zeile | Ergebnis |
|---|---|---|
| `spec/spezifikation.md`, `scripts/bau.sh`, `Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`, `.gitignore`, `.dockerignore`, `.d-check.yml`, `spec/architecture.md`, `AGENTS.md`, `docs/maintainer/quality.md`, `harness/README.md` | je eine Zeile | stimmt |
| `tools/messung/devcontainer-bau-lokal/*` (4), das Messprotokoll-Verzeichnis (4), `docs/messungen/README.md` | Sammelzeile | stimmt |
| Slice selbst, `docs/reviews/…-bau-lokal.md`, `…-nachlauf.md` | — (Slice und Review-Artefakte) | erwartet |
| `docs/plan/planning/open/README.md`, der Entwurf von [`slice-v1-devcontainer-instanz-label`](../plan/planning/open/slice-v1-devcontainer-instanz-label.md) | **keine** | F-4 |
| `docs/plan/adr/README.md`, `scripts/coverage-gate.sh`, `scripts/build-packages.sh`, `tools/mk/devcontainer.mk`, `harness/verification.md`, `.github/workflows/ci.yml`, `.github/workflows/integration.yml`, `CHANGELOG.md` | „unverändert“ | im Diff nicht enthalten; stimmt |

Gelesene Stellen zu Planaussagen: `Dockerfile` Zeilen 68, 82, 88, 136 rufen `scripts/bau.sh` (wie im Plan); `Dockerfile` Zeile 50 bis 52 setzt `GOFLAGS` in der Stage `deps` (Plan: „Zeile 51“, stimmt); `docs/plan/adr/README.md` Zeile 65 trägt die Indexzeile; `spec/architecture.md` Zeile 714; `scripts/bau.sh` setzt `CGO_ENABLED=0` an genau `compile`, `test`, `coverage`; `.gitignore` hat `/.cache/` (Zeile 53); die Messzahlen in Plan und Protokoll stimmen mit `roh.txt` und `roh-lokal-rest.txt` überein (9 Fehlschläge, 1627 gegen 1621 bestanden, 4 `SKIP`; 1621 + 9 + 1 = 1627 + 4).

Aussagen des Plans, die nicht zum Code passen (alle ohne Auswirkung auf das Verhalten):

- Der Plan sagt zu [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), „Release-Artefakte (`build-binaries`, `packages`) entstehen nur im Docker-Modus“. Der Code hat eine lokale Variante von `build-binaries` (`repo.mk`, `ifeq ($(BUILD_MODE),local)` vor `build-binaries`); die Spezifikation sagt zutreffend „entstehen in der Release-Pipeline im Docker-Modus“. Die Planzeile ist veraltet (F-3).
- Risiko 3 sagt, die Rezepte der `gates`-Targets lehnten `PKG`/`RUN` ab oder setzten sie leer. Der Code lehnt beim Einlesen über `MAKECMDGOALS` ab (`repo.mk` Zeilen 273 bis 275) und nur, wenn eines der vier Ziele auf der Kommandozeile steht (INFO, F-6).

**Bedingungen in [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) und ihr Pin** (AGENTS.md §3.10, letzter Punkt):

| Bedingung | Pin |
|---|---|
| Docker-Modus: Targets rufen `docker build`/`docker run` | `modus-pruefung`, Schleife über 8 Targets (nur Vorkommen, nicht „ausschließlich“; das war schon vorher so) |
| Lokal: acht Targets ohne `docker build` | `modus-pruefung`, `grep -q 'docker build'` |
| `build`, `run`, `image-scan`, `packages`, `docs-check` rufen Docker in beiden Modi | `modus-pruefung`; `packages` zusätzlich über die Rekursion `$(MAKE) BUILD_MODE=docker` (von mir geprüft: `make -n BUILD_MODE=local packages` zeigt `docker build … binary-archive`) |
| `BUILD_MODE` ausdrücklich, andere Werte Fehler | `modus-pruefung` (`BUILD_MODE=foo`) |
| `PKG`/`RUN` nur lokal, nur Kommandozeile, nicht mit den vier Nachweis-Zielen, die vier ohne sie angenommen | `modus-pruefung` (PKG im Docker-Modus, Umgebung, 4 Ziele × 2 Variablen, 4 Ziele ohne) |
| `test-docker` im verschachtelten Daemon nicht als Nachweis zugesagt | Nicht-Zusage, belegt durch Messprotokoll G |
| Befehle von `compile`, `lint`, `test`, `coverage` stehen in `scripts/bau.sh`, das auch die Dockerfile-Stages rufen | **kein direkter Pin**; nur mittelbar: CI und `make gates` im Docker-Modus führen die Stages aus (F-5) |
| „eine Erkennung von `go` im `PATH` gibt es nicht“ | **kein Pin** (negative Aussage; im Diff von `repo.mk` keine Erkennung gelesen) (F-5) |
| [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) Spalte *Zweck*, lokale Befehle (`go mod download`, `scripts/bau.sh …`) | **kein Pin** des genannten Befehls; `modus-pruefung` prüft nur das Fehlen von `docker build` (F-5) |
| Release-Pipeline im Docker-Modus | **kein Pin**; wahr, weil `.github/workflows/publish.yml` (ruft `make build-binaries`, Zeile 177) `BUILD_MODE` nicht setzt (`grep BUILD_MODE .github` ohne Treffer); `publish.yml` steht nicht in der „unverändert“-Zeile des Plans (F-3) |
| `RUN` im Docker-Modus ein Fehler | gleicher Codepfad wie `PKG`, aber nur `PKG` wird von `modus-pruefung` geprobt (F-5) |

### 3. ADR-0037 gegen den Code

| Entscheidung | Befund |
|---|---|
| 1 Zwei Modi, Docker Standard | umgesetzt (`repo.mk`, `BUILD_MODE ?= docker`) |
| 2 Schalter ausdrücklich, `ENV` im handgeschriebenen Teil, keine Erkennung | umgesetzt (`.devcontainer/Dockerfile` hinter `END U-BOOT MANAGED BLOCK`); „Zeilen überstehen `u-boot generate devcontainer`“ ist am 2026-10-31 gemessen (Plan), von mir nicht wiederholt |
| 3 Umfang: acht Targets lokal, `build`, `run`, `image-scan`, `packages` Docker | umgesetzt. **Abweichung A:** `test-docker` läuft lokal, ist aber laut [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) im verschachtelten Daemon nicht als Nachweis zugesagt (9 Fehlschläge, Messprotokoll G, Rohdaten `roh-lokal-rest.txt`). Die ADR nennt `test-docker` lokal, sagt aber nicht, dass es dort Nachweis ist; Entscheidung 6 macht `make gates` zum Nachweis, und `test-docker` gehört nicht zu `gates`. Daher kein Widerspruch zur Entscheidung, aber eine ungemessene Annahme der ADR (der verschachtelte Daemon trägt die Tests), die das Messprotokoll widerlegt. **Keine Folge-ADR nötig**; die Spezifikation trägt die Nicht-Zusage. Die Re-Evaluierungs-Trigger der ADR decken den Fall nicht ausdrücklich ab (F-2). **Abweichung B:** `packages` ist in der ADR „im Docker-Modus“; umgesetzt ist, dass sein `build-binaries` immer im Docker-Modus läuft (`$(MAKE) BUILD_MODE=docker build-binaries`), `nfpm` läuft wie vorher per `docker run`. Das ist die Entscheidung, nicht ihre Abweichung. |
| 4 Eine Quelle für die Stage-Befehle | umgesetzt (`scripts/bau.sh`; `Dockerfile` Zeilen 68, 82, 88, 136). Docker-Modus-Ergebnis unverändert: `make gates` grün, 92,10 % wie vorher (Summe 92,1 % in M0), 91,1 s gegen 89,9 und 89,3 s. Das Coverage-Verhalten (Bootstrap-Zweig, `pipefail`) steht gleichlautend in `bau.sh`; die Stage `coverage` setzt `SHELL` nicht mehr, `bau.sh` setzt `set -eo pipefail`. [`SPEC-025`](../../spec/spezifikation.md#spec-025--stages-des-multi-stage-dockerfiles) und [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) im Diff nicht berührt |
| 5 Pins im Devcontainer mit SHA-256 | umgesetzt, Prüfsummen gegen Quellen bestätigt (oben). Der Dockerfile-Kommentar sagt ehrlich, dass sie nicht gegen ein von Anfang an schädliches Release schützen |
| 6 Nachweis `make gates` im Modus der Umgebung | umgesetzt in `AGENTS.md` §3.1, `harness/README.md`, [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows); Teillauf-Sperre gepinnt |
| 7 Test-Cache; erst messen | gemessen (M0, Abschnitt A, B, E); Folge: zwei Targets, beide aus dem Cache; [`SPEC-032`](../../spec/spezifikation.md#spec-032--bootstrap-verhalten-der-coverage-stage) bleibt, Folge-ADR entfällt. Konsistent mit der ADR („legt danach fest“). Der Cache-Treffer des Coverage-Profils ist ein Lauf je Zelle |
| 8 `GOCACHE` im Workspace-Volume | umgesetzt (`GOCACHE=/workspaces/u-boot/.cache/go-build`, zusätzlich `GOMODCACHE` unter `.cache`; die ADR nennt nur `GOCACHE`: harmlose Erweiterung, in `.gitignore`, `.dockerignore`, `.d-check.yml` mitgezogen, keine Gate-Lockerung nach AGENTS.md §3.6, da kein Gate, keine Schwelle und kein Modul sich ändert). Der Pfad nennt den Workspace-Namen `u-boot` fest |
| Fitness Function `make modus-pruefung` | umgesetzt und deutlich schärfer als die ADR (Gegenproben) |

Ergebnis: keine Abweichung verlangt eine Folge-ADR. Abweichung A ist als Nicht-Zusage in der Spezifikation geführt und im Plan als Risiko; sie braucht in der Closure die angekündigte Beobachtung. Würde `test-docker` im Devcontainer als Nachweis gelten sollen, wäre eine Folge-ADR nötig.

### 4. Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | Exit 0; Coverage 92,10 %; `docs-check` 0 Befunde |
| `make modus-pruefung` | pass | Exit 0, siehe Abschnitt 1 |
| `make docs-check` | pass | 907 Dateien, 0 Befunde (Lauf nach Anlegen dieses Reports, siehe unten) |
| `make doc-trace` | pass | „84 Anforderung(en), 0 Waise(n)“ |
| `make doc-commits RANGE="9b01cb44..98699f3e"` | pass | 9 Commits, 0 Befunde |
| `make doc-immutable RANGE="9b01cb44..98699f3e"` | pass | 9 Commits, 0 Befunde (kein Accepted-ADR verändert; [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) liegt vor der Range) |
| Prüfsummen gegen go.dev/dl und GitHub-Release | pass | vier von vier gleich |
| `gh run list --branch main` | gelesen | letzter Lauf auf `main` ist `b371e262` (vor dem Code des Slice): `ci` und `integration-docker` grün; für den Branch gibt es keinen CI-Lauf |

### 5. Aussage der Messprotokolle zu `test-docker` (Abschnitt G)

Stimmt mit den Rohdaten überein. `roh-lokal-rest.txt`: zwei lokale Läufe Exit 2 (652,0 s, 646,1 s); `PASS=1621 FAIL=9 SKIP=1`; neun `--- FAIL`-Zeilen mit den im Protokoll genannten Namen; Docker-Modus am Host `exit=0`, `PASS=1627 FAIL=0 SKIP=4`, die vier `SKIP` sind genau `NestedPodman`, `NestedPodman_Build`, `NestedPodman_BuildEgress` und `StatusFiles_VFS`, die lokal unter den neun fehlschlagenden stehen (Summen passen: 1621 + 9 + 1 = 1627 + 4). Die beiden lokalen Läufe zählen verschieden (5025 gegen 1621); das Protokoll sagt das. Die Behauptung „die CI auf `main` (Stand `b371e262`) meldet `integration-docker` grün“ stimmt (`gh run list`: Lauf 41188840564, success, 14 min 40 s, `headSha` `b371e262…`), sie belegt den Docker-Modus vor dem Umbau des `Dockerfile`; für den Stand `98699f3e` gibt es keinen CI-Lauf. `test-docker-tools` und der Docker-Pfad von `test-docker` sind im Diff nicht berührt, daher ist der Schluss plausibel, aber er ist keine Messung (F-7). Die ungeklärten Ursachen der fünf Tests sind ehrlich als offen geführt.

### 6. Befunde

| ID | Schwere | Befund | Fundstelle | Folge |
|---|---|---|---|---|
| F-1 | MEDIUM | **Der Devcontainer selbst ist nicht gestartet.** Die DoD-Formulierung (Liefer-Punkt 2) und der Closure-Trigger (§5) sagen „`make gates` im Devcontainer ist grün“. Gemessen ist das Image aus `.devcontainer/Dockerfile` per `docker build` ohne die Features aus `devcontainer.json` (`node`, `claude-code`, `homebrew`), privilegiert gestartet, mit Bind-Mount statt Workspace-Volume (Messprotokoll, Grenzen). Offen bleibt damit: ob `u-boot devcontainer up --recreate` mit dem neuen Dockerfile baut (Risiko `node:2.1.1` mit `claude-code:1.0.5`), ob `/workspaces/u-boot/.cache` im Volume dem Benutzer `vscode` gehört und schreibbar ist (das Image legt nur `/workspaces/u-boot` an), ob `GOTOOLCHAIN=local` und die Feature-Reihenfolge `PATH` unverändert lassen. Ich habe es nicht gefahren (Hostabhängig, lange, braucht das `u-boot`-Binary und die Features; kein Docker-Pfad des Slice-Codes betroffen), halte es aber für relevant, weil genau der Bau-Pfad des Slice-Ziels dort läuft und das Risiko in §6 sonst keinen belegten Ausgang bekommt. | Slice §2, §5, §6; Messprotokoll „Grenzen“ | Vor der Closure einen Lauf von `u-boot devcontainer up` mit `make gates` darin nachholen, oder in der Closure-Notiz den Rest ausdrücklich als „weiter offen“ führen und die DoD-Formulierung entsprechend benennen (nicht stillschweigend abhaken) |
| F-2 | LOW | **Abweichung A (`test-docker`)**, siehe Abschnitt 3: ADR-Entscheidung 3 führt `test-docker` als lokal laufend, ohne die Nicht-Gleichwertigkeit zu kennen; die Spezifikation und das Handbuch tragen sie jetzt. Der Re-Evaluierungs-Trigger der ADR nennt den Fall nicht. | [`ADR-0037`](../plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md) Entscheidung 3; [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) Punkt „Lokaler Modus“ | Keine Folge-ADR. Beobachtung im Register in der Closure (Plan kündigt sie an; Nachlauf N-6 stellt fest, dass sie noch fehlt) |
| F-3 | LOW | **Plan-Aussagen über Code, die nicht (mehr) stimmen oder fehlen:** die Planzeile „Release-Artefakte (`build-binaries`, `packages`) entstehen nur im Docker-Modus“ (es gibt eine lokale Variante von `build-binaries`); `.github/workflows/publish.yml` (ruft `make build-binaries` und `make build`) steht nicht in der „unverändert“-Zeile, obwohl die Aussage „Release-Pipeline im Docker-Modus“ in der Spezifikation davon abhängt. | Slice §3, Zeilen `spec/spezifikation.md` und Workflow-Zeile; `.github/workflows/publish.yml` Zeile 177 | Planzeile berichtigen und `publish.yml` mit gelesener Stelle in die „unverändert“-Zeile aufnehmen (AGENTS.md §3.10, vierter Punkt) |
| F-4 | LOW | **Zwei Dateien im Diff ohne Plan-Zeile:** `docs/plan/planning/open/README.md` und der Entwurf unter `open/` (Commit `ce61058a`, Entwurf von [`slice-v1-devcontainer-instanz-label`](../plan/planning/open/slice-v1-devcontainer-instanz-label.md) zu [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Ein fremder Vorgang liegt in der Range dieses Slice; AGENTS.md §3.10 verlangt jede berührte Datei in der Plan-Tabelle. | `git diff --stat 9b01cb44..98699f3e` | Entweder eine Planzeile „fremder Entwurf, nicht Teil des Slice“ oder die Range des Closure-PR so ziehen, dass der Entwurf als eigener Vorgang erkennbar bleibt; kein Verhaltensbefund |
| F-5 | LOW | **Bedingungen in Spec ohne pinnenden Test** (AGENTS.md §3.10, letzter Punkt): (a) die Befehle der Stages stehen in `scripts/bau.sh`, das auch das `Dockerfile` ruft (nur mittelbar durch `make gates`/CI gedeckt); (b) „keine Erkennung von `go` im `PATH`“; (c) die in [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) genannten lokalen Befehle (`modus-pruefung` prüft nur das Fehlen von `docker build`, nicht den Befehl); (d) `RUN` im Docker-Modus ist nicht geprobt, nur `PKG`. Jede ist heute wahr (gelesen), aber nicht gepinnt. | [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), [`SPEC-008`](../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), `repo.mk` `modus-pruefung` | Je Bedingung einen `make -n`-Grep in `modus-pruefung` ergänzen oder die Aussage als Nicht-Zusage kennzeichnen |
| F-6 | INFO | Risiko 3 im Plan beschreibt „Rezepte lehnen ab oder setzen leer“; umgesetzt ist die Ablehnung beim Einlesen über `MAKECMDGOALS` für vier Ziele. Ein Teillauf bleibt mit anderen Zielen (`make lint test PKG=…`) erlaubt, ohne Nachweis zu sein; das deckt die Spezifikation („ist kein Nachweis“). | `repo.mk` Zeilen 267 bis 276 | Plan-Wortlaut angleichen |
| F-7 | INFO | Messprotokoll F: Die Zeile „Coverage 92,10 %“ steht nach einem Absatz außerhalb der Tabelle und rendert als Fließtext. Die Aussage, der Docker-Modus bleibe nach dem Umbau gleich, stützt sich auf einen lokalen Lauf (91,1 s) und nicht auf die CI; für `98699f3e` gibt es keinen CI-Lauf. Beides ändert das Verdikt nicht. | [Messprotokoll](../messungen/slice-v1-devcontainer-bau-lokal/README.md) Abschnitt F | Tabellenzeile zurücksetzen; CI-Lauf des Pull Requests ist ohnehin Closure-Bedingung (§5) |

Nicht als Befund geführt: Das Pinnen der Prüfsummen bleibt an drei Stellen ohne Abgleich (Risiko 2, vom Projektinhaber bewusst nicht gelöst, Beobachtung in der Closure); `govulncheck` lokal lädt sein Werkzeug über das Netz (akzeptiertes Negativ, Risiko 4); der lokale `test-docker` setzt kein `CGO_ENABLED=0` (der Docker-Modus-Pfad auch nicht, gleiches Verhalten).

Carveouts:
- Neu: none (keine Gate-Lockerung, `.d-check.yml` nimmt nur `.cache/**` aus dem Scan, siehe Abschnitt 3 Entscheidung 8)
- Gelöst: none
- Unverändert: alle bestehenden

Nicht ausgeführt:
- `make test-docker` — Docker-Pfad (`test-docker-tools`, Rezept) im Diff nicht berührt; ein Lauf dauert etwa 27 Minuten; Beleg stattdessen die Rohdaten des Implementers und der CI-Lauf auf `main` vor dem Umbau (F-7).
- `make verify-depguard` — keine Importregeln berührt (kein Go-Quelltext im Diff).
- `make ci`, `make govulncheck`, `make image-scan`, `make fullbuild` — Security-/Release-Pfad nicht berührt; die lokale Variante von `govulncheck` ist vom Implementer gemessen (Protokoll G, Exit 0).
- `u-boot devcontainer up` mit dem neuen `.devcontainer/Dockerfile` — nicht gefahren, siehe F-1.
- Mutationsproben an `modus-pruefung` — nicht selbst wiederholt; der Nachlauf-Report belegt sie; ich habe fünf Bedingungen durch `make -n` von Hand bestätigt (Abschnitt 1).

Commit / Artefakt: `98699f3e` (Range `9b01cb44..98699f3e`).
