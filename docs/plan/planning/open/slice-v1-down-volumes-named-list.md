# Slice V1: `u-boot down --volumes` Named-Volume-Liste

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Volume-Named-Liste-Carveout aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> §Out of Scope T0-(h). Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

`u-boot down --volumes` entfernt Named-Volumes via
`docker compose down -v` (LH-FA-UP-004 §1015 destructive op).
Heutiger Port-Vertrag
([`port/driving/down.go:80`](../../../../internal/hexagon/port/driving/down.go))
trägt `DownResponse.RemovedVolumes bool` als boolean Echo —
NICHT eine Liste der tatsächlich entfernten Volume-Namen.

Der Port-Kommentar verbietet das **bewusst**: *"No stop /
removed counters — `docker compose down` emits a human-
readable progress stream rather than a structured count, and
inventing an 'unknown' sentinel value would force every caller
to special-case it. If a future slice needs precise counts
(e.g. for `--json` output, LH-NFA-USE-004 V1), it would add a
`ComposePs` diff before/after the call rather than parse the
stderr stream."*

JSON-Konsument bekommt heute `data.removedVolumes: true/false`.
Für Audit-Logs / CI-Cleanup-Scripts wäre eine konkrete
Namen-Liste informativer.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Konsumenten-Bedarf** nach Volume-Namen-Liste
  (z. B. Audit-Log-Anforderung: "welche genau entfernt?",
  CI-Cleanup-Script: "removed volumes für post-`down` shell
  cleanup").
- **Cluster-T_close-Audit** fordert vollständige Volume-Lifecycle-
  Berichterstattung.

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen vor der Implementation:

1. **Docker-Adapter-Vertrag**: zwei Optionen für die
   `ComposeDown --volumes`-Effekt-Bestimmung:
   - (a) **ComposePs-Diff vor/nach**: ein
     `DockerEngine.ListVolumes(ctx, baseDir)`-Snapshot vor
     `ComposeDown`, ein zweiter danach, Difference =
     entfernte Volumen. Pattern-Vorbild Port-Kommentar.
   - (b) **Docker-Volume-API direkt**: `DockerEngine.
     ListVolumesByLabel(ctx, "com.docker.compose.project=<name>")`
     vor `ComposeDown` + Snapshot nach Down. Cleaner aber
     braucht Compose-Projekt-Label-Awareness.
2. **`downStatusData`-Carrier-Form-Migration**: heute
   `{removedVolumes bool}` → neu `{removedVolumes []string,
   removedVolumesEcho bool}` ODER pure `{removedVolumes
   []string}` mit Empty-Liste statt false-Echo. Kompatibilität
   für bestehende JSON-Konsumenten beachten (eventuell Empty-
   Array-Pin).
3. **Partial-Failure-Semantik**: was wenn ComposeDown
   Volume A entfernt, Volume B failed (in-use)? Snapshot-Diff
   würde Volume A als entfernt zeigen, Volume B als noch
   vorhanden. Error-Pfad-Form: `data.removedVolumes: [A]` +
   `diagnostics[]` mit B-Failure-Eintrag.

## Out of Scope

- **Backup-vor-Removal**: ein `--volumes --backup`-Flag der
  Volume-Inhalte vor Removal in `<project>/.u-boot-volume-
  backup-<timestamp>/<volume-name>.tar.gz` archiviert.
  Separater Slice falls Real-World-Druck (analog dem remove-
  Slice `slice-v1-volume-auto-removal` Out-of-Scope).
- **Volume-Lifecycle-Reporting** für `up` (welche Volumes
  wurden erstellt): wäre eigener Slice; up nutzt
  `RemovedVolumes` als Bool nicht.

## Spec-Bezug

- `LH-FA-UP-004` §1015 (Volume-Removal-Destructive-Op).
- `LH-NFA-USE-004` §1813 (JSON-Konsumenten-Vertrag).
