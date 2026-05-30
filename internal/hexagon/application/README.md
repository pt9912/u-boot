# internal/hexagon/application

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält
keine externe I/O (`LH-FA-ARCH-002`).

## Status

Stand M8 (MVP vollständig): sieben Use-Cases produktiv
verdrahtet (M3-T2 bis M8-T5), keine externe I/O, alle Ports
nil-tolerant via package-private `noop*`-Defaults.

## Inhalt

- `InitProjectService` — `port/driving.InitProjectUseCase`
  (`LH-FA-INIT-001..007`, M3). MVP-Closure-T1 ergänzt den
  `--devcontainer`-Flag (LH-AK-005): zwei zusätzliche Templates
  (`.devcontainer/devcontainer.json` + `Dockerfile`) durchlaufen
  die bestehende M3-`planFile`-Pipeline; `u-boot.yaml` bekommt
  `devcontainer.enabled: true`.
- `AddServiceService` — `port/driving.AddServiceUseCase`
  (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`, M5).
- `DoctorService` — `port/driving.DoctorUseCase`
  (`LH-FA-DIAG-001..004`, M4; 11 Checks). MVP-Closure-T2 ändert
  `compose.yaml.valid` no-services Severity Error → Warn
  (LH-AK-001-§2299-Konformität).
- `UpService` — `port/driving.UpUseCase` (`LH-FA-UP-001..003`, M6-T4).
  Polling-Loop mit `pollInterval=500ms` und `dialTimeout=300ms`,
  fail-safe `ContainerState`-Klassifikation (Dead-Allowlist,
  soft-Unknown, Restart-Loop-Counter mit Threshold 3),
  Healthcheck-dominanter Stabilisierungs-Vertrag mit TCP-Port-
  Probe als Warn-Diagnose (§141 / §968).
- `DownService` — `port/driving.DownUseCase` (`LH-FA-UP-004`, M6-T5).
  §T5-Truth-Table für den `--volumes`-Bestätigungs-Pfad (4 Zeilen
  × 2 Sub-Cases bei AssumeYes-und-NonInteractive); ruft
  `Confirmer.ConfirmRemoveVolumes` nur im interaktiven Pfad.
- `GenerateService` — `port/driving.GenerateUseCase`
  (`LH-FA-GEN-001..005`, M7). Vier Artefakt-Handler: env-example
  und readme über den shared `generateManagedFile`-Helper (4-State-
  Maschine mit Idempotenz-Pin), changelog mit konservativer
  User-Edit-Erkennung + `## [Unreleased]`-RepairedManual-Pfad
  (LH-AK-007), devcontainer mit atomarem Two-File-Plan +
  `forwardPorts`-Detection via Anti-Drift-Pin gegen die
  `Doctor::collectActiveServicePorts`-Quelle.
- `ConfigService` — `port/driving.ConfigUseCase`
  (`LH-FA-CONF-001..005`, M8). Get/Set/Show mit pfad-gesteuertem
  3-Element-Whitelist; Set durchläuft zweistufige Schema-
  Roundtrip-Validation (Struct-Unmarshal + Per-Pfad-Domain-Re-
  Validation) **vor** jedem WriteFile; `services.<svc>.enabled` ist
  Get-only (LH-FA-ADD-005-Lifecycle-Schutz).

Plus Helper:
- `parseComposePort` (M6-T4-fund) — pure 8-Syntax-Cases-Parser für
  Compose-`ports:`-Array-Elemente; nicht-probebare Formen
  (UDP/Range/Unknown) returnen `probable=false` für Warn-Diagnose-
  Pfad.
- `collectActiveServicePorts` / `activeServiceNames` (M5-T7, in
  `doctor.go`) — werden seit M7-T5 auch vom `GenerateService` für
  `forwardPorts` mitbenutzt (single source of truth, Anti-Drift-
  Pin).
- `managedblock/` — Marker-basiertes YAML-Blocksetting (M3+M5+M7);
  3 Marker-Stile (Hash / HTMLComment / DoubleSlash).
- `templates/` — Embedded Go-Templates für `u-boot init`/`add`/
  `generate` (M3 + M5-T4c + M7-T5).

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port`. **Nicht** erlaubt:
`internal/adapter/*`, externe I/O-Libraries.
