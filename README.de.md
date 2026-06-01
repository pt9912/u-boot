# u-boot

**Ein Bootloader für Entwicklungsumgebungen auf Docker-Basis.**

`u-boot` ist ein CLI-Tool, das reproduzierbare Entwicklungsumgebungen
aufsetzt: Projektstruktur, Docker-Compose-Stack, Devcontainer-
Konfiguration, Service-Add-ons (PostgreSQL, Keycloak, OpenTelemetry, …)
und wiederkehrende Artefakte (README, CHANGELOG, `.env.example`).

> **Language:** The English version of this README is at
> [`README.md`](README.md). The Lastenheft
> ([`spec/lastenheft.md`](spec/lastenheft.md)) is written in German;
> CLI-Ausgaben und erzeugte Dateien sind auf Englisch (`LH-LESE-002`).

## Status

**MVP vollständig — sieben Subkommandos verdrahtet (`init` + `doctor` + `add` + `up` + `down` + `generate` + `config`).**

Jeder MVP-priorisierte `LH-AK-*`-, `LH-FA-*`- und `LH-SA-*`-Eintrag
aus [`spec/lastenheft.md`](spec/lastenheft.md) ist geliefert.
**`v0.1.0` ist released (2026-05-31)** — siehe
[GitHub-Release](https://github.com/pt9912/u-boot/releases/tag/v0.1.0)
und das GHCR-Image `ghcr.io/pt9912/u-boot:0.1.0` (plus den stabilen
Floating-Tag `:latest`). Distributionsentscheidung in
[ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md). Audit-
Trail im
[MVP-Bilanz-Block der Roadmap](docs/plan/planning/in-progress/roadmap.md)
und im
[Release-Cut-Slice](docs/plan/planning/done/slice-v1-release-cut-v0.1.0.md).

**`v0.1.1` in Vorbereitung** — ergänzt einen container-aware
`doctor`
([`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md))
und eine host-native Binary-Distribution
([`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md),
T1 + T2 + T3 geliefert: `make build-binaries` für sechs Plattformen
(Linux/macOS/Windows × amd64/arm64), `publish.yml` lädt die Binaries
bei jedem `v*`-Tag an den GitHub-Release hoch, und der Binary-First-
Install-Block in der Quickstart unten). T4 (ADR-0007-Update +
carveouts-Reduktion + Slice-Closure) bleibt offen; der Tag-Push
bleibt Nutzer-Aktion — siehe
[`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md)
§Nächste Schritte.

- `u-boot init [name] [--devcontainer]` erzeugt die
  LH-FA-INIT-003-Projektstruktur plus `u-boot.yaml`
  (LH-FA-CONF-002) und initialisiert per Default ein Git-Repository
  (LH-FA-INIT-007). `--force` / `--backup` treiben den
  LH-FA-INIT-005-Überschreibschutz (Managed-Block-Edits vs
  Vollüberschreibung mit `.bak[.N]`-Sicherung); `--yes` /
  `--no-interactive` / `--assume-existing` sind die
  LH-FA-CLI-005A-Modi-Flags (letzteres treibt die
  LH-FA-INIT-004-Soft-Detection). `--devcontainer` (LH-AK-005)
  schreibt zusätzlich `.devcontainer/devcontainer.json` +
  `Dockerfile` und setzt `devcontainer.enabled: true` in
  `u-boot.yaml`.
- `u-boot doctor` führt 11 Diagnose-Checks gegen die lokale
  Umgebung und das Projekt aus (LH-FA-DIAG-002), klassifiziert
  Befunde als ok / warn / error (LH-FA-DIAG-003), gibt
  Reparaturhinweise (LH-FA-DIAG-004) und exited mit 11 bei Errors
  (oder Warns mit `--strict`). M5 ergänzt `services.enabled-key`,
  `devcontainer.forwardPorts.consistency` und die Severity-
  Eskalation über `devcontainer.enabled`. MVP-Closure-T2 ändert
  `compose.yaml.valid` no-services von Error auf Warn, damit ein
  frisches `init` + `doctor` LH-AK-001 §2299 erfüllt.
- `u-boot add <service>` fügt ein integriertes Service-Add-On in
  das aktuelle Projekt ein (LH-FA-ADD-001..002, LH-FA-ADD-005).
  Heute nur `postgres` im Katalog; Keycloak (LH-FA-ADD-003) und
  OpenTelemetry (LH-FA-ADD-004) folgen in V1. Idempotent, mit
  voller State-Machine: registrieren, reaktivieren, Block neu
  erzeugen, fehlende Artefakte reparieren, Abbruch bei
  Inkonsistenzen.
- `u-boot up [--timeout <sek>]` startet die Compose-Umgebung via
  `docker compose up -d` und pollt, bis jeder deklarierte Service
  `healthy` (mit Healthcheck) bzw. `running` (ohne) erreicht
  (LH-FA-UP-001..003). `--timeout 0` ist Fire-and-Forget (§970,
  kein `compose ps`-Follow-up). In `compose.yaml` deklarierte
  TCP-Ports werden auf `localhost` geprüft; mit Healthcheck
  emittiert ein Probe-Fehler nur eine Warn-Diagnose ohne Veto,
  ohne Healthcheck veto'd er die Stabilisierung (§968 +
  Slice §141). Exit-Codes per LH-FA-CLI-006: 11 wenn Docker nicht
  verfügbar (Pre-Flight), 12 bei Compose-Laufzeitfehler oder
  Stabilisierungs-Timeout, 10 bei fehlendem `u-boot.yaml` /
  `compose.yaml`.
- `u-boot down [--volumes]` stoppt die Compose-Umgebung
  (LH-FA-UP-004). `--volumes` entfernt zusätzlich Named-Volumes
  (§1015 destruktiv); der LH-FA-CLI-005A-§254-Bestätigungs-Gate
  bricht nicht-interaktive destruktive Ops ohne `--yes` mit
  Exit 10 ab und prompt't sonst mit sicherem Default-`N`.
- `u-boot generate <artifact>` erzeugt oder aktualisiert eines
  von vier Artefakten (LH-FA-GEN-001..005): `changelog`,
  `readme`, `env-example`, `devcontainer`. Idempotenter
  Block-Replace via den `U-BOOT MANAGED BLOCK: init`-Marker —
  User-Inhalt außerhalb des Blocks bleibt byte-identisch.
  `changelog` trägt den LH-AK-007-Pin (existierende Einträge
  werden nie zerstört; fehlender `## [Unreleased]`-Header wird
  vor der ersten Versions-Sektion ergänzt). Unbekannte Artefakte
  exiten mit 2 (Spec-Pflicht); managed-Block-Konflikte mit 10;
  FS-Fehler mit 14.
- `u-boot config [get <pfad> | set <pfad> <wert>]`
  (LH-FA-CONF-001..005). Ohne Argumente zeigt es die
  `u-boot.yaml` byte-identisch. `get` liefert den nackten Skalar
  an einem von drei whitelist-Pfaden (`project.name`,
  `devcontainer.enabled`, `services.<svc>.enabled`); `set`
  schreibt auf die ersten zwei, mit zweistufiger Schema-
  Roundtrip-Validation (Struct-Unmarshal + per-Pfad-Domain-
  Re-Validation), die vor jedem WriteFile bei
  Validierungsfehler abbricht. `services.<svc>.enabled` ist
  Get-only — das Toggeln geht über `u-boot add` / `remove`,
  damit die LH-FA-ADD-005-State-Machine atomar bleibt.
- `u-boot template list [--json]` (LH-FA-TPL-004, erstes V1-
  Template-Subkommando). Listet den eingebauten Projekt-Template-
  Katalog mit Name, Beschreibung und Version in tabellarischer
  Form; `--json` gibt ein strukturiertes Array mit der vollen
  LH-FA-TPL-002-Metadaten-Oberfläche (`supportedAddOns`,
  `generatedFiles`, `requiredTools`, `variables`) aus. Bootstrap-
  Katalog liefert ein Built-in: `basic`. Weitere Built-ins
  (`micronaut`, `sveltekit`, …) und der `u-boot init --template
  <name>`-Render-Pfad landen in eigenen ADR-0009-verankerten
  Slices (`slice-v1-template-init`, `slice-later-local-templates`).

| Phase | Status | Quelle |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architekturentscheidungen | 10 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementierung | M1–M8 ✅, MVP-Closure ✅ — **MVP vollständig; v0.1.0 released 2026-05-31** | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 1 temporär (LH-OPEN-002-Restwege mit benannten Trigger-Slices in ADR-0007), 8 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

## Quickstart

### Vorgefertigte Binary installieren (empfohlen)

Statisch gelinkte Single-File-Binaries werden ab **v0.1.1** mit jedem
`v*`-GitHub-Release für sechs Plattformen (Linux/macOS/Windows ×
amd64/arm64) ausgeliefert. Kein Docker-Daemon nötig — das ist die
host-native Form für `doctor`, `init` und die anderen host-seitigen
Subkommandos (gemäß
[ADR-0007 §Folgepunkte 1](docs/plan/adr/0007-distributionswege-ghcr.md),
Trigger aktiv via
[`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md)).

**Linux / macOS** (`<os>-<arch>` werden aus `uname` ermittelt):

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
u-boot --version
```

**Windows** (PowerShell — `amd64` oder `arm64` wählen):

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Eine bestimmte Version pinnst du mit
`https://github.com/pt9912/u-boot/releases/download/v0.1.1/u-boot-<os>-<arch>[.exe]`
statt `latest/download/`. `releases/latest/download/…` zeigt immer
auf den höchsten stabilen Tag — `v0.1.0` hatte noch keine Binary-
Assets, also funktioniert `latest` erst, sobald `v0.1.1` (oder ein
späterer Tag) gepusht ist.

### Pull von GHCR (alternativ — Container-/CI-Workflows)

```bash
docker pull ghcr.io/pt9912/u-boot:0.1.0    # gepinntes Tag
# oder
docker pull ghcr.io/pt9912/u-boot:latest   # stabiler Floating-Tag
```

Verifikation:

```bash
docker run --rm ghcr.io/pt9912/u-boot:0.1.0 --version
# → u-boot version 0.1.0
```

`u-boot init` gegen ein Host-Verzeichnis (Distroless läuft als non-
root UID 65532; `--user` matched die Host-UID, damit erzeugte Dateien
dir gehören):

```bash
mkdir /tmp/demo && \
  docker run --rm --user "$(id -u):$(id -g)" \
    -v /tmp/demo:/work -w /work \
    ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git
```

Ergebnis: `u-boot.yaml` (`schemaVersion: 1`), `compose.yaml`,
`README.md`, `CHANGELOG.md`, `.env.example`, `.gitignore`, plus die
Verzeichnisse `docker/`, `scripts/`, `docs/`.

Re-Init auf bestehendem Projekt (LH-FA-INIT-005) verlangt eine
explizite Strategie — kein stilles Überschreiben:

```bash
# Default: bestehende Dateien werden nicht angefasst
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git
# → Exit 10: "project already initialized"

# nur die U-BOOT MANAGED BLOCK-Regionen refreshen, User-Inhalt bleibt
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git --force

# Vollüberschreibung mit Sicherheits-Backup nach <datei>.bak[.N]
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  ghcr.io/pt9912/u-boot:0.1.0 init demo --no-git --force --backup
```

### `u-boot doctor` und die Container-Einschränkung

`doctor` ist für die **host-installierte** Form von u-boot
ausgelegt — es prüft `docker`, `docker compose` und `git` im
`$PATH`. Das Distroless-Image (`v0.1.0` und später) bringt keines
dieser Binaries mit (laut
[ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md)),
also können die Host-Probes aus einem `docker run …` heraus nicht
laufen.

Ab **`v0.1.1`** erkennt `doctor` die Container-Laufzeit via
`/.dockerenv` oder `/run/.containerenv` und emittiert für die vier
Host-Prerequisite-Checks eine `SeverityInfo`-Diagnostik
„skipped — running inside container" statt sie als Errors
fehlzudeuten. Exit-Code bei ansonsten gesundem Projekt ist `0`
(statt `11`). Designhintergrund:
[`slice-v0.1.1-doctor-container-awareness`](docs/plan/planning/done/slice-v0.1.1-doctor-container-awareness.md).

Für echte host-seitige Diagnostik `doctor` aus einer Host-
Installation laufen lassen, sobald die Binary-Distribution
gelandet ist
([`slice-v2-binary-distribution`](docs/plan/planning/done/slice-v2-binary-distribution.md),
ADR-0007 §Folgepunkte 1 Trigger jetzt aktiv). Die anderen
Subkommandos (`init`/`add`/`up`/`down`/`generate`/`config`)
funktionieren heute schon via Volume-Mount im Container.

### Build aus Quellen (Entwickler-Pfad)

Der Build ist **Docker-only** (`LH-FA-BUILD-007`): es wird keine
Go-Toolchain am Host benötigt. Nur Docker und `make` müssen
installiert sein.

```bash
make help                       # alle Targets auflisten
make build                      # Runtime-Image bauen (Distroless), Default VERSION=0.1.0-dev
make build VERSION=0.1.0        # Build mit gepinntem Version-Label
make run                        # Smoketest: docker run u-boot --help
make image-scan                 # lokaler Trivy-Scan (Parität mit CI image-scan-Job)
```

Inner-Loop-Quality-Gates (`LH-FA-BUILD-005` / `-006`):

```bash
make lint            # golangci-lint
make test            # go test ./...
make coverage-gate   # Coverage-Gate (bootstrap-aware, LH-FA-BUILD-008)
make gates           # lint + test + coverage-gate
make ci              # gates + govulncheck + image-scan
make fullbuild       # ci + build (vollständiger Closure-Lauf)
```

## Repository-Layout

```text
.
├── cmd/uboot/          # CLI-Entry-Point (`main.go`) — Wiring-Schicht
├── internal/           # hexagonales Layout (siehe spec/architecture.md)
│   ├── hexagon/{domain,application,port/{driving,driven}}/
│   └── adapter/{driving,driven}/
├── spec/               # Lastenheft + Architektur-Spezifikation
├── docs/               # ADRs, Planning, User-Doku (LH-FA-PROJDOCS-001)
├── scripts/            # Build-Helfer (coverage-gate.sh)
├── Dockerfile          # Multi-Stage-Build (LH-FA-BUILD-001)
├── Makefile            # Docker-only-Workflow (LH-FA-BUILD-005)
├── .dockerignore       # Build-Kontext-Filter (LH-FA-BUILD-004)
└── go.mod
```

Vollständiger Layout-Kontrakt:
[`LH-FA-BUILD-009` in `spec/lastenheft.md`](spec/lastenheft.md).

## Dokumentation

- **Lastenheft** (verbindliche Spezifikation):
  [`spec/lastenheft.md`](spec/lastenheft.md).
- **Architektur-Spezifikation:** [`spec/architecture.md`](spec/architecture.md)
  (hexagonales Pattern, Schicht-Regeln, depguard-Enforcement).
- **Quality Gates:** [`docs/user/quality.md`](docs/user/quality.md)
  (SOLID-nahes Lint-Profil §1.2, Carveouts §1.3, Tests §2,
  Coverage §3, Security §4).
- **Branch Protection:** [`docs/user/branch-protection.md`](docs/user/branch-protection.md)
  (LH-QA-003 PR-blockierende Checks, einmalige UI-Aktivierung).
- **Architecture Decision Records:**
  [`docs/plan/adr/`](docs/plan/adr/).
- **Planning-Artefakte (Slices, Tranchen):**
  [`docs/plan/planning/{open,next,in-progress,done}/`](docs/plan/planning/).
- **User-Dokumentation:** [`docs/user/`](docs/user/).

## Voraussetzungen

Für Konsumenten von `u-boot` (`LH-FA-DIAG-002`):

- Docker Engine ≥ 24.0.0 **oder** Podman ≥ 5.0 (Drop-in über
  `DOCKER_HOST=unix:///run/user/$UID/podman/podman.sock` und einen
  `docker → podman`-Symlink — siehe *Podman-Drop-in* unten).
- Docker Compose ≥ 2.20.0 **oder** `podman compose` (das
  containers/podman-compose-Plugin aus Podman 5.x).
- Git
- optional: VS Code mit der Dev-Containers-Extension

Für den Bau aus den Quellen (`LH-FA-BUILD-007`):

- Docker Engine (Podman funktioniert als Drop-in, ist aber heute
  nicht im CI abgedeckt — siehe „Podman-Drop-in" für die Caveats)
- GNU `make` (der einzige Carveout zu `LH-NFA-PORT-002`; Begründung
  siehe [`spec/lastenheft.md`](spec/lastenheft.md))

### Podman-Drop-in

u-boot ist auf Code-Ebene nicht Podman-aware — `DockerProbe`
shellt zu einer `docker`-Binary aus und parst Docker-Version-
Strings. Podman funktioniert als Drop-in, weil:

1. `podman` exposiert die gleiche CLI-Oberfläche, die u-boot
   braucht (`info`, `version`, `compose up/down/ps`, `build`,
   `push/pull`).
2. Die v0.1.1-Container-Detection (`slice-v0.1.1-doctor-container-
   awareness`) prüft bereits `/run/.containerenv` für Podman
   neben `/.dockerenv` für Docker.
3. Podman ≥ 4.0 liefert einen Docker-API-kompatiblen Socket;
   `DOCKER_HOST` darauf zeigen lassen, und jeder `docker`-CLI-
   Konsument spricht mit Podman.

Setup (typischer Linux-User):

```bash
# Rootless Podman-API-Socket starten.
systemctl --user enable --now podman.socket
export DOCKER_HOST=unix:///run/user/$(id -u)/podman/podman.sock

# Optional: docker→podman-Symlink für Tools, die exec("docker") machen.
sudo ln -sf "$(command -v podman)" /usr/local/bin/docker
```

Bekannte Caveats:

- `doctor` prüft `docker version` gegen die `LH-FA-DIAG-002`-
  Mindestwerte (24.0 / 2.20). Podmans Version-String ist
  parsbar, aber **eigene** Version (z. B. `5.3.1`), was heute
  als `Severity: warn — unrecognized version` klassifiziert
  wird statt `ok`. Funktional läuft `up`/`down`/`add` trotzdem.
- Keine CI-Matrix übt den Podman-Pfad aus; Bug-Reports gegen
  Podman sind willkommen, aber Blocking-Priorität ist Docker.
  Ein formaler Podman-Support-Slice landet bei konkretem
  Bedarf — siehe auch das v0.1.1 + ADR-0007 §Folgepunkte
  Trigger-Pattern.

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
