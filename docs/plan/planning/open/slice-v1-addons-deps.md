# Slice V1: Add-on-Abhängigkeiten (LH-FA-ADD-006)

## Auslöser

Dritter Slice des v0.3.0-Milestones. Voraussetzung für
`slice-v1-keycloak`: Keycloak deklariert eine optionale Postgres-
Dependency über `services.keycloak.persistence: external-postgres`
(Spec-Beispiel aus `LH-FA-ADD-006`). Ohne den Dependency-
Mechanismus müsste der Keycloak-Slice den Mechanismus mit-bauen
oder still-laufen lassen — beides würde den Keycloak-Scope
verwässern. Dieser Slice baut die Mechanik isoliert, mit dem
heutigen Postgres-Add-on als no-op-Pfad (Postgres hat keine
Deps).

Spec-Anforderungen aus `LH-FA-ADD-006` (V1):

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

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | Domain-Typ `domain.AddOnDependency` + Konstruktor + Validate (Pflichtfelder + Sentinel `domain.ErrInvalidAddOnDependency`). Application-Side-Table `dependenciesFor(svc) []domain.AddOnDependency` in addservice.go (returns nil für postgres). Tests: 100% Coverage auf Validate; Pure-Function-Tests für `dependenciesFor` (postgres → nil). KEINE Verhaltensänderung in `Add()` — nur Domain-Vorbereitung. |
| T2 | Application: Pure-Resolver-Funktion `resolveAddDependencies(svc, cfg ubootYAMLConfig, deps []AddOnDependency) []ServiceName` (gibt fehlende Service-Namen zurück basierend auf YAML-Path-Match + Service-Präsenz in cfg.Services). Integration in `AddServiceService.Add()` nach `detectServiceState` + vor `executeAdd`: bei `len(missing) > 0` und nicht-leerem `req.DependencyMode` (T2 fügt das Feld zur Request hinzu), entsprechend dispatchen. Neuer Driving-Sentinel `ErrDependenciesRequired` (Exit 10) für den Fail-Fast-Pfad. Tests: Pure-Function (Resolver mit synthetischen Inputs), Application-Integration mit Fake-Resolver (deps>0 → Sentinel/Response/Recursion-Marker). Postgres-Add-Tests bleiben unverändert (no-deps → kein Code-Pfad-Wechsel). |
| T3 | CLI: `--with-deps` BoolVar auf `add`; runAdd Vier-Modi-Dispatch (analog `--purge`-Pattern aus remove): WithDeps OR Yes → auto-install via rekursive `Add`-Calls für jede missing Dep; NoInteractive ohne Yes/WithDeps → ErrDependenciesRequired Exit 10; default interaktiv → confirmer.ConfirmAddDependency-Prompt → bei YES rekursive Add-Calls, bei NO ErrDependenciesRequired. Neue Confirmer-Methode `ConfirmAddDependency(ctx, svc, deps)` (extends `driven.Confirmer`-Port — bricht bestehende fakeConfirmer; alle Test-Helper updaten). isValidationError ergänzt ErrDependenciesRequired. CLI-Tests inkl. Vier-Modi + rekursive-Add-Verification. Smoke-Test gegen das gebaute Image: `add postgres` (no-deps) bleibt unverändert. |
| T4 | Closure: README.{md,de.md} `add`-Bullet erwähnt `--with-deps`; CHANGELOG `## [Unreleased]` Added-Eintrag; roadmap.md v0.3.0-Milestone-Tabelle markiert `slice-v1-addons-deps` ✅ und bumpt Stand auf 3/5; Slice-Plan `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. |

## Out of Scope

- **Keycloak-/OTel-Add-on-Implementation**: jeweils eigener
  Slice (`slice-v1-keycloak`, `slice-v1-otel`). Dieser Slice
  baut nur den Mechanismus.
- **Rekursive Dep-Auflösung über mehrere Ebenen**: heute reicht
  eine Ebene (Keycloak → Postgres). Falls Postgres selbst
  Deps bekäme (unwahrscheinlich), wäre das eigene Slice-Arbeit.
- **Dep-Removal beim `remove`**: `LH-FA-ADD-007` mentioned
  Dependency-Warn-on-Remove; das integriert sich später nach
  diesem Slice oder bei Keycloak.
- **OTel sample-config Generation für App-Services**: aus dem
  Spec-Text scheint das eher ein OTel-Feature-Detail zu sein
  als eine Dependency. Bleibt in `slice-v1-otel` zu klären.

## Bezug

- Spec: `LH-FA-ADD-006` (V1).
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
- Phase: V1 (nach `slice-v1-add-remove` und `slice-v1-audit-done`).
