# Slice v1: Fundament von `u-boot devcontainer shell` — Port-Methode `DevcontainerLauncher.Shell`, `docker exec` des Hosts mit Benutzer, Workspace und Umgebung, Probe der Shell, Docker-Tests mit Wächter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-shell-und-zweites-up`](../welle-devcontainer-shell-und-zweites-up.md) (Release `v0.14.0`). Dritter von sechs Slices; nach [`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md), vor [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md). Inhaltlich hängt dieser Slice nicht am Urteil nach dem Zustand (er berührt nur Port und Adapter); die Reihenfolge folgt der Welle, und der Befehl braucht beide.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Absatz zu `shell` in der Fassung des CR), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Übergabe zur Laufzeit per Umgebungsvariable, Wert in keiner Datei), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Klasse `12` der gescheiterten Probe), [`LH-NFA-PORT-002`](../../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) (keine neue Abhängigkeit auf dem Host); [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (Entscheidungen 1 *Sitzungsphase*, *Umgebung*, *Shell*, *Kein Socket während der Sitzung*; 4 Signale und Terminal mit den Prüfaufträgen; 5 Port; 6 Zusage-Menge; `Accepted`), [`ADR-0028`](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) (Entscheidung 4, Lebensdauer des Sockets: für die Startphase unverändert).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (*Umgebung und Zugangsdaten*, *Hilfs-Container*) und `architecture.md` §2.4 und §2.6 (Port `DevcontainerLauncher`, Adapter). Gelesen, nicht geändert; die Nachführung ist [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md). Der Verweis zeigt aufwärts.

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

**Ziel:** Der Driven-Port `DevcontainerLauncher` hat eine Methode `Shell`, und der Docker-Adapter öffnet damit eine Sitzung im laufenden Devcontainer über `docker exec` des Hosts: Benutzer, Workspace und Kennung kommen vom Aufrufer, `-t` nur mit Terminal, `-e <NAME>` ohne Wert nur für gesetzte Variablen, davor eine nicht interaktive Probe, die `bash` oder `sh` wählt und bei fehlendem Benutzer, Workspace oder Shell scheitert ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidungen 1 und 5). Kein Befehl, keine nutzersichtbare Änderung.

**Warum ein eigenes Fundament** (Schnitt des Planners nach dem Muster [`slice-v1-devcontainer-start-fundament`](../done/slice-v1-devcontainer-start-fundament.md)): Der ganze Befehl hätte fünf Liefer-Punkte (Port und Adapter, Application und Befehl, Docker-Tests am Adapter, Docker-Tests am Befehl, Dokumente), die Ziel-Form erlaubt drei. Der Schnitt trennt die Teile mit dem höchsten Risiko — die Semantik von `docker exec -e`, das Ende der Shell mit dem Ende des Clients und den Terminal-Pfad, die Prüfaufträge von Entscheidung 4 — vom Befehl, der auf ihnen aufsetzt. Widerlegt ein Docker-Test hier eine Annahme der ADR, geht das an den Architect, bevor ein Befehl darauf gebaut ist.

**Gelesener Stand (2026-10-21).** `internal/hexagon/port/driven/devcontainer_launcher.go`: Interface `DevcontainerLauncher` mit `Up` und `Stop` (Zeilen 104 bis 128); jede Methode weist ein ungültiges Argument als `ErrDevcontainerRuntime` ohne Docker-Aufruf ab und prüft dann die Umgebung (Zeilen 109 bis 116); `ValidateDevcontainerEnvName` (Zeilen 33 bis 48). `internal/adapter/driven/docker/devcontainer_launcher.go`: `preflight` (Zeilen 109 bis 126: Binary, Daemon, lokale Engine), `passedEnvNames` (Zeilen 265 bis 278: prüft die Namen, behält die gesetzten über `lookupEnv`), `containerIDPattern` in `Stop` (Zeile 149); kein `docker exec` im Adapter. Der Inspector führt nie ein Programm im Container aus (`internal/hexagon/port/driven/devcontainer_inspector.go`, Doc-Kommentar Zeilen 33 bis 41); das gilt weiter, die Probe gehört dem Launcher. Ein Test-Double des Ports steht in `internal/hexagon/application/devcontainer_start_test.go` (`fakeLauncher`, Zeilen 20 bis 48) und muss die neue Methode tragen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Application, Driving-Port und Befehl `devcontainer shell`, die Namensliste aus `devcontainer.json`** — [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md); die reine Funktion der Namensliste ist Application ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1, *Umgebung*), dieser Slice bekommt die Namen als Argument.
- **Die Startphase** (Hilfs-Container, Socket, Ergebniszeile) — unverändert `Up`; `Shell` startet keinen Hilfs-Container und bindet kein Socket ein.
- **Dokumente für Nutzer** (Handbuch, READMEs, `CHANGELOG.md`, JSON-Kontrakt) — keine nutzersichtbare Änderung; sie kommen mit dem Befehl.
- **`devcontainer exec` der CLI und `read-configuration`** — verworfen ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Optionen C und D).
- **Nachbildung von `userEnvProbe` und von `remoteEnv` in anderer Form** — nicht zugesagt ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 6).
- **Eine neue Modulabhängigkeit** (etwa eine Pseudo-Terminal-Bibliothek für den Test des Terminal-Pfads) — nur nach Entscheidung des Projektinhabers (Welle §5, offene Frage); bis dahin mit Mitteln des Test-Images oder als benannte Grenze.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Port und Adapter:** `DevcontainerLauncher.Shell` mit Optionen (Kennung, Benutzer, Workspace, Namensliste, Standardein- und -ausgabe, Fehlerausgabe, Terminal ja oder nein); der Doc-Vertrag sagt, dass die Methode ohne Docker-Aufruf ungültige Argumente als `ErrDevcontainerRuntime` abweist, dieselbe Vorab-Prüfung wie `Up` und `Stop` hat (`ErrDockerUnavailable`, `ErrContainerEngineNotLocal`), bei gescheiterter Probe `ErrDevcontainerRuntime` liefert und nach Beginn der Sitzung keinen Fehler aus dem Exit-Status der Shell macht ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 3). Der Adapter ruft `docker exec -i [-t] -u <Benutzer> -w <Workspace> [-e <NAME> …] <Kennung> <Shell>` (Login-Shell, `bash`, sonst `sh`) ohne Hilfs-Container, ohne Socket, ohne `--privileged` und ohne `-v`; im Terminal-Modus bricht er den Client nicht über den Kontext ab ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4). *Beleg:* Unit-Test der Argumentliste (Golden-Liste je Fall) und der Fehlerklassen.
- [ ] **Docker-Tests am Adapter mit Wächter** (Build-Tag `docker`): Sitzung ohne Terminal gegen einen laufenden Container mit einem Skript auf der Standardeingabe (`id -un`, `pwd`, `printenv` zweier Variablen): Benutzer und Workspace wie übergeben, eine gesetzte Variable mit ihrem Wert, eine auf dem Host nicht gesetzte fehlt; *Wächter:* dasselbe `docker exec` ohne `-e` trägt die gesetzte Variable nicht. Der Wert steht auf keiner Kommandozeile (Suche wie `cmdlineHolders`); *Wächter:* die Suche findet einen absichtlich in ein Argument gesetzten Wert. Probe: fehlender Benutzer, fehlender Workspace und ein Image ohne `bash` (Wahl `sh`) bzw. ohne beide Shells ergeben die erwarteten Fälle ohne Sitzung. **Prüfaufträge von Entscheidung 4:** die Shell im Container endet mit dem Ende des Clients (ohne Terminal; mit Terminal, wenn der Test-Container ein Pseudo-Terminal bereitstellen kann), Verhalten von Strg-C ohne Terminal; jedes Ergebnis steht als Befund in §7, ein nicht getesteter Pfad als benannte Grenze. *Beleg:* Lauf `integration-docker` mit Job-Log (sichtbar `PASS`, nicht übersprungen).
- [ ] **Doku für Entwickler:** Doc-Kommentare von Port und Adapter beschreiben den Ist-Stand ([`AGENTS.md`](../../../../AGENTS.md) §3.7); die Paket-READMEs `internal/adapter/driven/README.md` und `internal/hexagon/port/driven/README.md` nennen weder Launcher noch Devcontainer (gelesen, Suche ohne Treffer) und bleiben unverändert; kein Nutzer-Dokument.
- [ ] `make gates` grün; `make test-docker` grün als CI-Lauf auf dem Push-Tip ([`AGENTS.md`](../../../../AGENTS.md) §3.1); `make lint` ohne neue Ausnahme, `depguard` unverändert ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Fitness Function; `make verify-depguard` nur, wenn eine Regel berührt wird).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei Start und zurück auf „Nichts in Arbeit.“ bei der Closure; ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; widerlegt ein Test eine Annahme von [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (etwa die Gleichwertigkeit von `docker exec -e <NAME>` und `remoteEnv` für die erzeugte Konfiguration), steht der Befund in §7 und der Slice geht an den Architect zurück.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei oder Stelle, die der Planner nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/port/driven/devcontainer_launcher.go` | update | neuer Typ für die Optionen von `Shell` und die Methode im Interface `DevcontainerLauncher` (gelesen, Zeilen 104 bis 128); Doc-Kommentar des Interfaces („`Up` runs … `Stop` stops …“, „`Up`: project folder, then variable names; `Stop`: container ID“) um `Shell` und deren Argumentprüfung; ob Benutzer und Workspace eine eigene Prüffunktion neben `ValidateDevcontainerEnvName` (Zeilen 33 bis 48) bekommen, entscheidet der Implementer |
| `internal/hexagon/port/driven/devcontainer_launcher_test.go` | update (prüfen) | gelesen: `TestValidateDevcontainerProjectFolder` und `TestValidateDevcontainerEnvName`; nur bei einer neuen Prüffunktion ein Fall dazu |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | update | `Shell` neben `Up` und `Stop`: Argumentprüfung (Kennung über `containerIDPattern` wie in `Stop`, Zeile 149; Namen über `passedEnvNames`, Zeilen 265 bis 278), `preflight` (Zeilen 109 bis 126), Probe (nicht interaktives `docker exec -u -w <Kennung>`, Wahl `bash` oder `sh`), dann die Sitzung mit den Strömen der Optionen; Argumentliste als reine Funktion neben `runArgs` (Zeilen 280 bis 302), damit der Unit-Test sie pinnt. Im Terminal-Modus kein Abbruch über `cmd.Cancel` (anders als `runUp`, Zeile 210) |
| `internal/adapter/driven/docker/devcontainer_launcher_test.go` | update | Golden-Liste der Argumente je Fall nach Muster `TestRunArgs` (Zeile 135): `-i`, `-t` nur mit Terminal, `-u` und `-w`, `-e <NAME>` ohne Wert nur für gesetzte Namen, kein `--privileged`, kein Socket, kein `-v`; Fehlerklassen nach Muster `TestLauncher_Preflight` (Zeile 265), `TestLauncher_Stop` (Zeile 657) und `TestLauncher_Up_RejectsUnsafeInput` (Zeile 635); gescheiterte Probe ergibt `ErrDevcontainerRuntime`, ein Exit-Status der Shell ungleich `0` nach Beginn der Sitzung keinen Fehler |
| `internal/adapter/driven/docker/devcontainer_launcher_docker_test.go` | update | Docker-Tests der DoD; Hilfen `requireDocker` (Zeile 41), `mustDocker` (Zeile 69), `removeContainers` (Zeile 79), `cmdlineHolders` (Zeile 472) und `scanDuring` (Zeile 517) für die Suche nach dem Wert auf Kommandozeilen während der Sitzung; Testcontainer mit einem Image ohne `bash` für die Wahl `sh` (Image prüfen) |
| `internal/adapter/driven/docker/export_test.go` | update (prüfen) | Brücken für die Argumentliste, falls der Test sie aus dem Paket `docker_test` aufruft |
| `internal/hexagon/application/devcontainer_start_test.go` | update | `fakeLauncher` (gelesen, Zeilen 20 bis 48) bekommt die Methode `Shell`, damit das Paket mit dem erweiterten Interface baut; kein neuer Test hier |
| `Dockerfile` (Stufe `test-docker-tools`) | **unverändert**, prüfen | die Stufe baut auf `deps` (`golang`-Image) und installiert `docker-ce-cli` und `docker-compose-plugin` (gelesen, Zeilen 97 bis 111); ob `script` (Debian-Paket `bsdutils`) für einen Test mit Pseudo-Terminal im Image liegt, ist nicht gelesen; eine Änderung der Stufe ist ein Plan-Eintrag |
| `cmd/uboot/main.go` | **unverändert** | `docker.NewLauncher()` (gelesen, Zeile 98) verdrahtet den Adapter schon; der Kommentar Zeilen 94 bis 97 („builds the helper image and runs the Dev Container CLI …; `devcontainer down` stops the container through it“) bekommt den Satz zu `shell` erst mit dem Befehl |
| `internal/adapter/driven/README.md`, `internal/hexagon/port/driven/README.md` | **unverändert** | gelesen: keine Nennung von Launcher oder Devcontainer |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start, „Nichts in Arbeit.“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Zeile dieses Slice |
| diese Datei, die Welle-Datei, die übrigen Slices der Welle | update | Plan-Abgleich und Closure; Verweis-Nachzug durch `make slice-mv` |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence je angefallener Beobachtung |
| `spec/`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | keine nutzersichtbare Änderung; Dokumente mit dem Befehl, Spezifikation danach |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-cr-devcontainer-shell-und-zweites-up`](../in-progress/slice-cr-devcontainer-shell-und-zweites-up.md) liegt in `done/` (das Lastenheft sagt `shell` zu); [`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md) liegt in `done/` (Reihenfolge der Welle; ist er blockiert, darf dieser Slice vorgezogen werden, weil er ihn inhaltlich nicht braucht — die Umstellung steht dann in §5 der Welle-Datei); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Terminal-Pfad verlangt eine Änderung des Test-Images oder eine neue Abhängigkeit; dann wird er ein eigener Slice nach der Entscheidung des Projektinhabers, und dieser schließt mit der benannten Grenze.
- `in-progress` → `open` (blockiert): ein Docker-Test widerlegt eine Annahme der ADR (die Shell lebt nach dem Ende des Clients weiter, `docker exec -e <NAME>` setzt eine auf dem Host fehlende Variable leer statt gar nicht); dann Befund an den Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; die Docker-Tests am Adapter sind im CI-Lauf sichtbar bestanden und nicht übersprungen; jeder Prüfauftrag von [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4 hat in §7 einen Befund oder eine benannte Grenze; Review- und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Terminal-Pfad ist in der CI nicht testbar** (der Runner gibt kein Terminal; ein Pseudo-Terminal braucht ein Werkzeug im Test-Image oder eine Abhängigkeit). — **Ausgang:** weiter offen; Befund in §7, und ohne Test steht der Pfad in Handbuch und Spezifikation als „nicht zugesagt“ ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 10×). Die Frage an den Projektinhaber steht in der Welle-Datei §5.
- **Die Shell im Container lebt nach dem Ende des Clients weiter** (nicht belegt, Prüfauftrag). — **Ausgang:** weiter offen; Docker-Test, sonst Rückführung (§4).
- **Der Wert einer Variable erscheint auf einer Kommandozeile** (etwa weil die Probe die Umgebung mitgibt). — **Ausgang:** weiter offen; die Probe bekommt keine `-e`-Liste, der Docker-Test sucht den Wert während der Sitzung mit Wächter.
- **Ein Docker-Test hinterlässt Container oder ist nur mit lokalem Cache grün.** — **Ausgang:** weiter offen ([`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md), 2×; [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), 1×); Aufräumen über `removeContainers`, Beleg ist der CI-Lauf.
- **Das Paket `internal/e2e` und das Adapter-Paket laufen in einem `go test ./...` gleichzeitig** (der Fundament-Test der Vorwelle entfernte ein Image, das der andere nutzte; [`slice-v1-cli-devcontainer-up-down`](../done/slice-v1-cli-devcontainer-up-down.md) §3). — **Ausgang:** weiter offen; die Tests dieses Slice starten eigene Container und entfernen kein geteiltes Image.
- **Eine Aussage im Plan über den Code stimmt nicht** (die Zeilen „prüfen“). — **Ausgang:** weiter offen ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Host-Werkzeug außerhalb Docker** beim Untersuchen des Verhaltens (etwa ein Host-`docker exec` als Beleg). — **Ausgang:** weiter offen; Belege nur aus `make`-Läufen ([`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md), 7×).

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

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Befunde zu den Prüfaufträgen von Entscheidung 4 · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6. Die drei Paarungen prüft die Closure der Welle.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/port` (PORT, driven), `internal/adapter/driven` (ADP) und im Test-Double `internal/hexagon/application` (APP), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) alle Greenfield. Der Port-Pin `var _ driven.DevcontainerLauncher = (*Launcher)(nil)` (`devcontainer_launcher.go`, Zeile 85) macht eine Abweichung von Port und Adapter zum Build-Fehler.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-21; Zähler = Dateien in `evidence/`): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (10, offen) — Terminal-Pfad, Risiko in §6; [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (1) — die Aussage zu `devcontainer exec` in [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) stammt aus einer Messvariante; dieser Slice misst die ausgelieferte Form am Adapter; [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2) und [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) — Risiken in §6; [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (2) — die Tests setzen nur Platzhalter-Werte und schreiben keine Rohdaten; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (35, verkörpert); [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (7, offen). Evidenz-/Diskrepanz-Risiko: mittel (Prüfaufträge ohne Vorbefund).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
