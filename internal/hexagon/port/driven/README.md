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
- `Clock` — `Now`, `Sleep(d)`. Sleep seit M6-T4-fund load-bearing
  für den UpService-Polling-Loop; ohne Context, weil
  Production-Implementation non-blocking-now bzw. delegierend an
  time.Sleep ist (Convention im Paket-Doc).
- `ProgressPort` — `AffectedFiles(baseDir, rows)` für die
  LH-FA-INIT-005-§609-betroffenen-Pfade-Reports vor jedem
  Re-Init-Write. Presentation lebt im Adapter.
- `Confirmer` — zwei narrow-scoped Methoden (M4-Konvention "explicit
  names per question"):
  - `ConfirmTreatAsExisting(ctx, baseDir, indicators)` für die
    LH-FA-INIT-004-Soft-Existing-Detection-Prompts (M4).
  - `ConfirmRemoveVolumes(ctx, baseDir)` für den
    LH-FA-CLI-005A-§254-destruktive-Confirmation-Pfad von
    `u-boot down --volumes` (M6-T5).
- `Logger` — `Debug`/`Info`/`Warn`/`Error` (variadisch, slog-konform)
  als LH-QA-004-Logging-Port. Production-Adapter slog-basiert.
- `DockerProbe` — `Version`/`Info`/`ComposeVersion` für die
  LH-FA-DIAG-002-Probes (M4-doctor T3). Read-only; bewusst
  getrennt vom state-mutierenden `DockerEngine`-Port.
- `DockerEngine` — `ComposeUp(ctx, dir, opts)`,
  `ComposeDown(ctx, dir, opts)`, `ComposePs(ctx, dir)` für M6
  (`LH-FA-UP-001..004`, `LH-SA-DOCKER-001/-002`). Per-Call-
  Preflight-Vertrag: LookPath + Daemon-Roundtrip + Compose-
  Plugin-Probe vor jedem echten Call. Sentinels
  `ErrDockerUnavailable` (CLI-Code 11) und `ErrComposeRuntime`
  (CLI-Code 12) — `errors.Is` survival pin durch kontextuelle
  Application-Wraps (slice §Sentinel-Schichtung). `ComposeUp`
  liefert seit M6-Closure-Review eine leere `ComposeUpResult` —
  kein Follow-up `compose ps` mehr, um den §970 fire-and-forget-
  Vertrag bei `--timeout=0` zu wahren.
- `NetProbe` — `DialTCP(ctx, host, port, timeout)` (M6-T3) für die
  Reachability-Probes des UpService-Polling-Loops
  (LH-FA-UP-001 §968). `ctx.Err()` hat Vorrang vor Net-Error
  (Adapter nutzt `net.Dialer.DialContext`).

## Geplante Erweiterungen

- `Logs`, `Exec` als V1-Erweiterung von `DockerEngine` für
  `LH-FA-UP-005` (`u-boot logs`).

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driving`,
`internal/adapter/*`.
