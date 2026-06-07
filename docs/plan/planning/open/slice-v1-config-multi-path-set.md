# Slice V1: `u-boot config set` Multi-Path-Set (mehrere Pfade in einem Call)

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Multi-Path-Set-Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../next/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure des config-Slice trägt den
> Eintrag nach).

## Auslöser

`u-boot config set` heutige Surface (`cli/config.go:102`)
trägt `cobra.ExactArgs(2)` — ein Pfad-Wert-Paar pro
Set-Aufruf. CI-Use-Cases mit mehreren atomar zu setzenden
Werten (z. B. `set project.name X devcontainer.enabled true`)
brauchen heute zwei separate Aufrufe ohne Transaktions-
Semantik (Erfolg-Halb-Schreibe-Halb-Bruch möglich).

V1-Trade-off: schmale Surface > Multi-Path. Multi-Path mit
Transaktion (alle oder keine schreiben) wandert in diesen
Folge-Slice.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Druck nach atomarem Multi-Pfad-Set**: CI-Use-
  Case beschwert sich über partial-state bei sequentiellen
  Set-Aufrufen.
- **`config set` Schema-Erweiterung** die mehrere zusammen-
  hängende Felder atomar erfordert (z. B. Coupled-Path-
  Constraints).

## Lösungs-Skizze (vorläufig)

`cli/config.go` Args-Form auf `cobra.MinimumNArgs(2)` mit
Pair-Parser; `ConfigSetRequest.Paths []ConfigPathValue`
statt single `Path` + `Value`; Application-Layer transaktional:
alle Coerce + Schema-Validate VOR erstem `WriteFile`-Aufruf;
WriteFile als einzelner finaler Schreib-Akt.

## Spec-Bezug

- `LH-FA-CONF-001` (Config-Subcommand) — Spec listet Multi-
  Path nicht; Erweiterung ist Use-Case-Druck-Argument.
