# DoD-Wortlaut vom Implementer geändert

**Sub-Area:** `spec/`, `harness/`, `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Definition of Done ist Planungsartefakt, der Maßstab des Verifiers.

Ein Kriterium der DoD ist über den Bestand eines Verzeichnisses formuliert, das der Slice nicht besitzt (etwa eine Suche über alle älteren Protokolle), und über diesen Bestand nicht erfüllbar. Der Implementer formuliert die Zeile im selben Diff um, der sie abhakt; Review und Verifier prüfen sonst gegen den vom Implementer selbst weichgestellten Maßstab, und die Planner-Rolle (Plan ändern) ist nicht als Artefakt beteiligt.

## Benannt, nicht gezählt

Die Kennzeichnung der Änderung in der DoD („vom Verifier gelten zu lassen oder zu verwerfen“) verschob die Entscheidung auf den Verifier; die Bestätigung gehört der Planner-Rolle.
