# Messprotokoll: Ollama im Sandbox-Devcontainer auf `amd64`

Rohdaten (diese Dateien) und Skripte
([`tools/messung/ollama-amd64/`](../../../tools/messung/ollama-amd64/)) der Messung des Slice
[`slice-v1-devcontainer-ollama-amd64-messung`](../../plan/planning/in-progress/slice-v1-devcontainer-ollama-amd64-messung.md)
(Messfragen A0 bis A10; Ergebnisse und Befunde: dort §7). Die Messung ist **keine** Verifikation im Sinn von
[`ADR-0024`](../../plan/adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10: Sie ist ein Lauf auf einem Host und kein CI-Lauf.

## Messstand

| Größe | Wert (Fundstelle) |
|---|---|
| Datum | 2026-10-19 (`rohdaten-01.txt`) |
| Host | Linux x86_64, Kernel 6.8.0-139-generic, 20 CPUs (13th Gen Intel Core i9-13900H), 31,1 GiB RAM, native Docker Engine, kein Colima, keine VM |
| Docker | Client 29.9.1 und Server 29.9.1 (`linux/amd64`), Storage-Treiber `overlay2`, cgroup v2 |
| Host-Schalter | `kernel.apparmor_restrict_unprivileged_userns = 0`; `/dev/fuse` und `/dev/net/tun` vorhanden |
| Basis-Image | `mcr.microsoft.com/devcontainers/base:debian` (Debian 13), lokale Kopie vom 2026-09-10, `sha256:1f851004adcd3dff3776b4a1da86727cf52280b0fdc9d5b0568c9c7d74274286`; die Skripte ziehen es nicht neu |
| `u-boot` | `make build IMAGE=u-boot-messung-amd64-ub` aus dem Stand `7956c49` (Slice-Start), `u-boot version 0.11.0-dev` (`rohdaten-02.txt`) |
| Ollama | 0.40.6; SHA-256 des `amd64`-Archivs `c164ab8b…c5e0` (Pin der Vorlage; das heruntergeladene Archiv stimmt überein, `rohdaten-04.txt`) |
| Plattenplatz | 39 GB frei vor, 35 GB nach der Messung (`rohdaten-01.txt`, `rohdaten-11.txt`; der Host wird von anderen Arbeiten mitbenutzt, die Differenz ist nicht allein die der Messung) |
| Netzstandort | ein Host an einem Anschluss; der Ort ist nicht erfasst und kein Maß für die Stabilität des Blob-Hosts |
| Läufe | **ein** vollständiger Lauf (`00-alles.sh`) für die Rohdaten dieses Verzeichnisses, ein zweiter Lauf derselben Skripte zum Vergleich (Abschnitt *Wiederholung*); A7: 30 Läufe des Startscripts |

Zwei Größen für das Archiv und das Image: **Download** = `Content-Length` des Archivs, 1 438 364 593 Bytes (1,44 GB, `rohdaten-03.txt`,
`rohdaten-04.txt`); **Image** = `.Size` des gebauten Images, 1 013 158 921 Bytes (1,01 GB) mit dem Ollama-Block gegen 946 749 835 Bytes (0,95 GB)
ohne ihn (`rohdaten-03.txt`); entpackt sind es 61 MiB mit den Mustern und 2151 MiB ohne sie (`rohdaten-04.txt`).

## Voraussetzungen

Auf dem Host nur `docker` (Daemon erreichbar), `make` und die Shell mit ihren Standardwerkzeugen (`bash`, `date`, `tee`, `grep`, `sort`,
`diff`, `mkdir`, `cp`, `ls`, `wc`); jede Auswertung (`jq`, `sed`, `awk`, `sha256sum`, `tar`, `curl`) läuft in einem Container. Netzzugriff auf
`github.com`, `registry.ollama.ai`, Docker Hub und `deb.debian.org`; Plattenplatz für rund fünf Images zu je 1 bis 1,4 GB, das Archiv (1,44 GB) und
ein Entpacken (rund 2,2 GB). Ein Messordner **außerhalb des Repos** (`WORK`): Das Archiv, die erzeugten Projekte und Zwischenstände liegen dort.

## Aufruf

```bash
export WORK=<messordner-ausserhalb-des-repos>
tools/messung/ollama-amd64/00-alles.sh               # alle Schritte der Reihe nach
tools/messung/ollama-amd64/03-bau.sh                 # oder ein Schritt allein (Voraussetzungen unten)
OUT=<anderes-verzeichnis> tools/messung/ollama-amd64/00-alles.sh   # zweiter Lauf, Rohdaten an anderen Ort
tools/messung/ollama-amd64/vergleich.sh <ausgabe-1> <ausgabe-2>    # stabile Aussagen zweier Läufe vergleichen
```

`OUT` (Vorgabe: dieses Verzeichnis) ist der Ort der Rohdaten. Alle Docker-Objekte tragen das Präfix `u-boot-messung-amd64-`
(die Images der Dev-Container-CLI `vsc-u-boot-messung-amd64-…`). Jeder Schritt ist wiederholbar: Er entfernt zuerst seine Objekte
und legt sie neu an. Fehler brechen die Schritte nicht ab (`set -u`, kein `set -e`): Ein roter Teilschritt soll die übrigen Messungen
nicht verdecken; der Befund steht im Rohdatum.

| Schritt | Skript | Voraussetzung | Rohdaten | Messfrage |
|---|---|---|---|---|
| 1 | `01-messstand.sh` | — | `rohdaten-01.txt`, `rohdaten-11-vorher.txt` | A0, A10 (Liste vorher) |
| 2 | `02-projekte.sh` | 1 | `rohdaten-02.txt`, `rohdaten-02-erzeugt-<variante>.txt` | `make build`; sechs erzeugte Projekte: `oll` (ollama), `bas` (ohne), `pod`/`podn` (podman mit/ohne ollama), `dok`/`dokn` (docker mit/ohne ollama) |
| 3 | `03-bau.sh` | 2 | `rohdaten-03.txt`, `rohdaten-03-kalt.log`, `-cache.log`, `-bas.log` | A1 |
| 4 | `04-archiv.sh` (`Dockerfile.tar`) | 3 | `rohdaten-04.txt` | A2 |
| 5 | `05-laufzeit.sh` (`inner-05.sh`) | 3 | `rohdaten-05-discover.txt`, `rohdaten-05-restricted.txt` | A3, A4 |
| 6 | `06-volume.sh` | 3 | `rohdaten-06.txt` | A5 |
| 7 | `07-kombination.sh [podman\|docker]` | 2 | `rohdaten-07-podman.txt`, `rohdaten-07-docker.txt`, `rohdaten-07-<variante>.log` | A6 |
| 8 | `08-wiederholungslauf.sh [läufe]` (`inner-08.sh`) | 3 | `rohdaten-08.txt` | A7 (30 Läufe) |
| 9 | `09-tests.sh` | 2 | `rohdaten-09.txt` | A8 |
| 10 | `10-devcontainer-cli.sh` (`Dockerfile.cli`) | 2 | `rohdaten-10.txt` | A9 |
| 11 | `11-aufraeumen.sh` | 1 | `rohdaten-11.txt`, `rohdaten-11-nachher.txt` | A10 |

`11-aufraeumen.sh` bricht ab, wenn `rohdaten-11-vorher.txt` fehlt oder leer ist; ohne diese Liste entfernt es nichts. Es entfernt nur
Objekte mit dem Präfix, unbenannte Images, die in der Liste vorher fehlen, das von der Messung gezogene Image `docker:29-cli` und den
Zwischencontainer des fehlgeschlagenen Baus der Negativprobe (A8). Der Build-Cache und fremde Objekte bleiben.

## Wiederholung

Jede Messfrage schreibt stabile Aussagen als Zeilen `AUSSAGE <schlüssel>=<wert>`. Zwei vollständige Läufe auf demselben Host
ergaben dieselben 86 Aussagen (`vergleich.sh`, Exit 0). Ausgenommen sind zwei Gruppen, deren Zahlen je Lauf wechseln dürfen: die Zählungen
der Wiederholung A7 (26 und 28 Läufe mit Fenster in den beiden Läufen) und der Pull des verschachtelten Docker hinter der Restriktion
(A6, zeitabhängig: `rohdaten-07-docker.txt` zeigt einen Lauf, dessen erster Pull mit `i/o timeout` scheiterte und dessen Wiederholungsversuch nach 15 s
gelang; im zweiten, verglichenen Lauf gelang der erste Versuch. In zwei früheren Läufen derselben Skripte, deren Rohdaten die spätere Fassung des Skripts überschrieb,
scheiterte der erste Versuch ebenfalls).

## Lesehilfen

- `rohdaten-02-erzeugt-<variante>.txt`: die Dateien, die `u-boot generate devcontainer` erzeugt hat, als Text (Abschnitte `=== .devcontainer/<Datei> ===`).
  Das Dockerfile des Slice-Gegenstands steht in `rohdaten-02-erzeugt-oll.txt`.
- `rohdaten-05-*`: Die Fortschrittsanzeige von `ollama pull` enthält Steuerzeichen; die Skripte filtern sie. Im Modus `restricted` ersetzt ein
  protokollierender `dnsmasq` mit derselben Konfiguration den Resolver des erzeugten Startscripts (Verfahren des Vorgänger-Protokolls
  der `arm64`-Messung unter Colima); Hosts und Antworten sind dadurch mitgeschnitten.
- `rohdaten-08.txt`: Die Sonden `B1` bis `B8` (acht versetzte Schleifen, zusammen rund alle 7 ms ein Versuch) verbinden sich mit dem Gateway des Containers (Port 9, Frist 50 ms); ein Ausgang ungleich 124
  (Frist) heißt: die Verbindung wurde nicht verworfen, die Restriktion galt in dem Moment nicht. Das Fenster ist die Spanne vom ersten bis zum letzten
  Treffer eines Laufs und damit eine **Untergrenze** (zwischen zwei Sonden kann eine Lücke liegen). Es gibt keinen Treffer außerhalb der Läufe. Die
  Sonden decken IPv4; die IPv6-Tabelle des Skripts ist nicht gemessen.
- `rohdaten-10.txt`: Die Dev-Container-CLI vergibt Container und Images selbst; das Skript findet sie über das Label `devcontainer.local_folder`.
  Der erste `up` dauerte 2,7 s, weil das Image aus dem Docker-Cache kam; der zweite mit `--build-no-cache` baute das Image neu (50,6 s).
- Zeitstempel in den Rohdaten sind UTC (Containerzeit und `date -u`); Dateizeiten des Dateisystems sind Ortszeit.
- Die Listen `rohdaten-11-vorher.txt` und `rohdaten-11-nachher.txt` enthalten Namen aller Docker-Objekte des Hosts. Ihr Unterschied in
  `rohdaten-11.txt` enthält auch Änderungen fremder Arbeiten auf dem Host während der Messung (neu gebaute Images, anonyme Volumes,
  Container), die nicht aus dieser Messung stammen.

## Schwärzung

Die Rohdaten sind per Filter (`redact` in `lib.sh`) auf vorsignierte URLs (`X-Amz…`, `Signature`, `Credential`), `Authorization`/`Bearer`-Werte und den
Home-Pfad des Hosts geschwärzt und danach per Suche über das Verzeichnis geprüft (Suchmuster: `X-Amz`, `Authorization`, `Bearer`, `token`, `secret`,
`password`, `Signature=`, `Credential=`, `sk-ant`, `ghp_`, Schlüsselblöcke, Home-Pfad, URLs mit Query): Es gibt keinen Treffer außer Variablennamen und
Platzhaltern in den erzeugten Dateien (`GIT_TOKEN`, die Referenz `${localEnv:GIT_TOKEN}`) und dem Hinweis des Docker-Linters zu `PODMAN_SIGNING_KEY`
(der öffentliche Fingerabdruck). Das Schlüsselpaar `id_ed25519` des Volumes steht nirgends; `06-volume.sh` ermittelt nur die Prüfsumme des öffentlichen Teils.

## Nicht gemessen

CI (`make test-docker` ohne `-v`; der Lauf in `rohdaten-09.txt` ist keiner), Docker Desktop, Colima/`arm64` durch diese Messung, Linux hinter einer VM,
größere Modelle und Mindest-Speicher, `OLLAMA_NO_CLOUD`, Cloud-Modelle, die Stabilität des Blob-Hosts über Zeit und Region, ein Start des Servers durch
`postStartCommand`, die Aktion eines Editors (VS Code) zum Neuaufbau, die IPv6-Tabelle des Startscripts, ein zweiter Host.
