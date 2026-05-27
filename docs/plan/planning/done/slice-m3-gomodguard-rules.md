# Slice M3: gomodguard-Regeln definieren

> **Status:** Done
> **DoD:** Commit `201fb4b`

## Auslöser

`.golangci.yml` aktivierte `gomodguard_v2` mit `blocked: {}` (leerer
Block-Set) — das Profil ließ damit jedes externe Modul ohne Prüfung
zu und verletzte `LH-QA-004` („Profil insgesamt nicht aufweicht",
`LH-FA-PROJDOCS-005`).

## Aufhebung

Vier Anti-Module sind jetzt in `.golangci.yml`
`gomodguard_v2.blocked` gelistet, jeweils mit `recommendations`
(stdlib- bzw. project-bevorzugte Alternative) und `reason`:

| Modul | Recommendation | Begründung |
| --- | --- | --- |
| `gopkg.in/yaml.v2` | `gopkg.in/yaml.v3` | v2 feature-frozen; M3-T1 hat v3-Codec-Adapter etabliert |
| `github.com/pkg/errors` | `errors`, `fmt` | stdlib seit Go 1.13 vollständig (errors.Is/As/Unwrap + wrap-aware fmt.Errorf) |
| `github.com/sirupsen/logrus` | `log/slog` | Project-Default für Logging (slice-m4-logging-port) |
| `go.uber.org/zap` | `log/slog` | Gleiches Baseline wie logrus |

## Unerwarteter Beifang: golangci-lint-Bump v2.12.1 → v2.12.2

Beim ersten Roll-out fiel auf, dass `gomodguard_v2` in
`golangci-lint v2.12.1` **silent fails** — der Schema-Verify ist grün,
die Regel ist enabled, das verbotene Modul ist Direct-Dep, der
Verbose-Lint listet den Linter als aktiv, aber `0 issues`. Verifiziert
mit minimal-Config und full-Config gegen einen Fixtur-Import.

`v2.12.2` (Patch-Release) fixt den Bug; identische Config feuert dort
korrekt. Daher wurde der Pin in `Dockerfile` + `Makefile` von
`v2.12.1` auf `v2.12.2` gehoben — Routine-Bump, dokumentiert sich im
Commit-Body. Kein eigener Slice nötig, weil unmittelbare Voraussetzung
für die Carveout-Auflösung.

## Geliefert

- `.golangci.yml` `gomodguard_v2.blocked`: 4 Block-Regeln (s. Tabelle).
  v2-Schema-Form `[{module, recommendations, reason}, ...]` (statt v1
  `blocked.modules: [name: {...}]`).
- `Dockerfile` + `Makefile`: `GOLANGCI_LINT_VERSION` v2.12.1 → v2.12.2.
- Fixtur-Test einmalig manuell durchgeführt (LH-QA-004
  Akzeptanz-Kriterium): `go get pkg/errors && go mod tidy` + Smoke-
  Datei `internal/hexagon/domain/gomodguard_smoke.go` mit
  `errors.New(...)` → `make lint` produziert genau diese Meldung:

  ```
  internal/hexagon/domain/gomodguard_smoke.go:3:8: import of package
  `github.com/pkg/errors` is blocked because the module is in the
  blocked modules list. `errors` and `fmt` are recommended modules.
  stdlib (errors.Is/As/Unwrap + wrap-aware fmt.Errorf) covers all
  pkg/errors features since Go 1.13. (gomodguard_v2)
  ```

  Reverted nach dem Test (`git status` clean).
- `carveouts.md`: `gomodguard_v2.blocked: {}` Zeile entfernt.
- Roadmap: Carveout-Auflösungs-Slice → Done.
- READMEs: Carveout-Count 14 → 13.

## Out of Scope

- **Allowlist-Modus** statt Blocklist: für ein kleines Projekt mit
  3 Direct-Deps wäre der Maintenance-Overhead (jede neue Dep verlangt
  Config-Edit) zu hoch. Blocklist wächst bedarfsgetrieben.
- **`recommendations:` für `pkg/errors`** — gomodguard rendert die
  Reason inkl. printf-Platzhaltern, daher `%w` aus dem Reason-Text
  entfernt; sinngemäße Umschreibung („wrap-aware fmt.Errorf").
- **Reproduzierbares `verify-gomodguard.sh`** analog zu
  `verify-depguard.sh`: kostet mehr Engineering (go.mod-Mutation +
  Cleanup) als der einmalige Fixtur-Test. Bei künftiger Block-Liste-
  Erweiterung optional einführbar.

## Bezug

- Auslösende Konfig: `.golangci.yml` `gomodguard_v2.blocked: {}` (vor diesem Slice).
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Hängt von: M3-T1 (yaml.v3 in go.mod), M3-T3 (Cobra).
