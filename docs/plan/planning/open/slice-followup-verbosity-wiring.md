# Slice Follow-up: `--verbose` / `--debug` zum Logger-Level verdrahten

## Auslöser

`LH-FA-CLI-005` schreibt vier Verbosity-Stufen vor: `--quiet`,
Standard, `--verbose`, `--debug`. Mit M4-T7 wurden die zugehörigen
`PersistentFlags` auf der Root-Cobra eingebaut; `--quiet` ist
load-bearing (filtert OK-Items aus dem Doctor-Render), `--verbose`
und `--debug` werden aber heute **akzeptiert ohne sichtbaren
Effekt** — der `slog`-basierte `driven.Logger`-Adapter ist im
Wiring (`cmd/uboot/main.go`) fix auf `slog.LevelInfo` konstruiert
(`feedback-carveouts-need-plans`).

## Aufhebungsbedingung

Beim Doctor-CLI-Aufruf (und perspektivisch bei jedem
Application-Service mit Logger-Konsumenten) soll:

- `--debug` den Logger-Level auf `slog.LevelDebug` heben (alle
  service-internen Debug-Calls werden sichtbar);
- `--verbose` analog auf `slog.LevelDebug` heben (oder einen neuen
  „Verbose"-Pegel zwischen Info und Debug, falls künftig ein
  service-spezifischer Verbose-Logger-Pfad entsteht);
- `--quiet` den Logger-Level auf `slog.LevelWarn` heben (Logger-
  Output zusätzlich zur Doctor-Filter-Logik im Rendering reduzieren).

## Mechanik (Vorschlag)

1. `internal/adapter/driven/logger/logger.go`: `New` signature von
   `(out, format, level slog.Level)` auf `(out, format, level
   slog.Leveler)` heben. `*slog.LevelVar` implementiert `Leveler`
   und ist nach Konstruktion mutierbar.
2. `cmd/uboot/main.go`: `levelVar := new(slog.LevelVar);
   levelVar.Set(slog.LevelInfo)`; an `logger.New(...)` übergeben.
3. `cli.App` bekommt einen `*slog.LevelVar`-Field; Konstruktor um
   einen Parameter erweitert. `buildRootCommand` schreibt in
   `levelVar` per `PersistentPreRunE` (nach Cobra-Parse aber vor
   RunE):
   ```go
   root.PersistentPreRunE = func(*cobra.Command, []string) error {
       switch {
       case a.debug:   a.logLevel.Set(slog.LevelDebug)
       case a.verbose: a.logLevel.Set(slog.LevelDebug)
       case a.quiet:   a.logLevel.Set(slog.LevelWarn)
       }
       return nil
   }
   ```

## Akzeptanzkriterien

- `u-boot doctor --debug` zeigt mindestens einen Debug-Log-Eintrag
  vom DoctorService auf stderr (z. B. `soft-existing-detection`
  oder das in T2 eingebaute `"doctor: starting checks"`).
- `u-boot doctor --verbose` analog (gleicher Pegel).
- `u-boot doctor --quiet` unterdrückt die Info-Log-Zeile
  `"doctor: checks complete"` (Logger-Pegel ≥ Warn) zusätzlich zur
  Filter-Logik im Doctor-Render.
- `cli/doctor.go` `doctorFlags`-Struct kann auf `{Strict, Quiet}`
  reduziert werden (die `Verbose`/`Debug`-Felder werden dead-state-
  free, weil der Effekt jetzt am Logger-Level passiert).
- `make gates` grün.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf
  den Aufhebungs-Commit als gelöst markieren.

## Out of Scope

- Eigene Logger-Pegel jenseits von Debug/Info/Warn/Error.
- Per-Subcommand-Verbosity-Override (heute global persistent).

## Bezug

- Auslösende Spec: `LH-FA-CLI-005`.
- Slice-Vorgänger: [`slice-m4-doctor`](../done/slice-m4-doctor.md)
  T7 hat `--verbose`/`--debug` per Spec eingebaut, aber bewusst
  ohne Logger-Level-Wiring (Scope-Constraint).
- Inventar-Eintrag in [`carveouts.md`](../in-progress/carveouts.md).
