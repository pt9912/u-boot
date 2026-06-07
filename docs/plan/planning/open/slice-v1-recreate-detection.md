# Slice V1: `u-boot up` Recreate-Warnings-Detection

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Recreate-Detection-Carveout aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> §Out of Scope T0-(k). Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

`u-boot up` (LH-FA-UP-001..003) startet die Compose-Environment
und liefert per `--json` einen Status-Envelope mit
`data.services[]`. Heute fehlt jede Warning bei
**Container-Recreate** (Compose-Plan-Drift gegen laufenden
Stack):

- Wenn `image:`-Tag wechselt (z. B. `postgres:16-alpine` →
  `postgres:17-alpine`) erstellt Compose den Container neu —
  Datenverlust bei nicht-persistierten Volumes.
- Wenn `environment:`-Variablen sich ändern, kann Compose ein
  Recreate triggern (je nach Variable).
- Wenn `volumes:`-Mounts sich ändern, je nach Compose-Logik
  auch Recreate.

`u-boot up` zeigt heute keine WARN dafür. JSON-Konsument
(z. B. CI-Skript) bekommt einen "everything stabilized"-
Envelope obwohl die Compose-Operation destruktiv war.

`driving.WarningEntry`-Type ist **proaktiv** aus dem remove-T2
Cluster-Vorlauf (R12-LOW-F4) verfügbar — der `UpResponse.
Warnings []WarningEntry`-Field existiert seit up-down T2.
Heute leer auf dem Happy-Path; konkrete Detection wartet.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Beschwerde** über fehlende Replace-Warnings
  (z. B. CI-Bug "tests passed but `up` silently recreated the
  postgres container, data is gone").
- **Cluster-T_close-Audit** (slice-v1-cli-json-dry-run nach 9/9
  Folge-Slices) fordert vollständige WARN-Coverage als
  Vertrags-Schuld.
- **Compose-Plan-Pre-Walk-Refactor**: ein anderer Slice
  führt `docker compose config`-Parse für andere Zwecke ein
  (z. B. Validation-Pre-Check vor `compose up`), dann lohnt
  sich die Recreate-Detection als zusätzlicher Konsument
  desselben Pre-Walks.

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen vor der Implementation:

1. **Docker-Adapter-Vertrag**: `driven.DockerEngine` braucht
   eine `ComposeConfig(ctx, baseDir)`-Methode die die effektive
   Compose-Config zurückliefert (analog `docker compose
   config --format json`). Plus eine Methode für aktuelle
   Container-Hashes (oder Image-Digests) aus `compose ps`.
2. **Recreate-Detection-Algorithmus**: pre-`ComposeUp` ein
   Snapshot der laufenden Container ziehen, Compose-Plan
   parsen, Diff bilden, WARN für jedes Service mit Image-
   Hash-Wechsel ODER state-impacting `environment:`/`volumes:`-
   Diff emittieren.
3. **WARN-Form**: `WarningEntry{Code: "LH-FA-UP-???", Level:
   "warn", Message: "container 'postgres' will be replaced
   (image-digest changed)", Subject: "postgres"}`. Subject
   ist proaktiv im WarningEntry-Type aus remove T2 R12-LOW-F4.
   LH-Code-Anker noch offen (eigene Spec-Erweiterung oder
   Subsumtion unter `LH-FA-UP-003`).

## Out of Scope

- **Recreate-Plan-Dry-Run** (`u-boot up --plan-only` o. ä.):
  separater Subcommand-Slice, nicht Teil dieser Erweiterung.
- **Backup-vor-Recreate**: ein zukünftiger `up --backup`-Flag
  der Volumes vor Recreate snapshotted — eigener Slice falls
  Real-World-Druck.

## Spec-Bezug

- `LH-FA-UP-003` — Startstatus anzeigen (Erweiterungs-Kandidat
  für WARN-Anker).
- `LH-NFA-USE-004` §1813 — Minimalkontrakt-Vertrag (WARN-
  Diagnostic-Form ist konsistent).
