# Review-Report: slice-harness-werkzeug-schicht-ai-harness-init — 2026-10-16

**Review-Art:** Code — Diff gegen Plan, MR-011 und Hard Rules.

**Gegenstand:** Diff `28cb8fe..HEAD` (Commits `489eb89`, `45fb0ad`): Übernahme ausgewählter Teile der `ai-harness-init`-Werkzeug-Schicht.

**Skill:** `.harness/skills/reviewer.md` @ `45fb0ad` · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- `slice-harness-werkzeug-schicht-ai-harness-init` (Plan, Entscheidungen, DoD) <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->
- `MR-011` (Adaption), `MR-004`, `MR-005`, `MR-003`
- `AGENTS.md` (Hard Rules 3.1, 3.3, 3.6, 3.7, §4, §6), `harness/roles.md`
- Referenzquelle der Originale: `/tmp/ahi080` (nur gelesen)
- Läufe in einem Wegwerf-Klon unter `/private/tmp/rv` (bash 3.2.57 und 5.3): Guard, Stop-Hook, `record-gates`, `working-tree-hash`, Commit-Wächter, `make hooks-install`, `make slice-mv`, `make -n -j8 gates`. Keine Host-Go-Toolchain, nichts im Repo geändert. <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Kommentare nennen Dateien und Targets, die es im Repo nicht gibt: `make archive-welle` samt `KENNUNG=<K>` (zweimal), `test/slice-mv.bats` (ein `test/`-Verzeichnis existiert nicht), „Root-Aggregator (gen)". Ein Leser, der `make archive-welle` aufruft oder die Tests sucht, läuft ins Leere. | `AGENTS.md` §4 („Kein Target nennen, das im Makefile nicht existiert — auch nicht in Prosa"), §3.7 | `harness/mk/hooks-install.mk` · „`make slice-mv` und `make archive-welle` committen intern"; `tools/harness/slice-mv.sh` · „test/slice-mv.bats (Fälle `quelle:`)" | ja — `grep -rn "archive-welle\|slice-mv.bats" harness tools` | Verweis auf nicht vorhandenes Target/Testdatei in übernommenem Kommentar | <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->
| F-2 | LOW | Die Kopfkommentare behaupten Eigentümer und Verhalten, die MR-011 verneint: „emittiert von ai-harness-init", „wird bei jedem Bootstrap kanonisch neu geschrieben, eine Änderung an ihr überlebt den nächsten Lauf nicht", „mitemittiertes Regelwerk/Fragment". MR-011 sagt: kein Werkzeug schreibt diese Dateien neu, das Repo führt sie. Wer dem Kommentar glaubt, setzt Anpassungen an der falschen Stelle (Umgebung statt Datei) oder scheut die Änderung. | `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist); `MR-011` §Adaption | `tools/harness/slice-mv.sh` · „Diese Datei wird bei jedem Bootstrap kanonisch neu geschrieben"; `harness/mk/slice-mv.mk` · „wird bei jedem Bootstrap kanonisch neu geschrieben"; `harness/mk/enforce.mk` · „emittiert von ai-harness-init" | ja — Lesen von MR-011 gegen die Köpfe | Herkunfts-/Eigentümer-Kommentar widerspricht Adaption | <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->
| F-3 | LOW | `working-tree-hash.sh` sortiert mit `sort -zu` ohne `LC_ALL=C`; der Hash ist locale-abhängig (gemessen: `LC_ALL=C` und `en_US.UTF-8` liefern auf demselben Baum verschiedene Hashes). Läuft `make gates` in einer anderen Locale als der Stop-Hook, stimmt der Stempel nicht, und der Hook blockt trotz grünem Lauf (einmal; danach greift der Schleifen-Schutz). | `MR-011`; Kopplung `record-gates` ↔ Stop-Hook (Kommentar „keine Logik-Dopplung") | `tools/harness/working-tree-hash.sh` · „`| sort -zu \|`" | ja — Hash unter zwei Locales vergleichen | Locale-abhängiger Hash bei gemeinsamer Quelle zweier Leser |
| F-4 | LOW | Der PreToolUse-Guard prüft nur den ersten Token je Segment. Schlüsselwörter der Shell-Grammatik als Kopf lassen Toolchains durch: `if true; then go test; fi` und `for f in a; do go vet; done` wurden nicht geblockt (`ls`-Kopf nach `;` wird dagegen erkannt). Der Kopfkommentar nennt den Guard Stolperdraht, nicht Sandbox; `then`/`do`/`else` sind aber keine bewusst genannte Lücke. | `AGENTS.md` §3.1 | `.claude/hooks/pretooluse-command-guard.sh` · `PREFIXES="sudo env command exec nice time xargs eval"` | ja — Hook mit JSON-Eingabe füttern | Guard-Präfixliste ohne Shell-Schlüsselwörter |
| F-5 | INFO | `make slice-mv` ist nicht rundlauf-neutral: Hin- und Rückweg (`in-progress` → `next` → `in-progress`) hinterließ praefixlose Links als `](../in-progress/<datei>)` statt `](<datei>)` (3 Dateien, 4 Zeilen gegenüber dem Ausgangsstand). Die Links lösen auf; es ist kosmetisch. | `AGENTS.md` §3.3 | `tools/harness/slice-mv.sh` · `rewrite_incoming_bare_in_file` | ja — Klon, zwei `slice-mv`-Läufe, `git diff` | Nicht-idempotente Link-Umschreibung | <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->
| F-6 | INFO | `working-tree-hash.sh` setzt `sha256sum` voraus. Es liegt auf dem Host (`/sbin/sha256sum`), steht aber nicht als Host-Abhängigkeit in `harness/README.md`. Fehlt es, bricht `record-gates` (pipefail, kein Stempel) und der Stop-Hook blockt: fail-closed, kein stilles Offen. | `AGENTS.md` §3.1 (Host braucht nur Docker und `make`) | `tools/harness/working-tree-hash.sh` · „`sha256sum -- \"$f\"`" | ja — `PATH` ohne `sha256sum` | Ungenannte Host-Abhängigkeit (git, sha256sum) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `Makefile`: Ordnungskante `GATE_CHECKS` → `record-gates` → `gates` | geprüft, ohne Befund — `make -n -j8 gates` zeigt drei `docker build`, dann `docker run` (docs-check), erst danach `record-gates.sh` und das `[gates]`-echo; `record-gates` hängt an allen vier Checks als Prerequisite, läuft also nur nach grünen Checks. |
| `Makefile`: Gate-Liste (AGENTS.md §3.6) | geprüft, ohne Befund — `lint test coverage-gate docs-check` unverändert (`GATE_CHECKS`), `.DEFAULT_GOAL` bleibt `help`, `ci` und `fullbuild` ziehen `gates` unverändert nach. Keine Lockerung. |
| `harness/mk/*.mk`, `make help`, `.d-check.yml targets.makefiles` | geprüft, ohne Befund — die drei Fragmente stehen in `targets.makefiles`; `record-gates`, `slice-mv`, `hooks-install` erscheinen in `make help` und in `harness/README.md` §Sensors (eigene Werkzeug-Tabelle, „kein Gate", Bindung MR-011). | <!-- d-check:ignore (Make-Target bzw. Slice-Kennung im Zeitdokument) -->
| `tools/harness/record-gates.sh`, `.claude/hooks/stop-require-gates.sh` | geprüft, ohne Befund — Stempel und Hash gleiche Quelle; ohne neuen HEAD frei, mit `STOP_GATE_STRENG=1` Block bei ungedecktem Inhalt, Fehler enden fail-closed (Trap Exit 2), Schleifen-Schutz vorhanden; `head_wert` in beiden Fassungen gleich. |
| `.githooks/commit-msg`, `commit-msg-traceability.sh`, `make hooks-install` | geprüft, ohne Befund — unter bash 3.2 laufen Ablehnung ohne Kennung, `LH-`, `slice-<name>`, Merge-Ausnahme und die `slice-mv`-Messages wie zugesagt; `hooks-install` setzt `core.hooksPath` und bricht ohne Träger ab. | <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
| `tools/harness/slice-mv.sh`, `slice-mv.mk` | geprüft, ohne Befund — im Klon unter bash 3.2/BSD-sed Move- und Nachzug-Commit getrennt, reiner Move-Commit, Abbruch bei unsauberem Baum; §3.3 (Regelfall) eingehalten, bei `→ done` steht die Closure-Notiz laut Command vor dem Move. Ausgenommen: F-5. | <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
| `.claude/hooks/pretooluse-command-guard.sh`, `extract-command.awk` | geprüft, außer F-4 ohne Befund — `go`, `gofmt`, Pfad-Form, `VAR=`-Präfix, `bash -lc`, Pipes, `$(`/Klammern, `\u`-Escape und Garbage-Eingabe blocken (fail-closed); `make gates`, `git commit -m` mit Toolchain-Wort im Text und `docker run … go test` laufen durch; der dokumentierte quote-blinde False-Positive (`a; go b` in einer Message) tritt wie beschrieben auf. Läuft unter bash 3.2. |
| `.gitignore` | geprüft, ohne Befund — `.claude/settings.local.json`, `.claude/scheduled_tasks.lock` und `.harness/state/` bleiben ignoriert; `settings.json`, `hooks/`, `agents/`, `commands/` und `.harness/stop-gate-streng` sind versioniert (`git check-ignore` im Klon). Ausführrechte der Skripte gesetzt. |
| `.claude/settings.json` | geprüft, ohne Befund — `PreToolUse`(Bash) und `Stop` zeigen auf vorhandene Skripte über `$CLAUDE_PROJECT_DIR`. |
| `.claude/agents/*`, `.claude/commands/implement-slice.md` vs. `AGENTS.md` §6, `harness/roles.md` | geprüft, ohne Befund — keine Verweise auf `planner`, `plan-welle`, `close-welle`, Beobachtungs-Register, Span-Erfassung oder Archivierung; Rollenfolge, „kein Self-Review" und genannte Sensors (`make test-docker`, `make verify-depguard`) existieren. |
| `MR-011` und `harness/conventions.md`-Zeile, Plan-DoD | geprüft, ohne Befund — Ablehnungsliste deckt sich mit dem Plan §2; der Ist-Zustand (kein `harness/mk/ai-harness-init.md`, `baseline-verify` ersetzt, `doc-gate.mk` nicht übernommen) stimmt; die Kommentar-Abweichung steht in F-1/F-2. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Verweis auf nicht vorhandenes Target/Testdatei in übernommenem Kommentar · Herkunfts-/Eigentümer-Kommentar widerspricht Adaption · Locale-abhängiger Hash bei gemeinsamer Quelle zweier Leser · Guard-Präfixliste ohne Shell-Schlüsselwörter · Nicht-idempotente Link-Umschreibung · Ungenannte Host-Abhängigkeit (git, sha256sum)

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) verletzt den ausdrücklichen Satz aus `AGENTS.md` §4 (kein nicht vorhandenes Target nennen, auch nicht in Prosa); die Korrektur ist ein reiner Kommentar-Fix. Die mechanische Substanz (Ordnungskante, unveränderte Gate-Liste, fail-closed-Hooks, `.gitignore`) ist ohne Befund. Gate-Läufe (`make gates`) hat dieser Review nicht ausgeführt; das ist Sache des Implementers und Verifiers.

**Übergabe:** Findings an den Implementer; Finding-Klassen in die Slice-Closure §7. Der Report ist ein Lauf-Beleg und ersetzt keine Verifikation der DoD.
