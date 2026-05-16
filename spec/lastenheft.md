# Lastenheft: `u-boot` – Projekt-Bootstrapper für Docker-/Devcontainer-Stacks

| Dokument         | Lastenheft                                                         |
| ---------------- | ------------------------------------------------------------------ |
| Projektname      | `u-boot`                                                           |
| Kurzbeschreibung | CLI-Tool zum Bootstrapping reproduzierbarer Entwicklungsumgebungen |
| Zielplattform    | Linux, Docker, VS Code Dev Containers                              |
| Hauptnutzer      | Softwareentwickler, DevOps-Engineers, technische Teams             |
| Version          | 0.1.0                                                              |
| Status           | Entwurf                                                            |
| Datum            | 2026-05-16                                                         |

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
u-boot add keycloak
u-boot add otel
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
u-boot add keycloak
u-boot add otel
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

Das Produkt muss eine konfigurierbare Ausgabeverbosität unterstützen.

Mindestens müssen folgende Stufen unterstützt werden:

- `--quiet` – nur Fehler
- Standard – Statusmeldungen und Fehler
- `--verbose` – zusätzlich Detailinformationen
- `--debug` – zusätzlich interne Diagnoseausgaben

---

### LH-FA-CLI-006 – Exit Codes

Priorität: MVP

Das Produkt muss aussagekräftige Exit Codes liefern.

Mindestens:

- `0` – Erfolg
- `1` – allgemeiner Fehler
- `2` – fehlerhafte CLI-Nutzung
- `>2` – fachliche Fehler nach Kategorie (z. B. Docker nicht erreichbar)

---

### LH-FA-CLI-007 – Dry Run

Priorität: V1

Das Produkt muss für dateiverändernde Befehle einen Dry-Run-Modus unterstützen.

Beispiel:

```bash
u-boot add postgres --dry-run
```

Der Dry-Run muss anzeigen, welche Dateien erzeugt, geändert oder gelöscht würden, ohne Änderungen am Dateisystem vorzunehmen.

---

### LH-FA-CLI-008 – Diff-Ausgabe

Priorität: V1

Das Produkt soll bei dateiverändernden Befehlen eine Diff-Ausgabe unterstützen.

Beispiel:

```bash
u-boot add postgres --diff
```

Die Diff-Ausgabe muss Unterschiede zwischen aktuellem und geplantem Zustand der betroffenen Dateien zeigen.

---

## 4.2 Projektinitialisierung

### LH-FA-INIT-001 – Neues Projekt initialisieren

Priorität: MVP

Das Produkt muss mit folgendem Befehl ein neues Projekt initialisieren können:

```bash
u-boot init
```

---

### LH-FA-INIT-002 – Projektname

Priorität: MVP

Das Produkt muss bei der Initialisierung einen Projektnamen verwenden können.

Beispiel:

```bash
u-boot init my-service
```

---

### LH-FA-INIT-003 – Projektstruktur erzeugen

Priorität: MVP

Das Produkt muss eine grundlegende Projektstruktur erzeugen.

Mindestumfang:

```text
.
├── .devcontainer/
├── docker/
├── scripts/
├── docs/
├── README.md
├── CHANGELOG.md
├── compose.yaml
├── .env.example
└── .gitignore
```

---

### LH-FA-INIT-004 – Bestehendes Projekt erkennen

Priorität: MVP

Das Produkt muss erkennen, ob es in einem bestehenden Projektverzeichnis ausgeführt wird.

Falls bereits relevante Dateien vorhanden sind, darf das Produkt diese nicht kommentarlos überschreiben.

---

### LH-FA-INIT-005 – Überschreibschutz

Priorität: MVP

Das Produkt muss vor dem Überschreiben bestehender Dateien schützen.

Standardverhalten ohne Option:

- Abbruch mit Hinweis, welche Datei kollidiert

Zusätzliche Strategien über Option:

- `--backup` – bestehende Datei als `<name>.bak` sichern und ersetzen
- `--force` – bestehende Dateien ohne Rückfrage überschreiben; vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden

---

### LH-FA-INIT-006 – Projektnamen-Validierung

Priorität: MVP

Das Produkt muss den Projektnamen validieren.

Regeln:

- erlaubt sind Kleinbuchstaben, Ziffern und Bindestrich
- beginnt mit einem Kleinbuchstaben
- endet mit einem Kleinbuchstaben oder einer Ziffer
- maximale Länge: 63 Zeichen
- regulärer Ausdruck: `^[a-z][a-z0-9-]{0,61}[a-z0-9]$`

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

---

### LH-FA-DEV-004 – Benutzerrechte

Priorität: MVP

Der Devcontainer soll standardmäßig mit einem nicht-root Benutzer arbeiten.

---

### LH-FA-DEV-005 – Ports

Priorität: MVP

Das Produkt muss Ports aus aktivierten Services in der Devcontainer-Konfiguration berücksichtigen.

Konkret müssen die Ports der Services in `devcontainer.json` als `forwardPorts` eingetragen werden.

---

## 4.4 Docker- und Compose-Unterstützung

### LH-FA-DOC-001 – Compose-Datei erzeugen

Priorität: MVP

Das Produkt muss eine `compose.yaml` erzeugen können.

---

### LH-FA-DOC-002 – Dockerfile erzeugen

Priorität: V1

Das Produkt soll bei Bedarf ein Dockerfile für die Anwendungsentwicklung erzeugen können.

---

### LH-FA-DOC-003 – Netzwerk

Priorität: MVP

Das Produkt muss ein gemeinsames Docker-Netzwerk für Services definieren können.

---

### LH-FA-DOC-004 – Volumes

Priorität: MVP

Das Produkt muss persistente Volumes für zustandsbehaftete Dienste erzeugen können.

Beispiele:

- PostgreSQL-Daten
- Keycloak-Daten
- OpenTelemetry-Konfiguration

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
- Admin-Benutzer über `.env`
- Port-Konfiguration
- optional PostgreSQL-Anbindung
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

Ein bereits vorhandener Service darf nicht doppelt eingefügt werden.

---

### LH-FA-ADD-006 – Add-on-Abhängigkeiten

Priorität: V1

Das Produkt muss Abhängigkeiten zwischen Add-ons erkennen.

Beispiele:

- Keycloak kann optional PostgreSQL benötigen.
- OpenTelemetry kann Beispielkonfigurationen für bestehende App-Services erzeugen.

Verhalten bei erkannter Abhängigkeit:

- Im interaktiven Modus muss das Produkt nachfragen, ob das fehlende Add-on automatisch hinzugefügt werden soll.
- Im nicht-interaktiven Modus muss das Produkt mit einem Fehler abbrechen und auf die fehlende Abhängigkeit hinweisen.
- Über die Option `--with-deps` muss das Produkt fehlende Abhängigkeiten automatisch hinzufügen.

---

### LH-FA-ADD-007 – Service entfernen

Priorität: V1

Das Produkt muss einen Service wieder entfernen können.

Beispiel:

```bash
u-boot remove postgres
```

Mindestumfang:

- Service-Eintrag in `compose.yaml` entfernen
- zugehörige verwaltete Blöcke (z. B. in `.env.example`) entfernen
- Eintrag in `u-boot.yaml` auf `enabled: false` setzen
- Volumes nur auf explizite Anforderung (`--purge`) löschen
- Abhängigkeiten anderer Services prüfen und vor dem Entfernen warnen

---

## 4.6 Starten und Stoppen der Umgebung

### LH-FA-UP-001 – Umgebung starten

Priorität: MVP

Das Produkt muss die Entwicklungsumgebung starten können.

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

- Docker installiert (Mindestversion gemäß `LH-RISK-001`)
- Docker erreichbar
- Docker Compose verfügbar (Mindestversion)
- Git verfügbar
- Schreibrechte im Projektverzeichnis
- gültige `compose.yaml`, falls vorhanden
- gültige `u-boot.yaml`, falls vorhanden

---

### LH-FA-DIAG-003 – Fehlerklassifikation

Priorität: MVP

Die Diagnosefunktion muss Probleme nach Schweregrad klassifizieren.

Mögliche Stufen:

- `OK`
- `WARN`
- `ERROR`

Die Diagnosefunktion muss den Exit Code an die höchste festgestellte Stufe binden:

- nur `OK` → Exit Code `0`
- mindestens `WARN`, kein `ERROR` → Exit Code `0`
- mindestens ein `ERROR` → Exit Code ungleich `0`

---

### LH-FA-DIAG-004 – Reparaturhinweise

Priorität: MVP

Die Diagnosefunktion muss bei Problemen konkrete Reparaturhinweise ausgeben.

Beispiel:

```text
ERROR Docker daemon is not reachable.
Hint: Start Docker or check your user permissions for /var/run/docker.sock.
```

---

## 4.8 Generatoren

### LH-FA-GEN-001 – Generate-Befehl

Priorität: MVP

Das Produkt muss Generatoren über folgenden Befehl anbieten:

```bash
u-boot generate <artifact>
```

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

Das Produkt soll eine eigene Projektkonfigurationsdatei verwenden.

Beispiel:

```text
u-boot.yaml
```

---

### LH-FA-CONF-002 – Inhalt der Konfiguration

Priorität: MVP

Die Konfigurationsdatei soll mindestens enthalten:

```yaml
project:
  name: my-service
  template: micronaut-sveltekit

services:
  postgres:
    enabled: true
  keycloak:
    enabled: true
  otel:
    enabled: true

devcontainer:
  enabled: true
```

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

Priorität: V1

Das Produkt muss einen Befehl zum Anzeigen und Ändern der Konfiguration bereitstellen.

Beispiele:

```bash
u-boot config                       # gesamte Konfiguration anzeigen
u-boot config get project.name      # einzelnen Wert anzeigen
u-boot config set project.name foo  # Wert setzen
```

Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden.

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

Das Produkt darf keine fremden Skripte aus dem Internet ohne ausdrückliche Zustimmung ausführen.

---

## 5.6 Performance

### LH-NFA-PERF-001 – Schnelle CLI-Antwort

Priorität: MVP

Einfache Befehle müssen auf einem typischen Entwicklungsrechner innerhalb folgender Zeiten reagieren (gemessen ohne Docker-Kommunikation, Kaltstart):

- `u-boot --help`, `u-boot --version` – unter 200 ms
- `u-boot config get …` – unter 300 ms
- `u-boot doctor` (ohne Netz-Wartezeit) – unter 2 s

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
u-boot <command> [subcommand] [options]
```

---

### LH-SA-CLI-002 – Vorgesehene Befehle

Priorität: MVP

| Befehl                       | Zweck                              |
| ---------------------------- | ---------------------------------- |
| `u-boot init`                | Projekt initialisieren             |
| `u-boot add <service>`       | Service hinzufügen                 |
| `u-boot remove <service>`    | Service entfernen                  |
| `u-boot up`                  | Umgebung starten                   |
| `u-boot down`                | Umgebung stoppen                   |
| `u-boot doctor`              | Umgebung prüfen                    |
| `u-boot logs`                | Logs anzeigen                      |
| `u-boot generate <artifact>` | Artefakt erzeugen                  |
| `u-boot config`              | Konfiguration anzeigen oder ändern |
| `u-boot template list`       | Templates anzeigen                 |

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

---

### LH-SA-FILE-002 – Markierte verwaltete Bereiche

Priorität: MVP

Automatisch verwaltete Bereiche in Dateien sollen markiert werden.

Beispiel:

```yaml
# BEGIN U-BOOT MANAGED BLOCK: postgres
# ...
# END U-BOOT MANAGED BLOCK: postgres
```

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
- Vor der Migration muss eine Sicherungsdatei (`u-boot.yaml.bak`) erzeugt werden.
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

### LH-QA-003 – CI-Fähigkeit

Priorität: MVP

Das Projekt soll in einer CI-Umgebung testbar sein.

---

### LH-QA-004 – Linting

Priorität: V1

Das Projekt soll Linting für Quellcode und Konfigurationsdateien unterstützen.

---

## 9. Akzeptanzkriterien

### LH-AK-001 – Minimaler Init-Flow

Priorität: MVP

Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
mkdir demo
cd demo
u-boot init
u-boot doctor
```

Erwartetes Ergebnis:

- Projektstruktur wurde erzeugt
- `u-boot doctor` meldet keine kritischen Fehler
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
- Admin-Zugangsdaten werden über `.env.example` dokumentiert (z. B. `KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD`)
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
- Projekt kann in VS Code im Devcontainer geöffnet werden

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

- Mindestversionen dokumentieren
- `u-boot doctor` prüft Versionen

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
- `u-boot up`
- `u-boot down`
- Erzeugung von `compose.yaml`
- Erzeugung von `.env.example`
- Erzeugung von `README.md`
- Erzeugung von `CHANGELOG.md`
- grundlegende Devcontainer-Unterstützung

---

### LH-MVP-002 – Kann nach dem MVP folgen

Nach dem MVP können ergänzt werden:

- `u-boot add keycloak`
- `u-boot add otel`
- Template-System
- JSON-Ausgabe
- lokale Custom-Templates
- Diff-Vorschau
- Plugin-System

---

## 13. Traceability-Matrix

| Lastenheft-Kennung | Kurzbeschreibung               | Priorität | Spätere Ableitung im Pflichtenheft | Testfall     |
| ------------------ | ------------------------------ | --------- | ---------------------------------- | ------------ |
| LH-FA-CLI-002      | Hilfeausgabe                   | MVP       | PH-CLI-002                         | TC-CLI-002   |
| LH-FA-CLI-004      | Fehlerausgabe                  | MVP       | PH-CLI-004                         | TC-CLI-004   |
| LH-FA-CLI-005      | Verbosity und Logging          | MVP       | PH-CLI-005                         | TC-CLI-005   |
| LH-FA-CLI-006      | Exit Codes                     | MVP       | PH-CLI-006                         | TC-CLI-006   |
| LH-FA-CLI-007      | Dry Run                        | V1        | PH-CLI-007                         | TC-CLI-007   |
| LH-FA-INIT-001     | Projekt initialisieren         | MVP       | PH-INIT-001                        | TC-INIT-001  |
| LH-FA-INIT-003     | Projektstruktur erzeugen       | MVP       | PH-INIT-003                        | TC-INIT-003  |
| LH-FA-INIT-005     | Überschreibschutz              | MVP       | PH-INIT-005                        | TC-INIT-005  |
| LH-FA-INIT-006     | Projektnamen-Validierung       | MVP       | PH-INIT-006                        | TC-INIT-006  |
| LH-FA-INIT-007     | Git-Initialisierung            | MVP       | PH-INIT-007                        | TC-INIT-007  |
| LH-FA-DEV-001      | Devcontainer erzeugen          | MVP       | PH-DEV-001                         | TC-DEV-001   |
| LH-FA-DEV-005      | forwardPorts                   | MVP       | PH-DEV-005                         | TC-DEV-005   |
| LH-FA-DOC-001      | Compose-Datei erzeugen         | MVP       | PH-DOC-001                         | TC-DOC-001   |
| LH-FA-DOC-004      | Volumes                        | MVP       | PH-DOC-004                         | TC-DOC-004   |
| LH-FA-ADD-002      | PostgreSQL hinzufügen          | MVP       | PH-ADD-002                         | TC-ADD-002   |
| LH-FA-ADD-003      | Keycloak hinzufügen            | V1        | PH-ADD-003                         | TC-ADD-003   |
| LH-FA-ADD-004      | OpenTelemetry hinzufügen       | V1        | PH-ADD-004                         | TC-ADD-004   |
| LH-FA-ADD-005      | Mehrfaches Hinzufügen          | MVP       | PH-ADD-005                         | TC-ADD-005   |
| LH-FA-ADD-006      | Add-on-Abhängigkeiten          | V1        | PH-ADD-006                         | TC-ADD-006   |
| LH-FA-ADD-007      | Service entfernen              | V1        | PH-ADD-007                         | TC-ADD-007   |
| LH-FA-UP-001       | Umgebung starten               | MVP       | PH-UP-001                          | TC-UP-001    |
| LH-FA-UP-003       | Startstatus anzeigen           | MVP       | PH-UP-003                          | TC-UP-003    |
| LH-FA-UP-004       | Umgebung stoppen               | MVP       | PH-UP-004                          | TC-UP-004    |
| LH-FA-UP-005       | Logs anzeigen                  | V1        | PH-UP-005                          | TC-UP-005    |
| LH-FA-DIAG-001     | Doctor-Befehl                  | MVP       | PH-DIAG-001                        | TC-DIAG-001  |
| LH-FA-DIAG-002     | Voraussetzungen prüfen         | MVP       | PH-DIAG-002                        | TC-DIAG-002  |
| LH-FA-DIAG-003     | Exit Code nach Diagnose        | MVP       | PH-DIAG-003                        | TC-DIAG-003  |
| LH-FA-GEN-002      | Changelog erzeugen             | MVP       | PH-GEN-002                         | TC-GEN-002   |
| LH-FA-GEN-003      | README erzeugen                | MVP       | PH-GEN-003                         | TC-GEN-003   |
| LH-FA-GEN-005      | Idempotenz                     | MVP       | PH-GEN-005                         | TC-GEN-005   |
| LH-FA-CONF-003     | Konfiguration lesen            | MVP       | PH-CONF-003                        | TC-CONF-003  |
| LH-FA-CONF-004     | Konfiguration aktualisieren    | MVP       | PH-CONF-004                        | TC-CONF-004  |
| LH-FA-CONF-005     | Konfiguration ändern           | V1        | PH-CONF-005                        | TC-CONF-005  |
| LH-FA-TPL-004      | Templates auflisten            | V1        | PH-TPL-004                         | TC-TPL-004   |
| LH-DA-003          | Schema-Version                 | MVP       | PH-DA-003                          | TC-DA-003    |
| LH-DA-004          | Schema-Migration               | Later     | PH-DA-004                          | TC-DA-004    |
| LH-NFA-REL-001     | Kein stilles Überschreiben     | MVP       | PH-REL-001                         | TC-REL-001   |
| LH-NFA-REL-002     | Wiederholbare Ausführung       | MVP       | PH-REL-002                         | TC-REL-002   |
| LH-NFA-PERF-001    | Antwortzeiten                  | MVP       | PH-PERF-001                        | TC-PERF-001  |
| LH-NFA-SEC-001     | Keine Secrets einchecken       | MVP       | PH-SEC-001                         | TC-SEC-001   |
| LH-NFA-SEC-004     | Keine fremde Skriptausführung  | MVP       | PH-SEC-004                         | TC-SEC-004   |

---

## 14. Offene Punkte

### LH-OPEN-001 – Implementierungssprache

Die Implementierungssprache ist noch festzulegen.

Mögliche Optionen:

- Go
- Rust
- Python
- TypeScript/Node.js

---

### LH-OPEN-002 – Paketierung

Die spätere Verteilung ist noch festzulegen.

Mögliche Optionen:

- einzelnes Binary
- Container Image
- npm package
- pip package
- Homebrew
- Debian/RPM-Paket

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


