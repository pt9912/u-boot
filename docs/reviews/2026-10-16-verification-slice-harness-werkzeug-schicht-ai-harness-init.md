# Verification Evidence — Slice Werkzeug-Schicht ai-harness-init

Rolle: Verifier (frischer Kontext), Datum 2026-10-16, Range `28cb8fe..2fd8c65`.
Gegenstand: [`Slice`](../plan/planning/done/slice-harness-werkzeug-schicht-ai-harness-init.md),
Review: [`Report`](2026-10-16-slice-harness-werkzeug-schicht-ai-harness-init.md),
Adaption [`MR-011`](../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md).

## Verification Evidence

Scope:
- Slice: slice-harness-werkzeug-schicht-ai-harness-init <!-- d-check:ignore (Slice-Kennung im Zeitdokument) -->
- IDs: [`MR-011`](../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md), [`MR-005`](../../harness/conventions/MR-005-gate-haltung-docs-check.md), [`MR-004`](../../harness/conventions/MR-004-regelwerk-vendored.md); keine `LH-*`/`ADR-*` berührt (Slice: „Berührte Spec-Stellen: —")
- Artefakte: `Makefile`, `harness/mk/`, `tools/harness/`, `.githooks/`, `.claude/`, `.gitignore`, `.d-check.yml`, `harness/README.md`, `harness/conventions.md`

DoD-Abgleich:
- [x] Entscheidung je Fragment dokumentiert (§2) — bestätigt: Plan §2 und MR-011 (Übernahme-/Ablehnungsliste) decken sich; Ist-Diff enthält genau die übernommenen Teile (`slice-mv`, `hooks-install`, `enforce`, Guard samt `extract-command.awk`/`blocked/go`, Stop-Hook, vier Agenten plus `verifier`, `implement-slice`); abgelehnt und tatsächlich abwesend: `plan-welle`, `close-welle`, `planner`, `baseline-verify`, `doc-gate.mk`, `harness/mk/ai-harness-init.md`. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
- [x] Einbindung in `Makefile`, `harness/README.md` §Sensors, `.d-check.yml`, MR-011 — bestätigt: `include harness/mk/*.mk` (Makefile:207); `.d-check.yml` `targets.makefiles` nennt die drei Fragmente einzeln; §Sensors führt `record-gates`, `slice-mv`, `hooks-install` (README Z. 79–81) und `gates` mit Hinweis `record-gates`; MR-011 und `harness/conventions.md`-Zeile vorhanden. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
- [x] `make gates` grün und `make docs-check` grün — bestätigt, selbst gefahren (siehe Sensors). Gate-Liste unverändert: `gates: record-gates` mit `record-gates: $(GATE_CHECKS)`, Echo „lint + test + coverage-gate + docs-check"; Stempel `.harness/state/gates-passed.{head,diffsha}` vorhanden; `stop-require-gates.sh` antwortet `{"decision":"approve"}` (Exit 0).
- [x] Review durchgeführt; MEDIUM und die zwei LOW (Kopf/Hash) behoben — bestätigt im Code: `grep -rnE "archive-welle|bats|emittiert|kanonisch neu"` über `harness/mk`, `tools/harness`, `.githooks`, `.claude` liefert 0 Treffer (F-1 und F-2 behoben); `tools/harness/working-tree-hash.sh:14` nutzt `LC_ALL=C sort -zu` (F-3 behoben). Review-Zählung (0 HIGH, 1 MEDIUM, 3 LOW, 2 INFO) stimmt mit der Tabelle F-1..F-6 überein. F-4 (Guard erkennt nur erstes Wort) bleibt laut Slice §4 offen und ist dort als Risiko benannt; F-5/F-6 als entfallen begründet (nicht von mir neu bewertet, nur Konsistenz geprüft).
- [ ] `make hooks-install` im eigenen Klon — offen (bewusst): `git config core.hooksPath` ist leer. Closure-Notiz §5 „Offen." ebenfalls bewusst offen; Slice liegt korrekt in `in-progress/`.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0; `coverage-gate: OK — coverage 90.90% meets threshold 90%`; `d-check: 218 Datei(en) geprüft, 0 Befund(e)`; `[gates] lint + test + coverage-gate + docs-check green` |
| `make docs-check` | pass | `d-check: 218 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Berichts) |
| `tools/harness/fetch-baseline-cache.sh --verify` | pass | `verify ok (54 Dateien, vollständig)` |
| Stop-Hook `stop-require-gates.sh` | pass | `{"decision":"approve"}`, Arbeitsbaum sauber |

Traceability:
| ID | Beleg |
| --- | --- |
| `MR-011` | Datei und Ist-Diff stimmen überein (Fragmente repo-geführt, kein `ai-harness-init.md`, Ablehnungsliste); Kopfkommentare widersprechen ihr nicht mehr (grep leer) |
| `MR-005` | `d-check.mk` weiter eingebunden, `doc-gate.mk` nicht übernommen; docs-check grün |
| `MR-004` | `fetch-baseline-cache.sh --verify` grün |
| `AGENTS.md` §3.6 | Gate-Liste von `make gates` unverändert, `record-gates` nur angehängt |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker` — kein Docker-Runtime-/Produktpfad berührt (nur Harness-Werkzeug).
- `make verify-depguard` — keine Importregeln berührt.
- Funktionstest von `slice-mv`/Hook-Pfaden durch den Verifier — nicht wiederholt; Review belegt Klon-Erprobung, Gates laufen nicht über `slice-mv`. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->

Abweichungen / Hinweise:
- Keine DoD-Verletzung gefunden. Kein Plan-vs-Code-Widerspruch.
- Der docs-check-Lauf deckt diesen Bericht erst ab, wenn er committed/vorhanden ist; ein erneuter `make docs-check` danach ist beim Closure ohnehin fällig.
- Offene Risiken laut Slice §4: F-4 (Guard-Lücke bei `if`/`for`-Rümpfen), Commit-Hook weist kennungslose Messages ab.

Commit / Artefakt:
- Geprüfter Stand: `2fd8c65`
