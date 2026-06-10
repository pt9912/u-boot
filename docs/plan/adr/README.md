# docs/plan/adr

Architecture Decision Records (ADRs) für u-boot.

Format und Konventionen sind in [`LH-FA-PROJDOCS-002`](../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
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

## Index

| ADR | Status | Entscheidung |
| --- | --- | --- |
| [ADR 0001](0001-implementierungssprache-go.md) | Accepted | Implementierungssprache Go |
| [ADR 0002](0002-hexagonale-architektur.md) | Accepted | Hexagonale Architektur mit driving/driven-Split |
| [ADR 0003](0003-solid-nahes-lint-profil.md) | Accepted | SOLID-nahes Lint-Profil |
| [ADR 0004](0004-ci-system.md) | Accepted | CI-System mit GitHub Actions und Docker-only-Gates |
| [ADR 0005](0005-cli-framework-cobra.md) | Accepted | CLI-Framework Cobra |
| [ADR 0006](0006-revive-custom-rules.md) | Accepted | revive Custom-Rules-Profil |
| [ADR 0007](0007-distributionswege-ghcr.md) | Accepted | Distributionswege mit GHCR und Binary, Restwege vertagt/verworfen |
| [ADR 0008](0008-plugin-system-statisch.md) | Accepted | Add-on-System bleibt statisch |
| [ADR 0009](0009-template-format-yaml-files.md) | Accepted | Template-Format YAML-Metadaten plus `text/template` |
| [ADR 0010](0010-kein-http-driving-adapter.md) | Accepted | Kein HTTP-Driving-Adapter |
| [ADR 0011](0011-agent-harness-scaffolding.md) | Proposed | Agent-Harness-Scaffolding |
| [ADR 0012](0012-devcontainer-egress-firewall.md) | Proposed | Devcontainer-Egress-Firewall |
| [ADR 0013](0013-dokumentationsreferenzmodell.md) | Accepted | Dokumentationsreferenzmodell und normative Kanten |
