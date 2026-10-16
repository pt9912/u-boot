# Review-Report: slice-v1-cli-quiet-vorschau — 2026-10-16

**Review-Art:** Code — Diff gegen Plan, [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.10.0), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) und Hard Rules (Maintainability; die DoD-Closure prüft der Verifier)

**Gegenstand:** Commit-Range `16a9603..HEAD` (Implementer-Commits `acefc1d` Code und Tests, `27c6de5` Handbuch und CHANGELOG) für den Slice [`slice-v1-cli-quiet-vorschau`](../plan/planning/in-progress/slice-v1-cli-quiet-vorschau.md) (Welle `welle-cli-verbosity`)

**Skill:** `.harness/skills/reviewer.md` @ HEAD `27c6de5` · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- Slice [`slice-v1-cli-quiet-vorschau`](../plan/planning/in-progress/slice-v1-cli-quiet-vorschau.md) (§1 Regeltabelle, §2 DoD, §3 Plan, §6 Risiken)
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- `AGENTS.md` (Hard Rules §3.2, §3.7), `harness/conventions.md`, `harness/review.md`; frühere Reports `docs/reviews/2026-10-16-slice-v1-cli-quiet-nur-fehler*.md` (Herkunft F-2) <!-- d-check:ignore (Dateimuster als Zitat) -->
- Sensors: `make test` sechsmal (Baseline grün, fünf Mutationen rot, Arbeitsbaum jeweils per `git checkout --` zurückgesetzt, `git status` sauber), `make docs-check` (Lauf nach Schreiben des Reports, Ergebnis siehe Verdikt). Nicht gelaufen: `make gates`/`make lint` (der Diff trägt kein `//nolint`; kein Neu-Lauf), `make test-docker` (kein Docker-Pfad), `make verify-depguard` (keine Importregel berührt)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der neue Doc-Kommentar von `summaryVisible` steht ohne Leerzeile unmittelbar hinter dem Doc-Kommentar von `previewModeFromFlags`; beide bilden einen zusammenhängenden Kommentarblock. Dadurch beginnt der Doc von `summaryVisible` mit der Wahrheitstabelle von `previewModeFromFlags` (samt Herkunftschronik „Originally lived in add.go"), und `previewModeFromFlags` trägt keinen eigenen Doc mehr. | `AGENTS.md` §3.7 (Ein Kommentar beschreibt, was da ist); Maintainability | `internal/adapter/driving/cli/previewmode.go` · „package helpers, no RunE binding)." gefolgt von „// summaryVisible reports whether the plan or result summary" | ja — `go doc`/Lesen der Datei; kein Gate fängt die Zuordnung | Einfügung trennt Kommentar von seinem Gegenstand |
| F-2 | INFO | Der CHANGELOG-Eintrag sagt, `--dry-run --quiet` gebe den Plan aus „(the same text as without `--quiet`)". Für `init` filtert der neue Code die Warnungen im Plan auf die Ausnahme-Warnungen (`resp.Warnings = visibleWarnings(resp.Warnings, flags.Quiet)`), der Text ist dann nicht byte-gleich zum Lauf ohne `--quiet`. Das folgt dem Lastenheft („Alle anderen Hinweise und Warnungen entfallen"); nur die Formulierung „same text" ist darauf nicht eingeschränkt. | `LH-FA-CLI-005`; Doku-Präzision | `CHANGELOG.md` · „prints the plan (the same text as without `--quiet`)" | Review-only | CHANGELOG-Formulierung ohne Warnungs-Vorbehalt | <!-- d-check:ignore (ID als Zitat) -->
| F-3 | INFO | Unter `--dry-run --quiet` erscheint für `remove`, `generate` und `config set` die Vergangenheitsform („Removed service …", „Generated …", „config: a → b."), obwohl nichts geschrieben wurde; der Test pinnt genau diese Texte als `summary`. Der Slice benennt das als offenes Risiko (§6); keine Abweichung vom Lastenheft. | `LH-FA-CLI-007` | `internal/adapter/driving/cli/quiet_preview_test.go` · `summary: `Removed service "postgres".`` | Review-only | Vergangenheitsform im Dry-Run (bekannt, §6) | <!-- d-check:ignore (ID als Zitat) -->
| F-4 | INFO | Die Warn-Konsistenz unter `--quiet --dry-run` ist nur für `init` durch einen Test gepinnt (`TestExecute_Init_Quiet_DryRunPlanKeepsOnlyExceptionWarnings`). Für `generate` (Filter `visibleWarnings(resp.Warnings, flags.Quiet)` nach dem Summary-Zweig) gibt es keinen Dry-Run-Pin mit Warnungen; die Ausnahme-Warnung bleibt dort nur durch den Bestandstest ohne `--dry-run` belegt. Das Verhalten selbst ist korrekt. | `LH-FA-CLI-005` | `internal/adapter/driving/cli/quiet_preview_test.go` · „TestExecute_Init_Quiet_DryRunPlanKeepsOnlyExceptionWarnings" | ja — `make test` (Mutation des `generate`-Filters im Dry-Run-Zweig wäre rot nur über Bestandstests) | Warn-Pin nur für ein Subkommando | <!-- d-check:ignore (ID als Zitat) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Schreibstellen `internal/adapter/driving/cli/*.go` gegen Regeltabelle Slice §1 | geprüft, ohne Befund: je `init`, `generate`, `add`, `remove`, `config set` und `runConfigSetMany` steht `writeDiff` (bei `flags.Diff`) vor dem Summary-Zweig; der Summary hängt an `summaryVisible(quiet, dryRun)`. (a) `--dry-run`: Summary bleibt. (b) `--dry-run --diff`: Diff und Summary bleiben (Reihenfolge Diff, dann Summary). (c) `--diff` ohne `--dry-run`: Diff bleibt, Erfolgsmeldung entfällt (`PreviewAndApply` hat `DryRun=false`). (d) weder noch: nichts. Die Wahrheitstabelle in `previewModeFromFlags` (DryRun gewinnt) stimmt mit `summaryVisible` überein. |
| Überkorrektur (Erfolgsmeldung nach echtem Schreiben) | geprüft, ohne Befund: `summaryVisible` ist nur bei `dryRun` wahr; `TestExecute_Quiet_ActualWriteStaysSilent` prüft je Subkommando `--quiet` ohne Option (stdout leer) und `--quiet --diff` (Diff da, Erfolgszeile des Default-Laufs nicht). Mutation M-c unten belegt die Wirksamkeit. |
| Warnungen/Hinweise im Dry-Run (Konsistenz) | geprüft, ohne Befund: `init` filtert den Plan auf Ausnahme-Warnungen (`visibleWarnings(…, flags.Quiet)`) und gibt ohne Dry-Run nur diese aus; `generate` gibt Summary und danach nur die Ausnahme-Warnungen aus (gleiche Semantik); `add`, `config set` und `config set` mit mehreren Paaren schreiben in Human-Mode keine Warnungen; `remove` erzeugt unter `PreviewDryRun` keine Purge-Warnungen (`previewMode != driving.PreviewDryRun` im Hinweiszweig, `PurgeAttempted` nur bei echtem Entfernen), und bei echtem Schreiben entfallen sie unter `--quiet` weiter. Kein Subkommando gibt eine nicht-ausgenommene Warnung unter `--quiet` aus. |
| Abgebrochene Dry-Runs (No-op-Plan, Kollision, `--no-interactive`) | geprüft, ohne Befund: Kollision und Abbruch kommen als Fehler aus dem Use Case und laufen vor dem Summary-Zweig durch `reportError`/`reportErrorSub` (Fehler bleibt sichtbar, Exit-Code aus der zentralen Klassifikation, unverändert); der No-op-Plan ist ein Erfolg und zeigt seine „no changes"-Zeile unter `--dry-run --quiet` (Pin `TestExecute_Add_Quiet_DryRunNoOpPlanStays`, für `add`; die übrigen Subkommandos nutzen denselben Zweig). Exit-Code und Ausgabe sind konsistent: `summaryVisible` wird nie auf einem Fehlerpfad ausgewertet. |
| `--quiet --json` byte-gleich, Exit-Codes, Dual-Classifier | geprüft, ohne Befund: alle `flags.JSON`-Zweige liegen vor `writeDiff` und Quiet-Zweig und sind unverändert; `TestExecute_Quiet_PreviewJSONUnchanged` vergleicht `--json --dry-run --diff` mit und ohne `--quiet` byte-gleich für alle sechs Fälle. Kein neuer Sentinel, kein Eingriff in `ExitCode` oder in eine Kommando-Diagnostic-Abbildung; `TestExecute_Quiet_PreviewErrorStays` pinnt `errors.Is` je Subkommando und Vorschau-Option. |
| Testmatrix gegen Verhalten | geprüft, ohne Befund: `TestExecute_Quiet_PreviewMatrix` prüft je sechs Subkommandos × vier Optionsreihen Summary-Text, Diff-Marker und leeres stderr gegen die Regeltabelle (nicht tautologisch: die Erwartung ist die Tabelle, die Ausgabe kommt aus der Produktionspfad-Verdrahtung mit Fake-Use-Cases) und vergleicht bei `--dry-run`-Zeilen die Ausgabe byte-gleich mit dem Lauf ohne `--quiet`. |
| Wächter-Mutationen (selbst gefahren, über `make test`, jeweils zurückgesetzt) | M-a `previewmode.go`: `return !quiet \|\| dryRun` → `return !quiet` → rot. M-b `init.go`: Zeile `resp.Warnings = visibleWarnings(resp.Warnings, flags.Quiet)` entfernt → rot. M-c `generate.go`: `summaryVisible(flags.Quiet, flags.DryRun \|\| flags.Diff)` (Überkorrektur bei `--diff` ohne `--dry-run`) → rot. M-d `config.go` (`runConfigSetMany`): `!summaryVisible(…)` → `flags.Quiet` (Dry-Run-Plan unterdrückt) → rot. M-e `init.go`: `printInitSummary(out, resp, false)` (Dry-Run als „Initialized") → rot. Jeweils Exit 2 des Test-Stage; die Fehlerliste der letzten Mutation nennt `TestInitHumanDryRun_ShowsWouldInitiate`, `TestExecute_Init_Quiet_DryRunPlanKeepsOnlyExceptionWarnings` und die Matrix-Fälle `init/dry-run`, `init/dry-run_diff`. |
| Auslegung `--diff` ohne `--dry-run` (Slice §6) gegen Lastenheft-Wortlaut | geprüft, ohne Befund: `spec/lastenheft.md` · „Die Erfolgs- und Zusammenfassungsmeldung, die ein tatsächliches Schreiben abschließt, ist keine Vorschau und entfällt auch bei `--diff` ohne `--dry-run`" deckt Slice §1 Zeile 3 und den Code (Diff bleibt, Erfolgsmeldung entfällt) wörtlich. Die Slice-Aussage „vom Projektinhaber zu bestätigen" ist damit durch 0.10.0 überholt; das Risiko-Feld §6 trägt noch „offen" (Aufgabe der Closure, nicht des Diffs). |
| Handbuch und CHANGELOG gegen Lastenheft 0.10.0 | geprüft, ohne Befund über F-2 hinaus: beide Optionstabellen-Zeilen des Handbuchs (Globale Optionen und Option-Tabelle am Ende) sind identisch ergänzt um Vorschau als Ergebnis des Aufrufs und die entfallende Erfolgsmeldung bei `--diff` ohne `--dry-run`; der CHANGELOG nennt nur noch „the last verbosity option wins" als folgend. Keine weitere `--quiet`-Stelle in `docs/user/`, `README.md`, `README.de.md` (Grep; Treffer `--quiet ändert die JSON-Ausgabe nicht` ist unverändert zutreffend). |
| Architekturgrenzen, Suppression, Sprache | geprüft, ohne Befund: Änderungen nur im Driving-Adapter, keine neuen Imports (`previewmode.go` importiert weiter nur `port/driving`); kein `//nolint` im Diff; die Texte der Ausgabe sind unverändert Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)); der neue Helfer ist rein und ohne I/O. |
| §3.7-Kommentare im neuen Code | geprüft, ohne Befund über F-1 hinaus: die Kommentare in `add.go`, `config.go`, `generate.go`, `init.go`, `remove.go` stehen im Indikativ, ohne Chronik und ohne tote Symbolverweise („Diff", „plan of a dry run", „four exceptions" sind Vertragsbegriffe). |
| Referenzmodell, Lifecycle, Carveouts, Modus | geprüft, ohne Befund: keine Spec-/ADR-Dateien im Diff; der Slice liegt in `in-progress/` (reiner Move-Commit `591e457` vor den Code-Commits, [`AGENTS.md`](../../AGENTS.md) §3.3); Sub-Area CLI (Greenfield) trägt den Diff, keine Konvention entsteht nur im Code. |
| `make gates`, `make test-docker`, `make verify-depguard` | nicht ausgeführt: siehe Eingangs-Kontext. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Einfügung trennt Kommentar von seinem Gegenstand · CHANGELOG-Formulierung ohne Warnungs-Vorbehalt · Vergangenheitsform im Dry-Run (bekannt, §6) · Warn-Pin nur für ein Subkommando

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Das Verhalten folgt Lastenheft 0.10.0 und der Regeltabelle des Slice in allen sechs Schreibpfaden; die Mutationen M-a bis M-e sind rot. F-1 ist ein Kommentar-Zuordnungsfehler ohne Verhaltenswirkung (LOW, Anker `AGENTS.md` §3.7), F-2 bis F-4 sind Beobachtungen ohne Änderungspflicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Die DoD-Prüfung (Gates, Closure-Notiz, Risiko-Ausgänge in §6) liegt beim Verifier.
