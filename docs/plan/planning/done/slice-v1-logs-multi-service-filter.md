# Slice V1: `u-boot logs <svc1> <svc2>` Multi-Service-Filter

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `4db9f90`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
> §Out of Scope (Carveout-Eintrag entfernt).

## Auslöser

Heutige `u-boot logs`-Surface ([slice-v1-logs](../done/slice-v1-logs.md) §AK) nutzt
`cobra.MaximumNArgs(1)` — Single-Service oder Compose-Default
(alle Services). Multi-Service-Form `u-boot logs svc1 svc2`
würde Subset-Filter erlauben (nicht alle, aber mehr als einer).

Spec [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) spricht von "Service" im Singular —
Multi-Service ist Spec-Erweiterung.

## Trigger

- **Real-World-Konsumenten-Bedarf** nach Per-Service-Subset
  (z. B. CI-Use-Case mit zwei korrelierten Services).
- **Compose-CLI-Parity-Druck**: Docker-Compose-CLI unterstützt
  Multi-Service direkt.

## Lösungs-Skizze (vorläufig)

`cobra.MaximumNArgs(1)` → `cobra.ArbitraryArgs` mit Per-Arg
`domain.NewServiceName`-Validation. Application-Layer
`LogsRequest.Service string` → `Services []string`. Adapter
`ComposeLogsOptions.Services` ist bereits Slice — kein
Driven-Port-Refactor.

## Spec-Bezug

- [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) (Logs anzeigen) — Erweiterung von Singular
  auf Plural Args.

## Closure-Notiz

**Geliefert:** `u-boot logs [service…]` (`cobra.ArbitraryArgs`, Regex-Prüfung je Name, Deduplizierung); `LogsRequest.Service string` → `Services []string`; der Adapter reicht die Liste (war schon ein Slice) durch.

**Sensoren:** `make gates` grün; Docker-Integrationstest
`TestE2E_LHFAUP005_LogsFormatAndTimeRange` gegen echtes Compose
(Präfix, Zeitstempel, `--since 1h` behält / `--until 1h` leert die Boot-Zeilen);
Argv-Pin im Adapter-Test; CLI-Tests für Validierung und Exit-Codes.

**Doku:** `docs/user/cli-json-output.md` §6.8, Benutzerhandbuch §Logs, Beispiele, README, CHANGELOG.
