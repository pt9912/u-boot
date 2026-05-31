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
- Drei parallele Jobs (image-scan via `slice-v1-release-pipeline`
  T3 ergänzt, Distributionsentscheidung in
  [ADR-0007](0007-distributionswege-ghcr.md)):
  - `gates (lint + test + coverage-gate)` — `make gates`.
  - `security-gates (govulncheck)` — `make govulncheck`.
  - `image-scan (trivy HIGH+CRITICAL)` — `make image-scan` (lokale
    Reproduktion); CI nutzt `aquasecurity/trivy-action` mit
    identischem Severity-Profil.
- Alle drei PR-blockierend; Required-Status-Checks-Konfiguration
  erfolgt im GitHub-UI nach dem ersten grünen Lauf (kann nicht vom
  Workflow-File selbst gesetzt werden); die Required-Status-Check-
  Liste muss die verbose `name:`-Felder verwenden, nicht die kurzen
  `jobs.<key>`-Identifier — siehe
  [`docs/user/branch-protection.md`](../../user/branch-protection.md).
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

Bewusst **noch nicht** in diesem Slice (Stand 2026-05-31 mit
Verweisen auf die Nachfolge-Slices ergänzt):

- Image-Publish-Workflow (`publish.yml` / GHCR) — geliefert mit
  [`slice-v1-release-pipeline`](../planning/done/slice-v1-release-pipeline.md)
  T2 (`93b703e`), Distributionsentscheidung in
  [ADR-0007](0007-distributionswege-ghcr.md).
- Trivy-Image-Scan — geliefert mit `slice-v1-release-pipeline` T3
  (`8212889`) als dritter PR-blockierender CI-Job
  (`ci.yml::image-scan`).
- Compose-Stack-Smoke und Adapter-Integrationstests gegen die externe
  Docker-Engine — geliefert mit
  [`slice-m6-docker-integrationstests`](../planning/done/slice-m6-docker-integrationstests.md)
  als eigenständiger Workflow `.github/workflows/integration.yml`.
  Der Begriff „Cluster-Smoke" aus `k-deskflight` (Kubernetes-Smoke
  gegen `kind`) passt für u-boot nicht, weil u-boot Compose-Stacks
  orchestriert, nicht Kubernetes.
- DCO-Bot / Branch Protection — Branch-Protection-Checkliste in
  [`docs/user/branch-protection.md`](../../user/branch-protection.md)
  publiziert (`slice-v1-release-pipeline` Teilabschluss 2026-05-27);
  DCO-Bot bleibt Out of Scope ohne eigenen Trigger.

## Konsequenzen

Positiv:

- **Pflicht-Gates ab Tag 1 enforced**: PRs ohne grünes `gates` /
  `security-gates` / `image-scan` (letzteres seit
  `slice-v1-release-pipeline` T3) werden nicht gemerged, sobald die
  Branch-Protection-Required-Status-Checks aktiviert sind
  (siehe [`docs/user/branch-protection.md`](../../user/branch-protection.md)).
- **Low maintenance**: drei Jobs, alle starten von Make-Targets aus —
  `gates`/`security-gates` delegieren komplett (`make gates`,
  `make govulncheck`), `image-scan` läuft `make build` plus die
  `aquasecurity/trivy-action` mit eigenständigem `trivy-version`-Pin.
  Build-Pipeline-Änderungen lassen den YAML-Workflow weitgehend
  unberührt; Trivy-Versions-Bumps berühren ZWEI Pin-Stellen
  (`Makefile::TRIVY_VERSION` + `ci.yml::trivy-version`, beide auf
  derselben Trivy-Version in unterschiedlicher Schreibweise —
  Detail-Kommentare an beiden Pin-Stellen sowie in
  [`docs/user/quality.md`](../../user/quality.md) §4).
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
- **CI-Runtime**: `gates` läuft zwei sequentielle Docker-Builds
  (deps-Cache, dann lint/test/coverage). `security-gates` läuft
  einen einzelnen ephemeren `docker run golang:$(GO_VERSION) ...`
  mit `go install` + `govulncheck` (kein `docker build`).
  `image-scan` baut das Runtime-Image (`make build`) und scannt es
  mit Trivy. Akzeptabel für den Bootstrap-Stand; mit wachsender
  Codebase ggf. Cache-Strategien (BuildKit Remote Cache, GHA Cache
  Action) ergänzen.
- **Branch-Protection im UI** ist nicht im Repo versioniert. Schritt-
  für-Schritt-Aktivierung dokumentiert in
  [`docs/user/branch-protection.md`](../../user/branch-protection.md)
  (`slice-v1-release-pipeline` Teilabschluss 2026-05-27).

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
- ~~Image-Publish-Workflow in einem eigenen Slice (`LH-OPEN-002`)~~
  — geliefert mit `slice-v1-release-pipeline` T2 (`93b703e`,
  Distributionsentscheidung in ADR-0007).
- ~~Trivy-Image-Scan als optionaler dritter Job~~ — geliefert mit
  `slice-v1-release-pipeline` T3 (`8212889`) als PR-blockierender
  dritter CI-Job.
- ~~Branch-Protection-Konfiguration im README dokumentieren~~ —
  veröffentlicht in `docs/user/branch-protection.md`
  (`slice-v1-release-pipeline` Teilabschluss 2026-05-27).
