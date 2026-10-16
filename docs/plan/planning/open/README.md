# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-v1-sandbox-docker-nachholmessung-desktop`](slice-v1-sandbox-docker-nachholmessung-desktop.md) | Docker-Desktop-Messung für `nestedRuntime: docker` (Abnahmekriterium 6) |
| [`slice-v2-homebrew-smoke-haertung`](slice-v2-homebrew-smoke-haertung.md) | Härtung des Homebrew-Smoke im Tap (Build-from-source, `u-boot add`, gatendes Audit), Formel-`desc` und Workflow-Kommentar |
| [`slice-v1-cli-verbosity-letzte-option-gewinnt`](slice-v1-cli-verbosity-letzte-option-gewinnt.md) | Verbosity: die zuletzt angegebene Option gewinnt (Code folgt [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), statt fester Rangfolge) |
| [`slice-cr-cli-quiet-fehler-und-warnungen`](../done/slice-cr-cli-quiet-fehler-und-warnungen.md) | CR: `--quiet` bleibt nur Fehler, Ausnahme für Sicherheitslockerungen (Lastenheft 0.8.0) |
| [`slice-cr-cli-quiet-weitere-ausnahmen`](slice-cr-cli-quiet-weitere-ausnahmen.md) | CR: `--quiet` bleibt nur Fehler, Ausnahmen für Degradations-Fallbacks und Compose-Fortschritt (Lastenheft 0.9.0) |
| [`slice-v1-cli-quiet-nur-fehler`](slice-v1-cli-quiet-nur-fehler.md) | Verbosity: `--quiet` zeigt nur Fehler in allen Subkommandos, drei Ausnahmen bleiben; blockiert bis Lastenheft 0.9.0 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)) |
| [`slice-harness-ai-harness-init-adoption`](../done/slice-harness-ai-harness-init-adoption.md) | ai-harness-init als Werkzeug des Repos: Aggregator, Fragmente, `repo.mk`, Sensor-Dateien, d-check v0.86.1 |
| [`slice-cr-lastenheft-produktanforderungen`](../done/slice-cr-lastenheft-produktanforderungen.md) | CR: Lastenheft nur mit Produktanforderungen — BUILD-, PROJDOCS-, ARCH- und QA-Anforderungen als Grabsteine |
| [`slice-harness-wellen-betrieb`](../done/slice-harness-wellen-betrieb.md) | Wellen-Betrieb nach Baseline: Welle-Dateien, Archivierung, Beobachtungs-Register (MR-003-Entscheidung) |
