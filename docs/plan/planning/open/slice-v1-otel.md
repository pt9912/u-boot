# Slice V1: `u-boot add otel` (LH-FA-ADD-004 + LH-AK-004)

## Auslöser

`LH-FA-ADD-004` (V1) verlangt, dass `u-boot add otel` einen
OpenTelemetry Collector mit Compose-Service, Collector-Konfig-
Datei, Standard-OTLP-Ports und Beispielkonfiguration für Logs,
Metrics und Traces in ein initialisiertes Projekt einbaut.
`LH-AK-004` ist der zugehörige Acceptance-Flow (`init` + `add
otel` + `up` + TCP-Probe auf 4317 + 4318).

Fünfter und letzter Slice des v0.3.0-Milestones, parallel-
entwickelbar zu [`slice-v1-keycloak`](./slice-v1-keycloak.md):
keine Dependency zwischen Keycloak und OTel, beide nutzen das
M5-Postgres-Pattern + die mit `slice-v1-keycloak` T2 gebauten
Per-Service-Probe-Mechanismen.

Diese Slice führt zusätzlich einen neuen Catalogue-Feld-Typ ein:
`extraFiles` — pro Service eine Liste von **rendered Datei-
Artifacts** abseits von `compose.yaml`/`.env.example`/Volume-
Block. Bei OTel ist das die `otel-collector-config.yaml`.

## Aufhebungsbedingung

`u-boot add otel` in einem initialisierten Projekt liefert:

1. OpenTelemetry-Collector-Service in `compose.yaml` mit
   - Image `otel/opentelemetry-collector:0.108.0` (gepinnter
     Stable-Tag — T1-Decision dokumentiert),
   - Port-Mappings `4317:4317` (OTLP/gRPC) + `4318:4318`
     (OTLP/HTTP),
   - `command: ["--config=/etc/otel-collector-config.yaml"]`,
   - Volume-Mount der gerenderten Config-Datei (read-only
     Bind-Mount auf die Datei im Projektverzeichnis),
   - **Kein Healthcheck**: LH-AK-004 verlangt nur Container-Status
     `running` ODER `healthy` — Collector-Default-Healthcheck-
     Endpoint braucht eine separate Extension-Konfiguration, die
     im Mindestumfang nicht gefordert ist.
2. `otel-collector-config.yaml` im Projektverzeichnis mit
   Mindest-Setup:
   - Receivers: `otlp` (grpc auf 4317, http auf 4318),
   - Processors: `batch`,
   - Exporters: `debug` (loggt auf stdout — Beispiel-Setup ohne
     externe Senken),
   - Service-Pipelines: `logs`, `metrics`, `traces` (je
     receiver→processor→exporter).
3. `.env.example` bleibt unverändert (OTel-Default-Setup braucht
   keine Admin-Credentials).
4. `services.otel.enabled: true` in `u-boot.yaml`.

LH-AK-004-Akzeptanz: `init demo --no-git` + `add otel` + `up`
bringt den Collector-Container auf `running` oder `healthy`;
OTLP/gRPC und OTLP/HTTP auf `localhost:4317`/`:4318` erreichbar.

## Akzeptanzkriterien

- ✅ `u-boot add otel` (Service noch nicht aktiv) erzeugt drei
  Mutations: Compose-Block neu, `otel-collector-config.yaml` neu,
  `services.otel.enabled: true` in `u-boot.yaml`. Idempotent:
  zweiter Aufruf → „already active; no changes" + Exit 0; kein
  Repair-Loop (Vorbedingung erfüllt durch slice-v1-keycloak T2
  `volumeOptional` + diesen Slice's `extraFiles`-Detect-
  Erweiterung).
- ✅ `u-boot add otel` ohne `u-boot.yaml` failt mit
  `ErrProjectNotInitialized` + Exit 10 (vorhandener Code-Pfad).
- ✅ `u-boot add otel --with-deps`: `dependenciesFor(otel)`
  returnt heute `[]` (Spec-Beispielkonfigurationen für App-Services
  aus LH-FA-ADD-006 §909 sind out-of-scope, siehe unten); der
  `--with-deps`-Pfad läuft normal durch.
- ✅ `u-boot remove otel` funktioniert reziprok — entfernt
  Compose-Block UND die `otel-collector-config.yaml` (extraFiles-
  Removal). `services.otel.enabled: false` bleibt im
  `u-boot.yaml`.
- ✅ E2E-Acceptance-Test `otel_acceptance_docker_test.go` analog
  zu Postgres/Keycloak — nutzt `acceptance_helpers.go` aus
  slice-v1-keycloak T3 + `dialTCP` für 4317 und 4318.
- ✅ Doctor-Smoke: `add otel` darf keine neuen Errors auslösen.
  `devcontainer.forwardPorts`-Vorschlag enthält 4317 und 4318
  (`collectActiveServicePorts` ist generisch).
- ✅ Per-Service-Catalogue-Tabelle gewinnt das neue Feld
  `extraFiles []extraFileEntry` mit `{path string; tmpl string;
  managedBlock bool}` pro Eintrag; `renderServiceTemplates`
  liefert die Liste der gerenderten Bytes mit; `executeAdd` und
  `executeRemove` respektieren sie.

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | **Templates + extraFiles-Catalogue-Erweiterung.** Drei Templates neu — `templates/services/otel.compose.tmpl` (Image `otel/opentelemetry-collector:0.108.0` — T1-Decision Stable-Pin; Ports 4317+4318; `command: --config=/etc/otel-collector-config.yaml`; read-only Bind-Mount `./otel-collector-config.yaml:/etc/otel-collector-config.yaml`), `otel.env.tmpl` (leer), `otel.config.tmpl` (Receivers `otlp/grpc+http`, Processors `batch`, Exporters `debug`, Pipelines `logs`/`metrics`/`traces`). **Kein** `otel.volume.tmpl` (analog Keycloak; Config-File liegt im Projektverzeichnis, nicht in einem Named Volume). `serviceCatalogueEntry` erweitert um `extraFiles []extraFileEntry`; `renderServiceTemplates` liefert die zusätzlichen Bytes pro Service zurück. **T1-Sub-Decision:** `extraFileEntry.managedBlock = false` für OTel-Config (die ganze Datei wird vom Slice-Generator erzeugt; User-Edits außerhalb sind ohnehin sinnlos, weil die Datei ein einziges YAML-Dokument ist). Test-Helper-Erweiterung: `RenderServiceTemplatesForTest` returnt jetzt auch ExtraFiles. Tests: Template-Existenz-Pin (Postgres unverändert, Keycloak unverändert, OTel compose+env+config), Render-Smoke (OTel rendert ohne Fehler, config.yaml ist syntaktisch valides YAML via Test-only-Roundtrip durch yaml.Unmarshal), Postgres-Byte-Identity bleibt unverändert. Catalogue für `isSupportedService` **noch nicht** erweitert (gleicher Grund wie Keycloak T1 — vor T2 ist `executeAdd`/`executeRemove` noch nicht extraFiles-aware). |
| T2 | **executeAdd + executeRemove + Catalogue-Erweiterung.** `executeAdd` schreibt die `extraFiles` nach dem Compose-Patch und nimmt sie in `Changed` auf. Plan-Phase: file-mode-Erfassung (oder default 0o644 für create-paths), Two-Phase Plan-then-Write analog F1-Fix aus slice-v1-add-remove. `executeRemove`: für jeden `extraFiles`-Eintrag im Catalogue-Eintrag des Service die Datei löschen (idempotent — `os.Remove` mit `IsNotExist`-skip), in `Changed` aufnehmen. `isSupportedService("otel")` + `supportedServices()` erweitert auf `[postgres, keycloak, otel]`. **Designentscheidung T2:** `extraFiles` werden **nicht** in der Active-Repair-Detect-Pfad integriert (`detectActiveArtifacts`) — die Datei-Existenz pro extraFile wäre eine vierte Repair-Flag, aber für OTel ist Datei-Vorhandensein ein vollständiger Re-add-Trigger (state=Unregistered/Deactivated) statt eines partiellen Repair-Pfades. Dies ist eine bewusste Vereinfachung; ein Folge-Slice kann das nachholen wenn ein Add-on mehrere extraFiles + partielle Korruption realistisch wird. Tests: 100% Coverage auf neue executeAdd/Remove-Pfade, Postgres + Keycloak Snapshot-Pin (Byte-Identity), OTel-Idempotenz-Pin (`AddTwice_NoRepairLoop` analog Keycloak T2). |
| T3 | **E2E-Acceptance.** `otel_acceptance_docker_test.go` nutzt `acceptance_helpers.go` (slice-v1-keycloak T3). `acceptanceFlow.envKeys = nil` (OTel hat keine .env-Keys). `dialTCP` für `localhost:4317` (OTLP/gRPC) UND `localhost:4318` (OTLP/HTTP). UpService-Timeout 60s (Collector-Boot < 5s, deutlich schneller als Keycloak). **CI-Flake-Carveout-Erwartung:** falls docker.io-Pull des OTel-Image in CI ebenfalls flake-anfällig ist (analog Keycloak/Quay), erweitert dieser Slice den `acceptance_extended`-build-tag-Patch UND benennt `slice-v1-keycloak-ci-flake` um zu `slice-v1-acceptance-ci-flake` (oder kreiert eigenen Folge-Slice). Diagnose passiert erst in T3-CI-Run. |
| T4 | **Closure.** READMEs (`add <service>`-Reference erwähnt jetzt `postgres \| keycloak \| otel`), CHANGELOG `## [Unreleased]` Added-Eintrag mit LH-FA-ADD-004 + LH-AK-004 + Hinweis auf `otel-collector-config.yaml` als drittes Slice-Datei-Artifact, roadmap.md §v0.3.0-Tabelle markiert `slice-v1-otel` ✅ mit T1..T3-Hashes + Stand-Bump 5/5 = **v0.3.0-Milestone vollständig**. Slice-Plan `open/` → `done/`. v0.3.0-Release-Cut-Slice danach als Folge-Slice. |

## Out of Scope

- **Beispielkonfigurationen für bestehende App-Services**
  (LH-FA-ADD-006 §909 sagt „OpenTelemetry kann
  Beispielkonfigurationen für bestehende App-Services erzeugen").
  Das ist eine OTel-spezifische Erweiterung der Add-Phase, die
  Service-spezifische Receivers (Postgres-Receiver, Keycloak-
  Receiver, …) in die `otel-collector-config.yaml` einfügen würde.
  Heute noch nicht im Catalogue; eigenes Folge-Slice
  `slice-v1-otel-app-service-receivers` bei konkretem Bedarf.
- **Healthcheck via Collector-Health-Extension**: LH-AK-004
  verlangt nur Container-Status `running` ODER `healthy` — der
  Collector-Default ohne Health-Extension läuft auf `running`,
  was ausreicht. Health-Extension + entsprechender Healthcheck
  wäre eigene Slice-Arbeit.
- **OTLP-TLS-Konfiguration**: das Default-Setup nutzt unverschlüsselte
  OTLP-Endpoints, was für lokale Dev-Compose-Stacks angemessen
  ist. TLS-Setup wäre ein Production-Hardening-Slice.
- **Andere Add-ons**: parallel-entwickelbar zu Keycloak, kein
  Cross-Slice-Abhängigkeit.

## Bezug

- Spec: `LH-FA-ADD-004` (V1, §862-§881) + `LH-AK-004` (V1, §2356-§2374).
- Voraussetzungs-Slices:
  - [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) —
    M5-Pattern.
  - [`slice-v1-add-remove`](../done/slice-v1-add-remove.md) —
    reziproker Remove-Pfad.
  - [`slice-v1-addons-deps`](../done/slice-v1-addons-deps.md) —
    Dep-Mechanik (heute inaktiv für OTel).
  - [`slice-v1-keycloak`](../done/slice-v1-keycloak.md) —
    Per-Service Probe-Mechanismus, volumeOptional-Catalogue-Feld,
    acceptance_helpers.go-Extraktion. Voraussetzung für T1 + T3.
- Folge-Slices:
  - `slice-v1-otel-app-service-receivers` — LH-FA-ADD-006 §909
    Beispielkonfigurationen für bestehende App-Services.
  - Falls CI-Flake gemeinsam: ggf. Migration von
    `slice-v1-keycloak-ci-flake` zu allgemein
    `slice-v1-acceptance-ci-flake`.
  - `slice-v1-release-cut-v0.3.0` (nach Slice-Closure) —
    Release-Cut analog `slice-v1-release-cut-v0.2.0`.
- Milestone: v0.3.0 „Add-on Catalogue Expansion" (siehe
  [roadmap.md §v0.3.0](../in-progress/roadmap.md#v030)). Fünfter
  Slice — schließt den Milestone (4/5 ✅ vor diesem Slice).
- Phase: V1 (nach v0.2.0); kein Carveout.
