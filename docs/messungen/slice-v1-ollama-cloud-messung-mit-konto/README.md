# Messprotokoll: Cloud-Modelle von Ollama mit Claude Code im Devcontainer — mit Konto, über einen angemeldeten Ollama-Server des Hosts (Weg a) und im Container (Weg b)

Rohdaten und Ablauf der Messung des Slice
[`slice-v1-ollama-cloud-messung-mit-konto`](../../plan/planning/in-progress/slice-v1-ollama-cloud-messung-mit-konto.md).
Skripte und Hilfsdateien liegen unter `tools/messung/ollama-cloud-konto/` (Ablage: [`docs/messungen/README.md`](../README.md)).
Dieses Protokoll liefert **Messwerte und eine Liste des Ungemessenen**; es entscheidet nichts und enthält keinen Wortlaut für das Lastenheft
(Entscheidung: Folge-ADR der Welle; Anforderungen: [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)).
Das Protokoll des Vorgängers ohne Konto: [`slice-v1-ollama-cloud-messung`](../slice-v1-ollama-cloud-messung/README.md).

**Stand dieses Protokolls: Weg (a) und Weg (b) gemessen; der Projektinhaber hat sich im Container selbst angemeldet. Offen: das Entfernen des Geräteschlüssels in den Kontoeinstellungen (Schritt des Projektinhabers).** Jede Aussage unten gilt
**bedingungsgleich**: *mit dem Konto des Projektinhabers, auf einem Linux-`amd64`-Host mit Docker 29.9.2 (lokale Engine), ein Host-Server, zwei Modelle, eine Claude-Code-Version, am 2026-10-24
gemessen, headless (`claude -p`, kein interaktiver Ablauf)*. „Der Aufruf lieferte eine Antwort“ heißt nicht „Cloud-Modelle laufen“ und begründet keine Zusage.

**Antwort-Inhalte sind nicht aufgenommen.** Die Skripte weisen bei Aufrufen mit Prompt nur Exit-Code, Dauer, Zeilen- und Byte-Zahl der Ausgabe und die Zeilen aus, die nach einer Meldung aussehen
(Muster `MELDUNG` in `inner.sh`); „Antwort“ in den Tabellen hat **zwei Kriterien**, je nach Aufruf. `claude -p`: Exit 0 und eine Zeile ohne Meldungsmuster von höchstens zwei Zeichen (Prompt: ein Satz mit der Bitte um ein Wort).
`ollama run` und `ollama launch claude`: nur Exit 0 und mindestens eine Ausgabezeile ohne Meldungsmuster; **ob das die erbetene Ein-Wort-Antwort ist, belegen die Rohdaten nicht** (die Ausgabe von `ollama run` hat 4 bis 16 Zeilen
bis 75 Zeichen, die von `launch` zwei Zeilen bis 170 Zeichen; Denkausgabe, Banner oder Hinweise sind nicht ausgeschlossen). „Antwort“ heißt dort „Aufruf endete mit Exit 0 und ohne Fehlermeldung, mit Ausgabe“. Der Inhalt selbst
steht in keiner Rohdatei und in keinem Chat.

## Messstand

Linux `x86_64` (Ubuntu 24.04.5 LTS), Docker 29.9.2 (Client und Server, lokale Engine, `overlayfs`, cgroup v2) (`rohdaten-01.txt`); `u-boot version 0.15.1-dev` (`make build`, Repo-HEAD `f7b63acf` beim Bau,
`internal/`, `cmd/` und die Vorlagen unverändert; `rohdaten-02.txt`); Dev-Container-CLI 0.89.0 (`rohdaten-03.txt`); Basis-Image `mcr.microsoft.com/devcontainers/base:debian`;
Ollama-Client im Container `0.40.6` (Pin des erzeugten Dockerfiles); **Ollama-Server des Hosts `0.40.10`** (`api/version` über den Relay, `rohdaten-05-discover-1-111102.txt`, Fall `K1b-ollama-version`:
`ollama version is 0.40.10` und `Warning: client version is 0.40.6`); Claude Code `2.1.310` und Node `v24.22.1` in beiden Images (`rohdaten-05-discover-1-113322.txt`, `rohdaten-05-std-1-113340.txt`; nicht gepinnt, Stand beim Bau);
Modelle `glm-5.3-flash:cloud[1m]` (Modell 1) und `gpt-oss:120b-cloud` (Modell 2). Zeiten in den Rohdaten: UTC (die Dateizeit des Dateisystems: UTC+2). Der Host-Server lauscht nach `ss` **nur auf `127.0.0.1:11434`**
(`rohdaten-04-probe-113412.txt`); er wurde nie vom Host aus aufgerufen, nur aus Containern.

## Reihenfolge und Aufruf

`WORK=<messordner ausserhalb des Repos>`; jedes Skript bereinigt seine Umgebung im selben Prozess (Positivliste, `unset` aller anderen Namen), prüft die Namen (nie die Werte) und bricht ab, wenn das Artefakt von
`00-vorlauf.sh` fehlt oder zu einem anderen Stand von `lib.sh`/`redact-probe.sh` gehört. Rohdaten schreibt jedes Skript in **Stufe 1** (`$WORK/stufe1`, außerhalb des Repos); nach der Handlesung kopiert
`10-uebernehmen.sh --gelesen` (**Stufe 2**) nach hier und lässt `09-suche.sh` laufen.

| Schritt | Skript | Rohdaten | Messfrage |
|---|---|---|---|
| 0 | `00-vorlauf.sh` (ruft `redact-probe.sh`) | `rohdaten-00-vorlauf.txt` | Gegenprobe der Schwärzung, Umgebungsprüfung samt Gegenprobe, Hash des Stands; der Stand **vor dem ersten Prompt** (Hash `c9a5467b…`) liegt nur in der Historie von Commit `6bfc17e8` (2026-10-24 11:06:29 UTC; erster Prompt 11:07:26 UTC), `git show 6bfc17e8:docs/messungen/slice-v1-ollama-cloud-messung-mit-konto/rohdaten-00-vorlauf.txt`. **Die Datei im Baum ist der zweite Vorlauf** (11:27:20 UTC, nach 22 Prompts), geschrieben, weil die Hash-Prüfung nach der Änderung von `lib.sh` (Aufräumen) einen neuen Stand verlangte |
| 1 | `01-messstand.sh` | `rohdaten-01.txt`, `rohdaten-01-vorher.txt` | Messstand, Namen der Umgebung, Docker-Objekte vorher (vor allen Bauten) |
| 2 | `02-projekte.sh` | `rohdaten-02.txt`, `rohdaten-02-erzeugt-<variante>.txt` | `u-boot generate devcontainer` für `std`, `res`, `alw`; Handvariante `--add-host`; Egress-Listen |
| 3 | `03-lauf-vorbereitung.sh` (`Dockerfile.cli`) | `rohdaten-03.txt` | Hilfs-Image der Dev-Container-CLI |
| 3b | `05-lauf-host.sh bauen` | `rohdaten-05-bauen-1.txt` | Images der drei Projekte (kein Prompt) |
| 4 | `04-relay.sh probe`, `start`, `stop` | `rohdaten-04-probe-*.txt`, `rohdaten-04-relay-{start,stop}-*.txt` | K1 ohne Relay; Relay für die Dauer eines Falls |
| 5 | `05-lauf-host.sh <teil> <fall>…` mit `inner.sh` | `rohdaten-05-<teil>-<lauf>-<HHMMSS>.txt` | Weg (a), K1 bis K3; Teile `std`, `discover`, `restricted`, `allow`, `c`; Läufe 1 und 2 |
| 6 | `06-lauf-container.sh` | `rohdaten-06-*` | Weg (b): K4 bis K8 (Stufen `vorbereiten`, `messen discover`, `messen restricted`, `ende`) |
| 7 | `07-anmeldung.sh` | **keine** (nur Terminal des Projektinhabers) | Anmeldung im Container (vom Projektinhaber ausgeführt) |
| 8 | `08-aufraeumen.sh weg-a`, `08-aufraeumen.sh alles` | `rohdaten-08-weg-a{,-nachher}.txt`, `rohdaten-08-alles{,-nachher}.txt` | Aufräumen von Container, Volumes, Netzen; nach Weg (a) bleiben die Images für Weg (b), danach entfernt `alles` auch sie |
| 9 | `09-suche.sh`, `10-uebernehmen.sh` | (Konsole) | Suche mit Ausnahmeliste; Stufe 2 |
| — | Gegenproben | `rohdaten-09-gegenproben.txt` | rot gesehene Gegenstücke aller Wächter (Abschnitt *Gegenproben*) |

Varianten der erzeugten Projekte (`02-projekte.sh`; Ergebnis von `u-boot init`, `config set`, `generate devcontainer`): `std` Standardprofil mit `node` und `claude-code`, **kein Ollama im Container**;
`res` Sandbox-Profil mit Egress-Restriktion (`ollama`, `node`, `claude-code`; Liste des Generators, elf Namen); `alw` wie `res` mit `devcontainer.sandbox.egress.allow=host.docker.internal`
(über `config set`; `u-boot egress list` führt `host.docker.internal` mit Herkunft `user`; `rohdaten-02.txt`). `ollama.com` steht in keiner der Listen.

## Was Messvariante ist und was Auslieferung

**Wie ausgeliefert:** Projekt, Dockerfile, `devcontainer.json`, Egress-Skript und Liste stammen aus `u-boot generate devcontainer` des gebauten Binary; Image und Container entstehen durch `devcontainer up` der
Dev-Container-CLI 0.89.0 (Features `claude-code:1`, `node:1`, Ollama-Block, `postStartCommand` mit dem erzeugten Egress-Skript, Volume auf `/home/vscode/.ollama`).

**Handvarianten (jede ausgewiesen; keine begründet eine Zusage; die Auslieferung erzeugt keine davon):**

| Variante | Wo | Abweichung von der Auslieferung |
|---|---|---|
| `--add-host=host.docker.internal:host-gateway` in `runArgs` jedes Projekts | `02-projekte.sh` (sed im Werkzeug-Container nach `generate`) | `u-boot` erzeugt keinen solchen Eintrag; ohne ihn löst `host.docker.internal` im Container nicht auf (`rohdaten-04-probe-113412.txt`, Exit 2) |
| **Relay**: Hilfs-Container `--network host` mit `alpine/socat@sha256:603174b77eebaba25ba8f8ebb15e1234d6600409554e95d6fb203fd106150287` (Digest beim ersten Pull gelesen), bindet `172.17.0.1:11434` (die Adresse, die `host-gateway` im Container ergibt) und reicht an `127.0.0.1:11434` weiter; nie `0.0.0.0` (`ss`: `rohdaten-04-relay-start-110648.txt`); nur während eines Falls; vorher Prüfung auf fremde Container im Netz `bridge` (0), nachher `docker ps -a`: 0 Relay-Container (in jeder Fall-Datei, Zeile `relay: entfernt`) | `04-relay.sh` | nicht Produkt; solange er läuft, kann jeder Container im Netz `bridge` den angemeldeten Server des Hosts nutzen |
| Modus `discover`: zweiter Container aus dem Image von `res` per `docker run --cap-add=NET_ADMIN --add-host=…`, ohne `postStartCommand`, ohne Workspace- und Modell-Volume | `05-lauf-host.sh` | wie im Vorgänger-Protokoll; ein Weg (a) braucht kein Modell-Volume |
| Protokollierender `dnsmasq` ersetzt den Resolver (`discover`: Weiterreichen aller Namen; `restricted`: dieselbe Konfiguration wie das Egress-Skript plus `log-queries`) | `inner.sh setup` | das Egress-Skript lief unverändert als `postStartCommand`; Verfahren O4 des Vorgängers |
| **K3c**: Zusatzzeile `address=/host.docker.internal/172.17.0.1` im Resolver, und getrennt davon `nft add element ip u_boot_egress allowed { 172.17.0.1 }` von Hand | `05-lauf-host.sh c`, `inner.sh nftadd` | Eingriffe in die Egress-Restriktion, die `u-boot` nicht erzeugt; zeigen nur, ob der Weg bei erlaubter Adresse trägt |
| `docker exec -u vscode -w /home/vscode` statt Sitzung der CLI oder des Editors; `claude -p` und `ollama run` headless, ohne Terminal | `05-lauf-host.sh` | die `remoteEnv`-Einträge der erzeugten `devcontainer.json` wirken nur in Sitzungen der CLI; hier nicht angewandt |
| `ANTHROPIC_BASE_URL=http://host.docker.internal:11434` und der Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` im Aufruf IM Container (nie in der Umgebung des Hosts); `OLLAMA_HOST=http://host.docker.internal:11434` für `ollama run` und `ollama launch` | `inner.sh fall` | Nutzung nach Angabe des Projektinhabers; der Platzhalter ist kein Geheimnis |
| `devcontainer up` der Dev-Container-CLI statt `u-boot devcontainer up` | `05-lauf-host.sh` | derselbe Mechanismus darunter |
| Standardprofil: CLI mit dem **Host-Pfad** des Projekts eingehängt (Bind-Mount des Workspace braucht eine Quelle auf dem Docker-Host); kein Resolver-Protokoll (das Image trägt kein `dnsmasq`) | `05-lauf-host.sh` | im Standardprofil sind die Hosts, die der Container anfragt, **nicht gemessen** |
| Lauf 2: Container neu **und** Volumes des Projekts entfernt (frisches Volume); Images bleiben | `05-lauf-host.sh entfernen` | `devcontainer up --remove-existing-container` behielte die Volumes |

## Prompt-Zähler (Entscheidung des Projektinhabers: höchstens 40)

Gezählt wird jeder Aufruf, der einen Prompt an ein Cloud-Modell senden **könnte**, auch wenn er vorher an einer Meldung scheitert (`prompt_zaehlen` in `lib.sh`; Zähler außerhalb des Repos; die Obergrenze bricht mit Exit 99 ab, Gegenprobe
in `rohdaten-09-gegenproben.txt`). **Stand am Ende: 40 von 40 (Weg a: 22, Weg b: 18); die Obergrenze ist erreicht, nicht überschritten.**

| Nr. | Fall | Wirkung auf den Server |
|---|---|---|
| 1, 12 | `std` Lauf 1: `claude -p` Modell 1, Modell 2 | Antwort |
| 17, 18 | `std` Lauf 2: `claude -p` Modell 1, Modell 2 | Antwort |
| 2, 6 | `discover` Lauf 1: `claude -p` Modell 1, Modell 2 | Antwort |
| 3, 4 | `discover` Lauf 1: `ollama run`, `ollama launch claude` mit `glm-5.3-flash:cloud[1m]` | **Client lehnt den Namen ab** (`Error: invalid model name`, Exit 1), kein Verkehr; verbraucht im Zähler, nicht im Dienst |
| 5, 7 | `discover` Lauf 1: `run:o`, `launch:o` mit **leerem Modellnamen** (Fehler im Skript: `modell()` kannte den Schlüssel `o` nicht; `Error: model is required` und `headless --yes launch for claude requires --model <model>`, Exit 1) | kein Verkehr; **verschwendet**, Skript korrigiert (`modell_pruefen` bricht jetzt vor dem Zählen ab, Gegenprobe `rohdaten-05-std-1-111030.txt`) |
| 8, 9 | `discover` Lauf 1: `ollama run`, `ollama launch claude` mit Modell 2 | Antwort |
| 10, 11 | `discover` Lauf 1: `ollama run`, `ollama launch claude` mit `glm-5.3-flash:cloud` (ohne `[1m]`) | Antwort |
| 19, 20 | `discover` Lauf 2: `claude -p` Modell 1, Modell 2 | Antwort |
| 21, 22 | `discover` Lauf 2: `ollama run`, `ollama launch claude` mit `glm-5.3-flash:cloud` | Antwort |
| 13 | `restricted` (K3a) `claude -p` Modell 1 | **keine Antwort**, Exit 124 nach 120 s (die Verbindung zum Gateway läuft in die Zeitgrenze, `curl` Exit 28; Menge `allowed` leer) |
| 14 | `allow` (K3b) `claude -p` Modell 1 | **keine Antwort**, Exit 124 nach 120 s |
| 15, 16 | `c` (K3c, Adresse von Hand erlaubt) `claude -p` Modell 1, Modell 2 | Antwort |
| 23 bis 29 | Weg (b) `discover`: 23 `ollama run` Modell 1 (ohne `[1m]`), 24 `claude -p` Modell 1 (mit `[1m]`), 25 `claude -p` Modell 1 ohne `[1m]` (K6), 26 `ollama launch claude` Modell 1 (ohne `[1m]`), 27 `ollama run` Modell 2, 28 `claude -p` Modell 2, 29 `ollama launch claude` Modell 2 | Antwort |
| 30 bis 33 | Weg (b) `nocloud` und `apikey` mit dem Namen `…[1m]` (`run`, `launch` je Fall) | **Client lehnt den Namen ab**, kein Verkehr; wertlos, vom Entwurf des Skripts verschuldet (der Name mit `[1m]` war aus Weg (a) als ungültig bekannt) |
| 34, 35 | Weg (b) `OLLAMA_NO_CLOUD=1`: `ollama run`, `ollama launch claude` (ohne `[1m]`) | Meldungen „cloud is disabled“, Exit 1, kein Verkehr |
| 36 | Weg (b) `OLLAMA_API_KEY=platzhalter`: `ollama run` (ohne `[1m]`) | Antwort |
| 37 bis 40 | Weg (b) `restricted` mit Eintrag `ollama.com`: `ollama run` Modell 1, `claude -p` Modell 1, `ollama launch claude` Modell 1, `ollama run` Modell 2 | Antwort |

Aufteilung der 40 (Kriterien siehe oben; davon `claude -p` mit Kurzantwort 14 (Nr. 1, 2, 6, 12, 15 bis 20, 24, 25, 28, 38), `ollama run`/`launch` mit Exit 0 und Ausgabe ohne Fehlermeldung 14): Exit 0 mit Ausgabe 28 (Weg a: Nr. 1, 2, 6, 8 bis 12, 15 bis 22 = 16; Weg b: Nr. 23 bis 29, 36 bis 40 = 12), ohne Wirkung auf den Server 10 (Nr. 3, 4, 5, 7, 30 bis 35), durch die Egress-Restriktion blockiert 2 (Nr. 13, 14).
Nach Wegen: Weg (a) 22 (Plan: 20; die zwei Überschreitungen sind die verschwendeten Prompts 5 und 7), Weg (b) 18 (Plan: 14 + 2 + 2 + 2 Reserve = 20 abzüglich der 2, die Weg (a) aus der Reserve nahm). Vier Prompts (Nr. 30 bis 33) verbrauchte ich mit einem Namen, den Weg (a) schon als ungültig gezeigt hatte; deshalb entfielen in `restricted` `claude:2` und `launch:2` und im Fall `apikey` der `launch`.

## Aussagen Weg (a): Host-Server, Claude Code im Container

Notation: *Lauf* zählt die Läufe mit neuem Container und frischem Volume; `n/n` heißt: in allen Läufen gleich. Dateinamen sind Muster (`<HHMMSS>` in den Dateinamen der Rohdaten).
`discover` = Namen frei (Container ohne Egress-Skript), `restricted` = Liste des Generators (`res`), `allow` = Liste plus `host.docker.internal` (`alw`).

### K1: Erreichbarkeit des Host-Servers aus dem Container

| Aussage | Rohdaten | Läufe |
|---|---|---|
| Der Host-Server lauscht nur auf `127.0.0.1:11434`; vom Host-Netz-Container aus liefert `api/version` `{"version":"0.40.10"}` | `rohdaten-04-probe-113412.txt` | 1 |
| Ohne `--add-host` löst `host.docker.internal` im Container **nicht** auf (`getent`: Exit 2); mit `--add-host=…:host-gateway` ergibt es `172.17.0.1` (das Gateway des Netzes `bridge`) | `rohdaten-04-probe-113412.txt` | 1 |
| Mit `--add-host`, **ohne Relay**: `curl` auf `host.docker.internal:11434` scheitert (`Failed to connect … Could not connect to server`, Exit 7, 0 ms): die Annahme „der Dienst auf `127.0.0.1` ist über das Gateway nicht erreichbar“ ist **gemessen** | `rohdaten-04-probe-113412.txt`, `rohdaten-05-discover-{1-111102,2-112243}.txt` (Fall `K1-api-ohne-relay`) | 3 |
| **Mit Relay** (K1b): `api/version` liefert `0.40.10`; `ollama --version` im Container mit `OLLAMA_HOST` meldet die Version des **Servers** (`ollama version is 0.40.10`) und `Warning: client version is 0.40.6` | `rohdaten-05-discover-1-111102.txt`, `…-2-112243.txt` (Fälle `K1b-*`) | 2/2 |
| Abhängigkeit von der Engine: gemessen nur auf Linux mit lokaler Engine; die Adresse stammt aus `host-gateway` (Bridge-Gateway `172.17.0.1`), die Bindung des Dienstes auf `127.0.0.1` macht den Relay nötig. Docker Desktop, Colima, Podman, rootless Docker, `arm64`: **nicht gemessen** | — | — |

### K2: Standardprofil (`std`: nur `claude-code` und `node`, kein Ollama im Container)

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `claude --model glm-5.3-flash:cloud[1m] -p "<Prompt>"` mit `ANTHROPIC_BASE_URL=http://host.docker.internal:11434` und dem Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` (beide nur im Aufruf) endet mit Exit 0 in 1,0 bis 1,4 s und liefert eine Antwort (eine Zeile von höchstens zwei Zeichen); zusätzlich eine Zeile `[claude-code:…] {"model":"glm-5.3-flash:cloud[1m]","query_source":"sdk"}` | `rohdaten-05-std-1-110713.txt`, `rohdaten-05-std-2-112219.txt` (Fall `K-claude-m1`) | 2/2 |
| Dasselbe mit `gpt-oss:120b-cloud`: Exit 0 in 1,0 bis 1,5 s, Antwort; Claude Code meldet vorab `Claude Code doesn't recognize the model "gpt-oss:120b-cloud", so it assumes a context window of 200k tokens …` | `rohdaten-05-std-1-111134.txt`, `rohdaten-05-std-2-112219.txt` (Fall `K-claude-m2`) | 2/2 |
| Der Container trägt keine Zugangsvariable (`umgebung-container-verbotene-namen=0`); `ANTHROPIC_API_KEY` aus `remoteEnv` ist nicht gesetzt; der Platzhalter reicht | dieselben | 2/2 |
| Voraussetzung am Container: die **Handvariante** `--add-host` und ein **Relay** auf dem Host; ohne beides gibt es im Standardprofil keinen Weg zum Host-Server (K1) | K1 | — |

### K3: Sandbox-Profil

`discover` (Namen frei):

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `claude --model <m> -p` (Modell 1 und 2): Exit 0, Antwort, 0,9 bis 1,8 s | `rohdaten-05-discover-{1-110753,1-110835,2-112243}.txt` | 2/2 |
| `ollama run gpt-oss:120b-cloud "<Prompt>"` mit `OLLAMA_HOST` auf den Host-Server: Exit 0, Antwort (mehrere Zeilen, die Ausgabe des Modells ist länger als bei `claude -p`; Inhalt nicht aufgenommen) in 1,9 s; `ollama run glm-5.3-flash:cloud`: Exit 0, Antwort in 0,4 bis 0,8 s. **Eine Anmeldung im Container ist nicht nötig**; die Anmeldung liegt am Host-Server | `…-1-110835.txt` (Modell 2), `…-1-111102.txt`, `…-2-112243.txt` (Modell 1 ohne `[1m]`) | 1 (Modell 2), 2/2 (Modell 1) |
| `ollama launch claude --model <m> --yes` mit dem Prompt über stdin und `OLLAMA_HOST` auf den Host-Server: Exit 0 in 1,4 bis 2,0 s, Antwort; legt `~/.ollama/config.json` an (104/105 Byte, im Wurzel-Dateisystem des Containers, weil dort kein Volume liegt). Der Container hat keinen Server auf `127.0.0.1:11434`; die Antwort kam also vom Host-Server — **abgeleitet**, nicht an der Umgebung von `claude` gelesen | `…-1-110835.txt` (Modell 2), `…-1-111102.txt`, `…-2-112243.txt` | 1 (Modell 2), 2/2 (Modell 1) |
| Namen im Resolver-Log des Containers (K3d): **nur `api.anthropic.com`** (je eine Anfrage, bei `claude -p` und bei `launch`), bei `ollama run` keine; weder `ollama.com` noch ein Cloud-Host. Der Verkehr des Host-Servers zum Cloud-Dienst liegt außerhalb des Containers und ist **nicht** sichtbar | `rohdaten-05-discover-*` (Abschnitte „Namen im Resolver-Log“) | 2/2 |
| Dateien (Pfad, Größe, Rechte): `~/.claude.json` (852 Byte), `~/.claude/backups/…`, `~/.claude/projects/-home-vscode/<id>.jsonl` (131 bis 133 KB, Rechte 600), `~/.ollama/config.json`; alle im Wurzel-Dateisystem; kein Inhalt gelesen | dieselben (Abschnitte „Dateien neu/geaendert“) | 2/2 |

`restricted` und `allow` (Egress-Restriktion, Liste des Generators; der Host-Server liegt hinter dem Gateway `172.17.0.1`):

| Aussage | Rohdaten | Läufe |
|---|---|---|
| **K3a, Liste allein:** `host.docker.internal` löst über `/etc/hosts` auf `172.17.0.1` auf, der Resolver protokolliert zwei Anfragen mit `config error`; die nftables-Menge `allowed` ist leer; `curl` auf `host.docker.internal:11434` endet mit Zeitüberschreitung (Exit 28, 5 s), **auch mit Relay**; `ollama --version` mit `OLLAMA_HOST` läuft in die Zeitgrenze (Exit 124 nach 25 s). `claude -p`: **keine Antwort**, kein Ende innerhalb von 120 s (Exit 124), nur die Zeile mit dem Modellnamen | `rohdaten-05-restricted-1-111157.txt`, `…-2-112311.txt` (K1), `…-1-111639.txt` (Prompt) | 2/2 (K1), 1 (Prompt) |
| **K3b, Eintrag `host.docker.internal` in `egress.allow`** (erzeugt über `u-boot config set`, Name in der Liste und als `server=`/`nftset=` im Egress-Skript): der Resolver antwortet jetzt (`reply`, eine Anfrage), die Menge `allowed` bleibt **leer** (0 Elemente), `curl` Exit 28 mit und ohne Relay, `claude -p` Exit 124 nach 120 s ohne Antwort. Der Eintrag allein öffnet den Weg **nicht** (gemessen). Erklärung, nicht geprüft: die Adresse stammt aus `/etc/hosts` (Handvariante `--add-host`) und nicht aus einer Antwort des Resolvers; die Menge `allowed` füllt nur `nftset=` aus Antworten des Resolvers für den Namen | `rohdaten-05-allow-1-111256.txt`, `…-2-112411.txt` (K1), `…-1-111902.txt` (Prompt) | 2/2 (K1), 1 (Prompt) |
| **K3c, Handvariante A:** Zusatzzeile `address=/host.docker.internal/172.17.0.1` im Resolver (Projekt `alw`): der Resolver antwortet (`reply`), die Menge bleibt leer (0 Elemente), `curl` Exit 28 — die Zeile füllt die Menge nicht. (Im ersten Lauf lag die Zeile im Projekt `res` ohne Eintrag in der Liste: zwei Anfragen mit `config error`, Menge leer.) | `rohdaten-05-c-1-111355.txt` (`res`), `…-1-111518.txt` (`alw`) | 1 |
| **K3c, Handvariante B:** die Adresse `172.17.0.1` von Hand in die Menge eingetragen (`nft add element ip u_boot_egress allowed { 172.17.0.1 }`, Exit 0, 1 Element): `curl` auf `api/version` Exit 0 (`0.40.10`), `claude -p` mit Modell 1 Exit 0 in 8,1 s, mit Modell 2 Exit 0 in 1,2 s, je mit Antwort. Ist die Adresse erlaubt, trägt der Weg im Sandbox-Profil; erzeugen kann `u-boot` diese Erlaubnis **nicht** (K3b) | `rohdaten-05-c-1-112130.txt`, `…-2-112510.txt` (`curl` ohne Prompt) | 1 (Prompt), 2/2 (`curl`) |
| `ollama run` und `ollama launch` in `restricted`/`allow` mit dem Host-Server: **nicht gemessen** (Kontingent; `claude -p` zeigt dieselbe Sperre auf Verbindungsebene) | — | — |

### K6 (Teil): Modellname `…[1m]`

| Aussage | Rohdaten |
|---|---|
| `ollama run glm-5.3-flash:cloud[1m]` und `ollama launch claude --model glm-5.3-flash:cloud[1m] --yes`: `Error: invalid model name`, Exit 1, 0,0 s, **ohne Verkehr** — der Name mit `[1m]` ist für den Ollama-Client **kein gültiger Modellname** | `rohdaten-05-discover-1-110753.txt` (Nachtrag, siehe *Rohdaten-Eigenheiten*) |
| `claude --model glm-5.3-flash:cloud[1m] -p` liefert eine Antwort (Exit 0); Claude Code gibt den Namen mit `[1m]` aus. Ob Claude Code das Suffix vor dem Senden entfernt (Anthropic-Konvention für ein Kontextfenster von einer Million) und der Server `glm-5.3-flash:cloud` sieht, ist **nicht gemessen** (kein Server-Log des Host-Servers gelesen; Weg (b) liefert den Gegenbeleg) | `rohdaten-05-std-*`, `rohdaten-05-discover-*` |

## Aussagen Weg (b): angemeldeter Ollama-Server im Container (Sandbox-Profil)

Die Anmeldung führte der Projektinhaber am 2026-10-24 selbst an seinem Terminal aus (`07-anmeldung.sh`, Container `u-boot-messung-ocloud-konto-b` aus dem Image von `res`, Modell-Volume `…-b-modelle`); URL und Code
standen nur auf seinem Terminal, in keiner Datei und nicht im Gespräch. Server und Client im Container: Ollama `0.40.6`. Ein Lauf (Container neu, Volume neu). Die Messung nach der Anmeldung: `rohdaten-06-discover-1-*.txt`,
`rohdaten-06-restricted-1-114624.txt`; Vorbereitung `rohdaten-06-vorbereiten-1-114212.txt`; Rückabwicklung `rohdaten-06-ende-1-114643.txt`.

### K4: Anmeldung

| Aussage | Rohdaten |
|---|---|
| Vor der Anmeldung: der erste Start von `ollama serve` legt `id_ed25519` (387 Byte, Rechte 600) und `id_ed25519.pub` (81 Byte, 644) in `/home/vscode/.ollama` an, dazu `cache/model-recommendations.json` (2370 Byte); alle auf dem **Modell-Volume** (`df`-Einhängepunkt `/home/vscode/.ollama`, nicht im Wurzel-Dateisystem). Ein Workspace-Volume gibt es in diesem Container nicht (kein Pfad im Workspace) | `rohdaten-06-vorbereiten-1-114212.txt` |
| **Hosts des Anmeldeablaufs:** der Resolver des Containers zeichnete zwischen der Marke vor der Anmeldung und dem ersten Aufruf danach **keine** Anfrage auf (Log unverändert bei 10 Zeilen). Der erste Modellaufruf danach fragt `ollama.com` (2 Anfragen) und `registry.ollama.ai` (4) an | `rohdaten-06-discover-1-114452.txt` |
| **Dateien nach der Anmeldung:** gegenüber dem Zustand davor keine neue oder geänderte Datei (Pfad, Größe, Rechte von `id_ed25519` und `.pub` unverändert); die Anmeldung ändert den Geräteschlüssel im Volume also nicht sichtbar. Ob die Anmeldung als Zustand des Servers im Konto liegt, ist nicht gemessen | dieselbe |
| **Überleben der Anmeldung:** der Container wurde entfernt und mit **demselben Volume neu angelegt** (K5b, Modus `restricted`); `ollama run` lieferte ohne neue Anmeldung eine Antwort (Exit 0). Die Anmeldung hängt damit am Volume (am Geräteschlüssel `id_ed25519`), nicht am Container | `rohdaten-06-restricted-1-114624.txt` |
| **Geräteschlüssel und [`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer):** Ort: Modell-Volume `/home/vscode/.ollama` (Dateien `id_ed25519`, `id_ed25519.pub`); nicht im Workspace-Volume, nicht im Wurzel-Dateisystem; überlebt die Neuanlage des Containers; fehlt nach Entfernen des Volumes. Die Wertung (Zugangsdatum ja/nein) trifft die Folge-ADR | dieselben |
| Rückabwicklung: `ollama signout` im Container (`You have signed out of ollama.com`, Exit 0), Server beendet, Container und Volume entfernt (0 Objekte mit Präfix), Images entfernt. **Den Geräteschlüssel dieses Wegwerf-Containers in den Kontoeinstellungen entfernt der Projektinhaber** (Schritt offen bis zu seiner Bestätigung) | `rohdaten-06-ende-1-114643.txt`, `rohdaten-08-alles.txt` |

### K5 und K8: Läufe, `launch` in `restricted` gegen `discover`

| Aussage | Rohdaten | Läufe |
|---|---|---|
| `discover`: `ollama run` (Modell 1 ohne `[1m]`: 2,5 s; Modell 2: 1,3 s), `claude -p` (Modell 1 mit `[1m]` 1,4 s; ohne `[1m]` 1,3 s; Modell 2 1,2 s) und `ollama launch claude --model <m> --yes` (Modell 1 ohne `[1m]` 1,9 s; Modell 2 1,3 s): alle Exit 0 mit Antwort. Ein `ollama pull` ist nicht nötig; `run` zieht das Cloud-Modell selbst (Manifest 426 Byte und zwei Blobs von 264 und 307 bis 317 Byte, Namen `registry.ollama.ai`) | `rohdaten-06-discover-1-114452.txt`, `…-114507.txt` | 1 |
| Namen im Resolver-Log (`discover`): `ollama run` → `ollama.com` und `registry.ollama.ai`; `claude -p` und `launch` → nur `api.anthropic.com`; `run` ein zweites Mal → nur `registry.ollama.ai` | dieselben | 1 |
| `restricted` mit Liste des Generators **plus** Eintrag `ollama.com` (Handvariante: dieselben zwei Zeilen wie im Egress-Skript): `ollama run` (Modell 1 1,1 s, Modell 2 0,4 s), `claude -p` (1,0 s) und `ollama launch claude` (1,7 s) alle Exit 0 mit Antwort; im Resolver-Log nur `api.anthropic.com` (bei `claude` und `launch`). **K8:** in diesem Lauf verhält sich `launch` in `restricted` wie in `discover` (startet `claude`, Antwort); ohne `ollama.com` in der Liste ist `restricted` **nicht gemessen** (Kontingent) | `rohdaten-06-restricted-1-114624.txt` | 1 |
| Dateien: `~/.ollama/config.json` (104/105 Byte) und `~/.ollama/backup/config.json.<zeit>` (von `launch`), `cache/model-recommendations.json` (2041 Byte), Manifest und Blobs unter `models/`: alle auf dem **Modell-Volume**; `~/.claude.json`, `~/.claude/projects/…jsonl`: Wurzel-Dateisystem | `rohdaten-06-discover-1-*`, `rohdaten-06-ende-1-114643.txt` | 1 |
| Gemessen ist nur: **mit** dem Eintrag `ollama.com` (Handvariante) liefen die Aufrufe mit Exit 0. Im Resolver-Log dieser Aufrufe steht für `ollama run` und `ollama launch` **keine** Anfrage nach `ollama.com` oder `registry.ollama.ai` (nur `api.anthropic.com` bei `claude`/`launch`); dass die Aufrufe ohne neue Namensanfrage antworteten, ist **nicht erklärt** (naheliegend, nicht geprüft: der Server hielt Auflösung oder Verbindung aus dem Start der Sitzung). Ein Lauf **ohne** den Eintrag ist nicht gemessen; ob `ollama.com` für Weg (b) nötig ist, folgt aus dieser Messung **nicht** (der Vorgänger ohne Konto zeigte für `restricted` den Lookup-Fehler `ollama.com` bei `ollama run`, ohne Anmeldung) | `rohdaten-06-restricted-1-114624.txt` | 1 |

### K6: Modellname `…[1m]`

| Aussage | Rohdaten |
|---|---|
| `claude --model glm-5.3-flash:cloud[1m] -p` liefert eine Antwort und gibt **keinen** Hinweis zum Kontextfenster aus; `claude --model glm-5.3-flash:cloud -p` (ohne Suffix) liefert eine Antwort **mit** dem Hinweis `Claude Code doesn't recognize the model "glm-5.3-flash:cloud", so it assumes a context window of 200k tokens`. Gemessen ist nur dieser Unterschied der Ausgabe von Claude Code; was der Server als Modellname erhält und ob Claude Code das Suffix vor dem Senden entfernt, ist am Server nicht gelesen. `ollama run` und `ollama launch` lehnen den Namen mit `[1m]` ab (Weg a, Weg b) | `rohdaten-06-discover-1-114507.txt` |

### K7: Variablen

| Aussage | Rohdaten |
|---|---|
| `OLLAMA_NO_CLOUD=1` (Server und Client): `ollama run glm-5.3-flash:cloud` → `Error: ollama cloud is disabled: remote model details are unavailable`, Exit 1; `ollama launch claude --model … --yes` → `Warning: ignoring --model glm-5.3-flash:cloud because cloud is disabled` und `Error: headless --yes launch for claude requires --model <model>`, Exit 1. Beides ohne Verkehr, **trotz Anmeldung** | `rohdaten-06-discover-1-114559.txt` |
| `OLLAMA_API_KEY=platzhalter` (Server und Client, Platzhalter, kein gültiger Schlüssel): `ollama run glm-5.3-flash:cloud` liefert eine Antwort (Exit 0) wie ohne die Variable; keine geänderte Meldung. Die Wirkung eines **gültigen** Schlüssels ist nicht gemessen; der Platzhalter ändert das Verhalten bei bestehender Anmeldung nicht sichtbar. `ollama launch` mit der Variable ist nicht gemessen | dieselbe |
| Die ersten beiden Versuche beider Fälle (`rohdaten-06-discover-1-114526.txt`) benutzten den Namen mit `[1m]` und scheiterten am Client (`invalid model name`); sie sind **wertlos** und wurden mit dem Namen ohne Suffix wiederholt (vier verschwendete Prompts, Nr. 30 bis 33) | `rohdaten-06-discover-1-114526.txt` |


## Nicht gemessen

| Nicht gemessen | Bemerkung |
|---|---|
| Weg (b): `restricted` **ohne** den Eintrag `ollama.com`; `launch` und `claude -p` mit dem zweiten Modell in `restricted`; `ollama launch` unter `OLLAMA_API_KEY`; ein zweiter Lauf | Kontingent (40 Prompts erreicht) |
| Weg (b): Verhalten bei abgelaufener oder entfernter Anmeldung (nach dem Entfernen des Geräteschlüssels im Konto); ob der Server den Schlüssel im Konto prüft | nicht beobachtet; der Schlüssel wird erst nach dem Lauf entfernt |
| Weg (b): Verkehr des Servers im Container zum Cloud-Dienst | der Resolver des Containers sah beim Aufruf nur `ollama.com`/`registry.ollama.ai`; die Verbindungen selbst und weitere Hosts (Anmelde-Rückkanal) sind nicht beobachtet; die Anmeldung selbst erzeugte keine Namensanfrage im Container |
| `ollama run`/`ollama launch` hinter der Egress-Restriktion mit dem Host-Server | Kontingent; `claude -p` zeigt die Sperre |
| Hosts des Host-Servers zum Cloud-Dienst | liegen außerhalb des Containers |
| Hosts im Standardprofil | das Image trägt kein `dnsmasq`; kein Resolver-Protokoll |
| Größe und Zahl der Anfragen, die ein `claude -p` an den Dienst sendet; Kontingent-Verbrauch | nicht beobachtbar; ein Prompt kann mehrere Anfragen auslösen |
| Interaktiver Ablauf von Claude Code und von `ollama launch claude` (Terminal mit Tastatur) | alles headless |
| Wirkung eines gültigen `OLLAMA_API_KEY` | nur der Platzhalter wurde gemessen (Weg b, keine sichtbare Wirkung bei bestehender Anmeldung); mit gültigem Schlüssel nicht gemessen |
| Andere Modelle, andere Versionen von Ollama (Server `0.40.10`, Client `0.40.6`), Claude Code (`2.1.310`) | ein Lauf je Kombination |
| Ob Claude Code das Suffix `[1m]` vor dem Senden entfernt | siehe K6; am Server nicht gelesen |
| Docker Desktop, Colima, Podman, rootless Docker, `arm64` | Zusage-Menge ist Linux `amd64` mit lokaler Engine; ein anderer Host begründet keine Zusage |
| Mehrere gleichzeitige Nutzer des Host-Servers über den Relay; Absicherung des Relays | der Relay lief nur während der Fälle, ohne fremde Container im Netz |

## Nutzerpfad je Profil — Weg (a) (Eingabe für ADR und Validator)

| Profil | Schritt | Meldung / Ergebnis | Beleg |
|---|---|---|---|
| Standard | `u-boot config set devcontainer.features.claude-code.enabled true`, `…node.enabled true`, `u-boot generate devcontainer` | zwei Hinweise wie im Vorgänger; **kein** Eintrag für den Host-Server (Handvariante `--add-host` nötig) | `rohdaten-02.txt` |
| Standard | Relay auf dem Host starten (Handvariante) | — | `rohdaten-04-relay-start-110648.txt` |
| Standard | im Container: Variablen `ANTHROPIC_BASE_URL` und `ANTHROPIC_AUTH_TOKEN=ollama` setzen, `claude --model <m> -p "<Prompt>"` | Antwort in 1 bis 1,5 s; Hinweis zum Kontextfenster bei `gpt-oss:120b-cloud` | `rohdaten-05-std-*` |
| Sandbox, `discover` | wie Standard, zusätzlich `ollama run`/`launch` mit `OLLAMA_HOST` | Antwort; `[1m]` im Namen für `ollama run`/`launch` ungültig | `rohdaten-05-discover-*` |
| Sandbox, `restricted` | wie Standard, Liste des Generators | **keine Meldung**, `claude -p` hängt bis zur Zeitgrenze (120 s im Skript); `curl` Zeitüberschreitung | `rohdaten-05-restricted-*` |
| Sandbox, `restricted` | `u-boot config set devcontainer.sandbox.egress.allow host.docker.internal`, Neuanlage des Containers | unverändert: Zeitüberschreitung, Menge leer | `rohdaten-05-allow-*` |
| Sandbox, `restricted` | Adresse von Hand in die Menge eingetragen (Handvariante) | Antwort | `rohdaten-05-c-1-112130.txt` |
| Sandbox, Weg (b) | `ollama` und `claude-code` aktivieren, `generate devcontainer`, Container bauen (Image von `res`) | wie im Vorgänger; der Container enthält Ollama `0.40.6` | `rohdaten-02.txt` |
| Sandbox, Weg (b) | `ollama serve` im Container starten, `ollama signin` im Container | URL und Code nur auf dem Terminal des Projektinhabers (nicht aufgenommen); Bestätigung im Browser des Projektinhabers; im Container keine Namensanfrage und keine neue Datei | `rohdaten-06-discover-1-114452.txt` |
| Sandbox, Weg (b), `discover` | `ollama run <modell> "<Prompt>"`, `claude -p` mit `ANTHROPIC_BASE_URL=http://127.0.0.1:11434`, `ollama launch claude --model <m> --yes` | Antwort in 1 bis 2,5 s; `[1m]` im Namen nur für `claude` gültig | `rohdaten-06-discover-1-114452.txt`, `…-114507.txt` |
| Sandbox, Weg (b), `restricted` | Liste des Generators plus `ollama.com`, Container neu mit demselben Volume | Anmeldung bleibt, Antwort ohne neue Anmeldung | `rohdaten-06-restricted-1-114624.txt` |
| Sandbox, Weg (b) | `OLLAMA_NO_CLOUD=1` | „ollama cloud is disabled“, Exit 1 | `rohdaten-06-discover-1-114559.txt` |
| Sandbox, Weg (b) | Rückabwicklung: `ollama signout`, Volume entfernen, Schlüssel im Konto entfernen | `You have signed out of ollama.com`; der letzte Schritt liegt beim Projektinhaber | `rohdaten-06-ende-1-114643.txt` |

## Zugangsdaten und Rohdaten

- **Anmeldung im Container (Weg b):** nur der Projektinhaber, an seinem Terminal; die Anmelde-URL und der Code stehen in keiner Datei, keinem Skriptprotokoll und nicht im Gespräch; der Implementer hat die Anmeldung weder ausgeführt noch gesehen. Der Geräteschlüssel im Container wurde nie gelesen (nur Pfad, Größe, Rechte).
- **Zugang:** der Projektinhaber ist am Host-Server angemeldet; kein Skript liest, setzt oder reicht ein Zugangsdatum; **keine Datei des Hosts** (kein `~/.ollama`, kein `~/.claude`) wurde gelesen, kopiert oder eingebunden;
  der Host-Server wurde nur aus Containern über das Netz angesprochen. Der Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` steht nur im Aufruf im Container; die Umgebung des Hosts trug in allen Läufen
  keinen der geprüften Namen (`rohdaten-00-vorlauf.txt`, `rohdaten-01.txt`: `umgebung-pruefung exit=0`, Gegenprobe: sieben gesetzte Namen schlagen an).
- **Sicherung vor dem ersten Lauf (DoD 1):** (i) `redact()` mit Gegenprobe samt Artefakt (`rohdaten-00-vorlauf.txt`: 32 Formen und ein Schlüsselblock geschwärzt, 15 neutrale Zeilen unverändert; Hash des Stands;
  die Läufe brechen ohne das Artefakt oder bei anderem Stand ab); (ii) `umgebung_pruefen()` mit Ergebnis und Zeitstempel im Artefakt; **keine Ausnahme für `OLLAMA_API_KEY` nötig** (Platzhalter nur im Container-Aufruf von Weg (b));
  (iii) `09-suche.sh` mit benannter Ausnahmeliste (Platzhalter `ANTHROPIC_AUTH_TOKEN=ollama` und `OLLAMA_API_KEY=platzhalter`; die älteren Protokolle nach Pfad, je mit Grund im Skriptkopf); über dieses Verzeichnis und
  `tools/messung/ollama-cloud-konto/`: 0 Treffer; über ganz `docs/messungen/` und `tools/messung/` (`--alle`): 0 Treffer außerhalb der Pfad-Ausnahmen (37 Treffer in den benannten Pfaden der älteren Protokolle:
  Testnamen mit `Ollama_`, öffentliche Geräteschlüssel, GPG-Fingerabdrücke, `$DUMMY`-Zuweisungen — von Hand als kein Zugangsdatum gelesen); (iv) die Rohdaten sind in **zwei Stufen** entstanden (Stufe 1 außerhalb des Repos,
  Handlesung Zeile für Zeile, Stufe 2 durch `10-uebernehmen.sh --gelesen`); (v) **kein Antwort-Inhalt und keine Anmelde-URL** in den Rohdaten (Aufrufe mit Prompt geben nur Meldungszeilen aus; Weg (a) erzeugt keine Anmelde-URL).
- **Von Hand gelesen:** alle Dateien dieses Verzeichnisses, vor der Übernahme; die Zeile `[claude-code:[entfernt]]` ist eine Debug-Zeile von Claude Code, deren Zeichenkette hinter `claude-code:` die Schwärzung nach `code: <16 Zeichen>` traf (der Wert ist
  vermutlich eine Sitzungskennung, nicht gelesen); die erzeugten Dateien tragen nur `${localEnv:…}`-Verweise.
- **Dies ist der dritte Beleg** für [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../../plan/planning/observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md) (ein Lauf mit Zugang fand statt); die Closure trägt ihn ein.

## Gegenproben (rot gesehen)

`rohdaten-09-gegenproben.txt`: (1) je eine Regel von `redact()` entfernt (Schlüsselblock, lange Base64-Zeile, `ANTWORT`, JSON-Antwortfelder, `signed in as`, E-Mail, Anmelde-URL): die Probe endet jeweils mit Exit 1;
(2) die Vorlauf-Prüfung bricht ab, wenn das Artefakt fehlt oder `lib.sh` sich geändert hat; (3) `09-suche.sh` schlägt bei einer Token-Zuweisung an (Exit 1), auch wenn der benannte Platzhalter in derselben Zeile steht, und lässt
den Platzhalter allein durch (Exit 0); (4) `10-uebernehmen.sh` verweigert ohne `--gelesen` und entfernt Dateien, bei denen die Suche anschlägt; (5) die Prompt-Obergrenze bricht mit Exit 99 ab; (6) der Relay startet nicht,
solange ein fremder Container im Netz `bridge` hängt (Exit 3); (7) der unbekannte Modellschlüssel bricht vor dem Zählen ab (`rohdaten-05-std-1-111030.txt`); (8) `07-anmeldung.sh` bricht ohne Terminal ab (Exit 2, Konsole, keine Datei).
Das Wort „Schlüsselblock“ in der ersten Zeile von (1) steht dort ohne die Kopfzeile des Blocks, damit die Suche nicht anschlägt (von Hand ersetzt, nicht vom Skript geschrieben).

**Zeitbezug (ehrlich ausgewiesen):** `rohdaten-09-gegenproben.txt` trägt keine Zeitstempel; alle acht Läufe darin entstanden **nach Weg (a)**, also nach dem ersten Prompt (11:07:26 UTC), und sind für die Dokumentation neu gefahren. Vor dem ersten Prompt rot gesehen
und im Gespräch gelesen, aber nicht als Datei erhalten, wurden nur (1) die Regeln von `redact()` (Probe mit je einer entfernten Regel), (2) die Vorlauf-Prüfung (Artefakt fehlt; Stand geändert) und (3) die Suche. Die Gegenproben (4) bis (8) entstanden erst danach:
(5) die Prompt-Obergrenze wurde nur mit `PROMPT_MAX=22` bei Zählerstand 22 rot gesehen; **der Abbruch bei 40 wurde nie rot gesehen** (der 40. Prompt lief durch, die Grenze wurde erreicht, nicht überschritten); (7) der Wächter für den Modellschlüssel kam nach den
verschwendeten Prompts 5 und 7. Das Artefakt der Gegenprobe vor dem ersten Prompt gilt also für (1) bis (3), nicht für die übrigen.

**Änderung an einer Rohdatei:** `rohdaten-09-gegenproben.txt` wurde nach dem Schreiben von Hand und durch `redact()` verändert: Platzhalter-Werte erscheinen als `[entfernt]`, die Kopfzeile des Schlüsselblocks ist ersetzt (siehe oben), und Pfade des Messordners
und des Repos stehen als `<messordner>` und `<repo>` statt im Klartext (nachträglich von Hand ersetzt, damit keine umgebungsgebundenen Pfade des Agenten in der Rohdatei stehen). In `rohdaten-02.txt` steht der Repo-Pfad in der `make`-Zeile `Leaving directory` unverändert.

**Grenze der Suche (`09-suche.sh`):** Zeilen, die `[entfernt]`, `[Query entfernt` oder `PROBEWERT` tragen, werden übersprungen; eine Zeile mit einem echten Wert **und** einer solchen Marke würde die Suche nicht melden. Die Handlesung ersetzt das nicht. Die Rohdaten tragen die
öffentliche Form `https://ollama.com/connect?[Query entfernt: Anmelde-URL]` (aus der Gegenprobe der Schwärzung), keine Query, keinen Code, keinen Schlüssel.

## Aufräumen

`08-aufraeumen.sh weg-a` (`rohdaten-08-weg-a.txt`, `rohdaten-08-weg-a-nachher.txt`): entfernt Container, Volumes und Netze, deren Name `u-boot-messung-ocloud-konto` oder `vsc-u-boot-messung-ocloud-konto` beginnt, sowie die Container der CLI
über das Label; **die zehn Images mit dem Präfix bleiben** (Weg (b) braucht das Image von `res`; `08-aufraeumen.sh alles` entfernt sie). Fremde Objekte wurden nicht angefasst; die Zahlen vorher/nachher
(`container=0 images=85 volumes=899` gegen `container=1 images=95 volumes=914`) enthalten Arbeit anderer Prozesse auf dem Host (Volumes ohne Präfix, der Container ist der Filter der Schwärzung des laufenden Skripts oder ein fremder);
die zehn Mehr-Images sind die der Messung. Beim ersten Aufruf von `08-aufraeumen.sh` entfernte das Skript seinen eigenen Filter-Container der Schwärzung und brach ab (Skriptfehler, behoben: `lib.sh` lässt `<Präfix>-redact-…`
unberührt; Artefakt der Gegenprobe zum neuen Stand); der Aufruf wurde wiederholt.

`08-aufraeumen.sh alles` nach Weg (b) (`rohdaten-08-alles.txt`, `rohdaten-08-alles-nachher.txt`): `objekte-mit-praefix-nachher-alles=0`, `relay-container-nachher=0`; Images wieder 85 wie vorher. Containerzahl und Volumes (899 gegen 914) enthalten fremde Arbeit.
Der Geräteschlüssel des Wegwerf-Containers im Konto ist **nicht** entfernt: Schritt des Projektinhabers.

## Rohdaten-Eigenheiten (nur ausgewiesen)

- Weg (b): `rohdaten-06-vorbereiten-1-114212.txt` ist die Vorbereitung, die `07-anmeldung.sh` für den Projektinhaber aufrief; zwei frühere Testläufe von `vorbereiten` und `ende` (11:26 UTC, ohne Anmeldung) sind nicht übernommen. `nft set allowed nicht lesbar` im Modus `discover` ist erwartet (kein Egress-Skript).

- `rohdaten-05-discover-1-110753.txt` ist ein **Nachtrag**: die Rohdatei dieses Aufrufs wurde vom nächsten Aufruf desselben Teils und Laufs überschrieben (gleicher Dateiname; behoben: die Uhrzeit steht jetzt im Namen). Der Text ist die
  Konsolenausgabe des Aufrufs, wörtlich übernommen, die erste Zeile weist es aus. Der Aufruf `discover k1 nft` (11:06:40 UTC, vor dem Messen der Prompts) ist nicht erhalten; `K1` wurde im selben Lauf mit Datei wiederholt (`…-1-111102.txt`).
- Die Dateien `rohdaten-04-relay-{start,stop}-*.txt` sind bis auf das erste Paar nicht übernommen (gleiche Form, 18 Starts bis zum Ende von Weg (a)); der Nachweis „Relay entfernt“ steht in jeder Fall-Datei.
- `rohdaten-05-up-*` (Ende der Ausgabe von `devcontainer up`) sind nicht übernommen (Bauprotokoll, keine Messung).
- Der Fall `K-run-mo` trägt in den Dateien `…-1-110835.txt` den Befehl `ollama run  @PROMPT@` mit leerem Namen: der Fehler des Skripts (Prompts 5 und 7, siehe Zähler).
- Die Fälle heißen im Skript `K-*` (Weg a) und `B-*` (Weg b), ohne die Nummer des Messfalls; die Zuordnung steht in den Tabellen oben.

## Lesehilfen

`AUSSAGE <schlüssel>=<wert>` sind die stabilen Zeilen (`grep '^AUSSAGE' rohdaten-05-<teil>-<lauf>-*.txt`); Marken `---- <Zeit> <Fall>` trennen die Fälle. Zeiten der Aufrufe: `dauer=` (Sekunden).
