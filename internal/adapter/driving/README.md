# internal/adapter/driving

Konkrete Driver — Einstiegspunkte aus der Außenwelt
(`LH-FA-ARCH-002`).

Geplante Inhalte (M3+):

- `cli/` – Cobra-Commands `init`, `add`, `remove`, `up`, `down`,
  `doctor`, `logs`, `generate`, `config`, `template`. Bindet
  CLI-Flags an Driving-Port-Aufrufe.

Import-Regeln: `internal/hexagon/domain`, `internal/hexagon/port/driving`
und externe Libraries (z. B. Cobra). **Nicht** erlaubt:
`internal/adapter/driven` direkt — das Wiring erfolgt in `cmd/uboot/`.
