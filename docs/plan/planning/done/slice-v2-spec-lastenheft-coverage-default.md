# Slice V2: Coverage-Default im Lastenheft an die gelebte Praxis angleichen

> **Status:** **Done** (2026-10-02, erledigt durch die Kürzung von [`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap); Gegenstand entfallen) — Befund aus dem Review der Welle `welle-spec-technik-stratum` (F-4). ([`docs/reviews/2026-10-02-welle-spec-technik-stratum.md`](../../../reviews/2026-10-02-welle-spec-technik-stratum.md), F-4). Carveout-Plan-Anker.

**Welle:** ohne Welle.
**Bezug:** [`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap), [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik).
**Berührte Spec-Stellen:** `spec/lastenheft.md` (Bootstrap-Aussage zum Default-Schwellwert), `spec/spezifikation.md` (Build-Args).
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

[`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap) nennt als Bootstrap-Verhalten den Default-Schwellwert `0` („wird in einem Folge-Schritt angehoben“). Gelebt und im Build gesetzt ist `90` Prozent (`Dockerfile` `ARG COVERAGE_THRESHOLD=90`, `Makefile` `THRESHOLD ?= 90`, [`AGENTS.md`](../../../../AGENTS.md) §Quality Gates). Die Spezifikation führt seit der Welle den gelebten Wert (`90`); das Lastenheft trägt noch die Bootstrap-Aussage. Das Lastenheft ist ein Vertrag: Die Angleichung ist eine Vertragsänderung und gehört nicht in einen Struktur-Schnitt.

## Aufhebungsbedingung

Das Lastenheft nennt keinen vom Build abweichenden Default mehr (Bootstrap-Aussage als erledigt gekennzeichnet oder an den gelebten Wert gebunden); Version-Bump und Historie-Zeile nach `LH-FA-PROJDOCS`-Vertragsregeln.

## Akzeptanzkriterien

- Die Aussage zum Default-Schwellwert in [`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap) stimmt mit Dockerfile und Makefile überein; das Verhalten bei leerer Coverage-Eingabe (kein falsches Grün) bleibt erhalten.
- Version und Historie des Lastenhefts sind fortgeschrieben; `make gates` grün.

## Out of Scope

- Änderung des Schwellwerts selbst (bleibt 90 Prozent); Gates werden nicht gelockert.

## Closure

- **Erledigt:** Beim Verschlanken des Lastenhefts (Welle `welle-lastenheft-verschlankung`) wurde die Bootstrap-Angabe „Default-Schwellwert `0`“ aus der Anforderung entfernt; der Schwellwert ist überschreibbar, der gelebte Wert (`90`) steht in der Spezifikation. Die Änderung ist in der Historie-Zeile 0.4.0 des Lastenhefts festgehalten. Gates unverändert.
