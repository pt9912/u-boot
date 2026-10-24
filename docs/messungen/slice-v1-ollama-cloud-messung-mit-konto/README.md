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

**Stand dieses Protokolls: Weg (a) gemessen, Weg (b) wartet auf die Anmeldung durch den Projektinhaber (Abschnitt *Wartepunkt*).** Jede Aussage unten gilt
**bedingungsgleich**: *mit dem Konto des Projektinhabers, auf einem Linux-`amd64`-Host mit Docker 29.9.2 (lokale Engine), ein Host-Server, zwei Modelle, eine Claude-Code-Version, am 2026-10-24
gemessen, headless (`claude -p`, kein interaktiver Ablauf)*. „Der Aufruf lieferte eine Antwort“ heißt nicht „Cloud-Modelle laufen“ und begründet keine Zusage.

**Antwort-Inhalte sind nicht aufgenommen.** Die Skripte weisen bei Aufrufen mit Prompt nur Exit-Code, Dauer, Zeilen- und Byte-Zahl der Ausgabe und die Zeilen aus, die nach einer Meldung aussehen
(Muster `MELDUNG` in `inner.sh`); „Antwort vorhanden“ in den Tabellen heißt: Exit 0 und eine Zeile ohne Meldungsmuster von höchstens zwei Zeichen (`claude -p`, Prompt: ein Satz mit der Bitte um ein Wort). Der Inhalt selbst
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
| 0 | `00-vorlauf.sh` (ruft `redact-probe.sh`) | `rohdaten-00-vorlauf.txt` | Gegenprobe der Schwärzung, Umgebungsprüfung samt Gegenprobe, Hash des Stands; **vor dem ersten Lauf** (Commit `6bfc17e8`, 2026-10-24 11:06:29 UTC; erster Prompt 11:07:26 UTC) |
| 1 | `01-messstand.sh` | `rohdaten-01.txt`, `rohdaten-01-vorher.txt` | Messstand, Namen der Umgebung, Docker-Objekte vorher (vor allen Bauten) |
| 2 | `02-projekte.sh` | `rohdaten-02.txt`, `rohdaten-02-erzeugt-<variante>.txt` | `u-boot generate devcontainer` für `std`, `res`, `alw`; Handvariante `--add-host`; Egress-Listen |
| 3 | `03-lauf-vorbereitung.sh` (`Dockerfile.cli`) | `rohdaten-03.txt` | Hilfs-Image der Dev-Container-CLI |
| 3b | `05-lauf-host.sh bauen` | `rohdaten-05-bauen-1.txt` | Images der drei Projekte (kein Prompt) |
| 4 | `04-relay.sh probe`, `start`, `stop` | `rohdaten-04-probe-*.txt`, `rohdaten-04-relay-{start,stop}-*.txt` | K1 ohne Relay; Relay für die Dauer eines Falls |
| 5 | `05-lauf-host.sh <teil> <fall>…` mit `inner.sh` | `rohdaten-05-<teil>-<lauf>-<HHMMSS>.txt` | Weg (a), K1 bis K3; Teile `std`, `discover`, `restricted`, `allow`, `c`; Läufe 1 und 2 |
| 6 | `06-lauf-container.sh` | `rohdaten-06-*` | Weg (b): **noch nicht gefahren** |
| 7 | `07-anmeldung.sh` | **keine** (nur Terminal des Projektinhabers) | Anmeldung im Container |
| 8 | `08-aufraeumen.sh weg-a` | `rohdaten-08-weg-a.txt`, `rohdaten-08-weg-a-nachher.txt` | Aufräumen von Container, Volumes, Netzen; die Images bleiben für Weg (b) |
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
in `rohdaten-09-gegenproben.txt`). **Stand nach Weg (a): 22 von 40.**

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

Aufteilung der 22: mit Antwort 16 (Nr. 1, 2, 6, 8, 9, 10, 11, 12, 15 bis 22 ohne 13 und 14), ohne Wirkung auf den Server (abgelehnter oder leerer Name) 4 (Nr. 3, 4, 5, 7), durch die Egress-Restriktion blockiert 2 (Nr. 13, 14).
Nach Wegen: Weg (a) 22 (Plan: 20; die zwei Überschreitungen sind die verschwendeten Prompts 5 und 7), Weg (b) 0 (Plan: 14 + 2 + 2 + 2 Reserve; für Weg (b) bleiben 18 von 40).

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

## Wartepunkt: Weg (b) — Anmeldung im Container

Weg (b) (K4 bis K8: Anmeldung des Ollama-Servers im Container, Hosts nach der Anmeldung, Dateien, `OLLAMA_NO_CLOUD` mit Konto, `OLLAMA_API_KEY` mit Platzhalter, `launch` in `restricted` gegen `discover`,
Neuanlage des Containers mit erhaltenem Volume) ist **nicht gemessen**. Die Anmeldung führt der Projektinhaber an seinem Terminal aus (`07-anmeldung.sh`); bis dahin stehen im Messordner das Image von `res`
(`vsc-u-boot-messung-ocloud-konto-res-…-uid`) und die Skripte bereit, ein Container läuft nicht. Der Test der Stufe `vorbereiten` (Container, Resolver mit Protokoll, `ollama serve` ohne Anmeldung) lief ohne Prompt:
Server `0.40.6` im Container, `id_ed25519` und `id_ed25519.pub` entstehen beim ersten Start des Servers im Modell-Volume (Rechte 600 und 644, 387 und 81 Byte). Die Rohdaten dieses Tests sind **nicht übernommen**
(Stufe 1, außerhalb des Repos); die Stufe läuft für die Messung erneut und schreibt ihre eigenen Rohdaten.

## Nicht gemessen — Stand nach Weg (a)

| Nicht gemessen | Bemerkung |
|---|---|
| Weg (b) insgesamt (K4 bis K8) | wartet auf die Anmeldung (Wartepunkt) |
| `ollama run`/`ollama launch` hinter der Egress-Restriktion mit dem Host-Server | Kontingent; `claude -p` zeigt die Sperre |
| Hosts des Host-Servers zum Cloud-Dienst | liegen außerhalb des Containers |
| Hosts im Standardprofil | das Image trägt kein `dnsmasq`; kein Resolver-Protokoll |
| Größe und Zahl der Anfragen, die ein `claude -p` an den Dienst sendet; Kontingent-Verbrauch | nicht beobachtbar; ein Prompt kann mehrere Anfragen auslösen |
| Interaktiver Ablauf von Claude Code und von `ollama launch claude` (Terminal mit Tastatur) | alles headless |
| Wirkung eines gültigen `OLLAMA_API_KEY` | nur Platzhalter vorgesehen; nicht gemessen |
| Andere Modelle, andere Versionen von Ollama (Server `0.40.10`, Client `0.40.6`), Claude Code (`2.1.310`) | ein Lauf je Kombination |
| Ob `[1m]` von Claude Code vor dem Senden entfernt wird | siehe K6 |
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

## Zugangsdaten und Rohdaten

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

## Aufräumen

`08-aufraeumen.sh weg-a` (`rohdaten-08-weg-a.txt`, `rohdaten-08-weg-a-nachher.txt`): entfernt Container, Volumes und Netze, deren Name `u-boot-messung-ocloud-konto` oder `vsc-u-boot-messung-ocloud-konto` beginnt, sowie die Container der CLI
über das Label; **die zehn Images mit dem Präfix bleiben** (Weg (b) braucht das Image von `res`; `08-aufraeumen.sh alles` entfernt sie). Fremde Objekte wurden nicht angefasst; die Zahlen vorher/nachher
(`container=0 images=85 volumes=899` gegen `container=1 images=95 volumes=914`) enthalten Arbeit anderer Prozesse auf dem Host (Volumes ohne Präfix, der Container ist der Filter der Schwärzung des laufenden Skripts oder ein fremder);
die zehn Mehr-Images sind die der Messung. Beim ersten Aufruf von `08-aufraeumen.sh` entfernte das Skript seinen eigenen Filter-Container der Schwärzung und brach ab (Skriptfehler, behoben: `lib.sh` lässt `<Präfix>-redact-…`
unberührt; Artefakt der Gegenprobe zum neuen Stand); der Aufruf wurde wiederholt.

## Rohdaten-Eigenheiten (nur ausgewiesen)

- `rohdaten-05-discover-1-110753.txt` ist ein **Nachtrag**: die Rohdatei dieses Aufrufs wurde vom nächsten Aufruf desselben Teils und Laufs überschrieben (gleicher Dateiname; behoben: die Uhrzeit steht jetzt im Namen). Der Text ist die
  Konsolenausgabe des Aufrufs, wörtlich übernommen, die erste Zeile weist es aus. Der Aufruf `discover k1 nft` (11:06:40 UTC, vor dem Messen der Prompts) ist nicht erhalten; `K1` wurde im selben Lauf mit Datei wiederholt (`…-1-111102.txt`).
- Die Dateien `rohdaten-04-relay-{start,stop}-*.txt` sind bis auf das erste Paar nicht übernommen (gleiche Form, 18 Starts bis zum Ende von Weg (a)); der Nachweis „Relay entfernt“ steht in jeder Fall-Datei.
- `rohdaten-05-up-*` (Ende der Ausgabe von `devcontainer up`) sind nicht übernommen (Bauprotokoll, keine Messung).
- Der Fall `K-run-mo` trägt in den Dateien `…-1-110835.txt` den Befehl `ollama run  @PROMPT@` mit leerem Namen: der Fehler des Skripts (Prompts 5 und 7, siehe Zähler).
- Die Fälle heißen im Skript `K-*` (Weg a) und `B-*` (Weg b), ohne die Nummer des Messfalls; die Zuordnung steht in den Tabellen oben.

## Lesehilfen

`AUSSAGE <schlüssel>=<wert>` sind die stabilen Zeilen (`grep '^AUSSAGE' rohdaten-05-<teil>-<lauf>-*.txt`); Marken `---- <Zeit> <Fall>` trennen die Fälle. Zeiten der Aufrufe: `dauer=` (Sekunden).
