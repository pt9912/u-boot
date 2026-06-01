# Slice V1: Keycloak-Acceptance-Test in CI grün

## Auslöser

`slice-v1-keycloak` T3 hat den LH-AK-003-Acceptance-Test
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

`make test-docker` lässt LH-AK-003 wieder mitlaufen (build-tag
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

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | **Diagnose.** `docker compose --verbose up -d` aus dem GitHub-Actions-`test-docker-tools`-Container ziehen — entweder über einen extra CI-Step der nur bei Failure läuft (`if: failure()` mit `docker compose logs --no-color` + `docker compose --verbose ps`) oder lokal durch Reproduktion mit slow-network-Throttling. Failure-Cause eindeutig isolieren: (a) Manifest-Lookup-Fehler (Quay-Registry transient) → Pull-Retry-Wrapper; (b) Compose-Validate-Fehler (Image-Hashing, Plugin-Inkompatibilität) → Compose-Plugin-Version oder Compose-File-Syntax fixen; (c) Network-Access aus dem nested-docker-Container blockiert → Mirror notwendig. **T1-Decision:** ein Sub-Punkt aus a/b/c als root cause. Carveout-Hinweis dokumentiert. |
| T2 | **Fix.** Je nach T1-Decision: (a) UpService-Pull-Retry — neuer `driven.DockerEngine.PullImage(ctx, image, opts)`-Port mit konfigurierbarer Retry-Strategie (3 Versuche mit Exponential-Backoff bei 5xx/Timeout), UpService ruft ihn explizit vor `ComposeUp` für jedes Service-Image; (b) Image-Mirror — `ghcr.io/<org>/keycloak`-Mirror eingerichtet via `publish.yml` Cross-Push-Step oder Docker-Hub-Pull-Through-Cache via Actions-Service-Container; Template aktualisiert. Tests: Unit-Test für Pull-Retry-Logik (mock-Registry liefert erst 503, dann 200), Integration-Smoke gegen GHCR-Mirror-URL. |
| T3 | **Re-Activate.** `//go:build docker && acceptance_extended` zurück auf `//go:build docker`; Carveout-Kommentar-Block aus `keycloak_acceptance_docker_test.go` raus; carveouts.md-Zeile gelöscht; slice-v1-keycloak.md §Out-of-Scope-Punkt 1 gestrichen. `make test-docker` lokal grün; ein commit-and-push, dann drei aufeinanderfolgende `integration-docker`-Runs beobachten. Falls einer von drei rot → zurück zu T1. |
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
  [`slice-v1-otel`](../in-progress/roadmap.md#v030--milestone-tabelle-add-on-catalogue-expansion);
  falls dort dieselbe Quay-Klasse von Flake auftaucht, ggf. den
  Pull-Retry-Pfad aus diesem Slice mitnutzen.

## Bezug

- Auslösendes Slice:
  [`slice-v1-keycloak`](../done/slice-v1-keycloak.md) T3
  (Commit `beb222b` E2E + Helper-Extraktion; Commit `9d0be1c`
  CI-Flake-Carveout).
- Carveout-Eintrag:
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre Carveouts.
- Spec-Bezug: `LH-AK-003` Keycloak-Flow (V1) — Test existiert,
  läuft aber nicht in der CI-Pflicht-Lane bis dieser Slice
  schließt.
- Milestone: v0.3.0 oder v0.4.0 — abhängig davon ob T2-Fix vor
  oder nach dem v0.3.0-Release-Cut landet. Trigger: nächste
  produktive Compose-Run-Cycle (slice-v1-otel oder Sammel-
  Refactor).
- Phase: V1 (Test-Stabilisierungs-Slice; kein neues Spec-Feature).
