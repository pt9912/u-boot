# u-boot Roadmap

Übergreifendes Master-Dokument zum Stand aller Slices und Tranchen
(`LH-FA-PROJDOCS-003`). Wird laufend gepflegt und liegt deshalb dauerhaft
in `in-progress/`.

| Phase | Status | Beschreibung | Artefakt |
| ----- | ------ | ------------ | -------- |
| M0 Spec | Done | Lastenheft v0.1.0 (Sektionen 1–14, inkl. 4.11 Build-/CI-Infrastruktur, 4.12 Doku-Struktur) | [`spec/lastenheft.md`](../../../../spec/lastenheft.md) |
| M0 ADRs | Done | ADR-0001 Implementierungssprache Go | [`docs/plan/adr/0001-implementierungssprache-go.md`](../../adr/0001-implementierungssprache-go.md) |
| M1 Repo-Skeleton | Done | Multi-Stage Dockerfile, Makefile, .dockerignore, Repo-Layout (`LH-FA-BUILD-001..009`), Doku-Struktur (`LH-FA-PROJDOCS-001..003`), `u-boot --help` / `--version`-Stub | [`slice-m1-repo-skeleton`](../done/slice-m1-repo-skeleton.md) |
| M2 Architektur | Done | Hexagonale Architektur (`LH-FA-ARCH-001..003`), `spec/architecture.md`, ADR-0002, `internal/{hexagon,adapter}/`-Skeleton, depguard mit aktiven Schicht-Regeln (match nichts, bis erste Pakete in `./internal/...`) | [`slice-m2-hexagonale-architektur`](../done/slice-m2-hexagonale-architektur.md) |
| M2b SOLID-Lint | Done | SOLID-nahes Lint-Profil (`LH-QA-004` auf MVP gehoben), 5 Default-Linter + 24 SOLID-nahe Linter (inkl. `depguard`), `docs/user/quality.md`, ADR-0003 | [`slice-m2b-solid-lint-profil`](../done/slice-m2b-solid-lint-profil.md) |
| M2c CI | Done | GitHub-Actions-CI (`LH-QA-003` auf konkret gehoben), `.github/workflows/ci.yml` mit anfangs zwei PR-blockierenden Jobs (`gates` + `security-gates`), SHA-pinned Actions, Docker-only, ADR-0004. **Stand:** drei PR-blockierende Jobs — `image-scan (trivy HIGH+CRITICAL)` mit [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md) T3 (`8212889`) als dritter Job ergänzt; LH-QA-003-Pflichtmenge ist heute drei Jobs (siehe `LH-QA-003` in [`spec/lastenheft.md`](../../../../spec/lastenheft.md) + ADR-0004 §Entscheidung). | [`slice-m2c-ci-pipeline`](../done/slice-m2c-ci-pipeline.md) |
| M2d Carveouts | Done | Carveout-Disziplin (`LH-FA-PROJDOCS-005` MVP-Pflicht), Master-Inventar [`carveouts.md`](carveouts.md), 7 neue Slice-Pläne in [`open/`](../open/) für offene Carveouts; permanente Carveouts dokumentiert | [`slice-m2d-carveout-disziplin`](../done/slice-m2d-carveout-disziplin.md) |
| M3 `u-boot init` | Done | Projektstruktur erzeugen (`LH-FA-INIT-001..007`), `u-boot.yaml` schreiben, Git-Init, Re-Init mit `--force`/`--backup` (LH-FA-INIT-005) + Modi-Flags (LH-FA-CLI-005A). Coverage-, depguard- und gomodguard-Carveouts aufgelöst. Detail: [`slice-m3-init-flow.md`](../done/slice-m3-init-flow.md). **Stand:** T1..T4c ✅ (Commits siehe Slice-DoD); T5 ✅ `scripts/verify-depguard.sh` + `make verify-depguard`; M3-followup: [`slice-m3-build-polish`](../done/slice-m3-build-polish.md) (`987c164`, govulncheck-Pin + PROGRESS_FLAG) und [`slice-m3-gomodguard-rules`](../done/slice-m3-gomodguard-rules.md) (`201fb4b`, 4 Block-Regeln + golangci-lint v2.12.2) |
| M4 `u-boot doctor` | Done | Lokale Voraussetzungen prüfen (`LH-FA-DIAG-001..004`), Severity-Klassifikation, Repair-Hints. 9 Checks: write-permissions, git, docker (+reachable+compose-plugin), u-boot.yaml, compose.yaml, devcontainer.json/Dockerfile. CLI `doctor`-Subkommando mit `--strict`. Exit-Code 11 bei Errors (oder Warns + --strict). | [`slice-m4-doctor`](../done/slice-m4-doctor.md) |
| M5 `u-boot add postgres` | Done | PostgreSQL-Add-on (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`), services-Schema in u-boot.yaml, Compose split-block scaffold, `.env.example`-Block, Healthcheck mit `$${POSTGRES_USER:-postgres}`-Defaults für LH-AK-002, State-Machine, Active-Repair, CLI-Subcommand, doctor-Integration (services.enabled-key + devcontainer.forwardPorts + devcontainer.enabled-Severity-Eskalation). 11 doctor-Checks gesamt. | [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) |
| M6 `u-boot up` / `down` | Done | Compose-Wrapper (`LH-FA-UP-001..004`), Healthcheck-Polling, `--timeout`, `--volumes`, CLI-Subcommands + Status-Tabelle. Alle 7 Tranchen in `done/`. Carveout-Slice [`slice-m6-docker-integrationstests`](../done/slice-m6-docker-integrationstests.md) **Done** (Sub-T1..T4 + Audit-Härtung `41cab1b` + Stabilisierung `43b42e4`/`8865ca1` + Carveout-Entfernung). | [`slice-m6-up-down`](../done/slice-m6-up-down.md) |
| M7 `u-boot generate` | Done | `generate changelog`/`readme`/`env-example`/`devcontainer` (`LH-FA-GEN-001..005` + LH-FA-DEV-001/004/005 + LH-AK-007). Sechs Tranchen ✅ (`67fc181`/`3c5de48`/`037ab00`/`19c4110`/`294e492`/`d32a733`) + Review-Followup `27de9c5` (9 Findings S1..S4/N1..N5 adressiert, u. a. fenced-code-block-Schutz gegen Markdown-Korruption + CRLF-Normalisierung). Reuse von `managedblock` (3 Marker-Stile decken 4 Datei-Mappings), `generateManagedFile`-Helper für env-example/readme, atomarer Two-File-Plan für devcontainer, konservative User-Edit-Erkennung für changelog. CLI: `u-boot generate <artifact>`, Exit-Codes 0/2/10/14. | [`slice-m7-generate`](../done/slice-m7-generate.md) |
| M8 `u-boot config` | Done | `config get`/`set`/Anzeigen (`LH-FA-CONF-001..005`), Schema-Validierung. Fünf Tranchen ✅ (`f531e7e`/`d3fa294`/`23952b2`/`fbf3778`/`25cb123`): Whitelist, Port + Skeleton, Get + Show, Set mit Two-Stage-Validation, CLI-Subcommand. **Letzter MVP-blockierender Slice → MVP vollständig.** | [`slice-m8-config`](../done/slice-m8-config.md) |
| MVP-Closure | Done | Devcontainer-Mindestumfang (`LH-FA-DEV-001..005`), MVP-Acceptance-Flows (`LH-AK-001..002`, `LH-AK-005..007`). Drei Tranchen — T1 ✅ `bfe6416` `u-boot init --devcontainer` (LH-AK-005), T2 ✅ `8525c4c` LH-AK-001/-006-Pins in `acceptance_test.go` inkl. Doctor-Severity-Fix `compose.yaml.valid` (Error → Warn), T3 ✅ Slice-Closure + MVP-Bilanz. Alle 5 MVP-`LH-AK-*` gepinnt; alle MVP-`LH-FA-DEV-*` ausgeliefert. (M8 `u-boot config` ist die zweite MVP-Schließung, siehe nächste Zeile.) | [`slice-mvp-closure`](../done/slice-mvp-closure.md) |
| V1 Add-on Expansion (v0.3.0-Milestone) | In Progress (1/5) | Fünf Slices als v0.3.0-Scope: `slice-v1-audit-done` (Doku-Audit `LH-FA-BUILD-006`/`LH-NFA-MAINT-004`/`LH-NFA-PORT-003`), ✅ [`slice-v1-add-remove`](../done/slice-v1-add-remove.md) (`LH-FA-ADD-007`) geliefert 2026-06-01, `slice-v1-addons-deps` (`LH-FA-ADD-006`), `slice-v1-keycloak` (`LH-FA-ADD-003` + `LH-AK-003`), `slice-v1-otel` (`LH-FA-ADD-004` + `LH-AK-004`). Reihenfolge: audit-done → add-remove → addons-deps → Keycloak + OTel (parallel). | partial — siehe done/ |
| V1 Templates | Partial Done | `LH-FA-TPL-001..004` — Format via [ADR-0009](../../adr/0009-template-format-yaml-files.md) (YAML+`text/template`). `slice-v1-template-list` ✅ + `slice-v1-template-init` ✅ (für `basic`; Variable-Resolution defer-pflichtig); `slice-later-local-templates` (`LH-FA-TPL-003`) bleibt Later-Phase. | [`slice-v1-template-list`](../done/slice-v1-template-list.md) + [`slice-v1-template-init`](../done/slice-v1-template-init.md) |
| V1 Logs / Dry-Run / Diff | Open | `LH-FA-UP-005`, `LH-FA-CLI-007/008` | offen |
| Later Migration / Custom Templates | Open | `LH-FA-CONF-006`, `LH-FA-TPL-003` (in ADR-0009 §Folgepunkte als `slice-later-local-templates` benannt), `LH-DA-004` | offen |

## Carveout-Auflösungs-Slices

Slices, die ausschließlich offene Carveouts (`LH-FA-PROJDOCS-005`)
auflösen. Verbindlich verankert hier *und* in [`carveouts.md`](carveouts.md);
ein Carveout ohne Eintrag in beiden Quellen ist ein
Disziplin-Verstoß.

| Slice | Auslöser | Phase | Status |
| ----- | -------- | ----- | ------ |
| [`slice-m3-init-flow`](../done/slice-m3-init-flow.md) | `LH-FA-INIT-*` initialer Flow + zwei M3-Carveouts (Coverage ✅, depguard ✅) | M3 | Done |
| [`slice-m3-depguard-aktivierung-verifizieren`](../done/slice-m3-depguard-aktivierung-verifizieren.md) | `LH-FA-ARCH-003` depguard-Regeln matchen bisher nichts | M3-T5 | Done |
| [`slice-m3-gomodguard-rules`](../done/slice-m3-gomodguard-rules.md) | `gomodguard_v2.blocked: {}` leer; yaml.v3 schon drin, Cobra kommt mit T3 | M3-followup | Done |
| [`slice-m3-retroaktive-slice-plaene`](../done/slice-m3-retroaktive-slice-plaene.md) | Bootstrap-Slices (M1/M2/M2b/M2c/M2d) liegen nicht in `done/` | Done | Done |
| [`slice-m4-soft-existing-detection`](../done/slice-m4-soft-existing-detection.md) | `LH-FA-INIT-004` Soft-Erkennung + `--assume-existing` | M4-vorgezogen | Done |
| [`slice-m4-logging-port`](../done/slice-m4-logging-port.md) | `forbidigo.msg` referenziert nicht-existenten Logging-Port; `u-boot doctor` braucht strukturiertes Logging | M4-vorgezogen | Done |
| [`slice-m6-docker-integrationstests`](../done/slice-m6-docker-integrationstests.md) | `//go:build docker`-Pfad nur dokumentiert, kein CI-Job; erst mit Docker-Adapter sinnvoll | M6 | Done |
| [`slice-followup-verbosity-wiring`](../done/slice-followup-verbosity-wiring.md) | `--verbose`/`--debug` (LH-FA-CLI-005) waren persistent Cobra-Flags ohne Logger-Effekt | M4-followup | Done (`7c6fbce`) |
| [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md) | ADR-0004 Folgepunkte Image-Publish + Trivy; `LH-OPEN-002` Paketierung (GHCR-Anteil) | V1 | Done (T1 `0f64938`, T2 `93b703e`, T3 `8212889`, T4 `066917a`, T5 `bc487fc` — Branch-Protection-Teilabschluss 2026-05-27) |
| [`slice-v1-markdown-link-validator`](../done/slice-v1-markdown-link-validator.md) | Doku-/Link-Drift in `docs/`/`spec/` nicht maschinell geprüft | V1-vorgezogen | Done |
| [`slice-v1-backup-streaming-copy`](../done/slice-v1-backup-streaming-copy.md) | `LH-FA-INIT-005` Backup heute mit `ReadFile`+`WriteFile`; harter 256-MiB-Cap als MVP-Workaround | V1-vorgezogen | Done |
| [`slice-v1-plugin-system-entscheidung`](../done/slice-v1-plugin-system-entscheidung.md) | `LH-OPEN-003` Plugin-System offen | V1 | Done (Entscheidung in [ADR-0008](../../adr/0008-plugin-system-statisch.md): statisch) |
| [`slice-v1-template-format-entscheidung`](../done/slice-v1-template-format-entscheidung.md) | `LH-OPEN-004` Template-Format offen | V1 | Done (Entscheidung in [ADR-0009](../../adr/0009-template-format-yaml-files.md): YAML+`text/template`) |
| [`slice-v1-yaml-parse-error-sentinel`](../done/slice-v1-yaml-parse-error-sentinel.md) | M7-T5-Review-Followup N2: `YAMLCodec`-Port unterscheidet Parse- nicht von IO-Fehlern; Exit-Code-14-vs-10-Klassifikation reißt bei kaputter `compose.yaml` unter `u-boot generate devcontainer` | V1-vorgezogen | Done (`1008326`) |
| [`slice-v2-revive-custom-rules`](../done/slice-v2-revive-custom-rules.md) | ADR-0003 Folgepunkt revive-Custom-Rules | V2-vorgezogen | Done |
| [`slice-later-http-driving-adapter`](../done/slice-later-http-driving-adapter.md) | `spec/architecture.md` §7 HTTP-Driving-Adapter prospektiv | Later | Done (Entscheidung in [ADR-0010](../../adr/0010-kein-http-driving-adapter.md): wird nicht gebaut) |
| [`slice-v0.1.1-doctor-container-awareness`](../done/slice-v0.1.1-doctor-container-awareness.md) | `doctor` im distroless-Container findet docker/git nicht (Real-world-Befund 2026-05-31 post-v0.1.0) | v0.1.1-Followup | Done (T1 `9a99bbf`, T2 `c35360f`, T3 `111e725`, T4 schließt; Tag-Push bleibt Nutzer-Aktion analog v0.1.0-T4) |
| [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md) | ADR-0007 §Folgepunkte 1 Trigger (erste konkrete Cross-Plattform-Distributionsanfrage) durch `doctor`-Befund ausgelöst | V2 | Done — T1 ✅ `dc9a336` + `f3f1731` (`make build-binaries` für 6 Plattformen Linux/macOS/Windows × amd64/arm64), T2 ✅ `5e5166b` (`publish.yml` build + GitHub-Release-Upload), T3 ✅ `866f6fd` (READMEs Install-Block Binary-first + CHANGELOG `## [Unreleased]`), T4 schließt mit ADR-0007 §Entscheidung-Update (Binary „Vertagt → Gewählt") + carveouts.md `LH-OPEN-002`-Reduktion auf Homebrew+Debian/RPM + open→done. |

## Nächste Schritte

**MVP-Status: vollständig** (Audit-Trail in der MVP-Bilanz unten).
Keine MVP-blockierenden Slices mehr offen; alle ADR-getriebenen
V1/Later-Trigger-Slices sind in `done/`.

**v0.1.0 released — 2026-05-31.** Tag `v0.1.0` auf Commit `49ec464`,
`publish.yml`-Run `26717376068` grün (SemVer-Validate, GHCR-Login,
`make build VERSION=0.1.0`, OCI-Label-Verify, Live-`--version`-Smoke).
Images: `ghcr.io/pt9912/u-boot:0.1.0` und `:latest`. GitHub-Release
mit CHANGELOG-Auszug:
<https://github.com/pt9912/u-boot/releases/tag/v0.1.0>.
Audit-Trail im Release-Cut-Slice
[`done/slice-v1-release-cut-v0.1.0.md`](../done/slice-v1-release-cut-v0.1.0.md).
**Eine Nutzer-Aktion offen:** Branch-Protection-UI für `main`
spätestens vor erstem externem PR nach
[`docs/user/branch-protection.md`](../../../user/branch-protection.md)
aktivieren (Required-Status-Check-Liste: die drei verbose
Workflow-`name:`-Felder).

Aktuell offen sind nur trigger- oder nutzer-getriebene V1- und
Later-Folgen:

1. **~~v0.2.0-Tag-Push~~ — released 2026-06-01.** Tag `v0.2.0`
   auf Commit `595acdf`, publish.yml-Run `26740508466` grün
   (SemVer-Validate, GHCR-Login, GHCR-Push für `:0.2.0` + `:latest`,
   OCI-Label + Live-`--version`-Smoke, Cross-Compile + Upload
   für sechs Binary-Plattformen). Images: `ghcr.io/pt9912/u-boot:0.2.0`
   und `:latest`. GitHub-Release mit CHANGELOG-Auszug und
   sechs Binary-Assets: <https://github.com/pt9912/u-boot/releases/tag/v0.2.0>.
   Audit-Trail im Release-Cut-Slice
   [`done/slice-v1-release-cut-v0.2.0.md`](../done/slice-v1-release-cut-v0.2.0.md);
   enthaltene Slices:
   [`done/slice-v0.1.1-doctor-container-awareness.md`](../done/slice-v0.1.1-doctor-container-awareness.md),
   [`done/slice-v2-binary-distribution.md`](../done/slice-v2-binary-distribution.md),
   [`done/slice-v1-template-list.md`](../done/slice-v1-template-list.md),
   [`done/slice-v1-template-init.md`](../done/slice-v1-template-init.md).
   v0.1.1 wurde übersprungen (drei Features seit der ursprünglichen
   v0.1.1-Planung verlangten SemVer-MINOR-Bump). Branch-Protection-UI
   bleibt als Nutzer-One-Shot offen aus v0.1.0-Era (nicht
   release-blockierend).
2. **v0.3.0-Milestone — Add-on Catalogue Expansion (in progress, 1/5).**
   Fünf Slices als nächstes Release-Cluster, geordnet von klein
   nach groß:
   - **`slice-v1-audit-done`** — reiner Doku-Audit: drei vermutlich-
     erfüllte V1-IDs (`LH-FA-BUILD-006` Aggregator-Targets,
     `LH-NFA-MAINT-004` Dokumentierte Schnittstellen,
     `LH-NFA-PORT-003` Containerfreundlichkeit) verifizieren und
     in Phase-Table / MVP-Bilanz als ✅ markieren.
   - ✅ [`slice-v1-add-remove`](../done/slice-v1-add-remove.md)
     (`LH-FA-ADD-007`) — `u-boot remove <service>` geliefert
     2026-06-01 (T1 `ca1267f` Driving-Port + Skeleton, T2
     `e26cb42` State-Machine + `detectServiceState`-Extract,
     T3 `c508b4f` `--purge`-Confirmation-Gate, T4 `3cc2646`
     CLI + Wiring + E2E, T5 `764e737` Slice-Closure, Review-
     Followup `78ddcc6` F1..F6 — Two-Phase Plan-then-Write +
     InconsistentBlock-Convergence + Mode-Preservation +
     stderr-WARNING + Deactivated-Gate-Skip). E2E gegen das
     gebaute Image: `add postgres` → `enabled: true`;
     `remove postgres` → `enabled: false` + Changed-Liste;
     zweiter `remove` → idempotent No-Op. **T3-Decision**:
     actual Volume-Removal bleibt deferred (CLI surface'd
     `docker volume rm <name>`-Cleanup-Hinweis); eigener Folge-
     Slice, wenn der Docker-Engine-Port-Erweiterung gerechtfertigt
     ist.
   - **`slice-v1-addons-deps`** (`LH-FA-ADD-006`) — Add-on-
     Dependency-Resolution; Voraussetzung für Keycloak
     (`requires: [postgres]`). Domain-Modell für `requires`-Block
     im Add-on-Katalog + Validierung beim `add`/`remove`-Flow.
   - **`slice-v1-keycloak`** (`LH-FA-ADD-003` + `LH-AK-003`) —
     Keycloak-Add-on analog M5-Postgres-Pattern, mit Postgres-
     Dependency aus addons-deps.
   - **`slice-v1-otel`** (`LH-FA-ADD-004` + `LH-AK-004`) —
     OpenTelemetry-Add-on parallel zu Keycloak.
3. **V1-Templates-Implementation** — drei Slices aus
   [ADR-0009](../../adr/0009-template-format-yaml-files.md)
   §Folgepunkte.
   [`slice-v1-template-list`](../done/slice-v1-template-list.md)
   ✅ geliefert (T1 `65795b5` Domain+Driven-Port+embed.FS-Adapter
   inkl. `basic`-Bootstrap-Template, T2 `a099d63` Driving-Port +
   Application-Service, T3 `23bd91b` CLI `u-boot template list
   [--json]` + Wiring, T4 Slice-Closure mit ADR-0009-Pfad-
   Konsolidierung `external-templates/` → `externaltemplates/`).
   [`slice-v1-template-init`](../done/slice-v1-template-init.md)
   ✅ geliefert (T1 `9e81b02` `domain.TemplatePath` +
   `driven.TemplateFiles` + `Catalog.Open()`, T2 `65a1ce8`
   `TemplateInitUseCase` + Render-Loop, T3 `ed6d9a0`
   `basic`-Bootstrap-Content für die sechs `*.tmpl`-Files +
   Byte-Identity-Pin gegen Default-Init, T4 `daaaa9a` CLI-Flag
   `--template <name>` + `InitProjectService`-Delegation via
   `WithTemplateInit`-Option + E2E `diff -r`-Byte-Identity-Smoke,
   T5 Slice-Closure). Variable-Resolution für `--var key=value`
   bleibt out-of-scope, bis ein Built-in (z. B. `micronaut`)
   tatsächlich Variablen einführt. Offen: `slice-later-local-
   templates` (`--template ./pfad`-Auflösung, `LH-FA-TPL-003`).
4. **V1-Generators** — `u-boot logs` (`LH-FA-UP-005`),
   `--json`-/`--dry-run`-Output (`LH-FA-CLI-007/008`,
   `LH-NFA-USE-004`). Maschinen-Schnittstelle, auf die
   [ADR-0010](../../adr/0010-kein-http-driving-adapter.md) das
   "kein HTTP-Adapter"-Argument stützt.
5. **`LH-OPEN-002`-Distributionsweg-Restwege** — Homebrew und
   Debian/RPM bleiben vertagt mit Trigger-Slices aus
   [ADR-0007](../../adr/0007-distributionswege-ghcr.md)
   §Entscheidung. Binary-Distribution ist in
   [`done/slice-v2-binary-distribution.md`](../done/slice-v2-binary-distribution.md)
   vollständig geliefert: T1 `make build-binaries` für sechs
   Plattformen (`dc9a336`/`f3f1731`), T2 `publish.yml` baut +
   uploadet die Binaries pro Tag (`5e5166b`), T3 READMEs
   Install-Block Binary-first + CHANGELOG `## [Unreleased]`
   (`866f6fd`), T4 Slice-Closure mit ADR-0007 §Entscheidung
   „Vertagt → Gewählt". Greift ab v0.1.1 (v0.1.0 hatte noch
   keine Binary-Assets).
6. **Later** — Migration (`LH-FA-CONF-006`), Custom-Data-Sources
   (`LH-DA-004`).
7. **Podman-Drop-in als Container-Engine** — heute funktional
   per Symlink + `DOCKER_HOST=…/podman.sock`; v0.1.1-Container-
   Detection probes `/run/.containerenv` neben `/.dockerenv`,
   `LH-FA-DIAG-002` und `spec/architecture.md` §2.4 dokumentieren
   den Drop-in offiziell (`a504a36`). Formal getestete Variante
   (eigener `PodmanProbe`-Adapter + CI-Matrix) wird zum Slice,
   sobald ein konkreter Bedarf gemeldet wird.

### MVP-Bilanz — **MVP vollständig** (Stand `bc487fc`; M8-T5 `25cb123`)

**Alle 5 MVP-`LH-AK-*` gepinnt** mit benannten e2e-Tests:
LH-AK-001 (Init+Doctor) `8525c4c`, LH-AK-002 (Postgres-Flow)
`b537929`+`aa3a45c`, LH-AK-005 (Devcontainer-Init) `bfe6416`,
LH-AK-006 (Doppel-Add-Idempotenz) `8525c4c`, LH-AK-007
(Changelog) `19c4110`.

**Alle MVP-`LH-FA-*` ausgeliefert** über M1..M8 (Audit-Trail aus
M8-Review S4):

| Bereich (Spec-IDs)                                            | Slice               | Status |
| ------------------------------------------------------------- | ------------------- | ------ |
| ARCH (`LH-FA-ARCH-001..003`) + BUILD (`LH-FA-BUILD-001..009`) + PROJDOCS (`LH-FA-PROJDOCS-001..005`) | M1 + M2 + M2b..M2d  | ✅ |
| INIT (`LH-FA-INIT-001..007`)                                  | M3 + MVP-Closure-T1 | ✅ |
| DIAG (`LH-FA-DIAG-001..004`)                                  | M4                  | ✅ |
| ADD (`LH-FA-ADD-001/-002/-005` MVP)                           | M5                  | ✅ |
| UP / DOWN (`LH-FA-UP-001..004`)                               | M6                  | ✅ |
| DOC (`LH-FA-DOC-001/-003/-004` MVP — Compose / Network / Volumes) | M5 + M6 (Compose-Block-Output via `add` / `init`, Volumes via `add postgres`, Network via Compose-default-Netzwerk) | ✅ |
| GEN (`LH-FA-GEN-001..005`)                                    | M7                  | ✅ |
| DEV (`LH-FA-DEV-001/-002/-004/-005` MVP)                      | M7-T5 + MVP-Closure-T1 | ✅ |
| CONF (`LH-FA-CONF-001..005`)                                  | M8                  | ✅ (T5 `25cb123`) |
| CLI (`LH-FA-CLI-001..006` + `LH-FA-CLI-005A`)                 | M3..M8 inkrementell | ✅ |

**Software-Architecture-Schnittstellen** (`LH-SA-*`, alle MVP-
Priorität — M8-Review S5): cross-cutting, ohne dediziertes
Slice, abgedeckt durch die Implementierung:

- `LH-SA-CLI-001` Befehlsstruktur — Cobra-Layout (M3..M8).
- `LH-SA-CLI-002` Vorgesehene Befehle — alle MVP-Subkommandos
  vorhanden (siehe CLI-Zeile oben).
- `LH-SA-FILE-001` Erzeugte Dateien — M3 (init-Templates), M5
  (compose-Blocks), M7 (generate-Pfade).
- `LH-SA-FILE-002` Markierte verwaltete Bereiche — `managedblock`
  + 3 Marker-Stile (StyleHash/StyleHTMLComment/StyleDoubleSlash).
- `LH-SA-DOCKER-001` Docker Compose — DockerEngine + Probe Adapter (M6).
- `LH-SA-DOCKER-002` Containerstatus — UpService-Healthcheck-Polling +
  ComposePs-JSON-Parser (M6).

Damit ist **kein MVP-`LH-AK-*`, kein MVP-`LH-FA-*` und kein
MVP-`LH-SA-*` mehr offen**. Die Release-Maschinerie für den ersten
Schnitt (`v0.1.0` o. ä.) liegt seit
[`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md)
**bereit**: GHCR-Push via `.github/workflows/publish.yml` auf Tag `v*`,
Trivy als dritter PR-blockierender CI-Job, ADR-0007 setzt GHCR als
primären Distributionsweg. Der Tag-Push selbst bleibt Nutzer-Trigger.

### v0.3.0 — Milestone-Tabelle „Add-on Catalogue Expansion"

Aktiver Release-Scope. Fünf Slices, geordnet von klein nach groß
(audit-done → add-remove → addons-deps → Keycloak + OTel parallel).
Die Spec-ID-Spalte mappt jeden Slice auf die `LH-FA-*` / `LH-AK-*`-
Arbeitspunkte aus dem [Lastenheft](../../../../spec/lastenheft.md), die
er ausliefert; mehr V1-Arbeitspunkte (Generators, `--json`/
`--dry-run`, Logs) bleiben in §Nächste Schritte 4 als Post-v0.3.0-
Backlog. v0.3.0-Tag-Push wird wie v0.1.0/v0.2.0 Nutzer-Aktion sein
(Release-Cut-Slice analog `slice-v1-release-cut-v0.2.0` bei
Milestone-Schluss).

| Done | Slice | Spec-IDs (Lastenheft) | Status / Bezug |
| ---- | ----- | --------------------- | -------------- |
| [ ] | `slice-v1-audit-done` | `LH-FA-BUILD-006` Aggregator-Targets, `LH-NFA-MAINT-004` Dokumentierte Schnittstellen, `LH-NFA-PORT-003` Containerfreundlichkeit | Reiner Doku-Audit — drei vermutlich-erfüllte V1-IDs gegen den aktuellen Code-Stand verifizieren und in Phase-Tabelle / MVP-Bilanz als ✅ markieren. Bringt keine Code-Änderung; gute Aufwärm-Tranche vor den drei Implementations-Slices. |
| [x] | [`slice-v1-add-remove`](../done/slice-v1-add-remove.md) | `LH-FA-ADD-007` Service entfernen | Geliefert 2026-06-01 — T1 `ca1267f`, T2 `e26cb42`, T3 `c508b4f`, T4 `3cc2646`, T5 `764e737` + Review-Followup `78ddcc6` (F1..F6: Two-Phase + Mode-Preservation + stderr-WARNING + Deactivated-Gate-Skip). Dependency-Check und Volume-Removal defer-pflichtig auf `slice-v1-addons-deps` bzw. eigenen Folge-Slice. |
| [ ] | `slice-v1-addons-deps` | `LH-FA-ADD-006` Add-on-Abhängigkeiten | Voraussetzung für Keycloak (`requires: [postgres]`). Domain-Modell für `requires`-Block im Add-on-Katalog + Validierung beim `add`/`remove`-Flow. Liefert die fehlende „Dependency-Check"-Hälfte aus `LH-FA-ADD-007`. |
| [ ] | `slice-v1-keycloak` | `LH-FA-ADD-003` Keycloak hinzufügen, `LH-AK-003` Keycloak-Flow | Keycloak-Add-on analog M5-Postgres-Pattern; Postgres-Dependency über addons-deps deklariert. Acceptance-Test analog `LH-AK-002`. |
| [ ] | `slice-v1-otel` | `LH-FA-ADD-004` OpenTelemetry hinzufügen, `LH-AK-004` OpenTelemetry-Flow | OpenTelemetry-Add-on parallel zu Keycloak — kann gleichzeitig mit Keycloak entwickelt werden, weil keine Dependency zwischen den beiden. |

Stand: 1/5 ✅, 4/5 offen. Beim Schließen des Milestones folgt der
v0.3.0-Release-Cut-Slice mit CHANGELOG-Konsolidierung, Dev-Version-
Bump auf `0.3.0-dev`, READMEs-Sync und Tag-Push-Nutzer-Aktion
(Pattern aus `slice-v1-release-cut-v0.2.0`).

### v0.1.0 / v0.2.0 — Audit-Trail ausgelieferter Slices + offene Trigger-Restwege

Audit-Trail der nicht-MVP-Slices die zwischen v0.1.0 und v0.2.0
ausgeliefert wurden — ADR-getriebene Entscheidungs-Trigger,
vorgezogene Review-Followups, Release-Cuts, und ADR-§Folgepunkte-
Implementierungen — plus die drei noch unausgelösten Trigger-
Restwege (Homebrew, Debian/RPM, local-templates). Die Tabelle
zeigt den Stand mit Checkbox je Slice; noch nicht implementierte
V1-Folgen (Add-ons, Generators, Logs/Dry-Run/Diff) und der
laufende v0.3.0-Milestone stehen in §Nächste Schritte oben.

| Done | Slice | Kategorie | Audit-Trail / Bezug |
| ---- | ----- | --------- | ------------------- |
| [x] | [`slice-v1-plugin-system-entscheidung`](../done/slice-v1-plugin-system-entscheidung.md) | ADR-Entscheidung (V1) | [ADR-0008](../../adr/0008-plugin-system-statisch.md): Add-on-System statisch (keine Plugins); vier Re-Eval-Trigger in §Folgepunkte. |
| [x] | [`slice-v1-template-format-entscheidung`](../done/slice-v1-template-format-entscheidung.md) | ADR-Entscheidung (V1) | [ADR-0009](../../adr/0009-template-format-yaml-files.md): YAML-Metadaten + `text/template`-Files. Drei Implementierungs-Slices in §Folgepunkte. |
| [x] | [`slice-later-http-driving-adapter`](../done/slice-later-http-driving-adapter.md) | ADR-Entscheidung (Later) | [ADR-0010](../../adr/0010-kein-http-driving-adapter.md): kein HTTP-Adapter; CLI-only + `LH-NFA-USE-004`. Zwei Re-Eval-Trigger. |
| [x] | [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md) | ADR-Mechanik (V1) | [ADR-0004](../../adr/0004-ci-system.md) + [ADR-0007](../../adr/0007-distributionswege-ghcr.md). T1 `0f64938`, T2 `93b703e`, T3 `8212889`, T4 `066917a`, T5 `bc487fc`. |
| [x] | [`slice-v1-yaml-parse-error-sentinel`](../done/slice-v1-yaml-parse-error-sentinel.md) | V1-vorgezogen | M7-T5-N2 Review-Followup-Closure. `1008326`. |
| [x] | [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md) | Release-Cut | Version-Pin im Build + CHANGELOG-Bootstrap. T1 `056e4c6`, T2 `f176e95`, T3 `4fc93a9`. v0.1.0 released 2026-05-31. |
| [x] | [`slice-v0.1.1-doctor-container-awareness`](../done/slice-v0.1.1-doctor-container-awareness.md) | v0.1.0-Real-World-Feedback | T1 `9a99bbf`, T2 `c35360f`, T3 `111e725`, T4 `f3f1731` — `doctor` skipped Host-Prerequisite-Checks im Container-Modus. |
| [x] | [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md) | ADR-0007 §Folgepunkte 1 Trigger | T1 `dc9a336`/`f3f1731`, T2 `5e5166b`, T3 `866f6fd`, T4 `2f39511` — sechs Plattformen als GitHub-Release-Asset. |
| [x] | [`slice-v1-template-list`](../done/slice-v1-template-list.md) | ADR-0009 §Folgepunkte 1 | T1 `65795b5`, T2 `a099d63`, T3 `23bd91b`, T4 `a7e0d7b` + Review-Followup `c807cdb` (N1..N5). `LH-FA-TPL-004`. |
| [x] | [`slice-v1-template-init`](../done/slice-v1-template-init.md) | ADR-0009 §Folgepunkte 2 | T1 `9e81b02`, T2 `65a1ce8`, T3 `ed6d9a0`, T4 `daaaa9a`, T5 `133622f` + Review-Followup `7fe26e0` (F1..F5). `LH-FA-TPL-001`. |
| [x] | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) | Release-Cut | CHANGELOG-Konsolidierung + Dev-Version-Bump + Status-Sync. T1 `be139cb`, T2 `1823598`, T3 `6b9cc6c`. v0.2.0 released 2026-06-01. |
| [ ] | `slice-later-local-templates` | ADR-0009 §Folgepunkte 3 | `LH-FA-TPL-003` (Later). `--template ./pfad`-Auflösung; eigener Slice bei konkretem Bedarf. |
| [ ] | `slice-v2-homebrew-formula` | ADR-0007 Restwege | Trigger: erste macOS-Nutzer-Nachfrage. Plan-Anker in ADR-0007 §Entscheidung. |
| [ ] | `slice-v2-distro-pakete` | ADR-0007 Restwege | Trigger: konkrete Distro-Anfrage (`debhelper`/`rpmbuild`-Overhead). Plan-Anker in ADR-0007 §Entscheidung. |

Die noch offenen V1- und Later-Folgen ohne benannten Slice
(Add-ons-Implementation außerhalb des aktuellen v0.3.0-Milestones,
Generators, Logs/Dry-Run/Diff, Migration, Custom-Data-Sources)
sind oben in §Nächste Schritte aufgeführt.

## Lifecycle-Hinweis

Diese Datei ist die einzige zulässige Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und [`../README.md`](../README.md)).
