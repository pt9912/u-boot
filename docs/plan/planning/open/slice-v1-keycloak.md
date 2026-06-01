# Slice V1: `u-boot add keycloak` (LH-FA-ADD-003 + LH-AK-003)

## Auslöser

`LH-FA-ADD-003` (V1) verlangt, dass `u-boot add keycloak` einen
Keycloak-Service mit Compose-Block, Admin-Env-Block (Placeholder-
Secrets), Port-Konfiguration und Healthcheck in ein initialisiertes
Projekt einbaut. `LH-AK-003` ist der zugehörige Acceptance-Flow
(`init` + `add keycloak` + `up` + Endpoint-Probe auf Port 8080).

Vierter Slice des v0.3.0-Milestones, baut direkt auf dem
`slice-m5-add-postgres`-Pattern auf und ist der erste echte
Add-on-Service, der von der `slice-v1-addons-deps`-Mechanik
profitiert (optionaler Postgres-Dep bei externer Persistenz, siehe
LH-FA-ADD-006-Spec-Beispiel).

## Aufhebungsbedingung

`u-boot add keycloak` in einem initialisierten Projekt liefert:

1. Keycloak-Service-Block in `compose.yaml` mit
   - Image `quay.io/keycloak/keycloak:<pinned-tag>` (kein `:latest`),
   - Port-Mapping `8080:8080`,
   - Healthcheck (HTTP-Endpoint-Probe oder `kc.sh` falls verfügbar),
   - Env-Block mit Admin-Credentials aus `.env`,
   - Default-mode embedded persistence (H2/file-basiert) — kein
     postgres-Bedarf out-of-the-box.
2. `.env.example` mit `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN` +
   `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD` als
   Placeholder-Secrets (Spec-mandatiert: niemals reale Secrets im
   Repo).
3. `services.keycloak.enabled: true` in `u-boot.yaml`.

LH-AK-003-Akzeptanz: `init demo --no-git` + `add keycloak` + `up`
auf einem Docker-fähigen System bringt Keycloak-Endpoint
`http://localhost:8080/` mit HTTP 200 oder 302 hoch.

## Akzeptanzkriterien

- ✅ `u-boot add keycloak` (Service noch nicht aktiv) erzeugt drei
  Mutations: Compose-Block neu, env-Block neu, `services.keycloak.
  enabled: true` in `u-boot.yaml`. Idempotent: zweiter Aufruf →
  „already active; no changes" + Exit 0.
- ✅ `u-boot add keycloak` ohne `u-boot.yaml` failt mit
  `ErrProjectNotInitialized` + Exit 10 (vorhandener Code-Pfad).
- ✅ `u-boot add keycloak --with-deps`: heute no-op (Keycloak
  deklariert default keine harte Postgres-Dep, weil Default-mode
  embedded ist); der `--with-deps`-Flag bleibt voll funktional und
  ändert das Verhalten erst, wenn `services.keycloak.persistence:
  external-postgres` als Schema-Feld + Dep-Trigger landet
  (eigener Folge-Slice, siehe Out-of-Scope).
- ✅ `u-boot remove keycloak` funktioniert reziprok (existierender
  RemoveServiceService-Pfad greift automatisch, sobald Keycloak in
  der Catalogue ist).
- ✅ E2E-Acceptance-Test `keycloak_acceptance_docker_test.go`
  analog zur Postgres-Acceptance: Docker-only-Test (`//go:build
  docker`-Tag), init + add + up, HTTP-Probe gegen Port 8080
  erwartet 200 oder 302.
- ✅ Doctor-Checks: `services.keycloak.enabled` taucht in
  `services.enabled`-Diagnostic auf (keine Service-Hardcoding
  nötig); `devcontainer.forwardPorts` schlägt Port 8080 vor wenn
  Keycloak aktiv ist.
- ✅ Hexagonale Verdrahtung: kein neuer Driving-Port nötig (Reuse
  AddServiceUseCase + RemoveServiceUseCase); Catalogue-Erweiterung
  + Template-Render-Refactor sind die einzigen Code-Änderungen
  außerhalb der Test- und Template-Dateien.

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | Catalogue erweitern (`isSupportedService("keycloak")` + `supportedServices()` returns `["postgres", "keycloak"]`); drei Templates anlegen — `templates/services/keycloak.compose.tmpl` (Image-Pin, Port-Mapping, Healthcheck, Env-Block, embedded-persistence-default), `keycloak.env.tmpl` (KEYCLOAK_ADMIN + KEYCLOAK_ADMIN_PASSWORD Placeholder-Secrets), `keycloak.volume.tmpl` (heute leer/keine Volume; Datei-Existenz aus Template-Loader-Symmetrie). `embed.FS`-Eintrag für Template-Resolution erweitern. Tests: Catalogue-Pin (`TestSupportedServices`), Template-Existenz-Pin (`TestEmbeddedTemplatesContainAllServiceFiles`), Render-Smoke (Templates rendern fehlerfrei). |
| T2 | `renderPostgresTemplates` zu `renderServiceTemplates(svc)` generalisieren: Template-Pfad-Resolution per `services/<svc>.compose.tmpl` etc., Empty-Volume-Behandlung wenn `<svc>.volume.tmpl` leer ist. `detectActiveArtifacts` ebenfalls generisch oder per-Service-Probe (Switch über `svc.String()`). Postgres-Probe (`POSTGRES_USER`) bleibt unverändert; Keycloak-Probe checkt `KEYCLOAK_ADMIN`-Env-Var + `8080`-Port-Mapping. Tests: 100% Coverage auf den Generalisierungs-Pfad, Postgres-Pfad bleibt grün (no-behavior-change-Pin via Snapshot). |
| T3 | E2E-Acceptance-Test `keycloak_acceptance_docker_test.go` analog `postgres_acceptance_docker_test.go` (`//go:build docker`-Tag, MakefileSL test-docker-Target deckt ihn automatisch ab). Init + add + up, danach HTTP-GET auf `http://localhost:8080/` erwartet 200 oder 302. Doctor-Smoke: `add keycloak` darf doctor nicht ins Rote treiben (keine neuen Errors); `devcontainer.forwardPorts`-Vorschlag enthält 8080. Manuelle Smoke-Anleitung im CHANGELOG-Eintrag (für Repo-Owner-Spot-Check). |
| T4 | Closure: READMEs (`add <service>`-Subcommand-Reference erwähnt jetzt `postgres \| keycloak`), CHANGELOG `## [Unreleased]` Added-Eintrag mit LH-FA-ADD-003 + LH-AK-003-Bezug, roadmap.md §v0.3.0-Tabelle markiert `slice-v1-keycloak` ✅ mit T1..T3-Hashes + Stand-Bump 4/5, Slice-Plan `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. |

## Out of Scope

- **`services.keycloak.persistence`-Schema-Feld + externe Postgres-
  Dep-Auflösung**: Spec (LH-FA-ADD-003) sagt „optionale PostgreSQL-
  Anbindung bei konfigurierter persistenter externer Datenbank" —
  dafür braucht es ein neues Schema-Feld in `u-boot.yaml`
  (`services.keycloak.persistence: external-postgres`) + den
  Dep-Trigger in `dependenciesFor("keycloak")`. Heute hat das
  `LH-FA-CONF-005`-Schema dafür keinen Slot; der Schema-Slot + die
  CLI-Setterkonventionen sind eigene Slice-Arbeit. Plan-Anker:
  Folge-Slice `slice-v1-keycloak-external-postgres` (Trigger:
  Nutzer-Bedarf nach persistenter Keycloak-Datenbank).
- **Custom-Realm-Import** / Theme-Konfiguration: nicht in
  LH-FA-ADD-003 Mindestumfang. Nutzer kann nach `add keycloak`
  manuell zusätzliche Volumes / Env-Vars im managed-block-freien
  Bereich der compose.yaml ergänzen.
- **OpenTelemetry-Add-on**: eigener Slice `slice-v1-otel`
  (LH-FA-ADD-004 + LH-AK-004), parallel-entwickelbar weil keine
  Dep zwischen Keycloak und OTel.
- **Keycloak-Healthcheck-Spec-Recherche**: konkrete
  Healthcheck-Implementierung (`/health/ready`-Endpoint vs.
  `kc.sh show-config` vs. einfacher `curl localhost:9000/health`)
  wird in T1 entschieden, basierend auf der gewählten Keycloak-
  Image-Version. Falls keine zuverlässige Healthcheck-Lösung im
  Slice-Zeitfenster ergibt → Service ohne Healthcheck akzeptabel
  (LH-FA-ADD-003 sagt „soweit technisch sinnvoll").

## Bezug

- Spec: `LH-FA-ADD-003` (V1) + `LH-AK-003` (V1).
- Voraussetzungs-Slices: keine harten (Catalogue + Template-
  Refactor sind in diesem Slice enthalten); weiche Voraussetzungen:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  liefert das Pattern, [`slice-v1-add-remove`](../done/slice-v1-add-remove.md)
  liefert den reziproken Remove-Pfad, [`slice-v1-addons-deps`](../done/slice-v1-addons-deps.md)
  liefert die Dep-Mechanik (heute hier nur Reserve-Kapazität, ohne
  aktiven Trigger).
- Folge-Slices:
  - `slice-v1-otel` (LH-FA-ADD-004 + LH-AK-004) — parallel
    entwickelbar.
  - `slice-v1-keycloak-external-postgres` — eigener Trigger-Slice
    für die Schema-Erweiterung + Dep-Aktivierung (siehe
    Out-of-Scope).
- Milestone: v0.3.0 „Add-on Catalogue Expansion" (siehe
  [roadmap.md §v0.3.0](../in-progress/roadmap.md#v030)). Vierter
  Slice des Milestones (3/5 ✅ vor diesem Slice).
- Phase: V1 (nach v0.2.0); kein Carveout.
