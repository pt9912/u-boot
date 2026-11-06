# Review-Report (Nachlauf 3): slice-cr-devcontainer-remove, Folgeänderung Lastenheft 0.37.3 — 2026-11-06

**Review-Art:** Nachlauf-Review zu einer Änderung des Wortlauts nach dem zweiten Nachlauf-Review (Beobachtung `BEO-DOK/vertragsaenderung-nach-review`, DoD-Zeile 7 des Slice). Vorlauf: [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md).
**Gegenstand:** `git diff 7f568768..18f07fcc` (HEAD), 8 Dateien, 213 Einfügungen, 32 Löschungen. Betroffen sind `spec/lastenheft.md` (drei Zeilen), der CR-Vorschlag [`cr-vorschlag-devcontainer-remove.md`](../plan/cr/cr-vorschlag-devcontainer-remove.md), der Slice [`slice-cr-devcontainer-remove`](../plan/planning/in-progress/slice-cr-devcontainer-remove.md), [`open/README.md`](../plan/planning/open/README.md), die zwei Folge-Slices [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) und [`slice-v2-spec-devcontainer-remove`](../plan/planning/open/slice-v2-spec-devcontainer-remove.md), die Welle [`welle-devcontainer-remove`](../plan/planning/welle-devcontainer-remove.md) und der Report des Vorlaufs.
**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice (Kopf, §1 mit Ausgang H-1, DoD, §3 Tabelle und Abgleich, §4, §6, §8), CR-Vorschlag (Abschnitt „Folgeänderung 0.37.3“, Stellen (p) und (q), Prüftabellen), Vorlauf-Report.
- Im Lastenheft: [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Zeilen 200 bis 216), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 513 bis 524), [`LH-FA-UP-004`](../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) (Zeile 710), [`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (Zeile 668), Historie (Zeilen 1662 bis 1665).
- Als Beleg der Aussagen im CR-Vorschlag: `internal/hexagon/application/downservice.go` (`Down`, Aufruf `runConfirmationGate`), `spec/spezifikation.md` ([`SPEC-038`](../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes), Tabelle und Absatz zu `--json`; Abschnitte zu der Anforderung „Verhalten bei aktivierter Nicht-Interaktivität“ und „Auswertungslogik der Bestätigungsmodi“).

**Sensors:**

- `make docs-check`: gefahren auf HEAD `18f07fcc` (vor Anlage dieses Reports), `1169 Datei(en) geprüft, 0 Befund(e)`. Nach Anlage dieses Reports erneut gefahren: `1170 Datei(en) geprüft, 0 Befund(e)`.
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält nur Doku, keinen Code, keinen Docker-Pfad und keine Importregel.
- Eigene Proben, nur lesend:
  - Wortvergleich der Stellen (p) und (q) des CR-Vorschlags gegen `spec/lastenheft.md` (Link-Form und Fettdruck normiert, per Skript): 2 von 2 wortgleich.
  - `git diff -U0 7f568768 HEAD -- spec/lastenheft.md`: Versionsfeld, Zeile 214 (ein Halbsatz an Stelle des Punktes), Historie-Zeile 1665; sonst nichts.
  - Suche nach weiteren Aussagen zu Exit-Code `10` bei fehlender Freigabe im Lastenheft (`ohne `--yes``, `nicht-interaktiv`, `--no-interactive`): Treffer in den Zeilen 211, 212, 214, 327, 345, 436, 495, 517, 522, 651, 668, 1133 und Historie; keine trägt eine unbedingte Aussage zu `u-boot devcontainer remove --volumes` außer Zeile 517 (Satz 4 und 5, eingeschränkt durch Satz 10).
  - Zeilenangaben des Slice-Plans (9, 214, 1665; gelesen, nicht geändert: 210, 212, 517, 522, 668) gegen HEAD: stimmen.
  - `git ls-files` zum Report des Vorlaufs: der Report ist in HEAD versioniert (siehe LOW 1).

## Findings

### LOW `docs/plan/planning/in-progress/slice-cr-devcontainer-remove.md` · „ist ein neuer, nicht versionierter Bestandteil des Review-Laufs“ — Zustandsmarker mit Ablaufpunkt, im Commit abgelaufen

**Klasse:** Zustandsmarker mit Ablaufpunkt nicht entfernt (Plan-Aussage ohne Entsprechung im Stand)
**Quelle:** `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt wird beim Ablauf auf den Ist-Stand gesetzt; Abgleich der Tabelle gegen `git diff --stat` für den Stand nach der letzten Fix-Runde); Skill LOW mit Konventions-Anker
**Befund:** Der Abgleichssatz in §3 nennt den Report [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md) „ein neuer, nicht versionierter Bestandteil des Review-Laufs“ und führt sieben Dateien. Der Commit `18f07fcc` enthält den Report; `git diff --stat 7f568768..HEAD` zeigt acht Dateien, `git ls-files` führt den Report. Der Satz war vor dem Commit wahr und ist seither falsch.
**Risiko:** Der nächste Änderer wiederholt den Abgleich nach `AGENTS.md` §3.10 gegen `git diff --stat`, findet acht Dateien gegen sieben im Text und kann den Plan nicht auf einen Blick gegen den Diff halten; die Aussage „nicht versioniert“ hält niemand fest, weil kein Gate sie widerlegt.
**Verifizierbar:** Review-only (`git ls-files`, `git diff --stat`)

### LOW `docs/plan/planning/open/slice-v1-devcontainer-remove.md` · „**Bedingungen des Lastenhefts `0.37.2` (Nachlauf-Review des CR):**“ — Überschrift der DoD-Zeile trägt die Version der Bedingungen 1 und 2, Bedingung 3 gehört zu `0.37.3`

**Klasse:** Zustandsmarker mit Ablaufpunkt nicht entfernt (Versionsangabe nicht mit dem Inhalt der Zeile fortgeschrieben)
**Quelle:** `AGENTS.md` §3.10 (Aussagen im Plan nennen nur, was gelesen ist; Marker auf den Ist-Stand setzen); Skill LOW mit Konventions-Anker
**Befund:** Die DoD-Zeile führt unter der Überschrift „`0.37.2`“ jetzt drei Punkte; der neue Punkt (3) nennt selbst „Lastenheft `0.37.3`, Zeile 214“. Die Plan-Tabelle des CR-Slice (Zeile zu [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md)) nennt „Punkt 3 der Bedingungen des Lastenhefts“ ohne Version und trägt die Zeile damit korrekt; nur die Überschrift der Zeile ist nicht fortgeschrieben.
**Risiko:** Der Verifier des Umsetzungs-Slice sucht die Bedingungen der Version `0.37.3` und liest „`0.37.2`“; er kann Punkt 3 als überholt oder zu `0.37.2` gehörig lesen und den pinnenden Test des Falls (d) gegen Zeile 214 auslassen. Kein Gate prüft die Überschrift.
**Verifizierbar:** Review-only

### INFO `spec/lastenheft.md` · „`--no-interactive` bricht bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab“ — Spannung zu Exit-Code `10` besteht weiter, Rang 2 trägt beide Seiten

**Klasse:** Vorbestehende Spannung, keine Änderung durch den Diff (F-7 des Vorlaufs des CR-Reviews)
**Quelle:** [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
**Befund:** Zeile 212 (Exit-Code `2`) steht unverändert neben Zeile 214 (Exit-Code `10`). In der Spezifikation stehen beide Seiten: die Verfeinerungen zu [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Abschnitte „Verhalten bei aktivierter Nicht-Interaktivität“ und „Auswertungslogik der Bestätigungsmodi“) sagen `2` bei fehlender Bestätigung unter `--no-interactive`, die Tabelle in [`SPEC-038`](../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes) sagt für `down --volumes` mit `--no-interactive` ohne `--yes` `10`. Die Ausnahme der Version `0.37.3` berührt diese Spannung nicht (ohne Container gibt es keine offene Frage); der Slice führt sie in §6 als Frage an Architect und Projektinhaber.
**Risiko:** Keines aus dem Diff. Der Pin-Test von `remove --volumes` mit Container (nicht-interaktiv ohne `--yes`) hat bis zur Entscheidung zwei Soll-Werte im Rang 1 und 2; der Slice trägt `10`, wie der Code von `u-boot down --volumes` und der Wortlaut von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
**Verifizierbar:** Review-only

### INFO `spec/lastenheft.md` · „bei `u-boot devcontainer remove --volumes` gilt das nur, wenn es einen Container zu entfernen gibt“ — Bedingung braucht den pinnenden Test des Umsetzungs-Slice

**Klasse:** Bedingung in der Spec ohne pinnenden Test (offen im Umsetzungs-Slice)
**Quelle:** `AGENTS.md` §3.10 (Bedingung in Spec oder Handbuch bekommt einen Test oder wird gestrichen)
**Befund:** Die Bedingung ist neu in Zeile 214 und dieselbe wie in Satz 10 des Spiegelstrichs in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten). Ein Test existiert nicht (Code des Befehls fehlt noch). Den Auftrag tragen der CR-Vorschlag (Abschnitt „Prüfung der Folgeänderung 0.37.3 gegen die Messgrenzen“, letzter Absatz) und Punkt (3) der DoD von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md); er nennt beide Seiten (ohne Container `0`, mit Container `10`) und die zwei unveränderten Befehle (`down --volumes`, `remove --purge`: `10`).
**Risiko:** keines im Diff; die Pflicht liegt beim Verifier des Umsetzungs-Slice.
**Verifizierbar:** Review-only (Test im Umsetzungs-Slice)

## Auflösung der Findings des Vorlaufs

| Vorlauf | Stand im Diff |
|---|---|
| MEDIUM (Ausnahme der Version `0.37.2` widerspricht Zeile 214 von [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) | Aufgelöst. Zeile 214 trägt jetzt die Ausnahme am Satz, der den Befehl namentlich nennt, und verweist auf [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten). Die Zeilen 214 und 517 (Satz 10) und 522 (Fehlerklasse) sagen denselben Fall gleich: ohne Container zu entfernen kein Exit-Code `10`, Exit-Code `0`. Der Rang-1-interne Widerspruch ist damit nicht mehr im Wortlaut. |
| LOW (Zustandsfeld `open/README.md` trägt die Versionskette) | Aufgelöst. Die Zeile nennt den Zustand (Lastenheft `0.37.3`, Wortlaut freigegeben, Nachlauf-Review offen, Closure offen) und verweist für die Fassungen auf den CR-Vorschlag; die Chronik steht dort einmal. Gleiche Form in der Welle (§1, §3, §5) und im Kopf des Slice. |
| INFO (Image-Satz ohne pinnenden Test) | offen, unverändert; Auftrag in der DoD von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md) Punkt (1) und in der Plan-Tabelle von [`slice-v2-spec-devcontainer-remove`](../plan/planning/open/slice-v2-spec-devcontainer-remove.md) |
| INFO (Spannung Exit-Code `2` und `10`) | offen, unverändert (INFO oben) |

## Geprüft ohne Befund

- **Wortgleichheit CR-Vorschlag und Lastenheft:** Stelle (p) (Absatz zu destruktiven Operationen) und Stelle (q) (Historie-Zeile `0.37.3`) 2 von 2 wortgleich; das Versionsfeld nennt `0.37.3`. Der Lastenheft-Diff trägt nichts außerhalb dieser drei Zeilen.
- **Halbsatz gegen die übrigen Sätze von Zeile 214:** Der erste Satz („Löschung nur über den expliziten Freigabepfad“) bleibt unbedingt; ohne Container gibt es keine Löschung, die Ausnahme betrifft allein den zweiten Satz („Exit-Code `10`“). `u-boot down --volumes` und `u-boot remove --purge` fallen nicht unter den Halbsatz; für `down` stützt das der gelesene Code (`Down` ruft `runConfirmationGate` nach der Prüfung des Projekts und der Compose-Datei und vor jedem weiteren Zugriff; Zeile 63), für `remove --purge` die unveränderte Zeile 668.
- **Halbsatz gegen die Fälle von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten):** Die Ausnahme „gilt das nur, wenn es einen Container zu entfernen gibt“ betrifft den Exit-Code `10` wegen fehlender Freigabe. Der Fall „ohne `--instance`, nur Container von Instanzen“ endet nach Satz 8 mit Exit-Code `10` aus einem anderen Grund (Hinweis auf `--instance`), unabhängig von `--yes`; beide Sätze widersprechen sich nicht. `--instance` ohne Container der Instanz (Exit-Code `0`) deckt Satz 10 („im Projekt oder, mit `--instance`, in der Instanz“).
- **Zeilen 210 und 212 gegen die Aussage im CR-Vorschlag („nennen den Exit-Code `10` nicht“):** gelesen, stimmt; Zeile 211 nennt `10` nur für `init --assume-existing`, Zeile 216 nur für `--force`/`--backup`, beide ohne Bezug zu `remove`.
- **Prüftabelle „weitere Stellen“ im CR-Vorschlag:** Zeilen 668 ([`LH-FA-ADD-007`](../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)), 710 ([`LH-FA-UP-004`](../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen)), 189, 1190, 517 und 522 gelesen; die Aussagen der Tabelle stimmen mit dem Lastenheft überein. Die Suche unter „Eigene Proben“ fand keine weitere Stelle mit unbedingtem Satz zu `remove --volumes` außer Zeile 517.
- **Fehlerklassen (Zeile 522) gegen Zeile 214:** „bei einem Container, den der Befehl entfernen würde“ und der Halbsatz sagen dieselbe Bedingung; Exit-Code `11` vor `10` (Containersuche braucht die Engine) bleibt gegeben.
- **Referenzrichtung (`AGENTS.md` §3.4, Skill HIGH #2):** Die hinzugefügten Lastenheft-Zeilen verweisen nur auf Anforderungen (`#lh-fa-dev-010--…`), nicht auf ADR, Slice, Carveout, Roadmap oder Dateien. Die Historie-Zeile nennt „zweiten Nachlauf-Review“ als Anlass in Worten, ohne Link.
- **Anforderungsebene des Wortlauts:** keine Labelnamen, Pfade, Docker-Befehle oder Flag-Interna außer der Bedienoberfläche (`--volumes`, `--yes`).
- **Aussagen des CR-Vorschlags über Code:** `DownService.Down` (Reihenfolge Gate vor Engine-Zugriff) gelesen; die Aussage zu `remove --purge` („Reihenfolge nicht Gegenstand dieser Folgeänderung“) behauptet nichts über den Code.
- **Zustandsmarker nach dem Commit:** Slice-Kopf, DoD-Zeilen 6 und 7, §1 (Folgeänderung `0.37.3`, Ausgang H-1), §6 (Risiko zu Zeile 214 und zur Spannung), `open/README.md`, Welle (§1, §3, §4, §5), CR-Status und Abweichungsabsatz nennen für `0.37.3` den Ist-Stand (freigegeben, im Lastenheft; Nachlauf-Review offen). Die Marker „Nachlauf-Review offen“ sind mit diesem Report überholt (siehe Übergabe). Die Zeilennummer „DoD-Zeile 7“ stimmt mit der Zählung der DoD (Zeile 7 = Nachlauf-Review zu `0.37.3`).
- **Slice §3 Plan-Tabelle gegen `git diff --stat 7f568768..HEAD`:** dieselben sieben Plan-Dateien, jede mit Zeile der Tabelle; die Zeilenangaben zum Lastenheft stimmen mit HEAD. Die achte Datei (Report des Vorlaufs) ist ein Report und kein Plan-Eintrag; die Beschreibung ihres Zustands ist LOW 1.
- **Folge-Slices:** [`slice-v2-spec-devcontainer-remove`](../plan/planning/open/slice-v2-spec-devcontainer-remove.md) §3 „Bezug zum CR“ nennt `0.37.3` samt der Regel in [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung); die Plan-Zeile zu `spec/spezifikation.md` nennt `0.37.2` nur dort, wo `0.37.2` die Aussage einführt (Image-Satz, Reihenfolge von Suche und Freigabe). `spec/spezifikation.md`, `spec/architecture.md` und Nutzer-Doku tragen keine Aussage zu `devcontainer remove` (Suche), also keine Kopie von Zeile 214 für den Befehl.
- **Lifecycle (`AGENTS.md` §3.3):** keine Renames; der Slice bleibt in `in-progress/`, die Folge-Slices bleiben in `open/`.
- **Kommentar-Selbstcheck (`AGENTS.md` §3.9):** kein Code- oder Skript-Kommentar im Diff.
- **Accepted-ADR, Nutzer-Doku, CHANGELOG, Carveouts:** nicht berührt (`AGENTS.md` §3.5 und §5 Regeln 5 und 6).
- **Dual-Classifier, Sprachvertrag, Dateisicherheit, Importregeln:** nicht berührt; der Diff enthält keinen Code.

## Nicht geprüft

- **Code, Tests, Docker-Pfad:** im Diff nicht berührt. Die Sentinel- und Diagnostic-Abdeckung von `remove` ist Aufgabe des Review von [`slice-v1-devcontainer-remove`](../plan/planning/open/slice-v1-devcontainer-remove.md).
- **Ob der Projektinhaber den Wortlaut `0.37.3` im Sitzungsdialog gesehen und freigegeben hat:** liegt im Sitzungsdialog und ist nicht belegbar; die Planungsstellen führen es als erfolgt.
- **`--json` ohne `--yes` bei `remove --volumes` ohne Container:** der Lastenheft-Wortlaut sagt „nicht-interaktiver Modus“ und nennt `--json` nicht; [`SPEC-038`](../../spec/spezifikation.md#spec-038--bestätigungspfad-von-u-boot-down---volumes) behandelt `--json` für `down`. Für `remove` legt es die Spezifikation fest (Aufgabe von [`slice-v2-spec-devcontainer-remove`](../plan/planning/open/slice-v2-spec-devcontainer-remove.md)); nicht Gegenstand des Diffs, daher kein Finding.
- **`make gates`:** nicht gefahren (Review-Rolle, Doku-Diff).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zustandsmarker mit Ablaufpunkt nicht entfernt (2×: Abgleichssatz, Überschrift der DoD-Zeile) · Vorbestehende Spannung (Exit-Code `2` und `10`) · Bedingung in der Spec ohne pinnenden Test (offen im Umsetzungs-Slice).

## Verdikt

**Merge-blockierend:** nein. Kein HIGH und kein MEDIUM; der Rang-1-interne Widerspruch des Vorlaufs ist im Wortlaut aufgelöst, der Wortlaut `0.37.3` ist wortgleich mit dem freigegebenen CR-Vorschlag. Die Wortlaut-Kette `0.37.0` bis `0.37.3` endet in diesem Lauf ohne weitere Änderung des Lastenhefts.

**Übergabe:**

- **Implementer:** Die zwei LOW sind behebbar ohne Wortlaut-Änderung. Nach diesem Report sind die Marker „weiterer Nachlauf-Review offen“ und „steht aus“ auf den Ist-Stand zu setzen: DoD-Zeile 7 des Slice (abhaken, Report [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-3`](2026-11-06-slice-cr-devcontainer-remove-nachlauf-3.md) nennen), §6 des Slice (Risiko „Der Wortlaut ändert sich nach dem Review“), §8 (Beobachtung `BEO-DOK/vertragsaenderung-nach-review`), Kopf des Slice, CR-Vorschlag (Kopf zu `0.37.3`, Abweichungsabsatz), `open/README.md` und Welle (§1, §5). Der Abgleich der Plan-Tabelle gegen `git diff --stat` ist danach zu wiederholen (`AGENTS.md` §3.10); der Report dieses Laufs ist eine neue Datei unter `docs/reviews/`, kein Plan-Eintrag.
- **Architect und Projektinhaber:** Die Frage Exit-Code `2` oder `10` für destruktive Operationen im nicht-interaktiven Modus ohne `--yes` (INFO oben) bleibt die einzige offene Rang-1-Frage zu `remove --volumes`; sie ist kein Gegenstand dieser Änderung. Wird sie entschieden, ändert sich Zeile 212 oder 214, und das ist eine weitere Wortlaut-Änderung mit Freigabe und Nachlauf-Review.
- **Steering-Loop:** `BEO-DOK/vertragsaenderung-nach-review` steht über der Schwelle (der Slice führt „Schwelle erreicht, Ausgang beim Projektinhaber“); dieser Lauf fand keine neue Wortlaut-Änderung, die Evidence-Datei zum Lauf zählt der Implementer bei der Closure. Die Klasse „Zustandsmarker mit Ablaufpunkt nicht entfernt“ trat in diesem Lauf in Form zweier LOW auf, nicht in den README-/Welle-Zeilen, deren Form seit dem Vorlauf zustandsbezogen ist.
- **Verifier:** Der Image-Satz (Vorlauf, INFO), die Bedingung zu Zeile 214 (INFO oben) und die Pin-Tests der Fälle (d) und (e) gehören in den Abgleich des Umsetzungs-Slice. Dieser Report belegt nur den Nachlauf, nicht die Closure.
