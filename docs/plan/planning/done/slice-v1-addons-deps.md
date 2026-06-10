# Slice V1: Add-on-Abhängigkeiten ([`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten))

## Auslöser

Dritter Slice des v0.3.0-Milestones. Voraussetzung für
[`slice-v1-keycloak`](slice-v1-keycloak.md): Keycloak deklariert eine optionale Postgres-
Dependency über `services.keycloak.persistence: external-postgres`
(Spec-Beispiel aus [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten)). Ohne den Dependency-
Mechanismus müsste der Keycloak-Slice den Mechanismus mit-bauen
oder still-laufen lassen — beides würde den Keycloak-Scope
verwässern. Dieser Slice baut die Mechanik isoliert, mit dem
heutigen Postgres-Add-on als no-op-Pfad (Postgres hat keine
Deps).

Spec-Anforderungen aus [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) (V1):

1. Add-on-Abhängigkeiten erkennen (z. B. Keycloak →
   PostgreSQL bei `services.keycloak.persistence: external-postgres`).
2. Bei fehlender Dep darf der Aufruf nicht stillschweigend
   fortfahren.
3. Vier-Modi-Behandlung:
   - `--with-deps` → fehlende Deps automatisch hinzufügen
   - `--yes` (ohne `--with-deps`) → Default-Entscheidung
     „Abhängigkeit hinzufügen" deterministisch
   - `--no-interactive` (ohne `--yes`/`--with-deps`) → Exit 10
   - Default interaktiv → Confirmer-Prompt
4. `--with-deps` ist mit `--no-interactive` kombinierbar
   (deterministisch ohne Rückfrage).

## Aufhebungsbedingung

- Domain-Typ `AddOnDependency` modelliert eine Dependency-
  Deklaration (Pfad-bedingte Service-Abhängigkeit).
- Application-Phase erkennt fehlende Deps nach State-Detection
  und vor Execute; surface über Response-Field oder Sentinel.
- CLI-Flag `--with-deps` auf `add` ergänzt; Vier-Modi-Logik
  implementiert.
- Postgres-Pfad bleibt verhaltensidentisch (no deps → no-op);
  bestehende Tests laufen unverändert.
- Synthetische Tests demonstrieren die Vier-Modi für eine
  Fake-Add-on-mit-Dep-Konfiguration. Echter E2E-Test mit
  realer Dep-Konfiguration landet im Keycloak-Slice.

## Akzeptanzkriterien

- Domain `AddOnDependency` mit Konstruktor + Validate
  (Pflichtfelder: `Requires`, `WhenPath`, `EqualsValue`).
- Application: Dependency-Resolver (Funktion oder Methode) der
  current `ubootYAMLConfig` + Add-on-Deklarationen gegen die
  intendierte Add-Aktion prüft und eine Liste fehlender
  Services zurückgibt. Integration in `AddServiceService.Add()`
  zwischen Detect-Phase und Execute.
- Neuer Driving-Sentinel `ErrDependenciesRequired` (Exit 10)
  für den `--no-interactive`-Fail-Fast-Pfad.
- Neuer Confirmer-Method-Pattern `ConfirmAddDependency(ctx, svc,
  deps)` (oder Wiederverwendung des bestehenden Confirmer-Ports
  mit angepasstem Prompt — T3-Entscheidung).
- CLI `--with-deps`-Flag mit Vier-Modi-Logik analog zum
  `--purge`-Pattern aus M6/v1-add-remove. Rekursive
  Sub-`Add`-Calls für fehlende Deps via interner
  `AddServiceUseCase`-Wiederverwendung.
- Pure-Function-Tests für den Resolver (synthetic
  AddOnDependency-Inputs ohne Catalogue-Touch).
- Integration-Test mit Postgres (verifiziert no-deps-no-op:
  bestehende Add-Tests bleiben grün).
- Synthetic-Catalogue-Tests demonstrieren die Vier-Modi
  (Mock-Add-on mit Postgres-Dep, alle vier User-Modi durchspielen).
- E2E-Smoke gegen das gebaute Image: `init demo --no-git` +
  `add postgres` (no-deps-Pfad) → unverändert wie vor diesem
  Slice.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `23abd2b` | Domain-Typ `domain.AddOnDependency` (Requires/WhenPath/EqualsValue) mit `Validate` + Sentinel `domain.ErrInvalidAddOnDependency`. Application-Side-Table `dependenciesFor(svc) []domain.AddOnDependency` in addservice.go (Postgres → nil, MVP). `//nolint:unparam` mit Begründung „erste echte Row landet in [slice-v1-keycloak](slice-v1-keycloak.md)". Tests: Validate 100%, `dependenciesFor(postgres)` → nil, `dependenciesFor(ghost)` → nil. KEINE Verhaltensänderung in `Add()`. |
| T2 | `cd4f88c` | Pure-Resolver `resolveAddDependencies(cfg, deps) []ServiceName`: WhenPath-Match (Trigger-Bedingung) UND Service-NICHT-in-cfg.Services → missing-Liste, dedupliziert in Insertion-Order. Helper `resolveScalarPath` deckt `project.name`, `devcontainer.enabled`, `services.<svc>.enabled` ab (v0.3.0-Scope). Integration in `Add()` zwischen `detectServiceState` und State-Switch über `checkAddDependencies(baseDir, svc, deps)` (load + resolve + wrap Missing in `ErrDependenciesRequired`). Neuer Driving-Sentinel `ErrDependenciesRequired` (Exit 10, in `isServiceValidationError` aufgenommen). Tests via neuem `addservice_dependencies_test.go` (export_test-Seam `ResolveAddDependenciesForTest` / `ResolveScalarPathForTest` / `CheckAddDependenciesForTest`). Postgres-Tests bleiben grün (no-deps short-circuit). |
| T3 | `41b51ed` | CLI `--with-deps` BoolVar auf `add` + Plumbing in `AddServiceRequest{WithDeps, Yes, NoInteractive}`. Application: T2's `checkAddDependencies` zu Orchestrator refaktoriert + `findMissingDependencies` (load+resolve) + `handleMissingDependencies` (Vier-Modi-Dispatch): WithDeps OR Yes → autoInstall (rekursive `Add`-Calls; Flags vererben sich auf Sub-Requests); NoInteractive ohne Yes/WithDeps → Fail-Fast mit `ErrDependenciesRequired`; default → `Confirmer.ConfirmAddDependency`-Prompt → YES promotet zu autoInstall, NO → ErrDependenciesRequired. Neue Driven-Port-Method `ConfirmAddDependency(ctx, svc, missing []string)` (mirror von `ConfirmRemoveVolumes`); Production-Adapter + noopConfirmer + fakeConfirmer extended. **Breaking**: `NewAddServiceService` nimmt jetzt einen Confirmer zwischen yaml und logger; alle 8 Callsites (main.go, e2e, 4 application-Tests) angepasst. Tests: 7 Dispatch-Arme über neuen `HandleMissingDependenciesForTest` (ghost-service als unsupported sub-target ist Beweis-Marker für recursive-Add ohne Disk-Setup); 4 Production-Confirmer-Adapter-Tests; 4 CLI-Plumbing-Tests inkl. `--with-deps`/`--yes`/`--no-interactive`. lint + test + coverage (90.10%) grün. |
| T4 | dieser Commit | Closure: README.{md,de.md} subcommand-Reference erwähnt `--with-deps`; CHANGELOG `## [Unreleased]` Added-Eintrag mit [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten)-Bezug + Verweis auf diesen Slice-Plan; roadmap.md v0.3.0-Milestone-Tabelle markiert [`slice-v1-addons-deps`](slice-v1-addons-deps.md) ✅ mit T1..T3-Hashes und bumpt Stand auf 3/5; Slice-Plan `open/` → `done/`. `make docs-check` grün. |

## Out of Scope

- **Keycloak-/OTel-Add-on-Implementation**: jeweils eigener
  Slice ([`slice-v1-keycloak`](slice-v1-keycloak.md), [`slice-v1-otel`](slice-v1-otel.md)). Dieser Slice
  baut nur den Mechanismus.
- **Rekursive Dep-Auflösung über mehrere Ebenen**: heute reicht
  eine Ebene (Keycloak → Postgres). Falls Postgres selbst
  Deps bekäme (unwahrscheinlich), wäre das eigene Slice-Arbeit.
- **Dep-Removal beim `remove`**: [`LH-FA-ADD-007`](../../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) mentioned
  Dependency-Warn-on-Remove; das integriert sich später nach
  diesem Slice oder bei Keycloak.
- **OTel sample-config Generation für App-Services**: aus dem
  Spec-Text scheint das eher ein OTel-Feature-Detail zu sein
  als eine Dependency. Bleibt in [`slice-v1-otel`](slice-v1-otel.md) zu klären.

## Bezug

- Spec: [`LH-FA-ADD-006`](../../../../spec/lastenheft.md#lh-fa-add-006--add-on-abhängigkeiten) (V1).
- Voraussetzungs-Slice: keine (steht für sich).
- Wird-genutzt-von-Slice:
  [`slice-v1-keycloak`](../in-progress/roadmap.md) (kommt nach
  diesem Slice; Keycloak deklariert Postgres-Dep via diesem
  Mechanismus).
- M5-Vorbild:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  liefert das State-Machine-Pattern, in das die Dep-Phase
  zwischen Detect und Execute eingehängt wird.
- M6 + v1-add-remove-Vorbild für Confirmer-Gate:
  `ConfirmRemoveVolumes` ist das parallele Pattern; neue Methode
  `ConfirmAddDependency` folgt derselben Signatur-Familie.
- Milestone: v0.3.0 „Add-on Catalogue Expansion".
- Phase: V1 (nach [`slice-v1-add-remove`](slice-v1-add-remove.md) und [`slice-v1-audit-done`](slice-v1-audit-done.md)).
