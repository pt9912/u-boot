# Verifikations-Report: slice-v1-release-cut-v0.19.0 — 2026-11-02

**Rolle:** Verifier (`harness/roles.md` §Verifier — „built the thing right", frischer Kontext)
**Gegenstand:** [[`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)](../plan/planning/done/slice-v1-release-cut-v0.19.0.md), Cut `3564ac31` (Tag `v0.19.0`, live), Nachzüge `ec556f96` (Roadmap), `6daf52ac` (Review-Fix-Runde), `27bf4304` (F-5-Bestand-Zeile) — HEAD `27bf4304`
**Skill:** `.harness/skills/verifier` (Rollenvertrag [`roles.md`](../../harness/roles.md), Schema [`verification.md`](../../harness/verification.md)) · **Modell:** glm-5.3-flash · **Datum:** 2026-11-02
**Basis:** Review-Report [`2026-11-02-slice-v1-release-cut-v0.19.0.md`](2026-11-02-slice-v1-release-cut-v0.19.0.md) (3 MEDIUM/3 LOW/2 INFO), DoD §2 des Slice, [`releasing.md`](../maintainer/releasing.md) §§3–8

---

## Scope

- Slice: [[`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)](../plan/planning/done/slice-v1-release-cut-v0.19.0.md) (liegt in `in-progress/`, vor der Closure)
- IDs: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Eintragstyp, S1–S11), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (Cache-Werte unter `containerEnv`), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) (Pin-Politik), [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) (Grabstein, nicht belebt), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv), [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md), [`ADR-0004`](../plan/adr/0004-ci-system.md), [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md), [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md)
- Artefakte: [`CHANGELOG.md`](../../CHANGELOG.md) Block `[0.19.0]` · Version v0.19.0 (Tag `3564ac31`) · Release `u-boot v0.19.0` (12 Assets) · GHCR `pt9912/u-boot:0.19.0` == `:latest` · Tap-Commit `72cabea` in `pt9912/homebrew-u-boot`

## Lead gegen die Lastenheft-Stände 0.30.0 und SPEC-054 (Wortlaut-Sonden, Fix-Stand)

Gelesen am Fix-Stand [`CHANGELOG.md`](../../CHANGELOG.md) Block `[0.19.0]` gegen die
Fassung von [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
(Absatz zum Eintragstyp), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)
(Cache-Absatz, Zeile 385 des Lastenhefts) und
[`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) (Abschnitte Default-Pin/Übersteuerung, Abweisungen, Cache-Werte):

- **F-1 (Profil-Grenze) behoben:** der Lead nennt jetzt die Profil-Grenze in der Zusage —
  „so that in the sandbox profile the caches survive `u-boot devcontainer up --recreate`" —
  kongruent zur Vertragszeile „Dass die Caches … das Neuanlegen des Containers übersteht, ist
  nur im Sandbox-Profil …" ([`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)).
  Kein Leser liest den Lead mehr als Überstehen des `--recreate` ohne Profil-Bedingung.
  Residuum (kein Blocker): die Lead-Not-promised-Liste nennt das Standardprofil für den
  Cache-Effekt weiterhin nicht ausdrücklich (nur „other hosts and engines"); die Grenze steht
  vollständig im Träger-Eintrag (zweiter `Added`-Bullet: „Not promised: the standard profile
  for the cache effect") und die lead-eigene Zusage ist nun gescoped — die Klasse
  „Lead-Grenzliste kürzer als Vertragsliste" bleibt beobachtenswert, der Befund F-1 ist
  im Kern ausgeräumt.
- **F-2 (Lese-Moment) behoben:** der Lead sagt jetzt „the *generation* then reads the official
  checksum source of the chosen version … and the *build* verifies the archive against it" —
  Wortlaut-kongruent zu [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  („Der Generator liest bei der Erzeugung die Prüfsummen-Quelle …"). Der Netzgriff liegt im
  Erzeugungs-Schritt, der Build verifiziert; die Träger-Eintrag-Lesart („u-boot then reads …
  and builds the override against the read checksums") deckt der Lead nun.
- **F-3 (Abweisungsliste) behoben:** der Lead enumeriert jetzt alle vier Fälle — „An unreadable
  checksum source, a version that the source does not cover, an entry without a source, and the
  combination of the external `go` feature with `go-toolchain`, each end the call with exit
  code `10`" — deckend mit den Abweisungen der
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  (Quelle nicht lesbar oder überdeckt die Version nicht · Eintrag ohne Quellen-Angabe · externes
  `go` plus `go-toolchain`); der Promised-Zähler „the three exit codes" zählt die drei
  Abweisungs-Gruppen der Vertrags-Historie („die drei Abweisungen je Exit-Code `10`") — konsistent.
- Kein Punkt des Lead zugesagt, den Lastenheft `0.30.0` nicht trägt; die Promised-Liste
  (Linux `amd64`, lokale nicht verschachtelte Engine, Testbasis Unit/Golden/`docker`-Tag) und
  die Not-promised-Liste (Arbeitsverhalten übersteuerter Werkzeuge, Lint ohne Go-Installation,
  `BUILD_MODE`/`GOFLAGS`/`GOTOOLCHAIN`, Grenze jenseits der Messung) decken die
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)-Ausweisung
  der nicht zugesagten Punkte ab. — **geprüft, ohne Befund abseits des F-1-Residuums oben.**

## DoD-Abgleich

| DoD-Punkt | Stand | Evidence |
|---|---|---|
| Release-Text und Versionsstellen: `## [0.19.0] - 2026-11-02` mit Lead („Nineteenth release."), leerem `[Unreleased]`, Compare-Links; `0.19.0-dev` an den drei Stellen samt Kommentar; beide Readmes und Handbuch auf `v0.19.0`; `Stand` = Tag-Datum | ✓ erfüllt | [`CHANGELOG.md`](../../CHANGELOG.md) Zeile 16 (Block-Kopf), Zeile 24 leerer Anker `[Unreleased]`, Compare-Zeilen 1609/1610 (`[Unreleased]: …/compare/v0.19.0...HEAD`, `[0.19.0]: …/compare/v0.18.0...v0.19.0`); `0.19.0-dev` an `cmd/uboot/main.go:45`, `repo.mk:67` (Kommentare `repo.mk:63,66`), `Dockerfile:51` (Kommentare `Dockerfile:47,49`); Versionsreste `0.17.0`/`0.18.0` in den gelesenen Pfaden: 0 Treffer abseits der Compare-Links; `README.md:10,98`, `README.de.md:11,101`; `benutzerhandbuch.md:4,5` (v0.19.0, Stand 2026-11-02 = Tag-Datum); `git rev-parse v0.19.0` = `3564ac31a15d83354349abe9e84a3ef04656bb15` |
| Vorab-Prüfungen: `make gates`, `make ci`, `make fullbuild` lokal grün; Tap-Trockenlauf mit Gegenprobe; `ci`+`integration-docker` auf dem Push-Tip grün; `tap-check` grün; danach nach Freigabe Tag gesetzt und gepusht; Kontrolle §6 | ✓ erfüllt (M7 offen, s. u.) | `.harness/state/gates-passed.head` = `27bf4304…` = HEAD und `.harness/state/gates-passed.diffsha` = eigener Neulauf von `tools/harness/working-tree-hash.sh` (`7b93809a…`) — gates-Stempel deckt den exakten Baum des letzten Commits; `ci`‑/`fullbuild`‑Läufe des Hauptagenten stehen in der Cut-Message (Vorabprüfungen grün, 2026-11-02) — remote trägt der Build: der `publish`-Lauf baute und schob das Runtime-Image; CI selbst nachgelesen: `ci` `41460180976` success, `integration-docker` `41460180946` success, `tap-check` `41463190182` success, alle auf dem Tag-Commit `3564ac31`; danach `publish` `41463337361` success mit vier Jobs (`publish (ghcr.io image push)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)`, `tap (homebrew formula)`) je success; `gh release view v0.19.0`: draft=false, prerelease=false, created `2026-11-02T14:15:16Z`, published `2026-11-02T14:36:59Z`, 12 Assets, `targetCommitish` `3564ac31` |
| Doku nachgezogen: Roadmap-Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash; Bestand-Zeilen | ✓ erfüllt | [`roadmap.md`](../plan/planning/in-progress/roadmap.md) Zeile 117 (Tag `v0.19.0` auf `3564ac31`, Datum 2026-11-02, Kennung der Zeile mit Fix `6daf52ac` — dazu V-2 unten); Bestand-Zeile des Cut-Slices in [`open/README.md`](../plan/planning/open/README.md) Zeile 98 nachgetragen in `27bf4304` (Tag-Hash, `publish`-Lauf-ID, Zustand `in-progress/`) — Review F-5 damit behoben |
| `make gates` grün nach dem letzten Commit vor dem Handoff und nach der Closure | ✓ Teil 1 / **offen Teil 2** | Stempel deckt HEAD `27bf4304` und den Inhalts-Hash des aktuellen Baums (`record-gates` läuft nur als letztes gates-Prerequisite nach grünem Lauf); der Lauf nach der Closure (§7-Notiz + Roadmap + Move) ist mit dem Closure-Commit frisch zu fahren |
| Review (kein Self-Review) und Verification mit Reports unter `docs/reviews/` | ✓ erfüllt, sobald dieser Report vorliegt | Review-Report [`2026-11-02-slice-v1-release-cut-v0.19.0.md`](2026-11-02-slice-v1-release-cut-v0.19.0.md) (committet in `6daf52ac`); dieser Verifier-Report |
| Closure-Notiz mit Lerneintrag; Register fortgeschrieben oder „keine Beobachtung"; jedes Risiko mit Ausgang; die drei Paarungen getragen | **offen — ist Closure-Schritt** | §7 des Slice trägt noch den Plan-Marker (Review F-4, s. u.); bei der Closure zu ersetzen — der Lerneintrag muss die Schwelle-Klassen des Steering-Loop-Vermerks des Review (Lead-Grenzliste: 3. Vorgang; Zustandsmarker nach Ablauf) und die Paarungen der drei Träger tragen |
| M7 (Mac `brew upgrade u-boot`, `u-boot --version`) | **offen — Aktion des Projektinhabers** | nicht gelaufen; der CI-Job `tap` (success in `41463337361`) und der Tap-Commit `72cabea` ersetzen die macOS-Sondierung nicht — kein Teil dieses Verdicts |

## Review-Findings am Fix-Stand

| ID | Prüfung nach `6daf52ac`/`27bf4304` |
|---|---|
| F-1 (MEDIUM) | behoben — Profil-Grenze in der Lead-Zusage (Wortlaut-Sonde oben); Residuum in der Not-promised-Liste ohne Blocker-Wirkung dokumentiert. |
| F-2 (MEDIUM) | behoben — Erzeugung liest, Build verifiziert; gegen [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) gelesen. |
| F-3 (LOW) | behoben — vier-Fälle-Enumeration im Lead, deckend mit den Abweisungen der Spec; `go.dev/dl`/`CHECKSUMS`-Quellen benannt. |
| F-4 (MEDIUM) | nicht aufgehoben — ist mit der Closure selbst zu lösen: §7 des Slice trägt unverändert „Noch nicht geschlossen; der Slice liegt in `open/` (Plan)" (Fundstelle Zeile 65 des Slice). Bedingung der Move-Freigabe: die Closure-Notiz ersetzt den Marker **vor** dem `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3 Fall 2). Kein Blocker für dieses Verdict — die zwei MEDIUM-Behandlungen (`v0.17.0`/`v0.18.0`) zeigen denselben Weg. |
| F-5 (LOW) | behoben — die Plan-Zeile `docs/plan/planning/open/README.md` (§3 des Slice) hat jetzt den Diff: Bestand-Zeile in `27bf4304` nachgetragen; die Message von `6daf52ac` (F-5-Absatz unzutreffend) wurde in der Message von `27bf4304` ehrlich berichtigt — der Abgleich steht damit auf dem Ist-Stand. |
| F-6 (LOW) | teils behoben, **Reste-Befund V-2** (unten): der Apostroph der Zeile ist weg, aber die Anker-Klammer trägt jetzt einen **schließenden Backtick-Fehler** — ein offener Inline-Code (ungerade Backtick-Zahl 17 in der Zeile) statt des Apostrophs; Vorbild-Zeile ist der Fix von `27bf4304` an der Bestand-Zeile von [`open/README.md`](../plan/planning/open/README.md) (20 Backticks, ausbalanciert). |
| F-7 (LOW) | teils behoben, **Reste-Befund V-3** (unten): die Roadmap-Zeile 32 trägt `v0.20.0`; die Bestand-Zeile des JDK-Entwurfs in [`open/README.md`](../plan/planning/open/README.md) Zeile 86 trägt weiterhin „ohne Welle (Vorschlag `v0.19.0`)" — die zweite der zwei genannten Fundstellen ist nicht berichtigt. |
| F-8 (INFO) | korrekt offen — die Nachzugs-Marker ([`roadmap.md`](../plan/planning/in-progress/roadmap.md) Zeilen 24 und 117, Zustand `in-progress/`, Beleg-Anker Tag-Hash und `publish` `41463337361`) tragen Zustand und Beleg und laufen erst mit der Closure ab (Run-ID selbst gelesen: `publish` `41463337361` success). |
| F-9 (INFO) | keine Aktion nötig — [`releasing.md`](../maintainer/releasing.md) §5 verlangt keine Anhebung der Handbuch-Version; `Stand` 2026-11-02 = Tag-Datum bestätigt. |

## Plan-vs-Code

Abgeglichen für den Stand nach der letzten Berichtigung (`27bf4304`) gegen
`git diff --stat 9639db21..HEAD` (11 Dateien):

| Datei im Range | Plan-Tabelle §3 | Befund |
|---|---|---|
| `CHANGELOG.md` | Lead-Zeile | Deckung (`6daf52ac` Lead-Fix gehört zur Lead-Zeile) |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | Versionszeilen-Zeile | Deckung (`main.go:45`, `repo.mk:67`, `Dockerfile:51`) |
| `README.md`, `README.de.md` | Status-Zeile und `VERSION=`-Beispiel | Deckung |
| `docs/user/benutzerhandbuch.md` | Kopf, Stand, Versionsbeispiele | Deckung (`:4,5`) |
| `docs/plan/planning/in-progress/roadmap.md` | Zeile nach dem Tag mit Tag-Hash | Deckung (`ec556f96`/`6daf52ac`) |
| `docs/plan/planning/open/README.md` | Bestand-Zeile für diesen Slice | Deckung seit `27bf4304` (vorher F-5) |
| [`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md) | „dieser Slice — neu/Moves" | Deckung (Anlage `b4d56997`) |
| [`2026-11-02-slice-v1-release-cut-v0.19.0`](2026-11-02-slice-v1-release-cut-v0.19.0.md) | keine Plan-Zeile | Review-Produkt (Pflicht DoD-Zeile „Review … mit Reports"); kein Plan-Stand kann es vorwegnehmen — der Abgleich nennt es, wie hier geschehen; das Vorbild `v0.18.0` hat dasselbe in der Closure-Notiz ausgewiesen; analog für diesen Verifier-Report nach der Closure |
| `spec/`, `docs/plan/adr/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | „kein Diff" | Bestätigt — im ganzen Range `9639db21..HEAD` nicht berührt |

## Sensoren

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass (Stempel deckt den exakten Baum; kein Neulauf nötig) | Stempel `.harness/state/gates-passed.head` = `27bf43048bb513ec16783625ece38c28ddf28579` = `git rev-parse HEAD`; `.harness/state/gates-passed.diffsha` = `7b93809a48db21c5a21cf8930d23493db86862ff3d51ddff0d16068421bb1a83` = eigener Neulauf von `tools/harness/working-tree-hash.sh`; [`record-gates.sh`](../../tools/harness/record-gates.sh) schreibt den Stempel nur als letztes gates-Prerequisite nach grünem Lauf |
| `make test-docker` | not run | kein codeberührender Diff im Verifikations-Umfang: die Nachzüge `ec556f96`/`6daf52ac`/`27bf4304` berühren nur Lead-Text und Planning-Doku; der Cut berührt nur Versionszeilen. Der CI-Lauf `integration-docker` auf dem Tag-Commit (`41460180946`, success) trägt die Docker-Tests — selbst nachgelesen |
| `make verify-depguard` | not run | keine Importregel berührt — Produkt-Diff nur Versionszeilen (`cmd/uboot/main.go`), `repo.mk`, `Dockerfile` |
| CI nachgelesen (`gh run list`/`gh run view`) | pass | Tag-Commit `3564ac31`: `ci` `41460180976` success, `integration-docker` `41460180946` success, `tap-check` `41463190182` success, `publish` `41463337361` success (Jobs s. o.); Nachzugs-Head: [Nachtrag](#nachtrag-ci-auf-dem-nachzugs-head) — beide Läufe success |
| Release selbst | pass | `gh release view v0.19.0`: 12 Assets (Binaries ×6 Plattformen, `SHA256SUMS`, `.rpm` ×2, `.deb` ×2, `u-boot.rb`), draft=false, prerelease=false, published `2026-11-02T14:36:59Z` |
| Tap | pass | `gh api repos/pt9912/homebrew-u-boot/commits`: `72cabea` „u-boot 0.19.0" (Vorgänger `b52a0bc` „u-boot 0.18.0"); Formel `Formula/u-boot.rb` trägt die vier Download-URLs des Tags `v0.19.0`; drei Formel-Prüfsummen (`darwin-arm64`, `darwin-amd64`, `linux-arm64`) gegen das Release-Asset `SHA256SUMS` abgeglichen — identisch |
| GHCR | pass | Manifest-Digests `ghcr.io/pt9912/u-boot:0.19.0` == `:latest` == `sha256:814b94357e824075d82059547aac3f03edb35cbd290e5987bfd50a6cd1bfe43c` (Registry-API, `docker-content-digest`) |
| Remote-Image-Gegenprobe | pass | `docker run --rm ghcr.io/pt9912/u-boot:0.19.0 --version` → `u-boot version 0.19.0` |

## Befunde des Verifier-Laufs

| ID | Kategorie | Befund | Fundstelle |
|---|---|---|---|
| V-1 | INFO | Die Nachzugs-Commits `ec556f96`, `6daf52ac`, `27bf4304` sind mit einem Push auf `main` gelangt (Run-Gründung 2026-11-02T14:55:02Z); zum Ablesezeitpunkt liefen `ci` und `integration-docker` auf dem Head `27bf4304` — Endstand im [Nachtrag](#nachtrag-ci-auf-dem-nachzugs-head) unten: beide success. Die Nachzüge berühren Lead-Text und Planning-Doku, keinen Produktcode und kein Gate; die grünen Tag-Läufe tragen die Verifikation des Release-Inhalts. | `gh run list --commit 27bf4304…` |
| V-2 | LOW | Der F-6-Fix an der Roadmap-Zeile hat den Apostroph entfernt, aber die Anker-Klammer trägt jetzt einen offenen Inline-Code: `` [`slice-v1-release-cut-v0.19.0](…)` — das öffnende Backtick vor der Kennung bleibt ohne schließendes Partner-Backtick (ungerade Backtick-Zahl 17 in der Zeile); der Code-Span schluckt bis zum nächsten Backtick der Zeile. Die Kennung selbst ist jetzt textlich korrekt und der Link-Pfad trifft; der Fehler ist der letzte Rest desselben Malformed-Characters, den F-6 benannt hat — ein Fix-Zeichen (schließendes Backtick oder Löschung des öffnenden). `docs-check`-grün beweist, dass kein Sensor Code-Span-Balance liest. | [`roadmap.md`](../plan/planning/in-progress/roadmap.md) Zeile 117 |
| V-3 | LOW | F-7 ist nur zur Hälfte berichtigt: die Roadmap-Zeile 32 trägt `v0.20.0`, aber die Fundstelle, die der Review mit nannte — die Zeile „ohne Welle (Vorschlag `v0.19.0`)" des JDK-Entwurfs — steht unverändert in [`open/README.md`](../plan/planning/open/README.md) Zeile 86, obwohl der Tag `v0.19.0` jetzt gesetzt ist und den Entwurf nicht trägt (Vorschlagswert verbraucht). Die Zeile ist beim nächsten Zuordnungs-Commit auf die Folgerversion zu berichtigen; da sie kein aktiver Plan-Beleg des Cuts ist, ist es kein Blocker. | [`open/README.md`](../plan/planning/open/README.md) Zeile 86; [`roadmap.md`](../plan/planning/in-progress/roadmap.md) Zeile 32 (Trigger-Zelle nennt weiterhin `v0.18.0 veröffentlicht` — mit dem `v0.20.0`-Label zu prüfen, aber wahr und nicht abgelaufen) |
| V-4 | INFO | Plan-Tabelle §3 trägt keine Zeile für die beiden Reports unter `docs/reviews/` (Review-Report `6daf52ac`, dieser nach der Closure); beide sind Review-Pflicht-Produkte (DoD-Zeile „Review … mit Reports"), der Vorbild-Abgleich von `v0.18.0` nannte dieselbe Ausnahme in der Closure-Notiz — dieselbe Abgleich-Zeile gehört in §7. | [[`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)](../plan/planning/done/slice-v1-release-cut-v0.19.0.md) §3 |

## Traceability

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) | CHANGELOG Block `[0.19.0]`, erster `Added`-Eintrag (Eintragstyp samt Prüfsummen, Übersteuerung, drei Abweisungen, Ausweisung unter `--quiet`) gegen die S1–S11-Fassung gelesen — Lastenheft 0.30.0; Release als Lieferung |
| [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) | zweiter `Added`-Eintrag (Cache-Werte) gegen die Zeile 385 des Lastenhefts; Lead post-F-1 trägt die Profil-Grenze |
| [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette) | dritter `Added`-Eintrag (eine Quelle, vier Stellen, Test) und die Versionszeilen des Cuts (`main.go:45`, `repo.mk:67`, `Dockerfile:51`) |
| [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) | Lead post-F-2/F-3 gegen die Abschnitte Übersteuerung/Abweisungen gelesen |
| [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) | Hybrid-Pin als Träger des Blocks `Unreleased` → `[0.19.0]`; Pin-Politik-Eintrag |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | `tap-check` `41463190182` success, `publish`-Job `tap` success, Tap-Commit `72cabea`, Formel-Prüfsummen gegen `SHA256SUMS` identisch |
| [`ADR-0007`](../plan/adr/0007-distributionswege-ghcr.md) | Release `u-boot v0.19.0` (12 Assets), `publish` `41463337361` mit Image-Push-Job, `docker run --version` = `0.19.0`, Digest-Gleichheit |
| [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md) | wellenloser Cut (keine Welle-Datei); Roadmap-Zeile `v0.19.0` unter `Abgeschlossene Wellen` mit dem Format der MR-003-Schnittstelle; Vorschlagswerte `v0.20.0` berichtigt (teilweise, V-3) |

## Carveouts

- Neu: none
- Gelöst: none
- Unverändert: Cut und Nachzüge senken keine Schwelle, berühren kein Gate, keinen Pin-Wert außer der Versionszeile des Releases; die „kein Diff"-Zeilen der Plan-Tabelle sind im ganzen Range bestätigt.

## Nicht ausgeführt

- `make gates`-Neulauf — der Stempel deckt HEAD **und** den Inhalts-Hash des aktuellen Baums; Mechanismus von [`record-gates.sh`](../../tools/harness/record-gates.sh) gelesen. Ein frischer Lauf gehört zum Closure-Commit.
- `make test-docker` — kein codeberührender Diff in den geprüften Commits; der CI-Lauf `integration-docker` auf dem Tag-Commit ist der Nachweis und wurde selbst nachgelesen.
- `make verify-depguard` — keine Importregel berührt.
- `make ci`/`make fullbuild` erneut lokal — der Tag ist live, die remote-Läufe wurden gelesen; ein Zweitlauf brächte nichts Neues.
- M7 (macOS-Probe des Tap) — Aktion des Projektinhabers.

## Rest-Risiken

1. **M7 offener Punkt:** `brew upgrade u-boot` und `u-boot --version` auf einem Mac sind nicht gelaufen; der CI-Job `tap` und der Tap-Commit `72cabea` ersetzen das nicht.
2. **Closure-Schritte offen** (bewusst, in dieser Reihenfolge zu fahren): §7-Notiz (löst F-4) mit Lerneintrag — Schwelle-Klassen des Review-Vermerks tragen („Lead-Grenzliste kürzer als Vertragsliste", 3. Vorgang; „Zustandsmarker nach Ablauf"; das F-6/V-2-Residuum als neuer Fund) —, Risiko-Ausgänge §6, Plan-Abgleich-Vorbehalt V-4 in die Closure-Notiz, Berichtigung V-2/V-3 **vor** dem Move, Roadmap-Marker des F-8 auf den Ist-Stand, frisches `make gates`, reiner `git mv` nach `done/`, Ruheort-Berichtigung des Roadmap-Zeigers (Auftrag aus der Message von `ec556f96`).
3. **Der Lead bleibt der gelesenste Text des Release:** die Grenzlisten der `Added`-Einträge tragen die Vertragsaussagen vollständig; künftige Träger-Änderungen brauchen die Lead-Spiegelung (Klasse des Review-Vermerks, nun drei Vorgänge).

## Verdikt

**Der Release-Schnitt `v0.19.0` ist DoD-Konform verifiziert — closure-fähig nach Berichtigung der zwei Residuen.**
Jede DoD-Zeile trägt einen Beleg; die offenen Punkte sind genau die Closure-Schritte selbst (§7-Notiz inkl. F-4, Risiko-Ausgänge, Register, fresh `make gates`, reiner `git mv`) und die vom Projektinhaber getragene M7-Probe. Die Fix-Runde hat die drei Lead-Befunde (F-1/F-2/F-3) wortlaut-kongruent zur Vertragslage gelöst und F-5 vollständig; offene Reste sind **V-2** (Roadmap-Zeile 117, ein Zeichen) und **V-3** (`open/README.md` Zeile 86, Vorschlagswert) — beide berichtigt sein müssen, bevor der Slice den `done/`-Move fährt, da beide Fundstellen als Beleg-Anker des geschlossenen Stands gelesen werden. Alle Release-Nachweise sind im Verifier-Lauf selbst nachgelesen (Run-IDs, Release-Metadaten, Tap-Commit mit Prüfsummen-Abgleich, GHCR-Digests, Remote-`--version`), der Gate-Stempel deckt den exakten Baum von HEAD, und auch die CI-Läufe des Nachzugs-Heads sind success (Nachtrag unten). Der Lead sagt an keiner Stelle mehr, als Lastenheft `0.30.0` und [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) tragen.

## Nachtrag: CI auf dem Nachzugs-Head

Die Nachzugs-Commits sind mit einem Push auf `main` gelangt; die beiden Workflow-Läufe auf
dem Head `27bf4304` (Gründung 2026-11-02T14:55:02Z) sind im Verifier-Lauf bis zum Ende
mitgelesen und grün:

- `ci` `41466486010` — success (2026-11-02T14:55:02Z → 14:57:20Z)
- `integration-docker` `41466485866` — success (2026-11-02T14:55:02Z → 15:13:57Z, ~18:55 min, gegen die gemessene Marge 19:02 min von 30)

Damit ist auch der Nachzugs-Stand vollständig durch CI gedeckt; Rest-Risiko der laufenden
Läufe (V-1) entfällt. Für die Closure gilt unverändert: frische Läufe auf dem
Closure-Commit selbst (`make gates` als Stempel, CI auf dem Push).