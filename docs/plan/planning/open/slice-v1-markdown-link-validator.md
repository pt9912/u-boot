# Slice V1: Markdown Link Validator

## Motivation

Querverweise in `spec/`, `docs/plan/`, READMEs und ADRs verwenden
relative Markdown-Links. Verzeichnis-Moves im Planning-Lifecycle
(`LH-FA-PROJDOCS-003` — `open/` → `next/` → `in-progress/` → `done/`)
und Archivierungen (`LH-FA-PROJDOCS-004`) ändern Pfade laufend, was zu
toten Links führt.

Vorlage: `grid-gym/tools/check_refs.py` (Markdown-Link-Validator als
docs-check-Stage im Dockerfile, Trigger im `lint`-Stage bzw. eigenem
`docs-check`-Target).

## Scope

- Tooling in `tools/check_refs/` (Go-Programm oder Bash+ripgrep —
  beides akzeptabel).
- Scan über `**/*.md` in `spec/`, `docs/`, Repo-Root-READMEs.
- Auflösung relativer Links gegen das Dateisystem; fehlende Targets
  werden mit Pfad und Quellzeile gemeldet.
- Whitelist für externe Links (http(s)://, `mailto:`, Anker im selben
  Dokument).
- Integration als neuer Stage `docs-check` im `Dockerfile` und
  Make-Target `make docs-check`.
- Aufnahme in `make ci` (V1-Erweiterung von `LH-FA-BUILD-006`).

## Akzeptanzkriterien

- `make docs-check` erkennt einen absichtlich gebrochenen Link in
  einem Smoketest-Fixture und exitiert mit Non-Zero.
- `make docs-check` läuft auf dem aktuellen Repo grün.
- Dokumentation des Tools im `tools/check_refs/README.md`.
- Lastenheft-Eintrag (z. B. `LH-QA-005 – Markdown-Link-Validator`,
  Priorität V1) plus Traceability-Matrix-Eintrag.

## Out of Scope

- Validierung von Links innerhalb von Code-Beispielen (Codeblocks
  bleiben unberührt).
- Anker-Validierung (`#abschnitt`) — kann V2-Erweiterung sein.
- Validierung externer Links (kein Netzwerk-Check).

## Bezug

- Vorlage: `grid-gym/tools/check_refs.py`.
- Triggernder Review-Befund: M2-Review #11 — kein automatischer Schutz
  gegen Drift in Querverweisen.
- Erweitert: `LH-FA-BUILD-006` (Aggregator-Targets), `LH-QA-*`
  (neuer Eintrag).
