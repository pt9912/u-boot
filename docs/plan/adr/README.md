# docs/plan/adr

Architecture Decision Records (ADRs) für u-boot.

Format und Konventionen sind in `LH-FA-PROJDOCS-002`
([../../../spec/lastenheft.md](../../../spec/lastenheft.md)) verbindlich
festgelegt:

- Dateiname: `<NNNN>-<kebab-slug>.md`, beginnend bei `0001`, monoton steigend.
- Mindestabschnitte in dieser Reihenfolge:
  1. Dokumenttitel als `#`: `# ADR <Nr>: <Titel>`
  2. `## Status` – `Proposed` | `Accepted` | `Superseded by <NNNN>-<slug>` | `Deprecated`
  3. `## Datum` – `YYYY-MM-DD`
  4. `## Kontext`
  5. `## Entscheidung`
  6. `## Konsequenzen`
- ADR-Nummern werden nie wiederverwendet. Abgelöste ADRs bleiben mit
  Status `Superseded by <NNNN>-<slug>` erhalten.
