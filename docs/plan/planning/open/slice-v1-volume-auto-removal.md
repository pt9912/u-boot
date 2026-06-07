# Slice V1: `u-boot remove --purge` Volume-Auto-Removal

> **Status:** `open/`, on hold pending trigger. Cleanup-/Hardening-
> Slice zum Volume-Auto-Removal-Carveout aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> §Out of Scope. Carveout-Plan-Anker ([[feedback_carveouts_need_plans]]);
> verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

`u-boot remove <service> --purge` ist in v0.3.0/v0.4.0 als
**Confirmation-Gate-Plus-Status-Carrier** implementiert, aber
**ohne** echte `docker volume rm`-Calls:

- `RemoveServiceService.runPurgeGate`
  (`internal/hexagon/application/removeservice.go:158-178`) prompted
  via `driven.Confirmer`, returnt aber bei `confirmed=true` nur
  `nil` — kein FS- oder Docker-Side-Effect.
- `RemoveServiceResponse.VolumesPurged` ist hartcodiert `false`
  (T0-(h) Sub-Decision: deferred-Status, kein Code-Pfad setzt es
  jemals auf `true`).
- Human-Mode (`cli/remove.go:526-547` `printRemoveSummary`) druckt
  eine WARNING auf stderr inkl. manuellem
  `docker volume ls --filter label=…` + `docker volume rm <name>`-
  Cleanup-Hint.
- JSON-Mode emittiert die WARNING als `diagnostics[]`-Eintrag mit
  `code: "LH-FA-ADD-007"`, `level: "warn"` (T0-(g) WARN-Migration).

Spec §2602 (`LH-FA-ADD-007` Volume-Anforderung) ist damit
**partiell** erfüllt: `--purge` ist als opt-in CLI-Surface
vorhanden, der Confirmation-Gate ist konsistent mit `down
--volumes` (`LH-FA-CLI-005A` §254), aber die tatsächliche
Volume-Removal ist auf den User abgewälzt.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Beschwerde** über manuellen Cleanup-Schritt
  (z. B. User berichtet "ich hatte 12 stale postgres-data-Volumes
  bevor ich gemerkt habe dass `--purge` nichts entfernt").
- **Cluster-T_close-Audit** (slice-v1-cli-json-dry-run nach 9/9
  Folge-Slices) fordert vollständige `--purge`-Semantik als
  Vertrags-Schuld.
- **Docker-Compose-V3-Migration** oder ähnliche Compose-Spec-
  Erweiterung, die Volume-Lifecycle anders modelliert und einen
  Re-Audit triggert.

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen, vor der eigentlichen Implementierung
zu klären:

1. **Docker-Adapter-Verträge erweitern**: `driven.DockerClient`
   (oder ein neuer `driven.VolumeManager`-Port?) bekommt eine
   `RemoveVolumes(ctx, projectLabel, volumeNames)`-Methode. Implementer
   im Adapter ruft `docker volume rm <name>` oder das Compose-CLI-
   Äquivalent. Sub-Decision: ein dedizierter Port (saubere
   Schicht-Trennung) vs. Methode am existierenden `DockerClient`-
   Port (weniger Boilerplate).
2. **Catalog-Volume-Discovery**: `serviceCatalogueEntry.volumeRefLiteral`
   (`addservice_execute.go:190-224`) trägt heute den compose-
   internen Volume-Namen (`postgres-data` für postgres). Real-World-
   Naming hat ein `<project>_`-Prefix (Compose-Project-Label). Sub-
   Decision: `docker volume ls --filter label=com.docker.compose.
   project=<project>` zur Runtime ODER deterministische
   Projekt-Name-Konstruktion aus `u-boot.yaml`.
3. **Partial-Removal-Atomicity**: was passiert wenn Volume 1
   entfernt wird aber Volume 2 failt (Docker-Daemon-Race,
   Volume-in-use)? Analog
   [`slice-v2-generate-devcontainer-rollback-aware-write`](slice-v2-generate-devcontainer-rollback-aware-write.md)
   Half-Write-State: Recorder-Architektur kennt keine Roll-back-
   aware-Captures für Docker-Side-Effects (Cluster-T0-(b)
   Variante 3 verworfen). Sub-Decision: Best-Effort mit per-
   Volume-Diagnostic-Eintrag im Envelope, ODER pre-flight
   Volume-in-use-Check vor jedem Remove.

## Out of Scope

- **Backup-vor-Removal**: ein optionales `--purge --backup`-Flag,
  das die Volume-Inhalte vor Removal in `<project>/.u-boot-volume-
  backup-<timestamp>/<volume-name>.tar.gz` archiviert. Separate
  Sub-Decision, eigener Folge-Slice.
- **Multi-Service-Bulk-Purge**: `u-boot remove --purge --all`
  (alle deaktivierten Services purgen). Bricht die Single-Service-
  Semantik von `remove`; eigener Slice.
- **Rollback-aware-Recorder**-Erweiterung für Docker-Side-Effects
  (siehe Sub-Decision 3 oben).

## Spec-Bezug

- `LH-FA-ADD-007` §2602 — Volume-Anforderung (`--purge`-Opt-in-
  Form ist erfüllt; tatsächliche Removal ist offen).
- `LH-FA-CLI-005A` §254 — Confirmation-Gate für destruktive
  Operationen (heute bereits implementiert).
- `LH-NFA-REL-003` — FS-Failure-Klasse erbt auf Docker-Failure-
  Klassifikation; ggf. neuer `LH-NFA-REL-004`-Sentinel für
  Docker-Volume-Errors.
