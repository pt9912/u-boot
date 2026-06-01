# Slice V1: `u-boot remove <service>`

## Auslöser

`LH-FA-ADD-007` (V1) verlangt einen Mechanismus, ein Add-on
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
„nothing to do"-Meldung ohne Fehler. `--purge` entfernt
zusätzlich das benannte Volume (destruktiv, mit LH-FA-CLI-005A-
Confirmation-Gate analog `down --volumes`).

## Akzeptanzkriterien

- `u-boot remove postgres` (Service aktiv) erzeugt drei
  Removal-Actions: Compose-Block weg, env-Block weg,
  `services.postgres.enabled: false`. Volumes bleiben (kein
  Daten-Verlust ohne `--purge`).
- `u-boot remove postgres` (Service auf enabled: false) liefert
  idempotent eine Meldung („service is already disabled; no
  changes") + Exit 0; KEINE Dateien werden angefasst.
- `u-boot remove postgres` (Service nicht registriert / inkonsistent)
  failt mit klarem Sentinel + Exit 10
  (`ErrServiceUnregistered` analog M5
  `ErrServiceInconsistent`).
- `u-boot remove unknown-svc` failt mit `ErrServiceUnsupported`
  (Exit 10) — same as M5 `add` für unbekannte Services.
- `u-boot remove postgres --purge` zusätzlich: Volume gelöscht.
  Confirmation-Gate analog `down --volumes` (`LH-FA-CLI-005A` §254):
  non-interactive ohne `--yes` → `ErrConfirmationRequired`
  Exit 10; interaktiv → Default-`N`-Prompt.
- `u-boot remove postgres` ohne `u-boot.yaml` failt mit
  `ErrProjectNotInitialized` (Exit 10) — analog M5.
- `u-boot.yaml` nach erfolgreichem Remove zeigt
  `services.postgres.enabled: false` (NICHT entfernt) —
  konsistent zur State-Machine-Idee „registriert aber inaktiv".
- Hexagonale Verdrahtung: neuer Driving-Port
  `RemoveServiceUseCase` + Application `RemoveServiceService`
  (spiegelt M5 `AddServiceService`-Pattern) + CLI-Subkommando.
  Wiring in `cmd/uboot/main.go`.

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | Driving-Port `port/driving.RemoveServiceUseCase` mit `Remove(ctx, RemoveServiceRequest) (RemoveServiceResponse, error)`. Request `{BaseDir, Service domain.ServiceName, Purge bool, Yes bool, NoInteractive bool}`. Response `{Service, PriorState, State domain.ServiceState, Changed []string, VolumesPurged bool}`. Neue Sentinels: `ErrServiceUnregistered` (Exit 10, Service nie registriert), `ErrServicePurgeRequiresConfirm` (NICHT eigener Sentinel — wiederverwendet `ErrConfirmationRequired` aus M6). Application-Skeleton `RemoveServiceService` mit Konstruktor `NewRemoveServiceService(fs, yaml, confirmer, logger)`; Static-Check gegen den Port. Application-Unit-Tests für Konstruktor + Skeleton (nicht-State-Machine). |
| T2 | Application-State-Machine: Detect-Phase ermittelt PriorState via existierender `domain.ServiceState`-Konstanten (Unregistered/Active/Deactivated/EnabledUnset). Execute-Phase verzweigt: PriorState=Active → 3 Actions (Compose-Block-Remove via managedblock, env-Block-Remove, u-boot.yaml-Patch auf enabled: false), PriorState=Deactivated → No-Op-Response, PriorState=Unregistered/EnabledUnset → `ErrServiceUnregistered`-Wrap. Volumes bleiben unangetastet (ohne `--purge`). 8+ Application-Unit-Tests (jede ServiceState-Transition + Error-Pfade). |
| T3 | `--purge`-Flag und Confirmation-Gate: Wenn Purge && !Yes → `confirmer.ConfirmRemoveVolumes(ctx, service)` (existierende Confirmer-Methode aus M6 down). Bei Yes-Bestätigung: zusätzliche Action „remove volume <name>" über Compose-Block oder direkt via Docker-Engine-Port (T3-Decision: machen wir das via Compose-`down --volumes`-Equivalent oder über separaten `RemoveVolumes(ctx, names)`-Port?). Vorschlag: zunächst Compose-Block weg + Note „volumes still on disk; use `docker volume rm` to remove" — voll volume-removal kann in einem Folge-Slice landen, wenn Docker-Engine-Port-Erweiterung sich lohnt. Tests inkl. Confirmation-Gate-Pfade. |
| T4 | CLI-Subkommando `u-boot remove <service> [--purge]` in `internal/adapter/driving/cli/remove.go` (analog `add.go`-Pattern, BoolVar für Purge, persistent Yes/NoInteractive vom Root). `cli.New` 10. positional `removeUC driving.RemoveServiceUseCase` ergänzt; alle 8 bestehenden `newApp*`-Test-Helper + `fakeRemoveServiceUseCase`. ExitCode-Doc + `isValidationError` für `ErrServiceUnregistered`. CLI-Tests (Happy-Path, Idempotent-Path, Confirmation-Gate, ErrProjectNotInitialized). Wiring in `cmd/uboot/main.go`. Smoke-Test gegen das gebaute Image: `init` + `add postgres` + `remove postgres` + Verify `services.postgres.enabled=false`. |
| T5 | Closure: README.{md,de.md} bekommen `u-boot remove`-Bullet; `CHANGELOG.md ## [Unreleased]` Added-Eintrag; `roadmap.md` §Nächste Schritte 2 (v0.3.0-Milestone) markiert `slice-v1-add-remove` als ✅; Slice-Move `open/` → `done/` mit Tranchen-Commit-Spalte. `make docs-check` grün. |

## Out of Scope

- **Dependency-Check** („Abhängigkeiten anderer Services prüfen
  und vor dem Entfernen warnen" aus LH-FA-ADD-007): heute hat
  Postgres keine Dependents (es ist die einzige Add-on-Option).
  Mit Keycloak (`requires: [postgres]`) wird das relevant —
  eigener Slice `slice-v1-addons-deps` deckt die Mechanik ab
  und ergänzt sie hier nachträglich.
- **Echte Volume-Löschung** (statt nur Compose-Block-Removal):
  je nach T3-Decision; potentiell ein Folge-Slice
  `slice-v1-remove-purge-volumes` wenn der Docker-Engine-Port
  erweitert werden muss.
- **`u-boot remove` für Keycloak/OTel**: kommt automatisch
  durch denselben Pfad, sobald die Add-ons existieren — kein
  eigener Slice nötig, solange das `add`-State-Machine-Pattern
  symmetrisch zum `remove` ist.

## Bezug

- Spec: `LH-FA-ADD-007` (V1).
- M5-Vorbild:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  liefert das State-Machine-Pattern (`AddServiceService`,
  `detectState`, `executeAction`), das hier gespiegelt wird.
  M5-domain (`ServiceState`-Konstanten) bleibt wiederverwendet.
- M6-Vorbild:
  [`slice-m6-up-down`](../done/slice-m6-up-down.md) T6 liefert
  das `Confirmer.ConfirmRemoveVolumes`-Pattern, das hier für
  `--purge` wiederverwendet wird.
- Milestone: v0.3.0 „Add-on Catalogue Expansion" (siehe
  [roadmap.md §Nächste Schritte 2](../in-progress/roadmap.md)).
- Phase: V1 (nach v0.2.0); kein Carveout (`LH-FA-ADD-007` ist
  V1-Spec-ID, kein Carveout).
