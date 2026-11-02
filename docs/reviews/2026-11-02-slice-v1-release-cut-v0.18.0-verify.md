# Verifikations-Report: slice-v1-release-cut-v0.18.0 — 2026-11-02

**Rolle:** Verifier (`harness/roles.md` §Verifier — „built the thing right")
**Gegenstand:** [`slice-v1-release-cut-v0.18.0`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md), Diff `v0.17.0` (`f89424e1`)..HEAD; Cut-Commit `e6c87379` (Tag `v0.18.0`), Nachzüge `9a7db2f1` (Review-Fix-Runde, kein Release-Inhalt) und `724c5e5e` (d-check-Befund am Review-Report)
**Skill:** `.harness/skills/verifier` (Rollenvertrag [`roles.md`](../../harness/roles.md), Schema [`verification.md`](../../harness/verification.md)) · **Modell:** glm-5.3-flash · **Datum:** 2026-11-02
**Basis:** Review-Report [`2026-11-02-slice-v1-release-cut-v0.18.0.md`](2026-11-02-slice-v1-release-cut-v0.18.0.md) (0 HIGH/2 MEDIUM/8 LOW/1 INFO), DoD des Slice, [`releasing.md`](../maintainer/releasing.md) §§3–8

---

## Scope

- Slice: [`slice-v1-release-cut-v0.18.0`](../plan/planning/done/slice-v1-release-cut-v0.18.0.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration), [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache), [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Pins unverändert), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md), [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md), [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)
- Artefakte: [`CHANGELOG.md`](../../CHANGELOG.md) Block `[0.18.0]` · Version v0.18.0 (Tag `e6c87379`) · Release `u-boot v0.18.0` (12 Assets) · GHCR `pt9912/u-boot:0.18.0` · Tap-Commit `b52a0bc` im Repo `pt9912/homebrew-u-boot`

## Lead gegen die drei Lastenheft-Stände (Plan-Risiko 1, Review F-1/F-2/F-3)

Gelesen am Fix-Runde-Stand [`CHANGELOG.md`](../../CHANGELOG.md) §[0.18.0] gegen
[`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration),
[`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer),
[`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
und die Historie 0.27.0/0.28.0/0.29.0; Codefundstellen der Shell-Aussagen selbst gelesen
(`devcontainer_start.go:520–525` — die Abbildung von `ErrDevcontainerContainerNameConflict` auf den
Fehler der Startklasse; `TestExecute_DevcontainerShell_NameOfTheSessionContainer`,
`internal/adapter/driving/cli/devcontainer_shell_test.go:457`):

- **Shell- Satz (F-1 behoben):** der Lead sagt jetzt: der namens-zusammensetzende Befehl trägt die
  Fehlerklassen in `up` **und** `devcontainer shell`, die wie `up` ohne `--recreate` startet —
  Wortlaut-kongruent zur Vertragszeile „startet den Devcontainer wie `u-boot devcontainer up` ohne
  `--recreate`; für den Start gelten die Fehlerklassen von `up`"; Namesteil-/Konflikt-Exit `10` vor
  dem Hilfs-Container stimmt mit dem Namens-Absatz überein; `down` nur als Label-Suche benannt —
  keine „work as before"-Behauptung mehr, die der eigenen Lieferung widerspricht.
- **Transport (F-3 behoben):** „by name on the command line, the value in the environment of the
  `docker` client process" — entzweit zu [ADR-0040](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md);
  der Vollinhalt „nie in einem Image, Volume, keiner Kommandozeile, Ausgabe" steht im
  `Added`-Eintrag (auch „Volume" fehlt im Lead-Kompromiss, nicht im Eintrag — keine Zusage
  weiter als das Lastenheft).
- **Not-promised (F-2 behoben):** „the behaviour of a host `gh` with several accounts" steht jetzt
  in der Lead-Liste — auch dort im `Added`-Eintrag von 0.28.0; die übrigen Grenzen (echter Push,
  andere Hosts/Engines, Forge-Token, `git credential fill`, kleines Limit, Summe über Instanzen,
  Wirkung auf das verschachtelte Docker) decken die drei Stände ab.
- Kein Punkt des Lead zugesagt, den die Stände 0.27.0/0.28.0/0.29.0 nicht tragen; die
  Testbasis-Erweiterung der Promised-Liste (Unit-Tests, Fake-Sources neben Docker-Tests mit
  Build-Tag `docker`) erweitert nicht über „nur, soweit ein Test es belegt" und die
  Host-Scope-Bedingung (Linux `amd64`, lokale Engine) hinaus. — **geprüft, ohne Befund.**

## DoD-Abgleich

| DoD-Punkt | Stand | Evidence |
|---|---|---|
| Release-Text und Versionsstellen: `## [0.18.0] - 2026-11-02` mit Lead („Eighteenth release."), leerem `[Unreleased]`, Compare-Links; `0.18.0-dev` an den drei Stellen samt Kommentar; beide Readmes und Handbuch auf `v0.18.0`; `Stand` = Tag-Datum | ✓ erfüllt | [`CHANGELOG.md`](../../CHANGELOG.md) Zeile 13/1592/1593; kein `0.17.0`-Rest in `cmd/uboot/main.go:43`, `repo.mk:61`, `Dockerfile:45`; `README.md:10,98`, `README.de.md:11,101`; `benutzerhandbuch.md:4,5,84,110,136,3111`; `git rev-parse v0.18.0` = `e6c873798380991ff6c072766085a1e5e1f1590e` |
| Vorab-Prüfungen: `make gates`, `make ci`, `make fullbuild` lokal grün; Tap-Trockenlauf mit Gegenprobe; `ci`+`integration-docker` auf dem Push-Tip grün; `tap-check` grün; danach Tag gesetzt und gepusht; Kontrolle nach §6 | ✓ erfüllt (M7 offen, s. u.) | Stempel `.harness/state/gates-passed.head` = HEAD `724c5e5e` und Inhalts-Hash `7baa2ba9…` = eigener Neulauf von `tools/harness/working-tree-hash.sh`; `/tmp/c1.log` (Exit 0, „ci … green") und `/tmp/fb1.log` („fullbuild ci + runtime image green", 2026-11-02 06:53/06:55, nach dem Cut-Commit 06:50Z); CI selbst nachgelesen: `gh run view` — `ci` `41403124364` success, `integration-docker` `41403124438` success, `tap-check` `41403131805` success, alle drei auf `e6c87379`; Trockenlauf **im Verifier-Lauf neu gefahren** (s. u. Sensors); `gh release view v0.18.0`: draft=false, prerelease=false, published `2026-11-02T07:17:01Z`, 12 Assets |
| Doku nachgezogen: Roadmap-Zeiger → nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash; Bestands-Readmes-Zeilen | **offen — ist Closure-Schritt** | Die Roadmap-Zeile des Cuts steht ehrlich noch unter *Offene Wellen* ([`roadmap.md`](../plan/planning/in-progress/roadmap.md) Zeile 24, Freigabe-Zeile auf dem Ist-Stand); der Zug mit Tag-Hash gehört zur Closure (`releasing.md` §3 Punkt 4: „bis der Tag gesetzt ist; danach") — kein Befund, aber Bestandteil der noch ausstehenden Closure; F-5 (Änderungsdatum `2026-11-02`) und F-8 (geplante Version `v0.19.0`) in der Fix-Runde `9a7db2f1` berichtigt |
| `make gates` grün nach dem letzten Commit vor Handoff und nach der Closure | ✓ Teile 1 / **offen Teil 2** | Stempel deckt HEAD `724c5e5e` und den exakten Baum (Inhalts-Hash); der Lauf nach der Closure ist mit dem Closure-Commit zu fahren (`record-gates` stuft nur nach grünem Lauf ab) |
| Review (kein Self-Review) und Verification mit Reports unter `docs/reviews/` | ✓ erfüllt, sobald dieser Report vorliegt | Review-Report [`2026-11-02-slice-v1-release-cut-v0.18.0.md`](2026-11-02-slice-v1-release-cut-v0.18.0.md) (committet `9a7db2f1`, Nachzug `724c5e5e`); dieser Verifier-Report |
| Closure-Notiz mit Steering-Loop-Lerneintrag; Register fortgeschrieben oder „keine Beobachtung"; jedes Risiko mit Ausgang; Paarungen getragen | **offen — ist Closure-Schritt** | §7 des Slice trägt noch den Plan-Marker (Review F-4); bei der Closure zu ersetzen — der Lerneintrag muss die Schwelle-Klasse „Zustandsmarker nach Ablauf nicht abgebaut" (Review §Steering-Loop-Vermerk) und die Paarungen der drei Träger-Slices tragen |
| M7 (Mac `brew upgrade u-boot`, `u-boot --version`) | **offen — Aktion des Projektinhabers** | nicht gelaufen; der CI-Job `tap` (geprüft: `41405100611` Job `tap` success) ersetzt die macOS-Sondierung nicht — gehört nicht in dieses Verifikations-Verdict, als offener Punkt ausgewiesen |

## Review-Findings: jeder Befund am Fundstellen-Stand nach der Fix-Runde

| ID | Prüfung nach `9a7db2f1`/`724c5e5e` |
|---|---|
| F-1 (MEDIUM) | behoben — Lead-Satz trägt den Ist-Stand (Shell-Fehlerklassen, `init`/`generate` ändern geschriebene Dateien, `down`-Label-Suche); Fundstellen gegen Code und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) selbst gelesen (s. o.). |
| F-2 (MEDIUM) | behoben — Grenzliste um den `gh`-Mehr-Accounts-Punkt ergänzt, deckungsgleich mit dem `Added`-Eintrag (0.28.0). |
| F-3 (LOW) | behoben — Transportformulierung entzweit; Lead und `Added`-Eintrag tragen dieselbe Lesart. |
| F-4 (LOW) | nicht aufgehoben — ist mit der Closure selbst zu lösen: §7 des Slice trägt noch „Noch nicht geschlossen; der Slice liegt in `open/` (Plan)". Die Closure-Notiz ersetzt den Marker **vor** dem `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3 Fall 2). Bedingung der Move-Freigabe, kein Blocker für dieses Verdict. |
| F-5 (LOW) | behoben — Roadmap „Letzte Änderung: 2026-11-02" (Diff `9a7db2f1`). |
| F-6 (LOW) | behoben — DoD-Zeile berichtigt auf den gestrichenen Zustand (Diff `9a7db2f1`); die DoD-Tabelle oben liest die berichtigte Fassung. |
| F-7 (LOW) | behoben — [`releasing.md`](../maintainer/releasing.md) Punkt 5: „aus der Skelett-Vorlage `slice.template.md` anlegen" (Diff `9a7db2f1`); Residuum entfernt. |
| F-8 (LOW) | behoben — Roadmap *Nächste Wellen* führt die Katalog-Welle als `v0.19.0` (Diff `9a7db2f1`). |
| F-9 (LOW) | behoben — beide CR-Entwürfe ([`slice-v1-katalog-gepinnte-toolchain-entwurf`](../plan/planning/open/slice-v1-katalog-gepinnte-toolchain-entwurf.md) inkl. §5 Frage 5, [`slice-v1-katalog-feature-jdk-gepinnt`](../plan/planning/open/slice-v1-katalog-feature-jdk-gepinnt.md)) und `open/README.md` tragen `v0.19.0` (Diff `9a7db2f1`). |
| F-10 (LOW) | behoben — Kopfliste des Slice nennt die Host-git-Messungs-Commits und die zwei `tools/mk`-Commits `82682508`/`7bbb8185` (Diff `9a7db2f1`). |
| F-11 (INFO) | keine Aktion nötig — keine [Regel](../maintainer/releasing.md) §3/§5 verlangt die Anhebung der Handbuch-Version; `Stand` = 2026-11-02 = Tag-Datum bestätigt. |

Keine offene Lücke aus dem Review; die in der Fix-Runde nicht berührten Negativbefunde des
Review (Versionszeilen samt Kommentare, CHANGELOG-Fuß, Readmes, Handbuch, Plan-Moves,
[`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)-Nachzug `827601be`, `git mv`-Disziplin
§3.3) wurden am Fix-Runde-Stand stichprobenartig nachgelesen — unverändert (die Fix-Runde berührt
diese Stellen nicht).

## Plan-vs-Code

Abgeglichen für den Stand nach der letzten Fix-Runde (`724c5e5e`).

**Cut-Commit `e6c87379`** — `git show --stat`: 7 Dateien, jede mit Plan-Zeile:

| Datei | Plan-Tabellen-Zeile | Befund |
|---|---|---|
| `CHANGELOG.md` | „[Unreleased]-Inhalt unter `## [0.18.0]` mit Lead …" | Deckung (Lead post-Fix-Runde geprüft, s. o.) |
| `cmd/uboot/main.go` | „`0.17.0-dev` → `0.18.0-dev` …" | Deckung (`:43`) |
| `repo.mk` | analog | Deckung (`:57,60,61`) |
| `Dockerfile` | analog | Deckung (`:41,43,45`) |
| `README.md`, `README.de.md` | „Status-Zeile und `VERSION=`-Beispiel" | Deckung (`README.md:10,98`, `README.de.md:11,101`) |
| `docs/user/benutzerhandbuch.md` | „Kopf, Stand, Versionsbeispiele, Gültigkeitsbereich" | Deckung (`:4,5,84,110,136,3111`) |
| `releasing.md`, `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | „kein Diff" (im Cut) | Bestätigt — der Cut-Commit berührt sie nicht |

**Slice-Commits außerhalb des Cut-Commits (Abgleich gemäß Auftrag):**

| Commit | Dateien | gegen Plan-Tabelle |
|---|---|---|
| `827601be` (MR-003-Nachzug, vor dem Cut) | `MR-003…`, `harness/conventions.md`, `roadmap.md`, `planning/README.md`, `plan-welle.md`, `close-welle.md`, `planner.md`, Slice selbst | Deckung je Zeile der Plan-Tabelle (MR-Zeilen und roadmap/README/command-Zeilen) |
| `b160d4ac`, `77e88143`, `44c6b644`, `dcc3697c` (Moves + Ruheort-Nachzüge) | Planungs-Doku | Deckung (Plan-Zeile „dieser Slice · `make slice-mv` je Lifecycle-Schritt") |
| `9a7db2f1` (Review-Fix-Runde) | `CHANGELOG.md` (Plan-Zeile), `releasing.md` (1 Zeile, F-7), `roadmap.md` (F-5/F-8), Slice (F-6/F-10), `open/README.md` (F-9), zwei CR-Entwürfe (F-9), Review-Report | Deckung bis auf zwei Stellen s. u. |
| `724c5e5e` (d-check-Korrektur) | Review-Report, 1 Zeile | Abgleich-Artikel, kein Release-Inhalt |

**Abweichungen vom Plan-Abgleich — benannt, nicht als Befund gewertet:**

1. `docs/plan/planning/open/slice-v1-katalog-*.md` (2 Dateien) und der Review-Report stehen nicht
   in der Plan-Tabelle §3. Beide sind Produkte der Review-Runde (F-9 in den Planning-Doku außerhalb
   des Release, Report-Pflicht DoD-Zeile 5) — der Plan-Abgleich nach [`AGENTS.md`](../../AGENTS.md)
   §3.10 nennt sie, wie hier geschehen; der Plan-Stand vom 2026-11-01 konnte die Review-Produkte
   nicht vorwegnehmen.
2. Die Plan-Zeile „`releasing.md` — **im Cut kein Diff**" bleibt für den Cut-Commit wahr; die
   Fix-Runde berührt `releasing.md` (1 Zeile, F-7) — eine Review-Korrektur des Vertragstexts, nicht
   ein nachträgliches Nachziehen des Cut.

## Sensoren

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass (Docker-Modus) | Stempel `.harness/state/gates-passed.head` = HEAD `724c5e5e1e32068137135bb4bcb41a320e5a926a`; Inhalts-Hash `.harness/state/gates-passed.diffsha` = `7baa2ba9c445838d679b315d56e3478c7a95add96353417dc29604cb5935d790` = eigener Neulang von `tools/harness/working-tree-hash.sh` über den aktuellen Baum; `record-gates.sh` schreibt den Stempel nur als gates-Prerequisite nach grünem Lauf — der Nachweis deckt den aktuellen Tree exakt, ein Neulauf war nicht nötig |
| Tap-Trockenlauf (M-T) **im Verifier-Container neu gefahren** | pass inkl. beider Gegenproben | `docker run --network none` in `u-boot:test` mit `/repo:ro` und Stellvertreter-`curl`/`git`: `homebrew-formula-fill.sh v0.18.0` → Formel mit vier Digests; `tap-formula-check.sh v0.18.0` → Exit 0; `tap-formula-check.sh v0.17.0` → Exit 1 („does not carry four download URLs of v0.17.0"); `tap-nachzug.sh TAG=v0.18.0` (Dummy-Token, Stubs) → Protokoll `clone`, `add`, `commit -m u-boot 0.18.0`, `push`; `TAG=v0.17.0` → Exit 1 **vor `clone`** |
| `make doc-immutable RANGE="v0.17.0..HEAD"` | pass | d-check vcs: 98 Commits aufgelöst (96 beim Cut-Stand + die zwei Nachzugs-Commits des Standes `724c5e5e`), 0 Befund(e), Exit 0 |
| Remote-Image-Gegenprobe | pass | `docker run --rm ghcr.io/pt9912/u-boot:0.18.0 --version` → `u-boot version 0.18.0`; Manifest-Digests `ghcr.io/pt9912/u-boot:0.18.0` == `:latest` == `sha256:bef5b03dff71467cc439ed015e98f0b8552e827ddef79789f8c19ca82ad50383` (`docker buildx imagetools inspect`) |
| Lokale Build-Gegenprobe | pass | `docker inspect u-boot:latest` → OCI-Label `org.opencontainers.image.version` = `0.18.0`; `docker run --rm u-boot:latest --version` → `u-boot version 0.18.0` (Image aus `make build VERSION=0.18.0` des Hauptagenten) |
| CI-Läufe nachgelesen (`gh run view`) | pass | `ci` `41403124364`, `integration-docker` `41403124438`, `tap-check` `41403131805` — alle completed success auf `e6c87379`; nach dem Tag: `publish` `41405100611` success mit Jobs `publish (ghcr.io image push)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)`, `tap (homebrew formula)` — je success |
| Tap-Nachzug-Ergebnis | pass | Commit `b52a0bc479adf53483883e4c00330b6f4972fe20` in `pt9912/homebrew-u-boot` mit Message „u-boot 0.18.0" (`gh api`); der Job `tap` des Publish-Laufs ist der Beleg für den echten Push |
| Release-Metadaten | pass | `gh release view v0.18.0`: draft=false, prerelease=false, published `2026-11-02T07:17:01Z`, 12 Assets |
| `make ci` / `make fullbuild` | pass (Hauptagent) | `/tmp/c1.log` („[ci] gates + govulncheck + image-scan green", 2026-11-02 06:53), `/tmp/fb1.log` („[fullbuild] ci + runtime image green", 06:55) — nach dem Cut-Commit (06:50Z) gefahren; Logs des Hauptagenten vorgefunden und Endzeilen gelesen |
| `make test-docker` | not run | Entscheidung des Projektinhabers (2026-11-01, Plan §Nicht-in-Slice): der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis — selbst nachgelesen (success `41403124438` auf `e6c87379`); Plan-Schritt M-T/Trockenlauf deckt die Skript-Pfade des Tap-Jobs |
| `make verify-depguard` | not run | Der Slice berührt keine Importregeln (Produktcode unverändert, nur Versionszeilen in `cmd/`, `repo.mk`, `Dockerfile`) |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | CHANGELOG-Block `[0.18.0]` (Containername, Instanz-Aufrufe, Shell-Fehlerklassen) gegen den Vertrag und die Fundstellen `devcontainer_start.go:525`, `devcontainer_shell_test.go:457` gelesen; Release als Lastenheft `0.29.0` |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | CHANGELOG (`hostGitToken`, Identität, Grenzen) gegen den Vertrag; Lastenheft `0.28.0` |
| [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) | Ressourcen-Schlüssel im CHANGELOG gegen den Abschnitt `Inhalt der Konfiguration`; Lastenheft `0.27.0` |
| [`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache) | Hygiene-Eintrag des CHANGELOG, Beleg-Commit `fcbafc2e` im Range (Reviewer-Negativbefund) |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | Trockenlauf (Verifier-Lauf, s. o.), `tap-check` `41403131805` success, Tap-Commit `b52a0bc` |
| [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) | Release `u-boot v0.18.0` (12 Assets), GHCR `0.18.0` == `latest` (Digest), Package-Smoke-Jobs |
| [`ADR-0040`](../plan/adr/0040-host-git-quellen-fuer-token-und-identitaet.md) | Transport-Formulierung im Lead post-F-3 entzweit gegen den ADR |
| [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md) | Nachzug `827601be` — wellenlose Cut-Praxis als Adaption; Roadmap-Format-Regel auf dem Ist-Stand |

## Carveouts

- Neu: none
- Gelöst: none
- Unverändert: der Cut senkt keine Schwelle und berührt kein Gate; [`doc-immutable`](../maintainer/releasing.md) über den Range grün (98 Commits, 0 Befund) belegt, dass Accepted-ADRs unangetastet blieben.

## Nicht ausgeführt

- `make gates`-Neulauf — Stempel deckt HEAD **und** den Inhalts-Hash des aktuellen Baums; Mechanismus von [`record-gates.sh`](../../tools/harness/record-gates.sh) gelesen (Stempel nur nach grünem gates-Prerequisite).
- `make test-docker` — Entscheidung des Projektinhabers; CI-Lauf `integration-docker` `41403124438` success als Nachweis selbst nachgelesen.
- `make ci` / `make fullbuild` erneut — bereits grün beim Hauptagenten (Logs gelesen); ein Zweitlauf brächte keinen neuen Erkenntniswert und der Tag ist bereits live.
- `make verify-depguard` — keine Importregel berührt.
- M7 (macOS-Probe des Tap) — Aktion des Projektinhabers; kein Sensor des Repo.

## Rest-Risiken

1. **M7 offener Punkt:** `brew upgrade u-boot` und `u-boot --version` auf einem Mac sind nicht gelaufen. Der CI-Job `tap` (success in `41405100611`) und der Tap-Commit `b52a0bc` belegen, dass die Formel im Tap steht und aus den Release-Assets gefüllt wurde; die macOS-Installations-Sondierung bleibt ein Punkt des Projektinhabers.
2. **Closure-Schritte offen** (bewusst, vor dem `git mv` zu fahren): §7-Closure-Notiz (löst F-4) mit Steering-Loop-Lerneintrag — die Klasse „Zustandsmarker nach Ablauf nicht abgebaut" hat die Schwelle erreicht (Review §Steering-Loop-Vermerk) —, Risiken §6 mit je Ausgang, Register-Fortführung, Roadmap-Zug unter *Abgeschlossene Wellen* mit Tag-Hash, danach `make gates` frisch auf dem Closure-Stand, dann der reine `git mv`.
3. **Der Lead ist der gelesenste Text des Release:** die Grenzlisten der drei `Added`-Einträge tragen die Vertragsaussagen vollständig; jede künftige Änderung an den Träger-Zusagen braucht eine Lead-Spiegelung (Klasse „Lead-Grenzliste kürzer als Vertragsliste" — zweite Wiederholung, Review §Steering-Loop-Vermerk).

## Verdikt

**Der Release-Schnitt `v0.18.0` ist DoD-Konform verifiziert und closure-fähig.** Jede DoD-Zeile
trägt einen Beleg; die offenen Punkte sind genau die Closure-Schritte selbst (§7-Notiz, Roadmap-Zug,
gates nach der Closure, `git mv`) und die vom Projektinhaber getragene M7-Probe — der Slice selbst
weist sie als offen aus, nichts behauptet Erfüllung ohne Beleg. Die Review-Findings F-1 bis F-10
sind am Fix-Runde-Stand behoben (F-4 mit der Closure zu lösen), F-11 bedarf keiner Aktion. Die
Sensor-Evidence ist geprüft oder im Verifier-Lauferneuert: `make gates`-Stempel deckt den exakten
aktuellen Baum, die drei CI-Läufe auf dem Tag-Commit sind success, Publish-Lauf und Release-Metadaten
stimmen, GHCR-Digests stimmen überein, der Tap-Trockenlauf wurde im Verifier-Container samt
Gegenprobe grün nachgefahren, `doc-immutable` ist grün über den Range. Der Weg zur Closure steht
frei: Closure-Notiz §7 (inkl. Lerneintrag und Risiko-Ausgänge), Roadmap-Zug mit Tag-Hash
(`e6c873798380991ff6c072766085a1e5e1f1590e`), frisches `make gates`, reiner `git mv` nach `done/`,
M7 wie oben ausgewiesen.