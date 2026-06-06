# Slice V2: `generate devcontainer` Rollback-aware Multi-File-Write

> **Status:** `open/`, on hold pending trigger. Cleanup-/Hardening-
> Slice zum Devcontainer-Phase-2-Half-Write-Carveout aus
> [`slice-v1-cli-json-dry-run-generate`](slice-v1-cli-json-dry-run-generate.md)
> T0-(i). Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

`generate devcontainer` schreibt zwei Files (`.devcontainer/
devcontainer.json` + `.devcontainer/Dockerfile`) in einer **Two-
Phase-Architektur**:

- **Phase 1** (`planDevcontainerFiles`, `generate.go:618-624`):
  klassifiziert beide Files ohne FS-Mutation. Wenn auch nur ein
  File `present-no-block` oder `malformed` ist, returnt der
  Use-Case `ErrGenerateManualConflict` **ohne ein einziges
  WriteFile**. Phase 1 ist Pre-Write-Validation-atomar.
- **Phase 2** (`executeDevcontainerPlans`): schreibt nacheinander.
  Wenn File 2 mid-stream failt (Disk-Full, Permission, Race),
  ist File 1 bereits committed; das `.devcontainer/`-Verzeichnis
  bleibt in **halbgeschriebenem Zustand** auf Disk.

Plus: `applyAllowExternalFeatureSources` (`generate.go:670`) als
LAST-Schreib-Operation mutiert `u-boot.yaml` nach den
devcontainer-Files. Auch hier kann ein Mid-Failure einen
halbgeschriebenen Zustand zwischen drei Artefakten produzieren.

V1-Closure des generate-Slice akzeptiert diesen Half-State als
**bewussten Carveout** — der V1-Recorder ist nicht Roll-back-
aware (Cluster-T0-(b) Variante 3 ChangeSet-Pattern verworfen
für V1, weil Add/Init keinen Roll-back-Bedarf haben und ein
Cluster-übergreifender Pattern-Bruch nicht V1-würdig war).

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Beschwerde** über halbgeschriebenen
  `.devcontainer/`-Zustand (z. B. CI bricht, User berichtet
  „File 1 ist neu, File 2 fehlt komplett, wie repariere ich das").
- **Devcontainer-Schema-Erweiterung**: ein zukünftiger Slice
  fügt einen dritten oder vierten devcontainer-File (z. B.
  Dockerfile.dev, post-create-script) hinzu — Half-State-Risk
  skaliert mit File-Anzahl.
- **Cluster-T_close-Cleanup-Audit**: bei systematischem
  Carveout-Inventur kann die Half-State-Akzeptanz als
  Vertrags-Schuld neu bewertet werden.

## Lösungs-Skizze (vorläufig)

Drei Optionen mit unterschiedlicher Tiefe:

1. **Pre-Phase-2-Snapshot + Rollback-on-Failure**:
   `executeDevcontainerPlans` snapshotted alle existing Files
   in Buffer (oder als `.bak.<n>`), schreibt File 1, schreibt
   File 2 (+ ggf. `u-boot.yaml`-Mutation); bei Failure → ALLE
   commiteten Files aus Snapshot zurückschreiben oder löschen.
   Lokaler Fix ohne Recorder-Architektur-Eingriff. **Echte
   Multi-File-Atomicity** (alle-oder-keiner). Risiko:
   zweiter Failure beim Rollback-Write hinterlässt
   inkonsistenten State — Best-Effort-Rollback mit
   ErrGenerateFileSystem-Wrap, der den Rollback-Failure-State
   explizit signalisiert.
2. **Per-Use-Case Roll-back-aware Recorder**: ChangeSet-Pattern
   (Cluster-T0-(b) Variante 3) speziell für `generate devcontainer`.
   Schmaler als Cluster-weit, weil Init/Add weiterhin
   capture-only sind. Semantisch sauberste Lösung, aber Architektur-
   Eingriff (RecorderPort-Interface erweitern).
3. **~~Per-File Temp+Rename~~ — verworfen**: Files in
   `<file>.tmp.<n>` schreiben, dann je `os.Rename`. **Löst das
   Multi-File-Problem NICHT**: wenn Rename 1 succeeds und
   Rename 2 failt, bleibt File 1 committed, File 2 in tmp-
   Zustand — exakt der Half-State, den V2 vermeiden will.
   Per-File-Atomicity ≠ Multi-File-Atomicity. (R3-Finding
   gegen die ursprüngliche Stub-Empfehlung.)

**Bevorzugte Skizze**: **Option 1 (Snapshot + Rollback-on-
Failure)** — minimaler Architektur-Eingriff bei echter Multi-
File-Atomicity. Best-Effort-Rollback ist ehrlicher als
falsche per-File-Atomicity-Versprechen. Trigger-Slice
klärt: Snapshot-Form (In-Memory-Buffer vs. `.bak.<n>`-Files),
Rollback-Sequenz (LIFO über commit-Liste), Rollback-Failure-
Signalisierung (eigener Sentinel oder gewrappter
ErrGenerateFileSystem mit Hinweis-Message).

Failure-Injection-Pin im Trigger-Slice: „erste Datei committed,
zweite Rename/YAML-Write failt → Restore aktiviert, Disk-Zustand
nach Aufruf == Disk-Zustand vor Aufruf (oder explizite Roll-back-
Failure-Diagnostic)."

**Rollback-Scope (R4-Finding 4)**: „Disk-Zustand vor Aufruf"
umfasst MEHR als nur die zwei Devcontainer-Files. Phase 2
schreibt mindestens drei Side-Effects mit Cleanup-Pflicht:

1. **`.devcontainer/`-Verzeichnis selbst**: wird im Pre-Phase-2-
   Schritt via `MkdirAll` erzeugt (`generate.go:848`). Existierte
   das Verzeichnis vor dem Aufruf NICHT (Fresh-Project) und der
   gesamte Aufruf failt → Rollback muss das leere
   `.devcontainer/`-Dir wieder entfernen, sonst bleibt ein
   Scratch-Artefakt zurück (Tree-Diff zeigt extra Dir).
2. **`u-boot.yaml`-Allowlist-Mutation** (`generate.go:951`,
   `applyAllowExternalFeatureSources`): wird LAST geschrieben.
   Bei Failure NACH den zwei devcontainer-Files +
   YAML-Mutation müssen ALLE drei zurückgesetzt werden.
3. **Snapshot-Persistierung**: bei In-Memory-Buffer-Variante
   harmlos; bei `.bak.<n>`-File-Variante muss der Cleanup-Pfad
   die `.bak.<n>`-Files am Ende der Erfolgs-Sequenz löschen,
   sonst bleiben sie als Scratch-Artefakte.

**Echter Acceptance-Pin (Trigger-Slice T6)**: Failure-Injection
auf einem **frischen Projekt ohne `.devcontainer/`-Dir** und
ohne `feature-sources`-Block in `u-boot.yaml`. Pin:
- vor dem Aufruf: `tree` zeigt nur die u-boot.yaml + restliche
  Projekt-Files
- generate devcontainer mit `--allow-external-feature-sources
  https://...` (triggert alle drei Side-Effects)
- WriteFile-Injection-Spy lässt File 2 (`Dockerfile`) failen
- nach dem Aufruf: `tree` IDENTISCH zum Pre-State (kein leeres
  `.devcontainer/`-Dir, keine `.bak.<n>`-Files, u-boot.yaml
  byte-identisch zum Pre-State).

## Out of Scope (V1)

- **Cluster-weiter ChangeSet-Pattern-Recorder** — bewusst V1-
  out-of-scope, weil Add/Init kein Roll-back-Bedarf haben und
  ein gemeinsamer Pattern-Bruch unnötig wäre.
- **Backup/Snapshot-Persistierung über Use-Case-Grenzen
  hinaus** — Crash-Recovery zwischen `generate`-Aufrufen ist
  V3-Scope.

## Bezug

- Carveouts-Tracking:
  [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
  §Temporäre Carveouts — Eintrag verweist hierher.
- Generate-Slice T0-(i):
  [`slice-v1-cli-json-dry-run-generate`](slice-v1-cli-json-dry-run-generate.md)
  §Sub-Decisions T0-(i) Devcontainer-Atomicity-Klärung.
- Code-Anker:
  [`generate.go:618-690`](../../../../internal/hexagon/application/generate.go)
  (Phase-1-Comment + Phase-2-Implementation).
- Spec: `LH-FA-DEV-001` (Devcontainer-Render),
  `LH-NFA-REL-003` (technische Persistenz-Klasse).
- Phase: V2 (Hardening, post-V1-Cluster-Closure).
