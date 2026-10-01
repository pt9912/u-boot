# Slice V1: `u-boot up` Partial-Snapshot bei Mid-`ComposeUp`-Failure

> **Status:** **abgeschlossen** (2026-10-01, **Delivery-Hash: `c4c9b09`**). Cleanup-/Feature-
> Slice zum Carveout aus
> [`slice-v1-cli-json-dry-run-up-down`](../done/slice-v1-cli-json-dry-run-up-down.md)
> §Out of Scope (Carveout-Eintrag entfernt).

## Auslöser

`UpService.Up` returnt bei Mid-Failure-Szenarien eine
**Zero-Response** + Error:

- `ComposeUp`-Fehler (`upservice.go:76-80`): `driving.UpResponse
  {}` + `fmt.Errorf("up service: ComposeUp on %q: %w", ...)`.
- Poll-Failures (`upservice.go:200-202`): `driving.UpResponse{}` +
  `fmt.Errorf("up service: poll cancelled at t=...: %w", ...)`
  bzw. `ComposePs at t=...`.
- Stabilization-Timeout (`upservice.go:208-211`): Zero-Response +
  `ErrStabilizationTimeout`.
- Terminal-State-Failure (`upservice.go:197-200`): Zero-Response +
  `ErrComposeRuntime`-Wrap.

JSON-Konsument bekommt `data: null` mit Error-Diagnostic.
Für Mid-Failure-Debugging (z. B. CI-Diagnose "was lief schon
hoch, was nicht?") wäre ein **Partial-Snapshot** der teilweise
gestarteten Services hilfreich.

Heutige Architektur erlaubt das **nicht**:

- `domain.ContainerState`-Enum (`domain/serviceup.go:20-40`)
  kennt nur `unknown|starting|running|restarting|dead` — kein
  `failed`.
- `UpResponse`-Struct hat keinen `PartialServices`-Field.
- Use-Case-Code-Pfad bricht früh mit Zero-Response ab; kein
  Pre-Error-`ComposePs`-Snapshot.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Bedarf** nach Mid-Failure-Debugging-Daten
  (z. B. interaktive CI-Diagnose, Postmortem-Logs).
- **Domain-Enum-Erweiterungs-Slice** der `StateFailed` cluster-
  weit etabliert (z. B. für `doctor`-Erweiterung): dann lohnt
  sich der Partial-Snapshot als sekundärer Konsument.
- **Compose-`config`-Pre-Walk-Slice** (siehe
  [`slice-v1-recreate-detection`](../done/slice-v1-recreate-detection.md)): hat ohnehin `ComposePs`-
  Snapshot-Infrastruktur.

## Lösungs-Skizze (vorläufig)

Drei Sub-Entscheidungen vor der Implementation:

1. **`UpResponse`-Vertrags-Erweiterung**: neuer Field
   `PartialServices []domain.ServiceStatus` mit Doc-
   Convention "populated when error is non-nil; happy-path
   leaves it empty". Pattern-Erbe `Warnings`-Field aus T2
   (auch optional).
2. **Use-Case-Refactor-Stellen**: vor jeder `return driving.
   UpResponse{}, fmt.Errorf(...)` einen `s.engine.ComposePs(
   ctx, baseDir)`-Snapshot ziehen, in `PartialServices`
   verpacken, dann mit Error returnen. Drei Stellen:
   `upservice.go:80`, `:197-202`, `:208-211`.
3. **Domain-Enum-Erweiterung**: `domain.ContainerState`
   bekommt `StateFailed` mit Doc-Convention "terminal failure
   from Compose-side (exited with non-zero code) — distinct
   from `StateDead` which is Docker-API-side terminal".
   Alle Switch-Statements im Code (cli/statusview.go +
   application/upservice.go classify-Logik) müssen die neue
   State migriert behandeln (Pattern-Erbe enum-Add-Disziplin).

## Out of Scope

- **`failed-on`-Field** im `serviceStatus` (CLI-Layer-Carrier
  type, `cli/up.go`): pro Service ein "failed at step X" —
  wäre zusätzliche Sub-Klassifikation, eigener Slice falls
  Real-World-Druck.
- **Retry-Recovery-Hints** (`hint: "try docker compose logs <svc>"`):
  Hint-Generierung wäre Konsumenten-UX-Erweiterung, separat.

## Spec-Bezug

- [`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten) §966-§969 (Stabilisierungs-Semantik).
- [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003--abbruch-bei-kritischen-fehlern) (FS-Failure-Klasse, indirekt für
  Mid-Failure-Reporting).

## Closure-Notiz

**Geliefert:** `UpResponse.PartialServices` (ComposeUp-Fehler außer Daemon-nicht-erreichbar, Terminal-Zustand, Stabilisierungs-Timeout); der JSON-Fehler-Envelope trägt `data.services`. Abweichung: kein neues `StateFailed` — `exited` wird bereits als `dead` klassifiziert und reicht für die Diagnose.

**Sensoren:** `make gates` grün; Docker-Integrationstests gegen echtes Docker/Compose
(`TestE2E_UpDownCluster_NamesRecreateAndDown`, `TestE2E_RemovePurgeRemovesVolume`,
`TestE2E_RemovePurgeRemovesVolume_StackDown` — das Volume ist nach `--purge` verschwunden);
Unit-/CLI-/Adapter-Tests (Parser, Snapshot, Fehlerpfade).

**Doku:** `docs/user/cli-json-output.md` §6.7, Benutzerhandbuch §4.4/§4.5, CHANGELOG.
