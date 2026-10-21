# Slice v1: Befehl `u-boot devcontainer shell` (Application, Driving-Port, CLI, Ausgabe auf der Fehlerausgabe, Exit-Codes, Docker-Tests, Handbuch, READMEs, `CHANGELOG.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-shell-und-zweites-up`](../welle-devcontainer-shell-und-zweites-up.md) (Release `v0.14.0`). Vierter von sechs Slices; nach [`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md) (die Startphase urteilt nach dem Zustand) und [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md) (Port-Methode `Shell`), vor [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Absatz zu `shell` in der Fassung des CR), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Hilfs-Container nicht Teil des Ergebnisses), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (`0`, `2`, `10`, `11`, `12`, `14`; reservierte Codes), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Umschlag des Exit-Codes `2` mit `subcommand: shell`), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Fallback-Warnungen nach dem Start), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten in der Sitzung); [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (Entscheidungen 1 bis 6), [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) (Entscheidung 1 für die Startphase), [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (Entscheidung 3, Warntext; für `shell` die Variante „für die Dauer des Starts“), [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Entscheidungen 3 und 4).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (neuer Unterabschnitt für `shell`, *Ausweisung des Sockets*), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Ausnahme 1 und `--quiet` bei `shell`), [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Pflicht von `subcommand`, gelesen), `architecture.md` §2.2 bis §2.6. Gelesen, nicht geändert; die Nachführung ist [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md). Der Verweis zeigt aufwärts.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `u-boot devcontainer shell` startet den Devcontainer wie `up` ohne `--recreate` und öffnet danach über `DevcontainerLauncher.Shell` eine interaktive Shell mit Benutzer, Workspace und Kennung aus der Ergebniszeile und den Variablen der gleichnamigen `${localEnv:<NAME>}`-Einträge von `devcontainer.json` als Umgebung; die Sitzung hält keinen Socket; alles, was `u-boot` selbst ausgibt, steht auf der Standardfehlerausgabe; der Befehl endet mit `0`, sobald die Sitzung begonnen hat ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidungen 1 bis 6). Dazu die Docker-Tests am Befehl mit Wächter und die Dokumente.

**Herkunft.** Lücke L-1 des [Validator-Berichts](../../../reviews/2026-10-21-validation-welle-devcontainer-start-ohne-editor.md) (ohne Editor führt der einzige Weg in den Container über `docker exec`, dort gilt `remoteEnv` nicht, im Sandbox-Profil fehlt `GIT_TOKEN`); offene Entscheidung (h) der Closure von [`welle-devcontainer-start-ohne-editor`](../done/welle-devcontainer-start-ohne-editor.md) §7; entschieden in [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (Option B, Name `shell`, nur interaktive Shell, Benutzer und Workspace aus der Ergebniszeile, Lastenheft-CR).

**Gelesener Stand (2026-10-21).** `internal/adapter/driving/cli/devcontainer.go`: Gruppe `devcontainer` mit genau `up` und `down` (`newDevcontainerCommand`, Zeilen 62 bis 83, Langtext mit der Liste *Subcommands*); `ErrDevcontainerSubcommandRequired` nennt im Text `up` und `down` (Zeile 20), `ErrDevcontainerPreviewNotApplicable` ebenso (Zeile 26); `runDevcontainerUp` schreibt die Ausweisung des Sockets im Textmodus auf die **Standardausgabe** (`writeWarnings(stdout, …)`, Zeile 255) und den Fortschritt auf die Fehlerausgabe (Zeile 246) — `shell` braucht beides auf der Fehlerausgabe ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2). `internal/hexagon/application/devcontainer_start.go`: `localEnvNames` (Zeilen 41 bis 56) liest jeden Verweis `${localEnv:…}` aus dem Rohtext, auch in Kommentaren und außerhalb von `remoteEnv`; `socketRelaxationMessage` (Zeilen 31 bis 35) sagt „for the duration of the call“. Die erzeugte Konfiguration trägt `remoteEnv` nur im Sandbox-Profil mit gleichnamigen Einträgen (`internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`, Zeilen 33 bis 38). `cmd/uboot/main.go` setzt `signal.NotifyContext` mit `os.Interrupt` und `SIGTERM` (Zeile 51). Eine Erkennung, ob die Standardeingabe ein Terminal ist, gibt es im Code nicht (Suche nach `ModeCharDevice`, `IsTerminal`, `isatty` und `x/term` ohne Treffer; `go.mod` ohne Terminal-Bibliothek).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port-Methode und Adapter** — [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md); dieser Slice ruft sie auf.
- **`shell -- <Befehl>` mit durchgereichtem Exit-Code** — verworfen ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md), Entscheidung des Projektinhabers c); das Handbuch nennt für einen Befehl mit Exit-Code `docker exec`.
- **Benutzer und Workspace aus Labels der CLI** — verworfen, nicht gemessen ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md), Entscheidung des Projektinhabers d).
- **`--recreate` bei `shell`** — nicht entschieden; die Startphase ist `up` *ohne* `--recreate` ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1). Der Weg nach einem `failed` ist `u-boot devcontainer up --recreate`, danach `shell`.
- **`spec/`** — der CR liegt davor, die Spezifikation danach; jede Bedingung im Handbuch braucht einen Test aus diesem Slice oder entfällt ([`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Ob Hilfe und Ausgabe Anforderungs-Kennungen tragen** (offene Entscheidung (j) der Vorwelle) — der Langtext folgt der Form von `up` und `down` (Zeile „Exit codes (…)“), solange der Projektinhaber nicht anders entscheidet; Risiko in §6.
- **Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, entfernte Engine, Windows, Linux-`arm64`** — nicht zugesagt ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 6).
- **Die Kommentare zu den bare-Rejects in `root.go` und `cli.go`** (Hinweis F-3 aus [`slice-cr-json-subcommand-pflicht`](../done/slice-cr-json-subcommand-pflicht.md) §7) — der Slice ändert die Abweisung von bare `u-boot devcontainer` nicht; berührt er die Stelle doch, zieht er die Kommentare nach und nennt es in §7.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Befehl** nach [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md): Driving-Port (Methode `Shell` am `DevcontainerUseCase` oder ein eigener Use-Case, Wahl des Implementers; Request mit Strömen und Terminal ja oder nein, die Application bleibt frei von konkreter I/O), Application (Startphase wie `up` ohne `--recreate` einschließlich der Regel nach dem Zustand: `failed` ergibt `11` ohne Sitzung; danach `Shell` mit den Werten der Ergebniszeile; die Namensliste als reine Funktion neben `localEnvNames`: nur Einträge der Form `"<NAME>": "${localEnv:<NAME>}"`), CLI (`devcontainer shell` ohne Argumente; Terminal-Erkennung der Standardeingabe im Driving-Adapter; Ausweisung, Fortschritt, Fallback-Warnungen und Fehler auf der Standardfehlerausgabe, nichts auf der Standardausgabe; keine Erfolgszeile; `--quiet` hält den Fortschritt zurück; `--json`, `--dry-run`, `--diff` und Argumente ergeben `2`, mit `--json` als Umschlag mit `subcommand: shell`; `0` nach Beginn der Sitzung unabhängig vom Exit-Status der Shell; `12` bei gescheiterter Probe); Ausweisung mit der Variante „für die Dauer des Starts“. Jeder neue Sentinel steht in **beiden** Klassifikatoren (`ExitCode` und die Diagnostic-Abbildung von `shell`) mit Pin-Test je Sentinel. *Beleg:* Unit- und CLI-Tests nach der Fitness Function (Fälle in §3).
- [ ] **Docker-Tests am Befehl mit Wächter** (Fitness Function, Build-Tag `docker`): Sitzung ohne Terminal mit einem Skript auf der Standardeingabe gegen ein erzeugtes Sandbox-Projekt mit gesetztem Platzhalter für `GIT_TOKEN` und einer Variable eines Agenten: Benutzer und Workspace gleich der Ergebniszeile, beide Variablen gesetzt; eine auf dem Host nicht gesetzte Variable fehlt; *Wächter:* dasselbe `docker exec` ohne `-e` trägt sie nicht. Während der Sitzung besteht kein Container mit dem Namenspräfix des Hilfs-Containers; *Wächter:* die Suche findet einen von Hand gestarteten Container mit dem Präfix. Der Wert steht nicht in `Config.Env` des Devcontainers. Ein Container mit `failed` in einer konfigurierten Statusdatei öffnet keine Sitzung (`11`). Ein Skript, das mit `exit 3` endet, ergibt `0`. Alle sichtbar `PASS` im CI-Lauf, nicht übersprungen. *Beleg:* Lauf `integration-docker` mit Job-Log.
- [ ] **Doku:** Handbuch 4.9 (*Im Container arbeiten* mit `u-boot devcontainer shell`; `docker exec` für einen Befehl mit Exit-Code), 4.11 (*Zugangsdaten*) und jede Aufzählung der Devcontainer-Befehle (Fundstellen in §3), beide READMEs (Befehlstabelle, *Prerequisites*, Link-Liste), `CHANGELOG.md` `[Unreleased]` unter *Added*, `docs/maintainer/cli-json-contract.md` §6.11 (Wert `shell` von `subcommand`, Sentinel-Tabelle). Handbuch und READMEs ohne Anforderungs-, ADR- und Spec-Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); der Terminal-Pfad steht nur mit Test als zugesagt, sonst als Grenze ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4). *Beleg:* Beleg-Tabelle (Satz → Test) in der Übergabe.
- [ ] `make gates` grün; `make test-docker` grün als CI-Lauf auf dem Push-Tip ([`AGENTS.md`](../../../../AGENTS.md) §3.1); `depguard` unverändert ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 5; `make verify-depguard` nur, wenn eine Regel berührt wird).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei Start und zurück auf „Nichts in Arbeit.“ bei der Closure; ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; widerlegt ein Test eine Zeile der Zusage, steht der Befund in §7 und der Slice geht an den Architect zurück.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei oder Stelle, die der Planner nicht gelesen hat. Vor der Umsetzung sucht der Implementer im ganzen Baum nach `devcontainer up`, `devcontainer down`, `up|down`, `up`/`down` und nach den Werten von `subcommand` und trägt jede weitere Fundstelle hier ein, auch die unveränderte ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driving/devcontainer.go` | update | Request und Methode für `shell` (Muster `DevcontainerUpRequest`, gelesen Zeilen 9 bis 31: `BaseDir`, `ProgressSink`, `OnWarning`); Vertrag von `DevcontainerUseCase` (Zeilen 94 bis 114) um `shell`; ein Sentinel für die gescheiterte Probe ist nicht nötig, wenn `driven.ErrDevcontainerRuntime` durchgereicht wird (Wahl des Implementers, in beiden Klassifikatoren gepinnt) |
| `internal/hexagon/application/devcontainer_start.go` | update | `Shell`: Startphase wie `Up` ohne `Recreate` (gemeinsame private Funktion oder Aufruf, Wahl des Implementers), bei `ErrDevcontainerStartScriptFailed` keine Sitzung; dann `launcher.Shell` mit `ContainerID`, `RemoteUser`, `RemoteWorkspaceFolder` der Ergebniszeile; neue reine Funktion der Namensliste neben `localEnvNames` (gelesen, Zeilen 37 bis 56; Fälle der Fitness Function: Kommentar, abweichender Name, Verweis in `mounts`, nicht gesetzte Variable); Ausweisung mit „for the duration of the start“ statt „for the duration of the call“ (`socketRelaxationMessage`, Zeilen 31 bis 35; der Kommentar Zeilen 25 bis 30 nennt die Tests je Teilaussage und bekommt den Test der Variante); Doc-Kommentar von `DevcontainerStartService` („`u-boot devcontainer up` and `down`“, Zeile 59) und `wired` (Zeile 91) für den neuen Befehlsnamen |
| `internal/hexagon/application/devcontainer_start_test.go`, `internal/hexagon/application/export_test.go` | update | Namensliste (Fälle oben), Startphase von `shell` (Ausweisung vor dem Start mit der Variante, `failed` ergibt `11` ohne Aufruf von `Shell`, Validierungsfehler ohne Ausweisung), Weitergabe von Benutzer, Workspace, Kennung, Namen und Terminal an `Shell`; `fakeLauncher` zeichnet die Optionen von `Shell` auf (die Methode kommt aus dem Fundament); Brücke für die reine Funktion wie `LocalEnvNamesForTest` (prüfen) |
| `internal/adapter/driving/cli/devcontainer.go` | update | `newDevcontainerShellCommand` (Muster `newDevcontainerUpCommand`, gelesen Zeilen 85 bis 148); Eintrag in `newDevcontainerCommand` (Zeile 81) und in dessen Langtext (Zeilen 72 bis 75); Texte von `ErrDevcontainerSubcommandRequired` (Zeile 20) und `ErrDevcontainerPreviewNotApplicable` (Zeile 26) um `shell`; Abweisung von `--json` mit eigenem Sentinel oder über den Vorschau-Sentinel (Wahl des Implementers; `isDevcontainerUsageError`, Zeilen 31 bis 33); `runDevcontainerShell` schreibt alles auf `stderr` (anders als `runDevcontainerUp`, Zeile 255) und entscheidet das Terminal an `cmd.InOrStdin()`; Diagnostic-Abbildung für `shell` (Muster `mapDevcontainerUpErrorToDiagnostic`, Zeilen 351 bis 365) |
| `internal/adapter/driving/cli/cli.go` | update | Kommentar des Feldes `devcontainerUseCase` („implements `u-boot devcontainer up|down`“, gelesen Zeile 82); `ExitCode`-Kommentar (Zeilen 266 bis 273) und `isEnvironmentError` (Zeilen 333 bis 340) prüfen, ob sie `shell` nennen müssen; ein neuer Usage-Sentinel geht über `isGroupUsageError` (Zeilen 507 bis 509) in `isUsageError` |
| `internal/adapter/driving/cli/devcontainer_test.go` | update | Muster der bestehenden Tests (gelesen, Zeilen 84 bis 386): bare `devcontainer` nennt `shell`; `shell` mit `--json`, `--dry-run`, `--diff` und Argumenten ergibt `2` (mit `--json` Umschlag mit `subcommand: shell`); nichts auf der Standardausgabe; Ausweisung auch unter `--quiet` auf der Fehlerausgabe; `0` nach einer Sitzung mit Exit-Status ungleich `0`; `11` bei `failed` ohne Sitzung; `12` bei gescheiterter Probe; Exit-Code und Diagnose-Code je Sentinel |
| `internal/adapter/driving/cli/jsontestutil/` | **unverändert**, prüfen | `checkSubcommand` verlangt `subcommand` bei `devcontainer` ([`slice-cr-json-subcommand-pflicht`](../done/slice-cr-json-subcommand-pflicht.md)); der Wert `shell` braucht keine Änderung, wenn die Prüfung nur auf Vorhandensein prüft |
| `cmd/uboot/main.go`, `cmd/uboot/main_test.go` | update | Kommentare „`devcontainer down` stops the container through it“ (gelesen, Zeilen 94 bis 97) und „`devcontainer up|down`“ (Zeilen 170 bis 173) nennen `shell`; Wiring-Test ohne `docker` im `PATH` ergibt `11` nach Muster `TestRun_DevcontainerUp_NoDockerBinaryExits11` (Zeile 350) |
| `internal/e2e/devcontainer_start_docker_test.go` oder eine neue Datei im Paket | update / neu | Docker-Tests der DoD; Hilfen `dcGenerateProject` (Zeile 191, Sandbox-Profil), `dcRun` (Zeile 646, gebautes Binary mit Umgebung und Standardeingabe prüfen), `dcInspect` (Zeile 92, `Config.Env`), `dcRegisterCleanup` (Zeile 165); das Namenspräfix des Hilfs-Containers ist `helperContainerPrefix` = `u-boot-devcontainer-up-` (`internal/adapter/driven/docker/devcontainer_launcher.go`, Zeile 29) |
| `docs/user/benutzerhandbuch.md` | update | 4.9 *Im Container arbeiten* (gelesen, Zeilen 818 bis 833) und der Befehlsblock (Zeilen 746 bis 750); 4.11 *Zugangsdaten* (Zeilen 1907 bis 1912). Aufzählungen der Devcontainer-Befehle (gelesen): Voraussetzungen (Zeilen 57, 58, 62 und 63, 69), *Container-Image* (Zeilen 162 bis 164), Option `--quiet` (Zeile 323), Profil-Tabelle (Zeile 1123), Rechte-Tabelle (Zeile 2156) und Hinweis (Zeile 2163), §7 *Exit-Codes* (Zeilen 2179 und 2180), FAQ (Zeile 2384), Befehlsübersicht (Zeilen 2478 und 2479), Optionen je Befehl (Zeilen 2498 und 2499), JSON-Ausgabe `subcommand` (Zeile 2548; die Bedingung „nur bei“ der ganzen Zeile ist ungepinnter Bestand, V-1 aus [`slice-cr-json-subcommand-pflicht`](../done/slice-cr-json-subcommand-pflicht.md) §7 — berührt der Diff die Zeile, nennt der Plan die Pin-Lage), `data` (Zeilen 2590 und 2591: `shell` hat kein `data`), *Besonderheiten* (Zeilen 2600 bis 2606). Eine Umbenennung der Überschrift *Den Devcontainer ohne Editor starten und stoppen* (Zeile 744) bricht die Anker in beiden READMEs |
| `README.md`, `README.de.md` | update | Befehlstabelle (gelesen: `README.md` Zeilen 224 und 225, `README.de.md` Zeilen 227 und 228), *Prerequisites* (`README.md` Zeilen 232 und 235, `README.de.md` Zeilen 235 und 238), Link-Liste (`README.md` Zeile 264, `README.de.md` Zeile 269); Status-Block (Zeilen 11 bzw. 12) beschreibt `v0.13.0` und bleibt bis zum Release-Cut |
| `CHANGELOG.md` | update | `[Unreleased]` *Added* (gelesen am 2026-10-21: Zeile 12, Abschnitt leer; der Slice zum Urteil nach dem Zustand trägt dort vorher seinen Eintrag unter *Changed*) |
| `docs/maintainer/cli-json-contract.md` | update | §6.11 (gelesen, Zeilen 1244 bis 1302): Überschrift und Text um `shell` (nur der Umschlag des Exit-Codes `2`), Spalte `shell` der Sentinel-Tabelle; Absatz zu den Reject-Formen (Zeilen 390 bis 393) prüfen |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start, „Nichts in Arbeit.“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Zeile dieses Slice |
| diese Datei, die Welle-Datei, die übrigen Slices der Welle | update | Plan-Abgleich und Closure; Verweis-Nachzug durch `make slice-mv` |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence je angefallener Beobachtung |
| `spec/` <!-- d-check:ignore (Pfad-Zitat) --> | **unverändert** | CR davor, Spezifikation danach |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md) und [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md) liegen in `done/`; der Befund des Fundaments zum Terminal-Pfad liegt vor (Test oder benannte Grenze); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Startphase lässt sich nicht mit `Up` teilen, ohne dessen Vertrag zu ändern; dann zuerst ein Umbau der Application als eigener Slice.
- `in-progress` → `open` (blockiert): ein Docker-Test widerlegt eine Zeile der Zusage (etwa ein Hilfs-Container während der Sitzung, ein Wert in `Config.Env`); Befund an den Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; die Docker-Tests am Befehl sind im CI-Lauf von `make test-docker` sichtbar bestanden und nicht übersprungen; Review- und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Terminal-Pfad ohne Test.** Die Docker-Tests laufen ohne Terminal; die Nutzung mit Terminal ist der Regelfall. — **Ausgang:** weiter offen; folgt dem Befund des Fundaments: mit Test zugesagt, ohne Test in Handbuch und Spezifikation als Grenze, und der Lastenheft-Wortlaut „interaktiv“ geht als Frage an den Projektinhaber ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 10×).
- **Eine Aufzählung der Befehle bleibt ohne `shell`** (Handbuch, READMEs, Fehlertexte, Langtexte, Kommentare). — **Ausgang:** weiter offen; Suche vor der Umsetzung und nach der letzten Fix-Runde, Fundstellen in §3 ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), 3×; [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md), 2×).
- **Ein neuer Sentinel steht nur in einem Klassifikator.** — **Ausgang:** weiter offen; Pin-Test je Sentinel in `ExitCode` und in der Diagnostic-Abbildung von `shell` (Regel der Sicht-Spec, *Fehlermodelle*).
- **Die Hilfe von `shell` trägt die Zeile „Exit codes (…)“ mit Anforderungs-Kennung** wie `up` und `down`; findet der Validator dieser Welle sie wieder, erreicht [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) (2×, Belege je Welle) die Schwelle. — **Ausgang:** weiter offen; Produktfrage beim Projektinhaber (offene Entscheidung (j) der Vorwelle, Frage in der Welle-Datei §5).
- **Strg-C und Signalende verhalten sich anders als in [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4 angenommen** (etwa ein Client, der nach dem Ende von `u-boot` weiterläuft). — **Ausgang:** weiter offen; Befunde aus dem Fundament, am Befehl wiederholt, soweit ein Test ohne Terminal sie trägt.
- **Das e2e-Paket überschreitet das Zeitlimit** (weitere Tests mit gestartetem Sandbox-Container). — **Ausgang:** weiter offen; CI-Lauf ist Beleg, das Limit ist Gate-Gegenstand ([`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md), 1×).
- **Ein Docker-Test hinterlässt Container oder ist nur mit lokalem Cache grün.** — **Ausgang:** weiter offen ([`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md), 2×; [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), 1×).
- **Eine Bedingung im Handbuch ohne Test** (etwa „eine auf dem Host nicht gesetzte Variable fehlt in der Shell“, „kein Socket während der Sitzung“). — **Ausgang:** weiter offen; Beleg-Tabelle in der Übergabe ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert).
- **Ein Kommentar beschreibt Abwesendes oder einen alten Namen** (Kommentare mit `up|down`, „two subcommands“). — **Ausgang:** weiter offen; Selbstcheck nach [`AGENTS.md`](../../../../AGENTS.md) §3.9 ([`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md), verkörpert).
- **Eine Aussage im Plan über den Code stimmt nicht** (die Zeilen „prüfen“). — **Ausgang:** weiter offen ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6. Die drei Paarungen prüft die Closure der Welle.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT, driving), `internal/adapter/driving/cli` (CLI), `internal/e2e` (E2E), `cmd/uboot` (Wiring, keine Sub-Area) und `docs/` (DOK), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) alle Greenfield. Die Rückstufung der CLI-Sub-Area droht, wenn eine tragende Adapter-Konvention nur im Code entsteht: Die Regel „alles, was `u-boot` selbst ausgibt, auf der Standardfehlerausgabe“ und die Terminal-Erkennung gehören deshalb in die Spezifikation (Spezifikations-Slice), nicht nur in den Code.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-21; Zähler = Dateien in `evidence/`): [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (3, offen) und [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (2) — viele Aufzählungen, Risiko in §6; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (10, offen) — Terminal-Pfad; [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) (2) und [`BEO-CLI/hilfetext-backtick-wort-als-typname`](../observations/BEO-CLI/hilfetext-backtick-wort-als-typname/observation.md) (1; die Hilfe des neuen Befehls zeigt `--quiet up` wie jede andere, die Stelle liegt im Usage-Text von `--quiet`, nicht in diesem Slice); [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) (5, verkörpert); [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (1), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1); [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2) — nur Platzhalter-Werte, keine Rohdaten; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (35) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (14), beide verkörpert; [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (7, offen). Evidenz-/Diskrepanz-Risiko: mittel (Terminal-Pfad, Signale).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
