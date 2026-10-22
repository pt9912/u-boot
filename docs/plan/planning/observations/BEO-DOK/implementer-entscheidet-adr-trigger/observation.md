# Die Implementation stellt selbst fest, ob der Re-Evaluierungs-Trigger einer Accepted-ADR ausgelöst ist

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Plan benennt eine Rückführung vorab und die ADR ihren Trigger; wer sie feststellt, regelt die Rollenvorgabe ([`harness/roles.md`](../../../../../../harness/roles.md), *Implementation*: ein ADR- oder Spec-Konflikt wird nicht pragmatisch übergangen, es entsteht ein Übergabe-Artefakt an Architect oder Planner).

Ein Test der Umsetzung erfüllt den Wortlaut eines Re-Evaluierungs-Triggers einer `Accepted`-ADR und einer im Plan vorab benannten Rückführung. Die Implementation schreibt in den Slice, der Trigger sei nicht ausgelöst („Entscheidung des Implementers“), statt dem Architect den Befund mit ihrer Lesung vorzulegen. Bis ein Artefakt des Architects vorliegt, tragen zwei `Accepted`-ADRs einen widerlegten Satz, die DoD-Zeile zur Rückführung ist nicht abhakbar, und die Übergabe an die Folge-Slices fehlt. Kein Gate prüft, wer einen Trigger feststellt; der Review fand es.

**Trigger:** ein zweiter Vorgang, in dem die Umsetzung einen Trigger oder eine Rückführung selbst als „nicht ausgelöst“ führt, oder die Entscheidung des Projektinhabers, die Feststellung in die Vorlagen der Pläne zu schreiben.
