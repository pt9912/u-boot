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
| M2c CI | Done | GitHub-Actions-CI (`LH-QA-003` auf konkret gehoben), `.github/workflows/ci.yml` mit Jobs `gates` + `security-gates` (beide PR-blockierend), SHA-pinned Actions, Docker-only, ADR-0004 | [`slice-m2c-ci-pipeline`](../done/slice-m2c-ci-pipeline.md) |
| M2d Carveouts | Done | Carveout-Disziplin (`LH-FA-PROJDOCS-005` MVP-Pflicht), Master-Inventar [`carveouts.md`](carveouts.md), 7 neue Slice-Pläne in [`open/`](../open/) für offene Carveouts; permanente Carveouts dokumentiert | [`slice-m2d-carveout-disziplin`](../done/slice-m2d-carveout-disziplin.md) |
| M3 `u-boot init` | Done | Projektstruktur erzeugen (`LH-FA-INIT-001..007`), `u-boot.yaml` schreiben, Git-Init, Re-Init mit `--force`/`--backup` (LH-FA-INIT-005) + Modi-Flags (LH-FA-CLI-005A). Coverage-, depguard- und gomodguard-Carveouts aufgelöst. Detail: [`slice-m3-init-flow.md`](../done/slice-m3-init-flow.md). **Stand:** T1..T4c ✅ (Commits siehe Slice-DoD); T5 ✅ `scripts/verify-depguard.sh` + `make verify-depguard`; M3-followup: [`slice-m3-build-polish`](../done/slice-m3-build-polish.md) (`987c164`, govulncheck-Pin + PROGRESS_FLAG) und [`slice-m3-gomodguard-rules`](../done/slice-m3-gomodguard-rules.md) (`201fb4b`, 4 Block-Regeln + golangci-lint v2.12.2) |
| M4 `u-boot doctor` | Done | Lokale Voraussetzungen prüfen (`LH-FA-DIAG-001..004`), Severity-Klassifikation, Repair-Hints. 9 Checks: write-permissions, git, docker (+reachable+compose-plugin), u-boot.yaml, compose.yaml, devcontainer.json/Dockerfile. CLI `doctor`-Subkommando mit `--strict`. Exit-Code 11 bei Errors (oder Warns + --strict). | [`slice-m4-doctor`](../done/slice-m4-doctor.md) |
| M5 `u-boot add postgres` | In progress | PostgreSQL-Add-on (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`), services-Schema in u-boot.yaml, Compose-Block, `.env.example`-Block, Healthcheck, State-Machine für Re-Aktivierung + Inkonsistenz-Erkennung. Detail: [`slice-m5-add-postgres.md`](slice-m5-add-postgres.md). **Stand:** T1 ✅ `995726a` (services-Schema + ServiceName/ServiceState-Domain); T2..T7 offen | [`slice-m5-add-postgres`](slice-m5-add-postgres.md) |
| M6 `u-boot up` / `down` | Open | Compose-Wrapper (`LH-FA-UP-001..004`), Healthcheck-Polling, `--timeout`, `--volumes` | offen |
| M7 `u-boot generate` | Open | `generate changelog`/`readme`/`env-example`/`devcontainer` (`LH-FA-GEN-001..005`) | offen |
| M8 `u-boot config` | Open | `config get`/`set`/Anzeigen (`LH-FA-CONF-001..005`), Schema-Validierung | offen |
| MVP-Closure | Open | Devcontainer-Mindestumfang (`LH-FA-DEV-001..005`), MVP-Acceptance-Flows (`LH-AK-001..002`, `LH-AK-005..007`) | offen |
| V1 Keycloak / OTel | Open | `LH-FA-ADD-003`, `LH-FA-ADD-004`, `LH-AK-003`, `LH-AK-004` | offen |
| V1 Templates | Open | `LH-FA-TPL-001..004` | offen |
| V1 Logs / Dry-Run / Diff | Open | `LH-FA-UP-005`, `LH-FA-CLI-007/008` | offen |
| Later Migration / Custom Templates | Open | `LH-FA-CONF-006`, `LH-FA-TPL-003`, `LH-DA-004` | offen |

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
| [`slice-m6-docker-integrationstests`](../open/slice-m6-docker-integrationstests.md) | `//go:build docker`-Pfad nur dokumentiert, kein CI-Job; erst mit Docker-Adapter sinnvoll | M6 | Open |
| [`slice-v1-release-pipeline`](../open/slice-v1-release-pipeline.md) | ADR-0004 Folgepunkte Image-Publish + Trivy; `LH-OPEN-002` Paketierung (GHCR-Anteil) | V1 | Open |
| [`slice-v1-markdown-link-validator`](../done/slice-v1-markdown-link-validator.md) | Doku-/Link-Drift in `docs/`/`spec/` nicht maschinell geprüft | V1-vorgezogen | Done |
| [`slice-v1-backup-streaming-copy`](../done/slice-v1-backup-streaming-copy.md) | `LH-FA-INIT-005` Backup heute mit `ReadFile`+`WriteFile`; harter 256-MiB-Cap als MVP-Workaround | V1-vorgezogen | Done |
| [`slice-v1-plugin-system-entscheidung`](../open/slice-v1-plugin-system-entscheidung.md) | `LH-OPEN-003` Plugin-System offen | V1 | Open |
| [`slice-v1-template-format-entscheidung`](../open/slice-v1-template-format-entscheidung.md) | `LH-OPEN-004` Template-Format offen | V1 | Open |
| [`slice-v2-revive-custom-rules`](../done/slice-v2-revive-custom-rules.md) | ADR-0003 Folgepunkt revive-Custom-Rules | V2-vorgezogen | Done |
| [`slice-later-http-driving-adapter`](../open/slice-later-http-driving-adapter.md) | `spec/architecture.md` §7 HTTP-Driving-Adapter prospektiv | Later | Open |

## Nächste Schritte

1. **M5 add-postgres**: 7 Tranchen geplant (siehe [`slice-m5-add-postgres.md`](slice-m5-add-postgres.md)). Aktive Tranche: T1 (u-boot.yaml services-Schema + Domain-Types).

## Lifecycle-Hinweis

Diese Datei ist die einzige zulässige Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und [`../README.md`](../README.md)).
