# Slice V1: Release-Cut `v0.20.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv` (Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine).

**Welle:** keine Welle-Datei — Minor-Release `v0.20.0` von `main` nach
`v0.19.0` (Setzung [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md):
Release-Cuts sind keine Wellen). Der Release trägt, was seit dem Tag
`v0.19.0` (`3564ac31`) unter `## [Unreleased]` steht: der Katalog-Schlüssel
`java-toolchain` — der dritte Schlüssel des Eintragstyps „gepinnte
Werkzeugkette" ([`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md),
Temurin LTS 21, Same-Origin-Prüfsummen-Quelle, Konflikt-Abweisung external
`java` + `java-toolchain` je Exit-Code `10`), Lastenheft `0.31.0` — getragen
vom CR [`slice-cr-devcontainer-java-toolchain`](../done/slice-cr-devcontainer-java-toolchain.md)
(gelesen: die CR-Closure liegt **nach** dem Tag `v0.19.0`, im Diff
`v0.19.0..HEAD`) und der Umsetzung
[`slice-v1-katalog-feature-jdk-gepinnt`](../done/slice-v1-katalog-feature-jdk-gepinnt.md)
(geschlossen 2026-11-03 mit CI-Nachweis). Der Harness-Vorgang
[`slice-harness-carveouts-baseline-form`](../done/slice-harness-carveouts-baseline-form.md)
berührt keinen öffentlichen Vertrag und trägt keinen `[Unreleased]`-Eintrag.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
[`ADR-0004`](../../adr/0004-ci-system.md), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md),
[`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0041`](../../adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md).
Ablauf (einzig normativ): [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) §§3 bis 8 —
Form aus der vendored Vorlage
[`slice.template.md`](../../../../.harness/baseline/v6.18.1/templates/docs/plan/planning/slice.template.md).

**Verantwortlich:** pt9912 (Implementer; Umsetzung und Gate-Läufe durch den Hauptagenten im Auftrag)

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-03.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.20.0` liefert den Block `[0.20.0]` des CHANGELOG als GHCR-Image,
Binaries, Pakete und Homebrew-Formel. **Push und Tag führt der Hauptagent nach
Freigabe des Projektinhabers (Muster `v0.19.0`: Push nach grünen
Vorabprüfungen, Stand-Bericht vor dem Tag; die Freigabe zur Vorab-Phase steht
aus und wird beim M-Start eingeholt).**

**Nummer:** Minor (Zusage-Menge wächst: dritter Schlüssel des
Eintragstyps „gepinnte Werkzeugkette" samt Konflikt-Abweisung), keine
*Changed*-Einträge (das Verhalten bestehender Kommandos ändert sich nicht;
allein der Katalog wächst). Kein Vorab-Tag außer auf Wunsch.

**Ausdrücklich NICHT in diesem Slice:** Produktcode, `spec/`, ADRs (der Cut
ändert nur Release-Dokumente und drei Versionszeilen) · `make test-docker`
lokal (der CI-Lauf `integration-docker` auf dem Push-Tip ist der Nachweis,
Muster [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Docker-Tests) ·
Hebung der Pins (kein Anlass: [`LH-FA-DEV-011`](../../../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette)
bleibt bei den zwei Go-Schlüsseln, der Temurin-Default-Pin steht im Template
und hebt nur mit einem u-boot-Release, nicht mit diesem Cut selbst) ·
Tap-Formel-Inhalte über das vom Cut getragene Heben der Version hinaus.

**Vor dem Start zu klären (M-Start):** jeder Eintrag unter
`## [Unreleased]` trägt einen Beleg und einen Commit im Diff
`v0.19.0..HEAD` (gelesen: der java-toolchain-Eintrag, Umsetzung
`9faa3137`; keine weiteren Träger-Commits unter `[Unreleased]`).

## 2. Definition of Done

- [x] **Release-Text und Versionsstellen:** `CHANGELOG.md` trägt
  `## [0.20.0] - <Tag-Datum>` mit Lead (Stil `[0.19.0]`: „Twentieth
  release."), leerem Anker `[Unreleased]` und Compare-Links; Versionsstrings
  `0.20.0-dev` an den drei Stellen (`cmd/uboot/main.go`, `repo.mk`,
  `Dockerfile` — Fundstellen dieses Tages: `cmd/uboot/main.go` Zeile 45,
  `repo.mk` Zeilen 63, 66, 67, `Dockerfile` Zeilen 47, 49, 51; der
  Implementer gleicht mit `grep -rIn '0\.19\.0-dev'` nach); beide READMEs und
  das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen
  `v0.20.0`; Handbuch-`Stand` gleich dem Tag-Datum. **Beleg (Closure):** Schnitt `90f2a1d2` — die `[Unreleased]`-Wanderung ist byte-identisch gegen den Stand vor dem Schnitt (`git show 3bc9757f:CHANGELOG.md`, nur Header/Lead/Compare-Links, plus das F-4-Addendum der Berichtigung); die drei Versionsstellen samt Kommentar-Nennungen vom Verifier komplett gelesen; `grep -rIn '0\.19\.0-dev'` trägt nur Historie-Treffer; READMEs und Handbuch auf `v0.20.0` ohne Kennungen im Nutzertext; M0-Dreiklang 2026-11-03 (`git log -1 --format=%cs v0.20.0` = CHANGELOG = Handbuch-`Stand`).
- [x] **Vorab-Prüfungen:** `make gates`, `make ci`, `make fullbuild` lokal
  grün; Trockenlauf des Tap-Pfads nach `releasing.md` §4 (mit Gegenprobe auf
  den Vorgänger-Tag); `ci` und `integration-docker` auf dem Push-Tip grün;
  `tap-check` grün; **danach nur nach Freigabe**: Tag `v0.20.0` gesetzt und
  gepusht; Kontrolle nach `releasing.md` §6. **Beleg (Closure):** `make gates`/`make ci`/`make fullbuild` je Exit 0 am Schnitt-Stand (Stempel `90f2a1d2`); CI am Push-Tip: [`ci`] 41581784136 und [`integration-docker`] 41581784179 je **success**; [`tap-check`] 41581828207 **success**; `make build VERSION=0.20.0` meldet `u-boot version 0.20.0` (CLI) und OCI-Label `0.20.0`; Trockenlauf siehe §7 (Protokoll mit Befehl und Ausgabe, Verifier V-3); Kontrolle §6: [`publish`] 41583702257 **success** mit den vier Jobs (inklusive Job `tap`, Log `tap-nachzug: pt9912/homebrew-u-boot updated to v0.20.0.`), Release 12 Assets (6 Binaries, 2 `.deb`, 2 `.rpm`, `SHA256SUMS`, `u-boot.rb`, kein Draft/Prerelease), `docker run --rm ghcr.io/pt9912/u-boot:0.20.0 --version` = `0.20.0`, `:latest`-Pull-Digest `sha256:370bdfbe1b993bafebe5994c163df2826a9b323e170012f221b660ffbef058b5` gleich `:0.20.0`; M7 bestätigt (§7).
- [x] **Doku nachgezogen:** Roadmap (Zeile nach dem Tag *Abgeschlossene
  Wellen* mit Tag-Hash und Datum), Bestand-Zeilen (`open/README.md`, die
  Zeile dieses Slices, Muster der Zeile `slice-v1-release-cut-v0.19.0`). **Beleg (Closure):** die *Abgeschlossene-Wellen*-Zeile und die „geliefert"-Form der Bestand-Zeile trägt der Move-Runden-Commit (Reihenfolge [`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2 — Inhalt vor dem `git mv`); die Nächste-Wellen-Zeile ist entzeitlicht (Fix-Runde `f652ed4f`, V-1).
- [x] `make gates` grün nach dem letzten Commit vor dem Handoff und nach der
  Closure. **Beleg (Closure):** Stempel `90f2a1d2` (Schnitt) und `0a7f7e6b` (Report-/Vorschlag-Commits je eigener Lauf); die finalen Läufe folgen dem Closure- und Move-Zug (Verifikations-Zeile unten).
- [x] Review (kein Self-Review) und Verification mit Reports unter
  `docs/reviews/`; Closure-Notiz mit Lerneintrag; Register fortgeschrieben
  oder „keine Beobachtung"; jedes Risiko mit Ausgang; die drei Paarungen
  getragen. **Beleg (Closure):** [Review-Report](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0.md) (0 HIGH/2 MEDIUM/2 LOW/1 INFO, F-2/F-4 berichtigt, F-1/F-3 durch die Closure, V-2-Fundstellen-Fix der Berichtigung) · [Verify-Report](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0-verify.md) (1 MEDIUM/3 LOW/2 INFO, V-3 löst §7-Protokoll) · Register und Risiko-Ausgänge in §7.

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `CHANGELOG.md` | update | `[Unreleased]`-Inhalt unter `## [0.20.0] - <Tag-Datum>` mit Lead — Lead nennt die Nutzer-Wirkung samt Grenzen des Lastenheft-Stands `0.31.0` (Zusage-Menge Linux `amd64` mit lokaler, nicht verschachtelter Engine, JDK-Behaviour in einer übersteuerten Version nicht versprochen, Brew-Befund ohne zweite Installationsform); Fuß: `[Unreleased]: …/compare/v0.20.0...HEAD`, neue Zeile `[0.20.0]: …/compare/v0.19.0...v0.20.0` |
| `cmd/uboot/main.go`, `repo.mk`, `Dockerfile` | update | `0.19.0-dev` → `0.20.0-dev` je Stelle samt Kommentar (Fundstellen oben; Nachgleich per `grep -rIn '0\.19\.0-dev'`) |
| `README.md`, `README.de.md` | update | Status-Zeile und Versionsbeispiele; keine Kennungen, kein „released" |
| `docs/user/benutzerhandbuch.md` | update | Kopf `Software-Version` und `Stand`, Versionsbeispiele, Gültigkeitsbereich |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile nach dem Tag *Abgeschlossene Wellen* mit Tag-Hash und Datum; Zustandsmarker („Vorschlag Release"-Zeile der `Nächste Wellen`-Tabelle) nach dem Tag auf den Ist-Stand (Muster-Marker-Klasse, vier Auftritte) |
| `docs/plan/planning/open/README.md` | update | Bestand-Zeile für diesen Slice (mit dem Plan-Commit, nicht erst beim Start — M-F-5-Lektion `v0.19.0`) |
| dieser Slice | neu/Moves | `make slice-mv` je Lifecycle-Schritt; nach dem Tag nach `done/` |
| `docs/reviews/` — [Review-Report](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0.md) (committet `0a7f7e6b`, Link-Berichtigung samt V-2-Fundstellen-Fix in der Berichtigungs-Runde) und [Verify-Report](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0-verify.md) (committet mit der Berichtigungs-Runde) | neu | die Berichte der zwei Rollen — Review nach dem Schnitt (kein Self-Review), Verification am Tag-Stand (frischer Kontext); **Nachtrag nach Verify-Report V-5** — die Berichts-Dateien fehlen als §3-Zeilen (vierte Auftritts-Runde der Familie [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)) |
| `spec/`, `docs/plan/adr/`, `templates/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung (`AGENTS.md` §3.6); der Temurin-Default-Pin bleibt bei `21.0.13+9` |

**Reihenfolge der Maintainer-Aktionen** (aus `releasing.md` §§4 bis 6, Muster
`slice-v1-release-cut-v0.19.0`): M0 Datum-Abgleich vor jedem Tag · M-T
Trockenlauf des Tap-Pfads · M1 Push von `main` (Freigabe erteilt) · M2
CI-Läufe des Push-Tips lesen · M4 `gh workflow run tap-check` · M5 Tag
`v0.20.0` gesetzt und gepusht (nur nach Freigabe und Stand-Bericht) · M6
Kontrolle §6 · M7 Mac-Probe (Projektinhaber).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): der trägende Umsetzungs-Slice
[`slice-v1-katalog-feature-jdk-gepinnt`](../done/slice-v1-katalog-feature-jdk-gepinnt.md)
liegt in `done/` (2026-11-03 geschlossen); `Verantwortlich:` gesetzt; die CI-Läufe
des Diffs sind grün oder der Cut wartet (Push-Tip `7bce82b5`:
`integration-docker` 41576482662 und `ci` 41576482888 je success; der
Docs-Zug nach dem Tag des Cuts trägt keinen Code-Diff).

## 5. Closure-Trigger

DoD vollständig; Tag `v0.20.0` gesetzt und gepusht; `publish` und Job `tap`
kontrolliert; Closure-Notiz mit Tag-Hash, Lauf-IDs und Lerneintrag; der Slice
geht **erst nach dem Tag** nach `done/` (`AGENTS.md` §3.3 Fall 2).

## 6. Risiken und offene Punkte

1. **Der Lead sagt mehr als die Beleg-Tabelle** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
   — die Lead-Grenzlisten-Klasse steht auf der Drei-Schwelle aus `v0.18.0`
   (`F-2`), `v0.19.0` (`F-1`) und ihrem Vorgänger; Lehre des `v0.19.0`-Cuts:
   der Lead wird gegen die Vertragsliste (Lastenheft `0.31.0`, [`SPEC-054`](../../../../spec/spezifikation.md))
   gelesen, Sätze ohne Beleg-Träger stehen nicht im Lead). — **Ausgang:** bei Closure.
2. **Das Datum weicht vom Tag-Datum ab** (M0). — **Ausgang:** bei Closure.
3. **`integration-docker` am Push-Tip rot oder am Zeitlimit** (F-9-Klasse —
   gemessene Marge des JDK-Slices: 41576482662; transienter externer
   Netzfehler möglich, [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md)).
   — **Ausgang:** bei Closure.
4. **Job `tap` scheitert** — `tap-check` unmittelbar vor M5, Trockenlauf M-T. —
   **Ausgang:** bei Closure.
5. **Zustandsmarker nach dem Tag nicht auf den Ist-Stand** (Roadmap, `open/README.md`,
   Entwürfe — vierter und fünfter Auftritt der Marker-Klasse in den Cuts
   `v0.18.0` und `v0.19.0`). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

### Closure (Planner, 2026-11-03, nach Review-Fix `f652ed4f`, Verifier-Berichtigung und Tag `v0.20.0`)

**Evidence:** Review-Report [`2026-11-03-slice-v1-release-cut-v0.20.0`](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0.md) (0 HIGH/2 MEDIUM/2 LOW/1 INFO; `[Unreleased]`-Wanderung byte-identisch, READMEs/Handbuch kennungsfrei, Move-Disziplin §3.3) · Verify-Report [`…-verify`](../../../reviews/2026-11-03-slice-v1-release-cut-v0.20.0-verify.md) (frischer Kontext: DoD-Zeilen 1 und 5 belegt erfüllt, Release-Kette je Läufen selbst gelesen; **V-3** Trockenlauf-Beleg, **V-2** Fundstellen-Vertauschung des Review-Reports, **V-5** Berichts-Dateien ohne §3-Zeile, **V-1** verbrauchte Label-Version, **V-4/V-6** M7-Zitat und Push-Defizit — V-2/V-5/V-1 in der Berichtigungs-Runde berichtigt).

**Läufe (M-Start bis M7):** Vorab am Schnitt-Stand `90f2a1d2`: `make gates` (Exit 0, Stempel `90f2a1d2`), `make ci` (grün, Trivy „Clean"), `make fullbuild` (grün); `make build VERSION=0.20.0` → `u-boot version 0.20.0` (CLI) und OCI-Label `0.20.0`. **Tap-Trockenlauf nach `releasing.md` §4** (Container `u-boot:test`, `--network none`, Arbeitsbaum nur-lesen, Stubs für `curl`/`git` aus der Umgebung, Fixture aus `sumsFull` von `cmd/uboot/homebrew_scripts_test.go`; Befehl und Ausgabe, Report V-3): (1) `scripts/homebrew-formula-fill.sh v0.20.0 /fixture/SHA256SUMS-fixture scripts/homebrew-formula.rb.tmpl /out/u-boot.rb` → *„written (tag v0.20.0, four platform digests filled)"*; `scripts/tap-formula-check.sh v0.20.0 /out/u-boot.rb` → **Exit 0**; Gegenprobe mit `v0.19.0` → **Exit 1** (erwartet). (2) `TAP_TOKEN=dummy TAG=v0.20.0 scripts/tap-nachzug.sh` mit Stub-`curl` (schreibt die Fixture-Formel an `-o`) und Stub-`git` (protokolliert) → *„tap-nachzug: pt9912/homebrew-u-boot updated to v0.20.0."*; Protokoll zeigt `git clone --depth 1 …`, `git add Formula/u-boot.rb`, `git commit -m u-boot 0.20.0`, `git push origin HEAD`; Gegenprobe `TAG=v0.19.0` → bricht mit `tap-formula-check`-Fehl (Exit 1, vier URL-Zeilen tragen `0` von `v0.19.0`) **vor `clone`** ab. CI am Push-Tip: [`ci` 41581784136](https://github.com/pt9912/u-boot/actions/runs/41581784136) und [`integration-docker` 41581784179](https://github.com/pt9912/u-boot/actions/runs/41581784179) je **success**; [`tap-check` 41581828207](https://github.com/pt9912/u-boot/actions/runs/41581828207) **success**. **M5 Tag:** Freigabe des Projektinhabers 2026-11-03 („Push + Tag freigeben" im Stand-Bericht); M0-Datumscheck `date +%F` = 2026-11-03 = CHANGELOG = Handbuch-`Stand` = `git log -1 --format=%cs v0.20.0`; Tag `v0.20.0` auf `90f2a1d2`, gepusht. **Kontrolle §6:** [`publish` 41583702257](https://github.com/pt9912/u-boot/actions/runs/41583702257) **success** (Jobs `publish (ghcr.io image push)`, `package smoke (rpm, fedora)`, `package smoke (deb, ubuntu)`, `tap (homebrew formula)` je success); Release [v0.20.0](https://github.com/pt9912/u-boot/releases/tag/v0.20.0) mit 12 Assets (6 Binaries, 2 `.deb`, 2 `.rpm`, `SHA256SUMS`, `u-boot.rb`), kein Draft/Prerelease, veröffentlicht `2026-11-03T05:26:34Z`; Job-`tap`-Log trägt `tap-nachzug: pt9912/homebrew-u-boot updated to v0.20.0.` (die `::warning`-Zeilen im Log sind Schritt-Echo, nicht der gelaufene Zweig — Verifier-Nachzug); `docker run --rm ghcr.io/pt9912/u-boot:0.20.0 --version` = `u-boot version 0.20.0`; Pull-Digest `sha256:370bdfbe1b993bafebe5994c163df2826a9b323e170012f221b660ffbef058b5` je gleich für `:0.20.0` und `:latest`. **M7 (Mac-Probe, Projektinhaber 2026-11-03):** *„auf dem mac läuft es auch"* — zusätzlich vom Projektinhaber gemeldet (`u-boot --version` am Mac).

**Offengelegte Punkte:** (a) Der **getaggte Lead** `v0.20.0` trägt die engere Liste der Abweisungen (Review F-4 — der Fall „Eintrag des Typs ohne Prüfsummen-Quelle" fehlt im Lead, der Block-Eintrag trägt alle drei); der Tag wird nicht verschoben — die Berichtigung steht seit `f652ed4f` auf `main` (Muster der `v0.19.0`-Lead-Nachzug-Runde), der getaggte Block bleibt im Tag-Stand. (b) Die **Bestand-Zeile** `open/README.md` trug „nur Plan" bis zum Move (Review F-3, Marker-Klasse); die „geliefert"-Form trägt der Move-Runden-Commit. (c) Der §7-Marker („Noch nicht geschlossen") ist mit diesem Block abgelöst (Review F-1, dieselbe Klasse).

**Risiken aus §6 — je Ausgang:**
1. *Der Lead sagt mehr als die Beleg-Tabelle* — **eingetreten (abgemildert):** der Lead trug eine engere statt einer weiteren Zusage (F-4); für die Decken-Regel unproblematisch, gemeldet und auf `main` berichtigt; die Grenzlisten-Klasse ist im Verifier-Nachzug gegen Lastenheft 0.31.0 und [`SPEC-054`](../../../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) sauber.
2. *Das Datum weicht vom Tag-Datum ab* — **entfallen:** M0-Dreiklang je Läufen oben.
3. *`integration-docker` am Push-Tip rot oder am Zeitlimit* — **entfallen:** 41581784179 **success**; kein transienter Fall ([`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md) — kein neuer Beleg).
4. *Job `tap` scheitert* — **entfallen:** Trockenlauf §4 (Protokoll oben) und der reale Lauf im `publish`-Run, Log „updated to v0.20.0.".
5. *Zustandsmarker nach dem Tag nicht auf den Ist-Stand* — **eingetreten:** Roadmap-Zeile, Bestand-Zeile und §7-Marker (F-1/F-2/F-3, V-1) — alle in den Berichtigungs-Runden berichtigt; **dritte Auftritts-Generation der Marker-Klasse in den Release-Cuts in Folge** (nach `v0.18.0` und `v0.19.0`), Register unten; Sensor-Kandidat verbleibt beim Projektinhaber.

**Was ging anders als geplant:** (a) Die Marker-Klasse kam zum **dritten Cut in Folge** wieder (F-1/F-3) — trotz der F-5-Lehre des `v0.19.0`-Cuts (Bestand-Zeile mit dem Plan-Commit) und der V-4-Lehre des Verifier-Reports dieses Vorgangs; die Bestand-Zeile war da, trug aber den Start-Stand weiter. (b) Der getaggte Lead trug die engere Abweisungsliste (F-4) — der Nachzug auf `main` folgt dem Muster der Vorgänger-Cuts, weil ein bereits veröffentlichter Tag nicht verschoben wird (`releasing.md` §8). (c) Die Fundstellen des Review-Reports für `repo.mk` trugen 66/67 vertauscht (V-2, berichtigt). (d) Parallel zu diesem Cut ist der Vorschlag-Slice [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md) entstanden (Wunsch des Projektinhabers, 2026-11-03, ENV-Rest des eigenen Devcontainers in `u-boot.yaml`) — im Scope dieses Cuts nicht enthalten, nur Plan.

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre:* die Zustandsmarker-Klasse erreicht in den drei Release-Cuts in Folge immer wieder dieselben zwei Träger — §7-Marker des Plans und die Bestand-Zeile `open/README.md`. Bestehende Regeln tragen beide (der Plan §7 wird bei der Closure ohnehin umgeschrieben, die Bestand-Zeile trägt „nur Plan" als Beschreibung ihres Zustands); ein Sensor-Kandidat („Bestand-Zeile und §7-Marker nach dem Tag gegen die Verzeichnis-Lage") bleibt zur Entscheidung beim Projektinhaber. *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) ([`evidence/slice-v1-release-cut-v0.20.0.md`](../observations/BEO-DOK/plan-aussage-ueber-code/evidence/slice-v1-release-cut-v0.20.0.md): F-1/F-3 (Marker), V-2 (Fundstellen), V-5 (Bericht-Dateien); berichtigt `f652ed4f` und der Berichtigungs-Runde) · [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) ([`evidence/slice-v1-release-cut-v0.20.0.md`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/slice-v1-release-cut-v0.20.0.md): F-4 dritte Abweisung im Lead, die Grenzlisten-Klassenhistorie). *Nicht gezählt:* [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — kein Fund als Zusage-Wider (F-4 ist eine engere, nicht eine weitergehende Liste); [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md) — kein Fund.

**Carveouts:** keiner — kein Gate, keine Schwelle und kein Lint-Bestand sind im Range berührt (Verifier-`--name-status`-Lesung; `Dockerfile`/`repo.mk`-Diffs sind reine Versionszeilen samt Kommentar-Nennungen). **Folge-Slices:** der Weg-A/B/C-Entscheidung zu [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md) und [`slice-v1-katalog-features-quellen-und-freigabe`](../in-progress/slice-v1-katalog-features-quellen-und-freigabe.md) tragen die Restlinie der Feature-Katalog-Zeile; der Vorschlag-Slice [`slice-v1-katalog-projekt-containerenv`](../done/slice-v1-katalog-projekt-containerenv.md) liegt parallel (2026-11-03, gewünscht). **Weiter offen:** nichts aus der DoD.

**Drei Paarungen:** *Anker* — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert. *Folge-Slice* — der Vorschlag-Slice `containerenv` und die zwei `open/`-Reste sind Dateien mit Plänen (`ls`; keine Cuts versprochen). *Register* — jede `BEO-*`-Kennung in §7 ist ein Verzeichnis mit der Evidence-Datei dieser Closure (Verzeichnis-Lesung; die `state.md`-Einträge wurden am Schreibort geschrieben und folgen der Ruheort-Fassung nach dem Move).

**Verifikation** (Läufe nach Closure-Inhalt, Move, Register-Anker und Verweis-Nachzug; Docker-Modus): `make gates` Exit 0 am 2026-11-03 (`record-gates`-Stempel = HEAD nach dem Nachzug). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `CHANGELOG.md`, beide
READMEs, `docs/user/` (DOK) und `docs/plan/` (Planning, DOK); laut
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area)
alle Greenfield (DOK). Keine Code-Sub-Area ist im Diff (Plan oben).

**Vorgelagert — offene Beobachtungen sichten:** [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)
(Lead, Risiko 1), [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md)
(Risiko 3). Evidenz-/Diskrepanz-Risiko: niedrig — der Cut führt kein
Produktverhalten, nur Release-Prosa, Versionsteile und Marker.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield (kein
Modus-Begründungsblock nötig).