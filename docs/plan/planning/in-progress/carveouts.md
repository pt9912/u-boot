# u-boot Carveouts

Master-Inventar aller bewussten Carveouts in der u-boot-Codebase
(`LH-FA-PROJDOCS-005`). Wird laufend gepflegt und liegt deshalb dauerhaft
in `in-progress/`.

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
| `COVERAGE_THRESHOLD=0` Bootstrap (`Dockerfile` coverage-Stage, `scripts/coverage-gate.sh`, `LH-FA-BUILD-008`) | temporär | [`open/slice-m3-coverage-threshold-aktivieren.md`](../open/slice-m3-coverage-threshold-aktivieren.md) |
| `gomodguard_v2.blocked: {}` leer (`.golangci.yml`) | temporär | [`open/slice-v1-gomodguard-rules.md`](../open/slice-v1-gomodguard-rules.md) |
| `depguard`-Regeln aktiv, matchen nichts (alle 8 Schicht-Blöcke in `.golangci.yml`; `spec/architecture.md` §4) | temporär | scharf-Schalten ist gekoppelt an erste produktive Pakete pro Schicht → [`open/slice-m3-depguard-aktivierung-verifizieren.md`](../open/slice-m3-depguard-aktivierung-verifizieren.md) |
| `forbidigo.msg` referenziert nicht-existenten "configured logging port" (`.golangci.yml`) | temporär | [`open/slice-v1-logging-port.md`](../open/slice-v1-logging-port.md) |
| ADR-0004 Folgepunkt: Image-Publish nach GHCR und Trivy-Image-Scan fehlen (`.github/workflows/ci.yml` enthält sie nicht) | temporär | [`open/slice-v1-release-pipeline.md`](../open/slice-v1-release-pipeline.md) |
| ADR-0004 Folgepunkt: Branch-Protection im GitHub-UI ist nicht im Repo versioniert | temporär | [`open/slice-v1-branch-protection-checkliste.md`](../open/slice-v1-branch-protection-checkliste.md) |
| Build-Tag-Pfad `//go:build docker` für Adapter-Integrationstests (`spec/architecture.md` §5) ist nur beschrieben, kein CI-Pfad und kein Adapter-Test existiert | temporär | [`open/slice-v1-docker-integrationstests.md`](../open/slice-v1-docker-integrationstests.md) |
| Doku-/Link-Drift in `docs/`, `spec/`, READMEs ist heute nicht maschinell geprüft (M2-Review #11) | temporär | [`open/slice-v1-markdown-link-validator.md`](../open/slice-v1-markdown-link-validator.md) |
| Slice-Pläne für M1, M2, M2b, M2c liegen nicht in `done/` (M2-Review #10); Roadmap referenziert nur Commit-Hashes | temporär | [`open/slice-v1-retroaktive-slice-plaene.md`](../open/slice-v1-retroaktive-slice-plaene.md) |
| ADR-0001 Folgepunkt: CLI-Framework (`flag` vs. Cobra) ist offen — heute reicht `flag` für `--help`/`--version`, mit Subkommandos wird Cobra fällig | temporär | offener Folgepunkt in ADR-0001 (eigener ADR folgt, wenn der Subkommando-Slice startet, vermutlich M3) |
| ADR-0003 Folgepunkt: `revive`-Custom-Rules sind nicht konfiguriert (Default-Profil) | temporär | offener Folgepunkt in ADR-0003 (eigener ADR, wenn Default-Profil aufweicht) |

## Permanente Carveouts (kein Plan, im Inventar dokumentiert)

| Carveout | Status | Plan / Begründung |
| -------- | ------ | ----------------- |
| `errcheck.exclude-functions` für `fmt.Fprintln`/`Fprintf`/`Fprint` (`.golangci.yml`) | permanent | CLI-Writes auf stdout/stderr können nicht meaningful fehlschlagen; `_, _ =`-Prefix bringt keinen Wert. |
| `testpackage`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | Wiring-Schicht braucht `package main`; externe `_test`-Packages sind dort nicht erzwingbar. |
| `gochecknoglobals`-Ausnahme für `cmd/uboot/` (`.golangci.yml`) | permanent | `var version` wird per `-ldflags="-X main.version=…"` überschrieben — das ist der kanonische Go-Pattern für Build-Metadaten. |
| Test-Carveouts in `_test.go` (`cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam`, `revive(unused-parameter)`) | permanent | Tabellengetriebene Tests und Fakes erzeugen legitim hohe Komplexität / fehlenden Context; Profil-Schwellen passen für Production-Code, nicht für Tests. |
| `!**/*_test.go` als erste files-Pattern in jedem `depguard`-Regelblock (`.golangci.yml`, `spec/architecture.md` §4) | permanent | Tests müssen Fakes und Test-Libraries (`testify`, …) frei importieren können; Schicht-Regeln gelten production-only (`LH-FA-ARCH-003`). |
| GNU `make` als Host-Voraussetzung neben Docker (`LH-FA-BUILD-007`, `LH-NFA-PORT-002`) | permanent | Carveout zu `LH-NFA-PORT-002` (möglichst wenige Host-Deps); `make` ist überall verfügbar und der pragmatischste Wrapper für Docker-only-Workflows. |

## Disziplin

`LH-FA-PROJDOCS-005` verlangt: jeder neue temporäre Carveout bekommt
**parallel** zu seiner Entstehung einen Slice-Plan in `open/` und einen
Eintrag in dieser Tabelle. Permanente Carveouts kommen ohne Plan, aber
mit Begründung in das Inventar unten.

Diese Datei ist die einzige zulässige Ausnahme von der
`slice-`/`tranche-`-Konvention für Dateinamen in
`docs/plan/planning/` neben `roadmap.md` (siehe `LH-FA-PROJDOCS-003`
und [`../README.md`](../README.md)).
