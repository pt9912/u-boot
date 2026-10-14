# ADR 0018: Eigener Docker-Daemon im Sandbox-Devcontainer — Herkunft, Start und Lockerung

**Status:** Accepted

**Datum:** 2026-10-14

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0014](0014-nested-podman-sandbox-devcontainer.md), [ADR-0015](0015-sandbox-volumes-pro-instanz.md)

**Schärft:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) — Herkunft der Pakete, Start und Bereitschaftsprüfung des Daemons, Socket-Zugriff und Lockerungsprofil für `devcontainer.sandbox.nestedRuntime: docker`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) kennt seit Lastenheft 0.5.0 den Wert `docker`: ein eigener Docker-Daemon im Sandbox-Container mit BuildKit, Buildx und Compose, ohne Host-Socket, mit Speicher je Instanz und `--privileged` als ausgewiesener Lockerung. Anlass sind Repos, deren Gates BuildKit-Funktionen (`-o type=cacheonly`, `RUN --network=none`), `docker compose` und `docker run --privileged` für eigene dind-Sidecars voraussetzen; rootless Podman ([ADR-0014](0014-nested-podman-sandbox-devcontainer.md)) liefert das nicht. Der Spec lässt offen: (a) Herkunft des Daemons, (b) Start und Bereitschaftsprüfung, (c) Socket-Zugriff, (d) Storage-Driver, (e) Egress-Allowlist.

**Messung (2026-10-14).** Host: Colima 0.10.3 (Linux-VM Ubuntu 24.04.4, arm64, Kernel 6.8.0-117, cgroup v2), Docker 29.5.2. Probe-Image: `mcr.microsoft.com/devcontainers/base:debian` (Debian 13 „trixie“) plus die Debian-Pakete `docker.io` 26.1.5, `docker-cli`, `docker-buildx` 0.13.1, `docker-compose` 2.26.1, `iptables`; Benutzer `vscode` in der Gruppe `docker`. Der Daemon (`dockerd --group docker`) wird per `sudo` gestartet; Daten liegen in einem benannten Volume unter `/var/lib/docker`. Alle Prüfungen laufen als `vscode`.

| Prüfung | Ergebnis |
|---|---|
| Container mit `--privileged`: Daemon startet, `docker info` | bereit; Driver `overlay2`, cgroup v2 |
| `docker buildx build --target <stufe> -o type=cacheonly .` | läuft |
| Dockerfile mit `RUN --network=none` | baut (BuildKit ist Standard, sobald `docker-buildx` installiert ist) |
| `docker compose up` | läuft |
| `docker run --rm --privileged docker:dind docker --version` (dritte Ebene) | läuft; ein dauerhaft laufender `docker:dind` antwortet auf `docker exec … docker info` |
| Zwei Instanzen mit je eigenem Volume, gleichnamige Images und Container | keine Kollision, getrennte Daemon-IDs |
| Container **ohne** `--privileged` | `dockerd` bricht ab (`iptables … Permission denied (you must be root)`) und ist nicht erreichbar |
| `podman build -o type=cacheonly .` (Podman 5, `quay.io/podman/stable`) | legt ein Image an; das BuildKit-Verhalten „nur Cache“ liefert Podman nicht |

Die Pakete `docker.io`, `docker-cli`, `docker-buildx` und `docker-compose` liegen im Debian-Archiv (`deb.debian.org`); auf Debian 13 ist die CLI ein eigenes Paket (`docker-cli`), `docker.io` enthält nur den Daemon.

## Entscheidung

1. **Herkunft: Debian-Pakete** `docker.io`, `docker-cli`, `docker-buildx`, `docker-compose`, `iptables` im generierten Dockerfile auf `devcontainers/base:debian`; kein `docker-ce` aus dem Docker-APT-Repo und kein Devcontainer-Feature `docker-in-docker`. Damit entsteht keine zusätzliche externe Quelle ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) greifen nicht zusätzlich), und die Egress-Allowlist braucht kein `download.docker.com`.
2. **Start per `postStartCommand`** mit `sudo sh /usr/local/bin/u-boot-docker-init` (Skript per `COPY` im Image, wie bei Podman und Egress). Das Skript startet `dockerd --host=unix:///var/run/docker.sock --group docker` abgelöst im Hintergrund und wartet bis zu 60 Sekunden auf `docker info`. Endet der Daemon vorher oder wird er nicht bereit, gilt die Degradation: bei `warn` Warnung „no daemon available“ und Exit 0, bei `fail` Exit 11. Ein bereits laufender Daemon (Neustart-Pfad) wird nicht erneut gestartet. Bei gleichzeitiger Egress-Restriktion läuft deren Skript zuerst.
3. **Socket-Zugriff über die Gruppe `docker`:** `vscode` ist Mitglied (`usermod -aG docker`), der Daemon legt den Socket mit `--group docker` an. Kein `sudo` für `docker`-Aufrufe.
4. **Lockerungsprofil: `--privileged`** und sonst nichts; kein `--cap-add`-Teilsatz. Die Messung ohne `--privileged` zeigt, dass der Daemon iptables, cgroups und Mounts braucht; ein feineres Profil wurde nicht gemessen und nicht zugesagt. Die Befehlsausgabe weist `--privileged` als Warnung aus mit dem Hinweis, dass das Profil keine Isolationsgrenze gegenüber dem Container-Host mehr bietet.
5. **Storage:** Volume `<projekt>-docker-${devcontainerId}` auf `/var/lib/docker` (Benennung nach [ADR-0015](0015-sandbox-volumes-pro-instanz.md)); Driver `overlay2` (gemessen).
6. **Egress ([ADR-0012](0012-devcontainer-egress-firewall.md)):** dieselben Registry-Hosts wie bei `podman` (`registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io`). Die Restriktion ist mit `docker` nur ein Guardrail: der Verkehr verschachtelter Container und die iptables-Regeln des Daemons laufen an ihr vorbei, und `--privileged` hebt sie auf; `u-boot doctor` warnt.
7. **Verhältnis zu ADR-0014:** [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) bleibt für `podman` gültig. Sein Contra zu Alternative B („verstößt gegen [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)“) gilt für `docker` nicht mehr, weil [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) `--privileged` für diesen Wert ausdrücklich verlangt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Host-Socket mounten | kein Daemon im Container | geteilte Image-, Container-, Netz- und Volume-Namen; Socket ist Host-Root-Äquivalent |
| B — `docker-ce` aus dem Docker-APT-Repo | aktuellere Versionen | zusätzliche Quelle mit Schlüsselimport, Egress-Host `download.docker.com` |
| C — Feature `docker-in-docker` | etabliert | externe Quelle (Allowlist, [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)); Skript fremder Herkunft |
| D — rootless `dockerd` | keine Root-Rechte im Container | braucht im Container ebenfalls `--privileged` (Docker-Doku), keine Lockerung gewonnen; nicht gemessen |
| E — Podman mit `docker`-Alias | kleineres Profil | liefert `type=cacheonly`, `--privileged`-dind-Sidecars und die Zielwerkzeuge nicht |
| **F — Debian-Pakete, `dockerd` per `postStartCommand`, `--privileged`** | gemessen lauffähig, keine zusätzliche Quelle, Speicher je Instanz | keine Isolation gegenüber dem Container-Host; Docker 26.1 statt der neuesten Version |

## Konsequenzen

- Positiv: BuildKit-, Compose- und dind-basierte Gates laufen im Sandbox-Profil, parallel in mehreren Instanzen, ohne Host-Socket und ohne zusätzliche Paketquelle.
- Negativ: Mit `--privileged` trifft ein Container-Ausbruch den Container-Host (unter macOS die VM von Docker Desktop oder Colima). Jede Instanz hält ihre eigenen Images (mehrere GB) und startet kalt. Die Egress-Restriktion ist nur noch ein Guardrail.
- Debian-Paketstand (Docker 26.1.5) liegt hinter der Upstream-Version; neuere Funktionen sind nicht zugesagt.
- Folgepflicht: Integrationstest (`//go:build docker`) wiederholt die Messung als Regression; Nachholmessung unter Docker Desktop und unter Linux-Hosts mit rootful Docker.
- **Nicht verifiziert (Stand 2026-10-14):** Docker Desktop (macOS), Podman als Host-Engine, VS Code/Codespaces. Die Messung lief unter Colima (Linux-VM). Die Doku sagt für Docker Desktop erst nach der Nachholmessung Lauffähigkeit zu.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Integrationstest (`//go:build docker`) | `nestedRuntime: docker` baut mit `type=cacheonly` und `RUN --network=none`, startet `docker compose up` und `docker run --privileged docker:dind`; zwei Instanzen haben getrennte Daemons; ohne `--privileged` Exit 0 (`warn`) bzw. 11 (`fail`) | `make test-docker` |
| Golden Case | `devcontainer.json` enthält `--privileged` nur bei `docker`; `none` und `podman` unverändert | `make test` |

## Re-Evaluierungs-Trigger

Ein gemessenes, kleineres Lockerungsprofil als `--privileged` (etwa rootless `dockerd` oder ein Seccomp-Profil), ein Nutzerbericht oder CI-Befund von Docker Desktop, Podman-Host oder Codespaces, der das Profil ändert, oder ein Debian-Paketstand, der Buildx/Compose-Funktionen der Zielrepos nicht mehr trägt.

## Nachholmessung (Docker Desktop, Linux-Host)

Auf der jeweiligen Umgebung ausführen; Ergebnis (Engine-Version, VM-Typ, Kernel, cgroup-Version) in die Geschichte eintragen. Der Integrationstest `TestE2E_SandboxDevcontainer_NestedDocker` deckt die Prüfungen ab: `make test-docker` (Linux) beziehungsweise `go test -tags docker -run NestedDocker ./internal/e2e/` mit der Zielengine als `docker`-Kontext. Ein abweichendes Ergebnis ändert die Entscheidung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-14 | Entwurf mit Messreihe (Colima 0.10.3, Docker 29.5.2, Ubuntu 24.04.4, arm64) | Umsetzungs-Slices V1 Sandbox Docker |
| 2026-10-14 | `Accepted` durch den Projektinhaber; Docker Desktop, Podman-Host und VS Code/Codespaces bleiben bis zur Nachholmessung ungeprüft (§Nachholmessung) | Vereinbarung mit dem Projektinhaber |
