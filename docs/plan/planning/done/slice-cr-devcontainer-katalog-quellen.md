# Slice CR: CR-Schnitt Katalogquellen — Entfernung `cpp`/`postgres-client` und Abweisung unverarbeitbarer Formen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`; er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** ohne Welle-Datei (Rest nach dem Release `v0.20.0`); Anlass: die Messung [`docs/messungen/slice-v1-katalog-features-quellen-und-freigabe/`](../../../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md) und die zwei Entscheidungen des Projektinhabers vom 2026-11-03 (N1 entfernen, N2 Abweisung). Muster [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) (*Nicht Gegenstand* (a)/(b) — die Messung trägt die Entlassung).

**Verantwortlich:** pt9912 (Implementer; Wortlaut-Vorlage und Vertrags-Übernahme im Auftrag des Hauptagenten)

**Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** der freigegebene CR-Wortlaut ([`cr-vorschlag-katalog-quellen.md`](cr-vorschlag-katalog-quellen.md)) — die Entfernung der zwei Katalog-Einträge samt der Vertrags-Beispielzeile und die Abweisung unverarbeitbarer externer Quell-Formen je Exit-Code `10` — wird je Stelle freigegeben und wortgleich übernommen. Kein Produktcode in diesem CR (Muster der CR-Slices).

**Nicht in diesem Slice:** die Umsetzung (Katalog-Entfernung im Code samt Tabellentesten, Generator-Abweisung samt Klassifikatoren, Doku-Nachzug) — Umsetzungs-Slice nach der CR-Closure; die Tarball-Probe (bleibt Offenpunkt mit Historie-Vermerk).

## 2. Definition of Done

- [x] Vorschlags-Datei mit Alt/Neu je Stelle (a)–(d) samt Mess-Basis
- [x] Freigabe des Projektinhabers je Stelle (oder mit benannten Änderungen)
- [x] Vertrags-Übernahme wortgleich ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Beispiele-Zeile + Abweisungssätze; [`SPEC-023`](../../../../spec/spezifikation.md)/[`SPEC-040`](../../../../spec/spezifikation.md); Historie 0.33.0; Versionsfeld), Abgleich je Stelle in §7
- [x] `make gates` grün; Review + Verification mit Reports (`docs/reviews/`); Nachlauf-Review über jeden Diff von Vertragstext ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)); Closure-Notiz mit Lerneintrag; Register oder „keine Beobachtung"; je Risiko ein Ausgang; drei Paarungen getragen

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `cr-vorschlag-katalog-quellen.md` | neu | DoD-Zeile 1 |
| `spec/lastenheft.md` | update (nach Freigabe) | Stellen (a), (b), (d); Versionsfeld 0.33.0 |
| `spec/spezifikation.md` | update (nach Freigabe) | Stelle (c) — Katalogtabelle, Form-Tabelle, Pipeline-Schritt 6; Historie |
| dieser Slice, [`roadmap.md`](../in-progress/roadmap.md), `open/README.md` (Bestand) | update | Lifecycle je `make slice-mv` |
| `docs/reviews/` (Berichte) | neu | je Berichte-Runden |
| Produktcode, Templates, Tests, Handbuch, READMEs, `CHANGELOG.md` | **kein Diff in diesem CR** | der Umsetzungs-Slice trägt sie |

**Ansatz-Zeile (§3.10-Berichtigung nach dem Review, F-Findings zu den Move-Trägern):** die Lifecycle-Moves des Vorgangs zogen je Schritt die eingehenden Verweis-Träger mit `make slice-mv` nach (CR-Schnitt: ~22 Träger-Dateien über Planning-Abbage, Berichte und Roadmap — reine Moves samt Verweis-Nachzug je eigenem Commit; der §3.10-Abgleich gegen `git diff --name-status 137f2f92^..HEAD` läuft je Fix-und-Closure-Runde erneut).

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`.

**Rückführungen:** `in-progress` → `next` (zu groß — etwa ein dritter Weg); `in-progress` → `open` (blockiert — der Projektinhaber zieht einen Weg zurück; die Messung bliebe als Beleg).

## 5. Closure-Trigger

DoD vollständig; Vertrags-Diff wortgleich je Stelle; `make gates` grün; Reports; Closure-Notiz (§7).

## 6. Risiken und offene Punkte

1. Wortlaut nicht wortgleich übernommen — Abgleich je Stelle in §7 (Muster der CRs)
2. Eine Zusage geht weiter als die Messung (die Form-Tabelle trägt nur das Gemessene; Tarball-Weg als Grenze) — Vorsorge: Messbasis der Historie-Zeile; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
3. Eine Abweisung ohne Test-Auftrag — Beleg-Aufträge je Stelle (Muster [`SPEC-054`](../../../../spec/spezifikation.md)); [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
4. Plan-Aussagen über Code ohne Fundstelle — §3.10-Abgleich nach jeder Fix-Runde; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)

## 7. Closure-Notiz

### Closure (Planner, 2026-11-03, nach Review-Fix-Runden (`b122f697`, `283f815b`, `0637cb37`) und Vertrags-Übernahme `151a7976`)

**Evidence:** Review-Report [`2026-11-03-slice-cr-devcontainer-katalog-quellen`](../../../reviews/2026-11-03-slice-cr-devcontainer-katalog-quellen.md) (0 HIGH/4 MEDIUM/6 LOW/2 INFO je berichtigt) · Verify-Report [`…-verify`](../../../reviews/2026-11-03-slice-cr-devcontainer-katalog-quellen-verify.md) (DoD-Zeilen 1/2/5 belegt erfüllt; V-1 das Historie-Fragment — die gebrochene Edit-Kante der Umsetzungs-Runde — berichtigt vor der Closure; Wortlaut-Abweichungen je Stelle ausgeweisen, Ausweis-Ort dieser Notiz) · [Beobachtungs-Evidence](../observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-cr-devcontainer-katalog-quellen.md) (F-1/F-2, V-1).

**Abgleich je Stelle §§2:** (a) wortgleich (einziger Rest „PostgreSQL Client" = Historie-Chronik) · (b) Satz 1 zeichenweise; **Abweichung ausgewiesen:** Satz 2 (Tarball-Offenpunkt) ist im Vertrags-Bullet als Ende „Registry-Tarball-Form offen" verdichtet und als Grenze der Historie-Zeile getragen — inhaltlich gedeckt, Umformulierung je Verifier-Nachzug · (c) Katalog-Aufzählung, Metadaten-Satz, Form-Tabelle samt Generator-Abweisung (je Messbasis) und Pipeline-Schritt-6-Prüfung — der Beleg-Auftrag „Test-Namen mit dem Umsetzungs-Diff" trägt als Klammer der Spec-Historie-Zeile (Verifier V-3, INFO: Form-Abweichung zum [`SPEC-054`](../../../../spec/spezifikation.md)-Präzedenz, kein auflösender Träger-Fehl), die Umsetzungs-Beleg-Zeilen entstehen mit dem Diff · (d) Umformulierung (Chronik-Zellen), Inhalt gedeckt; Reihenfolge 0.32.0/0.33.0 getragen.

**Risiken aus §6 — je Ausgang:**
1. *Wortlaut nicht wortgleich übernommen* — **eingetreten (abgemildert):** die zwei Formulierungs-Abweichungen der Stellen (b)/(d) sind je ausgeweisen (oben); die Träger sind die Vertrags-Zeilen selbst (Verifier-Nachzüge).
2. *Eine Zusage geht weiter als die Messung* — **entfallen:** die Form-Tabelle trägt nur das Gemessene (Registry-Referenz verarbeitet; `oci://`/`https://` unverarbeitbar), die Tarball-Form steht ausdrücklich als Offenpunkt in der Historie und der Grenze ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — kein Fund wider die Messbasis).
3. *Eine Abweisung ohne Beleg-Auftrag* — **entfallen für den CR:** die Aufträge tragen die Form-Tabelle von [`SPEC-023`](../../../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen) und die Klammer der Historie-Zeile; die Code- und Test-Träger liefert der Umsetzungs-Diff ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md): kein neuer Fund).
4. *Plan-Aussagen über Code ohne Fundstelle* — **eingetreten (abgemildert):** die zustandsmarker-Familie und das Historie-Fragment traten (Review F-1/F-2, Verifier V-1) und sind je berichtigt; der §3.10-Abgleich je Fix-und-Closure-Runde läuft erneut ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), Evidence-Datei dieses Vorgangs).

**Was ging anders als geplant:** (a) die **WIP-Flohbewegung** (quellen-Slice zurück nach `open/`, CR-Schnitt voraus) brauchte die Link-Tiefe-Fixes der beiden Träger-Dateien (drei Berichtigungs-Commits `b122f697`, `283f815b` samt `open/README`-Verweis). (b) Der Vertrags-Commit `151a7976` führte die gebrochene Edit-Kante (Historie-Fragment) — vom Reviewer nicht gesehen (d-check liest Anker), vom Verifier per blame gefunden; die Klasse „Zellen-Inhalt keinem Gate zugänglich" ist ein Grenzvermerk, kein neuer Träger. (c) Die Messbasis trug schon die Entscheidungen, aber die Roadmap-Wahrheit kam erst mit der Fix-Runde (F-1).

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre 1:* die zustandsmarker-Familie am WIP-Wechsel (dritter Kontext in Folge: Release-Cuts, CR-Übernahme, WIP-Floh) — die Frage des Teilsensors bleibt beim Projektinhaber. *Lehre 2:* Historie-Zellen-Inhalte sind keinem Gate zugänglich — der Selbstcheck läuft nach jeder gebrochenen Edit-Kante (die Umsetzungs-Runde schuf das Fragment; die Berichtigungs-Runden fingen es). Kein Feld `liegt in`. *Benannte Spec-Lücke:* keine.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) ([`evidence/slice-cr-devcontainer-katalog-quellen.md`](../observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-cr-devcontainer-katalog-quellen.md): F-1, F-2, V-1) · [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — kein Fund als Zusage-Wider · [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) — kein neuer Fund.

**Carveouts:** keiner — kein Gate, keine Schwelle und kein Code im Range. **Folge-Slice:** die Umsetzung [`slice-v1-katalog-features-quellen-und-freigabe`](../in-progress/slice-v1-katalog-features-quellen-und-freigabe.md) nimmt den WIP-Platz der Closure; sie trägt den Katalog-Entfernungs-Diff samt Tabellentest, die Generator- und `config set`-Abweisung der unverarbeitbaren Formen (Klassifikatoren, Test-Namen), den Doku-Nachzug über Handbuch, READMEs und `CHANGELOG.md`. **Weiter offen:** der Registry-Tarball-Weg (Grenze, Historie-Zeile).

**Drei Paarungen:** *Anker* — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert. *Folge-Slice* — der Umsetzungs-Slice ist eine Datei in `open/` (Verzeichnis-Lesung; er übernimmt die Vertrags-Abweisungs-Sätze samt Beleg-Aufträge). *Register* — jede `BEO-*`-Kennung löst in ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls`).

**Verifikation** (Läufe nach Closure-Inhalt, Move und Marker; Docker-Modus): `make gates` Exit 0 am 2026-11-03 (`record-gates`-Stempel = HEAD nach dem Nachzug). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` (DOK), `docs/plan/` (DOK); alle Greenfield. **Vorgelagert — offene Beobachtungen sichten:** die vier [`BEO-DOK/*`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)-Einträge der Risiko-Zeilen oben. **Modus-Begründungsblock:** alle berührten Sub-Areas Greenfield.


**Belege (Closure 2026-11-03):** (1) [`cr-vorschlag-katalog-quellen.md`](cr-vorschlag-katalog-quellen.md) trägt die vier Stellen samt der Mess-Basis M1/M2. (2) Freigabe „wie vorgelegt" 2026-11-03 — Spuren konsistent (Vorschlags-Status, Roadmap-CR-Zeile, Übernahme-Commit, Plan; Verifier-Zeile 2). (3) Übernahme `151a7976` je Stelle — Abgleich: Stelle (a) wortgleich (Rest-Suche „PostgreSQL Client" nur Chronik); Stelle (b) Satz 1 zeichenweise (difflib des Verifiers), **Satz 2 (Tarball-Offenpunkt) im Vertrags-Bullet verdichtet** („Registry-Tarball-Form offen" — Träger je Verifier nachgelesen), Ausgewisener; Stelle (d) Umformulierung mit gedecktem Inhalt. (4) [`make gates`] grün je Runde (Stempel je HEAD; `gates-exit:0` am Fix-Stand), [Review-Report](../../../reviews/2026-11-03-slice-cr-devcontainer-katalog-quellen.md) (0 HIGH/4 MEDIUM/6 LOW, je berichtigt) / [Verify-Report](../../../reviews/2026-11-03-slice-cr-devcontainer-katalog-quellen-verify.md) (V-1 Fragment berichtigt); Register und Risiko-Ausgänge und Paarungen in §7.
