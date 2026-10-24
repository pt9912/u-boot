# Verifier-Report: slice-v1-devcontainer-ollama-cloud — 2026-10-24

**Rolle:** Verifier (Modul 8). Frage: Bauen wir es richtig — gegen Plan, DoD, Lastenheft, ADR und Sensor-Belege.

**Gegenstand:** `git diff 8ad49398..HEAD` (Commits `7f082a8a`, `5b983434`, `24711dd3`, `3cd2e700`, `377bc72b`, `490a1a58`), Slice
[`slice-v1-devcontainer-ollama-cloud`](../plan/planning/done/slice-v1-devcontainer-ollama-cloud.md)
(Welle [`welle-ollama-cloud-modelle`](../plan/planning/welle-ollama-cloud-modelle.md)).

**Maßstab:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Absatz `ollama-host`, Fassung 0.23.2),
[`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Ausnahme 1),
[`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Accepted: Entscheidungen 3, 4, 6, 8, O1 bis O3, Z13, Fitness Function),
Review-Report [`2026-10-24-slice-v1-devcontainer-ollama-cloud.md`](2026-10-24-slice-v1-devcontainer-ollama-cloud.md),
`harness/verification.md`, `harness/replay.md`. Nicht Gegenstand: Closure-Notiz (Planner danach).

## Urteil

**Konform.** DoD-Liefer-Punkte 1 bis 3 und der Gate-Punkt sind durch Code, Tests und selbst gefahrene Sensoren belegt. Kein Befund, der die
Closure blockiert; zwei Hinweise (V-1, V-2) und zwei Rest-Beobachtungen (V-3, V-4) unten. Die Review-Befunde F-1 bis F-6 sind behoben
beziehungsweise als bewusst benannt (Tabelle unten). Nicht geprüft: Wirkung gegen einen Dienst des Hosts (Validator).

## Sensoren (selbst gefahren, Belege gelesen)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | Exit 0 | `0 issues` (lint); `coverage-gate: OK — coverage 92.10% meets threshold 90%` (deckt sich mit der Angabe im Slice); `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; Arbeitsbaum danach sauber |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | `9 Commit(s) — OK`, `0 Befund(e)` |
| D1 `TestE2E_OllamaHost_DefaultProfileResolvesHostDockerInternal` (gezielt, `-tags docker -p 1`, `-v /tmp:/tmp`) | PASS, 10,46 s | `host.docker.internal -> 172.17.0.1`; Test prüft `ExtraHosts` exakt, `getent` Exit 0 und die Gegenprobe ohne Feature (Exit 2) |
| D2 `TestE2E_OllamaHost_SandboxWithEgressResolvesAndAddsNoHost` (gezielt) | PASS, 6,56 s | `hosts=`-Zeile des Container-Skripts ohne `ollama.com`, `host.docker.internal`, Adresse |
| Mutationen V01 bis V07 (eigene Patches, `git apply`, Rücknahme `git apply -R`) | je rot, Rücknahme sauber | siehe unten; `git diff` danach 0 Zeilen, `git status` leer |
| `make test-docker` (volles Paket) | nicht gefahren | Auftrag (Speicher, Paket am Zeitlimit); gezielte Läufe ersetzen es nur für die neuen Tests |
| `make verify-depguard` | nicht gefahren | keine Importregel berührt (Diff nur in `application`, Tests, Template; `go list`-Ebene unverändert) |
| `make ci`/`image-scan`/`govulncheck` | nicht gefahren | nicht Teil des Auftrags; kein Release-/Security-Pfad berührt |

Eigene Mutationen (Patch unter dem Scratchpad, Lauf `go test ./internal/hexagon/application/ ./cmd/uboot/` im Image `test-docker-tools`; Exit 1 = rot):

| Nr. | entspricht | Mutation | rote Tests |
|---|---|---|---|
| V01 | M07 (H3) | Template `{{- else if .RunArgs}}` → `{{- else}}` | `…OllamaHost_DefaultProfileRunArgs`, `…OllamaHost_InactiveIsByteIdentical`, bestehend `…SandboxEgress_NotGenerated`, `…SandboxPodman_NotGenerated` |
| V02 | M15 (H7) | `SecurityRelaxation: sandbox` → `false` | `…SandboxNoticeIsRelaxation`, `…EgressRestrictionIsAWarning`, `TestRun_Generate_OllamaHostNotices` (unter `--quiet`) |
| V03 | M22 (H11) | `featureVersionNeedsSource` ohne den neuen Typ | `…VersionWithoutSourceRejected`, `TestConfigSet_OllamaHostVersionWithoutSource_Rejected`, `TestRun_OllamaHostVersionWithoutSourceExits10` (Lese- und Schreibpfad) |
| V04 | M33 (Review F-2) | Filter `ollamaNotices` zurück auf `feature ollama` | `…NoticesDoNotOverlapWithOllama` |
| V05 | M24 (H12) | `hostGatewayFeatureActive` ignoriert `Enabled` | `…NotActive`, `…ReplayAndLaterActivation`, `…InactiveIsByteIdentical` |
| V06 | M29 (H15) | Relay-Satz gestrichen | `…NoticeText` |
| V07 | M03 (H1) | `host-gateway` → `host-gate` | `…DefaultProfileKeyOrder`, `…DefaultProfileRunArgs`, `…ReplayAndLaterActivation`, `…SandboxRunArgs`, `…TogetherWithOtherFeaturesAndNestedRuntime`, `TestRun_Generate_OllamaHostNotices` |

Sieben Mutationen aus 33 reproduziert (V01, V02 und V03 überlappen mit M07, M15 und M22 des Reviewers; V04 bis V07 kommen hinzu); alle Rotfärbungen decken sich mit den
Angaben in §7. Die übrigen Zeilen (M01 bis M06, M08 bis M14, M16 bis M21, M25 bis M32) habe ich **nicht** reproduziert; die Tabelle in §7 und die Testnamen sind gelesen, nicht gefahren.

## DoD, Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Katalogeintrag, Generator, Meldung, Warnung, Tests H1 bis H15 | erfüllt | `featureOllamaHost` (`containerOption`, kein `source`/`defaultVersion`/`publisher`/`relaxations`); `hostGatewayRunArgs` im Standardprofil vor dem frühen `return`, im Sandbox-Profil als letztes Element nach `nestedRuntime` und Egress (`devcontainerTemplateData`); Template `{{- else if .RunArgs}}`; Tests je Zeile H1 bis H15 vorhanden (28 Funktionen `…OllamaHost…`/`TestE2E_OllamaHost_*`), Mutation je Zeile in §7, sieben davon von mir reproduziert |
| 1 bestehende Konfigurationen byte-identisch | erfüllt | kein Golden-Datei-Diff, keine Änderung eines bestehenden Testerwartungswerts (nur Filter in `generate_ollama_test.go`, Zeile `ollama-host` in `TestFeatureCatalogue_KeysCoverSpecExamples`); V01 färbt die bestehenden Sandbox-Goldens und den neuen Literalpin; `initproject.go` unverändert, `hostGatewayFeatureActive(nil)` falsch |
| 1 Golden Case nach `replay.md` | erfüllt | Zeile in `harness/replay.md` (Fresh-State `…DefaultProfileRunArgs`/`…SandboxRunArgs`, Double-Run-NoOp und veralteter Block/Nutzerzeile `…ReplayAndLaterActivation`, kein Teilschreiben bei Exit 10 `TestRun_OllamaHostVersionWithoutSourceExits10` mit `--dry-run`/`--diff`); die genannten Testnamen existieren |
| 2 Lauf belegt (D1, D2) | erfüllt | von mir gefahren (Tabelle oben); Gegenprobe ohne Feature; keine Netzverbindung zu `ollama.com`, kein Konto, kein Token im Test |
| 3 kleine Doku | erfüllt | Handbuch: eine Katalogzeile `ollama-host` und `ollama-host` in der Konfigurationstabelle (`version` nur mit `source`); README.md/README.de.md: 17 statt 16 (Aufzählung, Satz zur Ablehnung von `version`), keine verbliebene „16“ im Baum (Suche in README, `docs/user`, `docs/maintainer`, CHANGELOG ohne Treffer); `CHANGELOG.md` `[Unreleased]` Added/Changed; keine Kennung in den Nutzerdoku-Diffs (`LH-`, `ADR-`, `SPEC-`, `slice-`, `welle-` ohne Treffer in den `+`-Zeilen von README, README.de.md, `docs/user`) |
| 4 `make gates` grün, Docker gezielt | erfüllt | selbst gefahren |
| 5 bis 9 (Review, Doku-Update, Closure, Register, Risiken, Paarungen) | offen, nicht Gegenstand | Review-Report liegt vor (Punkt 5 durch den Reviewer erfüllt); Closure-Pflichten gehören dem Planner/Implementer |

## Wirkung gegen den Vertrag

- **Eine Zeile, beide Profile, nichts sonst.** Der Eintrag liefert genau `--add-host=host.docker.internal:host-gateway` (`hostGatewayOption`). Es gibt keinen Code, der `containerEnv`, `forwardPorts`, `mounts`, `remoteEnv`, `EgressHosts` oder `agentCredentialEnv` für den Eintrag berührt (Diff gelesen). Die Abwesenheit pinnen `…WritesNoVariableAndNoRelay`, `…AddsNoEgressHost`, `…AddsNoCredentialVariable`, `TestDoctor_OllamaHost_NoNewFinding`, `TestRun_EgressList_OllamaHostAddsNoHost`; `egressFeatureHosts` ist unverändert (nur Kommentar). Deckt [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) Absatz `ollama-host` und Z13 der ADR. Die Fitness-Function-Zeile „kein neuer Pfad in der Konfigurations-Registry“ ist durch Abwesenheit im Diff erfüllt (keine Registry-Datei berührt), aber durch keinen eigenen Test gepinnt (V-3).
- **Byte-identisch ohne Feature.** Template-Zweig nur bei nicht leerer Liste; Literalpin „kein `runArgs`-Schlüssel“ je Fall; V01 rot.
- **Meldung.** Eine `WarningEntry` je Lauf; Sandbox: Code der Sandbox-Anforderung, `SecurityRelaxation: true`, unter `--quiet` sichtbar (`TestRun_Generate_OllamaHostNotices`, Nutzerpfad mit echten Adaptern, vier Profil-Fälle); Standardprofil: Code der Feature-Anforderung, ungeflaggt, unter `--quiet` entfällt sie. Egress-Zusatz nur bei Sandbox **und** `egressEnabled`, im selben Eintrag, Exit 0 (V02 und die Matrix). Der Text nennt `host.docker.internal`, die Bedingung (Loopback, Relay oder Bindung des Nutzers), „starts no relay“, nichts zu `ollama.com`, Cloud-Modellen oder Kennungen (`…NoticeText`, V06 rot). Entscheidung 8 („eine Meldung, auf das Gemessene begrenzt, Warnung kein Abbruch“) ist erfüllt.
- **Exit 10 für `version` ohne `source`.** Eine Regel `featureVersionNeedsSource`, gefragt von `validateDevcontainerFeatures` (Lesen, `ErrInvalidFeatureSource`) und `revalidateFeatureEntry` (Schreiben, `ErrConfigValueInvalid`); Nutzerpfad `generate`, `generate --dry-run`, `config set`, `--dry-run`, `--diff` jeweils Exit 10, Datei unverändert, Text ohne „Dockerfile“ (V03 rot in beiden Pfaden). Kein neuer Sentinel, die Dual-Classifier-Regel entfällt.
- **Aussagen in Handbuch/README/CHANGELOG gegen Code gelesen (`AGENTS.md` §3.10).** Handbuch-Zeile („eine Zeile in `runArgs`“, „Standard- und Sandbox-Profil“, „weist beim Erzeugen aus“, „Sicherheitslockerung auch unter `--quiet`“, „kein Relay, keine Variable, kein Host in der Egress-Liste“, „`version` nur mit `source`“) und CHANGELOG-Aussagen („Exit 10“, „kein Schlüssel unter `features`“, „byte-identisch ohne das Feature“, „that a service of the host is reachable is not tested“) haben je einen Test (oben). Der Rest bleibt wie geplant: Handbuch „ein Anschluss an einen Ollama-Dienst des Hosts … gibt es nicht“ und „Cloud-Modelle sind nicht enthalten“ stehen neben der Katalogzeile — der in [`carveouts.md`](../plan/planning/in-progress/carveouts.md) als Carveout 1 geführte, bis zum Spezifikations-Slice dokumentierte Widerspruch (V-1).
- **Docs-Straten.** `spec/` unverändert (`git diff --stat` ohne `spec/`), `harness/replay.md` mit Golden-Zeile, ADR 0036 unberührt (Accepted, immutable).

## Review-Befunde F-1 bis F-6

| ID | Stand | Beleg |
|---|---|---|
| F-1 (Plan-Tabelle `config.go`) | behoben | Plan-Tabelle zeigt `config.go` als update; Abgleich unten |
| F-2 (Filter `ollamaNotices`) | behoben | Filter `feature ollama `; `…NoticesDoNotOverlapWithOllama`; V04 rot |
| F-3 (Wortlaut nahe Zusage-Grenze) | INFO, bewusst | Text formuliert eine notwendige Bedingung, folgt Entscheidung 8; Lesepunkt für Validator und Spezifikations-Slice |
| F-4 (Literalpin an Fallnamen) | behoben | Feld je Fall in `…InactiveIsByteIdentical`; V01 rot |
| F-5 (unerreichbarer Zweig) | behoben | `hostGatewayFeatureActive` ohne `containerOption == ""`-Zweig (Diff gelesen) |
| F-6 (Zwischenwiderspruch, Komma) | Komma behoben, Widerspruch bewusst | README.de.md „`glab` …, `ollama` und `ollama-host`“ mit Komma; Widerspruch = Carveout 1 |

## Plan-Tabelle gegen `git diff --stat 8ad49398..HEAD`

Diff-Dateien (21), alle mit Plan-Eintrag: `devcontainer_features.go`, `devcontainer_sandbox.go`, `templates.go`, `devcontainer.json.tmpl`, `generate.go`, `devcontainer_egress.go`, `config.go`, `generate_ollama_host_test.go` (neu), `devcontainer_features_test.go`, `export_test.go`, `generate_ollama_test.go`, `cmd/uboot/main_test.go`, `internal/e2e/ollama_host_docker_test.go` (neu), `harness/replay.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `roadmap.md`, `docs/reviews/` (Review-Report; mit diesem Report zwei), und der Slice selbst (§7-Beleg, DoD-Haken).

Plan-Einträge ohne Diff in dieser Range: `welle-ollama-cloud-modelle.md`, `open/README.md`, `carveouts.md` (Plan: „update“) — ihre Angleichung liegt vor der Range (`89cc1a82` Anlage, `8ad49398` Verweise); `doctor_sandbox_test.go`, `doctor_features_drift_test.go` (Plan: „update (Plan) / unverändert (Ist)“, Abweichung in §7 genannt). Alle übrigen „unverändert“-Zeilen (`doctor*.go`, `initproject.go`, `warnings.go`, Dockerfile-/Skript-Templates, `spec/`, Makefile, `.golangci.yml`, `internal/adapter/driven`) erscheinen nicht im Diff. Keine Zeile ohne Plan-Eintrag, kein stiller Eintrag ohne Diff außer den drei Planungsdateien oben. Die Zeile „Plan-Abgleich nach der letzten Fix-Runde“ in §7 trägt der Slice als Abgleich in *Abweichungen* (Commit `490a1a58`), der Stand deckt sich mit meinem Abgleich.

## Hinweise und Rest-Beobachtungen

- **V-1 (Zwischenzustand, bewusst).** Handbuch und READMEs sagen für Cloud-Modelle/Host-Anschluss weiter „nicht enthalten“, neben der neuen Katalogzeile; [`SPEC-044`](../../spec/spezifikation.md#spec-044--features-sshd-und-claude-code-im-devcontainer-warnungen-und-zugangsdaten) und [`SPEC-047`](../../spec/spezifikation.md#spec-047--feature-ollama-im-sandbox-devcontainer-installation-volume-meldungen-und-grenzen) sagen noch das Alte. Das ist Carveout 1 mit Anker auf diesen Slice und den Spezifikations-Slice; der Release liegt nach diesem. Bei der Closure: Teil „Katalog im Code“ aus dem Carveout nehmen, Widerspruch zur Sandbox-Regel der Spezifikation ergänzen (§5 des Slice).
- **V-2 (Form der Mutationstabelle).** §7 trägt Befehl und Wiederherstellung als **Gesamtregel** statt je Zeile (Plan §3 *Mutationsprotokoll* verlangt je Mutation). Die Angabe ist belegt für die sieben von mir und die vier vom Reviewer reproduzierten Zeilen; die Patches/Logs der übrigen liegen nicht im Repository (Scratchpad des Implementers), die Zeilen sind daher nicht aus dem Repo allein nachprüfbar. Kein DoD-Verstoß, aber die Closure sollte die Gesamtregel als gewählte Form benennen.
- **V-3 (Fitness-Function-Zeile ohne Pin).** „Kein neuer Pfad in der Konfigurations-Registry“ (ADR 0036, Fitness Function, Zeile Unit): erfüllt, weil der Diff keine Registry-Datei berührt; es gibt keinen Test, der das festhält. Der Slice-Test-Auftrag (H1 bis H15) enthält keine eigene Zeile dafür. Hinweis an den Planner, nicht Blocker.
- **V-4 (Umgebung).** Während der Läufe war ein fremder Container (`eloquent_jackson`, Image `golang`, Mounts aus dem Verzeichnis eines anderen Repos) auf dem Docker-Host aktiv; er gehört nicht zu diesem Lauf und wurde nicht angefasst. Die Docker-Tests räumten ihre Container (Zahl vor und nach den Läufen gleich, 1 fremder Container).

## Traceability und Carveouts

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | `TestGenerateDevcontainer_OllamaHost_*`, `TestFeatureCatalogue_*`, `TestRun_OllamaHostVersionWithoutSourceExits10`, D1, D2 |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | `…SandboxNoticeIsRelaxation`, `…EgressRestrictionIsAWarning`, `…SandboxRunArgs` |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `TestRun_Generate_OllamaHostNotices` (Sichtbarkeit je Profil unter `--quiet`) |
| [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) | Entscheidungen 3, 4, 6, 8, O1 bis O3 und Z13 durch die oben genannten Tests; Fitness Function: alle drei Tabellenzeilen belegt (V-3 für die Registry-Aussage) |

Carveouts: neu keine; unverändert Carveout 1 (`carveouts.md`), siehe V-1. Neue Gate-Lockerung: keine (`.golangci.yml`, `Makefile`, Schwellen unverändert, kein `//nolint`).

Replay / Golden:
- Affected generators: `generate devcontainer` (Standard- und Sandbox-Profil)
- Golden cases added: Standard `runArgs` und Schlüsselfolge, Sandbox `runArgs` (allein, mit Egress), Kombination `ollama`+`gh`+`glab`+`ollama-host` mit `nestedRuntime` podman/docker und Egress, inaktiv byte-identisch, Replay/späte Aktivierung
- Golden cases updated: keine
- Golden cases replayed: alle bestehenden Golden-Tests im Paket (`make gates` grün, V01 färbt bestehende Sandbox-Goldens rot)
- Intentional output changes: nur bei aktivem `ollama-host` (neue `runArgs`-Zeile); ohne das Feature keine

## Nicht ausgeführt

- Volles `make test-docker`, `make verify-depguard`, `make ci` — Gründe oben.
- Mutationen außer V01 bis V07 — Aufwand; Angaben aus §7 gelesen, nicht reproduziert.
- Wirkung gegen einen Dienst des Hosts, Relay, Cloud-Modell — ausdrücklich nicht zugesagt (ADR *Nicht gemessen*), Aufgabe des Validators.

## Summary

DoD-Punkte 1 bis 4 erfüllt und von mir bestätigt (`make gates` Exit 0, `make doc-commits` Exit 0, D1 und D2 PASS, sieben Mutationen rot mit sauberer Rücknahme). Review-Befunde F-1 bis F-6 erledigt oder bewusst benannt. Offen für den Planner: V-1 (Carveout-Fortschreibung bei der Closure), V-2 (Form der Mutationstabelle), V-3 (Registry-Aussage ohne Pin).
