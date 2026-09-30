# ADR 0014: Nested rootless Podman im Sandbox-Devcontainer — Base-Image und Lockerungsprofil

**Status:** Accepted

**Datum:** 2026-09-30

**Autor:** pt9912

**Bezug:** [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0012](0012-devcontainer-egress-firewall.md)

**Schärft:** [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) — Base-Image, Paketumfang und das Set der ausgewiesenen Sicherheitslockerungen für `devcontainer.sandbox.nestedRuntime: podman`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) verlangt optional rootless Podman im Sandbox-Devcontainer, engine-neutral, mit ausgewiesenen Sicherheitslockerungen. Der Spec lässt Base-Image und konkrete Lockerungen offen. Zu klären sind (a) das Base-Image, (b) welche Lockerungen Docker minimal braucht, (c) ob `NET_ADMIN` für die Egress-Restriktion ([ADR-0012](0012-devcontainer-egress-firewall.md)) daneben besteht und (d) ob das Default-Profil (`nestedRuntime: none`) davon berührt ist.

**Messung (2026-09-30).** Probe-Image: `mcr.microsoft.com/devcontainers/base:debian` (Digest `sha256:1f851004…4286`) plus Debian-Pakete `podman` 5.4.2, `uidmap`, `fuse-overlayfs`, `slirp4netns`; Benutzer `vscode` (subuid/subgid `100000:65536` bereits im Base-Image vorhanden — ein zweiter Eintrag erzeugt überlappende Bereiche und `EINVAL`). Prüfung: `podman run --rm --network=host alpine echo …` im Container. Host: Ubuntu 24.04, Kernel 6.8, Docker 29.8.1 (rootful, cgroup v2, AppArmor `docker-default`, Seccomp builtin, `apparmor_restrict_unprivileged_userns=0`).

| Docker-Optionen (zusätzlich zum Default) | Ergebnis |
|---|---|
| keine | `cannot clone: Operation not permitted` |
| `--device /dev/fuse` | wie oben |
| + `seccomp=unconfined` | `newuidmap: write to uid_map failed: EPERM` |
| + `seccomp=unconfined` + `apparmor=unconfined` | wie oben |
| + `seccomp=unconfined` + `--cap-add NET_ADMIN` | wie oben |
| + `seccomp=unconfined` + `--cap-add SYS_ADMIN` (AppArmor default) | `mount overlay: permission denied` |
| `--cap-add ALL` + Seccomp/AppArmor unconfined | `crun: mount proc: Operation not permitted` |
| **`--cap-add SYS_ADMIN` + `seccomp=unconfined` + `apparmor=unconfined` + `systempaths=unconfined` + `--device /dev/fuse`** | **läuft** (`overlay`-Storage, Image-Pull, `run`) |
| dasselbe + `--cap-add NET_ADMIN` | **läuft** |
| dasselbe ohne `/dev/fuse` | `fuse-overlayfs: cannot mount` (harter Fehler); mit `--storage-driver vfs` **läuft** |
| `--privileged` | läuft |
| Rootless Podman als Host-Engine (4.9.3), Default / `--device /dev/fuse` / `label=disable` | `newuidmap … EPERM` (mit diesen Optionen nicht lauffähig; weitere Flags ungeprüft) |

Zusatzbefund: Ohne `--network=host` fehlt im Probe-Image `pasta` („could not find pasta"); das Paket `passt` gehört ins Image.

## Entscheidung

1. **Base-Image bleibt `mcr.microsoft.com/devcontainers/base:debian`**; Podman und Zubehör (`podman`, `uidmap`, `fuse-overlayfs`, `passt`) kommen als Distro-Pakete im generierten Dockerfile, nicht als Devcontainer-Feature. Damit greifen [`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) nicht zusätzlich. Bestehende subuid-/subgid-Einträge des Benutzers werden nicht dupliziert.
2. **Lockerungsprofil unter Docker** ist bei `nestedRuntime: podman` genau: `--cap-add SYS_ADMIN`, `seccomp=unconfined`, `apparmor=unconfined`, `systempaths=unconfined`, `--device /dev/fuse`. Kein `--privileged`. Jedes Element wird in der Befehlsausgabe einzeln ausgewiesen ([`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)).
3. **Default-Profil bleibt vertretbar** (keine Teilrevision des Spec): `nestedRuntime` steht auf `none`; die Lockerungen entstehen nur bei ausdrücklich gesetztem `podman`. Die Doku benennt, dass das Profil dann nahe an `--privileged` liegt.
4. **Storage:** Das Startscript wählt `vfs`, wenn `/dev/fuse` fehlt (Fallback aus der Degradationstabelle von [`LH-FA-DEV-007`](../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); `overlay` + `fuse-overlayfs` sonst. Ein blockierter User-Namespace ist Exit `11`, kein Fallback.
5. **Egress ([ADR-0012](0012-devcontainer-egress-firewall.md)):** `NET_ADMIN` besteht neben dem Profil. Weil `SYS_ADMIN` im Container vorhanden ist, kann ein Prozess die Egress-Regeln aufheben; die Restriktion ist mit `nestedRuntime: podman` nur noch ein Guardrail und wird so dokumentiert.
6. **Podman als Host-Engine und Colima/macOS:** ungeprüft; Nachholmessung durch den Projektinhaber auf der jeweiligen Umgebung. Rootless Podman als Host-Engine scheitert mit den getesteten Standard-Optionen; für Colima/macOS liegt keine Messung vor. Beides blockiert `Accepted` nicht: bis zur Nachholmessung gilt es als dokumentierte Einschränkung; das Ergebnis wird in der Geschichte dieses ADR nachgetragen.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Kein nested Runtime (nur Socket-Mount) | keine Lockerung | Docker-Socket = Host-Root-Äquivalent; widerspricht dem Sandbox-Ziel |
| B — `--privileged` | einfachste Konfiguration, läuft | maximale Lockerung, nicht ausweisbar fein; verstößt gegen [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) |
| C — `quay.io/podman/stable` (Fedora) als Base | für Podman gebaut | Bruch mit dem devcontainers-Ökosystem, ungetestet, zweite Distro im Template |
| D — Sidecar-Container mit Podman/Docker-in-Docker | Isolation vom Arbeitscontainer | Compose-Komplexität, braucht selbst `--privileged` |
| **E — Debian-Base + Distro-Podman, gemessenes Minimalprofil** | gemessen lauffähig, feiner als `--privileged`, kein Feature | `SYS_ADMIN` + drei `unconfined`-Optionen sind kaum weniger als `--privileged` |

## Konsequenzen

- Positiv: Image-Builds ohne Host-Socket; Lockerungen sind belegt und einzeln ausweisbar; das Default-Profil bleibt ungelockert.
- Negativ: Mit `nestedRuntime: podman` ist die Sandbox-Grenze deutlich schwächer; Egress-Guardrail hebelbar; rootless Podman als Host-Engine ungelöst.
- Folgepflicht: Umsetzungs-Slice erzeugt die Optionen in `devcontainer.json` (`runArgs`) und weist sie aus; Integrationstest (`//go:build docker`) wiederholt die Messung als Regression.
- **Nicht verifiziert (Stand 2026-09-30):** macOS+Colima und Podman als Host-Engine. Die Messung lief auf einem Linux-Host mit rootful Docker; die Umgebungen stehen dem Projektinhaber nur nacheinander zur Verfügung. Die Aussagen dieses ADR gelten bis zur Nachholmessung (Abschnitt „Nachholmessung“) nur für Docker unter Linux; die Doku darf für Colima/macOS bis dahin keine Lauffähigkeit zusagen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Integrationstest (`//go:build docker`) | Sandbox-Profil mit `nestedRuntime: podman` führt `podman run` im Container aus; `nestedRuntime: none` erzeugt keine der Lockerungen | `make test-docker` |
| Golden Case | erzeugte `devcontainer.json` enthält die Lockerungen nur bei `podman` | `make test` |

## Re-Evaluierungs-Trigger

Neue Docker-/Podman-/Kernel-Version, die eine kleinere Lockerung erlaubt (z. B. Seccomp-Profil statt `unconfined`, Docker-Userns-Remap), oder ein Nutzerbericht bzw. CI-Befund von Colima/macOS/Podman-Host, der das Profil ändert.

## Nachholmessung (macOS+Colima, Podman-Host)

Auf der jeweiligen Umgebung ausführen, Ergebnis (Zeile je Variante, Engine-/Colima-Version, VM-Typ, Kernel) in die Geschichte eintragen. `ENGINE` ist `docker` (unter Colima der Docker-Client gegen den Colima-Socket) bzw. `podman`.

```bash
mkdir -p /tmp/nested && cd /tmp/nested
cat > Dockerfile <<'EOT'
FROM mcr.microsoft.com/devcontainers/base:debian
RUN apt-get update && apt-get install -y --no-install-recommends podman uidmap fuse-overlayfs passt ca-certificates && rm -rf /var/lib/apt/lists/*
USER vscode
EOT
ENGINE=docker   # oder podman
$ENGINE build -t nested-probe .
P='podman run --rm --network=host docker.io/library/alpine:3 echo NESTED_OK'
$ENGINE run --rm nested-probe bash -c "$P"                                   # A: Default
$ENGINE run --rm --device /dev/fuse --cap-add SYS_ADMIN \
  --security-opt seccomp=unconfined --security-opt apparmor=unconfined \
  --security-opt systempaths=unconfined nested-probe bash -c "$P"            # B: Profil aus der Entscheidung
$ENGINE run --rm --cap-add SYS_ADMIN --security-opt seccomp=unconfined \
  --security-opt apparmor=unconfined --security-opt systempaths=unconfined \
  nested-probe bash -c "podman --storage-driver vfs run --rm --network=host docker.io/library/alpine:3 echo VFS_OK"  # C: ohne /dev/fuse
```

Erwartung nach Linux/Docker-Befund: A scheitert, B und C laufen. Unter Colima (Linux-VM) ist dasselbe Ergebnis zu erwarten; abweichende Ergebnisse ändern die Entscheidung. Unter `podman` als Host-Engine ist B zusätzlich mit `--userns=keep-id` und `--security-opt label=disable` zu versuchen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-30 | Entwurf mit Messreihe (Docker 29.8.1, Ubuntu 24.04) | Umsetzungs-Slice V1 Sandbox-Profil |
| 2026-09-30 | `Accepted` durch den Projektinhaber; macOS/Colima und Podman-Host bleiben bis zur Nachholmessung ungeprüft (§Nachholmessung) | Vereinbarung mit dem Projektinhaber |
