# Devcontainer-Sandbox-Profil — u-boot

| Dokument    | Sandbox-Profil-User-Guide |
| ----------- | ------------------------- |
| Projektname | `u-boot` |
| Bezug       | [LH-FA-DEV-004](../../spec/lastenheft.md#lh-fa-dev-004--benutzerrechte), [LH-FA-DEV-006](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [LH-FA-DEV-007](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [LH-FA-DEV-009](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) |
| ADR         | [ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md) (Messung, Lockerungsprofil) |
| Status      | Lastenheft 0.3.1 (noch nicht released) |

## Zweck

Das Sandbox-Profil erzeugt einen Devcontainer für den Einsatz autonomer
Agenten (ohne Rückfrage): Schadensbegrenzung, **keine** harte
Isolationsgrenze. u-boot startet keinen Agenten und setzt keinen
Berechtigungsmodus.

## 1. Aktivieren

```bash
u-boot init my-service --devcontainer --sandbox     # neues Projekt
u-boot generate devcontainer --sandbox              # bestehendes Projekt
u-boot config set devcontainer.profile sandbox      # alternativ, dann generate
```

`--sandbox` setzt `devcontainer.profile: sandbox` in `u-boot.yaml`. Ohne
`--devcontainer` (bei `init`) bzw. für andere Artefakte als `devcontainer`
(bei `generate`) ist das Flag ein Fehler.

## 2. Was das Profil erzeugt

| Eigenschaft | Sandbox-Profil | Default-Profil |
| --- | --- | --- |
| Workspace | benanntes Volume `<projekt>-workspace-${devcontainerId}` (pro Ordner eindeutig), Repo wird im Container geklont; `workspaceMount` ist leer | Bind-Mount des Projektordners |
| Docker-Socket des Hosts | nicht eingebunden | nicht eingebunden |
| `--privileged` / zusätzliche Capabilities | nein (außer bei `nestedRuntime: podman`, §4) | nein |
| Host-Dateien mit Geheimnissen (`~/.ssh`, `~/.aws`, …) | nicht eingebunden | nicht eingebunden |
| Benutzer | `vscode`, nicht root | `vscode`, nicht root |

**Clone-Quelle:** `devcontainer.sandbox.repository`, sonst die URL des
Remotes `origin` im Projekt-Repository. Der Clone läuft beim ersten Start
(`postCreateCommand`), nur wenn der Workspace noch leer ist.

Um im Container ein **anderes Repository** zu klonen (und daraus zu pullen
und zu pushen), setze die Quelle explizit; sie gilt auch ohne `origin`:

```bash
u-boot config set devcontainer.sandbox.repository git@github.com:org/anderes-repo.git
u-boot generate devcontainer
```

Dieselben Regeln gelten für die URL (keine Zugangsdaten, keine unsicheren
Zeichen, sonst Exit `10`). Pull/Push im Container laufen gegen das
geklonte Repository (dessen eigenes `origin`).

- **Weder `repository` noch Remote** (z. B. direkt nach `init`): kein Clone-Schritt, Warnung
  [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil). Nach `git remote add origin …` (oder `config set devcontainer.sandbox.repository …`) ergänzt
  `u-boot generate devcontainer` den Schritt.
- **URL mit Zugangsdaten oder unsicheren Zeichen** (`https://token@…`,
  `user:pw@`, `;`, `$(…)`): Fehler, Exit-Code `10`, es wird nichts
  geschrieben. Erlaubt sind `git@host:pfad`, `ssh://git@host/pfad` und
  `https://` ohne Benutzerinfo.
- Git-Worktrees und Submodule (`.git` als Datei) gelten als „kein Remote“.

## 3. UID des Container-Benutzers

```bash
u-boot config set devcontainer.user.uid 501      # z. B. macOS unter Colima
u-boot generate devcontainer
```

Zulässig sind `1` bis `65535` (`0` = root wird abgelehnt, Exit `10`).
Default `1000`: dann bleibt die erzeugte Ausgabe unverändert. Bei einem
anderen Wert erzeugt u-boot das Build-Argument `USER_UID` und einen
`usermod`-Schritt im Dockerfile (gilt für beide Profile).

## 4. Nested Podman (optional)

```bash
u-boot config set devcontainer.sandbox.nestedRuntime podman
u-boot generate devcontainer
```

Damit läuft im Container rootless Podman (mit `docker`-Alias), sodass
`docker build`/`docker run` ohne Host-Socket möglich sind. Dafür sind unter
Docker folgende Lockerungen nötig; **jede wird in der Befehlsausgabe
einzeln als Warnung [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) ausgewiesen**:

```text
--cap-add=SYS_ADMIN
--security-opt=seccomp=unconfined
--security-opt=apparmor=unconfined
--security-opt=systempaths=unconfined
--device=/dev/fuse
```

**Das ist nahe an `--privileged`.** Die Sandbox ist damit deutlich
schwächer; das Default-Profil (`nestedRuntime: none`) bleibt ungelockert.
Die Messung und die Alternativen stehen in
[ADR-0014](../plan/adr/0014-nested-podman-sandbox-devcontainer.md).

### Degradation: `devcontainer.sandbox.onUnavailable`

`warn` (Default) | `fail`. Das Startscript `.devcontainer/sandbox-init.sh`
(per `COPY` im Image) prüft beim Containerstart:

| Zustand | `warn` | `fail` |
| --- | --- | --- |
| `/dev/fuse` fehlt | Fallback auf `vfs`-Storage, Warnung | Exit `11` |
| Nested User-Namespaces blockiert (Seccomp/AppArmor) | Exit `11` | Exit `11` |

`u-boot doctor` (`devcontainer.sandbox.runtime`, [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) prüft `/dev/fuse` nur auf
Linux-Hosts; unter macOS/Colima läuft die Engine in einer VM und der Host
kann das nicht beurteilen (dort prüft das Startscript).

## 5. Git-Zugangsdaten ([LH-FA-DEV-009](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer))

Zugangsdaten stehen **nie** im Image, im Volume, in `u-boot.yaml` oder einer
erzeugten Datei. Übergabe zur Laufzeit:

- `devcontainer.json` enthält `remoteEnv: { "GIT_TOKEN": "${localEnv:GIT_TOKEN}" }`
  (nur die Referenz). Setze `GIT_TOKEN` in der Host-Umgebung, bevor du den
  Container startest.
- Ein System-Credential-Helper im Image liefert `$GIT_TOKEN` an `git` (nur
  für https-Remotes; bei `ssh` gibt es keinen Schlüssel im Container).

Empfehlung: **kurzlebige, auf das Repository begrenzte Tokens** (Private
Key nie im Container) und **Branch-Protection** auf dem Remote, damit ein
kompromittierter Agent nicht direkt auf geschützte Branches schreiben kann.

`u-boot doctor` (`devcontainer.sandbox.credentials`, nur `warn`) meldet
mögliche Klartext-Tokens in `u-boot.yaml`, `compose.yaml`, `.env.example`
und `.devcontainer/*` sowie eine fehlende Token-Quelle bei https-Clone.

## 6. Mehrere Instanzen ([ADR-0015](../plan/adr/0015-sandbox-volumes-pro-instanz.md))

Die Volumes heißen `<projekt>-workspace-${devcontainerId}` und (mit nested
Podman) `<projekt>-containers-${devcontainerId}`. `devcontainerId` löst das
Dev-Containers-Werkzeug pro Projektordner auf; das Workspace-Volume steht in
`mounts` (nicht in `workspaceMount`, dort wird die Variable nicht aufgelöst;
`workspaceMount` bleibt leer und schaltet den Bind-Mount ab).

- **Parallele Instanzen desselben Projekts:** aus **getrennten Ordnern**
  (Klone oder `git worktree`); jede Instanz hat eigene Volumes.
- **Derselbe Ordner** bleibt eine Instanz (gleiche `devcontainerId`).
- **Compose-Service-Ports** (`u-boot up`) mehrerer Instanzen kollidieren
  weiterhin; das ist ein eigener Vorgang.
- Geprüft mit der `devcontainer`-CLI (0.80.2) und Docker auf Linux;
  **VS Code und Codespaces ungeprüft**.
- Alte Volumes `<projekt>-workspace` / `<projekt>-containers` aus früheren
  Entwicklungsständen bleiben verwaist (`docker volume rm`).

## 7. Egress-Restriktion ([LH-FA-DEV-008](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [ADR-0012](../plan/adr/0012-devcontainer-egress-firewall.md))

```bash
u-boot config set devcontainer.sandbox.egress.enabled true
u-boot config set devcontainer.sandbox.egress.allow api.anthropic.com,example.org
u-boot generate devcontainer
```

Mechanismus (DNS-gesteuert, **Guardrail, keine Sandbox-Grenze**): Beim
Containerstart (`postStartCommand`, `sudo`) startet `.devcontainer/egress-init.sh`
einen lokalen `dnsmasq`, der **nur erlaubte Namen** auflöst und die Antworten
in ein nftables-Set einträgt; nftables lässt nur Loopback, bestehende
Verbindungen, DNS zu den Upstream-Resolvern und Ziele im Set zu. Alles andere
(IPv6 komplett) wird verworfen; nicht erlaubte Namen sind nicht auflösbar.
`--cap-add=NET_ADMIN` wird als Lockerung ausgewiesen.

**Default-Allowlist:**

| Quelle | Hosts |
| --- | --- |
| immer | `github.com`, `api.github.com`, `codeload.github.com`, `raw.githubusercontent.com`, `objects.githubusercontent.com`, `deb.debian.org`, `security.debian.org` |
| `nestedRuntime: podman` | `registry-1.docker.io`, `auth.docker.io`, `production.cloudflare.docker.com`, `ghcr.io` |
| Feature `node` | `registry.npmjs.org` |
| Feature `go` | `proxy.golang.org`, `sum.golang.org`, `storage.googleapis.com` |
| Feature `java` | `repo.maven.apache.org`, `repo1.maven.org` |
| Clone-Quelle | der Host von `devcontainer.sandbox.repository` bzw. `origin` |
| Nutzer | `devcontainer.sandbox.egress.allow` (kleingeschriebene Hostnamen, ohne Schema/Port/Wildcard) |

**Wichtig:** Hosts des Agenten selbst (z. B. dessen API) kennt u-boot nicht —
sie gehören in `egress.allow`. Subdomains eines erlaubten Namens sind
mitgemeint. Die Allowlist steht im Script; nach einer Änderung
`u-boot generate devcontainer` und Container neu starten.

**Grenzen:** Wer die Capability hat (und das `sudo` des Base-Images
nutzen kann), kann die Regeln aufheben. Der Verkehr verschachtelter
Podman-Container läuft über `FORWARD` und wird nicht erfasst.
**Degradation:** Ohne `NET_ADMIN` entfällt die Restriktion mit Warnung
(`onUnavailable: warn`) oder das Startscript endet mit Exit `11` (`fail`);
`u-boot doctor` (`devcontainer.sandbox.egress`) prüft nur die Konfiguration
(ohne Sandbox-Profil wirkungslos → `warn`).

## 8. Grenzen

- Geprüft unter Docker auf Linux (siehe ADR). **Nicht geprüft:**
  macOS/Colima und Podman als Host-Engine.
- Wechsel von `podman` zurück auf `none` lässt eine vorhandene
  `sandbox-init.sh` liegen (vom Dockerfile nicht mehr referenziert).
