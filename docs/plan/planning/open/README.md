# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-cr-devcontainer-agent-features-sandbox`](../in-progress/slice-cr-devcontainer-agent-features-sandbox.md) | CR: Sandbox-Regel für `sshd` und Claude Code und Zugangsdaten des Agenten, Sicherheitswarnungen unter `--quiet` ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)); Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), erster Slice; entfällt, wenn der Projektinhaber keine neue Zusage verlangt |
| [`slice-v1-cli-quiet-sicherheitswarnungen`](slice-v1-cli-quiet-sicherheitswarnungen.md) | `doctor --quiet` zeigt die Warn-Einträge, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1); Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), zweiter Slice |
| [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) | Katalog-Schlüssel `sshd` (offizielles Feature), Egress, Sandbox-Regel, Tests, Doku; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), dritter Slice |
| [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) | Katalog-Schlüssel `claude-code` (Feature von Anthropic, Drittanbieter), Egress, Zugangsdaten zur Laufzeit, Tests, Doku; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), vierter Slice |
| [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md) | Spezifikation zu Katalog, Egress und Sandbox-Verhalten beider Features; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), letzter Slice |
| [`slice-v1-sandbox-fallback-ausweisung`](slice-v1-sandbox-fallback-ausweisung.md) | Degradations-Fallbacks Daemon und Egress host-seitig ausweisen (Code) oder [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) präzisieren (CR); Entscheidung beim Projektinhaber |
| [`slice-v1-cli-dry-run-would-formulierung`](slice-v1-cli-dry-run-would-formulierung.md) | Dry-Run: `remove`, `generate` und `config set` sagen „Would …“ statt der Vergangenheitsform, falls [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) es verlangt (offene Frage an den Projektinhaber) |
