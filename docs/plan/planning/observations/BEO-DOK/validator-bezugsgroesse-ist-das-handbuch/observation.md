# Das Closure-Kriterium des Validators nimmt das Handbuch als Bezugsgröße, das Zielbild liegt in `spec/`

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Welle-Datei legt im Closure-Kriterium *Validator-Bericht* fest, wogegen der Validator das Ergebnis liest.

Ein Closure-Kriterium der Welle lässt den Validator den Nutzerpfad „aus Benutzerhandbuch und Hilfetext allein“ prüfen. Das Handbuch beschreibt den Ist-Zustand, ist abgeleitet und nicht normativ; gegen sich selbst gelesen kann es das Zielbild nicht verfehlen, und ein Bedarf, den Anforderung und Spezifikation ausschließen, bleibt unsichtbar. Das Zielbild tragen die Dateien in `spec/` (Anforderung und Spezifikation) und der Anlass des Auftraggebers; der Validator liest das Handbuch als Teil des Ergebnisses.
