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

**Autor:** pt9912 (Rolle Implementer, Umsetzung vor dem Plan). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Tests von `HostSource` antworten unabhängig vom Host: `GIT_TOKEN` leer, Git-Konfiguration ohne globale und Systemdatei; keine Fehlermeldung gibt einen Token-Wert aus.

**Ausdrücklich NICHT in diesem Slice:**

- **Kein Maskieren in Unit-Tests mit Fixture-Werten** (`devcontainer_start_test.go`, `devcontainer_hostgit_docker_test.go`): ihre Werte stammen aus Platzhaltern der Tests, nicht aus der Umgebung des Hosts.
- **Keine Änderung am Produktcode** (`hostsource.go`): das Lesen der Umgebung ist die Zusage des Adapters.

## 2. Definition of Done

- [x] **Isolation.** `isolateHostEnv` in `internal/adapter/driven/git/hostsource_test.go` setzt `GIT_TOKEN` leer, `GIT_CONFIG_GLOBAL` auf eine leere Datei und `GIT_CONFIG_NOSYSTEM=1`; sechs Tests rufen ihn auf.
- [x] **Kein Wert in Meldungen.** Drei Meldungen in `hostsource_test.go` und die Guard-Meldung in `internal/e2e/devcontainer_shell_docker_test.go` geben `set=<bool>` statt des Werts aus.
- [x] `make gates` grün im Devcontainer-Container (Commit `e3547e47`, Verifier-Lauf EXIT 0). *(nicht gezählt — Gate-Lauf)*
- [x] Review ([`2026-11-04-slice-v1-test-host-umgebung-isolieren`](../../../reviews/2026-11-04-slice-v1-test-host-umgebung-isolieren.md)), Verifikation ([`verification-2026-11-04-slice-v1-test-host-umgebung-isolieren`](../../../reviews/verification-2026-11-04-slice-v1-test-host-umgebung-isolieren.md)), Closure-Notiz, Register, Risiko-Ausgänge. *(nicht gezählt — Closure-Pflichten)*

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

- **Der geänderte Docker-Test lief nicht.** `make gates` enthält `test-docker` nicht. — **Ausgang:** entfallen: `make test-docker` lief am 2026-11-04 auf `e3547e47`; `TestDevcontainerShell_Docker_Session` ist `--- PASS` (16,73 s, Unterfall `guards_of_the_session` 1,68 s). Der Lauf endete mit Exit `2` wegen zwei fremder Fehlschläge, siehe nächstes Risiko.
- **`make test-docker` ist in diesem Devcontainer rot:** `TestE2E_Resources_DefaultProfileLimitsReachTheContainer` und `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` scheitern an `unable to apply cgroup configuration: cannot enter cgroupv2` — die verschachtelte Engine wendet keine cgroup-Grenze an (Herkunft im Plan von [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md)); die Tests berühren diesen Slice nicht. — **Ausgang:** weiter offen → [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md).
- **Weitere Tests lesen die Umgebung.** Der Review fand in `internal/adapter/driven/**` keine weitere Lücke, aber `devcontainer_launcher_test.go` gibt bei einem Fehlschlag die `GIT_TOKEN=…`-Zeile des Prozess-Environments aus (F-4). — **Ausgang:** weiter offen → `BEO-ADP/test-liest-host-umgebung` im Register.
- **Nicht gepinnt: `GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG_COUNT`/`KEY_n`/`VALUE_n`, `GIT_CONFIG_PARAMETERS`** (Review F-1); die Suite aus einem Git-Hook kann `ConfigValue` anders beantworten lassen. — **Ausgang:** weiter offen → `BEO-ADP/test-liest-host-umgebung` im Register.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Isolation (`isolateHostEnv`) macht `make gates` im Devcontainer mit gesetztem `GIT_TOKEN` und globaler Git-Identität grün (vorher rot); der Review fand 0 HIGH und 0 MEDIUM, die Verifikation bestätigte beide Liefer-Punkte.
- **Was ging anders als geplant:** Der Test-Fix entstand vor dem Slice-Plan (Rollenfolge verletzt, hier nachgeholt: Plan, `in-progress/`, Review, Verifikation). Beim Lauf des `docs-check` meldete `planning-drift` den Slice in `in-progress/` gegen den Ruhe-Marker der Roadmap; die Roadmap benennt den Slice jetzt.
- **Steering-Loop-Eintrag:** Ein Test, der die Umgebung des Hosts liest, ist nur in CI grün (Fund 1×, noch unter der Schwelle) — gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `BEO-ADP/test-liest-host-umgebung/` neu angelegt, Beleg `evidence/slice-v1-test-host-umgebung-isolieren.md` — Zähler steht bei 1×.
- **Folge-Slices:** keine; die Rest-Risiken F-1 und F-4 führt das Register.
- **Risiken aus §6:** je ein Ausgang, siehe §6.
- **Drei Paarungen:** Anker — kein Herkunfts-Anker gesetzt (nichts verkörpert); Folge-Slice — [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md) ist eine Datei in `open/`; Register — `BEO-ADP/test-liest-host-umgebung` liegt vor.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driven` (ADP) und `internal/e2e` (E2E) sind deklarierte Sub-Areas, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-04 — Treffer nur verwandt: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (lokal grün, Runner rot — die Gegenrichtung); keine Zählung dort.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
