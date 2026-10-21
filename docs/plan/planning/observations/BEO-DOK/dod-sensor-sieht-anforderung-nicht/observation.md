# DoD-Sensor sieht die Anforderung des Slice nicht

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Definition of Done des Slice-Plans ist der Ort, an dem der Sensor als Beleg steht.

Ein DoD-Punkt nennt einen Sensor als Beleg für die Anforderung des Slice, dessen Sicht die Anforderung nicht enthält: „`make doc-trace` ohne neue Waise“ für eine Anforderung außerhalb der Belegsicht, die nur die Familien `LH-FA-*` und `LH-QA-*` führt (`.d-check.yml`, Abschnitt `trace`). Der Lauf ist grün und der Punkt abgehakt, ohne dass er etwas über die Anforderung sagt. Kein Gate prüft, ob ein DoD-Sensor die Anforderung sieht.
