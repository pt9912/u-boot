# internal/hexagon/port/driving

Interfaces, über die u-boot **von außen angesprochen wird**
(`LH-FA-ARCH-002`).

Implementiert von Strukturen in `internal/hexagon/application/`,
verwendet von Adaptern in `internal/adapter/driving/` (z. B. CLI-Commands).

Geplante Inhalte (M3+):

- `InitProjectUseCase`, `AddServiceUseCase`, `RemoveServiceUseCase`,
  `LifecycleUseCase`, `DoctorUseCase`, `GenerateUseCase`,
  `ConfigUseCase`.

Import-Regeln: nur `internal/hexagon/domain` und Go-Standard-Library.
**Nicht** erlaubt: `internal/hexagon/application`,
`internal/hexagon/port/driven`, `internal/adapter/*`.
