# Messung: `docker run --env-file` — Eigenschaften, auf die `.devcontainer/devcontainer.env` sich stützt

Vorgang: [`slice-v1-devcontainer-env-datei`](../../plan/planning/done/slice-v1-devcontainer-env-datei.md). Skript: [`tools/messung/devcontainer-env-datei/messen.sh`](../../../tools/messung/devcontainer-env-datei/messen.sh); Rohdaten: [`roh.txt`](roh.txt).
Aufruf aus der Repo-Wurzel: `bash tools/messung/devcontainer-env-datei/messen.sh > <dieses Verzeichnis>/roh.txt`.

## Messstand

Datum 2026-11-04, Docker Engine 29.9.3 (`linux/amd64`), die **verschachtelte** Engine des Devcontainers dieses Repos, Image `busybox:1.36`, **ein Lauf je Fall**, nur Platzhalterwerte. Gemessen ist `docker run --env-file` am Client, nicht der Start über `u-boot devcontainer up`, die Dev-Container-CLI oder VS Code.

## Ergebnis

| Fall | Beobachtung |
|---|---|
| E1 fehlende Datei | der Start bricht ab, Exit-Code `125`, Meldung `--env-file: open …: no such file or directory` |
| E2 Datei nur mit einer Kommentarzeile (`#`) | Start ohne Fehler, Exit-Code `0` |
| E3 Werte | `NAME=WERT` setzt den Wert samt Leerzeichen; Anführungszeichen bleiben Teil des Werts (`Q=["mit anfuehrung"]`); `NAME=` setzt einen leeren Wert; eine Zeile `NAME` ohne `=` setzt die Variable nicht, wenn sie am Client fehlt |
| E4 Präfix `export` | der Start bricht ab, Exit-Code `125`, Meldung `variable 'export DEMO_X' contains whitespaces` |
| E5 relativer Pfad | wird gegen das Arbeitsverzeichnis des Clients aufgelöst; aus einem anderen Verzeichnis `no such file or directory` |
| E6 Vorrang | ein mit `-e` gesetzter Wert gewinnt gegen denselben Namen in der Datei |

## Nicht gemessen

- die Auflösung von `${localWorkspaceFolder}` im Hilfs-Container von `u-boot devcontainer up` und im Editor (VS Code, Dev-Container-CLI)
- andere Engines und Hosts (Docker Desktop, Colima, Podman, `arm64`), andere Engine-Versionen
- die Reihenfolge gegen `remoteEnv` und `containerEnv` der erzeugten Konfiguration
- Dateien mit Zeilenenden `CRLF`, mehrzeilige Werte, sehr große Dateien
- ein zweiter Lauf je Fall
