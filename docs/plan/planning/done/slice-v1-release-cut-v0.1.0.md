# Slice V1: Release-Cut `v0.1.0`

> **Status:** Done (T1 `056e4c6`, T2 `f176e95`, T3 dieser Commit;
> T4 Tag-Push bleibt Nutzer-Aktion).

## Auslöser

Mit `slice-v1-release-pipeline` (T1..T5 done, `bc487fc`) ist die
Maschinerie für GHCR-Image-Publish auf Tag `v*` bereit. Auf die
Frage „Können wir schon ein MVP-Release erstellen?" (2026-05-31)
hat ein Pre-Release-Check vier Blocker gefunden, davon einen
echten Code-Bug:

1. **Version-Verankerung fehlt im Build** (Code-Bug): `Dockerfile`
   baute `go build -ldflags="-s -w"` ohne `-X main.version=...`,
   sodass jeder Build den in-source-Fallback
   `var version = "0.1.0-dev"` (`cmd/uboot/main.go:37`)
   konservierte. `ghcr.io/pt9912/u-boot:0.1.0` hätte
   `u-boot --version` mit „0.1.0-dev" beantwortet — falsch.
2. **CHANGELOG.md im Repo-Root fehlt** (Convention-Gap): u-boot
   hat zwar einen `generate changelog`-Handler für Nutzer-Projekte
   (`LH-AK-007`), aber keinen eigenen Top-Level-Changelog.
3. **31 lokale Commits, nicht gepusht** (Prerequisite): alle
   Pipeline-Bausteine seit `e0d6c87` waren nur lokal.
4. **Branch-Protection-UI nicht aktiviert** (Nutzer-Aktion):
   `docs/user/branch-protection.md` beschreibt die Pflicht-
   Aktivierung, aber sie ist nicht automatisierbar.

## Aufhebungsbedingung

Die vier Blocker sind beseitigt:

- T1 fixt (1) — Build mit korrekter VERSION-Injection.
- T2 fixt (2) — CHANGELOG.md angelegt.
- T3 (dieser Slice) bündelt (1)+(2), dokumentiert die
  Nutzer-Aktionen (3) Push und (4) Branch-Protection + Tag-Push.

Erster `v0.1.0`-Tag-Push selbst bleibt **bewusst** Nutzer-Aktion,
nicht Auto-Trigger.

## Akzeptanzkriterien

- `make build` (ohne Override) liefert `u-boot --version`
  → `0.1.0-dev` und OCI-Label `org.opencontainers.image.version=0.1.0-dev`.
  Identisch für alle CI-Pfade.
- `make build VERSION=0.1.0` liefert `u-boot --version` → `0.1.0`
  und OCI-Label `0.1.0`.
- publish.yml hat zwei zusätzliche Pinning-Steps:
  (a) `org.opencontainers.image.version`-Label gegen
  `${{ steps.tag.outputs.version }}`,
  (b) Live-Smoketest `docker run --rm $REF --version` gegen
  denselben Wert. Eine VERSION-Drift (z. B. vergessener
  Build-Arg) bricht den Workflow vor dem Push.
- `CHANGELOG.md` im Keep-a-Changelog-Format mit `## [Unreleased]`
  und `## [0.1.0]`-Sektion (Added / Known limitations / Setup).
  Compare-Links am Ende für künftige Tag-Pushs erweiterbar.
- `make gates` und `make image-scan` lokal grün gegen den
  v0.1.0-Stand.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `056e4c6` | Dockerfile `ARG UBOOT_VERSION` (build + runtime stage), `-X main.version=...` in go build, `org.opencontainers.image.version`-Label; Makefile `VERSION ?= 0.1.0-dev` + `--build-arg UBOOT_VERSION=$(VERSION)`; publish.yml reicht `VERSION` an `make build` durch und pinnt das image.version-Label + live `--version`-Smoke gegen die Tag-VERSION. |
| T2 | `f176e95` | `CHANGELOG.md` im Repo-Root angelegt, vier Sektionen für `0.1.0` (Added Subcommands / Added CI+release infra / Added Architecture & docs / Known limitations + Setup). |
| T3 | dieser Commit | Slice-Plan + carveouts/roadmap + READMEs-Sync. |
| T4 | — | **Nutzer-Aktion:** (a) 33+ lokale Commits auf `origin/main` pushen, (b) Branch-Protection-Required-Status-Checks im GitHub-UI aktivieren gemäß `docs/user/branch-protection.md` (drei verbose `name:`-Felder), (c) ersten grünen CI-Lauf auf `main` abwarten, (d) `git tag v0.1.0 && git push origin v0.1.0` → `publish.yml` triggert GHCR-Image-Push. |

## Out of Scope

- Erweiterte Distributionswege (Binary, Homebrew, Debian/RPM) — in
  ADR-0007 §Entscheidung als vertagt mit Trigger-Slices benannt.
- V1-Add-ons (Keycloak, OTel), V1-Generators (`logs`, `--json`),
  Template-Implementation-Slices (`template-list`/`-init`/
  `local-templates`) — jeweils eigener V1-Slice.

## Bezug

- Auslöser: `slice-v1-release-pipeline` lieferte die Mechanik;
  dieser Slice füllt die letzten Lücken vor dem ersten Tag.
- Hängt von:
  [`slice-v1-release-pipeline`](slice-v1-release-pipeline.md)
  T2/T3,
  [ADR-0007](../../adr/0007-distributionswege-ghcr.md),
  [`docs/user/branch-protection.md`](../../../user/branch-protection.md).
- Phase: V1 (release-cut), keine Carveout-Auflösung — daher kein
  Eintrag in [`carveouts.md`](../in-progress/carveouts.md),
  sondern nur in der Roadmap.
