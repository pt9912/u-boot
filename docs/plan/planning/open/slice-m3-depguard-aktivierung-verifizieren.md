# Slice M3: depguard-Regelaktivierung verifizieren

## Auslöser

`LH-FA-ARCH-003` und `spec/architecture.md` §4 definieren acht
`depguard`-Regelblöcke für die hexagonalen Schicht-Imports. Sie sind
seit M2b in `.golangci.yml` aktiv, **matchen aber heute nichts**, weil
`./internal/...` keinen produktiven Code enthält. Die Regeln greifen
laut Konvention „automatisch mit dem ersten Paket pro Schicht" — das
ist heute nur Behauptung, kein Test (`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Mit dem ersten produktiven Slice (M3 `u-boot init`) liefern wir
mindestens je ein Paket in `internal/hexagon/domain/`,
`internal/hexagon/application/`, `internal/hexagon/port/{driving,driven}/`
und `internal/adapter/{driving,driven}/`. Jede `depguard`-Regel wird
einmal durch einen *bewusst falschen* Import verifiziert (Lint muss
fail), dann sauber gelassen (Lint muss grün sein).

## Akzeptanzkriterien

- Pro Regelblock dokumentierte Negativ-Verifikation (z. B. in einem
  Branch oder lokalen Edit), die zeigt, dass der Block den verbotenen
  Import abweist mit der erwarteten `desc:`-Begründung.
- M3-Slice-Plan in `done/` beschreibt das Vorgehen kurz.
- `make lint` läuft auf dem M3-Stand grün — alle Regeln matchen real,
  kein Block ist tot.
- Eintrag in `carveouts.md` von „temporär" auf aufgehoben verschoben.

## Out of Scope

- Strukturelle Änderungen an den Regelblöcken (Konvention sitzt seit
  M2b).
- Erweiterung um neue Schichten oder Regeln.

## Bezug

- Auslösende Spec: `LH-FA-ARCH-003`, `spec/architecture.md` §4.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `depguard`-Regeln matchen nichts.
- Hängt von: M3 `u-boot init` (erste produktive Pakete pro Schicht).
