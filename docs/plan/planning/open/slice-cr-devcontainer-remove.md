# Slice CR: `devcontainer remove` — Lastenheft-Änderung zu `LH-FA-DEV-010`: ein Befehl entfernt den Container des Projekts, mit `--volumes` auch seine Volumes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Zweiter von vier Slices; ihm geht [`slice-v1-devcontainer-remove-messung`](slice-v1-devcontainer-remove-messung.md) voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes: nur Anwendung, keine Änderung), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) (Volumes je Instanz).
Anlass: Projektinhaber, 2026-11-05 (`devcontainer down` entfernt den Container nicht und seine Volumes nicht; Weg „eigener Befehl `devcontainer remove`“ gewählt).

**Berührte Spec-Stellen:** `lastenheft.md`, Abschnitt zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Spiegelstrich zu `down`, neuer Spiegelstrich, Befehlsübersicht), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt dieser Slice nicht.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft sagt zu, dass `u-boot devcontainer remove [--instance <name>] [--volumes]` den Container des Projekts entfernt und mit `--volumes` zusätzlich seine Volumes, mit Abfrage und ohne dass `devcontainer down` seine Bedeutung („stoppen, nichts entfernen“) ändert; der Wortlaut trägt nur, was die Messung belegt.

Ausgangslage (gelesen):

- `devcontainer down` stoppt und entfernt weder Container noch Volumes ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zu `down`; `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Down`: `FindContainer`, dann `launcher.Stop`).
- `up --recreate` entfernt den Container und lässt benannte Volumes stehen (Benutzerhandbuch, Abschnitt „Stoppen, Neuanlage und Volumes“); Volume-Namen tragen `${devcontainerId}` ([`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md)).
- Das Aufräumen heute: `docker rm` und `docker volume rm` von Hand (Benutzerhandbuch, Abschnitt zu Instanzen). Das Muster einer zerstörerischen Aktion mit Abfrage liefert `u-boot down --volumes` (Bestätigung, `--yes`, `--json` ohne `--yes` abgelehnt, Liste der entfernten Volumes).

**Entscheidungen des Projektinhabers (Vorschlag; mit dem Wortlaut vorzulegen)**

1. **Name und Ort:** `u-boot devcontainer remove` (gewählt 2026-11-05). `u-boot remove <dienst>` bleibt unverändert und bedeutet etwas anderes (Dienst aus der Konfiguration nehmen); das Handbuch grenzt beide ab.
2. **Standard:** `remove` entfernt den Container; Volumes bleiben. `--volumes` entfernt zusätzlich die Volumes, die die Messung als Volumes dieses Containers belegt, mit Abfrage, `--yes` und der Liste der entfernten Volumes — wie `u-boot down --volumes`.
3. **Laufender Container:** Vorschlag: `remove` stoppt einen laufenden Container und entfernt ihn dann; Gegenposition: Abbruch mit Hinweis auf `down` (Messfrage M2).
4. **Volume der Ollama-Modelle:** Vorschlag: `--volumes` entfernt es mit, weil es zu diesem Container gehört; die Abfrage nennt es; Gegenposition: eigenes Flag (Messfrage M4).
5. **`--instance <name>`:** wie bei `down`: ohne Flag wirkt der Befehl nur auf den Container ohne Instanz; gibt es nur Container von Instanzen, endet er mit Exit-Code `10` und dem Hinweis auf `--instance`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Spezifikation und Code** — Folge-Slices [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) und [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md); die Spezifikation setzt den freigegebenen Vertrag voraus.
- **Änderung von `devcontainer down` und `up --recreate`** — beide behalten ihre Zusage; eine Änderung wäre eine Zusage-Verschärfung ohne Anlass.
- **`remove` für Compose-Dienste** — `u-boot down --volumes` deckt sie ab.

## 2. Definition of Done

- [ ] **CR-Vorschlag und Freigabe:** `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` (Form `cr-vorschlag.template.md`) trägt je Stelle Alt/Neu; der Projektinhaber gibt den Wortlaut frei; die Entscheidungen 1 bis 5 sind getroffen.
- [ ] **Übernahme:** `spec/lastenheft.md` trägt den Wortlaut (Spiegelstrich zu `remove`, Befehlsübersicht, Version `0.36.0`, Historie-Zeile mit Datum); Abweichungen der Übernahme je Stelle sind im CR-Vorschlag ausgewiesen.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` | neu | Vorschlags-Datei nach `cr-vorschlag.template.md` (bleibt dauerhaft im CR-Ordner, Ablage-Regel) |
| `spec/lastenheft.md` | update | Spiegelstrich zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Befehlsübersicht (Zeile `u-boot devcontainer down` im Abschnitt der Befehle), Versionsfeld, §16 Historie |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

Der Implementer liest `lastenheft.md` an den Stellen und ergänzt die Tabelle (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-v1-devcontainer-remove-messung`](slice-v1-devcontainer-remove-messung.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn die Messung eine Entscheidung 2 bis 4 anders trägt als vorgeschlagen; → `open`, solange die Freigabe des Projektinhabers aussteht.

## 5. Closure-Trigger

DoD vollständig, Wortlaut im Lastenheft, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der Wortlaut sagt mehr zu, als die Messung trägt (Zusage weiter als Messung). — **Ausgang:** offen; der Reviewer prüft jeden Satz gegen das Protokoll.
- `remove` und `u-boot remove <dienst>` werden verwechselt. — **Ausgang:** offen; das Handbuch grenzt beide ab, der Wortlaut nennt den Unterschied.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
