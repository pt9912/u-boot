# Slice V1: Release-Cut `v0.23.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`; er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.23.0` von `main` nach `v0.22.0` ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md): Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag `v0.22.0` (`f1711f3a`) unter `## [Unreleased]` steht: die Umgebungsdatei des Devcontainers (`.devcontainer/devcontainer.env`, *Added*) und den cgroup-Umbau im Start-Skript des Sandbox-Devcontainers mit nested Docker (*Changed*).

**Bezug:** [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`ADR-0044`](../../adr/0044-devcontainer-umgebungsdatei-env-file.md), [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-04.

## 1. Ziel und Abgrenzung

**Ziel:** `v0.23.0` liefert den Block `[0.23.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent erst nach ausdrücklicher Freigabe des Projektinhabers**, und der Tag erst nach grünem `integration-docker` auf dem Push-Tip; vorher steht der Stand-Bericht der Vorab-Prüfungen.

**Nummer:** Minor — die Zusage-Menge wächst um eine Anforderung (Lastenheft 0.35.0, [`LH-FA-DEV-012`](../../../../spec/lastenheft.md#lh-fa-dev-012--umgebungsdatei-des-devcontainers)), und die erzeugte `devcontainer.json` aller Projekte trägt eine neue `runArgs`-Zeile. Im Unterschied zum Cut `v0.22.0` ändern 12 Nicht-Test-Dateien unter `internal/` das Produkt (`git diff v0.22.0..HEAD --stat -- cmd internal ':!*_test.go'` am 2026-11-04 vor dem Cut-Commit); der Cut-Commit ändert in `cmd/` und `internal/` allein die Versionszeile in `cmd/uboot/main.go`.

**Ausdrücklich NICHT in diesem Slice:** Produktcode · `spec/`-Änderungen · die Neuanlage des Entwicklungs-Devcontainers und die Messung nach ihr ([`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md); Entscheidung des Projektinhabers 2026-11-04: nach dem Release) · `make test-docker` im Devcontainer dieser Sitzung (der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis) · Hebung der Pins.

**Vor dem Start zu klären (M-Start):** jeder `[Unreleased]`-Eintrag trägt einen Beleg und einen Commit im Diff `v0.22.0..HEAD` (gelesen: zwei Einträge — Umgebungsdatei → [`slice-v1-devcontainer-env-datei`](../done/slice-v1-devcontainer-env-datei.md); cgroup-Umbau → [`slice-v1-sandbox-nested-docker-cgroup`](../done/slice-v1-sandbox-nested-docker-cgroup.md)). Der Test-Fix [`slice-v1-test-host-umgebung-isolieren`](../done/slice-v1-test-host-umgebung-isolieren.md) liegt schon im Tag `v0.22.0` und betrifft nur Tests.

## 2. Definition of Done

- [ ] Release-Text und Versionsstellen: `CHANGELOG.md` trägt `## [0.23.0] - <Tag-Datum>` mit Lead („Twenty-third release.“), der nur sagt, was die zwei Einträge und der Diff belegen, leerem `[Unreleased]` und Compare-Links; Versionsstrings `0.23.0-dev` an den drei Stellen samt Kommentaren (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); READMEs und Handbuch auf `v0.23.0`.
- [ ] Vorab-Prüfungen: `make gates`, `make ci` lokal grün; Gegenprobe `make build VERSION=0.23.0`; Tap-Trockenlauf nach `releasing.md` §4 mit Vorgänger-Gegenprobe; `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach nur nach Freigabe:** Tag `v0.23.0` gesetzt und gepusht, `publish` und Job `tap` kontrolliert.
- [ ] Doku nachgezogen: Roadmap (*Abgeschlossene Wellen* mit Tag-Hash, Zustandsmarker), Bestand-Zeilen.
- [ ] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review) und Verification mit Reports; Closure-Notiz mit Lerneintrag; Register oder „keine Beobachtung“; jedes Risiko mit Ausgang; drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.23.0] - <Tag-Datum>` mit Lead; Compare-Links `[Unreleased]: …/compare/v0.23.0...HEAD`, `[0.23.0]: …/compare/v0.22.0...v0.23.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.22.0-dev` → `0.23.0-dev` je Stelle samt Kommentar (Fundstellen beim Cut mit `grep` gelesen) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status-Zeile, `VERSION=`-Beispiel, Handbuch-Kopf (`Software-Version`, `Stand`), `u-boot version`, `download/v…`, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile unter *Offene Wellen* bis zum Tag, danach *Abgeschlossene Wellen* mit Tag-Hash |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile dieses Slice (mit dem Plan-Commit) |
| dieser Slice | neu / Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff durch den Cut | kein Pin-Bump, keine Pipeline-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6): M0 Datum-Abgleich · M-T Tap-Trockenlauf · M1 Push (Freigabe) · M2 CI-Läufe lesen · M4 `gh workflow run tap-check` · M5 Tag (nur nach Freigabe, grünem `integration-docker` und Stand-Bericht) · M6 Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`next` → `in-progress`): die trägenden Vorgänge der zwei `[Unreleased]`-Einträge liegen in `done/`; `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`.

**Rückführungen:** `in-progress` → `next` (zu groß); `in-progress` → `open` (blockiert — Vorab rot, `tap-check` ungrün oder `integration-docker` rot; der Cut wartet).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.23.0` gesetzt und gepusht; `publish` und Job `tap` kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Menge** — der Lead liest gegen die zwei Einträge und gegen den Diff `v0.22.0..v0.23.0` nach dem Cut-Commit (beim Cut `v0.22.0` stimmte eine Aussage nicht, weil sie vor dem Cut-Commit gemessen war). — **Ausgang:** offen → bei Closure.
2. **Das Datum weicht vom Tag-Datum ab** (M0). — **Ausgang:** offen → bei Closure.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit.** Lokal ist `make test-docker` in diesem Devcontainer bis zur Neuanlage rot (`TestE2E_Resources_*`). — **Ausgang:** offen → bei Closure.
4. **Job `tap` scheitert** — Trockenlauf M-T, `tap-check` vor M5. — **Ausgang:** offen → bei Closure.
5. **Zustandsmarker nach dem Tag nicht auf den Ist-Stand.** — **Ausgang:** offen → bei Closure.
6. **Der Release trägt einen cgroup-Umbau, dessen Verhalten nach einer Neuanlage ungemessen ist** (Risiken 2 bis 4 des Vorgänger-Slice, geführt im Folge-Slice); [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) ist ab dem Push nicht mehr in place änderbar, ein Befund führt zu einer Folge-ADR mit `Supersedes`. Der CHANGELOG führt die Grenze unter *Not promised*. — **Ausgang:** offen → [`slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung`](../open/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung.md).
7. **Der Release enthält eine Ausgabe-Änderung für alle Projekte** (neue `runArgs`-Zeile) und in einem frischen Klon fehlt die Datei bis zu einem `u-boot generate devcontainer`. — **Ausgang:** offen → bei Closure (der CHANGELOG führt beides).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `CHANGELOG.md`, READMEs, `docs/user/` und `docs/plan/` (DOK); `cmd/uboot/main.go` (Versionszeile); alle Greenfield bzw. keine Sub-Area ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)).

**Vorgelagert — offene Beobachtungen sichten:** nicht durchgegangen; die Sichtung erfolgt beim M-Start, bis dahin steht hier kein Befund.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
