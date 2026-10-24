# Die Folgepflicht einer ADR nennt nicht die Vertragsstelle, auf die sich ihre Aussage stützt, oder nicht den Satz, der falsch wird

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die *Folgepflicht* einer ADR benennt den Bedarf an Lastenheft, Spezifikation, Handbuch und Code; der Wortlaut folgt in den Folge-Slices.

Die Folgepflicht listet die Stellen, die die Entscheidung *ergänzt*, und übersieht zwei andere Arten: die Vertragsstelle, auf die sich eine Aussage der ADR *stützt* (eine Sichtbarkeit unter `--quiet` stützt sich auf eine Aufzählung im Lastenheft, die den neuen Fall nicht nennt), und den Satz, der durch die Entscheidung *falsch wird* (eine Klammer oder ein Doc-Kommentar, nicht der Satz daneben). Kein Gate prüft die Liste gegen die Aussagen; der Folge-Slice folgt der Liste und findet die Lücke erst am Konflikt mit dem Vertrag oder am Test.

**Trigger:** ein dritter Vorgang, oder die Entscheidung des Projektinhabers, die Prüfung „jede Aussage mit Vertragsstützung und jeder durch die Entscheidung falsch werdende Satz steht in der Folgepflicht“ in den Reviewer-Skill zu schreiben.
