# ADR 0021: `podman build` im Sandbox-Devcontainer mit `nestedRuntime: podman` — Messgrenze, Abnahmetest, Zusage nach Messung

**Status:** Proposed

**Datum:** 2026-10-17

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0014](0014-nested-podman-sandbox-devcontainer.md), [ADR-0012](0012-devcontainer-egress-firewall.md)

**Schärft:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) — Abnahmekriterium und Beleg für `podman build` bei `devcontainer.sandbox.nestedRuntime: podman`; ergänzt [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) (Paketumfang und Lockerungsprofil bleiben unverändert).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Projektinhaber verlangt, dass im Sandbox-Devcontainer mit `nestedRuntime: podman` mindestens `podman build` funktioniert (2026-10-17); das Kommando `buildah` ist nicht gefordert. [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) hat nur `podman run --rm --network=host alpine echo …` gemessen; „Image-Builds ohne Host-Socket“ aus ihren Konsequenzen ist **nicht gemessen**. Der Abnahmepunkt „`docker build` läuft im Container ohne Docker-Socket“ des CR zum Sandbox-Devcontainer hat bisher keinen Beleg. Zu klären sind (a) ob das Paketset und die fünf Lockerungen für `podman build` genügen, (b) ob der Build mit Egress-Restriktion und Default-Allowlist läuft, (c) wo sich der Befund wiederholbar belegen lässt.

Randbedingung des Projektinhabers (2026-10-17): Die Colima-Engine-VM bleibt unverändert (kein `sysctl`, kein Eingriff); amd64 wird nur über den CI-Lauf gemessen.

**Messung (2026-10-17).** Probe-Image: das Erzeugnis des Generators (`u-boot init --devcontainer --sandbox`, `nestedRuntime: podman`, `generate devcontainer`), Basis `mcr.microsoft.com/devcontainers/base@sha256:6ee2b4097c75d33e9a94d941577a3286e5d629434fba384688b06daeeaaeac99` (Debian 13.7). Host: macOS 27.0.1 (arm64), Colima 0.10.3 (Virtualization.Framework), VM Ubuntu 24.04.4, Kernel 6.8.0-117-generic (aarch64), Docker-Client 29.9.1, Docker-Server 29.5.2; `kernel.apparmor_restrict_unprivileged_userns = 1`.

| Prüfung | Ergebnis |
|---|---|
| Paketstand im Probe-Image | `podman` 5.4.2+ds1-2+b2, `crun` 1.21-1, `conmon` 2.1.12-4, `netavark` 1.14.0-2, `passt` 0~git20250503, `uidmap` 1:4.17.4-2, `fuse-overlayfs` 1.14-1+b1, `slirp4netns` 1.2.1-1.1; nicht installiert: `buildah` (`command -v buildah` Exit-Code 1), `aardvark-dns`, `catatonit`, `runc` |
| Abhängigkeiten `podman` (Debian 13) | Depends: `conmon`, `crun` oder `runc`, `golang-github-containers-common`, `netavark`; Recommends unter anderem `buildah`, `catatonit` oder `tini`, `passt`, `slirp4netns`, `uidmap` |
| Startscript mit den fünf Lockerungen | **Exit-Code 11**: `unshare --user --map-root-user true` als `vscode` scheitert mit `write failed /proc/self/uid_map: Operation not permitted`; AppArmor-Profil `unconfined`. `--privileged` scheitert gleich |
| Vier Lockerungen ohne `apparmor=unconfined` (Profil `docker-default`) | Startscript Exit-Code 0; `podman build` scheitert: `overlay` → `mount …/storage/overlay: permission denied`, `--storage-driver vfs` → `remount /, flags: 0x44000: permission denied` |
| `podman build` rootless mit den fünf Lockerungen | **nicht messbar** auf dieser VM (siehe Zeile Startscript) |
| Ersatzlauf als `root` im Container, fünf Lockerungen (nicht das Produkt-Profil) | Pull von `docker.io/library/alpine:3` läuft; `RUN` scheitert an `setrlimit RLIMIT_NOFILE`, mit `--ulimit` an `cgroup.subtree_control` (read-only): rootful-spezifisch, kein Befund über rootless |
| Kurzname im Ersatzlauf (nur `podman pull`) | Alias `alpine` löst auf; `mycompany/img` scheitert (`no unqualified-search-registries are defined`) |
| Egress-Pull im Ersatzlauf (`NET_ADMIN`, `u-boot-egress-init` Exit-Code 0) | `docker.io/library/alpine:3` und `debian:stable-slim` scheitern, dreimal wiederholt: `lookup production.cloudfront.docker.com … server misbehaving`; die Default-Allowlist trägt `production.cloudflare.docker.com` |

## Entscheidung

1. **Paketumfang und Lockerungsprofil aus [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) bleiben unverändert.** Kein Messbefund verlangt ein weiteres Paket (auch nicht `buildah`) oder eine weitere Lockerung; die Messung des Basisfalls war nicht möglich. `apparmor=unconfined` ist gemessen nötig: mit dem Profil `docker-default` scheitert der Build an `mount`/`remount`.
2. **Der Beleg für `podman build` ist der Integrationstest** `TestE2E_SandboxDevcontainer_NestedPodman_Build` (`//go:build docker`): `FROM`, `RUN` und `COPY` per `podman build` als `vscode`, danach ein Build aus dem Layer-Cache (`Using cache`), ein Build mit `--no-cache` und ein Build über den `docker`-Alias; jeweils mit den `runArgs` des Generators. Auf Hosts mit `kernel.apparmor_restrict_unprivileged_userns=1` prüft der Test die Degradation (Exit-Code 11) und überspringt; ein übersprungener Test belegt den Build nicht.
3. **Zusage erst nach Messung.** Das Lastenheft sagt `podman build` erst zu, wenn der Test auf mindestens einem Host grün gelaufen ist; der Messstand (Host, Engine, Architektur) steht im Wortlaut. Bis dahin gilt: `podman build` ist **nicht gemessen und nicht zugesagt**. Handbuch und Spezifikation sagen das.
4. **Zwei Befunde gehen an den Projektinhaber, nicht an diese ADR:** (a) Die Default-Allowlist enthält `production.cloudfront.docker.com` nicht, auf das Docker Hub den Blob-Abruf umleitet; unter Restriktion scheitert damit jeder Pull von Docker Hub (betrifft `podman` und `docker`). Die Aufnahme des Hosts ändert [`SPEC-019`](../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) und ist eine Allowlist-Entscheidung. (b) [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) sagt, `podman` starte „unter Docker (inkl. Colima)“; die Colima-VM mit Ubuntu 24.04 blockiert verschachtelte User-Namespaces (Exit-Code 11).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, Aussage „Image-Builds ohne Host-Socket“ stehen lassen | kein Aufwand | die Aussage bleibt ungemessen; die Anforderung des Projektinhabers bleibt unbelegt |
| B — `sysctl kernel.apparmor_restrict_unprivileged_userns=0` in der Colima-VM oder im CI-Lauf setzen und messen | liefert den Beleg | die Colima-VM bleibt nach Entscheidung des Projektinhabers unverändert; der CI-Schritt ändert ein Gate und gehört nicht in diesen Vorgang |
| C — `podman build` im Lastenheft ohne Messung zusagen | erfüllt den Wortlaut der Anforderung | Zusage ohne Beleg; widerspricht dem Gebot, nur Gemessenes zuzusagen |
| D — Paket `buildah` und `aardvark-dns` vorsorglich installieren | `buildah`-Kommando verfügbar | nicht gemessen nötig; ändert Entscheidung 1 von [ADR-0014](0014-nested-podman-sandbox-devcontainer.md) ohne Befund; das Kommando ist nicht gefordert |
| **E — Test als Beleg, Paketset und Lockerungen unverändert, Zusage erst nach grünem Messlauf** | ehrlicher Stand; Messung wiederholbar, sobald ein Host sie zulässt | die Anforderung ist bis dahin nicht belegt |

## Konsequenzen

- Positiv: Der Beleg ist als Test im Repo; Handbuch und Spezifikation sagen nur Gemessenes; zwei Befunde (Egress-Host, Colima-Aussage) sind sichtbar.
- Negativ: `podman build` bleibt unbelegt, bis ein Host mit erlaubten verschachtelten User-Namespaces den Test grün durchläuft; auf Colima mit Ubuntu-24.04-VM und auf `ubuntu-latest` (sofern der Schalter dort 1 ist) überspringt der Test.
- Folgepflicht: Die Messung auf einem Host ohne Restriktion (oder nach Entscheidung des Projektinhabers über die Umgebung) wiederholen; danach Lastenheft-Wortlaut mit Messstand; Egress-Host-Entscheidung.
- **Nicht verifiziert (Stand 2026-10-17):** `podman build` rootless, auf jeder Plattform (arm64 und amd64); Egress mit Build; Speicherbedarf; `vfs`-Degradation beim Build; `docker build` über den Alias.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Integrationstest (`//go:build docker`) | Sandbox-Profil mit `nestedRuntime: podman` baut per `podman build` und `docker build` ein Dockerfile mit `FROM`, `RUN`, `COPY` (kalt, Cache, `--no-cache`); auf beschränkten Hosts Degradation Exit-Code 11 und Skip | `make test-docker` |

## Re-Evaluierungs-Trigger

Ein Messlauf auf einem Host mit erlaubten verschachtelten User-Namespaces (dann `Accepted` und Lastenheft-Wortlaut); eine neue `podman`-Version in Debian oder ein neues Base-Image; eine Änderung der Egress-Hosts von Docker Hub; ein Befund, dass `podman build` mehr als die fünf Lockerungen oder ein weiteres Paket braucht (dann Folge-ADR).

## Anhang: Messbefehle (Kurzform)

```bash
# Probe-Image aus dem Generator, Messlauf unter ~/ (Colima teilt nur das Home)
docker run --rm -v "$M/pod":/w -w /w --user 0 u-boot:latest --no-interactive init --no-git --devcontainer --sandbox
# … config set devcontainer.sandbox.nestedRuntime podman; generate devcontainer; docker build -t mpod-probe .devcontainer
RA="--cap-add=SYS_ADMIN --security-opt=seccomp=unconfined --security-opt=apparmor=unconfined --security-opt=systempaths=unconfined --device=/dev/fuse"
docker run --rm $RA mpod-probe sh /usr/local/bin/u-boot-sandbox-init             # Exit 11 auf der Colima-VM
docker run --rm --cap-add=SYS_ADMIN --security-opt=seccomp=unconfined --security-opt=systempaths=unconfined --device=/dev/fuse \
  mpod-probe bash -c 'podman build -t probe .'                                      # ohne apparmor=unconfined: permission denied
```

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-17 | Entwurf mit Messprotokoll (Colima 0.10.3, Docker 29.5.2, arm64); `Proposed`, bis ein Messlauf den Build belegt | Vereinbarung mit dem Projektinhaber |
