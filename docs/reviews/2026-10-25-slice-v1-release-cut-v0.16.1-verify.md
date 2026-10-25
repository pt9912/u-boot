# Verification-Report: slice-v1-release-cut-v0.16.1

* Rolle: Verifier (frischer Kontext), Datum 2026-10-25
* Gegenstand: `git diff f9cc5e87..HEAD` (`HEAD` = `28295c07`; Move `c238bc71`, Nachzug `cb62453c`, Start `1d092bd0`, Cut `bde533fc`, Umsetzung `34809217`, Review `53fffebf`/`9c101499`, Fix-Runde `28295c07`), Slice [`slice-v1-release-cut-v0.16.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.16.1.md), Review-Report [`2026-10-25-slice-v1-release-cut-v0.16.1`](2026-10-25-slice-v1-release-cut-v0.16.1.md)
* Frage: Wurde der Cut richtig gebaut, gemessen an Plan, DoD, Lastenheft [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) und den selbst gefahrenen Sensoren? Maßstab: [`harness/verification.md`](../../harness/verification.md).
* Nichts geändert außer diesem Report; nicht gepusht, nicht getaggt.

## Ergebnis

Der Cut ist bis zur Grenze „vor M1“ **konform**. DoD-Punkte 1 (Release-Text und Versionsstellen) erfüllt; die Punkte 2 und 3 (Tag, Kontrolle, Doku-Nachzug) und die Closure-Pflichten sind erwartungsgemäß offen und hängen an M1 bis M8. Review-Findings F-1, F-2, F-3 sind behoben. 0 Befunde mit Blockwirkung auf M1; 1 LOW (veralteter Satz im Umsetzungsvermerk), 1 Hinweis zur Lead-Breite.

## 1. Sensoren, selbst gefahren (Docker-only, Host `linux/amd64`)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | Exit 0 | `lint` 0 issues; `coverage-gate` 92,10 % gegen 90 %; `docs-check` 896 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge |
| `make ci` | Exit 0 | `govulncheck` „No vulnerabilities found.“; `image-scan` `u-boot:latest (debian 12.15)` 0, `usr/local/bin/u-boot` 0; „[ci] gates + govulncheck + image-scan green“ |
| `make doc-trace` | Exit 0 | 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0 | 9 Commits, 896 Dateien, 0 Befunde |
| `make doc-immutable RANGE="v0.16.0..HEAD"` | Exit 0 | 47 Commits, 0 Befunde (keine Accepted-ADR umgeschrieben) |
| `make build VERSION=0.16.1` | Exit 0 | `docker run --rm u-boot:latest --version` = `u-boot version 0.16.1`; OCI-Label `org.opencontainers.image.version` = `0.16.1`; Image `u-boot:latest` danach per `docker rmi` entfernt |
| M-T Tap-Trockenlauf | nachgefahren, deckungsgleich mit §3 | siehe Abschnitt 2 |
| `make test-docker`, `make fullbuild`, `make verify-depguard` | **nicht gelaufen** | `test-docker`: Auftrag (Speicher, Zeitlimit); der Produktcode-Diff ist eine Zeile (`var version`), die Docker-Pfade belegt M2 auf dem Push-Tip. `fullbuild` = `ci` + Runtime-Image-Build, beides einzeln gelaufen (`ci`, `build`). `verify-depguard`: keine Importregel berührt (Diff ohne Go-Import) |

Der Befund des Implementers (`gates` Exit 0, `ci` Exit 0, `build` Label) ist damit reproduziert, nicht nur übernommen.

## 2. M-T nachgefahren

Befehl: `docker run -i --rm --network none -v /Development/u-boot:/src:ro u-boot:test bash < <Skript im Scratchpad>`; das Repo wird nur gelesen (Skripte werden in ein Temp-Verzeichnis kopiert), `curl` und `git` sind Stellvertreter im `PATH`, die ihre Aufrufe protokollieren (`git diff --quiet` meldet immer „geändert“).

| Fall | Erwartung (Plan) | Ist |
|---|---|---|
| (1) `homebrew-formula-fill.sh v0.16.1` mit fünf Fixture-Digests | Exit 0, vier `url`-Zeilen mit `v0.16.1` | Exit 0, 4 `url`-Zeilen, 4 mit `/download/v0.16.1/` |
| (2a) `tap-formula-check.sh v0.16.1` | Exit 0 | Exit 0 |
| (2b) `tap-formula-check.sh v0.16.0` gegen dieselbe Formel | Exit 1 | Exit 1 („found 4 url lines, 0 of v0.16.0“) |
| (2c) gemischte Formel (eine URL `v0.16.0`) gegen `v0.16.1` | Exit 1 | Exit 1 („3 of v0.16.1“) |
| (3) `tap-nachzug.sh`, `TAP_TOKEN=dummy`, `TAG=v0.16.1` | `clone`, `add`, `commit -m u-boot 0.16.1`, `push` | Exit 0; Protokoll: `curl …/download/v0.16.1/u-boot.rb`, `git clone --depth 1`, `git diff --quiet`, `git add Formula/u-boot.rb`, `git commit -m u-boot 0.16.1`, `git push origin HEAD`; „updated to v0.16.1.“ |
| (3b) `tap-nachzug.sh`, `TAG=v0.16.0` gegen dieselbe Formel | bricht vor `clone` ab | Exit 1; Protokoll nur der `curl`-Aufruf, 0 `clone`, 0 `push` |

Rot gesehen: (2b), (2c), (3b). **Grenze (wie im Plan):** nicht bewiesen sind `is_stable`, der echte `curl` gegen das Release-Asset, der echte `git push` und das Token; `tap-check` (M4) bleibt offen.

## 3. Binary-Diff und Versionsstellen

* `git diff v0.16.0..HEAD --stat -- ':!docs' ':!spec' ':!*.md'`: vier Dateien, `Dockerfile` (3), `cmd/uboot/main.go` (1), `repo.mk` (3) und `internal/e2e/sandbox_devcontainer_docker_test.go` (3 Einfügungen, 2 Löschungen). Die Testdatei: gelesen, nur `//`-Zeilen (Kommentar vor `TestE2E_SandboxDevcontainer_Ollama`), Build-Tag `docker`. Die drei Versionsdateien: nur `0.16.0-dev` → `0.16.1-dev` samt Kommentar-Erwähnungen. „The code of the binary is unchanged; only its version string differs“ trägt. Kein „byte-identical“ im Text.
* `grep 0.16.0-dev` über den Baum (ohne `.git`, `.harness`): Treffer nur in Planungstext ([`slice-v1-release-cut-v0.16.0`](../plan/planning/done/slice-v1-release-cut-v0.16.0.md), dieser Slice) und Review-/Verify-Reports. Nicht in Code, Doku, Workflows.
* `0.16.0` in `README.md:212` und `README.de.md:215` ausschließlich als Tabellenzeile `v0.16.0` (Geschichte, laut Plan unberührt); `docs/maintainer/releasing.md` ohne Diff, Zeiger „zuletzt“ nennt `v0.16.0`, plan-konform bis M8.
* `0.16.1` an allen Plan-Stellen: `main.go`, `repo.mk`, `Dockerfile` (`…-dev`), README Status-Block, `VERSION=`, Tabellenzeile (beide Sprachen), Handbuch `:4`, `:82`, `:108`, `:135`, `:3207`.
* CHANGELOG: `[Unreleased]` leerer Anker, `## [0.16.1] - 2026-10-25`, Vergleichslinks `[Unreleased]: …v0.16.1...HEAD`, `[0.16.1]: …v0.16.0...v0.16.1`; `### Changed`-Eintrag wörtlich unter `[0.16.1]` (Diff).

## 4. M0: Datum

`date +%F` = `2026-10-25` = `CHANGELOG.md` (Überschrift `[0.16.1]`), Tabellenzeile beider READMEs, Handbuch-`Stand` (`:5`). **Keine Abweichung.** Der Hauptagent wiederholt M0 vor M1 und vor M5 (ein Lauf über Mitternacht verlangt einen weiteren Cut-Commit und M2 erneut).

## 5. Lead-Wortlaut gegen Lastenheft 0.24.0, Messprotokoll, CI-Befund

Lastenheft: Zeile 9 = `0.24.0`; Zeile 412 (`ollama`: Linux `amd64` mit lokaler Docker-Engine und Colima `arm64`); Zeilen 414 und 418 (Cloud-Modelle und `ollama-host` nur Linux `amd64`, lokale Engine; Colima ausdrücklich nicht); Zeile 472 (`podman build`: Docker-Hub-Images, „auf den Hosts, auf denen `podman` … startet“); Zeile 474 (Start: beide Hosts, verschachtelte User-Namespaces mit `kernel.apparmor_restrict_unprivileged_userns=0`, bei Colima in der VM; `/dev/fuse` und `/dev/net/tun`; Wert `1` → Exit-Code `11`).

| Aussage im Lead | Fundstelle | Befund |
|---|---|---|
| Lastenheft 0.24.0; Zusage-Mengen von `ollama` und `podman` in Worten | Zeilen 9, 412, 472, 474 | gedeckt |
| `podman`: nur bei verschachtelten User-Namespaces (`…=0`, bei Colima in der VM), `/dev/fuse`, `/dev/net/tun`; Wert `1` → Exit-Code `11` | Zeile 474 | gedeckt; F-1 behoben (Bedingung steht im Zusage-Absatz, nicht erst im Messsatz) |
| Cloud-Modelle und `ollama-host` nur Linux `amd64`, lokale Engine; zwei Mengen getrennt; „every other host and engine stays not promised“ | Zeilen 414, 418 | gedeckt |
| `ollama` auf `amd64`: zwei Läufe auf einem Host, ein CI-Pass (Ubuntu 24.04, 2026-10-25) | Handbuch `:2159`–`:2175`; Messprotokoll Zeile 17 (Image `ubuntu-24.04`), Zeile 34 (`--- PASS: …_Ollama`) | gedeckt; „Pass“ nicht breiter als die Logzeile |
| `podman` auf `amd64` gemessen mit `…=0`; CI ohne Gate, Runner Wert `1`, Tests übersprungen | Messprotokoll Zeilen 25 bis 28 (vier `SKIP`), Zeile 39 (Meldung zum Wert `1`) | gedeckt |
| Colima: Zusage bleibt, Probe-Image von Hand, erzeugtes Dockerfile dort nicht gemessen | Handbuch `:2171`–`:2174` | gedeckt; kein „never“ |
| kein neuer Befehl, keine neue Option; Code unverändert | Abschnitt 3 | gedeckt |

Nicht im Lead: Run-Identifier, CI-Aussage zum Push-Tip, „security“/„vulnerability“, „unified“ (nur in der Verneinung „not unified“).

**README-/Handbuch-Zeilen zum Lead (nach F-2):** Status-Block beider READMEs sagt „gemessen auf Linux `amd64` (zwei Läufe auf einem x86_64-Host, ein CI-Lauf); für Colima, arm64, beruht die Zusage auf einem von Hand nachgeschriebenen Probe-Image, und das erzeugte Dockerfile ist dort nicht gemessen“ — konsistent zu Lead und Handbuch. Tabellenzeile `v0.16.1` (EN wörtlich nach Plan, mit der Bedingung aus F-1; DE sinngleich, Name „Zusagen in Worten“) trägt Bedingung, CI-Grenze und die Colima-Grenze.

**Kennungen:** `git diff f9cc5e87..HEAD` in `README.md`, `README.de.md`, `docs/user/` enthält in den hinzugefügten Zeilen keine `LH-`/`ADR-`/`SPEC-`/Slice-Kennung und keinen Planverweis (Suche und Lesung der Diffs). `AGENTS.md` §5 Regel 6 gewahrt.

**Hinweis (kein Befund):** Der Lead sagt „promised … on Colima (arm64)“ für `podman` mit der Bedingung „inside the VM“; ob der Start dort jemals im Erzeugnis gemessen wurde, sagt der Lead nicht (Beleg-Tabelle §3: Halbsatz „Colima gemessen“ fehlt bewusst). Das ist die engere, nicht die breitere Lesart und deckt sich mit „nichts weiter als gemessen“.

## 6. Review-Findings

| Finding | Soll | Ist |
|---|---|---|
| F-1 MEDIUM (Lead/README ohne Bedingung der `podman`-Zusage) | Bedingung im Zusagesatz | **behoben** (`28295c07`, CHANGELOG-Lead und Tabellenzeile beider READMEs; gegen Zeile 472/474 gelesen) |
| F-2 MEDIUM (README-Klammer zu `ollama`/Colima) | Linux `amd64` gemessen, Colima Probe-Image | **behoben** (Status-Block beider READMEs, siehe Abschnitt 5) |
| F-3 LOW (Ablaufpunkt-Marker Roadmap) | Marker nennen ihren Ablaufpunkt | **behoben** („der Zeiger entfällt bei der Closure, M8“; Spalte Slice: „Zeile wandert bei der Closure, M8“). Merkposten für M8: beide Marker mit M1/M5 nicht vergessen |
| F-4, F-5, F-6 (INFO) | keine Änderung verlangt | unverändert; F-4 (Patch gegen Minor) bleibt Entscheidung des Projektinhabers |

## 7. Plan-vs-Code-Diff (Plan-Tabelle gegen `git diff --stat f9cc5e87..HEAD`, §3.10)

Elf Dateien im Diff vor diesem Report: `CHANGELOG.md`, `Dockerfile`, `README.md`, `README.de.md`, `cmd/uboot/main.go`, `repo.mk`, `docs/user/benutzerhandbuch.md`, `roadmap.md`, `open/README.md`, dieser Slice, Review-Report — jede mit Plan-Eintrag, jeder Plan-Eintrag mit Diff. Kein Diff, durch Lesung bestätigt: `docs/maintainer/releasing.md` (Zeiger „zuletzt“ `v0.16.0`), `spec/` (Lastenheft 0.24.0 Zeile 9), `docs/plan/adr/` (`doc-immutable` Exit 0), `.github/workflows`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `internal/` außer der geplanten Testkommentar-Datei (Diff liegt vor dem Plan-Commit). Dieser Report ist die im Plan genannte Zeile `…-verify.md`.

**LOW (V-1), Zustandsmarker §3.10:** Im Umsetzungsvermerk (Slice §3, *Ist gegen Plan*) steht, die Zeile 13 der `README.md` sei „**nicht** geändert“ (Colima-Klammer, „der Reviewer prüft“). Nach der Fix-Runde ist die Klammer geändert (Absatz *Fix-Runde nach dem Review* im selben Slice sagt es richtig). Der ältere Satz beschreibt einen überholten Zustand; bei der Closure berichtigen oder auf die Fix-Runde verweisen. Ebenso die Plan-Zeile `README.md` in §3 („der Rest des Status-Blocks bleibt“). Keine Wirkung auf Release-Text.

## 8. DoD-Abgleich

| DoD-Punkt | Stand | Beleg |
|---|---|---|
| 1 Release-Text und Versionsstellen | **erfüllt** (vorbehaltlich Datum bis M5) | Abschnitte 3 bis 5; Beleg-Tabelle §3 gelesen; Diff gegen `v0.16.0` außerhalb `docs/`, `spec/`, `*.md` nur Testdatei (Kommentar) und Versionszeilen |
| 2 Tag und Kontrolle | **offen** | `gates` und `ci` grün, `build` grün, `test-docker`-Auslassung begründet (Abschnitt 1), M-T festgehalten und nachgefahren; offen: Push, M2, M4, M5, M6 (Unterpunkte `publish`, GHCR, zwölf Assets, Release-Notes, Job `tap`) |
| 3 Doku nachgezogen (`releasing.md`-Zeiger, Roadmap-Zeile, `open/README.md`) | **offen bis M8**; `open/README.md` und Roadmap-Zeiger tragen den Ist-Stand | Abschnitt 7 |
| `make gates` grün | **erfüllt** für `HEAD` `28295c07` (vor diesem Report); erneut nach der Closure | Abschnitt 1 |
| Review-Report liegt vor | **erfüllt** | [`2026-10-25-slice-v1-release-cut-v0.16.1`](2026-10-25-slice-v1-release-cut-v0.16.1.md) |
| Verification-Report (vor M1) | **mit diesem Report erfüllt** | diese Datei |
| Closure-Notiz mit Steering-Loop-Eintrag, Register, Risiko-Ausgänge, drei Paarungen | **offen** (Closure, M8; §7 trägt noch den Platzhalter, §6 alle Ausgänge „offen bis zur Closure“) | Slice §6, §7 |

## 9. Was für M1 bis M8 bleibt

* **M1 Push** (Hauptagent, nach ausdrücklicher Freigabe): vorher M0 wiederholen, Lead im Wortlaut vorlegen (CHANGELOG `[0.16.1]`, Frage 1), Arbeitsbaum sauber, `make docs-check` lokal grün. `origin/main..HEAD` trägt 9 Commits plus den Commit dieses Reports.
* **M2** CI auf dem Push-Tip: `ci` (`gates`, `security-gates`, `image-scan`, keine Annotation „Unexpected input(s)“) und `integration-docker` mit `--- FAIL` 0, kein Zeitlimit-Abbruch, genau die vier `--- SKIP` (`NestedPodman_Build`, `…_BuildEgress`, `…_NestedPodman`, `StatusFiles_VFS`) mit der Meldung zum Wert `1`, und `--- PASS` für `…_Ollama`, `…_Ollama_PinMismatchBreaksBuild`, `TestE2E_OllamaHost_DefaultProfileResolvesHostDockerInternal`, `TestE2E_OllamaHost_SandboxWithEgressResolvesAndAddsNoHost`; Job-Log über REST-API, Abschnitt „Operating System“ (Runner-Image), Dauer gegen `timeout-minutes`. Weicht die Zählung ab: Lead, READMEs, Handbuch vor dem Tag nachziehen.
* **M4** `tap-check` grün, unmittelbar vor M5.
* **M5** Tag `v0.16.1` auf dem in M2 geprüften Commit, `git push origin v0.16.1` (nicht `--tags`), nur nach eigener Freigabe; M0 erneut.
* **M6** `publish`-Lauf (`publish`, `tap`, Paket-Smoke rpm/fedora, deb/ubuntu), `gh release view v0.16.1` (zwölf Assets, kein Prerelease/Draft), GHCR `--version` = `0.16.1` und `:0.16.1` = `:latest` (Digest), Log und Annotationen des Jobs `tap`.
* **M7** Mac: `brew install`/`upgrade`, `u-boot --version` = `0.16.1` (Projektinhaber).
* **M8** Roadmap-Zeile nach *Abgeschlossene Wellen* mit Tag-Hash und Datum, Zeiger entfernen und Ruhe-Marker zurück, `releasing.md`-Zeiger „zuletzt“, V-1 berichtigen, Closure-Notiz mit Steering-Loop-Lerneintrag, Register, Risiko-Ausgänge, erst Inhalt, dann reiner `git mv` nach `done/`.

## Verification Evidence (Kurzform)

Scope: Slice [`slice-v1-release-cut-v0.16.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.16.1.md); IDs [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../plan/adr/0017-linux-pakete-nfpm.md); Artefakte CHANGELOG, READMEs, Handbuch, drei Versionszeilen.
Carveouts: neu keine, gelöst keine, unverändert alle (kein Diff an `carveouts.md`).
Nicht ausgeführt: `make test-docker` (Speicher/Zeitlimit; Docker-Pfade in M2), `make fullbuild` (Teile einzeln gelaufen), `make verify-depguard` (keine Importregel berührt), M-T-Grenze (`is_stable`, echter `curl`/`push`, Token).
Commit / Artefakt: `HEAD` `28295c07`; lokales Image `u-boot:latest` (Version 0.16.1) nach der Prüfung entfernt.
