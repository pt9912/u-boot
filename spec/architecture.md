# Architektur — u-boot

| Dokument         | Architektur-Spezifikation                                     |
| ---------------- | -------------------------------------------------------------- |
| Projektname      | `u-boot`                                                       |
| Bezug            | `LH-FA-ARCH-001..003` in [`spec/lastenheft.md`](lastenheft.md) |
| ADR              | [`docs/plan/adr/0002-hexagonale-architektur.md`](../docs/plan/adr/0002-hexagonale-architektur.md) |
| Status           | Entwurf 0.1.0                                                  |
| Datum            | 2026-05-21                                                     |

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

Reine Datentypen und invariantes Verhalten ohne I/O.

- **Beispielinhalte:** `Project`, `Service`, `ComposeFile`, `EnvVar`, Value-Objects (`ProjectName`, `Port`, `ImageRef`), Validierungs-Methoden.
- **Erlaubte Imports:** ausschließlich Go-Standard-Library.
- **Verbotene Imports:** alle anderen `internal/`-Pakete, externe Libraries mit I/O.
- **Tests:** Unit-Tests mit `*_test.go` im selben Paket.

### 2.2 `hexagon/application`

Anwendungslogik (Use-Cases). Orchestriert Domäne und Ports, enthält keine externe I/O.

- **Beispielinhalte:** `InitProjectService`, `AddServiceService`, `RunDoctorService`, `RenderTemplateService`.
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port` (zum Konsumieren von Driven-Ports und Implementieren von Driving-Ports).
- **Verbotene Imports:** `adapter/*`, externe I/O-Libraries.
- **Tests:** Unit-Tests mit Test-Doubles für Driven-Ports (Fakes oder Mocks in `_test.go`).

### 2.3 `hexagon/port/driving`

Interfaces, über die u-boot von außen angesprochen wird.

- **Beispiel:** `InitProjectUseCase`, `AddServiceUseCase` — werden vom CLI-Adapter aufgerufen.
- **Implementiert von:** Strukturen in `hexagon/application`.
- **Verwendet von:** `adapter/driving/*` (z. B. `cli/`).

### 2.4 `hexagon/port/driven`

Interfaces, über die `hexagon/application` externe Systeme nutzt.

- **Beispiel:** `DockerEngine` (`Up`, `Down`, `Ps`), `FileSystem` (`ReadFile`, `WriteFile`, `Exists`), `YAMLCodec` (`Marshal`, `Unmarshal`), `Clock`.
- **Implementiert von:** Strukturen in `adapter/driven/*`.
- **Verwendet von:** `hexagon/application`.

### 2.5 `adapter/driving`

Konkrete Driver — Einstiegspunkte aus der Außenwelt.

- **Beispielinhalte:** `cli/` (Cobra-Commands `init`, `add`, `up`, `doctor`, …), perspektivisch ggf. ein HTTP- oder Daemon-Adapter.
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driving`, externe Libraries (z. B. Cobra).
- **Verbotene Imports:** `hexagon/application` und `adapter/driven`. Die Instanziierung von Application-Services und Driven-Adaptern erfolgt ausschließlich im Wiring (`cmd/uboot/`), das beide Welten zusammenfügt; der Driving-Adapter erhält fertig konstruierte Driving-Port-Implementierungen per Konstruktor.

### 2.6 `adapter/driven`

Konkrete externe Adapter — Implementierungen der Driven-Ports.

- **Beispielinhalte:** `docker/` (via `os/exec` oder Docker-SDK gegen die Docker Engine), `fs/` (Dateisystem-IO), `yaml/` (YAML-Codec via `gopkg.in/yaml.v3`), `clock/` (Real-Time, Test-Stub in Tests).
- **Erlaubte Imports:** `hexagon/domain`, `hexagon/port/driven`, externe Libraries.
- **Verbotene Imports:** `hexagon/application`, `adapter/driving`.

### 2.7 `cmd/uboot` — Wiring-Schicht

Einziger Ort, an dem `application` und `adapter` zusammen importiert werden.

- Erzeugt konkrete Driven-Adapter (`fsAdapter`, `dockerAdapter`, `yamlCodec`), injiziert sie in `application`-Services, registriert diese als Driving-Ports im CLI-Adapter.
- Hält keine Geschäftslogik. So klein wie sinnvoll möglich (Größenordnung 150–300 Zeilen `main.go` plus ein paar kleine Wiring-Helper); ab dieser Marke ist eine Aufteilung in mehrere Wiring-Pakete (`internal/wire/<feature>/`) zu erwägen.

---

## 3. Import-Regeln

| Schicht                  | darf importieren                                                                          | darf nicht importieren                                                |
| ------------------------ | ----------------------------------------------------------------------------------------- | --------------------------------------------------------------------- |
| `hexagon/domain`         | Go-Standard-Library                                                                       | alle anderen `internal/`-Pakete, I/O-Libraries                        |
| `hexagon/application`    | `hexagon/domain`, `hexagon/port`                                                          | `adapter/*`, externe I/O-Libraries                                    |
| `hexagon/port/driving`   | `hexagon/domain`                                                                          | `hexagon/application`, `hexagon/port/driven`, `adapter/*`             |
| `hexagon/port/driven`    | `hexagon/domain`                                                                          | `hexagon/application`, `hexagon/port/driving`, `adapter/*`            |
| `adapter/driving`        | `hexagon/domain`, `hexagon/port/driving`, externe Libraries (z. B. Cobra)                | `hexagon/application`, `adapter/driven`                               |
| `adapter/driven`         | `hexagon/domain`, `hexagon/port/driven`, externe Libraries (z. B. Docker-SDK)            | `hexagon/application`, `adapter/driving`                              |
| `cmd/uboot`              | `internal/...`, Standardbibliothek, externe Libraries                                         | (frei — Wiring-Schicht)                                               |

Begründung der Regeln:

- **Domain isoliert** halten ⇒ Domänenobjekte sind portabel, testbar, frei von Framework-Annahmen.
- **Application kennt nur Ports** ⇒ Use-Cases sind ohne reale Docker-Engine testbar (Fake-`DockerEngine`).
- **Ports sind kreuz-blind** (`driving` ↔ `driven`) ⇒ vermeidet versteckte Kopplungen über das Application.
- **Wiring in cmd/** ⇒ Austausch eines Adapters (z. B. von `os/exec` Docker auf das Docker-SDK) erfordert keine Änderung in `application` oder `hexagon/port`.

---

## 4. Enforcement via `golangci-lint depguard`

Die Regeln aus Abschnitt 3 werden im `lint`-Stage (`LH-FA-BUILD-001`) per `golangci-lint` mit dem `depguard`-Linter durchgesetzt. Konfiguration in [`.golangci.yml`](../.golangci.yml).

Schema (Beispiel — wird mit dem ersten produktiven Paket scharf geschaltet):

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
            - '**/internal/hexagon/domain/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: domain must not depend on application (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port
              desc: domain must not depend on port (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: domain must not depend on adapter (LH-FA-ARCH-003)

        application-no-adapter:
          files:
            - '**/internal/hexagon/application/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: application must depend on ports, not on adapter implementations (LH-FA-ARCH-003)

        port-no-application:
          files:
            - '**/internal/hexagon/port/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: port must not depend on application (LH-FA-ARCH-003)
            - pkg: github.com/pt9912/u-boot/internal/adapter
              desc: port must not depend on adapter (LH-FA-ARCH-003)

        port-driving-no-driven:
          files:
            - '**/internal/hexagon/port/driving/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port/driven
              desc: driving port must not depend on driven port (LH-FA-ARCH-003)

        port-driven-no-driving:
          files:
            - '**/internal/hexagon/port/driven/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/port/driving
              desc: driven port must not depend on driving port (LH-FA-ARCH-003)

        adapter-no-application:
          files:
            - '**/internal/adapter/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/hexagon/application
              desc: adapter must implement ports, not consume application (LH-FA-ARCH-003)

        adapter-driving-no-driven:
          files:
            - '**/internal/adapter/driving/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter/driven
              desc: driving adapter must not depend on driven adapter — wire via cmd/uboot (LH-FA-ARCH-003)

        adapter-driven-no-driving:
          files:
            - '**/internal/adapter/driven/**'
          deny:
            - pkg: github.com/pt9912/u-boot/internal/adapter/driving
              desc: driven adapter must not depend on driving adapter (LH-FA-ARCH-003)
```

`//nolint:depguard`-Pragmas sind verboten. Carveouts werden zentral in `.golangci.yml` mit `Why:`-Kommentar dokumentiert.

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

Diese Architektur ist der Stand vom 2026-05-21. Änderungen erfolgen über neue ADRs, die das ADR-0002 superseden (`LH-FA-PROJDOCS-002`).

Geplante Erweiterungen, die im aktuellen Dokument noch nicht abgebildet sind:

- HTTP-Driving-Adapter, falls u-boot perspektivisch eine Daemon-Variante bekommen soll.
- Plugin-System (`LH-OPEN-003`): voraussichtlich als zusätzlicher Driven-Port `PluginRegistry` mit dynamischer Adapter-Auswahl zur Laufzeit.
