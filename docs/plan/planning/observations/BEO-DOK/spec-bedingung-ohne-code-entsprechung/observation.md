# Spec- oder Handbuch-Bedingung ohne Code-Entsprechung

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — Spezifikation und Handbuch sind der Ort, an dem die Bedingung steht; der Code trägt sie nicht.

Spezifikation, Handbuch oder CHANGELOG nennen eine Bedingung („nur auf einem Linux-Host mit lokaler Engine“, „im Devcontainer erreicht der Weg keine Engine“, „`init` weist den Fallback aus“), die eine Eigenschaft des Codes behauptet; der Code prüft sie nicht, oder der Pfad ist in Produktion nicht erreichbar, und kein Test pinnt sie. Unterschied zu [`BEO-DOK/plan-aussage-ueber-code`](../plan-aussage-ueber-code/observation.md): Dort stimmt eine Aussage im **Plan** nicht; hier steht die Bedingung im **Vertrag** (Spezifikation, Handbuch), und der Befund ist das Fehlen von Code oder Test, der sie trägt.
