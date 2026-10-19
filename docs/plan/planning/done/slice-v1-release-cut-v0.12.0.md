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
- **Die Folge-Slices der Welle** ([`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../done/slice-v1-egress-hinweise-neuaufbau-und-messstand.md) für L-2/L-3 u. a.): sie folgen dem Cut, die ausgelieferten Meldungen tragen den Stand der Closure.
- **Änderung des Lastenhefts oder der Spezifikation:** Die Version `0.17.0` bleibt; die Bestätigung der Gates ändert keinen Wortlaut.
- **Umschreiben der Historie:** Entscheidung des Projektinhabers (oben).

## 2. Definition of Done

- [x] `v0.12.0` ist getaggt (leichter Tag auf `9c28c5c`, 2026-10-19), der `publish`-Lauf ist grün (Lauf 39224848402: `publish` mit GHCR-Push, `tap`, Paket-Smoke rpm/fedora und deb/ubuntu `success`; zwölf Assets, `:0.12.0` und `:latest` gleicher Digest); ein Vorab-Tag wurde nicht gesetzt (Entscheidung des Projektinhabers: „direkt v0.12.0 ohne rc“). Beleg: §7.
- [x] `make gates`, `make ci` und `make fullbuild` waren vor dem Tag grün (§3, Sensors; vom Verifier unabhängig gefahren, `aeb11d6`), der Trockenlauf des Tap-Pfads (M-T) ist ausgeführt (§3), der letzte CI-Lauf auf dem Push-Tip `9c28c5c` ist grün (`ci` 39223699156, `integration-docker` 39223699209); der Befund zum roten Test steht fest (§3) und die Aussage des Cuts zum Docker-Lauf ist ihm gleich: grün auf dem Runner heißt nicht, dass der Test dort lief (§7).
- [x] Doku-Update: CHANGELOG, READMEs, Handbuch (`fb4641c`, Fix-Runde `1015027`, `9c28c5c`) und Roadmap (`fddf45f`) tragen `v0.12.0`; die Release-Aussagen nennen den Messstand und „Cloud-Modelle nicht enthalten“.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`6536c07`, [`2026-10-19-slice-v1-release-cut-v0.12.0`](../../../reviews/2026-10-19-slice-v1-release-cut-v0.12.0.md)); Verification-Report `aeb11d6` ([`2026-10-19-verification-slice-v1-release-cut-v0.12.0`](../../../reviews/2026-10-19-verification-slice-v1-release-cut-v0.12.0.md)) — kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Beobachtungs-Register fortgeschrieben (§7, `cbd00d7` und Closure-Commit).
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (§7).

## 3. Plan (vor Code)

Jede Aussage in dieser Tabelle ist an der genannten Stelle gelesen (Stand 2026-10-19, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Zeilennummern nennen den Fundort beim Schreiben; der Implementer sucht sie vor der Änderung neu. Vor dem Handoff (nach der letzten Fix-Runde) gleicht er die Tabelle gegen `git diff --stat` ab; eine Datei, die er als „unverändert“ führt, belegt er durch eine gelesene Aussage, nicht durch eine Suche.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update | `:12` `## [Unreleased]` wird zum leeren Anker; der Inhalt wandert unter `## [0.12.0] - 2026-10-19` mit Lead-Absatz im Stil der Vorgänger („Eleventh release.“ steht unter `[0.11.0]`; hier „Twelfth release.“). Der Lead nennt Egress-Anzeige, `egress` im JSON-Enum, Wildcard-Klärung, `ollama` (nur CPU, nur Sandbox-Profil), „Cloud-Modelle nicht enthalten“, „Neuaufbau statt Neu-Laden“, Lastenheft 0.17.0 und den Messstand. Im Block `### Fixed` kommt ein Eintrag zum Handbuch-Satz über Cloud-Modelle (L-1) hinzu. Compare-Links am Fuß (`:1384` `[Unreleased]: …/compare/v0.11.0...HEAD` → `v0.12.0...HEAD`; neue Zeile `[0.12.0]: …/compare/v0.11.0...v0.12.0` vor `:1385`). |
| `cmd/uboot/main.go` | update | `:43` `var version = "0.11.0-dev"` → `0.12.0-dev`. |
| `repo.mk` | update | `:51` und `:54` (Kommentare) und `:55` `VERSION ?= 0.11.0-dev` → `0.12.0-dev`. |
| `Dockerfile` | update | `:35` und `:37` (Kommentar) und `:39` `ARG UBOOT_VERSION=0.11.0-dev` → `0.12.0-dev`. |
| `README.md` | update | Status-Block `:10`–`:14` (als Stand formuliert, ohne Datum und ohne „released“; ohne Kennungen): Version `v0.12.0`, `egress list`, `ollama` mit Messstand und „Cloud-Modelle nicht enthalten“; `:99` `VERSION=0.11.0` → `0.12.0`; neue Zeile in der Releases-Tabelle nach der `v0.11.0`-Zeile (`:200`). Der Satz zu `podman build` und `amd64` im Status-Block nennt nach der Fix-Runde die gemessene Grenze (Namensauflösung hinter der Beschränkung; §3 Umsetzung, Befund zum roten Test). Die Feature-Zahl „14 features“ bleibt (der Katalog ändert sich in diesem Cut nicht). |
| `README.de.md` | update | gleiches Muster: `:11`–`:15` Status, `:101` `VERSION=0.11.0`, Tabellenzeile nach `:203`. |
| `docs/user/benutzerhandbuch.md` | update | `:4` `Software-Version: v0.11.0`, `:5` `Stand: 2026-10-18` (→ Tag-Datum), `:80` Beispielausgabe `u-boot version 0.11.0`, `:106` `download/v0.11.0`, `:133` `VERSION=0.11.0`, `:2371` „beschreibt `u-boot` v0.11.0“; im Abschnitt *Ollama (optional)* ein Punkt zu Cloud-Modellen nach dem Punkt *Egress* (L-1; Ist-Zustand, Messstand des Validators, kennungsfrei). |
| `docs/maintainer/releasing.md` | update | `:57` nennt als „zuletzt“ den Cut `v0.10.0`; berichtigt auf den letzten abgeschlossenen Cut `v0.11.0`. Die übrigen Aussagen (§1 bis §8) stimmen mit dem Ist-Stand überein (gelesen). |
| `docs/plan/planning/in-progress/roadmap.md` | update (Cut `cc58225`, Closure `fddf45f`) | Beim Start: `Nichts in Arbeit.` (Ruhe-Marker, `.d-check.yml` `planning`) entfällt, stattdessen ein Zeiger auf diesen Slice im Block „Offene Wellen“; die Zeile der Welle `welle-sandbox-egress-ollama` unter „Abgeschlossene Wellen“ nennt „Release-Zuordnung und zwei Release-Gates offen“ — berichtigt auf den Ist-Stand (Gates bestätigt, 2026-10-19; Release `v0.12.0`); nach dem Tag: Marker zurück, neue Zeile `v0.12.0` nach `v0.11.0` mit Datum und Verweis auf diesen Slice, Tag-Commit-Hash erst nach dem Tag; `**Letzte Änderung:**` im Kopf. Die Zeile von `welle-ollama-cloud-modelle` unter „Nächste Wellen“ trägt den Trigger „Release der Welle … geschnitten“; er tritt mit dem Tag ein, die Roadmap nennt es dann (Planner). |
| `docs/plan/planning/observations/BEO-E2E/ci-ueberspringt-nested-podman-build/state.md` | update | Register: Trigger eingetreten (ohne `-v` kein Beleg), Befunde zum Workflow und zum Test (Auftrag 6). |
| `docs/plan/planning/open/slice-v1-egress-dns-upstream-schliessen.md` | update (Fix-Runde) | Anker für die Korrektur des roten Tests: Beleg im Gelesenen Stand, Abnahme in der DoD, Test in der Plan-Tabelle (Review F-3). |
| `docs/plan/planning/open/README.md` | kein Diff (gelesen) | Der Slice steht dort nicht als Bestandszeile; `make slice-mv` zog nichts nach (0 eingehende Verweise). |
| `cmd/uboot/package_asset_name_test.go`, `scripts/package-asset-name.sh`, `internal/adapter/driving/cli/cli_test.go` | kein Diff (gelesen) | `cli_test.go` `:380`–`:391`: die Versionen `0.11.0`, `0.12.0-rc.1` usw. sind Eingaben der Tabelle für die Handbuch-Links, kein Versionsstand. |
| `.github/workflows/*.yml`, `scripts/`, `packaging/`, `spec/`, `docs/plan/adr/` | kein Diff (gelesen) | Workflow-Änderungen sind Gate-Änderungen (§1, nicht Gegenstand); `spec/lastenheft.md` `:9` `Version 0.17.0` bleibt. |
| `docs/plan/planning/done/slice-v1-release-cut-v0.12.0.md` (jetzt `done/`) | neu, update | Dieser Slice: Plan (`ed4d51f`), Move nach `in-progress/` (`f1d0c94`, `e85c084`), Umsetzung im Plan (`ee577e3`), Fix-Runde, Closure und Move nach `done/` (Planner). Die Zeile fehlte im Stand vor der Closure (Verifier V-3). |
| `docs/plan/planning/done/welle-sandbox-egress-ollama.md`, `slice-cr-lastenheft-zielbild.md`, `slice-cr-sandbox-egress-anzeige.md`, `slice-cr-cli-json-command-enum.md` | update (datierte Vermerke) | Vermerk der Gates 1 und 2 vom 2026-10-19 in vier Dateien aus `332b171`, an bestehende Absätze angehängt, kein Ursprungstext gelöscht (Verifier, Diff `ecde149..HEAD -- docs/plan/planning/done`); Auftrag in §1 *Aufträge*, Punkt 1. |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/` (`state.md`, `evidence/slice-v1-release-cut-v0.12.0.md`) | update, neu (Closure) | Vorgang vierundzwanzig der Beobachtung: Befund V-3 des Verifiers (diese Tabelle). |
| `docs/reviews/2026-10-19-slice-v1-release-cut-v0.12.0.md`, `docs/reviews/2026-10-19-verification-slice-v1-release-cut-v0.12.0.md` | neu, durch Reviewer (`6536c07`) und Verifier (`aeb11d6`, `9c28c5c`) | Review- und Verification-Report des Cuts (nicht im Implementer-Diff). Kennungen in Tabellenzellen ohne Backticks im Linktext. |

**Aussage zum Messstand — Soll, in CHANGELOG, Handbuch und Release-Zeile der READMEs gleich:** Gemessen mit Docker 29.5.2 unter Colima 0.10.3 (arm64) und, für das Feature `ollama`, in zwei Läufen auf einem nativen Linux-x86_64-Host (Docker Engine 29.9.1, 2026-10-19; kein CI-Lauf). **Nicht gemessen:** Docker Desktop, Linux hinter einer VM, größere Modelle, die Stabilität des zweiten Hosts des Modell-Pull, `podman build` auf `amd64` hinter der Egress-Beschränkung scheitert an der Namensauflösung (zwei Läufe, ein Host; Befund zum roten Test).

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

### Umsetzung (Implementer, 2026-10-19; nichts gepusht, kein Tag)

Commits dieses Cuts (lokal): `332b171` (Vermerk der Gates 1 und 2 in der Welle-Datei und in den drei CR-Slices), `ed4d51f` (Plan, `open/`), `f1d0c94` und `e85c084` (`make slice-mv` nach `in-progress/`), `cc58225` (Roadmap: Zeiger, Zeile der Welle), `fb4641c` (CHANGELOG, Versionsstrings, READMEs, Handbuch, `releasing.md`), `51bf62e` (Register).

**Gates 1 und 2.** Der Vermerk (Wortlaut: bestätigt am 2026-10-19 durch den Projektinhaber im Gespräch; Wortlaut unverändert) steht in [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) bei Gate 1, Gate 2, §7 Punkt 8 und im *Auftrag an den Release-Cut*, ferner im Zustandsmarker „offen“ von [`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md) (§1 Absatz „Das Streichen ist keine Freigabe“ und §7 „Stand der Gates 1 und 2“; die Frage E-3 bleibt unverändert, Chronik) sowie im „Stand des Gates“ von [`slice-cr-sandbox-egress-anzeige`](../done/slice-cr-sandbox-egress-anzeige.md) und [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) (die Welle-Datei verlangt den Vermerk „in dieser Datei und in den beiden CR-Slices“). Die Zeile der Welle in der Roadmap trägt den Ist-Stand.

**Sensoren (Docker, Host `linux/amd64`, Docker Engine 29.9.1; Exit-Codes der `make`-Prozesse):**

| Sensor | Ergebnis | Laufzeit |
|---|---|---|
| `make gates` | Exit 0 (`baseline-verify` v6.18.1 OK, 54 Dateien; `d-check` 500 Dateien, 0 Befunde; Coverage 91,40 % gegen 90 %) | 38,1 s |
| `make docs-check` / `make doc-planning` | je 500 Dateien, 0 Befunde | — |
| `make doc-trace` | 83 Anforderungen, 0 Waisen | — |
| `make ci` | Exit 0 (`govulncheck`: No vulnerabilities found; `image-scan` HIGH/CRITICAL: grün; Trivy meldet nur, Version 0.76.0 sei verfügbar, gepinnt ist 0.70.0) | 1 min 43,7 s |
| `make fullbuild` | Exit 0 (`ci` + Runtime-Image); `docker run --rm u-boot:latest --version` = `u-boot version 0.12.0-dev`, OCI-Label `0.12.0-dev` | 1 min 01,8 s |
| `make build VERSION=0.12.0` | Exit 0; `--version` = `u-boot version 0.12.0`, OCI-Label `0.12.0` (lokales Image, nicht gepusht) | — |
| `make doc-immutable RANGE="v0.11.0..HEAD"` | 253 Commits, 0 Befunde (Accepted-ADRs unverändert) | — |
| `make test-docker` | **Exit 2**: ein roter Test, `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (26,94 s); alle übrigen Pakete `ok` (`internal/e2e` FAIL nach 594,958 s, wegen dieses einen Tests); der Aufruf läuft ohne `-v` | 10 min 05 s |
| `make verify-depguard` | **nicht gelaufen**: keine Importregel berührt, kein Produktcode im Diff | — |

**Plan-Tabelle gegen `git diff --stat ed4d51f..HEAD`:** elf Einträge (`CHANGELOG.md`, `Dockerfile`, `README.de.md`, `README.md`, `cmd/uboot/main.go`, `docs/maintainer/releasing.md`, `roadmap.md`, diese Datei, `state.md` von `BEO-E2E/ci-ueberspringt-nested-podman-build`, `docs/user/benutzerhandbuch.md`, `repo.mk`); jeder hat eine Zeile in der Tabelle. `open/README.md` steht nicht im Diff (der Slice hatte dort nie eine Bestandszeile; gelesen: keine Zeile für `slice-v1-release-cut-v0.12.0`). Gelesen als „kein Diff“: `cli_test.go` `:380`–`:391` (Eingaben der Link-Tabelle), `.github/workflows/`, `scripts/`, `packaging/`, `spec/lastenheft.md` `:9` (`Version 0.17.0`). Nicht in diesem Abgleich stehen die vier Dateien in `done/` aus `332b171` (Vermerke vor dem Plan-Commit; ihr Auftrag steht in §1 unter *Aufträge an den Release-Cut*, Punkt 1).

**Fix-Runde zum Review ([`2026-10-19-slice-v1-release-cut-v0.12.0.md`](../../../reviews/2026-10-19-slice-v1-release-cut-v0.12.0.md), 0 HIGH, 0 MEDIUM, 3 LOW).** F-1: Lead, Fixed-Eintrag und Handbuch sagen gleich, `ollama signin` ohne Konto habe `account unavailable` gemeldet und eine Anmeldung mit Konto sei nicht geprüft (Beleg: Validierungsbericht Zeile e). F-2: CHANGELOG-Lead, Fixed-Eintrag, beide READMEs (Status-Block, Zeile `v0.12.0`) und das Handbuch (Abschnitt *Grenze der Namensauflösung*, Grenzen, Gültigkeitsbereich) nennen die Messung zum verschachtelten `podman build` auf `amd64` (Schritt 3 scheitert, `example.com` löst auf, Verbindung gesperrt, mit `--dns 9.9.9.9` grün); nicht gemessen bleibt Colima/arm64 und die CI. F-3: Anker ist [`slice-v1-egress-dns-upstream-schliessen`](../open/slice-v1-egress-dns-upstream-schliessen.md) (Ursache gelesen: `egress-init.sh.tmpl` Zeilen 63 bis 66 ohne Bedingung an Prozess oder Benutzer, Gegenstand und Messfragen D1/D2 treffen den verschachtelten Pfad); der Slice nennt den Test im Gelesenen Stand, in der Plan-Tabelle und als Abnahme in der DoD; kein neuer Slice. **Lücke:** Ob eine Regel nach Benutzer den Pfad des verschachtelten Podman schließt (die `pasta`-Prozesse laufen als `vscode`), ist ungemessen; D2/D3 des Folge-Slice entscheiden, ob der Test grün wird oder angepasst wird.

**Plan-Tabelle gegen `git diff --stat ed4d51f..HEAD` nach der Fix-Runde (Stand der Fix-Runde, vor deren Commit):** die zwölf Dateien des Stands davor (die elf oben plus `docs/reviews/2026-10-19-slice-v1-release-cut-v0.12.0.md`, als neu durch den Reviewer geführt) tragen je eine Zeile; hinzu kommt `docs/plan/planning/open/slice-v1-egress-dns-upstream-schliessen.md` mit eigener Zeile. Berührt hat die Fix-Runde nur `CHANGELOG.md`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, diese Datei und den offenen Slice; `internal/`, `spec/`, `docs/plan/adr/` und `.github/` stehen weiter nicht im Diff.

**Plan-Tabelle gegen `git diff --stat ecde149..HEAD` (Closure, Stand `9c28c5c` vor den Closure-Commits).** 18 Dateien im Diff: `CHANGELOG.md`, `Dockerfile`, `README.de.md`, `README.md`, `cmd/uboot/main.go`, `docs/maintainer/releasing.md`, `repo.mk`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, dieser Slice, `state.md` von `BEO-E2E/ci-ueberspringt-nested-podman-build`, `slice-v1-egress-dns-upstream-schliessen.md` (`open/`), die beiden Reports und die vier Dateien in `done/`; jede trägt eine Zeile der Tabelle (Verifier V-3 ergänzt: Slice selbst, vier `done/`-Vermerke, Reports mit Datum). Die Closure berührt zusätzlich `roadmap.md`, `state.md` von `BEO-E2E`, `BEO-DOK/plan-aussage-ueber-code/` (`state.md`, Evidence), diesen Slice (Move nach `done/`) und Verweise, die `make slice-mv` nachzieht; `internal/`, `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/` und `tools/` stehen weiter nicht im Diff (Verifier, `git diff --stat ecde149..HEAD -- spec internal docs/plan/adr .github tools scripts packaging Makefile` leer).

**M-T ausgeführt (Trockenlauf des Tap-Pfads).** Befehl: `docker run -i --rm --network none -v "$PWD":/src:ro u-boot:test bash < <Skript>` (Image aus `make test`, Arbeitsbaum schreibgeschützt, ohne Netz; das Skript liegt im Scratchpad des Laufs, nicht im Repo). Ausgabe (Auszug, Exit-Codes):

```text
(1) homebrew-formula-fill.sh v0.12.0 <sums> scripts/homebrew-formula.rb.tmpl <Ausgabe>   exit=0   (vier url-Zeilen mit v0.12.0)
(2a) tap-formula-check.sh v0.12.0 <Ausgabe>   exit=0
(2b) tap-formula-check.sh v0.11.0 <Ausgabe>   exit=1
     tap-nachzug: the u-boot.rb asset of v0.11.0 does not carry four download URLs of v0.11.0 (found 4 url lines, 0 of v0.11.0).
(3)  tap-nachzug.sh, TAP_TOKEN=dummy-token, TAG=v0.12.0, curl- und git-Stellvertreter:   exit=0
     Protokoll: curl … -o <tmp>/u-boot.rb · git clone --depth 1 … · git diff --quiet · git … add · git … commit -m u-boot 0.12.0 · git push origin HEAD
(3b) tap-nachzug.sh mit TAG=v0.11.0 gegen dieselbe Formel: exit=1, Protokoll nur der curl-Aufruf, kein clone, kein push
```

Rot gesehen: (2b) und (3b) sind die rote Gegenprobe der Prüfung auf den Tag. **Grenze, wie in M-T benannt:** die Bedingung `is_stable == 'true'` und der Workflow-Schritt, der echte `curl` gegen das Release-Asset, der echte `git push` und das Token sind nicht geprobt; der Stellvertreter von `git` meldet immer eine Änderung. Befund: keiner; kein Skript und kein Workflow geändert.

**Befund zum roten Test (Auftrag 4) — gemessen mit eigenen Läufen in Wegwerf-Containern (Namen `u-boot-cut-*`, Image `u-boot-cut-img:sbx`, nach dem Lauf entfernt; das Projektverzeichnis lag im Scratchpad, nicht im Repo; keine Produktänderung):**

- **Reproduziert:** `docker run … u-boot:test-docker-tools go test -tags docker -count=1 -v -run TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress ./internal/e2e/` auf dem nativen `amd64`-Host: `--- FAIL` (31,16 s). `STEP 3/3: RUN ! getent hosts example.com …` gibt `2606:4700:10::ac42:93f3 example.com` und `2606:4700:10::6814:179a example.com` aus, der Bau bricht ab. Schritt 2 (`getent hosts deb.debian.org`) gelingt.
- **Ursache (Messung):** Im Sandbox-Container steht nach `u-boot-egress-init` `nameserver 127.0.0.1` (der lokale `dnsmasq`). Ein verschachtelter Container (Podman 6.1.4, `pasta`, `netavark`) bekommt `nameserver 169.254.1.1` (Weiterleitung von `pasta`), dazu die Rückfall-Nameserver `nameserver 8.8.8.8` und `nameserver 8.8.4.4`. Die Upstream-Liste des Docker-Hosts (`/etc/resolv.conf.u-boot-orig` im Container: `1.1.1.1`, `8.8.8.8`) enthält `8.8.8.8`; `egress-init.sh` lässt `ip daddr <Upstream> udp/tcp dport 53` durch. Der verschachtelte Container fragt also `8.8.8.8` direkt und bekommt Antworten für jeden Namen. Einzelmessungen: `--dns 169.254.1.1` allein: keine Antwort (der `dnsmasq` verweigert); `--dns 8.8.8.8`: Antwort; `--dns 127.0.0.1`: keine Antwort; `--network=none`: keine Antwort; Verbindung zur aufgelösten Adresse (`/dev/tcp/104.20.23.154/80`): Zeitüberschreitung, also gesperrt.
- **Gegenprobe (rot gesehen → grün):** derselbe Skript-Text des Tests im Sandbox-Container, gestartet mit `docker run --dns 9.9.9.9` (Upstream-Liste `9.9.9.9`): `STEP 3/3 … UNLISTED_BLOCKED`, `EGRESS_BUILD_OK`, Exit 0. Der Test hängt damit am Upstream-Resolver des Docker-Hosts, nicht an `amd64` und nicht an Podman 6.1.4. Auf Colima und auf dem GitHub-Runner ist der Upstream-Resolver nicht gemessen.
- **Einordnung:** Es ist die im Handbuch ausgewiesene Grenze („ein Prozess ohne `sudo` kann den Upstream-Resolver direkt fragen und erhält die Adresse eines Namens außerhalb der Liste; die Verbindung bleibt gesperrt“), hier für verschachtelte Container: der Rückfall `8.8.8.8` von Podman trifft einen Upstream des Hosts. Es ist kein Fehler des Cuts. Der Test unterstellt dagegen, dass ein unlisteter Name im `RUN` nicht auflöst; das gilt nur, wenn der Upstream-Resolver des Hosts kein Rückfall-Nameserver von Podman ist. Eine Korrektur (Test pinnen oder das Skript `egress-init.sh` ändern, sodass nur der lokale Resolver Port 53 darf) ist ein eigener Slice und eine Produktentscheidung; sie gehört nicht in den Cut. **Anker für die Korrektur:** [`slice-v1-egress-dns-upstream-schliessen`](../open/slice-v1-egress-dns-upstream-schliessen.md) (Pfad zum Upstream-Resolver für Prozesse ohne `sudo`); er schließt diese Ursache (die Regel zu Zeilen 63 bis 66 von `egress-init.sh.tmpl`) und führt den Test in Gegenstand und DoD. Der Cut sagt deshalb **nicht**, der Docker-Lauf sei grün; er nennt den einen roten Test und die Ursache.
- **Auf der CI:** `integration-docker` läuft bei `push` und `pull_request` auf `main` und per `workflow_dispatch` (`.github/workflows/integration.yml`, `on:`); der Job ruft `make test-docker` ohne `-v` auf (`repo.mk` `test-docker`). Die letzten Läufe auf `main` sind grün (`gh run list`, 2026-10-19: `integration-docker` 39187068499 auf `b9ae4e7`, `ci` 39187068493). Ob der Test dort läuft oder überspringt (`skipIfNestedUsernsRestricted` überspringt bei `kernel.apparmor_restrict_unprivileged_userns=1`), ist ohne `-v` und ohne abrufbares Log nicht zu entscheiden (`gh run view 39187068499 --log` lieferte keine Zeilen); ein Lauf auf dem GitHub-Runner ist nicht gemessen.

**CI-Hinweise an den Projektinhaber (Gate-Gegenstand, `AGENTS.md` §3.6; der Cut ändert kein Workflow):**

1. Der Kopfkommentar von `.github/workflows/integration.yml` beschreibt `continue-on-error: true` als zeitlich befristete, geltende Maßnahme; die Datei enthält das Flag nicht (der Kommentar stimmt nicht mit dem Workflow überein).
2. Dem Aufruf von `make test-docker` fehlt `-v` (`repo.mk`, `go test -tags docker -timeout 30m ./...`): ein Skip ist von einem Pass nicht zu trennen; Voraussetzung dafür, `amd64` im Sinn von [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10 als „verifiziert“ zu führen.
3. `integration-docker` ist kein Pflicht-Check, und `main` hat gar keine Branch Protection: `gh api repos/pt9912/u-boot/branches/main/protection/required_status_checks` meldet `Branch not protected` (HTTP 404, 2026-10-19). `docs/maintainer/quality.md` nennt drei Required Status Checks für `ci` (`gates`, `security-gates`, `image-scan`); `integration-docker` steht nicht darunter, und gesetzt ist laut der Abfrage keiner.
4. Das Zusammentreffen aus 2 und dem roten Test: lokal rot, auf der CI grün oder übersprungen, ohne `-v` nicht unterscheidbar.

**Aussage zum Docker-Lauf in diesem Cut:** `make gates`, `make ci` und `make fullbuild` sind grün; `make test-docker` ist **rot** (ein Test, Ursache gemessen, nicht Produkt-Defekt des Cuts). Eine grüne Closure des Docker-Laufs wird nicht behauptet.

**Befehle für Push und Tag (nach ausdrücklicher Freigabe des Projektinhabers, je Vorgang; nicht durch den Implementer):** siehe M0 bis M8 oben; M0 zuerst: `date +%F` muss `2026-10-19` sein, sonst Datum in `CHANGELOG.md`, beiden READMEs und im Handbuch-`Stand` nachziehen (ein weiterer Cut-Commit).

## 4. Trigger

**Start** (`next` → `in-progress`): Die Welle [`welle-sandbox-egress-ollama`](../done/welle-sandbox-egress-ollama.md) und die drei Slices ohne Welle ([`slice-cr-lastenheft-zielbild`](../done/slice-cr-lastenheft-zielbild.md), [`slice-v2-spec-zielbild`](../done/slice-v2-spec-zielbild.md), [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md)) liegen in `done/`; die Release-Gates 1 und 2 sind bestätigt (Vermerk 2026-10-19, `332b171`); `in-progress/` enthält nur die Master-Dokumente (WIP-Limit 1). Der Weg ist `open/` → `in-progress/` per reinem `git mv` (`make slice-mv SLICE=slice-v1-release-cut-v0.12.0 TO=in-progress`); der Roadmap-Marker geht in einen eigenen Commit direkt nach dem Move, beide im selben Push.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut ist ein Commit plus Tags.
- `in-progress` → `open` (blockiert): ein roter `image-scan` (Go-Pin heben), ein Befund des Trockenlaufs, der eine Änderung an Workflow oder Skript verlangt (eigener Vorgang davor), ein fehlschlagender `publish`-Lauf, dessen Ursache einen Patch-Release braucht, oder ein Befund des Reviewers, der eine Produktänderung verlangt.

## 5. Closure-Trigger

DoD vollständig: stabiler Tag gesetzt (M5), `publish`- und Tap-Lauf kontrolliert (M6), Trockenlauf M-T mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs und Steering-Loop-Eintrag, Verification Evidence (Verifier, getrennt von Implementer und Reviewer). Der Slice geht **erst nach dem Tag** nach `done/`; der Ruhe-Marker der Roadmap steht erst dann wieder.

## 6. Risiken und offene Punkte

Der Ausgang jedes Risikos ist bei der Closure festgeschrieben.

- **Der Job `tap` scheitert im stabilen Lauf, obwohl der Trockenlauf grün war** (der Trockenlauf probt die Skripte, nicht die Workflow-Bedingung, den echten Download, den Push und das Token). — **Ausgang:** nicht eingetreten: `tap-check` 39224782459 `success` (2026-10-19T08:43Z, vor dem Tag; Verifier V-5 erledigt) und der stabile `publish`-Lauf 39224848402 mit Job `tap` `success`; Log „tap-nachzug: pt9912/homebrew-u-boot updated to v0.12.0.“ (Secret aus Actions, keine Secret-Warnung) = der Job lief wirklich. Die Installation auf einem Mac (M7) steht beim Projektinhaber.
- **Der CI-Lauf `integration-docker` ist rot** (der Test `NestedPodman_BuildEgress` ist bekannt rot; ob er auf dem Runner läuft oder überspringt, ist offen; Anker für die Korrektur: [`slice-v1-egress-dns-upstream-schliessen`](../open/slice-v1-egress-dns-upstream-schliessen.md)). — **Ausgang:** nicht eingetreten: `integration-docker` 39223699209 auf dem Push-Tip `9c28c5c` `success` (Verifier V-2: der CI-Lauf auf dem Push-Tip ist grün). Ob der Test dort lief oder übersprang, bleibt ohne `-v` im Log unentscheidbar (Register `BEO-E2E/ci-ueberspringt-nested-podman-build`); der Test bleibt auf dem nativen `amd64`-Host rot, die Korrektur gehört dem offenen Slice.
- **Die Release-Aussage zu Cloud-Modellen fehlt oder ist zu weit.** — **Ausgang:** entfallen: CHANGELOG-Lead („Not included: cloud models of Ollama“), Fixed-Eintrag und Handbuch sagen es gleich (Review F-1 behoben, Verifier bestätigt, keine HIGH- oder MEDIUM-Befunde).
- **Das Tag-Datum weicht vom festen Datum ab** (CHANGELOG, Tabellenzeilen der READMEs, Handbuch-`Stand`). — **Ausgang:** entfallen: `git log -1 --format=%cs v0.12.0` = 2026-10-19; alle fünf Stellen tragen 2026-10-19, keine Datumsänderung bei der Closure.
- **Das Messprotokoll der Welle trägt Rohdaten mit Zugangsdaten** ([`BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten`](../observations/BEO-DOK/messprotokoll-rohdaten-mit-zugangsdaten/observation.md)); der Projektinhaber lässt die Historie unberührt. — **Ausgang:** entfallen: Entscheidung des Projektinhabers vom 2026-10-19 (§1); der Cut ändert nichts daran.
- **GHCR-Tags und Release früherer Vorab-Versionen** (`0.10.0-rc.1`, `0.10.0-rc.2`, `0.11.0-rc.1`) stehen noch; der Cut berührt sie nicht. — **Ausgang:** weiter offen beim Projektinhaber (Aufräumen von `v0.11.0-rc.1` samt GHCR-Tags und der GHCR-Tags `0.10.0-rc.1` und `0.10.0-rc.2`, Scope `delete:packages` oder Web-Oberfläche); ein Vorab-Tag `v0.12.0-rc.*` wurde nicht gesetzt, es entsteht keine neue Aufräumlast.
- **Befunde des Verifiers** (V-1 bis V-5, [`2026-10-19-verification-slice-v1-release-cut-v0.12.0`](../../../reviews/2026-10-19-verification-slice-v1-release-cut-v0.12.0.md)): V-1 (Satz „two runs“ zur Gegenprobe) erledigt in `9c28c5c` vor dem Push; V-2 erledigt (CI auf dem Push-Tip grün, siehe oben); V-3 (Plan-Tabelle) in §3 ergänzt und als Vorgang vierundzwanzig im Register `BEO-DOK/plan-aussage-ueber-code` geführt; V-4 ohne Wirkung; V-5 erledigt (`tap-check` 39224782459).

## 7. Closure-Notiz

- **Was hat funktioniert:** Der stabile Tag `v0.12.0` liegt auf `9c28c5c2e97de6e055385913353524abe81ed827` (leichter Tag, 2026-10-19, `git rev-parse v0.12.0`, `git log -1 --format=%cs v0.12.0`). Kette: Verifier-GO mit Auflagen ([`2026-10-19-verification-slice-v1-release-cut-v0.12.0`](../../../reviews/2026-10-19-verification-slice-v1-release-cut-v0.12.0.md), `aeb11d6`; Review [`2026-10-19-slice-v1-release-cut-v0.12.0`](../../../reviews/2026-10-19-slice-v1-release-cut-v0.12.0.md), 0 HIGH, 0 MEDIUM, 3 LOW behoben) → Push von `main` (`b9ae4e7..9c28c5c`) → `ci` 39223699156 `success` (gates, security-gates, image-scan) und `integration-docker` 39223699209 `success` auf dem Push-Tip → `gh workflow run tap-check` 39224782459 `success` (2026-10-19T08:43Z, vor dem Tag) → Tag `v0.12.0` → `publish` 39224848402 `success` (Job `publish` mit GHCR-Image-Push, Job `tap`, Paket-Smoke rpm/fedora, Paket-Smoke deb/ubuntu alle `success`). GitHub-Release: `prerelease=false`, `draft=false`, `publishedAt` 2026-10-19T08:45:16Z, zwölf Assets (`SHA256SUMS`, `u-boot-0.12.0-1.aarch64.rpm`, `u-boot-0.12.0-1.x86_64.rpm`, `u-boot_0.12.0_amd64.deb`, `u-boot_0.12.0_arm64.deb`, `u-boot-darwin-amd64`, `u-boot-darwin-arm64`, `u-boot-linux-amd64`, `u-boot-linux-arm64`, `u-boot.rb`, `u-boot-windows-amd64.exe`, `u-boot-windows-arm64.exe`). Image: `docker run --rm ghcr.io/pt9912/u-boot:0.12.0 --version` = `u-boot version 0.12.0`; `:0.12.0` und `:latest` tragen denselben Digest `sha256:9a625aa36ad12912fbe4d8809439142a32fc6e4e59ca3e658df177d0c75b843e`. Job `tap`: Log „tap-nachzug: pt9912/homebrew-u-boot updated to v0.12.0.“ (Secret aus Actions, keine Secret-Warnung) = wirklich gelaufen, nicht übersprungen und nicht „already at“. Die Läufe und Messungen von M2, M4 und M6 stammen vom Hauptagenten (Befehle: `gh run list`, `gh workflow run tap-check`, `gh release view`, `docker run … --version`); der Planner hat sie nicht selbst wiederholt. Das Tag-Datum (2026-10-19) stimmte an allen fünf Stellen aus M0; die Closure ändert kein Datum.
- **Freigabe von Push und Tag:** ausdrückliche Anweisung des Projektinhabers im Gespräch („pushen, dann direkt v0.12.0 ohne rc“); Push und Tag führte der Hauptagent, nicht der Implementer. Das Löschen von Vorab-Tags und GHCR-Tags ist nicht freigegeben und nicht geschehen.
- **Probe M-T:** Ergebnis in §3 festgehalten (Formel aus Fixture, `tap-formula-check.sh` Exit `0` mit `v0.12.0` und Exit `1` mit `v0.11.0`, `tap-nachzug.sh` mit Stellvertretern; vom Verifier unabhängig nachgefahren, gleiches Ergebnis). Die benannte Grenze (`is_stable == 'true'`, echter `curl`, echter `git push`, Token) schloss der echte stabile Lauf zusammen mit `tap-check`.
- **Was ging anders als geplant:** (a) Kein Vorab-Tag: Der Projektinhaber entschied „ohne rc“; der Verifier belegte die Voraussetzung (kein Diff seit `v0.11.0` in `.github`, `scripts`, `packaging`, `tools/harness`, `Makefile`; stabiler Lauf `v0.11.0` in allen vier Jobs grün) und der Lauf bestätigte sie. (b) Der Verifier-Befund V-1 (Satz zur Gegenprobe mit `--dns 9.9.9.9`) wurde vor dem Push berichtigt (`9c28c5c`), weil der CHANGELOG-Eintrag eines Tags danach nicht mehr geändert wird. (c) Beim Push meldete GitHub „Bypassed rule violations for refs/heads/main … Cannot update this protected ref“: `main` hat ein Ruleset, obwohl die klassische Branch-Protection-Abfrage 404 liefert (Hinweis unten). Kein Tag wurde verschoben, kein Patch-Release, keine Änderung an Workflow oder Skripten.
- **Was belegt ist und was nicht:** Belegt sind Image (`:0.12.0` = `:latest`, `--version`), Assets, Pakete (Smokes), der Tap-Nachzug (Log), CI auf dem Push-Tip. **Nicht belegt:** (1) M7: `brew tap pt9912/u-boot`, Installation von `u-boot` und `u-boot --version` = `0.12.0` auf einem Mac steht beim Projektinhaber; (2) ob `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` in der CI lief oder übersprang: `integration-docker` ist `success`, der Aufruf läuft ohne `-v`, das Log entscheidet es nicht; der Test ist auf dem nativen `amd64`-Host rot (Upstream-Resolver des Docker-Hosts, §3); (3) `ollama` und `podman build` auf `amd64` bleiben Einzelmessungen ohne CI-Beleg, wie im CHANGELOG ausgewiesen.
- **M6-Befunde:** zwölf Assets (Liste oben), kein Prerelease, kein Draft; `:0.12.0` und `:latest` gleicher Digest; `--version` = `0.12.0`; Job `tap` mit Log-Zeile `updated to v0.12.0.`. Kein Befund, der eine Nacharbeit am Release verlangt.
- **Steering-Loop-Eintrag:**
  - Lerneintrag (Lehre, nicht verkörpert): Ein Release ohne Vorab-Tag ist vertretbar, wenn drei Dinge belegt sind: (i) die Pipeline-Dateien haben seit dem letzten stabilen, grünen Lauf keinen Diff (Verifier: `git diff --stat v0.11.0..HEAD -- .github scripts packaging tools/harness Makefile` leer), (ii) der Trockenlauf M-T deckt die Skripte, (iii) `tap-check` unmittelbar vor dem Tag deckt das Token. Der Vorab-Tag hätte `is_stable == true` und den Job `tap` ohnehin nicht geprobt. Zweitens: Ein Befund des Verifiers zum Release-Text fällt **vor** dem Push (V-1), weil der CHANGELOG-Eintrag eines Tags danach nicht mehr geändert wird; der Verifier lief deshalb vor M1.
  - Gezählt, ohne `liegt in`-Feld (gezählt, nicht verkörpert): [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): vierundzwanzigster Vorgang (V-3: Plan-Tabelle ohne Zeilen für den Slice selbst, die vier Vermerke in `done/` und den Verification-Report; Rückfall nach der Verkörperung in `AGENTS.md` §3.10; Evidence-Datei angelegt, `state.md` fortgeschrieben); keine weitere Hard Rule und kein Sensor ohne Entscheidung des Projektinhabers (die Frage nach einem Teilsensor liegt dort). [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md): Trigger „CI-Lauf auf dem Push-Tip“ eingetreten, CI grün, Test-Verhalten in der CI ohne `-v` weiter unentscheidbar; keine neue Evidence-Datei, `state.md` fortgeschrieben, Beobachtung bleibt offen.
  - Benannte Spec-Lücke: keine.
- **Hinweis an den Projektinhaber, kein Eingriff (Gate-Gegenstand, [`AGENTS.md`](../../../../AGENTS.md) §3.6):** `main` hat ein Ruleset (GitHub meldete beim Push übergangene Regelverstöße), die klassische Branch-Protection-Abfrage lieferte 404; welche Checks das Ruleset verlangt, ist nicht gelesen. `docs/maintainer/quality.md` nennt nur die klassische Form.
- **Beobachtungs-Register (`../observations/`):** `BEO-E2E/ci-ueberspringt-nested-podman-build` (`state.md` fortgeschrieben, `cbd00d7`-Stand) und `BEO-DOK/plan-aussage-ueber-code` (Evidence, `state.md`); keine neue Beobachtung angefallen.
- **Folge-Slices** (alle in `open/`, nicht Teil dieses Cuts): [`slice-v1-egress-dns-upstream-schliessen`](../open/slice-v1-egress-dns-upstream-schliessen.md) (Test-Anker: `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` im Gelesenen Stand, in der Plan-Tabelle und in der DoD), [`slice-v1-egress-init-ohne-luecke`](../open/slice-v1-egress-init-ohne-luecke.md), [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../done/slice-v1-e2e-pin-mismatch-zwischencontainer.md), [`slice-v1-egress-hinweise-neuaufbau-und-messstand`](../done/slice-v1-egress-hinweise-neuaufbau-und-messstand.md), [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md), [`slice-v1-messskripte-nach-tools`](../done/slice-v1-messskripte-nach-tools.md). **Folge-Wellen** (Roadmap, *Nächste Wellen*): `welle-ollama-cloud-modelle` (Trigger „Release der Welle … geschnitten“ eingetreten) und `welle-forge-cli-gh-glab` (Trigger „Eröffnung nach dem Release-Cut“ eingetreten; beauftragt, noch nicht eröffnet).
- **Offen beim Projektinhaber:**
  - M7: Installation über den Tap auf einem Mac (`u-boot --version` = `0.12.0`).
    *Vermerk 2026-10-19:* Der Projektinhaber hat M7 im Gespräch bestätigt: Installation über den Tap auf einem Mac erfolgreich, `u-boot --version` = `0.12.0`. Zusätzlich lief auf Linux mit Linuxbrew ein Upgrade von `0.11.0` auf `0.12.0` mit verifiziertem Download (Ausgabe vom Projektinhaber, vom Orchestrator nicht selbst gemessen).
  - Aufräumen von `v0.11.0-rc.1` (GitHub-Release, Tag, GHCR-Tag `0.11.0-rc.1`) und der GHCR-Tags `0.10.0-rc.1` und `0.10.0-rc.2`.
  - Gate-Änderungen zur CI ([`AGENTS.md`](../../../../AGENTS.md) §3.6): `-v` im Aufruf von `make test-docker`; `integration-docker` als Pflicht-Check (nach Klärung des Rulesets auf `main`); der Kopfkommentar von `integration.yml` (nennt `continue-on-error: true`, die Datei enthält das Flag nicht).
  - Wortlaut des Lastenheft-CR in [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) (Freigabe je Stelle).
  - Folge-ADR zu [`ADR-0025`](../../adr/0025-egress-liste-im-laufenden-container.md), falls die Korrektur des DNS-Upstream-Pfads die Entscheidung der ADR berührt (nach den Messfragen des Folge-Slice).
  - Verzeichnisname der Messskripte in [`slice-v1-messskripte-nach-tools`](../done/slice-v1-messskripte-nach-tools.md) (Entscheidung (c), `tools/messung`).
  - Ob der Teilsensor für `BEO-DOK/plan-aussage-ueber-code` (Dateiliste der Plan-Tabelle gegen `git diff --stat`) kommt.
- **Risiken aus §6:** siehe §6 (Ausgänge festgeschrieben).
- **Drei Paarungen:** Anker: keine neue verkörperte Regel, kein `liegt in`-Feld, nichts zu paaren · Folge-Slice: die sechs Slices in `open/` und die beiden Folge-Wellen (oben), offene Punkte beim Projektinhaber · Register: `BEO-E2E/ci-ueberspringt-nested-podman-build` und `BEO-DOK/plan-aussage-ueber-code` fortgeführt, keine neue Beobachtung.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/`, `docs/maintainer/`, `docs/reviews/` und Release-Dateien (`repo.mk`, `Dockerfile`) — Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006). `cmd/uboot` ist laut Modus-Deklaration keine Sub-Area (Aspirantin); berührt wird dort eine Zeile (`var version`), keine Konvention.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen ([`../observations/`](../observations/README.md)): [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) ist `offen` und betrifft diesen Cut direkt (Auftrag 6); `BEO-DOK/plan-aussage-ueber-code` und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` sind `verkörpert` (`AGENTS.md` §3.10: Aussagen über Code mit Fundstelle, Bedingungen in Handbuch mit Test oder gestrichen — der neue Handbuch-Satz zu Cloud-Modellen nennt nur Gemessenes).

Alle berührten Sub-Areas GF.
