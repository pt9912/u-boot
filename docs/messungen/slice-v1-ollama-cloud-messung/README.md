# Messprotokoll: Cloud-Modelle von Ollama mit Claude Code im Sandbox-Devcontainer — Hosts, Fehlertext, Verhalten bis zur Anmeldung, ohne Konto

Rohdaten und Ablauf der Messung des Slice
[`slice-v1-ollama-cloud-messung`](../../plan/planning/in-progress/slice-v1-ollama-cloud-messung.md).
Skripte und Hilfsdateien liegen unter `tools/messung/ollama-cloud/` (Ablage: [`docs/messungen/README.md`](../README.md)).
Dieses Protokoll liefert **Messwerte und eine Liste des Ungemessenen**; es entscheidet nichts und enthält keinen Wortlaut für das Lastenheft
(Entscheidung: Folge-ADR der Welle; Anforderungen: [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)).

**Bedingung jeder Aussage unten: ohne Konto gemessen, Erfolgspfad nicht gemessen.** „Der Aufruf endet mit Meldung X“ heißt nicht „Cloud-Modelle
laufen“ oder „laufen nicht“; was ein Konto ändert, steht im Abschnitt *Nicht gemessen*.

**Messstand** (`rohdaten-01.txt`, `rohdaten-02.txt`, `rohdaten-03.txt`): Linux `x86_64` (`uname -sm`), Docker 29.9.2 (Client und Server, lokale Engine,
`overlayfs`, cgroup v2); `u-boot version 0.15.1-dev` (`make build`, Repo-HEAD `d322c1d0`, `internal/`, `cmd/` und die Vorlagen unverändert gegenüber diesem Stand);
Dev-Container-CLI 0.89.0 (Hilfs-Image, `Dockerfile.cli`); Basis-Image des Projekts `mcr.microsoft.com/devcontainers/base:debian`; Ollama `0.40.6` (Pin des erzeugten
Dockerfiles, `ARG OLLAMA_VERSION`); Claude Code `2.1.310` (Feature `claude-code:1`, nicht gepinnt, Stand beim Bau); Cloud-Modell **`gpt-oss:120b-cloud`**
(Name aus der Bibliothek von Ollama; ob das Modell in der Cloud verfügbar ist, ist nicht Gegenstand). Gemessen am 2026-10-24 (Zeiten in den Rohdaten: UTC),
Netzzugriff auf `ghcr.io`, `github.com`, `registry.ollama.ai`, `ollama.com`, `api.anthropic.com`, `registry.npmjs.org` und die Debian-Spiegel; **ohne jede Anmeldung und ohne
Zugangsdaten** (siehe *Zugangsdaten und Rohdaten*).

## Reihenfolge und Aufruf

`WORK=<messordner ausserhalb des Repos> tools/messung/ollama-cloud/00-alles.sh` führt 01 bis 03, dann `04-lauf.sh` zweimal (`LAUF=1`, `LAUF=2`) und zuletzt `05-aufraeumen.sh` aus;
`redact-probe.sh` läuft vor dem ersten Lauf, `06-suche.sh` nach dem Lauf. Jedes Skript startet sich mit bereinigter Umgebung neu (`env -i`) und schreibt seine Rohdaten
(geschwärzt) in dieses Verzeichnis (Vorgabe `OUT`). Voraussetzungen auf dem Host: Docker und GNU `make`; die **Skripte** rufen `sed`, `awk`, `curl` nur in Containern auf. Außerhalb der Skripte liefen während der Arbeit zwei Host-`sed`-Aufrufe (Verfahrensabweichung des Messenden); die Skripte, Rohdaten und deren Schwärzung gehen nicht durch sie.
Zu `rohdaten-01.txt` und `rohdaten-01-vorher.txt`: beide Dateien tragen den Zeitstempel 08:38 UTC, also einen Stand **nach** den Läufen (04: 08:18 bis 08:34) und nach dem Aufräumen (05: 08:34); Der Name „vorher“ stammt aus dem Skript (Schritt 1 der Reihenfolge); die Datei weist 0 Container und 0 Objekte mit dem Präfix aus, also den Zustand nach dem Aufräumen, nicht den Bestand vor der Messung. Schritt 1 wurde demnach nach den Läufen (erneut) ausgeführt; die Dateien sind Stand nach den Läufen. Die Umgebungsprüfung vor den Läufen belegen `rohdaten-04-*` (`umgebung-container-verbotene-namen=0`) und `umgebung_pruefen()` in `04-lauf.sh`.

| Schritt | Skript | Rohdaten | Messfrage |
|---|---|---|---|
| 0 | `redact-probe.sh` | (Konsole, siehe *Zugangsdaten und Rohdaten*) | Gegenprobe der Schwärzung vor dem ersten Lauf |
| 1 | `01-messstand.sh` | `rohdaten-01.txt`, `rohdaten-01-vorher.txt` | Messstand, Namen der Umgebungsvariablen samt Gegenprobe der Prüfung, Docker-Objekte vorher |
| 2 | `02-projekte.sh` | `rohdaten-02.txt`, `rohdaten-02-erzeugt-<variante>.txt`, `rohdaten-02-c7.txt` | `u-boot generate devcontainer` für drei Varianten (C7); Egress-Liste; Pin |
| 3 | `03-lauf-vorbereitung.sh` mit `Dockerfile.cli` | `rohdaten-03.txt` | Hilfs-Image der Dev-Container-CLI |
| 4 | `04-lauf.sh` mit `inner.sh`, `claude-stub.sh` | `rohdaten-04-<n>.txt`, `rohdaten-04-up-<variante>-<n>.txt`, `rohdaten-04-<modus>-<n>.txt` | C0 bis C6 in den Modi `discover`, `restricted`, `allow` |
| 5 | `05-aufraeumen.sh` | `rohdaten-05.txt`, `rohdaten-05-nachher.txt` | Aufräumen der Objekte mit dem Präfix `u-boot-messung-ocloud` |
| 6 | `06-suche.sh` | (Konsole) | Suche nach Zugangsdaten über Rohdaten und Skripte |

Varianten der erzeugten Projekte (`02-projekte.sh`; Features `ollama`, `node`, `claude-code` über `u-boot config set`, dann `u-boot generate devcontainer --yes`):
`std` Standardprofil (C7); `res` Sandbox-Profil mit Egress-Restriktion (Liste des Generators); `alw` wie `res` mit `devcontainer.sandbox.egress.allow=ollama.com`, gesetzt über
`u-boot config set` (C2b). Modi von `04-lauf.sh`: `discover` (kein Egress-Skript, protokollierender Resolver reicht alle Namen weiter), `restricted` (Projekt `res`),
`allow` (Projekt `alw`, im selben Skript wie `restricted`).

## Was Messvariante ist und was Auslieferung (Risiko 2 des Slice)

**Wie ausgeliefert:** Projekt, Dockerfile, `devcontainer.json`, Egress-Skript und Liste stammen aus `u-boot generate devcontainer` des gebauten Binary (`rohdaten-02-erzeugt-*.txt`); Image und
Container entstehen durch `devcontainer up` der Dev-Container-CLI 0.89.0 (Features `claude-code:1`, `node:1`, Ollama-Block, `postStartCommand` mit dem erzeugten Egress-Skript, Volume
auf `/home/vscode/.ollama`; Einhängepunkte in `rohdaten-04-<n>.txt`). Ollama `0.40.6` ist der Pin des Generators; **C0 ergab keine Abweichung**, die Messung läuft mit der ausgelieferten Version.

**Handvarianten (jede ausgewiesen, keine begründet eine Zusage):**

| Variante | Wo | Abweichung von der Auslieferung |
|---|---|---|
| `devcontainer up` der Dev-Container-CLI statt `u-boot devcontainer up` | `04-lauf.sh` | derselbe Mechanismus darunter; der Befehl von `u-boot` selbst wurde nicht gefahren |
| Modus `discover`: zweiter Container aus dem Image der CLI per `docker run --cap-add=NET_ADMIN`, ohne `postStartCommand`, ohne Workspace-Volume, mit benanntem Volume `u-boot-messung-ocloud-modelle` auf `/home/vscode/.ollama` | `04-lauf.sh` | `runArgs` wie im Projekt, aber kein Egress-Skript und kein `remoteEnv`; das Volume ersetzt das Volume der CLI |
| `docker exec -u vscode` in den Container statt Sitzung des Editors oder der CLI (`exec`) | `04-lauf.sh` | die `remoteEnv`-Einträge der erzeugten `devcontainer.json` (`${localEnv:ANTHROPIC_API_KEY}`, `${localEnv:GIT_TOKEN}`) wirken nur in Sitzungen der CLI; hier nicht angewandt, am Host waren sie ohnehin ungesetzt |
| Ein protokollierender `dnsmasq` ersetzt den Resolver (Verfahren O4 des Messprotokolls [`slice-v1-devcontainer-ollama-messung`](../slice-v1-devcontainer-ollama-messung/README.md)) | `inner.sh` | in `restricted`/`allow`: dieselbe Konfiguration wie das Egress-Skript plus `log-queries`; das Skript selbst lief unverändert als `postStartCommand` |
| `claude` durch einen Platzhalter ersetzt (`claude-stub.sh`), nur für den Aufruf `C3-launch-claude-stub`, über `PATH` | `inner.sh`, nur `restricted` und `allow` | zeigt Argumente und Variablen-**Namen**, die `ollama launch claude` übergibt; das echte `claude` startet dabei nicht |
| Aufrufe ohne Eingabe (`</dev/null`), mit `timeout -k 5`; ein Aufruf im Pseudo-Terminal über `script` | `inner.sh` | kein echtes Terminal mit Tastatur; ein interaktiver Ablauf ist **nicht** gemessen |
| Lauf 2 mit vorhandenem Volume aus Lauf 1 | `04-lauf.sh` | `up --remove-existing-container` behält die Volumes des Projekts; die Dateien in `~/.ollama` aus Lauf 1 waren in Lauf 2 schon da (Folge: Lauf 2 legt `id_ed25519` und die Blobs nicht neu an). Für Dateien gilt Lauf 1 |
| Modellname `gpt-oss:120b-cloud` | `04-lauf.sh` | ein Name von vielen; andere Cloud-Modelle sind nicht gemessen |

**Lauf 0:** ein erster Lauf von `04-lauf.sh` wurde abgebrochen, weil `timeout` ohne `-k` den Aufruf im Pseudo-Terminal (`ollama launch claude` unter `script`) nach 40 s nicht beendete (Prozess
nach über zehn Minuten noch vorhanden; Ursache nicht untersucht). Seine Rohdaten sind nicht erhalten (beim Neuaufsetzen der Rohdaten-Ablage versehentlich überschrieben); die Messwerte stammen aus Lauf 1 und 2.
Einzelheiten zu Lauf 0 (welche Aufrufe, wie oft gegen `ollama.com`/`registry.ollama.ai`) sind **nicht belegt**: keine Rohdaten erhalten, die Angabe der Vorgänge stünde auf Erinnerung und entfällt. Die Zahl der Cloud-Vorgänge belegen nur Lauf 1 und 2.

## Aussagen

Notation: *Lauf* zählt die gewerteten Läufe (1 und 2) je Modus; „2/2“ heißt: in beiden Läufen gleich. Die Dateinamen sind Muster (`<n>` = 1, 2); Zeilen sind per `AUSSAGE`-Zeilen und Marken
(`---- <Zeit> <Fall>`) auffindbar. `discover` = Namen frei, `restricted` = Liste des Generators, `allow` = Liste plus `ollama.com`.

### A. Erzeugtes Projekt (C7 und Grundlage der Läufe)

| Aussage (ohne Konto gemessen) | Rohdaten | Läufe |
|---|---|---|
| `generate devcontainer` mit `ollama`, `node`, `claude-code` im **Sandbox-Profil** erzeugt Dockerfile mit Ollama-Block (`ARG OLLAMA_VERSION=0.40.6`), `devcontainer.json` mit den OCI-Features `claude-code:1` und `node:1` und dem `remoteEnv` `GIT_TOKEN` und `ANTHROPIC_API_KEY` (nur als `${localEnv:…}`-Verweis), und das Egress-Skript | `rohdaten-02.txt`, `rohdaten-02-erzeugt-res.txt` | 1 je Variante |
| Egress-Liste `res` (11 Namen): `api.anthropic.com` (Herkunft `claude-code`), `registry.npmjs.org` (`claude-code`, `node`), `registry.ollama.ai` und der Blob-Host `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` (`ollama`) und die Basis-Hosts; **`ollama.com` steht nicht in der Liste** | `rohdaten-02.txt` (`u-boot egress list`) | 1 |
| Mit `egress.allow=ollama.com` über `config set` steht `ollama.com` mit Herkunft `user` in der Liste und als `server=`/`nftset=` im Egress-Skript | `rohdaten-02.txt`, `rohdaten-04-allow-<n>.txt` (Zeilen `conf:`) | 1 (Liste), 2/2 (im Container) |
| **C7 Standardprofil:** `generate devcontainer` mit `ollama` und `claude-code` ohne Sandbox-Profil aktualisiert den verwalteten Block der `.devcontainer/devcontainer.json` (344 Byte; Feature `claude-code:1`, kein Egress-Skript, keine `runArgs`); das Dockerfile (271 Byte, aus `init`) enthält keinen Ollama-Block (`ollama`-Treffer im Dockerfile: 0); Meldungen (zwei `Warning`-Zeilen, die erste mit der Kennung [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), die zweite mit [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)): `devcontainer feature claude-code comes from the third-party publisher anthropics …` und `devcontainer feature ollama has no effect without devcontainer.profile=sandbox; no Ollama installation was generated`; Exit 0. Ein zweiter Aufruf meldet `devcontainer already up to date; no changes.` und wiederholt beide Warnungen | `rohdaten-02.txt` (Abschnitt `Projekt std`), `rohdaten-02-c7.txt`, `rohdaten-02-erzeugt-std.txt` | 2 (erster und zweiter Aufruf) |

### B. C0: Unterbefehl und Version

| Aussage | Rohdaten | Läufe |
|---|---|---|
| Die ausgelieferte Version `ollama version is 0.40.6` kennt `launch`; `ollama launch --help` nennt `claude` (Claude Code) als Ziel und die Flags `--config`, `--model`, `--restore`, `-y/--yes`; **keine Pin-Hebung nötig** | `rohdaten-04-{discover,restricted,allow}-<n>.txt`, Abschnitt `C0` (`AUSSAGE c0-launch-in-help=1`, `c0-launch-claude-in-help=4`) | 6 (3 Modi × 2) |
| `claude` liegt im `PATH` des Containers (`/usr/local/share/nvm/current/bin/claude`), `claude --version` meldet `2.1.310 (Claude Code)` | dort, Abschnitt `C0` | 6 |
| `ollama --version` ohne laufenden Server meldet `Warning: could not connect to a running Ollama instance` und `Warning: client version is 0.40.6` auf stderr | dort | 6 |

### C. C1/C2: Aufrufe mit einem Cloud-Modell, Hosts und Fehlertext

Hosts = Namen im Resolver-Log (Anfragen A/AAAA), ohne Antwort-Inhalte.

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `ollama pull gpt-oss:120b-cloud` endet in **allen drei Modi mit Exit 0** (`success`, in `ollama list` mit Größe `-`); der Resolver sieht dabei **nur `registry.ollama.ai`** (2 bis 4 Anfragen), weder `ollama.com` noch den Blob-Host. Gezogen werden zwei kleine Blobs (264 und 307 Byte) und ein Manifest (426 Byte) | `…-<modus>-<n>.txt`, Abschnitt `C1-pull` | 6 |
| Beim Start von `ollama serve` fragt Ollama **`ollama.com`** an (2 Anfragen `discover`/`allow`, 4 Anfragen `restricted`); mit `OLLAMA_NO_CLOUD=1` keine | `…`, Abschnitt `Namen beim Start des Servers` | 6 |
| **`ollama run` im Modus `discover` und `allow`:** Meldung `You need to be signed in to Ollama to run Cloud models.`, **Exit 0**, Dauer 0,2 s; **keine** neue Namensanfrage im Aufruf | `rohdaten-04-discover-<n>.txt`, `rohdaten-04-allow-<n>.txt`, Abschnitt `C1-run` | 4 (2/2 je Modus) |
| **`ollama run` im Modus `restricted` (C2):** `Error: Post "https://ollama.com:443/api/show?ts=…": dial tcp: lookup ollama.com on 127.0.0.1:53: server misbehaving`, **Exit 1**, Dauer < 0,1 s; der Resolver zeichnet 4 Anfragen nach `ollama.com` auf, das Egress-Skript beantwortet sie mit `config error` | `rohdaten-04-restricted-<n>.txt`, Abschnitt `C1-run` | 2/2 |
| **C2b:** mit dem Eintrag `ollama.com` (Modus `allow`) ändert sich die Meldung von `run` auf den Text der Zeile davor (Anmeldung nötig, Exit 0). Das Hindernis, auf das der Aufruf als nächstes stößt, ist die Anmeldung; was **nach** einer Anmeldung folgt, ist nicht gemessen | `rohdaten-04-allow-<n>.txt` gegen `rohdaten-04-restricted-<n>.txt` | 2/2 |
| **`ollama launch claude --model gpt-oss:120b-cloud`, `discover` und `allow`:** `Error: gpt-oss:120b-cloud requires sign in`, Exit 1, 0,2 s; `claude` wird nicht gestartet (der Platzhalter-Aufruf in `allow` endet ebenso: Exit 1, kein `STUB`) | `…-discover-<n>.txt`, `…-allow-<n>.txt`, Abschnitte `C3-launch-claude`, `C3-launch-claude-stub` | 2/2 |
| **`ollama launch claude …`, `restricted`:** keine Anmeldeaufforderung; `ollama` startet das installierte `claude` (Namensanfrage `api.anthropic.com`, 2 Anfragen, erlaubt), dieses meldet `Error: Input must be provided either through stdin or as a prompt argument when using --print`, `ollama` meldet `Error: exit status 1`, Exit 1, 0,7 bis 0,8 s; daneben 12 abgelehnte Anfragen nach `ollama.com` | `rohdaten-04-restricted-<n>.txt`, Abschnitt `C3-launch-claude` | 2/2 |

### D. C3: `ollama launch claude` bis zum Hindernis

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `ollama launch claude` findet das installierte `claude` (`PATH`), startet es im Modus `restricted` (kein Zugriff auf `ollama.com`), in `discover` und `allow` nicht (Anmeldehinweis davor; warum der Ablauf davon abhängt, ist nicht gemessen), und lädt oder installiert selbst nichts (keine Namensanfrage außer `ollama.com`/`api.anthropic.com`; angelegte Dateien: siehe F) | siehe C | 2/2 |
| Mit dem Platzhalter (`restricted`) übergibt `ollama launch claude` das Argument `--model gpt-oss:120b-cloud` und setzt in der Umgebung von `claude` die **Namen** `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `ANTHROPIC_BASE_URL`, `ANTHROPIC_DEFAULT_HAIKU_MODEL`, `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL` sieben `CLAUDE_CODE_*`-Variablen und `ENABLE_CLAUDEAI_MCP_SERVERS` (ausgewiesen sind nur Namen mit `ANTHROPIC`, `CLAUDE`, `OLLAMA`, `API`, `TOKEN`, `KEY`); `ANTHROPIC_BASE_URL=http://127.0.0.1:11434` (Wert nur, weil lokal). **Werte der übrigen Variablen sind nicht ausgewiesen und nicht gemessen** | `rohdaten-04-restricted-<n>.txt`, Zeilen `STUB` (Handvariante) | 2/2 |
| Aufruf im Pseudo-Terminal ohne Eingabe: **keine Ausgabe** (0 Byte, bzw. 8 Byte `Killed` nach `SIGKILL`) und kein Ende innerhalb von 25 s; beendet nach `timeout` (Exit 124) oder nach `-k 5` (Exit 137, 30 s) — die Läufe unterscheiden sich in dieser Zahl. Was ein Nutzer am Terminal sieht, ist **nicht gemessen** | `…-<modus>-<n>.txt`, Abschnitt `C3 mit Pseudo-Terminal` | 6, Exit 124/137 wechselnd |

### E. C4/C5: `ollama signin`, Variablen, `OLLAMA_NO_CLOUD`

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `ollama signin` im Modus `discover` und `allow`: Exit 0; Text `You need to be signed in to Ollama to run Cloud models.` / `If your browser did not open, navigate to:` und eine Anmelde-URL der Form `https://ollama.com/connect` mit Query (Query geschwärzt); `ollama signout` danach: `You are not signed in to ollama.com`, Exit 0; in diesem Aufruf keine neue Namensanfrage; **die Anmeldung wurde nicht ausgeführt** | `…-discover-<n>.txt`, `…-allow-<n>.txt`, Abschnitte `C4-signin`, `C4-signout` | 4 |
| `ollama signin` im Modus `restricted`: `Error: account unavailable`, Exit 1 (4 abgelehnte Anfragen nach `ollama.com`); `ollama signout`: `Error: there was an error signing out`, Exit 1 | `…-restricted-<n>.txt` | 2/2 |
| Das Binary enthält die Zeichenketten `OLLAMA_API_KEY`, `OLLAMA_CLOUD_BASE_URL`, `OLLAMA_NO_CLOUD` und `OLLAMA_LAUNCH_DSH_API_KEY` (gelesen mit `grep -a` aus `/usr/local/bin/ollama`; Go-Zeichenketten liegen aneinandergereiht, daher auch verklebte Treffer). **`ollama --help`, `ollama launch --help` und `ollama help serve` nennen keine Variable für den Zugang**; `ollama help serve` nennt `OLLAMA_NO_CLOUD` (*Disable Ollama cloud features (remote inference and web search)*). Ob `OLLAMA_API_KEY` ein Zugang ist und wie er wirkt, ist **nicht gemessen** | `…-<modus>-<n>.txt`, `AUSSAGE c4-ollama-var-…`, Abschnitt `ollama help serve` | 6 |
| Der Generator setzt für `ollama` keine Zugangsvariable (`remoteEnv` nennt nur `GIT_TOKEN`, `ANTHROPIC_API_KEY`) | `rohdaten-02-erzeugt-res.txt` | 1 |
| **`OLLAMA_NO_CLOUD=1`** (Server **und** Client): `pull` unverändert (Exit 0, nur `registry.ollama.ai`); `run`: `Error: ollama cloud is disabled: remote model details are unavailable`, Exit 1, 0,0 s, **keine** Namensanfrage; `launch claude --model …`: `Warning: ignoring --model gpt-oss:120b-cloud because cloud is disabled` und `Error: model selection requires an interactive terminal; use --model to run in headless mode`, Exit 1; der Server startet ohne Anfrage an `ollama.com` und protokolliert `Ollama cloud disabled: true` | `…-<modus>-<n>.txt`, Abschnitte `C5-*` | 6 |
| `ollama signin` unter `OLLAMA_NO_CLOUD=1` verhält sich wie ohne die Variable (`discover`/`allow`: Anmeldehinweis, Exit 0, 2 Anfragen nach `ollama.com` im Aufruf; `restricted`: `account unavailable`, Exit 1, 4 abgelehnte Anfragen) | `…`, Abschnitt `C5-signin` | 6 |

### F. C6: Dateien, die die Aufrufe anlegen (Pfad und Größe; kein Inhalt gelesen)

Der Einhängepunkt `/home/vscode/.ollama` ist ein **benanntes Docker-Volume** des Projekts (`Einhängepunkte` in `rohdaten-04-<n>.txt`: `volume /workspaces/<projekt>` und `volume /home/vscode/.ollama`);
die Spalte *Volume* sagt, ob der Pfad darauf liegt (`df` auf den Pfad). Gilt für **Lauf 1** (frisches Volume).

| Pfad | Größe | Angelegt von | Volume |
|---|---|---|---|
| `~/.ollama/id_ed25519`, `~/.ollama/id_ed25519.pub` | 387 und 81 Byte | erster Start von `ollama serve` (in allen Modi, **vor** jedem Anmeldeversuch) | ja |
| `~/.ollama/cache/model-recommendations.json` | 2370 Byte | `ollama serve`/`launch`, wenn `ollama.com` erreichbar ist (`discover`, `allow`); in `restricted` nicht angelegt | ja |
| `~/.ollama/models/blobs/sha256-…` (2 Dateien), `~/.ollama/models/manifests/registry.ollama.ai/library/gpt-oss/120b-cloud` | 264, 307, 426 Byte | `ollama pull` | ja |
| `~/.ollama/config.json` | 104 Byte | `ollama launch claude` (nur `restricted`, wo `claude` startete) | ja |
| `~/.claude.json`, `~/.claude/backups/.claude.json.backup.<zeit>` | 852 und 84 Byte | `claude`, gestartet von `ollama launch claude` (`restricted`) | **nein** (Wurzel-Dateisystem des Containers) |
| `ollama signin` (`discover`, `allow`) | keine Datei neu oder geändert | — | — |

Zusammenhang mit [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten nicht im Workspace-Volume): kein Pfad aus den Aufrufen liegt im
Workspace-Volume (`/workspaces/<projekt>`); der Schlüssel `id_ed25519` liegt im **Modell-Volume**. Ob er ein Zugangsdatum im Sinne der Anforderung ist, ist nicht gemessen und nicht Gegenstand
dieses Protokolls (Inhalt nicht gelesen).

### G. Bestätigung der früheren Messung

Das Handbuch nennt als gemessen (Ollama `0.40.6`, `amd64`, hinter der Egress-Beschränkung): `ollama pull …-cloud` gelingt, `ollama run` scheitert an `ollama.com` (nicht auflösbar), `ollama signin` ohne Konto
meldet `account unavailable`. **Bestätigt** in `restricted` (2/2) für Modell `gpt-oss:120b-cloud`: `pull` Exit 0, `run` mit Lookup-Fehler `ollama.com`, `signin` `Error: account unavailable`. **Neu gegenüber der
früheren Messung:** das Verhalten in `discover`/`allow` (Anmeldehinweis statt Lookup-Fehler), `launch claude`, `OLLAMA_NO_CLOUD`, die angelegten Dateien und die Variablen-Namen.

## Nicht gemessen — und der Messplan für den Lauf mit Konto

Alles in dieser Tabelle ist **ohne Konto nicht messbar** und wurde **nicht gemessen und nicht geraten**. Die Spalte *Aufruf* ist das, was der bedingte Folge-Slice
(der bedingte Mess-Slice der Welle mit Konto, Kennung noch nicht vergeben) fährt; *Bedingung*: Entscheidung und Zugang des Projektinhabers (Welle-Datei, Frage 2); ein Zugang gehört **nur** in die
Umgebung des Hosts beim Lauf (nie in Datei, Skript, Protokoll oder Chat); die Hilfsmittel dieses Slice (`lib.sh` mit `redact`, `06-suche.sh`) sind der Ausgangspunkt, die Schwärzung
braucht vorher eine Gegenprobe für die Form des Zugangs (die Anmelde-URL und ihr Code sind eine eigene Form).

| Nicht gemessen | Aufruf im Folge-Slice | Bedingung / Hinweis |
|---|---|---|
| Anmeldung mit Konto | `ollama signin` im Container (`discover`), Freigabe in einem Browser außerhalb des Containers; danach `ollama signout` | Konto des Projektinhabers; Hosts des Anmeldeablaufs aus dem Resolver-Log; Dateien als Pfad und Größe (**kein Inhalt** von `id_ed25519`, `config.json`) |
| Lauf eines Cloud-Modells | `ollama run gpt-oss:120b-cloud "<kurzer Prompt>"` einmal, `discover`, dann `restricted` ohne und mit `ollama.com` | Konto; Antwort-Inhalt nicht ins Protokoll; Zahl der Läufe ausweisen (Kontingent) |
| Hosts nach der Anmeldung | wie zuvor, Namen im Resolver-Log, Abgleich mit der erzeugten Liste (Zusage-Menge der Hosts: Linux `amd64`, lokale Engine) | Konto |
| Zusammenspiel mit Claude Code | `ollama launch claude --model <cloud-modell> --yes` in einem Pseudo-Terminal mit Eingabe, sowie eine headless Fassung (`claude -p`-Aufruf über `launch`, falls möglich); Hosts inkl. `api.anthropic.com`; Verhältnis der von `launch` gesetzten Variablen zu `ANTHROPIC_API_KEY` aus `remoteEnv` | Konto; **ohne** den Anthropic-Zugang des Projektinhabers (nur der Zugang von Ollama) |
| Name und Wirkung der Zugangsvariable | `OLLAMA_API_KEY` als Kandidat aus dem Binary: Lauf mit dem Wert aus der Host-Umgebung ohne Anmeldung per `signin`; Namensprüfung in Container und Rohdaten | Konto/Schlüssel von Ollama; Wert nur als Umgebungsvariable des Hosts, durchgereicht ohne Wert auf der Kommandozeile |
| Interaktiver Ablauf von `ollama launch claude` | derselbe Aufruf im Terminal mit Tastatur (hier: ohne Ausgabe und ohne Ende in 25 s) | Mensch am Terminal oder ein Terminal-Treiber mit Eingabe; ohne Konto fragt der Ablauf möglicherweise die Anmeldung ab (nicht gemessen) |
| Wirkung von `OLLAMA_NO_CLOUD` mit Konto | `run`, `launch` unter der Variablen | Konto |
| Andere Cloud-Modelle als `gpt-oss:120b-cloud`, andere Ollama-Versionen, andere Claude-Code-Versionen | je Modell/Version ein Lauf | Verfügbarkeit des Modells im Konto; Pin-Hebung eigener Slice |
| `arm64`, Docker Desktop, Podman als Host-Engine | — | Zusage-Menge ist Linux `amd64` mit lokaler Engine; ein anderer Host begründet keine Zusage |
| Kaltstart von Lauf 2 | Lauf 2 mit frischem Volume | Lauf 2 lief mit vorhandenem Volume (siehe Handvarianten) |

## Zugangsdaten und Rohdaten

- **Kein Zugang verwendet.** Kein Skript liest, setzt oder reicht einen Token, eine Anmeldung oder eine Konfigurationsdatei eines Agenten durch; es gibt keinen Anmeldeabschluss (`ollama signin` bricht ohne Browser ab,
  der Container hat keinen). Keine Variable mit Zugang wurde angelegt.
- **Umgebung des Messlaufs, geprüft nach Namen:** jedes Skript startet sich mit `env -i` (nur `PATH`, `HOME`, `USER`, `WORK`, `OUT`, `ROOT`, gegebenenfalls `LAUF`) neu; `01-messstand.sh` weist die Namen aus
  (`HOME OUT PATH PWD ROOT SHLVL UB_BEREINIGT USER WORK _`), `umgebung_pruefen()` meldet jeden Namen mit `ANTHROPIC`, `CLAUDE`, `OLLAMA`, `GH_TOKEN`, `GITLAB_TOKEN`, `GIT_TOKEN`, `TOKEN`, `SECRET`, `AUTH`, `API_KEY` u. a. und bricht
  `04-lauf.sh` ab. Gegenprobe der Prüfung: gesetzte Namen (Platzhalterwert `platzhalter`) schlagen an (`rohdaten-01.txt`, Exit 1). Der Messcontainer trägt nur `HOME HOSTNAME LANG NVM_DIR NVM_SYMLINK_CURRENT PATH PWD SHLVL _`
  (`AUSSAGE umgebung-container-verbotene-namen=0`). Der Host des Messenden trug in der Sitzung Variablen wie `CLAUDE_CODE_*` und `GEMINI_CLI_IDE_AUTH_TOKEN`; sie erreichten wegen `env -i` weder die Skripte noch einen Container.
- **Schwärzung:** `redact()` (formgebunden: Header, Zuweisungen an Namen mit Zugangswort, JSON-Felder, URLs mit Zugang, Token-Präfixe, Anmelde-URLs mit Query, Code/Schlüssel nach `code`, `key`, `nonce`) mit der Gegenprobe
  `redact-probe.sh` **vor dem ersten Lauf**: 24 Formen geschwärzt, 5 neutrale Zeilen unverändert (Exit 0); rot gesehen: ohne die Regel für Anmelde-URLs lässt die Probe drei Zeilen durch (Exit 1). Die Zusage gilt den Formen der
  Probe, nicht jedem Token.
- **Suche vor dem Commit:** `06-suche.sh` über dieses Verzeichnis und `tools/messung/ollama-cloud/` (Tokenpräfixe, Schlüsselblöcke, Header, URL-Zugang, Zuweisungen, Anmelde-URLs): 0 Treffer, Exit 0 (Aufruf `06-suche.sh <dieses Verzeichnis> tools/messung/ollama-cloud`, gefahren nach dem letzten Schreiben dieser Datei; ein früherer Lauf der Fix-Runde endete mit Exit 1 wegen zweier Beispielzeilen in dieser Datei, die umformuliert wurden, kein Zugangsdatum); Gegenprobe der Suche: eine Datei mit
  einer Zuweisung an `GH_TOKEN` (Wert mit dem Präfix `ghp_`) schlägt an (Exit 1). Treffer über **ganz** `docs/messungen/` und `tools/messung/` (die Vorgabe des Skripts): ältere Protokolle enthalten die öffentlichen Geräteschlüssel, die `ollama serve` in Wegwerf-Containern ausgibt (`Your new public key is:` in
  [`slice-v1-devcontainer-ollama-messung`](../slice-v1-devcontainer-ollama-messung/README.md); öffentlich, der private Teil wurde nie ausgegeben), GPG-Fingerabdrücke
  ([`slice-v1-devcontainer-ollama-amd64-messung`](../slice-v1-devcontainer-ollama-amd64-messung/README.md), [`slice-v1-forge-cli-quellen-messung`](../slice-v1-forge-cli-quellen-messung/README.md)), Testnamen mit `Ollama_` und Platzhalter-Zuweisungen (`DUMMY` in `forge-cli`, `platzhalter-…` in
  `devcontainer-start-befehl`); alle von Hand als **kein Zugangsdatum** gelesen, nicht verändert.
- **Von Hand gelesen:** alle `rohdaten-04-*`-Abschnitte außerhalb der Hilfetexte, `rohdaten-01`, `rohdaten-03`, `rohdaten-05`; die erzeugten Dateien (`rohdaten-02-erzeugt-*`) tragen nur `${localEnv:…}`-Verweise, keinen Wert. Die Anmelde-URLs stehen als
  `https://ollama.com/connect?[Query entfernt: Anmelde-URL]`.
- Das Binary-Fundstück `OLLAMA_API_KEY` ist ein Name, kein Wert.
- **Rohdaten-Eigenheiten (nur ausgewiesen, nicht geändert):** In den Sammeldateien `rohdaten-04-<n>.txt` ist die Schwärzung doppelt angewandt und nicht idempotent (zerstückelte Marke `[Query entfernt: Anmelde-URL] entfernt: Anmelde-URL]`, kein Wert sichtbar); `rohdaten-04-allow-<n>.txt` trägt die Überschrift `Modus restricted`, weil `inner.sh` für `allow` den Modus `restricted` übergibt; `rohdaten-02-c7.txt` läuft nicht durch `redact()` (Ausgabe von `generate`, kein Zugang, von Hand gelesen).

## Aufräumen

`05-aufraeumen.sh` (`rohdaten-05.txt`, `rohdaten-05-nachher.txt`): entfernt Container, Volumes und Images, deren Name `u-boot-messung-ocloud` oder `vsc-u-boot-messung-ocloud` beginnt (nach dem Namen, nie nach Zeitfenster);
`objekte-mit-praefix-nachher=0`. Nicht entfernt: der Build-Cache und die Basis-Images (Tags ohne Präfix); fremde Objekte des Hosts wurden nicht angefasst (Anzahl Container vorher/nachher unterscheiden sich um fremde
Arbeit auf dem Host, `rohdaten-05.txt`). Der Messordner `WORK` (erzeugte Projekte) liegt außerhalb des Repos.

## Lesehilfen

`AUSSAGE <schlüssel>=<wert>` sind die stabilen Zeilen, die ein zweiter Lauf wiederholen muss: Vergleich `grep '^AUSSAGE' rohdaten-04-<modus>-1.txt` gegen `…-2.txt` — in den drei Dateipaaren (`discover`, `restricted`, `allow`, je Lauf 1 und 2) unterscheidet sich **nur** der
Exit-Code des Pseudo-Terminal-Aufrufs (`C3-launch-claude-tty-exit`: 124 oder 137). Zahlen der Resolver-Anfragen schwanken in Details (`registry.ollama.ai` 4 gegen 6 je nach Lauf); die Anzahl gesperrter Anfragen an `ollama.com` ist ein
Wiederholungszähler des Clients, kein Maß für Verbindungen. Zeiten der Aufrufe: `dauer=` (Sekunden); `devcontainer up` des ersten Laufs dauerte 539 s (kalter Bau der Features), der zweite Lauf 5,5 s.
