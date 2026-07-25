# internal/hexagon/port/driving

Interfaces, über die u-boot **von außen angesprochen wird**
([`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)).

Implementiert von Strukturen in `internal/hexagon/application/`,
verwendet von Adaptern in `internal/adapter/driving/` (z. B.
CLI-Commands).

## Status

Je Use-Case-Familie ein Interface, jedes mit dediziertem
`*Request` / `*Response`-Paar plus narrow-scoped Sentinels (`Err*`)
neben der Interface-Definition.

## Inhalt

- `InitProjectUseCase` — [`LH-FA-INIT-001`](../../../../spec/lastenheft.md#lh-fa-init-001--neues-projekt-initialisieren)..[`LH-FA-INIT-007`](../../../../spec/lastenheft.md#lh-fa-init-007--git-repository-initialisierung). Sentinels:
  `ErrProjectExists`, `ErrFileExists`, `ErrBaseDirMissing`,
  `ErrBackupSourceMissing`, `ErrBackupSuffixExhausted`,
  `ErrBackupUnsupportedKind`, `ErrForceRequiresBackup`.
  `ErrProjectNotInitialized` lebt hier und wird von mehreren
  Use-Case-Familien mitgenutzt.
- `DoctorUseCase` — [`LH-FA-DIAG-001`](../../../../spec/lastenheft.md#lh-fa-diag-001--doctor-befehl)..[`LH-FA-DIAG-004`](../../../../spec/lastenheft.md#lh-fa-diag-004--reparaturhinweise). Keine Sentinels
  (Befunde sind im `domain.DiagnosticReport`).
- `AddServiceUseCase` — [`LH-FA-ADD-001`](../../../../spec/lastenheft.md#lh-fa-add-001--add-on-befehl)..[`LH-FA-ADD-002`](../../../../spec/lastenheft.md#lh-fa-add-002--postgresql-hinzufügen), [`LH-FA-ADD-005`](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern).
  Sentinels: `ErrServiceUnsupported`, `ErrServiceInconsistent`.
- `UpUseCase` — [`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)..[`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003--startstatus-anzeigen). Sentinels:
  `ErrComposeFileMissing`, `ErrStabilizationTimeout`. Plus die
  driven-port-Sentinels `driven.ErrDockerUnavailable` (CLI-Code 11)
  und `driven.ErrComposeRuntime` (CLI-Code 12) durchgereicht via
  `errors.Is`-Wrap-Vertrag.
- `DownUseCase` — [`LH-FA-UP-004`](../../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen). Sentinel:
  `ErrConfirmationRequired` (CLI-Code 10 für §254-destruktive-Aborts).

## Geplante Erweiterungen

- `GenerateUseCase` — [`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001--generate-befehl)..[`LH-FA-GEN-005`](../../../../spec/lastenheft.md#lh-fa-gen-005--idempotenz).
- `ConfigUseCase` — [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration)..[`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern).
- `LogsUseCase` (V1) — [`LH-FA-UP-005`](../../../../spec/lastenheft.md#lh-fa-up-005--logs-anzeigen).

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driven`,
`internal/adapter/*`.
