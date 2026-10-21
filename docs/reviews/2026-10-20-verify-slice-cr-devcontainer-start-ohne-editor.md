# Verification-Report: slice-cr-devcontainer-start-ohne-editor — 2026-10-20

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur; es wurde kein Planungsdokument geändert und nichts committet.

**Gegenstand:** Lastenheft-CR [`slice-cr-devcontainer-start-ohne-editor`](../plan/planning/done/slice-cr-devcontainer-start-ohne-editor.md) (`in-progress/`, kein Produktcode), `HEAD` = `2981c43`, Diff-Basis `9e6933b..HEAD` (5 Commits, 13 Dateien). Maßstab: der freigegebene, unverändert zu übernehmende Wortlaut in [`slice-v1-devcontainer-start-ohne-cli`](../plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md) §7 (Stellen 1 bis 6, Entscheidungsvermerk Punkte 1 bis 3), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (beide `Accepted`).

**Eingang:** Review-Report [`2026-10-20-slice-cr-devcontainer-start-ohne-editor`](2026-10-20-slice-cr-devcontainer-start-ohne-editor.md) (0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO; Fix-Commit `2981c43`) als Eingang, nicht als Beweis; Wortlaut, Fix-Runde und Pläne unten am Ist-Stand nachgelesen.

**Gesamturteil: GO.** Der Lastenheft-Text ist wortgleich zum freigegebenen Vorschlag, trägt Version 0.18.0 samt Historie-Zeile, ist ein Zielbild ohne Prozesszustand, und seit dem Review-Stand `69280df` hat sich `spec/` nicht geändert; die Fix-Runde berührt nur Planungstext und erzeugt keinen Widerspruch. Offen sind allein die Closure-Zeilen des Planners (Freigabe-Vermerk in §7, Closure-Notiz, Register, Roadmap-Zeile bei Closure, Paarungen). Kein HIGH, kein MEDIUM; zwei LOW, drei INFO.

## 1. Sensoren (selbst gefahren, Docker-only)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass (Exit 0, sauberer Arbeitsbaum) | Lauf endet mit `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege` und Exit 0; `docs-check`, `lint`, `test`, `coverage-gate`, `baseline-verify` sind Teil von `gates` |
| `make doc-trace` | pass (Exit 0) | `84 Anforderung(en), 0 Waise(n).`; Zeile [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) `ok` (Belege: der CR-Slice und der Eltern-Slice, ADR-Spalte leer, siehe V-4); [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) tragen den CR-Slice als Beleg |
| `make doc-immutable RANGE="9e6933b..HEAD"` | pass (Exit 0) | `history-range-guard: Range '9e6933b..HEAD' aufgeloest, 5 Commit(s) — OK`; `d-check: 556 Datei(en) geprüft, 0 Befund(e)`; [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) nicht berührt |
| `git diff 69280df..HEAD -- spec/` | leer (erwartet) | keine Zeile Ausgabe: zwischen Review-Stand und `HEAD` ist der Vertragstext unverändert, ein Nachlauf-Review nach [`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) ist für `spec/` nicht ausgelöst |
| `make docs-check` | siehe Abschnitt 7 | nach dem Schreiben dieses Reports erneut gefahren |
| `make test-docker` | **nicht gelaufen** | kein Code, kein Test, kein Docker-Pfad im Diff (`git diff --stat 9e6933b..HEAD`: nur `spec/lastenheft.md`, `docs/plan/planning/`, `docs/reviews/`) |
| `make verify-depguard` | **nicht gelaufen** | keine Importregel berührt, kein Go-Code im Diff |

## 2. DoD-Abgleich (Slice §2, Zeile für Zeile)

Alle Kästchen im Slice sind leer; für einen Slice in `in-progress/` ist das der erwartete Zustand (Häkchen setzt der Planner bei der Closure). Bewertung nach Ist-Stand:

| # | DoD-Zeile | Urteil | Beleg |
|---|---|---|---|
| 1 | Wortlaut im Lastenheft nach Stelle 1 bis 5, freigegeben, Zielbild ohne Prozesszustand und Messstand | **erfüllt** | Abschnitt 3 (Stellen einzeln); Freigabe: Eltern-Slice §7 Entscheidungsvermerk Punkt 1 („freigegeben, unverändert“); Zielbild: Abschnitt 4 |
| 2 | Versionsfeld und Historie (§16) auf die nächste freie Minor-Version, eine Zeile; Priorität V2 in Anforderung und Befehlsübersicht gleich | **erfüllt** | `spec/lastenheft.md`: Versionsfeld `0.18.0` (vorher `0.17.0`, letzte Historie-Zeile war `0.17.0`); Historie-Zeile `0.18.0 · 2026-10-20` nennt alle fünf Stellen wortgleich zu Stelle 6; Anforderung `*Priorität: V2.*`, Befehlsübersicht beide Zeilen `V2` |
| 3 | Freigabe des Projektinhabers belegt **und in §7 vermerkt** | **offen** | Die Freigabe ist belegt (Eltern-Slice §7 Entscheidungsvermerk, Punkte 1 bis 3); §7 des CR-Slice trägt noch „*(Wird bei der Closure gefüllt.)*“. Closure-Aufgabe des Planners, kein Mangel des Wortlauts |
| 4 | `make gates` grün, `make docs-check` grün | **erfüllt** | Abschnitt 1; `docs-check` in `gates` (Exit 0) und nach diesem Report (Abschnitt 7) |
| 5 | Review durchgeführt, Report liegt unter `docs/reviews/` vor | **erfüllt** | [`2026-10-20-slice-cr-devcontainer-start-ohne-editor`](2026-10-20-slice-cr-devcontainer-start-ohne-editor.md) (im Commit `2981c43`, Skill `.harness/skills/reviewer.md`, eigener Kontext); Befunde in der Fix-Runde bearbeitet (Abschnitt 5) |
| 6 | Doku-Update: Roadmap-Zeile „In Arbeit“ bei der Closure; Handbuch und `CHANGELOG.md` unverändert | **teilweise erfüllt, Rest offen** | `roadmap.md` trägt „In Arbeit“ mit dem Verweis auf [`slice-cr-devcontainer-start-ohne-editor`](../plan/planning/done/slice-cr-devcontainer-start-ohne-editor.md) (seit Start); die Änderung bei der Closure steht aus. Handbuch und `CHANGELOG.md` im Diff unberührt (erfüllt) |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | §7 leer; Closure-Aufgabe des Planners |
| 8 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 | **offen** | im Diff kein Eintrag unter `observations/`; Closure-Aufgabe (Kandidat: [`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), Nachlauf-Befund Abschnitt 6) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **formal erfüllt** | alle vier Risiken tragen „Ausgang: weiter offen“ mit Begründung; die Ausgänge „eingetreten/entfallen“ setzt der Planner bei der Closure (Risiko 3 „Wortlaut ändert sich nach dem Review“: nach Abschnitt 1 nicht eingetreten für `spec/`) |
| 10 | Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen | **offen** | hängt an §7 (Steering-Loop-Eintrag); die Folge-Slice-Seite ist belegt: Fundament-, Befehl- und Spezifikations-Slice liegen in `open/` und verweisen auf den CR-Slice in `in-progress/` |

Zusammenfassung: erfüllt 1, 2, 4, 5, formal 9; teilweise 6; offen 3, 7, 8, 10 (alle Closure-Aufgaben des Planners).

## 3. Lastenheft-Text gegen den Vorschlag (Stellen 1 bis 6)

Verglichen wurde `git diff 9e6933b..HEAD -- spec/lastenheft.md` (25 Zeilen, nur `spec/lastenheft.md` in `spec/`) mit dem Vorschlag im Eltern-Slice §7.

| Stelle | Ort im Lastenheft | Urteil |
|---|---|---|
| 1 neue Anforderung | §4.3, nach [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Überschrift, `*Priorität: V2.*`, sieben Spiegelstriche | wortgleich. Zusage-Menge „unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`)“, Colima nur in der Liste „nicht zugesagt“; Exit `10`, `11` (inklusive Startscript-Abbruch) und `12`; `u-boot.yaml`-Grenze („ohne diese Datei … `12`“) und Container-Image-Satz vorhanden. Einzige Abweichung ist die Linkform (Anker statt Pfad plus Anker), im Lastenheft die übliche |
| 2 Befehlsübersicht | [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): zwei Zeilen nach `u-boot egress list`, vor `u-boot template list` | wortgleich, beide `V2` |
| 3 `--quiet` (a) | [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): `devcontainer up`, `devcontainer down` in der Liste der ändernden Subkommandos | wortgleich |
| 3 `--quiet` (b) | Ausnahme 1: Einbindung des Container-Runtime-Sockets in den Hilfs-Container; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in der Verweisliste; Ausnahmen 2 bis 4 unverändert | wortgleich; die Zahl „Vier Ausnahmen“ stimmt weiter (Satz in Ausnahme 1, keine neue Nummer) |
| 4 Enum `command` | [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run): `…, egress, devcontainer oder template` | wortgleich |
| 5 Abgrenzung | [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Spiegelstrich Socket: Zusatz zum Hilfs-Container mit Verweis auf [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | wortgleich |
| 6 Version und Historie | Versionsfeld `0.18.0`; §16 Zeile `0.18.0 · 2026-10-20 · Vereinbarung mit dem Projektinhaber` | wortgleich zur Zeile aus Stelle 6 (zusätzlich Verlinkung der Kennungen) |

Entscheidungen des Vermerks gegen den Text: Priorität V2 (Punkt 2), Colima nicht im Wortlaut (Punkt 3 und „Grenze des Wortlauts“), Exit `11` mit `u-boot.yaml`-Grenze nach [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md): alle drei im Lastenheft so umgesetzt. Nicht zugesagt werden Docker Desktop, rootless Docker, Podman als Host-Engine, entfernte Engine, Linux-`arm64`, Windows; keine Aussage geht über die Zusage-Menge hinaus. Kein Satz des Zusatzes widerspricht Ausnahme 2 (Fallback-Ausweisung bleibt unter `--quiet`).

## 4. Zielbild-Regel (AGENTS.md §3.4)

Gesucht in den hinzugefügten Zeilen von `spec/lastenheft.md` nach `ADR-`, `slice-`, `welle`, `roadmap`, `Carveout`, `Messstand`, Commit-Hashes und Prozesszuständen („ausstehend“, „geplant“, „noch nicht“): keine Fundstelle. Die Verweise der neuen Zeilen zeigen nur innerhalb des Lastenhefts (`#lh-…`). „Linux-Host `amd64`“ ist eine Zusage-Menge, kein Messstand. Der Lastenheft-Text nennt den CR-Slice nicht.

## 5. Plan-Tabelle gegen `git diff --stat` (AGENTS.md §3.10) und Fix-Runde

`git diff --stat 9e6933b..HEAD` (13 Dateien) gegen die Tabelle in §3 des CR-Slice:

| Datei im Diff | Tabellenzeile | Urteil |
|---|---|---|
| `spec/lastenheft.md` | Zeile 1 | ok |
| `in-progress/roadmap.md` | Zeile 2 („bei Start und bei der Closure“) | ok |
| `open/README.md` | Zeile 3 | ok |
| [`2026-10-20-slice-cr-devcontainer-start-ohne-editor`](2026-10-20-slice-cr-devcontainer-start-ohne-editor.md) (Review-Report) | Zeile 4 | ok |
| zwei bestehende Reports des Eltern-Slice (je eine Zeile, Linknachzug) | Zeile 4 („Reports des Eltern-Slice bleiben, der Verschiebe-Werkzeugaufruf zieht nur Verweise nach“) | ok |
| der CR-Slice selbst | Zeile 5 | ok |
| [`slice-v1-devcontainer-start-ohne-cli`](../plan/planning/done/slice-v1-devcontainer-start-ohne-cli.md), `welle-devcontainer-start-ohne-editor.md`, [`slice-v1-devcontainer-start-fundament`](../plan/planning/done/slice-v1-devcontainer-start-fundament.md), [`slice-v2-spec-devcontainer-start`](../plan/planning/open/slice-v2-spec-devcontainer-start.md) | Zeile 6 („nur Verweise“) | ok (Diffs gelesen: nur Pfad `open/` → `in-progress/`; in der Welle zusätzlich die Aussage zur Versionsnummer, von der Zeile genannt) |
| [`slice-v1-cli-devcontainer-up-down`](../plan/planning/done/slice-v1-cli-devcontainer-up-down.md) | Zeile 6 („update (nur Verweise)“) | **Abweichung**, siehe V-1 |
| [`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md) | Zeile 7 | ok |
| Dateien „unverändert“ (`spec/spezifikation.md`, `spec/architecture.md`, Handbuch, `CHANGELOG.md`, `internal/`, `cmd/`) | Zeile 8 | im Diff nicht enthalten, bestätigt |

Fix-Runde (`2981c43`, 3 Dateien): (a) Plan-Tabelle ergänzt, (b) S4-Test im Befehl-Slice [`slice-v1-cli-devcontainer-up-down`](../plan/planning/done/slice-v1-cli-devcontainer-up-down.md): „im Sandbox-Profil überlebt ein Marker im benannten Volume `down`, ein folgendes `up` und `--recreate`“ mit Wächter (`docker volume rm`), (c) der Review-Report. Der neue Satz deckt den Halbsatz „er entfernt weder den Container noch seine Volumes“ aus [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und widerspricht keiner anderen Stelle desselben Slice (DoD-Zeile „benanntes Volume mit Marker überlebt“ und die Rückführung „`--recreate` entfernt das Volume“ bleiben verträglich, vgl. V-3). Die Fix-Runde berührt `spec/` nicht (Abschnitt 1).

## 6. Findings

```text
LOW docs/plan/planning/done/slice-cr-devcontainer-start-ohne-editor.md · Plan-Tabelle §3, Zeile „update (nur Verweise)“ - die Zeile führt done/slice-v1-cli-devcontainer-up-down.md nur als Verweis-Nachzug
V-1  Klasse: plan-tabelle-ohne-diff-abgleich · Quelle: AGENTS.md §3.10
Befund: Der Fix-Commit 2981c43 hat in done/slice-v1-cli-devcontainer-up-down.md inhaltlich den Test S4 geändert (Marker im Volume nach down). Die Tabelle, die der Fix-Commit selbst ergänzt hat, nennt für diese Datei „nur Verweise“; der Abgleich „nach der letzten Fix-Runde“ ist damit für diese Datei nicht nachgeführt.
Risiko: gering; die Änderung ist im Diff sichtbar, die Tabelle ist als Maßstab für den nächsten Leser unvollständig.
Verifizierbar: git diff --stat 9e6933b..HEAD gegen §3; Behebung: Planner (Zeile um „S4: Marker im Volume nach down“ ergänzen), bei der Closure.
```

```text
LOW docs/reviews/2026-10-20-slice-cr-devcontainer-start-ohne-editor.md · Zeile 23 (erstes LOW-Finding, „Befund“) - zerschossener Link-Text
V-2  Klasse: report-darstellungsfehler · Quelle: harness/review.md (Output-Schema, lesbarer Report)
Befund: Der Dateiname der messstaende-Slice-Datei stand im Review-Report als ineinander verschachtelte Link-Syntax; `make docs-check` meldete keinen Befund, der Satz war aber nicht lesbar. Behoben: ein einfacher Link auf die Datei.
Risiko: keines für den Vertragstext; Lesbarkeit des Eingangs-Reports.
Verifizierbar: Review-only (Reviewer-Artefakt; der Verifier repariert nichts).
```

```text
INFO docs/plan/planning/done/slice-v1-cli-devcontainer-up-down.md · S4 und Rückführung - „--recreate: benannte Volumes bleiben erhalten“ ist nur im Sandbox-Profil gepinnt
V-3  Rest des zweiten LOW-Findings des Reviews: Das Standardprofil hat kein benanntes Volume; die Zusage im Lastenheft nennt kein Profil. Die Rückführung („--recreate entfernt das Volume“) nennt die down-Seite nicht ausdrücklich. Kein Widerspruch, kein Test fehlt für ein vorhandenes Volume; Hinweis für den Befehl-Slice.
```

```text
INFO spec/lastenheft.md / doc-trace · ADR-Spalte der Zeile LH-FA-DEV-010 leer
V-4  Wie im Review (INFO 4) bestätigt: Die Zeile ist `ok` über den Slice-Beleg. Kein Handlungsbedarf, solange keine Folge-ADR das Schärft-Feld setzt.
```

```text
INFO docs/plan/planning/done/slice-cr-devcontainer-start-ohne-editor.md · DoD „Freigabe … in §7 vermerkt“
V-5  Klasse: vertragsaenderung-nach-review · Die Freigabe ist im Eltern-Slice §7 belegt; der Vermerk im CR-Slice §7 entsteht mit der Closure. Nachlauf-Review: nicht ausgelöst, solange `git diff 69280df..HEAD -- spec/` leer bleibt (Stand: leer). Wird `spec/lastenheft.md` vor oder bei der Closure noch angefasst, ist ein Nachlauf-Review über genau diesen Diff nötig, bevor der Slice nach done/ geht.
```

Kein HIGH, kein MEDIUM.

## 7. Geprüft ohne Befund

- Wortlaut-Treue Stellen 1 bis 6 (Abschnitt 3), Entscheidungen Priorität V2, Zusage-Menge Linux-Host `amd64` ohne Colima, Exit `11` mit `u-boot.yaml`-Grenze
- Zielbild-Regel im Lastenheft-Diff (Abschnitt 4)
- Version 0.18.0 durchgängig: Verweise in der messstaende-Slice-Datei (Titel, Historie-Frage), `open/README.md` und Welle nennen „nach 0.18.0“ beziehungsweise „0.18.0 vergeben“; keine Stelle im Planungsbaum behauptet mehr, 0.18.0 sei für den Messstände-CR vorgesehen
- ADR-Immutabilität (`doc-immutable`), `spec/spezifikation.md` und `spec/architecture.md` unberührt

## 8. Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle.

## 9. Nicht ausgeführt / Grenzen

- `make test-docker`, `make verify-depguard`: kein Code-, Docker- oder Importregel-Pfad im Diff.
- Die Tragfähigkeit der Zusage (Linux `amd64`, Standardprofil mit UID ungleich 1000) gegen reale Läufe: Aufgabe der Docker-Tests des Befehl-Slice, nicht dieses CR; nicht Gegenstand dieser Verifikation.
- `make docs-check` nach dem Schreiben dieses Reports: Ergebnis im Handback an den Aufrufer.

Commit / Artefakt: `HEAD` = `2981c43`, kein Commit durch den Verifier.
