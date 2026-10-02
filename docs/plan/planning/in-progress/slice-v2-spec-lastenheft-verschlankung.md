# Slice V2: Lastenheft verschlanken — Verhaltensdetails in die Spezifikation

> **Status:** **in Arbeit** (seit 2026-10-02).

**Welle:** `welle-lastenheft-verschlankung` (siehe [`roadmap.md`](roadmap.md)).
**Bezug:** alle `LH-FA-*`-Anforderungen mit langen Regelblöcken (Scope siehe unten); Fundament: [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln.
**Berührte Spec-Stellen:** `spec/lastenheft.md` §4, §5, `spec/spezifikation.md` §1.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Nach dem ersten Schnitt ([`welle-spec-technik-stratum`](roadmap.md)) misst das Lastenheft 2636 Zeilen (vorher 2984), davon 1483 in §4. Der Projektinhaber bewertet das als „immer noch sehr groß“. Die erste Runde hat bewusst nur Formate, Schemata und Werte verschoben und beobachtbares Verhalten im Vertrag gelassen. Die Masse liegt jedoch in **ausführlichen Regelblöcken** (Entscheidungslogik der Bestätigungsmodi, Add-on-Zustandsregeln, Degradationsdetails, Prozessregeln der Doku-Anforderungen), die ein Pflichtenheft-Niveau haben.

## Ziel und Abgrenzung

**Ziel:** Jede Anforderung nennt im Lastenheft die Zusage in wenigen Sätzen; der ausführliche Regelblock (Entscheidungstabellen, Zustandsregeln, Randfälle, Prozessdetails) steht wörtlich als Verfeinerung `<Anforderung>.<Buchstabe>` in der Spezifikation. Das Lastenheft verliert dabei keine Zusage: Exit-Codes, Sicherheitsverhalten, Gates und Opt-in-Pflichten bleiben im Vertrag benannt.

**Ausdrücklich NICHT in diesem Slice:**

- **Änderung des Inhalts** — Übernahme wörtlich; Kürzungen im Lastenheft fassen nur zusammen.
- **Entfernen oder Umbenennen von `LH-*`-Kennungen** — Anker in ADRs, Slices, Code und Doku bleiben stabil.
- **Traceability-Matrix (§13), Akzeptanzkriterien (§9), Glossar** — reine Indizes bzw. Abnahmetexte; bleiben unverändert.

## Definition of Done

- [ ] Regelblöcke der Anforderungen (siehe Verbleib-Tabelle) in die Spezifikation verschoben, Lastenheft-Zusagen gekürzt; Verbleib-Tabelle vollständig, Gegenlesen ohne Verlust.
- [ ] `make gates` grün; Review durch eine andere Rolle (Report unter `docs/reviews/`).

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur:** Regelblöcke > 8 Zeilen je Anforderung bestimmen, Zusage-Satz je Anforderung festlegen. |
| T2 | **Übernahme:** Blöcke wörtlich als Verfeinerungen übernehmen, Lastenheft kürzen. |
| T3 | **Gegenlesen und Review:** Verlustprüfung, Lesbarkeit, Review-Report. |

## Risiken

- Eine Kürzung schwächt eine Zusage ab: je Anforderung bleibt Exit-Code, Sicherheits- und Opt-in-Verhalten im Vertrag benannt; das Review prüft das ausdrücklich.

## Closure

Offen.
