# docs/plan/planning

Planning-Artefakte (Slice- und Tranchen-Pläne) für u-boot.

## Lifecycle

Artefakte durchlaufen vier Verzeichnisse in dieser Reihenfolge
([`LH-FA-PROJDOCS-003`](../../../spec/lastenheft.md#lh-fa-projdocs-003-planning-lifecycle)):

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
- Vor dem Übergang nach `done/` muss der Slice eine
  Verification-Evidence nach
  [`../../../harness/verification.md`](../../../harness/verification.md)
  tragen oder auf ein eigenes Evidence-Artefakt verweisen.

## Referenzsemantik

Planning-Artefakte folgen [`LH-FA-PROJDOCS-006`](../../../spec/lastenheft.md#lh-fa-projdocs-006-dokumentationsreferenzmodell) und
[`ADR-0013`](../adr/0013-dokumentationsreferenzmodell.md):
Slices referenzieren `LH-*` und aktive ADRs normativ. Referenzen auf
andere Slices, Carveouts oder Roadmap/Wellen sind nur Trigger-,
Owner-, Closure- oder Orchestrierungskontext und erzeugen keine
Spezifikation.

## Dateiname-Konvention

Zwei verbindliche Formate, abhängig vom Artefakttyp:

- `slice-<phase>-<kebab-slug>.md` – Slice-Pläne pro Meilenstein-Phase.
  Beispiel: [`slice-m1-repo-skeleton.md`](done/slice-m1-repo-skeleton.md).
- `tranche-<nr>-<kebab-slug>.md` – Tranchen-Pläne innerhalb eines
  Slice. Beispiel: `tranche-01-init-flow.md`.

Ein Artefakt verwendet genau eines der beiden Formate. Slice-Pläne sind
für phasenübergreifende Vorhaben; Tranchen-Pläne zerlegen einen Slice in
inkrementell auslieferbare Stücke.

## Ausnahmen: Master-Dokumente

Zwei übergreifende Master-Dokumente folgen keinem der beiden Formate
und liegen dauerhaft in `in-progress/`:

- [`in-progress/roadmap.md`](in-progress/roadmap.md) — Stand aller
  Slices und Tranchen ([`LH-FA-PROJDOCS-003`](../../../spec/lastenheft.md#lh-fa-projdocs-003-planning-lifecycle)).
- [`in-progress/carveouts.md`](in-progress/carveouts.md) — Inventar
  aller bewussten Carveouts mit Plan-Verweis ([`LH-FA-PROJDOCS-005`](../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)).
