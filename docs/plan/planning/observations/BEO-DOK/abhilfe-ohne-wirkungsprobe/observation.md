# Eine Abhilfe in Meldung und Doku wird nicht ausgeführt, bevor sie im Text steht

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Abhilfe steht in der Meldung (Code), im Handbuch, im CHANGELOG und in der Spezifikation; keiner dieser Orte führt sie aus.

Eine Fehlermeldung nennt eine Abhilfe („schalte das Feature ab oder …“), und Handbuch, CHANGELOG und Spezifikation geben sie weiter. Die Abhilfe wird nicht ausgeführt, bevor sie im Text steht: Ein genannter Weg kann ohne Wirkung sein (der zweite Weg hob den Fehler nicht auf), die Texte können die Schritte in verschiedener Reihenfolge nennen, und ein Schritt, der nur manchmal nötig ist und ein Risiko trägt, kann als ausführbare Zeile eines Codeblocks stehen. Die Abhilfe steht an mehreren Stellen, und jede Fix-Runde muss alle nachziehen. Kein Gate prüft sie; der Reviewer fand die Lücke mit einer Probe, die er selbst ausführte.
