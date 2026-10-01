# Slice V1: Release-Cut `v0.6.0`

> **Status:** **abgeschlossen** (2026-10-01) — T1–T3 im Commit `76a8200`, T4: Tag `v0.6.0`
> gesetzt und `publish` durchgelaufen (Details in der Closure-Notiz; der Tap-Nachzug des Workflows
> scheiterte am Token und wurde von Hand ersetzt). Ablauf: [`docs/user/releasing.md`](../../../user/releasing.md).

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

- Homebrew-Core-Einreichung, Distro-Pakete ([`slice-v2-distro-pakete`](../in-progress/slice-v2-distro-pakete.md)).
- Der Test `keycloak-ci-flake` ([`slice-v1-keycloak-ci-flake`](../in-progress/slice-v1-keycloak-ci-flake.md)): kein Release-Blocker.

## Closure-Notiz (nach `done/`)

- **Tag:** `v0.6.0` auf `76a8200` (2026-10-01), nach grünem CI-Lauf und `make ci` / `make test-docker`.
- **`publish`-Lauf `36873768725`:** Job `publish` grün — Image `ghcr.io/pt9912/u-boot:0.6.0` und `:latest`
  melden `0.6.0`; sechs Binaries, `SHA256SUMS` und `u-boot.rb` am Release (Linux-Binary gegen die
  Summen geprüft). Job `tap` **rot**: `git push` ans Tap endete mit `403 Permission denied to pt9912`
  — das Secret `HOMEBREW_TAP_GITHUB_TOKEN` in `u-boot` hat (noch) kein Schreibrecht auf
  `homebrew-u-boot`. Die Formel wurde deshalb einmalig von Hand nachgezogen
  (`scripts/tap-nachzug.sh`, Tap-Commit `e191e30`).
- **Folgearbeit:** Der Workflow `tap-check` (`gh workflow run tap-check`, echter Push-Test) wurde
  nachgereicht; die API-Sicht allein hätte das falsche Token nicht erkannt (`permissions.push` zeigt
  bei klassischen PATs die Rechte des Benutzers). Ob das Secret korrekt ist, zeigt der nächste Release
  (Entscheidung des Projektinhabers: dort testen).
- **Offen:** macOS-Smoke im Tap (queued) und `brew install` auf einem Mac — siehe
  [`slice-v2-homebrew-formula`](../in-progress/slice-v2-homebrew-formula.md).
