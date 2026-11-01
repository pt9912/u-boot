# Slice V1: Den CI-Lauf `integration-docker` in mehrere Jobs aufteilen — jeder Job fährt einen Teil der Docker-Tests auf eigenem Runner und eigener Engine, Wege A bis C zur Entscheidung des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Entwurf; interner Werkzeug- und CI-Slice ohne eigene Closure-Bedingung ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)). Die Zuordnung zu einer Welle (Release-Version) ist offen und liegt beim Projektinhaber.

**Bezug:** [`ADR-0004`](../../adr/0004-ci-system.md) (CI-System; `integration.yml` als eigener Workflow), [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (lokaler Modus, CI unverändert). Kein `LH-*` im Scope: der Slice ändert Prozess und Harness-Vertrag, kein Produktverhalten. Anlass: Entscheidung des Projektinhabers, 2026-11-01 (Hebel 3 der CI-Zeitanalyse).

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) (Komposition der CI-Pipeline; Pflicht-Jobs sind die drei aus `ci.yml`) · [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Docker-only-Workflow; Docker-Tests im verschachtelten Daemon sind kein zugesagter Nachweis) · [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (nur, wenn `test-docker` Parameter bekommt, Weg A und B).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der CI-Lauf `integration-docker` besteht aus mehreren Jobs, die je einen Teil der Docker-Tests auf einem eigenen Runner mit eigener Engine fahren; die Gesamtzeit ist die des längsten Teils statt der Summe, und jede Prüfung läuft weiter hart.

**Der Slice wählt nicht.** Die Wege und Fragen in §3 sind Alternativen; die Wahl trifft der Projektinhaber vor dem Start (`open` → `next`). Bis dahin ist kein Satz dieses Entwurfs Zusage.

**Verhältnis zu [`slice-v1-test-docker-in-ci`](slice-v1-test-docker-in-ci.md):** Jener Entwurf legt fest, was der Nachweis eines Slice mit Docker-Verhalten ist und wann ein CI-Lauf Nachweis ist (Required Check, Skips). Dieser Slice ändert die Bauart des Laufs (Aufteilung). Beide hängen über die Required-Check-Frage (§3, Frage 4) und über das Zeitlimit zusammen; keiner setzt den anderen voraus. Die Reihenfolge entscheidet der Projektinhaber.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Hebel 1: die Skip-Bedingung der Nested-Podman-Tests vor den Image-Aufbau ziehen** — das ändert Testcode (Reihenfolge von Bau und Skip-Prüfung) und spart Zeit unabhängig von der Aufteilung (§3, Zeitwerte: vier `SKIP`-Tests kosten 86,0 s im CI-Lauf, 189,3 s im lokalen Log); ein eigener Vorgang mit eigener Messung. Ob er vor oder nach diesem Slice läuft, ist Frage 7.
- **Hebel 2: `t.Parallel()` und eindeutige Containernamen in den Tests** — ändert Testcode und die Isolation der Tests untereinander (`t.Parallel()` steht heute nur in `internal/e2e/sandbox_status_docker_test.go`, drei Stellen, gelesen per Suche); eigener Vorgang. Die Aufteilung auf getrennte Runner braucht ihn nicht.
- **Das Zeitlimit ändern** (`timeout-minutes: 30` in `integration.yml`, `-timeout 30m` in `repo.mk`) — Gate-Gegenstand ([`AGENTS.md`](../../../../AGENTS.md) §3.6) und Gegenstand der offenen Beobachtung [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md); der Slice nennt die Wechselwirkung (§6), entscheidet sie nicht. Jeder Teil-Job bekommt ein eigenes Limit; ob 30 Minuten bleiben, ist Frage 6.
- **Den Nachweis-Vertrag (lokal gegen CI) festlegen** — [`slice-v1-test-docker-in-ci`](slice-v1-test-docker-in-ci.md).
- **Tests streichen oder abschwächen** — die Aufteilung verschiebt, welcher Job einen Test fährt; die Menge der Tests und ihre Prüfungen bleiben (Vollständigkeits-Prüfung, DoD).
- **Produktcode** — nicht berührt.

## 2. Definition of Done

Höchstens drei Liefer-Punkte; die Wege A bis C (§3) entscheiden die Ausprägung.

- [ ] **Aufteilbares Target:** `make test-docker` fährt auf Anforderung einen Teil der Docker-Tests (Form nach Weg A oder B), im Docker-Modus und im lokalen Modus; der Aufruf ohne Teil bleibt der volle Lauf. Ein Test pinnt die Bedingungen von `make modus-pruefung` (neue Gegenproben).
- [ ] **Aufgeteilter Workflow:** `.github/workflows/integration.yml` fährt die Teile als eigene Jobs mit eigenem Runner; die Vereinigung der Teile ist die volle Menge (ein Sensor oder eine Prüfung im Lauf fängt ein Paket oder einen Test, den kein Teil fährt); der Required-Check-Vertrag ist nach Frage 4 festgelegt und in [`branch-protection.md`](../../../maintainer/branch-protection.md) beschrieben, soweit er im GitHub-UI liegt (dort als Schritt benannt).
- [ ] **Doku:** `harness/README.md` (Sensors-Tabelle, Zeile `make test-docker`), [`quality.md`](../../../maintainer/quality.md) §2.2 und, wenn betroffen, [`releasing.md`](../../../maintainer/releasing.md) §4 beschreiben die Aufteilung; ob eine ADR nötig ist, ist beantwortet (Frage 8).
- [ ] `make gates` grün (`make docs-check` 0 Befunde).
- [ ] Ein CI-Lauf mit allen Teil-Jobs ist grün; die Lauf-Kennung steht in der Closure-Notiz; die Gesamtzeit gegen die Zeiten aus §3 ist gemessen.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt: nur Wartungsdoku (`docs/maintainer/`), `CHANGELOG.md` nur, wenn der Projektinhaber es verlangt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§8); es wird kein Zähler gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

### Befunde (gelesen, mit Fundstelle)

- `.github/workflows/integration.yml`: Trigger `pull_request` (Ziel `main`), `push` (`main`), `workflow_dispatch`; ein Job `integration-docker` mit `name: integration-docker (//go:build docker)`, `runs-on: ubuntu-latest`, `timeout-minutes: 30`, Schritte Checkout und `make test-docker`. Der Kopfkommentar (Zeilen 13 und 22) nennt `continue-on-error: true` als befristete Maßnahme; im Job steht das Flag nicht ([`slice-v1-test-docker-in-ci`](slice-v1-test-docker-in-ci.md) §3 führt dasselbe als Teil seines Gegenstands).
- `repo.mk`, Ziel `test-docker`: Docker-Modus baut die Stage `test-docker-tools` und fährt `go test -v -tags docker -timeout 30m ./...` im Container mit `--network=host`, Socket- und `/tmp`-Mount (Zeilen 167 bis 201); lokaler Modus fährt denselben `go test`-Aufruf gegen den Daemon der Umgebung (Zeile 147 bis 148). Der Befehl ist in beiden Modi fest `./...`.
- `repo.mk` Zeilen 82 bis 93: `PKG` und `RUN` zählen nur von der Kommandozeile, sind im Docker-Modus ein Fehler und werden mit `gates`, `record-gates`, `ci`, `fullbuild` abgelehnt; `scripts/bau.sh` setzt `BAU_RUN` als `-run` (Zeilen 15 und 37 bis 38, nur dort gelesen; ob der Zweig `test` oder auch `test-docker` betrifft, ist nicht gelesen). `make modus-pruefung` pinnt das ([`harness/README.md`](../../../../harness/README.md), Sensors-Tabelle; [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md)). `Makefile` selbst enthält kein Ziel `test-docker` (Suche ohne Treffer).
- `docs/maintainer/branch-protection.md`: Required Status Checks sind die drei Jobs aus `ci.yml` (`gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`, `image-scan (trivy HIGH+CRITICAL)`, Namen in `ci.yml` Zeilen 40, 55, 70); `integration-docker` steht nicht darunter. Ob die Regel im GitHub-UI so gesetzt ist, ist aus dem Repo nicht lesbar und nicht geprüft.
- Beobachtungen: [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (3×, kein Ausgang; CI-Joblaufzeit 12 min 53 s auf dem Tag-Commit von `v0.13.0`), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (2×; vier Tests springen auf dem Runner, Grund `kernel.apparmor_restrict_unprivileged_userns=1`), [`BEO-E2E/test-zaehlt-container-paralleler-pakete`](../observations/BEO-E2E/test-zaehlt-container-paralleler-pakete/observation.md) (1×; ein Test zählte Container eines parallel laufenden Pakets), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2×).

### Zeitwerte (gelesen)

**CI-Lauf `integration-docker` 41285105546** auf `main` (Commit `a6264402`, 2026-11-01, grün, Joblaufzeit 15 min 14 s; `gh run list`, Job-Log): Paketzeiten `internal/e2e` 849,686 s, `internal/adapter/driven/docker` 190,026 s, `internal/hexagon/application` 21,769 s; alle übrigen Pakete unter einer Sekunde bis einigen Sekunden. Die Pakete laufen im selben Job gegen dieselbe Engine; die Joblaufzeit liegt über dem längsten Paket, weil Checkout, Image-Aufbau (`test-docker-tools`) und Übersetzen dazukommen. Der Wert „ca. 13 Minuten Jobzeit“ der Anfrage ist die Spanne der Beobachtung (12 min 53 s) und der letzten Läufe (12 bis 16 min), im Lauf oben 15 min 14 s.

**Lokales Log `test-docker.log`** (Scratchpad der Session, Änderungszeit 2026-11-01 05:03; **ein roter Lauf**: das Paket `internal/e2e` endet mit `FAIL`, drei Tests rot — `TestE2E_LHFAUP005_LogsTail`, `…_LogsFollow`, `TestE2E_SandboxDevcontainer_TwoInstancesIsolated` —, die Zeiten sind trotzdem die der gelaufenen Tests): 1651 Testzeilen oberster Ebene, Summe der Einzelzeiten 1528,2 s; 25 Tests mit mindestens 20 s ergeben 1121,9 s; Paketzeiten `internal/e2e` 1224,538 s, `internal/adapter/driven/docker` 176,144 s, `internal/hexagon/application` 17,453 s. Vier `SKIP`-Tests (`TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…_BuildEgress`, `…_NestedPodman`, `…_StatusFiles_VFS`) kosten dort 189,3 s, im CI-Lauf 86,0 s; die Skip-Bedingung wird erst nach dem Aufbau geprüft ([`slice-v1-test-docker-in-ci`](slice-v1-test-docker-in-ci.md) §3 nennt 236,6 s aus einem anderen Log). „Lokal ca. 26 Minuten“ der Anfrage ist **nicht** aus diesen Dateien gelesen.

**Folgerung für den Schnitt (aus den Zeiten, nicht gemessen):** Das Paket `internal/e2e` ist in beiden Quellen der längste Teil (850 s von 15 min im CI-Lauf). Eine Aufteilung nur nach Paketen ergibt Teile von rund 850 s (e2e) und rund 190 s (alles andere): die Gesamtzeit sinkt höchstens auf die Zeit des e2e-Teils plus Aufbau. Eine Aufteilung innerhalb von `internal/e2e` nach Testnamen drückt die Zeit weiter, kostet aber je Job den Aufbau von `test-docker-tools` erneut.

### Wege und Fragen (zur Entscheidung vor dem Start)

**Frage 1 — Schnitt der Teile.**

- **Weg A, nach Paketen:** Matrix über Paketgruppen (z. B. `internal/e2e` allein, `internal/adapter/driven/docker` allein, Rest). Einfach, `PKG` genügt als Parameter. Wirkung begrenzt auf die Differenz zwischen e2e und dem Rest (siehe Folgerung).
- **Weg B, nach Testnamen-Gruppen:** Matrix mit `-run` über Gruppen innerhalb von `internal/e2e` (Teile nach Familien, z. B. Devcontainer-Start/Shell/Instanz, Sandbox, Forge-CLI, Cluster/Compose). Größere Ersparnis; verlangt eine Gruppenzuordnung, die jeden Test genau einmal trifft, und eine Vollständigkeitsprüfung (eine neue Testfamilie darf nicht still in keinem Teil laufen).
- **Weg C, A und B kombiniert:** Pakete getrennt, `internal/e2e` zusätzlich nach Gruppen.

**Frage 2 — Filter für `make test-docker`.** Heute sind `PKG` und `RUN` im Docker-Modus ein Fehler (`repo.mk` Zeilen 88 bis 92; `make modus-pruefung`); die CI läuft im Docker-Modus. Ein Filter für `test-docker` braucht (a) eine Regel im Makefile, die `PKG`/`RUN` bei `test-docker` im Docker-Modus zulässt, (b) neue Gegenproben in `modus-pruefung`, (c) die Aussage in [`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) und [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), und zwar nur mit einem Test, der sie pinnt ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Ein Teillauf ist kein Nachweis ([`AGENTS.md`](../../../../AGENTS.md) §3.1): die Teile dürfen nicht mit `gates`, `ci`, `fullbuild` kombinierbar sein. Alternative: ein eigener Parameter (etwa `TEST_GROUP`) statt der Wiederverwendung von `PKG`/`RUN`, damit die bestehende Ablehnung unverändert bleibt. Ein direkter `go test -run …`-Aufruf im Workflow ist nach [`AGENTS.md`](../../../../AGENTS.md) §3.1 keine Lösung.

**Frage 3 — Kollisionen bei gleicher Engine.** Die Annahme der Anfrage: Mit getrennten Runnern entfallen Kollisionen (feste Containernamen) zwischen Teilen, weil jeder Job seine eigene Engine hat. Das ist **nicht gelesen**: Dass Tests feste Container- oder Projektnamen teilen, ist an keiner Stelle belegt (gefunden wurden Namen mit Zufallssuffix wie `dcSuffix()` und Präfixe wie `ubl-shellguard-`, `uboot-e2e-cluster-`, jeweils nur per Suche). Der Implementer prüft das vor dem Start; falls sie gelten, ist die Aufteilung auch der Hebel, der sie in der CI entschärft, ohne Hebel 2. Innerhalb eines Jobs bleibt die Beobachtung [`BEO-E2E/test-zaehlt-container-paralleler-pakete`](../observations/BEO-E2E/test-zaehlt-container-paralleler-pakete/observation.md) bestehen, solange mehrere Pakete in einem Teil laufen.

**Frage 4 — Required Check.** Heute ist `integration-docker` kein Required Check (`branch-protection.md`). Mit einer Matrix ist zu entscheiden:

- jeder Teil-Job einzeln als Required Check (Namen bei einer Matrix kommen aus der Matrix-Zeile; jede Änderung der Teile ändert die Namen und erfordert eine Anpassung im GitHub-UI), oder
- ein Sammel-Job (`needs:` auf alle Teile, läuft auch bei Rot der Teile, schlägt fehl, wenn einer fehlschlägt oder übersprungen wurde), der als einziger Required Check eingetragen ist; der Name bleibt bei einer Änderung der Teile stabil.

Der Planner empfiehlt nichts; die Wahl hängt an der Entscheidung zu [`slice-v1-test-docker-in-ci`](slice-v1-test-docker-in-ci.md) Frage 2(b) (ob `integration-docker` überhaupt Required wird). Ein Sammel-Job ist ein neuer Job in `integration.yml`; sein Verhalten bei übersprungenen Teilen ist ein Test der Workflow-Logik und nicht lokal lauffähig (Grenze, §6).

**Frage 5 — Runner-Minuten.** Die Summe der Minuten über alle Jobs steigt (jeder Teil zahlt Checkout und Image-Aufbau erneut), die Wanduhr-Zeit sinkt. Ob das Repo auf Runner-Minuten begrenzt ist, ist **nicht gelesen** (Sichtbarkeit und Plan des Repos stehen in keiner gelesenen Datei). Der Projektinhaber nennt die Grenze oder bestätigt, dass sie keine ist.

**Frage 6 — Zeitlimit je Teil.** Soll `timeout-minutes: 30` je Teil bleiben oder auf die erwartete Teilzeit gesenkt werden (eine Senkung ist Gate-Änderung im Sinn von [`AGENTS.md`](../../../../AGENTS.md) §3.6 und braucht den Plan-Anker)? Die Aufteilung gibt dem Limit mehr Reserve (der längste Teil liegt in beiden gelesenen Quellen unter dem Limit).

**Frage 7 — Hebel 1 und 2.** Out-of-Scope (§1). Zu entscheiden ist die Reihenfolge: Hebel 1 (Skip-Bedingung vor den Aufbau, 86 s im CI-Lauf) senkt die Zeit auch ohne Aufteilung und verändert die Zusammensetzung der Teile (Weg B); soll er vor diesem Slice laufen oder als eigener Slice danach? Hebel 2 ändert die Isolation der Tests und gehört in einen eigenen Entwurf.

**Frage 8 — ADR.** [`AGENTS.md`](../../../../AGENTS.md) §3.6: Gates dürfen nicht ohne Plan- oder ADR-Anker gelockert werden. Einschätzung des Planners, zu bestätigen vom Architect: Das Aufteilen verschiebt, welcher Job einen Test fährt; die Menge der Tests und ihre Pflicht (alle grün) bleiben, also keine Senkung — **nur wenn** die Vereinigung der Teile geprüft wird (DoD, Vollständigkeit) und der Required-Check-Vertrag (Frage 4) keinen Teil unverbindlich lässt. Eine neue Aussage zu den Required Checks ändert die Annahme in [`SPEC-026`](../../../../spec/spezifikation.md#spec-026--komposition-der-ci-pipeline-github-actions) (drei Jobs) nur, wenn `integration-docker` Required wird; das ist die Sache der Entscheidung im Nachbar-Entwurf. Ob es eine Folge-ADR zu [`ADR-0004`](../../adr/0004-ci-system.md) braucht, entscheidet der Architect.

### Plan-Tabelle (Kandidaten; gelesen, was jeweils berührt wird)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/integration.yml` | update | Matrix oder mehrere Jobs; ggf. Sammel-Job; Kopfkommentar auf den Ist-Stand (kein `continue-on-error` im Job); Zeitlimit je Teil nach Frage 6 |
| `repo.mk` (Ziel `test-docker`, beide Modi; Zeilen 82 bis 93, 147 bis 148, 167 bis 201, `MODUS_TARGETS` Zeile 317; Gegenproben von `modus-pruefung`) | update | Parameter für einen Teil (Frage 2), Ablehnung mit `gates`/`ci`/`fullbuild` bleibt |
| `scripts/bau.sh` | update (lokaler Modus, falls der Parameter dort ankommt) | `BAU_RUN`/`BAU_PKG` werden gelesen (Zeilen 15, 37 bis 38); ob `test-docker` lokal über das Skript geht, ist nicht gelesen |
| `Dockerfile` (Stage `test-docker-tools`) | nicht gelesen | jeder Teil-Job baut die Stage neu; ob Cache über Jobs hinweg möglich ist, ist offen (§6) |
| `harness/README.md` (Sensors-Tabelle, Zeile `make test-docker`) | update | Aufteilung und Parameter |
| `docs/maintainer/quality.md` §2.2, `docs/maintainer/branch-protection.md`, `docs/maintainer/releasing.md` §4 | update (je Frage 4 und Weg) | Beschreibung des Integration-Loops; Required-Check-Namen; Vorab-Prüfung |
| `spec/spezifikation.md` ([`SPEC-008`](../../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles), [`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)) | update, nur wenn ein Test die Aussage pinnt | Zusage zum Parameter |
| `internal/e2e/*_docker_test.go` | unverändert (Hebel 1 und 2 sind Out-of-Scope) | Gruppenzuordnung nach Weg B liest die Testnamen, ändert sie nicht; Namen der Tests nicht neu benannt |
| `docs/plan/planning/observations/BEO-E2E/…` | Evidence, Closure | §8 |
| ADR (Folge von [`ADR-0004`](../../adr/0004-ci-system.md)) | neu, nur wenn der Architect es bestätigt | Vorbedingung, eigener Slice |

Nicht gelesen und deshalb nicht als „unverändert“ geführt: `Dockerfile`, `ci.yml` über die Job-Namen hinaus, `.harness/skills/*.md`, `harness/roles.md`, `harness/review.md`, der Verifier-Auftrag unter `.claude/`.

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat Frage 1 (Schnitt), Frage 2 (Parameter), Frage 4 (Required Check), Frage 5 (Runner-Minuten), Frage 6 (Zeitlimit) und Frage 7 (Reihenfolge der Hebel) beantwortet; die ADR-Frage ist dem Architect gestellt, eine Folge-ADR ist, falls nötig, `Accepted`. WIP-Limit 1 pro Lauf. **Hinweis zur Umsetzung (Projektinhaber, 2026-11-01), kein Zeitplan:** Der Slice soll nach dem Release über einen Devcontainer mit `--instance` abgearbeitet werden. Der Docker-Test im verschachtelten Daemon des Devcontainers ist kein zugesagter Nachweis ([`SPEC-031`](../../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)); der Beleg der Aufteilung ist der CI-Lauf, nicht der Lauf im Devcontainer.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): die Gruppenzuordnung nach Weg B und die Vollständigkeitsprüfung wachsen zu einem dritten und vierten Liefer-Punkt — dann Weg A als erster Slice, Weg B als eigener.
- `in-progress` → `open` (blockiert): ein Teil-Job läuft auf dem Runner nicht grün aus einem Grund, der in der Aufteilung liegt (etwa geteilter Zustand zwischen Tests), oder die Required-Check-Einrichtung im GitHub-UI ist nicht möglich.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Report liegt vor, ein CI-Lauf `integration-docker` mit allen Teilen ist grün (Lauf-Kennung in der Closure-Notiz, die Gesamtzeit gelesen und gegen die Zeiten aus §3 gestellt; ein roter Lauf nur mit dokumentiertem Carveout, nie als stilles Rot), Closure-Notiz mit Lerneintrag in §7.

## 6. Risiken und offene Punkte

- **Ein Test fällt aus allen Teilen heraus** (neue Testfamilie ohne Gruppe bei Weg B, neues Paket bei Weg A): die Aufteilung würde dann Tests still nicht fahren, eine Lockerung des Gates ([`AGENTS.md`](../../../../AGENTS.md) §3.6). — **Ausgang:** offen; die Vollständigkeitsprüfung der DoD trägt es.
- **Die Ersparnis ist kleiner als erwartet:** Weg A senkt die Zeit höchstens auf die des e2e-Teils (siehe Folgerung §3); jeder Teil zahlt Aufbau und Übersetzen erneut. — **Ausgang:** offen; die Messung im Closure-Lauf entscheidet, ob Weg B folgt.
- **Der Sammel-Job oder die Namen der Matrix-Jobs verhalten sich im Branch-Schutz anders als angenommen** (Verhalten von GitHub bei übersprungenen oder umbenannten Jobs, nicht gelesen und nicht lokal prüfbar). — **Ausgang:** offen; Prüfung im ersten Lauf auf einem Branch.
- **Runner-Minuten** steigen (Frage 5). — **Ausgang:** offen bis zur Antwort des Projektinhabers.
- **Das Zeitlimit** gilt je Teil neu; der Zusammenhang mit [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (kein Ausgang zugewiesen) ändert sich, wenn der längste Teil kürzer wird. — **Ausgang:** weiter offen → diese Beobachtung (Evidence beim Closure).
- **Übersprungene Tests** (vier Nested-Podman-Tests und `StatusFiles_VFS` auf dem Runner) bleiben übersprungen; die Aufteilung ändert das nicht, und ein grüner Teil belegt sie nicht. — **Ausgang:** weiter offen → [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md).
- **Der Umsetzungsweg im Devcontainer** (Hinweis in §4) prüft die CI-Logik nicht: Workflow und Required Check sind dort nicht lauffähig; der Beleg ist der CI-Lauf auf dem Push. — **Ausgang:** offen bis zur Umsetzung.
- **Annahmen der Anfrage, nicht gelesen:** feste Containernamen (Frage 3), „lokal ca. 26 Minuten“, Runner-Minuten-Grenze (Frage 5). — **Ausgang:** offen; der Implementer liest oder misst vor dem Start.
- **Das lokale Log der Zeitwerte ist ein roter Lauf** (drei Tests rot, §3); die Verteilung der Zeiten ist damit für die Dauer der gelaufenen Tests belastbar, für ein grünes Gesamtbild nicht. — **Ausgang:** offen; der CI-Lauf 41285105546 ist die grüne Gegenprobe der Paketzeiten.

## 7. Closure-Notiz

Noch nicht geschlossen. Entwurf; wird mit der Closure gefüllt.

- **Was hat funktioniert:**
- **Was ging anders als geplant:**
- **Steering-Loop-Eintrag:**
- **Beobachtungs-Register (`../observations/`):**
- **Folge-Slices:**
- **Risiken aus §6:**
- **Drei Paarungen:**

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `.github/workflows/` (keine eigene Sub-Area-Zeile in `harness/conventions.md`; ausführbarer Vertrag, Source-Precedence Rang 6), `repo.mk` und `scripts/` (WZG), `harness/` und `docs/maintainer/`, `docs/plan/` (DOK, `docs/maintainer/` ohne eigene Zeile, dort als Doku-Pfad geführt, nicht neu geschnitten). Alle Zeilen der Modus-Deklaration sind Greenfield; keine zu grobe Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-E2E` durchgegangen (Verzeichnisliste; die vier unten genannten Zustandsdateien gelesen, die `observation.md` nicht):

- [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) — Treffer: Laufzeit des Pakets und Limit; Ausgang liegt beim Projektinhaber (Risiko 5). Zähler: aus den Dateien unter `evidence/` zu lesen (Implementer).
- [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) — Treffer: Skips bleiben (Risiko 6).
- [`BEO-E2E/test-zaehlt-container-paralleler-pakete`](../observations/BEO-E2E/test-zaehlt-container-paralleler-pakete/observation.md) — Treffer: Parallelität mehrerer Pakete gegen eine Engine; die Aufteilung trennt Pakete auf getrennte Engines (Weg A), innerhalb eines Teils bleibt sie.
- [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) — Treffer: Aufräumen auf einer geteilten Engine; auf getrennten Runnern entfällt die Wirkung zwischen den Teilen. Ob ein Test dort Reste hinterlässt, ist nicht gelesen.
- [`BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig`](../observations/BEO-E2E/test-docker-im-verschachtelten-daemon-nicht-gleichwertig/observation.md) — Treffer: der Umsetzungsweg im Devcontainer (§4).

Evidenz-/Diskrepanz-Risiko: mittel — die Zeitwerte sind aus zwei Läufen gelesen, davon einer rot (§3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

### Sub-Area: `repo.mk` / `scripts/` / Workflow (WZG)

- **Modus:** GF
- **Konventionen-Dichte:** `modus-pruefung` pinnt die Bedingungen der Modi von `BUILD_MODE` ([`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md)); Teillauf ist kein Nachweis.
- **Phase-Reife:** Phase 4 (Gate besteht, wird umgebaut).
- **Evidenz-/Diskrepanz-Risiko:** mittel (siehe oben).
- **Reconciliation-Aufwand:** gering; Folge-Slice: keiner angelegt.

*Alle übrigen berührten Sub-Areas (`harness/`, `docs/plan/`, `docs/maintainer/`): GF.*
