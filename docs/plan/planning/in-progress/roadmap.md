# u-boot Roadmap

Übergreifendes Master-Dokument zum aktuellen Stand der Slices und
Tranchen (`LH-FA-PROJDOCS-003`). Diese Datei liegt dauerhaft in
`in-progress/`, bleibt aber bewusst knapp: Sie steuert aktuelle Arbeit,
Trigger und nächste Entscheidungen. Historische Release-Details liegen
in [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md).

## Aktueller Snapshot

| Version | Status | Datum | Fokus | Detailquelle |
| --- | --- | --- | --- | --- |
| v0.1.0 | released | 2026-05-31 | MVP-Core, Release-Pipeline, GHCR | [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md) |
| v0.2.0 | released | 2026-06-01 | Container-aware `doctor`, Binary-Distribution, Template-Katalog | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | released | 2026-06-01 | Add-on Catalogue Expansion (`remove`, deps, Keycloak, OTel) | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| v0.4.0 | in progress | 2026-06-02 | `logs`, JSON-/Dry-Run-CLI, restliche V1/Later-Trigger | diese Datei |

## v0.4.0 Arbeitspakete

Diese Punkte sind Arbeitspakete für Version 0.4.0. Einige sind bereits
als Slice-Plan angelegt, andere bleiben bis zur Ausarbeitung als
benannte APs in dieser Roadmap.

| AP | Status | Entscheidung / nächster Schritt |
| --- | --- | --- |
| [`slice-v1-logs`](slice-v1-logs.md) | `in-progress/`; T0 Discovery sowie T1 bis T3 im Slice dokumentiert | T4 Docker-E2E/Spec-Pin und T5 Doku/Closure abschließen, dann nach `done/` bewegen. |
| `slice-v1-cli-json-dry-run` | noch kein Slice-Plan | Maschinenlesbare CLI (`LH-FA-CLI-007/008`, `LH-NFA-USE-004`); stützt ADR-0010 als Alternative zu einem HTTP-Adapter. |
| [`slice-v1-keycloak-ci-flake`](../open/slice-v1-keycloak-ci-flake.md) | `open/`, on hold | Keycloak-Acceptance-Flake analysieren, sobald CI-Logs/Quay- oder Mirror-Befund belastbar sind. |
| [`slice-v2-homebrew-formula`](../open/slice-v2-homebrew-formula.md) | `open/`, on hold | Erste konkrete macOS-/Homebrew-Nutzeranfrage. |
| [`slice-v2-distro-pakete`](../open/slice-v2-distro-pakete.md) | `open/`, on hold | Konkrete Debian-/RPM-Anfrage mit Bereitschaft für Packaging-Overhead. |
| `slice-later-local-templates` | noch kein Slice-Plan | `--template ./pfad` konkretisieren (`LH-FA-TPL-003`). |
| `slice-later-migration` | noch kein Slice-Plan | Konfigurationsmigration konkretisieren (`LH-FA-CONF-006`). |
| `slice-later-custom-data-sources` | noch kein Slice-Plan | Erweiterung jenseits YAML-Quellen konkretisieren (`LH-DA-004`). |
| `slice-vN-podman-formal` | noch kein Slice-Plan | Podman-first Probe-Adapter und CI-Matrix konkretisieren; heutiger Stand bleibt Docker-compatible Drop-in. |
| Branch-Protection-UI | Nutzeraktion, kein Code-Slice | Repo-Owner aktiviert Required Checks vor erstem externem PR; Anleitung in [`docs/user/branch-protection.md`](../../../user/branch-protection.md). |

## Bereits Geschlossen

Die abgeschlossenen Slices bleiben in [`done/`](../done/) die
Detailquelle. Für die Agentensteuerung reicht hier der Cluster-Überblick:

| Cluster | Ergebnis | Detailquelle |
| --- | --- | --- |
| MVP M1..M8 | Repo-Skeleton, Architektur, CI/Gates, `init`, `doctor`, `add postgres`, `up/down`, `generate`, `config` | [`done/`](../done/) und [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md) |
| v0.2.0 | Container-aware `doctor`, sechs Plattform-Binaries, Template-Katalog | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | Add-on Catalogue Expansion: `remove`, `--with-deps`, Keycloak, OTel, V1-Audit | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| Harness-Doku | Agent-Briefing, Harness-Einstieg, Rollentrennung | [`AGENTS.md`](../../../../AGENTS.md), [`harness/README.md`](../../../../harness/README.md), [`harness/roles.md`](../../../../harness/roles.md) |

## Verwandte Dokumente

- [`carveouts.md`](carveouts.md) — Master-Inventar aller temporären und
  permanenten Carveouts (`LH-FA-PROJDOCS-005`), plus Audit-Trail der
  Slices, die offene Carveouts geschlossen haben.
- [`README.md`](../README.md) — Slice-/Tranche-Konventionen für
  Dateinamen in `docs/plan/planning/` (`LH-FA-PROJDOCS-003`).
- [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md)
  — ausgelagerte Release-Historie.

## Pflege-Regeln

- Diese Roadmap beschreibt aktuelle Steuerung, nicht jede historische
  Tranche.
- Release-Details, lange Commit-Listen und retrospektive Tabellen
  gehören in `done/`-Slices oder nach `docs/archive/`.
- Neue v0.4.0-Arbeitspakete brauchen eine klare Entscheidung oder einen
  nächsten Schritt. Ohne ausgearbeiteten Plan bleiben sie als benanntes
  AP hier, nicht als halbfertiger Slice.
- Diese Datei ist die einzige zulässige Ausnahme von der
  `slice-`/`tranche-`-Konvention für Dateinamen in
  `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und
  [`../README.md`](../README.md)).
