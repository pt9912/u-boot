# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-v2-homebrew-smoke-haertung`](../in-progress/slice-v2-homebrew-smoke-haertung.md) | Härtung des Homebrew-Smoke im Tap (Build-from-source, `u-boot add`, gatendes Audit), Formel-`desc` und Workflow-Kommentar |
| [`slice-v1-sandbox-fallback-ausweisung`](slice-v1-sandbox-fallback-ausweisung.md) | Degradations-Fallbacks Daemon und Egress host-seitig ausweisen (Code) oder [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) präzisieren (CR); Entscheidung beim Projektinhaber |
| [`slice-v1-cli-dry-run-would-formulierung`](slice-v1-cli-dry-run-would-formulierung.md) | Dry-Run: `remove`, `generate` und `config set` sagen „Would …“ statt der Vergangenheitsform, falls [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) es verlangt (offene Frage an den Projektinhaber) |
