# Verification-Report: [`slice-cr-devcontainer-shell-und-zweites-up`](../plan/planning/in-progress/slice-cr-devcontainer-shell-und-zweites-up.md) — 2026-10-21

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Lastenheft-CR richtig gebaut, gemessen an Plan, DoD, den Entscheidungen des Projektinhabers und den Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Außer dieser Datei ist nichts geändert, nichts committet, nichts gepusht. Jede Mutation ist zurückgesetzt; `git status` war danach jeweils leer.

**Gegenstand:** `git diff 770113e8..5f643a81`. Commits des CR: `5e54f672`, `961e22cd`, `a5ede8a3`, `9e6c47d4`, `d01fc191`, `d20e4dd6`, `fc3151f0`, `a46dcf5b`. Fremde Planner-Commits in der Range, nur auf Verträglichkeit gelesen: `6289f5c6` und `5f643a81`. `HEAD` = `5f643a81`, der Arbeitsbaum war vor dieser Datei sauber. Das Lastenheft ändern nur `d01fc191`, `fc3151f0` und `a46dcf5b`; Code, Spezifikation, Sicht, ADRs, Handbuch, READMEs und `CHANGELOG.md` liegen nicht im Diff.

**Eingang (gelesen, aber nicht als Beweis genommen):**

- Slice §1 bis §8 samt *Umsetzung*, *Fix-Runde* und *Fix-Runde 2*, §7 mit den Entscheidungen 1 bis 12 und den Übergabe-Zeilen B1 bis B19 und A1 bis A7.
- Review-Report [`2026-10-21-slice-cr-devcontainer-shell-und-zweites-up`](2026-10-21-slice-cr-devcontainer-shell-und-zweites-up.md) (0 HIGH, 2 MEDIUM, 5 LOW, 4 INFO), erstes Nachlauf-Review [`…-nachlauf`](2026-10-21-slice-cr-devcontainer-shell-und-zweites-up-nachlauf.md) (0 HIGH, 2 MEDIUM, 2 INFO, NO-GO bis zur Entscheidung) und zweites Nachlauf-Review [`…-nachlauf-2`](2026-10-21-slice-cr-devcontainer-shell-und-zweites-up-nachlauf-2.md) (0 HIGH, 1 MEDIUM, 1 LOW, 1 INFO, GO für den Verifier).
- [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) und [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), vollständig.
- Lastenheft: Kopf, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und §16.
- Die Pläne der Folge-Slices in `open/` und die Welle-Datei, jeweils auf `5f643a81`.
- Entscheidungen des Projektinhabers vom 2026-10-21 (gesetzt, nicht neu verhandelt): Wortlaut der Stellen 1 bis 5, 6a und 8 freigegeben, 6b und 7 mitgenommen; 6b-Satz als letzter Satz des Absatzes und Historienzeile mit Klammer zu Stelle 7 bestätigt; Lastenheft enger (Prüfung vor der Shell, Socket des Hosts, dann Beitrag des Befehls); keine Lastenheft-Änderung für den Fall „Prüfung gelingt, Sitzung beginnt nicht“; Terminal-Test per `script`; Hilfe von `shell` ohne Kennung; `d20e4dd6` nicht umgeschrieben. Das Gespräch selbst kann ich nicht nachprüfen; Beleg sind die Vermerke in §7 und der Auftrag.

**Gesamturteil: GO für die Closure.** Ich habe keine DoD-Verletzung gefunden. Die drei Liefer-Punkte der DoD (Wortlaut, Version und Historie, Freigabe) sowie die Punkte Gates und Review sind erfüllt; das Doku-Update ist bis auf das Zurücksetzen der Roadmap erfüllt. Offen sind die Closure-Pflichten: Closure-Notiz, Register, Ausgang je Risiko und die drei Paarungen. Der Lastenheft-Text ist Zeichen für Zeichen die Summe aus den freigegebenen Stellen und den beiden Fix-Runden; die zwei Fix-Runden hatten je ein Nachlauf-Review über genau ihren Diff. Die drei Befunde V-1 bis V-3 sind LOW und betreffen Plan-Buchhaltung des CR-Slice; der Implementer zieht sie in den Closure-Commits nach. V-4 bis V-9 sind INFO. Ändert die Closure Lastenheft-Text, ist ein weiteres Nachlauf-Review und ein neuer Verifier-Lauf fällig.

## Verification Evidence

Scope:

- Slice: [`slice-cr-devcontainer-shell-und-zweites-up`](../plan/planning/in-progress/slice-cr-devcontainer-shell-und-zweites-up.md), Lastenheft `0.19.0` → `0.20.0`, Teil von [`welle-devcontainer-shell-und-zweites-up`](../plan/planning/welle-devcontainer-shell-und-zweites-up.md).
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging); Träger [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) und [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) (beide `Accepted`, unverändert).
- Artefakte: `spec/lastenheft.md`, Slice, Roadmap, `open/README.md`, Messstände-CR, Welle-Datei, drei Review-Reports; fremd (Planner): die Pläne der Folge-Slices.

### DoD-Abgleich (§2 des Slice, Stand `5f643a81`)

- [x] **Wortlaut im Lastenheft** — erfüllt. Basis `770113e8` plus die freigegebenen Stellen plus die beiden Fix-Runden ergibt `spec/lastenheft.md` auf `5f643a81` zeichengleich (Einzelnachweis 1). Kein Prozesszustand, kein Messstand, kein Verweis nach unten (Einzelnachweis 3). **Für die Closure:** Die DoD verlangt die Fundstellen „in §7“; sie stehen in §3 *Umsetzung* (V-8).
- [x] **Versionsfeld und Historie** — erfüllt. Zeile 9 trägt `0.20.0`; die Historie-Zeile `0.20.0` (Zeile 1592) nennt nur freigegebene Stellen; `V2` steht in Zeile 1140 und in Zeile 490 gleich.
- [x] **Freigabe, kein Wortlaut nach dem Review ohne Nachlauf-Review** — erfüllt. §7 Punkte 1 bis 6 vermerken die Freigabe je Stelle (6b und 7 einzeln). Die Fix-Runde `fc3151f0` änderte zwei Lastenheft-Zeilen, das erste Nachlauf-Review las `git diff d20e4dd6..fc3151f0 -- spec/lastenheft.md`; die Fix-Runde 2 `a46dcf5b` änderte zwei Zeilen, das zweite Nachlauf-Review las `git diff fc3151f0..a46dcf5b -- spec/lastenheft.md`. Beide Diffs habe ich nachgezogen: je 2 Einfügungen, 2 Löschungen, dieselben zwei Zeilen. Seit `a46dcf5b` ist `spec/` unverändert.
- [x] **`make gates`, `make docs-check`, `make doc-trace`** — selbst gefahren auf `5f643a81`, alle Exit 0 (siehe *Sensors*). [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) stehen in der Matrix auf `ok`. Für [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) trägt der Vergleich in Einzelnachweis 4. **Für die Closure:** Läufe nach den Closure-Commits mit Exit-Code in §7.
- [x] **Review** — erfüllt; drei Reports liegen im Repo. Die Ausgänge der Findings des zweiten Nachlauf-Reviews stehen noch nicht im Slice (V-3).
- [ ] **Doku-Update** — Teil erfüllt. Die Roadmap-Zeile „In Arbeit“ ist gesetzt (`9e6c47d4`); Handbuch, READMEs, `CHANGELOG.md` und ADR-Index sind unverändert (leerer Diff, Einzelnachweis 8). Offen (Closure): „Nichts in Arbeit.“ bei der Closure.
- [ ] **Closure-Notiz mit Steering-Loop-Lerneintrag** — offen (Closure). §7 trägt nur die Entscheidungen und die Übergabe-Tabelle.
- [ ] **Beobachtungs-Register** — offen (Closure). `observations/` liegt nicht im Diff. Eingaben in der *Prüfliste*.
- [ ] **Jedes Risiko aus §6 mit Ausgang** — offen (Closure). §6 hat sieben Risiken, eines steht auf „entfallen“, sechs auf „weiter offen“. Vorschläge in der *Prüfliste*.
- [ ] **Drei Paarungen** — offen (Closure der Welle).

### Einzelnachweise (selbst gelesen oder gemessen, nicht aus dem Slice übernommen)

**1. Wortlaut, zeichengenau.** Das Skript lag im Scratchpad und ist ein Text-Werkzeug, kein Gate. Es liest die *Alt*- und *Neu*-Wortlaute der Stellen 1 bis 8 aus §1 des Slice (Stand `5f643a81`, Verweise `../../../../spec/lastenheft.md#` als `#`), prüft, dass jeder *Alt*-Wortlaut im Lastenheft auf `770113e8` genau einmal steht, ersetzt ihn durch den *Neu*-Wortlaut, fügt die neue Zeile der Befehlsübersicht, den Spiegelstrich zu `shell` und die Historie-Zeile ein und vergleicht das Ergebnis mit `spec/lastenheft.md` auf `5f643a81`. Angewandt habe ich nur die Abweichungen, die der Slice selbst benennt: das Datum, die Auflösung der Klammer „(und, mit 6b, …)“, die Klammer zu Stelle 7 in der Historie-Zeile (§3 *Plan vor Code*, Punkt 2; bestätigt in §7 Punkt 7), und die Sätze der beiden Fix-Runden (§1 Stellen 2 und 8). Ergebnis: **26 Prüfungen, 0 Fehlschläge; das Ergebnis ist zeichengleich mit dem Lastenheft auf `5f643a81`.** Das schließt jede weitere, ungeplante Änderung im Lastenheft aus.

Drei Gegenkontrollen zum Skript. Erstens die Vorfrage, ob der Vorschlag im Plan seit `770113e8` geändert wurde: Im Abschnitt *Lastenheft-Wortlaut-Vorschlag* unterscheiden sich `770113e8` und `5f643a81` nur in Überschriften, Freigabe-Vermerken und den beiden Fix-Runden-Notizen; der *Neu*-Text und der Vorschlag der Historie-Zeile sind gleich. Zweitens die Sätze, die §7 Punkte 8 und 10 als Entscheidung des Projektinhabers wiedergeben: Die neuen stehen im Lastenheft genau einmal, die ersetzten nicht mehr. Drittens der Satz zur Klasse `11`: Er ist zeichengleich mit der Entscheidung c von [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md).

| Stelle | Fundstelle im Lastenheft | Ergebnis |
|---|---|---|
| 1 — Fehlerklassen | Zeile 497 | zeichengleich; Satz wörtlich Entscheidung c von [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) |
| 2 — Einleitung und Spiegelstrich zu `shell` | Zeilen 490 und 494 | zeichengleich, mit den Sätzen der Fix-Runde 1 zur `12` und der Fix-Runde 2 zum Socket; der Spiegelstrich steht nach dem zu `down` und vor dem zum Socket |
| 3 — Socket-Spiegelstrich | Zeile 495 | zeichengleich, Rest des Spiegelstrichs unverändert; „Dauer des Aufrufs“ steht nur noch hier, neben „Dauer des Starts“ |
| 4 — Befehlsübersicht | Zeile 1140 | zeichengleich, genau einmal; nach `down`, vor `template list`; alle Tabellenzeilen 81 Zeichen breit |
| 5 — Abgrenzung des Hilfs-Containers | Zeile 448 | zeichengleich |
| 6a — Ausnahme 1 | Zeile 191 | zeichengleich |
| 6b — Liste der ändernden Befehle und Satz | Zeilen 189 und 196 | zeichengleich; der Satz ist der letzte des Absatzes, „Ebenso“ bezieht sich weiter auf die Abfragebefehle |
| 7 — Zusage-Menge | Zeile 498 | zeichengleich |
| 8 — Versionsfeld und Historie | Zeilen 9 und 1592 | `0.20.0`, Spaltenbreite unverändert; Historie-Zeile nach `0.19.0`, Datum `2026-10-21` (Commit-Datum von `d01fc191`), Verweis „Vereinbarung mit dem Projektinhaber“ |
| übriges Lastenheft | — | Der Diff hat sechs Hunks (Version, Zeilen 189 bis 196, Zeile 448, Zeilen 490 bis 498, Zeile 1140, Zeile 1592); außerhalb der Stellen 1 bis 8 ist keine Zeile geändert |

**2. Innere Widerspruchsfreiheit und Bezug zu den ADRs.** Ich habe den Spiegelstrich zu `shell` im Zusammenhang mit den Nachbarstrichen gelesen, nicht nur als Diff.

- Start und Sitzung widersprechen sich nicht: Der Hilfs-Container und der Socket gehören zur Startphase („bei `shell` für die Dauer des Starts“, Zeile 495); die Sitzung läuft „ohne Hilfs-Container“, und der Befehl bindet für sie keinen Socket ein (Zeile 494). Der Nachbarsatz „der gestartete Devcontainer erhält den Socket durch diesen Befehl nicht“ (Zeile 495) hat dieselbe Lesart wie der neue Satz: Beide sagen zu, was der Befehl einbindet.
- Die Fehlerklassen passen zusammen: Zeile 494 verweist für den Start auf die Fehlerklassen von `up`; Zeile 497 führt `11` für den Abbruch eines Startscripts und, neu, für den früheren Abbruch; die `12` der Prüfung vor der Shell steht nur in Zeile 494 und kollidiert mit keiner anderen Zeile. Die Voraussetzung der Statusdatei in Zeile 498 gilt für „der Befehl“.
- [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): `devcontainer shell` unter den ändernden Befehlen, Ausnahme 1 nennt den Hilfs-Container von `shell`, der Satz zur Shell sagt, dass `--quiet` nur auf die Ausgabe des Starts wirkt. Ausnahme 2 (Fallback-Ausweisung) gilt über „wie `up`“ auch für den Start von `shell`.
- [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): `0` und `12` liegen im Wertebereich, kein reservierter Code wird genutzt.
- Gegen [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md): Der Spiegelstrich geht an keiner Stelle über die Entscheidungen 1 bis 3 hinaus. Er ist an zwei Stellen enger: Die ADR sagt „ohne Socket“, das Lastenheft „bindet der Befehl keinen … Socket des Hosts ein“ (Entscheidung des Projektinhabers 10); die ADR bindet die `12` an die Probe, das Lastenheft ebenso, nennt den Fall „Probe gelingt, Sitzung beginnt nicht“ aber nicht. Die ADR ist unverändert, das Lastenheft geht vor. Der Satz zu den Umgebungsvariablen nennt „die erzeugte Konfiguration“ wie der bestehende Spiegelstrich zu den Zugangsdaten.
- Gegen [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md): Stelle 1 ist wörtlich Entscheidung c; Stelle 7 ist Grenze (ii) von Entscheidung 1.

**3. Referenzrichtung** ([`AGENTS.md`](../../AGENTS.md) §3.4). Die zwölf hinzugefügten Zeilen in `spec/lastenheft.md` enthalten kein `slice-`, `welle`, `ADR`, `adr/`, `docs/plan`, `BEO-`, `roadmap`, `carveout`, `tranche`, keine `SPEC-`- oder `ARC-`Kennung und keinen Pfad mit `../`. Alle Link-Ziele sind Anker der eigenen Datei (`#lh-…`). Prozesszustände („Vorschlag“, „ausstehend“, „in Arbeit“) kommen in den hinzugefügten Zeilen nicht vor.

**4. [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (außerhalb der Belegsicht von `make doc-trace`).** Die Anforderung führt die Belegsicht nicht (nur `LH-FA-*` und `LH-QA-*`). Der Vergleich trägt hier: Stelle 4 ist zeichengleich (Einzelnachweis 1), die Zeile steht an der geforderten Stelle, sie trägt `V2` wie die Einleitung von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), und die Tabelle ist nicht verrutscht (alle Zeilen 81 Zeichen). Das Gate fängt einen verfälschten Anker in der Historie-Zeile (Gegenprobe G-2).

**5. Gegenproben des Docs-Gates** (die DoD verlangt sie; je ein Eingriff in `spec/lastenheft.md`, ein `make docs-check`, danach `git checkout -- spec/lastenheft.md`):

| # | Eingriff | Ergebnis | Ausgabe (Auszug) |
|---|---|---|---|
| G-1 | im Spiegelstrich zu `shell` ein Link auf die ADR-Datei von [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) | **rot, Exit 2** | `spec/lastenheft.md:494 … matrix-forbidden … Referenz contract-spec → adr ist nicht erlaubt` |
| G-2 | Anker der Historie-Zeile `0.20.0` auf [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) verfälscht (`…-befehlx`) | **rot, Exit 2** | `spec/lastenheft.md:1592 … anchor-missing … Anker entspricht keinem Heading-Slug` |
| G-3 | im selben Spiegelstrich ein Link auf die technische Spezifikation ([`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)) | **rot, Exit 2** | `spec/lastenheft.md:494 … matrix-forbidden … Referenz contract-spec → tech-spec ist nicht erlaubt` |

Nach jedem Eingriff war `git status` leer (`git diff --stat` zeigte vorher je 1 Einfügung, 1 Löschung). G-1 bis G-3 entsprechen den Mutationen M2, M1 und M3 des Implementers, mit gleichem Ergebnis. M0 (Stand `a5ede8a3`, Ruhe-Marker bei Slice in `in-progress/`) habe ich nicht nachgefahren; die Voraussetzung stimmt (auf `a5ede8a3` steht „Nichts in Arbeit.“ in der Roadmap, und der Slice liegt in `in-progress/`), `9e6c47d4` behebt sie.

**6. Fix-Runden: Ausgänge im Slice.**

*Review-Report 1* (`d20e4dd6`, Fix-Runde `fc3151f0`):

| Finding | Ausgang im Slice | Belegt |
|---|---|---|
| F-1 MEDIUM (Terminal-Test ohne Eigentümer) | Übergabe B3, Eigentümer Fundament; §6 nennt ihn | Fundament DoD 2 (Test unter `script`, Gegenprobe) und Befehl DoD 1 (Terminal-Erkennung) tragen es auf `5f643a81` |
| F-2 MEDIUM (Satz aus 6b ohne Pin) | Übergabe B15; §6 führt 6b, Stelle 7 und die Prüfung | Befehl DoD 1 und §3 (`devcontainer_test.go`, Sitzungsausgabe unter `--quiet` mit Wächter) |
| F-3 LOW (abgelaufene Marker in Folge-Plänen) | an den Planner | `6289f5c6`: Die Sätze „solange der Projektinhaber nicht anders entscheidet“, „ist nicht gelesen“ und „nur nach Entscheidung des Projektinhabers“ stehen in Fundament- und Befehl-Slice nicht mehr (auf `770113e8` je zweimal, auf `5f643a81` nullmal) |
| F-4 LOW (Aufträge nur im CR) | Übergabe A1 bis A7 | Spezifikations-Plan trägt A1 bis A3, Befehl-Plan A4 bis A7 |
| F-5 LOW, F-6 LOW (Wortlaut) | im Lastenheft behoben | Einzelnachweis 1 |
| F-7 LOW (Struktur-IDs in `d20e4dd6`) | nicht umgeschrieben, §3 *Fix-Runde* und §7 Punkt 9 | bestätigt (V-7); Sensor-Lücke für die Closure vorgemerkt |
| F-8, F-10, F-11 INFO | Übergabe B9, A7, B10 | in den Plänen als Hinweis bzw. Risiko (Einzelnachweis 9) |
| F-9 INFO (Zitat nicht wörtlich) | keine Änderung, Klasse für die Closure | bestätigt: Der Doc-Kommentar von `Stop` lautet „It neither removes the container nor starts a helper container.“; das Zitat „starts no helper container“ in §1 steht unverändert |

*Erstes Nachlauf-Review* (`fc3151f0`, Fix-Runde 2 `a46dcf5b`): F-1 MEDIUM → Lastenheft geändert (Entscheidung 10), Nachlauf 2 las den Diff. F-2 MEDIUM → keine Lastenheft-Änderung, §3 und §7 B8 berichtigt, neue Übergabe B18. F-3 INFO → B19. F-4 INFO → außerhalb, eigener Bugfix (§7 Punkt 12). Alle vier sind im Slice belegt.

*Zweites Nachlauf-Review* (`a46dcf5b`): F-1 MEDIUM, F-2 LOW, F-3 INFO haben **keinen Eintrag im Slice**. Ihre Ausgänge liegen im Planner-Commit `5f643a81` (siehe Einzelnachweis 7), nicht im Slice (V-3).

**7. Fremde Planner-Commits `6289f5c6` und `5f643a81`: Verträglichkeit.** Beide berühren nur `docs/plan/planning/` und, mit `5f643a81`, den Nachlauf-Report 2; `spec/`, Code und ADRs bleiben unverändert (numstat). Die Lastenheft-Zitate in den Plänen sind zeichengleich mit dem Lastenheft (geprüft: die Sätze zur `12`, zum Socket, zur Zusage-Menge, zu [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); je genau eine Fundstelle im Lastenheft und im Spezifikations-Plan). `5f643a81` ändert im CR-Slice genau eine Zeile, das Ziel in §1, und es stimmt mit dem Lastenheft überein („für die Sitzung bindet der Befehl keinen Socket des Hosts ein“). Die vier Plan-Aussagen aus Nachlauf-2 F-1 sind berichtigt: Welle-Datei (Ziel und Closure-Kriterium 2), Befehl-Slice (Ziel, DoD Doku, Handbuch-Zeile A4, Risiko „Bedingung im Handbuch ohne Test“), CR-Ziel. Nachlauf-2 F-2 ist berichtigt: Der Spezifikations-Plan nennt für den Fall „Prüfung gelingt, Sitzung beginnt nicht“ nur den nicht trennbaren Rest als Grenze. Nachlauf-2 F-3 ist aufgenommen: Fundament §4 schickt den Fall „Container endet zwischen Probe und Sitzung“, wenn nicht trennbar, an den Projektinhaber. Eine Suche nach „ohne Socket“ und „Sitzung ohne/kein Socket“ in Welle, `open/` und `in-progress/` findet außerhalb des CR-Slice nur die neuen Formulierungen (Beitrag des Befehls), den Verweis auf die Überschrift der ADR und den Argumentlisten-Pin des Adapters (Fundament DoD 1: „ohne Hilfs-Container, ohne Socket, ohne `--privileged` und ohne `-v`“). Der Pin beschreibt die Argumentliste, nicht den Zustand der Sitzung.

**8. Plan-Tabelle gegen `git diff --stat 770113e8..5f643a81`** ([`AGENTS.md`](../../AGENTS.md) §3.10, Stand nach der letzten Änderung). Der Diff hat 14 Zeilen in 13 Dateien (der Slice steht zweimal, Löschung unter `open/` und Neuanlage unter `in-progress/`; über die ganze Range liegt der Rename unter der Ähnlichkeitsschwelle, der Move ist der reine Commit `961e22cd`, `100 %`). Jede Datei hat einen Eintrag:

| Datei | Plan-Zeile | Herkunft |
|---|---|---|
| `spec/lastenheft.md` | ja | CR: `d01fc191`, `fc3151f0`, `a46dcf5b` |
| `in-progress/roadmap.md` | ja | CR: `9e6c47d4` |
| der Slice (Rename) | „diese Datei“ | CR; fremd: `5f643a81` (eine Zeile) |
| `open/README.md` | ja | CR: `d01fc191`, `a5ede8a3`; fremd: `5f643a81` (Zeile des Befehl-Slice) |
| Messstände-CR | ja | CR: `d01fc191` |
| Urteil-nach-dem-Zustand-Slice | übrige Slice-Dateien (nur Verweise) | CR: `a5ede8a3` |
| Fundament-Slice | dito | CR: `a5ede8a3`; fremd: `6289f5c6`, `5f643a81` |
| Spezifikations-Slice | dito | CR: `a5ede8a3`; fremd: `6289f5c6`, `5f643a81` |
| Befehl-Slice | dito | nur fremd: `6289f5c6`, `5f643a81` |
| Welle-Datei | ja | CR: `5e54f672`, `a5ede8a3`; fremd: `6289f5c6`, `5f643a81` |
| Review-Report und erstes Nachlauf-Review | `docs/reviews/` | CR: `fc3151f0`, `a46dcf5b` |
| zweites Nachlauf-Review | die Zeile `docs/reviews/` nennt es nicht | nur fremd: `5f643a81` |

Die als „unverändert“ geführten Pfade liegen nicht im Diff (`git diff --stat` über `spec/spezifikation.md`, `spec/architecture.md`, `docs/user`, `docs/maintainer`, beide READMEs, `CHANGELOG.md`, `internal`, `cmd`, `docs/plan/adr`, `Makefile`, `repo.mk`, `Dockerfile`, `.golangci.yml`, `.github`, `.d-check.yml`, `harness`, `AGENTS.md`, `.harness`: leer). Die gelesenen Aussagen der Zeile „unverändert“ habe ich nachgelesen: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeile 777 trägt den Satz „welcher Exit-Code für ein `failed` aus einem früheren Lauf gilt, ist nicht zugesagt“; Handbuch Zeile 818 *Im Container arbeiten*, Zeile 874 *Erneutes `up` nach einem Fehler oder Abbruch*; `newDevcontainerCommand` hängt in `devcontainer.go` Zeile 81 genau `up` und `down` an. Die Zeilennummern der Aufträge A1 bis A7 stimmen: Spezifikation 648, 936, 1181; Handbuch 770 bis 775, 811, 699, 742, 841; JSON-Kontrakt 1328 und 1329. Die Zahlen des Implementers stimmen ebenfalls (9, 11 und 13 Zeilen auf `d01fc191`, `fc3151f0`, `a46dcf5b`). Zum Stand `5f643a81`: V-2.

**9. Übergabe-Tabelle (§7), Spalte „im Plan“, gegen die Pläne in `open/` auf `5f643a81`.** Der Slice führt die Spalte auf dem Stand `fc3151f0` und sagt selbst, dass er den Nachzug des Planners nicht nachgeprüft hat. Das Ergebnis meiner Prüfung: 13 von 26 Zeilen sind falsch oder unvollständig, 12 stimmen, B17 ist „—“.

| Zeile | Spalte laut Slice | Stand `5f643a81` | Fundstelle im Plan |
|---|---|---|---|
| B1, B2 | ja | stimmt | Urteil-nach-dem-Zustand-Slice, DoD 1 und DoD 2 (1) mit Wächter |
| B3 | nein (F-1) | **ja** | Fundament DoD 2 (Terminal-Pfad unter `script`, Gegenprobe, Eigentümer) und Befehl DoD 1 (zweiter Pin); Welle §5 Frage 2 |
| B4, B5, B6, B7, B8 | ja | stimmen | Fundament DoD 1/2; Befehl DoD 1/2 und §3 |
| B9 | nein (Hinweis) | **Hinweis im Plan, kein Pin in einer DoD** | Befehl §3 (`devcontainer_start_test.go`, als Vorschlag des Planners), Spezifikations-Plan §1 und §6 |
| B10 | nein (Hinweis) | **Risiko mit Ausgang „weiter offen“** | Befehl §6 (`--yes`, `--no-interactive`), Welle §3 Punkt 3 (vi) |
| B11, B12, B13, B14 | ja | stimmen | Fundament DoD 1; Befehl DoD 1/2 und §3 |
| B15 | nein (F-2) | **ja** | Befehl DoD 1 und §3 (`devcontainer_test.go`: Sitzungsausgabe unter `--quiet`, Wächter) |
| B16 | ja (Befehl DoD 1) | stimmt; der Anteil „Verifier dieses CR“ ist mit Einzelnachweis 4 erledigt | Befehl DoD 1 |
| B17 | — | — | — |
| B18 | nein (Auftrag an den Planner) | **ja** | Fundament §1, DoD 1, DoD 2, §3, §4, §5, §6; Spezifikations-Plan §1, DoD 1, §6. Zusatz des Planners, der in B18 fehlt: Ist der Fall „Container endet zwischen Probe und Sitzung“ nicht trennbar, geht der Befund an den Projektinhaber (Fundament §4) |
| B19 | nein (Hinweis) | **ja** | Spezifikations-Plan §1 (letzter Satz der Zeile zur `12`) und DoD 1 |
| A1, A2, A3 | nein (F-4) | **ja** | Spezifikations-Plan §1: Zeilen 648, 936, 1181 |
| A4, A5, A6, A7 | nein (F-4) | **ja** | Befehl-Plan §3: Handbuch (A4, A5, A7) und JSON-Kontrakt (A6) |

**10. Commit-Messages und Hard Rules.** Alle zehn Messages der Range tragen eine Kennung (Slice, `LH-…` oder `ADR-…`); `tools/harness/commit-msg-traceability.sh` (Host-Skript als Text-Werkzeug, der Hook ist im Klon nicht aktiviert, `core.hooksPath` leer) endet für jede Message mit Exit 0, `make doc-commits` auf `770113e8..5f643a81` und auf `origin/main..HEAD` ebenfalls (siehe *Sensors*). Die CR-Commits tragen die Slice-Kennung; alle außer den Verschiebe-Commits `961e22cd` und `a5ede8a3` tragen zusätzlich [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), `d01fc191` alle vier `LH`-Kennungen. Keine Message trägt eine Struktur-Kennung außer `d20e4dd6` (drei: [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen), in einer Aufzählung der Aufträge). Das ist die dokumentierte Ausnahme (§3 F-7, §7 Punkt 9, Entscheidung des Projektinhabers) gegen [`AGENTS.md`](../../AGENTS.md) §5 Regel 1; `make doc-commits` meldet sie nicht (V-7). `961e22cd` ist ein reiner Move (`100 %`, 0 Zeilen); der Inhalt der Eröffnung kam vorher (`5e54f672`), der Verweis-Nachzug danach (`a5ede8a3`) — damit gilt §3.3 Punkt 2 für den Übergang nach `in-progress/` wie im Vorgänger-CR. Kommentare und Marker ([`AGENTS.md`](../../AGENTS.md) §3.7, §3.9, §3.10): Der Diff berührt keinen Code-Kommentar. Abgelaufene Marker siehe V-3, V-4 und V-8.

**11. Bekannte offene Punkte (Bewertung, keine Entscheidung).**

- **(i) §1, Begründung zu Stelle 3 (Zeile 58): „der Aufruf umfasst die Sitzung, in der kein Socket eingebunden ist“.** Bestätigt: Der Satz steht unverändert, `5f643a81` hat im CR-Slice nur das Ziel berichtigt. Er gehört zur Klasse von Nachlauf-1 F-1 (Zusage über den Zustand der Sitzung statt über den Beitrag des Befehls) und steht in keiner Fundstellenliste von Nachlauf-2 F-1. Er steht nicht im Lastenheft und ändert keinen Vertrag; er paraphrasiert [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1. Bewertung: INFO (V-4).
- **(ii) Spezifikations-Plan, Übergabe B9: `shell` ohne `u-boot.yaml` „nicht zugesagt“.** [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) bindet diesen Fall schwächer als B18. In B18 kann der Befehl bei einem Fehler *dieses Aufrufs* (der Client scheitert) mit `0` enden (Nachlauf-1 F-2); der Fall hat im Lastenheft keinen Code, die Bindung greift. Bei B9 geht es um einen Abbruch aus einem *früheren* Start, den ohne `u-boot.yaml` niemand erkennt. Das Lastenheft setzt dafür die Voraussetzung selbst (Zeile 498: die `11` „setzt die Statusdatei voraus“; „ein erneutes `up` meldet einen früheren Abbruch nicht“), und `shell` startet „wie `up`“ mit den Fehlerklassen von `up` (Zeile 494). Die Lücke ist also im Lastenheft für `up` ausdrücklich begrenzt und für `shell` über „wie `up`“ mitgemeint; es fehlt nur das Wort `shell` im zweiten Teilsatz. „Nicht zugesagt“ hebt daher keine Bindung auf. Besser als „nicht zugesagt“ wäre eine benannte Grenze mit Verweis auf die Voraussetzung der Statusdatei, mit dem Unit-Fall, den der Befehl-Plan vorschlägt (V-5).
- **(iii) Sicherheitsbefund `docker-cli`.** Bestätigt außerhalb dieses CR; der CR-Text sagt dazu nichts Falsches zu. Gelesen: `internal/hexagon/application/devcontainer_features.go` führt `"docker-cli"` mit der Quelle `ghcr.io/devcontainers/features/docker-outside-of-docker`, ohne Profilfilter; die Vorlage `devcontainer.json.tmpl` rendert den Block `features` (`{{- if .Features}}`, Zeile 43) außerhalb von `{{- if .Sandbox}}` (Zeile 11). Das Lastenheft sagt in `0.20.0` nur zu, was der Befehl einbindet („bindet der Befehl keinen … Socket des Hosts ein“, „erhält den Socket durch diesen Befehl nicht“). Die Klammer in [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) sagt nur, dass der Hilfs-Container nicht zum erzeugten Ergebnis gehört. Die bestehende Anforderung „keinen Container-Runtime-Socket des Hosts einbinden“ im Sandbox-Profil bleibt vom Katalog-Feature verletzt; sie ist Gegenstand des Bugfixes `v0.13.1`, nicht dieses CR. Ob das Feature den Socket tatsächlich einbindet, habe ich nicht gemessen (nur gelesen). Hinweis zur Reichweite (V-6): Der Auftrag in §7 Punkt 12 beschreibt den Bugfix für das Sandbox-Profil; die Vorlage rendert das Feature in beiden Profilen.

### Konformität mit den Entscheidungen (§7)

| Entscheidung | Befund |
|---|---|
| 1 Stellen 1 bis 5, 6a freigegeben wie vorgeschlagen | zeichengleich (Einzelnachweis 1) |
| 2 Stelle 6b, beide Teile | Liste (Zeile 189) und Satz (Zeile 196) zeichengleich |
| 3 Stelle 7 | Zeile 498 zeichengleich |
| 4 Stelle 8, Historie nur mit freigegebenen Stellen | Version und Historie wie berechnet; die Zeile nennt Stellen 1 bis 7, nichts darüber hinaus |
| 5 Reihenfolge mit dem Messstände-CR | `0.20.0` belegt; Titel, Frage 2, DoD und Zeile in `open/README.md` des Messstände-CR sagen „nach 0.20.0“; `0.19.0` steht in `open/` nur noch als Version des geschlossenen Vorgänger-CR |
| 6 Folge-Slices: `script`, Hilfe ohne Kennung, Frage 4 offen | Fundament DoD 2 und Befehl DoD 1 tragen den Terminal-Test und den Langtext ohne Anforderungskennung; Frage 4 bleibt in der Welle-Datei offen |
| 7 Bestätigung der zwei Lesarten | 6b-Satz als letzter Satz des Absatzes; Historie-Zeile mit Klammer zu Stelle 7, ohne den 6b-Satz zu nennen |
| 8 Lastenheft enger (F-5, F-6) | der erste Teil (Prüfung vor der Shell) steht unverändert im Lastenheft; der zweite Teil (Socket) ist durch Entscheidung 10 abgelöst, der ersetzte Satz kommt im Lastenheft nicht mehr vor |
| 9 `d20e4dd6` nicht umgeschrieben | Commit unverändert in der Range, Vermerk in §3 und §7 |
| 10 Nachlauf-F-1 | beide Sätze zeichengleich im Lastenheft, die ersetzten nicht mehr |
| 11 Nachlauf-F-2 ohne Lastenheft-Änderung | `spec/` seit `a46dcf5b` unverändert; Messauftrag im Fundament, Grenze im Spezifikations-Plan |
| 12 Nachlauf-F-4 außerhalb | der CR ändert weder Katalog noch Vorlage noch Handbuch; Pläne verweisen auf den Bugfix |

### Sensors

Alle Läufe habe ich selbst auf `5f643a81` gefahren (Docker, Host `linux/amd64`). Angegeben ist der Exit-Code des `make`-Prozesses.

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` (sauberer Baum) | pass, Exit 0, 75 s | `lint`: `0 issues.` · `test`: 24 Pakete `ok` · `coverage-gate: OK — coverage 91.70% meets threshold 90%` · `d-check: 622 Datei(en) geprüft, 0 Befund(e)` · `baseline-verify: v6.18.1 OK — 54 Dateien` · `ablage-pruefung: OK`, Wächter-Test `26 Faelle, 0 Fehlschlaege` · Stempel `gates-passed.head` = `5f643a81…` |
| `make docs-check` | pass, Exit 0 | `622 Datei(en) geprüft, 0 Befund(e)`; Lauf mit dieser Datei unter *Abschluss* |
| `make doc-trace` | pass, Exit 0 | `84 Anforderung(en), 0 Waise(n).`; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) je `ok`; Grenze für [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): Einzelnachweis 4 |
| `make doc-commits RANGE="770113e8..5f643a81"` | pass, Exit 0 | `10 Commit(s) — OK`, 0 Befunde |
| `make doc-commits RANGE="origin/main..HEAD"` | pass, Exit 0 | `13 Commit(s) — OK`, 0 Befunde |
| `make doc-immutable RANGE="282f04b6..5f643a81"` | pass, Exit 0 | `13 Commit(s) — OK`, 0 Befunde; dazu `git diff ec8fdf5a..5f643a81 -- docs/plan/adr` leer (beide ADRs führen den Status inline, das Gate sieht sie nicht) |
| `make doc-planning` | pass, Exit 0 | `622 Datei(en) geprüft, 0 Befund(e)` |
| `make docs-check` mit G-1 | fail (erwartet), Exit 2 | `matrix-forbidden … contract-spec → adr`, 1 Befund |
| `make docs-check` mit G-2 | fail (erwartet), Exit 2 | `anchor-missing`, 1 Befund |
| `make docs-check` mit G-3 | fail (erwartet), Exit 2 | `matrix-forbidden … contract-spec → tech-spec`, 1 Befund |
| `make docs-check` und `make gates` mit dieser Datei | siehe *Abschluss* | — |

Traceability:

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Wortlaut zeichengleich (Einzelnachweis 1), Zeilen 490, 494, 495, 497, 498; `doc-trace` `ok`; Nachlauf-Reviews 1 und 2 über die Fix-Runden |
| [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) | Zeile 1140, Vergleich in Einzelnachweis 4 (kein Gate) |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | Zeilen 448, 189, 191, 196; `doc-trace` `ok` |
| [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) | Stelle 1 wörtlich Entscheidung c von [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md); Spiegelstrich zu `shell` nirgends weiter als die Entscheidungen 1 bis 3 von [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md); beide ADRs seit `ec8fdf5a` unverändert |
| Referenzmodell | G-1 und G-3 (Gate fängt die Kanten nach unten) |

Carveouts:

- Neu: none
- Gelöst: none
- Unverändert: alle. `carveouts.md` liegt nicht im Diff. Kein Gate und keine Schwelle ist verändert.

Nicht ausgeführt:

- `make test-docker`: Der Diff berührt keinen Docker-Pfad (`git diff --stat 770113e8..5f643a81 -- internal cmd` leer).
- `make verify-depguard`: keine Importregel, kein Import berührt.
- `make ci`, `make govulncheck`, `make image-scan`: keine Abhängigkeit, kein `Dockerfile`, kein Workflow, kein Release-Pfad berührt (Diff über `Dockerfile`, `.github`, `Makefile`, `repo.mk` leer). Für Markdown-Änderungen ist nach der Tabelle in [`harness/verification.md`](../../harness/verification.md) `make docs-check` der Mindest-Sensor und `make gates` der Closure-Sensor.
- Mutation M0 des Implementers (Stand `a5ede8a3`): nicht nachgefahren, nur die Voraussetzung gelesen (Einzelnachweis 5). G-1 bis G-3 decken M1 bis M3.
- Die Freigaben im Gespräch: nicht nachprüfbar.
- Validierung gegen den Bedarf des Nutzers (`shell --no-interactive`, Arbeiten in der Shell): Aufgabe des Validators.
- Eine Messung von `docker-cli` oder von `docker exec`: nicht Gegenstand dieses CR.

Commit / Artefakt:

- Geprüft auf `5f643a81`; kein Release-Artefakt. Nichts gepusht (`origin/main` = `282f04b6`, 13 Commits voraus).

## Befunde und Diskrepanzen des Verifiers

Keine DoD-Verletzung. V-1 bis V-3 sind Closure-Auflagen für die Plan-Buchhaltung des Slice.

| ID | Klasse | Befund | Wirkung |
|---|---|---|---|
| V-1 | LOW | Die Spalte „im Plan“ der Übergabe-Tabelle in §7 ist auf `5f643a81` in 13 von 26 Zeilen falsch oder unvollständig: B3, B15, B18, B19 und A1 bis A7 stehen auf „nein“, sind aber im Plan getragen (Einzelnachweis 9); B9 und B10 stehen auf „nein (Hinweis)“, sind aber als Hinweis bzw. Risiko im Plan geführt. Der Slice nennt den Stand `fc3151f0` selbst und sagt „hier nicht nachgeprüft“. | Der Implementer setzt die Spalte in der Closure auf den Ist-Stand (Werte in Einzelnachweis 9) und ersetzt den Satz „Stand fc3151f0 … nicht nachgeprüft“. B18 bekommt den Zusatz aus Fundament §4. |
| V-2 | LOW | Die Plan-Tabelle und ihr Abgleich nennen auf `a46dcf5b` nur den fremden Commit `6289f5c6`. `5f643a81` fügt weitere Dateien in den Diff und ändert das Ziel des CR-Slice; die Zeile `docs/reviews/` führt zwei Reports („beide unverändert“), nicht den dritten, und diesen Report nicht. Der Diff hat jetzt 14 Zeilen in 13 Dateien (Einzelnachweis 8). | Closure: Abgleich nach der letzten Änderung wiederholen ([`AGENTS.md`](../../AGENTS.md) §3.10), `5f643a81` als fremden Commit nennen, die Zeile `docs/reviews/` um das zweite Nachlauf-Review und diesen Report ergänzen. |
| V-3 | LOW | §3 hat keine Fix-Runde-Einträge für die Findings F-1 bis F-3 des zweiten Nachlauf-Reviews, und zwei Marker mit Ablaufpunkt sind abgelaufen: „**Offen, Auftrag an den Planner**“ (§3 *Fix-Runde 2*, Punkt F-2; erledigt durch `5f643a81`) und „hier nicht nachgeprüft“ (§7, vor der Übergabe-Tabelle). Ausgänge: F-1 durch `5f643a81` (Welle, Befehl-Slice, CR-Ziel, `open/README.md`), F-2 durch `5f643a81` (Spezifikations-Plan §1, DoD, §6), F-3 durch Fundament §4 und B18. | Closure: Ausgang je Finding im Slice eintragen und die beiden Marker auf den Ist-Stand setzen ([`AGENTS.md`](../../AGENTS.md) §3.10). |
| V-4 | INFO | Bekannter Punkt (i): §1, Begründung zu Stelle 3 (Zeile 58) sagt „in der kein Socket eingebunden ist“ und beschreibt in „Ohne die Änderung sagte der Spiegelstrich …“ die verworfene Fassung. Nachlauf-2 F-1 nennt die Zeile nicht, `5f643a81` fasst sie nicht an. | Kein Vertragstext. Entscheidung beim Planner: auf den Beitrag des Befehls berichtigen („in der der Befehl keinen Socket einbindet“) oder als Begründung der freigegebenen Stelle belassen und in §7 vermerken. Eine Berichtigung braucht kein Nachlauf-Review, denn das Lastenheft bleibt unberührt. |
| V-5 | INFO | Bekannter Punkt (ii): B9. [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) bindet den Fall `shell` ohne `u-boot.yaml` schwächer als B18 (Einzelnachweis 11). Der Spezifikations-Plan fällt ohne Test auf „nicht zugesagt“ zurück. | Hinweis an den Planner: im Spezifikations-Plan die Grenze als Folge der im Lastenheft genannten Voraussetzung der Statusdatei formulieren und den Unit-Fall des Befehl-Plans (§3) in dessen DoD ziehen. Keine Lastenheft-Änderung. |
| V-6 | INFO | Bekannter Punkt (iii): Der CR-Text sagt zu `docker-cli` nichts Falsches. Der Bugfixauftrag in §7 Punkt 12 und B5 nennen das Sandbox-Profil; die Vorlage rendert den Block `features` (`devcontainer.json.tmpl`, `{{- if .Features}}`) außerhalb des Sandbox-Zweigs, und der Katalogeintrag hat keinen Profilfilter. Das zweite Nachlauf-Review merkt dasselbe für die Handbuch-Sätze zum Standardprofil an. | Hinweis für den Slice des Bugfixes `v0.13.1`: Reichweite beide Profile prüfen. Keine Änderung am CR. |
| V-7 | INFO | F-7 bestätigt: Der Body von `d20e4dd6` nennt [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) und [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen). `make doc-commits` meldet es nicht (Exit 0 über die Range); der Sensor prüft die Anwesenheit einer Kennung, nicht die Abwesenheit einer Struktur-Kennung. Der Vorgänger-CR hatte denselben Befund als Review F-7 („Verfeinerungs-Kennung im Body“). | Closure: Sensor-Lücke in §7 vermerken; das Register kennt dafür keine Beobachtung (Suche in `observations/`), der Planner entscheidet über einen neuen Eintrag. |
| V-8 | INFO | Zwei Ablaufmarker und eine Auflage der DoD: (a) Die DoD verlangt Fundstellen „in §7“; sie stehen in §3 *Umsetzung*. (b) Die Überschrift „Umsetzung (Implementer, 2026-10-21; nichts gepusht, kein Tag)“ trägt einen Marker, der mit dem Push abläuft; heute ist er wahr. (c) Die Commits `961e22cd`, `a5ede8a3` und `9e6c47d4` gehören in denselben Push: Nach dem Slice (M0) ist `a5ede8a3` allein rot (`planning-drift`; die Voraussetzung habe ich gelesen, den Lauf nicht wiederholt), und ein Verschiebe-Commit ohne seinen Verweis-Nachzug ist nach der Regel der Baseline kurz rot. | Closure: (a) §7 verweist auf §3; (b) Marker bei der Closure auf den Ist-Stand setzen oder entfernen; (c) alle Commits der Closure in denselben Push. |
| V-9 | INFO | Eigene Vorab-Suche in `internal/`: Der Doc-Kommentar von `startFallbackWording` (`internal/hexagon/application/devcontainer_fallback.go`, Zeile 270) lautet „the text of `devcontainer up`“. Wenn `shell` die Funktion für die Warnungen nach dem Start nutzt, nennt der Kommentar den neuen Befehl nicht. Kein Plan führt die Datei für diesen Grund (der Plan zum Urteil nach dem Zustand führt sie als „prüfen“ aus anderem Anlass). | Hinweis an den Befehl-Slice zu seiner Suche nach Aufzählungen im Code (dort §6, Risiko „Eine Aufzählung der Befehle bleibt ohne `shell`“); kein Befund am CR. |

Ich habe keine Behauptung ohne Beleg im Lastenheft-Text gefunden. Die Aussagen des Implementers zu Sensoren, Fundstellen, Zahlen der Plan-Tabelle und Mutationen M1 bis M3 habe ich nachgefahren und bestätigt.

## Prüfliste für die Closure (Planner / Implementer)

1. **Diesen Report committen.** Die Nachricht trägt die Slice-Kennung und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), keine `SPEC`- oder `ARC`-Kennung.
2. **§3 nachziehen (V-2, V-3):** Abschnitt zu den Ausgängen des zweiten Nachlauf-Reviews (F-1 bis F-3, Verweis auf `5f643a81`), Marker „Offen, Auftrag an den Planner“ auf den Ist-Stand, die Zeile `docs/reviews/` der Plan-Tabelle (drei Reports plus dieser), den Abgleich gegen `git diff --stat 770113e8..HEAD` mit den fremden Commits `6289f5c6` und `5f643a81` wiederholen.
3. **§7 Übergabe-Tabelle (V-1):** Spalte „im Plan“ nach Einzelnachweis 9; den Satz „Stand `fc3151f0` … hier nicht nachgeprüft“ ersetzen; B16: Verweis auf diesen Report.
4. **§1 Zeile 58 (V-4):** berichtigen oder in §7 als Chronik der freigegebenen Stelle vermerken.
5. **§7 Closure-Notiz füllen:** was funktioniert hat und was anders ging (zwei Fix-Runden, zwei Nachlauf-Reviews, der Fund `docker-cli`); Freigabe-Vermerk je Stelle (§7 Punkte 1 bis 12 stehen, ein Verweis genügt); Verweis auf §3 *Umsetzung* und die Fix-Runden für Fundstellen, Vorab-Suche und Sensoren (V-8); Gate-Läufe nach den Closure-Commits mit Exit-Code; Wortlaut-Vergleich: Verweis auf diesen Report; Hinweise an die Folge-Slices (F-8, F-10, F-11, V-5, V-9) und an den Bugfix `v0.13.1` (V-6); die Sensor-Lücke aus V-7 und die Grenzen von `doc-trace` ([`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) ist nicht in der Belegsicht) und der fehlenden Wortlaut-Prüfung (kein Gate vergleicht den Wortlaut mit der Freigabe).
6. **Steering-Loop-Lerneintrag** (Vorschläge, der Planner entscheidet; `liegt in` nur, wenn etwas verkörpert wird): (a) Ein Lastenheft-Satz, der eine Eigenschaft der Laufzeit zusagt, die auch die Konfiguration bestimmt, sagt den Beitrag des Befehls zu, nicht den Zustand (Nachlauf-1 F-1 und der Fix-Runde-2-Satz); (b) ändert eine Fix-Runde einen Vertragssatz, sucht der Implementer dessen Aussage vor dem Nachlauf-Review in allen Plänen derselben Welle (Nachlauf-2 F-1 fand vier Plan-Stellen, ich eine weitere in V-4); (c) eine Spalte „im Plan“ ist ein Zustandsfeld mit Ablaufpunkt und trägt ihren Stand (V-1); (d) die Sensor-Lücken: `doc-commits` ohne Prüfung auf Struktur-Kennungen (V-7), `doc-trace` ohne `LH-SA-*`, kein Gate für den Wortlaut gegen die Freigabe.
7. **Beobachtungs-Register** (Eingaben, kein Zähler wird gesetzt): [`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (zwei Nachläufe in einem Vorgang); [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) und [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Review F-1, F-2, Nachlauf-1 F-1); [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Review F-3, F-9, Nachlauf-2 F-1, V-4); [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../plan/planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (Stellen 6a und 6b, F-10, V-9); [`BEO-DOK/dod-sensor-sieht-anforderung-nicht`](../plan/planning/observations/BEO-DOK/dod-sensor-sieht-anforderung-nicht/observation.md) ([`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)); [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../plan/planning/observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (Stelle 7, B9); neu zu entscheiden: Struktur-Kennung in Commit-Message ohne Sensor (V-7, zweiter Beleg nach dem Vorgänger-CR).
8. **Risiken §6, Ausgang je Zeile** (Vorschläge aus meiner Prüfung):

   | Risiko | Vorschlag | Grund |
   |---|---|---|
   | Vorab-Suche übersieht eine Aufzählung | im Lastenheft nicht eingetreten; außerhalb eingetreten und übergeben | Reviewer und Verifier fanden im Lastenheft keine weitere Stelle; die Funde F-10 und V-9 gingen als Aufträge an den Befehl-Slice |
   | Wortlaut ändert sich in der Fix-Runde | eingetreten, aufgefangen | zwei Fix-Runden, zwei Nachlauf-Reviews; das erste fand die Lücke, die die erste Runde offen ließ |
   | Zusage reicht weiter als die Tests | weiter offen | Eigentümer benannt (Fundament, Befehl); die Tests liefern die Code-Slices |
   | Bedingung des Wortlauts ohne Test | weiter offen | Übergabe-Tabelle B1 bis B19; Beleg-Tabelle im Spezifikations-Slice |
   | Stelle 6a wird nicht freigegeben | entfallen | bereits eingetragen |
   | Versionsfolge mit dem Messstände-CR | entfallen | `0.20.0` belegt, die Verweise des Messstände-CR stehen auf „nach 0.20.0“ |
   | Aussage über den Code im Plan stimmt nicht | eingetreten | F-9 (Zitat), F-3 (Marker in Folge-Plänen), V-4 |

9. **DoD-Haken** setzen: alle zehn Punkte, wenn die Punkte 2 bis 8 und 10 getan sind.
10. **Roadmap und `open/README.md`:** „Nichts in Arbeit.“ zurücksetzen, sobald der Slice `in-progress/` verlässt (sonst `planning-drift`); die Zeile des Slice in `open/README.md` auf „geliefert (liegt in `done/`)“. Die Zeile des Messstände-CR („nach 0.20.0“) stimmt und bleibt.
11. **Welle-Datei:** keine inhaltliche Änderung nötig (Tabellenzeile, Vermerke zu Frage 1 und 2 und die Reihenfolge in §5 stimmen mit dem Slice überein); das Make-Target zum Verschieben von Slices zieht die Pfade nach (Welle, Folge-Pläne, `open/README.md`, Roadmap-Zeiger, Reports).
12. **Reihenfolge** ([`AGENTS.md`](../../AGENTS.md) §3.3): erst der Inhalt (Punkte 2 bis 9: Plan-Abgleich, Übergabe-Tabelle, Closure-Notiz, Register, Risiken, Haken), dann der reine Move nach `done/` per Make-Target, dann der Verweis-Nachzug, dann der Roadmap-Reset; alle Commits in denselben Push. Das Lastenheft bleibt unverändert; ändert die Closure doch Lastenheft-Text, ist ein drittes Nachlauf-Review und ein neuer Verifier-Lauf fällig.
13. **Danach:** `make docs-check` und `make gates` auf dem letzten Commit, Exit-Codes in §7. Die Start-Bedingung des Slices zum Urteil nach dem Zustand (CR in `done/`) ist dann erfüllt. Push und Tag nur nach ausdrücklicher Freigabe des Projektinhabers.

## Minimal Handoff

```text
Role: Verifier
Input context: Slice §1–§8 (Umsetzung, Fix-Runde 1 und 2), drei Review-Reports, ADR-0030 und ADR-0031, Lastenheft 0.20.0, git diff 770113e8..5f643a81, Pläne der Folge-Slices und Welle-Datei auf 5f643a81
Changed artefacts: dieser Verification-Report (neu, nicht committet)
Evidence: make gates, docs-check, doc-trace, doc-commits (zwei Ranges), doc-immutable, doc-planning grün; Wortlaut zeichengleich (26 Prüfungen, 0 Fehlschläge); Gegenproben G-1 bis G-3 rot, zurückgesetzt
Open risks: V-1 bis V-3 (Plan-Buchhaltung des Slice, Closure-Auflagen); V-4 bis V-9 INFO; kein Gate vergleicht den Wortlaut mit der Freigabe
Next role: Hauptagent (Commit dieser Datei), dann Planner/Implementer (Closure nach Prüfliste)
```

## Abschluss

Nach dem Anlegen dieser Datei liefen beide Sensoren auf dem Arbeitsbaum (`HEAD` `5f643a81` plus diese Datei, ungestagt):

- `make docs-check`: Exit 0, `d-check: 623 Datei(en) geprüft, 0 Befund(e)`. Ein erster Lauf fand sieben unverlinkte Kennungen (`id-unlinked`: je drei Struktur-Kennungen in Einzelnachweis 10 und in V-7, eine in der Tabelle der Gegenproben); sie sind verlinkt.
- `make gates`: Exit 0, 62 s; `lint` `0 issues.`, `coverage-gate` 91,70 %, `d-check` 623 Dateien und 0 Befunde, `baseline-verify` v6.18.1 OK, `ablage-pruefung` OK (26 Fälle). Die Datei ist ungestagt und färbt `ablage-pruefung` nicht.

`git status` zeigt nur diese Datei als ungetrackt; `spec/lastenheft.md` ist nach den Gegenproben unverändert.
