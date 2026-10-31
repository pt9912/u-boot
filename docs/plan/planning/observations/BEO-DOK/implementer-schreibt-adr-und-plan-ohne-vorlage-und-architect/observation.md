# Der Implementer schreibt ADR und Slice-Plan ohne Vorlage und ohne Architect

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — ADR und Slice-Plan sind Artefakte der Rollen Architect und Planner ([`harness/roles.md`](../../../../../../harness/roles.md)); die Vorlagen der Baseline sind die Form, die beim Anlegen kopiert wird ([`AGENTS.md`](../../../../../../AGENTS.md) §1).

Der Implementer-Kontext schreibt den ersten Entwurf einer ADR und eines Slice-Plans selbst, ohne sie per `cp` aus den Vorlagen anzulegen, und die ADR steht auf `Accepted`, bevor ein Architect sie liest. Der Projektinhaber bemerkt es; danach läuft der Architect nach, und der Planner bringt dessen Report in die Vorlagenform. Kein Gate prüft, aus welcher Rolle und in welcher Form ein Plan-Artefakt entsteht.

Verwandt, nicht gleich: [`BEO-DOK/dod-wortlaut-vom-implementer-geaendert`](../dod-wortlaut-vom-implementer-geaendert/observation.md) (der Implementer ändert den Maßstab des Verifiers) und [`BEO-DOK/implementer-entscheidet-adr-trigger`](../implementer-entscheidet-adr-trigger/observation.md) (der Implementer stellt eine Entscheidung des Architects selbst fest). Gemeinsam ist die Klasse „Artefakt einer anderen Rolle vom Implementer-Kontext“; jede der drei trägt ein anderes Artefakt.
