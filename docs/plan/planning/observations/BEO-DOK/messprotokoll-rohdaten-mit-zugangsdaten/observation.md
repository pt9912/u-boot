# Messprotokoll-Rohdaten tragen Zugangsdaten

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — das Messprotokoll liegt als Evidence im Repo; was in den Rohdaten steht, steht damit im Verlauf.

Ein Messskript gibt die Ausgabe eines Werkzeugs unverändert in eine Rohdatei, und diese enthält einen Zugang: eine vorsignierte Download-URL mit Schlüsselkennung und Signatur, ein Token, einen privaten Schlüssel. Die Rohdatei wandert mit dem Messprotokoll ins Repo; auch eine später geschwärzte Stelle bleibt im Git-Verlauf. Kein Gate prüft Rohdaten auf Zugangsdaten; ein Review findet es nur, wenn er die Rohdaten liest.
