# Slice M3: Retroaktive Slice-Pläne für M1/M2/M2b/M2c/M2d

> **Status:** Done
> **DoD:** Commit `afee7e9`

## Auslöser

Die abgeschlossenen Bootstrap-Slices M1, M2, M2b, M2c und M2d waren
nur in der Roadmap und in Commit-Messages dokumentiert. Es gab keine
zugehörigen Slice-Pläne in `docs/plan/planning/done/`, wie es
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003-planning-lifecycle) für den Standard-Lifecycle vorsieht
([`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)).

M2-Review #10 hatte das angeschnitten („Slice-Plan folgt mit nächster
Iteration"); danach blieb der Stand unverändert. Bei M3-Anker-Triage
(2026-05-27) wurde der Slice von „M3-T5 (Carveout-Cleanup)" auf
„Later / eigene Sitzung" re-phasiert — diese Sitzung holt die
Retro-Pläne nach, weil das Memory-Drift-Risiko sonst weiter wächst.

## Aufhebung

Pro abgeschlossenem Bootstrap-Slice (M1, M2, M2b, M2c, M2d) wurde ein
schlanker `slice-m<n>-<slug>.md` retrospektiv in
`docs/plan/planning/done/` angelegt, der Auslöser, Lieferumfang,
Akzeptanz und Commit-Hash benennt. Die Roadmap-Spalte „Artefakt"
verweist jeweils auf den Slice-Plan statt nur auf den Commit-Hash.

## Geliefert

- 5 retro-Slice-Pläne in `docs/plan/planning/done/`:
  - [`slice-m1-repo-skeleton`](slice-m1-repo-skeleton.md) (Commit `7da05c7`)
  - [`slice-m2-hexagonale-architektur`](slice-m2-hexagonale-architektur.md) (Commit `9d191a5`)
  - [`slice-m2b-solid-lint-profil`](slice-m2b-solid-lint-profil.md) (Commit `365e532`)
  - [`slice-m2c-ci-pipeline`](slice-m2c-ci-pipeline.md) (Commit `9a74e35`)
  - [`slice-m2d-carveout-disziplin`](slice-m2d-carveout-disziplin.md) (Commits `5b92a97` + `3dc0467` + `0124742` + `b51a4aa`)
- Roadmap-Spalte „Artefakt" für die 5 M1–M2d-Zeilen auf die neuen
  Slice-Pläne umgestellt (zuvor: Commit-Hash direkt).
- `carveouts.md`: Retro-Slice-Zeile entfernt.
- READMEs: Carveout-Count 12 → 11.

## Out of Scope

- **Rückwirkende Slice-Pläne für jeden einzelnen Spec-Review-Commit**
  (`dab4e45`, `2ea534b`, ...): die sind Politur und nicht
  Slice-relevant.
- **Tranchen-Pläne**: Bootstrap-Slices sind nicht in Tranchen
  unterteilt.

## Bezug

- Auslösende Spec: [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003-planning-lifecycle) Lifecycle,
  [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin) Carveout-Disziplin.
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle + M1–M2d-Artefakt-Spalte).
- Vorausgegangene Re-Phasierung: Commit `3b9f4eb` (M3-Anker-Triage)
  hatte den Slice von „M3-T5" auf „Later" gesetzt.
