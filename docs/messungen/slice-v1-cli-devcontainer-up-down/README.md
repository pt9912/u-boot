# Messprotokoll: Einzelmessung des Befehls `devcontainer up` und `down`

Ablauf und Rohdaten der Einzelmessung des ausgelieferten Befehls aus dem Slice
[`slice-v1-cli-devcontainer-up-down`](../../plan/planning/in-progress/slice-v1-cli-devcontainer-up-down.md).
Das Skript liegt unter `tools/messung/devcontainer-start-befehl/` (Ablage:
[`docs/messungen/README.md`](../README.md)). Keine Zahl dieses Protokolls ist eine Zusage.

## Skript und Aufruf

`01-einzelmessung.sh` nimmt das ausgelieferte `u-boot`-Binary des Hosts und zwei vorhandene Ordner und führt
je Projekt (Standardprofil, Sandbox-Profil) den Ablauf der Docker-Tests S1 und S4 aus: `init`, `devcontainer up`
(Dauer, Umschlag, Labels, Einhängungen des Containers, Benutzer), Marker im Workspace-Volume (nur Sandbox-Profil),
`devcontainer down` (zweimal), `up` (gleicher Container?), `up --recreate` (neuer Container, Marker erhalten?),
Aufräumen. Danach folgt die Probe der Grenze: ein Projekt in einem Ordner, den die Engine nicht sieht.

```bash
cd tools/messung/devcontainer-start-befehl
bash 01-einzelmessung.sh <u-boot-binary> <ordner-im-geteilten-bereich> <ordner-ausserhalb> \
  > ../../../docs/messungen/slice-v1-cli-devcontainer-up-down/rohdaten-01-<host>.txt 2>&1
```

Voraussetzung: Bash, Docker (Client und Engine), das Binary, Netzzugriff auf die Registries. Das Skript braucht weder
Go noch Node noch die `devcontainer`-CLI. Es setzt `GIT_TOKEN` auf einen Platzhalter und gibt keine Umgebung aus; die
Rohdaten tragen keine Zugangsdaten. Zeitangaben haben die Auflösung einer Sekunde (das Skript läuft auch unter macOS).

## Probelauf auf Linux (Skript geprüft, keine Colima-Messung)

[`rohdaten-00-probelauf-linux.txt`](rohdaten-00-probelauf-linux.txt): `x86_64`, Linux 6.8, Docker 29.9.1 (Client und Server), Engine
Ubuntu 24.04.5, Cgroup 2, `u-boot version 0.12.0-dev` (aus dem Quellstand gebaut, `CGO_ENABLED=0`), 2026-10-20, ein Lauf.
Beide Ordner lagen unter `/tmp` und waren für die Engine sichtbar; der Lauf prüft das Skript und den Ablauf, nicht Colima.

| Beobachtung | Wert |
|---|---|
| `up` Standardprofil (Basis-Image lokal vorhanden) | Exit 0, 6 s, `status: warn` |
| `up` Sandbox-Profil | Exit 0, 8 s, `status: warn` |
| Labels `devcontainer.local_folder` und `devcontainer.config_file` | Projektordner und `<Ordner>/.devcontainer/devcontainer.json` |
| Einhängungen mit `docker.sock` im gestarteten Container | 0 (beide Profile) |
| `down`, zweites `down` | Exit 0, Container bleibt (`Running` falsch), zweiter Aufruf „nothing to stop“ |
| `up` nach `down` | gleicher Container |
| `up --recreate` | neuer Container, alter entfernt; Marker im Workspace-Volume erhalten (Sandbox) |
| Benutzer im Container | `vscode`, UID 1000 |
| Ordner „außerhalb“ | auf Linux sichtbar: `up` Exit 0 (die Probe der Grenze ist nur auf einem Host mit VM aussagekräftig) |

## Einzelmessung macOS mit Colima

Nicht gelaufen. Sie braucht einen macOS-Host mit Colima und liegt beim Projektinhaber; der Aufruf steht oben,
`<host>` im Dateinamen der Rohdaten nennt den Host (zum Beispiel `colima-arm64`). Das Protokoll hält fest: Colima- und Docker-Version
und Architektur (das Skript schreibt sie im Abschnitt „Messstand“), den Projektpfad im geteilten Bereich und außerhalb. Trägt der Ablauf auf Colima
(Exit 0 und Fortgang wie im Probelauf, bei „außerhalb“ Exit 12 statt eines leeren Bind-Mounts), ist das Ergebnis die Grundlage eines Folge-Änderungsantrags, der Colima in die
Zusage aufnimmt; trägt es nicht, bleibt Colima im Handbuch „nicht zugesagt“.

## Lesehilfen

- Die Dauer eines ersten `up` hängt am Netz und am lokalen Bestand der Basis-Images (Probelauf: beide Basis-Images des Hilfs-Images und das
  Basis-Image des Projekts waren vorhanden). Die Zahlen für ein kaltes System stehen im Protokoll des Fundaments
  ([`slice-v1-devcontainer-start-fundament`](../slice-v1-devcontainer-start-fundament/README.md)).
- Das Skript räumt Container, Volumes, Images `vsc-<Ordner>-…` und das Projektverzeichnis auf; das Hilfs-Image `u-boot-devcontainer-cli:…` bleibt liegen.
