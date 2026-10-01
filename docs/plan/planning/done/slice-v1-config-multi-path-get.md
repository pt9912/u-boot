# Slice V1: `u-boot config get` Multi-Pfad-Get / `--json-array`

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `HASH`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../done/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope (Carveout-Eintrag entfernt).

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

- [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (Path-Whitelist) — Spec listet Multi-Path
  nicht; Erweiterung ist Use-Case-Druck-Argument.

## Closure-Notiz

**Geliefert:** `config get <p1> <p2>…` und `--json-array`. `ConfigUseCase.GetMany` (ein Lesezugriff, Reihenfolge, all-or-nothing); human: ein Wert je Zeile; `--json`: `data.entries[]` (ohne `omitempty`); ein Pfad ohne Flag behält die alte Form.

**Sensoren:** `make gates` (lint, test, coverage-gate, docs-check) grün; Funktionsprobe mit
dem gebauten Binary (Mehrfach-`get`/`set`, `list`, atomarer Abbruch, Hint-Envelope).
`make test-docker` für diese Änderung nicht gesondert ausgeführt (kein Docker-Pfad berührt).

**Doku:** `docs/user/cli-json-output.md` §6.9, Benutzerhandbuch §Konfiguration, Beispiele, README, CHANGELOG.
