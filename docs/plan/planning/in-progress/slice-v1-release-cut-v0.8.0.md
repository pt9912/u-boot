# Slice V1: Release-Cut `v0.8.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.8.0 (Docker im Sandbox-Devcontainer).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md),
[ADR-0007](../../adr/0007-distributionswege-ghcr.md),
[ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md),
[ADR-0017](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; Lastenheft 0.5.0 liefert der CR).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.8.0` liefert `nestedRuntime: docker` (Lastenheft 0.5.0, [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)) als Release: GHCR-Image, Binaries, Pakete, Homebrew-Formel.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Docker-Desktop-Messung — [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) trägt sie; die Release-Notizen und das Handbuch sagen Docker Desktop nicht zu.
- Der macOS-Smoke im Tap zur Formel 0.8.0 — [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md) führt die Tap-Verifikation.

## 2. Definition of Done

- [ ] `v0.8.0` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.8.0` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`).
- [ ] `make fullbuild` grün vor dem Tag.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Doku-Update: CHANGELOG, READMEs, Handbuch, Roadmap tragen `v0.8.0`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]` → `[0.8.0] - 2026-10-14`, Compare-Links |
| `cmd/uboot/main.go`, `Makefile`, `Dockerfile` | update | Versionsstrings `0.7.0-dev` → `0.8.0-dev` |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status, Releases-Tabelle, Version |
| `docs/plan/planning/in-progress/roadmap.md` | update | Welle nach „Abgeschlossene Wellen" nach dem Tag |

## 4. Trigger

**Start** (`next` → `in-progress`): die Sandbox-Docker-Slices sind geschlossen, `main` ist grün.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tag.
- `in-progress` → `open` (blockiert): ein roter `image-scan` oder CI-Lauf (Go-Pin heben, `releasing.md` §4) oder ein fehlschlagender `publish`-Lauf, dessen Ursache einen Patch-Release braucht (der Tag wird nicht verschoben).

## 5. Closure-Trigger

DoD vollständig: Tag gesetzt, `publish` und Tap-Lauf kontrolliert (`releasing.md` §6), Closure-Notiz mit Tag-Hash und Steering-Loop-Eintrag.

## 6. Risiken und offene Punkte

- Der `publish`-Lauf scheitert nach dem GHCR-Push — **Ausgang:** weiter offen: Behebung per Workflow-Rerun, nie per Tag-Verschiebung (`releasing.md` §8).
- Das Lastenheft sagt Docker Desktop zu, gemessen ist nur Colima — **Ausgang:** weiter offen: [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, Release-Dateien (Versionsstrings in `cmd/uboot/`, Makefile, Dockerfile) — Greenfield laut `harness/conventions.md`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.
