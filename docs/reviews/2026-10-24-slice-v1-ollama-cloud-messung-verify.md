# Verification-Report: slice-v1-ollama-cloud-messung

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff e60956c3..HEAD` (HEAD `952cb06c`), Slice
  [`slice-v1-ollama-cloud-messung`](../plan/planning/done/slice-v1-ollama-cloud-messung.md)
  (Welle [`welle-ollama-cloud-modelle`](../plan/planning/welle-ollama-cloud-modelle.md))
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
* Review: [`2026-10-24-slice-v1-ollama-cloud-messung.md`](2026-10-24-slice-v1-ollama-cloud-messung.md) (HIGH-1, MEDIUM-1, LOW und INFO in `952cb06c` behoben)
* Protokoll: [`README.md`](../messungen/slice-v1-ollama-cloud-messung/README.md)

## Urteil

**Schließbar (Messung und Verifikation), unter einer Planer-Bedingung und den noch offenen Closure-Pflichten.** Alle Sensoren, die ich fuhr, sind grün;
die rote Gegenprobe der Schwärzung ist reproduziert; die geprüften Aussagen des Protokolls decken sich mit den Rohdaten; keine Zugangsdaten im
Slice-Bestand, auch nicht im Verlauf des Diffs; Aufräumen belegt (0 Objekte). Die gekennzeichnete Planänderung DoD 3 lasse ich **gelten** (Begründung unten).
Bedingung vor `done/`: V-1 (Wortlaut §5 an DoD 3 angleichen, Entscheidung des Planners). Offen sind nur Closure-Pflichten, die nicht mir gehören
(DoD-Häkchen Review und Doku-Update, Closure-Notiz, Register, Risiken-Ausgänge, Paarungen).

## Sensoren (selbst gefahren, Exit-Codes sichtbar)

| Sensor | Ergebnis |
|---|---|
| `make gates` | `GATES_EXIT=0` (Log bis `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `d-check: 811 Datei(en) geprüft, 0 Befund(e)`); Baum danach unverändert (0 Zeilen `git status --short`) |
| `make ablage-pruefung` | Exit 0 |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen |
| `make doc-planning` | Exit 0, 0 Befunde |
| `make doc-commits RANGE=26bfd84e..HEAD` | Exit 0, 11 Commits, 0 Befunde; zusätzlich `RANGE=origin/main..HEAD` (Wortlaut der DoD 4) Exit 0, 0 Befunde |
| `redact-probe.sh` (Docker, wie vorgesehen) | Exit 0: „24 Formen geschwaerzt, neutrale Zeilen unveraendert“ |
| rote Gegenprobe | Patch (`git apply`, entfernt die Regel für Anmelde-URLs aus `redact()` in `lib.sh`): Exit 1, drei Zeilen `DURCHGELASSEN` (`connect?name=…&key=…`, `connect?name=host&key=…`, `signin#…`); Rücknahme mit `git apply -R`; Baum vorher und nachher sauber (0 Zeilen `git status --short`, 0 Zeilen `git diff --stat`). Deckt sich mit der Angabe im Protokoll (drei Zeilen, Exit 1) |
| `06-suche.sh` über Messordner plus `tools/messung/ollama-cloud/` | Exit 0, „keine Treffer“ |
| `06-suche.sh` Vorgabe (ganz `docs/messungen/` und `tools/messung/`) | Exit 1; Treffer gelesen, siehe DoD 3 |
| eigene Suche | Präfixe (`ghp_`, `github_pat_`, `glpat-`, `sk-ant-`, `sk-…`, `ollama_…`), `Bearer`, `Authorization`, `BEGIN … PRIVATE`, `ssh-ed25519`/`ssh-rsa`/`AAAAC3Nza`, `id_ed25519`-Inhalt, `connect?`/`signin?`-Query, `code=`/`key=`/`nonce=`/`token=`/`name=`, URL-Zugang, `eyJ…`, `AKIA…`, Zeichenketten ab 40 Zeichen über Messordner, Skripte und Review-Report: nur Hashes von Images, Volumes und Containern, Digests in `Dockerfile.cli`; Anmelde-URLs stehen als `https://ollama.com/connect?[Query entfernt: Anmelde-URL]`. Zusätzlich über **alle hinzugefügten Zeilen** `git log -p e60956c3..HEAD -- docs tools` (Verlauf, nicht nur Endstand): nur Muster in `06-suche.sh`/`lib.sh` und Platzhalter `PROBEWERT` in `redact-probe.sh` |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Produktcode (`internal/`, `cmd/` nicht im Diff), keine Importregel berührt; Speicher des Nachbar-Repos |
| `make ci`, `make fullbuild` | nicht gefahren; nicht Teil der DoD dieses Slice |
| Messung wiederholt | nicht gefahren (Auftrag: Speicher); keine Gegenprobe-Läufe nötig, die Rohdaten tragen zwei Läufe je Modus |

Hinweis zu meinem Vorgehen: zum Lesen habe ich in Bash einzelne `sed -n`/`awk`-Anzeigen auf Dateien benutzt (nur lesend, kein Schreiben, kein
Edit); alle Schreibvorgänge und der Patch liefen über `Write` und `git apply`.

## DoD, Punkt für Punkt

| # | Zusage | Befund | Beleg |
|---|---|---|---|
| 1 | C0 bis C7 am erzeugten Projekt, je Aussage Umgebung, Datum, Läufe, bedingungsgleich | **erfüllt** | Projekte aus `u-boot generate devcontainer` (`rohdaten-02.txt`, `-02-erzeugt-*`); Pin `OLLAMA_VERSION=0.40.6` = ausgeliefert, C0 ohne Abweichung; Messstand und Datum im Kopf des Protokolls; jede Aussagenzeile trägt Rohdaten und Läufe; Formulierung „ohne Konto gemessen, Erfolgspfad nicht gemessen“ im Kopf und in den Tabellen; Handvarianten (CLI statt `u-boot devcontainer up`, `discover`-Zweitcontainer, `docker exec`, Platzhalter `claude`, Pseudo-Terminal, Lauf 2 mit Volume, ein Modellname) ausgewiesen |
| 2 | Messplan und Liste „nicht gemessen“ | **erfüllt** | Tabelle *Nicht gemessen* mit Aufruf und Bedingung: Anmeldung, Lauf eines Cloud-Modells, Hosts nach der Anmeldung, Zusammenspiel mit Claude Code, Zugangsvariable, interaktiver Ablauf, `OLLAMA_NO_CLOUD` mit Konto, andere Modelle/Versionen, andere Plattformen, Kaltstart Lauf 2. Vollständig gegenüber §1 (die fünf genannten Punkte sind enthalten, dazu fünf weitere); kein Wortlaut für das Lastenheft |
| 3 | Keine Zugangsdaten | **erfüllt mit gekennzeichneter Planänderung (gelten lassen)** | siehe unten |
| 4 | `make gates`, `docs-check`, `ablage-pruefung`, `doc-commits` | **erfüllt, von mir wiederholt** | Tabelle Sensoren |
| 5 bis 9 (Review, Doku-Update, Closure-Notiz, Register, Risiken-Ausgänge, Paarungen) | **offen, nicht meine Pflicht**; Review-Report liegt vor, Häkchen steht aus | §7 trägt noch den Platzhalter „Der Slice ist noch nicht gestartet“, `observations/` ist nicht fortgeschrieben (Plan-Tabelle nennt beides für die Closure) |

### Entscheidung zur Planänderung DoD 3: gelten lassen

Die Fassung vom Plan-Commit `e60956c3` verlangte „Suche über `docs/messungen/` und `tools/messung/` leer“. Die Fassung im Slice verlangt sie leer über das Verzeichnis dieser Messung
und `tools/messung/ollama-cloud/` und liest die Treffer über den Gesamtbestand von Hand.

**Ich habe die älteren Treffer selbst gelesen** (`06-suche.sh` ohne Argument, Exit 1, 36 Zeilen, fünf Muster):

* `ghp_…|ollama_…`: vier Zeilen in zwei Rohdaten des [Protokolls amd64](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md) (`rohdaten-09.txt`, `rohdaten-lauf2-09.txt`): Go-Testname `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`, kein Token.
* `ssh-ed25519 AAAA…`: vier Zeilen im [Protokoll Ollama-Messung](../messungen/slice-v1-devcontainer-ollama-messung/README.md) (`rohdaten-05-restricted-1/2/3.txt`, `rohdaten-05-discover.txt`), jeweils direkt hinter „Your new public key is:“ (gelesen im Kontext): öffentlicher Schlüssel eines Wegwerf-Containers, kein privater Teil.
* `://…:…@`/`…TOKEN=`: Prosa im [Protokoll Forge-CLI](../messungen/slice-v1-forge-cli-quellen-messung/README.md) (Beschreibung der Schwärzungsformen) sowie Zuweisungen von `$DUMMY` (`MESSUNG-KEIN-ECHTES-TOKEN`, `tools/messung/forge-cli/lib.sh` und `inner-07.sh`) und der Wert `platzhalter-kein-zugangsdatum-0123456789abcdef` in `tools/messung/devcontainer-start-befehl/01-einzelmessung.sh`.
* `key=…`: drei GPG-Fingerabdrücke (`GITHUB_CLI_ARCHIVE_GPG_KEY`, `PODMAN_SIGNING_KEY`), öffentliche Kennungen von Signierschlüsseln.

Keiner der Treffer ist ein Zugangsdatum. Gründe, die Änderung gelten zu lassen:

1. Die ursprüngliche Fassung ist über den Gesamtbestand **nicht erfüllbar**, ohne fremde, abgeschlossene Protokolle zu verändern; dieser Slice berührt sie nicht (im Diff stehen nur `docs/messungen/README.md` und das eigene Verzeichnis).
2. Der Zweck der Zeile (dieser Slice bringt kein Zugangsdatum in das Repo) ist mit der engeren Suche **vollständig** getragen: leer über den eigenen Bestand, von mir wiederholt, plus Suche über alle neu hinzugefügten Zeilen im Verlauf und eigene Muster, die über die des Skripts hinausgehen.
3. Die Abschwächung ist nicht still: DoD 3 kennzeichnet sie im Wortlaut, nennt Fassung und Grund, und das Protokoll weist die Treffer nach Art aus.

Eingeschränkt gilt: Es ist eine Abschwächung des Maßstabs durch die Implementer-Rolle; sie wird mit **dieser** Entscheidung des Verifiers wirksam, nicht vorher. Die Planer-Rolle sollte sie
bestätigen und das Kennzeichen („vom Verifier gelten zu lassen oder zu verwerfen“) beim Schließen durch den Befund ersetzen (V-1).

## Stichprobe Protokoll gegen Rohdaten

Gelesen gegen `rohdaten-04-{discover,restricted,allow}-1/2.txt`, `rohdaten-02*.txt`, `rohdaten-01`, `rohdaten-05*`; 16 Aussagen, alle **bestätigt**:

| # | Aussage des Protokolls | Rohdaten |
|---|---|---|
| 1 | `restricted`, `ollama pull`: Exit 0, nur `registry.ollama.ai` (4 Anfragen in Lauf 1, 2 in `C5-pull`), kein `ollama.com`, kein Blob-Host; Dateien 264, 307, 426 Byte | `restricted-1` Z. 169–183; `discover-1` Z. 133–146 |
| 2 | Start von `ollama serve`: `ollama.com` 2 Anfragen (`discover`, `allow`), 4 (`restricted`); mit `OLLAMA_NO_CLOUD=1` keine Anfrage | `discover-1` Z. 98, `allow-1` Z. 138, `restricted-1` Z. 135, 252 (leer) |
| 3 | `run` in `restricted`: `Error: Post "https://ollama.com:443/api/show?ts=…": dial tcp: lookup ollama.com on 127.0.0.1:53: server misbehaving`, Exit 1, 4 Anfragen, `config error` ×4; Lauf 2 gleich | `restricted-1` Z. 184–191, `restricted-2` Z. 181–183 |
| 4 | `run` in `discover`/`allow`: `You need to be signed in to Ollama to run Cloud models.`, Exit 0, 0,2 s | `discover-1` Z. 147–150, `allow-1` Z. 188–191 |
| 5 | `launch claude` in `discover`/`allow`: `Error: gpt-oss:120b-cloud requires sign in`, Exit 1, 0,2 s; Platzhalter in `allow` ebenso Exit 1, kein `STUB` | `discover-1` Z. 157–160, `allow-1` Z. 199–217 |
| 6 | `launch claude` in `restricted`: 2 × `api.anthropic.com`, 12 × `ollama.com` abgelehnt, `Error: Input must be provided …`, `Error: exit status 1`, 0,7 s (Lauf 2: 0,8 s) | `restricted-1` Z. 197–206, `restricted-2` Z. 197 |
| 7 | `signin` in `restricted`: `Error: account unavailable`, Exit 1, 4 abgelehnte Anfragen; `signout`: `Error: there was an error signing out`, Exit 1 | `restricted-1` Z. 232–249 |
| 8 | `signin` in `discover`/`allow`: Anmeldehinweis, Exit 0, `signout`: `You are not signed in to ollama.com`, Exit 0; unter `OLLAMA_NO_CLOUD=1` 2 Anfragen `ollama.com` | `discover-1` Z. 171–182, 214–221; `allow-1` Z. 223–235, 271–278 |
| 9 | `OLLAMA_NO_CLOUD=1`: `run` → `Error: ollama cloud is disabled: remote model details are unavailable` (keine Anfrage); `launch` → `Warning: ignoring --model …` und `Error: model selection requires an interactive terminal …`; Server-Log `Ollama cloud disabled: true` | `restricted-1` Z. 250–297 |
| 10 | Hostzählung `ollama.com` gegen `registry.ollama.ai` über die Sitzung: `discover`/`allow` 4 gegen 6, `restricted` 52 gegen 6 (Lauf 1) und 52 gegen 4 (Lauf 2, „4 gegen 6 je nach Lauf“ im Protokoll) | `discover-1` Z. 223–225, `allow-1` Z. 282–284, `restricted-1` Z. 291–294, `restricted-2` Z. 287–289 |
| 11 | Dateien Lauf 1: `id_ed25519` 387 und `.pub` 81 (vor `pull`), `cache/model-recommendations.json` 2370 nur in `discover`/`allow`, `config.json` 104 nur in `restricted`, `~/.claude.json` 852 und `backups/…` 84 auf Wurzel-Dateisystem (`/`), alles andere im Volume `/home/vscode/.ollama` | `restricted-1` Z. 167–168, 208–210; `discover-1` Z. 130–132, 229–238 |
| 12 | `AUSSAGE`-Zeilen: Lauf 1 gegen 2 unterscheiden sich in allen drei Paaren **nur** im Exit-Code des Pseudo-Terminal-Aufrufs (124/137, wechselnd) | `diff` der `AUSSAGE`-Zeilen, selbst gefahren |
| 13 | Stub: sieben `CLAUDE_CODE_*`, sechs `ANTHROPIC_*`, `ENABLE_CLAUDEAI_MCP_SERVERS`; `ANTHROPIC_BASE_URL=http://127.0.0.1:11434` | `restricted-1` Z. 221–223 |
| 14 | C7 Standardprofil: `Updated devcontainer managed block`, zwei Warnungen ([`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) `claude-code`, [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) `ollama` ohne Wirkung), Exit 0; `devcontainer.json` 344 Byte, `Dockerfile` 271 Byte, Ollama-Block 0; zweiter Aufruf „already up to date“ plus beide Warnungen | `rohdaten-02.txt` Z. 11–27, 129–132; `rohdaten-02-c7.txt` |
| 15 | Egress-Liste `res`: 11 Namen, `ollama.com` nicht darin, Blob-Host `dd20bb89…r2.cloudflarestorage.com` und `registry.ollama.ai` als `feature:ollama`; Variante `alw`: `ollama.com` mit Herkunft `user`, als `server=`/`nftset=` im Container | `rohdaten-02.txt` Z. 79–117; `allow-1` Z. 27–29 |
| 16 | Umgebung: Namen `HOME OUT PATH PWD ROOT SHLVL UB_BEREINIGT USER WORK _`, Gegenprobe der Prüfung schlägt an (Exit 1, fünf Namen); Docker 29.9.2 `amd64`; Aufräumen `objekte-mit-praefix-nachher=0`, `container=0` | `rohdaten-01.txt`, `rohdaten-05.txt` Z. 35, `rohdaten-05-nachher.txt` |

Aufräumen zusätzlich selbst geprüft: `docker ps -a`, `docker images`, `docker volume ls`, `docker network ls` mit Präfix `u-boot-messung-ocloud` (und `vsc-`): jeweils 0 Objekte.

## Weitere Prüfungen

* **Zusage nur Gemessenes:** Das Protokoll liefert Messwerte und eine Liste des Ungemessenen, keinen Lastenheft-Wortlaut; „behebt den Fehler“ ist auf „die Meldung ändert sich“ begrenzt (C2b); die Standardformel steht im Kopf. Kein Satz „Cloud-Modelle laufen/laufen nicht“. Der Satz zur Standardprofil-Warnung (C7) ist auf die gemessene Ausgabe beschränkt.
* **Handvarianten:** acht Zeilen ausgewiesen; die weitere Verfahrensabweichung (zwei Host-`sed`-Aufrufe des Messenden außerhalb der Skripte) steht im Protokoll; Skripte rufen auf dem Host nur `docker`, `make`, `bash`, Coreutils, `grep`, `diff` (gegrept: kein `sed`, `awk`, `perl`, `python`, `jq`, `curl`, `sha256sum` im Skriptkopf).
* **Plan-Tabelle gegen `git diff --stat`:** `tools/messung/ollama-cloud/` (12 Dateien), das [Messverzeichnis](../messungen/slice-v1-ollama-cloud-messung/README.md) (Protokoll plus 22 Rohdaten), `docs/messungen/README.md`, Review-Report, Slice-Datei, Welle-Datei, `open/README.md`, `roadmap.md`: jede Diff-Zeile hat einen Plan-Eintrag; `internal/`, `cmd/`, `spec/`, `docs/user/`, `README*.md`, `CHANGELOG.md` im Diff nicht enthalten (Plan „unverändert“ stimmt). Offen laut Plan, noch nicht im Diff: Verify-Report (dieser), `observations/` (Closure).
* **Kommentare §3.7/§3.9:** `git diff` der Skripte und des Protokolls auf „früher/former/originally/no longer/ohne X wäre“ durchsucht: nur „früheren Messung“ (Abschnitt G, Gegenstand: die Handbuch-Messung, sachlich) und der Vermerk zum früheren Lauf der Suche in der Fix-Runde. Skriptkommentare tragen Zusage, Kopplung, Grenze; Verweise auf `lib.sh` und Dateinamen stimmen mit dem Bestand.
* **Anweisungen/Einstellungen in Dateien:** `claude-stub.sh` ist das einzige Treffer-Wort „claude“ im Diff nach Namen (Platzhalter-Skript, gibt Argumente und Variablen-**Namen** aus); keine Datei unter `.claude/`, `settings*`, `AGENTS.md`, `Makefile`, `.github/`, Hooks im Diff; keine Anweisung an Agenten im Text (Suche nach `settings.json`, „ignore“, „dangerously“, „disable hook/permission“: 0 Treffer). Die Warnung „settings-json“ des Harness ist ein Mustertreffer.
* **Sicherheitsrelevant, nur zur Kenntnis:** `Dockerfile.cli` und `04-lauf.sh` geben der Dev-Container-CLI den Docker-Socket des Hosts; das ist im Kommentar des Dockerfiles ausgewiesen („spricht über den Docker-Socket mit dem Daemon des Messhosts“) und für diese Messung gewollt (kein Host-Verzeichnis eines Agenten, `env -i`).

## Findings

### V-1 (LOW)

* `kategorie`: LOW (Plan-Konsistenz)
* `quelle`: `AGENTS.md` §3.10 (Zustandsmarker beim Ablauf berichtigen), Slice §2 DoD 3 gegen §5
* `pfad`: [der Slice](../plan/planning/done/slice-v1-ollama-cloud-messung.md), §5 *Closure-Trigger* („Suche über `docs/messungen/` und `tools/messung/` leer“) und DoD 3 (Kennzeichen „vom Verifier gelten zu lassen oder zu verwerfen“)
* `befund`: §5 trägt weiter die Fassung, die DoD 3 ersetzt, und ist mit dem Kriterium „DoD vollständig“ nicht erfüllbar (Suche über den Gesamtbestand: Exit 1). Das Kennzeichen in DoD 3 bleibt nach meiner Entscheidung stehen.
* `risiko`: Der Slice wird mit zwei widersprüchlichen Closure-Kriterien geschlossen; ein späterer Leser wertet §5 als nicht erfüllt.
* `empfehlung`: Planner gleicht §5 an DoD 3 an und ersetzt das Kennzeichen durch den Verweis auf diese Entscheidung (Abschnitt *Entscheidung zur Planänderung DoD 3* dieses Reports). Nicht von mir editiert.
* `verifizierbar`: Verifier (gelesen)

### V-2 (INFO)

* `befund`: Zwei Planänderungen sind nur teilweise gekennzeichnet. DoD 3 ist ausdrücklich gekennzeichnet; die Änderung des *Ansatzes* in §3 von „`env -u` der genannten Variablen“ zu „`env -i` mit Positivliste“ steht nicht als Planänderung da. Sie ist eine Verschärfung (strengere Umgebung) und im Protokoll belegt (Namen der Umgebung); ich lasse sie gelten, ohne Handlung.
* `verifizierbar`: Verifier (Diff gelesen)

### V-3 (INFO)

* `befund`: Die Gegenprobe der Schwärzung „vor dem ersten Lauf“ hat kein Artefakt (Konsole, im Protokoll als solche ausgewiesen); ebenso fehlt die Host-seitige Ausgabe von `umgebung_pruefen()` vor den Läufen (die Datei `rohdaten-01.txt` entstand nach den Läufen, im Protokoll ausgewiesen; Belege davor: `umgebung-container-verbotene-namen=0` in `rohdaten-04-*`). Ich habe die Gegenprobe heute grün und rot (Patch) gefahren und die Prüfung der Namen gegen die Skripte gelesen; die Reihenfolge „vor dem ersten Lauf“ selbst kann ich nicht belegen, nur die Fähigkeit des Werkzeugs.
* `empfehlung`: Im Folge-Slice Probe und Namensprüfung in eine Rohdatei schreiben (Zeitstempel vor dem ersten Lauf).
* `verifizierbar`: Verifier (Protokoll und Skripte gelesen)

## Übergabe an den Folge-Slice (Messung mit Konto)

Eingabe sind der Abschnitt *Nicht gemessen* des Protokolls und diese Punkte:

1. **Schwärzung vor dem ersten Lauf erweitern.** `redact()` ist formgebunden; die Anmelde-URL und ihr Code sind eine eigene Form (Probe hat sie, die Zusage gilt den Formen der Probe). Vor einem Lauf mit Konto eine Gegenprobe für die Form des Zugangs, die erst mit Konto sichtbar wird (Token von Ollama, Antwort-Inhalt, Kontodaten), mit Artefakt (V-3); das Register-Thema `BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten` steht bei 2 Belegen, ein Lauf mit Zugang wäre der dritte.
2. **Zugang nur als Umgebungsvariable des Hosts**, durchgereicht ohne Wert auf der Kommandozeile; `umgebung_pruefen()` blockt `OLLAMA*`, `TOKEN`, `API_KEY` — der Folge-Slice braucht eine bewusste, benannte Ausnahme für genau die eine Variable (Kandidat `OLLAMA_API_KEY` aus dem Binary, Wirkung ungemessen) und darf die Prüfung nicht abschalten.
3. **Offene Fragen aus den Rohdaten:** Warum `launch claude` ohne Konto in `restricted` das installierte `claude` startet (kein Anmeldehinweis), in `discover`/`allow` dagegen mit `requires sign in` abbricht (nur beobachtet, Ursache nicht gemessen); Ausgabe und Ende des Pseudo-Terminal-Aufrufs (0 oder 8 Byte, Exit 124/137); Wirkung von `OLLAMA_API_KEY`; Hosts und Dateien nach der Anmeldung (`id_ed25519` im Modell-Volume, ob Zugangsdatum im Sinne von [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): nicht gemessen, Inhalt nie lesen).
4. **Messhygiene:** Kaltstart von Lauf 2 mit frischem Volume; `timeout -k 5` beibehalten (Lauf 0 hing ohne `-k`); `rohdaten-04-allow-<n>.txt` trägt die Überschrift `Modus restricted` (Skript übergibt `restricted`) und die doppelte Schwärzung ist nicht idempotent: im Folge-Slice berichtigen, nicht in diesem.
5. **Suche:** `06-suche.sh` ohne Argument hat bekannte Altlast-Treffer (36 Zeilen aus vier älteren Messungen). Der Folge-Slice bekommt entweder eine benannte Ausnahmeliste im Skript oder nennt wie dieser Slice die engere Menge; V-1 klärt den Wortlaut.
6. **Zusage:** Messort Linux `amd64`, lokale Engine, Docker 29.9.2, Ollama `0.40.6`, Claude Code `2.1.310` (nicht gepinnt), ein Modell `gpt-oss:120b-cloud`; kein Satz im Lastenheft-Wortlaut aus diesem Protokoll.
