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
| v0.4.0 | in progress | 2026-06-02 | `logs`, JSON-/Dry-Run-CLI, restliche V1/Later-Trigger | diese Datei |

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
| [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md) | `in-progress/`, T0 ✅ (Stub `1d3f652` + Review-Findings R1 `5d651bf` + T0-Outcomes + R1/R2 `c6c3bb2`); Folge-Slice 1/9 Doctor **done** ([`done/slice-v1-cli-json-dry-run-doctor.md`](../done/slice-v1-cli-json-dry-run-doctor.md), T0-T5 DoD-Hashes im Slice-File); Folge-Slice 2/9 Add **done** ([`done/slice-v1-cli-json-dry-run-add.md`](../done/slice-v1-cli-json-dry-run-add.md), T0-T6 + Review-Round-6/7 DoD-Hash-Tabelle im Slice-File); Folge-Slice 3/9 Init **done** ([`done/slice-v1-cli-json-dry-run-init.md`](../done/slice-v1-cli-json-dry-run-init.md), T0-T8 + Review-Round-9 DoD-Hash-Tabelle im Slice-File); Folge-Slice 4/9 Generate **done** ([`done/slice-v1-cli-json-dry-run-generate.md`](../done/slice-v1-cli-json-dry-run-generate.md), T0-T8 + Review-Rounds 1-7 + 10 DoD-Hash-Tabelle im Slice-File); Folge-Slice 5/9 Remove **done** ([`done/slice-v1-cli-json-dry-run-remove.md`](../done/slice-v1-cli-json-dry-run-remove.md), T0-T8 + R1-R12 + R13-R15 DoD-Hash-Tabelle im Slice-File; Cluster-Stand 5/9 done, 4/9 offen — `up-down`, `logs`, `config`, `template`) | Cluster-Slice für maschinenlesbare CLI (`LH-FA-CLI-007/008`, `LH-NFA-USE-004`); Per-Command-Inkrementell-Strategie mit 9 Folge-Slices, V1-pünktlich wegen ADR-0010 Re-Eval-Trigger 2. **Closure-Hard-Rule (strict):** Cluster-Closure ausschließlich bei allen 9 Folge-Slices in `done/` — kein Quorum, kein Carveout als Closure-Alternative; Slip-Pfad ist Notfall-Restlauf, kein wählbarer Pfad. T0-Outcomes festgezurrt (Flag-Scope, RecordingFileSystem mit Passthrough, Common-Envelope, Diff-Renderer, doctor-vor-add). Add-Slice etabliert die Cluster-Infrastruktur als Pattern-Vorbild (RecordingFileSystem mit impliziter MkdirAll-Modellierung, Pure-Go Diff-Renderer, fsFactory-Closure pro AddPreviewMode, CountAdditions-Semantik per Spec §477) für die vier folgenden modifying-Slices. Init-Slice erbt das Pattern 1:1 und ergänzt `initGit`-Skip im Dry-Run, ProgressPort-Silencing im JSON-Mode, Template-Mutex zu Preview-Flags und `ErrInitFileSystem`-Multi-`%w`-Wrap mit Switch-Order-Pflicht. Generate-Slice erbt das Pattern von init 1:1 und ergänzt Multi-Artefakt-Envelope-Form (`command="generate"`, `data.artifact`/`data.action`), `cliJSONEnvelope.Data`-Feld + `newDataEnvelope`-Konstruktor (aus Template-Slice 9/9 vorgezogen — T9-T1 entfällt damit), per-Artefakt LH-Code-Tabelle in `mapGenerateErrorToDiagnostic(err, artifact)` und `ErrConfigValueInvalid`-Sentinel-Wrap für den LH-FA-DEV-003-URL-Reject-Pfad. Remove-Slice erbt von generate 1:1 (data-Envelope, Pattern-Erbe-Disziplin) und ergänzt Confirmer-Swap mit `noopConfirmer` für `--purge`-Gate-Silencing im JSON-Mode, WARN-Migration via `LH-FA-ADD-007`-Multi-Use (ERROR `ErrServiceUnregistered` + WARN deferred-Volumes, Disambiguation über `(code, level)`), Custom-`Args`-Validator `validateRemoveArgs(a *App)` mit `*App`-Closure für stdout-Envelope-Emission VOR Cobra-Return (Spec §1841/§1842-Konformität bei NoPositionalArg + TooManyArgs), `baseDirSanitizedError`-Wrapper für `diagnostic.message` (Path-Leak-Defense + Substring-Kollisions-Robustheit via `replaceBareBaseDir`-Word-Boundary), Dry-Run-WARN-Suppression in `printRemoveSummary` und `delete`-Action für `RemoveAll`-Captures. Nächster Schritt: Folge-Slice 6/9 (`up-down`) nach Cluster-T0-(e)-Reihenfolge. |
| [`slice-v1-keycloak-ci-flake`](../open/slice-v1-keycloak-ci-flake.md) | `open/`, on hold | Keycloak-Acceptance-Flake analysieren, sobald CI-Logs/Quay- oder Mirror-Befund belastbar sind. |
| [`slice-v2-homebrew-formula`](../open/slice-v2-homebrew-formula.md) | `open/`, on hold | Erste konkrete macOS-/Homebrew-Nutzeranfrage. |
| [`slice-v2-generate-devcontainer-rollback-aware-write`](../open/slice-v2-generate-devcontainer-rollback-aware-write.md) | `open/`, on hold pending trigger | Carveout-Plan-Anker für `generate devcontainer` Phase-2-Half-Write-State (siehe [`carveouts.md`](carveouts.md) §Temporäre Carveouts). Trigger: Real-World-Beschwerde über Half-State oder Devcontainer-Schema-Erweiterung. Bevorzugte Skizze Option 1 (Snapshot + Rollback). |
| [`slice-v2-distro-pakete`](../open/slice-v2-distro-pakete.md) | `open/`, on hold | Konkrete Debian-/RPM-Anfrage mit Bereitschaft für Packaging-Overhead. |
| `slice-later-local-templates` | noch kein Slice-Plan | `--template ./pfad` konkretisieren (`LH-FA-TPL-003`). |
| `slice-later-migration` | noch kein Slice-Plan | Konfigurationsmigration konkretisieren (`LH-FA-CONF-006`). |
| `slice-later-custom-data-sources` | noch kein Slice-Plan | Erweiterung jenseits YAML-Quellen konkretisieren (`LH-DA-004`). |
| `slice-vN-podman-formal` | noch kein Slice-Plan | Podman-first Probe-Adapter und CI-Matrix konkretisieren; heutiger Stand bleibt Docker-compatible Drop-in. |
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
