# Verification-Report: slice-v1-test-host-umgebung-isolieren — 2026-11-04

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-test-host-umgebung-isolieren`](../plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md) — Code-Commit `4995f65e`, Review-Fix-Runde `e3547e47`; Abgleichstand **nach der letzten Fix-Runde = `e3547e47`** (= Repo-Tip, Arbeitsbaum sauber). Der Backlog-Anlage-Commit `09d765f0` (Slice-Datei, Beobachtung) und die vier Move-Commits (`e9ceb890`, `60d3139e`, `1b7275ac`, `9d3f3eca`) sind Lifecycle-Bewegungen.

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel), §2 (DoD), §3 (Plan-Tabelle), §6 (zwei offene Risiken), §7 (Closure offen)
- [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (der Wert steht in keiner Ausgabe), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Host-Auflösung)
- Review-Report [`2026-11-04-slice-v1-test-host-umgebung-isolieren.md`](2026-11-04-slice-v1-test-host-umgebung-isolieren.md) (0 HIGH / 0 MEDIUM; F-1, F-4 als offene Risiken)
- [`harness/verification.md`](../../harness/verification.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.9, §3.10), [`harness/roles.md`](../../harness/roles.md)

---

## Scope

- Slice: [`slice-v1-test-host-umgebung-isolieren`](../plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md)
- IDs: [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); keine Spec-, ADR- oder Produktcode-Änderung
- Artefakte: `internal/adapter/driven/git/hostsource_test.go`, `internal/e2e/devcontainer_shell_docker_test.go` (Build-Tag `docker`)

## DoD-Abgleich

DoD-Quelle: Plan §2. Gate-Lauf und Closure-Pflichten sind im Plan als „nicht gezählt“ markiert; hier trotzdem je belegt.

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Isolation** | **erbracht** | `isolateHostEnv` (`hostsource_test.go`, Z. 32 ff.) setzt `GIT_TOKEN` auf `""`, `GIT_CONFIG_GLOBAL` auf eine leere Datei im `t.TempDir()` und `GIT_CONFIG_NOSYSTEM=1` (je `t.Setenv`). Aufrufer im Diff: genau sechs — `GitToken_FallsBackToGh`, `…GhMissingReadsAsNoSource`, `…GhFailingReadsAsNoSource`, `…GhEmptyOutputReadsAsNoSource`, `ConfigValue_ReadsTheEffectiveValue`, `…MissingKeyIsEmptyWithoutError`. `…HostVariableWins` setzt `GIT_TOKEN` selbst auf einen Platzhalter und braucht keine Isolation. |
| 2 | **Kein Wert in Meldungen** | **erbracht** | Drei Meldungen in `hostsource_test.go` (`value set=%t source=%q …`, Wert nur als `value != ""`) und die Guard-Meldung in `devcontainer_shell_docker_test.go` Z. 408 (`output set=%t, err=%v`). Suche `GIT_TOKEN=` in der e2e-Datei: 0 Treffer. Der Diff trägt keinen weiteren `%q`/`%s` auf einen Umgebungswert. |
| 3 | `make gates` grün im Devcontainer-Container | **erbracht** (eigener Lauf) | Siehe Sensors. |
| 4 | Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge, Drei-Paarungen | **teils erbracht** | Review: Report vorliegend. Verifikation: dieser Report. Offen sind die Closure-Pflichten selbst (Closure-Notiz, Register, Ausgänge, Paarungen). |

## Plan-vs-Code (§3 gegen `git diff --stat 09d765f0..HEAD`, Stand nach `e3547e47`)

| Plan-Zeile (§3) | Diff | Abgleich |
|---|---|---|
| `internal/adapter/driven/git/hostsource_test.go` | +29/−… (Helfer, sechs Aufrufe, drei Meldungen; verwaister `scriptBinary`-Kommentar entfernt) | deckt |
| `internal/e2e/devcontainer_shell_docker_test.go` | eine Zeile (Z. 408) | deckt |
| `…/observations/BEO-ADP/test-liest-host-umgebung/` | nicht im Range (angelegt mit `09d765f0`, wie die Zeile selbst sagt) | deckt |
| `…/in-progress/roadmap.md` | 4 Zeilen | deckt |
| `docs/reviews/2026-11-04-…md` | neu | deckt |

Diff-Dateien ohne Plan-Zeile (Lifecycle-Mechanik, keine Produktdatei): die Slice-Datei selbst (`in-progress/`, Rest-Stub unter `open/`/`next/` je nach Range) und `docs/plan/planning/open/README.md` (Index-Zeile durch den Move). Hinweis zur Sauberkeit, keine DoD-Verletzung. Plan-Zeilen ohne Diff: keine. Produktcode (`hostsource.go`) im Diff abwesend, wie §1 zusagt.

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer, eigener Lauf am 2026-11-04) | **pass** | EXIT=0; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: 1100 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — coverage 92.20% meets threshold 90%`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `tool-pruefung … OK`. Stempel `.harness/state/gates-passed.head` = `e3547e47…` = HEAD. Das Log wurde in einer Datei gefangen; die Umgebung wurde nicht ausgegeben. |
| `make test-docker` | **not run, begründet** | Der Daemon ist erreichbar, aber der Lauf der gesamten `//go:build docker`-Suite braucht laut CI-Referenz (Vorgänger-Report: `integration-docker` 16m35s) mehr als die Zeitgrenze von ~10 Minuten dieses Auftrags. Der geänderte Test ändert nur den Text einer Fehlermeldung im Fehlerpfad (`t.Errorf`-Argumente); der Bedingungspfad (`err == nil \|\| out != ""`) ist unverändert. Die Kompilierung des Tags ist durch `gates` nicht belegt. **Das Risiko aus Plan §6 („geänderter Docker-Test lief nicht“) bleibt offen** → CI-Lauf `integration-docker` auf dem Push-Tip oder `make test-docker` vor der Closure. |
| `make verify-depguard` | **not run, begründet** | Keine Importregel berührt; der Diff ändert nur Testcode ohne neue Importe. |
| `make ci` | **not run** | Nicht angefordert; `govulncheck`/`image-scan` vom Diff nicht berührt. |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | Beide Commit-Messages nennen die ID; Meldungen ohne Wert (DoD 2) |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Plan-Bezug; Host-Auflösung von Token und Identität unter Isolation getestet (DoD 1) |
| Slice/Beobachtung | Beide Code-/Fix-Commits nennen den Slice und `BEO-ADP/test-liest-host-umgebung` |

Docs-Links: `make docs-check` im `gates`-Lauf grün (0 Befunde).

## Carveouts

- Neu: **keine** (Carveout-Ablage im Diff unberührt); keine Lockerung von Gates, Schwellen oder Architekturregeln.

## Rest-Risiken / offene Punkte vor der Closure

1. **Geänderter Docker-Test ungelaufen** (s. Sensors) — Beleg vor der Closure nachzuholen.
2. **F-1:** nicht gepinnt bleiben `GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG_COUNT`/`KEY_n`/`VALUE_n`, `GIT_CONFIG_PARAMETERS`; Läufe aus einem Git-Hook können `ConfigValue` anders beantworten lassen. Offen im Register.
3. **F-4:** `devcontainer_launcher_test.go` gibt bei Fehlschlag die `GIT_TOKEN=…`-Zeile des Prozess-Environments aus (der Wert-Ausschluss aus [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) ist dort nicht pinnend); ausdrücklich außerhalb dieses Slice, im Register geführt.
4. **Plan-Text:** Slice-§7 sagt „der Slice liegt in `open/`“, er liegt in `in-progress/`; der Plan-Kopf nennt den Autor noch als offenen Rollenwechsel. Beides wird mit der Closure-Notiz berichtigt (der Verifier ändert sie nicht).

## Gesamturteil

**Closure-fähig unter einer Bedingung.** Beide Liefer-DoD-Zeilen sind am Diff erbracht, `make gates` ist am Tip `e3547e47` grün, Plan-Tabelle und Diff decken sich, Traceability trägt. Vor der Closure muss der Lauf des geänderten Docker-Tests (CI `integration-docker` oder `make test-docker`) nachgewiesen werden; F-1/F-4 gehen als Register-Ausgänge in die Closure-Notiz.

## Commit / Artefakt

- Slice: `4995f65e`, `e3547e47`
- HEAD/Stempel: `e3547e471528d96f2f084e1183e0be9fd9564c02`
- Lokal: `make gates` EXIT=0 am 2026-11-04
