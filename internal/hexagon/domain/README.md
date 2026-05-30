# internal/hexagon/domain

Reine Datentypen und invariantes Verhalten der u-boot-Domäne. Kein I/O,
keine externen Libraries (`LH-FA-ARCH-002`).

## Aktueller Inhalt

Project-Identität:
- `ProjectName` (M3) — validierter Value-Object-Typ; Regex aus
  `LH-FA-INIT-006`. `NewProjectName` liefert
  `ErrInvalidProjectName`-gewrappte Fehler.
- `NormalizeProjectName` (M3) — deterministische 6-Schritt-
  Normalisierung nach `LH-FA-INIT-002`.
- `Project` (M3) — Aggregat mit `Name ProjectName` und
  `SchemaVersion int`; `SchemaVersionCurrent = 1`.

Services (M5):
- `ServiceName` — validierter Value-Object-Typ analog `ProjectName`.
- `ServiceState` — 6-stelliger State-Machine-Enum (`Unregistered`,
  `Active`, `Deactivated`, `EnabledUnset`, `InconsistentYAML`,
  `InconsistentBlock`) für `LH-FA-ADD-005`.

Diagnostics (M4 + M6):
- `Severity` — 4-stufig: `Ok`, `Info` (M6-T1, für non-judgmental
  Hinweise wie `up.fire-and-forget`), `Warn`, `Error` (strict-
  monotone Ordering für `MaxSeverity`-Reduce).
- `Diagnostic`, `DiagnosticReport` mit Helper-Methoden
  (`MaxSeverity`, `HasErrors`, `HasWarnings`,
  `SortedByIssuesFirst`).

Service-Lifecycle (M6):
- `ContainerState` — Compose-State-Normalisierung mit
  Dead-Allowlist (`StateDead` ← exited/dead/removing/removed) und
  fail-safe `StateUnknown`-Default. `ParseContainerState` ist
  case-insensitive + whitespace-trimmend.
- `StabilizationOutcome` — `OutcomeRunningOnly` (Zero-Value, fail-
  safe), `OutcomeStabilized`, `OutcomeFailed`.
- `RestartLoopThreshold = 3` — Domain-Konstante (gepinnt).
- `ServiceStatus` — 4-Spalten-Shape für die LH-FA-UP-003-Status-
  Anzeige (Name/ContainerStatus/Port/Healthcheck).
- `UpResult` — `Services` + `Stabilized` + `Diagnostics`-Slice.

Generate-Artefakte (M7-T1):
- `Artifact` — 4-Element-Enum (`changelog`, `readme`,
  `env-example`, `devcontainer`) für `LH-FA-GEN-001`.
  Out-of-Range-`String()` rendert als `Artifact(N)` statt
  `unknown` (M7-Review-S4).

Config-Pfade (M8-T1):
- `ConfigPath` — typed Value-Object mit `Kind`/`Service`/
  `WriteAllowed` für die 3-Pfad-Whitelist von `u-boot config`
  (`project.name`, `devcontainer.enabled`,
  `services.<svc>.enabled`). `WriteAllowed=false` für
  `services.<svc>.enabled` schützt vor LH-FA-ADD-005-Lifecycle-
  Bypass (Toggling läuft über `u-boot add`/`remove`).

## Geplante Erweiterungen

- `Port`, `ImageRef` — Value-Objects für Service-Konfiguration
  (V1, wenn weitere Add-ons kommen).
- `ComposeFile`, `EnvVar` — strukturelle Modelle nur falls
  künftige Generators sie über die heutige `text/template`-Stufe
  hinaus brauchen (M7 reicht heute).

## Import-Regeln

Ausschließlich Go-Standard-Library. Verstöße werden durch
`golangci-lint depguard` im `lint`-Stage abgewiesen (`LH-FA-ARCH-003`).
