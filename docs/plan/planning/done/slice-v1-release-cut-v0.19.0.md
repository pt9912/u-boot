# Slice V1: Release-Cut `v0.19.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.19.0` von `main` nach `v0.18.0` (Setzung [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md), neu gefasst 2026-11-02: Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag `v0.18.0` (`e6c87379`) unter `## [Unreleased]` steht: der Katalog-Eintragstyp „gepinnte Werkzeugkette" (die zwei Schlüssel `go-toolchain`/`golangci-lint`, Hybrid-Pin per [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md)), die Pin-Politik [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) und die `containerEnv`-Cache-Form (Lastenheft `0.30.0`, [`SPEC-054`](../../../../spec/spezifikation.md), [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md), closed mit CI-Nachweis) — getragen vom CR [`slice-cr-devcontainer-gepinnte-werkzeugkette`](../done/slice-cr-devcontainer-gepinnte-werkzeugkette.md).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Grabstein, nicht belebt), [`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md). Ablauf (einzig normativ): [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) §§3 bis 8 — Form aus der vendored Vorlage [`slice.template.md`](../../../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.19.0` liefert den Block `[0.19.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent nach Freigabe des Projektinhabers (2026-11-02 im Gespräch: Push + Tag nach grünen Vorabprüfungen, Stand-Bericht vor dem Tag).**

**Nummer:** Minor (Zusage-Menge wächst: neuer Katalog-Eintragstyp mit zwei Schlüsseln), keine *Changed*-Einträge (das Verhalten bestehender Kommandos ändert sich nicht; das eigene Devcontainer-Projekt wechselt auf den erzeugten Block). Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice:** Produktcode, `spec/`, ADRs (der Cut ändert nur Release-Dokumente und drei Versionszeilen) · `make test-docker` lokal (der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis) · `tap-check` über den Auftraggeber... (Hauptagent `gh workflow run tap-check` wie bei `v0.18.0`) · Hebung der Pins (kein Anlass, [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)-Stellen bleiben).

**Vor dem Start zu klären (M-Start):** jeder Eintrag unter `## [Unreleased]` trägt einen Beleg und einen Commit im Diff `v0.18.0..HEAD` (gelesen: der Werkzeugketten-Eintrag, Commit `37062fdc`; keine weiteren Träger-Commits unter `[Unreleased]`).

## 2. Definition of Done

- [ ] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt `## [0.19.0] - <Tag-Datum>` mit Lead (Stil `[0.18.0]`: „Nineteenth release."), leerem Anker `[Unreleased]` und Compare-Links; Versionsstrings `0.19.0-dev` an den drei Stellen (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); beide READMEs und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.19.0`; Handbuch-`Stand` gleich dem Tag-Datum.
- [ ] **Vorab-Prüfungen:** `make gates`, `make ci`, `make fullbuild` lokal grün; Trockenlauf des Tap-Pfads nach `releasing.md` §4 (mit Gegenprobe auf den Vorgänger-Tag); `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach nur nach Freigabe**: Tag `v0.19.0` gesetzt und gepusht; Kontrolle nach `releasing.md` §6.
- [ ] **Doku nachgezogen:** Roadmap (Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash), Bestand-Zeilen.
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports unter `docs/reviews/`; Closure-Notiz mit Lerneintrag; Register fortgeschrieben oder „keine Beobachtung"; jedes Risiko mit Ausgang; die drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.19.0] - <Tag-Datum>` mit Lead — Lead nennt die Nutzer-Wirkung samt Grenzen des Lastenheft-Stands `0.30.0`; Fuß: `[Unreleased]: …/compare/v0.19.0...HEAD`, neue Zeile `[0.19.0]: …/compare/v0.18.0...v0.19.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.18.0-dev` → `0.19.0-dev` je Stelle samt Kommentar (Fundstellen beim Start mit `grep -rIn '0\.18\.0-dev'`) |
| `README.md`, `README.de.md` | update | Status-Zeile und `VERSION=`-Beispiel; keine Kennungen, kein „released" |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand`, Versionsbeispiele, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash und Datum |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile für diesen Slice |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung (AGENTS §3.6) |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6): M0 Datum-Abgleich vor jedem Tag · M-T Trockenlauf des Tap-Pfads · M1 Push von `main` (Freigabe erteilt) · M2 CI-Läufe des Push-Tips lesen · M4 `gh workflow run tap-check` · M5 Tag `v0.19.0` gesetzt und gepusht (nur nach Freigabe und Stand-Bericht) · M6 Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): der trägende Umsetzungs-Slice [`slice-v1-katalog-gepinnte-toolchain`](../done/slice-v1-katalog-gepinnte-toolchain.md) liegt in `done/`; `Verantwortlich:` gesetzt; die CI-Läufe des Diffs sind grün oder der Cut wartet.

## 5. Closure-Trigger

DoD vollständig; Tag `v0.19.0` gesetzt und gepusht; `publish` und Job `tap` kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` (AGENTS §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Tabelle** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) — der Reviewer liest den Lead gegen [`SPEC-054`](../../../../spec/spezifikation.md) und das Lastenheft `0.30.0`. — Ausgang: bei Closure.
2. **Das Datum weicht vom Tag-Datum ab** (M0). — Ausgang: bei Closure.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit** (F-9-Klasse — gemessene Marge 19:02 min von 30). — Ausgang: bei Closure.
4. **Job `tap` scheitert** — `tap-check` unmittelbar vor M5, Trockenlauf M-T. — Ausgang: bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Plan).
## 7A — Closure-Notiz (Planner, 2026-11-02, nach Fix-Runde `6daf52ac`, F-5-Nachzug `27bf4304` und Verification)

**Evidence:** Review-Report [`2026-11-02-slice-v1-release-cut-v0.19.0`](../../../reviews/2026-11-02-slice-v1-release-cut-v0.19.0.md) (3 MEDIUM/3 LOW/2 INFO — Lead-Grenzliste ohne Profil-Grenze, Lesung/Verifikation vertauscht, dritter Abweisungsfall unbenannt; kein HIGH) · Verification-Report [`…-verify`](../../../reviews/2026-11-02-slice-v1-release-cut-v0.19.0-verify.md) — **DoD-Konform, closure-fähig**; alle Release-Nachweise im Verifier-Lauf selbst gelesen: Tag-Commit-Läufe `ci` [`41460180976`](https://github.com/pt9912/u-boot/actions/runs/41460180976), `integration-docker` [`41460180946`](https://github.com/pt9912/u-boot/actions/runs/41460180946), `tap-check` [`41463190182`](https://github.com/pt9912/u-boot/actions/runs/41463190182) je success; `publish` [`41463337361`](https://github.com/pt9912/u-boot/actions/runs/41463337361) success mit den vier Jobs; Release `v0.19.0` (12 Assets, kein Draft/Prerelease, `2026-11-02T14:36:59Z`); Tap-Commit `72cabea` (`u-boot 0.19.0`, Formel-URLs auf den Tag); GHCR `0.19.0` == `latest` == `sha256:814b9435…`; `docker run --rm ghcr.io/pt9912/u-boot:0.19.0 --version` = `u-boot version 0.19.0`. Nachzugs-Heads grün (`ci` 41466486010, `integration-docker` 41466485866). Vorab: `gates`/`ci`/`fullbuild` Exit 0, Build-Gegenprobe `VERSION=0.19.0` (`--version` + OCI-Label), Tap-Trockenlauf mit Gegenprobe `v0.18.0` (bricht vor `clone` ab), `doc-immutable` 0 Befunde. **M7:** gelaufen — `u-boot --version` = `0.19.0` und die Bestätigung des Projektinhabers, dass die Version auf dem Mac läuft.

**Risiken aus §6:** Lead-Spanne — F-1 bis F-3 gefunden und berichtigt (`6daf52ac`), Register unten · Datum — M0 ✓ · `integration-docker` am Zeitlimit — ein transienter externer Netzfehler (`Ollama`-Pull 503→`PULL_FAILED`, Lauf 41466002061), Wiederholungs-Lauf grün, Register [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md) **drittes Auftreten, Schwelle erreicht** — Sensor-Kandidat beim Projektinhaber, nichts verkörpert · Job `tap` — `tap-check` grün, Job success ✓.

**Was ging anders als geplant:** (a) Der Lead trug die drei bekannten Klassen erneut (Profil-Grenze fehlte, Lesung/Verifikation vertauscht, Fall unbenannt) — die Grenzlisten-Klasse erreicht die **Drei-Schwelle**; Kandidat eines Pflichtschritts „Lead gegen die Vertragsliste lesen“ in `releasing.md` §3 (Entscheidung des Projektinhabers, wie bei den Vorgänger-Cuts offen). (b) Die Bestand-Zeile `open/README.md` fehlte in der Start-Runde trotz Plan-Zeile (F-5) — in `27bf4304` nachgeholt, der untertritte in der Commit-Message ehrlich berichtigt. (c) Der d-check-Lauf der Fix-Runde deckte die Report-Anker-Umfragereste (Apostroph, `MR-000`-Anker-Slug, wörtliches `(…)`) — der Report ist Bestand, kein Vertragsort.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (`evidence/slice-v1-release-cut-v0.19.0.md`: F-1/F-2/F-3, Lead-Sätze) · [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md) (`evidence/slice-v1-release-cut-v0.19.0.md`: 503/`PULL_FAILED`, drittes Auftreten). *Nicht gezählt:* [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — die Grenzliste ist die engere Lesart (kein Fund als Zusage-Wider); die Klasse der kurzen Grenzliste ist im Vermerk oben benannt. *V-2/V-3 (Verifier):* Roadmap-Label und Vorschlagswerte `v0.20.0` berichtigt (`b32e9c38`-Nachzug). *Plan-Tabelle:* die `docs/reviews/`-Zeilen fehlen in §3 (V-4, INFO) — im Abgleich-Vorbehalt dieser Notiz ausgewiesen.

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre 1:* **Die Lead-Grenzlisten-Klasse erreicht die Drei-Schwelle** — in den drei Cuts `v0.18.0` (F-2, LOW), `v0.19.0` (F-2, 2. Befund, Schwelle laut Skill) und diesem Cut (F-1, MEDIUM-Abstufung). Kandidaten: ein Pflichtschritt „Lead gegen die Vertragsliste lesen“ und/oder ein `d-check`-Prüflauf gegen die Not-promised-Grenzen; die Klassen „Lesung/Verifikation vertauscht“ und „Fall unbenannt“ folgen derselben Linse. *Lehre 2:* **Das e2e-Netzrisiko ist wiederholt** (dritter Ollama-Ausfall; Klasse `BEO-E2E/docker-test-scheitert-am-externen-netz`). *Beide* Setzungen liegen beim Projektinhaber; nichts verkörpert. *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine. *Folge-Slices:* keiner.

**Drei Paarungen** (geprüft): *Anker* — kein `liegt in`-Feld, nichts verkörpert · *Folge-Slice* — keiner versprochen; die Katalog-Restthemen laufen als Vorschlag `v0.20.0` (Roadmap-`Nächste-Wellen`-Zeile und open/README-Bestand, gepinnt durch die F-7-Berichtigung) · *Register* — jede genannte Kennung ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls`).

**Verifikation** (Läufe nach Closure-Inhalt, Register-Tiefen und Verweise-Nachzug; Docker-Modus): `make gates` Exit 0 am 2026-11-02 (`record-gates`-Stempel = HEAD nach dem Nachzug `ecf0adfa`). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.
