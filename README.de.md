# u-boot

[English](README.md) | **Deutsch**

`u-boot` ist ein CLI-Tool, das reproduzierbare Docker-basierte
Entwicklungsumgebungen aufsetzt — Projektstruktur, Docker-Compose-Stack,
Devcontainer-Konfiguration, Service-Add-Ons (PostgreSQL, Keycloak,
OpenTelemetry, …) und wiederkehrende Artefakte (README, CHANGELOG,
`.env.example`).

> **Stand:** aktuelle Version `v0.16.1` (GHCR-Image, sechs Plattform-Binaries,
> `.deb`-/`.rpm`-Pakete, Homebrew-Tap). Zusagen, Bedingungen und Grenzen stehen im
> [Benutzerhandbuch](docs/user/benutzerhandbuch.md); was jedes Release geändert hat,
> steht in [`CHANGELOG.md`](CHANGELOG.md) und auf der
> [GitHub-Release-Seite](https://github.com/pt9912/u-boot/releases).

Das verbindliche Lastenheft
([`spec/lastenheft.md`](spec/lastenheft.md)) ist auf Deutsch verfasst;
CLI-Ausgaben und erzeugte Dateien sind auf Englisch.

## Für wen ist es?

Entwickler, Teams und Berater, die ein reproduzierbares Docker-
basiertes Projekt-Skelett brauchen, ohne pro Projekt Compose-Stacks
von Hand zu schreiben. `u-boot` erzeugt die Boilerplate
(`u-boot.yaml`, `compose.yaml`, Devcontainer-Files, …), bedient
den Add-on-Katalog (PostgreSQL, Keycloak, OpenTelemetry) und
liefert idempotente State-Machine-Operationen für Re-Init, Add,
Remove und Managed-Block-Edits.

## Was kann ich heute tun?

Nach Installation des Binarys (siehe *Installation* unten):

```bash
u-boot init my-service                  # Projekt-Skelett + git init
u-boot add postgres                     # Postgres registrieren + Compose-Block
u-boot up                               # docker compose up + Healthcheck-Poll
u-boot doctor                           # 16 Diagnose-Checks gegen Host + Projekt
u-boot down --volumes                   # Stop + Named-Volume-Cleanup (bestätigt)
u-boot remove postgres                  # Spiegel von add — disable + Blocks raus
u-boot generate readme                  # Managed-Block-Artefakt aktualisieren
u-boot config set project.name renamed-service
u-boot template list                    # Eingebauten Template-Katalog browsen
u-boot init demo --template basic       # Projekt aus eingebautem Template rendern
u-boot init demo --template ./mein-tpl  # Projekt aus lokalem Template-Verzeichnis rendern
```

Alle Subkommandos verwenden dieselben Exit-Codes
(`0` / `2` / `10` / `11` / `12` / `14`). Die *Subkommando-Referenz*
unten listet jedes Subkommando.

## Was macht es vertrauenswürdig?

- **Schriftliche Anforderungen.** Was `u-boot` zusagt, steht im [Lastenheft](spec/lastenheft.md).
- **Hexagonale Architektur.** Schicht-Regeln werden bei jedem
  `make gates` durch `depguard` enforced; die Port/Adapter-Trennung
  ist eine festgehaltene Architekturentscheidung.
- **Festgehaltene Entscheidungen.** Sprache (Go), Docker-only-Build, CI, CLI-Framework (Cobra), Distribution (GHCR + Binary), Template-Format (YAML + `text/template`), statische Plugin-Policy und der Verzicht auf einen HTTP-Adapter sind festgelegt.
- **PR-blockierende CI.** Drei PR-blockierende GitHub-Actions-Jobs
  bei jedem Push: `gates (lint + test + coverage-gate)` (läuft auch
  den Markdown-Link-Validator `docs-check`),
  `security-gates (govulncheck)` und
  `image-scan (trivy HIGH+CRITICAL)`.
- **Docker-only Inner-Loop.** `make build` baut das Runtime-Image
  ohne Go-Toolchain am Host; `make gates` läuft Lint + Test +
  Coverage im selben pinned Image-Stack den CI verwendet.

## Installation

### Vorgefertigtes Binary (empfohlen)

Statisch gelinkte Single-File-Binaries werden ab `v0.2.0` mit jedem
`v*`-GitHub-Release für sechs Plattformen (Linux/macOS/Windows ×
amd64/arm64) ausgeliefert. Kein Docker-Daemon nötig — das ist die
host-native Form für `doctor`, `init` und alle anderen Subkommandos.

**Linux / macOS:**

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
u-boot --version
```

**macOS mit Homebrew** (stabile Releases; Apple Silicon und Intel):

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot     # Homebrew 5 verlangt, Drittanbieter-Taps ausdrücklich zu vertrauen
brew install u-boot
u-boot --version
```

**Debian / Ubuntu / Fedora / RHEL** (`.deb` und `.rpm` aus dem Release, amd64 und arm64):

```bash
VERSION=0.16.1   # gewünschtes Release
# Debian / Ubuntu
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot_${VERSION}_amd64.deb"
sudo apt install "./u-boot_${VERSION}_amd64.deb"
# Fedora / RHEL
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot-${VERSION}-1.x86_64.rpm"
sudo dnf install "./u-boot-${VERSION}-1.x86_64.rpm"
```

(Für arm64 `_arm64.deb` bzw. `.aarch64.rpm`. Es gibt kein APT-/DNF-Repository; ein Update erfolgt durch Installation des neueren Pakets.)

**Windows (PowerShell):**

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Eine bestimmte Version pinnst du mit
`releases/download/v0.3.0/u-boot-<os>-<arch>[.exe]` statt
`latest/download/`.

### Pull von GHCR (alternativ für Container-/CI-Workflows)

```bash
docker pull ghcr.io/pt9912/u-boot:0.3.0    # gepinntes Tag
docker pull ghcr.io/pt9912/u-boot:latest   # stabiler Floating-Tag
docker run --rm ghcr.io/pt9912/u-boot:0.3.0 --version
```

Das Distroless-Image läuft als non-root UID 65532; mountet euer
Projekt mit `--user "$(id -u):$(id -g)"`, damit erzeugte Dateien
euch gehören. `doctor` läuft ab v0.2.0 im container-aware Modus:
die vier Host-Prerequisite-Checks werden mit `SeverityInfo`
geskipped statt als False-Positives zu feuern.

## Quickstart

```bash
mkdir my-service && cd my-service
u-boot init my-service --no-git    # --no-git in einem bestehenden Repo
u-boot add postgres
u-boot up
```

Ergebnis: `u-boot.yaml`, `compose.yaml`, `README.md`, `CHANGELOG.md`,
`.env.example`, `.gitignore` sowie die Verzeichnisse `docker/`,
`scripts/`, `docs/` — plus ein gesunder Postgres-Container auf dem
deklarierten Port.

Eine Dev-Toolchain wird über den Devcontainer-Features-Katalog
ergänzt (17 Built-in-Features: `git`, `docker-cli`,
`node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, `python`, `sdkman`, `sshd`, `homebrew` und `claude-code` (beide Drittanbieter), `gh` und `glab` (GitHub- und GitLab-CLI), `ollama` (nur CPU, nur im Sandbox-Profil) sowie `ollama-host` (eine `runArgs`-Zeile für einen Dienst des Hosts, kein Relay, beide Profile)):

```bash
u-boot init my-service --devcontainer
u-boot config set devcontainer.features.node.enabled true
u-boot generate devcontainer
# → .devcontainer/devcontainer.json enthält
#   "ghcr.io/devcontainers/features/node:1": {}
```

Externe Feature-Quellen brauchen einen expliziten Allowlist-
Eintrag; siehe
[`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen)
für den `--allow-external-feature-sources`-Fluss; fremde Installationsskripte
laufen nie unbemerkt (`--yes` reicht nicht).

`u-boot doctor` bringt zwei Checks zur Feature-
Konfiguration mit: `devcontainer.features.allowlist` (Error,
wenn ein `source:`-Override nicht in der Allowlist steht oder wenn im
Sandbox-Profil das Feature `docker-cli` aktiv ist) und
`devcontainer.features.drift` (Warn, wenn `u-boot.yaml` und die
gerenderte `devcontainer.json`-Features-Map auseinanderlaufen —
Repair via `u-boot generate devcontainer`).

Für autonome Agenten rendert `u-boot init --devcontainer --sandbox` (bzw.
`generate devcontainer --sandbox`) einen Sandbox-Devcontainer: Workspace in
benanntem Volume (Clone aus `origin`), kein Bind-/Socket-Mount, optional
rootless Podman; siehe
[`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten).

Re-Init auf einem bestehenden Projekt verlangt eine explizite
Strategie (`--force` für Managed-Block-Edits, `--backup` für
Vollüberschreibung mit `.bak[.N]`-Sicherheitskopien).

---

## Subkommando-Referenz

| Subkommando | Kurz |
| ----------- | ---- |
| `init [name] [--devcontainer] [--template <name\|pfad>]` | Projekt-Skelett + `git init`. `--template` nimmt einen Katalog-Namen (`basic`) oder einen lokalen Verzeichnis-Pfad (`./mein-tpl`, `~/tpl`). |
| `doctor [--strict]` | 16 Diagnose-Checks; container-aware Skip für Host-Probes. |
| `add <service> [--with-deps]` | Idempotente State-Machine für Service-Add-Ons (`postgres`, `keycloak`, `otel`); `--with-deps` installiert fehlende Abhängigkeiten automatisch. |
| `remove <service> [--purge]` | Spiegel von `add` — disable + Managed-Blocks raus. |
| `up [--timeout <s>]` | Compose up + Healthcheck-Poll + TCP-Probe. |
| `down [--volumes]` | Compose down mit destruktiver Bestätigungs-Gate. |
| `logs [service…] [--follow] [--tail <n>] [--no-log-prefix] [--timestamps] [--since <t>] [--until <t>]` | Compose-Logs streamen (alle Services oder eine Auswahl); `--follow` beendet sauber mit Exit 0 bei Ctrl-C. |
| `generate <artifact>` | Idempotente Block-Ersetzung via `U-BOOT MANAGED BLOCK`-Marker. |
| `config [get\|set\|list] [<pfad>…]` | Whitelist-skopierte Reads/Writes mit zweistufiger Schema-Validierung. |
| `egress list [--json]` | Egress-Allowlist des Sandbox-Devcontainers mit der Herkunft je Host anzeigen; liest nur `u-boot.yaml` und `.git/config`. |
| `devcontainer up [--recreate]` | Devcontainer des Projekts ohne Editor, Node.js und Dev-Container-CLI auf dem Host starten (nur Docker): die CLI läuft in einem kurzlebigen Hilfs-Container, der das Engine-Socket hält; der Befehl sagt das vor dem Start, auch unter `--quiet`. |
| `devcontainer down` | Diesen Container stoppen; er entfernt weder den Container noch ein Volume. |
| `devcontainer shell` | Den Devcontainer wie `devcontainer up` starten und eine interaktive Shell darin öffnen, als der Benutzer und im Workspace, die der Start meldet, und mit den Zugangsdaten, die die erzeugte Konfiguration weitergibt; die Shell selbst läuft ohne Hilfs-Container und ohne ein Socket, das der Befehl einbindet. |
| `template list [--json]` | Eingebauten Template-Katalog browsen. |

## Voraussetzungen

Für Konsumenten von `u-boot`:

- Docker Engine ≥ 24.0.0 oder Podman ≥ 4.0 (Drop-in unterstützt; für `devcontainer up`/`down`/`shell` nicht zugesagt)
- Docker Compose ≥ 2.20.0 oder `podman compose`
- Git
- Optional, zum Starten des Devcontainers: VS Code mit der Dev-Containers-Extension, die Dev-Container-CLI oder Docker allein mit `u-boot devcontainer up` und `u-boot devcontainer shell`

Für den Bau aus den Quellen:

- Docker Engine
- GNU `make`

## Repository-Layout

```text
.
├── cmd/uboot/          # CLI-Entry-Point (main.go) — Wiring-Schicht
├── internal/           # hexagonales Layout
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft (Vertrag), Spezifikation (Technik), Architektur (Sicht)
├── docs/               # ADRs, Planning, User-Doku
├── Dockerfile          # Multi-Stage-Build
├── Makefile            # generierter Aggregator (ai-harness-init)
├── repo.mk             # u-boot-Targets, Docker-only-Workflow
└── go.mod
```

## Dokumentation

- **Lastenheft** (verbindlicher Vertrag):
  [`spec/lastenheft.md`](spec/lastenheft.md)
- **Quality Gates:**
  [`docs/user/quality.md`](docs/maintainer/quality.md)
- **Branch Protection:**
  [`docs/user/branch-protection.md`](docs/maintainer/branch-protection.md)
- **Devcontainer ohne Editor (`devcontainer up`/`down`/`shell`):**
  [`docs/user/benutzerhandbuch.md#den-devcontainer-ohne-editor-starten-und-stoppen`](docs/user/benutzerhandbuch.md#den-devcontainer-ohne-editor-starten-und-stoppen)
- **Devcontainer-Features:**
  [`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen)
- **Sandbox-Profil des Devcontainers:**
  [`docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten`](docs/user/benutzerhandbuch.md#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten)
- **Maschinen-lesbarer CLI-Vertrag (`--json`, `--dry-run`, `--diff`):**
  [`docs/maintainer/cli-json-contract.md`](docs/maintainer/cli-json-contract.md)
  (Nutzersicht: [JSON-Ausgabe](docs/user/benutzerhandbuch.md#json-ausgabe))
- **User-Dokumentation:** [`docs/user/`](docs/user/)

## Build, Test, Lint

Der Build ist Docker-only; es wird keine
Go-Toolchain am Host benötigt. Nur Docker und `make` müssen
installiert sein.

```bash
make help                       # alle Targets auflisten
make build                      # Runtime-Image bauen (Distroless)
make gates                      # lint + test + coverage-gate + docs-check + baseline-verify + ablage-pruefung
make ci                         # gates + govulncheck + image-scan
make fullbuild                  # ci + build (vollständiger Closure-Lauf)
```

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
