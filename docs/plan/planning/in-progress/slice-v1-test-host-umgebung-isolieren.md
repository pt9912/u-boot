# Slice V1: Adapter-Tests gegen die Umgebung des Hosts isolieren — `hostsource_test.go` pinnt `GIT_TOKEN` und die Git-Konfiguration, Fehlermeldungen geben keinen Token-Wert aus

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle (kleiner Testfix; keine Closure-Bedingung außerhalb der DoD).

**Anlass.** `make gates` war im Devcontainer rot (Projektinhaber-Sitzung 2026-11-04), in CI grün. Die Tests lasen `GIT_TOKEN` und die globale Git-Konfiguration des Hosts; drei Meldungen gaben den Token-Wert aus. Anweisung des Projektinhabers: der Token-Wert darf nicht ausgegeben werden.

**Bezug:** [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (der Wert steht in keiner Ausgabe), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Host-Auflösung von Token und Identität).

**Berührte Spec-Stellen:** `—` — nur Testcode.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Implementer, Umsetzung vor dem Plan; Rollenwechsel zum Reviewer offen). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Tests von `HostSource` antworten unabhängig vom Host: `GIT_TOKEN` leer, Git-Konfiguration ohne globale und Systemdatei; keine Fehlermeldung gibt einen Token-Wert aus.

**Ausdrücklich NICHT in diesem Slice:**

- **Kein Maskieren in Unit-Tests mit Fixture-Werten** (`devcontainer_start_test.go`, `devcontainer_hostgit_docker_test.go`): ihre Werte stammen aus Platzhaltern der Tests, nicht aus der Umgebung des Hosts.
- **Keine Änderung am Produktcode** (`hostsource.go`): das Lesen der Umgebung ist die Zusage des Adapters.

## 2. Definition of Done

- [x] **Isolation.** `isolateHostEnv` in `internal/adapter/driven/git/hostsource_test.go` setzt `GIT_TOKEN` leer, `GIT_CONFIG_GLOBAL` auf eine leere Datei und `GIT_CONFIG_NOSYSTEM=1`; sechs Tests rufen ihn auf.
- [x] **Kein Wert in Meldungen.** Drei Meldungen in `hostsource_test.go` und die Guard-Meldung in `internal/e2e/devcontainer_shell_docker_test.go` geben `set=<bool>` statt des Werts aus.
- [ ] `make gates` grün im Devcontainer-Container (lief grün; Nachweis im Handoff). *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge, Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

| Datei / Komponente | Art | Begründung / Fundstelle |
|---|---|---|
| `internal/adapter/driven/git/hostsource_test.go` | update | `isolateHostEnv`; Aufruf in `GitToken_FallsBackToGh`, `…GhMissingReadsAsNoSource`, `…GhFailingReadsAsNoSource`, `…GhEmptyOutputReadsAsNoSource`, `ConfigValue_ReadsTheEffectiveValue`, `…MissingKeyIsEmptyWithoutError`; drei Meldungen ohne Wert |
| `internal/e2e/devcontainer_shell_docker_test.go` | update | Guard-Meldung (Zeile ~408) ohne Ausgabe von `printenv GIT_TOKEN`; Docker-Test (Build-Tag `docker`), nicht ausgeführt |
| `docs/plan/planning/observations/BEO-ADP/test-liest-host-umgebung/` | neu | Beobachtung mit Evidence dieses Slice; angelegt mit der Backlog-Anlage `09d765f0`, nicht im Code-Commit |
| `docs/plan/planning/in-progress/roadmap.md` | update | der Ruhe-Marker „Nichts in Arbeit.“ weicht der Benennung dieses Slice (`planning-drift` des `docs-check`) |
| `docs/reviews/2026-11-04-slice-v1-test-host-umgebung-isolieren.md` | neu | Review-Report |

## 4. Trigger

**Start:** der Slice ist umgesetzt (Implementer-Lauf der Sitzung); es bleibt die Übergabe an den Reviewer. **Rückführung:** —.

## 5. Closure-Trigger

DoD vollständig; Review und Verifikation durch andere Rollen; Closure-Notiz; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

- **Der geänderte Docker-Test lief nicht.** `make gates` enthält `test-docker` nicht. — **Ausgang:** weiter offen → Lauf von `make test-docker` vor der Closure.
- **Weitere Tests lesen die Umgebung.** Der Review fand in `internal/adapter/driven/**` keine weitere Lücke, aber `devcontainer_launcher_test.go` gibt bei einem Fehlschlag die `GIT_TOKEN=…`-Zeile des Prozess-Environments aus (F-4). — **Ausgang:** weiter offen → `BEO-ADP/test-liest-host-umgebung` im Register.
- **Nicht gepinnt: `GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG_COUNT`/`KEY_n`/`VALUE_n`, `GIT_CONFIG_PARAMETERS`** (Review F-1); die Suite aus einem Git-Hook kann `ConfigValue` anders beantworten lassen. — **Ausgang:** weiter offen → `BEO-ADP/test-liest-host-umgebung` im Register.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driven` (ADP) und `internal/e2e` (E2E) sind deklarierte Sub-Areas, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-04 — Treffer nur verwandt: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (lokal grün, Runner rot — die Gegenrichtung); keine Zählung dort.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
