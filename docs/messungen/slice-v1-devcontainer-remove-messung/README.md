# Messung: Volumes eines Devcontainers, das Verhalten von `docker rm` und `docker volume rm` und `down` ohne `--instance` (Grundlage für `devcontainer remove`)

| Dokument    | Messprotokoll |
| ----------- | ------------- |
| Auslösender Vorgang | [`slice-v1-devcontainer-remove-messung`](../../plan/planning/in-progress/slice-v1-devcontainer-remove-messung.md) (erster von vier Slices zu `devcontainer remove`; Grundlage: [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Volume-Namen je Instanz: [`ADR-0015`](../../plan/adr/0015-sandbox-volumes-pro-instanz.md)) |
| Gelesen am  | 2026-11-05 |
| Skripte     | `tools/messung/slice-v1-devcontainer-remove-messung/` (`01-m1-volumes.sh`, `02-m2-entfernen.sh`, `03-m3-instanz.sh`, `04-m4-ueberleben.sh`, `04b-m5-down-ohne-instanz.sh`, `05-aufraeumen.sh`, gemeinsam `lib.sh`; `00-alles.sh` ruft 01 bis 04b) | <!-- d-check:ignore (Skriptpfad, keine Slice-Kennung) -->

Die Messung fährt den ausgelieferten Befehl `u-boot devcontainer up` und `down` (kein Nachbau des Aufrufs der Dev-Container-CLI) und die
Docker-Befehle, mit denen ein Nutzer heute aufräumt. Sie ändert keinen Produktionscode.

## Umgebung

Linux-`amd64`, Docker Engine 29.9.3, die **verschachtelte** Engine des Entwicklungs-Devcontainers dieses Repos (Speichertreiber `overlayfs`,
cgroup v2 mit `cgroupfs`, Kernel `6.8.0-147-generic`, „Debian GNU/Linux 13 (trixie) (containerized)“; Zeile `engine:` in `rohdaten-m1.txt`).
`u-boot devcontainer up` erkannte die Engine als lokal (alle `up` endeten mit Exit-Code 0). Das Binary `u-boot` (`u-boot version 0.23.0-dev`, der
Stand des Slice) ist mit `make compile BAU_COMPILE_OUT=<pfad>` gebaut. Das Hilfs-Image der Dev-Container-CLI
(`u-boot-devcontainer-cli:0.89.0-f1a94cce795e`) lag vor der Messung in der Engine; sein Bau gehört nicht zur Messung. Basis-Image der Projekte:
`mcr.microsoft.com/devcontainers/base:debian` (lokal), Netz offen. `GIT_TOKEN` ist ein Platzhalter; die Messung liest keine Zugangsdaten.

Alle Objekte der Messung tragen das Präfix `ubmess-rm-` (Projekt- und Volume-Namen), die Container `u-boot-ubmess-rm-…` (so benennt `u-boot` sie:
`u-boot-<projekt>[-<instanz>]`), die Images `vsc-ubmess-rm-…`. Jede Rohdatei der archivierten Fassung (`rohdaten-m1.txt` bis `rohdaten-m5.txt`) trägt die Inventur der
Engine vor und nach dem Lauf (Zeilen `inventur vorher` und `inventur nachher` mit Zahl und Namen der Container und Volumes) und die Differenz (`--- Differenz ps`, `--- Differenz vol`). In
allen fünf Dateien lauten beide Inventuren gleich: **7 Container und 2 Volumes ohne das Präfix** (Container `boring_noyce`, `nice_fermi`, `pedantic_spence`,
`relaxed_archimedes`, `stoic_aryabhata`, `stupefied_elion`, `unruffled_perlman`; Volumes `dcsdmn72acy22yw-workspace-…` und `dcsdmn7vfi1bhvb-workspace-…`),
die Differenz ist leer, „Reste mit dem Präfix: 0“. Das sind Objekte, die nicht zur Messung gehören; sie blieben unberührt. Die Skripte löschen nur Container, Volumes und Images, deren
Name mit dem Präfix und einem Bindestrich beginnt (und die Volumes, die das Register aus den `Mounts` solcher Container kennt). `05-aufraeumen.sh` entfernt die Images mit dem Präfix und
liefert `rohdaten-aufraeumen.txt` (Reste mit dem Präfix: 0). Nicht entfernt werden das Basis-Image, das Hilfs-Image, der Build-Cache, fremde Objekte und das Probe-Image
`alpine:3` (M2: Zweitcontainer, Lesen der Markerdatei, Probe-Container); ob `alpine:3` vor der ersten Messung in der Engine lag, ist nicht festgehalten, und `05-aufraeumen.sh` entfernt es nicht,
weil ein Image dieses Namens fremd sein kann.

**Läufe und ihre Herkunft.** Die Rohdaten `rohdaten-m1.txt` bis `rohdaten-m5.txt` und `rohdaten-aufraeumen.txt` stammen aus **einem** Lauf der archivierten Fassung der Skripte
(`00-alles.sh`, danach `05-aufraeumen.sh`; Stempel in den Dateien, 2026-11-05 13:25 bis 13:29 UTC, Build-Cache warm). Gegen die im Commit `8861376d` archivierten Rohdaten von M1 bis M4
(Lauf der in `8861376d` archivierten Skripte, die weder Inventurzeilen noch den Label-Schlüssel in M3 schrieben) sind die Zeilen `AUSSAGE` verglichen (52-stellige Namensenden gleichgesetzt, `diff`):
gleich, bis auf die Zeile zum Label-Schlüssel in M3, die es vorher nicht gab. Ältere Läufe (die Vorfassungen von M2 bis M4 vor dem ersten Commit) sind nicht archiviert; ein Vergleich mit ihnen
ist nicht prüfbar, und das Protokoll stützt sich auf keine ihrer Zeilen. `rohdaten-m1-erster-lauf.txt` ist die Ausgabe des ersten Laufs von M1 (2026-11-05 12:50 UTC, Dauern mit neu gebauten Images)
und stammt aus einer Skriptfassung, die **nicht** archiviert ist (sie schrieb die Kopfzeilen `engine:` und `hilfs-image von …` noch nicht); ihre Zeilen `AUSSAGE` sind mit denen von `rohdaten-m1.txt`
verglichen und gleich (`diff` leer, gleiche Gleichsetzung der Namensenden). Diese Datei trägt die Dauern des ersten Laufs (siehe M1) und sonst nichts, was das Protokoll zählt.

## Befunde aus dem Code (gelesen, nicht gemessen)

- `DevcontainerStartService.Down` (`internal/hexagon/application/devcontainer_start.go`) sucht den Container mit `FindContainer` und ruft nur
  `launcher.Stop` (`docker stop`, `internal/adapter/driven/docker/devcontainer_launcher.go`); kein `docker rm`, kein `docker volume rm`. Die Zweige von `Down`:
  `driven.ErrOnlyInstanceContainers` endet mit `onlyInstancesError("down")` (Exit-Code `10`, [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten));
  `ErrContainerEngineNotLocal` und ein Engine-Fehler enden mit einem Fehler; ohne Fund ist die Antwort leer (Exit-Code `0`, Meldung „nothing to stop“); ein gestoppter Container gibt
  seine ID ohne `Stopped` zurück; ein laufender wird mit `launcher.Stop` gestoppt. `up` und `shell` ohne `--instance` fragen denselben Suchweg vorab
  (`containerBeforeStart`) und enden bei `ErrOnlyInstanceContainers` mit demselben Fehler.
- `Inspector.FindContainer` (`internal/adapter/driven/docker/inspector.go`) findet den Container über `docker ps -a` mit dem Label
  `devcontainer.local_folder` und, mit Instanz, `u-boot.instance`. Ohne Instanz filtert `docker ps` nur nach dem Ordner-Label; der Adapter liest den Wert des Labels
  `u-boot.instance` je Zeile, der erste Container ohne Wert ist der Fund, und gibt es Treffer, aber keinen ohne Wert, ist das Ergebnis `ErrOnlyInstanceContainers`. Der Adapter ruft weder
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
| M5 | `down` ohne `--instance` stoppt nur den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen (laufend oder gestoppt), endet es mit Exit-Code `10` und stoppt nichts; ohne jeden Container endet es mit Exit-Code `0`. `up` ohne `--instance` endet bei nur Instanz-Containern ebenfalls mit Exit-Code `10`. | `rohdaten-m5.txt` |

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
  `${devcontainerId}`. Alle benannten Volumes **eines** Containers tragen dasselbe Ende (`verschiedene-ids-ueber-alle-benannten-volumes=1`). Im archivierten
  Lauf liefert derselbe Projektordner dieselben Namen (`…-sbx-workspace-1gju3f8…` in M1, M2 und M3, dort der Container ohne Instanz); `rohdaten-m1-erster-lauf.txt` trägt für `sbx` andere Namen
  (`…-sbx-workspace-15kf574…`), der Messordner dieses Laufs ist nicht festgehalten.
- **Labels:** `docker volume inspect` zeigt für die benannten Volumes `labels=null options=null`; das anonyme Volume trägt das Label
  `com.docker.volume.anonymous`. Die Labels des Containers (`dev.containers.*`, `devcontainer.*`, `version`) enthalten das Namensende nicht, die Umgebung
  des Containers auch nicht (`id-in-labels=nein id-in-umgebung=nein` in allen Varianten mit benannten Volumes). `Config.Volumes` des Containers
  nennt nur das `VOLUME` des Images (`anon`); das Basis-Image deklariert keins.
- **Das Hilfsverfahren legt nichts nebenbei an:** In keiner Variante entstand ein Volume, das nicht in `Mounts` des Containers steht.
- **Runtimes im verschachtelten Daemon:** `dock` und `full`: Statusdatei `daemon` `result=ok`, `docker info` im Container antwortet
  (Server 29.9.3). `pod`: Statusdatei `vfs` `result=ok`, `podman version 6.1.4`, `podman info` Exit-Code 0; ein Container im Podman wurde nicht gestartet.
- **Dauer und Platz (Risiko des Slice):** `up` mit Bau des Images: im ersten Lauf (`rohdaten-m1-erster-lauf.txt`; Images dieser Projekte neu, Basis-Image lokal, Schichten teils im
  Build-Cache) `std` 5 s, `sbx` 7 s, `dock` 10 s, `pod` 21 s, `oll` 46 s, `full` 48 s, `anon` 5 s; im Lauf der archivierten Fassung (Build-Cache warm) 4 bis 8 s je
  Variante (`rohdaten-m1.txt`). Die Images (`docker images`, Größe wie angezeigt, je Projekt zwei Images, `vsc-…` und `vsc-…-uid`): 1,41 GB (`std`, `sbx`, `anon`),
  1,5 GB (`oll`), 1,54 GB (`pod`), 1,91 GB (`dock`), 2 GB (`full`); die Anzeige zählt gemeinsame Schichten je Image. Der Platzbedarf der Messung ist **nicht
  belastbar gemessen**: Das Entfernen aller 14 Images änderte den freien Platz der Docker-Wurzel um −1.937.408 Byte (er sank; andere Last auf dem Host), weil der Build-Cache der Engine die Schichten hält
  (`rohdaten-aufraeumen.txt`). Die Volumes sind klein (Workspace 6 B, Docker-Speicher 278,6 kB, Ollama 2,5 kB mit Schlüsselpaar, ohne Modell; `docker system df -v` in M4).

### Eingang für den CR (M1)

**Darf tragen:**

- Die Volumes eines Containers, den `u-boot devcontainer up` aus der erzeugten Konfiguration startet, lassen sich über `Mounts` des Containers
  (Typ `volume`, Feld `Name`) vollständig ermitteln; die Menge deckte sich in jeder gemessenen Variante mit der Menge, die `up` in der Engine
  angelegt hat. Eine Auflösung von `${devcontainerId}` durch `u-boot` ist dafür nicht nötig.
- Im Standardprofil ohne Katalog-Feature hält der Container kein Volume, nur den Bind-Mount des Projektordners (gemessen: `std` und `anon`). Katalog-Features im Standardprofil sind nicht gemessen (unten).
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

**d — `u-boot` danach.** `u-boot devcontainer down` in einem Projekt ohne Container: Exit-Code 0, „No devcontainer container found for this project; nothing to
stop“ (Projekt ohne Instanzen; das Verhalten bei nur Instanz-Containern ist M5). `u-boot devcontainer up` legt einen neuen Container an (andere ID), hängt dasselbe Volume ein (gleicher Name), die Markerdatei ist im neuen
Container da. Nach `docker volume rm` des Volumes legt `up` ein Volume **gleichen Namens** neu an; die Markerdatei fehlt.

**b — Volume, das ein zweiter Container hält.** Der zweite Container (`alpine:3`, läuft) hängt das Workspace-Volume ein und sieht die Markerdatei. Nach
`docker rm -f -v` des Devcontainers: `docker volume rm <volume>` Exit-Code 1, `volume is in use - [<id des zweiten Containers>]`; `docker volume rm -f
<volume>` ebenso (Exit-Code 1, derselbe Text). Volume und Inhalt bleiben. Nach dem Entfernen des zweiten Containers gelingt `docker volume rm` (Exit-Code 0).

**c — gestoppter Container** (nach `u-boot devcontainer down`, Status `exited`). `docker inspect` des gestoppten Containers nennt dieselben Volumes in
`Mounts` wie beim laufenden. `docker volume rm` und `docker volume rm -f`: Exit-Code 1, `volume is in use - [<id des gestoppten Containers>]`.
`docker rm <id>` (ohne `-f`): Exit-Code 0; das Volume bleibt; danach gelingt `docker volume rm`. `u-boot devcontainer down` ohne Container im Projekt: Exit-Code 0.

**f — `u-boot devcontainer up --recreate` mit benanntem Volume** (Projekt `sbx`): nach `--recreate` ist der Container ein neuer (andere ID), der alte ist
nicht mehr da (`docker inspect` scheitert), das Workspace-Volume hat denselben Namen und ist vorhanden, die zuvor geschriebene Markerdatei ist im neuen
Container lesbar. Denselben Zusammenhang pinnt der Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (`internal/e2e/devcontainer_instance_docker_test.go`, Abschnitt
`--recreate` einer Instanz: Workspace-Volume unverändert, Datei bleibt) für Instanzen; M2f misst es am Container ohne Instanz.

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
- `u-boot devcontainer down` endet in einem Projekt ohne jeden Container mit Exit-Code 0 (mit und ohne `--instance`, M5). Gibt es im Projekt nur Container von Instanzen,
  endet `down` ohne `--instance` mit Exit-Code 10 (M5), nicht mit 0.

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
`label=u-boot.instance=a` genau einen; die Vorlage `{{.Label "u-boot.instance"}}` von `docker ps` liefert beim Container ohne Instanz die leere Zeichenkette. Die Label-Map des Containers
(`docker inspect`, `Config.Labels`) trägt den Schlüssel `u-boot.instance` beim Container ohne Instanz **nicht**, bei `a` und `b` mit dem Namen der Instanz als Wert
(`m3-sbx-label-u-boot.instance-schluessel: ohne-instanz=fehlt a=vorhanden b=vorhanden`). `docker volume ls --filter name=<projekt>-workspace-`
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
- Eine Instanz erkennt `u-boot` am Label `u-boot.instance` des Containers; der Container ohne Instanz trägt den Schlüssel in `Config.Labels` nicht (gemessen mit `docker inspect`;
  `docker ps` zeigt für ihn die leere Zeichenkette).
- Im Standardprofil hat auch eine Instanz keine Volumes; das Entfernen ihres Containers lässt den Projektordner und seinen Inhalt unberührt.

**Darf nicht tragen:**

- Dass sich die Volumes einer Instanz über das Namenspräfix `<projekt>-…-` bestimmen lassen: Das Präfix ist allen Instanzen gemeinsam. Ob ein Präfix eines
  Projekts das eines anderen Projekts trifft (zum Beispiel `app` und `app-workspace`), ist nicht gemessen.
- Isolation des Docker-Speichers oder der Modelle zwischen Instanzen (Namen verschieden, Inhalt nicht geprüft; [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt sie auch nicht zu).
- Drei oder mehr gleichzeitige Instanzen im selben Projekt (gemessen: höchstens drei Container, davon zwei Instanzen).
- Instanzen mit `nestedRuntime: podman` (nicht gemessen, nur ohne Instanz in M1).

## M4 — Was das Entfernen des Containers übersteht (`rohdaten-m4.txt`)

Projekt `full` (Sandbox, `nestedRuntime: docker`, `ollama`), ein Container ohne Instanz. In jedes der drei Volumes steht eine Markerdatei; im Ollama-Volume
legt `ollama serve` beim Start das Schlüsselpaar `id_ed25519` und `id_ed25519.pub` sowie die Verzeichnisse `models` und `cache` an. Ein Modell
wurde **nicht** geladen; die Markerdatei steht stellvertretend. Die Version `0.40.6` stammt aus `ollama --version`, einer Warnung des Clients (`Warning: client version is 0.40.6`; der Server war zu dem Zeitpunkt nicht gestartet); die Version des Servers steht in den Rohdaten nicht.

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
- Das Volume des Ollama-Verzeichnisses (`-ollama-`, Ziel `/home/vscode/.ollama`) überlebt das Entfernen des Containers (auch `docker rm -f -v`); ein neuer Container
  hängt es unter demselben Namen ein, und der öffentliche Schlüssel hat dieselbe Prüfsumme wie vorher. Gemessen ist der Inhalt: Markerdatei, Schlüsselpaar
  `id_ed25519`/`id_ed25519.pub` und die leeren Verzeichnisse `models` und `cache`.
- Nach dem Entfernen des Containers hat ein Volume weder Label noch Verweis auf seine Herkunft; die Unterscheidung (Workspace, Docker-Speicher, Modelle) trägt der
  Name (`-workspace-`, `-docker-`, `-ollama-`), die Zuordnung zu einem Container das gemeinsame 52-stellige Ende.
- Verwaiste Volumes erscheinen in `docker volume ls --filter dangling=true`.

**Darf nicht tragen:**

- Dass ein Modell das Entfernen des Containers übersteht und dass ein neuer Container mit demselben Volume das Modell listet: Es wurde keins geladen; gemessen sind Markerdatei und Schlüsselpaar.
- Dass die Modelle den **Neubau des Images** überstehen: Alle `up` nutzten dasselbe Image; das Lastenheft sagt es ebenfalls nicht zu.
- Dass die Volume-Namen stabil bleiben, wenn der Projektordner umzieht oder umbenannt wird: Gemessen sind gleiche Namen bei gleichem Pfad.
- Dass das Volume des Podman-Speichers sich so verhält (nicht entfernt und neu eingehängt gemessen).
- Die Entscheidung, welche Volumes ein späteres `--volumes` entfernt oder stehen lässt (Sache des CR).

## M5 — `down` ohne `--instance`: welcher Container, und was bei nur Instanz-Containern (`rohdaten-m5.txt`)

Ergänzung der Fix-Runde nach dem Review (F-7); die Messfrage steht im Slice unter M5. Projekt `sbx` (Sandbox-Profil), Befehle `u-boot devcontainer down` und `up` aus dem
ausgelieferten Binary; Zeilen `AUSSAGE m5-…`. Der Suchweg im Code ist derselbe wie hinter `up` und `shell` (siehe „Befunde aus dem Code“).

| Fall | Aufbau | Befehl | Ergebnis |
|---|---|---|---|
| A | Projekt ohne Container | `down`; `down --instance a` | Exit-Code 0, „No devcontainer container found for this project; nothing to stop“ (beide) |
| B | Container ohne Instanz, Instanz `a`, Instanz `b`, alle laufend | `down` | Exit-Code 0, „Devcontainer stopped: container 98aa6e5564b1 (not removed)“; gestoppt ist nur der Container ohne Instanz, `a` und `b` laufen weiter |
| C | nur `a` und `b` (der Container ohne Instanz mit `docker rm -f -v` entfernt), beide laufend | `down` | **Exit-Code 10**, „the project has only containers of instances; pass --instance <name> to act on one of them“; `a` und `b` laufen weiter |
| C | derselbe Aufbau | `up` | **Exit-Code 10**, derselbe Hinweis; kein Container ohne Instanz angelegt; die Warnung zur Einhängung des Engine-Sockets in den Hilfs-Container steht davor in der Ausgabe |
| D | `a` gestoppt (`down --instance a`, Exit-Code 0), `b` läuft | `down` | **Exit-Code 10**, derselbe Hinweis; `b` läuft weiter |
| E | `a` und `b` gestoppt | `down` | **Exit-Code 10**, derselbe Hinweis |

Die Fälle D und E zeigen, dass die Suche auch gestoppte Container findet (der Code ruft `docker ps -a`): Ein gestoppter Container der Instanz zählt für die Frage „gibt es einen Container ohne Instanz“ wie ein
laufender. Der Code (gelesen) unterscheidet den Container ohne Instanz am leeren Wert des Labels `u-boot.instance`; dass der Schlüssel dort fehlt, misst M3.

### Eingang für den CR (M5)

**Darf tragen:**

- `u-boot devcontainer down` ohne `--instance` wirkt nur auf den Container ohne Instanz; sind Instanz-Container daneben vorhanden, bleiben sie unberührt (Fall B).
- In einem Projekt, in dem es nur Container von Instanzen gibt (laufend, gestoppt oder gemischt), endet `down` ohne `--instance` mit Exit-Code 10 und dem Hinweis auf `--instance`, und es stoppt keinen
  der Container (Fälle C bis E); `up` ohne `--instance` endet dort ebenso mit Exit-Code 10 und legt keinen Container ohne Instanz an (Fall C). Das stimmt mit
  [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) überein.
- In einem Projekt ohne jeden Container endet `down` mit Exit-Code 0, mit und ohne `--instance` (Fall A).

**Darf nicht tragen:**

- Dass ein `remove` ohne `--instance` dieselbe Auswahl trifft: `remove` gibt es nicht; die Messung zeigt nur, was `down` tut, und das Bild des Suchwegs, den ein `remove` nutzen könnte.
- Aussagen zu `shell` ohne `--instance` (interaktiv, nicht gemessen; der Code ruft denselben Suchweg, `containerBeforeStart`) und zum Standardprofil (gemessen ist das Sandbox-Profil; der gelesene Code des Suchwegs wertet nur Labels aus).
- Aussagen zu genau einer Instanz oder zu drei und mehr Instanz-Containern ohne Container ohne Instanz (gemessen: zwei Instanz-Container) und zu Containern anderer Projektordner mit demselben Namen.
- Aussagen zu `paused`, `restarting` und `created` (nicht gemessen).

**Frage an den Architect (und den CR, Entscheidung beim Projektinhaber), nicht entschieden:** Welcher Container ist bei `remove` ohne `--instance` gemeint? Der Suchweg von `up`, `down` und `shell` kennt
genau eine Antwort („der Container ohne Instanz, sonst Exit-Code 10 mit Hinweis“); für eine zerstörerische Aktion kann der CR dieselbe Antwort wählen oder eine andere Auswahl vorschreiben (etwa eine
ausdrückliche Wahl statt der stillen Auswahl). Das Protokoll liefert die Grundlage (Fälle A bis E), nicht die Entscheidung.

## Aufruf

```bash
# Binary des Slice-Stands (im Devcontainer-Container):
make compile BAU_COMPILE_OUT=<pfad>/u-boot
UBOOT=<pfad>/u-boot WORK=<messordner-ausserhalb-des-repos> tools/messung/slice-v1-devcontainer-remove-messung/00-alles.sh   # M1 bis M5
UBOOT=<pfad>/u-boot WORK=<messordner> tools/messung/slice-v1-devcontainer-remove-messung/05-aufraeumen.sh
```

Je Messung einzeln: `01-m1-volumes.sh`, `02-m2-entfernen.sh`, `03-m3-instanz.sh`, `04-m4-ueberleben.sh`, `04b-m5-down-ohne-instanz.sh` (je eine Rohdatei `rohdaten-m1.txt` bis
`rohdaten-m5.txt`). Lesehilfe: Zeilen `AUSSAGE …` sind die stabilen Aussagen, die ein zweiter Lauf wiederholen muss; Dauern, Container-IDs und
Zeitstempel streuen. Die Datei `rohdaten-m1-erster-lauf.txt` ist die Ausgabe des ersten Laufs von M1 (Dauern mit neu gebauten Images).

## Grenzen

- Ein Host, eine Engine (die verschachtelte des Devcontainers), Docker 29.9.3; ein Lauf je Messpunkt, M1 mit einem zusätzlichen ersten Lauf (nicht archivierte Skriptfassung, siehe „Läufe und ihre Herkunft“). Eine native Engine auf einem
  Linux-Host ist nicht gemessen; wie weit der Aufbau abweicht, ist unbekannt. `u-boot` hat die Engine als lokal behandelt.
- Gemessen sind die von `u-boot` erzeugten Konfigurationen (Standard- und Sandbox-Profil) mit den Optionen `nestedRuntime` und dem Feature `ollama`; die übrigen
  Features des Katalogs, Features aus externen Quellen und von Hand geschriebene Konfigurationen sind nicht gemessen (Ausnahme: ein von Hand angehängtes
  `VOLUME` im Dockerfile).
- Der Platzbedarf der Images ist nicht belastbar gemessen (Build-Cache); die Bauzeiten stehen mit der Einschränkung ihres Cache-Zustands in M1.
- Ein Modell von Ollama wurde nicht geladen; das Netz diente dem Bau der Images.
- Ob `docker volume prune`, `docker system prune` oder `docker compose down -v` anders wirken, ist nicht gemessen.
