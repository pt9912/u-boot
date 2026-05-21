# Slice V1: gomodguard-Regeln definieren

## Auslöser

`.golangci.yml` aktiviert `gomodguard_v2` mit `blocked: {}` (leerer
Block-Set). Die Konvention `LH-QA-004` verlangt, dass das Profil
insgesamt nicht aufweicht; ein leerer Block-Set lässt jedes externe
Modul zu (`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Sobald u-boot die ersten externen Modul-Dependencies bekommt
(voraussichtlich `spf13/cobra` mit dem CLI-Slice und `gopkg.in/yaml.v3`
mit dem Konfigurations-Slice), wird die `gomodguard_v2`-Konfiguration
mit konkreten Block-/Allow-Regeln versehen.

## Akzeptanzkriterien

- `.golangci.yml` `gomodguard_v2.blocked.modules` enthält mindestens
  eine Block-Regel für ein bekanntes Anti-Modul (z. B. veraltete
  yaml-Bindings, alte logging-Libraries) oder eine `allowed.modules`-
  Liste, die nur explizit freigegebene Module erlaubt.
- Begründung (`Why:`) pro Regel als Kommentar in `.golangci.yml`.
- `docs/user/quality.md` §1.2: gomodguard-Zeile bekommt einen
  Spalten-Hinweis auf die Regel-Quelle.
- `make lint` läuft grün auf dem ersten Commit mit externen Deps.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf den Aufhebungs-Commit als gelöst markieren.

## Bezug

- Auslösende Konfig: `.golangci.yml` `gomodguard_v2.blocked: {}`.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `gomodguard_v2` leer.
- Hängt von: erster Commit mit externen Modul-Dependencies (vermutlich
  M3 oder M4).
