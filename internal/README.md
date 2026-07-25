# internal/

Nicht-exportierbare Go-Pakete für u-boot. Strukturiert nach dem
hexagonalen Architektur-Pattern ([`LH-FA-ARCH-001`](../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement); Detail in
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
  fahren ([`LH-AK-002`](../spec/lastenheft.md#lh-ak-002--postgresql-flow) PostgreSQL-Acceptance,
  [`LH-FA-UP-004`](../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) §1015 Volume-Removal). Laufen ausschließlich
  über `make test-docker` — siehe
  [`docs/user/quality.md`](../docs/user/quality.md) §2.2.
- `acceptance_test.go` — benannte Spec-Pins für
  [`LH-AK-001`](../spec/lastenheft.md#lh-ak-001--minimaler-init-flow) (Init+Doctor) und [`LH-AK-006`](../spec/lastenheft.md#lh-ak-006--idempotenz) (Doppel-Add-Idempotenz);
  [`LH-AK-007`](../spec/lastenheft.md#lh-ak-007--changelog-generator) lebt im `generate_test.go` neben den Set-Helpern;
  [`LH-AK-002`](../spec/lastenheft.md#lh-ak-002--postgresql-flow) ist in der Docker-tagged `e2e/`-Suite gepinnt.

## CLI-Subcommands

| Command | Spec |
| ------- | ---- |
| `u-boot init [name] [--devcontainer]` | [`LH-FA-INIT-001`](../spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../spec/lastenheft.md#lh-fa-init-007--git-repository-initialisierung) + [`LH-AK-005`](../spec/lastenheft.md#lh-ak-005--devcontainer-flow) |
| `u-boot doctor` | [`LH-FA-DIAG-001`](../spec/lastenheft.md#lh-fa-diag-001--doctor-befehl)..[`LH-FA-DIAG-004`](../spec/lastenheft.md#lh-fa-diag-004--reparaturhinweise) |
| `u-boot add <service>` | [`LH-FA-ADD-001`](../spec/lastenheft.md#lh-fa-add-001--add-on-befehl)..[`LH-FA-ADD-006`](../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) |
| `u-boot remove <service>` | [`LH-FA-ADD-007`](../spec/lastenheft.md#lh-fa-add-007--service-entfernen) |
| `u-boot up` | [`LH-FA-UP-001`](../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)..[`LH-FA-UP-003`](../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen) |
| `u-boot down` | [`LH-FA-UP-004`](../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) |
| `u-boot logs [service]` | [`LH-FA-UP-005`](../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) |
| `u-boot generate <artifact>` | [`LH-FA-GEN-001`](../spec/lastenheft.md#lh-fa-gen-001--generate-befehl)..[`LH-FA-GEN-005`](../spec/lastenheft.md#lh-fa-gen-005--idempotenz), [`LH-AK-007`](../spec/lastenheft.md#lh-ak-007--changelog-generator) |
| `u-boot config [get/set]` | [`LH-FA-CONF-001`](../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)..[`LH-FA-CONF-005`](../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) |
| `u-boot template list` | [`LH-FA-TPL-001`](../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen)..[`LH-FA-TPL-004`](../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten) |

Alle Subkommandos tragen `--json`, `--dry-run` und `--diff`.

## Coverage

`./internal/...` ist der Coverage-Scope ([`LH-FA-BUILD-008`](../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap),
[`LH-FA-BUILD-009`](../spec/lastenheft.md#lh-fa-build-009--repository-layout)); `./cmd/...` ist ausgeschlossen. Der Schwellwert ist
aktiv und wird über `make coverage-gate` durchgesetzt.

## Import-Regeln

Verbindliche Schicht-Regeln in [`LH-FA-ARCH-003`](../spec/lastenheft.md)
und [`../spec/architecture.md`](../spec/architecture.md). Enforcement
über `golangci-lint depguard` im `lint`-Stage; `//nolint:depguard` ist
verboten.
