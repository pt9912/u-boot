# u-boot

[English](README.md) | **Deutsch**

`u-boot` ist ein CLI-Tool, das reproduzierbare Docker-basierte
Entwicklungsumgebungen aufsetzt — Projektstruktur, Docker-Compose-Stack,
Devcontainer-Konfiguration, Service-Add-Ons (PostgreSQL, Keycloak,
OpenTelemetry, …) und wiederkehrende Artefakte (README, CHANGELOG,
`.env.example`).

> **Stand:** `v0.2.0` released 2026-06-01 (GHCR + sechs Plattform-
> Binaries); `v0.3.0`-Milestone „Add-on Catalogue Expansion" in
> Bearbeitung (1/5 Slices). Vollständige Release-Tabelle unten.

Das verbindliche Lastenheft
([`spec/lastenheft.md`](spec/lastenheft.md)) ist auf Deutsch verfasst;
CLI-Ausgaben und erzeugte Dateien sind auf Englisch (`LH-LESE-002`).

## Für wen ist es?

Entwickler, Teams und Berater, die ein reproduzierbares Docker-
basiertes Projekt-Skelett brauchen, ohne pro Projekt Compose-Stacks
von Hand zu schreiben. `u-boot` erzeugt die Boilerplate
(`u-boot.yaml`, `compose.yaml`, Devcontainer-Files, …), bedient
den Add-on-Katalog (PostgreSQL heute; Keycloak und OpenTelemetry
folgen in v0.3.0) und liefert idempotente State-Machine-Operationen
für Re-Init, Add, Remove und Managed-Block-Edits.

## Was kann ich heute tun?

Nach Installation des Binarys (siehe *Installation* unten):

```bash
u-boot init my-service                  # Projekt-Skelett + git init
u-boot add postgres                     # Postgres registrieren + Compose-Block
u-boot up                               # docker compose up + Healthcheck-Poll
u-boot doctor                           # 11 Diagnose-Checks gegen Host + Projekt
u-boot down --volumes                   # Stop + Named-Volume-Cleanup (bestätigt)
u-boot remove postgres                  # Spiegel von add — disable + Blocks raus
u-boot generate readme                  # Managed-Block-Artefakt aktualisieren
u-boot config set project.name renamed-service
u-boot template list                    # Eingebauten Template-Katalog browsen
u-boot init demo --template basic       # Projekt aus einem Template rendern
```

Alle Subkommandos respektieren die LH-FA-CLI-006-Exit-Codes
(`0` / `2` / `10` / `11` / `12` / `14`). Die *Subkommando-Referenz*
unten mappt jedes Subkommando auf seine Lastenheft-IDs.

## Was macht es vertrauenswürdig?

- **MVP mit v0.1.0 geschlossen.** Jede MVP-prioritäre `LH-AK-*`-,
  `LH-FA-*`- und `LH-SA-*`-ID aus
  [`spec/lastenheft.md`](spec/lastenheft.md) ist ausgeliefert —
  Audit-Trail in der
  [roadmap §MVP-Bilanz](docs/plan/planning/in-progress/roadmap.md).
- **Hexagonale Architektur.** Schicht-Regeln werden bei jedem
  `make gates` durch `depguard` enforced; Port/Adapter-Trennung
  formalisiert in
  [`ADR-0002`](docs/plan/adr/0002-hexagonale-architektur.md).
- **ADR-getriebene Entscheidungen.** 10 Architecture Decision
  Records decken Sprache (Go), Build (Docker-only), CI, CLI-
  Framework (Cobra), Distribution (GHCR + Binary), Template-Format
  (YAML + `text/template`), Plugin-Policy (statisch) und die
  „kein HTTP-Adapter"-Entscheidung ab.
- **PR-blockierende CI.** Drei PR-blockierende GitHub-Actions-Jobs
  (`gates (lint + test + coverage-gate)`,
  `security-gates (govulncheck)`,
  `image-scan (trivy HIGH+CRITICAL)`) plus ein Markdown-Link-
  Validator bei jedem Push.
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

**Windows (PowerShell):**

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Eine bestimmte Version pinnst du mit
`releases/download/v0.2.0/u-boot-<os>-<arch>[.exe]` statt
`latest/download/`.

### Pull von GHCR (alternativ für Container-/CI-Workflows)

```bash
docker pull ghcr.io/pt9912/u-boot:0.2.0    # gepinntes Tag
docker pull ghcr.io/pt9912/u-boot:latest   # stabiler Floating-Tag
docker run --rm ghcr.io/pt9912/u-boot:0.2.0 --version
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

Re-Init auf einem bestehenden Projekt verlangt eine explizite
Strategie (`--force` für Managed-Block-Edits, `--backup` für
Vollüberschreibung mit `.bak[.N]`-Sicherheitskopien). Siehe den
[init-Slice](docs/plan/planning/done/slice-m3-init-flow.md) für die
`LH-FA-INIT-005`-State-Machine.

---

## Status

| Release | Datum | Highlights |
| ------- | ----- | ---------- |
| `v0.1.0` | 2026-05-31 | MVP vollständig — sieben Subkommandos (`init`, `doctor`, `add`, `up`, `down`, `generate`, `config`), alle MVP-prioritären Lastenheft-IDs geliefert. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.1.0). |
| `v0.2.0` | 2026-06-01 | Container-aware `doctor`, Six-Plattform-Binary-Distribution, `template list` + `init --template basic`. [GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.2.0). |
| `v0.3.0` | in Bearbeitung (1/5) | Milestone „Add-on Catalogue Expansion" — Keycloak + OpenTelemetry + Add-on-Dependency-Resolution + `remove <service>` (geliefert). Milestone-Tabelle in [`roadmap.md §v0.3.0`](docs/plan/planning/in-progress/roadmap.md). |

Die Roadmap
([`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md))
hat den vollständigen Audit-Trail: Phase-Tabelle (M1..M8 +
Closure + V1-Cluster), per-Release-Milestone-Tabellen, Carveout-
Auflösungs-Slices und §Nächste Schritte für das laufende Backlog.

## Subkommando-Referenz

| Subkommando | Spec-IDs | Kurz |
| ----------- | -------- | ---- |
| `init [name] [--devcontainer] [--template <name>]` | `LH-FA-INIT-001..007`, `LH-FA-TPL-001` | Projekt-Skelett + `git init`. |
| `doctor [--strict]` | `LH-FA-DIAG-001..004` | 11 Diagnose-Checks; container-aware Skip für Host-Probes. |
| `add <service>` | `LH-FA-ADD-001..002`, `LH-FA-ADD-005` | Idempotente State-Machine für Service-Add-Ons. |
| `remove <service> [--purge]` | `LH-FA-ADD-007` | Spiegel von `add` — disable + Managed-Blocks raus. |
| `up [--timeout <s>]` | `LH-FA-UP-001..003` | Compose up + Healthcheck-Poll + TCP-Probe. |
| `down [--volumes]` | `LH-FA-UP-004` | Compose down mit destruktiver Bestätigungs-Gate. |
| `generate <artifact>` | `LH-FA-GEN-001..005` | Idempotente Block-Ersetzung via `U-BOOT MANAGED BLOCK`-Marker. |
| `config [get\|set] [<pfad> [<wert>]]` | `LH-FA-CONF-001..005` | Whitelist-skopierte Reads/Writes mit zweistufiger Schema-Validierung. |
| `template list [--json]` | `LH-FA-TPL-004` | Eingebauten Template-Katalog browsen. |

## Voraussetzungen

Für Konsumenten von `u-boot` (`LH-FA-DIAG-002`):

- Docker Engine ≥ 24.0.0 oder Podman ≥ 5.0 (Drop-in unterstützt;
  siehe [`spec/architecture.md §2.4`](spec/architecture.md))
- Docker Compose ≥ 2.20.0 oder `podman compose`
- Git
- Optional: VS Code mit der Dev-Containers-Extension

Für den Bau aus den Quellen (`LH-FA-BUILD-007`):

- Docker Engine
- GNU `make` (der einzige permanente Carveout zu
  `LH-NFA-PORT-002`)

## Repository-Layout

```text
.
├── cmd/uboot/          # CLI-Entry-Point (main.go) — Wiring-Schicht
├── internal/           # hexagonales Layout (siehe spec/architecture.md)
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft + Architektur-Spezifikation
├── docs/               # ADRs, Planning, User-Doku (LH-FA-PROJDOCS-001)
├── Dockerfile          # Multi-Stage-Build (LH-FA-BUILD-001)
├── Makefile            # Docker-only-Workflow (LH-FA-BUILD-005)
└── go.mod
```

Vollständiger Layout-Kontrakt:
[`LH-FA-BUILD-009` in `spec/lastenheft.md`](spec/lastenheft.md).

## Dokumentation

- **Lastenheft** (verbindliche Spezifikation):
  [`spec/lastenheft.md`](spec/lastenheft.md)
- **Architektur-Spezifikation:**
  [`spec/architecture.md`](spec/architecture.md) (hexagonales
  Pattern, Schicht-Regeln, Podman-Drop-in §2.4)
- **Architecture Decision Records:**
  [`docs/plan/adr/`](docs/plan/adr/)
- **Roadmap, Slices, Carveouts:**
  [`docs/plan/planning/`](docs/plan/planning/)
- **Quality Gates:**
  [`docs/user/quality.md`](docs/user/quality.md)
- **Branch Protection:**
  [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
- **User-Dokumentation:** [`docs/user/`](docs/user/)

## Build, Test, Lint

Der Build ist Docker-only (`LH-FA-BUILD-007`); es wird keine
Go-Toolchain am Host benötigt. Nur Docker und `make` müssen
installiert sein.

```bash
make help                       # alle Targets auflisten
make build                      # Runtime-Image bauen (Distroless)
make gates                      # lint + test + coverage-gate + docs-check
make ci                         # gates + govulncheck + image-scan
make fullbuild                  # ci + build (vollständiger Closure-Lauf)
```

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
