# Verification-Report: slice-v1-release-cut-v0.24.0 — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)

**Gegenstand:** [`slice-v1-release-cut-v0.24.0`](../plan/planning/done/slice-v1-release-cut-v0.24.0.md), Commit-Spanne `e38d6461..8aca6bef` (Slice in Umsetzung bis lokaler Tip). Der Cut-Commit ist `58411aaf` (8 Dateien, 28 Einfügungen, 19 Löschungen); `8aca6bef` ändert nur die Lead-Zeile des CHANGELOG-Blocks (1 Zeile) und legt den Review-Report ab (`git diff --stat 58411aaf..HEAD`: 2 Dateien). **Abgleichstand nach der letzten Fix-Runde = `8aca6bef`**; der Arbeitsbaum war vor und nach allen Läufen sauber (`git status --short` leer).

**Push-Stand (gelesen am 2026-11-06):** `origin/main` steht auf `68321ff6`; `main` liegt sieben Commits davor (`git log origin/main..HEAD`: `6d68a2a3` bis `8aca6bef`). Es gibt keinen Lauf von `ci`, `integration-docker` oder `tap-check` auf `8aca6bef` (`gh run list --commit 8aca6bef…` liefert `[]`); der Tag `v0.24.0` existiert weder lokal (`git tag -l 'v0.24*'` leer) noch ist er gepusht. Der Push und der Tag stehen nach Plan §1 aus und gehören dem Projektinhaber.

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel, Bedingung für Push und Tag, M-Start (a)/(b)), §2 (DoD), §3 (Plan-Tabelle, Reihenfolge der Maintainer-Aktionen), §6 (Risiken 1 bis 10)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md), [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), [`ADR-0047`](../plan/adr/0047-brew-formeln-lastenheft-stellen-und-verteilung-lastenheft-spezifikation.md)
- Review-Report [`2026-11-06-slice-v1-release-cut-v0.24.0`](2026-11-06-slice-v1-release-cut-v0.24.0.md) (0 HIGH / 1 MEDIUM / 1 LOW; Eingang, nicht Evidence)
- [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §4 und §5; [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.6, §3.7, §3.9, §3.10, §5)
- Vorbild der Form: [`verification-2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo`](verification-2026-11-06-slice-v1-devcontainer-brew-ripgrep-im-repo.md)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-release-cut-v0.24.0`](../plan/planning/done/slice-v1-release-cut-v0.24.0.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md), [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md)
- Artefakte: `CHANGELOG.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md`

### DoD-Abgleich (Plan §2)

Die Häkchen im Slice stehen alle offen; sie gehören dem Planner. Der Zustand je Punkt:

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1a | `CHANGELOG.md`: `## [0.24.0] - 2026-11-06`, Lead „Twenty-fourth release.“, leeres `[Unreleased]`, Compare-Links | **erbracht** | `CHANGELOG.md` Zeile 12 (`## [Unreleased]`, leer), Zeile 14 (`## [0.24.0] - 2026-11-06`), Zeile 16 (Lead); Zeilen 1670 und 1671: `[Unreleased]: …compare/v0.24.0...HEAD`, `[0.24.0]: …compare/v0.23.0...v0.24.0`. Die Zeilennummern 1662/1663 des Plans sind durch den Block verschoben (Plan-Stand, kein Befund) |
| 1b | Unter `### Changed` der Eintrag zur Gate-Reihenfolge mit Zusage und Grenze unter `make -j`; nichts zu `BUILD_MODE` und zum eigenen Devcontainer mit `ripgrep` | **erbracht** | Eintrag „The gate run of the repository checks its tool chain first“ mit **Promised** (Stellung in der Prerequisite-Liste von `record-gates`, rote Probe des Wächter-Tests) und **Not promised** (Reihenfolge unter `make -j`). `grep -c -i -w 'BUILD_MODE'` über die CHANGELOG-Zeilen 13 bis 24 (Block `[0.24.0]` bis zum Ende der *Added*-Einträge): 0 Treffer; der *Changed*-Eintrag (Zeile 25) wurde gelesen und nennt `BUILD_MODE` nicht. Der Eintrag ist durch den eigenen Lauf gedeckt: `make gates` Zeile 3 `gate-reihenfolge: tool-pruefung is the first check of record-gates — OK`, Zeilen 4 bis 6 Wächter-Test `2 Faelle, 0 Fehlschlaege` (grün auf Kopie, rot bei `+=`) |
| 1c | Lead nennt alle drei Einträge und sagt nur, was Einträge und Diff belegen | **erbracht, mit Beleg zur Fix-Runde** | Der Lead nennt `devcontainer remove`, `brew` und die Gate-Reihenfolge. Die zwei Review-Findings sind in `8aca6bef` aufgelöst, siehe *Prüfung der Review-Übergaben*. `8aca6bef` hat **kein Nachlauf-Review**; ich habe die geänderte Lead-Zeile selbst gegen den *Added*-Eintrag, die Lastenheft-Aussage und die Tests gelesen (V-2) |
| 1d | Versionsstrings `0.24.0-dev` an den drei Stellen samt Kommentaren | **erbracht** | `git diff e38d6461..HEAD`: `cmd/uboot/main.go` `var version = "0.24.0-dev"`; `repo.mk` Kommentar (zwei Stellen) und `VERSION ?= 0.24.0-dev`; `Dockerfile` Kommentar (zwei Stellen) und `ARG UBOOT_VERSION=0.24.0-dev`. Suche nach `0.23.0` außerhalb von `CHANGELOG.md`, `docs/plan`, `docs/messungen`, `docs/reviews`, `.harness`: einzige Fundstelle ist die Versionstabelle von `spec/lastenheft.md` (Eintrag 0.23.0 des Lastenhefts, kein Release-Bezug) |
| 1e | READMEs und Handbuch auf `v0.24.0`; Handbuch-`Stand` = Tag-Datum | **erbracht, Tag-Datum noch ungeprüft** | `README.md` Status-Zeile und `VERSION=`, `README.de.md` ebenso; Handbuch: Kopf (`Software-Version: v0.24.0`, `Stand: 2026-11-06`), `u-boot version 0.24.0`, `download/v0.24.0`, `VERSION=0.24.0`, Gültigkeitsbereich. `date +%F` = `2026-11-06` = Datum von CHANGELOG, Handbuch und `git log -1 --format=%cs HEAD`. Das Tag-Datum ist erst nach dem Tag lesbar (Risiko 3, offen) |
| 2a | `make gates` lokal grün (Devcontainer-Container) | **erbracht, eigener Lauf** | siehe *Sensors* |
| 2b | `make ci` lokal grün | **erbracht, eigener Lauf** | siehe *Sensors* |
| 2c | Gegenprobe `make build VERSION=0.24.0`: `--version` und OCI-Label melden `0.24.0` | **erbracht, eigener Lauf** | siehe *Sensors* |
| 2d | Tap-Trockenlauf nach `releasing.md` §4 mit Gegenprobe zum Tag des Vorgängers | **erbracht, eigener Lauf** | siehe *Tap-Trockenlauf* |
| 2e | `ci` und `integration-docker` auf dem Push-Tip grün, `tap-check` grün | **OFFEN — nicht erbracht** | Der Push steht aus; es gibt keinen Lauf auf `8aca6bef`. Die Läufe auf `origin/main` (`68321ff6`) belegen einen älteren Stand und zählen nicht (V-1) |
| 2f | Tag `v0.24.0` auf dem Push-Tip, `publish` mit Job `tap` kontrolliert (zwölf Assets, Image `0.24.0` und `:latest`, `tap-nachzug: … updated to v0.24.0`) | **OFFEN — nicht erbracht** | Bedingt durch zwei getrennte Freigaben des Projektinhabers und durch 2e |
| 3a | Roadmap: *Abgeschlossene Wellen* mit Tag-Hash und Datum, Zustandsmarker; Bestand-Zeile | **OFFEN — nach dem Tag** | `roadmap.md` Zeile 45 führt den Slice unter *Offene Wellen*, `open/README.md` Zeile 14 den Bestand — der Zustand vor dem Tag, wie der Plan ihn vorsieht |
| 3b | `releasing.md` §4: Kommentarzeile zu `make gates` nennt `tool-pruefung` | **erbracht** | `git diff e38d6461..HEAD -- docs/maintainer/releasing.md`: `… + ablage-pruefung + tool-pruefung`; deckt sich mit der Gate-Zeile in `harness/README.md` |
| 3c | Datum in CHANGELOG und Handbuch stimmt nach dem Tag mit `git log -1 --format=%cs v0.24.0` überein | **OFFEN — nach dem Tag** | Risiko 3 |

Pflichten, die nicht mitzählen (Plan §2):

| Pflicht | Zustand | Evidence |
|---|---|---|
| `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure | **erbracht für `8aca6bef`; „nach der Closure“ offen** | eigener Lauf am Tip, Exit `0` (*Sensors*) |
| Review (kein Self-Review, Report unter `docs/reviews/`) auf dem Cut-Commit vor dem Push | **erbracht** | [Review-Report](2026-11-06-slice-v1-release-cut-v0.24.0.md) liegt vor, Gegenstand `58411aaf`, vor dem Push. Verification mit Report: dieser Bericht |
| Closure-Notiz mit Lerneintrag, Tag-Hash, Lauf-IDs; Register; Risiko-Ausgänge; drei Paarungen | **offen** | §7 des Slice trägt den Platzhalter „Wird bei der Closure geschrieben“ — richtig vor dem Tag |

### Prüfung der Review-Übergaben

| Finding | Auflösung | Beleg |
|---|---|---|
| F-1 (MEDIUM, Lead nennt Exit `10` ohne die Bedingung „es gibt einen Container“) | **behoben in `8aca6bef`** | Der Lead sagt jetzt „exit code `10` without `--yes` in non-interactive mode and under `--json`, when there is a container to remove; with no container the command asks nothing and ends with exit code `0`“. Gegen den *Added*-Eintrag gelesen („With no container to remove the command asks nothing and ends with exit 0, also `--volumes` without `--yes` in non-interactive mode“): deckt sich. Die Bedingung ist gepinnt: `internal/hexagon/application/devcontainer_remove_test.go` Unterfall `no container` (`Volumes: true, NonInteractive: true` → kein Fehler, `Removed` falsch, kein `removeDevcontainerCalls`) und `TestDevcontainerRemove_ReleaseTable` (mit Container: `no-interactive` und `json` ohne `--yes` → `ErrConfirmationRequired`); CLI-Test `internal/adapter/driving/cli/devcontainer_remove_test.go` Zeile 319 („with no container the command ends with exit 0 and asks nothing“); Docker-Test `TestDevcontainerRemove_Docker_NoContainer` (`internal/e2e/devcontainer_remove_docker_test.go` Zeile 485; läuft erst als `integration-docker`) |
| F-2 (LOW, Lead ohne „serial“) | **behoben in `8aca6bef`** | Der Lead sagt jetzt „`tool-pruefung` stands before all other checks of the serial `make gates`, not promised under `make -j`“; deckt sich mit dem *Changed*-Eintrag und dem Kommentar in `tools/mk/tool-pruefung.mk` |

### Prüfung des Risikos 2 (keine Aussage über `rg` im Container)

`grep -i -w 'rg\|ripgrep'` über die hinzugefügten Zeilen von `git diff e38d6461..HEAD` (ohne `docs/reviews/`): 0 Treffer; über die Lead-Zeile 16 des CHANGELOG: 0 Treffer. Das Wort `ripgrep` steht im Block nur in der Testformel der *Added*-Zeile zu `brew` („the build with `ripgrep` and a tap formula runs in both profiles“) — sie war vor dem Cut Teil von `[Unreleased]` und ist im Cut nur verschoben. Der Cut sagt nichts über den Container dieses Repos aus. Der Nachweis bleibt beim Projektinhaber im Folge-Slice [`slice-v1-devcontainer-brew-ripgrep-nachweis`](../plan/planning/open/slice-v1-devcontainer-brew-ripgrep-nachweis.md).

### Entscheidungs-Konformität (Auftragspunkt: ADR und Spec)

| Quelle | Ergebnis | Beleg |
|---|---|---|
| Plan §1: „der Cut-Commit ändert in `cmd/` und `internal/` allein die Versionszeile“ | **konform** | `git show 58411aaf --stat`: die einzige Datei unter `cmd/` oder `internal/` ist `cmd/uboot/main.go` (1 Zeile). Der Diff `v0.23.0..HEAD` über Nicht-Test-Dateien trägt 27 Dateien (`git diff v0.23.0..HEAD --stat -- cmd internal ':!*_test.go'`), wie im Plan gelesen |
| Plan §3, Zeile „kein Diff durch den Cut“ (`spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `d-check.mk`, `.golangci.yml`, `go.mod`) | **konform** | `git diff e38d6461..HEAD --stat` über diese Pfade und `internal`: leer. `git diff v0.23.0..HEAD -- Makefile d-check.mk go.mod`: 0 Zeilen. (`.github/workflows/ci.yml` trägt seit `v0.23.0` eine Zeile weniger; sie stammt nicht aus dem Cut und steht im Plan §1 (b) als `BUILD_MODE`-Zeile) |
| [`AGENTS.md`](../../AGENTS.md) §3.6 (keine Gate-Lockerung, kein Pin-Bump) | **bestätigt** | `golangci`, Coverage-Schwelle 90 % und `.d-check.yml` unberührt; Coverage 92,40 % |
| [`AGENTS.md`](../../AGENTS.md) §3.2 (Suppression) | **bestätigt** | kein `//nolint` im Diff |
| [`AGENTS.md`](../../AGENTS.md) §3.5 (Accepted-ADRs unverändert) | **bestätigt** | keine Datei unter `docs/plan/adr/` im Diff |
| [`AGENTS.md`](../../AGENTS.md) §5 Regel 6 (Nutzer-Doku ohne Kennungen) | **bestätigt** | die Diffs von `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` tragen ausschließlich Versionsstrings und das Datum |
| [`AGENTS.md`](../../AGENTS.md) §3.9 (Kommentare) | **bestätigt** | geänderte Kommentare in `repo.mk` und `Dockerfile` nennen die geltende Zusage; der alte Wert `0.23.0-dev` steht nirgends im Baum |
| [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only) | **bestätigt** | alle Nachweise sind `make`-Targets im Devcontainer-Container; kein Host-`go test` |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) / [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) (Changed-Eintrag) | **konform** | `tool-pruefung` als erster Check von `record-gates` im eigenen Lauf belegt (DoD 1b); die Pin-Gruppen laufen grün (`d-check: pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`) |
| Lastenheft-Bezug des Leads („Lastenheft 0.36.0 to 0.37.3“) | **konform** | `spec/lastenheft.md` Zeile 9: `0.37.3`; Änderungstabelle führt 0.36.0, 0.36.1, 0.37.0 bis 0.37.3 |

### Plan-vs-Code (§3 gegen `git diff --stat e38d6461..HEAD`, Stand nach der letzten Fix-Runde)

| Diff-Datei | Plan-Zeile (§3) | Abgleich |
|---|---|---|
| `CHANGELOG.md` (+10/−1 über beide Commits) | `CHANGELOG.md` update | deckt (Block, Lead, *Changed*-Eintrag, zwei Compare-Links) |
| `cmd/uboot/main.go` (+1/−1) | update | deckt |
| `repo.mk` (+3/−3) | update | deckt (Zeilen für Kommentar und `VERSION ?=`) |
| `Dockerfile` (+3/−3) | update | deckt |
| `README.md`, `README.de.md` (je +2/−2) | update | deckt |
| `docs/user/benutzerhandbuch.md` (+6/−6) | update | deckt (Kopf, `Stand`, `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich) |
| `docs/maintainer/releasing.md` (+1/−1) | update | deckt |
| [Review-Report](2026-11-06-slice-v1-release-cut-v0.24.0.md) (+80, neu) | `docs/reviews/<Datum>-slice-…` neu | deckt |
| dieser Report (neu) | `docs/reviews/verification-<Datum>-slice-…` neu | deckt |
| `roadmap.md`, `open/README.md`, Folge-Slice, Evidence-Datei | Plan-Zeilen „mit dem Plan-Commit“ bzw. Zeile unter *Offene Wellen* | nicht im Diff `e38d6461..HEAD` (sie kamen mit dem Plan-Commit und `e38d6461`); `roadmap.md` Zeile 45 und `open/README.md` Zeilen 14 bis 16 gelesen: Zustand stimmt mit dem Plan überein |
| als „kein Diff“ geführt: `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `d-check.mk`, `.golangci.yml`, `go.mod` | „kein Diff durch den Cut“ | im Diff nicht enthalten, deckt |

Befund zur Tabelle: keine Zeile ohne Plan-Eintrag, kein Plan-Eintrag ohne Diff (außer den nach dem Tag fälligen Zeilen: Roadmap-Wechsel nach *Abgeschlossene Wellen*, Slice-Move nach `done/`, Beobachtungs-Register).

### Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer-Container, `/.dockerenv` vorhanden, Docker 29.9.3 verschachtelt; Tip `8aca6bef`) | **pass** | Exit `0`. `tool-pruefung: … (go 1.27.3, golangci-lint 2.14.0, make >= 4.3) … — OK`; `gate-reihenfolge: tool-pruefung is the first check of record-gates — OK`; `gate-reihenfolge-test: 2 Faelle, 0 Fehlschlaege`; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`; `d-check: 1198 Datei(en) geprüft, 0 Befund(e)`; `0 issues.` (lint); `coverage-gate: OK — coverage 92.40% meets threshold 90%`; `ablage-pruefung: OK`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`. Der Lauf endete mit `record-gates` ohne Stempelfehler |
| `make ci` (gleicher Kontext und Tip) | **pass** | Exit `0`. `gates` wie oben; `govulncheck` (`go install golang.org/x/vuln/cmd/govulncheck@v1.8.0`): `No vulnerabilities found.`; `image-scan` (Trivy): `u-boot:latest (debian 12.15)` 0 Vulnerabilities, `usr/local/bin/u-boot` (gobinary) 0 Vulnerabilities; Schlusszeile `[ci] gates + govulncheck + image-scan green` |
| `make build VERSION=0.24.0` mit Gegenprobe | **pass** | Exit `0`. `docker run --rm u-boot:latest --version` → `u-boot version 0.24.0` (Exit `0`); `docker inspect u-boot:latest --format '{{index .Config.Labels "org.opencontainers.image.version"}}'` → `0.24.0`. Das Image trägt danach den Tag `latest` mit der Version `0.24.0`; es ist lokal, nicht gepusht |
| Tap-Trockenlauf | **pass** | siehe *Tap-Trockenlauf* |
| `make docs-check` nach Anlage dieses Reports | siehe Abschnitt *Sensor-Lauf* | — |
| Lead gegen Einträge und Diff gelesen (M-L) | **pass** | Abschnitt *Prüfung der Review-Übergaben* und *Prüfung des Risikos 2* |
| `make test-docker` | not run | Plan §1 schließt ihn aus (verschachtelter Daemon trägt die Docker-Tests nicht gleichwertig); der Nachweis ist `integration-docker` auf dem Push-Tip — er steht aus (2e). Der Diff des Cuts enthält keinen Go-Code außer der Versionszeile |
| `make verify-depguard` | not run | keine Importregel berührt |
| `make fullbuild` | not run | `make ci` plus der Runtime-Image-Build durch `make build VERSION=0.24.0` sind gefahren; `fullbuild` ist keine DoD-Forderung des Plans |
| `gh workflow run tap-check` | not run | Plan §3 Schritt 6 (M4) setzt den Push voraus; das Schreibrecht des Tokens belegt nur dieser Lauf, nicht der Trockenlauf (`releasing.md` §4) |

### Tap-Trockenlauf (DoD 2d)

Ausführung nach `releasing.md` §4: im Image `u-boot:test-docker-tools` (enthält `bash`), `--network none`, Arbeitsbaum als schreibgeschützter Mount (`-v "$PWD:/repo:ro"`), Fixture `sumsFull` (Muster aus `cmd/uboot/homebrew_scripts_test.go`), Stellvertreter für `curl` (kopiert die Fixture-Formel an das Ziel von `-o`) und `git` (protokolliert, schreibt nichts); `TAP_TOKEN=dummy`. Das Image `u-boot:test` aus `make test` stand nicht zur Verfügung; die Wahl des Images ändert den Beleg nicht, denn der Lauf braucht nur `bash` und Standardwerkzeuge.

| Schritt | Erwartet | Gelesen |
|---|---|---|
| `homebrew-formula-fill.sh v0.24.0 … ` | Exit `0` | Exit `0`: `u-boot.rb written (tag v0.24.0, four platform digests filled)` |
| `tap-formula-check.sh v0.24.0 u-boot.rb` | Exit `0` | Exit `0` |
| `tap-formula-check.sh v0.23.0 u-boot.rb` (Tag des Vorgängers) | Exit `1` | Exit `1`: `does not carry four download URLs of v0.23.0 (found 4 url lines, 0 of v0.23.0)` |
| `tap-nachzug.sh`, `TAG=v0.24.0` | Protokoll `clone`, `add`, `commit -m u-boot 0.24.0`, `push` | Exit `0`; Protokoll: `clone --depth 1 …homebrew-u-boot.git`, `diff --quiet`, `add Formula/u-boot.rb`, `commit -m u-boot 0.24.0`, `push origin HEAD`; Schlusszeile `pt9912/homebrew-u-boot updated to v0.24.0` |
| `tap-nachzug.sh`, `TAG=v0.23.0` | Abbruch vor `clone` | Exit `1`; Protokoll hat 0 Zeilen (kein `git`-Aufruf) |

**Er beweist nicht** (Wortlaut `releasing.md` §4): die Bedingung `is_stable == 'true'` und den Schritt im Workflow, den echten `curl` gegen das Release-Asset, den echten `git push` und das Token. Das Token belegt allein `tap-check` (offen, 2e).

### Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Lead und *Added*-Eintrag zu `devcontainer remove`; die Bedingung „ohne Container Exit `0`“ pinnen `internal/hexagon/application/devcontainer_remove_test.go` (`no container`, `TestDevcontainerRemove_ReleaseTable`) und `internal/adapter/driving/cli/devcontainer_remove_test.go`; `make gates` grün (Coverage 92,40 %). Der Cut fügt keinen Test hinzu |
| [`LH-FA-DEV-013`](../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) | *Added*-Eintrag zu `brew`; die Tests der trägenden Slices liegen in `internal/` und `internal/e2e/brew_docker_test.go`; der Docker-Teil läuft erst als `integration-docker` auf dem Push-Tip (offen) |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md) | *Changed*-Eintrag; `gate-reihenfolge.sh` und `gate-reihenfolge-test.sh` im eigenen `make gates` (Zeilen 3 bis 6 der Ausgabe) |
| [`ADR-0046`](../plan/adr/0046-homebrew-formeln-im-devcontainer-schluessel-brew-und-tap-freigabe.md), [`ADR-0047`](../plan/adr/0047-brew-formeln-lastenheft-stellen-und-verteilung-lastenheft-spezifikation.md) | *Added*-Eintrag zu `brew` (Bezugsklammer); der Cut ändert weder ADR noch Spec |

### Carveouts

- Neu: none
- Gelöst: none
- Unverändert: none berührt (kein Eintrag unter `docs/plan/carveouts/` im Diff)

### Nicht ausgeführt

- **`ci`, `integration-docker`, `tap-check` auf dem Push-Tip** — der Push steht aus (Plan §1, erste Freigabe des Projektinhabers). Es gibt keinen Lauf auf `8aca6bef`; eine Erfolgsmeldung für den Tip wird nicht behauptet.
- **`make test-docker`** — Plan §1 schließt ihn im Devcontainer-Container aus; Ersatz ist `integration-docker` auf dem Push-Tip (oben).
- **`make verify-depguard`** — keine Importregel berührt.
- **Tag, `publish`, Job `tap`, Kontrolle nach `releasing.md` §6** — der Tag steht aus.
- **Aussage über `rg` im Entwicklungs-Devcontainer** — nicht Gegenstand dieses Slice.

### Commit / Artefakt

- Cut-Commit `58411aaf`; Lead-Fix und Review-Report `8aca6bef` (lokaler Tip, nicht gepusht); `origin/main` auf `68321ff6`.
- Lokales Image `u-boot:latest` aus `make build VERSION=0.24.0` (`--version` und Label `0.24.0`), nicht gepusht.
- Kein Tag, kein Release-Asset.

---

## Findings des Verifiers

| ID | Kategorie | Befund | Auflösung |
|---|---|---|---|
| V-1 | MEDIUM | **Die Push-Tip-Belege fehlen.** DoD-Punkt 2 und die zweite Bedingung des Tags verlangen `ci`, `integration-docker` und `tap-check` auf dem Push-Tip; `8aca6bef` ist ungepusht und ohne Lauf. Die Läufe auf `origin/main` (`68321ff6`) gelten einem anderen Commit. Plan §1 und Risiko 9 binden den Tag an `git rev-parse HEAD` gleich der `headSha` des grünen Laufs. Lokal ist nur der Teil belegt, den der Plan als Vorab-Prüfung führt. | Nach der ersten Freigabe des Projektinhabers pushen, dann `ci` und `integration-docker` je mit Conclusion und Log lesen (Risiko 4: die acht Tests auf `--- PASS`, nicht `--- SKIP`; Log über `gh api …/jobs/<id>/logs`), `tap-check` starten und lesen; Stand-Bericht; erst danach die zweite Freigabe für den Tag. Jeder weitere Commit vor dem Tag macht diese Läufe wertlos. |
| V-2 | LOW | **Die Lead-Zeile aus `8aca6bef` hat kein Nachlauf-Review.** Der Review-Report liest `58411aaf`; die Korrektur beider Findings ist ein Text im Cut-Block (1 Zeile). Plan §6 Risiko 1 nennt als Messstelle M-L und M-R „vor dem Push“. Ich habe die Zeile gegen *Added*-Eintrag, Tests und *Changed*-Eintrag gelesen (siehe *Prüfung der Review-Übergaben*): keine Abweichung. Das ist eine Verifier-Lesung, kein Review. | Die Entscheidung, ob ein Nachlauf-Review der Zeile vor dem Push gebraucht wird, liegt beim Planner. Aus dem Diff ergibt sich keine Abweichung. |
| V-3 | INFO | **Das Datum ist heute gelesen, nicht am Tag-Tag.** `date +%F`, CHANGELOG, Handbuch-`Stand` und `git log -1 --format=%cs HEAD` zeigen `2026-11-06`. Läuft der Tag an einem späteren Tag, ist nach `releasing.md` §5 ein weiterer Cut-Commit fällig, und jeder weitere Commit entwertet V-1. | M0 am Tag-Tag wiederholen, vor der zweiten Freigabe. |
| V-4 | INFO | **Der Trockenlauf lief im Image `u-boot:test-docker-tools`, nicht in `u-boot:test`** (`releasing.md` §4 nennt `u-boot:test` „zum Beispiel“). Der Beleg zählt: Skripte, Fixture und Stellvertreter sind die Soll-Form, nur das Basis-Image ist ein anderes. | nur nennen. |
| V-5 | INFO | **Die Zeilennummern im Plan sind überholt** (Compare-Links Zeilen 1662/1663 im Plan, im Cut 1670/1671; die *Added*-Einträge Zeilen 16/17 im Plan, im Cut 20/21). Das ist der Plan-Stand vor dem Cut, kein Befund am Cut. | Bei der Closure die Zeilenangaben entfernen oder auf den Ist-Stand setzen (`AGENTS.md` §3.10, Zustandsmarker mit Ablaufpunkt). |

## Verdikt

**DoD nicht vollständig — der Slice bleibt in `in-progress/`, der Tag wartet.** Erbracht und belegt: Liefer-Punkt 1 (Release-Text, Versionsstellen an drei Stellen samt Kommentaren, READMEs und Handbuch), der lokale Teil von Liefer-Punkt 2 (`make gates`, `make ci`, `make build VERSION=0.24.0` samt Label-Gegenprobe, Tap-Trockenlauf mit Gegenprobe), `releasing.md` §4 aus Liefer-Punkt 3, beide Review-Findings (F-1, F-2) aufgelöst und die Zusage des Cut-Commits (nur die Versionszeile im Produktcode, kein Pipeline-, Spec- oder ADR-Diff) bestätigt.

**Offen, und damit den Tag blockierend:**

1. **Push-Tip-Belege (V-1)** — `ci`, `integration-docker` (acht Docker-Tests auf `--- PASS`), `tap-check` auf dem Commit, den der Tag tragen soll. Setzt die erste, getrennte Freigabe des Projektinhabers voraus.
2. **Zweite Freigabe des Projektinhabers** für den Tag, nach dem Stand-Bericht.
3. **Nach dem Tag:** `publish` und Job `tap`, Datumsabgleich mit `git log -1 --format=%cs v0.24.0`, Roadmap-Wechsel, Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, drei Paarungen.

**Übergabe an den Planner:** Der Cut ist für den Push bereit, soweit er lokal prüfbar ist. Reihenfolge laut Plan §3: erste Freigabe → Push von `main` → `ci` und `integration-docker` auf dem Push-Tip lesen → `tap-check` → Stand-Bericht (mit den Exit-Codes dieses Reports als Vorab-Belege, Plan §6 Risiko 8) → zweite Freigabe → Tag auf genau diesem Tip → `releasing.md` §6 → Closure. Zu entscheiden: V-2 (Nachlauf-Review der Lead-Zeile, ja oder nein) vor dem Push. Keine Erfolgsmeldung für den Tip, keinen Tag und keinen Push habe ich ausgelöst.

## Sensor-Lauf

`make docs-check` nach Anlage dieses Reports: `1199 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf fand eine unverlinkte Kennung in der Tabelle *Plan-vs-Code*; berichtigt).
