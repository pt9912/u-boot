# internal/hexagon/domain

Reine Datentypen und invariantes Verhalten der u-boot-Domäne. Kein I/O,
keine externen Libraries (`LH-FA-ARCH-002`).

## Aktueller Inhalt (M3-T1)

- `ProjectName` — validierter Value-Object-Typ; Regex aus
  `LH-FA-INIT-006` (`^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`). Konstruktor
  `NewProjectName` liefert `ErrInvalidProjectName`-gewrappte Fehler.
- `NormalizeProjectName` — deterministische 6-Schritt-Normalisierung
  nach `LH-FA-INIT-002` (lowercase, map-to-dash, collapse-dashes,
  trim, clamp-63, re-trim).
- `Project` — Aggregat mit `Name ProjectName` und
  `SchemaVersion int`; `SchemaVersionCurrent = 1` (`LH-DA-003`,
  `LH-FA-CONF-002`).

## Geplante Erweiterungen

- `Service` — Service-Add-on (PostgreSQL, Keycloak, OTel) — mit
  M4/M5 (`LH-FA-ADD-*`).
- `Port`, `ImageRef` — Value-Objects für Service-Konfiguration.
- `ComposeFile`, `EnvVar` — strukturelle Modelle für erzeugte
  Artefakte.

## Import-Regeln

Ausschließlich Go-Standard-Library. Verstöße werden durch
`golangci-lint depguard` im `lint`-Stage abgewiesen (`LH-FA-ARCH-003`).
