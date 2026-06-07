# Slice V1: `u-boot logs <svc1> <svc2>` Multi-Service-Filter

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Multi-Service-Filter-Carveout aus
> [`slice-v1-cli-json-dry-run-logs`](../in-progress/slice-v1-cli-json-dry-run-logs.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure trägt den Eintrag nach).

## Auslöser

Heutige `u-boot logs`-Surface (slice-v1-logs §AK) nutzt
`cobra.MaximumNArgs(1)` — Single-Service oder Compose-Default
(alle Services). Multi-Service-Form `u-boot logs svc1 svc2`
würde Subset-Filter erlauben (nicht alle, aber mehr als einer).

Spec `LH-FA-UP-005` spricht von "Service" im Singular —
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

- `LH-FA-UP-005` (Logs anzeigen) — Erweiterung von Singular
  auf Plural Args.
