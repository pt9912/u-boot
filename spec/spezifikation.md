# Spezifikation — u-boot

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
[`lastenheft.md`](lastenheft.md) formulierten Anforderungen (`LH-*`-Kennungen). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — verbindlich, aber ohne Änderung des Lastenhefts
fortschreibbar. Sie verweist nur aufwärts auf das Lastenheft, nie auf ADRs, Slices,
Carveouts oder die Roadmap.

**Kennungen:** Eine Verfeinerung einer einzelnen Anforderung trägt deren Kennung mit
Buchstabensuffix (`<Anforderungs-Kennung>.a`, `.b`, …). Alles, was keine einzelne
Anforderung verfeinert (Schemata, Defaults, Fehler-Codes, Metrik-Felder, externe
Verträge), trägt `SPEC-<NNN>`: dreistellig, fortlaufend je Datei, nicht je Abschnitt.
Eine `SPEC-*` ist eine Adresse, keine Anforderung.

---

## 1. Algorithmen und Datenflüsse

Verfeinerungen einzelner Anforderungen: Ablauf, Zustandsfolgen, Entscheidungsregeln.
Tatsächlicher Code gehört nicht hierher.

### LH-FA-INIT-002.a — Normalisierung des abgeleiteten Projektnamens

Verfeinert [`LH-FA-INIT-002`](lastenheft.md#lh-fa-init-002--projektname).

1. Der Basisname des Arbeitsverzeichnisses wird auf Kleinbuchstaben gesetzt.
2. Alle Zeichen außer `a-z`, `0-9` und `-` werden auf `-` abgebildet.
3. aufeinanderfolgende `-` werden zu einem einzelnen `-` zusammengeführt.
4. führende und nachgestellte `-` sowie Leerzeichen werden entfernt.
5. Die Länge wird auf 1 bis 63 Zeichen begrenzt.
6. Nach Kürzung auf 63 Zeichen wird erneut auf führende/nachgestellte `-` geprüft und diese notfalls entfernt.
7. Anschließend wird der Name gegen die Validierung in [`LH-FA-INIT-006`](lastenheft.md#lh-fa-init-006--projektnamen-validierung) geprüft.

### LH-FA-DEV-004.a — Übergabe der Benutzer-ID an den Image-Build

Verfeinert [`LH-FA-DEV-004`](lastenheft.md#lh-fa-dev-004--benutzerrechte).

- Der Wert wird als Build-Argument `USER_UID` an den Image-Build übergeben; der Container-Benutzer wird mit dieser UID angelegt.

### LH-FA-DIAG-002.a — Prüfung von Docker und Docker-kompatiblen Drop-ins

Verfeinert [`LH-FA-DIAG-002`](lastenheft.md#lh-fa-diag-002--lokale-voraussetzungen-prüfen).

- Docker installiert (Mindestversion: 24.0.0 oder neuer) **oder**
  ein Docker-API-kompatibler Drop-in (z. B. Podman ≥ 4.0 mit
  aktivem `podman.socket` und `DOCKER_HOST` darauf gezeigt). Die
  Checks `docker.installed` / `docker.reachable` /
  `docker.compose.installed` shellen aus zum `docker`-Binary;
  Drop-ins, die nicht den Docker-Version-Format-Pin
  (`<major>.<minor>.<patch>`-Bereich 24.0 / 2.20) treffen,
  werden als `Severity: warn` („unrecognized version") emittiert,
  ohne den Exit-Code zu eskalieren. Eine formal getestete
  Podman-Variante folgt in einem eigenen Slice bei konkretem
  Bedarf — die heutige MVP-Pflicht ist Docker.

## 2. Datenstrukturen und Schemas

Formate und Schemata (`u-boot.yaml`, JSON-Ausgabe, CLI-Tabellen). Jede Struktur trägt
eine `SPEC-<NNN>`.

### SPEC-001 — JSON-Schema der Vorschau-Ausgabe (`--dry-run --json`)

Gilt für `--dry-run --json` und `--diff --json`; die Pflichtfelder nennt die Anforderung.

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

### SPEC-002 — Beispielinstanz einer Vorschau-Ausgabe (`add --dry-run --json`)

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

### SPEC-003 — Beispielinstanz einer Diff-Ausgabe (`add --diff --json`, ohne `--dry-run`)

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

### SPEC-004 — Projektnamen-Muster

Gilt für den Projektnamen (`project.name`) bei `init` und `config set`.

- regulärer Ausdruck: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`

### SPEC-005 — Schema der Projektkonfiguration `u-boot.yaml`

Gilt als Mindestinhalt der Konfigurationsdatei; die Hinweise zu Schlüsseln und Werten nennt die Anforderung.

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
#   user:
#     uid: 1000                    # 1..65535
#   profile: default               # default | sandbox
#   sandbox:
#     nestedRuntime: none          # none | podman
#     onUnavailable: warn          # warn | fail
#     repository: ''               # Clone-Quelle statt origin (URL)
#     egress:
#       enabled: false
#       allow: []
```

## 3. Defaults und Konstanten

Werte, die im Produkt fest sind (Standardwerte, Grenzwerte, Versionsuntergrenzen).


## 4. Fehler-Codes und Logging-Felder

Verbindliche Diagnose- und Fehler-Codes sowie Logging-Felder.


## 5. Metriken und Tracing-Felder

Verbindliche Telemetrie-Felder pro Span.


## 6. Externe Verträge

Schnittstellen zu Drittsystemen mit Versionsannahme (Docker, Compose, Devcontainer).


## 7. Historie

| Datum | Änderung |
|---|---|
| 2026-10-02 | Angelegt als Gefäß des Technik-Stratums; die Inhalte werden schrittweise aus dem Lastenheft übernommen. |
