# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
(`LH-FA-ARCH-002`).

## Status

Stand M6: fünf Cobra-Subcommands plus die LH-FA-CLI-005A-Modi-Flags
am Root verdrahtet. Wiring erfolgt zentral in `cmd/uboot/main.go`.

## Inhalt

- `cli/` — Cobra-basierte Commands:
  - `init [name]` (M3) — `LH-FA-INIT-001..007`.
  - `doctor` (M4) — `LH-FA-DIAG-001..004`, `--strict`.
  - `add <service>` (M5) — `LH-FA-ADD-001..002`/`-005`.
  - `up [--timeout <sek>]` (M6) — `LH-FA-UP-001..003`.
  - `down [--volumes]` (M6) — `LH-FA-UP-004`.
  - Persistente Root-Flags: `--yes`/`--no-interactive`
    (LH-FA-CLI-005A), `--quiet`/`--verbose`/`--debug`
    (LH-FA-CLI-005). `--yes` gilt explizit auch für
    `down --volumes` (Spec §237).
- `cli/statusview.go` — tabwriter-basierter LH-FA-UP-003-Status-
  Renderer plus Down-Success-Renderer mit asymmetrischem
  `--quiet`-Vertrag (up suppress't Tabelle+Diagnostics, down
  suppress't die Erfolgsmeldung; Progress-Stream auf stderr
  bleibt in beiden Fällen unangetastet — LH-NFA-PERF-002).
- `cli/cli.go` — `ExitCode`-Klassifikation per LH-FA-CLI-006:
  - 0 Erfolg
  - 2 CLI-Validation (`--timeout=-1`, `--yes`+`--no-interactive`,
    unbekannte Flags)
  - 10 fachliche Validation (Projekt nicht initialisiert,
    `compose.yaml` fehlt, destruktive Bestätigung verweigert)
  - 11 Umgebung (Docker/Compose-Plugin via Pre-Probe)
  - 12 Ausführung (Compose-Runtime, Stabilisierungs-Timeout)
  - 14 technischer FS-Fehler (Backup-Suffix exhausted etc.)

## Geplante Erweiterungen

- `generate` (M7), `config` (M8), `template`, `logs` (V1).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driving` und
externe Libraries (z. B. Cobra). **Nicht** erlaubt:
`internal/adapter/driven` direkt und `internal/hexagon/application`
— das Wiring erfolgt in `cmd/uboot/`.
