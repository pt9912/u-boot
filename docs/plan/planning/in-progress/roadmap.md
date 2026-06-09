# u-boot Roadmap

Übergreifendes Master-Dokument zum aktuellen Stand der Slices und
Tranchen (`LH-FA-PROJDOCS-003`). Diese Datei liegt dauerhaft in
`in-progress/`, bleibt aber bewusst knapp: Sie steuert aktuelle Arbeit,
Trigger und nächste Entscheidungen. Historische Release-Details liegen
in [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md).

## Aktueller Snapshot

| Version | Status | Datum | Fokus | Detailquelle |
| --- | --- | --- | --- | --- |
| v0.1.0 | released | 2026-05-31 | MVP-Core, Release-Pipeline, GHCR | [`slice-v1-release-cut-v0.1.0`](../done/slice-v1-release-cut-v0.1.0.md) |
| v0.2.0 | released | 2026-06-01 | Container-aware `doctor`, Binary-Distribution, Template-Katalog | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | released | 2026-06-01 | Add-on Catalogue Expansion (`remove`, deps, Keycloak, OTel) | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| v0.4.0 | released | 2026-06-08 | Maschinenlesbare CLI (`--json`/`--dry-run`/`--diff` für alle 10 Subcommands), `logs`, Devcontainer-Features | [`slice-v1-release-cut-v0.4.0`](../done/slice-v1-release-cut-v0.4.0.md) (Tag `v0.4.0` auf `bce886f`) |

## v0.4.0 Arbeitspakete

Diese Punkte sind Arbeitspakete für Version 0.4.0. Einige sind bereits
als Slice-Plan angelegt, andere bleiben bis zur Ausarbeitung als
benannte APs in dieser Roadmap.

Done v0.4.0-Slices (z. B. [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md),
[`slice-v1-logs`](../done/slice-v1-logs.md)) wandern beim v0.4.0-
Release-Cut zusammenfassend in §Bereits Geschlossen — analog zu
v0.2.0/v0.3.0. Bis dahin bleibt der per-Tranche-Audit-Trail
im jeweiligen Slice-File die Quelle der Wahrheit; diese Tabelle
listet nur das **offene v0.4.0-Backlog**.

| AP | Status | Entscheidung / nächster Schritt |
| --- | --- | --- |
| [`slice-v1-cli-json-dry-run`](../done/slice-v1-cli-json-dry-run.md) | **`done/` — Cluster vollständig abgeschlossen (9/9 Folge-Slices + T_close)**. T0 ✅ (Stub `1d3f652` + Review-Findings R1 `5d651bf` + T0-Outcomes + R1/R2 `c6c3bb2`); Folge-Slice 1/9 Doctor **done** ([`done/slice-v1-cli-json-dry-run-doctor.md`](../done/slice-v1-cli-json-dry-run-doctor.md), T0-T5 DoD-Hashes im Slice-File); Folge-Slice 2/9 Add **done** ([`done/slice-v1-cli-json-dry-run-add.md`](../done/slice-v1-cli-json-dry-run-add.md), T0-T6 + Review-Round-6/7 DoD-Hash-Tabelle im Slice-File); Folge-Slice 3/9 Init **done** ([`done/slice-v1-cli-json-dry-run-init.md`](../done/slice-v1-cli-json-dry-run-init.md), T0-T8 + Review-Round-9 DoD-Hash-Tabelle im Slice-File); Folge-Slice 4/9 Generate **done** ([`done/slice-v1-cli-json-dry-run-generate.md`](../done/slice-v1-cli-json-dry-run-generate.md), T0-T8 + Review-Rounds 1-7 + 10 DoD-Hash-Tabelle im Slice-File); Folge-Slice 5/9 Remove **done** ([`done/slice-v1-cli-json-dry-run-remove.md`](../done/slice-v1-cli-json-dry-run-remove.md), T0-T8 + R1-R12 + R13-R15 DoD-Hash-Tabelle im Slice-File); Folge-Slice 6/9 Up-Down **done** ([`done/slice-v1-cli-json-dry-run-up-down.md`](../done/slice-v1-cli-json-dry-run-up-down.md), T0-T8 + R1-R6 + T7-Adressierung + T8-Bestätigungsrunde DoD-Hash-Tabelle im Slice-File); Folge-Slice 7/9 Logs **done** ([`done/slice-v1-cli-json-dry-run-logs.md`](../done/slice-v1-cli-json-dry-run-logs.md), T0-T8 + R1-R3 + Pre-T6-Review + Pre-T8-Bestätigungsrunde DoD-Hash-Tabelle im Slice-File); Folge-Slice 8/9 Config **done** ([`done/slice-v1-cli-json-dry-run-config.md`](../done/slice-v1-cli-json-dry-run-config.md), T0–T8 + drei Review-Runden R-T4-1/R-IR-1/R-CLI-1 DoD-Hash-Tabelle im Slice-File — zwei neue Sentinels `ErrConfigWriteRejected`/`ErrConfigPostPatchSanityFailed` (T0-(m)-Split) + Multi-`%w` + SilenceLogger + Orphan-WARN-Migration + PreviewMode-Dry-Run-FS-Routing + vollständige CLI-Neufassung (3 Data-Carrier, Subcommand-Pflicht inkl. Error-Pfad via subcommand-bewusste `reportErrorSub`, Allowlist 3 Forms/Reject 4→1, Mapper Switch-Order T0-(f), `configArgsValidator`, Voll-Schema/Dry-Run/Diff, bare/get-Reject, WARN→diagnostics); zwei Review-Runden vor T5 (R-T4-1 + unabhängig R-IR-1, beide HIGH gefixt); R1+R2+R3 durchlaufen; 16 T0-Sub-Decisions festgezurrt; drei Review-Runden (R-T4-1 + R-IR-1 + R-CLI-1, alle gefixt); vier Folge-Carveout-Stubs in `open/`; LOC ~1500-1900; T8-Closure: CHANGELOG + `cli-json-output.md` §6.9 + carveouts + `done/`-Move); Cluster-Stand **vollständig done: 9/9 Folge-Slices + T_close (`3a35d58`)** — Übergangs-Mechanik abgebaut, bare-`template`-RunE-Reject, `LH-NFA-USE-004`-Surface vollständig ausgeliefert (alle zehn Spec-Enum-Forms tragen `--json`); keine Folge-ADR (SD-1) | Cluster-Slice für maschinenlesbare CLI (`LH-FA-CLI-007/008`, `LH-NFA-USE-004`); Per-Command-Inkrementell-Strategie mit 9 Folge-Slices, V1-pünktlich wegen ADR-0010 Re-Eval-Trigger 2. **Closure-Hard-Rule (strict):** Cluster-Closure ausschließlich bei allen 9 Folge-Slices in `done/` — kein Quorum, kein Carveout als Closure-Alternative; Slip-Pfad ist Notfall-Restlauf, kein wählbarer Pfad. T0-Outcomes festgezurrt (Flag-Scope, RecordingFileSystem mit Passthrough, Common-Envelope, Diff-Renderer, doctor-vor-add). Add-Slice etabliert die Cluster-Infrastruktur als Pattern-Vorbild (RecordingFileSystem mit impliziter MkdirAll-Modellierung, Pure-Go Diff-Renderer, fsFactory-Closure pro AddPreviewMode, CountAdditions-Semantik per Spec §477) für die vier folgenden modifying-Slices. Init-Slice erbt das Pattern 1:1 und ergänzt `initGit`-Skip im Dry-Run, ProgressPort-Silencing im JSON-Mode, Template-Mutex zu Preview-Flags und `ErrInitFileSystem`-Multi-`%w`-Wrap mit Switch-Order-Pflicht. Generate-Slice erbt das Pattern von init 1:1 und ergänzt Multi-Artefakt-Envelope-Form (`command="generate"`, `data.artifact`/`data.action`), `cliJSONEnvelope.Data`-Feld + `newDataEnvelope`-Konstruktor (aus Template-Slice 9/9 vorgezogen — T9-T1 entfällt damit), per-Artefakt LH-Code-Tabelle in `mapGenerateErrorToDiagnostic(err, artifact)` und `ErrConfigValueInvalid`-Sentinel-Wrap für den LH-FA-DEV-003-URL-Reject-Pfad. Remove-Slice erbt von generate 1:1 (data-Envelope, Pattern-Erbe-Disziplin) und ergänzt Confirmer-Swap mit `noopConfirmer` für `--purge`-Gate-Silencing im JSON-Mode, WARN-Migration via `LH-FA-ADD-007`-Multi-Use (ERROR `ErrServiceUnregistered` + WARN deferred-Volumes, Disambiguation über `(code, level)`), Custom-`Args`-Validator `validateRemoveArgs(a *App)` mit `*App`-Closure für stdout-Envelope-Emission VOR Cobra-Return (Spec §1841/§1842-Konformität bei NoPositionalArg + TooManyArgs), `baseDirSanitizedError`-Wrapper für `diagnostic.message` (Path-Leak-Defense + Substring-Kollisions-Robustheit via `replaceBareBaseDir`-Word-Boundary), Dry-Run-WARN-Suppression in `printRemoveSummary` und `delete`-Action für `RemoveAll`-Captures. Up-Down-Slice (6/9) ist Read-only-Klasse: nur `--json` (kein `--dry-run`/`--diff`, Cluster Z. 464-467); ergänzt `SilenceProgress`-Bool-Field-Pattern (Pattern-Erbe init's Interface-Swap aber mit `io.Writer`-Variable, nicht 1:1), Request-time Confirmer-Gate-Branch für `down --volumes` ohne `downMu`-Mutex (race-frei by construction), zwei neue Read-spezifische FS-Sentinels (`ErrUpFileSystem`/`ErrDownFileSystem` mit `"filesystem read failed"`-Message), shared `mapComposeRuntimeSentinel`-Helper in neuem `cli/composesentinel.go`, `(code, exitCode)`-Tupel-Disambiguation-Pattern für Multi-`%w`-Wraps (Mapper FS-first ↔ ExitCode-Helper Driven-first als by-design Klassen-vs-Sub-Klassen-Trennung) und Sanitizer-Helper-Extraktion `cli/remove.go:465-538` → `cli/sanitize.go` für cluster-weite Wiederverwendung. Logs-Slice (7/9) erbt das Read-only-Klasse-Pattern 1:1 und ergänzt T0-(a) **Single-Envelope + `--follow --json` Reject** als Cluster-Konsens (Spec-§1841 honoriert; NDJSON-Stream-Form Cluster-weit nicht vorgesehen) mit neuem `ErrFollowJSONNotSupported` (LH-FA-CLI-006/Exit 2), `ErrLogsFileSystem` als dritter Read-FS-Sentinel des Read-only-Trios, `logsStatusData.lines []string` ohne `omitempty` (Empty-Array-Pin: `[]`, nicht `null`), T0-(i) Validation-Reihenfolge-Pin (`--follow --json` schlägt `--tail=-1`), Trailing-Newline-Strip in `splitLogLines` und FS+Docker-Multi-`%w`-Switch-Order-Defense-Pin mit `_ByDesign`-Suffix (Pattern-Erbe up-down §6.7). CRLF-Edge-Case als bekannte Limitation in §6.8 dokumentiert (LOW-2). Drei Folge-Stubs in `open/` (`logs-format-flags`, `logs-multi-service-filter`, `logs-time-range-filter`). Config-Slice (8/9) ist erste Read-only+Modifying-Hybrid-Klasse: bündelt `config`/`config get` (Read-only) + `config set` (Modifying mit `--dry-run`/`--diff`); T0–T8 done (Drei-Klassen-Sentinel-Split T0-(m), Multi-`%w`, SilenceLogger, Orphan-WARN-Migration, Dry-Run-FS-Routing, vollständige `--json`-CLI mit Subcommand-Pflicht + Voll-Schema/Dry-Run/Diff, Acceptance-Suite, drei Review-Runden, T8-Closure inkl. §6.9-Doku). Reject-Liste 4→1 (nur `template (bare)` offen). Folge-Slice 9/9 template **done** ([`done/slice-v1-cli-json-dry-run-template.md`](../done/slice-v1-cli-json-dry-run-template.md), T0→T2→T4 DoD-Hash-Tabelle; `template list --json` Array→Minimalkontrakt-Envelope migriert, Breaking-Change CHANGELOG `### Changed`; T3 bare-Reject nach T_close verschoben; kleinster Slice ~60 LOC). **T_close abgeschlossen (`3a35d58`)**: Übergangs-Mechanik (`jsonallowlist.go`, Gate, `ErrJSONNotImplemented`) entfernt, bare-`template`-RunE-Reject (`ErrTemplateSubcommandRequired`), Public-Doku nachgezogen, netto −90 LOC. **Damit ist die gesamte 9-teilige JSON-CLI-Serie + ADR-0010-Trigger-2 ausgeliefert.** |
| [`slice-v1-keycloak-ci-flake`](../open/slice-v1-keycloak-ci-flake.md) | `open/`, on hold | Keycloak-Acceptance-Flake analysieren, sobald CI-Logs/Quay- oder Mirror-Befund belastbar sind. |
| [`slice-v2-homebrew-formula`](../open/slice-v2-homebrew-formula.md) | `open/`, on hold | Erste konkrete macOS-/Homebrew-Nutzeranfrage. |
| [`slice-v2-generate-devcontainer-rollback-aware-write`](../open/slice-v2-generate-devcontainer-rollback-aware-write.md) | `open/`, on hold pending trigger | Carveout-Plan-Anker für `generate devcontainer` Phase-2-Half-Write-State (siehe [`carveouts.md`](carveouts.md) §Temporäre Carveouts). Trigger: Real-World-Beschwerde über Half-State oder Devcontainer-Schema-Erweiterung. Bevorzugte Skizze Option 1 (Snapshot + Rollback). |
| [`slice-v2-distro-pakete`](../open/slice-v2-distro-pakete.md) | `open/`, on hold | Konkrete Debian-/RPM-Anfrage mit Bereitschaft für Packaging-Overhead. |
| [`slice-later-local-templates`](../done/slice-later-local-templates.md) | **`done/`** (Later) | ✅ **`LH-FA-TPL-003` ausgeliefert** — `u-boot init --template ./pfad` (lokale FS-Templates). T1–T5 + Pre-T5-Review (`66c347d`/`87a8704`/`5031b5f`/`adaafbe`/`7d63532`/`4fdf1a6`); Carveout in [`carveouts.md`](carveouts.md) §Auflösungs-Slices aufgelöst. |
| `slice-later-migration` | noch kein Slice-Plan | Konfigurationsmigration konkretisieren (`LH-FA-CONF-006`). |
| `slice-later-custom-data-sources` | noch kein Slice-Plan | Erweiterung jenseits YAML-Quellen konkretisieren (`LH-DA-004`). |
| `slice-vN-podman-formal` | noch kein Slice-Plan | Podman-first Probe-Adapter und CI-Matrix konkretisieren; heutiger Stand bleibt Docker-compatible Drop-in. |
| `slice-vN-harness-bootstrap-scaffold` | noch kein Slice-Plan; Entscheidung offen in [ADR-0011](../../adr/0011-agent-harness-scaffolding.md) (`Proposed`) | u-boot scaffoldet **Agent-Harness-Artefakte** (`AGENTS.md`/`harness/`/`spec/`/`docs/plan/`/ADR) opt-in und GF/BF-aware — nutzt die bestehende `LH-FA-INIT-004`-Detection (GF=frisch, BF=bestehend); u-boot legt nur das Skelett, der Agent reconciled. Idee aus dem `ai-harness-course` (Modul 2, GF/BF-Bootstrap). Vor Code: ADR ratifizieren + Spec-Erweiterung (`LH-FA-*` / `LH-ZB-002`) + Lizenz-Check. |
| `slice-vN-devcontainer-egress-firewall` | noch kein Slice-Plan; Entscheidung offen in [ADR-0012](../../adr/0012-devcontainer-egress-firewall.md) (`Proposed`) | **Network-hardened Devcontainer** (opt-in): `init-firewall.sh` (iptables+ipset, default-DROP + Allowlist) + `runArgs --cap-add=NET_ADMIN` + `postCreate`, Allowlist als Config (`devcontainer.firewall.allow`) mit Ökosystem-Defaults, plus doctor-Check + graceful degradation ohne `NET_ADMIN`. Runtime-Pendant zur Build-Time-Allowlist `LH-FA-DEV-003`; Guardrail, kein Sandbox. Vor Code: ADR ratifizieren + Spec-Erweiterung (`LH-FA-DEV-*` / ggf. `LH-NFA-SEC-*`). |
| Branch-Protection-UI | Nutzeraktion, kein Code-Slice | Repo-Owner aktiviert Required Checks vor erstem externem PR; Anleitung in [`docs/user/branch-protection.md`](../../../user/branch-protection.md). |

## Bereits Geschlossen

Die abgeschlossenen Slices bleiben in [`done/`](../done/) die
Detailquelle. Für die Agentensteuerung reicht hier der Cluster-Überblick:

| Cluster | Ergebnis | Detailquelle |
| --- | --- | --- |
| MVP M1..M8 | Repo-Skeleton, Architektur, CI/Gates, `init`, `doctor`, `add postgres`, `up/down`, `generate`, `config` | [`done/`](../done/) und [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md) |
| v0.2.0 | Container-aware `doctor`, sechs Plattform-Binaries, Template-Katalog | [`slice-v1-release-cut-v0.2.0`](../done/slice-v1-release-cut-v0.2.0.md) |
| v0.3.0 | Add-on Catalogue Expansion: `remove`, `--with-deps`, Keycloak, OTel, V1-Audit | [`slice-v1-release-cut-v0.3.0`](../done/slice-v1-release-cut-v0.3.0.md) |
| Harness-Doku | Agent-Briefing, Harness-Einstieg, Rollentrennung | [`AGENTS.md`](../../../../AGENTS.md), [`harness/README.md`](../../../../harness/README.md), [`harness/roles.md`](../../../../harness/roles.md) |

## Verwandte Dokumente

- [`carveouts.md`](carveouts.md) — Master-Inventar aller temporären und
  permanenten Carveouts (`LH-FA-PROJDOCS-005`), plus Audit-Trail der
  Slices, die offene Carveouts geschlossen haben.
- [`README.md`](../README.md) — Slice-/Tranche-Konventionen für
  Dateinamen in `docs/plan/planning/` (`LH-FA-PROJDOCS-003`).
- [`docs/archive/roadmap-history-v0.1-v0.3.md`](../../../archive/roadmap-history-v0.1-v0.3.md)
  — ausgelagerte Release-Historie.

## Pflege-Regeln

- Diese Roadmap beschreibt aktuelle Steuerung, nicht jede historische
  Tranche.
- Release-Details, lange Commit-Listen und retrospektive Tabellen
  gehören in `done/`-Slices oder nach `docs/archive/`.
- Neue v0.4.0-Arbeitspakete brauchen eine klare Entscheidung oder einen
  nächsten Schritt. Ohne ausgearbeiteten Plan bleiben sie als benanntes
  AP hier, nicht als halbfertiger Slice.
- Diese Datei ist die einzige zulässige Ausnahme von der
  `slice-`/`tranche-`-Konvention für Dateinamen in
  `docs/plan/planning/` (siehe `LH-FA-PROJDOCS-003` und
  [`../README.md`](../README.md)).
