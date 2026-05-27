# internal/hexagon/port/driven

Interfaces, über die `internal/hexagon/application` **externe Systeme
nutzt** (`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/adapter/driven/`.

## Aktueller Inhalt

- `FileSystem` — `Exists`, `ReadFile`, `WriteFile`,
  `WriteFileExclusive`, `Mkdir`, `MkdirAll`, `Rename`, `ReadDir`,
  `Lstat`, `RemoveAll`, `Copy`, `CopyExclusive`. Streaming-Primitive
  (`Copy`/`CopyExclusive`) seit `slice-v1-backup-streaming-copy`;
  `interfacebloat`-Limit für diese eine Schnittstelle bewusst
  aufgeweicht (siehe `carveouts.md`).
- `YAMLCodec` — `Marshal`, `Unmarshal`. Schlanke Surface für
  `LH-FA-CONF-001..003`.
- `Git` — `IsRepository`, `Init`, `Version`. Alle mit
  `context.Context` als erstem Parameter (Adapter shellt zum
  `git`-Binary, das blockieren kann; Application-Layer muss
  cancellable bleiben).
- `Clock` — `Now`. Ohne Context, weil die Implementierung
  non-blocking ist (Convention im Paket-Doc).
- `ProgressPort` — `AffectedFiles(baseDir, rows)` für die
  LH-FA-INIT-005-§609-betroffenen-Pfade-Reports vor jedem
  Re-Init-Write. Presentation lebt im Adapter.
- `Confirmer` — `ConfirmTreatAsExisting(ctx, baseDir, indicators)`
  für die LH-FA-INIT-004-Soft-Existing-Detection-Prompts.
- `Logger` — `Debug`/`Info`/`Warn`/`Error` (variadisch, slog-konform)
  als LH-QA-004-Logging-Port. Production-Adapter slog-basiert.
- `DockerProbe` — `Version`/`Info`/`ComposeVersion` für die
  LH-FA-DIAG-002-Probes (M4-doctor T3). Read-only; getrennt vom
  geplanten `DockerEngine`-Port.

## Geplante Erweiterungen

- `DockerEngine` — `Up`, `Down`, `Ps`, `Logs`, `Exec` für **M6**
  (`LH-FA-UP-001..004`, `LH-FA-UP-005`, `LH-SA-DOCKER-001/-002`).
  Bewusst separater Port von `DockerProbe`: state-mutierend vs.
  read-only.

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driving`,
`internal/adapter/*`.
