# Eine Messung an der Variante belegt nicht die ausgelieferte Umsetzung

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Messung ist Beleg einer Aussage in Plan, ADR oder Handbuch; ob sie die ausgelieferte Form trägt, entscheidet der Abgleich von Messvariante und Umsetzung.

Eine Messung läuft an einer Wegwerf-Variante des Produkts (etwa eine `sed`-Kopie des erzeugten Scripts im Wegwerf-Container), die Umsetzung entsteht danach und weicht in Einzelheiten ab (Konfigurationszeile, Prüfungen, Abbruchzweige). Eine Aussage der Messung („grün mit der Variante“) gilt dann für die Variante; für die Umsetzung gilt sie erst, wenn dieselbe Probe an ihr läuft. Kein Gate prüft, an welcher Form eine Aussage im Plan gemessen wurde.
