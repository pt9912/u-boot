# u-boot

**English** | [Deutsch](README.de.md)

`u-boot` is a CLI that bootstraps reproducible Docker-based development
environments — project structure, Docker Compose stack, devcontainer
configuration, service add-ons (PostgreSQL, Keycloak, OpenTelemetry, …),
and the usual recurring artefacts (README, CHANGELOG, `.env.example`).

> **Status:** current version `v0.16.1` (GHCR image, six-platform binaries,
> `.deb`/`.rpm` packages, Homebrew tap). Promises, conditions and limits are in the
> [user manual](docs/user/benutzerhandbuch.md); what each release changed is in
> [`CHANGELOG.md`](CHANGELOG.md) and on the
> [GitHub releases page](https://github.com/pt9912/u-boot/releases).

The normative requirements ([`spec/lastenheft.md`](spec/lastenheft.md))
are written in German; CLI output and generated files are English.

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
u-boot doctor                           # 16 diagnostic checks against host + project
u-boot down --volumes                   # stop + named-volume cleanup (confirmed)
u-boot remove postgres                  # mirror of add — disable + cut blocks
u-boot generate readme                  # refresh a managed-block artefact
u-boot config set project.name renamed-service
u-boot template list                    # browse the built-in template catalogue
u-boot init demo --template basic       # render a project from a built-in template
u-boot init demo --template ./my-tpl    # render from a local template directory
```

All subcommands use the same exit codes
(`0` / `2` / `10` / `11` / `12` / `14`). The *Subcommand reference*
table below lists each subcommand. End-to-end
tasks (Postgres stack, Keycloak+OTel, devcontainer, templates, CI/JSON)
are described in the [user manual](docs/user/benutzerhandbuch.md) (German).

## What makes it trustworthy?

- **Written requirements.** What `u-boot` promises is specified in the [Lastenheft](spec/lastenheft.md) (German).
- **Hexagonal architecture.** Layer rules enforced by `depguard` at
  every `make gates`; the ports/adapters split is a recorded
  architecture decision.
- **Recorded decisions.** Language (Go), Docker-only build, CI, CLI framework (Cobra), distribution (GHCR + binary), template format (YAML + `text/template`), static plugin policy and the absence of an HTTP adapter are fixed decisions.
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

**macOS with Homebrew** (stable releases; Apple Silicon and Intel):

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot     # Homebrew 5 refuses third-party taps until they are trusted
brew install u-boot
u-boot --version
```

**Debian / Ubuntu / Fedora / RHEL** (`.deb` and `.rpm` from the release, amd64 and arm64):

```bash
VERSION=0.16.1   # the release you want
# Debian / Ubuntu
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot_${VERSION}_amd64.deb"
sudo apt install "./u-boot_${VERSION}_amd64.deb"
# Fedora / RHEL
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot-${VERSION}-1.x86_64.rpm"
sudo dnf install "./u-boot-${VERSION}-1.x86_64.rpm"
```

(For arm64 use `_arm64.deb` / `.aarch64.rpm`. There is no APT/DNF repository; update by installing the newer package.)

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
(17 built-in features: `git`, `docker-cli`, `node`,
`java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, `python`, `sdkman`, `sshd`, `homebrew` and `claude-code` (both third-party), `gh` and `glab` (GitHub and GitLab CLI), `ollama` (CPU only, sandbox profile only), and `ollama-host` (one `runArgs` line for a service on the host, no relay, both profiles)):

```bash
u-boot init my-service --devcontainer
u-boot config set devcontainer.features.node.enabled true
u-boot generate devcontainer
# → .devcontainer/devcontainer.json carries
#   "ghcr.io/devcontainers/features/node:1": {}
```

External feature sources need an explicit allowlist entry; see
[`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen)
for the `--allow-external-feature-sources` flow; foreign install scripts never run
unnoticed (`--yes` is not sufficient).

`u-boot doctor` adds two checks against the feature
configuration: `devcontainer.features.allowlist` (Error when a
`source:` override is not in the allowlist, or when the sandbox profile
is on and the feature `docker-cli` is active) and
`devcontainer.features.drift` (Warn when `u-boot.yaml` and the
rendered `devcontainer.json` features map disagree — repair via
`u-boot generate devcontainer`).

For autonomous agents, `u-boot init --devcontainer --sandbox` (or
`generate devcontainer --sandbox`) renders a sandbox devcontainer: named-volume
workspace cloned from `origin`, no host bind/socket mount, optional rootless
Podman; see
[`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten).

Re-init on an existing project requires an explicit strategy
(`--force` for managed-block edits, `--backup` for full overwrite with
`.bak[.N]` safety copies).

---

## Subcommand reference

| Subcommand | Brief |
| ---------- | ----- |
| `init [name] [--devcontainer] [--template <name\|path>]` | Scaffold project + `git init`. `--template` takes a catalogue name (`basic`) or a local directory path (`./my-tpl`, `~/tpl`). |
| `doctor [--strict]` | 16 diagnostic checks; container-aware skip for host probes. |
| `add <service> [--with-deps]` | Idempotent state-machine for service add-ons (`postgres`, `keycloak`, `otel`); `--with-deps` auto-installs missing dependencies. |
| `remove <service> [--purge]` | Mirror of `add` — disable + cut managed blocks. |
| `up [--timeout <s>]` | Compose up + healthcheck-poll + TCP probe. |
| `down [--volumes]` | Compose down with destructive-confirmation gate. |
| `logs [service…] [--follow] [--tail <n>] [--no-log-prefix] [--timestamps] [--since <t>] [--until <t>]` | Stream Compose logs (all services or a selection); `--follow` exits 0 on Ctrl-C; format flags and time range pass through to Compose. |
| `generate <artifact>` | Idempotent block-replace via `U-BOOT MANAGED BLOCK` marker. |
| `config [get\|set\|list] [<path>…]` | Whitelist-scoped reads/writes with two-stage schema validation; multi-path `get`, atomic multi-pair `set`, `list`. |
| `egress list [--json]` | Show the egress allowlist of the sandbox devcontainer with the origin of each host; reads `u-boot.yaml` and `.git/config` only. |
| `devcontainer up [--recreate]` | Start the devcontainer of the project without an editor, Node.js or the Dev Container CLI on the host (Docker only): the CLI runs in a short-lived helper container that holds the engine socket, which the command says before it starts (also under `--quiet`). |
| `devcontainer down` | Stop that container; it removes neither the container nor a volume. |
| `devcontainer shell` | Start the devcontainer like `devcontainer up` and open an interactive shell in it, as the user and in the workspace the start reports and with the credentials the generated configuration passes on; the shell itself runs without a helper container and without a socket that the command mounts. |
| `template list [--json]` | Browse the built-in template catalogue. |

## Prerequisites

For consumers of `u-boot`:

- Docker Engine ≥ 24.0.0 or Podman ≥ 4.0 (drop-in supported; not promised for `devcontainer up`/`down`/`shell`)
- Docker Compose ≥ 2.20.0 or `podman compose`
- Git
- Optional, to start the devcontainer: VS Code with the Dev Containers extension, the Dev Container CLI, or Docker alone with `u-boot devcontainer up` and `u-boot devcontainer shell`

For building from source:

- Docker Engine
- GNU `make`

## Repository layout

```text
.
├── cmd/uboot/          # CLI entry point (main.go) — wiring layer
├── internal/           # hexagonal layout
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft (contract), Spezifikation (technical), architecture (view)
├── docs/               # ADRs, planning, user docs
├── Dockerfile          # multi-stage build
├── Makefile            # generated aggregator (ai-harness-init)
├── repo.mk             # u-boot targets, docker-only workflow
└── go.mod
```

## Documentation

- **Lastenheft** (German, normative contract): [`spec/lastenheft.md`](spec/lastenheft.md)
- **Quality gates:** [`docs/user/quality.md`](docs/maintainer/quality.md)
- **Branch protection:**
  [`docs/user/branch-protection.md`](docs/maintainer/branch-protection.md)
- **Devcontainer without an editor (`devcontainer up`/`down`/`shell`):**
  [`docs/user/benutzerhandbuch.md#den-devcontainer-ohne-editor-starten-und-stoppen`](docs/user/benutzerhandbuch.md#den-devcontainer-ohne-editor-starten-und-stoppen)
- **Devcontainer features:**
  [`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen)
- **Devcontainer sandbox profile:**
  [`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten)
- **Machine-readable CLI contract (`--json`, `--dry-run`, `--diff`):**
  [`docs/maintainer/cli-json-contract.md`](docs/maintainer/cli-json-contract.md)
  (user view: [JSON output](docs/user/benutzerhandbuch.md#json-ausgabe))
- **User documentation:** [`docs/user/`](docs/user/)

## Build, Test, Lint

The build is Docker-only; no Go toolchain on the
host is required. Only Docker and `make` need to be installed.

```bash
make help                       # list all targets
make build                      # build runtime image (distroless)
make gates                      # lint + test + coverage-gate + docs-check + baseline-verify + ablage-pruefung
make ci                         # gates + govulncheck + image-scan
make fullbuild                  # ci + build (full closure)
```

## License

MIT — see [`LICENSE`](LICENSE).
