# Spezifikation — u-boot

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
[`lastenheft.md`](lastenheft.md) formulierten Anforderungen (`LH-*`-Kennungen). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — verbindlich, aber ohne Änderung des Lastenhefts
fortschreibbar. Sie verweist nur aufwärts auf das Lastenheft, nie auf ADRs, Slices,
Carveouts oder die Roadmap.

**Kennungen:** Eine Verfeinerung einer einzelnen Anforderung trägt deren Kennung mit
Buchstabensuffix (`<Anforderungs-Kennung>.a`, `.b`, …). Alles, was keine einzelne
Anforderung verfeinert (Schemata, Defaults, Fehler-Codes, Metrik-Felder, externe
Verträge), trägt `SPEC-<NNN>`: dreistellig, fortlaufend je Datei, nicht je Abschnitt.
Eine `SPEC-*` ist eine Adresse, keine Anforderung.

---

## 1. Algorithmen und Datenflüsse

Verfeinerungen einzelner Anforderungen: Ablauf, Zustandsfolgen, Entscheidungsregeln.
Tatsächlicher Code gehört nicht hierher.

### LH-FA-INIT-002.a — Normalisierung des abgeleiteten Projektnamens

Verfeinert [`LH-FA-INIT-002`](lastenheft.md#lh-fa-init-002--projektname).

1. Der Basisname des Arbeitsverzeichnisses wird auf Kleinbuchstaben gesetzt.
2. Alle Zeichen außer `a-z`, `0-9` und `-` werden auf `-` abgebildet.
3. aufeinanderfolgende `-` werden zu einem einzelnen `-` zusammengeführt.
4. führende und nachgestellte `-` sowie Leerzeichen werden entfernt.
5. Die Länge wird auf 1 bis 63 Zeichen begrenzt.
6. Nach Kürzung auf 63 Zeichen wird erneut auf führende/nachgestellte `-` geprüft und diese notfalls entfernt.
7. Anschließend wird der Name gegen die Validierung in [`LH-FA-INIT-006`](lastenheft.md#lh-fa-init-006--projektnamen-validierung) geprüft.

### LH-FA-DEV-004.a — Übergabe der Benutzer-ID an den Image-Build

Verfeinert [`LH-FA-DEV-004`](lastenheft.md#lh-fa-dev-004--benutzerrechte).

- Der Wert wird als Build-Argument `USER_UID` an den Image-Build übergeben; der Container-Benutzer wird mit dieser UID angelegt.

### LH-FA-DIAG-002.a — Prüfung von Docker und Docker-kompatiblen Drop-ins

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen).

- Docker installiert (Mindestversion: 24.0.0 oder neuer) **oder**
  ein Docker-API-kompatibler Drop-in (z. B. Podman ≥ 4.0 mit
  aktivem `podman.socket` und `DOCKER_HOST` darauf gezeigt). Die
  Checks `docker.installed` / `docker.reachable` /
  `docker.compose.installed` shellen aus zum `docker`-Binary;
  Drop-ins, die nicht den Docker-Version-Format-Pin
  (`<major>.<minor>.<patch>`-Bereich 24.0 / 2.20) treffen,
  werden als `Severity: warn` („unrecognized version") emittiert,
  ohne den Exit-Code zu eskalieren. Eine formal getestete
  Podman-Variante folgt in einem eigenen Slice bei konkretem
  Bedarf — die heutige MVP-Pflicht ist Docker.

### LH-FA-BUILD-001.a — Stages des Multi-Stage-Dockerfiles

Verfeinert [`LH-FA-BUILD-001`](lastenheft.md#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo).

- BuildKit-Direktive in der ersten Zeile: `# syntax=docker/dockerfile:1.7`.
- Pflicht-Stages im MVP:
  - `deps` – Modulauflösung (`go mod download`) als Cache-Layer.
  - `compile` – schnelles Compile-Feedback (`go build`) ohne Tests/Lint.
  - `test` – `go test ./...`.
  - `lint` – `golangci-lint run ./...`.
  - `coverage` – `go test -coverprofile` + Coverage-Gate gegen `COVERAGE_THRESHOLD`; Bootstrap-Verhalten und Scope sind in [`LH-FA-BUILD-008`](lastenheft.md#lh-fa-build-008--coverage-bootstrap) und [`LH-FA-BUILD-009`](lastenheft.md#lh-fa-build-009--repository-layout) definiert.
  - `build` – statisch gelinktes Binary (`CGO_ENABLED=0`, `-ldflags="-s -w"`).
  - `runtime` – minimales Endimage ([`LH-FA-BUILD-002`](lastenheft.md#lh-fa-build-002--runtime-stage-pflichten)).
- Jeder Stage ist ein eigenständiges Build-Ziel und wird per `docker build --target <stage>` einzeln baubar.

### LH-FA-BUILD-002.a — Eigenschaften der Runtime-Stage

Verfeinert [`LH-FA-BUILD-002`](lastenheft.md#lh-fa-build-002--runtime-stage-pflichten).

- Base-Image: `gcr.io/distroless/static-debian12:nonroot` (oder gleichwertig minimal und ohne Shell).
- Non-root-Benutzer (Distroless-`nonroot`-User, `USER 65532:65532`).
- `ENTRYPOINT` zeigt auf das im `build`-Stage erzeugte Binary; der konkrete Pfad (Empfehlung: `/usr/local/bin/u-boot`) ist im Dockerfile dokumentiert.
- OCI Image Labels gesetzt:
  - `org.opencontainers.image.source`
  - `org.opencontainers.image.description`
  - `org.opencontainers.image.licenses`
  - `org.opencontainers.image.title`
- Keine Build-Toolchain im Endimage; alle Build-Artefakte werden aus dem `build`-Stage per `COPY --from=build` übernommen.

### LH-QA-003.a — Komposition der CI-Pipeline (GitHub Actions)

Verfeinert [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions).

- Workflow-Datei: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- Trigger: `pull_request` und `push` auf den Branch `main`.
- Drei Jobs, parallel, alle PR-blockierend (Required-Status-Checks im GitHub-UI nach dem ersten grünen Lauf zu konfigurieren; die Required-Status-Check-Liste muss die tatsächlichen Workflow-`name:`-Felder verwenden, nicht die Kurz-Keys). Die drei delegierten Make-Targets sind in [`LH-FA-BUILD-005`](lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets) (MVP) und [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets) (V1) definiert; die PR-Blocking-**Pflicht** für alle drei kommt aus diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)-Eintrag (MVP). Damit ist `make govulncheck` und `make image-scan` über [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions) MVP-bindend, auch wenn die jeweiligen Make-Target-Definitionen unter [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets) (V1) liegen:
  - `gates (lint + test + coverage-gate)` — führt `make gates` aus (lint + test + coverage-gate, [`LH-FA-BUILD-005`](lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets)/[`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets)).
  - `security-gates (govulncheck)` — führt `make govulncheck` aus ([`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets), MVP-Pflicht via diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)).
  - `image-scan (trivy HIGH+CRITICAL)` — führt `make image-scan` aus (Trivy gegen das Runtime-Image, severity `HIGH,CRITICAL`, exit-code `1`; [`LH-FA-BUILD-006`](lastenheft.md#lh-fa-build-006--aggregator-targets), MVP-Pflicht via diesem [`LH-QA-003`](lastenheft.md#lh-qa-003--ci-fähigkeit-github-actions)).
- Runner: `ubuntu-latest`. Keine Host-Go-Toolchain ([`LH-FA-BUILD-007`](lastenheft.md#lh-fa-build-007--docker-only-workflow)); der Runner braucht nur das vorinstallierte Docker + BuildKit.
- Actions sind **SHA-gepinnt** mit Tag-Kommentar (Supply-Chain-Härtung gegen Tag-Move). Pin-Hebung ist Routine; neuer Commit-SHA via `gh api repos/<owner>/<repo>/git/refs/tags/<tag>`.
- Top-Level `permissions: {}` (alle Tokens entzogen); jeder Job lockert auf das Minimum (Defense-in-Depth).
- Jeder Job mit `timeout-minutes` versehen (Empfehlung: 20).
- Die konkrete Workflow-Ausprägung muss die hier genannten
  Anforderungen und die Build-Target-Verträge erfüllen.

### LH-QA-004.a — Zusammensetzung des Lint-Profils

Verfeinert [`LH-QA-004`](lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil).

Profil-Komposition (29 Linter insgesamt):

- 5 Default-Linter (`govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`).
- 24 SOLID-nahe Zusatz-Linter (Komplexitäts-, Funktionslänge-, Interface-, Kopplungs- und Boundary-Signale). **`depguard`** für die Schicht-Regeln aus [`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) ist Teil dieser 24.

## 2. Datenstrukturen und Schemas

Formate und Schemata (`u-boot.yaml`, JSON-Ausgabe, CLI-Tabellen). Jede Struktur trägt
eine `SPEC-<NNN>`.

### SPEC-001 — JSON-Schema der Vorschau-Ausgabe (`--dry-run --json`)

Gilt für `--dry-run --json` und `--diff --json`; die Pflichtfelder nennt die Anforderung.

Für `--dry-run --json` ist die Ausgabe mindestens wie folgt als maschinenlesbares JSON zu liefern:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"],
  "properties": {
    "subcommand": {
      "type": "string",
      "description": "Unterkommando bei gruppierten Hauptkommandos wie `template` oder `config`"
    },
    "status": {
      "type": "string",
      "enum": ["ok", "warn", "error"]
    },
    "command": {
      "type": "string",
      "enum": ["init", "add", "remove", "up", "down", "doctor", "logs", "generate", "config", "template"]
    },
    "dryRun": {
      "type": "boolean"
    },
    "diff": {
      "type": "boolean"
    },
    "plannedFiles": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "action"],
        "properties": {
          "path": { "type": "string" },
          "action": {
            "type": "string",
            "enum": ["create", "modify", "delete"]
          }
        },
        "additionalProperties": true
      }
    },
    "changes": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "count"],
        "properties": {
          "path": { "type": "string" },
          "count": { "type": "integer", "minimum": 0 }
        },
        "additionalProperties": true
      }
    },
    "diagnostics": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["level", "code", "message"],
        "properties": {
          "level": { "type": "string", "enum": ["warn", "error"] },
          "code": { "type": "string" },
          "message": { "type": "string" },
          "file": { "type": "string" }
        },
        "additionalProperties": true
      }
    },
    "exitCode": {
      "type": "integer",
      "minimum": 0
    }
  },
  "allOf": [
    {
      "if": {
        "properties": {
          "command": { "const": "template" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    },
    {
      "if": {
        "properties": {
          "command": { "const": "config" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    }
  ],
  "additionalProperties": true
}
```

### SPEC-002 — Beispielinstanz einer Vorschau-Ausgabe (`add --dry-run --json`)

Beispielinstanz:

```json
{
  "status": "warn",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "create" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 12 }
  ],
  "diagnostics": [
    { "level": "warn", "code": "LH-FA-CLI-007", "message": "Geplante Datei fehlt bereits" }
  ],
  "exitCode": 0
}
```

### SPEC-003 — Beispielinstanz einer Diff-Ausgabe (`add --diff --json`, ohne `--dry-run`)

Beispiel für `--diff --json` ohne `--dry-run` (Vorschau mit anschließendem Schreiben):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": false,
  "diff": true,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "modify" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 6 }
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

### SPEC-004 — Projektnamen-Muster

Gilt für den Projektnamen (`project.name`) bei `init` und `config set`.

- regulärer Ausdruck: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`

### SPEC-005 — Schema der Projektkonfiguration `u-boot.yaml`

Gilt als Mindestinhalt der Konfigurationsdatei; die Hinweise zu Schlüsseln und Werten nennt die Anforderung.

```yaml
schemaVersion: 1
project:
  name: my-service

services:
  postgres:
    enabled: false

devcontainer:
  enabled: false

# Optionale, V1-relevante Felder:
# services:
#   keycloak:
#     enabled: false
#     persistence: embedded   # embedded | external-postgres
#   otel:
#     enabled: false
#
# devcontainer:
#   featureSources:
#     allow:
#       - https://ghcr.io/devcontainers/features/node
#   user:
#     uid: 1000                    # 1..65535
#   profile: default               # default | sandbox
#   sandbox:
#     nestedRuntime: none          # none | podman
#     onUnavailable: warn          # warn | fail
#     repository: ''               # Clone-Quelle statt origin (URL)
#     egress:
#       enabled: false
#       allow: []
```

### SPEC-008 — Eigenschaften und Pflicht-Targets des Makefiles

Pflicht-Eigenschaften:

- `.DEFAULT_GOAL := help`
- `.PHONY` für alle Targets gesetzt
- `help`-Target mit Übersicht über alle verfügbaren Targets
- Variablen mit `?=`-Defaults für Overridability (`IMAGE`, `GO_VERSION`, `GOLANGCI_LINT_VERSION`, `THRESHOLD`)

MVP-Pflicht-Targets:

| Target          | Zweck                                                           |
| --------------- | --------------------------------------------------------------- |
| `help`          | Übersicht aller Targets                                         |
| `deps`          | `docker build --target deps`                                    |
| `compile`       | `docker build --target compile`                                 |
| `lint`          | `docker build --target lint`                                    |
| `test`          | `docker build --target test`                                    |
| `coverage`      | Alias auf `coverage-gate`                                       |
| `coverage-gate` | `docker build --target coverage --build-arg COVERAGE_THRESHOLD` |
| `build`         | `docker build --target runtime`                                 |
| `run`           | `docker run --rm <image> --help` (Smoketest); Dependency: `build` |
| `clean`         | lokale Artefakte und gebaute Images entfernen                   |

### SPEC-009 — Mindestlayout des u-boot-Repositories

Mindestlayout:

```text
.
├── cmd/
│   └── uboot/
│       └── main.go              # Entry point der CLI (Wiring-Schicht)
├── internal/                    # nicht-exportierbare Implementierung
│   ├── hexagon/                 # innere Schichten (LH-FA-ARCH-002)
│   │   ├── domain/
│   │   ├── application/
│   │   └── port/{driving,driven}/
│   └── adapter/                 # äußere Schichten (LH-FA-ARCH-002)
│       ├── driving/
│       └── driven/
├── spec/                        # Lastenheft, weitere Spezifikationen
│   ├── lastenheft.md
│   └── <weitere-spezifikation>.md
├── docs/                        # Doku-Struktur (LH-FA-PROJDOCS-001)
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
├── .dockerignore
├── .gitignore
├── LICENSE
└── README.md
```

### SPEC-010 — Mindest-Verzeichnisstruktur unter `docs/`

```text
docs/
├── archive/                  # abgelöste oder ersetzte Inhalte
├── plan/
│   ├── adr/                  # Architecture Decision Records
│   └── planning/
│       ├── open/             # Backlog
│       ├── next/             # priorisiert für nächsten Schritt
│       ├── in-progress/      # aktiv bearbeitet
│       └── done/             # abgeschlossen
└── user/                     # User-facing Dokumentation
```

### SPEC-011 — Format der Architecture Decision Records

- Dateiname beginnt mit vierstelliger Nummer, beginnend bei `0001` und monoton steigend: `0001-<slug>.md`, `0002-<slug>.md`; Slug in Kebab-Case (z. B. `0001-implementierungssprache-go.md`).
- Dokumenttitel als `#`-Überschrift: `# ADR <Nr>: <Titel>`.
- Direkt darunter die Kopf-Felder als **fette Inline-Felder** (nicht als `##`-Überschriften):
  - `**Status:**` – einer aus `Proposed`, `Accepted`, `Deprecated`, `Superseded by <NNNN>-<slug>`.
  - `**Datum:**` – Entscheidungsdatum im Format `YYYY-MM-DD`.
  - `**Autor:**` – verantwortliche Rolle oder Person.
  - `**Bezug:**` – betroffene `LH-*`- und ggf. Vorgänger-ADR-IDs als Markdown-Links (optional, wenn zutreffend).
  - `**Schärft:**` – welche Spec-Stelle (`architecture.md §N`) diese ADR verbindlich macht, als Aufwärts-Deklaration der Änderungskopplung (wer die ADR ändert, zieht von hier die Spec-Stellen nach); `—`, wenn Prozess-ADR ohne Spec-Bezug.
- Danach die Abschnitte, jeweils als `##`-Überschrift, in dieser Reihenfolge:
  1. `## Kontext` – Ausgangslage, auslösende Anforderung, tragende Annahmen.
  2. `## Entscheidung` – die Wahl, eindeutig.
  3. `## Verglichene Alternativen` – Optionen mit Pro/Contra (auch „nichts tun").
  4. `## Konsequenzen` – kurz- und langfristige Folgen, positiv und negativ, inkl. Folgepflichten.
  5. `## Fitness Function` – die maschinell prüfbare Regel, falls die Entscheidung sich in einer Code-Eigenschaft niederschlägt (sonst entfällt der Abschnitt).
  6. `## Re-Evaluierungs-Trigger` – wann die Entscheidung erneut zu prüfen ist.
  7. `## Geschichte` – Tabelle Datum/Ereignis/Verweis (`Proposed`, `Accepted`, …).
- ADR-Nummern werden nie wiederverwendet; abgelöste ADRs bleiben mit Status `Superseded by <NNNN>-<slug>` erhalten und verweisen auf den Nachfolger über den vollen Dateinamen-Stamm (ohne `.md`), als klickbaren Link.

### SPEC-012 — Schichten und Verzeichnislayout unter `internal/`

```text
internal/
├── hexagon/
│   ├── domain/         # reine Datentypen + invariantes Verhalten, keine I/O
│   ├── application/    # Use-Cases; ruft ausschließlich Ports auf
│   └── port/
│       ├── driving/    # Interfaces, die von außen (CLI/HTTP) konsumiert werden
│       └── driven/     # Interfaces, die das Application nach außen ruft
└── adapter/
    ├── driving/        # konkrete Driver (z. B. cli/ mit Cobra-Commands)
    └── driven/         # konkrete Adapter (z. B. docker/, fs/, yaml/)
```

### SPEC-013 — Import-Regel-Tabelle der Schichten

| Schicht | darf importieren | darf nicht importieren |
| --- | --- | --- |
| `hexagon/domain` | Go-Standard-Library | alle anderen `internal/`-Pakete, I/O-Libraries |
| `hexagon/application` | `hexagon/domain`, `hexagon/port/driving`, `hexagon/port/driven` | `adapter/*`, externe I/O-Libraries |
| `hexagon/port/driving` | `hexagon/domain` | `hexagon/application`, `hexagon/port/driven`, `adapter/*` |
| `hexagon/port/driven` | `hexagon/domain` | `hexagon/application`, `hexagon/port/driving`, `adapter/*` |
| `adapter/driving` | `hexagon/domain`, `hexagon/port/driving`, externe Libraries | `hexagon/application`, `adapter/driven` |
| `adapter/driven` | `hexagon/domain`, `hexagon/port/driven`, externe Libraries | `hexagon/application`, `adapter/driving` |
| `cmd/uboot` | `internal/...`, Standardbibliothek, externe Libraries | keine Einschränkung; Wiring-Schicht |

### SPEC-014 — Markierungsformat verwalteter Bereiche je Dateityp

- YAML, `.env`, `Dockerfile`, Shell-Skripte (`#`-Kommentare):

  ```yaml
  # BEGIN U-BOOT MANAGED BLOCK: postgres
  # ...
  # END U-BOOT MANAGED BLOCK: postgres
  ```

- Markdown (`README.md`, `CHANGELOG.md`) als HTML-Kommentar:

  ```markdown
  <!-- BEGIN U-BOOT MANAGED BLOCK: postgres -->
  ...
  <!-- END U-BOOT MANAGED BLOCK: postgres -->
  ```

- JSONC (z. B. `.devcontainer/devcontainer.json`):

  ```jsonc
  // BEGIN U-BOOT MANAGED BLOCK: postgres
  // ...
  // END U-BOOT MANAGED BLOCK: postgres
  ```

### SPEC-015 — Beispielinstanz einer Minimalkontrakt-Ausgabe (`doctor --json`)

Beispiel:

```json
{
  "status": "ok",
  "command": "doctor",
  "diagnostics": [],
  "exitCode": 0
}
```

## 3. Defaults und Konstanten

Werte, die im Produkt fest sind (Standardwerte, Grenzwerte, Versionsuntergrenzen).

### SPEC-006 — Build-Args des Dockerfiles und Pin-Politik

Gilt für das Dockerfile im Repo-Root; Overrides per `docker build --build-arg <NAME>=<value>` bzw. Makefile-Variable.

- `ARG GO_VERSION` – mit Default-Pin (z. B. `1.26.3`); Hebung ist Routine ohne separaten Spec-Eintrag.
- `ARG GOLANGCI_LINT_VERSION` – mit Default-Pin; gleiche Pin-Politik.
- `ARG COVERAGE_THRESHOLD` – mit Default `0` (bootstrap) und Override-Pfad `make coverage-gate THRESHOLD=…`.

### SPEC-007 — Mindest-Ausschlüsse der `.dockerignore`

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

### SPEC-016 — Go-Toolchain: Mindestversion und Dockerfile-Pin

Stand zum Entscheidungsdatum 2026-05-21; der aktuelle Pin steht im Dockerfile (`ARG GO_VERSION`).

Mindest-Toolchain: Go 1.26 oder neuer (`go 1.26.0` in `go.mod`, analog Referenzprojekt `k-deskflight`); Default-Pin im Dockerfile als `ARG GO_VERSION` (aktuell `1.26.3`, die aktuelle Stable-Version am Entscheidungsdatum). Pin-Hebung ist Routine ohne separaten Spec-Eintrag.

### SPEC-018 — Standardwerte der Build- und Laufzeitumgebung

| Name | Wert | Quelle |
|---|---|---|
| Coverage-Schwellwert (`COVERAGE_THRESHOLD` / `THRESHOLD`) | 90 Prozent | Dockerfile-Build-Arg, Makefile-Variable |
| Polling-Intervall von `u-boot up` | 500 ms | Stabilisierungsschleife (`docker compose ps`) |
| Standard-Wartezeit von `u-boot up` | 60 s | `--timeout` |
| UID des Container-Benutzers | 1000 | `devcontainer.user.uid` |

### SPEC-019 — Default-Allowlist der Egress-Restriktion

Die wirksame Liste ist die Vereinigung aus Basis, Ergänzungen und Nutzerliste, sortiert und ohne Duplikate.

| Quelle | Hosts |
|---|---|
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` |
| `nestedRuntime: podman` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io` |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Clone-Quelle | Host von `devcontainer.sandbox.repository` beziehungsweise `origin` |
| Nutzer | Einträge von `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema, Port und Wildcard) |

## 4. Fehler-Codes und Logging-Felder

Verbindliche Diagnose- und Fehler-Codes sowie Logging-Felder.

### SPEC-017 — Diagnose-Codes der `doctor`-Prüfungen (Code-Registry)

Jede Prüfung von `u-boot doctor` trägt einen stabilen Code in `diagnostics[].code`. Die Tabelle ist die kanonische Registry; der Quelltext (`DefaultAllowedCodes`) und diese Tabelle werden von einem Test symmetrisch abgeglichen. Andere Befehle tragen als Code die Kennung der verursachenden Anforderung.

<!-- code-registry:start -->

| Code | Bedeutung |
| --- | --- |
| `fs.write-permissions` | Schreib-Permission im Working Directory |
| `git.installed` | Git-Binary verfügbar |
| `docker.installed` | Docker-Binary verfügbar |
| `docker.reachable` | Docker-Daemon erreichbar |
| `docker.compose.installed` | Compose-Plugin verfügbar |
| `uboot.yaml.valid` | `u-boot.yaml` syntaktisch valide |
| `compose.yaml.valid` | `compose.yaml` syntaktisch valide |
| `devcontainer.json.valid` | `.devcontainer/devcontainer.json` syntaktisch valide |
| `devcontainer.dockerfile.valid` | `.devcontainer/Dockerfile` parsebar |
| `services.enabled-key` | `u-boot.yaml` `services`-Block konsistent |
| `devcontainer.forwardPorts.consistency` | `devcontainer.json` `forwardPorts` konsistent |
| `devcontainer.features.allowlist` | `devcontainer` Features auf Allowlist |
| `devcontainer.features.drift` | `devcontainer` Features ohne Drift |
| `devcontainer.sandbox.runtime` | Sandbox: nested Podman (`/dev/fuse`, Profil-Konsistenz), [`LH-FA-DEV-007`](lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) |
| `devcontainer.sandbox.egress` | Sandbox: Egress-Restriktion konsistent konfiguriert, [`LH-FA-DEV-008`](lastenheft.md#lh-fa-dev-008--egress-restriktion) |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Git-Zugangsdaten, Token-Quelle, [`LH-FA-DEV-009`](lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |

<!-- code-registry:end -->

## 5. Metriken und Tracing-Felder

Verbindliche Telemetrie-Felder pro Span.

### SPEC-020 — Telemetrie des Produkts

`u-boot` erzeugt selbst keine Metriken oder Traces und hat keine Tracing-Felder. Der Add-on `otel` richtet lediglich einen OpenTelemetry-Collector für das Zielprojekt ein (siehe Add-on-Katalog); dessen Standardports stehen dort.

## 6. Externe Verträge

Schnittstellen zu Drittsystemen mit Versionsannahme (Docker, Compose, Devcontainer).

### SPEC-021 — Docker, Docker Compose und Podman

| System | Mindestversion | Hinweis |
|---|---|---|
| Docker Engine | 24.0.0 | Prüfung durch `u-boot doctor` |
| Docker Compose (Plugin) | 2.20.0 | Prüfung durch `u-boot doctor` |
| Podman (als Docker-Ersatz) | 4.0 | mit aktivem `podman.socket` und `DOCKER_HOST` auf den Socket; nicht erkannte Versionen werden als Warnung gemeldet |

### SPEC-022 — Add-on-Katalog: Images und Ports

| Add-on | Image | Host-Ports |
|---|---|---|
| `postgres` | `postgres:16-alpine` | 5432 |
| `keycloak` | `quay.io/keycloak/keycloak:26.0` | 8080 |
| `otel` | `otel/opentelemetry-collector:0.108.0` | 4317 (OTLP/gRPC), 4318 (OTLP/HTTP) |

### SPEC-023 — Devcontainer: Basis-Image und Feature-Quellen

- Basis-Image des erzeugten Dockerfiles: `mcr.microsoft.com/devcontainers/base:debian`.
- Eingebaute Features: `git`, `docker-cli`, `node`, `java`, `go`, `cpp`, `kubectl-helm`, `postgres-client`, jeweils als `ghcr.io/devcontainers/features/<name>:<version>` mit der Standardversion `1`; jede andere Quelle gilt als extern und braucht eine Freigabe in `devcontainer.featureSources.allow`.
- Externe Feature-Quellen müssen als URL mit Schema `http://`, `https://` oder `oci://` angegeben werden.

## 7. Historie

| Datum | Änderung |
|---|---|
| 2026-10-02 | Angelegt als Technik-Stratum; Schemata, Beispielinstanzen, Algorithmen, Build-/CI-, Doku- und Architektur-Details sowie Markierungsformate wörtlich aus dem Lastenheft übernommen (`SPEC-001`..`SPEC-016`, Verfeinerungen); neu aus dem Bestand von Code und Konfiguration: Standardwerte, Egress-Default-Allowlist, Code-Registry der `doctor`-Prüfungen, Telemetrie-Aussage, externe Verträge (`SPEC-017`..`SPEC-023`). |
