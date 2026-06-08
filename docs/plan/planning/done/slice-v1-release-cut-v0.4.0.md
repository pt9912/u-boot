# Slice V1: Release-Cut `v0.4.0`

> **Status:** `done/` — T1–T3 ✅ (2026-06-08), `make gates` grün.
> T4 ist die **Nutzer-Aktion** (Tag-Push). Vorbild
> [`slice-v1-release-cut-v0.3.0`](slice-v1-release-cut-v0.3.0.md)
> (identisches T1–T4-Pattern: drei Doku-Tranchen + eine
> Nutzer-Tranche).

## Auslöser

Seit `v0.3.0` (2026-06-01) ist der v0.4.0-Milestone **„Maschinen-
lesbare CLI"** auf `main` feature-complete. Kern-Deliverable ist der
vollständige `slice-v1-cli-json-dry-run`-Cluster plus der
Konsolidierungs-Folge-Slice:

- [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md) —
  9-teiliger Cluster-Slice (doctor/add/init/generate/remove/up-down/
  logs/config/template) + **T_close** (`3a35d58`). Liefert `--json`
  (`LH-NFA-USE-004`-Minimalkontrakt) für **alle zehn** Spec-Enum-
  Subcommands plus `--dry-run`/`--diff`-Voll-Schema (`LH-FA-CLI-007/
  008`) für die dateiverändernden Formen. Übergangs-Reject-Mechanik
  beim T_close abgebaut; `ADR-0010`-Re-Eval-Trigger-2 erfüllt.
- [`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md)
  — R15-Cross-Slice-1-Konsolidierung: add/init/generate adoptieren
  den geteilten `jsonArgsValidator` (Envelope-Symmetrie §1841 +
  Path-Leak-Defense).
- [`slice-v1-logs`](slice-v1-logs.md) — `u-boot logs [service]
  [--follow] [--tail <n>]` (`LH-FA-OBS-001`).
- [`slice-v1-devcontainer-features`](slice-v1-devcontainer-features.md)
  + [`slice-followup-devcontainer-features-drift-doctor`](slice-followup-devcontainer-features-drift-doctor.md)
  — Devcontainer-Features-Allowlist + Katalog + `devcontainer.
  features.drift`-Doctor-Check (Doctor-Total 12→13).

Die CHANGELOG-`## [Unreleased]`-Sektion sammelt die zugehörigen
`### Added`/`### Changed`/`### Fixed`-Einträge.

## Aufhebungsbedingung

`v0.4.0` ist veröffentlicht. Die Release-Maschinerie aus
[`slice-v1-release-pipeline`](slice-v1-release-pipeline.md) +
[`slice-v2-binary-distribution`](slice-v2-binary-distribution.md)
ist seit v0.2.0 unverändert in `publish.yml` aktiv. T1–T3 bereiten
Doku + Versionsstrings vor; T4 ist die Nutzer-Aktion (CHANGELOG-Datum
final + Tag-Push).

## Akzeptanzkriterien

- `CHANGELOG.md`: `## [Unreleased]` bleibt als leerer Anker, darunter
  eine `## [0.4.0] - 2026-06-08`-Sektion mit Milestone-Lead-Absatz,
  die die gesammelten Einträge trägt. Compare-Links um `v0.4.0`
  erweitert (`[Unreleased]: …v0.4.0...HEAD`, neu `[0.4.0]:
  …v0.3.0...v0.4.0`). ✅
- Dev-Default-Versionsstrings `0.3.0-dev` → `0.4.0-dev` in
  `cmd/uboot/main.go`, `Makefile`, `Dockerfile` (kosmetisch; tagged
  Releases überschreiben via `publish.yml` `make build VERSION=…`). ✅
- `README.md` + `README.de.md` Status-Block auf „v0.4.0 ready to
  tag" + Releases-Tabelle um die v0.4.0-Zeile erweitert. ✅
- `roadmap.md` §Snapshot v0.4.0-Zeile „in progress" → „🔖 ready to
  tag" + Verlinkung dieses Slice-Plans. ✅
- `make gates` grün gegen den v0.4.0-Stand. ✅

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1–T3 | (dieser Commit) | **CHANGELOG** `## [Unreleased]` → `## [0.4.0] - 2026-06-08` + Lead-Absatz + Compare-Links. **Versionsstrings** `0.3.0-dev`→`0.4.0-dev` (main.go/Makefile/Dockerfile). **README.{md,de.md}** Status + Releases-Tabelle; **roadmap** §Snapshot. Slice-Doc direkt in `done/`. `make gates` grün. |
| T4 | — | **Nutzer-Aktion** (analog v0.3.0-T4): (a) falls Tag-Push an anderem Datum: `## [0.4.0] - <Datum>` vor dem Push aktualisieren. (b) Lokale Commits auf `origin/main` pushen. (c) Ersten grünen CI-Lauf auf `main` abwarten. (d) `git tag v0.4.0 && git push origin v0.4.0` → `publish.yml` triggert GHCR-Push (`ghcr.io/pt9912/u-boot:0.4.0` + `:latest`) + Binary-Upload (sechs Plattformen). (e) Post-Push: roadmap §Snapshot auf „released" + Tag-Commit-Hash + Datum; CHANGELOG-`## [Unreleased]` bleibt leerer Anker für die nächste Iteration. |

## Out of Scope

- **Offene `open/`-Trigger-Stubs** (keycloak-ci-flake, config-list,
  config-structured-hint, volume-auto-removal, recreate-detection,
  multi-port, …): alle on-hold ohne gefeuerten Trigger — kein
  v0.4.0-Blocker, gelistet in [`roadmap.md`](../in-progress/roadmap.md)
  §v0.4.0+ Backlog und [`carveouts.md`](../in-progress/carveouts.md).
- **V2-Distribution** (Homebrew, Distro-Pakete) + **slice-later-***
  (local-templates, migration, custom-data-sources, podman-formal):
  Post-v0.4.0.
- **Branch-Protection-UI-Aktivierung**: user-getriebener One-Shot aus
  der v0.1.0-Era ([`docs/user/branch-protection.md`](../../../user/branch-protection.md)).
- **CHANGELOG-Kosmetik**: die `## [0.4.0]`-Sektion trägt zwei
  `### Added`-Strata (chronologisch geschichtet aus dem Unreleased-
  Block, Vorbestand) — vollständig + lesbar; eine optionale
  Zusammenführung ist Doku-Pflege, kein Release-Blocker.

## Bezug

- Vorbild: [`slice-v1-release-cut-v0.3.0`](slice-v1-release-cut-v0.3.0.md)
  (T1–T4-Struktur).
- Hängt von: [`slice-v1-release-pipeline`](slice-v1-release-pipeline.md)
  + [`slice-v2-binary-distribution`](slice-v2-binary-distribution.md)
  (`publish.yml`-Mechanik, seit v0.2.0 aktiv).
- Phase: V1-Release-Cut, keine Carveout-Auflösung — daher kein
  `carveouts.md`-Eintrag, nur Roadmap.
