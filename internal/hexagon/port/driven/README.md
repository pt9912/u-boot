# internal/hexagon/port/driven

Interfaces, über die `internal/hexagon/application` **externe Systeme
nutzt** (`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/adapter/driven/`.

## Aktueller Inhalt (M3-T1)

- `FileSystem` — `Exists`, `ReadFile`, `WriteFile`, `MkdirAll`,
  `Rename`, `ReadDir`. Genug für `LH-FA-INIT-003`/`-005`.
- `YAMLCodec` — `Marshal`, `Unmarshal`. Schlanke Surface für
  `LH-FA-CONF-001..003`; managed-block-aware Edits
  (`LH-SA-FILE-002`) folgen.
- `Git` — `IsRepository`, `Init`. Beide mit `context.Context`
  als erstem Parameter (Adapter shellt zum `git`-Binary, das
  blockieren kann; Application-Layer muss cancellable bleiben).
- `Clock` — `Now`. Ohne Context, weil die Implementierung
  non-blocking ist (Convention im Paket-Doc).

## Geplante Erweiterungen

- `DockerEngine` — `Up`, `Down`, `Ps`, `Logs`, `Exec` für M4
  (`LH-FA-UP-001..004`, `LH-FA-UP-005`, `LH-SA-DOCKER-001/-002`).
- Erweiterungen am `Git`-Interface, falls Operationen jenseits
  `Init`/`IsRepository` gebraucht werden.

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driving`,
`internal/adapter/*`.
