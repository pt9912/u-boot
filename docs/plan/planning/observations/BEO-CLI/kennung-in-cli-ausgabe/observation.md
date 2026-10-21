# Anforderungs-Kennungen in der Ausgabe der CLI

**Sub-Area:** `internal/adapter/driving/cli` (Kürzel CLI, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Hilfe-Texte und die Textform der Diagnosen stehen im Adapter.

Hilfe und Textausgabe der CLI tragen Anforderungs-Kennungen, die ein Nutzer ohne Lastenheft nicht auflösen kann: Warnungen und Ausweisungen mit dem Code als Präfix, Hilfe-Zeilen der Form „Exit codes ([`LH-FA-CLI-006`](../../../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes))“. Die Nutzer-Doku trägt nach [`AGENTS.md`](../../../../../../AGENTS.md) §5 Regel 6 keine Kennungen; für die Ausgabe der CLI gibt es keine solche Regel.
