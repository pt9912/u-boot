# Slice V1: `u-boot remove <service>`

## Auslöser

[`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (V1) verlangt einen Mechanismus, ein Add-on
wieder zu entfernen. Heute kann man Postgres via `u-boot add
postgres` aktivieren — aber nicht sauber entfernen; nur
`u-boot config get services.postgres.enabled` zeigt den State,
ein Toggle ist nicht erlaubt (M8: services.<svc>.enabled ist
get-only, „toggling goes through `u-boot add` / `remove`"). Dieser
Slice schließt die fehlende Hälfte.

Erster Slice des v0.3.0-Milestones („Add-on Catalogue Expansion").

## Aufhebungsbedingung

`u-boot remove postgres` in einem initialisierten Projekt mit
aktivem Postgres-Add-on entfernt den Compose-Block, den
env-Block und setzt `services.postgres.enabled: false` in
`u-boot.yaml`. Idempotent: bereits-disabled liefert eine
„nothing to do"-Meldung ohne Fehler. `--purge` triggert das
[`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)-§254-Confirmation-Gate analog `down --volumes`
(actual Volume-Removal in v0.3.0 deferred, CLI surface'd den
Cleanup-Hint).

## Akzeptanzkriterien

- ✅ `u-boot remove postgres` (Service aktiv) erzeugt drei
  Removal-Actions: Compose-Block weg, env-Block weg,
  `services.postgres.enabled: false`. Volumes bleiben (kein
  Daten-Verlust ohne `--purge`).
- ✅ `u-boot remove postgres` (Service auf enabled: false) liefert
  idempotent eine Meldung („Service \"postgres\" is already
  disabled; no changes") + Exit 0; KEINE Dateien werden angefasst.
- ✅ `u-boot remove postgres` (Service nicht registriert) failt
  mit `ErrServiceUnregistered` + Exit 10 ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)).
- ✅ `u-boot remove unknown-svc` failt mit `ErrServiceUnsupported`
  + Exit 10.
- ✅ `u-boot remove postgres --purge`: Confirmation-Gate analog
  `down --volumes` — non-interactive ohne `--yes` → Exit 10
  (`ErrConfirmationRequired`); interaktiv → Confirmer-Prompt;
  Yes-Approval → executeRemove läuft; deklinte Bestätigung →
  Exit 10. Actual Volume-Löschung deferred (T3-Decision); CLI
  surface'd den `docker volume rm <name>`-Manuell-Hinweis.
- ✅ `u-boot remove postgres` ohne `u-boot.yaml` failt mit
  `ErrProjectNotInitialized` + Exit 10.
- ✅ `u-boot.yaml` nach erfolgreichem Remove zeigt
  `services.postgres.enabled: false` (NICHT entfernt) —
  konsistent zur State-Machine-Idee „registriert aber inaktiv".
- ✅ Hexagonale Verdrahtung: Driving-Port
  `RemoveServiceUseCase` + Application `RemoveServiceService`
  (spiegelt M5 `AddServiceService`-Pattern) + CLI-Subkommando.
  Wiring in `cmd/uboot/main.go`. `detectServiceState` als
  Package-Free-Function extrahiert (geteilte Nutzung mit Add).

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `ca1267f` | Driving-Port `port/driving.RemoveServiceUseCase` mit Request `{BaseDir, ServiceName, Purge, Yes, NoInteractive}` und Response `{ServiceName, PriorState, State, Changed, VolumesPurged}`. `ErrServiceUnregistered`-Sentinel neu (Exit 10) — abgegrenzt von `ErrServiceUnsupported`-Catalogue-Miss. Reuse: `ErrServiceUnsupported`, `ErrServiceInconsistent`, `ErrProjectNotInitialized`, `ErrConfirmationRequired`. Application-Skeleton `RemoveServiceService` mit Konstruktor (fs+yaml mandatory, confirmer+logger nil-tolerant), Static-Check `var _ driving.RemoveServiceUseCase = (*RemoveServiceService)(nil)`. Remove() T1-Stub: BaseDir-Empty-Check + „not yet implemented"-Error für T4-CLI-Wiring-Vorbereitung. 3 Tests. |
| T2 | `e26cb42` | Refactor: `(s *AddServiceService).detectServiceState` zu Package-Level-Free-Function `detectServiceState(fs, yaml, baseDir, svc)` extrahiert, Add+Remove nutzen sie. Volle State-Machine in `RemoveServiceService.Remove`: Catalogue-Check via `isSupportedService`, Detect-Phase liefert PriorState, Branch-Logik (Unregistered → `ErrServiceUnregistered`; Inconsistent → `ErrServiceInconsistent`; Deactivated → idempotent No-Op; Active/EnabledUnset → `executeRemove`). `executeRemove` macht drei Mutations in der Reihenfolge compose-Block → env-Block → u-boot.yaml-PatchScalar(enabled=false) — u-boot.yaml zuletzt damit Mid-Flight-Fail die enabled-Flag unangetastet lässt. `removeBlock` per `managedblock.Replace(content, marker, nil)`. 11 Tests inkl. InconsistentBlock, Malformed-Block, FS-Error-Propagation, Snapshot-basiertem No-Op-Pin. |
| T3 | `c508b4f` | `--purge`-Confirmation-Gate via neuer `runPurgeGate(ctx, req)`-Methode — exakt parallel zu `DownService.runConfirmationGate` (M6-T5), Truth-Table komplett: `!Purge` → proceed; `Purge && Yes` → proceed; `Purge && NoInteractive && !Yes` → `ErrConfirmationRequired`; `Purge && interactive` → `Confirmer.ConfirmRemoveVolumes(ctx, baseDir)` → proceed oder Refuse. Remove() umstrukturiert: Gate-Call NACH Reject-States (Unregistered, Inconsistent) und VOR proceeding-Branch — User bekommt informative Fehler statt Confirmation-Prompt für nichts. Gate feuert auch für `Deactivated + Purge` (Spec-mandatiert; Approval führt zu idempotent No-Op). **T3-Decision**: actual Volume-Removal bleibt out-of-scope; `VolumesPurged` bleibt `false`; T4-CLI surface'd den Cleanup-Hint. 6 neue Tests inkl. Snapshot-Pins „no-side-effect on refuse" + Reihenfolge-Pin „UnregisteredSkipsGate". |
| T4 | `3cc2646` | CLI-Subkommando `u-boot remove <service> [--purge]` in `internal/adapter/driving/cli/remove.go` (analog `add.go`). `removeFlags{Purge, Yes, NoInteractive}` mit Yes/NoInteractive read-through vom Root, `--purge` als lokale BoolVar. `runRemove` macht Mutex-Check (`ErrConflictingModeFlags` → Exit 2), `domain.NewServiceName`-Validation (invalid → Exit 10), Delegation an UseCase. `printRemoveSummary` mit drei Shapes (No-Op / Transition / Transition+Purge-NOTE). Wiring: `cli.New` 10. positional `removeUC`; alle 8 bestehenden `newApp*`-Test-Helper + `fakeRemoveServiceUseCase` + neuer `newAppWithRemove`-Helper. `isValidationError` refaktoriert: `ErrServiceUnsupported`/`ErrServiceInconsistent`/`ErrServiceUnregistered` in neuem `isServiceValidationError`-Helper gebündelt (gocyclo-Carve-Out parallel zu `isConfigValidationError`). `cmd/uboot/main.go`: `NewRemoveServiceService` konstruiert (selber `confirmAdapter` wie init/down), an `cli.New` durchgereicht. 8 CLI-Tests + **E2E-Smoketest** gegen das gebaute Image: `init demo --no-git` + `add postgres` (enabled: true) + `remove postgres` (enabled: false, 3 Changed-Files) + zweiter `remove` (idempotent „already disabled"). |
| T5 | `764e737` | Slice-Plan nach `done/`; README.{md,de.md} `add`-Bullet erwähnt `remove`; `CHANGELOG.md ## [Unreleased]` Added-Eintrag; `roadmap.md` §Nächste Schritte 2 (v0.3.0-Milestone) markiert [`slice-v1-add-remove`](slice-v1-add-remove.md) ✅ mit T1..T4-Hashes. `make docs-check` grün. |
| Review | dieser Commit | Code-Review-Followup: sechs Findings (F1..F6) direkt am Slice gepatcht — Two-Phase Plan-then-Write + InconsistentBlock-Convergence (F1), Doc-Comment-Truth (F2), File-Mode-Preservation via Lstat (F3), WARNING auf stderr statt NOTE auf stdout (F4+F5), Gate-Skip für Deactivated+Purge (F6). Siehe Review-Followup-Tabelle unten. |

## Review-Followup (post-T5)

Code-Review (xhigh-effort, 5 Angles, gezielter Scope auf
`detectServiceState`-Refactor / Partial-Write / `--purge`-UX) hat
sechs Findings ergeben. Im Followup-Commit gefixt (gleicher Branch,
kein eigener Slice — selber Pattern wie post-template-list/init):

| F | Datei | Befund | Fix |
| - | ----- | ------ | --- |
| F1 | `application/removeservice.go` | Mid-Flight-Write-Failure ließ Projekt in `InconsistentBlock` zurück; Retry war blockiert weil das Dispatch `InconsistentBlock` als Reject-State behandelte. | (a) Two-Phase Plan-then-Write: alle Reads + managedblock.Replaces + PatchScalar passieren in-memory; Phase 2 schreibt erst nach Erfolg aller Plan-Steps. Plan-time-Fehler (Malformed-Block, Patch-Fehler) short-circuiten vor jedem Disk-Write. (b) `InconsistentBlock` aus Reject-Liste raus, in proceeding-Switch aufgenommen — Remove konvergiert forwards (removeBlock ist idempotent; PatchScalar setzt enabled=false). Asymmetrisch zu Add, das InconsistentBlock weiter rejected (Add hat keinen klaren Konvergenz-Ziel). Neuer Test `InconsistentBlockConverges` ersetzt den alten Reject-Pin. |
| F2 | `application/removeservice.go` | Doc-Comment versprach "project stays self-consistent for a retry" — Lüge, weil compose-removed + env-kept + enabled=true ist `InconsistentBlock`. | Comment auf den tatsächlichen Mechanismus umgeschrieben: F1's InconsistentBlock-Convergence macht den Retry funktionsfähig; die "self-consistent"-Behauptung ist gestrichen. |
| F3 | `application/removeservice.go` | File-Mode (0o644 hardcoded) statt Original-Mode preserved — Hardened-Setups mit 0o600 wären silent downgraded. Asymmetrisch zu Add das via `loadForPatch` preserved. | Neuer `fileMode(path, fallback)`-Helper liest die Original-Mode via `s.fs.Lstat(path).Mode().Perm()` und gibt sie ans `WriteFile` durch. Fallback nur bei `iofs.ErrNotExist` (TOCTOU-Race; Apply-Loop würde sowieso scheitern). |
| F4 | `cli/remove.go` | NOTE-Block ging auf stdout — würde zukünftige `--json`-Pipes ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) / [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) verpesten. | `runRemove`-Signatur um `errOut io.Writer` erweitert; `printRemoveSummary` schreibt Success-Summary auf `out`, WARNING-Block auf `errOut`. Cobra-Closure übergibt `cmd.ErrOrStderr()`. Smoketest gegen das Image: stdout enthält nur die Success-Zeilen, stderr den WARNING-Block. |
| F5 | `cli/remove.go` | "NOTE:"-Präfix war zu mild; passive Sprache ("was requested; volume removal is not auto-handled") konnte überlesen werden. | "WARNING:" als Präfix; explizite Sprache "volume removal is NOT yet automated in v0.3.0" + "are still on disk and untouched". Volume-`ls`-Empfehlung um `--filter label=com.docker.compose.project=<...>` ergänzt, damit der User nicht alle Volumes des Hosts gezeigt bekommt. |
| F6 | `application/removeservice.go` | Gate feuerte für `Deactivated + Purge` und konsumierte User-Confirmation für eine No-Op — User sagte YES zu einem Prompt, der ins Nichts führte. | Gate-Position aus dem zentralen Pfad raus, ins proceeding-Branch geschoben: Gate fires nur bei Active/EnabledUnset/InconsistentBlock (Pfade die wirklich State transitionen). Deactivated-Branch ist gate-frei. Bestehender Test `PurgeOnDeactivatedStateRequiresGate` invertiert: zu `PurgeOnDeactivatedSkipsGate` (Confirmer wird NICHT gerufen). |

Gates nach Followup grün: lint + test + coverage-gate 90% +
docs-check. E2E-Smoketest gegen das gebaute Image:
- `init demo --no-git` + `add postgres` + `remove postgres --yes
  --purge`
- stdout: "Removed service \"postgres\"" + Changed-Liste (compose/
  env/u-boot.yaml) — KEIN WARNING
- stderr: WARNING-Block mit `docker volume ls --filter` +
  `docker volume rm` Cleanup-Hinweis

## Out of Scope

- **Dependency-Check** („Abhängigkeiten anderer Services prüfen
  und vor dem Entfernen warnen" aus [`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen)): heute hat
  Postgres keine Dependents (es ist die einzige Add-on-Option).
  Mit Keycloak (`requires: [postgres]`) wird das relevant —
  eigener Slice [`slice-v1-addons-deps`](slice-v1-addons-deps.md) deckt die Mechanik ab
  und ergänzt sie hier nachträglich.
- **Echte Volume-Löschung** (statt nur Confirmation-Gate +
  CLI-Cleanup-Hint): T3-Decision verschiebt die Volume-Entfernung
  auf einen Folge-Slice, weil sie eine Docker-Engine-Port-
  Erweiterung (`RemoveVolumes(ctx, names)` oder `ComposeDown
  --volumes`-Wrapping) braucht und in v0.3.0 ohne klaren
  Trigger nicht lohnt. T4-CLI sagt dem Nutzer explizit
  `docker volume rm <name>`.
- **`u-boot remove` für Keycloak/OTel**: kommt automatisch
  durch denselben Pfad, sobald die Add-ons existieren — kein
  eigener Slice nötig, solange das `add`-State-Machine-Pattern
  symmetrisch zum `remove` ist.

## Bezug

- Spec: [`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (V1) — vollständig geliefert für
  `postgres`; Dependency-Check (Sub-Punkt der Spec)
  deferred an [`slice-v1-addons-deps`](slice-v1-addons-deps.md); Volume-Removal
  deferred an einen Folge-Slice (T3-Decision).
- M5-Vorbild:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  liefert das State-Machine-Pattern. `detectServiceState` ist
  mit T2 zu einer Package-Free-Function refaktoriert, beide
  Services nutzen sie ohne Duplikation.
- M6-Vorbild:
  [`slice-m6-up-down`](../done/slice-m6-up-down.md) T6 liefert
  das `Confirmer.ConfirmRemoveVolumes`-Pattern, das mit T3
  hier wiederverwendet ist.
- Milestone: v0.3.0 „Add-on Catalogue Expansion" (siehe
  [roadmap.md §Nächste Schritte 2](../in-progress/roadmap.md)).
  Erster Slice des Milestones; vier weitere folgen.
- Phase: V1 (nach v0.2.0); kein Carveout.
