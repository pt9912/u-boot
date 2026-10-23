# Ein Docker-Test mit Bau im Test scheitert lokal am Netzweg zu einem externen Host, ohne dass der Code sich geändert hat

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Test baut im Lauf ein Image, dessen `apt-get install` Pakete von einem externen Host zieht (`TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate`, Pakete von `download.docker.com`).

Der Bau im Test hängt vom Netzweg zum Paket-Host ab. Bricht die TLS-Verbindung mitten im Abruf ab (`SSL routines::unexpected eof while reading`, `E: Failed to fetch …`), färbt sich der Test rot, obwohl weder Test noch Generator noch Image-Vorlage sich geändert haben. Der Einzellauf wiederholt den Fehler oder besteht; ein `curl` auf dasselbe Paket im Container lief durch. Kein Gate trennt einen Netzfehler dieser Art von einem Fehler des Codes, und die Rate ist nicht bekannt.

**Trigger:** ein zweiter Vorgang, in dem derselbe oder ein anderer Test mit Bau im Test am Netzweg zu einem externen Host rot wird (dann eine Entscheidung, ob der Abruf im Test wiederholt wird oder der Host als bekannte Störquelle in `make test-docker` ausgewiesen ist); ein roter Lauf dieses Tests in der CI. Leser ist der Verifier jedes Slice, dessen Docker-Lauf diesen Test enthält: Er prüft die Ursache über die Dateiliste des Diffs, nicht über eine Wortsuche.
