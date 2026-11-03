# Slice V1: Release-Cut `v0.21.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`; er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.21.0` von `main` nach `v0.20.0` ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md): Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag `v0.20.0` (`90f2a1d2`) unter `## [Unreleased]` steht: der Konfig-Schlüssel `devcontainer.containerEnv` (die Umsetzung [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md), CR [`slice-cr-devcontainer-containerenv`](../done/slice-cr-devcontainer-containerenv.md), Lastenheft 0.32.0) und die zwei *Changed*-Träger des Quellen-Wegs (Katalog-Einträge `cpp`/`postgres-client` entfernt, Abweisung unverarbeitbarer externer Quell-Formen; CR [`slice-cr-devcontainer-katalog-quellen`](../done/slice-cr-devcontainer-katalog-quellen.md), Lastenheft 0.33.0, Umsetzung [`slice-v1-katalog-features-quellen-und-freigabe`](../done/slice-v1-katalog-features-quellen-und-freigabe.md)).

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md). Ablauf (einzig normativ): [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) §§3 bis 8; Form aus [`slice.template.md`](../../../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md). Muster [`slice-v1-release-cut-v0.20.0`](../done/slice-v1-release-cut-v0.20.0.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** `v0.21.0` liefert den Block `[0.21.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent nach Freigabe des Projektinhabers** (Muster `v0.20.0`: Push nach grünen Vorabprüfungen, Stand-Bericht vor dem Tag; freigegeben als Auftrag „danach das Release v0.21.0 schneiden" am 2026-11-03).

**Nummer:** Minor (Zusage-Menge wächst: neuer Konfig-Schlüssel; *Changed*: zwei Katalog-Einträge entfernt, neues Abweisungs-Verhalten je Exit-Code `10`). Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice:** Produktcode, `spec/`-Änderungen (der Cut ändert nur Release-Dokumente und drei Versionszeilen; die Lastenheft-0.33.0-Stellen sind mit den Träger-Commits geschlossen) · `make test-docker` lokal (CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis) · Hebung der Pins · der Restbestand der Feature-Katalog-Zeile (docker-cli-bau, Katalog-Features-Quellen — Backlog).

**Vor dem Start zu klären (M-Start):** jeder `[Unreleased]`-Eintrag trägt einen Beleg und einen Commit im Diff `v0.20.0..HEAD` (gelesen: drei Einträge — `devcontainer.containerEnv` → `06e81ac4`-Vorgänger-Kette des CR `d62aad5e`/Umsetzung; zwei Changed → `151a7976`. Nachzug `git show` beim M-Start).

## 2. Definition of Done

- [ ] Release-Text und Versionsstellen: `CHANGELOG.md` trägt `## [0.21.0] - <Tag-Datum>` mit Lead (Stil `[0.20.0]`: „Twenty-first release."), leerem `[Unreleased]` und Compare-Links; Versionsstrings `0.21.0-dev` an den drei Stellen samt Kommentare (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile` — Nachgleich `grep -rIn '0\.20\.0-dev'`); beide READMEs und das Handbuch tragen `v0.21.0`, Handbuch-`Stand` gleich dem Tag-Datum.
- [ ] Vorab-Prüfungen: `make gates`, `make ci`, `make fullbuild` lokal grün; VERSION-Gegenprobe `make build VERSION=0.21.0`; Tap-Trockenlauf nach `releasing.md` §4 (mit Vorgänger-Gegenprobe); `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach nur nach Freigabe**: Tag `v0.21.0` gesetzt und gepusht; Kontrolle §6.
- [ ] Doku nachgezogen: Roadmap (*Abgeschlossene Wellen* mit Tag-Hash, Zustandsmarker), Bestand-Zeilen.
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports; Closure-Notiz mit Lerneintrag; Register oder „keine Beobachtung"; jedes Risiko mit Ausgang; drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.21.0] - <Tag-Datum>` mit Lead (Nutzer-Wirkung samt Grenzen: Lastenheft 0.32.0/0.33.0, Messgrenze der Quellen-Formen, CLI-Verweigerung der Tokens, Tarball-Weg offen); Compare-Links `[Unreleased]: …/compare/v0.21.0...HEAD`, `[0.21.0]: …/compare/v0.20.0...v0.21.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.20.0-dev` → `0.21.0-dev` je Stelle samt Kommentar (Fundstellen beim M-Start mit `grep`) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status-Zeile, `VERSION=`-Beispiel, `18 built-in features`-Zählung unverändert, Kopf/`Stand` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash; Zustandsmarker am Ist |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile dieses Slices (mit dem Plan-Commit — M-F-5-Lehre) |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung (`AGENTS.md` §3.6) |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6): M0 Datum-Abgleich · M-T Tap-Trockenlauf · M1 Push (Freigabe) · M2 CI-Läufe lesen · M4 `gh workflow run tap-check` · M5 Tag (nur nach Freigabe und Stand-Bericht) · M6 Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`next` → `in-progress`): die trägenden Vorgänge der drei `[Unreleased]`-Einträge liegen in `done/`; `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`.

**Rückführungen:** `in-progress` → `next` (zu groß — etwa die Vorab-Prüfungen verlangen mehr als den Cut); `in-progress` → `open` (blockiert — Vorab rot oder `tap-check` ungrün; der Cut wartet).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.21.0` gesetzt und gepusht; `publish` und Job `tap` kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` (`AGENTS.md` §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Menge** (Lead-Grenzlisten-Klasse des Vorgänger-Quellen-Vertrags) — der Lead liest gegen die Vertragsliste (Lastenheft 0.32.0/0.33.0, [`SPEC-049`](../../../../spec/spezifikation.md)/[`SPEC-023`](../../../../spec/spezifikation.md)); Grenzsätze je Beleg — Ausgang: bei Closure.
2. **Das Datum weicht vom Tag-Datum ab** (M0). — Ausgang: bei Closure.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit** (F-9-Klasse) — Ausgang: bei Closure.
4. **Job `tap` scheitert** — Trockenlauf M-T, `tap-check` vor M5. — Ausgang: bei Closure.
5. **Zustandsmarker nach dem Tag nicht auf den Ist-Stand** (Marker-Klasse der Cuts) — Ausgang: bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `CHANGELOG.md`/READMEs/`docs/user/` (DOK) und `docs/plan/` (DOK); alle Greenfield. **Vorgelagert — offene Beobachtungen sichten:** [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 1), [`BEO-E2E/test-docker-nachlauf-ohne-stempel`](../observations/BEO-E2E/test-docker-nachlauf-ohne-stempel/observation.md) (kein neues Träger-Diff — der Cut führt keinen Code). Evidenz-/Diskrepanz-Risiko: niedrig. **Modus-Begründungsblock:** alle berührten Sub-Areas Greenfield.