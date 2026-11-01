# Slice V1: Release-Cut `v0.18.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** keine Welle-Datei — Minor-Release `v0.18.0` von `main` nach `v0.17.0` (Release-Zuordnung und Einstufung entschieden: Minor, *Changed* — Projektinhaber 2026-11-01 im Slice [`slice-v1-devcontainer-container-name`](../in-progress/slice-v1-devcontainer-container-name.md), §6 Frage 8; wie die Cuts ohne Welle-Datei, geführt in der [Roadmap](../in-progress/roadmap.md), bis der Tag gesetzt ist). Der Release trägt, was seit dem Tag `v0.17.0` (`f89424e1`) unter `## [Unreleased]` steht: die Auflösung von Git-Zugangsdaten und Identität des Hosts (opt-in, Lastenheft `0.28.0`, [`slice-v1-devcontainer-host-git-uebernehmen`](../done/slice-v1-devcontainer-host-git-uebernehmen.md), closed mit CI-Nachweis 41328395477), die Schlüssel `devcontainer.resources.memory`/`.cpus` (Lastenheft `0.27.0`, [`slice-v1-devcontainer-ressourcenlimits`](../done/slice-v1-devcontainer-ressourcenlimits.md)) und der Containername `u-boot-<project>[-<instance>]` (Lastenheft `0.29.0`, [`slice-v1-devcontainer-container-name`](../in-progress/slice-v1-devcontainer-container-name.md) — **noch offen: Review, Verification, Closure**). Planungs-Commits ohne Release-Wirkung: der CR-Entwurf [`slice-v1-katalog-gepinnte-toolchain-entwurf`](slice-v1-katalog-gepinnte-toolchain-entwurf.md) und die Messungs-Commits des Container-Name-Slices (Messprotokoll, kein Produktcode).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Containername und Instanz-Aufrufe), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Zugangsdaten/Identität), [`LH-FA-CONF-002`](../../../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Ressourcen-Keys, HostGitToken-Schalter), [`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Pins unverändert), [`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md), [ADR-0040](../../adr/0040-host-git-quellen-fuer-token-und-identitaet.md). Ablauf: [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md). Muster: [`slice-v1-release-cut-v0.17.0`](../done/slice-v1-release-cut-v0.17.0.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; dieses Lastenheft `0.29.0` lieferte der CR des Container-Name-Slices; dieser Slice ändert `spec/` nicht).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.18.0` liefert den Block `[0.18.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. Der Slice bereitet vor; **Push und Tag führt nur der Hauptagent nach ausdrücklicher Freigabe des Projektinhabers im Gespräch**.

**Nummer:** Minor (Zusage-Menge wächst: zwei neue Konfigurationsbereiche und ein opt-in Schalter), mit *Changed* für den Container-Verhalten-Wechsel (`up` startet künftig benannte Container); entschieden 2026-11-01. Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tag und alles mit Außenwirkung durch den Implementer** ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*); `tap-check` startet der Auftraggeber mit dem PAT.
- **`make test-docker` lokal** — Entscheidung des Projektinhabers (2026-11-01): der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis; der Monitor der beiden Läufe `41335734851`/`41335734895` trägt den Stand des Diffs.
- **Produktcode, `spec/`, ADRs** — der Cut ändert nur Release-Dokumente und drei Versionszeilen.
- **Offene Entwürfe** ([`slice-v1-katalog-gepinnte-toolchain-entwurf`](slice-v1-katalog-gepinnte-toolchain-entwurf.md), [`slice-v1-katalog-feature-jdk-gepinnt`](slice-v1-katalog-feature-jdk-gepinnt.md), [`slice-v1-docker-cli-bau-auf-standard-basis-image`](slice-v1-docker-cli-bau-auf-standard-basis-image.md), die beiden Test-docker-Slices) — sind nicht Inhalt des Release und geben ihm keine Bedingung.
- **Hebung von `GO_VERSION`/`GOLANGCI_LINT_VERSION`** — nur Anlass bei rotem `image-scan` ([`releasing.md`](../../../maintainer/releasing.md) §4, §8).

**Vor dem Start zu klären** (Implementer, M-Schritt): jeder Eintrag unter `## [Unreleased]` trägt einen Beleg (`LH-*`, Slice) und einen Commit im Diff `v0.17.0..HEAD`; ein Eintrag ohne beides ist ein Befund, der vor dem Cut geklärt wird (gelesen: drei Einträge mit Trägern; die Hygiene-Einträge der Hilfe-Texte gehören zu einem älteren Commit und brauchen diesen Abgleich, nicht Annahme).

## 2. Definition of Done

- [ ] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt `## [0.18.0] - <Tag-Datum>` mit Lead (Stil `[0.17.0]`: „Eighteenth release.“), leerem Anker `[Unreleased]` und Compare-Links; Versionsstrings `0.18.0-dev` an den drei Stellen (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); beide READMEs und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.18.0`; Handbuch-`Stand` gleich dem Tag-Datum.
- [ ] **Vorab-Prüfungen:** `make gates`, `make ci`, `make fullbuild` lokal grün; Trockenlauf des Tap-Pfads (Muster `v0.17.0` §3 M-T); `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach und nur nach Freigabe**: Tag `v0.18.0` gesetzt und gepusht; Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6.
- [ ] **Doku nachgezogen:** Roadmap (Ruhe-Marker → Zeiger → nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash), `releasing.md`-Zeiger „zuletzt", Zeilen der Bestands-Readmes.
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports unter `docs/reviews/`; Nachlauf-Review über jeden Diff von Vertragstext.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register fortgeschrieben oder „keine Beobachtung" notiert; jedes Risiko mit Ausgang; die drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung (Fundstellen) |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt (gelesen: drei Zusage-Einträge, [`AGENTS.md`](../../../../AGENTS.md) §3.10) unter `## [0.18.0] - <Tag-Datum>` mit Lead — Lead nennt die drei Nutzer-Wirkungen samt Grenzen der Lastenheft-Stände `0.27.0`/`0.28.0`/`0.29.0`; Fuß: `[Unreleased]: …/compare/v0.18.0...HEAD`, neue Zeile `[0.18.0]: …/compare/v0.17.0...v0.18.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.17.0-dev` → `0.18.0-dev` je Stelle samt Kommentar (Fundstellen beim Start mit `grep -rIn '0\.17\.0-dev'` erheben) |
| `README.md`, `README.de.md` | update | Status-Zeile und `VERSION=`-Beispiel (README.md Zeilen 10 und 98, README.de.md Zeilen 11 und 101); keine Kennungen, kein „released" |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand` (heut `2026-11-01`), Versionsbeispiele, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeiger auf die Umsetzung (Container-Name) weg nach Closure; Zeiger auf diesen Slice unter *Offene Wellen*; nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash und Datum |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile für diesen Slice |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `docs/maintainer/releasing.md` | update (Fix-Runde) | Zeiger „zuletzt" (Zeilen 56 bis 57 nennen `v0.17.0`) → dieser Slice |
| `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |

### Reihenfolge der Maintainer-Aktionen (Muster `v0.17.0`)

M0 Datum-Abgleich vor jedem Tag (`date +%F`, CHANGELOG-Überschrift und Handbuch-`Stand`) · M-T Trockenlauf des Tap-Pfads · M1 Push von `main` (nach Freigabe) · M2 CI-Läufe des Push-Tips lesen · M3 Vorab-Tag entfällt außer auf Wunsch · M4 `gh workflow run tap-check` (Auftraggeber, unmittelbar vor M5) · M5 Tag `v0.18.0` gesetzt und gepusht (nur nach Freigabe des Projektinhabers) · M6 Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 · M7 Mac-Probe des Tap.

**Sensoren und Nachweise** wie beim Muster `v0.17.0` §3 (lokale Läufe mit Exit-Code, `make build VERSION=0.18.0` mit `--version`- und OCI-Label-Gegenprobe, `make doc-immutable RANGE="v0.17.0..HEAD"`, `doc-trace`, `doc-commits`, `doc-planning`; `verify-depguard` nur bei berührten Importregeln).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): die getragenen Slices liegen in `done/` — **der Container-Name-Slice folgt mit Review, Verification und Closure** (Review läuft); `Verantwortlich:` gesetzt; `in-progress/` trägt sonst nichts (WIP-Limit 1); die CI-Läufe des Diffs `33e9cd5d` sind grün oder der Cut wartet.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): ein `--- FAIL` rot in `integration-docker` auf dem Push-Tip (dann ein eigener Fix-Slice vor dem Cut); ein roter `image-scan` (Go-Pin heben, eigener Vorgang); ein roter `tap-check`; ein Träger-Eintrag im CHANGELOG ohne Beleg (§1 Vor dem Start).
- `in-progress` → `next`: nicht vorgesehen (der Cut ist ein Slice mit drei Textstellen).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.18.0` gesetzt und gepusht (M5, nur mit Freigabe); `publish` und Job `tap` kontrolliert (M6); Trockenlauf-Serge festgehalten; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2). Ohne Freigabe bleibt er in `in-progress/`.

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als das Lastenheft** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). Vorsorge: der Lead trägt je Eintrag die Bedingungen des Lastenheft-Stands und die Liste des Nicht-Zugesagten im Wortlaut des CHANGELOG-Eintrags; der Reviewer liest ihn gegen die drei Stände (`0.27.0`, `0.28.0`, `0.29.0`). — **Ausgang:** bei Closure.
2. **Ein Unreleased-Eintrag ohne Träger** (etwa der Hygiene-Eintrag der Hilfetexte). Vorsorge: §1-Abgleich mit `git diff v0.17.0..HEAD`. — **Ausgang:** bei Closure.
3. **Das Datum weicht vom Tag-Datum ab** (Mitternacht). Vorsorge: M0 vor M1 und M5. — **Ausgang:** bei Closure.
4. **Ein neuer Advisory macht `image-scan` rot.** Vorsorge: `make ci` lokal, Läufe des Push-Tips. — **Ausgang:** bei Closure.
5. **Job `tap` scheitert im stabilen Lauf.** Vorsorge: `tap-check` unmittelbar vor M5; Trockenlauf M-T. — **Ausgang:** bei Closure.
6. **`integration-docker` am Push-Tip rot oder am Zeitlimit.** Vorsorge: Monitor läuft; ein Fix-Runde gehört dem trägenden Slice. — **Ausgang:** bei Closure.
7. **Der `up`-Verhalten-Wechsel mit den benannten Containern überrascht Skripte, die den Container-Namen parsen** (`docker ps --format` mit Zufallsnamen-Konvention). Vorsorge: der Lead nennt die neue Form als *Changed* mit dem Hinweis auf die Such-Identität (Labels unverändert). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Plan).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `CHANGELOG.md`, READMEs, Handbuch und Planung (alle DOK laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)); die drei Versionszeilen liegen in `cmd/` und dem Build-Werkzeug — `cmd/` führt keine Sub-Area (Modus-Deklaration, Zeile Hexagon-Schichtungs-Linie), die Builds der Tool-Dateien sind keine Sub-Area-Änderung, weil keine Konvention berührt wird.

**Vorgelagert — offene Beobachtungen sichten** (Register 2026-11-01 durchgesehen, nur zitierte Einträge gelesen): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 1), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (CI-Nachweis statt test-docker lokal). Evidenz-/Diskrepanz-Risiko: gering (der Cut trägt geschlossene Zusagen; der Reviewer liest den Lead gegen die Lastenheft-Stände). Alle berührten Sub-Areas Greenfield — kein Modus-Begründungsblock nötig.