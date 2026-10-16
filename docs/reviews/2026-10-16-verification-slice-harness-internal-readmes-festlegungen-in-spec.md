# Verification Evidence: slice-harness-internal-readmes-festlegungen-in-spec — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Slice [`slice-harness-internal-readmes-festlegungen-in-spec`](../plan/planning/done/slice-harness-internal-readmes-festlegungen-in-spec.md), Commit-Range `145fc42..HEAD` (Implementer `b636076`..`db27b3a`, Review-Fix `16cf4b4`, `58fd64d`, `e75655b`, Review-Report `fceea11`), Welle [`welle-harness-wartung`](../plan/planning/done/welle-harness-wartung.md) §3 (Closure-Trigger).

## Verification Evidence

Scope:
- Slice: [`slice-harness-internal-readmes-festlegungen-in-spec`](../plan/planning/done/slice-harness-internal-readmes-festlegungen-in-spec.md)
- IDs: [`LH-LESE-003`](../../spec/lastenheft.md#lh-lese-003--dokumentenordnung), [`SPEC-012`](../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (nur Abweichungsbefund)
- Artefakte: `spec/spezifikation.md` (`SPEC-036`, `LH-FA-UP-001.a`, `SPEC-038`..`SPEC-042`), `spec/architecture.md` §2.1, sieben `internal/**/README.md`, `harness/conventions.md`, `internal/hexagon/port/driven/clock.go` (Kommentar) <!-- d-check:ignore (Kennungen und Pfade als Zitat) -->

DoD-Abgleich (Slice §2):
- [x] Inventur — Evidence: Slice §9 deckt alle sieben READMEs (application, domain, port/driving, port/driven, adapter/driven, adapter/driving, `internal/README.md`), je Aussage Einordnung F/R/A/I und benannte Abweichungen README gegen Code. Der Slice-Haken in §2 steht noch aus (Implementer).
- [x] Spezifikation — Evidence: `SPEC-036`, `LH-FA-UP-001.a`, `SPEC-038`..`SPEC-042` an den verfeinerten Anforderungen verankert, 7 Einträge für 7 Anforderungsgruppen (gebündelt, nicht je Konstante). Stichprobe gegen den Code: `SPEC-039` Changelog-Reparatur entspricht `changelogRepairUnreleased`/`hasChangelogUnreleased` in `internal/hexagon/application/generate.go` (ganze Datei außer Codeblöcken, Release-Abschnitt nur nach END-Markierung); Review-Finding F-1 damit korrekt behoben. Zahlenwerte aus dem Code laut Review-Negativbefunden (500 ms, 300 ms, Schwelle 3, 14 Pfad-Arten). <!-- d-check:ignore (Kennungen als Zitat) -->
- [x] READMEs gekürzt (490 auf 109 Zeilen, 7 Dateien, je 12..22 Zeilen); `make docs-check` grün — Evidence: Sensors. Stichprobe `application/README.md`: nur Zweck, Verweise auf Architektur und `SPEC-*`, Hinweis auf `make doc-trace`, keine Zahlen, Schwellen oder Tabellen. Grep über alle sieben READMEs ohne Zahlenwert-Treffer (einziger Treffer: Nennung „Restart-Schwelle" als Begriff in `domain/README.md` mit Verweis).
- [x] `make gates` grün — Evidence: Sensors.
- [x] Review-Report — Evidence: `docs/reviews/2026-10-16-slice-harness-internal-readmes-festlegungen-in-spec.md` (0 HIGH, 1 MEDIUM, 2 LOW, 5 INFO), Verfahren mit Reviewer-Skill. Behebungen stichprobenhaft verifiziert: F-1 (siehe oben), F-2 (`SPEC-037` heißt `LH-FA-UP-001.a`, Nummer `037` ist nirgends belegt, Lücke), F-4 (`Clock.Sleep`-Begründung als Kommentar in `clock.go`, einzige Code-Änderung, 4 Kommentarzeilen, kein Wertwechsel), F-5 (Folge-Slice-Ergänzung `e75655b`). <!-- d-check:ignore (Pfad und Kennungen als Zitat) -->
- [ ] Closure-Notiz — offen: Slice §7 lautet „Offen."; Haken in §2 ungesetzt. Auflage unten.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: alle drei Risiken stehen weiter mit „Ausgang: offen". Sachlich sind sie erledigt (Werte aus dem Code; Bündelung je Anforderung; Modus-Zeile in `harness/conventions.md` mit `db27b3a` nachgezogen). Auflage unten.

Plan-vs-Code-Diff (Slice §3):
| Plan-Zeile | Ist |
|---|---|
| `internal/README.md` und sechs Paket-READMEs: update | erfüllt, sieben Dateien gekürzt |
| `spec/spezifikation.md`: neue Einträge | erfüllt (`SPEC-036`..`SPEC-042` mit `SPEC-037` als `LH-FA-UP-001.a`) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `spec/architecture.md`: prüfen | Lücke gefunden und geschlossen (§2.1 `SeverityInfo`), eine Zeile, ohne Slice-/ADR-Bezug |
| `internal/**/doc.go`: neu, wo nötig | nicht angelegt; Implementierungshinweise blieben am Code, ein Kommentar ergänzt (`clock.go`) |
| `harness/conventions.md`: prüfen | Zeile der Sub-Area nachgezogen (`db27b3a`), Modus unverändert Greenfield |
| Nicht-Ziele: kein Code-/Wertwechsel, keine neuen Anforderungen | eingehalten: `git diff` über `*.go` zeigt nur den Kommentar in `clock.go`; Lastenheft unverändert |

Regelkonformität:
- Lücken `SPEC-011` und `SPEC-037` sind regelkonform: Es gibt keine Überschrift mit diesen Nummern; keine lebende Referenz im Repo außerhalb der Review-Reports; Baseline `grundlagen-source-precedence.md` §Vergabe: Lücken werden nicht nachbelegt. <!-- d-check:ignore (Kennungen und Pfad als Zitat) -->
- Decken-Regel: laut Review keine neue Zeile in `spec/` mit Verweis nach unten; `spec/architecture.md` trägt keinen Slice-/ADR-Bezug.
- Welle-Closure-Trigger §3: der Review-Report (`fceea11`) liegt nach dem Move von `slice-harness-dcheck-bump-reviewer-skill` nach `done/` (`c9b1f4e`), und der Reviewer-Skill (`3ec4535`) stammt aus diesem Slice; Reihenfolge belegt per `git log`. `make doc-trace`: 83 Anforderung(en), 0 Waise(n); `LH-FA-UP-001` und `LH-FA-UP-003` tragen den Slice als Beleg. <!-- d-check:ignore (Kennungen als Zitat) -->

Abweichung Lastenheft gegen Code (Verbosity): kein Teil dieses Slice; sie steht korrekt nicht in der Spezifikation. Zwei Folge-Slices liegen in `open/`: [`slice-v1-cli-verbosity-letzte-option-gewinnt`](../plan/planning/in-progress/slice-v1-cli-verbosity-letzte-option-gewinnt.md) und [`slice-cr-cli-quiet-fehler-und-warnungen`](../plan/planning/done/slice-cr-cli-quiet-fehler-und-warnungen.md). Im Slice selbst fehlt der Verweis: §9 spricht nur von „offener Befund für den Projektinhaber", §6 hat keinen Eintrag. Auflage unten.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (zweimal); `coverage-gate: OK — coverage 90.90% meets threshold 90%`; lint, test, docs-check, baseline-verify Teil von `gates` |
| `make docs-check` | pass | selbst gefahren: 249 Dateien, 0 Befunde (vor dem Anlegen dieses Berichts) |
| `make doc-trace` | pass | selbst gefahren: 83 Anforderung(en), 0 Waise(n) |
| `make test-docker` | not run | kein Docker-Laufzeit-Pfad berührt (nur Doku und ein Kommentar) |
| `make verify-depguard` | not run | keine Importregel berührt |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-LESE-003`](../../spec/lastenheft.md#lh-lese-003--dokumentenordnung) | Festlegungen stehen in `spec/spezifikation.md`, READMEs verweisen; Slice §9 |
| [`SPEC-012`](../../spec/spezifikation.md#spec-012--schichten-und-verzeichnislayout-unter-internal) | Verweis im Zweck-Satz der READMEs |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Abweichungsbefund Slice §9; Folge-Slices in `open/` |

Verdikt: **bestanden unter Auflagen** — Inhalt und Sensors tragen, keine DoD-Verletzung in der Sache; formale Closure-Schritte fehlen.

Auflagen für die Closure (Implementer/Planner):
1. Slice §7 Closure-Notiz mit Steering-Loop-Lerneintrag schreiben (Finding-Klassen des Reviews aus dem Report übernehmen; F-3 als Frage an den Architect, Kopf-Regel für Mehr-Anforderungs-Verfeinerung, benennen).
2. Slice §6: jedem der drei Risiken einen Ausgang geben; neues Risiko „Verbosity-Abweichung" mit Ausgang auf die beiden Folge-Slices in `open/` aufnehmen; §9 den Verweis „offener Befund" durch die Slice-Links ersetzen.
3. §2-DoD-Haken setzen (nach Closure-Notiz), dann reiner `git mv` nach `done/` (`make slice-mv`), danach Welle §3 Haken und Welle-Closure. <!-- d-check:ignore (Make-Target als Zitat) -->
4. Nach der Closure `make docs-check` und `make gates` erneut fahren (Move zieht Verweise nach).
