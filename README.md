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

**MVP in progress — `u-boot init` is fully wired, including re-init.**
The first functional subcommand is shipped end-to-end (M3 ✅):
`u-boot init [name]` creates the LH-FA-INIT-003 project structure plus
`u-boot.yaml` (LH-FA-CONF-002) and runs `git init` by default
(LH-FA-INIT-007); re-running on an existing project supports the
LH-FA-INIT-005 overwrite-protection (`--force` for managed-block-only
edits, `--backup` for full overwrite with `.bak[.N]` safety copies)
and the LH-FA-CLI-005A mode flags (`--yes` / `--no-interactive`
exclusive, `--assume-existing` accepted for M4 soft-detection).
Subsequent MVP subcommands (`add`, `up`, `down`, `doctor`, `generate`,
`config`) follow in M4+; planning is tracked under
[`docs/plan/planning/`](docs/plan/planning/).

| Phase | Status | Source |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architecture decisions | 5 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementation | M1–M3 ✅, M4 next | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 12 temporär (10 mit Slice-Plan, 1 Slice deckt 3), 7 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

## Quickstart

The build is **Docker-only** (`LH-FA-BUILD-007`): no Go toolchain on the
host is required. Only Docker and `make` need to be installed.

```bash
make help            # list all targets
make build           # build the runtime image (distroless static, nonroot)
make run             # smoke test: docker run u-boot --help
```

Real `u-boot init` against a host directory (distroless runs as
non-root UID 65532; `--user` matches the host UID so written files
are owned by you):

```bash
mkdir /tmp/demo && \
  docker run --rm --user "$(id -u):$(id -g)" \
    -v /tmp/demo:/work -w /work \
    u-boot:latest init demo --no-git
```

Result: `u-boot.yaml` (`schemaVersion: 1`), `compose.yaml`, `README.md`,
`CHANGELOG.md`, `.env.example`, `.gitignore`, plus `docker/`, `scripts/`,
`docs/` directories.

Re-init on an existing project (LH-FA-INIT-005) requires an explicit
strategy — no silent overwrites:

```bash
# default: refuse to touch existing files
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git
# → exit 10: "project already initialized"

# refresh only the U-BOOT MANAGED BLOCK regions, keep user content
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git --force

# full overwrite with safety backup to <file>.bak[.N]
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git --force --backup
```

Inner-loop quality gates (`LH-FA-BUILD-005` / `-006`):

```bash
make lint            # golangci-lint
make test            # go test ./...
make coverage-gate   # coverage gate (bootstrap-aware, LH-FA-BUILD-008)
make gates           # lint + test + coverage-gate
make ci              # gates + govulncheck
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
- **Architecture Decision Records:**
  [`docs/plan/adr/`](docs/plan/adr/).
- **Planning artefacts (slices, tranches):**
  [`docs/plan/planning/{open,next,in-progress,done}/`](docs/plan/planning/).
- **User documentation:** [`docs/user/`](docs/user/) (empty during the
  bootstrap phase).

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
