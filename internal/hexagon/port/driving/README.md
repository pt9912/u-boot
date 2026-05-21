# internal/hexagon/port/driving

Interfaces, über die u-boot **von außen angesprochen wird**
(`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/hexagon/application/`,
verwendet von Adaptern in `internal/adapter/driving/` (z. B.
CLI-Commands).

## Status

Noch leer; erste Inhalte mit M3-T2 (siehe
[`docs/plan/planning/in-progress/slice-m3-init-flow.md`](../../../../docs/plan/planning/in-progress/slice-m3-init-flow.md)).

## Geplante Inhalte (M3-T2 und später)

- `InitProjectUseCase` — M3-T2.
- `AddServiceUseCase`, `RemoveServiceUseCase` — M4/M5.
- `LifecycleUseCase`, `DoctorUseCase` — M4.
- `GenerateUseCase`, `ConfigUseCase` — V1+.

## Import-Regeln

Nur `internal/hexagon/domain` und Go-Standard-Library. **Nicht**
erlaubt: `internal/hexagon/application`, `internal/hexagon/port/driven`,
`internal/adapter/*`.
