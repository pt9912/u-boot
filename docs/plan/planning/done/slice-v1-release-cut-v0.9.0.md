# Slice V1: Release-Cut `v0.9.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.9.0 (Devcontainer-Features Homebrew, SDKMAN, Python).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md),
[ADR-0007](../../adr/0007-distributionswege-ghcr.md),
[ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md),
[ADR-0017](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.9.0` liefert die Katalog-Features `python`, `sdkman` und `homebrew` ([ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md)) als Minor-Release: GHCR-Image, Binaries, Pakete, Homebrew-Formel.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Docker-Desktop-Messung — [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) trägt sie; die Release-Notizen und das Handbuch sagen Docker Desktop nicht zu.
- Der macOS-Smoke im Tap zur Formel 0.9.0 — [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md) führt die Tap-Verifikation.

## 2. Definition of Done

- [x] `v0.9.0` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.9.0` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`).
- [x] `make fullbuild` und `make test-docker` grün vor dem Tag.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update: CHANGELOG, READMEs, Handbuch, Roadmap tragen `v0.9.0`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]` → `[0.9.0] - 2026-10-15`, Compare-Links |
| `cmd/uboot/main.go`, `Makefile`, `Dockerfile` | update | Versionsstrings `0.8.1-dev` → `0.9.0-dev` |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status, Releases-Tabelle, Version |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/in-progress/README.md` | update | Welle nach „Abgeschlossene Wellen" nach dem Tag; Bestand |

## 4. Trigger

**Start** (`next` → `in-progress`): die Slices [`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md) und [`slice-v1-devcontainer-feature-python`](../done/slice-v1-devcontainer-feature-python.md) sind geschlossen, `main` ist grün.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tag.
- `in-progress` → `open` (blockiert): ein roter `image-scan` oder CI-Lauf (Go-Pin heben, `releasing.md` §4) oder ein fehlschlagender `publish`-Lauf, dessen Ursache einen Patch-Release braucht (der Tag wird nicht verschoben).

## 5. Closure-Trigger

DoD vollständig: Tag gesetzt, `publish`- und Tap-Lauf kontrolliert (`releasing.md` §6), Closure-Notiz mit Tag-Hash und Steering-Loop-Eintrag, Verification-Evidence.

## 6. Risiken und offene Punkte

- Der `publish`-Lauf scheitert — **Ausgang:** entfallen: `publish`-Lauf `38689531814` lief grün (vier Jobs).
- Das Zustimmungsmodell für Drittanbieter-Features im Katalog (Kennzeichnung durch Warnung bei `generate`) ist nur in [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) festgelegt, das Lastenheft unverändert — **Ausgang:** weiter offen: Entscheidung des Projektinhabers, ob die Kennzeichnung per Warnung als Zustimmung genügt (Review-Befund F-1 des Slices [`slice-v1-devcontainer-features-homebrew-sdkman`](../done/slice-v1-devcontainer-features-homebrew-sdkman.md)).

## 7. Closure-Notiz

- **Was hat funktioniert:** Tag `v0.9.0` auf `7bfd8d6` (2026-10-15), nach grünem `ci` und `integration-docker` auf `main` sowie `make fullbuild` und `make test-docker` lokal. `publish`-Lauf `38689531814`: alle vier Jobs grün (`publish`, `tap`, Paket-Smoke `.deb` und `.rpm`); Release mit zwölf Assets, kein Prerelease; `ghcr.io/pt9912/u-boot:0.9.0` und `:latest` zeigen nach frischem Pull auf denselben Digest und melden `0.9.0`.
- **Was ging anders als geplant:** Nichts am Ablauf; die Review-Finding des Cuts (ein LOW zum CHANGELOG-Wortlaut) ist vor dem Tag behoben.
- **Steering-Loop-Eintrag:** Vor dem Tag das Datum im Cut gegen das Tag-Datum prüfen (`date +%F`), weil es an sechs Stellen fest steht.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md); die Tap-Verifikation führt [`slice-v2-homebrew-formula`](../done/slice-v2-homebrew-formula.md).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Zustimmungsmodell Drittanbieter-Features) → Entscheidung des Projektinhabers · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, Release-Dateien (Versionsstrings in `cmd/uboot/`, Makefile, Dockerfile) — Greenfield laut `harness/conventions.md`.

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Verifier-Lauf 2026-10-15, nur lesend; Rolle getrennt von Implementer und Reviewer.

Scope:
- Slice: [slice-v1-release-cut-v0.9.0](slice-v1-release-cut-v0.9.0.md)
- IDs: [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md), [ADR-0007](../../adr/0007-distributionswege-ghcr.md), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md), [ADR-0017](../../adr/0017-linux-pakete-nfpm.md)
- Artefakte: Tag `v0.9.0`, GitHub-Release, `ghcr.io/pt9912/u-boot:0.9.0` und `:latest`, `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, `cmd/uboot/main.go`, `Makefile`, `Dockerfile`, `docs/reviews/2026-10-15-release-cut-v0.9.0.md`

DoD-Abgleich:
- [x] Tag, `publish`-Lauf grün — Evidence: `git rev-parse v0.9.0^{commit}` = `7bfd8d6f65df35162558fda09c55c53d239c11d8` (leichtgewichtiger Tag, Typ `commit`); `gh run view 38689531814` auf `7bfd8d6`: `success`, Jobs `publish (ghcr.io image push)`, `tap (homebrew formula)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)` je `success`; `gh release view v0.9.0`: `isDraft=false`, `isPrerelease=false`, zwölf Assets (`SHA256SUMS`, `u-boot.rb`, sechs Binaries linux/darwin/windows je amd64/arm64, zwei `.deb`, zwei `.rpm`); `:0.9.0` und `:latest` melden `u-boot version 0.9.0`, beide Digest `sha256:eb60945f8ecf2f227aae02ee77bec9060cea3a78f6b781506fa9ddeede5c9931` (der Pull meldete „up to date", die Images lagen lokal schon; der Digest-Abgleich gegen die Registry ist belegt, ein Pull ohne Cache nicht).
- [x] `make fullbuild` und `make test-docker` grün vor dem Tag — Evidence: lokal nur vom Implementer gelaufen, vom Verifier nicht selbst belegbar; belegt sind die CI-Läufe auf `7bfd8d6`: `ci` `38688911910` und `integration-docker` `38688911942`, je `success`.
- [x] Review ohne Self-Review — Evidence: `docs/reviews/2026-10-15-release-cut-v0.9.0.md` liegt vor (0 HIGH, 0 MEDIUM, 1 LOW zum CHANGELOG-Wortlaut); behoben, `CHANGELOG.md:25` lautet „the hosts of all three features".
- [x] Doku-Update — Evidence: Versionsstrings `0.9.0-dev` in `cmd/uboot/main.go:43`, `Makefile:27`, `Dockerfile:39`; `git grep 0.8.1-dev` findet nur Beschreibungen (dieser Slice §3 Zeile 49, Slice und Report des Cuts v0.8.1), keine Quelle. `CHANGELOG.md:12` `[Unreleased]` leer, `:14` `[0.9.0] - 2026-10-15`, Compare-Links `:1181-1182`. `README.md:10,105,206`, `README.de.md:11,108,210`, `docs/user/benutzerhandbuch.md:4,80,106,132,1760` tragen `0.9.0`. Roadmap: „Aktuelle Welle" `roadmap.md:16` = Homebrew-Verifikation, „Abgeschlossene Wellen" `:85` (homebrew/sdkman), `:86` (python), `:87` (Cut, Tag `7bfd8d6`).
- [x] Closure-Notiz mit Steering-Loop-Eintrag — Evidence: §7 vorhanden, Tag-Hash und Steering-Loop-Eintrag gelesen; Lauf-ID `38689531814` und Tag-Commit stimmen mit den Befunden oben überein.
- [x] Beobachtungs-Register — Evidence: §7 „keine Beobachtung angefallen".
- [x] Risiken aus §6 mit Ausgang — Evidence: §6 trägt zwei Ausgänge (entfallen mit Lauf-ID; weiter offen mit Entscheidung des Projektinhabers).
- [x] Drei Paarungen — Evidence: §7 verweist auf [slice-v1-sandbox-docker-nachholmessung-desktop](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) und [slice-v2-homebrew-formula](../done/slice-v2-homebrew-formula.md); beide Dateien existieren. Die Anker-Paarung (Zustimmungsmodell) führt keinen Folge-Slice, sondern die Entscheidung des Projektinhabers; sie bleibt offen.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `gh run view 38689531814` (publish) | pass | vier Jobs `success`, Head `7bfd8d6` |
| `gh run list` (ci `38688911910`, integration-docker `38688911942`) | pass | beide `success`, Head `7bfd8d6` |
| `gh release view v0.9.0` | pass | kein Draft, kein Prerelease, zwölf Assets |
| `docker pull` + `docker run --version` (`:0.9.0`, `:latest`) | pass | `u-boot version 0.9.0`, identischer Digest |
| `make docs-check` | fail (erwartet) | 208 Dateien, 1 Befund: `roadmap.md:87` Ziel `../done/slice-v1-release-cut-v0.9.0.md` `target-missing`; löst sich mit dem Move nach `done/` auf (Lifecycle-Übergang, `AGENTS.md` 3.3 Nr. 2); nach dem Move erneut ausführen. Kein weiterer Befund. |

Traceability:
| ID | Beleg |
| --- | --- |
| [LH-FA-DEV-003](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | Release-Inhalt im CHANGELOG `[0.9.0]`; Umsetzung in den geschlossenen Slices homebrew/sdkman und python; `integration-docker` `38688911942` grün |
| [ADR-0020](../../adr/0020-drittanbieter-features-im-katalog.md) | CHANGELOG, README-Releases-Tabelle, Review-Report des Cuts; Zustimmungsmodell offen |
| [ADR-0007](../../adr/0007-distributionswege-ghcr.md) | GHCR `:0.9.0` und `:latest`, gleicher Digest |
| [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) | Release-Asset `u-boot.rb`, Job `tap` grün; macOS-Smoke offen |
| [ADR-0017](../../adr/0017-linux-pakete-nfpm.md) | vier Pakete im Release, Paket-Smoke deb/rpm grün |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: none (nicht geändert, nicht geprüft über den Slice-Inhalt hinaus)

Nicht ausgeführt (vom Verifier):
- `make fullbuild`, `make test-docker`, `make gates` — Docker-/Make-Läufe sind dem Verifier-Auftrag nicht erlaubt; lokal vom Implementer ausgeführt, belegt durch die CI-Läufe `38688911910` und `38688911942` auf `7bfd8d6`.
- Offen, nicht Teil dieses Cuts: macOS-Smoke im Tap zur Formel 0.9.0 — [slice-v2-homebrew-formula](../done/slice-v2-homebrew-formula.md).
- Offen, nicht Teil dieses Cuts: Docker-Desktop-Messung — [slice-v1-sandbox-docker-nachholmessung-desktop](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md); Docker Desktop wird nicht zugesagt.
- Offen: Zustimmungsmodell für Drittanbieter-Features (Warnung bei `generate` als Zustimmung) — Entscheidung des Projektinhabers.
- Pull-/Run-Läufe auf arm64-Host mit linux/amd64-Image (Emulation); nur `--version` geprüft.

Commit / Artefakt:
- Tag-Commit `7bfd8d6f65df35162558fda09c55c53d239c11d8` (Tag `v0.9.0`)
- `publish`-Lauf `38689531814`; `ci` `38688911910`; `integration-docker` `38688911942`
- `ghcr.io/pt9912/u-boot:0.9.0` = `:latest` = `sha256:eb60945f8ecf2f227aae02ee77bec9060cea3a78f6b781506fa9ddeede5c9931`
