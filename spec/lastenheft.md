# Lastenheft: `u-boot` – Projekt-Bootstrapper für Docker-/Devcontainer-Stacks

| Dokument         | Lastenheft                                                         |
| ---------------- | ------------------------------------------------------------------ |
| Projektname      | `u-boot`                                                           |
| Kurzbeschreibung | CLI-Tool zum Bootstrapping reproduzierbarer Entwicklungsumgebungen |
| Zielplattform    | Linux, Docker, VS Code Dev Containers                              |
| Hauptnutzer      | Softwareentwickler, DevOps-Engineers, technische Teams             |
| Version          | 0.14.0                                                             |
| Status           | Accepted                                                           |
| Datum            | 2026-05-21 (Erstfassung; Änderungen siehe §16 Historie)            |

---

## 0. Lesehinweise

### LH-LESE-001 – Modalverben

In diesem Dokument haben Modalverben folgende Bedeutung (in Anlehnung an RFC 2119):

- **muss** – verbindliche Anforderung (Pflicht).
- **soll** – Empfehlung; Abweichungen müssen begründet werden.
- **kann** / **darf** – optionale Eigenschaft oder ausdrückliche Erlaubnis.

### LH-LESE-002 – Sprache

Die Spezifikation ist auf Deutsch verfasst.

CLI-Ausgaben, Fehlermeldungen und erzeugte Dateien (Kommentare, Beispielwerte, README-Vorlagen) sind auf Englisch.

### LH-LESE-003 – Dokumentenordnung

Dieses Lastenheft ist der Vertrag: Es legt fest, **was** das Produkt leistet. Technische Festlegungen (Schemata, Algorithmen, Defaults, Fehler-Codes, Build- und Dokumentationsdetails) stehen in einem eigenen technischen Dokument, das dieses Lastenheft präzisiert, aber nie erweitert; bei Widerspruch gilt das Lastenheft. Das Lastenheft verweist nicht auf dieses technische Dokument.

Anforderungen an den Entwicklungsprozess der u-boot-Codebase (Bau, Dokumentation, Architektur-Vorgaben, Qualitätssicherung) gehören nicht in dieses Lastenheft; sie stehen im adoptierten Baseline-Regelwerk.

### LH-LESE-004 – Prüfbarkeit

Jede funktionale Anforderung soll durch mindestens einen Akzeptanztest überprüfbar sein.

---

## 1. Zielbestimmung

### LH-ZB-001 – Projektziel

`u-boot` soll ein CLI-Tool werden, das vollständige Entwicklungsumgebungen für Docker-basierte Softwareprojekte erzeugt, erweitert, prüft und startet.

Das Tool soll insbesondere Projektstrukturen, Docker-Konfigurationen, Devcontainer-Setups, optionale Infrastrukturservices und wiederkehrende Entwicklungsartefakte automatisch bereitstellen.

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

### LH-ZB-003 – Repo-Beschreibung

```text
u-boot: A developer environment bootloader for Docker-based projects.
```

---

## 2. Produkteinsatz

### LH-PE-001 – Anwendungsbereich

`u-boot` soll für Softwareprojekte eingesetzt werden, die lokal oder in Devcontainern entwickelt werden und Docker beziehungsweise Docker Compose als zentrale Laufzeitumgebung verwenden.

### LH-PE-002 – Zielgruppen

Das Produkt richtet sich an:

- Softwareentwickler
- DevOps-Engineers
- technische Projektleiter
- Entwicklerteams mit Docker-basierten Entwicklungsumgebungen
- Teams, die reproduzierbare lokale Setups benötigen

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

*Priorität: MVP.* Das Produkt muss als Kommandozeilenprogramm mit dem Namen `u-boot` aufrufbar sein.

### LH-FA-CLI-002 – Hilfeausgabe

*Priorität: MVP.* Das Produkt muss eine Hilfeausgabe bereitstellen.

Die Hilfeausgabe muss mindestens enthalten:

- verfügbare Befehle
- Kurzbeschreibung je Befehl
- Optionen je Befehl
- Beispiele

### LH-FA-CLI-003 – Versionsausgabe

*Priorität: MVP.* Das Produkt muss die installierte Version ausgeben können (`u-boot --version`).

### LH-FA-CLI-004 – Fehlerausgabe

*Priorität: MVP.* Das Produkt muss verständliche Fehlermeldungen ausgeben.

Fehlermeldungen müssen enthalten:

- Ursache
- betroffener Befehl oder betroffene Datei
- empfohlene Korrekturmaßnahme

### LH-FA-CLI-005 – Verbosity und Logging

*Priorität: MVP.* Das Produkt muss ein konfigurierbares Ausgabeverbose-Level (Verbosity) unterstützen.

Mindestens müssen folgende Stufen unterstützt werden:

- `--quiet` – nur Fehler
- Standard – Statusmeldungen und Fehler
- `--verbose` – zusätzlich Detailinformationen
- `--debug` – zusätzlich interne Diagnoseausgaben

Unter `--quiet` geben die Subkommandos, die etwas ändern (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`), nur Fehler aus; ihre Erfolgs- und Zusammenfassungsmeldungen entfallen. `u-boot doctor` zeigt nur Fehler-Einträge. Vier Ausnahmen bleiben unter `--quiet` sichtbar:

1. die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) und eines im Container aktivierten Login-Dienstes (Passwort-Login) in der Befehlsausgabe ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion)) sowie die Warn-Einträge von `u-boot doctor`, die eine Sicherheitseigenschaft des Sandbox-Profils betreffen (Quelle der Zugangsdaten nach [`LH-FA-DEV-009`](#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), Egress-Restriktion nach [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion), Sicherheitslockerungen des Docker-Daemons nach [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); dazu gehört die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift; die Meldung, dass die Container-Runtime ohne Sandbox-Profil keine Wirkung hat, ist eine Konfigurationsmeldung und gehört nicht dazu;
2. die Ausweisung der Degradations-Fallbacks (Container ohne Docker-Daemon, `vfs`-Storage, Wegfall der Egress-Restriktion) in der Befehlsausgabe und in `u-boot doctor` ([`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer));
3. der Fortschritt einzelner Services bei `u-boot up` (Pull, Create, Start, Healthcheck; [`LH-NFA-PERF-002`](#lh-nfa-perf-002--startzeit-abhängig-von-docker));
4. die Zusammenfassung der betroffenen Pfade, die `--force` vor dem Überschreiben ausgibt ([`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-INIT-005`](#lh-fa-init-005--überschreibschutz)).

Alle anderen Hinweise und Warnungen entfallen unter `--quiet`. Die Ausgabe der Abfragebefehle, die nichts verändern (`config get`, `config list`, `egress list`, `logs`, `template list`), ist das Ergebnis des Befehls und bleibt. Ebenso ist die Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` ([`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](#lh-fa-cli-008--diff-ausgabe)) — der geplante Änderungsplan und die Diff-Darstellung — das Ergebnis des Aufrufs und bleibt. Die Erfolgs- und Zusammenfassungsmeldung, die ein tatsächliches Schreiben abschließt, ist keine Vorschau und entfällt auch bei `--diff` ohne `--dry-run`.

Werden mehrere Verbosity-Optionen gleichzeitig angegeben (z. B. `--quiet --verbose`), gewinnt die zuletzt auf der Kommandozeile angegebene Option. Eine Validierungsabweisung wegen Mehrfachangabe erfolgt nicht.

### LH-FA-CLI-005A – Interaktivität und Automatisierung

*Priorität: MVP.* Das Produkt muss nicht-interaktive Ausführung unterstützen.

Es muss mindestens folgende Optionen bieten:

- `--yes` – Standardfragen automatisch bejahen (für CI/Skripte).
- `--no-interactive` – keine Rückfragen stellen; erforderliche Bestätigungen mit klarem Fehler abbrechen.
- `--yes` und `--no-interactive` sind exklusiv. Bei gleichzeitiger Nutzung ist ein CLI-Fehler mit Exit-Code `2` ([`LH-FA-CLI-006`](#lh-fa-cli-006--exit-codes)) zu erzeugen.
- Für deterministisches Verhalten in Skripten und CI sind beide Modi einzeln nutzbar.
- Die Optionen sind auf Befehle anzuwenden, die Bestätigungsentscheidungen benötigen (insb. `u-boot init`, `u-boot add`, `u-boot remove`, `u-boot config set`, `u-boot down --volumes`).
- Für `u-boot init` gilt zusätzlich `--assume-existing` (nur für diesen Befehl): Ohne dieses Flag wird eine implizite Erkennung als bestehendes Projekt im nicht-interaktiven Modus abgelehnt (fachlicher Fehler, Exit-Code `10`); `--yes` genügt dafür nicht.
- Bei aktivierter Nicht-Interaktivität entsteht keine neue Rückfrage: `--no-interactive` bricht bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab, `--yes` führt die vorgesehene Standardentscheidung deterministisch aus.

Bei destruktiven Operationen (insb. `u-boot down --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad (`--yes` oder aktiv bestätigten interaktiven Pfad) erfolgen. Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.

`--force` und `--backup` sind in nicht-interaktiven Läufen zulässig und arbeiten deterministisch; `--force` erzeugt keine zusätzlichen Rückfragen und gibt immer eine Zusammenfassung der betroffenen Pfade aus. Ist eine sichere automatische Abarbeitung nicht möglich (z. B. fehlender verwalteter Block ohne `--backup`), bricht der Befehl mit Exit-Code `10` ab.

### LH-FA-CLI-006 – Exit Codes

*Priorität: MVP.* Das Produkt muss aussagekräftige Exit Codes liefern.

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

Für alle fachlichen Fehler ist die Verwendung von `10`, `11` oder `12` bindend. Typische Zuordnung (Empfehlung): `10` bei Struktur-, Namens- und Konfigurationsvalidierung, `11` bei Umgebungsproblemen, `12` bei Laufzeitfehlern von Docker/Compose. Nicht-fachliche Fehler dürfen mit `1` codiert werden, `13` bis `15` nur mit dokumentierter Bedeutung; `16` bis `19` sind nicht zu verwenden.

### LH-FA-CLI-007 – Dry Run

*Priorität: V1.* Das Produkt muss für dateiverändernde Befehle einen Dry-Run-Modus unterstützen.

Der Dry-Run muss anzeigen, welche Dateien erzeugt, geändert oder gelöscht würden, ohne Änderungen am Dateisystem vorzunehmen.

`--dry-run` darf mit `--diff` kombiniert werden. Bei gleichzeitiger Nutzung darf das Tool keine Datei schreiben; die Ausgabe besteht ausschließlich aus dem geplanten Änderungsplan und der Diff-Darstellung.

Bei gleichzeitiger Verwendung von `--dry-run` und `--json` muss die Ausgabe streng maschinenlesbar (`JSON`) erfolgen und keine unstrukturierten Text-UI-Zeilen enthalten.

Für `--dry-run --json` ist die Ausgabe als maschinenlesbares JSON mit mindestens den Pflichtfeldern `status`, `command`, `dryRun`, `diff`, `plannedFiles`, `changes`, `diagnostics` und `exitCode` zu liefern; `command` ist einer der Werte `init`, `add`, `remove`, `up`, `down`, `doctor`, `logs`, `generate`, `config` oder `template`. Jeder Eintrag in `plannedFiles` trägt `path` und `action` (`create`, `modify` oder `delete`), jeder Eintrag in `changes` trägt `path` und `count`, jeder Eintrag in `diagnostics` trägt `level` (`warn` oder `error`), `code` und `message` (optional `file`).

Bei gruppierten Befehlen wie `command == "template"` oder `command == "config"` muss das Feld `subcommand` gesetzt sein (z. B. `list`, `get`, `set`).

Weitere Felder sind erlaubt.

Für `diagnostics[*].code` und die Kopplung von `status` an den höchsten `level` gilt die Regel aus [`LH-NFA-USE-004`](#lh-nfa-use-004--maschinenlesbare-ausgabe).

### LH-FA-CLI-008 – Diff-Ausgabe

*Priorität: V1.* Das Produkt soll bei dateiverändernden Befehlen eine Diff-Ausgabe unterstützen.

Die Diff-Ausgabe muss Unterschiede zwischen aktuellem und geplantem Zustand der betroffenen Dateien zeigen.

Wird `--diff` ohne `--dry-run` gesetzt, zeigt sie den geplanten Endzustand als Vorschau.
Wird `--diff` mit `--dry-run` kombiniert, gilt dieselbe Vorschau bei vollständigem Schreibschutz.

Bei Kombination von `--diff` mit `--json` ist die komplette JSON-Struktur inkl. Pflichtfeldern aus dem in [`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run) definierten Schema auszugeben. Die Felder `dryRun` und `diff` sind dabei korrekt auf den konkreten Ausführungsmodus gesetzt (`dryRun` je nach Aufruf, `diff` immer `true`).

Beispiel: `u-boot add postgres --diff --json` ohne `--dry-run` (Vorschau mit anschließendem Schreiben) liefert `dryRun: false` und `diff: true`.

Für reine Vorschau-Workflows gelten die selben Exit-Codes wie bei der Nicht-Diff-Ausgabe.

---

## 4.2 Projektinitialisierung

### LH-FA-INIT-001 – Neues Projekt initialisieren

*Priorität: MVP.* Das Produkt muss ein neues Projekt mit `u-boot init` initialisieren können (in einem bestehenden Verzeichnis zusätzlich mit `--assume-existing`).

### LH-FA-INIT-002 – Projektname

*Priorität: MVP.* Das Produkt muss bei der Initialisierung einen Projektnamen verwenden können.

Wird kein Name explizit angegeben, verwendet das Tool standardmäßig den aktuellen Verzeichnisnamen als Basis.

Der abgeleitete Name wird deterministisch normalisiert:

Kleinbuchstaben, Zeichen außerhalb von `a-z`, `0-9` und `-` werden zu `-`, aufeinanderfolgende `-` zusammengeführt, führende und nachgestellte `-` entfernt, Länge auf 1 bis 63 Zeichen begrenzt; anschließend gilt die Validierung aus [`LH-FA-INIT-006`](#lh-fa-init-006--projektnamen-validierung).

Ist kein gültiger Name ableitbar oder angegeben, muss der Befehl mit einer klaren Fehlermeldung abbrechen und auf die explizite Übergabe eines Namens (`u-boot init <name>`) verweisen.

### LH-FA-INIT-003 – Projektstruktur erzeugen

*Priorität: MVP.* Das Produkt muss eine grundlegende Projektstruktur erzeugen.

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

Bei aktivierter Devcontainer-Unterstützung (siehe [`LH-FA-DEV-001`](#lh-fa-dev-001--devcontainer-erzeugen)) zusätzlich:

```text
.devcontainer/devcontainer.json
.devcontainer/Dockerfile
```

### LH-FA-INIT-004 – Bestehendes Projekt erkennen

*Priorität: MVP.* Das Produkt muss erkennen, ob es in einem bestehenden Projektverzeichnis ausgeführt wird.

Relevante Dateien sind die Projektsteuerdateien (`u-boot.yaml`, `compose.yaml`, `.env.example`) und die Elemente des Mindestumfangs der Projektstruktur ([`LH-FA-INIT-003`](#lh-fa-init-003--projektstruktur-erzeugen)).

Wenn mindestens eine der Projektsteuerdateien (`u-boot.yaml`, `compose.yaml`, `.env.example`) vorhanden ist, ist das Verzeichnis als bestehendes Projekt zu behandeln.
Liegt keine Projektsteuerdatei vor, gilt das Verzeichnis nur als wahrscheinliches bestehendes Projekt, wenn mindestens drei Elemente aus dem Mindestumfang der Projektstruktur bereits vorhanden sind.
In diesem Fall muss `u-boot init` im interaktiven Modus explizit nachfragen, ob das Verzeichnis als bestehendes Projekt behandelt werden soll.
`--assume-existing` ist nur für `u-boot init` gültig.
Das genaue Verhalten im nicht-interaktiven Modus (mit/ohne `--assume-existing`, Exit-Code-Vergabe) ist verbindlich in [`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung) definiert; diese Anforderung wiederholt es nicht.
Bestehende Dateien dürfen auch bei impliziter oder expliziter Annahme als bestehendes Projekt nicht kommentarlos überschrieben werden; es gilt der Überschreibschutz aus [`LH-FA-INIT-005`](#lh-fa-init-005--überschreibschutz).

### LH-FA-INIT-005 – Überschreibschutz

*Priorität: MVP.* Das Produkt muss vor dem Überschreiben bestehender Dateien schützen.

Standardverhalten ohne Option:

- Abbruch mit Hinweis, welche Datei kollidiert

Zusätzliche Strategien über Option:

- `--backup` – bestehende Datei (bei Verzeichnissen der gesamte Baum) als `<name>.bak` sichern und ersetzen, ohne vorhandene Backups zu überschreiben (numerische Suffixe); bei Fehlern erfolgt ein Rollback.
- `--force` – bestehende Dateien ohne Rückfrage überschreiben; vor dem Schreiben muss eine Zusammenfassung der betroffenen Pfade ausgegeben werden

Für strukturierte Konfigurationsdateien bleiben nicht verwaltete Inhalte bei `--force` erhalten; nur ein erkannter `U-BOOT MANAGED BLOCK` wird verändert ([`LH-SA-FILE-002`](#lh-sa-file-002--markierte-verwaltete-bereiche)). Fehlt der verwaltete Block, wird mit `--backup` vor dem vollständigen Überschreiben der gesamte Inhalt gesichert, ohne `--backup` mit Exit-Code `10` abgebrochen.

Für `--force`, `--backup` und nicht-interaktive Modi gilt zusätzlich die in [`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung) definierte Entscheidungslogik für Bestätigungen.

### LH-FA-INIT-006 – Projektnamen-Validierung

*Priorität: MVP.* Das Produkt muss den Projektnamen validieren.

Die Validierung gilt für den explizit übergebenen und den automatisch aus dem Arbeitsverzeichnis abgeleiteten Projektnamen.

Regeln:

- erlaubt sind Kleinbuchstaben, Ziffern und Bindestrich
- beginnt mit einem Kleinbuchstaben
- endet mit einem Kleinbuchstaben oder einer Ziffer (entfällt bei einstelligen Namen)
- minimale Länge: 1 Zeichen
- maximale Länge: 63 Zeichen

Ungültige Namen müssen mit einer klaren Fehlermeldung abgelehnt werden.

### LH-FA-INIT-007 – Git-Repository-Initialisierung

*Priorität: MVP.* Das Produkt muss Git-Initialisierung als Teil von `u-boot init` unterstützen.

Verhalten:

- Standardverhalten: aktiviert – ein neues Git-Repository wird angelegt, sofern noch keines vorhanden ist.
- Abschaltbar über `--no-git`.
- Ein bereits vorhandenes Repository darf nicht erneut initialisiert werden.

---

## 4.3 Devcontainer-Unterstützung

### LH-FA-DEV-001 – Devcontainer erzeugen

*Priorität: MVP.* Das Produkt muss eine Devcontainer-Konfiguration erzeugen können.

Die Erzeugung muss sowohl bei `u-boot init` über eine Option als auch nachträglich auslösbar sein (`u-boot init --devcontainer` bzw. `u-boot generate devcontainer`).

Mindestdateien:

```text
.devcontainer/devcontainer.json
.devcontainer/Dockerfile
```

### LH-FA-DEV-002 – VS-Code-Kompatibilität

*Priorität: MVP.* Die erzeugte Devcontainer-Konfiguration muss mit VS Code Dev Containers kompatibel sein.

### LH-FA-DEV-003 – Devcontainer-Features

*Priorität: V1.* Das Produkt soll optionale Devcontainer-Features unterstützen.

Beispiele:

- Git
- Docker CLI
- Node.js
- Java
- SDKMAN
- PostgreSQL Client
- Kubernetes Tools

Für optionale externe Feature-Quellen gilt:

- Keine fremden Skripte dürfen ohne Zustimmung ausgeführt werden ([`LH-NFA-SEC-004`](#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)).
- Standardmäßig sind nur lokal hinterlegte oder ausdrücklich freigegebene Features erlaubt.
- Die Freigabe erfolgt als klarer, protokollierter Schritt: interaktiv durch Bestätigung, nicht-interaktiv nur über `--allow-external-feature-sources <quelle>[,<quelle>...]` (gültig nur für `u-boot init --devcontainer`, `u-boot generate devcontainer` und `u-boot config set devcontainer.featureSources.allow`).
- Die zugelassenen Quellen werden als explizit freigegebene Liste in der Projektkonfiguration gespeichert.
- Ohne explizit erlaubte Quelle führt der Versuch, externe Quellen zu nutzen, zu einem fachlichen Fehler (`code [LH-FA-DEV-003](#lh-fa-dev-003--devcontainer-features)`, Exit-Code `10`).
- `--yes` allein gilt nicht als Zustimmung für externe Quellen.

### LH-FA-DEV-004 – Benutzerrechte

*Priorität: MVP.* Der Devcontainer soll standardmäßig mit einem nicht-root Benutzer arbeiten.

Die UID dieses Benutzers muss an den Host anpassbar sein (z. B. `501` unter macOS mit Colima):

- Konfigurationsschlüssel `devcontainer.user.uid` (Ganzzahl, optional, Default `1000`).
- Zulässig sind Ganzzahlen von `1` bis `65535`. `0` (root), negative oder nicht numerische Werte führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).
- Ohne `devcontainer.user.uid` bleibt das erzeugte Ergebnis unverändert (Default `1000`).

### LH-FA-DEV-005 – Ports

*Priorität: MVP.* Das Produkt muss Ports aus aktivierten Services in der Devcontainer-Konfiguration berücksichtigen.

Konkret müssen die Ports der Services in `devcontainer.json` als `forwardPorts` eingetragen werden.
Ist keine aktive Port-Exposition in der aktuellen Projektkonfiguration vorhanden, darf `forwardPorts` fehlen.

### LH-FA-DEV-006 – Sandbox-Profil

*Priorität: V1.* Das Produkt soll ein opt-in Sandbox-Profil für Devcontainer erzeugen können, das den Einsatz autonomer Agenten (ohne Rückfrage an den Menschen) im Container auf Schadensbegrenzung auslegt. Das Profil ist eine Schadensbegrenzung und keine harte Isolationsgrenze.

Aktivierung über das Flag `--sandbox` (bei `u-boot init --devcontainer` und `u-boot generate devcontainer`)
oder über die Projektkonfiguration `devcontainer.profile: sandbox` (Werte: `default` | `sandbox`, Default `default`). `--sandbox` setzt den Konfigurationsschlüssel; `--sandbox` ohne aktivierbaren Devcontainer führt zu einem fachlichen Fehler (Exit-Code `10`).

Das erzeugte Ergebnis im Sandbox-Profil muss:

- einen nicht-root Benutzer verwenden ([`LH-FA-DEV-004`](#lh-fa-dev-004--benutzerrechte));
- das Host-Arbeitsverzeichnis nicht per Bind-Mount einbinden: Der Workspace liegt in einem benannten Volume, das Repository wird im Container geklont (Quelle: `devcontainer.sandbox.repository`, sonst `origin`; URLs mit Zugangsdaten oder unzulässigen Zeichen sind ein fachlicher Fehler, Exit-Code `10`; ohne Quelle entfällt der Clone-Schritt mit Warnung);
- keinen Container-Runtime-Socket des Hosts einbinden;
- kein `--privileged` und keine zusätzlichen Capabilities setzen, sofern nicht durch [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) oder [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion) ausdrücklich verlangt;
- keine Host-Dateien mit Geheimnissen (z. B. `~/.ssh`, `~/.aws`, Credential-Stores) einbinden;
- jede Sicherheitslockerung (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) in der Befehlsausgabe ausweisen;
- einen im Container aktivierten Login-Dienst (Passwort-Login) in der Befehlsausgabe ausweisen.

Nicht-Ziele: `u-boot` startet keinen Agenten und setzt keinen Berechtigungsmodus des Agenten.

Für Devcontainer-Features oder externe Skripte, die das Profil einbindet, gelten [`LH-FA-DEV-003`](#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) unverändert; das Profil gibt keine externe Quelle implizit frei.

### LH-FA-DEV-007 – Container-Runtime im Sandbox-Devcontainer

*Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) soll optional eine Container-Runtime im Container bereitgestellt werden, damit Image-Builds und Container-Läufe ohne Zugriff auf einen Host-Socket möglich sind.

- Konfigurationsschlüssel `devcontainer.sandbox.nestedRuntime` (`podman` | `docker` | `none`, Default `none`).
- Bei `podman`: rootless Podman mit `docker`-Kompatibilität im Container. `podman build` — und über den `docker`-Alias `docker build` — baut im Container ein Dockerfile (`FROM`, `RUN`, `COPY`, voll qualifizierte Image-Namen) ohne Zugriff auf einen Host-Socket, auch hinter der Egress-Restriktion ([`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion)) mit ihrer Default-Allowlist (gemessen mit `docker.io/library/alpine:3` und `debian:stable-slim`); das gilt auf den gemessenen Hosts (Stand der Messung: Docker 29.5.2 auf Colima 0.10.3, arm64, VM Ubuntu 24.04 mit `kernel.apparmor_restrict_unprivileged_userns=0`; andere Architekturen und Engines nicht gemessen). BuildKit-Erweiterungen (`docker buildx`, `# syntax=`-Frontends) und `docker compose` sind nicht zugesagt.
- Bei `docker`: ein eigener Docker-Daemon im Container mit BuildKit sowie Buildx- und Compose-Plugin. Die Docker-CLI im Container spricht ausschließlich mit diesem Daemon, sein Socket ist dem Container-Benutzer zugänglich. Der Daemon-Speicher liegt in einem benannten Volume je Instanz. `docker` verlangt `--privileged`; die Befehlsausgabe weist das als Lockerung aus, mit dem Hinweis, dass das Profil damit keine Isolationsgrenze gegenüber dem Container-Host mehr bietet.
- Bei `podman` startet das Ergebnis unter Docker auf Hosts, die verschachtelte User-Namespaces erlauben (gemessen: Colima, arm64, mit `kernel.apparmor_restrict_unprivileged_userns=0` in der VM als Voraussetzung) und wenn die Devices `/dev/fuse` und `/dev/net/tun` auf dem Container-Host vorhanden sind; mit dem Wert `1`, dem Standard einer Ubuntu-24.04-VM, endet das Startscript mit Exit-Code `11`, ohne eines der beiden Devices startet der Container nicht (Fehler der Container-Engine). Podman als Host-Engine ist für `podman` nicht zugesagt. Bei `docker` startet es unter Docker auf Linux und unter Colima (Linux-VM); Docker Desktop und Podman als Host-Engine sind für `docker` nicht zugesagt.

**Degradation:** `devcontainer.sandbox.onUnavailable` (`warn` | `fail`, Default `warn`) steuert den Umgang mit fehlenden Fähigkeiten: Bei `warn` führen ein im Container fehlendes `/dev/fuse` (der Container wurde ohne das Device gestartet, etwa weil die erzeugten `runArgs` geändert wurden) zu einem Fallback auf `vfs`-Storage, eine nicht gewährbare Egress-Capability zum Wegfall der Restriktion und ein Docker-Daemon, der beim Containerstart nicht betriebsbereit wird, zu einem Container ohne Daemon, jeweils mit Warnung; bei `fail` sind alle drei ein Umgebungsproblem (Exit-Code `11`). Blockierte verschachtelte User-Namespaces (Seccomp/AppArmor) bei `nestedRuntime: podman` sind in beiden Modi ein Umgebungsproblem (Exit-Code `11`). Ein fehlendes `/dev/fuse` oder `/dev/net/tun` auf dem Container-Host ist keine Degradationsstufe: die Container-Engine lehnt den Start ab; `u-boot doctor` prüft `/dev/net/tun` nicht und weist ein fehlendes `/dev/fuse` nur aus, soweit der Host es bestimmen kann, als Umgebungsproblem (Exit-Code `11`). Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

### LH-FA-DEV-008 – Egress-Restriktion

*Priorität: V2.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) soll eine Allowlist für ausgehenden Netzwerkverkehr aktivierbar sein.

- Aktivierung über `devcontainer.sandbox.egress.enabled: true` (Default `false`); keine Flag-Variante, damit `--yes`/`--no-interactive` ([`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung)) ohne Sonderfall bleiben.
- Erlaubte Ziele: `devcontainer.sandbox.egress.allow` (Liste von Hostnamen); die Default-Allowlist (gemeinsame Basis plus Ergänzungen je gewähltem Stack) ist dokumentiert.
- Die wirksame Allowlist und die Herkunft jedes Hosts (Basis, Container-Runtime, Feature, Clone-Quelle, Nutzerliste) sind mit dem Befehl `u-boot egress list` abfragbar. Der Befehl liest nur die Projektkonfiguration (`u-boot.yaml`) und die Git-Konfiguration des Projekts, braucht weder Netz noch Docker noch einen laufenden Container, ändert nichts und zeigt dieselbe Liste, die `u-boot generate devcontainer` in das Startscript schreibt. Er weist aus, ob die Restriktion für die Konfiguration aktiv ist (Sandbox-Profil und `devcontainer.sandbox.egress.enabled`); ist sie es nicht, zeigt er die Liste, die bei Aktivierung gälte, und nennt das. Ein Name in der Liste schließt seine Subdomains ein.
- Die Allowlist ist unabhängig von `devcontainer.featureSources.allow` ([`LH-FA-DEV-003`](#lh-fa-dev-003--devcontainer-features)): jene steuert erlaubte Build-Quellen, diese die Laufzeit-Ziele.
- Die Restriktion ist ein Guardrail und keine Sandbox-Grenze; Prozesse mit der nötigen Capability können sie aufheben. Die Dokumentation muss das ausdrücklich sagen.
- Ist die nötige Capability nicht gewährbar, greift die Degradation aus [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (geprüft beim Containerstart; `u-boot doctor` prüft die Konfiguration).

### LH-FA-DEV-009 – Git-Zugangsdaten im Sandbox-Devcontainer

*Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten und Zugangsdaten eines Agenten (etwa der Zugangsschlüssel eines Coding-Agenten, den ein aktiviertes Devcontainer-Feature nutzt) weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.

- Übergabe zur Laufzeit per Umgebungsvariable; Git-Zugangsdaten alternativ als schreibgeschützter Secret-Mount. Die erzeugte Konfiguration trägt nur die Referenz auf die Umgebungsvariable des Hosts, nie den Wert.
- Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote; für den Agenten nennt sie die Umgebungsvariable, aus der der Zugangsschlüssel stammt.
- `u-boot doctor` prüft die Quelle der Zugangsdaten (Git-Token-Quelle, bei aktivem Agenten-Feature die Umgebungsvariable des Agenten): fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

---

## 4.4 Docker- und Compose-Unterstützung

### LH-FA-DOC-001 – Compose-Datei erzeugen

*Priorität: MVP.* Das Produkt muss eine `compose.yaml` erzeugen können.

### LH-FA-DOC-002 – Dockerfile erzeugen

*Priorität: V1.* Das Produkt soll bei Bedarf ein Dockerfile für die Anwendungsentwicklung erzeugen können.

Das Minimum bei aktivem Devcontainer ist die Erzeugung von `.devcontainer/Dockerfile`.

Zusätzlich kann optional ein separates Anwendungs-Dockerfile erzeugt werden:

- Standardpfad: `docker/Dockerfile`
- Konfigurierbar über Template-/Add-on-Konfiguration

### LH-FA-DOC-003 – Netzwerk

*Priorität: MVP.* Das Produkt muss ein gemeinsames Docker-Netzwerk für Services definieren können.

### LH-FA-DOC-004 – Volumes

*Priorität: MVP.* Das Produkt muss für aktivierte zustandsbehaftete Dienste persistente Volumes erzeugen können.

Beispiele:

- PostgreSQL-Daten (MVP)
- Keycloak-Daten (V1)
- OpenTelemetry-Konfiguration (V1)

### LH-FA-DOC-005 – Compose-Validierung

*Priorität: V1.* Das Produkt soll erzeugte Compose-Dateien auf syntaktische Gültigkeit prüfen können.

---

## 4.5 Service-Add-ons

### LH-FA-ADD-001 – Add-on-Befehl

*Priorität: MVP.* Das Produkt muss Services mit `u-boot add <service>` hinzufügen können.

Der Befehl ist nur in einem initialisierten `u-boot`-Projekt nutzbar (`u-boot.yaml` vorhanden).  
Ist keine gültige Projektkonfiguration vorhanden, ist mit klarer Fehlermeldung und Hinweis auf `u-boot init` abzubrechen.

### LH-FA-ADD-002 – PostgreSQL hinzufügen

*Priorität: MVP.* Das Produkt muss PostgreSQL als Service hinzufügen können (`u-boot add postgres`).

Mindestumfang:

- PostgreSQL-Service in `compose.yaml`
- Volume für Daten
- `.env.example`-Einträge
- Port-Konfiguration
- Healthcheck

### LH-FA-ADD-003 – Keycloak hinzufügen

*Priorität: V1.* Das Produkt muss Keycloak als Service hinzufügen können (`u-boot add keycloak`).

Mindestumfang:

- Keycloak-Service in `compose.yaml`
- Admin-Benutzer über `.env.example` als eindeutig markierte Beispielwerte (niemals reale Secrets), z. B.:
  - `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN`
  - `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`
- Port-Konfiguration
- optionale PostgreSQL-Anbindung bei konfigurierter persistenter externer Datenbank
- Healthcheck, soweit technisch sinnvoll

### LH-FA-ADD-004 – OpenTelemetry hinzufügen

*Priorität: V1.* Das Produkt muss OpenTelemetry-Komponenten hinzufügen können (`u-boot add otel`).

Mindestumfang:

- OpenTelemetry Collector
- Collector-Konfigurationsdatei
- Compose-Service
- Standardports für OTLP
- Beispielkonfiguration für Logs, Metrics und Traces

### LH-FA-ADD-005 – Mehrfaches Hinzufügen verhindern

*Priorität: MVP.* Das Produkt muss erkennen, ob ein Service bereits vorhanden ist.

Ein bereits vorhandener Service wird nicht doppelt eingefügt; `add` ist idempotent und reaktiviert einen deaktivierten Service. Ein Service gilt als registriert, sobald `services.<name>` in `u-boot.yaml` existiert, und als aktiv mit `enabled: true` und verwaltetem Eintrag in `compose.yaml`. `services.<name>.enabled` ist immer explizit zu setzen (fehlt es, gilt der Service als deaktiviert und `u-boot doctor` warnt). Inkonsistenzen werden nie stillschweigend ignoriert: Ein Block ohne Registrierung bricht den Befehl mit klarer Diagnose ab, eine Registrierung ohne Block erzeugt ihn deterministisch neu.

### LH-FA-ADD-006 – Add-on-Abhängigkeiten

*Priorität: V1.* Das Produkt muss Abhängigkeiten zwischen Add-ons erkennen.

Beispiele:

- Keycloak kann optional PostgreSQL benötigen, wenn `services.keycloak.persistence: external-postgres` in `u-boot.yaml` gesetzt ist.
- OpenTelemetry kann Beispielkonfigurationen für bestehende App-Services erzeugen.

Bei erkannter Abhängigkeit (z. B. `services.keycloak.persistence: external-postgres` ohne PostgreSQL) fährt der Aufruf nie stillschweigend fort: Interaktiv wird nachgefragt, mit `--no-interactive` ohne `--with-deps` bricht er mit Exit-Code `10` ab, mit `--with-deps` bzw. `--yes` werden Abhängigkeiten deterministisch und ohne Rückfrage hinzugefügt.

### LH-FA-ADD-007 – Service entfernen

*Priorität: V1.* Das Produkt muss einen Service wieder entfernen können (`u-boot remove <service>`).

Der Befehl ist nur in einem initialisierten `u-boot`-Projekt nutzbar (`u-boot.yaml` vorhanden).  
Ist keine gültige Projektkonfiguration vorhanden, ist mit klarer Fehlermeldung und Hinweis auf `u-boot init` abzubrechen.

Mindestumfang:

- Service-Eintrag in `compose.yaml` entfernen
- zugehörige verwaltete Blöcke (z. B. in `.env.example`) entfernen
- Eintrag in `u-boot.yaml` auf `enabled: false` setzen
- Volumes nur auf explizite Anforderung (`--purge`) löschen
- Abhängigkeiten anderer Services prüfen und vor dem Entfernen warnen

`--purge` ist eine destruktive Operation; die Bestätigungs- und Modi-Regeln (`--yes`, `--no-interactive`, interaktive Rückfrage) sind verbindlich in [`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung) definiert. Im nicht-interaktiven Modus ohne `--yes` muss der Aufruf mit Exit-Code `10` abgebrochen werden.

Ist der Service bereits auf `enabled: false`, darf der Aufruf idempotent als No-Op mit klarer Meldung beendet werden.

---

## 4.6 Starten und Stoppen der Umgebung

### LH-FA-UP-001 – Umgebung starten

*Priorität: MVP.* Das Produkt muss die Entwicklungsumgebung starten können.

`u-boot up` muss standardmäßig auf den Stabilisierungspfad der aktivierten Dienste warten, bevor der Befehl endet.

- Die Standard-Wartezeit beträgt 60 Sekunden; nach Ablauf erfolgt der Abbruch mit Fehler.
- Die maximale Wartezeit kann über `--timeout <sekunden>` überschrieben werden.
- `--timeout` akzeptiert nur nicht-negative Ganzzahlen (`>= 0`). Negative Werte führen zu Exit-Code `2` und klarer Validierungsfehlermeldung.
- Für Dienste mit Healthcheck ist `healthy` als Zielzustand erforderlich.
- Für Dienste ohne Healthcheck ist `running` als Zielzustand ausreichend.
- Bei definierten Ports wird auf Erreichbarkeit auf `localhost` geprüft, sofern es sich um TCP-basierten Zugriff handelt.
- Für nicht-TCP oder nicht eindeutig probebare Ports darf `up` nicht mit Fehler abbrechen; es ist ein strukturiertes `warn`-Diagnoseergebnis auszugeben.
- Mit `--timeout=0` wird auf das Warten verzichtet; `up` beendet nach Initiierung der Compose-Aktionen.

### LH-FA-UP-002 – Docker Compose verwenden

*Priorität: MVP.* Der Befehl `u-boot up` muss intern Docker Compose verwenden können.

### LH-FA-UP-003 – Startstatus anzeigen

*Priorität: MVP.* Nach dem Start muss das Produkt den Status der relevanten Services anzeigen.

Mindestangaben:

- Service-Name
- Containerstatus
- Port
- Healthcheck-Status, falls vorhanden

### LH-FA-UP-004 – Umgebung stoppen

*Priorität: MVP.* Das Produkt muss die Umgebung stoppen können.

Das Produkt muss zwischen einem regulären Stopp (Container stoppen) und einem vollständigen Aufräumen (Container und Volumes entfernen) unterscheiden (`u-boot down` bzw. `u-boot down --volumes`).

### LH-FA-UP-005 – Logs anzeigen

*Priorität: V1.* Das Produkt soll Logs anzeigen können (`u-boot logs` für alle Services, `u-boot logs <service>` für einen).

Mindestens müssen folgende Optionen unterstützt werden:

- `--follow` – Logs fortlaufend anzeigen
- `--tail <n>` – nur die letzten n Zeilen anzeigen

---

## 4.7 Diagnose

### LH-FA-DIAG-001 – Doctor-Befehl

*Priorität: MVP.* Das Produkt muss eine Diagnosefunktion bereitstellen (`u-boot doctor`).

### LH-FA-DIAG-002 – Lokale Voraussetzungen prüfen

*Priorität: MVP.* Die Diagnosefunktion muss mindestens prüfen:

- Docker installiert (Mindestversion: 24.0.0 oder neuer) **oder** ein Docker-API-kompatibler Drop-in (z. B. Podman ≥ 4.0 mit aktivem `podman.socket` und `DOCKER_HOST` darauf gezeigt); Drop-ins, deren Version nicht erkannt wird, werden als `warn` gemeldet, ohne den Exit-Code zu eskalieren.
- Docker erreichbar
- Docker Compose verfügbar (Mindestversion: 2.20.0 oder neuer)
- Git verfügbar
- Schreibrechte im Projektverzeichnis
- gültige `compose.yaml`, falls vorhanden
- gültige `u-boot.yaml`, falls vorhanden
- falls Devcontainer-Dateien vorhanden sind: syntaktische Gültigkeit von `.devcontainer/devcontainer.json`, Mindestkompatibilität mit VS Code Dev Containers und `forwardPorts`-Konsistenz zu aktivierten Services (Schweregrad `error` bei `devcontainer.enabled == true`, sonst `warn`)
- falls `.devcontainer/Dockerfile` vorhanden ist: Lesbarkeit und erkennbare Build-Basisstruktur (`FROM` vorhanden)

### LH-FA-DIAG-003 – Fehlerklassifikation

*Priorität: MVP.* Die Diagnosefunktion muss Probleme nach Schweregrad klassifizieren.

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

### LH-FA-DIAG-004 – Reparaturhinweise

*Priorität: MVP.* Die Diagnosefunktion muss bei Problemen konkrete Reparaturhinweise ausgeben.

Beispiel:

```text
error: Docker daemon is not reachable.
hint: Start Docker or check your user permissions for /var/run/docker.sock.
```

---

## 4.8 Generatoren

### LH-FA-GEN-001 – Generate-Befehl

*Priorität: MVP.* Das Produkt muss Generatoren mit `u-boot generate <artifact>` anbieten.

Erlaubte Werte für `<artifact>`:

- `changelog`
- `readme`
- `env-example` (`.env.example`)
- `devcontainer`

Bei unbekanntem Artefakt muss der Befehl mit Exit Code `2` abbrechen und die erlaubten Werte explizit zurückgeben.

### LH-FA-GEN-002 – Changelog erzeugen

*Priorität: MVP.* Das Produkt muss ein Changelog erzeugen oder aktualisieren können (`u-boot generate changelog`).

### LH-FA-GEN-003 – README erzeugen

*Priorität: MVP.* Das Produkt muss eine README-Datei erzeugen können (`u-boot generate readme`).

### LH-FA-GEN-004 – Beispiel-ENV erzeugen

*Priorität: MVP.* Das Produkt muss eine `.env.example` erzeugen oder aktualisieren können (`u-boot generate env-example`).

### LH-FA-GEN-005 – Idempotenz

*Priorität: MVP.* Generatoren müssen möglichst idempotent arbeiten.

Das bedeutet:

- mehrfaches Ausführen erzeugt keine unnötigen Duplikate
- bestehende manuelle Inhalte werden möglichst erhalten
- automatisch verwaltete Bereiche sind eindeutig markiert

---

## 4.9 Template-System

### LH-FA-TPL-001 – Projektvorlagen

*Priorität: V1.* Das Produkt soll Projektvorlagen unterstützen (`u-boot init --template <name>`, z. B. `basic`, `micronaut`, `sveltekit`, `micronaut-sveltekit`).

### LH-FA-TPL-002 – Template-Metadaten

*Priorität: V1.* Jedes Template soll Metadaten besitzen.

Mindestangaben:

- Name
- Beschreibung
- Version
- unterstützte Add-ons
- erzeugte Dateien
- benötigte Tools

### LH-FA-TPL-003 – Eigene Templates

*Priorität: Later.* Das Produkt soll später eigene lokale Templates unterstützen können (`u-boot init --template ./pfad`).

### LH-FA-TPL-004 – Templates auflisten

*Priorität: V1.* Das Produkt muss verfügbare Templates auflisten können (`u-boot template list`).

Die Ausgabe muss mindestens enthalten:

- Name
- Beschreibung
- Version

Die Ausgabe muss optional auch maschinenlesbar erfolgen können (`--json`).

---

## 4.10 Konfigurationsdatei

### LH-FA-CONF-001 – Projektkonfiguration

*Priorität: MVP.* Das Produkt muss eine eigene Projektkonfigurationsdatei verwenden.

Beispiel:

```text
u-boot.yaml
```

Die Konfiguration muss über den Konfigurationsbefehl gepflegt werden können (`u-boot config get` und `u-boot config set`).
Die Migrationsfunktion ist in [LH-FA-CONF-006](#lh-fa-conf-006--konfiguration-migrieren) separat beschrieben.

### LH-FA-CONF-002 – Inhalt der Konfiguration

*Priorität: MVP.* Die Konfigurationsdatei muss mindestens `schemaVersion`, `project.name`, je Dienst `services.<name>.enabled` und `devcontainer.enabled` enthalten. Optionale V1-Felder sind `services.keycloak.persistence` (`embedded` | `external-postgres`), `services.otel.enabled`, `devcontainer.featureSources.allow`, `devcontainer.user.uid` (1 bis 65535), `devcontainer.profile` (`default` | `sandbox`) sowie `devcontainer.sandbox.nestedRuntime` (`none` | `podman`), `devcontainer.sandbox.onUnavailable` (`warn` | `fail`), `devcontainer.sandbox.repository` und `devcontainer.sandbox.egress.enabled` / `.allow`.

Hinweise:

- `services.keycloak.persistence` ist optional und kann `embedded` oder `external-postgres` sein. Fehlt der Wert, gilt der Default `embedded`.
- `template` ist optional und Bestandteil des V1-Template-Systems.
- `devcontainer.featureSources.allow` ist optional; fehlt das Feld, ist die Liste leer (nur lokale Features erlaubt).
- Erlaubte Einträge in `devcontainer.featureSources.allow` müssen gültige, non-empty Quell-Strings (z. B. `https://ghcr.io/devcontainers/features/node`) sein.
- Beim Schreiben wird die Liste dedupliziert.
- Die Schlüssel `devcontainer.user.uid`, `devcontainer.profile` und `devcontainer.sandbox.*` sind optional und in [`LH-FA-DEV-004`](#lh-fa-dev-004--benutzerrechte) bis [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion) definiert; ungültige Werte führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).
- Bei ungültigen oder nicht zugelassenen Quellen ist ein fachlicher Validierungsfehler mit Code `10` zu melden.
- Nicht-mandatorische Add-ons dürfen im MVP auf `enabled: false` stehen.
- `services.<name>.enabled` ist immer explizit zu setzen; siehe [`LH-FA-ADD-005`](#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) für die Default-Konvention.
- `enabled: false` bedeutet, dass der Service deaktiviert ist und bei erneutem `u-boot add <service>` wieder aktiviert werden kann.

### LH-FA-CONF-003 – Konfiguration lesen

*Priorität: MVP.* Das Produkt muss die Konfiguration lesen und bei Befehlen berücksichtigen können.

### LH-FA-CONF-004 – Konfiguration aktualisieren

*Priorität: MVP.* Das Produkt muss die Konfiguration aktualisieren können, wenn Add-ons hinzugefügt oder entfernt werden.

### LH-FA-CONF-005 – Konfiguration anzeigen und ändern

*Priorität: MVP.* Das Produkt muss einen Befehl zum Anzeigen und Ändern der Konfiguration bereitstellen (`u-boot config` zeigt die gesamte Konfiguration, `u-boot config get <pfad>` einen Wert, `u-boot config set <pfad> <wert>` setzt einen Wert).

Beim Setzen muss die geänderte Konfiguration auf Schema-Konformität geprüft werden.

### LH-FA-CONF-006 – Konfiguration migrieren

*Priorität: Later.* Das Produkt muss ein Schema-Migrationskommando bereitstellen (`u-boot config migrate`).

Die Migration muss mit einem klaren Fehler auf unbekannte Zukunftsversionen reagieren und bei Migrationen mit älteren Versionen eine Sicherung anlegen.

---

## 4.11 Build- und CI-Infrastruktur des u-boot-Projekts

Die Anforderungen dieses Abschnitts beschreiben den Bau der **u-boot-Codebase selbst** und damit den Entwicklungsprozess, nicht das Produkt. Sie sind seit 0.6.0 entfallen; die Kennungen bleiben als Grabsteine bestehen und werden nicht wieder vergeben.

### LH-FA-BUILD-001 – Multi-Stage Dockerfile (u-boot-Repo)

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-002 – Runtime-Stage Pflichten

*Entfallen (0.6.0).* Die Produktzusage zum ausgelieferten Image steht in [`LH-NFA-PORT-004`](#lh-nfa-port-004--ausgeliefertes-container-image); die Stage-Mechanik steht im technischen Dokument.

### LH-FA-BUILD-003 – Build-Args und Pin-Politik

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-004 – `.dockerignore` Pflicht

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-005 – Makefile mit Standard-Targets

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-006 – Aggregator-Targets

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-007 – Docker-only-Workflow

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-008 – Coverage-Bootstrap

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

### LH-FA-BUILD-009 – Repository-Layout

*Entfallen (0.6.0).* Anforderung an den Bau der u-boot-Codebase, kein Produktmerkmal; sie steht im technischen Dokument und im adoptierten Baseline-Regelwerk.

---

## 4.12 Doku-Struktur des u-boot-Projekts

Die Anforderungen dieses Abschnitts beschreiben die Doku- und Planungsregeln der **u-boot-Codebase selbst** und damit den Entwicklungsprozess, nicht das Produkt. Sie sind seit 0.6.0 entfallen; die Kennungen bleiben als Grabsteine bestehen und werden nicht wieder vergeben.

### LH-FA-PROJDOCS-001 – Mindeststruktur

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-FA-PROJDOCS-002 – ADR-Format

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-FA-PROJDOCS-003 – Planning-Lifecycle

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-FA-PROJDOCS-005 – Carveout-Disziplin

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-FA-PROJDOCS-006 – Dokumentationsreferenzmodell

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-FA-PROJDOCS-004 – Archivierung

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

---

## 4.13 Architektur des u-boot-Projekts

Die Anforderungen dieses Abschnitts beschreiben das Architektur-Pattern der **u-boot-Codebase selbst** und damit eine Entwurfsvorgabe, kein Produktmerkmal. Sie sind seit 0.6.0 entfallen; die Kennungen bleiben als Grabsteine bestehen und werden nicht wieder vergeben.

### LH-FA-ARCH-001 – Hexagonales Pattern

*Entfallen (0.6.0).* Entwurfsvorgabe für die u-boot-Codebase, kein Produktmerkmal; sie steht in der Architektur-Sicht und im technischen Dokument.

### LH-FA-ARCH-002 – Schichten und Verzeichnislayout

*Entfallen (0.6.0).* Entwurfsvorgabe für die u-boot-Codebase, kein Produktmerkmal; sie steht in der Architektur-Sicht und im technischen Dokument.

### LH-FA-ARCH-003 – Import-Regeln und Enforcement

*Entfallen (0.6.0).* Entwurfsvorgabe für die u-boot-Codebase, kein Produktmerkmal; sie steht in der Architektur-Sicht und im technischen Dokument.

---

## 5. Nichtfunktionale Anforderungen

## 5.1 Benutzbarkeit

### LH-NFA-USE-001 – Verständliche Bedienung

*Priorität: MVP.* Das Produkt muss ohne tiefes Vorwissen über die interne Implementierung bedienbar sein.

### LH-NFA-USE-002 – Klare Befehle

*Priorität: MVP.* Befehle müssen sprechend, konsistent und kurz sein.

### LH-NFA-USE-003 – Lesbare Ausgaben

*Priorität: MVP.* CLI-Ausgaben müssen klar strukturiert und gut lesbar sein.

### LH-NFA-USE-004 – Maschinenlesbare Ausgabe

*Priorität: V1.* Das Produkt soll optional maschinenlesbare Ausgabe unterstützen.

Für `--dry-run`/`--diff`-Kombinationen gilt zusätzlich die JSON-Ausgabe in [`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run) und [`LH-FA-CLI-008`](#lh-fa-cli-008--diff-ausgabe).

Für alle `--json`-Ausgaben gilt ergänzend ein gemeinsames Minimalkontrakt-Schema:

- `status` (`ok`/`warn`/`error`)
- `command` (Hauptbefehl als in [`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run) definiertes Enum)
- optional `subcommand` (für gruppierte Befehle wie `template` oder `config`)
- `diagnostics` (Liste von Objekten mit mind. `level`, `code`, `message`, optional `file`)
- `exitCode` (vgl. [`LH-FA-CLI-006`](#lh-fa-cli-006--exit-codes))

Für `--json`-Antworten gilt zusätzlich: `diagnostics.level` ist `warn` oder `error`; enthält `diagnostics` keinen Eintrag, darf es als `[]` ausgegeben werden und `status` ist `ok`, andernfalls folgt `status` dem höchsten `level` (`error` → `error`, `warn` ohne `error` → `warn`); `diagnostics.file` ist optional; bei `template` und `config` ist `subcommand` verpflichtend, und `diagnostics.code` trägt die Kennung der verursachenden Anforderung oder einen dokumentierten tool-internen Code.

Für normale (`--json` ohne `--dry-run`/`--diff`) Ausgaben ist der obige Minimalkontrakt bindend.
Für `--dry-run`- oder `--diff`-Ausgaben mit `--json` gilt zusätzlich das vollständige Schema aus [`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run) als bindender Pflichtkontrakt (inkl. `plannedFiles`, `changes`, `dryRun`, `diff`).

Beispiel: `u-boot doctor --json` im Erfolgsfall liefert `status: ok`, `command: doctor`, eine leere `diagnostics`-Liste und `exitCode: 0`.

---

## 5.2 Zuverlässigkeit

### LH-NFA-REL-001 – Kein stilles Überschreiben

*Priorität: MVP.* Das Produkt darf bestehende Dateien nicht stillschweigend überschreiben.

### LH-NFA-REL-002 – Wiederholbare Ausführung

*Priorität: MVP.* Wiederholte Ausführung desselben Befehls darf das Projekt nicht beschädigen.

### LH-NFA-REL-003 – Abbruch bei kritischen Fehlern

*Priorität: MVP.* Bei kritischen Fehlern muss das Produkt abbrechen und eine klare Fehlermeldung ausgeben.

### LH-NFA-REL-004 – Validierung erzeugter Dateien

*Priorität: MVP.* Das Produkt soll erzeugte Dateien validieren, soweit passende Validatoren verfügbar sind.

Beispiele:

- YAML
- JSON
- Docker Compose

---

## 5.3 Wartbarkeit

### LH-NFA-MAINT-001 – Modulare Architektur

*Priorität: MVP.* Das Produkt muss modular aufgebaut sein.

Insbesondere sollen Add-ons, Templates und Generatoren voneinander getrennt implementiert werden.

### LH-NFA-MAINT-002 – Erweiterbarkeit

*Priorität: MVP.* Neue Services müssen mit geringem Aufwand ergänzt werden können.

### LH-NFA-MAINT-003 – Testbarkeit

*Priorität: MVP.* Die fachlichen Funktionen müssen automatisiert testbar sein.

### LH-NFA-MAINT-004 – Dokumentierte Schnittstellen

*Priorität: V1.* Interne Schnittstellen für Add-ons und Templates sollen dokumentiert werden.

---

## 5.4 Portabilität

### LH-NFA-PORT-001 – Linux-Unterstützung

*Priorität: MVP.* Das Produkt muss Linux als primäre Plattform unterstützen.

### LH-NFA-PORT-002 – Keine unnötigen Systemabhängigkeiten

*Priorität: MVP.* Das Produkt soll möglichst wenige externe Systemabhängigkeiten benötigen.

### LH-NFA-PORT-003 – Containerfreundlichkeit

*Priorität: V1.* Das Produkt soll selbst in einem Container oder Devcontainer ausführbar sein können.

### LH-NFA-PORT-004 – Ausgeliefertes Container-Image

*Priorität: MVP.* Das ausgelieferte Container-Image ([`LH-OPEN-002`](#lh-open-002--paketierung)) ist minimal und enthält keine Shell und keine Build-Toolchain, läuft als Non-root-Benutzer und trägt OCI-Image-Labels; sein `ENTRYPOINT` ist das u-boot-Binary.

---

## 5.5 Sicherheit

### LH-NFA-SEC-001 – Keine Secrets einchecken

*Priorität: MVP.* Das Produkt darf keine echten Secrets in erzeugte Dateien schreiben.

### LH-NFA-SEC-002 – Beispielwerte markieren

*Priorität: MVP.* Beispielwerte in `.env.example` müssen eindeutig als Beispielwerte erkennbar sein.

### LH-NFA-SEC-003 – Sichere Defaults

*Priorität: MVP.* Das Produkt soll sichere Standardwerte verwenden, soweit dies mit lokaler Entwicklung vereinbar ist.

### LH-NFA-SEC-004 – Keine verdeckte Ausführung fremder Skripte

*Priorität: MVP.* Das Produkt darf keinen externen ausführbaren Code aus nicht freigegebenen Quellen ohne ausdrückliche Zustimmung des Nutzers ausführen.

Konkretisierung des Begriffs "externer Code aus nicht freigegebenen Quellen":

- Devcontainer-Features, Templates oder andere Skripte, die nicht lokal im Repository liegen oder nicht ausdrücklich in `u-boot.yaml` freigegeben sind (siehe [`LH-FA-DEV-003`](#lh-fa-dev-003--devcontainer-features), `devcontainer.featureSources.allow`).
- ad-hoc geladene Shell-, Python- oder ähnliche Skripte über HTTP(S) oder andere Netzwerk-Quellen.

Nicht erfasst sind:

- Docker-Images, die in `compose.yaml`-Services oder im Devcontainer-Build explizit konfiguriert sind – sie werden durch `docker pull` regulär aus konfigurierten Registries bezogen und gelten als bewusst gewählte Abhängigkeit des Projekts.
- Pakete, die innerhalb einer Image-Build-Pipeline (z. B. in einem Dockerfile via Paketmanager) installiert werden.

Die Zustimmung ist im interaktiven Modus durch explizite Rückfrage und im nicht-interaktiven Modus durch entsprechende Flag-Optionen (z. B. `--allow-external-feature-sources`) einzuholen.

---

## 5.6 Performance

### LH-NFA-PERF-001 – Schnelle CLI-Antwort

*Priorität: MVP.* Einfache Befehle müssen auf einem typischen Entwicklungsrechner innerhalb folgender Zeiten reagieren (gemessen ohne Docker-Kommunikation, Kaltstart):

MVP:

- `u-boot --help`, `u-boot --version` – unter 200 ms
- `u-boot doctor` (ohne Netz-Wartezeit) – unter 2 s

V1:

- `u-boot config get …` – unter 300 ms

### LH-NFA-PERF-002 – Startzeit abhängig von Docker

*Priorität: MVP.* Die Startzeit von `u-boot up` darf von Docker-Images und Services abhängen, muss aber transparent dargestellt werden.

Insbesondere muss der Fortschritt einzelner Services (Pull, Create, Start, Healthcheck) sichtbar sein.

---

## 6. Schnittstellenanforderungen

## 6.1 Kommandozeilenschnittstelle

### LH-SA-CLI-001 – Befehlsstruktur

*Priorität: MVP.* Die CLI soll die Grundstruktur `u-boot <command> [subcommand|args...] [options]` verwenden.
`subcommand` ist für kommandospezifische Unterbefehle reserviert (z. B. `template`, `config`).
Positionsargumente (z. B. `postgres`, `project.name`) stehen ebenfalls vor den Optionen.

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
| `u-boot egress list`         | Egress-Liste mit Herkunft anzeigen | V2        |
| `u-boot template list`       | Templates anzeigen                 | V1        |

---

## 6.2 Dateischnittstellen

### LH-SA-FILE-001 – Erzeugte Dateien

*Priorität: MVP.* Das Produkt soll folgende Dateien erzeugen oder aktualisieren können:

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

Optional, sobald ein Anwendungs-Dockerfile ([`LH-FA-DOC-002`](#lh-fa-doc-002--dockerfile-erzeugen)) erzeugt wird:

```text
.dockerignore
```

### LH-SA-FILE-002 – Markierte verwaltete Bereiche

*Priorität: MVP.* Automatisch verwaltete Bereiche in Dateien sollen markiert werden.

Markierungsformat: Der Anfang eines verwalteten Bereichs trägt die Markierung `BEGIN U-BOOT MANAGED BLOCK: <name>`, das Ende `END U-BOOT MANAGED BLOCK: <name>`; die Markierung steht als Kommentar der jeweiligen Dateiart (`#` bei YAML, `.env`, `Dockerfile` und Shell-Skripten, HTML-Kommentar bei Markdown, `//` bei JSONC).

- Strikte JSON-Dateien ohne Kommentar-Support werden nicht inline markiert; die gesamte Datei gilt als verwaltet, und der verwaltete Status ist in `u-boot.yaml` zu hinterlegen.

---

## 6.3 Docker-Schnittstelle

### LH-SA-DOCKER-001 – Docker Compose

*Priorität: MVP.* Das Produkt muss Docker Compose aufrufen oder kompatible Compose-Dateien erzeugen können.

### LH-SA-DOCKER-002 – Containerstatus

*Priorität: MVP.* Das Produkt muss den Status laufender Container auslesen können.

---

## 7. Datenanforderungen

### LH-DA-001 – Projektmetadaten

*Priorität: MVP.* Das Produkt muss Projektmetadaten speichern können.

Beispiele:

- Projektname
- Template
- aktivierte Services
- Ports
- Version des `u-boot`-Schemas

### LH-DA-002 – Service-Metadaten

*Priorität: MVP.* Das Produkt muss Informationen über aktivierte Services speichern können.

Beispiele:

- Name
- Image
- Ports
- Volumes
- Environment-Variablen
- Healthchecks

### LH-DA-003 – Schema-Version

*Priorität: MVP.* Die Projektkonfiguration muss eine Schema-Version enthalten.

Beispiel:

```yaml
schemaVersion: 1
```

### LH-DA-004 – Schema-Migration

*Priorität: Later.* Das Produkt muss mit älteren Schema-Versionen umgehen können.

Anforderungen:

- Eine ältere `schemaVersion` muss erkannt und gemeldet werden.
- Das Produkt muss eine automatische Migration anbieten (`u-boot config migrate`).
- Vor der Migration muss eine Sicherungsdatei nach der Backup-Konvention aus [`LH-FA-INIT-005`](#lh-fa-init-005--überschreibschutz) erzeugt werden: primär `u-boot.yaml.bak`; ist bereits ein Backup vorhanden, wird der kleinste freie numerische Suffix verwendet (`u-boot.yaml.bak.1`, `u-boot.yaml.bak.2`, ...) ohne bestehende Backups zu überschreiben.
- Eine unbekannte (zu neue) `schemaVersion` muss zu einem klaren Fehler führen und das Tool darf in diesem Fall keine Dateien verändern.

---

## 8. Qualitätsanforderungen

### LH-QA-001 – Automatisierte Tests

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-QA-002 – Testbare Akzeptanzkriterien

*Entfallen (0.6.0).* Als Lesehinweis erhalten: [`LH-LESE-004`](#lh-lese-004--prüfbarkeit).

### LH-QA-003 – CI-Fähigkeit (GitHub Actions)

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

### LH-QA-004 – Linting (SOLID-nahes Lint-Profil)

*Entfallen (0.6.0).* Regel des Entwicklungsprozesses, kein Produktmerkmal; sie steht im adoptierten Baseline-Regelwerk.

---

## 9. Akzeptanzkriterien

### LH-AK-001 – Minimaler Init-Flow

*Priorität: MVP.* Vorbedingung: eine erreichbare Docker-Engine und Docker Compose in den jeweils geforderten Mindestversionen ([`LH-FA-DIAG-002`](#lh-fa-diag-002--lokale-voraussetzungen-prüfen), [`LH-RISK-001`](#lh-risk-001--docker-versionen)).

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

### LH-AK-002 – PostgreSQL-Flow

*Priorität: MVP.* Folgender Ablauf muss erfolgreich ausführbar sein:

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

### LH-AK-003 – Keycloak-Flow

*Priorität: V1.* Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init
u-boot add keycloak
u-boot up
```

Erwartetes Ergebnis:

- Keycloak-Service ist in `compose.yaml` vorhanden
- Admin-Zugangsdaten werden über `.env.example` mit Platzhaltern dokumentiert (z. B. `KEYCLOAK_ADMIN=CHANGEME_KEYCLOAK_ADMIN`, `KEYCLOAK_ADMIN_PASSWORD=CHANGEME_KEYCLOAK_ADMIN_PASSWORD`)
- Web-Oberfläche ist über den konfigurierten Port (Standard: `8080`) auf `localhost` erreichbar (HTTP 200 oder 302 auf `/`)

### LH-AK-004 – OpenTelemetry-Flow

*Priorität: V1.* Folgender Ablauf muss erfolgreich ausführbar sein:

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

### LH-AK-005 – Devcontainer-Flow

*Priorität: MVP.* Folgender Ablauf muss erfolgreich ausführbar sein:

```bash
u-boot init --devcontainer
```

Erwartetes Ergebnis:

- `.devcontainer/devcontainer.json` existiert
- `.devcontainer/Dockerfile` existiert
- `devcontainer.json` enthält mindestens `name`, und mindestens eines aus `build` oder `image`
- `devcontainer.json` enthält `forwardPorts`, sofern mindestens ein Add-on aktive Ports exportiert
- `u-boot doctor` enthält keinen `error` zu `devcontainer`-Konfiguration oder Feature-Quellen

### LH-AK-006 – Idempotenz

*Priorität: MVP.* Folgender Ablauf darf keine Duplikate erzeugen:

```bash
u-boot add postgres
u-boot add postgres
```

Erwartetes Ergebnis:

- PostgreSQL ist nur einmal in der Konfiguration vorhanden
- das Tool gibt eine verständliche Meldung aus

### LH-AK-007 – Changelog-Generator

*Priorität: MVP.* Folgender Ablauf muss erfolgreich ausführbar sein:

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

### LH-ABG-002 – Keine IDE-Abhängigkeit

`u-boot` darf VS Code Dev Containers unterstützen, soll aber nicht ausschließlich davon abhängig sein.

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

### LH-RISK-002 – Überschreiben manueller Änderungen

Automatische Generatoren können manuelle Änderungen beschädigen.

Maßnahme:

- verwaltete Blöcke
- Backups
- Diff-Anzeige
- `--force` nur explizit

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
- ausgeliefertes Container-Image nach [`LH-NFA-PORT-004`](#lh-nfa-port-004--ausgeliefertes-container-image)

### LH-MVP-002 – Kann nach dem MVP folgen

Nach dem MVP können ergänzt werden:

- `u-boot add keycloak`
- `u-boot add otel`
- Template-System
- lokale Custom-Templates
- Plugin-System

---

## 13. Traceability

Die Rückverfolgbarkeit von Anforderung zu Umsetzung und Test wird aus den `LH-*`-Kennungen abgeleitet (Überschriften dieses Dokuments, Verweise in Planung, Code und Tests); eine separat geführte Matrix gibt es nicht.

---

## 14. Offene Punkte und Entscheidungen

### LH-OPEN-001 – Implementierungssprache (entschieden)

Status: entschieden am 2026-05-21.
Sprache: **Go**.

Mindest-Toolchain: Go 1.26 oder neuer. Die Pin-Hebung im Dockerfile ist Routine ohne separaten Spec-Eintrag.

### LH-OPEN-002 – Paketierung

Status: **GHCR, Binary, Homebrew und Debian/RPM entschieden**, npm
und pip verworfen. Formell offen, bis die Auslieferung aller
gewählten Wege belegt ist.

| Option | Status | Normative Setzung |
| ------ | ------ | ----------------- |
| Container Image (GHCR `ghcr.io/pt9912/u-boot`) | **Gewählt** | Primärer Distributionsweg. |
| Einzelnes Binary | **Gewählt** | Zusätzliches Host-natives Distributionsartefakt. |
| Homebrew | **Gewählt** | Eigener Tap (`pt9912/homebrew-u-boot`), Formel aus den Release-Assets des Tags; stabile Releases, vier Plattformen. |
| Debian/RPM | **Gewählt** | `.deb` und `.rpm` (amd64, arm64) als Release-Assets des Tags; kein gehostetes APT-/DNF-Repository. |
| npm package | Verworfen | Sprach-Ökosystem-Mismatch. |
| pip package | Verworfen | Sprach-Ökosystem-Mismatch. |

### LH-OPEN-003 – Plugin-System (entschieden)

Status: entschieden am 2026-05-31.
Entscheidung: **statisch eingebaute Add-ons** (kein Plugin-Loader).

Re-Evaluation erfolgt nur bei konkretem externem Add-on-Bedarf,
Release-Frequenz-Druck, Contributor-Bedarf für nicht-mainline Add-ons
oder produktiver Sicherheitsanforderung an isoliertes Sub-Tooling.

### LH-OPEN-004 – Template-Format (entschieden)

Status: entschieden am 2026-05-31.
Entscheidung: **YAML-Metadaten + `text/template`-Files**
(Option „YAML-Metadaten plus Dateivorlagen" aus der ursprünglichen
Optionsliste).

Konkretisierung erfolgt über die `LH-FA-TPL-*`-Einträge. Eine spätere
Format-Hebung (z. B. auf Cookiecutter-Kompatibilität oder
OCI-Distribution) braucht konkrete Nutzer- oder Distributions-Trigger.

Vergleich der vier Optionen:

- **YAML-Metadaten plus Dateivorlagen** — gewählt.
- Cookiecutter-kompatible Templates — verworfen (Python-Toolchain + Jinja2-Code-Eval).
- Eigenes Template-System — verworfen (Pflegeaufwand ohne Mehrwert).
- OCI-basierte Template-Pakete — verworfen (prospektive Architektur ohne Use-Case-Trigger).

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

---

## 16. Historie

Dieser Abschnitt ist der **Fußabdruck angenommener Vertragsänderungen**. Er
entsteht nicht aus interner Arbeit: Weder eine ADR noch ein Slice darf eine
`LH-*`-Anforderung ändern — beide referenzieren nur. Eine Änderung an einer
angenommenen Anforderung ist eine Vereinbarung mit dem Projektinhaber; im Repo
hinterlässt sie genau drei Spuren: das Versions-Feld oben, eine Zeile hier und
die geänderte Anforderung selbst.

Die Spalte **Verweis** benennt den *externen* Vorgang (Vereinbarung, Ticket,
Vertragsanhang), nicht das interne Artefakt, das die Änderung ausgeführt hat —
das Lastenheft verweist nie abwärts auf Planung
(Decken-Regel des adoptierten Baseline-Regelwerks).

| Version | Datum | Änderung | Verweis |
| ------- | ---------- | ------------------------------------------------------------ | ---------------------------------------------- |
| 0.1.0 | 2026-05-21 | Initiale Fassung. Der Bestand bis 2026-07-23 entstand in der Entwurfsphase (Status `Entwurf`) und trägt deshalb keine Einzeleinträge — in dieser Phase steuern die IDs noch keine Verbindlichkeit. | — (Entwurfsphase) |
| 0.2.0 | 2026-07-24 | [`LH-FA-PROJDOCS-002`](#lh-fa-projdocs-002--adr-format) auf die MADR-/Nygard-Template-Form umgestellt (Inline-Kopf-Felder inkl. `Schärft`; Pflicht-Abschnitte Alternativen, Fitness Function, Re-Evaluierungs-Trigger, Geschichte). Zum Umstellungszeitpunkt `Accepted` ADRs bleiben in der leanen Form und unveränderlich (Grandfathering). | Vereinbarung mit dem Projektinhaber, ausgelöst durch die Adoption des externen Betriebsregelwerks |
| 0.3.0 | 2026-09-30 | Neue Anforderungen [`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil) bis [`LH-FA-DEV-009`](#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Sandbox-Profil, Container-Runtime mit Degradationstabelle, Egress-Restriktion, Git-Zugangsdaten); Ergänzung von [`LH-FA-DEV-004`](#lh-fa-dev-004--benutzerrechte) (UID-Anpassbarkeit über `devcontainer.user.uid`). | Vereinbarung mit dem Projektinhaber |
| 0.3.1 | 2026-09-30 | [`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil): Verhalten ohne Git-Remote präzisiert (kein Clone-Schritt und Warnung statt Fehler; Zugangsdaten in der Remote-URL bleiben ein Fehler), weil `u-boot init --sandbox` nie einen Remote vorfindet. | Vereinbarung mit dem Projektinhaber |
| 0.3.2 | 2026-09-30 | [`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil): neuer optionaler Schlüssel `devcontainer.sandbox.repository` als Clone-Quelle statt `origin`, damit im Sandbox-Container ein anderes Repository als das Projekt-Repository geklont und bearbeitet werden kann. | Vereinbarung mit dem Projektinhaber |
| 0.3.3 | 2026-09-30 | [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion): Prüfung der Capability beim Containerstart statt im `u-boot doctor` präzisiert (die Capability ist vom Host aus nicht zuverlässig bestimmbar; `doctor` prüft die Konfiguration). | Vereinbarung mit dem Projektinhaber |
| 0.3.4 | 2026-10-01 | [`LH-OPEN-002`](#lh-open-002--paketierung): Homebrew von „vertagt mit Trigger“ auf „gewählt“ gesetzt (eigener Tap, Formel aus den Release-Assets). | Vereinbarung mit dem Projektinhaber |
| 0.3.5 | 2026-10-01 | [`LH-OPEN-002`](#lh-open-002--paketierung): Debian/RPM von „vertagt mit Trigger“ auf „gewählt“ gesetzt (`.deb`/`.rpm` für amd64 und arm64 als Release-Assets, kein gehostetes Repository). | Vereinbarung mit dem Projektinhaber |
| 0.4.0 | 2026-10-02 | Strukturänderung ohne neue Produktzusage: Technische Festlegungen (Schemata, Beispielinstanzen, Algorithmen, Defaults, Build-/CI-, Doku- und Architektur-Details, Markierungsformate) wurden in ein eigenes technisches Dokument überführt; Projektkapitel §4.11–§4.13 und einzelne Anforderungen sind auf Vorgaben gekürzt; neue Lesehinweis-Anforderung [`LH-LESE-003`](#lh-lese-003--dokumentenordnung); die Anforderungen zu ADR-Format, Planning-Lifecycle, Carveout-Disziplin und Dokumentationsreferenzmodell sind auf Kurzzusagen gekürzt (die Einzelregeln führt das adoptierte Baseline-Regelwerk), Aufrufbeispiele und ausführliche Regelblöcke stehen in der technischen Spezifikation; die Traceability-Matrix entfällt (die Rückverfolgbarkeit folgt aus den Kennungen); die veraltete Bootstrap-Angabe „Default-Schwellwert 0“ in [`LH-FA-BUILD-008`](#lh-fa-build-008--coverage-bootstrap) entfällt, der Schwellwert ist überschreibbar. Alle Anforderungs-Kennungen und Überschriften sind unverändert. | Vereinbarung mit dem Projektinhaber |
| 0.5.0 | 2026-10-14 | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): dritter Wert `docker` für `devcontainer.sandbox.nestedRuntime` (eigener Docker-Daemon im Container, `--privileged` als ausgewiesene Lockerung, Speicher je Instanz, Degradationsfall „Daemon nicht betriebsbereit“). | Vereinbarung mit dem Projektinhaber |
| 0.6.0 | 2026-10-16 | Strukturänderung ohne neue Produktzusage: Das Lastenheft beschreibt nur noch Produktanforderungen. 22 Anforderungen an den Entwicklungsprozess der u-boot-Codebase (`LH-FA-BUILD-001`..`009`, `LH-FA-PROJDOCS-001`..`006`, `LH-FA-ARCH-001`..`003`, `LH-QA-001`..`004`) sind entfallen und bleiben als Grabsteine mit unveränderter Kennung und Überschrift bestehen; die Produktzusage zum ausgelieferten Image steht in [`LH-NFA-PORT-004`](#lh-nfa-port-004--ausgeliefertes-container-image), die Prüfbarkeit als Lesehinweis in [`LH-LESE-004`](#lh-lese-004--prüfbarkeit). Der Entwicklungsprozess steht im adoptierten Baseline-Regelwerk. | Vereinbarung mit dem Projektinhaber |
| 0.7.0 | 2026-10-16 | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Zusage für `nestedRuntime: docker` auf die gemessenen Hosts eingegrenzt: Docker auf Linux und Colima (Linux-VM). Docker Desktop ist nicht mehr zugesagt, weil es nicht gemessen ist; Podman als Host-Engine bleibt ausgeschlossen. | Vereinbarung mit dem Projektinhaber |
| 0.8.0 | 2026-10-16 | [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): Ausnahme für Sicherheitslockerungen unter `--quiet` ergänzt: Die Ausweisung von Sicherheitslockerungen (zusätzliche Capabilities, `--privileged`, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) nach [`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion) bleibt sichtbar, ein Skript mit `--quiet` schaltet sie nicht stumm; alle anderen Hinweise entfallen. Die Stufenbeschreibungen („`--quiet` – nur Fehler“) und die Regel „die zuletzt angegebene Option gewinnt“ bleiben unverändert. | Vereinbarung mit dem Projektinhaber |
| 0.9.0 | 2026-10-16 | [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): Vier Ausnahmen unter `--quiet` (Sicherheitslockerungen, Ausweisung der Degradations-Fallbacks nach [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Fortschritt von `u-boot up` nach [`LH-NFA-PERF-002`](#lh-nfa-perf-002--startzeit-abhängig-von-docker), `--force`-Zusammenfassung nach [`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-FA-INIT-005`](#lh-fa-init-005--überschreibschutz)); die ändernden Subkommandos geben unter `--quiet` nur Fehler aus, `u-boot doctor` nur Fehler-Einträge, die Abfragebefehle (`config get`, `config list`, `logs`, `template list`) behalten ihre Ausgabe; die Stufen und die Mehrfachangabe-Regel bleiben unverändert. | Vereinbarung mit dem Projektinhaber |
| 0.10.0 | 2026-10-16 | [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): Die Ausgabe der Vorschau-Optionen `--dry-run` und `--diff` ([`LH-FA-CLI-007`](#lh-fa-cli-007--dry-run), [`LH-FA-CLI-008`](#lh-fa-cli-008--diff-ausgabe)) ist das Ergebnis des Aufrufs und bleibt unter `--quiet`; die Erfolgs- und Zusammenfassungsmeldung eines tatsächlichen Schreibens entfällt auch bei `--diff` ohne `--dry-run`; die Stufen, die vier Ausnahmen und die Mehrfachangabe-Regel bleiben unverändert. | Vereinbarung mit dem Projektinhaber |
| 0.11.0 | 2026-10-17 | [`LH-FA-DEV-009`](#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): Zusage auf Zugangsdaten eines Agenten erweitert (nie in Image, Workspace-Volume, `u-boot.yaml` oder erzeugter Datei; Übergabe per Umgebungsvariable; `doctor` prüft die Quelle); [`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil): Ausweisung eines im Container aktivierten Login-Dienstes (Passwort-Login) im Sandbox-Profil; [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): Ausnahme 1 unter `--quiet` nennt auch diesen Login-Dienst und die Warn-Einträge von `u-boot doctor` zu Sicherheitseigenschaften des Sandbox-Profils, darunter die Meldung, dass die Egress-Restriktion ohne Sandbox-Profil nicht greift. | Vereinbarung mit dem Projektinhaber |
| 0.12.0 | 2026-10-17 | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Zusage für `nestedRuntime: podman` erweitert: `podman build` und `docker build` (Alias) bauen im Container ein Dockerfile mit `FROM`, `RUN`, `COPY`, auch hinter der Egress-Restriktion, gemessen auf Docker 29.5.2 unter Colima 0.10.3 (arm64, VM Ubuntu 24.04 mit `kernel.apparmor_restrict_unprivileged_userns=0`); BuildKit-Erweiterungen und `docker compose` nicht zugesagt. Die Host-Zusage für `podman` („startet unter Docker (inkl. Colima) und Podman“) ist auf Hosts mit erlaubten verschachtelten User-Namespaces begrenzt: mit dem Wert `1` endet das Startscript mit Exit-Code `11`; Podman als Host-Engine ist nicht zugesagt. | Vereinbarung mit dem Projektinhaber |
| 0.13.0 | 2026-10-18 | [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer): Ein fehlendes `/dev/fuse` auf dem Container-Host ist wie `/dev/net/tun` keine Degradationsstufe, sondern ein Startfehler der Container-Engine (gemessen: Docker 29.5.2 auf Colima 0.10.3, `docker run` Exit-Code `127`, `docker start` Exit-Code `1`); der Fallback auf `vfs`-Storage gilt nur für einen Container, der ohne das Device gestartet wurde (etwa nach geänderten `runArgs`); `u-boot doctor` weist ein fehlendes `/dev/fuse` auf dem Container-Host als Umgebungsproblem aus, nicht als Fallback. | Vereinbarung mit dem Projektinhaber |
| 0.14.0 | 2026-10-18 | [`LH-FA-DEV-008`](#lh-fa-dev-008--egress-restriktion): Die wirksame Egress-Allowlist samt Herkunft jedes Hosts ist abfragbar (`u-boot egress list`; offline, ohne laufenden Container, ändert nichts); [`LH-FA-CLI-005`](#lh-fa-cli-005--verbosity-und-logging): der Abfragebefehl behält unter `--quiet` seine Ausgabe; [`LH-SA-CLI-002`](#lh-sa-cli-002--vorgesehene-befehle): die Befehlsübersicht führt den Befehl mit Priorität V2 (aus dem freigegebenen Befehl abgeleitet; Bestätigung des Projektinhabers vor dem Release ausstehend). | Vereinbarung mit dem Projektinhaber |

**Status-Wechsel `Entwurf` → `Accepted` (2026-07-25).** Bis dahin trug dieses
Dokument formal `Entwurf`, obwohl seine IDs bereits als bindend behandelt
wurden (Gates, Exit-Code-Vertrag, Traceability-Matrix). Der Wechsel zieht die
Deklaration an die gelebte Praxis nach; ab hier ist jede Änderung einer
Anforderung vertragspflichtig und bekommt eine Zeile in dieser Tabelle.
Rückwirkend aufgenommen ist nur die eine Änderung, die bereits ausdrücklich als
Vertragsänderung geführt wurde (0.2.0).