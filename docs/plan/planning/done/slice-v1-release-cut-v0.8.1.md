# Slice V1: Release-Cut `v0.8.1`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.8.1 (aktueller Docker im Sandbox-Devcontainer).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md),
[ADR-0007](../../adr/0007-distributionswege-ghcr.md),
[ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md),
[ADR-0017](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.8.1` liefert den aktuellen Docker (Docker-APT-Repo, [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md)) für `nestedRuntime: docker` als Patch-Release: GHCR-Image, Binaries, Pakete, Homebrew-Formel.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Docker-Desktop-Messung — [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) trägt sie; die Release-Notizen und das Handbuch sagen Docker Desktop nicht zu.
- Der macOS-Smoke im Tap zur Formel 0.8.1 — [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) führt die Tap-Verifikation.

## 2. Definition of Done

- [x] `v0.8.1` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.8.1` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`).
- [x] `make fullbuild` und `make test-docker` grün vor dem Tag.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update: CHANGELOG, READMEs, Handbuch, Roadmap tragen `v0.8.1`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]` → `[0.8.1] - 2026-10-14`, Compare-Links |
| `cmd/uboot/main.go`, `Makefile`, `Dockerfile` | update | Versionsstrings `0.8.0-dev` → `0.8.1-dev` |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status, Releases-Tabelle, Version |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/in-progress/README.md` | update | Welle nach „Abgeschlossene Wellen" nach dem Tag; Bestand |
| `docs/reviews/2026-10-14-sandbox-docker-ce*.md` | update | Slice-Links nach `done/` |

## 4. Trigger

**Start** (`next` → `in-progress`): der Slice [`slice-v1-sandbox-docker-ce`](../done/slice-v1-sandbox-docker-ce.md) ist geschlossen, `main` ist grün.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tag.
- `in-progress` → `open` (blockiert): ein roter `image-scan` oder CI-Lauf (Go-Pin heben, `releasing.md` §4) oder ein fehlschlagender `publish`-Lauf, dessen Ursache einen Patch-Release braucht (der Tag wird nicht verschoben).

## 5. Closure-Trigger

DoD vollständig: Tag gesetzt, `publish`- und Tap-Lauf kontrolliert (`releasing.md` §6), Closure-Notiz mit Tag-Hash und Steering-Loop-Eintrag, Verification-Evidence.

## 6. Risiken und offene Punkte

- Der `publish`-Lauf scheitert — **Ausgang:** entfallen: `publish`-Lauf `38600679867` lief grün (vier Jobs).
- Das Lastenheft sagt Docker Desktop zu, gemessen ist nur Colima — **Ausgang:** weiter offen: [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** Tag `v0.8.1` auf `5aadaf3` (2026-10-14), nach grünem `ci` und `integration-docker` auf `main` sowie `make fullbuild` und `make test-docker` lokal. `publish`-Lauf `38600679867`: alle vier Jobs grün (`publish`, `tap`, Paket-Smoke `.deb` und `.rpm`); Release mit zwölf Assets, kein Prerelease; `ghcr.io/pt9912/u-boot:0.8.1` und `:latest` zeigen auf denselben Digest und melden `0.8.1`.
- **Was ging anders als geplant:** Nichts am Ablauf; die Review-Findings des Cuts (zwei LOW) sind in der Nacharbeit erledigt.
- **Steering-Loop-Eintrag:** Vor dem Cut ein `:latest`-Vergleich nur nach frischem `docker pull` — ein lokal gecachtes `:latest` meldete zunächst die Vorversion.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md); die Tap-Verifikation führt [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, Release-Dateien (Versionsstrings in `cmd/uboot/`, Makefile, Dockerfile) — Greenfield laut `harness/conventions.md`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Verifier-Lauf 2026-10-14, nur lesend; Rolle getrennt von Implementer und Reviewer.

Scope:
- Slice: [slice-v1-release-cut-v0.8.1](slice-v1-release-cut-v0.8.1.md)
- IDs: [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md), [ADR-0007](../../adr/0007-distributionswege-ghcr.md), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../../adr/0017-linux-pakete-nfpm.md)
- Artefakte: Tag `v0.8.1`, GitHub-Release, `ghcr.io/pt9912/u-boot:0.8.1` und `:latest`, `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, `cmd/uboot/main.go`, `Makefile`, `Dockerfile`

DoD-Abgleich:
- [x] Tag, `publish`-Lauf grün — Evidence: `git rev-parse v0.8.1^{commit}` = `5aadaf3fc60429c70b195e4f2cd3a26f7b9c5e62` (Commit-Objekt); `gh run view 38600679867` auf `5aadaf3`: `success`, Jobs `publish (ghcr.io image push)`, `tap (homebrew formula)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)` je `success`; `gh release view v0.8.1`: `isDraft=false`, `isPrerelease=false`, zwölf Assets (`SHA256SUMS`, `u-boot.rb`, sechs Binaries linux/darwin/windows je amd64/arm64, zwei `.deb`, zwei `.rpm`); frisch gepullt melden `:0.8.1` und `:latest` `u-boot version 0.8.1`, beide Digest `sha256:3fe6fafcf39d7f171079a4e00a5bf971206ecb54b8dee4d78af38e38a3bdfddf`.
- [x] `make fullbuild` und `make test-docker` grün vor dem Tag — Evidence: lokal nur vom Implementer gelaufen, vom Verifier nicht selbst belegbar; belegt sind die CI-Läufe auf `5aadaf3`: `ci` `38599726829` (`gates`, `image-scan`, `security-gates` je `success`) und `integration-docker` `38599726948` (`success`).
- [x] Review ohne Self-Review — Evidence: `docs/reviews/2026-10-14-release-cut-v0.8.1.md` liegt vor (0 HIGH, 0 MEDIUM, 2 LOW; beide laut Slice-§7 nachgezogen, Roadmap-Zeile für den Slice `sandbox-docker-ce` in `roadmap.md:83`, Cut-Zeile `roadmap.md:84`).
- [x] Doku-Update — Evidence: Versionsstrings `0.8.1-dev` in `cmd/uboot/main.go:43`, `Makefile:27`, `Dockerfile:39`; `git grep 0.8.0-dev` findet nur Slice-/Report-Beschreibungen (historisch). `CHANGELOG.md:12` `[Unreleased]` leer, `:14` `[0.8.1] - 2026-10-14`, Compare-Links `:1164-1165`. `README.md:10,206`, `README.de.md:11,210`, `docs/user/benutzerhandbuch.md:4,80,106,132,1749` tragen `0.8.1`. Roadmap: „Aktuelle Welle" `roadmap.md:16` = Homebrew-Verifikation, „Abgeschlossene Wellen" `:83` (Slice) und `:84` (Cut, Tag `5aadaf3`).
- [x] Closure-Notiz mit Steering-Loop-Eintrag — Evidence: §7 vorhanden, Tag-Hash und Steering-Loop-Eintrag gelesen.
- [x] Beobachtungs-Register — Evidence: §7 „keine Beobachtung angefallen".
- [x] Risiken aus §6 mit Ausgang — Evidence: §6 trägt zwei Ausgänge (entfallen mit Lauf-ID; weiter offen mit Folge-Slice).
- [x] Drei Paarungen — Evidence: §7 verweist auf [slice-v1-sandbox-docker-nachholmessung-desktop](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) und [slice-v2-homebrew-formula](../done/slice-v2-homebrew-formula.md); beide Dateien existieren.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `gh run view 38600679867` (publish) | pass | vier Jobs `success`, Head `5aadaf3` |
| `gh run view 38599726829` (ci) | pass | `gates`, `image-scan`, `security-gates` `success` |
| `gh run view 38599726948` (integration-docker) | pass | `success` |
| `docker pull` + `docker run --version` (`:0.8.1`, `:latest`) | pass | `u-boot version 0.8.1`, identischer Digest |
| `make docs-check` | fail (erwartet) | 201 Dateien, 1 Befund: `roadmap.md:84` Ziel `../done/slice-v1-release-cut-v0.8.1.md` `target-missing`; löst sich mit dem Move nach `done/` (Lifecycle-Übergang, `AGENTS.md` 3.3 Nr. 2) auf; vor dem Move erneut ausführen. |

Traceability:
| ID | Beleg |
| --- | --- |
| [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | Release-Inhalt im CHANGELOG `[0.8.1]`; Umsetzung im geschlossenen Slice `sandbox-docker-ce`; `integration-docker` `38599726948` grün |
| [ADR-0019](../../adr/0019-docker-daemon-aus-docker-apt-repo.md) | CHANGELOG, README-Releases-Tabelle, Review-Report des Cuts |
| [ADR-0007](../../adr/0007-distributionswege-ghcr.md) | GHCR `:0.8.1` und `:latest`, gleicher Digest |
| [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) | Release-Asset `u-boot.rb`, Job `tap` grün; macOS-Smoke offen |
| [ADR-0017](../../adr/0017-linux-pakete-nfpm.md) | vier Pakete im Release, Paket-Smoke deb/rpm grün |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none (nicht geändert, nicht geprüft über den Slice-Inhalt hinaus)

Nicht ausgeführt (vom Verifier):
- `make fullbuild`, `make test-docker`, `make gates` — Docker-Läufe sind dem Verifier-Auftrag nicht erlaubt; lokal vom Implementer ausgeführt, belegt durch die CI-Läufe `38599726829` und `38599726948` auf `5aadaf3`.
- macOS-Smoke im Tap zur Formel 0.8.1 — gehört zu [slice-v2-homebrew-formula](../done/slice-v2-homebrew-formula.md); offen.
- Docker-Desktop-Messung — gehört zu [slice-v1-sandbox-docker-nachholmessung-desktop](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md); offen, Docker Desktop wird nicht zugesagt.
- Pull-Läufe liefen auf arm64-Host mit linux/amd64-Image (Emulation); nur `--version` geprüft.

Commit / Artefakt:
- Tag-Commit `5aadaf3fc60429c70b195e4f2cd3a26f7b9c5e62` (Tag `v0.8.1`)
- `publish`-Lauf `38600679867`; `ci` `38599726829`; `integration-docker` `38599726948`
- `ghcr.io/pt9912/u-boot:0.8.1` = `:latest` = `sha256:3fe6fafcf39d7f171079a4e00a5bf971206ecb54b8dee4d78af38e38a3bdfddf`
