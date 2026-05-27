# u-boot Carveouts

Master-Inventar aller bewussten Carveouts in der u-boot-Codebase
(`LH-FA-PROJDOCS-005`). Wird laufend gepflegt und liegt deshalb dauerhaft
in `in-progress/`. Die Konvention ist zusätzlich als persistente
Claude-Memory `feedback-carveouts-need-plans` hinterlegt, damit
künftige Sessions sie nicht vergessen.

Spalten:

- **Carveout** — wo (Datei/Sektion) und was (kurze Beschreibung).
- **Status** — `temporär` (mit Aufhebungsplan) oder `permanent`
  (begründet, kein Plan).
- **Plan / Begründung** — bei `temporär`: Verweis auf den Slice-Plan in
  `docs/plan/planning/{open,next,in-progress}/`. Bei `permanent`: kurze
  Begründung.

## Temporäre Carveouts (Plan-Pflicht)

| Carveout | Status | Plan / Begründung |
| -------- | ------ | ----------------- |
| ADR-0004 Folgepunkt: Image-Publish nach GHCR und Trivy-Image-Scan fehlen (`.github/workflows/ci.yml` enthält sie nicht) | temporär | [`open/slice-v1-release-pipeline.md`](../open/slice-v1-release-pipeline.md) |
| Build-Tag-Pfad `//go:build docker` für Adapter-Integrationstests (`spec/architecture.md` §5) ist nur beschrieben, kein CI-Pfad und kein Adapter-Test existiert | temporär | [`open/slice-m6-docker-integrationstests.md`](../open/slice-m6-docker-integrationstests.md) |
| `LH-OPEN-002` Paketierung ist offen (`spec/lastenheft.md` §14) — Distributionswege (Binary-Release, Homebrew, Debian/RPM, npm/pip) sind nicht festgelegt | temporär | [`open/slice-v1-release-pipeline.md`](../open/slice-v1-release-pipeline.md) (GHCR-Anteil); weitere Distributionswege bekommen eigene Slices beim ersten konkreten Bedarf |
| `LH-OPEN-003` Plugin-System ist offen (`spec/lastenheft.md` §14, auch `spec/architecture.md` §7 als „geplante Erweiterung") — keine Entscheidung zwischen fest-eingebauten Add-ons und nachladbaren Plugins | temporär | [`open/slice-v1-plugin-system-entscheidung.md`](../open/slice-v1-plugin-system-entscheidung.md) |
| `LH-OPEN-004` Template-Format ist offen (`spec/lastenheft.md` §14) — YAML+Dateien vs. Cookiecutter vs. eigenes Format vs. OCI-Pakete | temporär | [`open/slice-v1-template-format-entscheidung.md`](../open/slice-v1-template-format-entscheidung.md) |
| HTTP-Driving-Adapter ist als „geplante Erweiterung" in `spec/architecture.md` §7 erwähnt, aber nicht spezifiziert oder gefordert | temporär | [`open/slice-later-http-driving-adapter.md`](../open/slice-later-http-driving-adapter.md) |
| `maxBackupFileSize = 256 << 20` als harter Cap in `internal/hexagon/application/backup.go` (LH-FA-INIT-005) — Backup lädt heute via `ReadFile`+`WriteFile` ins RAM, multi-GB-Assets würden OOM erzeugen; Cap zieht `driving.ErrBackupTooLarge` (Exit-Code 14) | temporär | [`open/slice-v1-backup-streaming-copy.md`](../open/slice-v1-backup-streaming-copy.md) |

## Permanente Carveouts (kein Plan, im Inventar dokumentiert)

| Carveout | Status | Plan / Begründung |
| -------- | ------ | ----------------- |
| `errcheck.exclude-functions` für `fmt.Fprintln`/`Fprintf`/`Fprint` (`.golangci.yml`) | permanent | CLI-Writes auf stdout/stderr können nicht meaningful fehlschlagen; `_, _ =`-Prefix bringt keinen Wert. |
| `testpackage`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | Wiring-Schicht braucht `package main`; externe `_test`-Packages sind dort nicht erzwingbar. |
| `gochecknoglobals`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | `var version` wird per `-ldflags="-X main.version=…"` überschrieben — das ist der kanonische Go-Pattern für Build-Metadaten. |
| Test-Carveouts in `_test.go` (`cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam`, `revive(unused-parameter)`) | permanent | Tabellengetriebene Tests und Fakes erzeugen legitim hohe Komplexität / fehlenden Context; Profil-Schwellen passen für Production-Code, nicht für Tests. |
| `!**/*_test.go` als erste files-Pattern in jedem `depguard`-Regelblock (`.golangci.yml`, `spec/architecture.md` §4) | permanent | Tests müssen Fakes und Test-Libraries (`testify`, …) frei importieren können; Schicht-Regeln gelten production-only (`LH-FA-ARCH-003`). |
| GNU `make` als Host-Voraussetzung neben Docker (`LH-FA-BUILD-007`, `LH-NFA-PORT-002`) | permanent | Carveout zu `LH-NFA-PORT-002` (möglichst wenige Host-Deps); `make` ist überall verfügbar und der pragmatischste Wrapper für Docker-only-Workflows. |
| `contextcheck`-Ausnahme für `internal/adapter/driving/cli/` (`.golangci.yml`) | permanent | Cobras `RunE`-Signatur (`func(cmd, args) error`) kennt keinen Context-Parameter; der Closure muss `cmd.Context()` extrahieren und durchreichen. contextcheck sieht die Closure-Grenze nicht. Strikte Propagation passiert eine Ebene tiefer in `runInit` (Context als erster Parameter). |

## Disziplin

`LH-FA-PROJDOCS-005` verlangt: jeder neue temporäre Carveout bekommt
**parallel** zu seiner Entstehung einen Slice-Plan in `open/` und einen
Eintrag in der oberen Tabelle. Permanente Carveouts kommen ohne Plan,
aber mit Begründung in die zweite Tabelle.

Neben [`roadmap.md`](roadmap.md) ist diese Datei die zweite zulässige
Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in
`docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und
[`../README.md`](../README.md)).
