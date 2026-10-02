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
- `generate devcontainer` hat zusätzliche Optionen (siehe 4.9 bis 4.11).
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
- Das Standardprofil bindet den Projektordner ein, aber weder das Docker-Socket
  des Hosts noch Geheimnis-Verzeichnisse wie `~/.ssh` oder `~/.aws`.
- `u-boot doctor` meldet, wenn `devcontainer.json` und `u-boot.yaml`
  auseinanderlaufen (Warnung); `u-boot generate devcontainer` behebt das.
- Werkzeuge im Container (Node.js, Java, Go, …) fügen Sie als Features hinzu
  (siehe 4.10). Für autonome Agenten gibt es ein eigenes Profil (siehe 4.11).

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

Ohne weitere Freigabe aktivierbar. Die Standardversion ist `1`.

| Name | Inhalt |
|---|---|
| `git` | Git |
| `docker-cli` | Docker-CLI (nutzt den Docker des Hosts: „docker-outside-of-docker“) |
| `node` | Node.js |
| `java` | Java und SDKMAN |
| `go` | Go-Toolchain |
| `cpp` | C++-Toolchain |
| `kubectl-helm` | kubectl, helm und minikube |
| `postgres-client` | PostgreSQL-Client |

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

#### Prüfungen durch `u-boot doctor`

| Schweregrad | Auslöser | Abhilfe |
|---|---|---|
| Fehler | `source` ist gesetzt, steht aber nicht in `featureSources.allow` | URL freigeben (siehe oben) |
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
   dem Host.

#### Ergebnis

`--sandbox` setzt `devcontainer.profile: sandbox` in `u-boot.yaml`.

| Eigenschaft | Sandbox-Profil | Standardprofil |
|---|---|---|
| Workspace | benanntes Volume `<projekt>-workspace-${devcontainerId}`; das Repository wird im Container geklont | Bind-Mount des Projektordners |
| Docker-Socket des Hosts | nicht eingebunden | nicht eingebunden |
| `--privileged` und zusätzliche Fähigkeiten | nein (außer mit verschachteltem Podman) | nein |
| Geheimnis-Verzeichnisse des Hosts (`~/.ssh`, `~/.aws`, …) | nicht eingebunden | nicht eingebunden |
| Benutzer | `vscode`, nicht root | `vscode`, nicht root |

#### Hinweise

- `--sandbox` ohne `--devcontainer` (bei `init`) oder bei einem anderen Artefakt
  als `devcontainer` (bei `generate`) ist ein Fehler.

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

Im Container läuft dann rootless Podman mit einem `docker`-Alias, sodass
`docker build` und `docker run` ohne Host-Socket möglich sind. Unter Docker sind
dafür diese Lockerungen nötig; **`u-boot` weist jede einzeln als Warnung aus**:

```text
--cap-add=SYS_ADMIN
--security-opt=seccomp=unconfined
--security-opt=apparmor=unconfined
--security-opt=systempaths=unconfined
--device=/dev/fuse
```

Das ist nahe an `--privileged`; der Schutz ist damit deutlich schwächer. Das
Standardprofil (`nestedRuntime: none`) bleibt ungelockert.

Was beim Containerstart passiert, wenn Fähigkeiten fehlen, steuert
`devcontainer.sandbox.onUnavailable` (`warn` ist der Standard, `fail` bricht ab):

| Zustand | `warn` | `fail` |
|---|---|---|
| `/dev/fuse` fehlt | Rückfall auf `vfs`-Speicher, Warnung | Abbruch mit Exit-Code 11 |
| Verschachtelte User-Namespaces blockiert (Seccomp oder AppArmor) | Abbruch mit Exit-Code 11 | Abbruch mit Exit-Code 11 |

Unter **Ubuntu 24.04** (und auf Hosts mit
`kernel.apparmor_restrict_unprivileged_userns=1`) blockiert dieser Schalter
verschachtelte User-Namespaces auch mit `apparmor=unconfined`; der Start endet
dann mit Exit-Code 11. Abhilfe: auf dem Container-Host
`sysctl kernel.apparmor_restrict_unprivileged_userns=0` setzen oder
`nestedRuntime: none` wählen.

`u-boot doctor` prüft `/dev/fuse` nur auf Linux-Hosts. Unter macOS mit Colima
läuft die Engine in einer VM, die der Host nicht beurteilen kann; dort prüft das
Startscript im Container.

#### Internet-Beschränkung (Egress)

```bash
u-boot config set devcontainer.sandbox.egress.enabled true
u-boot config set devcontainer.sandbox.egress.allow api.example.org,registry.npmjs.org
u-boot generate devcontainer
```

Beim Containerstart startet ein lokaler DNS-Server, der **nur erlaubte Namen**
auflöst. Eine Firewall lässt nur Loopback, bestehende Verbindungen, DNS und die
Adressen der erlaubten Namen durch; alles andere (IPv6 ganz) wird verworfen. Das
ist eine Leitplanke, keine Sandbox-Grenze. Der Container braucht dafür die
Fähigkeit `NET_ADMIN`, die `u-boot` als Lockerung ausweist.

Standardmäßig erlaubt sind:

| Quelle | Hosts |
|---|---|
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` |
| mit `nestedRuntime: podman` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io` |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Clone-Quelle | der Host von `devcontainer.sandbox.repository` beziehungsweise `origin` |
| Ihre Ergänzungen | `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema, Port oder Wildcard) |

- **Wichtig:** Die Hosts des Agenten selbst (zum Beispiel dessen API) kennt
  `u-boot` nicht. Tragen Sie sie in `egress.allow` ein. Subdomains eines
  erlaubten Namens sind mitgemeint.
- Nach einer Änderung führen Sie `u-boot generate devcontainer` aus und starten
  den Container neu.
- Wer die Fähigkeit hat und das `sudo` des Base-Images nutzen kann, kann die
  Regeln aufheben. Der Verkehr verschachtelter Podman-Container wird nicht
  erfasst.
- Ohne `NET_ADMIN` entfällt die Beschränkung mit einer Warnung (`warn`) oder das
  Startscript endet mit Exit-Code 11 (`fail`). `u-boot doctor` prüft nur die
  Konfiguration; ohne Sandbox-Profil ist sie wirkungslos und `doctor` warnt.

#### Git-Zugangsdaten

Zugangsdaten stehen **nie** im Image, im Volume, in `u-boot.yaml` oder in einer
erzeugten Datei:

- `devcontainer.json` enthält nur die Referenz
  `"GIT_TOKEN": "${localEnv:GIT_TOKEN}"`. Setzen Sie `GIT_TOKEN` in der
  Host-Umgebung, bevor Sie den Container starten.
- Ein Credential-Helper im Image liefert `$GIT_TOKEN` an `git` (nur für
  `https`-Remotes; bei `ssh` gibt es keinen Schlüssel im Container).
- Empfohlen: **kurzlebige, auf das Repository begrenzte Tokens**, nie ein
  privater Schlüssel im Container, und **Branch-Protection** im Remote, damit ein
  kompromittierter Agent nicht auf geschützte Branches schreiben kann.
- `u-boot doctor` warnt vor möglichen Klartext-Tokens in `u-boot.yaml`,
  `compose.yaml`, `.env.example` und `.devcontainer/*` sowie vor einer fehlenden
  Token-Quelle bei einem `https`-Clone.

#### Mehrere Instanzen

Die Volumes heißen `<projekt>-workspace-${devcontainerId}` und (mit
verschachteltem Podman) `<projekt>-containers-${devcontainerId}`. Das
Dev-Containers-Werkzeug löst `devcontainerId` pro Projektordner auf.

- Parallele Instanzen desselben Projekts starten Sie aus **getrennten Ordnern**
  (Klone oder `git worktree`); jede Instanz hat eigene Volumes.
- Derselbe Ordner bleibt eine Instanz.
- Die Compose-Ports (`u-boot up`) mehrerer Instanzen kollidieren weiterhin.
- Alte Volumes `<projekt>-workspace` und `<projekt>-containers` aus früheren
  Entwicklungsständen bleiben verwaist; entfernen Sie sie mit `docker volume rm`.

#### Grenzen

- Geprüft unter Docker auf Linux mit der `devcontainer`-CLI. **Nicht geprüft:**
  macOS mit Colima, Podman als Host-Engine, VS Code und Codespaces.
- Der Wechsel von `podman` zurück auf `none` lässt eine vorhandene
  `sandbox-init.sh` im Ordner liegen; das Dockerfile verwendet sie dann nicht
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
3. Weitere Abhilfen stehen in Abschnitt [4.11](#411-einen-sandbox-devcontainer-für-autonome-agenten-einrichten) unter „Verschachteltes Podman".

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
`u-boot <befehl> --help` zeigt sie; die Optionen aller Befehle stehen auch im
[Anhang](#optionen-je-befehl). Die JSON-Ausgabe beschreibt der Anhang
[JSON-Ausgabe](#json-ausgabe).

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
| `command` | Befehl: `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config`, `template` |
| `subcommand` | nur bei `config` (`get`, `set`, `list`) und `template` (`list`) |
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
{"status":"error","command":"add","diagnostics":[{"level":"error","code":"LH-FA-ADD-002","message":"service not supported: \"redis\" is not in the built-in catalogue [postgres keycloak otel]"}],"exitCode":10}
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
| `template list` | Liste der Vorlagen mit `name`, `description`, `version`, `supportedAddOns`, `generatedFiles`, `requiredTools`, `variables` |

Bei Fehlern von `config get` und `config set` enthält `data.hint` einen
Reparaturvorschlag: `command` (der Befehl zum Kopieren), `action` und `argument`.

**Besonderheiten:**

- Ein Aufruf von `u-boot template --json` ohne `list` endet mit Exit-Code 2 und
  gibt kein JSON aus. `u-boot logs --follow --json` ist ebenfalls ein
  Nutzungsfehler (Exit-Code 2).
- `doctor` endet mit `warn` und Exit-Code 0, solange nur Warnungen vorliegen;
  mit `--strict` oder bei einem Fehler lautet der Exit-Code 11.
- Bei `down --volumes` ohne `--yes` wird im JSON-Modus nicht gefragt, sondern
  abgelehnt (Exit-Code 10).

**Kennungen in `diagnostics[].code`:**

- Prüfungen von `doctor` tragen den Namen der Prüfung (siehe Tabelle).
- Andere Fehler tragen eine feste Kennung, die je Fehlerklasse stabil bleibt (im
  zweiten Beispiel oben für einen unbekannten Dienst). Für die Fehlerbehandlung in
  Skripten genügt in der Regel der `exitCode`.

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
| `devcontainer.features.allowlist` | externe Features sind freigegeben |
| `devcontainer.features.drift` | Features in `u-boot.yaml` und `devcontainer.json` stimmen überein |
| `devcontainer.sandbox.runtime` | Sandbox: verschachteltes Podman passend konfiguriert |
| `devcontainer.sandbox.egress` | Sandbox: Internet-Beschränkung passend konfiguriert |
| `devcontainer.sandbox.credentials` | Sandbox: keine Klartext-Zugangsdaten, Token-Quelle vorhanden |

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

**Lizenz:** MIT, siehe `LICENSE` im Projektarchiv.

**Gültigkeitsbereich:** Dieses Handbuch beschreibt `u-boot` v0.7.0. Prüfen Sie
mit `u-boot --version`, welche Version Sie einsetzen. Weicht sie ab, zeigt
`u-boot <befehl> --help` den Stand Ihrer Version.
