# u-boot

**English** | [Deutsch](README.de.md)

`u-boot` is a CLI that bootstraps reproducible Docker-based development
environments — project structure, Docker Compose stack, devcontainer
configuration, service add-ons (PostgreSQL, Keycloak, OpenTelemetry, …),
and the usual recurring artefacts (README, CHANGELOG, `.env.example`).

> **Status:** `v0.3.0` released 2026-06-01 (GHCR + six-platform
> binaries). Adds `add keycloak` + `add otel` + `remove <service>` +
> `--with-deps` to the catalogue. Full release table below.

The normative requirements ([`spec/lastenheft.md`](spec/lastenheft.md))
are written in German; CLI output and generated files are English
(`LH-LESE-002`).

## Who is it for?

Developers, teams, and consultants who need a reproducible Docker-based
project skeleton without hand-rolling Compose stacks per project.
`u-boot` generates the boilerplate (`u-boot.yaml`, `compose.yaml`,
devcontainer files, …), wires the add-on catalogue (PostgreSQL,
Keycloak, OpenTelemetry), and provides idempotent state-machine
operations for re-init, add, remove, and managed-block edits.

## What can I do today?

After installing the binary (see *Install* below):

```bash
u-boot init my-service                  # scaffold project + git init
u-boot add postgres                     # register Postgres + write compose block
u-boot up                               # docker compose up + healthcheck poll
u-boot doctor                           # 11 diagnostic checks against host + project
u-boot down --volumes                   # stop + named-volume cleanup (confirmed)
u-boot remove postgres                  # mirror of add — disable + cut blocks
u-boot generate readme                  # refresh a managed-block artefact
u-boot config set project.name renamed-service
u-boot template list                    # browse the built-in template catalogue
u-boot init demo --template basic       # render a project from a template
```

All subcommands respect LH-FA-CLI-006 exit codes
(`0` / `2` / `10` / `11` / `12` / `14`). The *Subcommand reference*
table below maps each subcommand to its Lastenheft IDs.

## What makes it trustworthy?

- **Spec-driven releases.** Three tagged releases (`v0.1.0`, `v0.2.0`,
  `v0.3.0`) deliver every MVP and v0.3.0 V1-add-on Spec-ID listed in
  [`spec/lastenheft.md`](spec/lastenheft.md); the release table below
  maps each slice to its `LH-FA-*` / `LH-AK-*` anchor.
- **Hexagonal architecture.** Layer rules enforced by `depguard` at
  every `make gates`; ports/adapters split formalised in
  [`ADR-0002`](docs/plan/adr/0002-hexagonale-architektur.md).
- **ADR-driven decisions.** 10 Architecture Decision Records cover
  language (Go), build (Docker-only), CI, CLI framework (Cobra),
  distribution (GHCR + binary), template format (YAML + `text/template`),
  plugin policy (static), and the no-HTTP-adapter stance.
- **PR-blocking CI.** Three required GitHub-Actions jobs on every
  push: `gates (lint + test + coverage-gate)` (which also runs the
  Markdown-link validator `docs-check`), `security-gates (govulncheck)`,
  and `image-scan (trivy HIGH+CRITICAL)`.
- **Docker-only inner-loop.** `make build` builds the runtime image
  without any Go toolchain on the host; `make gates` runs lint + test
  + coverage in the same pinned image stack CI uses.

## Install

### Pre-built binary (recommended)

Single-file static binaries are attached to every `v*` GitHub Release
for six platforms (Linux/macOS/Windows × amd64/arm64) starting with
`v0.2.0`. No Docker daemon required — this is the host-native form for
`doctor`, `init`, and all other subcommands.

**Linux / macOS:**

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
u-boot --version
```

**Windows (PowerShell):**

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Pin a specific version with
`releases/download/v0.3.0/u-boot-<os>-<arch>[.exe]` instead of
`latest/download/`.

### Pull from GHCR (alternative for container/CI workflows)

```bash
docker pull ghcr.io/pt9912/u-boot:0.3.0    # pinned tag
docker pull ghcr.io/pt9912/u-boot:latest   # stable-floating
docker run --rm ghcr.io/pt9912/u-boot:0.3.0 --version
```

The distroless image runs as non-root UID 65532; mount your project
with `--user "$(id -u):$(id -g)"` so written files are owned by you.
`doctor` runs in container-aware mode (since v0.2.0): the four
host-prerequisite checks are skipped with `SeverityInfo` instead of
firing as false positives.

## Quickstart

```bash
mkdir my-service && cd my-service
u-boot init my-service --no-git    # use --no-git inside an existing repo
u-boot add postgres
u-boot up
```

Result: `u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`,
`.env.example`, `.gitignore`, plus `docker/`, `scripts/`, `docs/`
directories — and a healthy Postgres container ready at the declared
port.

Add a development toolchain via the devcontainer features catalogue
(`LH-FA-DEV-003`, 8 built-in features: `git`, `docker-cli`, `node`,
`java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`):

```bash
u-boot init my-service --devcontainer
u-boot config set devcontainer.features.node.enabled true
u-boot generate devcontainer
# → .devcontainer/devcontainer.json carries
#   "ghcr.io/devcontainers/features/node:1": {}
```

External feature sources need an explicit allowlist entry; see
[`docs/user/devcontainer-features.md`](docs/user/devcontainer-features.md)
for the `--allow-external-feature-sources` flow and the
`LH-NFA-SEC-004` discipline (`--yes` is not sufficient).

Re-init on an existing project requires an explicit strategy
(`--force` for managed-block edits, `--backup` for full overwrite with
`.bak[.N]` safety copies). See the
[init slice](docs/plan/planning/done/slice-m3-init-flow.md) for the
`LH-FA-INIT-005` state machine.

---

## Status

| Release | Date | Highlights |
| ------- | ---- | ---------- |
| `v0.1.0` | 2026-05-31 | MVP complete — seven subcommands (`init`, `doctor`, `add`, `up`, `down`, `generate`, `config`), all MVP-priority Lastenheft IDs delivered. [GitHub release](https://github.com/pt9912/u-boot/releases/tag/v0.1.0). |
| `v0.2.0` | 2026-06-01 | Container-aware `doctor`, six-platform binary distribution, `template list` + `init --template basic`. [GitHub release](https://github.com/pt9912/u-boot/releases/tag/v0.2.0). |
| `v0.3.0` | 2026-06-01 | "Add-on Catalogue Expansion" milestone — `u-boot add keycloak` (LH-FA-ADD-003), `add otel` (LH-FA-ADD-004), `add <service> --with-deps` (LH-FA-ADD-006), `remove <service> [--purge]` (LH-FA-ADD-007), plus a doku-audit closure for three V1 spec-IDs. [GitHub release](https://github.com/pt9912/u-boot/releases/tag/v0.3.0). |

The roadmap ([`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md))
has the full audit trail: Phase table (M1..M8 + Closure + V1
clusters), per-release milestone tables, carveout-resolution slices,
and §Nächste Schritte for the in-progress backlog.

## Subcommand reference

| Subcommand | Spec IDs | Brief |
| ---------- | -------- | ----- |
| `init [name] [--devcontainer] [--template <name>]` | `LH-FA-INIT-001..007`, `LH-FA-TPL-001` | Scaffold project + `git init`. |
| `doctor [--strict]` | `LH-FA-DIAG-001..004`, `LH-FA-DEV-003` | 12 diagnostic checks (the 12th is `devcontainer.features.allowlist` from v0.4.0); container-aware skip for host probes. |
| `add <service> [--with-deps]` | `LH-FA-ADD-001..006` | Idempotent state-machine for service add-ons (`postgres`, `keycloak`, `otel`); `--with-deps` auto-installs missing dependencies. |
| `remove <service> [--purge]` | `LH-FA-ADD-007` | Mirror of `add` — disable + cut managed blocks. |
| `up [--timeout <s>]` | `LH-FA-UP-001..003` | Compose up + healthcheck-poll + TCP probe. |
| `down [--volumes]` | `LH-FA-UP-004` | Compose down with destructive-confirmation gate. |
| `generate <artifact>` | `LH-FA-GEN-001..005` | Idempotent block-replace via `U-BOOT MANAGED BLOCK` marker. |
| `config [get\|set] [<path> [<value>]]` | `LH-FA-CONF-001..005` | Whitelist-scoped reads/writes with two-stage schema validation. |
| `template list [--json]` | `LH-FA-TPL-004` | Browse the built-in template catalogue. |

## Prerequisites

For consumers of `u-boot` (`LH-FA-DIAG-002`):

- Docker Engine ≥ 24.0.0 or Podman ≥ 5.0 (drop-in supported; see
  [`spec/architecture.md §2.4`](spec/architecture.md))
- Docker Compose ≥ 2.20.0 or `podman compose`
- Git
- Optional: VS Code with the Dev Containers extension

For building from source (`LH-FA-BUILD-007`):

- Docker Engine
- GNU `make` (single permanent carveout to `LH-NFA-PORT-002`)

## Repository layout

```text
.
├── cmd/uboot/          # CLI entry point (main.go) — wiring layer
├── internal/           # hexagonal layout (see spec/architecture.md)
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft + architecture spec
├── docs/               # ADRs, planning, user docs (LH-FA-PROJDOCS-001)
├── Dockerfile          # multi-stage build (LH-FA-BUILD-001)
├── Makefile            # docker-only workflow (LH-FA-BUILD-005)
└── go.mod
```

Full layout contract: [`LH-FA-BUILD-009` in `spec/lastenheft.md`](spec/lastenheft.md).

## Documentation

- **Lastenheft** (German, normative): [`spec/lastenheft.md`](spec/lastenheft.md)
- **Architecture specification:** [`spec/architecture.md`](spec/architecture.md)
  (hexagonal pattern, layer rules, Podman drop-in §2.4)
- **Architecture Decision Records:** [`docs/plan/adr/`](docs/plan/adr/)
- **Roadmap, slices, carveouts:**
  [`docs/plan/planning/`](docs/plan/planning/)
- **Quality gates:** [`docs/user/quality.md`](docs/user/quality.md)
- **Branch protection:**
  [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
- **Devcontainer features:**
  [`docs/user/devcontainer-features.md`](docs/user/devcontainer-features.md)
- **User documentation:** [`docs/user/`](docs/user/)

## Build, Test, Lint

The build is Docker-only (`LH-FA-BUILD-007`); no Go toolchain on the
host is required. Only Docker and `make` need to be installed.

```bash
make help                       # list all targets
make build                      # build runtime image (distroless)
make gates                      # lint + test + coverage-gate + docs-check
make ci                         # gates + govulncheck + image-scan
make fullbuild                  # ci + build (full closure)
```

## License

MIT — see [`LICENSE`](LICENSE).
