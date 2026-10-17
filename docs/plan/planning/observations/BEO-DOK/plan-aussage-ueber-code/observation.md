# Plan-Aussage über Code stimmt nicht

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Planung ist der Ort, an dem die Aussage steht.

Ein Slice-Plan, eine Closure-Notiz oder eine Spec nennt Dateien, Funktionen, Werte oder Verhalten des Codes, die der Autor nicht an der Stelle im Code gelesen hat: ein Item, das es nicht gibt, eine falsche Testdatei, eine Datei im Diff, die die Plan-Tabelle (§3) nicht listet. Kein Gate prüft Plan-Aussagen; Review oder Verifier finden sie erst nach dem Bau.
