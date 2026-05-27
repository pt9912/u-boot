# Slice M4: Logging-Port einführen

> **Status:** Done
> **DoD:** Commit `<TBD>` (wird im Folge-Fixup konkretisiert)

## Auslöser

Die `forbidigo`-Regel in `.golangci.yml` verbot `fmt.Print*` mit der
Begründung *„use log/slog (LH-QA-004); a project-specific logging
port may replace this once it exists"*. Der referenzierte
„configured logging port" existierte nicht — es war prospektive Doku
(`LH-FA-PROJDOCS-005`). Mit M4-doctor steht ein Subbefehl an, der
nennenswert strukturiertes Logging produziert; der Port muss vor
diesem Konsumenten stehen.

## Aufhebung

u-boot hat einen eigenen Logging-Driven-Port (`driven.Logger`) und
einen slog-basierten Adapter. Application referenziert nur den Port;
`log/slog` wird ausschließlich im Adapter importiert.

**Konventionsabweichung vom ursprünglichen Slice-Text:** der
Slice-Plan suggerierte `internal/hexagon/port/driven/logger/` als
Subpackage. Übernommen wurde stattdessen die flache Layout-
Konvention der anderen Ports (`Clock`, `Confirmer`, `FileSystem` —
alle einer Datei in `port/driven/`): `port/driven/logger.go`,
package `driven`, interface `Logger`. Konsumenten schreiben
`driven.Logger`, parallel zu `driven.Confirmer` und `driven.Clock`.

## Geliefert

- **Driven-Port** `port/driven/logger.go` — `Logger` mit
  `Debug`/`Info`/`Warn`/`Error`-Methoden (`...any` für slog-konforme
  Key-Value-Pairs). Alle Methoden best-effort, kein Error-Return.
- **Adapter** `adapter/driven/logger/logger.go` — `New(out, format,
  level)` retourniert `driven.Logger`. `Format` als Enum
  (`FormatText` / `FormatJSON`); slog-Level pass-through.
- **5 Adapter-Tests**: Text-Output mit `key=value`, JSON-Output mit
  validem Schema, Level-Filter, alle 4 Levels, Fallback bei
  unbekanntem Format.
- **Service-Integration**: `InitProjectService`-Konstruktor um
  `logger driven.Logger` (6. Driven-Port) erweitert; `nil` → internal
  `noopLogger` (alle Levels Discard). Soft-Detection emittiert einen
  Debug-Log oberhalb der Schwelle (`baseDir`, `indicators`,
  `threshold`).
- **2 Service-Tests** für die Logger-Integration: Debug feuert über
  Threshold, Logger bleibt unterhalb der Schwelle still.
  `fakeLogger` in `fakes_test.go` zeichnet alle Calls auf.
- **`cmd/uboot/main.go`**: `logger.New(stderr, FormatText,
  slog.LevelInfo)` konstruiert + injiziert. Verbose/JSON-Flag-Wiring
  kommt mit M4-doctor.
- **`.golangci.yml`**: `forbidigo.msg` aktualisiert auf *„use the
  logging port at internal/hexagon/port/driven (driven.Logger)
  (LH-QA-004)"*. `gomodguard_v2`-Regeln für `logrus`/`zap` zeigen
  jetzt auf den existierenden Adapter statt auf den prospektiven
  Slice.
- **`carveouts.md`**: `forbidigo.msg`-Zeile entfernt.
- **Roadmap**: Slice → Done, Phase „M4-vorgezogen".
- **READMEs**: Carveout-Count 10 → 9.

## Out of Scope

- **Logging-Backend-Wahl jenseits von slog** (zerolog, zap, ...):
  separater ADR, falls je nötig. `gomodguard_v2` blockiert die
  Alternativen heute aktiv.
- **Strukturierte Telemetrie / OTel**: gehört zum OTel-Add-on-Slice
  (`LH-FA-ADD-004`).
- **CLI-Flags `--verbose` / `--debug` / `--json`**: kommen mit dem
  M4-doctor-Slice; der Adapter ist heute auf `Info` + `Text`
  hartcodiert.

## Bezug

- Auslösende Konfig: `.golangci.yml` `forbidigo.msg`.
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Hängt von: nichts; vorgezogen vor M4-doctor, damit der Konsument
  bei seinem Start auf einen stabilen Port aufsetzt.
