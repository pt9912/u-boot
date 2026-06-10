# ADR 0013: Dokumentationsreferenzmodell und normative Kanten

## Status

Accepted

## Datum

2026-06-10

## Kontext

[`LH-FA-PROJDOCS-001`](../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..[`LH-FA-PROJDOCS-005`](../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin) definieren Struktur, ADR-Format,
Planning-Lifecycle und Carveout-Disziplin. Bisher war aber nicht
verbindlich geregelt, welche Dokument-Referenzen normative Kraft tragen
und welche nur Planungskontext oder Buchfuehrung sind.

Das Stable Dependencies Principle ist als Denkmodell nuetzlich, reicht
aber fuer Dokumente nicht aus: Stabilitaet und Autoritaet fallen
auseinander, sobald ein eingefrorenes Artefakt seine normative Kraft
verliert. Eine superseded ADR ist stabil, aber nicht mehr autoritativ.
Das Referenzmodell muss deshalb Normativitaet explizit regeln.

[`LH-FA-PROJDOCS-006`](../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) legt fest, dass das Lastenheft die normative Decke
des Projektmodells ist. Externe Normen, Gesetze, Standards,
Upstream-Vertraege oder Produktvorgaben wirken im Repo nur ueber
explizite `LH-*`-Anforderungen normativ.
`spec/architecture.md` ist ein Sicht-Stratum: Es visualisiert und
praezisiert den durch das Lastenheft begrenzten Architekturstand, ist
aber nicht die Wurzel des Modells.

## Entscheidung

u-boot verwendet folgendes Dokumentationsreferenzmodell:

> Normative Kraft existiert nur auf aufwaertsgerichteten
> Inter-Layer-Kanten plus der ADR-internen Lineage-Kante. Alles andere
> ist Kontext.

Die Referenzmatrix:

| Dokument ↓ referenziert → | Lastenheft | ADR | Slice | Carveout | Roadmap/Welle |
| --- | --- | --- | --- | --- | --- |
| **Lastenheft** | **Normativ**: nur intra-`LH-*` | ❌ | ❌ | ❌ | ❌ |
| **ADR** | **Normativ**: `LH-*`-Grundlage | **Normativ/Lineage**: aktive ADRs als Grundlage; superseded ADRs nur ADR-interne Historie | ❌ | ❌ | ❌ |
| **Slice** | **Normativ**: `LH-*`-Scope | **Normativ**: nur aktive ADRs | **Kontext**: `triggered-by`, `blocked-by`, `follow-up-of` | **Kontext**: eigener/offener Carveout, Debt-/Closure-Rueckverweis | **Kontext**: Einordnung in Welle/Roadmap |
| **Carveout** | **Normativ**: betroffene `LH-*` | **Normativ**: betroffene aktive ADRs | **Kontext/Traceability**: owner/verursachender/schliessender Slice | **Kontext**: ersetzt/zusammengefuehrt/abhaengig von anderem Carveout | **Kontext**: Welle/Planungseinordnung |
| **Roadmap/Welle** | **Kontext**: Zielbild/Scope, keine Spezifikation | **Kontext**: Architekturhintergrund, keine Entscheidungskraft | **Kontext**: Orchestrierung/Sequenzierung | **Kontext**: Risiko-/Debt-Uebersicht | **Kontext**: Hierarchie/Sequenz |

Zusaetzliche Regeln:

1. **Autoritaet schlaegt Stabilitaet.** Slices und Carveouts
   referenzieren nur aktive ADRs. Eine superseded ADR darf nur innerhalb
   der ADR-Lineage referenziert werden.
2. **Spec-Straten referenzieren nicht abwaerts.** Technische oder
   Sicht-Specs duerfen auf das Lastenheft und innerhalb ihres Stratums
   referenzieren, aber keine ADRs, Slices, Carveouts oder Roadmap/Wellen
   als bindenden Text verlinken. Die ADR deklariert aufwaerts, welche
   `LH-*`- oder Spec-Stellen sie begruendet oder schaerft.
3. **ADR-Lineage ist normativ, Carveout-Lineage nicht.** `ADR -> ADR`
   kann normative Entscheidungslinie tragen. `Carveout -> Carveout`
   bleibt reine Buchfuehrung, weil Carveouts Schuld und Ablaufzustand
   dokumentieren.
4. **Slice- und Roadmap-Kanten sind Kontext.** Sie duerfen Trigger,
   Reihenfolge, Blocker, Owner oder Closure dokumentieren, erzeugen aber
   keine Spezifikation.
5. **Carveout-Rueckkanten sind Traceability.** `Slice -> Carveout` und
   `Carveout -> Slice` sind erlaubt, aber ausschliesslich als Rueckverweis
   auf bekannte Schuld, offene Ausnahme, Owner oder Closure-Arbeit.
6. **Keine offene Normativitaet nach aussen.** Externe Autoritaet wird
   in `LH-*` importiert, versioniert und pruefbar gemacht. Direkte
   externe Referenzen koennen Herkunft oder Begruendung liefern, aber
   keine offene Normquelle im Repo sein.

## Konsequenzen

- Lastenheft-Passagen duerfen keine ADRs, Slices, Carveouts oder
  Roadmap/Wellen als Quelle ihrer Normativitaet referenzieren.
- `spec/architecture.md` darf keine ADRs als bindenden Text
  referenzieren. Aenderungskopplung lebt in der ADR: die ADR nennt die
  betroffenen `LH-*`- oder Spec-Stellen, die nachzuziehen sind.
- ADRs tragen Entscheidungen und duerfen keine konkreten Slices als
  Umsetzung, Trigger oder Closure-Quelle referenzieren. Umsetzungsspur,
  Trigger und Closure gehoeren in Slice-, Roadmap- und
  Carveout-Artefakte.
- Slices duerfen `LH-*` und aktive ADRs normativ referenzieren. Peer-
  Slices, Roadmap/Wellen und Carveouts bleiben Kontext.
- Carveouts muessen ihre fachliche Begruendung ueber `LH-*` oder aktive
  ADRs tragen; Slice-Bezuege bleiben Owner-/Trigger-/Closure-
  Buchfuehrung.
- `docs-check` prueft Markdown-Link-Pfade, Heading-Anker, nackte
  `ADR-*`-Kennungen, `LH-*`-Kennungen in Spec-Straten, README-Dateien,
  `harness/`, `docs/archive/` und `docs/user/`. Im Lastenheft selbst
  sind `LH-*`-Selbstreferenzen ausserhalb von Ueberschriften
  linkpflichtig, wenn die Ziel-ID einen eigenen Heading-Anker hat.
  Konkrete
  Slice-/Tranche-IDs duerfen in Lastenheft, Sicht-Specs und ADRs nicht
  auftauchen; in allen anderen gescannten Markdown-Dateien sind
  eindeutig aufloesbare Planning-IDs linkpflichtig. Markdown-
  Ueberschriften sind fuer Planning-IDs ausgenommen, damit bestehende
  Section-Anker stabil bleiben. Konkrete `PH-*`- und `TC-*`-Kennungen
  sind im Lastenheft Traceability-Aliase auf die zugehoerige `LH-*`-
  Anforderung derselben Matrixzeile; konkrete `CO-*`-Kennungen sind
  linkpflichtig, sobald sie eingefuehrt werden.
