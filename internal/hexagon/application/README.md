# internal/hexagon/application

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält
keine externe I/O (`LH-FA-ARCH-002`).

## Status

Ein Service je Use-Case-Familie, keine externe I/O, alle Ports
nil-tolerant via package-private `noop*`-Defaults.

## Inhalt

- `InitProjectService` — `port/driving.InitProjectUseCase`
  (`LH-FA-INIT-001..007`). Der `--devcontainer`-Flag ([`LH-AK-005`](../../../spec/lastenheft.md#lh-ak-005--devcontainer-flow)): zwei zusätzliche Templates
  (`.devcontainer/devcontainer.json` + `Dockerfile`) durchlaufen
  dieselbe `planFile`-Pipeline; `u-boot.yaml` bekommt
  `devcontainer.enabled: true`.
- `AddServiceService` — `port/driving.AddServiceUseCase`
  (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`).
- `DoctorService` — `port/driving.DoctorUseCase`
  (`LH-FA-DIAG-001..004`; 11 Checks). `compose.yaml.valid` stuft den
  no-services-Fall als Warn statt Error ein
  ([`LH-AK-001`](../../../spec/lastenheft.md#lh-ak-001--minimaler-init-flow)-§2299-Konformität).
- `UpService` — `port/driving.UpUseCase` (`LH-FA-UP-001..003`).
  Polling-Loop mit `pollInterval=500ms` und `dialTimeout=300ms`,
  fail-safe `ContainerState`-Klassifikation (Dead-Allowlist,
  soft-Unknown, Restart-Loop-Counter mit Threshold 3),
  Healthcheck-dominanter Stabilisierungs-Vertrag mit TCP-Port-
  Probe als Warn-Diagnose (§141 / §968).
- `DownService` — `port/driving.DownUseCase` (`LH-FA-UP-004`).
  §T5-Truth-Table für den `--volumes`-Bestätigungs-Pfad (4 Zeilen
  × 2 Sub-Cases bei AssumeYes-und-NonInteractive); ruft
  `Confirmer.ConfirmRemoveVolumes` nur im interaktiven Pfad.
- `RemoveServiceService` — `port/driving.RemoveServiceUseCase`
  ([`LH-FA-ADD-007`](../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)); der
  destruktive `--purge`-Pfad laeuft nur ueber das Bestaetigungs-Gate.
- `LogsService` — `port/driving.LogsUseCase`
  ([`LH-FA-UP-005`](../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen)).
- `TemplateListService` / `TemplateInitService` —
  `port/driving.TemplateListUseCase` bzw. `TemplateInitUseCase`
  (`LH-FA-TPL-001..004`). Der Init-Pfad wird vom `InitProjectService`
  konsumiert, nicht vom Adapter direkt: Er teilt Verzeichnis-,
  Erkennungs- und Git-Logik mit dem Standard-Pfad und uebernimmt nur
  das Datei-Rendering.
- `GenerateService` — `port/driving.GenerateUseCase`
  (`LH-FA-GEN-001..005`). Vier Artefakt-Handler: env-example
  und readme über den shared `generateManagedFile`-Helper (4-State-
  Maschine mit Idempotenz-Pin), changelog mit konservativer
  User-Edit-Erkennung + `## [Unreleased]`-RepairedManual-Pfad
  ([`LH-AK-007`](../../../spec/lastenheft.md#lh-ak-007--changelog-generator)), devcontainer mit atomarem Two-File-Plan +
  `forwardPorts`-Detection via Anti-Drift-Pin gegen die
  `Doctor::collectActiveServicePorts`-Quelle.
- `ConfigService` — `port/driving.ConfigUseCase`
  (`LH-FA-CONF-001..005`). Get/Set/Show mit pfad-gesteuertem
  3-Element-Whitelist; Set durchläuft zweistufige Schema-
  Roundtrip-Validation (Struct-Unmarshal + Per-Pfad-Domain-Re-
  Validation) **vor** jedem WriteFile; `services.<svc>.enabled` ist
  Get-only ([`LH-FA-ADD-005`](../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern)-Lifecycle-Schutz).

Plus Helper:
- `parseComposePort` — pure 8-Syntax-Cases-Parser für
  Compose-`ports:`-Array-Elemente; nicht-probebare Formen
  (UDP/Range/Unknown) returnen `probable=false` für Warn-Diagnose-
  Pfad.
- `collectActiveServicePorts` / `activeServiceNames` (in `doctor.go`) — werden auch vom `GenerateService` für
  `forwardPorts` mitbenutzt (single source of truth, Anti-Drift-
  Pin).
- `managedblock/` — Marker-basiertes YAML-Blocksetting;
  3 Marker-Stile (Hash / HTMLComment / DoubleSlash).
- `templates/` — Embedded Go-Templates für `u-boot init`/`add`/
  `generate`.

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port`. **Nicht** erlaubt:
`internal/adapter/*`, externe I/O-Libraries.
