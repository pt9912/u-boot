# Messprotokoll: Egress-Liste im laufenden Container

Skripte, Dockerfiles und Rohdaten der Messung des Slice
[`slice-v1-sandbox-egress-neuladen-entwurf`](../../../../../done/slice-v1-sandbox-egress-neuladen-entwurf.md)
(Ergebnisse und Grenzen: dort §3 *Messung*). Gemessen am 2026-10-18.

**Messstand:** Docker 29.5.2 (Server, linux/arm64) auf Colima 0.10.3 (VM Ubuntu 24.04.4, Kernel 6.8.0-117-generic, 12 CPUs),
Docker-Client 29.9.1, Dev-Container-CLI 0.89.0, Basis-Image
`mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99` (Debian 13.7, arm64),
`dnsmasq` 2.91, `nftables` 1.1.3, `sudo` 1.9.16p2. Das Projekt (`u-boot init --devcontainer --sandbox`, `egress.enabled: true`) erzeugt das
Binary dieses Repos (`u-boot version 0.11.0-dev`, `make build-binaries`, `bin/u-boot-linux-arm64`, Stand des Commits `073b3dd`); die
Messung läuft gegen das **erzeugte** Dockerfile und das **erzeugte** `egress-init.sh`, nicht gegen eine Nachbildung. Netzzugriff:
öffentliche Hosts der Default-Allowlist (`api.github.com`, `deb.debian.org`), `example.org`, `example.com`, `1.1.1.1` (ein UDP-Paket je
Sonden-Durchlauf an Port 9, bei offenem Pfad; und je ein TCP-Verbindungsversuch nach `https://1.1.1.1`).

## Voraussetzungen am Host (nicht im Repo)

- Docker (Colima) mit Zugriff auf `debian:stable-slim`, `mcr.microsoft.com/devcontainers/base:debian` und das Netz; `bash`, `mktemp`, `id`,
  `grep`, `awk`, `head`, `tail`, `cut` und `sed` am Host (`04-m3-zustellung.sh` und `05-m4-neustart.sh` nutzen `head`, `tail`, `cut` und `sed`
  in Pipelines auf dem Host; gemessen auf macOS, dessen `grep`, `awk` und `sed` nicht die GNU-Fassungen sind); keine Go-Toolchain, kein Python.
- **Docker-Socket des Hosts im Messcontainer:** `05-m4-neustart.sh` startet die Dev-Container-CLI als Container mit `-v /var/run/docker.sock:/var/run/docker.sock`
  (Image aus `Dockerfile.cli`); der Socket der Engine muss dort liegen.
- Das Binary `bin/u-boot-linux-arm64` aus `make build-binaries` (`bin/` steht nicht im Repo).
- **Images** (die Schritte 1 bis 3 brauchen die ersten beiden, Schritt 4 das dritte; Aufruf im Verzeichnis dieser Datei, nach Schritt 0):
  `docker build -t u-boot-messung-egress-neuladen <ziel>/.devcontainer` (das von Schritt 0 erzeugte Sandbox-Image),
  danach `docker build -f Dockerfile.messung -t u-boot-messung-egress-neuladen-dig .` (plus `dig`; das Produkt-Image hat kein `dig`)
  und `docker build -f Dockerfile.cli -t u-boot-messung-devcontainer-cli .` (Dev-Container-CLI).
- **Arbeitsverzeichnis außerhalb des Repos:** Schritt 0 und Schritt 4 erzeugen ihr Probe-Projekt in dem Ordner, den der Aufrufer wählt (bei der
  Messung `~/egress-neuladen-messung/proj` und `~/egress-neuladen-messung/m4`); beide löschen den übergebenen Ordner und legen ihn neu an. Schritt 4 legt sein
  Projekt selbst an und braucht Schritt 0 nicht. Die Schritte 1 bis 3 brauchen keinen Ordner.

## Reihenfolge und Aufruf

Jedes Skript wird im eigenen Verzeichnis aufgerufen; die Ausgabe entspricht der Rohdatei in der Spalte rechts.

| Schritt | Skript | Rohdaten | Messfrage |
|---|---|---|---|
| 0 | `01-projekt-erzeugen.sh <repo> <ziel>` | — | das Sandbox-Projekt mit Egress, erzeugt von `u-boot` (Eingabe für den Bau des Images in Schritt 0) |
| 1 | `02-m1-rechte.sh` mit `inner-m1.sh` | `rohdaten-02-m1.txt` | M1: was ein Prozess ohne und mit `sudo` an der Beschränkung ändern kann |
| 2 | `03-m2-neuladen.sh` mit `inner-m2.sh`, `probe.pl`, `reload-atomar.sh` | `rohdaten-03-m2.txt` (Abschnitte A bis G), `rohdaten-03-m2-wiederholung.txt` (vier weitere Läufe von A, vier von E) | M2: erneuter Lauf des Scripts: Fenster ohne Filter, DNS, Startrennen, laufende Verbindung, Set; atomare Variante |
| 3 | `04-m3-zustellung.sh` | `rohdaten-04-m3.txt` | M3: `docker cp` und `docker exec` vom Host, Eigentümer, gestoppter Container, Zustand nach `docker start` |
| 4 | `05-m4-neustart.sh <repo> <arbeitsverzeichnis>` | `rohdaten-05-m4.txt` | M4: Liste ändern und Container neu starten oder neu anlegen, mit der Dev-Container-CLI |

`03-m2-neuladen.sh` nimmt die Abschnittsbuchstaben als Argumente (Standard: `A E B F G C D`); die Wiederholung entstand mit
`03-m2-neuladen.sh A A A A E E E E`.

## Lesehilfen zu den Rohdaten

- **`reload-atomar.sh` ist Messvariante, kein Produktcode.** Es schreibt dieselben Regeln wie das erzeugte Script, in einer nft-Transaktion
  (`add table`, `delete table`, `add table` und alle Regeln in einer Datei, `nft -f`) und startet den neuen `dnsmasq` erst, wenn der alte beendet
  oder Zombie ist. Es ist kein Vorschlag für einen Wortlaut und nicht im Generator.
- **`probe.pl` zeigt Zustandswechsel mit Epochenzeit.** *egress* ist ein UDP-`send` an `1.1.1.1:9`: ein `drop` im OUTPUT-Haken liefert dem
  Aufrufer `EPERM`, ein offener Pfad Erfolg; „offen“ heißt, dass in diesem Augenblick kein Filter-Regelwerk den Verkehr hält. Die Sonde
  durchläuft etwa 180 000 Mal je Sekunde die Schleife (gemessen: `iterationen=` in der letzten Zeile der Egress-Sonde, rund 1,1 Millionen in 6 s), also alle
  rund 5 µs; ein Fenster, das kürzer wäre, bliebe unbeobachtet. *dns* fragt den lokalen `dnsmasq`; `timeout` ist 20 ms ohne Antwort (kann auch die Antwortzeit
  des Upstream-Resolvers bei abgelaufenem Cache sein: die Zeilen `dns timeout` und gleich danach `antwort`, vor dem Neu-Lauf, sind keine Lücke);
  `leseFehler (Connection refused)` heißt, kein Listener; `sendfehler (Operation not permitted)` heißt, dass die Loopback-Regel noch nicht stand.
- **Zombies in `rohdaten-03-m2.txt` und `rohdaten-03-m2-wiederholung.txt`:** Dort läuft der Container mit `sleep infinity` als Prozess 1, der keine
  Waisen einsammelt; der beendete `dnsmasq` bleibt als `<defunct>` stehen. Das ist ein Artefakt dieses Messaufbaus. `kill -0` meldet einen Zombie
  als lebend, deshalb wartet `reload-atomar.sh` auf den Prozesszustand `Z`. Unter der Dev-Container-CLI (`rohdaten-05-m4.txt`) gab es null Zombies.
- **`rohdaten-04-m3.txt`, Abschnitt 1 und 2:** `docker cp` von einer Datei des Hosts legte sie als `501 dialout` ab (UID und GID des Hosts, auf diesem Stand
  als `501` und `20` gelesen), nicht als `root`. Abschnitt 7 zeigt, dass der Eigentümer vom Aufruf abhängt (`docker cp -a`: `vscode`; Tar-Strom mit `owner=0`:
  `root`; mit `owner=1000`: `vscode`). Auf einem Host mit der UID `1000` kann ein einfacher `docker cp` dem Container-Benutzer `vscode` damit schreibbare Dateien
  liefern; gemessen ist nur dieser macOS-Host mit Colima.
- **`rohdaten-05-m4.txt`:** Der Wert hinter `Image` ist die Image-Kennung des Containers; `up --remove-existing-container` zeigt eine neue Kennung, `up` ohne die
  Option dieselbe wie zuvor. Die Zeilen `https://example.org (neuer Eintrag)` und `https://deb.debian.org` sind HTTPS-Aufrufe aus dem Container (`http=000`: nicht
  auflösbar oder nicht erreichbar).
- Zeitstempel in `rohdaten-*` sind UTC (Epochenzeit oder Containerzeit).

## Nicht gemessen

Docker Desktop, Linux-Host ohne VM, `amd64`, eine entfernte Container-Engine (`DOCKER_HOST`), VS Code (die Aussagen zu „Neustart“ und „Neu aufbauen“ gelten für die
Dev-Container-CLI 0.89.0; welche VS-Code-Aktion ihnen entspricht, ist nicht gemessen), ein Container mit `nestedRuntime`, die Last des Hosts als Auslöser des Startrennens
(50 Läufe sequenziell ohne Last), der Durchsatz eines Datenkanals über die Namensauflösung, mehrere Läufe je Fall bei M1, M3 und M4 (ein Lauf; M2 A: fünf, E: fünf).
