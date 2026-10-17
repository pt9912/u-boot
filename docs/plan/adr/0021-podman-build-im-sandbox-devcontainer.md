# ADR 0021: `podman build` im Sandbox-Devcontainer mit `nestedRuntime: podman` — sechste Lockerung, Egress-Host, Abnahmetest

**Status:** Accepted

**Datum:** 2026-10-17

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0014](0014-nested-podman-sandbox-devcontainer.md), [ADR-0012](0012-devcontainer-egress-firewall.md)

**Schärft:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) — Lockerungsprofil, Egress-Hosts und Abnahmekriterium für `podman build` bei `devcontainer.sandbox.nestedRuntime: podman`; **ersetzt** Entscheidung 2 von [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) (Lockerungsprofil „genau“ fünf Optionen) und ergänzt die Default-Allowlist aus [ADR-0012](0012-devcontainer-egress-firewall.md) Entscheidung 4. Alle übrigen Entscheidungen von ADR-0014 bleiben gültig.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Projektinhaber verlangt, dass im Sandbox-Devcontainer mit `nestedRuntime: podman` mindestens `podman build` funktioniert (2026-10-17); das Kommando `buildah` ist nicht gefordert. [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) hat nur `podman run --rm --network=host alpine echo …` gemessen — mit `--network=host`, das Netzwerk der verschachtelten Container blieb außen vor. „Image-Builds ohne Host-Socket“ aus ihren Konsequenzen war nicht gemessen.

**Messung (2026-10-17).** Probe-Image: das Erzeugnis des Generators (`u-boot init --devcontainer --sandbox`, `nestedRuntime: podman`, `generate devcontainer`; für Egress zusätzlich `egress.enabled`), Basis `mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99` (Debian 13.7), `podman` 5.4.2+ds1-2+b2, `crun` 1.21-1, `netavark` 1.14.0-2, `passt` 0~git20250503. Host: macOS 27.0.1 (arm64), Colima 0.10.3 (Virtualization.Framework), VM Ubuntu 24.04.4, Kernel 6.8.0-117-generic (aarch64), Docker-Client 29.9.1, Docker-Server 29.5.2.

| Prüfung | `kernel.apparmor_restrict_unprivileged_userns = 1` (Standard der VM) | `= 0` (vom Projektinhaber gesetzt) |
|---|---|---|
| Startscript mit den fünf Lockerungen | **Exit-Code 11** (`unshare`: `write failed /proc/self/uid_map`; AppArmor-Profil `unconfined`); `--privileged` gleich | Exit-Code 0; `overlay`, rootless, `crun`, `netavark` |
| Vier Lockerungen ohne `apparmor=unconfined` (`docker-default`) | Start Exit-Code 0, `podman build` scheitert: `overlay` → `mount … permission denied`, `vfs` → `remount /, flags: 0x44000: permission denied` | nicht wiederholt |
| `podman build` (`FROM`, `RUN`, `COPY`) mit den fünf Lockerungen | nicht messbar | Pull läuft, **`RUN` scheitert**: `pasta failed with exit code 1: Failed to open() /dev/net/tun`, Exit-Code 125 |
| dasselbe plus `--device=/dev/net/tun` | nicht messbar | **läuft** (Exit-Code 0); `podman run` des Images liefert den `COPY`- und den `RUN`-Inhalt |
| `docker build` über den Alias; `FROM alpine` | — | Exit-Code 0; Kurzname `alpine` löst auf, `mycompany/img` nicht (`no unqualified-search-registries are defined`) |
| Layer-Cache, `--no-cache` | — | zweiter Build 2× `Using cache`; mit `--no-cache` keiner |
| Debian-Slim mit Paketinstallation (`RUN` mit Netz) | — | kalt 7 s, Cache 0 s, `--no-cache` 4 s; Rootless-Storage nach drei Builds 142 MB |
| Ohne `/dev/fuse` | — | `vfs`-Fallback, Build und `docker build` laufen |
| Egress: `FROM docker.io/library/alpine:3` mit Default-Allowlist | (Ersatzlauf als `root`) scheitert, dreimal: `lookup production.cloudfront.docker.com … server misbehaving` | scheitert gleich; **nach Aufnahme des Hosts** läuft Pull und Build, `RUN apt-get update` läuft, `getent hosts example.com` löst nicht auf, `deb.debian.org` löst auf |
| Integrationstests | übersprungen (Degradation Exit-Code 11 geprüft) | `…NestedPodman_Build` und `…_BuildEgress` grün, ohne Skip |

Nicht gemessen: amd64; andere Engines; große Images; Kurznamen außer `alpine`.

## Entscheidung

1. **Sechste Lockerung: `--device=/dev/net/tun`.** Das Lockerungsprofil bei `nestedRuntime: podman` ist `--cap-add SYS_ADMIN`, `seccomp=unconfined`, `apparmor=unconfined`, `systempaths=unconfined`, `--device /dev/fuse` und `--device /dev/net/tun`; kein `--privileged`, kein `NET_ADMIN` (nicht nötig, gemessen). Das Device braucht `pasta` für das Netzwerk der verschachtelten Container; ohne es scheitert jedes `RUN`. Jede Option erscheint einzeln in der Warnungsausgabe ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)). Das Device muss auf dem Container-Host vorhanden sein, sonst startet der Container nicht.
2. **Paketumfang unverändert** ([ADR-0014](0014-nested-podman-sandbox-devcontainer.md) Entscheidung 1): kein weiteres Paket, auch nicht `buildah`; `podman build` nutzt die eingebaute Bibliothek. `apparmor=unconfined` bleibt nötig (gemessen: mit `docker-default` scheitert der Build an `mount`/`remount`).
3. **Egress ([ADR-0012](0012-devcontainer-egress-firewall.md)):** `production.cloudfront.docker.com` gehört bei `nestedRuntime: podman` und `docker` zur Default-Allowlist, weil Docker Hub Blob-Abrufe dorthin umleitet; ohne den Host scheitert jeder Pull hinter der Restriktion.
4. **Beleg und Zusage:** Die Integrationstests `TestE2E_SandboxDevcontainer_NestedPodman_Build` (kalt, Cache, `--no-cache`, `docker`-Alias) und `…_BuildEgress` (Pull und `RUN` hinter der Restriktion) sind die Abnahme; auf Hosts mit `kernel.apparmor_restrict_unprivileged_userns=1` prüfen sie die Degradation (Exit-Code 11) und überspringen, was den Build nicht belegt. Das Lastenheft 0.12.0 sagt `podman build` mit dem Messstand zu; die Zusage „startet unter Docker (inkl. Colima)“ wird auf Hosts begrenzt, die verschachtelte User-Namespaces erlauben, weil der Standardwert der Ubuntu-24.04-VM von Colima sie blockiert.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, nur `podman run --network=host` zusagen | keine Lockerung mehr | erfüllt die Anforderung nicht; `RUN` ohne Netz unbrauchbar |
| B — Der Build läuft nur mit `--network=host` (kein `/dev/net/tun`) | keine sechste Option | `podman build` bekäme das Flag in jedem Aufruf; nicht der Aufruf der Zusage; nicht gemessen als Alternative |
| C — `--privileged` | läuft | maximale Lockerung, nicht einzeln ausweisbar |
| D — Paket `buildah` und weitere vorsorglich installieren | `buildah`-Kommando verfügbar | nicht nötig (gemessen), nicht gefordert |
| **E — sechste Lockerung `--device=/dev/net/tun`, Paketset unverändert, Egress-Host ergänzt** | gemessen lauffähig, einzeln ausgewiesen, kleinste Änderung | ein weiteres Device im Profil, das nahe an `--privileged` liegt |

## Konsequenzen

- Positiv: `podman build` und `docker build` laufen im Container, auch hinter der Egress-Restriktion; die Messung ist als Tests wiederholbar; Pulls von Docker Hub hinter der Restriktion funktionieren.
- Negativ: Ein weiteres Device im Profil; `/dev/net/tun` muss auf dem Host existieren (sonst Startfehler des Containers, ohne Degradationsstufe wie bei `/dev/fuse`); auf Ubuntu-24.04-VMs mit Standardwert 1 läuft nichts davon.
- Folgepflicht: amd64-Messung über den CI-Lauf; falls der Schalter auf dem Runner 1 ist, entscheidet der Projektinhaber über den Weg (Gate-Änderung).
- **Nicht verifiziert (Stand 2026-10-17):** amd64; Docker Desktop; Podman als Host-Engine; Speicherbedarf großer Images.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Golden Case / Unit | `runArgs` bei `podman` enthält genau die sechs Optionen, je eine Warnung; die Allowlist enthält `production.cloudfront.docker.com` bei `podman` und `docker` | `make test` |
| Integrationstest (`//go:build docker`) | `podman build` und `docker build` (kalt, Cache, `--no-cache`) sowie Build hinter der Egress-Restriktion; auf beschränkten Hosts Degradation Exit-Code 11 und Skip | `make test-docker` |

## Re-Evaluierungs-Trigger

Eine neue `podman`-Version in Debian oder ein neues Base-Image; eine Änderung der Blob-Hosts von Docker Hub; ein Messbefund auf amd64 oder auf einer anderen Engine, der das Profil ändert; ein kleineres Profil (etwa ein `pasta`-Modus ohne `/dev/net/tun`).

## Anhang: Messbefehle (Kurzform)

```bash
# Probe-Image aus dem Generator, Messlauf unter ~/ (Colima teilt nur das Home)
docker run --rm -v "$M/pod":/w -w /w --user 0 u-boot:latest --no-interactive init --no-git --devcontainer --sandbox
# … config set devcontainer.sandbox.nestedRuntime podman; generate devcontainer; docker build -t mpod-probe .devcontainer
RA="--cap-add=SYS_ADMIN --security-opt=seccomp=unconfined --security-opt=apparmor=unconfined --security-opt=systempaths=unconfined --device=/dev/fuse"
docker run --rm $RA mpod-probe sh /usr/local/bin/u-boot-sandbox-init               # Exit 11 bei Schalter 1
docker run --rm $RA mpod-probe bash -c 'podman build -t probe .'                    # Schalter 0: RUN scheitert an /dev/net/tun
docker run --rm $RA --device=/dev/net/tun mpod-probe bash -c 'podman build -t probe .'   # läuft
```

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-17 | Entwurf mit Messprotokoll in zwei Läufen (Colima 0.10.3, Docker 29.5.2, arm64; Schalter 1, dann 0); `Accepted` durch den Projektinhaber (Lastenheft 0.12.0) | Vereinbarung mit dem Projektinhaber |
