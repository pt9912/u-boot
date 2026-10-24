# Eine Regel prüft beim Lesen der Datei, aber nicht beim Schreiben desselben Wertes

**Sub-Area:** `internal/hexagon/application` (Kürzel APP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Lesestelle `validateDevcontainer` und die Schreibstelle `revalidateConfigDomain` der Konfiguration liegen dort (`devcontainer_features.go`, `config.go`).

Die Schreibpipeline von `config set` sagt in Stufe 5, das Ergebnis müsse sich wieder als Projektkonfiguration einlesen lassen. Der Code las es nur durch `yaml.Unmarshal` zurück; eine Regel, die erst die Lesestelle (`generate`, `doctor`, jedes weitere `config set`) prüft, lief beim Schreiben nicht. Das Setzen meldete Erfolg, das nächste Lesen sperrte das Projekt (`devcontainer.features.<feature>.version` ohne `source`, Befund M-2 des Validators zu `gh` und `glab`). Geschlossen hat es erst ein Slice, der die Regel an **einer** Stelle hält und beide Pfade fragen lässt.

**Trigger:** der nächste Slice, der ein Feld oder einen Pfad zu `config set` hinzufügt oder eine Regel in `validateDevcontainer*` ergänzt. **Grenze:** Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der die Schreibpipeline oder die Validierung der Projektkonfiguration berührt. Kein Carveout: Der bekannte Fall ist behoben; es gibt keinen offenen Mangel, nur die Klasse.
