# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
(`LH-FA-ARCH-002`).

## Status

Noch leer; erste Inhalte mit M3-T3 (siehe
[`docs/plan/planning/in-progress/slice-m3-init-flow.md`](../../../docs/plan/planning/in-progress/slice-m3-init-flow.md)).

## Geplante Inhalte (M3-T3 und später)

- `cli/` — Cobra-basierte Commands (`init`, später `add`, `remove`,
  `up`, `down`, `doctor`, `logs`, `generate`, `config`, `template`).
  Cobra wird mit M3-T3 als externe Modul-Dependency in `go.mod`
  landen und löst dann den ADR-0001-Folgepunkt CLI-Framework auf.

## Import-Regeln

`internal/hexagon/domain`, `internal/hexagon/port/driving` und
externe Libraries (z. B. Cobra). **Nicht** erlaubt:
`internal/adapter/driven` direkt und `internal/hexagon/application`
— das Wiring erfolgt in `cmd/uboot/`.
