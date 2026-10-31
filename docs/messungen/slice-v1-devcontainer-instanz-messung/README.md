# Messprotokoll: mehrere Devcontainer-Instanzen desselben Projektordners (zwei, drei und vier Instanzen)

Messung zum Slice [`slice-v1-devcontainer-instanz-messung`](../../plan/planning/in-progress/slice-v1-devcontainer-instanz-messung.md) ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Gemessen ist, was die Dev-Container-CLI des Hilfs-Images tut, wenn N gleichzeitige Instanzen desselben Projektordners über ein zweites Id-Label unterschieden werden. Es lief **nicht** `u-boot devcontainer up`: Das Skript spiegelt den `docker run` von `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`) und hängt `--id-label` an.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-10-31 (Läufe 14:53 bis 16:02 UTC) |
| Host | Linux 7.0.0-42-generic, `x86_64`, Ubuntu 26.04.1 LTS, 8 CPUs, 15,24 GiB RAM, 16 GiB Swap (vor allen Läufen ungenutzt), `/` 301 GB mit rund 217 GB frei (Docker-Wurzel `/var/lib/docker` liegt auf `/`) |
| Docker | Client und Server 29.9.2 (`linux/amd64`), Engine lokal (Context `default`), Storage `overlayfs` (containerd-Snapshotter), cgroup v2 (systemd), `MemoryLimit` und `SwapLimit` wahr |
| Dev-Container-CLI | `0.89.0` (`devcontainer --version` im Mess-Image `ubmess-devcontainer-cli:0.89.0`; Kontext = `Dockerfile`, `package.json`, `package-lock.json` aus `internal/adapter/driven/docker/devcontainercli/`) |
| Hilfs-Image | Mess-Image `ubmess-devcontainer-cli:0.89.0` (ID je Lauf in den Rohdaten, z. B. `c87d21426419`); das Bestands-Image `u-boot-devcontainer-cli:0.89.0-f1a94cce795e` hat die ID `ed569ec37d56`, die Mess-Image-IDs weichen ab (anderer Tag, gleicher Kontext; die ID-Gleichheit ist nicht geprüft) |
| `u-boot` | `u-boot version 0.16.1-dev` (`bin/u-boot-linux-amd64`; erzeugt nur die Projektordner) |
| `kernel.apparmor_restrict_unprivileged_userns` | `1` |
| Profil | Sandbox mit `nestedRuntime: docker` und Feature `ollama` (Mounts: Workspace, Docker-Speicher, Modelle); zweite Probe im Standardprofil |
| Messbedingung je Instanz | `runArgs` der `devcontainer.json` im Wegwerf-Ordner um `--memory=3g --memory-swap=3g --cpus=1` ergänzt; **zurückgelesen** an jeder Instanz: `HostConfig.Memory=3221225472`, `MemorySwap=3221225472`, `NanoCpus=1000000000`. Die Annahme „je Instanz“ wurde auftragsgemäß so gesetzt (Projektinhaber-Vorgabe 2026-10-31) |
| Lockerung | Der Hilfs-Container bekommt den Docker-Socket des Hosts (wie `runArgs`). Kein Devcontainer trägt eine Einhängung `docker.sock` (in jeder Instanz geprüft) |
| Durchführung | Implementer, Skript unten |

## Skript und Aufruf

[`tools/messung/devcontainer-instanz/messen.sh`](../../../tools/messung/devcontainer-instanz/messen.sh), aus der Repo-Wurzel:

```bash
bash tools/messung/devcontainer-instanz/messen.sh <2|3|4|namen> bin/u-boot-linux-amd64 \
  > docs/messungen/slice-v1-devcontainer-instanz-messung/rohdaten-stufe<N>-lauf<L>.txt 2>&1
```

Ablauf je Stufe: M0, Vorprüfung, Mess-Image bauen, Sandbox-Phase (N Instanzen nacheinander starten, M1 bis M5, Inhaltsprobe, M6, M7, M2-Bezug), Standard-Phase (Bind-Mount, M6 umgekehrt), bei N = 2 zusätzlich die Sandbox-Phase ohne Limit; Aufräumen per `trap`, Inventur nachher. Alle Objekte tragen `ubmess-`/`vsc-ubmess-`; die CLI legt ihre Container mit **zufälligen Namen** an (`bold_lehmann`), sie sind nur am Label `devcontainer.local_folder=<Messordner>` zu erkennen; das Aufräumen sucht darum nach diesem Label, nicht nach dem Präfix.

Zeitlimits (nach dem Hänger im Lauf `namen`, siehe unten): je `devcontainer up` 300 s (`namen`: 90 s) mit `timeout -k 10` und anschließendem `docker rm -f` des Hilfs-Containers, je Lauf 900 s. Gewählt aus den Läufen: ein `up` dauerte 1 bis 7 s, ein Stufenlauf 67 bis 83 s (Stufe 4: rund 80 s). Ein Abbruch durch ein Limit steht als `BEFUND` in den Rohdaten, die Stufe gilt dann nicht als vollständig.

## Läufe und Rohdaten

| Datei | Inhalt |
|---|---|
| `rohdaten-stufe2-lauf1.txt`, `-lauf2.txt` | Stufe 2, vollständig, 14:56 und 14:57 UTC |
| `rohdaten-stufe3-lauf1.txt`, `-lauf2.txt` | Stufe 3, vollständig, 14:59 und 15:00 UTC |
| `rohdaten-stufe4-lauf1.txt`, `-lauf2.txt` | Stufe 4, vollständig, 15:41 und 15:43 UTC, **mit gelockerter Vorprüfung** (siehe Abweichung) |
| `rohdaten-stufe4-alte-schwelle-lauf1.txt`, `-lauf2.txt` | Stufe 4 mit der Schwelle des Slice: Abbruch in der Vorprüfung (14:53, 15:01 UTC), nichts angelegt |
| `rohdaten-stufe2-vorlauf-skriptfehler.txt` | Vorlauf mit Skriptfehler (M6 entfernte den Container von Instanz 2, M7 deshalb unvollständig); M1 bis M5 und M6 sind darin gültig, M7 nicht. Zählt nicht als Lauf |
| `rohdaten-namen-lauf1.txt`, `-lauf2.txt` | M8, 15:56 und 15:59 UTC |
| `rohdaten-namen-abbruch-haenger-lauf1.txt`, `-lauf2.txt`, `-lauf3.txt` | Drei Läufe von M8 mit der ersten Skriptfassung: jedes Mal Hänger bei dem Namen `a=b`, von mir abgebrochen; Läufe 1 und 2 sind unvollständig, Lauf 3 enthält das Limit noch nicht |

Die Läufe 1 und 2 jeder Stufe liefen mit warmem Build-Cache der Engine (der Host trägt 41 GB Build-Cache aus früheren Arbeiten; ein Lauf löst ihn nicht): Lauf 2 ist **kein** Lauf auf frischem Docker-Stand im Sinn des Slice (§1 „Grenzen vorab“). Ein kalter Bau des Images (Ollama-Archiv, Docker-Pakete) ist **nicht gemessen**; der Cache gehört dem Host und wurde nicht geleert (kein `builder prune`).

## Abweichung von der Vorprüfung des Slice (Stufe 4)

Vorprüfung (d) des Slice verlangt für Stufe 4 `MemAvailable` ≥ 4 × 3 GiB + 1,5 GiB = 13,5 GiB (14.155.776 kB; GiB, weil `--memory=3g` GiB bedeutet). Gemessen waren 13.525.120 kB: **nicht erfüllt**, der Lauf brach ab (`rohdaten-stufe4-alte-schwelle-*`). Auf Entscheidung des Projektinhabers (2026-10-31, „der Rechner ist nicht ausgelastet“, Swap 16 GiB ungenutzt) gilt **nur für Stufe 4**: Start erlaubt, wenn `MemAvailable` ≥ 9 GiB und `SwapFree` ≥ 8 GiB. Die Instanzgrenze bleibt `--memory=3g --memory-swap=3g`. Das Skript überwacht im Sekundentakt `MemAvailable`, `SwapFree`, `vmstat` `si`/`so` und ob die vorhandenen Container (hier `exciting_boyd`) noch laufen; Abbruch bei `si+so` ≥ 5000 KiB/s in fünf aufeinanderfolgenden Proben (ein ruhender Host zeigt 0, Auslagern im Megabyte-Takt über zehn Sekunden ist dauerhafter Druck, eine Spitze beim Start nicht) oder bei einem vorhandenen Container, der nicht mehr läuft oder `OOMKilled` meldet. **Ergebnis beider Läufe:** Host-Swap unverändert (16.383 MiB frei, Wachstum 0 MiB), `si+so` maximal 0 KiB/s, kleinstes `MemAvailable` 12.814 MiB (Lauf 1) und 12.772 MiB (Lauf 2), keine OOM-Meldung, keine Abbruchbedingung. Stufe 4 ist damit **nicht** „unter Swap-Druck“ gelaufen; ihre Dauer- und Platzwerte tragen nur den Vermerk der gelockerten Vorprüfung. Fremde Last wurde nicht beendet. Der Speicherbedarf der Instanzen lag weit unter der Grenze (siehe M5), die Schwelle der Vorprüfung ist eine Reserve für die Grenze, nicht ein Verbrauchswert.

## Ergebnis je Frage

### M0 Messstand
Siehe Tabelle oben und je Rohdatei der Abschnitt `M0`.

### M1 Container und Labels (N = 2, 3, 4)
Alle Container-IDs paarweise verschieden (2 von 2, 3 von 3, 4 von 4). Eine Instanz trägt 9 Labels, darunter `devcontainer.local_folder=<Ordner>` und `u-boot.instance=<Name>`; **`devcontainer.config_file` fehlt** an den Instanzen (leer). Der Container ohne Flag trägt `devcontainer.local_folder` und `devcontainer.config_file=<Ordner>/.devcontainer/devcontainer.json` (kein `u-boot.instance`). `docker ps -a --filter label=devcontainer.local_folder=<Ordner>` trifft alle Instanzen (zuletzt angelegte zuerst). Keine Einhängung `docker.sock` in einer Instanz.

### M2 Wechselt `${devcontainerId}`?
**Ja, für alle gemessenen Stufen.** Je Instanz ein eigener Wert (4 von 4, 3 von 3, 2 von 2 paarweise verschieden; je ein 52-stelliger Wert, z. B. Stufe 4 Lauf 1: `01ao4rjr…`, `0smvms16…`, `0vtvr7up…`, `1butg38m…`). Bezug: ohne Flag ergibt einen anderen Wert als mit nur `--id-label devcontainer.local_folder=<Ordner>` (alle sechs Läufe „Gleich: nein“); beide unterscheiden sich von den Werten der Instanzen. M2b entfällt (der Wert wechselt).

### M3 Volume-Namen und Inhaltsprobe
Schema je Instanz: `ubmess-<Nr>-workspace-<id>`, `ubmess-<Nr>-docker-<id>`, `ubmess-<Nr>-ollama-<id>` mit `<id>` = `devcontainerId` der Instanz. Zahl der Volumes: 3 × N (6, 9, 12), alle Namen je Familie paarweise verschieden. Stufe 4, Lauf 1 (Beispiel):

| Instanz | Workspace | Docker-Speicher | Modelle (`ollama`) |
|---|---|---|---|
| i1 | `ubmess-1-workspace-01ao4rjrnavbfb8l637nlc8fo74sujmb19dhd7vtltlmcr1kgeie` | `ubmess-1-docker-01ao4rjr…` | `ubmess-1-ollama-01ao4rjr…` |
| i2 | `ubmess-1-workspace-0smvms16ijk3fle7g9m2072l4kadptg2bcei8ffie7qtb98nme94` | `ubmess-1-docker-0smvms16…` | `ubmess-1-ollama-0smvms16…` |
| i3 | `ubmess-1-workspace-0vtvr7uprpntmehidjlkkgv88uqtml4lsct8u9r11a24dckkmpob` | `ubmess-1-docker-0vtvr7up…` | `ubmess-1-ollama-0vtvr7up…` |
| i4 | `ubmess-1-workspace-1butg38m48jh6pmes1dmjnvojscmgijlhdd7rml29g1ih87cfhct` | `ubmess-1-docker-1butg38m…` | `ubmess-1-ollama-1butg38m…` |

Die vollen Namen aller sechs Läufe stehen in den Rohdaten (Abschnitt `M2/M3`). **Inhaltsprobe über alle Paare** (Datei im Workspace-Volume, Datei im Modell-Volume, Image `ubmess-probe-<i>:1` im verschachtelten Docker): in allen sechs Läufen 0 Paare mit sichtbarem Fremdinhalt. **Podman-Speicher: nicht messbar auf diesem Host** (`kernel.apparmor_restrict_unprivileged_userns = 1`, Host nicht umgestellt); die Familie `containers` steht ohne Wert. **Standardprofil:** keine Volumes, ein Bind-Mount des Projektordners; alle N Instanzen und der Host sehen dieselbe Datei im Projektordner (alle sechs Läufe).

### M4 Image
Je Projekt ein Image `vsc-ubmess-<Nr>-<Hash>-uid` (1,99 GB laut `docker system df -v`, Sandbox; Standardprofil 1,41 GB). Jedes `up` baut das Image über den Build-Cache neu und setzt den Tag `…-uid:latest` um: Die Container der Instanzen tragen verschiedene Image-IDs (z. B. Stufe 2: `3c8cc762fde6`, `db6bbce97eec`), derselbe Tag `Config.Image`, der eindeutige Anteil je Image rund 25 kB. Das zweite bis vierte `up` dauerte 3 bis 5 s (Cache warm), das erste 4 bis 7 s (Cache der Engine ebenfalls warm; kalter Bau nicht gemessen).

### M5 Dauer, Platz, Speicher, CPU (Limit je Instanz)
| Stufe | Start der Instanzen (s, Lauf 1 / Lauf 2) | `docker stats` je Instanz (Speicher, CPU) | kleinstes Host-`MemAvailable` | Platz frei vor/nach |
|---|---|---|---|---|
| 2 | 4, 4 / 5, 5 | 59 bis 66 MiB von 3 GiB, 0,3 bis 0,5 % | 12.996 MiB (Lauf 1, Stichproben `[host]`) | 217 G / 217 G |
| 3 | 4, 5, 5 / 3, 4, 4 | 55 bis 62 MiB von 3 GiB, 0,4 bis 0,55 % | 12.958 MiB (Lauf 1, Stichproben `[host]`) | 217 G / 217 G |
| 4 (gelockerte Vorprüfung) | 5, 5, 4, 5 / 4, 4, 4, 4 | 56 bis 62 MiB von 3 GiB, 0,3 bis 0,5 % | 12.814 / 12.772 MiB | 217 G / 217 G |

**Die Grenze wurde in keinem Lauf erreicht:** `OOMKilled=false`, `ExitCode=0`, alle Instanzen laufen bis zum Ende. Gemessen ist der Zustand **direkt nach dem Start** (Startscripts inklusive `dockerd`, ein `docker pull busybox:1.36` im verschachtelten Docker); ein Go-Build, ein Modell-Pull und kalter Go-Cache sind **nicht** gemessen. Platzbedarf: `docker system df`: Images wuchsen um 0,09 GB (14,90 → 14,99 GB), Volumes um 28 MB (Stufe 4: 280 → 308 MB; je Instanz ein Docker-Speicher mit 7 MB, Workspace und Modelle einige Byte). Bezug ohne Limit (nur N = 2, beide Läufe): `HostConfig.Memory=0`, `NanoCpus=0`, `docker stats` 62 bis 65 MiB je Instanz; kein Unterschied im Start (4 s). Bei N = 3 und 4 „nicht gemessen, Host-Schutz“ (so im Slice).

### M6 `up` ohne Flag
**Neben laufenden Instanzen (alle sechs Läufe):** `up` ohne Flag legt **keinen** neuen Container an, es **übernimmt den zuletzt angelegten Container einer Instanz** (Stufe 2: i2, Stufe 3: i3, Stufe 4: i4; Dauer 1 bis 2 s; dieselben Volumes und Labels, `u-boot.instance` bleibt gesetzt). Die Instanz-Container tragen kein `devcontainer.config_file`; die CLI fällt ohne Treffer auf beiden Labels auf die Suche allein nach `devcontainer.local_folder` zurück und nimmt den ersten (neuesten) Treffer. `docker ps --filter label=devcontainer.local_folder=<Ordner>` trifft N Container. Ohne Instanz legt `up` ohne Flag einen Container mit `devcontainer.local_folder` und `devcontainer.config_file` an.
**Umgekehrt (Bestandscontainer ohne Flag zuerst, dann Instanz):** die Instanz (`--id-label devcontainer.local_folder=… --id-label u-boot.instance=i1`) legt einen **neuen** Container an (alle sechs Läufe), übernimmt den Bestandscontainer nicht; ein erneutes `up` ohne Flag trifft weiter den Bestandscontainer (zwei Container gesamt). `docker ps -a --filter label=devcontainer.local_folder=<Ordner>` listet die Instanz zuerst (neuer): die Suche nach dem ersten Treffer (`FindContainer`) liefert dann die **Instanz**, nicht den Bestandscontainer.

### M7 `--remove-existing-container`, `docker stop`
`up --remove-existing-container` mit den Labels von Instanz 1 entfernt **nur** deren Container (alter Container weg, neuer angelegt); die übrigen N − 1 Container bleiben derselbe Container und laufen. `devcontainerId` und alle Volumes bleiben gleich; Workspace-Datei, Modell-Datei und das Image im verschachtelten Docker von Instanz 1 sind nach der Neuanlage da. `docker stop` eines Containers (i2) lässt die anderen laufen (`Exited (0)` nur dort); alle 3 × N Volumes bleiben (alle sechs Läufe). Der Vorlauf mit Skriptfehler zählt hierfür nicht.

### M8 Namen
Zwei Läufe, gleiches Ergebnis. Als Label-Wert von `--id-label u-boot.instance=<Name>` nimmt die CLI und gibt unverändert am Container zurück (23 von 25 Namen, `gleich dem Eingabewert: ja`) und der Filter `docker ps --filter label=u-boot.instance=<Name>` trifft genau einen Container: Länge 1, 32, 63, 64 und 128; Großbuchstaben (`Abc`); `.`, `_`, `-` im Wert, am Anfang und am Ende (`.abc`, `_abc`, `-abc`, `abc.`, `abc_`, `abc-`); Leerzeichen; `"`; `'`; Nicht-ASCII (`ä-ö`); `,`; `*`; `/`; `:`. **Zwei Ausnahmen:** (1) der Wert `a=b` (Gleichheitszeichen im Wert): `devcontainer up` kehrt **nicht zurück** (siehe Hänger); (2) die leere Zeichenkette: die CLI lehnt ab mit Exit-Code 1 und `Unmatched argument format: id-label must match <name>=<value>`. Ein Test mit einem gewöhnlichen Container (`--label x=a=b`) zeigt, dass `docker ps --filter label=x=a=b` und `docker events --filter label=x=a=b` den Wert finden; die Docker-Filter sind also nicht die Ursache.

### M9 Aufräumen
Alle Läufe: Inventur nachher gegen vorher (`docker ps -a`, `docker volume ls`, `docker images`) **ohne Differenz**; keine Reste mit Präfix; Platz nachher 217 bis 218 G. Der erste Lauf `namen` und die abgebrochenen Läufe: die Container der CLI (zufällige Namen `bold_lehmann`, `elated_faraday`) wurden über das Label des Messordners entfernt, nicht über das Präfix; ein Teil der Abbrüche wurde von mir ausgelöst (s. u.), die Inventur danach war leer. Der Build-Cache der Engine wächst durch Läufe; er ist nicht Teil der Inventur und wurde nicht geräumt.

## Hänger im Lauf `namen` (Befund)

Die ersten drei Läufe von M8 hingen bei dem Namen `a=b` (`rohdaten-namen-abbruch-haenger-*`), einer 34 Minuten, einer 9 Minuten. Das Hilfs-Container-Protokoll endet nach `Container started`; im Hilfs-Container laufen dann `node … devcontainer up`, `docker events --format {{json .}} --filter event=start` und `docker run --sig-proxy=false -a STDOUT -a STDERR …`. Der Messcontainer läuft (`Up`), mit Label `u-boot.instance=a=b`. **Ursache nicht geklärt:** reproduziert in jedem Lauf mit einem Gleichheitszeichen im Wert (4 von 4), in keinem anderen der 24 Namen, in keiner Instanz der Stufen 2 bis 4 (deren Werte `i1` bis `i4` kein `=` tragen). Der Unterschied zu den Stufen ist also der Wert, nicht das Profil. Die Docker-Filter vertragen den Wert (Test oben); die Vermutung, dass die CLI den Wert beim Zusammensetzen ihres Wartens auf das Start-Ereignis falsch zerlegt, ist **nicht belegt**. Das erste Zeitlimit (`timeout 90`) griff nicht: `timeout` sendet SIGTERM, `docker run` reicht es an den Container weiter und `node` als PID 1 ignoriert es; mit `timeout -k 10` und `docker rm -f` des Hilfs-Containers greift das Limit (exit 137, Rohdaten `rohdaten-namen-lauf1.txt`, Abschnitt `BEFUND`). Für den Code-Slice heißt das: ein Instanzname mit `=` darf die CLI nicht erreichen (die Regel des CR sollte `=` ausschließen, Begründung: gemessen hängt `up`).

## Urteil zur Kernfrage

**Fassung (1): `${devcontainerId}` wechselt mit dem zweiten Id-Label; die gemessenen Volume-Familien (Workspace, Docker-Speicher, Modelle) sind je Instanz getrennt — für zwei, drei und vier gleichzeitige Instanzen**, je zwei Läufe (Stufe 4 mit gelockerter Vorprüfung, ohne Swap-Druck). Die Podman-Familie ist **nicht messbar auf diesem Host** und damit nicht belegt. Im Standardprofil teilen alle Instanzen den Projektordner (Bind-Mount).

## Nicht gemessen (offen)

- Podman-Speicher (`nestedRuntime: podman`): Host-Sysctl `kernel.apparmor_restrict_unprivileged_userns = 1`.
- Kalter Bau des Images (Cache war warm), Last unter einem Go-Build oder Modell-Pull, kalter Go-Cache: das Erreichen der Grenze 3 GiB ist daher nicht belegt, nur der Zustand nach dem Start.
- Mehr als vier Instanzen, andere Hosts (Colima, Docker Desktop, rootless, `arm64`).
- Bezug ohne Limit bei N = 3 und 4 (Host-Schutz).
- Ursache des Hängers bei `=` im Label-Wert.
- Ein Lauf auf frischem Docker-Stand (Cache geleert) je Stufe.
- Annahme „Limit je Instanz“ ist die Vorgabe des Projektinhabers, nicht zurückgefragt.
- ID-Gleichheit des Mess-Images mit dem Bestands-Image des Hilfs-Containers.
