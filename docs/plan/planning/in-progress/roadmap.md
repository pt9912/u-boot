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
| M5 `u-boot add postgres` | Done | PostgreSQL-Add-on (`LH-FA-ADD-001..002`, `LH-FA-ADD-005`), services-Schema in u-boot.yaml, Compose split-block scaffold, `.env.example`-Block, Healthcheck mit `$${POSTGRES_USER:-postgres}`-Defaults für LH-AK-002, State-Machine, Active-Repair, CLI-Subcommand, doctor-Integration (services.enabled-key + devcontainer.forwardPorts + devcontainer.enabled-Severity-Eskalation). 11 doctor-Checks gesamt. | [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) |
| M6 `u-boot up` / `down` | Done | Compose-Wrapper (`LH-FA-UP-001..004`), Healthcheck-Polling, `--timeout`, `--volumes`, CLI-Subcommands + Status-Tabelle. Alle 7 Tranchen in `done/`. Carveout-Slice [`slice-m6-docker-integrationstests`](../done/slice-m6-docker-integrationstests.md) **Done** (Sub-T1..T4 + Audit-Härtung `41cab1b` + Stabilisierung `43b42e4`/`8865ca1` + Carveout-Entfernung). | [`slice-m6-up-down`](../done/slice-m6-up-down.md) |
| M7 `u-boot generate` | Done | `generate changelog`/`readme`/`env-example`/`devcontainer` (`LH-FA-GEN-001..005` + LH-FA-DEV-001/004/005 + LH-AK-007). Sechs Tranchen ✅ (`67fc181`/`3c5de48`/`037ab00`/`19c4110`/`294e492`/`d32a733`) + Review-Followup `27de9c5` (9 Findings S1..S4/N1..N5 adressiert, u. a. fenced-code-block-Schutz gegen Markdown-Korruption + CRLF-Normalisierung). Reuse von `managedblock` (3 Marker-Stile decken 4 Datei-Mappings), `generateManagedFile`-Helper für env-example/readme, atomarer Two-File-Plan für devcontainer, konservative User-Edit-Erkennung für changelog. CLI: `u-boot generate <artifact>`, Exit-Codes 0/2/10/14. | [`slice-m7-generate`](../done/slice-m7-generate.md) |
| M8 `u-boot config` | Open | `config get`/`set`/Anzeigen (`LH-FA-CONF-001..005`), Schema-Validierung | offen |
| MVP-Closure | Done | Devcontainer-Mindestumfang (`LH-FA-DEV-001..005`), MVP-Acceptance-Flows (`LH-AK-001..002`, `LH-AK-005..007`). Drei Tranchen — T1 ✅ `bfe6416` `u-boot init --devcontainer` (LH-AK-005), T2 ✅ `8525c4c` LH-AK-001/-006-Pins in `acceptance_test.go` inkl. Doctor-Severity-Fix `compose.yaml.valid` (Error → Warn), T3 ✅ Slice-Closure + MVP-Bilanz. Alle 5 MVP-`LH-AK-*` gepinnt; alle MVP-`LH-FA-DEV-*` ausgeliefert. **Verbleibender MVP-Item: `LH-FA-CONF-001..005` (M8 `u-boot config`)**. | [`slice-mvp-closure`](../done/slice-mvp-closure.md) |
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
| [`slice-m6-docker-integrationstests`](../done/slice-m6-docker-integrationstests.md) | `//go:build docker`-Pfad nur dokumentiert, kein CI-Job; erst mit Docker-Adapter sinnvoll | M6 | Done |
| [`slice-followup-verbosity-wiring`](../done/slice-followup-verbosity-wiring.md) | `--verbose`/`--debug` (LH-FA-CLI-005) waren persistent Cobra-Flags ohne Logger-Effekt | M4-followup | Done (`7c6fbce`) |
| [`slice-v1-release-pipeline`](../open/slice-v1-release-pipeline.md) | ADR-0004 Folgepunkte Image-Publish + Trivy; `LH-OPEN-002` Paketierung (GHCR-Anteil) | V1 | Open |
| [`slice-v1-markdown-link-validator`](../done/slice-v1-markdown-link-validator.md) | Doku-/Link-Drift in `docs/`/`spec/` nicht maschinell geprüft | V1-vorgezogen | Done |
| [`slice-v1-backup-streaming-copy`](../done/slice-v1-backup-streaming-copy.md) | `LH-FA-INIT-005` Backup heute mit `ReadFile`+`WriteFile`; harter 256-MiB-Cap als MVP-Workaround | V1-vorgezogen | Done |
| [`slice-v1-plugin-system-entscheidung`](../open/slice-v1-plugin-system-entscheidung.md) | `LH-OPEN-003` Plugin-System offen | V1 | Open |
| [`slice-v1-template-format-entscheidung`](../open/slice-v1-template-format-entscheidung.md) | `LH-OPEN-004` Template-Format offen | V1 | Open |
| [`slice-v1-yaml-parse-error-sentinel`](../open/slice-v1-yaml-parse-error-sentinel.md) | M7-T5-Review-Followup N2: `YAMLCodec`-Port unterscheidet Parse- nicht von IO-Fehlern; Exit-Code-14-vs-10-Klassifikation reißt bei kaputter `compose.yaml` unter `u-boot generate devcontainer` | V1 | Open |
| [`slice-v2-revive-custom-rules`](../done/slice-v2-revive-custom-rules.md) | ADR-0003 Folgepunkt revive-Custom-Rules | V2-vorgezogen | Done |
| [`slice-later-http-driving-adapter`](../open/slice-later-http-driving-adapter.md) | `spec/architecture.md` §7 HTTP-Driving-Adapter prospektiv | Later | Open |

## Nächste Schritte

1. **M6 up/down**: **Done** (siehe [`slice-m6-up-down.md`](../done/slice-m6-up-down.md)). Alle 7 Tranchen abgeschlossen: T1 ✅ `9f8badd`, T2 ✅ `84a676c`, T3 ✅ `1e5ef18`, T4 ✅ `1351cfb` (+ fund `9101bdc`, + review `d1deee5`), T5 ✅ `a46bec3`, T6 ✅ `4a7e60d`, T7 ✅ `6d9aa88` (+ review `adeea13` für `up --quiet`-Vertrag). Coverage 91.20%.
2. **Verbosity-Wiring**: **Done** (`7c6fbce`, siehe [`slice-followup-verbosity-wiring.md`](../done/slice-followup-verbosity-wiring.md)). `--quiet` → `slog.LevelWarn`, `--verbose`/`--debug` → `slog.LevelDebug` via `*slog.LevelVar` und Cobra-`PersistentPreRunE`. Carveouts-Eintrag entfernt; temporäre Carveouts jetzt 6 statt 7.
3. **M6-docker-int**: **Done** (siehe [`slice-m6-docker-integrationstests.md`](../done/slice-m6-docker-integrationstests.md)). Iterations-Bilanz: PR #1 ohne Merge geschlossen 2026-05-30 wegen (a) deterministisch roten `TestE2E_LHAK002_PostgresAcceptanceFlow` (LH-FA-INIT-006-Regex-Reject bei `t.TempDir()`-Counter-Leaf, gefixt `aa3a45c`) und (b) Audit-Heuristik, die nur workflow-level `conclusion` prüfte (gehärtet `41cab1b` mit Job-level-Pflicht). PR #2 (`43b42e4` → Merge `8865ca1`) entfernte `continue-on-error: true` nach drei job-level-grünen Läufen auf `main` (`41cab1b`/`379fe21`/`2fa46fd`). Erster Lauf ohne `continue-on-error` grün auf `8865ca1` (https://github.com/pt9912/u-boot/actions/runs/26679225340). Carveout entfernt.
4. **M7 generate**: **Done** (siehe [`slice-m7-generate.md`](../done/slice-m7-generate.md)).
   Alle 6 Tranchen abgeschlossen: T1 ✅ `67fc181` Port+Skeleton,
   T2 ✅ `3c5de48` env-example, T3 ✅ `037ab00` readme,
   T4 ✅ `19c4110` changelog (LH-AK-007), T5 ✅ `294e492`
   devcontainer (LH-FA-DEV-001/004/005), T6 ✅ `d32a733`
   CLI-Subcommand + ExitCode-Wiring (ErrArtifactUnknown→2,
   ErrGenerateManualConflict→10, ErrGenerateFileSystem→14).
   Review-Followup `27de9c5`: 9 Findings (S1..S4 + N1..N5) aus
   Post-Merge-Code-Review adressiert, u. a. fenced-code-block-
   Schutz im Changelog-Handler (verhindert Markdown-Korruption bei
   dokumentierten Versions-Beispielen) und CRLF-Normalisierung im
   bytes.Equal-Heuristik (CRLF-Files registrieren als fresh statt
   silent CRLF→LF zu rewritern). Coverage 90.20 %.
5. **MVP-Closure**: **Done** (siehe [`slice-mvp-closure.md`](../done/slice-mvp-closure.md)).
   Drei Tranchen abgeschlossen: T1 ✅ `bfe6416` `init
   --devcontainer` (LH-AK-005), T2 ✅ `8525c4c` LH-AK-001/-006
   benannte e2e-Pins + Doctor-Severity-Fix für
   `compose.yaml.valid` no-services (Error → Warn,
   LH-AK-001-§2299-Konformität), T3 ✅ Closure + MVP-Bilanz.

### MVP-Bilanz (Stand `c1fcd58`/`8525c4c`)

**Was MVP-komplett ist:**

- **Alle 5 MVP-`LH-AK-*` gepinnt** mit benannten e2e-Tests:
  LH-AK-001 (Init+Doctor) `8525c4c`, LH-AK-002 (Postgres-Flow)
  `b537929`+`aa3a45c`, LH-AK-005 (Devcontainer-Init) `bfe6416`,
  LH-AK-006 (Doppel-Add-Idempotenz) `8525c4c`, LH-AK-007
  (Changelog) `19c4110`.
- **Alle MVP-`LH-FA-*` aus M1..M7 + MVP-Closure ausgeliefert**:
  ARCH/BUILD/PROJDOCS (M1/M2), INIT (M3), DIAG (M4), ADD (M5),
  UP (M6), GEN + DEV (M7 + MVP-T1), CLI durchgehend
  inkrementell. CONF-002 (YAML-Schema) ist seit M3
  ausgeliefert.

**Was MVP noch fehlt:**

- **`M8 u-boot config`**: `LH-FA-CONF-001/-003/-004/-005`
  verlangen `u-boot config get/set/list`-Kommandos. Heute
  existiert kein Subcommand; der User muss `u-boot.yaml`
  manuell editieren. Das ist ein eigener Slice (siehe Roadmap-
  Tabelle, Phase M8), kein MVP-Closure-Scope. **MVP ist erst
  vollständig nach M8.**

### V1-Phase: nicht release-blockierend, Trigger-getrieben

Vier offene V1-Slices warten auf konkrete Trigger:

- [`slice-v1-release-pipeline`](../open/slice-v1-release-pipeline.md):
  Trigger erster GHCR-Tag oder erster externer PR.
- [`slice-v1-plugin-system-entscheidung`](../open/slice-v1-plugin-system-entscheidung.md):
  Trigger erster externer Plugin-Vorschlag.
- [`slice-v1-template-format-entscheidung`](../open/slice-v1-template-format-entscheidung.md):
  Trigger erster externer Template-Vorschlag.
- [`slice-v1-yaml-parse-error-sentinel`](../open/slice-v1-yaml-parse-error-sentinel.md):
  Trigger User-Bug zu Exit-14-vs-10 bei kaputter compose.yaml,
  oder nächster YAMLCodec-Port-Touch.

Plus die V1-Add-ons (LH-AK-003 Keycloak, LH-AK-004 OTel) und
V1-Generators (`u-boot logs`, `--json`-Output).

## Lifecycle-Hinweis

Diese Datei ist die einzige zulässige Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und [`../README.md`](../README.md)).
