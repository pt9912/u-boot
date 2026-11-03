# Slice CR: CR-Schnitt Katalogquellen — Entfernung `cpp`/`postgres-client` und Abweisung unverarbeitbarer Formen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`; er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** ohne Welle-Datei (Rest nach dem Release `v0.20.0`); Anlass: die Messung [`docs/messungen/slice-v1-katalog-features-quellen-und-freigabe/`](../../../messungen/slice-v1-katalog-features-quellen-und-freigabe/README.md) und die zwei Entscheidungen des Projektinhabers vom 2026-11-03 (N1 entfernen, N2 Abweisung). Muster [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0032`](../../adr/0032-katalog-features-socket-und-lockerungen.md) (*Nicht Gegenstand* (a)/(b) — die Messung trägt die Entlassung).

**Verantwortlich:** pt9912 (Implementer; Wortlaut-Vorlage und Vertrags-Übernahme im Auftrag des Hauptagenten)

**Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** der freigegebene CR-Wortlaut ([`cr-vorschlag-katalog-quellen.md`](../open/cr-vorschlag-katalog-quellen.md)) — die Entfernung der zwei Katalog-Einträge samt der Vertrags-Beispielzeile und die Abweisung unverarbeitbarer externer Quell-Formen je Exit-Code `10` — wird je Stelle freigegeben und wortgleich übernommen. Kein Produktcode in diesem CR (Muster der CR-Slices).

**Nicht in diesem Slice:** die Umsetzung (Katalog-Entfernung im Code samt Tabellentesten, Generator-Abweisung samt Klassifikatoren, Doku-Nachzug) — Umsetzungs-Slice nach der CR-Closure; die Tarball-Probe (bleibt Offenpunkt mit Historie-Vermerk).

## 2. Definition of Done

- [ ] Vorschlags-Datei mit Alt/Neu je Stelle (a)–(d) samt Mess-Basis
- [ ] Freigabe des Projektinhabers je Stelle (oder mit benannten Änderungen)
- [ ] Vertrags-Übernahme wortgleich ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Beispiele-Zeile + Abweisungssätze; [`SPEC-023`](../../../../spec/spezifikation.md)/[`SPEC-040`](../../../../spec/spezifikation.md); Historie 0.33.0; Versionsfeld), Abgleich je Stelle in §7
- [ ] `make gates` grün; Review + Verification mit Reports (`docs/reviews/`); Nachlauf-Review über jeden Diff von Vertragstext ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)); Closure-Notiz mit Lerneintrag; Register oder „keine Beobachtung"; je Risiko ein Ausgang; drei Paarungen getragen

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `cr-vorschlag-katalog-quellen.md` | neu | DoD-Zeile 1 |
| `spec/lastenheft.md` | update (nach Freigabe) | Stellen (a), (b), (d); Versionsfeld 0.33.0 |
| `spec/spezifikation.md` | update (nach Freigabe) | Stelle (c) — Katalogtabelle, Form-Tabelle, Pipeline-Schritt 6; Historie |
| dieser Slice, [`roadmap.md`](../in-progress/roadmap.md), `open/README.md` (Bestand) | update | Lifecycle je `make slice-mv` |
| `docs/reviews/` (Berichte) | neu | je Berichte-Runden |
| Produktcode, Templates, Tests, Handbuch, READMEs, `CHANGELOG.md` | **kein Diff in diesem CR** | der Umsetzungs-Slice trägt sie |

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

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` (DOK), `docs/plan/` (DOK); alle Greenfield. **Vorgelagert — offene Beobachtungen sichten:** die vier [`BEO-DOK/*`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)-Einträge der Risiko-Zeilen oben. **Modus-Begründungsblock:** alle berührten Sub-Areas Greenfield.
