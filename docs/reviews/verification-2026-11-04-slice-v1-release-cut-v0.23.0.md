# Verification-Report: slice-v1-release-cut-v0.23.0 — 2026-11-04

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-release-cut-v0.23.0`](../plan/planning/done/slice-v1-release-cut-v0.23.0.md) — Cut-Commit `84a1ff3e` (Tag `v0.23.0`, veröffentlicht), Review-Fix-Runde `83801980` (nur Text im Block `[0.23.0]` auf `main` und Review-Report; der Tag trägt den früheren Wortlaut); Abgleichstand **nach der letzten Fix-Runde = `83801980`** (= Repo-Tip, Arbeitsbaum sauber).

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel), §2 (DoD), §3 (Plan-Tabelle), §6 (sieben Risiken), §7 (Closure offen)
- [`LH-FA-DEV-012`](../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`ADR-0044`](../plan/adr/0044-devcontainer-umgebungsdatei-env-file.md), [`ADR-0045`](../plan/adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)
- Review-Report [`2026-11-04-slice-v1-release-cut-v0.23.0.md`](2026-11-04-slice-v1-release-cut-v0.23.0.md) (0 HIGH, 0 MEDIUM, 1 LOW F-1 im HEAD berichtigt, 1 INFO F-2 an den Verifier übergeben)
- [`harness/verification.md`](../../harness/verification.md), [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 bis §7, [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.10, §5)

---

## Scope

- Slice: [`slice-v1-release-cut-v0.23.0`](../plan/planning/done/slice-v1-release-cut-v0.23.0.md)
- IDs: [`LH-FA-DEV-012`](../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`ADR-0044`](../plan/adr/0044-devcontainer-umgebungsdatei-env-file.md), [`ADR-0045`](../plan/adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)
- Artefakte: `CHANGELOG.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, Roadmap

## DoD-Abgleich

DoD-Quelle: Plan §2. Die Haken stehen offen (Closure steht aus); hier ist der Zustand je Punkt gegen den Beleg gelesen.

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | Release-Text und Versionsstellen | **erbracht** | `CHANGELOG.md`: `## [0.23.0] - 2026-11-04` unter leerem `## [Unreleased]`, Lead „Twenty-third release.“; Compare-Links `[Unreleased]: …compare/v0.23.0...HEAD`, `[0.23.0]: …compare/v0.22.0...v0.23.0`, `[0.22.0]` darunter unverändert. `0.23.0-dev` an `cmd/uboot/main.go` (`var version`), `repo.mk` (`VERSION ?=` und zwei Kommentare), `Dockerfile` (`ARG UBOOT_VERSION` und zwei Kommentare); `grep 0.22.0` in `cmd`, `repo.mk`, `Dockerfile`, READMEs und Handbuch ohne Treffer. READMEs (Status-Zeile, `VERSION=`) und Handbuch (`Software-Version`, `u-boot version 0.23.0`, `download/v0.23.0`, `VERSION=`, Gültigkeitsbereich) auf `0.23.0`. Datum: `git log -1 --format=%cs v0.23.0` = `2026-11-04` = CHANGELOG-Datum = Handbuch-`Stand`. |
| 2 | Vorab-Prüfungen und Tag | **teils belegt** | Belegt: `ci`, `integration-docker`, `tap-check`, `publish` samt Tag und Release (siehe unten). Nicht als Beleg im Repo: `make gates` des Implementers am Cut-Commit, `make ci`, `make build VERSION=0.23.0` mit Gegenprobe, Tap-Trockenlauf mit Vorgänger-Gegenprobe (siehe Abschnitt Vorab-Belege). |
| 3 | Doku nachgezogen | **offen (Closure)** | Roadmap-Zeile unter *Offene Wellen* trägt „Push und Tag stehen aus“; der Tag steht. *Abgeschlossene Wellen* mit Tag-Hash `84a1ff3e` fehlt noch. |
| 4 | `make gates` grün nach dem letzten Commit | **erbracht für HEAD** | Eigener Lauf, siehe Sensors; „nach der Closure“ steht aus. |
| 5 | Review, Verification, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **teils** | Review-Report vorliegend (0 HIGH, 0 MEDIUM); Verification: dieser Report. Closure-Notiz (§7), Register-Sichtung (§8 „nicht durchgegangen“), Risiko-Ausgänge (§6 Risiken 1 bis 5 und 7 „offen → bei Closure“), Paarungen: offen. |

### Lead-Aussage gegen den Code (Prüfpunkt 2)

Lead im HEAD (Tag-Wortlaut mit F-1-Abweichung, siehe unten). Gelesen gegen `git diff v0.22.0..v0.23.0 -- cmd internal go.mod Dockerfile repo.mk` (33 Dateien, davon die Nicht-Test-Dateien: `cmd/uboot/main.go`, `Dockerfile`, `repo.mk` je nur Versionswert und Kommentare, `go.mod` **unberührt**, zwölf Dateien unter `internal/`):

| Aussage | Fundstelle | Urteil |
|---|---|---|
| `generate devcontainer` und `init --devcontainer` schreiben `--env-file=${localWorkspaceFolder}/.devcontainer/devcontainer.env` in `runArgs` | `internal/hexagon/application/devcontainer_sandbox.go` (Diff +29); `init --devcontainer` existiert (`internal/adapter/driving/cli/init.go`, Flag-Doku Z. 92 und Beispiel Z. 144) | trägt |
| legt die Datei „with comment lines only“ nur an, wenn sie fehlt | `templates/devcontainer/devcontainer.env.tmpl` (10 Zeilen), `generate.go` (+32) und `initproject.go` (+34) im Diff; das Anlegen-nur-wenn-fehlend ist durch die zwei Docker-Tests und die Unit-Tests der Träger-Slices belegt, im Cut nicht neu gelesen | trägt (Eintrag im Block, vom Review gelesen) |
| neue `doctor`-Prüfung warnt bei fehlender Datei | `doctor.go` (`checkIDDevcontainerEnvFile = "devcontainer.envfile"`), `doctor_envfile.go` (neu, 37 Zeilen) | trägt |
| „In a fresh clone of a project that does not track the file (the `.gitignore` of `u-boot init` ignores it) … exit code 125, measured with `docker run`“ (berichtigter Wortlaut F-1) | `templates/gitignore.tmpl` Z. 12 `.devcontainer/devcontainer.env`, `basic/.gitignore.tmpl` (+1); Messung E1 „fehlende Datei … Exit-Code `125`“ in [`docs/messungen/slice-v1-devcontainer-env-datei/README.md`](../messungen/slice-v1-devcontainer-env-datei/README.md) Z. 14 | trägt; Bedingung und Messmittel jetzt genannt |
| `docker-init.sh` „of every sandbox project with `nestedRuntime: docker`“ ändert sich | `generate.go` Z. 1074 (`docker-init.sh` nur im NestedDocker-Zweig), `docker-init.sh.tmpl` (+45) | trägt |
| baut die Root-cgroup vor dem nested `dockerd` um, best-effort | `docker-init.sh.tmpl` Z. 46 bis 72 (`cgroup_prepare`, jeder Fehlerpfad `cgroup_prepare_failed` und `return 0`), Aufruf Z. 87 vor `setsid nohup dockerd` Z. 89 | trägt |
| „measured in one container in one run“, Wirkung nach Neuanlage „not measured“ | Slice §6 Risiko 6; Folge-Slice [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../plan/planning/open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md) in `open/` | trägt; der Runner-Lauf (siehe unten) kam nach dem Cut und ist im *Changed*-Eintrag des HEAD ergänzt |
| „Lastenheft 0.35.0 (new requirement for the environment file)“ | `spec/lastenheft.md` Z. 9 (Version 0.35.0) und Versionstabelle Z. 1649 | trägt |

**Berichtigung F-1 an der Stelle:** Der HEAD-Lead und der *Added*-Eintrag tragen „In a fresh clone of a project that does not track the file … (exit code 125, measured with `docker run`)“; der *Changed*-Eintrag trägt zusätzlich „one run on the GitHub-hosted runner `ubuntu-latest` in the CI job `integration-docker`, passed without a skip“. `git diff 84a1ff3e 83801980 -- CHANGELOG.md` ändert genau drei Zeilen (Lead, *Added*, *Changed*), sonst nichts; kein Code, keine Compare-Links, kein Datum. Die Aussage „passed without a skip“ ist **am Lauf gelesen und trägt** (siehe `integration-docker`). Der Tag `v0.23.0` zeigt weiter auf `84a1ff3e` und trägt den früheren Wortlaut; `main` und Tag weichen im Block `[0.23.0]` ab (Ausgang zu Risiko 1 in der Closure-Notiz festhalten, wie bei `v0.22.0`).

`go.mod` ist im Diff `v0.22.0..v0.23.0` unberührt, der Lead sagt dazu nichts. Die Lead-Aussagen tragen keine Zählwerte, die nach dem Cut-Commit abweichen könnten.

## Belege nach Push und Tag (Prüfpunkt 3, lesend abgefragt; kein Token ausgegeben)

| Beleg | Ergebnis |
|---|---|
| `ci` Run 41837587543, push, Head `84a1ff3e` | success; Jobs `gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`, `image-scan (trivy HIGH+CRITICAL)` je success (13:23:04 → 13:25:09) |
| `integration-docker` Run 41837587498, push, Head `84a1ff3e` | success; Job `integration-docker (//go:build docker)` 13:23:04 → 13:40:33 (≈17 min). Log gelesen (12 860 Zeilen), siehe Testtabelle |
| `tap-check` Run 41837628526, workflow_dispatch, Head `84a1ff3e` | success, Job `check` (13:23:23 → 13:23:36); vor dem Tag (Lauf `publish` 14:07:16) — Reihenfolge M4 vor M5 gewahrt |
| `publish` Run 41844401037, push, Head `84a1ff3e` | success; Jobs `publish (ghcr.io image push)`, `package smoke (rpm, fedora)`, `tap (homebrew formula)`, `package smoke (deb, ubuntu)` je success. Log: `Resolved version=0.23.0, is_stable=true`, Build mit `UBOOT_VERSION=0.23.0`, Push `:0.23.0` und `:latest`. Job `tap` **wirklich gelaufen**: Schritt „Pull the formula into the tap“ zeigt `Cloning`, `[main bade812] u-boot 0.23.0`, `52dfd50..bade812  HEAD -> main`, `tap-nachzug: pt9912/homebrew-u-boot updated to v0.23.0.`; die Warnung „HOMEBREW_TAP_GITHUB_TOKEN is not configured“ steht nur im Skript-Echo des Run-Blocks, nicht als Ausgabe. Schließt Risiko 4 sachlich |
| `gh release view v0.23.0` | veröffentlicht (nicht Draft, nicht Prerelease), 2026-11-04T14:09:23Z; 12 Assets: `SHA256SUMS`, `u-boot.rb`, sechs Binaries (linux/darwin/windows × amd64/arm64), zwei `.rpm` (x86_64, aarch64), zwei `.deb` (amd64, arm64) |
| Tag | `v0.23.0` ist ein leichtgewichtiger Tag (`git cat-file -t` = `commit`), `git rev-parse v0.23.0^{commit}` = `84a1ff3e24fc2bd4be40a242bea08fbbcce07ec1` = Cut-Commit, Commit-Zeit 13:20:52Z |
| `docker run ghcr.io/pt9912/u-boot:0.23.0 --version` (§6) | **nicht ausgeführt** (kein Docker-Pull in diesem Lauf); der Build-Schritt des `publish`-Logs zeigt `UBOOT_VERSION=0.23.0` |
| Mac-Probe (M7, `brew install`) | **nicht belegt:** Aussage des Projektinhabers („Auf dem Mac läuft es auch“), kein Protokoll. Nicht geprüft; Tap-Commit `bade812` im Log belegt nur den Nachzug |

### Tests im Log von `integration-docker` (Prüfpunkt 3, Skip-Pfad)

| Test | Ergebnis im Log |
|---|---|
| `TestE2E_EnvFile_VariableReachesContainer` | `--- PASS` (0.43s) |
| `TestE2E_EnvFile_VariableReachesContainerViaDevcontainerUp` | `--- PASS` (1.26s) |
| `TestE2E_SandboxDevcontainer_NestedDocker_CgroupLimitsReachNestedContainer` | `--- PASS` (47.46s), **kein `--- SKIP`** (der `t.Skipf`-Pfad in `internal/e2e/sandbox_devcontainer_docker_test.go` Z. 1230, „`memory` is already in cgroup.subtree_control“, wurde nicht genommen) |
| `TestE2E_Resources_DefaultProfileLimitsReachTheContainer` | `--- PASS` (2.58s) |
| `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` | `--- PASS` (14.47s) |

`--- FAIL`: **0**. `--- SKIP`: 4, alle nicht aus der Prüfmenge: `TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…NestedPodman_BuildEgress`, `…NestedPodman`, `…StatusFiles_VFS`. Damit trägt der Runner-Lauf die Aussage „passed without a skip“ im HEAD-Text (F-2 des Reviews erledigt); er belegt **nicht** die Neuanlage-Messung (Folge-Slice).

## Vorab-Belege ohne Log im Repo (Prüfpunkt 7)

| Vorab-Beleg (DoD 2) | Was der Verifier belegen kann | Zustand |
|---|---|---|
| `make gates` des Implementers am Cut-Commit | nur indirekt: `ci` Run 41837587543 (gates-Job: lint, test, coverage-gate) auf `84a1ff3e` grün; eigener Lauf am HEAD grün | Sitzungsprotokoll; Behauptung ersetzt durch CI und eigenen Lauf |
| `make ci` lokal | `ci` auf dem Runner grün (govulncheck, image-scan); der lokale Lauf liegt nicht vor | nur Sitzungsprotokoll, nicht wiederholt (Auftrag) |
| `make build VERSION=0.23.0` mit Gegenprobe | `publish`-Build mit `UBOOT_VERSION=0.23.0` grün; die lokale Gegenprobe (`--version` und OCI-Label, Standardbau meldet `0.23.0-dev`) liegt nicht vor | nur Sitzungsprotokoll, nicht wiederholt |
| Tap-Trockenlauf (M-T) mit Vorgänger-Gegenprobe | `tap-check` (Schreibrecht des Tokens) und der echte Job `tap` grün; der Trockenlauf mit Fixture und Vorgänger-Tag liegt nicht vor | nur Sitzungsprotokoll; der echte Lauf ersetzt ihn sachlich, nicht als Beleg des Schritts |

In der Closure-Notiz benennen: diese drei lokalen Belege stehen nur im Sitzungsprotokoll, getragen vom CI- und Release-Nachweis.

## Plan-vs-Code (§3 gegen `git diff --stat`)

Gegen den Cut-Commit `84a1ff3e` (8 Dateien, +23/−18):

| Plan-Zeile (§3) | Diff | Abgleich |
|---|---|---|
| `CHANGELOG.md` | +7/−1: Block, Lead, zwei Compare-Links (Zeile `[Unreleased]` ersetzt) | deckt |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | 1 + 3 + 3 Zeilen, je Wert samt Kommentar | deckt |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | je 2/2/5 Zeilen (4/4/10 Diffzeilen) | deckt |
| `docs/plan/planning/in-progress/roadmap.md` | 1 Zeile (Offene Wellen) | deckt; *Abgeschlossene Wellen* folgt mit der Closure |
| `docs/plan/planning/open/README.md` | „mit dem Plan-Commit“ | stimmt (nicht im Cut-Commit; Bestand-Zeile Z. 20 steht) |
| `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | „kein Diff durch den Cut“ | stimmt (keine dieser Dateien im Cut-Commit; `.github`, `scripts`, `packaging`, `Makefile`, `go.mod` auch im Range `v0.22.0..HEAD` ohne Diff; `spec/lastenheft.md` und `spec/spezifikation.md` tragen den Diff der Träger-Slices) |

Diff-Dateien des Cut-Commits ohne Plan-Zeile: keine. Plan-Zeilen ohne Diff: nur die als „kein Diff“ geführten.

Nach dem Tag (`84a1ff3e..HEAD`): `CHANGELOG.md` (drei Textzeilen, Plan-Zeile „update“ trägt es) und der Review-Report [`2026-11-04-slice-v1-release-cut-v0.23.0.md`](2026-11-04-slice-v1-release-cut-v0.23.0.md) (neu). **Die Plan-Tabelle hat für den Review-Report keine Zeile** (und diesen Verification-Report ebenso); beim Release-Cut `v0.22.0` hat die Tabelle die Fix-Runde übernommen. Abgleich nach der letzten Fix-Runde nach `AGENTS.md` §3.10: eine Zeile „`docs/reviews/` (Review- und Verification-Report)“ nachtragen (Hinweis, siehe Rest-Risiken). Range `v0.22.0..HEAD`: 83 Dateien, +3185/−497; die über den Cut hinausgehenden Dateien sind die der Träger-Slices (Produktcode, Templates, Tests, Spec, [`ADR-0044`](../plan/adr/0044-devcontainer-umgebungsdatei-env-file.md) und [`ADR-0045`](../plan/adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md), Register, Reports), von der letzten Plan-Zeile und §1 abgedeckt.

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer, eigener Lauf am 2026-11-04, HEAD `83801980`, vor diesem Report) | **pass** | EXIT=0; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`; `d-check: 1125 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — coverage 92.10% meets threshold 90%`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `tool-pruefung … OK (go 1.27.3, golangci-lint 2.14.0)`. Stempel `.harness/state/gates-passed.head` = `83801980a1ee1b617fe8fdbead4a72f08923a8a1` = HEAD. |
| `make docs-check` nach diesem Report | im Handback des Verifiers gemeldet | siehe Handback |
| `make test-docker` | **not run, begründet** | Auftrag (Docker-Last); Ersatznachweis ist `integration-docker` Run 41837587498 auf dem Push-Tip (grün, Tests oben). Lokal in dieser Sitzung bis zur Neuanlage rot (Plan §6 Risiko 3). |
| `make verify-depguard` | **not run, begründet** | Keine Importregel berührt; der Cut ändert in `cmd/` und `internal/` allein die Versionszeile. |
| `make ci`, `make fullbuild`, `make build VERSION=0.23.0`, Tap-Trockenlauf | **not run** | Auftrag; Ersatz teils durch `ci` und `publish` am selben Tip (siehe Vorab-Belege). |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-012`](../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | Cut-Commit `84a1ff3e` und Fix-Commit `83801980` nennen beide IDs und den Slice [`slice-v1-release-cut-v0.23.0`](../plan/planning/done/slice-v1-release-cut-v0.23.0.md) im Betreff |
| Commits `v0.22.0..HEAD` | 53 Commits; alle nennen eine `LH-*`-, `ADR-*`-, Slice- oder Wellen-Kennung (geprüft per Muster auf den Betreff; kein Commit ohne) |
| Slice-Moves | 18 Verschiebe-Commits des Make-Targets für Slice-Moves im Range; für diesen Slice `979417a9` (reiner Move `next/` → `in-progress/`) und `911cc938` (Verweis-Nachzug als eigener Commit, §3.3) |

Docs-Links: `make docs-check` im `gates`-Lauf grün (0 Befunde; nach diesem Report erneut gefahren).

## Carveouts

- Neu: **keine**; keine Lockerung von Gates, Schwellen oder Architekturregeln (kein Pin-Bump, keine Pipeline-Änderung im Cut; `.github`, `scripts`, `Makefile`, `go.mod` im Range ohne Diff).

## Rest-Risiken / offene Punkte vor der Closure

1. **Veröffentlichter Lead weicht ab.** Der Tag `v0.23.0` (`84a1ff3e:CHANGELOG.md`) trägt den früheren Lead ohne die Bedingung „does not track the file“ und ohne das Messmittel, `*Changed*` ohne den Runner-Lauf; `main` ist berichtigt. Der Tag bleibt. Die Closure-Notiz hält die Abweichung Tag-Stand gegen `main` als Ausgang zu Risiko 1 fest.
2. **Mac-Probe ohne Protokoll** (M7): nur Aussage des Projektinhabers; in der Closure als „ohne Protokoll“ führen, nicht als belegt.
3. **Lokale Vorab-Belege** (`make gates` am Cut-Commit, `make ci`, `make build VERSION=0.23.0` mit Gegenprobe, Tap-Trockenlauf) stehen nur im Sitzungsprotokoll; CI und Release tragen den Nachweis, nicht den Schritt selbst.
4. **Plan-Tabelle §3 ohne Zeile für `docs/reviews/`** (Review- und Verification-Report) nach der Fix-Runde; vor der Closure nachtragen, danach `git diff --stat` erneut abgleichen (`AGENTS.md` §3.10). Hinweis, kein Finding der Cut-Diffs.
5. **Closure-Pflichten offen** (kein Finding): §7 leer; DoD-Haken offen; Roadmap-Marker „Push und Tag stehen aus“ abgelaufen → auf Ist-Stand setzen und *Abgeschlossene Wellen* mit Tag-Hash `84a1ff3e` ergänzen; §6 Risiken 1 bis 5 und 7 ohne Ausgang (Belege: Datum = Tag-Datum; `integration-docker` grün ohne Skip der Prüfmenge; `tap` grün und gelaufen; Ausgabe-Änderung steht im CHANGELOG); §8 Beobachtungs-Sichtung „nicht durchgegangen“; Verschiebung nach `done/` erst nach der Notiz (§3.3 Fall 2); `make gates` nach der Closure.
6. **Risiko 6** (cgroup-Verhalten nach Neuanlage ungemessen): bleibt beim Folge-Slice in `open/`; der Runner-Lauf belegt es nicht.
7. **§6 der Release-Anleitung:** `docker run ghcr.io/pt9912/u-boot:0.23.0 --version` und der `brew`-Teil sind nicht protokolliert.

## Gesamturteil

**Closure-fähig unter Bedingungen.** Der Cut-Diff deckt Plan §2 Punkt 1 und die Plan-Tabelle; die Lead-Aussagen tragen gegen `v0.22.0..v0.23.0` und die zwei Einträge, die Berichtigung F-1 ist an der Stelle gelesen (drei Textzeilen, kein Code); alle vier CI-/Release-Läufe am Tip `84a1ff3e` sind grün, die fünf geforderten Tests stehen auf `--- PASS` ohne `--- SKIP`, es gibt kein `--- FAIL`; der Job `tap` hat nachgezogen (`bade812`); das Release trägt die erwarteten zwölf Assets, der Tag steht auf dem Cut-Commit; `make gates` ist am HEAD `83801980` grün. Vor der Closure: Closure-Notiz mit Ausgängen zu den Rest-Risiken 1 bis 3, Plan-Tabelle um `docs/reviews/` ergänzen, Roadmap auf den Ist-Stand, `make gates` danach.

## Commit / Artefakt

- Slice: `84a1ff3e` (Cut, Tag `v0.23.0`), `83801980` (Review-Fix-Runde)
- HEAD/Stempel: `83801980a1ee1b617fe8fdbead4a72f08923a8a1`
- Lokal: `make gates` EXIT=0 am 2026-11-04
