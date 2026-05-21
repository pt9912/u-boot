# Slice M3: Coverage-Schwellwert aktivieren

## Auslöser

`LH-FA-BUILD-008` definiert einen Bootstrap-Modus für den Coverage-Gate:
solange `./internal/...` keinen produktiven Code enthält, läuft
`scripts/coverage-gate.sh` mit Schwellwert `0` und akzeptiert leeren
Coverage-Input ohne Fail. Das vermeidet während M1–M2c falsche
Grün-Signale, ist aber ein temporärer Carveout
(`LH-FA-PROJDOCS-005`).

## Aufhebung

Mit M3-T1 ([`slice-m3-init-flow.md`](slice-m3-init-flow.md)) sind die
ersten produktiven Pakete unter `./internal/...` entstanden;
`scripts/coverage-gate.sh` läuft seitdem im Production-Pfad. Die reale
Coverage liegt bei 93.90 %.

Die Schwelle wurde direkt auf **90 %** gehoben (statt des ursprünglich
geplanten Zwischenschritts 80) — User-Entscheidung: der aktuelle
Puffer (93.90 % vs. 90 %) ist ausreichend, und etwaige Coverage-Drops
in T2/T3 sollen mit zusätzlichen Tests aufgefangen werden, nicht durch
Schwellen-Absenken.

## Geliefert

- `Makefile`: `THRESHOLD ?= 0` → `THRESHOLD ?= 90`.
- `Dockerfile`: `ARG COVERAGE_THRESHOLD=0` → `ARG COVERAGE_THRESHOLD=90`.
- `docs/user/quality.md` §3 aktualisiert: Schwellwert ist verbindlich
  90 %, Bootstrap-Pfad im Skript bleibt als Fallback erhalten.
- `carveouts.md`-Eintrag `COVERAGE_THRESHOLD=0` ist aufgehoben.
- `make gates` grün bei 93.90 % gegen Schwelle 90 %.

## Bezug

- Auslösende Spec: `LH-FA-BUILD-008` Coverage-Bootstrap.
- Aufhebung dokumentiert in: [`carveouts.md`](carveouts.md).
- Hängt von: M3-T1 (erste produktive Pakete) — erfüllt mit Commit
  `132d1a1`.
- Aktivierung erfolgt in: diesem Commit.
