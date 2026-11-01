# Messprotokoll: Host-Quellen für Git-Token und Identität — `gh auth token`, `git config`, `git credential fill`, Transport über `up`/`shell` und Wert-Stellen nach dem Start

Messung zum Slice [`slice-v1-devcontainer-host-git-uebernehmen`](../../plan/planning/done/slice-v1-devcontainer-host-git-uebernehmen.md) ([`LH-FA-DEV-009`](../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Gemessen sind die acht Messungen des Slice §6 vor dem Code: Git-Identität aus `GIT_AUTHOR_*`/`GIT_COMMITTER_*` (1), `remoteEnv` mit `${localEnv:…}` in Lebenszyklus-Befehlen (2), `gh auth token` auf dem Host (3), `git credential fill` (4), `git config`-Verhalten (5), Transport über `u-boot devcontainer up`/`shell` in beiden Profilen und mit Instanz (6), Fake-Quelle für den Docker-Test (7, Probe) und die Stellen, an denen der Wert nach dem Start liegt (8).

## Messstand

| Eigenschaft | Wert |
|---|---|
| Datum | 2026-11-01 (Läufe 14:17 bis 14:33 UTC) |
| Host | Linux 7.0.0-42-generic, `x86_64` |
| Docker | Client und Server 29.9.2 (`linux/amd64`), lokale Engine |
| git | 2.53.0 auf dem Host |
| gh | 2.46.0, angemeldet an `github.com` (Token aus `hosts.yml`, `~/.config/gh`; gleiche Anmeldung wie die des Projektinhabers) |
| Dev-Container-CLI | 0.89.0 (Meldung im Lauf) |
| Hilfs-Image | `u-boot-devcontainer-cli:0.89.0-f1a94cce795e` (von `u-boot devcontainer up` selbst gebaut); die reine CLI-Probe der M2 baute `ubmess-devcontainer-cli:0.89.0` aus demselben Kontext `internal/adapter/driven/docker/devcontainercli/` |
| `u-boot` | `0.17.0-dev` (`bin/u-boot-linux-amd64`, Bau mit `make build-binaries PLATFORMS="linux/amd64"`, docker-Modus) |
| Werte | nur Platzhalter (`platzh-m8-git-token-0123456789abcdef`, `h-name`/`h-email`); der echte Token des Hosts wird in **keiner** Rohdaten-Datei geführt — der M3-Lauf belegt ihn nur als Länge (93 Zeichen) und SHA-256-Präfix |
| Durchführung | Implementer, Skripte unten |

## Skripte und Aufruf

[`tools/messung/devcontainer-host-git-uebernehmen/messen-leicht.sh`](../../../tools/messung/devcontainer-host-git-uebernehmen/messen-leicht.sh) (M1, M3, M4, M5, M7-Probe; keine Docker-Voraussetzung) und [`tools/messung/devcontainer-host-git-uebernehmen/messen-schwer.sh`](../../../tools/messung/devcontainer-host-git-uebernehmen/messen-schwer.sh) (M2, M6, M8; Gruppen m2/m6/m8), aus der Repo-Wurzel je Lauf:

```bash
bash tools/messung/devcontainer-host-git-uebernehmen/messen-leicht.sh \
  > docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/rohdaten-m135-leicht-lauf<L>.txt 2>&1
bash tools/messung/devcontainer-host-git-uebernehmen/messen-schwer.sh <m2|m6|m8> bin/u-boot-linux-amd64 \
  > docs/messungen/slice-v1-devcontainer-host-git-uebernehmen/rohdaten-<Gruppe>-lauf<L>.txt 2>&1
```

Beide Skripte tragen WIEDERHOLUNGEN (Vorgabe 2) als Umgebungsvariable; die Läufe unten liefen mit 2, der leichte Lauf zusätzlich mit 3 vollständigen Wiederholungen des ganzen Skripts.

## Ergebnisse

### M1 — Git-Identität aus Umgebungsvariablen (Reihe für D1; D2 gewählt)

| Zelle | Ergebnis (je 3 Läufe) |
|---|---|
| alle vier `GIT_AUTHOR_*`/`GIT_COMMITTER_*` | Commit gelingt; Autor und Committer exakt aus den Variablen (`h-name|h-email|c-name|c-email`) |
| nur `GIT_AUTHOR_NAME`/`EMAIL` | Exit 128, „Committer identity unknown“ — der Commit braucht auch die Committer-Identität |
| nur Name oder nur E-Mail | Exit 128 („Author identity unknown“ beziehungsweise „empty ident name … not allowed“) |
| keine Variable, keine Konfiguration | Exit 128, „Please tell me who you are“ |
| `git var GIT_AUTHOR_IDENT` mit Variablen | Exit 0; ohne Exit 128 |

### M3 — `gh auth token` auf dem Host (Quelle 2)

| Zelle | Ergebnis |
|---|---|
| unveränderte Umgebung | Exit 0, Ausgabe nichtleer (93 Zeichen, SHA-256-Präfix in Rohdaten; Wert nie ausgegeben) |
| `gh` nicht im PATH | Exit 127 („command not found“) |
| nicht angemeldet (leere `GH_CONFIG_DIR`) | Exit 1, „no oauth token found for github.com“ |
| `--hostname github.com` | Exit 0 |
| `GH_TOKEN` in der Umgebung | der Umgebungs-Wert wird ausgegeben (31 Zeichen Platzhalter) — gh gibt die Variable, wenn gesetzt, vor der hosts.yml-Anmeldung |

### M4 — `git credential fill` ohne und mit Helper (Beleg für den Verzicht auf Quelle 3)

| Zelle | Ergebnis |
|---|---|
| ohne Helper, `GIT_TERMINAL_PROMPT=0` | Exit 128, „could not read Username …: terminal prompts disabled“, 10–11 ms — kein Hänger, aber nichts brauchbar |
| mit Helper (store, Platzhalter) | Exit 0, liefert `username`/`password` — die Quelle wäre mit Helper grundsätzlich tragbar |
| Helper-Datei leer | Exit 128, „No such device or address“ |

Die Entscheidung entfällt auf Quelle 3 (siehe Slice §6, entschieden 2026-11-01); die Messung bleibt als Beleg dafür.

### M5 — `git -C <Ordner> config user.name` / `user.email`

| Zelle | Ergebnis |
|---|---|
| Wert nur in global | wirksamer Wert, Exit 0 |
| Wert in local und global | local gewinnt |
| Wert in system und global | global vor system |
| Schlüssel nirgends | Exit 1, leere Ausgabe |
| Ordner ohne Repository | Exit 1 |
| leerer Wert gesetzt | Exit 0, Ausgabe leer (0 Zeichen) — vom fehlenden Schlüssel (Exit 1) unterscheidbar |
| `includeIf` auf den Projektordner | Bedingung greift, der eingefügte Wert gewinnt (Korrektur der Zelle: der erste Lauf trug ein Muster ungefiltert `$TMP`, beurteilt nur die Läufe 2 und 3) |
| keine Konfigurationsdatei überhaupt | Exit 1 |

### M2 — `remoteEnv` mit `${localEnv:…}` in Lebenszyklus-Befehlen (Referenz für Weg B, verworfen)

Minimaler Devcontainer (Alpine) mit `"remoteEnv": {"GIT_TOKEN": "${localEnv:GIT_TOKEN}"}` und `postCreateCommand`, das die Länge der Variable schreibt; reine CLI-Probe über das Hilfs-Image (Projektordner als Bind-Mount, wie `runArgs`).

| Zelle | Ergebnis (je 2 Läufe) |
|---|---|
| `GIT_TOKEN` im Hilfs-Container-Umfeld gesetzt | `postCreateCommand` liest den Wert: 36+1 Zeichen (Platzhalter `platzh-…`, 36 + Zeilentrenner), Exit 0 |
| ungesetzt | 0+1 Zeichen |

Die beiden Vorläufe (`rohdaten-m2-vorlauf-sonde-fehler.txt`) tragen eine fehlerhafte Sonde (verschränkte Quotes im `postCreateCommand`, danach `gt-vorhanden` an den Exit von `wc` gebunden); sie gelten nicht als Zelle.

### M6 — Transport über `u-boot devcontainer up` und `shell` (Platzhalter)

| Zelle | Ergebnis (je 2 Läufe) |
|---|---|
| A: Sandbox, Variable gesetzt, `up` | `up` Exit 0; `docker exec printenv GIT_TOKEN` → 0 Bytes: der Wert steht **nicht** in `Config.Env` des Containers; `git credential fill` im Exec ohne Sitzungsumgebung → 128 |
| B: Sandbox, gesetzt, `shell` | in der Sitzung 36+1 (Platzhalter + Trenner) und der Credential-Helper des Images liefert den Wert in der Sitzung (`fill-exit=0`) |
| C: Sandbox, ungesetzt, `shell` | in der Sitzung 0, Helper liefert nichts (128) — Warn-Fall des Ist-Zustands |
| D: Standardprofil, Variable gesetzt, `shell` | in der Sitzung 0, Helper nichts — ohne `remoteEnv`-Referenz trägt der Standardprofil-Aufruf nichts |
| E: Sandbox, gesetzt, `--instance ubmess-i1` | `up` Exit 0 (Labels siehe unten), Sitzung wie B |

### M7 — Probe der Fake-Quelle für den Docker-Test

| Zelle | Ergebnis |
|---|---|
| Fake-`gh` (Skript im PATH) | gibt den Platzhalter aus, Exit 0 |
| Fake-`gh`, leere Ausgabe | Exit 0, nichts |
| Fake-`gh`, Fehlerausgabe, Exit 1 | Exit 1, Fehlerzeile |

### M8 — Stellen des Werts nach dem Start (Platzhalter)

| Ort | Befund |
|---|---|
| Argumentlisten der `docker`-Prozesse des Hosts | nur `-e GIT_TOKEN` (Name); der Platzhalter steht in **keiner** Argumentliste |
| `/proc/<pid>/environ` des `docker`-Client-Prozesses des Hilfs-Containers | der Platzhalter ist dort (der Prozess-Umfeld-Kanal, wie der Port-Doc sagt) |
| `Config.Env` des gestarteten Devcontainers | Platzhalter nicht enthalten (nur `PATH`/`LANG`) |
| Label `devcontainer.metadata` | Platzhalter nicht enthalten; die **Referenz** `${localEnv:GIT_TOKEN}` ist dort (als Teil der erzeugten Konfiguration) |
| Prozesslisten im Hilfs-Container | die Proben (4-s-Raster, 10 Proben) trafen den inneren `docker`-Aufruf des Hilfs-Containers nicht (ein schneller Lauf, 2 bis 4 s); der M2-Lauf belegt, dass die CLI den aufgelösten Wert in die Befehle ihrer Container-Befehle schreibt — für die Länge des Befehls, nicht für die Dauer des Containers |

## Grenzen

- Ein echter Push mit dem echten Token ist nicht gefahren; die Zusage der M3 bleibt auf die Ausgabe der Anmeldung bezogen, nicht auf das Schreibrecht des Tokens.
- Das Verhalten ohne `GIT_TERMINAL_PROMPT=0` (ob `git credential fill` auf diesem Host fragen würde) ist nicht gemessen; die Quelle ist durch die Entscheidung entfallen.
- Mehrere gleichzeitige gh-Anmeldungen, `gh auth token --hostname` auf einen zweiten Host und ein Token ohne Schreibrecht: nicht gemessen.
- Der Hilfs-Container der M6/M8-Läufe war in den Prozesslisten nicht zu fassen (Lauf unter zwei Sekunden nach warmem Bau); die Probe ist ein Raster mit 4-s-Abstand — der Fund über `/proc/<pid>/environ` des Hosts stützt sich auf den `docker`-Client des Hilfscontainers.
- Nur dieser Host (Linux `amd64`, lokale Engine 29.9.2), Docker 28.0.4 und entfernter Engines: nicht gemessen.
- Die M8-Rohdaten der zwei Läufe stehen für eine Zelle; mehr Wiederholungen unter der Zeitbeschränkung liefen nicht.
- Beim Aufräumen blieben Objekte der Wegwerf-Projekte (Container mit dem Label `devcontainer.local_folder`, deren Volumes und die Devcontainer-Images) bis zur Prüfung durch den Verifier auf der Engine; der Nachweis der Bereinigung ist der Inventur-Teil der Skripte (`aufraeumen`, per trap), nicht dieser Protokoll-Lauf.