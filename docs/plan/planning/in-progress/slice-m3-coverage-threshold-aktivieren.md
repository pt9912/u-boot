# Slice M3: Coverage-Schwellwert anheben

## Auslöser

`LH-FA-BUILD-008` definiert einen Bootstrap-Modus für den Coverage-Gate:
solange `./internal/...` keinen produktiven Code enthält, läuft
`scripts/coverage-gate.sh` mit Schwellwert `0` und akzeptiert leeren
Coverage-Input ohne Fail. Das vermeidet während M1–M2c falsche
Grün-Signale, ist aber ein temporärer Carveout
(`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Sobald der erste produktive Slice (M3 `u-boot init`) mindestens ein
Paket unter `./internal/...` (vermutlich `internal/hexagon/domain/` oder
`internal/hexagon/application/`) liefert, wechselt der Bootstrap-Pfad
in `Dockerfile`-Stage `coverage` von `COVERPKG=""` auf `COVERPKG=$(go
list ./internal/...)` automatisch. Der Schwellwert muss dann manuell
gehoben werden.

## Akzeptanzkriterien

- `make coverage-gate THRESHOLD=80` läuft grün gegen den M3-Stand.
- `Dockerfile`-`ARG COVERAGE_THRESHOLD` Default bleibt `0` (Bootstrap-
  Schutz für seltene Edge-Cases), aber das Makefile-Default
  `THRESHOLD ?= 0` wird auf `80` gehoben.
- `.github/workflows/ci.yml`: `gates`-Job setzt `make coverage-gate
  THRESHOLD=80` oder ruft `make gates THRESHOLD=80`.
- `docs/user/quality.md` §3 Coverage: Empfehlung „80 % als erster Wert"
  wird zur dokumentierten Pflicht-Schwelle; langfristiger Zielwert
  90 % (analog m-trace/k-deskflight) bleibt als Folgepunkt.
- Eintrag in `carveouts.md` von „temporär" auf „aufgehoben mit Commit
  <hash>" verschoben (oder Zeile entfernt) und der Slice in `done/`.

## Bezug

- Auslösende Spec: `LH-FA-BUILD-008` Coverage-Bootstrap.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `COVERAGE_THRESHOLD=0` Bootstrap.
- Hängt von: M3 `u-boot init` (erste produktive Pakete).
- Verwandt: `slice-m3-depguard-aktivierung-verifizieren.md` (gleicher
  Trigger).
