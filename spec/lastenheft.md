# Lastenheft: `u-boot` – Projekt-Bootstrapper für Docker-/Devcontainer-Stacks

| Dokument         | Lastenheft                                                         |
| ---------------- | ------------------------------------------------------------------ |
| Projektname      | `u-boot`                                                           |
| Kurzbeschreibung | CLI-Tool zum Bootstrapping reproduzierbarer Entwicklungsumgebungen |
| Zielplattform    | Linux, Docker, VS Code Dev Containers                              |
| Hauptnutzer      | Softwareentwickler, DevOps-Engineers, technische Teams             |
| Version          | 0.4.0                                                              |
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
- jede Sicherheitslockerung (Capabilities, Seccomp-/AppArmor-Anpassungen, zusätzliche Devices) in der Befehlsausgabe ausweisen.

Nicht-Ziele: `u-boot` startet keinen Agenten und setzt keinen Berechtigungsmodus des Agenten.

Für Devcontainer-Features oder externe Skripte, die das Profil einbindet, gelten [`LH-FA-DEV-003`](#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) unverändert; das Profil gibt keine externe Quelle implizit frei.

### LH-FA-DEV-007 – Container-Runtime im Sandbox-Devcontainer

*Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) soll optional eine rootless Container-Runtime im Container bereitgestellt werden, damit Image-Builds ohne Zugriff auf einen Host-Socket möglich sind.

- Konfigurationsschlüssel `devcontainer.sandbox.nestedRuntime` (`podman` | `none`, Default `none`).
- Bei `podman`: rootless Podman mit `docker`-Kompatibilität im Container.
- Das Ergebnis ist engine-neutral und startet unter Docker (inkl. Colima) und Podman.

**Degradation:** `devcontainer.sandbox.onUnavailable` (`warn` | `fail`, Default `warn`) steuert den Umgang mit fehlenden Fähigkeiten: Bei `warn` führen ein fehlendes `/dev/fuse` zu einem Fallback auf `vfs`-Storage und eine nicht gewährbare Egress-Capability zum Wegfall der Restriktion, jeweils mit Warnung; bei `fail` sind beide ein Umgebungsproblem (Exit-Code `11`). Blockierte verschachtelte User-Namespaces (Seccomp/AppArmor) bei `nestedRuntime: podman` sind in beiden Modi ein Umgebungsproblem (Exit-Code `11`). Eine ausdrücklich angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Ungültige Werte der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

### LH-FA-DEV-008 – Egress-Restriktion

*Priorität: V2.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) soll eine Allowlist für ausgehenden Netzwerkverkehr aktivierbar sein.

- Aktivierung über `devcontainer.sandbox.egress.enabled: true` (Default `false`); keine Flag-Variante, damit `--yes`/`--no-interactive` ([`LH-FA-CLI-005A`](#lh-fa-cli-005a--interaktivität-und-automatisierung)) ohne Sonderfall bleiben.
- Erlaubte Ziele: `devcontainer.sandbox.egress.allow` (Liste von Hostnamen); die Default-Allowlist (gemeinsame Basis plus Ergänzungen je gewähltem Stack) ist dokumentiert.
- Die Allowlist ist unabhängig von `devcontainer.featureSources.allow` ([`LH-FA-DEV-003`](#lh-fa-dev-003--devcontainer-features)): jene steuert erlaubte Build-Quellen, diese die Laufzeit-Ziele.
- Die Restriktion ist ein Guardrail und keine Sandbox-Grenze; Prozesse mit der nötigen Capability können sie aufheben. Die Dokumentation muss das ausdrücklich sagen.
- Ist die nötige Capability nicht gewährbar, greift die Degradation aus [`LH-FA-DEV-007`](#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (geprüft beim Containerstart; `u-boot doctor` prüft die Konfiguration).

### LH-FA-DEV-009 – Git-Zugangsdaten im Sandbox-Devcontainer

*Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](#lh-fa-dev-006--sandbox-profil)) dürfen Git-Zugangsdaten weder im Image, im Workspace-Volume noch in `u-boot.yaml` oder einer anderen erzeugten Datei stehen.

- Übergabe zur Laufzeit per Umgebungsvariable oder als schreibgeschützter Secret-Mount.
- Die Dokumentation beschreibt kurzlebige, auf das Repository begrenzte Tokens (privater Schlüssel nie im Container) und die Forderung nach Branch-Protection auf dem Remote.
- `u-boot doctor` prüft die Token-Quelle: fehlt sie oder liegt sie im Klartext in einer Projektdatei, `warn` (kein Abbruch).

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

Diese Sektion definiert die Build- und CI-Infrastruktur für die **u-boot-Codebase selbst**. Sie ist nicht zu verwechseln mit den Anforderungen aus 4.4 (`LH-FA-DOC-*`), die das Verhalten der Compose-/Dockerfile-Generatoren in den per `u-boot init` erzeugten Zielprojekten beschreiben.

Bezug:

- Implementierungssprache: [`LH-OPEN-001`](#lh-open-001--implementierungssprache-entschieden) (Go).
- Vorlage: das Referenzprojekt `k-deskflight` (Docker-only-Workflow, Multi-Stage Dockerfile, Distroless-Runtime).

### LH-FA-BUILD-001 – Multi-Stage Dockerfile (u-boot-Repo)

*Priorität: MVP.* Die u-boot-Codebase muss ein Multi-Stage `Dockerfile` im Repo-Root bereitstellen.

Die Pflicht-Stages `deps`, `compile`, `test`, `lint`, `coverage`, `build` und `runtime` sind je einzeln per `docker build --target <stage>` baubar.

### LH-FA-BUILD-002 – Runtime-Stage Pflichten

*Priorität: MVP.* Der `runtime`-Stage des u-boot-Dockerfiles muss folgende Eigenschaften erfüllen:

Das Endimage ist minimal und shell-los, läuft als Non-root-Benutzer, trägt OCI-Image-Labels und enthält keine Build-Toolchain.

### LH-FA-BUILD-003 – Build-Args und Pin-Politik

*Priorität: MVP.* Das u-boot-Dockerfile muss versions- und schwellwertbezogene Build-Args bereitstellen: die Go-Version, die golangci-lint-Version und den Coverage-Schwellwert, jeweils mit Default. Die Hebung der Pins ist Routine ohne separaten Spec-Eintrag; der Coverage-Schwellwert lässt sich per `make coverage-gate THRESHOLD=…` überschreiben.

Overrides erfolgen per `docker build --build-arg` bzw. Makefile-Variable.

### LH-FA-BUILD-004 – `.dockerignore` Pflicht

*Priorität: MVP.* Das u-boot-Repo muss eine `.dockerignore` im Repo-Root bereitstellen.

Mindestens auszuschließen sind Versionsverwaltung, IDE- und Agent-Verzeichnisse sowie lokale Build-Artefakte und Caches.

Die `.dockerignore` selbst gehört nicht ins Image und ist daher auszuschließen, sofern sie nicht von einem Stage-Build benötigt wird.

### LH-FA-BUILD-005 – Makefile mit Standard-Targets

*Priorität: MVP.* Das u-boot-Repo muss ein `Makefile` im Repo-Root bereitstellen.

Pflicht-Eigenschaften: `help` als Default-Ziel, alle Targets als `.PHONY`, überschreibbare Variablen mit `?=`-Defaults. Pflicht-Targets decken Hilfe, Abhängigkeitsauflösung, Compile, Lint, Test, Coverage-Gate, Runtime-Image-Build, Smoke-Test und Aufräumen ab.

### LH-FA-BUILD-006 – Aggregator-Targets

*Priorität: V1.* Das Makefile soll Aggregator-Targets bereitstellen:

- `gates` – Inner-Loop-Pflichtgates (`lint` + `test` + `coverage-gate`), PR-blockierend.
- `ci` – `gates` plus `govulncheck` (bei Go-Stack aus [`LH-OPEN-001`](#lh-open-001--implementierungssprache-entschieden) zwingend) plus `image-scan` (Trivy-Image-Scan gegen das Runtime-Image, durch [`LH-QA-003`](#lh-qa-003--ci-fähigkeit-github-actions) als dritter PR-blockierender Job verbindlich — der `make ci`-Aggregator selbst ist V1, sein `image-scan`-Bestandteil ist über [`LH-QA-003`](#lh-qa-003--ci-fähigkeit-github-actions) MVP-Pflicht); SBOM-Erzeugung bleibt optional.
- `fullbuild` – `ci` plus `build`; vollständiger Closure-Lauf.

Aggregator-Targets müssen bei Fehler eines untergeordneten Targets mit Non-Zero-Exit abbrechen und die Fehlerursache klar benennen.

### LH-FA-BUILD-007 – Docker-only-Workflow

*Priorität: MVP.* Der Standard-Build-/Test-Workflow muss ohne hostseitige Sprach-Toolchain auskommen.

Pflicht-Targets rufen ausschließlich `docker build`, `docker run` oder andere solche Targets auf. Voraussetzung am Host sind Docker Engine und `make`; `make` ist ein bewusster Carveout zu [`LH-NFA-PORT-002`](#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten). Weitere Carveouts stehen im `Makefile`-Header.

### LH-FA-BUILD-008 – Coverage-Bootstrap

*Priorität: MVP.* Der `coverage`-Stage muss in der Bootstrap-Phase (noch keine produktiven Pakete in `./internal/...`) deterministisch mit einer leeren Coverage-Eingabe umgehen können.

- Der Schwellwert lässt sich per `make coverage-gate THRESHOLD=…` überschreiben.
- Leere Coverage darf in der Bootstrap-Phase nicht zu einem falschen Grün führen, das echte Test-Failures maskiert.

### LH-FA-BUILD-009 – Repository-Layout

*Priorität: MVP.* Das u-boot-Repo muss folgendem Go-Layout folgen:

- Implementierung unter `./internal/...`, CLI-Entry-Points unter `./cmd/<binary>/` (Binary `uboot`, ausgeliefert als `u-boot`), Unit-Tests neben dem produktiven Code im selben Paket.
- Coverage-Messung ([`LH-FA-BUILD-001`](#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo), [`LH-FA-BUILD-008`](#lh-fa-build-008--coverage-bootstrap)) bezieht sich auf `./internal/...`; `./cmd/...` ist ausgeschlossen.

---

## 4.12 Doku-Struktur des u-boot-Projekts

Diese Sektion definiert die Verzeichnisstruktur unter `docs/` für die **u-boot-Codebase selbst**. Sie ist nicht zu verwechseln mit der `docs/`-Erzeugung in Zielprojekten (siehe [`LH-FA-INIT-003`](#lh-fa-init-003--projektstruktur-erzeugen), [`LH-SA-FILE-001`](#lh-sa-file-001--erzeugte-dateien)), die nur das Top-Level-Verzeichnis anlegt.

Vorlage: die Referenzprojekte `k-deskflight` und `grid-gym` (Basis-Pattern: archive + plan/adr + plan/planning-Lifecycle + user).

### LH-FA-PROJDOCS-001 – Mindeststruktur

*Priorität: MVP.* Das u-boot-Repo muss `docs/` mit den Unterverzeichnissen `archive/`, `plan/adr/`, `plan/planning/` (mit `open/`, `next/`, `in-progress/`, `done/`) und `user/` bereitstellen; weitere Unterverzeichnisse sind zulässig.

Jedes Unterverzeichnis muss mindestens eine `README.md` mit kurzer Zweckbeschreibung enthalten, damit Git die Struktur trackt und Newcomer den Verzeichnisstandard ohne externe Erklärung erfassen können. `.gitkeep` ist als Ersatz unzureichend, weil er den Zweck nicht kommuniziert.

Abgrenzung zu Zielprojekten: Für per `u-boot init` erzeugte Zielprojekte ist nur `docs/` als Top-Level Pflicht ([`LH-FA-INIT-003`](#lh-fa-init-003--projektstruktur-erzeugen)). Ob diese Unterstruktur auch in Zielprojekten erzeugt wird, ist eine spätere Entscheidung (z. B. via Template oder Flag) und gehört nicht zum MVP-Umfang.

### LH-FA-PROJDOCS-002 – ADR-Format

*Priorität: MVP.* Architecture Decision Records in `docs/plan/adr/` folgen dem vendorten MADR-/Nygard-Template des adoptierten Baseline-Regelwerks (`.harness/baseline/<tag>/templates/docs/plan/adr/NNNN-titel.template.md`): vierstellige, nie wiederverwendete Nummer, Kopf-Felder als Inline-Felder (Status, Datum, Autor, Bezug, Schärft) und die vorgeschriebene Abschnittsfolge. Abgelöste ADRs bleiben mit dem Status „Superseded by“ erhalten.

### LH-FA-PROJDOCS-003 – Planning-Lifecycle

*Priorität: MVP.* Planning-Artefakte durchlaufen den Lifecycle `open → next → in-progress → done` des adoptierten Baseline-Regelwerks: Übergang per `git mv`, kein Artefakt in mehreren Verzeichnissen, Inhalte in `done/` nur korrigierend änderbar (substanzielle Änderungen erzeugen ein neues Artefakt in `open/` oder `next/` mit Verweis auf den vorhergehenden Stand).

Dateinamen in `planning/`: `slice-<phase>-<kebab-slug>.md` für Slice-Pläne und `tranche-<nr>-<kebab-slug>.md` für Tranchen-Pläne; die Wahl ist im `README.md` von `docs/plan/planning/` dokumentiert. Übergreifende Master-Dokumente (`roadmap.md`, `carveouts.md`) liegen dauerhaft in `in-progress/` und folgen keinem der beiden Formate.

### LH-FA-PROJDOCS-005 – Carveout-Disziplin

*Priorität: MVP.* Jeder temporäre Carveout in der u-boot-Codebase bekommt parallel zu seiner Entstehung einen Slice-Plan in `docs/plan/planning/open/` mit Auslöser, Aufhebungsbedingung und Akzeptanzkriterien (Carveout-Disziplin des adoptierten Baseline-Regelwerks) und ist im Master-Inventar `carveouts.md` und in der Roadmap sichtbar. Als temporärer Carveout zählen auch Bootstrap-Schwellwerte, bewusst leere Regelblöcke in der Tooling-Konfiguration, prospektive Doku-Phrasen und bewusst weggelassene CI-/Build-Pflichten; Spec-Open-Punkte (`LH-OPEN-*`) und ADR-Folgepunkte gelten ebenfalls. Permanente Carveouts stehen mit Begründung im Master-Inventar und brauchen keinen Aufhebungsplan.

### LH-FA-PROJDOCS-006 – Dokumentationsreferenzmodell

*Priorität: V1.* Das Repo wendet das Dokumentationsreferenzmodell (Referenz-Richtung, Decken-Regel) des adoptierten Baseline-Regelwerks an; `docs-check` erzwingt es. Das Lastenheft ist die normative Decke: Externe Normen, Standards und Vorgaben wirken im Repo nur über explizite `LH-*`-Anforderungen normativ.

### LH-FA-PROJDOCS-004 – Archivierung

*Priorität: V1.* Abgelöste oder veraltete Inhalte aus `user/`, `plan/` oder anderen `docs/`-Bereichen werden nach `docs/archive/` verschoben, statt sie zu löschen.

- Beim Verschieben wird ein kurzer Hinweis am Anfang des Zielfiles ergänzt (z. B. `> Archiviert am YYYY-MM-DD; ersetzt durch [<Pfad>](<pfad>).`).
- Querverweise in lebendiger Doku werden auf das neue Ziel umgebogen oder explizit als historisch markiert.
- Das Verschieben erfolgt per `git mv`, damit die Historie erhalten bleibt.

---

## 4.13 Architektur des u-boot-Projekts

Diese Sektion definiert das Architektur-Pattern für die **u-boot-Codebase selbst** und trägt die Pflichten, die für jede Code-Änderung gelten.

Vorlage: die Referenzprojekte `k-deskflight` (Go, flach), `m-trace` (TypeScript, driving/driven-Split) und `grid-gym` (Python, driving/driven-Split). Die konkrete Variante ist durch [`LH-FA-ARCH-001`](#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](#lh-fa-arch-003--import-regeln-und-enforcement) festgelegt.

### LH-FA-ARCH-001 – Hexagonales Pattern

*Priorität: MVP.* Die u-boot-Codebase muss dem hexagonalen Architektur-Pattern (Ports & Adapters) folgen.

Pflichten:

- Trennung zwischen reiner Domäne, Anwendungslogik, Ports (Interfaces) und Adaptern (Implementierungen).
- Keine direkten Abhängigkeiten von Anwendungslogik oder Domäne zu externen Bibliotheken (Docker-SDK, YAML-Parser, Dateisystem).
- Externe Zugriffe laufen ausschließlich über Driven-Ports, die in Adapter-Paketen implementiert werden.

### LH-FA-ARCH-002 – Schichten und Verzeichnislayout

*Priorität: MVP.* Das u-boot-Repo muss unter `internal/` die Schichten `hexagon` (mit `domain`, `application` und `port/driving`, `port/driven`) und `adapter` (mit `driving` und `driven`) bereitstellen.

Die Wiring-Schicht (`cmd/uboot/`) ist die einzige Stelle, an der `application` und `adapter` zusammen importiert werden dürfen.

### LH-FA-ARCH-003 – Import-Regeln und Enforcement

*Priorität: MVP.* Die Import-Regeln der Schichten sind verbindlich. Kernaussagen: `domain` importiert nur die Go-Standardbibliothek; `application` kennt keine konkreten Adapter; die Port-Pakete `driving` und `driven` kennen einander nicht; Adapter importieren nicht `application` und nicht den jeweils anderen Adapter-Typ; nur `cmd/uboot` verbindet `application` und Adapter.

Pflichten:

- Die Regeln werden im `lint`-Stage ([`LH-FA-BUILD-001`](#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)) per `golangci-lint` mit `depguard` durchgesetzt; Verstöße sind PR-blockierend.
- Die `depguard`-Konfiguration ist deckungsgleich mit den Import-Regeln zu halten; `//nolint:depguard` ist verboten, Carveouts stehen zentral in `.golangci.yml` mit Begründung.
- Die Regeln gelten nur für Produktivcode; Tests sind ausgenommen.

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

*Priorität: MVP.* Für zentrale Funktionen müssen automatisierte Tests vorhanden sein.

Mindestumfang:

- CLI-Befehle
- Dateigeneratoren
- Template-Verarbeitung
- Add-on-Erzeugung
- Konfigurationsparser

### LH-QA-002 – Testbare Akzeptanzkriterien

*Priorität: MVP.* Jede funktionale Anforderung soll durch mindestens einen Akzeptanztest überprüfbar sein.

### LH-QA-003 – CI-Fähigkeit (GitHub Actions)

*Priorität: MVP.* Das u-boot-Repo muss eine CI-Pipeline auf GitHub Actions führen.

Pflicht-Komposition: Die Pipeline läuft bei `pull_request` und `push` auf `main` in drei parallelen, PR-blockierenden Jobs: Gates (`make gates`), Security-Gates (`make govulncheck`) und Image-Scan (`make image-scan`, Trivy gegen das Runtime-Image, Severity HIGH und CRITICAL). Die PR-Blocking-Pflicht aller drei folgt aus diesem Eintrag, auch wenn die Make-Target-Definitionen unter [`LH-FA-BUILD-006`](#lh-fa-build-006--aggregator-targets) liegen. Der Runner braucht nur Docker und BuildKit, keine Host-Go-Toolchain ([`LH-FA-BUILD-007`](#lh-fa-build-007--docker-only-workflow)); Actions sind gepinnt, Token-Rechte minimal gehalten und jeder Job hat ein Zeitlimit. Die Required-Status-Check-Liste im GitHub-UI muss die tatsächlichen Job-Namen des Workflows verwenden.

### LH-QA-004 – Linting (SOLID-nahes Lint-Profil)

*Priorität: MVP.* Die u-boot-Codebase muss ein verschärftes Lint-Profil führen, das über die Default-Linter hinausgeht.

Das Profil besteht aus den Default-Lintern und SOLID-nahen Zusatz-Lintern (Komplexitäts-, Funktionslänge-, Interface-, Kopplungs- und Boundary-Signale); `depguard` für die Schicht-Regeln aus [`LH-FA-ARCH-003`](#lh-fa-arch-003--import-regeln-und-enforcement) ist Teil davon.

Pflichten:

- Die Konfiguration lebt in `.golangci.yml`; `//nolint`-Pragmas sind verboten, Carveouts stehen zentral dort mit Begründung.
- Verstöße brechen den `lint`-Stage ([`LH-FA-BUILD-001`](#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo)) und damit `make gates`, `make ci` und `make fullbuild`; die Linter-Auswahl deckt die Architekturgrenzen aus [`LH-FA-ARCH-003`](#lh-fa-arch-003--import-regeln-und-enforcement) ab.

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
- Build-Infrastruktur der u-boot-Codebase selbst:
  - Multi-Stage `Dockerfile` ([`LH-FA-BUILD-001`](#lh-fa-build-001--multi-stage-dockerfile-u-boot-repo))
  - `Makefile` mit MVP-Pflicht-Targets ([`LH-FA-BUILD-005`](#lh-fa-build-005--makefile-mit-standard-targets))
  - `.dockerignore` ([`LH-FA-BUILD-004`](#lh-fa-build-004--dockerignore-pflicht))
  - Docker-only-Workflow ([`LH-FA-BUILD-007`](#lh-fa-build-007--docker-only-workflow))
  - Repository-Layout nach [`LH-FA-BUILD-009`](#lh-fa-build-009--repository-layout)
- Doku-Struktur der u-boot-Codebase nach [`LH-FA-PROJDOCS-001`](#lh-fa-projdocs-001--mindeststruktur), inkl. ADR-Format ([`LH-FA-PROJDOCS-002`](#lh-fa-projdocs-002--adr-format)), Planning-Lifecycle ([`LH-FA-PROJDOCS-003`](#lh-fa-projdocs-003--planning-lifecycle)), Carveout-Disziplin ([`LH-FA-PROJDOCS-005`](#lh-fa-projdocs-005--carveout-disziplin)) und Dokumentationsreferenzmodell ([`LH-FA-PROJDOCS-006`](#lh-fa-projdocs-006--dokumentationsreferenzmodell))
- Architektur-Pattern (hexagonal, driving/driven-Split) nach [`LH-FA-ARCH-001`](#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](#lh-fa-arch-003--import-regeln-und-enforcement), mit Import-Enforcement via `golangci-lint depguard`
- SOLID-nahes Lint-Profil nach [`LH-QA-004`](#lh-qa-004--linting-solid-nahes-lint-profil); die Konfiguration liegt in `.golangci.yml`, die abgeleitete Quality-Doku erklärt Linter und Carveouts
- CI-Pipeline nach [`LH-QA-003`](#lh-qa-003--ci-fähigkeit-github-actions) (GitHub Actions, drei PR-blockierende Jobs)

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
([`LH-FA-PROJDOCS-006`](#lh-fa-projdocs-006--dokumentationsreferenzmodell)).

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

**Status-Wechsel `Entwurf` → `Accepted` (2026-07-25).** Bis dahin trug dieses
Dokument formal `Entwurf`, obwohl seine IDs bereits als bindend behandelt
wurden (Gates, Exit-Code-Vertrag, Traceability-Matrix). Der Wechsel zieht die
Deklaration an die gelebte Praxis nach; ab hier ist jede Änderung einer
Anforderung vertragspflichtig und bekommt eine Zeile in dieser Tabelle.
Rückwirkend aufgenommen ist nur die eine Änderung, die bereits ausdrücklich als
Vertragsänderung geführt wurde (0.2.0).