# Slice V1: Release-Cut `v0.3.0`

## Auslöser

Seit `v0.2.0` (Tag-Push am 2026-06-01) ist der v0.3.0-Milestone
„Add-on Catalogue Expansion" auf `main` feature-complete: fünf
Slices, drei davon mit substantiellen Code-Erweiterungen, plus
ein Doku-Audit-Slice und ein Pure-Refactor-Slice für die
Dependency-Mechanik.

- [`slice-v1-audit-done`](../done/slice-v1-audit-done.md) —
  Doku-Audit für `LH-FA-BUILD-006`/`LH-NFA-MAINT-004`/
  `LH-NFA-PORT-003`. Pure-Doku.
- [`slice-v1-add-remove`](../done/slice-v1-add-remove.md) — fünf
  Tranchen + Review-Followup `78ddcc6` (F1..F6). Liefert
  `u-boot remove <service> [--purge]` analog M5-Postgres-Pattern,
  extrahiert `detectServiceState` zur Package-Free-Function.
- [`slice-v1-addons-deps`](../done/slice-v1-addons-deps.md) — vier
  Tranchen. Domain `AddOnDependency` + Resolver +
  `Confirmer.ConfirmAddDependency`-Port-Extension +
  `--with-deps`-CLI mit Vier-Modi-Dispatch. Breaking refactor:
  `NewAddServiceService` nimmt jetzt einen `Confirmer`.
- [`slice-v1-keycloak`](../done/slice-v1-keycloak.md) — vier
  Tranchen. `u-boot add keycloak` (LH-FA-ADD-003 / LH-AK-003);
  per-Service Probe-Mechanismus (`requiredEnvKeys` /
  `volumeRefLiteral` / `volumeOptional` Catalogue-Felder);
  `acceptance_helpers.go`-Extraktion (init+add+up-Pipeline für
  alle künftigen Acceptance-Docker-Tests).
- [`slice-v1-otel`](../done/slice-v1-otel.md) — vier Tranchen.
  `u-boot add otel` (LH-FA-ADD-004 / LH-AK-004); zweites
  Catalogue-Pattern-Update: `extraFiles []extraFileEntry` für
  whole-file artefacts + `healthcheckOptional`. Makefile-Patch
  `test-docker` `-v /tmp:/tmp` für Compose-Bind-Mount-Auflösung
  in CI (echte Diagnose statt Carveout-Eskalation).

Plus zwei begleitende Slice-Pläne in `open/`, die als Folge-Slices
ohne aktive Tranchen liegen:

- [`slice-v1-keycloak-ci-flake`](slice-v1-keycloak-ci-flake.md) —
  offener Trigger-Slice für die `acceptance_extended`-Carveout-
  Auflösung des Keycloak-Acceptance-Tests. Nicht v0.3.0-blocking.
- [`slice-v2-homebrew-formula.md`](slice-v2-homebrew-formula.md)
  + [`slice-v2-distro-pakete.md`](slice-v2-distro-pakete.md) —
  proaktive Trigger-Stubs für `LH-OPEN-002`-Restwege. Bewusst
  ohne v0.3.0-Bezug.

## Aufhebungsbedingung

`v0.3.0` ist veröffentlicht. Die Release-Maschinerie aus
[`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
plus die Binary-Asset-Erweiterung aus
[`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
sind seit v0.2.0 aktiv und unverändert. T1..T3 dieses Slices
bereiten Doku und Versionsstrings vor; T4 ist die Nutzer-Aktion
(CHANGELOG-Datum + Tag-Push).

## Akzeptanzkriterien

- `CHANGELOG.md` hat eine einzige `## [0.3.0] - <Datum>`-Sektion,
  die den vorherigen `## [Unreleased]`-Block zusammenführt
  (drei „Added"-Einträge für add-remove, addons-deps, keycloak,
  otel + ein „Verified"-Eintrag für audit-done). Der
  `## [Unreleased]`-Header bleibt leer als Anker für die nächste
  Iteration. Die Compare-Links am Dateiende sind um `v0.3.0`
  erweitert.
- `cmd/uboot/main.go`, `Makefile` und `Dockerfile` haben die
  Dev-Default-Versionsstrings von `0.2.0-dev` auf `0.3.0-dev`
  gehoben (kosmetisch — bei tagged Releases überschreibt
  `publish.yml` ohnehin via `make build VERSION=…`).
- `README.md` und `README.de.md` Status-Block sind von
  „v0.2.0 released" auf „v0.3.0 ready to tag" (pre-Tag-Push)
  bzw. „v0.3.0 released" (post-Tag-Push) umgemünzt. Die
  Feature-Liste nennt die fünf v0.3.0-Slices explizit, plus
  `u-boot add keycloak / otel` in der Subcommand-Reference
  (schon mit slice-v1-keycloak T4 + slice-v1-otel T4 gelandet).
- `roadmap.md` `§Releases`-Tabelle: v0.3.0-Zeile von
  „feature-complete (5/5)" auf „✅ released" + Tag-Commit-Hash
  + Datum gehoben. `§v0.1.0/v0.2.0 — Audit-Trail`-Sektion
  umbenannt zu `§v0.1.0 / v0.2.0 / v0.3.0 — Audit-Trail`,
  v0.3.0-Vorspann mit Tag-Push-Details ergänzt.
  `§v0.3.0 — Milestone-Tabelle` bleibt unverändert (zeigt jetzt
  endgültigen Stand 5/5 ✅).
- `make gates` grün gegen den v0.3.0-Stand.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | `CHANGELOG.md` konsolidiert: der `## [Unreleased]`-Block (drei Added-Einträge für `u-boot add otel`, `u-boot add keycloak`, `u-boot add <service> --with-deps`, ein Added-Eintrag für `u-boot remove`, plus ein Verified-Eintrag für audit-done) zu einer einzigen `## [0.3.0] - TBD`-Sektion gemergt. Reihenfolge: Added (4 Einträge) → Verified (1) → Notes (falls nötig). Lead-Paragraph erklärt den Milestone-Scope „Add-on Catalogue Expansion" und benennt die drei Add-ons explizit. Compare-Links: `[Unreleased]: v0.3.0...HEAD`, neue Zeile `[0.3.0]: v0.2.0...v0.3.0`, `[0.2.0]` + `[0.1.0]` unverändert. |
| T2 | Dev-Default-Versionsstrings `0.2.0-dev` → `0.3.0-dev` in drei Files: `cmd/uboot/main.go` `var version`, `Makefile` `VERSION ?=`, `Dockerfile` `ARG UBOOT_VERSION=` (jeweils plus Doc-Comments). Bei tagged Releases überschreibt `publish.yml` ohnehin via `make build VERSION=<tag>`; betrifft nur lokale Outer-Loop-Aufrufe und untagged CI-Runs. Smoketest: `make build && docker run --rm u-boot --version` meldet `0.3.0-dev`. |
| T3 | README.{md,de.md} Status-Block (Zeile ~31) von „v0.2.0 released" auf „v0.3.0 ready to tag" mit Feature-Bullet-Liste (5 Slices); `roadmap.md` §Releases-Tabelle v0.3.0-Zeile auf „🔖 ready to tag" + verlinkt den eigenen Release-Cut-Slice; Slice-Move `open/` → `done/`. `make docs-check` grün. |
| T4 | — **Nutzer-Aktion** (analog v0.2.0-T4): (a) Wenn der Tag-Push an einem anderen Tag als T1-T3-Commits erfolgt, `## [0.3.0] - <Datum>` in `CHANGELOG.md` vor dem Push aktualisieren. (b) Lokale Commits auf `origin/main` pushen, falls nicht schon geschehen. (c) Ersten grünen CI-Lauf auf `main` abwarten (inkl. `integration-docker` mit dem T3-`-v /tmp:/tmp`-Fix aktiv). (d) `git tag v0.3.0 && git push origin v0.3.0` → `publish.yml` triggert GHCR-Push (Image `ghcr.io/pt9912/u-boot:0.3.0` + `:latest`) + Binary-Upload (sechs Plattformen als GitHub-Release-Asset). (e) Post-Push: T1-Datum + Releases-Tabelle in roadmap auf „✅ released" + Tag-Commit-Hash + Datum hieven (eigener Doku-Commit), und `## [Unreleased]`-Block in CHANGELOG für die nächste Iteration leeren. |

## Out of Scope

- **`slice-v1-keycloak-ci-flake`-Auflösung**: der Keycloak-
  Acceptance-Test ist seit `9d0be1c` hinter `acceptance_extended`
  versteckt; der CI-Bind-Mount-Fix aus slice-v1-otel T3 könnte
  ihn vermutlich befreien, das ist aber nicht v0.3.0-blocking.
  Trigger bleibt offen; bei der Diagnose-Runde nach v0.3.0
  könnte der Test schnell zurück in die Pflicht-Lane (die echte
  Failure-Cause wurde damals nicht isoliert, deshalb der
  Carveout — heute könnte sie es).
- **Branch-Protection-UI-Aktivierung** — bleibt als
  user-getriebener One-Shot offen aus der v0.1.0-Era; nicht
  v0.3.0-spezifisch (`docs/user/branch-protection.md`).
- **V1-Generators** (`u-boot logs`, `--json`/`--dry-run`),
  **`slice-later-local-templates`** (`--template ./pfad`),
  **Migration / Custom-Data-Sources**, **Homebrew / Distro-
  Pakete** — alles Post-v0.3.0, gelistet in
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0+ Backlog.

## Bezug

- Auslöser: vier Feature-Slices + ein Doku-Audit-Slice seit
  v0.2.0, Milestone „Add-on Catalogue Expansion" feature-
  complete.
- Vorbild:
  [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md)
  T1-T4-Struktur — identisches Pattern, drei Doku-Tranchen
  + eine Nutzer-Tranche.
- Hängt von:
  [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
  T2/T3 (publish.yml-Mechanik) + ADR-0007 §Entscheidung +
  [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md)
  T2 (Binary-Asset-Upload). Beide seit v0.2.0 unverändert in
  `publish.yml` aktiv.
- Milestone:
  [v0.3.0](../in-progress/roadmap.md#v030) — letzter Slice des
  Milestones.
- Phase: V1-Release-Cut, keine Carveout-Auflösung — daher kein
  Eintrag in `carveouts.md`, nur in der Roadmap.
