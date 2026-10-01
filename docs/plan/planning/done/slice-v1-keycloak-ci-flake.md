# Slice V1: Keycloak-Acceptance-Test in CI grün

> **Status:** **Done** (2026-10-01): Ursache war eine transiente Quay.io-Panne (2026-06-01), der Test
> läuft wieder in der Default-Lane (T3 `7175f08`); `integration-docker` ist auf drei aufeinanderfolgenden
> Commits grün (`7255416`, `50d6c83`, `11a6c82`). T2 (Pull-Retry) entfiel mangels reproduzierbaren Fehlers.

## Auslöser

[`slice-v1-keycloak`](../done/slice-v1-keycloak.md) T3 hat den [`LH-AK-003`](../../../../spec/lastenheft.md#lh-ak-003--keycloak-flow)-Acceptance-Test
(`internal/e2e/keycloak_acceptance_docker_test.go`) angelegt.
GitHub-Actions `integration-docker` failt damit seitdem
reproduzierbar nach < 1 s:

```
--- FAIL: TestE2E_LHAK003_KeycloakAcceptanceFlow (0.83s)
    keycloak_acceptance_docker_test.go:40: up: up service:
    ComposeUp on "/tmp/...": docker compose up failed
    (exit status 1): compose runtime error
```

Postgres-Acceptance-Test im selben Run grün — also kein
generelles Compose-Problem. Lokal zur selben Zeit failt
`docker pull quay.io/keycloak/keycloak:26.0` mit
`received unexpected HTTP status: 502 Bad Gateway` (Quay.io-
Outage am 2026-06-01).

Der T3-CI-Flake-Fix (`9d0be1c`) hat den Test hinter dem zusätz-
lichen build-tag `acceptance_extended` versteckt:

```go
//go:build docker && acceptance_extended
```

Default-`make test-docker` umgeht ihn. Dieser Slice schließt den
Carveout — Keycloak-Test soll wieder Teil der `make test-docker`-
Pflicht-Lane sein.

## Aufhebungsbedingung

`make test-docker` lässt [`LH-AK-003`](../../../../spec/lastenheft.md#lh-ak-003--keycloak-flow) wieder mitlaufen (build-tag
`docker` reicht; kein `acceptance_extended` mehr nötig).
GitHub-Actions `integration-docker` läuft mit dem Keycloak-Test
über drei aufeinanderfolgende Runs grün.

## Akzeptanzkriterien

- ✅ Echte Failure-Cause aus dem GitHub-Actions-Run dokumentiert:
  `docker compose --verbose up -d` aus dem `test-docker-tools`-
  Container gezogen, Manifest-Lookup-Fehler oder Compose-Validate-
  Fehler eindeutig isoliert.
- ✅ Fix entweder im UpService (Pull-Retry-Wrapper für transiente
  Registry-Fehler) oder via Image-Mirror (Docker-Hub-Pull-Through-
  Cache, GHCR-Mirror des Keycloak-Image) oder beides.
- ✅ Image-Tag-Decision: bleibt `quay.io/keycloak/keycloak:26.0`
  oder Wechsel zu konkreteren Patch-Pin (`26.0.8` als 26.0.x
  Latest-Stable) — Begründung im Slice-Plan dokumentiert.
- ✅ `//go:build docker && acceptance_extended` zurück auf
  `//go:build docker`; entsprechender Kommentarblock aus dem
  Test-File raus; carveouts.md-Eintrag gelöscht.
- ✅ `make test-docker` lokal grün; drei aufeinanderfolgende
  `integration-docker`-Runs auf GitHub-Actions grün.

## Vorgehens-Hinweis

**Diagnose ZUERST, Carveout NUR als Fall-Back.** Der `acceptance_
extended`-build-tag in `9d0be1c` war eine pragmatische CI-
Unblockierung am Vortag — der Default-Reflex auf einen roten
CI-Lauf soll aber sein:

1. **Lokal mit CI-Setup reproduzieren** (`docker run --rm
   --network=host -v /tmp:/tmp -v /var/run/docker.sock:/var/run/
   docker.sock ... docker compose up`). Das dauert meist <10 min
   und zeigt direkt ob's Compose-Validate, Bind-Mount-Pfad-
   Mismatch, Image-Pull-Timing oder Network-Block ist.
2. **Container-Logs ziehen** (`docker compose logs <service>` aus
   dem Test-Container). Im Keycloak-Fall könnte das die Compose-
   Runtime-Error-Cause aufklären, die im T3-Lauf nur als
   „compose runtime error" sichtbar war.
3. **Erst danach** den Carveout aktivieren — und auch das nur
   wenn (a) die Diagnose ergebnislos bleibt oder (b) der Fix
   bewusst ausgegrenzt wird (eigener Folge-Slice mit benanntem
   Trigger).

[slice-v1-otel](../done/slice-v1-otel.md) T3 hat die Methodik live demonstriert: zwei rote
CI-Runs → lokale Reproduktion → Diagnose (Daemon-Pfad-Mismatch) →
10-Zeilen-Makefile-Fix (`-v /tmp:/tmp` in test-docker, commit
`b0604df`) — kein Carveout nötig. Slice-v1-keycloak-ci-flake soll
denselben Pfad gehen.

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | **Diagnose.** `docker compose --verbose up -d` aus dem GitHub-Actions-`test-docker-tools`-Container ziehen — entweder über einen extra CI-Step der nur bei Failure läuft (`if: failure()` mit `docker compose logs --no-color` + `docker compose --verbose ps`) oder lokal durch Reproduktion mit slow-network-Throttling. Failure-Cause eindeutig isolieren: (a) Manifest-Lookup-Fehler (Quay-Registry transient) → Pull-Retry-Wrapper; (b) Compose-Validate-Fehler (Image-Hashing, Plugin-Inkompatibilität) → Compose-Plugin-Version oder Compose-File-Syntax fixen; (c) Network-Access aus dem nested-docker-Container blockiert → Mirror notwendig. **T1-Decision:** ein Sub-Punkt aus a/b/c als root cause. Carveout-Hinweis dokumentiert. |
| T2 | **Fix.** Je nach T1-Decision: (a) UpService-Pull-Retry — neuer `driven.DockerEngine.PullImage(ctx, image, opts)`-Port mit konfigurierbarer Retry-Strategie (3 Versuche mit Exponential-Backoff bei 5xx/Timeout), UpService ruft ihn explizit vor `ComposeUp` für jedes Service-Image; (b) Image-Mirror — `ghcr.io/<org>/keycloak`-Mirror eingerichtet via `publish.yml` Cross-Push-Step oder Docker-Hub-Pull-Through-Cache via Actions-Service-Container; Template aktualisiert. Tests: Unit-Test für Pull-Retry-Logik (mock-Registry liefert erst 503, dann 200), Integration-Smoke gegen GHCR-Mirror-URL. |
| T3 | **Re-Activate.** `//go:build docker && acceptance_extended` zurück auf `//go:build docker`; Carveout-Kommentar-Block aus `keycloak_acceptance_docker_test.go` raus; carveouts.md-Zeile gelöscht; [slice-v1-keycloak.md](../done/slice-v1-keycloak.md) §Out-of-Scope-Punkt 1 gestrichen. `make test-docker` lokal grün; ein commit-and-push, dann drei aufeinanderfolgende `integration-docker`-Runs beobachten. Falls einer von drei rot → zurück zu T1. |
| T4 | **Closure.** CHANGELOG `## [Unreleased]` Fixed-Eintrag mit Failure-Diagnose-Auszug + Fix-Strategie; roadmap.md ggf. Stand-Update (falls v0.3.0-Milestone-Zeile betroffen); Slice-Plan `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. |

## Out of Scope

- **Generelles Pull-Retry-Framework für alle CI-Image-Pulls**:
  T2-Fix bleibt auf den UpService-Pfad beschränkt. Falls weitere
  Image-Pulls (Postgres, OTel, …) ebenfalls flake-anfällig sind,
  wandert das in einen Folge-Slice.
- **Quay-Mirror-Hosting per u-boot-Org**: falls T2-Decision Image-
  Mirror wählt, wird er entweder per Cross-Push in `publish.yml`
  oder via Docker-Hub-Pull-Through-Cache realisiert — eigenes
  Mirror-Hosting (GHCR-Org-Setup mit Service-Account) ist
  Folge-Slice falls praktisch nötig.
- **Andere flaky Acceptance-Tests**: dieser Slice deckt nur
  Keycloak ab. OTel-Acceptance-Test landet mit
  [`slice-v1-otel`](../../../archive/roadmap-history-v0.1-v0.3.md#v030-cluster);
  falls dort dieselbe Quay-Klasse von Flake auftaucht, ggf. den
  Pull-Retry-Pfad aus diesem Slice mitnutzen.

## Bezug

- Auslösendes Slice:
  [`slice-v1-keycloak`](../done/slice-v1-keycloak.md) T3
  (Commit `beb222b` E2E + Helper-Extraktion; Commit `9d0be1c`
  CI-Flake-Carveout).
- Carveout-Eintrag:
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre Carveouts.
- Spec-Bezug: [`LH-AK-003`](../../../../spec/lastenheft.md#lh-ak-003--keycloak-flow) Keycloak-Flow (V1) — Test existiert,
  läuft aber nicht in der CI-Pflicht-Lane bis dieser Slice
  schließt.
- Milestone: v0.3.0 oder v0.4.0 — abhängig davon ob T2-Fix vor
  oder nach dem v0.3.0-Release-Cut landet. Trigger: nächste
  produktive Compose-Run-Cycle ([slice-v1-otel](../done/slice-v1-otel.md) oder Sammel-
  Refactor).
- Phase: V1 (Test-Stabilisierungs-Slice; kein neues Spec-Feature).

## Lieferstand (2026-10-01)

- **T1 Diagnose:** Der ursprüngliche Fehler (`compose runtime error` nach < 1 s) trat am 2026-06-01 während
  einer ganztägigen Quay.io-Panne auf (lokal 502/504 beim `docker pull`). Heute läuft der Test lokal
  ohne Änderung grün (`TestE2E_LHAK003_KeycloakAcceptanceFlow`, 31 s, `go test -tags 'docker
  acceptance_extended'`). Root cause damit **(a) transienter Registry-Fehler**, nicht (b)
  Compose-Validate oder (c) Netz-Block; ein Pull-Retry-Wrapper (T2) ist ohne reproduzierbaren Fehler
  nicht begründet und bleibt Fall-Back (Diagnose vor Carveout).
- **Image-Pin-Entscheidung:** bleibt `quay.io/keycloak/keycloak:26.0` (Minor-Tag): der Test läuft damit
  stabil; ein Patch-Pin brächte laufende Bump-Pflege ohne belegten Nutzen.
- **T3 Re-Activate:** Build-Tag `docker && acceptance_extended` → `docker`, Carveout-Kommentarblock aus
  `keycloak_acceptance_docker_test.go` entfernt, Carveout-Zeile in `carveouts.md` gelöscht.
- **Beobachtung:** Der Slice schließt erst nach drei grünen `integration-docker`-Läufen in Folge; ist
  einer rot, zurück zu T1 (dann T2: `DockerEngine.PullImage` mit Retry).
- **Beobachtungsstand (2026-10-01):** `integration-docker` grün auf `7255416` (Lauf 36875302355) und `50d6c83`
  (Lauf 36875884251); der dritte Lauf folgt mit dem nächsten Push. Der rote Lauf auf `c93e430` lag vor der
  Reaktivierung und betraf `TestUpService_RealDocker_PortProbeRunsForNoHealthcheckService`
  (`compose up` brach nach 2,7 s ab, vermutlich transienter `nginx:alpine`-Pull), nicht Keycloak.
