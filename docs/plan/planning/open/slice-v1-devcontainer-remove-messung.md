# Slice V1: `devcontainer remove` — Messung, welche Container-Objekte und Volumes ein Devcontainer hält und was `docker rm` und `docker volume rm` daran tun

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Erster von vier Slices. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope: Grundlage der Änderung), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) (Volume-Namen je Instanz).

**Berührte Spec-Stellen:** — (Messung; kein Vertrag wird geändert).

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Messprotokoll belegt, welche Volumes ein Devcontainer-Container hält, woran `u-boot` sie erkennt und wie `docker rm` und `docker volume rm` sich bei laufendem, gestopptem und instanziertem Container verhalten, damit der Wortlaut des CR nur Gemessenes trägt.

Ausgangslage (gelesen):

- `devcontainer down` stoppt und entfernt weder Container noch Volumes ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zu `down`; `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Down`: `FindContainer`, dann `launcher.Stop`).
- `up --recreate` entfernt den Container und lässt benannte Volumes stehen (Benutzerhandbuch, Abschnitt „Stoppen, Neuanlage und Volumes“); Volume-Namen tragen `${devcontainerId}` ([`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md)).
- Das Aufräumen heute: `docker rm` und `docker volume rm` von Hand (Benutzerhandbuch, Abschnitt zu Instanzen). Das Muster einer zerstörerischen Aktion mit Abfrage liefert `u-boot down --volumes` (Bestätigung, `--yes`, `--json` ohne `--yes` abgelehnt, Liste der entfernten Volumes).

Messfragen:

- **M1** Welche Volumes hält der Container (benannte und anonyme; Standard- und Sandbox-Profil, mit und ohne Features `ollama` und Podman-Runtime), und trägt `docker inspect` sie vollständig in `Mounts`? Folgt daraus eine Erkennung über den Container statt über die Auflösung von `${devcontainerId}`?
- **M2** Wie verhält sich `docker rm` bei einem laufenden Container (ohne und mit `-f`), und bleibt ein Volume, das ein zweiter Container hält, beim `docker volume rm` stehen?
- **M3** Trägt die Instanz (`--instance`) eigene Volumes (Namen, Isolation) und löscht das Entfernen einer Instanz nur deren Volumes?
- **M4** Welche Volumes überleben das Entfernen des Containers absichtlich (das Volume der Ollama-Modelle, laut Lastenheft), und wie sind sie zu unterscheiden?

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code für `remove`** — erst der Wortlaut (CR) und die Spezifikation legen den Vertrag fest; die Messung liefert nur die Grundlage.
- **Entscheidung, welche Volumes `--volumes` löscht** — gehört dem Projektinhaber im CR (Folge-Slice [`slice-cr-devcontainer-remove`](slice-cr-devcontainer-remove.md)).
- **Messung auf Hosts außer Linux-`amd64` mit lokaler Engine** — die Zusagen des Devcontainer-Bereichs gelten nur dort; weitere Hosts wären ein eigener Vorgang.

## 2. Definition of Done

- [ ] **Messprotokoll:** `docs/messungen/slice-v1-devcontainer-remove-messung/README.md` beantwortet M1 bis M4 je mit Befehl, Ausgabe, Datum und Grenzen (kein Wert ohne Messung); das Skript liegt unter `tools/messung/slice-v1-devcontainer-remove-messung/` (Ablage-Regel in `harness/conventions.md`).
- [ ] **Eingang für den CR:** das Protokoll nennt je Frage den Satz, den der CR-Wortlaut tragen darf, und die Sätze, die er nicht tragen darf.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/slice-v1-devcontainer-remove-messung/` | neu | Messskript und Hilfsdateien |
| `docs/messungen/slice-v1-devcontainer-remove-messung/README.md` | neu | Protokoll M1 bis M4 |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

Der Implementer ergänzt die Tabelle und gleicht sie vor dem Handoff gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber priorisiert; `in-progress/` ist leer (WIP-Limit 1).

**Rückführungen:** `in-progress` → `next`, wenn eine Messfrage ein Gerät oder eine Engine braucht, die der Devcontainer nicht stellt; → `open`, wenn die Engine im Devcontainer-Container nicht zur Verfügung steht.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Messung im verschachtelten Docker des Devcontainer-Containers weicht von einer lokalen Engine ab. — **Ausgang:** offen bis Messung; das Protokoll nennt den Aufbau.
- Der Messaufbau verbraucht Image-Bauzeit und Plattenplatz. — **Ausgang:** offen; das Protokoll nennt beides.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `tools/`, `scripts/`, `docs/messungen/` (WZG) und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
