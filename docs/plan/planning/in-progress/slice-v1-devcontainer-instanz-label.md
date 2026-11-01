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

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, `down` und `shell`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Workspace im benannten Volume je Instanz), [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (Weg A, zweites Id-Label, Verhalten ohne Flag; Proposed, der Projektinhaber setzt `Accepted` vor dem Start), [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) (Volume-Namen über `${devcontainerId}`; unverändert, von der neuen ADR geschärft), [ADR-0028](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [ADR-0031](../../adr/0031-devcontainer-shell-ohne-editor.md) (aktive ADRs zu `up` und `shell`; die ADR [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) ergänzt sie).

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

**Entschieden (Projektinhaber, 2026-10-31): Weg A.** Entschieden sind außerdem das Flag `--instance <name>` und die Namensregel (2026-10-31). Entschieden ist auch `up`, `down` und `shell` ohne `--instance` (Fassung (i): nur der Container ohne Instanz; bei nur Instanz-Containern Exit-Code `10`, geprüft vor dem Aufruf der CLI). Offen ist der Wortlaut des Lastenheft-CR (Freigabe durch den Projektinhaber, vor dem Code); er liegt im CR ([`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md)), die Messung im Mess-Slice (siehe Kopf und §6).

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

- [ ] **Annahmen des Mess-Slice geprüft** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Ein Docker-Test fährt `u-boot devcontainer up --instance` für N gleichzeitige Instanzen desselben Projektordners und pinnt die Messwerte aus [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) §7 (eigene Containerkennungen, verschiedene Volume-Namen für die dort gemessenen Familien, Inhaltsprobe: Datei in einer Instanz in keiner anderen, Verhalten von `up` ohne Flag neben einer Instanz); N ist die Zahl, die [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) zusagt (⟦M9⟧ dort; Anforderung des Projektinhabers: Isolation mit drei und vier Instanzen), **mindestens drei**, nicht nur zwei; alle Volume-Namen und Containerkennungen sind paarweise verschieden und die Inhaltsprobe läuft über alle Paare. Eine Abweichung vom Messwert ist ein Befund an den Mess-Slice, kein Anpassen des Tests. Der Test ist ein Docker-Test (`make test-docker`) und trägt das Host-Limit des Mess-Slice (Speicher, CPU, Platz); er bricht mit Befund ab und lässt die Stufe nicht weg. Die Messung selbst (Skript, Protokoll) gehört nicht mehr in diesen Slice.
- [ ] **Der gewählte Weg** (A, B oder C): Bei A akzeptieren `up`, `down` und `shell` den Instanznamen, `FindContainer` findet den Container der Instanz, ohne Angabe ändert sich nichts (Bestandstest), ein ungültiger Name (Entscheidung des Projektinhabers 2026-10-31: `[a-z0-9-]`, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32; leer und `=` ungültig) endet mit Exit-Code `10` nach [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); nach der Entscheidung des Projektinhabers (2026-10-31, [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md), Entscheidung 4) prüft `u-boot` vor dem Aufruf der CLI, ob es im Projekt nur Container von Instanzen gibt, und endet bei `up`, `down` und `shell` ohne `--instance` dann mit Exit-Code `10` und einem Hinweis auf `--instance` (je ein Test), statt die CLI eine Instanz übernehmen zu lassen; Tests je Bedingung und je Mutation, ein Docker-Test mit N Instanzen (Zahl des CR, mindestens drei). Bei B oder C steht der Ablauf mit Datum und Messstand im Handbuch.
- [ ] **Spezifikation, Handbuch und `CHANGELOG.md` nach dem Code** (das Lastenheft liefert [`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md) vorher) mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] `make gates` grün (im Modus der Umgebung).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die öffentliche Schnittstelle (CLI-Flag, Handbuch, README, `CHANGELOG.md`), soweit ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

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
| `docs/plan/adr/0038-devcontainer-instanz-zweites-id-label.md`, `docs/plan/adr/README.md` | neu (geschrieben vom Architect, Status `Proposed`) | die ADR für Weg A und die Zusage-Menge; sie schärft den Satz „zwei Instanzen aus **demselben** Ordner bleiben eine Instanz“ (Variante D) in [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md); [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) selbst wird nicht umgeschrieben (`AGENTS.md` §3.5). Der Slice startet erst mit `Accepted` |
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
- **Neu:** `internal/hexagon/application/devcontainer_instance_test.go`, `internal/adapter/driving/cli/devcontainer_instance_test.go`, `internal/e2e/devcontainer_instance_docker_test.go`.
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
| `doctor`, `generate`, `init`: nur Instanzen = nicht ermittelbar | `TestDevcontainerInstance_OnlyInstancesIsNotDeterminableForTheStatusReading` |
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

- **Mehrdeutige Suche:** Tragen zwei Container dasselbe Label `devcontainer.local_folder`, kann der Treffer der falsche sein ([`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md)); ein Instanz-Label löst das nur, wenn **alle** drei Aufrufer von `FindContainer` es mitgeben. — **Ausgang:** bei Closure.
- **Volumes und Image je Instanz** sind im Entwurf angenommen; gemessen werden sie im Mess-Slice, hier prüft der Docker-Test sie. **Anforderung des Projektinhabers (2026-10-31):** Zwei Instanzen teilen sich kein Volume; jede läuft isoliert (eigenes Workspace-Volume, eigener Docker-Speicher). Wechselt die `devcontainerId` nicht mit dem Id-Label, ist das kein hinnehmbarer Zustand, sondern eine Lücke, die der Slice schließt (etwa indem `u-boot` den Instanznamen in die Volume-Namen der erzeugten Konfiguration aufnimmt) oder bei der `--instance` mit Exit-Code `10` abbricht, solange die Isolation nicht gilt. — **Ausgang:** bei Closure (die Messung liegt im Vorgänger, der Wortlaut im CR).
- **Bestandscontainer ohne Instanz-Label:** Ein Container ohne Flag trägt `devcontainer.local_folder` und `devcontainer.config_file` (gemessen im Mess-Slice, Protokoll M1), die Instanzen tragen kein `devcontainer.config_file`; ohne Flag muss der Bestandscontainer weiter gefunden werden, auch wenn eine Instanz neuer ist (gemessen: der erste Treffer von `docker ps -a` ist dann die Instanz). Nach der Entscheidung zu `up` ohne `--instance` schließt die Suche Container mit Instanz-Label aus (ob `docker ps` ein fehlendes Label filtern kann, ist nicht gemessen) und `up` ohne Flag endet bei nur Instanz-Containern mit Exit-Code `10`, geprüft vor dem Aufruf der CLI; Test-Aufträge: Bestandscontainer plus neuere Instanz (ohne Flag trifft `up`, `down`, `shell` den Bestandscontainer), nur Instanzen (`up`, `down` und `shell` ohne Flag: je ein Test, Exit-Code `10` mit Hinweis auf `--instance`, keine CLI gestartet, soweit der Befehl sie sonst aufriefe), ungültiger Name (Exit-Code `10`); dazu die Fälle, die [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) als nicht gemessen führt: wiederholtes `up --instance` ohne `--recreate`, `down --instance` mit unbekanntem Namen, Bestandscontainer plus neuere Instanz im Sandbox-Profil und für `down` und `shell`. — **Ausgang:** bei Closure.
- **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): der Docker-Test prüft die Zahl, die der CR zusagt (mindestens drei, höchstens vier), und kein Test oder Text sagt mehr zu; zwei Instanzen (Bestandstest) sind kein Beleg, und eine Zahl auf einem Host ist kein Beleg für mehr Instanzen oder für andere Hosts. — **Ausgang:** bei Closure.
- **Ressourcen:** Jede Instanz läuft mit `--privileged`, hat einen verschachtelten Docker mit eigenem Speicher und einen kalten Go-Cache im Workspace-Volume; Platz und Zeit sind nicht gemessen. — **Ausgang:** bei Closure.
- **Ausweisung des Sockets bei „nur Instanzen“:** `up` und `shell` weisen die Socket-Lockerung des Hilfs-Containers aus, bevor die Prüfung „nur Instanzen“ mit Exit-Code `10` endet; es startet dann kein Hilfs-Container (Abweichung 2 in §3). Der Text sagt „für die Dauer des Aufrufs“ und bleibt wahr, beschreibt aber einen Hilfs-Container, der nicht entsteht. Alternative: die Abfrage vor die Ausweisung ziehen und den Bestandstest `TestDevcontainerUp_AlreadyRunning` ändern (Entscheidung beim Review). — **Ausgang:** bei Closure.
- **Name der Instanz:** Zeichensatz und Länge sind entschieden (Projektinhaber, 2026-10-31: `[a-z0-9-]`, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32, ungültig → Exit-Code `10`; leer und `=` ausgeschlossen, gemessen: leer lehnt die CLI ab, bei `=` hängt `devcontainer up`); der Wortlaut steht im CR, der Code setzt die freigegebene Regel um und pinnt sie mit Tests (je Rand, Länge 32 und 33, `=`, leer, Großbuchstabe). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

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
