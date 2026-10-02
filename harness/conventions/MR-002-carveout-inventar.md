# MR-002 — Carveout-Inventar an fester Stelle statt `docs/plan/carveouts/`

- **Datum:** 2026-07-24
- **Geltungsbereich:** Carveout-Ablage;
  [`docs/plan/planning/in-progress/carveouts.md`](../../docs/plan/planning/in-progress/carveouts.md),
  [`AGENTS.md`](../../AGENTS.md) §3.15, `.d-check.yml` `matrix`.
- **Ersetzt-Baseline-Regel:** [`modul-07-carveouts.md` §Ziel-Form: Carveout](../../.harness/baseline/v6.13.0/regelwerk/modul-07-carveouts.md#ziel-form-carveout)
- **Adaption:** Carveouts werden **inventarisiert** in der einen Datei
  `docs/plan/planning/in-progress/carveouts.md` (mit Plan-Anker je Eintrag),
  nicht als je eine Datei `docs/plan/carveouts/CO-<NNN>-<titel>.md`. Das
  `CO-*`-Schema bleibt fuer die ID-Vergabe gueltig.
- **Begruendung:** Etablierte u-boot-Struktur, bereits in `AGENTS.md` und der
  d-check-`matrix` (Klasse `carveout`) verankert; eine zusaetzliche
  Ein-Datei-pro-Carveout-Ebene braechte keinen Mehrwert und erzeugte Drift.
- **Aufloesungs-Trigger:** permanent, solange Carveouts zentral inventarisiert
  werden.
