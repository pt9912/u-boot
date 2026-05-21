# Slice V1: Logging-Port einführen

## Auslöser

Die `forbidigo`-Regel in `.golangci.yml` verbietet `fmt.Print*` mit
der Begründung *„use log/slog (LH-QA-004); a project-specific logging
port may replace this once it exists"*. Der referenzierte „configured
logging port" existiert heute nicht — es ist prospektive Doku
(`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

u-boot bekommt einen eigenen Logging-Driven-Port in
`internal/hexagon/port/driven/logger/` (Interface) und einen
`slog`-basierten Adapter in `internal/adapter/driven/logger/`. Application
und Adapter loggen ausschließlich über den Port; `log/slog` selbst
wird nur im Adapter importiert.

## Akzeptanzkriterien

- `internal/hexagon/port/driven/logger/logger.go` definiert ein
  knappes Interface (z. B. `Logger` mit `Debug`/`Info`/`Warn`/`Error`,
  strukturierte Felder). Paketname `logger` (nicht `log`) vermeidet
  Kollision mit der Standard-Library.
- `internal/adapter/driven/logger/slog.go` implementiert das Interface
  via `log/slog` mit konfigurierbarem Handler (Text vs. JSON je nach
  CLI-Flag, `LH-FA-CLI-005`/`-007`).
- Wiring in `cmd/uboot/main.go` instantiiert den Adapter und injiziert
  ihn.
- `.golangci.yml` `forbidigo.msg` wird auf die finale Form aktualisiert:
  *„use the logging port at internal/hexagon/port/driven/logger
  (LH-QA-004)"*.
- `make lint` läuft grün.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf den Aufhebungs-Commit als gelöst markieren.

## Out of Scope

- Logging-Backend-Wahl jenseits von `log/slog` (z. B. zerolog, zap) —
  bleibt separater ADR-Wert, falls je nötig.
- Strukturierte Telemetrie / OTel — gehört zum OTel-Add-on-Slice
  (`LH-FA-ADD-004`).

## Bezug

- Auslösende Konfig: `.golangci.yml` `forbidigo.msg`.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `forbidigo.msg` referenziert nicht-existenten Port.
- Hängt von: M3 oder erstem Adapter-Slice (sobald
  `internal/adapter/driven/` ohnehin Inhalte bekommt).
