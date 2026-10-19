# Aufräum-Skript eines Messprotokolls löscht über das Präfix der Messung hinaus

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — Messskripte laufen auf dem Docker-Host des Projektinhabers, der von anderen Arbeiten mitbenutzt wird; ein Aufräum-Schritt, der anders als nach dem Namen auswählt, kann fremde Objekte entfernen.

Ein Aufräum-Schritt entfernt neben den Objekten mit dem Präfix der Messung Objekte nach einem anderen Kriterium (etwa jedes unbenannte Image, das in der Liste vorher fehlt, oder jeden beendeten Container, dessen Befehl ein Merkmal der Messung trägt). Auf einem geteilten Host entstehen solche Objekte auch durch fremde Arbeiten; die Aussage „fremde Objekte wurden nicht entfernt“ folgt dann aus der Liste vorher und nicht aus einer Schutzgrenze im Skript. Kein Gate prüft Aufräum-Skripte; ein Review findet es nur, wenn er das Skript liest.

**Trigger:** ein weiterer Messslice oder ein Skript, das Docker-Objekte entfernt; zweites Auftreten in einem anderen Vorgang (dann die Frage, ob eine Regel „Löschen nur nach Präfix im Namen, ohne zweites Kriterium“ in die Vorlage der Messprotokolle gehört). Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der ein Messprotokoll mit Aufräum-Schritt anlegt.
