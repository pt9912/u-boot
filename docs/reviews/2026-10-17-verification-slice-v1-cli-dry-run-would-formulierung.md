# Verification Evidence: slice-v1-cli-dry-run-would-formulierung — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-v1-cli-dry-run-would-formulierung`](../plan/planning/in-progress/slice-v1-cli-dry-run-would-formulierung.md), Lösung A (Entscheidung des Projektinhabers 2026-10-17), Commits `73f1f4b`, `1536c8e`, `27723e2` (Range `874b6b8..HEAD`, ohne Welle). Maßstab: §1, §2 (DoD), §3 (Plan), §6, [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe). Review-Report: [`2026-10-17-slice-v1-cli-dry-run-would-formulierung.md`](2026-10-17-slice-v1-cli-dry-run-would-formulierung.md).

## Verification Evidence

Scope:
- Slice: [`slice-v1-cli-dry-run-would-formulierung`](../plan/planning/in-progress/slice-v1-cli-dry-run-would-formulierung.md)
- IDs: [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)
- Artefakte: `internal/adapter/driving/cli/{remove,generate,config}.go`, Tests `dry_run_would_test.go`, `quiet_preview_test.go`, `remove_acceptance_test.go`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `roadmap.md` <!-- d-check:ignore (Pfad-Zitate) -->

DoD-Abgleich (§2 des Slice; im Slice sind die Kästchen noch alle leer, die Closure hakt sie ab):
- [x] Auslegung entschieden (Lösung A), Tests je Subkommando Happy/Boundary/Negative — Evidence: Code gelesen. `printRemoveSummary` schaltet auf `previewMode == PreviewDryRun`, `printGenerateSummary` und `printConfigSetSummary` auf `flags.DryRun`, `runConfigSetMany` reicht `flags.DryRun` durch. Tests: Happy `TestDryRun_SummarySaysWould` (je Subkommando, `--dry-run` und `--dry-run --diff`, mit und ohne `--quiet`), Boundary `TestDryRun_ActualWriteKeepsPastForm` (ohne Option und `--diff` allein), Negative `TestDryRun_JSONAndErrorsUnchanged` (`--json` ohne Would-Zeile mit `"dryRun":true`, Fehler bleibt Fehler). Zusätzlich je `GenerateAction` (`TestDryRun_GenerateActions`, Drift-Wächter `TestDryRun_GenerateEveryActionHasADryRunForm`) und der Mischfall (`TestDryRun_ConfigSetManyMixedPairs`). **Hinweis:** §7 des Slice hält die Entscheidung des Projektinhabers noch nicht fest („Offen.“); das gehört in die Closure-Notiz.
- [x] Doku-Update — Evidence: Handbuch §„Vorschau statt Überraschung“ nennt die `Would …`-Zeilen, die Vergangenheitsform ohne `--dry-run`, `--diff` allein und die unveränderten No-op-/`--json`-Fälle; stimmt mit dem beobachteten Verhalten überein (siehe Praxislauf). CHANGELOG `[Unreleased]` · Changed nennt alle Formen, den Quiet-Zusatz und den Skript-Hinweis auf `--json`.
- [x] `make gates` grün — selbst gefahren, Exit 0, `coverage-gate: OK — coverage 91.20% meets threshold 90%`.
- [x] Review durchgeführt, Report liegt vor (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO); Findings siehe unten.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 steht auf „Offen.“); gehört in die Closure.
- [ ] Beobachtungs-Register fortgeschrieben oder Vermerk in §7 — offen (Closure).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: das erste Risiko steht noch auf „weiter offen: Entscheidung des Projektinhabers“, obwohl entschieden; in der Closure auf „entfallen“ setzen.
- [ ] Die drei Paarungen — offen (Closure).

Review-Findings am Code geprüft:
| Finding | Befund am Code |
| --- | --- |
| F-1 CHANGELOG-Satz | behoben: „`--quiet --dry-run` still shows the summary …, now with the `Would …` wording“; Test `quiet`-Läufe pinnen es, der Praxislauf bestätigt es |
| F-2 `default`-Zweig nicht dry-run-aware | Zweig selbst unverändert (in §3 als Grenze benannt); Drift-Wächter `TestDryRun_GenerateEveryActionHasADryRunForm` läuft über alle `GenerateAction` bis `String()=="unknown"` und fordert Präfix `Would ` nur im Dry-Run; fängt eine neue Aktion ohne Dry-Run-Form |
| F-3 Mischfall `config set` | behoben: `TestDryRun_ConfigSetManyMixedPairs`; im Praxislauf bestätigt (`already demo; no changes.` plus `Would set devcontainer.enabled (unset) → true.`) |
| F-4 JSON-Negative zu breit | behoben: der Test prüft statt „Would“ überall die konkrete Dry-Run-Zeile des Subkommandos und `"dryRun":true`; Praxislauf: JSON ohne Would-Text |
| F-5 Host-`python3` | Verstoß gegen `AGENTS.md` §3.1 ist in §3 des Slice benannt; der Diff von `remove_acceptance_test.go` umfasst genau Kommentar-Berichtigung und eine Headline-Assertion; Nachweis lief über `make test`/`make gates`. Ohne Wirkung auf das Ergebnis, in die Closure-Notiz übernehmen |

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0, Coverage 91,20 % |
| `make doc-trace` | pass | selbst gefahren: `83 Anforderung(en), 0 Waise(n)` |
| `make build` + Praxislauf (Runtime-Image, Wegwerf-Verzeichnis, danach gelöscht) | pass | `init`, `add`, `remove`, `generate` (changelog/readme/env-example/devcontainer), `config set` einzeln und mit mehreren Paaren, jeweils mit und ohne `--dry-run`, mit `--diff`, `--quiet`, `--json` — Ergebnisse unten |
| `make docs-check` | siehe Übergabe | nach Anlage dieser Datei gefahren |

Praxislauf (Docker-Image `u-boot:latest`, Verzeichnis `~/uboot-verify-dryrun`, anschließend entfernt):
- `init --dry-run`: `Would initialize … Would create:`; ohne: `Initialized … Created:` (unverändert). `add --dry-run`: `Would add service … Would change:`; ohne: `Added service …` (unverändert).
- `remove postgres --dry-run`, `--dry-run --diff` (Diff vorab, Zusammenfassung gleich), `--dry-run --quiet`: `Would remove service "postgres".` + `Would change:`; `remove postgres` und `remove postgres --diff`: `Removed service "postgres".` + `Changed:`; bereits entfernt: `… already disabled; no changes.` in beiden Modi; `--dry-run --json`: Envelope mit `dryRun:true`, kein Wort „Would“ im Text; unbekannter Dienst und nicht initialisiertes Projekt: Fehler, Exit 10, unverändert (auch `--json`).
- `generate readme/env-example --dry-run`: `Would update … managed block (…)`; `generate devcontainer --dry-run`: `Would generate devcontainer (…)`; ohne `--dry-run`: `Updated … managed block`, `Generated devcontainer …`; `--diff` allein: Vergangenheitsform; No-op identisch in beiden Modi; `--json`: `"action":"no-op"`/Envelope unverändert. Die Aktion `Repaired` (`Would repair`) war im Lauf nicht herstellbar; sie ist durch `TestDryRun_GenerateActions/repaired-manual` und den Drift-Wächter belegt.
- `config set` einzeln: `config: Would set project.name demo → demo2.`, mit `--diff` und `--quiet` ebenso; ohne: `config: project.name demo2 → demo3.`; `--json` ohne Would-Text. Mehrere Paare: `Would set` je geändertes Paar, `already demo; no changes.` für unverändertes; ohne `--dry-run` Vergangenheitsform; `--quiet --dry-run` zeigt die Zeilen.

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) | `dry_run_would_test.go` (Happy/Boundary/Negative, Drift-Wächter), `quiet_preview_test.go`, `TestRemove_PurgeDryRun_HumanMode_NoVolumeWarningPollution`; Praxislauf |
| [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) | `--dry-run --diff` und `--diff` allein in `TestDryRun_SummarySaysWould` / `TestDryRun_ActualWriteKeepsPastForm`; Praxislauf |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `--json`-Negativ-Test; CHANGELOG-Skript-Hinweis |

Plan-Tabelle §3 gegen `git diff --stat 874b6b8..HEAD` (§3.10): 11 Dateien. Die zehn Dateien außerhalb des Slice-Plans (`CHANGELOG.md`, `roadmap.md`, `benutzerhandbuch.md`, `config.go`, `generate.go`, `remove.go`, `dry_run_would_test.go`, `quiet_preview_test.go`, `remove_acceptance_test.go` sowie der Review-Report) haben bis auf den Report je einen Plan-Eintrag; der Report ist DoD-Artefakt und braucht keine Zeile. Kein Plan-Eintrag ohne Diff. Kein `nolint` im Diff. Keine Diskrepanz.

Diskrepanzen und Beobachtungen:
- INFO (vorbestehend, nicht im Diff, kein Slice-Gegenstand): Im Dry-Run erscheinen auf stderr weiterhin `level=INFO`-Logzeilen in Vergangenheitsform (`config set: updated`, `generate: updated block`, `generate devcontainer: created`). Kandidat für das Beobachtungs-Register, nicht für diesen Slice.
- Die Roadmap trägt keinen Ruhe-Marker mehr („Nichts in Arbeit.“ entfernt, Plan §3), solange der Slice in `in-progress/` liegt.

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: alle

Nicht ausgeführt:
- `make test-docker` — kein Docker-Runtime-Pfad berührt.
- `make verify-depguard` — keine Importregel berührt.
- `make ci` (govulncheck, image-scan) — kein Security-/Release-Pfad berührt.
- Eigene Mutationen — nicht wiederholt; die Mutationsläufe M1–M6 des Reviews wurden nicht erneut gefahren, der Praxislauf deckt dieselben Zweige am Binary ab.

Commit / Artefakt:
- `73f1f4b`, `1536c8e`, `27723e2` (HEAD zum Prüfzeitpunkt)

## Urteil

Die Code- und Doku-Lieferung der Lösung A ist erfüllt; die Prozess-Punkte der DoD (Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen) sind offen und gehören der Closure.

Auflagen für die Closure:
1. §7 füllen: Entscheidung des Projektinhabers (Lösung A, 2026-10-17), Finding-Klassen aus dem Review, Host-`python3`-Vorfall, Vermerk zu den INFO-Logzeilen im Dry-Run (Register oder „keine Beobachtung“), Steering-Loop-Eintrag, Paarungen.
2. §6, erstes Risiko, auf „entfallen: entschieden“ setzen.
3. DoD-Kästchen abhaken (erst nach der Closure-Notiz), danach reiner `git mv` nach `done/` (`AGENTS.md` §3.3).
4. Nach dem Move: Ruhe-Marker („Nichts in Arbeit.“) in `roadmap.md` unter „Offene Wellen“ vom Planner wieder setzen; `make doc-planning` bzw. `make docs-check` danach prüfen.
