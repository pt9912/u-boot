# Verification-Report: slice-v1-release-cut-v0.22.0 — 2026-11-04

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`)
**Gegenstand:** [`slice-v1-release-cut-v0.22.0`](../plan/planning/done/slice-v1-release-cut-v0.22.0.md) — Cut-Commit `f1711f3a` (Tag `v0.22.0`, veröffentlicht), Review-Fix-Runde `353bef50`; Abgleichstand **nach der letzten Fix-Runde = `353bef50`** (= Repo-Tip, Arbeitsbaum sauber).

**Eingangs-Kontext:**

- Slice-Plan §1 (Ziel), §2 (DoD), §3 (Plan-Tabelle), §6 (sechs Risiken), §7 (Closure offen)
- [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Review-Report [`2026-11-04-slice-v1-release-cut-v0.22.0.md`](2026-11-04-slice-v1-release-cut-v0.22.0.md) (0 HIGH, 1 MEDIUM F-1, im HEAD berichtigt)
- [`harness/verification.md`](../../harness/verification.md), [`docs/maintainer/releasing.md`](../maintainer/releasing.md) §3 bis §7, [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.10, §5)

---

## Scope

- Slice: [`slice-v1-release-cut-v0.22.0`](../plan/planning/done/slice-v1-release-cut-v0.22.0.md)
- IDs: [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../plan/adr/0043-nur-devcontainer-bau-pfad.md)
- Artefakte: `CHANGELOG.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, Roadmap

## DoD-Abgleich

DoD-Quelle: Plan §2. Die Haken stehen offen (Closure steht aus); hier ist der Zustand je Punkt gegen den Beleg gelesen.

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 1 | Release-Text und Versionsstellen | **erbracht** | `CHANGELOG.md`: `## [0.22.0] - 2026-11-04` unter leerem `## [Unreleased]`, Lead „Twenty-second release.“; Compare-Links `[Unreleased]: …compare/v0.22.0...HEAD`, `[0.22.0]: …compare/v0.21.0...v0.22.0`. `0.22.0-dev` an `cmd/uboot/main.go` (`var version`), `repo.mk` (`VERSION ?=` und beide Kommentare), `Dockerfile` (`ARG UBOOT_VERSION` und beide Kommentare). READMEs (Status-Zeile, `VERSION=`-Beispiel) und Handbuch (Kopf `Software-Version`/`Stand`, `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich) auf `0.22.0`. Datum: `git log -1 --format=%cs v0.22.0` = `2026-11-04` = CHANGELOG-Datum = Handbuch-`Stand`. |
| 2 | Vorab-Prüfungen und Tag | **teils belegt** | CI und Release-Läufe siehe unten (alle grün). `make gates` des Implementers am Cut-Commit: im Eingang nur behauptet, nicht als Log vorliegend; eigener Lauf am HEAD grün (Sensors). `make ci`, `make fullbuild`, `make build VERSION=0.22.0` und der Tap-Trockenlauf (M-T) liegen dem Verifier nicht als Belege vor und wurden nicht wiederholt (siehe Sensors). |
| 3 | Doku nachgezogen | **offen (Closure)** | Roadmap-Zeile unter *Offene Wellen* trägt „Push und Tag stehen aus“; der Tag steht. *Abgeschlossene Wellen* mit Tag-Hash fehlt noch. |
| 4 | `make gates` grün nach dem letzten Commit | **erbracht für HEAD** | Eigener Lauf, siehe Sensors; „nach der Closure“ steht aus. |
| 5 | Review, Verification, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **teils** | Review-Report vorliegend (0 HIGH); Verification: dieser Report. Closure-Notiz (§7), Register-Sichtung (§8 „nicht durchgegangen“), Risiko-Ausgänge (§6 alle „offen“), Paarungen: offen. |

### Lead-Aussage gegen den Code (Prüfpunkt 2)

Lead im HEAD: „the only non-test file under `cmd/` or `internal/` that differs from 0.21.0 is `cmd/uboot/main.go`, and only in `var version`; in `go.mod` the `go` directive is raised to `1.27.0` (the toolchain pin stays `1.27.3`), and `cobra` and `pflag` move from indirect to direct requirements at the same versions.“

- `git diff --stat v0.21.0..v0.22.0 -- cmd internal go.mod`: `cmd/uboot/main.go` (1 Zeile), `go.mod` (12 Zeilen) und fünf `*_test.go`. Nicht-Test-Dateien: nur `main.go`, die Änderung ist allein `var version = "0.21.0-dev"` → `"0.22.0-dev"`. **Aussage trägt.**
- `go.mod`: `go 1.26.0` → `go 1.27.0`; `cobra v1.10.2` und `pflag v1.0.9` ohne `// indirect`, gleiche Versionen; `mousetrap` bleibt indirect. **Aussage trägt.**
- `GO_VERSION ?= 1.27.3` in `repo.mk` (Z. 59). **Aussage trägt.**
- `git diff f1711f3a..HEAD --stat`: CHANGELOG (1 Zeile = Lead), Slice (Plan-Tabelle), Review-Report. Kein Code nach dem Tag.

## Belege nach Push und Tag (Prüfpunkt 3, lesend abgefragt)

| Beleg | Ergebnis |
|---|---|
| `ci` Run 41816259305, push, Head `f1711f3a` | success; Jobs `gates (lint + test + coverage-gate)`, `security-gates (govulncheck)`, `image-scan (trivy HIGH+CRITICAL)` je success (2m02s) |
| `integration-docker` Run 41816259276, push, Head `f1711f3a` | success; Job `integration-docker (//go:build docker)` 10:45:16 → 11:04:25 (≈19 min). Das ist der Nachweis für `make test-docker` (Plan §1) und schließt Risiko 3 sachlich |
| `tap-check` Run 41818441675, workflow_dispatch, Head `f1711f3a` | success, Job `check` (12 s); 11:04:57, vor dem `publish`-Lauf (11:06:30) — Reihenfolge M4 vor M5 gewahrt |
| `publish` Run 41818642441, push, Head `f1711f3a` | success; Jobs `publish (ghcr.io image push)`, `tap (homebrew formula)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)` je success; kein Schritt mit anderem Ergebnis als success. Schließt Risiko 4 sachlich |
| `gh release view v0.22.0` | veröffentlicht (nicht Draft, nicht Prerelease), 2026-11-04T11:08:54Z; 12 Assets: `SHA256SUMS`, `u-boot.rb`, sechs Binaries (linux/darwin/windows × amd64/arm64), zwei `.rpm` (x86_64, aarch64), zwei `.deb` (amd64, arm64) |
| Tag | `v0.22.0` ist ein leichtgewichtiger Tag, `git rev-parse v0.22.0^{}` = `f1711f3a272ac2854332b7f3228ce85c2b22d67b` = Cut-Commit |
| Mac-Probe (M7) | **nicht belegt:** Aussage des Projektinhabers („läuft auch auf dem Mac“), kein Protokoll. Nicht geprüft |

## Plan-vs-Code (§3 gegen `git diff --stat`)

Gegen den Cut-Commit `f1711f3a` (8 Dateien):

| Plan-Zeile (§3) | Diff | Abgleich |
|---|---|---|
| `CHANGELOG.md` | +7/−1: Block, Lead, zwei Compare-Links | deckt |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | je Wert und Kommentare (1 + 3 + 3 Zeilen) | deckt |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | je 4/4/12 Zeilen | deckt |
| `docs/plan/planning/in-progress/roadmap.md` | 1 Zeile (Offene Wellen) | deckt; *Abgeschlossene Wellen* folgt mit der Closure |
| `docs/plan/planning/open/README.md` | „kein Diff durch den Cut“ | stimmt (nicht im Cut-Commit) |
| `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | „kein Diff durch den Cut“ | stimmt (keine dieser Dateien im Cut-Commit; der Diff gegen `v0.21.0` ist der der Träger-Slices) |

Diff-Dateien ohne Plan-Zeile im Cut-Commit: keine. Plan-Zeilen ohne Diff: nur die bewusst als „kein Diff“ geführten. Gegen `git diff --stat v0.21.0..v0.22.0` (107 Dateien) tragen die übrigen den Diff der zwei Träger-Slices, wie die letzte Plan-Zeile sagt. Nach dem Tag (`f1711f3a..HEAD`): CHANGELOG-Lead, Slice-Plan, Review-Report — die Plan-Tabelle hat die Fix-Runde übernommen (Review F-3).

## Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` (Devcontainer, eigener Lauf am 2026-11-04, HEAD `353bef50`) | **pass** | EXIT=0; `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: pinsync: 3 Pin(s) geprüft, 11 Quelle(n) gelesen`; `d-check: 1104 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — coverage 92.20% meets threshold 90%`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`; `tool-pruefung … OK (go 1.27.3, golangci-lint 2.14.0, make >= 4.3)`. Stempel `.harness/state/gates-passed.head` = `353bef5039433efc0f160ebd55883a9c6ee26354` = HEAD. |
| `make test-docker` | **not run, begründet** | Lokal rot an der cgroup ([`slice-v1-sandbox-nested-docker-cgroup`](../plan/planning/open/slice-v1-sandbox-nested-docker-cgroup.md)); Ersatznachweis ist `integration-docker` Run 41816259276 auf dem Push-Tip (grün). |
| `make verify-depguard` | **not run, begründet** | Keine Importregel berührt; der Code-Diff des Cuts ist die Versionszeile. |
| `make ci`, `make fullbuild`, `make build VERSION=0.22.0`, Tap-Trockenlauf (M-T) | **not run** | Nicht angefordert. Ersatz teils durch CI: `ci` (govulncheck, image-scan) und `publish` (Image, Binaries, Pakete, `tap`) am selben Tip grün. Der lokale Vorab-Beleg (DoD 2) liegt dem Verifier nicht vor und ist in der Closure-Notiz zu benennen. |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) | Cut-Commit und Fix-Commit nennen beide IDs und den Slice im Betreff |
| Commits `v0.21.0..HEAD` | 66 Commits; 65 nennen eine `LH-*`-, `ADR-*`- oder Slice-Kennung. Ausnahme: `97bebdc1` („docs/plan/cr: der CR-Ablage-Ort …“) trägt keine Kennung; er liegt vor dem Cut und gehört nicht zum Cut-Diff (Hinweis, kein Finding dieses Slice) |
| Slice-Moves | `472bf3e0`, `f4c9e51e`, `5d9c3d66`, `c0033c7e`: Move-Commit rein, Verweis-Nachzug als eigener Commit (§3.3) |

Docs-Links: `make docs-check` im `gates`-Lauf grün (0 Befunde, vor und nach diesem Report zu halten).

## Carveouts

- Neu: **keine**; keine Lockerung von Gates, Schwellen oder Architekturregeln (kein Pin-Bump, keine Pipeline-Änderung im Cut).

## Rest-Risiken / offene Punkte vor der Closure

1. **Veröffentlichter Lead weicht ab.** Der Tag `v0.22.0` (`f1711f3a:CHANGELOG.md`) trägt den früheren Lead „no non-test file under `cmd/` or `internal/` differs from 0.21.0“, was `main.go` nicht deckt; HEAD ist berichtigt. Der Tag bleibt, die Release-Seite verweist auf `CHANGELOG.md` von `main` (dort der berichtigte Text). Die Closure-Notiz hält die Abweichung zwischen Tag-Stand und `main` als Ausgang zu Risiko 1 fest.
2. **Mac-Probe ohne Protokoll** (M7): nur Aussage des Projektinhabers; in der Closure als „ohne Protokoll“ führen, nicht als belegt.
3. **Lokale Vorab-Belege** (`make ci`, `make fullbuild`, `make build VERSION=0.22.0`, M-T-Trockenlauf) sind nicht als Log vorgelegt; der CI-Nachweis (ci, integration-docker, tap-check, publish mit tap) steht dafür, deckt aber nicht den Trockenlauf-Schritt selbst.
4. **Closure-Pflichten offen** (kein Finding): §7 leer; DoD-Haken offen; Roadmap-Marker „Push und Tag stehen aus“ abgelaufen → auf Ist-Stand setzen und *Abgeschlossene Wellen* mit Tag-Hash `f1711f3a` ergänzen; §6 Risiken 1 bis 6 ohne Ausgang (Belege für 2 bis 4 siehe oben: Datum = Tag-Datum, `integration-docker` grün, `tap` grün); §8 Beobachtungs-Sichtung „nicht durchgegangen“; Verschiebung nach `done/` erst nach der Notiz (§3.3 Fall 2); `make gates` nach der Closure.
5. **Versionsnummer** (Risiko 6, Minor statt Patch): Entscheidung des Projektinhabers, Bestätigung im Eingang nicht belegt; Ausgang in der Closure-Notiz festhalten.

## Gesamturteil

**Closure-fähig unter Bedingungen.** Der Cut-Diff deckt Plan §2 Punkt 1 und die Plan-Tabelle; die Lead-Aussage im HEAD trägt gegen `v0.21.0..v0.22.0`; alle vier CI-/Release-Läufe am Tip `f1711f3a` sind grün, das Release trägt die erwarteten zwölf Assets, der Tag steht auf dem Cut-Commit; `make gates` ist am HEAD `353bef50` grün. Vor der Closure: Closure-Notiz mit Ausgängen zu den Rest-Risiken 1 bis 3 und 5, Roadmap auf den Ist-Stand, `make gates` danach.

## Commit / Artefakt

- Slice: `f1711f3a` (Cut, Tag `v0.22.0`), `353bef50` (Review-Fix-Runde)
- HEAD/Stempel: `353bef5039433efc0f160ebd55883a9c6ee26354`
- Lokal: `make gates` EXIT=0 am 2026-11-04
