# Verification-Report: slice-v1-release-cut-v0.17.0

* Rolle: Verifier (frischer Kontext), Datum 2026-11-01
* Gegenstand: Release-Cut `v0.17.0` im Hauptbaum, `HEAD` = Tag `v0.17.0` = `f89424e1` (veröffentlicht 2026-11-01), Slice [`slice-v1-release-cut-v0.17.0`](../plan/planning/done/slice-v1-release-cut-v0.17.0.md), Muster [`slice-v1-release-cut-v0.16.1`](../plan/planning/done/slice-v1-release-cut-v0.16.1.md), Ablauf [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 bis §7
* Frage: Wurde der Cut richtig gebaut und das Veröffentlichte richtig kontrolliert, gemessen an Plan, DoD, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md) und eigenen Sensoren? Maßstab: [`harness/verification.md`](../../harness/verification.md).
* Nur lesende Befehle gegen GitHub, GHCR und das Tap; nichts geschrieben, nichts getaggt, nichts gepusht. Einzige Änderung: dieser Report.

**Skill/Modell:** Verifier, Claude Sonnet 5.5. **Review-Art:** Code/Release (Plan-vs-Ist, DoD, Veröffentlichtes).

## Ergebnis

Der Release ist **konform veröffentlicht**: Tag, Datum, 12 Assets, `publish` inkl. Job `tap` und Paket-Smokes, GHCR-Image, Formel im Tap und Digests stimmen mit der Zusage überein. DoD-Punkt 1 ist erfüllt. Die Punkte 2 bis 9 sind zum Prüfzeitpunkt **erwartungsgemäß offen oder teilerfüllt**: Der Slice liegt noch in `in-progress/`, §7 ist „Folgt nach dem Tag“, die Nacharbeit (§7 in `releasing.md`) ist nicht vollzogen. Das ist keine Abweichung, sondern der Stand vor dem Nachzug. 0 Befunde mit Blockwirkung auf den Release selbst.

## 1. Veröffentlichtes, selbst geprüft

| Prüfung | Soll | Ist |
|---|---|---|
| `git rev-parse v0.17.0` | `f89424e1` | `f89424e1724647a9398536562b612ef63a1a18e8` (= `HEAD`; leichtgewichtiger Tag, `git cat-file -t` = `commit`) |
| `git log -1 --format=%cs v0.17.0` | = CHANGELOG-Datum = Handbuch-`Stand` | `2026-11-01`; `CHANGELOG.md:14` `## [0.17.0] - 2026-11-01`; `docs/user/benutzerhandbuch.md:5` `Stand: 2026-11-01` |
| `gh release view v0.17.0` | 12 Assets, kein Prerelease | 12 Assets (`SHA256SUMS`, 6 Binaries, 2 `.deb`, 2 `.rpm`, `u-boot.rb`); `isPrerelease=false`, `isDraft=false`, `targetCommitish=f89424e1`, publiziert 2026-11-01T09:21:23Z |
| `gh run view 41292635297` | `publish` success inkl. `tap`, Paket-Smokes | `workflowName=publish`, `event=push`, `headSha=f89424e1`, Conclusion success; Jobs: `publish (ghcr.io image push)`, `tap (homebrew formula)`, `package smoke (deb, ubuntu)`, `package smoke (rpm, fedora)` alle success |
| Log Job `tap` (`gh api …/jobs/124807340208/logs`) | `updated to v0.17.0.`, keine Secret-Warnung | Zeile `tap-nachzug: pt9912/homebrew-u-boot updated to v0.17.0.`; `[main cb47bc3] u-boot 0.17.0`, Push `78bb2ec..cb47bc3`. Der Text `HOMEBREW_TAP_GITHUB_TOKEN is not configured` steht im Log nur als Quelltext des `run`-Blocks (Gruppenkopf „Run if [ -z "$TAP_TOKEN" ]“), nicht als ausgegebene Warnung; `TAP_TOKEN: ***` war gesetzt |
| `docker run --rm ghcr.io/pt9912/u-boot:0.17.0 --version` und `:latest` | `0.17.0`, gleiche Image-ID | beide `u-boot version 0.17.0`; beide `sha256:a18ffc55bf3c1a5a6fe21438e7eb94011501fa154904668f3cdd31bea7edfbaa` |
| Formel im Tap (`Formula/u-boot.rb`) | vier URLs auf `v0.17.0` | vier `url`-Zeilen (`darwin-arm64`, `darwin-amd64`, `linux-arm64`, `linux-amd64`), alle `…/download/v0.17.0/…`; byte-gleich (`cmp`) mit dem Release-Asset `u-boot.rb` |
| `SHA256SUMS` gegen `u-boot.rb` | vier Digests gleich | `darwin-arm64 9e19f116…`, `darwin-amd64 0e6596d1…`, `linux-arm64 c22921fc…`, `linux-amd64 70d0127b…` in `SHA256SUMS` und in der Formel identisch |
| `sha256sum -c` (Download in leeres Scratchpad-Verzeichnis `…/scratchpad/dl`, nicht in den Baum) | Assets passen zu `SHA256SUMS` | 6 Binaries und 2 `.rpm`: OK. Die zwei `.deb` habe ich nicht geladen (mein Filter `u-boot-*` traf `u-boot_*` nicht); ihre Digests stehen in `SHA256SUMS`, die Datei selbst ist ungeprüft. Der Smoke-Job `deb, ubuntu` ist success |
| CI auf `f89424e1` | `ci`, `integration-docker` grün; `tap-check` grün | `ci` 41290870228 success, `integration-docker` 41290870231 success, `tap-check` 41290880081 success (`gh run list --branch main`) |

**Hinweis zu M4 (Reihenfolge `tap-check` vor Tag):** Nach Lauf-ID (41290880081 < 41292635297) lag `tap-check` vor `publish`. Die Zeitstempel der API sind dagegen nicht widerspruchsfrei (`createdAt` von `tap-check` 09:55:45 gegenüber `publish` 09:19:33; beim `ci`-Lauf liegt `updatedAt` vor `createdAt`). Ich stütze die Reihenfolge nur auf die IDs und werte die Zeitstempel nicht als Beleg. Der Job `tap` hat das Token im echten Lauf ohnehin bewiesen.

## 2. DoD, Punkt für Punkt (Slice §2)

| # | DoD-Punkt | Befund | Beleg |
|---|---|---|---|
| 1 | Release-Text und Versionsstellen | **erfüllt** | `CHANGELOG.md:12` `## [Unreleased]` (leer), `:14` `## [0.17.0] - 2026-11-01` mit Lead; Fuß `:1573` `[Unreleased]: …/compare/v0.17.0...HEAD`, `:1574` `[0.17.0]: …/compare/v0.16.1...v0.17.0`. `cmd/uboot/main.go:43` `0.17.0-dev`, `repo.mk:57` (Kommentar) und `:61` `VERSION`-Bereich ohne `0.16.1`, `Dockerfile:41/:43/:45` `0.17.0-dev`. `README.md:10`/`:98`, `README.de.md:11`/`:101` `v0.17.0`/`0.17.0`. Handbuch `:4`, `:84`, `:110`, `:136`, `:3038` `0.17.0`; kein `0.16.1` mehr im Handbuch. Datum in CHANGELOG und Handbuch = Tag-Datum (Abschnitt 1) |
| 2 | Tag und Kontrolle (M5, M6) | **Tag und Kontrolle sachlich erfüllt, Evidence-Eintrag offen** | Tag und `publish` siehe Abschnitt 1. `make gates`/`ci`/`fullbuild`/Trockenlauf stehen in §3 des Slice mit Exit-Codes (vom Implementer berichtet); `make gates` habe ich selbst gefahren (Abschnitt 3). §7 des Slice („Folgt nach dem Tag“) trägt die Evidence noch nicht: das ist der Nacharbeit-Schritt, nicht der Cut |
| 3 | Doku nachgezogen (Roadmap, `releasing.md`, `open/README.md`) | **offen** | `roadmap.md:20` führt `v0.17.0` noch unter *Offene Wellen* mit „(Push und Tag stehen aus)“ — dieser Marker ist seit dem Tag überholt (§3.10 Ablaufpunkt); keine Zeile in *Abgeschlossene Wellen*; `releasing.md:56` Zeiger „zuletzt“ nennt weiter `v0.16.1` (plan-konform bis zur Closure) |
| 4 | `make gates` grün | **erfüllt für den Ist-Stand** | Abschnitt 3: Exit 0, selbst gefahren |
| 5 | Review durchgeführt, Report unter `docs/reviews/` | **offen** | Beim Prüfen lag kein `…-release-cut-v0.17.0.md`-Review-Report im Ordner; ein Reviewer legt ihn parallel an. Nicht von mir bewertet |
| 6 | Closure-Notiz mit Steering-Loop-Eintrag | **offen** | Slice §7 „Folgt nach dem Tag“ |
| 7 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ | **offen** | steht in der Closure (§7) aus |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **offen** | alle sechs Risiken in §6 stehen auf „offen bis zur Closure“. Sachlicher Stand nach meiner Prüfung: Datum ≠ Tag-Datum: nicht eingetreten; `image-scan`-Advisory: `ci` auf dem Tip grün; Job `tap`: grün und gelaufen; `integration-docker`: auf dem Tip grün; Lead-Breite und fehlende Nutzer-Änderung: siehe Abschnitt 4 |
| 9 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen** | an die Closure gebunden |

## 3. Sensoren, selbst gefahren (Docker-Modus, Host `linux/amd64`)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (Docker-Modus, `BUILD_MODE=docker`) | **Exit 0** | `lint` 0 issues; `coverage-gate` 92,20 % gegen 90 %; `docs-check` 980 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung` OK, `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge. Der Lauf begann vor Anlage dieses Reports; ein transienter Befund durch den parallel angelegten Review-Report ist nicht aufgetreten |
| `make ci`, `make fullbuild`, `make build VERSION=0.17.0`, `doc-*` | **nicht von mir gefahren** | die Veröffentlichung belegt dasselbe unabhängig: `ci` grün auf `f89424e1` in der CI, `publish` baute das Image mit `--version` = Tag (GHCR-Prüfung, Abschnitt 1). Die Exit-0-Zeilen des Implementers in Slice §3 sind übernommen, nicht reproduziert |
| `make test-docker` | **nicht gelaufen** | Entscheidung des Projektinhabers (2026-11-01); CI `integration-docker` auf `f89424e1` success (Lauf 41290870231) |
| `make verify-depguard` | **nicht gelaufen** | der Cut berührt keine Importregel (Diff: Versionszeilen und Doku) |
| Trockenlauf des Tap-Pfads (M-T) | **nicht nachgefahren** | Fixture-Läufe sind durch den echten Lauf (Job `tap` grün, `updated to v0.17.0.`) überholt; die Zeilen (2b)/(3b) des Slice sind die Gegenprobe und bleiben Aussage des Implementers |
| `brew tap`/`brew install`, `u-boot --version` auf einem Mac (M7) | **nicht gelaufen** | kein Mac; DoD-fern, im Plan als Aktion des Projektinhabers (M7) geführt |

## 4. Lead gegen Lastenheft

Lastenheft `0.26.0` (`spec/lastenheft.md:9`). Zusage in `spec/lastenheft.md:504`: Workspace-Volumes, zwei gleichzeitig laufende Instanzen, Sandbox-Profil, direkt nach dem Start, Linux `amd64`, lokale Docker-Engine. Der Lead (`CHANGELOG.md:18`) nennt dieselben Bedingungen und eine Nicht-Zugesagt-Liste (drei und mehr Instanzen, Docker-Speicher, `ollama`-Modelle, Podman, CPU/Speicher-Grenzen, Verhalten unter Last, andere Hosts und Engines, Standardprofil). Gedeckt.

* **INFO (V-1):** Die Nicht-Zugesagt-Liste im Lead lässt „kalter Cache/kalter Bau“ und „Dauer und Platzbedarf des Starts“ aus `lastenheft.md:504` weg. Das ist die engere Lesart (nichts Zugesagtes fehlt) und kein Verstoß gegen [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md). Kein Handlungsbedarf, wenn der Reviewer es nicht anders wertet.
* Kennungen: `CHANGELOG.md` trägt `LH-`/`ADR-`-Links (zulässig). `README.md`, `README.de.md` und Handbuch: die Versionszeilen tragen keine Kennung (gelesen an `README.md:10`, `README.de.md:11`, Handbuch `:4`, `:84`).

## 5. Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | LOW | Die Roadmap führt `v0.17.0` nach dem Tag noch unter *Offene Wellen* mit dem Marker „(Push und Tag stehen aus)“; die Zeile in *Abgeschlossene Wellen* fehlt. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt), `docs/maintainer/releasing.md` §7 | `docs/plan/planning/in-progress/roadmap.md:20` · „(Push und Tag stehen aus)“ | ja, `grep -n 'Push und Tag stehen aus' roadmap.md` | Zustandsmarker nach Ablaufpunkt nicht aktualisiert (erwartet vor der Nacharbeit) |
| V-2 | INFO | Die API-Zeitstempel der Läufe sind in sich widersprüchlich; die Reihenfolge `tap-check` vor `publish` ist nur über die Lauf-IDs belegt. | `docs/maintainer/releasing.md` §4 (Token vor dem stabilen Tag prüfen) | `gh run view 41290880081 --json createdAt` | ja, `gh run view` | Beleg-Zeitstempel der Plattform uneindeutig |
| V-3 | INFO | Zwei `.deb`-Dateien sind nicht heruntergeladen und nicht gegen `SHA256SUMS` geprüft (Filterfehler des Verifiers); der Paket-Smoke `deb, ubuntu` ist success. | `docs/maintainer/releasing.md` §1 (Paket-Smoke) | `u-boot_0.17.0_amd64.deb`, `u-boot_0.17.0_arm64.deb` | ja, `gh release download -p '*.deb'` und `sha256sum -c` | Prüfumfang des Verifiers |

Keine Findings der Klassen HIGH oder MEDIUM.

## 6. Plan-vs-Code-Diff

`git diff --stat f58dc581..HEAD` (Cut-Commit bis Tag) enthält die Plan-Tabelle-Dateien (Slice §3, Abgleich dort): `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, `open/README.md`, der Slice. Der letzte Commit `f89424e1` ändert nur den Slice (`git show --stat HEAD`: eine Datei, 40 Einfügungen). Kein Diff, gelesen: `docs/maintainer/releasing.md:56` (Zeiger „zuletzt“ `v0.16.1`), `.github/workflows/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `spec/`, `docs/plan/adr/`. Plan-Aussage „Produktcode gegenüber `v0.16.1` geändert“ (§1, 70 Dateien) habe ich nicht erneut gezählt; der Lead behauptet entsprechend kein „unverändert“.

## 7. Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: [`slice-v1-release-cut-v0.17.0`](../plan/planning/done/slice-v1-release-cut-v0.17.0.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md)
- Artefakte: Tag `v0.17.0` (`f89424e1`), GitHub-Release mit 12 Assets, `ghcr.io/pt9912/u-boot:0.17.0` und `:latest` (`sha256:a18ffc55…`), Tap-Commit `cb47bc3`

DoD-Abgleich:
- [x] Release-Text und Versionsstellen — Evidence: Abschnitt 2 Punkt 1
- [ ] Tag und Kontrolle — Tag/`publish`/`tap`/GHCR/Formel/Digests erfüllt (Abschnitt 1); Evidence-Eintrag §7 des Slice steht aus
- [ ] Doku nachgezogen — offen: Nacharbeit nach `releasing.md` §7 (V-1)
- [x] `make gates` grün — Evidence: Exit 0, Abschnitt 3 (nach der Closure zu wiederholen)
- [ ] Review-Report — offen: parallel durch den Reviewer
- [ ] Closure-Notiz, Register, Risiko-Ausgänge, Paarungen — offen: an die Closure gebunden

Sensors: siehe Abschnitt 3.

Traceability: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) — Lead gegen `lastenheft.md:504` gelesen (Abschnitt 4); [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) — Job `tap` grün, Formel im Tap = Asset; [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md) — vier Pakete im Release, Paket-Smokes grün; [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) — Image `0.17.0` = `:latest`; [`ADR-0004`](../plan/adr/0004-ci-system.md) — `ci` und `publish` grün auf dem Tag-Commit.

Carveouts: Neu none · Gelöst none · Unverändert: nicht berührt (kein Diff an `carveouts.md`).

Nicht ausgeführt: `make test-docker` (Projektinhaber-Entscheidung, CI `integration-docker` grün), `make verify-depguard` (keine Importregel berührt), `make ci`/`fullbuild`/`build` (durch CI und `publish` belegt, nicht reproduziert), Trockenlauf (durch echten Lauf überholt), `brew tap`/`brew install` (kein Mac, nicht gelaufen), `.deb`-Digests (V-3).

Commit / Artefakt: `f89424e1`, Tag `v0.17.0`, Run `41292635297`, Tap-Commit `cb47bc3`.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Tag-Hash, Tag-Datum gegen CHANGELOG und Handbuch | geprüft, ohne Befund |
| Release-Assets (Anzahl, Namen, Prerelease-Flag) | geprüft, ohne Befund |
| `publish`-Lauf (alle vier Jobs) und Log des Jobs `tap` | geprüft, ohne Befund |
| GHCR `:0.17.0` und `:latest` (Version, Image-ID) | geprüft, ohne Befund |
| Formel im Tap gegen Asset und `SHA256SUMS` | geprüft, ohne Befund |
| Versionsstellen in Code, Doku, Compare-Links | geprüft, ohne Befund |
| Lead gegen Lastenheft `0.26.0` Zeile 504 | geprüft, ohne Befund (V-1 als INFO) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zustandsmarker nach Ablaufpunkt nicht aktualisiert · Beleg-Zeitstempel der Plattform uneindeutig · Prüfumfang des Verifiers

## Verdikt

**Merge-blockierend:** nein — der Release ist konform; die offenen DoD-Punkte sind die Nacharbeit (`releasing.md` §7) und die Closure, keine Abweichung. Vor `done/`: Roadmap-Marker (V-1) setzen, Zeiger „zuletzt“ in `releasing.md`, §7 mit Hash `f89424e1`, Läufen `41292635297` (`publish`), `41290870228`/`41290870231` (`ci`/`integration-docker`), `41290880081` (`tap-check`), den Review-Report verlinken, danach `make gates` wiederholen. M7 (Mac) bleibt Aktion des Projektinhabers.

**Übergabe:** an den Planner/Implementer für die Nacharbeit; DoD- und Entscheidungs-Konformität: hergestellt für den Veröffentlichungsteil. Der Report ersetzt weder Review noch Validierung.
