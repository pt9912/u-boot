# Kopplung nicht nachgezogen: eine Änderung berührt ein abhängiges Skript nicht

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — dort liegen die Skripte, die an Artefakten oder Namen hängen, die ein anderer Ort erzeugt.

Ein Skript prüft oder bildet etwas aus einer Form, die ein anderes Artefakt erzeugt (Formel-Inhalt, Dateiname eines Pakets). Ändert sich die Form am Erzeuger, zieht niemand das Skript nach, und der Zweig des Skripts, der die Form liest, läuft nur in einem echten Lauf. Kein Test und kein Vorab-Lauf fängt es, wenn der Zweig dort nicht läuft.
