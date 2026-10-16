# Review-Report: slice-v1-cli-verbosity-letzte-option-gewinnt — 2026-10-16

**Review-Art:** Code — gegen Plan, Lastenheft 0.10.0, `SPEC-043` und Hard Rules. <!-- d-check:ignore (Kennung als Zitat) -->

**Gegenstand:** Commit-Range `321f365..0bfe1fa` (`65b1fa6` Code und Tests, `c56d134` Handbuch und CHANGELOG, `0bfe1fa` `SPEC-043`); Slice `slice-v1-cli-verbosity-letzte-option-gewinnt` (Welle `welle-cli-verbosity`, letzter Slice). <!-- d-check:ignore (Kennung als Zitat) -->

**Skill:** `.harness/skills/reviewer.md` @ `cf2f3e7` · <!-- d-check:ignore (Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-cli-verbosity-letzte-option-gewinnt` (Scope, DoD, Nicht-Ziele) <!-- d-check:ignore (Kennung als Zitat) -->
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.10.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)
- [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet)
- `AGENTS.md` (Hard Rules §3.2, §3.4, §3.6, §3.7)
- Sensor-Läufe des Reviewers: `make test` (nur für drei Mutationen, siehe Negativbefunde); `make docs-check` am Ende (siehe Verdikt). `make test-docker` und `make verify-depguard` nicht gelaufen (siehe „Nicht geprüft").

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Testkommentar beschreibt abwesenden Text: die „former fixed order" sei weg. Er trägt keine der Kommentar-Klassen (Zusage, Kopplung, Abgrenzung, Rang-Zeiger, Grenze); die geltende Zusage („the option given last decides") steht schon im Satz davor. | `AGENTS.md` §3.7 (Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen") | `internal/adapter/driving/cli/verbosity_test.go` · "the former fixed order (--debug before --quiet) is gone." | Review-only | Kommentar beschreibt abwesenden Text |
| F-2 | LOW | Der Wiring-Kommentar nennt für `--quiet` den Pegel Warn; Code und `SPEC-043` setzen `slog.LevelError`. Der Slice-Plan nannte die Fundstelle nur als cli.go-Kommentar; die Aussage in `main.go` blieb stehen. | `AGENTS.md` §3.7; `SPEC-043` (Pegelabbildung) | `cmd/uboot/main.go` · "Debug, --quiet lowers to Warn." | Review-only | Kommentar widerspricht Code | <!-- d-check:ignore (Kennung als Zitat) -->
| F-3 | INFO | `SPEC-043` und Handbuch legen die Wirkung von `--flag=false` nicht fest (Code: `--quiet --verbose=false` bleibt `--quiet`; `--quiet --quiet=false` kehrt zum Default zurück) und nicht den Nicht-Boolean-Wert (`--quiet=maybe` ist ein Parse-Fehler der Klasse „invalid argument", Exit `2` nach `LH-FA-CLI-006`, pinnt der Test `TestExecute_Verbosity_InvalidValueIsRejected` nur als „err != nil", nicht als Exit-Code). Die Lastenheft-Zusage „keine Validierungsabweisung wegen Mehrfachangabe" ist davon nicht berührt (ungültiger Wert ist keine Mehrfachangabe). | `LH-FA-CLI-005`, `LH-FA-CLI-006` | `internal/adapter/driving/cli/verbosity_last_wins_test.go` · "--quiet=maybe was accepted" | Review-only | Randfall nur im Code festgelegt | <!-- d-check:ignore (Kennung als Zitat) -->
| F-4 | INFO | Die Reihenfolge „Option nach dem Subkommando" (`--quiet doctor --verbose`) ist nur im Logger-Pegel-Test gepinnt, nicht in den Lesestellen-Tests (Zusammenfassung, Warnungen, `doctor`, `down`-Puffer). Mechanismus ist ein gemeinsamer Zustand, das Risiko ist gering. | `LH-FA-CLI-005` | `internal/adapter/driving/cli/verbosity_test.go` · "option-after-subcommand-counts" | `make test` | Lesestellen-Pin ohne Position-nach-Subkommando | <!-- d-check:ignore (Kennung als Zitat) -->

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Mechanismus `verbosity.go` (`verbosityFlag`, `Set`, `String`, `Type`, `NoOptDefVal`) | geprüft, ohne Befund: die drei Flags teilen ein `*Flag`-Wertobjekt je Name über den Persistent-Flagset; pflag/Cobra rufen `Set` je Vorkommen in Parse-Reihenfolge auf, auch für Optionen nach dem Subkommando. Fälle `--quiet --verbose`, `--verbose --quiet`, `--debug --quiet`, `--quiet --debug`, `--quiet --verbose --quiet`, `--quiet=false`, `--verbose=false` nach `--quiet` sind in Tests mit den im Code ableitbaren Werten gepinnt. `Type()` ist `bool`, `NoOptDefVal` `true`. |
| Zustand zwischen `Execute`-Läufen, Parallelität | geprüft, ohne Befund: `buildRootCommand` setzt `a.verbosity = verbosityDefault`; Zustand je `App`, kein Package-Global; Test `TestExecute_Verbosity_StateIsPerExecute` pinnt das. Parallele Tests nutzen je eigene `App`. |
| Kollision/Shadowing | geprüft, ohne Befund: kein Subkommando definiert eigene `quiet`/`verbose`/`debug`-Flags; kein zweiter `PersistentPreRunE`; keine `Changed`-Abfragen; keine Umgebungsvariable als Pegelquelle (bestätigt `SPEC-043`). | <!-- d-check:ignore (Kennung als Zitat) -->
| Lesestellen | geprüft, ohne Befund: jede Schreibstelle (`init`, `generate`, `add`, `remove`, `config set/get/list/show`, `up`, `down`, `doctor`, `logs`) setzt `flags.Quiet` aus `a.quiet()`; `summaryVisible`, `visibleWarnings`, `renderUpDiagnostics`, `renderDownSuccess`, `writeDoctorReport`, `down`-Puffer lesen nur `flags.Quiet`. Kein Pfad leitet `--quiet` aus anderem Zustand ab. |
| Pegel/Logger | geprüft, ohne Befund: `slogLevel()` Error/Debug/Debug/Info; `--debug --quiet` ergibt Pegel Error und `quiet()` wahr aus demselben Zustand, konsistent. |
| Exit-Code / Dual-Classifier (Skill HIGH #1) | geprüft, ohne Befund: kein neuer Sentinel; ein Nicht-Boolean-Wert bleibt ein Cobra-Parse-Fehler (Präfix „invalid argument", Klasse Exit `2`). Exit-Codes sonst unverändert. |
| Wächterqualität | geprüft, ohne Befund: Tests prüfen Ausgabe-Verhalten (Zusammenfassung, Warnungsfilter, `doctor`-Bericht, `down`-Strom, `--json` unverändert). Der geänderte Pin (`quiet-after-debug-wins`) setzt den Vertrag um, nichts gelockert. Drei eigene Mutationen, alle rot: (1) `Set` ignoriert `--flag=false` (Zweig `f.app.verbosity == f.mode` deaktiviert) → `TestExecute_Verbosity_LastWins_InitSummary` und Logger-Test rot; (2) `doctor.go` liest `a.verbosity != verbosityDefault` statt `a.quiet()` → `…LastWins_DoctorReport` plus mehrere bestehende `doctor`-Quiet-Tests rot; (3) `NoOptDefVal` entfernt → rot (Doctor-/Config-Quiet-Tests). Alle Mutationen zurückgesetzt (`git status` sauber). |
| `SPEC-043` | geprüft, ohne Befund bis auf F-3: Pegelabbildung, Rangfolge, vier Ausnahmen (Pfad-Zusammenfassung nur bei `init --force`, dem einzigen `--force`), Abfragebefehle, `--dry-run`/`--diff`, `down`-Puffer, `--json`-Wirkung stimmen mit Code und Lastenheft 0.10.0 überein; nur Aufwärtsverweise (`LH-*`), keine Slice-/ADR-/Roadmap-Verweise, keine Welle; Historie-Zeile vorhanden. | <!-- d-check:ignore (Kennung als Zitat) -->
| Handbuch (beide Optionstabellen), CHANGELOG, Hilfetexte | geprüft, ohne Befund: beide Tabellen gleich; Alias-Beispiel (`alias ub='u-boot --quiet'`, `ub --verbose`) stimmt mit dem Verhalten; CHANGELOG nennt nichts mehr als folgend, nennt die Verhaltensänderung `--debug --quiet`; Hilfetexte Englisch. |
| Architektur, Suppression, Sprache | geprüft, ohne Befund: keine neuen Imports außerhalb `slog`/`strconv`/cobra; kein `//nolint`; Doc-Kommentare stehen direkt über ihren Funktionen, keine toten Symbolverweise (`[verbosityFlag]`, `[verbosity.slogLevel]`, `[App.quiet]` existieren). |

Nicht geprüft:

- `make test-docker` — kein Docker-Pfad berührt (nur Driving-Adapter-Flags).
- `make verify-depguard` — keine Importregeln geändert.
- `make gates` im Ganzen — Verifier-Aufgabe; der Reviewer fuhr nur `make test` (Mutationen) und `make docs-check`.
- Hilfeausgabe von `--quiet` (Platzhalter durch Backticks im Usage-Text) — nicht ausgeführt; kein Host-Go, kein Test dazu.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar beschreibt abwesenden Text · Kommentar widerspricht Code · Randfall nur im Code festgelegt · Lesestellen-Pin ohne Position-nach-Subkommando

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (HIGH nach Skill-Anker, Failure-Szenario: der nächste Änderer der Pegelfälle liest „former fixed order … is gone" als Hinweis auf einen noch zu entfernenden Rest und sucht oder reaktiviert eine Rangfolge, die kein Gate widerlegt). F-2 bis F-4 blockieren nicht.

**Übergabe:** Findings an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation.
