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

**MVP-Bootstrap.** Dieser Commit liefert die Build-Infrastruktur, das
Projekt-Skelett und einen `u-boot --help` / `u-boot --version`-Stub.
Subkommandos sind noch nicht implementiert; sie folgen in späteren
Slices, getrackt unter
[`docs/plan/planning/`](docs/plan/planning/).

| Phase | Status | Quelle |
| ----- | ------ | ------ |
| Lastenheft | Entwurf 0.1.0 | [`spec/lastenheft.md`](spec/lastenheft.md) |
| Architekturentscheidungen | 1 ADR | [`docs/plan/adr/`](docs/plan/adr/) |
| Implementierung | nur Bootstrap | [`docs/plan/planning/`](docs/plan/planning/) |

## Quickstart

Der Build ist **Docker-only** (`LH-FA-BUILD-007`): es wird keine
Go-Toolchain am Host benötigt. Nur Docker und `make` müssen installiert
sein.

```bash
make help            # alle Targets auflisten
make build           # Runtime-Image bauen (Distroless static, nonroot)
make run             # Smoketest: docker run u-boot --help
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
