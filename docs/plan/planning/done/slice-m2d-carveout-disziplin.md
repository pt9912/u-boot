# Slice M2d: Carveout-Disziplin

> **Status:** Done
> **DoD:** Commit `5b92a97` (feat) + `3dc0467` (12 Review-Findings) + `0124742` (Discipline-Gap-Close) + `b51a4aa` (Re-Phase 5 Slices)
> **Retro-Plan:** Retroaktiv geschrieben 2026-05-27 (siehe [`slice-m3-retroaktive-slice-plaene`](slice-m3-retroaktive-slice-plaene.md))

## Auslöser

User-Feedback bei M2c-Review (2026-05-21): *„wichtig dass keine
Carve-Outs vergessen werden. Immer gleich dazu einen Plan schreiben."*

Bis M2c waren über die Bootstrap-Slices hinweg eine Reihe bewusster
Aufweichungen entstanden (Coverage-Schwellwert 0, leerer
`depguard`-Regel-Block, leerer `gomodguard_v2`-Block, `forbidigo.msg`
referenziert nicht-existenten Logging-Port, ADR-0004-Folgepunkte wie
Image-Publish/Trivy/Branch-Protection, ...), ohne disziplinierte
Plan-Spur. Risiko: vergessen werden, sobald das Projekt fachlich wächst.

## Lieferumfang

- **Neue Spec-Anforderung `LH-FA-PROJDOCS-005` (MVP-Pflicht)**: jeder
  temporäre Carveout (Bootstrap-Schwellwert, leerer Regelblock,
  prospektive Doku-Phrase, ADR-Folgepunkt) bekommt **parallel** zur
  Entstehung einen Slice-Plan in `open/` und einen Eintrag in
  `carveouts.md`. Permanente Carveouts kommen ins Inventar **ohne** Plan,
  aber mit Begründung. `LH-MVP-001` ergänzt; Traceability-Matrix-Zeile
  angelegt.
- **Master-Inventar** `docs/plan/planning/in-progress/carveouts.md`:
  zwei Sektionen (temporär = Plan-Pflicht, permanent = nur Begründung).
  Initial 11 Einträge, davon 9 mit Plan-Verweis; weitere kamen
  fortlaufend dazu.
- **`planning/README.md`** um die zweite Master-Datei-Ausnahme
  (`carveouts.md`) erweitert.
- **7 neue Slice-Pläne** in `open/` (ergänzten den bestehenden
  markdown-link-validator):
  - `slice-m3-coverage-threshold-aktivieren.md` (jetzt in `done/`),
  - `slice-m3-depguard-aktivierung-verifizieren.md` (jetzt in `done/`),
  - `slice-v1-gomodguard-rules.md` (re-phasiert auf M3-followup, jetzt
    in `done/`),
  - `slice-v1-logging-port.md` (re-phasiert auf M4, in `open/`),
  - `slice-v1-release-pipeline.md` (in `open/`),
  - `slice-v1-branch-protection-checkliste.md` (absorbiert in
    `slice-v1-release-pipeline.md`, Datei gelöscht),
  - `slice-v1-docker-integrationstests.md` (re-phasiert auf M6,
    umbenannt zu `slice-m6-docker-integrationstests.md`, in `open/`).
- **`slice-v1-retroaktive-slice-plaene.md`** (re-phasiert auf Later,
  umbenannt zu `slice-m3-retroaktive-slice-plaene.md`) — meta-Slice für
  M1/M2/M2b/M2c/M2d-Retro-Pläne; mit dem aktuellen Commit (2026-05-27)
  abgeschlossen.

## Akzeptanz

- `carveouts.md` existiert, gepflegt, mit getrennten Sektionen
  (`temporär` / `permanent`).
- `LH-FA-PROJDOCS-005` ist MVP-Pflicht und in Traceability-Matrix.
- Jeder bewusste Carveout in `.golangci.yml`, Dockerfile, Makefile,
  Specs hat einen Plan-Verweis in `carveouts.md`.
- Memory-Anker [[feedback-carveouts-need-plans]] (User-Vorgabe
  2026-05-21) festgehalten.

## Bezug

- Auslösende Spec: `LH-FA-PROJDOCS-005` (in M2d entstanden).
- Vorgänger: [`slice-m2c-ci-pipeline`](slice-m2c-ci-pipeline.md).
- Nachfolger: M3 (`u-boot init`-Flow) — der erste fachliche Slice unter
  voller Carveout-Disziplin.
- Folge-Commits nach Review: `3dc0467` (12 Review-Findings),
  `0124742` (Discipline-Gap: jeder Carveout zu Slice, jeder Slice zu
  Roadmap), `b51a4aa` (5 Carveout-Slices auf natürliche Milestones
  re-phasiert).
