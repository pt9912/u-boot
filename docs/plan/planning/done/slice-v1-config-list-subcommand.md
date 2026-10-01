# Slice V1: `u-boot config list` als eigener Subcommand (strukturiertes Path-Value-Listing)

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `HASH`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../done/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope (Carveout-Eintrag entfernt).

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

- [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration) (Config-Subcommand) — Spec listet `list`
  nicht; Erweiterung ist Konsument-Komfort-Argument.

## Closure-Notiz

**Geliefert:** `u-boot config list` (strukturiertes Pfad-Wert-Listing). `ConfigUseCase.List`: Whitelist-Pfade mit Wert, plus `services.<svc>.enabled` und `devcontainer.features.<name>.*` aus dem Dokument, nach Pfad sortiert; human `pfad=wert`, JSON `data.entries[]`; `--dry-run`/`--diff` abgelehnt (Exit 2).

**Sensoren:** `make gates` (lint, test, coverage-gate, docs-check) grün; Funktionsprobe mit
dem gebauten Binary (Mehrfach-`get`/`set`, `list`, atomarer Abbruch, Hint-Envelope).
`make test-docker` für diese Änderung nicht gesondert ausgeführt (kein Docker-Pfad berührt).

**Doku:** `docs/user/cli-json-output.md` §6.9, Benutzerhandbuch §Konfiguration, Beispiele, README, CHANGELOG.
