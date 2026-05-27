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
- `clock/` — `Clock`-Adapter mit `time.Now()` in UTC
  (projektweiter Default für deterministische Doku-Timestamps).
- `progress/` — `ProgressPort`-Adapter (Text-Output für
  `LH-FA-INIT-005`-§609-Reports).
- `confirm/` — `Confirmer`-Adapter (`bufio.Scanner` über stdin,
  Prompt auf stderr; Default `[y/N]`).
- `logger/` — `Logger`-Adapter via `log/slog` (Text- und
  JSON-Format).
- `docker/` — `DockerProbe`-Adapter via `os/exec docker version`
  und `docker compose version --short` (read-only diagnostics
  für `LH-FA-DIAG-002`, M4-doctor). Bewusst NICHT der
  DockerEngine-Adapter — der kommt mit M6.

## Geplante Erweiterungen

- `DockerEngine`-Adapter (`Up`/`Down`/`Ps`/`Logs`/`Exec`) via
  `os/exec docker compose` für **M6** (`LH-SA-DOCKER-001`,
  `LH-SA-DOCKER-002`). Wird im selben `docker/`-Verzeichnis
  als zweiter Adapter-Typ landen oder ein eigenes Package
  bekommen — Entscheidung mit dem M6-Slice.
- Test-Strecke für `docker/` (Daemon-Tests, Compose-Smoke) mit
  Build-Tag `//go:build docker` (siehe
  [`docs/plan/planning/open/slice-m6-docker-integrationstests.md`](../../../docs/plan/planning/open/slice-m6-docker-integrationstests.md)).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driven` und externe
Libraries (z. B. `yaml.v3`). **Nicht** erlaubt:
`internal/hexagon/application`, `internal/adapter/driving`.
