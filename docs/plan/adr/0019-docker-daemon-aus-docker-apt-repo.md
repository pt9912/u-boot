# ADR 0019: Docker-Daemon im Sandbox-Devcontainer aus dem Docker-APT-Repo

**Status:** Accepted

**Datum:** 2026-10-14

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0018](0018-nested-docker-daemon-sandbox-devcontainer.md), [ADR-0012](0012-devcontainer-egress-firewall.md)

**Schärft:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) — Herkunft der Pakete für `devcontainer.sandbox.nestedRuntime: docker`; **ersetzt** Entscheidung 1 und ergänzt Entscheidung 6 von [ADR-0018](0018-nested-docker-daemon-sandbox-devcontainer.md). Alle übrigen Entscheidungen von ADR-0018 bleiben gültig.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0018](0018-nested-docker-daemon-sandbox-devcontainer.md) nahm die Debian-Pakete (`docker.io` 26.1.5, `docker-buildx` 0.13.1, `docker-compose` 2.26.1), um keine zusätzliche Paketquelle einzuführen, und nannte den Preis: ein Paketstand hinter dem Upstream. Dieser Preis ist für die Zielrepos zu hoch; der Daemon im Container soll dem aktuellen Docker entsprechen.

**Messung (2026-10-14).** Host: Colima 0.10.3 (Linux-VM Ubuntu 24.04.4, arm64, Kernel 6.8.0-117, cgroup v2), Docker 29.5.2. Probe-Image: `mcr.microsoft.com/devcontainers/base:debian` (Debian 13 „trixie“; `curl` und `gpg` sind enthalten) plus das Docker-APT-Repo `https://download.docker.com/linux/debian trixie stable`, Schlüssel `https://download.docker.com/linux/debian/gpg` mit dem Fingerprint `9DC858229FC7DD38854AE2D88D81803C0EBFCD88`. Der Daemon wird wie in ADR-0018 per `sudo` mit `--privileged` und einem Daemon-Volume gestartet.

| Prüfung | Ergebnis |
|---|---|
| Paketstand `docker-ce`, `docker-ce-cli` | 5:29.9.0 |
| `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin` | 2.4.2, 0.38.0, 5.6.0 |
| `docker info`: Storage | `overlayfs` (containerd-Image-Store, Snapshotter `io.containerd.snapshotter.v1`) statt `overlay2`; cgroup v2 |
| Abnahmekriterien aus dem CR (Integrationstest `TestE2E_SandboxDevcontainer_NestedDocker`, `_Degradation`, `_Restart`) | laufen unverändert: `type=cacheonly`, `RUN --network=none`, `docker compose up`, dritte Ebene `--privileged`, zwei getrennte Instanzen, Neustart auf demselben Volume |

## Entscheidung

1. **Herkunft: Docker-APT-Repo** (`download.docker.com/linux/debian`, Suite = `VERSION_CODENAME` des Base-Images, Kanal `stable`). Pakete: `docker-ce`, `docker-ce-cli`, `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin`, dazu `iptables`. Das Dockerfile lädt den Schlüssel und bricht den Build ab, wenn dessen Fingerprint nicht `9DC858229FC7DD38854AE2D88D81803C0EBFCD88` ist (Pin im erzeugten Dockerfile). Die Pakete selbst sind nicht auf eine Version gepinnt; ein Image-Build holt den aktuellen `stable`-Stand.
2. **Kein Skript fremder Herkunft:** Es wird nur ein APT-Repo mit geprüftem Schlüssel eingetragen, kein heruntergeladenes Installationsskript ausgeführt; damit bleibt [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) gewahrt. Das Feature `docker-in-docker` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)) bleibt verworfen.
3. **Egress (ergänzt ADR-0018, Entscheidung 6):** `download.docker.com` gehört bei `nestedRuntime: docker` zur Default-Allowlist, weil die Paketquelle im Image eingetragen bleibt und `apt-get update` im eingeschränkten Container sie braucht. Weitere Hosts ergänzt der Nutzer über `devcontainer.sandbox.egress.allow`.
4. **Storage:** Der Daemon-Speicher bleibt das Volume `<projekt>-docker-${devcontainerId}` auf `/var/lib/docker`. Welcher Treiber darauf läuft, entscheidet der Daemon (`overlayfs` mit dem containerd-Image-Store bei Docker 29); der Vertrag macht dazu keine Zusage, und Tests akzeptieren `overlay2` und `overlayfs`.
5. **Unverändert aus ADR-0018:** `--privileged`, Start per `postStartCommand`, Bereitschaftsprüfung und Degradation, Gruppe `docker`, Verhältnis zu ADR-0014.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Debian-Pakete (ADR-0018) | keine zusätzliche Quelle | Docker 26.1, Buildx 0.13, Compose 2.26: zu alt |
| **B — Docker-APT-Repo, Schlüssel per Fingerprint gepinnt** | aktueller Stand, signierte Pakete, kein Fremdskript | zusätzliche Quelle (`download.docker.com`) im Build und in der Egress-Allowlist; Image-Builds sind nicht bitgleich reproduzierbar |
| C — Feature `docker-in-docker` | etabliert | externe Quelle mit Freigabeliste ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)), fremdes Installationsskript |
| D — Binär-Tarball von `download.docker.com/linux/static` | keine Paketverwaltung | kein Update-Pfad, eigene Dienstverwaltung, Prüfsummenpflege |

## Konsequenzen

- Positiv: Docker 29.9, Buildx 0.38 und Compose 5.6 im Container; BuildKit-Funktionen der Zielrepos sind verfügbar.
- Negativ: Der Image-Build braucht Netz zu `download.docker.com`; bei einer Schlüsselrotation bricht der Build mit der Fingerprint-Prüfung ab, bis der Pin angehoben ist (Pin-Hebung per Folge-ADR). Zwei Builds zu verschiedenen Zeiten können verschiedene Docker-Versionen liefern.
- Folgepflicht: Integrationstest (`//go:build docker`) prüft weiter die Abnahmekriterien des CR; die Nachholmessung unter Docker Desktop bleibt offen ([ADR-0018](0018-nested-docker-daemon-sandbox-devcontainer.md), §Nachholmessung).
- **Nicht verifiziert (Stand 2026-10-14):** Docker Desktop (macOS), Podman als Host-Engine, VS Code/Codespaces; die Messung lief unter Colima (arm64) und in der CI auf Linux.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Golden Case | Dockerfile bei `docker` installiert `docker-ce` aus `download.docker.com` mit dem Fingerprint-Pin; bei `none` und `podman` kein `docker-ce`; Allowlist enthält `download.docker.com` nur bei `docker` | `make test` |
| Integrationstest (`//go:build docker`) | Image baut, Daemon startet, Abnahmekriterien des CR laufen mit dem Docker-Repo-Stand | `make test-docker` |

## Re-Evaluierungs-Trigger

Docker rotiert den Paketschlüssel, ändert den Repo-Pfad für Debian oder stellt `stable` für `VERSION_CODENAME` des Base-Images nicht bereit; ein Nutzerbericht, dass eine neuere Docker-Hauptversion das Verhalten der Zielrepos bricht (dann Versions-Pin der Pakete); ein gemessenes, kleineres Lockerungsprofil als `--privileged`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-14 | Entwurf mit Messreihe (Colima 0.10.3, Docker 29.5.2, arm64); `Accepted` durch den Projektinhaber | Vereinbarung mit dem Projektinhaber |
