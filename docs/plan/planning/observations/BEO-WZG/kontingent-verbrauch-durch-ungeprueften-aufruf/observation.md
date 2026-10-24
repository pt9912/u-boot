# Kontingent verbraucht durch einen Aufruf, dessen Eingabe nicht geprüft war

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — Messskripte zählen und verbrauchen ein begrenztes Kontingent.

Ein Messskript zählt einen Aufruf gegen eine Obergrenze (hier: Prompts an einen Cloud-Dienst), bevor die Eingabe des Aufrufs geprüft ist. Der Aufruf scheitert am Client (leerer Modellname, ein Name, den der Client als ungültig kennt) und verbraucht Kontingent ohne Messwert. Der Zähler und die Obergrenze selbst hielten; die Prüfung der Eingabe stand hinter dem Zählen. Kein Gate prüft die Reihenfolge von Eingabeprüfung und Zählen in einem Messskript.

## Benannt, nicht gezählt

Kein weiteres Vorkommen außerhalb des ersten Vorgangs.
