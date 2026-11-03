# Beobachtung — `BEO-WZG/tree-hash-collation`

**Zähler-Form (Stand):** 1 Vorgang (`evidence/`) — unter der Schwelle.

**Klasse:** der Gate-Nachweis-Hash ist umgebungsgetragen, wenn die
Sortier-Lesart der Dateiliste nicht gepinnt ist: `tools/harness/working-tree-hash.sh`
pinnte die Collation nicht; der Host dieses Kontexts trägt uutils-`sort`
(0.10.0), das Dev-Container-Image GNU-`sort` (9.7) — die Sort-Zeilen weichen
byte-weise ab (786 Zeilen), und der im Container gestempelte Nachweis liest
am Host anders (nicht-streng verschleiert, streng falsch-rot). Berichtigt mit
einemzeiligen `LC_ALL=C`-Pin; die Hash-Gleichung je Beweis (beide Kontexte
byte-gleich nach dem Pin).

**Sensor-Form:** der Pin ist selbst die Verkörperung — für Wiederholungen
gilt: jede neue Sortier-/Hash-Stelle im Werkzeug-Bestand liest denselben
Pin; ein Sensor über neue Hash-Stellen ist kein Träger (die Stelle ist
einzeln gelesen).

**Zugehörige Vorgänge:** `evidence/slice-v1-bau-nur-devcontainer.md`.