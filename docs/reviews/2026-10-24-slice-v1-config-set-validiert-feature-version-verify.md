# Verify-Report: slice-v1-config-set-validiert-feature-version

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff a67acfc0..HEAD` (HEAD `7f355889`), Slice
  [`slice-v1-config-set-validiert-feature-version`](../plan/planning/in-progress/slice-v1-config-set-validiert-feature-version.md),
  Review [`2026-10-24-slice-v1-config-set-validiert-feature-version`](2026-10-24-slice-v1-config-set-validiert-feature-version.md)
* Frage: Bauen wir es richtig (Plan und DoD)? Reparaturen: keine, der Baum war vor und nach jeder Probe sauber (`git status --short` leer).

## Urteil

**Schließbar** für den Code- und Doku-Teil. Offen sind allein Pflichten, die der Planner bei der Closure erfüllt und die kein Code-Mangel sind (siehe DoD-Tabelle: Closure-Notiz, Register, LOW-1 des Reviews).

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit-Code 0 (`EXIT=0` im Log; Lint, Test, Coverage-Gate, docs-check, baseline-verify, ablage-pruefung) |
| `make doc-trace` | Exit-Code 0, „84 Anforderung(en), 0 Waise(n)“ |
| `make doc-planning` | Exit-Code 0, 790 Dateien, 0 Befunde |
| `make doc-immutable RANGE="4655bcff..HEAD"` | Exit-Code 0, 89 Commits, 0 Befunde (Accepted-ADRs unverändert) |
| Gezielte Tests (`u-boot:test-docker-tools`, `-run`) | grün: `DockerfileBlockVersionWithoutSource_Rejected`, `VersionWithSourceOrWithoutDockerfileBlock_Succeeds`, `SetMany_DockerfileBlockVersionNeedsSourceFirst`, `SetMany_RejectsFileWithForgeVersionWithoutSource`, `ForgeVersionWithoutSource_*` |
| `make build` | Exit-Code 0 (Binary für die Ende-zu-Ende-Probe) |

Nicht gelaufen, mit Grund: `make test-docker` (kein Docker-Pfad berührt, so auch in der DoD), `make verify-depguard` (keine Importregel berührt, keine neuen Importe in `internal/hexagon/application`), `make ci`/`image-scan`/`govulncheck` (nicht Teil der DoD dieses Slice). Die Coverage-Zahl des Gates habe ich nicht einzeln ausgewertet; das Gate selbst war grün.

## Mutationen (eigene; Rücknahme per `git apply -R`, Baum danach sauber)

| # | Mutation | Ergebnis |
|---|---|---|
| M1 | Prüfung in `revalidateFeatureEntry` aus (`if false && …`) | rot: `TestConfigSet_DockerfileBlockVersionWithoutSource_Rejected`, `TestConfigSetMany_DockerfileBlockVersionNeedsSourceFirst` |
| M2 | Aufruf `revalidateConfigDomain` in `SetMany` (`applySetItem`) entfernt | rot: `TestConfigSetMany_DockerfileBlockVersionNeedsSourceFirst`, `TestConfigSetMany_OrderedDependency` |
| M3 | `dockerfileBlock`-Bedingung im Helfer `featureVersionNeedsSource` entfernt | rot: `TestConfigSet_VersionWithSourceOrWithoutDockerfileBlock_Succeeds` (`java`) und neun weitere |
| M4 | `Source == ""` im Helfer entfernt (Ablehnung auch mit `source`) | rot: `…VersionWithSourceOrWithoutDockerfileBlock_Succeeds`, `…VersionNeedsSourceFirst` |
| M5 | Reihenfolge-Semantik: `SetMany` prüft Paare mit `version` nicht pro Paar | rot: `TestConfigSetMany_DockerfileBlockVersionNeedsSourceFirst` (Fall „version before source“) |

Alle fünf werden gefangen. Der Test der Preview-Modi (`--dry-run` = `PreviewDryRun`, `--diff` = `PreviewAndApply`) liegt auf Service-Ebene; M1 färbt ihn rot.

## Ende-zu-Ende am gebauten Binary (`u-boot:latest`, `docker run` mit Projektverzeichnis-Mount, ohne Host-Toolchain)

| Probe | Ergebnis |
|---|---|
| `config set devcontainer.features.gh.version 2.0.0` | Exit 10, Meldung nennt Pfad und `…gh.source` als Weg, Datei byte-identisch (sha256) |
| dasselbe mit `--dry-run`, `--diff`, `--dry-run --diff` | je Exit 10, Datei unverändert (schließt INFO-2 des Reviews am Binary) |
| `config set devcontainer.features.java.version 21` ohne `source` | Exit 0, Eintrag geschrieben |
| Mehrere Paare `featureSources.allow`, `gh.source`, `gh.version` in einem Aufruf | Exit 0, alle drei Einträge in der Datei |
| Umgekehrte Reihenfolge (`glab.version` vor `glab.source`) | Exit 10, Datei unverändert |
| Lesepfad: gesperrte Datei (`gh.version` ohne `source`, von Hand): `generate devcontainer` | Exit 10 (Text „managed block missing or malformed: …“ unverändert, vom Slice ausdrücklich ausgeklammert) |
| dieselbe Datei: `config set project.name x` | Exit 10 (`schema validation failed`) |
| dieselbe Datei: `doctor` | Exit 11 (Fehler in `uboot.yaml.valid`) |

## DoD Punkt für Punkt

| DoD-Punkt | Befund |
|---|---|
| Liefer-Punkt 1, Code: `Set` und `SetMany` lehnen `version` ohne `source` für `ollama`, `gh`, `glab` mit Exit 10 ab, Datei byte-identisch, auch `--dry-run`/`--diff`; mit `source` gelingt es, Reihenfolge `source` vor `version`; Mutation färbt rot; T4 des Schwester-Slice gedreht | erfüllt. Eine Regel an einer Stelle (`featureVersionNeedsSource`), beide Pfade fragen sie (`devcontainer_features.go` Zeile 641, `config.go` Zeile 586). Sentinel `ErrConfigValueInvalid` → Exit 10 (kein Klassifikator berührt, der Dual-Classifier-Befund entfällt). Mutationen M1–M5, Binary-Probe siehe oben. Der alte Charakterisierungstest `…SetItselfSucceeds` ist durch `…Rejected` ersetzt. |
| Liefer-Punkt 2, Doku: [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Wertetabelle, Stufe 6), [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) (Absatz `version` und `source`), Handbuch (A11 zurückgenommen, Pfad-Tabelle), jede Bedingung mit pinnendem Test, Lastenheft unverändert | erfüllt. `spec/lastenheft.md` nicht im Diff. Bedingungen und Pins: Ablehnung ohne `source` (M1), „nur Dockerfile-Block“ (M3), „mit `source` gelingt“ (M4), Reihenfolge (M5), `ollama`/`gh`/`glab` je Fall in beiden Testfunktionen. Die Handbuch-Aussage „auch mit `--dry-run` und `--diff`“ ist am Binary gemessen und im Service-Test gepinnt. |
| `make gates` grün (inkl. `docs-check`, `doc-trace`), Fehlen von `make test-docker` in der Closure-Notiz | Gates grün, selbst gefahren. Die Begründung für das Fehlen von `test-docker` ist in diesem Bericht dokumentiert; **die Closure-Notiz steht noch aus** (Planner). |
| Review durchgeführt, Report unter `docs/reviews/` | erfüllt (Report liegt vor, Kein HIGH/MEDIUM). |
| Doku-Update: CHANGELOG `[Unreleased]` *Fixed*, Handbuch, Spezifikation | erfüllt. `[Unreleased]` trägt nur `### Fixed`, kein Zeilenabzug im CHANGELOG-Diff; die Texte unter `[0.15.0]` sind unberührt. |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** (§7 trägt noch den Platzhalter; Planner bei der Closure) |
| Beobachtungs-Register fortgeschrieben (oder „keine Beobachtung“ in §7) | **offen** (Planner bei der Closure; die Klasse „Prüfung beim Lesen, nicht beim Schreiben“ ist der vom Plan genannte Kandidat) |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen** (Planner; Ausgangslage: Reihenfolge-Risiko durch M5 und Binary entfallen, V-3 durch `…SetMany_RejectsFileWithForgeVersionWithoutSource` entfallen, gesperrte Bestandsdatei bleibt offen, L-2/L-3 und Wortlaut `generate` weiter offen) |
| Drei Paarungen | mit der Closure zu prüfen (Planner) |

## Entscheidungslage

* [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (Schema-Konformität beim Setzen) und [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (ungültiger Wert: Exit 10): vom Code jetzt gehalten; keine Lastenheft-Änderung, keine Projektinhaber-Entscheidung berührt.
* [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) Stufe 6 und Wertetabelle sowie [`SPEC-051`](../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) tragen die Regel; der Verweis der zweiten auf die erste Spec-Stelle geht im selben Stratum, nicht abwärts.
* Hard Rules: §3.4 gewahrt (`spec/architecture.md` nicht im Diff), §3.5 gewahrt (`doc-immutable` grün), kein `//nolint`, kein Gate gelockert (`.golangci.yml`, Makefile nicht im Diff), §5 Regel 6: der Handbuch-Diff trägt keine Kennung und keinen Verweis auf Planungsdokumente.
* Tag `v0.15.0` zeigt weiter auf `a895ffd2320cccb11916051879b1e61ecf65129d`.

## Plan-vs-Code-Diff

Jede Code- und Doku-Zeile der Plan-Tabelle (§3) hat ihren Diff: `config.go`, `devcontainer_features.go`, `config_features_test.go`, `config_many_test.go`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `roadmap.md`, `open/README.md`, Review-Report. Die Zeilen „unverändert“ (`cli/config_test.go`, `generate_forge_cli_test.go`, `spec/lastenheft.md`, `spec/architecture.md`) stehen nicht im Diff. Abweichung: die Plan-Zeile `cli/config_test.go` führt noch den ursprünglichen Eintrag als Zusatz; der Test auf Adapter-Ebene (B1 „`ExitCode`-Pin am Adapter“, B3) entstand nicht neu, weil der bestehende Pin `ErrConfigValueInvalid` → 10 trägt. Dieses Verhalten ist durch den Binary-Lauf belegt.

## Findings

Kein HIGH, kein MEDIUM, kein neues LOW.

* **Bestätigt, nicht neu:** LOW-1 des Reviews (Plan-Tabelle nennt den Verweis-Nachzug unter `done/` und die zwei Evidence-Dateien nicht; Report-Dateiname 10-24, Kopf nennt 10-23). Der Verweis-Nachzug unter `done/` ist nach Diff-Lesen rein ein Pfadwechsel `open/` → `in-progress/` in Link-Zeilen; kein Inhalt geändert. Behebung durch den Planner bei der Closure.
* **INFO (Lücke, ungemessen als Automat):** Es gibt keinen Test, der den CLI-Befehl `config set … --dry-run` durchlaufen lässt; die Zusage ist am Binary einmal gemessen (siehe oben) und auf Service-Ebene gepinnt. Das schließt INFO-2 des Reviews für diese Closure, ersetzt aber keinen Regressionstest.
* **INFO:** Der Slice-Text sagt in B1 „`ExitCode`-Pin am Adapter“; ein neuer Pin wurde nicht angelegt (bestehender Pin trägt, siehe oben).

## Ungemessen

* `--json`-Ausgabe der neuen Ablehnung (`data.hint`): nicht geprüft; der Slice sieht keinen Hinweis vor.
* Mac, Podman, andere Hosts: nicht gemessen; der Fix liegt in der Anwendungsschicht ohne Plattformpfad.
* `make test-docker`, `make verify-depguard`, `make ci`: nicht gelaufen (Gründe oben).
