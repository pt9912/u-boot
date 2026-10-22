# Die Rücksetzung einer Mutation per `git checkout` verwirft ungesicherte Änderungen

**Sub-Area:** `tools/`, `scripts/` (Kürzel WZG, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Mutationsläufe der Closure-Pflichten gehören zum Werkzeug-Verbund, nicht zu einer Go-Sub-Area.

Ein Mutationslauf ändert Quelltext im Arbeitsbaum und setzt ihn mit `git checkout -- <Pfad>` zurück (der Slice führt `git checkout -- internal cmd`, danach `git diff -- internal cmd` leer). Der Befehl stellt den Stand des letzten Commits her, nicht den Stand vor der Mutation: Eine Änderung, die zur Zeit der Mutation nicht committet war, geht dabei verloren, ohne dass `git` es meldet. Kein Gate prüft, ob der Arbeitsbaum vor der Mutation sauber war; die Verification des Slice hielt den Zustand ausdrücklich fest („Arbeitsbaum vorher sauber“), der Implementer nicht.

**Trigger:** ein zweiter Mutationslauf, der ungesicherte Arbeit verwirft, oder die Entscheidung, die Rücksetzung auf einen Patch (`git apply -R`) oder eine Kopie des Baums zu stützen und die Vorbedingung „Arbeitsbaum sauber“ in die Vorlage der Mutationstabelle zu schreiben. Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der Mutationen fährt.
