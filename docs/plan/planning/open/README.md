# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-v1-cli-cleanup-add-preview-mode-alias.md`](slice-v1-cli-cleanup-add-preview-mode-alias.md) | Cleanup: `AddPreviewMode`-Alias entfernen |
| [`slice-v1-cli-json-dry-run-template.md`](slice-v1-cli-json-dry-run-template.md) | Folge-Slice 9/9 des Cluster-Slice `slice-v1-cli-json-dry-run`: `template list --json` Envelope-Migration |
| [`slice-v1-cli-json-envelope-consolidation.md`](slice-v1-cli-json-envelope-consolidation.md) | Konsolidierung: CLI-JSON-Envelope-Pattern add/init/generate (R15-Cross-Slice-1) |
| [`slice-v1-down-volumes-named-list.md`](slice-v1-down-volumes-named-list.md) | Cleanup: `u-boot down --volumes` Named-Volume-Liste |
| [`slice-v1-keycloak-ci-flake.md`](slice-v1-keycloak-ci-flake.md) | Hardening: Keycloak-Acceptance-Test (`LH-AK-003`) in CI grün |
| [`slice-v1-logs-format-flags.md`](slice-v1-logs-format-flags.md) | Cleanup: `u-boot logs --no-log-prefix`/`--timestamps` Format-Flags |
| [`slice-v1-logs-multi-service-filter.md`](slice-v1-logs-multi-service-filter.md) | Cleanup: `u-boot logs <svc1> <svc2>` Multi-Service-Filter |
| [`slice-v1-logs-time-range-filter.md`](slice-v1-logs-time-range-filter.md) | Cleanup: `u-boot logs --since`/`--until` Time-Range-Filter |
| [`slice-v1-multi-port-services.md`](slice-v1-multi-port-services.md) | Cleanup: Strukturierte Multi-Port-Liste für `u-boot up --json` |
| [`slice-v1-recreate-detection.md`](slice-v1-recreate-detection.md) | Cleanup: `u-boot up` Recreate-Warnings-Detection |
| [`slice-v1-up-partial-snapshot-on-failure.md`](slice-v1-up-partial-snapshot-on-failure.md) | Cleanup: `u-boot up` Partial-Snapshot bei Mid-`ComposeUp`-Failure |
| [`slice-v1-volume-auto-removal.md`](slice-v1-volume-auto-removal.md) | Cleanup: `u-boot remove --purge` Volume-Auto-Removal |
| [`slice-v2-distro-pakete.md`](slice-v2-distro-pakete.md) | Plan-Stub: Debian-/RPM-Pakete (`LH-OPEN-002`-Restweg) |
| [`slice-v2-generate-devcontainer-rollback-aware-write.md`](slice-v2-generate-devcontainer-rollback-aware-write.md) | Cleanup: `generate devcontainer` Rollback-aware Multi-File-Write |
| [`slice-v2-homebrew-formula.md`](slice-v2-homebrew-formula.md) | Plan-Stub: Homebrew-Formula (`LH-OPEN-002`-Restweg) |
