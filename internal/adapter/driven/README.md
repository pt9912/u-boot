# internal/adapter/driven

Konkrete externe Adapter — Implementierungen der Driven-Ports aus
`internal/hexagon/port/driven/` (`LH-FA-ARCH-002`).

Jeder Adapter pinnt sein Port-Interface per
`var _ driven.X = (*Adapter)(nil)` im Production-Code; ein Drift
zwischen Port und Adapter bricht damit den Package-Build, nicht erst
den Test-Build.

## Aktueller Inhalt

- `fs/` — `FileSystem`-Adapter via Go-Standard-Library
  (`os.ReadFile`/`WriteFile`/`MkdirAll`/`Rename`/`ReadDir`,
  Streaming-Copy via `os.Open` + `io.Copy`).
- `yaml/` — `YAMLCodec`-Adapter via `gopkg.in/yaml.v3` (erste
  externe Modul-Dependency in `go.mod`).
- `git/` — `Git`-Adapter shellt zum `git`-Binary via `os/exec`
  mit Context-Support. `IsRepository` klassifiziert nur exit
  code 128 als "no repo"; alle anderen Fehler propagieren.
  `Version` parsed das `git version <X.Y.Z>`-Format und liefert
  die bare semver.
- `clock/` — `Clock`-Adapter mit `time.Now()` in UTC und
  `time.Sleep`. Sleep ist load-bearing für
  den UpService-Polling-Loop; Tests ersetzen ihn durch einen
  manual-advance Fake (Vorgabe aus dem Slice-Plan: "kein reales
  time.Sleep in Tests").
- `progress/` — `ProgressPort`-Adapter (Text-Output für
  `LH-FA-INIT-005`-§609-Reports).
- `confirm/` — `Confirmer`-Adapter (`bufio.Scanner` über stdin,
  Prompt auf stderr; Default `[y/N]`).
  `ConfirmRemoveVolumes` deckt den destruktiven `down --volumes`-
  Pfad ([`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) §254).
- `logger/` — `Logger`-Adapter via `log/slog` (Text- und
  JSON-Format).
- `docker/` — zwei Adapter im selben Package:
  - `Probe` implementiert `DockerProbe` via `os/exec docker
    version` + `docker compose version --short` (read-only
    diagnostics für `LH-FA-DIAG-002`).
  - `Engine` implementiert `DockerEngine` via `docker
    compose -f compose.yaml up -d` / `down [-v]` / `ps --format
    json`. Jeder Call durchläuft den `preflight`-Pfad (LookPath
    + Probe.Info + Probe.ComposeVersion) zur deterministischen
    `ErrDockerUnavailable`-vs-`ErrComposeRuntime`-Klassifikation
    (CLI-Codes 11 vs. 12 per `LH-FA-CLI-006`).
- `netprobe/` — `NetProbe`-Adapter via
  `net.Dialer.DialContext` mit `Dialer.Timeout`. Wraps mit
  noctx-Lint-konformen `*net.OpError`, sodass `errors.Is`
  context.Canceled / context.DeadlineExceeded durchläuft.
  Directory bewusst `netprobe/` (statt `net/`), um stdlib-`net`-
  Aliasing-Konflikte am Aufrufer zu vermeiden.

## Test-Build-Tag-Pfad

- Build-Tagged Adapter-Integrationstests (`//go:build docker`)
  laufen via `make test-docker` gegen einen echten Docker-Daemon
  (mountet `/var/run/docker.sock`). Die Verhaltens-Pins für
  [`LH-AK-002`](../../../spec/lastenheft.md#lh-ak-002--postgresql-flow) und
  [`LH-NFA-PERF-002`](../../../spec/lastenheft.md#lh-nfa-perf-002--startzeit-abhängig-von-docker)
  liegen in der `e2e/`-Suite.

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driven` und externe
Libraries (z. B. `yaml.v3`). **Nicht** erlaubt:
`internal/hexagon/application`, `internal/adapter/driving`.
