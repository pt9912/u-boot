# Quality Gates — u-boot

| Dokument         | Quality-Gates-Übersicht                                       |
| ---------------- | -------------------------------------------------------------- |
| Projektname      | `u-boot`                                                       |
| Bezug            | `LH-QA-001..004`, `LH-FA-BUILD-001..009` in [`spec/lastenheft.md`](../../spec/lastenheft.md) |
| ADR              | [`docs/plan/adr/0003-solid-nahes-lint-profil.md`](../plan/adr/0003-solid-nahes-lint-profil.md) |
| Status           | Entwurf 0.1.0                                                  |
| Datum            | 2026-05-21                                                     |

## Zweck

Übersicht über die verbindlichen Quality-Gates der u-boot-Codebase: was
läuft im `lint`-Stage, was im `test`-Stage, welche Schwellen gelten,
welche Carveouts sind dokumentiert. Das Dokument ist die Detail-Doku
zu `LH-QA-004`; die Pflichtaussagen leben im Lastenheft, die konkreten
Konfigurations- und Schwellwerte hier.

---

## 1. Statische Analyse (`golangci-lint`)

Statische Analyse läuft Docker-basiert über die `lint`-Stage des
Top-Level-`Dockerfile` (`LH-FA-BUILD-001`):

```bash
docker build --target lint -t u-boot:lint .
# oder bequemer:
make lint
```

Die Stage führt `golangci-lint run ./...` mit Default-Lintern und dem
SOLID-nahen Zusatzprofil aus:

| Linter        | Zweck                                            |
| ------------- | ------------------------------------------------ |
| `govet`       | semantische Korrektheit (z. B. printf-Argumente) |
| `errcheck`    | Fehlerwerte werden nicht ignoriert               |
| `staticcheck` | klassische Bug-Patterns + Style                  |
| `unused`      | toter Code                                       |
| `ineffassign` | unwirksame Zuweisungen                           |

Die vollständige Konfiguration (5 Defaults + 24 SOLID-nahe Linter
aus §1.2; `depguard` ist als Schicht-Regel-Linter Teil dieser 24,
siehe `LH-FA-ARCH-003`) lebt in
[`.golangci.yml`](../../.golangci.yml). Damit sind 29 Linter aktiv. `//nolint`-Suppressions bleiben
ausgeschlossen — falls ein Linter auf einem Pfad designseitig keinen
Sinn ergibt (z. B. `testpackage` im `cmd/uboot`-Wiring), wird der Pfad
per `issues.exclude-rules` mit `Why:`-Kommentar ausgenommen; dort
dokumentierte Scope-Definitionen sind keine Suppressions, sondern
bewusste Profil-Entscheidungen. Verstöße brechen den Build.

Verbindliche Make-Targets (`LH-FA-BUILD-005`/`-006`):

```bash
make lint            # nur statische Analyse
make gates           # lint + test + coverage-gate
make ci              # gates + govulncheck
```

### 1.1 (reserviert)

Erweitert die Lint-Sektion, sobald u-boot um weitere Sprach-Stacks
(z. B. Skript-Suite, Shell-Helper unter `scripts/`) ergänzt wird.
Aktuell ist u-boot ein reines Go-Projekt; analog zur m-trace-Vorlage
ist dieser Slot für TypeScript/Shell/etc. reserviert.

### 1.2 Go: SOLID-nahe Linter

Die folgenden `golangci-lint`-Linter sind keine offizielle
SOLID-Kategorie. Sie sind die verbindliche Projektauswahl für
SOLID-nahe Designsignale: geringe Komplexität und kleine
Verantwortlichkeiten (SRP), schlanke Interfaces (ISP), stabile
Import-/Modulgrenzen (DIP) oder reduzierte globale Kopplung.

| Linter             | Kurzbeschreibung                                 | Teil von SOLID |
| ------------------ | ------------------------------------------------ | -------------- |
| `containedctx`     | `context.Context` nicht in Structs speichern      | Y              |
| `contextcheck`     | Context korrekt weiterreichen                     | Y              |
| `cyclop`           | Zyklomatische Komplexität                         | Y              |
| `depguard`         | Import-Regeln/Layer-Grenzen (`LH-FA-ARCH-003`)    | Y              |
| `dupl`             | Code-Duplikate                                    | Y              |
| `fatcontext`       | Context in Loops/Closures                         | Y              |
| `forbidigo`        | Verbotene Identifier/APIs                         | Y              |
| `funlen`           | Zu lange Funktionen                               | Y              |
| `gochecknoglobals` | Keine globalen Variablen                          | Y              |
| `gochecknoinits`   | Keine `init()`-Funktionen                         | Y              |
| `gocognit`         | Kognitive Komplexität                             | Y              |
| `gocyclo`          | Zyklomatische Komplexität                         | Y              |
| `gomodguard_v2`    | Modul-Allow-/Blocklist                            | Y              |
| `iface`            | Interface-Pollution vermeiden                     | Y              |
| `inamedparam`      | Interface-Parameter benennen                      | Y              |
| `interfacebloat`   | Zu große Interfaces                               | Y              |
| `ireturn`          | Interfaces annehmen, konkrete Typen zurückgeben   | Y              |
| `maintidx`         | Maintainability Index                             | Y              |
| `nestif`           | Tiefe `if`-Verschachtelung                        | Y              |
| `noctx`            | HTTP-Aufrufe ohne Context                         | Y              |
| `reassign`         | Package-Variablen nicht neu zuweisen              | Y              |
| `revive`           | Konfigurierbarer Stil-/Design-Linter              | Y              |
| `testpackage`      | Externe `_test`-Packages (`<name>_test`)          | Y              |
| `unparam`          | Ungenutzte Parameter                              | Y              |

Schwellen (in `.golangci.yml` zentral; Vorlage m-trace `apps/api`):

| Linter           | Setting                | Wert |
| ---------------- | ---------------------- | ---- |
| `cyclop`         | `max-complexity`       | 15   |
| `dupl`           | `threshold`            | 150  |
| `funlen`         | `lines`                | 100  |
| `funlen`         | `statements`           | 60   |
| `gocognit`       | `min-complexity`       | 20   |
| `gocyclo`        | `min-complexity`       | 15   |
| `interfacebloat` | `max`                  | 10   |
| `maintidx`       | `under`                | 20   |
| `nestif`         | `min-complexity`       | 5    |

Alle nicht in der Tabelle aufgeführten Linter (z. B. `containedctx`,
`contextcheck`, `fatcontext`, `gochecknoglobals`, `gochecknoinits`,
`iface`, `inamedparam`, `reassign`, `testpackage`, …) laufen mit ihren
Default-Schwellen.

Wert-Hebung ist Routine (Commit-Body begründet, kein eigener ADR
nötig), solange das Profil insgesamt nicht aufweicht.

### 1.3 Pflicht-Carveouts

Die folgenden Scope-Definitionen sind in `.golangci.yml` unter
`issues.exclude-rules` zentral dokumentiert. Sie sind **keine**
Suppressions, sondern bewusste Profil-Entscheidungen:

| Pfad / Pattern        | Ausgenommene Linter                                          | Warum |
| --------------------- | ------------------------------------------------------------ | ----- |
| `_test\.go$`          | `cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`          | Tabellengetriebene Tests haben oft hohe Komplexität ohne Designschaden. |
| `_test\.go$`          | `noctx`, `unparam`, `revive` (`unused-parameter`)            | Tests gegen `httptest.Server`/Fakes brauchen die Linter nicht. |
| `cmd/uboot/`          | `testpackage`, `gochecknoglobals` (für `var version`)         | Wiring-Schicht braucht `main`-Package; `version` wird per `-ldflags` überschrieben. |
| `errcheck` global     | `fmt.Fprintln`, `fmt.Fprintf`, `fmt.Fprint` (stdout/stderr)  | CLI-Writes können nicht meaningful fehlschlagen; `_, _ =`-Prefix bringt keinen Wert. |

---

## 2. Tests

### 2.1 Inner-Loop: `make test`

Tests laufen Docker-basiert über die `test`-Stage:

```bash
make test
```

Die Stage führt `go test ./...` aus. Build-Tags wie
`//go:build docker` werden hier NICHT gesetzt — entsprechende
Integrationstests bleiben aus dem Default-Pfad ausgeschlossen
(`internal/adapter/driven/docker/engine_docker_test.go`,
`internal/hexagon/application/upservice_*_docker_test.go`,
`internal/e2e/*_docker_test.go`).

### 2.2 Integration-Loop: `make test-docker`

Adapter- und e2e-Integrationstests gegen eine echte Docker-Engine
laufen unter dem `docker`-Build-Tag (`spec/architecture.md` §5,
slice `slice-m6-docker-integrationstests`):

```bash
make test-docker
```

Das Target baut zuerst die Dockerfile-Stage `test-docker-tools`
(golang + `docker-ce-cli` + `docker-compose-plugin`) und startet
dann den Test-Container mit `--network=host` plus gemountetem
Docker-Socket. **Beide** Bedingungen sind nötig: das Test-Binary
muss den Docker-Daemon erreichen UND im selben Network-Namespace
wie der Daemon laufen, damit `NetProbe.DialTCP("localhost", ...)`-
Pins gegen Compose-veröffentlichte Ports funktionieren.

CI-Pfad: separater Workflow `.github/workflows/integration.yml`
(siehe §6). Aktuell mit `continue-on-error: true` als
Stabilisierungs-Maßnahme; Aufhebung nach drei aufeinanderfolgenden
grünen Läufen auf `main` (`run_attempt=1, event=push`).

Pin-Inventar (Stand M6-docker-int Sub-T3):

| Spec-ID | Test-Datei |
| ------- | ---------- |
| LH-NFA-PERF-002 | `internal/adapter/driven/docker/engine_progressstream_docker_test.go` |
| LH-FA-DIAG-002 | `internal/adapter/driven/docker/engine_psjsonschema_docker_test.go` |
| LH-FA-UP-001 §966 | `internal/hexagon/application/upservice_healthcheck_docker_test.go` |
| LH-FA-UP-001 §968 | `internal/hexagon/application/upservice_portprobe_docker_test.go` |
| LH-AK-002 | `internal/e2e/postgres_acceptance_docker_test.go` |
| LH-FA-UP-004 §1015 | `internal/e2e/down_volumes_docker_test.go` |

---

## 3. Coverage

Coverage-Messung über die `coverage`-Stage, bootstrap-aware
(`LH-FA-BUILD-008`):

```bash
make coverage-gate
```

- Schwellwert: **90 %** (Makefile-Default `THRESHOLD ?= 90`,
  Dockerfile-Build-Arg `ARG COVERAGE_THRESHOLD=90`). Aktiviert mit
  M3-T1, sobald produktive Pakete unter `./internal/...` existieren.
- Override pro Aufruf, z. B. zur lokalen Diagnose:
  `make coverage-gate THRESHOLD=80`.
- Bootstrap-Pfad in `scripts/coverage-gate.sh` bleibt erhalten für
  den Fall, dass `./internal/...` jemals wieder leer wäre — wirkt
  aber nicht mehr in der Produktion (`internal/` ist seit M3-T1
  bestückt).

---

## 4. Security

```bash
make govulncheck     # Go-Modul-CVEs gegen die installierten Versionen
make ci              # gates + govulncheck
```

Trivy-Image-Scan und SBOM-Erzeugung sind nach `LH-FA-BUILD-006` als
optionale Erweiterungen vorgesehen und folgen mit dem Release-Slice.

---

## 5. Architektur-Enforcement

`depguard` (Teil des SOLID-nahen Profils) erzwingt die Schicht-Regeln
aus `LH-FA-ARCH-003`. Detail in [`spec/architecture.md`](../../spec/architecture.md) §4.

Die `depguard`-Regelblöcke sind heute aktiv, matchen aber nichts,
solange `./internal/...` keinen produktiven Code enthält. Mit dem
ersten Paket pro Schicht greift die jeweilige Regel automatisch.

---

## 6. CI-Pipeline (GitHub Actions)

CI läuft auf GitHub Actions; Konfiguration in
[`.github/workflows/ci.yml`](../../.github/workflows/ci.yml).
Verbindliche Setzungen aus `LH-QA-003`, Begründung in
[ADR-0004](../plan/adr/0004-ci-system.md).

Pflichten:

- **Trigger:** `pull_request` und `push` auf `main`.
- **Jobs (beide PR-blockierend):**
  - `gates` — `make gates` (lint + test + coverage-gate).
  - `security-gates` — `make govulncheck`.
- **Runner:** `ubuntu-latest` mit vorinstalliertem Docker + BuildKit.
- **Keine Host-Toolchain:** Docker-only (`LH-FA-BUILD-007`); der
  Workflow installiert weder Go noch `golangci-lint` am Runner.
- **Actions SHA-gepinnt** mit Tag-Kommentar, z. B.
  `uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2`.
  Pin-Hebung ist Routine.
- **Permissions:** Top-Level `permissions: {}`; jeder Job lockert auf
  das Minimum (typisch `contents: read`).
- **Timeout:** `timeout-minutes: 20` pro Job.

Required-Status-Checks für `gates` und `security-gates` werden im
GitHub-UI nach dem ersten grünen Lauf gesetzt (Repository → Settings →
Branches → Branch protection rules → `main`).

Bewusst noch nicht enthalten (Folge-Slices, ADR-0004 Folgepunkte):
Image-Publish nach GHCR, Trivy-Image-Scan, Cluster-/Integrations-
Smoke gegen die echte Docker-Engine.
