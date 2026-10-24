# Eine ADR verweist auf Welle, Roadmap oder Slice und trägt Aufgaben des Planners

**Sub-Area:** `docs/plan/` (Kürzel DOK, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die ADR steht im Referenzmodell ([`ADR-0013`](../../../../adr/0013-dokumentationsreferenzmodell.md)) über der Planung; ihr Verweis zeigt aufwärts, und sie ist nach `Accepted` unveränderlich ([`AGENTS.md`](../../../../../../AGENTS.md) §3.5).

Ein Entwurf der ADR verlinkt die Welle-Datei als Eingabe einer Entscheidung und trägt in der Folgepflicht Aufgaben des Planners („Closure-Kriterien angleichen“, „Zeilen der Folge-Slices und Roadmap nachziehen“, die Reihenfolge der Folge-Slices). Wandert die Welle nach `done/`, zieht der Verschiebe-Befehl den Verweis in einer unveränderlichen Datei nach oder lässt ihn toten Link werden. Der Link-Check des Repos meldet die Kante nicht (Review: `make docs-check` 0 Befunde); die ungelinkte Nennung von Rollen und Plan-Begriffen findet nur eine Wortsuche.

**Trigger:** ein dritter Vorgang, oder die Entscheidung des Projektinhabers, die Kante ADR nach Planung als Wächterregel zu führen.
