# Messprotokoll liegt nur außerhalb des Repos

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Messung ist Beleg einer Aussage in Plan, Handbuch oder Evidence; der Ort der Skripte und Rohdaten entscheidet, ob ein Dritter sie wiederholen kann.

Ein Messslice nennt Messstand und Ergebnis in Plan und Evidence, hält aber Skripte und Rohdaten in einem Verzeichnis außerhalb des Repos (etwa im Home-Verzeichnis des Implementers). Der Befund steht dann im Text, nicht im Repo: Fälle sind nicht reproduzierbar, und Artefakte der Rohdaten (etwa eine Fehlermeldung eines Hilfsprogramms im Messlauf) sind ohne Lesehilfe nicht deutbar. Kein Gate prüft den Ort der Rohdaten; ein Review fängt es nur, wenn er die Dateien außerhalb liest.
