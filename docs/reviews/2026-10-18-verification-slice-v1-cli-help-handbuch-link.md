# Verification: slice-v1-cli-help-handbuch-link — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Slice liegt bei Closure in done/) -->

**Rolle:** Verifier — „Bauen wir es richtig?" gegen Plan, DoD und Entscheidungen.
**Gegenstand:** Commit-Range `2ab2ee6..HEAD` (HEAD `8297111`), Slice `slice-v1-cli-help-handbuch-link` (`docs/plan/planning/in-progress/`). <!-- d-check:ignore (Slice-Kennung) -->
**Eingang:** Review [`2026-10-18-slice-v1-cli-help-handbuch-link.md`](2026-10-18-slice-v1-cli-help-handbuch-link.md) (0 HIGH, 0 MEDIUM, F-1/F-2 LOW, F-3/F-4 INFO).

## Verification Evidence

Scope:
- Slice: `slice-v1-cli-help-handbuch-link` <!-- d-check:ignore (Slice-Kennung) -->
- IDs: [`LH-FA-CLI-002`](../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe), [`LH-FA-CLI-003`](../../spec/lastenheft.md#lh-fa-cli-003--versionsausgabe), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort)
- Artefakte: `internal/adapter/driving/cli/root.go` (`rootLongText`, `releaseVersionPattern`, `repoURL`, `manualPath`), `internal/adapter/driving/cli/cli_test.go` (`helpForVersion`, `TestExecute_RootHelpManualLinks`), `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md`, Slice-Plan, Review-Report

DoD-Abgleich (Stand des Slice: alle Häkchen im Plan noch offen, das ist Closure-Sache):
- [x] Code: Tag-Link nur bei `X.Y.Z` und `X.Y.Z-rc.N`, sonst nur `main` — Evidence: `root.go` gelesen; Praxislauf Runtime-Image: `0.11.0` und `0.12.0-rc.1` zeigen README-, Tag- und `main`-Link; `0.12.0-rc1` und Default `0.11.0-dev` zeigen README und `main`, keinen Tag-Link.
- [x] Test: 13 Fälle in `TestExecute_RootHelpManualLinks` (`0.11.0`, `1.20.3`, `0.12.0-rc.1` mit Tag-Link; `0.11.0-dev`, `0.0.0-test`, `0.12.0-rc1`, `0.12.0-beta.1`, `0.12.0-rc.01`, `01.2.3`, `1.2`, `1.2.3+meta`, `v0.11.0`, leer ohne `/blob/v`); erwartete Strings hartkodiert. `TestExecute_RootHelpPointsToUserDocs` unverändert und grün. Mutationen: nicht neu gefahren, vom Reviewer belegt (6 von 7 gekillt, eine bewusst nicht gepinnt); die Fälle decken die DoD-Mutationen (kein `v`, Regex gelockert, Verzweigung entfernt) strukturell ab.
- [x] Doku: `CHANGELOG.md` `[Unreleased]` `### Added` mit Halbsatz zur Grenze gelöschter Vorab-Tags; Handbuch-FAQ „Wo finde ich alle Optionen eines Befehls?" trägt den Satz.
- [x] `releasing.md` und Cut-Vorlage unverändert — Evidence: nicht im `git diff --stat`; §7-Notiz zum Cut-Handarbeitsschritt ist Closure-Auflage.
- [x] `make gates` grün, `make doc-trace` ohne Waise — Evidence: siehe Sensors.
- [x] Review-Report liegt vor (`docs/reviews/2026-10-18-slice-v1-cli-help-handbuch-link.md`). <!-- d-check:ignore (Slice-Kennung) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: §7 ist „Wird bei der Closure gefüllt."
- [ ] Beobachtungs-Register fortgeschrieben oder §7-Notiz — offen (siehe Auflagen).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: alle fünf Risiken stehen auf „offen" bzw. ohne Closure-Ausgang.
- [ ] Drei Paarungen (Anker · Folge-Slice · Register) — offen bis zur Closure.

Review-Findings am Code und Text:
- F-1 (LOW) behoben: Kommentare zu `releaseVersionPattern` („subset of what the publish workflow accepts") und `rootLongText` („because a tag of that form is not promised") stehen im Indikativ und decken sich mit `publish.yml`.
- F-2 (LOW) behoben: Plan §2 (Code- und Test-Punkt) und §6 (erstes Risiko) nennen `0.12.0-rc1` und `0.0.0-test` als „nur `main`"; Volltextsuche nach `rc1`, `-dev`, `0.0.0-test`, `Tag-Link` im Slice ohne Widerspruch zur Festlegung (Tag-Link nur bei `X.Y.Z` und `X.Y.Z-rc.N`).
- F-3 (INFO) benannt: CHANGELOG-Halbsatz und Plan §3 „Hinweise aus dem Review"; Ausgang in §6/§7 bei Closure.
- F-4 (INFO) benannt: Plan §3 „Hinweise" (Regex je Aufruf kompiliert, 200 ms nicht gemessen).

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0; `coverage-gate: OK — coverage 91.10% meets threshold 90%` |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n)." |
| `make doc-immutable RANGE="2ab2ee6..HEAD"` | pass | 9 Commits, „0 Befund(e)" |
| `make docs-check` | pass | Teil von `make gates`; nach Anlage dieser Datei erneut gefahren (siehe Bericht) |
| Praxislauf Runtime-Image (`UBOOT_VERSION` `0.11.0`, `0.12.0-rc.1`, `0.12.0-rc1`, Default `-dev`) | pass | `--help` je Fall wie oben; `--version` unverändert (`u-boot version <v>`); `init --help` ohne `blob`/`manual`-Zeile; Images danach per `docker rmi` entfernt |
| `gh api repos/pt9912/u-boot/contents/docs/user/benutzerhandbuch.md?ref=v0.11.0` | pass | Exit 0, Pfad vorhanden: Tag-Link nicht tot |

Plan-Tabelle §3 gegen `git diff --stat 2ab2ee6..HEAD` (§3.10): Diff berührt sieben Dateien — `root.go`, `cli_test.go`, `roadmap.md`, `CHANGELOG.md`, `benutzerhandbuch.md` (alle in der Tabelle als „update"), den Slice-Plan (Gegenstand) und den Review-Report (Tabellenzeile „neu, durch den Reviewer"). „Unverändert/prüfen" stimmt: `export_test.go` (`WalkRootCommandPathsForTest` ruft `buildRootCommand(a)`, Signatur unverändert), `releasing.md`, Lastenheft, Spezifikation, READMEs nicht im Diff. Keine Zeile ohne Plan-Eintrag, kein Eintrag ohne Diff.

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-002`](../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) | `TestExecute_RootHelpManualLinks`, `TestExecute_RootHelpPointsToUserDocs`, Praxislauf `--help` |
| [`LH-FA-CLI-003`](../../spec/lastenheft.md#lh-fa-cli-003--versionsausgabe) | Praxislauf `--version` unverändert; Version ist die Quelle des Tag-Links |
| [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | englischer Text, „(German)" benannt; Test auf `LH-`-Abwesenheit |
| [`LH-NFA-PERF-001`](../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort) | nicht gemessen (benannt in Plan §3 „Hinweise"); reine String-Logik, kein Netz |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: alle

Nicht ausgeführt:
- `make test-docker` — kein Template-, E2E- oder Port-Pfad im Diff.
- `make verify-depguard` — keine Importregel berührt (neu nur `fmt`, `regexp`).
- Mutationen — nicht wiederholt; Reviewer-Beleg übernommen, Fälle gelesen.
- Messung der 200-ms-Grenze — Docker-Lauf ungeeignet; im Plan benannt.

Commit / Artefakt:
- `8297111` (HEAD, Verifikations-Stand); Feature-Commit `ab931a2`.

## Diskrepanzen

Keine zwischen Zusage und Code. Der Slice-Plan führt alle DoD-Zeilen noch ungehakt; das ist erwartet, bis die Closure sie setzt.

## Auflagen für die Closure

1. §7 füllen: Closure-Notiz, Steering-Loop-Lerneintrag, Finding-Klassen aus dem Review (Kommentar-Zusage weicht von Quelle ab · Plan-Aussage nach Entscheidungswechsel nicht nachgezogen · Vorab-Tag-Aufräumen macht Tag-Link tot), Hinweis, dass der Cut keinen Handarbeitsschritt trägt.
2. §6: jedem der fünf Risiken einen Ausgang geben (F-3 gehört zum Risiko „Link zum Handbuch einer Version ohne Tag").
3. Register `BEO-CLI` und `BEO-DOK`: Sichtung nachholen und in §8 eintragen (Plan §8 verlangt sie vom Implementer; der Eintrag fehlt noch); sonst in §7 festhalten, dass keine Beobachtung anfiel.
4. DoD-Häkchen setzen, dann (§3.3 Fall 2) reiner `git mv` nach `done/`; Roadmap-Marker „In Arbeit" nach dem Move zurücksetzen und `docs/plan/planning/open/README.md` (Zeile verlinkt auf `in-progress/`) auf `done/` nachziehen.
5. Pfade im Plan so schreiben, dass sie vom Ruheort `done/` auflösen (Ruheort-Regel).
