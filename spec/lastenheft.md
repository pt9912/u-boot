# Lastenheft: `u-boot` – Projekt-Bootstrapper für Docker-/Devcontainer-Stacks

| Dokument         | Lastenheft                                                         |
| ---------------- | ------------------------------------------------------------------ |
| Projektname      | `u-boot`                                                           |
| Kurzbeschreibung | CLI-Tool zum Bootstrapping reproduzierbarer Entwicklungsumgebungen |
| Zielplattform    | Linux, Docker, VS Code Dev Containers                              |
| Hauptnutzer      | Softwareentwickler, DevOps-Engineers, technische Teams             |
| Version          | 0.1.0                                                              |
| Status           | Entwurf                                                            |
| Datum            | 2026-05-21                                                         |

---

## 0. Lesehinweise

### LH-LESE-001 – Modalverben

In diesem Dokument haben Modalverben folgende Bedeutung (in Anlehnung an RFC 2119):

- **muss** – verbindliche Anforderung (Pflicht).
- **soll** – Empfehlung; Abweichungen müssen begründet werden.
- **kann** / **darf** – optionale Eigenschaft oder ausdrückliche Erlaubnis.

---

### LH-LESE-002 – Sprache

Die Spezifikation ist auf Deutsch verfasst.

CLI-Ausgaben, Fehlermeldungen und erzeugte Dateien (Kommentare, Beispielwerte, README-Vorlagen) sind auf Englisch.

---

## 1. Zielbestimmung

### LH-ZB-001 – Projektziel

`u-boot` soll ein CLI-Tool werden, das vollständige Entwicklungsumgebungen für Docker-basierte Softwareprojekte erzeugt, erweitert, prüft und startet.

Das Tool soll insbesondere Projektstrukturen, Docker-Konfigurationen, Devcontainer-Setups, optionale Infrastrukturservices und wiederkehrende Entwicklungsartefakte automatisch bereitstellen.

---

### LH-ZB-002 – Produktvision

`u-boot` soll sich wie ein **Bootloader für Entwicklungsumgebungen** verhalten.

Ein neues oder bestehendes Projekt soll mit wenigen Befehlen in einen lauffähigen, reproduzierbaren Entwicklungszustand versetzt werden können.

Beispiel:

```bash
u-boot init
u-boot add postgres
# optional in späteren Versionen:
# u-boot add keycloak
# u-boot add otel
u-boot up
```

---

### LH-ZB-003 – Repo-Beschreibung

```text
u-boot: A developer environment bootloader for Docker-based projects.
```

---

## 2. Produkteinsatz

### LH-PE-001 – Anwendungsbereich

`u-boot` soll für Softwareprojekte eingesetzt werden, die lokal oder in Devcontainern entwickelt werden und Docker beziehungsweise Docker Compose als zentrale Laufzeitumgebung verwenden.

---

### LH-PE-002 – Zielgruppen

Das Produkt richtet sich an:

- Softwareentwickler
- DevOps-Engineers
- technische Projektleiter
- Entwicklerteams mit Docker-basierten Entwicklungsumgebungen
- Teams, die reproduzierbare lokale Setups benötigen

---

### LH-PE-003 – Betriebsumgebung

Die primäre Betriebsumgebung ist:

- Linux
- Docker Engine oder kompatible Docker-Laufzeit
- Docker Compose
- Git
- optional: VS Code mit Dev Containers Extension

Sekundäre Betriebsumgebungen können später ergänzt werden:

- macOS
- Windows mit WSL2

---

## 3. Produktübersicht

### LH-PÜ-001 – Grundfunktion

`u-boot` soll als Kommandozeilenwerkzeug bereitgestellt werden.

Das Tool soll über Befehle wie die folgenden bedient werden:

```bash
u-boot init
u-boot up
u-boot doctor
u-boot add postgres
# optional in späteren Versionen:
# u-boot add keycloak
# u-boot add otel
u-boot generate changelog
```

---

### LH-PÜ-002 – Hauptmodule

Das Produkt soll mindestens folgende fachliche Module besitzen:

| Kennung    | Modul                  | Beschreibung                                                  |
| ---------- | ---------------------- | ------------------------------------------------------------- |
| LH-MOD-001 | Projektinitialisierung | Erzeugt neue Projektstruktur                                  |
| LH-MOD-002 | Devcontainer-Generator | Erzeugt `.devcontainer`-Konfiguration                         |
| LH-MOD-003 | Docker-Stack-Generator | Erzeugt Dockerfile und Compose-Dateien                        |
| LH-MOD-004 | Service-Add-ons        | Fügt Dienste wie PostgreSQL, Keycloak und OpenTelemetry hinzu |
| LH-MOD-005 | Umgebungsprüfung       | Prüft lokale Voraussetzungen                                  |
| LH-MOD-006 | Stack-Start            | Startet Entwicklungsumgebung                                  |
| LH-MOD-007 | Generatoren            | Erzeugt Zusatzdateien wie Changelog, README oder Configs      |
| LH-MOD-008 | Template-System        | Verwaltet wiederverwendbare Projektvorlagen                   |

---

## 4. Funktionale Anforderungen

## 4.1 CLI-Grundverhalten

### LH-FA-CLI-001 – CLI-Aufruf

Priorität: MVP

Das Produkt muss als Kommandozeilenprogramm mit dem Namen `u-boot` aufrufbar sein.

Beispiel:

```bash
u-boot --help
```

---

### LH-FA-CLI-002 – Hilfeausgabe

Priorität: MVP

Das Produkt muss eine Hilfeausgabe bereitstellen.

Die Hilfeausgabe muss mindestens enthalten:

- verfügbare Befehle
- Kurzbeschreibung je Befehl
- Optionen je Befehl
- Beispiele

---

### LH-FA-CLI-003 – Versionsausgabe

Priorität: MVP

Das Produkt muss die installierte Version ausgeben können.

Beispiel:

```bash
u-boot --version
```

---

### LH-FA-CLI-004 – Fehlerausgabe

Priorität: MVP

Das Produkt muss verständliche Fehlermeldungen ausgeben.

Fehlermeldungen müssen enthalten:

- Ursache
- betroffener Befehl oder betroffene Datei
- empfohlene Korrekturmaßnahme

---

### LH-FA-CLI-005 – Verbosity und Logging

Priorität: MVP

Das Produkt muss ein konfigurierbares Ausgabeverbose-Level (Verbosity) unterstützen.

Mindestens müssen folgende Stufen unterstützt werden:

- `--quiet` – nur Fehler
- Standard – Statusmeldungen und Fehler
- `--verbose` – zusätzlich Detailinformationen
- `--debug` – zusätzlich interne Diagnoseausgaben

Werden mehrere Verbosity-Optionen gleichzeitig angegeben (z. B. `--quiet --verbose`), gewinnt die zuletzt auf der Kommandozeile angegebene Option. Eine Validierungsabweisung wegen Mehrfachangabe erfolgt nicht.

---

### LH-FA-CLI-005A – Interaktivität und Automatisierung

Priorität: MVP

Das Produkt muss nicht-interaktive Ausführung unterstützen.

Es muss mindestens folgende Optionen bieten:

- `--yes` – Standardfragen automatisch bejahen (für CI/Skripte).
- `--no-interactive` – keine Rückfragen stellen; erforderliche Bestätigungen mit klarem Fehler abbrechen.
- `--yes` und `--no-interactive` sind exklusiv. Bei gleichzeitiger Nutzung ist ein CLI-Fehler mit Exit-Code `2` (`LH-FA-CLI-006`) zu erzeugen.
- Für deterministisches Verhalten in Skripten und CI sind beide Modi einzeln nutzbar.
- Die Optionen sind auf Befehle anzuwenden, die Bestätigungsentscheidungen benötigen (insb. `u-boot init`, `u-boot add`, `u-boot remove`, `u-boot config set`, `u-boot down --volumes`).
- Für `u-boot init` ist zusätzlich das Flag `--assume-existing` definiert (nicht global, nur für diesen Befehl):
  - Ohne `--assume-existing` wird eine implizite Erkennung als bestehendes Projekt im nicht-interaktiven Modus nicht automatisch akzeptiert.
  - Mit `--assume-existing` wird die implizite Erkennung als bestehendes Projekt in nicht-interaktiven Läufen akzeptiert.
  - `--yes` ist für diesen Sonderfall **nicht** ausreichend; die implizite Erkennung bleibt abgelehnt, wenn keine `--assume-existing` gesetzt ist.
  - Ohne `--assume-existing` und bei nicht-interaktivem Lauf ist die implizite Erkennung zwingend ablehnend und erzeugt einen fachlichen Fehler.
  - Der Fehlercode für diese Abweisung ist `10`.
- Bei aktivierter Nicht-Interaktivität darf keine neue Rückfrage erzeugt werden:
  - mit `--no-interactive` bricht der Aufruf bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab,
  - mit `--yes` wird die vorgesehene Standardentscheidung deterministisch ausgeführt.
- Für bereits deterministische Ausführungspfade (keine relevante Rückfrage) ist das Verhalten in beiden Modi unverändert.

Bei `u-boot init` gilt zusätzlich die feste Auswertungsreihenfolge im nicht-interaktiven Modus:

- ohne `--assume-existing`: keine implizite Annahme einer bestehenden Projekterkennung, deterministisch abbrechen (Exit-Code `10` bei bestehendem Projekt),
- mit `--assume-existing`: implizite Annahme als bestehendes Projekt (soweit kompatibel mit den übrigen Validierungen).

Bei destruktiven Operationen (insb. `u-boot down --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad (`--yes` oder aktiv bestätigten interaktiven Pfad) erfolgen. Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.

Deterministische Auswertungslogik für bestätigungsrelevante Modi:

- `--yes` und `--no-interactive` sind exklusiv.
- `--no-interactive` erlaubt keinerlei Rückfragen. Alle Entscheidungswege müssen deterministisch sein oder mit `LH-FA-CLI-006`-Code `2` abbrechen, wenn eine notwendige Bestätigung fehlt.
- `--yes` erlaubt deterministische Standardpfade ohne Nutzerinteraktion.
- `--force` und/oder `--backup` sind in nicht-interaktiven Läufen explizit zulässig, weil beide Modi deterministisch arbeiten.
- `--no-interactive` + `--force` erlaubt das Überschreiben ohne Rückfrage; dabei ist immer eine vollständige Zusammenfassung der betroffenen Pfade auszugeben.
- `--force` darf keine zusätzlichen Rückfragen erzeugen; die Sicherheitslogik beschränkt sich auf die Validierung der Eingabedaten.
- `--backup` ist optional. Wenn `--backup` gesetzt ist, dürfen Dateischutz-Szenarien mit automatischer Sicherung deterministisch abgearbeitet werden.
- Bei fehlender Möglichkeit zur sicheren automatischen Abarbeitung (z. B. fehlender verwalteter Block ohne `--backup` bei vollständig kontrolliertem Überschreiben) muss der Befehl mit Fehlercode `10` abbrechen.

---

### LH-FA-CLI-006 – Exit Codes

Priorität: MVP

Das Produkt muss aussagekräftige Exit Codes liefern.

Mindestens:

- `0` – Erfolg
- `1` – allgemeiner Fehler
- `2` – fehlerhafte CLI-Nutzung
- `10` – fachlicher Validierungsfehler (z. B. ungültige Konfiguration)
- `11` – fachlicher Umgebungs-/Prüfungsfehler (z. B. Docker nicht erreichbar)
- `12` – fachlicher Ausführungsfehler (z. B. Compose-Startfehler)
- `3` bis `9` – reserviert (nicht verwenden)
- `13` – technischer Infrastruktur-/Feature-Source-Fehler (z. B. Laden externer Devcontainer-Features fehlgeschlagen)
- `14` – technischer Persistenz- oder Dateisystemfehler (z. B. unerwartete IO-/Permissions-Probleme)
- `15` – technischer Ausführungsfehler außerhalb der fachlichen Domäne
- `16` bis `19` – reserviert (nicht verwenden)

Empfehlung für typische Fehlerzuordnung:

- `10` bei Struktur-, Namens- oder Konfigurationsvalidierungsfehlern
- `11` bei Umgebungsproblemen (z. B. fehlende Tools, Versionsinkompatibilität)
- `12` bei Laufzeitfehlern beim Ausführen von Docker/Compose-Operationen

Für alle fachlichen Fehler ist die Verwendung von `10`, `11` oder `12` bindend.
Nicht-fachliche Fehler dürfen standardmäßig mit `1` codiert werden, wenn eine feinere technische Klassifikation nicht sinnvoll ist.
`13` bis `15` dürfen zusätzlich verwendet werden, wenn deren Bedeutung für den aufrufenden Kontext explizit dokumentiert ist.
`16` bis `19` sind in der aktuellen Spezifikation nicht zu verwenden.

---

### LH-FA-CLI-007 – Dry Run

Priorität: V1

Das Produkt muss für dateiverändernde Befehle einen Dry-Run-Modus unterstützen.

Beispiel:

```bash
u-boot add postgres --dry-run
```

Der Dry-Run muss anzeigen, welche Dateien erzeugt, geändert oder gelöscht würden, ohne Änderungen am Dateisystem vorzunehmen.

`--dry-run` darf mit `--diff` kombiniert werden. Bei gleichzeitiger Nutzung darf das Tool keine Datei schreiben; die Ausgabe besteht ausschließlich aus dem geplanten Änderungsplan und der Diff-Darstellung.

Bei gleichzeitiger Verwendung von `--dry-run` und `--json` muss die Ausgabe streng maschinenlesbar (`JSON`) erfolgen und keine unstrukturierten Text-UI-Zeilen enthalten.

Für `--dry-run --json` ist die Ausgabe mindestens wie folgt als maschinenlesbares JSON zu liefern:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"],
  "properties": {
    "subcommand": {
      "type": "string",
      "description": "Unterkommando bei gruppierten Hauptkommandos wie `template` oder `config`"
    },
    "status": {
      "type": "string",
      "enum": ["ok", "warn", "error"]
    },
    "command": {
      "type": "string",
      "enum": ["init", "add", "remove", "up", "down", "doctor", "logs", "generate", "config", "template"]
    },
    "dryRun": {
      "type": "boolean"
    },
    "diff": {
      "type": "boolean"
    },
    "plannedFiles": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "action"],
        "properties": {
          "path": { "type": "string" },
          "action": {
            "type": "string",
            "enum": ["create", "modify", "delete"]
          }
        },
        "additionalProperties": true
      }
    },
    "changes": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "count"],
        "properties": {
          "path": { "type": "string" },
          "count": { "type": "integer", "minimum": 0 }
        },
        "additionalProperties": true
      }
    },
    "diagnostics": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["level", "code", "message"],
        "properties": {
          "level": { "type": "string", "enum": ["warn", "error"] },
          "code": { "type": "string" },
          "message": { "type": "string" },
          "file": { "type": "string" }
        },
        "additionalProperties": true
      }
    },
    "exitCode": {
      "type": "integer",
      "minimum": 0
    }
  },
  "allOf": [
    {
      "if": {
        "properties": {
          "command": { "const": "template" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    },
    {
      "if": {
        "properties": {
          "command": { "const": "config" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    }
  ],
  "additionalProperties": true
}
```

Bei gruppierten Befehlen wie `command == "template"` oder `command == "config"` muss das Feld `subcommand` gesetzt sein (z. B. `list`, `get`, `set`).

Beispielinstanz:

```json
{
  "status": "warn",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "create" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 12 }
  ],
  "diagnostics": [
    { "level": "warn", "code": "LH-FA-CLI-007", "message": "Geplante Datei fehlt bereits" }
  ],
  "exitCode": 0
}
```

Weitere Felder sind erlaubt.

Konvention für `diagnostics[*].code`: LH-Kennung der verursachenden Anforderung (z. B. `LH-FA-DEV-003`, `LH-FA-CLI-007`). Tool-interne Codes ohne LH-Bezug dürfen nur dann verwendet werden, wenn ihre Bedeutung in der Dokumentation festgehalten ist.

Das Feld `status` ist an den höchsten in `diagnostics` enthaltenen `level` gekoppelt: enthält `diagnostics` mindestens einen `error`-Eintrag, ist `status == "error"`; enthält es mindestens einen `warn`-Eintrag (und keinen `error`), ist `status == "warn"`; andernfalls `status == "ok"`. Diese Regel gilt für alle `--json`-Ausgaben (`LH-FA-CLI-007`, `LH-FA-CLI-008`, `LH-NFA-USE-004`).

---

### LH-FA-CLI-008 – Diff-Ausgabe

Priorität: V1

Das Produkt soll bei dateiverändernden Befehlen eine Diff-Ausgabe unterstützen.

Beispiel:

```bash
u-boot add postgres --diff
```

Die Diff-Ausgabe muss Unterschiede zwischen aktuellem und geplantem Zustand der betroffenen Dateien zeigen.

Wird `--diff` ohne `--dry-run` gesetzt, zeigt sie den geplanten Endzustand als Vorschau.
Wird `--diff` mit `--dry-run` kombiniert, gilt dieselbe Vorschau bei vollständigem Schreibschutz.

Bei Kombination von `--diff` mit `--json` ist die komplette JSON-Struktur inkl. Pflichtfeldern aus dem in `LH-FA-CLI-007` definierten Schema auszugeben. Die Felder `dryRun` und `diff` sind dabei korrekt auf den konkreten Ausführungsmodus gesetzt (`dryRun` je nach Aufruf, `diff` immer `true`).

Beispiel für `--diff --json` ohne `--dry-run` (Vorschau mit anschließendem Schreiben):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": false,
  "diff": true,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "modify" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 6 }
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

Für reine Vorschau-Workflows gelten die selben Exit-Codes wie bei der Nicht-Diff-Ausgabe.

---

## 4.2 Projektinitialisierung

### LH-FA-INIT-001 – Neues Projekt initialisieren

Priorität: MVP

Das Produkt muss mit folgendem Befehl ein neues Projekt initialisieren können:

```bash
u-boot init
u-boot init --assume-existing
```

---

### LH-FA-INIT-002 – Projektname

Priorität: MVP

Das Produkt muss bei der Initialisierung einen Projektnamen verwenden können.

Wird kein Name explizit angegeben, verwendet das Tool standardmäßig den aktuellen Verzeichnisnamen als Basis.

Der abgeleitete Name wird deterministisch normalisiert:

1. Der Basisname des Arbeitsverzeichnisses wird auf Kleinbuchstaben gesetzt.
2. Alle Zeichen außer `a-z`, `0-9` und `-` werden auf `-` abgebildet.
3. aufeinanderfolgende `-` werden zu einem einzelnen `-` zusammengeführt.
4. führende und nachgestellte `-` sowie Leerzeichen werden entfernt.
5. Die Länge wird auf 1 bis 63 Zeichen begrenzt.
6. Nach Kürzung auf 63 Zeichen wird erneut auf führende/nachgestellte `-` geprüft und diese notfalls entfernt.
7. Anschließend wird der Name gegen die Validierung in `LH-FA-INIT-006` geprüft.

Beispiel:

```bash
u-boot init
u-boot init my-service
```

Ist kein gültiger Name ableitbar oder angegeben, muss der Befehl mit einer klaren Fehlermeldung abbrechen und auf die explizite Übergabe eines Namens (`u-boot init <name>`) verweisen.

---

### LH-FA-INIT-003 – Projektstruktur erzeugen

Priorität: MVP

Das Produkt muss eine grundlegende Projektstruktur erzeugen.

Mindestumfang:

```text
.
├── docker/
├── scripts/
├── docs/
├── README.md
├── CHANGELOG.md
├── compose.yaml
├── .env.example
├── u-boot.yaml
└── .gitignore
```

Bei aktivierter Devcontainer-Unterstützung (siehe `LH-FA-DEV-001`) zusätzlich:

```text
.devcontainer/devcontainer.json
.devcontainer/Dockerfile
```

---

### LH-FA-INIT-004 – Bestehendes Projekt erkennen

Priorität: MVP

Das Produkt muss erkennen, ob es in einem bestehenden Projektverzeichnis ausgeführt wird.

Relevante Dateien sind mindestens:

- `u-boot.yaml`
- `compose.yaml`
- `.env.example`
- `README.md`
- `CHANGELOG.md`
- `.gitignore`
- `docs/`
- `scripts/`
- `docker/`
- `.devcontainer/devcontainer.json`

Wenn mindestens eine der Projektsteuerdateien (`u-boot.yaml`, `compose.yaml`, `.env.example`) vorhanden ist, ist das Verzeichnis als bestehendes Projekt zu behandeln.
Liegt keine Projektsteuerdatei vor, gilt das Verzeichnis nur als wahrscheinliches bestehendes Projekt, wenn mindestens drei Elemente aus dem Mindestumfang der Projektstruktur bereits vorhanden sind.
In diesem Fall muss `u-boot init` im interaktiven Modus explizit nachfragen, ob das Verzeichnis als bestehendes Projekt behandelt werden soll.
`--assume-existing` ist nur für `u-boot init` gültig.
Das genaue Verhalten im nicht-interaktiven Modus (mit/ohne `--assume-existing`, Exit-Code-Vergabe) ist verbindlich in `LH-FA-CLI-005A` definiert; diese Anforderung wiederholt es nicht.
Bestehende Dateien dürfen auch bei impliziter oder expliziter Annahme als bestehendes Projekt nicht kommentarlos überschrieben werden; es gilt der Überschreibschutz aus `LH-FA-INIT-005`.

---

### LH-FA-INIT-005 – Überschreibschutz

Priorität: MVP

Das Produkt muss vor dem Überschreiben bestehender Dateien schützen.

Standardverhalten ohne Option:

- Abbruch mit Hinweis, welche Datei kollidiert

Zusätzliche Strategien über Option:

- `--backup` – bestehende Datei als `<name>.bak` sichern und ersetzen; ist `<name>.bak` bereits vorhanden, wird automatisch `<name>.bak.1`, `<name>.bak.2`, ... verwendet (kleinster freier numerischer Suffix), ohne vorhandene Backups zu überschreiben.
- Für bestehende Verzeichnisse (z. B. `docs/`, `scripts/`, `docker/`, `.devcontainer/`) wird der komplette Verzeichnisbaum rekursiv als `<name>.bak*` gesichert und innerhalb derselben Operation ersetzt; bei Fehlern während des Ersetzens muss ein Rollback auf den ursprünglichen Zustand durchgeführt werden (POSIX-Atomarität für rekursive Bäume wird nicht garantiert).
- `--force` – bestehende Dateien ohne Rückfrage überschreiben; vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden

Zusätzliche Schutzregeln für strukturierte Konfigurationsdateien (`compose.yaml`, `.env.example`, `README.md`, `CHANGELOG.md`, `.devcontainer/devcontainer.json`):

- bestehende, nicht verwaltete Inhalte bleiben in `--force`-Ausführung erhalten.
- wird ein `U-BOOT MANAGED BLOCK` erkannt, darf bei `--force` nur dieser Block verändert werden. Das Markierungsformat pro Dateityp ist in `LH-SA-FILE-002` definiert.
- für `.devcontainer/devcontainer.json` gilt der JSONC-Markerstil (`// BEGIN U-BOOT MANAGED BLOCK: <name>` / `// END U-BOOT MANAGED BLOCK: <name>`); für strikte JSON-Dateien ohne Kommentar-Support wird die gesamte Datei als verwaltet behandelt und in `u-boot.yaml` referenziert.
- fehlt ein verwalteter Block in einer vorhandenen Datei:
  - ist `--backup` gesetzt, wird vor jedem vollständigen Überschreiben der komplette Dateiinhalt gesichert und danach ersetzt.
  - ist `--backup` nicht gesetzt, wird der Vorgang mit einem fachlichen Fehler (Code `10`) abgebrochen; es erfolgt ein klarer Hinweis auf die nötige Option `--backup`.
- bei vollständiger Überschreibung ohne verwalteten Block gilt ein vollständiges Backup vor dem Schreiben als Pflicht.

Für `--force`, `--backup` und nicht-interaktive Modi gilt zusätzlich die in `LH-FA-CLI-005A` definierte Entscheidungslogik für Bestätigungen.

---

### LH-FA-INIT-006 – Projektnamen-Validierung

Priorität: MVP

Das Produkt muss den Projektnamen validieren.

Die Validierung gilt für den explizit übergebenen und den automatisch aus dem Arbeitsverzeichnis abgeleiteten Projektnamen.

Regeln:

- erlaubt sind Kleinbuchstaben, Ziffern und Bindestrich
- beginnt mit einem Kleinbuchstaben
- endet mit einem Kleinbuchstaben oder einer Ziffer (entfällt bei einstelligen Namen)
- minimale Länge: 1 Zeichen
- maximale Länge: 63 Zeichen
- regulärer Ausdruck: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`

Ungültige Namen müssen mit einer klaren Fehlermeldung abgelehnt werden.

---

### LH-FA-INIT-007 – Git-Repository-Initialisierung

Priorität: MVP

Das Produkt muss Git-Initialisierung als Teil von `u-boot init` unterstützen.

Verhalten:

- Standardverhalten: aktiviert – ein neues Git-Repository wird angelegt, sofern noch keines vorhanden ist.
- Abschaltbar über `--no-git`.
- Ein bereits vorhandenes Repository darf nicht erneut initialisiert werden.

---

## 4.3 Devcontainer-Unterstützung

### LH-FA-DEV-001 – Devcontainer erzeugen

Priorität: MVP

Das Produkt muss eine Devcontainer-Konfiguration erzeugen können.

Die Erzeugung muss sowohl bei `u-boot init` über eine Option als auch nachträglich auslösbar sein:

```bash
u-boot init --devcontainer
u-boot generate devcontainer
```

Mindestdateien:

```text
.devcontainer/devcontainer.json
.devcontainer/Dockerfile
```

---

### LH-FA-DEV-002 – VS-Code-Kompatibilität

Priorität: MVP

Die erzeugte Devcontainer-Konfiguration muss mit VS Code Dev Containers kompatibel sein.

---

### LH-FA-DEV-003 – Devcontainer-Features

Priorität: V1

Das Produkt soll optionale Devcontainer-Features unterstützen.

Beispiele:

- Git
- Docker CLI
- Node.js
- Java
- SDKMAN
- PostgreSQL Client
- Kubernetes Tools

Für optionale externe Feature-Quellen gilt:

- Keine fremden Skripte dürfen ohne Zustimmung ausgeführt werden (`LH-NFA-SEC-004`).
- Standardmäßig sind nur lokal hinterlegte oder ausdrücklich freigegebene Features erlaubt.
- Die Freigabe erfolgt als klarer, protokollierter Schritt im interaktiven Modus oder im Skriptmodus nur über die explizite Option:
  - `--allow-external-feature-sources <quelle>[,<quelle>...]` (`interaktiv`: Quelle bei Nachfrage bestätigen, `nicht-interaktiv`: alle Quellen als Flag-Argumente übergeben).
  - Die Option ist nur für diese Befehle gültig:
    - `u-boot init --devcontainer`
    - `u-boot generate devcontainer`
    - `u-boot config set devcontainer.featureSources.allow`
- Ein einzelnes `--allow-external-feature-sources` kann mehrere explizit erlaubte Quellen über Komma trennen.
- Die zugelassenen Quellen werden als explizit freigegebene Liste in der Projektkonfiguration gespeichert.
- Ohne explizit erlaubte Quelle führt der Versuch, externe Quellen zu nutzen, zu einem fachlichen Fehler (`code LH-FA-DEV-003`, Exit-Code `10`).
- `--yes` allein gilt nicht als Zustimmung für externe Quellen.

---

### LH-FA-DEV-004 – Benutzerrechte

Priorität: MVP

Der Devcontainer soll standardmäßig mit einem nicht-root Benutzer arbeiten.

---

### LH-FA-DEV-005 – Ports

Priorität: MVP

Das Produkt muss Ports aus aktivierten Services in der Devcontainer-Konfiguration berücksichtigen.

Konkret müssen die Ports der Services in `devcontainer.json` als `forwardPorts` eingetragen werden.
Ist keine aktive Port-Exposition in der aktuellen Projektkonfiguration vorhanden, darf `forwardPorts` fehlen.

---

## 4.4 Docker- und Compose-Unterstützung

### LH-FA-DOC-001 – Compose-Datei erzeugen

Priorität: MVP

Das Produkt muss eine `compose.yaml` erzeugen können.

---

### LH-FA-DOC-002 – Dockerfile erzeugen

Priorität: V1

Das Produkt soll bei Bedarf ein Dockerfile für die Anwendungsentwicklung erzeugen können.

Das Minimum bei aktivem Devcontainer ist die Erzeugung von `.devcontainer/Dockerfile`.

Zusätzlich kann optional ein separates Anwendungs-Dockerfile erzeugt werden:

- Standardpfad: `docker/Dockerfile`
- Konfigurierbar über Template-/Add-on-Konfiguration

---

### LH-FA-DOC-003 – Netzwerk

Priorität: MVP

Das Produkt muss ein gemeinsames Docker-Netzwerk für Services definieren können.

---

### LH-FA-DOC-004 – Volumes

Priorität: MVP

Das Produkt muss für aktivierte zustandsbehaftete Dienste persistente Volumes erzeugen können.

Beispiele:

- PostgreSQL-Daten (MVP)
- Keycloak-Daten (V1)
- OpenTelemetry-Konfiguration (V1)

---

### LH-FA-DOC-005 – Compose-Validierung

Priorität: V1

Das Produkt soll erzeugte Compose-Dateien auf syntaktische Gültigkeit prüfen können.

---

## 4.5 Service-Add-ons

### LH-FA-ADD-001 – Add-on-Befehl

Priorität: MVP

Das Produkt muss Services über folgenden Befehl hinzufügen können:

```bash
u-boot add <service>
```

Der Befehl ist nur in einem initialisierten `u-boot`-Projekt nutzbar (`u-boot.yaml` vorhanden).  
Ist keine gültige Projektkonfiguration vorhanden, ist mit klarer Fehlermeldung und Hinweis auf `u-boot init` abzubrechen.

---

### LH-FA-ADD-002 – PostgreSQL hinzufügen

Priorität: MVP

Das Produkt muss PostgreSQL als Service hinzufügen können.

Beispiel:

```bash
u-boot add postgres
```

Mindestumfang:

- PostgreSQL-Service in `compose.yaml`
- Volume für Daten
- `.env.example`-Einträge
- Port-Konfiguration
- Healthcheck

---

### LH-FA-ADD-003 – Keycloak hinzufügen

Priorität: V1

Das Produkt muss Keycloak als Service hinzufügen können.

Beispiel:

```bash
u-boot add keycloak
```

Mindestumfang:

- Keycloak-Service in `compose.yaml`
- Admin-Benutzer über `.env.example` als eindeutig markierte Beispielwerte (niemals reale Secrets), z. B.:
  - `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN`
  - `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`
- Port-Konfiguration
- optionale PostgreSQL-Anbindung bei konfigurierter persistenter externer Datenbank
- Healthcheck, soweit technisch sinnvoll

---

### LH-FA-ADD-004 – OpenTelemetry hinzufügen

Priorität: V1

Das Produkt muss OpenTelemetry-Komponenten hinzufügen können.

Beispiel:

```bash
u-boot add otel
```

Mindestumfang:

- OpenTelemetry Collector
- Collector-Konfigurationsdatei
- Compose-Service
- Standardports für OTLP
- Beispielkonfiguration für Logs, Metrics und Traces

---

### LH-FA-ADD-005 – Mehrfaches Hinzufügen verhindern

Priorität: MVP

Das Produkt muss erkennen, ob ein Service bereits vorhanden ist.

- Ein bereits vorhandener Service darf nicht doppelt eingefügt werden.
- Ein Service gilt als registriert, sobald `services.<name>` in `u-boot.yaml` existiert.
- Er gilt als aktiv vorhanden, wenn `services.<name>.enabled` explizit auf `true` steht **und** ein verwalteter Eintrag in `compose.yaml` existiert.
- `services.<name>.enabled` ist immer explizit zu setzen. Ein registrierter Service ohne expliziten `enabled`-Schlüssel gilt als deaktiviert (`false`) und führt bei `u-boot doctor` zu einer `warn`-Diagnose, die das explizite Setzen empfiehlt.
- Liegt `services.<name>.enabled: false` vor, gilt der Service als deaktiviert (weiterhin registriert), und `u-boot add <service>` darf ihn idempotent reaktivieren.
- Besteht `services.<name>` nicht in `u-boot.yaml`, aber ein verwalteter Block in `compose.yaml`, darf die Inkonsistenz nicht stillschweigend ignoriert werden. Der Befehl muss mit klarer Diagnose abbrechen und auf manuelle Bereinigung oder Re-Konfiguration verweisen.
- Besteht `services.<name>` in `u-boot.yaml` mit `enabled: true`, aber der verwaltete Compose-Eintrag fehlt, muss das Verhalten deterministisch sein: `u-boot add <service>` erzeugt den fehlenden Compose-Block wieder.

---

### LH-FA-ADD-006 – Add-on-Abhängigkeiten

Priorität: V1

Das Produkt muss Abhängigkeiten zwischen Add-ons erkennen.

Beispiele:

- Keycloak kann optional PostgreSQL benötigen, wenn `services.keycloak.persistence: external-postgres` in `u-boot.yaml` gesetzt ist.
- OpenTelemetry kann Beispielkonfigurationen für bestehende App-Services erzeugen.

Verhalten bei erkannter abhängiger Konfiguration:

- Ist `services.keycloak.persistence: external-postgres` gesetzt und PostgreSQL nicht vorhanden, darf der Aufruf nicht stillschweigend fortfahren.
- Ist die optionale Abhängigkeit nicht aktiv, darf Keycloak ohne PostgreSQL angelegt werden.

- Im interaktiven Modus (Standardmodus) muss das Produkt nachfragen, ob das fehlende Add-on automatisch hinzugefügt werden soll.
- Im nicht-interaktiven Modus (`--no-interactive`) ohne `--with-deps` muss das Produkt mit Exit-Code `10` abbrechen und auf die fehlende Abhängigkeit hinweisen.
- Über die Option `--with-deps` muss das Produkt fehlende Abhängigkeiten automatisch hinzufügen. `--with-deps` ist mit `--no-interactive` kombinierbar; in dem Fall werden Abhängigkeiten deterministisch und ohne Rückfrage installiert.
- Mit `--yes` (ohne `--with-deps`) wird die Standardentscheidung "Abhängigkeit hinzufügen" deterministisch ausgeführt, ohne dass eine Rückfrage gestellt wird.
- Mit `--yes` oder `--no-interactive` (jeweils exklusiv) muss das Verhalten in Skript-/CI-Umgebungen deterministisch und nicht-blockierend sein.

---

### LH-FA-ADD-007 – Service entfernen

Priorität: V1

Das Produkt muss einen Service wieder entfernen können.

Beispiel:

```bash
u-boot remove postgres
```

Der Befehl ist nur in einem initialisierten `u-boot`-Projekt nutzbar (`u-boot.yaml` vorhanden).  
Ist keine gültige Projektkonfiguration vorhanden, ist mit klarer Fehlermeldung und Hinweis auf `u-boot init` abzubrechen.

Mindestumfang:

- Service-Eintrag in `compose.yaml` entfernen
- zugehörige verwaltete Blöcke (z. B. in `.env.example`) entfernen
- Eintrag in `u-boot.yaml` auf `enabled: false` setzen
- Volumes nur auf explizite Anforderung (`--purge`) löschen
- Abhängigkeiten anderer Services prüfen und vor dem Entfernen warnen

`--purge` ist eine destruktive Operation; die Bestätigungs- und Modi-Regeln (`--yes`, `--no-interactive`, interaktive Rückfrage) sind verbindlich in `LH-FA-CLI-005A` definiert. Im nicht-interaktiven Modus ohne `--yes` muss der Aufruf mit Exit-Code `10` abgebrochen werden.

Ist der Service bereits auf `enabled: false`, darf der Aufruf idempotent als No-Op mit klarer Meldung beendet werden.

---

## 4.6 Starten und Stoppen der Umgebung

### LH-FA-UP-001 – Umgebung starten

Priorität: MVP

Das Produkt muss die Entwicklungsumgebung starten können.

`u-boot up` muss standardmäßig auf den Stabilisierungspfad der aktivierten Dienste warten, bevor der Befehl endet.

- Die Standard-Wartezeit beträgt 60 Sekunden; nach Ablauf erfolgt der Abbruch mit Fehler.
- Die maximale Wartezeit kann über `--timeout <sekunden>` überschrieben werden.
- `--timeout` akzeptiert nur nicht-negative Ganzzahlen (`>= 0`). Negative Werte führen zu Exit-Code `2` und klarer Validierungsfehlermeldung.
- Für Dienste mit Healthcheck ist `healthy` als Zielzustand erforderlich.
- Für Dienste ohne Healthcheck ist `running` als Zielzustand ausreichend.
- Bei definierten Ports wird auf Erreichbarkeit auf `localhost` geprüft, sofern es sich um TCP-basierten Zugriff handelt.
- Für nicht-TCP oder nicht eindeutig probebare Ports darf `up` nicht mit Fehler abbrechen; es ist ein strukturiertes `warn`-Diagnoseergebnis auszugeben.
- Mit `--timeout=0` wird auf das Warten verzichtet; `up` beendet nach Initiierung der Compose-Aktionen.

Beispiel:

```bash
u-boot up
```

---

### LH-FA-UP-002 – Docker Compose verwenden

Priorität: MVP

Der Befehl `u-boot up` muss intern Docker Compose verwenden können.

---

### LH-FA-UP-003 – Startstatus anzeigen

Priorität: MVP

Nach dem Start muss das Produkt den Status der relevanten Services anzeigen.

Mindestangaben:

- Service-Name
- Containerstatus
- Port
- Healthcheck-Status, falls vorhanden

---

### LH-FA-UP-004 – Umgebung stoppen

Priorität: MVP

Das Produkt muss die Umgebung stoppen können.

Beispiel:

```bash
u-boot down
```

Das Produkt muss zwischen einem regulären Stopp (Container stoppen) und einem vollständigen Aufräumen (Container und Volumes entfernen) unterscheiden:

```bash
u-boot down --volumes
```

---

### LH-FA-UP-005 – Logs anzeigen

Priorität: V1

Das Produkt soll Logs anzeigen können.

Beispiel:

```bash
u-boot logs
u-boot logs postgres
```

Mindestens müssen folgende Optionen unterstützt werden:

- `--follow` – Logs fortlaufend anzeigen
- `--tail <n>` – nur die letzten n Zeilen anzeigen

---

## 4.7 Diagnose

### LH-FA-DIAG-001 – Doctor-Befehl

Priorität: MVP

Das Produkt muss eine Diagnosefunktion bereitstellen.

Beispiel:

```bash
u-boot doctor
```

---

### LH-FA-DIAG-002 – Lokale Voraussetzungen prüfen

Priorität: MVP

Die Diagnosefunktion muss mindestens prüfen:

- Docker installiert (Mindestversion: 24.0.0 oder neuer)
- Docker erreichbar
- Docker Compose verfügbar (Mindestversion: 2.20.0 oder neuer)
- Git verfügbar
- Schreibrechte im Projektverzeichnis
- gültige `compose.yaml`, falls vorhanden
- gültige `u-boot.yaml`, falls vorhanden
- falls Devcontainer-Dateien vorhanden sind:
  - Ist `u-boot.yaml` vorhanden und `devcontainer.enabled == true`, müssen diese Prüfungen mit `error` bewertet werden:
    - syntaktische Gültigkeit von `.devcontainer/devcontainer.json`
    - Mindestkompatibilität mit VS Code Dev Containers (`name` gesetzt; mindestens `image` oder `build` vorhanden)
    - `forwardPorts`-Konsistenz zu aktivierten Services, falls Portangaben existieren
  - Ist `u-boot.yaml` vorhanden und `devcontainer.enabled == false`, sind die obigen Prüfungen optional (`warn`, keine harte Validierungspflicht).
  - Ist keine `u-boot.yaml` vorhanden, werden die obigen Prüfungen als ergänzende Qualitätsdiagnosen mit `warn` ausgegeben.
- falls `.devcontainer/Dockerfile` vorhanden ist: Lesbarkeit und erkennbare Build-Basisstruktur (`FROM` vorhanden)
- `forwardPorts`-Konsistenzregeln:
  - Für jeden aktivierten Service mit expliziter `ports`-Zuordnung (TCP) ist der Host-Port in `forwardPorts` enthalten.
  - bei mehreren TCP-Ports werden eindeutige Portzahlen eingetragen (Duplikate dedupliziert).
  - UDP- oder nicht eindeutig auflösbare Portangaben dürfen in `forwardPorts` fehlen; dafür ist ein `warn`-Diagnoseeintrag zulässig.

---

### LH-FA-DIAG-003 – Fehlerklassifikation

Priorität: MVP

Die Diagnosefunktion muss Probleme nach Schweregrad klassifizieren.

Mögliche Stufen:

- `ok`
- `warn`
- `error`

Die Diagnosefunktion muss den Exit Code an die höchste festgestellte Stufe binden:

- nur `ok` → Exit Code `0`
- mindestens `warn`, kein `error` → Exit Code `0`
- mindestens ein `error` → Exit Code ungleich `0`

`error` ist der einzige kritische Schweregrad.

Optional:

- Mit `--strict` muss mindestens ein `warn` zu einem Exit Code ungleich `0` führen.

---

### LH-FA-DIAG-004 – Reparaturhinweise

Priorität: MVP

Die Diagnosefunktion muss bei Problemen konkrete Reparaturhinweise ausgeben.

Beispiel:

```text
error: Docker daemon is not reachable.
hint: Start Docker or check your user permissions for /var/run/docker.sock.
```

---

## 4.8 Generatoren

### LH-FA-GEN-001 – Generate-Befehl

Priorität: MVP

Das Produkt muss Generatoren über folgenden Befehl anbieten:

```bash
u-boot generate <artifact>
```

Erlaubte Werte für `<artifact>`:

- `changelog`
- `readme`
- `env-example` (`.env.example`)
- `devcontainer`

Bei unbekanntem Artefakt muss der Befehl mit Exit Code `2` abbrechen und die erlaubten Werte explizit zurückgeben.

---

### LH-FA-GEN-002 – Changelog erzeugen

Priorität: MVP

Das Produkt muss ein Changelog erzeugen oder aktualisieren können.

Beispiel:

```bash
u-boot generate changelog
```

---

### LH-FA-GEN-003 – README erzeugen

Priorität: MVP

Das Produkt muss eine README-Datei erzeugen können.

Beispiel:

```bash
u-boot generate readme
```

---

### LH-FA-GEN-004 – Beispiel-ENV erzeugen

Priorität: MVP

Das Produkt muss eine `.env.example` erzeugen oder aktualisieren können.

Beispiel:

```bash
u-boot generate env-example
```

---

### LH-FA-GEN-005 – Idempotenz

Priorität: MVP

Generatoren müssen möglichst idempotent arbeiten.

Das bedeutet:

- mehrfaches Ausführen erzeugt keine unnötigen Duplikate
- bestehende manuelle Inhalte werden möglichst erhalten
- automatisch verwaltete Bereiche sind eindeutig markiert

---

## 4.9 Template-System

### LH-FA-TPL-001 – Projektvorlagen

Priorität: V1

Das Produkt soll Projektvorlagen unterstützen.

Beispiele:

```bash
u-boot init --template basic
u-boot init --template micronaut
u-boot init --template sveltekit
u-boot init --template micronaut-sveltekit
```

---

### LH-FA-TPL-002 – Template-Metadaten

Priorität: V1

Jedes Template soll Metadaten besitzen.

Mindestangaben:

- Name
- Beschreibung
- Version
- unterstützte Add-ons
- erzeugte Dateien
- benötigte Tools

---

### LH-FA-TPL-003 – Eigene Templates

Priorität: Later

Das Produkt soll später eigene lokale Templates unterstützen können.

Beispiel:

```bash
u-boot init --template ./my-template
```

---

### LH-FA-TPL-004 – Templates auflisten

Priorität: V1

Das Produkt muss verfügbare Templates auflisten können.

Beispiel:

```bash
u-boot template list
```

Die Ausgabe muss mindestens enthalten:

- Name
- Beschreibung
- Version

Die Ausgabe muss optional auch maschinenlesbar erfolgen können (`--json`).

---

## 4.10 Konfigurationsdatei

### LH-FA-CONF-001 – Projektkonfiguration

Priorität: MVP

Das Produkt muss eine eigene Projektkonfigurationsdatei verwenden.

Beispiel:

```text
u-boot.yaml
```

Die Konfiguration muss über den Konfigurationsbefehl gepflegt werden können:

```bash
u-boot config get project.name
u-boot config set project.name my-service
```

Die Migrationsfunktion ist in LH-FA-CONF-006 separat beschrieben.

---

### LH-FA-CONF-002 – Inhalt der Konfiguration

Priorität: MVP

Die Konfigurationsdatei muss mindestens enthalten:

```yaml
schemaVersion: 1
project:
  name: my-service

services:
  postgres:
    enabled: false

devcontainer:
  enabled: false

# Optionale, V1-relevante Felder:
# services:
#   keycloak:
#     enabled: false
#     persistence: embedded   # embedded | external-postgres
#   otel:
#     enabled: false
#
# devcontainer:
#   featureSources:
#     allow:
#       - https://ghcr.io/devcontainers/features/node
```

Hinweise:

- `services.keycloak.persistence` ist optional und kann `embedded` oder `external-postgres` sein. Fehlt der Wert, gilt der Default `embedded`.
- `template` ist optional und Bestandteil des V1-Template-Systems.
- `devcontainer.featureSources.allow` ist optional; fehlt das Feld, ist die Liste leer (nur lokale Features erlaubt).
- Erlaubte Einträge in `devcontainer.featureSources.allow` müssen gültige, non-empty Quell-Strings (z. B. `https://ghcr.io/devcontainers/features/node`) sein.
- Beim Schreiben wird die Liste dedupliziert.
- Bei ungültigen oder nicht zugelassenen Quellen ist ein fachlicher Validierungsfehler mit Code `10` zu melden.
- Nicht-mandatorische Add-ons dürfen im MVP auf `enabled: false` stehen.
- `services.<name>.enabled` ist immer explizit zu setzen; siehe `LH-FA-ADD-005` für die Default-Konvention.
- `enabled: false` bedeutet, dass der Service deaktiviert ist und bei erneutem `u-boot add <service>` wieder aktiviert werden kann.

---

### LH-FA-CONF-003 – Konfiguration lesen

Priorität: MVP

Das Produkt muss die Konfiguration lesen und bei Befehlen berücksichtigen können.

---

### LH-FA-CONF-004 – Konfiguration aktualisieren

Priorität: MVP

Das Produkt muss die Konfiguration aktualisieren können, wenn Add-ons hinzugefügt oder entfernt werden.

---

### LH-FA-CONF-005 – Konfiguration anzeigen und ändern

Priorität: MVP

Das Produkt muss einen Befehl zum Anzeigen und Ändern der Konfiguration bereitstellen.

Beispiele:

```bash
u-boot config                       # gesamte Konfiguration anzeigen
u-boot config get project.name      # einzelnen Wert anzeigen
u-boot config set project.name foo  # Wert setzen
```

Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden.

---

### LH-FA-CONF-006 – Konfiguration migrieren

Priorität: Later

Das Produkt muss ein Schema-Migrationskommando bereitstellen.

Beispiel:

```bash
u-boot config migrate
```

Die Migration muss mit einem klaren Fehler auf unbekannte Zukunftsversionen reagieren und bei Migrationen mit älteren Versionen eine Sicherung anlegen.

---

## 4.11 Build- und CI-Infrastruktur des u-boot-Projekts

Diese Sektion definiert die Build- und CI-Infrastruktur für die **u-boot-Codebase selbst**. Sie ist nicht zu verwechseln mit den Anforderungen aus 4.4 (`LH-FA-DOC-*`), die das Verhalten der Compose-/Dockerfile-Generatoren in den per `u-boot init` erzeugten Zielprojekten beschreiben.

Bezug:

- Implementierungssprache: `LH-OPEN-001` (Go).
- Vorlage: das Referenzprojekt `k-deskflight` (Docker-only-Workflow, Multi-Stage Dockerfile, Distroless-Runtime).

---

### LH-FA-BUILD-001 – Multi-Stage Dockerfile (u-boot-Repo)

Priorität: MVP

Die u-boot-Codebase muss ein Multi-Stage `Dockerfile` im Repo-Root bereitstellen.

Mindestumfang:

- BuildKit-Direktive in der ersten Zeile: `# syntax=docker/dockerfile:1.7`.
- Pflicht-Stages im MVP:
  - `deps` – Modulauflösung (`go mod download`) als Cache-Layer.
  - `compile` – schnelles Compile-Feedback (`go build`) ohne Tests/Lint.
  - `test` – `go test ./...`.
  - `lint` – `golangci-lint run ./...`.
  - `coverage` – `go test -coverprofile` + Coverage-Gate gegen `COVERAGE_THRESHOLD`; Bootstrap-Verhalten und Scope sind in `LH-FA-BUILD-008` und `LH-FA-BUILD-009` definiert.
  - `build` – statisch gelinktes Binary (`CGO_ENABLED=0`, `-ldflags="-s -w"`).
  - `runtime` – minimales Endimage (`LH-FA-BUILD-002`).
- Jeder Stage ist ein eigenständiges Build-Ziel und wird per `docker build --target <stage>` einzeln baubar.

---

### LH-FA-BUILD-002 – Runtime-Stage Pflichten

Priorität: MVP

Der `runtime`-Stage des u-boot-Dockerfiles muss folgende Eigenschaften erfüllen:

- Base-Image: `gcr.io/distroless/static-debian12:nonroot` (oder gleichwertig minimal und ohne Shell).
- Non-root-Benutzer (Distroless-`nonroot`-User, `USER 65532:65532`).
- `ENTRYPOINT` zeigt auf das im `build`-Stage erzeugte Binary; der konkrete Pfad (Empfehlung: `/usr/local/bin/u-boot`) ist im Dockerfile dokumentiert.
- OCI Image Labels gesetzt:
  - `org.opencontainers.image.source`
  - `org.opencontainers.image.description`
  - `org.opencontainers.image.licenses`
  - `org.opencontainers.image.title`
- Keine Build-Toolchain im Endimage; alle Build-Artefakte werden aus dem `build`-Stage per `COPY --from=build` übernommen.

---

### LH-FA-BUILD-003 – Build-Args und Pin-Politik

Priorität: MVP

Das u-boot-Dockerfile muss versions- und schwellwertbezogene Build-Args bereitstellen:

- `ARG GO_VERSION` – mit Default-Pin (z. B. `1.26.3`); Hebung ist Routine ohne separaten Spec-Eintrag.
- `ARG GOLANGCI_LINT_VERSION` – mit Default-Pin; gleiche Pin-Politik.
- `ARG COVERAGE_THRESHOLD` – mit Default `0` (bootstrap) und Override-Pfad `make coverage-gate THRESHOLD=…`.

Overrides erfolgen über `docker build --build-arg <NAME>=<value>` bzw. die korrespondierende Makefile-Variable.

---

### LH-FA-BUILD-004 – `.dockerignore` Pflicht

Priorität: MVP

Das u-boot-Repo muss eine `.dockerignore` im Repo-Root bereitstellen.

Mindestens auszuschließen:

- `.git`, `.gitignore`, `.github`
- IDE-Verzeichnisse: `.idea`, `.vscode`
- Agent-Verzeichnisse: `.claude`, `.codex`, `.agents`
- lokale Build-Artefakte und Caches (z. B. `dist/`, `coverage*`, `*.log`)

Die `.dockerignore` selbst gehört nicht ins Image und ist daher auszuschließen, sofern sie nicht von einem Stage-Build benötigt wird.

---

### LH-FA-BUILD-005 – Makefile mit Standard-Targets

Priorität: MVP

Das u-boot-Repo muss ein `Makefile` im Repo-Root bereitstellen.

Pflicht-Eigenschaften:

- `.DEFAULT_GOAL := help`
- `.PHONY` für alle Targets gesetzt
- `help`-Target mit Übersicht über alle verfügbaren Targets
- Variablen mit `?=`-Defaults für Overridability (`IMAGE`, `GO_VERSION`, `GOLANGCI_LINT_VERSION`, `THRESHOLD`)

MVP-Pflicht-Targets:

| Target          | Zweck                                                           |
| --------------- | --------------------------------------------------------------- |
| `help`          | Übersicht aller Targets                                         |
| `deps`          | `docker build --target deps`                                    |
| `compile`       | `docker build --target compile`                                 |
| `lint`          | `docker build --target lint`                                    |
| `test`          | `docker build --target test`                                    |
| `coverage`      | Alias auf `coverage-gate`                                       |
| `coverage-gate` | `docker build --target coverage --build-arg COVERAGE_THRESHOLD` |
| `build`         | `docker build --target runtime`                                 |
| `run`           | `docker run --rm <image> --help` (Smoketest); Dependency: `build` |
| `clean`         | lokale Artefakte und gebaute Images entfernen                   |

---

### LH-FA-BUILD-006 – Aggregator-Targets

Priorität: V1

Das Makefile soll Aggregator-Targets bereitstellen:

- `gates` – Inner-Loop-Pflichtgates (`lint` + `test` + `coverage-gate`), PR-blockierend.
- `ci` – `gates` plus mindestens `govulncheck` (bei Go-Stack aus `LH-OPEN-001` zwingend); weitere Prüfungen (z. B. Trivy-Image-Scan, SBOM) sind optional.
- `fullbuild` – `ci` plus `build`; vollständiger Closure-Lauf.

Aggregator-Targets müssen bei Fehler eines untergeordneten Targets mit Non-Zero-Exit abbrechen und die Fehlerursache klar benennen.

---

### LH-FA-BUILD-007 – Docker-only-Workflow

Priorität: MVP

Der Standard-Build-/Test-Workflow muss ohne hostseitige Sprach-Toolchain auskommen.

- Alle MVP- und V1-Pflicht-Targets aus `LH-FA-BUILD-005`/`LH-FA-BUILD-006` müssen ausschließlich `docker build`, `docker run` oder die Aggregation anderer solcher Targets aufrufen.
- Voraussetzung am Host: Docker Engine und `make`. `make` ist ein bewusster Carveout zu `LH-NFA-PORT-002` (weit verbreitet, einzige zusätzliche Host-Abhängigkeit neben Docker). Eine Go-Toolchain am Host darf für Standard-Targets nicht vorausgesetzt werden.
- Carveouts (z. B. ein Bash-Skript, das nicht containerisiert wird) sind im `Makefile`-Header explizit zu dokumentieren.

---

### LH-FA-BUILD-008 – Coverage-Bootstrap

Priorität: MVP

Der `coverage`-Stage muss in der Bootstrap-Phase (noch keine produktiven Pakete in `./internal/...`) deterministisch mit einer leeren Coverage-Eingabe umgehen können.

- Default-Schwellwert `0` (`ARG COVERAGE_THRESHOLD=0`).
- Sobald `./internal/...` produktive Pakete enthält, wird die Schwelle in einem Folge-Schritt angehoben; der Override-Pfad `make coverage-gate THRESHOLD=…` muss funktionieren.
- Leere Coverage darf in der Bootstrap-Phase nicht zu einem falschen Grün führen, das echte Test-Failures maskiert; der `go test`-Exit-Code wird über `set -o pipefail` o. ä. an die Gate-Logik durchgereicht.

---

### LH-FA-BUILD-009 – Repository-Layout

Priorität: MVP

Das u-boot-Repo muss folgendem Go-Layout folgen:

- Modul-Pfad in `go.mod`: `github.com/pt9912/u-boot`.
- Implementierungspakete leben unter `./internal/...`; öffentlich konsumierbare Pakete unter `./pkg/...` werden im MVP nicht erzeugt.
- CLI-Entry-Points unter `./cmd/<binary>/`; das primäre Binary heißt `uboot` (Verzeichnis `./cmd/uboot/`, Go-konform ohne Bindestrich), wird beim Build aber als `u-boot` ausgeliefert (`-o /out/u-boot`).
- Unit-Tests stehen als `*_test.go` neben dem produktiven Code im selben Paket.
- Coverage-Messung (`LH-FA-BUILD-001`, `LH-FA-BUILD-008`) bezieht sich auf `./internal/...`; `./cmd/...` ist bewusst ausgeschlossen, weil dort nur dünne Wireup-Logik liegt.

Mindestlayout:

```text
.
├── cmd/
│   └── uboot/
│       └── main.go              # Entry point der CLI (Wiring-Schicht)
├── internal/                    # nicht-exportierbare Implementierung
│   ├── hexagon/                 # innere Schichten (LH-FA-ARCH-002)
│   │   ├── domain/
│   │   ├── application/
│   │   └── port/{driving,driven}/
│   └── adapter/                 # äußere Schichten (LH-FA-ARCH-002)
│       ├── driving/
│       └── driven/
├── spec/                        # Lastenheft, weitere Spezifikationen
│   ├── lastenheft.md
│   └── architecture.md          # Architektur-Detailspec (LH-FA-ARCH-*)
├── docs/                        # Doku-Struktur (LH-FA-PROJDOCS-001)
├── go.mod
├── go.sum
├── Dockerfile
├── Makefile
├── .dockerignore
├── .gitignore
├── LICENSE
└── README.md
```

---

## 4.12 Doku-Struktur des u-boot-Projekts

Diese Sektion definiert die Verzeichnisstruktur unter `docs/` für die **u-boot-Codebase selbst**. Sie ist nicht zu verwechseln mit der `docs/`-Erzeugung in Zielprojekten (siehe `LH-FA-INIT-003`, `LH-SA-FILE-001`), die nur das Top-Level-Verzeichnis anlegt.

Vorlage: die Referenzprojekte `k-deskflight` und `grid-gym` (Basis-Pattern: archive + plan/adr + plan/planning-Lifecycle + user).

---

### LH-FA-PROJDOCS-001 – Mindeststruktur

Priorität: MVP

Das u-boot-Repo muss folgende `docs/`-Unterstruktur bereitstellen:

```text
docs/
├── archive/                  # abgelöste oder ersetzte Inhalte
├── plan/
│   ├── adr/                  # Architecture Decision Records
│   └── planning/
│       ├── open/             # Backlog
│       ├── next/             # priorisiert für nächsten Schritt
│       ├── in-progress/      # aktiv bearbeitet
│       └── done/             # abgeschlossen
└── user/                     # User-facing Dokumentation
```

Jedes Unterverzeichnis muss mindestens eine `README.md` mit kurzer Zweckbeschreibung enthalten, damit Git die Struktur trackt und Newcomer den Verzeichnisstandard ohne externe Erklärung erfassen können. `.gitkeep` ist als Ersatz unzureichend, weil er den Zweck nicht kommuniziert.

Abgrenzung zu Zielprojekten: Für per `u-boot init` erzeugte Zielprojekte ist nur `docs/` als Top-Level Pflicht (`LH-FA-INIT-003`). Ob diese Unterstruktur auch in Zielprojekten erzeugt wird, ist eine spätere Entscheidung (z. B. via Template oder Flag) und gehört nicht zum MVP-Umfang.

---

### LH-FA-PROJDOCS-002 – ADR-Format

Priorität: MVP

Architecture Decision Records in `docs/plan/adr/` müssen folgenden Konventionen folgen:

- Dateiname beginnt mit vierstelliger Nummer, beginnend bei `0001` und monoton steigend: `0001-<slug>.md`, `0002-<slug>.md`.
- Slug nach der Nummer in Kebab-Case (z. B. `0001-implementierungssprache-go.md`).
- Mindestabschnitte im Dokument, in dieser Reihenfolge, jeweils als `##`-Überschrift:
  1. Dokumenttitel als `#`-Überschrift: `# ADR <Nr>: <Titel>`.
  2. `## Status` – einer aus `Proposed`, `Accepted`, `Superseded by <NNNN>-<slug>`, `Deprecated`.
  3. `## Datum` – Entscheidungsdatum im Format `YYYY-MM-DD`.
  4. `## Kontext` – warum die Entscheidung nötig wird.
  5. `## Entscheidung` – was beschlossen wird.
  6. `## Konsequenzen` – kurz- und langfristige Auswirkungen, inkl. Trade-offs.
- ADR-Nummern werden nie wiederverwendet; abgelöste ADRs bleiben mit Status `Superseded by <NNNN>-<slug>` erhalten und verweisen auf den Nachfolger über den vollen Dateinamen-Stamm (ohne `.md`).

---

### LH-FA-PROJDOCS-003 – Planning-Lifecycle

Priorität: MVP

Planning-Artefakte (Slices, Tranchen, Tickets) durchlaufen die Verzeichnisse `open → next → in-progress → done` in dieser Reihenfolge.

- Ein Artefakt darf nicht in mehreren Lifecycle-Verzeichnissen gleichzeitig liegen.
- Übergänge zwischen Lifecycle-Stufen erfolgen per `git mv` (Move statt Kopie), damit die Datei-Historie erhalten bleibt.
- Inhalte in `done/` dürfen nachträglich nur korrigierend (Tippfehler, Querverweise, Archiv-Hinweise) verändert werden; substanzielle inhaltliche Änderungen erzeugen ein neues Artefakt in `open/` oder `next/` mit Verweis auf den vorhergehenden Stand.
- Dateinamen in `planning/` folgen einem der zwei verbindlichen Formate, abhängig vom Artefakttyp:
  - `slice-<phase>-<kebab-slug>.md` für Slice-Pläne (z. B. `slice-m1-repo-skeleton.md`).
  - `tranche-<nr>-<kebab-slug>.md` für Tranchen-Pläne (z. B. `tranche-01-init-flow.md`).
  Die Wahl zwischen Slice- und Tranchen-Format ist im `README.md` von `docs/plan/planning/` dokumentiert; ein Artefakt verwendet genau eines der beiden Formate.
- Ausnahme für übergreifende Master-Dokumente: eine `roadmap.md` darf direkt unter `docs/plan/planning/in-progress/` liegen und folgt keinem der beiden Formate. Sie fasst Slices und Tranchen lebendig zusammen und wird laufend gepflegt.

---

### LH-FA-PROJDOCS-005 – Carveout-Disziplin

Priorität: MVP

Jeder **temporäre Carveout** in der u-boot-Codebase muss parallel zu seiner Entstehung einen Slice-Plan in `docs/plan/planning/open/` bekommen, der die Aufhebungsbedingung benennt. Sobald der Slice priorisiert wird, wandert er per `git mv` nach `next/` (`LH-FA-PROJDOCS-003`).

Als temporärer Carveout zählt insbesondere:

- Bootstrap-Schwellwerte (z. B. `COVERAGE_THRESHOLD=0` bis erste produktive Pakete existieren).
- Bewusst leere Regelblöcke in Linter-/Tooling-Konfiguration (z. B. `gomodguard_v2.blocked: {}` bis externe Modul-Dependencies vorhanden sind, `depguard rules: {}` in Bootstrap-Phasen).
- Prospektive Doku-Phrasen ("scharf zu schalten mit M3", "wird mit V1 ergänzt", "Logging-Port kommt später", "Folgepflicht im GitHub-UI").
- Bewusst weggelassene Pflichten in einem CI-/Build-Setup, deren Aufhebung in einem ADR-Folgepunkt vermerkt ist (z. B. Image-Publish, Image-Scan, Branch-Protection).

Pflichten:

- Der Slice-Plan folgt der Dateiname-Konvention aus `LH-FA-PROJDOCS-003` (`slice-<phase>-<slug>.md`).
- Der Plan benennt mindestens: Auslöser (was wurde wo bewusst weggelassen), Aufhebungsbedingung (was muss passieren), Akzeptanzkriterien.
- Wo der Carveout in einer Spec-Anforderung dokumentiert ist (z. B. `LH-FA-BUILD-008` für Coverage-Bootstrap), referenziert der Spec-Text den Slice-Plan.
- **Doppelte Verankerung:** jeder temporäre Carveout ist sowohl in [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md) als auch in [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) als Slice-Zeile sichtbar. Carveouts ohne Roadmap-Eintrag oder Slice-Pläne ohne Carveout-Inventar-Verweis sind Verstoß gegen diese Anforderung.
- Auch Spec-Open-Punkte (`LH-OPEN-*`) und ADR-Folgepunkte gelten als temporäre Carveouts und brauchen einen Slice-Plan — kein „bleibt offen bis MVP-Closure" als Inventar-Eintrag.
- Ein Master-Inventar in `carveouts.md` listet alle aktuellen Carveouts mit Status (`temporär` + Plan-Verweis vs. `permanent` + Begründung). Diese Datei lebt analog zur `roadmap.md` dauerhaft in `in-progress/`.

Permanente Carveouts (z. B. `errcheck.exclude-functions` für CLI-Writes, `testpackage`/`gochecknoglobals` für die Wiring-Schicht `cmd/uboot/`) sind ebenfalls im Master-Inventar zu listen, brauchen aber keinen Aufhebungsplan; sie tragen den Status `permanent` mit kurzer Begründung.

---

### LH-FA-PROJDOCS-004 – Archivierung

Priorität: V1

Abgelöste oder veraltete Inhalte aus `user/`, `plan/` oder anderen `docs/`-Bereichen werden nach `docs/archive/` verschoben, statt sie zu löschen.

- Beim Verschieben wird ein kurzer Hinweis am Anfang des Zielfiles ergänzt (z. B. `> Archiviert am YYYY-MM-DD; ersetzt durch [<Pfad>](<pfad>).`).
- Querverweise in lebendiger Doku werden auf das neue Ziel umgebogen oder explizit als historisch markiert.
- Das Verschieben erfolgt per `git mv`, damit die Historie erhalten bleibt.

---

## 4.13 Architektur des u-boot-Projekts

Diese Sektion definiert das Architektur-Pattern für die **u-boot-Codebase selbst**. Detail-Spezifikation (Schichten, Import-Regeln, Beispiele, `depguard`-Konfiguration) liegt in [`spec/architecture.md`](architecture.md); diese Sektion definiert die Pflichten, die für jede Code-Änderung gelten.

Vorlage: die Referenzprojekte `k-deskflight` (Go, flach), `m-trace` (TypeScript, driving/driven-Split) und `grid-gym` (Python, driving/driven-Split). Begründung der konkreten Variante siehe ADR-0002.

---

### LH-FA-ARCH-001 – Hexagonales Pattern

Priorität: MVP

Die u-boot-Codebase muss dem hexagonalen Architektur-Pattern (Ports & Adapters) folgen.

Pflichten:

- Trennung zwischen reiner Domäne, Anwendungslogik, Ports (Interfaces) und Adaptern (Implementierungen).
- Keine direkten Abhängigkeiten von Anwendungslogik oder Domäne zu externen Bibliotheken (Docker-SDK, YAML-Parser, Dateisystem).
- Externe Zugriffe laufen ausschließlich über Driven-Ports, die in Adapter-Paketen implementiert werden.

Detail und Anti-Patterns: `spec/architecture.md`, ADR-0002.

---

### LH-FA-ARCH-002 – Schichten und Verzeichnislayout

Priorität: MVP

Das u-boot-Repo muss folgende Schichten unter `internal/` bereitstellen:

```text
internal/
├── hexagon/
│   ├── domain/         # reine Datentypen + invariantes Verhalten, keine I/O
│   ├── application/    # Use-Cases; ruft ausschließlich Ports auf
│   └── port/
│       ├── driving/    # Interfaces, die von außen (CLI/HTTP) konsumiert werden
│       └── driven/     # Interfaces, die das Application nach außen ruft
└── adapter/
    ├── driving/        # konkrete Driver (z. B. cli/ mit Cobra-Commands)
    └── driven/         # konkrete Adapter (z. B. docker/, fs/, yaml/)
```

Die Wiring-Schicht (`cmd/uboot/`) ist die einzige Stelle, an der `application` und `adapter` zusammen importiert werden dürfen.

---

### LH-FA-ARCH-003 – Import-Regeln und Enforcement

Priorität: MVP

Die verbindliche Import-Regel-Tabelle, die Begründung der einzelnen Schicht-Pflichten und die Anti-Patterns sind in [`spec/architecture.md`](architecture.md) §3 definiert (Single Source of Truth).

Pflichten:

- Die Regeln werden im `lint`-Stage (`LH-FA-BUILD-001`) per `golangci-lint` mit `depguard` durchgesetzt; Verstöße sind PR-blockierend.
- Die `depguard`-Konfiguration in `.golangci.yml` ist deckungsgleich mit der Regel-Tabelle aus `spec/architecture.md` §3 zu halten; Drift wird im Review zurückgewiesen.
- `//nolint:depguard`-Pragmas sind verboten. Carveouts werden zentral in `.golangci.yml` mit `Why:`-Kommentar dokumentiert.
- `depguard`-Regeln gelten production-only; `*_test.go`-Dateien sind ausgenommen, damit Tests Fakes und Test-Libraries (`testify`, …) frei nutzen können.
- Solange einzelne Schichten noch keine produktiven Pakete enthalten, dürfen `depguard`-Regelblöcke aktiv sein und nichts treffen — die Schicht-Regeln greifen automatisch, sobald das erste produktive Paket angelegt wird.

---

## 5. Nichtfunktionale Anforderungen

## 5.1 Benutzbarkeit

### LH-NFA-USE-001 – Verständliche Bedienung

Priorität: MVP

Das Produkt muss ohne tiefes Vorwissen über die interne Implementierung bedienbar sein.

---

### LH-NFA-USE-002 – Klare Befehle

Priorität: MVP

Befehle müssen sprechend, konsistent und kurz sein.

Beispiele:

```bash
u-boot init
u-boot add postgres
u-boot doctor
u-boot up
```

---

### LH-NFA-USE-003 – Lesbare Ausgaben

Priorität: MVP

CLI-Ausgaben müssen klar strukturiert und gut lesbar sein.

---

### LH-NFA-USE-004 – Maschinenlesbare Ausgabe

Priorität: V1

Das Produkt soll optional maschinenlesbare Ausgabe unterstützen.

Beispiel:

```bash
u-boot doctor --json
```

Für `--dry-run`/`--diff`-Kombinationen gilt zusätzlich die JSON-Ausgabe in `LH-FA-CLI-007` und `LH-FA-CLI-008`.

Für alle `--json`-Ausgaben gilt ergänzend ein gemeinsames Minimalkontrakt-Schema:

- `status` (`ok`/`warn`/`error`)
- `command` (Hauptbefehl als in `LH-FA-CLI-007` definiertes Enum)
- optional `subcommand` (für gruppierte Befehle wie `template` oder `config`)
- `diagnostics` (Liste von Objekten mit mind. `level`, `code`, `message`, optional `file`)
- `exitCode` (vgl. `LH-FA-CLI-006`)

Für `--json`-Antworten gilt zusätzlich:

- `diagnostics`, wenn leer, darf als `[]` ausgegeben werden.
- `diagnostics.level` darf nur `warn` oder `error` enthalten.
- `diagnostics.code` folgt der Konvention: LH-Kennung der verursachenden Anforderung (z. B. `LH-FA-DEV-003`). Tool-interne Codes ohne LH-Bezug dürfen nur dann verwendet werden, wenn ihre Bedeutung in der Dokumentation festgehalten ist (Verweis: `LH-FA-CLI-007`).
- `diagnostics.file` ist optional.
- `status` ist an den höchsten in `diagnostics` enthaltenen `level` gekoppelt: `error` → `status == "error"`; `warn` ohne `error` → `status == "warn"`; sonst `status == "ok"`.
- Bei `command == "template"` oder `command == "config"` ist `subcommand` verpflichtend.
- Die Felder `status`, `command`, `diagnostics` und `exitCode` sind minimal verpflichtend und sollten mit anderen Feldern ergänzt werden.

Für normale (`--json` ohne `--dry-run`/`--diff`) Ausgaben ist der obige Minimalkontrakt bindend.
Für `--dry-run`- oder `--diff`-Ausgaben mit `--json` gilt zusätzlich das vollständige Schema aus `LH-FA-CLI-007` als bindender Pflichtkontrakt (inkl. `plannedFiles`, `changes`, `dryRun`, `diff`).

Beispiel:

```json
{
  "status": "ok",
  "command": "doctor",
  "diagnostics": [],
  "exitCode": 0
}
```

---

## 5.2 Zuverlässigkeit

### LH-NFA-REL-001 – Kein stilles Überschreiben

Priorität: MVP

Das Produkt darf bestehende Dateien nicht stillschweigend überschreiben.

---

### LH-NFA-REL-002 – Wiederholbare Ausführung

Priorität: MVP

Wiederholte Ausführung desselben Befehls darf das Projekt nicht beschädigen.

---

### LH-NFA-REL-003 – Abbruch bei kritischen Fehlern

Priorität: MVP

Bei kritischen Fehlern muss das Produkt abbrechen und eine klare Fehlermeldung ausgeben.

---

### LH-NFA-REL-004 – Validierung erzeugter Dateien

Priorität: MVP

Das Produkt soll erzeugte Dateien validieren, soweit passende Validatoren verfügbar sind.

Beispiele:

- YAML
- JSON
- Docker Compose

---

## 5.3 Wartbarkeit

### LH-NFA-MAINT-001 – Modulare Architektur

Priorität: MVP

Das Produkt muss modular aufgebaut sein.

Insbesondere sollen Add-ons, Templates und Generatoren voneinander getrennt implementiert werden.

---

### LH-NFA-MAINT-002 – Erweiterbarkeit

Priorität: MVP

Neue Services müssen mit geringem Aufwand ergänzt werden können.

---

### LH-NFA-MAINT-003 – Testbarkeit

Priorität: MVP

Die fachlichen Funktionen müssen automatisiert testbar sein.

---

### LH-NFA-MAINT-004 – Dokumentierte Schnittstellen

Priorität: V1

Interne Schnittstellen für Add-ons und Templates sollen dokumentiert werden.

---

## 5.4 Portabilität

### LH-NFA-PORT-001 – Linux-Unterstützung

Priorität: MVP

Das Produkt muss Linux als primäre Plattform unterstützen.

---

### LH-NFA-PORT-002 – Keine unnötigen Systemabhängigkeiten

Priorität: MVP

Das Produkt soll möglichst wenige externe Systemabhängigkeiten benötigen.

---

### LH-NFA-PORT-003 – Containerfreundlichkeit

Priorität: V1

Das Produkt soll selbst in einem Container oder Devcontainer ausführbar sein können.

---

## 5.5 Sicherheit

### LH-NFA-SEC-001 – Keine Secrets einchecken

Priorität: MVP

Das Produkt darf keine echten Secrets in erzeugte Dateien schreiben.

---

### LH-NFA-SEC-002 – Beispielwerte markieren

Priorität: MVP

Beispielwerte in `.env.example` müssen eindeutig als Beispielwerte erkennbar sein.

---

### LH-NFA-SEC-003 – Sichere Defaults

Priorität: MVP

Das Produkt soll sichere Standardwerte verwenden, soweit dies mit lokaler Entwicklung vereinbar ist.

---

### LH-NFA-SEC-004 – Keine verdeckte Ausführung fremder Skripte

Priorität: MVP

Das Produkt darf keinen externen ausführbaren Code aus nicht freigegebenen Quellen ohne ausdrückliche Zustimmung des Nutzers ausführen.

Konkretisierung des Begriffs "externer Code aus nicht freigegebenen Quellen":

- Devcontainer-Features, Templates oder andere Skripte, die nicht lokal im Repository liegen oder nicht ausdrücklich in `u-boot.yaml` freigegeben sind (siehe `LH-FA-DEV-003`, `devcontainer.featureSources.allow`).
- ad-hoc geladene Shell-, Python- oder ähnliche Skripte über HTTP(S) oder andere Netzwerk-Quellen.

Nicht erfasst sind:

- Docker-Images, die in `compose.yaml`-Services oder im Devcontainer-Build explizit konfiguriert sind – sie werden durch `docker pull` regulär aus konfigurierten Registries bezogen und gelten als bewusst gewählte Abhängigkeit des Projekts.
- Pakete, die innerhalb einer Image-Build-Pipeline (z. B. in einem Dockerfile via Paketmanager) installiert werden.

Die Zustimmung ist im interaktiven Modus durch explizite Rückfrage und im nicht-interaktiven Modus durch entsprechende Flag-Optionen (z. B. `--allow-external-feature-sources`) einzuholen.

---

## 5.6 Performance

### LH-NFA-PERF-001 – Schnelle CLI-Antwort

Priorität: MVP

Einfache Befehle müssen auf einem typischen Entwicklungsrechner innerhalb folgender Zeiten reagieren (gemessen ohne Docker-Kommunikation, Kaltstart):

MVP:

- `u-boot --help`, `u-boot --version` – unter 200 ms
- `u-boot doctor` (ohne Netz-Wartezeit) – unter 2 s

V1:

- `u-boot config get …` – unter 300 ms

---

### LH-NFA-PERF-002 – Startzeit abhängig von Docker

Priorität: MVP

Die Startzeit von `u-boot up` darf von Docker-Images und Services abhängen, muss aber transparent dargestellt werden.

Insbesondere muss der Fortschritt einzelner Services (Pull, Create, Start, Healthcheck) sichtbar sein.

---

## 6. Schnittstellenanforderungen

## 6.1 Kommandozeilenschnittstelle

### LH-SA-CLI-001 – Befehlsstruktur

Priorität: MVP

Die CLI soll folgende Grundstruktur verwenden:

```bash
u-boot <command> [subcommand|args...] [options]
```

`subcommand` ist für kommandospezifische Unterbefehle reserviert (z. B. `template`, `config`).
Positionsargumente (z. B. `postgres`, `project.name`) stehen ebenfalls vor den Optionen.

---

### LH-SA-CLI-002 – Vorgesehene Befehle

Priorität: MVP/V1 gemischt (siehe Spalte)

| Befehl                       | Zweck                              | Priorität |
| ---------------------------- | ---------------------------------- | --------- |
| `u-boot init`                | Projekt initialisieren             | MVP       |
| `u-boot add <service>`       | Service hinzufügen                 | MVP       |
| `u-boot remove <service>`    | Service entfernen                  | V1        |
| `u-boot up`                  | Umgebung starten                   | MVP       |
| `u-boot down`                | Umgebung stoppen                   | MVP       |
| `u-boot doctor`              | Umgebung prüfen                    | MVP       |
| `u-boot logs`                | Logs anzeigen                      | V1        |
| `u-boot generate <artifact>` | Artefakt erzeugen                  | MVP       |
| `u-boot config`              | Konfiguration anzeigen oder ändern | MVP       |
| `u-boot config migrate`      | Konfigurationsschema migrieren     | Later     |
| `u-boot template list`       | Templates anzeigen                 | V1        |

---

## 6.2 Dateischnittstellen

### LH-SA-FILE-001 – Erzeugte Dateien

Priorität: MVP

Das Produkt soll folgende Dateien erzeugen oder aktualisieren können:

```text
README.md
CHANGELOG.md
compose.yaml
.env.example
.gitignore
u-boot.yaml
.devcontainer/devcontainer.json
.devcontainer/Dockerfile
docker/
scripts/
docs/
```

Optional, sobald ein Anwendungs-Dockerfile (`LH-FA-DOC-002`) erzeugt wird:

```text
.dockerignore
```

---

### LH-SA-FILE-002 – Markierte verwaltete Bereiche

Priorität: MVP

Automatisch verwaltete Bereiche in Dateien sollen markiert werden.

Markierungsformat je Dateityp:

- YAML, `.env`, `Dockerfile`, Shell-Skripte (`#`-Kommentare):

  ```yaml
  # BEGIN U-BOOT MANAGED BLOCK: postgres
  # ...
  # END U-BOOT MANAGED BLOCK: postgres
  ```

- Markdown (`README.md`, `CHANGELOG.md`) als HTML-Kommentar:

  ```markdown
  <!-- BEGIN U-BOOT MANAGED BLOCK: postgres -->
  ...
  <!-- END U-BOOT MANAGED BLOCK: postgres -->
  ```

- JSONC (z. B. `.devcontainer/devcontainer.json`):

  ```jsonc
  // BEGIN U-BOOT MANAGED BLOCK: postgres
  // ...
  // END U-BOOT MANAGED BLOCK: postgres
  ```

- Strikte JSON-Dateien ohne Kommentar-Support werden nicht inline markiert; die gesamte Datei gilt als verwaltet, und der verwaltete Status ist in `u-boot.yaml` zu hinterlegen.

---

## 6.3 Docker-Schnittstelle

### LH-SA-DOCKER-001 – Docker Compose

Priorität: MVP

Das Produkt muss Docker Compose aufrufen oder kompatible Compose-Dateien erzeugen können.

---

### LH-SA-DOCKER-002 – Containerstatus

Priorität: MVP

Das Produkt muss den Status laufender Container auslesen können.

---

## 7. Datenanforderungen

### LH-DA-001 – Projektmetadaten

Priorität: MVP

Das Produkt muss Projektmetadaten speichern können.

Beispiele:

- Projektname
- Template
- aktivierte Services
- Ports
- Version des `u-boot`-Schemas

---

### LH-DA-002 – Service-Metadaten

Priorität: MVP

Das Produkt muss Informationen über aktivierte Services speichern können.

Beispiele:

- Name
- Image
- Ports
- Volumes
- Environment-Variablen
- Healthchecks

---

### LH-DA-003 – Schema-Version

Priorität: MVP

Die Projektkonfiguration muss eine Schema-Version enthalten.

Beispiel:

```yaml
schemaVersion: 1
```

---

### LH-DA-004 – Schema-Migration

Priorität: Later

Das Produkt muss mit älteren Schema-Versionen umgehen können.

Anforderungen:

- Eine ältere `schemaVersion` muss erkannt und gemeldet werden.
- Das Produkt muss eine automatische Migration anbieten (`u-boot config migrate`).
- Vor der Migration muss eine Sicherungsdatei nach der Backup-Konvention aus `LH-FA-INIT-005` erzeugt werden: primär `u-boot.yaml.bak`; ist bereits ein Backup vorhanden, wird der kleinste freie numerische Suffix verwendet (`u-boot.yaml.bak.1`, `u-boot.yaml.bak.2`, ...) ohne bestehende Backups zu überschreiben.
- Eine unbekannte (zu neue) `schemaVersion` muss zu einem klaren Fehler führen und das Tool darf in diesem Fall keine Dateien verändern.

---

## 8. Qualitätsanforderungen

### LH-QA-001 – Automatisierte Tests

Priorität: MVP

Für zentrale Funktionen müssen automatisierte Tests vorhanden sein.

Mindestumfang:

- CLI-Befehle
- Dateigeneratoren
- Template-Verarbeitung
- Add-on-Erzeugung
- Konfigurationsparser

---

### LH-QA-002 – Testbare Akzeptanzkriterien

Priorität: MVP

Jede funktionale Anforderung soll durch mindestens einen Akzeptanztest überprüfbar sein.

---

### LH-QA-003 – CI-Fähigkeit (GitHub Actions)

Priorität: MVP

Das u-boot-Repo muss eine CI-Pipeline auf GitHub Actions führen.

Pflicht-Komposition:

- Workflow-Datei: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- Trigger: `pull_request` und `push` auf den Branch `main`.
- Zwei Jobs, parallel, beide PR-blockierend (Required-Status-Checks im GitHub-UI nach dem ersten grünen Lauf zu konfigurieren):
  - `gates` — führt `make gates` aus (lint + test + coverage-gate, `LH-FA-BUILD-005`/`-006`).
  - `security-gates` — führt `make govulncheck` aus (`LH-FA-BUILD-006`).
- Runner: `ubuntu-latest`. Keine Host-Go-Toolchain (`LH-FA-BUILD-007`); der Runner braucht nur das vorinstallierte Docker + BuildKit.
- Actions sind **SHA-gepinnt** mit Tag-Kommentar (Supply-Chain-Härtung gegen Tag-Move). Pin-Hebung ist Routine; neuer Commit-SHA via `gh api repos/<owner>/<repo>/git/refs/tags/<tag>`.
- Top-Level `permissions: {}` (alle Tokens entzogen); jeder Job lockert auf das Minimum (Defense-in-Depth).
- Jeder Job mit `timeout-minutes` versehen (Empfehlung: 20).
- Begründung der konkreten Setzungen in [`docs/plan/adr/0004-ci-system.md`](../docs/plan/adr/0004-ci-system.md).

---

### LH-QA-004 – Linting (SOLID-nahes Lint-Profil)

Priorität: MVP

Die u-boot-Codebase muss ein verschärftes Lint-Profil führen, das über die Default-Linter hinausgeht.

Profil-Komposition (29 Linter insgesamt):

- 5 Default-Linter (`govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`).
- 24 SOLID-nahe Zusatz-Linter (Komplexitäts-, Funktionslänge-, Interface-, Kopplungs- und Boundary-Signale; vollständige Liste in [`docs/user/quality.md`](../docs/user/quality.md) §1.2). **`depguard`** für die Schicht-Regeln aus `LH-FA-ARCH-003` ist Teil dieser 24.

Pflichten:

- Die Konfiguration lebt in `.golangci.yml` (v2-Schema).
- Schwellen und Linter-Settings sind in [`docs/user/quality.md`](../docs/user/quality.md) §1.2 dokumentiert; bei Drift zwischen Doku und Config gewinnt die Doku, Config ist anzupassen.
- `//nolint`-Pragmas sind verboten. Pro-Pfad-Carveouts (z. B. Tests, `cmd/uboot`) werden zentral in `.golangci.yml` unter `issues.exclude-rules` mit `Why:`-Kommentar dokumentiert.
- Verstöße brechen den `lint`-Stage (`LH-FA-BUILD-001`) und damit `make gates`/`make ci`/`make fullbuild`.
- Begründung der konkreten Linter-Auswahl ist in [`docs/plan/adr/0003-solid-nahes-lint-profil.md`](../docs/plan/adr/0003-solid-nahes-lint-profil.md) festgehalten.

---

## 9. Akzeptanzkriterien

### LH-AK-001 – Minimaler Init-Flow

Priorität: MVP

Vorbedingung: eine erreichbare Docker-Engine und Docker Compose in den jeweils geforderten Mindestversionen (`LH-FA-DIAG-002`, `LH-RISK-001`).

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
mkdir demo
cd demo
u-boot init
u-boot doctor
```

Erwartetes Ergebnis:

- Projektstruktur wurde erzeugt
- `u-boot doctor` meldet keinen `error`-Eintrag
- vorhandene Dateien wurden nicht ungewollt überschrieben

---

### LH-AK-002 – PostgreSQL-Flow

Priorität: MVP

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init
u-boot add postgres
u-boot up
```

Erwartetes Ergebnis:

- PostgreSQL-Service ist in `compose.yaml` vorhanden
- `.env.example` enthält PostgreSQL-Variablen (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`)
- Container ist gestartet und erreicht den Healthcheck-Status `healthy` innerhalb von 60 Sekunden
- der konfigurierte Port (Standard: `5432`) ist auf `localhost` erreichbar

---

### LH-AK-003 – Keycloak-Flow

Priorität: V1

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init
u-boot add keycloak
u-boot up
```

Erwartetes Ergebnis:

- Keycloak-Service ist in `compose.yaml` vorhanden
- Admin-Zugangsdaten werden über `.env.example` mit Platzhaltern dokumentiert (z. B. `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`)
- Web-Oberfläche ist über den konfigurierten Port (Standard: `8080`) auf `localhost` erreichbar (HTTP 200 oder 302 auf `/`)

---

### LH-AK-004 – OpenTelemetry-Flow

Priorität: V1

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init
u-boot add otel
u-boot up
```

Erwartetes Ergebnis:

- OpenTelemetry Collector ist konfiguriert
- Collector-Konfigurationsdatei wurde erzeugt und ist syntaktisch gültig
- OTLP/gRPC ist auf `localhost:4317` erreichbar
- OTLP/HTTP ist auf `localhost:4318` erreichbar
- Collector-Container erreicht innerhalb von 60 Sekunden den Status `running` oder `healthy`

---

### LH-AK-005 – Devcontainer-Flow

Priorität: MVP

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init --devcontainer
```

Erwartetes Ergebnis:

- `.devcontainer/devcontainer.json` existiert
- `.devcontainer/Dockerfile` existiert
- `devcontainer.json` enthält mindestens `name`, und mindestens eines aus `build` oder `image`
- `devcontainer.json` enthält `forwardPorts`, sofern mindestens ein Add-on aktive Ports exportiert
- `u-boot doctor` enthält keinen `error` zu `devcontainer`-Konfiguration oder Feature-Quellen

---

### LH-AK-006 – Idempotenz

Priorität: MVP

Folgender Ablauf darf keine Duplikate erzeugen:

```bash
u-boot add postgres
u-boot add postgres
```

Erwartetes Ergebnis:

- PostgreSQL ist nur einmal in der Konfiguration vorhanden
- das Tool gibt eine verständliche Meldung aus

---

### LH-AK-007 – Changelog-Generator

Priorität: MVP

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot generate changelog
```

Erwartetes Ergebnis:

- `CHANGELOG.md` existiert
- vorhandene Inhalte werden nicht zerstört
- neuer Abschnitt wird korrekt ergänzt oder vorbereitet

---

## 10. Abgrenzung

### LH-ABG-001 – Kein vollständiges Deployment-System

`u-boot` ist in der ersten Version kein vollständiges Produktionsdeployment-System.

Nicht im Kernumfang enthalten:

- Kubernetes-Produktionsdeployment
- Cloud-Provisioning
- Terraform-Management
- Secret-Management für Produktion

---

### LH-ABG-002 – Keine IDE-Abhängigkeit

`u-boot` darf VS Code Dev Containers unterstützen, soll aber nicht ausschließlich davon abhängig sein.

---

### LH-ABG-003 – Kein Ersatz für Docker Compose

`u-boot` soll Docker Compose nicht ersetzen, sondern erzeugen, konfigurieren und komfortabel verwenden.

---

## 11. Risiken und Annahmen

### LH-RISK-001 – Docker-Versionen

Unterschiedliche Docker- und Compose-Versionen können zu Kompatibilitätsproblemen führen.

Maßnahme:

- Mindestversionen dokumentieren (`Docker Engine >= 24.0.0`, `Docker Compose >= 2.20.0`)
- `u-boot doctor` prüft Versionen
- Die Mindestversionen sind als harte Voraussetzung im Lastenheft, in der README und in der CLI-Hilfe dokumentiert.

---

### LH-RISK-002 – Überschreiben manueller Änderungen

Automatische Generatoren können manuelle Änderungen beschädigen.

Maßnahme:

- verwaltete Blöcke
- Backups
- Diff-Anzeige
- `--force` nur explizit

---

### LH-RISK-003 – Zu großer Funktionsumfang

Das Projekt kann durch zu viele Templates und Services unübersichtlich werden.

Maßnahme:

- MVP klar begrenzen
- Add-on-System modular halten
- stabile Kernbefehle priorisieren

---

## 12. MVP-Umfang

### LH-MVP-001 – Muss im MVP enthalten sein

Der MVP muss enthalten:

- `u-boot init`
- `u-boot doctor`
- `u-boot add postgres`
- `u-boot generate changelog`
- `u-boot generate readme`
- `u-boot generate env-example`
- `u-boot generate devcontainer`
- `u-boot up`
- `u-boot down`
- `u-boot config` (`get`, `set`, gesamte Konfiguration anzeigen)
- Erzeugung von `compose.yaml`
- Erzeugung von `.env.example`
- Erzeugung von `README.md`
- Erzeugung von `CHANGELOG.md`
- grundlegende Devcontainer-Unterstützung
- Build-Infrastruktur der u-boot-Codebase selbst:
  - Multi-Stage `Dockerfile` (`LH-FA-BUILD-001`)
  - `Makefile` mit MVP-Pflicht-Targets (`LH-FA-BUILD-005`)
  - `.dockerignore` (`LH-FA-BUILD-004`)
  - Docker-only-Workflow (`LH-FA-BUILD-007`)
  - Repository-Layout nach `LH-FA-BUILD-009`
- Doku-Struktur der u-boot-Codebase nach `LH-FA-PROJDOCS-001`, inkl. ADR-Format (`LH-FA-PROJDOCS-002`), Planning-Lifecycle (`LH-FA-PROJDOCS-003`) und Carveout-Disziplin (`LH-FA-PROJDOCS-005`)
- Architektur-Pattern (hexagonal, driving/driven-Split) nach `LH-FA-ARCH-001..003`, mit Detail-Spezifikation in `spec/architecture.md` und Import-Enforcement via `golangci-lint depguard`
- SOLID-nahes Lint-Profil nach `LH-QA-004` (5 Default-Linter + 24 SOLID-nahe Linter inkl. `depguard`, 29 Linter gesamt); Konfiguration in `.golangci.yml`, Doku in `docs/user/quality.md` §1.2 / §1.3, Begründung in ADR-0003
- CI-Pipeline nach `LH-QA-003` (GitHub Actions, `.github/workflows/ci.yml`, Jobs `gates` + `security-gates`, beide PR-blockierend); Begründung in ADR-0004

---

### LH-MVP-002 – Kann nach dem MVP folgen

Nach dem MVP können ergänzt werden:

- `u-boot add keycloak`
- `u-boot add otel`
- Template-System
- lokale Custom-Templates
- Plugin-System

---

## 13. Traceability-Matrix

`-` bedeutet in der aktuellen Fassung noch keine getrennte Pflichtenheft-/Testfallableitung.

| Lastenheft-Kennung | Kurzbeschreibung               | Priorität | Spätere Ableitung im Pflichtenheft | Testfall        |
| ------------------ | ------------------------------ | --------- | ---------------------------------- | --------------- |
| LH-LESE-001        | Modalverben                   | -         | -                                  | -               |
| LH-LESE-002        | Sprache                       | -         | -                                  | -               |
| LH-ZB-001          | Projektziel                   | -         | -                                  | -               |
| LH-ZB-002          | Produktvision                 | -         | -                                  | -               |
| LH-ZB-003          | Repo-Beschreibung             | -         | -                                  | -               |
| LH-PE-001          | Anwendungsbereich             | -         | -                                  | -               |
| LH-PE-002          | Zielgruppen                   | -         | -                                  | -               |
| LH-PE-003          | Betriebsumgebung              | -         | -                                  | -               |
| LH-PÜ-001          | Grundfunktion                 | -         | -                                  | -               |
| LH-PÜ-002          | Hauptmodule                   | -         | -                                  | -               |
| LH-MOD-001         | Projektinitialisierung        | -         | -                                  | -               |
| LH-MOD-002         | Devcontainer-Generator       | -         | -                                  | -               |
| LH-MOD-003         | Docker-Stack-Generator       | -         | -                                  | -               |
| LH-MOD-004         | Service-Add-ons              | -         | -                                  | -               |
| LH-MOD-005         | Umgebungsprüfung             | -         | -                                  | -               |
| LH-MOD-006         | Stack-Start                  | -         | -                                  | -               |
| LH-MOD-007         | Generatoren                  | -         | -                                  | -               |
| LH-MOD-008         | Template-System              | -         | -                                  | -               |
| LH-FA-CLI-001      | CLI-Aufruf                     | MVP       | PH-CLI-001                         | TC-CLI-001      |
| LH-FA-CLI-002      | Hilfeausgabe                   | MVP       | PH-CLI-002                         | TC-CLI-002      |
| LH-FA-CLI-003      | Versionsausgabe                | MVP       | PH-CLI-003                         | TC-CLI-003      |
| LH-FA-CLI-004      | Fehlerausgabe                  | MVP       | PH-CLI-004                         | TC-CLI-004      |
| LH-FA-CLI-005      | Verbosity und Logging          | MVP       | PH-CLI-005                         | TC-CLI-005      |
| LH-FA-CLI-005A     | Interaktivität und Automatisierung | MVP    | PH-CLI-005A                        | TC-CLI-005A     |
| LH-FA-CLI-006      | Exit Codes                     | MVP       | PH-CLI-006                         | TC-CLI-006      |
| LH-FA-CLI-007      | Dry Run                        | V1        | PH-CLI-007                         | TC-CLI-007      |
| LH-FA-CLI-008      | Diff-Ausgabe                   | V1        | PH-CLI-008                         | TC-CLI-008      |
| LH-FA-INIT-001     | Neues Projekt initialisieren    | MVP       | PH-INIT-001                        | TC-INIT-001     |
| LH-FA-INIT-002     | Projektname                    | MVP       | PH-INIT-002                        | TC-INIT-002     |
| LH-FA-INIT-003     | Projektstruktur erzeugen        | MVP       | PH-INIT-003                        | TC-INIT-003     |
| LH-FA-INIT-004     | Bestehendes Projekt erkennen    | MVP       | PH-INIT-004                        | TC-INIT-004     |
| LH-FA-INIT-005     | Überschreibschutz              | MVP       | PH-INIT-005                        | TC-INIT-005     |
| LH-FA-INIT-006     | Projektnamen-Validierung       | MVP       | PH-INIT-006                        | TC-INIT-006     |
| LH-FA-INIT-007     | Git-Repository-Initialisierung | MVP       | PH-INIT-007                        | TC-INIT-007     |
| LH-FA-DEV-001      | Devcontainer erzeugen          | MVP       | PH-DEV-001                         | TC-DEV-001      |
| LH-FA-DEV-002      | VS-Code-Kompatibilität         | MVP       | PH-DEV-002                         | TC-DEV-002      |
| LH-FA-DEV-003      | Devcontainer-Features          | V1        | PH-DEV-003                         | TC-DEV-003      |
| LH-FA-DEV-004      | Benutzerrechte                 | MVP       | PH-DEV-004                         | TC-DEV-004      |
| LH-FA-DEV-005      | Ports                          | MVP       | PH-DEV-005                         | TC-DEV-005      |
| LH-FA-DOC-001      | Compose-Datei erzeugen         | MVP       | PH-DOC-001                         | TC-DOC-001      |
| LH-FA-DOC-002      | Dockerfile erzeugen            | V1        | PH-DOC-002                         | TC-DOC-002      |
| LH-FA-DOC-003      | Netzwerk                       | MVP       | PH-DOC-003                         | TC-DOC-003      |
| LH-FA-DOC-004      | Volumes                        | MVP       | PH-DOC-004                         | TC-DOC-004      |
| LH-FA-DOC-005      | Compose-Validierung            | V1        | PH-DOC-005                         | TC-DOC-005      |
| LH-FA-ADD-001      | Add-on-Befehl                  | MVP       | PH-ADD-001                         | TC-ADD-001      |
| LH-FA-ADD-002      | PostgreSQL hinzufügen           | MVP       | PH-ADD-002                         | TC-ADD-002      |
| LH-FA-ADD-003      | Keycloak hinzufügen            | V1        | PH-ADD-003                         | TC-ADD-003      |
| LH-FA-ADD-004      | OpenTelemetry hinzufügen       | V1        | PH-ADD-004                         | TC-ADD-004      |
| LH-FA-ADD-005      | Mehrfaches Hinzufügen verhindern| MVP       | PH-ADD-005                         | TC-ADD-005      |
| LH-FA-ADD-006      | Add-on-Abhängigkeiten          | V1        | PH-ADD-006                         | TC-ADD-006      |
| LH-FA-ADD-007      | Service entfernen              | V1        | PH-ADD-007                         | TC-ADD-007      |
| LH-FA-UP-001       | Umgebung starten               | MVP       | PH-UP-001                          | TC-UP-001       |
| LH-FA-UP-002       | Docker Compose verwenden        | MVP       | PH-UP-002                          | TC-UP-002       |
| LH-FA-UP-003       | Startstatus anzeigen           | MVP       | PH-UP-003                          | TC-UP-003       |
| LH-FA-UP-004       | Umgebung stoppen               | MVP       | PH-UP-004                          | TC-UP-004       |
| LH-FA-UP-005       | Logs anzeigen                  | V1        | PH-UP-005                          | TC-UP-005       |
| LH-FA-DIAG-001     | Doctor-Befehl                  | MVP       | PH-DIAG-001                        | TC-DIAG-001     |
| LH-FA-DIAG-002     | Lokale Voraussetzungen prüfen   | MVP       | PH-DIAG-002                        | TC-DIAG-002     |
| LH-FA-DIAG-003     | Fehlerklassifikation           | MVP       | PH-DIAG-003                        | TC-DIAG-003     |
| LH-FA-DIAG-004     | Reparaturhinweise              | MVP       | PH-DIAG-004                        | TC-DIAG-004     |
| LH-FA-GEN-001      | Generate-Befehl                | MVP       | PH-GEN-001                         | TC-GEN-001      |
| LH-FA-GEN-002      | Changelog erzeugen             | MVP       | PH-GEN-002                         | TC-GEN-002      |
| LH-FA-GEN-003      | README erzeugen                | MVP       | PH-GEN-003                         | TC-GEN-003      |
| LH-FA-GEN-004      | Beispiel-ENV erzeugen          | MVP       | PH-GEN-004                         | TC-GEN-004      |
| LH-FA-GEN-005      | Idempotenz                     | MVP       | PH-GEN-005                         | TC-GEN-005      |
| LH-FA-TPL-001      | Projektvorlagen                | V1        | PH-TPL-001                         | TC-TPL-001      |
| LH-FA-TPL-002      | Template-Metadaten             | V1        | PH-TPL-002                         | TC-TPL-002      |
| LH-FA-TPL-003      | Eigene Templates               | Later     | PH-TPL-003                         | TC-TPL-003      |
| LH-FA-TPL-004      | Templates auflisten            | V1        | PH-TPL-004                         | TC-TPL-004      |
| LH-FA-CONF-001     | Projektkonfiguration           | MVP       | PH-CONF-001                        | TC-CONF-001     |
| LH-FA-CONF-002     | Inhalt der Konfiguration       | MVP       | PH-CONF-002                        | TC-CONF-002     |
| LH-FA-CONF-003     | Konfiguration lesen            | MVP       | PH-CONF-003                        | TC-CONF-003     |
| LH-FA-CONF-004     | Konfiguration aktualisieren    | MVP       | PH-CONF-004                        | TC-CONF-004     |
| LH-FA-CONF-005     | Konfiguration anzeigen/ändern  | MVP       | PH-CONF-005                        | TC-CONF-005     |
| LH-FA-CONF-006     | Konfiguration migrieren        | Later     | PH-CONF-006                        | TC-CONF-006     |
| LH-FA-BUILD-001    | Multi-Stage Dockerfile (u-boot-Repo) | MVP | PH-BUILD-001                       | TC-BUILD-001    |
| LH-FA-BUILD-002    | Runtime-Stage Pflichten        | MVP       | PH-BUILD-002                       | TC-BUILD-002    |
| LH-FA-BUILD-003    | Build-Args und Pin-Politik     | MVP       | PH-BUILD-003                       | TC-BUILD-003    |
| LH-FA-BUILD-004    | `.dockerignore` Pflicht        | MVP       | PH-BUILD-004                       | TC-BUILD-004    |
| LH-FA-BUILD-005    | Makefile mit Standard-Targets  | MVP       | PH-BUILD-005                       | TC-BUILD-005    |
| LH-FA-BUILD-006    | Aggregator-Targets             | V1        | PH-BUILD-006                       | TC-BUILD-006    |
| LH-FA-BUILD-007    | Docker-only-Workflow           | MVP       | PH-BUILD-007                       | TC-BUILD-007    |
| LH-FA-BUILD-008    | Coverage-Bootstrap             | MVP       | PH-BUILD-008                       | TC-BUILD-008    |
| LH-FA-BUILD-009    | Repository-Layout              | MVP       | PH-BUILD-009                       | TC-BUILD-009    |
| LH-FA-PROJDOCS-001 | docs/-Mindeststruktur (u-boot-Repo) | MVP  | PH-PROJDOCS-001                    | TC-PROJDOCS-001 |
| LH-FA-PROJDOCS-002 | ADR-Format                     | MVP       | PH-PROJDOCS-002                    | TC-PROJDOCS-002 |
| LH-FA-PROJDOCS-003 | Planning-Lifecycle             | MVP       | PH-PROJDOCS-003                    | TC-PROJDOCS-003 |
| LH-FA-PROJDOCS-004 | Archivierung                   | V1        | PH-PROJDOCS-004                    | TC-PROJDOCS-004 |
| LH-FA-PROJDOCS-005 | Carveout-Disziplin             | MVP       | PH-PROJDOCS-005                    | TC-PROJDOCS-005 |
| LH-FA-ARCH-001     | Hexagonales Pattern            | MVP       | PH-ARCH-001                        | TC-ARCH-001     |
| LH-FA-ARCH-002     | Schichten und Verzeichnislayout | MVP      | PH-ARCH-002                        | TC-ARCH-002     |
| LH-FA-ARCH-003     | Import-Regeln und Enforcement  | MVP       | PH-ARCH-003                        | TC-ARCH-003     |
| LH-DA-003          | Schema-Version                 | MVP       | PH-DA-003                          | TC-DA-003       |
| LH-DA-004          | Schema-Migration               | Later     | PH-DA-004                          | TC-DA-004       |
| LH-SA-CLI-001      | Befehlsstruktur                | MVP       | PH-SA-CLI-001                      | TC-SA-CLI-001   |
| LH-SA-CLI-002      | Vorgesehene Befehle            | MVP/V1    | PH-SA-CLI-002                      | TC-SA-CLI-002   |
| LH-SA-FILE-001     | Erzeugte Dateien               | MVP       | PH-SA-FILE-001                     | TC-SA-FILE-001  |
| LH-SA-FILE-002     | Markierte verwaltete Bereiche  | MVP       | PH-SA-FILE-002                     | TC-SA-FILE-002  |
| LH-SA-DOCKER-001    | Docker Compose                 | MVP       | PH-SA-DOCKER-001                   | TC-SA-DOCKER-001 |
| LH-SA-DOCKER-002    | Containerstatus                | MVP       | PH-SA-DOCKER-002                   | TC-SA-DOCKER-002 |
| LH-NFA-USE-001     | Verständliche Bedienung        | MVP       | PH-NFA-USE-001                     | TC-NFA-USE-001  |
| LH-NFA-USE-002     | Klare Befehle                 | MVP       | PH-NFA-USE-002                     | TC-NFA-USE-002  |
| LH-NFA-USE-003     | Lesbare Ausgaben              | MVP       | PH-NFA-USE-003                     | TC-NFA-USE-003  |
| LH-NFA-USE-004     | Maschinenlesbare Ausgabe       | V1        | PH-NFA-USE-004                     | TC-NFA-USE-004  |
| LH-NFA-REL-001     | Kein stilles Überschreiben     | MVP       | PH-NFA-REL-001                     | TC-NFA-REL-001  |
| LH-NFA-REL-002     | Wiederholbare Ausführung       | MVP       | PH-NFA-REL-002                     | TC-NFA-REL-002  |
| LH-NFA-REL-003     | Abbruch bei kritischen Fehlern | MVP       | PH-NFA-REL-003                     | TC-NFA-REL-003  |
| LH-NFA-REL-004     | Validierung erzeugter Dateien  | MVP       | PH-NFA-REL-004                     | TC-NFA-REL-004  |
| LH-NFA-MAINT-001   | Modulare Architektur           | MVP       | PH-NFA-MAINT-001                   | TC-NFA-MAINT-001|
| LH-NFA-MAINT-002   | Erweiterbarkeit                | MVP       | PH-NFA-MAINT-002                   | TC-NFA-MAINT-002|
| LH-NFA-MAINT-003   | Testbarkeit                    | MVP       | PH-NFA-MAINT-003                   | TC-NFA-MAINT-003|
| LH-NFA-MAINT-004   | Dokumentierte Schnittstellen    | V1        | PH-NFA-MAINT-004                   | TC-NFA-MAINT-004|
| LH-NFA-PORT-001    | Linux-Unterstützung            | MVP       | PH-NFA-PORT-001                    | TC-NFA-PORT-001 |
| LH-NFA-PORT-002    | Keine unnötigen Systemabhängigkeiten | MVP       | PH-NFA-PORT-002                    | TC-NFA-PORT-002 |
| LH-NFA-PORT-003    | Containerfreundlichkeit        | V1        | PH-NFA-PORT-003                    | TC-NFA-PORT-003 |
| LH-NFA-SEC-001     | Keine Secrets einchecken       | MVP       | PH-NFA-SEC-001                     | TC-NFA-SEC-001  |
| LH-NFA-SEC-002     | Beispielwerte markieren        | MVP       | PH-NFA-SEC-002                     | TC-NFA-SEC-002  |
| LH-NFA-SEC-003     | Sichere Defaults               | MVP       | PH-NFA-SEC-003                     | TC-NFA-SEC-003  |
| LH-NFA-SEC-004     | Keine verdeckte Ausführung fremder Skripte | MVP       | PH-NFA-SEC-004                     | TC-NFA-SEC-004  |
| LH-NFA-PERF-001    | Schnelle CLI-Antwort           | MVP       | PH-NFA-PERF-001                    | TC-NFA-PERF-001 |
| LH-NFA-PERF-002    | Startzeit abhängig von Docker   | MVP       | PH-NFA-PERF-002                    | TC-NFA-PERF-002 |
| LH-DA-001          | Projektmetadaten               | MVP       | PH-DA-001                          | TC-DA-001       |
| LH-DA-002          | Service-Metadaten              | MVP       | PH-DA-002                          | TC-DA-002       |
| LH-QA-001          | Automatisierte Tests           | MVP       | PH-QA-001                          | TC-QA-001       |
| LH-QA-002          | Testbare Akzeptanzkriterien    | MVP       | PH-QA-002                          | TC-QA-002       |
| LH-QA-003          | CI-Fähigkeit                  | MVP       | PH-QA-003                          | TC-QA-003       |
| LH-QA-004          | Linting (SOLID-nahes Profil)   | MVP       | PH-QA-004                          | TC-QA-004       |
| LH-AK-001          | Minimaler Init-Flow            | MVP       | PH-AK-001                          | TC-AK-001       |
| LH-AK-002          | PostgreSQL-Flow                | MVP       | PH-AK-002                          | TC-AK-002       |
| LH-AK-003          | Keycloak-Flow                  | V1        | PH-AK-003                          | TC-AK-003       |
| LH-AK-004          | OpenTelemetry-Flow             | V1        | PH-AK-004                          | TC-AK-004       |
| LH-AK-005          | Devcontainer-Flow              | MVP       | PH-AK-005                          | TC-AK-005       |
| LH-AK-006          | Idempotenz                     | MVP       | PH-AK-006                          | TC-AK-006       |
| LH-AK-007          | Changelog-Generator            | MVP       | PH-AK-007                          | TC-AK-007       |
| LH-ABG-001         | Kein vollständiges Deployment-System | -      | PH-ABG-001                         | TC-ABG-001      |
| LH-ABG-002         | Keine IDE-Abhängigkeit         | -                                  | PH-ABG-002                         | TC-ABG-002      |
| LH-ABG-003         | Kein Ersatz für Docker Compose  | -                                  | PH-ABG-003                         | TC-ABG-003      |
| LH-RISK-001        | Docker-Versionen               | -                                  | PH-RISK-001                        | TC-RISK-001     |
| LH-RISK-002        | Überschreiben manueller Änderungen | -                               | PH-RISK-002                        | TC-RISK-002     |
| LH-RISK-003        | Zu großer Funktionsumfang      | -                                  | PH-RISK-003                        | TC-RISK-003     |
| LH-MVP-001         | Muss im MVP enthalten sein     | MVP                               | -                                  | -               |
| LH-MVP-002         | Kann nach dem MVP folgen       | -                                  | -                                  | -               |
| LH-OPEN-001        | Implementierungssprache (Go, entschieden 2026-05-21) | - | -                          | -               |
| LH-OPEN-002        | Paketierung (GHCR entschieden 2026-05-31 via ADR-0007; Restwege vertagt/verworfen) | -                                  | -                                  | -               |
| LH-OPEN-003        | Plugin-System                  | -                                  | -                                  | -               |
| LH-OPEN-004        | Template-Format                | -                                  | -                                  | -               |

---

## 14. Offene Punkte und Entscheidungen

### LH-OPEN-001 – Implementierungssprache (entschieden)

Status: entschieden am 2026-05-21.
Sprache: **Go**.
Begründung und Konsequenzen: siehe `docs/plan/adr/0001-implementierungssprache-go.md` (`LH-FA-PROJDOCS-002`).

Mindest-Toolchain: Go 1.26 oder neuer (`go 1.26.0` in `go.mod`, analog Referenzprojekt `k-deskflight`); Default-Pin im Dockerfile als `ARG GO_VERSION` (aktuell `1.26.3`, die aktuelle Stable-Version am Entscheidungsdatum). Pin-Hebung ist Routine ohne separaten Spec-Eintrag.

---

### LH-OPEN-002 – Paketierung

Status: **GHCR entschieden** am 2026-05-31 (ADR-0007), Restwege
vertagt/verworfen. Formell offen, bis alle Restwege entschieden sind.

| Option | Status | Verweis |
| ------ | ------ | ------- |
| Container Image (GHCR `ghcr.io/pt9912/u-boot`) | **Gewählt** | [ADR-0007](../docs/plan/adr/0007-distributionswege-ghcr.md), Umsetzung in [`slice-v1-release-pipeline`](../docs/plan/planning/done/slice-v1-release-pipeline.md) T2/T3 |
| Einzelnes Binary | Vertagt mit Trigger | `slice-v2-binary-distribution.md` (bei Auslösung) |
| Homebrew | Vertagt mit Trigger | `slice-v2-homebrew-formula.md` (bei Auslösung) |
| Debian/RPM | Vertagt mit Trigger | `slice-v2-distro-pakete.md` (bei Auslösung) |
| npm package | Verworfen | ADR-0007 §Entscheidung — Sprach-Ökosystem-Mismatch |
| pip package | Verworfen | ADR-0007 §Entscheidung — Sprach-Ökosystem-Mismatch |

---

### LH-OPEN-003 – Plugin-System

Es ist zu klären, ob Add-ons langfristig fest eingebaut oder als Plugins nachladbar sein sollen.

---

### LH-OPEN-004 – Template-Format

Das genaue Format für Templates ist noch festzulegen.

Mögliche Optionen:

- YAML-Metadaten plus Dateivorlagen
- Cookiecutter-kompatible Templates
- eigenes Template-System
- OCI-basierte Template-Pakete

---

## 15. Glossar

| Begriff              | Bedeutung                                                                                 |
| -------------------- | ----------------------------------------------------------------------------------------- |
| CLI                  | Command Line Interface                                                                    |
| Bootstrapper         | Werkzeug zur initialen Bereitstellung einer Umgebung oder Projektstruktur                 |
| Compose              | Kurzform für Docker Compose; auch Bezeichnung für die Datei `compose.yaml`                |
| Docker Compose       | Werkzeug zum Definieren und Starten mehrerer Container                                    |
| Devcontainer         | Containerisierte Entwicklungsumgebung, häufig mit VS Code verwendet                       |
| Devcontainer-Feature | Optionaler, wiederverwendbarer Baustein für Devcontainer (z. B. Node.js, Docker CLI)      |
| Add-on               | Erweiterbarer Service-Baustein wie PostgreSQL oder Keycloak                               |
| Template             | Wiederverwendbare Projektvorlage                                                          |
| Healthcheck          | Prüfung, ob ein Service technisch funktionsfähig ist                                      |
| Idempotenz           | Mehrfaches Ausführen führt zum gleichen stabilen Ergebnis                                 |
| Managed Block        | Automatisch verwalteter Bereich in einer Datei (Markierung: `BEGIN U-BOOT MANAGED BLOCK`) |
| OTLP                 | OpenTelemetry Protocol – Protokoll zur Übertragung von Logs, Metrics und Traces           |
