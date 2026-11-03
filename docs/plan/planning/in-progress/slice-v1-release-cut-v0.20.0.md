# Slice V1: Release-Cut `v0.20.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.20.0` von `main` nach
`v0.19.0` (Setzung [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md):
Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag
`v0.19.0` (`3564ac31`) unter `## [Unreleased]` steht: der Katalog-Schlüssel
`java-toolchain` — der dritte Schlüssel des Eintragstyps „gepinnte
Werkzeugkette" ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md),
Temurin LTS 21, Same-Origin-Prüfsummen-Quelle, Konflikt-Abweisung external
`java` + `java-toolchain` je Exit-Code `10`), Lastenheft `0.31.0` — getragen
vom CR [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md)
(gelesen: die CR-Closure liegt **nach** dem Tag `v0.19.0`, im Diff
`v0.19.0..HEAD`) und der Umsetzung
[`slice-v1-katalog-feature-jdk-gepinnt`](../done/slice-v1-katalog-feature-jdk-gepinnt.md)
(geschlossen 2026-11-03 mit CI-Nachweis). Der Harness-Vorgang
[`slice-harness-carveouts-baseline-form`](../done/slice-harness-carveouts-baseline-form.md)
berührt keinen öffentlichen Vertrag und trägt keinen `[Unreleased]`-Eintrag.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
[`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md),
[`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md).
Ablauf (einzig normativ): [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) §§3 bis 8 —
Form aus der vendored Vorlage
[`slice.template.md`](../../../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-03.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.20.0` liefert den Block `[0.20.0]` des CHANGELOG als GHCR-Image,
Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent nach
Freigabe des Projektinhabers (Muster `v0.19.0`: Push nach grünen
Vorabprüfungen, Stand-Bericht vor dem Tag; die Freigabe zur Vorab-Phase steht
aus und wird beim M-Start eingeholt).**

**Nummer:** Minor (Zusage-Menge wächst: dritter Schlüssel des
Eintragstyps „gepinnte Werkzeugkette" samt Konflikt-Abweisung), keine
*Changed*-Einträge (das Verhalten bestehender Kommandos ändert sich nicht;
allein der Katalog wächst). Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice:** Produktcode, `spec/`, ADRs (der Cut
ändert nur Release-Dokumente und drei Versionszeilen) · `make test-docker`
lokal (der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis,
Muster [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Docker-Tests) ·
Hebung der Pins (kein Anlass: [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
bleibt bei den zwei Go-Schlüsseln, der Temurin-Default-Pin steht im Template
und hebt nur mit einem u-boot-Release, nicht mit diesem Cut selbst) ·
Tap-Formel-Inhalte über das vom Cut getragene Heben der Version hinaus.

**Vor dem Start zu klären (M-Start):** jeder Eintrag unter
`## [Unreleased]` trägt einen Beleg und einen Commit im Diff
`v0.19.0..HEAD` (gelesen: der java-toolchain-Eintrag, Umsetzung
`9faa3137`; keine weiteren Träger-Commits unter `[Unreleased]`).

## 2. Definition of Done

- [ ] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt
  `## [0.20.0] - <Tag-Datum>` mit Lead (Stil `[0.19.0]`: „Twentieth
  release."), leerem Anker `[Unreleased]` und Compare-Links; Versionsstrings
  `0.20.0-dev` an den drei Stellen (`cmd/uboot/main.go`, `repo.mk`,
  `Dockerfile` — Fundstellen dieses Tages: `cmd/uboot/main.go` Zeile 45,
  `repo.mk` Zeilen 63, 66, 67, `Dockerfile` Zeilen 47, 49, 51; der
  Implementer gleicht mit `grep -rIn '0\.19\.0-dev'` nach); beide READMEs und
  das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen
  `v0.20.0`; Handbuch-`Stand` gleich dem Tag-Datum.
- [ ] **Vorab-Prüfungen:** `make gates`, `make ci`, `make fullbuild` lokal
  grün; Trockenlauf des Tap-Pfads nach `releasing.md` §4 (mit Gegenprobe auf
  den Vorgänger-Tag); `ci` und `integration-docker` auf dem Push-Tip grün;
  `tap-check` grün; **danach nur nach Freigabe**: Tag `v0.20.0` gesetzt und
  gepusht; Kontrolle nach `releasing.md` §6.
- [ ] **Doku nachgezogen:** Roadmap (Zeile nach dem Tag *Abgeschlossene
  Wellen* mit Tag-Hash und Datum), Bestand-Zeilen (`open/README.md`, die
  Zeile dieses Slices, Muster der Zeile `slice-v1-release-cut-v0.19.0`).
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der
  Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports unter
  `docs/reviews/`; Closure-Notiz mit Lerneintrag; Register fortgeschrieben
  oder „keine Beobachtung"; jedes Risiko mit Ausgang; die drei Paarungen
  getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.20.0] - <Tag-Datum>` mit Lead — Lead nennt die Nutzer-Wirkung samt Grenzen des Lastenheft-Stands `0.31.0` (Zusage-Menge Linux `amd64` mit lokaler, nicht verschachtelter Engine, JDK-Behaviour in einer übersteuerten Version nicht versprochen, Brew-Befund ohne zweite Installationsform); Fuß: `[Unreleased]: …/compare/v0.20.0...HEAD`, neue Zeile `[0.20.0]: …/compare/v0.19.0...v0.20.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.19.0-dev` → `0.20.0-dev` je Stelle samt Kommentar (Fundstellen oben; Nachgleich per `grep -rIn '0\.19\.0-dev'`) |
| `README.md`, `README.de.md` | update | Status-Zeile und Versionsbeispiele; keine Kennungen, kein „released" |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand`, Versionsbeispiele, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash und Datum; Zustandsmarker („Vorschlag Release"-Zeile der `Nächste Wellen`-Tabelle) nach dem Tag auf den Ist-Stand (Muster-Marker-Klasse, vier Auftritte) |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile für diesen Slice (mit dem Plan-Commit, nicht erst beim Start — M-F-5-Lektion `v0.19.0`) |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `spec/`, `docs/plan/adr/`, `templates/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung (`AGENTS.md` §3.6); der Temurin-Default-Pin bleibt bei `21.0.13+9` |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6, Muster
`slice-v1-release-cut-v0.19.0`): M0 Datum-Abgleich vor jedem Tag · M-T
Trockenlauf des Tap-Pfads · M1 Push von `main` (Freigabe erteilt) · M2
CI-Läufe des Push-Tips lesen · M4 `gh workflow run tap-check` · M5 Tag
`v0.20.0` gesetzt und gepusht (nur nach Freigabe und Stand-Bericht) · M6
Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): der trägende Umsetzungs-Slice
[`slice-v1-katalog-feature-jdk-gepinnt`](../done/slice-v1-katalog-feature-jdk-gepinnt.md)
liegt in `done/` (2026-11-03 geschlossen); `Verantwortlich:` gesetzt; die CI-Läufe
des Diffs sind grün oder der Cut wartet (Push-Tip `7bce82b5`:
`integration-docker` 41576482662 und `ci` 41576482888 je success; der
Docs-Zug nach dem Tag des Cuts trägt keinen Code-Diff).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.20.0` gesetzt und gepusht; `publish` und Job `tap`
kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice
geht **erst nach dem Tag** nach `done/` (`AGENTS.md` §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Tabelle** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
   — die Lead-Grenzlisten-Klasse steht auf der Drei-Schwelle aus `v0.18.0`
   (`F-2`), `v0.19.0` (`F-1`) und ihrem Vorgänger; Lehre des `v0.19.0`-Cuts:
   der Lead wird gegen die Vertragsliste (Lastenheft `0.31.0`, [`SPEC-054`](../../../../spec/spezifikation.md))
   gelesen, Sätze ohne Beleg-Träger stehen nicht im Lead). — **Ausgang:** bei Closure.
2. **Das Datum weicht vom Tag-Datum ab** (M0). — **Ausgang:** bei Closure.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit** (F-9-Klasse —
   gemessene Marge des JDK-Slices: 41576482662; transienter externer
   Netzfehler möglich, [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md)).
   — **Ausgang:** bei Closure.
4. **Job `tap` scheitert** — `tap-check` unmittelbar vor M5, Trockenlauf M-T. —
   **Ausgang:** bei Closure.
5. **Zustandsmarker nach dem Tag nicht auf den Ist-Stand** (Roadmap, `open/README.md`,
   Entwürfe — vierter und fünfter Auftritt der Marker-Klasse in den Cuts
   `v0.18.0` und `v0.19.0`). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Plan).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `CHANGELOG.md`, beide
READMEs, `docs/user/` (DOK) und `docs/plan/` (Planning, DOK); laut
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)
alle Greenfield (DOK). Keine Code-Sub-Area ist im Diff (Plan oben).

**Vorgelagert — offene Beobachtungen sichten:** [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
(Lead, Risiko 1), [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md)
(Risiko 3). Evidenz-/Diskrepanz-Risiko: niedrig — der Cut führt kein
Produktverhalten, nur Release-Prosa, Versionsteile und Marker.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein
Modus-Begründungsblock nötig).