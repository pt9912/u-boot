# Verification-Report: slice-v1-devcontainer-remove-messung — 2026-11-05

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-devcontainer-remove-messung`](../plan/planning/done/slice-v1-devcontainer-remove-messung.md) — Commit-Range `642aefe7..94f65083` (2 Commits: Messung `8861376d`, Fix-Runde nach dem Review `94f65083`; 19 Dateien, 2225 Einfügungen / 10 Löschungen). Abgleichstand **nach der letzten Fix-Runde = `94f65083`** (= Repo-Tip, Arbeitsbaum vor diesem Report sauber; `main` ist 7 Commits vor `origin/main`, nichts davon gepusht).

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel, Messfragen M1 bis M5, Nicht-Ziele), §2 (DoD, Closure-Pflichten), §3 (Plan-Tabelle, „Gelesen für die Messung“), §6 (Risiken), §7 (leer), §8
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Spiegelstriche zu `up --recreate`, `down`, `shell`, Isolation der Volumes, Fehlerklassen), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Volume der Ollama-Modelle), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Volume-Namen je Instanz)
- Messprotokoll [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) mit `rohdaten-m1.txt` bis `rohdaten-m5.txt`, `rohdaten-aufraeumen.txt`, `rohdaten-m1-erster-lauf.txt`; Skripte unter `tools/messung/<Kennung des Slice>/` (`00-alles.sh`, `01` bis `04b`, `05-aufraeumen.sh`, `lib.sh`)
- Review-Report [`2026-11-05-slice-v1-devcontainer-remove-messung`](2026-11-05-slice-v1-devcontainer-remove-messung.md) (0 HIGH / 0 MEDIUM / 7 LOW / 3 INFO; Eingang, nicht Evidence); sein Gegenstand endet bei `8861376d`, die Fix-Runde `94f65083` hat kein Review gelesen
- Vorbild der Form: [`verification-2026-11-05-slice-v1-devcontainer-homebrew-pakete`](verification-2026-11-05-slice-v1-devcontainer-homebrew-pakete.md); [`harness/verification.md`](../../harness/verification.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.6, §3.7, §3.9, §3.10, §5 Regel 6), [`harness/roles.md`](../../harness/roles.md)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-devcontainer-remove-messung`](../plan/planning/done/slice-v1-devcontainer-remove-messung.md) (Messslice, kein Go-Code)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md)
- Artefakte: Messskripte, Messprotokoll mit Rohdaten (M1 bis M4 laut DoD, M5 als Zusatz der Fix-Runde), Zeile im Bestand [`docs/messungen/README.md`](../messungen/README.md)

### DoD-Abgleich (Plan §2)

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | **Messprotokoll** M1 bis M4 je mit Befehl, Ausgabe, Datum und Grenzen; Skript unter `tools/messung/<Kennung des Slice>/` (Haken `[x]`) | **erbracht, Haken zu Recht gesetzt** | Befehl: in den Rohdaten als `$ …`-Zeilen und als Ablauf in den Skripten, im Protokoll je Frage beschrieben. Ausgabe: je Frage eine Rohdatei mit Zeitstempel (13:25:42Z bis 13:29:05Z, Aufräumen 13:29:32Z; Datum im Protokollkopf „2026-11-05“). Grenzen: Abschnitt „Grenzen“ und je Frage „Darf nicht tragen“. Ablage: Skripte (alle `100755`) nur unter `tools/`, nur `.md`/`.txt` unter `docs/messungen/`; `make ablage-pruefung` grün. **Abgleich Zahl für Zahl** (eigene Probe, ohne Docker): Umgebung (Docker 29.9.3, `overlayfs`, cgroup v2, Kernel `6.8.0-147-generic`, `u-boot version 0.23.0-dev`, Hilfs-Image `0.89.0-f1a94cce795e`); M1-Tabelle (7 Varianten: Zahl der Volumes in `Mounts`, Bind-Mounts, Ziele, `mounts-vollstaendig=ja`); 52-stellige Namensenden, `verschiedene-ids=1`, `labels=null`, anonymes Label `com.docker.volume.anonymous`; Bauzeiten (archiviert 4 bis 8 s, erster Lauf 5/7/10/21/46/48/5 s); 14 Images, Größen 1,41 bis 2 GB, Platzdifferenz −1.937.408 Byte (66.302.042.112 − 66.303.979.520); Runtimes (`daemon`/`vfs` `result=ok`, `podman version 6.1.4`); M2 (Exit-Codes 1/0, Fehlertexte, `volume is in use` mit Container-ID, Tabelle zu anonymen Volumes, `--recreate`); M3 (drei Container, Filter, `gemeinsame=0`, Standardprofil); M4 (`links=0`, 278,6 kB / 6 B / 2,519 kB, Infix `-ollama-` ein Treffer, Prüfsumme des öffentlichen Schlüssels gleich): jede Zahl und jeder Name steht in einer Rohdatei. **Gegenprobe der Läufe-Herkunft** (`git show 8861376d:…/rohdaten-m1..m4.txt`, Zeilen `AUSSAGE`, 52-stellige Namensenden gleichgesetzt, `diff`): gleich, einzige Abweichung die neue Zeile `m3-sbx-label-u-boot.instance-schluessel`, wie das Protokoll es sagt; `rohdaten-m1-erster-lauf.txt` gegen `rohdaten-m1.txt`: `diff` leer, Kopfzeilen `engine:`/`hilfs-image` fehlen dort, Dauern wie im Protokoll. **Vorbehalte:** V-1, V-2, V-3 (drei Sätze unter „Darf tragen“ breiter als die Messung). |
| 2 | **Eingang für den CR:** je Frage der Satz, den der CR tragen darf, und die Sätze, die er nicht tragen darf (Haken `[x]`) | **erbracht in der Form; drei Sätze über dem Gemessenen (V-1 bis V-3) — Haken vorbehaltlich** | Abschnitte „Eingang für den CR“ mit „Darf tragen“ / „Darf nicht tragen“ stehen zu M1, M2, M3, M4 und (Zusatz) M5. Jeder Satz „Darf tragen“ wurde gegen die `AUSSAGE`-Zeilen gelesen: M1 (Erkennung über `Mounts`, Vollständigkeit in 7/7 Varianten, Namensfamilien, anonymes Volume), M2 (`docker rm` laufend/gestoppt, `volume rm` mit Halter, `-f` ändert nichts, anonyme nur mit `-v`, `down` ohne Container Exit 0), M3 (Isolation der Workspace-Volumes bei drei Containern, Entfernen einer Instanz über deren `Mounts` trifft die anderen nicht), M4 (alle benannten Volumes überleben `docker rm -f -v`, Wiedereinhängen mit Inhalt, Unterscheidung nur über den Namen, `dangling=true`) sind gedeckt. Die „Darf nicht tragen“-Listen nennen jeweils die nicht gemessene Menge (andere Hosts, `prune`, `paused`/`created`, drei und mehr Instanzen, Podman-Instanzen, Neubau des Images, Umzug des Projektordners, `--recreate` für Docker-Speicher und Modelle, `shell` ohne `--instance`) und halten sie aus „Darf tragen“ heraus. Die Entscheidung über `--volumes` bleibt ausdrücklich beim CR (Nicht-Ziel §1). |

Pflichten der Closure (Plan §2, zählen nicht mit):

| Pflicht | Zustand | Evidence |
|---|---|---|
| `make gates` grün (Haken `[x]`) | **erbracht, Haken zu Recht gesetzt** | eigener Lauf auf `94f65083`, siehe Sensors |
| Review durchgeführt (Haken offen) | **erbracht für `8861376d`, Haken steht aus** | Report vorhanden und mit `94f65083` eingecheckt; die Fix-Runde (M5-Skript, M5-Rohdaten, 133 Zeilen Protokoll, neu gelaufene Rohdaten) hat kein Nachlauf-Review gelesen — vom Verifier gelesen (siehe V-6). Der Haken gehört dem Planner. |
| Verification-Report (Haken offen) | **mit diesem Report erbracht** | Haken gehört dem Planner |
| Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang (Haken offen) | **offen** | §7 trägt den Platzhalter; beide Risiken in §6 stehen auf „**Ausgang:** offen“ mit Inhalt. Was die Notiz tragen muss: Abschnitt „Was die Closure-Notiz tragen muss“ unten. |

### M5 (Zusatz der Fix-Runde, nicht in der DoD) gegen `rohdaten-m5.txt`

| Fall | Aussage im Protokoll | Rohdaten | Ergebnis |
|---|---|---|---|
| A | Projekt ohne Container: `down` und `down --instance a` Exit 0, „nothing to stop“ | `m5-a-down-ohne-container-exit=0`, `m5-a-down-instanz-a-ohne-container-exit=0` | gedeckt |
| B | drei Container laufen, `down` stoppt nur den ohne Instanz (Exit 0) | `m5-b-down-ohne-instanz-exit=0`, `ohne-instanz-laeuft=false a-laeuft=true b-laeuft=true` | gedeckt |
| C | nur `a` und `b` laufend: `down` Exit **10** mit Hinweis, nichts gestoppt; `up` Exit **10**, kein Container ohne Instanz angelegt; Warnung zum Engine-Socket steht vor der Meldung | `m5-c-down-…-exit=10`, `a-laeuft=true b-laeuft=true`, `m5-c-up-…-exit=10`, `m5-c-container-ohne-instanz-angelegt=nein`; `log: Warning: this command mounts the host container-engine socket …` | gedeckt |
| D | `a` gestoppt, `b` läuft: `down` Exit 10, `b` läuft weiter | `m5-d-down-ohne-instanz-exit=10`, `a-laeuft=false b-laeuft=true` | gedeckt |
| E | beide gestoppt: `down` Exit 10 | `m5-e-down-ohne-instanz-nur-gestoppte-instanzen-exit=10` (Zustand `exited`, `instance=a/b`) | gedeckt |
| — | „Fälle D und E zeigen, dass die Suche auch gestoppte Container findet (`docker ps -a`)“ | `state=exited` in den `zustand`-Zeilen; Code `FindContainer`: `ps -a` | gedeckt (Messung und gelesener Code) |
| — | Exit 10 stimmt mit [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) überein | Fehlerklassen-Spiegelstrich: „wenn `up`, `down` oder `shell` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen … gibt“ | konform |

Das Inventar (7 fremde Container, 2 fremde Volumes) steht vor und nach dem Lauf in `rohdaten-m5.txt`, Differenz leer, „Reste mit dem Präfix: 0“. Der Fall „Container ohne Instanz gestoppt, daneben Instanzen“ ist nicht gemessen und im Protokoll nicht beansprucht (nur der gelesene Code, `!container.Running` → Antwort mit `ContainerID`, ohne `Stopped`). M5 ändert die DoD nicht.

### Prüfung der Review-Behebung (F-1 bis F-10 des Reviews)

| Finding | Behoben/geführt? | Beleg |
|---|---|---|
| F-1 (LOW, Zahlen und Läufe ohne Rohdaten) | **ja, mit Rest V-5** | Die Inventur (Zahl und Namen der Container und Volumes) steht in allen fünf Rohdateien M1 bis M5 vor und nach dem Lauf, Differenz leer. Das Protokoll hat den Abschnitt „Läufe und ihre Herkunft“: die archivierte Fassung ist **ein** Lauf (13:25 bis 13:29), die Vorfassungen M2 bis M4 sind als nicht archiviert benannt und stützen nichts, `rohdaten-m1-erster-lauf.txt` ist als Lauf einer nicht archivierten Skriptfassung ausgewiesen (Gegenprobe oben). Rest: V-5 (die Namen selbst). |
| F-2 (LOW, „Kein Code des Repos liest `docker inspect`“) | **ja** | Slice §3 trennt jetzt Produktions- und Testcode. Eigene Suche in `internal/` und `cmd/` (ohne `_test.go`): nur `context inspect` (`inspector.go`, Zeile 199) und `image inspect` (`devcontainer_launcher.go`, Zeile 309). Testcode: `dcInspect`, `dcInspectDoc` (`internal/e2e/devcontainer_start_docker_test.go`, Feld `Mounts`), Nutzung in `devcontainer_instance_docker_test.go` und `forge_cli_docker_test.go`, `inspectDoc` im Paket `docker` (Launcher-Test), `sandbox_status_docker_test.go` Zeile 421 (`inspect --format {{.Id}}`) — alle vorhanden. |
| F-3 (LOW, Zustandsmarker in §8) | **ja, mit INFO V-8** | §8 trägt Datum und Ergebnis der Sichtung; die drei genannten Beobachtungen liegen im Register (`BEO-WZG/aufraeum-skript-ueber-praefix-hinaus`, `BEO-DOK/zusage-weiter-als-messung`, `BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`); Plan-Tabelle führt jetzt „§6 und §8“. |
| F-4 (LOW, Label-Satz) | **ja** | `03-m3-instanz.sh` misst den Schlüssel in `Config.Labels` (`hat_label`: Schlüsselnamen der Label-Map, `grep -qxF`); Rohdaten `ohne-instanz=fehlt a=vorhanden b=vorhanden`; das Protokoll trennt den Schlüssel (M3) von der leeren Zeichenkette in `docker ps`; M5 nennt „der Code (gelesen) unterscheidet am leeren Wert; dass der Schlüssel fehlt, misst M3“. |
| F-5 (LOW, Ollama-Satz) | **teilweise** | „bestätigt, was die Spezifikation sagt“ und „neben den Modellen“ sind weg; „Darf nicht tragen“ nennt „Es wurde keins geladen“; die Version `0.40.6` ist als Version des Clients ausgewiesen (Rohdaten: `Warning: client version is 0.40.6`). **Neu:** der Satz nennt „die leeren Verzeichnisse `models` und `cache`“ — die Rohdaten tragen das nicht (V-1). |
| F-6 (LOW, Standardprofil-Satz) | **teilweise** | „Im Standardprofil ohne Katalog-Feature …“ (M1). Dieselbe Verallgemeinerung steht unverändert in M3 „Darf tragen“ und in den Kurzfassungen (V-3). |
| F-7 (LOW, Fehlerzweig von `Down`) | **ja** | Slice §3 und Protokoll („Befunde aus dem Code“) nennen `ErrOnlyInstanceContainers` → `onlyInstancesError("down")` (Exit 10). Gegen den Code gelesen (`devcontainer_start.go`, `Down`, Zeilen 781 bis 812). Messung M5 (Fälle A bis E) schließt die Lücke und geht über den Review-Auftrag hinaus (Zusatz, nicht in der DoD). |
| F-8 (INFO, Kennung im Pfad) | **unverändert, kein Änderungszwang** | `d-check:ignore` an zwei Stellen. |
| F-9 (INFO, Kleinigkeiten) | **ja** | (a) „Zeilen 154 bis 186“ ersetzt durch Symbolangabe (`devcontainerTemplateData`, Zweig `if !sandbox`); (b) `alpine:3` im Protokoll und im Kommentar von `05-aufraeumen.sh` genannt; (c) M2f nennt den Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (vorhanden, `devcontainer_instance_docker_test.go` Zeile 112); (d) Lesehilfe steht im Absatz „Aufruf“, der Kommentar in `lib.sh` verweist auf „Abschnitt Aufruf“. |
| F-10 (INFO, Zuständigkeit) | **hiermit übernommen** | DoD-Abgleich, `make gates`, Risiko-Ausgänge: dieser Report. |

### Entscheidungs-Konformität

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) `down`: stoppt, entfernt weder Container noch Volumes | **bestätigt** | M2c/M2d: nach `down` Status `exited`, Container und Volume da; Code: `Down` ruft nur `launcher.Stop` (`docker stop`) |
| `--recreate` entfernt den Container, benannte Volumes bleiben | **bestätigt** | M2f: neuer Container, alter nicht mehr da, Volume gleichen Namens vorhanden, Markerdatei lesbar |
| Isolation der Workspace-Volumes zweier Instanzen (Zusage-Menge) | **bestätigt** | M3 `workspace-volumes-verschieden=ja`, `marker-isoliert=ja` (drei Container); Spezifikation sagt Docker-Speicher/Modelle ausdrücklich nicht zu, das Protokoll auch nicht |
| Ohne `--instance` nur der Container ohne Instanz; Exit 10 bei nur Instanz-Containern | **bestätigt** | M5 B bis E (zusätzlich zu `up`) |
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) (Volume-Namen je Instanz) | **bestätigt** | M3: gleiche Namensfamilie, verschiedenes 52-stelliges Ende je Container |
| Messung ändert keinen Produktionscode, kein Gate gelockert ([`AGENTS.md`](../../AGENTS.md) §3.2, §3.6) | **bestätigt** | `git diff --stat 642aefe7..HEAD`: kein Eintrag in `Makefile`, `repo.mk`, `.golangci.yml`, `.d-check.yml`, `.github/`, `internal/`, `cmd/`, `spec/`, `docs/plan/adr/`, `docs/plan/carveouts/`; die zwei `d-check:ignore` sind Zeilen-Ausnahmen mit Begründung. Kein ADR berührt (Accepted-ADRs unverändert, [`AGENTS.md`](../../AGENTS.md) §3.5). |
| Nutzer-Doku-Regel ([`AGENTS.md`](../../AGENTS.md) §5 Regel 6) | **bestätigt** | `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` nicht im Diff |
| Docker-only ([`AGENTS.md`](../../AGENTS.md) §3.1) | **bestätigt** | Skripte nutzen `docker`, das Binary `u-boot` (`make compile`), Bash und Standardwerkzeuge; keine Host-Go-Toolchain |
| Zugangsdaten | **bestätigt** | `GIT_TOKEN` ist ein Platzhalter (`lib.sh`), Suche in den Rohdaten nach `token`/`secret`: nur Platzhalter-Erwähnung im Protokoll |

### Aussagen über Code (§3.10), Stichprobe gegen den Code

| Aussage (Slice/Protokoll) | Fundstelle | Ergebnis |
|---|---|---|
| `Down`: `FindContainer`, bei `ErrOnlyInstanceContainers` `onlyInstancesError("down")`, Engine-Fehler → Fehler, ohne Fund leere Antwort, gestoppter Container → `ContainerID` ohne `Stopped`, sonst `launcher.Stop` | `internal/hexagon/application/devcontainer_start.go`, `Down` (781 bis 812), `onlyInstancesError` (366) | stimmt |
| `containerBeforeStart` gibt denselben Suchweg für `up` und `shell` ohne `--instance` | `devcontainer_start.go` Zeile 493 bis 506; Aufrufer `start` mit `command: "up"` (Zeile 208) und `"shell"` (Zeile 223) | stimmt |
| `FindContainer`: `docker ps -a --no-trunc` mit Ordner-Label, Instanz-Label nur mit Instanz; ohne Instanz erster Container mit leerem Label-Wert, sonst `ErrOnlyInstanceContainers`; Doc-Kommentar nennt `docker ps`, `docker cp`, `docker context inspect` | `internal/adapter/driven/docker/inspector.go` 43 bis 46, 85 bis 122 | stimmt |
| `Launcher.Stop` ist `docker stop <id>`; `runArgs`: `--remove-existing-container` bei `Recreate`, zwei `--id-label` nur mit Instanz | `devcontainer_launcher.go` 282 bis 290, 503 bis 508 | stimmt |
| Standardprofil erzeugt kein `mounts`; Sandbox: `-workspace`, `-containers` (podman), `-docker`, `-ollama`, jeweils `volumeMount` | `devcontainer_sandbox.go`, `devcontainerTemplateData` (Zweig `if !sandbox`, Zuweisungen an `data.Mounts`) | stimmt |
| Kein Produktionscode liest `docker inspect` eines Containers oder `Mounts` | Suche über `internal/`, `cmd/` ohne `_test.go` | stimmt (nur `image inspect`, `context inspect`) |
| Keine `VOLUME`-Anweisung in den Templates | `templates/devcontainer/Dockerfile.tmpl` (Suche: kein Treffer) | stimmt |
| Zeilen 45 (`roadmap.md`) und 23 (`open/README.md`) nennen den Slice in `in-progress/` / „in Umsetzung“ und sind „unverändert“ | gelesen | stimmt; im Diff nicht berührt |

### Plan-vs-Code (§3 gegen `git diff --stat 642aefe7..94f65083`, Stand nach der letzten Fix-Runde)

| Diff-Gruppe | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `tools/messung/…/` (8 Dateien: `00`, `01`, `02`, `03`, `04`, `04b`, `05`, `lib.sh`) | „Messskripte M1 bis M5 (`01` bis `04b`), `05-aufraeumen.sh`, `00-alles.sh`, `lib.sh`“ | deckt |
| `docs/messungen/…/README.md` | „Protokoll M1 bis M5“ | deckt |
| `docs/messungen/…/rohdaten-*.txt` (7: `m1`–`m5`, `aufraeumen`, `m1-erster-lauf`) | „Rohdaten der Läufe …“ mit Aufzählung | deckt |
| `docs/messungen/README.md` | Zeile „update“ | deckt |
| Slice-Datei | „Plan-Tabelle, DoD-Häkchen, §1, Gelesen, Zustandsmarker in §6 und §8“ | deckt |
| der Review-Report [`2026-11-05-slice-v1-devcontainer-remove-messung`](2026-11-05-slice-v1-devcontainer-remove-messung.md) (neu, 80 Zeilen, in `94f65083`) | **keine Zeile** | **Abweichung V-4** |
| dieser Verification-Report (neu) | **keine Zeile** | kommt zu V-4 hinzu |
| `docs/plan/planning/open/README.md`, `in-progress/roadmap.md` | „unverändert“ mit gelesenen Fundstellen | im Diff nicht berührt, Fundstellen nachgelesen: deckt |

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | **pass (eigener Lauf)** | Im Devcontainer-Container (`/.dockerenv`, Docker 29.9.3), Exit 0 auf `94f65083` (Arbeitsbaum sauber): `tool-pruefung` OK (go 1.27.3, golangci-lint 2.14.0), `gate-reihenfolge` OK (2 Fälle), `baseline-verify` v6.18.1 OK (54 Dateien), `docs-check` 1160 Dateien / 0 Befunde (`pinsync` 3 Pins), Lint und Go-Tests grün (Cache-Treffer sind inhaltsgleiche frühere Läufe; der Slice berührt keinen Go-Code), `coverage-gate` 92,30 % (Schwelle 90 %), `ablage-pruefung` OK, `ablage-pruefung-test` 26 Fälle / 0 Fehlschläge. |
| `make ablage-pruefung` | **pass** | Teil von `make gates` (oben): Skripte nur unter `tools/`, nur `.md`/`.txt` unter `docs/messungen/` |
| `make docs-check` mit diesem Report | **pass** | 1161 Dateien, 0 Befunde (siehe „Nachtrag Sensor“) |
| `make test-docker` | **nicht gefahren** | Der Diff berührt keinen Go-Code und keinen Test; die Messskripte sind keine Tests und laufen nicht unter `make`. Keine Docker-Messung wiederholt (Auftrag). |

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Messprotokoll M1 bis M5 mit Rohdaten (Verhalten von `down`, `up --recreate`, Instanzen, Exit-Code 10 ohne `--instance`); bestehender Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` pinnt die Isolation und `--recreate` (nicht von diesem Slice) |
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | M4: Volume des Features `ollama` (Name, Ziel `/home/vscode/.ollama`, Überleben, Prüfsumme des öffentlichen Schlüssels) |
| [`ADR-0015`](../plan/adr/0015-sandbox-volumes-pro-instanz.md) | M1/M3: Namensfamilien und Ende je Container, `gemeinsame=0` |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: [`CO-001`](../plan/carveouts/CO-001-crlf-logs.md) (Diff berührt `docs/plan/carveouts/` nicht)

Nicht ausgeführt:
- `make test-docker` — kein Go-Code, kein Test im Diff (Begründung oben).
- `make verify-depguard` — keine Importregel berührt.
- `make doc-immutable RANGE=…` — kein ADR im Diff (`git diff --stat`: `docs/plan/adr/` nicht berührt).
- `make ci`, `make fullbuild` — kein Security-/Release-Pfad.
- Wiederholung der Messungen M1 bis M5 — Auftrag: keine Docker-Messungen wiederholen; die Aussagen sind gegen die archivierten Rohdaten gelesen, nicht gegen einen neuen Lauf. Folge: Ob die Rohdaten die Skripte des Tips erzeugt haben, ist nur durch Zeitstempel und Diff belegt (Änderungszeit von `lib.sh`, `03`, `04b`, `00` 13:25 vor dem Lauf 13:25:42; `05-aufraeumen.sh` danach nur in Kommentarzeilen geändert, `git diff 8861376d..HEAD`), nicht durch einen Wiederholungslauf.
- Native Engine, `arm64`, andere Docker-Versionen — nicht Gegenstand (§1 Nicht-Ziele).

Commit / Artefakt:
- `94f65083` (Tip); Messung `8861376d`; Basis `642aefe7`

## Abweichungen und Hinweise

| ID | Gewicht | Befund | Fundstelle |
|---|---|---|---|
| V-1 | LOW (Vorbehalt Haken 2) | Unter M4 „Darf tragen“ steht „Gemessen ist der Inhalt: Markerdatei, Schlüsselpaar `id_ed25519`/`id_ed25519.pub` und die leeren Verzeichnisse `models` und `cache`.“ Die Rohdaten tragen „leer“ nicht: `ls -la` zeigt `models` mit **drei** Links (`drwxr-xr-x 3 …  models`), das Verzeichnis enthält also mindestens ein Unterverzeichnis; nur `cache` (zwei Links) ist als leer lesbar, und auch das ist eine Folgerung aus der Linkzahl, keine Messung des Inhalts. Dieselbe Klasse wie Review F-5 (Zusage-Satz nennt mehr, als gemessen). Entweder „die Verzeichnisse `models` und `cache`“ (ohne „leer“) oder die Messung nachholen. | Protokoll M4, „Darf tragen“, zweiter Spiegelstrich; `rohdaten-m4.txt`, Zeilen 37 bis 44 |
| V-2 | LOW (Vorbehalt Haken 2) | M1 „Darf tragen“, letzter Spiegelstrich: „Weder die Labels des Containers noch die der benannten Volumes tragen die `devcontainerId` **oder einen Verweis zwischen Container und Volume**.“ Gemessen ist nur, dass das 52-stellige Namensende weder in den Label-Werten noch in der Umgebung des Containers vorkommt (`id-in-labels=nein`, `id-in-umgebung=nein`) und dass die Volumes keine Labels haben (`labels=null`). Von den Labels des Containers sind die **Namen** archiviert, nicht die Werte; darunter `devcontainer.metadata`; die erzeugte `devcontainer.json` enthält `source=<projekt>-workspace-${devcontainerId}` (gelesen in den Rohdaten, Zeilen `json:`), und ob die Dev-Container-CLI solche Konfigurationsteile in den Wert dieses Labels schreibt, steht in keiner Rohdatei. Dieser Verifier hat es nicht gemessen (kein Docker-Lauf) und behauptet es nicht; die Aussage „kein Verweis“ ist damit ungeprüft. Auf „nennen das aufgelöste Namensende nicht“ einengen oder den Wert des Labels in die Rohdaten aufnehmen. Die Folgerung des CR (Erkennung über `Mounts`) bleibt davon unberührt. | Protokoll M1, „Darf tragen“ letzter Spiegelstrich; `rohdaten-m1.txt` (Zeilen `Labels des Containers (Namen)`) |
| V-3 | LOW | Review F-6 ist nur für M1 berichtigt. In M3 „Darf tragen“ („Im Standardprofil hat auch eine Instanz keine Volumes …“) und in den Kurzfassungen M1 („Das Standardprofil hält keins.“) und M3 steht dieselbe Verallgemeinerung ohne „ohne Katalog-Feature“; gemessen ist das Standardprofil ohne Katalog-Feature (`std`, `anon`). Der Abschnitt „Grenzen“ trägt die Einschränkung allgemein, die Sätze selbst nicht. | Protokoll Kurzfassung M1/M3; M3 „Darf tragen“, vierter Spiegelstrich |
| V-4 | LOW | Plan-Tabelle gegen den Diff nach der letzten Fix-Runde: Der Review-Report (neu in `94f65083`) hat keine Zeile in §3; dieser Report ist die nächste Datei ohne Zeile. Beide ergänzen („Review-Report und Verification-Report, neu“) und den Abgleich nach der Closure-Änderung wiederholen ([`AGENTS.md`](../../AGENTS.md) §3.10). Dieselbe Klasse wie [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Rückfall einer verkörperten Regel). | Slice §3 Plan-Tabelle |
| V-5 | LOW (Entscheidung nötig: Planner/Projektinhaber) | Die Fix-Runde hat zu Review F-1 die **Namen** des Fremdinventars in fünf Rohdateien und in das Protokoll geschrieben (7 Container: `boring_noyce`, `nice_fermi`, `pedantic_spence`, `relaxed_archimedes`, `stoic_aryabhata`, `stupefied_elion`, `unruffled_perlman`; 2 Volumes `dcsdmn72acy22yw-workspace-…`, `dcsdmn7vfi1bhvb-workspace-…`). Für die Aussage „unberührt“ genügt die leere Differenz plus Zahl; die Namen sind nicht nötig. Der Vorgänger [`slice-v1-devcontainer-ollama-amd64-messung`](../plan/planning/done/slice-v1-devcontainer-ollama-amd64-messung.md) hat dieselbe Klasse (Inventar des Messhosts in Rohdaten, Review F-3, MEDIUM) so behoben, dass `lib.sh` nur Anzahlen und Präfix-Objekte ausgibt; Register: [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (Klasse „Fremdinventar statt Zugangsdaten“). Gelesen: die Namen sind von Docker vergebene Zufallsnamen und ein Volume-Name aus Zufallszeichen, ohne Projekt- oder Personenbezug; **die Rohdaten tragen keine Zugangsdaten** (Suche leer). Die Commits `94f65083` und `8861376d` sind **nicht gepusht**, ein Umschreiben vor dem Push ist möglich. Zu entscheiden: Namen behalten (bewusste Abweichung, in der Closure-Notiz ausweisen) oder `inventur_zeigen` auf Zahl und Prüfsumme der sortierten Namensliste umstellen. | `tools/messung/<Kennung des Slice>/lib.sh` (`inventur_zeigen`); Protokoll „Umgebung“; `rohdaten-m1.txt` Zeilen 5 bis 14 |
| V-6 | INFO (Entscheidung Planner) | Die Fix-Runde `94f65083` ist nicht nachgeprüft worden: neues Skript `04b-m5-down-ohne-instanz.sh`, neue Rohdaten M5, alle Rohdaten M1 bis M4 und `rohdaten-aufraeumen.txt` neu gelaufen (13:25 bis 13:29), 133 geänderte Zeilen im Protokoll mit neuen „Darf tragen“-Sätzen, neue Frage an den Architect. Der Verifier hat M5 und die geänderten Sätze gegen die Rohdaten gelesen (oben); V-1 bis V-3 und V-5 liegen genau in diesen Änderungen. Ob ein Nachlauf-Review nötig ist (Vorbild: `…-nachlauf`-Reports in diesem Verzeichnis), entscheidet der Planner. | `git diff --stat 8861376d..94f65083` |
| V-7 | INFO | Ursachen- und Zustandsaussagen im Protokoll ohne Messstelle, keine davon unter „Darf tragen“: (a) „(er sank; andere Last auf dem Host), weil der Build-Cache der Engine die Schichten hält“ (Platz; gemessen ist nur die Differenz −1.937.408 Byte; Slice §6 Risiko 2 wiederholt das „weil“); (b) „legt `ollama serve` beim Start das Schlüsselpaar … an“ (das Listing folgt auf `ollama --version` **und** `ollama serve`; nichts davor); (c) „das Workspace-Volume war neu und leer“ (gemessen: die Markerdatei fehlt; „neu“ folgt aus dem vorherigen Entfernen); (d) „Images dieser Projekte neu, Basis-Image lokal, Schichten teils im Build-Cache“ für den ersten Lauf (Cache-Zustand nicht festgehalten). Als Vermutung kennzeichnen oder streichen ([`AGENTS.md`](../../AGENTS.md) §3.10). | Protokoll M1 „Dauer und Platz“, M4; Slice §6 |
| V-8 | INFO | Slice §8 schließt: „Die übrigen Titel betreffen diesen Slice nicht“, gelesen sind nur die Titel. Mindestens ein Titel betrifft Messslices und ist für V-5 einschlägig: [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (zweiter Beleg: Fremdinventar des Messhosts). Weitere Titel mit Messbezug ([`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../plan/planning/observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md): die Messung fährt den ausgelieferten Befehl, so im Protokoll ausgewiesen; [`BEO-DOK/messprotokoll-ausserhalb-des-repos`](../plan/planning/observations/BEO-DOK/messprotokoll-ausserhalb-des-repos/observation.md): Protokoll und Skripte im Repo) sind eingehalten. Der Satz in §8 auf „gelesen: Titel; einschlägig: …“ ergänzen oder die Beobachtung ergänzen. | Slice §8 |
| V-9 | INFO | Rohdaten und Skripte: `rohdaten-aufraeumen.txt` stammt von einer früheren Fassung von `05-aufraeumen.sh`; die Datei wurde um 13:32 nur in Kommentarzeilen geändert (`git diff 8861376d..HEAD -- …/05-aufraeumen.sh`), der Code ist gleich. Die Register-Datei der Volumes überdauert die Läufe von M1 bis M5 (die Ausgabe in `rohdaten-m2.txt` listet auch Volumes aus M1 als „entfernt“); das ist für die Aussagen unschädlich, weil nur das Präfix und exakt registrierte Namen gelöscht werden (`objekte_entfernen`, Muster mit Bindestrich). | `lib.sh`, `05-aufraeumen.sh`, `rohdaten-m2.txt` Zeilen 156 bis 175 |
| V-10 | INFO | Offene Messpunkte (kein Mangel, aber Rest für die Closure-Notiz): native Engine (der Aufbau ist die verschachtelte Engine des Devcontainers); `docker volume prune`/`system prune`/`compose down -v`; `paused`, `restarting`, `created`; drei und mehr Instanzen; Instanzen mit `nestedRuntime: podman`; Katalog-Features außer `ollama` im Standardprofil und Features aus externen Quellen; von Hand geschriebene `mounts`; Neubau des Images; Umzug des Projektordners; `--recreate` für Docker-Speicher und Modelle; Inhalt von Docker-Speicher und Modellen zwischen Instanzen; ein geladenes Modell; `shell` ohne `--instance`; M5 im Standardprofil und mit genau einer Instanz; gestoppter Container ohne Instanz neben Instanzen. | Protokoll „Grenzen“ und alle „Darf nicht tragen“ |

## Beobachtungs-Register

Der Verifier schreibt nicht ins Register; Hinweise für die Closure-Notiz:

| Beobachtung | Befund an diesem Slice |
|---|---|
| [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md) | weiterer Beleg: Review F-4, F-5, F-6; in der Fix-Runde berichtigt, mit Rest V-1, V-2, V-3 (derselbe Slice, keine neue Klasse: „Satz unter ‚darf tragen‘ breiter als die Messung“). Die Beobachtung steht über der Schwelle, Ausgang beim Projektinhaber. |
| [`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../plan/planning/observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md) | Stand: 2 Belege. Review F-1 (Zahl der unberührten Fremdobjekte ohne Messstelle) ist ein dritter Beleg; die Behebung folgt dem Kandidaten der Beobachtung (Zustand vor und nach dem Lauf lesen, in die Rohdaten schreiben, Rückabwicklung gegen den Lesestand bestätigen) und **hielt**: Differenz leer in fünf Rohdateien. Mit dem dritten Beleg ist die Schwelle erreicht; Entscheidung des Planners. Spannung zu V-5. |
| [`BEO-WZG/aufraeum-skript-ueber-praefix-hinaus`](../plan/planning/observations/BEO-WZG/aufraeum-skript-ueber-praefix-hinaus/observation.md) | Der Rest der Beobachtung (Namensmuster ohne abschließendes Trennzeichen) ist in diesem Slice umgesetzt: `lib.sh` löscht mit `^(u-boot-)?ubmess-rm-`, `^ubmess-rm-`, `^vsc-ubmess-rm-`; `anfang` meldet Reste mit dem Präfix als `BEFUND`. Kein neues Auftreten; Beleg „befolgt“. |
| [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | Rückfall einer in [`AGENTS.md`](../../AGENTS.md) §3.10 verkörperten Regel: Review F-2 (Aussage „Kein Code des Repos“), F-3 (Zustandsmarker in §8), F-7 (Beschreibung ohne Fehlerzweig), V-4 (Plan-Tabelle ohne Zeile für den Report). |
| [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) | V-5: Fremdinventar-Namen in den Rohdaten, bewusst geschrieben; Zugangsdaten: keine. |

## Was die Closure-Notiz tragen muss (Planner)

1. **Risiko-Ausgänge §6** (beide stehen auf „offen“):
   - Risiko 1 (verschachtelte Engine weicht von einer lokalen ab): Ausgang **offen mit Grenze, weitergereicht** — der Aufbau steht im Protokoll (Docker 29.9.3, `overlayfs`, cgroup v2, verschachtelt, `u-boot` behandelt die Engine als lokal); eine Messung auf einer nativen Engine ist nicht beauftragt. Ein CR-Wortlaut darf nur die Zusage-Menge von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) tragen (Linux-`amd64`, lokale Engine); die Abweichung ist Grenze, nicht Mangel.
   - Risiko 2 (Bauzeit, Platz): Ausgang **gemessen, mit Grenze** — Bauzeiten 4 bis 8 s (Cache warm) und 5 bis 48 s (erster Lauf), Image-Größen 1,41 bis 2 GB bei 14 Images; der Platzbedarf ist nicht belastbar gemessen (V-7 a: die Ursache nicht als Messung führen).
   - Neue Risiken aus diesem Lauf, die einen Ausgang brauchen: Fix-Runde ohne Nachlauf-Review (V-6); Fremdinventar-Namen in unveröffentlichten Commits (V-5).
2. **Eingang für die Folge-Slices** (Verifier-Lesung, Entscheidung beim Planner und Projektinhaber): Gegen die Entscheidungen 2 bis 5 von [`slice-cr-devcontainer-remove`](../plan/planning/in-progress/slice-cr-devcontainer-remove.md) trägt die Messung den Vorschlag (keine Rückführung nach `open/` nötig nach dem Wortlaut von dessen §4): Volumes eines Containers sind vollständig über `Mounts` bestimmbar, aber **nur vor dem Entfernen des Containers** (M1, M4); ein gestoppter Container hält sein Volume bis zum Entfernen (M2c); `docker rm` scheitert bei laufendem Container ohne `-f` (M2a, Entscheidung 3); das Ollama-Volume ist ein benanntes Volume wie die anderen, nur über Namen (`-ollama-`) und Ziel (`/home/vscode/.ollama`) zu unterscheiden (M4, Entscheidung 4); ohne `--instance` gilt der Suchweg von `up`/`down`/`shell` (M5, Entscheidung 5). Die **Frage an den Architect** aus M5 („Welcher Container ist bei `remove` ohne `--instance` gemeint?“) steht nur im Protokoll; der CR-Slice nennt sie weder in §1 noch in §6 — in die Closure-Notiz als Übergabe aufnehmen.
3. **Rest-/Folgepunkte:** die Liste V-10; V-1 bis V-5 als Berichtigungen vor oder mit der Closure (Textänderungen, keine neue Messung, außer V-1/V-2, wenn der Satz bleiben soll).
4. **Steering-Loop-Lerneintrag:** (a) Rückfall von `plan-aussage-ueber-code` trotz verkörperter Regel (V-4, Review F-2, F-3, F-7); (b) `zusage-weiter-als-messung`: weiterer Beleg (Review F-4 bis F-6), in der Fix-Runde nur teilweise berichtigt (V-1, V-2, V-3), (c) `plan-annahme-ueber-fremdzustand-ohne-messstelle`: dritter Beleg, Wirksamkeit der Behebung belegt (Inventur vor/nach in den Rohdaten), (d) Kandidat für den Lerneintrag der Spannung: Eine Inventur, die „unberührt“ belegt, trägt **Zahl und Prüfsumme der sortierten Namensliste** statt der Namen, damit der Beleg nicht zugleich Fremdinventar in den Verlauf schreibt (V-5; Vorgänger: [`slice-v1-devcontainer-ollama-amd64-messung`](../plan/planning/done/slice-v1-devcontainer-ollama-amd64-messung.md) Review F-3).
5. **Reihenfolge der Closure** ([`AGENTS.md`](../../AGENTS.md) §3.3): erst der Inhalt (Haken Review und Verification, Closure-Notiz, Plan-Tabelle mit Review- und Verification-Report, Register-Evidence, Nachführung von `roadmap.md` Zeile 45 und `open/README.md` Zeile 23), danach der reine `git mv` nach `done/`; der Zustandsmarker in Roadmap und `open/README.md` wechselt mit der Closure.

## Urteil

**Verifiziert mit Auflagen; Freigabe der Closure vorbehaltlich V-1 bis V-5.** Die Messung hält der Prüfung stand: Jede Zahl und jeder Name des Protokolls steht in einer Rohdatei mit Befehl, Ausgabe und Zeitstempel; die Läufe-Herkunft (archivierte Fassung gegen `8861376d` und gegen den ersten Lauf) ist nachgeprüft und stimmt mit dem Protokoll überein; M5 ist durch `rohdaten-m5.txt` gedeckt und stimmt mit [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und dem gelesenen Code überein; die Aussagen über Code stimmen; kein Gate und kein Produktionscode ist berührt. `make gates` (eigener Lauf im Devcontainer-Container, Exit 0) ist grün. Die DoD-Haken „Messprotokoll“ und „`make gates` grün“ stehen zu Recht; der Haken „Eingang für den CR“ steht **vorbehaltlich** der Berichtigung der drei Sätze V-1, V-2 und V-3 (Textänderung auf das Gemessene, genau die Klasse, die der Reviewer als F-4 bis F-6 gefunden hat — der CR soll sich auf diese Sätze stützen). Die Review-Findings F-1 bis F-4, F-7, F-9 sind behoben, F-5 und F-6 nur teilweise (V-1, V-3), F-8 bleibt INFO, F-10 ist mit diesem Report übernommen.

Vor dem Verschieben nach `done/` fehlen:

1. **V-1, V-2, V-3:** Sätze unter „Darf tragen“ auf das Gemessene einengen (oder messen).
2. **V-4:** Plan-Tabelle um Review- und Verification-Report ergänzen und nach der Closure-Änderung neu gegen `git diff --stat` abgleichen.
3. **V-5:** Entscheidung zu den Namen des Fremdinventars (behalten und ausweisen oder auf Zahl und Prüfsumme umstellen; vor dem Push möglich).
4. **V-6:** Entscheidung, ob ein Nachlauf-Review der Fix-Runde nötig ist.
5. Haken Review und Verification, Closure-Notiz mit Ausgängen zu §6 und Lerneintrag (Abschnitt oben), Register-Evidence, Nachführung von Roadmap und `open/README.md`; danach der reine `git mv`.

Eine Rückführung nach `next/` oder `open/` (§4) ist nicht nötig: Die Engine im Devcontainer-Container stand zur Verfügung, und keine Messfrage brauchte ein Gerät, das der Devcontainer nicht stellt.

## Nachtrag Sensor

`make gates` lief vor Anlage dieses Reports (Stand `94f65083`, 1160 Dateien); die Datei selbst ist neu und im Lauf nicht enthalten. `make docs-check` mit diesem Report: 1161 Dateien, 0 Befunde (die erste Fassung trug zwei ungelinkte Kennungen und wurde berichtigt). `make ablage-pruefung` nach Anlage des Reports: OK, 26 Fälle / 0 Fehlschläge.
