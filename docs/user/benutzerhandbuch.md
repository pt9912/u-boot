# Benutzerhandbuch: u-boot

Handbuch-Version: 1.0
Software-Version: v0.4.0
Stand: 2026-07-25

---

## Inhalt

1. [Einleitung](#1-einleitung)
2. [Installation](#2-installation)
3. [Erste Schritte](#3-erste-schritte)
4. [Aufgaben](#4-aufgaben)
5. [Konfiguration](#5-konfiguration)
6. [Fehlerbehebung](#6-fehlerbehebung)
7. [FAQ](#7-faq)
8. [Glossar](#8-glossar)
9. [Support und Lizenz](#9-support-und-lizenz)
10. [Änderungshistorie](#10-änderungshistorie)

---

## 1. Einleitung

### Zweck der Software

`u-boot` legt Entwicklungsumgebungen an, die auf jedem Rechner gleich
aussehen. Statt eine Projektstruktur, eine `compose.yaml`, eine
Devcontainer-Konfiguration und die üblichen Begleitdateien von Hand zu
schreiben, erzeugt und pflegt `u-boot` sie für Sie — und fügt Dienste wie
PostgreSQL oder Keycloak auf Zuruf hinzu.

Der Kern ist ein einzelner Befehl: `u-boot init` erzeugt ein vollständiges,
lauffähiges Projekt. Danach begleitet Sie dasselbe Werkzeug durch den Alltag:
Dienste hinzufügen und entfernen, die Umgebung starten und stoppen, Logs
ansehen, Artefakte wie README oder CHANGELOG aktualisieren.

`u-boot` ist ausdrücklich **kein** Deployment-Werkzeug und kein Ersatz für
Docker Compose. Es erzeugt und pflegt die Dateien, mit denen Sie arbeiten; das
Starten übernimmt darunter weiterhin Compose.

### Zielgruppe dieses Handbuchs

Entwicklerinnen und Entwickler, die ein neues Projekt aufsetzen oder ein
bestehendes `u-boot`-Projekt betreuen. Sie brauchen kein Vorwissen über den
inneren Aufbau von `u-boot`. Grundkenntnisse in der Kommandozeile und ein
grobes Verständnis von Docker Compose genügen.

Wenn Sie an `u-boot` selbst mitentwickeln, ist [`quality.md`](quality.md) und
die Spezifikation unter `spec/` der richtige Einstieg — nicht dieses Handbuch.

### Voraussetzungen

| Voraussetzung | Wofür |
|---|---|
| `u-boot`-Binary oder Container-Image | alle Befehle |
| Docker Engine (ab 24.0) | `up`, `down`, `logs` und die Prüfungen von `doctor` |
| Docker-Compose-Plugin (ab 2.20) | dieselben Befehle |
| `git` | nur für `init` ohne `--no-git` |

**Ohne Docker** funktionieren `init`, `generate`, `config` und `template list`
vollständig. Nur die Befehle, die eine laufende Umgebung brauchen, setzen einen
erreichbaren Docker-Daemon voraus.

Podman ab 4.0 funktioniert als Ersatz für Docker, wenn `docker` auf `podman`
zeigt und `DOCKER_HOST` auf den Podman-Socket gesetzt ist. `u-boot doctor`
meldet Podman-Versionen als Warnung („unrecognized version"), bricht aber nicht
ab.

---

## 2. Installation

### Binary (empfohlen)

Ein einzelnes, statisch gelinktes Programm ohne weitere Abhängigkeiten.

**Linux und macOS:**

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
curl -sSL -o u-boot \
  "https://github.com/pt9912/u-boot/releases/latest/download/u-boot-${OS}-${ARCH}"
chmod +x u-boot && sudo mv u-boot /usr/local/bin/
u-boot --version
```

**Windows (PowerShell):**

```powershell
Invoke-WebRequest `
  -Uri https://github.com/pt9912/u-boot/releases/latest/download/u-boot-windows-amd64.exe `
  -OutFile u-boot.exe
.\u-boot.exe --version
```

**Ergebnis:** `u-boot --version` gibt die installierte Version aus, zum
Beispiel `u-boot version 0.4.0`.

### Container-Image

```bash
docker run --rm -v "$PWD:/work" -w /work ghcr.io/pt9912/u-boot:latest init mein-projekt
```

Zwei Einschränkungen im Container:

- **`git` fehlt im Image.** `u-boot init` bricht dort mit einer git-Meldung ab.
  Nutzen Sie `--no-git` und legen Sie das Repository anschließend selbst an:

  ```bash
  docker run --rm -v "$PWD:/work" -w /work ghcr.io/pt9912/u-boot:latest \
    init mein-projekt --no-git
  git init
  ```

- **`u-boot doctor` sieht Ihre Werkzeuge nicht.** Die vier Host-Prüfungen
  (Docker, Compose-Plugin, Docker-Erreichbarkeit, git) werden übersprungen und
  mit `?` markiert.

Für den Alltag ist das Binary die bequemere Form.

---

## 3. Erste Schritte

### Ein Projekt in zwei Minuten

#### Voraussetzung

`u-boot` ist installiert. Sie stehen in einem leeren Verzeichnis.

#### Vorgehen

1. Prüfen Sie zuerst Ihre Umgebung:

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

4. Starten Sie die Umgebung:

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

Nach Schritt 4 läuft Ihr Stack; `u-boot up` wartet, bis jeder Dienst stabil
ist, und zeigt eine Statustabelle.

#### Hinweise

- Ohne Namen (`u-boot init`) leitet `u-boot` den Projektnamen aus dem
  Verzeichnisnamen ab.
- `u-boot init` legt auch ein Git-Repository an. Mit `--no-git` unterbleibt das.

### Grundlegende Bedienkonzepte

Drei Konzepte erklären das meiste Verhalten:

**Die Projektdatei `u-boot.yaml`** ist die Quelle der Wahrheit. Sie hält
Projektnamen und aktivierte Dienste. Alles andere — `compose.yaml`,
`.env.example`, Devcontainer — wird daraus abgeleitet.

**Verwaltete Bereiche.** In Dateien, die Ihnen gehören, ändert `u-boot` nur den
Abschnitt zwischen den Markierungen `BEGIN U-BOOT MANAGED BLOCK` und
`END U-BOOT MANAGED BLOCK`. Ihre eigenen Ergänzungen außerhalb dieser
Markierungen bleiben unangetastet.

**Nie stilles Überschreiben.** Trifft `u-boot` auf eine bestehende Datei, die
es nicht über einen verwalteten Bereich anfassen kann, bricht es ab und
erklärt, was fehlt — `--force` (Bereich ersetzen) oder `--backup` (vorher
sichern).

### Vorschau statt Überraschung

Jeder schreibende Befehl kennt zwei Vorschau-Flags:

| Flag | Wirkung |
|---|---|
| `--dry-run` | zeigt die geplanten Änderungen, schreibt nichts |
| `--diff` | zeigt die geplanten Änderungen als Unified Diff |

```bash
u-boot add postgres --dry-run
```

Nutzen Sie das, wann immer Sie unsicher sind. Es kostet nichts und ist
folgenlos.

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

Es entstehen neun Einträge (drei Verzeichnisse, sechs Dateien) und die
Projektdatei:

```yaml
schemaVersion: 1
project:
    name: mein-projekt
```

#### Hinweise

- Projektnamen dürfen Kleinbuchstaben, Ziffern und Bindestriche enthalten.
  Großbuchstaben und Unterstriche werden normalisiert; ist der Name danach
  ungültig, bricht `u-boot` mit einer Erklärung ab.
- Mit `--devcontainer` entstehen zusätzlich `.devcontainer/devcontainer.json`
  und ein zugehöriges `Dockerfile`.
- Mit `--template <name>` rendern Sie aus einer Vorlage statt aus dem
  Standardablauf. `u-boot template list` zeigt den Katalog.

### 4.2 Ein bestehendes Projekt erneut initialisieren

#### Voraussetzung

Im Verzeichnis liegen bereits Dateien, möglicherweise aus einem früheren
`u-boot init`.

#### Vorgehen

1. Sehen Sie sich zuerst an, was passieren würde:

   ```bash
   u-boot init --dry-run
   ```

2. Entscheiden Sie sich für einen der beiden Wege:

   ```bash
   u-boot init --force     # verwaltete Bereiche ersetzen
   u-boot init --backup    # bestehende Dateien vorher sichern
   ```

#### Ergebnis

Betroffene Pfade werden vor dem Schreiben aufgelistet. Mit `--backup` entstehen
Sicherungen nach dem Muster `<name>.bak`, `<name>.bak.1`, `<name>.bak.2` und so
fort — eine vorhandene Sicherung wird nie überschrieben.

#### Hinweise

- Erkennt `u-boot` Spuren eines bestehenden Projekts, ohne sicher zu sein,
  fragt es nach. In Skripten und CI beantworten Sie das mit
  `--assume-existing` — `--yes` genügt für diesen Fall bewusst **nicht**.
- Ohne `--force` und ohne `--backup` bricht der Lauf ab, statt zu raten.

### 4.3 Einen Dienst hinzufügen

#### Voraussetzung

Ein `u-boot`-Projekt (es gibt eine `u-boot.yaml`).

#### Vorgehen

1. Fügen Sie den Dienst hinzu:

   ```bash
   u-boot add postgres
   ```

   Verfügbar sind `postgres`, `keycloak` und `otel`.

2. Prüfen Sie die Änderung:

   ```bash
   cat u-boot.yaml
   ```

#### Ergebnis

```text
Added service "keycloak".

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
  keycloak:
    enabled: true
```

#### Hinweise

- Der Befehl ist wiederholbar: Ein zweiter `u-boot add postgres` ändert nichts
  und meldet keinen Fehler.
- Braucht ein Dienst einen anderen, installiert `--with-deps` das Fehlende
  automatisch mit, ohne Rückfrage.
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

#### Hinweise

- `--purge` fordert zusätzlich das Entfernen der Datenvolumes an. Das ist
  **destruktiv** und löst eine Sicherheitsabfrage aus; in nicht-interaktiven
  Läufen bricht der Befehl ohne `--yes` ab.
- In Version 0.4.0 entfernt `--purge` die Volumes **nicht** selbst. Die
  Zusammenfassung nennt Ihnen die passenden `docker volume rm`-Aufrufe für die
  manuelle Bereinigung.

### 4.5 Die Umgebung starten und stoppen

#### Voraussetzung

Docker läuft, und im Projekt liegt eine `compose.yaml`.

#### Vorgehen

```bash
u-boot up                # startet und wartet auf Stabilisierung
u-boot up --timeout 120  # wartet länger
u-boot up --timeout 0    # startet und kehrt sofort zurück
u-boot down              # stoppt die Umgebung
u-boot down --volumes    # stoppt und entfernt die Datenvolumes
```

#### Ergebnis

`u-boot up` gibt eine Statustabelle je Dienst aus (Name, Container-Status,
Port, Healthcheck). Ohne Zeitlimit (`--timeout 0`) entfällt die Wartezeit; Sie
bekommen stattdessen einen Hinweis, den Status später selbst zu prüfen.

#### Hinweise

- `--volumes` löscht Daten. Die Sicherheitsabfrage lässt sich in Skripten nur
  mit `--yes` überspringen.
- Der Standardwert für `--timeout` ist 60 Sekunden.

### 4.6 Logs ansehen

#### Vorgehen

```bash
u-boot logs                  # alle Dienste
u-boot logs postgres         # nur ein Dienst
u-boot logs --tail 100       # nur die letzten 100 Zeilen je Dienst
u-boot logs --follow         # laufend mitlesen, Abbruch mit Strg-C
```

#### Ergebnis

Die Log-Ausgabe der Compose-Dienste erscheint auf der Konsole. `--follow`
läuft, bis Sie mit Strg-C abbrechen; das gilt als normales Ende.

#### Hinweise

- `--follow` und `--json` schließen sich aus. Für maschinenlesbare Ausschnitte
  nutzen Sie `--tail <n> --json`.

### 4.7 Die Umgebung prüfen

#### Vorgehen

```bash
u-boot doctor
u-boot doctor --strict   # jede Warnung gilt als Fehler
```

#### Ergebnis

Ein Bericht über 13 Prüfungen mit vier Zuständen: `✓` in Ordnung, `⚠` Warnung,
`✗` Fehler, `?` übersprungen.

```text
Diagnostic report for /projekt
──────────────────────────────────────
⚠  uboot.yaml.valid    u-boot.yaml not present — directory is not a u-boot project.
   → Run `u-boot init` to create one (LH-FA-INIT-001).

Summary: 1 error, 2 warn, 6 ok
```

Jede Meldung enthält eine Zeile mit `→`, die den nächsten Schritt nennt.

#### Hinweise

- `--strict` ist für CI gedacht: Dort soll eine Warnung den Lauf anhalten.
- Läuft `u-boot` im Container, werden die vier Host-Prüfungen übersprungen
  (`?`). Das ist kein Fehler.

### 4.8 Artefakte erzeugen und aktualisieren

#### Vorgehen

```bash
u-boot generate readme
u-boot generate changelog
u-boot generate env-example
u-boot generate devcontainer
```

#### Ergebnis

Das jeweilige Artefakt entsteht oder wird aktualisiert. Bei bestehenden Dateien
ändert `u-boot` nur den verwalteten Bereich.

#### Hinweise

- Wiederholte Läufe mit unverändertem Projekt ändern nichts.
- Haben Sie innerhalb eines verwalteten Bereichs von Hand editiert, meldet
  `u-boot` den Konflikt, statt Ihre Arbeit zu überschreiben.
- Devcontainer-Features und der zugehörige Drift-Check sind in
  [`devcontainer-features.md`](devcontainer-features.md) beschrieben.

### 4.9 Vorlagen nutzen

#### Vorgehen

```bash
u-boot template list                 # Katalog ansehen
u-boot init mein-projekt --template basic
```

#### Ergebnis

`template list` zeigt eine Tabelle mit `NAME`, `DESCRIPTION` und `VERSION`:

```text
NAME   DESCRIPTION                                        VERSION
basic  Minimal u-boot project skeleton — same files …     0.1.0
```

`init --template` rendert das Projekt aus der gewählten Vorlage.

#### Hinweise

- `--template` gilt nur für frische Projekte und lässt sich nicht mit
  `--devcontainer`, `--force`, `--backup`, `--dry-run` oder `--diff`
  kombinieren.

### 4.10 In Skripten und CI verwenden

#### Voraussetzung

Der Lauf darf nicht auf eine Eingabe warten.

#### Vorgehen

Wählen Sie **einen** der beiden Modi — sie schließen sich aus:

```bash
u-boot add postgres --yes              # Rückfragen automatisch bejahen
u-boot add postgres --no-interactive   # bei jeder Rückfrage abbrechen
```

Maschinenlesbare Ausgabe erhalten Sie mit `--json`:

```bash
u-boot doctor --json
```

#### Ergebnis

`--json` liefert einen einheitlichen Umschlag mit Status, Daten und
Diagnosemeldungen. Das vollständige Schema samt Exit-Code-Matrix je
Subkommando steht in [`cli-json-output.md`](cli-json-output.md).

#### Hinweise

- Werten Sie in Skripten den **Exit-Code** aus, nicht den Ausgabetext (siehe
  [Abschnitt 6](#6-fehlerbehebung)).
- `--yes` und `--no-interactive` gleichzeitig ist ein Nutzungsfehler und
  endet mit Exit-Code 2.

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
  enabled: false
```

`schemaVersion` gehört `u-boot` und sollte nicht von Hand geändert werden.

### Werte lesen und ändern

```bash
u-boot config                            # ganze Datei anzeigen
u-boot config get project.name           # einen Wert lesen
u-boot config set project.name neuer-name
```

Schreibbar sind genau zwei Pfade:

| Pfad | Bedeutung | Schreibbar |
|---|---|---|
| `project.name` | Projektname | ja |
| `devcontainer.enabled` | Devcontainer-Unterstützung | ja |
| `services.<dienst>.enabled` | Dienst aktiv | **nein**, nur lesbar |

Dienste schalten Sie über `u-boot add` und `u-boot remove`, nicht über
`config set`. Der Grund: Ein Dienst besteht nicht nur aus einem Schalter,
sondern auch aus Einträgen in `compose.yaml` und `.env.example` — ein direkt
gesetzter Schalter würde diese auseinanderlaufen lassen.

Jeder `config set`-Aufruf wird gegen das Schema geprüft, **bevor** geschrieben
wird. Ungültige Werte führen zu Exit-Code 10, ohne die Datei anzufassen.

### Ausführlichkeit der Ausgabe

| Flag | Wirkung |
|---|---|
| `--quiet` | nur das Nötigste; Warnungen und Fehler bleiben |
| `--verbose` | zusätzliche Details |
| `--debug` | interne Diagnoseausgabe |

---

## 6. Fehlerbehebung

### Exit-Codes

Werten Sie in Skripten diese Codes aus:

| Code | Bedeutung |
|---|---|
| `0` | Erfolg |
| `1` | allgemeiner Fehler |
| `2` | falsche Benutzung (unbekanntes Flag, fehlendes Argument, widersprüchliche Modus-Flags) |
| `10` | fachlicher Validierungsfehler (ungültiger Name, fehlendes Projekt, verweigerte Bestätigung) |
| `11` | Umgebungsproblem (Docker nicht erreichbar, Compose-Plugin fehlt) |
| `12` | Laufzeitfehler beim Ausführen von Compose |
| `14` | Dateisystem- oder Persistenzfehler |

### Fehler: „u-boot.yaml not present — directory is not a u-boot project"

#### Ursache

Sie stehen nicht in einem `u-boot`-Projekt, oder die Projektdatei fehlt.

#### Lösung

1. Prüfen Sie Ihr Verzeichnis: `pwd` und `ls u-boot.yaml`.
2. Wechseln Sie in das Projektverzeichnis — oder legen Sie mit
   `u-boot init` ein neues Projekt an.

### Fehler: Docker nicht erreichbar (Exit-Code 11)

#### Ursache

Der Docker-Daemon läuft nicht, oder Ihr Benutzer darf nicht auf ihn zugreifen.

#### Lösung

1. Führen Sie `u-boot doctor` aus und lesen Sie die `→`-Zeile der
   fehlgeschlagenen Prüfung.
2. Prüfen Sie den Daemon: `docker version`.
3. Prüfen Sie das Compose-Plugin: `docker compose version`.
4. Unter Linux: Ist Ihr Benutzer in der Gruppe `docker`?

### Fehler: Abbruch statt Überschreiben

#### Ursache

Es gibt bereits eine Datei, die `u-boot` nicht über einen verwalteten Bereich
ändern kann.

#### Lösung

1. Sehen Sie sich die geplanten Änderungen an: `u-boot init --diff`.
2. Entscheiden Sie sich:
   - `--force` ersetzt den verwalteten Bereich,
   - `--backup` sichert die Datei vorher.
3. Führen Sie den Befehl mit dem gewählten Flag erneut aus.

### Fehler: Bestätigung im nicht-interaktiven Lauf (Exit-Code 10 oder 2)

#### Ursache

Der Befehl braucht eine Bestätigung, kann sie aber nicht einholen.

#### Lösung

- Für destruktive Operationen (`down --volumes`, `remove --purge`): `--yes`
  ergänzen — bewusst, denn dabei gehen Daten verloren.
- Für die Erkennung eines bestehenden Projekts bei `init`:
  `--assume-existing` ergänzen. `--yes` genügt hier **nicht**.

### Fehler: „--tail with invalid value" (Exit-Code 2)

#### Ursache

`--tail` hat einen negativen oder nicht-numerischen Wert bekommen.

#### Lösung

Geben Sie eine positive ganze Zahl an: `u-boot logs --tail 200`.

### Wenn nichts davon hilft

1. Wiederholen Sie den Befehl mit `--debug`.
2. Führen Sie `u-boot doctor` aus und halten Sie die Ausgabe bereit.
3. Notieren Sie Version (`u-boot --version`), Betriebssystem und den genauen
   Befehl.
4. Eröffnen Sie ein Issue (siehe [Abschnitt 9](#9-support-und-lizenz)).

---

## 7. FAQ

**Brauche ich Docker, um `u-boot` zu benutzen?**
Nur für `up`, `down`, `logs` und die Docker-Prüfungen von `doctor`. `init`,
`generate`, `config` und `template list` laufen ohne.

**Überschreibt `u-boot` meine Änderungen?**
Nein. Außerhalb der verwalteten Bereiche wird nichts angefasst, und innerhalb
nur mit `--force` oder nach einer Sicherung mit `--backup`.

**Kann ich `u-boot` auf ein bestehendes Projekt anwenden?**
Ja. Nutzen Sie `u-boot init --dry-run`, um vorher zu sehen, was passieren
würde.

**Warum kann ich `services.<dienst>.enabled` nicht mit `config set` ändern?**
Weil ein Dienst mehr ist als ein Schalter. `u-boot add` und `u-boot remove`
pflegen zusätzlich `compose.yaml` und `.env.example`.

**Was passiert, wenn ich einen Dienst zweimal hinzufüge?**
Nichts. Der zweite Aufruf ist wirkungslos und meldet keinen Fehler.

**Funktioniert Podman?**
Ja, ab Version 4.0, wenn `docker` auf `podman` zeigt und `DOCKER_HOST` gesetzt
ist. `doctor` meldet die Podman-Version als Warnung, blockiert aber nicht.

**Warum überspringt `doctor` im Container vier Prüfungen?**
Aus dem Container heraus sind die Werkzeuge Ihres Rechners nicht sichtbar. Ein
Fehlschlag wäre irreführend, deshalb der Zustand `?`.

**Wo finde ich die genaue Ausgabe eines Befehls?**
`u-boot <befehl> --help` ist die verbindliche Referenz. Beispielabläufe stehen
in [`examples.md`](examples.md).

---

## 8. Glossar

| Begriff | Bedeutung |
|---|---|
| Add-on | ein zuschaltbarer Dienst wie `postgres`, `keycloak` oder `otel` |
| Compose | Docker Compose; auch die Datei `compose.yaml` |
| Devcontainer | containerisierte Entwicklungsumgebung, meist mit VS Code genutzt |
| Exit-Code | Rückgabewert eines Befehls; `0` bedeutet Erfolg |
| Healthcheck | Prüfung, ob ein Dienst technisch bereit ist |
| Idempotenz | mehrfaches Ausführen führt zum selben Ergebnis |
| Verwalteter Bereich | Abschnitt zwischen `BEGIN`/`END U-BOOT MANAGED BLOCK`, den `u-boot` pflegt |
| Projektdatei | `u-boot.yaml` — Quelle der Wahrheit für Name und Dienste |
| Stabilisierung | Wartezeit nach `up`, bis jeder Dienst als bereit gilt |

---

## 9. Support und Lizenz

**Fragen und Fehler:** <https://github.com/pt9912/u-boot/issues>

Legen Sie einem Fehlerbericht bei: die Ausgabe von `u-boot --version`, die
Ausgabe von `u-boot doctor`, Ihr Betriebssystem und den genauen Befehl.

**Weiterführende Dokumentation:**

- [`examples.md`](examples.md) — Beispielabläufe als Kommando-Rezepte
- [`cli-json-output.md`](cli-json-output.md) — JSON-Schema und Exit-Code-Matrix
- [`devcontainer-features.md`](devcontainer-features.md) — Devcontainer-Features
- [`quality.md`](quality.md) — Qualitäts-Gates (für Mitwirkende)

**Lizenz:** siehe `LICENSE` im Projektarchiv.

---

## 10. Änderungshistorie

| Handbuch-Version | Software-Version | Datum | Änderung |
|---|---|---|---|
| 1.0 | v0.4.0 | 2026-07-25 | Erstfassung |

**Gültigkeitsbereich:** Dieses Handbuch beschreibt `u-boot` v0.4.0. Prüfen Sie
mit `u-boot --version`, welche Version Sie einsetzen. Weicht sie ab, ist
`u-boot <befehl> --help` die verbindliche Auskunft.

**Noch nicht enthalten:** Mit der nächsten Version löst `--template` auch
Dateisystempfade auf (`u-boot init --template ./meine-vorlage`), nicht nur
Katalognamen. Bis dahin gilt die Beschreibung in
[Abschnitt 4.9](#49-vorlagen-nutzen).
