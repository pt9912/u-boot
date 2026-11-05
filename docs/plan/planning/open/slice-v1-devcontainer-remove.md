# Slice V1: `devcontainer remove` — der Befehl: Container entfernen, mit `--volumes` auch seine Volumes (Abfrage, `--yes`, `--instance`, `--json`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle, Release-Zuordnung offen. Letzter von vier Slices; ihm gehen [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) und davor der CR voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](../next/slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (und die Stelle, die der Spec-Slice anlegt); `architecture.md` — keine.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot devcontainer remove [--instance <name>] [--volumes] [--yes]` entfernt den Container des Projekts und mit `--volumes` seine Volumes, nach dem Wortlaut des Lastenhefts und der Spezifikation.

Ausgangslage (gelesen): Der Befehl `devcontainer down` liegt im CLI-Adapter (`internal/adapter/driving/cli/devcontainer.go`: `newDevcontainerDownCommand`, `runDevcontainerDown`, `mapDevcontainerDownErrorToDiagnostic`) und im Dienst `DevcontainerStartService.Down` (`internal/hexagon/application/devcontainer_start.go`), der den Container mit `inspector.FindContainer` sucht und mit `launcher.Stop` stoppt. Die Bestätigung bei zerstörerischen Aktionen liefert der Port in `internal/hexagon/port/driven/confirm.go`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung von `devcontainer down`** — behält seine Zusage (stoppen, nichts entfernen).
- **Vertrag und Spezifikation** — liegen in den zwei vorangehenden Slices; dieser Slice setzt sie um und ändert sie nicht.
- **Entfernen von Images und des Hilfs-Images** — nicht Gegenstand des Vertrags; es wäre ein eigener Vorgang mit eigener Abfrage.

## 2. Definition of Done

- [ ] **Befehl und Ports:** `remove` (Container, `--volumes`, `--instance`, `--yes`, `--json`) arbeitet nach dem Vertrag; die Ports in `internal/hexagon/port/` bekommen die nötigen Operationen (Entfernen von Container und Volumes, Auflistung der Volumes eines Containers), der Docker-Adapter setzt sie um; Exit-Codes und Diagnostic-Abbildung stehen in der zentralen Klassifikation **und** in der Abbildung des neuen Subkommandos (Dual-Classifier-Regel, `.harness/skills/reviewer.md`). Tests nach [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Happy, Boundary, Negative, ohne Container, laufender Container, Instanz, Abfrage abgelehnt).
- [ ] **Docker-Test:** ein Test mit Build-Tag `docker` belegt gegen eine echte Engine, dass `remove --volumes` genau die Volumes des Containers entfernt und die eines anderen Containers stehen lässt.
- [ ] **Doku:** Benutzerhandbuch (Abgrenzung zu `u-boot remove <dienst>` und `devcontainer down`, Beispiel, Warnung „löscht Daten“; ohne Kennungen), CHANGELOG, Hilfetext.

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapter/driving/cli/devcontainer.go` | update | Subkommando `remove`; `mapDevcontainerRemoveErrorToDiagnostic`; zentrale Klassifikation |
| `internal/hexagon/application/devcontainer_start.go` | update | `Remove` neben `Down` |
| `internal/hexagon/port/driving/`, `internal/hexagon/port/driven/` | update | Request/Response und Port-Operationen (die Dateien nennt der Implementer nach Lesen) |
| `internal/adapter/driven/docker/` | update | Entfernen von Container und Volumes, Auflistung der Mounts |
| Tests dazu und `internal/e2e` (Build-Tag `docker`) | update / neu | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Doku-Pflicht (`AGENTS.md` §5) |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

Der Implementer liest die Dateien, ergänzt die Tabelle und gleicht sie vor dem Handoff gegen `git diff --stat` ab (`AGENTS.md` §3.10).

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn die Port-Änderung mehr als einen Slice trägt; → `open`, wenn die Messung eine Operation nicht trägt, die der Vertrag verlangt.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, `make test-docker` gelaufen oder begründet nicht, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- `--volumes` löscht Daten, die der Nutzer nicht erwartet (Ollama-Modelle, Workspace im Sandbox-Profil). — **Ausgang:** offen; die Abfrage nennt jedes Volume, `--json` ohne `--yes` wird abgelehnt.
- Die Volume-Erkennung über die Mounts erfasst ein Volume, das ein anderer Container teilt. — **Ausgang:** offen bis Messung M2.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driving/cli` (CLI), `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT), `internal/adapter/driven` (ADP), `internal/e2e` (E2E) und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
