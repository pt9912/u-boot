# Slice V2: Lastenheft verschlanken, Runde 2 — Beispiele, Restverhalten, Regelwerk-Dopplungen

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `d024a1b`**).

**Welle:** `welle-lastenheft-verschlankung` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** [`slice-v2-spec-lastenheft-verschlankung`](../done/slice-v2-spec-lastenheft-verschlankung.md) (Runde 1), Fundament [`slice-v2-spec-technik-stratum-grundlage`](../done/slice-v2-spec-technik-stratum-grundlage.md) §Schnittregeln.
**Berührte Spec-Stellen:** `spec/lastenheft.md` §4, `spec/spezifikation.md` §1 und §2.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Der Projektinhaber meldet nach Runde 1: Das Lastenheft ist „immer noch sehr groß“, die `bash`-Beispiele sollen aus den Anforderungen, es stehe noch viel Technisches darin, und das Dokumentationsreferenzmodell stehe bereits im Regelwerk.

## Ziel und Abgrenzung

**Ziel:** (1) Aufrufbeispiele aus den Anforderungen entfernen, (2) restliche Mechanismus-Absätze als Verfeinerungen in die Spezifikation, (3) Anforderungen, die Prozessregeln des adoptierten Baseline-Regelwerks wiederholen, auf eine Vertragszusage kürzen (die Kopien in der Spezifikation entfallen, weil das Regelwerk sie führt).

**Ausdrücklich NICHT in diesem Slice:**

- **Entfernen der Traceability-Matrix (§13), der Akzeptanzkriterien (§9) oder von Kennungen** — braucht eine Entscheidung des Projektinhabers (berührt `.d-check.yml`, `harness/conventions.md` und Verweise).
- **Kürzen der projektspezifischen Build- und Architektur-Anforderungen** — sie sind repo-eigene Zusagen, keine Regelwerk-Wiederholung.

## Definition of Done

- [x] `bash`-Beispiele aus §4–§7 entfernt (gesammelt in der Spezifikation), Restverhalten und Regelwerk-Dopplungen gekürzt; Gegenlesen ohne Verlust (Ausnahme: bewusst entfernte Regelwerk-Dopplungen).
- [x] `make gates` grün.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | Beispiele: 31 Blöcke gesammelt in die Spezifikation, Sätze mit Block-Bezug inline umformuliert. |
| T2 | Restverhalten: sechs weitere Mechanismus-Absätze als Verfeinerungen. |
| T3 | Regelwerk-Dopplungen: ADR-Format, Lifecycle, Carveout-Disziplin und Referenzmodell auf Kurzzusagen; Kopien in der Spezifikation entfernt. |

## Risiken

- Eine Kürzung schwächt eine Zusage ab: je Anforderung bleiben Exit-Code, Sicherheits- und Opt-in-Verhalten im Vertrag benannt.

## Closure

- **Geliefert (`d024a1b`):** 31 `bash`-Beispielblöcke aus §4–§7 in einem Eintrag der Spezifikation gesammelt (acht Sätze mit Block-Bezug inline umformuliert); sechs weitere Mechanismus-Absätze als Verfeinerungen übernommen (Zuordnung von Fehlern zu Exit-Codes, Konvention für Diagnose-Codes, relevante Dateien der Projekterkennung, Backup-Verfahren, Ausstattung bei `nestedRuntime: podman`, Verhalten bei nicht gewährbarer Capability); die Anforderungen zu ADR-Format, Planning-Lifecycle, Carveout-Disziplin und Dokumentationsreferenzmodell sind auf Kurzzusagen gekürzt, ihre Kopien in der Spezifikation (Format der ADRs, zwei Verfeinerungen der Doku-Anforderungen) entfallen, weil das adoptierte Baseline-Regelwerk sie führt.
- **Größe:** Lastenheft 2290 → **2039** Zeilen (vor der Welle 2984); Spezifikation 819.
- **Gegenlesen:** Skript-Abgleich: nicht wiederfindbar sind nur die Fence-Zeilen der Beispielblöcke (Inhalt steht gesammelt in der Spezifikation) und die bewusst entfernten Regelwerk-Dopplungen.
- **Sensoren:** `make gates` grün. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review (nachgeholt, unabhängige Rolle, beide Runden):** `docs/reviews/2026-10-02-lastenheft-verschlankung-r2.md` — 0 HIGH, 1 MEDIUM, 7 LOW, 5 INFO. Behoben: abgeschwächte Exit-Code-Zusage bei blockierten User-Namespaces (MEDIUM), Empfehlungs-Charakter der Exit-Code-Zuordnung, Status-/Diagnose-Regeln, inhaltstragende Aufrufangaben, Leerzeilen-Artefakte, fehlende Zusagen in den gekürzten Doku-Anforderungen, Historie-Zeile, Spezifikations-Verweis. Bewusst offen: Verweis in der akzeptierten ADR zur Degradationstabelle (ADR unveränderlich; die Tabelle steht in der Verfeinerung der Anforderung).
- **Befund:** Die verbleibenden rund 2000 Zeilen bestehen zu über der Hälfte aus Struktur (144 Überschriften mit Priorität, Leerzeilen, Traceability-Matrix 164 Zeilen, Akzeptanzkriterien 132); echte Technik ist in Anforderungstexten kaum noch übrig.
