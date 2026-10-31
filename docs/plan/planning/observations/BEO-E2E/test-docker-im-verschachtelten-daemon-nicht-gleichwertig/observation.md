# `make test-docker` im verschachtelten Daemon ist nicht gleichwertig zum Docker-Modus am Host

**Sub-Area:** `internal/e2e` (Kürzel E2E, Greenfield laut [`harness/conventions.md`](../../../../../../harness/conventions.md#mr-006)) — die Docker-Tests (Build-Tag `docker`) sind die Test-Infrastruktur, deren Ergebnis vom Daemon abhängt, gegen den sie laufen.

`make test-docker` im lokalen Modus des Devcontainers läuft gegen den verschachtelten Docker-Daemon des Containers. Dort schlagen 9 Tests fehl, die im Docker-Modus am selben Host bestehen oder `SKIP` melden: bei vier Tests meldet das Startscript Exit 0 statt 11 (`kernel.apparmor_restrict_unprivileged_userns=1`, der Test überspringt sonst), bei fünf ist die Ursache nicht geklärt (Bau im verschachtelten Daemon bricht ab, der Test schlägt dennoch fehl). Kein Gate trennt einen Fehlschlag des Daemons von einem Fehler des Codes.
