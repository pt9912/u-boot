# Test liest die Umgebung des Hosts

**Sub-Area:** `internal/adapter/driven` (Kürzel ADP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — ein Adapter-Test, dessen Ergebnis vom Zustand des Hosts abhängt, ist in CI grün und im Devcontainer des Entwicklers rot; der Fehler liegt im Test, nicht im Produkt.

Ein Adapter, der `os.Getenv` oder `git config` liest, antwortet im Test mit dem Zustand des Hosts, wenn der Test ihn nicht pinnt: ein globales `user.name`/`user.email` oder ein gesetztes `GIT_TOKEN` überlagert die erwartete leere Antwort. Das Gate fängt das nur auf einem Host ohne diese Werte (CI); ein Entwickler-Container mit Identität und Token färbt `make gates` rot. Schlägt ein solcher Test fehl und gibt den Wert aus, steht ein Secret im Log.
