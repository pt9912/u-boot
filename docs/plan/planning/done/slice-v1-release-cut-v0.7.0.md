# Slice V1: Release-Cut `v0.7.0`

> **Status:** **abgeschlossen** (2026-10-01) — T1–T3 im Commit `cb8f937`, T4: Tag `v0.7.0`
> gesetzt, `publish` und Tap-Nachzug liefen grün (Closure-Notiz). Ablauf:
> [`docs/user/releasing.md`](../../../maintainer/releasing.md).

## Auslöser

Seit `v0.6.0` (2026-10-01) ist ein Paket gelandet: **Debian-/RPM-Pakete** per `nfpm` als Release-Assets
([`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung), [ADR-0017](../../adr/0017-linux-pakete-nfpm.md),
[`slice-v2-distro-pakete`](slice-v2-distro-pakete.md)); dazu der Keycloak-Test zurück in der Default-Lane
([`slice-v1-keycloak-ci-flake`](../done/slice-v1-keycloak-ci-flake.md)).

Das ist ein Minor-Release (`0.7.0`): neues Distributionsartefakt, kein Bruch bestehender Verträge.
Der Tag ist zugleich die Probe für zwei Dinge, die nur auf einem echten Tag belegbar sind: die
Paket-Smoke-Jobs (`.deb`/`.rpm`) und den Tap-Nachzug mit dem Secret `HOMEBREW_TAP_GITHUB_TOKEN`.

## Aufhebungsbedingung

`git tag v0.7.0 && git push origin v0.7.0` triggert `publish.yml` (GHCR `:0.7.0` + `:latest`, sechs
Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Smoke-Jobs, Job `tap`); `u-boot --version` des
Release-Images zeigt `0.7.0`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1–T3 | (ein Commit) **CHANGELOG** `[Unreleased]` → `[0.7.0] - 2026-10-01` mit Lead und Compare-Links; **Versionsstrings** `0.6.0-dev` → `0.7.0-dev` (`main.go`, `Makefile`, `Dockerfile`); **READMEs** Status + Releases-Zeile; **Benutzerhandbuch** 1.3 / v0.7.0; Roadmap. `make ci` und `make test-docker` grün. |
| T4 | `git tag v0.7.0` auf dem grünen `main`-Stand, `git push origin v0.7.0`; `publish`-Lauf, Assets, Paket-Smoke und Tap kontrollieren; Closure-Notiz. Danach können [`slice-v2-distro-pakete`](slice-v2-distro-pakete.md) und [`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md) schließen. |

## Out of Scope

- Homebrew-Core-Einreichung, gehostetes APT-/DNF-Repository, Paket-Signierung.

## Closure-Notiz

- **Tag:** `v0.7.0` auf `cb8f937` (2026-10-01), nach grünem `ci`, `integration-docker` und `make ci`.
- **`publish`-Lauf `36894700620`:** alle vier Jobs grün — `publish` (Image `ghcr.io/pt9912/u-boot:0.7.0` und
  `:latest` zeigen auf denselben Digest und melden `0.7.0`; sechs Binaries, vier Pakete, `SHA256SUMS`,
  `u-boot.rb`), `package-smoke-deb`, `package-smoke-rpm` und `tap` (Tap-Commit `3059299` „u-boot 0.7.0“;
  damit ist das schreibberechtigte Secret `HOMEBREW_TAP_GITHUB_TOKEN` im echten Lauf belegt).
- **Offen:** macOS-Smoke im Tap zur Formel 0.7.0 (Runner-Warteschlange) — siehe
  [`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md).
