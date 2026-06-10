# Slice V1: Markdown-Link-Validator

> **Status:** Done
> **DoD:** Commit `2f8242c`

## Auslöser

Relative `[text](path)`-Links in `docs/`, `spec/` und Root-READMEs sind
heute nicht maschinell geprüft (M2-Review #11, [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin)).
Nach den vielen Slice-Moves in M3-Closure und M3-Anker-Triage ist das
Drift-Risiko konkret: heutiger erster Lauf hat 3 Findings produziert,
darunter zwei echte Drift-Bugs (Pfad nach Slice-Rename, Pfad nach
Slice-Move).

Eigener Slice statt Carveout-Sammlung, weil das Tool stdlib-only ist
und sich gut Docker-only verpacken lässt.

## Aufhebung

`tools/check_refs.py` adaptiert aus `c-hsm-doc/tools/check_refs.py`
(gleiche Build-Familie, gleiche Docker-only-Philosophie). Scant alle
`.md` unter `docs/`, `spec/` und Root, prüft jeden `[text](path)`-Link
(inkl. Image-Refs `![alt](src)`), refused Symlinks defensiv und
ignoriert externe Refs (`http://`, `mailto:`, `#anchor`).

Verbesserung gegenüber c-hsm-doc-Original: Fenced-Code-Blocks
(```` ``` ```` und `~~~`) werden korrekt übersprungen — sonst hätte
`docs/archive/README.md` einen Falsch-Positiv-Treffer im
Markdown-Beispiel-Block produziert.

## Geliefert

- `tools/check_refs.py` (172 Zeilen, stdlib-only).
- `make docs-check`-Target via Docker-encapsulated Python
  (`PYTHON_VERSION ?= 3.13-slim`); kein Host-Python-Requirement.
- `docs-check` in `gates` aufgenommen — jede `make gates`-Iteration
  validiert die Doc-Links als vierte Gate-Stufe.
- Drei initiale Drift-Findings gefixt:
  - `docs/plan/adr/0004-ci-system.md:71`: stale Pfad
    `slice-v1-docker-integrationstests.md` → `slice-m6-docker-...`
    (Slice-Rename bei Phase-Re-Anker).
  - `docs/plan/planning/done/[slice-m3-coverage-threshold-aktivieren.md](slice-m3-coverage-threshold-aktivieren.md):40`:
    same-dir-Link auf `carveouts.md` → `../in-progress/carveouts.md`
    (Datei lebt dauerhaft in `in-progress/`).
  - `docs/archive/README.md:14`: Falsch-Positiv durch Beispiel im
    fenced-markdown-Block; via Validator-Patch (Fenced-Code-Stripping)
    aufgelöst, kein Doc-Edit nötig.
- `carveouts.md`: Doku-/Link-Drift-Zeile entfernt.
- Roadmap: Carveout-Auflösungs-Slice → Done, Phase „V1-vorgezogen".
- READMEs: Carveout-Count 13 → 12.

## Out of Scope

- **Kennungs-Auflösung (`LH-*`-Querverweise)** wie c-hsm-doc-Vorlage:
  spätere Erweiterung, falls semantische Spec-Cross-Refs gepflegt
  werden müssen.
- **§-Sektions-Verweise** (z.B. „§4 Architektur"): nicht prüfbar ohne
  Heading-Index. Out of Scope.
- **Indented Code Blocks** (4-Space-eingerückte Markdown-Code-Blocks,
  ohne Fence): heute selten verwendet, kein Falsch-Positiv-Risiko
  bisher. Bei Bedarf erweitern.
- **Markdown-Reference-Style-Links** (`[text][ref]` mit `[ref]: url`
  separat): u-boot verwendet nur inline-Style.

## Bezug

- Auslösende Konvention: [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005-carveout-disziplin) Carveout-Disziplin +
  M2-Review #11.
- Aufhebung dokumentiert in: [`carveouts.md`](../in-progress/carveouts.md)
  (Zeile entfernt) und [`roadmap.md`](../in-progress/roadmap.md)
  (Carveout-Auflösungs-Slice-Tabelle).
- Vorlage: `../c-hsm-doc/tools/check_refs.py`.
