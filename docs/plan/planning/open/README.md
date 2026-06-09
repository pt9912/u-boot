# docs/plan/planning/open

Backlog: Slice- und Tranchen-Pläne, die noch nicht für die nächste
Iteration eingeplant sind.

Übergang nach `next/` per `git mv`, sobald ein Artefakt priorisiert
wird. Siehe `docs/plan/planning/README.md` für Lifecycle und
Dateiname-Konventionen.

## Bestand

| Datei | Gegenstand |
| ----- | ---------- |
| [`slice-v1-release-cut-v0.5.0.md`](slice-v1-release-cut-v0.5.0.md) | **Ready-to-execute:** Release-Cut v0.5.0 (Ein-Feature-Minor: local-templates) — T1–T3 morgen, T4 Tag-Push |
| [`slice-v1-cli-cleanup-add-preview-mode-alias.md`](slice-v1-cli-cleanup-add-preview-mode-alias.md) | Cleanup: `AddPreviewMode`-Alias entfernen |
| [`slice-v1-config-list-subcommand.md`](slice-v1-config-list-subcommand.md) | Cleanup: `u-boot config list` als eigener Subcommand mit strukturiertem Path-Value-Tree |
| [`slice-v1-config-multi-path-get.md`](slice-v1-config-multi-path-get.md) | Cleanup: `u-boot config get` Multi-Pfad-Get mit `--json-array` |
| [`slice-v1-config-multi-path-set.md`](slice-v1-config-multi-path-set.md) | Cleanup: `u-boot config set` Multi-Path-Set (atomar mehrere Pfade) |
| [`slice-v1-config-structured-hint.md`](slice-v1-config-structured-hint.md) | Cleanup: `config` strukturiertes `data.hint{action, argument}`-Field |
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
