# Pin an drei Stellen ohne Abgleich

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — Go und `golangci-lint` sind Werkzeug-Pins der Build-Verträge (`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`).

Die Versionen von Go und `golangci-lint` stehen an drei Stellen: `Dockerfile`, `repo.mk` und `.devcontainer/Dockerfile` (dort mit SHA-256 je Architektur). Eine Hebung an einer Stelle lässt die beiden anderen stehen, und kein Sensor meldet die Abweichung; die Stellen tragen nur einen Kopplungskommentar „Hebung betrifft drei Stellen“. Das Modul für den Abgleich liefert d-check (Issue 14 in `pt9912/d-check`, laut Plan des Vorgangs); im Repo gibt es dazu kein Skript, der Projektinhaber will keines.

Verwandt, nicht gleich: [`BEO-WZG/kopplung-nicht-nachgezogen`](../kopplung-nicht-nachgezogen/observation.md) (dort ist die Abweichung eingetreten, hier fehlt nur der Sensor).
