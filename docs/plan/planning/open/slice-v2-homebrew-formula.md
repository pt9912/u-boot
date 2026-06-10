# Slice V2: Homebrew-Formula für u-boot ([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restweg)

> **Status:** on hold — Trigger noch nicht gefeuert. Plan-Stub
> existiert, damit [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)-Disziplin den Carveout-
> Anker erfüllt (siehe
> [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
> Carveouts, [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
> §Entscheidung Tabelle „Homebrew").

## Auslöser

`spec/lastenheft.md` §14 listet sechs Distributionswege als
mögliche Optionen für [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung). Drei sind in [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
gewählt (GHCR + Binary), zwei verworfen (npm, pip), zwei vertagt
mit eigenem Trigger-Slice — Homebrew und Debian/RPM.

Homebrew ist der natürliche Folgeweg nach Binary-Distribution:
sechs Plattformen sind via
[`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
T2 (`5e5166b`) ab v0.1.1 als GitHub-Release-Asset verfügbar; eine
Homebrew-Formula zieht das macOS-arm64/-amd64-Binary daraus und
bündelt es zu `brew install u-boot`. Ohne Binary keine
Homebrew-Formula — Voraussetzung ist also bereits erfüllt.

## Trigger

**Erste macOS-Nutzer-Nachfrage.** Solange das nicht passiert,
bleibt der Wartungs-Overhead (eigene Tap-Repo unter
`pt9912/homebrew-tap`, SHA256-Pin pro Release, CI-Smoke gegen
`brew install`-Pfad) ohne Mehrwert.

## Aufhebungsbedingung

`brew install pt9912/tap/u-boot` (oder analoger Pfad) installiert
das neueste `v*`-Tag-Binary auf einem frisch aufgesetzten macOS;
`u-boot --version` zeigt die korrekte Version; `u-boot doctor`
läuft ohne Errors.

## Akzeptanzkriterien

- ✅ Homebrew-Formula in einem `pt9912/homebrew-tap`-Repository
  (oder analoger Tap-Pfad), die das passende Binary aus dem
  GitHub-Release zieht — SHA256-Pin pro Plattform.
- ✅ `publish.yml` (oder ein zusätzlicher Workflow) aktualisiert
  die Formula automatisch pro Tag-Push — neuer Tag → neuer
  Formula-Commit mit Version + SHA-Updates.
- ✅ README-Install-Block (EN + DE) listet `brew install` als
  zweite-empfohlene Variante nach der Binary-Direct-Variante.
- ✅ Tap-Repo hat einen Smoke-Workflow, der nach jedem Formula-
  Update `brew install --build-from-source u-boot` ausführt + die
  drei `LH-AK-*`-Pre-Checks (init / add / doctor) durchläuft.

## Tranchen (vorgeschlagen, wird beim Trigger ausgearbeitet)

| T | Inhalt (Skizze) |
| - | --------------- |
| T1 | Tap-Repo `pt9912/homebrew-tap` anlegen + initiale Formula-Datei für die aktuelle v0.x.y-Version (manuelle SHA-Pins). Formula-Syntax-Check (`brew audit --strict`). |
| T2 | Automatisierung: `publish.yml` (oder neuer `homebrew-bump.yml`) erkennt v*-Tag-Push, holt SHA256 für Linux/macOS amd64/arm64 aus den GitHub-Release-Assets, committet ins Tap-Repo via GitHub-App-Token oder PAT. |
| T3 | Tap-Repo-Smoke-Workflow (macOS-Runner): `brew install` + `u-boot init demo --no-git` + `u-boot doctor`. Tap-Repo-README + main-Repo-READMEs (EN + DE) mit `brew install`-Install-Block. |
| T4 | Closure: CHANGELOG `## [Unreleased]` Added-Eintrag, carveouts.md [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restweg-Zeile reduziert (nur noch Debian/RPM offen), [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung „Vertagt → Gewählt" für Homebrew. Slice-Plan `open/` → `done/`. |

## Out of Scope

- **Homebrew-Core-Submission**: das ist ein eigener Antragsprozess
  beim Homebrew-Maintainerteam mit zusätzlichen Qualitätsregeln.
  Erst sinnvoll wenn das Tap-Setup stabil läuft und es eine
  nicht-triviale Nutzerbasis gibt.
- **Linux-Brew (`brew install` unter Linux)**: technisch unter-
  stützt, aber bietet keinen Mehrwert über das direkte Binary
  oder das GHCR-Image für Linux-User.

## Bezug

- Spec: [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) §14 (offene Distributionswege).
- ADR: [ADR-0007 §Entscheidung Tabelle „Homebrew"](../../adr/0007-distributionswege-ghcr.md)
  — verbindlicher Plan-Anker bis Trigger.
- Voraussetzungs-Slice:
  [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
  — Binaries existieren seit T2 `5e5166b` als GitHub-Release-Asset.
- Carveout:
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts, [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile.
- Roadmap:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0+ Backlog.
- Phase: V2 (nach v0.3.0-Milestone, Trigger-getrieben).
