# Slice V2: Homebrew-Formula für u-boot ([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Restweg)

> **Status:** **repo-seitig geliefert, Tap-Einrichtung offen** (2026-10-01,
> **Delivery-Hash: `0cab536`**). Die Entscheidung steht ([ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md)),
> Skeleton, Fill-Skript, Nachzug-Skript, Workflow-Änderung und Tap-Vorlagen sind im
> Repo; das Tap-Repo selbst und das Secret legt der Projektinhaber an (siehe
> §Offene Einrichtung unten). Der Slice bleibt bis dahin in `in-progress/`.

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

**Erste macOS-Nutzer-Nachfrage** (am 2026-10-01 durch den Projektinhaber ersetzt: Referenzen `pt9912/homebrew-ai-harness-init` und `pt9912/homebrew-d-migrate` geliefert, Umsetzung angewiesen). Solange das nicht passiert,
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
  [`carveouts.md`](carveouts.md) §Temporäre
  Carveouts, [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung)-Zeile.
- Roadmap:
  [`roadmap.md`](roadmap.md) §v0.4.0+ Backlog.
- Phase: V2 (nach v0.3.0-Milestone, Trigger-getrieben).

## Lieferstand (2026-10-01)

Vorbild: `pt9912/ai-harness-init` (Formel-Skeleton + Fill-Skript + Release-Asset + Tap-Nachzug-Job);
`d-migrate` nutzt dagegen die Drittanbieter-Action `homebrew-releaser` (Begründung der Wahl: [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md)).

| T | Stand |
| - | ----- |
| T1 | **geliefert (repo-seitig; Tap-Repo angelegt 2026-10-01):** `scripts/homebrew-formula.rb.tmpl` (Klasse `UBoot`, vier Plattformen, Test-Block `--version`), Fill-Skript `scripts/homebrew-formula-fill.sh` (fail-closed; Go-Test `cmd/uboot/homebrew_scripts_test.go`). **Offen:** `brew audit --strict` auf einem Mac. |
| T2 | **geliefert:** `publish.yml` erzeugt `SHA256SUMS`, füllt für stabile Tags die Formel und hängt sie als Asset `u-boot.rb` an; Job `tap` → `scripts/tap-nachzug.sh` (überspringt ohne Secret mit Warnung). **Offen:** Secret `HOMEBREW_TAP_GITHUB_TOKEN` (PAT mit Schreibrecht auf das Tap) setzen. |
| T3 | **geliefert (Vorlagen):** `packaging/homebrew-tap/README.md` und `.github/workflows/smoke.yml` (macOS: install, `--version`, `init`, `doctor` ohne Absturz) zum Kopieren ins Tap; README EN + DE mit `brew install`-Block. Im Tap-Repo (2026-10-01). |
| T4 | **teilweise:** Spec 0.3.4 ([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung): Homebrew gewählt), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) (statt Umschreiben des Accepted [ADR-0007](../../adr/0007-distributionswege-ghcr.md)), CHANGELOG, Carveout-Zeile angepasst. **Offen:** Verschiebung nach `done/` nach der ersten Tap-Installation. |

## Offene Einrichtung (nur der Projektinhaber)

1. ~~Repo `pt9912/homebrew-u-boot` anlegen~~ — **erledigt 2026-10-01** (öffentlich, `main`; README und
   `smoke.yml` aus `packaging/homebrew-tap/` hineinkopiert; die Formel folgt per Nachzug).
2. PAT (Contents: Read & Write auf dieses Repo) als Secret `HOMEBREW_TAP_GITHUB_TOKEN` im Repo
   `pt9912/u-boot` hinterlegen.
3. Nächsten **stabilen** Tag setzen; der Job `tap` füllt `Formula/u-boot.rb` und der Smoke-Workflow
   im Tap installiert sie auf macOS. Danach `brew install pt9912/u-boot/u-boot` auf einem Mac
   prüfen und diesen Slice nach `done/` verschieben.
