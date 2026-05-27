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

**MVP in Arbeit — `u-boot init` ist vollständig verdrahtet inkl. Re-Init.**
Das erste fachliche Subkommando ist end-to-end ausgeliefert (M3 ✅):
`u-boot init [name]` erzeugt die LH-FA-INIT-003-Projektstruktur plus
`u-boot.yaml` (LH-FA-CONF-002) und initialisiert per Default ein
Git-Repository (LH-FA-INIT-007); ein zweiter Lauf auf bestehendem
Projekt nutzt den LH-FA-INIT-005-Überschreibschutz (`--force` für
Managed-Block-only-Edits, `--backup` für Vollüberschreibung mit
`.bak[.N]`-Sicherung) plus die LH-FA-CLI-005A-Modi-Flags (`--yes` /
`--no-interactive` exklusiv, `--assume-existing` durchgereicht für
M4-Soft-Detection). Die weiteren MVP-Subkommandos (`add`, `up`, `down`,
`doctor`, `generate`, `config`) folgen in M4+; Planung in
[`docs/plan/planning/`](docs/plan/planning/).

| Phase | Status | Quelle |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architekturentscheidungen | 5 ADRs | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementierung | M1–M3 ✅, M4 als Nächstes | [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) |
| Carveouts | 11 temporär (9 mit Slice-Plan, 1 Slice deckt 3), 7 permanent | [`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md) |

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
- **Architecture Decision Records:**
  [`docs/plan/adr/`](docs/plan/adr/).
- **Planning-Artefakte (Slices, Tranchen):**
  [`docs/plan/planning/{open,next,in-progress,done}/`](docs/plan/planning/).
- **User-Dokumentation:** [`docs/user/`](docs/user/) (während des
  Bootstrap leer).

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
