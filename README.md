# u-boot

**A developer environment bootloader for Docker-based projects.**

`u-boot` is a CLI that bootstraps reproducible development environments:
project structure, Docker Compose stack, devcontainer configuration,
service add-ons (PostgreSQL, Keycloak, OpenTelemetry, …), and the usual
recurring artefacts (README, CHANGELOG, `.env.example`).

> **Sprachversion:** Die deutsche Variante dieses README liegt unter
> [`README.de.md`](README.de.md). Das Lastenheft
> ([`spec/lastenheft.md`](spec/lastenheft.md)) ist auf Deutsch verfasst;
> CLI-Ausgaben und erzeugte Dateien sind auf Englisch (`LH-LESE-002`).

## Status

**MVP vollständig — seven subcommands fully wired (`init` + `doctor` + `add` + `up` + `down` + `generate` + `config`).**

Every MVP-priority `LH-AK-*`, `LH-FA-*` and `LH-SA-*` ID from
[`spec/lastenheft.md`](spec/lastenheft.md) is delivered. **`v0.1.0`
is released (2026-05-31)** — see
[GitHub Release](https://github.com/pt9912/u-boot/releases/tag/v0.1.0)
and the GHCR image at `ghcr.io/pt9912/u-boot:0.1.0` (plus the
stable-floating `:latest`). Distribution policy in
[ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md). Audit
trail in the
[roadmap MVP-Bilanz block](docs/plan/planning/in-progress/roadmap.md)
and the
[release-cut slice](docs/plan/planning/done/slice-v1-release-cut-v0.1.0.md).

- `u-boot init [name] [--devcontainer]` creates the
  LH-FA-INIT-003 project structure plus `u-boot.yaml`
  (LH-FA-CONF-002) and runs `git init` by default
  (LH-FA-INIT-007). `--force` / `--backup` drive the LH-FA-INIT-005
  overwrite-protection (managed-block edits vs full overwrite with
  `.bak[.N]`); `--yes` / `--no-interactive` / `--assume-existing`
  are the LH-FA-CLI-005A mode flags (the latter drives the
  LH-FA-INIT-004 soft-detection). `--devcontainer` (LH-AK-005) also
  writes `.devcontainer/devcontainer.json` + `Dockerfile` and sets
  `devcontainer.enabled: true` in `u-boot.yaml`.
- `u-boot doctor` runs 11 diagnostic checks against the local
  environment + project (LH-FA-DIAG-002), classifies findings as
  ok / warn / error (LH-FA-DIAG-003), prints repair hints
  (LH-FA-DIAG-004) and exits 11 on any error (or any warn with
  `--strict`). M5 adds `services.enabled-key`,
  `devcontainer.forwardPorts.consistency`, and severity escalation
  based on `devcontainer.enabled`. MVP-Closure-T2 retargets
  `compose.yaml.valid` no-services from Error → Warn so a fresh
  `init` + `doctor` is clean per LH-AK-001 §2299.
- `u-boot add <service>` adds an integrated service add-on
  (LH-FA-ADD-001..002, LH-FA-ADD-005). Today only `postgres` is in
  the catalogue; Keycloak (LH-FA-ADD-003) and OpenTelemetry
  (LH-FA-ADD-004) land in V1. Idempotent state machine: register,
  reactivate, rebuild block, repair stale artefacts, abort on
  inconsistencies.
- `u-boot up [--timeout <sec>]` starts the Compose environment via
  `docker compose up -d` and polls until every declared service
  reaches `healthy` (with healthcheck) or `running` (without)
  (LH-FA-UP-001..003). `--timeout 0` is fire-and-forget (§970, no
  `compose ps` follow-up). TCP ports declared in `compose.yaml` are
  probed on `localhost`; healthcheck-services treat a failed probe
  as warn, healthcheck-less services as a stabilization veto
  (§968 + slice §141). Exit codes per LH-FA-CLI-006: 11 Docker
  unavailable (pre-flight), 12 Compose runtime failure or
  stabilization timeout, 10 missing `u-boot.yaml` / `compose.yaml`.
- `u-boot down [--volumes]` stops the Compose environment
  (LH-FA-UP-004). `--volumes` also removes named volumes (§1015
  destructive); the LH-FA-CLI-005A §254 confirmation gate refuses
  non-interactive destructive ops without `--yes` (exit 10) and
  prompts with safe default-`N` otherwise.
- `u-boot generate <artifact>` produces or refreshes one of four
  artefacts (LH-FA-GEN-001..005): `changelog`, `readme`,
  `env-example`, `devcontainer`. Idempotent block-replace via the
  `U-BOOT MANAGED BLOCK: init` marker — user content outside the
  block survives byte-identically. `changelog` carries the
  LH-AK-007 pin (existing entries are never destroyed; a missing
  `## [Unreleased]` header is added before the first release
  section). Unknown artefacts exit 2 (spec-mandated); managed-block
  conflicts exit 10; FS errors exit 14.
- `u-boot config [get <path> | set <path> <value>]`
  (LH-FA-CONF-001..005). Without arguments prints the full
  `u-boot.yaml` byte-identically. `get` returns the bare scalar at
  one of three whitelisted paths (`project.name`,
  `devcontainer.enabled`, `services.<svc>.enabled`); `set` writes
  to the first two, with a two-stage schema-roundtrip (struct
  unmarshal + per-path domain re-validation) that aborts before
  any WriteFile when validation fails. `services.<svc>.enabled` is
  Get-only — toggling goes through `u-boot add` / `remove` to keep
  the LH-FA-ADD-005 state machine atomic.

| Phase | Status | Source |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architecture decisions | 10 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementation | M1–M8 ✅, MVP-Closure ✅ — **MVP vollständig; v0.1.0 released 2026-05-31** | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 1 temporär (LH-OPEN-002-Restwege mit benannten Trigger-Slices in ADR-0007), 8 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

## Quickstart

### Pull from GHCR (recommended)

```bash
docker pull ghcr.io/pt9912/u-boot:0.1.0    # pinned tag, recommended
# or
docker pull ghcr.io/pt9912/u-boot:latest   # stable-floating
```

Verify:

```bash
docker run --rm ghcr.io/pt9912/u-boot:0.1.0 --version
# → u-boot version 0.1.0
```

`u-boot init` against a host directory (the distroless image runs as
non-root UID 65532; `--user` matches the host UID so written files
are owned by you):

```bash
mkdir /tmp/demo && \
  docker run --rm --user "$(id -u):$(id -g)" \
    -v /tmp/demo:/work -w /work \
    ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git
```

Result: `u-boot.yaml` (`schemaVersion: 1`), `compose.yaml`, `README.md`,
`CHANGELOG.md`, `.env.example`, `.gitignore`, plus `docker/`, `scripts/`,
`docs/` directories.

Re-init on an existing project (LH-FA-INIT-005) requires an explicit
strategy — no silent overwrites:

```bash
# default: refuse to touch existing files
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git
# → exit 10: "project already initialized"

# refresh only the U-BOOT MANAGED BLOCK regions, keep user content
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git --force

# full overwrite with safety backup to <file>.bak[.N]
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git --force --backup
```

### `u-boot doctor` and the container caveat

`doctor` is designed for the **host-installed** form of u-boot — it
checks `docker`, `docker compose` and `git` on `$PATH`. The
`v0.1.0` distroless container ships none of these binaries (per
[ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md)), so
`docker run … ghcr.io/pt9912/u-boot:0.1.0 doctor` will mis-report
the host's tooling as missing. A v0.1.1-follow-up
([`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/open/slice-v0.1.1-doctor-container-awareness.md))
adds container-detection + skip semantics; medium-term a binary
distribution
([`slice-v2-binary-distribution`](docs/plan/planning/open/slice-v2-binary-distribution.md),
ADR-0007 §Folgepunkte 1 trigger now active) provides a host-native
install path. Until then, run `doctor` from a host install, or use
the `init`/`add`/`up`/`down`/`generate`/`config` subcommands —
those work via volume-mount in the container.

### Build from source (developer path)

The build is **Docker-only** (`LH-FA-BUILD-007`): no Go toolchain on the
host is required. Only Docker and `make` need to be installed.

```bash
make help                       # list all targets
make build                      # build runtime image (distroless), default VERSION=0.1.0-dev
make build VERSION=0.1.0        # build with a pinned version label
make run                        # smoke test: docker run u-boot --help
make image-scan                 # local Trivy scan (parity with CI image-scan job)
```

Inner-loop quality gates (`LH-FA-BUILD-005` / `-006`):

```bash
make lint            # golangci-lint
make test            # go test ./...
make coverage-gate   # coverage gate (bootstrap-aware, LH-FA-BUILD-008)
make gates           # lint + test + coverage-gate
make ci              # gates + govulncheck + image-scan
make fullbuild       # ci + build (full closure)
```

## Repository layout

```text
.
├── cmd/uboot/          # CLI entry point (`main.go`) — wiring layer
├── internal/           # hexagonal layout (see spec/architecture.md)
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft + architecture spec
├── docs/               # ADRs, planning, user docs (LH-FA-PROJDOCS-001)
├── scripts/            # build helpers (coverage-gate.sh)
├── Dockerfile          # multi-stage build (LH-FA-BUILD-001)
├── Makefile            # docker-only workflow (LH-FA-BUILD-005)
├── .dockerignore       # build context filter (LH-FA-BUILD-004)
└── go.mod
```

Full layout contract: [`LH-FA-BUILD-009` in
`spec/lastenheft.md`](spec/lastenheft.md).

## Documentation

- **Lastenheft** (verbindliche Spezifikation, Deutsch):
  [`spec/lastenheft.md`](spec/lastenheft.md).
- **Architecture specification:** [`spec/architecture.md`](spec/architecture.md)
  (hexagonal pattern, layer rules, depguard enforcement).
- **Quality gates:** [`docs/user/quality.md`](docs/user/quality.md)
  (SOLID-near lint profile §1.2, carveouts §1.3, tests §2,
  coverage §3, security §4).
- **Branch protection:** [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
  (LH-QA-003 PR-blocking-checks setup, one-time UI activation).
- **Architecture Decision Records:**
  [`docs/plan/adr/`](docs/plan/adr/).
- **Planning artefacts (slices, tranches):**
  [`docs/plan/planning/{open,next,in-progress,done}/`](docs/plan/planning/).
- **User documentation:** [`docs/user/`](docs/user/).

## Prerequisites

For consumers of `u-boot` (`LH-FA-DIAG-002`):

- Docker Engine ≥ 24.0.0
- Docker Compose ≥ 2.20.0
- Git
- optional: VS Code with the Dev Containers extension

For building from source (`LH-FA-BUILD-007`):

- Docker Engine
- GNU `make` (the single carveout to `LH-NFA-PORT-002` —
  see [`spec/lastenheft.md`](spec/lastenheft.md) for the rationale)

## License

MIT — see [`LICENSE`](LICENSE).
