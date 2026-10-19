# Slice V1: Release-Cut `v0.12.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** v0.12.0 (Egress-Anzeige `u-boot egress list`, Subdomain-/Wildcard-Klärung, Neuaufbau statt Neu-Laden, Katalog-Feature `ollama`, `egress` im Wertebereich von `command`, Messstand `amd64` als Einzelmessung; Release-Zuordnung des Projektinhabers). Der Release trägt alles, was seit dem Tag `v0.11.0` auf `main` liegt: die Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) (neun Slices, Closure geschlossen), die drei Slices ohne Welle [`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md), [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md) und [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) sowie die beiden Hygiene-Slices [`slice-v1-cli-dry-run-log-would`](../done/slice-v1-cli-dry-run-log-would.md) und [`slice-v1-cli-help-handbuch-link`](../done/slice-v1-cli-help-handbuch-link.md) (und [`slice-v1-harness-gofmt-und-depguard`](../done/slice-v1-harness-gofmt-und-depguard.md), ohne Produktänderung); Slices ohne Produktänderung (Harness, Planung) ändern den Inhalt nicht.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run),
[`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle),
[`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe),
[`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md),
[`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md),
[`ADR-0007`](../../adr/0007-distributionswege-ghcr.md),
[`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md),
[`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; Lastenheft 0.17.0 liefern die Slices, dieser Slice ändert `spec/` nicht).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.12.0` liefert als Minor-Release (neue Funktionen seit `v0.11.0`, SemVer) den Block `[Unreleased]` des CHANGELOG — den Befehl `u-boot egress list`, den Wert `egress` im Feld `command` der `--json`-Ausgabe, die Klärung von Subdomain und Wildcard (Exit `10` mit Hinweis), das Katalog-Feature `ollama` (nur CPU, nur Sandbox-Profil) mit dem Messstand `amd64` als Einzelmessung, die Handbuch-Links in `--help` und die Bedingungsform der `--dry-run`-Logzeilen — als GHCR-Image, Binaries, Pakete und Homebrew-Formel. Dieser Slice bereitet vor und lässt den Tag zu; Push und Tags führt nur der Hauptagent nach Freigabe des Projektinhabers.

### Inhalt des Release (gelesen am 2026-10-19, `CHANGELOG.md` Block `[Unreleased]`)

| Eintrag im CHANGELOG | Art | Hinweis für den Lead |
|---|---|---|
| Devcontainer-Feature `ollama` (Ollama 0.40.6 aus gepinntem, geprüftem Archiv; nur CPU; nur mit Sandbox-Profil; Volume für Modelle; zwei Hosts in der Default-Allowlist) samt Messstand (Colima/arm64 und zwei Läufe auf einem nativen Linux-x86_64-Host, kein CI-Lauf) | Added | **Cloud-Modelle nicht enthalten** (Auftrag aus der Welle-Closure, Handbuch-Satz L-1); der zweite Host des Modell-Pull ist nicht auf Stabilität gemessen |
| `u-boot --help` verlinkt das Benutzerhandbuch (Tag der Version, aktuelles auf `main`) | Added | Dev-Builds nur `main` |
| `u-boot egress list` (Herkunft je Host, offline, `--json`, `--dry-run`/`--diff` mit Exit `2` abgewiesen, Exit `10` bei verweigerter `origin`-URL) | Added | Wert `egress` im Feld `command` der `--json`-Ausgabe; Lastenheft 0.17.0 trägt die beiden vom Projektinhaber bestätigten Wortlaute (Release-Gates 1 und 2, 2026-10-19) |
| Wildcard-Eintrag in `egress.allow` weiter Exit `10`, Fehlertext nennt den einzutragenden Host; Handbuch trägt das gemessene Subdomain-Verhalten | Changed | |
| `--dry-run`: INFO-Zeilen von `config set` und `generate` im Konditional | Changed | Skripte, die die alte Zeile auf stderr abglichen, nutzen `--json` |
| Fixed (vier Einträge zum Handbuch: Neuaufbau statt Neustart, Subdomain-Aussage getrennt in „Test“ und „Messung“, zwei Grenzen der Restriktion ohne `sudo`, `amd64` gemessen mit zwei Grenzen) | Fixed | Messstände stehen im Eintrag |
| Handbuch: Cloud-Modelle von Ollama (L-1) | Fixed | neu in diesem Cut, ein Satz im Ist-Zustand |

### Aufträge an den Release-Cut (gelesen in [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) §7, *Auftrag an den Release-Cut*, Punkte 1 bis 6)

1. **Gates 1 und 2:** erfüllt — der Projektinhaber hat die beiden Wortlaute am 2026-10-19 im Gespräch bestätigt (Wortlaut unverändert, Lastenheft `0.17.0`); der datierte Vermerk steht in der Welle-Datei (Commit `332b171`). Keine Änderung am Lastenheft.
2. **CHANGELOG zuordnen:** `[Unreleased]` → `[0.12.0]`; Lead mit „Cloud-Modelle nicht enthalten“ und dem Weg „Liste ändern = Neuaufbau, kein Neu-Laden“ (Entscheidung R0); Messstand Colima/arm64 und Linux-`amd64` als Einzelmessung.
3. **Handbuch-Satz zu Cloud-Modellen (L-1):** ziehen gelingt, Lauf nicht; Hosts `ollama.com` und Anmeldung fehlen.
4. **Der rote Test** `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` ist zu klären, bevor der Cut sagt, der Docker-Lauf sei grün (§3, Umsetzung, Befund).
5. **CI-Hinweise an den Projektinhaber** (kein Eingriff, `AGENTS.md` §3.6; §3, Umsetzung).
6. **Register:** `BEO-E2E/ci-ueberspringt-nested-podman-build` fortschreiben (`state.md`).

**Entscheidung des Projektinhabers (2026-10-19, nur hier festgehalten):** Die Historie wird nicht umgeschrieben („ich habe jetzt keine Geheimnisse gesehen – wir lassens so“). Gemeint ist die Git-Historie nach dem Befund [`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md); der Cut ändert weder Commits noch Rohdaten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tags und jede Aktion mit Außenwirkung durch den Implementer.** Push von `main`, jeder Tag, `gh release` und jedes Löschen sind extern wirksam und lassen sich nicht sauber zurücknehmen ([`releasing.md`](../../../maintainer/releasing.md) Zweck). **Push und Tags führt nur der Hauptagent, nach der ausdrücklichen Freigabe des Projektinhabers, je Vorgang**; der Implementer bereitet vor und führt keinen der Schritte M1 bis M8 aus.
- **Eine Änderung der Pipeline** (Workflows, Tap-Skripte, `-v` im Aufruf von `make test-docker`, Required Checks): Gate-Änderung nach [`AGENTS.md`](../../../../AGENTS.md) §3.6; der Cut meldet die CI-Befunde an den Projektinhaber.
- **Eine Produktänderung zum roten Test:** Der Cut untersucht die Ursache nur lesend und mit Messläufen außerhalb des Repos; eine Korrektur ist ein eigener Slice.
- **Die Folge-Slices der Welle** ([`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../open/slice-v1-egress-hinweise-neuaufbau-und-messstand.md) für L-2/L-3 u. a.): sie folgen dem Cut, die ausgelieferten Meldungen tragen den Stand der Closure.
- **Änderung des Lastenhefts oder der Spezifikation:** Die Version `0.17.0` bleibt; die Bestätigung der Gates ändert keinen Wortlaut.
- **Umschreiben der Historie:** Entscheidung des Projektinhabers (oben).

## 2. Definition of Done

- [ ] `v0.12.0` ist getaggt, der `publish`-Lauf ist grün (GHCR `:0.12.0` und `:latest`, sechs Binaries, vier Pakete, `SHA256SUMS`, `u-boot.rb`, Paket-Smoke, Job `tap`); ein Vorab-Tag ist optional (Entscheidung des Projektinhabers) und, falls gesetzt, aufgeräumt oder als offen vermerkt.
- [ ] `make gates`, `make ci` und `make fullbuild` sind vor dem Tag grün (oder der Grund für ein nicht lauffähiges Target steht in §3), der Trockenlauf des Tap-Pfads (M-T) ist ausgeführt, der letzte CI-Lauf auf dem Push-Tip ist grün; **der Befund zum roten Test steht fest** und die Aussage des Cuts zum Docker-Lauf ist ihm gleich.
- [ ] Doku-Update: CHANGELOG, READMEs, Handbuch und Roadmap tragen `v0.12.0`; die Release-Aussagen nennen den Messstand und „Cloud-Modelle nicht enthalten“.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Jede Aussage in dieser Tabelle ist an der genannten Stelle gelesen (Stand 2026-10-19, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Zeilennummern nennen den Fundort beim Schreiben; der Implementer sucht sie vor der Änderung neu. Vor dem Handoff (nach der letzten Fix-Runde) gleicht er die Tabelle gegen `git diff --stat` ab; eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage, nicht durch eine Suche.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update | `:12` `## [Unreleased]` wird zum leeren Anker; der Inhalt wandert unter `## [0.12.0] - 2026-10-19` mit Lead-Absatz im Stil der Vorgänger („Eleventh release.“ steht unter `[0.11.0]`; hier „Twelfth release.“). Der Lead nennt Egress-Anzeige, `egress` im JSON-Enum, Wildcard-Klärung, `ollama` (nur CPU, nur Sandbox-Profil), „Cloud-Modelle nicht enthalten“, „Neuaufbau statt Neu-Laden“, Lastenheft 0.17.0 und den Messstand. Im Block `### Fixed` kommt ein Eintrag zum Handbuch-Satz über Cloud-Modelle (L-1) hinzu. Compare-Links am Fuß (`:1384` `[Unreleased]: …/compare/v0.11.0...HEAD` → `v0.12.0...HEAD`; neue Zeile `[0.12.0]: …/compare/v0.11.0...v0.12.0` vor `:1385`). |
| `cmd/uboot/main.go` | update | `:43` `var version = "0.11.0-dev"` → `0.12.0-dev`. |
| `repo.mk` | update | `:51` und `:54` (Kommentare) und `:55` `VERSION ?= 0.11.0-dev` → `0.12.0-dev`. |
| `Dockerfile` | update | `:35` und `:37` (Kommentar) und `:39` `ARG UBOOT_VERSION=0.11.0-dev` → `0.12.0-dev`. |
| `README.md` | update | Status-Block `:10`–`:14` (als Stand formuliert, ohne Datum und ohne „released“; ohne Kennungen): Version `v0.12.0`, `egress list`, `ollama` mit Messstand und „Cloud-Modelle nicht enthalten“; `:99` `VERSION=0.11.0` → `0.12.0`; neue Zeile in der Releases-Tabelle nach der `v0.11.0`-Zeile (`:200`). Der Satz zu `podman build` und `amd64` im Status-Block bleibt (gelesen: Handbuch `:1695` „`podman build` auf `amd64` und CI sind nicht gemessen“). Die Feature-Zahl „14 features“ bleibt (der Katalog ändert sich in diesem Cut nicht). |
| `README.de.md` | update | gleiches Muster: `:11`–`:15` Status, `:101` `VERSION=0.11.0`, Tabellenzeile nach `:203`. |
| `docs/user/benutzerhandbuch.md` | update | `:4` `Software-Version: v0.11.0`, `:5` `Stand: 2026-10-18` (→ Tag-Datum), `:80` Beispielausgabe `u-boot version 0.11.0`, `:106` `download/v0.11.0`, `:133` `VERSION=0.11.0`, `:2371` „beschreibt `u-boot` v0.11.0“; im Abschnitt *Ollama (optional)* ein Punkt zu Cloud-Modellen nach dem Punkt *Egress* (L-1; Ist-Zustand, Messstand des Validators, kennungsfrei). |
| `docs/maintainer/releasing.md` | update | `:57` nennt als „zuletzt“ den Cut `v0.10.0`; berichtigt auf den letzten abgeschlossenen Cut `v0.11.0`. Die übrigen Aussagen (§1 bis §8) stimmen mit dem Ist-Stand überein (gelesen). |
| `docs/plan/planning/in-progress/roadmap.md` | update | Beim Start: `Nichts in Arbeit.` (Ruhe-Marker, `.d-check.yml` `planning`) entfällt, stattdessen ein Zeiger auf diesen Slice im Block „Offene Wellen“; die Zeile der Welle `welle-sandbox-egress-ollama` unter „Abgeschlossene Wellen“ nennt „Release-Zuordnung und zwei Release-Gates offen“ — berichtigt auf den Ist-Stand (Gates bestätigt, 2026-10-19; Release `v0.12.0`); nach dem Tag: Marker zurück, neue Zeile `v0.12.0` nach `v0.11.0` mit Datum und Verweis auf diesen Slice, Tag-Commit-Hash erst nach dem Tag; `**Letzte Änderung:**` im Kopf. Die Zeile von `welle-ollama-cloud-modelle` unter „Nächste Wellen“ trägt den Trigger „Release der Welle … geschnitten“; er tritt mit dem Tag ein, die Roadmap nennt es dann (Planner). |
| `docs/plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/state.md` | update | Register: Trigger eingetreten (ohne `-v` kein Beleg), Befunde zum Workflow und zum Test (Auftrag 6). |
| `docs/plan/planning/open/README.md` | update | Verweis-Nachzug durch `make slice-mv` (der Slice steht dort nicht als Bestandszeile; gelesen). |
| `cmd/uboot/package_asset_name_test.go`, `scripts/package-asset-name.sh`, `internal/adapter/driving/cli/cli_test.go` | kein Diff (gelesen) | `cli_test.go` `:380`–`:391`: die Versionen `0.11.0`, `0.12.0-rc.1` usw. sind Eingaben der Tabelle für die Handbuch-Links, kein Versionsstand. |
| `.github/workflows/*.yml`, `scripts/`, `packaging/`, `spec/`, `docs/plan/adr/` | kein Diff (gelesen) | Workflow-Änderungen sind Gate-Änderungen (§1, nicht Gegenstand); `spec/lastenheft.md` `:9` `Version 0.17.0` bleibt. |
| `docs/reviews/<Datum>-slice-v1-release-cut-v0.12.0.md`, `docs/reviews/<Datum>-verification-slice-v1-release-cut-v0.12.0.md` | neu, durch Reviewer und Verifier | Review- und Verification-Report des Cuts (nicht im Implementer-Diff). Kennungen in Tabellenzellen ohne Backticks im Linktext. |

**Aussage zum Messstand — Soll, in CHANGELOG, Handbuch und Release-Zeile der READMEs gleich:** Gemessen mit Docker 29.5.2 unter Colima 0.10.3 (arm64) und, für das Feature `ollama`, in zwei Läufen auf einem nativen Linux-x86_64-Host (Docker Engine 29.9.1, 2026-10-19; kein CI-Lauf). **Nicht gemessen:** Docker Desktop, Linux hinter einer VM, größere Modelle, die Stabilität des zweiten Hosts des Modell-Pull, `podman build` auf `amd64`.

**Reihenfolge der Maintainer-Aktionen.** Der Implementer bereitet vor (Commits lokal, Gates lokal) und **führt keinen der Schritte M1 bis M8 aus**. **Push und Tags führt nur der Hauptagent nach Freigabe des Projektinhabers, je Vorgang.** Der Slice bleibt bis zum Tag in `in-progress/`.

| # | Aktion | Wer | Prüfung / Erwartung |
|---|---|---|---|
| M0 | `date +%F`; Datum in CHANGELOG, READMEs (Tabellenzeile), Handbuch (`Stand`) gegen das Tag-Datum halten | Implementer, vor jedem Push | Das Datum steht an mehreren Stellen fest (`releasing.md` §5) |
| M-T | Trockenlauf des Tap-Pfads (`releasing.md` §4), im Docker, `--network none`, ohne Schreiben ins Tap | Implementer; Ergebnis hier | Erwartet: Formel aus Fixture, `tap-formula-check.sh` Exit `0` mit `v0.12.0` und Exit `1` mit `v0.11.0`, `tap-nachzug.sh` mit Stellvertretern. **Beweist nicht:** `is_stable`-Bedingung, echter `curl`, echter `git push`, Token |
| M1 | `git push origin main` | Hauptagent, nach Freigabe des Projektinhabers | Review-Report und Cut-Commit liegen vor; `make docs-check` lokal grün |
| M2 | `gh run list --limit 3`: `ci` und `integration-docker` auf dem Push-Tip | Maintainer | Roter `image-scan`: Go-Pin heben (`releasing.md` §4); **roter `integration-docker` wegen `NestedPodman_BuildEgress`: bekannter Befund (§3, Umsetzung), kein Pflicht-Check** |
| M3 | Optional: `git tag v0.12.0-rc.1 && git push origin v0.12.0-rc.1` | Hauptagent, nach Freigabe des Projektinhabers | Vorab-Tag probt `is_stable == true` und den Job `tap` nicht |
| M4 | `gh workflow run tap-check` (Token darf auf das Tap schreiben) | Maintainer | `releasing.md` §2, §4 |
| M5 | `git switch main && git pull --ff-only`, `git tag v0.12.0` auf dem in M2 geprüften Commit, `git push origin v0.12.0` | Hauptagent, nach Freigabe des Projektinhabers | Nach dem Push wird der Tag nicht verschoben (`releasing.md` §8) |
| M6 | `releasing.md` §6: `gh run watch`; `gh release view v0.12.0 --json assets`; `docker run --rm ghcr.io/pt9912/u-boot:0.12.0 --version`; `:latest` gleicher Digest; Job `tap` wirklich gelaufen | Maintainer | zwölf Assets, kein Prerelease |
| M7 | Auf einem Mac: `brew tap pt9912/u-boot`, `brew install u-boot`, `u-boot --version` = `0.12.0` | Projektinhaber | `releasing.md` §6 |
| M8 | Nacharbeit (`releasing.md` §7): Roadmap-Zeile mit Tag-Hash und Datum, Closure-Notiz, Verification Evidence, `git mv` nach `done/` | Planner (Closure), Implementer (Roadmap) | Closure erst nach M6 |

## 4. Trigger

**Start** (`next` → `in-progress`): Die Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) und die drei Slices ohne Welle ([`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md), [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md), [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md)) liegen in `done/`; die Release-Gates 1 und 2 sind bestätigt (Vermerk 2026-10-19, `332b171`); `in-progress/` enthält nur die Master-Dokumente (WIP-Limit 1). Der Weg ist `open/` → `in-progress/` per reinem `git mv` (`make slice-mv SLICE=slice-v1-release-cut-v0.12.0 TO=in-progress`); der Roadmap-Marker geht in einen eigenen Commit direkt nach dem Move, beide im selben Push.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tags.
- `in-progress` → `open` (blockiert): ein roter `image-scan` (Go-Pin heben), ein Befund des Trockenlaufs, der eine Änderung an Workflow oder Skript verlangt (eigener Vorgang davor), ein fehlschlagender `publish`-Lauf, dessen Ursache einen Patch-Release braucht, oder ein Befund des Reviewers, der eine Produktänderung verlangt.

## 5. Closure-Trigger

DoD vollständig: stabiler Tag gesetzt (M5), `publish`- und Tap-Lauf kontrolliert (M6), Trockenlauf M-T mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs und Steering-Loop-Eintrag, Verification Evidence (Verifier, getrennt von Implementer und Reviewer). Der Slice geht **erst nach dem Tag** nach `done/`; der Ruhe-Marker der Roadmap steht erst dann wieder.

## 6. Risiken und offene Punkte

Der Ausgang jedes Risikos wird bei der Closure festgeschrieben.

- **Der Job `tap` scheitert im stabilen Lauf, obwohl der Trockenlauf grün war** (der Trockenlauf probt die Skripte, nicht die Workflow-Bedingung, den echten Download, den Push und das Token). — **Ausgang:** offen bis zur Closure.
- **Der CI-Lauf `integration-docker` ist rot** (der Test `NestedPodman_BuildEgress` ist bekannt rot; ob er auf dem Runner läuft oder überspringt, ist Gegenstand der Untersuchung, §3 Umsetzung). — **Ausgang:** offen bis zur Closure.
- **Die Release-Aussage zu Cloud-Modellen fehlt oder ist zu weit.** — **Ausgang:** offen bis zur Closure (Review und Verifier prüfen Lead, Handbuch und Messprotokoll).
- **Das Tag-Datum weicht vom festen Datum ab** (CHANGELOG, Tabellenzeilen der READMEs, Handbuch-`Stand`). — **Ausgang:** offen bis zur Closure.
- **Das Messprotokoll der Welle trägt Rohdaten mit Zugangsdaten** ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)); der Projektinhaber lässt die Historie unberührt. — **Ausgang:** entfallen: Entscheidung des Projektinhabers vom 2026-10-19 (§1); der Cut ändert nichts daran.
- **GHCR-Tags und Release früherer Vorab-Versionen** (`0.10.0-rc.1`, `0.10.0-rc.2`, `0.11.0-rc.1`) stehen noch; der Cut berührt sie nicht. — **Ausgang:** offen bis zur Closure (beim Projektinhaber).

## 7. Closure-Notiz

Wird vom Planner bei der Closure gefüllt (nach dem Tag).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, `docs/maintainer/`, `docs/reviews/` und Release-Dateien (`repo.mk`, `Dockerfile`) — Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006). `cmd/uboot` ist laut Modus-Deklaration keine Sub-Area (Aspirantin); berührt wird dort eine Zeile (`var version`), keine Konvention.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen ([`../observations/`](../observations/README.md)): [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) ist `offen` und betrifft diesen Cut direkt (Auftrag 6); `BEO-DOK/plan-aussage-ueber-code` und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` sind `verkörpert` (`AGENTS.md` §3.10: Aussagen über Code mit Fundstelle, Bedingungen in Handbuch mit Test oder gestrichen — der neue Handbuch-Satz zu Cloud-Modellen nennt nur Gemessenes).

Alle berührten Sub-Areas GF.
