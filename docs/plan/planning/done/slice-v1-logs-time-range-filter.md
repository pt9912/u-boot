# Slice V1: `u-boot logs --since` / `--until` Time-Range-Filter

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `4db9f90`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-logs`](../done/slice-v1-cli-json-dry-run-logs.md)
> §Out of Scope (Carveout-Eintrag entfernt).

## Auslöser

Heutige `u-boot logs`-Surface unterstützt nur `--tail` (Anzahl
Zeilen). Docker-Compose-CLI selbst trägt `--since` und
`--until` für Zeitbereichs-Filter (`--since="1h"`,
`--until="2026-06-07T12:00"`). Bei Post-Hoc-Debugging
(Container ist nicht mehr aktiv, aber Logs sind im
Compose-Verbund noch verfügbar) sind Zeitbereiche
informativer als reine Zeilen-Anzahl.

Spec [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) listet die zwei Flags nicht.

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

- [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) (Logs anzeigen) — Erweiterung um Time-Range-
  Format.

## Closure-Notiz

**Geliefert:** `--since`/`--until` mit CLI-Validierung (positive Dauer oder Zeitstempel in RFC 3339 / `YYYY-MM-DD[THH:MM[:SS]]`, sonst `ErrInvalidLogsTime` → Exit 2) und Pass-Through an Compose.

**Sensoren:** `make gates` grün; Docker-Integrationstest
`TestE2E_LHFAUP005_LogsFormatAndTimeRange` gegen echtes Compose
(Präfix, Zeitstempel, `--since 1h` behält / `--until 1h` leert die Boot-Zeilen);
Argv-Pin im Adapter-Test; CLI-Tests für Validierung und Exit-Codes.

**Doku:** `docs/user/cli-json-output.md` §6.8, Benutzerhandbuch §Logs, Beispiele, README, CHANGELOG.
