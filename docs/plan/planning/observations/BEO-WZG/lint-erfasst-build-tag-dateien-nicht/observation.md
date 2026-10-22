# `make lint` prüft Dateien mit Build-Tag nicht

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — das Lint-Profil (`.golangci.yml`, `make lint`) gehört zum Werkzeug-Verbund und ist keine einzelne Go-Sub-Area; die Befunde liegen in zwei Go-Sub-Areas (ADP, E2E), eine Aufteilung nach Fundort ergäbe zwei Pfade für dieselbe Beobachtung.

`.golangci.yml` führt keinen Build-Tag (`run.build-tags` fehlt); `make lint` und `make test` lesen Dateien mit `//go:build docker` deshalb nicht. Ein Verstoß gegen das Lint-Profil in einer solchen Datei (die Docker-Tests liegen alle dort) bleibt bis zum Review unsichtbar, der ihn nur findet, wenn er das Profil mit dem Tag in einer Kopie fährt. Kein Gate erfasst die Dateien; die Aufnahme des Tags wäre eine Erweiterung des Gates und braucht einen Plan-Anker ([`AGENTS.md`](../../../../../../AGENTS.md) §3.6).

**Trigger:** die Entscheidung des Projektinhabers, den Tag in das Lint-Gate aufzunehmen (dann sind die vorhandenen Befunde zu beheben oder als Carveout zu führen), oder ein weiterer Verstoß in einer getaggten Datei, den nur ein Review findet.
