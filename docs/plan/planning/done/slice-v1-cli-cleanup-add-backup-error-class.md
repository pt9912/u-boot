# Slice V1: `add` Backup-Sentinel LH-Code-Harmonisierung

> **Status:** ✅ **done** — Cleanup-Folge-Slice aus
> [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)
> T7 Review-Round-9 Finding #5.
>
> **DoD:** Commit `1cbc3a7` (`fix(cli)`: `mapAddErrorToDiagnostic`
> Backup-Sentinels von [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)
> umgestellt; `erroremission_internal_test.go`-Cases nachgezogen;
> add-Pfad-Defense-Branch dokumentiert; CHANGELOG `### Fixed`-
> Eintrag plus Code-Liste in add-Slice-Eintrag auf
> `LH-FA-INIT-{004,006}` korrigiert).

## Auslöser

Beim Review der [`slice-v1-cli-json-dry-run-init`](slice-v1-cli-json-dry-run-init.md)-Implementierung
(Review-Round-9 Finding #5, Schweregrad: low / cross-slice) ist
eine Divergenz zwischen dem `add`- und dem `init`-Diagnostic-
Mapper aufgefallen:

- [`internal/adapter/driving/cli/add.go`](../../../../internal/adapter/driving/cli/add.go)
  `mapAddErrorToDiagnostic` mappt `ErrBackupSuffixExhausted` /
  `ErrBackupSourceMissing` zu **[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)** (Validation-
  Klasse, Code-Class-Implication: Exit 10).
- [`internal/adapter/driving/cli/init.go`](../../../../internal/adapter/driving/cli/init.go)
  `mapInitErrorToDiagnostic` mappt dieselben Sentinels zu
  **[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)** (FS-Klasse, Code-Class-Implication: Exit 14).

Der tatsächliche Exit-Code wird in beiden Fällen via
[`cli.go`](../../../../internal/adapter/driving/cli/cli.go)
`isFilesystemError` zu **14** klassifiziert (Backup-Sentinels sind
in der FS-Klasse-Liste). Damit:

- **`init`-Envelope**: Code [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) + Exit 14 → konsistent.
- **`add`-Envelope**: Code [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) + Exit 14 → **inkonsistent**
  (Diagnostic-Klasse sagt Validation, Exit-Klasse sagt FS).

Die `init`-Variante ist die spec-konforme — [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz) ist die
Validation-Cluster-ID (Flag-Mutex, Confirmation-Required), während
[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) die technische Persistenz-Cluster-ID ist, aus der
der Exit-Code stammt. Der `add`-Mapper wurde im
[`slice-v1-cli-json-dry-run-add`](../done/slice-v1-cli-json-dry-run-add.md)
T5 mit dem [`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)-Code geseedet, ohne zu prüfen ob der
Backup-Pfad reell unter Validation oder Persistenz fällt.

## Scope

**Pflicht:**

- [`internal/adapter/driving/cli/add.go`](../../../../internal/adapter/driving/cli/add.go)
  `mapAddErrorToDiagnostic`: `ErrBackupSuffixExhausted` /
  `ErrBackupSourceMissing` Case-Branch nach [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) umstellen
  (analog `init.go:294-295`).
- [`internal/adapter/driving/cli/erroremission_internal_test.go`](../../../../internal/adapter/driving/cli/erroremission_internal_test.go)
  `TestMapAddErrorToDiagnostic_AllCases`: zwei Fälle (`ErrBackupSuffix-
  Exhausted`, `ErrBackupSourceMissing`) auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) aktualisieren.
- Optional: zusätzlicher Add-Acceptance-Test mit Backup-Failure-Mid-
  Write zur Pinnung der Voll-Envelope (Code [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern), Exit 14).

**Out-of-Scope:**

- `add`-Slice carveouts/done-Plan nicht editieren — der Fix steht in
  diesem Folge-Slice und Cross-Linking via Commit-Message reicht.
- Andere Diagnostic-Code-Drifts (z.B. zwischen `doctor` und
  `generate`) — der Scope dieses Slices ist genau die add↔init
  Backup-Symmetrie.

## Done-Definition

- Exit-Code-Tabelle für `u-boot add` mit Backup-Failure liefert
  Envelope {"diagnostics":[{"code":"[`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern)"}], "exitCode":14}
  (statt heute {"code":"[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)", "exitCode":14}).
- `make gates` bleibt grün; Coverage-Gate weiterhin ≥ 90%.
- CHANGELOG-Eintrag unter `Fixed`.

## Risiko

Minimal. Die Änderung berührt nur die Envelope-Diagnostic-Code-Feld;
Exit-Code, Sentinel-Identität und Cobra-Dispatch bleiben unverändert.
Externe Konsumenten (JSON-Pipeline-Skripte), die auf
diagnostics[0].code == "[`LH-FA-INIT-005`](../../../../spec/lastenheft.md#lh-fa-init-005-überschreibschutz)" für Backup-Failures matchen,
müssen auf [`LH-NFA-REL-003`](../../../../spec/lastenheft.md#lh-nfa-rel-003-abbruch-bei-kritischen-fehlern) umstellen — aber das ist konsistent zum
ohnehin schon korrekten Exit-Code-14-Match und matched semantisch
besser.

## Reihenfolge im Cluster

Außerhalb des [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)-Cluster — ein
Cleanup-Slice, der nicht das `--json`-Schema selbst migriert,
sondern eine Pattern-Konsistenz repariert. Kann unabhängig vom
Cluster-T_close-Lauf gemerged werden.
