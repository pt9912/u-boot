# Verification Evidence: slice-cr-cli-quiet-weitere-ausnahmen — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der CR richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** CR-Slice [`slice-cr-cli-quiet-weitere-ausnahmen`](../plan/planning/in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md) (Delivery `0052e1c`, Review-Fixes `591acb7`), Welle [`welle-cli-verbosity`](../plan/planning/welle-cli-verbosity.md). Maßstab: Slice §1 (CR-Kern, Abnahmekriterien, Historie-Zeile), §2 (DoD), §3, §6. Review-Report: [`2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md`](2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md).

## Verification Evidence

Scope:
- Slice: [`slice-cr-cli-quiet-weitere-ausnahmen`](../plan/planning/in-progress/slice-cr-cli-quiet-weitere-ausnahmen.md)
- IDs: [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Änderung), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)
- Artefakte: `spec/lastenheft.md` (Versionsfeld, `LH-FA-CLI-005`, §16), `CHANGELOG.md` `[Unreleased]`, Planungsdateien; kein Go-Code im Diff. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2):
- [x] Lastenheft 0.9.0 — Evidence: Versionsfeld `0.9.0`; `LH-FA-CLI-005` ist nach `591acb7` absatzweise zeichengleich zum CR-Kern in Slice §1 (per Skript geprüft: Einleitungssatz, Ausnahmen 1 bis 4, Schlusssatz, Mehrfachangabe-Satz; Link-Ziele normalisiert); §16-Zeile 0.9.0 mit „Vereinbarung mit dem Projektinhaber" und allen vier Ausnahmen. <!-- d-check:ignore (Kennung als Zitat) -->
- [x] Umsetzungsbedarf einem Slice zugeordnet — Evidence: [`slice-v1-cli-quiet-nur-fehler`](../plan/planning/open/slice-v1-cli-quiet-nur-fehler.md) liegt in `open/`; DoD trägt Kriterium 1 (Tests je Subkommando, `remove --quiet`), Kriterium 2 (je Ausnahme ein Pin, einschließlich der `--force`-Zusammenfassung bei `init` und `generate`) und Kriterium 3 (Handbuch und CHANGELOG „mit den vier Ausnahmen"). Befund V-2 zur Stimmigkeit im Rest der Datei.
- [x] Doku-Update CHANGELOG — Evidence: `[Unreleased]` nennt vier Ausnahmen, ändernde Subkommandos, `doctor`, Abfragebefehle; kein Satz widerspricht dem Lastenheft (siehe Bewertung F-3).
- [x] `make gates` grün — siehe Sensors (selbst gefahren).
- [x] Review-Report liegt vor — Evidence: `docs/reviews/2026-10-16-slice-cr-cli-quiet-weitere-ausnahmen.md` (0 HIGH, 3 MEDIUM, 1 LOW, 3 INFO), Reviewer-Kontext getrennt vom Autor. <!-- d-check:ignore (Pfad als Zitat) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: §7 trägt „Wird bei der Closure gefüllt." (Auflage 1).
- [ ] Beobachtungs-Register — offen: §8 notiert „Register leer" für die Sichtung; Fortschreibung oder „keine Beobachtung" in §7 steht mit der Closure-Notiz aus (Auflage 1).
- [~] Jedes Risiko aus §6 trägt einen Ausgang — alle drei Risiken tragen „entfallen" mit Begründung; der CR-Kern trägt jetzt vier Ausnahmen, §6 nennt weiter nur die drei Fallbacks aus `-007` (korrekt, kein Widerspruch). <!-- d-check:ignore (Kennung als Zitat) -->
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Folge-Slice getragen; Anker und Register entstehen mit §7 (Auflage 1).
- Die DoD-Häkchen im Slice sind alle noch ungesetzt; sie werden mit der Closure-Notiz gesetzt (Auflage 1).

CR-Kern gegen Lastenheft (wörtlich, Prüfung am Text selbst):
- Zeichengleich, Ergebnis siehe DoD-Punkt 1. Die Historie-Zeile in Slice §1 weicht von der geltenden §16-Zeile ab (sie trägt noch den Satz „Die geltende Zeile steht in spec/lastenheft.md §16." und „Vier Ausnahmen" ohne Links); sie ist Entwurf, die §16-Zeile gilt (Abwärts-Verweis zu Recht entfallen). Als Entwurf zu kennzeichnen oder gleichzuziehen (Auflage 2). <!-- d-check:ignore (Pfad als Zitat) -->
- Nicht-Ziele (Code, „letzte Option gewinnt", Spezifikation) eingehalten: Diff berührt weder `internal/`, noch `spezifikation.md`, noch `architecture.md`, noch `docs/user/`.

Bewertung der drei MEDIUM-Behebungen gegen den Lastenheft-Text:
- **F-1 (Ergebnis-Satz, Kriterium):** behoben. Der Absatz trennt jetzt „Subkommandos, die etwas ändern" (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`: nur Fehler) von „Abfragebefehlen, die nichts verändern" (`config get`, `config list`, `logs`, `template list`: Ausgabe bleibt) und ordnet `doctor` ausdrücklich zu (nur Fehler-Einträge). Alle Subkommandos der CLI (`init`, `generate`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `config get/set/list`, `template list`; geprüft an den Cobra-`Use:`-Zeilen in `internal/adapter/driving/cli/`) sind einer Gruppe eindeutig zugeordnet; `config set`, `add` und `init` stehen in der ändernden Gruppe. Begriffe jetzt einheitlich („Erfolgs- und Zusammenfassungsmeldungen", „Hinweise und Warnungen"). Das Kriterium „nichts verändert" ist am Vertrag prüfbar. <!-- d-check:ignore (Pfad als Zitat) -->
- **F-2 (`--force`-Zusammenfassung):** behoben und konsistent. Ausnahme 4 nennt genau die Zusammenfassung der betroffenen Pfade, die `--force` vor dem Überschreiben ausgibt, und verweist auf `LH-FA-CLI-005A` („gibt immer eine Zusammenfassung der betroffenen Pfade aus") und `LH-FA-INIT-005` („vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden"); beide Quelltexte sind am Lastenheft gegengelesen. Die Zählung stimmt in Absatz, §16 und CHANGELOG. Rest-Beobachtung (nicht blockierend, vom Review als schwächer eingestuft): „Abhängigkeiten prüfen und vor dem Entfernen warnen" (`LH-FA-ADD-007`) und die Purge-Warnungen von `remove` fallen unter „alle anderen Hinweise und Warnungen entfallen"; der Code-Slice behandelt sie ausdrücklich so („soweit sie keine Fehler sind"), der Vertrag nennt sie nicht als Ausnahme. Das ist eine bewusste Setzung des Slice, kein Widerspruch, aber an der Anforderung sichtbar zu halten. <!-- d-check:ignore (Kennung als Zitat) -->
- **F-3 (CHANGELOG):** behoben. Der Eintrag behauptet keinen Ist-Stand mehr, der dem Lastenheft widerspricht: der ältere Absatz („warn items disappear too", „drop every warning") ist entfernt, der verbleibende nennt vier Ausnahmen, ändernde Subkommandos, `doctor` nur Fehler-Einträge, Abfragebefehle behalten ihre Ausgabe und „the implementation follows in a later release". Beobachtung: Der entfernte Absatz beschrieb bereits gelieferten Code (Commit `72006d3`); der Code-Slice trägt den Eintrag in seiner DoD und §3 neu nach. Auflage 4 (kein Release-Cut vor dem Code-Slice) bleibt. Die Zeile ist ein sehr langer Absatz, rein kosmetisch. <!-- d-check:ignore (Kennung als Zitat) -->
- F-4 (LOW, §16-Zeile): behoben; die Zeile nennt die Abfragebefehle und alle vier Ausnahmen. F-5 bis F-7 (INFO): F-5 offen (siehe V-1), F-6 durch die vierte Ausnahme gegenstandslos, F-7 ohne Änderungspflicht.

Befunde dieses Laufs (Planungsdateien, kein Vertragsfehler):
- **V-1 (CR-Slice, alte Zählung):** `slice-cr-cli-quiet-weitere-ausnahmen.md` Z. 1 (Titel „mit drei benannten Ausnahmen") und Z. 39 bis 43 (Ziel: „benennt **drei** Ausnahmen statt einer", Aufzählung ohne `--force`) widersprechen dem CR-Kern (vier) und der Historie-Zeile desselben Slice; Z. 45 „Beide neuen Ausnahmen" meint seit der vierten Ausnahme drei neue. Die Anweisung der Beauftragung („durchgängig vier") ist hier nicht erfüllt. <!-- d-check:ignore (Pfad als Zitat) -->
- **V-2 (`slice-v1-cli-quiet-nur-fehler`, alte Zählung und überholter Prüfpunkt):** §1 Ziel (Z. 44 bis 48) sagt „drei Ausnahmen bleiben sichtbar" und zählt nur drei auf; die Liste der Subkommandos in Ziel und DoD-Punkt 1 lässt `config set` aus. Der „Prüfpunkt — Ergebnis-Ausgaben" nennt die Zuordnung von `config set` und `template`-Unterbefehlen als noch zu entscheiden, obwohl das Lastenheft sie jetzt entscheidet (`config set` ändernd, nur Fehler; `template list` Abfragebefehl). Titel, DoD-Punkt 2, §3, §4 und §6 sagen „vier". Die Abnahmekriterien 1 bis 3 des CR sind getragen (siehe DoD-Punkt 2). <!-- d-check:ignore (Pfad als Zitat) -->
- **V-3 (`open/README.md`):** Z. 19 „drei Ausnahmen bleiben"; Z. 18 nennt für den CR nur Degradations-Fallbacks und Compose-Fortschritt. <!-- d-check:ignore (Pfad als Zitat) -->
- **V-4 (Welle `welle-cli-verbosity`):** §1 und die §4-Zeile des Code-Slices sagen „vier"; die §4-Zeile des CR (Z. 51: „Ausnahmen Degradations-Fallbacks und Compose-Fortschritt") und der Trigger in §2 (Z. 31 bis 33: „Degradations-Fallbacks und Compose-Fortschritt", „vorgeschlagene Fassung") nennen die `--force`-Zusammenfassung nicht und tragen den Vorschlags-Stand. In §1 steht die `--force`-Zusammenfassung vor der Klammer mit `LH-NFA-PERF-002`, die nur zum `up`-Fortschritt gehört (Lesbarkeit). <!-- d-check:ignore (Pfad als Zitat) -->
- **V-5 (`slice-v1-cli-verbosity-letzte-option-gewinnt`):** keine Zählung („vier"/„drei") vorhanden; die Aufzählung der `--quiet`-Wirkungen (Z. 67 bis 70) nennt nur Sandbox-Lockerungs-Hinweise als sichtbar. Informativ, kein Widerspruch (Verweis auf den Code-Slice). <!-- d-check:ignore (Kennung als Zitat) -->

Welle-Reihenfolge:
- CR in `in-progress/`, beide Code-Slices in `open/`; der Code-Slice ist bis zur Schließung des CR blockiert (§4 und §6 sagen es). Konsistent.
- Welle-Closure-Trigger betreffen diesen Slice nicht; Spezifikation im Diff unberührt.
- Roadmap: Wegfall der Zeile „Nichts in Arbeit." bei gelisteter Welle ohne Befund (Zustandsfeld).
- Bekannt, ohne Bezug zum CR: Der Code-Slice ist laut seinem Ausgangsstand teilweise umgesetzt (Fehler-Pegel, `doctor`) und liegt dennoch in `open/`; Lifecycle-Entscheidung des Planners.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0 (lint, test, coverage-gate, docs-check, baseline-verify) |
| `make doc-trace` | pass | selbst gefahren: „83 Anforderung(en), 0 Waise(n)."; `LH-FA-CLI-005`-Zeile trägt `slice-cr-cli-quiet-weitere-ausnahmen` als Beleg | <!-- d-check:ignore (Kennung als Zitat) -->
| `make docs-check` | pass | selbst gefahren vor Anlegen dieser Datei: 257 Dateien, 0 Befunde; nach Anlegen siehe Abschluss |
| `make test-docker` | not run | kein Docker-Pfad berührt (Diff: Lastenheft, CHANGELOG, Planung, Review) |
| `make verify-depguard` | not run | keine Importregel berührt, kein Go-Code im Diff |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Lastenheft 0.9.0 (Absatz, §16), CHANGELOG; Umsetzungsbeleg (Tests) liegt im Code-Slice, hier nicht Ziel |
| [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) | als Quellen der vierten Ausnahme zitiert; Anforderungen selbst unverändert |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-NFA-PERF-002`](../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker) | als Quellen der Ausnahmen 2 und 3 zitiert; unverändert |

Plan-vs-Code-Diff:
- Plan §3 `spec/lastenheft.md` (Anforderung, Versionsfeld, §16-Zeile) — Ist: gesetzt (`0052e1c`, verfeinert in `591acb7`).
- Plan §3 `CHANGELOG.md` — Ist: `[Unreleased]`-Eintrag, englisch, ohne Widerspruch zum Lastenheft.
- Nicht im Plan, im Diff: Wegfall „Nichts in Arbeit." in `roadmap.md`; Welle-, Code-Slice- und Review-Dateien in `591acb7` nachgezogen (V-1 bis V-4: unvollständig).
- Reihenfolge „Lastenheft ändern → Umsetzung" eingehalten (Code unberührt).

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt:
- `make test-docker`, `make verify-depguard` — kein berührter Pfad (siehe Sensors).

Commit / Artefakt:
- Delivery `0052e1c`; Review-Fixes samt Review-Report `591acb7`.

## Verdikt

**Bestanden, mit Auflagen für die Closure.** Der Vertrag (Lastenheft 0.9.0) ist in sich stimmig, deckt sich wörtlich mit dem CR-Kern, alle drei MEDIUM-Findings sind am Lastenheft-Text selbst behoben, und der CHANGELOG widerspricht dem Vertrag nicht. Die Gates sind grün. Offen sind Planungs-Nachzüge und die Closure-Notiz; keine sperrt die Wahrheit des Vertrags, aber V-1 verletzt die Vorgabe „durchgängig vier Ausnahmen" und ist vor `done/` zu beheben.

Auflagen:
1. Vor dem Move nach `done/` (Inhalt zuerst, dann reiner `git mv`, `AGENTS.md` §3.3): DoD-Häkchen setzen, §7 Closure-Notiz mit Steering-Loop-Lerneintrag (Finding-Klassen des Reviews: Normsatz ohne Abgrenzungskriterium, Ausgabepflicht-Zusage ohne Vorrang gegen `--quiet`, CHANGELOG-Ist-Stand gegen Spec-Soll, Historie-Zeile unvollständig), Register-Vermerk und Paarungen belegen. <!-- d-check:ignore (Pfad als Zitat) -->
2. V-1 beheben (Titel, Ziel, „Beide neuen Ausnahmen"); Historie-Zeilen-Entwurf in §1 an die geltende §16-Zeile angleichen oder als Entwurf kennzeichnen.
3. V-2 bis V-4 beheben, bevor der Code-Slice startet: `slice-v1-cli-quiet-nur-fehler` Ziel (vier Ausnahmen, `config set` in der Liste) und Prüfpunkt (Zuordnung entscheidet jetzt das Lastenheft; der Code-Slice pinnt sie nur), `open/README.md` Z. 18 und 19, Welle Z. 31 bis 33 (angenommene statt vorgeschlagene Fassung, vier Ausnahmen) und die CR-Zeile in §4. <!-- d-check:ignore (Kennung als Zitat) -->
4. Kein Release-Cut vor `slice-v1-cli-quiet-nur-fehler` und `slice-v1-cli-verbosity-letzte-option-gewinnt` (CHANGELOG-Präsens). Im Code-Slice die Purge-Warnungen und die Abhängigkeits-Warnung von `remove` (`LH-FA-ADD-007`) ausdrücklich entscheiden und pinnen. <!-- d-check:ignore (Kennung als Zitat) -->
