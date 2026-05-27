# Slice M4: `u-boot doctor`-Flow

> **Status:** In progress
> **DoD:** T1 ✅ `bcf684f` / T2 ✅ `19423d9` / T3 offen / T4 offen / T5 offen / T6 offen / T7 offen

## Auslöser

Nach M3 (`u-boot init`) ist das zweite MVP-Subkommando dran:
`u-boot doctor` liefert die Diagnose-Funktion aus `LH-FA-DIAG-001..004`.

Pflichten der Spec:

- **`LH-FA-DIAG-001`** Doctor-Befehl muss existieren.
- **`LH-FA-DIAG-002`** Mindest-Checks: Docker (≥24.0.0), Docker-Erreichbarkeit,
  Docker Compose (≥2.20.0), Git, Schreibrechte im Projektverzeichnis,
  `compose.yaml`-Gültigkeit, `u-boot.yaml`-Gültigkeit, Devcontainer-
  Konsistenz (.devcontainer/devcontainer.json + forwardPorts).
- **`LH-FA-DIAG-003`** Severity-Klassifikation `ok`/`warn`/`error` mit
  Exit-Code-Mapping (`error` → ≠0). `--strict` macht `warn`
  ebenfalls zu ≠0.
- **`LH-FA-DIAG-004`** Reparaturhinweise bei Problemen.

Plus aus angrenzenden Spec-Punkten:

- **`LH-FA-CLI-005`** Verbosity-Flags (`--quiet`, `--verbose`, `--debug`)
  — werden mit doctor erstmals load-bearing.
- **`LH-NFA-USE-004`** JSON-Output (`--json`) — kann auf V1 vertagt
  werden, wenn der Text-Output zuerst stabilisiert wird.

## Vorbereitende Slices (alle bereits abgeschlossen)

- [`slice-m4-soft-existing-detection`](../done/slice-m4-soft-existing-detection.md)
  — `Confirmer`-Port, `LH-FA-INIT-004` aufgelöst.
- [`slice-m4-logging-port`](../done/slice-m4-logging-port.md) — `Logger`-
  Port + slog-Adapter; doctor nutzt den Port intensiv (jeder Check
  emittiert Debug/Info-Events).

## Tranchen-Schnitt

Jede Tranche eigener Commit, je grün durch alle 4 Gates (lint + test +
coverage-gate + docs-check):

1. **T1 — Domain + Driving-Port + Tests.**
   - `internal/hexagon/domain/diagnostic.go`: `Severity`-Enum
     (`SeverityOK`/`SeverityWarn`/`SeverityError`), `Diagnostic`-
     Struct (`{ID, Severity, Message, Hint}`), `DiagnosticReport`
     (`Items []Diagnostic`) mit Helper-Methoden (`MaxSeverity()`,
     `HasErrors()`, `HasWarnings()`).
   - `internal/hexagon/port/driving/doctor.go`: `DoctorRequest`
     (`BaseDir`), `DoctorResponse` (`Report DiagnosticReport`),
     `DoctorUseCase`-Interface mit `Check(ctx, req)`.
   - Domain-Tests mit Tabellen-Tests für die Severity-Helpers.

2. **T2 — Application-Service-Skeleton + Write-Perms-Check.**
   - `internal/hexagon/application/doctor.go`: `DoctorService` mit
     `fs driven.FileSystem` + `logger driven.Logger`. `Check(ctx, req)`
     aggregiert eine Liste interner Check-Funktionen.
   - Erster Check `checkWritePermissions`: legt ein Sentinel-File
     via `WriteFileExclusive` + `RemoveAll` an; klassifiziert
     `permission denied` → `error`, andere Fehler → `error`, Erfolg →
     `ok`. Repair-Hint auf `chmod`/Owner.
   - Tests mit `fakeFS` (success-Pfad + write-error-Pfad).

3. **T3 — Externe-Binär-Probes (Git availability + Docker + Compose).**
   - `Git`-Port um `Version(ctx) (string, error)` erweitert; Adapter
     ruft `git --version` via `os/exec`; Fake gibt vorkonfigurierten
     Wert/Fehler zurück.
   - Neuer Driven-Port `DockerProbe`:
     `Version(ctx) (string, error)`, `Info(ctx) error`, ggf.
     `ComposeVersion(ctx) (string, error)`.
   - Adapter `internal/adapter/driven/docker/` mit `os/exec docker
     version --format`-Auswertung + semver-Min-Check (Docker
     `≥24.0.0`, Compose `≥2.20.0`).
   - Drei neue Checks im DoctorService.
   - Fixture-Tests mit Fake-Probe für die semver-Schwellen.

4. **T4 — `u-boot.yaml`-Validierung.**
   - `u-boot.yaml`-Schema-Check: Existenz, `schemaVersion: 1`,
     `project.name`-Format (regex aus `LH-FA-INIT-006`). YAML-Codec-
     Port reicht; nutzt den `Unmarshal`-Pfad.
   - Klassifikation: fehlende Datei → `warn` (kein u-boot-Projekt),
     ungültige Datei → `error`.

5. **T5 — `compose.yaml`-Validierung.**
   - Parse-Success-Check; minimal Top-Level-Shape (`services:`-Key
     vorhanden). Tiefere Compose-Schema-Validierung wäre eigener
     Slice.

6. **T6 — Devcontainer-Validierung.**
   - JSONC-Parser für `.devcontainer/devcontainer.json` (Go hat
     keine stdlib-JSONC; entweder mit eigenem Vor-Strip oder einer
     externen Lib — Entscheidung im T6-Slice).
   - Mindestkompat-Checks (`name` gesetzt; `image` ODER `build`
     vorhanden).
   - `forwardPorts`-Konsistenz vs. aktivierten Services in
     `u-boot.yaml`.
   - Bedingte Logik je nach `devcontainer.enabled` in u-boot.yaml.

7. **T7 — CLI-Subcommand + Output-Layer + --strict.**
   - `internal/adapter/driving/cli/doctor.go`: Cobra-Subcommand
     `doctor` mit `--strict`-Flag.
   - Severity-Output mit Glyphen/Labels, Repair-Hints, Sortierung
     nach Severity dann ID.
   - `--strict`: `Warn` → Exit-Code ≠0 (sonst `0`).
   - Verbosity-Flags `--quiet`/`--verbose`/`--debug` werden im selben
     T7 oder einem T8 verdrahtet — Entscheidung beim Schreiben.
   - JSON-Output (`--json`, `LH-NFA-USE-004`) bewusst out-of-scope —
     eigener V1-Slice, sobald der Text-Output stabilisiert ist.

## Akzeptanzkriterien (Slice-Level)

- `LH-AK-001` doctor-Anteil: `mkdir demo && cd demo && u-boot init &&
  u-boot doctor` läuft grün (alle Checks `ok`).
- `LH-FA-DIAG-001..004` abgehakt.
- `make gates` grün.
- Carveouts in `carveouts.md` für M4-Folgepunkte bekommen Slice-
  Pläne (z. B. JSON-Output, falls als V1 vertagt).

## Out of Scope

- **`--json`-Output (`LH-NFA-USE-004`)**: separater V1-Slice nach M4-
  Text-Stabilisierung.
- **Tiefe Compose-Schema-Validierung**: nur Parse-Success +
  Top-Level-Shape im T5.
- **Network-Reachability** (z. B. Registry-Probes): nicht in
  LH-FA-DIAG-002, eigener Slice falls je gewollt.
- **Auto-Repair**: doctor diagnostiziert, repariert nicht.
  Reparatur-Hints sind Text, keine Aktionen.

## Bezug

- Auslösende Spec: `LH-FA-DIAG-001..004` (`spec/lastenheft.md` §4.7),
  `LH-FA-CLI-005` (Verbosity, Mindest-CLI-Layer).
- Vorgänger: [`slice-m4-soft-existing-detection`](../done/slice-m4-soft-existing-detection.md),
  [`slice-m4-logging-port`](../done/slice-m4-logging-port.md).
- Nachfolger: M5 (`u-boot add postgres`).
