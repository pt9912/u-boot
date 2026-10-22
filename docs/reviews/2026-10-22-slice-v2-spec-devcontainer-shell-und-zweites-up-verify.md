# Verification-Report: [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/done/slice-v2-spec-devcontainer-shell-und-zweites-up.md) — 2026-10-22

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD, [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) (als Quelle: in Spezifikation und Sicht dürfen sie nicht stehen), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)?

**Gegenstand:** `git diff --stat 0d653453..HEAD` (`HEAD` = `803ee13c`, 15 Commits ab `origin/main`, 25 Dateien, 488 Einfügungen, 138 Löschungen; Arbeitsbaum zu Beginn sauber, nichts gepusht). Eingang: Slice §1 bis §8 mit *Entscheidungen beim Start*, *Vorab-Suche*, *Ist gegen Plan*, Beleg-Tabelle (B1 bis B29), *Streichungen*, *Sensoren und Abgleich*, *Fix-Runde*; Review-Report [`2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up`](2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up.md) und [Nachlauf](2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up-nachlauf.md). Die Sensor-Zeilen des Implementers habe ich nicht als Beweis genommen; die Läufe unten sind eigene.

**Gesamturteil: GO für die Closure.** Keine Verletzung einer DoD-Zeile, einer Entscheidung oder einer Hard Rule gefunden. Drei Befunde, alle INFO oder LOW, keiner blockiert (V-1 bis V-3); der Planner entscheidet. V-1 beantwortet die offene Frage I-1 des Reviewers.

## Verification Evidence

Scope:

- Slice: [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/done/slice-v2-spec-devcontainer-shell-und-zweites-up.md), fünfter von sechs der Welle [`welle-devcontainer-shell-und-zweites-up`](../plan/planning/welle-devcontainer-shell-und-zweites-up.md).
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen), [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md).
- Artefakte: `spec/spezifikation.md`, `spec/architecture.md`, der Slice, `roadmap.md`, die zwei Reports des Reviews; der Verweis-Nachzug des Move berührt fünf Slices in `done/`, die Welle-Datei, `open/README.md`, elf Reports und eine Evidence-Datei.

### DoD-Abgleich (§2 des Slice, Stand `803ee13c`)

- [x] **Spezifikation** (DoD 1) — Beleg: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) mit Unterabschnitt *`shell`: Startphase und Sitzung* (Zeilen 825 bis 845 der Datei), Urteil nach dem Zustand, Zeile `already running`, Feld `alreadyRunning`, Absatz zu den Lebenszyklus-Befehlen, Tabelle der Exit-Codes mit Spalte `shell`, *Reichweite* mit drei Unterbefehlen; [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (Kopf, *Ermittlung durch den Host*, Label, Container-Erkennung, *Abbildung*), [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) (*Abbruch*), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (vier Stellen: Zeilen 584, 588, 590, 595), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Absatz, Zeile 980), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) (Zeile 1225), Änderungshistorie (Zeile 1410). Gelesen gegen Lastenheft [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 488 bis 499): nichts darin weiter als das Lastenheft, bis auf die Spiegelung von V-1.
- [x] **Sicht** (DoD 2) — Beleg: `spec/architecture.md` §2.2 bis §2.6, die Zeile der Komponententabelle, §7.1 und §7.2 geändert; Diff (+52/−26) auf Verweise nach unten durchsucht: 0 Treffer (siehe Stufe Strata). Der Titel von [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und sein Anker sind unverändert (Zeile 756), keine neue `SPEC-<NNN>`.
- [x] **Beleg-Tabelle** (DoD 3) — Beleg: 60 Testnamen in B1 bis B29; `grep -rn "^func <Name>(" internal cmd` je Name ergibt genau eine Definition (60 von 60, eigener Lauf). Zehn Docker-Tests und der Test unter `script` liefen eigen und bestanden; 18 Mutationen färben die zugehörigen Tests rot (unten).
- [x] `make gates` grün — Beleg: eigener Lauf Exit 0 (`coverage-gate` 92,00 % gegen 90 %, `docs-check` 714 Dateien 0 Befunde, `ablage-pruefung` 26 Fälle 0 Fehlschläge); `make doc-trace` 84 Anforderungen, 0 Waisen.
- [ ] Review durchgeführt — erfüllt (zwei Reports, GO); der Haken ist Sache der Closure.
- [ ] Doku-Update Roadmap („In Arbeit“ gesetzt, Rücksetzung bei der Closure), Closure-Notiz (§7 ist der Platzhalter), Register, Risiken aus §6 (alle „weiter offen“), Paarungen — **bewusst offen**: Closure des Planners.

Sensoren (eigene Läufe):

| Sensor | Ergebnis | Beleg |
| --- | --- | --- |
| `make gates` | pass, Exit 0 | `gates-exit=0`, Coverage 92,00 %, `docs-check` 0 Befunde |
| `make doc-trace` | pass, Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-immutable RANGE="v0.13.1..HEAD"` | pass, Exit 0 | 111 Commits, 0 Befunde |
| `make doc-commits RANGE="origin/main..HEAD"` | pass, Exit 0 | 15 Commits, 0 Befunde |
| `make doc-planning` | pass, Exit 0 | 0 Befunde |
| `make docs-check` (nach dem Report) | pass, Exit 0 | 715 Dateien, 0 Befunde (der Report trägt jede Kennung verlinkt) |
| Docker-Tests (Image `u-boot:test-docker-tools`, `-tags docker`, `-run` auf die zehn Namen der Tabelle über `./internal/e2e/... ./internal/adapter/...`) | pass, Exit 0 | `internal/e2e` 70,8 s, `internal/adapter/driven/docker` 35,9 s; `TestDevcontainerShell_Docker_Session` (drei Teile, darunter `terminal_under_script`), `…_FailedStatusOpensNoSession`, `TestDevcontainerUp_Docker_S8_FailurePicture` (sieben Teile), `…_S8_SecondUpAfterStartScriptAbort` (zwei Teile), `…_AlreadyRunningLifecycleCounters`, `TestShell_Docker_Probe` (fünf), `…_NoTerminal_UserWorkspaceAndEnvironment`, `…_B18_SessionDoesNotBegin` (acht), `…_Terminal_UnderScript` (zwei), `…_ClientEnd` (vier) alle `PASS`, 0 `SKIP`; Docker 29.9.2 (Engine und Client), Linux `amd64` |
| `make test-docker` als Ganzes | nicht gelaufen | kein Code im Diff; das e2e-Paket nähert sich dem Zeitlimit, der gezielte Lauf deckt jeden Docker-Test der Beleg-Tabelle ab |
| DinD-Engine 28.0.4 | nicht gelaufen | kein Code im Diff; die Aussagen zur Engine 28.0.4 stehen im Fundament-Slice auf beiden Engines gemessen und werden hier nicht neu behauptet |
| `make verify-depguard` | nicht gelaufen | keine Importregel, kein Code berührt (`git diff --stat 0d653453..HEAD -- internal cmd …` leer) |
| `make ci` / `make fullbuild` | nicht gelaufen | Schritt der Welle-Closure |

**Reine Doku.** `git diff --stat 0d653453..HEAD -- internal cmd docs/user README.md README.de.md CHANGELOG.md spec/lastenheft.md docs/plan/adr` ist leer (eigener Lauf). `git log 0d653453..HEAD -- internal cmd` ist leer, die Spezifikations-Commits liegen damit nach dem Move des Befehl-Slice nach `done/` (Closure-Trigger §5). Der Move `b7260cd1` ist ein reiner Rename (0 Zeilen), der Inhalt folgt in eigenen Commits (`AGENTS.md` §3.3).

## Plan gegen Diff (`AGENTS.md` §3.10)

| Datei im Diff | Plan-Tabelle §3 |
| --- | --- |
| `spec/spezifikation.md`, `spec/architecture.md` | Zeilen 1 und 2, mit „Ist (2026-10-22)“ (inklusive der Funde §2.3, `ARC-008`, §7.1, §7.2) |
| `docs/plan/planning/in-progress/roadmap.md` | Zeile 3 (Roadmap-Zeile „In Arbeit“) |
| der Slice (R, 59 % Ähnlichkeit im Gesamtdiff, Move selbst rein), `open/README.md`, die Welle-Datei | Zeilen 4 und 5 |
| fünf Slices in `done/` (je 8 bis 26 Zeilen), elf Reports, eine Evidence-Datei | Zeile 5 „Ist (2026-10-22)“: Nachzug des Move (Commit `322f0ab9`), elf Reports und die Evidence-Datei genannt; die fünf `done/`-Slices fallen unter „die übrigen Slices der Welle“ (alle fünf in der Welle-Datei geführt) |
| zwei neue Reports unter `docs/reviews/` | Zeile 6 |
| `spec/lastenheft.md`, `docs/plan/adr/`, `internal/`, `cmd/`, `docs/user/`, beide READMEs, `CHANGELOG.md` | Zeile 8 „unverändert“ — Diff leer, bestätigt |

Jede Datei im Diff hat einen Plan-Eintrag, jeder Eintrag mit Diff-Anspruch hat einen Diff; der Eintrag „Observations (update, Closure)“ ist noch ohne Diff und gehört zur offenen Closure (kein Mangel).

## Stufe 1: Entscheidungen (ADR-0030, ADR-0031, ADR-0033, ADR-0034) als Quelle

| Entscheidung | Befund |
| --- | --- |
| Urteil nach dem Zustand in beiden Pfaden ([`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) 1) | in der Spezifikation (Zeile 796) und in der Sicht (§2.2, §7.1, §7.2); am Binary bestätigt (Stufe Binary), Mutationen A4 und A5 rot |
| Zeile `already running`, Feld `alreadyRunning`, Abfrage mit Grenze ([`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md) 2) | Zeile 804 bis 807; A1, A2, A15 rot, Binary bestätigt |
| `shell` = `up` ohne `--recreate` ([`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) 1), Socket nur für den Start (2) | Zeile 800 und 825; `TestDevcontainerShell_StartIsUpWithoutRecreate` und `…SocketDisclosureNamesTheStart` vorhanden und im Gates-Lauf grün; Ausweisung am Binary „for the duration of the start“ |
| Folgepflicht [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md): Satz „keine Lebenszyklus-Befehle“ nicht übernommen, Befehle einzeln mit Messgrenzen | Zeile 792: drei Befehle mit Verhalten, Messgrenzen (gepinnte CLI, Zeichenketten, Linux `amd64`, lokale Engine, Docker 29.9.2) und „nicht zugesagt“-Liste (`initializeCommand`, `onCreateCommand`, `updateContentCommand`, Array- und Objektform, scheiternde Befehle, andere Hosts); `grep -i "keine Lebenszyklus"` in `spec/`: 0 Treffer; Mutation A7 färbt den Zähler-Test rot |
| Folgepflicht [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md): „Endet `u-boot` durch ein Signal, endet der Client mit ihm“ und „Login-Shell“ nicht übernommen | `grep -i "endet der Client mit ihm"` 0 Treffer; „Login-Shell“ steht nur verneint, im Satz zur Art der Shell mit Terminal, den `TestShell_Docker_Terminal_UnderScript` trägt (`shopt -q login_shell`, `$0`, `~/.bashrc`, `~/.profile`); Zusage-Menge nur mit Test (`ClientEnd` vier Fälle, `UnderScript`) und Messbedingungen (Docker 29.9.2, `amd64`, `sleep 3600`, `script`, ein Lauf je Fall) |

Eine Verletzung einer Entscheidung wurde nicht gefunden.

## Stufe 2: Strata-Regeln (`AGENTS.md` §3.4)

- Diff von `spec/` (`git diff -U0 0d653453..HEAD -- spec/`, nur Zeilen mit `+`): 0 Treffer für `ADR-<Zahl>`, `slice-`, `welle-`, `roadmap`, `BEO-`, `v0.<n>`, `Meilenstein`, `M<n>-T`, `LOC`, `Closure`, `Review`, `Verifier`, `Nachlauf`, `noch nicht`, `ausstehend`, `bisher`, `nicht mehr`, `no longer`, `former`, `originally`. Treffer „Phase“ (1): *Startphase*/*Sitzungsphase* als Begriff des Befehls, kein Meilenstein. Treffer „früher“ (5): Zustand „aus einem früheren Start des Containers“, den das Lastenheft nennt, und die Zeile der Änderungshistorie. Sicht (`spec/architecture.md`): 0 Treffer für ADR, Slice, Welle, Roadmap, Release-Version, Phase, Review, früher.
- Änderungshistorie (Zeile 1410, `2026-10-22`): ein Eintrag, „deckungsgleich mit Lastenheft 0.21.0 und Code“ (die Fassung des Lastenhefts steht dort im Muster der Vorgänger-Zeilen, `0.18.0`, `0.19.0`), nennt alle geänderten Stellen und die nicht zugesagten Fälle; stimmt mit dem Diff überein. Die Aufzählung „[`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) und [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten)“ entspricht den Hunks (Zeilen 207, 584 bis 595, 643 bis 688, 752, 980, 1225).
- Hinweis zur Form (INFO, kein Befund): Zeile 842 trägt „(aus der Form abgeleitet, nicht gefahren)“; das ist eine Aussage über den Beleg im Text der Zusage. Sie benennt eine Grenze ehrlich und ist durch die Regeln der Spezifikation (Zielbild, Messbedingung als Bedingung) gedeckt.
- Nutzer-Doku ohne Diff: `docs/user/`, beide READMEs, `CHANGELOG.md` unverändert (Regel 6 von `AGENTS.md` §5 nicht berührt).

## Mutationen (selbst gefahren, 18 Änderungen, alle rot)

Unit-Läufe: `go test -count=1` über das jeweilige Paket im Image `u-boot:test-docker-tools`; Docker-Läufe mit `-tags docker`. Jede Änderung per Sicherung zurückgespielt, nach jedem Lauf `git diff -- <Datei>` leer; am Ende `git status --short` leer und `git diff --stat -- internal cmd` leer (Arbeitsbaum zu Beginn sauber, daher kein Verlust uncommitteter Änderungen).

| Nr | Spec-Aussage | Änderung am Code | Ergebnis | Rot in |
| --- | --- | --- | --- | --- |
| A1 | Zeile `Devcontainer already running` | `outcome = "already running"` → `"started"` (`devcontainer.go`, CLI) | rot | `TestExecute_DevcontainerUp_AlreadyRunning/human` |
| A2 | Feld `alreadyRunning` | Feld immer `false` (`devcontainer_start.go`) | rot | `TestDevcontainerUp_AlreadyRunning/running_before,_same_ID` |
| A15 | Abfrage entfällt mit `--recreate` | Schutz `if recreate` entfernt | rot | `TestDevcontainerUp_AlreadyRunning/recreate_skips_the_query` |
| A3 | `--json` bei `shell`: Exit 2, Umschlag mit `subcommand: "shell"` | Umschlag mit `"up"` statt `"shell"` | rot | `TestExecute_DevcontainerShell_RejectsFlagsAndArguments` (`json`, `json_before_the_command`, `diff`, `dry-run`) |
| A4 | `failed` → 11 auch nach gelungenem Lauf und aus früherem Start | `case failed:` → `case false:` | rot | `TestDevcontainerUp_FailedStatusAfterSuccess` (sechs Teile), `…QueryBeforeTheLauncherDoesNotSpoilTheStatusReading`, `TestDevcontainerUp_ErrorsNameUp`, `TestDevcontainerShell_FailedStatusOpensNoSession/after_a_successful_run_of_the_CLI` |
| A5 | `failed` → 11 statt 12 nach Fehler der CLI | `case launchErr != nil && failed:` → `case false:` | rot | `TestDevcontainerUp_RuntimeFailureClassification` (vier Teile), `TestDevcontainerShell_FailedStatusOpensNoSession/after_a_failed_run_of_the_CLI` |
| A6 | `vfs` vor `egress` | `firstFailedScript` in umgekehrter Reihenfolge | rot | `TestDevcontainerUp_FailedStatusAfterSuccess/the_first_failed_capability_is_named` |
| A20 | Meldung nennt „this call did not run the start script“ | Zeichenkette in `failedAfterStartText` geändert (global) | rot | `TestDevcontainerUp_FailedStatusAfterSuccess` (`vfs_failed`, `daemon_failed`, `egress_failed`) |
| A20d | dasselbe am Prozess | wie A20, Docker-Lauf | rot | `TestDevcontainerUp_Docker_S8_SecondUpAfterStartScriptAbort/with_u-boot.yaml` |
| A7 | Zähler der Lebenszyklus-Befehle (`postCreateCommand` und `postStartCommand` nicht erneut gegen einen laufenden Container) | `--remove-existing-container` immer gesetzt (`devcontainer_launcher.go`) | rot | `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters` (Docker-Lauf) |
| A8 | `-e NAME` ohne Wert | `-e NAME=<Wert>` (`shellArgs`) | rot | `TestShellArgs/variables_by_name`, `…/terminal_with_a_variable`, `TestLauncher_Shell_PassesVariablesByNameOnly` |
| A9 | Benutzer der Sitzung aus der Ergebniszeile | `User: "root"` | rot | `TestDevcontainerShell_OpensTheSessionWithTheResultOfTheStart` |
| A10 | OCI-Zeile nur bei Status 127 (Status 1 ist Sitzung) | Bedingung `exit != 0` | rot | `TestLauncher_Shell_SessionEndIsNoError/status_1,_the_OCI_line_as_the_only_output` |
| A11 | Terminal-Erkennung je Dateiart (Nullgerät ausgenommen) | Prüfung auf das Nullgerät entfällt | rot | `TestIsTerminal`, `TestExecute_DevcontainerShell_Terminal/the_default_probe` |
| A12 | Grenze (e): „einzige Zeile der Standardausgabe“ | `stdout.singleLine()` entfernt | rot | `TestLauncher_Shell_SessionEndIsNoError` (zwei Wächterfälle zu weiterer Ausgabe) |
| A13 | Grenze (e): Doppelpunkt als Namensgrenze | Doppelpunkt aus `missingUserLead` entfernt | rot | `TestLauncher_Shell_SessionEndIsNoError` (zwei `vscode2`-Wächterfälle) |
| A14 | Umgebung der Sitzung ohne Kommentare, Form `"NAME": "${localEnv:NAME}"` | `withoutComments` nicht aufgerufen | rot | `TestSessionEnvNames` (Zeilenkommentar, Blockkommentar, unterminierter Blockkommentar) |
| A16 | Standardausgabe leer (ohne Terminal) bei `Error response from daemon: …` | Schutz am Kopf von `clientMessageOnStderr` entfernt | rot | `TestLauncher_Shell_SessionEndIsNoError/status_1,_the_message_names_this_container_but_the_session_wrote_to_standard_output` |
| A17 | `Cannot connect …` nur, wenn die Engine auch bei erneuter Frage nicht antwortet | `engineUnreachable(ctx)` → `true` | rot | `TestLauncher_Shell_SessionEndIsNoError/status_1,_the_engine_message_while_the_engine_answers` |

Zur Ehrlichkeit des Verfahrens: A11 brach im ersten Versuch am Bau ab (`null` ungenutzt), A20 traf im ersten Versuch den Kommentar über der Funktion statt der Zeichenkette und blieb grün; beide sind mit geänderter Ersetzung wiederholt (Zeilen oben). Beim Lauf des Pakets `internal/adapter/driven/docker` ohne eingehängten Docker-Socket fällt `TestProbe_Version_HappyPath` rot, auch unverändert (Umgebung des Ad-hoc-Laufs, nicht der Mutation; `make gates` und die Docker-Läufe mit Socket sind grün).

## Binary (eigener Lauf)

Binary per `go build ./cmd/uboot` im Image `u-boot:test-docker-tools`, Wegwerf-Projekt `u-boot init vbproj --no-git --devcontainer --sandbox` unter dem Scratchpad (Pfad unter `/tmp`, damit die Engine ihn sieht):

- `devcontainer shell --json`: Standardausgabe `{"status":"error","command":"devcontainer","subcommand":"shell","diagnostics":[{"level":"error","code":<Kennung der Exit-Code-Anforderung>,…}],"exitCode":2}` (Code der Diagnose: [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)), Prozess-Exit 2.
- erstes `up`: Zeile `Devcontainer started: container 67a7b7544312, remote user vscode, workspace /workspaces/vbproj`, Exit 0.
- zweites `up`: Zeile `Devcontainer already running: container 67a7b7544312, remote user vscode, workspace /workspaces/vbproj`, Exit 0; die Ausweisung des Sockets steht davor („for the duration of the call“). `--quiet`: Zeile entfällt, die Ausweisung bleibt, Exit 0. `--json`: `status: warn`, `exitCode: 0`, `data.alreadyRunning: true`, `diagnostics` mit dem Code der Anforderung [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Ausweisung).
- `shell` mit Skript auf der Standardeingabe (`id -un; pwd; exit 3`): `vscode`, `/workspaces/vbproj`, Exit 0 trotz `exit 3`; Standardfehlerausgabe trägt die Ausweisung mit „for the duration of the start“.
- `failed` aus einem früheren Start: `u-boot.yaml` mit `egress.enabled` und `nestedRuntime: podman`, im laufenden Container die Statusdateien `vfs` und `egress` mit `result=failed` und `at=2026-10-22T20:00:00Z` (Datum vor dem Lauf) angelegt. `up --json`: Exit 11; `up` (Text): Exit 11 mit „…/dev/fuse was not available… (reported at 2026-10-22T20:00:00Z); the container is running and this call did not run the start script; `u-boot devcontainer up --recreate` creates the container anew“ — die erste Fähigkeit der Reihenfolge (`vfs`) ist genannt, nicht `egress`. `shell`: Exit 11, `echo SESSION-OPENED` wurde nicht ausgeführt (keine Sitzung).
- Aufräumen: Container, Volume `vbproj-workspace-…` und die zwei Images `vsc-vbp-…` entfernt; Container-Liste und Volume-Liste gleich dem Stand vor dem Lauf (zwei fremde Container `pgr-gp`, `pgr-dev2` unberührt); die Docker-Tests (zehn Namen) hinterließen nichts; Wegwerf-Verzeichnisse und das Cache-Volume des Mutationslaufs entfernt.

## Streichungen und Grenzen (a) bis (e)

Die Streichungstabelle des Slice habe ich gegen die Zeilen der Spezifikation gelesen:

- **Nichts Getestetes verschenkt:** Die Sätze, die ein Test trägt, stehen in der Spezifikation (Erkennung je Meldung, Wächterfälle Status 1 und Doppelpunkt, Messbedingungen, Zusage-Menge zu Signalen, Art der Shell). Nicht übernommen sind nur Aussagen, die der Baum nicht trägt (Mutationen M36 bis M44 und Wegwerf-Messungen des Fundament-Verifiers, `ioctl`, `: unknown`, `126` für die OCI-Zeile). Der Code lässt `126` für die OCI-Zeile zu (`exit == 126 || exit == 127`), die Spezifikation nennt nur `127`: schwächer als der Code, deshalb vertretbar.
- **Nichts Ungetestetes behalten:** Die Aussage zur Umgebung der Sitzung (`"NAME": "${localEnv:NAME}"` mit gleichem Namen in Schlüssel und Verweis, sortiert, ohne Duplikate, Kommentare) trägt `TestSessionEnvNames` (Fälle `another name in the reference`, `a default`, `a reference inside a longer value`, `a reference in mounts`, `a key that is no plain name`, Kommentarfälle; Mutation A14 rot, A8 für die Weitergabe als Name). `shell` ohne `u-boot.yaml` (Zeile 845) trägt `TestDevcontainerShell_FailedStatusNobodyReadsDoesNotStopTheSession` und der Wächter in `TestDevcontainerShell_Docker_FailedStatusOpensNoSession` (beide grün).
- **Grenze (e)** (Zeile 842) und Spalte `0` der Exit-Code-Tabelle: die Form pinnen die Wächterfälle von `TestLauncher_Shell_SessionEndIsNoError` und der positive Fall von `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` (A12 und A13 rot). Der Erzeuger in der Sitzung ist abgeleitet; die Spezifikation sagt es im Text.

## Offene Punkte des Reviewers

- **I-1 (Spalte `0` und Lastenheft):** siehe V-1.
- **I-2 (Grenze (e) steht in der Streichungstabelle):** bestätigt. Die Zeile steht unter „Streichungen“, obwohl die Grenze in der Spezifikation steht (Spalte „Stand in der Spezifikation“ sagt „als Grenze in der Spezifikation … und als Ausnahme in der Spalte `0`“). Die Lage ist richtig, die Ablage etwas unscharf; INFO, keine Änderung nötig.
- **I-3 (`Cannot connect` bei nicht leerer Standardausgabe ohne Test):** bestätigt als Lücke im Test, nicht in der Zusage. Der Code (`clientMessageOnStderr`) verlangt ohne Terminal bei jeder Meldung der Fehlerausgabe eine leere Standardausgabe; Mutation A16 (Schutz entfernt) färbt nur den Fall `Error response from daemon: …` rot, ein Fall `Cannot connect` mit nicht leerer Standardausgabe besteht nicht. Die Spezifikation sagt es nur beim ersten Spiegelstrich und deckt die übrigen über „eine Meldung des Clients nach anderer Ausgabe der Sitzung gilt als Sitzung“: schwächer als der Code, nicht stärker. INFO (V-3).

## Befunde

**V-1 (LOW, I-1 beantwortet; Planner entscheidet, blockiert nicht): Grenze (e) schränkt eine Zusage des Lastenhefts ein, ohne dass das Lastenheft eine Ausnahme nennt.**
[`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt: „Hat die Shell begonnen, endet der Befehl mit `0`, unabhängig von ihrem Exit-Status.“ Die Spalte `0` der Tabelle in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) nimmt „die Grenzen der Erkennung“ aus. Für die Grenzen (a) bis (d) ist das kein Widerspruch: Dort gilt der Satz des Lastenhefts zu `12` („Prüfung scheitert vor dem Öffnen“) und die Bindung von [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) an `10`, `11` oder `12`; die Spezifikation führt dort nur den nicht trennbaren Rest, und der Code endet für die nicht erkannten Fälle mit `0`, was der Satz „Sitzung hat begonnen“ entspricht. Die Grenze (e) ist die einzige Stelle, an der die Spezifikation für eine **begonnene** Sitzung `12` statt `0` nennt (ein `docker exec -u <Benutzer> <anderer Container>` in der Sitzung gegen eine Engine der Form 28.0.4 als einzige Ausgabe, Status 126): Das ist eine Einschränkung der Zusage in einem abgeleiteten, nicht gefahrenen Randfall. Wertung: kein Widerspruch im Sinn einer Erweiterung (die Spezifikation behauptet nichts, was das Lastenheft verneint, und sie benennt die Abweichung ehrlich); formal gilt bei Konflikt das Lastenheft (Rang 1). Empfehlung: Das Lastenheft braucht keinen Hinweis für die Closure dieses Slice; wenn der Projektinhaber die Zusage „`0` nach Beginn der Shell“ ohne Ausnahme halten will, ist das ein Code-Slice (Erkennung enger fassen), wenn er die Ausnahme will, ein CR mit einem Satz im Lastenheft. Ich ändere das Lastenheft nicht.

**V-2 (INFO): Dieser Slice hat keine Mutation für die Sicht.** Die Sicht ist beschreibend (kein Test); die Fundstellen der Sicht-Tabelle (Symbole in `devcontainer_start.go`, `devcontainer.go`, `terminal.go`, `devcontainer_shell.go`, `devcontainer_launcher.go`, Port-Dateien) habe ich per Suche am Baum bestätigt: `Shell`, `start`, `startOutcome`, `runningContainerID`, `firstFailedScript`, `failedAfterStartText`, `runSession`, `probeShell`, `clientErrorBeforeSession`, `headBuffer`, `isTerminal`, `ErrDevcontainerShellJSONNotApplicable` existieren. Kein Befund, nur Abgrenzung.

**V-3 (INFO, I-3): ungepinnte Randbedingung der Erkennung.** `Cannot connect` mit nicht leerer Standardausgabe ohne Terminal ist im Code Sitzung, im Test nicht gepinnt; die Spezifikation behauptet das Gegenteil nicht (siehe oben). Wenn der Projektinhaber den Fall pinnen will, ist es ein Test in einem Code-Slice, kein Spezifikationsmangel.

## Nicht geprüft, mit Grund

- Dev-Container-CLI-Verhalten für Lebenszyklus-Befehle jenseits des Zähler-Tests (Messgrenzen der Spezifikation; nicht zugesagt).
- Engine 28.0.4 (DinD): keine Code-Änderung in diesem Slice; die Aussagen dazu sind im Fundament-Slice gemessen und gestützt, hier nicht neu gefahren.
- CI-Lauf auf dem Push-Tip: Stand lokal und ungepusht.
- `docs/user/`, READMEs, `CHANGELOG.md`: unverändert (Diff leer), nicht Gegenstand.

## Übergaben (g)

Die Pläne in `open/` ([`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md), [`slice-v1-devcontainer-start-colima-messung`](../plan/planning/open/slice-v1-devcontainer-start-colima-messung.md), [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../plan/planning/open/slice-v1-docker-cli-bau-auf-standard-basis-image.md), [`slice-v1-katalog-features-quellen-und-freigabe`](../plan/planning/open/slice-v1-katalog-features-quellen-und-freigabe.md)) sind im Diff unverändert (`git diff --stat 0d653453..HEAD -- docs/plan/planning/open docs/plan/planning/next`: nur `README.md`, Zeile dieses Slice, Pfad). Die Welle-Datei ändert zwei Zeilen (Pfad des Slice in der Tabelle und in der Reihenfolge, `open/` → `in-progress/`); §3 (Nutzerpfad als Prüfgegenstand) und die Validator-Prüfpunkte (Punkt 7: Bericht des Validators vor dem Release-Cut) sind unverändert und vollständig. Die Aussagen der Spezifikation, die der Validator am Binary prüft (zweites `up` nach `11`, `shell` mit Exit 11 ohne Sitzung), habe ich in „Binary“ vorab bestätigt; das ersetzt den Validator nicht.

## Verification Evidence (Kurzform nach `harness/verification.md`)

Traceability:

| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Beleg-Tabelle B1 bis B29, Docker-Tests, Mutationen A1 bis A20d, Binary |
| [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) | B21 und B27 (`TestExecute_DevcontainerShell_Quiet`, `…Up_Quiet`), Binary `--quiet` |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Exit 2 (A3), 11 (A4, A5, Binary), 12 (Grenzen A10, A12, A13, A17) |
| [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) | Stufe 1: in Spezifikation und Sicht nicht genannt, Folgepflichten erfüllt |

Carveouts: neu — keine; gelöst — keine; unverändert — keine berührt.

Commit / Artefakt: Stand `803ee13c`; dieser Report wird als einziger Commit ergänzt.

**GO/NO-GO für die Closure: GO.** Offen für den Planner: Haken Review/Doku-Update, Closure-Notiz mit Steering-Loop-Lerneintrag, Register, Ausgang der Risiken aus §6, Move nach `done/` (erst Inhalt, dann reiner `git mv`), Entscheidung zu V-1.
