# Slice V1: `u-boot config set` Multi-Path-Set (mehrere Pfade in einem Call)

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `acc0233`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../done/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope (Carveout-Eintrag entfernt).

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

- [`LH-FA-CONF-001`](../../../../spec/lastenheft.md#lh-fa-conf-001--projektkonfiguration) (Config-Subcommand) — Spec listet Multi-
  Path nicht; Erweiterung ist Use-Case-Druck-Argument.

## Closure-Notiz

**Geliefert:** `config set <p1> <v1> <p2> <v2>…` atomar. `ConfigUseCase.SetMany`: alle Paare der Reihe nach auf einem In-Memory-Dokument (Coerce, Schema-Re-Validierung, Allowlist-Prüfung je Paar), danach genau ein `WriteFile`; Fehler → Datei byte-identisch; ungerade Argumentzahl → Exit 2; Listenpfade (`featureSources.allow`, `egress.allow`) über Plan-Funktionen (Marshal-Rewrite).

**Sensoren:** `make gates` (lint, test, coverage-gate, docs-check) grün; Funktionsprobe mit
dem gebauten Binary (Mehrfach-`get`/`set`, `list`, atomarer Abbruch, Hint-Envelope).
`make test-docker` für diese Änderung nicht gesondert ausgeführt (kein Docker-Pfad berührt).

**Doku:** `docs/user/cli-json-output.md` §6.9, Benutzerhandbuch §Konfiguration, Beispiele, README, CHANGELOG.
