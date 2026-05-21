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

## Status (M3-T1)

Mit M3-T1 ([`docs/plan/planning/in-progress/slice-m3-init-flow.md`](../docs/plan/planning/in-progress/slice-m3-init-flow.md))
sind die ersten produktiven Pakete entstanden:

- `hexagon/domain/`: `Project`, `ProjectName` (mit `NormalizeProjectName`).
- `hexagon/port/driven/`: `FileSystem`, `YAMLCodec`, `Git`, `Clock`.
- `adapter/driven/{fs,yaml,git,clock}/`: konkrete Implementierungen.

Noch leer und folgen mit M3-T2 / M3-T3:

- `hexagon/application/` (M3-T2): `InitProjectService`.
- `hexagon/port/driving/` (M3-T2): `InitProjectUseCase`.
- `adapter/driving/` (M3-T3): CLI-Commands (Cobra).

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
