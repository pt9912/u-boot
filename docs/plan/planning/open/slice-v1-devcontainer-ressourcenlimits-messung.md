# Slice V1: Grenzen für Arbeitsspeicher und CPU des Devcontainers — die offenen Messungen auf einer Engine, die Grenzen anwendet, und die Entscheidungen, die an ihnen hängen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Slice von [`slice-v1-devcontainer-ressourcenlimits`](../done/slice-v1-devcontainer-ressourcenlimits.md), das mit eingeschränktem Gegenstand geschlossen wurde (Entscheidung des Projektinhabers, 2026-11-01); Priorisierung und Release-Zuordnung beim Projektinhaber. Eine Closure-Bedingung, die von der DoD dieses Slice verschieden ist, gibt es nicht.

**Bezug:** [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage nur Linux-Host `amd64` mit lokaler Engine).

**Berührte Spec-Stellen:** [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Beleg-Tabelle, Abschnitt *Nicht zugesagt*), je nach Entscheidung Lastenheft-CR.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Messungen 1, 2, 4 und 5 aus §1 von [`slice-v1-devcontainer-ressourcenlimits`](../done/slice-v1-devcontainer-ressourcenlimits.md) sind auf einer Engine gefahren, die Grenzen anwendet (Linux-Host `amd64`, lokale, nicht verschachtelte Engine, zum Beispiel ein CI-Runner), mit mindestens zwei Läufen je Zelle; Messung 3 bleibt Abgrenzung (andere Hosts und Engines: nicht zugesagt). Danach entscheidet der Projektinhaber über Wortlaut der Zusage, Diagnose-Code und Werteformat.

**Herkunft.** Die verschachtelte Engine des Devcontainers dieses Repos kann keine Grenze anwenden (Exit `125` an der cgroup); die Messungen 1, 2, 4 und 5 sind dort nicht gefahren, Messung 6 nur für Meldungen vor dem Start; der CI-Lauf des Jobs `integration-docker` hat die beiden Docker-Tests einmal bestanden (Protokoll [`slice-v1-devcontainer-ressourcenlimits`](../../../messungen/slice-v1-devcontainer-ressourcenlimits/README.md)).

**Wiederverwendet:** das Skript [`tools/messung/devcontainer-ressourcenlimits/messen-engine.sh`](../../../../tools/messung/devcontainer-ressourcenlimits/messen-engine.sh) und das Protokoll-Verzeichnis; kein neuer Produktcode ohne Entscheidung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Weitere `docker run`-Optionen, Plattenquoten, Compose-Grenzen, GPU, Grenzen des Hosts oder der VM** — dieselben Ausschlüsse wie im Eltern-Slice (§1 dort); kein Bedarf über Arbeitsspeicher und CPU hinaus.
- **Docker Desktop, Colima, Podman, Linux-`arm64`, Windows** — nicht zugesagt ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); eine Probe dort wäre eine eigene Messung.
- **Eine Änderung des Werteformats ohne Freigabe** — das Format bleibt wie im Lastenheft 0.27.0, bis der Projektinhaber entscheidet.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung:** Messungen 1, 2, 4 und 5 sind auf einer Engine, die Grenzen anwendet, gefahren (Linux `amd64`, lokale Engine, mindestens zwei Läufe je Zelle, datiert, mit Umgebung; beide Profile, mit und ohne `nestedRuntime: docker`; ein zu kleines Limit; zwei Instanzen), Messung 6 für Annahme und Abweisung der Werte durch die Engine ergänzt; Protokoll im bestehenden Verzeichnis unter `docs/messungen/slice-v1-devcontainer-ressourcenlimits/`, ohne Zugangsdaten in den Rohdaten.
- [ ] **Entscheidung des Projektinhabers** zu drei Punkten, mit dem Befund der Messung als Grundlage: Wortlaut der Zusage-Menge (Wirkung auf den verschachtelten Docker, zu kleines Limit, mehrere Instanzen), eigener Diagnose-Code oder die Abweichung gilt, Freigabe des Werteformats.
- [ ] **Folgen der Entscheidung:** Lastenheft (CR, falls der Wortlaut sich ändert), [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs), Handbuch und `CHANGELOG.md` nennen nur das Gemessene; jede neue Bedingung hat einen pinnenden Test oder entfällt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die öffentliche Schnittstelle, soweit ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10). Gelesen ist nur das Messprotokoll und das Skript des Eltern-Slice; der Rest ist „prüfen“.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/messungen/slice-v1-devcontainer-ressourcenlimits/README.md`, neue Rohdaten (`roh-*.txt`) | update, neu | Ergebnis je Messung und Lauf (gelesen: Abschnitt *Messung an der Engine*, *Nicht gemessen*) |
| `tools/messung/devcontainer-ressourcenlimits/messen-engine.sh` | prüfen | misst Messungen 1, 2, 4 und 6 (gelesen: Kopfkommentar); Messung 5 (zwei Instanzen) braucht ein Skript oder eine Erweiterung |
| `spec/lastenheft.md`, `spec/spezifikation.md` ([`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)), `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update, nur nach Entscheidung | Zusage-Menge folgt dem Gemessenen |
| `internal/e2e/resources_docker_test.go` | prüfen | Docker-Test des Eltern-Slice; neue Zellen nur, wenn der Projektinhaber sie in die Zusage aufnimmt |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; eine Engine, die Grenzen anwendet, steht für die Läufe bereit (Linux-Host `amd64`, lokale Engine, oder ein CI-Job dafür); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): keine Engine verfügbar, die Grenzen anwendet.

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Zusage in Lastenheft, Spezifikation und Handbuch nennt nur das Gemessene; Review-Report und Verification-Report liegen vor; Closure-Notiz mit Steering-Loop-Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

1. **Die Grenze trägt im Sandbox-Profil oder im verschachtelten Docker nicht.** Dann schrumpft die Zusage (nur Standardprofil, oder nur der äußere Container) und der Projektinhaber entscheidet neu. — **Ausgang:** bei Closure.
2. **Auch die CI-Runner wenden Grenzen nicht an oder ändern sich.** Dann bleibt die Messung offen, der Slice geht zurück nach `open/`. — **Ausgang:** bei Closure.
3. **Ein zu kleines Limit bricht ein Startscript ab** (Exit-Code von `up`: `11` oder `12`, [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); die Meldung dazu ist nicht gemessen. — **Ausgang:** bei Closure.
4. **Die Zusage reicht weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): jede Zählung und jede Mengenaussage („beide Profile“, „zwei Instanzen“) wird Satz für Satz gegen das Protokoll gelesen. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK) sowie `tools/` (WZG) und gegebenenfalls `internal/e2e` (E2E); alle sind in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert.

**Vorgelagert — offene Beobachtungen sichten:** nach Titeln gesichtet: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 4) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md); Inhalt beim Start zu lesen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
