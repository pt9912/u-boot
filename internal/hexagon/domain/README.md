# internal/hexagon/domain

Reine Datentypen und invariantes Verhalten der u-boot-Domäne. Kein I/O,
keine externen Libraries (`LH-FA-ARCH-002`).

Geplante Inhalte (M3+):

- `Project` – Aggregat für ein u-boot-Projekt.
- `Service` – Service-Add-on (PostgreSQL, Keycloak, OTel).
- `ProjectName`, `Port`, `ImageRef` – Value-Objects mit Validierung
  (siehe `LH-FA-INIT-006`).
- `ComposeFile`, `EnvVar` – strukturelle Modelle für die erzeugten
  Artefakte.

Import-Regeln: ausschließlich Go-Standard-Library. Verstöße werden
durch `golangci-lint depguard` im `lint`-Stage abgewiesen
(`LH-FA-ARCH-003`).
