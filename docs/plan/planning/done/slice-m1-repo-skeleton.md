# Slice M1: Repo-Skeleton

> **Status:** Done
> **DoD:** Commit `7da05c7`
> **Retro-Plan:** Retroaktiv geschrieben 2026-05-27 (siehe [`slice-m3-retroaktive-slice-plaene`](slice-m3-retroaktive-slice-plaene.md))

## Auslöser

u-boot startete von einem leeren Repo. M1 musste das gesamte Repo-
Skelett, die Docker-only-Build-Infrastruktur und die Doku-Struktur
gleichzeitig bringen, damit die nachfolgenden fachlichen Slices
(M2 Hexagonale Architektur, M3 `u-boot init`) auf einer
konsistenten Grundlage aufsetzen können.

Adressiert die Lastenheft-Anforderungen [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)..[`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout)
(Multi-Stage-Build, Pin-Politik, .dockerignore, Coverage-Bootstrap,
Make-Wrapper) und [`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) (`docs/`-Struktur,
ADR-Format, Roadmap als Master-Dokument).

## Lieferumfang

- `go.mod` mit Modul-Pfad `github.com/pt9912/u-boot`, Go 1.26.0.
- `cmd/uboot/main.go` mit CLI-Stub: `--help`, `--version`, Exit-Code 2
  für unbekannte Subkommandos ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)). Tests decken alle
  Stub-Pfade.
- `Dockerfile` mit BuildKit-Syntax, Multi-Stage-Pipeline
  (deps/compile/lint/test/coverage/build/runtime); Runtime-Image
  `gcr.io/distroless/static-debian12:nonroot`, OCI-Labels,
  USER 65532:65532 ([`LH-FA-BUILD-002`](../../../../spec/lastenheft.md#lh-fa-build-002--runtime-stage-pflichten)).
- `Makefile` mit Docker-only Inner-Loop-Targets (deps/compile/lint/
  test/coverage-gate/build/run/clean) und Aggregatoren (gates/ci/
  fullbuild).
- `scripts/coverage-gate.sh` bootstrap-aware ([`LH-FA-BUILD-008`](../../../../spec/lastenheft.md#lh-fa-build-008--coverage-bootstrap)).
- `.dockerignore` ([`LH-FA-BUILD-004`](../../../../spec/lastenheft.md#lh-fa-build-004--dockerignore-pflicht)), `.gitignore`.
- `.golangci.yml` v2-Schema mit 5 Default-Lintern; volles SOLID-Profil
  folgte in M2b.
- `docs/{archive,user,plan/{adr,planning/{open,next,in-progress,done}}}/`
  mit READMEs.
- `docs/plan/adr/0001-implementierungssprache-go.md` (ADR-Format
  nach [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format), referenziert von [`LH-OPEN-001`](../../../../spec/lastenheft.md#lh-open-001--implementierungssprache-entschieden)).
- `docs/plan/planning/in-progress/roadmap.md` als Master-Dokument
  (Spec-Ausnahme zur `slice-`/`tranche-`-Namens-Konvention dokumentiert).
- `README.md` (Englisch) + `README.de.md` (Deutsch).

## Akzeptanz

- `make gates` grün (lint/test/coverage-gate im Bootstrap-Modus,
  `./internal/...` leer → bootstrap-OK).
- `docker run` über das Runtime-Image liefert `u-boot --help` /
  `--version`-Output.
- [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)..[`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout) und [`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) abgehakt.

## Bezug

- Auslösende Spec: [`LH-FA-BUILD-001`](../../../../spec/lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)..[`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout), [`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle),
  [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes).
- ADR: `0001-implementierungssprache-go.md`.
- Nachfolger: M2 (Hexagonale Architektur) baut auf der `internal/`-
  Struktur auf, die hier als leeres Skelett vorbereitet wurde.
