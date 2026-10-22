# Verification-Report: [`slice-v1-cli-devcontainer-shell`](../plan/planning/done/slice-v1-cli-devcontainer-shell.md) — 2026-10-22

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), [`ADR-0032`](../plan/adr/0032-katalog-features-socket-und-lockerungen.md), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) und [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)?

**Gegenstand:** `git diff --stat 41e0b262..HEAD` (`HEAD` = `302b3a7b`, 15 Commits ab `origin/main`, 41 Dateien, 2580 Einfügungen, 213 Löschungen; Arbeitsbaum zu Beginn sauber, nichts gepusht). Eingang: Slice §1 bis §8 mit *Entscheidungen beim Start*, *Umsetzung* (Beleg-Tabelle, Mutationen, Sensoren, Abgleich), *Fix-Runde*; Review-Report [`2026-10-22-slice-v1-cli-devcontainer-shell`](2026-10-22-slice-v1-cli-devcontainer-shell.md) und [Nachlauf](2026-10-22-slice-v1-cli-devcontainer-shell-nachlauf.md). Die Reports und die Sensor-Zeilen des Implementers habe ich nicht als Beweis genommen; die Sensoren und Docker-Läufe unten sind eigene Läufe.

**Gesamturteil: GO für die Closure**, mit drei Berichtigungen an Plan und Kontrakt (V-1 bis V-3, alle LOW, der Planner entscheidet) und einem offenen Nachweis (CI auf dem Push-Tip, nicht prüfbar). Keine Verletzung einer Entscheidung, einer Hard Rule oder einer DoD-Zeile des Baus gefunden; kein Nutzertext trägt eine Aussage ohne Test oder eigene Messung.

## Verification Evidence

Scope:

- Slice: [`slice-v1-cli-devcontainer-shell`](../plan/planning/done/slice-v1-cli-devcontainer-shell.md), vierter von sechs der Welle [`welle-devcontainer-shell-und-zweites-up`](../plan/planning/welle-devcontainer-shell-und-zweites-up.md).
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md).
- Artefakte: `Shell` in `internal/hexagon/application/devcontainer_start.go` (gemeinsame Startphase `start`, `sessionEnvNames`, `withoutComments`), Port `internal/hexagon/port/driving/devcontainer.go`, CLI `internal/adapter/driving/cli/devcontainer.go` und `terminal.go`, Docker-Tests `internal/e2e/devcontainer_shell_docker_test.go`, Handbuch, READMEs, `CHANGELOG.md`, `docs/maintainer/cli-json-contract.md`.

### DoD-Abgleich (§2 des Slice, Stand `302b3a7b`)

- [x] **Befehl** (DoD 1) — Beleg: Unit- und CLI-Tests (Beleg-Tabelle des Slice; alle benannten Tests gelesen, 21 von 22 Namen existieren, V-2), eigene Läufe `make gates` Exit 0 und 14 Mutationen rot (unten). Am echten Binary gemessen: Start wie `up`, Benutzer `vscode`, Workspace `/workspaces/vbproj`, `GIT_TOKEN` und Agenten-Variable gesetzt, ungesetzte Variable fehlt, Exit 0 nach `exit 3` der Shell, Ausweisung nur auf der Fehlerausgabe, kein Hilfs-Container nach dem Lauf. `--json` ergibt Exit 2 mit Umschlag `subcommand: shell` auf der Standardausgabe, `--dry-run`/`--diff` und ein Argument Exit 2, ohne `docker` im `PATH` Exit 11 mit Ausweisung, ohne `devcontainer.json` Exit 10, `--no-interactive` ohne Wirkung (Exit 0). Beide Klassifikatoren: `TestExitCode_DevcontainerShellUsageSentinels`, `TestExecute_DevcontainerShell_SentinelsAndExitCodes` (prüft `ExitCode` und den Weg über den Befehl je Fall).
- [x] **Docker-Tests am Befehl mit Wächter** (DoD 2) — Beleg: eigener Lauf, `TestDevcontainerShell_Docker_Session` (drei Teile, darunter Terminal unter `script` und Gegenprobe Pipe) und `TestDevcontainerShell_Docker_FailedStatusOpensNoSession` sichtbar `PASS`, kein `SKIP` (Log: 0 Treffer). Wächter wirksam: Mutation `--rm` entfernt ist rot. Der CI-Teil ("sichtbar PASS im CI-Lauf") ist offen, siehe Sensoren.
- [x] **Doku** (DoD 3) — Beleg: Handbuch 4.9, 4.11 und alle Aufzählungen, beide READMEs, `CHANGELOG.md` `[Unreleased]` *Added*, Kontrakt §6.11 und §7 geändert; in den hinzugefügten Zeilen von `docs/user/` und beiden READMEs keine Kennung (Suche nach `LH-`, `ADR-`, `SPEC-`, `ARC-`, `slice-`, `welle-`, `spec/`, `docs/plan`: 0 Treffer). Jede Nutzeraussage gegen Test oder eigene Messung gelesen (unten); Fehler nur in der Tabelle des Kontrakts (V-1).
- [ ] `make gates` grün (erfüllt, eigener Lauf Exit 0); `make test-docker` als CI-Lauf auf dem Push-Tip — **offen**: der Stand ist lokal und ungepusht, `gh run list --branch main` zeigt ihn nicht. Der Tip `41e0b262` hatte `ci` grün (Angabe des Auftrags, nicht von mir geprüft). `depguard` unverändert: `.golangci.yml`, `Makefile`, `repo.mk`, `Dockerfile`, `.github/` im Diff leer; `make verify-depguard` trotzdem gelaufen (Exit 0, alle 9 Regeln).
- [ ] Review durchgeführt — erfüllt (zwei Reports vorhanden, GO); der Haken ist Sache der Closure.
- [ ] Doku-Update Roadmap, Closure-Notiz, Register, Risiken, Paarungen — **bewusst offen**: §7 ist der Platzhalter, die Roadmap-Zeile steht auf "In Arbeit", `open/README.md` zeigt auf den in-progress-Pfad. Das ist die Closure des Planners.

Sensoren (eigene Läufe):

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` | pass, Exit 0 | `coverage-gate` 92,00 % gegen 90 %, `docs-check` 0 Befunde, `ablage-pruefung` 26 Fälle ohne Fehlschlag |
| `make docs-check` (nach dem Report) | pass, Exit 0 | |
| `make doc-trace` | pass, Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-immutable RANGE="v0.13.1..HEAD"` | pass, Exit 0 | 76 Commits, 0 Befunde |
| `make doc-commits RANGE="origin/main..HEAD"` | pass, Exit 0 | 15 Commits, 0 Befunde |
| `make verify-depguard` | pass, Exit 0 | alle 9 Regeln feuern |
| Docker-Tests (Image `u-boot:test-docker-tools`, Build-Befehl des `test-docker`-Targets, `-run 'TestDevcontainerShell_Docker_\|TestShell_Docker_\|TestDevcontainerUp_Docker_'` über `./internal/e2e/... ./internal/adapter/...`) | pass, Exit 0 | `internal/e2e` 227,7 s, `internal/adapter/driven/docker` 35,9 s; alle Tests und Teile `PASS`, 0 `SKIP`; Linux `amd64`, Engine und Client 29.9.2 |
| `make test-docker` als Ganzes | nicht gelaufen | Zeitlimit des Pakets ([`BEO-E2E/e2e-paket-am-zeitlimit`](../plan/planning/observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); der gezielte Lauf deckt die Tests des Slice und die Nachbarn ab, das Ganze ist CI-Gegenstand |
| CI `integration-docker`/`ci` auf dem Push-Tip | nicht prüfbar | lokaler, ungepushter Stand |
| DinD `docker:28.0.4-dind` | nicht gelaufen, begründet | siehe Engine 28.0.4 unten |

**Engine 28.0.4 (CI-Runner).** Ich habe `internal/e2e/devcontainer_shell_docker_test.go` auf Prüfungen von Engine-Text gelesen: Die Tests suchen nur nach Texten von `u-boot` selbst (Ausweisung, `for the duration of the start)`, `devcontainer shell:`, "the container is running and this call did not run the start script", `u-boot devcontainer up --recreate`) und nach Werten aus `docker inspect`/`docker ps` über Formatvorlagen; kein Meldungstext der Engine wird gepinnt. Deshalb habe ich DinD nicht gefahren; die Erkennung der Client-Fehler (Engine-Text) liegt im Fundament und ist dort geprüft. Restrisiko: ein anderes Verhalten von Engine 28.0.4 jenseits des Textes, nicht gemessen.

**Aufräumen.** Vorher/nachher geprüft (Container, Volumes, Images): Der Docker-Lauf ließ drei anonyme Volumes zurück, die Mutation ohne `--rm` vier Hilfs-Container (`u-boot-devcontainer-up-*`), der Binary-Lauf einen Container, ein Image `vsc-vbproj-*` (zwei Tags) und ein Volume `vbproj-workspace-*`; alles entfernt, Container- und Volume-Liste gleich dem Stand vorher. Images: Die Gates bauten die Stufen `u-boot:lint`, `:test`, `:coverage` und `:test-docker-tools` neu (andere IDs, alte dangling); nicht entfernt, Build-Cache. Wegwerf-Verzeichnisse unter `/tmp` entfernt.

## Mutationen (selbst gefahren, 14, alle rot)

Unit: `go test` über `./internal/hexagon/application/ ./internal/adapter/driving/cli/ ./cmd/uboot/` im Image `u-boot:test-docker-tools`; jede Änderung danach per `git checkout -- <Datei>` zurückgesetzt (Arbeitsbaum vorher sauber), `git diff -- internal cmd` leer.

| Nr | Änderung | Ergebnis | Rot in |
| --- | --- | --- | --- |
| 1 | Recreate in der Startphase von `Shell` (`recreate: true`) | rot | `TestDevcontainerShell_StartIsUpWithoutRecreate` |
| 2 | Namensgleichheit in `sessionEnvNames` entfernt | rot | `TestSessionEnvNames/another_name_in_the_reference` |
| 3 | `--json` bei `shell` nicht abgewiesen | rot | `TestExecute_DevcontainerShell_RejectsFlagsAndArguments` (`json`, `json_before_the_command`, `the_rejection_wins…`), `TestRun_DevcontainerShell_NoDockerBinaryExits11` |
| 4 | `--dry-run`/`--diff` nicht abgewiesen | rot | `…RejectsFlagsAndArguments/dry-run`, `/diff` |
| 5 | Warnungen auf die Standardausgabe | rot | `TestExecute_DevcontainerShell_HumanOutput`, `…_Quiet` (zwei Teile), `TestRun_DevcontainerShell_NoDockerBinaryExits11` |
| 6 | Terminal immer `true` | rot | `…HumanOutput`, `TestExecute_DevcontainerShell_Terminal` (zwei Teile) |
| 7 | Terminal immer `false` | rot | `TestExecute_DevcontainerShell_Terminal` (zwei Teile) |
| 8 | Fehler der Startphase ignoriert (`failed` ausgenommen) | rot | `TestDevcontainerShell_FailedStatusOpensNoSession` (zwei Teile); die rohe Variante (`_ = err`) ergibt Panik in `…UnwiredPortsAreErrors` und `TestRun_DevcontainerShell_NoDockerBinaryExits11` |
| 9 | Hilfe-Satz zur `u-boot.yaml`-Bedingung entfernt | rot | `TestExecute_DevcontainerShell_Help` |
| 10 | `--quiet` hält den Fortschritt nicht zurück | rot | `TestExecute_DevcontainerShell_Quiet` |
| 11 | Ausweisung "call" statt "start" | rot | `TestDevcontainerShell_SocketDisclosureNamesTheStart`, `TestRun_DevcontainerShell_NoDockerBinaryExits11` |
| 12 | Fallback-Warnungen nicht an `OnWarning` | rot | `TestDevcontainerShell_FallbackWarningsComeBeforeTheSession/fallback` |
| 13 | `visibleWarnings(…, false)` in `shell` | rot | `TestExecute_DevcontainerShell_Quiet/a_notice_that_is_neither_relaxation_nor_fallback_is_dropped` |
| 14 | `--rm` aus dem Hilfs-Container entfernt (Docker-Lauf) | rot | `TestDevcontainerShell_Docker_Session/session_without_a_terminal` |

Anmerkung zu 3: Mein erster Versuch ersetzte `case flags.JSON:` im ganzen File und färbte auch `TestExecute_DevcontainerUp_JSON`; die Zeile in der Tabelle gilt für die auf `shell` begrenzte Fassung.

## Echtes Binary (eigener Lauf)

Binary per `go build` im Image `u-boot:test-docker-tools`, Wegwerf-Projekt `u-boot init vbproj --no-git --devcontainer --sandbox` plus ein zusätzlicher `remoteEnv`-Eintrag; Lauf ohne Terminal mit Skript auf der Standardeingabe (`id -un`, `pwd`, `$GIT_TOKEN`, `$AGENT_VB`, ungesetzte Variable, `[ -t 0 ]`, `exit 3`):

- Standardausgabe der Sitzung exakt: `vscode`, `/workspaces/vbproj`, `TOK=placeholder-token-123`, `AG=agentval`, `UNSET=absent`, `TTY=no`; Exit 0.
- Standardfehlerausgabe: die Warnung "this command mounts the host container-engine socket …" mit der Anforderungskennung von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und der Fortschritt der Dev-Container-CLI; unter `--quiet` bleibt die Warnung, der Fortschritt entfällt, die Sitzungsausgabe (`hi`) ist unverändert.
- Exit 11 bei `failed`: das Binary im Docker-Test `TestDevcontainerShell_Docker_FailedStatusOpensNoSession` (erster und zweiter Aufruf, Wächter ohne `u-boot.yaml` öffnet die Sitzung); zusätzlich Exit 11 ohne `docker` im `PATH`.
- Exit 2: `--json` (Umschlag mit `"subcommand":"shell"`, `exitCode: 2`), `--dry-run`, `--diff`, Argument.

## Nutzeraussagen gegen Test oder Messung (Handbuch 4.9/4.11, READMEs, `CHANGELOG.md`, Hilfe, Kontrakt)

Gelesen werden die hinzugefügten Zeilen des Diffs; jede Zusage hat einen Beleg:

- "Start wie `up` ohne `--recreate`", Benutzer und Workspace der Ergebniszeile, Umgebung aus `"NAME": "${localEnv:NAME}"`, ungesetzte Variable fehlt, Verweise anderer Form nicht: `TestDevcontainerShell_StartIsUpWithoutRecreate`, `TestSessionEnvNames`, Docker-Test `Session`, mein Binary-Lauf.
- Ausweisung vor dem Hilfs-Container, auch unter `--quiet`, "Dauer des Starts"; Ausgabe nur auf der Fehlerausgabe, keine Erfolgszeile; `--quiet` hält Fortschritt zurück, Sitzungsausgabe unverändert: `TestDevcontainerShell_SocketDisclosureNamesTheStart`, `TestExecute_DevcontainerShell_HumanOutput`/`_Quiet`, mein Lauf. Die `--quiet`-Zelle des Handbuchs (Ausnahmen 1 und 2, Absatz *Fallbacks nach dem Start*) hat als Wortlaut keinen Test (Review-only, vom Plan so benannt); das Verhalten dahinter ist gepinnt.
- Terminal an der Standardeingabe angeschlossen, sonst Eingabe von der Standardeingabe: `TestExecute_DevcontainerShell_Terminal`, `TestIsTerminal`, Docker-Test `terminal_under_script` mit Gegenprobe über eine Pipe, im Fundament `TestShell_Docker_Terminal_UnderScript`. Die Grenze der Erkennung (ein anderes Zeichengerät als ein Terminal) steht im Code-Kommentar und im Plan, nicht im Handbuch; zugesagt wird sie nicht.
- `failed` → Exit 11 ohne Shell, auch bei laufendem Container; ohne `u-boot.yaml` öffnet die Shell; `up --recreate` als Weg; Prüfung vor der Shell → Exit 12: `FailedStatus…`-Tests (Application und Docker), `…SentinelsAndExitCodes`, `TestShell_Docker_Probe` des Fundaments.
- Exit 0 nach Beginn der Sitzung, was immer der Status war (Docker-Test: Skript mit `exit 3`; mein Lauf): belegt.
- "Beenden Sie die Sitzung mit `exit`. Beendet sich `u-boot` auf anderem Weg (etwa mit `kill -9`), kann die Shell … weiterlaufen": Wortlaut nach [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Entscheidung 4, getragen von `TestShell_Docker_ClientEnd` und `TestShell_Docker_Terminal_UnderScript` (beide in meinem Lauf `PASS`). Keine Aussage zu Signalen, Login-Shell oder Lebensdauer im Handbuch, in der Hilfe oder im `CHANGELOG.md`: bestätigt.
- Lebenszyklus-Befehle: das Handbuch sagt bei `shell` nichts zum `postAttachCommand` (Suche in `docs/user/`, beiden READMEs und `CHANGELOG.md`: nur die bestehenden Sätze zu `up`, getragen vom Test zu `up`). Die in [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md) geforderte Zurückhaltung ist eingehalten.
- `--json` ergibt Umschlag mit `subcommand: shell` (Zeile "nur im Umschlag der Abweisung"), `--dry-run`, `--diff`, Argumente → 2; `--yes`/`--no-interactive` ohne Wirkung: `…RejectsFlagsAndArguments`, `…YesAndNoInteractiveHaveNoEffect`, mein Lauf. Die Handbuch-Bedingung "nur bei `config`, `egress`, `devcontainer`, `template`" der Zeile `subcommand` bleibt ungepinnt (Bestand, vom Slice benannt).
- Hilfe von `shell` ohne Anforderungskennung: `TestExecute_DevcontainerShell_Help`; mein Aufruf von `--help` zeigt keine Kennung, auch nicht in der Exit-Code-Zeile.
- READMEs und `CHANGELOG.md`: Aussagen decken sich mit dem Handbuch; "Promised on a Linux host (`amd64`)" entspricht der Messbasis ([`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 6).

## Plan-vs-Code-Diff (§3 gegen `git diff --stat 41e0b262..HEAD`)

- Jede Datei des Diffs hat eine Zeile der Plan-Tabelle oder steht in *Ist gegen Plan*: Port, Application (`devcontainer_start.go`, `devcontainer_fallback.go`, `export_test.go`), CLI (`devcontainer.go`, `cli.go`, `terminal.go`, `export_test.go`), `cmd/uboot`, Handbuch, READMEs, `CHANGELOG.md`, Kontrakt, Roadmap, `open/README.md`. Unter *Ist gegen Plan* geführt statt in der Tabelle: `terminal.go`, die drei neuen Testdateien (`devcontainer_shell_test.go` in Application und CLI, `devcontainer_shell_docker_test.go`), beide `export_test.go`. Die Plan-Tabelle selbst trägt diese Zeilen nicht (offen aus dem Nachlauf, siehe unten).
- Keine Datei im Diff außerhalb Plan und *Ist gegen Plan*, kein Eintrag der Tabelle ohne Diff, außer den bewusst späteren Zeilen (Beobachtungs-Register, Welle-Datei, Closure) und `spec/`, `.golangci.yml`, `Makefile`, `Dockerfile`, `.github/`: unverändert, bestätigt.
- Die Dateien unter `done/`, `observations/`, `docs/reviews/` und die zwei Slices unter `open/` stammen vom Verweis-Nachzug des Make-Targets für Slice-Moves; die Verschiebung der Datei selbst ist ein reiner Move (`7110bf46`).
- `//nolint`, `TODO`, `FIXME` im Diff von `internal` und `cmd`: 0 Treffer.

## Befunde je Stufe

**HIGH/MEDIUM:** keine.

**LOW**

- **V-1 · Kontrakt §7, Tabelle der Mutations-Matrix hat falsche Zellenzahl** (`docs/maintainer/cli-json-contract.md`, Zeilen der Zeilen `devcontainer down` und `devcontainer shell` der Tabelle in §7). Kopf: Name plus 8 Spalten. Die Zeile `devcontainer down` hat im Diff eine Zelle verloren (8 statt 9 Zellen, vorher 9), die neue Zeile `devcontainer shell` hat eine zu viel (10 statt 9). In der Darstellung verrutscht `down` um eine Spalte, und die Zelle zu viel bei `shell` wird je nach Renderer abgeschnitten. Kein Gate fängt Tabellen-Zellenzahlen; der Review hat es nicht gefunden. Berichtigung: Zelle von `shell` nach `down` zurückgeben. Ob Folgen für die Aussage: keine (alle Zellen sind `—`), nur Darstellung.
- **V-2 · Beleg-Tabelle des Slice nennt einen Test, den es nicht gibt.** Die Zeile "Prüfung vor der Shell scheitert → `12`" verweist auf `TestExecute_DevcontainerShell_SessionErrorPassesThrough`; keine Datei in `internal` oder `cmd` definiert ihn (Suche nach dem Funktionsnamen: 0 Treffer). Die Aussage ist belegt, aber durch den Fall "the check before the session fails" in `TestExecute_DevcontainerShell_SentinelsAndExitCodes` (Exit 12, ohne Standardausgabe). Berichtigung: Namen in der Tabelle ersetzen. Der Rest der Tabelle (21 Namen) löst auf.
- **V-3 · Offen aus dem Nachlauf (L-1), im Plan nicht behoben:** Der Abschnitt *Fix-Runde* in §3 sagt "siehe den Bericht an den Hauptagenten", das ist kein auflösbares Artefakt, und die beiden `devcontainer_shell_test.go` fehlen in der Plan-Tabelle (sie stehen nur unter *Ist gegen Plan*). Der Planner behebt beides in der Closure; ich bestätige nur, dass es beim Stand `302b3a7b` noch so steht.

**INFO**

- **I-1 · Übergabe `--quiet` an die Spezifikation:** Der Plan dieses Slice trägt zu I-1 des Nachlaufs (die Spezifikation führt `shell` in den `--quiet`-Aufzählungen nicht) keine eigene Zeile; die Übergabe steht im Plan [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/open/slice-v2-spec-devcontainer-shell-und-zweites-up.md) in der Zeile zur Verbosity-Spezifikation ("Liste der ändernden Subkommandos", "Ausnahme 1", "`DegradationFallback` auch bei …", "`shell` nach dem freigegebenen Wortlaut" von [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)). Sie ist also getragen, aber pauschal; sie nennt nicht ausdrücklich, dass `shell` den Fortschritt wie `up` hält und die Sitzungsausgabe unverändert lässt. Empfehlung: die Closure-Notiz dieses Slice nennt den Folge-Slice in §7 mit dieser Einzelheit. Der Stand der Spezifikation ([`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) ohne `shell`) ist durch den Plan des Slice ("`spec/` unverändert") gedeckt, kein Verstoß.
- **I-2 · Hilfe von `up` und `down`:** unverändert (Entscheidung (j) der Welle); die Hilfe von `shell` ist die einzige ohne Kennung in der Gruppe. Kein Befund, nur Kontext für den Validator.
- **I-3 · Die globale Option `--quiet` zeigt in `--help` den Typnamen `up`** ("`--quiet up`"): bekannter Bestand ([`BEO-CLI/hilfetext-backtick-wort-als-typname`](../plan/planning/observations/BEO-CLI/hilfetext-backtick-wort-als-typname/observation.md)), nicht Gegenstand.
- **I-4 · Terminal-Pfad am Stück:** weiterhin kein einzelner Test von einem echten Terminal des Nutzers bis zur Sitzung; drei Pins (CLI-Test der Weitergabe, Test unter `script` am Binary, Test unter `script` am Adapter) tragen die Zusage, der Nutzerpfad ist Sache des Validators (Risiko in §6 bleibt "weiter offen").

## Traceability

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Docker-Tests am Binary, Unit-Tests der Beleg-Tabelle, Handbuch 4.9 |
| [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) | `TestExecute_Devcontainer_NamesShell`, Befehlsübersicht im Handbuch, `--help` |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | `TestExecute_DevcontainerShell_Quiet`, Mutationen 5, 10, 13, mein Lauf unter `--quiet` |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | `TestExitCode_DevcontainerShellUsageSentinels`, `TestExecute_DevcontainerShell_SentinelsAndExitCodes`, Exit 0/2/10/11 am Binary |
| [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) | Entscheidungen 1 bis 6 gegen Code und Tests (Start ohne `--recreate`, Ausgabe auf der Fehlerausgabe, Exit 0 nach Sitzungsbeginn, Terminal, Wiring ohne neue Regel, Messbasis) |
| [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) | Satz zum Ende im Handbuch nach Entscheidung 4, Tests im Fundament (eigener Lauf `PASS`) |
| [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md) | keine Aussage zu `postAttachCommand` bei `shell` |
| [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) | `FailedStatus…`-Tests, Mutation 8 |
| [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), [`ADR-0032`](../plan/adr/0032-katalog-features-socket-und-lockerungen.md) | Text der Ausweisung ("duration of the start", Teilaussage zu `docker-cli`) am Binary gelesen; Handbuch und READMEs sagen zur Sitzung nur, was der Befehl einbindet |

Carveouts: Neu: keine. Gelöst: keine. Unverändert: alle (Inventar im Diff nicht berührt).

Nicht ausgeführt: CI-Lauf auf dem Push-Tip (nicht verfügbar, siehe oben); `make test-docker` als Ganzes (Zeitlimit, gezielter Lauf statt dessen); DinD mit Engine 28.0.4 (kein Engine-Text gepinnt); `make image-scan`, `make govulncheck` (kein Security- oder Release-Pfad berührt, `make ci` ist Gegenstand der CI).

Commit / Artefakt: Stand `302b3a7b`, dieser Report.

## Hinweis an den Planner (für die Closure)

1. V-1 (Tabelle in §7 des Kontrakts), V-2 (Name in der Beleg-Tabelle), V-3 (Plan-Nachtrag und Plan-Tabelle) vor dem Move nach `done/` berichtigen; V-1 ist die einzige Berichtigung außerhalb des Slice-Textes.
2. Die offenen DoD-Zeilen (CI `integration-docker` und `ci` auf dem Push-Tip, Haken, Roadmap "Nichts in Arbeit.", Register, Ausgang je Risiko, drei Paarungen) bleiben der Closure; der Nachweis der CI ist erst nach dem Push möglich.
3. Closure-Notiz: Hinweise an den Projektinhaber aus der *Umsetzung* (Zeichengerät-Grenze der Erkennung, Kommentare zählen bei `sessionEnvNames` nicht, "Variable eines Agenten" als zusätzlicher `remoteEnv`-Eintrag) und I-1.

**GO für die Closure** nach V-1 bis V-3.
