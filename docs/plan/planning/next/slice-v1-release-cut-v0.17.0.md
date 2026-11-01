# Slice V1: Release-Cut `v0.17.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** keine Welle-Datei — Minor-Release `v0.17.0` von `main` nach `v0.16.1` (Release-Version nach [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md); wie `v0.16.1` ohne Welle-Datei, geführt in der [Roadmap](../in-progress/roadmap.md) unter *Offene Wellen*, bis der Tag gesetzt ist). Der Release trägt, was seit dem Tag `v0.16.1` auf `main` liegt und in `CHANGELOG.md` unter `## [0.17.0]` steht: `--instance` für `u-boot devcontainer up`, `down` und `shell` ([`slice-v1-devcontainer-instanz-label`](../done/slice-v1-devcontainer-instanz-label.md)), die Referenz `CLAUDE_CODE_OAUTH_TOKEN` für `claude-code` ([`slice-v1-devcontainer-claude-oauth-token`](../done/slice-v1-devcontainer-claude-oauth-token.md)), das Lastenheft `0.26.0` ([`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md), [`slice-cr-devcontainer-instanz-zusage-getestet`](../done/slice-cr-devcontainer-instanz-zusage-getestet.md)) und die Doku-Änderung zum Homebrew-Tap unter Linux. Weitere Slices im Release, ohne eigenen CHANGELOG-Eintrag: [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md), [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md), [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (nach dem Revert nur ein Kommentar in `devcontainer_shell.go`).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope: die Instanz-Zusage, die der Release ausliefert), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Referenz `CLAUDE_CODE_OAUTH_TOKEN`), [`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md). Ablauf: [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md). Muster: [`slice-v1-release-cut-v0.16.1`](../done/slice-v1-release-cut-v0.16.1.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; das Lastenheft `0.26.0` lieferten die CR-Slices; dieser Slice ändert `spec/` nicht).

**Verantwortlich:** pt9912 (Rolle Implementer).

**Autor:** pt9912 (Rolle Implementer, im Auftrag des Hauptagenten). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.17.0` liefert den Block `[0.17.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. Dieser Slice bereitet vor; **Push und Tag führt nur der Hauptagent nach ausdrücklicher Freigabe des Projektinhabers im Gespräch**.

**Nummer:** Minor, weil `--instance` eine neue Option ist (SemVer §2); die Roadmap führte den Release als `v0.17.0`. Freigabe des Projektinhabers: „Release nach releasing.md erstellen“ (2026-11-01, Wortlaut des Auftrags an den Implementer). Kein Vorab-Tag (Vorschlag des Implementers wie bei `v0.13.0` bis `v0.16.1`; der Projektinhaber entscheidet, ob er einen `-rc` will).

**Unterschied zu `v0.16.1`:** Der Code des Binary ist **nicht** unverändert. `git diff v0.16.1..HEAD --stat -- . ':!docs' ':!spec' ':!*.md'` trifft vor dem Cut 70 Dateien, darunter Produktcode (`internal/adapter/driving/cli/devcontainer.go`, `internal/adapter/driven/docker/`, `internal/hexagon/`), Tests, Messskripte unter `tools/messung/`, das Build-Werkzeug (`repo.mk`, `scripts/bau.sh`, `.devcontainer/`, `Dockerfile`). Der Lead sagt deshalb nichts über einen unveränderten Binary-Stand. Das Werkzeug (`BUILD_MODE=local`, Devcontainer des Repos) ist für Nutzer des Binary ohne Wirkung und hat keinen CHANGELOG-Eintrag.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tag und jede Aktion mit Außenwirkung durch den Implementer.** Extern wirksam, nicht sauber zurückzunehmen ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*); nie `git push --tags`, ein gesetzter Tag wird nicht verschoben (§8 dort).
- **`gh workflow run tap-check`.** Der Maintainer startet es mit dem PAT ([`releasing.md`](../../../maintainer/releasing.md) §2, §4); der Implementer hat kein Token.
- **`make test-docker` lokal.** Entscheidung des Projektinhabers (2026-11-01): der CI-Lauf `integration-docker` auf `main` ist der Nachweis (Vorbild: M2 im Cut `v0.16.1`).
- **Produktcode, Spezifikation, Lastenheft, ADRs.** Die Aufträge sind Text in Release-Dokumenten und drei Versionszeilen.
- **Der Slice [`slice-v1-devcontainer-ressourcenlimits`](../open/slice-v1-devcontainer-ressourcenlimits.md)** (Entwurf mit „Vorschlag `v0.17.0`“): er liegt in `open/` und ist nicht Inhalt des Release.

## 2. Definition of Done

- [ ] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt `## [0.17.0] - 2026-11-01` mit Lead, leerem Anker `[Unreleased]` und Vergleichslinks; Versionsstrings `0.17.0-dev` an den drei Stellen; beide READMEs und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.17.0`; das Datum steht in CHANGELOG und Handbuch-`Stand` gleich dem Tag-Datum. — Evidence: §3.
- [ ] **Tag und Kontrolle:** `make gates`, `make ci`, `make fullbuild` und der Trockenlauf des Tap-Pfads sind vor dem Tag grün (§3); `ci` und `integration-docker` auf dem Push-Tip sind grün (M2); `tap-check` ist grün (M4); danach und nur nach ausdrücklicher Freigabe: Tag `v0.17.0` gesetzt und gepusht (M5); Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 (M6). — Evidence: folgt nach dem Tag in §7.
- [ ] **Doku nachgezogen:** Roadmap (Zeiger unter *Offene Wellen* weg, Zeile in *Abgeschlossene Wellen* mit Tag-Hash und Datum, Ruhe-Marker), `releasing.md` (Zeiger „zuletzt“) und die Zeile in `open/README.md`. — folgt nach dem Tag.
- [ ] `make gates` grün (nach dem letzten Commit vor dem Handoff und nach der Closure).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag. — folgt nach dem Tag in §7.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (kein Welle-Betrieb).

## 3. Plan (vor Code)

Jede Aussage ist an der genannten Stelle gelesen (Stand 2026-11-01, `HEAD` = `f58dc581` vor dem Cut, [`AGENTS.md`](../../../../AGENTS.md) §3.10).

**Versionsstellen** (gelesen mit `grep -rIn -e '0\.16\.1' --exclude-dir=.git --exclude-dir=.harness --exclude-dir=reviews --exclude-dir=messungen .`): `cmd/uboot/main.go:43`, `repo.mk:57`, `:60`, `:61`, `Dockerfile:41`, `:43`, `:45`, `README.md:10`, `:98`, `README.de.md:11`, `:101`, `docs/user/benutzerhandbuch.md:4`, `:84`, `:110`, `:136`, `:3038`. **Keine Versionsstelle** (gelesen): `docs/plan/adr/0037-bauen-im-devcontainer-lokaler-modus.md:27` („In `u-boot` 0.16.1 setzt …“, Aussage über jenen Stand in einer Accepted-ADR, unveränderlich), `docs/plan/planning/open/slice-v1-devcontainer-ressourcenlimits.md:39` und `done/slice-v1-devcontainer-instanz-messung.md` (Messstand „u-boot 0.16.1“), der Abschnitt `[0.16.1]` und die Link-Zeile im CHANGELOG, `done/slice-v1-release-cut-v0.16.1.md`.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update (Cut) | `## [Unreleased]` bleibt leerer Anker; Inhalt (`### Added`, `### Changed`) wandert unter `## [0.17.0] - 2026-11-01` mit Lead (Stil von `[0.16.0]`, „Seventeenth release.“); Fuß: `[Unreleased]: …/compare/v0.17.0...HEAD`, neue Zeile `[0.17.0]: …/compare/v0.16.1...v0.17.0`. Die Einträge bleiben im Wortlaut. |
| `cmd/uboot/main.go` | update | `var version = "0.16.1-dev"` → `0.17.0-dev`. |
| `repo.mk` | update | Kommentare (zwei Zeilen) und `VERSION ?= 0.16.1-dev` → `0.17.0-dev`. |
| `Dockerfile` | update | Kommentare (zwei Zeilen) und `ARG UBOOT_VERSION=0.16.1-dev` → `0.17.0-dev`. |
| `README.md`, `README.de.md` | update | Status-Zeile (`current version` / `aktuelle Version`) und `VERSION=`-Beispiel; keine Kennungen, kein „released“, keine Tabelle (gelesen: die READMEs tragen keine Release-Tabelle mehr). |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand` (2026-11-01), Versionsbeispiel `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich. `Handbuch-Version: 2.1` bleibt: sie zählt die Fassung des Dokuments, `releasing.md` §3 und §5 verlangen sie nicht. |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile `v0.17.0` aus *Nächste Wellen* entfernt; unter *Offene Wellen* ein Zeiger auf diesen Slice an Stelle des Ruhe-Markers (`Nichts in Arbeit.` gilt nur, wenn `in-progress/` keinen Slice trägt). Nach dem Tag: Zeiger weg, Marker zurück, Zeile in *Abgeschlossene Wellen*. |
| `docs/plan/planning/open/README.md` | update | Zeile für diesen Slice (Muster: Zeile `v0.16.1`); die Moves ziehen den Link per `make slice-mv` nach. |
| `docs/plan/planning/done/slice-v1-release-cut-v0.17.0.md` (dieser Slice; Anlage in `open/`) | neu | Moves nach `next/` und `in-progress/` per `make slice-mv`; nach dem Tag nach `done/`. |
| `docs/maintainer/releasing.md` | kein Diff im Cut | `:56`–`:57` nennt als „zuletzt“ den Cut `v0.16.1` (gelesen): richtig bis zum Abschluss; die Closure setzt den Zeiger. |
| `.github/workflows/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `spec/`, `docs/plan/adr/` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6). |

### Reihenfolge der Maintainer-Aktionen

| # | Aktion | Wer |
|---|---|---|
| M0 | `date +%F` vor M1 und vor M5; CHANGELOG-Überschrift und Handbuch-`Stand` gegen das Tag-Datum halten; weicht es ab, ein weiterer Cut-Commit und M2 auf dem neuen Tip | Hauptagent (Implementer: 2026-11-01) |
| M-T | Trockenlauf des Tap-Pfads ([`releasing.md`](../../../maintainer/releasing.md) §4), unten | Implementer |
| M1 | `git push origin main` | Hauptagent, nach Freigabe |
| M2 | `gh run list`: `ci` und `integration-docker` auf dem Push-Tip; Logs über REST-API und Annotationen lesen | Hauptagent |
| M3 | Vorab-Tag `v0.17.0-rc.1`: entfällt, außer der Projektinhaber verlangt ihn | — |
| M4 | `gh workflow run tap-check` (PAT) | Auftraggeber, unmittelbar vor M5 |
| M5 | `git tag v0.17.0` auf dem in M2 geprüften Commit, `git push origin v0.17.0` (nicht `--tags`) | Hauptagent, nach Freigabe |
| M6 | [`releasing.md`](../../../maintainer/releasing.md) §6: `gh run watch`, `gh release view`, `docker run … --version`, Log des Jobs `tap` | Hauptagent |
| M7 | Auf einem Mac `brew upgrade u-boot`, `u-boot --version` = `0.17.0` | Projektinhaber |

### Sensoren (Implementer, Docker-Modus, Host `linux/amd64`)

AUSSTEHEND-SENSOREN

### Nacharbeit nach dem Tag ([`releasing.md`](../../../maintainer/releasing.md) §7)

- [ ] Roadmap: Zeile `v0.17.0` nach *Abgeschlossene Wellen* mit Tag-Hash und Datum; Zeiger unter *Offene Wellen* weg, Ruhe-Marker zurück.
- [ ] Datum gegen `git log -1 --format=%cs v0.17.0` prüfen (CHANGELOG, Handbuch-`Stand`).
- [ ] Job `tap` im Log des Laufs `publish` prüfen („updated to v0.17.0.“ oder „already at“, keine Warnung zum fehlenden Secret).
- [ ] Zeiger „zuletzt“ in `releasing.md` auf diesen Slice.
- [ ] Closure-Notiz (§7) mit Tag-Hash, Lauf-IDs und Sensoren; Slice nach `done/` (erst Inhalt, dann reiner `git mv`, [`AGENTS.md`](../../../../AGENTS.md) §3.3).

## 4. Trigger

**Start** (`open` → `in-progress`): die Nummer `v0.17.0` gilt (Auftrag); die Inhalts-Slices liegen in `done/` (gelesen); `Verantwortlich:` ist gesetzt; `in-progress/` enthält nur die Master-Dokumente (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): nicht vorgesehen.
- `in-progress` → `open` (blockiert): ein roter `image-scan` durch einen neuen Advisory (Go-Pin heben, eigener Vorgang davor); ein Befund des Trockenlaufs, der eine Änderung an Workflow oder Skript verlangt; ein `--- FAIL` oder Zeitlimit-Abbruch in `integration-docker` auf dem Push-Tip; ein roter `tap-check`.

## 5. Closure-Trigger

DoD vollständig: Tag `v0.17.0` gesetzt und gepusht (M5), `publish` und Job `tap` kontrolliert (M6), Trockenlauf mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag, Verification Evidence (Verifier). Der Slice geht **erst nach dem Tag** nach `done/`. Kein Tag ohne ausdrückliche Freigabe: bleibt sie aus, bleibt der Slice in `in-progress/` oder geht mit benanntem Grund nach `open/` zurück.

## 6. Risiken und offene Punkte

- **Der Lead sagt mehr, als der Vertrag zusagt** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). Gegenmittel: der Lead nennt die Zusage mit den Bedingungen von `spec/lastenheft.md` (Zeile 504: zwei Instanzen, Workspace-Volumes, Sandbox-Profil, direkt nach dem Start, Linux `amd64`, lokale Docker-Engine) und die Liste des Nicht-Zugesagten; der Reviewer liest ihn gegen das Lastenheft. — **Ausgang:** offen bis zur Closure.
- **Der Lead lässt eine Nutzer-Änderung aus.** Gelesen: der Diff gegen `v0.16.1` im Produktcode besteht aus `--instance`, der Prüfung des Namens, der Referenz `CLAUDE_CODE_OAUTH_TOKEN` samt `doctor`-Regel und der Konstante `domain.ConfigFileName`; der Commit zu `shell` ist nach dem Revert ein Kommentar. Ein weiterer Nutzer-Effekt ist nicht gelesen worden. — **Ausgang:** offen bis zur Closure.
- **Das Datum weicht vom Tag-Datum ab** (Mitternacht). Gegenmittel: M0. — **Ausgang:** offen bis zur Closure.
- **Ein neuer Advisory macht `image-scan` rot.** Gegenmittel: `make ci` lokal und `ci` auf dem Push-Tip. — **Ausgang:** offen bis zur Closure.
- **Der Job `tap` scheitert im stabilen Lauf**, obwohl der Trockenlauf grün war (er beweist weder Token noch `is_stable`). Gegenmittel: `tap-check` vor M5. — **Ausgang:** offen bis zur Closure.
- **`integration-docker` auf dem Push-Tip ist rot oder am Zeitlimit** (`make test-docker` lief lokal nicht, Entscheidung des Projektinhabers). — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Folgt nach dem Tag (Tag-Hash, Lauf-IDs, M0 bis M7, Steering-Loop-Eintrag, Register, Folge-Slices, Risiken).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/plan/`, `docs/user/`, `docs/maintainer/` (DOK), dazu die Versionszeilen in `cmd/uboot/main.go` (keine Sub-Area, [`harness/conventions.md`](../../../../harness/conventions.md)), `repo.mk` und `Dockerfile` (WZG). Keine Sub-Area wird neu geschnitten.

**Vorgelagert — offene Beobachtungen sichten:** einschlägig für einen Cut: [`zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), [`ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md), [`unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md), [`host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md). Zähler stehen in den `state.md` und sind hier nicht wiederholt.

Alle berührten Sub-Areas GF.
