# internal/

Nicht-exportierbare Go-Pakete für u-boot. Strukturiert nach dem
hexagonalen Architektur-Pattern (`LH-FA-ARCH-001..003`; Detail in
[`../spec/architecture.md`](../spec/architecture.md)).

## Layout

```
internal/
├── hexagon/
│   ├── domain/              # reine Datentypen, keine I/O
│   ├── application/         # Use-Cases; ruft nur Ports auf
│   └── port/
│       ├── driving/         # Interfaces, die CLI/HTTP konsumiert
│       └── driven/          # Interfaces, die Application nach außen ruft
└── adapter/
    ├── driving/             # konkrete Driver (cli/, …)
    └── driven/              # konkrete Adapter (docker/, fs/, yaml/, …)
```

## Status (M1–M6 Done)

Alle hexagonalen Schichten sind seit M6 produktiv besetzt. Die
Package-READMEs pflegen ihren Detail-Stand jeweils selbst; Kurz-
Inventar:

- `hexagon/domain/` — Project + Service Value-Objects,
  Diagnostic-Severities (4-stufig inkl. `SeverityInfo` aus M6),
  ContainerState / StabilizationOutcome / UpResult für
  `u-boot up` (M6).
- `hexagon/application/` — fünf Use-Cases verdrahtet:
  `InitProjectService` (M3), `DoctorService` (M4),
  `AddServiceService` (M5), `UpService` + `DownService` (M6).
- `hexagon/port/driving/` — fünf Use-Case-Interfaces mit narrow-
  scoped Sentinels (`ErrProjectExists`, …,
  `ErrStabilizationTimeout`, `ErrConfirmationRequired`).
- `hexagon/port/driven/` — `FileSystem`, `YAMLCodec`, `Git`,
  `Clock` (mit `Sleep` seit M6-T4-fund), `ProgressPort`,
  `Confirmer` (2 Methoden), `Logger`, `DockerProbe` (read-only,
  M4), `DockerEngine` (state-mutierend, M6) und `NetProbe` (M6).
- `adapter/driven/` — konkrete Implementierungen, plus
  `docker/engine.go` + `netprobe/probe.go` + `Clock.Sleep` als
  M6-Zugänge.
- `adapter/driving/cli/` — fünf Cobra-Subcommands plus
  Status-Renderer; persistente `--quiet`/`--verbose`/`--debug`-
  Flags steuern seit [`slice-followup-verbosity-wiring`](../docs/plan/planning/done/slice-followup-verbosity-wiring.md)
  den `slog.Level` zur Laufzeit (`PersistentPreRunE` mutiert ein
  `*slog.LevelVar`, das mit dem Logger-Adapter geteilt wird).
- `e2e/` — `//go:build docker`-Integrationstests, die mehrere
  Application-Services in Sequenz gegen eine echte Compose-Engine
  fahren (`LH-AK-002` PostgreSQL-Acceptance,
  `LH-FA-UP-004` §1015 Volume-Removal). Laufen ausschließlich
  über `make test-docker` — siehe
  [`docs/user/quality.md`](../docs/user/quality.md) §2.2 und
  [`slice-m6-docker-integrationstests`](../docs/plan/planning/in-progress/slice-m6-docker-integrationstests.md).

## CLI-Subcommands

| Command | Slice | Spec |
| ------- | ----- | ---- |
| `u-boot init [name]` | M3 | `LH-FA-INIT-001..007` |
| `u-boot doctor` | M4 | `LH-FA-DIAG-001..004` |
| `u-boot add <service>` | M5 | `LH-FA-ADD-001..002`, `-005` |
| `u-boot up` | M6 | `LH-FA-UP-001..003` |
| `u-boot down` | M6 | `LH-FA-UP-004` |

## Coverage

`./internal/...` ist seit M3-T1 der Coverage-Scope (`LH-FA-BUILD-008`,
`LH-FA-BUILD-009`). Bootstrap-Modus ist ab T1 verlassen; aktuelle
Messung liegt über 90 %. Schwellwert wird in M3-T5 von `0` auf `80`
gehoben (siehe
[`slice-m3-coverage-threshold-aktivieren.md`](../docs/plan/planning/in-progress/slice-m3-coverage-threshold-aktivieren.md)).

## Import-Regeln

Verbindliche Schicht-Regeln in [`LH-FA-ARCH-003`](../spec/lastenheft.md)
und [`../spec/architecture.md`](../spec/architecture.md). Enforcement
über `golangci-lint depguard` im `lint`-Stage; `//nolint:depguard` ist
verboten.
