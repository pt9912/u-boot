# Slice M2: Hexagonale Architektur

> **Status:** Done
> **DoD:** Commit `9d191a5`
> **Retro-Plan:** Retroaktiv geschrieben 2026-05-27 (siehe [`slice-m3-retroaktive-slice-plaene`](slice-m3-retroaktive-slice-plaene.md))

## Auslöser

M1 lieferte das Build-Skelett, aber die Architektur war noch unbestimmt.
Der erste fachliche Slice (M3 `u-boot init`) brauchte einen klaren
Layer-Schnitt für Tests (Fakes auf Driven-Ports) und Linter-Enforcement
(`depguard`-Regeln pro Schicht), bevor produktiver Code entstehen
konnte. Adressiert [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement).

## Lieferumfang

- **Spec-Erweiterung** in `spec/lastenheft.md`:
  - Sektion 4.13 mit [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) (hexagonales Pattern,
    Layer-Schnitt mit driving/driven-Split, Import-Regeln + Enforcement
    via `depguard`).
  - [`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout) Layout-Tree um `internal/{hexagon,adapter}/`
    erweitert.
  - [`LH-MVP-001`](../../../../spec/lastenheft.md#lh-mvp-001--muss-im-mvp-enthalten-sein) und Traceability-Matrix um die ARCH-Einträge ergänzt.
- **Architektur-Detailspec** `spec/architecture.md` (neu): Layer-Diagramm
  (ASCII), Verantwortlichkeiten je Schicht, Import-Regel-Tabelle, die 8
  `depguard`-Regelblöcke vollständig (scharf zu schalten mit M3),
  Test-Pattern, Anti-Patterns, Evolution.
- **[ADR-0002](../../adr/0002-hexagonale-architektur.md)** (`docs/plan/adr/0002-hexagonale-architektur.md`)
  begründet die Wahl gegenüber flachem Layout (k-deskflight) und
  alternativen Patterns (Clean / Onion); Trade-offs Boilerplate vs.
  Testbarkeit/Wartbarkeit explizit.
- **Code-Skelett**:
  - `internal/{hexagon/{domain,application,port/{driving,driven}},
    adapter/{driving,driven}}/` angelegt, je mit `README.md` mit
    Zweck/Inhalt/Import-Regeln.
  - `internal/README.md` auf das Layout aktualisiert.
- **`.golangci.yml`**: `depguard` aktiviert mit den 8 Layer-Regelblöcken
  (match-nichts-im-Bootstrap, scharf-werden-mit-M3).

## Akzeptanz

- `make gates` grün mit aktiviertem `depguard` (matchte im Bootstrap
  nichts; reale Verifikation kam mit M3-T5, siehe
  [`slice-m3-depguard-aktivierung-verifizieren`](slice-m3-depguard-aktivierung-verifizieren.md)).
- `spec/architecture.md` existiert und ist über READMEs verlinkt.
- [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) abgehakt.

## Bezug

- Auslösende Spec: [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement), [`LH-FA-BUILD-009`](../../../../spec/lastenheft.md#lh-fa-build-009--repository-layout).
- ADR: `0002-hexagonale-architektur.md`.
- Vorgänger: [`slice-m1-repo-skeleton`](slice-m1-repo-skeleton.md).
- Nachfolger: M2b (SOLID-Lint-Profil verschärft das Lint-Profil
  passend zur Architektur).
