# Slice V1: `AddPreviewMode`-Alias entfernen

> **Status:** geplant für v0.4.0+ — Cleanup-Folge-Slice aus
> [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
> T0-(c) Carveout (Carveout-Plan-Pflicht, MEMORY.md
> [[feedback_carveouts_need_plans]]).

## Auslöser

Der Init-Slice T1-F (`9ed3e34`) hat den kanonischen Enum
`driving.PreviewMode` etabliert (vorher `driving.AddPreviewMode`)
und das Alt-Symbol als **type-Alias** erhalten:

```go
// internal/hexagon/port/driving/addservice.go
type AddPreviewMode = PreviewMode
```

Der Gleichheits-Zeichen-Alias (statt `type AddPreviewMode
PreviewMode`) ist load-bearing — er hält die Factory-Signatur
`func(driving.AddPreviewMode) (driven.FileSystem,
driven.RecorderPort)` assignment-kompatibel zu
`func(driving.PreviewMode) ...`, damit die add-T1-Tests
(`addservice_factory_test.go`) ohne Signatur-Cast grün bleiben.

Der Alias ist seit Init-Slice der **einzige** Service-Prefix-
Mode-Alias. Folge-Slices (`generate`/`remove`/`config set`)
referenzieren `driving.PreviewMode` direkt; weitere
`XxxPreviewMode`-Aliases entstehen nicht (Init T0-(c) Alias-
Lebensdauer-Pflicht).

## Scope

**Pflicht:**

- [`internal/hexagon/port/driving/addservice.go`](../../../../internal/hexagon/port/driving/addservice.go):
  `type AddPreviewMode = PreviewMode`-Deklaration (Z. ~159) plus
  den umgebenden Carveout-Doku-Block (Z. ~145-158) entfernen.
- [`internal/hexagon/port/driving/previewmode_test.go`](../../../../internal/hexagon/port/driving/previewmode_test.go):
  komplett entfernen. Die Datei pinnt **nur** die Alias-Identität
  und die Funktions-Typ-Kompatibilität; nach Alias-Removal hat
  sie keinen Sinn mehr.
- [`internal/hexagon/application/addservice_factory_test.go`](../../../../internal/hexagon/application/addservice_factory_test.go)
  Z. ~25, ~125: `driving.AddPreviewMode` → `driving.PreviewMode`.
- [`internal/adapter/driving/cli/previewmode_internal_test.go`](../../../../internal/adapter/driving/cli/previewmode_internal_test.go)
  Z. ~19: `driving.AddPreviewMode` → `driving.PreviewMode`.
- Weitere `grep -rn "AddPreviewMode"`-Treffer (sollten keine
  Production-Sites sein nach Init T1-F — Sweep zur Sicherheit).

**Out-of-Scope:**

- Andere init-Slice-Carveouts (`mapErrorToDiagnostic`-Registry,
  `previewFSFactory`-Generalisierung) — die haben eigene Slices.
- Renaming des kanonischen `PreviewMode`-Types selbst — der ist
  jetzt der Vertrag aller modifying-Subcommands.

## Done-Definition

- `grep -rn "AddPreviewMode" internal/ cmd/` liefert null
  Treffer (`scripts/`/`docs/` ausgenommen — historische Bezüge
  in done/-Plänen bleiben erhalten).
- `make gates` grün; Coverage-Gate weiterhin ≥ 90%.
- CHANGELOG-Eintrag unter `Changed` (Breaking-Change für
  externe API-Konsumenten der port/driving-Types — falls
  vorhanden; intern ist u-boot der einzige Konsument, deshalb
  kein Migrations-Block nötig).

## Reihenfolge im Cluster

Außerhalb des `slice-v1-cli-json-dry-run`-Folge-Slice-Blocks —
reine Code-Cleanup, kein Schema-Migrations-Anteil. Kann
unabhängig vom Cluster-T_close-Lauf gemerged werden,
**vorzugsweise nach** mindestens einem weiteren Folge-Slice
(z. B. `generate`), damit verifiziert ist, dass das
Pattern-Erbe ohne den Alias funktioniert.

## Risiko

Minimal. Der Alias hat null Production-Konsumenten in u-boot
selbst; alle Reste sind Test-Code. Externe Konsumenten der
`port/driving`-Types existieren nicht (u-boot ist Self-Contained
CLI ohne separat distributiertes API-Package). LOC-Schätzung:
~30 LOC entfernt (Alias + Carveout-Doku-Block +
`previewmode_test.go`), ~10 LOC angepasst (Test-Identifier-
Rename).
