# Slice V1: `devcontainer remove` — Messung, welche Container-Objekte und Volumes ein Devcontainer hält und was `docker rm` und `docker volume rm` daran tun

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Erster von vier Slices. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](../open/slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](../open/slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md).

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
- **Entscheidung, welche Volumes `--volumes` löscht** — gehört dem Projektinhaber im CR (Folge-Slice [`slice-cr-devcontainer-remove`](../open/slice-cr-devcontainer-remove.md)).
- **Messung auf Hosts außer Linux-`amd64` mit lokaler Engine** — die Zusagen des Devcontainer-Bereichs gelten nur dort; weitere Hosts wären ein eigener Vorgang.

## 2. Definition of Done

- [x] **Messprotokoll:** `docs/messungen/slice-v1-devcontainer-remove-messung/README.md` beantwortet M1 bis M4 je mit Befehl, Ausgabe, Datum und Grenzen (kein Wert ohne Messung); das Skript liegt unter `tools/messung/slice-v1-devcontainer-remove-messung/` (Ablage-Regel in `harness/conventions.md`).
- [x] **Eingang für den CR:** das Protokoll nennt je Frage den Satz, den der CR-Wortlaut tragen darf, und die Sätze, die er nicht tragen darf.

Pflichten der Closure (zählen nicht mit):

- [x] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/slice-v1-devcontainer-remove-messung/` | neu | Messskript und Hilfsdateien |
| `docs/messungen/slice-v1-devcontainer-remove-messung/README.md` | neu | Protokoll M1 bis M4 |
| `docs/messungen/slice-v1-devcontainer-remove-messung/rohdaten-*.txt` | neu | Rohdaten der Läufe (archivierte Ausgabe der Skripte) |
| `docs/messungen/README.md` | update | Zeile im Bestand (Verzeichnis und Skriptpfad; Form der Zeilen oberhalb) |
| `docs/plan/planning/in-progress/slice-v1-devcontainer-remove-messung.md` | update | Plan-Tabelle (diese), DoD-Häkchen, Zustandsmarker in §6 |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | unverändert | Gelesen: `roadmap.md` Zeile 45 und `open/README.md` Zeile 23 nennen den Slice „in `in-progress/`“ und „in Umsetzung“; beides ist der Ist-Stand während der Messung und gilt bis zur Closure. Der Zustandsmarker wechselt mit der Closure, nicht mit der Messung. |

**Gelesen für die Messung (Fundstellen, `AGENTS.md` §3.10):**

- `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Down`: sucht den Container mit `inspector.FindContainer(ctx, BaseDir, Instance)`; ohne Fund leere Antwort; bei gestopptem Container Antwort mit `ContainerID` und ohne `Stopped`; sonst `launcher.Stop`. Kein Aufruf von `docker rm` oder `docker volume rm`.
- `internal/adapter/driven/docker/devcontainer_launcher.go`, `Launcher.Stop`: `docker stop <id>`; `runArgs`: `--remove-existing-container` bei `Recreate`, zwei `--id-label` (`devcontainer.local_folder`, `u-boot.instance`) nur mit Instanz.
- `internal/adapter/driven/docker/inspector.go`, `Inspector.FindContainer`: `docker ps -a --no-trunc --filter label=devcontainer.local_folder=<Ordner>` und mit Instanz zusätzlich `label=u-boot.instance=<Name>`; der Doc-Kommentar der Struktur nennt als einzige Befehle `docker ps`, `docker cp` und `docker context inspect`. Kein Code des Repos liest `docker inspect` oder die `Mounts` eines Containers (Suche nach `.Mounts` und `"inspect"` in `internal/`: nur `image inspect`/`context inspect` im Launcher und Inspector).
- `internal/hexagon/domain/containername.go`, `ComposeDevcontainerContainerName`: Name `u-boot-<projekt>` und mit Instanz `u-boot-<projekt>-<instanz>`.
- `internal/hexagon/application/devcontainer_sandbox.go`, Zeilen 154 bis 186 und `volumeMount`: das Standardprofil kehrt vor dem Zweig des Sandbox-Profils zurück und erzeugt kein `mounts`; das Sandbox-Profil trägt `<projekt>-workspace-${devcontainerId}`, mit `nestedRuntime: podman` zusätzlich `<projekt>-containers-${devcontainerId}`, mit `nestedRuntime: docker` `<projekt>-docker-${devcontainerId}`, mit dem Feature `ollama` `<projekt>-ollama-${devcontainerId}` (jeweils `type=volume`). Die Templates `devcontainer/Dockerfile.tmpl` und `devcontainer.json.tmpl` enthalten keine `VOLUME`-Anweisung; das Basis-Image trägt keine (`docker image inspect … .Config.Volumes` ist `null`, Messung unten).
- Nicht im Code gefunden und deshalb Messgegenstand: ob die Dev-Container-CLI selbst Volumes anlegt, ob `--remove-existing-container` anonyme Volumes mitnimmt, ob die Mounts eines Containers vollständig in `docker inspect` stehen.

Der Implementer gleicht die Tabelle vor dem Handoff gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber priorisiert; `in-progress/` ist leer (WIP-Limit 1).

**Rückführungen:** `in-progress` → `next`, wenn eine Messfrage ein Gerät oder eine Engine braucht, die der Devcontainer nicht stellt; → `open`, wenn die Engine im Devcontainer-Container nicht zur Verfügung steht.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Messung im verschachtelten Docker des Devcontainer-Containers weicht von einer lokalen Engine ab. — **Ausgang:** offen; das Protokoll nennt den Aufbau (Docker 29.9.3, `overlayfs`, cgroup v2, verschachtelt; `u-boot` behandelte die Engine als lokal). Die Abweichung gegen eine native Engine ist nicht gemessen und steht als Grenze im Protokoll.
- Der Messaufbau verbraucht Image-Bauzeit und Plattenplatz. — **Ausgang:** offen; das Protokoll nennt die Bauzeiten (4 bis 48 s je `up` mit Bau, je nach Zustand des Build-Caches) und die Image-Größen (1,41 bis 2 GB). Der Platzbedarf ist nicht belastbar gemessen, weil der Build-Cache der Engine die Schichten hält.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `tools/`, `scripts/`, `docs/messungen/` (WZG) und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
