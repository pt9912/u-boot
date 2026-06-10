# Slice V1: Release-Cut `v0.2.0`

## Auslöser

Seit `v0.1.0` (Tag-Push am 2026-05-31) sind drei substantielle
Features auf `main` gelandet, plus ein Patch-Slice und zwei
Review-Followups:

- [`slice-v0.1.1-doctor-container-awareness`](slice-v0.1.1-doctor-container-awareness.md) (4 Tranchen) —
  `doctor` skipped Host-Prerequisite-Checks im Container-Modus.
- [`slice-v2-binary-distribution`](slice-v2-binary-distribution.md) (4 Tranchen) — `make
  build-binaries` + `publish.yml`-Asset-Upload für sechs
  Plattformen (Linux/macOS/Windows × amd64/arm64).
- [`slice-v1-template-list`](slice-v1-template-list.md) (4 Tranchen + Review-Followup
  `c807cdb`) — `u-boot template list [--json]` ([`LH-FA-TPL-004`](../../../../spec/lastenheft.md#lh-fa-tpl-004-templates-auflisten))
  plus `basic`-Bootstrap-Metadaten.
- [`slice-v1-template-init`](slice-v1-template-init.md) (5 Tranchen + Review-Followup
  `7fe26e0`) — `u-boot init --template <name>` ([`LH-FA-TPL-001`](../../../../spec/lastenheft.md#lh-fa-tpl-001-projektvorlagen)
  für `basic`).

Die ursprüngliche Planung sah `v0.1.1` als reinen Patch-Release
für Doctor-Container vor. Mit den drei zwischenzeitlich
gelandeten Features (Binary, Template-List, Template-Init)
überschreitet der Scope die PATCH-Grenze; strikte Semver verlangt
einen MINOR-Bump. Daher **v0.2.0 statt v0.1.1**.

## Aufhebungsbedingung

`v0.2.0` ist veröffentlicht. Die Release-Maschinerie aus
[`slice-v1-release-pipeline`](slice-v1-release-pipeline.md) (T1..T5 done) plus die Binary-
Asset-Erweiterung aus [`slice-v2-binary-distribution`](slice-v2-binary-distribution.md) T2 sind
bereits in `publish.yml` aktiv; T1..T3 dieses Slices bereiten
die Dokumentation und Versionsstrings vor; T4 ist die
Nutzer-Aktion (CHANGELOG-Datum + Tag-Push).

## Akzeptanzkriterien

- `CHANGELOG.md` hat eine einzige `## [0.2.0] - <Datum>`-Sektion,
  die den vorherigen `## [Unreleased]`-Block und die noch nicht
  veröffentlichte `## [0.1.1] - TBD`-Sektion zusammenführt. Der
  `## [Unreleased]`-Header bleibt leer als Anker für die nächste
  Iteration. Die Compare-Links am Dateiende sind um `v0.2.0`
  erweitert.
- `cmd/uboot/main.go`, `Makefile` und `Dockerfile` haben die
  Dev-Default-Versionsstrings von `0.1.0-dev` auf `0.2.0-dev`
  gehoben (kosmetisch — bei tagged Releases überschreibt
  `publish.yml` ohnehin via `make build VERSION=…`).
- `README.md` und `README.de.md` Status-Block (Zeile ~31) sind
  von „`v0.1.1` in preparation" auf „v0.2.0 released" (post-
  Tag-Push) bzw. „v0.2.0 ready to tag" (pre-Tag-Push)
  umgemünzt. Die Feature-Liste nennt template-list, template-
  init, binary-distribution und doctor-container-awareness
  explizit.
- `roadmap.md` §Nächste Schritte 1 ist auf „v0.2.0-Tag-Push
  offen" umgestellt; der alte „v0.1.1-Tag-Push offen"-Wortlaut
  ist konsumiert. Der Punkt verweist auf den eigenen Slice
  (`done/[slice-v1-release-cut-v0.2.0.md](slice-v1-release-cut-v0.2.0.md)` nach T4).
- `make gates` grün gegen den v0.2.0-Stand.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `be139cb` | `CHANGELOG.md` konsolidiert: alter `## [Unreleased]`-Block (binary-distribution, template-list, template-init, Quickstart-Install) und `## [0.1.1] - TBD`-Block (RuntimeEnvironment-port, doctor-container-skip) in eine einzige `## [0.2.0] - TBD`-Sektion gemergt. Reihenfolge: Added (5 Einträge) → Changed (2) → Notes (1, latest/download-Caveat auf v0.2.0 angepasst). Lead-Paragraph erklärt warum v0.1.1 übersprungen wird. Compare-Links: `[Unreleased]: v0.2.0...HEAD`, neue Zeile `[0.2.0]: v0.1.0...v0.2.0`, `[0.1.0]` unverändert. |
| T2 | `1823598` | Dev-Default-Versionsstrings `0.1.0-dev` → `0.2.0-dev` in drei Files: `cmd/uboot/main.go` `var version`, `Makefile` `VERSION ?=`, `Dockerfile` `ARG UBOOT_VERSION=` (jeweils plus Doc-Comments). Bei tagged Releases überschreibt `publish.yml` ohnehin via `make build VERSION=<tag>`; betrifft nur lokale Outer-Loop-Aufrufe und untagged CI-Runs. Smoketest: `make build && docker run --rm u-boot --version` meldet `0.2.0-dev`. |
| T3 | dieser Commit | README.{md,de.md} Status-Block (Zeile ~31) von „`v0.1.1` in preparation" auf „`v0.2.0` ready to tag" mit Feature-Bullet-Liste (4 Slices) umgemünzt; `roadmap.md` §Nächste Schritte 1 von „v0.1.1-Tag-Push offen" auf „v0.2.0-Tag-Push offen" mit dem expandierten Scope und Verweis auf den eigenen Release-Cut-Slice. Slice-Move `open/` → `done/`. `make docs-check` grün. |
| T4 | — | **Nutzer-Aktion** (analog v0.1.0-T4): (a) Wenn der Tag-Push an einem anderen Tag als heute erfolgt, `## [0.2.0] - <Datum>` in `CHANGELOG.md` vor dem Push aktualisieren. (b) Lokale Commits auf `origin/main` pushen, falls nicht schon geschehen. (c) Ersten grünen CI-Lauf auf `main` abwarten. (d) `git tag v0.2.0 && git push origin v0.2.0` → `publish.yml` triggert GHCR-Push (Image `ghcr.io/pt9912/u-boot:0.2.0` + `:latest`) + Binary-Upload (sechs Plattformen als GitHub-Release-Asset). |

## Out of Scope

- **Branch-Protection-UI-Aktivierung** — bleibt als
  user-getriebener One-Shot offen aus der v0.1.0-Era; nicht
  v0.2.0-spezifisch, kann unabhängig nachgeholt werden
  (`docs/user/branch-protection.md`).
- **v0.1.1-Tag-Cut auf separatem Branch** — semver-strikt wäre
  ein Patch-Only-Release vom Doctor-Container-Commit möglich,
  aber Cherry-Pick + Branch-Management lohnt sich für ein
  Solo-Projekt nicht. v0.1.1 wird übersprungen.
- **V1-Add-ons** (Keycloak [`LH-AK-003`](../../../../spec/lastenheft.md#lh-ak-003-keycloak-flow), OTel [`LH-AK-004`](../../../../spec/lastenheft.md#lh-ak-004-opentelemetry-flow)),
  **V1-Generators** (`logs`, `--json`/`--dry-run`),
  **[`slice-later-local-templates`](slice-later-local-templates.md)** (`--template ./pfad`),
  **Template-Init-Variables** — alles Post-v0.2.0.

## Bezug

- Auslöser: drei Feature-Slices + ein Patch-Slice seit v0.1.0.
- Vorbild: [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md)
  T1-T4-Struktur; T1 dort enthielt einen Code-Bug-Fix (Version-
  Injection), den wir hier nicht haben, also ist unser T2 rein
  kosmetisch.
- Hängt von:
  [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
  T2/T3 (publish.yml-Mechanik) + [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung +
  [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
  T2 (Binary-Asset-Upload).
- Phase: V1-Release-Cut, keine Carveout-Auflösung — daher kein
  Eintrag in `carveouts.md`, nur in der Roadmap.
