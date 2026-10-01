# Slice V1: Release-Cut `v0.6.0`

> **Status:** T1–T3 am 2026-10-01 ausgeführt (ein Commit); **T4** (Tag-Push)
> durch den Projektinhaber freigegeben und von der Session ausgeführt — Closure-Notiz
> am Dateiende. Ablauf: [`docs/user/releasing.md`](../../../user/releasing.md).

## Auslöser

Seit `v0.5.0` (2026-07-25) sind zwei große Pakete gelandet:

1. **Devcontainer-Sandbox-Profil** ([`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte), [`-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)..[`-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)),
   [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md), [ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md), [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).
2. **V1-Cleanup-Backlog:** Config-, Logs-, Up/Down-Cluster, rollback-aware
   `generate devcontainer`; dazu Homebrew ([ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md)) und der Go-Toolchain-Bump
   auf 1.27.1 (acht HIGH-stdlib-Advisories).

Das ist ein Minor-Release (`0.6.0`): neue Features, additive JSON-Felder, kein
Bruch bestehender Verträge.

## Aufhebungsbedingung

`git tag v0.6.0 && git push origin v0.6.0` triggert `publish.yml` (GHCR
`:0.6.0` + `:latest`, sechs Binaries, `SHA256SUMS`, `u-boot.rb`, Job `tap`);
`u-boot --version` des Release-Images zeigt `0.6.0`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1–T3 | (ein Commit) **CHANGELOG** `[Unreleased]` → `[0.6.0] - 2026-10-01` mit Lead und Compare-Links; **Versionsstrings** `0.5.0-dev` → `0.6.0-dev` (`main.go`, `Makefile`, `Dockerfile`); **READMEs** Status + Releases-Zeile; **Benutzerhandbuch** 1.2 / v0.6.0; Roadmap. `make ci` und `make test-docker` grün. |
| T4 | `git tag v0.6.0` auf dem grünen `main`-Stand, `git push origin v0.6.0`; `publish`-Lauf, Assets und Tap kontrollieren (siehe Closure-Notiz). |

## Out of Scope

- Homebrew-Core-Einreichung, Distro-Pakete ([`slice-v2-distro-pakete`](../open/slice-v2-distro-pakete.md)).
- Der Test `keycloak-ci-flake` ([`slice-v1-keycloak-ci-flake`](../open/slice-v1-keycloak-ci-flake.md)): kein Release-Blocker.

## Closure-Notiz (nach `done/`)

<!-- Nach dem Tag-Push füllen: Tag-Commit, publish-Lauf, Assets, Tap-Ergebnis. -->
