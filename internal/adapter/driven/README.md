# internal/adapter/driven

Konkrete externe Adapter — Implementierungen der Driven-Ports aus
`internal/hexagon/port/driven/` (`LH-FA-ARCH-002`).

Jeder Adapter pinnt sein Port-Interface per
`var _ driven.X = (*Adapter)(nil)` im Production-Code; ein Drift
zwischen Port und Adapter bricht damit den Package-Build, nicht erst
den Test-Build.

## Aktueller Inhalt (M3-T1)

- `fs/` — `FileSystem`-Adapter via Go-Standard-Library
  (`os.ReadFile`/`WriteFile`/`MkdirAll`/`Rename`/`ReadDir`).
- `yaml/` — `YAMLCodec`-Adapter via `gopkg.in/yaml.v3` (erste
  externe Modul-Dependency in `go.mod`).
- `git/` — `Git`-Adapter shellt zum `git`-Binary via `os/exec`
  mit Context-Support. `IsRepository` klassifiziert nur exit
  code 128 als "no repo"; alle anderen Fehler propagieren.
- `clock/` — `Clock`-Adapter mit `time.Now()` in UTC
  (projektweiter Default für deterministische Doku-Timestamps).

## Geplante Erweiterungen

- `docker/` — `DockerEngine`-Adapter via `os/exec docker compose`
  (`LH-SA-DOCKER-001`, `LH-SA-DOCKER-002`) — M4.
- Test-Strecke für `docker/` mit Build-Tag `//go:build docker`
  (siehe
  [`docs/plan/planning/open/slice-v1-docker-integrationstests.md`](../../../docs/plan/planning/open/slice-v1-docker-integrationstests.md)).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driven` und externe
Libraries (z. B. `yaml.v3`). **Nicht** erlaubt:
`internal/hexagon/application`, `internal/adapter/driving`.
