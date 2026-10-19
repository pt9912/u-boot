# u-boot

[English](README.md) | **Deutsch**

`u-boot` ist ein CLI-Tool, das reproduzierbare Docker-basierte
Entwicklungsumgebungen aufsetzt — Projektstruktur, Docker-Compose-Stack,
Devcontainer-Konfiguration, Service-Add-Ons (PostgreSQL, Keycloak,
OpenTelemetry, …) und wiederkehrende Artefakte (README, CHANGELOG,
`.env.example`).

> **Stand:** aktuelle Version `v0.12.0` (GHCR-Image, sechs Plattform-Binaries,
> `.deb`-/`.rpm`-Pakete, Homebrew-Tap). Der Devcontainer-Katalog hat 14 Features, darunter `sshd`, `claude-code` und `ollama` (nur CPU, nur im Sandbox-Profil; Cloud-Modelle von Ollama sind nicht enthalten; gemessen auf Colima, arm64, und in zwei Läufen auf einem Linux-x86_64-Host, nicht in der CI).
> `u-boot egress list` zeigt die Egress-Allowlist eines Projekts mit der Herkunft jedes Hosts; eine geänderte Liste wirkt nach einem Neuaufbau des Containers.
> `podman build` läuft im Sandbox-Devcontainer, der Podman 6.1.4 enthält (gemessen auf Colima, arm64; auf amd64 scheitert ein `podman build` hinter der Egress-Beschränkung, weil ein verschachtelter Container Namen außerhalb der Allowlist weiter auflöst, siehe Handbuch; Linux ohne VM, Docker Desktop und CI nicht gemessen); bestehende
> `nestedRuntime: podman`-Projekte brauchen `/dev/fuse` und `/dev/net/tun` auf dem Container-Host: Eine Container-Engine ohne sie startet den Container nicht, und `u-boot doctor` meldet ein fehlendes `/dev/fuse` auf einem Linux-Host außerhalb eines Containers als Fehler (Exit-Code `11`). `nestedRuntime: docker` wird unter Docker auf Linux und auf Colima zugesagt, nicht auf Docker Desktop.
> `u-boot doctor` und `u-boot generate devcontainer` weisen aus, wenn ein Sandbox-Container auf einem Rückfall läuft (kein Docker-Daemon, `vfs`-Speicher, keine Egress-Beschränkung). Die `--dry-run`-Zusammenfassungen sagen, was geschehen würde („Would …“). `--quiet` zeigt nur Fehler und vier Ausweisungen, die sichtbar bleiben. Verlauf: Release-Tabelle unten und [`CHANGELOG.md`](CHANGELOG.md).

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
VERSION=0.12.0   # gewünschtes Release
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
ergänzt (14 Built-in-Features: `git`, `docker-cli`,
`node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, `python`, `sdkman`, `sshd`, `homebrew` und `claude-code` (beide Drittanbieter) sowie `ollama` (nur CPU, nur im Sandbox-Profil)):

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
wenn ein `source:`-Override nicht in der Allowlist steht) und
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

## Status

| Release | Datum | Highlights |
| ------- | ----- | ---------- |
| `v0.1.0` | 2026-05-31 | MVP vollständig — sieben Subkommandos (`init`, `doctor`, `add`, `up`, `down`, `generate`, `config`), alle MVP-prioritären Anforderungen geliefert. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.1.0). |
| `v0.2.0` | 2026-06-01 | Container-aware `doctor`, Six-Plattform-Binary-Distribution, `template list` + `init --template basic`. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.2.0). |
| `v0.3.0` | 2026-06-01 | Milestone „Add-on Catalogue Expansion" — `u-boot add keycloak`, `add otel`, `add <service> --with-deps`, `remove <service> [--purge]` plus Doku-Audit. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.3.0). |
| `v0.6.0` | 2026-10-01 | Milestone „Sandbox-Devcontainer + CLI-Cleanup" — `--sandbox`-Devcontainer-Profil für autonome Agenten (Named-Volume-Workspace je Instanz, optional rootless Podman, DNS-gesteuerte Egress-Allowlist, Git-Zugangsdaten nur zur Laufzeit), atomares Multi-Pfad-`config get/set` + `config list`, ausgebautes `logs`/`up`/`down`/`remove --purge`, rollback-aware `generate devcontainer`, Homebrew-Formel, Go 1.27.1. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.6.0). |
| `v0.7.0` | 2026-10-01 | Milestone „Linux-Pakete“ — `.deb` und `.rpm` (amd64, arm64) als Release-Assets per `nfpm`; Keycloak-Acceptance-Test zurück in der Standard-Integrations-Lane. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.7.0). |
| `v0.8.0` | 2026-10-14 | Milestone „Docker im Sandbox-Devcontainer“ — `devcontainer.sandbox.nestedRuntime: docker` gibt dem Sandbox-Devcontainer einen eigenen Docker-Daemon mit BuildKit, Buildx und Compose; verlangt `--privileged`, als Lockerung ausgewiesen. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.8.0). |
| `v0.8.1` | 2026-10-14 | Patch „Aktueller Docker im Sandbox-Devcontainer“ — mit `nestedRuntime: docker` läuft im Container Docker 29.9, Buildx 0.38 und Compose 5.6 aus dem Docker-APT-Repository (Schlüssel gepinnt); `download.docker.com` gehört zur Default-Allowlist. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.8.1). |
| `v0.9.0` | 2026-10-15 | Milestone „Mehr Devcontainer-Features“ — eingebaute Katalog-Schlüssel `python`, `sdkman` und `homebrew` (Community-Feature, als Drittanbieter ausgewiesen) mit ihren Egress-Hosts. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.9.0). |
| `v0.10.0` | 2026-10-17 | „sshd, Claude Code und `podman build`“ — eingebaute Katalog-Schlüssel `sshd` und `claude-code` (Drittanbieter, ausgewiesen); `nestedRuntime: podman` baut Images, mit einer sechsten Lockerung `--device=/dev/net/tun` (gemessen auf Colima, arm64; Docker Desktop, amd64 und CI nicht gemessen); `--quiet` lässt vier Ausweisungen sichtbar und unterdrückt Erfolgsmeldungen; Docker Desktop wird nicht mehr zugesagt. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.10.0). |
| `v0.11.0` | 2026-10-18 | „Podman 6, `/dev/fuse` und Ausweisung der Rückfälle“ — der Sandbox-Devcontainer mit `nestedRuntime: podman` enthält Podman 6.1.4; ein fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler, und `u-boot doctor` meldet es auf einem Linux-Host außerhalb eines Containers (Exit-Code `11`, auch ohne `--strict`); `doctor` und `generate devcontainer` weisen Rückfälle aus (kein Daemon, `vfs`, keine Egress-Beschränkung); `--dry-run`-Zusammenfassungen sagen „Would …“ (gemessen auf Colima, arm64; Linux ohne VM, Docker Desktop, amd64 und CI nicht gemessen). [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.11.0). |
| `v0.12.0` | 2026-10-19 | „Egress-Liste und Ollama“ — `u-boot egress list` zeigt die wirksame Egress-Allowlist mit der Herkunft jedes Hosts; ein Wildcard-Eintrag wird mit dem Hinweis abgewiesen, den übergeordneten Namen einzutragen, und das Handbuch nennt das gemessene Verhalten für Subdomains; eine geänderte Liste wirkt nach einem Neuaufbau des Containers; der neue Katalog-Schlüssel `ollama` installiert Ollama (nur CPU, nur im Sandbox-Profil, Cloud-Modelle nicht enthalten); `--help` verlinkt das Benutzerhandbuch (gemessen auf Colima, arm64, und `ollama` in zwei Läufen auf einem Linux-x86_64-Host; Docker Desktop, Linux hinter einer VM und CI nicht gemessen; auf amd64 scheitert ein verschachtelter `podman build` hinter der Egress-Beschränkung an der Namensauflösung, siehe Handbuch). [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.12.0). |
| `v0.5.0` | 2026-07-25 | „Sicherheit + lokale Templates" — Go-Toolchain-Bump schließt CVE-2026-39822 im Runtime-Image (HIGH); `u-boot init --template ./pfad` rendert ein Projekt aus einem lokalen Verzeichnis, neben dem eingebauten Katalog. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.5.0). |
| `v0.4.0` | 2026-06-08 | Milestone „Maschinenlesbare CLI" — `--json` / `--dry-run` / `--diff` für alle zehn Subkommandos, `u-boot logs`, Devcontainer-Features mit Drift-Doctor-Check. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.4.0). |

Die vollständige Änderungshistorie jedes Releases steht in [`CHANGELOG.md`](CHANGELOG.md).

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
| `template list [--json]` | Eingebauten Template-Katalog browsen. |

## Voraussetzungen

Für Konsumenten von `u-boot`:

- Docker Engine ≥ 24.0.0 oder Podman ≥ 5.0 (Drop-in unterstützt)
- Docker Compose ≥ 2.20.0 oder `podman compose`
- Git
- Optional: VS Code mit der Dev-Containers-Extension

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
- **Devcontainer-Features:**
  [`docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen`](docs/user/benutzerhandbuch.md#410-devcontainer-features-nutzen)
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
