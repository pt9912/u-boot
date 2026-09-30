# Slice V1: Sandbox-Profil für Devcontainer umsetzen

## Auslöser

Der CR [`slice-cr-sandbox-devcontainer`](../done/slice-cr-sandbox-devcontainer.md)
hat Lastenheft 0.3.0 angenommen. Dieser Slice setzt das V1-Paket um:
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
und die Ergänzung von
[`LH-FA-DEV-004`](../../../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte)
(`devcontainer.user.uid`). `-008` (V2, Egress) ist ausgenommen.

## Umfang (Skizze, vor Tranchen-Zerlegung)

1. Umsetzungs-ADR: Base-Image für Podman-im-Container; Verifikation
   NET_ADMIN vs. rootless User-Namespaces und Seccomp/AppArmor unter
   Docker (auch Negativbefunde).
2. Config-Keys `devcontainer.user.uid`, `devcontainer.profile`,
   `devcontainer.sandbox.nestedRuntime|onUnavailable` inkl. Validierung
   (Exit 10) und Exit-11-Pfade mit Sentinel-Pins.
3. `--sandbox` in `init --devcontainer` / `generate devcontainer`;
   Golden Cases (Fresh-State, Idempotenz, Safety) nach `harness/replay.md`.
4. Doctor-Checks (Token-Quelle, Degradationszustände).
5. Integrationstest (`//go:build docker`) unter Docker (Linux); Colima/macOS und Podman-Host per Nachholmessung des Projektinhabers (Verfahren in ADR 0014) — bis dahin als ungeprüfte Einschränkung dokumentiert;
   `docs/user/`, README, CHANGELOG.

## Trigger

`open` → `next`: erfolgt am 2026-09-30 (Projektinhaber). `next` → `in-progress`: Beginn der Umsetzung; Umsetzungs-ADR 0014 liegt als `Proposed` vor.
