# Messprotokoll: ENV- und Cache-Mechanik der erzeugten Dev-Container-Konfiguration (Workspace-Variable, Block-ENV gegen `containerEnv`, Cache unter dem Workspace-Volume)

Messung zum Entwurf [`slice-v1-katalog-gepinnte-toolchain-entwurf`](../../plan/planning/done/slice-v1-katalog-gepinnte-toolchain-entwurf.md) (DoD Zeile 1), für [`LH-FA-BUILD-003`](../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (§4 Frage 3) und [`LH-FA-DEV-001`](../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)/[`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Form der erzeugten Konfiguration). Gemessen ist: löst die Dev-Container-CLI eine Workspace-Pfad-Form in `containerEnv`/`remoteEnv` auf, welchen Wert ein `ENV` aus dem Dockerfile-Block und `containerEnv` der Konfiguration bei Konkurrenz haben, und übersteht der Cache unter dem Workspace-Volume ein `up --recreate` — in beiden ENV-Trägerformen (Dockerfile-Block `ENV` und, als Zusatzreihe nach der Entscheidung zu §4 Frage 3, `containerEnv` mit `${containerWorkspaceFolder}`). Es lief **nicht** `u-boot devcontainer up`; das Skript spiegelt den `docker run` der CLI-Adapter (`internal/adapter/driven/docker/devcontainer_launcher.go`, Funktion `runArgs`, Zeilen 481 bis 511) und ruft die gepinnte CLI direkt.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-11-02 (Läufe 07:59 bis 08:10 UTC) |
| Host | Linux 7.0.0-42-generic, `x86_64` |
| Docker | Server 29.9.2 (`linux/amd64`), Engine lokal |
| Dev-Container-CLI | `0.89.0` (`devcontainer --version` im Mess-Image) |
| Basis-Image der Messprojekte | `mcr.microsoft.com/devcontainers/base:debian` (per `docker pull`; die CLI baut daraus je Lauf das Feature-Image, siehe Rohdaten-Logs) |
| Umgebung | linux-host `amd64`, lokale Docker-Engine (DoD-Umgebung) |

**Weg der CLI.** Das Mess-Image `ub-messung-devcontainer-cli:0.89.0` ist per `docker build` aus `internal/adapter/driven/docker/devcontainercli/` (Dockerfile, `package.json`, `package-lock.json`) gebaut — dieselben drei Dateien, die der Launcher als Build-Kontext einbettet (`devcontainer_launcher.go` Zeilen 64 bis 68, `go:embed`). Es läuft also derselbe CLI-Pin (`@devcontainers/cli 0.89.0`, `npm ci --ignore-scripts`) wie beim Produkt; der Mess-Tag trägt nur die Version und nicht den Content-Hash des Launcher-Tags. Benötigt werden für den Build die npm-Registry (`@devcontainers/cli@0.89.0`) und die zwei im Dockerfile per Index-Digest gepinnten Basis-Images; beide Quellen sind dieselben wie beim Produkt-Bau.

## Fundstellen (vor dem Messen gelesen, [`AGENTS.md`](../../../AGENTS.md) §3.10)

| Stelle | Was sie für die Messung trägt |
|---|---|
| `internal/adapter/driven/docker/devcontainer_launcher.go` Zeilen 481 bis 511 | Startweg des Helpers: Engine-Socket und Projekt-Ordner als Bind-Mount am selben Pfad, `devcontainer up --workspace-folder <folder>`; recreate ist `--remove-existing-container` (Zeilen 502 bis 504) |
| ebenda Zeilen 275 bis 280 | Session-Wege: `docker exec -u <RemoteUser>` (der Launcher-Weg) — `remoteEnv` ist auf diesem Weg nicht wirksam, das ist die Messstellen-Trennung der Zelle m1d |
| `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` Zeilen 12 bis 13 | Sandbox-Form: `workspaceMount: ""`, `workspaceFolder` unter `/workspaces/` — die Form, der die Mess-Zellen m1b/m1c/m3 folgen |
| ebenda Zeilen 33 bis 41 | Die Pipeline rendert heute `remoteEnv` mit `${localEnv:...}` und `runArgs`; **keine `containerEnv`-Stelle** |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` Zeilen 21 bis 24 und 151 bis 157 | ARG-Pin-Muster der Blöcke; `mkdir -p <workspaceFolder> && chown vscode:vscode` als Vorbereitung des Volume-Besitzers |
| `internal/hexagon/application/devcontainer_features.go` Zeilen 157 bis 171 und 305 bis 315 | Muster `dockerfileBlock`: der Pin des Blocks trägt im Dockerfile-Template; Einträge `ollama`, `gh`, `glab` |
| `.devcontainer/Dockerfile` Zeilen 75 bis 117 (ENV-Block ab Zeile 113) | Hand-written-Teil, den der Katalog-Eintragstyp ersetzen soll: Go- und `golangci-lint`-Install plus `ENV PATH/BUILD_MODE/GOFLAGS/GOTOOLCHAIN/GOCACHE/GOMODCACHE` mit festen `/workspaces/u-boot/.cache/…`-Pfaden |

## Skript und Aufruf

[`tools/messung/katalog-gepinnte-toolchain/`](../../../tools/messung/katalog-gepinnte-toolchain/) — [`lib.sh`](../../../tools/messung/katalog-gepinnte-toolchain/lib.sh) (Mess-Image-Bau, Helper-Aufruf, Messstellen, Aufräumen) und je Messung ein Skript, aus der Repo-Wurzel; die Rohdaten schreibt das Skript selbst in dieses Verzeichnis:

```bash
bash tools/messung/katalog-gepinnte-toolchain/m1-workspace-variablen.sh
bash tools/messung/katalog-gepinnte-toolchain/m2-env-block-vs-containerenv.sh
bash tools/messung/katalog-gepinnte-toolchain/m3-cache-recreate.sh
bash tools/messung/katalog-gepinnte-toolchain/m3z-cache-recreate-containerenv.sh
```

- Je Zelle zwei Läufe; jeder Lauf legt ein **frisches Wegwerf-Projekt** an (`mktemp -d /tmp/ub-messung-<zelle>-lauf<L>-XXXXXX`) — frischer Ordner heißt frisches `devcontainer.local_folder`-Label, also zwei unabhängige Starts statt Bestands-Übernahme (Bestand-Verhalten ist in [`slice-v1-devcontainer-container-name`](../slice-v1-devcontainer-container-name/README.md) gemessen).
- Namespace: Helper-Container `ub-messung-cli-…`, Volumes `ub-messung-<zelle>-lauf<L>-ws`, Projekt-Ordner und `vsc-…`-Feature-Images tragen das Präfix; jede Rohdaten-Datei endet mit der Inventar-Ausgabe (`leftovers after teardown`) der Container und Volumes.
- Es liegt ein Protokoll-Verwerfen vor (siehe M3 unten).

## Befunde

### M1 — Workspace-Variable in `containerEnv`/`remoteEnv` (vier Zellen, je zwei Läufe)

| Zelle | Form | Konformations-Form | Ergebnis je beide Läufe |
|---|---|---|---|
| m1a | `containerEnv` mit `${localWorkspaceFolder}` | Standard (CLI bind-mountet den Workspace selbst) | **Aufgelöst**: `.Config.Env` und `docker exec env` zeigen `GOCACHE=<Host-Ordner>/.cache/go-build` |
| m1b | `containerEnv` mit `${localWorkspaceFolder}` | Sandbox (`workspaceMount: ""`, named Volume) | **Aufgelöst — zum Host-Pfad**: `GOCACHE=<Host-Ordner>/.cache/go-build`; der Host-Ordner ist im Container **nicht erreichbar** (kein Bind-Mount), der Wert zeigt damit nicht aufs Workspace-Volume |
| m1c | `containerEnv` mit `${containerWorkspaceFolder}` | Sandbox | **Aufgelöst — zum Container-Pfad**: `GOCACHE=/workspaces/<Projekt>/.cache/go-build`, exakt der Pfad des named Volumes |
| m1d | `remoteEnv` mit `${localWorkspaceFolder}` | Standard | **Nicht auf dem Launcher-Weg**: `docker exec env` kennt kein `GOCACHE`; nur `devcontainer exec` (CLI-Weg) zeigt `GOCACHE=<Host-Ordner>/.cache/go-build` |

Lesarten mit Rohdatabeleg:

- Die CLI 0.89.0 löst Variablen in `containerEnv` und in `remoteEnv` **beim `up`**; kein Wert blieb als Literal stehen.
- `${localWorkspaceFolder}` ist der Host-Ordner (`--workspace-folder`-Wert). In der Sandbox-Form (m1b) zeigt ein daraus gebauter Cache-Pfad **nicht** unter das Workspace-Volume, sondern auf einen Pfad, den es im Container nicht gibt — genau die Grenze, über die der Slice §1 Punkt "ContainerEnv" schreibt: die Cache-Pfade sind im Katalog-Eintrag nicht statisch formulierbar, und dieser Variablenweg löst das im Sandbox-Profil nicht.
- `${containerWorkspaceFolder}` (m1c) löst zum Workspace-Ordner im Container — der einzige gemessene Variablenweg, der unter dem Workspace-Volume liegt und ohne Projekt-spezifische Pfad-Hardcodierung formulierbar bleibt.
- `remoteEnv` (m1d) wirkt nur auf dem `devcontainer exec`-Weg der CLI. Der Container-Umfang (`docker exec -u vscode …`) des Adapters bekommt die Werte nicht (Rohdaten m1d: Abschnitt `docker exec` ohne `GOCACHE`, Abschnitt `devcontainer exec` mit aufgelöstem Wert).

### M2 — `ENV` aus dem Dockerfile-Block gegen `containerEnv` (drei Zellen, je zwei Läufe)

| Zelle | Dockerfile-Block | `containerEnv` | Wirksamer Wert (`.Config.Env` und `docker exec env`) |
|---|---|---|---|
| m2a | `GOCACHE=/ub-messung-block/go-build` | — | Block-Wert |
| m2b | Block-Wert | gleicher Wert | Block-Wert (`= config-Wert`) |
| m2c | Block-Wert | `GOCACHE=/ub-messung-config/go-build` | **`containerEnv`-Wert** — `containerEnv` übersteuert den Block-`ENV` bei gleichem Namen |

- Ein `ENV` eines Dockerfile-Blocks greift im Container (m2a), belegt über dieselbe Messstelle wie die andern Zellen (`docker inspect .Config.Env`, `docker exec -u vscode env`).
- Bei Konkurrenz trägt `.Config.Env` nur **einen** Wert je Namen: der `containerEnv`-Eintrag ersetzt den Bild-`ENV`, die Liste zeigt keinen zweiten Eintrag (Rohdaten m2c).
- Vorab-Lauf verworfen: Die erste M2-Reihe (07:59, alle sechs Läufe) trug einen Skript-Fehler — die Zellen-Übergabe (`mode`-Wort) paßte nicht zum Vergleich im Skript, alle Zellen liefen ohne `containerEnv` (zusammen also sechs weitere m2a-Bestätigungen). Die aufbewahrten Rohdaten stammen aus der korrigierten Reihe; die erste Reihe wurde von derselben Dateibasis überschrieben (siehe M3, gleicher Verwerfen-Mechanismus).

### M3 — Cache unter dem Workspace-Volume über `up --recreate` (Dockerfile-ENV-Form, zwei Läufe)

Ablauf je Lauf: frisches Sandbox-Form-Projekt (named Volume am `workspaceFolder`), Dockerfile-Block setzt `GOCACHE=/workspaces/m/.cache/go-build` und `GOMODCACHE=/workspaces/m/.cache/gomod`; erster `up`, `docker exec` als `vscode` schreibt Marker + 64-KiB-Pad in beide Caches, dann `up --remove-existing-container` (die recreate-Form des Adapters), Prüfung im neuen Container:

| Prüfstelle | Lauf 1 | Lauf 2 |
|---|---|---|
| Container-ID gewechselt | ja | ja |
| Marker-Datei `go-build` im Volume nach recreate | ja (`go-build-lauf1 2026-11-02T08:01:25+00:00`) | ja (`go-build-lauf2 2026-11-02T08:01:32+00:00`) |
| Marker-Datei `gomod` im Volume nach recreate | ja | ja |
| Inhalt (Marker + Pad, 65536 B) nach recreate | vollständig | vollständig |
| Volume nach recreate | gleiches Volume (`ub-messung-m3-lauf1-ws`, Besitzerrecht `vscode:vscode` in den `ls -la` der Rohdaten) | gleiches Volume |

- Protokoll-Verwerfen: Der erste M3-Lauf der Reihe brach an einem Skript-Fehler ab (Dockerfile- und Konfigurationspfad standen nicht unter `.devcontainer/` und die CLI fand die Konfiguration nicht, 08:00 UTC); das abgebrochene Projekt-Verzeichnis blieb und wurde nachträglich geräumt. Die aufbewahrten Läufe stammen aus der korrigierten Reihe (08:01 UTC). Weil `set -euo pipefail` vor dem ersten `up` abbrach, entstanden dabei keine Container oder Volumes.

### M3-Zusatz — Cache über `up --recreate` in der containerEnv-Form (zwei Zellen, je zwei Läufe)

Der Gegenlauf zur `m1c`-Form, die der Projektinhaber in §4 Frage 3 gewählt hat (Entscheidung 2026-11-02, Commit `b89943ea`): das Dockerfile trägt **keinen** `ENV`, die ENV kommt allein aus `containerEnv` — je Zelle ein Name (`GOCACHE` beziehungsweise `GOMODCACHE`) mit dem Wert `${containerWorkspaceFolder}/.cache`; sonst dieselbe Sandbox-Form und derselbe Ablauf wie M3 (Marker + 64-KiB-Pad, dann `up --remove-existing-container`, Prüfung im neuen Container):

| Prüfstelle | m3z-gocache Lauf 1 / Lauf 2 | m3z-gomod Lauf 1 / Lauf 2 |
|---|---|---|
| `containerEnv`-Variable aufgelöst (vor recreate) | `GOCACHE=/workspaces/m/.cache` (beide Läufe) | `GOMODCACHE=/workspaces/m/.cache` (beide Läufe) |
| Container-ID gewechselt | ja / ja | ja / ja |
| Variable nach recreate erneut aufgelöst | ja / ja | ja / ja |
| Marker-Datei im Volume nach recreate | ja / ja (Marker `m3z-gocache-lauf1`/`-lauf2` mit Zeitstempel) | ja / ja (Marker `m3z-gomod-lauf1`/`-lauf2`) |
| Inhalt (Marker + Pad, 65536 B) nach recreate | vollständig / vollständig | vollständig / vollständig |
| Eigentumsrechte nach recreate | `vscode:vscode` (beide, `ls -la` der Rohdaten) | `vscode:vscode` (beide) |
| Volume nach recreate | gleiches Volume (beide) | gleiches Volume (beide) |

- Der Pfad der `m1c`-Messung bleibt unter dem Volume erhalten, und die Form übersteht das Recreate mit denselben Prüfstellen wie das Block-`ENV`-Vorbild von M3 — andere Trägerform der ENV, gleiches Ergebnis.

## Rohdaten

| Datei | Inhalt |
|---|---|
| `rohdaten-m1{a,b,c,d}-lauf{1,2}.txt` + gleiche Namen mit `.log` | M1, je Zelle je Lauf die Konfiguration, die `.Config.Env`-Liste und die env der beiden exec-Wege; das Log trägt den stdout/stderr des `up` |
| `rohdaten-m2{a,b,c}-lauf{1,2}.txt` + `.log` | M2 |
| `rohdaten-m3-lauf{1,2}.txt` + `.log` | M3 komplett (Marker-Listen vor und nach dem recreate) |
| `rohdaten-m3z-gocache-lauf{1,2}.txt` + `.log` | M3-Zusatz, Zelle containerEnv/GOCACHE |
| `rohdaten-m3z-gomod-lauf{1,2}.txt` + `.log` | M3-Zusatz, Zelle containerEnv/GOMODCACHE |

## Aufräumen (Beobachtung BEO-E2E/docker-test-hinterlaesst-zwischencontainer)

- Container und Volumes: je Lauf im Skript geräumt; die Inventar-Abschnitte in jedem Rohdaten-Log zeigen nach dem Aufräumen je **leer**.
- Projekt-Ordner: je Lauf entfernt; das eine Verzeichnis des verworfenen M3-Vorlaufs wurde nach dem Abbruch manuell geräumt (`ls /tmp | grep ub-messung` leer).
- `vsc-ub-messung-*`-Feature-Images, die die CLI je Lauf baute (60 Stück): nach der Messung per `docker rmi` entfernt, Restbefund 0.
- Mess-Image `ub-messung-devcontainer-cli:0.89.0`: nach der Messung entfernt; die Skripte bauen es bei Bedarf neu (siehe Abschnitt Weg der CLI).
- Die CLI hinterlässt Build-Cache der Feature-Image-Bauten auf der Engine; der wurde nicht geprüft und nicht geräumt (Grenze unten).

## Grenzen der Aussage

- Gemessen ist die CLI 0.89.0 auf Docker Engine 29.9.2, Linux `amd64`, lokale Engine, Basis-Image `mcr.microsoft.com/devcontainers/base:debian`. Andere CLI-Versionen, andere Engines (Podman, Colima) und andere Basis-Images wurden nicht gemessen.
- Es lief kein `u-boot devcontainer up`; der Helper spiegelt den `docker run` des Adapters. Der Content-Hash-Tag des Launcher-Images wurde nicht nachgerechnet — der Kontext ist aus denselben drei Embedded-Dateien gebaut.
- M1 misst je Variablenform nur die eine Kombination aus ENV-Schlüssel (`containerEnv`/`remoteEnv`) und Konformations-Form, wie die Zellen-Tabelle sie nennt; `remoteEnv` mit `${containerWorkspaceFolder}` wurde nicht gemessen, ebensowenig `containerEnv` mit `${containerWorkspaceFolder}` im Standardprofil.
- M3 misst die Dockerfile-ENV-Form mit `$GOCACHE`/`$GOMODCACHE` **unterhalb** des gemounteten Workspace-Ordners; die `containerEnv`/`${containerWorkspaceFolder}`-Form (M1c) ist im M3-Zusatz mit denselben Prüfstellen gemessen. Was die Messung nicht trennt: die beiden Formen werden nicht am selben Projekt gemessen; der Gleichlauf folgt den Paar-Messungen, nicht einem Doppelträger-Vergleich in einem Lauf.
- Nur zwei Läufe je Zelle auf einem Host; keine Statistik, eine Reproduktion.