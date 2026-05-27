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

Adressiert die Lastenheft-Anforderungen `LH-FA-BUILD-001..009`
(Multi-Stage-Build, Pin-Politik, .dockerignore, Coverage-Bootstrap,
Make-Wrapper) und `LH-FA-PROJDOCS-001..003` (`docs/`-Struktur,
ADR-Format, Roadmap als Master-Dokument).

## Lieferumfang

- `go.mod` mit Modul-Pfad `github.com/pt9912/u-boot`, Go 1.26.0.
- `cmd/uboot/main.go` mit CLI-Stub: `--help`, `--version`, Exit-Code 2
  für unbekannte Subkommandos (LH-FA-CLI-006). Tests decken alle
  Stub-Pfade.
- `Dockerfile` mit BuildKit-Syntax, Multi-Stage-Pipeline
  (deps/compile/lint/test/coverage/build/runtime); Runtime-Image
  `gcr.io/distroless/static-debian12:nonroot`, OCI-Labels,
  USER 65532:65532 (LH-FA-BUILD-002).
- `Makefile` mit Docker-only Inner-Loop-Targets (deps/compile/lint/
  test/coverage-gate/build/run/clean) und Aggregatoren (gates/ci/
  fullbuild).
- `scripts/coverage-gate.sh` bootstrap-aware (LH-FA-BUILD-008).
- `.dockerignore` (LH-FA-BUILD-004), `.gitignore`.
- `.golangci.yml` v2-Schema mit 5 Default-Lintern; volles SOLID-Profil
  folgte in M2b.
- `docs/{archive,user,plan/{adr,planning/{open,next,in-progress,done}}}/`
  mit READMEs.
- `docs/plan/adr/0001-implementierungssprache-go.md` (ADR-Format
  nach LH-FA-PROJDOCS-002, referenziert von `LH-OPEN-001`).
- `docs/plan/planning/in-progress/roadmap.md` als Master-Dokument
  (Spec-Ausnahme zur `slice-`/`tranche-`-Namens-Konvention dokumentiert).
- `README.md` (Englisch) + `README.de.md` (Deutsch).

## Akzeptanz

- `make gates` grün (lint/test/coverage-gate im Bootstrap-Modus,
  `./internal/...` leer → bootstrap-OK).
- `docker run` über das Runtime-Image liefert `u-boot --help` /
  `--version`-Output.
- LH-FA-BUILD-001..009 und LH-FA-PROJDOCS-001..003 abgehakt.

## Bezug

- Auslösende Spec: `LH-FA-BUILD-001..009`, `LH-FA-PROJDOCS-001..003`,
  `LH-FA-CLI-006`.
- ADR: `0001-implementierungssprache-go.md`.
- Nachfolger: M2 (Hexagonale Architektur) baut auf der `internal/`-
  Struktur auf, die hier als leeres Skelett vorbereitet wurde.
