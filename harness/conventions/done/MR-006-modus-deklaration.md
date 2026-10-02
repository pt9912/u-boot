# MR-006 — Modus-Deklaration pro Sub-Area

- **Datum:** 2026-07-24
- **Geltungsbereich:** Abschnitt Modus-Deklaration in [`conventions.md`](../../conventions.md).
- **Ersetzt-Baseline-Regel:** [`grundlagen-bootstrap.md` §Modus pro Sub-Area: Greenfield vs Brownfield](../../../.harness/baseline/v6.13.0/regelwerk/grundlagen-bootstrap.md#modus-pro-sub-area-greenfield-vs-brownfield)
- **Adaption:** u-boot traegt Bestandscode (`hexagon/`, `cmd/`, `internal/`) neben
  den Doku-Sub-Areas. Der Abschnitt Modus-Deklaration unten ordnet jede Sub-Area
  als GF/BF/Hybrid ein; jede BF-/Hybrid-Markierung traegt eine
  Graduation-Bedingung. **Audit ausgefuehrt (2026-07-25):** Der Erst-Pass
  (drei grobe Sub-Areas, pauschal GF) ist durch eine auditierte Einordnung
  ersetzt - Drei-Achsen-Inklusion je Kandidat, dann vier Modus-Kriterien je
  qualifizierter Sub-Area. Ergebnis: acht Sub-Areas statt drei, davon eine
  Hybrid (`internal/adapter/driving/cli`) und eine Brownfield
  (`internal/**/README.md`); `cmd/uboot` faellt auf Sub-Area-Aspirantin zurueck.
  **Graduation im selben Durchlauf:** Die Hybrid-Aussage zu
  `internal/adapter/driving/cli` hat ihre Bedingung noch in dieser Welle
  erfuellt (beide Adapter-Konventionen stehen in der Sicht-Spec) und ist auf GF
  gesetzt - eine Graduation, die stattgefunden hat, bleibt nicht als offene
  Ausnahme stehen.
- **Begruendung:** Das Regelwerk verlangt eine Modus-Aussage pro qualifizierter
  Sub-Area; eine BF-Sub-Area ohne Graduation-Plan waere "permanente Ausnahme als
  temporaer getarnt".
- **Aufloesungs-Trigger:** Audit erledigt; Delivery-Verweis im Slice
  [`slice-harness-sub-area-modus-audit`](../../../docs/plan/planning/done/slice-harness-sub-area-modus-audit.md)
  §9. Re-evaluieren bei jeder neuen Pfad-Familie im Produktivcode sowie beim
  Erfuellen einer der beiden Graduation-Bedingungen oben.
