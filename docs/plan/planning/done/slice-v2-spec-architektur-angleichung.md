# Slice V2: Architektur-Sicht an die Vorlage angleichen

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `5661be3`**).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) (Scope).
**Berührte Spec-Stellen:** `spec/architecture.md` (Kopf, Komponentenkennungen `ARC-<NNN>`), Vorlage `.harness/baseline/v6.13.0/templates/spec/architecture.template.md`.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

`spec/architecture.md` (648 Zeilen) ist die Sicht der Welle. Gegen die Vorlage zu prüfen: Kopf (`Letzte Änderung`, keine Historie),
Komponenten- und Schnittstellenkennungen `ARC-<NNN>`, Sprach- und Meilensteinfreiheit, keine eigenen Anforderungen.

## Ziel und Abgrenzung

**Ziel:** `spec/architecture.md` entspricht der Vorlage; Komponenten und Schnittstellen sind per `ARC-<NNN>` adressierbar; Slices können aufwärts darauf zeigen.

**Ausdrücklich NICHT in diesem Slice:**

- **Importregeln ändern** — nur Form, kein Inhalt; die Regeln bleiben bindend (depguard).
- **Technik aus dem Lastenheft aufnehmen** — Technik gehört in die Spezifikation; die Sicht visualisiert nur.


## Definition of Done

- [x] Vorlagen-Abgleich dokumentiert (Abschnitte, Kopf, Kennungen); Abweichungen entweder behoben oder als `MR-<NNN>` begründet.
- [x] `ARC-<NNN>` für Komponenten und Schnittstellen vergeben; `make docs-check` grün.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Abgleich** gegen die Vorlage, Liste der Abweichungen. |
| T2 | **Angleichung:** Kopf, Kennungen, Gliederung; keine inhaltliche Änderung der Regeln. |
| T3 | **Prüfen:** `make docs-check`, `make lint` (depguard unberührt), Review. |

## Risiken

- Das Anpassen der Kennungen bricht Anker, auf die ADRs und Slices zeigen: T2 prüft Links mit `make docs-check` und passt Ziele an, nicht den Inhalt.

## Closure

- **Geliefert (`5661be3`):** `spec/architecture.md` entspricht der Vorlage im Inhalt: Rolle und erweiterte Hard Rule (keine ADR-Bezüge, keine Historie), `ARC-001`..`ARC-007` für die Komponenten, `ARC-008`..`ARC-012` für die externen Berührungspunkte, Komponentenverweis je Schicht-Abschnitt, Verweis der Import-Regel-Tabelle aufwärts auf die Spezifikation ([`SPEC-013`](../../../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten)). Abweichung (erweiterte Abschnittsfolge, Anker-Stabilität) ist als `MR-010` begründet.
- **Sensoren:** `make gates` grün; `depguard` und Lint-Konfiguration unberührt. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** in das unabhängige Review der Welle eingeschlossen (keine Befunde zur Sicht außer Wortlaut „Folge-Slice“, behoben).
