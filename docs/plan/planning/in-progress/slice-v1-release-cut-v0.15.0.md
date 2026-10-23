# Slice V1: Release-Cut `v0.15.0`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-forge-cli-gh-glab`](../welle-forge-cli-gh-glab.md), sechster von sechs Slices (Release-Version nach [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md): `v0.15.0`, vom Projektinhaber am 2026-10-23 bestätigt, Kopf und §5 der Welle-Datei). Der Release trägt alles, was seit dem Tag `v0.14.0` (`68f3c94f`) auf `main` liegt und in `CHANGELOG.md` unter `## [Unreleased]` steht: die Features `gh` und `glab` ([`slice-v1-devcontainer-feature-gh-glab`](../done/slice-v1-devcontainer-feature-gh-glab.md)) mit Lastenheft `0.22.0` ([`slice-cr-devcontainer-forge-cli`](../done/slice-cr-devcontainer-forge-cli.md)), Spezifikation und Handbuch ([`slice-v2-spec-forge-cli`](../done/slice-v2-spec-forge-cli.md)) und die Nachbesserung nach dem Validator ([`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md)); Grundlage ist die Messung [`slice-v1-forge-cli-quellen-messung`](../done/slice-v1-forge-cli-quellen-messung.md) mit [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md). Stand beim Schreiben (2026-10-23, gelesen): `HEAD` = `07353780`, Arbeitsbaum sauber, `in-progress/` ohne Slice, Lastenheft `0.22.0` (`spec/lastenheft.md` Zeile 9).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
[`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion),
[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte),
[`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md),
[`ADR-0004`](../../adr/0004-ci-system.md),
[`ADR-0007`](../../adr/0007-distributionswege-ghcr.md),
[`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md),
[`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md). Ablauf:
[`docs/maintainer/releasing.md`](../../../maintainer/releasing.md). Muster:
[`slice-v1-release-cut-v0.14.0`](../done/slice-v1-release-cut-v0.14.0.md). Aufträge: Welle-Datei §5 (*Übergaben*) und der Validator-Bericht
[`2026-10-23-welle-forge-cli-gh-glab-validierung`](../../../reviews/2026-10-23-welle-forge-cli-gh-glab-validierung.md)
(Urteil GO für den Cut mit Auflagen; A-1 dieses Slice, A-2 und A-3 liegen im Handbuch-Slice).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; Lastenheft `0.22.0` liefern die Slices; dieser Slice ändert `spec/` nicht).

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-23); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.15.0` liefert als Minor-Release den Block `[Unreleased]` des CHANGELOG — die Katalog-Features `gh` und `glab` im Devcontainer (beide Profile; Zusage nur Linux-Host `amd64` mit lokaler Docker-Engine, **kein authentifizierter Aufruf gemessen**) — als GHCR-Image, Binaries, Pakete und Homebrew-Formel; der Release-Text trägt die Auflage A-1 des Validators und die Präzisierung zu `version:` ohne `source:` (M-2). Dieser Slice bereitet vor und lässt den Tag zu; **Push und Tag führt nur der Hauptagent, nach ausdrücklicher Freigabe des Projektinhabers im Gespräch, je Vorgang** (§3, *Reihenfolge der Maintainer-Aktionen*).

### Versionsnummer: `v0.15.0` (bestätigt durch den Projektinhaber, Welle-Datei §5, Frage 1)

Gelesen: `CHANGELOG.md` Zeilen 12 bis 20 (`[Unreleased]`: ein Eintrag *Added*, ein Eintrag *Changed*, kein *Removed*, kein *Security*), `spec/lastenheft.md` Zeile 9 (`0.22.0`), Welle-Datei (Kopf).

- **Minor, nicht Patch:** neue Funktion (zwei Katalog-Features, SemVer 2.0 §7).
- **Kein Major:** Die Hauptversion ist `0` (SemVer §4). Kein Befehl und keine Option entfällt. Eine Verhaltensänderung für bestehende Konfigurationen gibt es nicht: ohne `gh`/`glab` rendert alles byte-identisch (`CHANGELOG.md` Zeile 20, *Changed*); der Lead nennt das als solches.
- Widerspricht der Projektinhaber beim Start, entsteht ein neuer Slice mit passender Kennung und dieser geht mit `Gegenstand:` nach `done/`.

### Inhalt des Release (gelesen am 2026-10-23, `CHANGELOG.md` Block `[Unreleased]`)

| Eintrag im CHANGELOG (Zeile) | Art | Hinweis für den Lead |
|---|---|---|
| Katalog-Features `gh` und `glab` (16) | Added | Zusage-Menge und „kein authentifizierter Aufruf gemessen“ (A-1); der Satz „`version:` without `source:` ends with exit code `10`“ gilt nur für `generate devcontainer` (M-2); Nicht-Zusagen sind im Eintrag schon genannt (Aufrufe mit gültigem Token, Hosts nach Anmeldung, Telemetrie- und Update-Hosts, selbst gehostete Instanzen, `glab duo`, `glab orbit`) |
| Dockerfile des Standardprofils trägt mit `gh`/`glab` `USER root` vor dem Block (20) | Changed | nur mit den neuen Features; Konfigurationen ohne sie rendern byte-identisch |

### Aufträge an den Release-Cut

Jede neue Aussage in Handbuch, CHANGELOG und READMEs braucht einen Beleg nach [`AGENTS.md`](../../../../AGENTS.md) §3.10 (Beleg-Tabelle in §3); was keinen hat, steht nicht im Text oder steht als „nicht gemessen“. Kein Punkt verlangt Produktcode.

1. **A-1 — Lead nennt die Zusage-Menge und die Grenze der Messung** (Validator, *Auflagen an den Cut*): zugesagt nur Linux-Host `amd64` mit lokaler Docker-Engine; **kein authentifizierter Aufruf von `gh` oder `glab` gemessen** (alle Aufrufe mit Variable liefen mit Platzhaltern und endeten in 401 oder `Bad credentials`). Dazu die Grenzen aus dem Validator-Bericht (*Was ohne gültiges Token ungeprüft blieb*), soweit der Lead sie braucht, ohne das Handbuch zu wiederholen.
2. **M-2 — Satz zu `version:` ohne `source:` präzisieren** (Welle-Datei §5, *Übergaben*; Handbuch-Slice): Der Satz nennt `u-boot generate devcontainer` als Ort des Exit-Codes `10`. Hinweis im Lead: `u-boot config set …features.<feature>.version` prüft das noch nicht (es schreibt die `version`; danach lehnen `generate`, `doctor` und jedes weitere `config set` das Projekt ab, Abhilfe: den Eintrag `version` in `u-boot.yaml` von Hand entfernen; das Handbuch sagt es, Zeile 2060); die Prüfung ist Folge ([`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md), nicht vor dem Cut). Der Cut sagt die Lage ehrlich im Release-Text, ohne Produktcode.
3. **Nicht-Zusagen im Lead stehen lassen und bündeln:** `glab duo` und `glab orbit` (nicht zugesagt), Telemetrie nicht beeinflusst und im Handbuch genannt ([`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkt e, Standard aus der Empfehlung), selbst gehostete Instanzen, `arm64`-Bau, Standardprofil: keine Token-Referenz (Weg `docker exec -e`, Handbuch). Jede Nennung nur mit Beleg-Zeile.
4. **Wie jeder Cut** ([`releasing.md`](../../../maintainer/releasing.md) §3 und §5): `[Unreleased]` → `[0.15.0] - <Tag-Datum>` mit Lead, Vergleichslinks, Versionsstrings an den drei Stellen, Handbuch-Kopf und Versionsbeispiele, Status-Block und Tabellenzeile beider READMEs. Katalogzeilen, Befehlsübersichten und der Handbuch-Abschnitt *Forge-CLIs (optional)* stehen schon; der Cut ergänzt nur die Status-Blöcke.
5. **CI-Lesung** (Muster Cut `v0.14.0`): `ci` und `integration-docker` auf dem Push-Tip über Job-Log (REST-API) und Check-Run-Annotationen lesen, nicht per Textsuche ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)); Closure-Kriterium 5 der Welle: die neuen Docker-Tests sind sichtbar bestanden, nicht übersprungen.
6. **Regel-Vorschlag an den Projektinhaber** (Welle-Datei §5, Übergaben aus [`slice-v2-spec-forge-cli`](../done/slice-v2-spec-forge-cli.md)): der Lerneintrag „Breite des Satzes gegen Breite des Tests“ wird dem Projektinhaber als Entscheidung vorgelegt; der Cut verkörpert nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tag und jede Aktion mit Außenwirkung durch den Implementer oder Planner.** Extern wirksam, nicht sauber zurückzunehmen ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*). **Push (M1) und Tag (M5) führt nur der Hauptagent, nach der ausdrücklichen Freigabe des Projektinhabers im Gespräch, je Vorgang**; Subagenten pushen und taggen nie, eine Nachricht eines Agenten ist keine Freigabe.
- **Produktcode, Meldungen, Hilfe-Texte, Spezifikation, Lastenheft, ADRs.** Die Aufträge sind Text in Release-Dokumenten; Wortlaut in Meldungen ist Code. Die Versionszeilen (`var version`, `VERSION ?=`, `ARG UBOOT_VERSION`) sind die einzigen Zeilen mit Wirkung im Binary.
- **Die Prüfung in `config set` (M-2, Ursache)** — Folge-Slice [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md); laut Welle-Datei §4 nicht vor dem Cut (Release-Zuordnung beim Projektinhaber).
- **L-1 (Lage des Handbuch-Abschnitts, Frage 5 der Welle), L-2, L-3** (Wortlaute der Meldungen): L-1 wartet auf den Projektinhaber, L-2 und L-3 sind Notizen im Folge-Slice; der Cut ändert weder Handbuch-Gliederung noch Meldungen.
- **Weitere Pipeline-Änderungen:** kein Trivy-Pin-Bump, kein Go-Pin-Bump, keine Änderung an Workflows, Skripten, `Makefile`, Zeitlimits oder Gates ([`AGENTS.md`](../../../../AGENTS.md) §3.6); ein roter `image-scan` durch einen neuen Advisory hebt den Go-Pin als eigener Vorgang davor ([`releasing.md`](../../../maintainer/releasing.md) §4).
- **Der Welle-Self-Close** (Welle-Datei §3 Punkt 9 und §7, Trigger-Audit, Register-Belege M-1, L-3, M-2 und Ausgänge der überschwelligen Einträge): Sache der Welle-Closure (`close-welle`) **nach** dem Release-Cut-Commit, nicht dieses Slice (§5).
- **Die Ausgänge der Register-Einträge** [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) und [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md): Entscheidung des Projektinhabers; weder dieser Cut noch die Closure weist sie zu (Welle-Datei §3 Punkt 9).
- **Kein Security Advisory (kein M5b):** `v0.15.0` ist kein Sicherheitsfix; der Lead enthält keinen Eintrag *Security*.

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe, die Kontrollen unter Punkt 2 und die Closure-Pflichten darunter zählen nicht mit.

- [ ] **Release-Text und Handbuch:** `CHANGELOG.md` trägt `## [0.15.0] - <Tag-Datum>` mit Lead nach A-1, M-2 und den Nicht-Zusagen und den Vergleichslinks; Versionsstrings `0.15.0-dev` an den drei Stellen; beide READMEs (Status-Block, Tabellenzeile, `VERSION=`) und das Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich) tragen `v0.15.0`. Jede neue Aussage steht mit ihrem Beleg in der Beleg-Tabelle (§3); Handbuch und READMEs tragen keine Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] **Tag und Kontrolle:** Vor dem Tag sind `make gates`, `make ci` und `make fullbuild` grün, `make test-docker` ist gelaufen (Ergebnis mit Zahl der übersprungenen Tests in §3), der Trockenlauf des Tap-Pfads (M-T) ist festgehalten, `tap-check` ist grün, und `ci` und `integration-docker` auf dem Push-Tip sind grün mit den neuen Docker-Tests sichtbar bestanden (M2). **Erst danach und nur nach ausdrücklicher Freigabe des Projektinhabers im Gespräch:** Tag `v0.15.0` gesetzt und gepusht (M5). Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 (M6).
  - [ ] Release-Workflow `publish`: Jobs `publish`, `tap`, Paket-Smoke rpm/fedora und deb/ubuntu `success` (Lauf-ID in §7).
  - [ ] GHCR-Image: `docker run --rm ghcr.io/pt9912/u-boot:0.15.0 --version` meldet `0.15.0`; `:0.15.0` und `:latest` tragen denselben Digest.
  - [ ] Binaries und Pakete: `gh release view v0.15.0 --json assets,isPrerelease,isDraft` zeigt sechs Binaries, vier Pakete, `SHA256SUMS` und `u-boot.rb` (zwölf Assets), kein Prerelease, kein Draft.
  - [ ] Release-Notes: der GitHub-Release verweist auf `CHANGELOG.md` §[0.15.0] (Auto-Notes) oder trägt vom Maintainer geschriebene Notes; der Abschnitt `[0.15.0]` im CHANGELOG auf `main` ist der des Tags.
  - [ ] Job `tap` wirklich gelaufen (Log „tap-nachzug: … updated to v0.15.0.“ oder „already at“, keine Warnung zum fehlenden Secret, Annotationen gelesen).
- [ ] **Befunde des Validators haben einen Ausgang:** A-1 und M-2 (Hinweis) sind im Text, mit Beleg; A-2 und A-3 liegen im Handbuch-Slice (`done/`), L-1 bis L-3 und I-1, I-2 sind Folge oder ohne Aufgabe laut Welle-Datei §5 (§6, Closure-Notiz), ohne dass der Cut sie entscheidet.
- [ ] `make gates` grün (nach dem letzten Commit vor dem Handoff und nach der Closure).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); Verification-Report getrennt davon (Verifier), vor M1.
- [ ] Doku-Update: CHANGELOG, READMEs, Handbuch, `releasing.md` (Zeiger „zuletzt“), Roadmap (Zeilen der Welle und `v0.15.0` mit Tag-Hash) und die abgelaufenen Zustandsmarker (Welle-Datei §4, `open/README.md`) tragen den Ist-Stand.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (§7); die Welle-Closure prüft sie für die Welle erneut.

## 3. Plan (vor Code)

Jede Aussage in dieser Tabelle ist an der genannten Stelle gelesen (Stand 2026-10-23, `HEAD` = `07353780`, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Zeilennummern nennen den Fundort beim Schreiben; der Implementer sucht sie vor der Änderung neu. Vor dem Handoff und nach jeder späteren Fix-Runde gleicht er die Tabelle gegen `git diff --stat <Plan-Commit>..HEAD` ab — einschließlich dieser Datei, der Vermerke in `done/`, der Reports und der Verweise, die `make slice-mv` nachzieht; eine Datei, die er als „kein Diff“ führt, belegt er durch eine gelesene Aussage, nicht durch eine Suche.

**Versionsstellen** (gelesen mit `grep -rn "0\.14\.0"` über `*.go`, `*.mk`, `*.md`, `*.yml`, `Makefile`, `Dockerfile` außerhalb von `CHANGELOG.md`, `docs/plan/`, `docs/reviews/`, `docs/messungen/`, `.harness/`): `cmd/uboot/main.go:43`, `repo.mk:51`, `:54`, `:55`, `Dockerfile:35`, `:37`, `:39`, `README.md:10`, `:102`, `:208`, `README.de.md:11`, `:104`, `:211`, `docs/user/benutzerhandbuch.md:4`, `:82`, `:108`, `:135`, `:2985`. Der Implementer wiederholt die Suche vor und nach dem Cut (auch mit `-e "0.14.0-dev"` für die Kommentare) und prüft `docs/maintainer/releasing.md:56`–`:57` (Zeiger „zuletzt“ auf den Cut `v0.14.0`, richtig bis zur Closure).

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update (Cut) | `:12` `## [Unreleased]` wird zum leeren Anker; der Inhalt (`:14`–`:20`) wandert unter `## [0.15.0] - <Datum M0>` mit Lead „Fifteenth release.“ (Stil von `[0.14.0]`). Eintrag *Added* (`:16`): der Satz „`version:` without `source:` ends with exit code `10`“ wird auf `u-boot generate devcontainer` präzisiert, mit dem Hinweis zu `config set` (Auftrag 2). Lead nach *Aussage des Release*. Vergleichslinks am Fuß: `:1504` `[Unreleased]: …/compare/v0.14.0...HEAD` → `v0.15.0...HEAD`; neue Zeile `[0.15.0]: …/compare/v0.14.0...v0.15.0` vor `:1505` (`[0.14.0]`). |
| `cmd/uboot/main.go` | update | `:43` `var version = "0.14.0-dev"` → `0.15.0-dev`. Einzige Zeile mit Wirkung im Binary; ein Test pinnt die Nummer nicht (im Cut `v0.14.0` an `main_test.go` gelesen; der Implementer liest es erneut). |
| `repo.mk` | update | `:51` und `:54` (Kommentare) und `:55` `VERSION ?= 0.14.0-dev` → `0.15.0-dev`. `TRIVY_VERSION` bleibt. |
| `Dockerfile` | update | `:35` und `:37` (Kommentare) und `:39` `ARG UBOOT_VERSION=0.14.0-dev` → `0.15.0-dev`. Kein weiterer Diff (kein Go-Pin-Bump, §1). |
| `README.md` | update | Status-Block ab `:10`: `:10` Version `v0.15.0`; die Katalogzeile (`:13`, „The devcontainer catalogue has 16 features, including … `gh`, `glab` …“, gelesen) bekommt einen Satz zu `gh`/`glab` in Worten für den Benutzer (Zusage-Menge Linux `amd64`, lokale Engine; kein Aufruf mit gültigem Token gemessen), ohne Datum, ohne „released“, ohne Kennungen; `:102` `VERSION=0.14.0` → `0.15.0`; neue Zeile in der Release-Tabelle nach `:208` (Zeile `v0.14.0`; die Tabelle läuft aufsteigend, `:209` ist `v0.5.0`), Soll unten. |
| `README.de.md` | update | gleiches Muster: `:11` Version, Katalogzeile (`:14`), `:104` `VERSION=0.14.0`, Tabellenzeile nach `:211`. |
| `docs/user/benutzerhandbuch.md` | update | Kopf `:4` `Software-Version: v0.14.0` → `v0.15.0`, `:5` `Stand: 2026-10-23` → Tag-Datum; `:82` `u-boot version 0.14.0`, `:108` `download/v0.14.0`, `:135` `VERSION=0.14.0`, `:2985` „beschreibt `u-boot` v0.14.0“ → `0.15.0`. **Kein Diff:** Abschnitt *Forge-CLIs (optional)* und Katalog-/Egress-Zeilen (stehen aus den Slices der Welle; die Abhilfe zur ungültigen `version` steht bei `:2060`, gelesen); vor dem Cut liest der Implementer Handbuch und Lead gegeneinander. Kennungsfrei. |
| `docs/maintainer/releasing.md` | kein Diff im Cut; update bei der Closure | `:56`–`:57` nennt als „zuletzt“ den Cut `v0.14.0` (gelesen): richtig bis zum Abschluss; die Closure setzt den Zeiger auf diesen Slice (Muster: Closure von [`slice-v1-release-cut-v0.14.0`](../done/slice-v1-release-cut-v0.14.0.md)). |
| `docs/plan/planning/in-progress/roadmap.md` | update (Plan-Commit: nur `:20`; Start, Cut, Closure) | `:20` führt die Welle unter „Offene Wellen“ mit „(die Kennung des sechsten vorläufig)“ — Plan-Commit: nennt den Slice als Datei in `open/`. `:22` `Nichts in Arbeit.` (Ruhe-Marker, `.d-check.yml` `planning`). Start: der Marker entfällt, an seine Stelle tritt ein Zeiger auf diesen Slice (eigener Commit direkt nach dem Move). Nach dem Tag: unter *Abgeschlossene Wellen* nach der Zeile `v0.14.0` (`:104`) die Zeile der Welle und die Zeile `v0.15.0` mit Datum, Tag-Hash und Verweis auf diesen Slice (Muster `:103`–`:104`); Closure: Zeiger weg, Ruhe-Marker zurück. *Nächste Wellen* ohne Diff (gelesen: `welle-ollama-cloud-modelle` wartet auf Beauftragung und Reihenfolge, beides Sache des Projektinhabers). |
| `docs/plan/planning/welle-forge-cli-gh-glab.md` | update (Planner, Plan-Commit; Verweis-Nachzug) | Abgelaufene Zustandsmarker: §4 Zeile 6 (`slice-v1-release-cut-v0.15.0` „(Datei erst beim Trigger)“, `:71`) und der Absatz darunter (`:73`, „Die Kennung des Slice in Zeile 6 ist vorläufig …“). Der Trigger ist eingetreten (die fünf Slices in `done/`, Validator-Bericht vom 2026-10-23); die Zeile verlinkt den Slice in `open/`. §7 (Closure-Notiz) bleibt der Closure (`close-welle`). |
| `docs/plan/planning/open/README.md` | update (Plan-Commit) | neue Zeile für diesen Slice am Ende der Tabelle (nach der Zeile von [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md)); die Zeile von [`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md) trägt den abgelaufenen Marker „in Arbeit: Handbuch und Tests umgesetzt“ und wird auf den Ist-Stand „geliefert (liegt in `done/`)“ gesetzt. Die Moves ziehen den Link per `make slice-mv` nach. |
| `docs/plan/planning/done/slice-v1-release-cut-v0.15.0.md` (dieser Slice; Plan in `open/`, später `done/`) | neu, update | Plan (Planner), Moves per `make slice-mv` nach `in-progress/` und nach `done/` (reine Commits, [`AGENTS.md`](../../../../AGENTS.md) §3.3), Umsetzung in §3 (Implementer), Closure (Planner): §2, §3, §6, §7. |
| `docs/reviews/<Datum>-slice-v1-release-cut-v0.15.0.md`, `…-verify.md` | neu, durch Reviewer und Verifier; Verweis-Nachzug | Review-Report und Verification-Report (vor M1; Dateiname `<Datum>-<Slice-Kennung>-verify.md`). |
| `docs/plan/planning/observations/…` | nur Closure | Stand-Zeilen und Evidence-Dateien schreibt die Closure; **nicht im Cut:** die Evidence-Datei `welle-forge-cli-gh-glab.md` zu `BEO-DOK/vertrag-schweigt-zum-nutzerpfad` (M-1), `BEO-CLI/kennung-in-cli-ausgabe` (L-3) und `BEO-DOK/spec-bedingung-ohne-code-entsprechung` (M-2) schreibt die Welle-Closure (Welle-Datei §5). |
| `spec/`, `docs/plan/adr/` | kein Diff (gelesen) | Lastenheft `0.22.0` bleibt (`spec/lastenheft.md` Zeile 9). `make doc-immutable RANGE="v0.14.0..HEAD"` belegt, dass keine Accepted-ADR umgeschrieben wurde. |
| `.github/workflows/*`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `tools/`, `d-check.mk`, `docs/maintainer/quality.md`, `docs/maintainer/cli-json-contract.md`, `internal/`, übrige `cmd/uboot/` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung, kein Produktcode (§1). Der Implementer belegt „kein Diff“ vor dem Handoff durch eine gelesene Aussage (Statuszeile, Pin-Zeile), nicht durch eine Suche. |

### Beleg-Tabelle der neuen Aussagen ([`AGENTS.md`](../../../../AGENTS.md) §3.10)

Der Implementer schreibt nur, was hier einen Beleg hat, in der Form der letzten Spalte; eine Aussage ohne Zeile steht nicht im Text.

| Aussage (Soll) | Ort | Beleg | Form im Text |
|---|---|---|---|
| Zusage-Menge: Linux `amd64`, lokale Docker-Engine; nicht gemessen und nicht zugesagt: macOS mit Colima (auch für Homebrew-Nutzer), Docker Desktop, rootless Docker, Podman als Host-Engine, entfernte Engine, Linux `arm64`, Windows | `CHANGELOG.md`-Lead (A-1), Status-Blöcke | Welle-Datei §5 (Zusage-Menge), Validator-Bericht, *Messumgebung* (native Docker Engine 29.9.2, `amd64`) und *Was ohne gültiges Token ungeprüft blieb* (`arm64`, Docker Desktop, Colima, macOS) | Zusage mit Grenze, wie im Lead von `v0.14.0` |
| Kein authentifizierter Aufruf von `gh` oder `glab` ist gemessen; Aufrufe mit Platzhalter endeten in 401 oder `Bad credentials` | `CHANGELOG.md`-Lead (A-1) | Validator-Bericht, *Was ohne gültiges Token ungeprüft blieb* und Punkt 5 | „nicht gemessen“, nicht „nicht unterstützt“ |
| `gh` 2.103.0 und `glab` 1.124.0 sind ein Stand (gepinnte Release-Archive, SHA-256-Prüfung vor dem Entpacken; keine Signatur nach Suchmuster gefunden, Attestierung von `gh` nicht geprüft) | `CHANGELOG.md`-Lead (kurz) | `CHANGELOG.md` Zeile 16 (besteht schon); [`BEO-APP/podman-static-pin-alterung`](../observations/BEO-APP/podman-static-pin-alterung/observation.md) für die Alterung der Pins | Stand, mit „by search pattern, not exhaustive“ |
| `version:` ohne `source:` endet bei `generate devcontainer` mit Exit `10`; `config set …version` prüft es nicht und sperrt danach das Projekt (Abhilfe: Eintrag `version` in `u-boot.yaml` von Hand entfernen) | `CHANGELOG.md` `:16` und Lead (M-2) | Validator-Bericht, Befund M-2 (ein Lauf, Linux `amd64`); Handbuch `:2060`; Plan [`slice-v1-config-set-validiert-feature-version`](../open/slice-v1-config-set-validiert-feature-version.md) (am Code gelesen: `config set` schreibt, `generate` lehnt ab) | Messung „ein Lauf“; der Folge-Slice wird als Folge genannt, nicht als Zusage |
| Standardprofil: keine Token-Referenz in `remoteEnv`; Weg zum Token `docker exec -e` (mit Platzhalter gemessen, ein Lauf) | Lead nur als Verweis aufs Handbuch | Validator-Bericht, Punkt 5; [`slice-v1-handbuch-forge-cli-nachbesserung`](../done/slice-v1-handbuch-forge-cli-nachbesserung.md) (Messung A4) | kurz, Messung „ein Lauf“; die Messung steht dort ohne Protokoll unter `docs/messungen/` |
| Nicht zugesagt: `glab duo`, `glab orbit`, selbst gehostete Instanzen, Hosts nach einer Anmeldung, Telemetrie- und Update-Hosts | `CHANGELOG.md`-Lead | `CHANGELOG.md` Zeile 16 (besteht schon); [`ADR-0035`](../../adr/0035-forge-clis-gh-und-glab-im-devcontainer.md) Punkte e und h | Nicht-Zusage, ohne neue Behauptung |
| Mit `gh`/`glab` trägt das Dockerfile des Standardprofils `USER root` vor dem Block; ohne die Features rendert alles byte-identisch | `CHANGELOG.md` `:20` (besteht schon) | `CHANGELOG.md` Zeile 20; Golden Cases des Code-Slice | Satz wie im Eintrag |
| CI: Docker-Integrationstests grün, die vier Nested-Podman-Tests übersprungen | `CHANGELOG.md`-Lead, Closure | Lauf `integration-docker` 39890008034 auf `68f3c94f` (Cut `v0.14.0`, §3); gilt für diesen Release erst nach M2 auf dem neuen Push-Tip | **Kein CI-Satz im Lead**, solange M2 aussteht (Muster Cut `v0.14.0`, Abweichung 1); die Status-Blöcke behalten ihren bestehenden Satz unverändert |

### Aussage des Release — Soll für den CHANGELOG-Lead

Tabellenzeile und Status-Block beider READMEs tragen eine Kurzform ohne Kennungen; sie widersprechen dem Lead nicht und sagen nichts, was er nicht sagt. Der Lead steht im Stil von `[0.14.0]` (englisch, `Fifteenth release.` und fette Kernaussage).

- **Neu:** die Katalog-Features `gh` (GitHub CLI) und `glab` (GitLab CLI): `u-boot config set devcontainer.features.gh.enabled true` (und `glab`), dann `u-boot generate devcontainer`; das Werkzeug liegt nach dem Bau im Container, in beiden Profilen, im Sandbox-Profil mit seinen Hosts in der Egress-Liste und ausgewiesen; ein Token kommt nur als Umgebungsvariable des Hosts (`GH_TOKEN`, `GITLAB_TOKEN`) über `${localEnv:…}` in den Container, nie als Wert in eine Datei (Sandbox-Profil; das Standardprofil trägt keine Referenz).
- **Zusage-Menge (A-1):** nur ein Linux-Host `amd64` mit lokaler Docker-Engine; **nicht zugesagt:** macOS mit Colima (auch für Homebrew-Nutzer), Docker Desktop, rootless Docker, Podman als Host-Engine, entfernte Engine, Linux `arm64`, Windows.
- **Gemessen ist, was ohne gültiges Token möglich ist (A-1):** Bau, Version, Ausweisung, Egress-Hosts, Verhalten mit und ohne Variable. **Kein authentifizierter Aufruf von `gh` oder `glab` ist gemessen** (alle Aufrufe mit Platzhalter endeten in 401 oder `Bad credentials`).
- **Nicht zugesagt:** `glab duo`, `glab orbit`, selbst gehostete Instanzen, Hosts nach einer Anmeldung, Telemetrie- und Update-Hosts (nicht beeinflusst, im Handbuch genannt).
- **Grenze (M-2):** `version:` ohne `source:` weist `u-boot generate devcontainer` mit Exit `10` ab; `u-boot config set …features.<feature>.version` prüft das noch nicht und sperrt das Projekt (Abhilfe: Eintrag von Hand aus `u-boot.yaml` entfernen); die Prüfung ist Folge.
- **Unverändert:** Konfigurationen ohne `gh` und `glab` rendern byte-identisch. Lastenheft `0.22.0` (nur im CHANGELOG).

Tabellenzeile der READMEs (Muster `v0.14.0`): Name „Forge CLIs in the devcontainer“ / „Forge-CLIs im Devcontainer“, dann `gh` und `glab`, Plattform, „kein Aufruf mit gültigem Token gemessen“, `version` ohne `source` in Worten, Link auf den GitHub-Release.

### Reihenfolge der Maintainer-Aktionen

Der Implementer bereitet vor (Commits lokal, Gates lokal) und **führt M1 und M5 nicht aus**. **Push (M1) und Tag (M5) führt nur der Hauptagent, nach ausdrücklicher Freigabe des Projektinhabers im Gespräch, je Vorgang.**

> **Rückfrage-Punkt (Hauptagent, nicht Planner und nicht Implementer):** Vor **M1** (Push des Cuts) und vor **M5** (Tag) fragt der Hauptagent den Projektinhaber und wartet auf die ausdrückliche Freigabe im Gespräch. Eine Freigabe für M1 ist keine für M5. Der Slice bleibt bis zum Tag in `in-progress/`.

| # | Aktion | Wer | Prüfung / Erwartung |
|---|---|---|---|
| M0 | `date +%F` (Ortszeit); Datum in `CHANGELOG.md`, Tabellenzeile beider READMEs und Handbuch-`Stand` gegen das Tag-Datum halten | Implementer vor dem Handoff; Hauptagent vor M1 und vor M5 | Weicht es ab: ein weiterer Cut-Commit vor dem Tag und M2 auf dem neuen Push-Tip erneut ([`releasing.md`](../../../maintainer/releasing.md) §5); Cut, Push und Tag am selben Tag |
| M-T | Trockenlauf des Tap-Pfads (`releasing.md` §4), im Docker, `--network none`, ohne Schreiben ins Tap; Muster Cut `v0.14.0` §3 | Implementer; Ergebnis in §3; der Verifier fährt ihn nach | `tap-formula-check.sh` Exit `0` mit `v0.15.0` und Exit `1` mit `v0.14.0`; `tap-nachzug.sh` mit Stellvertretern protokolliert mit `v0.15.0` `clone`, `add`, `commit -m u-boot 0.15.0` und `push` und bricht mit `v0.14.0` vor `clone` ab. **Beweist nicht:** `is_stable`, echter `curl`, echter `git push`, Token |
| M1 | `git push origin main` (alle Commits des Cuts in einem Push) | Hauptagent, **nach Freigabe des Projektinhabers** | Review- und Verification-Report liegen vor; `make docs-check` lokal grün; Arbeitsbaum sauber |
| M2 | `gh run list --limit 4`: `ci` und `integration-docker` auf dem Push-Tip; Logs über `gh api repos/pt9912/u-boot/actions/jobs/<id>/logs`, Warnungen über `gh api repos/pt9912/u-boot/check-runs/<id>/annotations` ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)) | Hauptagent | `ci` `success`: Jobs `gates`, `security-gates`, `image-scan`; keine Annotation „Unexpected input(s)“. `integration-docker` `success` (Konklusion **und** Zählung): `--- FAIL` 0, kein `panic: test timed out`; `--- SKIP` genau die vier bekannten Tests (`TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_NestedPodman`, `…_StatusFiles_VFS`); die neuen Docker-Tests der Welle sichtbar `--- PASS`: `TestE2E_ForgeCLI_BuildsInBothProfiles`, `…_PinMismatchBreaksBuild`, `…_EgressHosts`, `…_TokenReferencesCarryNoValue` (`internal/e2e/forge_cli_docker_test.go` Zeilen 47, 75, 133, 221); der Implementer ergänzt die Tests, die die Nachbesserung ohne Produktcode hinzufügte, soweit sie Docker-Tests sind. Dauer des Jobs gegen `timeout-minutes: 30` festhalten; Abschnitt „Operating System“ der Logs lesen. Abweichend: Release-Aussage vor dem Tag nachziehen |
| M3 | Vorab-Tag `v0.15.0-rc.1`: entfällt | — | Vorschlag des Planners (wie `v0.13.0` bis `v0.14.0`: direkt stabil); bestätigt der Projektinhaber beim Start (§7) |
| M4 | `gh workflow run tap-check` | Hauptagent, unmittelbar vor M5 | grün = das Token darf auf das Tap schreiben (`releasing.md` §2, §4) |
| M5 | `git switch main && git pull --ff-only`, `git tag v0.15.0` auf dem in M2 geprüften Commit, `git push origin v0.15.0` (nicht `--tags`) | **Hauptagent, nur nach ausdrücklicher Freigabe des Projektinhabers im Gespräch** | `HEAD` = der Commit aus M2; nach dem Push wird der Tag nicht verschoben (`releasing.md` §8) |
| M5b | Security Advisory: **entfällt** | — | kein Sicherheitsfix (§1) |
| M6 | `releasing.md` §6: `gh run watch`; `gh release view v0.15.0 --json assets,isPrerelease,isDraft`; `docker run --rm ghcr.io/pt9912/u-boot:0.15.0 --version`; Digest `:0.15.0` = `:latest`; Log und Annotationen des Jobs `tap`; Tap-Commit und Smoke-Workflow im Tap | Hauptagent, unmittelbar nach dem Push des Tags | DoD Punkt 2, Kontrollzeilen |
| M7 | Auf einem Mac: `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot`, `brew install u-boot` (oder `brew upgrade u-boot`), `u-boot --version` = `0.15.0` | Projektinhaber | `releasing.md` §6; ersetzt keine Zusage für `gh`/`glab` auf macOS (nicht zugesagt) |
| M8 | Nacharbeit (`releasing.md` §7): Roadmap-Zeilen mit Tag-Hash und Datum, Ruhe-Marker; Datum gegen `git log -1 --format=%cs v0.15.0`; Zeiger „zuletzt“ in `releasing.md`; Closure-Notiz, Register, `git mv` nach `done/`; danach `close-welle` | Planner (Closure), Implementer (Roadmap) | Closure erst nach M6 |

### Sensoren (vom Implementer vor dem Handoff, Exit-Codes der `make`-Prozesse in §3 *Umsetzung*)

`make gates`, `make ci`, `make fullbuild`, `make test-docker` (mit `-v`: Zahl von PASS, SKIP und FAIL aus dem Log; ein lokaler Lauf ersetzt den CI-Lauf nicht; am Paket-Zeitlimit von 30 Minuten endete der lokale Lauf in drei Cuts mit Exit 2, [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md); ein Netzfehler im Bau ist ein eigener Befund, [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md)), `make build VERSION=0.15.0` mit `docker run --rm u-boot:latest --version` = `u-boot version 0.15.0` (lokal, nicht gepusht), `make doc-immutable RANGE="v0.14.0..HEAD"`, `make doc-trace`, `make doc-commits RANGE="origin/main..HEAD"`, `make doc-planning`; M-T. `make verify-depguard` entfällt (keine Importregel berührt), mit Begründung zu nennen.

### Umsetzung (Implementer)

*(Wird vom Implementer gefüllt: Commits, M0, Abweichungen, Sensoren mit Exit-Codes, M-T-Protokoll, Abgleich der Plan-Tabelle, Kommentar-Selbstcheck nach [`AGENTS.md`](../../../../AGENTS.md) §3.9.)*

## 4. Trigger

**Start** (`open` → `in-progress`): (1) Die fünf übrigen Slices der Welle [`welle-forge-cli-gh-glab`](../welle-forge-cli-gh-glab.md) liegen in `done/` und der Validator-Bericht liegt vor (Welle-Datei §3 Punkt 7; eingetreten am 2026-10-23, `07353780`); (2) die Nummer `v0.15.0` gilt (Projektinhaber, Welle-Datei §5, Frage 1); (3) `Verantwortlich:` ist gesetzt; (4) `in-progress/` enthält nur die Master-Dokumente (WIP-Limit 1; gelesen am 2026-10-23: `carveouts.md`, `README.md`, `roadmap.md`); (5) die Entscheidungen des Starts in §7 sind getroffen (Nummer, kein Vorab-Tag). Der Weg ist `open/` → `in-progress/` per reinem `git mv` (`make slice-mv SLICE=slice-v1-release-cut-v0.15.0 TO=in-progress`); der Roadmap-Zeiger geht in einen eigenen Commit direkt nach dem Move, beide im selben Push.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht vorgesehen; der Cut sind wenige Commits plus Tag.
- `in-progress` → `open` (blockiert): ein roter `image-scan` (neuer Advisory: Go- oder Trivy-Pin heben, eigener Vorgang davor); ein Befund des Trockenlaufs, der eine Änderung an Workflow oder Skript verlangt; ein `--- FAIL` oder ein Zeitlimit-Abbruch in `integration-docker` auf dem Push-Tip; eines der neuen Docker-Tests ist `--- SKIP` statt `--- PASS`; ein Befund des Reviewers oder Verifiers, der Produktcode verlangt; ein `publish`-Lauf, dessen Fehler einen Patch-Release braucht.

## 5. Closure-Trigger

DoD vollständig: Tag `v0.15.0` gesetzt (M5) und gepusht, `publish` und Job `tap` kontrolliert (M6), Trockenlauf M-T mit Ergebnis und Grenze festgehalten, Closure-Notiz mit Tag-Hash, Lauf-IDs, Nacharbeit (M8) und Steering-Loop-Lerneintrag, Verification Evidence (Verifier, getrennt von Implementer und Reviewer). Der Slice geht **erst nach dem Tag** nach `done/` (erst Inhalt, dann reiner `git mv`, [`AGENTS.md`](../../../../AGENTS.md) §3.3); der Ruhe-Marker der Roadmap steht erst dann wieder. Kein Tag ohne ausdrückliche Freigabe: Gibt der Projektinhaber den Tag nicht frei, bleibt der Slice in `in-progress/` oder geht mit benanntem Grund nach `open/` zurück.

**Danach: Welle-Self-Close (`close-welle`), nach dem Release-Cut-Commit und dem Move dieses Slice nach `done/`; nicht Teil dieses Slice.** Die Welle schließt mit diesem Slice als Detailquelle ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)). Der Planner plant sie so ein:

1. Welle-Datei §7 füllen: Ergebnis mit Zeiger auf diesen `done/`-Slice, jeder Closure-Trigger (§3 Punkte 1 bis 9) mit Beleg; Trigger-Audit (jeder Trigger eingetreten, beobachtbar), Zeiger vom Ruheort `done/`.
2. Lese-Schritt des Beobachtungs-Registers, Evidence-Datei `welle-forge-cli-gh-glab.md` (je Vorgang zählt die Welle einmal): [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (Beleg M-1), [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) (Beleg L-3), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Beleg M-2); L-2 ohne Eintrag. Die Zählung folgt aus den Dateien; die Closure setzt keinen Zähler.
3. Die Einträge **über der Schwelle** (Welle-Datei §8: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md), `vertrag-schweigt-zum-nutzerpfad`, dazu `kennung-in-cli-ausgabe`, `e2e-paket-am-zeitlimit`): Der Ausgang liegt beim Projektinhaber; die Closure legt ihm die Vorschläge vor, weist keinen zu und trägt sonst die benannte Grenze „Entscheidung des Projektinhabers steht aus“ (Welle-Datei §3 Punkt 9).
4. Steering-Loop-Eintrag der Welle (Lerneintrag: geschärfte Regel · neuer Sensor · benannte Spec-Lücke), einschließlich des Regel-Vorschlags „Breite des Satzes gegen Breite des Tests“ als Entscheidung des Projektinhabers; verkörpert wird nichts ohne dessen Entscheidung.
5. Roadmap (Zeilen der Welle und von `v0.15.0` unter *Abgeschlossene Wellen*, Zeiger weg, Ruhe-Marker zurück), `git mv` der Welle-Datei nach `done/`, Verweis-Nachzug (eigener Commit), `make gates` grün.

## 6. Risiken und offene Punkte

Der Ausgang jedes Risikos wird bei der Closure festgeschrieben.

- **Der Lead sagt mehr, als gemessen ist** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10). Besonders A-1: „kein authentifizierter Aufruf gemessen“, der Kalt-Bau mit frischem Download ist im Validator-Lauf nicht gefahren (Cache-Treffer; Punkt 2 des Berichts). Gegenmittel: Beleg-Tabelle (§3); Formulierungen „nicht gemessen“, „ein Lauf“. — **Ausgang:** *(bei der Closure)*
- **Der Satz zu M-2 sagt mehr als der Code** („`config set` prüft es noch nicht“ ist eine Aussage über Verhalten). Gegenmittel: Die Aussage stützt sich auf eine Messung des Validators (ein Lauf) und den gelesenen Plan des Folge-Slice; der Implementer liest den Pfad `config set` → Schreiben noch einmal am Code (Plan-Notiz V-3 zu `SetMany`) und nennt die Fundstelle; der Satz benennt es als Messung, nicht als Zusage. — **Ausgang:** *(bei der Closure)*
- **Die Release-Aussage zur CI weicht vom Lauf auf dem Push-Tip ab** (Zahl übersprungener Tests, Runner-Wechsel auf Ubuntu 26; [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)). Gegenmittel: M2 zählt `--- SKIP`; kein CI-Satz im Lead vor M2. — **Ausgang:** *(bei der Closure)*
- **Zeit- und Netz-Risiko der neuen Docker-Tests** (vier Bau-Tests laden Archive von `github.com` und `gitlab.com`; [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md), [`BEO-E2E/docker-test-scheitert-am-externen-netz`](../observations/BEO-E2E/docker-test-scheitert-am-externen-netz/observation.md), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md)). Gegenmittel: M2 hält die Dauer fest; kein Eingriff in Limits ([`AGENTS.md`](../../../../AGENTS.md) §3.6); ein lokaler Netzfehler wird mit Einzellauf und Dateiliste belegt, nicht mit Wortsuche (Lehre des Cuts `v0.14.0`). — **Ausgang:** *(bei der Closure)*
- **Ein neuer Advisory macht `image-scan` vor dem Tag rot.** Gegenmittel: `make ci` lokal und `image-scan` auf dem Push-Tip; Rückführung nach `open` (§4). — **Ausgang:** *(bei der Closure)*
- **Der Job `tap` scheitert im stabilen Lauf**, obwohl M-T grün war. Gegenmittel: `tap-check` unmittelbar vor dem Tag (M4). — **Ausgang:** *(bei der Closure)*
- **Das Tag-Datum weicht vom festen Datum ab** (vier Stellen). Gegenmittel: M0; Cut, Push und Tag am selben Tag. — **Ausgang:** *(bei der Closure)*
- **Die Plan-Tabelle ist unvollständig oder trägt eine Aussage, die der Code nicht trägt** ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10). Gegenmittel: Abgleich nach der letzten Fix-Runde. — **Ausgang:** *(bei der Closure)*
- **Host-Werkzeuge außerhalb von Docker** ([`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md)): Textänderungen über das Edit-Werkzeug, kein `sed -i` und kein Interpreter auf dem Host; M-T läuft im Container. — **Ausgang:** *(bei der Closure)*
- **Der Welle-Self-Close liegt vor dem Cut-Commit oder ohne ihn** (Reihenfolge: Closure-Trigger 8 der Welle verlangt den geschnittenen Release). Gegenmittel: §5 ordnet ihn nach dem Move dieses Slice nach `done/`. — **Ausgang:** *(bei der Closure)*

## 7. Closure-Notiz

**Entscheidungen beim Start** (2026-10-23, Implementer nach dem Auftrag des Hauptagenten; der Wortlaut einer Freigabe des Projektinhabers liegt nicht vor): (1) Nummer `v0.15.0` (bestätigt, Welle-Datei §5); (2) kein Vorab-Tag, M3 entfällt (Vorschlag des Planners, Planvorgabe); (3) kein CI-Satz im Lead vor M2 (Vorschlag des Planners, Planvorgabe); (4) **Reichweite:** Push (M1) und Tag (M5) führt nur der Hauptagent, je Vorgang nach ausdrücklicher Freigabe des Projektinhabers; der Implementer pusht und taggt nie, auch nicht lokal.

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders · Steering-Loop-Eintrag (Lerneintrag; `liegt in` nur bei Verkörperung) · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `docs/plan/`, `docs/user/`, `docs/maintainer/` (nur `releasing.md`, bei der Closure) und `docs/reviews/` (DOK), dazu die Versionszeilen in `cmd/uboot/main.go` (`cmd/uboot` ist keine Sub-Area, Modus-Deklaration in [`harness/conventions.md`](../../../../harness/conventions.md)), `repo.mk` und `Dockerfile` (Werkzeug-Dateien, WZG). `docs/user/` und `docs/maintainer/` sind in der Modus-Deklaration nicht eigens geführt; die früheren Cuts ordneten sie wie hier DOK zu. Keine Sub-Area wird neu geschnitten.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen ([`../observations/`](../observations/README.md); gelesen am 2026-10-23). Einschlägig und in §6 geführt: `zusage-weiter-als-messung`, `spec-bedingung-ohne-code-entsprechung`, `plan-aussage-ueber-code`, `ci-log-falsch-gelesen`, `host-werkzeug-ausserhalb-docker`, `ci-ueberspringt-nested-podman-build`, `e2e-paket-am-zeitlimit`, `docker-test-scheitert-am-externen-netz`, `test-gruen-nur-mit-lokalem-cache`, `podman-static-pin-alterung` (Pins als Stand). Die Einträge über der Schwelle erhalten bei diesem Cut keinen Ausgang (§5, Welle-Self-Close, Punkt 3). Zähler bei der Closure nicht gesetzt.

Alle berührten Sub-Areas GF.
