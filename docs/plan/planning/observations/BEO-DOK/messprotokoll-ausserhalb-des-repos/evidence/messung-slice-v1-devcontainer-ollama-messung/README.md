# Messprotokoll: Ollama im Sandbox-Devcontainer

Skripte, Dockerfiles und Rohdaten der Messung des Slice
[`slice-v1-devcontainer-ollama-messung`](../../../../../in-progress/slice-v1-devcontainer-ollama-messung.md)
(Ergebnisse und Grenzen: dort §3 *Messung*). Messstand: Docker 29.5.2 (Server, arm64) auf Colima 0.10.3
(VM Ubuntu 24.04.4, Kernel 6.8.0-117, 12 CPUs, 23,4 GiB), Docker-Client 29.9.1, Dev-Container-CLI 0.89.0,
Basis-Image `mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99`
(Debian 13, arm64), Ollama 0.40.6. Gemessen am 2026-10-18; Netzzugriff nur auf öffentliche Release-Seiten und Registries.

## Reihenfolge und Aufruf

Jedes Skript wird im eigenen Verzeichnis aufgerufen; die Ausgabe entspricht der Rohdatei in der Spalte rechts.
Voraussetzungen: Docker (Colima), `gh` mit Anmeldung, `jq`, `curl` auf dem Host; für Schritt 3 die Binaries aus `make build-binaries`
(`bin/`, nicht im Repo).

| Schritt | Skript | Rohdaten | Messfrage |
|---|---|---|---|
| 1 | `01-basis-und-optionen-c-d.sh` | `rohdaten-01.txt` | Messstand, Basis-Image, Option D (Debian-Paket), Option C (Docker-Hub-Image, Manifest) |
| 2 | `02-build-b.sh` mit `Dockerfile.b` | `rohdaten-02.txt` | Option B: Build, Download, Größe (mit CUDA-Verzeichnissen) |
| 3 | `03-projekt-erzeugen.sh <repo> <ziel>` | — | das Sandbox-Projekt mit Egress, erzeugt von `u-boot` |
| 4 | `04-build-probe.sh <ziel>` mit `Dockerfile.probe` | `rohdaten-04.txt` | Probe-Image (erzeugtes Dockerfile plus Ollama-Block ohne CUDA) |
| 5 | `05-laufzeit.sh` mit `inner-05.sh` | `rohdaten-05-discover.txt`, `rohdaten-05-restricted-1.txt` bis `-3.txt` | Laufzeit: Bindung, Pull, Run, Speicher, DNS-Mitschnitt; mit und ohne Egress-Restriktion |
| 6 | `06-option-a.sh` mit `Dockerfile.cli` | `rohdaten-06.txt`, `rohdaten-06-build-a.log` | Option A auf dem Basis-Image: Build bricht ab |
| 7 | `07-option-a-quelltext.sh` | `rohdaten-07.txt` | Feature-Quelltext (OCI-Artefakt 1.0.1) |
| 8 | `08-lieferkette.sh` | `rohdaten-08.txt` | Signatur, Provenance, Prüfsummen, Veränderlichkeit des Releases |
| 9 | `09-option-a-mit-zstd.sh` | `rohdaten-09.txt` | Option A mit zstd im Basis-Image: Build, Entrypoint, Benutzer |
| 10 | `10-volume-und-start.sh` | `rohdaten-10.txt` | Modellspeicher (Volume), Start nach dem Egress-Skript, Verhalten bei fehlendem `NET_ADMIN` |
| 11 | `11-archivinhalt.sh` | `rohdaten-11.txt` | entpackte Größe je Verzeichnis, beide Architekturen |
| 12 | `12-negativprobe-pin.sh` | `rohdaten-12.txt` | falscher SHA-256-Pin bricht den Build ab |

## Lesehilfen zu den Rohdaten

- `rohdaten-05-*`: Die Fortschrittsanzeige von `ollama pull` enthält Steuerzeichen; die Skripte filtern sie, die Zeilen mit
  `pulling …: 100%` bleiben. Die Zeile `config error` unter „DNS-Antworten“ ist ein Eintrag des Resolver-Protokolls (`dnsmasq`); sie erscheint nur in den
  Läufen, in denen ein benötigter Name nicht freigegeben war (der Aufrufer sieht `server misbehaving` beziehungsweise den Abbruch
  des Pulls). Ihr genauer Wortlaut ist nicht weiter untersucht.
- `rohdaten-05-restricted-2.txt`: Der Pull bricht ab, weil der Blob-Host fehlt; die 66 s sind die Wiederholungen des Clients.
- `rohdaten-06.txt`: Der Abbruch (`requires zstd`) ist das Messergebnis für Option A auf dem Basis-Image.
- `rohdaten-10.txt`, Abschnitt `S`: Je Richtlinie (`warn`, `fail`) läuft ein frischer Container.
- Zeitstempel in `rohdaten-*` sind UTC (Containerzeit); Dateizeiten des Dateisystems sind Ortszeit.
- Nicht gemessen und im Slice ausgewiesen: `amd64` (nur Archiv-Prüfsumme und -Inhalt, kein Lauf), Docker Desktop, Linux-Host,
  Auto-Start mit der Dev-Container-CLI (nachgestellt mit `docker exec`).
