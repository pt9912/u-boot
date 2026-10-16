# Verification Evidence: slice-cr-cli-quiet-vorschau — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der CR richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** CR-Slice [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md) (Delivery `49f459a`, Review-Fix `d7c2e3a`), Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md). Maßstab: Slice §1 (CR-Kern, Abnahmekriterien, Historie-Zeile, Nicht-Ziele), §2 (DoD), §3, §6. Review-Report: [`2026-10-16-slice-cr-cli-quiet-vorschau.md`](2026-10-16-slice-cr-cli-quiet-vorschau.md).

## Verification Evidence

Scope:
- Slice: [`slice-cr-cli-quiet-vorschau`](../plan/planning/done/slice-cr-cli-quiet-vorschau.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe)
- Artefakte: `spec/lastenheft.md` (Versionsfeld, `LH-FA-CLI-005`, §16), `CHANGELOG.md` `[Unreleased]`, Slice-Entwurf der Historie-Zeile, Review-Report; kein Go-Code im Diff. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2):
- [x] Lastenheft 0.10.0 — Evidence: Versionsfeld `0.10.0` (Spaltenbreite erhalten); der Absatz in `LH-FA-CLI-005` ist nach Link-Normalisierung zeichengleich zum kursiven CR-Kern in Slice §1 (per Skript geprüft, Ergebnis `True`); §16-Zeile 0.10.0 mit Datum, Link auf die Anforderung und „Vereinbarung mit dem Projektinhaber". Stufen, vier Ausnahmen, Abfragebefehle-Satz und Mehrfachangabe-Regel im Diff unverändert. <!-- d-check:ignore (Kennung als Zitat) -->
- [x] Umsetzungsbedarf einem Slice zugeordnet — Evidence: [`slice-v1-cli-quiet-vorschau`](../plan/planning/done/slice-v1-cli-quiet-vorschau.md) (in `open/`) trägt in seiner DoD: Punkt 1 deckt Abnahmekriterien 1 und 2 des CR (Tests je Subkommando mit `--dry-run --quiet`, `--dry-run --diff --quiet`, `--diff --quiet`; Boundary „Erfolgsmeldung nach dem Schreiben bei `--diff --quiet` entfällt"; Negative: Fehler bleiben), Punkt 2 deckt Kriterium 3 (`benutzerhandbuch.md` und `CHANGELOG.md`). Rest-Risiken dort (Vergangenheitsform, `--diff` ohne `--dry-run`) tragen „offen bis Entscheidung" und sind dort zu pinnen.
- [x] Doku-Update CHANGELOG — Evidence: `[Unreleased]` nennt Lastenheft 0.10.0, vier Ausnahmen und, nach `d7c2e3a`, „The output of `--dry-run` and `--diff` is the result of the call and stays under `--quiet`. The implementation of that rule and the rule that the last verbosity option wins follow in a later release." Kein Satz widerspricht dem Lastenheft 0.10.0 (siehe Bewertung F-1).
- [x] `make gates` grün — selbst gefahren, Exit 0.
- [x] Review-Report liegt vor — `docs/reviews/2026-10-16-slice-cr-cli-quiet-vorschau.md` (0 HIGH, 1 MEDIUM, 0 LOW, 4 INFO), Reviewer-Kontext getrennt vom Autor. <!-- d-check:ignore (Pfad als Zitat) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: §7 trägt „Offen." (Auflage 1).
- [ ] Beobachtungs-Register — offen: §8 notiert „Register leer" für die Sichtung; Fortschreibung oder „keine Beobachtung" in §7 steht aus (Auflage 1).
- [~] Jedes Risiko aus §6 trägt einen Ausgang — Risiko 2 und 3 „entfallen" mit Begründung; Risiko 1 (`--diff` ohne `--dry-run`) trägt „offen bis zur Bestätigung durch den Projektinhaber" (Auflage 2). Risiko 2 verlangt, dass der Verifier die Suche wiederholt: erledigt, siehe unten.
- [ ] Drei Paarungen — Folge-Slice getragen; Anker und Register entstehen mit §7 (Auflage 1).
- Alle DoD-Häkchen im Slice sind noch ungesetzt; sie werden mit der Closure gesetzt (Auflage 1).

CR-Kern gegen Lastenheft (wörtlich):
- Zeichengleich (siehe DoD-Punkt 1). Die Historie-Zeile in Slice §1 ist nach `d7c2e3a` inhaltlich an die geltende §16-Zeile angeglichen (Erfolgs- und Zusammenfassungsmeldung, „auch bei `--diff` ohne `--dry-run`", Reihenfolge „Stufen, vier Ausnahmen, Mehrfachangabe-Regel"). Verbleibender Unterschied rein formal: der Entwurf trägt Kennungen als Klartext und die Klammer „(Änderungsplan und Diff-Darstellung, nach LH-FA-CLI-007 und LH-FA-CLI-008)", die geltende Zeile Links; kein Gehaltsunterschied, die §16-Zeile gilt. Nicht blockierend. <!-- d-check:ignore (Kennung als Zitat) -->
- Nicht-Ziele eingehalten: Diff berührt weder `internal/`, `spezifikation.md`, `architecture.md` noch `docs/user/`; Compose-Strom von `down` und CHANGELOG-Widerspruch des Code-Slices nicht angefasst.

Wiederholte Suche nach älteren kollidierenden Zusagen im Lastenheft (Begriffe `--dry-run`, `--diff`, `Vorschau`, `LH-FA-CLI-007`, `LH-FA-CLI-008`, `--quiet`, `Zusammenfassung`; selbst gefahren, nicht aus dem Review übernommen): <!-- d-check:ignore (Kennung als Zitat) -->
- `--quiet` / Verbosity: nur `LH-FA-CLI-005` (Z. 184 bis 198). <!-- d-check:ignore (Kennung als Zitat) -->
- `--dry-run` / `--diff` / Vorschau: `LH-FA-CLI-007` (Plan „muss anzeigen", `--dry-run` mit `--diff` „ausschließlich … Änderungsplan und Diff-Darstellung", `--dry-run --json` streng maschinenlesbar), `LH-FA-CLI-008` (Vorschau; Beispiel „Vorschau mit anschließendem Schreiben"; gleiche Exit-Codes) und der JSON-Abschnitt (Schema-Absätze, nur `--json`). <!-- d-check:ignore (Kennung als Zitat) -->
- „Zusammenfassung": `LH-FA-CLI-005` (Ausnahme 4), `LH-FA-CLI-005A` und `--force` in `LH-FA-INIT-005`/Optionsliste von `init` (Pfad-Zusammenfassung vor dem Überschreiben, Ausnahme 4, unabhängig von der Vorschau). <!-- d-check:ignore (Kennung als Zitat) -->
- Ergebnis: Keine weitere Zusage kollidiert mit „`--quiet` zeigt nur Fehler"; die Kollisionen mit `LH-FA-CLI-007` und `-008` löst der neue Satz. Deckt sich mit der Tabelle in Slice §1 und dem Review. Risiko-2-Ausgang „entfallen" ist belegt. <!-- d-check:ignore (Kennung als Zitat) -->

Bewertung der Review-Findings:
- **F-1 (MEDIUM), CHANGELOG:** behoben in `d7c2e3a`. Der Satz sagt nicht mehr „bleibt" und „folgt später" zur selben Ausgabe; „follow in a later release" bezieht sich auf die Umsetzung. Die Plan-Zusammenfassung wird nicht mehr als zurückgestellt genannt. Der vorangehende Satz „Success and summary messages of the subcommands that change something disappear under `--quiet`" bleibt neben dem Vorschau-Satz lesbar (Vorschau ist Ergebnis, Erfolgsmeldung nicht). Kein Widerspruch zum Lastenheft 0.10.0; Hinweis: ein Release-Cut vor `slice-v1-cli-quiet-vorschau` liefert einen Eintrag aus, dessen Zusage der Code noch nicht erfüllt (siehe Auflage 3). <!-- d-check:ignore (Kennung als Zitat) -->
- **F-2 (INFO), „Änderungsplan" gegen „Zusammenfassung":** nicht blockierend. Der Vertrag nennt „geplanter Änderungsplan" (Wortlaut aus `LH-FA-CLI-007`); die Abbildung auf die heutigen „Would …"- bzw. Vergangenheitsform-Ausgaben ist Umsetzung und steht als Risiko im Code-Slice. Folge: dort pinnen. <!-- d-check:ignore (Kennung als Zitat) -->
- **F-3 (INFO), Historie-Entwurf:** erledigt durch `d7c2e3a` (siehe oben).
- **F-4 (INFO), `--quiet --json` im Lastenheft nicht genannt:** nicht blockierend. Der Satz berührt `--json` nicht; `LH-FA-CLI-007` hält `--dry-run --json` streng maschinenlesbar; kein widersprüchlicher Fall. <!-- d-check:ignore (Kennung als Zitat) -->
- **F-5 (INFO), Benutzerhandbuch folgt dem Code:** nicht blockierend. Zuständig ist der Folge-Slice (DoD Punkt 2, Zeilen um 299 und 1280 in seinem Plan §3); kein Handbuch-Update in diesem CR geplant (Plan §3 nennt nur Lastenheft und CHANGELOG).

Konsistenz der Planungsdateien (0.10.0 und „vier Ausnahmen"):
- `welle-cli-verbosity.md` §1: „vier Ausnahmen", Vorschau-Satz mit Links auf CLI-007/-008, Version 0.10.0 (Z. 43, 44), CR-Zeile in der Tabelle (Z. 56) und Trigger-Eintrag in §2 stimmig; keine „drei Ausnahmen". Befund V-1: Z. 55 trägt für den fertigen Code-Slice `slice-v1-cli-quiet-nur-fehler` noch „(Lastenheft 0.9.0; in Arbeit)", obwohl er in `done/` liegt (Zustandsfeld). Das ist 0.9.0 als zutreffende Version dieses Slice, aber „in Arbeit" ist veraltet. <!-- d-check:ignore (Kennung als Zitat) -->
- `open/README.md`: Z. 18 „vier Ausnahmen … Lastenheft 0.9.0" (korrekt, historisch), Z. 19 „vier Ausnahmen … blockiert bis Lastenheft 0.9.0" (Verbform veraltet, der Slice liegt in `done/`; vorbestehend), Z. 20 und 21 nennen 0.10.0 und die Folge-Slices. Keine alte Zählung („drei") mehr in diesen Dateien. Befund V-2: Z. 19 „blockiert bis" und Z. 20 verweist auf einen Slice in `in-progress/` in einem `open/`-Index (vorbestehendes Muster der Datei; mit den Closures zu bereinigen).
- `open/slice-v1-cli-verbosity-letzte-option-gewinnt.md`: Z. 72 „die vier Ausnahmen", Z. 87 „Lastenheft 0.10.0 (nach den CRs)"; keine alte Zählung oder Fassung. Kein Befund. <!-- d-check:ignore (Kennung als Zitat) -->
- `done/slice-v1-cli-quiet-vorschau.md`: durchgängig 0.10.0 und Verweise auf den CR in `in-progress/`; kein Befund. Befund V-3: sein Start-Trigger („CR in `done/`") ist bis zur CR-Closure nicht erfüllt, das stimmt mit dem Plan überein. <!-- d-check:ignore (Kennung als Zitat) -->
- Alle drei Zeilen mit „0.9.0" außerhalb der erledigten Slices sind historisch zutreffende Versionsangaben der früheren CRs; keine ist als geltende Fassung geschrieben.

Welle-Reihenfolge:
- CR in `in-progress/`, Code-Slice `slice-v1-cli-quiet-vorschau` und `slice-v1-cli-verbosity-letzte-option-gewinnt` in `open/`; Reihenfolge „Lastenheft ändern → Umsetzung" eingehalten (Code unberührt). Spezifikation im Diff unberührt (Welle-Closure-Trigger). <!-- d-check:ignore (Kennung als Zitat) -->

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (lint, test, coverage-gate, docs-check, baseline-verify) |
| `make doc-trace` | pass | selbst gefahren: „83 Anforderung(en), 0 Waise(n)."; `LH-FA-CLI-005`, `-007`, `-008` je „ok" und tragen `slice-cr-cli-quiet-vorschau` als Beleg | <!-- d-check:ignore (Kennung als Zitat) -->
| `make docs-check` | pass | Teil von `make gates` vor Anlegen dieser Datei; nach Anlegen siehe Abschluss |
| `make test-docker` | not run | kein Docker-Pfad berührt (Diff: Lastenheft, CHANGELOG, Planung, Review) |
| `make verify-depguard` | not run | keine Importregel berührt, kein Go-Code im Diff |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Lastenheft 0.10.0 (Absatz, §16), CHANGELOG; Umsetzungsbeleg (Tests) liegt im Code-Slice, hier nicht Ziel |
| [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) | im Satz zitiert; Wortlaut unverändert |

Plan-vs-Code-Diff:
- Plan §3 `spec/lastenheft.md` (Anforderung, Versionsfeld, §16-Zeile) — Ist: gesetzt (`49f459a`). <!-- d-check:ignore (Pfad als Zitat) -->
- Plan §3 `CHANGELOG.md` — Ist: Version 0.10.0, Vorschau-Satz; Korrektur in `d7c2e3a`. <!-- d-check:ignore (Pfad als Zitat) -->
- Nicht im Plan, im Diff: Review-Report und Angleichung des Historie-Entwurfs im Slice (`d7c2e3a`), beides Review-Folgen innerhalb des Slice.
- Plan nicht umgesetzt: nichts.

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — kein berührter Pfad (siehe Sensors).

Commit / Artefakt:
- Delivery `49f459a`; Review-Fix samt Review-Report `d7c2e3a`.

## Verdikt

**Bestanden, mit Auflagen für die Closure.** Der Vertrag (Lastenheft 0.10.0) deckt sich wörtlich mit dem CR-Kern, die wiederholte Suche findet keine weitere kollidierende Zusage, F-1 ist behoben und F-2 bis F-5 blockieren nicht. Gates und RTM sind grün. Offen sind die Closure-Notiz, die Bestätigung des Projektinhabers und kleine Planungs-Nachzüge.

Auflagen:
1. Vor dem Move nach `done/` (Inhalt zuerst, dann reiner `git mv`, `AGENTS.md` §3.3): DoD-Häkchen setzen, §7 Closure-Notiz mit Steering-Loop-Lerneintrag (Finding-Klassen des Reviews: Changelog-Zusage widerspricht sich selbst, Vertragsbegriff „Änderungsplan" nicht von „Zusammenfassung" getrennt, Slice-Entwurf der Historie-Zeile weicht vom gesetzten Text ab), Register-Vermerk und drei Paarungen belegen. <!-- d-check:ignore (Pfad als Zitat) -->
2. Der Projektinhaber bestätigt vor der Closure (a) die Version 0.10.0 statt 0.9.1 und (b) den Schlusssatz „Erfolgs- und Zusammenfassungsmeldung … entfällt auch bei `--diff` ohne `--dry-run`" (Slice §1 und §6 Risiko 1); dann §6 Risiko 1 auf „entfallen" oder „eingetreten" setzen. Bei Gegenentscheidung entfällt der Schlusssatz in Lastenheft, §16-Zeile und Slice (Rückführung nach §4).
3. Kein Release-Cut vor `slice-v1-cli-quiet-vorschau` und `slice-v1-cli-verbosity-letzte-option-gewinnt` (CHANGELOG nennt eine Zusage, die der Code noch nicht erfüllt). Im Code-Slice die Abbildung „Änderungsplan" auf die „Would …"-/Vergangenheitsform-Ausgaben pinnen (F-2). <!-- d-check:ignore (Kennung als Zitat) -->
4. Planungs-Nachzüge (nicht vertragsrelevant, bei der Closure mitzuführen): `welle-cli-verbosity.md` Z. 55 „in Arbeit" bei `slice-v1-cli-quiet-nur-fehler` (liegt in `done/`) ersetzen; `open/README.md` Z. 19 „blockiert bis Lastenheft 0.9.0" und die Einträge für Slices, die nicht mehr in `open/` liegen, beim nächsten Pflegeschritt bereinigen; Historie-Entwurf in Slice §1 optional auf Links angleichen (rein formal). Folgende Dateien tragen bereits durchgängig 0.10.0 und „vier Ausnahmen": `open/slice-v1-cli-verbosity-letzte-option-gewinnt.md`, `done/slice-v1-cli-quiet-vorschau.md`, Welle §1, §2, Tabelle. <!-- d-check:ignore (Pfade als Zitat) -->
