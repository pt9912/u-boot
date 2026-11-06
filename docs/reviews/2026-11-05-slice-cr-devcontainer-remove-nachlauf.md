# Review-Report (Nachlauf): slice-cr-devcontainer-remove, Folgeänderung Lastenheft 0.37.1 — 2026-11-05

**Review-Art:** Nachlauf-Review zu einer Änderung des Wortlauts nach dem Review (Beobachtung `BEO-DOK/vertragsaenderung-nach-review`, Risiko 6 des Slice, F-9 des Vorlauf-Reports [`2026-11-05-slice-cr-devcontainer-remove`](2026-11-05-slice-cr-devcontainer-remove.md)).
**Gegenstand:** `git diff 7f5110cc..a7914c99`, 7 Dateien, 146 Einfügungen, 50 Löschungen. Betroffen sind `spec/lastenheft.md`, der CR-Vorschlag [`cr-vorschlag-devcontainer-remove.md`](../plan/cr/cr-vorschlag-devcontainer-remove.md), der Slice [`slice-cr-devcontainer-remove`](../plan/planning/in-progress/slice-cr-devcontainer-remove.md), `open/README.md`, die zwei Folge-Slices und die Welle [`welle-devcontainer-remove`](../plan/planning/welle-devcontainer-remove.md).
**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5. Der Report ist vom Aufrufer unverändert aus der Ausgabe des Reviewers übernommen, weil die Rolle keine Report-Dateien schreibt.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice §1 bis §8 (Diff in Wortform) und CR-Vorschlag (Stellen (a) bis (k)).
- Vorlauf-Report [`2026-11-05-slice-cr-devcontainer-remove`](2026-11-05-slice-cr-devcontainer-remove.md).
- [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) im Lastenheft.
- Code: `internal/hexagon/application/devcontainer_start.go` (`start`, `validInstance`, `Down`), `downservice.go` (Zeile 202, „--volumes declined by user“) und `internal/adapter/driving/cli/terminal.go`.
- Stichproben der neu in die Plan-Tabellen der Folge-Slices eingetragenen Fundstellen. Sie stimmen:
  - `README.md` Zeilen 214, 215, 222 und 225 sowie `README.de.md` Zeilen 218, 219, 226 und 229;
  - `docs/maintainer/zusagen-und-messstaende.md` Zeilen 58 bis 60;
  - `spec/spezifikation.md` Zeilen 619, 631, 866, 870, 875, 876, 938, 1354 und 1602;
  - `spec/architecture.md` Zeilen 299, 303, 314, 381, 651 und 656;
  - `docs/maintainer/cli-json-contract.md` Zeilen 1244 und 1356 bis 1358;
  - `docs/user/benutzerhandbuch.md` Zeilen 59, 60, 65, 164, 324, 845, 1042, 1043, 1057 und 2259.

**Sensors:**

- `make docs-check`: gefahren, `1167 Datei(en) geprüft, 0 Befund(e)`. Der Stand ist HEAD `a7914c99`; der Arbeitsbaum ist sauber.
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält nur Doku, keinen Code, keinen Docker-Pfad und keine Importregel.
- Eigene Proben, nur lesend:
  - Wortvergleich der Stellen (h), (i), (j) und der Historie-Zeile (k) des CR-Vorschlags gegen `git show HEAD:spec/lastenheft.md`. Das Ergebnis ist 4 von 4 wortgleich, mit Link-Form-Normierung wie bei der Übernahme.
  - `git diff -U0 7f5110cc HEAD -- spec/lastenheft.md` hat 5 Einfügungen und 4 Löschungen. Das sind Versionsfeld, drei Spiegelstriche und die Historie-Zeile, sonst nichts.
  - Suche der hinzugefügten Lastenheft-Zeilen nach `adr-`, `slice`, `planning`, `roadmap`, `docs/`, `SPEC-`, `.go` und `label`: 0 Treffer.

## Findings

### MEDIUM `docs/plan/cr/cr-vorschlag-devcontainer-remove.md` · Zustandsmarker nennen einen Stand, den der Commit beendet hat

**Klasse:** Zustandsmarker mit Ablaufpunkt nicht entfernt (dritter Lauf in Folge; Steering-Loop-Schwelle 3 erreicht)
**Quelle:** `AGENTS.md` §3.7 und §3.10 (Marker mit Ablaufpunkt, Abgleich nach der letzten Fix-Runde); Skill MEDIUM „ein LOW-Muster wiederholt sich“
**Befund:** Der Commit `a7914c99` trägt `spec/lastenheft.md` in Version `0.37.1`, der Status ist `Accepted`. Mehrere Planungsstellen sagen weiter „im Arbeitsbaum“ und „nicht committet“. Zwei davon sind als Anweisung oder Abgleich falsch:

- Der CR-Vorschlag sagt „Alt ist der Wortlaut von `0.37.0` (`git show HEAD:spec/lastenheft.md`), Neu der Arbeitsbaum.“ `git show HEAD:` liefert jetzt den Neu-Wortlaut.
- Der Slice sagt in §3 „Stand `HEAD` `7f5110cc` und Arbeitsbaum, nicht committet“. `git diff --stat HEAD` ist jetzt leer.

Weitere Fundstellen mit demselben Muster:

- CR-Vorschlag: „Neu der Arbeitsbaum“ (Stellen (h) bis (k)).
- Slice [`docs/plan/planning/in-progress/slice-cr-devcontainer-remove.md`](../plan/planning/in-progress/slice-cr-devcontainer-remove.md): „Folgeänderung `0.37.1` im Arbeitsbaum, nicht committet“ (§3), „ist dem Projektinhaber noch vorzulegen und nicht committet“ (Kopf).
- `docs/plan/planning/open/README.md`: „Folgeänderung `0.37.1` nach dem Review im Arbeitsbaum“.
- `docs/plan/planning/welle-devcontainer-remove.md`: „(Lastenheft `0.37.0`, CR-Vorschlag“ (Tabelle) und „im Lastenheft `0.37.0` übernommen, die Closure steht aus“ (Reihenfolge). Beide nennen die Version `0.37.0`, das Lastenheft trägt `0.37.1`.

Wahr bleibt „Freigabe des Wortlauts steht aus“. Der Slice, der CR-Vorschlag und die Welle halten es so fest, und die DoD-Zeile 3 ist offen.
**Risiko:** Der nächste Änderer folgt der falschen Aussage, weil kein Gate sie widerlegt. Wer die `git show`-Anweisung ausführt, hält den Neu-Wortlaut für den Alt-Wortlaut. Wer „Arbeitsbaum“ liest, sucht Änderungen, die schon im Commit stehen. Zum Zustand des Vertrags selbst: Rang 1 trägt `Accepted`-Wortlaut, dessen Freigabe die Planungsstellen als ausstehend führen, und die DoD-Zeile 3 verlangt Freigabe vor Übernahme. Dieselbe Klasse stand in den zwei vorherigen Läufen ([`2026-11-05-slice-cr-devcontainer-brew-pakete`](2026-11-05-slice-cr-devcontainer-brew-pakete.md) F-5, [`2026-11-05-slice-cr-devcontainer-remove`](2026-11-05-slice-cr-devcontainer-remove.md) F-6). Ob dieser Zustand zulässig ist oder die Freigabe nachzuholen ist, entscheidet der Projektinhaber über den Architect, nicht der Reviewer.
**Verifizierbar:** Review-only (`git status`, `git show HEAD:spec/lastenheft.md`)

### MEDIUM `spec/lastenheft.md` · „Ohne Freigabe (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage) endet der Befehl mit Exit-Code `10` und es geschieht nichts“ — Reihenfolge von Freigabe und Containersuche bleibt offen

**Klasse:** Anforderung lässt einen Nutzerpfad offen (`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`)
**Quelle:** [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code-Zuordnung bindend); [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung); Skill MEDIUM „Fehlerbehandlung … keiner Exit-Code-Kategorie eindeutig zugeordnet“
**Befund:** Die Anforderung stellt die Abfrage „bevor er etwas stoppt oder entfernt“. Der Befehl muss die Volumes vor der Abfrage nennen, also zuerst den Container suchen. Zwei andere Sätze gelten für den Fall ohne Container:

- „gibt es keinen Container dieser Instanz, entfernt er nichts und endet mit Exit-Code `0`, auch wenn es andere Container des Projekts gibt“
- „Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`“

Für `remove --volumes` im nicht-interaktiven Modus ohne `--yes` liefern „Exit-Code `10`“ und „Exit-Code `0`“ verschiedene Antworten, sobald es keinen Container gibt. Die Fehlerklassen-Klausel „wenn `remove --volumes` die Freigabe nicht erhält“ löst das nicht auf. Der CR-Vorschlag prüft den Fall weder in der Satz-Lesart noch in der Messgrenzen-Tabelle.
**Risiko:** Ein Aufräum-Skript `remove --volumes --instance a --no-interactive` läuft ein zweites Mal gegen eine schon entfernte Instanz.

- Prüft der Implementer zuerst die Freigabe, wie es das Vorbild `down --volumes` tut (`downservice.go` bricht vor dem Engine-Aufruf ab), endet der zweite Lauf mit `10`.
- Sucht er zuerst den Container, endet er mit `0`.

Beide Lesarten sind mit dem Wortlaut vereinbar. Das Skript bricht oder läuft durch, je nach Implementer, und der pinnende Test des Umsetzungs-Slice (`AGENTS.md` §3.10) hat kein eindeutiges Soll.
**Verifizierbar:** Review-only (Test im Umsetzungs-Slice)

### LOW `spec/lastenheft.md` · „läuft er nicht“ im Satz „Im ausgelieferten Container-Image (…) läuft er nicht“ — Reichweite des Satzes für `remove` ungeklärt

**Klasse:** Zusage-Menge des neuen Befehls undefiniert (zweite Hälfte von F-4 des Vorlaufs)
**Quelle:** `AGENTS.md` §3.10 (Bedingung im Lastenheft bekommt einen Test oder wird gestrichen)
**Befund:** Der Halbsatz „und `u-boot devcontainer remove` entfernt,“ nimmt `remove` in die Zusage-Menge „Linux-Host (`amd64`), lokale Engine“ auf. Der letzte Satz des Spiegelstrichs steht wörtlich da und bezieht sich grammatisch auf „Der Befehl“, den Start. Ob `remove` im ausgelieferten Image ebenfalls nicht läuft, sagt weder der Spiegelstrich noch der CR-Vorschlag. Dort steht nur „bleibt wörtlich“. Die Plan-Tabellen der Folge-Slices tragen dafür keinen Test.
**Risiko:** Bei der Spezifikation und beim Test des Umsetzungs-Slice hat der Implementer keinen Anhalt, ob die Aussage `remove` mitmeint. Er liest sie in die eine oder andere Richtung oder lässt sie ungepinnt.
**Verifizierbar:** Review-only

### INFO `spec/lastenheft.md` · „`--no-interactive` bricht bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab“ — Spannung zu Exit-Code `10` besteht weiter

**Klasse:** Vorbestehende Spannung, keine Änderung durch den Diff (F-7 des Vorlaufs)
**Quelle:** [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
**Befund:** Die 0.37.1-Fassung bindet die abgelehnte Abfrage und „nicht-interaktiver Modus ohne `--yes`“ an Exit-Code `10`. Die Zeile in [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) mit Exit-Code `2` für `--no-interactive` bleibt daneben stehen. Der Slice führt die Frage in §6 an den Architect und den Projektinhaber. Sie ist damit benannt und nicht gelöst.
**Risiko:** Keines aus dem Diff; die Frage betrifft ebenso `down --volumes` und `remove --purge`.
**Verifizierbar:** Review-only

### INFO `docs/plan/planning/open/slice-v1-devcontainer-remove.md` · „(e) abgelehnte Abfrage von `remove --volumes`“ — Machbarkeit des Falls (e) auf Ebene des gebauten Binary nicht belegt

**Klasse:** Plan-Aussage über ein Testmittel ohne gelesene Stelle
**Quelle:** `AGENTS.md` §3.10; Zuständigkeit Implementer und Verifier des Umsetzungs-Slice
**Befund:** Die DoD verlangt einen Docker-Test des gebauten Binary, der eine abgelehnte Abfrage belegt. Die Interaktivität bestimmt `internal/adapter/driving/cli/terminal.go` über `ModeCharDevice`. In `internal/e2e` fand die Suche nach `creack/pty`, `openpty`, `ptmx` und `script -qc` keinen Treffer. Ob der Fall dort ohne neues Testmittel ausführbar ist, ist nicht gelesen. Das ist keine Aussage über den Fall, sondern ein offener Punkt für die Umsetzung.
**Risiko:** keines im Diff
**Verifizierbar:** Review-only

## Auflösung der Findings des Vorlaufs

| Vorlauf | Stand im Diff |
|---|---|
| F-1 (Fehlerklassen) | Aufgelöst. Abgelehnte Abfrage führt zu `10`, Fehler beim Stoppen oder Entfernen zu `12`, und die Klauseln tragen ein Befehls-Subjekt. Die Subjekte stimmen mit dem Code überein. `Down` (Zeilen 781 bis 812) ruft weder `readConfig`, `containerName` noch den Startscript-Pfad auf. Diese laufen nur im gemeinsamen `start` von `up` und `shell`. `validInstance` gilt für beide. Siehe zweites MEDIUM zur offenen Reihenfolge. |
| F-2 (Träger der Kopien) | Aufgelöst. Die Plan-Tabellen der Folge-Slices tragen sie, die Fundstellen stimmen (Stichprobe oben). |
| F-3 (Nutzerpfade a und b) | Satz für Satz aufgelöst (Satz 4, 5 und 7). Das Zusammenspiel mit Satz 9 ist das zweite MEDIUM dieses Laufs. |
| F-4 (Zusage-Menge) | Zur Hälfte aufgelöst: `remove` ist aufgenommen, die Reichweite des Image-Satzes bleibt offen (LOW oben). |
| F-5 (Abgrenzung) | Aufgelöst. Der Satz stimmt mit [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) überein (Compose-Eintrag, `--purge`) und trägt die Gemeinsamkeit. |
| F-6 (Zustandsmarker) | Der alte Marker ist entfernt. Mit 0.37.1 sind neue Marker entstanden, die der Commit überholt (erstes MEDIUM). |
| F-7 (Spannung Exit-Code 2 und 10) | offen, im Slice §6 an den Architect adressiert |
| F-8, F-9 | Testauftrag steht in [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) §2. F-9 ist dieser Report. |

## Geprüft ohne Befund

- **Wortgleichheit CR-Vorschlag und Lastenheft:** Stellen (h), (i), (j) und die Historie-Zeile (k) sind 4 von 4 wortgleich. Die Stellen (a) bis (g) blieben, soweit (h) bis (k) sie nicht ändern, wie übernommen: Der Lastenheft-Diff trägt nichts außerhalb von Version, drei Spiegelstrichen und der Historie-Zeile.
- **Referenzrichtung (`AGENTS.md` §3.4, Skill HIGH #2):** Die hinzugefügten Lastenheft-Zeilen tragen keinen Verweis nach unten, auch nicht in der Historie-Zeile.
- **Anforderungsebene des Wortlauts:** keine Labelnamen, Pfade, Docker-Befehle, Messgrenzen oder Flag-Interna. `--purge`, `--volumes`, `--yes` und `--instance` sind Bedienoberfläche. Die Zusage-Menge ist Entscheidung F4 des Projektinhabers.
- **Abgrenzungssatz** gegen [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) Zeilen 653 bis 670 gelesen: stimmt.
- **Slice §3 Plan-Tabelle gegen `git diff --stat 7f5110cc..HEAD`:** dieselben 7 Dateien, jede mit Zeile. Die Zeilenangaben des Slice zum Lastenheft stimmen mit HEAD (Version 9, Spiegelstrich 517, Fehlerklassen 522, Zusage-Menge 523, Historie 1663). Ausnahme sind die Zustandsmarker (erstes MEDIUM).
- **Lifecycle (`AGENTS.md` §3.3):** keine Renames, der Slice bleibt in `in-progress/`, die Folge-Slices bleiben in `open/`.
- **`roadmap.md` Zeile 45:** nennt die Welle ohne Zustandsmarker, keine Chronik.
- **Kommentar-Selbstcheck (`AGENTS.md` §3.9):** kein Code- oder Skript-Kommentar im Diff. Die Alt/Neu-Erzählung im CR-Vorschlag und in den Entscheidungs-Ausgängen des Slice ist dort ihr Zweck.
- **Accepted-ADR:** keine ADR im Diff.
- **Nutzer-Doku und CHANGELOG:** Der Diff berührt sie nicht (`AGENTS.md` §5 Regeln 5 und 6).

## Nicht geprüft

- **Code, Tests, Dual-Classifier, Docker-Pfad, Importregeln:** im Diff nicht berührt. Die Sentinel- und Diagnostic-Abdeckung von `remove` ist Aufgabe des Review von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md).
- **Ob der Projektinhaber den Wortlaut 0.37.1 gesehen hat:** liegt im Sitzungsdialog und ist nicht belegbar. Sicher ist nur, dass die Planungsstellen die Freigabe als ausstehend führen (siehe erstes MEDIUM).
- **`make gates`:** nicht gefahren (Review-Rolle, Doku-Diff).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zustandsmarker mit Ablaufpunkt nicht entfernt (dritter Lauf in Folge) · Anforderung lässt einen Nutzerpfad offen (Freigabe gegen Containersuche) · Zusage-Menge des neuen Befehls undefiniert (Image-Satz) · Plan-Aussage über ein Testmittel ohne gelesene Stelle.

## Verdikt

**Merge-blockierend:** ja, für die zwei MEDIUM-Findings. Sie sind vor der Closure zu beheben oder durch den Architect oder Projektinhaber mit Begründung zu entscheiden. Ein Widerspruch des Implementers stuft sie nicht herab.

**Übergabe:**

- **Implementer:** Die Marker (erstes MEDIUM) gehen an den Implementer. Er gleicht dabei die Planungsstellen, die Version in der Welle und die Abgleich-Zeile des Slice gegen den Ist-Stand ab.
- **Architect und Projektinhaber:** Das erste MEDIUM betrifft zusätzlich die Freigabe des bereits committeten Wortlauts. Das zweite MEDIUM und die LOW-Zeile betreffen die Reihenfolge der Abfrage und die Reichweite des Image-Satzes im Wortlaut. Jede Wortlaut-Änderung läuft über die Freigabe des Projektinhabers und einen weiteren Nachlauf-Review, nicht über eine stille Korrektur.
- **Steering-Loop:** Die Klasse „Zustandsmarker mit Ablaufpunkt nicht entfernt“ hat dreimal in Folge getroffen. Das ist die Schwelle: Klassifikation in Skill und `harness/review.md` schärfen und prüfen, ob ein Sensor oder eine Hard Rule möglich ist (Skill §Pflege).
- **Verifier:** DoD-Zeile 3 des Slice ist offen und gehört in seinen Abgleich. Dieser Report belegt nur den Nachlauf, nicht die Closure.
