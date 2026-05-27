# Architektur — u-boot

| Dokument    | Architektur-Spezifikation                                                                         |
| ----------- | ------------------------------------------------------------------------------------------------- |
| Projektname | `u-boot`                                                                                          |
| Bezug       | `LH-FA-ARCH-001..003` in [`spec/lastenheft.md`](lastenheft.md)                                    |
| ADR         | [`docs/plan/adr/0002-hexagonale-architektur.md`](../docs/plan/adr/0002-hexagonale-architektur.md) |
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

Pfeile zeigen die **Aufruf-/Datenfluss-Richtung** zur Laufzeit. Die **Import-Richtung** ist nicht überall identisch: `application` importiert nur Ports (Interfaces) und kennt die konkreten Adapter nicht; Dependency Injection findet im Wiring (`cmd/uboot/`) statt. Die innere Welt (`hexagon/`) kennt die äußere Welt (`adapter/`) **nicht** — das wird per `depguard` durchgesetzt (`LH-FA-ARCH-003`, siehe §4).

---

## 2. Schichten und Verzeichnisse

### 2.1 `hexagon/domain`

Reine Datentypen und invariantenhaltige Verhaltensregeln ohne I/O.

- **Aktuelle Inhalte (M3-T1):** `Project` (Aggregat mit `SchemaVersion`), `ProjectName` (Value-Object mit Regex aus `LH-FA-INIT-006`), `NormalizeProjectName` (deterministische Normalisierung nach `LH-FA-INIT-002`), `ErrInvalidProjectName`-Sentinel.
- **Geplante Erweiterungen:** `Service`, `Port`, `ImageRef`, `ComposeFile`, `EnvVar` (folgen mit M4/M5 als Add-on-Slices).
- **Erlaubte Imports:** ausschließlich Go-Standard-Library.
- **Verbotene Imports:** alle anderen `internal/`-Pakete, externe Libraries mit I/O.
- **Tests:** Unit-Tests mit `*_test.go` im selben Paket; pure Validierung ohne Mocks.

### 2.2 `hexagon/application`

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält keine externe I/O.

- **Aktuelle Inhalte (M3-T4):** `InitProjectService` orchestriert `FileSystem`/`YAMLCodec`/`Git`/`ProgressPort` zum LH-FA-INIT-001..007-Flow inklusive Re-Init-Pfaden nach LH-FA-INIT-005 (`--force`/`--backup`); Templates für die erzeugten Dateien via `embed.FS` + `text/template` (Templates leben unter `application/templates/*.tmpl`, die §611-strukturierten Configs wrappen ihren Inhalt in `BEGIN/END U-BOOT MANAGED BLOCK: init`-Marker); `ubootYAMLConfig`-Struct als Schema-Repräsentation für `u-boot.yaml` (LH-FA-CONF-002). Re-Init folgt einem strikten Plan-and-Execute-Split: `planFile` entscheidet pro Datei (`actionWrite`/`actionReplaceBlock`/`actionOverwriteFull`/Abort-Sentinel), Plan-Fehler verhindern jeden Side-Effect.
- **Hilfs-Pakete:** `application/managedblock/` (LH-SA-FILE-002-Marker-Parser: `Find`/`Has`/`Replace`, drei Comment-Styles Hash/HTMLComment/DoubleSlash, Sentinel `ErrBlockNotFound`/`ErrBlockMalformed`); `application/backup.go` mit `BackupPath` (kleinster-freier-Suffix-Algorithmus für `<path>.bak[.N]`, File + rekursive Verzeichnisse, TOCTOU-sicher via `WriteFileExclusive`/`Mkdir`, Rollback bei partiellem Tree-Copy, Mode- und Symlink-Reject per Lstat).
- **Geplante Erweiterungen:** `AddServiceService` (LH-FA-ADD-*), `RunDoctorService` (LH-FA-DIAG-*), `UpService`/`DownService` (LH-FA-UP-*), `GenerateService` (LH-FA-GEN-*).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driving`, `hexagon/port/driven` (zum Konsumieren von Driven-Ports und Implementieren von Driving-Ports).
- **Verbotene Imports:** `adapter/*`, externe I/O-Libraries.
- **Tests:** Unit-Tests mit Test-Doubles für Driven-Ports (Fakes oder Mocks in `_test.go`); Test-Library-Imports (z. B. `yaml.v3` für Fake-YAMLCodec) sind über den `*_test.go`-Carveout in `LH-FA-ARCH-003` erlaubt.

### 2.3 `hexagon/port/driving`

Interfaces, über die u-boot von außen angesprochen wird.

- **Aktuelle Inhalte (M3-T4):** `InitProjectUseCase` mit `InitProjectRequest` (`Name`/`BaseDir`/`SkipGit` aus T1–T3 plus `Force`/`Backup`/`AssumeExisting` aus T4) und `InitProjectResponse` (`Project`/`Created`/`Backups []BackupAction`). Sentinels für Re-Init und LH-FA-CLI-006-Mapping:
  - **Code 10 (Validierung):** `ErrProjectExists` (LH-FA-INIT-004 Marker u-boot.yaml/compose.yaml/.env.example), `ErrFileExists` (Non-Marker-Kollision), `ErrBaseDirMissing` (LH-AK-001), `ErrForceRequiresBackup` (LH-FA-INIT-005 §619), `ErrBackupUnsupportedKind` (Symlink-Reject).
  - **Code 14 (Technischer FS-Fehler):** `ErrBackupSourceMissing` (Race zwischen Caller-Check und Backup), `ErrBackupSuffixExhausted` (.bak[.0..999] alle belegt), `ErrBackupTooLarge` (Datei > MVP-Cap 256 MiB).
- **Geplante Erweiterungen:** `AddServiceUseCase`, `RemoveServiceUseCase`, `LifecycleUseCase`, `DoctorUseCase`, `GenerateUseCase`, `ConfigUseCase`.
- **Implementiert von:** Strukturen in `hexagon/application`.
- **Verwendet von:** `adapter/driving/*` (z. B. `cli/`).

### 2.4 `hexagon/port/driven`

Interfaces, über die `hexagon/application` externe Systeme nutzt.

- **Aktuelle Inhalte (M3-T4):**
  - `FileSystem` (`Exists`/`ReadFile`/`WriteFile`/`WriteFileExclusive`/`Mkdir`/`MkdirAll`/`Rename`/`ReadDir`/`Lstat`/`RemoveAll`). Die T4a-Erweiterungen folgen `os.*`-Konventionen: `Lstat` (no-follow für Symlink-Detection und Mode-Preservation), `WriteFileExclusive` (O_CREATE|O_EXCL für TOCTOU-sichere Backup-Slot-Reservierung), `Mkdir` (analog für Dir-Slots), `RemoveAll` (Rollback bei partiellem Tree-Copy).
  - `YAMLCodec` (`Marshal`/`Unmarshal`), `Git` (`IsRepository`/`Init`, jeweils mit `context.Context` als erstem Parameter), `Clock` (`Now`). Context-Konvention: nur Ports, deren Adapter blockieren können (Git via `os/exec`), nehmen Context; FS/YAML/Clock bleiben Context-frei (im Paket-Doc begründet).
  - `ProgressPort` (T4c-Review): `AffectedFiles(baseDir string, rows []AffectedFile)` zum strukturierten Reporting der LH-FA-INIT-005 §609 / LH-FA-CLI-005A §262 betroffenen Pfade vor jedem Re-Init-Write. `AffectedFile` trägt `Path`/`Action AffectedAction`/`Backup bool`; `AffectedAction` enumeriert `AffectedReplaceBlock`/`AffectedOverwriteFull`. Presentation lebt im Adapter; das Application-Paket bleibt I/O-Text-frei.
- **Aktuelle Inhalte (Forts., M4-Vorbereitung):** `Confirmer` (`ConfirmTreatAsExisting`, LH-FA-INIT-004 soft-existing-detection; siehe [`slice-m4-soft-existing-detection`](../docs/plan/planning/done/slice-m4-soft-existing-detection.md)) und `Logger` (LH-QA-004 Logging-Port mit `Debug`/`Info`/`Warn`/`Error`, slog-Adapter in `adapter/driven/logger`; siehe [`slice-m4-logging-port`](../docs/plan/planning/done/slice-m4-logging-port.md)).
- **Geplante Erweiterungen:** `DockerEngine` (`Up`/`Down`/`Ps`/`Logs`/`Exec`, M6).
- **Implementiert von:** Strukturen in `adapter/driven/*`.
- **Verwendet von:** `hexagon/application`.

### 2.5 `adapter/driving`

Konkrete Driver — Einstiegspunkte aus der Außenwelt.

- **Aktuelle Inhalte (M3-T4):** `cli/` mit Cobra v1.10.2 (siehe [`ADR-0005`](../docs/plan/adr/0005-cli-framework-cobra.md)). `init`-Subkommando mit `[name]`-Positional plus den Flags:
  - **Lokal:** `--no-git` (LH-FA-INIT-007), `--force` (LH-FA-INIT-005), `--backup` (LH-FA-INIT-005), `--assume-existing` (LH-FA-CLI-005A §238 — init-only).
  - **Persistent (Root):** `--yes`, `--no-interactive` (LH-FA-CLI-005A — gelten für alle bestätigungs-relevanten Subbefehle, heute init, künftig add/remove/config-set/down).
  Konflikt-Check `--yes` + `--no-interactive` → `ErrConflictingModeFlags` (CLI-internes Sentinel) → Exit-Code 2.
- `ExitCode(err)` bündelt die LH-FA-CLI-006-Klassifikation (0 / 2 / 10 / 14 / 1); `isValidationError` und `isFilesystemError` mappen die in §2.3 gelisteten Driving-Sentinels.
- `cli.App` mit Functional-Options-Pattern (`WithGetwd` als Test-Seam); `App.yes`/`noInteractive` halten die persistenten Flag-Werte, die `BoolVar` beim Re-Build der Root-Cobra pro `Execute` zurücksetzt — kein Flag-Leak zwischen Aufrufen.
- **Geplante Erweiterungen:** weitere Subkommandos (`add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `template`) folgen pro Use-Case-Slice; perspektivisch HTTP-/Daemon-Adapter (siehe [`slice-later-http-driving-adapter`](../docs/plan/planning/open/slice-later-http-driving-adapter.md)).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driving`, externe Libraries (z. B. Cobra).
- **Verbotene Imports:** `hexagon/application` und `adapter/driven`. Die Instanziierung von Application-Services und Driven-Adaptern erfolgt ausschließlich im Wiring (`cmd/uboot/`), das beide Welten zusammenfügt; der Driving-Adapter erhält fertig konstruierte Driving-Port-Implementierungen per Konstruktor.
- **Permanenter Carveout:** `contextcheck`-Ausnahme in `.golangci.yml`, weil Cobras `RunE`-Signatur (`func(cmd, args) error`) keinen Context-Parameter kennt — die Closure muss `cmd.Context()` extrahieren und an `runInit` durchreichen. Strikte Propagation passiert eine Ebene tiefer.

### 2.6 `adapter/driven`

Konkrete externe Adapter — Implementierungen der Driven-Ports.

- **Aktuelle Inhalte (M3-T4):** `fs/` (stdlib `os.*` — heute mit den T4a-Methoden `Lstat`/`WriteFileExclusive`/`Mkdir`/`RemoveAll` ergänzt; `WriteFileExclusive` nutzt `os.OpenFile` mit `O_CREATE|O_EXCL|O_WRONLY`), `yaml/` (`gopkg.in/yaml.v3`-Wrapper), `git/` (`os/exec git` mit `WithBinary`-Test-Override und Exit-Code-128-Klassifikation als „not a repo"), `clock/` (`time.Now()` in UTC), `progress/` (TextWriter rendert `driven.ProgressPort`-Events auf einen `io.Writer`; M3-T4c-Review). Jeder Adapter pinnt sein Port-Interface per `var _ driven.X = (*Adapter)(nil)` im Production-Code; Drift bricht den Package-Build.
- **Geplante Erweiterungen:** `docker/` (Docker-Engine via `os/exec docker compose`, M6), `logger/` (M4, `log/slog`-basiert), `progress/json` (für `LH-NFA-USE-004 --json` in V1).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driven`, externe Libraries.
- **Verbotene Imports:** `hexagon/application`, `adapter/driving`.
- **Test-Pfad:** `t.TempDir()` für FS, echte `git`-Binary via `os/exec.LookPath`-Skip (CI-Runner ohne git skippen sauber).

### 2.7 `cmd/uboot` — Wiring-Schicht

Einziger Ort, an dem `application` und `adapter` zusammen importiert werden.

- **Aktueller Inhalt (M3-T4):** `main.go` instantiiert `fs.New()`/`yaml.New()`/`git.New()`/`progress.NewText(stdout)`, baut den `InitProjectService` mit dem ProgressPort und übergibt ihn dem `cli.New(version, initSvc)`-Konstruktor. Plus signal-aware Context via `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)` und Error→Exit-Code-Mapping über `cli.ExitCode(err)`. Aktuelle Größe: ~60 Zeilen.
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

Die Regeln aus Abschnitt 3 werden im `lint`-Stage (`LH-FA-BUILD-001`) per `golangci-lint` mit dem `depguard`-Linter durchgesetzt. Konfiguration aktiv in [`.golangci.yml`](../.golangci.yml); das untenstehende Schema ist deckungsgleich mit der dortigen Konfiguration. Bei Änderungen müssen beide Quellen synchron gehalten werden.

Konventionen für jeden Regelblock:

- `list-mode: lax` — `deny`-only-Auswertung (Imports ohne `deny`-Treffer sind erlaubt, kein impliziter `allow`-Filter).
- `files` enthält als erste Pattern `!**/*_test.go`, um Tests vom Enforcement auszunehmen (Tests dürfen Fakes und Test-Libraries frei importieren; `LH-FA-ARCH-003`).
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

Mit M3 (T1..T4c) ist produktiver Code in allen sieben Schichten gelandet — jede `depguard`-Regel matcht jetzt mindestens ein Paket. Die Pro-Schicht-Verifikation ist mit M3-T5 abgeschlossen: `scripts/verify-depguard.sh` (Target `make verify-depguard`) injiziert pro Regel einen deklariert verbotenen Import, prüft `make lint` auf das erwartete `desc:` und revertiert. Slice-Plan: [`slice-m3-depguard-aktivierung-verifizieren`](../docs/plan/planning/done/slice-m3-depguard-aktivierung-verifizieren.md).

`//nolint:depguard`-Pragmas sind verboten. Carveouts werden zentral in `.golangci.yml` mit `desc` dokumentiert (`LH-FA-ARCH-003`).

---

## 5. Tests

- Unit-Tests stehen als `*_test.go` neben dem produktiven Code im selben Paket.
- **Domäne:** klassische Property/Value-Tests; keine Mocks nötig.
- **Application:** Fake-Implementierungen der Driven-Ports im `_test.go`-Paket; keine echte Docker-Engine.
- **Adapter (driven):** Integrationstests gegen echte Systeme, soweit lokal verfügbar (z. B. Docker-Engine für `adapter/driven/docker`). Ohne Docker-Engine werden diese Tests via Build-Tag (`//go:build docker`) ausgeschlossen. Build-Tag-Konvention:
  - Default ist *aus*: `make test` (Stage `test` im Dockerfile, `LH-FA-BUILD-001`) führt Tag-getaggte Tests nicht aus und bleibt damit auch ohne Docker-Socket grün.
  - Lokal mit verfügbarer Docker-Engine: `go test -tags docker ./...`.
  - In CI: ein separater Stage / ein separates Make-Target (Folge-Slice) aktiviert das Tag und mountet das Docker-Socket; dieser Pfad ist nicht Bestandteil von `make gates`, sondern ergänzt `make ci` als optionales Integrations-Smoketest-Ziel.
  - Pro Test-Datei mit dem entsprechenden Tag: erste Zeile `//go:build docker`, leere Zeile, dann `package …`.
- **Adapter (driving):** Tabellengetriebene Tests gegen den Driving-Port mit Fake-Application.
- Coverage-Messung (`LH-FA-BUILD-008`) bezieht sich auf `./internal/...`; `./cmd/...` ist ausgeschlossen.

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

Diese Architektur ist der Stand vom 2026-05-22. Änderungen erfolgen über neue ADRs, die das ADR-0002 superseden (`LH-FA-PROJDOCS-002`).

Geplante Erweiterungen, die im aktuellen Dokument noch nicht abgebildet sind (beide auch im Carveout-Inventar [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md) gelistet, `LH-FA-PROJDOCS-005`):

- HTTP-Driving-Adapter, falls u-boot perspektivisch eine Daemon-Variante bekommen soll.
- Plugin-System (`LH-OPEN-003`): voraussichtlich als zusätzlicher Driven-Port `PluginRegistry` mit dynamischer Adapter-Auswahl zur Laufzeit.
