# internal/hexagon/application

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält
keine externe I/O (`LH-FA-ARCH-002`).

## Status

Stand M6: fünf Use-Cases produktiv verdrahtet (M3-T2 bis M6-T5),
keine externe I/O, alle Ports nil-tolerant via package-private
`noop*`-Defaults.

## Inhalt

- `InitProjectService` — `port/driving.InitProjectUseCase`
  (`LH-FA-INIT-001..007`, M3).
- `AddServiceService` — `port/driving.AddServiceUseCase`
  (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`, M5).
- `DoctorService` — `port/driving.DoctorUseCase`
  (`LH-FA-DIAG-001..004`, M4; 11 Checks).
- `UpService` — `port/driving.UpUseCase` (`LH-FA-UP-001..003`, M6-T4).
  Polling-Loop mit `pollInterval=500ms` und `dialTimeout=300ms`,
  fail-safe `ContainerState`-Klassifikation (Dead-Allowlist,
  soft-Unknown, Restart-Loop-Counter mit Threshold 3),
  Healthcheck-dominanter Stabilisierungs-Vertrag mit TCP-Port-
  Probe als Warn-Diagnose (§141 / §968).
- `DownService` — `port/driving.DownUseCase` (`LH-FA-UP-004`, M6-T5).
  §T5-Truth-Table für den `--volumes`-Bestätigungs-Pfad (4 Zeilen
  × 2 Sub-Cases bei AssumeYes-und-NonInteractive); ruft
  `Confirmer.ConfirmRemoveVolumes` nur im interaktiven Pfad.

Plus Helper:
- `parseComposePort` (M6-T4-fund) — pure 8-Syntax-Cases-Parser für
  Compose-`ports:`-Array-Elemente; nicht-probebare Formen
  (UDP/Range/Unknown) returnen `probable=false` für Warn-Diagnose-
  Pfad.
- `managedblock/` — Marker-basiertes YAML-Blocksetting (M3+M5).
- `templates/` — Embedded Go-Templates für `u-boot init`/`add`.

## Geplante Erweiterungen

- `GenerateService` — `port/driving.GenerateUseCase`
  (`LH-FA-GEN-001..005`, M7).
- `ConfigService` — `port/driving.ConfigUseCase`
  (`LH-FA-CONF-001..005`, M8).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port`. **Nicht** erlaubt:
`internal/adapter/*`, externe I/O-Libraries.
