# internal/hexagon/port/driven

Interfaces, über die `internal/hexagon/application` **externe Systeme
nutzt** (`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/adapter/driven/`.

Geplante Inhalte (M3+):

- `DockerEngine` – `Up`, `Down`, `Ps`, `Logs`, `Exec`
  (`LH-FA-UP-001..004`, `LH-FA-UP-005`).
- `FileSystem` – `ReadFile`, `WriteFile`, `Exists`, `Mkdir`, `Move`
  (mit Backup-Konvention aus `LH-FA-INIT-005`).
- `YAMLCodec` – `Marshal`, `Unmarshal`, managed-block-aware Edits
  (`LH-SA-FILE-002`).
- `Clock` – `Now`; testbar über Fake-Implementierung.
- `Git` – Repository-Initialisierung (`LH-FA-INIT-007`).

Import-Regeln: nur `internal/hexagon/domain` und Go-Standard-Library.
**Nicht** erlaubt: `internal/hexagon/application`,
`internal/hexagon/port/driving`, `internal/adapter/*`.
