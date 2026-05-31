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
aus [`spec/lastenheft.md`](spec/lastenheft.md) ist geliefert. Die
Release-Pipeline liegt bereit — GHCR-Image-Push auf `v*`-Tags via
[`.github/workflows/publish.yml`](.github/workflows/publish.yml),
Trivy als dritter PR-blockierender Job, Distributionsentscheidung
in [ADR-0007](docs/plan/adr/0007-distributionswege-ghcr.md). Der
erste Tag-Push selbst bleibt Nutzer-Trigger. Audit-Trail im
[MVP-Bilanz-Block der Roadmap](docs/plan/planning/in-progress/roadmap.md).

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

| Phase | Status | Quelle |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architekturentscheidungen | 6 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementierung | M1–M8 ✅, MVP-Closure ✅ — **MVP vollständig** | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 5 temporär (4 Slice-Pläne — Release-Pipeline-Slice deckt 2 Carveouts ab), 8 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

## Quickstart

Der Build ist **Docker-only** (`LH-FA-BUILD-007`): es wird keine
Go-Toolchain am Host benötigt. Nur Docker und `make` müssen installiert
sein.

```bash
make help            # alle Targets auflisten
make build           # Runtime-Image bauen (Distroless static, nonroot)
make run             # Smoketest: docker run u-boot --help
```

Echtes `u-boot init` gegen ein Host-Verzeichnis (Distroless läuft als
non-root UID 65532; `--user` matched die Host-UID, damit erzeugte
Dateien dir gehören):

```bash
mkdir /tmp/demo && \
  docker run --rm --user "$(id -u):$(id -g)" \
    -v /tmp/demo:/work -w /work \
    u-boot:latest init demo --no-git
```

Ergebnis: `u-boot.yaml` (`schemaVersion: 1`), `compose.yaml`,
`README.md`, `CHANGELOG.md`, `.env.example`, `.gitignore`, plus die
Verzeichnisse `docker/`, `scripts/`, `docs/`.

Re-Init auf bestehendem Projekt (LH-FA-INIT-005) verlangt eine
explizite Strategie — kein stilles Überschreiben:

```bash
# Default: bestehende Dateien werden nicht angefasst
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git
# → Exit 10: "project already initialized"

# nur die U-BOOT MANAGED BLOCK-Regionen refreshen, User-Inhalt bleibt
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git --force

# Vollüberschreibung mit Sicherheits-Backup nach <datei>.bak[.N]
docker run --rm --user "$(id -u):$(id -g)" -v /tmp/demo:/work -w /work \
  u-boot:latest init demo --no-git --force --backup
```

Inner-Loop-Quality-Gates (`LH-FA-BUILD-005` / `-006`):

```bash
make lint            # golangci-lint
make test            # go test ./...
make coverage-gate   # Coverage-Gate (bootstrap-aware, LH-FA-BUILD-008)
make gates           # lint + test + coverage-gate
make ci              # gates + govulncheck
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

- Docker Engine ≥ 24.0.0
- Docker Compose ≥ 2.20.0
- Git
- optional: VS Code mit der Dev-Containers-Extension

Für den Bau aus den Quellen (`LH-FA-BUILD-007`):

- Docker Engine
- GNU `make` (der einzige Carveout zu `LH-NFA-PORT-002`; Begründung
  siehe [`spec/lastenheft.md`](spec/lastenheft.md))

## Lizenz

MIT — siehe [`LICENSE`](LICENSE).
