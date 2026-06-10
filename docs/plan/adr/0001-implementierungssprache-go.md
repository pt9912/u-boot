# ADR 0001: Implementierungssprache Go

## Status

Accepted

## Datum

2026-05-21

## Kontext

Das Lastenheft ([`LH-OPEN-001`](../../../spec/lastenheft.md#lh-open-001-implementierungssprache-entschieden)) ließ die Implementierungssprache zunächst
offen und nannte vier Kandidaten: Go, Rust, Python, TypeScript/Node.js.
Vor Beginn der Codebase muss eine Entscheidung getroffen werden, damit
Build-Pipeline ([`LH-FA-BUILD-001`](../../../spec/lastenheft.md#lh-fa-build-001-multi-stage-dockerfile-u-boot-repo)..[`LH-FA-BUILD-009`](../../../spec/lastenheft.md#lh-fa-build-009-repository-layout)), Paketierung
([`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung)) und Toolchain-Pins festgezurrt werden können.

Entscheidungsrelevante Anforderungen aus dem Lastenheft:

- [`LH-NFA-PORT-001`](../../../spec/lastenheft.md#lh-nfa-port-001-linux-unterstützung) – Linux als primäre Plattform.
- [`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002-keine-unnötigen-systemabhängigkeiten) – möglichst wenige Systemabhängigkeiten am Host.
- [`LH-NFA-PERF-001`](../../../spec/lastenheft.md#lh-nfa-perf-001-schnelle-cli-antwort) – `u-boot --help`/`--version` unter 200 ms (Kaltstart).
- [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung) – Paketierung als Single-Binary bevorzugt.
- Vorlage `k-deskflight` (Go) für Multi-Stage Dockerfile + Makefile-Pattern.

## Entscheidung

**Go** als Implementierungssprache, mit folgenden konkreten Setzungen:

- Mindest-Toolchain: Go 1.26 (`go 1.26.0` in `go.mod`, analog
  `k-deskflight`).
- Default-Pin im Dockerfile (`ARG GO_VERSION`): `1.26.3` (aktuelle
  Stable-Version am Entscheidungsdatum).
- Modul-Pfad: `github.com/pt9912/u-boot`.
- Repository-Layout nach [`LH-FA-BUILD-009`](../../../spec/lastenheft.md#lh-fa-build-009-repository-layout): `cmd/uboot/`, `internal/`,
  Tests neben Production-Code. Die Substruktur unter `internal/`
  (hexagonale Schichten mit driving/driven-Split) ist in
  [`ADR-0002`](0002-hexagonale-architektur.md) festgelegt.
- Runtime-Image: `gcr.io/distroless/static-debian12:nonroot`, statisch
  gelinktes Binary (`CGO_ENABLED=0`, `-ldflags="-s -w"`).

## Konsequenzen

Positiv:

- Single statisch gelinktes Binary ohne Sprach-Laufzeit am Zielsystem
  (erfüllt [`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002-keine-unnötigen-systemabhängigkeiten) und vereinfacht [`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung)).
- Sehr schnelle Startzeit ([`LH-NFA-PERF-001`](../../../spec/lastenheft.md#lh-nfa-perf-001-schnelle-cli-antwort)); CLI-Frameworks
  (`spf13/cobra`, `urfave/cli`) sind etabliert.
- Erstklassige Standard-Library für YAML-/JSON-Verarbeitung,
  HTTP/`os/exec`-Aufrufe gegen Docker, sowie Cross-Compilation.
- Build-/CI-Pattern aus `k-deskflight` (Docker-only, Multi-Stage,
  Distroless-Runtime) ist 1:1 übernehmbar.

Negativ / Trade-offs:

- Generics in Go sind weniger ausdrucksstark als in Rust oder
  TypeScript; abstrakte Add-on-/Template-Interfaces erfordern
  konventionelle Patterns.
- Fehlerbehandlung über explizite Returns ist verbose; mit
  `errors.Is`/`errors.As` und benutzerdefinierten Fehlertypen
  beherrschbar.
- Reflection-basierte YAML-Bindings sind langsamer als handgeschriebene
  Encoder; für die zu erwartenden Konfigurationsgrößen unkritisch.

Offene Folgepunkte (eigene ADRs bei Bedarf):

- Paketierung im Detail ([`LH-OPEN-002`](../../../spec/lastenheft.md#lh-open-002-paketierung)): GHCR-Image, GitHub Release mit
  Binary-Artefakten, später ggf. Homebrew/Debian. Die konkrete
  Distributionsentscheidung erfolgt in eigenen ADRs.

Geschlossene Folgepunkte:

- `golangci-lint`-Profil und Pflicht-Linter-Set — entschieden mit
  [`ADR-0003`](0003-solid-nahes-lint-profil.md) (SOLID-nahes Profil
  mit 29 Lintern).
- CLI-Framework — entschieden mit [`ADR-0005`](0005-cli-framework-cobra.md)
  (Cobra v1.10.2).
