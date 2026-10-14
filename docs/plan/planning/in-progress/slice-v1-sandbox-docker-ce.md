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
