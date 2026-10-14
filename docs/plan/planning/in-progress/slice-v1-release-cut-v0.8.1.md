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

- Die Docker-Desktop-Messung — [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) trägt sie; die Release-Notizen und das Handbuch sagen Docker Desktop nicht zu.
- Der macOS-Smoke im Tap zur Formel 0.8.1 — [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md) führt die Tap-Verifikation.

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
- Das Lastenheft sagt Docker Desktop zu, gemessen ist nur Colima — **Ausgang:** weiter offen: [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** Tag `v0.8.1` auf `5aadaf3` (2026-10-14), nach grünem `ci` und `integration-docker` auf `main` sowie `make fullbuild` und `make test-docker` lokal. `publish`-Lauf `38600679867`: alle vier Jobs grün (`publish`, `tap`, Paket-Smoke `.deb` und `.rpm`); Release mit zwölf Assets, kein Prerelease; `ghcr.io/pt9912/u-boot:0.8.1` und `:latest` zeigen auf denselben Digest und melden `0.8.1`.
- **Was ging anders als geplant:** Nichts am Ablauf; die Review-Findings des Cuts (zwei LOW) sind in der Nacharbeit erledigt.
- **Steering-Loop-Eintrag:** Vor dem Cut ein `:latest`-Vergleich nur nach frischem `docker pull` — ein lokal gecachtes `:latest` meldete zunächst die Vorversion.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md); die Tap-Verifikation führt [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, Release-Dateien (Versionsstrings in `cmd/uboot/`, Makefile, Dockerfile) — Greenfield laut `harness/conventions.md`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.
