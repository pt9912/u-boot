# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
(`LH-FA-ARCH-002`).

## Status

Ein Cobra-Command je Subkommando, dazu die
[`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)-Modi-Flags
am Root. Wiring erfolgt zentral in `cmd/uboot/main.go`; der Constructor
`cli.New` nimmt je Use-Case-Familie einen Parameter plus den optionalen
`WithLogLevel`-Hook.

## Inhalt

- `cli/` — Cobra-basierte Commands:
  - `init [name] [--devcontainer]` —
    `LH-FA-INIT-001..007` + `LH-AK-005`.
  - `doctor` — `LH-FA-DIAG-001..004`, `--strict`.
  - `add <service>` — `LH-FA-ADD-001..002`/`-005`.
  - `up [--timeout <sek>]` — `LH-FA-UP-001..003`.
  - `down [--volumes]` — `LH-FA-UP-004`.
  - `remove <service>` — `LH-FA-ADD-007`, inkl. `--purge`.
  - `logs [service]` — `LH-FA-UP-005`.
  - `generate <artifact>` — `LH-FA-GEN-001..005`.
  - `template list` — `LH-FA-TPL-001..004`.
  - `config [get/set]` — `LH-FA-CONF-001..005`. Drei
    Cobra-Shapes: parent-`config` läuft Show via `Args: NoArgs +
    RunE`, `get`/`set` sind Children mit `ExactArgs(1)` /
    `ExactArgs(2)`.
  - Persistente Root-Flags: `--yes`/`--no-interactive`
    ([`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)), `--quiet`/`--verbose`/`--debug`
    ([`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)). `--yes` gilt explizit auch für
    `down --volumes` (Spec §237). Die Verbosity-Flags steuern
    seit [`slice-followup-verbosity-wiring`](../../../docs/plan/planning/done/slice-followup-verbosity-wiring.md)
    zusätzlich den `slog.Level` zur Laufzeit (`PersistentPreRunE`
    flippt ein per `WithLogLevel` injiziertes `*slog.LevelVar`):
    `--debug`/`--verbose` → `Debug`, `--quiet` → `Warn`, sonst
    `Info`.
- `cli/statusview.go` — tabwriter-basierter [`LH-FA-UP-003`](../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen)-Status-
  Renderer plus Down-Success-Renderer mit asymmetrischem
  `--quiet`-Vertrag (up suppress't Tabelle+Diagnostics, down
  suppress't die Erfolgsmeldung; Progress-Stream auf stderr
  bleibt in beiden Fällen unangetastet — [`LH-NFA-PERF-002`](../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)).
- `cli/cli.go` — `ExitCode`-Klassifikation per [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes):
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
