# Verification-Report: slice-v1-sandbox-egress-neuladen-entwurf — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-sandbox-egress-neuladen-entwurf` (`docs/plan/planning/in-progress/`), Commit-Range `2533ecc..HEAD` (8 Commits, HEAD `606157e`, 35 Dateien), [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) (`Proposed`). Entscheidung des Projektinhabers (2026-10-18): Variante R0, Guardrail-Status quo, kein Neu-Laden-Befehl, keine Folge-Slices. <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-sandbox-egress-neuladen-entwurf.md` (F-1 MEDIUM, F-2..F-5 LOW, F-6..F-10 INFO), Fix-Runde `606157e`, DoD-Haken und Sensor-Behauptung `make gates` Exit 0 aus dem Slice. <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflage.** Die Entscheidung R0 ist in Slice, ADR, Welle, Handbuch und CHANGELOG konsistent umgesetzt, die Handbuch-Aussagen tragen gegen die Rohdaten, kein Produktcode, keine Spec, kein Lastenheft, keine Accepted-ADR ist berührt. Kein Review-Finding ist offen, wo es behebbar war. Der Slice ist **nicht schließbar**, bis der Projektinhaber [`ADR-0025`](../plan/adr/0025-egress-liste-im-laufenden-container.md) auf `Accepted` setzt und die offene DNS-Frage (Slice §6) beantwortet. Auflagen siehe Abschnitt 8. <!-- d-check:ignore (Kennung als Zitat) -->

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** (Exit-Code des `make`-Prozesses selbst erfasst, nicht der von `tail`) | `coverage-gate: OK — coverage 91.20% meets threshold 90%`; Lauf im Verifier-Kontext nach HEAD `606157e` |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ |
| `make doc-immutable RANGE=2533ecc..HEAD` | pass, Exit 0 | Range aufgelöst, 8 Commits; „417 Datei(en) geprüft, 0 Befund(e)“; im Diff nur die neue `ADR-0025` und eine Index-Zeile, keine Accepted-ADR verändert | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make test-docker`, `make verify-depguard` | nicht gelaufen | kein Produktcode, keine Importregeln, kein Docker-Pfad des Produkts: der Diff berührt weder `internal/` noch `spec/` noch `Makefile` (aus `git diff --name-status` ablesbar) |
| Neuer Lauf der Messskripte | nicht gelaufen | braucht Colima mit `NET_ADMIN`, Netz und das Dev-Container-CLI-Image; Docker-only (`AGENTS.md` §3.1); Rohdaten und Skripte gelesen und gegeneinander gehalten |
| Host-Werkzeuge | nicht benutzt | Greps, Diffs und Lesen mit den Werkzeugen der Umgebung; kein Host-Go, kein Host-Textersetzer |

## 2. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff 2533ecc..HEAD --stat`)

35 Dateien im Diff: 21 neu, 1 gelöscht (`open/`-Datei des Moves), 13 geändert. Zählung stimmt mit dem Abgleich im Slice (Ende von *Messung*) überein.

| Diff-Gruppe | Plan-Zeile | Befund |
|---|---|---|
| `ADR-0025` (neu) und `docs/plan/adr/README.md` (Zeile `Proposed`) | Zeile 1 und Zeile 5 (Fix-Runde, `Schärft:`) | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/user/benutzerhandbuch.md` | Zeile 2 | ok: drei Stellen (Aufzählungspunkt Neuaufbau, Absatz *Grenze der Namensauflösung*, Absatz *Grenze* unter *Fallbacks sehen*) und der `sshd`-Verweis; die Tabelle nennt alle vier |
| `CHANGELOG.md` | Zeile 6 | ok: `Fixed` mit zwei Einträgen (Neuaufbau, zwei gemessene Grenzen) |
| `open/slice-cr-devcontainer-ollama.md`, `open/slice-v2-spec-sandbox-egress-ollama.md`, `open/slice-v1-cli-egress-anzeige.md`, `open/README.md` | Zeile 3 (Fix-Runde, F-2) | ok: Diff vorhanden, Zustandsmarker auf R0 | <!-- d-check:ignore (Kennungen als Zitat) -->
| `welle-sandbox-egress-ollama.md` | Zeile 3 (Fix-Runde) und Zeile 8 | ok: Entscheidung R0 in Tabelle §4, Closure-Trigger, Zielbild, Versionen |
| `BEO-WZG/host-werkzeug-ausserhalb-docker/` (Evidence neu, `state.md`) | Zeile 4 | ok |
| Messprotokoll-Verzeichnis, 17 Dateien (5 Skripte, `inner-m1.sh`, `inner-m2.sh`, `probe.pl`, `reload-atomar.sh`, 2 Dockerfiles, README, 5 Rohdaten) | Zeile 7 | ok, `ls` zählt 17 |
| `roadmap.md` | Zeile 9 | ok: Arbeitszeile; Rückstellung auf „Nichts in Arbeit.“ ist Closure-Aufgabe |
| Verweis-Nachzug des Moves (`done/slice-v1-egress-wildcard-pruefung.md`, `done/slice-cr-sandbox-egress-anzeige.md`, `BEO-DOK/…/state.md` und drei weitere) | Zeile 10 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| Slice-Datei (Move, Löschung in `open/`) | Zeile 11 | ok |
| `docs/reviews/` (Review-Report, dieser Report) | Zeile 12 | ok; dieser Report ist die letzte Datei der Range, Zeile 12 deckt ihn |
| `internal/**`, `spec/**`, Lastenheft, `Makefile` | „kein Produktcode“ | durch den Diff bestätigt: keine Datei dort |
| Folge-Slices `slice-cr-sandbox-egress-neuladen`, `slice-v1-sandbox-egress-neuladen` | Zeile 14 „entfallen“ | bestätigt: `ls open/` findet keine Datei „neuladen“ außer dem Entwurfs-Slice | <!-- d-check:ignore (Kennungen als Zitat) -->

Kein Diff ohne Plan-Zeile, keine Plan-Zeile mit Datei-Eintrag ohne Diff (ausgenommen die ausgewiesenen Reviewer-, Planer- und „entfallen“-Zeilen).

## 3. DoD-Abgleich (Slice §2)

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| Messung M1..M4 „soweit in §3 ausgewiesen“ | **erfüllt, der Haken trägt** | Wortlaut ist die engere Fassung (Review F-4): nicht gemessen sind M3 bei entfernter Engine, M2 Lesen der Statusdatei nach dem Neu-Lauf, Mehrfachläufe bei M1, M3, M4; jeweils in §3 und im ADR (Absatz *Nicht gemessen*) benannt. Messstand mit Docker-, Colima-Version, Architektur und Digest des Basis-Images steht in README, §3 und ADR gleichlautend (Digest `6ee2b409…`). |
| Entscheidungsvorlage und Folge-ADR | **erfüllt, Status `Proposed` zu Recht** | `ADR-0025` MADR-Form (Kontext, Entscheidung, Verglichene Alternativen, Konsequenzen, Fitness Function, Re-Evaluierungs-Trigger, Geschichte); `ADR-0012` und `ADR-0022` unverändert (`doc-immutable` grün) | <!-- d-check:ignore (Kennungen als Zitat) -->
| Verzweigung ausgeführt | offen, zu Recht | in der Welle vollzogen (§4: „Folge-Slices zum Neu-Laden gibt es nicht“), Handbuch berichtigt; Vermerk in §7 folgt mit der Closure |
| `make gates` grün | erfüllt | Abschnitt 1, ein eigener Lauf nach der Fix-Runde |
| Review, Report liegt vor | erfüllt, Haken noch nicht gesetzt | Report vorhanden; Behandlung der Findings Abschnitt 4; Haken ist Closure-Schritt |
| Doku-Update (Handbuch, ADR-Index, CHANGELOG) | erfüllt | Handbuch, CHANGELOG `Fixed` und Index-Zeile im Diff |
| Closure-Notiz §7, Steering-Loop-Eintrag | offen | §7 „Wird bei der Closure gefüllt.“ |
| Beobachtungs-Register | teils, als offen ausgewiesen | `BEO-WZG` auf drei Belege fortgeschrieben (vorgezogen, ausgewiesen); `BEO-DOK/messprotokoll…/state.md` im Diff; Verkörperung und Verknüpfung in §7 stehen aus |
| Risiken §6 mit Ausgang | erfüllt | jede Zeile trägt „eingetreten / entfallen / weiter offen“ mit Beleg |
| Drei Paarungen | offen | Welle-Closure |

Entscheidungs-Konformität R0 (Stichproben): kein neuer Port und kein Programmaufruf im Container (`internal/` nicht berührt); `SPEC-045` und `LH-FA-DEV-008` unverändert; R1 bis R4 und S nur dokumentiert, nicht entschieden (ADR Entscheidung 3, Tabelle der Alternativen); Welle §4 und Closure-Trigger tragen R0; der Spezifikations-Slice führt `SPEC-045` als „bleibt unverändert“ (Diff gelesen). <!-- d-check:ignore (Kennungen als Zitat) -->

## 4. Review-Findings F-1..F-10

Der Review-Report trägt keine Spalte „Behandlung“; die Zuordnung zur Fix-Runde `606157e` steht in ihrer Commit-Message und im Slice (Plan-Tabelle, §6). Verifier-Prüfung am Baum:

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 MEDIUM Handbuch „nur erlaubte Namen“ | **behoben** | Absatz *Grenze der Namensauflösung* im Egress-Abschnitt (Diff gelesen): direkte Abfrage des Upstream-Resolvers, Adresse eines Namens außerhalb der Liste, ohne `sudo`, Verbindung bleibt gesperrt, Datenkanal nicht gemessen, Messstand und „je ein Lauf“. Der Satz „lokaler DNS-Server, der nur erlaubte Namen auflöst“ steht weiter, trägt aber den Absatz darunter als Einschränkung; vertretbar (V-2). Absatz *Grenze* unter *Fallbacks sehen* ergänzt (Statusdatei ersetzbar). CHANGELOG `Fixed`, zweiter Eintrag. |
| F-2 Zustandsmarker „bedingt auf V2“ | **behoben** | Suche nach `V2`, `Neu-Laden`, `reload`, `bedingt` über Welle und `open/`: verbleibende Treffer nennen R0 („entsteht nach R0 nicht“) oder sind fachfremd (Roadmap-Meilenstein V2, „bedingt“-Zeilen der Anzeige-Plan-Tabelle mit eigenem Abschlusssatz). Keine Spur von „bedingt auf V2“, kein Neu-Laden-CR als Plan. |
| F-3 `Schärft:` | **behoben** | `ADR-0025` `Schärft:` nennt nur den Guardrail-Satz und weist aus, dass der Weg „Liste ändern“ im Handbuch steht; `SPEC-045` als nachzuführende Stelle. | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 Haken-Wortlaut | **behoben** | Haken „soweit in §3 ausgewiesen“ mit Aufzählung des Nicht-Gemessenen; „Status ok“ aus der M2-Frage genommen. |
| F-5 Host-Werkzeug | **behoben im Rahmen des Möglichen** | Slice §6 benennt den Regelverstoß; Evidence-Datei und `state.md` (drei Belege, Schwelle erreicht, sobald der Vorgang in `done/` liegt, Verkörperung nicht entschieden). Zahl und Erfolg der Aufrufe bleiben aus dem Repo nicht ableitbar (so ausgewiesen). Der Verifier hat selbst kein Host-Werkzeug zur Textänderung benutzt. |
| F-6 Voraussetzungen | **behoben** | README „Voraussetzungen am Host“ nennt `head`, `tail`, `cut`, `sed`, macOS-Fassungen, den Docker-Socket im Messcontainer und die Images. |
| F-7 betroffene Stellen | **ausgewiesen, offen beim Projektinhaber** | Slice §6 „Offene Frage an den Projektinhaber“ (DNS-Pfad am Upstream, `SPEC-045` *Grenze*, `ADR-0022` Entscheidung 8, `ADR-0012`); nichts geändert, Accepted-ADRs unberührt. Auflage 2. | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-8 Evidence-Verweis aus ADR | **ok** | `make docs-check` (Teil von `make gates`) meldet 0 Befunde; Muster wie `ADR-0024`. | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-9 Host-Benutzername | **teilweise behoben** | `rohdaten-04-m3.txt` Zeile 3 trägt „<Host-Benutzer geschwärzt>“; Messprotokoll und ADR frei von Benutzername. **Rest (V-3):** der Review-Report selbst zitiert den Namen zweimal (Findings-Tabelle F-9 und Negativbefund), und die Commits `169cbad` bis `e88d6ab` tragen ihn im Verlauf. |
| F-10 Verifier-Hinweise | **berücksichtigt** | (a) R0 nur als Vereinbarung belegt: siehe Auflage 1; (b) Abgleich nach der Fix-Runde im Slice wiederholt und von mir gegen den Diff bestätigt (Abschnitt 2); (c) M4 je ein Lauf, Handbuch trägt Messstand und „Editor-Aktion nicht gemessen“. |

## 5. Handbuch und Aussagen gegen Rohdaten

**Kennungen im Handbuch (`AGENTS.md` §5 Zeile 6):** Suche nach `LH-`, `ADR`, `SPEC-`, `slice-`, `welle-`, `docs/plan`, `spec/spezifikation`, `spec/architecture` über `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`: kein Treffer; im Diff des Handbuchs keine hinzugefügte Zeile mit Kennung oder Planungsverweis. **Ohne Befund.**

Stichproben (fünfzehn), Handbuch, CHANGELOG, ADR und Slice gegen die Rohdaten im Messprotokoll:

| # | Aussage | Beleg | Ergebnis |
|---|---|---|---|
| 1 | „bauen den Container neu auf“; Neustart führt altes Script mit alter Liste aus | `rohdaten-05-m4.txt` Abschnitt 2: gleicher Container `83790fb46dca`, gleiches Image `8c6bfba3dc84`, gebackene Liste ohne `example.org`, `http=000`, Statusdatei neu (`18:20:39Z` gegen `18:20:33Z`) | stimmt |
| 2 | `devcontainer up --remove-existing-container` setzt neue Liste durch, Image neu gebaut | Abschnitt 4: Container `5fbdc8f6cc34`, Image `ab9d55a6b8f9`, `example.org` in der gebackenen Liste, `http=200` | stimmt |
| 3 | `devcontainer up` gegen laufenden Container ändert nichts | Abschnitt 3: Container und Image unverändert, Liste alt | stimmt |
| 4 | Host-Datei trägt die neue Liste schon vor dem Stopp | Abschnitt 1: Host-Datei mit `example.org`, Container unverändert | stimmt |
| 5 | Statusdatei: Datei `root`, Verzeichnis `vscode`, `vscode` ersetzt sie (Handbuch *Grenze*) | `rohdaten-02-m1.txt`: `drwxr-xr-x vscode vscode …/status`, `-rw-r--r-- root root egress`, `ersetzen rc=0`, neue Datei `vscode vscode` mit `at=2000-01-01T00:00:00Z` | stimmt |
| 6 | Prozess ohne `sudo` bekommt über den Upstream-Resolver die Adresse eines Namens außerhalb der Liste (Handbuch *Grenze der Namensauflösung*) | `rohdaten-02-m1.txt`: Abfrage an `192.168.5.1` liefert `104.20.23.154`, `172.66.147.243` | stimmt |
| 7 | „Verbindung zu dieser Adresse bleibt gesperrt“ | gleiche Datei: `Connection timed out after 5006 milliseconds` zur Adresse `104.20.23.154` | stimmt (ein Lauf, so ausgewiesen) |
| 8 | „Wer die Regeln ändern will, braucht weiter `sudo`“ | `nft` ohne `sudo` „Operation not permitted“, `CapEff` 0, fünf Dateien „NICHT schreibbar“, `kill` „Operation not permitted“; mit `sudo` zwei `delete` rc 0 | stimmt |
| 9 | Passwortloses `sudo` (`NOPASSWD: ALL`) | `sudo -n true rc=0`, `(root) NOPASSWD: ALL` | stimmt |
| 10 | Fenster ohne Filter 24 bis 45 ms; Läufe 50 bis 78 ms; atomare Variante ohne offenes Fenster | `rohdaten-03-m2.txt`: `dauer_s` 0.051 (erster Lauf), Zustandswechsel „egress offen“ von `…389.710` an; atomarer Lauf 0.023 s ohne „egress offen“ | stimmt (Einzelwerte der übrigen Läufe aus dem Review übernommen, Stichprobe der Läufe 1 und atomar 1 selbst gelesen) |
| 11 | Startrennen: 50 Läufe, 0 auffällig; alter `dnsmasq` endet 1,9 bis 3,4 ms nach `SIGTERM` | Abschnitte B und F: „50 Laeufe, 0 mit auffaelligem Zustand danach, 0 mit Fehlermeldung“; Abschnitt G: 1932 µs bis 3419 µs, 20 Läufe | stimmt |
| 12 | Laufende Verbindung überlebt, neue nicht; Set danach leer | Abschnitt C: Bytes 1 302 528 auf 3 465 216 über 9 s, neue Verbindung „Could not resolve host“; Abschnitt D: Set mit `140.82.121.6 expires 59m59s997ms` vorher, danach leer, ohne DNS `http=000`, mit DNS `200` | stimmt |
| 13 | `docker cp` Eigentümer `501 dialout`, `-a` ergibt `vscode`, Tar `owner=0` ergibt `root`; gestoppter Container: `exec` scheitert, `cp` geht; nach `docker start` keine Tabellen, nach `postStart` 000 | `rohdaten-04-m3.txt` Abschnitte 1, 5, 6, 7 | stimmt |
| 14 | Messstand (Docker 29.5.2, Colima 0.10.3, arm64, Digest `6ee2b409…`) in ADR, Slice, README gleich | `README.md` des Protokolls, ADR *Messung*, Slice §3 | stimmt |
| 15 | Statusdatei nach dem Neu-Lauf (M2) „nicht gemessen“ | keine Statuszeile in `rohdaten-03-m2*.txt` nach dem Neu-Lauf | stimmt (Ausweis korrekt) |

**Aussagen über Code (`AGENTS.md` §3.10):** `egress-init.sh.tmpl` („re-run safe“, `kill` ohne Warten, `nft delete table`), `Dockerfile.tmpl` (`COPY egress-init.sh …`, `chown vscode:vscode /var/lib/u-boot/status`), `devcontainer_sandbox.go` (`egressRelaxation`, `sudo sh /usr/local/bin/u-boot-egress-init` im `postStartCommand`), `devcontainer_inspector.go` („never runs a program in the container“): die vom ADR und vom Slice genannten Symbole und Zitate sind im Review an der Stelle gelesen belegt; der Verifier hat sie nicht erneut im Code gelesen, sondern die Rohdaten gegen die Aussagen gehalten (Abschnitt 9). Die Handbuch-Zusage „ein Neustart führt das alte Script aus“ ist Verhalten der Container-Engine, im ADR als „Ein Test zu M4 entsteht nicht“ begründet (Messstand statt Test, siehe V-4).

**Weitere Stellen „Container neu starten“ für die Egress-Liste:** Suche nach `neu starten`, `Neustart`, `restart` über `docs/user`, `README*`, `CHANGELOG.md`, `spec`, `internal`, ADRs, `open/`, `in-progress/`: Treffer nur in der berichtigten Handbuch-Stelle, im CHANGELOG-Eintrag, im ADR (Zitat des alten Satzes in der Messtabelle), in der Plan-Prosa des Slice (Frage M4) und im Kommentar „restart-safe“ des Script-Templates (betrifft `resolv.conf`, nicht die Liste). Der Handbuch-Hinweis zur Feature-Liste (`sshd`, Schritt 1) verweist nun auf den gemessenen Egress-Absatz; der Anker `#internet-beschränkung-egress` löst (`make docs-check` grün). **Kein verbleibender Satz „neu starten genügt“.**

**Folge-Slices ohne „bedingt auf V2“-Reste:** `slice-cr-devcontainer-ollama` (Versionszeile ohne Neu-Laden-CR), `slice-v2-spec-sandbox-egress-ollama` (`SPEC-045` „bleibt unverändert“, „Ein Vertrag, vier Stellen“ auf „Neuaufbau“), `slice-v1-cli-egress-anzeige` (kein `reload`), `open/README.md` (Zeile mit Ergebnis R0 und Link nach `in-progress/`): im Diff gelesen, ohne Rest. <!-- d-check:ignore (Slice- und SPEC-Kennungen als Zitat) -->

## 6. Hygiene, ADR-0025

**Hygiene** (Suche über den gesamten Diff und das Messprotokoll): kein Token, kein `X-Amz`, kein privater Schlüssel, kein `ghp_`/`gho_`, kein `Bearer`, kein Home-Pfad im Messprotokoll, in der ADR, im Slice und in der Evidence. **Einzige Treffer:** der Host-Benutzername im Review-Report (zweimal, als Zitat des Findings F-9), siehe V-3. UID 501 und Gruppe `dialout` sind keine Personendaten.

**`ADR-0025`:** Status `Proposed`, Index-Zeile `Proposed`; Optionen R0, R1, R2, R3 (verworfen, begründet), R4 und S (nicht gemessen, ausgewiesen) mit gemessenen Zellen; `Schärft:` nur der Satz, den das Lastenheft trägt, plus nachzuführende `SPEC-045`; ergänzt, ersetzt nichts. **Decke-Regel:** die ADR (Rang 4) verlinkt aufwärts auf Lastenheft und Spezifikation und abwärts auf das Messprotokoll (Evidence, kein Slice, Carveout oder Roadmap); sie nennt „Spezifikations-Slice der Welle“ ohne Kennung; `spec/architecture.md` unverändert; `make docs-check` (Referenzmatrix) 0 Befunde. Entscheidung 1 zusagt „kein Befehl, der ein Programm im Container ausführt“ und stützt sich auf den Ist-Stand der Ports. Der Absatz *Nicht gemessen* ist vollständig (Docker Desktop, Linux-Host, `amd64`, entfernte Engine, VS Code, `nestedRuntime`, Last, Datenkanal über DNS, Mehrfachläufe). <!-- d-check:ignore (Kennungen als Zitat) -->

## 7. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | INFO | Plan-Tabelle Zeile 4 bündelt Register-Fortschreibung, Messprotokoll-README und Schwärzung in einem Eintrag; vollständig, nur lang. |
| V-2 | LOW | Handbuch: „lokaler DNS-Server, der **nur erlaubte Namen** auflöst“ steht unverändert vor dem neuen Absatz *Grenze der Namensauflösung*. Der Absatz korrigiert die Lesart unmittelbar danach; wer nur den ersten Satz liest, behält die alte. Zu entscheiden (Projektinhaber, Doku-Autor): Satz präzisieren oder lassen. Kein Blocker. |
| V-3 | LOW | Der Host-Benutzername steht noch im Review-Report `docs/reviews/2026-10-18-slice-v1-sandbox-egress-neuladen-entwurf.md` (zwei Stellen, als Zitat von F-9) und im Git-Verlauf (`169cbad` bis `e88d6ab`). Die Schwärzung in `rohdaten-04-m3.txt` ist wirksam, aber unvollständig. Kein Geheimnis, aber ein Personenbezug, den die Fix-Runde nicht schließt. Der Verlauf wird nicht umgeschrieben (nicht beauftragt). Zu entscheiden: Report schwärzen (eigener Commit). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-4 | INFO | Die Handbuch-Zusage „Neustart führt das alte Script mit der alten Liste aus“ hat keinen Test (begründet im ADR: Verhalten der Dev-Container-CLI und des Docker-Caches). Die Bedingungsregel `AGENTS.md` §3.10 (letzter Spiegelstrich) verlangt für Eigenschaften des Codes einen Test; hier ist es keine Eigenschaft des Codes, sondern eine Messung mit Messstand und Grenze, so ausgewiesen. Hinweis an den Planner, falls die Lesart enger gezogen wird. |
| V-5 | INFO | Ruheort-Regel: Das README des Messprotokolls und die Evidence-Datei von `BEO-WZG` verlinken den Slice mit `in-progress/`-Pfad (die Evidence-Datei benennt das ausdrücklich). Der Move nach `done/` läuft über `make slice-mv`, das Verweise nachzieht; sonst bricht `docs-check`. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-6 | INFO | Der Review-Report trägt keine Spalte für den Behandlungsstand; die Zuordnung steht nur in Commit-Message und Slice. Dieser Report schließt die Lücke (Abschnitt 4). |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, keine Abweichung von R0.

## 8. Auflagen für die Closure

1. **Closure erst nach `Accepted` von `ADR-0025`** durch den Projektinhaber (Bestätigung des Wortlauts; im Gespräch zu belegen und in §7 zu vermerken: Entscheidung R0 vom 2026-10-18, Bestätigung des ADR-Wortlauts, Datum). Wählt er statt R0 eine andere Option, sind ADR, Welle und die zurückgenommenen Zustandsmarker neu zu prüfen. <!-- d-check:ignore (Kennung als Zitat) -->
2. **Antwort zur DNS-Frage** (Slice §6, Review F-7): Der Projektinhaber entscheidet, ob ein eigener Folge-Slice den offenen DNS-Pfad am Upstream schließt oder ob der Ist-Stand mit den neuen Handbuch-Absätzen genügt. Die Antwort steht in §7; die Spec-Nachführung (`SPEC-045` *Grenze*) übernimmt der Spezifikations-Slice der Welle. Bis dahin bleibt der Risiko-Ausgang „weiter offen“. <!-- d-check:ignore (Kennungen als Zitat) -->
3. §7 füllen: Messergebnis M1 bis M4, Entscheidung (V1/R0), `ADR-0025` mit Status, **keine Folge-Slices**, Steering-Loop-Eintrag, Register (`BEO-WZG/host-werkzeug-ausserhalb-docker`: dritter Beleg, Schwelle erreicht mit `done/`; Verkörperung entscheidet der Planner mit dem Projektinhaber), `Verantwortlich:` (bleibt `—`, wenn der Projektinhaber es nicht setzt; in §7 festhalten). Die „weiter offen“-Ausgänge in §6 nach der Antwort zu 2 aktualisieren. <!-- d-check:ignore (Kennungen als Zitat) -->
4. DoD-Haken Verzweigung, Review, Closure-Notiz, Register, Paarungen setzen; Plan-Tabelle gegen `git diff --stat` **nach der letzten Änderung** wiederholen (`AGENTS.md` §3.10; dieser Report ist die letzte Datei der Range; die Zahl 35 ändert sich durch ihn und durch die Closure).
5. Roadmap auf „Nichts in Arbeit.“ zurücksetzen und den Link der Nächste-Wellen-Zeile auf `done/` nachziehen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3), über `make slice-mv`. <!-- d-check:ignore (Kennungen als Zitat) -->
6. Optional: V-3 (Host-Benutzername im Review-Report schwärzen), V-2 (Satz im Handbuch präzisieren).

## 9. Nicht verifiziert

Wiederholung der Messung (Colima, Netz, Dev-Container-CLI-Image); Messstände außerhalb Colima/arm64 (Docker Desktop, Linux-Host, `amd64`, entfernte Engine, VS Code); Rennfreiheit des Startscripts unter Last; Datenkanal über die Namensauflösung; erneutes Lesen der Code-Stellen (`generate.go`, `Dockerfile.tmpl`, `egress-init.sh.tmpl`, `devcontainer_inspector.go`): die Aussagen über Code sind aus dem Review und aus den Rohdaten (erzeugtes Script, gemessenes Verhalten) gestützt, nicht neu am Quelltext gelesen; ob die nicht abgewiesenen Host-Aufrufe (`perl`, `sed -i`) des Implementers gewirkt haben (Aussage des Auftrags, nicht aus dem Repo ableitbar); `make test-docker` und `make verify-depguard` (kein Produktcode).
