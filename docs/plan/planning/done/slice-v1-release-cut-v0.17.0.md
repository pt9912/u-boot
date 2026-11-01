# Slice V1: Release-Cut `v0.17.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** keine Welle-Datei — Minor-Release `v0.17.0` von `main` nach `v0.16.1` (Release-Version nach [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md); wie `v0.16.1` ohne Welle-Datei, geführt in der [Roadmap](../in-progress/roadmap.md) unter *Offene Wellen*, bis der Tag gesetzt ist). Der Release trägt, was seit dem Tag `v0.16.1` auf `main` liegt und in `CHANGELOG.md` unter `## [0.17.0]` steht: `--instance` für `u-boot devcontainer up`, `down` und `shell` ([`slice-v1-devcontainer-instanz-label`](../done/slice-v1-devcontainer-instanz-label.md)), die Referenz `CLAUDE_CODE_OAUTH_TOKEN` für `claude-code` ([`slice-v1-devcontainer-claude-oauth-token`](../done/slice-v1-devcontainer-claude-oauth-token.md)), das Lastenheft `0.26.0` ([`slice-cr-devcontainer-instanz`](../done/slice-cr-devcontainer-instanz.md), [`slice-cr-devcontainer-instanz-zusage-getestet`](../done/slice-cr-devcontainer-instanz-zusage-getestet.md)) und die Doku-Änderung zum Homebrew-Tap unter Linux. Weitere Slices im Release, ohne eigenen CHANGELOG-Eintrag: [`slice-v1-devcontainer-bau-lokal`](../done/slice-v1-devcontainer-bau-lokal.md) (`BUILD_MODE=local`, `scripts/bau.sh`, `.devcontainer/`, `repo.mk`, `Dockerfile`; Werkzeug ohne Wirkung für Nutzer des Binary), [`slice-v1-konfigdatei-name-als-konstante`](../done/slice-v1-konfigdatei-name-als-konstante.md) (`domain.ConfigFileName`, Wert `u-boot.yaml` unverändert; beide laut Review im Diff `v0.16.1..v0.17.0` gelesen), [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md), [`slice-v1-devcontainer-instanz-test-leicht`](../done/slice-v1-devcontainer-instanz-test-leicht.md), [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (nach dem Revert nur ein Kommentar in `devcontainer_shell.go`).

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

- [x] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt `## [0.17.0] - 2026-11-01` mit Lead, leerem Anker `[Unreleased]` und Vergleichslinks; Versionsstrings `0.17.0-dev` an den drei Stellen; beide READMEs und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.17.0`; das Datum steht in CHANGELOG und Handbuch-`Stand` gleich dem Tag-Datum. — Evidence: §3; Verification-Report Abschnitt 2 Punkt 1 (Stellen gelesen), Datum `2026-11-01` an den drei Stellen.
- [x] **Tag und Kontrolle:** `make gates`, `make ci`, `make fullbuild` und der Trockenlauf des Tap-Pfads sind vor dem Tag grün (§3); `ci` und `integration-docker` auf dem Push-Tip sind grün (M2); `tap-check` ist grün (M4); danach und nur nach ausdrücklicher Freigabe: Tag `v0.17.0` gesetzt und gepusht (M5); Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 (M6). — Evidence: §7 (*Läufe*) und Verification-Report Abschnitt 1. **Abweichung vom Wortlaut, ausgewiesen:** `make test-docker` lief lokal nicht (Entscheidung des Projektinhabers); `integration-docker` 41290870231 auf dem Push-Tip `success`, 0 `--- FAIL`. Der Trockenlauf des Tap-Pfads ist durch den echten Lauf des Jobs `tap` überholt, nicht nachgefahren (Verifier).
- [x] **Doku nachgezogen:** Roadmap (Zeiger unter *Offene Wellen* weg, Zeile in *Abgeschlossene Wellen* mit Tag-Hash und Datum, Ruhe-Marker), `releasing.md` (Zeiger „zuletzt“) und die Zeile in `open/README.md`. — Evidence: §7 (*Doku nachgezogen*).
- [x] `make gates` grün (nach dem letzten Commit vor dem Handoff und nach der Closure). — Handoff: §3 *Sensoren*; Verifier: Exit 0; Closure: *Verifikation* in §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review. — [Review](../../../reviews/2026-11-01-slice-v1-release-cut-v0.17.0.md) (0 HIGH, 1 MEDIUM, 5 LOW, 3 INFO; Nachlauf Fix-Runde: 1 LOW, 1 INFO) und [Verification](../../../reviews/2026-11-01-slice-v1-release-cut-v0.17.0-verify.md) (0 Befunde mit Blockwirkung).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. — §7.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung“ in §7 notiert. — §7 (*Beobachtungs-Register*).
- [x] Jedes Risiko aus §6 trägt einen Ausgang. — §6.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (kein Welle-Betrieb). — §7.

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
| `docs/maintainer/releasing.md` | kein Diff im Cut; update in der Fix-Runde `38805ed1` | Der Cut ließ `:56`–`:57` ("zuletzt" = `v0.16.1`) stehen; die Fix-Runde setzt den Zeiger auf diesen Slice (gelesen: `:56`–`:57`). |
| `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/open/README.md`, `open/slice-v1-devcontainer-ressourcenlimits.md`, `open/slice-v1-konfigdatei-alternativer-name.md` | update (Fix-Runde `38805ed1`) | Review F-1 bis F-6: Handbuch-Unterabschnitt *Mehrere Instanzen* (`benutzerhandbuch.md:2278`), Grenzliste "Not promised" im Lead, abgelaufene Marker `Vorschlag v0.17.0` in den zwei Entwürfen und in `open/README.md`. |
| `docs/reviews/2026-11-01-slice-v1-release-cut-v0.17.0.md`, `…-verify.md` | neu (Reviewer, Verifier) | Review- und Verification-Report. |
| `docs/plan/planning/observations/BEO-DOK/{plan-aussage-ueber-code,spec-bedingung-ohne-code-entsprechung}/` | update (Closure) | je eine Evidence-Datei `slice-v1-release-cut-v0.17.0.md` und eine Zeile in `state.md`. |
| `harness/README.md`, `tools/mk/devcontainer.mk` | nicht Teil des Release | Commit `7bbb8185` liegt nach dem Tag `v0.17.0` (`git diff --stat 38805ed1 HEAD`: nur diese zwei Dateien); er ist weder in `v0.17.0` enthalten noch Inhalt dieses Slice. |
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

**M0:** `date +%F` = `2026-11-01` (Systemuhr UTC, NTP-synchron); CHANGELOG-Überschrift `[0.17.0]` und Handbuch-`Stand` tragen dasselbe Datum. Der Hauptagent wiederholt M0 vor M1 und vor M5.

Lauf auf dem Stand nach dem Move nach `in-progress/` (Docker-Modus, Exit-Codes der `make`-Prozesse):

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0 (zweimal gelaufen: vor `ci` und nach dem Füllen dieser Tabelle); `coverage-gate` 92,20 % gegen 90 %; `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge; `docs-check` und `baseline-verify` Teil des Ziels, grün |
| `make ci` | Exit 0: `govulncheck` und `image-scan`: `u-boot:latest (debian 12.15)` 0 Schwachstellen, `usr/local/bin/u-boot` 0 |
| `make fullbuild` | Exit 0 (`ci` + Runtime-Image-Build) |
| `make build VERSION=0.17.0` | Exit 0; `docker run --rm u-boot:latest --version` = `u-boot version 0.17.0`; OCI-Label `org.opencontainers.image.version` = `0.17.0` (lokal, nicht gepusht). Der Standardbau von `make fullbuild` meldete `u-boot version 0.17.0-dev` |
| `make doc-immutable RANGE="v0.16.1..HEAD"` | Exit 0, 980 Dateien, 0 Befunde |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0, 980 Dateien, 0 Befunde |
| `make doc-planning` | Exit 0, 980 Dateien, 0 Befunde |
| `make test-docker` | **nicht gelaufen**, Entscheidung des Projektinhabers (2026-11-01): der CI-Lauf `integration-docker` auf `main` ist der Nachweis (M2). Das Zählen der `--- SKIP` und `--- PASS` des Tips gehört dem Hauptagenten |
| `make verify-depguard` | **nicht gelaufen**: der Cut berührt keine Importregel |
| `gh workflow run tap-check` | **nicht gelaufen**: startet der Auftraggeber mit dem PAT (M4) |

**M-T ausgeführt (Trockenlauf des Tap-Pfads).** Befehl: `docker run -i --rm --network none -v /Development/u-boot:/src:ro u-boot:test bash < <Skript>` (Image aus `make gates`, Arbeitsbaum schreibgeschützt; Skript im Scratchpad, nicht im Repo; Fixture nach `sumsFull` in `cmd/uboot/homebrew_scripts_test.go`; Stellvertreter für `curl` und `git` im `PATH`; kein echtes Token, `TAP_TOKEN=dummy`). Exit-Code des `docker run`: 0. Ausgabe (gekürzt auf die Aussagezeilen):

```text
(1)  homebrew-formula-fill.sh v0.17.0 <sums> scripts/homebrew-formula.rb.tmpl <Ausgabe>   exit=0   (vier url-Zeilen mit v0.17.0)
(2a) tap-formula-check.sh v0.17.0 <Ausgabe>   exit=0
(2b) tap-formula-check.sh v0.16.1 <Ausgabe>   exit=1
     tap-nachzug: the u-boot.rb asset of v0.16.1 does not carry four download URLs of v0.16.1 (found 4 url lines, 0 of v0.16.1).
(3)  tap-nachzug.sh, TAP_TOKEN=dummy, TAG=v0.17.0:   exit=0
     Protokoll: git clone --depth 1 … · git diff --quiet -- Formula/u-boot.rb · git … add Formula/u-boot.rb · git … commit -m u-boot 0.17.0 · git push origin HEAD
     tap-nachzug: pt9912/homebrew-u-boot updated to v0.17.0.
(3b) tap-nachzug.sh mit TAG=v0.16.1 gegen dieselbe Formel: exit=1, 0 clone-Aufrufe, 0 push-Aufrufe
```

Rot gesehen: (2b) und (3b) sind die rote Gegenprobe der Prüfung auf den Tag. **Grenze, wie in [`releasing.md`](../../../maintainer/releasing.md) §4 benannt:** `is_stable == 'true'` und der Workflow-Schritt, der echte `curl` gegen das Release-Asset, der echte `git push` und das Token sind nicht geprobt; der Stellvertreter von `git` meldet immer eine Änderung. Befund: keiner; kein Skript geändert.

**Binary-Unterschied:** `git diff v0.16.1..HEAD --stat -- . ':!docs' ':!spec' ':!*.md'` nach dem Cut: 70 Dateien (3438 Einfügungen, 213 Löschungen). Die Aussage „Code unverändert“ steht deshalb nicht im Lead.

**Kommentar-Selbstcheck** ([`AGENTS.md`](../../../../AGENTS.md) §3.9): (a) die Kommentare in `repo.mk` und `Dockerfile` nennen den geltenden Wert `0.17.0-dev`; (b) der alte Wert `0.16.1-dev` steht im Baum nicht mehr außerhalb von Planung, Reviews und Messprotokollen (`grep -rIn '0\.16\.1-dev'`); (c) keine neue Funktion.

**Rot gesehen — Grenze:** Kein Test pinnt die Nummer; eine falsche Versionszeile färbt kein Gate. Die Wächter sind die Prüfung `--version` = Tag im Workflow `publish` (nach dem Push) und die lokale Gegenprobe: `make build VERSION=0.17.0` meldet `0.17.0`, der Standardbau `0.17.0-dev`. Für die Lead-Sätze gibt es keinen Wächter im Repo.

**Abgleich der Plan-Tabelle** (`git diff --stat f58dc581..HEAD`): `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, `open/README.md` und dieser Slice haben je einen Plan-Eintrag und einen Diff; kein Diff, gelesen: `docs/maintainer/releasing.md:56`–`:57` (Zeiger „zuletzt“ nennt `v0.16.1`), `.github/workflows/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `spec/lastenheft.md` Zeile 9 (`0.26.0`), `docs/plan/adr/` (`make doc-immutable` Exit 0).

### Nacharbeit nach dem Tag ([`releasing.md`](../../../maintainer/releasing.md) §7)

- [x] Roadmap: Zeile `v0.17.0` nach *Abgeschlossene Wellen* mit Tag-Hash und Datum; Zeiger unter *Offene Wellen* weg, Ruhe-Marker zurück.
- [x] Datum gegen `git log -1 --format=%cs v0.17.0` prüfen (CHANGELOG, Handbuch-`Stand`).
- [x] Job `tap` im Log des Laufs `publish` prüfen („updated to v0.17.0.“ oder „already at“, keine Warnung zum fehlenden Secret).
- [x] Zeiger „zuletzt“ in `releasing.md` auf diesen Slice.
- [x] Closure-Notiz (§7) mit Tag-Hash, Lauf-IDs und Sensoren; Slice nach `done/` (erst Inhalt, dann reiner `git mv`, [`AGENTS.md`](../../../../AGENTS.md) §3.3).

## 4. Trigger

**Start** (`open` → `in-progress`): die Nummer `v0.17.0` gilt (Auftrag); die Inhalts-Slices liegen in `done/` (gelesen); `Verantwortlich:` ist gesetzt; `in-progress/` enthält nur die Master-Dokumente (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): nicht vorgesehen.
- `in-progress` → `open` (blockiert): ein roter `image-scan` durch einen neuen Advisory (Go-Pin heben, eigener Vorgang davor); ein Befund des Trockenlaufs, der eine Änderung an Workflow oder Skript verlangt; ein `--- FAIL` oder Zeitlimit-Abbruch in `integration-docker` auf dem Push-Tip; ein roter `tap-check`.

## 5. Closure-Trigger

DoD vollständig: Tag `v0.17.0` gesetzt und gepusht (M5), `publish` und Job `tap` kontrolliert (M6), Trockenlauf mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag, Verification Evidence (Verifier). Der Slice geht **erst nach dem Tag** nach `done/`. Kein Tag ohne ausdrückliche Freigabe: bleibt sie aus, bleibt der Slice in `in-progress/` oder geht mit benanntem Grund nach `open/` zurück.

## 6. Risiken und offene Punkte

- **Der Lead sagt mehr, als der Vertrag zusagt** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). Gegenmittel: der Lead nennt die Zusage mit den Bedingungen von `spec/lastenheft.md` (Zeile 504: zwei Instanzen, Workspace-Volumes, Sandbox-Profil, direkt nach dem Start, Linux `amd64`, lokale Docker-Engine) und die Liste des Nicht-Zugesagten; der Reviewer liest ihn gegen das Lastenheft. — **Ausgang:** nicht eingetreten. Review und Verification lasen "Promised" gegen Zeile 504 des Lastenheft: wörtlich gedeckt, nichts Zugesagtes zu breit. Review F-2 (LOW): die Liste "Not promised" war kürzer als die des Lastenheft (engere Lesart, kein Verstoß; Verifier V-1 INFO); die Fix-Runde `38805ed1` hat "kalter Cache oder kalter Bau" und "Dauer und Platzbedarf des Starts" ergänzt (`CHANGELOG.md`, gelesen).
- **Der Lead lässt eine Nutzer-Änderung aus.** Gelesen: der Diff gegen `v0.16.1` im Produktcode besteht aus `--instance`, der Prüfung des Namens, der Referenz `CLAUDE_CODE_OAUTH_TOKEN` samt `doctor`-Regel und der Konstante `domain.ConfigFileName`; der Commit zu `shell` ist nach dem Revert ein Kommentar. Ein weiterer Nutzer-Effekt ist nicht gelesen worden. — **Ausgang:** nicht eingetreten. Review F-7 las den Produktcode-Diff `v0.16.1..v0.17.0` und fand keinen weiteren Effekt; jede Nutzer-Wirkung steht im CHANGELOG. Das Build-Werkzeug berührt die Stufen `build` und `runtime` nicht (Review F-7, `--version`-Gegenprobe §3).
- **Das Datum weicht vom Tag-Datum ab** (Mitternacht). Gegenmittel: M0. — **Ausgang:** nicht eingetreten. Cut, Push und Tag am `2026-11-01`; `git log -1 --format=%cs v0.17.0` = `2026-11-01` = CHANGELOG-Überschrift = Handbuch-`Stand` (Verification Abschnitt 1).
- **Ein neuer Advisory macht `image-scan` rot.** Gegenmittel: `make ci` lokal und `ci` auf dem Push-Tip. — **Ausgang:** nicht eingetreten. `make ci` lokal Exit 0 (§3 *Sensoren*); `ci` 41290870228 auf `f89424e1` `success` (Verification Abschnitt 1, Review F-9).
- **Der Job `tap` scheitert im stabilen Lauf**, obwohl der Trockenlauf grün war (er beweist weder Token noch `is_stable`). Gegenmittel: `tap-check` vor M5. — **Ausgang:** nicht eingetreten. `tap-check` 41290880081 `success`; Job `tap` im Lauf `publish` 41292635297 `success`, Log (Verifier gelesen): "tap-nachzug: pt9912/homebrew-u-boot updated to v0.17.0.", Tap-Commit `cb47bc3`, `TAP_TOKEN` gesetzt; Formel im Tap byte-gleich mit dem Release-Asset `u-boot.rb`. Die Reihenfolge `tap-check` vor `publish` stützt sich auf die Lauf-IDs, nicht auf die Zeitstempel (siehe §7, *GitHub-Zeitstempel*).
- **`integration-docker` auf dem Push-Tip ist rot oder am Zeitlimit** (`make test-docker` lief lokal nicht, Entscheidung des Projektinhabers). — **Ausgang:** nicht eingetreten. `integration-docker` 41290870231 auf dem Push-Tip `success`; Review F-9 las das Job-Log über die REST-API: 3471 `--- PASS`, 4 `--- SKIP` (alle nested Podman), 0 `--- FAIL`, Paket `internal/e2e` `ok` nach 930,9 s (Log-Zeiten 08:55:39 bis 09:12:07, rund 17 Minuten). `make test-docker` lokal bleibt ausgewiesen nicht gelaufen.

## 7. Closure-Notiz

**Evidence der Nacharbeit** (Implementer, nach dem Tag; Bericht des Hauptagenten, nicht vom Implementer erhoben):

- Tag `v0.17.0` auf `f89424e1` (`git rev-parse v0.17.0^{commit}`, vom Implementer gelesen), Tag-Datum `2026-11-01` (`git log -1 --format=%cs v0.17.0`); CHANGELOG-Überschrift und Handbuch-`Stand` tragen dasselbe Datum.
- Läufe: `ci` 41290870228, `integration-docker` 41290870231, `tap-check` 41290880081, `publish` 41292635297, Tap-Commit `cb47bc3`.
- M7 (Mac, `brew upgrade u-boot`, `u-boot --version` = `0.17.0`): nicht gelaufen, Aktion des Projektinhabers.
- Fix-Runde nach dem Tag zu Review und Verification: Handbuch (Unterabschnitt *Mehrere Instanzen*), Lead-Grenzliste, Zustandsmarker, zwei Entwürfe in `open/`, Inhaltsliste.

**Closure (Planner, 2026-11-01, nach M6; Nacharbeit nach [`releasing.md`](../../../maintainer/releasing.md) §7).**

- **Erhebung und Lesung.** M0 bis M6 hat der Hauptagent erhoben. Der Planner hat die Reports von Review und Verification, den Slice, `CHANGELOG.md` (Lead, Grenzliste), den Handbuch-Unterabschnitt *Mehrere Instanzen* (`benutzerhandbuch.md:2278`), die Roadmap und `releasing.md:56` gelesen; er hat keine Außenwirkung ausgelöst und Lauf-IDs, Digests und Tap-Commit nicht gegen GitHub neu erhoben (übernommen aus Verification Abschnitt 1 und Review F-9). Tag und Release hat er nicht angefasst.
- **Läufe:** Tag `v0.17.0` (leichtgewichtig, `git cat-file -t` = `commit`) auf `f89424e1`, Tag-Datum `2026-11-01`; `ci` 41290870228, `integration-docker` 41290870231, `tap-check` 41290880081, `publish` 41292635297 (Jobs `publish`, `tap`, Paket-Smokes deb/ubuntu und rpm/fedora `success`); Release mit zwölf Assets, kein Prerelease, kein Draft, publiziert 2026-11-01T09:21:23Z; Image `ghcr.io/pt9912/u-boot:0.17.0` und `:latest` mit demselben Digest `sha256:a18ffc55bf3c1a5a6fe21438e7eb94011501fa154904668f3cdd31bea7edfbaa`, `--version` = `0.17.0`; Tap-Commit `cb47bc3`, Formel byte-gleich mit dem Asset `u-boot.rb`. Quelle für alles in diesem Absatz: Verification Abschnitt 1 (vom Verifier gelesen) und Hauptagent-Bericht.
- **M5-Form:** Der Tag ist leichtgewichtig, im Unterschied zum annotierten Tag `v0.16.1`. Nach [`releasing.md`](../../../maintainer/releasing.md) §8 wird er nicht verschoben; der Unterschied ist kein Befund, die Form hat keinen Sensor.
- **Doku nachgezogen:** Roadmap-Zeile in *Abgeschlossene Wellen* (`roadmap.md:109`, Tag-Hash und Datum); `releasing.md:56` "zuletzt" nennt diesen Slice; `open/README.md` führt die Zeile des Slice (Ortsverweis nach dem Move durch `make slice-mv`, Zustand "geschlossen, liegt in `done/`" im Folge-Commit). Der Zeiger unter *Offene Wellen* und der Ruhe-Marker folgen nach dem Move, weil `Nichts in Arbeit.` nur bei leerem `in-progress/` gilt.
- **Was hat funktioniert:** Der Cut selbst hielt: Versionsstellen, Datum, Compare-Links und Lead-Bedingungen trugen die Gegenlesung (Review Negativbefunde, Verification Abschnitt 2). Die Aussage "Code des Binary geändert, deshalb kein 'unverändert'" (§1) hat der Review gegen den Produktcode-Diff bestätigt (F-7). M0 hielt das Datum; der Job `tap` lief im ersten stabilen Lauf, nachdem `tap-check` das Token geprüft hatte.
- **Was ging anders als geplant:** (a) Das Handbuch widersprach sich zu Instanzen (Review F-1, MEDIUM); der Cut hatte nur Kopf und Versionsstellen berührt. (b) Zwei Entwürfe in `open/` und `open/README.md` trugen den Marker "Vorschlag `v0.17.0`" über den Ablauf hinaus, und die Inhaltsliste nannte zwei Slices nicht (F-3 bis F-6). (c) Die Fix-Runde fiel nach dem Tag; sie berührt nur Dokumente auf `main`, nicht den Tag-Stand (siehe *Verifikation*, Diff gegen `v0.17.0`).
- **Befunde von Review und Verification, Ausgang je Befund:**

  | Befund | Ausgang |
  |---|---|
  | F-1 MEDIUM (Handbuch *Mehrere Instanzen* widerspricht dem Abschnitt zu `--instance`) | behoben in `38805ed1`; der Unterabschnitt sagt jetzt, dass `devcontainerId` den Projektordner und mit `--instance` den Namen einbezieht (`benutzerhandbuch.md:2278`–`2285`, gelesen); Beleg im Register; der Fix-Text trägt N-1 |
  | N-1 LOW (Nachlauf Fix-Runde, [Review](../../../reviews/2026-11-01-slice-v1-release-cut-v0.17.0.md) *Nachlauf Fix-Runde*: der Satz "Jede Instanz hat damit eigene Volumes" im Unterabschnitt *Mehrere Instanzen* nennt alle Volumes der Liste davor eigen; gepinnt sind nur die Workspace-Volumes, die Spec nimmt Docker-Speicher und Modelle aus) | **behoben** nach dem Cut auf `main`: der Satz sagt jetzt, dass sich die Namen der Volumes von Instanz zu Instanz unterscheiden (`benutzerhandbuch.md` Unterabschnitt *Mehrere Instanzen*), wie der Abschnitt zu `--instance` es für Docker-Speicher und Modelle schon tat. Der Release-Text (`CHANGELOG.md`, Lead) behauptet nicht mehr als die Spec. Beleg im Register (Evidence `spec-bedingung-ohne-code-entsprechung`) |
  | N-2 INFO (Nachlauf: `INSTANCE` unquotiert im Rezept von `tools/mk/devcontainer.mk`) | kein Handlungsbedarf (lokaler Komfortpfad; `u-boot` prüft den Namen, Exit `10`) |
  | F-2 LOW (Grenzliste im Lead kürzer als im Lastenheft) | behoben in `38805ed1` (zwei Punkte ergänzt, `CHANGELOG.md` gelesen); engere Lesart, kein Register-Beleg (siehe unten) |
  | F-3 LOW (Roadmap-Marker "Push und Tag stehen aus") | behoben in `38805ed1`; mit dieser Closure ist auch der Zeiger weg |
  | F-4, F-5 LOW (Marker/Zeiger "Vorschlag `v0.17.0`" in zwei Entwürfen und `open/README.md`) | behoben in `38805ed1` (`open/README.md`: "`v0.17.0` ist veröffentlicht; Zuordnung offen") |
  | F-6 LOW (Inhaltsliste unvollständig) | behoben in `38805ed1` (Roadmap-Zeiger nennt neun Slices) |
  | F-7 INFO (Binary geändert, Lead sagt nichts über Unveränderung) | kein Handlungsbedarf |
  | F-8 INFO (Handbuch-Version `2.1`) | begründet belassen |
  | F-9 INFO (CI-Nachweis gegengelesen) | bestätigt, 0 `--- FAIL` |
  | V-1 LOW (Roadmap-Marker nach dem Tag) | mit F-3 behoben |
  | V-1 INFO (Grenzliste) | mit F-2 behoben |
  | V-2 INFO (API-Zeitstempel widersprüchlich) | Beobachtung ohne Register-Dach (siehe unten) |
  | V-3 INFO (zwei `.deb` nicht heruntergeladen) | weiter offen: Paket-Smoke `deb, ubuntu` `success`, die Digests stehen in `SHA256SUMS`; die Prüfung der zwei Dateien selbst ist nicht nachgeholt |

- **Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`):
  - *Lehre 1 (Muster, nicht verkörpert):* **Ein Handbuch-Unterabschnitt blieb beim Cut mit altem Wortlaut stehen.** "Derselbe Ordner bleibt eine Instanz" und "Dev-Containers-Werkzeug löst `devcontainerId` pro Projektordner auf" waren vor `--instance` wahr; mit dem Release gelten sie nicht mehr, und der Abschnitt zu `--instance` steht im selben Handbuch. Der Cut liest Kopf, Versionsstellen und Gültigkeitsbereich, nicht die Handbuchabschnitte, die der Release ausliefert, gegeneinander. Klasse: Aussage über Verhalten des Codes ohne pinnenden Test (`BEO-DOK/spec-bedingung-ohne-code-entsprechung`, Altbestand, durch den neuen Code überholt) und Cut ohne Abgleich der ausgelieferten Handbuchstellen (`BEO-DOK/plan-aussage-ueber-code`). Gefangen hat es der Review, kein Gate. Ob daraus ein Pflichtschritt "ausgelieferte Handbuchabschnitte gegen den Lead lesen" in [`releasing.md`](../../../maintainer/releasing.md) §3 wird, entscheidet der Projektinhaber; verkörpert wird nichts.
  - *Lehre 2 (Muster, nicht verkörpert):* Zustandsmarker mit Ablaufpunkt ("Push und Tag stehen aus", "Vorschlag `v0.17.0`") standen nach dem Tag in Roadmap, `open/README.md` und zwei Entwürfen. Gleiche Klasse wie `BEO-DOK/plan-aussage-ueber-code` (Rückfall nach der Verkörperung in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
  - *Neuer Sensor:* keiner. Kein Gate liest Handbuchaussagen gegeneinander oder gegen den Lead. *Benannte Spec-Lücke:* keine.
- **Beobachtungs-Register (`../observations/`):** Kein Zähler gesetzt; die Zahl folgt aus den Dateien.
  - *Evidence neu (ein Vorgang zählt einmal):* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (`evidence/slice-v1-release-cut-v0.17.0.md`: Zustandsmarker nach Ablauf, Inhaltsliste und Entwurf-Zeiger, Review F-3 bis F-6, Verification V-1; behoben) · [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (`evidence/slice-v1-release-cut-v0.17.0.md`: Handbuch-Aussage über `devcontainerId` ohne pinnenden Test, vom neuen Code überholt, Review F-1; behoben).
  - *Nicht gezählt:* [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md): F-2 ist die engere Lesart (der Lead sagt weniger als das Lastenheft, nicht mehr); Verifier V-1 wertet es als INFO. [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md): Prüfweg über REST-API und Annotationen angewandt (Verification Zeile `tap`: der Quelltext der Warnung im Run-Block ist nicht die ausgegebene Warnung), kein Fund. [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md): `tap-check` lief vor dem Tag, kein Fund. [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md): Review und Verification melden keinen Fund; die Closure selbst hat ein Python-Skript (Abhaken der DoD-Kästchen, Skript im Scratchpad) und ein `sed -i` (eine Zeile in einer `state.md`) auf dem Host gefahren, zwei Vorgänge dieser Klasse, die kein Reviewer gesehen hat: hier ausgewiesen, nicht als Evidence angelegt, weil der Befund nicht durch einen Report belegt ist; das Ergebnis ist im Diff gelesen.
  - *GitHub-Zeitstempel (V-2, Review F-9):* `gh run list` zeigt "0s", `completedAt` vor `startedAt`, `updatedAt` vor `createdAt`; die Log-Zeiten tragen die Dauer. Ein Register-Dach dafür gibt es nicht (gesucht: `observations/` nach "Zeitstempel" und "createdAt"; `BEO-WZG/messbeleg-fehlt-oder-kommt-nach-dem-lauf` ist ein anderer Gegenstand). Beobachtung ohne Eintrag: Lauf-Reihenfolge nur über Lauf-IDs und Log-Zeiten belegen, nicht über `createdAt`/`updatedAt`. Ein erster Vorgang; eine neue Beobachtung legt die Closure nicht an.
- **Carveouts:** keiner neu, keiner gelöst, keiner verändert.
- **Folge-Kandidaten (keine Dateien):** Pflichtschritt in [`releasing.md`](../../../maintainer/releasing.md) §3 "ausgelieferte Handbuchabschnitte gegeneinander und gegen den Lead lesen" (Lehre 1); Prüfung der zwei `.deb` gegen `SHA256SUMS` (V-3). Die Anlage entscheidet der Projektinhaber; diese Notiz verspricht keinen Slice.
- **Folge-Slices** (Dateien): [`slice-v1-devcontainer-ressourcenlimits`](../open/slice-v1-devcontainer-ressourcenlimits.md) und [`slice-v1-konfigdatei-alternativer-name`](../open/slice-v1-konfigdatei-alternativer-name.md) liegen als Entwürfe in `open/`, ohne Zuordnung zu einem Release; kein neuer Slice aus dieser Closure.
- **Weiter offen (nicht abgehakt):** **M7**: Mac-Test `brew tap pt9912/u-boot`, `brew install` beziehungsweise `brew upgrade u-boot`, `u-boot --version` = `0.17.0`. **Nicht gelaufen**, Aktion des Projektinhabers (§3, Tabelle M7); der Cut und der Job `tap` ersetzen die Bestätigung nicht. **V-3**: die zwei `.deb`-Dateien nicht heruntergeladen. **Nicht belegt außerdem:** `make test-docker` lokal (Entscheidung des Projektinhabers; CI-Lauf ist der Nachweis), `make verify-depguard` (keine Importregel berührt).
- **Risiken aus §6:** je mit Ausgang in §6 (alle sechs nicht eingetreten).
- **Verifikation** (Läufe des Planners am 2026-11-01, Exit-Code des `make`-Prozesses; nach Closure-Inhalt, Move, Link-Nachzug, Roadmap-Ruhe-Marker und `open/README.md`, Docker-Modus): `make gates` Exit 0 (`coverage-gate` 92,20 % gegen 90 %; `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge; `lint`, `docs-check` und `baseline-verify` Teil des Ziels, Zeilen im Log nicht einzeln ausgelesen). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa. Nicht gelaufen: `make ci`, `make fullbuild`, `make test-docker`, `make verify-depguard` (die Closure ändert nur Planungs-Dokumente, Register, Verweise und Roadmap; die Läufe vor dem Tag stehen in §3 und im Verification-Report).
- **Abgleich der Plan-Tabelle** (`git diff --stat f58dc581 HEAD` auf dem Arbeitsbaum der Closure, nach der Fix-Runde `38805ed1`): jede Datei hat einen Eintrag in §3 oder einen ausgewiesenen Nicht-Eintrag (`harness/README.md`, `tools/mk/devcontainer.mk` aus `7bbb8185` nach dem Tag, nicht Teil des Release); `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `benutzerhandbuch.md` (Cut und Fix-Runde), `releasing.md` (Zeiger), `roadmap.md`, `open/README.md`, die zwei Entwürfe, dieser Slice, die zwei Reports, im Register zwei Evidence-Dateien und zwei `state.md`. Die Zeilen "kein Diff" sind durch die Lesung in Review (Negativbefund Plan-Tabelle) und Verification Abschnitt 6 belegt (`spec/lastenheft.md` Zeile 9 `0.26.0`; `releasing.md:56` vor der Fix-Runde "zuletzt" `v0.16.1`).
- **Drei Paarungen** (geprüft auf den Inhalt dieser Notiz): Anker — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert, nichts zu paaren. Folge-Slice — kein Slice versprochen; die zwei genannten Entwürfe liegen als Dateien in `open/` (`ls` bei der Closure). Register — jede genannte Kennung `BEO-<KUERZEL>/<slug>` ist ein Verzeichnis mit nicht leerem `evidence/` (`ls` bei der Closure, nach den Evidence-Dateien dieser Closure).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/plan/`, `docs/user/`, `docs/maintainer/` (DOK), dazu die Versionszeilen in `cmd/uboot/main.go` (keine Sub-Area, [`harness/conventions.md`](../../../../harness/conventions.md)), `repo.mk` und `Dockerfile` (WZG). Keine Sub-Area wird neu geschnitten.

**Vorgelagert — offene Beobachtungen sichten:** einschlägig für einen Cut: [`zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), [`ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md), [`unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md), [`host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md). Zähler stehen in den `state.md` und sind hier nicht wiederholt.

Alle berührten Sub-Areas GF.
