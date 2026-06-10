# Architektur — u-boot

| Dokument    | Architektur-Spezifikation                                                                         |
| ----------- | ------------------------------------------------------------------------------------------------- |
| Projektname | `u-boot`                                                                                          |
| Bezug       | [`LH-FA-ARCH-001`](lastenheft.md#lh-fa-arch-001-hexagonales-pattern)..[`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) |
| Status      | Entwurf 0.1.0                                                                                     |
| Datum       | 2026-05-22                                                                                        |

---

## 1. Überblick

u-boot folgt dem **hexagonalen Architektur-Pattern** (auch: *Ports & Adapters*, Alistair Cockburn, 2005).

Sechs Schichten plus Wiring, klare Verantwortungen und einseitig gerichtete Abhängigkeiten:

```
            ┌──────────────────────────────────────────────────┐
            │                cmd/uboot (Wiring)                │
            │   (instanziiert Application + Adapter; main.go)  │
            └────────────────────┬─────────────────────────────┘
                                 │
        ┌────────────────────────┴─────────────────────────────┐
        ▼                                                      ▼
┌──────────────────┐                            ┌──────────────────────┐
│  adapter/driving │ → ruft AppService an  →    │  adapter/driven      │
│   (CLI-Commands) │                            │  (Docker, FS, YAML)  │
└──────────────────┘                            └──────────────────────┘
        │                                                      ▲
        │   ruft via Port-Interface                            │   wird via Port-Interface
        ▼                                                      │   aus Application gerufen
┌──────────────────┐    ┌─────────────────────────┐   ┌────────┴────────────┐
│ hexagon/         │    │ hexagon/                │   │ hexagon/            │
│   port/driving   │ →  │   application           │ → │   port/driven       │
└──────────────────┘    │   (Use-Cases)           │   └─────────────────────┘
                        └────────────┬────────────┘
                                     │
                                     ▼
                          ┌──────────────────────┐
                          │ hexagon/domain       │
                          │ (reine Datentypen)   │
                          └──────────────────────┘
```

Pfeile zeigen die **Aufruf-/Datenfluss-Richtung** zur Laufzeit. Die **Import-Richtung** ist nicht überall identisch: `application` importiert nur Ports (Interfaces) und kennt die konkreten Adapter nicht; Dependency Injection findet im Wiring (`cmd/uboot/`) statt. Die innere Welt (`hexagon/`) kennt die äußere Welt (`adapter/`) **nicht** — das wird per `depguard` durchgesetzt ([`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement), siehe §4).

---

## 2. Schichten und Verzeichnisse

### 2.1 `hexagon/domain`

Reine Datentypen und invariantenhaltige Verhaltensregeln ohne I/O.

- **Inhalt:** `Project` (Aggregat mit `SchemaVersion`), `ProjectName` (Value-Object mit Regex aus [`LH-FA-INIT-006`](lastenheft.md#lh-fa-init-006-projektnamen-validierung)), `NormalizeProjectName` (deterministische Normalisierung nach [`LH-FA-INIT-002`](lastenheft.md#lh-fa-init-002-projektname)), `ErrInvalidProjectName`-Sentinel; `ServiceName` (Value-Object für Add-on-Identifier mit eigener Regex, Sentinel `ErrInvalidServiceName`) und `ServiceState`-Enum (Active/Deactivated/EnabledUnset/Unregistered/InconsistentYAML/InconsistentBlock) für die [`LH-FA-ADD-005`](lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern)-State-Machine; `DiagnosticReport` mit `Severity`-Enum (`SeverityOK`/`SeverityWarn`/`SeverityError`) und `Diagnostic{ID, Severity, Message, Hint}` für die Doctor-Use-Cases ([`LH-FA-DIAG-003`](lastenheft.md#lh-fa-diag-003-fehlerklassifikation)).
- **Vorgesehene Erweiterungen:** `Service`, `Port`, `ImageRef`, `ComposeFile`, `EnvVar` für Add-on-Use-Cases.
- **Erlaubte Imports:** ausschließlich Go-Standard-Library.
- **Verbotene Imports:** alle anderen `internal/`-Pakete, externe Libraries mit I/O.
- **Tests:** Unit-Tests mit `*_test.go` im selben Paket; pure Validierung ohne Mocks.

### 2.2 `hexagon/application`

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält keine externe I/O.

- **Inhalt:**
  - `InitProjectService` orchestriert `FileSystem`/`YAMLCodec`/`Git`/`ProgressPort`/`Confirmer`/`Logger` zum [`LH-FA-INIT-001`](lastenheft.md#lh-fa-init-001-neues-projekt-initialisieren)..[`LH-FA-INIT-007`](lastenheft.md#lh-fa-init-007-git-repository-initialisierung)-Flow inklusive Re-Init-Pfaden nach [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005-überschreibschutz) (`--force`/`--backup`) und [`LH-FA-INIT-004`](lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) Soft-Existing-Detection. Templates für die erzeugten Dateien via `embed.FS` + `text/template` (Templates leben unter `application/templates/*.tmpl`; die §611-strukturierten Configs wrappen ihren Inhalt in `BEGIN/END U-BOOT MANAGED BLOCK: init`-Marker). `ubootYAMLConfig`-Struct als Schema-Repräsentation für `u-boot.yaml` ([`LH-FA-CONF-002`](lastenheft.md#lh-fa-conf-002-inhalt-der-konfiguration)). Re-Init folgt einem strikten Plan-and-Execute-Split: `planFile` entscheidet pro Datei (`actionWrite`/`actionReplaceBlock`/`actionOverwriteFull`/Abort-Sentinel), Plan-Fehler verhindern jeden Side-Effect.
  - `DoctorService` orchestriert `FileSystem`/`Git`/`DockerProbe`/`Logger` zu den [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen)-Checks (write-permissions, git availability, docker installed/reachable, compose installed, später u-boot.yaml/compose.yaml-Validierung, Devcontainer-Konsistenz). Stdlib-Semver-Min-Check (`parseSemverMajorMinor` + `classifyVersionAtLeast`) für die Mindestversionen 24.0 (Docker) / 2.20 (Compose). Service ist severity-agnostisch — failures sind `SeverityError`-Diagnostics im Report, kein Go-error.
- **Hilfs-Pakete:** `application/managedblock/` ([`LH-SA-FILE-002`](lastenheft.md#lh-sa-file-002-markierte-verwaltete-bereiche)-Marker-Parser: `Find`/`Has`/`Replace`, drei Comment-Styles Hash/HTMLComment/DoubleSlash, Sentinel `ErrBlockNotFound`/`ErrBlockMalformed`); `application/backup.go` mit `BackupPath` (kleinster-freier-Suffix-Algorithmus für `<path>.bak[.N]`, File + rekursive Verzeichnisse, TOCTOU-sicher via `WriteFileExclusive`/`Mkdir`, Rollback bei partiellem Tree-Copy, Mode- und Symlink-Reject per Lstat, Streaming-Copy via `FileSystem.Copy`/`CopyExclusive`).
- **Vorgesehene Erweiterungen:** `AddServiceService` (LH-FA-ADD-*), `UpService`/`DownService` (LH-FA-UP-*), `GenerateService` (LH-FA-GEN-*).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driving`, `hexagon/port/driven` (zum Konsumieren von Driven-Ports und Implementieren von Driving-Ports).
- **Verbotene Imports:** `adapter/*`, externe I/O-Libraries.
- **Tests:** Unit-Tests mit Test-Doubles für Driven-Ports (Fakes oder Mocks in `_test.go`); Test-Library-Imports (z. B. `yaml.v3` für Fake-YAMLCodec) sind über den `*_test.go`-Carveout in [`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) erlaubt.

### 2.3 `hexagon/port/driving`

Interfaces, über die u-boot von außen angesprochen wird.

- **Inhalt:**
  - `InitProjectUseCase` mit `InitProjectRequest` (`Name`/`BaseDir`/`SkipGit`/`Force`/`Backup`/`AssumeExisting`/`NoInteractive`) und `InitProjectResponse` (`Project`/`Created`/`Backups []BackupAction`).
  - `DoctorUseCase` mit `DoctorRequest` (`BaseDir`) und `DoctorResponse` (`Report domain.DiagnosticReport`). Per Kontrakt liefert `Check` immer einen Report; check-failures sind `SeverityError`-Diagnostics, kein Go-error. Severity-Klassifikation + Exit-Code-Mapping (`--strict`) übernimmt der CLI-Adapter.
  - `AddServiceUseCase` mit `AddServiceRequest` (`BaseDir`/`ServiceName`) und `AddServiceResponse` (`ServiceName`/`PriorState`/`State`/`Changed []string`). Idempotenz-garantiert: Zweit-Add mit gleichen Args ist no-op + nil-error (`PriorState=Active`, `Changed=nil`).
- **Sentinels** für die [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006-exit-codes)-Exit-Code-Klassifikation (liegen im `driving`-Paket statt im `application`-Paket, damit der CLI-Adapter via `errors.Is` auf sie verzweigt, ohne `application` zu importieren — [`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)):
  - **Code 10 (Validierung):** `ErrProjectExists` ([`LH-FA-INIT-004`](lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) Marker u-boot.yaml/compose.yaml/.env.example), `ErrFileExists` (Non-Marker-Kollision), `ErrBaseDirMissing` ([`LH-AK-001`](lastenheft.md#lh-ak-001-minimaler-init-flow) oder leeres `BaseDir`-Feld; geteilt zwischen `InitProjectUseCase` und `DoctorUseCase`), `ErrForceRequiresBackup` ([`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005-überschreibschutz) §619), `ErrBackupUnsupportedKind` (Symlink-Reject), `ErrProjectNotInitialized` ([`LH-FA-ADD-001`](lastenheft.md#lh-fa-add-001-add-on-befehl) — kein/unparsbares u-boot.yaml), `ErrServiceUnsupported` ([`LH-FA-ADD-002`](lastenheft.md#lh-fa-add-002-postgresql-hinzufügen) — ServiceName syntaktisch valide aber nicht im built-in catalog), `ErrServiceInconsistent` ([`LH-FA-ADD-005`](lastenheft.md#lh-fa-add-005-mehrfaches-hinzufügen-verhindern) §896 — orphan compose-block ohne YAML-Anker). Plus die `domain`-Validierungs-Sentinels `ErrInvalidProjectName` und `ErrInvalidServiceName`.
  - **Code 14 (Technischer FS-Fehler):** `ErrBackupSourceMissing` (Race zwischen Caller-Check und Backup), `ErrBackupSuffixExhausted` (.bak[.0..999] alle belegt).
- **Vorgesehene Erweiterungen:** `RemoveServiceUseCase`, `LifecycleUseCase` (Up/Down), `GenerateUseCase`, `ConfigUseCase`.
- **Implementiert von:** Strukturen in `hexagon/application`.
- **Verwendet von:** `adapter/driving/*` (z. B. `cli/`).

### 2.4 `hexagon/port/driven`

Interfaces, über die `hexagon/application` externe Systeme nutzt.

- **Inhalt:**
  - `FileSystem` (`Exists`/`ReadFile`/`WriteFile`/`WriteFileExclusive`/`Mkdir`/`MkdirAll`/`Rename`/`ReadDir`/`Lstat`/`RemoveAll`/`Copy`/`CopyExclusive`). Folgt `os.*`-Konventionen: `Lstat` (no-follow für Symlink-Detection und Mode-Preservation), `WriteFileExclusive` (O_CREATE|O_EXCL für TOCTOU-sichere Backup-Slot-Reservierung), `Mkdir` (analog für Dir-Slots), `RemoveAll` (Rollback bei partiellem Tree-Copy), `Copy`/`CopyExclusive` (Streaming-Backup via `io.Copy` ohne RAM-Cap).
  - `YAMLCodec` (`Marshal`/`Unmarshal`).
  - `Git` (`IsRepository`/`Init`/`Version`) — alle mit `context.Context` als erstem Parameter (Adapter shellt zum `git`-Binary, das blockieren kann). `Version` liefert die bare semver (Adapter strippt das `git version `-Prefix).
  - `Clock` (`Now`) — ohne Context (Implementierung non-blocking).
  - **Context-Konvention:** nur Ports, deren Adapter blockieren können (Git, Docker via `os/exec`), nehmen Context; FS/YAML/Clock bleiben Context-frei (im Paket-Doc begründet).
  - `ProgressPort` (`AffectedFiles(baseDir, rows)`) zum strukturierten Reporting der [`LH-FA-INIT-005`](lastenheft.md#lh-fa-init-005-überschreibschutz) §609 / [`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung) §262 betroffenen Pfade vor jedem Re-Init-Write. `AffectedFile` trägt `Path`/`Action AffectedAction`/`Backup bool`; `AffectedAction` enumeriert `AffectedReplaceBlock`/`AffectedOverwriteFull`. Presentation lebt im Adapter.
  - `Confirmer` (`ConfirmTreatAsExisting(ctx, baseDir, indicators)`) für die [`LH-FA-INIT-004`](lastenheft.md#lh-fa-init-004-bestehendes-projekt-erkennen) Soft-Existing-Detection-Prompts. Narrowly scoped per Confirm-Kontext.
  - `Logger` (`Debug`/`Info`/`Warn`/`Error`, slog-konforme `...any`-Variadic) als [`LH-QA-004`](lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil)-Logging-Port.
  - `DockerProbe` (`Version`/`Info`/`ComposeVersion`) für die read-only [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen)-Probes (`docker version --format`, `docker compose version --short`). Bewusst getrennt vom state-mutierenden `DockerEngine` (siehe Erweiterungen). Backend-Annahme: ein Docker-API-kompatibles Engine-Binary auf `$PATH`. Heute Docker; Podman ≥ 4.0 funktioniert als Drop-in (`docker → podman`-Symlink + `DOCKER_HOST` auf `podman.socket`) — die Version-Klassifikation pinnt die Docker-Mindestwerte (24.0 / 2.20), Podman-Versionen werden vorerst als `Severity: warn` ("unrecognized version") emittiert, kein Exit-Code-Eskalation. Ein dedizierter Podman-Probe-Pfad ist ein eigener Slice (Trigger: erster konkreter Bedarf).
  - `RuntimeEnvironment` (`InContainer() bool`) für die best-effort Container-Self-Detection via `/.dockerenv` (Docker) / `/run/.containerenv` (Podman/CRI-O/buildah). Treibt das `doctor`-Skip-Verhalten für die vier Host-Prerequisite-Checks im distroless-Container-Run, ohne dass die Adapter im Container fehlschlagen müssten.
- **Vorgesehene Erweiterungen:** `DockerEngine` (`Up`/`Down`/`Ps`/`Logs`/`Exec`) für die Compose-Lifecycle-Operationen — explizit getrennt von `DockerProbe`, weil state-mutierend.
- **Implementiert von:** Strukturen in `adapter/driven/*`.
- **Verwendet von:** `hexagon/application`.

### 2.5 `adapter/driving`

Konkrete Driver — Einstiegspunkte aus der Außenwelt.

- **Inhalt:** `cli/` mit Cobra. Pro Subkommando ein eigenes Cobra-Command in einer eigenen Datei.
  - Lokale Flags pro Subkommando (z. B. `init`: `--no-git`/`--force`/`--backup`/`--assume-existing`).
  - **Persistente Flags (Root):** `--yes`, `--no-interactive` ([`LH-FA-CLI-005A`](lastenheft.md#lh-fa-cli-005a-interaktivität-und-automatisierung) — gelten für alle bestätigungs-relevanten Subbefehle). Konflikt-Check `--yes` + `--no-interactive` → `ErrConflictingModeFlags` (CLI-internes Sentinel) → Exit-Code 2.
  - `ExitCode(err)` bündelt die [`LH-FA-CLI-006`](lastenheft.md#lh-fa-cli-006-exit-codes)-Klassifikation (0 / 2 / 10 / 14 / 1); `isValidationError` und `isFilesystemError` mappen die in §2.3 gelisteten Driving-Sentinels.
  - `cli.App` mit Functional-Options-Pattern (`WithGetwd` als Test-Seam); persistente Flag-Werte werden beim Re-Build der Root-Cobra pro `Execute` zurückgesetzt — kein Flag-Leak zwischen Aufrufen.
- **Vorgesehene Erweiterungen:** weitere Subkommandos (`add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `template`). Ein HTTP-/Daemon-Adapter ist nicht vorgesehen; u-boot bleibt CLI-only (siehe §7).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driving`, externe Libraries (z. B. Cobra).
- **Verbotene Imports:** `hexagon/application` und `adapter/driven`. Die Instanziierung von Application-Services und Driven-Adaptern erfolgt ausschließlich im Wiring (`cmd/uboot/`), das beide Welten zusammenfügt; der Driving-Adapter erhält fertig konstruierte Driving-Port-Implementierungen per Konstruktor.
- **Permanenter Carveout:** `contextcheck`-Ausnahme in `.golangci.yml`, weil Cobras `RunE`-Signatur (`func(cmd, args) error`) keinen Context-Parameter kennt — die Closure muss `cmd.Context()` extrahieren und an `runInit` durchreichen. Strikte Propagation passiert eine Ebene tiefer.

### 2.6 `adapter/driven`

Konkrete externe Adapter — Implementierungen der Driven-Ports.

- **Inhalt:**
  - `fs/` — `FileSystem`-Adapter via stdlib `os.*` (`os.ReadFile`/`WriteFile`/`MkdirAll`/`Rename`/`ReadDir`/`Lstat`/`RemoveAll`; `WriteFileExclusive` mit `O_CREATE|O_EXCL|O_WRONLY`; Streaming `Copy`/`CopyExclusive` über `os.Open` + `io.Copy`).
  - `yaml/` — `YAMLCodec`-Adapter via `gopkg.in/yaml.v3`.
  - `git/` — `Git`-Adapter via `os/exec git` mit `WithBinary`-Test-Override und Exit-Code-128-Klassifikation als „not a repo".
  - `clock/` — `Clock`-Adapter via `time.Now()` in UTC.
  - `progress/` — `ProgressPort`-Adapter (TextWriter rendert Events auf einen `io.Writer`).
  - `confirm/` — `Confirmer`-Adapter (`bufio.Scanner` über stdin, Prompt auf stderr, Default `[y/N]`).
  - `logger/` — `Logger`-Adapter via `log/slog` (Text + JSON-Format konfigurierbar).
  - `docker/` — `DockerProbe`-Adapter via `os/exec docker` (read-only diagnostics für [LH-FA-DIAG-002](lastenheft.md#lh-fa-diag-002-lokale-voraussetzungen-prüfen)).
  - Jeder Adapter pinnt sein Port-Interface per `var _ driven.X = (*Adapter)(nil)` im Production-Code; Drift bricht den Package-Build.
- **Vorgesehene Erweiterungen:** `docker/`-Erweiterung um den `DockerEngine`-Adapter (`Up`/`Down`/`Ps`/`Logs`/`Exec` via `docker compose`); `progress/json` für [`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe) `--json`.
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driven`, externe Libraries.
- **Verbotene Imports:** `hexagon/application`, `adapter/driving`.
- **Test-Pfad:** `t.TempDir()` für FS, echte `git`-Binary via `os/exec.LookPath`-Skip (CI-Runner ohne git skippen sauber).

### 2.7 `cmd/uboot` — Wiring-Schicht

Einziger Ort, an dem `application` und `adapter` zusammen importiert werden.

- **Inhalt:** `main.go` instantiiert die Driven-Adapter (`fs.New()`, `yaml.New()`, `git.New()`, `progress.NewText(stdout)`, `confirm.New(os.Stdin, stderr)`, `logger.New(stderr, ...)`, `docker.New()`), konstruiert die Application-Services (`InitProjectService`, `DoctorService`) mit den nötigen Ports und übergibt sie dem `cli.New(version, ...)`-Konstruktor. Plus signal-aware Context via `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` und Error→Exit-Code-Mapping über `cli.ExitCode(err)`.
- Hält keine Geschäftslogik. So klein wie sinnvoll möglich (Größenordnung 150–300 Zeilen `main.go` plus ein paar kleine Wiring-Helper); ab dieser Marke ist eine Aufteilung in mehrere Wiring-Pakete (`internal/wire/<feature>/`) zu erwägen.

---

## 3. Import-Regeln

| Schicht                | darf importieren                                                              | darf nicht importieren                                     |
| ---------------------- | ----------------------------------------------------------------------------- | ---------------------------------------------------------- |
| `hexagon/domain`       | Go-Standard-Library                                                           | alle anderen `internal/`-Pakete, I/O-Libraries             |
| `hexagon/application`  | `hexagon/domain`, `hexagon/port/driving`, `hexagon/port/driven`               | `adapter/*`, externe I/O-Libraries                         |
| `hexagon/port/driving` | `hexagon/domain`                                                              | `hexagon/application`, `hexagon/port/driven`, `adapter/*`  |
| `hexagon/port/driven`  | `hexagon/domain`                                                              | `hexagon/application`, `hexagon/port/driving`, `adapter/*` |
| `adapter/driving`      | `hexagon/domain`, `hexagon/port/driving`, externe Libraries (z. B. Cobra)     | `hexagon/application`, `adapter/driven`                    |
| `adapter/driven`       | `hexagon/domain`, `hexagon/port/driven`, externe Libraries (z. B. Docker-SDK) | `hexagon/application`, `adapter/driving`                   |
| `cmd/uboot`            | `internal/...`, Standardbibliothek, externe Libraries                         | (frei — Wiring-Schicht)                                    |

Begründung der Regeln:

- **Domain isoliert** halten ⇒ Domänenobjekte sind portabel, testbar, frei von Framework-Annahmen.
- **Application kennt nur Ports** ⇒ Use-Cases sind ohne reale Docker-Engine testbar (Fake-`DockerEngine`).
- **Ports sind kreuz-blind** (`driving` ↔ `driven`) ⇒ vermeidet versteckte Kopplungen über das Application.
- **Wiring in cmd/** ⇒ Austausch eines Adapters (z. B. von `os/exec` Docker auf das Docker-SDK) erfordert keine Änderung in `application` oder `hexagon/port`.

---

## 4. Enforcement via `golangci-lint depguard`

Die Regeln aus Abschnitt 3 werden im `lint`-Stage ([`LH-FA-BUILD-001`](lastenheft.md#lh-fa-build-001-multi-stage-dockerfile-u-boot-repo)) per `golangci-lint` mit dem `depguard`-Linter durchgesetzt. Konfiguration aktiv in [`.golangci.yml`](../.golangci.yml); das untenstehende Schema ist deckungsgleich mit der dortigen Konfiguration. Bei Änderungen müssen beide Quellen synchron gehalten werden.

Konventionen für jeden Regelblock:

- `list-mode: lax` — `deny`-only-Auswertung (Imports ohne `deny`-Treffer sind erlaubt, kein impliziter `allow`-Filter).
- `files` enthält als erste Pattern `!**/*_test.go`, um Tests vom Enforcement auszunehmen (Tests dürfen Fakes und Test-Libraries frei importieren; [`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)).
- `deny`-Einträge nennen den blockierten Modul-Pfad und in `desc` die LH-Kennung als Begründung.

```yaml
linters:
  enable:
    - depguard

  settings:
    depguard:
      rules:
        domain-isoliert:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/hexagon/domain/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: domain must not depend on application (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port
              desc: domain must not depend on port (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: domain must not depend on adapter (LH-FA-ARCH-003)

        application-no-adapter:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/hexagon/application/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: application must depend on ports, not on adapter implementations (LH-FA-ARCH-003)

        port-no-application:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/hexagon/port/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: port must not depend on application (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: port must not depend on adapter (LH-FA-ARCH-003)

        port-driving-no-driven:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/hexagon/port/driving/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port/driven
              desc: driving port must not depend on driven port (LH-FA-ARCH-003)

        port-driven-no-driving:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/hexagon/port/driven/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port/driving
              desc: driven port must not depend on driving port (LH-FA-ARCH-003)

        adapter-no-application:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/adapter/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: adapter must implement ports, not consume application (LH-FA-ARCH-003)

        adapter-driving-no-driven:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/adapter/driving/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter/driven
              desc: driving adapter must not depend on driven adapter — wire via cmd/uboot (LH-FA-ARCH-003)

        adapter-driven-no-driving:
          list-mode: lax
          files:
            - '!**/*_test.go'
            - '**/internal/adapter/driven/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter/driving
              desc: driven adapter must not depend on driving adapter (LH-FA-ARCH-003)
```

Jede `depguard`-Regel matcht mindestens ein Paket im Produktiv-Code; die Pro-Schicht-Verifikation läuft per `scripts/verify-depguard.sh` (Target `make verify-depguard`), das pro Regel einen deklariert verbotenen Import injiziert, `make lint` auf das erwartete `desc:` prüft und die Stub-Datei revertiert.

`//nolint:depguard`-Pragmas sind verboten. Carveouts werden zentral in `.golangci.yml` mit `desc` dokumentiert ([`LH-FA-ARCH-003`](lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)).

---

## 5. Tests

- Unit-Tests stehen als `*_test.go` neben dem produktiven Code im selben Paket.
- **Domäne:** klassische Property/Value-Tests; keine Mocks nötig.
- **Application:** Fake-Implementierungen der Driven-Ports im `_test.go`-Paket; keine echte Docker-Engine.
- **Adapter (driven):** Integrationstests gegen echte Systeme, soweit lokal verfügbar (z. B. Docker-Engine für `adapter/driven/docker`). Ohne Docker-Engine werden diese Tests via Build-Tag (`//go:build docker`) ausgeschlossen. Build-Tag-Konvention:
  - Default ist *aus*: `make test` (Stage `test` im Dockerfile, [`LH-FA-BUILD-001`](lastenheft.md#lh-fa-build-001-multi-stage-dockerfile-u-boot-repo)) führt Tag-getaggte Tests nicht aus und bleibt damit auch ohne Docker-Socket grün.
  - Lokal mit verfügbarer Docker-Engine: `go test -tags docker ./...`.
  - In CI: ein separater Stage / ein separates Make-Target (Folge-Slice) aktiviert das Tag und mountet das Docker-Socket; dieser Pfad ist nicht Bestandteil von `make gates`, sondern ergänzt `make ci` als optionales Integrations-Smoketest-Ziel.
  - Pro Test-Datei mit dem entsprechenden Tag: erste Zeile `//go:build docker`, leere Zeile, dann `package …`.
- **Adapter (driving):** Tabellengetriebene Tests gegen den Driving-Port mit Fake-Application.
- Coverage-Messung ([`LH-FA-BUILD-008`](lastenheft.md#lh-fa-build-008-coverage-bootstrap)) bezieht sich auf `./internal/...`; `./cmd/...` ist ausgeschlossen.

---

## 6. Anti-Patterns

Die folgenden Muster sind verboten und werden im Review abgelehnt:

- **God-Service:** ein `application`-Service, der alle Use-Cases bündelt. Stattdessen ein Service pro Use-Case-Familie.
- **Anämische Domäne:** Domänentypen ohne Verhalten, die nur Daten halten. Domänen-Invarianten gehören in die Domäne.
- **Adapter ruft Adapter:** `adapter/driving` importiert `adapter/driven` direkt. Wiring gehört in `cmd/uboot`.
- **Port importiert Application:** zyklische Abhängigkeit, verbietet sich aus Architektur und ist `depguard`-blockiert.
- **`//nolint:depguard`** zur Umgehung einer Schicht-Regel. Es gibt keinen legitimen Carveout im Fachcode; wenn eine Regel im Weg steht, gehört die Schicht-Definition korrigiert.
- **Externe Library im `domain`-Paket** (`yaml.v3`, Docker-SDK, Cobra, …). Domäne bleibt I/O-frei.

---

## 7. Evolution

Diese Architektur ist der Stand vom 2026-05-22. Änderungen erfolgen über
neue Architekturentscheidungen und anschließende Spec-Nachführung
([`LH-FA-PROJDOCS-002`](lastenheft.md#lh-fa-projdocs-002-adr-format)).

Geplante Erweiterungen, die im aktuellen Dokument noch nicht abgebildet sind: keine.

**Nicht** geplant:

- HTTP-Driving-Adapter (Daemon-Variante): u-boot bleibt CLI-only;
  Maschinen-Schnittstellen laufen über `--json`/`--dry-run`-Flags
  ([`LH-NFA-USE-004`](lastenheft.md#lh-nfa-use-004-maschinenlesbare-ausgabe)).
- Plugin-System ([`LH-OPEN-003`](lastenheft.md#lh-open-003-plugin-system-entschieden)):
  das Add-on-System bleibt statisch; kein `PluginRegistry`-Driven-Port.
