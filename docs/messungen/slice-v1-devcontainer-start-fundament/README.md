# Messprotokoll: Hilfs-Image des Devcontainer-Starts

Rohdaten und Ablauf der Einmal-Messung des Slice
[`slice-v1-devcontainer-start-fundament`](../../plan/planning/in-progress/slice-v1-devcontainer-start-fundament.md)
(DoD: Größe, Bauzeit kalt und warm, Dauer des ersten `up`, Trivy HIGH/CRITICAL). Skripte liegen unter
`tools/messung/devcontainer-start-fundament/` (Ablage: [`docs/messungen/README.md`](../README.md)).

**Messstand.** Docker 29.9.1 (Client und Server), Linux-Host `x86_64`, Server `linux/amd64`, ein Lauf je Zahl. Dev-Container-CLI 0.89.0 im Image.
Gemessen am 2026-10-20. Das Hilfs-Image ist das aus dem eingebetteten Kontext gebaute (`internal/adapter/driven/docker/devcontainercli/`,
Basis-Images und CLI als Pins im Dockerfile). **Nicht gemessen:** `arm64` (Colima), Docker Desktop, ein kalter Pull der beiden Basis-Images des
Hilfs-Images (beide waren lokal vorhanden), der Weg über den Go-Adapter (Messung 2 ruft `docker run` mit der Argumentliste des Adapters auf; der
Adapter selbst läuft in den Docker-Tests des Slice). Keine Zahl dieses Protokolls ist eine Zusage.

## Reihenfolge und Aufruf

Jedes Skript wird im Verzeichnis `tools/messung/devcontainer-start-fundament/` aufgerufen; seine Ausgabe wird in die Rohdatei der Tabelle
geschrieben. `P` steht für den Pfad dieses Verzeichnisses, von dort aus gesehen `../../../docs/messungen/` plus Verzeichnisname.

**Voraussetzung:** ein Linux-Host mit Docker, Bash und den GNU-Werkzeugen, die die Skripte aufrufen (`tar` mit `--mtime`, `--owner` und `--group`,
`date +%N`, `stat -c`, `awk`, `sed`, `grep`, `cut`, `mktemp`, `head`, `tr`); Netzzugriff auf die Registries, auf `mirror.gcr.io` (Trivy-Datenbank) und für
den Pull des Trivy-Images. Gelaufen ist es nur auf Linux (`x86_64`). Auf macOS laufen `date +%N`, `stat -c` und die `tar`-Optionen nicht; ein Lauf dort
ist nicht belegt und braucht GNU-Werkzeuge oder eine Anpassung der Skripte.

| Schritt | Aufruf | Rohdaten | Messfrage |
|---|---|---|---|
| 1 | `bash 01-hilfs-image.sh ../../.. > $P/rohdaten-01-hilfs-image.txt 2>&1` | `rohdaten-01-hilfs-image.txt` | Plattformen der Basis-Images, Bauzeit kalt und warm, Größe, Trivy |
| 2 | `bash 02-erster-up.sh ../../.. /tmp/<arbeitsverzeichnis> > $P/rohdaten-02-erster-up.txt 2>&1` | `rohdaten-02-erster-up.txt` | Dauer von `devcontainer up`: erster Start mit kaltem Basis-Image des Projekts, Wiederholung, Neustart, Neuanlage |

Die Rohdaten tragen keine Zugangsdaten; die Skripte setzen keine Umgebung durch.

## Ergebnisse (Einzelmessung, `amd64`, ein Lauf)

| Größe | Wert | Quelle |
|---|---|---|
| Größe des Hilfs-Images | 1 799 126 108 Byte (1716 MiB) | Schritt 1, Abschnitt 3 |
| Bauzeit kalt (`--no-cache`, Basis-Images lokal) | 7,8 s | Schritt 1, Abschnitt 2 |
| Bauzeit warm, Tag fehlt (Layer im Cache) | 1,7 s | Schritt 1, Abschnitt 2 |
| Bauzeit warm, Tag vorhanden (Bau ohne Änderung) | 1,2 s | Schritt 1, Abschnitt 2 |
| Erster `up` (Basis-Image `base:ubuntu` wird gezogen, Container wird angelegt) | 52,0 s | Schritt 2, Abschnitt 1 |
| `up` gegen den laufenden Container | 1,2 s | Schritt 2, Abschnitt 2 |
| `up` nach `docker stop` (gleicher Container) | 1,8 s | Schritt 2, Abschnitt 3 |
| `up --remove-existing-container` (neuer Container) | 4,6 s | Schritt 2, Abschnitt 4 |
| Trivy 0.70.0, HIGH/CRITICAL, Betriebssystem-Pakete (Debian 12.15) | 597 Funde (davon 21 CRITICAL) | Schritt 1, Abschnitt 4 |
| Trivy, HIGH/CRITICAL, übrige Ziele mit Funden | `docker-buildx` (Go-Binary) 7; Pakete von `npm` im Node-Image: `brace-expansion` 5 und je 1 in `@npmcli/metavuln-calculator`, `http-cache-semantics`, `ip-address`, `pacote`, `picomatch`, `sigstore`; die Dev-Container-CLI selbst 0 | Schritt 1, Abschnitt 4 |

## Lesehilfen

- **Streuung.** Ein früherer Lauf von Schritt 1 (die Rohdatei wurde danach überschrieben, weil sie die ungefilterte Trivy-Tabelle mit
  rund 550 KB trug) ergab kalt 19,2 s, warm mit fehlendem Tag 2,3 s und warm mit Tag 1,5 s; die Größe war gleich bis auf 147 Byte. Die kalte Bauzeit hängt
  am Netz zur npm-Registry und schwankt damit stark.
- **„Kalt“** heißt hier ohne Layer-Cache des Bauens, nicht ohne die Basis-Images (Node-Image, Docker-CLI-Image): deren Pull ist nicht in der
  Bauzeit. Das Node-Image (`node:22-bookworm`) ist das volle Image und trägt die Funde des Betriebssystems; der Dev-Container-CLI-Anteil im Image
  (`opt/devcontainer-cli`) hat 0 Funde.
- **Erster `up`** (Schritt 2): Das Projekt nutzt `mcr.microsoft.com/devcontainers/base:ubuntu`, das das Skript vorab nicht lokal haben darf; die 52,0 s
  enthalten den Pull dieses Images und den Bau des UID-Images der CLI. Der Wert gilt für dieses Netz und diesen Host.
- **Pull des Trivy-Images** und der **Trivy-Datenbank** laufen im Skript mit; ihre Dauer steht nicht in der Tabelle.
- **Schritt 2 vs. Adapter.** `up` in `02-erster-up.sh` hängt den Projektordner schreibend ein und reicht keine Umgebung durch, wie die Voreinstellung des
  Adapters; die Docker-Tests des Slice führen dieselbe Argumentliste über den Go-Adapter aus.
