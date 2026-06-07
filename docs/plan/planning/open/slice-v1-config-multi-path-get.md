# Slice V1: `u-boot config get` Multi-Pfad-Get / `--json-array`

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum config-multi-path-get-Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../next/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure des config-Slice trägt den
> Eintrag nach).

## Auslöser

`u-boot config get` heutige Surface (`cli/config.go:81`)
trägt `cobra.ExactArgs(1)` — ein Pfad pro Get-Aufruf. CI-
Use-Cases die mehrere Werte gleichzeitig brauchen
(`get project.name devcontainer.enabled`) führen heute
mehrere Aufrufe + Output-Concat. Mit `--json` wäre eine
Array-Form (`data.entries: [{path, value}, ...]`) natürlicher.

V1-Trade-off: Single-Path-Form folgt Cluster-Slice T0-(c).
Multi-Path-Get plus `--json-array`-Variante wandert in diesen
Folge-Slice.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Druck nach Batch-Get**: CI-Use-Case beschwert
  sich über Multi-Call-Latenz oder über die Output-Concat-
  Pflicht.
- **Multi-Path-Set-Slice** ([`slice-v1-config-multi-path-set`](slice-v1-config-multi-path-set.md))
  geht live: symmetrische Surface-Erweiterung erwartet.

## Lösungs-Skizze (vorläufig)

`cli/config.go` Args-Form auf `cobra.MinimumNArgs(1)`;
`ConfigGetRequest.Paths []ConfigPath` statt single `Path`;
Application-Layer Loop pro Pfad mit Sammel-Response
`ConfigGetResponse.Entries []ConfigPathValue`. JSON-Envelope
trägt `data.entries []` ohne `omitempty` (Empty-Pin).

## Spec-Bezug

- `LH-FA-CONF-005` (Path-Whitelist) — Spec listet Multi-Path
  nicht; Erweiterung ist Use-Case-Druck-Argument.
