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
aus §1.2 + `depguard` aus `LH-FA-ARCH-003`) lebt in
[`.golangci.yml`](../../.golangci.yml). `//nolint`-Suppressions bleiben
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

Tests laufen Docker-basiert über die `test`-Stage:

```bash
make test
```

Die Stage führt `go test ./...` aus. Im MVP-Bootstrap deckt der
Stub-Code nur `cmd/uboot` ab (CLI-Flag-Parsing, Exit-Codes).

Mit dem ersten produktiven Slice (M3) folgen:

- Unit-Tests in `internal/hexagon/{domain,application}` mit
  Fake-Implementierungen der Driven-Ports (`spec/architecture.md` §5).
- Integrationstests in `internal/adapter/driven/{docker,fs,yaml}` mit
  Build-Tag-getriggerten Pfaden (`//go:build docker` etc., siehe
  `spec/architecture.md` §5).

---

## 3. Coverage

Coverage-Messung über die `coverage`-Stage, bootstrap-aware
(`LH-FA-BUILD-008`):

```bash
make coverage-gate
```

- Solange `./internal/...` keinen produktiven Code enthält, läuft der
  Gate im Bootstrap-Modus mit Schwellwert `0`.
- Mit dem ersten produktiven Paket wird der Schwellwert in einem
  Folge-Commit angehoben (Empfehlung: 80 % als erster Wert, langfristig
  90 % analog m-trace/k-deskflight).
- Override pro Aufruf: `make coverage-gate THRESHOLD=80`.

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
