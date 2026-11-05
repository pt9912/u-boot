# Messung: Volumes eines Devcontainers und das Verhalten von `docker rm` und `docker volume rm` (Grundlage für `devcontainer remove`)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-devcontainer-remove-messung`](../../plan/planning/in-progress/slice-v1-devcontainer-remove-messung.md) (erster von vier Slices zu `devcontainer remove`; Grundlage: [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Volume-Namen je Instanz: [`ADR-0015`](../../plan/adr/0015-sandbox-volumes-pro-instanz.md)) |
| Gelesen am  | 2026-11-05 |
| Skripte     | `tools/messung/slice-v1-devcontainer-remove-messung/` (`01-m1-volumes.sh`, `02-m2-entfernen.sh`, `03-m3-instanz.sh`, `04-m4-ueberleben.sh`, `05-aufraeumen.sh`, gemeinsam `lib.sh`; `00-alles.sh` ruft 01 bis 04) | <!-- d-check:ignore (Skriptpfad, keine Slice-Kennung) -->

Die Messung fährt den ausgelieferten Befehl `u-boot devcontainer up` und `down` (kein Nachbau des Aufrufs der Dev-Container-CLI) und die
Docker-Befehle, mit denen ein Nutzer heute aufräumt. Sie ändert keinen Produktionscode.

## Umgebung

Linux-`amd64`, Docker Engine 29.9.3, die **verschachtelte** Engine des Entwicklungs-Devcontainers dieses Repos (Speichertreiber `overlayfs`,
cgroup v2 mit `cgroupfs`, Kernel `6.8.0-147-generic`, „Debian GNU/Linux 13 (trixie) (containerized)“; Zeile `engine:` in `rohdaten-m1.txt`).
`u-boot devcontainer up` erkannte die Engine als lokal (alle `up` endeten mit Exit-Code 0). Das Binary `u-boot` (`u-boot version 0.23.0-dev`, der
Stand des Slice) ist mit `make compile BAU_COMPILE_OUT=<pfad>` gebaut. Das Hilfs-Image der Dev-Container-CLI
(`u-boot-devcontainer-cli:0.89.0-f1a94cce795e`) lag vor der Messung in der Engine; sein Bau gehört nicht zur Messung. Basis-Image der Projekte:
`mcr.microsoft.com/devcontainers/base:debian` (lokal), Netz offen. `GIT_TOKEN` ist ein Platzhalter; die Messung liest keine Zugangsdaten.

Alle Objekte der Messung tragen das Präfix `ubmess-rm` (Projekt- und Volume-Namen), die Container `u-boot-ubmess-rm-…` (so benennt `u-boot` sie:
`u-boot-<projekt>[-<instanz>]`), die Images `vsc-ubmess-rm-…`. Vor und nach jedem Skript steht die Inventur der Engine; die Differenz der
Container und Volumes ist in allen vier Rohdateien leer, die sieben vorhandenen fremden Container und zwei fremden Volumes blieben unberührt.
`05-aufraeumen.sh` entfernt die Images mit dem Präfix und liefert `rohdaten-aufraeumen.txt` (Reste mit dem Präfix: 0).

Läufe: M1 zweimal (`rohdaten-m1-erster-lauf.txt` und `rohdaten-m1.txt`; die Zeilen `AUSSAGE` sind gleich, `diff` leer); M2 bis M4 je ein Lauf in
der archivierten Fassung der Skripte. Davor lief je eine Vorfassung (bei M2 ohne den Abschnitt M2f) mit denselben Zeilen `AUSSAGE` (verglichen mit `diff`);
die Vorfassungen sind nicht archiviert.

## Befunde aus dem Code (gelesen, nicht gemessen)

- `DevcontainerStartService.Down` (`internal/hexagon/application/devcontainer_start.go`) sucht den Container mit `FindContainer` und ruft nur
  `launcher.Stop` (`docker stop`, `internal/adapter/driven/docker/devcontainer_launcher.go`); kein `docker rm`, kein `docker volume rm`.
- `Inspector.FindContainer` (`internal/adapter/driven/docker/inspector.go`) findet den Container über `docker ps -a` mit den Labels
  `devcontainer.local_folder` und, mit Instanz, `u-boot.instance`; der Container ohne Instanz ist der ohne dieses Label. Der Adapter ruft weder
  `docker inspect` auf einem Container noch liest er dessen `Mounts`.
- Die Volume-Namen entstehen aus `volumeMount` (`internal/hexagon/application/devcontainer_sandbox.go`) als `<projekt>-workspace-${devcontainerId}`,
  mit `nestedRuntime: podman` `<projekt>-containers-…`, mit `nestedRuntime: docker` `<projekt>-docker-…`, mit dem Feature `ollama`
  `<projekt>-ollama-…`; das Standardprofil erzeugt kein `mounts`. Die Auflösung von `${devcontainerId}` leistet die Dev-Container-CLI, nicht `u-boot`.

## Kurzfassung

| Frage | Befund (Messung) | Rohdaten |
|---|---|---|
| M1 | Alle Volumes eines Containers von `devcontainer up` stehen vollständig in `Mounts` von `docker inspect` (sieben Varianten, Menge gleich der Menge, die `up` in der Engine neu angelegt hat). Das Standardprofil hält keins. Die Volumes tragen keine Labels außer dem Marker anonymer Volumes; weder Labels noch Umgebung des Containers nennen die `devcontainerId`. | `rohdaten-m1*.txt` |
| M2 | `docker rm` eines laufenden Containers scheitert (Exit-Code 1), mit `-f` gelingt es; ein benanntes Volume bleibt stehen. `docker volume rm` eines Volumes, das ein weiterer Container hält (laufend oder gestoppt), scheitert (Exit-Code 1, auch mit `-f`); das Volume und sein Inhalt bleiben. Anonyme Volumes entfernt nur `docker rm -v`. | `rohdaten-m2.txt` |
| M3 | Jeder Container (mit und ohne Instanz) hat im Sandbox-Profil eigene Volumes; das Entfernen von Container und Volumes einer Instanz, bestimmt aus deren `Mounts`, lässt die anderen unberührt. Im Standardprofil teilen sich die Container den Bind-Mount des Projektordners. | `rohdaten-m3.txt` |
| M4 | `docker rm -f -v` entfernt **kein** benanntes Volume: Workspace-, Docker-Speicher- und Ollama-Volume bleiben, ein neuer Container hängt sie unter denselben Namen ein, der Inhalt ist da. Nach dem Entfernen des Containers unterscheidet nur noch der Name die Volumes. | `rohdaten-m4.txt` |

## M1 — Volumes des Containers und ihre Erkennung (`rohdaten-m1.txt`, `rohdaten-m1-erster-lauf.txt`)

Befehle je Variante (Projekt im Messordner, Sandbox-Profil mit `init --devcontainer --sandbox`, Einstellungen per `config set`, dann
`generate devcontainer`): `u-boot devcontainer up`; danach `docker inspect --format '{{range .Mounts}}…'` auf dem Container,
`docker volume ls -q` vor und nach `up`, `docker volume inspect` je Volume. Verglichen wird die Menge der Volume-Namen in `Mounts` mit der
Menge, die `up` in der Engine neu angelegt hat (`mounts-vollstaendig`).

| Variante | Konfiguration | Volumes in `Mounts` | Bind-Mounts | `mounts-vollstaendig` |
|---|---|---|---|---|
| `std` | Standardprofil | 0 | 1 (Projektordner nach `/workspaces/<projekt>`) | ja |
| `sbx` | Sandbox-Profil | 1 benannt: `<projekt>-workspace-<id>` | 0 | ja |
| `dock` | Sandbox, `nestedRuntime: docker` | 2 benannt: `-workspace-`, `-docker-` (Ziel `/var/lib/docker`) | 0 | ja |
| `pod` | Sandbox, `nestedRuntime: podman` | 2 benannt: `-workspace-`, `-containers-` (Ziel `/home/vscode/.local/share/containers`) | 0 | ja |
| `oll` | Sandbox, Feature `ollama` | 2 benannt: `-workspace-`, `-ollama-` (Ziel `/home/vscode/.ollama`) | 0 | ja |
| `full` | Sandbox, `nestedRuntime: docker` und `ollama` | 3 benannt: `-workspace-`, `-docker-`, `-ollama-` | 0 | ja |
| `anon` | Standardprofil, im Dockerfile des Projekts von Hand `VOLUME /ubmess-anon-daten` | 1 **anonym** (64 Hexzeichen, Ziel `/ubmess-anon-daten`) | 1 | ja |

- **Namen:** Bei den benannten Volumes folgt auf `<projekt>-<familie>-` ein Ende aus 52 Zeichen (Ziffern und Buchstaben); es ist die aufgelöste
  `${devcontainerId}`. Alle benannten Volumes **eines** Containers tragen dasselbe Ende (`verschiedene-ids-ueber-alle-benannten-volumes=1`). Dieselben
  Projektordner liefern in jedem Lauf dieselben Namen (zum Beispiel `…-sbx-workspace-15kf574…` in M1, M2 und M3).
- **Labels:** `docker volume inspect` zeigt für die benannten Volumes `labels=null options=null`; das anonyme Volume trägt das Label
  `com.docker.volume.anonymous`. Die Labels des Containers (`dev.containers.*`, `devcontainer.*`, `version`) enthalten das Namensende nicht, die Umgebung
  des Containers auch nicht (`id-in-labels=nein id-in-umgebung=nein` in allen Varianten mit benannten Volumes). `Config.Volumes` des Containers
  nennt nur das `VOLUME` des Images (`anon`); das Basis-Image deklariert keins.
- **Das Hilfsverfahren legt nichts nebenbei an:** In keiner Variante entstand ein Volume, das nicht in `Mounts` des Containers steht.
- **Runtimes im verschachtelten Daemon:** `dock` und `full`: Statusdatei `daemon` `result=ok`, `docker info` im Container antwortet
  (Server 29.9.3). `pod`: Statusdatei `vfs` `result=ok`, `podman version 6.1.4`, `podman info` Exit-Code 0; ein Container im Podman wurde nicht gestartet.
- **Dauer und Platz (Risiko des Slice):** `up` mit Bau des Images: im ersten Lauf (Images dieser Projekte neu, Basis-Image lokal, Schichten teils im
  Build-Cache) `std` 5 s, `sbx` 7 s, `dock` 10 s, `pod` 21 s, `oll` 46 s, `full` 48 s, `anon` 5 s; im zweiten Lauf (Build-Cache warm) 4 bis 8 s je
  Variante. Die Images (`docker images`, Größe wie angezeigt, je Projekt zwei Images, `vsc-…` und `vsc-…-uid`): 1,41 GB (`std`, `sbx`, `anon`),
  1,5 GB (`oll`), 1,54 GB (`pod`), 1,91 GB (`dock`), 2 GB (`full`); die Anzeige zählt gemeinsame Schichten je Image. Der Platzbedarf der Messung ist **nicht
  belastbar gemessen**: Das Entfernen aller Images gab nur 139.264 Byte frei, weil der Build-Cache der Engine die Schichten hält
  (`rohdaten-aufraeumen.txt`). Die Volumes sind klein (Workspace 6 B, Docker-Speicher 278,6 kB, Ollama 2,5 kB mit Schlüsselpaar, ohne Modell; `docker system df -v` in M4).

### Eingang für den CR (M1)

**Darf tragen:**

- Die Volumes eines Containers, den `u-boot devcontainer up` aus der erzeugten Konfiguration startet, lassen sich über `Mounts` des Containers
  (Typ `volume`, Feld `Name`) vollständig ermitteln; die Menge deckte sich in jeder gemessenen Variante mit der Menge, die `up` in der Engine
  angelegt hat. Eine Auflösung von `${devcontainerId}` durch `u-boot` ist dafür nicht nötig.
- Im Standardprofil hält der Container kein Volume, nur den Bind-Mount des Projektordners.
- Im Sandbox-Profil hält er das Workspace-Volume, mit `nestedRuntime: docker` das Volume des Docker-Speichers, mit `nestedRuntime: podman` das
  Volume der Container, mit dem Feature `ollama` das Volume der Modelle (Namen wie im Code).
- Ein Volume, das das Image mit `VOLUME` deklariert, ist anonym und steht ebenfalls in `Mounts`.
- Weder die Labels des Containers noch die der benannten Volumes tragen die `devcontainerId` oder einen Verweis zwischen Container und Volume.

**Darf nicht tragen:**

- Dass `u-boot` die Volumes eines schon entfernten Containers ermitteln kann: Ohne Container fehlen `Mounts` (M4); nur der Name bleibt.
- Dass sich die Volume-Namen aus `${devcontainerId}` berechnen lassen: Das Verfahren der CLI ist weder gelesen noch gemessen.
- Dass u-boot-Container keine anonymen Volumes haben: Gemessen sind die erzeugten Konfigurationen (keines); ein von Hand ergänztes `VOLUME` erzeugt eins
  (`anon`). Weitere Features des Katalogs außer `ollama`, Features aus externen Quellen, von Hand geschriebene `mounts` und `docker-cli` sind nicht gemessen.
- Dass die Runtime im Container „funktioniert“: Gemessen sind die Statusdateien und `docker info`/`podman info`, kein Lauf von Containern darin.
- Jede Aussage für andere Hosts und Engines (Docker Desktop, Colima, rootless, Podman als Host-Engine, entfernt, `arm64`, Windows).

## M2 — `docker rm` und `docker volume rm` (`rohdaten-m2.txt`)

Projekt `sbx` (Sandbox-Profil, ein Workspace-Volume); als zweiter Halter ein Container `alpine:3`; Projekt `anon` (Standardprofil mit `VOLUME`).

**a — laufender Container.** `docker rm <id>`: Exit-Code 1, `Error response from daemon: cannot remove container "…": container is running: stop the
container before removing or force remove`; der Container läuft weiter. `docker rm -f <id>`: Exit-Code 0, Container weg; das Workspace-Volume ist da,
eine zuvor geschriebene Markerdatei liest ein Fremd-Container aus dem Volume.

**d — `u-boot` danach.** `u-boot devcontainer down` ohne Container: Exit-Code 0, „No devcontainer container found for this project; nothing to
stop“. `u-boot devcontainer up` legt einen neuen Container an (andere ID), hängt dasselbe Volume ein (gleicher Name), die Markerdatei ist im neuen
Container da. Nach `docker volume rm` des Volumes legt `up` ein Volume **gleichen Namens** neu an; die Markerdatei fehlt.

**b — Volume, das ein zweiter Container hält.** Der zweite Container (`alpine:3`, läuft) hängt das Workspace-Volume ein und sieht die Markerdatei. Nach
`docker rm -f -v` des Devcontainers: `docker volume rm <volume>` Exit-Code 1, `volume is in use - [<id des zweiten Containers>]`; `docker volume rm -f
<volume>` ebenso (Exit-Code 1, derselbe Text). Volume und Inhalt bleiben. Nach dem Entfernen des zweiten Containers gelingt `docker volume rm` (Exit-Code 0).

**c — gestoppter Container** (nach `u-boot devcontainer down`, Status `exited`). `docker inspect` des gestoppten Containers nennt dieselben Volumes in
`Mounts` wie beim laufenden. `docker volume rm` und `docker volume rm -f`: Exit-Code 1, `volume is in use - [<id des gestoppten Containers>]`.
`docker rm <id>` (ohne `-f`): Exit-Code 0; das Volume bleibt; danach gelingt `docker volume rm`. `u-boot devcontainer down` ohne Container: Exit-Code 0.

**f — `u-boot devcontainer up --recreate` mit benanntem Volume** (Projekt `sbx`): nach `--recreate` ist der Container ein neuer (andere ID), der alte ist
nicht mehr da (`docker inspect` scheitert), das Workspace-Volume hat denselben Namen und ist vorhanden, die zuvor geschriebene Markerdatei ist im neuen
Container lesbar.

**e — anonyme Volumes.**

| Fall | anonymes Volume danach | benanntes Volume danach |
|---|---|---|
| Probe-Container (`docker run -v /ubmess-anon -v <benannt>:/…`), `docker rm -f` | **vorhanden** (`dangling=true`) | vorhanden |
| derselbe Aufbau, `docker rm -f -v` | entfernt | vorhanden |
| Projekt `anon`, `u-boot devcontainer up --recreate` | des alten Containers **vorhanden**; der neue Container hat ein neues anonymes Volume | – |
| Projekt `anon`, `docker rm -f` | vorhanden | – |
| Projekt `anon`, `docker rm -f -v` | entfernt | – |

### Eingang für den CR (M2)

**Darf tragen:**

- `docker rm` entfernt einen gestoppten Container; bei einem laufenden Container scheitert es (Exit-Code 1), solange nicht `-f` gesetzt ist.
- Das Entfernen eines Containers löscht kein benanntes Volume, mit oder ohne `-v`; Name und Inhalt bleiben, und ein folgendes `up` hängt dasselbe Volume
  wieder ein. Das gilt auch für `u-boot devcontainer up --recreate` (Workspace-Volume im Sandbox-Profil).
- `docker volume rm` scheitert (Exit-Code 1, „volume is in use“ mit der Container-ID), solange irgendein Container, laufend oder gestoppt, das Volume
  hält; `-f` ändert das nicht. Ein Volume ist also erst nach dem Entfernen des **letzten** Halters entfernbar; ein von einem zweiten Container gehaltenes
  Volume bleibt mit seinem Inhalt stehen.
- Anonyme Volumes entfernt `docker rm -v`; ohne `-v` bleiben sie als verwaiste Volumes zurück, auch nach `u-boot devcontainer up --recreate`.
- `u-boot devcontainer down` endet ohne Container mit Exit-Code 0.

**Darf nicht tragen:**

- Dass `docker volume rm -f` ein von einem Container gehaltenes Volume entfernt (es tut es nicht).
- Dass ein gestoppter Container sein Volume freigibt: Er hält es bis zum Entfernen des Containers.
- Dass `docker rm -v` oder `-f` benannte Volumes entfernt.
- Aussagen zu `docker volume prune`, zu Containern im Zustand `paused`, `restarting` oder `created`, zu zwei gleichzeitigen Aufrufen und zu anderen
  Docker-Versionen: nicht gemessen. Der zweite Halter war ein Fremd-Container (`alpine:3`); zwei Devcontainer auf **einem** Volume sind nicht gemessen (M3: sie
  haben getrennte Volumes).

## M3 — Instanzen (`rohdaten-m3.txt`)

**Projekt `sbx` (Sandbox-Profil), drei Container:** ohne Instanz (`u-boot-ubmess-rm-sbx`), `--instance a`, `--instance b`. Jeder hat ein eigenes
Workspace-Volume: die drei Namen unterscheiden sich im 52-stelligen Ende (`workspace-volumes-verschieden=ja`). Eine Markerdatei je Container
liest nur dieser Container (`marker-isoliert=ja`). `docker ps -a --filter label=devcontainer.local_folder=<ordner>` liefert alle drei; der Filter
`label=u-boot.instance=a` genau einen; beim Container ohne Instanz ist der Wert des Labels leer. `docker volume ls --filter name=<projekt>-workspace-`
liefert alle drei Volumes; der Name sagt nicht, zu welchem Container ein Volume gehört. **Entfernen von Instanz a** (`docker rm -f` des Containers,
dann `docker volume rm` der Volumes aus `Mounts` **von a**): a und sein Volume sind weg; der Container ohne Instanz und `b` laufen weiter, ihre
Volumes sind da, die Markerdateien lesbar. `u-boot devcontainer down --instance a` danach: Exit-Code 0 („nothing to stop“); `down --instance b` stoppt
nur `b`.

**Projekt `full` (`nestedRuntime: docker` und `ollama`), Instanzen a und b:** je Instanz drei benannte Volumes (`-workspace-`, `-ollama-`, `-docker-`), keine
gemeinsamen (`gemeinsame=0`). Das Entfernen von a und der drei Volumes aus deren `Mounts` lässt die drei Volumes von b stehen und b laufen. Der Inhalt
der Volumes des Docker-Speichers und der Modelle ist hier nicht gegengeprüft (nur die Namen und das Entfernen).

**Projekt `std` (Standardprofil):** der Container ohne Instanz und `--instance x` haben 0 Volumes und denselben Bind-Mount (`bind-quelle-gleich=ja`);
eine Datei, die `x` in den Workspace schreibt, steht im Projektordner des Hosts und im Container ohne Instanz; nach `docker rm -f -v` von `x` bleibt
sie im Projektordner, der Container ohne Instanz läuft weiter.

### Eingang für den CR (M3)

**Darf tragen:**

- Im Sandbox-Profil hat jeder Container, auch jede Instanz, eigene Volumes: gleiche Namensfamilie, verschiedenes 52-stelliges Ende; die Workspace-Volumes
  zweier Instanzen und des Containers ohne Instanz sind gegeneinander isoliert (gemessen mit drei Containern).
- Wer Container und Volumes einer Instanz entfernt und die Volumes aus `Mounts` **dieses Containers** bestimmt, trifft die Volumes anderer Instanzen und des
  Containers ohne Instanz nicht (gemessen: zwei Instanzen mit Docker-Speicher und Modellen, drei Container ohne diese).
- Eine Instanz erkennt `u-boot` am Label `u-boot.instance` des Containers; der Container ohne Instanz trägt das Label nicht.
- Im Standardprofil hat auch eine Instanz keine Volumes; das Entfernen ihres Containers lässt den Projektordner und seinen Inhalt unberührt.

**Darf nicht tragen:**

- Dass sich die Volumes einer Instanz über das Namenspräfix `<projekt>-…-` bestimmen lassen: Das Präfix ist allen Instanzen gemeinsam. Ob ein Präfix eines
  Projekts das eines anderen Projekts trifft (zum Beispiel `app` und `app-workspace`), ist nicht gemessen.
- Isolation des Docker-Speichers oder der Modelle zwischen Instanzen (Namen verschieden, Inhalt nicht geprüft; [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt sie auch nicht zu).
- Drei oder mehr gleichzeitige Instanzen im selben Projekt (gemessen: höchstens drei Container, davon zwei Instanzen).
- Instanzen mit `nestedRuntime: podman` (nicht gemessen, nur ohne Instanz in M1).

## M4 — Was das Entfernen des Containers übersteht (`rohdaten-m4.txt`)

Projekt `full` (Sandbox, `nestedRuntime: docker`, `ollama`), ein Container ohne Instanz. In jedes der drei Volumes steht eine Markerdatei; im Ollama-Volume
legt `ollama serve` (Version 0.40.6) beim Start das Schlüsselpaar `id_ed25519` und `id_ed25519.pub` sowie die Verzeichnisse `models` und `cache` an. Ein Modell
wurde **nicht** geladen; die Markerdatei steht stellvertretend.

- `docker rm -f -v <id>`: Exit-Code 0. Workspace-, Docker-Speicher- und Ollama-Volume sind alle **vorhanden** (`docker volume ls` zeigt sie als
  `dangling=true`). Die Engine kennt zu den verwaisten Volumes nur Name, Treiber `local` und Anlagezeit (`labels=null options=null`); `docker system df -v`
  zeigt `links=0`.
- `u-boot devcontainer up`: neuer Container (andere ID), dieselben drei Volume-Namen; die drei Markerdateien sind lesbar, und der öffentliche Schlüssel
  des Ollama-Volumes hat dieselbe Prüfsumme wie vorher.
- **Unterscheidung nach dem Entfernen:** Allein über den Namen. Der Infix `-ollama-` trifft genau ein Volume; die drei Volumes tragen dasselbe
  52-stellige Ende (ein Container). Vor dem Entfernen unterscheidet das Ziel in `Mounts` (`/home/vscode/.ollama`).
- **Auswahl über den Namen:** Wurden alle Volumes des Projekts außer dem mit dem Infix `-ollama-` entfernt (`docker volume rm`, Exit-Code 0), blieb das
  Ollama-Volume stehen; ein neuer Container hängte es ein und las die Markerdatei, das Workspace-Volume war neu und leer.

### Eingang für den CR (M4)

**Darf tragen:**

- Das Entfernen des Containers (auch `docker rm -f -v`) lässt **alle** benannten Volumes stehen, nicht nur das der Ollama-Modelle; ein folgendes `u-boot devcontainer up`
  hängt sie unter denselben Namen wieder ein, der Inhalt ist da. `up --recreate` ist für das Workspace-Volume in M2 gemessen, nicht für die Volumes des Docker-Speichers
  und der Modelle.
- Das Volume der Ollama-Modelle enthält neben den Modellen das Schlüsselpaar des Servers (bestätigt, was die Spezifikation sagt) und überlebt das Entfernen
  des Containers.
- Nach dem Entfernen des Containers hat ein Volume weder Label noch Verweis auf seine Herkunft; die Unterscheidung (Workspace, Docker-Speicher, Modelle) trägt der
  Name (`-workspace-`, `-docker-`, `-ollama-`), die Zuordnung zu einem Container das gemeinsame 52-stellige Ende.
- Verwaiste Volumes erscheinen in `docker volume ls --filter dangling=true`.

**Darf nicht tragen:**

- Dass ein Modell überlebt: Es wurde keins geladen; gemessen sind Markerdatei und Schlüsselpaar.
- Dass die Modelle den **Neubau des Images** überstehen: Alle `up` nutzten dasselbe Image; das Lastenheft sagt es ebenfalls nicht zu.
- Dass die Volume-Namen stabil bleiben, wenn der Projektordner umzieht oder umbenannt wird: Gemessen sind gleiche Namen bei gleichem Pfad.
- Dass das Volume des Podman-Speichers sich so verhält (nicht entfernt und neu eingehängt gemessen).
- Die Entscheidung, welche Volumes ein späteres `--volumes` entfernt oder stehen lässt (Sache des CR).

## Aufruf

```bash
# Binary des Slice-Stands (im Devcontainer-Container):
make compile BAU_COMPILE_OUT=<pfad>/u-boot
UBOOT=<pfad>/u-boot WORK=<messordner-ausserhalb-des-repos> tools/messung/slice-v1-devcontainer-remove-messung/00-alles.sh
UBOOT=<pfad>/u-boot WORK=<messordner> tools/messung/slice-v1-devcontainer-remove-messung/05-aufraeumen.sh
```

Je Messung einzeln: `01-m1-volumes.sh`, `02-m2-entfernen.sh`, `03-m3-instanz.sh`, `04-m4-ueberleben.sh` (je eine Rohdatei `rohdaten-m1.txt` bis
`rohdaten-m4.txt`). Lesehilfe: Zeilen `AUSSAGE …` sind die stabilen Aussagen, die ein zweiter Lauf wiederholen muss; Dauern, Container-IDs und
Zeitstempel streuen. Die Datei `rohdaten-m1-erster-lauf.txt` ist die Ausgabe des ersten Laufs von M1 (Dauern mit neu gebauten Images).

## Grenzen

- Ein Host, eine Engine (die verschachtelte des Devcontainers), Docker 29.9.3; ein Lauf je Messpunkt außer M1. Eine native Engine auf einem
  Linux-Host ist nicht gemessen; wie weit der Aufbau abweicht, ist unbekannt. `u-boot` hat die Engine als lokal behandelt.
- Gemessen sind die von `u-boot` erzeugten Konfigurationen (Standard- und Sandbox-Profil) mit den Optionen `nestedRuntime` und dem Feature `ollama`; die übrigen
  Features des Katalogs, Features aus externen Quellen und von Hand geschriebene Konfigurationen sind nicht gemessen (Ausnahme: ein von Hand angehängtes
  `VOLUME` im Dockerfile).
- Der Platzbedarf der Images ist nicht belastbar gemessen (Build-Cache); die Bauzeiten stehen mit der Einschränkung ihres Cache-Zustands in M1.
- Ein Modell von Ollama wurde nicht geladen; das Netz diente dem Bau der Images.
- Ob `docker volume prune`, `docker system prune` oder `docker compose down -v` anders wirken, ist nicht gemessen.
