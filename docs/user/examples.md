# u-boot — Beispiel-Workflows

Kommando-Rezepte für typische u-boot-Abläufe. **Bewusst nur Befehle,
kein committeter Output:** der byte-genaue Output ist in den Acceptance-/
e2e-Tests gepinnt (Single-Source-of-Truth), nicht hier dupliziert — so
veralten diese Rezepte nicht bei jeder Template-Änderung. Pro-Command-
Details liefert `u-boot <command> --help`; das Verhalten ist verbindlich
in [`spec/lastenheft.md`](../../spec/lastenheft.md) festgelegt.

Alle Subcommands liefern LH-FA-CLI-006-Exit-Codes (`0` Erfolg · `2`
CLI-Fehlnutzung · `10` Validierung · `11` Umgebung · `12` Ausführung ·
`14` IO/Persistenz). Maschinenlesbare Ausgabe + Dry-Run/Diff sind in
[`cli-json-output.md`](cli-json-output.md) dokumentiert.

## 1. PostgreSQL-Stack von Grund auf (MVP-Kern)

```bash
u-boot init my-service      # Projekt-Skelett (docker/, compose.yaml, u-boot.yaml, …) + git init
u-boot add postgres         # PostgreSQL-Service + Volume + .env.example-Block + Healthcheck
u-boot doctor               # Docker/Compose/Git-Voraussetzungen + compose.yaml/u-boot.yaml-Validität
u-boot up                   # docker compose up, wartet auf Healthcheck/TCP-Erreichbarkeit
u-boot logs postgres        # Logs eines einzelnen Service (alle: `u-boot logs`)
u-boot down                 # Container stoppen (Volumes bleiben erhalten)
```

`u-boot add postgres` ist idempotent: ein zweiter Aufruf reaktiviert
einen deaktivierten Service oder meldet „bereits vorhanden", ohne
doppelt einzufügen (`LH-FA-ADD-005`).

## 2. Keycloak + OpenTelemetry

```bash
u-boot add keycloak --with-deps   # Keycloak; --with-deps installiert deklarierte Abhängigkeiten automatisch
u-boot add otel                   # OpenTelemetry Collector + Collector-Config + OTLP-Ports
u-boot up
```

Ist in `u-boot.yaml` `services.keycloak.persistence: external-postgres`
deklariert, hängt Keycloak von PostgreSQL ab (`LH-FA-ADD-006`).
`--with-deps` zieht das fehlende PostgreSQL deterministisch nach; im
nicht-interaktiven Modus **ohne** `--with-deps` bricht der Aufruf mit
Exit `10` und Hinweis auf die fehlende Abhängigkeit ab.

## 3. Devcontainer

```bash
u-boot init my-service --devcontainer    # Projekt inkl. .devcontainer/devcontainer.json + Dockerfile
u-boot generate devcontainer             # nachträglich in ein bestehendes Projekt

# Externe Devcontainer-Feature-Quellen sind standardmäßig gesperrt
# (LH-FA-DEV-003 / LH-NFA-SEC-004) und nur explizit freigebbar:
u-boot generate devcontainer --allow-external-feature-sources ghcr.io/owner/feature
```

`u-boot doctor` prüft die erzeugte `devcontainer.json` auf
VS-Code-Mindestkompatibilität und `forwardPorts`-Konsistenz zu
aktivierten Services.

## 4. Projekt aus einem Template rendern

```bash
u-boot template list                       # eingebauten Katalog browsen
u-boot init my-service --template basic     # aus einem Katalog-Template
u-boot init my-service --template ./my-tpl  # aus einem lokalen Template-Verzeichnis
```

`--template` nimmt einen Katalog-Namen (`basic`) **oder** einen
Dateisystem-Pfad (`./my-tpl`, `/abs/tpl`, `~/tpl`,
`slice-later-local-templates`). Ein lokales Template ist ein Verzeichnis
mit gültiger `template.yaml` plus `*.tmpl`-Dateien (`{{ .Name }}` wird
mit dem Projektnamen gerendert) und beliebigen 1:1-kopierten Dateien.

## 5. CI / maschinenlesbar

```bash
u-boot add postgres --json --dry-run   # geplante Änderungen als JSON, kein Schreiben aufs FS
u-boot add postgres --diff             # Diff-Vorschau des geplanten Endzustands
u-boot up --no-interactive             # keine Rückfragen; offene Bestätigung → Exit 2
u-boot down --volumes --yes            # destruktiv, deterministisch bestätigt (kein Prompt)
```

`--yes` und `--no-interactive` sind exklusiv (`LH-FA-CLI-005A`); bei
destruktiven Operationen (`down --volumes`, `remove --purge`) bricht der
nicht-interaktive Modus ohne `--yes` mit Exit `10` ab. JSON-Envelope-
Schema + Exit-Code-Matrix: [`cli-json-output.md`](cli-json-output.md).

## 6. Service entfernen und aufräumen

```bash
u-boot remove postgres          # deaktivieren + verwaltete Blöcke entfernen (Volumes bleiben)
u-boot remove postgres --purge  # zusätzlich Volumes löschen (destruktiv, Bestätigung)
u-boot down --volumes           # Container + Volumes der gesamten Umgebung entfernen
```

## 7. Konfiguration und Zusatz-Artefakte

```bash
u-boot config get project.name
u-boot config set project.name renamed-service
u-boot config set devcontainer.enabled true
u-boot generate changelog       # changelog | readme | env-example | devcontainer
```

`config set` ist whitelist-begrenzt: nur explizit freigegebene Pfade
(z. B. `project.name`, `devcontainer.enabled`) sind beschreibbar; ein
unbekannter Pfad wird mit Exit `10` abgelehnt (`LH-FA-CONF-*`).
