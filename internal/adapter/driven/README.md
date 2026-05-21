# internal/adapter/driven

Konkrete externe Adapter — Implementierungen der Driven-Ports aus
`internal/hexagon/port/driven/` (`LH-FA-ARCH-002`).

Geplante Inhalte (M3+):

- `docker/` – Docker-Engine-Adapter via `os/exec docker compose`
  (`LH-SA-DOCKER-001`, `LH-SA-DOCKER-002`).
- `fs/` – Dateisystem-Adapter mit Backup-/Atomar-Rename-Logik
  (`LH-FA-INIT-005`).
- `yaml/` – YAML-Codec via `gopkg.in/yaml.v3`, managed-block-aware
  (`LH-SA-FILE-002`).
- `clock/` – Real-Time-Quelle; Test-Fakes leben in `_test.go`.
- `git/` – Git-Adapter via `os/exec git` (`LH-FA-INIT-007`).

Import-Regeln: `internal/hexagon/domain`,
`internal/hexagon/port/driven` und externe Libraries (z. B.
`yaml.v3`). **Nicht** erlaubt: `internal/hexagon/application`,
`internal/adapter/driving`.
