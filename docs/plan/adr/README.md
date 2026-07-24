# docs/plan/adr

Architecture Decision Records (ADRs) für u-boot.

Format und Konventionen sind in [`LH-FA-PROJDOCS-002`](../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
verbindlich festgelegt (MADR-/Nygard-Form, vendortes Template, Regelwerk-Stand
v3.5.1). Kurzfassung:

- Dateiname: `<NNNN>-<kebab-slug>.md`, beginnend bei `0001`, monoton steigend.
- Titel `# ADR <Nr>: <Titel>`; darunter die Kopf-Felder als **fette
  Inline-Felder**: `**Status:**` (`Proposed` | `Accepted` | `Deprecated` |
  `Superseded by <NNNN>-<slug>`), `**Datum:**`, `**Autor:**`, `**Bezug:**`,
  `**Schärft:**` (Aufwärts-Kopplung Spec/ADR).
- Abschnitte (`##`) in dieser Reihenfolge: Kontext, Entscheidung, Verglichene
  Alternativen, Konsequenzen, Fitness Function (falls zutreffend),
  Re-Evaluierungs-Trigger, Geschichte.
- ADR-Nummern werden nie wiederverwendet; abgelöste ADRs bleiben mit
  `Superseded by <NNNN>-<slug>` erhalten (klickbarer Link).

**Grandfathering.** Die zum v3.5.1-Format-Wechsel bereits `Accepted` ADRs
(`0001`–`0010`, `0013`) bleiben in der leanen Vorform und sind unveränderlich;
das MADR-Format gilt für neu angelegte ADRs sowie für die noch mutable
`Proposed`-ADRs (`0011`, `0012`) beim nächsten inhaltlichen Anfassen.

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
