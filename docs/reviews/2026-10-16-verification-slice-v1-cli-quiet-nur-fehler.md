# Verification Evidence: slice-v1-cli-quiet-nur-fehler — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/done/slice-v1-cli-quiet-nur-fehler.md), Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md). Commits: `c277c47`, `72006d3`, `7a244b9` (Erstlauf); `4c68a86`, `1e35794`, `82af811` (Ausbau); `b31451e`, `2b16cf2` (Review-Fixes). Reviews: [Erstlauf](2026-10-16-slice-v1-cli-quiet-nur-fehler.md) (F-1 bis F-7), [zweiter Lauf](2026-10-16-slice-v1-cli-quiet-nur-fehler-2.md) (0 HIGH, 3 MEDIUM, 3 INFO). Entscheidungen des Projektinhabers vom 2026-10-16 sind als gegeben genommen: (1) `down --quiet` puffert, (2) Vorschau-Frage geht in die Folge-Slices, (3) Ausnahme 2 ist im Betrieb nur für `vfs` im `doctor` erreichbar.

## Verification Evidence

Scope:
- Slice: [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/done/slice-v1-cli-quiet-nur-fehler.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Lastenheft 0.9.0), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation), [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)
- Artefakte: `internal/adapter/driving/cli/` (`root.go`, `cli.go`, `init.go`, `generate.go`, `add.go`, `remove.go`, `up.go`, `down.go`, `config.go`, `logs.go`, `doctor.go`, `warnings.go`, neue Tests `quiet_test.go`, `quiet_modifying_test.go`), `internal/hexagon/{domain,port/driving,application}`, `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/cli-json-contract.md`. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2; die Häkchen im Slice sind noch ungesetzt, sie folgen mit der Closure-Notiz):
- [x] `LH-FA-CLI-005` (0.9.0) für `init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `doctor`, Abfragebefehle — Evidence: `quiet_modifying_test.go`, `quiet_test.go`, `cli_test.go`; Mutation M-2 und M-3 (unten) fallen rot. `remove --quiet` existiert (Review-2 Negativbefund, `remove.go`). <!-- d-check:ignore (Kennung und Pfade als Zitat) -->
- [x] Purge- und Abhängigkeits-Warnungen von `remove` entfallen, ausdrücklich gepinnt — Evidence: `…Remove_Quiet_…PurgeWarningsAreDropped`, `…PurgeRefusedStaysAnError` (Review-2: Mutation M-c rot). <!-- d-check:ignore (Testnamen als Zitat) -->
- [x] Vier Ausnahmen je einzeln gepinnt, mit zwei benannten Grenzen — Ausnahme 1: Test pinnt Podman/Docker/Egress aus `sandboxRelaxationWarnings` für `init`, `generate`, `up`. Ausnahme 2: `TestExecute_Quiet_KeepsEachDegradationFallback` (neun Teilfälle, nur über Fake-Use-Cases) und `…Doctor_Quiet_FallbackItemAndErrorKept_OtherWarnsDropped` (echtes Item `vfs`; M-2 rot). Ausnahme 3: `up`-Stream auf stderr. Ausnahme 4: `…Init_Quiet_DoesNotSilenceForcePathSummary` prüft nur `SilenceProgress` im Request, nicht die tatsächliche Ausgabe (schwacher Pin, siehe Befund V-3). <!-- d-check:ignore (Testnamen als Zitat) -->
- [x] Review-Befunde behoben (F-1 Kommentar, F-3 `--strict` mit Exit 11 und Ursache auf stderr, F-4 Verweis `visibleWarnings`) — Evidence: Review 2 „Stand der Findings des Erstlaufs", selbst gegengelesen (`warnings.go` existiert, Verweis löst auf).
- [x] Doku: Handbuch an beiden Stellen und CHANGELOG — Evidence: Diff `2b16cf2`; beide Optionstabellen nennen „nur Fehler", die vier Ausnahmen, `down`-Puffer, `doctor`, Abfragebefehle, `--json`-No-op. Ehrlich bei der Lücke („soweit vom Host aus ermittelbar").
- [x] `make gates` grün — selbst gefahren, Exit 0, Coverage 91,00 % (Schwelle 90 %).
- [x] Review durchgeführt — zwei Reports vorhanden.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: §7 trägt „Wird bei der Closure gefüllt." (Auflage 1).
- [ ] Beobachtungs-Register — offen: `observations/` leer; Fortschreibung oder „keine Beobachtung" in §7 (Auflage 1, 4).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: die vier Risiken tragen „entfallen", aber zwei neue sind nicht erfasst (Auflage 2).
- [ ] Drei Paarungen (Anker · Folge-Slice · Register) — offen bis zur Closure; Folge-Slices für Vorschau und Fallback-Ausweisung (Auflage 3).

Bewertung der Slice-Text-Korrektur (F-4/F-5 des zweiten Reviews, Commit `2b16cf2`):
- **Ehrlich, keine Zusage weggelassen.** §1 Ausgangslage und DoD-Punkt 3 nennen jetzt den Ist-Stand: nur der `vfs`-Fall hat ein host-seitiges `doctor`-Item; Daemon- und Egress-Fallback entstehen im Startscript. Der Lastenheft-Wortlaut (CLI-005 Ausnahme 2) bleibt wörtlich im DoD („in der Befehlsausgabe und in `doctor`"), mit dem Zusatz „soweit der Host sie ermittelt". Vor der Korrektur behauptete §1 Items in `doctor_sandbox.go`, die es nicht gibt; das ist nun behoben. Die Einschränkung ist eine Beschreibung des Ist-Stands, keine Abschwächung der Zusage von `LH-FA-CLI-005`: Der Satz verlangt die Sichtbarkeit dessen, was ermittelt wird, nicht die Ermittlung. <!-- d-check:ignore (Kennung als Zitat) -->
- `--force` bei `generate`: der Slice-Text streicht `generate` aus der Ausnahme-4-Zeile mit der Begründung „kein `--force`-Pfad und kein `AffectedFiles`-Produzent". Das stimmt mit Review 2 überein (nur `InitProjectService.emitSummary` ruft `AffectedFiles`) und ist durch das Lastenheft gedeckt: Ausnahme 4 verlangt nur die Zusammenfassung, die `--force` ausgibt. <!-- d-check:ignore (Bezeichner als Zitat) -->
- **Nicht nachgezogen (Befund V-1):** Der Slice-Text trägt die vier Folge-Entscheidungen aus Review 2 nicht. `down --quiet` puffert (Commit `b31451e`), aber §1 („Ansatz"), §2 und §3 (Zeile `down.go`: „prüfen") kennen den Puffer nicht; er steht nur im Handbuch, CHANGELOG und Code. Die Vorschau-Frage (Review-F-2) fehlt in §1 „Ausdrücklich NICHT" und in §6. Die Bezugszeile nennt weiter „F-1 bis F-6" und nur den ersten Review.

Bewertung F-3 des zweiten Reviews (CHANGELOG „letzte Option gewinnt"): **behoben.** Der Eintrag sagt jetzt „`--debug` still takes precedence over `--quiet`" als geltend und die Regel nur noch als folgend („follows in a later release"); der widersprüchliche Klammer-Zusatz ist entfernt. Ein Grep nach „last"/„letzte" im CHANGELOG findet nur diese Stelle. Der Satz stimmt mit dem Test `debug-takes-precedence-over-quiet` überein. Der Eintrag ist wieder ein sehr langer Absatz (kosmetisch).

Bewertung Review-F-1 (`down`): **behoben.** `down.go` puffert den Compose-Strom unter `--quiet` (nicht unter `--json`), schreibt ihn nur im Fehlerfall nach stderr, verwirft ihn im Erfolgsfall; Hilfetext, Handbuch und CHANGELOG stimmen überein. `TestExecute_Down_Quiet_ComposeStream` pinnt beide Fälle (M-1 und M-3 rot). Vertragsbezug: Das Lastenheft nennt `down` unter den ändernden Subkommandos mit „nur Fehler"; der Puffer erhält die Fehlerdetails, die der Adapter nicht in die Fehlermeldung übernimmt. Er ist eine Setzung des Projektinhabers und im Lastenheft nicht ausdrücklich beschrieben (Ausnahme 3 nennt nur `up`); der Fehlerfall bleibt unter „nur Fehler" gedeckt.

Bewertung Review-F-2 (Vorschau unter `--quiet`): **nicht Teil dieses Slice, korrekt ausgelagert.** Die Folge-Slices [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md) (Lastenheft 0.10.0) und [`slice-v1-cli-quiet-vorschau`](../plan/planning/open/slice-v1-cli-quiet-vorschau.md) liegen in `open/`, stehen in der Welle (§3, §4) und im `doc-trace`-Beleg von `LH-FA-CLI-005`. Der Ist-Stand in diesem Slice (Diff läuft vor dem Quiet-Zweig, ein reines `--dry-run --quiet` ist leer) bleibt bis dahin bestehen und gehört in §6/§7 dieses Slice als offenes Risiko mit Verweis. <!-- d-check:ignore (Kennung als Zitat) -->

Bewertung Entscheidung (3): **vom Lastenheft-Satz gedeckt, aber ein Vertragsloch an anderer Stelle.**
- `LH-FA-CLI-005` Ausnahme 2 ist eine Sichtbarkeits-Ausnahme: Was die Befehlsausgabe oder `doctor` an Fallbacks ausweist, bleibt unter `--quiet` sichtbar. Fallbacks, die der Host nicht ermittelt, werden von `--quiet` nicht versteckt; `--quiet` verursacht die Lücke nicht (Startscripte laufen im Container). Der Slice verletzt `LH-FA-CLI-005` damit nicht. <!-- d-check:ignore (Kennung als Zitat) -->
- `LH-FA-DEV-007` verlangt dagegen eigenständig: „jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen" (drei Fallbacks). Produktiv ist nur `vfs` im `doctor` erreichbar; für Daemon und Egress gibt es keine host-seitige Quelle, und für `init`/`generate`/`up` erzeugt keine produktive Quelle einen Fallback-Eintrag. Das ist eine Lücke zwischen `LH-FA-DEV-007` und der Umsetzung, unabhängig von `--quiet`. Sie ist hier nur im Slice-Text „benannt"; ein Anker mit Slice-Kennung fehlt. <!-- d-check:ignore (Kennungen als Zitat) -->
- **Folge-Arbeit nötig (nicht angelegt):** Titel „Degradations-Fallbacks Daemon und Egress in Befehlsausgabe und `doctor` ausweisen, oder `LH-FA-DEV-007` auf das host-seitig Ermittelbare präzisieren (CR)"; vorgeschlagene Kennung `slice-v1-sandbox-fallback-ausweisung` (alternativ als CR `slice-cr-sandbox-fallback-ausweisung`, wenn der Vertrag statt des Codes geändert wird); Anker: `LH-FA-DEV-007` (Absatz Degradation), `LH-FA-CLI-005` Ausnahme 2, Welle `welle-cli-verbosity` §6 dieses Slice. Ob Code (Startscript meldet an den Host, etwa über ein Statusfile, das `doctor` liest) oder CR gewählt wird, entscheidet der Projektinhaber. <!-- d-check:ignore (vorgeschlagene Kennungen, noch nicht angelegt) -->

Plan-vs-Code-Diff (Slice §3):
- `root.go`, `cli.go`: Fehler-Pegel und Kommentare — Ist: gesetzt; `[visibleWarnings]` löst auf.
- `init.go`, `generate.go`, `add.go`, `remove.go`, `config.go`: Zusammenfassungen unter `--quiet` unterdrückt, `remove --quiet` — Ist: gesetzt, Fehlerpfade vor den Quiet-Zweigen.
- `up.go`: Tabelle und Diagnosen unterdrückt, Lockerungen und Fallback-Einträge sichtbar, Compose-Strom bleibt — Ist: gesetzt.
- `down.go`: Plan „prüfen" — Ist: **Code geändert** (Compose-Puffer, `b31451e`); Plan nicht nachgezogen (V-1).
- `doctor.go`: Fehler- und Fallback-Items, Kopf/Summary entfallen, F-3-Meldung — Ist: gesetzt (`strictWarnFailure`, Sentinel bleibt `errors.Is`-fähig, Exit 11 unverändert, Wire-Form unverändert).
- `removeservice.go` (`WarningEntry`), `doctor_sandbox.go`, `devcontainer_sandbox.go`, `domain/diagnostic.go`: Kennzeichen — Ist: gesetzt, `json:"-"` am `WarningEntry`, keine neuen Imports außer `strings` in `doctor.go`. Plan-Abweichung: `devcontainer_sandbox.go` und `diagnostic.go` stehen nicht in §3 (zu Recht im Diff, nachzutragen); `domain.Diagnostic.DegradationFallback` ist ohne Pin gegen das Wire-Format (Review-2 F-6, INFO).
- Tests, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` — Ist: gesetzt. Nicht im Plan, im Diff: `docs/maintainer/cli-json-contract.md` (konsistent), `doctor_json_test.go` (Pin des Fallback-Flags).
- Kein `//nolint`, keine Lockerung von Gates, Schwellen oder Bestandspins (`debug-takes-precedence-over-quiet` unverändert).

Mutationen (selbst gefahren, nur über `make test`, danach `git checkout -- <Datei>`; `git status` sauber):
| Nr | Mutation | Ergebnis |
| --- | --- | --- |
| M-1 | `down.go`: gepufferten Strom im Fehlerfall nicht ausgeben | rot: `TestExecute_Down_Quiet_ComposeStream/quiet-failure` |
| M-2 | `doctor.go`: `|| item.DegradationFallback` aus dem Quiet-Filter entfernt | rot: `TestExecute_Doctor_Quiet_FallbackItemAndErrorKept_OtherWarnsDropped` |
| M-3 | `down.go`: Puffer nie aktiv (`flags.Quiet && !flags.JSON` → `false`) | rot: `TestExecute_Down_Quiet_ComposeStream/quiet-success` |

Die Pins prüfen Verhalten. Schwach bleibt allein der `--force`-Pin (nur Request-Feld).

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (lint, test, coverage-gate 91,00 %, docs-check, baseline-verify) |
| `make doc-trace` | pass | selbst gefahren: „83 Anforderung(en), 0 Waise(n)."; die `LH-FA-CLI-005`-Zeile trägt diesen Slice und beide Vorschau-Slices | <!-- d-check:ignore (Kennung als Zitat) -->
| `make test` (3 Mutationen) | rot wie erwartet | siehe Tabelle; Arbeitsbaum zurückgesetzt |
| `make docs-check` | pass | nach Anlegen dieser Datei: siehe Abschluss in der Übergabe |
| `make test-docker` | not run | kein Docker-Pfad berührt: der Diff ändert CLI-Rendering, Fake-Use-Case-Tests und Doku; der Compose-Strom von `down` ist über den Sink-Port gepinnt, nicht über eine echte Engine |
| `make verify-depguard` | not run | keine Importregel geändert (neu nur `bytes`/`strings` in `cli`-Dateien, kein Port- oder Layer-Import) |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `quiet_test.go`, `quiet_modifying_test.go`, `cli_test.go`, `verbosity_test.go`; Handbuch, CHANGELOG; `doc-trace` ok |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Pin `wrapped ErrDoctorFailures` → 11, `--strict` mit Warn-Items Exit 11, Fehlermeldung auf stderr |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `TestExecute_QuietJSON_FallbackFlagNotInWireForm`, `doctor_json_test.go`; `--quiet --json` unverändert |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | nur `vfs`-Fall produktiv (`doctor_sandbox_test.go`); Daemon/Egress offen (Lücke, siehe oben) |
| [`LH-FA-DIAG-003`](../../spec/lastenheft.md#lh-fa-diag-003--fehlerklassifikation) | `--strict`-Meldung nennt Flag, Anzahl, IDs |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — siehe Sensors.

Commit / Artefakt:
- Slice-Commits `c277c47`, `72006d3`, `7a244b9`, `4c68a86`, `1e35794`, `82af811`, `b31451e`, `2b16cf2`; HEAD zum Prüfzeitpunkt `d5a25d2`.

## Verdikt

**Bestanden, mit Auflagen für die Closure.** Code und Verhalten decken `LH-FA-CLI-005` (0.9.0) und die drei Review-Fixes ab; Gates grün, Waisen null, drei Mutationen fallen rot. Offen sind Nachzüge im Slice-Text, die Closure-Notiz und eine Folge-Arbeit für `LH-FA-DEV-007`. <!-- d-check:ignore (Kennung als Zitat) -->

Auflagen:
1. Closure-Notiz §7 mit Steering-Loop-Lerneintrag; Finding-Klassen aus beiden Reviews (Ausgabepfad an `--quiet` nicht durchgesetzt, Vertragslücke zwischen zwei Anforderungen, CHANGELOG nennt noch nicht gelieferte Regel als geltend, Plan nennt Quellen, die es nicht gibt, Ausnahmepfad nur über Request-Feld gepinnt, Wire-Ausschluss ohne Pin). Register: Beobachtung eintragen oder „keine Beobachtung" notieren. Vor `done/`: DoD-Häkchen setzen, dann reiner `git mv` (`AGENTS.md` §3.3). <!-- d-check:ignore (Pfad als Zitat) -->
2. §1, §3 und §6 nachziehen (V-1): `down`-Puffer im „Ansatz" und in der §3-Zeile `down.go`; `devcontainer_sandbox.go` und `domain/diagnostic.go` in §3; Bezugszeile auf beide Reviews. Neue Risiken in §6 mit Ausgang: (a) **Vorschau unter `--quiet`** (`--dry-run --quiet` leer, Exit 0) — Ausgang: weiter offen, Träger [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md) und [`slice-v1-cli-quiet-vorschau`](../plan/planning/open/slice-v1-cli-quiet-vorschau.md); (b) **Fallbacks Daemon/Egress ohne host-seitige Ausweisung** — Ausgang: weiter offen, Träger die Folge-Arbeit unten; (c) **`down`-Puffer als Setzung ohne ausdrückliche Lastenheft-Zeile** — Ausgang: entfallen (Entscheidung des Projektinhabers, gepinnt, Fehlerfall durch „nur Fehler" gedeckt), oder bei Bedarf in den Vorschau-CR aufnehmen.
3. Folge-Arbeit für `LH-FA-DEV-007` mit Slice-Kennung anlegen (Planner, nicht dieser Lauf): Titel und Anker siehe Bewertung Entscheidung (3). In §7 als Folge-Slice nennen; bis dahin trägt „Lücke benannt" keinen auflösbaren Anker (Paarung offen). <!-- d-check:ignore (Kennung als Zitat) -->
4. Den `--force`-Pin stärken oder als Grenze in §7 benennen: Der Test prüft nur `SilenceProgress` im Request. Zusätzlich in §7 festhalten: Pin für das Doctor-Item gegen das Wire-Format fehlt (Review-2 F-6) — `doctor_json_test.go` pinnt nach `b31451e` das Fallback-Flag; Stand beim Closure gegenprüfen.
5. Kein Release-Cut vor `slice-v1-cli-verbosity-letzte-option-gewinnt` (CHANGELOG-Präsens „follows in a later release") und vor den Vorschau-Slices, sofern der Release die Vorschau-Zusage berührt. <!-- d-check:ignore (Kennung als Zitat) -->
