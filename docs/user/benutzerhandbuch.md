# Benutzerhandbuch: u-boot

Handbuch-Version: 2.0
Software-Version: v0.13.1
Stand: 2026-10-22

---

## Inhalt

1. [Einleitung](#1-einleitung)
2. [Installation](#2-installation)
3. [Erste Schritte](#3-erste-schritte)
4. [Aufgaben](#4-aufgaben)
5. [Konfiguration](#5-konfiguration)
6. [Rollen und Rechte](#6-rollen-und-rechte)
7. [Fehlerbehebung](#7-fehlerbehebung)
8. [FAQ](#8-faq)
9. [Glossar](#9-glossar)
10. [Anhang](#10-anhang)
11. [Support und Lizenz](#11-support-und-lizenz)

---

## 1. Einleitung

### Zweck der Software

`u-boot` legt Entwicklungsumgebungen an, die auf jedem Rechner gleich
aussehen. Statt eine Projektstruktur, eine `compose.yaml`, eine
Devcontainer-Konfiguration und die üblichen Begleitdateien von Hand zu
schreiben, erzeugt und pflegt `u-boot` sie für Sie. Dienste wie PostgreSQL,
Keycloak oder einen OpenTelemetry-Collector fügen Sie mit einem Befehl hinzu.

Der Kern ist ein einzelner Befehl: `u-boot init` erzeugt ein vollständiges,
lauffähiges Projekt. Danach begleitet Sie dasselbe Werkzeug durch den Alltag:
Dienste hinzufügen und entfernen, die Umgebung starten und stoppen, Logs
ansehen, Artefakte wie README, CHANGELOG oder Devcontainer aktualisieren.

`u-boot` ist **kein** Deployment-Werkzeug und kein Ersatz für Docker Compose.
Es erzeugt und pflegt die Dateien, mit denen Sie arbeiten; das Starten
übernimmt darunter weiterhin Compose.

### Zielgruppe dieses Handbuchs

Entwicklerinnen und Entwickler, die ein neues Projekt aufsetzen oder ein
bestehendes `u-boot`-Projekt betreuen. Sie brauchen kein Wissen über den
inneren Aufbau von `u-boot`. Grundkenntnisse in der Kommandozeile und ein
grobes Verständnis von Docker Compose genügen.

### Voraussetzungen

| Voraussetzung | Wofür |
|---|---|
| `u-boot` (Binary, Paket, Homebrew oder Container-Image) | alle Befehle |
| Docker Engine ab 24.0 | `up`, `down`, `logs` und die Docker-Prüfungen von `doctor` |
| Docker-Compose-Plugin ab 2.20 | dieselben Befehle (nicht `devcontainer up` und `devcontainer down`) |
| Docker (Client und lokale Engine; gemessen mit 29.9.1) | `devcontainer up` und `devcontainer down` (Abschnitt 4.9) |
| `git` | `init` (außer mit `--no-git`) |

**Ohne Docker** funktionieren `init`, `add`, `remove`, `generate`, `config`,
`egress list` und `template list`. Nur `up`, `down`, `logs`, `devcontainer up` und
`devcontainer down` brauchen einen erreichbaren Docker-Daemon. `u-boot doctor` meldet
fehlende Werkzeuge als Befund.

Podman ab 4.0 funktioniert als Ersatz für Docker, wenn `docker` auf `podman`
zeigt und `DOCKER_HOST` auf den Podman-Socket gesetzt ist. `u-boot doctor`
meldet Podman-Versionen als Warnung („unrecognized version"), bricht aber nicht
ab. Für `devcontainer up` und `devcontainer down` ist Podman als Host-Engine nicht zugesagt.

---

## 2. Installation

Wählen Sie **eine** der folgenden Möglichkeiten. Danach prüfen Sie immer mit

```bash
u-boot --version
```

**Ergebnis:** Die installierte Version erscheint, zum Beispiel
`u-boot version 0.13.1`.

### Binary (Linux, macOS, Windows)

Ein einzelnes, statisch gelinktes Programm ohne weitere Abhängigkeiten.

**Linux und macOS:**

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
```

**Windows (PowerShell):**

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

Für eine bestimmte Version ersetzen Sie `latest/download` durch
`download/v0.13.1`. Die Prüfsummen aller Dateien stehen in `SHA256SUMS` am
jeweiligen Release.

### Homebrew (macOS)

Für stabile Releases gibt es einen Homebrew-Tap für macOS (Apple Silicon und
Intel). Unter Linux nutzen Sie das Binary oder das Container-Image; Linux-Brew
ist nicht zugesagt:

```bash
brew tap pt9912/u-boot
brew trust pt9912/u-boot
brew install u-boot
```

`brew trust` ist nötig, weil Homebrew 5 Programme aus Drittanbieter-Taps erst
installiert, nachdem Sie den Tap ausdrücklich als vertrauenswürdig markiert
haben. Ein Update holen Sie mit `brew upgrade u-boot`. Vorabversionen
(`vX.Y.Z-rc.1`) landen nicht im Tap.

### Debian- und RPM-Pakete (Linux)

Zu jedem Release gibt es `.deb`- und `.rpm`-Pakete für `amd64` und `arm64`.
Ein APT- oder DNF-Repository gibt es nicht; Sie laden das Paket und
installieren es lokal:

```bash
VERSION=0.13.1
# Debian / Ubuntu
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot_${VERSION}_amd64.deb"
sudo apt install "./u-boot_${VERSION}_amd64.deb"
# Fedora / RHEL
curl -fsSLO "https://github.com/pt9912/u-boot/releases/download/v${VERSION}/u-boot-${VERSION}-1.x86_64.rpm"
sudo dnf install "./u-boot-${VERSION}-1.x86_64.rpm"
```

Für `arm64` verwenden Sie `_arm64.deb` beziehungsweise `.aarch64.rpm`. Das
Programm liegt danach unter `/usr/bin/u-boot`. Ein Update erfolgt durch die
Installation des neueren Pakets.

### Container-Image

```bash
docker run --rm -v "$PWD:/work" -w /work ghcr.io/pt9912/u-boot:latest init mein-projekt --no-git
```

Zwei Einschränkungen im Container:

- **`git` fehlt im Image.** Nutzen Sie `init --no-git` und legen Sie das
  Repository anschließend selbst an (`git init`).
- **`u-boot doctor` sieht Ihre Werkzeuge nicht.** Die vier Host-Prüfungen
  (Docker, Compose-Plugin, Docker-Erreichbarkeit, git) werden übersprungen und
  mit `?` markiert. Auch die Prüfung des Geräts `/dev/fuse` entfällt (Eintrag
  `devcontainer.sandbox.runtime`, Info): das Gerät im Image gehört dem
  u-boot-Container, nicht dem Container-Host. Auch `up`, `down`, `logs`,
  `devcontainer up` und `devcontainer down` laufen im Image nicht, weil dort kein Docker
  vorhanden ist (Exit-Code 11).

Für den Alltag ist ein installiertes Programm bequemer als das Image.

### Shell-Vervollständigung (optional)

```bash
source <(u-boot completion bash)
```

Unterstützt werden `bash`, `zsh`, `fish` und `powershell`
(`u-boot completion <shell> --help` zeigt den dauerhaften Einrichtungsweg der
jeweiligen Shell).

---

## 3. Erste Schritte

### Ein Projekt in zwei Minuten

#### Voraussetzung

`u-boot` und Docker sind installiert. Sie stehen in einem leeren Verzeichnis.

#### Vorgehen

1. Prüfen Sie Ihre Umgebung:

   ```bash
   u-boot doctor
   ```

2. Legen Sie das Projekt an:

   ```bash
   u-boot init demo-app
   ```

3. Fügen Sie eine Datenbank hinzu:

   ```bash
   u-boot add postgres
   ```

4. Kopieren Sie die Beispiel-Umgebung und tragen Sie ein eigenes Passwort ein:

   ```bash
   cp .env.example .env
   ```

   Ersetzen Sie in `.env` den Platzhalter `CHANGEME_POSTGRES_PASSWORD`.

5. Starten Sie die Umgebung:

   ```bash
   u-boot up
   ```

#### Ergebnis

Nach Schritt 2 meldet `u-boot`:

```text
Initialized u-boot project "demo-app".

Created:
  - docker/
  - scripts/
  - docs/
  - README.md
  - CHANGELOG.md
  - compose.yaml
  - .env.example
  - .gitignore
  - u-boot.yaml
```

Nach Schritt 5 läuft Ihr Stack. `u-boot up` wartet, bis jeder Dienst stabil ist,
und zeigt eine Statustabelle:

```text
SERVICE   CONTAINER  PORT                  HEALTH
postgres  running    5432:5432, 5432:5432  healthy
```

(Der Port erscheint doppelt, weil Docker ihn für IPv4 und IPv6 meldet.)

#### Hinweise

- Ohne Namen (`u-boot init`) leitet `u-boot` den Projektnamen aus dem
  Verzeichnisnamen ab.
- `u-boot init` legt auch ein Git-Repository an. Mit `--no-git` unterbleibt das.
- Aufräumen: `u-boot down` stoppt die Umgebung (Abschnitt 4.5).

### Grundlegende Bedienkonzepte

Vier Konzepte erklären das meiste Verhalten:

**Die Projektdatei `u-boot.yaml`** ist die Quelle der Wahrheit. Sie hält
Projektnamen, aktivierte Dienste und Devcontainer-Einstellungen. Alles andere
– `compose.yaml`, `.env.example`, `.devcontainer/` – wird daraus abgeleitet.

**Verwaltete Bereiche.** In Dateien, die Ihnen gehören, ändert `u-boot` nur den
Abschnitt zwischen den Markierungen `BEGIN U-BOOT MANAGED BLOCK` und
`END U-BOOT MANAGED BLOCK`. Ihre eigenen Ergänzungen außerhalb dieser
Markierungen bleiben unangetastet.

**Nie stilles Überschreiben.** Trifft `u-boot` auf eine bestehende Datei, die
es nicht über einen verwalteten Bereich anfassen kann, bricht es ab und
erklärt, was zu tun ist – `--force` (Bereich ersetzen) oder `--backup`
(vorher sichern).

**Wiederholbar.** Dieselben Befehle noch einmal auszuführen ist ungefährlich:
Ist nichts zu tun, meldet `u-boot` das und ändert keine Datei.

### Vorschau statt Überraschung

Jeder schreibende Befehl kennt zwei Vorschau-Flags:

| Flag | Wirkung |
|---|---|
| `--dry-run` | zeigt die geplanten Änderungen, schreibt nichts |
| `--diff` | zeigt die geplanten Änderungen als Unified Diff |

```bash
u-boot remove otel --dry-run --diff
```

Beide lassen sich kombinieren. Nutzen Sie das, wann immer Sie unsicher sind: Es
kostet nichts und ist folgenlos.

Die Zusammenfassung eines `--dry-run` sagt, was geschehen würde: `init`, `add`,
`remove`, `generate` und `config set` melden `Would …` (zum Beispiel
`Would remove service "otel".` mit der Liste `Would change:`,
`Would generate changelog (CHANGELOG.md).` oder
`config: Would set project.name alt → neu.`). Ohne `--dry-run` steht dieselbe
Zeile in der Vergangenheitsform (`Removed service …`, `Generated …`,
`config: project.name alt → neu.`). Dasselbe gilt für `--dry-run --diff`; `--diff`
allein schreibt und meldet die Vergangenheitsform. Die Meldung „nichts zu tun“
(`… already up to date; no changes.`) und die `--json`-Ausgabe sind in beiden
Fällen gleich.

Auch die Logzeilen auf stderr, die ein Schreiben melden, behaupten im `--dry-run`
nichts als geschehen: `config set` und `generate` schreiben sie im Konjunktiv
(`config set: would update`, `generate: would create`, `generate: would update block`,
`generate devcontainer: would update allowlist`). Ohne `--dry-run` und bei `--diff`
allein stehen sie in der Vergangenheitsform (`config set: updated`, `generate: created`).
Die Pegel bleiben: Die Zeilen erscheinen auf der Standardstufe und entfallen
mit `--quiet`.

### Allgemeine Optionen

Diese Optionen gelten für jeden Befehl:

| Option | Wirkung |
|---|---|
| `--yes` | Rückfragen automatisch bejahen |
| `--no-interactive` | nie fragen; bei einer nötigen Rückfrage abbrechen |
| `--json` | maschinenlesbare Ausgabe |
| `--quiet` | nur Fehler: Erfolgs- und Zusammenfassungsmeldungen von `init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up` und `devcontainer down` entfallen (bei `down` hält `--quiet` die Compose-Ausgabe zurück und gibt sie nur aus, wenn `down` fehlschlägt; bei `devcontainer up` hält es den Fortschritt zurück und gibt ihn nur bei einem Fehler aus), ebenso die übrigen Hinweise und Warnungen (bei `remove` auch die Purge-Warnungen). Vier Ausnahmen bleiben sichtbar: (1) Sicherheitslockerungen des Sandbox-Profils, ein im Container aktivierter Login-Dienst (Feature `sshd`), die Lockerungen des Features `go` und die Warnung je Feature aus einer externen Quelle; in der Befehlsausgabe erzeugen nur `generate devcontainer` und `devcontainer up` (die Einbindung des Docker-Sockets in den Hilfs-Container) diese Warnungen, dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift, während die Meldung zu `nestedRuntime` ohne Sandbox-Profil entfällt; (2) Degradations-Fallbacks (Container ohne Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion; in `doctor` und in der Warnungsausgabe von `generate devcontainer`, `init` und `devcontainer up`, soweit der Host sie aus der Statusdatei des Containers lesen kann; ein nicht ermittelbarer Zustand ist kein Fallback und entfällt); (3) der Compose-Fortschritt von `up`; (4) die Pfad-Zusammenfassung, die `--force` vor dem Überschreiben ausgibt. `doctor --quiet` zeigt Fehler-Einträge, Degradations-Fallbacks und die Warn-Einträge zu einer Sicherheitseigenschaft des Sandbox-Profils, ohne Kopf- und Summary-Zeile (der Exit-Code trägt das Ergebnis; Einzelheiten in Abschnitt [4.7](#47-die-umgebung-prüfen)). Die Abfragebefehle `config get`, `config list`, `egress list`, `logs` und `template list` behalten ihre Ausgabe, ebenso die Vorschau von `--dry-run` (geplante Änderungen) und `--diff` (Diff): sie ist das Ergebnis des Aufrufs. Die Erfolgsmeldung eines tatsächlichen Schreibens entfällt dagegen auch bei `--diff` ohne `--dry-run`. Mit `--json` ändert `--quiet` nichts. Bei mehreren der Optionen `--quiet`, `--verbose` und `--debug` gilt die zuletzt auf der Kommandozeile angegebene, ohne Abweisung: `--quiet --verbose` zeigt zusätzliche Details und hebt `--quiet` auf (auch für Zusammenfassungen und Warnungen), `--verbose --quiet` zeigt nur Fehler. So lässt sich ein Alias wie `alias ub='u-boot --quiet'` mit `ub --verbose` überstimmen |
| `--verbose` | zusätzliche Details; gilt nur, wenn es die zuletzt angegebene der Optionen `--quiet`, `--verbose`, `--debug` ist |
| `--debug` | interne Diagnoseausgabe; gilt nur, wenn es die zuletzt angegebene der Optionen `--quiet`, `--verbose`, `--debug` ist (`--verbose` und `--debug` wirken gleich) |
| `-v`, `--version` | Version anzeigen |
| `-h`, `--help` | Hilfe zum Befehl |

`--yes` und `--no-interactive` schließen sich aus.

---

## 4. Aufgaben

### 4.1 Ein neues Projekt anlegen

#### Voraussetzung

Sie stehen in dem Verzeichnis, in dem das Projekt entstehen soll.

#### Vorgehen

1. Führen Sie aus:

   ```bash
   u-boot init mein-projekt
   ```

2. Prüfen Sie das Ergebnis:

   ```bash
   ls
   cat u-boot.yaml
   ```

#### Ergebnis

Es entstehen drei Verzeichnisse (`docker/`, `scripts/`, `docs/`), sechs Dateien
(`README.md`, `CHANGELOG.md`, `compose.yaml`, `.env.example`, `.gitignore`,
`u-boot.yaml`) und ein Git-Repository. Die Projektdatei beginnt so:

```yaml
schemaVersion: 1
project:
    name: mein-projekt
```

#### Hinweise

- Projektnamen bestehen aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnen
  mit einem Buchstaben und enden mit einem Buchstaben oder einer Ziffer (bis zu
  63 Zeichen). Großbuchstaben und Unterstriche werden normalisiert; ist der Name
  danach ungültig, bricht `u-boot` mit einer Erklärung ab (Exit-Code 10).
- `--no-git` unterdrückt das Anlegen des Git-Repositories.
- `--devcontainer` erzeugt zusätzlich `.devcontainer/devcontainer.json` und ein
  `Dockerfile` (siehe 4.9). `--devcontainer --sandbox` erzeugt das
  Sandbox-Profil (siehe 4.11).
- `--template <name>` rendert das Projekt aus einer Vorlage (siehe 4.12).

### 4.2 Ein bestehendes Projekt erneut initialisieren

#### Voraussetzung

Im Verzeichnis liegen bereits Dateien, möglicherweise aus einem früheren
`u-boot init`.

#### Vorgehen

1. Sehen Sie sich zuerst an, was passieren würde:

   ```bash
   u-boot init --dry-run
   ```

2. Entscheiden Sie sich für einen Weg:

   ```bash
   u-boot init --force            # verwaltete Bereiche ersetzen
   u-boot init --backup           # bestehende Dateien vorher sichern
   u-boot init --force --backup   # Bereiche ersetzen und vorher sichern
   ```

#### Ergebnis

Betroffene Pfade werden vor dem Schreiben aufgelistet. Mit `--backup` entstehen
Sicherungen nach dem Muster `<name>.bak`, `<name>.bak.1`, `<name>.bak.2` und so
fort – eine vorhandene Sicherung wird nie überschrieben.

#### Hinweise

- Erkennt `u-boot` mindestens drei Strukturelemente eines Projekts
  (`README.md`, `CHANGELOG.md`, `docs/`, `scripts/`, `docker/`,
  `.devcontainer/devcontainer.json`), ohne sicher zu sein, fragt es nach
  („Treat as an existing u-boot project?"). In Skripten und CI beantworten Sie
  das mit `--assume-existing` – `--yes` genügt dafür bewusst **nicht**. Mit
  `--no-interactive` entfällt die Erkennung ganz (frisches Anlegen).
- Ohne `--force` und ohne `--backup` bricht der Lauf bei einer vorhandenen Datei
  ab („file exists … pass --backup or --force"), statt zu raten.

### 4.3 Einen Dienst hinzufügen

#### Voraussetzung

Ein `u-boot`-Projekt (es gibt eine `u-boot.yaml`).

#### Vorgehen

1. Fügen Sie den Dienst hinzu:

   ```bash
   u-boot add postgres
   ```

   Verfügbar sind `postgres`, `keycloak` und `otel` (siehe Dienstekatalog im
   [Anhang](#dienstekatalog)).

2. Prüfen Sie die Änderung:

   ```bash
   cat u-boot.yaml
   ```

3. Tragen Sie in `.env` (Kopie von `.env.example`) eigene Werte für die
   Platzhalter `CHANGEME_…` ein.

#### Ergebnis

```text
Added service "postgres".

Changed:
  - u-boot.yaml
  - compose.yaml
  - .env.example
```

Die Projektdatei führt den Dienst jetzt:

```yaml
schemaVersion: 1
project:
  name: demo-app
services:
  postgres:
    enabled: true
```

#### Hinweise

- Der Befehl ist wiederholbar: Ein zweiter `u-boot add postgres` ändert nichts
  und meldet „already active; no changes".
- Ein unbekannter Dienst führt zu Exit-Code 10 mit der Liste der verfügbaren
  Dienste: `service not supported: "redis" is not in the built-in catalogue
  [postgres keycloak otel]`.
- Dienste lassen sich kombinieren, zum Beispiel `u-boot add keycloak` und
  `u-boot add otel` im selben Projekt; danach startet `u-boot up` alle.
- `--with-deps` installiert Dienste mit, die ein anderer Dienst voraussetzt, ohne
  Rückfrage. Aktuell setzt keiner der eingebauten Dienste einen anderen voraus;
  der Schalter ist daher ohne Wirkung und bricht nicht ab.
- `.env.example` bleibt eine Vorlage. Echte Zugangsdaten gehören nur in `.env`,
  das nicht ins Git-Repository gehört.

### 4.4 Einen Dienst entfernen

#### Voraussetzung

Der Dienst ist im Projekt aktiviert.

#### Vorgehen

```bash
u-boot remove postgres
```

#### Ergebnis

Der Dienst verschwindet aus `u-boot.yaml`, `compose.yaml` und `.env.example`.
Ist er schon entfernt, meldet `u-boot` das und ändert nichts.

#### Hinweise

- `--purge` entfernt zusätzlich das benannte Datenvolume des Dienstes (zum
  Beispiel das Postgres-Datenvolume). Das ist **destruktiv** und löst eine
  Sicherheitsabfrage aus; in nicht-interaktiven Läufen bricht der Befehl ohne
  `--yes` mit Exit-Code 10 ab („confirmation required"):

  ```bash
  u-boot remove postgres --purge --yes
  ```

- Wird das Volume noch von einem laufenden Container benutzt, bleibt es
  bestehen und Sie erhalten eine Warnung je Volume. Stoppen Sie dann den Stack
  mit `u-boot down` und entfernen Sie das Volume mit `docker volume rm <name>`.
  Die Zusammenfassung nennt die entfernten Volumes.

### 4.5 Die Umgebung starten und stoppen

#### Voraussetzung

Docker läuft, und im Projekt liegt eine `compose.yaml` mit mindestens einem
Dienst.

#### Vorgehen

```bash
u-boot up                # startet und wartet auf Stabilisierung
u-boot up --timeout 120  # wartet bis zu 120 Sekunden
u-boot up --timeout 0    # startet und kehrt sofort zurück
u-boot down              # stoppt die Umgebung
u-boot down --volumes    # stoppt und entfernt die Datenvolumes (nennt sie)
```

#### Ergebnis

`u-boot up` zeigt den Fortschritt (Image laden, Container erstellen und starten)
und danach eine Statustabelle je Dienst: Name, Container-Status, Port,
Healthcheck. Ohne Zeitlimit (`--timeout 0`) entfällt die Wartezeit; Sie bekommen
stattdessen einen Hinweis, den Status später selbst zu prüfen.

`u-boot down --volumes` meldet am Ende, welche Volumes entfernt wurden
(`Removed volumes: …`).

#### Hinweise

- Der Standardwert für `--timeout` ist 60 Sekunden. Negative Werte sind ein
  Nutzungsfehler (Exit-Code 2).
- Würde Compose einen Container neu erstellen (geändertes Image, geänderte
  Umgebung oder Volumes), warnt `u-boot up` **vorher** je Container: Nicht
  persistierte Daten im Container gehen dabei verloren.
- Scheitert `u-boot up` nach dem Start, trägt die Fehlerausgabe mit `--json` in
  `data.services` den Stand der bis dahin gestarteten Dienste.
- `--volumes` löscht Daten. Die Sicherheitsabfrage lässt sich in Skripten nur
  mit `--yes` überspringen; mit `--json` wird ohne `--yes` abgelehnt.
- Belegt ein anderer Prozess einen Port (zum Beispiel 5432), startet der
  Dienst nicht; siehe [Fehlerbehebung](#fehler-ein-port-ist-bereits-belegt).

### 4.6 Logs ansehen

#### Voraussetzung

Die Umgebung wurde mit `u-boot up` gestartet.

#### Vorgehen

```bash
u-boot logs                      # alle Dienste
u-boot logs postgres             # ein Dienst
u-boot logs postgres keycloak    # mehrere ausgewählte Dienste
u-boot logs --tail 100           # nur die letzten 100 Zeilen je Dienst
u-boot logs --follow             # laufend mitlesen, Abbruch mit Strg-C
u-boot logs --timestamps --no-log-prefix   # Zeitstempel statt Dienst-Präfix
u-boot logs --since 30m          # nur die letzten 30 Minuten
u-boot logs --since 2026-10-01T12:00:00Z --until 2026-10-01T13:00:00Z
```

#### Ergebnis

Die Log-Ausgabe der Compose-Dienste erscheint auf der Konsole, mit dem
Dienstnamen als Präfix (`postgres-1  | …`). `--follow` läuft, bis Sie mit
Strg-C abbrechen; das gilt als normales Ende (Exit-Code 0).

#### Hinweise

- `--since` und `--until` nehmen eine Dauer (`30m`, `1h`) oder einen
  Zeitstempel (`2026-10-01T12:00:00Z`). Ungültige Werte sind ein
  Nutzungsfehler (Exit-Code 2).
- `--follow` und `--json` schließen sich aus. Für maschinenlesbare Ausschnitte
  nutzen Sie `--tail <n> --json`.
- Ein Dienstname, den Compose nicht kennt, endet mit Exit-Code 12.

### 4.7 Die Umgebung prüfen

#### Vorgehen

```bash
u-boot doctor
u-boot doctor --strict   # jede Warnung gilt als Fehler
```

#### Ergebnis

Ein Bericht über 16 Prüfungen mit vier Zuständen: `✓` in Ordnung, `⚠` Warnung,
`✗` Fehler, `?` übersprungen oder nicht ermittelbar. Je konfigurierter
Sandbox-Fähigkeit (verschachteltes Podman oder Docker, Egress-Beschränkung)
kommt ein Eintrag zum Zustand beim Containerstart hinzu (siehe
[Fallbacks sehen](#fallbacks-sehen)). Der Bericht endet mit einer
Zusammenfassung:

```text
Diagnostic report for /projekt
──────────────────────────────────────
⚠  uboot.yaml.valid    u-boot.yaml not present — directory is not a u-boot project.
   → Run `u-boot init` to create one.
✓  fs.write-permissions           BaseDir is writable.
✓  docker.reachable               docker daemon is reachable.

Summary: 0 error, 1 warn, 15 ok
```

Jede Meldung mit Befund enthält eine Zeile mit `→`, die den nächsten Schritt
nennt.

Geprüft werden: Schreibrechte im Verzeichnis, `git`, Docker (Version ab 24.0,
Daemon erreichbar), Compose-Plugin (ab 2.20), `u-boot.yaml`, `compose.yaml`,
`.devcontainer/devcontainer.json` und `.devcontainer/Dockerfile` (wenn
vorhanden), Devcontainer-Features (Freigabeliste, im Sandbox-Profil das
Feature `docker-cli` als Fehler und Abweichungen zur Projektdatei),
Weiterleitungs-Ports, die drei Sandbox-Prüfungen (Laufzeit,
Zugangsdaten, Egress) und je konfigurierter Sandbox-Fähigkeit der gemeldete
Zustand beim Containerstart (`vfs`-Speicher, Daemon, Egress).

#### Hinweise

- Bei einem Fehler endet `doctor` mit Exit-Code 11; mit `--strict` auch bei
  einer Warnung. `--strict` ist für CI gedacht.
- `u-boot --quiet doctor` zeigt nur Fehler, Degradations-Fallbacks und die
  Warnungen zu einer Sicherheitseigenschaft des Sandbox-Profils: Klartext-Token
  oder -Schlüssel in einer Projektdatei, `https`-Clone ohne Token-Quelle,
  aktives Agenten-Feature ohne Referenz auf seine Umgebungsvariable,
  Egress-Restriktion
  zusammen mit `nestedRuntime: docker` und `egress.enabled` ohne Sandbox-Profil
  (die Beschränkung greift dann nicht). Alle anderen Warnungen entfallen, auch
  die Meldung, dass `nestedRuntime` ohne Sandbox-Profil wirkungslos ist; Kopf-
  und Summary-Zeile entfallen ebenfalls. Der Exit-Code ändert sich nicht: Mit
  `--strict` endet auch `--quiet` bei einer dieser Warnungen mit `11`, und die
  gezeigten Zeilen nennen den Grund. `--json` bleibt unverändert.
- Ein nicht ermittelbarer Zustand beim Containerstart erscheint nur im Bericht
  (`?`), nicht in `--json`, und verändert den Exit-Code nicht.
- Läuft `u-boot` im Container, werden die vier Host-Prüfungen übersprungen
  (`?`). Das ist kein Fehler. Die Einträge zum Containerstart liest `doctor`
  dort direkt aus dem Container, wenn es im Devcontainer selbst läuft; aus dem
  Runtime-Image ist der Zustand nicht ermittelbar.
- Im Container ist außerdem die Prüfung des Geräts `/dev/fuse` übersprungen
  (`devcontainer.sandbox.runtime`, Info, Exit-Code 0): das Gerät im Container ist
  nicht das des Container-Hosts. Die Prüfung des Schalters
  `kernel.apparmor_restrict_unprivileged_userns` im selben Eintrag gilt weiter:
  bei dem Wert 1 meldet er auch im Container einen Fehler (Exit-Code 11). Ein Devcontainer, der ohne das Gerät gestartet
  wurde, erscheint im Eintrag `devcontainer.sandbox.fallback.vfs`, soweit seine
  Statusdatei lesbar ist (sonst „nicht ermittelbar“); das Gerät des
  Container-Hosts prüft `u-boot doctor` auf dem Host.

### 4.8 Artefakte erzeugen und aktualisieren

#### Voraussetzung

Ein `u-boot`-Projekt.

#### Vorgehen

```bash
u-boot generate readme
u-boot generate changelog
u-boot generate env-example
u-boot generate devcontainer
```

#### Ergebnis

Das jeweilige Artefakt entsteht oder wird aktualisiert. Bei bestehenden Dateien
ändert `u-boot` nur den verwalteten Bereich; Ihre eigenen Ergänzungen bleiben
Byte für Byte erhalten.

#### Hinweise

- Wiederholte Läufe mit unverändertem Projekt ändern nichts.
- Haben Sie innerhalb eines verwalteten Bereichs von Hand editiert, meldet
  `u-boot` den Konflikt, statt Ihre Arbeit zu überschreiben. Bei `changelog`
  ergänzt `u-boot` höchstens eine fehlende Überschrift `## [Unreleased]`.
- Ein unbekannter Artefaktname endet mit Exit-Code 2 und nennt die erlaubten
  Namen.
- `generate devcontainer` hat zusätzliche Optionen (siehe 4.9 bis 4.11).
  `--dry-run` und `--diff` stehen bei allen vier Artefakten zur Verfügung.

### 4.9 Einen Devcontainer einrichten

#### Voraussetzung

Ein `u-boot`-Projekt. Zum Starten des Devcontainers brauchen Sie einen
Devcontainer-fähigen Editor (zum Beispiel VS Code), die `devcontainer`-CLI oder
Docker allein mit `u-boot devcontainer up` (siehe unten in diesem Abschnitt).

#### Vorgehen

1. Neues Projekt mit Devcontainer:

   ```bash
   u-boot init mein-projekt --devcontainer
   ```

   Bestehendes Projekt:

   ```bash
   u-boot generate devcontainer
   ```

2. Optional: eine andere Benutzer-ID im Container (zum Beispiel unter macOS mit
   Colima):

   ```bash
   u-boot config set devcontainer.user.uid 501
   u-boot generate devcontainer
   ```

#### Ergebnis

`.devcontainer/devcontainer.json` und `.devcontainer/Dockerfile` entstehen. Der
Container läuft als Benutzer `vscode` (nicht als root). Die Weiterleitungs-Ports
ergeben sich aus den aktivierten Diensten.

#### Hinweise

- Zulässige Benutzer-IDs: 1 bis 65535 (Standard 1000). Die ID 0 (root) wird
  abgelehnt (Exit-Code 10). Bei einem anderen Wert als 1000 erzeugt `u-boot` ein
  Build-Argument `USER_UID` und einen `usermod`-Schritt im Dockerfile; das gilt
  für beide Profile.
- Das Standardprofil bindet den Projektordner ein, aber keine
  Geheimnis-Verzeichnisse wie `~/.ssh` oder `~/.aws`. Das Docker-Socket des
  Hosts bindet es ein, wenn ein Feature es tut, zum Beispiel `docker-cli`
  (siehe *Sicherheit: das Docker-Socket* weiter unten).
- `u-boot doctor` meldet, wenn `devcontainer.json` und `u-boot.yaml`
  auseinanderlaufen (Warnung); `u-boot generate devcontainer` behebt das.
- Werkzeuge im Container (Node.js, Java, Go, …) fügen Sie als Features hinzu
  (siehe 4.10). Für autonome Agenten gibt es ein eigenes Profil (siehe 4.11).
- Den Container starten Sie mit einem Editor, mit der `devcontainer`-CLI oder ohne
  beides mit `u-boot devcontainer up` (siehe unten).

#### Den Devcontainer ohne Editor starten und stoppen

```bash
u-boot devcontainer up              # startet den Devcontainer des Projekts
u-boot devcontainer up --recreate   # entfernt den Container und legt ihn neu an
u-boot devcontainer down            # stoppt den Container
```

Auf dem Host genügen Docker (Client und lokale Engine). Ein Editor, Node.js und
die `devcontainer`-CLI sind nicht nötig. `up` liest `.devcontainer/devcontainer.json`
im aktuellen Verzeichnis und führt die Dev-Container-CLI in einem kurzlebigen
Hilfs-Container aus.

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

**Sicherheit: das Docker-Socket.** Der Hilfs-Container bindet für die Dauer des Aufrufs
das Engine-Socket des Hosts ein (`/var/run/docker.sock`). Das ist Root-Äquivalenz auf der
Engine. `u-boot devcontainer up` sagt das, bevor es den Hilfs-Container startet, auch
unter `--quiet`; mit `--json` steht es als Eintrag `warn` in `diagnostics`.

Der Devcontainer, den `u-boot` erzeugt, bindet das Socket ein, wenn ein Feature es tut:
Das Feature `docker-cli` bindet im Standardprofil das Docker-Socket des Hosts in den
Container ein; der Zugriff auf die Engine des Hosts ist Root-Äquivalenz auf der Engine.
Das Sandbox-Profil weist `docker-cli` ab (siehe 4.10 und 4.11). Ein Feature aus einer
externen Quelle kann das Socket ebenfalls einbinden; `u-boot` liest dessen Metadaten nicht
und warnt im Sandbox-Profil davor (siehe 4.10). Eine handgeschriebene
`.devcontainer/devcontainer.json` kann es einbinden, und ihr
`initializeCommand` läuft mit dem Zugriff des Hilfs-Containers. **Lesen Sie die Datei vor
dem ersten Start, wenn Dritte oder ein Agent im Projektordner schreiben dürfen.** Das
betrifft das Standardprofil: Im Sandbox-Profil liegt der Workspace in einem Volume, der
Projektordner des Hosts ist im Container nicht eingehängt.

Der Projektordner ist im Hilfs-Container unter demselben Pfad wie auf dem Host
schreibend eingehängt, weil ein `initializeCommand` in ihn schreiben darf. Das
Home-Verzeichnis und Geheimnis-Verzeichnisse wie `~/.ssh` sind nicht eingehängt. Ein
Projektpfad mit Komma, Anführungszeichen oder Zeilenumbruch wird abgewiesen (Exit-Code 10);
ein Ordner, den die Engine nicht sieht, ergibt Exit-Code 12.

**Zugangsdaten.** Der Befehl reicht nur Variablen durch, auf die `devcontainer.json`
mit `${localEnv:NAME}` verweist und die auf dem Host gesetzt sind, und zwar als Name,
nie als Wert auf der Kommandozeile von `u-boot`. Bei der Form `remoteEnv` steht der Wert
weder in einem Image noch in den Labels des Devcontainers, in einer Datei im Projektordner
oder in den Volumes des Sandbox-Profils (geprüft). Zwei Grenzen gehören dazu: `remoteEnv` ist
die für Zugangsdaten vorgesehene Form. Bei `containerEnv` steht der aufgelöste Wert in den
Metadaten (`Config.Env`) des Devcontainers, solange dieser besteht; für die Dauer eines
Lebenszyklus-Befehls steht er auf einer Kommandozeile im Hilfs-Container (`docker exec -e`
bei `remoteEnv`, `docker run -e` bei `containerEnv`), sichtbar für Prozesse mit Zugriff auf
den Daemon (zum Beispiel `docker top`). Die Maskierung ersetzt den Wert nur im Fortschritt
und in den Fehlermeldungen von `u-boot`, und nur, wenn er unverändert vorkommt und
mindestens vier Byte lang ist; kodierte Formen und kürzere Werte bleiben stehen. Andere
Quellen für `${localEnv:…}` (Features, weitere Konfigurationsdateien) liest der Befehl nicht;
Proxy-Variablen und Zugangsdaten für private Basis-Images reicht er nicht durch.

**Ausgabe.**

- `up` meldet nach dem Start den Container (die ersten 12 Zeichen der ID), den Benutzer
  im Container und den Workspace: `Devcontainer started: container <ID>, remote user <Benutzer>,
  workspace <Pfad>`. Lief der Container schon vor dem Aufruf und meldet die Dev-Container-CLI
  dieselbe Kennung, lautet die Zeile `Devcontainer already running: …` mit denselben Angaben;
  der Container lief schon, dieser Aufruf hat ihn nicht gestartet. Mit `--recreate` lautet sie
  immer `started`. Der
  Fortschritt (Bau des Hilfs-Images, Ausgabe der CLI) geht auf die Fehlerausgabe.
- `--quiet`: Die Erfolgsmeldung und der Fortschritt entfallen (bei einem Fehler wird der
  Fortschritt nachträglich ausgegeben); die Ausweisung des Sockets und die Fallbacks
  (siehe unten) bleiben sichtbar.
- `--json`: genau ein Umschlag auf der Standardausgabe, `command` ist `devcontainer`,
  `subcommand` ist `up` oder `down`. Der Fortschritt entfällt. **Der Umschlag eines
  erfolgreichen `up` trägt `status: warn`**, weil die Ausweisung des Sockets eine
  Warnung ist (der Exit-Code ist trotzdem `0`); prüfen Sie in Skripten den Exit-Code, nicht
  `status: ok`. Die Felder von `data` stehen im Abschnitt „JSON-Ausgabe“ im Anhang.
- `--dry-run` und `--diff` werden mit Exit-Code 2 abgewiesen, ebenso `u-boot devcontainer`
  ohne Unterbefehl: Der Befehl schreibt keine Projektdatei.

**Im Container arbeiten.** Ohne Editor öffnen Sie eine Shell im Container mit `docker exec`.
Container, Benutzer und Workspace stehen in der Erfolgsmeldung von `up`:

```bash
docker exec -it -u <Benutzer> -w <Workspace> <Container> bash
```

In dieser Shell gilt `remoteEnv` aus `devcontainer.json` nicht: Im Sandbox-Profil ist
`GIT_TOKEN` dort nicht gesetzt (beobachtet in einem Lauf, Linux `amd64`, Docker 29.9.1). Eine
Variable geben Sie mit `docker exec -e <NAME>` ausdrücklich mit; `docker` übernimmt dann den Wert
aus Ihrer Umgebung, und ist die Variable dort nicht gesetzt, fehlt sie auch im Container
(geprüft in einem Lauf mit `GIT_TOKEN`, Linux `amd64`, Docker 29.9.2):

```bash
docker exec -it -e GIT_TOKEN -u <Benutzer> -w <Workspace> <Container> bash
```

**Fallbacks nach dem Start.** Wie `doctor` und `generate devcontainer` weist `up` die
Degradations-Fallbacks des Sandbox-Containers aus (Container ohne Daemon, `vfs`-Storage,
Wegfall der Egress-Restriktion), mit dem Stand des letzten Starts. Der Befehl findet die
Statusdateien über die Fähigkeiten, die `u-boot.yaml` im Projektordner anfordert; **ohne
`u-boot.yaml` weist `up` nichts aus** und endet trotzdem mit Exit-Code 0. Meldet eine
Statusdatei `failed`, ist das kein Fallback, sondern Exit-Code 11 (Tabelle unten).

**Stoppen, Neuanlage und Volumes.** `u-boot devcontainer down` stoppt den Container und
entfernt weder ihn noch ein Volume; es braucht keinen Hilfs-Container und bindet kein
Socket ein. Ohne Container oder bei einem schon gestoppten Container endet es mit
Exit-Code 0. Es heißt anders als `u-boot down`: Das baut die Compose-Umgebung ab,
`devcontainer down` stoppt nur. Ein folgendes `u-boot devcontainer up` startet denselben
Container wieder. `up --recreate` entfernt den Container und legt ihn neu an; benannte
Volumes (zum Beispiel der Workspace des Sandbox-Profils) bleiben erhalten. Die Hinweise von
`u-boot` zu einer geänderten Egress-Liste nennen `devcontainer up --remove-existing-container`;
`u-boot devcontainer up --recreate` bewirkt dasselbe.

**Exit-Codes von `devcontainer up`:**

| Code | Fall |
|---|---|
| `2` | `u-boot devcontainer` ohne Unterbefehl; `--dry-run` oder `--diff` |
| `10` | keine `.devcontainer/devcontainer.json` im aktuellen Verzeichnis; Pfad mit Komma, Anführungszeichen oder Zeilenumbruch; ein `${localEnv:NAME}` mit einem Namen, der kein einfacher Variablenname ist |
| `11` | kein `docker` im `PATH` (so auch im Container-Image von `u-boot`), der Daemon antwortet nicht, die Engine ist nicht lokal (zum Beispiel `DOCKER_HOST=tcp://…`), oder ein Startscript des Containers bricht ab und meldet `failed` in seiner Statusdatei (unter `onUnavailable: fail`; das Startscript der Egress-Beschränkung auch unter `warn`, wenn der Benutzer `nobody` im Image fehlt, die nftables-Regeln sich nicht anwenden lassen oder der lokale DNS-Server nicht startet oder nicht unter diesem Benutzer läuft; siehe [Fallbacks sehen](#fallbacks-sehen)), auch wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet und die Dev-Container-CLI selbst gelingt (siehe „Erneutes `up` nach einem Fehler oder Abbruch“); der letzte Fall braucht die `u-boot.yaml` des Projekts, sonst lautet der Exit-Code `12`, und ein erneutes `up` meldet einen früheren Abbruch nicht; bei `devcontainer down` gelten die ersten drei Fälle |
| `12` | jeder andere Fehler beim Bau oder Start: der Bau des Images scheitert, ein Lebenszyklus-Befehl wie `postStartCommand` endet mit einem Fehler, die Engine sieht den Projektordner nicht; bei `down` ein Fehler von `docker stop` |
| `14` | Lesefehler im Dateisystem |

Die Fehlermeldung von `up` enthält bei Exit-Code 12 das Ende der Ausgabe der CLI.

**Was zurückbleibt.** Neben dem Hilfs-Image (`u-boot-devcontainer-cli:…`) legt die
Dev-Container-CLI je Projekt Images `vsc-<Ordnername>-…` an (gemessen auch `…-uid`); sie
sind nach `down` und nach `--recreate` weiter vorhanden (getestet; `--recreate` baut sie neu).
Aufräumen mit `docker image rm` oder `docker image prune`.

**Abbruch mit Strg-C.** Beobachtet: Die CLI im Hilfs-Container beendet sich auf das Signal
nicht früher; die Ursache ist nicht untersucht. `u-boot` wartet bis zu zehn Sekunden (gemessen:
10,1 Sekunden in der CI, 10,5 Sekunden lokal bis zur Rückkehr) und entfernt dann den
Hilfs-Container. Es bleibt kein Hilfs-Container zurück. Der Dev-Container, den die CLI
bis dahin angelegt hat, bleibt bestehen.

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
  `u-boot devcontainer up --recreate`. Getestet mit dem Startscript der Egress-Beschränkung unter
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

Ein Exit-Code 0 des zweiten `up` heißt also nicht, dass das `postCreateCommand` und das
`postStartCommand` durchgelaufen sind: Gegen einen laufenden Container führt `up` sie nicht erneut
aus. Abhilfe: `u-boot devcontainer up --recreate` entfernt den Container und legt ihn neu an;
benannte Volumes bleiben erhalten.

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

### 4.10 Devcontainer-Features nutzen

Features sind Bausteine des Devcontainers (zum Beispiel Node.js oder Go).
`u-boot` bringt einen Katalog mit und erlaubt externe Quellen nur nach
ausdrücklicher Freigabe.

#### Voraussetzung

Ein Projekt mit Devcontainer (siehe 4.9).

#### Vorgehen

1. Ein Feature aus dem Katalog aktivieren (`config set` akzeptiert `true`,
   `false`, `1`, `0`):

   ```bash
   u-boot config set devcontainer.features.node.enabled true
   ```

2. Optional: eine bestimmte Version festlegen:

   ```bash
   u-boot config set devcontainer.features.java.enabled true
   u-boot config set devcontainer.features.java.version 21
   ```

3. `devcontainer.json` neu erzeugen:

   ```bash
   u-boot generate devcontainer
   ```

#### Ergebnis

Das Feature steht im verwalteten Bereich von `.devcontainer/devcontainer.json`:

```jsonc
// BEGIN U-BOOT MANAGED BLOCK: init
{
  "name": "demo",
  "features": {
    "ghcr.io/devcontainers/features/node:1": {}
  },
  "remoteUser": "vscode"
}
// END U-BOOT MANAGED BLOCK: init
```

Mit `version 21` heißt der Eintrag `"ghcr.io/devcontainers/features/java:21": {}`.
In `u-boot.yaml` steht:

```yaml
devcontainer:
  enabled: true
  features:
    java:
      enabled: true
      version: "21"
```

#### Eingebauter Katalog

Ohne weitere Freigabe aktivierbar (`homebrew` und `claude-code` sind Drittanbieter-Features: `u-boot generate devcontainer` weist sie aus). Die Standardversion ist `1`, bei `homebrew` `2`.

| Name | Inhalt |
|---|---|
| `git` | Git |
| `docker-cli` | Docker-CLI („docker-outside-of-docker“): bindet das Docker-Socket des Hosts in den Container ein (Root-Äquivalenz auf der Engine). Im Sandbox-Profil nicht wählbar: `u-boot generate devcontainer` endet dann mit Exit-Code 10 (siehe [Features im Sandbox-Profil](#features-im-sandbox-profil)) |
| `node` | Node.js |
| `java` | Java und SDKMAN |
| `go` | Go-Toolchain; das Feature verlangt laut seinen Metadaten die Fähigkeit `SYS_PTRACE` und `seccomp=unconfined`, die `u-boot generate devcontainer` im Sandbox-Profil als Sicherheitslockerung ausweist (siehe [Features im Sandbox-Profil](#features-im-sandbox-profil)) |
| `cpp` | C++-Toolchain |
| `kubectl-helm` | kubectl, helm und minikube |
| `postgres-client` | PostgreSQL-Client |
| `python` | Python (das Betriebssystem-Python; `pip install` außerhalb eines virtuellen Environments lehnt es nach PEP 668 ab: `python -m venv .venv` verwenden) |
| `sdkman` | SDKMAN (das offizielle `java`-Feature installiert es; zusammen mit `java` erscheint es einmal) |
| `sshd` | SSH-Server (offizielles Feature `sshd`, Login-Dienst auf Port 2222 mit Passwort- und Schlüssel-Login, Root-Login erlaubt; siehe [Login-Dienst `sshd`](#login-dienst-sshd-optional)); keine Optionen (`gatewayPorts`, `version` kennt der Katalog nicht) |
| `homebrew` | Homebrew; **Drittanbieter-Feature** (`meaningful-ooo`, Version `2`): `u-boot generate devcontainer` weist es aus, sein Installationsskript läuft beim Image-Build |
| `claude-code` | Claude Code CLI; **Drittanbieter-Feature** (`anthropics`, Version `1`): `u-boot generate devcontainer` weist es aus, sein Installationsskript läuft beim Image-Build; braucht das Feature `node` (siehe [Claude Code](#claude-code-optional)); keine Optionen |
| `ollama` | Ollama (lokale Sprachmodelle, nur CPU); wirkt nur im Sandbox-Profil. Kein OCI-Feature: Das erzeugte Dockerfile installiert ein gepinntes Release-Archiv (SHA-256 geprüft, vom Upstream nicht signiert; `u-boot generate devcontainer` weist das aus), dazu ein benanntes Volume für die Modelle; keine Optionen, `version` nur zusammen mit `source` (siehe [Ollama](#ollama-optional)) |

#### Externe Feature-Quellen freigeben

Jedes Feature außerhalb des Katalogs braucht eine Freigabe in
`devcontainer.featureSources.allow`. `--yes` ersetzt sie bewusst nicht.

1. Tragen Sie die URL in die Freigabeliste ein. Alle drei Wege sind
   gleichwertig:

   ```bash
   # bei init
   u-boot init --devcontainer \
     --allow-external-feature-sources https://ghcr.io/orgX/features/custom-rust
   # vor jedem generate (kann mehrfach angegeben werden)
   u-boot generate devcontainer \
     --allow-external-feature-sources https://ghcr.io/orgX/features/custom-rust
   # direkt per config set; mehrere URLs kommagetrennt
   u-boot config set devcontainer.featureSources.allow \
     https://ghcr.io/orgX/features/custom-rust
   ```

2. Aktivieren Sie das Feature:

   ```bash
   u-boot config set devcontainer.features.custom-rust.source \
     https://ghcr.io/orgX/features/custom-rust
   u-boot config set devcontainer.features.custom-rust.enabled true
   u-boot generate devcontainer
   ```

Danach steht in `u-boot.yaml`:

```yaml
devcontainer:
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust
```

Die Freigabeliste wird bei jeder Angabe ergänzt (nicht ersetzt); doppelte
Einträge werden still ausgelassen. Mehrere Aufrufe in einem `config set` sind
atomar (siehe Abschnitt 5).

#### Hinweise

- Die Quelle muss **Zeichen für Zeichen** in der Freigabeliste stehen. Ein
  abschließender Schrägstrich (`https://x/y` gegenüber `https://x/y/`) und die
  Groß-/Kleinschreibung des Hosts sind entscheidend.
- Ohne Freigabe endet `config set …source` mit Exit-Code 10:
  `external source "…" is not in devcontainer.featureSources.allow`. Die
  Fehlermeldung nennt den passenden `config set`-Befehl.
- `--allow-external-feature-sources` gibt es nur bei `init --devcontainer`,
  `generate devcontainer` und `config set …featureSources.allow`; bei anderen
  Befehlen wird es abgelehnt.
- Es erscheinen nur Features mit `enabled: true`, sortiert nach Quelle. Dadurch
  ist die Ausgabe immer gleich, und ein zweiter `generate devcontainer` ändert
  nichts.
- Eigene Features aus einem lokalen Verzeichnis (statt einer URL) werden nicht
  unterstützt.
- Im Sandbox-Profil weist `u-boot generate devcontainer` jedes aktivierte
  Feature aus einer externen Quelle mit einer Warnung aus (siehe
  [Features im Sandbox-Profil](#features-im-sandbox-profil)).

#### Features im Sandbox-Profil

Das Sandbox-Profil (Abschnitt 4.11) bindet das Docker-Socket des Hosts nicht
ein. Für die Features gilt:

- **`docker-cli` ist im Sandbox-Profil nicht wählbar.** Ist das Profil
  `sandbox` (in `u-boot.yaml` oder mit `--sandbox`) und das Feature aktiv, endet
  `u-boot generate devcontainer` mit Exit-Code 10 und schreibt nichts, auch nicht
  mit `--dry-run`; `--sandbox` wird dann nicht in `u-boot.yaml` übernommen.
  `u-boot doctor` meldet dasselbe als Fehler. Die Meldung nennt die Abhilfe in
  dieser Reihenfolge: das Feature abschalten, wer Docker im Container braucht
  zusätzlich `nestedRuntime: docker` setzen, dann neu erzeugen und einen
  vorhandenen Container neu anlegen.

  ```bash
  u-boot config set devcontainer.features.docker-cli.enabled false
  ```

  Brauchen Sie Docker im Container, setzen Sie zusätzlich (das Profil bekommt
  dann `--privileged`, ein eigener Docker-Daemon, Abschnitt 4.11):

  ```bash
  u-boot config set devcontainer.sandbox.nestedRuntime docker
  ```

  Danach erzeugen Sie den Devcontainer neu und legen einen vorhandenen
  Container neu an:

  ```bash
  u-boot generate devcontainer
  u-boot devcontainer up --recreate   # legt einen vorhandenen Container neu an
  ```

  `nestedRuntime: docker` allein hebt den Fehler nicht auf. `config set` bleibt
  auch in einem solchen Projekt möglich. Im Standardprofil bleibt `docker-cli`
  wählbar.
- **`go`:** Das Feature verlangt laut seinen Metadaten die Fähigkeit
  `SYS_PTRACE` und `seccomp=unconfined`. Im Sandbox-Profil weist
  `u-boot generate devcontainer` beide als Sicherheitslockerung aus, je eine
  Warnung, auch unter `--quiet`; `u-boot` schreibt sie nicht in
  `devcontainer.json`. Ohne Sandbox-Profil weist `u-boot` sie nicht aus.
- **Features aus einer externen Quelle** (siehe oben, *Externe Feature-Quellen
  freigeben*): Im Sandbox-Profil weist `u-boot generate devcontainer` jedes
  aktivierte Feature mit einer Warnung aus, auch unter `--quiet`. `u-boot` liest
  die Metadaten eines solchen Features nicht, und das Profil sagt für dieses
  Feature keine seiner Eigenschaften zu, etwa den Verzicht auf das Docker-Socket.

#### Prüfungen durch `u-boot doctor`

| Schweregrad | Auslöser | Abhilfe |
|---|---|---|
| Fehler | `source` ist gesetzt, steht aber nicht in `featureSources.allow` | URL freigeben (siehe oben) |
| Fehler | Das Profil ist `sandbox`, und das Feature `docker-cli` ist aktiv | `docker-cli` abschalten (wer Docker im Container braucht, setzt zusätzlich `nestedRuntime: docker`), `u-boot generate devcontainer` ausführen und den Container neu anlegen (siehe [Features im Sandbox-Profil](#features-im-sandbox-profil)) |
| Warnung | Name steht nicht im Katalog und hat keine `source`; er wird beim Erzeugen übersprungen | `source` setzen oder den Eintrag entfernen |
| Warnung | `enabled:` fehlt bei einem Feature-Eintrag | `enabled` setzen |
| Warnung | Feature ist aktiv, fehlt aber in `devcontainer.json` (oder die Datei fehlt) | `u-boot generate devcontainer` |
| Warnung | Feature ist ausgeschaltet, steht aber noch in `devcontainer.json` | `u-boot generate devcontainer` |
| Warnung | `devcontainer.json` enthält ein Feature, das in `u-boot.yaml` fehlt | Eintrag in `u-boot.yaml` ergänzen oder aus der JSON-Datei entfernen |

Treffen mehrere Befunde zu, zeigt `doctor` den schwersten. Ohne konfigurierte
Features oder bei nicht lesbarer `u-boot.yaml` oder `devcontainer.json` schweigt
die Prüfung; dafür sind die Prüfungen `uboot.yaml.valid` und
`devcontainer.json.valid` zuständig.

#### Beispiel: vollständiger Ablauf

```bash
mkdir myproj && cd myproj
u-boot init --devcontainer

u-boot config set devcontainer.features.git.enabled true
u-boot config set devcontainer.features.node.enabled true
u-boot config set devcontainer.features.java.enabled true
u-boot config set devcontainer.features.java.version 21

u-boot config set devcontainer.featureSources.allow \
  https://ghcr.io/orgX/features/custom-rust
u-boot config set devcontainer.features.custom-rust.source \
  https://ghcr.io/orgX/features/custom-rust
u-boot config set devcontainer.features.custom-rust.enabled true

u-boot generate devcontainer
u-boot doctor
```

Ergebnis: `devcontainer.json` mit vier Features, `doctor` ohne Befund.

### 4.11 Einen Sandbox-Devcontainer für autonome Agenten einrichten

Das Sandbox-Profil erzeugt einen Devcontainer, in dem ein autonomer Agent ohne
Rückfrage arbeiten kann. Es ist **Schadensbegrenzung, keine harte
Isolationsgrenze**. `u-boot` startet keinen Agenten und setzt keinen
Berechtigungsmodus.

#### Voraussetzung

Ein `u-boot`-Projekt, am besten mit einem Git-Remote `origin` (oder einer
konfigurierten Clone-Quelle, siehe unten).

#### Vorgehen

1. Profil aktivieren:

   ```bash
   u-boot init mein-projekt --devcontainer --sandbox     # neues Projekt
   u-boot generate devcontainer --sandbox                # bestehendes Projekt
   u-boot config set devcontainer.profile sandbox        # alternativ, danach generate
   ```

2. Optional: weitere Einstellungen (siehe die folgenden Abschnitte) und danach
   neu erzeugen:

   ```bash
   u-boot generate devcontainer
   ```

3. Zugangsdaten für `https`-Remotes geben Sie zur Laufzeit mit, nicht in Dateien:
   Setzen Sie vor dem Start des Containers die Umgebungsvariable `GIT_TOKEN` auf
   dem Host. Mit dem Feature `claude-code` gilt dasselbe für `ANTHROPIC_API_KEY`
   (siehe [Zugangsdaten](#zugangsdaten-git-und-agent)).

#### Ergebnis

`--sandbox` setzt `devcontainer.profile: sandbox` in `u-boot.yaml`.

| Eigenschaft | Sandbox-Profil | Standardprofil |
|---|---|---|
| Workspace | benanntes Volume `<projekt>-workspace-${devcontainerId}`; das Repository wird im Container geklont | Bind-Mount des Projektordners |
| Docker-Socket des Hosts (im erzeugten Devcontainer; der Hilfs-Container von `u-boot devcontainer up` bindet es ein, siehe 4.9) | nicht eingebunden; das Feature `docker-cli` wird abgewiesen (Exit-Code 10), für ein Feature aus einer externen Quelle sagt das Profil nichts zu (Warnung) | eingebunden, wenn ein Feature es tut, zum Beispiel `docker-cli` |
| `--privileged` und zusätzliche Fähigkeiten | nein (außer mit verschachteltem Podman oder Docker, Docker verlangt `--privileged`, und außer den Lockerungen, die ein Feature wie `go` mitbringt und `u-boot` ausweist) | nein (was ein Feature wie `go` mitbringt, weist `u-boot` hier nicht aus) |
| Geheimnis-Verzeichnisse des Hosts (`~/.ssh`, `~/.aws`, …) | nicht eingebunden | nicht eingebunden |
| Benutzer | `vscode`, nicht root | `vscode`, nicht root |

#### Hinweise

- `--sandbox` ohne `--devcontainer` (bei `init`) oder bei einem anderen Artefakt
  als `devcontainer` (bei `generate`) ist ein Fehler.
- Ist das Feature `docker-cli` aktiv, endet `u-boot generate devcontainer` im
  Sandbox-Profil mit Exit-Code 10 und schreibt nichts; die Meldung nennt die
  Abhilfe (siehe [Features im Sandbox-Profil](#features-im-sandbox-profil)).

#### Clone-Quelle

Der Clone läuft beim ersten Start, nur wenn der Workspace noch leer ist. Quelle
ist `devcontainer.sandbox.repository`, sonst die Adresse des Remotes `origin`.
Um im Container ein **anderes Repository** zu bearbeiten und daraus zu pullen und
zu pushen:

```bash
u-boot config set devcontainer.sandbox.repository git@github.com:org/anderes-repo.git
u-boot generate devcontainer
```

- Weder `repository` noch `origin` vorhanden (zum Beispiel direkt nach `init`):
  Es gibt keinen Clone-Schritt, `u-boot` warnt. Nach `git remote add origin …`
  (oder `config set …repository …`) ergänzt `u-boot generate devcontainer` den
  Schritt.
- Adressen mit Zugangsdaten oder unsicheren Zeichen (`https://token@…`,
  `user:pw@`, `;`, `$(…)`) werden abgelehnt (Exit-Code 10); es wird nichts
  geschrieben. Erlaubt sind `git@host:pfad`, `ssh://git@host/pfad` und `https://`
  ohne Benutzerangabe.
- Git-Worktrees und Submodule (`.git` als Datei) gelten als „kein Remote".

#### Verschachteltes Podman (optional)

```bash
u-boot config set devcontainer.sandbox.nestedRuntime podman
u-boot generate devcontainer
```

Im Container läuft dann rootless Podman mit einem `docker`-Alias (ein Symlink
auf `podman`), ohne Zugriff auf einen Host-Socket.

**Herkunft und Version.** Im Container liegt Podman 6.1.4, nicht das Debian-Paket:
das statisch gelinkte Release-Archiv des Community-Projekts `podman-static` mit
`crun`, `conmon`, `netavark`, `aardvark-dns`, `passt` und `fuse-overlayfs` unter
`/usr/local`. Version, Prüfsumme je Architektur und der Fingerprint des
Signaturschlüssels stehen im erzeugten Dockerfile; der Bau bricht vor dem
Entpacken ab, wenn die Prüfsumme nicht stimmt, der Schlüssel nicht importierbar
ist oder die Signatur nicht von diesem Schlüssel stammt. Aus Debian kommen nur
`uidmap` und `ca-certificates`. Das Archiv ist ein Binärpaket, kein Skript;
Vertrauen in die Quelle heißt Vertrauen in dieses Projekt und seine
Build-Pipeline. Die Version ist festgelegt: eine neuere Podman-Version kommt mit
einer neuen `u-boot`-Version, nicht von selbst. Archiv und Schlüssel lädt der Bau
des Images, nicht der laufende Container: Der Bau braucht `github.com`,
`objects.githubusercontent.com` und `keys.openpgp.org` offen (er läuft auf dem
Container-Host, nicht hinter der Egress-Restriktion); die Default-Allowlist des
laufenden Containers braucht keinen dieser Hosts zusätzlich. Bestehende Projekte mit `nestedRuntime: podman` bekommen die Version
nach `u-boot generate devcontainer` und einem Neubau des Containers.

Unter Docker sind dafür diese Lockerungen nötig; **`u-boot` weist jede einzeln
als Warnung aus**:

```text
--cap-add=SYS_ADMIN
--security-opt=seccomp=unconfined
--security-opt=apparmor=unconfined
--security-opt=systempaths=unconfined
--device=/dev/fuse
--device=/dev/net/tun
```

`/dev/net/tun` braucht das Netzwerk der verschachtelten Container (`pasta`);
ohne das Device scheitert jedes `RUN` eines `podman build`. `/dev/fuse` dient dem
Speicher der verschachtelten Container. Beide Devices müssen auf dem
Container-Host vorhanden sein; fehlt eines, lehnt Docker den Start ab (gemessen: `error gathering device information while adding custom device … no such file or directory`), das Startscript läuft dann nicht. `u-boot doctor` prüft `/dev/net/tun` nicht und meldet ein fehlendes `/dev/fuse` auf einem Linux-Host, nicht aus einem Container, als Fehler (Exit-Code 11). Abhilfe: die Kernelmodule auf dem Host laden (`modprobe fuse`, `modprobe tun`; `tun` nicht gemessen) oder `nestedRuntime: none` wählen. Bestehende Projekte mit `nestedRuntime: podman` brauchen das Device nach `u-boot generate devcontainer` ebenfalls.

Das ist nahe an `--privileged`; der Schutz ist damit deutlich schwächer. Das
Standardprofil (`nestedRuntime: none`) bleibt ungelockert.

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

**Gemessen** (Podman 6.1.4; Docker 29.5.2 auf Colima 0.10.3, arm64; in der VM
`kernel.apparmor_restrict_unprivileged_userns=0`, auf dem Container-Host `/dev/net/tun`): `podman build` mit `FROM`,
`RUN` und `COPY` baut kalt, aus dem Layer-Cache und mit `--no-cache`; `docker
build` über den Alias ebenso; ein `RUN` mit Netz zu einem erlaubten Host
(`apt-get update`) läuft auch hinter der Egress-Restriktion mit der
Default-Allowlist (gemessen mit `debian:stable-slim`); für `alpine:3` hinter der Restriktion sind Pull und Build gemessen, ein `RUN` mit Netz nicht. Der Blob-Host von Docker Hub kann sich ändern; ein nicht gelisteter Name löst hinter der Restriktion nicht auf. **Nicht
gemessen:** andere Hosts, andere Architekturen (amd64), CI, Docker Desktop, Podman als
Host-Engine. **Nicht zugesagt:** `buildx`, `docker compose` und `# syntax=`-Frontends.

Was beim Containerstart passiert, wenn Fähigkeiten fehlen, steuert
`devcontainer.sandbox.onUnavailable` (`warn` ist der Standard, `fail` bricht ab):

| Zustand | `warn` | `fail` |
|---|---|---|
| `/dev/fuse` fehlt im Container (ohne das Device gestartet) | Rückfall auf `vfs`-Speicher, Warnung | Abbruch mit Exit-Code 11 |
| Verschachtelte User-Namespaces blockiert (Seccomp oder AppArmor), nur `podman` | Abbruch mit Exit-Code 11 | Abbruch mit Exit-Code 11 |
| Docker-Daemon wird nicht betriebsbereit (`docker`, zum Beispiel Start ohne `--privileged`) | Container ohne Daemon, Warnung | Abbruch mit Exit-Code 11 |
| Egress-Beschränkung: Benutzer `nobody` fehlt, die nftables-Regeln lassen sich nicht anwenden, lokaler DNS-Server startet nicht oder läuft unter einem anderen Benutzer | Abbruch mit Exit-Code 11 | Abbruch mit Exit-Code 11 |

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

Ein fehlendes `/dev/fuse` auf dem Container-Host ist keine Stufe dieser Tabelle:
Docker lehnt den Start des Containers ab. `u-boot doctor` meldet es auf Linux-Hosts
als Fehler (Exit-Code 11), unabhängig von `onUnavailable`. Unter macOS mit Colima
läuft die Engine in einer VM, die der Host nicht beurteilen kann; dort lehnt Docker
den Start ab, die Prüfung von `doctor` zum Gerät meldet nichts, und das
Startscript im Container prüft die User-Namespaces. Läuft `u-boot doctor` in einem
Container (Runtime-Image, Devcontainer), prüft es das Gerät nicht: das `/dev/fuse`
dort gehört dem Container, nicht dem Container-Host. Der Eintrag
`devcontainer.sandbox.runtime` ist dann Info (Exit-Code 0), es sei denn, die
Prüfung des Schalters `kernel.apparmor_restrict_unprivileged_userns` davor
meldet einen Fehler (Wert 1, Exit-Code 11); ein ohne das Gerät
gestarteter Devcontainer erscheint im Eintrag `devcontainer.sandbox.fallback.vfs`,
soweit seine Statusdatei lesbar ist (sonst „nicht ermittelbar“), das Gerät des Container-Hosts prüft `u-boot doctor` auf dem Host.

#### Fallbacks sehen

Jedes Startscript schreibt bei jedem Lauf eine kleine Statusdatei im Container:
`/var/lib/u-boot/status/vfs` (`nestedRuntime: podman`), `daemon`
(`nestedRuntime: docker`) und `egress` (Egress-Beschränkung). Sie hat drei
Zeilen: das Format (`u-boot-status 1`), das Ergebnis (`ok`, `fallback` oder
`failed`) und den Zeitpunkt (UTC). Exit-Codes und Meldungen der Startscripte
bleiben, wie sie waren.

`u-boot doctor` liest die Dateien vom Host über die Container-Engine (nur
`docker ps` und `docker cp`, nie `docker exec`) und zeigt je konfigurierter
Fähigkeit einen Eintrag:

| Eintrag | Fähigkeit |
|---|---|
| `devcontainer.sandbox.fallback.vfs` | `nestedRuntime: podman`: Rückfall auf `vfs`-Speicher |
| `devcontainer.sandbox.fallback.daemon` | `nestedRuntime: docker`: Container ohne Daemon |
| `devcontainer.sandbox.fallback.egress` | Egress-Beschränkung: Wegfall der Beschränkung |

| Gemeldet | Anzeige in `doctor` |
|---|---|
| `ok` | in Ordnung |
| `fallback` | Warnung; bleibt unter `--quiet` sichtbar; mit `--strict` Exit-Code 11 |
| `failed` (das Startscript endete mit Exit-Code 11: `onUnavailable: fail`, bei der Egress-Beschränkung auch unter `warn`, wenn der Benutzer `nobody` fehlt, die nftables-Regeln sich nicht anwenden lassen oder der lokale DNS-Server nicht startet oder nicht unter ihm läuft) | Fehler, Exit-Code 11 |
| nicht ermittelbar | Hinweiszeile (`?`) mit dem Grund; kein Fallback; unter `--quiet` stumm |
| ungültige Datei | Warnung |

**Nicht ermittelbar** heißt: Es gibt keinen Container zu diesem Projektordner
(das Dev-Container-Werkzeug muss ihn mit dem Ordner dieses Projekts erzeugt
haben), die Engine ist nicht erreichbar oder zeigt auf einen anderen Rechner
(`DOCKER_HOST` oder Docker-Kontext), `u-boot` läuft aus dem Runtime-Image statt
im Devcontainer, oder der Container hat noch keine Statusdatei. **Bestehende
Container ohne Statusdateien** (vor dieser Funktion gebaut oder nie gestartet)
brauchen `u-boot generate devcontainer` und einen Neubau; bis zum ersten Start
danach bleibt der Zustand nicht ermittelbar.

Die Datei beschreibt den letzten Start: `daemon` und `egress` schreibt jeder
Start neu, den das Script erreicht (bei blockierten User-Namespaces entsteht
keine `vfs`-Datei), `vfs` einmal je
Container. Ein gestoppter Container wird gelesen; der
Eintrag nennt den Zeitpunkt und „container stopped“.

Nach `u-boot generate devcontainer` fragt `u-boot` denselben Weg ab und warnt,
wenn ein **bestehender** Container zuletzt `fallback` gemeldet hat (`init
--devcontainer` fragt auf demselben Weg, schreibt aber weder eine verschachtelte
Runtime noch eine Egress-Beschränkung und hat daher nichts abzufragen).
Die Warnung nennt den Stand des letzten Starts des bestehenden Containers, nicht
den der neu erzeugten Dateien, und bleibt unter `--quiet` sichtbar. Ohne Engine,
ohne Container, ohne Datei, bei einer ungültigen Datei und bei `ok` oder
`failed` bleibt die Abfrage still, und bei `--dry-run` findet sie nicht statt.
`up` fragt nicht.

**Grenze:** Die Datei stammt aus einem Container, in dem ein Agent läuft, und ist
dort mit Root-Rechten schreibbar; gemessen ist außerdem, dass ein Prozess ohne
`sudo` sie ersetzen kann (die Datei gehört `root`, aber das Verzeichnis
gehört dem Benutzer `vscode`; ein Lauf auf Colima, arm64). Sie ist ein Hinweis,
kein Beleg: **„Nicht gemeldet“ heißt nicht „nicht aufgetreten“.** `u-boot` gibt nie den Inhalt der
Datei aus, nur die geprüften Felder; eine Datei mit anderem Format, mit mehr als
512 Byte oder die keine reguläre Datei ist, gilt als ungültig.


#### Eigener Docker-Daemon (optional)

```bash
u-boot config set devcontainer.sandbox.nestedRuntime docker
u-boot generate devcontainer
```

Im Container läuft dann ein **eigener Docker-Daemon** (Pakete `docker-ce`,
`docker-ce-cli`, `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin` aus dem
Docker-APT-Repository; der Image-Build prüft den Schlüssel-Fingerprint und braucht
Netz zu `download.docker.com`) mit BuildKit, Buildx und
Compose. Die `docker`-CLI spricht nur mit diesem Daemon; der Socket ist dem
Benutzer `vscode` über die Gruppe `docker` zugänglich, `sudo` ist nicht nötig.
Der Daemon startet beim Containerstart und liegt in einem Volume je Instanz
(`<projekt>-docker-${devcontainerId}`). Das brauchen Repos, deren Prüfungen
`docker buildx build -o type=cacheonly`, `RUN --network=none`, `docker compose`
oder `docker run --privileged docker:dind` nutzen.

**Warnung:** `docker` verlangt `--privileged`. `u-boot` weist das als Warnung
aus. Das Profil bietet damit **keine Isolationsgrenze gegenüber dem
Container-Host** mehr; ein Ausbruch trifft unter Linux den Host und unter macOS
die Colima-VM. Das Profil begrenzt weiterhin Workspace,
Geheimnisse und Docker-Socket des Hosts.

- Jede Instanz hat einen eigenen Daemon: gleichnamige Images, Container, Netze
  und Volumes mehrerer Instanzen kollidieren nicht. Dafür hält jede Instanz
  ihre Images selbst (mehrere GB) und startet kalt.
- Wird der Daemon nicht innerhalb von 60 Sekunden bereit (zum Beispiel ohne
  `--privileged`), gilt `devcontainer.sandbox.onUnavailable`: Warnung und
  Container ohne Daemon (`warn`) oder Exit-Code 11 (`fail`). Podman ersetzt den
  Daemon nie. `u-boot doctor` weist das aus der Statusdatei aus (siehe
  [Fallbacks sehen](#fallbacks-sehen)).
- Das Standardprofil und `podman` bleiben unverändert.
- Unter Podman als Host-Engine ist `docker` nicht zugesagt.
- Der Host `production.cloudfront.docker.com` der Default-Allowlist gilt auch hier, ist aber für `nestedRuntime: docker` nicht gemessen (gemessen ist der Pull mit `podman`).

#### Internet-Beschränkung (Egress)

```bash
u-boot config set devcontainer.sandbox.egress.enabled true
u-boot config set devcontainer.sandbox.egress.allow api.example.org,registry.npmjs.org
u-boot generate devcontainer
```

Beim Containerstart startet ein lokaler DNS-Server, den der Container als
Standard-Resolver benutzt und der **nur erlaubte Namen** auflöst; auch die
DNS-Server des Netzes (die Upstream-Resolver) beantworten einem Prozess ohne
`sudo` keinen anderen Namen (siehe unten). Eine Firewall lässt nur Loopback,
bestehende Verbindungen, DNS vom lokalen DNS-Server zu den Upstream-Resolvern und
die Adressen der erlaubten Namen durch; alles andere (IPv6 ganz) wird verworfen. Das
ist eine Leitplanke, keine Sandbox-Grenze. Der Container braucht dafür die
Fähigkeit `NET_ADMIN`, die `u-boot` als Lockerung ausweist.

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

Das schließt die Regel nicht:

- Alles mit `sudo`: Wer `sudo` nutzt, hebt die Regeln auf oder fragt als `nobody`
  (der Docker-Test fragt so den Upstream-Resolver und bekommt eine Antwort).
- Verkehr zu den Adressen erlaubter Namen auf jedem Port und Protokoll: Das ist
  die Liste selbst.
- Anfragen nach Unternamen erlaubter Namen: Der lokale DNS-Server leitet sie an den
  Upstream-Resolver weiter. Das ist ein Weg, Daten zum zuständigen DNS-Server des
  erlaubten Namens zu schicken (im Laboraufbau gemessen: höchstens 3 650 bis 5 640
  Anfragen je Sekunde, ohne Test, der es pinnt); bei den Standard-Namen betreibt ihn
  kein Dritter, bei von Ihnen ergänzten Namen kann es anders sein (nicht gemessen).
  Der Betreiber des Upstream-Resolvers sieht diese Anfragen.

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

Standardmäßig erlaubt sind:

| Quelle | Hosts |
|---|---|
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` |
| mit `nestedRuntime: podman` oder `docker` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `production.cloudfront.docker.com`, `ghcr.io` |
| mit `nestedRuntime: docker` | zusätzlich `download.docker.com` |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Feature `python` | `pypi.org`, `files.pythonhosted.org` |
| Feature `sdkman` | `api.sdkman.io`, `broker.sdkman.io`, `repo.maven.apache.org`, `repo1.maven.org` |
| Feature `sshd` | keine zusätzlichen Hosts |
| Feature `homebrew` | `formulae.brew.sh`, `ghcr.io`, `pkg-containers.githubusercontent.com` |
| Feature `claude-code` | `registry.npmjs.org`, `api.anthropic.com` (Hosts der Anmeldung nicht: siehe [Claude Code](#claude-code-optional)) |
| Feature `ollama` | `registry.ollama.ai`, `dd20bb891979d25aebc8bec07b2b3bbc.r2.cloudflarestorage.com` (Modell-Pull; Stabilität des zweiten Namens nicht gemessen: siehe [Ollama](#ollama-optional)) |
| Clone-Quelle | der Host von `devcontainer.sandbox.repository` beziehungsweise `origin` |
| Ihre Ergänzungen | `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema, Port oder Wildcard) |

**Die Liste anzeigen:** `u-boot egress list` zeigt die Liste, die
`u-boot generate devcontainer` in das Startscript schreibt, und je Host seine
Herkunft:

```bash
u-boot egress list
u-boot egress list --json
```

Auszug der Ausgabe (drei von mehreren Hosts, ohne den Block `Notes:` am Ende):

```text
Egress restriction: active.

HOST                ORIGIN
api.github.com      base
ghcr.io             runtime, feature:homebrew
registry.npmjs.org  feature:claude-code, feature:node
```

- **Herkunft:** `base` (immer), `runtime` (`nestedRuntime: podman` oder
  `docker`), `docker-apt` (`nestedRuntime: docker`), `feature:<name>` (aktives
  Feature), `clone` (Host von `devcontainer.sandbox.repository` beziehungsweise
  `origin`) und `user` (`egress.allow`). Ein Host, den mehrere Quellen liefern,
  nennt alle (zum Beispiel `ghcr.io` bei `podman` oder `docker` zusammen mit
  `homebrew`).
- **Was der Befehl liest:** nur `u-boot.yaml` und die Git-Konfiguration des
  Projekts (`.git/config`). Er braucht weder Netz noch Docker noch einen
  laufenden Container und ändert nichts. Ein Test vergleicht seine Liste mit der
  Zeile `hosts=` des erzeugten `egress-init.sh`; beide stammen aus derselben
  Berechnung.
- **Aktiv oder nicht:** Die erste Zeile nennt, ob die Beschränkung für die
  Konfiguration aktiv ist (Sandbox-Profil und `egress.enabled` gelten beide). Ist
  sie es nicht, nennt die Zeile den Grund, und die Liste ist die, die bei
  Aktivierung gälte.
- **Clone-Quelle:** Ohne Remote `origin` in `.git/config` und ohne
  `devcontainer.sandbox.repository` hat die Liste keinen Clone-Host; die Ausgabe
  sagt das. Der Befehl löst die Clone-Quelle auf wie
  `u-boot generate devcontainer --sandbox`, auch wenn das Sandbox-Profil nicht
  gesetzt oder die Beschränkung nicht aktiv ist (die Liste ist die, die bei
  Aktivierung gälte). Eine `origin`-URL, die diese Auflösung abweist
  (Zugangsdaten oder unzulässige Zeichen), weist deshalb auch `egress list` ab
  (Exit-Code 10), ohne Sandbox-Profil und ohne aktive Beschränkung ebenso. Ein
  `devcontainer.sandbox.repository` mit Zugangsdaten weist schon die Prüfung der
  Konfiguration ab (Exit-Code 10).
- **Nicht der Zustand des Containers:** Die Liste wird aus den Projektdateien
  berechnet, nicht aus einem Container gelesen. Ein Container behält die Liste des
  Images, aus dem er gebaut wurde (siehe den Aufzählungspunkt zum Neuaufbau
  weiter unten).
- **Optionen:** `--json` liefert `active`, `sandboxProfile`, `egressEnabled`,
  `cloneSource` und `hosts` (je Host `host` und `sources`). `--quiet` ändert die
  Ausgabe nicht. `--dry-run` und `--diff` weist der Befehl mit Exit-Code 2 ab.

- **Wichtig:** Die Hosts eines Agenten kennt `u-boot` nur für das Feature
  `claude-code` (Zeile in der Tabelle). Die Hosts jedes anderen Agenten (zum
  Beispiel dessen API) tragen Sie in `egress.allow` ein.
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
- Für `claude-code` nimmt die Liste nur den API-Host und die npm-Registry auf; die
  Hosts der Anmeldung sowie von Telemetrie, Fehlerberichten, MCP-Connectors,
  Updates und Plugins sind **nicht aufgenommen** und werden bei Bedarf in
  `egress.allow` ergänzt (Einzelheiten im Abschnitt
  [Claude Code](#claude-code-optional)).
- Die Features `sdkman` und `homebrew` bringen nur ihre eigenen Hosts mit. Die
  Download-Hosts einzelner SDKMAN-Kandidaten (zum Beispiel `services.gradle.org`
  für Gradle, `github.com` und `release-assets.githubusercontent.com` für
  GitHub-Releases) und die Hosts weiterer Quellen von Homebrew-Formeln tragen
  Sie bei Bedarf in `egress.allow` ein.
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
- Wer die Fähigkeit hat und das `sudo` des Base-Images nutzen kann, kann die
  Regeln aufheben. Der Verkehr verschachtelter Container und die
  iptables-Regeln eines verschachtelten Docker-Daemons umgehen die Beschränkung;
  `u-boot doctor` warnt bei `docker` mit Beschränkung. Gemessen auf Linux (`amd64`) in
  zwei Läufen: Ein verschachtelter Docker-Container erreichte `example.com` und eine
  gesperrte Adresse; ein verschachtelter Podman-Container nicht (Zeitüberschreitung).
  Die Namensauflösung eines verschachtelten Podman-Containers läuft über den lokalen DNS-Server
  oder wird verworfen (siehe *Namensauflösung und ihre Grenze* im Abschnitt zur Egress-Beschränkung).
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
- Ohne `NET_ADMIN` entfällt die Beschränkung mit einer Warnung (`warn`) oder das
  Startscript endet mit Exit-Code 11 (`fail`). `u-boot doctor` prüft die
  Konfiguration (ohne Sandbox-Profil ist sie wirkungslos und `doctor` warnt) und
  weist den Wegfall der Beschränkung aus der Statusdatei aus (siehe
  [Fallbacks sehen](#fallbacks-sehen)).

#### Login-Dienst `sshd` (optional)

Das offizielle Feature `sshd` installiert einen SSH-Server im Sandbox-Container.
Er öffnet einen Login-Dienst mit Passwort- und Schlüssel-Login, Root-Login ist
erlaubt; lesen Sie vor Schritt 1 den [Sicherheitshinweis](#sicherheitshinweis-sshd)
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

**Schritt 1: Feature aktivieren und den Container neu aufbauen**

```bash
u-boot config set devcontainer.features.sshd.enabled true
u-boot generate devcontainer
```

Das Installationsskript eines Features läuft beim Image-Build. Bauen Sie den
Container deshalb nach `generate devcontainer` neu auf. Nicht gemessen ist, ob
ein bloßer Neustart des bestehenden Containers nach einer geänderten Feature-Liste
genügt; für die Egress-Liste ist gemessen, dass er nicht genügt (Abschnitt
[Internet-Beschränkung (Egress)](#internet-beschränkung-egress)).

**Schritt 2: Zugang einrichten**

Die Konten `root` und `vscode` sind nach dem Build gesperrt. Ein Login geht erst,
wenn Sie einen Schlüssel eintragen oder ein Passwort setzen; `u-boot` setzt kein
Passwort und legt keine Schlüssel an. Bevorzugt ist der Schlüssel; das Passwort
ist die Alternative.

- **Schlüssel (bevorzugt):** Der Projektordner ist ein Volume, `~/.ssh` des
  Hosts ist im Container nicht eingebunden; der Schlüssel kommt also nicht über
  einen Mount in den Container.
  1. Auf dem Host ein Schlüsselpaar erzeugen:

     ```bash
     ssh-keygen -t ed25519 -f ~/.ssh/uboot_sshd
     ```

  2. Den Inhalt der öffentlichen Datei `~/.ssh/uboot_sshd.pub` in
     `~/.ssh/authorized_keys` des Container-Benutzers `vscode` eintragen, in einem
     Terminal im Container (Editor-Terminal oder `docker exec`; `docker exec`
     nicht gemessen):

     ```bash
     mkdir -p ~/.ssh && chmod 700 ~/.ssh
     echo '<Inhalt von uboot_sshd.pub>' >> ~/.ssh/authorized_keys
     chmod 600 ~/.ssh/authorized_keys
     ```

- **Passwort (Alternative):** Im Container `sudo passwd vscode` (bzw. den Namen
  Ihres Nutzers) ausführen, mit einem starken Passwort: der Dienst erlaubt
  Passwort-Login und Root-Login. `sudo passwd` ohne Argument setzt das Passwort
  von `root`, dessen Login erlaubt ist.

**Schritt 3: Verbinden**

`u-boot` veröffentlicht Port 2222 nicht und leitet ihn nicht weiter:
`forwardPorts` in `devcontainer.json` stammt nur aus den Diensten, die `runArgs`
des Sandbox-Profils enthalten keine Port-Veröffentlichung. Sie leiten den Port
selbst auf den Host weiter, über die Port-Weiterleitung Ihres Dev-Container-Werkzeugs
(nicht gemessen, wie ein Editor oder die CLI das tut) oder eine Port-Option der
Container-Engine (nicht gemessen). Veröffentlichen Sie den Port nur an der
Loopback-Adresse des Hosts, nach dem Muster `127.0.0.1:2222:2222`, nicht als
`2222:2222`: sonst öffnen Passwort-Login und erlaubter Root-Login den Dienst auf
allen Schnittstellen des Hosts. Danach verbinden Sie sich über den Host-Port, hier
mit der Weiterleitung auf denselben Port:

```bash
ssh -i ~/.ssh/uboot_sshd -p 2222 vscode@localhost
```

Beim Passwort-Login entfällt `-i`. Der Aufruf ist im Näherungsaufbau gemessen,
nicht mit dem Feature selbst.

<a id="sicherheitshinweis-sshd"></a>

**Sicherheitshinweis**

- Die Konfiguration des Features lässt Passwort-Login (`PasswordAuthentication yes`),
  Schlüssel-Login (`PubkeyAuthentication yes`, sobald Sie einen Schlüssel eintragen)
  und `PermitRootLogin yes` zu; `u-boot` ändert das nicht. Der Port lauscht im
  Container auf allen Schnittstellen. Veröffentlichen oder leiten Sie den Port nur
  an die Loopback-Adresse des Hosts (`127.0.0.1`), nicht an alle Schnittstellen.
- Das Sandbox-Profil begrenzt Schaden und ist keine Isolationsgrenze. Das Feature
  bringt weder `--privileged` noch zusätzliche Fähigkeiten, Mounts oder Devices
  mit, nur einen Entrypoint, der `sshd` startet. Mit der Egress-Beschränkung
  braucht es keinen zusätzlichen Host (der Build holt Pakete von `deb.debian.org`,
  das immer erlaubt ist).
- Die Warnung von `u-boot generate devcontainer` nennt Passwort- und Schlüssel-Login
  und den erlaubten Root-Login; sie bleibt unter `--quiet` sichtbar.

#### Claude Code (optional)

```bash
u-boot config set devcontainer.features.node.enabled true
u-boot config set devcontainer.features.claude-code.enabled true
u-boot generate devcontainer
```

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
  Variable auf dem Host nicht gesetzt, bleibt sie im Container leer; die CLI
  behandelt das wie einen fehlenden Schlüssel (gemessen: „Not logged in“).
  `u-boot` reicht nur diese eine Variable durch; andere Wege (etwa ein Token mit
  `CLAUDE_CODE_OAUTH_TOKEN` oder ein Login mit `/login` im Container) gehören nicht zur
  Zusage und sind nicht gemessen. Wohin ein Login im Container seine Zugangsdaten
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
  Der Weg über einen Abo-Login oder `CLAUDE_CODE_OAUTH_TOKEN` ist nicht zugesagt
  und nicht gemessen.
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

#### Ollama (optional)

```bash
u-boot config set devcontainer.features.ollama.enabled true
u-boot generate devcontainer
```

Das Feature `ollama` installiert Ollama (lokale Sprachmodelle) im Sandbox-Container,
nur für die CPU. Es ist kein OCI-Feature und erscheint nicht unter `features` in
`devcontainer.json`: Das erzeugte Dockerfile lädt beim Image-Build das Archiv
`ollama-linux-<arch>.tar.zst` aus dem Release `0.40.6` von `github.com/ollama/ollama`,
prüft die SHA-256-Prüfsumme **vor** dem Entpacken (bei einer Abweichung bricht der Build
ab) und entpackt es nach `/usr/local`. Version und Prüfsummen stehen an einer Stelle im
Dockerfile. Eine andere Version gibt es nicht: `version` ohne `source` weist `u-boot` mit
Exit-Code `10` ab, und eine `source:`-Überschreibung macht den Schlüssel zu einem
externen Feature (Freigabe wie unter
[Externe Feature-Quellen freigeben](#externe-feature-quellen-freigeben)) ohne
Dockerfile-Block.

- **Nur im Sandbox-Profil:** Ohne Sandbox-Profil erzeugt `u-boot` keinen
  Dockerfile-Block und kein Volume und meldet, dass das Feature keine Wirkung hat.
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

Modell holen und nutzen (im Terminal des Containers):

```bash
ollama serve > /tmp/ollama.log 2>&1 &
ollama pull smollm2:135m
ollama run smollm2:135m "Reply with the single word OK."
```

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
  `ollama.com`, `example.com` und `huggingface.co` ließen sich nicht auflösen. `ollama.com` ist nicht aufgenommen: gesperrt
  meldet der Server beim Start zwei Warnzeilen, Start, Pull und Lauf funktionierten
  (gemessen).
- **Cloud-Modelle sind nicht enthalten:** Ein Cloud-Modell von Ollama (Modellname mit der
  Endung `-cloud`, etwa `gpt-oss:120b-cloud`) lässt sich ziehen: `ollama pull` meldet `success`,
  und das Modell steht in `ollama list`. Es läuft aber nicht: `ollama run` scheitert, weil sich
  `ollama.com` nicht auflösen lässt; der Host steht nicht in der Allowlist, und `u-boot` nimmt ihn
  für `ollama` nicht auf. `ollama signin` meldete ohne Konto `account unavailable`; eine Anmeldung mit Konto
  und ob ein Cloud-Modell mit Konto läuft, sind nicht gemessen. Gemessen in einem Lauf auf Linux (`amd64`,
  native Docker Engine 29.9.1, Ollama `0.40.6`) hinter der Egress-Beschränkung.
- **Grenzen der Messung:** Zwei Messumgebungen (auf dem `amd64`-Host zwei Läufe), Ollama `0.40.6`:
  Docker 29.5.2 auf Colima 0.10.3 (`arm64`) und, am 2026-10-19, Linux x86_64 (`amd64`) mit
  einer nativen Docker Engine 29.9.1 ohne VM (Kernel 6.8, `kernel.apparmor_restrict_unprivileged_userns=0`,
  `/dev/fuse` und `/dev/net/tun` vorhanden). Auf dem `amd64`-Host baute das erzeugte
  Dockerfile (Prüfsumme bestanden), der Server lauschte nur auf `127.0.0.1:11434` und meldete
  `library=cpu`, Pull und Lauf der beiden Modelle gelangen hinter der Beschränkung, die Modelle
  überstanden das Löschen des Containers und den Neuaufbau, und die Kombination mit
  `nestedRuntime: podman` und mit `nestedRuntime: docker` baute, ließ die Paketliste unverändert
  und führte einen kleinen verschachtelten Container aus, während der Server und ein Pull liefen. Das
  sind **zwei Messläufe auf einem Host, kein Lauf in der CI**: Ein CI-Lauf von `amd64` ist nicht
  gemessen. Nicht gemessen: Docker Desktop, Linux hinter einer VM, größere
  Modelle und Mindest-Speicher, `OLLAMA_NO_CLOUD`, ein Start des Servers über die
  `devcontainer.json` (gemessen ist der Start von Hand und der Egress-Start durch die
  Dev-Container-CLI). GPU-Beschleunigung, ein Anschluss an einen Ollama-Dienst des Hosts
  und die Verwaltung von Modellen durch `u-boot` gibt es nicht.

#### Zugangsdaten (Git und Agent)

Zugangsdaten stehen **nie** im Image, im Volume, in `u-boot.yaml` oder in einer
erzeugten Datei:

- `devcontainer.json` enthält nur die Referenz
  `"GIT_TOKEN": "${localEnv:GIT_TOKEN}"`. Setzen Sie `GIT_TOKEN` in der
  Host-Umgebung, bevor Sie den Container starten. Mit dem Feature `claude-code`
  kommt `"ANTHROPIC_API_KEY": "${localEnv:ANTHROPIC_API_KEY}"` hinzu; ohne das
  Feature bleibt die Datei unverändert.
- Ein Credential-Helper im Image liefert `$GIT_TOKEN` an `git` (nur für
  `https`-Remotes; bei `ssh` gibt es keinen Schlüssel im Container).
- Starten Sie den Container mit `u-boot devcontainer up`, reicht der Befehl nur den Namen
  der Variable durch; die Grenzen dieses Wegs (Metadaten bei `containerEnv`, Kommandozeile
  im Hilfs-Container während eines Lebenszyklus-Befehls) stehen in Abschnitt 4.9. In einer
  Shell, die Sie ohne Editor mit `docker exec` öffnen, gilt `remoteEnv` nicht, und `GIT_TOKEN`
  ist dort nicht gesetzt; geben Sie die Variable mit `docker exec -e GIT_TOKEN` mit (Abschnitt
  4.9, „Im Container arbeiten“).
- Empfohlen: **kurzlebige, auf das Repository begrenzte Tokens**, nie ein
  privater Schlüssel im Container, und **Branch-Protection** im Remote, damit ein
  kompromittierter Agent nicht auf geschützte Branches schreiben kann.
- `u-boot doctor` warnt vor möglichen Klartext-Tokens und -Schlüsseln in `u-boot.yaml`,
  `compose.yaml`, `.env.example` und `.devcontainer/*` (auch Anthropic-Schlüssel
  mit dem Präfix `sk-ant-`, nach dem Ausdruck, den die Claude Code CLI selbst zum
  Erkennen eines Schlüssels in Text verwendet, sowie ein Wert für
  `ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN` oder `CLAUDE_CODE_OAUTH_TOKEN`), vor einer
  fehlenden Token-Quelle bei einem `https`-Clone und, bei aktivem Feature
  `claude-code`, vor einer fehlenden Referenz auf `ANTHROPIC_API_KEY` in
  `devcontainer.json`. Die Prüfung liest nur Projektdateien; sie sieht weder die
  Umgebung des Hosts noch einen Schlüssel in einer anderen Datei. Bekannte Grenze:
  Der Ausdruck für eine Zuweisung kann einen Platzhalter nicht von einem Wert
  unterscheiden; ein Platzhalter wie `ANTHROPIC_API_KEY=your-api-key-here` (ab zwölf
  Zeichen) in `.env.example` löst die Warnung aus.

#### Mehrere Instanzen

Die Volumes heißen `<projekt>-workspace-${devcontainerId}` und (mit
verschachteltem Podman) `<projekt>-containers-${devcontainerId}` beziehungsweise (mit
verschachteltem Docker) `<projekt>-docker-${devcontainerId}`, mit dem Feature `ollama` zusätzlich
`<projekt>-ollama-${devcontainerId}`. Das
Dev-Containers-Werkzeug löst `devcontainerId` pro Projektordner auf.

- Parallele Instanzen desselben Projekts starten Sie aus **getrennten Ordnern**
  (Klone oder `git worktree`); jede Instanz hat eigene Volumes.
- Derselbe Ordner bleibt eine Instanz.
- Die Compose-Ports (`u-boot up`) mehrerer Instanzen kollidieren weiterhin.
- Alte Volumes `<projekt>-workspace` und `<projekt>-containers` aus früheren
  Entwicklungsständen bleiben verwaist; entfernen Sie sie mit `docker volume rm`.

#### Grenzen

- Geprüft unter Docker auf Linux mit der `devcontainer`-CLI; `nestedRuntime: docker`
  zusätzlich unter Colima (Linux-VM); nur diese beiden Hosts sagt das Lastenheft zu.
  **Nicht geprüft und nicht zugesagt:** Docker Desktop,
  Podman als Host-Engine, VS Code und Codespaces.
- `nestedRuntime: podman`: Podman 6.1.4 und `podman build` sind unter Docker auf
  Colima (arm64) mit erlaubten verschachtelten User-Namespaces gemessen; für
  amd64 sind Download, Prüfsumme, Signatur und Versionsabfrage unter Emulation
  geprüft. Ein verschachtelter Lauf (`podman run` eines kleinen Images, ohne und
  hinter der Egress-Beschränkung, mit einem Pull hinter ihr) ist auf Linux x86_64
  (`amd64`, native Docker Engine 29.9.1, 2026-10-19) in zwei Läufen gemessen, zusammen mit dem
  Feature `ollama`; `podman build` auf `amd64` hinter der Egress-Beschränkung läuft durch, braucht aber
  bei einem Namen außerhalb der Liste etwa 36 Sekunden je Name länger
  (siehe *Namensauflösung und ihre Grenze* im Abschnitt zur Egress-Beschränkung); die CI ist nicht gemessen.
  Podman 6 unterstützt nur cgroups v2 (laut Release-Notes, nicht gemessen);
  Hosts mit cgroups v1 sind davon betroffen.
  Unter einer Ubuntu-24.04-VM mit dem Schalter
  `kernel.apparmor_restrict_unprivileged_userns=1` startet verschachteltes
  Podman nicht (Exit-Code 11). Nicht verfügbar sind `buildah`, `docker buildx`,
  `docker compose` und BuildKit-Frontends.
- Der Wechsel von `podman` oder `docker` zurück auf `none` lässt eine vorhandene
  `sandbox-init.sh` beziehungsweise `docker-init.sh` im Ordner liegen; das Dockerfile verwendet sie dann nicht
  mehr.

### 4.12 Vorlagen nutzen

#### Vorgehen

```bash
u-boot template list                                  # Katalog ansehen
u-boot init mein-projekt --template basic             # aus Katalogvorlage
u-boot init mein-projekt --template ./meine-vorlage   # aus eigenem Verzeichnis
```

#### Ergebnis

`template list` zeigt eine Tabelle mit `NAME`, `DESCRIPTION` und `VERSION`:

```text
NAME   DESCRIPTION                                          VERSION
basic  Minimal u-boot project skeleton — same files …       0.1.0
```

`init --template` rendert das Projekt aus der gewählten Vorlage.

#### Hinweise

- Ein eigenes Vorlagenverzeichnis muss eine `template.yaml` enthalten.
  Symlinks im Vorlagenbaum werden abgelehnt.
- `--template` gilt nur für frische Projekte und lässt sich nicht mit
  `--devcontainer`, `--force`, `--backup`, `--dry-run` oder `--diff`
  kombinieren.
- `u-boot template list --json` liefert den Katalog als JSON-Liste.

### 4.13 In Skripten und CI verwenden

#### Voraussetzung

Der Lauf darf nicht auf eine Eingabe warten.

#### Vorgehen

Wählen Sie **einen** der beiden Modi:

```bash
u-boot add postgres --yes              # Rückfragen automatisch bejahen
u-boot add postgres --no-interactive   # bei jeder Rückfrage abbrechen
```

Maschinenlesbare Ausgabe erhalten Sie mit `--json`:

```bash
u-boot doctor --json
u-boot config get project.name --json
```

#### Ergebnis

`--json` liefert einen einheitlichen Umschlag mit Status, Befehl, Daten,
Diagnosemeldungen und Exit-Code:

```json
{"status":"ok","command":"config","subcommand":"get","diagnostics":[],"exitCode":0,"data":{"path":"project.name","value":"demo-app"}}
```

Alle Felder, die Ausgabe je Befehl und die Prüfkennungen beschreibt der Anhang
[JSON-Ausgabe](#json-ausgabe).

#### Hinweise

- Werten Sie in Skripten den **Exit-Code** aus, nicht den Ausgabetext (siehe
  [Exit-Codes](#exit-codes)).
- `--yes` und `--no-interactive` gleichzeitig ist ein Nutzungsfehler und endet
  mit Exit-Code 2.
- In CI empfiehlt sich `u-boot doctor --strict`, damit auch Warnungen den Lauf
  anhalten.

---

## 5. Konfiguration

### Die Projektdatei `u-boot.yaml`

```yaml
schemaVersion: 1
project:
  name: demo-app
services:
  postgres:
    enabled: true
devcontainer:
  enabled: true
```

`schemaVersion` gehört `u-boot` und sollte nicht von Hand geändert werden.

### Werte lesen

```bash
u-boot config                                          # ganze Datei anzeigen
u-boot config get project.name                         # einen Wert lesen
u-boot config get project.name services.postgres.enabled   # mehrere, je Zeile ein Wert
u-boot config list                                     # alle gesetzten Pfade als pfad=wert
```

Beispielausgabe von `u-boot config list`:

```text
project.name=demo-app
services.keycloak.enabled=true
services.postgres.enabled=true
```

Mit `--json` liefern `config get` bei mehreren Pfaden (oder mit `--json-array`)
und `config list` die Form `data.entries[{path, value}]`. Bei mehreren Pfaden gilt
„alles oder nichts“: Ist ein Pfad unbekannt oder nicht gesetzt (zum Beispiel
`devcontainer.enabled` in einem Projekt ohne Devcontainer), endet der Aufruf mit
Exit-Code 10 und gibt nichts aus.

### Werte ändern

```bash
u-boot config set project.name neuer-name
u-boot config set project.name neuer-name devcontainer.enabled true   # atomar
u-boot config set devcontainer.enabled true --dry-run                 # nur Vorschau
```

Mehrere Pfad-Wert-Paare in einem `config set` werden **atomar** geschrieben:
Schlägt eine Prüfung fehl, bleibt `u-boot.yaml` unverändert. Die Paare gelten
der Reihe nach; ein späteres Paar darf auf ein früheres aufbauen (zum Beispiel
erst `devcontainer.featureSources.allow`, dann die Feature-Quelle).

Jeder Aufruf wird gegen das Schema geprüft, **bevor** geschrieben wird.
Ungültige Werte führen zu Exit-Code 10, ohne die Datei anzufassen. Fehler tragen
mit `--json` einen strukturierten Hinweis in `data.hint` (`command`, `action`,
`argument`), der den Reparaturbefehl nennt.

### Verfügbare Pfade

| Pfad | Bedeutung | Werte | Schreibbar |
|---|---|---|---|
| `project.name` | Projektname | Kleinbuchstaben, Ziffern, Bindestriche | ja |
| `devcontainer.enabled` | Devcontainer-Unterstützung | `true` / `false` | ja |
| `devcontainer.profile` | Devcontainer-Profil | `sandbox` | ja |
| `devcontainer.user.uid` | Benutzer-ID im Container | 1 bis 65535 (Standard 1000) | ja |
| `devcontainer.features.<name>.enabled` | Feature aktiv | `true` / `false` | ja |
| `devcontainer.features.<name>.version` | Feature-Version | zum Beispiel `21` | ja |
| `devcontainer.features.<name>.source` | externe Feature-Quelle | URL aus der Freigabeliste | ja |
| `devcontainer.featureSources.allow` | Freigabeliste externer Quellen | URLs, kommagetrennt | ja |
| `devcontainer.sandbox.nestedRuntime` | Container-Laufzeit im Sandbox | `none` / `podman` / `docker` | ja |
| `devcontainer.sandbox.onUnavailable` | Verhalten ohne Fähigkeit | `warn` / `fail` | ja |
| `devcontainer.sandbox.repository` | Clone-Quelle statt `origin` | Git-Adresse ohne Zugangsdaten | ja |
| `devcontainer.sandbox.egress.enabled` | Internet-Beschränkung | `true` / `false` | ja |
| `devcontainer.sandbox.egress.allow` | zusätzlich erlaubte Hosts | Hostnamen, kommagetrennt, klein geschrieben | ja |
| `services.<dienst>.enabled` | Dienst aktiv | `true` / `false` | **nein**, nur lesbar |

Dienste schalten Sie über `u-boot add` und `u-boot remove`, nicht über
`config set`. Der Grund: Ein Dienst besteht nicht nur aus einem Schalter,
sondern auch aus Einträgen in `compose.yaml` und `.env.example` – ein direkt
gesetzter Schalter würde diese auseinanderlaufen lassen. Der Versuch endet mit
Exit-Code 10 und nennt den passenden Befehl.

Nach Änderungen an `devcontainer.*` führen Sie `u-boot generate devcontainer`
aus, damit die Dateien in `.devcontainer/` zur Projektdatei passen.

### Umgebungsvariablen

| Variable | Wirkung |
|---|---|
| `DOCKER_HOST` | Adresse des Docker- oder Podman-Daemons |
| `GIT_TOKEN` | Zugangsdaten für `https`-Remotes im Sandbox-Devcontainer (nur zur Laufzeit, auf dem Host gesetzt) |
| `ANTHROPIC_API_KEY` | API-Schlüssel für das Feature `claude-code` im Sandbox-Devcontainer (nur zur Laufzeit, auf dem Host gesetzt; `devcontainer.json` trägt nur die Referenz) |

### Ausführlichkeit der Ausgabe

| Option | Wirkung |
|---|---|
| `--quiet` | nur Fehler; Wirkung, die vier sichtbaren Ausnahmen und die Vorrang-Regel stehen unter [Allgemeine Optionen](#allgemeine-optionen) |
| `--verbose` | zusätzliche Details; gilt nur, wenn es die zuletzt angegebene der Optionen `--quiet`, `--verbose`, `--debug` ist |
| `--debug` | interne Diagnoseausgabe; gilt nur, wenn es die zuletzt angegebene der Optionen `--quiet`, `--verbose`, `--debug` ist (`--verbose` und `--debug` wirken gleich) |

---

## 6. Rollen und Rechte

`u-boot` kennt keine Benutzerrollen. Es braucht nur die Rechte, die auch die
Handarbeit bräuchte:

| Aufgabe | Benötigtes Recht |
|---|---|
| Dateien anlegen und ändern (`init`, `add`, `remove`, `generate`, `config set`) | Schreibrecht im Projektverzeichnis (`doctor` prüft das: `fs.write-permissions`) |
| `up`, `down`, `logs`, `devcontainer up`, `devcontainer down` | Zugriff auf den Docker-Daemon; unter Linux meist Mitgliedschaft in der Gruppe `docker` |
| Installation in `/usr/local/bin` oder per Paket | Administratorrechte (`sudo`) |
| `git`-Funktionen von `init` | `git` installiert und konfiguriert |

Hinweis: Im Devcontainer läuft Ihre Arbeit als Benutzer `vscode`, nicht als
root. Ihre Geheimnis-Verzeichnisse bindet weder das Standard- noch das
Sandbox-Profil ein. Das Docker-Socket des Hosts bindet das Standardprofil ein, wenn ein
Feature es tut, zum Beispiel `docker-cli`; das Sandbox-Profil weist dieses Feature ab
(siehe 4.9 und 4.10). Das gilt für den erzeugten Devcontainer;
der Hilfs-Container von `u-boot devcontainer up` ist davon ausgenommen (siehe 4.9).

---

## 7. Fehlerbehebung

### Exit-Codes

Werten Sie in Skripten diese Codes aus:

| Code | Bedeutung |
|---|---|
| `0` | Erfolg (auch „nichts zu tun") |
| `1` | allgemeiner Fehler |
| `2` | falsche Benutzung (unbekanntes Flag, fehlendes Argument, widersprüchliche Modus-Flags, ungültiger `--timeout`, `--tail`, `--since`) |
| `10` | fachlicher Fehler (ungültiger Name oder Wert, fehlendes Projekt, nicht unterstützter Dienst, verweigerte Bestätigung, vorhandene Datei) |
| `11` | Umgebungsproblem (Docker nicht erreichbar, Compose-Plugin fehlt; auch `doctor` mit Fehlern; bei `devcontainer up` auch ein abgebrochenes Startscript des Containers, auch aus einem früheren Start desselben Containers) |
| `12` | Laufzeitfehler von Compose (auch Zeitüberschreitung bei `up`); bei `devcontainer up` und `devcontainer down` ein Fehler beim Bau oder Start (Abschnitt 4.9) |
| `14` | Dateisystem- oder Persistenzfehler |

### Fehler: „project not initialized: u-boot.yaml missing"

#### Ursache

Sie stehen nicht in einem `u-boot`-Projekt, oder die Projektdatei fehlt. Der
Befehl endet mit Exit-Code 10.

#### Lösung

1. Prüfen Sie Ihr Verzeichnis: `pwd` und `ls u-boot.yaml`.
2. Wechseln Sie in das Projektverzeichnis – oder legen Sie mit `u-boot init`
   ein neues Projekt an.

### Fehler: „docker daemon unreachable" (Exit-Code 11)

#### Ursache

Der Docker-Daemon läuft nicht, Ihr Benutzer darf nicht auf ihn zugreifen, oder
`DOCKER_HOST` zeigt auf eine falsche Adresse.

#### Lösung

1. Führen Sie `u-boot doctor` aus und lesen Sie die `→`-Zeile der
   fehlgeschlagenen Prüfung.
2. Prüfen Sie den Daemon: `docker version`.
3. Prüfen Sie das Compose-Plugin: `docker compose version`.
4. Prüfen Sie `echo $DOCKER_HOST`; setzen Sie die Variable zurück, wenn sie
   nicht beabsichtigt ist.
5. Unter Linux: Ist Ihr Benutzer in der Gruppe `docker`?

### Fehler: „file exists: … pass --backup or --force" (Exit-Code 10)

#### Ursache

Es gibt bereits eine Datei, die `u-boot` nicht über einen verwalteten Bereich
ändern kann. `u-boot` überschreibt sie nicht von sich aus.

#### Lösung

1. Sehen Sie sich die geplanten Änderungen an: `u-boot init --diff`.
2. Entscheiden Sie sich:
   - `--force` ersetzt den verwalteten Bereich,
   - `--backup` sichert die Datei vorher.
3. Führen Sie den Befehl mit dem gewählten Flag erneut aus.

### Fehler: „confirmation required" (Exit-Code 10)

#### Ursache

Der Befehl braucht eine Bestätigung (`down --volumes`, `remove --purge`), kann
sie aber nicht einholen, weil `--no-interactive` oder `--json` gesetzt ist.

#### Lösung

- Für destruktive Operationen: `--yes` ergänzen – bewusst, denn dabei gehen
  Daten verloren.
- Für die Erkennung eines bestehenden Projekts bei `init`: `--assume-existing`
  ergänzen. `--yes` genügt hier **nicht**.

### Fehler: „--yes and --no-interactive are mutually exclusive" (Exit-Code 2)

#### Ursache

Beide Modus-Optionen wurden gleichzeitig angegeben.

#### Lösung

Wählen Sie eine: `--yes` (automatisch bejahen) oder `--no-interactive` (bei
Rückfrage abbrechen).

### Fehler: „service not supported: … is not in the built-in catalogue" (Exit-Code 10)

#### Ursache

Der Dienstname existiert nicht im Katalog.

#### Lösung

Verwenden Sie einen der genannten Namen (`postgres`, `keycloak`, `otel`) und
achten Sie auf die Schreibweise in Kleinbuchstaben.

### Fehler: „invalid project name"

#### Ursache

Bei `init` oder `config set project.name` passt der Name nicht zum erlaubten
Muster: Kleinbuchstaben, Ziffern und Bindestriche, beginnt mit einem Buchstaben,
endet mit Buchstabe oder Ziffer, bis zu 63 Zeichen (Exit-Code 10).

#### Lösung

Wählen Sie einen Namen wie `mein-projekt`. Unterstriche und Großbuchstaben
ersetzen Sie durch Bindestriche und Kleinbuchstaben.

### Fehler: „… is not writable via `u-boot config set`" (Exit-Code 10)

#### Ursache

Sie wollten `services.<dienst>.enabled` mit `config set` ändern. Dienste
verwalten `add` und `remove`.

#### Lösung

Nutzen Sie `u-boot add <dienst>` oder `u-boot remove <dienst>`.

### Fehler: „--tail must be a non-negative integer" oder „--timeout must be >= 0" (Exit-Code 2)

#### Ursache

Ein Zahlenwert ist negativ oder keine Zahl.

#### Lösung

Geben Sie eine ganze Zahl ab 0 an, zum Beispiel `u-boot logs --tail 200` oder
`u-boot up --timeout 120`.

### Fehler: Ein Port ist bereits belegt

#### Ursache

Ein anderer Prozess oder Container nutzt schon einen Port der Umgebung (zum
Beispiel 5432 für `postgres`, 8080 für `keycloak`, 4317 und 4318 für `otel`).
`u-boot up` meldet einen Compose-Laufzeitfehler (Exit-Code 12).

#### Lösung

1. Finden Sie den Verursacher: `docker ps` oder `ss -ltnp | grep 5432`.
2. Stoppen Sie ihn – bei einer anderen `u-boot`-Umgebung mit `u-boot down` in
   deren Projektordner – und starten Sie neu.

### Fehler: `up` meldet eine Zeitüberschreitung (Exit-Code 12)

#### Ursache

Ein Dienst wurde innerhalb von `--timeout` (Standard 60 Sekunden) nicht
„healthy" oder „running". Das passiert oft beim ersten Start, wenn Images noch
geladen werden, oder bei langsamen Diensten wie `keycloak`.

#### Lösung

1. Sehen Sie in die Logs: `u-boot logs <dienst> --tail 100`.
2. Erhöhen Sie die Wartezeit: `u-boot up --timeout 180`.
3. Prüfen Sie `.env`: Fehlen Werte oder stehen noch `CHANGEME_…`-Platzhalter?

### Fehler: Devcontainer und `u-boot.yaml` laufen auseinander

#### Ursache

Sie haben `devcontainer.*` in `u-boot.yaml` geändert oder `devcontainer.json`
von Hand bearbeitet. `u-boot doctor` meldet eine Warnung
(`devcontainer.features.drift` oder `devcontainer.forwardPorts.consistency`).

#### Lösung

Führen Sie `u-boot generate devcontainer` aus. Eigene Änderungen außerhalb des
verwalteten Bereichs bleiben erhalten.

### Fehler: Sandbox-Profil mit dem Feature `docker-cli` (Exit-Code 10)

#### Ursache

Das Profil ist `sandbox`, und das Feature `docker-cli` ist aktiv. Das Feature
bindet das Docker-Socket des Hosts in den Container ein, was das Sandbox-Profil
nicht tut; `u-boot generate devcontainer` endet deshalb mit Exit-Code 10 und
schreibt nichts, auch nicht mit `--dry-run`. `u-boot doctor` meldet dasselbe als
Fehler (Exit-Code 11).

#### Lösung

1. Schalten Sie das Feature ab:

   ```bash
   u-boot config set devcontainer.features.docker-cli.enabled false
   ```

   Brauchen Sie Docker im Container, setzen Sie zusätzlich
   `u-boot config set devcontainer.sandbox.nestedRuntime docker` (ein eigener
   Docker-Daemon im Container, Abschnitt 4.11). `nestedRuntime: docker` allein
   hebt den Fehler nicht auf.

2. Erzeugen Sie den Devcontainer neu: `u-boot generate devcontainer`.
3. Legen Sie einen vorhandenen Container neu an: `u-boot devcontainer up --recreate`.

### Fehler: Sandbox-Container startet nicht (Exit-Code 11 im Startscript)

#### Ursache

Verschachteltes Podman (`nestedRuntime: podman`) braucht Fähigkeiten, die der
Host nicht bietet, zum Beispiel blockierte User-Namespaces unter Ubuntu 24.04
oder ein im Container fehlendes `/dev/fuse` unter `onUnavailable: fail`. Ein
auf dem Container-Host fehlendes `/dev/fuse` oder `/dev/net/tun` verhindert den
Start des Containers durch Docker, das Startscript läuft dann nicht. Der Docker-Daemon (`nestedRuntime: docker`) wird
nicht bereit, wenn der Container ohne `--privileged` startet; das Log
`/var/log/u-boot-dockerd.log` im Container nennt den Grund. Das Startscript der
Egress-Beschränkung endet mit Exit-Code 11, auch unter `warn`, wenn der Benutzer
`nobody` im Image fehlt, die nftables-Regeln sich nicht anwenden lassen oder der
lokale DNS-Server nicht startet oder nicht unter diesem Benutzer läuft; die Meldung nennt den Grund.

#### Lösung

1. Lesen Sie die Meldung des Startscripts im Container-Log. `u-boot doctor`
   zeigt in den Einträgen `devcontainer.sandbox.fallback.*`, was die
   Startscripte gemeldet haben (siehe [Fallbacks sehen](#fallbacks-sehen)); ein
   nicht ermittelbarer Zustand heißt nicht, dass nichts vorlag.
2. Setzen Sie `devcontainer.sandbox.nestedRuntime` auf `none` und erzeugen Sie
   den Devcontainer neu, wenn Sie kein verschachteltes Podman brauchen.
3. Haben Sie den Container mit `u-boot devcontainer up` gestartet, legen Sie ihn nach einer
   Änderung mit `u-boot devcontainer up --recreate` neu an. Gegen einen laufenden Container
   führt ein erneutes `up` ohne `--recreate` das `postCreateCommand` und das `postStartCommand`
   nicht erneut aus; ist der Container gestoppt, läuft das `postStartCommand` nach `down` und
   `up` wieder. Meldet nach dem Aufruf die Statusdatei eines Startscripts `failed`, endet `up`
   wieder mit Exit-Code 11. Beheben Sie die Ursache und legen Sie den Container mit `--recreate`
   neu an (Abschnitt [4.9](#49-einen-devcontainer-einrichten), „Erneutes `up` nach einem Fehler
   oder Abbruch“).
4. Weitere Abhilfen stehen in Abschnitt [4.11](#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten) unter „Verschachteltes Podman".

### Wenn nichts davon hilft

1. Wiederholen Sie den Befehl mit `--debug`.
2. Führen Sie `u-boot doctor` aus und halten Sie die Ausgabe bereit.
3. Notieren Sie Version (`u-boot --version`), Betriebssystem und den genauen
   Befehl.
4. Eröffnen Sie ein Issue (siehe [Support und Lizenz](#11-support-und-lizenz)).

---

## 8. FAQ

**Brauche ich Docker, um `u-boot` zu benutzen?**
Nur für `up`, `down`, `logs`, `devcontainer up`, `devcontainer down` und die
Docker-Prüfungen von `doctor`. `init`, `add`, `remove`, `generate`, `config`,
`egress list` und `template list` laufen ohne.

**Überschreibt `u-boot` meine Änderungen?**
Nein. Außerhalb der verwalteten Bereiche wird nichts angefasst, und innerhalb
nur mit `--force` oder nach einer Sicherung mit `--backup`.

**Kann ich `u-boot` auf ein bestehendes Projekt anwenden?**
Ja. Nutzen Sie `u-boot init --dry-run`, um vorher zu sehen, was passieren würde.

**Wo trage ich Passwörter und Zugangsdaten ein?**
In `.env` (Kopie von `.env.example`), nicht in `.env.example` und nicht in
`u-boot.yaml`. `.env` gehört nicht ins Git-Repository.

**Warum kann ich `services.<dienst>.enabled` nicht mit `config set` ändern?**
Weil ein Dienst mehr ist als ein Schalter. `u-boot add` und `u-boot remove`
pflegen zusätzlich `compose.yaml` und `.env.example`.

**Was passiert, wenn ich einen Dienst zweimal hinzufüge?**
Nichts. Der zweite Aufruf ist wirkungslos und meldet keinen Fehler.

**Welche Dienste gibt es?**
`postgres`, `keycloak` und `otel` (OpenTelemetry-Collector).

**Kann ich eigene Einträge in `compose.yaml` ergänzen?**
Ja, außerhalb der Markierungen `BEGIN`/`END U-BOOT MANAGED BLOCK`. Solche
Einträge bleiben bei allen Befehlen erhalten.

**Funktioniert Podman?**
Ja, ab Version 4.0, wenn `docker` auf `podman` zeigt und `DOCKER_HOST` gesetzt
ist. `doctor` meldet die Podman-Version als Warnung, blockiert aber nicht.

**Warum überspringt `doctor` im Container vier Prüfungen?**
Aus dem Container heraus sind die Werkzeuge Ihres Rechners nicht sichtbar. Ein
Fehlschlag wäre irreführend, deshalb der Zustand `?`.

**Ist das Sandbox-Profil eine sichere Isolation?**
Nein. Es begrenzt Schäden (kein Host-Docker-Socket im Container, kein `--privileged`, optional
eine Liste erlaubter Internet-Ziele), ist aber keine harte Sicherheitsgrenze.
Das gilt besonders mit verschachteltem Podman.

**Wo finde ich alle Optionen eines Befehls?**
`u-boot <befehl> --help` zeigt sie; die Optionen aller Befehle stehen auch im
[Anhang](#optionen-je-befehl). Die JSON-Ausgabe beschreibt der Anhang
[JSON-Ausgabe](#json-ausgabe). `u-boot --help` verweist auf dieses Handbuch (die
Fassung Ihrer Version, bei einer selbst gebauten oder Entwicklungsversion nur die
aktuelle) und auf das README.

**Wie aktualisiere ich `u-boot`?**
Je nach Installation: Binary neu herunterladen, `brew upgrade u-boot`, das neue
`.deb`/`.rpm` installieren oder `docker pull ghcr.io/pt9912/u-boot:latest`.
Bestehende Projekte bleiben nutzbar; `u-boot generate …` bringt Artefakte auf den
neuen Stand.

---

## 9. Glossar

| Begriff | Bedeutung |
|---|---|
| Add-on | ein zuschaltbarer Dienst wie `postgres`, `keycloak` oder `otel` |
| Artefakt | eine von `u-boot` erzeugte Datei (README, CHANGELOG, `.env.example`, Devcontainer) |
| Compose | Docker Compose; auch die Datei `compose.yaml` |
| Devcontainer | containerisierte Entwicklungsumgebung, meist mit VS Code genutzt |
| Egress | ausgehender Netzwerkverkehr eines Containers |
| Exit-Code | Rückgabewert eines Befehls; `0` bedeutet Erfolg |
| Feature | Baustein eines Devcontainers (zum Beispiel `node` oder `go`) |
| Healthcheck | Prüfung, ob ein Dienst technisch bereit ist |
| Idempotenz | mehrfaches Ausführen führt zum selben Ergebnis |
| Projektdatei | `u-boot.yaml` – Quelle der Wahrheit für Name, Dienste und Devcontainer |
| Sandbox-Profil | Devcontainer-Variante für autonome Agenten (Volume-Workspace, ohne Host-Socket) |
| Stabilisierung | Wartezeit nach `up`, bis jeder Dienst als bereit gilt |
| Verwalteter Bereich | Abschnitt zwischen `BEGIN`/`END U-BOOT MANAGED BLOCK`, den `u-boot` pflegt |
| Volume | von Docker verwalteter Datenspeicher, der das Neustarten des Containers überlebt |

---

## 10. Anhang

### Befehlsübersicht

| Befehl | Zweck | Abschnitt |
|---|---|---|
| `u-boot init [name]` | Projekt anlegen oder erneut initialisieren | 4.1, 4.2 |
| `u-boot add <dienst>` | Dienst hinzufügen | 4.3 |
| `u-boot remove <dienst>` | Dienst entfernen | 4.4 |
| `u-boot up` | Umgebung starten | 4.5 |
| `u-boot down` | Umgebung stoppen | 4.5 |
| `u-boot logs [dienst…]` | Logs ansehen | 4.6 |
| `u-boot doctor` | Umgebung prüfen | 4.7 |
| `u-boot generate <artefakt>` | Artefakt erzeugen oder aktualisieren | 4.8–4.11 |
| `u-boot config [get\|set\|list]` | Konfiguration lesen und ändern | 5 |
| `u-boot egress list` | Egress-Liste mit Herkunft je Host anzeigen | 4.11 |
| `u-boot devcontainer up` | Devcontainer ohne Editor starten | 4.9 |
| `u-boot devcontainer down` | Devcontainer stoppen (entfernt nichts) | 4.9 |
| `u-boot template list` | Vorlagenkatalog anzeigen | 4.12 |
| `u-boot completion <shell>` | Shell-Vervollständigung erzeugen | 2 |

### Optionen je Befehl

| Befehl | Optionen |
|---|---|
| `init` | `--no-git`, `--force`, `--backup`, `--assume-existing`, `--devcontainer`, `--sandbox`, `--template`, `--allow-external-feature-sources`, `--dry-run`, `--diff` |
| `add` | `--with-deps`, `--dry-run`, `--diff` |
| `remove` | `--purge`, `--dry-run`, `--diff` |
| `up` | `--timeout <sekunden>` |
| `down` | `--volumes` |
| `logs` | `--follow`, `--tail <n>`, `--since <t>`, `--until <t>`, `--timestamps`, `--no-log-prefix` |
| `doctor` | `--strict` |
| `generate` | `--dry-run`, `--diff`; nur für `devcontainer`: `--sandbox`, `--allow-external-feature-sources` |
| `config get` | `--json-array` |
| `config set` | `--dry-run`, `--diff`, `--allow-external-feature-sources` |
| `egress list` | keine eigenen; `--dry-run` und `--diff` werden mit Exit-Code 2 abgewiesen |
| `devcontainer up` | `--recreate`; `--dry-run` und `--diff` werden mit Exit-Code 2 abgewiesen |
| `devcontainer down` | keine eigenen; `--dry-run` und `--diff` werden mit Exit-Code 2 abgewiesen |

### Dienstekatalog

| Dienst | Image | Host-Port(s) | Zugangsdaten in `.env` |
|---|---|---|---|
| `postgres` | `postgres:16-alpine` | 5432 | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` |
| `keycloak` | `quay.io/keycloak/keycloak:26.0` | 8080 | `KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD` |
| `otel` | `otel/opentelemetry-collector:0.108.0` | 4317 (gRPC), 4318 (HTTP) | keine; Konfiguration in `otel-collector-config.yaml` |

Die Werte `CHANGEME_…` in `.env.example` sind Platzhalter. Ersetzen Sie sie in
`.env` durch eigene Werte.

### Dateien eines Projekts

| Datei oder Ordner | Inhalt | Wer pflegt sie |
|---|---|---|
| `u-boot.yaml` | Projektdatei (Name, Dienste, Devcontainer) | `u-boot` (`config set`, `add`, `remove`) |
| `compose.yaml` | Docker-Compose-Stack | `u-boot` im verwalteten Bereich, Sie außerhalb |
| `.env.example` | Vorlage der Umgebungsvariablen | `u-boot` im verwalteten Bereich |
| `.env` | Ihre echten Werte | nur Sie (nicht im Git) |
| `README.md`, `CHANGELOG.md` | Begleitdokumente | `u-boot generate …` im verwalteten Bereich, Sie außerhalb |
| `otel-collector-config.yaml` | Collector-Konfiguration (nur mit `otel`) | `u-boot add otel`, danach Sie |
| `.devcontainer/` | Devcontainer (nur mit Devcontainer) | `u-boot generate devcontainer` |
| `docker/`, `scripts/`, `docs/` | Arbeitsordner für Ihre Dateien | Sie |

### Grenzwerte und Mindestversionen

| Gegenstand | Wert |
|---|---|
| Docker Engine | ab 24.0 |
| Docker-Compose-Plugin | ab 2.20 |
| Podman (als Ersatz) | ab 4.0 |
| Projektname | bis 63 Zeichen: `a–z`, `0–9`, `-`; Anfang Buchstabe, Ende Buchstabe oder Ziffer |
| `up --timeout` | Standard 60 Sekunden; `0` = nicht warten; negativ ist ungültig |
| `devcontainer.user.uid` | 1 bis 65535, Standard 1000 |
| Sicherungsdateien (`--backup`) | `<name>.bak`, `<name>.bak.1`, `<name>.bak.2`, … |

### JSON-Ausgabe

Mit `--json` gibt jeder Befehl genau ein JSON-Objekt auf der Standardausgabe
aus. Bei Fehlern erscheint der Fehlertext zusätzlich auf der Fehlerausgabe.

**Immer vorhanden:**

| Feld | Inhalt |
|---|---|
| `status` | `ok`, `warn` oder `error`; folgt dem schwersten Eintrag in `diagnostics` |
| `command` | Befehl: `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `egress`, `devcontainer`, `template` |
| `subcommand` | nur bei `config` (`show` für `u-boot config` ohne Unterbefehl, `get`, `set`, `list`), `egress` (`list`), `devcontainer` (`up`, `down`) und `template` (`list`) |
| `diagnostics` | Liste der Befunde, leer (`[]`), wenn alles in Ordnung ist |
| `exitCode` | derselbe Wert wie der Exit-Code des Prozesses |
| `data` | Ergebnisdaten des Befehls (siehe Tabelle unten), nicht bei jedem Befehl |

Jeder Eintrag in `diagnostics` hat `level` (`warn` oder `error`), `code`,
`message` und optional `file`. Meldungen mit dem Zustand „in Ordnung“ erscheinen
im JSON nicht; `--quiet` ändert die JSON-Ausgabe nicht.

**Zusätzlich bei `--dry-run` und `--diff`** (Vorschau von schreibenden Befehlen):

| Feld | Inhalt |
|---|---|
| `dryRun` | `true`, wenn nichts geschrieben wurde |
| `diff` | `true`, wenn Diffs enthalten sind |
| `plannedFiles` | Liste mit `path` und `action` (`create`, `modify` oder `delete`); mit `--diff` zusätzlich `hunks` |
| `changes` | Liste mit `path` und `count` (Anzahl geänderter Zeilen oder Einträge) |

**Beispiele:**

```json
{"status":"ok","command":"add","dryRun":true,"diff":false,"plannedFiles":[{"path":"u-boot.yaml","action":"modify"},{"path":"compose.yaml","action":"modify"},{"path":".env.example","action":"modify"}],"changes":[{"path":"u-boot.yaml","count":4},{"path":"compose.yaml","count":27},{"path":".env.example","count":6}],"diagnostics":[],"exitCode":0}
```

```json
{"status":"warn","command":"doctor","diagnostics":[{"level":"warn","code":"uboot.yaml.valid","message":"u-boot.yaml not present — directory is not a u-boot project."}],"exitCode":0}
```

**`data` je Befehl:**

| Befehl | `data` |
|---|---|
| `init`, `add`, `doctor` | keine |
| `remove` | `service`, `priorState`, `state` (zum Beispiel `deactivated`), `volumesPurged` |
| `up` | `services`: Liste mit `name`, `state`, `port`, `ports`, `healthcheck`; auch bei einem Fehler nach dem Start (Stand der bis dahin gestarteten Dienste) |
| `down` | `removedVolumes` (Wahrheitswert), `removedVolumeNames` (Liste) |
| `logs` | `lines`: Liste der Log-Zeilen (nur mit `--tail`, nicht mit `--follow`) |
| `generate` | `artifact` (zum Beispiel `readme`), `action` (zum Beispiel `created` oder `no-op`) |
| `config get` | `path` und `value`; bei mehreren Pfaden oder `--json-array` `entries` |
| `config list` | `entries`: Liste mit `path` und `value` |
| `config set` | `path`, `oldValue`, `newValue`, `noOp` |
| `egress list` | `active`, `sandboxProfile`, `egressEnabled`, `cloneSource` (Wahrheitswerte) und `hosts`: Liste mit `host` und `sources` (Herkunft) |
| `devcontainer up` | `containerId` (volle ID), `remoteUser`, `remoteWorkspaceFolder`, `alreadyRunning` (Wahrheitswert; `true`, wenn der Container vor dem Aufruf lief und die Dev-Container-CLI dieselbe Kennung meldet); nicht bei einem Fehler |
| `devcontainer down` | `containerId` (leer, wenn es keinen Container gibt) und `stopped` (Wahrheitswert; `false`, wenn nichts zu stoppen war) |
| `template list` | Liste der Vorlagen mit `name`, `description`, `version`, `supportedAddOns`, `generatedFiles`, `requiredTools`, `variables` |

Bei Fehlern von `config get` und `config set` enthält `data.hint` einen
Reparaturvorschlag: `command` (der Befehl zum Kopieren), `action` und `argument`.

**Besonderheiten:**

- Ein Aufruf von `u-boot template --json` oder `u-boot egress --json` ohne `list`
  endet mit Exit-Code 2 und gibt kein JSON aus; `u-boot devcontainer` ohne `up`
  oder `down` ebenso, auch ohne `--json`. `u-boot logs --follow --json` ist ebenfalls ein
  Nutzungsfehler (Exit-Code 2).
- Der Umschlag von `devcontainer up` trägt auch im Erfolgsfall `status: warn`, weil die
  Ausweisung des Docker-Sockets (Abschnitt 4.9) ein Eintrag `warn` in `diagnostics` ist;
  der Exit-Code ist `0`. Bei einem Fehler, nachdem der Hilfs-Container gestartet werden sollte,
  steht diese Warnung vor dem Fehler-Eintrag.
- `doctor` endet mit `warn` und Exit-Code 0, solange nur Warnungen vorliegen;
  mit `--strict` oder bei einem Fehler lautet der Exit-Code 11.
- Bei `down --volumes` ohne `--yes` wird im JSON-Modus nicht gefragt, sondern
  abgelehnt (Exit-Code 10).

**Kennungen in `diagnostics[].code`:**

- Prüfungen von `doctor` tragen den Namen der Prüfung (siehe Tabelle).
- Bei anderen Fehlern ist `code` eine interne Kennung der Fehlerklasse. Werten Sie
  für die Fehlerbehandlung in Skripten den `exitCode` aus und zeigen Sie `message`
  an.

| Prüfung (`code`) | Bedeutung |
|---|---|
| `fs.write-permissions` | Schreibrecht im Verzeichnis |
| `git.installed` | `git` vorhanden |
| `docker.installed` | Docker vorhanden und neu genug |
| `docker.reachable` | Docker-Daemon erreichbar |
| `docker.compose.installed` | Compose-Plugin vorhanden und neu genug |
| `uboot.yaml.valid` | `u-boot.yaml` gültig |
| `compose.yaml.valid` | `compose.yaml` gültig |
| `services.enabled-key` | alle Dienste tragen einen `enabled`-Schlüssel |
| `devcontainer.json.valid` | `devcontainer.json` gültig |
| `devcontainer.dockerfile.valid` | `.devcontainer/Dockerfile` lesbar |
| `devcontainer.forwardPorts.consistency` | Weiterleitungs-Ports passen zu den Diensten |
| `devcontainer.features.allowlist` | externe Features sind freigegeben; im Sandbox-Profil ist das Feature `docker-cli` nicht aktiv |
| `devcontainer.features.drift` | Features in `u-boot.yaml` und `devcontainer.json` stimmen überein |
| `devcontainer.sandbox.runtime` | Sandbox: verschachtelte Runtime (Podman oder Docker) passend konfiguriert |
| `devcontainer.sandbox.egress` | Sandbox: Internet-Beschränkung passend konfiguriert |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Zugangsdaten (Git-Token, Anthropic-Schlüssel), Quelle der Git- und Agenten-Zugangsdaten vorhanden |

### Unterstützte Plattformen

Binary: Linux, macOS und Windows, je `amd64` und `arm64`. Homebrew: macOS,
`amd64` und `arm64`. Pakete: Linux `amd64` und `arm64` (`.deb`,
`.rpm`). Container-Image: `ghcr.io/pt9912/u-boot`.

---

## 11. Support und Lizenz

**Fragen und Fehler:** <https://github.com/pt9912/u-boot/issues>

Legen Sie einem Fehlerbericht bei: die Ausgabe von `u-boot --version`, die
Ausgabe von `u-boot doctor`, Ihr Betriebssystem und den genauen Befehl. Entfernen
Sie vorher Passwörter, Tokens und interne Adressen aus der Ausgabe.

**Lizenz:** MIT, siehe `LICENSE` im Projektarchiv.

**Gültigkeitsbereich:** Dieses Handbuch beschreibt `u-boot` v0.13.1. Prüfen Sie
mit `u-boot --version`, welche Version Sie einsetzen. Weicht sie ab, zeigt
`u-boot <befehl> --help` den Stand Ihrer Version.
