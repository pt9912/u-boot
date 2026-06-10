# Slice V1: Strukturierte Multi-Port-Liste für `u-boot up --json`

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum Multi-Port-Format-Carveout aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> §Out of Scope T0-(g). Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts.

## Auslöser

`domain.ServiceStatus.Port` (`domain/serviceup.go:170-174`) ist
heute ein **Komma-getrennter Display-String** (`"5432:5432,
127.0.0.1:9091:9091"`). Multi-Port-Reporting **existiert**
schon; nur das Format ist nicht strukturiert.

`u-boot up --json` liefert pro Service:

```json
{"name": "postgres", "state": "running",
 "port": "5432:5432, 127.0.0.1:9091:9091",
 "healthcheck": "healthy"}
```

JSON-Konsument muss CSV-Parsing machen. Für
Per-Port-Health-Reporting oder Per-Port-Filter in CI-Scripts
wäre ein strukturiertes Array natürlicher:

```json
{"name": "postgres", "state": "running",
 "ports": ["5432:5432", "127.0.0.1:9091:9091"],
 "healthcheck": "healthy"}
```

**Funktional gleichwertig**, nur Konsumenten-Parse-Last
reduziert.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Konsumenten-Druck** nach strukturierter Form
  (z. B. CI-Skript-Autor beschwert sich über CSV-Parsing-
  Edge-Cases).
- **Per-Port-Erweiterung** (z. B. Per-Port-Probe-Results im
  `doctor`-Slice): braucht ohnehin strukturierte Form.
- **Cluster-T_close-Audit** fordert Konsistenz mit
  doctor-Status-Format das Per-Port-Strukturen tragen kann.

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen vor der Implementation:

1. **Domain-Refactor**: `domain.ServiceStatus.Port string` →
   `Ports []string`. Bricht den heutigen Vertrag —
   alle Konsumenten (`cli/statusview.go`, `cli/up.go`,
   tests) müssen migriert werden. Plus der Compose-Parser im
   Application-Layer der heute den CSV-String aus
   `compose ps`-Output baut.
2. **CLI-Wire-Form-Migration**: `serviceStatus.Port string`
   in `cli/up.go` → `serviceStatus.Ports []string`. Empty-
   Array-Pin (T0-(j) Pattern aus up-down): nil-Slice MUST
   serialize as `[]`. JSON-Konsument-Backward-Compatibility
   ist nicht möglich (Schlüssel-Rename `port` → `ports`)
   außer mit Dual-Field-Übergangsphase (`port` als deprecated
   `omitempty`).
3. **Human-Mode-Status-Tabelle**: `renderUpStatus`
   (`cli/statusview.go:24`) heute `port string` → muss
   `ports []string` joinen. Display-Pattern: CSV bleibt
   (`strings.Join(ports, ", ")` mit `"-"` für leer).

## Out of Scope

- **Port-Mapping-Sub-Klassifikation** (`{public: 5432,
  internal: 5432, protocol: "tcp"}`): wäre eine **dritte**
  Form (Struct-Array statt String-Array). Eigener Slice
  falls Real-World-Druck (z. B. Per-Port-Protocol-Awareness
  für probing).
- **Port-Range-Form** (`5432-5435:5432-5435`): Compose
  unterstützt Ranges; heutige CSV-Form trägt sie als
  Single-String. Strukturierte Form müsste entscheiden ob
  Ranges aufgespalten oder als atomarer String belassen
  werden.

## Spec-Bezug

- [`LH-FA-UP-003`](../../../../spec/lastenheft.md#lh-fa-up-003-startstatus-anzeigen) — Mindestangabe "Port" (Singular im Spec-
  Wortlaut; strukturierte Liste ist eine Erweiterung, keine
  Spec-Verletzung).
- [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) §1813 — JSON-Konsumenten-Vertrag.
