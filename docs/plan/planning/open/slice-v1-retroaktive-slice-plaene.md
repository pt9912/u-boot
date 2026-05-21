# Slice V1: Retroaktive Slice-Pläne für M1/M2/M2b/M2c/M2d

## Auslöser

Die abgeschlossenen Bootstrap-Slices M1, M2, M2b, M2c und M2d sind
nur in der Roadmap und in Commit-Messages dokumentiert. Es gibt keine
zugehörigen Slice-Pläne in `docs/plan/planning/done/`, wie es
`LH-FA-PROJDOCS-003` für den Standard-Lifecycle vorsieht
(`LH-FA-PROJDOCS-005`).

M2-Review #10 hatte das schon angeschnitten („Slice-Plan folgt mit
nächster Iteration"); seitdem ist der Stand unverändert.

## Aufhebungsbedingung

Pro abgeschlossenem Bootstrap-Slice (M1, M2, M2b, M2c, M2d) wird ein
schlanker `slice-m<n>-<slug>.md` retrospektiv in
`docs/plan/planning/done/` angelegt, der den Auslöser, das gelieferte
Artefakt, die Akzeptanz und den Commit-Hash benennt. Die Roadmap
verlinkt jeweils auf den Slice-Plan statt nur auf den Commit-Hash.

## Akzeptanzkriterien

- `docs/plan/planning/done/slice-m1-repo-skeleton.md`
- `docs/plan/planning/done/slice-m2-hexagonale-architektur.md`
- `docs/plan/planning/done/slice-m2b-solid-lint-profil.md`
- `docs/plan/planning/done/slice-m2c-ci-pipeline.md`
- `docs/plan/planning/done/slice-m2d-carveout-disziplin.md`
- Jeweils Mindestabschnitte: Auslöser, Lieferumfang, Akzeptanz,
  Commit-Hash.
- Roadmap-Spalte „Artefakt" verweist auf die Slice-Pläne.
- Eintrag in `carveouts.md` von „temporär" auf aufgehoben verschoben.

## Out of Scope

- Rückwirkende Slice-Pläne für jeden einzelnen Spec-Review-Commit
  (`dab4e45`, `2ea534b`, …) — die sind Politur und nicht
  Slice-relevant.
- Tranchen-Pläne — Bootstrap-Slices sind nicht in Tranchen
  unterteilt.

## Bezug

- Auslösende Spec: `LH-FA-PROJDOCS-003` Lifecycle.
- Inventar-Eintrag: [`carveouts.md`](../in-progress/carveouts.md) →
  Slice-Pläne fehlen.
- Hängt von: nichts; kann jederzeit abgearbeitet werden, am besten
  vor dem ersten externen Contributor.
