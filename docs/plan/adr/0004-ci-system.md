# ADR 0004: CI-System (GitHub Actions, schlank, Docker-only)

## Status

Accepted

## Datum

2026-05-21

## Kontext

`LH-QA-003` verlangt CI-Fähigkeit. Bisher (M1–M2b) ist die Anforderung
abstrakt erfüllt (`make gates` / `make ci` laufen Docker-only auf
jedem Host mit Docker + `make`), aber kein konkretes CI-System mit
PR-blockierenden Status-Checks ist angeschlossen.

Vorlagen:

- `k-deskflight` (Go, ADR 0012 §2.11): zwei Workflows
  (`ci.yml` + `cluster-smoke.yml`), schlank, Docker-only, SHA-pinned
  Actions, top-level `permissions: {}` mit Per-Job-Lockerung.
- `m-trace`: 9 Workflows, fein granuliert (benchmark, fuzz, mutation,
  publish-images/packages, security-audit, …). Für ein Projekt im
  MVP-Bootstrap zu reich.
- `grid-gym`: noch kein `.github/workflows/` — auch in der Aufbau-
  Phase.

Lastenheft-Bezug:

- `LH-QA-003` – CI-Fähigkeit (in diesem Commit von "soll testbar" auf
  konkrete Pflichten verschärft).
- `LH-FA-BUILD-005..007` – Make-Targets und Docker-only-Workflow.
- `LH-NFA-PORT-002` – möglichst wenige Systemabhängigkeiten am Host;
  Docker-only bleibt auch im CI-Runner Pflicht.

## Entscheidung

**GitHub Actions** als CI-System, mit dem schlanken k-deskflight-
Muster:

- Eine Workflow-Datei: `.github/workflows/ci.yml`.
- Zwei parallele Jobs:
  - `gates` — `make gates` (lint + test + coverage-gate).
  - `security-gates` — `make govulncheck`.
- Beide PR-blockierend; Required-Status-Checks-Konfiguration erfolgt im
  GitHub-UI nach dem ersten grünen Lauf (kann nicht vom Workflow-File
  selbst gesetzt werden).
- Trigger: `pull_request` und `push` auf `main`.
- Runner: `ubuntu-latest` mit vorinstalliertem Docker + BuildKit.
- Actions **SHA-gepinnt** mit Tag-Kommentar
  (`uses: actions/checkout@<sha> # v6.0.2`). Hebung Routine, neuer
  Commit-SHA via `gh api repos/actions/checkout/git/refs/tags/<tag>`.
- Top-Level `permissions: {}`, jeder Job lockert auf das Minimum
  (`contents: read` für reine Build-Jobs).
- `timeout-minutes: 20` pro Job (analog k-deskflight).
- `DOCKER_BUILDKIT=1` als Job-Env (BuildKit-Cache + Multi-Stage-Cache-
  Filter aus `LH-FA-BUILD-005`).

Bewusst **noch nicht** in diesem Slice:

- Image-Publish-Workflow (`publish.yml` / GHCR) — kommt mit dem
  Release-Slice (`LH-OPEN-002` Paketierung).
- Trivy-Image-Scan — Folge-Slice analog `LH-FA-BUILD-006`-Erweiterung.
- Compose-Stack-Smoke und Adapter-Integrationstests gegen die externe
  Docker-Engine — separater Slice mit Build-Tag-getriggerten Pfaden
  (`spec/architecture.md` §5). Der Begriff „Cluster-Smoke" aus
  `k-deskflight` (Kubernetes-Smoke gegen `kind`) passt für u-boot
  nicht, weil u-boot Compose-Stacks orchestriert, nicht Kubernetes;
  der u-boot-Pendant ist die Integrationsstrecke aus
  [`slice-m6-docker-integrationstests.md`](../planning/done/slice-m6-docker-integrationstests.md).
- DCO-Bot / Branch Protection — Folgepflicht im GitHub-UI.

## Konsequenzen

Positiv:

- **Pflicht-Gates ab Tag 1 enforced**: PRs ohne grünes `gates`/
  `security-gates` werden nicht gemerged.
- **Low maintenance**: zwei Jobs, beide rufen Make-Targets — wenn sich
  die Build-Pipeline ändert (neuer Stage, neuer Aggregator), läuft der
  Workflow weiter ohne YAML-Edit.
- **Supply-Chain-Härtung**: SHA-pinned Actions verhindern den
  klassischen Tag-Move-Angriff; explizite `permissions: {}` blockt
  versehentlich neu hinzukommende Steps mit schreibendem Token.
- **Docker-only-Konsistenz** (`LH-FA-BUILD-007`): identische
  Build-/Lint-/Test-Pfade lokal und in CI; keine "läuft nur in CI"-
  Drift.

Negativ / Trade-offs:

- **Vendor-Bindung an GitHub Actions**. Kompensiert dadurch, dass die
  eigentliche Logik in `Makefile`/`Dockerfile` lebt; ein Wechsel
  (GitLab CI, Forgejo, …) bedeutet nur den Workflow-Wrapper neu zu
  schreiben.
- **CI-Runtime**: zwei sequentielle Docker-Builds pro Job
  (deps-Cache, dann lint/test/coverage). Akzeptabel für den
  Bootstrap-Stand; mit wachsender Codebase ggf. Cache-Strategien
  (BuildKit Remote Cache, GHA Cache Action) ergänzen.
- **Branch-Protection im UI** ist nicht im Repo versioniert. Wird im
  README-Setup-Abschnitt einer späteren Iteration dokumentiert.

Alternativen (verworfen):

- **m-trace-Stil (9 Workflows)**: zu reich für den aktuellen Umfang;
  Workflows entstehen aus konkretem Bedarf (Release, Fuzz, Mutation),
  nicht prospektiv.
- **Eigenes CI-System (Forgejo Actions, Drone)**: erhöht den
  Bootstrap-Aufwand ohne klaren Nutzen; u-boot lebt auf GitHub.
- **CI ohne SHA-Pinning**: deutlich entspannter, aber bekanntes
  Supply-Chain-Risiko (Tag-Move auf populäre Actions).

## Folgepunkte

- Sobald M3 (`u-boot init`) merged ist und produktive Pakete in
  `./internal/...` liegen: Coverage-Schwellwert per
  `make coverage-gate THRESHOLD=…` im Workflow setzen
  (`LH-FA-BUILD-008`).
- Image-Publish-Workflow in einem eigenen Slice (`LH-OPEN-002`).
- Trivy-Image-Scan als optionaler dritter Job, sobald das
  Runtime-Image regelmäßig gebaut wird.
- Branch-Protection-Konfiguration im README dokumentieren
  (Folge-Slice oder direkt mit dem ersten PR).
