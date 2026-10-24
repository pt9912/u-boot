# Verify-Report: slice-v1-generate-prueft-freigabeliste

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff 750139e5..HEAD` (HEAD `5b5016b9`), Slice
  [`slice-v1-generate-prueft-freigabeliste`](../plan/planning/in-progress/slice-v1-generate-prueft-freigabeliste.md),
  Review [`2026-10-24-slice-v1-generate-prueft-freigabeliste`](2026-10-24-slice-v1-generate-prueft-freigabeliste.md) (LOW-1 bis LOW-3 in `77083695` und `5b5016b9` behoben)
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
* Frage: Bauen wir es richtig (Plan und DoD)? Reparaturen: keine. Der Baum war vor und nach jeder Probe sauber (`git status --short` leer; Mutationen per Diff angewandt und mit `git apply -R` zurückgenommen, nie `git checkout`).

## Urteil

**Schließbar** für den Code- und Doku-Teil. Offen sind allein Pflichten, die der Planner bei der Closure erfüllt (Closure-Notiz mit Steering-Loop-Lerneintrag, Register, Risiko-Ausgänge §6, drei Paarungen, Nachlauf-Review-Vermerk, Move nach `done/`); sie sind kein Code-Mangel.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit-Code 0 (`EXIT=0` im Log); docs-check 0 Befunde; Coverage-Gate „92.10% meets threshold 90%“ |
| `make doc-trace` | Exit-Code 0, „84 Anforderung(en), 0 Waise(n)“ |
| `make doc-planning` | Exit-Code 0, 798 Dateien, 0 Befunde |
| `make doc-commits RANGE=750139e5..HEAD` | Exit-Code 0, 8 Commits, 0 Befunde |
| `make doc-immutable RANGE=4655bcff..HEAD` | Exit-Code 0, 105 Commits, 0 Befunde |
| Gezielte Tests (`u-boot:test-docker-tools`, Repo gemountet) | grün: `Freigabe\|ThreeCallersAgree\|UnlistedFeatureSource` (14 PASS in Application und `cmd/uboot`), dazu Bestandstests der berichtigten Fixtures (`Homebrew`, `ClaudeCodeNodeNotice`, `HostSocketFeature`, `SandboxFeatureWarnings`) |
| `make build` | Exit-Code 0 (Image `u-boot:latest`, `--version` 0.15.0-dev) |

Nicht gelaufen, mit Grund: `make test-docker` (kein Docker-Pfad berührt, Anweisung Speicher schonen), `make verify-depguard` (keine Importregel berührt, keine neuen Importe im Diff), `make ci`/`image-scan`/`govulncheck` (nicht Teil der DoD).

## Eigene Mutationen

Jede per Patch angewandt, Tests `Freigabe|ThreeCallersAgree|UnlistedFeatureSource` in `internal/hexagon/application` und `cmd/uboot`, danach `git apply -R`.

| Mutation | Datei | Rot gesehen |
|---|---|---|
| M1 Prüfung aus (`if false && len(unlisted) > 0`) | `generate.go` | `FreigabeRejectsUnlistedSource`, `BeforeSandboxRejection`, `RejectsAlsoWithUnrelatedFlag`, `SeveralUnlistedSortedByName`, `MessageParts`, `SandboxProfileOnlyUnlistedSource`, `ThreeCallersAgree` (alle Teilfälle) |
| M3 Flag nicht gemerged | `devcontainer_features.go` | `FreigabeByFlagOfSameCall`, `FreigabeByFlagDryRun`, `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| M4 `enabled`-Filter weg | `devcontainer_features.go` | `FreigabeIgnoresInactiveEntries`, `TestRun_Generate_UnlistedFeatureSourceIsRejected` |
| M13 Präfix-Vergleich im gemeinsamen Prädikat | `config.go` | `ThreeCallersAgree/trailing_slash_in_source`, `…/query_in_source` |
| M10 Sandbox-Profil ausgesetzt | `generate.go` | `BeforeSandboxRejection`, `SandboxProfileOnlyUnlistedSource`, `TestRun_…IsRejectedInSandboxProfile` |
| M7a Zweig `ErrConfigValueInvalid` in `mapGenerateErrorToDiagnostic` weg | `cli/generate.go` | `TestRun_Generate_UnlistedFeatureSourceIsRejected`, `…InSandboxProfile` |
| M7b Sentinel nicht im `ExitCode` (`isConfigValidationError`) | `cli/cli.go` | `TestRun_Generate_UnlistedFeatureSourceIsRejected`, `…InSandboxProfile` |

Nicht von mir nachgefahren (Belege des Implementers und des Reviewers, nicht eigen gesehen): M1b, M2, M5, M6, M8a bis M8g, M9, M11, M12.

## Ende-zu-Ende am gebauten Binary

`u-boot:latest` aus `make build`, `docker run` mit Projektmount und Host-Benutzer, Eintrag `evil` mit `enabled: true` und `source: https://evil.example.org/feat`, keine `featureSources.allow`.

| Fall | Standardprofil | Sandbox-Profil |
|---|---|---|
| `generate devcontainer` | Exit 10, Meldung nennt Feature, Quelle, Freigeben, Abschalten, `u-boot.yaml`-Hash gleich, kein `.devcontainer` | Exit 10, dasselbe |
| `--dry-run` | Exit 10, nichts geschrieben | Exit 10, nichts geschrieben |
| `--json` | Code [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), `exitCode` 10 | Code [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), `exitCode` 10 |
| Abhilfe `config set …featureSources.allow` | `config set` Exit 0, `generate` Exit 0, Quelle in `devcontainer.json` | dasselbe |
| Abhilfe `--allow-external-feature-sources` | Exit 0, Quelle danach in der Liste | Exit 0, Quelle in der Liste, Warnung [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) bleibt |
| Abhilfe `…enabled false` | Exit 0, kein Schlüssel | Exit 0, kein Schlüssel |

Zusätzlich: Katalogschlüssel `node` ohne `source` → `generate` Exit 0. `config set …evil.source <unlisted>` → Exit 10 wie zuvor. `doctor` meldet `devcontainer.features.allowlist` weiter als Fehler.

## DoD, Punkt für Punkt

| DoD | Befund |
|---|---|
| 1 Freigabeprüfung (`[x]`) | Erfüllt. Exit 10 und Code [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) im Fehlerpfad und im `--json`-Umschlag, nichts geschrieben auch unter `--dry-run`, wirksame Freigabe = Liste plus Flag desselben Aufrufs ohne Vorab-Schreiben (`unlistedFeatureSources`, `normaliseFeatureSources`). Meldungsreihenfolge Freigeben, Abschalten, Vergleichssatz; beide Abhilfen am Binary ausgeführt. Ein Prädikat `featureSourceInAllow` trägt `revalidateFeatureEntry`, `classifyFeatureEntries` und die neue Prüfung; `doctor.go` im Diff nicht berührt, M13 färbt `ThreeCallersAgree`. Prüfung steht vor der Sandbox-Abweisung (`generate.go` nach `collectDevcontainerFeatures`). |
| 2 Vertragstext (`[x]`) | Erfüllt. [`SPEC-050`](../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) Absatz *Freigabe der Quelle (beide Profile)*, Aufzählung im Abschnitt zu [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), Historie-Zeile; Handbuch-Absatz ohne Kennungen (Diff der Handbuchzeilen auf `LH-`, `ADR-`, `SPEC-` gesucht: keine Treffer); Handbuchsatz „doctor meldet sie weiter“ gestrichen. CHANGELOG `[Unreleased]`: eine neue Zeile unter *Fixed*, englisch, sonst keine Zeile geändert, `[0.15.0]` und frühere Texte unberührt. Lastenheft, `docs/plan/adr/`, `spec/architecture.md`, `Makefile`, `.golangci.yml`, `.d-check.yml`, `README.md`, `README.de.md` laut `git diff --stat` unberührt (§3.4 gewahrt: `architecture.md` unverändert). Tag `v0.15.0` zeigt weiter auf `a895ffd2`. |
| 3 Gates (`[x]`) | Erfüllt: `make gates` Exit 0, docs-check 0 Befunde, `doc-trace` ohne Waise, `doc-commits` 0 Befunde (Range ab `750139e5` statt `origin/main`, gleiche Commits des Slice). `make test-docker` entfällt wie begründet. |
| 4 Review | Report liegt vor; LOW-1 bis LOW-3 behoben und im Slice belegt. Der Haken ist offen, solange der Nachlauf-Vermerk fehlt: Die Fix-Runde änderte Tests und einen Handbuchsatz, **nicht** Spezifikation (`git diff 84d943ea..HEAD --stat` lt. Slice: Tests, Handbuch, Slice-Datei), damit greift die Nachlauf-Review-Pflicht für Vertragstext nur für den Handbuchsatz (gestrichen, verkleinert). |
| 5 bis 8 | Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen: ausstehend, Aufgabe des Planners; §6 trägt noch „offen bis zur Closure“, §7 den Platzhalter. |

Plan-Tabelle gegen Diff: Die 14 Dateien des Diffs sind in §3 bzw. im Umsetzungsstand geführt (`harness/replay.md` und `roadmap.md` mit gelesener Begründung). Vier berichtigte Fixtures (`FeaturesHomebrew_ExternalOverrideNoWarning`, `ClaudeCodeNodeNotice` Fall „own source“, `SandboxHostSocketFeature_OnlyActiveEntryIsRejected` Fall „source override“, `TestRun_Generate_SandboxFeatureWarnings`): jede ergänzt die Freigabeliste um die tatsächliche Quelle; der Vergleich (`featureSourceInAllow`) ist im Diff nur um einen Kommentar geändert, kein Test überspringt oder schaltet die Prüfung ab (im Diff kein `Skip` außer dem bestehenden `testing.Short()`-Muster).

## Ungemessen

* Whitespace am Rand einer Quelle in `u-boot.yaml` (der Slice nennt: nicht gepinnt, der YAML-Parser schneidet).
* `doctor` bei fehlendem `enabled` (meldet `enabledKeyMissing`, kein Doctor-Test für die Freigabe-Verletzung; der Handbuchsatz dazu ist gestrichen).
* `make test-docker`, Image-Scan, `govulncheck`.
* Ob die Dev-Container-CLI einen Schlüssel einer freigegebenen externen Quelle baut (N2 des Ursprungs-Slice, nicht dieser Slice).
* Das Image `u-boot:latest` kam aus dem Docker-Build-Cache (Schichten unverändert); der Quellstand entspricht HEAD, da `COPY . .` bei Änderung neu gebaut hätte. Das Verhalten stimmt mit dem Code (Exit 10 gesehen).
* Reihenfolge der Abhilfen 1 und 2 in der Meldung: Pin steht laut Slice im Test, Rotsehen nicht gefahren (auch nicht von mir).
