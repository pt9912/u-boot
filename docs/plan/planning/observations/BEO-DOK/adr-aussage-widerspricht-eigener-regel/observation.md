# Zwei Aussagen derselben ADR widersprechen sich, die Regel der ADR und ihre Tabelle oder zwei Stellen zur selben Eigenschaft

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — eine ADR trägt Regeln, Tabellen und Entscheidungen, die sich aufeinander beziehen; nach `Accepted` ist sie unveränderlich ([`AGENTS.md`](../../../../../../AGENTS.md) §3.5).

Eine ADR legt eine Regel fest (zum Beispiel für die Pfad-Spalte einer Zusage-Tabelle oder für die Sichtbarkeit unter `--quiet`) und führt an anderer Stelle eine Zeile, die gegen sie steht, oder zwei Aussagen zur selben Eigenschaft ohne Auflösung, oder eine „Menge“, die ihre Tabelle nicht deckt. Kein Gate liest eine ADR gegen sich selbst; der Reviewer findet es durch Satz-für-Satz-Lesen, und eine Fix-Runde kann den Widerspruch neu erzeugen.

**Trigger:** ein dritter Vorgang, oder die Entscheidung des Projektinhabers, eine Selbstlese-Prüfung (Regel gegen Tabelle, Kontext gegen Entscheidung) in den Reviewer-Skill zu schreiben.
