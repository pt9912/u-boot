# Slice CR: Lastenheft-Historie ohne Prozesszustand — die Klausel „Bestätigung ausstehend“ entfällt (Lastenheft 0.17.0, `LH-SA-CLI-002`, `LH-FA-CLI-007`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Release-Zuordnung: der nächste Release-Cut `v0.12.0`, der die Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) ausliefert (deren Closure ist geschlossen; der Slice gehört nicht zu ihr und hält sie nicht auf). Er hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er ist der **erste** einer verbindlichen Folge von drei Slices (1 → 2 → 3): erst der Vertrag (dieser Slice), dann die Spezifikation [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md), dann die Messung [`slice-v1-devcontainer-ollama-amd64-messung`](../open/slice-v1-devcontainer-ollama-amd64-messung.md).

**Bezug:** [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht; Historie-Zeile 0.14.0), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Wertebereich von `command`; Historie-Zeile 0.16.0), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Gegenstand der Zeile 0.14.0). Das Lastenheft führt seine Änderungen in [§16 Historie](../../../../spec/lastenheft.md#16-historie). Kein ADR.
Anlass: Äußerung des Projektinhabers vom 2026-10-19 — „Was hat „bestätigung ausstehend“ im lastenheft verloren“ und „specs sind Zielbild“. Ein Prozesszustand gehört in die Welle-Datei und den Slice, nicht in den Vertrag (`AGENTS.md` §3.7 *Zustandsfelder*, §3.10 *Zustandsmarker mit Ablaufpunkt*).

**Berührte Spec-Stellen:** — (das Lastenheft ist Rang 1 und ändert nur seine Historie; `SPEC-*` und `ARC-*` bleiben unverändert, gelesen in §3). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Historie-Zeilen 0.14.0 und 0.16.0 von `spec/lastenheft.md` tragen die Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“ beziehungsweise die Variante „aus dem freigegebenen Befehl abgeleitet; Bestätigung … ausstehend“ nicht mehr, ohne dass sich der Wortlaut einer Anforderung ändert; die neue Zeile 0.17.0 (nächste freie Minor-Version zum Zeitpunkt des Commits; zuletzt 0.16.0) nennt diese Änderung ohne Prozesszustand.

**Das Streichen ist keine Freigabe.** Die Wortlaute selbst — die Zeile `u-boot egress list` in der Befehlsübersicht ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Historie 0.14.0; Release-Gate 1) und der Wert `egress` im Wertebereich von `command` ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Historie 0.16.0; Release-Gate 2) — bleiben, was sie sind: aus dem freigegebenen Befehl abgeleitet, vom Projektinhaber nicht ausdrücklich bestätigt. Beide Gates bleiben in [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) §7 (*Auftrag an den Release-Cut*, Punkt 1) offen, bis der Projektinhaber ausdrücklich bestätigt oder ändert. Den Zustand „unbestätigt“ führen nach diesem Slice die Welle-Datei, die Roadmap-Zeile der Welle, die Bestandszeilen in `open/README.md` und dieser Slice; das Lastenheft führt ihn nicht.

### Gefundener Stand (gelesen am 2026-10-19; Docker-`grep` in `alpine:3` über `spec/`, `docs/`, `harness/`, `internal/`, `cmd/`, `CHANGELOG.md`, `README.md`, `README.de.md`, `AGENTS.md` nach den Mustern `Bestätigung … ausstehend`, `ausstehend … Bestätigung` und `aus dem freigegebenen Befehl`; die Treffer-Absätze gelesen)

Im Lastenheft (`spec/lastenheft.md`):

- **Versionsfeld im Kopf:** `0.16.0`.
- **§16 Historie, Zeile 0.14.0** (Datum 2026-10-18): endet mit „[`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): die Befehlsübersicht führt den Befehl mit Priorität V2 (aus dem freigegebenen Befehl abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend).“
- **§16 Historie, Zeile 0.16.0** (Datum 2026-10-19): „… (nach [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) für alle `--json`-Ausgaben; aus dem freigegebenen Befehl `u-boot egress list` abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend).“
- **Keine weitere Fundstelle im Lastenheft.** Die übrigen Treffer des Wortes „Bestätigung“ meinen die Rückfragen der CLI ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), `--purge`, Feature-Quellen), keine Klausel. Die Anforderungstexte selbst — die Tabellenzeile `u-boot egress list` mit Priorität V2 in [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und die Wertemenge in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) — tragen keinen Vermerk.

Außerhalb des Lastenhefts (Fundstellen nach Datei und Absatz):

| Fundstelle | Aussage (gelesen) | Einordnung |
|---|---|---|
| `docs/maintainer/cli-json-contract.md`, §6.10, Absatz „Vertrag:“ | „Der Wortlaut ist aus dem freigegebenen Befehl abgeleitet; die Bestätigung des Projektinhabers steht vor dem Release aus.“ | Zustandsmarker in einem Dokument mit Vertragsrang (Maintainer-Doku); wird gestrichen |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md` §3, Release-Gate 1 und Release-Gate 2 | Gate 1: „Das Lastenheft trägt die Klausel, bis er bestätigt; danach setzt eine Änderung genau dieser Klausel sie auf den Ist-Stand oder streicht sie“. Gate 2: „… oder die Historie-Zeile trägt keine Klausel „Bestätigung ausstehend“ mehr. Zustand: … (Historie-Zeile trägt die Klausel …)“ | **Kern des Risikos:** Die Alternative „keine Klausel mehr“ erfüllte Gate 2 dem Buchstaben nach durch das Streichen; sie ist nicht gemeint (siehe oben und §6). Die „Zustand:“-Zeilen tragen einen Marker, der mit diesem Slice abläuft |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md` §3, Slice-Tabelle, Zeile von `slice-cr-cli-json-command-enum` | „aus dem freigegebenen Befehl abgeleitet, Bestätigung des Projektinhabers: Release-Gate in §3“ | Gate-Verweis, bleibt wahr |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md` §7, Punkt 8 und *Auftrag an den Release-Cut* (1) | „das Lastenheft trägt die Klauseln … unverändert“; „… bis der Projektinhaber bestätigt oder die Klauseln im Lastenheft auf den Ist-Stand setzt“ | Marker läuft ab; die zweite Alternative ist die Lesart, die dieser Slice ausschließt |
| `docs/plan/planning/done/slice-cr-sandbox-egress-anzeige.md`, `done/slice-cr-cli-json-command-enum.md`, `done/slice-v1-cli-egress-anzeige.md`, `done/slice-v2-spec-sandbox-egress-ollama.md` | Beschreiben den Stand zum Zeitpunkt ihrer Closure (mit Commit-Hash und Datum, „der Lastenheft-Commit führt …“, „Stelle 4: Bestätigung offen, Gate vor Release“) | Chronik in `done/`; bleibt unverändert (`docs/plan/planning/README.md`: in `done/` nur Korrekturen); das Gate selbst bleibt offen, die Aussage „Gate vor Release“ bleibt wahr |
| `docs/reviews/` (sechs Berichte, acht Treffer: `2026-10-18-slice-v1-cli-egress-anzeige`, `2026-10-18-verification-slice-cr-sandbox-egress-anzeige` (drei Treffer), `2026-10-18-verification-slice-v1-cli-egress-anzeige`, `2026-10-19-slice-cr-cli-json-command-enum`, `2026-10-19-verification-slice-cr-cli-json-command-enum`, `2026-10-19-verification-slice-v2-spec-sandbox-egress-ollama`) | Zeitdokumente; Zitate des damaligen Stands, darunter die Auflage des Verifiers „danach Klausel … auf den Ist-Stand setzen oder streichen“ (V-1 der Verifikation zum Anzeige-CR) | unverändert (Zeitdokument) |
| `docs/plan/planning/open/README.md`, Bestandszeilen von `slice-cr-sandbox-egress-anzeige` und `slice-cr-cli-json-command-enum` | „Stelle 4 … wartet auf die Bestätigung des Projektinhabers, Gate vor dem Release“; „Bestätigung des Projektinhabers: Gate vor dem Release“ | benennen das Gate, nicht den Lastenheft-Zustand; bleiben wahr und unverändert |
| `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/`, `harness/`, `AGENTS.md`, `internal/`, `cmd/`, `spec/spezifikation.md`, `spec/architecture.md` | kein Treffer | unverändert; die Nutzer-Doku trägt keine Kennungen (`AGENTS.md` §5 Regel 6) |

Die Historie der Spezifikation (`spec/spezifikation.md` §7) nennt die Lastenheft-Version 0.16.0 in der Zeile zu [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) („deckungsgleich mit Lastenheft 0.16.0 und Code“): die Aussage bleibt wahr, weil der Wortlaut der Anforderungen unverändert ist.

### Wortlaut-Vorschlag (alt/neu)

Das Lastenheft nennt keine Slice- oder Wellen-Kennung; Verweise stehen als Anker `#…` der Datei. Der Wortlaut ist **kein** freigegebener Text der Anforderungen — er ändert nur die Buchführung; die Zustimmung zur Streichung ist die Äußerung des Projektinhabers vom 2026-10-19.

**Stelle 1 — §16 Historie, Zeile 0.14.0.** Alt (Ende der Zeile): „… die Befehlsübersicht führt den Befehl mit Priorität V2 (aus dem freigegebenen Befehl abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend).“ Neu: „… die Befehlsübersicht führt den Befehl mit Priorität V2.“

**Stelle 2 — §16 Historie, Zeile 0.16.0.** Alt: „… (nach [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) für alle `--json`-Ausgaben; aus dem freigegebenen Befehl `u-boot egress list` abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend).“ Neu: „… (nach [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) für alle `--json`-Ausgaben).“

**Stelle 3 — Versionsfeld und neue Historie-Zeile.** Versionsfeld `0.16.0` → `0.17.0`. Neue Zeile (Ist-Stand im Lastenheft: angelegt in Commit `2b137a6`, Endfassung nach der Fix-Runde in Commit `2100532`):

> | 0.17.0 | 2026-10-19 | §16 Historie: Die Zeilen 0.14.0 ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)) und 0.16.0 ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)) nennen die Änderung der Anforderung. Der Wortlaut der Anforderungen bleibt unverändert. | Vereinbarung mit dem Projektinhaber |

Die Version `0.17.0` ist die nächste freie; die Textkorrektur der Zeile 0.17.0 in der Fix-Runde hebt sie nicht an (0.17.0 ist nicht released).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Freigabe der beiden Wortlaute** — Entscheidung des Projektinhabers; der Slice streicht eine Buchführungs-Klausel, er bestätigt nichts. Die Gates 1 und 2 sind weder erfüllt noch entfallen.
- **Andere Prozesszustände oder Messstände im Lastenheft** (etwa „gemessen auf Docker 29.5.2 unter Colima 0.10.3“ in den Zeilen 0.12.0 und 0.13.0 oder im Absatz zu `ollama` in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) — der Projektinhaber hat die Klausel „Bestätigung ausstehend“ benannt; ein Messstand im Vertrag ist die gewollte „Zusage nur für Gemessenes“ ([`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)), kein Prozesszustand.
- **Spezifikation und Messstände dort** — [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md); die Spezifikation folgt dem Lastenheft.
- **Der Wortlaut der Anforderungen** (Befehlsübersicht, Enum, `subcommand`-Pflicht) — unverändert; eine Änderung wäre ein eigener CR nach Äußerung des Projektinhabers.
- **Umschreiben der geschlossenen `done/`-Dateien und Berichte** — Chronik; nur der Vermerk in der Welle-Datei (§3) wird gesetzt, weil dort das Gate selbst steht.

## 2. Definition of Done

- [x] **Lastenheft 0.17.0** (Commits `2b137a6` und `2100532`; die Endfassung der Zeile 0.17.0 steht seit `2100532`, vom Verifier gegen den Vorschlag gelesen): Stellen 1 bis 3 im Wortlaut dieses Vorschlags (Versionsfeld, Historie-Zeilen 0.14.0 und 0.16.0 ohne die Klausel, neue Zeile 0.17.0 mit Datum und Verweis „Vereinbarung mit dem Projektinhaber“); `git diff` des Lastenhefts zeigt allein diese Hunks (Wortgleichheit der übrigen Zeilen Satz für Satz gelesen, nicht nur per Werkzeug).
- [x] **Vorab-Suche und Nachzug** (Suche vor dem Lastenheft-Commit am 2026-10-19 wiederholt, Docker-`grep` in `alpine:3` über `spec/`, `docs/`, `harness/`, `internal/`, `cmd/`, `CHANGELOG.md`, `README.md`, `README.de.md`, `AGENTS.md`, ganzer Baum, nicht nur das Lastenheft: keine Fundstelle außerhalb der Tabelle in §1; Maintainer-Doku und Welle-Datei im Commit nach `2b137a6`): `docs/maintainer/cli-json-contract.md` §6.10 trägt den Marker nicht mehr; die Welle-Datei trägt den Vermerk, dass das Streichen keine Freigabe ist und die Gates 1 und 2 offen bleiben (datierte Vermerke in `done/`, der Wortlaut bleibt); die Roadmap-Zeile der Welle nennt den Gate-Stand ohne Bezug auf das Lastenheft.
- [x] **Review und Verifikation** (Review-Report `88bc15b`, Verification-Report `7a711cc`): Review-Report und Verification-Report unter `docs/reviews/` (`.harness/skills/reviewer.md`; Rollenwechsel nach Schritt 8, kein Self-Review). Der Verifier prüft ausdrücklich, dass kein Text — weder das Lastenheft noch die Welle-Datei noch die Roadmap — die Streichung als Bestätigung liest.
- [x] `make gates` grün; `make doc-trace` ohne neue Waise; `make doc-planning` und `make docs-check` ohne Befund. Läufe und Exit-Codes in §7.
- [x] Doku-Update: Handbuch, README und CHANGELOG bleiben unverändert (kein Produktverhalten, keine Kennungen in der Nutzer-Doku; gelesen in §1).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb für diesen Slice von der nächsten Welle-Closure geprüft (Stand in §7).

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff berührt; der Implementer gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien; „unverändert“ ist an der Fundstelle belegt.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 bis 3: Versionsfeld im Kopf, §16 Historie (Zeilen 0.14.0 und 0.16.0, neue Zeile 0.17.0); kein Anforderungstext |
| `docs/maintainer/cli-json-contract.md` | update | §6.10 Absatz „Vertrag:“: der Satz „Der Wortlaut ist aus dem freigegebenen Befehl abgeleitet; die Bestätigung des Projektinhabers steht vor dem Release aus.“ entfällt; der Rest des Absatzes (Enum, [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung), `jsontestutil`) bleibt, je Aussage neu gelesen |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md` | update (korrigierender Hinweis) | §3 Gate 1 und Gate 2 (Zeilen „Zustand:“) und §7 Punkt 8 sowie *Auftrag an den Release-Cut* (1): Marker „das Lastenheft trägt die Klausel“ auf den Ist-Stand, Hinweis auf Lastenheft 0.17.0 und dass die Alternative „keine Klausel mehr“ die Gates nicht erfüllt. Die Welle-Datei behält ihren Wortlaut und trägt je Stelle einen datierten Vermerk (*Vermerk 2026-10-19*); Begründung in §6 |
| `docs/plan/planning/in-progress/roadmap.md` | update | Arbeitszeile bei Start und Closure („Nichts in Arbeit.“ wieder herstellen); Tabellenzeile der Welle `welle-sandbox-egress-ollama` („zwei Release-Gates offen“): nennt bei der Closure beide Gates (Befehlsübersicht `egress list`, Enum `egress`; Verifier E-1); Verweis-Nachzug der Moves |
| `docs/plan/planning/open/README.md` | update | Bestandszeile dieses Slice (bei Anlage); Verweis-Nachzug von `make slice-mv` |
| diese Slice-Datei (Ort wechselt mit dem Lifecycle: `open/`, `next/`, `in-progress/`, `done/`) | update | Haken, Gefundener Stand (Wiederholung der Suche), Closure (§7) |
| `docs/reviews/` (Review-Report, Verification-Report) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/BEO-DOK/` (`vertragsaenderung-nach-review`, `vorab-suche-uebersieht-befehlsuebersicht`, `plan-aussage-ueber-code`: `state.md`, `evidence/`) | update, neu | Register bei der Closure (§7), soweit der Review ein Wiederauftreten belegt |
| Verweis-Nachzüge von `make slice-mv` (`open` → `in-progress` → `done`) | update (nur Verweise) | `open/README.md`, `roadmap.md` und jede Datei, die diesen Slice verlinkt (die beiden Folge-Slices) |
| `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `spec/spezifikation.md`, `spec/architecture.md`, `harness/` | unverändert | gelesen am 2026-10-19: `CHANGELOG.md` `[Unreleased]` (die Egress-Einträge nennen keinen Prozesszustand des Lastenhefts); `benutzerhandbuch.md` Abschnitt *Internet-Beschränkung (Egress)* zu `egress list` und die Tabelle der Pflichtfelder der `--json`-Ausgabe (Eintrag `command` mit `egress`): kein Vermerk zu einer Bestätigung; `spec/spezifikation.md` [`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung) und §7 Historie, `spec/architecture.md` (Abschnitte zu `EgressListService` und zu den Driving-Adaptern): keine Klausel; Docker-`grep` über `README.md`, `README.de.md`, `harness/`, `AGENTS.md`: kein Treffer |

**Abgleich (Implementer, Stand `e8498a2..HEAD` nach der Fix-Runde):** `git diff --stat` nennt `spec/lastenheft.md`, `docs/maintainer/cli-json-contract.md`, `docs/plan/planning/done/welle-sandbox-egress-ollama.md`, diese Slice-Datei und den Review-Report unter `docs/reviews/`; jede Datei hat einen Eintrag in der Tabelle. Einträge ohne Diff in dieser Phase: `docs/plan/planning/in-progress/roadmap.md` (Arbeitszeile bei Start schon in `e8498a2`; die Tabellenzeile der Welle gelesen: „zwei Release-Gates offen“ bleibt wahr; Zurücksetzen der Arbeitszeile bei der Closure), `open/README.md` und die Verweis-Nachzüge (mit dem Move bei Start erledigt; weitere beim Move nach `done/`), `docs/reviews/` und das Register (Review, Verifikation, Closure). Die als „unverändert“ geführten Dateien sind nicht berührt.

Die Vorab-Suche (Pflicht vor jeder Vertragsänderung; Lehre aus [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)) läuft über den **ganzen Baum**, nicht über das Lastenheft allein, und sucht die Klausel in ihren beiden Formen sowie den Satz „aus dem freigegebenen Befehl“; zusätzlich nach dem Gate-Wortlaut („keine Klausel“, „Ist-Stand“) in `docs/plan/planning/`.

## 4. Trigger

**Start** (`next` → `in-progress`, zuvor `open` → `next`): `Verantwortlich:` gesetzt; `in-progress/` enthält nur `README.md`, `roadmap.md` und `carveouts.md` (WIP-Limit 1); kein Vorgänger (erster Slice der Folge 1 → 2 → 3). Der Projektinhaber hat die Streichung mit der Äußerung vom 2026-10-19 beauftragt; eine Bestätigung der Wortlaute ist **keine** Start-Bedingung.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Suche findet eine Kopie der Klausel in einer Datei, die eine eigene Entscheidung verlangt (etwa ein Dokument mit eigener Freigabe), oder der Projektinhaber verlangt zusätzlich eine Änderung am Wortlaut einer Anforderung (eigener CR).
- `in-progress` → `open` (blockiert): Der Projektinhaber widerruft die Streichung oder verlangt vorher die Bestätigung beider Wortlaute; das Lastenheft bleibt unverändert.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). **Reihenfolge der Folge:** [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) startet erst, wenn dieser Slice in `done/` liegt; [`slice-v1-devcontainer-ollama-amd64-messung`](../open/slice-v1-devcontainer-ollama-amd64-messung.md) erst nach jenem. Die Gates 1 und 2 sind kein Closure-Kriterium dieses Slice (sie hängen am Release-Cut).

## 6. Risiken und offene Punkte

- **Das Streichen wird als Freigabe gelesen** — die Welle-Datei bietet als Alternative „oder die Historie-Zeile trägt keine Klausel … mehr“, die dem Buchstaben nach erfüllt wäre. — **Ausgang:** weiter offen: die Gates 1 und 2 bleiben Auftrag des Release-Cut; der Verifier hat die Lesart geprüft und keinen Text gefunden, der das Streichen als Freigabe liest (Verification-Report `7a711cc`, Abschnitt 3).
- **Das Gate verliert seinen Träger im Vertrag:** Nach der Streichung zeigt das Lastenheft nicht mehr, dass die Wortlaute unbestätigt sind; wer nur den Vertrag liest, sieht keinen Vorbehalt. — **Ausgang:** weiter offen: gewollt (Äußerung des Projektinhabers); der Träger ist die Welle-Datei §3 und §7 samt Vermerken und die Roadmap-Zeile der Welle, die beide Gates benennt; der Release-Cut liest beide.
- **Der Hinweis in einer `done/`-Datei überschreitet „nur korrigierend“** (`docs/plan/planning/README.md` §Lifecycle: „Inhalte in `done/` dürfen nachträglich nur korrigierend (Tippfehler, Querverweise, Archiv-Hinweise) geändert werden“). — **Ausgang:** eingetreten und behoben in der Fix-Runde: Die erste Fassung strich die Gate-2-Alternative und schrieb Gate 1 und §7 um; das ändert das Kriterium und ist keine Korrektur im Sinne der Regel. Jetzt bleibt der Wortlaut der Welle-Datei unverändert (Stand `e8498a2`), und je Stelle (Gate 1, Gate 2, §7 Punkt 8, *Auftrag an den Release-Cut* (1)) steht ein datierter Vermerk — ein Archiv-Hinweis, den die Regel nennt. Die Auslegung, ob der Vermerk, der die Alternative für unwirksam erklärt, ein Archiv-Hinweis oder eine Umdeutung des Gate-Kriteriums ist, liegt beim Projektinhaber (zur Entscheidung vorgelegt, §7); die Regel entscheidet sie nicht eindeutig.
- **Die Vorab-Suche übersieht eine Kopie** (zweites Auftreten im Register: [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)). — **Ausgang:** entfallen: Die Suche ist vor dem Lastenheft-Commit am 2026-10-19 wiederholt (Docker-`grep`, ganzer Baum), keine Fundstelle außerhalb der Tabelle in §1; der Review hat unabhängig gesucht (`docs/reviews/2026-10-19-slice-cr-lastenheft-zielbild.md`), die Gegenprobe des Verifiers bestätigt es.
- **Ein früherer Historie-Eintrag wird umgeschrieben, obwohl §16 ein Fußabdruck ist.** — **Ausgang:** getragen: Die Zeilen 0.14.0 und 0.16.0 sind gekürzt, die neue Zeile 0.17.0 hält die Änderung fest, der Anforderungstext ändert sich nicht; der Verifier hat den Diff gegengelesen (drei Hunks, wortgleich zum Vorschlag).
- **Weitere Zustandsmarker im Lastenheft** (Messstände in Historie und Anforderungstext). — **Ausgang:** weiter offen: kein Gegenstand dieses Slice (§1), Entscheidung des Projektinhabers; der Messstand gehört zur Zusage („Zusage nur für Gemessenes“).

## 7. Closure-Notiz

**Gegenstand und Commits.** Lastenheft 0.17.0: Die Historie-Zeilen 0.14.0 und 0.16.0 tragen die Klausel „Bestätigung ausstehend“ nicht mehr; die neue Zeile 0.17.0 nennt die Änderung ohne Prozesszustand. Der Wortlaut der Anforderungen ist unverändert. Commits: `2b137a6` (Lastenheft, Maintainer-Doku), `d73829e` (Vermerke in der Welle-Datei), `156729e` (DoD-Haken, Abgleich), `88bc15b` (Review-Report), `2100532` (Fix-Runde F-1 bis F-4), `7a711cc` (Verification-Report).

**Belege.** Review-Report [`2026-10-19-slice-cr-lastenheft-zielbild`](../../../reviews/2026-10-19-slice-cr-lastenheft-zielbild.md) (`88bc15b`; kein HIGH, kein MEDIUM, 4 LOW, 2 INFO; F-1 bis F-4 behoben in `2100532`, F-5 und F-6 übergeben beziehungsweise ausgewiesen) und Verification-Report [`2026-10-19-verification-slice-cr-lastenheft-zielbild`](../../../reviews/2026-10-19-verification-slice-cr-lastenheft-zielbild.md) (`7a711cc`; GO mit Auflagen). Die Auflagen A-1 bis A-6 sind mit dieser Closure abgearbeitet; V-1 (Lastenheft-Text nach dem Review geändert) steht im Register (unten), V-2 bis V-4 sind INFO.

**Stand der Gates 1 und 2 (Release-Gates, nicht Gate vor `done/`).** Beide sind **offen** und bleiben Auftrag des Release-Cut ([`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) §7, *Auftrag an den Release-Cut*, Punkt 1): Gate 1 betrifft die Zeile `u-boot egress list` in der Befehlsübersicht ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)), Gate 2 den Wert `egress` im Wertebereich von `command` ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run)). Das Streichen der Klausel im Lastenheft ist keine Bestätigung und erfüllt keines der beiden Gates; eine Bestätigung des Projektinhabers liegt nicht vor.

**Zur Entscheidung vorgelegt (Projektinhaber), nicht entschieden:**

1. **Auslegung der Vermerke in `done/` (E-2).** Die Welle-Datei trägt vier datierte Vermerke (Gate 1, Gate 2, §7 Punkt 8, *Auftrag an den Release-Cut* (1)); der Vermerk bei Gate 2 erklärt die Alternative „oder die Historie-Zeile trägt keine Klausel … mehr“ für unwirksam. Ob das ein Archiv-Hinweis im Sinne von `docs/plan/planning/README.md` §Lifecycle ist oder eine Umdeutung des Gate-Kriteriums, entscheidet die Regel nicht eindeutig. Empfehlung des Verifiers (Verification-Report, Abschnitt 3.2): die Vermerke bleiben; der Projektinhaber bestätigt einmal ausdrücklich, dass die Gates 1 und 2 nur durch seine Bestätigung oder Änderung der Wortlaute erfüllt werden und die Alternative „Klausel entfällt“ nicht gilt. Diese Bestätigung steht aus; Datum und Herkunft werden hier nachgetragen, sobald sie vorliegt.
2. **Verkörperung zu `BEO-DOK/vertragsaenderung-nach-review` (Schwelle erreicht).** Eine neue Hard Rule in `AGENTS.md` oder ein Sensor liegt beim Projektinhaber. Vorschlag aus dem Verification-Report (A-3): Ändert die Fix-Runde einen Vertragstext, löst das einen Nachlauf-Review über genau diesen Diff aus, oder der Verifier führt es als eigenen Prüfpunkt. Der Ausgang (`verkörpert` oder `geplant`) ist im Register nicht zugewiesen.
3. **Gates 1 und 2 (E-3):** bestätigen oder ändern; bei Änderung eines Wortlauts ein eigener CR mit Review des Diffs.

**Roadmap (E-1).** Die Tabellenzeile der Welle nennt jetzt beide offenen Gates (Befehlsübersicht `egress list`, Enum `egress`) und zeigt auf Welle-Datei §7; sie nennt das Lastenheft nicht.

**Was hat funktioniert:** Die Form „Vermerk statt Umschreiben“ hielt den Wortlaut der Welle-Datei byte-genau (Verifier: nur Anhänge am Zeilenende). Die Gegenprobe über den ganzen Baum fand keine Kopie der Klausel außerhalb von Chronik und Zeitdokumenten. Der Review fand die erste Fassung der Welle-Datei (Alternative gestrichen, Sätze umgeschrieben) als Überschreitung von „nur korrigierend“ (F-2).

**Was ging anders als geplant:** Die erste Fassung der Welle-Datei (`d73829e`) änderte das Gate-Kriterium; die Fix-Runde `2100532` stellte sie auf Vermerke um. Die Fix-Runde änderte außerdem den Lastenheft-Text (Zeile 0.17.0, F-4) nach dem Review; Review-Report und Haken bezogen sich auf den Vorzustand, den Endtext las erst der Verifier (V-1).

**Steering-Loop-Eintrag:** Lerneintrag als **geschärfte Regel (Kandidat, Schwelle erreicht), nicht verkörpert**: Eine Fix-Runde, die einen Vertragstext (Lastenheft, Rang 1) ändert, ist von keinem Review gedeckt; der Text braucht einen Nachlauf-Review über genau diesen Diff oder einen benannten Prüfpunkt des Verifiers. Register [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): drei Belege (`slice-cr-devcontainer-agent-features-sandbox`, `slice-cr-sandbox-egress-anzeige`, `slice-cr-lastenheft-zielbild`). Ob eine Regel oder ein Sensor entsteht, liegt beim Projektinhaber (oben, Punkt 2); ohne seine Entscheidung wird nichts in `AGENTS.md` verkörpert. Kein Feld `liegt in` in diesem Eintrag.

**Beobachtungs-Register (`../observations/`):** [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md): `evidence/slice-cr-lastenheft-zielbild.md` ergänzt (dritter Beleg), `state.md` nachgeführt. Gezählt nach den Evidence-Dateien: drei. Kein weiterer Eintrag: Der Verifier hat in diesem Slice kein Wiederauftreten von `BEO-DOK/plan-aussage-ueber-code`, `BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht` und `BEO-WZG/host-werkzeug-ausserhalb-docker` festgestellt.

**Folge-Slices:** [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) startet erst, nachdem dieser Slice in `done/` liegt (Reihenfolge 1 → 2 → 3); [`slice-v1-devcontainer-ollama-amd64-messung`](../open/slice-v1-devcontainer-ollama-amd64-messung.md) erst nach jenem. Beide liegen in `open/`. Neue Folge-Slices entstehen nicht.

**Risiken aus §6:** jedes trägt in §6 genau einen Ausgang: weiter offen (Streichen als Freigabe gelesen: Release-Gates; Gate ohne Träger im Vertrag: gewollt; weitere Zustandsmarker im Lastenheft: kein Gegenstand), eingetreten und behoben (Hinweis in `done/`; Auslegung beim Projektinhaber), entfallen (Vorab-Suche), getragen (Historie-Zeile).

**Grenzen:** (1) Die Wortlaute von [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) sind nicht bestätigt (Gates 1 und 2). (2) Die Wortgleichheit des Lastenheft-Diffs ist von Hand gelesen, nicht mit einem Zeichen-Diff-Werkzeug. (3) `make test-docker` und `make verify-depguard` liefen nicht (kein Produktcode, keine Importregel, kein Docker-Pfad im Diff). (4) Die Laufzeit-Ausgabe von `egress list --json` ist nicht neu gemessen. (5) Die Paarungen prüft die nächste Welle-Closure.

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: `slice-v2-spec-zielbild` und `slice-v1-devcontainer-ollama-amd64-messung` liegen in `open/`. Register: `vertragsaenderung-nach-review` ist ein Verzeichnis im Register mit drei nicht leeren Evidence-Dateien.

**Sensor-Läufe der Closure (Exit-Code des `make`-Prozesses selbst, Stand `7a711cc` vor den Closure-Commits):** `make gates` Exit 0 (`baseline-verify` 54 Dateien OK; `docs-check` 476 Dateien, 0 Befunde; Coverage 91,40 % gegen 90 %); `make doc-trace` Exit 0 (83 Anforderungen, 0 Waisen); `make docs-check` Exit 0 (476 Dateien, 0 Befunde); `make doc-planning` Exit 0 (0 Befunde); `make doc-immutable RANGE="e8498a2..HEAD"` Exit 0 (6 Commits, 0 Befunde). Der Lauf nach dem Move steht im Abgleich unten. Host-Werkzeuge: `git`, `grep`, `sed`, `python3` für Text- und Diff-Arbeit, nicht für Gates.

**Abgleich der Plan-Tabelle nach der Closure:** `git diff --stat e8498a2..HEAD` nach dem Move und dem Nachzug: 12 Einträge, jeder mit Zeile in der Plan-Tabelle §3 — `spec/lastenheft.md`, `docs/maintainer/cli-json-contract.md`, `done/welle-sandbox-egress-ollama.md` (Vermerke und Slice-Link), diese Slice-Datei, `roadmap.md` (Arbeitszeile, Gate-Halbsatz), Register (`state.md`, Evidence), die beiden Berichte unter `docs/reviews/` sowie `open/README.md`, `open/slice-v1-devcontainer-ollama-amd64-messung.md` und `done/slice-v2-spec-zielbild.md` (Verweis-Nachzug von `make slice-mv`). Jede Zeile mit Datei-Eintrag hat einen Diff; die als „unverändert“ geführten Dateien (`CHANGELOG.md`, READMEs, Handbuch, `spec/spezifikation.md`, `spec/architecture.md`, `harness/`) stehen nicht im Diff. Der Nachzug des Moves traf den Link der Welle-Datei auf diesen Slice und beide Berichte (`done/`-Ziel); `make docs-check` nach dem Move: Exit 0, 477 Dateien, 0 Befunde.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `spec/` (Lastenheft, DOK), `docs/maintainer/` und `docs/plan/` (DOK) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Sub-Area `spec/`, `harness/`, `docs/plan/` erfüllt die Schwelle aus zwei von drei Achsen laut Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (offen, 2×: ein Review über den vorgelegten Wortlaut deckt eine spätere Änderung nicht ab; der Review dieses Slice läuft über den Diff nach der Streichung), [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (offen, 2×; die Suche läuft über den ganzen Baum), [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (offen, zwei Belege; Welle-Datei, Roadmap und Maintainer-Doku zusammen nachziehen), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (beide verkörpert in `AGENTS.md` §3.10). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
