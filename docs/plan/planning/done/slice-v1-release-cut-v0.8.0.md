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

- [x] `v0.8.0` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.8.0` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`).
- [x] `make fullbuild` grün vor dem Tag.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update: CHANGELOG, READMEs, Handbuch, Roadmap tragen `v0.8.0`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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

- Der `publish`-Lauf scheitert nach dem GHCR-Push — **Ausgang:** entfallen: `publish`-Lauf `38583042514` lief grün (vier Jobs).
- Das Lastenheft sagt Docker Desktop zu, gemessen ist nur Colima — **Ausgang:** weiter offen: [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** Tag `v0.8.0` auf `1f2b41b` (2026-10-14), nach grünem `ci` und `integration-docker` auf `main` sowie `make fullbuild` und `make test-docker` lokal. `publish`-Lauf `38583042514`: alle vier Jobs grün (`publish`, `tap`, Paket-Smoke `.deb` und `.rpm`); Release mit sechs Binaries, vier Paketen, `SHA256SUMS` und `u-boot.rb`; `ghcr.io/pt9912/u-boot:0.8.0` und `:latest` zeigen auf denselben Digest und melden `0.8.0`.
- **Was ging anders als geplant:** Der erste CI-Lauf nach dem Report-Commit schlug fehl (Links im Review-Report, `docs-check`); mit `1f2b41b` behoben, vor dem Tag. Vor dem Cut zeigte `image-scan` drei HIGH-Advisories der Go-Stdlib; der Go-Pin wurde auf 1.27.2 gehoben (`2a888f7`).
- **Steering-Loop-Eintrag:** Vor dem Push `make docs-check` auch für neue Review-Reports laufen lassen; Reports verlinken Kennungen in Tabellenzellen ohne Backticks im Linktext (der Sensor wertet sie sonst als ungelinkt).
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md); die Tap-Verifikation führt [`slice-v2-homebrew-formula`](slice-v2-homebrew-formula.md).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, Release-Dateien (Versionsstrings in `cmd/uboot/`, Makefile, Dockerfile) — Greenfield laut `harness/conventions.md`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Scope:
- Slice: [slice-v1-release-cut-v0.8.0](slice-v1-release-cut-v0.8.0.md)
- IDs: [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md), [ADR-0007](../../adr/0007-distributionswege-ghcr.md), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../../adr/0017-linux-pakete-nfpm.md)
- Artefakte: Tag `v0.8.0`, GitHub-Release, `ghcr.io/pt9912/u-boot:0.8.0` und `:latest`, `CHANGELOG.md`, `cmd/uboot/main.go`, `Makefile`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md`, `docs/reviews/2026-10-14-release-cut-v0.8.0.md`

DoD-Abgleich:
- [x] Tag und `publish`-Lauf — Evidence: `git rev-parse v0.8.0^{commit}` = `1f2b41b882289645e9279e3b3e24524a081b3e24` (Tag ist ein leichtgewichtiger Tag auf diesen Commit); `publish`-Lauf `38583042514` auf diesem SHA: `conclusion: success`, Jobs `publish`, `tap`, Paket-Smoke deb und rpm je success; `gh release view v0.8.0`: `isPrerelease: false`, `isDraft: false`, zwölf Assets (sechs Binaries `u-boot-{darwin,linux}-{amd64,arm64}`, `u-boot-windows-{amd64,arm64}.exe`; vier Pakete `u-boot_0.8.0_{amd64,arm64}.deb`, `u-boot-0.8.0-1.{x86_64,aarch64}.rpm`; `SHA256SUMS`; `u-boot.rb`); `docker run ghcr.io/pt9912/u-boot:0.8.0 --version` und `:latest` melden je `u-boot version 0.8.0`, beide Tags lösen auf denselben Digest auf.
- [x] `make fullbuild` vor dem Tag — Evidence (indirekt): vom Implementer lokal ausgeführt, vom Verifier nicht reproduzierbar; belegt sind die CI-Läufe auf `1f2b41b`: `ci` `38582043790` success (govulncheck, gates, image-scan) und `integration-docker` `38582043740` success.
- [x] Review — Evidence: Report `docs/reviews/2026-10-14-release-cut-v0.8.0.md` liegt vor (Commits `a2ae07c`, `1f2b41b`); Rollentrennung nicht Gegenstand dieser Prüfung.
- [x] Doku-Update — `CHANGELOG.md` (`[0.8.0] - 2026-10-14`, Compare-Links `v0.7.0...v0.8.0` und `v0.8.0...HEAD`), `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` tragen `v0.8.0`; die Roadmap führt die Welle nach dem Tag unter „Abgeschlossene Wellen“ (Verifier-Lücke vom Implementer nachgezogen, Closure-Commit).
- [x] Versionsstrings — Evidence: `0.8.0-dev` in `cmd/uboot/main.go:43`, `Makefile:27`, `Dockerfile:39`; kein `0.7.0-dev` mehr vorhanden.
- [x] Closure-Notiz mit Steering-Loop-Eintrag — Evidence: §7 vorhanden, nennt Tag-Hash `1f2b41b` und Lerneintrag.
- [x] Beobachtungs-Register — Evidence: §7 „keine Beobachtung angefallen".
- [x] Risiken aus §6 mit Ausgang — Evidence: beide Risiken tragen einen Ausgang (`entfallen`, `weiter offen` mit Folge-Slice).
- [x] Drei Paarungen — Evidence: §7 „Drei Paarungen" verweist auf [slice-v1-sandbox-docker-nachholmessung-desktop](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) und das Register.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` | pass | `d-check: 195 Datei(en) geprüft, 0 Befund(e)` (mit angehängter Evidence) |
| CI `ci` (`38582043790`) | pass | govulncheck, gates (lint, test, coverage-gate), image-scan success auf `1f2b41b` |
| CI `integration-docker` (`38582043740`) | pass | success auf `1f2b41b` |
| `publish` (`38583042514`) | pass | vier Jobs success auf `1f2b41b` |
| `gh release view v0.8.0` | pass | zwölf Assets, kein Pre-Release, kein Draft |
| `docker run ghcr.io/pt9912/u-boot:0.8.0 --version` | pass | `u-boot version 0.8.0` (Image ist `linux/amd64`, Host arm64: Emulations-Warnung) |
| `docker run ghcr.io/pt9912/u-boot:latest --version` | pass | `u-boot version 0.8.0`, gleicher Digest |

Traceability:
| ID | Beleg |
| --- | --- |
| [LH-FA-DEV-007](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | Release-Inhalt laut CHANGELOG `[0.8.0]`; Gates und `integration-docker` auf dem Tag-Commit grün |
| [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | wie oben; Docker-Desktop-Messung offen (siehe unten) |
| [ADR-0007](../../adr/0007-distributionswege-ghcr.md) | GHCR `:0.8.0` und `:latest` mit gleichem Digest, `--version` = 0.8.0 |
| [ADR-0017](../../adr/0017-linux-pakete-nfpm.md) | vier Pakete im Release, Paket-Smoke deb und rpm success |
| [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) | `u-boot.rb` im Release, Job `tap` success; macOS-Smoke im Tap nicht geprüft |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none berührt

Nicht ausgeführt (vom Verifier):
- `make fullbuild`, `make test-docker` — Verifier führt keine Docker-Läufe und keine `make`-Targets außer `make docs-check` aus; Ersatzbeleg sind die CI-Läufe `ci` und `integration-docker` (beide success auf `1f2b41b`). Der lokale Lauf des Implementers ist nicht verifiziert.
- Digest-Vergleich gegen die Release-Assets und `SHA256SUMS` — nicht geprüft.
- Offen, nicht Teil dieses Slice: macOS-Smoke im Tap zur Formel 0.8.0 — [slice-v2-homebrew-formula](slice-v2-homebrew-formula.md); Docker-Desktop-Messung — [slice-v1-sandbox-docker-nachholmessung-desktop](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md).

Commit / Artefakt:
- Tag-Commit `1f2b41b882289645e9279e3b3e24524a081b3e24`
- `publish`-Lauf `38583042514`
- Image `ghcr.io/pt9912/u-boot@sha256:08fc86a5699594ee6101fc97030c0ad996a075c3ff277925fa5f317fd630f26c` (`:0.8.0` und `:latest`)
- Release: https://github.com/pt9912/u-boot/releases/tag/v0.8.0
