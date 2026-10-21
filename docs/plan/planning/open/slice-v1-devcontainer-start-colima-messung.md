# Slice v1: Einzelmessung des Befehls `devcontainer up` und `down` auf macOS mit Colima (nach dem Release)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Punkt von [`welle-devcontainer-start-ohne-editor`](../welle-devcontainer-start-ohne-editor.md) und von [`slice-v1-cli-devcontainer-up-down`](../in-progress/slice-v1-cli-devcontainer-up-down.md), eigener Slice nach dem Release (er hält die Closure der Welle und den Release nicht auf; es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre). Priorität V2.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zusage-Menge: Linux-Host `amd64` mit lokaler Engine; Colima „nicht zugesagt“); [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 (Zeile „macOS, Colima, `arm64`“ als Einzelmessung des Projektinhabers) und Freigabe-Punkt h; [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten).

**Berührte Spec-Stellen:** — (der Slice misst; ändert er Zusage und Lastenheft, ist das der Folge-CR, siehe §1).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der ausgelieferte Befehl `u-boot devcontainer up` und `down` ist auf macOS mit Colima (`arm64`) einmal gemessen, und das Protokoll entscheidet nach einer vorab festgelegten Regel, ob Colima in die Zusage aufgenommen wird.

**Herkunft.** Der Befehl-Slice [`slice-v1-cli-devcontainer-up-down`](../in-progress/slice-v1-cli-devcontainer-up-down.md) hatte die Einzelmessung in seiner DoD. Entscheidung des Projektinhabers (2026-10-21): „Die Colima-Messung wird nach dem Release nachgeholt.“ Die DoD-Zeile ging dort per Plan-Eintrag (§6 und §7 des Befehl-Slice) hierher; Colima bleibt in Lastenheft, Handbuch und `CHANGELOG.md` „nicht zugesagt“, das Release behauptet nur Linux/`amd64`.

**Vorhanden** (aus dem Befehl-Slice, nicht Teil dieses Slice): das Messskript `tools/messung/devcontainer-start-befehl/01-einzelmessung.sh`, das Protokoll-Gerüst `docs/messungen/slice-v1-cli-devcontainer-up-down/README.md` mit der Entscheidungsregel und ein Probelauf auf Linux (`rohdaten-00-probelauf-linux.txt`).

**Entscheidungsregel** (steht im Protokoll-Gerüst; hier die Fassung für den Slice): Die Messung *trägt*, wenn auf Colima (a) `up` im Standard- und im Sandbox-Profil mit Exit 0 endet und Fortgang wie im Probelauf zeigt (gleicher Container nach `down` und `up`, neuer Container nach `--recreate`, Marker im Workspace-Volume erhalten), (b) kein `docker.sock` im gestarteten Dev-Container steht und (c) ein Projekt außerhalb des geteilten Bereichs mit Exit 12 endet statt mit einem leeren Bind-Mount. Trägt sie, folgt ein Lastenheft-CR, der Colima in die Zusage aufnimmt; trägt sie nicht, bleibt Colima „nicht zugesagt“, und der Befund steht in §7 und, wenn er den Nutzer betrifft, im Handbuch.

**Spannung, die der CR auflöst.** [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 führt Colima als „zugesagt als Einzelmessung“; das Lastenheft (Rang 1) ist enger („Colima … nicht zugesagt“). Das Handbuch folgt dem Lastenheft. Die ADR ist `Accepted` und wird nicht umgeschrieben (`AGENTS.md` §3.5); der Folge-CR benennt die Abweichung und löst sie auf der Seite des Lastenhefts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Lastenheft-CR selbst** — er setzt ein tragendes Protokoll voraus und braucht die Freigabe des Wortlauts durch den Projektinhaber; bei einem tragenden Befund legt dieser Slice den CR-Slice in `open/` an (Kennung dann), schreibt aber keinen Lastenheft-Text.
- **Code des Befehls, Docker-Tests, Skriptänderungen** — der Befehl ist mit dem Befehl-Slice geliefert; zeigt die Messung einen Fehler im Befehl, ist das ein Befund an den Architect und ein eigener Fix-Slice, kein stiller Eingriff hier.
- **Docker Desktop, Podman, entfernte Engine, Windows, rootless Docker, SELinux, Linux-`arm64`, `amd64` unter macOS** — nicht gemessen und nicht zugesagt; eine Probe dafür wäre eine eigene Messung.
- **Eine Wiederholung auf mehreren Colima-Versionen** — Einzelmessung, keine Zahl wird zugesagt.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung:** `01-einzelmessung.sh` ist auf macOS mit Colima (Projektinhaber) gelaufen; die Rohdaten liegen als `docs/messungen/slice-v1-cli-devcontainer-up-down/rohdaten-01-colima-arm64.txt` ohne Zugangsdaten (Suche nach Zugangsdaten-Mustern ohne Treffer), und das Protokoll nennt Colima- und Docker-Version, Architektur und Projektpfad (im geteilten Home und, als Probe der Grenze, außerhalb).
- [ ] **Auswertung nach der Entscheidungsregel** (§1) im Protokoll `docs/messungen/slice-v1-cli-devcontainer-up-down/README.md`: *trägt* oder *trägt nicht*, je Punkt (a) bis (c) mit Beleg aus den Rohdaten.
- [ ] **Folgen:** trägt die Messung, liegt der Folge-CR-Slice in `open/` (Wortlaut-Entwurf zur Freigabe des Projektinhabers; die Abweichung zu [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 ist benannt); trägt sie nicht, bleiben Handbuch und `CHANGELOG.md` bei „nicht zugesagt“, und die Grenze, die die Messung zeigt, steht im Handbuch in Worten ohne Kennungen. Lastenheft, ADRs und Spezifikation sind von diesem Slice unverändert.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Autor nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/messungen/slice-v1-cli-devcontainer-up-down/rohdaten-01-colima-arm64.txt` | neu | Ausgabe des Skripts auf Colima; Aufruf steht im Protokoll-Gerüst (gelesen: Abschnitt *Skript und Aufruf*) |
| `docs/messungen/slice-v1-cli-devcontainer-up-down/README.md` | update | Abschnitt *Einzelmessung macOS mit Colima* (gelesen: „Nicht gelaufen“) wird die Auswertung nach der Entscheidungsregel; das Gerüst verweist bis dahin auf diesen Slice |
| `tools/messung/devcontainer-start-befehl/01-einzelmessung.sh` | **unverändert**, prüfen | läuft das Skript unter macOS nicht, ist die Änderung ein Plan-Eintrag; der Kopfkommentar nennt die Auflösung einer Sekunde für macOS (nicht gelesen auf weitere BSD-Abweichungen) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | **unverändert** bei einem tragenden Befund (Colima kommt erst mit dem CR in die Zusage); update nur, wenn die Messung eine Grenze zeigt, die der Nutzer kennen muss | Handbuch ohne Kennungen (`AGENTS.md` §5 Regel 6); jede Bedingung braucht einen Test oder entfällt (`AGENTS.md` §3.10) |
| Folge-CR-Slice in `open/` (Kennung bei der Anlage) | neu, nur bei einem tragenden Befund | Lastenheft-CR für Colima; Muster [`slice-cr-devcontainer-start-ohne-editor`](../done/slice-cr-devcontainer-start-ohne-editor.md) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile „In Arbeit“ bei Start, zurück auf „Nichts in Arbeit.“ bei der Closure |
| `spec/`, `docs/plan/adr/`, `Dockerfile`, Quelltext | **unverändert** | Slice misst; Lastenheft über den Folge-CR, ADR immutable |

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; **das Release ist geschnitten** (das nächste Minor nach `v0.12.0`; der Release-Cut-Slice liegt in `done/`); der Projektinhaber steht mit einem macOS-Host mit Colima bereit; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): kein macOS-Host mit Colima verfügbar, oder das Skript läuft unter macOS nicht und die Änderung braucht eine Entscheidung des Architect.
- `in-progress` → `next` (zu groß): die Messung verlangt eine Reihe über mehrere Hosts statt einer Einzelmessung.

## 5. Closure-Trigger

DoD vollständig; Rohdaten und Auswertung liegen im Protokoll, der Befund steht nach der Entscheidungsregel fest, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben. Ein *trägt nicht* ist ein vollständiger Ausgang des Slice, kein Fehlschlag.

## 6. Risiken und offene Punkte

- **Das Release geht ohne Colima-Messung hinaus.** — **Ausgang:** weiter offen; Handbuch, `CHANGELOG.md` und Lastenheft sagen „nicht zugesagt“, das Release behauptet nur Linux/`amd64`; das ist bis zur Messung die ehrliche Aussage (Verification des Befehl-Slice, V-1).
- **Das Skript ist nur unter Linux gelaufen** (Probelauf); BSD-Werkzeuge unter macOS können abweichen. — **Ausgang:** weiter offen; der Kopfkommentar des Skripts nennt die Auflösung einer Sekunde, mehr ist nicht gelesen.
- **Grenzen, die nur eine VM-Umgebung zeigt** (Pfade außerhalb des geteilten Bereichs, Spiegelpfad unter Colima, `arm64`-Verhalten der Dev-Container-CLI im Hilfs-Container, Socket-Pfad). — **Ausgang:** weiter offen; sie sind der Gegenstand der Messung.
- **`arm64`-Bau des Hilfs-Images:** Digest-Pins gelten für Index-Manifeste, ein `arm64`-Bau ist bisher nicht gelaufen (nur die Plattformliste gelesen, [`slice-v1-devcontainer-start-fundament`](../done/slice-v1-devcontainer-start-fundament.md) §6). — **Ausgang:** weiter offen; der Bau auf Colima ist der erste `arm64`-Lauf. Der Anker, den der Fundament-Slice auf den Befehl-Slice setzte, liegt damit hier.
- **Spannung zwischen [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 und dem Lastenheft.** — **Ausgang:** weiter offen; der Folge-CR löst sie auf, die ADR bleibt unverändert.
- **Das Protokoll trägt Zugangsdaten.** — **Ausgang:** weiter offen; das Skript setzt `GIT_TOKEN` auf einen Platzhalter und gibt keine Umgebung aus; Suche nach Mustern vor dem Commit ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/` (DOK, Messprotokoll) und gegebenenfalls `tools/` (WZG, Skript), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) beide Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** bei Start zu lesen. Einschlägig: [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (die Messung gilt dem ausgelieferten Befehl, nicht einer Variante), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (der Folge-CR sagt nur zu, was das Protokoll trägt). Evidenz-/Diskrepanz-Risiko: mittel (VM-Grenze, `arm64`).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
