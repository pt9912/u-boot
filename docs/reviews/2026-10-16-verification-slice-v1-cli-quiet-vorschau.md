# Verification Evidence: slice-v1-cli-quiet-vorschau — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** [`slice-v1-cli-quiet-vorschau`](../plan/planning/done/slice-v1-cli-quiet-vorschau.md), Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md). Commits: `acefc1d` (Code und Tests), `27c6de5` (Handbuch, CHANGELOG), `7ee6c76` (Review-Fixes F-1/F-2 und Report). Review: [Report](2026-10-16-slice-v1-cli-quiet-vorschau.md) (0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO).

## Verification Evidence

Scope:
- Slice: [`slice-v1-cli-quiet-vorschau`](../plan/planning/done/slice-v1-cli-quiet-vorschau.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.10.0), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)
- Artefakte: `internal/adapter/driving/cli/` (`init.go`, `generate.go`, `add.go`, `remove.go`, `config.go`, `previewmode.go`, neuer Test `quiet_preview_test.go`), `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2; die Häkchen im Slice sind noch ungesetzt, sie folgen mit der Closure-Notiz):
- [x] `LH-FA-CLI-005` (0.10.0), `LH-FA-CLI-007`, `LH-FA-CLI-008` je Subkommando mit Vorschau belegt — Evidence: `quiet_preview_test.go`. Das Lastenheft 0.10.0 (Zeile 196) deckt die Regeltabelle wörtlich: Vorschau bleibt, Erfolgsmeldung eines tatsächlichen Schreibens entfällt auch bei `--diff` ohne `--dry-run`. <!-- d-check:ignore (Kennungen und Pfad als Zitat) -->
- [x] Testmatrix §1 gegen Tests: `TestExecute_Quiet_PreviewMatrix` läuft sechs Fälle (`init`, `generate`, `add`, `remove`, `config set`, `config set` mit mehreren Paaren) × vier Zeilen (`--dry-run`, `--dry-run --diff`, `--diff`, keine Option) und prüft Summary, Diff-Marker und leeres stderr; bei `--dry-run`-Zeilen ist die Ausgabe byte-gleich zum Lauf ohne `--quiet`. Vollständig: jede Zeile × jedes Subkommando. Boundary: `TestExecute_Quiet_ActualWriteStaysSilent` (Erfolgsmeldung entfällt, auch mit `--diff`), `TestExecute_Add_Quiet_DryRunNoOpPlanStays` (No-op-Plan, nur für `add`; die übrigen Subkommandos teilen den Zweig). Negative: `TestExecute_Quiet_PreviewErrorStays` (6 × 3 Optionssätze, `errors.Is`), `TestExecute_Quiet_PreviewJSONUnchanged` (byte-gleich, 6 Fälle), `TestExecute_Init_Quiet_DryRunPlanKeepsOnlyExceptionWarnings`. <!-- d-check:ignore (Testnamen als Zitat) -->
- [x] Doku: Handbuch (beide `--quiet`-Zeilen) und CHANGELOG nennen die Vorschau als Ergebnis; der CHANGELOG nennt nur noch „The rule that the last verbosity option wins follows in a later release“ als folgend (Grep: einzige Stelle). F-2 (Warnungs-Vorbehalt „same summary … notices … stay dropped“) in `7ee6c76` behoben.
- [x] `make gates` grün — selbst gefahren, Exit 0 (Coverage 91,10 %, Schwelle 90 %; d-check 266 Dateien, 0 Befunde).
- [x] Review durchgeführt — Report vorhanden; F-1 (Kommentarblock getrennt) und F-2 behoben, Diff `7ee6c76` auf `previewmode.go` gegengelesen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Offen.“; Auflage 1).
- [ ] Beobachtungs-Register — offen (Auflage 1).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: zwei Risiken stehen auf „offen“ (Auflage 2).
- [ ] Drei Paarungen — offen bis zur Closure (Auflage 3).

Plan-vs-Code-Diff (Slice §3):
- `init.go`, `add.go`, `remove.go`, `generate.go`, `config.go` (`runConfigSet`, `runConfigSetMany`) — Ist: gesetzt; Summary hängt an `summaryVisible(quiet, dryRun)`, Diff läuft davor, Fehlerpfade vor dem Zweig, JSON-Zweige unberührt. `init` filtert den Plan auf Ausnahme-Warnungen; `remove` erzeugt unter `PreviewDryRun` keine Purge-Warnungen (Review-Negativbefund, übernommen).
- `previewmode.go` „prüfen“ — Ist: Helfer `summaryVisible` (wie im Plan vorgeschlagen), Wahrheitstabelle unverändert.
- `quiet_modifying_test.go` „update“ — Abweichung: Plan nennt eine Änderung der Datei, der Diff legt stattdessen `quiet_preview_test.go` neu an; Bestandstests unverändert (wie im Plan verlangt). Der Plan-Text ist bei der Closure auf die neue Datei zu berichtigen (kleine Abweichung, kein Befund gegen die Zusage).
- Handbuch, CHANGELOG — Ist: gesetzt. Keine weiteren Dateien im Diff; kein `//nolint`, keine Gate-Lockerung, keine neuen Imports.

Bewertung der Risiken aus §6:
- **Auslegung `--diff` ohne `--dry-run`:** Ausgang „entfallen: Lastenheft 0.10.0 entscheidet, bestätigt vom Projektinhaber“ ist **korrekt**. Das Lastenheft 0.10.0 sagt ausdrücklich „… entfällt auch bei `--diff` ohne `--dry-run`“ (Änderungsvermerk 0.10.0: „Vereinbarung mit dem Projektinhaber“); Code und Pin (`ActualWriteStaysSilent`, Matrix-Zeile 3) folgen dem Wortlaut. Zusätzlich ist §1 („vom Projektinhaber zu bestätigen“) und die Überschrift der Regel auf „Auslegung … Lastenheft 0.10.0“ nachzuziehen, damit der Slice nicht eine offene Bestätigung behauptet.
- **Vergangenheitsform im Dry-Run (`remove`, `generate`, `config set`):** Ausgang „weiter offen“ mit **Träger-Slice**. Ein Nicht-Ziel ohne Slice-Kennung ist kein zulässiger Ausgang (DoD: „Folge-Slice mit Kennung“); „entfallen“ passt nicht, der Befund besteht fort. Nötig ist ein Folge-Slice, vorgeschlagene Kennung `slice-v1-cli-dry-run-would-formulierung` (noch nicht angelegt, Anlegen ist Sache des Planers); Anker `LH-FA-CLI-007`, Ausgang in §6 „weiter offen: Träger `<Kennung>`“, Eintrag in §7 Folge-Slices. <!-- d-check:ignore (vorgeschlagene Kennung, noch nicht angelegt) -->
- Konflikt mit [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/done/slice-v1-cli-quiet-nur-fehler.md) und Zwischenstand der Spec: „entfallen“ bleibt zutreffend (jener Slice liegt in `done/`; der Spec-Eintrag liegt im letzten Slice).

Bewertung der offenen Review-Findings:
- **F-3 (Vergangenheitsform):** nicht blockierend. Keine Abweichung vom Lastenheft (`--quiet` gibt aus, was ohne `--quiet` ausgegeben wird, byte-gleich gepinnt); der Befund ist vorbestehend. Er braucht den Träger-Slice oben.
- **F-4 (kein Dry-Run-Pin mit Warnungen für `generate`):** nicht blockierend. Verhalten korrekt, durch Lesen bestätigt (Ausnahme-Warnungen nach dem Summary-Zweig, wie ohne Dry-Run). Empfehlung: in §7 als benannte Grenze festhalten (oder, billiger, einen Pin nachziehen); ein Folge-Slice dafür ist unverhältnismäßig.

Mutationen (selbst gefahren, nur über `make test`, danach `git checkout -- <Datei>`; `git status` sauber):
| Nr | Mutation | Ergebnis |
| --- | --- | --- |
| M-1 | `remove.go`: `summaryVisible(…)` → `!flags.Quiet` (Dry-Run-Plan unterdrückt) | rot: `TestExecute_Quiet_PreviewMatrix/remove/dry-run`, `…/dry-run_diff` |
| M-2 | `add.go`: `!summaryVisible(…)` → `flags.Quiet && !flags.Diff` (Überkorrektur bei `--diff`, Dry-Run-Plan weg) | rot: `TestExecute_Add_Quiet_DryRunNoOpPlanStays`, `TestExecute_Quiet_ActualWriteStaysSilent/add`, `…PreviewMatrix/add/dry-run` |

Der Review-Report nennt fünf weitere Mutationen (M-a bis M-e); sie wurden hier nicht wiederholt.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (lint, test, coverage-gate 91,10 %, docs-check 0 Befunde, baseline-verify) |
| `make doc-trace` | pass | selbst gefahren: „83 Anforderung(en), 0 Waise(n).“; `LH-FA-CLI-005`, `-007` und `-008` je `ok` und tragen diesen Slice | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make test` (2 Mutationen) | rot wie erwartet | siehe Tabelle |
| `make test-docker` | not run | kein Docker-Pfad berührt: Änderungen im Driving-Adapter-Rendering, Fake-Use-Case-Tests und Doku |
| `make verify-depguard` | not run | keine Importregel geändert (`previewmode.go` importiert weiter nur `port/driving`) |
| `make docs-check` | pass | nach Anlegen dieser Datei erneut gefahren, siehe Übergabe |

Folge-Slices in der Welle: [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../plan/planning/in-progress/slice-v1-cli-verbosity-letzte-option-gewinnt.md) ist der einzige verbleibende Slice der Welle (Welle §3, Tabelle; Handbuch und CHANGELOG nennen nur noch „letzte Option gewinnt“ als folgend). Hinweis: `open/README.md` §Bestand führt die Zeile des Slice `slice-v1-cli-quiet-vorschau` noch, obwohl er in `in-progress/` liegt (Link zeigt korrekt dorthin); bei der Closure den Bestandseintrag entfernen bzw. mit dem `git mv` nachziehen. <!-- d-check:ignore (Dateiname und Ortsangabe als Zitat) -->

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `quiet_preview_test.go`, Handbuch, CHANGELOG; `doc-trace` ok |
| [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) | Matrix-Zeilen `--dry-run` und `--dry-run --diff` je Subkommando |
| [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) | Matrix-Zeilen mit `--diff`, `ActualWriteStaysSilent` |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `TestExecute_Quiet_PreviewJSONUnchanged` (byte-gleich, sechs Fälle) |

Carveouts: neu none · gelöst none · unverändert none berührt.

Commit / Artefakt: Slice-Commits `acefc1d`, `27c6de5`, `7ee6c76`; HEAD zum Prüfzeitpunkt `7ee6c76`.

## Verdikt

**Bestanden, mit Auflagen für die Closure.** Code und Tests decken die Regeltabelle in allen sechs Schreibpfaden und vier Optionszeilen ab; Gates grün, Waisen null, zwei Mutationen fallen rot, keine Verletzung von Lastenheft 0.10.0, Hard Rules oder Plan-Zusage.

Auflagen:
1. Closure-Notiz §7 mit Steering-Loop-Lerneintrag; Finding-Klassen aus dem Review (Einfügung trennt Kommentar von Gegenstand, CHANGELOG ohne Warnungs-Vorbehalt, Vergangenheitsform im Dry-Run, Warn-Pin nur für ein Subkommando); Register: Beobachtung eintragen oder „keine Beobachtung“ notieren; DoD-Häkchen, dann reiner `git mv` (`AGENTS.md` §3.3). <!-- d-check:ignore (Pfad als Zitat) -->
2. §6: Risiko `--diff` ohne `--dry-run` auf „entfallen: Lastenheft 0.10.0 entscheidet, bestätigt vom Projektinhaber“; Risiko Vergangenheitsform auf „weiter offen: Träger `<Kennung>`“ (Folge-Slice anlegen, Planer). §1 „vom Projektinhaber zu bestätigen“ berichtigen; §3-Zeile `quiet_modifying_test.go` auf `quiet_preview_test.go` berichtigen. <!-- d-check:ignore (Dateinamen als Zitat) -->
3. §7 Folge-Slices: `slice-v1-cli-verbosity-letzte-option-gewinnt` und der neue Slice für die „Would“-Formulierung; als Grenze benennen: kein Dry-Run-Pin mit Warnungen für `generate` (Review F-4), No-op-Pin nur für `add`. <!-- d-check:ignore (Kennung als Zitat) -->
4. Paarungen: `open/README.md` §Bestand bereinigen (Eintrag des Slice aus `in-progress/`); kein Release-Cut vor `slice-v1-cli-verbosity-letzte-option-gewinnt` (CHANGELOG „follows in a later release“). <!-- d-check:ignore (Pfad und Kennung als Zitat) -->
