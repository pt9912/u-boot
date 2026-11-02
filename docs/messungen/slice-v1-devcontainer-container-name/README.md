# Messprotokoll: Containernamen des Entwicklungscontainers (Namensquellen, Engine-Regeln, Kollision, Identität)

Messung zum Slice [`slice-v1-devcontainer-container-name`](../../plan/planning/done/slice-v1-devcontainer-container-name.md) ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Gemessen ist, wie der Entwicklungscontainer, den `u-boot devcontainer up` startet, seinen Namen bekommt: nimmt die gepinnte Dev-Container-CLI einen Namen an (M1), löst sie Variablen in `runArgs` (M2), bleibt damit die Identität des Containers unverändert (M3, M6), welche Namensregeln die Engine anwendet (M4) und was bei Kollisionen und Bestand geschieht (M5). Es lief **nicht** `u-boot devcontainer up`; das Skript spiegelt den `docker run` von `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`) und ruft die CLI je Frage direkt.

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-11-01 (Läufe 17:27 bis 17:33 UTC) |
| Host | Linux 7.0.0-42-generic, `x86_64`, Ubuntu 26.04.1 LTS, 8 CPUs, rund 15,2 GiB RAM nutzbar, 16 GiB Swap, `/` 301 GB mit rund 244 GB frei (Docker-Wurzel liegt auf `/`) |
| Docker | Client und Server 29.9.2 (`linux/amd64`), Engine lokal, Treiber `overlayfs`, cgroup v2 (systemd) |
| Dev-Container-CLI | `0.89.0` (`devcontainer --version` im Mess-Image; Kontext = `internal/adapter/driven/docker/devcontainercli/`) |
| `u-boot` | `u-boot version 0.17.0-dev` (`bin/u-boot-linux-amd64`; erzeugt nur die Projektordner) |
| Profile | Sandbox mit `nestedRuntime: docker` (ohne Feature `ollama`) und Standardprofil; Projektordner wie stets Wegwerf-Ordner mit Präfix `ubmess-` |
| Lockerung | Der Hilfs-Container bekommt den Docker-Socket des Hosts (wie in `runArgs`) |

## Skript und Aufruf

[`tools/messung/devcontainer-container-name/messen.sh`](../../../tools/messung/devcontainer-container-name/messen.sh), aus der Repo-Wurzel, ein Modus je Messung:

```bash
bash tools/messung/devcontainer-container-name/messen.sh <cli|var|identitaet|engine|bestand|config> \
  bin/u-boot-linux-amd64 > docs/messungen/slice-v1-devcontainer-container-name/rohdaten-<modus>-lauf<L>.txt 2>&1
```

Alle Objekte tragen `ubmess-`/`vsc-ubmess-`; Aufräumen per `trap` nach Label `devcontainer.local_folder` bzw. nach Namen, Inventur vorher und nachher mit Differenz — in jedem Rohdatenlauf leer abgesehen vom Mess-Image, das der Lauf baut und wieder entfernt. Je Modus zwei Läufe; jeder Lauf trägt in der letzten Zeile `Stand des Modus: vollstaendig`.

## Befunde

### M1 — Namensquellen der CLI (Modus `cli`, zwei Läufe)

- Die Hilfe zu `devcontainer up` (CLI 0.89.0) kennt **keine Option** für den Containernamen; Zeilen mit „name“ sind nur `--id-label` und `--remote-env`.
- Die Eigenschaft **`containerName`** der `devcontainer.json` wird **nicht angewendet**: Ohne Eigenschaft heißt der Container `ecstatic_bell` respektive `stoic_meitner` (Zufallsname der Engine), mit Eigenschaft `sweet_burnell` respektive `heuristic_carson`.
- **`runArgs` mit `--name=<name>` wird angewendet**: Der Container heißt genau `ubmess-name-runargs`.

### M2 — Variablen in `runArgs` (Modus `var`, zwei Läufe)

- `runArgs` mit `--name=ubmess-name-var-${devcontainerId}` ergibt am Container `ubmess-name-var-<52-stellige devcontainerId>`; die CLI löst **`${devcontainerId}` in `runArgs` auf**. Der Name ist damit je Instanz eindeutig, aber nicht sprechend.

### M3 — Identität mit gesetztem Namen (Modus `identitaet`, zwei Läufe)

Vier Zellen (Sandbox und Standardprofil, jeweils mit Id-Label `u-boot.instance=i1` und ohne), Ablauf je Zelle: Referenzlauf ohne Namen (Container samt Volumes entfernt), Lauf mit `runArgs`-Name, Vergleich:

| Zelle | devcontainerId | Volume-Namen | Label-Wert | Name am Container |
|---|---|---|---|---|
| Sandbox, mit Instanz | gleich | gleich | gleich (`i1`) | `ubmess-name-sandbox-mit` |
| Sandbox, ohne Instanz | gleich | gleich | gleich (leer) | `ubmess-name-sandbox-ohne` |
| Standard, mit Instanz | gleich (Bind-Mounts, beide leer) | gleich (leer) | gleich (`i1`) | `ubmess-name-standard-mit` |
| Standard, ohne Instanz | gleich (leer) | gleich (leer) | gleich (leer) | `ubmess-name-standard-ohne` |

Ein gesetzter Name verändert die Identität des Containers nicht ([ADR-0038](../../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidung 3 bleibt unberührt); die Sandbox-Zellen tragen die tatsächlichen 52-stelligen Werte in den Rohdaten (paarweise je Zelle identisch).

### M4 — Namensregeln der Engine (Modus `engine`, zwei Läufe mit Fix)

Direkt mit `docker run --name …` gegen die lokale Engine 29.9.2, ohne Hilfs-Container:

- **Erlaubt** (exit=0): `ubmess-name-proj-inst`; Ziffer am Anfang (`9ubmess`); Großbuchstabe (`Uubmess`); `.` und `_` innerlich (`ubmess.1`, `ubmess_1`); Längen 64, 65, 128, 129, 256 und 257 — auf diesen Längen wurde kein Zurückweisen gemessen.
- **Zurückgewiesen** (exit=125, Wortlaut `docker: Error response from daemon: invalid container name (<Name>): only [a-zA-Z0-9][a-zA-Z0-9_.-] are allowed`): führender Bindestrich, Punkt und Unterstrich (`-ubmess`, `.ubmess`, `_ubmess`); Leerzeichen, `"` und `/` und `:` und `@` im Namen; Umlaut (`übmess`).
- **Kollision**: laufender **und** gestoppter Container gleichen Namens liefern exit=125 mit `Conflict. The container name "/<Name>" is already in use by container "<ID>". You have to remove (or rename) that container to be able to reuse that name.`; der haltende Container bleibt bei beiden Proben vorhanden (je 1 Treffer).

Vorablauf verworfen: Die ersten beiden `engine`-Läufe der Reihe (17:28:40 und 17:28:53 UTC) enthielten einen Skriptfehler — das Aufräumen der Probe entfernte den Container, der den Namen hielt, womit die gestoppte Kollision nicht gültig gemessen war. Die aufbewahrten Rohdaten stammen aus den Wiederholungsläufen mit dem Fix (Aufräumen nur bei exit=0), 17:32:04 und 17:32:31 UTC.

### M5 — Kollision und Bestand über die CLI (Modus `bestand`, zwei Läufe)

Mit `runArgs`-Name `ubmess-name-konflikt` im Standardprofil:

- **Bestand, gleicher Projektordner, gleicher Name**: Der zweite `devcontainer up` übernimmt denselben Container (exit=0, identische Container-ID) statt einer Kollision.
- **Fremder Container gleichen Namens, laufend** (anderes Label `devcontainer.local_folder`, ohne Dev-Container-Bezug): `devcontainer up` endet exit=1 (CLI meldet `{"outcome":"error","message":"Command failed: docker run …"}` mit Stack-Dump); der fremde Container bleibt vorhanden. Auch mit `--remove-existing-container` endet der Lauf exit=1, und der fremde Container bleibt — die CLI entfernt nach Label, nicht nach Name, und benennt den fremden Container nicht um.
- **Gestoppter Container gleichen Namens** (ohne Labels): `up` endet exit=1; `--remove-existing-container` endet exit=1, der gestoppte Container bleibt.
- Grenze der Rohdaten: Die `Ergebnis`-Zeile ist auf 260 Zeichen geschnitten; der Docker-Wortlaut **hinter** der CLI-Meldung (hier `Conflict … is already in use`) steht in den Rohdaten nicht. Der Wortlaut ist an der Engine direkt in M4 gemessen; für die Kollisionsmeldung zählt hier der Exit-Weg, nicht der Wortlaut.

### M6 — Name je Aufruf über `--override-config` (Modus `config`, zwei Läufe)

Der Hilfs-Aufruf könnte eine je Aufruf gerenderte Kopie der erzeugten `devcontainer.json` mit dem Namen in `runArgs` an die CLI übergeben. Gemessen am Sandbox-Profil (mit Id-Label `u-boot.instance=i1`), Referenzlauf ohne Namen:

- Der Lauf mit `--override-config <Kopie>` und `--name=ubmess-name-cfg-i1` in `runArgs` setzt den Namen am Container (`ubmess-name-cfg-i1` — die Form `u-boot-<projekt>-<instanz>` ist damit je Instanz über den Hilfs-Aufruf erreicht, ohne die erzeugte Datei zu berühren).
- Identität unverändert: `devcontainerId` gleich (Rohwerte je Lauf paarweise identisch, siehe Rohdaten), Volume-Namen gleich, Label-Wert gleich (`i1`).
- Die erzeugte `devcontainer.json` enthält den Namen nach dem Lauf fern (Gegenprobe: 0 Treffer).
- Vorlauf verworfen: der erste `config`-Lauf (17:32:08 UTC) brach an einer `set -u`-Zeile des Skripts ab (`dir: unbound variable`, keine Objekte angelegt); der aufbewahrte Lauf ist die Wiederholung mit dem Fix, 17:32:58 UTC.
- Sonde zur Aufbewahrungsort-Frage (17:49 UTC, nicht als Rohdaten aufbewahrt, Protokoll hier): Overlay-Datei in einem **fremden Verzeichnis** neben dem Projekt (wie der Hilfs-Aufruf sie bei einem `/tmp`-Ablageort hält) — `up` mit `--workspace-folder <Projekt> --override-config <fremde Datei>` endet exit=0, der Container heißt `ubmess-sonde-ov1`, und das Workspace-Volume (`<projekt>-workspace-<devcontainerId>`) ist identisch zum Referenzlauf ohne Namen. Die CLI löst die relativen Pfadangaben der Konfiguration unabhängig vom Aufbewahrungsort der Overlay-Datei auf.

## Rohdaten

| Datei | Inhalt |
|---|---|
| `rohdaten-cli-lauf1.txt`, `-lauf2.txt` | M1, vollständig |
| `rohdaten-var-lauf1.txt`, `-lauf2.txt` | M2, vollständig |
| `rohdaten-identitaet-lauf1.txt`, `-lauf2.txt` | M3, vier Zellen, vollständig |
| `rohdaten-engine-lauf1.txt`, `-lauf2.txt` | M4, Wiederholungsläufe mit Fix (siehe oben), vollständig |
| `rohdaten-bestand-lauf1.txt`, `-lauf2.txt` | M5, vollständig |
| `rohdaten-config-lauf1.txt` | M6, Wiederholungslauf mit Fix (siehe oben), vollständig |
| `rohdaten-config-lauf2.txt` | M6, vollständig |

## Grenzen der Aussage

- Gemessen ist die gepinnte CLI 0.89.0 auf Docker Engine 29.9.2, Linux `amd64` (lokale Engine). Andere CLI-Versionen, andere Engines (Podman, Colima) wurden nicht gemessen.
- M1 bis M3 und M6 setzen den Namen über `runArgs` der von der CLI gelesenen Konfiguration; ob `u-boot` die erzeugte `devcontainer.json` verändert oder eine je Aufruf gerenderte Konfiguration übergibt, ist eine Entscheidung des Slices, keine Messfrage.
- Die Längen-Messreihe in M4 prüfte 64 bis 257 Zeichen mit Buchstabe-Ziffer-Füllung; eine Obergrenze wurde auf diesen Längen nicht erreicht, eine Grenze oberhalb wurde nicht gesucht.
- Für die Kollisionsmeldung von `up` gilt die Grenze unter M5 (Wortlaut hinter der 260-Zeichen-Kürzung): Der Exit-Weg ist gemessen, der exakte Wortlaut der CLI stammt aus M4 an der Engine.