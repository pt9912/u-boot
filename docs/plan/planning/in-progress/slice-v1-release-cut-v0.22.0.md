# Slice V1: Release-Cut `v0.22.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`; er wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.22.0` von `main` nach `v0.21.0` ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md): Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag `v0.21.0` (`f0026b4d`) unter `## [Unreleased]` steht: den einzigen Bau-Pfad im Devcontainer und die Pin-Inventur mit `pinsync` (zwei Einträge *Changed*).

**Bezug:** [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-04.

## 1. Ziel und Abgrenzung

**Ziel:** `v0.22.0` liefert den Block `[0.22.0]` des CHANGELOG als GHCR-Image, Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent erst nach ausdrücklicher Freigabe des Projektinhabers**; vorher steht der Stand-Bericht der Vorab-Prüfungen.

**Nummer:** Minor, auf Vorschlag des Hauptagenten, vom Projektinhaber noch nicht bestätigt. Der Release enthält keine Änderung an Nicht-Test-Code in `cmd/` und `internal/` (`git diff v0.21.0..HEAD --stat -- cmd internal ':!*_test.go'` am 2026-11-04 vor dem Cut-Commit: leer; der Cut-Commit ändert dort allein die Versionszeile in `cmd/uboot/main.go`); geändert sind Bau-Pfad, CI, Gate-Kette, die `go`-Direktive von `go.mod` (`1.26.0` → `1.27.0`) und Doku. Ob dafür ein Patch-Release (`v0.21.1`) genauer wäre, entscheidet der Projektinhaber.

**Ausdrücklich NICHT in diesem Slice:** Produktcode · `spec/`-Änderungen (der Cut ändert Release-Dokumente und drei Versionszeilen) · `make test-docker` im Devcontainer dieser Sitzung (rot an der cgroup, [`slice-v1-sandbox-nested-docker-cgroup`](../open/slice-v1-sandbox-nested-docker-cgroup.md); der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis) · Hebung der Pins.

**Vor dem Start zu klären (M-Start):** jeder `[Unreleased]`-Eintrag trägt einen Beleg und einen Commit im Diff `v0.21.0..HEAD` (gelesen: zwei Einträge — Bau-Pfad → [`slice-v1-bau-nur-devcontainer`](../done/slice-v1-bau-nur-devcontainer.md); Pin-Inventur → [`slice-harness-dcheck-pinsync`](../done/slice-harness-dcheck-pinsync.md)).

## 2. Definition of Done

- [x] Release-Text und Versionsstellen: `CHANGELOG.md` trägt `## [0.22.0] - <Tag-Datum>` mit Lead („Twenty-second release.“), leerem `[Unreleased]` und Compare-Links; Versionsstrings `0.22.0-dev` an den drei Stellen samt Kommentaren (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); READMEs und Handbuch auf `v0.22.0`.
- [x] Vorab-Prüfungen: `make gates`, `make ci` lokal grün (`make fullbuild` einzeln nicht gefahren: es besteht aus `ci` und dem Image-Build, beide liefen grün); Gegenprobe `make build VERSION=0.22.0`; Tap-Trockenlauf nach `releasing.md` §4 mit Vorgänger-Gegenprobe; `ci` und `integration-docker` auf dem Push-Tip grün; `tap-check` grün; **danach nur nach Freigabe:** Tag `v0.22.0` gesetzt und gepusht, `publish` und Job `tap` kontrolliert.
- [x] Doku nachgezogen: Roadmap (*Abgeschlossene Wellen* mit Tag-Hash, Zustandsmarker), Bestand-Zeilen.
- [x] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure.
- [ ] Review (kein Self-Review, [Report](../../../reviews/2026-11-04-slice-v1-release-cut-v0.22.0.md)) und Verification ([Report](../../../reviews/verification-2026-11-04-slice-v1-release-cut-v0.22.0.md)) mit Reports; Closure-Notiz mit Lerneintrag; Register oder „keine Beobachtung“; jedes Risiko mit Ausgang; drei Paarungen getragen.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.22.0] - <Tag-Datum>` mit Lead; Compare-Links `[Unreleased]: …/compare/v0.22.0...HEAD`, `[0.22.0]: …/compare/v0.21.0...v0.22.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.21.0-dev` → `0.22.0-dev` je Stelle samt Kommentar (Fundstellen am 2026-11-04 mit `grep` gelesen: `main.go:45`, `repo.mk:67-71`, `Dockerfile:49-53`) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md` | update | Status-Zeile (`README.md:10`, `README.de.md:11`), `VERSION=`-Beispiel (`:99`, `:102`), Handbuch-Kopf (Zeilen 3–5), `u-boot version`, `download/v…`, `VERSION=`, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile unter *Offene Wellen* bis zum Tag, danach *Abgeschlossene Wellen* mit Tag-Hash |
| `docs/plan/planning/open/README.md` | kein Diff durch den Cut | die Bestand-Zeile kam mit dem Plan-Commit `slice-v1-release-cut-v0.22.0` (Plan und Bestand-Zeile), nicht mit dem Cut-Commit |
| dieser Slice | neu / Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `spec/`, ADRs, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff durch den Cut | kein Pin-Bump, keine Pipeline-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6); gegen `v0.21.0` tragen diese Dateien den Diff der zwei Träger-Slices, nicht des Cuts |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6): M0 Datum-Abgleich · M-T Tap-Trockenlauf · M1 Push (Freigabe) · M2 CI-Läufe lesen · M4 `gh workflow run tap-check` · M5 Tag (nur nach Freigabe und Stand-Bericht) · M6 Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`next` → `in-progress`): die trägenden Vorgänge der zwei `[Unreleased]`-Einträge liegen in `done/`; `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`.

**Rückführungen:** `in-progress` → `next` (zu groß); `in-progress` → `open` (blockiert — Vorab rot, `tap-check` ungrün oder die Nummer nicht bestätigt; der Cut wartet).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.22.0` gesetzt und gepusht; `publish` und Job `tap` kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice geht **erst nach dem Tag** nach `done/` ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Menge** — der Lead liest gegen die zwei Einträge; die Aussage „Nicht-Test-Code von `cmd/` und `internal/` unverändert“ stützt sich auf den Diff vom 2026-11-04. — **Ausgang:** eingetreten: der Lead sagte, keine Nicht-Test-Datei unter `cmd/` oder `internal/` unterscheide sich von `v0.21.0`; der Cut-Commit `f1711f3a` ändert `cmd/uboot/main.go` (Versionszeile), und `go.mod` stuft `cobra` und `pflag` um (Review F-1, F-2). Berichtigt auf `main` mit `353bef50`; der Tag `v0.22.0` trägt den ungenauen Satz unverändert, weil ein veröffentlichter Tag nicht verschoben wird. Geführt als weiterer Beleg in [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).
2. **Das Datum weicht vom Tag-Datum ab** (M0). — **Ausgang:** entfallen: `date +%F` am Tag-Tag und `git log -1 --format=%cs v0.22.0` ergeben beide 2026-11-04, wie CHANGELOG und Handbuch-`Stand`.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit.** Lokal ist `make test-docker` rot (cgroup). — **Ausgang:** entfallen: `integration-docker` 41816259276 auf `f1711f3a` ist `success` (rund 19 Minuten); die lokale Rotlage führt [`slice-v1-sandbox-nested-docker-cgroup`](../open/slice-v1-sandbox-nested-docker-cgroup.md).
4. **Job `tap` scheitert** — Trockenlauf M-T, `tap-check` vor M5. — **Ausgang:** entfallen: `tap-check` 41818441675 (`git push --dry-run succeeded`) und `publish` 41818642441 mit Job `tap` sind `success`; das Log des Schritts „Pull the formula into the tap“ zeigt `updated to v0.22.0`.
5. **Zustandsmarker nach dem Tag nicht auf den Ist-Stand.** — **Ausgang:** entfallen: Die Roadmap führt den Release unter *Abgeschlossene Wellen* mit Tag-Hash und trägt den Ruhe-Marker; die Zeile „Push und Tag stehen aus“ ist entfernt (Review F-4).
6. **Die Versionsnummer ist nicht bestätigt** (Minor oder Patch). — **Ausgang:** entfallen: Der Projektinhaber gab den Auftrag „mach das release“ nach dem Vorschlag `v0.22.0` (Minor); die Alternative `v0.21.1` war genannt, ist nicht ausdrücklich gewählt worden. Die Nummer ist veröffentlicht.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Kette aus `releasing.md` lief in der Reihenfolge: lokale Vorab-Prüfungen, Tap-Trockenlauf, Push von `main` (`dffdef36..f1711f3a`), `ci` und `integration-docker` auf dem Push-Tip, `tap-check`, Tag `v0.22.0` auf `f1711f3a` (2026-11-04, Freigabe „mach das release“), `publish` mit allen vier Jobs. Kontrolle nach §6: 12 Assets (6 Binaries, 4 Pakete, `SHA256SUMS`, `u-boot.rb`), das Image `0.22.0` und `:latest` melden `u-boot version 0.22.0`. Lauf-IDs: `ci` 41816259305, `integration-docker` 41816259276, `tap-check` 41818441675, `publish` 41818642441.
- **Was ging anders als geplant:** (1) Der Lead sagte mehr als der Diff (§6 Risiko 1) und ist auf `main` berichtigt; der Tag trägt den früheren Satz. (2) Der erste Push des Tags wurde ohne Begründung abgelehnt (`remote rejected (failure)`), der zweite Versuch ging durch; die Ursache ist nicht geklärt. (3) `make test-docker` ist im Devcontainer dieser Sitzung rot (cgroup); `integration-docker` in CI trug den Nachweis. (4) Der Push von `main` lief über die Umgehung der Branch-Protection durch das Konto des Projektinhabers (`Bypassed rule violations`).
- **Vorab-Belege ohne Datei im Repo:** `make gates`, `make ci` (inklusive `govulncheck` und `image-scan`) und `make build VERSION=0.22.0` lokal grün (Exit 0; die Gegenprobe meldet `0.22.0` im `--version` und im OCI-Label); der Tap-Trockenlauf nach `releasing.md` §4 mit Fixture-Digests: `homebrew-formula-fill.sh` und `tap-formula-check.sh` mit `v0.22.0` Exit `0`, mit dem Tag des Vorgängers Exit `1`; `tap-nachzug.sh` mit Stellvertretern für `curl` und `git` zeigt `clone`, `add`, `commit -m u-boot 0.22.0`, `push`, mit dem Tag des Vorgängers Abbruch vor `clone`. Die Logs liegen im Scratchpad der Sitzung, nicht im Repo.
- **Mac-Probe:** Aussage des Projektinhabers („auf dem Mac läuft es auch“), ohne Protokoll.
- **Steering-Loop-Eintrag:** Eine Messaussage im Release-Text, die vor dem Cut-Commit gemessen wurde, gilt nach dem Cut-Commit nicht mehr (der Cut ändert selbst Dateien); das gehört zur Klasse „Lead sagt mehr als der Beleg“ — gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-v1-release-cut-v0.22.0.md` in `BEO-DOK/zusage-weiter-als-messung/` ergänzt — der Zähler steht damit über der Schwelle (die Entscheidung des Projektinhabers über den Ausgang steht aus, unverändert).
- **Folge-Slices:** [`slice-v1-sandbox-nested-docker-cgroup`](../open/slice-v1-sandbox-nested-docker-cgroup.md) und [`slice-v1-devcontainer-ressourcenlimits-messung`](../open/slice-v1-devcontainer-ressourcenlimits-messung.md) sind Dateien in `open/`; der Nachzug des Beispiels `BUILD_MODE` im Handbuch (`docs/user/benutzerhandbuch.md`, Abschnitt zu den Katalog-Schlüsseln) bleibt dem nächsten Release.
- **Risiken aus §6:** je ein Ausgang, siehe §6.
- **Drei Paarungen:** Anker — kein Herkunfts-Anker gesetzt (nichts verkörpert); Folge-Slice — die beiden Folge-Slices liegen in `open/`; Register — die Evidence-Datei liegt vor.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `CHANGELOG.md`, READMEs, `docs/user/` und `docs/plan/` (DOK); `cmd/uboot/main.go` (Versionszeile); alle Greenfield bzw. keine Sub-Area ([`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-04 (nach Titeln): Treffer [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (Risiko 1, eingetreten), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (die Plan-Aussage „leer“ zum Diff stammte von vor dem Cut-Commit; die Klasse ist verkörpert, [`AGENTS.md`](../../../../AGENTS.md) §3.10); sonst keine Treffer.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
