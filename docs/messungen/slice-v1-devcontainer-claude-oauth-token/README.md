# Messung: `CLAUDE_CODE_OAUTH_TOKEN` im Sandbox-Devcontainer (slice-v1-devcontainer-claude-oauth-token)

**Datum:** 2026-11-01. **Umgebung:** Linux 7.0.0-42-generic `x86_64`, Docker Engine 29.9.2 (Client und Server, `linux/amd64`, lokale Engine), `u-boot` aus dem Stand des Branches `worktree-agent-aa69746b53d530dc4`, Dev-Container-CLI 0.89.0 (Meldung im Lauf). Lauf als Teil von `make test-docker` (gesamtes Paket `internal/e2e`).

**Skript:** keines unter `tools/messung/`; die Messung ist der Test `TestE2E_ClaudeOAuthToken_ReachesTheSession` in `internal/e2e/claude_oauth_token_docker_test.go` (Build-Tag `docker`). Aufruf: `make test-docker`; der Test legt ein Sandbox-Projekt mit den Features `node` und `claude-code` an, startet den Container mit `u-boot devcontainer shell` und liest die Variable in der Sitzung.

**Wert:** ein erfundener Platzhalter (`sk-ant-oat01-placeholder…`), kein gültiges Token; er steht in keiner Rohdaten-Datei.

## Ergebnis

| Zelle | Läufe | Ergebnis |
|---|---|---|
| Host-Variable gesetzt (Platzhalter) | 2 (je eine Sitzung `devcontainer shell`) | Wert kommt unverändert in der Sitzung an; `ANTHROPIC_API_KEY` bleibt ungesetzt |
| Host-Variable ungesetzt | 2 | in der Sitzung ungesetzt |
| Host-Variable leer gesetzt | 2 | in der Sitzung gesetzt und leer |
| `claude -p` mit dem Platzhalter | **1** | Antwort des Dienstes `Failed to authenticate. API Error: 401 OAuth access token is invalid.`, Exit 1 des Aufrufs; kein Hänger |

Die drei Variablen-Zellen stehen im Test als Zusicherungen (Lauf grün); die Rohdaten (`rohdaten-01.txt`) tragen nur die Ausgabe des `claude -p`-Aufrufs, weil der Test die übrigen Zellen nicht ausgibt. Alle Läufe einer Zelle gehören zu einem Container (eine Zelle nach der anderen, derselbe Containerstart).

## Grenzen

- Ein Lauf mit einem **gültigen** Token ist nicht gefahren; die Antwort des Dienstes auf ein gültiges Token ist nicht belegt.
- `claude -p` wurde einmal gefahren, nicht zweimal.
- Verhalten bei gesetztem `CLAUDE_CODE_OAUTH_TOKEN` im Ollama-Ablauf: nicht gemessen.
- Eine Engine auf macOS, `arm64`, Podman: nicht gemessen.
- Im selben Lauf von `make test-docker` schlugen drei fremde Tests fehl (`TestDevcontainerShell_Docker_Session` mit einem fremden Hilfs-Container `u-boot-devcontainer-up-…` auf der gemeinsam genutzten Engine, `TestE2E_ForgeCLI_PinMismatchBreaksBuild/GLAB` beim Entfernen eines Containers, `TestE2E_LHAK003_KeycloakAcceptanceFlow` mit `compose runtime error`); keiner berührt `agentCredentialEnv`.
