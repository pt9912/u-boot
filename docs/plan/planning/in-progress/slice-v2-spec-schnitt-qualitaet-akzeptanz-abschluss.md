# Slice V2: Lastenheft-Schnitt 4 und Abschluss: Qualität, Akzeptanz, Traceability, Version 0.4.0

> **Status:** **in Arbeit** (seit 2026-10-02).

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §5 Nichtfunktionale Anforderungen, §8 Qualitätsanforderungen, §9 Akzeptanzkriterien, §10–§16 (Abgrenzung, Risiken, MVP, Traceability, Offene Punkte, Glossar, Historie) (Scope).
**Berührte Spec-Stellen:** `spec/lastenheft.md` §5, §8–§16 (Version, Historie); `spec/spezifikation.md` §1–§7.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Die Qualitätsanforderungen tragen Messmethoden, die teils technische Festlegungen enthalten (Schwellen, Messaufbau). Zum Abschluss
bekommt das Lastenheft den neuen Vertragsstand: Version 0.4.0 mit Historie-Zeile, Lesehinweis auf die drei Straten und
konsistente Traceability-Matrix.

## Ziel und Abgrenzung

**Ziel:** §5 und §8–§16 sind geschnitten; das Lastenheft trägt Version 0.4.0 mit Historie-Zeile; die Spezifikation ist vollständig befüllt und frei von Platzhaltern.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Änderung von Anforderungen oder Schwellenwerten** — der Schnitt verschiebt nur; Coverage-Schwelle, Exit-Codes usw. bleiben wertgleich.
- **Architektur-Sicht** — eigener Slice.


## Definition of Done

- [ ] §5, §8–§16 geschnitten, Verbleib-Tabelle der gesamten Welle vollständig (jede Quelle genau ein Ziel, nichts verwaist); Version 0.4.0, Historie-Zeile, Traceability-Matrix konsistent.
- [ ] Spezifikation vollständig (§1–§7, Historie), `SPEC-<NNN>` fortlaufend, keine Platzhalter; `LH-*`-Anker unverändert.
- [ ] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [ ] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur und Verbleib-Tabelle** für §5, §8–§16. |
| T2 | **Übernahme** der technischen Messaufbauten und Schwellen-Festlegungen; Wertgleichheit der Schwellen gegen die Gates prüfen. |
| T3 | **Vertragsstand:** Version 0.4.0, Historie-Zeile, Lesehinweis auf die drei Straten, Traceability-Matrix. |
| T4 | **Gegenlesen der gesamten Welle,** `make gates`, Review, Welle-Closure-Notiz. |

## Risiken

- Wertdrift bei Schwellen (zum Beispiel Coverage 90 %): T2 vergleicht die Zahlen mit `Makefile` und `.golangci.yml`.

## Closure

Offen. Beim Schließen: Verifikation-Evidence nach `harness/verification.md` (DoD, Sensoren, nicht ausgeführte
Sensoren, Carveouts), Steering-Loop-Lerneintrag, Delivery-Hash im Kopf.
