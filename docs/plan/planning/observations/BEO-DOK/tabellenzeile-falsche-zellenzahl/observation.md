# Eine Tabellenzeile trägt eine Zelle zu viel oder zu wenig

**Sub-Area:** `docs/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — betroffen ist die Mutations-Matrix in `docs/maintainer/cli-json-contract.md` §7, eine Tabelle mit neun Spalten, in die jeder neue Befehl eine Zeile bekommt.

Ein Diff, der eine Zeile an eine Tabelle hängt, nimmt dabei eine Zelle der Nachbarzeile mit oder fügt eine zu viel an: Die Zeile des Vorgängers hat acht statt neun Zellen, die neue zehn. Die Darstellung verrutscht um eine Spalte oder schneidet die überzählige Zelle ab; der Inhalt (hier jeweils `—`) verrät den Fehler nicht. Kein Gate zählt Zellen einer Markdown-Tabelle; der Review las die Tabelle nicht auf die Zellenzahl, der Verifier fand es.

**Trigger:** ein zweiter Befund dieser Art an einer Tabelle, oder die Entscheidung für einen Zellen-Zähler im Doku-Gate. Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der eine Zeile an eine Tabelle des Kontrakts hängt.
