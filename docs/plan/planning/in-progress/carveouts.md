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
- **Plan / Begründung** — bei `temporär`: entweder Verweis auf einen
  Slice-Plan in `docs/plan/planning/{open,next,in-progress}/` **oder**
  Verweis auf ein ADR mit benannten Trigger-Slices in §Folgepunkte
  (gleichwertig — das ADR ist hier das Plan-Artefakt, die
  Trigger-Slice-Namen sind die zugesagten künftigen `open/`-Pläne;
  vgl. ADR-0007 Distributions-Restwege). Bei `permanent`: kurze
  Begründung.

## Temporäre Carveouts (Plan-Pflicht)

| Carveout | Status | Plan / Begründung |
| -------- | ------ | ----------------- |
| `LH-OPEN-002` Paketierung-Restwege offen (`spec/lastenheft.md` §14): Homebrew / Debian/RPM — npm/pip durch ADR-0007 verworfen, GHCR durch ADR-0007 + [`done/slice-v1-release-pipeline.md`](../done/slice-v1-release-pipeline.md) entschieden und ausgeliefert (v0.1.0 am 2026-05-31), Binary durch ADR-0007 + [`done/slice-v2-binary-distribution.md`](../done/slice-v2-binary-distribution.md) entschieden und ausgeliefert (sechs Plattformen Linux/macOS/Windows × amd64/arm64 als GitHub-Release-Asset ab v0.1.1) | temporär | Restwege Homebrew + Debian/RPM haben jeweils einen Plan-Stub in `open/` ([`slice-v2-homebrew-formula.md`](../open/slice-v2-homebrew-formula.md), [`slice-v2-distro-pakete.md`](../open/slice-v2-distro-pakete.md)) mit `on hold pending trigger`-Status; Tranchen werden bei Trigger-Feuer (macOS-Anfrage bzw. Distro-Anfrage) ausgearbeitet. Verbindlich dokumentiert in [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Entscheidung (Tabelle pro Option) und §Folgepunkte. |
| Keycloak-Acceptance-Test (`internal/e2e/keycloak_acceptance_docker_test.go`, slice-v1-keycloak T3) unter zusätzlichem build-tag `acceptance_extended` — default `make test-docker` läuft ihn nicht. `docker compose up` failt in GitHub-Actions reproduzierbar nach <1 s mit „compose runtime error"; lokal Quay.io ebenfalls 502/504 zur selben Zeit. | temporär | Folge-Slice `slice-v1-keycloak-ci-flake` (Trigger: Compose-Verbose-Logs aus CI ziehen, dann entweder Pull-Retry-Wrapper im UpService oder Quay-Mirror via Docker-Hub-Pull-Through-Cache). Plan-Anker im Slice-Plan [`done/slice-v1-keycloak.md`](../done/slice-v1-keycloak.md) §Out of Scope. |
| Doctor-Drift-Check für Devcontainer-Features (`devcontainer.features.drift`, [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md) §AK „Doctor-Integration Teil B") nicht im Parent-T5 implementiert. Carveout-Trigger gefeuert nach T4+Followup: T1-T4-Real-LOC ≈ 1009 > 800-Schwelle gemäß Parent-Plan-Vertrag. | temporär | Folge-Slice [`slice-followup-devcontainer-features-drift-doctor`](../open/slice-followup-devcontainer-features-drift-doctor.md) (Status: ready — Trigger gefeuert, Start nach Parent-Slice-T5-Abschluss). |

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

## Carveout-Auflösungs-Slices (historisch)

Slices, die ausschließlich offene Carveouts (`LH-FA-PROJDOCS-005`)
aufgelöst haben. Quelle der Wahrheit für das Carveout-Inventar
sind die zwei Tabellen oben; diese Tabelle ist der Audit-Trail der
Slices, die das Inventar geschrumpft haben — verbindlich verankert
in [`roadmap.md`](roadmap.md) ist nur noch der Pointer auf diese
Sektion, nicht die Tabelle selbst.

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
| [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md) | ADR-0004 Folgepunkte Image-Publish + Trivy; `LH-OPEN-002` Paketierung (GHCR-Anteil) | V1 | Done (T1 `0f64938`, T2 `93b703e`, T3 `8212889`, T4 `066917a`, T5 `bc487fc` — Branch-Protection-Teilabschluss 2026-05-27) |
| [`slice-v1-markdown-link-validator`](../done/slice-v1-markdown-link-validator.md) | Doku-/Link-Drift in `docs/`/`spec/` nicht maschinell geprüft | V1-vorgezogen | Done |
| [`slice-v1-backup-streaming-copy`](../done/slice-v1-backup-streaming-copy.md) | `LH-FA-INIT-005` Backup heute mit `ReadFile`+`WriteFile`; harter 256-MiB-Cap als MVP-Workaround | V1-vorgezogen | Done |
| [`slice-v1-plugin-system-entscheidung`](../done/slice-v1-plugin-system-entscheidung.md) | `LH-OPEN-003` Plugin-System offen | V1 | Done (Entscheidung in [ADR-0008](../../adr/0008-plugin-system-statisch.md): statisch) |
| [`slice-v1-template-format-entscheidung`](../done/slice-v1-template-format-entscheidung.md) | `LH-OPEN-004` Template-Format offen | V1 | Done (Entscheidung in [ADR-0009](../../adr/0009-template-format-yaml-files.md): YAML+`text/template`) |
| [`slice-v1-yaml-parse-error-sentinel`](../done/slice-v1-yaml-parse-error-sentinel.md) | M7-T5-Review-Followup N2: `YAMLCodec`-Port unterscheidet Parse- nicht von IO-Fehlern; Exit-Code-14-vs-10-Klassifikation reißt bei kaputter `compose.yaml` unter `u-boot generate devcontainer` | V1-vorgezogen | Done (`1008326`) |
| [`slice-v2-revive-custom-rules`](../done/slice-v2-revive-custom-rules.md) | ADR-0003 Folgepunkt revive-Custom-Rules | V2-vorgezogen | Done |
| [`slice-later-http-driving-adapter`](../done/slice-later-http-driving-adapter.md) | `spec/architecture.md` §7 HTTP-Driving-Adapter prospektiv | Later | Done (Entscheidung in [ADR-0010](../../adr/0010-kein-http-driving-adapter.md): wird nicht gebaut) |
| [`slice-v0.1.1-doctor-container-awareness`](../done/slice-v0.1.1-doctor-container-awareness.md) | `doctor` im distroless-Container findet docker/git nicht (Real-world-Befund 2026-05-31 post-v0.1.0) | v0.1.1-Followup | Done (T1 `9a99bbf`, T2 `c35360f`, T3 `111e725`, T4 schließt; Tag-Push bleibt Nutzer-Aktion analog v0.1.0-T4) |
| [`slice-v2-binary-distribution`](../done/slice-v2-binary-distribution.md) | ADR-0007 §Folgepunkte 1 Trigger (erste konkrete Cross-Plattform-Distributionsanfrage) durch `doctor`-Befund ausgelöst | V2 | Done — T1 ✅ `dc9a336` + `f3f1731` (`make build-binaries` für 6 Plattformen Linux/macOS/Windows × amd64/arm64), T2 ✅ `5e5166b` (`publish.yml` build + GitHub-Release-Upload), T3 ✅ `866f6fd` (READMEs Install-Block Binary-first + CHANGELOG `## [Unreleased]`), T4 `2f39511` Slice-Closure mit ADR-0007 §Entscheidung-Update (Binary „Vertagt → Gewählt") + carveouts.md `LH-OPEN-002`-Reduktion auf Homebrew+Debian/RPM + open→done. |

## Disziplin

`LH-FA-PROJDOCS-005` verlangt: jeder neue temporäre Carveout bekommt
**parallel** zu seiner Entstehung einen Plan-Anker und einen Eintrag
in der oberen Tabelle. Plan-Anker ist entweder ein Slice-Plan in
`open/` (Standardform) **oder** ein ADR mit benannten
Trigger-Slices in §Folgepunkte (für ADR-getriebene Vertagungen
mit konkretem Trigger statt offener Implementierung). Permanente
Carveouts kommen ohne Plan, aber mit Begründung in die zweite
Tabelle.

Neben [`roadmap.md`](roadmap.md) ist diese Datei die zweite zulässige
Ausnahme von der `slice-`/`tranche-`-Konvention für Dateinamen in
`docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und
[`../README.md`](../README.md)).
