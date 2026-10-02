# Slice V2: Lastenheft verschlanken, Runde 3 — Matrix, Prioritätszeilen, Build- und Lint-Details

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `05212da`**).

**Welle:** `welle-lastenheft-verschlankung` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** [`slice-v2-spec-lastenheft-verschlankung-runde-2`](../done/slice-v2-spec-lastenheft-verschlankung-runde-2.md).
**Berührte Spec-Stellen:** `spec/lastenheft.md` §4.11, §4.13, §8, §13; `spec/spezifikation.md` §1; `.d-check.yml` (Kennungs-Muster).
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Der Projektinhaber: „Traceability-Matrix streichen — das Lastenheft ist immer noch viel zu groß“, und der Hinweis, dass Regelwerk und `MR-*` zählen (Prozess- und Docker-Harness-Regeln stehen dort).

## Ziel und Abgrenzung

**Ziel:** Traceability-Matrix (§13) streichen, Prioritätszeilen in den ersten Satz verdichten, Build-, Layout- und Lint-Anforderungen auf die Zusage kürzen (Detailzeilen wörtlich als Verfeinerungen in die Spezifikation; Multi-Stage-, Docker-only- und Bootstrap-Disziplin führt das Baseline-Regelwerk).

**Ausdrücklich NICHT in diesem Slice:**

- **Zusammenfassen oder Entfernen von Anforderungen** — würde Kennungen und Anker brechen; braucht eine eigene Entscheidung.
- **Änderung von Gates oder Schwellwerten** — unverändert; nur die veraltete Angabe „Default 0“ im Coverage-Bootstrap entfällt.

## Definition of Done

- [x] Matrix gestrichen (Kennungs-Muster in `.d-check.yml` und Verweise in Konventionen angepasst), Prioritätszeilen verdichtet, Detailzeilen verschoben.
- [x] `make gates` grün.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | §13 auf einen Satz, `PH-*`/`TC-*`-Muster und Verweise entfernen. |
| T2 | Prioritätszeilen in den ersten Satz der Anforderung ziehen (Überschriften und Anker unverändert). |
| T3 | Build-, Layout-, Enforcement- und Lint-Details als Verfeinerungen; Coverage-Bootstrap ohne veralteten Default. |

## Risiken

- Wegfall der Matrix schwächt eine Prüfung ab: Es gab keine Verifikationsdaten in ihr (alle `PH-`/`TC-`-Spalten leer); `make doc-trace` leitet aus den Kennungen ab.

## Closure

- **Geliefert (`05212da`):** Traceability-Matrix (§13, 164 Zeilen) gestrichen, Abschnitt auf einen Satz; `PH-*`/`TC-*`-Muster in `.d-check.yml` und Verweise in Konventionen und Gate-Beschreibungen entfernt; Prioritätszeilen aller Anforderungen in den ersten Satz gezogen (Überschriften und Anker unverändert, 240 Zeilen); Build-, Layout-, Enforcement- und Lint-Detailzeilen wörtlich als Verfeinerungen in die Spezifikation, im Lastenheft nur noch die Zusagen; die veraltete Angabe „Default-Schwellwert 0“ im Coverage-Bootstrap entfällt (Folgepunkt aus dem Review der ersten Welle erledigt).
- **Größe:** Lastenheft 2984 (vor der Welle) → 2039 → **1626** Zeilen (−45 %).
- **Gegenlesen:** Das Skript meldet erwartungsgemäß die verdichteten Prioritätszeilen und die gestrichene Matrix als „nicht in der Spezifikation“; inhaltliche Zeilen der Detailverschiebung sind vollständig übernommen.
- **Sensoren:** `make gates` grün. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review:** Eigener Review-Lauf für diese dritte Runde steht aus (kleine, mechanische Änderung; Runden 1 und 2 sind reviewt).
