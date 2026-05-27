# Slice M3: depguard-Regelaktivierung verifizieren

> **Status:** Done
> **DoD:** M3-T5-Commit (siehe `git log --grep=m3-t5`)

## Auslöser

`LH-FA-ARCH-003` und `spec/architecture.md` §4 definieren acht
`depguard`-Regelblöcke für die hexagonalen Schicht-Imports. Sie waren
seit M2b in `.golangci.yml` aktiv, **matchten aber nichts**, weil
`./internal/...` keinen produktiven Code enthielt. Die Regeln greifen
laut Konvention „automatisch mit dem ersten Paket pro Schicht" — das
war bis M3 nur Behauptung, kein Test (`LH-FA-PROJDOCS-005`).

## Aufhebung

Mit dem M3-Init-Flow (Slices T1..T4c) liegt jetzt mindestens ein
Paket pro Schicht in `./internal/...`:

- `internal/hexagon/domain/` (domain)
- `internal/hexagon/application/` (application + managedblock)
- `internal/hexagon/port/driving/` + `internal/hexagon/port/driven/`
- `internal/adapter/driving/cli/` + `internal/adapter/driven/{fs,yaml,git,clock,progress}/`

`scripts/verify-depguard.sh` automatisiert die Negativ-Verifikation:
für jede der 8 Regeln wird ein deklariert verbotener Import in einer
temporären `verify_depguard_violation.go` im passenden Layer
eingebaut, `make lint` läuft (depguard fail erwartet, `desc:` muss
matchen), die Stub-Datei wird wieder entfernt. `trap` sichert Cleanup
auch bei Abbruch.

## Geliefert

- `scripts/verify-depguard.sh` mit Cycle-Safety-Mapping aus der
  realen Import-Graphen-Analyse (jeder Forbidden-Import darf keinen
  Go-Import-Cycle erzeugen, sonst maskiert der Typecheck-Fehler die
  depguard-Regel).
- `make verify-depguard`-Target (manuell / on-demand, nicht Teil von
  `make gates`, weil pro Regel ein `make lint` Docker-Build läuft —
  Volldurchlauf ~3–5 min).
- Verifikation aller 8 Regeln grün:

  | Regel | Forbidden Import | desc-Substring |
  | --- | --- | --- |
  | `domain-isoliert` | `internal/hexagon/port/driven` | `domain must not depend on port` |
  | `application-no-adapter` | `internal/adapter/driven/clock` | `application must depend on ports, not on adapter implementations` |
  | `port-no-application` | `internal/adapter/driven/clock` (von `port/driving/`) | `port must not depend on adapter` |
  | `port-driving-no-driven` | `internal/hexagon/port/driven` | `driving port must not depend on driven port` |
  | `port-driven-no-driving` | `internal/hexagon/port/driving` | `driven port must not depend on driving port` |
  | `adapter-no-application` | `internal/hexagon/application` (von `adapter/driven/fs/`) | `adapter must implement ports, not consume application` |
  | `adapter-driving-no-driven` | `internal/adapter/driven/fs` (von `adapter/driving/cli/`) | `driving adapter must not depend on driven adapter` |
  | `adapter-driven-no-driving` | `internal/adapter/driving/cli` (von `adapter/driven/clock/`) | `driven adapter must not depend on driving adapter` |

- `carveouts.md`-Eintrag „depguard-Regeln matchen nichts" entfernt.
- `make lint` auf dem aktuellen Stand ist 0 issues.

## Bezug

- Auslösende Spec: `LH-FA-ARCH-003`, `spec/architecture.md` §4.
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Hängt von: M3-T1..T4c — erste produktive Pakete pro Schicht.
