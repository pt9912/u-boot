# Changelog

All notable changes to **u-boot** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Internal `u-boot generate changelog` (`LH-FA-GEN-001..005`, `LH-AK-007`)
maintains a Keep-a-Changelog-formatted changelog for user projects;
this file is the same format applied to u-boot itself.

## [Unreleased]

### Added

- **Cross-platform binary distribution** for six platforms
  (Linux/macOS/Windows × amd64/arm64). `make build-binaries`
  cross-compiles every supported `GOOS`/`GOARCH` combination via
  the pinned `golang:$(GO_VERSION)` builder image (CGO disabled,
  `-ldflags "-s -w -X main.version=$(VERSION)"`, output to
  `bin/u-boot-<os>-<arch>[.exe]`). `.github/workflows/publish.yml`
  builds the same set after the GHCR push on every `v*` tag and
  attaches them as GitHub-Release assets via `gh release upload`.
  See
  [`slice-v2-binary-distribution`](docs/plan/planning/open/slice-v2-binary-distribution.md)
  — ADR-0007 §Folgepunkte 1 trigger pulled forward by the v0.1.1
  doctor-container-awareness feedback.
- Quickstart in `README.md` / `README.de.md` gets a host-native
  install block (`curl -sSL … | chmod +x` for Linux/macOS,
  `Invoke-WebRequest` for Windows) as the primary recommended path;
  the GHCR `docker run …` block is demoted to "alternative for
  container/CI workflows".

### Notes

`releases/latest/download/u-boot-<os>-<arch>[.exe]` resolves to the
highest stable tag — since `v0.1.0` predates binary assets, the
`latest`-shortcut starts working once `v0.1.1` (or any later tag)
has been pushed.

## [0.1.1] - TBD

Targeted patch addressing the real-world feedback from the first
`v0.1.0` GHCR pull (2026-05-31): `docker run ghcr.io/pt9912/u-boot:0.1.0
doctor` reported four false-positive errors against a healthy host
because the distroless image bundles no `docker` / `git` binaries
to probe.

### Added

- `internal/hexagon/port/driven.RuntimeEnvironment` port plus
  `internal/adapter/driven/runtime.FileEnv` adapter:
  best-effort container detection via `/.dockerenv` (Docker Engine /
  Desktop) and `/run/.containerenv` (Podman / CRI-O / buildah).

### Changed

- `u-boot doctor` now skips the four host-prerequisite checks
  (`git.installed`, `docker.installed`, `docker.reachable`,
  `docker.compose.installed`) when running inside a container, with
  a `SeverityInfo` diagnostic and a hint that points at
  [`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md)
  for the rationale. Effect: `docker run --rm
  ghcr.io/pt9912/u-boot:0.1.1 doctor` no longer mis-reports a
  healthy host as 4 errors; exit code on an otherwise-clean project
  goes from 11 to 0.
- Host installations are unaffected — `runtime.FileEnv` returns
  `false` outside containers, so the existing
  `LH-FA-DIAG-002`-classified errors / warnings remain.

### Notes

The medium-term fix is a host-native binary distribution
([`slice-v2-binary-distribution`](docs/plan/planning/open/slice-v2-binary-distribution.md),
ADR-0007 §Folgepunkte 1 trigger now active); the v0.1.1 skip is the
short-term ergonomic patch.

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

[Unreleased]: https://github.com/pt9912/u-boot/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/pt9912/u-boot/releases/tag/v0.1.0
