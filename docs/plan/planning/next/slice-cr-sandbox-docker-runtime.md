# Slice CR: Eigener Docker-Daemon im Sandbox-Devcontainer (`nestedRuntime: docker`)

> **CR-Vehikel.** u-boot führt kein separates CR-Artefakt — dieser Slice
> **ist** der Change Request am Vertrags-Stratum (Precedent:
> [`done/slice-cr-sandbox-devcontainer`](../done/slice-cr-sandbox-devcontainer.md)).
> Anlass: Repos, deren Gates auf BuildKit und einen echten Docker-Daemon
> gebaut sind (Stufen-Builds mit `-o type=cacheonly`, `RUN --network=none`,
> privilegierte dind-Sidecars, `docker compose`), sollen im Sandbox-Profil
> laufen — parallel in mehreren Instanzen, ohne Host-Socket. Rootless Podman
> ([ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md)) deckt
> diese Gates nicht ab.

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** Spec-Wartung: Docker als verschachtelte Runtime (Lastenheft 0.5.0);
CR am Vertrags-Stratum.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Änderung),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
[ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), §16
(Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der
CR nicht; Folgeänderungen dort entstehen in den Umsetzungs-Slices.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Lastenheft-Änderung **0.4.0 → 0.5.0**: `devcontainer.sandbox.nestedRuntime`
bekommt den dritten Wert `docker` — ein eigener Docker-Daemon im
Sandbox-Container. Danach Umsetzung samt Folge-ADR.

**Warum nicht Podman:** Die Zielrepos nutzen Funktionen, die rootless Podman
im Container nicht oder anders liefert — Build-Ausgabe `type=cacheonly`
(Buildah kennt sie nicht; in der ADR-Messung zu belegen),
`docker run --privileged` für eigene dind-Sidecars (rootless nicht möglich),
`docker buildx imagetools`, `docker save` in einen fremden Daemon. Diese Repos
auf Podman umzustellen hieße, ihre Gate-Kette umzubauen; der Bedarf ist ein
Docker-Daemon, nicht irgendeine Runtime.

**Warum nicht Socket-Mount:** ein geteilter Host-Daemon bedeutet geteilte
Image-, Container-, Netz- und Volume-Namen. Parallele Instanzen desselben
Repos kollidieren dort; außerdem ist der Socket Host-Root-Äquivalent
([ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md), Alternative A).

### CR-Kern (vorgeschlagene Fassung von `LH-FA-DEV-007`)

> ### [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) – Container-Runtime im Sandbox-Devcontainer
>
> *Priorität: V1.* Im Sandbox-Profil ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) soll optional eine
> Container-Runtime im Container bereitgestellt werden, damit Image-Builds
> und Container-Läufe ohne Zugriff auf einen Host-Socket möglich sind.
>
> - Konfigurationsschlüssel `devcontainer.sandbox.nestedRuntime`
>   (`podman` | `docker` | `none`, Default `none`).
> - Bei `podman`: rootless Podman mit `docker`-Kompatibilität im Container.
> - Bei `docker`: ein eigener Docker-Daemon im Container mit BuildKit sowie
>   Buildx- und Compose-Plugin. Die Docker-CLI im Container spricht
>   ausschließlich mit diesem Daemon, sein Socket ist dem Container-Benutzer
>   zugänglich. Der Daemon-Speicher liegt in einem benannten Volume je
>   Instanz. `docker` verlangt `--privileged`; die Befehlsausgabe weist das
>   als Lockerung aus, mit dem Hinweis, dass das Profil damit keine
>   Isolationsgrenze gegenüber dem Container-Host mehr bietet.
> - Bei `podman` ist das Ergebnis engine-neutral und startet unter Docker
>   (inkl. Colima) und Podman. Bei `docker` startet es unter Docker
>   (inkl. Colima und Docker Desktop); Podman als Host-Engine ist für
>   `docker` nicht zugesagt.
>
> **Degradation:** `devcontainer.sandbox.onUnavailable` (`warn` | `fail`,
> Default `warn`) steuert den Umgang mit fehlenden Fähigkeiten: Bei `warn`
> führen ein fehlendes `/dev/fuse` zu einem Fallback auf `vfs`-Storage, eine
> nicht gewährbare Egress-Capability zum Wegfall der Restriktion und ein
> Docker-Daemon, der beim Containerstart nicht betriebsbereit wird, zu einem
> Container ohne Daemon, jeweils mit Warnung; bei `fail` sind alle drei ein
> Umgebungsproblem (Exit-Code `11`). Blockierte verschachtelte
> User-Namespaces (Seccomp/AppArmor) bei `nestedRuntime: podman` sind in
> beiden Modi ein Umgebungsproblem (Exit-Code `11`). Eine ausdrücklich
> angeforderte Runtime wird nie stillschweigend ersetzt; jeder Fallback wird
> in der Befehlsausgabe und in `u-boot doctor` ausgewiesen. Ungültige Werte
> der Schlüssel führen zu einem fachlichen Validierungsfehler (Exit-Code `10`).

**Unverändert, mit Begründung:**

- [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Die Klausel „kein `--privileged` … sofern
  nicht durch [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) … ausdrücklich verlangt" trägt den neuen Wert
  bereits; die Ausweispflicht für Lockerungen gilt unverändert.
- [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion): „Prozesse mit der nötigen Capability können sie
  aufheben" deckt `--privileged` ab. Die Default-Allowlist ist laut Text
  *dokumentiert*, nicht spezifiziert — die Registry-Hosts für `docker`
  ergänzt die Doku, nicht das Lastenheft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung** — sie liegt in den drei Umsetzungs-Slices (§2); dieser Slice
  liefert Vertrag und Entscheidung, nicht Code.

- **Kein rootless dockerd** — braucht laut Docker-Doku im Container
  ebenfalls `--privileged` und gewinnt damit keine Lockerung; die ADR darf es
  messen, der Vertrag sagt es nicht zu.
- **Kein geteilter Image-Cache zwischen Instanzen** — getrennte Daemons sind
  der Zweck; Speicherverbrauch je Instanz wird dokumentiert, nicht gelöst.
- **Keine Isolationszusage bei `docker`** — Schadensbegrenzung entfällt
  gegenüber dem Container-Host weitgehend; das Profil begrenzt weiterhin
  Workspace, Secrets und Socket des Hosts.
- **Keine Änderung an `podman` und `none`.**

**Abnahmekriterien (für die Umsetzung, nicht für den Spec-Change):**

1. `nestedRuntime: docker` erzeugt `--privileged` in `runArgs` und weist es
   als Warnung aus; kein Socket-Mount, kein Bind-Mount des Host-Ordners,
   keine Host-Secrets. Daemon-Speicher im Volume
   `<projekt>-docker-${devcontainerId}` (Benennung nach [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md)).
2. Im Container als `vscode`, ohne `sudo`: `docker info` nennt den inneren
   Daemon; `docker buildx build --target <stufe> -o type=cacheonly .` läuft;
   ein Dockerfile mit `RUN --network=none` baut; `docker compose up` läuft;
   `docker run --rm --privileged docker:dind` startet (dritte Ebene).
3. Zwei Instanzen aus zwei Ordnern (Klon oder `git worktree`) haben getrennte
   Daemons: gleichnamige Images, Container, Netze und Volumes kollidieren nicht.
4. Golden Cases: `none` und `podman` erzeugen dieselben Dateien wie vorher;
   nur `docker` erzeugt `--privileged`.
5. Degradation: ein Start ohne `--privileged` endet bei `warn` mit Warnung
   und ohne Daemon, bei `fail` mit Exit-Code `11`; nie mit Podman als Ersatz.
6. Integrationstest (`//go:build docker`) wiederholt 2 und 3 auf Linux;
   Nachholmessung auf macOS mit Docker Desktop und mit Colima, Ergebnis in die
   Geschichte der Folge-ADR.
7. `u-boot doctor`: `nestedRuntime=docker` ohne Sandbox-Profil warnt
   (wirkungslos, wie bei `podman`); mit Egress-Restriktion warnt er, dass der
   Verkehr verschachtelter Container und die iptables-Regeln des Daemons die
   Leitplanke unterlaufen.

**Historie-Zeile (§16, bei Spec-Change zu setzen):**

```text
| 0.5.0 | (Annahmedatum) | LH-FA-DEV-007: dritter Wert `docker` für
`devcontainer.sandbox.nestedRuntime` (eigener Docker-Daemon im Container,
`--privileged` als ausgewiesene Lockerung, Speicher je Instanz,
Degradationsfall „Daemon nicht betriebsbereit"). | Vereinbarung mit dem
Projektinhaber |
```

## 2. Definition of Done

- [ ] **Lastenheft 0.5.0**: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) in der Fassung oben,
  Versionsfeld, Historie-Zeile.
- [ ] **Folge-ADR** (nächste freie Nummer) entscheidet die Umsetzung (§3),
  mit Messung wie [ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md); [ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md) bleibt für `podman` gültig und wird
  nicht abgelöst — sein Contra zu Alternative B („verstößt gegen DEV-006")
  gilt für `docker` nicht mehr, weil DEV-007 es nun ausdrücklich verlangt;
  die neue ADR nennt das. ADR-Index aktualisiert.
- [x] **Umsetzungs-Slice(s)** mit Kennung angelegt, die die Abnahmekriterien
  1–7 tragen: [`slice-v1-sandbox-docker-generator`](../open/slice-v1-sandbox-docker-generator.md) (1, 4),
  [`slice-v1-sandbox-docker-degradation-doctor`](../open/slice-v1-sandbox-docker-degradation-doctor.md) (5, 7),
  [`slice-v1-sandbox-docker-integrationstest`](../open/slice-v1-sandbox-docker-integrationstest.md) (2, 3, 6).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Doku-Update: öffentliche Verträge nachgezogen (CHANGELOG, ADR-Index,
      Roadmap).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register: keine Beobachtung angefallen, oder fortgeschrieben (§7).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Versionsfeld 0.5.0, Historie-Zeile |
| `docs/plan/adr/` (Folge-ADR, Index) | neu / update | Umsetzungsentscheidung samt Messung |
| `CHANGELOG.md` | update | öffentlicher Vertrag berührt |

**Fragen an die Folge-ADR:**

- **Herkunft des Daemons:** Debian-Pakete (`docker.io`, `docker-buildx`,
  `docker-compose`) · `docker-ce` aus dem Docker-APT-Repo (zusätzliche Quelle,
  Egress) · Devcontainer-Feature `docker-in-docker` (externe Quelle →
  Allowlist nach [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)).
- **Start:** Entrypoint-Skript vs. `postStartCommand` mit `sudo`;
  Bereitschaftsprüfung mit Frist (Grundlage für die Degradation).
- **Socket-Zugriff:** Gruppe `docker` für `vscode` vs. Socket-Rechte.
- **Storage-Driver** auf dem Volume (`overlay2`), cgroup v2 verschachtelt
  unter Docker Desktop und Colima.
- **Egress-Allowlist** für `docker`: dieselben Registry-Hosts wie bei
  `podman`, plus ggf. `download.docker.com`.

## 4. Trigger

**Start** (`next` → `in-progress`): Projektinhaber gibt den CR frei; die
Zustimmung zu `--privileged` liegt vor (§6, Punkt 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zeigt,
  dass Abnahmekriterium 2 (dritte Ebene `--privileged`) unter Docker Desktop
  nicht erreichbar ist — dann CR auf die Ebenen 1–2 schneiden und die dritte
  als eigenen Punkt führen.
- `in-progress` → `open` (blockiert): der Projektinhaber nimmt die Zustimmung
  zu `--privileged` im Sandbox-Profil zurück.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7). Lastenheft 0.5.0 und Folge-ADR
sind angenommen; die drei Umsetzungs-Slices liegen in `open/` oder weiter.

## 6. Risiken und offene Punkte

1. **`--privileged` im Sandbox-Profil überhaupt?** Der Projektinhaber (pt9912)
   stimmt `--privileged` für `nestedRuntime: docker` zu (2026-10-14).
   — **Ausgang:** entfallen: Zustimmung liegt vor.
2. **Container-Ausbruch:** mit `--privileged` trivial. Unter macOS trifft er
   die VM von Docker Desktop/Colima, unter Linux den Host. Die Doku muss das
   so deutlich sagen wie die Podman-Warnung.
   — **Ausgang:** weiter offen: Doku im Slice
   [`slice-v1-sandbox-docker-integrationstest`](../open/slice-v1-sandbox-docker-integrationstest.md).
3. **Speicher und Kaltstart je Instanz:** jeder Daemon zieht und baut alles
   neu; mehrere GB je Instanz.
   — **Ausgang:** weiter offen: Doku im Slice
   [`slice-v1-sandbox-docker-integrationstest`](../open/slice-v1-sandbox-docker-integrationstest.md).
4. **macOS, Docker Desktop und VS Code ungeprüft** (wie bei [ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md)) — die
   Zielumgebung des Anlasses; Zusage erst nach Nachholmessung.
   — **Ausgang:** weiter offen: Nachholmessung im Slice
   [`slice-v1-sandbox-docker-integrationstest`](../open/slice-v1-sandbox-docker-integrationstest.md).
5. **Egress-Leitplanke** wird durch verschachtelte Container und
   Daemon-iptables unterlaufen.
   — **Ausgang:** weiter offen: `doctor`-Warnung im Slice
   [`slice-v1-sandbox-docker-degradation-doctor`](../open/slice-v1-sandbox-docker-degradation-doctor.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** [`slice-v1-sandbox-docker-generator`](../open/slice-v1-sandbox-docker-generator.md), [`slice-v1-sandbox-docker-degradation-doctor`](../open/slice-v1-sandbox-docker-degradation-doctor.md), [`slice-v1-sandbox-docker-integrationstest`](../open/slice-v1-sandbox-docker-integrationstest.md)
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `spec/`, `harness/`,
`docs/plan/` (eine Sub-Area laut `harness/conventions.md`, Modus Greenfield).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für diese Sub-Area.

Alle berührten Sub-Areas GF.
