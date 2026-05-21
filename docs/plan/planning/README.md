# docs/plan/planning

Planning-Artefakte (Slice- und Tranchen-Pläne) für u-boot.

## Lifecycle

Artefakte durchlaufen vier Verzeichnisse in dieser Reihenfolge
(`LH-FA-PROJDOCS-003`):

```
open/ → next/ → in-progress/ → done/
```

- Übergänge erfolgen per `git mv`, damit die Datei-Historie erhalten
  bleibt.
- Ein Artefakt darf nicht in mehreren Lifecycle-Verzeichnissen
  gleichzeitig liegen.
- Inhalte in `done/` dürfen nachträglich nur korrigierend (Tippfehler,
  Querverweise, Archiv-Hinweise) geändert werden. Substanzielle
  inhaltliche Änderungen erzeugen ein neues Artefakt in `open/` oder
  `next/` mit Verweis auf den vorhergehenden Stand.

## Dateiname-Konvention

Zwei verbindliche Formate, abhängig vom Artefakttyp:

- `slice-<phase>-<kebab-slug>.md` – Slice-Pläne pro Meilenstein-Phase.
  Beispiel: `slice-m1-repo-skeleton.md`.
- `tranche-<nr>-<kebab-slug>.md` – Tranchen-Pläne innerhalb eines
  Slice. Beispiel: `tranche-01-init-flow.md`.

Ein Artefakt verwendet genau eines der beiden Formate. Slice-Pläne sind
für phasenübergreifende Vorhaben; Tranchen-Pläne zerlegen einen Slice in
inkrementell auslieferbare Stücke.

## Ausnahme: Roadmap

`docs/plan/planning/in-progress/roadmap.md` ist ein übergreifendes
Master-Dokument, das laufend gepflegt wird (`LH-FA-PROJDOCS-003`). Es
folgt keinem der beiden Formate und fasst den Stand aller Slices und
Tranchen zusammen.
