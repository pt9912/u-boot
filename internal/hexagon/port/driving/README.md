# internal/hexagon/port/driving

Interfaces, über die u-boot **von außen angesprochen wird**
(`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/hexagon/application/`,
verwendet von Adaptern in `internal/adapter/driving/` (z. B.
CLI-Commands).

## Status

Stand M6: fünf Use-Cases produktiv, jeder mit dediziertem
`*Request` / `*Response`-Paar plus narrow-scoped Sentinels (`Err*`)
neben der Interface-Definition.

## Inhalt

- `InitProjectUseCase` (M3) — `LH-FA-INIT-001..007`. Sentinels:
  `ErrProjectExists`, `ErrFileExists`, `ErrBaseDirMissing`,
  `ErrBackupSourceMissing`, `ErrBackupSuffixExhausted`,
  `ErrBackupUnsupportedKind`, `ErrForceRequiresBackup`.
  `ErrProjectNotInitialized` lebt hier und ist M5+M6 mitgenutzt.
- `DoctorUseCase` (M4) — `LH-FA-DIAG-001..004`. Keine Sentinels
  (Befunde sind im `domain.DiagnosticReport`).
- `AddServiceUseCase` (M5) — `LH-FA-ADD-001..002`, `LH-FA-ADD-005`.
  Sentinels: `ErrServiceUnsupported`, `ErrServiceInconsistent`.
- `UpUseCase` (M6-T1/T4) — `LH-FA-UP-001..003`. Sentinels:
  `ErrComposeFileMissing`, `ErrStabilizationTimeout`. Plus die
  driven-port-Sentinels `driven.ErrDockerUnavailable` (CLI-Code 11)
  und `driven.ErrComposeRuntime` (CLI-Code 12) durchgereicht via
  `errors.Is`-Wrap-Vertrag.
- `DownUseCase` (M6-T1/T5) — `LH-FA-UP-004`. Sentinel:
  `ErrConfirmationRequired` (CLI-Code 10 für §254-destruktive-Aborts).

## Geplante Erweiterungen

- `GenerateUseCase` (M7) — `LH-FA-GEN-001..005`.
- `ConfigUseCase` (M8) — `LH-FA-CONF-001..005`.
- `LogsUseCase` (V1) — `LH-FA-UP-005`.

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driven`,
`internal/adapter/*`.
