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

**`v0.1.1` in preparation** — adds container-aware `doctor`
([`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md))
and a host-native binary distribution
([`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md),
T1 + T2 + T3 shipped: `make build-binaries` for six platforms
(Linux/macOS/Windows × amd64/arm64), `publish.yml` uploading the
binaries to the GitHub Release on every `v*` tag, and the
binary-first install block in the Quickstart below). T4 (ADR-0007
update + carveouts reduction + slice closure) remains; tag push
remains a user action — see
[`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md)
§Nächste Schritte.

- `u-boot init [name] [--devcontainer] [--template <name>]` creates the
  LH-FA-INIT-003 project structure plus `u-boot.yaml`
  (LH-FA-CONF-002) and runs `git init` by default
  (LH-FA-INIT-007). `--force` / `--backup` drive the LH-FA-INIT-005
  overwrite-protection (managed-block edits vs full overwrite with
  `.bak[.N]`); `--yes` / `--no-interactive` / `--assume-existing`
  are the LH-FA-CLI-005A mode flags (the latter drives the
  LH-FA-INIT-004 soft-detection). `--devcontainer` (LH-AK-005) also
  writes `.devcontainer/devcontainer.json` + `Dockerfile` and sets
  `devcontainer.enabled: true` in `u-boot.yaml`. `--template
  <name>` (LH-FA-TPL-001) renders the project from an external
  template catalogued by `u-boot template list`; the `basic`
  bootstrap template ships byte-identical content to the
  no-flag default flow (mutex with `--devcontainer`/`--force`/
  `--backup` in v1).
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
- `u-boot template list [--json]` (LH-FA-TPL-004, first V1
  template subcommand). Enumerates the built-in project-template
  catalog with name, description, and version in a tabular
  layout; `--json` emits a structured array with the full
  LH-FA-TPL-002 metadata surface (`supportedAddOns`,
  `generatedFiles`, `requiredTools`, `variables`). Bootstrap
  catalog ships one built-in, `basic`. Further built-ins
  (`micronaut`, `sveltekit`, …) and the `u-boot init --template
  <name>` render path land in their own ADR-0009-anchored slices
  (`slice-v1-template-init`, `slice-later-local-templates`).

| Phase | Status | Source |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architecture decisions | 10 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementation | M1–M8 ✅, MVP-Closure ✅ — **MVP vollständig; v0.1.0 released 2026-05-31** | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 1 temporär (LH-OPEN-002-Restwege mit benannten Trigger-Slices in ADR-0007), 8 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

## Quickstart

### Install pre-built binary (recommended)

Statically linked single-file binaries are attached to every `v*`
GitHub Release for six platforms (Linux/macOS/Windows × amd64/arm64)
starting with **v0.1.1**. No Docker daemon required — this is the
host-native form intended for `doctor`, `init`, and the other host-
side subcommands (per
[ADR-0007 §Folgepunkte 1](docs/plan/adr/0007-distributionswege-ghcr.md),
trigger active via
[`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md)).

**Linux / macOS** (`<os>-<arch>` auto-detected from `uname`):

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
u-boot --version
```

**Windows** (PowerShell — pick `amd64` or `arm64`):

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Pin a specific version with
`https://github.com/pt9912/u-boot/releases/download/v0.1.1/u-boot-<os>-<arch>[.exe]`
instead of `latest/download/`. `releases/latest/download/…` resolves
to the highest stable tag — `v0.1.0` predates binary assets, so
`latest` works only once `v0.1.1` (or any later tag) has been pushed.

### Pull from GHCR (alternative — container/CI workflows)

```bash
docker pull ghcr.io/pt9912/u-boot:0.1.0    # pinned tag
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
distroless image (`v0.1.0` and later) ships none of these binaries
(per [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md)), so
those host probes cannot run from a `docker run …` invocation.

Starting with **`v0.1.1`**, `doctor` detects container runtime via
`/.dockerenv` or `/run/.containerenv` and emits a `SeverityInfo`
"skipped — running inside container" diagnostic for the four
host-prerequisite checks instead of mis-reporting them as errors.
Exit code on an otherwise-clean project is `0` (not `11`). See
[`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md)
for the design rationale.

For real host-side diagnostics, run `doctor` from a host install
once the binary distribution lands
([`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md),
ADR-0007 §Folgepunkte 1 trigger now active). The other subcommands
(`init`/`add`/`up`/`down`/`generate`/`config`) work fine via
volume-mount in the container today.

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

- Docker Engine ≥ 24.0.0 **or** Podman ≥ 5.0 (drop-in via
  `DOCKER_HOST=unix:///run/user/$UID/podman/podman.sock` and a
  `docker → podman` symlink — see *Podman drop-in* below).
- Docker Compose ≥ 2.20.0 **or** `podman compose` (the
  containers/podman-compose plugin shipped with Podman 5.x).
- Git
- optional: VS Code with the Dev Containers extension

For building from source (`LH-FA-BUILD-007`):

- Docker Engine (Podman works as a drop-in but is not exercised
  in CI today — see ”Podman drop-in“ for the caveats)
- GNU `make` (the single carveout to `LH-NFA-PORT-002` —
  see [`spec/lastenheft.md`](spec/lastenheft.md) for the rationale)

### Podman drop-in

u-boot is not Podman-aware at the code level — `DockerProbe`
shells out to a `docker` binary and parses Docker version
strings. Podman works as a drop-in because:

1. `podman` exposes the same CLI surface u-boot needs
   (`info`, `version`, `compose up/down/ps`, `build`, `push/pull`).
2. The v0.1.1 container-detection (`slice-v0.1.1-doctor-container-
   awareness`) already probes `/run/.containerenv` for Podman in
   addition to `/.dockerenv` for Docker.
3. Podman ≥ 4.0 ships a Docker-API-compatible socket; pointing
   `DOCKER_HOST` at it lets every `docker`-CLI consumer talk to
   Podman.

Setup (typical Linux user):

```bash
# Start the rootless Podman API socket.
systemctl --user enable --now podman.socket
export DOCKER_HOST=unix:///run/user/$(id -u)/podman/podman.sock

# Optional: docker→podman symlink for tools that exec("docker").
sudo ln -sf "$(command -v podman)" /usr/local/bin/docker
```

Known caveats:

- `doctor` checks `docker version` against the
  `LH-FA-DIAG-002` minimums (24.0 / 2.20). Podman's version
  string is parseable but is **its own** version (e.g.
  `5.3.1`), which today classifies as
  `Severity: warn — unrecognized version` rather than `ok`.
  Functionally `up`/`down`/`add` still work.
- No CI matrix exercises the Podman path; bug reports against
  Podman are welcome but blocking-priority is Docker. A formal
  Podman-support slice will land when there is a concrete
  request — see also the v0.1.1 + ADR-0007 §Folgepunkte
  trigger pattern.

## License

MIT — see [`LICENSE`](LICENSE).
