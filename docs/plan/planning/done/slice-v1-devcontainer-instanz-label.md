# Slice V1: Mehrere Devcontainer je Projekt — ein Instanzname an `up`, `down` und `shell` über das Id-Label der Dev-Container-CLI

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Release `v0.17.0` (geplant, Roadmap *Nächste Wellen*; die Welle-Datei entsteht bei der Eröffnung). Die Closure-Bedingung der Welle (alle Slices in `done/`, `make gates` und Replay grün, Release-Cut mit Tag) ist von der DoD dieses Slice verschieden. Anlass: Der Projektinhaber will mehrere Agenten getrennt arbeiten lassen und statt Git-Worktrees mehrere Devcontainer desselben Projekts starten (Gespräch am 2026-10-31).

**Vorgänger (Messung → Lastenheft-CR → Code, Vorbild [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)):** dritter von drei Slices. Die Messung der Volume-Namen, des Id-Labels und des Verhaltens ohne Flag liegt in [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md), der Lastenheft-CR (Wortlaut, Version, Freigabe des Projektinhabers) in [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md). Dieser Slice trägt den Code samt Tests, Docker-Test, Spezifikation, Handbuch und `CHANGELOG.md`.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, `down` und `shell`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Workspace im benannten Volume je Instanz), [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (Weg A, zweites Id-Label, Verhalten ohne Flag; `Accepted`), [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) (Volume-Namen über `${devcontainerId}`; unverändert, von der neuen ADR geschärft), [ADR-0028](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [ADR-0031](../../adr/0031-devcontainer-shell-ohne-editor.md) (aktive ADRs zu `up` und `shell`; die ADR [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) ergänzt sie).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz *Aufruf*: `u-boot devcontainer up [--recreate]`, `down`, `shell`).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Projekt kann mehrere Devcontainer gleichzeitig haben, jeden mit eigenem Workspace-Volume, und `u-boot devcontainer up`, `down` und `shell` wählen den Container über einen Instanznamen; ohne Angabe verhalten sich alle drei wie heute.

**Befund** (gelesen am 2026-10-31, Fundstellen je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Ein Container je Projektordner.** `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`, Zeilen 284 bis 303) ruft die Dev-Container-CLI als `devcontainer up --workspace-folder <Ordner>`, bei `--recreate` mit `--remove-existing-container`, und gibt kein `--id-label` mit. `FindContainer` (`internal/adapter/driven/docker/inspector.go`, Filter in Zeile 86, Konstante `devcontainer.local_folder` in Zeile 21) sucht den Container allein über dieses Label und liefert bei mehreren Treffern den zuletzt angelegten (Port `DevcontainerInspector`, `internal/hexagon/port/driven/devcontainer_inspector.go`, Zeilen 44 bis 49).
- **Die CLI kann mehr.** Die Dev-Container-CLI 0.89.0 im Hilfs-Image von `u-boot` kennt `devcontainer up --id-label name=value`: Die Labels werden am Container gesetzt und dienen der Suche nach einem bestehenden Container; ohne Angabe leitet die CLI sie aus dem Workspace-Ordner ab (`devcontainer up --help`, gemessen am 2026-10-31). Ob die Namen der Volumes (`${devcontainerId}` in `.devcontainer/devcontainer.json`) mit dem Id-Label mitwechseln, ist aus der Spezifikation der CLI angenommen und **nicht gemessen**; die Messung ist Gegenstand von [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md), dieser Slice prüft ihre Annahmen im Docker-Test.
- **`u-boot` reicht nichts durch.** `u-boot devcontainer up` kennt nur `--recreate` neben den globalen Flags (`u-boot devcontainer up --help`; `internal/adapter/driving/cli/devcontainer.go`, Flag in Zeile 164, `down` ab Zeile 173, `shell` ab Zeile 209).
- **Die Suche ist an drei Stellen verdrahtet:** `internal/hexagon/application/devcontainer_start.go` (Zeilen 297 und 429) und `devcontainer_fallback.go` (Zeile 155) rufen `FindContainer(ctx, <Ordner>)`. Das Beobachtungs-Register führt den Mangel bereits: [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (ein Beleg, unter der Schwelle; Trigger: „Messung mit zwei Containern am selben Label“ und „eine Änderung am Port `DevcontainerInspector`“ — beides tritt mit diesem Slice ein).
- **Zusage heute:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) beschreibt `up`, `down` und `shell` für den Devcontainer „des Projekts“, ohne Instanz; zugesagt auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine.

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, `--instance <name>` an `up`, `down` und `shell`:** `u-boot` übergibt der CLI zwei Id-Labels (`devcontainer.local_folder=<Ordner>` und ein eigenes Label mit dem Instanznamen, Arbeitsname `u-boot.instance`), und `FindContainer` filtert auf beide. Ohne Flag bleibt es beim heutigen Aufruf. Neue öffentliche Flags: Lastenheft-CR ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt keine Instanz) in [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md).
- **B, Git-Worktrees als Projektordner, kein Produktcode:** Jeder Worktree ist ein eigener Projektordner und damit ein eigener Container (eigenes `devcontainer.local_folder`); das Handbuch beschreibt den Weg. Der Container klont weiter aus `origin`, der Worktree liefert nur die Identität. Nicht gemessen.
- **C, nichts im Produkt.** Das Problem bleibt.

**Entschieden (Projektinhaber, 2026-10-31): Weg A.** Entschieden sind außerdem das Flag `--instance <name>` und die Namensregel (2026-10-31). Entschieden ist auch `up`, `down` und `shell` ohne `--instance` (Fassung (i): nur der Container ohne Instanz; bei nur Instanz-Containern Exit-Code `10`, geprüft vor dem Aufruf der CLI). Der Wortlaut des Lastenheft-CR liegt im CR ([`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md)), die Messung im Mess-Slice (beide in `done/`; siehe Kopf und §6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Welchen Zweig eine Instanz bearbeitet.** Der Clone folgt `devcontainer.sandbox.repository` oder `origin` und checkt den Standardzweig aus; die Wahl des Zweigs ist ein anderer Vorgang mit eigener Entscheidung (ein Feld je Instanz in `u-boot.yaml` oder ein Startparameter).
- **Agenten starten oder verteilen.** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt als Nicht-Ziel, dass `u-boot` keinen Agenten startet; ein Ablauf über mehrere Instanzen bleibt beim Aufrufer.
- **Der Compose-Stack je Instanz.** Der Befehl startet laut [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nicht den Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)).
- **Last und Grenzen bei vielen Instanzen.** Gemessen wird mit zwei, drei und vier Instanzen (Mess-Slice); ein Urteil über zehn wäre ein Schluss über die Messung hinaus ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Andere Hosts als Linux `amd64`.** Sie sind für `up` und `shell` heute nicht zugesagt; der Slice ändert das nicht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Annahmen des Mess-Slice geprüft** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Ein Docker-Test fährt `u-boot devcontainer up --instance` für N gleichzeitige Instanzen desselben Projektordners und pinnt die Messwerte aus [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) §7 (eigene Containerkennungen, verschiedene Volume-Namen für die dort gemessenen Familien, Inhaltsprobe: Datei in einer Instanz in keiner anderen, Verhalten von `up` ohne Flag neben einer Instanz); N ist die Zahl, die [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) zusagt (⟦M9⟧ dort; Anforderung des Projektinhabers: Isolation mit drei und vier Instanzen), **mindestens drei**, nicht nur zwei; alle Volume-Namen und Containerkennungen sind paarweise verschieden und die Inhaltsprobe läuft über alle Paare. Eine Abweichung vom Messwert ist ein Befund an den Mess-Slice, kein Anpassen des Tests. Der Test ist ein Docker-Test (`make test-docker`) und trägt das Host-Limit des Mess-Slice (Speicher, CPU, Platz); er bricht mit Befund ab und lässt die Stufe nicht weg. Die Messung selbst (Skript, Protokoll) gehört nicht mehr in diesen Slice.
- [x] **Der gewählte Weg** (A, B oder C): Bei A akzeptieren `up`, `down` und `shell` den Instanznamen, `FindContainer` findet den Container der Instanz, ohne Angabe ändert sich nichts (Bestandstest), ein ungültiger Name (Entscheidung des Projektinhabers 2026-10-31: `[a-z0-9-]`, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32; leer und `=` ungültig) endet mit Exit-Code `10` nach [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); nach der Entscheidung des Projektinhabers (2026-10-31, [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md), Entscheidung 4) prüft `u-boot` vor dem Aufruf der CLI, ob es im Projekt nur Container von Instanzen gibt, und endet bei `up`, `down` und `shell` ohne `--instance` dann mit Exit-Code `10` und einem Hinweis auf `--instance` (je ein Test), statt die CLI eine Instanz übernehmen zu lassen; Tests je Bedingung und je Mutation, ein Docker-Test mit N Instanzen (Zahl des CR, mindestens drei). Bei B oder C steht der Ablauf mit Datum und Messstand im Handbuch.
- [x] **Spezifikation, Handbuch und `CHANGELOG.md` nach dem Code** (das Lastenheft liefert [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) vorher) mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [x] `make gates` grün (im Modus der Umgebung).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für die öffentliche Schnittstelle (CLI-Flag, Handbuch, README, `CHANGELOG.md`), soweit ein öffentlicher Vertrag berührt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; die Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Welche Zeilen der Weg berührt, entscheidet der Projektinhaber; die Tabelle nennt Weg A.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/devcontainer-instanz/`, `docs/messungen/slice-v1-devcontainer-instanz-messung/` | entfällt hier | Skript und Protokoll liegen im Mess-Slice [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) |
| `internal/adapter/driving/cli/devcontainer.go` | update | `up` (Flag `--recreate`, Zeile 164), `down` (ab Zeile 173) und `shell` (ab Zeile 209) erhalten `--instance <name>`; die Hilfetexte mit den Aufrufbeispielen (ab Zeile 84) folgen |
| `internal/hexagon/port/driven/devcontainer_launcher.go` | update | `DevcontainerUpOptions` (ab Zeile 113: `ProjectFolder`, `EnvNames`, `Recreate`, `ReadOnlyProject`, `Progress`) bekommt ein Feld für die Instanz |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | update | `runArgs` (Zeilen 284 bis 303) hängt mit Instanz zwei `--id-label` an den Aufruf von `devcontainer up` |
| `internal/hexagon/port/driven/devcontainer_inspector.go` | update | `FindContainer(ctx, localFolder)` (Zeile 49) wird um die Instanz erweitert oder bekommt eine zweite Methode |
| `internal/adapter/driven/docker/inspector.go` | update | Filter `label=devcontainer.local_folder=<Ordner>` (Zeile 86) wird um das Instanz-Label ergänzt |
| `internal/hexagon/application/devcontainer_start.go`, `devcontainer_fallback.go` | update | Aufrufe von `FindContainer` (Zeilen 297, 429 und 155); `down` und `shell` suchen den Container der Instanz. Die Tests daneben (`devcontainer_start_test.go`, `devcontainer_shell_test.go`) folgen |
| `internal/e2e/devcontainer_start_docker_test.go` und Nachbarn | prüfen | Docker-Test mit mindestens drei Instanzen (die Zahl des CR), der den Befehl fährt; nicht gelesen. Gelesen: `TestE2E_SandboxDevcontainer_TwoInstancesIsolated` (`internal/e2e/sandbox_devcontainer_docker_test.go`, ab Zeile 557) bleibt als Bestand unverändert, ist aber **kein Beleg** für die Isolation durch den Befehl; er ersetzt `${devcontainerId}` von Hand durch `idaaaa` und `idbbbb` und fährt weder den Befehl noch die CLI; er belegt die Isolation bei verschiedenen Werten, nicht den Wechsel des Werts |
| `spec/spezifikation.md` | update | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (ab Zeile 829, Absatz *Aufruf*): Flag, Fehlerklassen, Verhalten ohne Flag |
| `spec/lastenheft.md` | unverändert (CR-Slice davor) | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 507) ändert [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) mit Freigabe des Projektinhabers; dieser Slice liest den Stand nach dem CR |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | update (prüfen) | Abschnitt zu `up`, `down` und `shell` und die Release-Hinweise; nicht gelesen |
| `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md`, `docs/plan/adr/README.md` | neu (geschrieben vom Architect, Status `Accepted`) | die ADR für Weg A und die Zusage-Menge; sie schärft den Satz „zwei Instanzen aus **demselben** Ordner bleiben eine Instanz“ (Variante D) in [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md); [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) selbst wird nicht umgeschrieben (`AGENTS.md` §3.5) |
| `internal/hexagon/port/driving/devcontainer.go` | update | `DevcontainerUpRequest`, `DevcontainerShellRequest`, `DevcontainerDownRequest` (Zeilen 11, 56, 89) bekommen die Instanz; zwei neue Sentinels neben `ErrDevcontainerProjectFolderInvalid` (Zeile 114): ungültiger Name und „nur Instanzen“ |
| `internal/hexagon/port/driven/devcontainer_launcher.go` | update (zusätzlich) | die Prüffunktion der Namensregel neben `ValidateDevcontainerProjectFolder` (Zeile 28); der Adapter `docker` importiert `domain` nicht, darum liegt sie im Port-Paket |
| `internal/hexagon/application/noop.go` | update | `noopDevcontainerInspector.FindContainer` (Zeile 51) folgt der neuen Signatur |
| `internal/adapter/driving/cli/devcontainer.go`, `cli.go` | update (zusätzlich) | `mapDevcontainerDownErrorToDiagnostic` kennt heute keinen Sentinel der Klasse `10` und braucht den Zweig; `isDevcontainerValidationError` (Zeile 548) nimmt die zwei Sentinels, `isValidationError` (`cli.go` Zeile 386) bindet sie ein; `mapDevcontainerShellErrorToDiagnostic` bleibt (nur Klasse `2`, `--json` wird vor dem Use-Case abgewiesen). Tests: `devcontainer_test.go`, `devcontainer_shell_test.go` |
| Test-Doubles und Tests der Signatur | update | `fakeInspector.FindContainer` (`internal/hexagon/application/devcontainer_fallback_test.go`, Zeile 71) und `devcontainer_start_test.go` (Zeile 634 liest `insp.findFolders`); `devcontainer_shell_test.go` (Zeile 179 liest `insp.findFolders`) nur, falls das Feld seine Form ändert; `export_test.go` (Paket `application`) enthält kein `FindContainer` und bleibt außen vor (gelesen); `internal/adapter/driven/docker/inspector_test.go` (Argv-Test, Tests der Treffer ohne Instanz); die Aufrufer in `devcontainer_launcher_docker_test.go`, `internal/e2e/devcontainer_start_docker_test.go` und `internal/e2e/sandbox_status_docker_test.go` (Signatur) |
| `docs/maintainer/cli-json-contract.md` (§6.11), `docs/maintainer/zusagen-und-messstaende.md` | update (prüfen) | Exit-Code `10` für `down` und das Flag; die Zusage-Menge der Instanzen; beide nur nach Lesen der Stellen |
| `harness/replay.md` | unverändert | gelesen: nennt `devcontainer up`, `down` und `shell` nicht (nur `generate devcontainer`); `up`, `down` und `shell` erzeugen keine Dateien |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

**Abgleich mit `git diff --stat` nach der Umsetzung** (Implementer, Stand vor dem Review; nach jeder späteren Änderung zu wiederholen, `AGENTS.md` §3.10). Die Tabelle oben bleibt der Plan; hier steht, was der Diff wirklich berührt und wo er vom Plan abweicht.

- **Berührt wie geplant:** `internal/adapter/driving/cli/devcontainer.go` (Flag `--instance` je Unterbefehl, Hilfetexte, `emptyInstanceError`, `isDevcontainerInstanceError`, Zweig in `mapDevcontainerDownErrorToDiagnostic`; `isDevcontainerValidationError` nimmt die zwei Sentinels), `internal/hexagon/port/driving/devcontainer.go` (Feld `Instance` in den drei Anfragen, `ErrDevcontainerInstanceInvalid`, `ErrDevcontainerOnlyInstances`), `internal/hexagon/port/driven/devcontainer_launcher.go` (`DevcontainerUpOptions.Instance`, `ValidateDevcontainerInstance`), `internal/hexagon/port/driven/devcontainer_inspector.go`, `internal/adapter/driven/docker/devcontainer_launcher.go` (`runArgs`, Prüfung in `Up`), `internal/adapter/driven/docker/inspector.go`, `internal/hexagon/application/devcontainer_start.go`, `devcontainer_fallback.go`, `noop.go`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/cli-json-contract.md`, `docs/maintainer/zusagen-und-messstaende.md`.
- **Tests und Doubles wie geplant:** `fakeInspector` (`devcontainer_fallback_test.go`, Feld `findInstances` neu), `devcontainer_start_test.go`, `inspector_test.go`, `devcontainer_launcher_docker_test.go`, `internal/e2e/devcontainer_start_docker_test.go` (Signatur und die Felder der Grenze in `dcInspectDoc`), `internal/e2e/sandbox_status_docker_test.go`.
- **Berührt, im Plan nicht oder anders geführt:** `internal/hexagon/application/devcontainer_shell_test.go` (Plan: „nur, falls das Feld seine Form ändert“; die Erwartung „`shell` fragt den Inspector nicht“ gilt nicht mehr, siehe Abweichung 3), `internal/adapter/driving/cli/devcontainer_test.go` und `devcontainer_shell_test.go` (Sentinel-Tabellen), `internal/adapter/driven/docker/devcontainer_launcher_test.go` (`TestRunArgs`, Ablehnung), `internal/hexagon/port/driven/devcontainer_launcher_test.go` (Namensregel), `docs/plan/planning/observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md` (Symbol `runningContainerID` heißt jetzt `containerBeforeStart`; `AGENTS.md` §3.9 (b)). [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) nennt `runningContainerID` weiter; sie ist `Accepted` und wird nicht umgeschrieben (`AGENTS.md` §3.5).
- **Nach der Umsetzung hinzugekommen (Stand nach Review und Verifikation):** `.d-check.yml` (Commit `d888dd8c`, MR-005: `docs-check` wertet die gitignorierten Arbeitsbäume unter `.claude/worktrees` nicht; gelesen im Diff: eine Ausnahme für Kopien des Repos, keine Lockerung für Projektdateien; im Plan nicht geführt, Verification V-2), die zwei Reports `docs/reviews/2026-11-01-slice-v1-devcontainer-instanz-label.md` und `…-verify.md`, `devcontainer_instance_test.go` und `devcontainer_fallback_test.go` (Fix-Runde `81a31dff`, `requireNoInstanceAsked`) sowie `docs/plan/planning/open/README.md` und `open/slice-v1-test-docker-in-ci.md` (Commit `3d633863`, nicht Teil dieses Slice).
- **Neu:** `internal/hexagon/application/devcontainer_instance_test.go`, `internal/adapter/driving/cli/devcontainer_instance_test.go`, `internal/e2e/devcontainer_instance_docker_test.go`.
- **Plan-Zeile `cli.go` entfällt:** `isValidationError` (`internal/adapter/driving/cli/cli.go`) bindet `isDevcontainerValidationError` schon ein (gelesen: `return isServiceValidationError(err) || isDevcontainerValidationError(err)`), die zwei Sentinels kommen darüber an; `cli.go` bleibt unverändert.
- **Unverändert wie geplant:** `spec/lastenheft.md`, `harness/replay.md`, `internal/hexagon/application/export_test.go` (gelesen: kein `FindContainer`), `docs/plan/adr/0038-…` und `docs/plan/adr/README.md` (liegen aus dem Architect-Commit vor), `docs/plan/planning/open/README.md` (Zeile liegt seit dem Start vor).

**Abweichungen vom Plan** (Entscheidungen des Implementers innerhalb der ADR; keine widerspricht einer Entscheidung):

1. `FindContainer(ctx, localFolder, instance)` erhält einen dritten Parameter statt einer zweiten Methode; „nur Instanzen“ ist der Sentinel `driven.ErrOnlyInstanceContainers` (die ADR lässt die Form offen). Die Aufrufer `doctor`, `generate devcontainer` und `init` laufen über `fallbackReader` mit Instanz `""` und lesen den Sentinel als „nicht ermittelbar“.
2. Die Prüfung „nur Instanzen“ steht **nach** der Ausweisung des Sockets und vor dem Launcher, weil `TestDevcontainerUp_AlreadyRunning` die Ausweisung vor jeder Abfrage der Engine festhält (die Abfrage soll nicht von einer langsamen Engine aufgehalten werden). Folge: Bei „nur Instanzen“ erscheint die Ausweisung des Hilfs-Containers, obwohl keiner startet (Risiko unten). Der ungültige Name steht davor.
3. `up --recreate` ohne `--instance` und `shell` fragen den Inspector jetzt einmal vor dem Launcher (für die Prüfung); `alreadyRunning` bleibt `false`. Der Teilfall `recreate_skips_the_query` heißt `recreate_reports_no_running_state`.
4. `--instance` mit leerem Wert prüft der CLI-Adapter selbst (`Changed("instance")`), weil die Anfrage `""` als „keine Instanz“ trägt.

**Beleg-Tabelle der Bedingungen in Spec und Handbuch** (Bedingung → Test; `AGENTS.md` §3.10). Jede Zeile nennt die Mutation, die den Test rot färbt, im Bericht des Implementers.

| Bedingung ([`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Absatz *Instanzen*; Handbuch 4.9) | Test |
|---|---|
| Namensregel (Zeichen, Rand, 1 bis 32, leer, `=`, Großbuchstabe) | `TestValidateDevcontainerInstance` |
| ungültiger Name: Exit 10 vor Engine, Ausweisung, Hilfs-Container | `TestDevcontainerInstance_InvalidNameEndsBeforeAnyPort`, `TestDevcontainerInstance_DisclosureComesBeforeTheOnlyInstancesCheck`, `TestLauncher_Up_RejectsUnsafeInput`, `TestInspector_FindContainer_RejectsBadInstance`, Docker: `TestDevcontainerInstance_Docker_InvalidNamesEndWithExit10` |
| leeres `--instance` ist ungültig, nicht „keine Instanz“ | `TestExecute_Devcontainer_EmptyInstanceIsInvalid` |
| zwei `--id-label` nur mit Instanz, keine ohne | `TestRunArgs` (Fälle `defaults`, `instance`, `instance with recreate`) |
| Suche mit zwei Filtern / ohne Instanz Label je Zeile | `TestInspector_FindContainer_InstanceArgv`, `TestInspector_FindContainer_Instances` |
| nur Instanzen: Exit 10, Hinweis, kein Launcher, kein `docker stop` (auch `up --recreate`) | `TestDevcontainerInstance_OnlyInstancesEndsBeforeTheLauncher`, Docker: `…FourInstancesAreIsolated` (Binary, kein Hilfs-Container) |
| Ausweisung vor der Prüfung „nur Instanzen“ | `TestDevcontainerInstance_DisclosureComesBeforeTheOnlyInstancesCheck` |
| kein Container: `up`/`shell` starten, `down` Exit 0 | `TestDevcontainerDown/no container`, `TestDevcontainerInstance_DownWithoutAMatchingContainer`, `TestDevcontainerUp_AlreadyRunning/no_container_before` |
| Bestandscontainer neben neuerer Instanz: ohne Flag der Bestand, mit Flag die Instanz | `TestDevcontainerInstance_ExistingContainerAndNewerInstance`; Docker: `…FourInstancesAreIsolated` (Sandbox, `up`, `down`, `shell`), `…StandardProfileSharesTheProjectFolder` (`up`) |
| zweites `up --instance` ohne `--recreate` trifft denselben Container | `…ExistingContainerAndNewerInstance/up_with_the_instance`; Docker: `…FourInstancesAreIsolated` |
| `up --instance --recreate` ersetzt nur diesen Container, Volumes bleiben | Docker: `…FourInstancesAreIsolated` |
| `down --instance` stoppt nur diese Instanz; unbekannter Name: Exit 0 | `…DownWithoutAMatchingContainer`, `…ExistingContainerAndNewerInstance/down_with_the_instance`; Docker: `…FourInstancesAreIsolated` |
| Statusdateien aus dem Container der Instanz | `TestDevcontainerInstance_StatusReadingNamesTheInstance` |
| `doctor`, `generate`, `init`: nur Instanzen = nicht ermittelbar | `TestDevcontainerInstance_OnlyInstancesIsNotDeterminableForTheStatusReading` (Reader); die Aufrufer fragen mit Instanz `""`: `requireNoInstanceAsked` in `TestDoctor_SandboxFallback_OnlyConfiguredCapabilities`, `TestGenerateDevcontainer_ExistingContainerFallbackWarning`, `TestInit_ExistingContainerFallbackWarning` (rot bei Instanz `"x"` im `fallbackReader` von `doctor_sandbox.go`, `generate.go`, `initproject.go`) |
| Klasse 10 und Diagnose-Code für beide Sentinels, `up`, `down`, `shell` | `TestExecute_Devcontainer_SentinelsPerSubcommand`, `TestExecute_DevcontainerShell_SentinelsAndExitCodes` |
| Flag erreicht die Anfrage; Hilfetexte | `TestExecute_Devcontainer_InstanceFlagReachesTheRequest`, `TestExecute_Devcontainer_HelpNamesInstance` |
| eigene Container und Volumes je Instanz, Inhaltsprobe über alle Paare, Labels, Grenze zurückgelesen (vier Instanzen, Sandbox) | Docker: `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` |
| Standardprofil: Instanzen teilen den Projektordner | Docker: `TestDevcontainerInstance_Docker_StandardProfileSharesTheProjectFolder` |
| Container ohne Instanz bleibt ohne Label von `u-boot` | `TestRunArgs/defaults`; Docker: `TestDevcontainerUp_Docker_S1_S9_StandardProfile` |
| Race zwischen Prüfung und CLI „nicht abgesichert“ | keine Behauptung über Verhalten, kein Test |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; beide Vorgänger liegen in `done/`: [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) (Messwerte und Urteil zur Kernfrage) und [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) (Lastenheft mit freigegebenem Wortlaut); Weg A ist entschieden (Projektinhaber, 2026-10-31); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Docker-Test widerlegt einen Messwert des Mess-Slice (etwa: `${devcontainerId}` wechselt nicht mit dem Id-Label und die Volumes würden geteilt): dann gilt die Anforderung des Projektinhabers (Isolation, §6) und der Befund geht an den Mess-Slice und den CR; ändert sich dafür das Schema der Mounts, ist das ein eigener Vorgang vor diesem Slice.
- `in-progress` → `open` (blockiert): Der Docker-Test kann die Zahl paralleler Instanzen, die der CR zusagt, nicht herstellen (Host trägt sie nicht): Befund an den Mess-Slice und den CR, die Zahl wird nicht im Test gesenkt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; die Zusage in Lastenheft, Spezifikation und Handbuch nennt nur das Gemessene (die Zahl der Instanzen aus dem Mess-Slice, höchstens vier, Linux `amd64`); Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Mehrdeutige Suche:** Tragen zwei Container dasselbe Label `devcontainer.local_folder`, kann der Treffer der falsche sein ([`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md)); ein Instanz-Label löst das nur, wenn **alle** drei Aufrufer von `FindContainer` es mitgeben. — **Ausgang:** entfallen für Container mit Instanz-Label: der Instanzname erreicht `FindContainer` aus `up`, `down` und `shell`; `doctor`, `generate` und `init` fragen mit Instanz `""` und lesen „nur Instanzen“ als nicht ermittelbar (Review Punkt 9, `requireNoInstanceAsked` in drei Tests, Verification §6). Der Fall einer zweiten, von `u-boot` nicht erzeugten Konfigurationsdatei am selben Label bleibt wie im Register geführt.
- **Volumes und Image je Instanz** sind im Entwurf angenommen; gemessen werden sie im Mess-Slice, hier prüft der Docker-Test sie. **Anforderung des Projektinhabers (2026-10-31):** Zwei Instanzen teilen sich kein Volume; jede läuft isoliert (eigenes Workspace-Volume, eigener Docker-Speicher). Wechselt die `devcontainerId` nicht mit dem Id-Label, ist das kein hinnehmbarer Zustand, sondern eine Lücke, die der Slice schließt (etwa indem `u-boot` den Instanznamen in die Volume-Namen der erzeugten Konfiguration aufnimmt) oder bei der `--instance` mit Exit-Code `10` abbricht, solange die Isolation nicht gilt. — **Ausgang:** entfallen: `${devcontainerId}` wechselt mit dem Id-Label (Mess-Slice), die Mounts bleiben unverändert; `TestDevcontainerInstance_Docker_FourInstancesAreIsolated` bestand im Lauf des Implementers und im Lauf des Verifiers (vier Instanzen, Inhaltsprobe über alle Paare; Verification §1).
- **Bestandscontainer ohne Instanz-Label:** Ein Container ohne Flag trägt `devcontainer.local_folder` und `devcontainer.config_file` (gemessen im Mess-Slice, Protokoll M1), die Instanzen tragen kein `devcontainer.config_file`; ohne Flag muss der Bestandscontainer weiter gefunden werden, auch wenn eine Instanz neuer ist (gemessen: der erste Treffer von `docker ps -a` ist dann die Instanz). Nach der Entscheidung zu `up` ohne `--instance` schließt die Suche Container mit Instanz-Label aus (ob `docker ps` ein fehlendes Label filtern kann, ist nicht gemessen) und `up` ohne Flag endet bei nur Instanz-Containern mit Exit-Code `10`, geprüft vor dem Aufruf der CLI; Test-Aufträge: Bestandscontainer plus neuere Instanz (ohne Flag trifft `up`, `down`, `shell` den Bestandscontainer), nur Instanzen (`up`, `down` und `shell` ohne Flag: je ein Test, Exit-Code `10` mit Hinweis auf `--instance`, keine CLI gestartet, soweit der Befehl sie sonst aufriefe), ungültiger Name (Exit-Code `10`); dazu die Fälle, die [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) als nicht gemessen führt: wiederholtes `up --instance` ohne `--recreate`, `down --instance` mit unbekanntem Namen, Bestandscontainer plus neuere Instanz im Sandbox-Profil und für `down` und `shell`. — **Ausgang:** entfallen: jeder Test-Auftrag trägt einen Test (Beleg-Tabelle in §3); das Format `{{.Label "u-boot.instance"}}` gegen die echte Engine belegt der Docker-Test (Verification §1 und §6).
- **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): der Docker-Test prüft die Zahl, die der CR zusagt (mindestens drei, höchstens vier), und kein Test oder Text sagt mehr zu; zwei Instanzen (Bestandstest) sind kein Beleg, und eine Zahl auf einem Host ist kein Beleg für mehr Instanzen oder für andere Hosts. — **Ausgang:** entfallen: der Docker-Test fährt vier Instanzen (`dcInstanceCount = 4`), die Zusage nennt höchstens vier, Linux `amd64`, Sandbox-Profil; Review und Verifier fanden keine Aussage darüber hinaus. Grenze des Tests: der freie Platz des Docker-Roots ist im Test-Container nicht lesbar und bleibt ungeprüft (Verification V-3).
- **Ressourcen:** Jede Instanz läuft mit `--privileged`, hat einen verschachtelten Docker mit eigenem Speicher und einen kalten Go-Cache im Workspace-Volume; Platz und Zeit sind nicht gemessen. — **Ausgang:** weiter offen: gemessen ist nur der Zustand direkt nach dem Start (Mess-Slice), der Platz im Docker-Test nicht lesbar (Verification V-3). Anker: [`docs/maintainer/zusagen-und-messstaende.md`](../../../maintainer/zusagen-und-messstaende.md), Abschnitt *Mehrere Devcontainer je Projekt*, führt „das Verhalten unter Last“ als nicht zugesagt; kein Folge-Slice.
- **Ausweisung des Sockets bei „nur Instanzen“:** `up` und `shell` weisen die Socket-Lockerung des Hilfs-Containers aus, bevor die Prüfung „nur Instanzen“ mit Exit-Code `10` endet; es startet dann kein Hilfs-Container (Abweichung 2 in §3). Der Text sagt „für die Dauer des Aufrufs“ und bleibt wahr, beschreibt aber einen Hilfs-Container, der nicht entsteht. — **Ausgang:** akzeptiert, Reihenfolge bleibt (Review F-1). Gründe: [ADR-0030](../../adr/0030-devcontainer-up-failed-nach-zustand.md) führt für Exit 11 ebenfalls „Ausweisung, danach Fehler“ vor dem Start eines Hilfs-Containers; Umordnen bricht den Schutz von `TestDevcontainerUp_AlreadyRunning` (eine langsame Engine darf die Ausweisung nicht aufhalten) und die in [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und `cli-json-contract.md` gepinnte Reihenfolge. Eine wahrheitsgetreue Ausweisung braucht eine Folge-Entscheidung zu [ADR-0030](../../adr/0030-devcontainer-up-failed-nach-zustand.md).
- **Name der Instanz:** Zeichensatz und Länge sind entschieden (Projektinhaber, 2026-10-31: `[a-z0-9-]`, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32, ungültig → Exit-Code `10`; leer und `=` ausgeschlossen, gemessen: leer lehnt die CLI ab, bei `=` hängt `devcontainer up`); der Wortlaut steht im CR, der Code setzt die freigegebene Regel um und pinnt sie mit Tests (je Rand, Länge 32 und 33, `=`, leer, Großbuchstabe). — **Ausgang:** entfallen: `TestValidateDevcontainerInstance` pinnt Rand, Länge 32 und 33, `=`, leer und Großbuchstabe (Verification §1).
- **Docker-Nachweis des Gesamtlaufs (bei der Closure ergänzt):** Der Lauf des Implementers von `make test-docker` war grün (1540 s, die drei neuen Instanz-Tests bestanden). Im Lauf des Verifiers bestanden die drei neuen Tests ebenfalls; drei Bestandstests waren rot (`TestE2E_LHFAUP005_LogsTail`, `TestE2E_LHFAUP005_LogsFollow`, `TestE2E_SandboxDevcontainer_TwoInstancesIsolated` mit Namenskonflikt `uboot-e2e-inst-idaaaa`; Verification V-1). Vermutete Ursache: ein zeitgleicher zweiter Lauf auf derselben Engine, **nicht bewiesen**. Auf Entscheidung des Projektinhabers wird `make test-docker` nicht lokal wiederholt. — **Ausgang:** weiter offen: der Lauf der CI auf `main` (`integration.yml`, harter Job, kein Required Check) entscheidet; Anker: [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) (Entwurf). Der Gesamtlauf wird nicht als erfüllt geführt.

## 7. Closure-Notiz

Reports: [Review](../../../reviews/2026-11-01-slice-v1-devcontainer-instanz-label.md) (0 HIGH, 0 MEDIUM, 4 LOW, 2 INFO; Fix-Runde `81a31dff` behob F-1 bis F-4), [Verification](../../../reviews/2026-11-01-slice-v1-devcontainer-instanz-label-verify.md) (0 HIGH, 1 MEDIUM V-1, 1 LOW V-2, 2 INFO; `make gates` Exit 0 auf `c361601a`, `coverage-gate` 92,20 %).

- **Was hat funktioniert:** Weg A ist wie in [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) entschieden umgesetzt (Verification §6: alle neun Entscheidungen konform). Der neue Docker-Test mit vier Instanzen (`TestDevcontainerInstance_Docker_FourInstancesAreIsolated`) und seine zwei Nachbarn bestanden im Lauf des Implementers und im Lauf des Verifiers; damit sind die Annahmen des Mess-Slice (eigene Container und Volumes je Instanz, Inhaltsprobe über alle Paare) und das Format `{{.Label "u-boot.instance"}}` gegen die echte Engine belegt. Die Beleg-Tabelle in §3 ordnete jede Bedingung in Spec und Handbuch einem Test zu; der Review fand genau eine Lücke (F-4), die die Fix-Runde mit `requireNoInstanceAsked` schloss.
- **Was ging anders als geplant:** (1) Die vier Abweichungen des Implementers in §3 (dritter Parameter statt zweiter Methode, „nur Instanzen“ nach der Ausweisung, `shell` und `up --recreate` fragen den Inspector einmal, leeres `--instance` im CLI-Adapter). (2) Die Ausweisung der Socket-Lockerung erscheint bei „nur Instanzen“, obwohl kein Hilfs-Container startet; akzeptiert (Risiko in §6, Review F-1). (3) Der Plan-Abgleich nannte `.d-check.yml` nicht (Verification V-2, im Abgleich in §3 ergänzt). (4) **Der Docker-Nachweis des Gesamtlaufs bleibt offen:** im Lauf des Verifiers waren drei Bestandstests rot (§6, letztes Risiko; Verification V-1); der Slice führt ihn nicht als erfüllt, die CI auf `main` entscheidet.
- **Steering-Loop-Eintrag:** Zwei Lehren, beide gezählt und nicht verkörpert (kein `liegt in`). (a) Eine Aussage über die Aufrufer einer geänderten Schnittstelle („`doctor`, `generate`, `init` betrachten nur den Container ohne Instanz“) in Spec und Handbuch brauchte einen Test an jedem Aufrufer, nicht nur am Pfad, über den der Test lief (Review F-4): erneutes Auftreten der Klasse in `AGENTS.md` §3.10 (Bedingung ohne pinnenden Test), Rückfall nach der Verkörperung. (b) Ein Slice mit Docker-Verhalten hat auf einer von mehreren Agenten geteilten Engine keinen belastbaren lokalen Gesamtlauf: drei Bestandstests wurden rot, ohne dass der Diff ihren Code berührte, und die Ursache ist ohne Wiederholung nicht zu trennen. Der Weg dorthin (der CI-Lauf als Nachweis) ist als Entwurf [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) geführt, nicht verkörpert. Dazu ein Spec-Befund ohne Lücke: [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) nennt `runningContainerID`, der Code `containerBeforeStart`; die ADR ist `Accepted` und bleibt (`AGENTS.md` §3.5).
- **Beobachtungs-Register (`../observations/`):** keine neue Beobachtung angelegt; drei Evidence-Dateien `evidence/slice-v1-devcontainer-instanz-label.md` ergänzt, `state.md` je nachgeführt, kein Zähler gesetzt:
  - [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md): Trigger „eine Änderung am Port `DevcontainerInspector`“ eingetreten; der Instanz-Fall ist behoben, der Rest (zweite Konfigurationsdatei) nicht. Stand: **3 Dateien** unter `evidence/`, Schwelle erreicht, Ausgang nicht zugewiesen.
  - [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md): Review F-4 (behoben); Rückfall nach der Verkörperung, **28 Dateien**.
  - [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): Review F-2 und F-3, Verification V-2 (behoben); Rückfall nach der Verkörperung, **64 Dateien**.
  - Nicht fortgeschrieben: die Beobachtungen zu Docker-Tests unter `BEO-E2E/` (die Ursache der drei roten Tests ist unbewiesen, ein Beleg wäre eine Vermutung) und [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Review und Verifier fanden keine Aussage über die Messung hinaus).
- **Folge-Slices:** [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) (Entwurf in `open/`, Anker des offenen Docker-Nachweises). Kein weiterer Slice.
- **Risiken aus §6:** sieben Risiken mit Ausgang: fünf entfallen (Mehrdeutige Suche, Volumes und Image, Bestandscontainer, Zusage weiter als die Messung, Name der Instanz), eines akzeptiert (Ausweisung des Sockets), **zwei weiter offen** (Ressourcen, Anker `docs/maintainer/zusagen-und-messstaende.md`; Docker-Nachweis des Gesamtlaufs, Anker `slice-v1-test-docker-in-ci`).
- **Drei Paarungen:** Anker: kein `liegt in` in diesem Slice, nichts zu prüfen. Folge-Slice: [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md) liegt in `open/`. Register: die drei Evidence-Dateien liegen unter `../observations/` (Ergebnis nach dem `git mv` geprüft, siehe Häkchen in §2).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/adapter/driving/cli` (CLI), `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT), `internal/adapter/driven` (ADP), `internal/e2e` (E2E) sowie `spec/`, `harness/` und `docs/plan/` (DOK); `tools/` (WZG) und `docs/messungen/` berührt der Mess-Slice. Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen. Gelesen: [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (APP, ein Beleg, unter der Schwelle; sein Trigger tritt mit diesem Slice ein, ein zweiter Beleg wäre die Messung mit zwei Containern). Nur nach Titel gesichtet, Inhalt beim Start zu lesen: [`BEO-E2E/test-zaehlt-container-paralleler-pakete`](../observations/BEO-E2E/test-zaehlt-container-paralleler-pakete/observation.md) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Docker-Tests mit mehreren Containern), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
