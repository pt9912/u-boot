# Slice V1: `u-boot logs --no-log-prefix` / `--timestamps` Format-Flags

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Format-Flags-Carveout aus
> [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure trägt den Eintrag nach).

## Auslöser

`u-boot logs` heutige Surface ([slice-v1-logs](../done/slice-v1-logs.md) §T0-(d)) exponiert
nur `--follow` und `--tail`. Docker-Compose-CLI selbst trägt
zusätzlich `--no-log-prefix` (Service-Prefix unterdrücken) und
`--timestamps` (Zeitstempel ergänzen). Beide Flags würden in
CI-Use-Cases mit eigener Timestamp-Schicht Wert bringen
(Konsument hat schon Log-Aggregation und braucht keine
Prefix-Doppelung).

V1-Trade-off: Spec-treu ([`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) listet die zwei Flags
nicht) und schmale Surface > Compose-CLI-Parity. Die zwei
Flags wandern in diesen Folge-Slice.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Druck nach Format-Kontrolle**: CI-Use-Case mit
  eigener Timestamp-Schicht beschwert sich über
  Compose-Prefix-Doppelung.
- **`logs --json` NDJSON-Erweiterung**
  ([`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
  T0-(a) Option B): wenn der Slice (B) wählt und das
  Per-Line-`time`-Feld aus `--timestamps`-Form gefüllt werden
  soll, ist `--timestamps` Pflicht-Voraussetzung.

## Lösungs-Skizze (vorläufig)

Zwei neue lokale Flags im `cli/logs.go` Cobra-Command + Pass-
Through an `driven.DockerEngine.ComposeLogs`-Adapter. Keine
Application-Layer-Änderung (Stream-Pfad bleibt unverändert).

## Spec-Bezug

- [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) (Logs anzeigen) — Spec listet die zwei Flags
  nicht; Erweiterung ist Compose-CLI-Parity-Argument.
