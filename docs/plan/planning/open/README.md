# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-v1-devcontainer-feature-sshd`](slice-v1-devcontainer-feature-sshd.md) | Katalog-Schlüssel `sshd` (offizielles Feature), Egress, Sandbox-Regel, Tests, Doku; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), dritter Slice |
| [`slice-v1-devcontainer-feature-claude-code`](slice-v1-devcontainer-feature-claude-code.md) | Katalog-Schlüssel `claude-code` (Feature von Anthropic, Drittanbieter), Egress, Zugangsdaten zur Laufzeit, Tests, Doku; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), vierter Slice |
| [`slice-v2-spec-devcontainer-sshd-claude-code`](slice-v2-spec-devcontainer-sshd-claude-code.md) | Spezifikation zu Katalog, Egress und Sandbox-Verhalten beider Features; Welle [`welle-devcontainer-sshd-claude-code`](../welle-devcontainer-sshd-claude-code.md), letzter Slice |
| [`slice-v1-sandbox-fallback-ausweisung`](slice-v1-sandbox-fallback-ausweisung.md) | Degradations-Fallbacks Daemon und Egress host-seitig ausweisen (Code) oder [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) präzisieren (CR); Entscheidung beim Projektinhaber |
| [`slice-v1-cli-dry-run-would-formulierung`](slice-v1-cli-dry-run-would-formulierung.md) | Dry-Run: `remove`, `generate` und `config set` sagen „Would …“ statt der Vergangenheitsform, falls [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) es verlangt (offene Frage an den Projektinhaber) |
