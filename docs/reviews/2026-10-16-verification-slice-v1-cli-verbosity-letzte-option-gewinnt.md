# Verification Evidence: slice-v1-cli-verbosity-letzte-option-gewinnt — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../plan/planning/done/slice-v1-cli-verbosity-letzte-option-gewinnt.md), letzter Slice der Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md). Commits: `65b1fa6` (Code und Tests), `c56d134` (Handbuch, CHANGELOG), `0bfe1fa` (`SPEC-043`), `4260372` (Review-Fixes F-1/F-2 und Report). Review: [Report](2026-10-16-slice-v1-cli-verbosity-letzte-option-gewinnt.md) (1 HIGH, 0 MEDIUM, 1 LOW, 2 INFO; HIGH und LOW behoben). <!-- d-check:ignore (Commit-Hashes und SPEC-Kennung als Zitat) -->

## Verification Evidence

Scope:
- Slice: [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../plan/planning/done/slice-v1-cli-verbosity-letzte-option-gewinnt.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.10.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- Artefakte: `internal/adapter/driving/cli/` (neu `verbosity.go`, `verbosity_last_wins_test.go`; geändert `root.go`, `cli.go`, `verbosity_test.go` und elf Lesestellen), `cmd/uboot/main.go` (Kommentar), `spec/spezifikation.md` (`SPEC-043`), `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`. <!-- d-check:ignore (Pfade und Kennung als Zitat) -->

DoD-Abgleich (Slice §2; die Häkchen im Slice sind noch ungesetzt, sie folgen mit der Closure-Notiz):
- [x] `LH-FA-CLI-005` erfüllt, Happy/Boundary/Negative: `verbosity_test.go` (`TestExecute_VerbosityFlag_AdjustsLoggerLevel`, 12 neue Zeilen: alle vier Paarordnungen, Wiederholung `--quiet --verbose --quiet` und `--verbose --quiet --verbose`, gleiche Option doppelt, Option nach dem Subkommando, `--flag=false`/`=true`) und `verbosity_last_wins_test.go` (Lesestellen: `init`-Zusammenfassung, Warnungsfilter, `doctor`-Bericht, `down`-Compose-Strom; `--json` unverändert; Zustand je `Execute`; ungültiger Wert abgewiesen). Der Pin `debug-takes-precedence-over-quiet` ist zu `quiet-after-debug-wins` umgekehrt, nichts gelockert. Das Slice-Risiko „`a.quiet` driftet vom Pegel“ ist strukturell erledigt: ein Zustand `verbosity`, `a.quiet()` leitet daraus ab; alle zwölf Lesestellen (`grep`) nutzen `a.quiet()`. <!-- d-check:ignore (Testnamen und Pfade als Zitat) -->
- [x] Doku: Handbuch (beide Optionstabellen, gleiche Texte, `--verbose`/`--debug`-Zeilen und Alias-Beispiel) und CHANGELOG `[Unreleased]` nennen „letzte Option gewinnt“.
- [x] Spezifikation: `SPEC-043` (siehe unten, Aussage für Aussage geprüft), an `LH-FA-CLI-005` verankert, Werte aus dem Code. <!-- d-check:ignore (Kennungen als Zitat) -->
- [x] `make gates` grün — selbst gefahren, Exit 0 (Coverage 91,10 %, Schwelle 90 %; d-check 269 Dateien, 0 Befunde; baseline-verify v6.18.1 OK, 54 Dateien).
- [x] Review durchgeführt — Report vorhanden; F-1 (HIGH) und F-2 (LOW) in `4260372` behoben, F-3/F-4 INFO offen (Bewertung unten).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen (§7 „Wird bei der Closure gefüllt“; Auflage S1).
- [ ] Beobachtungs-Register — offen (Auflage S1; `observations/` enthält nur die README).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: Risiko 1 steht auf „offen bis zur Umsetzung“ (Auflage S2); Risiko 2 trägt „entfallen“, zutreffend.
- [ ] Drei Paarungen — offen bis zur Closure (Auflage S3).

### Prüfung 1: `SPEC-043` gegen Code und Lastenheft 0.10.0

- Pegelabbildung (`--quiet` Error, keine Option Info, `--verbose` und `--debug` Debug): deckungsgleich mit `verbosity.slogLevel()`; gesetzt im `PersistentPreRunE` nach dem Parsen (`root.go`). „Keine Quelle außerhalb der drei Optionen“: bestätigt durch Review-Negativbefund (kein Env-/Config-Pfad) und Slice §1.
- Rangfolge: ein gemeinsamer Zustand `App.verbosity`, jede Angabe überschreibt (`verbosityFlag.Set`); pflag ruft `Set` je Vorkommen in Parse-Reihenfolge, auch hinter dem Subkommando (Test `option-after-subcommand-counts`). „Nicht abgewiesen, Exit-Code unverändert“: Lastenheft Zeile 198 wörtlich. `a.quiet()` nur bei `verbosityQuiet`: so formuliert und durch `LastWins_*`-Tests gepinnt.
- Wirkung auf die Ausgabe: `writeDoctorReport` (Error-Items und `DegradationFallback`, kein Kopf, keine Summary), `visibleWarnings` (Sicherheitslockerung oder Degradations-Fallback), `summaryVisible(quiet, dryRun)` (Vorschau bleibt, Erfolgsmeldung eines Schreibens entfällt auch bei `--diff`), `down` puffert die Compose-Ausgabe außer bei `--json` (`down.go`), `remove` ohne Purge-Warnungen unter `--quiet`: alles mit dem Code identisch und mit Lastenheft 0.10.0 (Zeilen 189 bis 198) deckungsgleich, vier Ausnahmen in gleicher Reihenfolge und gleichem Wortlaut. „`init --force`“ als einziger Träger der Pfad-Zusammenfassung stimmt mit dem Review-Negativbefund.
- Referenzrichtung: Der Abschnitt enthält nur Verweise auf `LH-*` (Grep auf Slice, ADR, Welle, Tranche, Carveout, CR: kein Treffer); AGENTS.md §3.4 und die Decken-Regel gewahrt.
- Historie-Zeile (2026-10-16, „Verbosity festgelegt“) vorhanden; Nummerierung: `SPEC-043` steht direkt nach `SPEC-042`, die Lücken `SPEC-011` und `SPEC-037` bleiben (Grep der Überschriften). <!-- d-check:ignore (Kennungen als Zitat) -->
- Offene Randfälle: siehe F-3 unten (nicht in der Spec).

### Prüfung 2: Welle-Closure-Trigger (Welle §3) und `doc-trace`

- Die Welle führt **sechs** Slices (Tabelle §4; der Auftrag nannte sieben): fünf liegen in `done/` (`slice-cr-cli-quiet-fehler-und-warnungen`, `slice-cr-cli-quiet-weitere-ausnahmen`, `slice-v1-cli-quiet-nur-fehler`, `slice-cr-cli-quiet-vorschau`, `slice-v1-cli-quiet-vorschau`; per `find` bestätigt), der sechste ist dieser Slice in `in-progress/`. Die Welle-Closure-Bedingung „alle in `done/`“ ist erst nach dem `git mv` dieses Slice erfüllt; Lastenheft steht auf 0.10.0. <!-- d-check:ignore (Slice-Kennungen als Zitat, Links stehen im Slice) -->
- Commit-Reihenfolge: Die Spezifikation (`0bfe1fa`) entstand nach den Schließungen aller fünf anderen Slices (`slice-mv` nach `done/` bis `ea3b020`/`668725c`, danach Start `2b2cec0`, Code `65b1fa6`, Doku `c56d134`, Spec `0bfe1fa`): Reihenfolge belegt, Risiko 2 zu Recht „entfallen“. <!-- d-check:ignore (Kennungen als Zitat) -->
- `make doc-trace`: selbst gefahren, „83 Anforderung(en), 0 Waise(n).“; `LH-FA-CLI-005` ist `ok` und trägt diesen Slice. Hinweis: dieselbe Zeile trägt `slice-v1-sandbox-fallback-ausweisung` (Folge-Slice in `open/`, Anker `LH-FA-CLI-005`). <!-- d-check:ignore (Kennungen als Zitat) -->

### Prüfung 3: CHANGELOG `[Unreleased]` (Release-Cut-Szenario)

Der Eintrag nennt Lastenheft 0.10.0, die vier Ausnahmen, `--quiet`-Reichweite (ändernde Subkommandos, `doctor`, `down`-Puffer, Abfragebefehle unverändert, `--json` No-op), Vorschau von `--dry-run`/`--diff`, `doctor --strict` mit Warn-Items und die letzte Option gewinnt samt Wegfall der festen Reihenfolge und `--debug --quiet` als Verhaltensänderung. Nichts steht mehr als „folgt“ (Grep `later release`: kein Treffer). Nicht genannt ist Lastenheft 0.9.0 (der Eintrag zitiert nur 0.10.0, das 0.9.0 mit einschließt; kein Befund). Die Dichte eines einzigen Aufzählungspunkts ist lang, aber inhaltlich vollständig. Der Release-Cut ist durch den CHANGELOG nicht mehr blockiert (die Auflage 4 der Evidence zu `slice-v1-cli-quiet-vorschau` ist erledigt). <!-- d-check:ignore (Kennung und Pfad als Zitat) -->

### Prüfung 4: Folge-Arbeit und Beobachtungen

- `slice-v1-sandbox-fallback-ausweisung` und `slice-v1-cli-dry-run-would-formulierung` liegen in `open/`, jeweils „ohne Welle“, mit eigener Herkunft (nicht aus diesem Slice): Sie sind Folge-Slices der **Welle**, nicht dieses Slice. Dieser Slice hat keinen Befund, der sie begründete; §7 dieses Slice nennt sie nur als „Welle-Folge, unberührt“. Die Welle-Closure muss sie in Welle §7 als Folge-Slices tragen (Paarung Folge-Slice). <!-- d-check:ignore (Kennungen als Zitat) -->
- Kommentar-Befunde nach `AGENTS.md` §3.7 in dieser Welle: F-1 des Slice `slice-v1-cli-quiet-nur-fehler`, F-1 und F-2 hier, F-1 des Slice `slice-v1-cli-quiet-vorschau` = **vier Auftreten**, Schwelle 3× überschritten; `observations/` ist leer. Auflage W2. <!-- d-check:ignore (Kennungen als Zitat) -->

### Bewertung der offenen Review-Findings

- **F-3 (`--flag=false` und `--quiet=maybe` nur im Code):** nicht blockierend. Das Lastenheft sagt „keine Validierungsabweisung wegen Mehrfachangabe“; ein Nicht-Boolean-Wert ist keine Mehrfachangabe, und der Parse-Fehler gehört zur Klasse Exit `2` nach `LH-FA-CLI-006`. Beide Randfälle sind im Test gepinnt (`explicit-false-*`, `InvalidValueIsRejected`; letzterer nur als `err != nil`, nicht als Exit-Code). Genügt als benannte Grenze in §7 („Randfälle `--flag=false` und ungültiger Wert nur im Code und Test, nicht in `SPEC-043`/Handbuch“); alternativ ein Satz in `SPEC-043`, dann aber Änderung nach Spec-Regeln, kein Muss. <!-- d-check:ignore (Kennung als Zitat) -->
- **F-4 (`--quiet doctor --verbose` nur im Logger-Test):** nicht blockierend. Ein gemeinsamer Zustand schließt Drift aus; der Pegel-Pin trägt die Reihenfolge-Eigenschaft, mein Mutationslauf M-1 (unten) stützt das. Als Grenze in §7 benennen.

Plan-vs-Code-Diff (Slice §3):
- `root.go` — Ist: `BoolVar` ersetzt durch `addVerbosityFlag` (`VarPF`, `NoOptDefVal`), `PersistentPreRunE` ruft `slogLevel()`, Kommentar ersetzt. Wie geplant.
- `cli.go` — Ist: drei Felder durch `verbosity` ersetzt; Kommentar beschreibt den geltenden Stand.
- `verbosity_test.go` — Ist: Pin umgekehrt, neue Fälle. Wie geplant.
- Lesestellen — Plan: „konsistent halten“; Ist: elf Dateien lesen jetzt `a.quiet()` statt des Feldes (mechanisch, notwendig wegen des Wegfalls des Feldes). `doctor_json_test.go`, `up_acceptance_test.go`, `logs_acceptance_test.go` unverändert, wie „prüfen“ vorsah.
- Abweichung ohne Plan-Zeile: neue Dateien `verbosity.go` und `verbosity_last_wins_test.go`, `cmd/uboot/main.go` (Kommentar). Kleine Abweichung; der Plan §3 ist bei der Closure um diese Zeilen zu berichtigen.
- Handbuch, CHANGELOG — Ist: gesetzt (Zeilenangabe „um 299 und 1280“ im Plan; beide Tabellen geändert). Kein `//nolint`, keine Gate-Lockerung, keine neuen Imports außer `strconv`/`slog` im Adapter.

Bewertung der Risiken aus §6:
- **`a.quiet` driftet vom Pegel:** tatsächlich **entfallen** (ein Zustand, keine zwei Felder), belegt durch `TestExecute_Verbosity_LastWins_*` für Zusammenfassung, Warnungen, `doctor` und `down`-Strom. In §6 so eintragen (Auflage S2).
- **Reihenfolge-Abhängigkeit:** „entfallen“ korrekt (Prüfung 2).

Mutationen (selbst gefahren, nur über `make test`, danach `git checkout -- <Datei>`; `git status` sauber):
| Nr | Mutation | Ergebnis |
| --- | --- | --- |
| M-1 | `verbosity.go`: `Set` überschreibt nur bei Default (erste Option gewinnt) | rot: 22 `--- FAIL`-Zeilen, u. a. `LastWins_InitSummary`, `LastWins_WarningFilter`, `LastWins_DoctorReport`, `LastWins_DownComposeStream` |
| M-2 | `verbosity.go`: `quiet()` wahr bei jedem Nicht-Default (`!= verbosityDefault`, Lesestellen driften vom Pegel) | rot: 12 `--- FAIL`-Zeilen, vier `LastWins_*`-Tests (jeweils die `--verbose`/`--debug`-Fälle) |

Der Review-Report nennt drei weitere Mutationen; sie wurden hier nicht wiederholt.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (lint, test, coverage-gate 91,10 %, docs-check 0 Befunde, baseline-verify) |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“; `LH-FA-CLI-005` ok, trägt diesen Slice | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make test` (2 Mutationen) | rot wie erwartet | Tabelle oben |
| `make test-docker` | not run | kein Docker-Pfad berührt: Driving-Adapter-Flags, Fake-Use-Case-Tests, Doku |
| `make verify-depguard` | not run | keine Importregel geändert |
| `make docs-check` | pass | nach Anlegen dieser Datei erneut gefahren, siehe Übergabe |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `verbosity_test.go`, `verbosity_last_wins_test.go`, `SPEC-043`, Handbuch, CHANGELOG; `doc-trace` ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | kein Exit-Code verändert; Parse-Fehler bei ungültigem Wert bleibt Klasse Exit `2` |

Carveouts: neu none · gelöst none · unverändert none berührt.

Commit / Artefakt: Slice-Commits `65b1fa6`, `c56d134`, `0bfe1fa`, `4260372`; HEAD zum Prüfzeitpunkt `4260372`. <!-- d-check:ignore (Commit-Hashes als Zitat) -->

## Verdikt

**Bestanden, mit Auflagen für die Closure.** Code, Tests und `SPEC-043` decken Lastenheft 0.10.0 deckungsgleich ab; Gates grün, Waisen null, zwei Mutationen fallen rot, keine Verletzung von Hard Rules oder Plan-Zusage. F-3 und F-4 sind nicht blockierend. <!-- d-check:ignore (Kennung als Zitat) -->

Auflagen für die Slice-Closure:
- S1. Closure-Notiz §7 mit Steering-Loop-Lerneintrag (Klasse: Kommentar beschreibt abwesenden Text bzw. falsche Stufe, Review F-1/F-2); Register: Beobachtung eintragen oder „keine Beobachtung“ notieren (siehe W2 für den Wellen-Befund); Grenzen F-3 und F-4 in §7; DoD-Häkchen; dann reiner `git mv` (`AGENTS.md` §3.3). <!-- d-check:ignore (Pfad als Zitat) -->
- S2. §6: Risiko „`a.quiet` driftet“ auf „entfallen: ein Zustand, `a.quiet()` abgeleitet, Tests `LastWins_*`“; §3 um `verbosity.go`, `verbosity_last_wins_test.go`, `cmd/uboot/main.go` ergänzen. <!-- d-check:ignore (Dateinamen als Zitat) -->
- S3. Paarungen: `open/README.md` §Bestand führt die Zeile dieses Slice noch (Link zeigt auf `in-progress/`); beim `git mv` bereinigen. <!-- d-check:ignore (Pfad als Zitat) -->

Auflagen für die Welle-Closure:
- W1. Welle §3 verlangt „alle sechs Slices in `done/`“: erst nach dem `git mv` dieses Slice erfüllbar; danach Welle §7 füllen, Welle-Results nach der Vorlage, `roadmap.md` nachführen.
- W2. Beobachtungs-Register: vier Kommentar-Befunde nach `AGENTS.md` §3.7 in einer Welle (3×-Schwelle überschritten). Anzulegen ist eine Beobachtung `BEO-CLI/<slug>` (Kürzel CLI, vier Evidence-Dateien, je Review-Finding), danach Entscheidung bei der Welle-Closure, ob die bestehende Regel §3.7 samt Reviewer-HIGH-Liste genügt (dann Zielort `liegt in` samt Anker `seit welle-cli-verbosity`) oder ein Sensor nötig ist. Dieser Lauf legt nichts an. <!-- d-check:ignore (Pfad-Schema und Kennung als Zitat) -->
- W3. Folge-Slices der Welle in Welle §7 nennen: `slice-v1-sandbox-fallback-ausweisung` und `slice-v1-cli-dry-run-would-formulierung` (beide `open/`, ohne Welle). <!-- d-check:ignore (Kennungen als Zitat) -->
- W4. Release-Cut: CHANGELOG ist vollständig (Prüfung 3); Lastenheft-Version 0.10.0 und Handbuch passen; der Cut ist durch diese Welle nicht mehr blockiert.
