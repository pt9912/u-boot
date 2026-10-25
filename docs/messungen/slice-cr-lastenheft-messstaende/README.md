# Messprotokoll: CI-Lauf der Integrationstests auf `amd64` (Ollama, verschachteltes Podman und Docker)

Beleg M5 des Slice
[`slice-cr-lastenheft-messstaende`](../../plan/planning/done/slice-cr-lastenheft-messstaende.md)
(Einordnung der Zusage-Menge: dort, §1 *Messlage* und *Herkunft jedes Satzes*). Die Messung ist das Lesen eines
vorhandenen CI-Logs; es gibt kein Skript und keine Rohdatei außer dem Auszug unten. Gelesen am 2026-10-25 mit
`gh api repos/pt9912/u-boot/actions/jobs/121082471598/logs`; GitHub bewahrt das Log nur begrenzt auf, der Auszug
ist deshalb die Ablage.

## Messstand

| | |
|---|---|
| Lauf, Job | Run 40164437737, Job 121082471598 (`integration-docker (//go:build docker)`) |
| Zeit | 2026-10-25, Start 05:08:31Z, Ende 05:23:47Z; Ergebnis des Jobs `success` |
| Stand des Quelltexts | Commit `f0d85e5b4779c8a7badeea06b5a9e96a4b84d09f` |
| Runner | GitHub-Actions-Runner 2.338.0, Image `ubuntu-24.04` (Version 20261012.334.1) |
| Aufruf | `make test-docker` (Integrationstests mit Build-Tag `docker`, Ausgabe mit `-v`, Zeile des Logs: `Run make test-docker`) |
| Docker | Docker-CLI 29.9.2 im Test-Image (`docker-ce-cli` im Log); die Version der Engine des Runners steht nicht im Log und ist nicht gelesen |
| Paket | `github.com/pt9912/u-boot/internal/e2e`: `ok`, 856,974 s |

## Auszug (Zeilen `--- PASS` und `--- SKIP` der betroffenen Tests, Zeitstempel weggelassen)

```text
--- SKIP: TestE2E_SandboxDevcontainer_NestedPodman_Build (22.38s)
--- SKIP: TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress (29.03s)
--- SKIP: TestE2E_SandboxDevcontainer_NestedPodman (22.35s)
--- SKIP: TestE2E_SandboxDevcontainer_StatusFiles_VFS (30.89s)
--- PASS: TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild (30.52s)
--- PASS: TestE2E_SandboxDevcontainer_NestedDocker (37.84s)
--- PASS: TestE2E_SandboxDevcontainer_NestedDocker_Degradation (30.89s)
--- PASS: TestE2E_SandboxDevcontainer_NestedDocker_Restart (24.32s)
--- PASS: TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate (28.99s)
--- PASS: TestE2E_SandboxDevcontainer_Ollama (42.97s)
--- PASS: TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild (16.95s)
```

Das Log hat vier `SKIP`-Zeilen mit derselben Meldung: die drei Podman-Tests (`sandbox_devcontainer_docker_test.go:691`, `:737` und `:778`) und `TestE2E_SandboxDevcontainer_StatusFiles_VFS` (`sandbox_status_docker_test.go:173`; der Auszug oben nennt es wegen des verschachtelten Podman mit):
`host restricts unprivileged user namespaces (kernel.apparmor_restrict_unprivileged_userns=1); nested Podman cannot run here, degradation verified`.
Der Runner hat also den Wert `1`; die Tests prüfen dort die Degradation (nach ihrer Meldung) und überspringen den Lauf von
`podman build` und den Start von `podman`.

## Lesehilfen

- **Belegt** ist der Pass von `TestE2E_SandboxDevcontainer_Ollama`. Das Log trägt für ihn nur die Zeile `--- PASS`; was der Pass
  einschließt (Bau des erzeugten Dockerfiles, Server auf der Loopback-Adresse, `library=cpu`, Pull hinter der Egress-Restriktion,
  Modell überlebt das Löschen des Containers), ist aus den Assertions des Testcodes abgeleitet
  (`internal/e2e/sandbox_devcontainer_docker_test.go`, Prüfliste um die Zeilen 1307 und 1324), nicht aus Markern im Log. Belegt ist ebenso der Pass von
  `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` und der `NestedDocker*`-Tests auf einem `amd64`-Runner mit Docker.
- **Nicht belegt** ist durch die CI der Lauf von `podman build` und der Start von `podman` auf `amd64`: der Runner hat
  `kernel.apparmor_restrict_unprivileged_userns=1`. Diese Zusagen stützen sich auf die Messungen
  [`slice-v1-egress-dns-upstream-schliessen`](../slice-v1-egress-dns-upstream-schliessen/README.md) (Linux `amd64`, Wert `0`)
  und auf die Messung unter Colima (`arm64`).
- Das Log sagt nichts über Colima, `arm64`, Docker Desktop oder eine Engine hinter einer VM.
