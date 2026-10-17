# Verification: slice-v1-cli-egress-ohne-profil-sichtbar — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (Frage: Bauen wir es richtig? — gegen Plan, DoD, Spezifikation)

**Gegenstand:** Slice [`slice-v1-cli-egress-ohne-profil-sichtbar`](../plan/planning/done/slice-v1-cli-egress-ohne-profil-sichtbar.md), Commits `3063cb1..HEAD` (`d08146d`), Welle [`welle-devcontainer-sshd-claude-code`](../plan/planning/done/welle-devcontainer-sshd-claude-code.md), sechster Slice. Eingang: Review-Report [`2026-10-17-slice-v1-cli-egress-ohne-profil-sichtbar`](2026-10-17-slice-v1-cli-egress-ohne-profil-sichtbar.md) (0 HIGH/MEDIUM/LOW, 4 INFO). Reviewer-Findings sind keine Evidence; DoD, Spec und Sensors wurden separat geprüft.

## Verification Evidence

Scope:

- Slice: `slice-v1-cli-egress-ohne-profil-sichtbar` <!-- d-check:ignore (Slice-Kennung, Link oben) -->
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1, Lastenheft ungeändert), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten)
- Artefakte: `internal/hexagon/application/devcontainer_sandbox.go` (`sandboxRelaxationWarnings`), `internal/adapter/driving/cli/quiet_test.go`, `internal/hexagon/application/generate_sandbox_test.go`, `internal/hexagon/application/devcontainer_sandbox_relaxation_test.go`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`

## Plan gegen Diff

`git diff --stat 3063cb1..HEAD`: 10 Dateien, 96 Einfügungen, 27 Löschungen. Jede Datei steht in §3 des Slice; jeder Plan-Eintrag mit Art „update“ hat einen Diff. Die als „lesen“ geführten Zeilen (`export_test.go`, `harness/replay.md`, Golden-Fälle) haben keinen Diff, wie geplant. Die zehnte Datei ist der Slice selbst (3 Zeilen, Tabellen-Abgleich). Weder `spec/lastenheft.md` noch `spec/architecture.md` noch `harness/replay.md` sind im Diff (0 Treffer).

| Plan-Zeile (§3) | Diff | Befund |
|---|---|---|
| `devcontainer_sandbox.go` | `warn` zu `relax` nur im Zweig `egressEnabled(dc)` ohne Sandbox; Meldungstext unverändert; Doc-Kommentar nachgezogen | konform |
| `devcontainer_sandbox_relaxation_test.go` | Fall `egress-without-sandbox` 1 gekennzeichnet / 0; Spalte `flaggedText`; neue Prüfung `nestedRuntime` „no effect“ nie gekennzeichnet; Doc-Kommentar berichtigt | konform |
| `quiet_test.go` | `TestExecute_Generate_Quiet_KeepsEgressNoEffectDropsNestedRuntimeNoEffect`: `--quiet` zeigt Egress, unterdrückt `nestedRuntime`; `--quiet --json` gleich `--json`; Kennzeichen nicht in Text und Wire-Form | konform |
| `generate_sandbox_test.go` | Egress-Test prüft `SecurityRelaxation` am echten `GenerateService`; Podman-/Docker-Test prüfen `nestedRuntime` ungekennzeichnet | konform |
| `roadmap.md` | „Nichts in Arbeit.“ entfernt | konform (Auflage 1) |
| `open/README.md` | drei Zeilen entfernt, die nicht mehr in `open/` liegen | konform |
| `spec/spezifikation.md` | SPEC-043 Einstufungsliste und Historie; SPEC-044 unberührt (führt die Egress-Meldung nicht) | konform | <!-- d-check:ignore (Kennung im Fliesstext, Link in der Scope-Zeile) -->
| Handbuch | Zeilen 299 und 1417 wortgleich nachgezogen | konform (siehe D-2) |
| `CHANGELOG.md` | `[Unreleased]` nennt Egress sichtbar, `nestedRuntime` entfällt, nur `generate devcontainer` erzeugt, Text und `--json` unverändert | konform |

## DoD-Abgleich

- [x] **Code** — Evidence: `devcontainer_sandbox.go` (`relax(…)` für die Egress-Meldung ohne Profil); `nestedRuntime`-Meldungen weiter `warn(…)`; Relaxation-Test und CLI-Test wie in der DoD gefordert (siehe Diff-Tabelle). `make gates` grün.
- [x] **Spezifikation und Doku** — Evidence: SPEC-043 (Zeile Gekennzeichnet, Zeile Ungekennzeichnet, Historie), Handbuch 299/1417 (`generate devcontainer` als einziger Erzeuger, Egress sichtbar, `nestedRuntime` entfällt), CHANGELOG `[Unreleased]`. Alle drei sagen dasselbe wie Code und das ungeänderte Lastenheft (Ausnahme 1: „dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift“). SPEC-044 trägt die Egress-Meldung nicht und bleibt zu Recht ungeändert. <!-- d-check:ignore (Kennung im Fliesstext, Link in der Scope-Zeile) -->
- [x] **`make gates` grün** — Evidence: eigener Lauf, Exit 0 (siehe Sensors).
- [x] **Review durchgeführt, Report liegt vor** — Evidence: Review-Report oben, 0 HIGH/MEDIUM/LOW.
- [x] **Doku-Update (öffentlicher Vertrag)** — Evidence: Handbuch und CHANGELOG im Diff.
- [ ] **Closure-Notiz mit Steering-Loop-Lerneintrag** — offen: §7 steht auf „Wird bei der Closure gefüllt.“ (Closure-Schritt, Auflage 2).
- [ ] **Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7** — offen: Closure (Auflage 3).
- [ ] **Jedes Risiko aus §6 trägt einen Ausgang** — offen: beide Risiken stehen auf „weiter offen“ (Auflage 4). Befund aus Verifier-Sicht: Risiko 1 (Kennzeichen in `--json` oder Golden-Fall) ist **entfallen** (CLI-Test belegt, dass das Feld nicht in die Wire-Form gelangt; `mapWarningsToDiagnostics` bildet nur Code, Level, Message ab; Replay-Fälle prüfen Dateiinhalte). Risiko 2 (Plan-Aussage über Code) ist **nicht eingetreten**: Tabelle gegen `git diff --stat` abgeglichen, keine Abweichung (Verweis auf BEO-DOK/plan-aussage-ueber-code <!-- d-check:ignore (Beobachtungs-Kennung) --> im Register).
- [ ] **Drei Paarungen (Anker, Folge-Slice, Register)** — offen: Closure-Schritt.

Die offenen Punkte sind Closure-Arbeit des Planners und kein Mangel am Gelieferten; die Liefer-Punkte (Code, Spezifikation/Doku) sind erfüllt.

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | eigener Lauf auf HEAD, Exit 0 (lint, test, coverage-gate, docs-check, baseline-verify) |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass | siehe Abschluss unten (Lauf nach Anlage dieser Datei) |
| Mutationslauf Reviewer | Beleg gesichtet | sechs Mutationen M1..M6, fünf rot, M6 äquivalent; vom Verifier nicht wiederholt (Begründung: Gate-Lauf grün, Review-Evidenz mit benannten Testnamen) |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | CLI-Test `TestExecute_Generate_Quiet_KeepsEgressNoEffectDropsNestedRuntimeNoEffect`; `visibleWarnings` in `warnings.go` wertet `SecurityRelaxation` aus |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `TestGenerateDevcontainer_SandboxEgress_NotGenerated` (Kennzeichen am echten Service), Relaxation-Test Fall `egress-without-sandbox` |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `…SandboxPodman_NotGenerated`, `…SandboxDocker_NotGenerated`, Relaxation-Test (ungekennzeichnet) |
| [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) | Einstufungsliste und Historie im Diff; deckungsgleich mit Code und Handbuch |

Carveouts: Neu: none. Gelöst: none. Unverändert: alle.

Nicht ausgeführt: `make test-docker` (kein Docker-Pfad berührt), `make verify-depguard` (keine Importregel berührt), `make ci` (kein Security-/Release-Pfad).

Commit / Artefakt: `3063cb1..d08146d`.

## Diskrepanzen

| ID | Befund | Gewicht |
|---|---|---|
| D-1 | Keine Diskrepanz zwischen Plan, Code, Spezifikation, Handbuch, CHANGELOG und Lastenheft. | — |
| D-2 | Handbuch-`--quiet`-Zelle (Zeilen 299 und 1417) ist wortgleich doppelt und sehr lang (Review F-4). Sie ist sachlich wahr und DoD-konform. | Hinweis |

## Entscheidung zu den INFO-Findings

- **F-4 (überlanger Handbuch-Klammersatz, doppelte Zelle):** berührt die DoD nicht (beide Zeilen tragen den Zustand korrekt), also keine Closure-Bedingung. Vorschlag: **als Hinweis führen**, nicht in der Closure berichtigen. Begründung: Die Doppelung (Globale Optionen und Befehlsabschnitt) und die Länge stammen aus dem Bestand; eine Kürzung ändert öffentlichen Doku-Text ohne Anlass in diesem Slice und würde die Closure um einen Doku-Umbau erweitern (Überschreitung der Liefer-Punkte). Der Hinweis gehört als Beobachtung ins Register (Vorschlag: neuer Eintrag oder Zähler-Erhöhung nur, falls eine bestehende Klasse „Doku-Zelle zu lang / doppelt“ existiert; sonst Zähler 1 neu) und, falls der Projektinhaber es will, als Folge-Slice nach `open/`.
- **F-1 (kein End-to-End-Lauf APP↔CLI):** berührt die Closure nicht. Beide Hälften sind einzeln mutationsrot (M3, M4, M5), die Kopplung ist das Feld `SecurityRelaxation`; derselbe Schnitt trägt die übrigen Sandbox-Warnungen. Hinweis, keine Auflage.
- **F-2, F-3:** keine Closure-Wirkung (F-2 vom Verifier mit grünem `make gates` bestätigt; F-3 äquivalente Mutation).

## Auflagen für die Closure

1. **Roadmap-Ruhe-Marker:** Nach dem Move nach `done/` in `docs/plan/planning/in-progress/roadmap.md` unter „Offene Wellen“ „Nichts in Arbeit.“ **nicht** wiederherstellen, solange die Welle `welle-devcontainer-sshd-claude-code` dort als offen steht. Es gilt die Regel des Planning-Lifecycle (d-check `planning`): der Marker steht, wenn keine Slice-Datei in `in-progress/` liegt und keine offene Welle gelistet ist. Beim Schließen der Welle (mit diesem Slice schließbar, §5) entfällt der Welle-Eintrag, dann kehrt der Ruhe-Marker zurück. Reihenfolge nach `AGENTS.md` §3.3 Fall 2: erst Inhalt (DoD-Häkchen, §7), dann reiner `git mv` (`make slice-mv`). <!-- d-check:ignore (Kennung im Fliesstext, Link in der Scope-Zeile) -->
2. §7 Closure-Notiz mit Steering-Loop-Lerneintrag füllen; DoD-Häkchen setzen (Code, Spezifikation und Doku, `make gates`, Review, Doku-Update; die übrigen nach Closure).
3. **Register-Zähler:** BEO-DOK/plan-aussage-ueber-code <!-- d-check:ignore (Beobachtungs-Kennung) --> ist verkörpert (`AGENTS.md` §3.10); in diesem Slice trat die Klasse nicht erneut auf (Tabelle gegen Diff ohne Abweichung), der Zähler bleibt unverändert, in §7 „nicht eingetreten“ notieren. `BEO-CLI/kommentar-beschreibt-abwesendes` <!-- d-check:ignore (Beobachtungs-Kennung) --> ebenso: §3.9-Selbstcheck laut Review ohne Befund, Zähler unverändert. Neu zu notieren: Review F-4 (Handbuch-Zelle doppelt und überlang) als Zähler 1, falls der Planner die Klasse führen will; sonst in §7 „keine neue Beobachtung“.
4. §6: Risiko 1 auf „entfallen“, Risiko 2 auf „nicht eingetreten“ setzen (Belege oben).
5. Die Welle ist nach diesem Slice schließbar (§5): Welle-Closure (`welle-…-results.md`, Paarungen Anker, Folge-Slice, Register) ist ein eigener Schritt; der Slice-Move ist davon getrennt.
6. `open/README.md` bleibt wie im Diff (nur noch Slices, die in `open/` liegen).

## Fazit

DoD-Liefer-Punkte (Code, Spezifikation und Doku, `make gates`, Review, Doku-Update) **erfüllt**; die Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen) sind offen und liegen beim Planner. Keine Diskrepanz zwischen Plan und Code. Closure freigegeben unter den Auflagen oben.
