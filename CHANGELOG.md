# Changelog

All notable changes to **u-boot** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Internal `u-boot generate changelog` (`LH-FA-GEN-001..005`, `LH-AK-007`)
maintains a Keep-a-Changelog-formatted changelog for user projects;
this file is the same format applied to u-boot itself.

## [Unreleased]

### Added

- `feat(cli): u-boot doctor --json plus Root-PersistentFlag
  --json (LH-NFA-USE-004 / LH-FA-CLI-007)` — Pattern-Vorbild für
  die maschinen-lesbare CLI-Surface (Cluster-Slice
  `slice-v1-cli-json-dry-run`, Folge-Slice 1/9). **Doctor**
  emittiert mit `--json` einen Spec-§1841-Minimalkontrakt-
  Envelope (`status`/`command`/`diagnostics`/`exitCode`);
  All-OK-Fall ergibt `diagnostics: []` gemäß Lastenheft-Beispiel
  §1846-1852. `SeverityOK` und `SeverityInfo` werden gefiltert
  (Spec §1834 erlaubt `level` nur `warn|error`). `--quiet --json`
  ist semantisch identisch zu reinem `--json`; `--strict --json`
  upgraded Warn auf Exit-Code 11 ohne `status`-Drift (Spec §1837
  koppelt `status` an höchsten `level`, nicht an `--strict`).
  Broken-Pipe-Resistenz: fachlicher Exit-Code 11 hat Vorrang vor
  Write-Fehlern. **Root-PersistentFlag `--json`** für alle 10
  Spec-Enum-Subcommands; nicht-migrierte Forms rejecten mit
  `ErrJSONNotImplemented` (Exit-Code 2, `LH-FA-CLI-006`-Klasse)
  plus Folge-Slice-Verweis. Zentrale Allowlist in
  `cli/root.go` via `PersistentPreRunE`; `--help` als
  read-only Escape-Hatch durchgelassen. **`u-boot template list
  --json`** wandert vom lokalen Flag aufs Root-Flag (beide
  Schreibweisen `template list --json` und `--json template list`
  identisches Output); Envelope-Migration folgt mit
  `slice-v1-cli-json-dry-run-template`, dokumentiert im
  `carveouts.md`-Temporär-Eintrag. **Common-Envelope
  `cliJSONEnvelope`** (`internal/adapter/driving/cli/jsonenvelope.go`)
  trägt Minimalkontrakt- und Voll-Schema-Felder über zwei
  Konstruktoren (`newMinimalEnvelope`/`newFullEnvelope`); Voll-
  Schema-Felder via Pointer-Wrapping (`*bool`/`*[]T`) als
  Anti-Drift gegen `omitempty`-Semantik-Refactor. **Schema-
  Helper-Sub-Package** `internal/adapter/driving/cli/jsontestutil/`
  mit zwei Modi `AssertMinimalEnvelope`/`AssertFullEnvelope`
  (Options-Pattern, kein neuer Dep) plus `DefaultAllowedCodes`-
  Registry für die 13 Doctor-Check-Codes. **Drei aktive Drift-
  Gates** schützen die Code-Registry: (1) `application.DoctorCheckIDs()`
  ↔ Map-Vollständigkeit, (2) Map ↔ Markdown-Roundtrip-Parser
  auf `docs/user/cli-json-output.md` §5.1 mit HTML-Marker-
  Sektion-Begrenzung, (3) Helper-Reject im Acceptance-Pfad
  für undokumentierte Codes. **Schema-Vertrag-Doku**
  ([`docs/user/cli-json-output.md`](docs/user/cli-json-output.md))
  zitiert Minimalkontrakt und Voll-Schema verbatim, dokumentiert
  Code-Registry und Per-Command-Migrations-Reihenfolge.
- `feat(logs): u-boot logs [service] [--follow] [--tail <n>]
  (LH-FA-UP-005)` — neuer Subcommand streamt Compose-Logs als
  V1-Erweiterung der `up`/`down`-Familie. Ohne Service-Argument
  läuft `docker compose logs` über alle Services aus
  `compose.yaml` (T0-(a) Compose-Facade-Semantik, kein
  u-boot.yaml-Filter); mit Service-Argument nur diesen einen
  (Format-Validation via `domain.NewServiceName`, Existenz-Check
  delegiert an Compose). `--follow` blockt bis Ctrl-C und
  beendet via SIGINT-Vertrag-Schicht-1+2+3 (Adapter gibt
  `ctx.Err()` unverdeckt zurück, Use-Case übersetzt zu
  `(LogsResponse{}, nil)`, CLI exit-code 0). `--tail <n>` mit
  Stage-1-Validation auf nicht-negative Ganzzahlen (Default
  leer → Use-Case-Normalisierung zu Compose-`"all"`). Exit-
  Code-Mapping analog `up`/`down`: 10 (User/Project-State),
  11 (`ErrDockerUnavailable`), 12 (`ErrComposeRuntime`),
  14 (FS), 2 (CLI-Usage). Docker-tag E2E-Tests gegen echten
  postgres-Stack pinnen `--tail`-Buffer-Content und
  `--follow`-SIGINT-Vertrag.
- `feat(devcontainer): Devcontainer-Features-Allowlist und Katalog
  (LH-FA-DEV-003)` — 8 Built-in-Features (`git`, `docker-cli`,
  `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`)
  plus External-Source-Allowlist via
  `devcontainer.featureSources.allow`. CLI:
  `u-boot config set devcontainer.features.<name>.{enabled,source,version}`
  plus `--allow-external-feature-sources <url>[,<url>...]` auf
  den drei Spec-§714-717-Pfaden (`init --devcontainer`,
  `generate devcontainer`, `config set devcontainer.featureSources.allow`).
  Doctor-Check `devcontainer.features.allowlist` (Error bei
  Allowlist-Violation, Warn bei Orphan-Activation oder fehlendem
  `enabled:`-Key). User-Doku in
  [`docs/user/devcontainer-features.md`](docs/user/devcontainer-features.md).
- `feat(doctor): devcontainer.features.drift Check` — über-Spec
  Drift-Erkennung zwischen `u-boot.yaml`'s Features-Map und den
  Keys im gerenderten `.devcontainer/devcontainer.json`. Drei
  Warn-Cases: aktiviertes Feature fehlt im JSON (Case 1, inkl.
  Datei-fehlt-Disziplin), deaktiviertes Feature noch im JSON
  (Case 2a), JSON-Key ohne cfg-Pendant (Case 2b, Hand-Edit-Hint).
  Repair-Hint `u-boot generate devcontainer`. Doctor-Total
  steigt 12→13.

## [0.3.0] - 2026-06-01

Third release. Completes the V1 „Add-on Catalogue Expansion"
milestone (5/5): the catalogue now ships three integrated
service add-ons — Postgres (since MVP), Keycloak (LH-FA-ADD-003
/ LH-AK-003) and OpenTelemetry (LH-FA-ADD-004 / LH-AK-004) —
plus the matching `u-boot remove <service>` mirror, the
LH-FA-ADD-006 `--with-deps` dependency-resolution mechanism, and
a doku-only audit closure for three V1 spec-IDs. Architectural
side-effect: the per-service catalogue pattern grew from a flat
`(compose, env, volume)`-tuple in M5 to a declarative entry with
`requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional` /
`healthcheckOptional` / `extraFiles` — any new add-on plugs in
by adding one catalogue row and three templates.

### Verified

- **Three V1 spec-IDs audit-closed** by
  [`slice-v1-audit-done`](docs/plan/planning/done/slice-v1-audit-done.md)
  — Doku-only verification that the existing code/doc state
  already satisfies the requirements:
  `LH-FA-BUILD-006` (Aggregator-Targets `gates`/`ci`/`fullbuild` in
  the Makefile),
  `LH-NFA-MAINT-004` (Add-on and template interfaces documented via
  ADR-0008/-0009 + driving/driven port doc-comments + slice docs),
  `LH-NFA-PORT-003` (u-boot itself runs in container / devcontainer:
  GHCR distroless image + container-aware `doctor` since v0.2.0 +
  six-platform binary distribution + `init --devcontainer`-generated
  files).

### Added

- **`u-boot add otel`** — LH-FA-ADD-004 / LH-AK-004. Third and
  final add-on of the v0.3.0 milestone catalogue. Compose-Service
  mit Image-Pin `otel/opentelemetry-collector:0.108.0` (Stable),
  Port-Mappings `4317:4317` (OTLP/gRPC) + `4318:4318` (OTLP/HTTP),
  `command: --config=/etc/otel-collector-config.yaml`, Bind-Mount
  der gerenderten Config-Datei. Kein Healthcheck im Mindest-
  Setup (LH-AK-004 §2374 toleriert `running` ODER `healthy`).
  Mindest-Collector-Config in `otel-collector-config.yaml`:
  Receivers `otlp/grpc+http`, Processors `batch`, Exporters
  `debug` (stdout), Pipelines `logs`/`metrics`/`traces` — alle
  drei Signal-Typen aus LH-FA-ADD-004 §880.
  
  Internal: Catalogue-Pattern wächst um drei Felder pro Service —
  `extraFiles []extraFileEntry` für whole-file artefacts abseits
  von compose+env+volume (für OTel die Collector-Config-Datei),
  plus `envOptional` (implizit via leerem `envTmpl`) und
  `healthcheckOptional` für Services, die das Standard-Pattern
  legitim nicht brauchen. `executeAdd` schreibt extraFiles als
  vierten Slot nach yaml/compose/env; `executeRemove` löscht sie
  symmetrisch. `serviceComplete` skipt healthcheck-presence für
  `healthcheckOptional`; explicit `healthcheck.disable: true`
  bleibt hart abgelehnt. Acceptance-Helper-Reuse aus
  `slice-v1-keycloak` T3 (`acceptance_helpers.go`) — OTel-E2E
  bleibt ~30 Zeilen. **Makefile-Patch**: `test-docker`-Target
  mountet jetzt `/tmp` host-shared, damit Compose-Bind-Mount-
  Pfade vom Daemon (Host) aufgelöst werden können — sonst sieht
  der Daemon nur den Container-Pfad `t.TempDir()` nicht und
  erstellt einen leeren Verzeichnis-Mount, der den Collector
  beim Config-Read crasht. See
  [`slice-v1-otel`](docs/plan/planning/done/slice-v1-otel.md).
- **`u-boot add keycloak`** — LH-FA-ADD-003 / LH-AK-003. Second
  add-on in the catalogue after Postgres. Compose-Service mit
  Image-Pin `quay.io/keycloak/keycloak:26.0` (LTS), Port-Mapping
  `8080:8080`, `command: start-dev` für LH-AK-003-Boot, Healthcheck
  via `/dev/tcp/localhost/9000` (bash-builtin, kein curl im
  Image) gegen `/health/ready`. Admin-Credentials via Placeholder-
  Env-Block (`KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` +
  `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`).
  **Persistenz: flüchtige H2-In-Container-Datenbank** — kein
  Volume, nach `docker compose down` weg; LH-AK-003 verlangt nur
  Endpoint-200/302. Persistente externe Postgres-Anbindung
  (LH-FA-ADD-003 §857) bleibt als eigener Folge-Slice
  (`slice-v1-keycloak-external-postgres`, Trigger: Nutzer-Bedarf).
  Internal refactor: `renderPostgresTemplates` → generischer
  `renderServiceTemplates(svc)` über neue Service-Catalogue-
  Tabelle; `hasRequiredEnvKeys` / `contentScanState` /
  `inspectVolumeArtefact` / `patchTargetsFor` werden per-Service
  über `requiredEnvKeys` / `volumeRefLiteral` / `volumeOptional`
  parametrisiert, damit Keycloak's volume-loser Pfad nicht in
  den Postgres-Repair-Loop läuft. Test-Helper-Extraktion:
  `internal/e2e/acceptance_helpers.go` teilt die init+add+up-
  Pipeline mit dem LH-AK-002-Postgres-Test (Boot-Zeit-Carveout
  für Keycloak: 4 min UpService-Timeout vs. 90 s Postgres). See
  [`slice-v1-keycloak`](docs/plan/planning/done/slice-v1-keycloak.md).
- **`u-boot add <service> --with-deps`** — LH-FA-ADD-006 add-on
  dependency mechanism. New domain type `AddOnDependency` (path-
  conditional service dependency declaration) + per-service
  catalogue side-table `dependenciesFor(svc)` (Postgres has none
  today; first non-nil row lands with `slice-v1-keycloak`). When
  the requested add-on declares a dep that is not yet registered
  in `u-boot.yaml`, the four-mode dispatch decides what happens:
  `--with-deps` auto-installs the chain (recursive `Add` calls,
  flag inherited so transitive deps follow); `--yes` has the same
  effect; `--no-interactive` (without `--yes`/`--with-deps`)
  fails fast with the new `ErrDependenciesRequired` sentinel
  (exit 10); default-interactive prompts via the new
  `Confirmer.ConfirmAddDependency(ctx, svc, missing)` driven-port
  method (mirror of `ConfirmRemoveVolumes` from M6). Postgres-
  only flows are unchanged — the no-deps short-circuit keeps the
  load+resolve cost out of the MVP catalogue path. Breaking
  refactor in the application layer: `NewAddServiceService`
  now takes a `Confirmer` between `yaml` and `logger`; all eight
  callsites updated in lock-step. See
  [`slice-v1-addons-deps`](docs/plan/planning/done/slice-v1-addons-deps.md).
- **`u-boot remove <service> [--purge]`** — first slice of the
  v0.3.0 milestone ("Add-on Catalogue Expansion"). Mirror of
  `u-boot add`: detects the LH-FA-ADD-005 service state, strips
  the `service.<name>` managed block from `compose.yaml` and
  `.env.example`, then sets `services.<name>.enabled: false` in
  `u-boot.yaml`. Idempotent: removing an already-disabled service
  is a no-op with a clear message. Inconsistent project state
  (orphan block, missing entry) surfaces as `ErrServiceInconsistent`
  with a manual-cleanup hint. New driving sentinel
  `ErrServiceUnregistered` (exit 10) distinguishes "service was
  never added" from "service name not in the catalogue"
  (`ErrServiceUnsupported`). LH-FA-ADD-007 §"Volumes nur auf
  explizite Anforderung": `--purge` opts in destructively and
  triggers the LH-FA-CLI-005A §254 confirmation gate (mirror of
  `u-boot down --volumes`); auto-removal of volumes is deferred
  to a follow-up slice — v0.3.0's `--purge` summary points at
  `docker volume rm <name>` for the manual cleanup. Internal:
  `detectServiceState` extracted from the M5 add path to a
  package-level function so both add and remove share it without
  duplication. See
  [`slice-v1-add-remove`](docs/plan/planning/done/slice-v1-add-remove.md).

## [0.2.0] - 2026-06-01

Second release. Adds the first two V1 template features
(`template list` + `init --template`), a cross-platform binary
distribution (six platforms as GitHub-Release assets), and a
container-aware `doctor` that no longer mis-reports a healthy host
as 4 errors when run from inside the distroless image. v0.1.1 was
originally planned as a patch-only tag for the doctor fix but is
skipped in favour of this minor bump — three features landed before
the tag-push and strict SemVer wants a MINOR bump for them.

### Added

- **`u-boot template list [--json]`** — first V1 template
  subcommand (LH-FA-TPL-004). Enumerates the built-in project-
  template catalog with name, description, and version in a
  tabwriter-aligned table; `--json` emits a structured array
  with the full LH-FA-TPL-002 metadata surface (`supportedAddOns`,
  `generatedFiles`, `requiredTools`, `variables`). Bootstrap
  built-in: `basic` (one template; further built-ins follow on
  demand per ADR-0009 §Folgepunkte 4). Fully hexagonal:
  `domain.TemplateMetadata` + `Validate()` (kebab-case-name
  regex, `ErrInvalidTemplate` sentinel), driven port
  `TemplateCatalog`, embed.FS-backed `externaltemplates` adapter,
  application `TemplateListService` (multi-`%w` so the original
  `domain.ErrInvalidTemplate` chain survives), CLI
  `template list` rendering. Adapter directory consolidated to
  `internal/adapter/driven/externaltemplates/` (no hyphen) for
  consistency with the existing `driven/`-adapter naming; ADR-0009
  §Entscheidung updated to match. See
  [`slice-v1-template-list`](docs/plan/planning/done/slice-v1-template-list.md).
- **`u-boot init <name> --template <name>`** — second V1 template
  feature, the render path of LH-FA-TPL-001 / LH-FA-TPL-002. The
  init service delegates file rendering to the new
  `TemplateInitService` when `--template` is set; project structure
  directories and `git init` stay with the InitProjectService so
  the user-observable flow is one command. Byte-identity
  guarantee: `u-boot init demo --template basic` produces a
  project byte-identical to `u-boot init demo` for the six default
  files (`u-boot.yaml`, `compose.yaml`, `README.md`,
  `CHANGELOG.md`, `.env.example`, `.gitignore`) — pinned by an
  E2E `diff -r` test against the production catalog. Render engine:
  Go `text/template` for `*.tmpl` files, 1:1 copy for non-`.tmpl`
  files (per ADR-0009 §Entscheidung); `template.yaml` metadata is
  skipped. Two-phase render-then-write: a render error in any file
  short-circuits before the first disk write, so a buggy template
  no longer leaves a half-populated project. New
  `domain.TemplatePath` validator rejects `..` segments, absolute
  paths, Windows drive letters, backslashes, NUL bytes, and empty
  strings (LH-FA-CLI-006 exit 10 via `ErrInvalidTemplatePath`).
  Mutex with `--devcontainer`/`--force`/`--backup`: surfaces as
  `ErrTemplateConflictsWithFlag` (exit 2) — v1 is fresh-init-only.
  Soft-existing-detection is skipped on the template path because
  `--template` resolves the "is this an existing project?"
  ambiguity by intent; the hard-existing check
  (`u-boot.yaml` present → `ErrProjectExists`) remains the
  safety net. Variable resolution + `--var key=value` deferred to
  a future slice (basic has no variables). See
  [`slice-v1-template-init`](docs/plan/planning/done/slice-v1-template-init.md).
- **Cross-platform binary distribution** for six platforms
  (Linux/macOS/Windows × amd64/arm64). `make build-binaries`
  cross-compiles every supported `GOOS`/`GOARCH` combination via
  the pinned `golang:$(GO_VERSION)` builder image (CGO disabled,
  `-ldflags "-s -w -X main.version=$(VERSION)"`, output to
  `bin/u-boot-<os>-<arch>[.exe]`). `.github/workflows/publish.yml`
  builds the same set after the GHCR push on every `v*` tag and
  attaches them as GitHub-Release assets via `gh release upload`.
  See
  [`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md)
  — ADR-0007 §Folgepunkte 1 trigger pulled forward by the
  doctor-container-awareness feedback.
- Quickstart in `README.md` / `README.de.md` gets a host-native
  install block (`curl -sSL … | chmod +x` for Linux/macOS,
  `Invoke-WebRequest` for Windows) as the primary recommended path;
  the GHCR `docker run …` block is demoted to "alternative for
  container/CI workflows".
- `internal/hexagon/port/driven.RuntimeEnvironment` port plus
  `internal/adapter/driven/runtime.FileEnv` adapter: best-effort
  container detection via `/.dockerenv` (Docker Engine / Desktop)
  and `/run/.containerenv` (Podman / CRI-O / buildah). Drives the
  doctor-container-awareness change below.

### Changed

- `u-boot doctor` now skips the four host-prerequisite checks
  (`git.installed`, `docker.installed`, `docker.reachable`,
  `docker.compose.installed`) when running inside a container, with
  a `SeverityInfo` diagnostic and a hint that points at
  [`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md)
  for the rationale. Effect: `docker run --rm
  ghcr.io/pt9912/u-boot:0.2.0 doctor` no longer mis-reports a
  healthy host as 4 errors; exit code on an otherwise-clean project
  goes from 11 to 0. This addresses real-world feedback from the
  first `v0.1.0` GHCR pull (2026-05-31) where the distroless image's
  lack of bundled `docker` / `git` binaries surfaced as false-
  positive errors.
- Host installations are unaffected — `runtime.FileEnv` returns
  `false` outside containers, so the existing
  `LH-FA-DIAG-002`-classified errors / warnings remain.

### Notes

`releases/latest/download/u-boot-<os>-<arch>[.exe]` resolves to the
highest stable tag — since `v0.1.0` predates binary assets, the
`latest`-shortcut starts working with `v0.2.0` (or any later tag).

## [0.1.0] - 2026-05-31

First public release. Closes the MVP scope from
[`spec/lastenheft.md`](spec/lastenheft.md) MVP-priority IDs: all
`LH-AK-*`, `LH-FA-*` and `LH-SA-*` items are delivered (audit trail
in [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md)
§MVP-Bilanz).

### Added — Subcommands

- `u-boot init [name] [--devcontainer]` — generate project skeleton
  (`u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`,
  `.env.example`, `.gitignore`, `docker/`, `scripts/`, `docs/`)
  and run `git init` (`LH-FA-INIT-001..007`, `LH-AK-001`, `LH-AK-005`).
  Mode flags `--yes` / `--no-interactive` / `--assume-existing`
  (`LH-FA-CLI-005A`); re-init with `--force` / `--backup` for
  managed-block edits vs. full overwrite (`LH-FA-INIT-005`).
- `u-boot doctor [--strict]` — 11 diagnostic checks against the
  local environment + project, severity-classified
  (ok / warn / error), repair-hint output, exit-code 11 on errors
  (or warns under `--strict`) (`LH-FA-DIAG-001..004`).
- `u-boot add <service>` — idempotent state-machine for service
  add-ons; today's catalogue: `postgres` only (`LH-FA-ADD-001/-002/-005`,
  `LH-AK-002`, `LH-AK-006`). Keycloak (`LH-AK-003`) and
  OpenTelemetry (`LH-AK-004`) follow in V1.
- `u-boot up [--timeout <sec>]` and `u-boot down [--volumes]` —
  Compose wrapper with healthcheck polling and TCP port probes
  (`LH-FA-UP-001..004`). `--timeout 0` is fire-and-forget.
- `u-boot generate <changelog|readme|env-example|devcontainer>` —
  idempotent block-replace via the `U-BOOT MANAGED BLOCK: init`
  marker; user content outside the managed region is preserved
  byte-identically. `changelog` carries the `LH-AK-007` pin
  (no destructive edits to existing entries). Exit codes
  `0` / `2` / `10` / `14` per `LH-FA-CLI-006` (`LH-FA-GEN-001..005`,
  `LH-FA-DEV-001/004/005`).
- `u-boot config [get <path> | set <path> <value>]` — whitelist-
  scoped reads/writes with two-stage schema validation (struct
  unmarshal + per-path domain re-validation) before any
  `WriteFile`. `services.<svc>.enabled` is get-only; toggling
  happens through `add` / `remove` to keep the add-on state
  machine atomic (`LH-FA-CONF-001..005`).

### Added — CI & release infrastructure

- GitHub Actions CI workflow `.github/workflows/ci.yml` with three
  PR-blocking jobs (`LH-QA-003`): `gates (lint + test +
  coverage-gate)`, `security-gates (govulncheck)`,
  `image-scan (trivy HIGH+CRITICAL)`. All actions SHA-pinned;
  Docker-only runner (`LH-FA-BUILD-007`); per-job minimal
  permissions.
- GitHub Actions release workflow
  `.github/workflows/publish.yml` triggered on `v*` tags. Strict
  SemVer-2.0 validation (rejects leading-zero numeric prereleases
  and build-metadata `+...` tags), GHCR image push to
  `ghcr.io/pt9912/u-boot:<version>` (plus `:latest` for stable
  tags), OCI label verification, and live `--version` smoke test
  against the tag-derived `VERSION`.
- Local outer/inner-loop parity: `make image-scan` reproduces the
  `image-scan` CI job using the same Trivy version
  (`TRIVY_VERSION ?= 0.70.0`) the action installs.
- Multi-stage distroless runtime image (`gcr.io/distroless/static-debian12:nonroot`)
  built via `make build`; CGO-disabled static binary; version
  injected at build time as `-X main.version=<UBOOT_VERSION>`
  and as the `org.opencontainers.image.version` label.

### Added — Architecture & documentation

- Hexagonal architecture (`LH-FA-ARCH-001..003`, ADR-0002):
  `internal/hexagon/{domain,application,port/{driving,driven}}`
  + `internal/adapter/{driving,driven}`. `depguard` enforces
  layer rules in CI.
- 10 ADRs cover language (Go), architecture (hexagonal), lint
  profile (SOLID-near), CI system, CLI framework (Cobra), revive
  custom rules, distribution path (GHCR), plugin system (static —
  no plugins), template format (YAML + Go `text/template`), and
  the HTTP adapter (not built; CLI-only).
- User-facing setup docs:
  [`docs/user/quality.md`](docs/user/quality.md) (quality-gates
  overview) and
  [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
  (one-time GitHub UI activation of required status checks).
- German `spec/lastenheft.md` (~3000 lines, 14 sections + 4 open
  points all decided) is the single source of truth; English
  `README.md` / German `README.de.md` are equivalent.

### Known limitations and deliberate carve-outs

- **Add-on catalogue is intentionally small:** only `postgres`
  ships in v0.1.0. Keycloak and OpenTelemetry are V1.
- **Templates implementation is V1.** Format is decided
  (ADR-0009: YAML + `text/template`); the three implementation
  slices (`slice-v1-template-list`, `slice-v1-template-init`,
  `slice-later-local-templates`) follow on demand.
- **JSON / machine-readable output is V1.** `--json` and
  `--dry-run` flags (`LH-FA-CLI-007/008`, `LH-NFA-USE-004`) are
  not yet shipped; ADR-0010 (no HTTP adapter) explicitly relies
  on this V1 track landing.
- **Distribution is GHCR-only.** Binary, Homebrew, Debian/RPM
  paths are deferred with explicit trigger slices in ADR-0007.
  `npm` / `pip` are rejected (ecosystem mismatch).
- **No plugin loader.** Add-on system stays statically compiled
  into u-boot (ADR-0008). Four re-evaluation triggers documented
  in ADR-0008 §Folgepunkte.
- **CLI-only.** No HTTP / daemon adapter (ADR-0010); programmatic
  consumers use subprocess + `--json` once V1 lands.
- **Inner-loop is Docker-only** (`LH-FA-BUILD-007`). GNU `make`
  remains the single non-Docker host dependency (permanent
  carve-out to `LH-NFA-PORT-002`).

### Setup — required one-time GitHub UI activation

Before merging external PRs against `main`, activate the three
required status checks in GitHub UI per
[`docs/user/branch-protection.md`](docs/user/branch-protection.md):
the exact match strings are the workflow `name:` fields
(`gates (lint + test + coverage-gate)`,
`security-gates (govulncheck)`,
`image-scan (trivy HIGH+CRITICAL)`), not the shorter `jobs.<key>`
identifiers.

[Unreleased]: https://github.com/pt9912/u-boot/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/pt9912/u-boot/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pt9912/u-boot/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pt9912/u-boot/releases/tag/v0.1.0
