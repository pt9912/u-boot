# Verification-Report: slice-cr-lastenheft-zielbild — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-cr-lastenheft-zielbild` (`docs/plan/planning/in-progress/`), Lastenheft 0.16.0 auf 0.17.0, Commit-Range `e8498a2..HEAD` (5 Commits, HEAD `2100532`, 5 Dateien). Inhalts-Commits: `2b137a6` (Lastenheft, Maintainer-Doku), `d73829e` (Welle-Datei), `156729e` (Slice-Haken), `88bc15b` (Review-Report), `2100532` (Fix-Runde). Kein Produktcode, keine Spezifikation, keine Architektur-Sicht, keine ADR im Diff. <!-- d-check:ignore (Kennung als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-19-slice-cr-lastenheft-zielbild.md` (kein HIGH/MEDIUM; F-1 bis F-4 LOW, F-5, F-6 INFO), Slice §1 (Gefundener Stand, Wortlaut-Vorschlag), §2 (DoD), §3 (Plan-Tabelle), §6 (Risiken). <!-- d-check:ignore (Kennung als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Der Lastenheft-Diff hat genau die geplanten drei Hunks und steht wortgleich zum Vorschlag; kein Anforderungstext ist berührt. Kein Text im Baum liest das Streichen als Freigabe der Release-Gates 1 und 2. `make gates` und alle weiteren aufgetragenen Sensoren sind grün. Die Auflagen sind Closure-Arbeit des Planners (A-1 bis A-6) und zwei Punkte für den Projektinhaber (Abschnitt 4). Eine DoD-Verletzung liegt nicht vor; die Abweichung V-1 (Lastenheft-Text nach dem Review geändert) ist ein Register-Fall, kein Mangel des Textes.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | `EXIT=0` des `make`-Prozesses in der Log-Datei; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: 475 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — coverage 91.40% meets threshold 90%`; HEAD `2100532`; Arbeitsbaum vorher sauber |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ (keine neue Waise) |
| `make docs-check` | pass, Exit 0 | 475 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-planning` | pass, Exit 0 | 475 Dateien, 0 Befunde |
| `make doc-immutable RANGE="e8498a2..HEAD"` | pass, Exit 0 | 5 Commits aufgelöst, 0 Befunde (Accepted-ADRs unverändert) |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | kein Produktcode, keine Importregel, kein Docker-Pfad im Diff (`git diff --stat` ohne `internal/`, `cmd/`, Dockerfile, Makefile) |

## 2. Lastenheft-Diff und Wortlaut (Prüfpunkt 1)

`git diff e8498a2..HEAD -- spec/lastenheft.md`: **drei Hunks**, 3 Einfügungen, 2 Löschungen (Versionsfeld; die Zeilen 0.14.0 und 0.16.0 gekürzt; neue Zeile 0.17.0). Gegen Slice §1 „Wortlaut-Vorschlag“ von Hand gelesen:

| Stelle | Ergebnis |
|---|---|
| Stelle 1, Zeile 0.14.0 | endet „… die Befehlsübersicht führt den Befehl mit Priorität V2.“; wortgleich zu „Neu“ |
| Stelle 2, Zeile 0.16.0 | endet „(nach `LH-NFA-USE-004` für alle `--json`-Ausgaben).“; wortgleich zu „Neu“ (Link-Form: Anker der eigenen Datei, wie im Lastenheft üblich) | <!-- d-check:ignore (Kennung als Zitat) -->
| Stelle 3, Versionsfeld | `0.16.0` auf `0.17.0`, Spaltenbreite unverändert |
| Stelle 3, neue Zeile | `0.17.0`, `2026-10-19`, „§16 Historie: Die Zeilen 0.14.0 (…) und 0.16.0 (…) nennen die Änderung der Anforderung. Der Wortlaut der Anforderungen bleibt unverändert.“, Quelle „Vereinbarung mit dem Projektinhaber“; wortgleich zur Vorlage |
| Anforderungstext | `LH-SA-CLI-002` (Befehlstabelle) und `LH-FA-CLI-007` (Enum, Zeile 248 mit `egress`) stehen nicht im Diff | <!-- d-check:ignore (Kennung als Zitat) -->
| §3.4 / §3.7 | keine Slice-, Wellen-, ADR-Kennung, kein Prozesszustand, kein Abwärtsverweis; nur Anker der eigenen Datei; Indikativ, kein „nur noch“ mehr |

**Befund V-1 (Reihenfolge Review und Text):** Die Zeile 0.17.0 stand in `2b137a6` mit „führen nur noch die Änderung der Anforderung; …“. Der Review (`88bc15b`, F-4) hat genau diesen Text beanstandet; die Fix-Runde `2100532` hat den **Lastenheft-Text nach dem Review** geändert (`spec/lastenheft.md` 1 Zeile). Der Review-Report und der Slice-Haken („Commit `2b137a6`“) beziehen sich auf den Vorzustand; der Endtext entsteht erst in `2100532`. Der Endtext ist von mir gegen den Vorschlag gelesen und stimmt; ein Nachlauf-Review ist für diese Ein-Zeilen-Textkorrektur nicht nötig. Es ist aber ein weiteres Auftreten von `BEO-DOK/vertragsaenderung-nach-review` (siehe A-3).

## 3. Das Streichen ist keine Freigabe (Prüfpunkt 2)

### 3.1 Lesart der Texte

Welle-Datei `done/welle-sandbox-egress-ollama.md` (Diff `e8498a2..HEAD`: 4 Zeilen, **ausschließlich angehängte Vermerke**; der übrige Wortlaut jeder Zeile ist unverändert):

| Stelle | Vermerk | Lesart |
|---|---|---|
| §3 Gate 1 | „*Vermerk 2026-10-19:* Lastenheft 0.17.0 hat die Klausel aus der Historie gestrichen; das Streichen ist keine Bestätigung und erfüllt dieses Gate nicht, die Aussage ‚Das Lastenheft trägt die Klausel‘ gilt nicht mehr.“ | Gate bleibt offen; „Zustand: Bestätigung offen“ steht unverändert im selben Absatz |
| §3 Gate 2 | „… erfüllt dieses Gate nicht, die Alternative ‚oder die Historie-Zeile trägt keine Klausel … mehr‘ gilt nicht.“ | Gate bleibt offen; die Buchstaben-Lesart der Alternative ist ausdrücklich verneint |
| §7 Punkt 8 | „Lastenheft 0.17.0 hat die Klauseln gestrichen (siehe §3); das Streichen erfüllt die Gates nicht.“ | „offen, Bestätigung steht aus“ und „Diese Closure bestätigt sie nicht und entscheidet sie nicht“ stehen unverändert |
| *Auftrag an den Release-Cut* (1) | „Die Klauseln sind in Lastenheft 0.17.0 gestrichen; das Streichen ist keine Bestätigung, die Alternative ‚Klauseln auf den Ist-Stand setzen‘ erfüllt die Gates nicht.“ | „Der Cut schneidet nicht, bis …“ bleibt; der Release-Cut liest den Vermerk in demselben Punkt |

Weitere Texte:

- **Roadmap** (`in-progress/roadmap.md`): Tabellenzeile der Welle „Release-Zuordnung und zwei Release-Gates offen“, im Diff unverändert; kein Bezug auf das Lastenheft, kein „erledigt“. Die Arbeitszeile nennt den Slice als „In Arbeit“. Wahr.
- **`open/README.md`**, Bestandszeilen der beiden CR-Slices: „Gate vor dem Release“; wahr; die Zeile dieses Slice sagt ausdrücklich „das Streichen ist keine Freigabe“.
- **Slice** §1, §5, §6: konsistent („keine Freigabe“, „Gates 1 und 2 kein Closure-Kriterium“).
- **`docs/maintainer/cli-json-contract.md` §6.10**: der Satz zu Ableitung und Bestätigung entfällt; der Rest („`command: "egress"` liegt damit im Vertrag“) ist eine Aussage über das Enum, keine über die Freigabe. Kein Text liest dort eine Bestätigung.

**Ergebnis:** Kein Text lässt die Gates 1 und 2 als erledigt lesen. Restschwäche (kein Mangel): Die Hauptsätze der Gates tragen weiter den überholten Ist-Zustand („Das Lastenheft trägt die Klausel …“, „Zustand: … Historie-Zeile trägt die Klausel …“, „das Lastenheft trägt die Klauseln … unverändert“). Wer nur bis zum Satzende liest, liest einen falschen Stand; der Vermerk steht jeweils am Ende desselben Absatzes oder Punkts und verneint ihn ausdrücklich. Das ist der Preis der Vermerk-Form (Abschnitt 3.2) und im Slice §6 benannt.

### 3.2 Vereinbarkeit der Änderung in `done/` mit der Lifecycle-Regel

Regel (`docs/plan/planning/README.md` §Lifecycle): Inhalte in `done/` dürfen nachträglich nur korrigierend (Tippfehler, Querverweise, Archiv-Hinweise) geändert werden.

- **Form:** vereinbar. Der Wortlaut der Welle-Datei ist byte-genau erhalten (alle Änderungen sind Anhänge am Zeilenende, datiert, mit Verweis auf den Slice); es entsteht kein neues Artefakt und kein `git mv`. Das ist ein Archiv-Hinweis im Sinne der Regel. Die erste Fassung (`d73829e`) hatte dagegen die Alternative gestrichen und Sätze umgeschrieben; das war nicht korrigierend. Die Fix-Runde `2100532` hat das behoben (Review F-2).
- **Inhalt, die offene Frage des Implementers:** Der Vermerk erklärt die Alternative von Gate 2 („oder die Historie-Zeile trägt keine Klausel … mehr“) für unwirksam. Streng gelesen engt ein Vermerk, der eine Alternative eines Abnahmekriteriums für unwirksam erklärt, das Kriterium ein; das ist Auslegung, nicht Tippfehler. Dagegen spricht: (a) Gate 1 im selben Absatz führt die Streichung der Klausel ausdrücklich erst **nach** der Bestätigung („danach setzt eine Änderung genau dieser Klausel sie auf den Ist-Stand oder streicht sie“); die Absicht beider Gates ist die Bestätigung des Projektinhabers, die Klausel ist ihr Anzeiger, nicht ihr Ersatz; (b) das Gate-Kriterium im Kopf („Der Projektinhaber hat den Wortlaut … bestätigt“) bleibt wörtlich stehen; (c) der Vermerk setzt kein neues Kriterium, er schließt eine Lesart aus, die der Projektinhaber mit dem Auftrag „specs sind Zielbild“ nicht gemeint hat (Slice §1, Anlass). Die Regel entscheidet die Auslegungsfrage nicht eindeutig.

**Empfehlung an den Projektinhaber:** Die Vermerke **bleiben** (Form regelkonform, Wirkung gewollt, Alternativen schwächer: ein Hinweis nur in der Roadmap-Zeile würde den Release-Cut nicht erreichen, der die Welle-Datei liest; ein Ersatzartefakt widerspräche der Regel „kein neues Artefakt für eine Korrektur“ nicht, trüge das Gate aber an zwei Orten). Der Projektinhaber bestätigt **einmal ausdrücklich**, dass Gate 2 nur durch seine Bestätigung oder Änderung des Wortlauts erfüllt wird und die Alternative „Klausel entfällt“ nicht gilt; der Planner hält diese Bestätigung mit Datum und Herkunft (Gespräch) in der Closure-Notiz §7 des Slice fest. Damit ist die Auslegung keine Umdeutung des Projektinhabers durch Dritte, sondern seine eigene Setzung.

### 3.3 Verbleibende Träger des Vorbehalts

Das Lastenheft trägt keinen Vorbehalt mehr (gewollt, Slice §6). Träger sind: Welle-Datei §3 und §7 samt Vermerken, Roadmap-Zeile („zwei Release-Gates offen“, ohne Benennung der beiden), `open/README.md`-Bestandszeilen, die Done-Slices der beiden CR (Chronik, „Gate vor Release“ bleibt wahr), dieser Slice. Der Release-Cut liest Welle-Datei §7 *Auftrag an den Release-Cut*; die Roadmap-Zeile benennt die Gates nicht und verweist über die Welle-Datei (Empfehlung E-1 unten, optional).

## 4. Baum-Suche (Prüfpunkt 3)

Gesucht in `alpine:3` (Docker, kein Host-Go) über `spec/`, `docs/`, `harness/`, `internal/`, `cmd/`, `CHANGELOG.md`, `README.md`, `README.de.md`, `AGENTS.md`, `.harness/skills`; Muster: `Bestätigung … ausstehend`, `ausstehend … Bestätigung`, `aus dem freigegebenen Befehl`, `keine Klausel`, `Klausel Bestätigung`, `steht vor dem Release aus`; zweiter Lauf mit `unbestätigt`, `nicht bestätigt`, `Bestätigung offen`, `abgeleitet`, `Bestätigung des Projektinhabers`.

- **Kein Treffer** in `spec/`, `harness/`, `internal/`, `cmd/`, `CHANGELOG.md`, den READMEs, `AGENTS.md`, `.harness/skills`, `docs/user/`, `docs/maintainer/`, ADRs. Die Klausel steht auch nicht mehr in `cli-json-contract.md`.
- **Treffer nur in** `docs/plan/planning/`: dieser Slice (Zitate im Gefundenen Stand), `open/slice-v2-spec-zielbild.md` (Anlass-Zitat, Chronik), `open/README.md` (Gate-Verweise), `done/` (Welle-Datei mit Vermerken, `slice-cr-cli-json-command-enum`, `slice-cr-sandbox-egress-anzeige`, `slice-v1-cli-egress-anzeige`, `slice-v2-spec-sandbox-egress-ollama`; Chronik) und `docs/reviews/` (Zeitdokumente, ausgenommen). Alle Treffer sind Chronik, Gate-Verweise oder Zitate; keiner behauptet eine geltende Klausel im Vertrag. <!-- d-check:ignore (Kennung als Zitat) -->
- Weitere Treffer von „Bestätigung“ im Lastenheft und in `spec/spezifikation.md` meinen Rückfragen der CLI, keine Klausel.

Nicht aufgeführt im Gefundenen Stand des Slice (Review F-6 a, Register-Zeitdokument, bleibt unverändert zu Recht): `observations/BEO-DOK/vertragsaenderung-nach-review/evidence/slice-cr-sandbox-egress-anzeige.md`. Kein Mangel. <!-- d-check:ignore (Kennung als Zitat) -->

## 5. `cli-json-contract.md` §6.10 gegen Code (Prüfpunkt 4)

Diff: ein Hunk, der Satz „Der Wortlaut ist aus dem freigegebenen Befehl abgeleitet; die Bestätigung des Projektinhabers steht vor dem Release aus.“ entfällt, der Absatz bleibt grammatisch heil. Gelesen gegen den Code:

| Aussage in §6.10 | Fundstelle | Ergebnis |
|---|---|---|
| `command: "egress"`, `subcommand: "list"` | `internal/adapter/driving/cli/egress.go:158` `newDataEnvelope("egress", "list", toEgressListData(resp), nil, 0)` | stimmt |
| bare `u-boot egress --json`: RunE-Reject ohne Envelope, `ErrEgressSubcommandRequired` | `egress.go:21`, `:73` | stimmt |
| `--dry-run`/`--diff`: `ErrEgressPreviewNotApplicable`, Exit 2 | `egress.go:25`, `:144` | stimmt |
| Enum in `LH-FA-CLI-007` nennt `egress` | `spec/lastenheft.md:248` | stimmt | <!-- d-check:ignore (Kennung als Zitat) -->
| `SPEC-001` führt `egress` im Enum | `spec/spezifikation.md:743` | stimmt | <!-- d-check:ignore (Kennung als Zitat) -->
| `jsontestutil` prüft das Enum nicht, `checkSubcommand` verlangt `subcommand` nur für `template` und `config` | `jsontestutil.go:237-245` | stimmt |

Kein Befund.

## 6. DoD-Abgleich (Slice §2), Punkt für Punkt

| DoD-Zeile | Haken | Stand | Beleg |
|---|---|---|---|
| Lastenheft 0.17.0 (Stellen 1 bis 3) | gesetzt | **erfüllt, Haken trägt mit Berichtigung** | Abschnitt 2; der Haken nennt „Commit `2b137a6`“, die Endfassung der Zeile 0.17.0 entsteht in `2100532` (A-2) |
| Vorab-Suche und Nachzug | gesetzt | **erfüllt** | Abschnitt 4 (eigene Gegenprobe, ganzer Baum); §6.10 ohne Marker (Abschnitt 5); Welle-Datei mit datierten Vermerken, Wortlaut erhalten (3.1); Roadmap-Zeile ohne Lastenheft-Bezug, unverändert und wahr |
| Review und Verifikation | offen | **erfüllt mit diesem Report**; Haken beim Planner | Review-Report liegt vor; dieser Report prüft ausdrücklich Lastenheft, Welle-Datei, Roadmap (Abschnitte 2 und 3) |
| `make gates` grün; `make doc-trace`; `make doc-planning`; `make docs-check` | offen | **erfüllt** (Abschnitt 1); Läufe und Exit-Codes in §7 nachzutragen | Haken und Eintrag beim Planner |
| Doku-Update: Handbuch, README, CHANGELOG unverändert | offen | **erfüllt** | keine der Dateien im Diff; Suche ohne Treffer (Abschnitt 4) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | offen | offen, zu Recht | §7 trägt den Platzhalter (A-1) |
| Beobachtungs-Register fortgeschrieben | offen | offen, zu Recht | `observations/` nicht im Diff (A-3) |
| Jedes Risiko aus §6 trägt einen Ausgang | offen | offen, zu Recht | Stand in A-4 |
| Drei Paarungen | offen | offen, zu Recht | Welle-Closure prüft sie (A-6) |

Entscheidungs-Konformität: keine Accepted-ADR im Diff (`AGENTS.md` §3.5, `make doc-immutable` 0 Befunde); keine Gate-Lockerung (§3.6); kein `git mv` im Diff (§3.3 nicht berührt); Kommentar-Muster (§3.7, §3.9) in den geänderten Zeilen gelesen: Indikativ, kein „früher“/„nur noch“ im Lastenheft. Hinweis: Die Vermerke in der Welle-Datei sind Aussagen über den Ist-Stand mit Datum („Vermerk 2026-10-19“, „gilt nicht mehr“); „gilt nicht mehr“ in Gate 1 ist das Muster §3.9 (a) in einem datierten Archiv-Vermerk, dessen Gegenstand der Zustandswechsel ist; ich werte es als zulässig (Chronik-Funktion des Vermerks), nicht als Befund.

## 7. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff --stat e8498a2..HEAD`, `AGENTS.md` §3.10)

5 Dateien im Diff (Stand `2100532`); mit diesem Report: 6.

| Diff-Datei | Plan-Zeile | Befund |
|---|---|---|
| `spec/lastenheft.md` (7 Zeilen) | Zeile 1 | ok (Abschnitt 2) |
| `docs/maintainer/cli-json-contract.md` (3 Zeilen) | Zeile 2 | ok (Abschnitt 5) |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md` (8 Zeilen) | Zeile 3 | ok; Plan sagt „datierter Vermerk je Stelle“, Diff bestätigt vier Anhänge |
| diese Slice-Datei | Zeile „diese Slice-Datei“ | ok |
| `docs/reviews/2026-10-19-slice-cr-lastenheft-zielbild.md` (neu) | Zeile `docs/reviews/` | ok | <!-- d-check:ignore (Kennung als Zitat) -->
| `docs/reviews/2026-10-19-verification-slice-cr-lastenheft-zielbild.md` (dieser Report) | Zeile `docs/reviews/` | ok | <!-- d-check:ignore (Kennung als Zitat) -->

Einträge ohne Diff, begründet im „Abgleich“-Absatz des Slice: `roadmap.md` (Arbeitszeile stammt aus `e8498a2`; Zurücksetzen bei der Closure), `open/README.md` und Verweis-Nachzüge (mit dem Move bei Start erledigt, weitere beim Move nach `done/`), Register (Closure). Die „unverändert“-Zeilen (`CHANGELOG.md`, READMEs, Handbuch, `spec/spezifikation.md`, `spec/architecture.md`, `harness/`) sind nicht im Diff; die Fundstellen der Tabelle decken sich mit meiner Suche (Abschnitt 4). Der Abgleich nach der Fix-Runde steht im Slice (Absatz „Abgleich (Implementer, Stand …)“) und stimmt mit dem Diff überein. Er ist nach der Closure (Move, Nachzüge, Register) zu wiederholen (A-5).

## 8. Review-Findings F-1 bis F-6 einzeln

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 LOW absolute Aussage „steht nur noch hier …“ | **behoben** | Die Aussage ist weder in der Welle-Datei noch im Slice (§1) als „nur noch“ vorhanden; Slice §1 nennt als Träger Welle-Datei, Roadmap-Zeile, `open/README.md`-Bestandszeilen und den Slice (ohne Alleinanspruch); Suche nach „nur noch“ in den geänderten Dateien ohne diese Aussage |
| F-2 LOW Änderung in `done/` an der Grenze „nur korrigierend“, Entscheidung nicht festgehalten | **behoben, Auslegungsrest übergeben** | Welle-Datei: Wortlaut erhalten, nur datierte Vermerke (Abschnitt 3.2); Slice §6 hält den Vorgang als „eingetreten und behoben“ und die Auslegungsfrage als offen für den Projektinhaber fest. Die Entscheidung des Projektinhabers steht aus (E-2) |
| F-3 LOW Zustandsmarker im Slice (Datumsplatzhalter, „weiter offen bis …“, doppelter Haken) | **behoben** | Platzhalter durch `2026-10-19` und Commit ersetzt; die Risiko-Ausgänge „Hinweis in `done/`“, „Vorab-Suche“, „Historie umgeschrieben“ sind auf Ist-Stände gesetzt; der Haken der Vorab-Suche führt die Aussage einmal. Rest: Ausgang „getragen … Gegenprüfung durch den Verifier steht aus“ ist mit diesem Report abgelaufen (A-4) |
| F-4 LOW „nur noch“ in der Zeile 0.17.0 | **behoben** | Lastenheft-Zeile nennt „nennen die Änderung der Anforderung. Der Wortlaut der Anforderungen bleibt unverändert.“ (Abschnitt 2); Text nach dem Review geändert (V-1) |
| F-5 INFO Gate-Gegenstand ohne Träger im Vertrag | **übergeben, getragen** | Slice §6 (Risiko 2, „weiter offen: gewollt“); Vermerke tragen den Vorbehalt (Abschnitt 3); Roadmap-Zeile siehe E-1 |
| F-6 INFO (a) Register-Evidence nicht in der Suche-Tabelle; (b) Link `../in-progress/…` in der Welle-Datei bricht beim Move | **(a) ausgewiesen, kein Mangel; (b) an die Closure übergeben** | (a) Zeitdokument, Abschnitt 4. (b) Der Link steht in `done/welle-sandbox-egress-ollama.md` und zeigt auf `../in-progress/slice-cr-lastenheft-zielbild.md`; `make slice-mv ... TO=done` zieht Verweise nach, `make docs-check` ist der Wächter (A-5) | <!-- d-check:ignore (Kennung als Zitat) -->

## 9. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | LOW | Der Lastenheft-Text (Zeile 0.17.0) wurde in der Fix-Runde nach dem Review geändert; Review-Report und Haken „Commit `2b137a6`“ beziehen sich auf den Vorzustand. Endtext verifiziert (Abschnitt 2). Register-Fall `BEO-DOK/vertragsaenderung-nach-review`: nach den bisherigen zwei Belegen (Slice §8) ist dies der dritte, die 3×-Schwelle ist erreicht. |
| V-2 | INFO | Hauptsätze der Gates 1 und 2 und §7 Punkt 8 in der Welle-Datei beschreiben den überholten Ist-Stand („Das Lastenheft trägt die Klausel“); der Vermerk im selben Absatz verneint ihn. Kein Mangel der Vermerk-Form; Lesbarkeitsrest, vom Slice §6 getragen. |
| V-3 | INFO | Die Roadmap-Zeile der Welle sagt „zwei Release-Gates offen“ ohne zu benennen, welche; die Benennung steht nur in der Welle-Datei. Wahr, aber der Träger ist dünn (E-1). |
| V-4 | INFO (Release-Gate, kein Mangel) | Gates 1 und 2 (Befehlsübersicht `egress list`, Enum `egress`) sind beim Projektinhaber offen, nicht Closure-Kriterium dieses Slice (§5), Auftrag des Release-Cut. Nicht im Verifier-Zugriff. |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, keine Abweichung vom Plan-Wortlaut.

## 10. Auflagen

**Beim Planner vor `done/`:**

- **A-1 Closure-Notiz §7** füllen: Läufe und Exit-Codes (Abschnitt 1), Stand der Gates 1 und 2 („offen, Auftrag des Release-Cut; das Streichen ist keine Bestätigung“), Steering-Loop-Eintrag, Folge-Slices (`slice-v2-spec-zielbild` startet erst nach `done/`, Reihenfolge der Folge), Herkunft und Datum der Bestätigung des Projektinhabers zu E-2. <!-- d-check:ignore (Kennung als Zitat) -->
- **A-2 Haken berichtigen:** DoD-Zeile „Lastenheft 0.17.0“ und Stelle 3 nennen „Commit `2b137a6`“; die Endfassung der Zeile 0.17.0 entsteht in `2100532` (V-1). Beide Commits nennen.
- **A-3 Register:** `BEO-DOK/vertragsaenderung-nach-review` fortschreiben (dritter Beleg: Lastenheft-Text in der Fix-Runde nach dem Review geändert, vom Verifier gegen den Vorschlag gelesen); die 3×-Schwelle ist erreicht, der Steering-Loop-Eintrag in §7 entscheidet, ob eine Regel verkörpert wird (Vorschlag: Änderung eines Vertragstextes in der Fix-Runde löst einen Nachlauf-Review über genau diesen Diff aus oder wird vom Verifier als eigener Prüfpunkt geführt). Mitzählen: `BEO-DOK/plan-aussage-ueber-code` und `BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht` nur, wenn der Planner sie im Lauf belegt findet; ich habe in diesem Slice kein Wiederauftreten festgestellt.
- **A-4 Risiken §6** auf Ist-Stände setzen: Risiko 1 („als Freigabe gelesen“): weiter offen, Gates beim Release-Cut, Verifier hat die Lesart geprüft (kein Text liest eine Freigabe); Risiko 2 (Gate-Träger): weiter offen, gewollt; Risiko 3 (`done/`-Hinweis): eingetreten und behoben, Auslegung beim Projektinhaber; Risiko 4 (Vorab-Suche): entfallen, durch eigene Gegenprobe des Verifiers bestätigt; Risiko 5 (Historie umgeschrieben): „Gegenprüfung steht aus“ ist abgelaufen, auf „getragen, vom Verifier gegengelesen“ setzen; Risiko 6: weiter offen.
- **A-5 Closure-Mechanik:** Roadmap-Arbeitszeile auf „Nichts in Arbeit.“ zurücksetzen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3, `make slice-mv`); danach `make docs-check` (der Link `../in-progress/…` in der Welle-Datei und die Verweise in `open/README.md`, `slice-v2-spec-zielbild`, `slice-v1-devcontainer-ollama-amd64-messung`, Roadmap müssen nachgezogen sein); Plan-Tabelle nach der letzten Änderung erneut gegen `git diff --stat` abgleichen (§3.10). <!-- d-check:ignore (Kennung als Zitat) -->
- **A-6 Restliche Haken** setzen (Review und Verifikation, Sensoren, Doku-Update, Closure, Register, Risiken); die Paarungen (Anker · Folge-Slice · Register) von der nächsten Welle-Closure prüfen lassen, wie der Slice es festhält.

**Für den Projektinhaber (keine Mängel):**

- **E-1 (optional):** Die Roadmap-Zeile der Welle benennt die beiden offenen Gates nicht („zwei Release-Gates offen“). Da das Lastenheft den Vorbehalt nicht mehr trägt, genügt ein Halbsatz „(Gate 1 Befehlsübersicht, Gate 2 Enum; Welle-Datei §3)“, damit der Release-Cut den Träger findet. Das ist eine Entscheidung über die Roadmap, keine Voraussetzung dieser Closure.
- **E-2:** Auslegung der Vermerke in `done/` bestätigen (Abschnitt 3.2): Empfehlung **Vermerke bleiben**; der Projektinhaber bestätigt einmal, dass die Gates 1 und 2 nur durch seine ausdrückliche Bestätigung oder Änderung der Wortlaute erfüllt werden.
- **E-3 (Release-Gate):** Gates 1 und 2 bestätigen oder ändern. Bei Änderung eines Wortlauts: eigener CR mit Review des Diffs.

## 11. Nicht verifiziert

Die Bestätigung der Wortlaute (Gates 1 und 2, Entscheidung des Projektinhabers); Zeichen-Diff-Werkzeug (Wortgleichheit von Hand gelesen); die Laufzeit-Ausgabe von `egress list --json` (keine Tests, nur Code und Doku gelesen); `make test-docker`, `make verify-depguard` (kein Produktcode, keine Importregel im Diff); Inhalt der abgeschlossenen Slices und der Review-Reports über die Treffer der Suche hinaus (Zeitdokumente); `make docs-check` nach Anlage dieses Reports: 476 Dateien, 0 Befunde.
