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
| `LH-OPEN-002` Paketierung-Restwege offen (`spec/lastenheft.md` §14): Binary-Release / Homebrew / Debian/RPM — npm/pip durch ADR-0007 verworfen, GHCR durch ADR-0007 + [`done/slice-v1-release-pipeline.md`](../done/slice-v1-release-pipeline.md) entschieden und ausgeliefert | temporär | Restwege haben jeweils einen benannten Trigger-Slice-Plan (`slice-v2-binary-distribution.md`, `slice-v2-homebrew-formula.md`, `slice-v2-distro-pakete.md`), der bei Auslösung in `open/` angelegt wird; benannte Restwege + Trigger sind verbindlich dokumentiert in [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung (Tabelle pro Option) und §Folgepunkte. |
| HTTP-Driving-Adapter ist als „geplante Erweiterung" in `spec/architecture.md` §7 erwähnt, aber nicht spezifiziert oder gefordert | temporär | [`open/slice-later-http-driving-adapter.md`](../open/slice-later-http-driving-adapter.md) |

## Permanente Carveouts (kein Plan, im Inventar dokumentiert)

| Carveout | Status | Plan / Begründung |
| -------- | ------ | ----------------- |
| `errcheck.exclude-functions` für `fmt.Fprintln`/`Fprintf`/`Fprint` (`.golangci.yml`) | permanent | CLI-Writes auf stdout/stderr können nicht meaningful fehlschlagen; `_, _ =`-Prefix bringt keinen Wert. |
| `testpackage`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | Wiring-Schicht braucht `package main`; externe `_test`-Packages sind dort nicht erzwingbar. |
| `gochecknoglobals`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | `var version` wird per `-ldflags="-X main.version=…"` überschrieben — das ist der kanonische Go-Pattern für Build-Metadaten. |
| Test-Carveouts in `_test.go` (`cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam`, `revive(unused-parameter)`) | permanent | Tabellengetriebene Tests und Fakes erzeugen legitim hohe Komplexität / fehlenden Context; Profil-Schwellen passen für Production-Code, nicht für Tests. |
| `!**/*_test.go` als erste files-Pattern in jedem `depguard`-Regelblock (`.golangci.yml`, `spec/architecture.md` §4) | permanent | Tests müssen Fakes und Test-Libraries (`testify`, …) frei importieren können; Schicht-Regeln gelten production-only (`LH-FA-ARCH-003`). |
| GNU `make` als Host-Voraussetzung neben Docker (`LH-FA-BUILD-007`, `LH-NFA-PORT-002`) | permanent | Carveout zu `LH-NFA-PORT-002` (möglichst wenige Host-Deps); `make` ist auf Linux/macOS und unter WSL/Git-Bash auf Windows praktisch durchgängig verfügbar und der pragmatischste Wrapper für Docker-only-Workflows. Windows-PowerShell-natives Setup (ohne WSL/Git-Bash) wäre eine eigene LH-NFA-PORT-002-Erweiterung mit eigenem Slice. |
| `contextcheck`-Ausnahme für `internal/adapter/driving/cli/` (`.golangci.yml`) | permanent | Cobras `RunE`-Signatur (`func(cmd, args) error`) kennt keinen Context-Parameter; der Closure muss `cmd.Context()` extrahieren und durchreichen. contextcheck sieht die Closure-Grenze nicht. Strikte Propagation passiert eine Ebene tiefer in `runInit` (Context als erster Parameter). |
| `interfacebloat`-Ausnahme für `internal/hexagon/port/driven/filesystem.go` (`.golangci.yml`) | permanent | `driven.FileSystem` ist die zentrale FS-Abstraktion mit 12 Methoden (Exists/Lstat/ReadFile/WriteFile/WriteFileExclusive/Mkdir/MkdirAll/Rename/ReadDir/RemoveAll/Copy/CopyExclusive). Eine künstliche Aufspaltung würde Test-Fakes verkomplizieren ohne semantischen Mehrwert; das interfacebloat-Limit (10) bewusst aufgeweicht für genau diese Schnittstelle. |

## Disziplin

`LH-FA-PROJDOCS-005` verlangt: jeder neue temporäre Carveout bekommt
**parallel** zu seiner Entstehung einen Slice-Plan in `open/` und einen
Eintrag in der oberen Tabelle. Permanente Carveouts kommen ohne Plan,
aber mit Begründung in die zweite Tabelle.

Neben [`roadmap.md`](roadmap.md) ist diese Datei die zweite zulässige
Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in
`docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und
[`../README.md`](../README.md)).
