# internal/hexagon/application

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält
keine externe I/O (`LH-FA-ARCH-002`).

Geplante Inhalte (M3+):

- `InitProjectService` – implementiert
  `port/driving.InitProjectUseCase` (`LH-FA-INIT-001..007`).
- `AddServiceService` – `port/driving.AddServiceUseCase`
  (`LH-FA-ADD-001..006`).
- `RunDoctorService` – `port/driving.DoctorUseCase`
  (`LH-FA-DIAG-001..004`).
- `UpService` / `DownService` – `port/driving.LifecycleUseCase`
  (`LH-FA-UP-001..004`).

Import-Regeln: `internal/hexagon/domain`, `internal/hexagon/port`.
**Nicht** erlaubt: `internal/adapter/*`, externe I/O-Libraries.
