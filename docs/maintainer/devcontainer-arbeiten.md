# Devcontainer-Arbeit am u-boot-Repo

Das Repo entwickelt sich selbst im Devcontainer: das erzeugte
`.devcontainer/` ([`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen))
trägt das Sandbox-Profil mit verschachteltem Docker und den gepinnten
Werkzeugen (Go, golangci-lint, make, Ollama, GitHub CLI). Die drei
Bedien-Targets sind in [`tools/mk/devcontainer.mk`](../../tools/mk/devcontainer.mk)
definiert; sie prüfen nichts und hängen nicht an der Gate-Kette.

## Start, Stop, Shell

| Target | Ruft auf | Optionen |
|---|---|---|
| `make dev-up` | `u-boot devcontainer up` | `RECREATE=1` legt den Container neu an (`--recreate`), `INSTANCE=<name>` wählt eine Instanz (`--instance <name>`) |
| `make dev-shell` | `u-boot devcontainer shell` | startet den Container wie `up` und öffnet eine Shell darin; `INSTANCE=<name>` |
| `make dev-down` | `u-boot devcontainer down` | stoppt den Container; `INSTANCE=<name>` stoppt nur diese Instanz |

```bash
make dev-up                     # Devcontainer ohne Instanz starten
make dev-up INSTANCE=agent-1    # zweiter Container desselben Projekts
make dev-shell INSTANCE=agent-1 # Shell in die Instanz agent-1
make dev-down INSTANCE=agent-1  # nur agent-1 stoppen
make dev-up RECREATE=1          # Container ohne Instanz neu anlegen
```

Die Details der Instanzen (eigene Container und Volumes, Verhalten der
Befehle ohne `--instance`, Exit-Code 10 ohne Container-Bestand) führt das
Benutzerhandbuch in Abschnitt 4.9 („Den Devcontainer ohne Editor starten
und stoppen“), Absatz „Mehrere Devcontainer je Projekt (Instanzen)“ —
[Benutzerhandbuch](../user/benutzerhandbuch.md#den-devcontainer-ohne-editor-starten-und-stoppen);
die Make-Targets reichen die Werte unverändert als Flaggen durch.

## Voraussetzungen und Umgebung

- Das Binary `u-boot` muss auf dem Host-Pfad liegen (die Targets brauchen es
  als Bedienung, kein Gate-Vertrag); es fehlt, endet das Target mit Code 2
  (überschreibbar über die Variable `U_BOOT`).
- `dev-up` und `dev-shell` setzen `GIT_TOKEN` und `GH_TOKEN` aus
  `gh auth token`, wenn die Variable am Host nicht steht; ein gesetzter Wert
  hat Vorrang, ohne `gh` bleiben sie leer. Die Werte gehen über die
  Umgebung des Rezepts, nie über die Befehlszeile.
- Der Container läuft im Sandbox-Profil; sein Workspace ist das
  Volume-Clone von origin (`postCreateCommand` der erzeugten
  `devcontainer.json`), **nicht der Arbeitsbaum des Hosts** — eine Änderung
  am Host-Baum kommt per Commit/Push und Pull im Container.
- Der verschachtelte Docker-Daemon startet beim Container-Start
  (`postStartCommand`, Startscript `u-boot-docker-init`, Policy `warn` —
  fällt der Start aus, meldet der Status `fallback` und `u-boot doctor`
  zeigt ihn).

## Gatter im Container

Die `containerEnv` setzt `BUILD_MODE=local`, `GOTOOLCHAIN=local`,
`GOFLAGS` und die Go-Caches unter dem Workspace (je [`SPEC-031`](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)
führen dort die Targets die Toolchain des Containers statt
Docker-Stage-Bauten). `make gates` im Container ist damit der lokale
Modus-Nachweis der Umgebung; eine Messung des Gesamtzugs (Image-Bau,
verschachtelter Daemon, Gates im Container, `test-docker` im
Nested-Daemon und die CI-Runner-Analogie) liegt unter
[`docs/messungen/slice-v1-bau-nur-devcontainer/`](../messungen/slice-v1-bau-nur-devcontainer/README.md).
Die Grenze bleibt die des Vertrags: die Docker-Integrationstests
(`make test-docker`) sind im verschachtelten Daemon nicht als Nachweis
zugesagt — ihr Nachweis bleibt der Docker-Modus an einem Host mit lokaler
Engine und die CI.