# Pull mit verschachteltem Docker direkt nach dem Start der Egress-Beschränkung läuft in eine Zeitüberschreitung

**Sub-Area:** `internal/hexagon/application` (Kürzel APP, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — das Startscript der Beschränkung (Template des Generators) und der Start des verschachtelten Docker-Daemons folgen einander im selben Container.

Mit `nestedRuntime: docker` und Egress-Restriktion lief der Pull eines Images direkt nach dem Start der Beschränkung in einem von zwei vollständigen Läufen in `i/o timeout`; ein Versuch 15 Sekunden später gelang, im anderen Lauf gelang der erste Versuch. Die Ursache ist nicht untersucht; das Handbuch weist sie als nicht untersucht aus. Messdaten: [`rohdaten-lauf2-07-docker.txt`](../../../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-lauf2-07-docker.txt) (Zeitüberschreitung) und [`rohdaten-07-docker.txt`](../../../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-07-docker.txt) (erster Versuch gelungen).

**Trigger:** ein zweites Auftreten in einem anderen Vorgang (dann eine Diagnose als eigener Slice); ein Slice, der die Reihenfolge von Beschränkung und Daemon-Start oder die Namensauflösung des Containers ändert. Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der Egress und `nestedRuntime: docker` zusammen berührt.
