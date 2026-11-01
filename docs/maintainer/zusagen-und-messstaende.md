# Zusagen und Messstände zum Benutzerhandbuch

Dieses Dokument trägt, was das [Benutzerhandbuch](../user/benutzerhandbuch.md) nicht trägt: wo `u-boot` zugesagt ist, was
nicht zugesagt ist und was in welcher Umgebung gemessen oder nicht gemessen ist. Das Handbuch beschreibt nur den Ist-Zustand
und die Voraussetzungen, die Nutzer erfüllen müssen. Die Aussagen stehen in der Reihenfolge der Handbuch-Abschnitte; der
Wortlaut ist aus der Handbuch-Version 2.0 (Commit `d0878f1c`) übernommen.

## 1. Rang und Lesart

- **Verbindlich ist das Lastenheft.** Die Zusage-Mengen stehen dort:
  [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Features),
  [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil),
  [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (verschachtelte Laufzeit),
  [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Egress-Beschränkung),
  [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten) und
  [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Devcontainer ohne Editor starten).
  Dieses Dokument wiederholt sie nicht und widerspricht ihnen nicht; laufen ein Wortlaut unten und das Lastenheft auseinander,
  gilt das Lastenheft.
- **Zugesagt, nicht zugesagt:** Zugehörigkeit zur Zusage-Menge des Lastenhefts. „Nicht zugesagt“ sagt nichts darüber, ob etwas
  funktioniert.
- **Gemessen, nicht gemessen:** Ein Lauf hat es beobachtet (Messprotokoll unter [`docs/messungen/`](../messungen/README.md)) oder
  es gibt keinen solchen Lauf. „Nicht gemessen“ heißt weder „funktioniert“ noch „funktioniert nicht“.
- **Fundstelle:** Je Abschnitt steht das Messprotokoll, in dem die Versionsnummer, der Host oder das Stichwort des Wortlauts
  vorkommt (Suche im Verzeichnis `docs/messungen/`). Die Zuordnung einer einzelnen Aussage zu einer Rohdatei ist damit nicht
  belegt; wo die Suche nichts fand, steht das.
- **Voraussetzung im Handbuch:** Wo eine Aussage für Nutzer eine Voraussetzung ist, führt das Handbuch sie als Voraussetzung
  (Abschnitt 4); die Prüflücke dazu steht nur hier.

## 2. Messprotokolle

| Verzeichnis | Gegenstand |
|---|---|
| [`slice-v1-devcontainer-start-fundament`](../messungen/slice-v1-devcontainer-start-fundament/README.md) | Hilfs-Image des Devcontainer-Starts (Größe, Bauzeit, erster `up`, Trivy) |
| [`slice-v1-cli-devcontainer-up-down`](../messungen/slice-v1-cli-devcontainer-up-down/README.md) | Einzelmessung des Befehls `devcontainer up` und `down` |
| [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md) | mehrere Devcontainer-Instanzen desselben Projektordners |
| [`slice-v1-devcontainer-ollama-messung`](../messungen/slice-v1-devcontainer-ollama-messung/README.md) | Ollama im Sandbox-Devcontainer (Docker 29.5.2 auf Colima 0.10.3, `arm64`) |
| [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md) | Ollama im Sandbox-Devcontainer auf `amd64` |
| [`slice-v1-sandbox-egress-neuladen-entwurf`](../messungen/slice-v1-sandbox-egress-neuladen-entwurf/README.md) | Egress-Liste im laufenden Container (Docker 29.5.2 auf Colima 0.10.3) |
| [`slice-v1-egress-wildcard-pruefung`](../messungen/slice-v1-egress-wildcard-pruefung/README.md) | Egress-Allowlist, Subdomains und Wildcard |
| [`slice-v1-egress-dns-upstream-schliessen`](../messungen/slice-v1-egress-dns-upstream-schliessen/README.md) | Pfad zum Upstream-Resolver für Prozesse ohne `sudo` |
| [`slice-v1-forge-cli-quellen-messung`](../messungen/slice-v1-forge-cli-quellen-messung/README.md) | Forge-CLIs `gh` und `glab`: Quelle, Lieferkette, Hosts, Token-Umgang |
| [`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md) | Cloud-Modelle von Ollama mit Claude Code, ohne Konto |
| [`slice-v1-ollama-cloud-messung-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md) | Cloud-Modelle von Ollama mit Claude Code, mit Konto |
| [`slice-cr-lastenheft-messstaende`](../messungen/slice-cr-lastenheft-messstaende/README.md) | CI-Lauf der Integrationstests auf `amd64` (Ollama, verschachteltes Podman und Docker) |
| [`slice-v1-devcontainer-bau-lokal`](../messungen/slice-v1-devcontainer-bau-lokal/README.md) | Messlauf M0 zum lokalen Bauen im Devcontainer |
| [`slice-v1-devcontainer-ressourcenlimits`](../messungen/slice-v1-devcontainer-ressourcenlimits/README.md) | Hand-Ergänzung von `--memory` und `--cpus` und erneutes `generate devcontainer` |

## 3. Aussagen je Handbuch-Abschnitt

### 3.1 Voraussetzungen (Handbuch, Abschnitt 1)

- **Zusage:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
- **Fundstelle:** [`slice-v1-devcontainer-start-fundament`](../messungen/slice-v1-devcontainer-start-fundament/README.md) und
  [`slice-v1-cli-devcontainer-up-down`](../messungen/slice-v1-cli-devcontainer-up-down/README.md) (Docker 29.9.1).
- **Wortlaut:**
  - Zeile der Voraussetzungs-Tabelle: „Docker (Client und lokale Engine; gemessen mit 29.9.1)“.
  - „Für `devcontainer up`, `devcontainer down` und `devcontainer shell` ist Podman als Host-Engine nicht zugesagt.“
  - README.md: „Docker Engine ≥ 24.0.0 or Podman ≥ 4.0 (drop-in supported; not promised for `devcontainer up`/`down`/`shell`)“;
    README.de.md: „Docker Engine ≥ 24.0.0 oder Podman ≥ 4.0 (Drop-in unterstützt; für `devcontainer up`/`down`/`shell` nicht zugesagt)“.

### 3.2 Den Devcontainer ohne Editor starten und stoppen (Handbuch, Abschnitt 4.9)

- **Zusage:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
- **Fundstelle:** Erster Aufruf (Größe, Bauzeiten, Trivy): [`slice-v1-devcontainer-start-fundament`](../messungen/slice-v1-devcontainer-start-fundament/README.md)
  (Docker 29.9.1, 2026-10-20). Erneutes `up`, Statusdatei, Instanzen: [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md)
  (Docker 29.9.2, Dev-Container-CLI 0.89.0); die Zuordnung der einzelnen Aussage ist nicht geprüft. Keine Fundstelle gefunden
  (Suche nach den Zahlen und Stichwörtern): die Zeiten von Strg-C (10,1 und 10,5 Sekunden), `u-boot` 0.13.0, die Benutzer-ID 4321,
  `Devcontainer already running`.
- **Wortlaut:**

**Erster Aufruf**

**Erster Aufruf.** `u-boot` baut das Hilfs-Image lokal aus gepinnten Quellen: die
Dev-Container-CLI aus der npm-Registry (feste Version, Lockdatei mit Prüfsummen,
Installation ohne Installationsskripte der Pakete) und die Basis-Images `node` und
`docker` (jeweils mit Digest). Der erste Aufruf braucht deshalb Netzzugriff und baut;
spätere Aufrufe nutzen das Image weiter. Ein erster Aufruf ohne Netz ist nicht
zugesagt. Einzelmessung (Linux, `amd64`, Docker 29.9.1, keine Zusage): Das Image ist
1716 MiB groß. Der Bau dauerte bei vorhandenen Basis-Images 8 bis 19 Sekunden. Der
erste `up` eines Projekts, dessen Basis-Image erst gezogen wurde, dauerte 52 Sekunden,
weitere Starts 1 bis 5 Sekunden. Trivy meldete im Betriebssystem-Anteil des Node-Images
597 Funde (HIGH und CRITICAL), in der Dev-Container-CLI selbst keinen. Das Image entsteht
bei Ihnen und wird weder mit `u-boot` ausgeliefert noch dauerhaft geprüft; alte
Hilfs-Images räumt `u-boot` nicht ab (`docker image prune`).

**Zugangsdaten** (Satz zu `remoteEnv`)

nie als Wert auf der Kommandozeile von `u-boot`. Bei der Form `remoteEnv` steht der Wert
weder in einem Image noch in den Labels des Devcontainers, in einer Datei im Projektordner
oder in den Volumes des Sandbox-Profils (geprüft). Zwei Grenzen gehören dazu: `remoteEnv` ist

**Im Container arbeiten** (Shell ohne `remoteEnv`)

In dieser Shell gilt `remoteEnv` aus `devcontainer.json` nicht: Im Sandbox-Profil ist
`GIT_TOKEN` dort nicht gesetzt (beobachtet in einem Lauf, Linux `amd64`, Docker 29.9.1). Eine
Variable geben Sie mit `docker exec -e <NAME>` ausdrücklich mit; `docker` übernimmt dann den Wert
aus Ihrer Umgebung, und ist die Variable dort nicht gesetzt, fehlt sie auch im Container
(geprüft in einem Lauf mit `GIT_TOKEN`, Linux `amd64`, Docker 29.9.2):

```bash
docker exec -it -e GIT_TOKEN -u <Benutzer> -w <Workspace> <Container> bash

**Was zurückbleibt; Abbruch mit Strg-C**

**Was zurückbleibt.** Neben dem Hilfs-Image (`u-boot-devcontainer-cli:…`) legt die
Dev-Container-CLI je Projekt Images `vsc-<Ordnername>-…` an (gemessen auch `…-uid`); sie
sind nach `down` und nach `--recreate` weiter vorhanden (getestet; `--recreate` baut sie neu).
Aufräumen mit `docker image rm` oder `docker image prune`.

**Abbruch mit Strg-C.** Beobachtet: Die CLI im Hilfs-Container beendet sich auf das Signal
nicht früher; die Ursache ist nicht untersucht. `u-boot` wartet bis zu zehn Sekunden (gemessen:
10,1 Sekunden in der CI, 10,5 Sekunden lokal bis zur Rückkehr) und entfernt dann den
Hilfs-Container. Es bleibt kein Hilfs-Container zurück. Der Dev-Container, den die CLI
bis dahin angelegt hat, bleibt bestehen.


**Erneutes `up` nach einem Fehler oder Abbruch**

**Erneutes `up` nach einem Fehler oder Abbruch.** `up` urteilt nach dem Zustand des
Containers, nicht nach dem Ablauf des einzelnen Aufrufs. Ohne `--recreate` verwendet ein
erneutes `up` den vorhandenen Container:

- **Gegen einen laufenden Container führt die Dev-Container-CLI das `postCreateCommand` und das
  `postStartCommand` nicht erneut aus**, das `postAttachCommand` aber bei jedem `up` (getestet mit
  der Dev-Container-CLI 0.89.0 an einem Projekt, dessen drei Befehle je Lauf eine Zeile in eine
  Datei im Container schreiben; Linux `amd64`, lokale Engine, Docker 29.9.2, ein Lauf). Ist der
  Container gestoppt, läuft nach `u-boot devcontainer down` und `up` das `postStartCommand`
  wieder, im selben Container, das `postCreateCommand` nicht; `up --recreate` legt den Container
  neu an. Die Zeile von `up` sagt bei einem laufenden Container `Devcontainer already running`,
  im JSON-Umschlag ist `alreadyRunning` `true`.
- **Meldet nach dem Aufruf die Statusdatei eines Startscripts `failed`, endet jedes erneute `up`
  ohne `--recreate` wieder mit Exit-Code 11**, auch wenn die Dev-Container-CLI selbst gelingt. Die
  Meldung nennt die Ursache der Fähigkeit (wie `doctor`), den Zeitpunkt der Statusdatei, dass der
  Container läuft und dieser Aufruf das Startscript nicht ausgeführt hat, und
  `u-boot devcontainer up --recreate`. `--recreate` hilft nach Behebung der Ursache; besteht sie
  fort, endete auch das neu angelegte `up` wieder mit Exit-Code 11 (ein Lauf, Linux `amd64`,
  Docker 29.9.2). Getestet mit dem Startscript der Egress-Beschränkung unter
  `onUnavailable: fail` ohne die Capability `NET_ADMIN`: das erste, das zweite (mit `--json`:
  `exitCode` 11) und ein drittes `up` endeten mit Exit-Code 11 (für die Startscripte von
  `nestedRuntime: podman` und `docker` nur im Unit-Test mit nachgebildeter Engine). Das braucht
  die `u-boot.yaml`: ohne sie meldet ein erneutes `up` einen früheren Abbruch nicht und endet
  mit Exit-Code 0 (getestet: dasselbe Projekt ohne `u-boot.yaml`, erstes `up` Exit-Code 12,
  zweites 0).
- **Ohne Statusdatei bleibt ein Abbruch unsichtbar.** Bricht ein Lebenszyklus-Befehl ab, der keine
  Statusdatei schreibt (ein eigenes `postStartCommand` oder `postCreateCommand`), merkt ein
  erneutes `up` es nicht. Beobachtet in je einem Lauf (Linux `amd64`, Docker 29.9.1, `u-boot`
  0.13.0): Endete `up` mit Exit-Code 12, weil ein `postStartCommand` mit einem Fehler endete,
  lief der Devcontainer weiter, und das folgende `up` endete mit Exit-Code 0, ohne das
  `postStartCommand` erneut auszuführen. Nach Strg-C während eines `postCreateCommand` ebenso:
  Der Devcontainer lief weiter, das folgende `up` endete mit Exit-Code 0, und das abgebrochene
  `postCreateCommand` lief nicht erneut. Andere Fehler (etwa ein gescheiterter Bau) sind nicht
  gemessen.

**Benutzer-ID; Features und Reichweite**

**Benutzer-ID.** Die Benutzer-ID im Container wird nicht an den Besitzer des Projektordners
angepasst (beobachtet bei einem Ordner mit der UID 2000: der Benutzer behielt 1000). Die CLI
läuft im Hilfs-Container als `root`; ob das der Grund ist, ist nicht geprüft. Dateien, die der Benutzer im
eingebundenen Projektordner anlegt, gehören im Standardprofil der Benutzer-ID aus
`devcontainer.user.uid` (getestet mit 4321; Standard 1000). Setzen Sie sie auf Ihre
Host-UID, wenn Sie die Dateien als Host-Benutzer bearbeiten wollen.

**Features und Reichweite.** Als Katalog-Feature ist `node` im Test installiert worden und
nutzbar; andere Features sind nicht geprüft. Zugesagt ist der Befehl auf einem Linux-Host mit
lokaler Docker-Engine auf `amd64`. **Nicht zugesagt:** macOS mit Colima, Docker Desktop,
rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux auf `arm64` und Windows.
Der Befehl selbst ist auf keinem dieser Hosts gemessen; auf Colima ist nur der Messweg
(die Dev-Container-CLI in einem Wegwerf-Container) gelaufen. `u-boot up` startet den Devcontainer nicht mit.

**Mehrere Devcontainer je Projekt (Instanzen)** (Handbuch, Abschnitt 4.9)

- **Zusage:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Absatz zu den Instanzen): bis zu
  vier gleichzeitig laufende Instanzen, jede mit höchstens einer CPU und 3 GiB (Grenze samt Swap), im Sandbox-Profil mit
  `nestedRuntime: docker` und dem Feature `ollama`, im Zustand direkt nach dem Start, Linux `amd64`, lokale Engine. Nicht zugesagt:
  mehr als vier, `nestedRuntime: podman`, Instanzen ohne Grenze, das Verhalten unter Last, andere Hosts, die Isolation im Standardprofil.
- **Fundstelle:** Messwerte (Container, Labels, Volume-Namen, Inhaltsprobe, `up` ohne Flag):
  [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md) (Docker 29.9.2,
  Dev-Container-CLI 0.89.0). Das Verhalten von `u-boot` selbst prüft der Docker-Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated`
  (zwei Instanzen und der Container ohne Instanz, Sandbox-Profil mit Standard-Einstellungen: ein Workspace-Volume je Container;
  ohne `nestedRuntime`, ohne `ollama`, ohne Grenze, ohne Prüfung der Hostressourcen). Die Volumes für Docker-Speicher und Modelle,
  die Grenze von einer CPU und 3 GiB und vier Instanzen läuft die Suite nicht; sie stehen nur in der Messung, die Mount-Texte
  der Volumes in Unit-Tests. Im Standardprofil
  `TestDevcontainerInstance_Docker_StandardProfileSharesTheProjectFolder` (drei Container).
- **Handbuch:** Abschnitt 4.9 nennt keine Zahl gleichzeitiger Instanzen und keine Grenze je Instanz; beides steht nur im Lastenheft und in diesem Abschnitt.

### 3.3 Katalog-Features und ihre Lieferkette (Handbuch, Abschnitt 4.10, Feature-Tabelle)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features).
- **Fundstelle:** `gh`, `glab`: [`slice-v1-forge-cli-quellen-messung`](../messungen/slice-v1-forge-cli-quellen-messung/README.md);
  `ollama`: [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md) und
  [`slice-v1-devcontainer-ollama-messung`](../messungen/slice-v1-devcontainer-ollama-messung/README.md).
- **Wortlaut** (Zellen der Feature-Tabelle):
  - `ollama`: „SHA-256 geprüft, vom Upstream nicht signiert; `u-boot generate devcontainer` weist das aus“.
  - `gh`: „SHA-256 geprüft, nur das Programm `gh` entpackt; keine Signatur gefunden (nach Suchmuster, nicht abschließend), die
    GitHub-Attestierung des Projekts prüft `u-boot` nicht; `u-boot generate devcontainer` weist das aus“.
  - `glab`: „SHA-256 geprüft, nur das Programm `glab` entpackt; keine Signatur gefunden (nach Suchmuster, nicht abschließend);
    `u-boot generate devcontainer` weist das aus“.

### 3.4 Verschachteltes Podman (Handbuch, Abschnitt 4.11, Podman)

- **Zusage:** [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
- **Fundstelle:** [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md)
  (Podman 6.1.4, `kernel.apparmor_restrict_unprivileged_userns`, `/dev/fuse`, `registries.conf`),
  [`slice-v1-devcontainer-ollama-messung`](../messungen/slice-v1-devcontainer-ollama-messung/README.md) (Colima 0.10.3),
  [`slice-v1-egress-dns-upstream-schliessen`](../messungen/slice-v1-egress-dns-upstream-schliessen/README.md) und
  [`slice-v1-sandbox-egress-neuladen-entwurf`](../messungen/slice-v1-sandbox-egress-neuladen-entwurf/README.md) (`debian:stable-slim`),
  [`slice-cr-lastenheft-messstaende`](../messungen/slice-cr-lastenheft-messstaende/README.md) (CI-Lauf 2026-10-25).
- **Wortlaut:**

**Geräte**

`/dev/net/tun` braucht das Netzwerk der verschachtelten Container (`pasta`);
ohne das Device scheitert jedes `RUN` eines `podman build`. `/dev/fuse` dient dem
Speicher der verschachtelten Container. Beide Devices müssen auf dem
Container-Host vorhanden sein; fehlt eines, lehnt Docker den Start ab (gemessen: `error gathering device information while adding custom device … no such file or directory`), das Startscript läuft dann nicht. `u-boot doctor` prüft `/dev/net/tun` nicht und meldet ein fehlendes `/dev/fuse` auf einem Linux-Host, nicht aus einem Container, als Fehler (Exit-Code 11). Abhilfe: die Kernelmodule auf dem Host laden (`modprobe fuse`, `modprobe tun`; `tun` nicht gemessen) oder `nestedRuntime: none` wählen. Bestehende Projekte mit `nestedRuntime: podman` brauchen das Device nach `u-boot generate devcontainer` ebenfalls.

**Images bauen**

**Images bauen.** `podman build` und `docker build` sind im Container derselbe
Aufruf (der `docker`-Alias ist ein Symlink auf `podman`); gebaut wird mit der in
Podman eingebauten Buildah-Bibliothek. Nicht vorhanden sind das Kommando
`buildah`, `docker buildx`, `docker compose` und BuildKit-Frontends
(`# syntax=`); wer sie braucht, wählt `nestedRuntime: docker`. Verwenden Sie in
Dockerfiles voll qualifizierte Namen (`FROM docker.io/library/alpine:3`): Namen
ohne Registry sucht Podman nur in Docker Hub (`registries.conf` mit
`unqualified-search-registries = ['docker.io']`; gemessen: `alpine`
löst auf, `mycompany/img` nicht; andere Kurznamen sind nicht gemessen). Images und
Schichten liegen im Volume
`<projekt>-containers-${devcontainerId}`; gemessen: das Image `debian:stable-slim`
mit einer Paketinstallation und drei Builds belegte rund 140 MB, das Volumen
wächst mit jedem Image. Fehlt `/dev/fuse` im Container (der Container wurde ohne
das Device gestartet, etwa nach geänderten `runArgs`), richtet das Startscript
`vfs`-Speicher ein (Benutzer-Konfiguration unter `~/.config/containers/`; der
temporäre Platz für Image-Kopien liegt im Volume), und `podman run`, `podman build`
und `docker build` laufen damit; das kopiert aber jede Schicht vollständig.
Gemessen ist das mit eingehängtem Volume.

**Gemessen und zugesagt**

**Gemessen** (Podman 6.1.4; Docker 29.5.2 auf Colima 0.10.3, arm64; in der VM
`kernel.apparmor_restrict_unprivileged_userns=0`, auf dem Container-Host `/dev/net/tun`): `podman build` mit `FROM`,
`RUN` und `COPY` baut kalt, aus dem Layer-Cache und mit `--no-cache`; `docker
build` über den Alias ebenso; ein `RUN` mit Netz zu einem erlaubten Host
(`apt-get update`) läuft auch hinter der Egress-Restriktion mit der
Default-Allowlist (gemessen mit `debian:stable-slim`); für `alpine:3` hinter der Restriktion sind Pull und Build gemessen, ein `RUN` mit Netz nicht. Der Blob-Host von Docker Hub kann sich ändern; ein nicht gelisteter Name löst hinter der Restriktion nicht auf. Auf Linux x86_64 (`amd64`, native Docker Engine 29.9.1) läuft ein `podman build` hinter der Restriktion ebenfalls durch (Abschnitt zur Egress-Beschränkung).
**Zugesagt** ist `podman build` für Images von Docker Hub (`docker.io`) auf einem Linux-Host mit lokaler Engine (`amd64`) und unter Colima (`arm64`). **Nicht
gemessen:** andere Hosts, Docker Desktop, Podman als Host-Engine; die CI übersprang am 2026-10-25 die Docker-Tests zum Bau mit `podman build` und zum Start von `podman`, weil der Runner
`kernel.apparmor_restrict_unprivileged_userns=1` hat (der Test zum falschen Pin mit `nestedRuntime: podman` bestand). **Nicht zugesagt:** `buildx`, `docker compose` und `# syntax=`-Frontends.

**Ubuntu 24.04**

Unter **Ubuntu 24.04** (und auf Hosts mit
`kernel.apparmor_restrict_unprivileged_userns=1`) blockiert dieser Schalter
verschachtelte User-Namespaces auch mit `apparmor=unconfined`; der Start endet
dann mit Exit-Code 11. Abhilfe: auf dem Container-Host
`sysctl kernel.apparmor_restrict_unprivileged_userns=0` setzen oder
`nestedRuntime: none` wählen. Gemessen: Die Linux-VM von Colima 0.10.3 (Ubuntu
24.04.4, Kernel 6.8, Docker 29.5.2) trägt standardmäßig den Wert 1; dort endet
der Start mit den Lockerungen mit Exit-Code 11. Ohne `apparmor=unconfined`
startet das Skript, aber `podman build` scheitert an `mount` und `remount`
(Profil `docker-default`). Mit dem Wert 0 in der VM läuft der Build (siehe oben).
Einen dauerhaften Wert setzen Sie in der VM (`colima ssh`) in einer Datei unter
`/etc/sysctl.d/`.

### 3.5 Fallbacks sehen (Handbuch, Abschnitt 4.11, Statusdateien)

- **Zusage:** [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
- **Fundstelle:** keine gefunden (Suche nach dem Ersetzen der Statusdatei ohne `sudo` unter Colima: ohne Treffer).
- **Wortlaut:**

**Grenze:** Die Datei stammt aus einem Container, in dem ein Agent läuft, und ist
dort mit Root-Rechten schreibbar; gemessen ist außerdem, dass ein Prozess ohne
`sudo` sie ersetzen kann (die Datei gehört `root`, aber das Verzeichnis
gehört dem Benutzer `vscode`; ein Lauf auf Colima, arm64). Sie ist ein Hinweis,
kein Beleg: **„Nicht gemeldet“ heißt nicht „nicht aufgetreten“.** `u-boot` gibt nie den Inhalt der
Datei aus, nur die geprüften Felder; eine Datei mit anderem Format, mit mehr als
512 Byte oder die keine reguläre Datei ist, gilt als ungültig.

### 3.6 Eigener Docker-Daemon (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
- **Fundstelle:** `production.cloudfront.docker.com`: [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md) und
  [`slice-v1-egress-dns-upstream-schliessen`](../messungen/slice-v1-egress-dns-upstream-schliessen/README.md).
- **Wortlaut:**

- Unter Podman als Host-Engine ist `docker` nicht zugesagt.
- Der Host `production.cloudfront.docker.com` der Default-Allowlist gilt auch hier, ist aber für `nestedRuntime: docker` nicht gemessen (gemessen ist der Pull mit `podman`).

### 3.7 Internet-Beschränkung (Handbuch, Abschnitt 4.11, Egress)

- **Zusage:** [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).
- **Fundstelle:** [`slice-v1-egress-dns-upstream-schliessen`](../messungen/slice-v1-egress-dns-upstream-schliessen/README.md)
  (Docker 29.9.1, 2026-10-19; 3 650 bis 5 640 Anfragen je Sekunde),
  [`slice-v1-egress-wildcard-pruefung`](../messungen/slice-v1-egress-wildcard-pruefung/README.md) (Subdomains, `CNAME`),
  [`slice-v1-sandbox-egress-neuladen-entwurf`](../messungen/slice-v1-sandbox-egress-neuladen-entwurf/README.md)
  (Docker 29.5.2 auf Colima 0.10.3, `dnsmasq` 2.91, Dev-Container-CLI 0.89.0, 2026-10-18),
  [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md) (`i/o timeout` beim Pull).
- **Wortlaut:**

**Namensauflösung und ihre Grenze**

**Namensauflösung und ihre Grenze:** Die Firewall lässt DNS (Port 53, UDP und
TCP) zu den Upstream-Resolvern nur für den Benutzer des lokalen DNS-Servers zu
(`nobody`; das Startscript legt ihn in der Konfiguration des DNS-Servers fest und
leitet die Regel aus demselben Wert ab). Ein Prozess ohne `sudo` bekommt deshalb
über die Upstream-Resolver keine Antwort für einen Namen außerhalb der Liste;
erlaubte Namen löst der lokale DNS-Server weiter auf. Ein Docker-Test belegt das
(`make test-docker`). Messstand: nativer Linux-x86_64-Host, Docker 29.9.1, zwei
Läufe; **nicht gemessen und nicht zugesagt:** Colima, Docker Desktop, arm64, die
CI, `nestedRuntime: docker`, ein benutzerdefiniertes Docker-Netz, Podman-Netzwerkmodi
außer dem Standard und `--network=host`, Upstream-Resolver mit IPv6-Adresse.

**Unternamen erlaubter Namen**

- Anfragen nach Unternamen erlaubter Namen: Der lokale DNS-Server leitet sie an den
  Upstream-Resolver weiter. Das ist ein Weg, Daten zum zuständigen DNS-Server des
  erlaubten Namens zu schicken (im Laboraufbau gemessen: höchstens 3 650 bis 5 640
  Anfragen je Sekunde, ohne Test, der es pinnt); bei den Standard-Namen betreibt ihn
  kein Dritter, bei von Ihnen ergänzten Namen kann es anders sein (nicht gemessen).
  Der Betreiber des Upstream-Resolvers sieht diese Anfragen.


**Verschachtelte Container**

**Verschachtelte Container:** Ein verschachtelter Podman-Container fragt als Benutzer
`vscode`; seine Anfragen an die Upstream-Resolver des Docker-Hosts werden mit der Regel
verworfen. Ein Name außerhalb der Liste braucht darin etwa 36 Sekunden bis zur
Fehlermeldung statt etwa 20 Millisekunden (der erste Nameserver lehnt ab, die
verworfenen Rückfälle laufen in die Zeitgrenze); ein `podman build` mit einem
`RUN`-Schritt, der einen solchen Namen auflöst, dauert um diese Zeit länger
(nachgestellter Ablauf: 43 statt 7 Sekunden). Diese Zeiten und das Verhalten von `--dns`,
der Standard-Namensauflösung und `--network=host` sind gemessen; kein Test pinnt sie.
Dass der Bau hinter der Beschränkung durchläuft, pinnt ein Docker-Test (Linux `amd64`, native Docker Engine 29.9.1). `podman run --dns 8.8.8.8`
(ein Upstream-Resolver des Messhosts) löst nichts mehr auf, auch keinen erlaubten Namen;
die Standard-Namensauflösung und `--network=host` bleiben nutzbar. Auf Colima (arm64)
und in der CI ist das nicht gemessen.

**Abbruch des Startscripts**

Das Startscript bricht mit Exit-Code 11 ab und schreibt `failed` in die Statusdatei,
auch bei `onUnavailable: warn`, wenn der Benutzer `nobody` im Image fehlt, die
nftables-Regeln sich nicht anwenden lassen, der lokale DNS-Server nicht startet oder nicht
unter diesem Benutzer läuft: Ein Weiterlaufen würde die Beschränkung still aufweichen oder
den DNS-Server sperren. Der Container bleibt danach in einem von drei Zuständen (ein
Docker-Test zählt dort die Tabellen, die Kette mit der Richtlinie `drop` und den Eintrag in `/etc/resolv.conf`):

- **Benutzer fehlt:** Das Script bricht ab, bevor es Regeln anlegt. Im ersten Lauf hat der
  Container dann **keine** Beschränkung und sein DNS ist unverändert; die Meldung steht auf
  der Standardfehlerausgabe des Startscripts (Container-Log), und `u-boot doctor` zeigt `failed`.
- **Die nftables-Regeln lassen sich nicht anwenden:** Das Script ersetzt die Regeln in einer
  einzigen Transaktion und ändert bei einem Fehler nichts. Im ersten Lauf hat der Container
  dann **keine** Beschränkung und sein DNS ist unverändert; in einem wiederholten Lauf bleiben
  die Regeln und der DNS-Server des vorigen Laufs unverändert stehen (ein Docker-Test prüft
  beides mit einer ungültigen Regelzeile). Ein Host, auf dem `nft` die IPv6-Tabelle nicht
  anlegen kann, endet ebenso ohne Regel mit Exit-Code 11 (gelesen im Script; auf einem
  solchen Host nicht gelaufen).
- **DNS-Server startet nicht oder läuft unter einem anderen Benutzer:** Die Regeln mit der
  Richtlinie `drop` stehen. `/etc/resolv.conf` wird nicht auf den lokalen
  DNS-Server gesetzt (ein Docker-Test prüft das); in einem wiederholten Lauf bleibt der Eintrag
  des vorigen Laufs stehen (gelesen, kein Test).

Das Startscript der Beschränkung ist im `postStartCommand` mit `&&` vor das Startscript des
verschachtelten Docker-Daemons gesetzt (gelesen, kein Test mit Abbruch): Bei einem Abbruch startet
mit `nestedRuntime: docker` kein Daemon, und die Statusdatei `daemon` bleibt die des
vorigen Starts.

**Tabelle der Standard-Hosts, Zeile `ollama`**

| Feature `ollama` | `registry.ollama.ai`, `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` (Modell-Pull; Stabilität des zweiten Namens nicht gemessen: siehe [Ollama](../user/benutzerhandbuch.md#ollama-optional)) |

**`u-boot egress list`**

- **Was der Befehl liest:** nur `u-boot.yaml` und die Git-Konfiguration des
  Projekts (`.git/config`). Er braucht weder Netz noch Docker noch einen
  laufenden Container und ändert nichts. Ein Test vergleicht seine Liste mit der
  Zeile `hosts=` des erzeugten `egress-init.sh`; beide stammen aus derselben
  Berechnung.

**Subdomains, `CNAME`, Messstand**

- **Subdomains:** Ein Test ruft die Namen auf und liest den HTTP-Status (`000`:
  der Name löst nicht auf oder die Adresse wird verworfen; der Test trennt beides
  nicht). Mit dem Eintrag `example.org` sind `example.org` und `www.example.org`
  erreichbar, mit dem Eintrag `amazonaws.com` auch `s3.us-east-1.amazonaws.com`
  (zwei Ebenen darunter), `example.com` ist es nicht. Mit dem Eintrag
  `www.example.org` ist `www.example.org` erreichbar und `example.org` bei der
  Verbindung per Name nicht. Gemessen (je ein Lauf, Stand unter *Messstand*): Beim
  Eintrag `api.github.com` lösten `github.com` und `www.github.com` nicht auf, und
  mit dem Eintrag `example.org` hatten `example.org` und `www.example.org`
  dieselben IP-Adressen. Der Eintrag `www.example.org` selbst ist nicht gemessen;
  die Aussage zu `example.org` gilt analog zum Eintrag `api.github.com`. Ob eine Verbindung zu `example.org` über die Adresse möglich ist,
  wenn es sich mit dem eingetragenen Namen eine Adresse teilt, ist nicht geprüft
  und nicht ausgeschlossen. Einen Eintrag der Form `*.example.org` gibt es nicht:
  Tragen Sie `example.org` ein, `u-boot` weist den Eintrag mit diesem Hinweis ab.
- **Namen hinter einem CNAME:** Zeigt die DNS-Antwort eines erlaubten Namens
  über einen `CNAME` auf einen Namen außerhalb der Liste, war der erlaubte Name
  in den drei gemessenen Fällen erreichbar, ohne dass der Zielname eingetragen
  wurde; der Zielname selbst löst nicht auf. Die drei Fälle (`www.mozilla.org`,
  `files.pythonhosted.org`, `deb.debian.org`) liegen alle hinter demselben
  CDN-Anbieter und wurden einmal gemessen. Ob jeder `CNAME`-Fall so läuft, ist
  nicht gemessen, und es gibt keinen Test dafür.
- **Messstand:** Die Aussagen dieses Abschnitts sind auf Docker 29.5.2 mit Colima 0.10.3 (arm64) und
  `dnsmasq` 2.91 gemessen. Zusätzlich gibt es Messungen auf einem nativen Linux-x86_64-Host mit Docker 29.9.1
  (2026-10-19, zwei Läufe): der Pull von Modellen hinter der Beschränkung, das erneute Ausführen
  des Startscripts (frühere Fassung) und der Verkehr verschachtelter Container (siehe unten und den Abschnitt zu
  Ollama). Docker Desktop ist nicht gemessen.

**Neuaufbau nach einer Änderung**

- Nach einer Änderung führen Sie `u-boot generate devcontainer` aus und bauen
  den Container neu auf. Die Liste steht im Startscript, das beim Bau in das
  Image kopiert wird; ein Neustart des bestehenden Containers führt das alte
  Script mit der alten Liste aus. Gemessen mit der Dev-Container-CLI 0.89.0:
  Stoppen und erneutes Starten (`devcontainer up`) ließ die alte Liste gelten;
  erst `devcontainer up --remove-existing-container` baute das Image neu und
  setzte die neue Liste durch. Gemessen in einem Lauf auf Docker 29.5.2 mit
  Colima 0.10.3 (arm64); nicht gemessen: Docker Desktop, Linux ohne VM, `amd64`
  und die Aktion Ihres Editors (etwa VS Code), die dem Neuaufbau entspricht.
  Ändert `u-boot generate devcontainer` ein bestehendes `egress-init.sh`, gibt
  der Befehl dazu eine Warnung aus: Ein bestehender Container behält die
  Egress-Liste des Images, aus dem er gebaut wurde, und ist neu aufzubauen
  (`devcontainer up --remove-existing-container`). Die Warnung erscheint nur, wenn der Block eines bestehenden
  `egress-init.sh` ersetzt wird, also nicht beim ersten Schreiben der Datei (auch
  nicht bei der ersten Aktivierung der Beschränkung, bei der der Befehl das Script
  anlegt und `devcontainer.json` und das Dockerfile ändert), nicht beim Abschalten
  (der Befehl ändert `devcontainer.json` und das Dockerfile und lässt das Script
  stehen), nicht ohne Änderung und nicht unter `--dry-run`; unter `--quiet` entfällt sie, mit `--json` steht sie in
  `diagnostics[]`, und der Umschlag trägt dann den Status `warn` (Exit-Code 0).

**Umgehung durch `sudo` und verschachtelte Container**

- Wer die Fähigkeit hat und das `sudo` des Base-Images nutzen kann, kann die
  Regeln aufheben. Der Verkehr verschachtelter Container und die
  iptables-Regeln eines verschachtelten Docker-Daemons umgehen die Beschränkung;
  `u-boot doctor` warnt bei `docker` mit Beschränkung. Gemessen auf Linux (`amd64`) in
  zwei Läufen: Ein verschachtelter Docker-Container erreichte `example.com` und eine
  gesperrte Adresse; ein verschachtelter Podman-Container nicht (Zeitüberschreitung).
  Die Namensauflösung eines verschachtelten Podman-Containers läuft über den lokalen DNS-Server
  oder wird verworfen (siehe *Namensauflösung und ihre Grenze* im Abschnitt zur Egress-Beschränkung).

**Startscript erneut ausführen; Pull direkt nach dem Start**

- **Startscript erneut ausführen:** Führen Sie das Startscript in einem laufenden
  Container mit aktiver Beschränkung erneut aus (mit `sudo`), ersetzt es beide
  Tabellen in einer einzigen nftables-Transaktion, tauscht danach den lokalen
  DNS-Server aus und schreibt `/etc/resolv.conf` zuletzt. Das ist die Struktur
  des Scripts; ein Test pinnt sie. Ein Docker-Test führt das Script zwölfmal
  hintereinander aus: Jeder der ersten elf Läufe endet mit Exit-Code 0 und dem Status `ok`, es
  läuft genau ein DNS-Server unter dem Benutzer der Regeln, beide Tabellen
  bestehen (die IPv6-Tabelle mit einer verwerfenden Kette und der Regel für
  Loopback), und ein erlaubter Name löst auf; nach dem zwölften Lauf ist die Liste
  der erreichbaren Adressen leer. Währenddessen versucht eine Sonde laufend eine
  Verbindung zum Gateway des Containers (Port 9); die Beschränkung hält sie auf.
  Ein Treffer färbt den Test nur, wenn das Gateway vorher ohne Beschränkung auf die
  Sonde geantwortet hat; antwortet es nicht, meldet der Test nur, dass die Sonde auf
  diesem Rechner nichts aussagt, und bleibt grün. **Nicht zugesagt** ist, dass nie ein
  Zeitraum ohne Beschränkung entsteht oder die Auflösung eines erlaubten Namens
  nie aussetzt: Eine Zählung von null heißt „kein Treffer“, nicht „kein Fenster“,
  weil die Sonde nur in Abständen von Millisekunden abtastet. Schlägt die
  Transaktion fehl, meldet das Script das, schreibt den Status `failed` und endet
  mit Exit-Code 11, auch bei `warn`; das bisherige Regelwerk und der bisherige
  DNS-Server bleiben bestehen (ein Test pinnt das mit einer ungültigen
  Regelzeile).
  **Nicht gemessen** für diese Fassung: die Zeiten auf `amd64`, Docker Desktop und
  in der CI, die Zeiten unter Last des Hosts, mit verschachtelter Laufzeit und die
  IPv6-Tabelle (geprüft ist nur, dass beide Tabellen bestehen und die
  IPv6-Kette verwirft und Loopback zulässt, nicht ein Fenster). Gemessen ist nur die **frühere Fassung**, die ihre Tabellen löschte
  und danach neu aufbaute (Linux, `amd64`, ein Host, zwei Messungen zu je 30 Läufen):
  In 27 und 28 Läufen wurde währenddessen ein Verbindungsversuch zum Gateway
  (Port 9) nicht verworfen, sondern vom Gateway abgelehnt, bis zu 48 und 50 ms
  lang (eine Untergrenze); in 30 und 27 Läufen ließ sich ein erlaubter Name für
  30 bis 76 Millisekunden nicht auflösen. Diese Zahlen sagen nichts über die
  neue Fassung. Die Beschränkung bleibt eine Leitplanke und keine
  Sicherheitsgrenze.
- Der Pull eines Images mit verschachteltem Docker lief direkt nach dem Start der Beschränkung in
  einem von zwei vollständigen Läufen in eine Zeitüberschreitung (`i/o timeout`) und gelang 15 Sekunden später;
  im anderen Lauf gelang er sofort. Die Ursache ist nicht untersucht.

### 3.8 Login-Dienst `sshd` (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features).
- **Fundstelle:** keine gefunden (Suche nach `sshd`, `2222`, `ssh-keygen` und `authorized_keys` in `docs/messungen/`: ohne Treffer).
- **Wortlaut:**

**Port und Messweg**

Das offizielle Feature `sshd` installiert einen SSH-Server im Sandbox-Container.
Er öffnet einen Login-Dienst mit Passwort- und Schlüssel-Login, Root-Login ist
erlaubt; lesen Sie vor Schritt 1 den [Sicherheitshinweis](../user/benutzerhandbuch.md#sicherheitshinweis-sshd)
am Ende des Abschnitts. Der Server lauscht beim Containerstart auf Port 2222
(gemessen mit der Dev-Container-CLI 0.89.0 auf dem Basis-Image
`mcr.microsoft.com/devcontainers/base:debian`). Im Sandbox-Profil ist das Feature
erlaubt; `u-boot generate devcontainer` weist den Login-Dienst mit einer Warnung
aus, die auch unter `--quiet` sichtbar bleibt (anders als die Drittanbieter-Warnung
eines anderen Features). Ohne Sandbox-Profil erscheint diese Warnung nicht.

Die Schritte stehen in der Reihenfolge der Ausführung. Gemessen sind:

- am Feature selbst: Port 2222, die Werte der SSH-Konfiguration und die gesperrten
  Konten;
- an einem Näherungsaufbau (Wegwerf-Container aus dem Basis-Image mit dem
  Debian-Paket `openssh-server` auf Port 2222 statt des Feature-Skripts, ohne
  Editor- oder Engine-Portweiterleitung, ohne Dev-Container-CLI): die Schlüsselerzeugung
  mit `ssh-keygen -t ed25519`, die Rechte `700` und `600`, und der Aufruf
  `ssh -i … -p 2222 vscode@localhost` ohne Passwort.

**Nicht gemessen** sind ein SSH-Login mit dem Feature selbst, die
Portweiterleitung von Editor oder Container-Engine und `docker exec`; sie sind im
Text so gekennzeichnet.


**Neuaufbau, Schlüssel, Weiterleitung, Aufruf**

Das Installationsskript eines Features läuft beim Image-Build. Bauen Sie den
Container deshalb nach `generate devcontainer` neu auf. Nicht gemessen ist, ob
ein bloßer Neustart des bestehenden Containers nach einer geänderten Feature-Liste
genügt; für die Egress-Liste ist gemessen, dass er nicht genügt (Abschnitt
[Internet-Beschränkung (Egress)](../user/benutzerhandbuch.md#internet-beschränkung-egress)).

     `~/.ssh/authorized_keys` des Container-Benutzers `vscode` eintragen, in einem
     Terminal im Container (Editor-Terminal oder `docker exec`; `docker exec`
     nicht gemessen):

des Sandbox-Profils enthalten keine Port-Veröffentlichung. Sie leiten den Port
selbst auf den Host weiter, über die Port-Weiterleitung Ihres Dev-Container-Werkzeugs
(nicht gemessen, wie ein Editor oder die CLI das tut) oder eine Port-Option der
Container-Engine (nicht gemessen). Veröffentlichen Sie den Port nur an der
Loopback-Adresse des Hosts, nach dem Muster `127.0.0.1:2222:2222`, nicht als

Beim Passwort-Login entfällt `-i`. Der Aufruf ist im Näherungsaufbau gemessen,
nicht mit dem Feature selbst.

### 3.9 Claude Code (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
  [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).
- **Fundstelle:** keine gefunden für Claude Code `2.1.302`, Feature `1.0.5` und „Invalid API key“ (Suche in `docs/messungen/`:
  ohne Treffer). Zu Claude Code `2.1.310` siehe [`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md).
- **Wortlaut:**

Das Feature `claude-code` installiert die Claude Code CLI aus dem Feature von
Anthropic (`ghcr.io/anthropics/devcontainer-features/claude-code`, Tag `1`, zum
Messzeitpunkt Version 1.0.5). Es ist ein **Drittanbieter-Feature**: `u-boot
generate devcontainer` weist es aus (Herausgeber `anthropics`; die Meldung
entfällt unter `--quiet`), und sein Installationsskript läuft beim Image-Build.
Im Sandbox-Profil ist es erlaubt. `u-boot` startet den Agenten nicht und setzt
keinen Berechtigungsmodus.


- **Das Feature `node` ist nötig:** Das Feature installiert die CLI mit `npm`.
  Gemessen auf dem Basis-Image `mcr.microsoft.com/devcontainers/base:debian`
  (Debian trixie) bricht der Build ohne Node.js ab, weil der Versuch des Features,
  Node.js selbst nachzuinstallieren, Node.js ohne `npm` liefert. `u-boot generate
  devcontainer` weist darauf hin, wenn `claude-code` ohne `node` aktiv ist.

- **Zugangsdaten:** Im Sandbox-Profil trägt `devcontainer.json` bei aktivem Feature
  die Referenz `"ANTHROPIC_API_KEY": "${localEnv:ANTHROPIC_API_KEY}"` neben
  `GIT_TOKEN`, nie einen Wert. Setzen Sie `ANTHROPIC_API_KEY` in der Host-Umgebung,
  bevor Sie den Container starten (Quelle: Dokumentation von Anthropic, Stand
  2026-10-17: der Schlüssel hat Vorrang vor dem Login mit einem Abonnement). Ist die
  Variable auf dem Host nicht gesetzt, ist sie im Container nicht gesetzt (gemessen mit
  `u-boot devcontainer shell`, ein Lauf auf Linux `amd64`: weder in der Shell noch in der
  Konfiguration des Containers). Ohne Schlüssel meldet die CLI „Not logged in“ (gemessen
  mit der Dev-Container-CLI auf Colima, `arm64`, ohne Variable und mit leerer Variable;
  mit `u-boot devcontainer shell` nicht gefahren). Für
  Cloud-Modelle von Ollama gilt die Zusage nur bei nicht gesetztem `ANTHROPIC_API_KEY` im
  Aufruf; ob eine leere Variable dabei wie eine nicht gesetzte wirkt, ist nicht gemessen
  (siehe [Ollama](../user/benutzerhandbuch.md#ollama-optional)).
  Neben `ANTHROPIC_API_KEY` trägt `devcontainer.json` die Referenz auf
  `CLAUDE_CODE_OAUTH_TOKEN`; mit einem Platzhalter kommt die Variable über
  `u-boot devcontainer shell` im Container an, ungesetzt kommt sie ungesetzt an, leer
  gesetzt leer (je 2 Läufe), und `claude -p` mit dem Platzhalter endet mit der Antwort
  „401 OAuth access token is invalid“ (1 Lauf; Linux `amd64`, lokale Engine, 2026-11-01,
  `docs/messungen/slice-v1-devcontainer-claude-oauth-token/`). Ein Lauf mit einem
  gültigen Token ist nicht gefahren; mit gesetztem Token im Ollama-Ablauf ist nicht
  gemessen. Andere Variablen reicht `u-boot` nicht durch; ein Login mit `/login` im
  Container gehört nicht zur Zusage und ist nicht gemessen. Wohin ein Login im Container seine Zugangsdaten
  schreibt, ist nicht gemessen; das Home von `vscode` liegt nicht im
  Workspace-Volume (gemessen: das einzige Volume des Containers ist
  `/workspaces/<projekt>`).
- **Was beim Start entsteht (gemessen):** Die CLI legt `~/.claude` und
  `~/.claude.json` im Home von `vscode` an; mit dem Schlüssel aus der Umgebung
  stand der Schlüssel in keiner dieser Dateien. Das Feature bringt kein `mounts`,
  `capAdd`, `privileged` und keinen Entrypoint mit (`installsAfter`: das Feature
  `node`), und im Image gibt es kein `~/.claude`; das Feature installiert außerdem
  ein Firewall-Skript `/usr/local/bin/init-firewall.sh`, das **nicht ausgeführt wird**.

- **Egress:** Mit der Beschränkung erlaubt `u-boot` für `claude-code` nur zwei Hosts
  (Entscheidung des Projektinhabers): `api.anthropic.com` (gemessen: der API-Host,
  den die CLI zur Laufzeit abfragt) und `registry.npmjs.org`. Die Beschränkung wirkt
  erst beim Containerstart, der Image-Build davor ist nicht beschränkt; den Host
  `registry.npmjs.org` bringt auch das nötige Feature `node` mit, ein Bedarf der
  CLI zur Laufzeit ist nicht gemessen. Die Hosts der **Anmeldung** (`claude.ai`,
  `platform.claude.com`, `claude.com`; laut Dokumentation von Anthropic, nicht
  gemessen: ein echter Login ist ohne Konto nicht messbar) stehen **nicht** in der
  Liste: Wer sich im Container anmeldet, ergänzt sie in `devcontainer.sandbox.egress.allow`.
  Der Weg über einen Abo-Login im Container ist nicht zugesagt und nicht gemessen;
  für `CLAUDE_CODE_OAUTH_TOKEN` bedeutet das: Ob ein gültiges Token hinter der
  Egress-Restriktion mit der Default-Liste angenommen wird, ist nicht gemessen.
  Mit dem generierten Allowlist-Container erreichte die CLI die API (Antwort
  „Invalid API key“ auf einen absichtlich falschen Schlüssel), `example.com`
  blieb gesperrt.
  **Nicht aufgenommen** (optional, laut Dokumentation von Anthropic): Telemetrie und
  Fehlerberichte (`http-intake.logs.us5.datadoghq.com`,
  `browser-intake-us5-datadoghq.com`; gemessen wird `http-intake.logs.us5.datadoghq.com`
  angefragt, sobald ein Schlüssel verwendet wird; abschaltbar mit
  `DISABLE_TELEMETRY` und `DISABLE_ERROR_REPORTING`), MCP-Connectors
  (`mcp-proxy.anthropic.com`), Update und Plugins (`downloads.claude.ai`,
  `storage.googleapis.com`, `github.com`, `raw.githubusercontent.com`). `github.com` und
  `raw.githubusercontent.com` sind ohnehin in der Basisliste. Wer einen dieser Hosts
  braucht, trägt ihn in `devcontainer.sandbox.egress.allow` ein. Im Test lief die
  CLI mit gesperrtem Telemetrie-Host.

- **Grenzen der Messung:** Dev-Container-CLI 0.89.0, Colima, arm64, Basis-Image
  `mcr.microsoft.com/devcontainers/base:debian`, Claude Code 2.1.302, Feature `1.0.5`;
  der Pin `1` schwimmt, spätere Versionen können abweichen. Ein echter Login und
  eine echte Anfrage mit gültigem Schlüssel wurden nicht ausgeführt.

### 3.10 Ollama (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).
- **Fundstelle:** Installation: [`slice-v1-devcontainer-ollama-messung`](../messungen/slice-v1-devcontainer-ollama-messung/README.md)
  (Docker 29.5.2 auf Colima 0.10.3, `arm64`, 2026-10-18) und
  [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md)
  (nativer Linux-Host, Docker 29.9.1, 2026-10-19); CI-Lauf: [`slice-cr-lastenheft-messstaende`](../messungen/slice-cr-lastenheft-messstaende/README.md).
  Cloud-Modelle: [`slice-v1-ollama-cloud-messung`](../messungen/slice-v1-ollama-cloud-messung/README.md) (ohne Konto) und
  [`slice-v1-ollama-cloud-messung-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md)
  (Docker 29.9.2, Claude Code `2.1.310`, Node `v24.22.1`, Dienst des Hosts `0.40.10`).
- **Wortlaut:**

**Lieferkette, Auto-Start, Loopback, nur CPU**

- **Die Lieferkette hat eine Grenze:** Der Upstream veröffentlicht für das Archiv keine
  Signatur und keine Herkunftsbescheinigung (gemessen für das Release `0.40.6`). Die
  Prüfsumme stammt aus demselben Release. Sie schützt davor, dass eine Datei des Releases
  später ausgetauscht wird; sie schützt nicht vor einem Release, das von Anfang an
  schadhaft war. `u-boot generate devcontainer` weist das aus. Die Meldung ist keine
  Sicherheitslockerung des Profils und entfällt unter `--quiet`, ebenso die Meldung, dass
  das Feature ohne Sandbox-Profil keine Wirkung hat.
- **Kein Auto-Start:** Der Container startet den Server nicht; `devcontainer.json`
  trägt keinen Startbefehl dafür. Sie starten ihn im Container selbst (gemessen: bereit
  nach höchstens zwei Sekunden).
- **Nur Loopback:** `u-boot` setzt weder `OLLAMA_HOST` noch eine Portweiterleitung. Der
  Server lauscht dann nur auf `127.0.0.1:11434`; ein Aufruf über die Adresse des
  Containers scheitert. Das gilt für die gepinnte Version; setzen Sie `OLLAMA_HOST`
  selbst, ist das Ihre Entscheidung und nicht zugesagt.
- **Nur CPU:** Das Dockerfile entpackt die Verzeichnisse `cuda_*` und `vulkan` des
  Archivs nicht; der Server meldet `library=cpu` (gemessen auf `arm64` und `amd64`). Der Download des Archivs bleibt bei jedem
  Image-Build groß (Größe des Downloads, gemessen: 1,56 GB für `arm64`, 1,44 GB für `amd64`);
  das gebaute Image ist dagegen kleiner (Größe des Images, gemessen auf `amd64`: 1,01 GB mit
  dem Ollama-Block gegen 0,95 GB ohne ihn). Der Bau läuft ohne Egress-Beschränkung.

**Modell-Speicher und Egress**

- **Modell-Speicher:** Die Modelle liegen unter `/home/vscode/.ollama` in dem benannten
  Volume `<projekt>-ollama-${devcontainerId}` und überstehen das **Löschen des Containers**
  (ein neuer Container mit demselben Volume listet das Modell). Dass sie den Neubau des
  Images überstehen, ist nicht zugesagt; gemessen ist es auf Linux (`amd64`) in zwei Läufen: Mit
  der Dev-Container-CLI 0.89.0 (`devcontainer up --remove-existing-container --build-no-cache`)
  behielt das Volume seinen Namen, und das Modell war danach gelistet. Die Aktion eines Editors
  (etwa VS Code) für den Neuaufbau ist nicht gemessen. Das Volume enthält außerdem
  das Schlüsselpaar `id_ed25519`, das der Server beim ersten Start anlegt; geben Sie das
  Volume nicht weiter. Wofür Ollama den Schlüssel verwendet, ist nicht gemessen.
- **Egress:** Mit der Beschränkung erlaubt `u-boot` für `ollama` zwei Hosts:
  `registry.ollama.ai` (Manifeste) und `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com`
  (Modell-Dateien). Damit lief der Pull von `qwen2.5:0.5b` und `smollm2:135m`. **Nicht
  gemessen:** ob der Name des zweiten Hosts über Zeit, Region und weitere Modelle
  derselbe bleibt. Wechselt er, bricht der Pull ab (gemessen ohne diesen Host: nach 66
  Sekunden `max retries exceeded`); dann tragen Sie den neuen Namen in
  `devcontainer.sandbox.egress.allow` ein. Auf Linux (`amd64`) lief der Pull beider Modelle im
  ersten Lauf (8,7 und 10,5 Sekunden, an einem Netzstandort) über genau diese zwei Hosts, die dort
  mitgeschnitten wurden; im zweiten Lauf gelang der Pull ebenfalls, die Hosts sind dort nicht
  mitgeschnitten;
  `ollama.com`, `example.com` und `huggingface.co` ließen sich nicht auflösen. Für lokale Modelle ist
  `ollama.com` nicht aufgenommen: gesperrt meldet der Server beim Start zwei Warnzeilen, Start, Pull
  und Lauf funktionierten (gemessen). Für Cloud-Modelle siehe den folgenden Punkt.

**Cloud-Modelle: Messung und Zusage**

  Gemessen: Nach der Anmeldung endeten alle drei Aufrufe mit Exit-Code 0 und einer Ausgabe, mit
  zwei Modellen (`glm-5.3-flash:cloud` und `gpt-oss:120b-cloud`), in einem Lauf. Die Zusage gilt
  nur für diesen Aufbau: ein Container aus dem Image, das `u-boot` mit den Features `ollama`,
  `node` und `claude-code` im Sandbox-Profil erzeugt, **den Sie selbst starten** (mit der
  zusätzlichen Capability `NET_ADMIN`, ohne das Egress-Skript, ohne Workspace-Volume, mit einem
  Volume für die Modelle), Aufrufe ohne Terminal, auf einem Linux-Host (`amd64`) mit lokaler
  Docker-Engine. **Nicht zugesagt:** der Ablauf mit Anmeldung in dem Container, den
  `u-boot devcontainer up` aus der erzeugten Konfiguration startet (insbesondere im Sandbox-Profil
  ohne Egress-Restriktion und damit ohne `NET_ADMIN`), und ein Aufruf mit Anmeldung hinter der
  Egress-Restriktion.


**Cloud-Modelle: `ANTHROPIC_API_KEY`, Anmeldung und Geräteschlüssel, Rückbau**

  - **`ANTHROPIC_API_KEY`:** Die Zusage für `claude -p` gilt nur, wenn die Variable bei dem
    Aufruf nicht gesetzt ist; die Messung lief in einer Umgebung ohne sie. Das erzeugte
    Sandbox-Profil mit `claude-code` trägt die Referenz auf die Variable des Hosts (siehe
    [Claude Code](../user/benutzerhandbuch.md#claude-code-optional)). Haben Sie auf dem Host einen Schlüssel in dieser
    Variablen, reicht die Referenz ihn in den Container durch, und die Zusage gilt dann nicht
    (abgeleitet aus der Referenz, nicht mit einem Schlüssel gemessen; gemessen mit
    `u-boot devcontainer shell`, je ein Lauf auf Linux `amd64`: eine ungesetzte Host-Variable
    kommt ungesetzt an, eine leer gesetzte als leere, gesetzte Variable der Länge 0). Ob
    `unset ANTHROPIC_API_KEY` oder `env -u ANTHROPIC_API_KEY` im Aufruf die Bedingung
    herstellt, ist nicht gemessen. **Nicht gemessen:** der Aufruf bei gesetztem
    `ANTHROPIC_API_KEY` und bei einer leeren Variable im Container.
  - **Anmeldung und Geräteschlüssel:** Die Anmeldung führen Sie aus; `u-boot` reicht keine
    Umgebungsvariable des Hosts für Ollama durch und nimmt `OLLAMA_API_KEY` nicht auf. Der Server
    legt beim ersten Start das Schlüsselpaar `id_ed25519` und `id_ed25519.pub` im Volume der
    Modelle (`/home/vscode/.ollama`) an, **vor** jeder Anmeldung und nicht im Workspace-Volume;
    die Anmeldung änderte gegenüber diesem Zustand keine Datei (Pfad, Größe, Rechte; den Inhalt
    hat die Messung nie gelesen). Geben Sie das Volume nicht weiter. Den Rückbau auf der
    Schlüsselseite Ihres Kontos bei Ollama machen Sie selbst. Gemessen ist das im beschriebenen
    Aufbau und in einem Container, den die Dev-Container-CLI aus der erzeugten Konfiguration des
    Sandbox-Profils startet; für den Container, den `u-boot devcontainer up` ohne
    Egress-Restriktion startet, ist es nicht gemessen. **Nicht gemessen:** die Wirkung von
    `ollama signout` im Konto (der Befehl meldete `You have signed out of ollama.com`, Exit-Code
    0), die Rolle des Schlüssels bei den Cloud-Anfragen und der Zustand nach der Neuanlage des
    Containers.

    **Rückbau:** Auf der Seite Ihres Kontos bei Ollama, die Ihre Schlüssel führt, entfernen Sie den
    Eintrag selbst; der öffentliche Teil liegt als `id_ed25519.pub` im Volume, ob und wie diese
    Seite ihn anzeigt, ist nicht gemessen. Lokal entfernen Sie den Schlüssel mit dem Volume:
    `u-boot devcontainer down` entfernt weder den Container noch das Volume, und
    `docker volume rm` lehnt ein Volume ab, das ein Container noch benutzt, auch ein gestoppter
    (gemessen mit einem Container von `u-boot devcontainer up` nach `down`: `volume is in use`,
    die Meldung nennt den Container). Entfernen Sie den Container mit `docker rm`, dann das Volume
    mit `docker volume rm <projekt>-ollama-<Kennung>`; `docker volume ls` zeigt den vollständigen
    Namen. Mit dem Volume gehen die geladenen Modelle verloren.

**Cloud-Modelle: Egress, `OLLAMA_NO_CLOUD`, Modellnamen mit `[1m]`, Neuanlage**

  - **Hinter der Egress-Restriktion:** Die Liste führt `ollama.com` nicht; die Zeile `ollama`
    bleibt bei zwei Hosts. Gemessen ohne Konto in zwei Läufen je Fall, in einem Container, den die
    Dev-Container-CLI aus der erzeugten Konfiguration startet, mit einem protokollierenden
    Resolver: Ohne den Eintrag scheitert `ollama run` eines Cloud-Modells mit Exit-Code 1 und
    einem Fehler zur Namensauflösung (`dial tcp: lookup ollama.com on 127.0.0.1:53: server
    misbehaving`). Mit dem Eintrag `ollama.com` in `devcontainer.sandbox.egress.allow` (der Eintrag
    schließt Subdomains ein) wechselt die Meldung zur Aufforderung `You need to be signed in to
    Ollama to run Cloud models.`. So setzen und wenden Sie den Eintrag an (gemessen wurde der
    Neuaufbau mit der Dev-Container-CLI, nicht mit `u-boot devcontainer up --recreate`):

    ```bash
    u-boot config set devcontainer.sandbox.egress.allow ollama.com
    u-boot generate devcontainer
    u-boot devcontainer up --recreate   # legt einen vorhandenen Container neu an
    ```

    Mit Konto lief ein Aufruf mit diesem Eintrag, von Hand gesetzt, in einem Lauf; ob der Eintrag
    mit Konto **nötig** ist, folgt daraus nicht. **Nicht zugesagt:** `ollama launch claude` und
    `claude -p` hinter der Restriktion ohne Anmeldung; die Hosts, die der Server im Container zum
    Cloud-Dienst braucht.
  - **`OLLAMA_NO_CLOUD`:** `u-boot` setzt die Variable nicht. Mit `OLLAMA_NO_CLOUD=1` an Server
    **und** Client scheiterte `ollama run` mit `Error: ollama cloud is disabled: remote model
    details are unavailable` und `ollama launch claude` mit `Error: headless --yes launch for
    claude requires --model <model>` (davor `Warning: ignoring --model glm-5.3-flash:cloud because
    cloud is disabled`), jeweils mit Exit-Code 1 und ohne Namensanfrage im Container, trotz Anmeldung.
    **Nicht gemessen:** die Variable allein am Server oder allein am Client.
  - **Modellnamen mit `[1m]`:** `u-boot` wählt und verwaltet keine Modelle. Gemessen:
    `claude --model glm-5.3-flash:cloud[1m] -p` lieferte eine Antwort ohne Hinweis zum
    Kontextfenster; ohne das Suffix kam der Hinweis `Claude Code doesn't recognize the model
    "glm-5.3-flash:cloud", so it assumes a context window of 200k tokens`. `ollama run` und
    `ollama launch` lehnten den Namen mit `[1m]` ab (`Error: invalid model name`, Exit-Code 1,
    ohne Namensanfrage); gemessen mit einem Dienst des Hosts und im Container mit eigenem
    Server (dort nur bei gesetztem `OLLAMA_NO_CLOUD` beziehungsweise `OLLAMA_API_KEY`). **Nicht gemessen:** ob Claude Code das Suffix vor dem Senden entfernt.
  - **Neuanlage mit demselben Volume:** zugesagt ist nichts. Im Handaufbau der Messung
    (hinter der Egress-Restriktion, der Eintrag `ollama.com` von Hand gesetzt) kam nach der
    Neuanlage des Containers mit demselben Volume die Antwort ohne neue Anmeldung (ein Lauf);
    welcher Teil des Volume-Inhalts sie trägt, ist nicht gemessen.

**Grenzen der Messung (Cloud-Modelle) und der Installation**

  - **Grenzen der Messung (Cloud-Modelle):** Linux `amd64` (Ubuntu 24.04.5) mit lokaler
    Docker-Engine 29.9.2, Ollama-Client `0.40.6` (Pin des erzeugten Dockerfiles), Claude Code
    `2.1.310` und Node `v24.22.1` (nicht gepinnt, Stand beim Bau), die Modelle
    `glm-5.3-flash:cloud` und `gpt-oss:120b-cloud`, ohne Terminal, mit einem Konto des Projekts,
    am 2026-10-24, 40 Prompts. Ein bis drei Läufe je Aussage,
    ein Lauf zählt einen Container mit neuem Volume; die Läufe mit Konto liefen als Handvarianten
    (Aufruf der Dev-Container-CLI statt `u-boot devcontainer up`, `docker exec`, Platzhalter im
    Aufruf). **Nicht gemessen:** der interaktive Ablauf von Claude Code und von `ollama launch
    claude`, Kontingent und Kosten, `arm64`, Docker Desktop, Colima, Podman, rootless Docker, ein
    gültiger `OLLAMA_API_KEY`, das Sandbox-Profil, wie `u-boot` es ohne Egress-Restriktion
    erzeugt, andere Modelle und Versionen, abgelaufene oder entfernte Anmeldung.
- **Grenzen der Messung:** Zugesagt ist die Installation auf einem Linux-Host mit lokaler
  Docker-Engine (`amd64`) und unter Colima (`arm64`); gemessen sind zwei Umgebungen
  (auf dem `amd64`-Host zwei Läufe), Ollama `0.40.6`:
  Docker 29.5.2 auf Colima 0.10.3 (`arm64`) und, am 2026-10-19, Linux x86_64 (`amd64`) mit
  einer nativen Docker Engine 29.9.1 ohne VM (Kernel 6.8, `kernel.apparmor_restrict_unprivileged_userns=0`,
  `/dev/fuse` und `/dev/net/tun` vorhanden). Auf dem `amd64`-Host baute das erzeugte
  Dockerfile (Prüfsumme bestanden), der Server lauschte nur auf `127.0.0.1:11434` und meldete
  `library=cpu`, Pull und Lauf der beiden Modelle gelangen hinter der Beschränkung, die Modelle
  überstanden das Löschen des Containers und den Neuaufbau, und die Kombination mit
  `nestedRuntime: podman` und mit `nestedRuntime: docker` baute, ließ die Paketliste unverändert
  und führte einen kleinen verschachtelten Container aus, während der Server und ein Pull liefen. Das
  sind **zwei Messläufe auf einem Host**; in der CI bestand der Docker-Test der Installation
  (Bau des erzeugten Dockerfiles, Server auf der Loopback-Adresse, `library=cpu`, Pull, Modell
  überlebt das Löschen des Containers) am 2026-10-25 auf einem `amd64`-Runner mit Ubuntu 24.04;
  die Docker-Tests zum Bau mit `podman build` und zum Start von `podman` übersprang die CI, weil der Runner
  `kernel.apparmor_restrict_unprivileged_userns=1` hat (der Test zum falschen Pin mit `nestedRuntime: podman` bestand). Unter Colima (`arm64`) ist ein von Hand
  nachgeschriebenes Probe-Image gemessen, der Bau des erzeugten Dockerfiles nicht. Nicht gemessen: Docker Desktop, Linux hinter einer VM, größere
  Modelle und Mindest-Speicher, ein Start des Servers über die
  `devcontainer.json` (gemessen ist der Start von Hand und der Egress-Start durch die
  Dev-Container-CLI). GPU-Beschleunigung und die Verwaltung von Modellen durch `u-boot`
  gibt es nicht. Die Messung der Cloud-Modelle steht im Punkt *Cloud-Modelle von Ollama* oben,
  die des Anschlusses an einen Dienst des Hosts im folgenden Abschnitt.


### 3.11 Anschluss an einen Ollama-Dienst des Hosts (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion).
- **Fundstelle:** [`slice-v1-ollama-cloud-messung-mit-konto`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md)
  (`host.docker.internal`, Dienst des Hosts `0.40.10`, Docker 29.9.2).
- **Wortlaut:**

- **Relay und Bindung liegen bei Ihnen:** `u-boot` startet keinen Relay. Lauscht der Dienst
  nur auf der Loopback-Adresse des Hosts (gemessen: Linux, lokale Docker-Engine, Dienst auf
  `127.0.0.1:11434`), ist er ohne einen Relay oder eine geänderte Bindung durch Sie aus dem
  Container nicht erreichbar: Ohne die Zeile löst `host.docker.internal` im Container nicht
  auf; mit der Zeile und ohne Relay scheitert der Aufruf (drei Läufe, mit einer von Hand
  gesetzten Zeile). Für einen Dienst mit anderer Bindung gibt es keine Zusage, auch nicht für
  eine Bindung, die Sie ändern. Ein Verfahren für den Relay steht hier nicht.
- **Wirkung eines Relay:** Macht ein Relay den Dienst über die Adresse des Gateways
  (`host-gateway`) erreichbar, ist der angemeldete Dienst für **jeden Container im Netz
  `bridge`** erreichbar, solange er läuft (aus der Bindung abgeleitet, mit einem fremden
  Container **nicht geprüft**). Im Sandbox-Profil stellt das dem Agenten im Container das
  Kontingent Ihres Kontos zur Verfügung (Folgerung, nicht gemessen). Entscheiden Sie das
  vor dem Einsatz eines Relay.

  Gemessen: zwei Läufe im Standardprofil und zwei von zwei Läufen in einem Container aus dem
  Image des Sandbox-Profils, den Sie selbst starten (mit der zusätzlichen Capability
  `NET_ADMIN`, ohne das Egress-Skript), jeweils mit den Features `claude-code` und `node` und
  mit beiden Modellen. **In der Messung war die Zeile von Hand gesetzt**, nicht von `u-boot`
  erzeugt; die Wirkung der erzeugten gegenüber der von Hand gesetzten Zeile ist nicht zugesagt. Der
  Integrationstest von `u-boot` prüft nur, dass der Name im Container auflöst, nicht dass der
  Dienst antwortet.

- **Hinter der Egress-Restriktion trägt der Anschluss nicht:** `host.docker.internal` löst
  über `/etc/hosts` auf, aber ein Aufruf von `claude -p` läuft in die Zeitgrenze (gemessen:
  Exit-Code 124 nach 120 Sekunden); ein Eintrag `host.docker.internal` in `devcontainer.sandbox.egress.allow`
  ändert das nicht. `u-boot` trägt keine Adresse ein und warnt beim Erzeugen, ohne
  abzubrechen. **Nicht gemessen:** ein Weg, der die Adresse erlaubt, und `ollama run` und
  `ollama launch` mit dem Dienst des Hosts hinter der Restriktion.
- **Grenzen der Messung:** Linux `amd64` (Ubuntu 24.04.5) mit lokaler Docker-Engine 29.9.2,
  Ollama-Client `0.40.6`, Dienst des Hosts `0.40.10` auf `127.0.0.1:11434`, Claude Code
  `2.1.310`, die Modelle `glm-5.3-flash:cloud` und `gpt-oss:120b-cloud`, ohne Terminal, mit
  einem Konto des Projekts, am 2026-10-24. Ein bis drei Läufe je Aussage, als Handvarianten
  (die Zeile von Hand gesetzt, der Relay ein Hilfs-Container, Platzhalter im Aufruf).
  **Nicht gemessen:** Docker Desktop, Colima, Podman, rootless Docker, `arm64`, Windows, ein
  Dienst mit anderer Bindung, die Erreichbarkeit durch einen fremden Container, das
  Sandbox-Profil, wie `u-boot` es ohne Egress-Restriktion erzeugt, der interaktive Ablauf von
  Claude Code, Kontingent und Kosten, ein Aufruf bei gesetztem `ANTHROPIC_API_KEY`.


### 3.12 Forge-CLIs (Handbuch, Abschnitt 4.11)

- **Zusage:** [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer).
- **Fundstelle:** [`slice-v1-forge-cli-quellen-messung`](../messungen/slice-v1-forge-cli-quellen-messung/README.md)
  (Docker 29.9.2, Dev-Container-CLI 0.89.0, `gh` 2.103.0, `glab` 1.124.0).
- **Wortlaut:**

**Anmeldung in kurzer Form.** `u-boot` meldet nichts an und speichert kein Token; die
Werkzeuge lesen es aus einer Umgebungsvariable, `GH_TOKEN` für `gh` und `GITLAB_TOKEN` für
`glab`. Im Sandbox-Profil trägt `devcontainer.json` die Referenz auf die Variable; setzen
Sie sie auf dem Host, bevor Sie `u-boot devcontainer shell` aufrufen, dann gibt `shell` sie
weiter. Im Standardprofil steht keine Token-Referenz in `devcontainer.json`, und `shell` gibt die
Variablen nicht weiter; dort geben Sie die Variable mit `docker exec -e GH_TOKEN -u <Benutzer>
<Container> …` mit. Gemessen mit einem Platzhalter ist nur, dass die Variable dann im
Container ankommt; eine Anfrage von `gh` oder `glab` auf diesem Weg ist nicht gemessen.

Ohne Variable endet `gh api` mit Exit-Code 4, `glab` liest ein öffentliches Projekt weiter;
mit einem Platzhalter antworten die Werkzeuge mit `401`, ein Aufruf mit einem gültigen
Token ist nicht gemessen. `gh auth login` und `glab auth login` im Container sind nicht
geprüft. Die Spiegelstriche unten tragen die Belege und Messvorbehalte.

- **Die Lieferkette hat eine Grenze:** Die Prüfsumme stammt aus demselben Release wie das
  Archiv. Sie schützt davor, dass eine Datei des Releases später ausgetauscht wird; sie
  schützt nicht vor einem Release, das von Anfang an schadhaft war. Für beide Archive hat
  die Messung keine Signatur gefunden (nach dem Suchmuster der Messung, nicht
  abschließend). Für `gh` veröffentlicht das Projekt eine GitHub-Attestierung; `u-boot`
  prüft sie nicht. `u-boot generate devcontainer` weist das aus, je Feature eine Meldung in
  jedem Profil. Die Meldung ist keine Sicherheitslockerung des Profils und entfällt unter
  `--quiet`; die Blöcke im Dockerfile bleiben.
- **Token-Variablen:** Für `gh` gilt `GH_TOKEN`, für `glab` `GITLAB_TOKEN`. Im
  Sandbox-Profil trägt `devcontainer.json` bei aktivem Feature die Referenz
  `"GH_TOKEN": "${localEnv:GH_TOKEN}"` beziehungsweise
  `"GITLAB_TOKEN": "${localEnv:GITLAB_TOKEN}"` neben `GIT_TOKEN`, nie einen Wert. Setzen
  Sie die Variable in der Host-Umgebung, bevor Sie den Container starten. `u-boot` führt
  weder `gh auth login` noch `glab auth login` aus und speichert kein Token; nutzen Sie
  kurzlebige, auf das Repository begrenzte Token (siehe
  [Zugangsdaten](../user/benutzerhandbuch.md#zugangsdaten-git-und-agent)). Die Werkzeuge lesen auch `GITHUB_TOKEN`,
  `GITLAB_ACCESS_TOKEN` und `OAUTH_TOKEN` (gemessen mit einem Platzhalter: sie versuchen
  die Anfrage mit dem Wert); `u-boot` reicht nur `GH_TOKEN` und `GITLAB_TOKEN` durch.
  Im Sandbox-Profil öffnet `u-boot devcontainer shell` die Shell mit den beiden Variablen,
  soweit sie auf dem Host gesetzt sind; eine auf dem Host ungesetzte Variable fehlt in der
  Sitzung. Im Standardprofil gibt `shell` sie nicht weiter (siehe unten).
- **Ohne Variable** (gemessen mit einer in der Umgebung des Werkzeugs nicht gesetzten
  Variable): `gh api` endet mit Exit-Code 4 und dem Hinweis, `gh auth login` auszuführen oder
  `GH_TOKEN` zu setzen; `gh auth status` endet mit Exit-Code 1. `glab` liest ein
  öffentliches Projekt ohne Token (`glab api`, `glab repo view`, `glab release list`,
  `glab release download`: Exit-Code 0); `glab auth status` endet mit Exit-Code 1, und
  `glab repo clone` scheitert (`401 Unauthorized`). **Nicht gemessen:** eine Variable, die auf
  dem Host ungesetzt ist und im Container als leere Zeichenkette ankommt.
- **Mit Variable, ohne Anmeldung durch `u-boot`:** Das Werkzeug versucht die Anfrage mit
  dem Token aus der Umgebung. Gemessen mit einem **Platzhalter** statt eines gültigen
  Tokens: Die Antwort ist `401`, und bei `glab` scheitern damit auch `glab api` und `glab repo view`
  gegen ein öffentliches Projekt, die ohne Variable gelingen (`glab release list` und
  `glab release download` sind mit Variable nicht gemessen). `glab auth status` meldet, dass es ein
  Token in der Umgebungsvariable gefunden hat, mit maskiertem Wert. Ein Aufruf mit einem
  **gültigen** Token ist nicht gemessen und nicht zugesagt.
- **Standardprofil, keine Token-Referenz:** Das Standardprofil erzeugt kein `remoteEnv`
  und damit keine Token-Referenz; `u-boot doctor` prüft die Variable dort nicht. Das
  Feature wirkt trotzdem: Das Programm liegt im Image. Die gemessene Quelle des Tokens ist
  die Umgebungsvariable des Werkzeugs. `u-boot devcontainer shell` gibt nur Variablen
  weiter, auf die `devcontainer.json` verweist; im Standardprofil also weder `GH_TOKEN` noch
  `GITLAB_TOKEN`. Der gemessene Weg ist `docker exec -e GH_TOKEN -u <Benutzer> <Container> …`
  (ebenso `-e GITLAB_TOKEN`): Ist die Variable auf dem Host gesetzt, kommt sie im Aufruf an,
  ohne `-e` fehlt sie, und ist sie auf dem Host ungesetzt, fehlt sie auch im Container
  (gemessen mit einem Platzhalter statt eines Tokens, Linux `amd64`, lokale Docker Engine
  29.9.2, ein Lauf, in einem Container ohne die Forge-CLIs; die Antwort der Werkzeuge auf einen
  Platzhalter steht im vorigen Spiegelstrich). Ein Aufruf von `gh` oder `glab` mit einem
  gültigen Token auf diesem Weg ist nicht gemessen und nicht zugesagt. `gh auth login` und
  `glab auth login` im Container sind nicht geprüft: Kein Aufruf mit einem Token, und ob
  ein dabei abgelegtes Token das Neuanlegen des Containers übersteht, ist nicht geprüft.
  Das Sandbox-Profil reicht die Variable des Hosts durch.

- **Git-Protokoll von `glab`:** Das Standardprotokoll für Git ist `ssh` (gemessen:
  `glab config get git_protocol` ergibt `ssh`). Im Container liegt kein Schlüssel, und nur
  `https`-Remotes bekommen Zugangsdaten (siehe [Zugangsdaten](../user/benutzerhandbuch.md#zugangsdaten-git-und-agent)).
  `glab config set git_protocol https --global` stellt auf `https` um (gemessen: Exit-Code
  0; der Clone-Weg danach ist nicht gemessen). `gh` führt `https` als Standard
  (`gh config list`).
- **`gh auth setup-git`:** Führen Sie den Befehl im Container selbst aus (gemessen mit
  `--hostname github.com` und einem Platzhalter in `GH_TOKEN`, ein Lauf), schreibt er
  `~/.gitconfig` mit je einem leeren und einem `gh auth git-credential`-Eintrag für
  `github.com` und `gist.github.com`, ohne ein Token im Inhalt. Der Eintrag
  setzt für diese Hosts `gh` an die Stelle des Credential-Helpers des Images, der
  `GIT_TOKEN` liefert; die Wirkung auf einen Git-Aufruf ist nicht gemessen.
- **Telemetrie:** `glab` führt `telemetry` und `check_update` in
  `~/.config/glab-cli/config.yml` (gemessen: beide `true`); `gh` nennt in `gh help
  environment` die Variablen `GH_TELEMETRY` und `DO_NOT_TRACK`. `u-boot` setzt keine
  Umgebungsvariable dafür: Das erzeugte Dockerfile hat kein `ENV`, und `remoteEnv` trägt
  nur Referenzen auf Variablen des Hosts. **Nicht gemessen:** die Wirkung der Schalter, und welche Hosts
  Telemetrie und Aktualisierungsprüfung in einer interaktiven Sitzung ansprechen (die
  Aufrufe der Messung liefen ohne Terminal).
- **Egress:** Mit der Beschränkung erlaubt `u-boot` `api.github.com` und `github.com`
  (beide stehen schon in der Basisliste) für `gh` und `gitlab.com` für `glab`. **Nicht
  aufgenommen:** Hosts, die auf eine Anmeldung folgen, Hosts von Telemetrie und
  Aktualisierung sowie die Hosts von `glab duo` und `glab orbit`; wer sie braucht, trägt
  sie in `devcontainer.sandbox.egress.allow` ein. Ohne `gitlab.com` in der Liste scheitert
  `glab` an der Namensauflösung (gemessen: `lookup gitlab.com`), und der Client wiederholt
  die Anfragen: Der Containerlauf der Messung dauerte 127,8 statt 20,0 Sekunden mit dem Host
  in der Liste (ganzer Lauf samt Start, ein Host, ein Lauf, keine Dauer eines einzelnen Aufrufs).
- **Nicht zugesagt:** Aufrufe mit gültigem Token und deren Hosts; Telemetrie- und
  Update-Hosts einer interaktiven Sitzung; GitHub Enterprise und GitLab Self-Managed;
  andere Versionen als die gepinnten; `glab duo`, `glab orbit` und das Nachladen weiterer
  Programme durch `glab`; ein Bau auf `arm64` (die Prüfsummen für `arm64` stehen im
  Dockerfile, der Bau ist nicht gelaufen); die Kombination mit `ollama` und
  `nestedRuntime` über die Reihenfolge der Blöcke im Dockerfile, die Volumes und die Hosts
  hinaus.
- **Grenzen der Messung:** Linux `amd64` (x86_64), lokale Docker Engine 29.9.2,
  Dev-Container-CLI 0.89.0, Basis-Image `mcr.microsoft.com/devcontainers/base:debian`,
  `gh` 2.103.0, `glab` 1.124.0, gemessen am 2026-10-23, ein gewerteter Lauf je Aussage.
  Es wurde kein Token angelegt oder verwendet, die Aufrufe liefen ohne Terminal. Die
  Aufrufe der Werkzeuge liefen in Containern, deren Dockerfile den Block von Hand trug
  (dasselbe Archiv wie im erzeugten Block); am von `u-boot` erzeugten Dockerfile laufen in
  einem Docker-Test der Bau in beiden Profilen, die Versionen, ein Aufruf von `glab` gegen
  ein öffentliches Projekt hinter der Egress-Beschränkung und die Token-Referenzen.

### 3.13 Grenzen der verschachtelten Laufzeit (Handbuch, Abschnitt 4.11, „Grenzen“)

- **Zusage:** [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).
- **Fundstelle:** wie in Abschnitt 3.4; die Läufe vom 2026-10-19 (Linux x86_64, `amd64`, native Docker Engine 29.9.1):
  [`slice-v1-devcontainer-ollama-amd64-messung`](../messungen/slice-v1-devcontainer-ollama-amd64-messung/README.md).
- **Wortlaut:**

- Geprüft unter Docker auf Linux mit der `devcontainer`-CLI; `nestedRuntime: docker`
  zusätzlich unter Colima (Linux-VM); nur diese beiden Hosts sagt das Lastenheft zu.
  **Nicht geprüft und nicht zugesagt:** Docker Desktop,
  Podman als Host-Engine, VS Code und Codespaces.
- `nestedRuntime: podman`: zugesagt unter Docker auf einem Linux-Host mit lokaler Engine
  (`amd64`) und unter Colima (`arm64`). Podman 6.1.4 und `podman build` sind unter Docker auf
  Colima (arm64) mit erlaubten verschachtelten User-Namespaces gemessen; für
  amd64 sind Download, Prüfsumme, Signatur und Versionsabfrage unter Emulation
  geprüft. Ein verschachtelter Lauf (`podman run` eines kleinen Images, ohne und
  hinter der Egress-Beschränkung, mit einem Pull hinter ihr) ist auf Linux x86_64
  (`amd64`, native Docker Engine 29.9.1, 2026-10-19) in zwei Läufen gemessen, zusammen mit dem
  Feature `ollama`; `podman build` auf `amd64` hinter der Egress-Beschränkung läuft durch, braucht aber
  bei einem Namen außerhalb der Liste etwa 36 Sekunden je Name länger
  (siehe *Namensauflösung und ihre Grenze* im Abschnitt zur Egress-Beschränkung); die CI übersprang am 2026-10-25 die Docker-Tests zum Bau mit `podman build` und zum Start von `podman`, weil der Runner `kernel.apparmor_restrict_unprivileged_userns=1` hat (der Test zum falschen Pin mit `nestedRuntime: podman` bestand).
  Podman 6 unterstützt nur cgroups v2 (laut Release-Notes, nicht gemessen);
  Hosts mit cgroups v1 sind davon betroffen.
  Unter einer Ubuntu-24.04-VM mit dem Schalter
  `kernel.apparmor_restrict_unprivileged_userns=1` startet verschachteltes
  Podman nicht (Exit-Code 11). Nicht verfügbar sind `buildah`, `docker buildx`,
  `docker compose` und BuildKit-Frontends.

### 3.14 README.md und README.de.md

- **Zusage:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Abschnitt 3.1).
- **Wortlaut** (Status-Block und Abschnitt „Was macht es vertrauenswürdig?“; Wörter „Zusagen“, „zusagt“, „promises“, „promised“):
  - README.md: „Promises, conditions and limits are in the user manual“.
  - README.md: „What `u-boot` promises is specified in the Lastenheft (German).“
  - README.de.md: „Zusagen, Bedingungen und Grenzen stehen im Benutzerhandbuch“.
  - README.de.md: „Was `u-boot` zusagt, steht im Lastenheft.“

### 3.15 Messbedingungen der Spezifikation ([`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes))

Die Spezifikation ([`spezifikation.md`](../../spec/spezifikation.md)) beschreibt das Zielbild und trägt keine Messbedingungen.
Die Sätze unten sind wörtlich übernommen. Die drei Stellen, die der Projektinhaber am 2026-10-31 entschieden hat, sind in der
Spezifikation umgesetzt (Wortlaut vorher und nachher am Ende des Abschnitts). Die Verhaltensaussagen und die Grenzen der
Zusage bleiben in der Spezifikation.

- **Zusage:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
- **Zusage-Menge der Engines (Entscheidung des Projektinhabers, 2026-11-01):** Zugesagt ist nur die Engine 29.9.2 mit dem
  Client 29.9.2. Die Engine 28.0.4 ist nicht zugesagt; die Form „Zeile ohne Präfix auf der Standardausgabe, Status 126“
  stammt aus einer früheren Messung ohne Protokoll unter `docs/messungen/`, die
  [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md)
  nicht bestätigen konnte (keine Engine 28.0.4 verfügbar). Die Erkennung dieser Form steht als Bestandscode im Launcher
  (`userMissingOnStdout`), von Einheitstests mit Fake gepinnt. Liefert der GitHub-Runner weiter die Engine 28.0.4
  (`slice-v1-devcontainer-shell-fundament`, V-16: im README zweier Runner-Images gelesen), wäre eine Nachmessung dort
  möglich; sie ist nicht beschlossen und die Engine bleibt bis dahin nicht zugesagt.
- **Fundstelle:** Die Aussagen zu Sitzungsende, Signalen und Art der Shell stammen aus dem
  Planungsdokument [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md)
  (Docker 29.9.2, Basis-Image `devcontainers/base:debian`, Test unter `script`; die Runner-Engine 28.0.4 steht dort ebenfalls und gilt nicht als Zusage). Die Lebenszyklus-Aussage
  stammt aus dem Test `TestDevcontainerUp_Docker_AlreadyRunningLifecycleCounters`
  (`internal/e2e/devcontainer_start_docker_test.go`). **Für die Shell-Aussagen gibt es kein Messprotokoll unter
  `docs/messungen/`.** Das Protokoll [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md)
  (Docker 29.9.2, Dev-Container-CLI 0.89.0) trägt sie nicht; es misst Namen, Abbruch und Hänger.
- **Wortlaut, Aufzählungspunkt „Messbedingungen“ (Aussagen zum Ende der Sitzung, zu den Signalen und zur Art der Shell):**
  „**Messbedingungen** der Aussagen zum Ende der Sitzung, zu den Signalen und zur Art der Shell: Docker 29.9.2 als Engine und
  Client, Linux `amd64`, lokale Engine; der Container ist ein `sleep 3600` des Basis-Images
  `mcr.microsoft.com/devcontainers/base:debian`, nicht ein von `devcontainer up` erzeugter; das Terminal ist ein
  Pseudo-Terminal von `script`; ein Lauf je Fall. Andere Engines, Hosts und Terminals sind nicht zugesagt.“
- **Wortlaut, Absatz „Lebenszyklus-Befehle bei `up` gegen einen laufenden Container“ (zweiter Satz):**
  „Gemessen ist das mit der gepinnten Version der CLI, mit Befehlen in der Form einer Zeichenkette, die gelingen, unter Linux
  `amd64` mit lokaler Engine (Docker 29.9.2), ein Lauf je Aufruf.“
- **Wortlaut, Aufzählungspunkt „Grenzen der Erkennung“ (Ende des Punkts „Engine-Abhängigkeit“):**
  „Messbedingungen: zwei Engines (29.9.2 und 28.0.4), Client 29.9.2, Linux `amd64`, ein Lauf je Fall; andere Engine-Versionen
  sind nicht zugesagt, und die Meldungen sind für den Client 29.9.2 gemessen.“ Stand seit 2026-11-01: In der Spezifikation
  steht „Zugesagt ist die Engine 29.9.2 mit dem Client 29.9.2 unter Linux `amd64`“; die Engine 28.0.4 steht unter „nicht zugesagt“.
  Ungepinnt bleibt der Fall „die Sitzung schreibt die Zeile selbst“ (Lesung des Codes, kein Test), er ist in der Spezifikation als
  solche gekennzeichnet.
- **Wortlaut, Aufzählungspunkt „Grenzen der Erkennung“ (Zeile zum fehlenden Benutzer auf der Engine 28.0.4):**
  „(aus der Form abgeleitet, nicht gefahren)“ — nach der Aussage, der Befehl endete dann mit `12` statt `0`.
- **Umgesetzt nach Entscheidung des Projektinhabers (2026-10-31).** Die drei Stellen stehen so in der Spec; der Wortlaut
  vorher steht oben.
  - Aufzählungspunkt „Messbedingungen“ (Sitzungsende, Signale, Art der Shell): jetzt „**Reichweite** der Aussagen zum Ende der
    Sitzung, zu den Signalen und zur Art der Shell: Sie sind zugesagt für Container mit dem Basis-Image
    `mcr.microsoft.com/devcontainers/base:debian`, deren Hauptprozess ein Dauerlauf ist; für einen von `devcontainer up`
    erzeugten Container sind sie nicht zugesagt, für andere Basis-Images ebenfalls nicht. Andere Engines, Hosts und
    Terminals sind nicht zugesagt.“ Grund der Nachschärfung: Die Signal- und Shell-Tests
    (`internal/adapter/driven/docker/devcontainer_shell_docker_test.go`, `newShellContainer`) starten einen einfachen Container
    mit `docker run -d <Image> sleep 3600`; kein Test läuft gegen einen von `devcontainer up` erzeugten Container. Fassung vor
    der Nachschärfung: „… zugesagt für Container mit dem Basis-Image `mcr.microsoft.com/devcontainers/base:debian`; für andere
    Basis-Images sind sie nicht zugesagt. Andere Engines, Hosts und Terminals sind nicht zugesagt.“ Das Basis-Image steht im Planungsdokument `slice-v1-devcontainer-shell-fundament` (Zeile
    „Messung vor dem Bau“: Basis-Images `mcr.microsoft.com/devcontainers/base:debian` und `busybox:1.36`; welche Aussage auf
    welchem Image lief, nennt die Zeile nicht). Gestrichen sind Docker 29.9.2, `sleep 3600`, das Pseudo-Terminal von `script` und „ein Lauf je
    Fall“.
  - Lebenszyklus-Absatz: jetzt „Das gilt für die gepinnte Version der CLI, für Befehle in der Form einer Zeichenkette, die
    gelingen, unter Linux `amd64` mit lokaler Engine.“ Der Absatz endet weiter mit „… und andere Hosts und Engines“ in der
    Aufzählung des Nicht-Zugesagten. Fundstelle der
    Zusage-Menge (Linux-Host, lokale Engine, `amd64`): [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
  - Punkt „Grenzen der Erkennung“: der Satz „Diese Zeile ist auch das, was ein `docker exec -u <Benutzer der Sitzung> <anderer
    Container>` innerhalb der Sitzung gegen eine Engine dieser Form als einzige Ausgabe erzeugen würde; der Launcher trennt
    beides nicht, der Befehl endete dann mit `12` statt `0` (aus der Form abgeleitet, nicht gefahren).“ ist aus der Spec
    gestrichen. Er bleibt hier als **unbelegte Beobachtung**: Kein Test pinnt ihn (ein Einheitstest in
    `devcontainer_shell_test.go` prüft die Zeile eines anderen Benutzers, nicht den Fall innerhalb der Sitzung).
    Stand nach `slice-v1-shell-erkennung-fehlender-benutzer-eng`: Ein `docker exec -u <Benutzer> <anderer Container>` gegen eine
    Engine 28.0.4 ist nicht gefahren. Auf der Engine 29.9.2 ist er gefahren
    ([Messprotokoll](../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md)): Die Zeile hat dort das Präfix
    `Error response from daemon: `, steht auf der Fehlerausgabe und endet mit Status 1. Die Spec führt den Fall jetzt als Grenze
    der Erkennung (Lesung des Codes, kein gefahrener Lauf des Launchers); die Einstufung dieser Zeile pinnt der Fall „the user
    does not exist after the probe“ in `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure`, der Fall der Sitzung, die
    die Zeile selbst schreibt, ist durch keinen Test gepinnt.

**Historie von `spec/spezifikation.md`: gestrichene Messaussagen im Wortlaut vor der Änderung** (Fundstelle: `git show
HEAD:spec/spezifikation.md`, Tabellenzeilen 1474 bis 1477, alle vom 2026-10-19):

- Zeile 1474 (`SPEC-046`): „… und Abgrenzung zum laufenden Container (ein Neustart behält die alte Liste, gemessen in einem
  Lauf).“ Gestrichen: „, gemessen in einem Lauf“.
- Zeile 1475 (`SPEC-019`): „… und die Subdomain-Aussage mit Messstand: gepinnte und gemessene Fälle getrennt, der übergeordnete
  Name bei geteilter IP-Adresse als nicht geprüft ausgewiesen, der offene Upstream-Resolver ohne `sudo` als gemessen; …“ Die
  Zeile steht in der Spec unverändert im alten Wortlaut (F-5 des Reviews).
- Zeile 1476 (`SPEC-047`): „… `SPEC-019` die Zeile `ollama` (zwei Hosts, Stabilität des Blob-Hosts ungemessen), …“ Gestrichen:
  „, Stabilität des Blob-Hosts ungemessen“.
- Zeile 1477 (`SPEC-045`): „… dass ein Prozess ohne `sudo` die Statusdatei ersetzen kann (Verzeichnis gehört `vscode`, gemessen
  in einem Lauf).“ Gestrichen: „, gemessen in einem Lauf“.

## 4. Was das Handbuch als Voraussetzung führt

Die Aussage ist hier als Prüflücke geführt; das Handbuch nennt nur die Voraussetzung, die Nutzer erfüllen müssen.

| Handbuch-Stelle | Voraussetzung im Handbuch | Wortlaut der Aussage |
|---|---|---|
| 4.9, Ende des Abschnitts „Den Devcontainer ohne Editor starten und stoppen“ | Der Befehl braucht Docker mit einer lokalen Engine (sonst Exit-Code 11). „Linux-Host“ und `amd64` sind die zugesagte Menge des Lastenhefts und keine Voraussetzung, die das Handbuch aufstellt: ob der Befehl auf `arm64` oder auf anderen Hosts läuft, ist nicht gemessen. | Abschnitt 3.2, „Features und Reichweite“ |
| 4.11, Podman, „Voraussetzungen und Verhalten des Baus“ | `podman build` holt Images von Docker Hub und setzt verschachtelte User-Namespaces voraus. | Abschnitt 3.4, „Gemessen und zugesagt“ |
| 4.11, „Grenzen“ | Die verschachtelten Laufzeiten setzen Docker als Engine auf dem Container-Host voraus. | Abschnitt 3.13, erster Punkt |
| 4.11, Ollama, Cloud-Modelle | Der Ablauf setzt einen selbst gestarteten Container aus dem Sandbox-Image mit `NET_ADMIN`, ohne Egress-Skript, ohne Workspace-Volume, auf Linux `amd64` mit lokaler Docker-Engine voraus; `claude -p` setzt voraus, dass `ANTHROPIC_API_KEY` nicht gesetzt ist. | Abschnitt 3.10, „Cloud-Modelle: Messung und Zusage“ |
| 4.11, Ollama, verschachtelte Laufzeit | `nestedRuntime: podman` setzt erlaubte verschachtelte User-Namespaces sowie `/dev/fuse` und `/dev/net/tun` voraus. | Abschnitt 3.10, „Grenzen der Messung“ |
| 4.11, Anschluss an einen Ollama-Dienst des Hosts | Ein Ollama-Dienst auf dem Host (`127.0.0.1:11434`), ein Relay als Hilfs-Container und `ANTHROPIC_API_KEY` nicht gesetzt beim Aufruf von `claude -p`. | Abschnitt 3.11 |
| 4.11, Claude Code | Der Eintrag legt nur die Hauptversion `1` des Features fest. | Abschnitt 3.9 |
