# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-cr-sandbox-fallback-ausweisung`](slice-cr-sandbox-fallback-ausweisung.md) | CR: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) auf das host-seitig Ermittelbare präzisieren (Daemon- und Egress-Fallback gibt das Startscript im Container aus); wartet auf die Wortlaut-Freigabe des Projektinhabers |
| [`slice-v1-sandbox-fallback-host-ausweisung`](slice-v1-sandbox-fallback-host-ausweisung.md) | Daemon-, Egress- und `vfs`-Fallback host-seitig ermitteln und in `u-boot doctor` und, wo ein Befehl sie feststellt, in der Befehlsausgabe ausweisen ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Lastenheft unverändert); wartet auf die Folge-ADR des Architect zum Kanal (Rollenfolge Planner → Architect → Implementer → Reviewer → Verifier) |
