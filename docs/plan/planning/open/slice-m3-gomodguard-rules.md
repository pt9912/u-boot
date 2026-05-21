# Slice M3: gomodguard-Regeln definieren

## Auslöser

`.golangci.yml` aktiviert `gomodguard_v2` mit `blocked: {}` (leerer
Block-Set). Die Konvention `LH-QA-004` verlangt, dass das Profil
insgesamt nicht aufweicht; ein leerer Block-Set lässt jedes externe
Modul zu (`LH-FA-PROJDOCS-005`).

## Aufhebungsbedingung

Mit M3-T1 ist `gopkg.in/yaml.v3` in `go.mod` gelandet; M3-T3 bringt
zusätzlich Cobra. Spätestens mit M3-T5 (Carveout-Cleanup) müssen
konkrete `gomodguard_v2`-Regeln stehen, die diese erlaubten Module
explizit zulassen und bekannte Anti-Module blockieren.

## Akzeptanzkriterien

- `.golangci.yml` `gomodguard_v2.blocked.modules` enthält mindestens
  eine Block-Regel für ein bekanntes Anti-Modul (z. B. veraltete
  yaml-Bindings, alte logging-Libraries) oder eine `allowed.modules`-
  Liste, die nur explizit freigegebene Module erlaubt.
- Begründung (`Why:`) pro Regel als Kommentar in `.golangci.yml`.
- `docs/user/quality.md` §1.2: gomodguard-Zeile bekommt einen
  Spalten-Hinweis auf die Regel-Quelle.
- `make lint` läuft grün gegen die aktuellen u-boot-Deps und gegen
  einen Fixtur-Test, der einen verbotenen Import einführt und das
  rote Lint-Ergebnis abnimmt.
- Zeile in `carveouts.md` entweder entfernen oder mit Verweis auf
  den Aufhebungs-Commit als gelöst markieren.
- Roadmap-Eintrag in der Carveout-Tabelle als Done markieren.

## Bezug

- Auslösende Konfig: `.golangci.yml` `gomodguard_v2.blocked: {}`.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  `gomodguard_v2` leer.
- Hängt von: M3-T1 (yaml.v3 schon drin) + M3-T3 (Cobra) — beide
  liefern reale Modul-Dependencies, die in den Regeln explizit
  vorkommen müssen.
- Phase: M3-T5 (Carveout-Cleanup).
