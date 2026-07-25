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

## Inventar

Alle hexagonalen Schichten sind produktiv besetzt. Die Package-READMEs
pflegen ihren Detail-Stand jeweils selbst; Kurz-Inventar:

- `hexagon/domain/` — Value-Objects für Projekt- und Service-Namen,
  Diagnostic-Severities (4-stufig inkl. `SeverityInfo`),
  ContainerState / StabilizationOutcome / UpResult für
  `u-boot up`, `Artifact` für `u-boot generate`,
  `ConfigPath` mit `WriteAllowed`-Flag für `u-boot config`,
  Vorlagen-Metadaten für `u-boot template`.
- `hexagon/application/` — ein Service je Use-Case-Familie:
  `InitProjectService`, `DoctorService`, `AddServiceService`,
  `RemoveServiceService`, `UpService`, `DownService`, `LogsService`,
  `GenerateService`, `ConfigService`, `TemplateListService`,
  `TemplateInitService`.
- `hexagon/port/driving/` — je Use-Case-Familie ein Interface mit
  eigenem Request-/Response-Paar plus narrow-scoped Sentinels für die
  Exit-Code-Klassifikation.
- `hexagon/port/driven/` — `FileSystem`, `YAMLCodec` (mit
  `PatchScalar` + `PatchMappingEntryYAML` + `LocateMarkedEntry`
  und `ErrYAMLParse`), `Git`,
  `Clock` (mit `Sleep`), `ProgressPort`, `Confirmer`, `Logger`,
  `DockerProbe` (read-only), `DockerEngine` (state-mutierend),
  `NetProbe`, `TemplateCatalog`, `TemplateFiles`, `RecorderPort`,
  `RuntimeEnvironment`.
- `adapter/driven/` — konkrete Implementierungen aller Driven-
  Ports; jeder Adapter pinnt sein Port-Interface im Produktivcode,
  sodass Drift den Package-Build bricht. Der YAML-Adapter wrappt
  Parse-Fehler über vier Codec-Methoden hinweg konsistent mit
  `driven.ErrYAMLParse`.
- `adapter/driving/cli/` — ein Cobra-Command je Subkommando plus
  Status-Renderer; persistente `--quiet`/`--verbose`/`--debug`-
  Flags steuern den `slog.Level` zur Laufzeit (`PersistentPreRunE`
  mutiert ein `*slog.LevelVar`, das mit dem Logger-Adapter geteilt
  wird).
- `e2e/` — `//go:build docker`-Integrationstests, die mehrere
  Application-Services in Sequenz gegen eine echte Compose-Engine
  fahren (`LH-AK-002` PostgreSQL-Acceptance,
  `LH-FA-UP-004` §1015 Volume-Removal). Laufen ausschließlich
  über `make test-docker` — siehe
  [`docs/user/quality.md`](../docs/user/quality.md) §2.2.
- `acceptance_test.go` — benannte Spec-Pins für
  [`LH-AK-001`](../spec/lastenheft.md#lh-ak-001--minimaler-init-flow) (Init+Doctor) und [`LH-AK-006`](../spec/lastenheft.md#lh-ak-006--idempotenz) (Doppel-Add-Idempotenz);
  [`LH-AK-007`](../spec/lastenheft.md#lh-ak-007--changelog-generator) lebt im `generate_test.go` neben den Set-Helpern;
  [`LH-AK-002`](../spec/lastenheft.md#lh-ak-002--postgresql-flow) ist in der Docker-tagged `e2e/`-Suite gepinnt.

## CLI-Subcommands

| Command | Spec |
| ------- | ---- |
| `u-boot init [name] [--devcontainer]` | `LH-FA-INIT-001..007` + `LH-AK-005` |
| `u-boot doctor` | `LH-FA-DIAG-001..004` |
| `u-boot add <service>` | `LH-FA-ADD-001..006` |
| `u-boot remove <service>` | `LH-FA-ADD-007` |
| `u-boot up` | `LH-FA-UP-001..003` |
| `u-boot down` | `LH-FA-UP-004` |
| `u-boot logs [service]` | `LH-FA-UP-005` |
| `u-boot generate <artifact>` | `LH-FA-GEN-001..005`, `LH-AK-007` |
| `u-boot config [get/set]` | `LH-FA-CONF-001..005` |
| `u-boot template list` | `LH-FA-TPL-001..004` |

Alle Subkommandos tragen `--json`, `--dry-run` und `--diff`.

## Coverage

`./internal/...` ist der Coverage-Scope (`LH-FA-BUILD-008`,
`LH-FA-BUILD-009`); `./cmd/...` ist ausgeschlossen. Der Schwellwert ist
aktiv und wird über `make coverage-gate` durchgesetzt.

## Import-Regeln

Verbindliche Schicht-Regeln in [`LH-FA-ARCH-003`](../spec/lastenheft.md)
und [`../spec/architecture.md`](../spec/architecture.md). Enforcement
über `golangci-lint depguard` im `lint`-Stage; `//nolint:depguard` ist
verboten.
