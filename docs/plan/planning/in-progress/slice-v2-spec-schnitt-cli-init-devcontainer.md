# Slice V2: Lastenheft-Schnitt 1: CLI-Grundverhalten, Projektinitialisierung, Devcontainer

> **Status:** **in Arbeit** (seit 2026-10-02) — Fundament geliefert.

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** Lastenheft §4.1 CLI-Grundverhalten, §4.2 Projektinitialisierung, §4.3 Devcontainer-Unterstützung, §6.1 Kommandozeilenschnittstelle (Scope).
**Berührte Spec-Stellen:** `spec/spezifikation.md` §1–§4, §6; `spec/lastenheft.md` §4.1–§4.3, §6.1.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Größte Technik-Dichte im Lastenheft: JSON-Schema der `--dry-run`/`--diff`-Ausgabe, Diagnose-Codes, Ablaufbeschreibungen der
Init-Zustandsmaschine, Devcontainer-Defaults. Dieser Slice ist der **Pilot** des Schnitts: Er legt Form und Granularität
der Verfeinerungen fest, an der die übrigen Schnitt-Slices sich orientieren.

## Ziel und Abgrenzung

**Ziel:** Die Kapitel §4.1–§4.3 und §6.1 enthalten nur noch Vertrag (Was, Akzeptanzkriterien, Out-of-Scope); jede technische Festlegung steht wörtlich in der Spezifikation mit Kennung.

**Schnittregeln:** siehe [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln (bleibt im Lastenheft / zieht in die Spezifikation, Decken-Regel, Verbleib-Tabelle, Zeilenverweise).

**Ausdrücklich NICHT in diesem Slice:**

- **Kapitel §4.4 und folgende** — eigene Schnitt-Slices; ein Pilot, der alles anfasst, ließe sich nicht mehr gegenlesen.
- **Neuformulierung von Anforderungen** — der Schnitt entscheidet nichts neu; Unklares bleibt im Lastenheft (siehe Offen).


## Definition of Done

- [ ] §4.1–§4.3, §6.1 geschnitten; Verbleib-Tabelle vollständig (jede Quelle genau ein Ziel); `LH-*`-Überschriften und Anker unverändert.
- [ ] Spezifikation um die übernommenen Stellen ergänzt (Verfeinerungen `.a`, `SPEC-<NNN>`), ohne ADR- und Slice-Verweis; Zeilenverweise auf die berührten Stellen umgestellt.
- [ ] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [ ] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Inventur.** Je Anforderung in §4.1–§4.3 und §6.1 klassifizieren (bleibt / zieht um), Verbleib-Tabelle im Slice anlegen. |
| T2 | **Übernahme CLI.** §4.1 und §6.1: Schemata, Codes, Beispielinstanzen in die Spezifikation (§2, §4). |
| T3 | **Übernahme Init und Devcontainer.** §4.2 und §4.3: Zustandsfolgen, Defaults, Dateiformate. |
| T4 | **Gegenlesen.** Verbleib-Tabelle gegen Diff prüfen, `make gates`, Review. |

## Risiken

- Anforderungstext und Technik sind in einem Absatz verwoben (zum Beispiel Akzeptanzkriterium mit eingebettetem Schema): Schnitt an der Satzgrenze, Rest bleibt im Lastenheft; Zweifel unter „Offen“ statt raten.
- Link-Ziele auf verschobene Beispiele (Code-Kommentare, Tests) brechen: T4 sucht per `grep` nach `lastenheft.md:`-Zeilenverweisen.

## Closure

Offen. Beim Schließen: Verifikation-Evidence nach `harness/verification.md` (DoD, Sensoren, nicht ausgeführte
Sensoren, Carveouts), Steering-Loop-Lerneintrag, Delivery-Hash im Kopf.
