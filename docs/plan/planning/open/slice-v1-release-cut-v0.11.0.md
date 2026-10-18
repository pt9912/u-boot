# Slice V1: Release-Cut `v0.11.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.11.0 (Dry-Run-Formulierung, fehlendes `/dev/fuse` als Startfehler, host-seitige Ausweisung der Degradations-Fallbacks, Handbuch-Anleitung für `sshd`; Release-Zuordnung des Projektinhabers vom 2026-10-18). Der Release trägt alles, was seit dem Tag `v0.10.0` auf `main` liegt: [`slice-v1-cli-dry-run-would-formulierung`](../done/slice-v1-cli-dry-run-would-formulierung.md), [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md), [`slice-v1-sandbox-fallback-host-ausweisung`](../done/slice-v1-sandbox-fallback-host-ausweisung.md) und [`slice-v1-handbuch-sshd-anleitung`](../done/slice-v1-handbuch-sshd-anleitung.md); Slices ohne Produktänderung (Harness, Planung) ändern den Inhalt nicht.

**Bezug:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run),
[`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging),
[`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md),
[`ADR-0007`](../../adr/0007-distributionswege-ghcr.md),
[`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md),
[`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; Lastenheft 0.13.0 liefern die Slices, dieser Slice ändert `spec/` nicht).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.11.0` liefert als Minor-Release das Dry-Run-Vokabular „Would …“ für `remove`, `generate` und `config set`, das fehlende `/dev/fuse` als Startfehler, die host-seitige Ausweisung der Degradations-Fallbacks in `doctor` und `generate devcontainer` und die Handbuch-Anleitung für den `sshd`: GHCR-Image, Binaries, Pakete, Homebrew-Formel.

### Inhalt des Release (gelesen am 2026-10-18, `CHANGELOG.md` Block `[Unreleased]`)

| Eintrag im CHANGELOG | Art | Hinweis für den Lead |
|---|---|---|
| Sandbox devcontainer: Fallbacks (Daemon, `vfs`, Egress) host-seitig ausgewiesen (`doctor`, `generate devcontainer`; `init` fragt, hat aber nichts abzufragen) | Added | Bestehende Projekte bekommen die Statusdateien nach `u-boot generate devcontainer` und einem Neubau des Containers, bis dahin ist der Zustand „nicht ermittelbar“; Messstand Colima, arm64 |
| Fehlendes `/dev/fuse` auf dem Container-Host ist ein Startfehler der Engine; `doctor` meldet es auf einem Linux-Host als Fehler, Exit-Code `11` | Changed | **Verhaltensänderung:** Skripte, die `doctor` ohne `--strict` auswerten, scheitern dort künftig; Messstand Colima, arm64 |
| `--dry-run`: Zusammenfassung von `remove`, `generate`, `config set` im „Would …“-Wortlaut | Changed | Skripte, die die alte Zeile abglichen, nutzen `--json` |

Nicht im CHANGELOG und nicht im Lead: `-timeout 30m` im Ziel `test-docker` von `repo.mk` — das Zeitbudget gehört zur CI und zum Entwicklungs-Lauf, es ändert kein Verhalten des Produkts (gelesen: `go test -tags docker -timeout 30m ./...` in `repo.mk`, Ziel `test-docker`).

**Entscheidung des Projektinhabers (2026-10-18):** `v0.11.0` enthält die vier Inhalte oben. Die beiden Hygiene-Slices [`slice-v1-cli-dry-run-log-would`](slice-v1-cli-dry-run-log-would.md) und [`slice-v1-harness-gofmt-und-depguard`](slice-v1-harness-gofmt-und-depguard.md) sind **nicht** Teil von `v0.11.0`; sie folgen danach.

### Lehren aus dem Cut `v0.10.0`, die dieser Slice einplant (gelesen in [`slice-v1-release-cut-v0.10.0`](../done/slice-v1-release-cut-v0.10.0.md) §6 und §7)

- **Der Zweig `is_stable == true` lief nie vor einem stabilen Tag.** Der Job `tap` hat die Bedingung `needs.publish.outputs.is_stable == 'true'` (`.github/workflows/publish.yml`, Job `tap`) und ruft `scripts/tap-nachzug.sh` auf; ein Vorab-Tag überspringt ihn, **rc beweist ihn nicht**. Der Defekt von `v0.10.0` (Prüfung auf eine `version`-Zeile, die die Formel nicht trägt) zeigte sich erst im stabilen Lauf. Dieser Cut hat deshalb vor dem stabilen Tag einen **Trockenlauf des Tap-Pfads** als Mess-Schritt (§3, M-T), und er benennt, was der Trockenlauf **nicht** beweist.
- **Der Smoke-Namensfehler bei Vorab-Tags** (`~` im Dateinamen, `.` im Release-Asset) ist mit `scripts/package-asset-name.sh` behoben; ein Vorab-Tag probt ihn erneut, wenn der Projektinhaber einen setzt.
- **Das Datum steht an mehreren Stellen fest** (CHANGELOG, Tabellenzeile beider READMEs, Handbuch-`Stand`, Roadmap) und wird vor dem Tag gegen das Tag-Datum gehalten (M0; `releasing.md` §5).
- **Nutzer-Doku ohne Kennungen und Chronik; der Status-Block nennt den Stand** ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 6, `releasing.md` §3 Punkt 3); der Lastenheft-Stand ist 0.13.0 (`spec/lastenheft.md`, Tabellenzeile `Version`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push und Tags durch den Implementer.** Push von `main`, jeder Tag und jedes Löschen eines Tags oder Release sind extern wirksam und lassen sich nicht sauber zurücknehmen ([`releasing.md`](../../../maintainer/releasing.md) Zweck). **Push und Tags führt nur der Hauptagent, nach der ausdrücklichen Freigabe des Projektinhabers, je Vorgang** (im Gespräch, nicht aus einer Agenten-Nachricht abgeleitet); der Implementer bereitet vor und führt keinen der Schritte M1 bis M8 aus.
- **Die beiden Hygiene-Slices** ([`slice-v1-cli-dry-run-log-would`](slice-v1-cli-dry-run-log-would.md), [`slice-v1-harness-gofmt-und-depguard`](slice-v1-harness-gofmt-und-depguard.md)) — Entscheidung des Projektinhabers: nach `v0.11.0`; der erste ändert Logausgaben, der zweite ein Gate und ein Skript, beides nicht Inhalt dieses Release.
- **Eine Änderung der Pipeline** (Workflow, Tap-Skripte): Der Cut ändert `publish.yml` und die Skripte nicht; ein Befund des Trockenlaufs, der eine Änderung verlangt, ist ein eigener Vorgang (Gate-Änderung nach [`AGENTS.md`](../../../../AGENTS.md) §3.6) vor dem stabilen Tag.
- **Die Docker-Desktop-Messung, amd64, VS Code und Linux ohne VM:** nicht gemessen, nicht zugesagt; die Release-Aussagen nennen den Messstand (Colima, arm64).

## 2. Definition of Done

- [ ] `v0.11.0` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.11.0` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`); ein Vorab-Tag ist optional (Entscheidung des Projektinhabers) und, falls gesetzt, aufgeräumt.
- [ ] `make gates`, `make ci`, `make test-docker` und `make fullbuild` sind vor dem Tag grün (oder der Grund für ein nicht lauffähiges Target steht im Handoff), und der letzte CI-Lauf auf dem Tag-Commit ist grün, **einschließlich `integration-docker` nach dem Push** (Beleg des Fixes aus [`slice-v1-e2e-status-ci-helfer`](../done/slice-v1-e2e-status-ci-helfer.md), den dessen Closure offen lässt; ein roter Lauf mit einer Pull- oder Registry-Meldung ist ein neuer Befund, kein Beleg gegen den Fix); [`slice-v1-podman-aktuelle-version`](../in-progress/slice-v1-podman-aktuelle-version.md) liegt in `done/`; der Trockenlauf des Tap-Pfads (M-T) ist ausgeführt und sein Ergebnis samt Grenze in §7 festgehalten.
- [ ] Doku-Update: CHANGELOG, READMEs, Handbuch und Roadmap tragen `v0.11.0`; die Release-Aussagen nennen den Messstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Jede Aussage in dieser Tabelle ist an der genannten Stelle gelesen (Stand 2026-10-18, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Zeilennummern nennen den Fundort beim Schreiben; der Implementer sucht sie vor der Änderung neu. Vor dem Handoff (nach der letzten Fix-Runde) gleicht er die Tabelle gegen `git diff --stat` ab; eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage, nicht durch eine Suche.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update | `:12` `## [Unreleased]` (Inhalt bis vor `## [0.10.0] - 2026-10-17`) wird zum leeren Anker; Inhalt unter `## [0.11.0] - <Tag-Datum>` (vorläufig 2026-10-18, vor dem Tag mit `date +%F` abgleichen) mit Lead-Absatz im Stil der Vorgänger („Tenth release. **…**“ steht unter `[0.10.0]`; hier „Eleventh release.“); der Lead nennt alle drei Einträge der Tabelle in §1 samt Verhaltensänderung (`doctor` Exit-Code `11` ohne `--strict`), „Lastenheft 0.13.0“ und den Messstand (Colima 0.10.3, Docker 29.5.2, arm64; nicht gemessen: Linux ohne VM, amd64, Docker Desktop, VS Code). Compare-Links am Fuß (`:1291` `[Unreleased]: …/compare/v0.10.0...HEAD` → `v0.11.0...HEAD`; neue Zeile `[0.11.0]: …/compare/v0.10.0...v0.11.0` vor `:1292`). Eine Zeile zum Handbuch nur, wenn der Projektinhaber es will (§6). |
| `cmd/uboot/main.go` | update | `:43` `var version = "0.10.0-dev"` → `0.11.0-dev`. Der Beispielaufruf im Kommentar darüber (`-X main.version=v0.2.0`) ist kein Versionsstand und bleibt. |
| `repo.mk` | update | `:51` und `:54` (Kommentare „`var version = "0.10.0-dev"`“, „"0.10.0-dev" binary“) und `:55` `VERSION ?= 0.10.0-dev` → `0.11.0-dev`. |
| `Dockerfile` | update | `:35` und `:37` (Kommentar) und `:39` `ARG UBOOT_VERSION=0.10.0-dev` → `0.11.0-dev`. Weitere Treffer von `0.10.0-dev` im Baum gibt es nach dem Stand dieses Plans nicht (gesucht in Quellen, Skripten, Workflows und Doku außer `done/`, Reviews und Baseline; die Versionsangaben `0.10.0.rc.1` in `scripts/package-asset-name.sh` und `cmd/uboot/package_asset_name_test.go` sind Beispiele für ein Asset-Namensformat, kein Versionsstand, und bleiben). |
| `README.md` | update | Status-Block `:10`–`:14` (als Stand formuliert, ohne Datum und ohne „released“): Version `v0.11.0`; er nennt den Stand zur Ausweisung der Fallbacks, zum Startfehler bei fehlendem `/dev/fuse` und zum Dry-Run-Wortlaut, mit Messstand; `:99` `VERSION=0.10.0` → `0.11.0`; neue Zeile in der Releases-Tabelle unmittelbar nach der `v0.10.0`-Zeile (`:199`; die Tabelle ist nicht nach Version sortiert), Link `releases/tag/v0.11.0`; ohne Anforderungs-, ADR- oder Spec-Kennungen. Die Feature-Zahl „13 features“ (`:11`) bleibt: der Katalog ändert sich in diesem Release nicht (gelesen: `featureCatalogue()` in `internal/hexagon/application/devcontainer_features.go`, Schlüssel unverändert). |
| `README.de.md` | update | gleiches Muster: `:11`–`:15` Status, `:101` `VERSION=0.10.0`, Tabellenzeile nach `:202`. |
| `docs/user/benutzerhandbuch.md` | update | `:4` `Software-Version: v0.10.0`, `:5` `Stand: 2026-10-17` (→ Tag-Datum), `:80` Beispielausgabe `u-boot version 0.10.0`, `:106` `download/v0.10.0`, `:133` `VERSION=0.10.0`, `:2032` „beschreibt `u-boot` v0.10.0“. Die Anleitung für den `sshd` kommt aus [`slice-v1-handbuch-sshd-anleitung`](../done/slice-v1-handbuch-sshd-anleitung.md) und wird hier nicht geschrieben. Keine Kennungen, keine Chronik im Text. |
| `docs/maintainer/releasing.md` | prüfen, nur bei Abweichung ändern | Gelesen §1–§8. §3 nennt als „zuletzt“ den Cut `v0.9.0` (`slice-v1-release-cut-v0.9.0`); der letzte Cut ist `v0.10.0`, nach diesem `v0.11.0` → berichtigen. Die übrigen Aussagen stimmen mit dem Ist-Stand überein (Tap-Prüfung über Download-URLs, Asset-Name aus der Asset-Liste). Eine Aussage zum Trockenlauf des Tap-Pfads kommt in §4 (Vorab-Prüfung) nur, wenn der Projektinhaber ihn als Dauerschritt will (§6). |
| `docs/plan/planning/in-progress/roadmap.md` | update | Beim Start: `Nichts in Arbeit.` (Ruhe-Marker, gilt nur ohne Slice in `in-progress/`, `.d-check.yml` `planning`) entfällt, stattdessen ein Zeiger auf diesen Slice im Block „Offene Wellen“; nach dem Tag: Marker zurück, neue Zeile in „Abgeschlossene Wellen“ nach der Zeile `v0.10.0`: `v0.11.0 (…)` mit Datum und Verweis auf diesen Slice als Detailquelle, Tag-Commit-Hash erst nach dem Tag; `**Letzte Änderung:**` im Kopf. Der Hinweis „Nach v0.11.0: …“ für die beiden Hygiene-Slices bleibt unverändert stehen. |
| `docs/plan/planning/in-progress/README.md`, `docs/plan/planning/open/README.md` | update | Bestandszeile für diesen Slice beim Start (`in-progress/README.md`) und beim Move nach `done/` entfernt; `open/README.md` verliert die Zeile beim Start. |
| `docs/reviews/<Datum>-slice-v1-release-cut-v0.11.0.md`, `docs/reviews/<Datum>-verification-slice-v1-release-cut-v0.11.0.md` | neu, durch Reviewer und Verifier | Review- und Verification-Report des Cuts (nicht im Implementer-Diff, aber im Cut-Diff). Kennungen in Tabellenzellen ohne Backticks im Linktext, `make docs-check` vor dem Push. |
| `.github/workflows/publish.yml`, `scripts/tap-nachzug.sh`, `scripts/tap-formula-check.sh`, `scripts/homebrew-formula-fill.sh`, `packaging/` | kein Diff (gelesen) | Gelesen: `publish.yml` Job `tap` (Bedingung `is_stable == 'true'`, Schritt „Pull the formula into the tap“), `tap-nachzug.sh` (lädt `u-boot.rb`, ruft `tap-formula-check.sh`, klont, schreibt, pusht), `tap-formula-check.sh` (genau vier `url`-Zeilen mit `/releases/download/<tag>/`). Der Cut ändert sie nicht; der Trockenlauf (M-T) liest sie nur. |
| `docs/plan/planning/observations/` | je nach Befund | Ein Befund des Trockenlaufs oder der Tag-Läufe, der unter `BEO-WZG/kopplung-nicht-nachgezogen` zählt (`evidence/slice-v1-release-cut-v0.11.0.md`), kommt in der Closure; kein Eintrag im Implementer-Lauf. |

**Aussage zum Messstand — Soll, in CHANGELOG, Handbuch und Release-Zeile der READMEs gleich:** Gemessen mit Docker 29.5.2 unter Colima 0.10.3 (arm64, VM Ubuntu 24.04). **Nicht gemessen:** Linux ohne VM, amd64, Docker Desktop, VS Code und CI. Der Implementer schreibt nur, was das Messprotokoll von [`ADR-0022`](../../adr/0022-fallback-ausweisung-statusdateien-im-container.md) und der Review-Report von [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) belegen.

**Reihenfolge der Maintainer-Aktionen.** Der Implementer bereitet vor (Commits lokal, Gates lokal) und **führt keinen der Schritte M1 bis M8 aus**. **Push und Tags führt nur der Hauptagent nach Freigabe des Projektinhabers, je Vorgang.** Der Slice bleibt bis zum Tag in `in-progress/`.

| # | Aktion | Wer | Prüfung / Erwartung |
|---|---|---|---|
| M0 | `date +%F`; Datum in CHANGELOG, READMEs (Tabellenzeile), Handbuch (`Stand`) und Roadmap gegen das Tag-Datum halten | Implementer, vor jedem Push | Das Datum steht an mehreren Stellen fest ([`slice-v1-release-cut-v0.10.0`](../done/slice-v1-release-cut-v0.10.0.md) §6, `releasing.md` §5) |
| M-T | **Trockenlauf des Tap-Pfads, vor M1 und vor jedem Tag, ohne Netz und ohne Schreiben ins Tap:** (1) `scripts/homebrew-formula-fill.sh v0.11.0 <SHA256SUMS-Fixture> scripts/homebrew-formula.rb.tmpl <Ausgabe>` erzeugt eine Formel aus einer Fixture (Muster `sumsFull` in `cmd/uboot/homebrew_scripts_test.go`); (2) `scripts/tap-formula-check.sh v0.11.0 <Ausgabe>` endet mit Exit `0`, mit `v0.10.0` mit Exit `1`; (3) `scripts/tap-nachzug.sh` läuft mit Stellvertretern für `curl` und `git` im `PATH` (der Stellvertreter von `curl` liefert die Fixture-Formel, der von `git` protokolliert `clone`, `add`, `commit`, `push` und schreibt nichts) und mit gesetzten `TAP_TOKEN` und `TAG`. Alles im Docker (Image oder `make test`), nicht auf dem Host | Implementer; Ergebnis in §7 | Erwartet: (1) und (2) grün; (3) zeigt, dass das Skript die Formel prüft und den Commit-Befehl mit dem Tag-Namen bildet. **Was der Trockenlauf nicht beweist:** die Bedingung `needs.publish.outputs.is_stable == 'true'` und den Schritt im Workflow (sie laufen nur in einem echten stabilen Tag-Lauf), den echten `curl` gegen das Release-Asset, den echten `git push` und das Token. Trägt (3) mit Stellvertretern nicht, steht das als Grenze in §7; (1) und (2) sind der Mindestumfang. Ein Befund ist ein eigener Vorgang vor dem stabilen Tag (§1) |
| M1 | `git push origin main` | Hauptagent, nach Freigabe des Projektinhabers | Review-Report und Cut-Commit liegen vor; `make docs-check` lokal grün |
| M2 | `gh run list --limit 3`: `ci` und `integration-docker` auf dem Push-Tip grün | Maintainer | Roter `image-scan`: Go-Pin heben (`releasing.md` §4), neuer Commit, M1 wiederholen |
| M3 | Optional: `git tag v0.11.0-rc.1 && git push origin v0.11.0-rc.1` | Hauptagent, nach Freigabe des Projektinhabers | `publish`-Lauf: Log-Zeile `Resolved version=0.11.0-rc.1, is_stable=false`; Job `tap` übersprungen; Release als Prerelease, kein `:latest`, keine `u-boot.rb`; Paket-Smoke grün (probt den Asset-Namen mit `~` und `.`). **rc.1 beweist den Zweig `is_stable == true` nicht** |
| M4 | rc.1 aufräumen (nur wenn gesetzt): `gh release delete`, `git push --delete origin`, `git tag -d`, GHCR-Tag entfernen (`releasing.md` §5, §7) | Hauptagent, nach Freigabe des Projektinhabers je Löschung | Löschen ist nicht umkehrbar; GHCR-Löschung braucht Scope `delete:packages` |
| M5 | `git switch main && git pull --ff-only`, `git tag v0.11.0` auf dem in M2 geprüften Commit, `git push origin v0.11.0` | Hauptagent, nach Freigabe des Projektinhabers | Der Tag-Commit ist der Commit mit grünem CI; danach wird der Tag nicht verschoben (`releasing.md` §8) |
| M6 | `releasing.md` §6: `gh run watch`; `gh release view v0.11.0 --json assets`; `docker run --rm ghcr.io/pt9912/u-boot:0.11.0 --version`; `:latest` gleicher Digest; Job `tap` | Maintainer | Erwartet: zwölf Assets (sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`), kein Prerelease; **erster echter Lauf des Zweigs `is_stable == true` nach dem Trockenlauf:** `tap` läuft |
| M7 | Auf einem Mac: `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot`, `brew install u-boot`, `u-boot --version` = `0.11.0`; Smoke-Lauf im Tap abwarten | Projektinhaber | `brew audit --strict --online` gatet den Smoke (`releasing.md` §6) |
| M8 | Nacharbeit (`releasing.md` §7): Roadmap-Zeile mit Tag-Hash und Datum, Closure-Notiz, Verification Evidence, `git mv` nach `done/` | Planner (Closure), Implementer (Roadmap) | Closure erst nach M6 |

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-v1-handbuch-sshd-anleitung`](../done/slice-v1-handbuch-sshd-anleitung.md) liegt in `done/`; die Slices [`slice-v1-cli-dry-run-would-formulierung`](../done/slice-v1-cli-dry-run-would-formulierung.md), [`slice-cr-sandbox-fuse-startfehler`](../done/slice-cr-sandbox-fuse-startfehler.md) und [`slice-v1-sandbox-fallback-host-ausweisung`](../done/slice-v1-sandbox-fallback-host-ausweisung.md) liegen in `done/` (erfüllt, 2026-10-18). **Zwei Blocker** (Entscheidung des Projektinhabers vom 2026-10-18): [`slice-v1-e2e-status-ci-helfer`](../done/slice-v1-e2e-status-ci-helfer.md) liegt in `done/` (erfüllt mit seiner Closure; sein CI-Beleg steht in der DoD unten) und [`slice-v1-podman-aktuelle-version`](../in-progress/slice-v1-podman-aktuelle-version.md) liegt in `done/` (offen); kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Weg ist `open/` → `next/` → `in-progress/` per reinem `git mv` je Schritt (`make slice-mv SLICE=slice-v1-release-cut-v0.11.0 TO=<next|in-progress>`); der Roadmap-Marker (§3) geht in einen eigenen Commit direkt nach dem Move, beide im selben Push.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tags.
- `in-progress` → `open` (blockiert): ein roter `image-scan` oder CI-Lauf (Go-Pin heben, `releasing.md` §4), ein Befund des Trockenlaufs M-T, der eine Änderung an Workflow oder Skript verlangt (eigener Vorgang davor), ein fehlschlagender Vorab- oder `publish`-Lauf, dessen Ursache einen Patch-Release oder eine Korrektur der Pipeline vor dem stabilen Tag braucht (ein veröffentlichter Tag wird nicht verschoben), oder ein Befund des Reviewers, der eine Produktänderung verlangt (dann ein eigener Slice).

## 5. Closure-Trigger

DoD vollständig: stabiler Tag gesetzt (M5), `publish`- und Tap-Lauf kontrolliert (M6, `releasing.md` §6), Trockenlauf M-T mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs und Steering-Loop-Eintrag, Verification Evidence (Verifier, getrennt von Implementer und Reviewer). Der Slice geht **erst nach dem Tag** nach `done/` (Rollenreihenfolge: Implementer → Reviewer → Verifier → Maintainer-Aktionen → Closure durch den Planner); der Ruhe-Marker der Roadmap steht erst dann wieder.

## 6. Risiken und offene Punkte

Der Ausgang jedes Risikos wird bei der Closure festgeschrieben.

- **Der Job `tap` scheitert im stabilen Lauf, obwohl der Trockenlauf grün war:** Der Trockenlauf probt die Skripte, nicht die Workflow-Bedingung, den echten Download, den Push und das Token (M-T, „Was der Trockenlauf nicht beweist“). — **Ausgang:** offen bis zum stabilen Lauf; bei Fehler Behebung auf `main`, der Tag bleibt, Tap manuell nachziehen (`releasing.md` §8).
- **Der Paket-Smoke scheitert an der Namensbildung eines Vorab-Tags:** behoben mit `scripts/package-asset-name.sh`; ein Vorab-Tag probt es erneut. — **Ausgang:** offen; entfällt, wenn ein Vorab-Tag grün ist oder kein Vorab-Tag gesetzt wird (dann bleibt nur der stabile Lauf).
- **Das Release trägt eine Verhaltensänderung** (`doctor` Exit-Code `11` ohne `--strict` bei fehlendem `/dev/fuse` auf einem Linux-Host); ein Lead, der nur die Neuerungen nennt, verschweigt sie. — **Ausgang:** offen; Lead und `### Changed` tragen sie, der Review des Cuts prüft es.
- **Release-Aussagen ohne Messstand:** Die Fallback-Ausweisung ist auf Colima (arm64) gemessen; ohne Nennung läse sich die Aussage als allgemein. — **Ausgang:** offen; CHANGELOG, Handbuch und READMEs nennen den Messstand und das Nicht-Gemessene.
- **Das Tag-Datum weicht vom festen Datum an den Stellen aus M0 ab.** — **Ausgang:** offen; Prüfung in M0 vor jedem Push und in §7 von `releasing.md` danach.
- **Der Review-Report bricht `make docs-check`** (Linkform in Tabellenzellen). — **Ausgang:** offen; `make docs-check` vor M1.
- **Handbuch-Zeile im CHANGELOG** (ob die `sshd`-Anleitung einen Eintrag bekommt) und **Trockenlauf als Dauerschritt in `releasing.md` §4**: Entscheidungen des Projektinhabers, die der Plan nicht trifft. — **Ausgang:** offen.
- **GHCR-Tags früherer Vorab-Versionen** (`0.10.0-rc.1`, `0.10.0-rc.2`) stehen noch (Scope der `gh`-Anmeldung, [`slice-v1-release-cut-v0.10.0`](../done/slice-v1-release-cut-v0.10.0.md) §7); der Cut berührt sie nicht. — **Ausgang:** weiter offen beim Projektinhaber.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, `docs/maintainer/`, `docs/reviews/` und Release-Dateien (`repo.mk`, `Dockerfile`) — Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006). `cmd/uboot` ist laut Modus-Deklaration keine Sub-Area (Aspirantin); berührt wird dort eine Zeile (`var version`), keine Konvention.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen ([`../observations/`](../observations/README.md)): `BEO-WZG/kopplung-nicht-nachgezogen` ist `offen`, Zähler 1 (Skript, das einer geänderten Form nicht folgt; der Trockenlauf M-T ist die Antwort auf den Tap-Teil), `BEO-E2E/ci-ueberspringt-nested-podman-build` ist `offen`, Zähler 1 (betrifft die Release-Aussage zu amd64 und CI: sie bleibt „nicht gemessen“), `BEO-DOK/plan-aussage-ueber-code` und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` sind `verkörpert` (`AGENTS.md` §3.10: Aussagen über Code mit Fundstelle).

Alle berührten Sub-Areas GF.
