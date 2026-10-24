# Eine Annahme des Plans über den Zustand eines Fremdsystems hat keine Messstelle

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Slice-Plan nennt Annahmen und Rückabwicklungen.

Ein Plan nimmt an, dass ein Schritt einen Zustand in einem System erzeugt, das der Slice nicht besitzt (hier: die Anmeldung eines Ollama im Container trägt einen Geräteschlüssel im Konto des Projektinhabers ein, `signout` oder der Nutzer baut ihn zurück), und plant die Rückabwicklung und deren Bestätigung auf dieser Annahme, ohne den Zustand vor und nach dem Schritt mit Zeitangabe zu lesen. Danach lässt sich ein Befund im Fremdsystem (ein Eintrag, der fehlt, oder einer, der später erscheint) keinem Schritt zuordnen. Kein Gate prüft, ob eine Annahme über Fremdzustand eine Messstelle trägt.

## Benannt, nicht gezählt

Kein weiteres Vorkommen außerhalb des ersten Vorgangs.
