# internal/hexagon/port/driven

Interfaces, über die `internal/hexagon/application` **externe Systeme
nutzt** ([`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)).

Implementiert von Strukturen in `internal/adapter/driven/`.

## Aktueller Inhalt

- `FileSystem` — `Exists`, `ReadFile`, `WriteFile`,
  `WriteFileExclusive`, `Mkdir`, `MkdirAll`, `Rename`, `ReadDir`,
  `Lstat`, `RemoveAll`, `Copy`, `CopyExclusive`. Streaming-Primitive
  (`Copy`/`CopyExclusive`) seit [`slice-v1-backup-streaming-copy`](../../../../docs/plan/planning/done/slice-v1-backup-streaming-copy.md);
  `interfacebloat`-Limit für diese eine Schnittstelle bewusst
  aufgeweicht (siehe `carveouts.md`).
- `YAMLCodec` — `Marshal`, `Unmarshal`. Schlanke Surface für
  [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)..[`LH-FA-CONF-003`](../../../../spec/lastenheft.md#lh-fa-conf-003--konfiguration-lesen).
- `Git` — `IsRepository`, `Init`, `Version`. Alle mit
  `context.Context` als erstem Parameter (Adapter shellt zum
  `git`-Binary, das blockieren kann; Application-Layer muss
  cancellable bleiben).
- `Clock` — `Now`, `Sleep(d)`. Sleep load-bearing
  für den UpService-Polling-Loop; ohne Context, weil
  Production-Implementation non-blocking-now bzw. delegierend an
  time.Sleep ist (Convention im Paket-Doc).
- `ProgressPort` — `AffectedFiles(baseDir, rows)` für die
  [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz)-§609-betroffenen-Pfade-Reports vor jedem
  Re-Init-Write. Presentation lebt im Adapter.
- `Confirmer` — zwei narrow-scoped Methoden (Konvention "explicit
  names per question"):
  - `ConfirmTreatAsExisting(ctx, baseDir, indicators)` für die
    [`LH-FA-INIT-004`](../../../../spec/lastenheft.md#lh-fa-init-004--bestehendes-projekt-erkennen)-Soft-Existing-Detection-Prompts.
  - `ConfirmRemoveVolumes(ctx, baseDir)` für den
    [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)-§254-destruktive-Confirmation-Pfad von
    `u-boot down --volumes`.
- `Logger` — `Debug`/`Info`/`Warn`/`Error` (variadisch, slog-konform)
  als [`LH-QA-004`](../../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil)-Logging-Port. Production-Adapter slog-basiert.
- `DockerProbe` — `Version`/`Info`/`ComposeVersion` für die
  [`LH-FA-DIAG-002`](../../../../spec/lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen)-Probes. Read-only; bewusst
  getrennt vom state-mutierenden `DockerEngine`-Port.
- `DockerEngine` — `ComposeUp(ctx, dir, opts)`,
  `ComposeDown(ctx, dir, opts)`, `ComposePs(ctx, dir)` für den
  Compose-Lifecycle
  ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)..[`LH-FA-UP-004`](../../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen), [`LH-SA-DOCKER-001`](../../../../spec/lastenheft.md#lh-sa-docker-001--docker-compose)/[`LH-SA-DOCKER-002`](../../../../spec/lastenheft.md#lh-sa-docker-002--containerstatus)). Per-Call-
  Preflight-Vertrag: LookPath + Daemon-Roundtrip + Compose-
  Plugin-Probe vor jedem echten Call. Sentinels
  `ErrDockerUnavailable` (CLI-Code 11) und `ErrComposeRuntime`
  (CLI-Code 12) — `errors.Is` survival pin durch kontextuelle
  Application-Wraps (slice §Sentinel-Schichtung). `ComposeUp`
  liefert eine leere `ComposeUpResult` —
  kein Follow-up `compose ps` mehr, um den §970 fire-and-forget-
  Vertrag bei `--timeout=0` zu wahren.
- `NetProbe` — `DialTCP(ctx, host, port, timeout)` für die
  Reachability-Probes des UpService-Polling-Loops
  ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten) §968). `ctx.Err()` hat Vorrang vor Net-Error
  (Adapter nutzt `net.Dialer.DialContext`).

## Geplante Erweiterungen

- `Logs`, `Exec` als V1-Erweiterung von `DockerEngine` für
  [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen) (`u-boot logs`).

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driving`,
`internal/adapter/*`.
