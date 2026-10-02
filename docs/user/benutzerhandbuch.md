# Benutzerhandbuch: u-boot

Handbuch-Version: 2.0
Software-Version: v0.7.0
Stand: 2026-10-02

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

Zusätzliche Anleitungen für Spezialthemen:

- [`examples.md`](examples.md) – Beispielabläufe als Kommando-Rezepte
- [`devcontainer-features.md`](devcontainer-features.md) – Devcontainer-Features
- [`devcontainer-sandbox.md`](devcontainer-sandbox.md) – Sandbox-Devcontainer für autonome Agenten
- [`cli-json-output.md`](cli-json-output.md) – maschinenlesbare Ausgabe

### Voraussetzungen

| Voraussetzung | Wofür |
|---|---|
| `u-boot` (Binary, Paket, Homebrew oder Container-Image) | alle Befehle |
| Docker Engine ab 24.0 | `up`, `down`, `logs` und die Docker-Prüfungen von `doctor` |
| Docker-Compose-Plugin ab 2.20 | dieselben Befehle |
| `git` | `init` (außer mit `--no-git`) |

**Ohne Docker** funktionieren `init`, `add`, `remove`, `generate`, `config` und
`template list`. Nur `up`, `down` und `logs` brauchen einen erreichbaren
Docker-Daemon. `u-boot doctor` meldet fehlende Werkzeuge als Befund.

Podman ab 4.0 funktioniert als Ersatz für Docker, wenn `docker` auf `podman`
zeigt und `DOCKER_HOST` auf den Podman-Socket gesetzt ist. `u-boot doctor`
meldet Podman-Versionen als Warnung („unrecognized version"), bricht aber nicht
ab.

---

## 2. Installation

Wählen Sie **eine** der folgenden Möglichkeiten. Danach prüfen Sie immer mit

```bash
u-boot --version
```

**Ergebnis:** Die installierte Version erscheint, zum Beispiel
`u-boot version 0.7.0`.

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
`download/v0.7.0`. Die Prüfsummen aller Dateien stehen in `SHA256SUMS` am
jeweiligen Release.

### Homebrew (macOS und Linux)

Für stabile Releases gibt es einen Homebrew-Tap für macOS und Linux (je `amd64`
und `arm64`; Windows trägt Homebrew nicht):

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
VERSION=0.7.0
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
  mit `?` markiert. Auch `up`, `down` und `logs` laufen im Image nicht, weil dort
  kein Docker vorhanden ist.

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

### Allgemeine Optionen

Diese Optionen gelten für jeden Befehl:

| Option | Wirkung |
|---|---|
| `--yes` | Rückfragen automatisch bejahen |
| `--no-interactive` | nie fragen; bei einer nötigen Rückfrage abbrechen |
| `--json` | maschinenlesbare Ausgabe |
| `--quiet` | nur das Nötigste; Warnungen und Fehler bleiben |
| `--verbose` | zusätzliche Details |
| `--debug` | interne Diagnoseausgabe |
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
  Sandbox-Profil (siehe 4.10).
- `--template <name>` rendert das Projekt aus einer Vorlage (siehe 4.11).

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
- Braucht ein Dienst einen anderen, installiert `--with-deps` das Fehlende
  automatisch mit, ohne Rückfrage.
- `.env.example` bleibt eine Vorlage. Echte Zugangsdaten gehören nur in `.env`,
  das nicht ins Git-Repository gehört.
- Beispielabläufe für die einzelnen Dienste stehen in
  [`examples.md`](examples.md).

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
`✗` Fehler, `?` übersprungen. Er endet mit einer Zusammenfassung:

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
vorhanden), Devcontainer-Features (Freigabeliste und Abweichungen zur
Projektdatei), Weiterleitungs-Ports sowie die drei Sandbox-Prüfungen (Laufzeit,
Zugangsdaten, Egress).

#### Hinweise

- Bei einem Fehler endet `doctor` mit Exit-Code 11; mit `--strict` auch bei
  einer Warnung. `--strict` ist für CI gedacht.
- Läuft `u-boot` im Container, werden die vier Host-Prüfungen übersprungen
  (`?`). Das ist kein Fehler.

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
- `generate devcontainer` hat zusätzliche Optionen (siehe 4.9 und 4.10).
  `--dry-run` und `--diff` stehen bei allen vier Artefakten zur Verfügung.

### 4.9 Einen Devcontainer einrichten

#### Voraussetzung

Ein `u-boot`-Projekt. Zum Starten des Devcontainers brauchen Sie einen
Devcontainer-fähigen Editor (zum Beispiel VS Code) oder die
`devcontainer`-CLI.

#### Vorgehen

1. Neues Projekt mit Devcontainer:

   ```bash
   u-boot init mein-projekt --devcontainer
   ```

   Bestehendes Projekt:

   ```bash
   u-boot generate devcontainer
   ```

2. Optional: ein Feature aktivieren und neu erzeugen:

   ```bash
   u-boot config set devcontainer.features.node.enabled true
   u-boot generate devcontainer
   ```

3. Optional: eine andere Benutzer-ID im Container (zum Beispiel unter macOS mit
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

- Verfügbare eingebaute Features: `git`, `docker-cli`, `node`, `java`, `go`,
  `cpp`, `kubectl-helm`, `postgres-client`. Eine abweichende Version setzen Sie
  mit `devcontainer.features.<name>.version`.
- Externe Feature-Quellen brauchen eine ausdrückliche Freigabe in
  `devcontainer.featureSources.allow` (`--yes` ersetzt sie nicht). Die
  Einzelheiten stehen in [`devcontainer-features.md`](devcontainer-features.md).
- Zulässige Benutzer-IDs: 1 bis 65535 (Standard 1000). Die ID 0 (root) wird
  abgelehnt.
- `u-boot doctor` meldet, wenn `devcontainer.json` und `u-boot.yaml`
  auseinanderlaufen (Warnung); `u-boot generate devcontainer` behebt das.

### 4.10 Einen Sandbox-Devcontainer für autonome Agenten einrichten

#### Voraussetzung

Ein `u-boot`-Projekt, am besten mit einem Git-Remote `origin`. Das Sandbox-Profil
ist **Schadensbegrenzung, keine harte Isolationsgrenze**.

#### Vorgehen

1. Profil aktivieren:

   ```bash
   u-boot init mein-projekt --devcontainer --sandbox     # neues Projekt
   u-boot generate devcontainer --sandbox                # bestehendes Projekt
   ```

2. Optional: ein anderes Repository im Container klonen:

   ```bash
   u-boot config set devcontainer.sandbox.repository git@github.com:org/anderes-repo.git
   u-boot generate devcontainer
   ```

3. Optional: den Zugriff ins Internet auf eine Liste erlaubter Hosts beschränken:

   ```bash
   u-boot config set devcontainer.sandbox.egress.enabled true
   u-boot config set devcontainer.sandbox.egress.allow api.example.org,registry.npmjs.org
   u-boot generate devcontainer
   ```

4. Optional: Container-Programme ohne Zugriff auf den Host-Docker (verschachteltes
   Podman):

   ```bash
   u-boot config set devcontainer.sandbox.nestedRuntime podman
   u-boot generate devcontainer
   ```

5. Zugangsdaten für `https`-Remotes geben Sie zur Laufzeit mit, nicht in Dateien:
   Setzen Sie vor dem Start des Containers die Umgebungsvariable `GIT_TOKEN` auf
   dem Host.

#### Ergebnis

Der Workspace liegt in einem benannten Volume statt in einem Bind-Mount; das
Repository wird beim ersten Start im Container geklont. Das Docker-Socket des
Hosts, `--privileged` und Ihre Geheimnis-Verzeichnisse (`~/.ssh`, `~/.aws`) sind
**nicht** eingebunden. Mit `nestedRuntime: podman` weist `u-boot` jede nötige
Lockerung einzeln als Warnung aus; der Schutz ist dann deutlich schwächer als im
Standardprofil.

#### Hinweise

- `--sandbox` ohne `--devcontainer` (bei `init`) oder bei einem anderen Artefakt
  als `devcontainer` (bei `generate`) ist ein Fehler.
- Repository-Adressen mit Zugangsdaten (`https://token@…`) oder unsicheren
  Zeichen werden abgelehnt (Exit-Code 10); es wird nichts geschrieben. Erlaubt
  sind `git@host:pfad`, `ssh://git@host/pfad` und `https://` ohne Benutzerangabe.
- Verwenden Sie kurzlebige, auf das Repository begrenzte Tokens und schützen Sie
  wichtige Branches im Remote.
- Mehrere Instanzen starten Sie aus getrennten Ordnern (Klone oder
  `git worktree`); jede hat eigene Volumes.
- Vollständige Beschreibung mit Grenzen und Sonderfällen:
  [`devcontainer-sandbox.md`](devcontainer-sandbox.md).

### 4.11 Vorlagen nutzen

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

### 4.12 In Skripten und CI verwenden

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

Das vollständige Schema samt Exit-Code-Matrix je Befehl steht in
[`cli-json-output.md`](cli-json-output.md).

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
u-boot config get project.name devcontainer.enabled    # mehrere, je Zeile ein Wert
u-boot config list                                     # alle gesetzten Pfade als pfad=wert
```

Beispielausgabe von `u-boot config list`:

```text
project.name=demo-app
services.keycloak.enabled=true
services.postgres.enabled=true
```

Mit `--json` liefern `config get` bei mehreren Pfaden (oder mit `--json-array`)
und `config list` die Form `data.entries[{path, value}]`. Ein nicht gesetzter
oder unbekannter Pfad endet mit Exit-Code 10.

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
| `devcontainer.sandbox.nestedRuntime` | Container-Laufzeit im Sandbox | `none` / `podman` | ja |
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

### Ausführlichkeit der Ausgabe

| Option | Wirkung |
|---|---|
| `--quiet` | nur das Nötigste; Warnungen und Fehler bleiben |
| `--verbose` | zusätzliche Details |
| `--debug` | interne Diagnoseausgabe |

---

## 6. Rollen und Rechte

`u-boot` kennt keine Benutzerrollen. Es braucht nur die Rechte, die auch die
Handarbeit bräuchte:

| Aufgabe | Benötigtes Recht |
|---|---|
| Dateien anlegen und ändern (`init`, `add`, `remove`, `generate`, `config set`) | Schreibrecht im Projektverzeichnis (`doctor` prüft das: `fs.write-permissions`) |
| `up`, `down`, `logs` | Zugriff auf den Docker-Daemon; unter Linux meist Mitgliedschaft in der Gruppe `docker` |
| Installation in `/usr/local/bin` oder per Paket | Administratorrechte (`sudo`) |
| `git`-Funktionen von `init` | `git` installiert und konfiguriert |

Hinweis: Im Devcontainer läuft Ihre Arbeit als Benutzer `vscode`, nicht als
root. Weder das Standard- noch das Sandbox-Profil bindet das Docker-Socket des
Hosts oder Ihre Geheimnis-Verzeichnisse ein.

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
| `11` | Umgebungsproblem (Docker nicht erreichbar, Compose-Plugin fehlt; auch `doctor` mit Fehlern) |
| `12` | Laufzeitfehler von Compose (auch Zeitüberschreitung bei `up`) |
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

### Fehler: Sandbox-Container startet nicht (Exit-Code 11 im Startscript)

#### Ursache

Verschachteltes Podman (`nestedRuntime: podman`) braucht Fähigkeiten, die der
Host nicht bietet, zum Beispiel blockierte User-Namespaces unter Ubuntu 24.04
oder fehlendes `/dev/fuse`.

#### Lösung

1. Lesen Sie die Meldung des Startscripts im Container-Log.
2. Setzen Sie `devcontainer.sandbox.nestedRuntime` auf `none` und erzeugen Sie
   den Devcontainer neu, wenn Sie kein verschachteltes Podman brauchen.
3. Weitere Abhilfen stehen in [`devcontainer-sandbox.md`](devcontainer-sandbox.md).

### Wenn nichts davon hilft

1. Wiederholen Sie den Befehl mit `--debug`.
2. Führen Sie `u-boot doctor` aus und halten Sie die Ausgabe bereit.
3. Notieren Sie Version (`u-boot --version`), Betriebssystem und den genauen
   Befehl.
4. Eröffnen Sie ein Issue (siehe [Support und Lizenz](#11-support-und-lizenz)).

---

## 8. FAQ

**Brauche ich Docker, um `u-boot` zu benutzen?**
Nur für `up`, `down`, `logs` und die Docker-Prüfungen von `doctor`. `init`,
`add`, `remove`, `generate`, `config` und `template list` laufen ohne.

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
Nein. Es begrenzt Schäden (kein Host-Docker-Socket, kein `--privileged`, optional
eine Liste erlaubter Internet-Ziele), ist aber keine harte Sicherheitsgrenze.
Das gilt besonders mit verschachteltem Podman.

**Wo finde ich alle Optionen eines Befehls?**
`u-boot <befehl> --help` zeigt sie. Beispielabläufe stehen in
[`examples.md`](examples.md), die JSON-Ausgabe in
[`cli-json-output.md`](cli-json-output.md).

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
| `u-boot generate <artefakt>` | Artefakt erzeugen oder aktualisieren | 4.8–4.10 |
| `u-boot config [get\|set\|list]` | Konfiguration lesen und ändern | 5 |
| `u-boot template list` | Vorlagenkatalog anzeigen | 4.11 |
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

### Unterstützte Plattformen

Binary: Linux, macOS und Windows, je `amd64` und `arm64`. Homebrew: macOS und
Linux, je `amd64` und `arm64`. Pakete: Linux `amd64` und `arm64` (`.deb`,
`.rpm`). Container-Image: `ghcr.io/pt9912/u-boot`.

---

## 11. Support und Lizenz

**Fragen und Fehler:** <https://github.com/pt9912/u-boot/issues>

Legen Sie einem Fehlerbericht bei: die Ausgabe von `u-boot --version`, die
Ausgabe von `u-boot doctor`, Ihr Betriebssystem und den genauen Befehl. Entfernen
Sie vorher Passwörter, Tokens und interne Adressen aus der Ausgabe.

**Weiterführende Dokumentation:**

- [`examples.md`](examples.md) – Beispielabläufe als Kommando-Rezepte
- [`cli-json-output.md`](cli-json-output.md) – JSON-Schema und Exit-Code-Matrix
- [`devcontainer-features.md`](devcontainer-features.md) – Devcontainer-Features
- [`devcontainer-sandbox.md`](devcontainer-sandbox.md) – Sandbox-Profil für autonome Agenten

**Lizenz:** MIT, siehe `LICENSE` im Projektarchiv.

**Gültigkeitsbereich:** Dieses Handbuch beschreibt `u-boot` v0.7.0. Prüfen Sie
mit `u-boot --version`, welche Version Sie einsetzen. Weicht sie ab, zeigt
`u-boot <befehl> --help` den Stand Ihrer Version.
