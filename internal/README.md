# internal/

Nicht-exportierbare Go-Pakete für u-boot. Strukturiert nach dem
hexagonalen Architektur-Pattern (`LH-FA-ARCH-001..003`; Detail in
[`../spec/architecture.md`](../spec/architecture.md)).

## Layout

```
internal/
├── hexagon/
│   ├── domain/              # reine Datentypen, keine I/O
│   ├── application/         # Use-Cases; ruft nur Ports auf
│   └── port/
│       ├── driving/         # Interfaces, die CLI/HTTP konsumiert
│       └── driven/          # Interfaces, die Application nach außen ruft
└── adapter/
    ├── driving/             # konkrete Driver (cli/, …)
    └── driven/              # konkrete Adapter (docker/, fs/, yaml/, …)
```

## Status

Im MVP-Bootstrap leer (nur READMEs je Verzeichnis). Erste Pakete
entstehen mit dem ersten fachlichen Slice (M3 `u-boot init`).

Solange `./internal/...` keinen produktiven Code enthält, läuft das
Coverage-Gate im Bootstrap-Modus mit Schwellwert `0`
(`LH-FA-BUILD-008`). Mit dem ersten produktiven Paket wird die Schwelle
in einem Folge-Commit angehoben.

## Import-Regeln

Verbindliche Schicht-Regeln in [`LH-FA-ARCH-003`](../spec/lastenheft.md)
und [`../spec/architecture.md`](../spec/architecture.md). Enforcement
über `golangci-lint depguard` im `lint`-Stage; `//nolint:depguard` ist
verboten.
