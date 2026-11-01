# Slice V1: Docker-Tests in die CI — der volle Lauf von `make test-docker` ist der Nachweis der CI, lokal genügt der engste Docker-Test, Wege A bis C zur Entscheidung des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Entwurf; Priorisierung und Einordnung beim Projektinhaber.

**Bezug:** [`ADR-0004`](../../adr/0004-ci-system.md) (CI-System; nennt `integration.yml` als eigenen Workflow), [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (lokaler Modus, CI unverändert). Kein `LH-*` im Scope: der Slice ändert Prozess und Harness-Vertrag, kein Produktverhalten.

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) (Komposition der CI-Pipeline; die Pflicht-Jobs sind drei, `integration-docker` steht nicht darunter) · [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Regeln des Docker-only-Workflows; der lokale Modus sagt die Docker-Tests als Nachweis nicht zu, Nachweis sind der Docker-Modus an einem Host mit lokaler Engine und die CI) · [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (nur Weg D unten: Parameter von `test-docker`).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Projektinhaber hat am 2026-11-01 entschieden, den vollen Lauf von `make test-docker` nicht mehr bei jedem Slice lokal zu fahren, sondern der CI zu überlassen. Dieser Slice macht daraus einen Vertrag des Harness: Er legt fest, was der lokale Nachweis eines Slice mit Docker-Verhalten ist, was in der CI gelten muss, damit deren Lauf Nachweis sein kann, und ob die Verlagerung ohne ADR geschieht.

**Der Slice wählt nicht.** Die Wege unten sind Fragen und Alternativen; die Wahl trifft der Projektinhaber vor dem Start (`open` → `next`). Bis dahin ist kein Satz dieses Entwurfs Zusage.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Laufzeit der Tests senken (Parallelisierung, Kürzung, Image-Cache) — das ist ein eigener Vorgang mit eigener Messung; hier wird nur die Zeitanalyse (§3, Frage 5) als Grundlage gelesen. Eine Verkürzung ändert den Gegenstand der Prüfung, die Verlagerung nicht.
- Das Zeitlimit `-timeout 30m` in `repo.mk` und `timeout-minutes: 30` in `integration.yml` ändern — beides ist Gate-Gegenstand ([`AGENTS.md`](../../../../AGENTS.md) §3.6) und Gegenstand der offenen Beobachtung [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md); der Slice nennt die Wechselwirkung (§6), entscheidet sie nicht.
- Die Nachweis-Regel für die Docker-Tests im verschachtelten Daemon des Devcontainers — sie steht fest in [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und in [`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`](../observations/BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig/observation.md); der Slice ändert sie nicht.
- Einen Release-Cut ändern — [`releasing.md`](../../../maintainer/releasing.md) §4 führt `make test-docker` in der Vorab-Prüfung; ob dort der volle lokale Lauf bleibt, ist Frage 1 (Folgepflicht), nicht ein eigener Slice, solange die Antwort eine Zeile ist.

## 2. Definition of Done

Höchstens drei Liefer-Punkte; welche davon gelten, hängt vom gewählten Weg ab (§3).

- [ ] Der Slice-Nachweis für Docker-/Compose-/E2E-Verhalten steht in [`harness/verification.md`](../../../../harness/verification.md) (Tabelle *Sensor-Auswahl*) so, wie der Projektinhaber ihn entschieden hat, und `harness/README.md` (Sensors-Tabelle) und [`quality.md`](../../../maintainer/quality.md) §2.2 sagen dasselbe.
- [ ] Die Bedingung, unter der ein CI-Lauf von `integration-docker` Nachweis ist, ist festgelegt und — soweit sie im Repo liegt (Workflow, Dokumentation) — umgesetzt; was im GitHub-UI liegt (Branch-Schutz), ist als Schritt benannt und vom Projektinhaber bestätigt.
- [ ] Die Frage, ob eine ADR nötig ist, ist beantwortet: entweder Folge-ADR (Rolle Architect, eigener Slice, Vorbedingung dieses Slice) oder begründet keine.
- [ ] `make gates` grün (`make docs-check` 0 Befunde).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt: nur Wartungsdoku (`docs/maintainer/`), keine Nutzer-Doku.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben (§8).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

### Befunde (gelesen, mit Fundstelle)

- `harness/verification.md`, Tabelle *Sensor-Auswahl*, Zeile *Docker-/Compose-/E2E-Verhalten*: Mindest-Sensor „engster Docker-Test oder `make test-docker`“, normaler Closure-Sensor „`make test-docker` plus `make gates`“.
- `repo.mk`, Ziel `test-docker`: Docker-Modus baut die Stage `test-docker-tools` und fährt `go test -v -tags docker -timeout 30m ./...` im Container mit `--network=host`, Socket- und `/tmp`-Mount; lokaler Modus fährt denselben `go test`-Aufruf gegen den Daemon der Umgebung. Der Befehl ist in beiden Modi fest `./...`, ohne `PKG`/`RUN`.
- `repo.mk` (Zeilen um 82 bis 93, `TEST_PKG`/`TEST_RUN`): `PKG`/`RUN` sind nur im lokalen Modus erlaubt und wirken nur auf `test`; im Docker-Modus ist beides ein Fehler, und `gates`, `record-gates`, `ci`, `fullbuild` lehnen sie ab. Das prüft `make modus-pruefung` (`harness/README.md`, Sensors-Tabelle).
- `.github/workflows/integration.yml`: Trigger `pull_request` (Ziel `main`), `push` (`main`) und `workflow_dispatch`; ein Job `integration-docker` (`name: integration-docker (//go:build docker)`), `timeout-minutes: 30`, Schritt `make test-docker`. **Der Kopfkommentar nennt `continue-on-error: true` als befristete Maßnahme; im Job steht das Flag nicht** (`grep continue-on-error .github/` findet nur die zwei Kommentarzeilen 13 und 22). Der Commit `8865ca13` („remove integration-docker continue-on-error flag“) hat es entfernt. Der Kommentar beschreibt damit einen Zustand, der nicht mehr gilt (§3.7 verletzt) und ist Teil dieses Slice. Der Eintrag zu `slice-m6-docker-integrationstests` im Carveout-Inventar (`in-progress/carveouts.md`, Zeile 58) steht auf *Done*.
- `docs/maintainer/branch-protection.md`: die Liste der Required Status Checks nennt drei Namen aus `ci.yml` (`gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`, `image-scan (trivy HIGH+CRITICAL)`); `integration-docker` steht nicht darunter. Ob die Regel im GitHub-UI so gesetzt ist, ist aus dem Repo nicht lesbar und nicht geprüft.
- [`ADR-0004`](../../adr/0004-ci-system.md), Abschnitt *Bewusst noch nicht Teil der CI-Grundentscheidung*: Adapter-Integrationstests gegen die externe Docker-Engine sind ein eigener Workflow `integration.yml`; die ADR sagt dort nichts über Blockierung. [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions): „Drei Jobs, parallel, alle PR-blockierend“ (die Jobs von `ci.yml`).
- [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), Absatz *Lokaler Modus*: die Docker-Tests (`test-docker`) sind im verschachtelten Daemon nicht als Nachweis zugesagt; den Nachweis liefern der Docker-Modus an einem Host mit lokaler Engine und die CI. Die CI ist damit in der Spec schon als Nachweis-Ort genannt.
- Beobachtungen im Register: [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (Schwelle von drei erreicht, Ausgang beim Projektinhaber; lokale Läufe endeten dreimal am Zeitlimit von 30 min, CI-Läufe des Pakets 673 bis 885 s) und [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (ob die Build-Tests auf `ubuntu-latest` überspringen, ist dort ungemessen). Beide betreffen die Frage, was ein grüner CI-Lauf belegt.

### Zeitanalyse (Frage 5, gelesen)

Gelesen ist `test-docker-2.log` im Scratchpad der Session (Änderungszeit 2026-11-01 03:56; der in der Anfrage genannte Name `test-docker.log` existiert nicht). Das Log endet mit `PASS`, enthält kein `--- FAIL`; die Annahme, dass es der grüne Lauf des Verifiers ist, ist **nicht belegt**. Aus dem Log gelesen:

- Paketzeiten: `internal/e2e` 1540,129 s, `internal/adapter/driven/docker` 165,130 s, `internal/hexagon/application` 16,837 s (die Pakete laufen parallel gegen dieselbe Engine; die Gesamtzeit ist nicht die Summe).
- `internal/e2e`: 37 Tests der Form `TestE2E_*` auf oberster Ebene, Summe ihrer Einzelzeiten 1471,9 s. Die 15 langsamsten summieren sich auf 927,5 s, die fünf langsamsten sind `…ForgeCLI_BuildsInBothProfiles` (106,6 s), `…SandboxDevcontainer_Ollama` (92,3 s), `…NestedDocker_Degradation` (73,1 s), `…NestedPodman_BuildEgress` (71,5 s) und `…NestedDocker` (68,7 s).
- Vier Tests enden im `SKIP` und kosten dennoch 236,6 s (`…NestedPodman_Build` 57,2 s, `…NestedPodman_BuildEgress` 71,5 s, `…NestedPodman` 58,1 s, `…StatusFiles_VFS` 49,8 s): die Bedingung des Skips wird erst nach dem Aufbau geprüft. Der Skip-Grund ist im Log nicht ausgewertet.
- Parallelität: in `internal/e2e` steht `t.Parallel()` nur in `sandbox_status_docker_test.go` (drei Stellen: Zeilen 163, 232, 302).

Nicht gelesen und nicht übernommen: die Anzahl „56 Docker-Tests“ der Anfrage (gezählt wurden nur die 37 oberster Ebene in `internal/e2e`) und der Anteil der Zeit auf Image-Bau gegenüber Test-Logik.

### Wege und Fragen (zur Entscheidung vor dem Start)

**Frage 1 — Slice-Nachweis lokal (`harness/verification.md`).** Die Zeile *Docker-/Compose-/E2E-Verhalten* wird umgestellt. Was ist dann der Nachweis?

- **Weg A:** lokal der engste Docker-Test (Familie der geänderten Stelle) plus `make gates`; der volle Lauf ist der CI-Lauf `integration-docker` auf dem Push-Tip von `main`, nach dem Push gelesen. Der Slice geht bei Docker-Verhalten erst nach grünem CI-Lauf nach `done/` — oder mit dem Lauf als nachgetragener Evidence (offene Frage: Reihenfolge Closure und Push; der Slice-Zyklus schließt heute vor dem Push).
- **Weg B:** lokal der engste Docker-Test plus `make gates`; der CI-Lauf ist kein Closure-Kriterium, sondern ein Befund danach: ein roter Lauf öffnet einen Folge-Slice oder ein Finding. Kleinster Eingriff, aber der Nachweis ist dann lokal und teilweise.
- **Weg C:** wie A, zusätzlich läuft `integration-docker` bei Pull-Requests als Required Check (Frage 2); der Push erfolgt über einen Branch, der CI-Lauf ist vor dem Merge grün. Die schärfste Form; sie verlangt den PR-Weg, den der Solo-Betrieb heute nicht erzwingt (nicht geprüft, wie der Projektinhaber auf `main` schreibt).

Folgepflichten aus jedem Weg: `releasing.md` §4 (Zeile `make test-docker` in der Vorab-Prüfung), `harness/README.md` (Zeile `make test-docker`), `docs/maintainer/quality.md` §2.2 und die Beobachtung `BEO-E2E/e2e-paket-am-zeitlimit` (der Slice weist ihr einen Ausgang zu oder vermerkt, dass er den lokalen Lauf als Nachweis entfernt).

**Frage 2 — Was muss gelten, damit ein CI-Lauf Nachweis ist?** Befund oben: das Flag `continue-on-error` ist nicht gesetzt, der Job ist also hart, aber nicht Required Check. Zu entscheiden:

- (a) Der Kopfkommentar in `integration.yml` wird auf den Ist-Zustand gebracht (Zusage statt Chronik), das Flag steht nicht mehr da.
- (b) Der Job wird Required Status Check in `main` (Name `integration-docker (//go:build docker)`) — Schritt im GitHub-UI, dokumentiert in `branch-protection.md`; im Repo nicht prüfbar. Ohne (b) ist ein roter Lauf sichtbar, blockiert aber nichts.
- (c) Stabilität: die Läufe waren nicht durchweg grün (die Beobachtungen nennen einen roten Lauf wegen eines Docker-Tests des Adapters und einen wegen eines Container-Zählers paralleler Pakete). Welche Anzahl grüner Läufe in Folge vor (b) verlangt der Projektinhaber? Der frühere Maßstab des Workflow-Kommentars war drei aufeinanderfolgende grüne Läufe auf `main`.
- (d) Skips: ein grüner Lauf mit übersprungenen Tests belegt diese nicht ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)). Gilt „CI grün“ als Nachweis für Tests, die dort überspringen, oder bleibt für sie ein gesonderter Weg? Ohne Antwort verlagert die Entscheidung eine Prüfung in einen Lauf, der sie nicht ausführt.
- (e) Zeitreserve: Job-Limit `timeout-minutes: 30` bei bisher 12 bis 16 Minuten Joblaufzeit (Belege der Beobachtung `e2e-paket-am-zeitlimit`); wird die CI zum einzigen vollen Lauf, trägt jede Verlängerung der Tests die Gefahr eines Abbruchs ohne roten Test. Ein Hinweis, keine Änderung (§1).

**Frage 3 — Braucht es eine ADR?** [`AGENTS.md`](../../../../AGENTS.md) §3.6 verlangt für das Lockern eines Gates einen Plan- oder ADR-Anker. Einschätzung des Planners, zu bestätigen vom Architect: Die Schwelle der Prüfung (alle Docker-Tests müssen grün sein) bleibt; sie wandert vom Slice-Handoff in die CI. Das ist keine Senkung, **wenn** die CI-Prüfung hart ist — also nur zusammen mit Frage 2(b) und (d). Ohne (b) wäre die Verlagerung faktisch eine Lockerung und bräuchte den ADR-Anker. Unabhängig davon ändert Weg C die Annahme „Required Checks sind die drei Jobs von `ci.yml`“ aus [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) und [`ADR-0004`](../../adr/0004-ci-system.md) — dort wäre eine Folge-ADR (Ergänzung zu [`ADR-0004`](../../adr/0004-ci-system.md), wie [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) sie ließ) sauberer als ein Plan-Anker. Weg A/B: Plan-Anker in diesem Slice genügt nach dieser Einschätzung, wenn Frage 2(b) mitgeliefert wird. Die ADR-Frage entscheidet der Architect, nicht dieser Plan.

**Weg D (Frage 4, optional) — Filter für den engsten Test.** `make test-docker` bekommt einen Filter `RUN` (und `PKG`), damit der engste Docker-Test lokal ohne Umgehung des Make-Targets läuft. Befund: heute ist `PKG`/`RUN` für `test-docker` nicht definiert, im Docker-Modus ein Fehler (`repo.mk`, `modus-pruefung`); der Docker-Test-Lauf des engsten Tests geschieht daher mit `go test -tags docker -run …` im Container des Targets von Hand — das ist nach [`AGENTS.md`](../../../../AGENTS.md) §3.1 ohne Nachweiskraft. Folgen: `repo.mk` (beide Modi), `scripts/bau.sh` (nur lokal, nicht gelesen), `make modus-pruefung` (neue Gegenproben: `RUN` im Docker-Modus bei `test-docker`), [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) und [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Zusage zum Filter, mit Test nach §3.10), `gates`/`ci`/`fullbuild` lehnen `PKG`/`RUN` weiter ab. Ohne Weg D bleibt der „engste Docker-Test“ in Weg A bis C ohne make-konformen Aufruf. Ein Entscheidungspunkt: Weg D gehört entweder in diesen Slice (dann ist er der dritte Liefer-Punkt) oder in einen eigenen Folge-Slice.

### Plan-Tabelle (Kandidaten; gelesen, was jeweils berührt wird)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/verification.md` (Tabelle *Sensor-Auswahl*) | update | Zeile *Docker-/Compose-/E2E-Verhalten* nach Frage 1 |
| `harness/README.md` (Sensors-Tabelle, Zeile `make test-docker`) | update | Bindung/Vertrag nennt den CI-Lauf als Nachweis (nach Weg) |
| `docs/maintainer/quality.md` §2.2 | update | Beschreibung des Integration-Loops und des CI-Laufs |
| `docs/maintainer/releasing.md` §4 | update (je Weg) | Zeile `make test-docker` der Vorab-Prüfung |
| `.github/workflows/integration.yml` | update | Kopfkommentar auf den Ist-Zustand (kein Flag); Zeitlimit unverändert (§1) |
| `docs/maintainer/branch-protection.md` | update (Weg C / Frage 2(b)) | Required Check `integration-docker (//go:build docker)` |
| `docs/plan/planning/observations/BEO-E2E/e2e-paket-am-zeitlimit/` | Ausgang/Evidence | Ausgang zuweisen oder fortschreiben |
| `repo.mk`, `scripts/bau.sh` | update (nur Weg D) | Filter `RUN` für `test-docker` |
| `spec/spezifikation.md` ([`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)) | update (nur Weg D) | Zusage zum Filter, nur mit Test |
| `Makefile:modus-pruefung` (in `repo.mk`) | update (nur Weg D) | Gegenproben für `RUN` bei `test-docker` |
| ADR (Folge von [`ADR-0004`](../../adr/0004-ci-system.md)) | neu (nur wenn der Architect es bestätigt) | Vorbedingung, eigener Slice |

Nicht gelesen und deshalb in der Tabelle nicht als „unverändert“ geführt: `.harness/skills/*.md` (keine Zeile zu `test-docker` gefunden), `harness/roles.md` und `harness/review.md` (keine Fundstelle zu `test-docker`); ob der Verifier-Auftrag (`.claude/`) den Sensor nennt, ist nicht geprüft.

## 4. Trigger

**Start** (`next` → `in-progress`): Der Projektinhaber hat Weg A, B oder C zu Frage 1 gewählt, Frage 2 (a) bis (e) beantwortet, über Weg D entschieden und die ADR-Frage dem Architect gestellt; ist eine Folge-ADR nötig, ist sie `Accepted`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Weg D wächst zu einem dritten und vierten Liefer-Punkt — dann Weg D als eigener Slice.
- `in-progress` → `open` (blockiert): die Required-Check-Einrichtung im GitHub-UI ist nicht möglich oder nicht bestätigt, oder die Läufe sind nicht stabil (Frage 2(c)); dann bleibt der Nachweis lokal und die Entscheidung ruht.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Bei Weg A oder C gehört zum Beleg ein grüner CI-Lauf `integration-docker` auf dem Push-Tip, in der Evidence mit Lauf-Kennung genannt. Lerneintrag: die Regel zum Slice-Nachweis bei Docker-Verhalten (geschärfte Regel in `harness/verification.md`) und der Befund zum veralteten Workflow-Kommentar (Kommentar folgte dem Zustand nicht).

## 6. Risiken und offene Punkte

- Der CI-Lauf ist kein harter Nachweis, solange `integration-docker` kein Required Check ist: ein roter Lauf blockiert nichts, die Verlagerung wäre dann eine stille Lockerung (§3.6) — **Ausgang:** offen bis zur Entscheidung zu Frage 2(b) und 3.
- Übersprungene Tests belegen nichts; die vier `SKIP`-Tests des gelesenen Laufs (`NestedPodman_Build`, `…_BuildEgress`, `NestedPodman`, `StatusFiles_VFS`) würden in der CI womöglich ebenfalls überspringen — **Ausgang:** offen → [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md).
- Ein Docker-Fehler wird später entdeckt (erst nach dem Push statt vor dem Handoff), der Fix-Zyklus verlängert sich um einen CI-Durchlauf von etwa 13 bis 16 Minuten — **Ausgang:** offen; Mitigation ist der engste Docker-Test lokal (Weg A bis C).
- Das Zeitlimit von 30 Minuten (Job und `go test`) liegt bei einem vollen Lauf, der lokal 1540 s im Paket `internal/e2e` brauchte, dicht an der Grenze; in der CI sind es bisher 12 bis 16 Minuten — **Ausgang:** offen → [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md).
- Die Anzahl der Docker-Tests (Anfrage: 56) und die Identität des gelesenen Logs sind nicht verifiziert (§3, Zeitanalyse) — **Ausgang:** offen; der Implementer zählt und misst vor dem Start neu.

## 7. Closure-Notiz

Noch nicht geschlossen. Entwurf; wird mit der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `harness/` und `docs/plan/` (DOK, Greenfield) und `docs/maintainer/` (liegt in der Zeile DOK nicht ausdrücklich; Zuordnung zu DOK als Doku-Pfad ohne eigene Zeile, nicht neu geschnitten), `.github/workflows/` (keine eigene Sub-Area-Zeile; Gate-Datei nach Source Precedence Rang 6) und bei Weg D `tools/`/`scripts/` und `repo.mk` (WZG, Greenfield). Alle Zeilen laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006) Greenfield; keine zu grobe Sub-Area erkannt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (`BEO-E2E`, `BEO-WZG`, `BEO-DOK`, per Verzeichnisliste, die E2E-Beobachtungen gelesen). Treffer: [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (offen, 3×, Ausgang beim Projektinhaber), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen), [`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`](../observations/BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig/observation.md) (offen, 1×); beim `BEO-WZG`- und `BEO-DOK`-Bestand wurden nur die Namen gelesen, kein Treffer nach Titel. Ein Vermerk im Register fällt mit der Closure an.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
