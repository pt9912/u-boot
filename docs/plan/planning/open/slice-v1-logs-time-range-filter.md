# Slice V1: `u-boot logs --since` / `--until` Time-Range-Filter

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Time-Range-Filter-Carveout aus
> [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure trägt den Eintrag nach).

## Auslöser

Heutige `u-boot logs`-Surface unterstützt nur `--tail` (Anzahl
Zeilen). Docker-Compose-CLI selbst trägt `--since` und
`--until` für Zeitbereichs-Filter (`--since="1h"`,
`--until="2026-06-07T12:00"`). Bei Post-Hoc-Debugging
(Container ist nicht mehr aktiv, aber Logs sind im
Compose-Verbund noch verfügbar) sind Zeitbereiche
informativer als reine Zeilen-Anzahl.

Spec [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005-logs-anzeigen) listet die zwei Flags nicht.

## Trigger

- **Real-World-Druck** nach Time-Range-Filter (z. B.
  Post-Mortem-Analyse "was lief gestern zwischen 14:00 und
  15:00?").
- **Compose-CLI-Parity-Druck**: Compose unterstützt es bereits.

## Lösungs-Skizze (vorläufig)

Zwei neue lokale Flags `--since` + `--until` mit
`time.ParseDuration`/`time.Parse`-Validation. Pass-Through an
`driven.DockerEngine.ComposeLogs`-Adapter via
`ComposeLogsOptions.Since/Until`-Field-Erweiterung. Keine
Application-Layer-Logik-Änderung.

## Spec-Bezug

- [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005-logs-anzeigen) (Logs anzeigen) — Erweiterung um Time-Range-
  Format.
