# Verification-Report: [`slice-v1-devcontainer-instanz-label`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-label.md) — 2026-11-01

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) und den ausgeführten Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Außer dieser Datei ist nichts geändert, nichts committet. Die Datei entstand per `cp` aus `.harness/baseline/v6.18.1/templates/docs/reviews/review-report.template.md` (`diff -q` direkt nach dem `cp`: identisch; es gibt keine eigene Verifikations-Vorlage, Vorbild ist der Report [`2026-10-31-slice-v1-devcontainer-instanz-messung-verify`](2026-10-31-slice-v1-devcontainer-instanz-messung-verify.md)) und wurde mit dem Write-Werkzeug gefüllt.

**Gegenstand:** `git diff 71ba20a4^..HEAD` (HEAD `c361601a`; acht Commits `71ba20a4` bis `c361601a`), ohne `.claude/worktrees`. Eingang gelesen, nicht als Beweis genommen: Review-Report [`2026-11-01-slice-v1-devcontainer-instanz-label`](2026-11-01-slice-v1-devcontainer-instanz-label.md), Beleg-Tabelle des Slice.

**Verdikt:** Code, Tests der Ebenen Unit/Adapter/CLI, Spec, Doku und `make gates`: **erfüllt**. Die [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) ist konform umgesetzt. Der neue Docker-Test des Slice (vier Instanzen) **besteht** im eigenen `make test-docker`-Lauf. Der Lauf ist aber **insgesamt rot**: drei Bestands-Docker-Tests schlugen fehl (Abschnitt 4). Keiner berührt Code dieses Slice; die Ursache ist nicht bewiesen (Verdacht: parallele Docker-Nutzung der zwei Worktree-Agenten auf demselben Daemon). Ein zweiter Lauf ist auf Weisung des Projektinhabers ausgeschlossen. **Bereit für die Closure: bedingt**; die drei roten Tests sind vom Planner/Implementer gegen den grünen Implementer-Lauf (1540 s) und den CI-Lauf auf `main` einzuordnen, nicht von mir als „grün“ zu belegen. Keine HIGH, eine MEDIUM (V-1), zwei INFO.

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-instanz-label`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-label.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- Artefakte: 34 Dateien des Diffs (+1795/−121), siehe Abschnitt 3

### 1. DoD-Abgleich (gegen Code und Tests, gelesene Fundstellen)

- [x] **Annahmen des Mess-Slice geprüft** (Docker-Test, N ≥ 3). Fundstelle `internal/e2e/devcontainer_instance_docker_test.go`: `dcInstanceCount = 4` (Zeile 36); `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` (Zeile 212) fährt `up --instance` für vier Instanzen neben dem Container ohne Instanz, prüft Container-IDs und die Volumes der drei Familien paarweise verschieden, Inhaltsprobe über alle geordneten Paare, Labels, zurückgelesene Grenze (1 CPU, 3 GiB), `down`/`shell`/`up --recreate`, „nur Instanzen“ per Binary. Host-Limit: `dcRequireInstanceHost` (Zeile 81) bricht mit `FINDING` ab (Memory/Swap-Limit, CPUs ≥ 4, `MemAvailable`, `SwapFree`, Platz), senkt die Stufe nicht. **Selbst gefahren: PASS** (`--- PASS: …FourInstancesAreIsolated (41.50s)`, `…InvalidNamesEndWithExit10 (1.97s)`, `…StandardProfileSharesTheProjectFolder (9.88s)`, Log Zeilen 4349 bis 4353). Eine Anmerkung des Tests: Platz des Docker-Roots „nicht lesbar … nicht geprüft“ (Test läuft im Container, `/var/lib/docker` fehlt dort) — Teil der Prüfung entfällt, vom Test selbst ausgewiesen (Log Zeile 4342).
- [x] **Der gewählte Weg A.**
  - Flag an `up`, `down`, `shell`: `internal/adapter/driving/cli/devcontainer.go`, `registerDevcontainerInstanceFlag` in allen drei Unterbefehlen; Hilfetexte mit Beispiel; `InstanceGiven` über `cmd.Flags().Changed("instance")`. Test: `TestExecute_Devcontainer_InstanceFlagReachesTheRequest`, `…HelpNamesInstance` (PASS im Log).
  - Namensregel `[a-z0-9-]`, Rand, 1 bis 32, leer und `=` ungültig: `driven.ValidateDevcontainerInstance` (`internal/hexagon/port/driven/devcontainer_launcher.go`), Test `TestValidateDevcontainerInstance` (PASS). Leeres `--instance`: `emptyInstanceError`, Test `TestExecute_Devcontainer_EmptyInstanceIsInvalid` (6 Teilfälle, PASS). Exit-Code 10: `isDevcontainerValidationError` nimmt `isDevcontainerInstanceError` (beide Sentinels), `mapDevcontainerDownErrorToDiagnostic` hat den neuen Zweig; Tests `TestExecute_Devcontainer_SentinelsPerSubcommand`, `…ShellSentinelsAndExitCodes`.
  - `FindContainer` findet die Instanz: `internal/adapter/driven/docker/inspector.go` (zwei `--filter label=…`; ohne Instanz `{{.Label "u-boot.instance"}}` je Zeile, `ErrOnlyInstanceContainers` bei nur Instanzen). Tests `TestInspector_FindContainer_InstanceArgv`, `…_Instances`, `…_RejectsBadInstance` (PASS); gegen die echte Engine trägt es der Docker-Test (PASS, damit ist die in [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidung 4 offene Frage zum Format `{{.Label}}` belegt).
  - Ohne Angabe unverändert (Bestandstest): `TestRunArgs/defaults` (kein `--id-label`), `TestDevcontainerInstance_ExistingContainerAndNewerInstance`, `TestDevcontainerUp_Docker_S1_S9_StandardProfile` (Teil des Docker-Laufs, im grünen Bereich: kein FAIL dort).
  - „Nur Instanzen“: Exit 10 mit Hinweis, kein Launcher, kein `docker stop`, je `up` (auch `--recreate`), `down`, `shell`: `containerBeforeStart` und `Down` in `internal/hexagon/application/devcontainer_start.go`; Test `TestDevcontainerInstance_OnlyInstancesEndsBeforeTheLauncher` (PASS), Docker-Test (Binary).
  - Docker-Test mit N Instanzen: s. o.
- [x] **Spezifikation, Handbuch, `CHANGELOG.md`** mit Beleg-Tabelle. `spec/spezifikation.md` [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Absatz *Aufruf* (Zeile 833) und *Instanzen* (Zeilen 837 bis 844) tragen Flag, Namensregel, Identität, Suche, „nur Instanzen“, Exit-Code-Tabelle (Zeile 900). Die Beleg-Tabelle in Slice §3 (Zeilen 133 bis 156) ordnet jede Bedingung einem Test zu; die in der Fix-Runde `81a31dff` ergänzten Pins (`requireNoInstanceAsked` in den Tests von `doctor`, `generate`, `init`) liegen in `devcontainer_fallback_test.go` und `devcontainer_instance_test.go` (Diff gelesen). Handbuch/READMEs: `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md` im Diff; Stichprobe `grep -nE "LH-|ADR-|SPEC-|slice-" ` auf diese drei Dateien im Diff siehe Abschnitt 3 (keine Kennung neu eingeführt).
- [x] **`make gates` grün** — selbst gefahren (Abschnitt 2).
- [ ] **Review** — Report liegt vor (0 HIGH/MEDIUM, 4 LOW, 2 INFO), Fix-Runde `81a31dff` behoben (Abschnitt 5). Das Häkchen setzt der Planner.
- [ ] Doku-Update (Häkchen), Closure-Notiz, Register, Risiko-Ausgänge, Paarungen — offen (Planner), nicht Gegenstand dieser Prüfung. Slice §7 ist Platzhalter wie vorgesehen; im Slice sind alle DoD-Häkchen `[ ]`.

### 2. Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` (Docker-Modus am Host) | pass | Exit 0 auf `c361601a`; `coverage-gate: OK — coverage 92.20% meets threshold 90%`; `d-check: 951 Datei(en) geprüft, 0 Befund(e)`; `ablage-pruefung: OK`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `baseline-verify: v6.18.1 OK — 54 Dateien` |
| `make test-docker` | **fail (3 Bestandstests rot, neue Tests grün)** | Exit 2; `ok` für alle Pakete außer `internal/e2e` (`FAIL internal/e2e 1224.538s`); Details Abschnitt 4 |
| `make verify-depguard` | not run | keine Importregel berührt (Diff ohne `.golangci.yml`, keine neue Schichtgrenze; Review-Report „Architekturgrenzen“) |

**Testdauer.** `make test-docker` Gesamtlauf rund 26 Minuten; Paket `internal/e2e` 1224,5 s, Paket `internal/adapter/driven/docker` 176,1 s. **Vollständige `-v`-Ausgabe: `/tmp/claude-1000/-Development-u-boot/c9a21911-8809-487d-8cda-2b4e733058ba/scratchpad/test-docker.log`** (Kopie von `…/scratchpad/docker.log`, dort mit der Zeile `EXIT docker 2` am Ende; `make gates`: `…/scratchpad/gates.log`). Es lief genau ein Lauf; er hatte von Beginn an Log.

### 3. Plan-vs-Code-Diff (`git diff --stat 71ba20a4^..HEAD -- . ':!.claude/worktrees'`: 34 Dateien)

Jede Datei hat eine Plan-Zeile (§3 Tabelle oder Abgleich):
- Plan-Zeilen mit Diff (geplant): `cli/devcontainer.go`, `port/driving/devcontainer.go`, `port/driven/devcontainer_launcher.go`, `…devcontainer_inspector.go`, `docker/devcontainer_launcher.go`, `docker/inspector.go`, `application/devcontainer_start.go`, `…devcontainer_fallback.go`, `noop.go`, `spec/spezifikation.md`, Handbuch, `README.md`, `README.de.md`, `CHANGELOG.md`, `cli-json-contract.md`, `zusagen-und-messstaende.md`, Slice-Datei, Tests und Doubles (`fakeInspector`, `inspector_test.go`, `devcontainer_launcher_docker_test.go`, `devcontainer_start_docker_test.go`, `sandbox_status_docker_test.go`).
- Neu, im Plan als „Neu“ geführt: `application/devcontainer_instance_test.go`, `cli/devcontainer_instance_test.go`, `e2e/devcontainer_instance_docker_test.go`; dazu die zwei Review-Reports.
- Abweichend geführt (Abgleich im Slice benennt sie): `devcontainer_shell_test.go` (application und cli), `cli/devcontainer_test.go`, `docker/devcontainer_launcher_test.go`, `port/driven/devcontainer_launcher_test.go`, `observation.md` (Umbenennung `runningContainerID` → `containerBeforeStart`).
- Nicht in der Plan-Tabelle, aber im Diff: **`.d-check.yml`** (Commit `d888dd8c`, „docs-check wertet die Arbeitsbäume unter `.claude/worktrees` nicht“) — siehe V-2.
- Plan-Zeile ohne Diff: `cli.go` (Plan-Zeile entfällt; in der Fix-Runde berichtigt, Slice Abgleich „Plan-Zeile `cli.go` entfällt“). `spec/lastenheft.md`, `harness/replay.md`, `export_test.go`, `docs/plan/adr/0038-…` und `adr/README.md`: im Diff unberührt (Lastenheft/ADR liegen aus früheren Commits vor, [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) `Accepted` und seither unverändert, `git log` der Datei: `5953d4f7`).

Hard Rules: kein `//nolint` im Go-Diff (`grep` leer); keine Gate-, Schwellen- oder Linter-Änderung außer der `.d-check.yml`-Ausnahme (V-2); keine Accepted-ADR verändert; `spec/spezifikation.md` ohne Slice-/Welle-Verweis (Review-Report, Negativbefund bestätigt durch Diffstat: nur Abschnitt [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) und [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)); Arbeitsbaum zu Beginn sauber, `make gates` verändert keine getrackte Datei bis auf diese Datei.

### 4. `make test-docker`: drei rote Bestandstests

| Test | Dauer | Meldung (Log) | Berührt der Slice? |
|---|---|---|---|
| `TestE2E_LHFAUP005_LogsTail` | 1,83 s | `up: up service: ComposeUp … docker compose up failed (exit status 1): compose runtime error` (Zeile 4546 f.) | nein (Compose, `logs`) |
| `TestE2E_LHFAUP005_LogsFollow` | 1,20 s | gleiche Meldung | nein |
| `TestE2E_SandboxDevcontainer_TwoInstancesIsolated` | 41,56 s | `start instance idaaaa: docker: Error response from daemon: Conflict. The container name "/uboot-e2e-inst-idaaaa" is already in use by container "52947d09…"` (`sandbox_devcontainer_docker_test.go:589`) | nein (Bestandstest, ruft den Befehl nicht; Datei im Diff unberührt; er ist im Plan §3 als „unverändert“ genannt) |

Einordnung (Beobachtung, nicht bewiesen): Der Namenskonflikt zeigt einen Container desselben festen Namens aus einem anderen, gleichzeitigen oder abgebrochenen Lauf; nach meinem Lauf lag kein solcher Container mehr vor (`docker ps -a`: `exciting_boyd`, `nostalgic_nash`, `u-boot-devcontainer-up-7c21bec5`, das letzte „Up 5 seconds“ — also läuft auf diesem Daemon Fremdlast der Worktree-Agenten). Die zwei Compose-Fehler (`compose runtime error`) tragen im Log keine Ursache; Verdacht gleicher Art (geteilte Ports, Namen oder Compose-Projekt auf dem geteilten Daemon). **Nicht belegt**, weil nicht wiederholt (Weisung des Projektinhabers: `make test-docker` kein zweites Mal lokal). Gegenhalt: Der Implementer-Lauf war nach Angabe des Koordinators grün (1540 s); den Wert habe ich nicht nachgeprüft. Der Docker-Nachweis des Slice (die drei `TestDevcontainerInstance_Docker_*`) ist in meinem Lauf grün; die Bestandstests, die ein Wiederholen bräuchten, bleiben offen.

Aufträge daraus: (a) CI-Lauf auf `main` abwarten und die drei Namen dort prüfen; (b) falls die Drei dort rot sind, ist es ein Befund gegen den Slice nur, wenn eine Ursache im Diff liegt — bisher keine gefunden (kein Compose- und kein `TwoInstances`-Code im Diff).

### 5. Fix-Runde `81a31dff` (F-1 bis F-4 des Reviews)

| Finding | Ergebnis | Beleg |
|---|---|---|
| F-1 Ausweisung bei „nur Instanzen“ | behoben (Ausgang akzeptiert) | Slice §6 Risiko *Ausweisung des Sockets*: Ausgang „akzeptiert, Reihenfolge bleibt (Review F-1)“ mit Gründen ([`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), `TestDevcontainerUp_AlreadyRunning`, [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)); Test-Kommentar ohne ADR-Bezug; [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeile 842 nennt die Reihenfolge |
| F-2 Zustandsmarker | behoben | Slice-Diff (`81a31dff`, 11 Zeilen) |
| F-3 `cli.go`-Zeile | behoben | Slice Abgleich, Zeile „Plan-Zeile `cli.go` entfällt“ |
| F-4 Pin an `doctor`/`generate`/`init` | behoben | `requireNoInstanceAsked` in `devcontainer_fallback_test.go`; Tests `TestDoctor_SandboxFallback_OnlyConfiguredCapabilities`, `TestGenerateDevcontainer_ExistingContainerFallbackWarning`, `TestInit_ExistingContainerFallbackWarning` im Gate-Lauf grün (Mutation selbst nicht wiederholt) |

### 6. Entscheidungs-Konformität ADR-0038

| Entscheidung | Befund |
|---|---|
| 1 Flag an drei Unterbefehlen | konform (CLI-Diff) |
| 2 zwei `--id-label` mit Instanz; `FindContainer` mit zwei Filtern | konform: `runArgs` hängt `devcontainerFolderLabel=<Ordner>` und `devcontainerInstanceLabel=<Name>` nur unter `opts.Instance != ""` an; Inspector wie beschrieben; Labelname `u-boot.instance` |
| 3 kein Label für den Container ohne Instanz | konform: `TestRunArgs/defaults`; Docker-Test liest Labels zurück |
| 4 ohne Flag nur der Container ohne Instanz, Label je Zeile aus `docker ps` | konform: `{{.Label "u-boot.instance"}}`, `splitPsLine`; echte Engine durch Docker-Test belegt |
| 5 Prüfung vor dem CLI-Aufruf | konform; Abweichung 2 des Slice (nach Ausweisung, vor Launcher) ist keine Verletzung der ADR, die Reihenfolge ist dort offen; als Risiko ausgewiesen |
| 6 Namensregel einmal im Port, Aufrufer wenden sie an | konform: `ValidateDevcontainerInstance`; Application, Inspector, Launcher rufen sie |
| 7 Schnitt durch die Schichten, zwei Sentinels in `driving` | konform: `ErrDevcontainerInstanceInvalid`, `ErrDevcontainerOnlyInstances`; Port-Sentinel `ErrOnlyInstanceContainers` (Form laut ADR dem Slice überlassen) |
| 8 Dual-Classifier | konform: zentrale Klassifikation und Diagnostic-Abbildung von `up` und `down`; `shell` über die zentrale Klassifikation |
| 9 Urteil nach dem Start bei der Instanz | konform: `startOutcome` bekommt die Anfrage mit Instanz in den `fallbackReader` |

Eine Notiz: Code und Slice nennen `containerBeforeStart`, [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) weiter `runningContainerID` (Accepted, nicht umschreibbar; im Slice-Abgleich genannt).

### 7. Befunde des Verifiers

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-1 | MEDIUM | `make test-docker` ist in diesem Lauf rot: drei Bestandstests (Abschnitt 4). Ursache offen, Verdacht Fremdlast auf dem geteilten Daemon; keine Wiederholung zulässig. Die Closure („`make gates` grün“ ist DoD, `test-docker` ist es nicht) wird dadurch nicht blockiert, aber der Docker-Nachweis des Verifiers ist unvollständig | `test-docker.log` Zeilen 4546 bis 4551, 4608 |
| V-2 | LOW | Der Diff berührt `.d-check.yml` (Ausnahme für `.claude/worktrees`, Commit `d888dd8c`, Begründung MR-005), das der Plan-Abgleich des Slice nicht nennt (§3.10 verlangt jede berührte Datei). Es ist keine Gate-Lockerung für Projektdateien (gitignorierte Kopien des Repos), und `docs-check` ist grün (951 Dateien, 0 Befunde); die Plan-Tabelle ist nachzuziehen | `git diff --stat`: `.d-check.yml | 4 +-`; Slice §3 Abgleich ohne diese Datei |
| V-3 | INFO | Der Docker-Test meldet den freien Platz des Docker-Roots als „nicht lesbar, nicht geprüft“ (Test läuft im Container, ohne `/var/lib/docker`); die DoD verlangt das Host-Limit „Speicher, CPU, Platz“, Platz bleibt hier ungeprüft. Speicher und CPU wirken (Test sonst `FINDING`) | Log Zeile 4342 |
| V-4 | INFO | Der Slice führt alle DoD-Häkchen `[ ]` und §7 als Platzhalter; erwartet vor der Closure, hier nur festgehalten | Slice §2, §7 |

Summe: HIGH 0 · MEDIUM 1 (V-1) · LOW 1 (V-2) · INFO 2.

### 8. Traceability und Carveouts

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Lastenheft Zeilen 501 bis 509 (Flag, Exit 10, bis zu vier Instanzen); Tests `TestDevcontainerInstance_*` (Unit), `TestDevcontainerInstance_Docker_*` (Docker, grün) |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | Volumes je Instanz: `…FourInstancesAreIsolated` (grün) |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | `TestExecute_Devcontainer_SentinelsPerSubcommand`, `…ShellSentinelsAndExitCodes` |
| [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) | Abschnitt 6, Review-Report |

Carveouts: Neu: none · Gelöst: none · Unverändert: alle.

Nicht ausgeführt:
- `make verify-depguard` — keine Importregel berührt.
- `make ci` (`govulncheck`, `image-scan`) — kein Security- oder Release-Pfad berührt.
- Zweiter `make test-docker`-Lauf, auch gezielt für die drei roten Tests — Weisung des Projektinhabers (kein zweiter lokaler Lauf); der Docker-Nachweis ist damit „lokal gelaufen mit drei roten Bestandstests“, CI-Lauf auf `main` folgt.

Commit / Artefakt: `c361601a` (HEAD); Logs im Scratchpad: `…/scratchpad/test-docker.log`, `…/scratchpad/gates.log`.
