# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
(`LH-FA-ARCH-002`).

## Status

Stand M8 (MVP vollständig): sieben Cobra-Subcommands plus die
LH-FA-CLI-005A-Modi-Flags am Root verdrahtet. Wiring erfolgt
zentral in `cmd/uboot/main.go`; der Constructor `cli.New` trägt
seit M8-T5 sieben Use-Case-Parameter (Init/Doctor/Add/Up/Down/
Generate/Config) plus den optionalen `WithLogLevel`-Hook.

## Inhalt

- `cli/` — Cobra-basierte Commands:
  - `init [name] [--devcontainer]` (M3 + MVP-Closure-T1) —
    `LH-FA-INIT-001..007` + `LH-AK-005`.
  - `doctor` (M4) — `LH-FA-DIAG-001..004`, `--strict`.
  - `add <service>` (M5) — `LH-FA-ADD-001..002`/`-005`.
  - `up [--timeout <sek>]` (M6) — `LH-FA-UP-001..003`.
  - `down [--volumes]` (M6) — `LH-FA-UP-004`.
  - `generate <artifact>` (M7) — `LH-FA-GEN-001..005`.
  - `config [get/set]` (M8) — `LH-FA-CONF-001..005`. Drei
    Cobra-Shapes: parent-`config` läuft Show via `Args: NoArgs +
    RunE`, `get`/`set` sind Children mit `ExactArgs(1)` /
    `ExactArgs(2)`.
  - Persistente Root-Flags: `--yes`/`--no-interactive`
    (LH-FA-CLI-005A), `--quiet`/`--verbose`/`--debug`
    (LH-FA-CLI-005). `--yes` gilt explizit auch für
    `down --volumes` (Spec §237). Die Verbosity-Flags steuern
    seit [`slice-followup-verbosity-wiring`](../../../docs/plan/planning/done/slice-followup-verbosity-wiring.md)
    zusätzlich den `slog.Level` zur Laufzeit (`PersistentPreRunE`
    flippt ein per `WithLogLevel` injiziertes `*slog.LevelVar`):
    `--debug`/`--verbose` → `Debug`, `--quiet` → `Warn`, sonst
    `Info`.
- `cli/statusview.go` — tabwriter-basierter LH-FA-UP-003-Status-
  Renderer plus Down-Success-Renderer mit asymmetrischem
  `--quiet`-Vertrag (up suppress't Tabelle+Diagnostics, down
  suppress't die Erfolgsmeldung; Progress-Stream auf stderr
  bleibt in beiden Fällen unangetastet — LH-NFA-PERF-002).
- `cli/cli.go` — `ExitCode`-Klassifikation per LH-FA-CLI-006:
  - 0 Erfolg
  - 2 CLI-Validation (`--timeout=-1`, `--yes`+`--no-interactive`,
    unbekannte Flags, `generate <unknown-artifact>`,
    `config get/set` mit zu wenigen Args)
  - 10 fachliche Validation (Projekt nicht initialisiert,
    `compose.yaml` fehlt, destruktive Bestätigung verweigert,
    `generate`-Managed-Block-Konflikt + Schema-Konflikt,
    `config`-Pfad/Wert/Schema-/NotSet-Fehler). Die config-
    Familie ist über `isConfigValidationError` carve-outed,
    damit `isValidationError` unter dem gocyclo-Limit bleibt.
  - 11 Umgebung (Docker/Compose-Plugin via Pre-Probe)
  - 12 Ausführung (Compose-Runtime, Stabilisierungs-Timeout)
  - 14 technischer FS-Fehler (Backup-Suffix exhausted,
    `ErrGenerateFileSystem`, `ErrConfigFileSystem`)

## Geplante Erweiterungen

- `template`, `logs`, `--json`-Output (V1).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driving` und
externe Libraries (z. B. Cobra). **Nicht** erlaubt:
`internal/adapter/driven` direkt und `internal/hexagon/application`
— das Wiring erfolgt in `cmd/uboot/`.
