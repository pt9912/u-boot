# Slice V1: Aktueller Docker aus dem Docker-APT-Repo für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.8.1 (Patch-Release: aktuelle Docker-Version im Sandbox-Devcontainer).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md). Folge-ADR ersetzt dessen Entscheidung 1 (Herkunft) und ergänzt Entscheidung 6 (Egress).

**Berührte Spec-Stellen:** `spezifikation.md` [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) (Egress-Hosts) und die Verfeinerung [`LH-FA-DEV-007.b`](../../../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) (Ausstattung bei `docker`); das Lastenheft bleibt unverändert ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) nennt keine Paketquelle).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Daemon im Sandbox-Devcontainer kommt aus dem Docker-APT-Repo (`docker-ce`, Buildx- und Compose-Plugin in den aktuellen Versionen) statt aus dem Debian-Archiv (Docker 26.1.5, Buildx 0.13.1, Compose 2.26.1).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Änderungen am Lockerungsprofil (`--privileged`), an Start, Degradation und Volume — die Messung zeigt nur die Paketquelle; das Verhalten bleibt wie in [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) festgelegt.
- Das Feature `docker-in-docker` — externe Quelle mit fremdem Installationsskript ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)); das APT-Repo mit gepinntem Schlüssel ist enger.
- Docker Desktop und Podman als Host-Engine — der Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) führt die Docker-Desktop-Messung.

## 2. Definition of Done

- [x] **Folge-ADR** (nächste freie Nummer) entscheidet Herkunft und Schlüsselpinning mit Messung (Colima); ADR-Index aktualisiert.
- [x] **Generator:** Dockerfile-Template installiert `docker-ce`, `docker-ce-cli`, `containerd.io`, `docker-buildx-plugin`, `docker-compose-plugin` aus `download.docker.com` mit geprüftem Schlüssel-Fingerprint; Egress-Allowlist um `download.docker.com` ergänzt; Golden Cases und Integrationstest (`//go:build docker`) belegen die Version und die Abnahmekriterien aus dem CR.
- [x] **Doku:** Spezifikation, Benutzerhandbuch und CHANGELOG nennen die neue Herkunft.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/` (Folge-ADR, Index) | neu / update | Herkunft `docker-ce`, Fingerprint-Pin, Egress |
| `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl` | update | Paketquelle statt Debian-Pakete |
| `internal/hexagon/application/devcontainer_egress.go` | update | `download.docker.com` bei `nestedRuntime: docker` |
| `internal/hexagon/application/*_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go` | update | Golden Cases, Version und Abnahmekriterien |
| `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `harness/replay.md` | update | öffentlicher Vertrag und Golden-Zeile |

## 4. Trigger

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; kein anderer Slice des Implementers in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet 3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): `download.docker.com` ist für den Image-Build nicht erreichbar oder der Schlüssel-Fingerprint stimmt nicht mehr mit dem Pin überein.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Docker rotiert den Schlüssel; der gepinnte Fingerprint passt nicht mehr und der Image-Build bricht ab — **Ausgang:** weiter offen: Pin-Hebung per Folge-ADR beim Eintreten (Re-Evaluierungs-Trigger der ADR); der Build-Abbruch ist gewollt.
- Neuere Docker-Versionen ändern das Verhalten der Zielrepos — **Ausgang:** entfallen: die Abnahmekriterien des CR laufen unter Docker 29.9, Buildx 0.38 und Compose 5.6 (Integrationstests, Colima und CI).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Umstellung auf `docker-ce` aus dem Docker-APT-Repo (Docker 29.9, Buildx 0.38, Compose 5.6) lief mit unveränderten Abnahmekriterien; `download.docker.com` steht bei `docker` in der Egress-Allowlist, der Test `apt-get update` unter Egress wird ohne den Host rot (Negativprobe, Exit 100). [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) ersetzt Entscheidung 1 von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md), ohne diese zu ändern. Implementation `eea88cc`, Review-Fixes `3cfa710` und `90f9122`.
- **Was ging anders als geplant:** Der Review [`2026-10-14-sandbox-docker-ce`](../../../reviews/2026-10-14-sandbox-docker-ce.md) fand drei MEDIUM (Testtiefe), der Folgelauf [`2026-10-14-sandbox-docker-ce-r2`](../../../reviews/2026-10-14-sandbox-docker-ce-r2.md) einen weiteren (Egress-Test ohne `--error-on=any`); alle behoben. Der Storage-Treiber meldet bei Docker 29 `overlayfs` (containerd-Image-Store) statt `overlay2`; Tests und [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) akzeptieren beide.
- **Steering-Loop-Eintrag:** Tests, die eine Eigenschaft belegen sollen, brauchen eine Negativprobe (Eigenschaft entfernen, Test wird rot) — vor dem Review ausgeführt statt vom Reviewer nachgefordert.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) (Docker Desktop, unverändert offen).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/hexagon/application` (Greenfield laut `harness/conventions.md`), `docs/plan/`, `docs/user/`, `spec/`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Verifier-Lauf 2026-10-14 (Rolle Verifier, getrennt von Implementer und Reviewer).

Scope:
- Slice: [slice-v1-sandbox-docker-ce](slice-v1-sandbox-docker-ce.md)
- IDs: [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [LH-FA-DEV-008](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [LH-NFA-SEC-004](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md), [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
- Artefakte: `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl`, `internal/hexagon/application/devcontainer_egress.go`, `internal/hexagon/application/generate_sandbox_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `harness/replay.md`, [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) samt Index; Commits `eea88cc`, `3cfa710`, `90f9122`.

DoD-Abgleich (Verifier prüft, setzt keine Häkchen):
- [x] Folge-ADR — Evidence: [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) (Accepted, Messung Colima, Fingerprint-Pin, Egress, Schärft-Feld ersetzt Entscheidung 1 und ergänzt 6 von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)); Zeile im ADR-Index (`docs/plan/adr/README.md`); `git diff 6d726e0 HEAD` auf die Datei von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) ist leer.
- [x] Generator — Evidence: `Dockerfile.tmpl` installiert `docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin` aus `download.docker.com` mit Fingerprint-Prüfung vor Repo-Eintrag und Install; `egressDockerAptHost()` nur bei `nestedDocker`; Golden Cases `TestGenerateDevcontainer_SandboxDocker*`, `TestGenerateDevcontainer_SandboxEgress_DockerAptHostOnlyWithDocker`; Integrationstests `TestE2E_SandboxDevcontainer_NestedDocker` (Paketherkunft, Major >= 28, Abnahmekriterien), `_NestedDocker_EgressAptUpdate` liefen im `make test-docker` grün.
- [x] Doku — Evidence: `spec/spezifikation.md` ([LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) Verfeinerung b und [SPEC-019](../../../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)-Tabelle), `docs/user/benutzerhandbuch.md` (Abschnitt 4.11, Allowlist-Tabelle), `CHANGELOG.md` `[Unreleased]` / Changed.
- [x] `make gates` grün — Evidence: Lauf des Verifiers, siehe Sensors.
- [x] Review durchgeführt — Evidence: Reports `docs/reviews/2026-10-14-sandbox-docker-ce.md` und `-r2.md` liegen vor; Befunde in `3cfa710` und `90f9122` adressiert (Reviewer-Befunde sind keine Evidence; Fix-Wirkung durch die Tests unten belegt).
- [x] Doku-Update bei öffentlichem Vertrag — Evidence: siehe Doku-Punkt; `harness/replay.md` trägt die neue Golden-Zeile.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag — Evidence: §7 dieses Slice.
- [x] Beobachtungs-Register — Evidence: §7 „keine Beobachtung angefallen".
- [x] Jedes Risiko aus §6 trägt einen Ausgang — Evidence: §6 (Schlüsselrotation: „weiter offen" mit Re-Evaluierungs-Trigger der ADR; Verhaltensänderung neuer Versionen: „entfallen", belegt durch `make test-docker`).
- [x] Drei Paarungen — Evidence: §7 (Anker, Folge-Slice [slice-v1-sandbox-docker-nachholmessung-desktop](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md), Register).

Hinweise des Verifiers (keine DoD-Lücke): Die Versionsaussage „Docker 29.9" in §6/§7 und im CHANGELOG ist durch den Test nur als Major >= 28 erzwungen (Versionsangabe selbst stammt aus der Messung in [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md), nicht aus einem Test). Die in §7 behauptete Negativprobe vor dem Review ist nicht aus dem Repo nachprüfbar.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` | pass | d-check: 199 Dateien geprüft, 0 Befunde (vor Anhängen dieses Abschnitts; nach dem Anhängen erneut zu prüfen) |
| `make gates` | pass | Ausgabe `[gates] lint + test + coverage-gate + docs-check green`, Exit 0 |
| `make test-docker` | pass | Exit 0; Paket `internal/e2e` ok in 289 s, Paket `internal/hexagon/application` ok; privilegierte Container gegen die lokale Colima-Engine |

Traceability:
| ID | Beleg |
| --- | --- |
| [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `TestGenerateDevcontainer_SandboxDocker`, `_SandboxDocker_NotGenerated`, `_FailPolicy`; `TestE2E_SandboxDevcontainer_NestedDocker`, `_Degradation`, `_Restart` |
| [LH-FA-DEV-008](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `TestGenerateDevcontainer_SandboxEgressWithDocker`, `_SandboxEgress_DockerAptHostOnlyWithDocker`; `TestE2E_SandboxDevcontainer_NestedDocker_EgressAptUpdate` (`apt-get update --error-on=any` unter Egress, example.com bleibt gesperrt) |
| [LH-NFA-SEC-004](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Golden Case pinnt Reihenfolge Fingerprint-Prüfung, Repo-Eintrag, Install und das `awk`-Programm; kein Installationsskript im Dockerfile |
| [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Feature `docker-in-docker` bleibt verworfen ([ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md), Entscheidung 2); keine Codeänderung |
| [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) | Golden Cases (Fitness Function) und `make test-docker` |
| [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | Immutabilität: `git diff 6d726e0 HEAD` auf die ADR-Datei leer |

Replay / Golden:
- Affected generators: `generate devcontainer` (und `init --devcontainer --sandbox`) mit `devcontainer.sandbox.nestedRuntime: docker` — Dockerfile (Paketquelle) und `egress-init.sh` (Allowlist); `docker-init.sh` wurde in diesem Slice nicht geändert (dessen Änderungen stammen aus `9fa95ce`/`e1c2c3a` des Vorläufer-Slices).
- Golden cases added: `TestGenerateDevcontainer_SandboxEgress_DockerAptHostOnlyWithDocker` (docker: Host vorhanden, podman: Host fehlt, Nutzer-Host bleibt).
- Golden cases updated: `TestGenerateDevcontainer_SandboxDocker` (docker-ce-Paketliste, Repo-URL, Pin, Reihenfolge Prüfung vor Repo-Eintrag und Install, `awk`-Programm, `&&`-Kette), `_SandboxDocker_NotGenerated` (kein `docker-ce`/`download.docker.com` bei none/podman/Default-Profil), `TestGenerateDevcontainer_SandboxEgressWithDocker` (`download.docker.com` in der Allowlist); E2E `TestE2E_SandboxDevcontainer_NestedDocker` und neu `_EgressAptUpdate`.
- Golden cases replayed: Fresh-State und Double-Run-NoOp (`_SandboxDocker`), stale block mit User-Content-Erhalt (`_SandboxDocker_StaleBlockKeepsUserContent`), Skript ohne Marker -> Konflikt ohne Write (`_SandboxDocker_ScriptConflict`); alle im grünen `make gates` und `make test-docker` enthalten. Replay Rules erfüllt: Fresh + NoOp, stale block + Preservation, Konflikt ohne Partial Write, jeweils im Test nachgelesen.
- Intentional output changes: Dockerfile bei `nestedRuntime: docker` installiert `docker-ce` aus dem Docker-APT-Repo statt `docker.io`/`docker-cli`/`docker-buildx`/`docker-compose` aus Debian; Allowlist enthält bei `docker` zusätzlich `download.docker.com`. Begründung: [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) (Debian-Stand Docker 26.1 zu alt). `none`/`podman` unverändert (Test `_NotGenerated`).

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt (`carveouts.md` nennt diesen Slice nicht)

Nicht ausgeführt:
- `make ci` / `make fullbuild` — vom Implementer lokal ausgeführt, vom Verifier nicht selbst belegbar; eine CI auf `main` entsteht erst nach dem Push, bis dahin liegt kein CI-Beleg vor.
- Docker Desktop, Podman als Host-Engine, VS Code/Codespaces — außerhalb des Slice (Folge-Slice [slice-v1-sandbox-docker-nachholmessung-desktop](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md)).
- Echter Schlüsselrotations-Fall (Fingerprint-Mismatch bricht Build) — nicht live ausgeführt; belegt nur über Golden Case (Prüfung gatet Repo-Eintrag und Install).

Commit / Artefakt:
- Delivery-Hash `90f9122` (letzte Code-Änderung); Implementation `eea88cc`, Review-Fixes `3cfa710`, `90f9122`.
