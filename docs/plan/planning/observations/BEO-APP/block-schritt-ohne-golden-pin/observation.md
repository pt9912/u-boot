# Der Golden Case eines Dockerfile-Blocks pinnt, was im Block steht, nicht was darin fehlen muss

**Sub-Area:** `internal/hexagon/application` (Kürzel APP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Golden Cases der Katalog-Blöcke liegen dort (`generate_forge_cli_test.go`, Muster `generate_ollama_test.go`); die Blöcke selbst in `templates/devcontainer/Dockerfile.tmpl`.

Ein Golden Case prüft Zeichenfolgen des erzeugten Dockerfiles (Pin, Prüfsumme, Reihenfolge der Schritte, Abwesenheit von `ENV` und `auth login`). Zwei Änderungen am Block, die eine Zusage oder das Aufräumen berühren, färbten keinen Golden: `chown vscode:vscode /usr/local/bin/gh` im `gh`-Block (die Zusage „Binary-Eigentümer `root`“ hängt dann allein am Docker-Test, der nicht in `make gates` läuft) und `cd /` → `cd /tmp` vor `rm -rf "$work"` (alle vier Blöcke der Datei, funktional folgenlos). Gefunden haben es die Mutationen von Verifier und Reviewer, nicht die Tests.

**Trigger:** der nächste Slice, der einen Dockerfile-Block ändert oder hinzufügt (Pin-Hebung, weiterer Katalogeintrag mit Block). **Grenze:** Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der `internal/hexagon/application` berührt. Kein Carveout: Die Zusage ist durch den Docker-Test belegt, es fehlt allein der Wächter im inneren Gate.
