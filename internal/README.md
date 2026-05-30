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

## Status (M1–M8 Done, MVP vollständig)

Alle hexagonalen Schichten sind seit M6 produktiv besetzt; M7
ergänzte den Generate-Pfad, MVP-Closure den `init --devcontainer`-
Flag, M8 schließt mit `u-boot config`. Die Package-READMEs pflegen
ihren Detail-Stand jeweils selbst; Kurz-Inventar:

- `hexagon/domain/` — Project + Service-Name Value-Objects,
  Diagnostic-Severities (4-stufig inkl. `SeverityInfo`),
  ContainerState / StabilizationOutcome / UpResult für
  `u-boot up` (M6), `Artifact` für `u-boot generate` (M7-T1),
  `ConfigPath` mit `WriteAllowed`-Flag für `u-boot config` (M8-T1).
- `hexagon/application/` — sieben Use-Cases verdrahtet:
  `InitProjectService` (M3 + MVP-Closure-T1), `DoctorService` (M4),
  `AddServiceService` (M5), `UpService` + `DownService` (M6),
  `GenerateService` (M7), `ConfigService` (M8).
- `hexagon/port/driving/` — sieben Use-Case-Interfaces mit
  narrow-scoped Sentinels (M5 add-/M6 up-/M7 generate-/M8 config-
  Familien).
- `hexagon/port/driven/` — `FileSystem`, `YAMLCodec` (mit
  `PatchScalar` + `PatchMappingEntryYAML` + `LocateMarkedEntry`
  und seit V1-yaml-parse-Sentinel auch `ErrYAMLParse`), `Git`,
  `Clock` (mit `Sleep`), `ProgressPort`, `Confirmer`, `Logger`,
  `DockerProbe` (read-only, M4), `DockerEngine` (state-mutierend,
  M6), `NetProbe` (M6).
- `adapter/driven/` — konkrete Implementierungen aller Driven-
  Ports; der YAML-Adapter wrappt seit V1-yaml-parse-Sentinel
  Parse-Fehler über vier Codec-Methoden hinweg konsistent mit
  `driven.ErrYAMLParse`.
- `adapter/driving/cli/` — sieben Cobra-Subcommands plus
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
  [`slice-m6-docker-integrationstests`](../docs/plan/planning/done/slice-m6-docker-integrationstests.md) (Done).
- `acceptance_test.go` (M8-Closure) — benannte Spec-Pins für
  LH-AK-001 (Init+Doctor) und LH-AK-006 (Doppel-Add-Idempotenz);
  LH-AK-007 lebt im `generate_test.go` neben den Set-Helpern;
  LH-AK-002 ist in der Docker-tagged `e2e/`-Suite gepinnt.

## CLI-Subcommands

| Command | Slice | Spec |
| ------- | ----- | ---- |
| `u-boot init [name] [--devcontainer]` | M3 + MVP-Closure-T1 | `LH-FA-INIT-001..007` + `LH-AK-005` |
| `u-boot doctor` | M4 (+ MVP-Closure-T2 Severity-Fix) | `LH-FA-DIAG-001..004` |
| `u-boot add <service>` | M5 | `LH-FA-ADD-001..002`, `-005` |
| `u-boot up` | M6 | `LH-FA-UP-001..003` |
| `u-boot down` | M6 | `LH-FA-UP-004` |
| `u-boot generate <artifact>` | M7 | `LH-FA-GEN-001..005`, `LH-AK-007` |
| `u-boot config [get/set]` | M8 | `LH-FA-CONF-001..005` |

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
