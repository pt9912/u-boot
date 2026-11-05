# Slice CR: `devcontainer remove` — Lastenheft-Änderung zu `LH-FA-DEV-010`: ein Befehl entfernt den Container des Projekts, mit `--volumes` auch seine Volumes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Zweiter von vier Slices; ihm geht [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](../open/slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes: nur Anwendung, keine Änderung), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) (Volumes je Instanz).
Anlass: Projektinhaber, 2026-11-05 (`devcontainer down` entfernt den Container nicht und seine Volumes nicht; Weg „eigener Befehl `devcontainer remove`“ gewählt).

**Berührte Spec-Stellen:** `lastenheft.md`: Abschnitt zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Einleitungssatz, Spiegelstrich zu `down`, neuer Spiegelstrich zu `remove`, Fehlerklassen, Nicht-Ziele), die Aufzählungen in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), die Befehlsübersicht [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt dieser Slice nicht. Der Wortlaut steht im CR-Vorschlag [`cr-vorschlag-devcontainer-remove`](../../cr/cr-vorschlag-devcontainer-remove.md) (Zustand: freigegeben wie vorgelegt, übernommen).

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft sagt zu, dass `u-boot devcontainer remove [--instance <name>] [--volumes]` den Container des Projekts entfernt und mit `--volumes` zusätzlich seine Volumes, mit Abfrage und ohne dass `devcontainer down` seine Bedeutung („stoppen, nichts entfernen“) ändert; der Wortlaut trägt nur, was die Messung belegt.

Ausgangslage (gelesen):

- `devcontainer down` stoppt und entfernt weder Container noch Volumes ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zu `down`; `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Down`: `FindContainer`, dann `launcher.Stop`).
- `up --recreate` entfernt den Container und lässt benannte Volumes stehen (Benutzerhandbuch, Abschnitt „Stoppen, Neuanlage und Volumes“); Volume-Namen tragen `${devcontainerId}` ([`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md)).
- Das Aufräumen heute: `docker rm` und `docker volume rm` von Hand (Benutzerhandbuch, Abschnitt zu Instanzen). Das Muster einer zerstörerischen Aktion mit Abfrage liefert `u-boot down --volumes` (Bestätigung, `--yes`, `--json` ohne `--yes` abgelehnt, Liste der entfernten Volumes).

**Entscheidungen des Projektinhabers** (getroffen 2026-11-05; Quelle: Antwort des Projektinhabers im Sitzungsdialog 2026-11-05; der Wortlaut dazu steht im CR-Vorschlag und ist vom Projektinhaber am 2026-11-05 im Sitzungsdialog „freigegeben wie vorgelegt“ worden)

1. **Name und Ort:** `u-boot devcontainer remove`. — **Ausgang:** entschieden 2026-11-05. `u-boot remove <service>` bleibt unverändert und bedeutet etwas anderes (Dienst aus der Konfiguration nehmen); der Wortlaut nennt den Unterschied, das Handbuch grenzt beide ab.
2. **Standard:** `remove` entfernt den Container; Volumes bleiben. `--volumes` entfernt zusätzlich die Volumes dieses Containers, mit Abfrage, `--yes` und der Liste der entfernten Volumes — wie `u-boot down --volumes`. — **Ausgang:** entschieden 2026-11-05.
3. **Laufender Container:** `remove` stoppt einen laufenden Container und entfernt ihn dann. — **Ausgang:** entschieden 2026-11-05 (Gegenposition „Abbruch mit Hinweis auf `down`“ verworfen; Messfrage M2: `docker rm` eines laufenden Containers scheitert ohne `-f`, die Messung trägt beide Wege).
4. **Volume der Ollama-Modelle:** `--volumes` entfernt es mit; die Abfrage nennt es. — **Ausgang:** entschieden 2026-11-05 (Gegenposition „eigenes Flag“ verworfen; Messfrage M4).
5. **`--instance <name>`:** wie bei `down`: ohne Flag wirkt der Befehl nur auf den Container ohne Instanz; gibt es nur Container von Instanzen, endet er mit Exit-Code `10` und dem Hinweis auf `--instance`. — **Ausgang:** entschieden 2026-11-05.

**Frage aus der Messung (M5), beantwortet:** Welcher Container ist bei `remove` ohne `--instance` gemeint? — **Ausgang:** „wie `down`“ (Antwort des Projektinhabers im Sitzungsdialog 2026-11-05, bestätigt). Die Messung [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) (M5) belegt für `down` und `up`: ohne `--instance` gilt nur der Container ohne Instanz; gibt es im Projekt nur Container von Instanzen (laufend, gestoppt oder gemischt), endet der Befehl mit Exit-Code `10` und dem Hinweis auf `--instance`. Für `remove` belegt sie nichts; die Antwort ist eine Anforderung, keine Aussage der Messung.

**Auslegungen, die E1 bis E5 offen ließen** — je mit Gegenposition in [`cr-vorschlag-devcontainer-remove`](../../cr/cr-vorschlag-devcontainer-remove.md); **Ausgang:** F1 bis F7 entschieden am 2026-11-05 (Antwort des Projektinhabers im Sitzungsdialog: „freigegeben wie vorgelegt“, die Vorschläge gelten):

- **F1** (Folge einer abgelehnten oder fehlenden Freigabe): Container und Volumes bleiben bestehen.
- **F2** (Aufruf ohne Container): Exit-Code `0`, nichts wird entfernt.
- **F3** (Volumes nach einem `remove` ohne `--volumes`): `u-boot` entfernt sie danach nicht mehr; der zweite Satz in den Nicht-Zielen steht im Lastenheft, das Handbuch nennt die Reihenfolge.
- **F4** (Zusage-Menge): Spiegelstrich zur Zusage-Menge bleibt unverändert; Plattform- und Messgrenzen von `remove` führt die Spezifikation.
- **F5** (Volume, das ein anderer Container hält): nicht im Lastenheft; die Spezifikation legt fest, was `remove --volumes` dann tut.
- **F6** (Priorität): `V2`.
- **F7** (Version): `0.37.0`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Spezifikation und Code** — Folge-Slices [`slice-v2-spec-devcontainer-remove`](../open/slice-v2-spec-devcontainer-remove.md) und [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md); die Spezifikation setzt den freigegebenen Vertrag voraus.
- **Änderung von `devcontainer down` und `up --recreate`** — beide behalten ihre Zusage; eine Änderung wäre eine Zusage-Verschärfung ohne Anlass.
- **`remove` für Compose-Dienste** — `u-boot down --volumes` deckt sie ab.

## 2. Definition of Done

- [x] **CR-Vorschlag und Freigabe:** `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` (Form `cr-vorschlag.template.md`) trägt je Stelle Alt/Neu; der Projektinhaber gibt den Wortlaut frei; die Entscheidungen 1 bis 5 sind getroffen (2026-11-05, §1); die Auslegungen F1 bis F7 des CR-Vorschlags sind beantwortet.
- [x] **Übernahme:** `spec/lastenheft.md` trägt den Wortlaut (die Stellen (a) bis (g) des CR-Vorschlags: Spiegelstrich zu `remove` samt Einleitungssatz, Fehlerklassen und Nicht-Ziele, Aufzählungen in [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), Befehlsübersicht, Version `0.37.0`, Historie-Zeile mit Datum); Abweichungen der Übernahme je Stelle sind im CR-Vorschlag ausgewiesen.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` | neu (Commit `7e6e6c9f`), update bei der Übernahme | Vorschlags-Datei nach `cr-vorschlag.template.md` (bleibt dauerhaft im CR-Ordner, Ablage-Regel); Stellen (a) bis (g), Alt/Neu je Stelle; Übernahme: Status „freigegeben wie vorgelegt, übernommen“, Abweichungen (keine), Ausgang F1 bis F7, Platzhalter in Stelle (g) aufgelöst |
| `spec/lastenheft.md` | update (nach der Freigabe, Lastenheft `0.37.0`) | Stellen laut CR-Vorschlag; gelesen am HEAD vor der Änderung: Versionsfeld Zeile 9 (`0.36.1`); [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeile 189 (Liste der ändernden Subkommandos unter `--quiet` endet auf `devcontainer shell`); [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeilen 210 und 214 (Aufzählungen „insb. … `u-boot down --volumes`“; Exit-Code `10` im nicht-interaktiven Modus ohne `--yes` für destruktive Operationen); [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile 512 (Einleitungssatz „starten, … öffnen und ihn stoppen“), 515 (Spiegelstrich `down`), 516 (Spiegelstrich `shell`), 517 (Volumes zweier Instanzen), 521 (Fehlerklassen), 522 (Zusage-Menge, unverändert), 523 (Nicht-Ziele „räumt keine auf“); [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) Zeilen 1186 bis 1188 (Befehlsübersicht); §16 Historie Zeile 1659 (letzte Zeile `0.36.1`) |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zustandsmarker „Wortlaut wird vorgelegt“ auf „Wortlaut freigegeben und im Lastenheft `0.37.0` übernommen, Closure steht aus“ (gelesen: `open/README.md` Zeile 24 und `roadmap.md` Zeile 45 nennen den Slice in `in-progress/`); bei der Closure wechseln beide Zeilen auf „geschlossen“ |
| `docs/plan/planning/in-progress/slice-cr-devcontainer-remove.md` | update | dieser Slice: Version `0.37.0` in der DoD, §1 (Entscheidungen mit Ausgang), §3, §6, §8, Zustandsmarker |
| `docs/plan/planning/observations/BEO-DOK/*/evidence/` | neu bei der Closure | Register-Evidence der Closure-Notiz (nicht Teil dieses Entwurfs) |

Gelesen, nicht berührt: [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) Zeilen 218 bis 236 (nur Anwendung); [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) Zeile 244 (Wert `devcontainer` im Feld `command` steht schon in der Aufzählung); [`LH-FA-UP-004`](../../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) Zeile 709 (Compose: `u-boot down --volumes`); `docs/maintainer/cli-json-contract.md` Zeile 1244 (Abschnitt 6.11) und Zeile 1357 (Matrix) tragen Kopien der Befehle und ziehen im Spezifikations-Slice nach, nicht in diesem Slice.

Der Abgleich der Tabelle gegen `git diff --stat HEAD` nach der Freigabe und der Übernahme ist Pflicht (`AGENTS.md` §3.10); der Abgleich gilt für den Stand nach der letzten Fix-Runde. **Abgleich nach der Übernahme (Stand vor dem Commit der Übernahme):** der Diff berührt `spec/lastenheft.md`, `docs/plan/cr/cr-vorschlag-devcontainer-remove.md`, `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` und diese Datei; jede Datei hat eine Zeile der Tabelle, jede Zeile außer der Register-Evidence der Closure hat einen Diff. Nach einer Fix-Runde ist der Abgleich zu wiederholen.

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn die Messung eine Entscheidung 2 bis 4 anders trägt als vorgeschlagen; → `open`, solange die Freigabe des Projektinhabers aussteht.

## 5. Closure-Trigger

DoD vollständig, Wortlaut im Lastenheft, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der Wortlaut sagt mehr zu, als die Messung trägt (Zusage weiter als Messung). — **Ausgang:** offen; der Reviewer prüft jeden Satz gegen das Protokoll.
- `remove` und `u-boot remove <dienst>` werden verwechselt. — **Ausgang:** offen; das Handbuch grenzt beide ab, der Wortlaut nennt den Unterschied (letzter Satz des Spiegelstrichs in Stelle (c) des CR-Vorschlags).
- Eine Stelle der Befehlsübersicht oder einer Aufzählung wird übersehen, das Lastenheft widerspricht sich selbst ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)). — **Ausgang:** offen; die Vorab-Suche (Befehlsnamen, `--volumes`, `--quiet`, Fehlerklassen, Nicht-Ziele, Befehlsübersicht) ergab die Stellen (a) bis (g); der Reviewer wiederholt sie.
- Der Nicht-Ziel-Satz „räumt keine [Instanz] auf“ widerspricht `remove --instance` und die Grenze „keine Volumes eines nicht mehr vorhandenen Containers“ überrascht den Benutzer (die Messung trägt die Bestimmung der Volumes nur am vorhandenen Container). — **Ausgang:** offen; Stelle (d2) und F3 des CR-Vorschlags sind freigegeben (Grenze steht im Lastenheft); der Reviewer prüft den Satz gegen das Protokoll.
- Der Wortlaut sagt `remove --volumes` ohne Plattformgrenze und ohne die Ausnahme eines von einem anderen Container gehaltenen Volumes zu, die Messung ist ein Lauf auf der verschachtelten Engine. — **Ausgang:** offen; Plattform- und Messgrenzen führt die Spezifikation (Anweisung des Projektinhabers); F4 und F5 sind wie vorgeschlagen freigegeben (kein Halbsatz im Lastenheft), der Reviewer prüft die Zusage gegen das Protokoll.
- Der Wortlaut ändert sich nach dem Review ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). — **Ausgang:** offen; die Übernahme hat den freigegebenen Wortlaut nicht geändert (Abweichungen je Stelle: keine, im CR-Vorschlag ausgewiesen), ein Nachlauf-Review liest den Diff jeder Fix-Runde.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/` durchgegangen am 2026-11-05 (alle Titel der Bereiche DOK, APP, CLI, WZG und E2E gelesen; die Zustandsdateien und Beobachtungen der folgenden Treffer gelesen); einschlägig:

- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (39 Evidence-Dateien, offen, kein Ausgang zugewiesen): der CR-Vorschlag enthält eine Tabelle „Prüfung des Wortlauts gegen die Messgrenzen“, Zeile für Zeile gegen die Abschnitte „Darf nicht tragen“ (M1 bis M5); der Wortlaut nennt keine Plattform-, Namens- und Messgrenze (Risiko 1).
- [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (18 Evidence-Dateien, offen, Schwelle erreicht, Ausgang beim Projektinhaber): die Übernahme ändert den freigegebenen Wortlaut nicht; Abweichungen stehen im CR-Vorschlag je Stelle; Nachlauf-Review über den Diff jeder Fix-Runde (Risiko 5).
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (43 Evidence-Dateien, verkörpert in `AGENTS.md` §3.10): die Bedingungen des Wortlauts („nur auf den Container ohne Instanz“, „ohne `--yes` Exit-Code `10`“, „ohne Container Exit-Code `0`“) bekommen im Umsetzungs-Slice [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md) je einen pinnenden Test.
- [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (5 Evidence-Dateien, offen, Schwelle erreicht): die Vorab-Suche lief über Befehlsnamen (`devcontainer up|down|shell`), `--volumes`, die `--quiet`-Liste, die Fehlerklassen, die Nicht-Ziele und die Befehlsübersicht und fand Stellen außerhalb von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (zwei Stellen) und [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle); die Kopien in `docs/maintainer/cli-json-contract.md` sind Gegenstand des Spezifikations-Slice.
- [`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`](../observations/BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf/observation.md) (1 Evidence-Datei, offen): der CR-Vorschlag zeigt den Wortlaut je Stelle im vollen Text, nicht verdichtet; §7 überträgt ihn zeichengleich.
- [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (5 Evidence-Dateien, offen, Schwelle erreicht): Nutzerpfade nach dem Befehl geprüft — wiederholter Aufruf ohne Container (F2), `remove` ohne `--volumes` und die danach übrigen Volumes (F3), abgelehnte Abfrage (F1), Aufruf ohne `--instance` bei nur Instanz-Containern (Wortlaut).
- [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (10 Evidence-Dateien, offen): die Messung fuhr den ausgelieferten Befehl `u-boot devcontainer up`/`down`; `remove` gibt es nicht, der Wortlaut trägt für `remove` eine Anforderung, keine Messaussage (M5-Frage im Abschnitt „Darf nicht tragen“ des Protokolls).
- [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (3 Evidence-Dateien, offen): der Container wird nur über das Label des Projektordners (und das Instanz-Label) gefunden; für eine zerstörerische Aktion gibt der Spezifikations-Slice die Suche ausdrücklich an. Das Lastenheft schweigt zur Suche (Spezifikation).
- [`BEO-DOK/dod-wortlaut-vom-implementer-geaendert`](../observations/BEO-DOK/dod-wortlaut-vom-implementer-geaendert/observation.md) (2 Evidence-Dateien, offen): die Änderung am DoD-Wortlaut dieses Entwurfs macht der Planner, nicht ein Implementer, und er weist sie aus: DoD-Zeile 2 (Version `0.36.0` → `0.37.0`, weil `0.36.0` an [`LH-FA-DEV-013`](../../../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) vergeben ist) und DoD-Zeile 1 (Auslegungen F1 bis F7 beantwortet).

Die übrigen Titel betreffen diesen Slice nicht (gelesen: Titel).

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
