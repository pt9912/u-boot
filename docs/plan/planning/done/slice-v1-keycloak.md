# Slice V1: `u-boot add keycloak` (LH-FA-ADD-003 + LH-AK-003)

## Auslöser

`LH-FA-ADD-003` (V1) verlangt, dass `u-boot add keycloak` einen
Keycloak-Service mit Compose-Block, Admin-Env-Block (Placeholder-
Secrets), Port-Konfiguration und Healthcheck in ein initialisiertes
Projekt einbaut. `LH-AK-003` ist der zugehörige Acceptance-Flow
(`init` + `add keycloak` + `up` + Endpoint-Probe auf Port 8080).

Vierter Slice des v0.3.0-Milestones, baut direkt auf dem
`slice-m5-add-postgres`-Pattern auf. Die optionale Postgres-
Anbindung aus `LH-FA-ADD-003` §857 („optionale PostgreSQL-
Anbindung bei konfigurierter persistenter externer Datenbank")
ist **out-of-scope** für diesen Slice und an einen Folge-Slice
abgegeben (siehe Out-of-Scope).

## Aufhebungsbedingung

`u-boot add keycloak` in einem initialisierten Projekt liefert:

1. Keycloak-Service-Block in `compose.yaml` mit
   - Image `quay.io/keycloak/keycloak:26.0` (gepinnter LTS-Tag —
     T1-Decision dokumentiert),
   - Port-Mapping `8080:8080`,
   - Healthcheck (T1-Decision: `curl --fail http://localhost:9000/
     health/ready` falls Image die Management-Port-9000-Probe
     mitbringt, sonst entfällt der Healthcheck per LH-FA-ADD-003
     §858 „soweit technisch sinnvoll"),
   - Env-Block mit Admin-Credentials aus `.env`,
   - **Flüchtige H2-In-Container-Persistenz, kein Volume** —
     Keycloak's Default-Embed-DB läuft im Container-Filesystem und
     ist nach `docker compose down` weg; LH-AK-003 verlangt nur
     Endpoint-200/302, keine Persistenz. Persistente externe
     Postgres-Anbindung ist eigener Folge-Slice.
2. `.env.example` mit `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` +
   `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD` als
   Placeholder-Secrets (Spec-mandatiert: niemals reale Secrets im
   Repo).
3. `services.keycloak.enabled: true` in `u-boot.yaml`.

LH-AK-003-Akzeptanz: `init demo --no-git` + `add keycloak` + `up`
auf einem Docker-fähigen System bringt Keycloak-Endpoint
`http://localhost:8080/` mit HTTP 200 oder 302 hoch (Boot-Zeit
30–90 s JVM-Init, siehe T3-Timeout-Carveout).

## Akzeptanzkriterien

- ✅ `u-boot add keycloak` (Service noch nicht aktiv) erzeugt zwei
  Mutations: Compose-Block neu, env-Block neu, `services.keycloak.
  enabled: true` in `u-boot.yaml`. **Kein** Volume-Mutation
  (Default-Persistenz embedded; siehe oben). Idempotent: zweiter
  Aufruf → „already active; no changes" + Exit 0; **insbesondere
  KEIN actionRepairArtifacts-Trigger wegen fehlendem volume.
  keycloak-Marker** (Heutiger inspectVolumeArtefact-Pfad ist
  postgres-spezifisch — T2 behebt das, siehe T2).
- ✅ `u-boot add keycloak` ohne `u-boot.yaml` failt mit
  `ErrProjectNotInitialized` + Exit 10 (vorhandener Code-Pfad).
- ✅ `u-boot add keycloak --with-deps`: `dependenciesFor(keycloak)`
  returnt heute `[]` (Default-Persistenz embedded → kein
  Postgres-Bedarf); der `--with-deps`-Pfad läuft normal durch und
  installiert null Zusatz-Services. `--with-deps` bleibt
  funktional vorbereitet für den Folge-Slice, der die externe
  Postgres-Anbindung aktiviert.
- ✅ `u-boot remove keycloak` funktioniert reziprok (existierender
  RemoveServiceService-Pfad greift automatisch, sobald Keycloak in
  der Catalogue ist — Voraussetzung: T2 hat die Volume-Skip-
  Behandlung umgesetzt, sonst meldet executeRemove einen
  fehlenden Volume-Block als InconsistentBlock).
- ✅ E2E-Acceptance-Test `keycloak_acceptance_docker_test.go`
  analog zur Postgres-Acceptance: Docker-only-Test (`//go:build
  docker`-Tag), init + add + up, HTTP-Probe gegen Port 8080
  erwartet 200 oder 302.
- ✅ Doctor-Checks: `services.keycloak.enabled` taucht in
  `services.enabled`-Diagnostic auf (kein Service-Hardcoding
  nötig — `collectActiveServicePorts` (`doctor.go:986`) ist
  bereits generisch über Compose-Service-Namen);
  `devcontainer.forwardPorts.consistency` warnt höchstens für
  Port 8080-Mismatch, eskaliert nie zu Error.
- ✅ Hexagonale Verdrahtung: kein neuer Driving-Port nötig (Reuse
  AddServiceUseCase + RemoveServiceUseCase); Catalogue-Erweiterung
  + Template-Render-Generalisierung + **Detect-Pfad-Generalisierung**
  sind die Code-Änderungen außerhalb der Test- und Template-
  Dateien.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `2708606` | **Templates + Render-Generalisierung + T1-Sub-Decisions.** Drei Templates anlegen — `templates/services/keycloak.compose.tmpl` (Image `quay.io/keycloak/keycloak:26.0` — T1-Decision LTS-Pin; Port-Mapping `8080:8080`; Healthcheck-Entscheidung dokumentiert: bevorzugt `curl --fail http://localhost:9000/health/ready` mit `--health-cmd` analog Postgres-Template, sonst weglassen; Env-Block; `command: ["start-dev"]` als Compose-Command für die LH-AK-003-Dev-Mode-Boot-Variante), `keycloak.env.tmpl` (`KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` + `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`), **kein** `keycloak.volume.tmpl` (siehe T2 Skip-Volume). Service-Catalogue-Tabelle in `addservice_execute.go` (oder neue `servicecatalogue.go`) mit `composeTmpl`, `envTmpl`, `volumeTmpl` (optional) pro Service einführen; `renderPostgresTemplates(svc)` zu `renderServiceTemplates(svc)` umbenennen und über den Catalogue-Lookup auflösen. **Catalogue für `isSupportedService` wird NOCH NICHT erweitert** — Keycloak bleibt nach T1 weiterhin `ErrServiceUnsupported`, weil sonst `u-boot add keycloak` durch den Postgres-only Detect-Pfad läuft (F1-Befund: Endlos-Repair durch `inspectVolumeArtefact` + falsche `hasRequiredEnvKeys`). Catalogue-Erweiterung kommt mit T2 zusammen mit der Detect-Generalisierung. **T1-Sub-Decision: Template-Parametrisierung bleibt inline-hardcoded** wie heute Postgres (`templateData{}` leer), kein `{{ .Name }}`-Refactor — falls späterer Image-Tag-Override (z. B. via `services.keycloak.imageTag`) gewünscht, ist das eigener Slice. `embed.FS`-Eintrag in `templates.go` für Template-Resolution erweitern (Volume-Template wird **per Service** optional). Tests: Template-Existenz-Pin pro Service (Postgres alle 3, Keycloak 2), Render-Smoke für beide Catalogue-Einträge, **expliziter Pin auf `isSupportedService("keycloak") == false`** (T2-Voraussetzungs-Pin); Postgres-Render-Output-Byte-Identity-Pin (no-behavior-change). |
| T2 | `861f231` | **Catalogue-Erweiterung + Per-Service-Probe-Mechanismus.** Heute ist der Detect-Pfad an drei Stellen Postgres-hardcoded — alle drei müssen pro Service konfigurierbar werden: (a) `hasRequiredEnvKeys` (`addservice_detect.go:182`) prüft literal `POSTGRES_USER`/`POSTGRES_PASSWORD`/`POSTGRES_DB` — wird zu `hasRequiredEnvKeysFor(svc, blockBody)` mit Service-Catalogue-Lookup (`postgres → [POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB]`, `keycloak → [KEYCLOAK_ADMIN, KEYCLOAK_ADMIN_PASSWORD]`); (b) `contentScanState` / `feedSubBlockEntry` matchen `postgres-data` als Volume-Ref-Literal und `POSTGRES_*` als Environment-Keys — der Scan-State wird per Service parametrisiert (`requiredEnvKeys`, `volumeRefLiteral`, `skipVolume`-Flag); (c) `inspectVolumeArtefact` (`addservice_detect.go:98`) liefert heute für jeden Service ohne Volume `needsRepair=true` → für Keycloak Endlos-Repair-Schleife; Fix: Service-Catalogue-Eintrag `volumeOptional: true` → `inspectVolumeArtefact` skippt komplett, gibt `false, nil` (no-repair) zurück. **Erst hier** wird `isSupportedService("keycloak")` und `supportedServices()` erweitert. **Designentscheidung T2:** der T1-Service-Catalogue wächst um `requiredEnvKeys []string`, `volumeRefLiteral string`, `volumeOptional bool`. `executeRemove` (RemoveService) muss ebenfalls die Volume-Skip-Logik respektieren — heute liest es blind den Volume-Block; T2-Carveout: Volume-Removal bei `volumeOptional` skippen statt ErrServiceInconsistent. Tests: 100% Coverage auf den neuen Catalogue-Lookup, Catalogue-Pin (`TestSupportedServices` enthält jetzt Keycloak), **Postgres-Snapshot-Pin (Byte-Identity der Postgres-Render-Outputs unverändert)**, Keycloak-Probe-Pin (Active-Detection für Keycloak ohne Volume → kein `needsRepair`). |
| T3 | `beb222b` | **E2E-Acceptance + Test-Helper-Extraktion.** Beim zweiten Acceptance-Docker-Test landet typisch Copy-Paste — **T3-Decision:** Postgres-spezifische Compose-Up + Endpoint-Probe-Logik aus `postgres_acceptance_docker_test.go` in `internal/e2e/acceptance_helpers.go` extrahieren (signed-off via expliziter T3-Commit-Message), Postgres-Test wird mit-refaktoriert. Keycloak-Test ruft Helper auf, definiert nur die Keycloak-spezifische Probe (`probeEndpoint("http://localhost:8080/", 200, 302)`). **Boot-Zeit-Carveout:** Keycloak JVM-Boot 30–90 s (in CI gerne mehr) vs. Postgres ~5 s — UpService-Default-Healthcheck-Timeout im Test übersteuern; falls `upservice_healthcheck_docker_test.go`-Konventionen ein generisches Timeout-Override nicht erlauben, T3-Commit liefert `WithHealthcheckTimeout`-Option an UpService (sonst Test-spezifischen Override per `t.Cleanup` mit längerem Context-Timeout). Manuelle Smoke-Anleitung im CHANGELOG-Eintrag für Repo-Owner-Spot-Check. Doctor-Smoke: `add keycloak` darf keine neuen Errors auslösen (verifiziert: `collectActiveServicePorts` ist generisch, `forwardPorts.consistency` warnt höchstens). |
| T4 | dieser Commit | **Closure.** READMEs (`add <service>`-Subcommand-Reference erwähnt jetzt `postgres \| keycloak`), CHANGELOG `## [Unreleased]` Added-Eintrag mit LH-FA-ADD-003 + LH-AK-003-Bezug + Hinweis auf flüchtige Default-Persistenz + Folge-Slice-Verweis (`slice-v1-keycloak-external-postgres`), roadmap.md §v0.3.0-Tabelle markiert `slice-v1-keycloak` ✅ mit T1..T3-Hashes + Stand-Bump 4/5, Slice-Plan `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. |

## Out of Scope

- **`services.keycloak.persistence: external-postgres`-Schema-Feld
  + Postgres-Dep-Aktivierung**: `LH-FA-ADD-003` §857 sagt
  „optionale PostgreSQL-Anbindung bei konfigurierter persistenter
  externer Datenbank". Heute hat das `LH-FA-CONF-005`-Schema dafür
  keinen Slot; der Schema-Slot + die CLI-Setter +
  `dependenciesFor("keycloak")` mit dem konkreten `WhenPath` /
  `EqualsValue` sind eigene Slice-Arbeit. Plan-Anker: Folge-Slice
  `slice-v1-keycloak-external-postgres` (Trigger: Nutzer-Bedarf
  nach persistenter Keycloak-Datenbank).
- **Persistentes Default-Volume für Keycloak**: könnte als
  alternative Out-of-the-box-Persistenz (gegen den
  flüchtige-H2-Default) später eingeführt werden — bedeutet aber
  Festlegung auf einen Persistenz-Mechanismus, der heute nicht im
  Spec-Mindestumfang steht. Folge-Slice wenn der Bedarf
  konkretisiert wird; greift dann auf den T2-Catalogue-Slot
  `volumeOptional: false` zurück.
- **Custom-Realm-Import / Theme-Konfiguration**: nicht in
  LH-FA-ADD-003 Mindestumfang. Nutzer kann nach `add keycloak`
  manuell zusätzliche Volumes / Env-Vars im managed-block-freien
  Bereich der compose.yaml ergänzen.
- **CI-Stabilisierung des Keycloak-Acceptance-Tests**: in GitHub-
  Actions failt `docker compose up` für Keycloak reproduzierbar
  nach <1 s mit „compose runtime error" — vermutlich Quay.io-
  Manifest-Lookup (Outage heute 2026-06-01 den ganzen Tag; lokal
  ebenfalls 502/504 gesehen). T3-Carveout: der Test ist mit
  zusätzlichem build-tag `acceptance_extended` opt-in, default
  `make test-docker` umgeht ihn. Folge-Slice
  `slice-v1-keycloak-ci-flake` (Trigger: Compose-Verbose-Logs aus
  CI ziehen, dann entweder Pull-Retry-Wrapper im UpService oder
  Quay-Mirror via Docker-Hub-Pull-Through-Cache).
- **OpenTelemetry-Add-on**: eigener Slice `slice-v1-otel`
  (LH-FA-ADD-004 + LH-AK-004), parallel-entwickelbar weil keine
  Dep zwischen Keycloak und OTel.
- **`{{ .Name }}`-Template-Parametrisierung**: T1 wählt explizit
  inline-hardcoded Service-Namen (Variante a aus dem
  `renderPostgresTemplates`-Doc-Comment). Variable-Templates +
  Image-Tag-Override (z. B. `services.keycloak.imageTag`) sind
  eigene Slice-Arbeit, weil sie das `LH-FA-CONF-005`-Schema
  erweitern.

## Bezug

- Spec: `LH-FA-ADD-003` (V1, §838-§859) + `LH-AK-003` (V1, §2336-§2353).
- Voraussetzungs-Slices: keine harten (Catalogue +
  Template-Refactor + Detect-Generalisierung sind in diesem Slice
  enthalten); weiche Voraussetzungen:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  liefert das Pattern, [`slice-v1-add-remove`](../done/slice-v1-add-remove.md)
  liefert den reziproken Remove-Pfad, [`slice-v1-addons-deps`](../done/slice-v1-addons-deps.md)
  liefert die Dep-Mechanik (heute hier inaktiv — leere Dep-Liste
  für Keycloak).
- Folge-Slices:
  - `slice-v1-otel` (LH-FA-ADD-004 + LH-AK-004) — parallel
    entwickelbar.
  - `slice-v1-keycloak-external-postgres` — Schema-Erweiterung
    (`services.keycloak.persistence`) + Dep-Aktivierung (siehe
    Out-of-Scope).
- Milestone: v0.3.0 „Add-on Catalogue Expansion" (siehe
  [roadmap.md §v0.3.0](../in-progress/roadmap.md#v030)). Vierter
  Slice des Milestones (3/5 ✅ vor diesem Slice).
- Phase: V1 (nach v0.2.0); kein Carveout.
