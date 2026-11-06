# Verification-Report: slice-v1-devcontainer-remove — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-devcontainer-remove`](../plan/planning/done/slice-v1-devcontainer-remove.md) — Commit-Range `f5567aee..6d8f521a` (HEAD, zwei Commits: `c750ff31` Umsetzung, `6d8f521a` Fix-Runde nach dem Review; 37 Dateien, +3421 −154). Der Slice liegt in `in-progress/`.

**Eingangs-Kontext:**

- Slice-Plan (Kopf, §1 bis §6, §8), DoD (sieben inhaltliche Zeilen, vier Pflichten der Closure), Plan-Tabelle §3
- [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) (ganz, mit der Tabelle *Belege* in zwölf Zeilen), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Lastenheft `0.37.1` bis `0.37.3`), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- Review-Report (Eingang, nicht Evidence): [`2026-11-06-slice-v1-devcontainer-remove`](2026-11-06-slice-v1-devcontainer-remove.md) (0 HIGH / 1 MEDIUM / 2 LOW / 3 INFO, merge-blockierend wegen F-1)
- Diff der Fix-Runde `c750ff31..6d8f521a` (sieben Dateien) gegen F-1, F-2, F-3, F-5 des Reviews
- Tests und Docker-Test (Quelltext gelesen, dann ausgeführt), `Dockerfile` (Runtime-Stage)
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.9, §3.10, §5)

**Eingang des Implementers:** Der Auftrag trug keine Sensor-Belege des Implementers und kein Handoff mit den Fundstellen je Zeile der Tabelle *Belege*. Der Stempel unter `.harness/state/` nannte vor diesem Lauf den Commit `c750ff31`, nicht HEAD. Alle Sensoren dieses Reports sind eigene Läufe des Verifiers auf `6d8f521a`; die Fundstellen der Tabelle *Belege* hat der Verifier selbst aus den Tests zusammengesetzt (Abschnitt Traceability).

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-remove`](../plan/planning/done/slice-v1-devcontainer-remove.md) (Welle [`welle-devcontainer-remove`](../plan/planning/done/welle-devcontainer-remove.md))
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image), [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Artefakte: `u-boot devcontainer remove` (`internal/adapter/driving/cli/devcontainer.go`, `internal/hexagon/application/devcontainer_remove.go`, die vier Port-Dateien, `internal/adapter/driven/docker/`, `internal/adapter/driven/confirm/`, `cmd/uboot/main.go`), `internal/e2e/devcontainer_remove_docker_test.go`, Handbuch, README, CHANGELOG, `spec/spezifikation.md`, `spec/architecture.md`, `docs/maintainer/`

### DoD-Abgleich (Plan §2)

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | Befehl und Ports; Exit-Codes in beiden Klassifikatoren; Tests Happy, Boundary, Negative, ohne Container, laufender Container, Instanz, Abfrage abgelehnt (Abfrage vor dem Stoppen, `12`, `0` für `--instance` ohne Container, `0` ohne Abfrage ohne Container im nicht-interaktiven Modus) | erfüllt | Ports `RemoveContainer`, `RemoveVolume`, `ContainerVolumes`, `ExistingVolumes`, `ConfirmRemoveDevcontainer` im Diff; `TestExecute_DevcontainerRemove_ExitCodesInBothClassifiers` pinnt `ExitCode` und `mapDevcontainerRemoveErrorToDiagnostic` je Zeile (Mensch und `--json`); Reihenfolge `TestDevcontainerRemove_VolumesRunInTheOrderOfTheSpec`; Docker-Test (Einzellauf, Sensor 5) |
| 2 | Belege der Spezifikation: jede der zwölf Zeilen der Tabelle *Belege* hat einen Test der genannten Ebene; Marker „Zielbild“ entfernt | erfüllt, ein Vorbehalt (V-3) | Tabelle in Abschnitt Traceability; Wortsuche `Zielbild` in `spec/spezifikation.md` (Treffer nur in Historie-Zeilen anderer Vorgänge und in der ersten Historie-Zeile zu [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes)), `docs/maintainer/cli-json-contract.md`, `docs/user/`, den READMEs: keine Marker zu `remove` mehr |
| 3 | Docker-Test am gebauten Binary, Fälle (a) bis (f), Lesestand vor und nach dem Lauf, Beleg als Einzellauf, kein Überspringen | erfüllt | Einzellauf der vier `TestDevcontainerRemove_Docker_*` durch den Verifier: PASS, kein SKIP (Sensor 5); Zuordnung der Fälle unten |
| 4 | Grenzen von „die Volumes dieses Containers“ (CR F-8): Volume, das ein anderer Container hält; anonyme Volumes; nicht gemessene Profile | erfüllt | Gehaltenes Volume: `TestDevcontainerRemove_Docker_OllamaVolumeAndVolumeHeldByAnotherContainer` (Exit-Code `12`, Meldung nennt entferntes und bestehendes Volume, Halter läuft weiter). Anonym: `TestDevcontainerRemove_Docker_StandardProfileAndAnonymousVolume` (`docker rm -v` an gestopptem Container ohne `-f`; ohne `--volumes` bleibt es) und `TestDevcontainerRemove_AnonymousVolumeGoneWithTheContainer`; die Wahl legt [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes) fest. Nicht gemessene Profile: Standardprofil und `ollama` getestet; die übrigen Katalog-Features, externe Features und handgeschriebene Einhängungen stehen als Nicht-Zusage in der Spezifikation (Absatz *Nicht gemessene Profile*) und im Handbuch (Absatz *Grenzen*) |
| 5 | Bedingungen `0.37.2` und `0.37.3`: (1) Image-Satz, (2) Terminal für Fall (e), (3) Ausnahme zu [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) | (1) erfüllt mit Vorbehalt (V-1), (2) erfüllt, (3) erfüllt | (1): kein pinnender Test des Image-Inhalts; Handbuch („im Image nicht zugesagt“, Exit-Code 11 ohne `docker` im `PATH`) und CHANGELOG nach der Fix-Runde als Nicht-Zusage, Spezifikation (Absatz *Reichweite*) mit „Der Beleg ist der Aufruf ohne `docker` im `PATH`, nicht das Image selbst“; Test am Aufruf: `TestRun_DevcontainerRemove_NoDockerBinaryExits11` (`cmd/uboot/main_test.go:723`). (2): `rmRunTerminal` in `internal/e2e/devcontainer_remove_docker_test.go` führt das Binary unter `script` aus (Pseudo-Terminal), tippt die Antwort nach `[y/N]`; fehlendes `script` bricht ab (`t.Fatalf`), überspringt nicht; abgelehnte Abfrage (Fall (e)) und bestätigte Abfrage laufen im Einzellauf. (3): Fall (d) in `TestDevcontainerRemove_Docker_NoContainer` (`--volumes` ohne `--yes`, mit `--no-interactive` vor und nach dem Unterbefehl: Exit-Code `0`, keine Abfrage); `10` bei einem Container: `TestDevcontainerRemove_Docker_ContainersVolumesAndSelection`. `u-boot down --volumes` und `u-boot remove --purge` sind von diesem Slice nicht berührt (kein Diff in ihren Dateien) |
| 6 | Bedingungen des Wortlauts ohne eigenen Fall (g) bis (l) | erfüllt | (g1)/(g2) Auswahl ohne `--instance`, nur Instanzen laufend, gemischt, gestoppt: `ContainersVolumesAndSelection` (Abschnitt „Only containers of instances“); (h) Ollama-Volume: `OllamaVolumeAndVolumeHeldByAnotherContainer` (ohne `--volumes` bleibt es und der neue Container hängt es wieder ein; mit `--volumes --yes` geht es); (i) Nennung vor und nach: `VolumesRunInTheOrderOfTheSpec` am Fake-Port und `ContainersVolumesAndSelection` (Abfrage nennt genau das Volume des Containers, `Volumes removed` stimmt mit dem Lesestand überein); (j) sechs ungültige Namen (leer, `=`, Großbuchstabe, 33 Zeichen, Bindestrich am Anfang und am Ende) Exit-Code `10`, 32 Zeichen `0`: `ContainersVolumesAndSelection`; (k) Exit-Code `11` mit leerem `PATH` im Docker-Test und `TestRun_DevcontainerRemove_NoDockerBinaryExits11`; (l) `remove` ohne `--volumes`, dann `remove --volumes --yes` ohne Container: Exit-Code `0`, Volume bleibt (`ContainersVolumesAndSelection`, Instanz `i1`) |
| 7 | Doku: Handbuch (Abgrenzung, Beispiel, Warnung), CHANGELOG, Hilfetext, Kopien der Befehlsliste (README, `zusagen-und-messstaende.md`) | erfüllt | Handbuch-Absatz „Den Container entfernen“ mit Abgrenzung zu `u-boot remove <dienst>` und `devcontainer down`, drei Beispielen, Warnung „löscht Daten“; Wortsuche nach `LH-`, `SPEC-`, `ADR-` in den hinzugefügten Zeilen von `docs/user/`, `README.md`, `README.de.md`: null Treffer (`AGENTS.md` §5 Regel 6); Hilfetext: `TestExecute_DevcontainerRemove_Help`; `make docs-check` grün (Sensor 2) |
| P1 | `make gates` grün | erfüllt | Sensor 1, auf HEAD `6d8f521a`; der Stempel `.harness/state/gates-passed.head` nennt jetzt HEAD |
| P2 | Review durchgeführt, Report unter `docs/reviews/` | erfüllt | [`2026-11-06-slice-v1-devcontainer-remove`](2026-11-06-slice-v1-devcontainer-remove.md); F-1 bis F-3 und F-5 sind in `6d8f521a` behoben (siehe unten) |
| P3 | Verification-Report unter `docs/reviews/` | mit dieser Datei erfüllt | dieser Report |
| P4 | Closure-Notiz mit Steering-Loop-Lerneintrag; Risiken aus §6 mit Ausgang | offen | §7 des Slice trägt „(Wird bei der Closure gefüllt.)“; beide Risiken in §6 tragen „Ausgang: offen“. Die Häkchen in §2 sind nicht gesetzt. Closure ist Sache von Planner und Implementer nach diesem Report |

### Zuordnung der Fälle (a) bis (f) der DoD zum Docker-Test

| Fall | Belegt durch | Lesestand |
|---|---|---|
| (a) laufender Container entfernt, Volumes bleiben | `ContainersVolumesAndSelection`, `remove --instance i1` (Ausgabe `(stopped first)`, Volume mit Inhalt `marker-i1` danach lesbar, `up` hängt es wieder ein) | Kennung, Lauf-Zustand, Volumes nach Name und Inhalt, `docker images`, vor und nach |
| (b) `--volumes --yes` entfernt genau die Volumes dieses Containers, der zweite Container und seine Volumes bleiben | `ContainersVolumesAndSelection`, `remove --volumes --yes --instance i2` mit `--json`; `removedVolumeNames` gleich dem Lesestand der Einhängungen vor dem Lauf; `state.require` für alle übrigen Container und Volumes | wie (a) |
| (c) `--volumes` ohne `--yes` ohne Terminal: Exit-Code `10`, nichts verändert, laufender Container läuft | `ContainersVolumesAndSelection` (Stdin ist ein leeres Rohr; stdout leer; stderr nennt Container, Volume, `[y/N]`, „neither stopped nor removed“; nur das Volume dieses Containers wird genannt) | `state.require` nach dem Lauf: alle laufen, alle Volumes da |
| (d) ohne Container: Exit-Code `0`, auch `--volumes` ohne `--yes` ohne Terminal | `NoContainer` (vier Aufrufformen, `--json` mit `containerId` leer und `removedVolumeNames: []`, fehlendes `docker` mit `11`, Projektordner als Baum mit Hash vor und nach) | Anzahl der Container des Projektordners, Baum des Projektordners |
| (e) abgelehnte Abfrage unter Terminal: Exit-Code `10`, Container läuft weiter, kein Volume entfernt | `ContainersVolumesAndSelection` (Abschnitt „A declined question under a terminal“, `rmRunTerminal`, Antwort `n`) und für die bestätigte Abfrage `StandardProfileAndAnonymousVolume` (Antwort `y`) | `state.require` nach dem Lauf |
| (f) `remove --instance a` ohne Container der Instanz neben anderen Containern: Exit-Code `0`, andere unberührt | `ContainersVolumesAndSelection` (Abschnitt „An instance without a container“, mit `--volumes --yes`) | `state.require` nach dem Lauf |

### Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| 1. `make gates` (Devcontainer-Container, HEAD `6d8f521a`, Exit 0) | pass | `tool-pruefung` OK (go 1.27.3, golangci-lint 2.14.0), `baseline-verify` OK (54 Dateien), `docs-check` „1181 Datei(en) geprüft, 0 Befund(e)“, `lint` „0 issues.“, `test` alle Pakete `ok` (aus dem Testcache, dessen Schlüssel der Inhalt ist), `coverage-gate` „92.40% meets threshold 90%“, `ablage-pruefung` 26 Fälle, 0 Fehlschläge, `gate-reihenfolge-test` 2 Fälle, 0 Fehlschläge; Stempel nennt `6d8f521a` |
| 2. `make docs-check` (mit diesem Report) | pass | Lauf nach dem Schreiben dieses Reports: „d-check: 1182 Datei(en) geprüft, 0 Befund(e)“ (ein erster Lauf fand zwei ungelinkte Kennungen in diesem Report, sie sind verlinkt) |
| 3. `make verify-depguard` | not run | Begründet: der Diff ändert keine Importregel (`.golangci.yml` nicht im Diff, kein neues Paket; `lint` mit `depguard` ist in Sensor 1 grün) |
| 4. `make test-docker` (volles Paket) | not run | Begründet: Auftrag des Aufrufers; das volle Paket lief in diesem Container schon einmal ins 30-Minuten-Limit ([`BEO-E2E/e2e-paket-am-zeitlimit`](../plan/planning/observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); ersetzt durch Sensor 5 |
| 5. Einzellauf der Docker-Tests: Rezept von `make -n test-docker` (`docker build --target test-docker-tools`, `docker run --rm --network=host` mit Quellbaum, `/var/run/docker.sock` und `/tmp`), `go test -v -tags docker -timeout 30m` mit `-run 'TestDevcontainerRemove_Docker_'` und `./internal/e2e/` statt `./...`, im verschachtelten Daemon dieses Containers (Docker Engine 29.9.3, `linux x86_64`, `overlayfs`) | pass | `ContainersVolumesAndSelection` PASS (34.19 s), `NoContainer` PASS (0.62 s), `StandardProfileAndAnonymousVolume` PASS (20.28 s), `OllamaVolumeAndVolumeHeldByAnotherContainer` PASS (13.76 s); `ok github.com/pt9912/u-boot/internal/e2e 68.849s`, Exit 0; Zählung `SKIP` und `FAIL` im Log: 0; nach dem Lauf bleibt kein Container zurück (`docker ps -a` leer), `git status` sauber |
| 6. CI-Lauf `integration-docker` auf dem Push-Tip | not run | Begründet: die sieben Commits seit `origin/main` sind nicht gepusht; der Einzellauf (Sensor 5) ist der Beleg, den die DoD alternativ zulässt |
| 7. Replay / Golden | nicht anwendbar | kein Generator geändert |

### Prüfung der Fix-Runde (`6d8f521a`) gegen das Review

| Finding | Zustand | Beleg |
|---|---|---|
| F-1 (MEDIUM, Image-Aussage im Handbuch ohne Pin und ohne Nicht-Zusage) | behoben | Handbuch: „sind im Image nicht zugesagt; ohne `docker` im `PATH` enden sie mit Exit-Code 11“; Exit-Code-Absatz: „im Container-Image von `u-boot` ist der Befehl nicht zugesagt“; CHANGELOG: „(the command is not promised inside the released container image)“. Rest in V-1 |
| F-2 (LOW, Feldkommentar `devcontainerUseCase`) | behoben | `internal/adapter/driving/cli/cli.go`: „`u-boot devcontainer up\|down\|shell\|remove`“ |
| F-3 (LOW, Kommentar `inspectorOutputLimit`) | behoben | `internal/adapter/driven/docker/inspector.go`: nennt `docker ps`, `docker context inspect`, `docker inspect`, `docker volume ls` und die Auswertung durch `ExistingVolumes` |
| F-4 (INFO, Exit-Code `14` in nur einem Klassifikator gepinnt) | unverändert, Reviewer: kein Vertragsbruch | Pfad unerreichbar (`TestDevcontainerRemove_ReadsNoProjectFileAndStartsNoHelper`, `TestRemoveCalls_StartNoHelperAndMountNothing`, `TestRun_DevcontainerRemove_NoDockerBinaryExits11`); unten als Hinweis |
| F-5 (INFO, Zustandsmarker der Slice-Zeile) | angepasst, läuft mit diesem Report ab | `docs/plan/planning/open/README.md` nennt jetzt „umgesetzt, Review erfolgt, wartet auf Verifikation“ (V-4) |
| F-6 (INFO, Zuständigkeit Verifier: Stempel, Einzellauf, Fundstellen, Closure) | beantwortet | Stempel nennt jetzt HEAD (Sensor 1); Einzellauf Sensor 5; Fundstellen in Traceability; Closure-Notiz offen (DoD P4) |

### Traceability

| ID | Beleg |
|---|---|
| [`SPEC-058`](../../spec/spezifikation.md#spec-058--container-des-devcontainers-entfernen-u-boot-devcontainer-remove-aufruf-suche-volumes-abfrage-ausgabe-und-exit-codes), Belege Zeile 1 (Ablauf, Reihenfolge, Argumente) | Dienst: `TestDevcontainerRemove_VolumesRunInTheOrderOfTheSpec` (`internal/hexagon/application/devcontainer_remove_test.go:254`), `TestDevcontainerRemove_WithoutVolumesTouchesNoVolume` (:287); Adapter: `TestLauncher_RemoveContainerAndVolumeArguments` (`internal/adapter/driven/docker/devcontainer_remove_test.go:31`), `TestLauncher_RemoveRejectsInvalidArguments` (:58) |
| Zeile 2 (Freigabe in sechs Zeilen, Antworten, keine Terminalprüfung) | Dienst: `TestDevcontainerRemove_ReleaseTable` (:326), `TestDevcontainerRemove_WithoutConfirmerDeclines` (:434); Confirmer: `TestConfirmer_RemoveDevcontainer_Answers` (`internal/adapter/driven/confirm/confirm_test.go:252`), `TestConfirmer_RemoveDevcontainer_Prompt` (:273), `TestConfirmer_RemoveDevcontainer_ReadErrorPropagates` (:298); Verdrahtung: `TestRun_DevcontainerRemove_QuestionAndAnswerAtTheWiring` (`cmd/uboot/main_test.go:824`); Binary: `ContainersVolumesAndSelection` (Rohr und Pseudo-Terminal) |
| Zeile 3 (Auswahl des Containers) | Dienst: `TestDevcontainerRemove_SelectionOfTheContainer` (:450), `TestDevcontainerRemove_ActsOnTheContainerTheSearchReturns` (:529); Binary: `ContainersVolumesAndSelection` |
| Zeile 4 (Menge der Volumes) | Adapter: `TestInspector_ContainerVolumes` (`internal/adapter/driven/docker/devcontainer_remove_test.go:133`), `TestInspector_ExistingVolumes` (:188); Dienst: `TestDevcontainerRemove_SetOfTheVolumes` (:549), `TestDevcontainerRemove_AnonymousVolumeGoneWithTheContainer` (:580); Binary: `StandardProfileAndAnonymousVolume`, `OllamaVolumeAndVolumeHeldByAnotherContainer` |
| Zeile 5 (Fehler beim Entfernen, gehaltenes Volume) | Dienst: `TestDevcontainerRemove_FailuresKeepTheStateOfTheTable` (:604); Binary: `OllamaVolumeAndVolumeHeldByAnotherContainer` |
| Zeile 6 (Text, `--quiet`, `--json`) | CLI: `TestExecute_DevcontainerRemove_Output` (`internal/adapter/driving/cli/devcontainer_remove_test.go:81`) |
| Zeile 7 (Exit-Codes in beiden Klassifikatoren) | CLI: `TestExecute_DevcontainerRemove_ExitCodesInBothClassifiers` (:217), `TestExecute_DevcontainerRemove_UsageErrors` (:169) |
| Zeile 8 (Aufruf, Name der Instanz vor jedem `docker`-Aufruf) | CLI: `TestExecute_DevcontainerRemove_UsageErrors`, `TestExecute_DevcontainerRemove_FlagsReachTheRequest` (:135); Dienst: `TestDevcontainerRemove_InstanceNameIsCheckedBeforeTheEngine` (:754); Binary: `ContainersVolumesAndSelection` |
| Zeile 9 (ohne `docker` im `PATH` Exit-Code `11`) | Binary: `TestRun_DevcontainerRemove_NoDockerBinaryExits11` (`cmd/uboot/main_test.go:723`); Dienst: `TestDevcontainerRemove_EngineNotAskable` (:715); Docker-Test mit leerem `PATH` in `NoContainer` und `ContainersVolumesAndSelection` |
| Zeile 10 (keine Projektdatei, kein Hilfs-Container, kein Socket, kein Exit-Code `14`) | Dienst: `TestDevcontainerRemove_ReadsNoProjectFileAndStartsNoHelper` (:777); Adapter: `TestRemoveCalls_StartNoHelperAndMountNothing` (:96); CLI: `TestExecute_DevcontainerRemove_ReadsNoProjectFile` (:268) |
| Zeile 11 (Zustand nach dem Entfernen) | Binary: `ContainersVolumesAndSelection` (neuer Container, altes Volume mit Inhalt, neues leeres Volume nach `--volumes`, Images erhalten), `OllamaVolumeAndVolumeHeldByAnotherContainer`, `StandardProfileAndAnonymousVolume` |
| Zeile 12 (Inhalt der Meldungen `12` und `10`) | Dienst: `TestDevcontainerRemove_RefusalMessages` (:395), `TestDevcontainerRemove_FailuresKeepTheStateOfTheTable`; Binary: Stderr und JSON-Meldung in `ContainersVolumesAndSelection` und `OllamaVolumeAndVolumeHeldByAnotherContainer` (siehe V-3) |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | beide Klassifikatoren je Pfad der Tabelle: `TestExecute_DevcontainerRemove_ExitCodesInBothClassifiers`; Prozess-Exit-Codes `0`, `2`, `10`, `11`, `12` im Binary |
| [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) | Abfrage nur mit `--volumes` und einem Container; `--yes` zusammen mit `--no-interactive` ergibt `2`; `--json` fragt nie |
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) | Volumes je Instanz bleiben getrennt: Kennungs- und Namenstrennung der Instanzen `i1`, `i2` und der Basis in `ContainersVolumesAndSelection` (`dcRequireDistinct`) |
| Architektur ([`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)) | `make lint` mit `depguard` grün (Sensor 1) |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: `CO-001` (der Diff berührt `docs/plan/carveouts/` nicht)

Nicht ausgeführt:
- `make test-docker` (volles Paket) — begründet in Sensor 4
- CI-Lauf `integration-docker` auf dem Push-Tip — begründet in Sensor 6
- `make verify-depguard` — begründet in Sensor 3
- `make ci` / `make fullbuild` (Security-/Release-Pfad) — der Diff berührt weder Abhängigkeiten noch Dockerfile noch Release; nicht gelaufen

Commit / Artefakt:
- `6d8f521a7636715afda0f33e7c9cbf69f64b325d` (HEAD, Basis `f5567aee`); keine Release-Artefakte

---

## Befunde des Verifiers

Keine HIGH-Verstöße gegen DoD oder Spezifikation; die DoD-Punkte 1 bis 7 und die Pflichten P1 bis P3 sind erfüllt, P4 steht offen (Closure). Die Punkte unten sind Hinweise an den Planner, keine Rückgabe an den Implementer.

| ID | Kategorie | Befund | Zuständig |
|---|---|---|---|
| V-1 | LOW | **Image-Satz: Lastenheft stärker als Handbuch und Spezifikation.** Lastenheft `0.37.3` (Zusage-Menge von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)) sagt: „Im ausgelieferten Container-Image … läuft er nicht, und `u-boot devcontainer remove` läuft dort ebenfalls nicht.“ Handbuch und CHANGELOG sagen jetzt „nicht zugesagt“. Die Spezifikation (Absatz *Reichweite*) behauptet als Tatsache „enthält kein `docker`; der Befehl endet dort mit Exit-Code `11`“ und begrenzt den Beleg auf den Aufruf ohne `docker` im `PATH`. Das Wort „Nicht-Zusage“ steht in der Spezifikation nicht. Gelesen: die Runtime-Stage des `Dockerfile` (`gcr.io/distroless/static-debian12:nonroot`, einzige Datei `COPY … /usr/local/bin/u-boot`) trägt die Tatsachenbehauptung, ein Test liest sie nicht. DoD-Bedingung (1) verlangt Pin oder Nicht-Zusage „in der Spezifikation und im Handbuch“: das Handbuch erfüllt es wörtlich, die Spezifikation der Sache nach. Die Entscheidung, ob die Spezifikation „nicht zugesagt“ ausdrücklich sagt oder das Lastenheft den Satz streicht (Rang 1, Entscheidung des Projektinhabers), liegt beim Planner. | Planner |
| V-2 | INFO | **Lauf-Umgebung des Einzellaufs.** Der Einzellauf lief im verschachtelten Daemon dieses Containers (Linux `x86_64`, Docker Engine 29.9.3, `overlayfs`), nicht an einem Docker außerhalb des Devcontainer-Containers (Sockelform, [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)). Er entspricht der Umgebung der Messung (Engine 29.9.3, `amd64`, `overlayfs`) und der Zusage-Menge (Linux `amd64`, lokale Engine); eine native Engine auf einem Linux-Host ist, wie die Spezifikation sagt, nicht gemessen. Der CI-Lauf auf dem Push-Tip steht nach dem Push aus. | Planner (Closure-Notiz nennt die Umgebung) |
| V-3 | INFO | **Ebene der Zeile 12 der Tabelle *Belege*.** Die Zeile verlangt neben dem Test des Dienstes einen „Test des CLI-Adapters (Text der Meldung)“. Es gibt keinen CLI-Test, der den Wortlaut der Meldungen `10` und `12` pinnt (`TestExecute_DevcontainerRemove_ExitCodesInBothClassifiers` prüft Code, Klasse, Umschlag und dass kein Projektpfad austritt, die Texte kommen aus dem Fake-Fehler). Der Wortlaut ist am Dienst (`TestDevcontainerRemove_RefusalMessages`, `FailuresKeepTheStateOfTheTable`) und am gebauten Binary (Stderr und JSON-Meldung in den Docker-Tests) gepinnt; die Aussage der Zeile ist belegt, die genannte Ebene des Adapters nur durchgereicht. Entweder die Zeile nennt „Docker-Test am gebauten Binary“, oder der Implementer ergänzt einen CLI-Test. | Planner |
| V-4 | INFO | **Zustandsmarker mit Ablaufpunkt.** `docs/plan/planning/open/README.md` führt den Slice mit „Backlog, … liegt in `in-progress/`, umgesetzt, Review erfolgt, wartet auf Verifikation“. „wartet auf Verifikation“ endet mit diesem Report (`AGENTS.md` §3.10); „Backlog“ widerspricht dem Ort `in-progress/` im selben Satz. Beide bei der Closure auf den Ist-Stand setzen. | Planner / Implementer (Closure) |
| V-5 | INFO | **Handbuch nennt ein Volume, das kein `remove`-Test erzeugt.** Die Warnung im Handbuch nennt „mit `nestedRuntime: docker` der Docker-Speicher“ als Volume, das `--volumes` entfernt. Kein Docker-Test von `remove` startet ein Projekt mit `nestedRuntime: docker`; die Aussage stützt die Messung M1 (Einhängung) und die allgemeine Regel „alle Einhängungen vom Typ `volume`“, die `TestInspector_ContainerVolumes` und die Docker-Tests am Sandbox-Profil mit Standardeinstellungen und `ollama` pinnen. Die Spezifikation führt es als gemessene Einhängung, nicht als Test von `remove`. Kein Vertragsbruch; als Aussage über Verhalten hat sie Messung, aber keinen Test (`AGENTS.md` §3.10 Bedingung in Handbuch ohne Pin: die Regel ist gepinnt, ihr Ergebnis für dieses Profil nicht). | Planner |
| V-6 | INFO | **F-4 des Reviews bleibt als Rand.** Der Sentinel `ErrDevcontainerFileSystem` hat in `ExitCode` Exit `14`, in der Abbildung von `remove` den allgemeinen Code; der Pfad ist unerreichbar (Tests in Zeile 10), der Test pinnt nur die Abbildung. Kein Verifier-Befund gegen die DoD. | — |

## Verdikt

**DoD:** erfüllt für die Punkte 1 bis 7 und die Pflichten P1 bis P3; offen ist allein P4 (Closure-Notiz mit Steering-Loop-Eintrag, Ausgänge der Risiken in §6, Häkchen in §2, Verschiebung nach `done/` als reiner `git mv` nach `AGENTS.md` §3.3 Punkt 2).

**Einzellauf der Docker-Tests:** vier von vier PASS, kein Überspringen, auf `6d8f521a`, mit Lesestand vor und nach dem Lauf je Fall.

**Plan-vs-Code:** jede der 37 Dateien des Diffs hat einen Plan-Eintrag oder ist ein Bericht (Review-Report, dieser Report); die zwei Zeilen „gelesen, unverändert“ (`docs/maintainer/devcontainer-arbeiten.md`, `docs/plan/planning/in-progress/roadmap.md`) haben keinen Diff; die Tabelle *Belege* hat zwölf Zeilen wie die Spezifikation.

**Entscheidungs-Konformität:** [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Volumes je Instanz) und [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) (Nachweis im Devcontainer-Container) gehalten; kein `//nolint` im Diff, kein Carveout, keine Lockerung eines Gates, keine Importregel geändert.

**Übergabe an den Planner:** V-1, V-3 und V-5 entscheiden über Wortlaut und Ebene (kein Code); V-4 beim Schließen; V-2 in die Closure-Notiz. Kein Befund erfordert eine neue Runde des Implementers vor der Closure.
