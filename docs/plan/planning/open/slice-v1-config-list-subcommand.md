# Slice V1: `u-boot config list` als eigener Subcommand (strukturiertes Path-Value-Listing)

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum config-list-Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../next/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure des config-Slice trägt den
> Eintrag nach).

## Auslöser

`u-boot config` (bare) liefert byte-identisch das gesamte
`u-boot.yaml`-File (`ConfigShowResponse.Body []byte`). Ein
strukturierter Pfad-Wert-Tree (`[{path: "project.name",
value: "demo"}, ...]`) wäre konsument-freundlicher für JSON-
Pipelines die alle gesetzten Pfade ohne YAML-Parsing
enumerieren wollen.

V1-Trade-off: schmale Surface > strukturiertes Listing.
`u-boot config list` als eigener Subcommand wandert in
diesen Folge-Slice.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Druck nach Pfad-Enumeration**: CI-Use-Case
  beschwert sich über YAML-Parse-Pflicht in der Konsument-
  Pipeline.
- **`config get` Multi-Pattern-Support**: wenn Glob-Patterns
  in `get` landen (`config get "services.*.enabled"`), ist
  `list` der natürliche Vorgänger.

## Lösungs-Skizze (vorläufig)

Neuer `cli/config.go` `newConfigListCommand(a *App)` analog
`newConfigGetCommand`; `ConfigListResponse.Entries
[]ConfigPathValue`; Application-Layer enumeriert per
`domain.AllConfigPaths()` mit Lenient-Extract pro Pfad.

## Spec-Bezug

- `LH-FA-CONF-001` (Config-Subcommand) — Spec listet `list`
  nicht; Erweiterung ist Konsument-Komfort-Argument.
