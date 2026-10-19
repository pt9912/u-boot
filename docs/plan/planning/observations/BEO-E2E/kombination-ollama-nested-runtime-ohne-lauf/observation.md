# Katalog-Eintrag mit Dockerfile-Block ohne Lauftest in Kombination mit `nestedRuntime`

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — der Docker-Test des Feature-Slice baut `ollama` allein; die Kombination mit `devcontainer.sandbox.nestedRuntime` (`podman`, `docker`) hat keinen Docker-Test.

Der Block `ollama` steht im erzeugten Dockerfile hinter dem Block der Container-Runtime und entfernt `zstd` mit `apt-get purge -y --auto-remove`, was andere automatisch installierte Pakete entfernen kann. Das Dockerfile der Kombination ist als Text gepinnt (Golden Case), der Verifier hat Bau und Versionsprobe einmal gemessen (`ollama --version` `0.40.6`, `podman version 6.1.4` bzw. `Docker version 29.9.1`, im Worktree, nicht committet); der **Lauf** von Podman oder Docker in der Kombination und ein Build im Container sind nicht gemessen. Das Handbuch sagt über die Kombination nichts zu.

**Trigger:** ein Slice, der eine Zusage über die Kombination macht, oder ein weiterer Katalogeintrag mit Dockerfile-Block (zweiter Beleg für Blöcke, die einander im selben Dockerfile berühren). Leser ist die Slice-Planung (§8, Sichtungs-Schritt) jedes Slice, der `internal/e2e` oder die Dockerfile-Blöcke berührt.
