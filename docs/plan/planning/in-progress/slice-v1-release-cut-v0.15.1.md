# Slice V1: Release-Cut `v0.15.1`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** keine Welle-Datei — Patch-Release `v0.15.1` von `main` (Release-Version nach [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md); wie `v0.13.1` und `v0.8.1` ohne Welle-Datei, Zeile unter *Nächste Wellen* der [Roadmap](../in-progress/roadmap.md)). Entscheidung des Projektinhabers: ein Patch-Release (Antwort „ok patch“, 2026-10-24). Der Release trägt die zwei Fixes, die in `CHANGELOG.md` unter `## [Unreleased]` stehen (gelesen: `CHANGELOG.md` `:12`–`:17`, zwei Einträge *Fixed*, kein *Security*): [`slice-v1-generate-prueft-freigabeliste`](../done/slice-v1-generate-prueft-freigabeliste.md) und [`slice-v1-config-set-validiert-feature-version`](../done/slice-v1-config-set-validiert-feature-version.md). Stand beim Schreiben (2026-10-24, gelesen): `HEAD` = `d98f5755`, Arbeitsbaum sauber, `origin/main..HEAD` leer (`git rev-list --count`), 37 Commits seit dem Tag `v0.15.0` (`a895ffd2`, `git rev-parse v0.15.0^{commit}`), `in-progress/` ohne Slice, Lastenheft `0.22.0` (`spec/lastenheft.md` Zeile 9, im Verify-Report beider Slices als unverändert gelesen). Der Produktcode unterscheidet sich von `v0.15.0` nur durch die zwei Fixes: `git diff --stat v0.15.0..HEAD -- cmd internal` trifft neun Dateien (drei Nicht-Test-Dateien: `internal/hexagon/application/config.go`, `devcontainer_features.go`, `generate.go`). `.github/`, `scripts/`, `packaging/`, `Dockerfile`, `repo.mk`, `Makefile`, `go.mod`, `tools/` und `d-check.mk` haben seit `v0.15.0` keinen Diff; in `harness/` steht eine Zeile in `replay.md`, in `docs/maintainer/` die zwei Zeilen von `releasing.md` (Zeiger „zuletzt“) (gelesen: `git diff --stat v0.15.0..HEAD -- <diese Pfade>`).

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Freigabe externer Feature-Quellen; Scope: die Zusage, die der Release ausliefert), [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (`config set` prüft wie `generate`), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code `10`), [`ADR-0004`](../../adr/0004-ci-system.md) (CI), [`ADR-0007`](../../adr/0007-distributionswege-ghcr.md), [`ADR-0016`](../../adr/0016-homebrew-distribution-per-tap.md), [`ADR-0017`](../../adr/0017-linux-pakete-nfpm.md). Ablauf: [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md). Muster: [`slice-v1-release-cut-v0.13.1`](../done/slice-v1-release-cut-v0.13.1.md) (Patch-Cut) und [`slice-v1-release-cut-v0.15.0`](../done/slice-v1-release-cut-v0.15.0.md) (letzter Cut, Versionsstellen). Belege der Inhalte: [Verify-Report `generate`](../../../reviews/2026-10-24-slice-v1-generate-prueft-freigabeliste-verify.md) und [Verify-Report `config set`](../../../reviews/2026-10-24-slice-v1-config-set-validiert-feature-version-verify.md) (beide „schließbar“).

**Berührte Spec-Stellen:** — (Spec-Version ≠ Release-Version; dieser Slice ändert `spec/` nicht; [`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) und [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) tragen die Fixes schon, gelesen in `git diff v0.15.0..HEAD -- spec/spezifikation.md`).

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-24); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

**Ziel:** `v0.15.1` liefert als Patch-Release von `main` zwei Fixes der Prüfung von Feature-Einträgen in `u-boot.yaml` — `u-boot generate devcontainer` weist die Quelle eines aktivierten Features ab, die nicht freigegeben ist; `u-boot config set …features.<feature>.version` weist bei `ollama`, `gh` und `glab` eine `version` ohne `source` ab — als GHCR-Image, Binaries, Pakete und Homebrew-Formel, mit einem Eintrag *Fixed* im CHANGELOG, der nicht weiter reicht, als gelesen und gemessen ist. Dieser Slice bereitet vor und lässt den Tag zu; **Push (M1) und Tag (M5) führt nur der Hauptagent, je mit eigener ausdrücklicher Freigabe des Projektinhabers im Gespräch** (§3, *Reihenfolge der Maintainer-Aktionen*).

### Versionsnummer: `v0.15.1` (Entscheidung des Projektinhabers, 2026-10-24)

Gelesen: `CHANGELOG.md` `:12`–`:17`, die beiden Verify-Reports, `spec/lastenheft.md` Zeile 9.

- **Patch, nicht Minor:** kein neuer Befehl, keine neue Option; beide Fixes stellen eine schon zugesagte Prüfung her ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)). Eine **Verhaltensänderung für Skripte** gibt es: `generate devcontainer` endet bei einer nicht freigegebenen Quelle mit `10` statt `0`, `config set` bei `version` ohne `source` mit `10` statt `0`; der Lead nennt sie als solche.
- **Einstufung *Fixed*, kein *Security* und kein Advisory (M5b entfällt):** Der Reviewer sah keinen Bruch einer Vertrauensgrenze (Review des Fix-Slice `generate`, INFO-1; gelesen in dessen §7, *Einschätzung des Reviewers zur Sicherheit*: `featureSources.allow` und `features.<n>.source` stehen in derselben Datei `u-boot.yaml`; die Prüfung schützt vor dem Versehen und vor dem Umweg an `config set` und `doctor` vorbei). **Die Einstufung liegt beim Projektinhaber.** Will er *Security* und ein Advisory, sagt er es; dann kommen ein Eintrag `### Security` nach einer Beleg-Tabelle und der Schritt M5b hinzu (Muster [`slice-v1-release-cut-v0.13.1`](../done/slice-v1-release-cut-v0.13.1.md) §3); dieser Plan ändert sich dann, der Implementer entscheidet das nicht.
- **Kein Vorab-Tag:** direkt stabil (M3 entfällt, Vorschlag des Planners wie bei `v0.13.1` und `v0.15.0`; bestätigt der Projektinhaber beim Start, §7).

### Inhalt des Release (gelesen am 2026-10-24, `CHANGELOG.md` Block `[Unreleased]`)

| Eintrag im CHANGELOG (Zeile) | Art | Hinweis für den Lead |
|---|---|---|
| `generate devcontainer` prüft die Freigabeliste der Feature-Quellen: Exit `10`, Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), nichts geschrieben, auch mit `--dry-run`, in beiden Profilen; vorher Exit `0` und die unlistete Quelle in `devcontainer.json`; zwei Auswege (`:16`) | Fixed | Verhaltensänderung für Skripte; Projekte mit einem solchen Eintrag brauchen einen der zwei Auswege; die Messung am Binary (Verify `generate`, *Ende-zu-Ende*) ist Linux `amd64` |
| `config set …features.<feature>.version` (auch mehrere Paare) prüft wie `generate`: `version` ohne `source` bei `ollama`, `gh`, `glab` endet mit `10`, Datei unverändert, auch mit `--dry-run` und `--diff`; eine schon gesperrte Datei bleibt abgelehnt, Abhilfe von Hand (`:17`) | Fixed | schließt die Grenze, die der Lead von `[0.15.0]` nennt; der Text dort bleibt (Geschichte jenes Tags) |

### Aufträge an den Release-Cut

Jede neue Aussage in CHANGELOG, READMEs und Handbuch braucht einen Beleg nach [`AGENTS.md`](../../../../AGENTS.md) §3.10 (Beleg-Tabelle in §3); was keinen hat, steht nicht im Text oder steht als „nicht gemessen“. Kein Auftrag verlangt Produktcode.

1. **M1-Vorbereitung — Versionsstrings und Kopf des CHANGELOG:** `0.15.0-dev` → `0.15.1-dev` in `cmd/uboot/main.go`, `repo.mk` und `Dockerfile`; `[Unreleased]` → `[0.15.1] - <Datum M0>` mit Lead (Soll in §3), `[Unreleased]` bleibt leerer Anker, Vergleichslinks; READMEs (Status-Block, Tabellenzeile, `VERSION=`), Handbuch (Kopf, Versionsbeispiele, Gültigkeitsbereich, `Stand`).
2. **Zusage-Menge und Grenzen der Messung** im Lead: gemessen am gebauten Binary auf einem Linux-Host (Verify-Reports; ohne `amd64`, siehe Umsetzung Punkt 2) ist getrennt je Fix: `config set` für `gh`, `glab` und `java`, `ollama` nur durch Tests auf Service-Ebene gepinnt; `generate` mit einem Eintrag aus einer nicht freigegebenen Quelle (`evil`), im Standard- und im Sandbox-Profil, mit den drei Auswegen (Freigabe per `config set`, Flag `--allow-external-feature-sources`, Feature abschalten) und dem Katalogschlüssel `node` ohne `source` (Verify-Report `generate`, *Ende-zu-Ende*; `gh`, `glab`, `java` kommen dort nicht vor) (Trennung in der Fix-Runde nach LOW-V1); Mutationstests der Verifier; **ungemessen** (Verify `generate`, *Ungemessen*): Whitespace am Rand einer Quelle in `u-boot.yaml`, `doctor` bei fehlendem `enabled`; (Verify `config set`, *Ungemessen*): `--json`-Ausgabe der neuen Ablehnung, Mac, Podman, andere Hosts. Der Lead sagt nichts darüber hinaus.
3. **CI-Lesung** (Muster Cut `v0.15.0`): `ci` und `integration-docker` auf dem Push-Tip über Job-Log (REST-API) und Check-Run-Annotationen lesen, nicht per Textsuche ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)). **Keine CI-Aussage im Release-Text vor M2.**
4. **Wie jeder Cut** ([`releasing.md`](../../../maintainer/releasing.md) §3 bis §5): `releasing.md`-Zeiger „zuletzt“ und Roadmap sind Closure-Schritte (§3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Push, Tag und jede Aktion mit Außenwirkung durch den Implementer oder Planner.** Extern wirksam, nicht sauber zurückzunehmen ([`releasing.md`](../../../maintainer/releasing.md) *Zweck*). **M1 (Push) und M5 (Tag) führt nur der Hauptagent, je Vorgang mit eigener Freigabe des Projektinhabers im Gespräch;** eine Freigabe für M1 ist keine für M5, eine Nachricht eines Agenten ist keine Freigabe. Der Tag `v0.15.0` bleibt unberührt (`a895ffd2`).
- **Ein Eintrag *Security*, ein Advisory (M5b).** Die Einstufung liegt beim Projektinhaber (oben); der Plan nimmt sie nicht vorweg. Ein Wort „security“ oder „vulnerability“ steht im Lead nicht.
- **Produktcode, Meldungen, Hilfe-Texte, Spezifikation, Lastenheft, ADRs.** Die Aufträge sind Text in Release-Dokumenten. Die Versionszeilen (`var version`, `VERSION ?=`, `ARG UBOOT_VERSION`) sind die einzigen Zeilen mit Wirkung im Binary.
- **Nachträgliche Änderung der Texte unter `[0.15.0]`, der Tabellenzeile `v0.15.0` beider READMEs und der Welle `welle-forge-cli-gh-glab`:** sie beschreiben den Stand jenes Tags („`config set` prüft es noch nicht“) und bleiben wahr; `v0.15.1` korrigiert im eigenen Eintrag. Ein Tag wird nicht verschoben ([`releasing.md`](../../../maintainer/releasing.md) §8).
- **Die offenen Folgepunkte der Fixes:** der Wortlaut „managed block missing or malformed“ bei einer gesperrten Datei, `config unset`, L-2 und L-3 ([`slice-v1-config-set-validiert-feature-version`](../done/slice-v1-config-set-validiert-feature-version.md) §7) sowie N1 und N2 ([`slice-v1-katalog-features-quellen-und-freigabe`](../open/slice-v1-katalog-features-quellen-und-freigabe.md)) und N4 ([`slice-v1-docker-cli-bau-auf-standard-basis-image`](../open/slice-v1-docker-cli-bau-auf-standard-basis-image.md)): andere Ursache, Entscheidung des Projektinhabers, kein Teil der Fixes.
- **Pins und Pipeline:** kein Trivy-Pin-Bump, kein Go-Pin-Bump, keine Änderung an Workflows, Skripten, `Makefile`, Zeitlimits oder Gates ([`AGENTS.md`](../../../../AGENTS.md) §3.6); ein roter `image-scan` durch einen neuen Advisory hebt den Go-Pin als eigener Vorgang davor ([`releasing.md`](../../../maintainer/releasing.md) §4).

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe, die Kontrollen unter Punkt 2 und die Closure-Pflichten darunter zählen nicht mit.

- [x] **Release-Text und Dokumente:** `CHANGELOG.md` trägt `## [0.15.1] - <Tag-Datum>` mit Lead nach §3 (*Aussage des Release*) und den Vergleichslinks, `## [Unreleased]` bleibt als leerer Anker; Versionsstrings `0.15.1-dev` an den drei Stellen; beide READMEs (Status-Block, Tabellenzeile, `VERSION=`) und das Handbuch (Kopf, Versionsbeispiele, `Stand`, Gültigkeitsbereich) tragen `v0.15.1`; das Datum steht an den vier Stellen (CHANGELOG, Tabellenzeile beider READMEs, Handbuch-`Stand`) gleich dem Tag-Datum. Jede neue Aussage steht mit ihrem Beleg in der Beleg-Tabelle (§3); Handbuch und READMEs tragen keine Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] **Tag und Kontrolle:** Vor dem Tag sind `make gates`, `make ci` und `make fullbuild` grün, `make test-docker` ist lokal nicht gelaufen (Planänderung gegenüber der Fassung „gelaufen vor dem Tag“, Review-Fix-Runde 2026-10-24: Paket-Zeitlimit und Speicher des Hosts, Produktcode-Diff gegenüber `v0.15.0` eine Zeile, Punkt 8 in §3); die Docker-Tests laufen in M2 im CI auf dem Push-Tip, dessen Zahl der übersprungenen Tests in §3 steht, der Trockenlauf des Tap-Pfads (M-T) ist festgehalten, `tap-check` ist grün, und `ci` und `integration-docker` auf dem Push-Tip sind grün (M2). **Erst danach und nur nach ausdrücklicher Freigabe des Projektinhabers im Gespräch:** Tag `v0.15.1` gesetzt und gepusht (M5). Kontrolle nach [`releasing.md`](../../../maintainer/releasing.md) §6 (M6).
  - [ ] Release-Workflow `publish`: Jobs `publish`, `tap`, Paket-Smoke rpm/fedora und deb/ubuntu `success` (Lauf-ID in §7).
  - [ ] GHCR-Image: `docker run --rm ghcr.io/pt9912/u-boot:0.15.1 --version` meldet `0.15.1`; `:0.15.1` und `:latest` tragen denselben Digest.
  - [ ] Binaries und Pakete: `gh release view v0.15.1 --json assets,isPrerelease,isDraft` zeigt sechs Binaries, vier Pakete, `SHA256SUMS` und `u-boot.rb` (zwölf Assets), kein Prerelease, kein Draft.
  - [ ] Release-Notes: der GitHub-Release verweist auf `CHANGELOG.md` §[0.15.1] (Auto-Notes) oder trägt vom Maintainer geschriebene Notes; der Abschnitt `[0.15.1]` im CHANGELOG auf `main` ist der des Tags.
  - [ ] Job `tap` wirklich gelaufen (Log „tap-nachzug: … updated to v0.15.1.“ oder „already at“, keine Warnung zum fehlenden Secret, Annotationen gelesen).
- [ ] **Doku nachgezogen:** `releasing.md` (Zeiger „zuletzt“) und Roadmap (Zeile `v0.15.1` aus *Nächste Wellen* nach *Abgeschlossene Wellen* mit Tag-Hash und Datum, Ruhe-Marker) tragen den Ist-Stand; die abgelaufenen Zustandsmarker in `open/README.md` (§3) auch.
- [ ] `make gates` grün (nach dem letzten Commit vor dem Handoff und nach der Closure).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); Verification-Report getrennt davon (Verifier), vor M1.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Jede Aussage in dieser Tabelle ist an der genannten Stelle gelesen (Stand 2026-10-24, `HEAD` = `d98f5755`, [`AGENTS.md`](../../../../AGENTS.md) §3.10). Zeilennummern nennen den Fundort beim Schreiben; der Implementer sucht sie vor der Änderung neu. Vor dem Handoff und nach jeder späteren Fix-Runde gleicht er die Tabelle gegen `git diff --stat <Plan-Commit>..HEAD` ab — einschließlich dieser Datei, der Vermerke in `done/` und `open/`, der Reports und der Verweise, die `make slice-mv` nachzieht; eine Datei, die er als „kein Diff“ führt, belegt er durch eine gelesene Aussage, nicht durch eine Suche.

**Versionsstellen** (gelesen mit `grep -rIn -e "0\.15\.0"` über den Baum außerhalb von `.git`, `.harness`, `docs/plan/planning`, `docs/reviews`, `docs/messungen` und `CHANGELOG.md`): `cmd/uboot/main.go:43`, `repo.mk:51`, `:54`, `:55`, `Dockerfile:35`, `:37`, `:39`, `README.md:10`, `:102`, `:209`, `README.de.md:11`, `:104`, `:212`, `docs/user/benutzerhandbuch.md:4`, `:82`, `:108`, `:135`, `:2998`; `docs/maintainer/releasing.md:57` (Zeiger „zuletzt“, Closure). **Keine Versionsstelle** (gelesen): `spec/lastenheft.md:1591` (Historie, Lastenheft-Version `0.15.0`) und `spec/spezifikation.md:1435` (Verweis auf Lastenheft `0.15.0`) sind Spec-Versionen, nicht die Release-Version. Der Implementer wiederholt die Suche vor und nach dem Cut, auch mit `-e "0.15.0-dev"` für die Kommentare.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `CHANGELOG.md` | update (Cut) | `:12` `## [Unreleased]` wird zum leeren Anker; der Inhalt (`:16`–`:17`) wandert unter `## [0.15.1] - <Datum M0>` mit Lead „Patch release.“ (Stil von `[0.13.1]`; Soll unten), die zwei Einträge *Fixed* bleiben wörtlich. Fuß: `:1518` `[Unreleased]: …/compare/v0.15.0...HEAD` → `v0.15.1...HEAD`; neue Zeile `[0.15.1]: …/compare/v0.15.0...v0.15.1` vor `:1519` (`[0.15.0]`). Der Text unter `[0.15.0]` (Lead-Absatz „Limit of `version:`“, Eintrag *Added*) bleibt (gelesen: er beschreibt den Stand jenes Tags). |
| `cmd/uboot/main.go` | update | `:43` `var version = "0.15.0-dev"` → `0.15.1-dev`. Einzige Zeile mit Wirkung im Binary; `main_test.go` pinnt die Nummer nicht (gelesen: `:18`–`:25` und `:68`–`:69` prüfen nur „nicht leer“; der Implementer liest es erneut). |
| `repo.mk` | update | `:51` und `:54` (Kommentare) und `:55` `VERSION ?= 0.15.0-dev` → `0.15.1-dev`. |
| `Dockerfile` | update | `:35` und `:37` (Kommentare) und `:39` `ARG UBOOT_VERSION=0.15.0-dev` → `0.15.1-dev`. Kein weiterer Diff (kein Go-Pin-Bump, §1). |
| `README.md` | update | Status-Block `:10`: Version `v0.15.0` → `v0.15.1`; die Katalogzeile (`:13`) bleibt (gelesen: sie nennt `gh`/`glab`-Zusage, nichts zu `version` und Freigabe); ein Satz in Worten für den Benutzer zu den zwei Fixes, ohne Datum, ohne „released“, ohne Kennungen (Soll unten). `:102` `VERSION=0.15.0` → `0.15.1`; neue Zeile `v0.15.1` in der Releases-Tabelle nach `:209` (Zeile `v0.15.0`; die Tabelle läuft aufsteigend bis `v0.15.0`, dann `v0.5.0`). **Kein Diff:** die Tabellenzeile `v0.15.0` (`:209`, beschreibt den Tag-Stand). |
| `README.de.md` | update | gleiches Muster: `:11` Version, Satz im Status-Block, `:104` `VERSION=0.15.0`, Tabellenzeile nach `:212`; Tabellenzeile `v0.15.0` bleibt. |
| `docs/user/benutzerhandbuch.md` | update | Kopf `:4` `Software-Version: v0.15.0` → `v0.15.1`, `:5` `Stand: 2026-10-23` → Tag-Datum; `:82` `u-boot version 0.15.0`, `:108` `download/v0.15.0`, `:135` `VERSION=0.15.0`, `:2998` „beschreibt `u-boot` v0.15.0“ → `0.15.1`. **Kein Diff, gelesen (kein Nachzug nötig):** der Absatz zur Freigabe bei `generate devcontainer` (`:1118`–`:1127`, „Exit-Code 10, schreibt nichts (auch mit `--dry-run` nicht) … Standard- und Sandbox-Profil“), die Absätze zu `version` ohne `source` (`:1965`, `:2062`–`:2075`, mit `config set` Exit `10`, Datei unverändert, auch mit `--dry-run` und `--diff`, und der Abhilfe für eine schon gesperrte Datei) und die Pfad-Tabelle (`:2417`, „bei `ollama`, `gh` und `glab` nur mit `source`“) tragen beide Fixes schon (geändert in `0cdec3b2` und `00e3ca3d`). Kennungsfrei. |
| `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md` | kein Diff (gelesen) | `lastenheft.md` Zeile 9 `0.22.0` bleibt; `spezifikation.md` trägt die Fixes ([`SPEC-040`](../../../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) Wertetabelle und Stufe 6, [`SPEC-050`](../../../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen) Absatz *Freigabe der Quelle (beide Profile)*, [`SPEC-051`](../../../../spec/spezifikation.md#spec-051--features-gh-und-glab-im-devcontainer-installation-meldungen-hosts-und-zugangsdaten) Absatz `version` und `source`, Historie-Zeile 2026-10-24, gelesen im Diff `v0.15.0..HEAD`); `architecture.md` ohne Diff. |
| `docs/maintainer/releasing.md` | kein Diff im Cut; update bei der Closure | `:56`–`:57` nennt als „zuletzt“ den Cut `v0.15.0` (gelesen): richtig bis zum Abschluss; die Closure setzt den Zeiger auf diesen Slice (Muster: Closure von [`slice-v1-release-cut-v0.15.0`](../done/slice-v1-release-cut-v0.15.0.md)). Die übrigen Aussagen §1 bis §8 stimmen mit dem Ist-Stand (gelesen; §3 bis §7 kennen kein Advisory, nicht Teil dieses Cuts). |
| `docs/plan/planning/in-progress/roadmap.md` | update (Plan-Commit: Zeile unter *Nächste Wellen*; Start; Closure) | Plan-Commit: neue Zeile `v0.15.1` unter *Nächste Wellen* (Tabelle `:24`–`:28`; Muster: die Zeile `v0.13.1`, die dort bis zum Tag stand), Trigger „Projektinhaber: Patch-Release (2026-10-24)“, Slice dieser Cut. `:20` `Nichts in Arbeit.` (Ruhe-Marker, `.d-check.yml` `planning`). Start: der Marker entfällt, an seine Stelle tritt ein Zeiger auf diesen Slice (eigener Commit direkt nach dem Move). Nach dem Tag (Closure): Zeiger weg, `Nichts in Arbeit.` zurück, die Zeile wandert nach *Abgeschlossene Wellen* nach der Zeile `v0.15.0` (`:104`) mit Datum, Tag-Hash und Verweis auf diesen Slice. `:3` nur, wenn das Datum abweicht. Kein Eintrag unter *Historische Trigger-Verschiebungen* (keine Umplanung). |
| `docs/plan/planning/open/README.md` | update (Plan-Commit; Closure) | neue Zeile für diesen Slice am Ende der Tabelle (nach `:68`, Zeile `slice-v1-release-cut-v0.15.0`). Abgelaufene Marker: Zeile `slice-v1-generate-prueft-freigabeliste` (`:58`) „Release-Zuordnung offen (beim Projektinhaber)“ → „Release `v0.15.1`“; Zeile `slice-v1-config-set-validiert-feature-version` (`:67`) „**nicht vor dem Release-Cut `v0.15.0`**“ → „Release `v0.15.1` (nach dem Cut `v0.15.0`)“. Die Moves ziehen die Links der Zeile dieses Slice nach. Closure: „liegt in `done/`“, Tag. |
| `docs/plan/planning/done/slice-v1-generate-prueft-freigabeliste.md`, `…/slice-v1-config-set-validiert-feature-version.md` | kein Diff (gelesen) | Ihre §7 nennen die Release-Zuordnung als offen und die Einstufung beim Projektinhaber; das ist der Stand ihrer Closure (2026-10-24) und bleibt als Geschichte. Der Cut führt die Zuordnung in diesem Slice und in der Roadmap. |
| `docs/plan/planning/done/slice-v1-release-cut-v0.15.1.md` (dieser Slice; Plan in `open/`, später `done/`) | neu, update | Plan (Planner), Moves per `make slice-mv` nach `in-progress/` und nach `done/` (reine Commits, [`AGENTS.md`](../../../../AGENTS.md) §3.3), Umsetzung in §3 (Implementer), Closure (Planner): §2, §3, §6, §7. |
| `docs/reviews/<Datum>-slice-v1-release-cut-v0.15.1.md`, `…-verify.md` | neu, durch Reviewer und Verifier; Verweis-Nachzug | Review-Report und Verification-Report (vor M1; Dateiname `<Datum>-<Slice-Kennung>[-verify].md`). |
| `docs/plan/planning/observations/…` | nur Closure | Stand-Zeilen und Evidence-Dateien schreibt die Closure; Kandidaten aus dem Register (§8). |
| `.github/workflows/*`, `scripts/`, `packaging/`, `Makefile`, `go.mod`, `tools/`, `d-check.mk`, `docs/maintainer/quality.md`, `internal/`, übrige `cmd/uboot/` | kein Diff | Kein Pin-Bump, keine Pipeline-Änderung, kein Produktcode (§1). Der Implementer belegt „kein Diff“ vor dem Handoff durch eine gelesene Aussage (Statuszeile, Pin-Zeile), nicht durch eine Suche. |

### Beleg-Tabelle der neuen Aussagen ([`AGENTS.md`](../../../../AGENTS.md) §3.10)

Der Implementer schreibt nur, was hier einen Beleg hat, in der Form der letzten Spalte; eine Aussage ohne Zeile steht nicht im Text.

| Aussage (Soll) | Ort | Beleg | Form im Text |
|---|---|---|---|
| `generate devcontainer`: Quelle eines aktivierten Features außerhalb der Freigabe (Liste oder Flag desselben Aufrufs) → Exit `10`, Code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), nichts geschrieben, auch mit `--dry-run`, in beiden Profilen | Lead, Status-Blöcke | Verify `generate`, *Ende-zu-Ende am gebauten Binary* (Tabelle: Standard- und Sandbox-Profil, `--dry-run`, `--json`), Mutationen M1, M3, M4, M10 | Zusage mit Messort („Linux `amd64`, gebautes Binary“) |
| Vorher: Exit `0` und die Quelle in `devcontainer.json`; `doctor` und `config set …source` wiesen sie ab | Lead (steht schon im Eintrag `:16`) | `CHANGELOG.md` `:16`; Slice `generate` §7 (am Release-Image `0.15.0` mit Exit `0` gezeigt, 2026-10-24) | wie im Eintrag; „am Release-Image `0.15.0`“ nur, wenn der Slice es so trägt (gelesen am Slice, nicht aus dem Gedächtnis) |
| Projekte mit solchem Eintrag brauchen einen der zwei Auswege: Quelle freigeben oder Feature abschalten | Lead | Verify `generate`, *Ende-zu-Ende* (Abhilfen `config set …featureSources.allow`, `--allow-external-feature-sources`, `…enabled false`, je Exit `0`) | Verhaltensänderung für Skripte, als solche genannt |
| `config set …version` (auch mehrere Paare) ohne `source` bei `ollama`, `gh`, `glab` → Exit `10`, Datei unverändert, auch mit `--dry-run` und `--diff`; `java` ohne `source` geht | Lead | Verify `config set`, *Ende-zu-Ende* (Hash gleich, vier Aufrufformen, `java` Exit `0`, Reihenfolge `source` vor `version`), Mutationen M1 bis M5 | Zusage mit Messort |
| Schon gesperrte Datei bleibt abgelehnt (`generate`, `config set`, `doctor` Exit `11`); Abhilfe: `version`-Eintrag von Hand entfernen | Lead (steht schon im Eintrag `:17`) | Verify `config set`, Lesepfad-Probe (Exit `10`, `10`, `11`); Handbuch `:2072`–`:2075` | wie im Eintrag |
| Gemessen: Linux-Host, gebautes Binary, `docker run` mit Projektmount; Mutationstests der Verifier; Aufzählung getrennt je Fix (`config set`: `gh`, `glab`, `java`; `generate`: Eintrag aus nicht freigegebener Quelle, beide Profile, Auswege, `node`) | Lead (kurz) | beide Verify-Reports (Abschnitt *Ende-zu-Ende* und *Mutationen*); Plattform: nur „Linux“ ist im Report nicht ausdrücklich als `amd64` benannt — der Implementer liest die Messumgebung (Host-Architektur, Docker-Version) in den Reports nach; steht sie dort nicht, schreibt der Lead „one Linux host“ ohne `amd64` | Messort so genau, wie der Report ihn trägt |
| Ungemessen: Whitespace am Rand einer Quelle; `doctor` bei fehlendem `enabled`; `--json` der Ablehnung bei `config set`; macOS, Podman, andere Hosts | Lead (Grenzen) | Verify `generate`, *Ungemessen*; Verify `config set`, *Ungemessen* | „not measured“, nicht „nicht unterstützt“ |
| Lastenheft `0.22.0` unverändert | Lead | `spec/lastenheft.md` Zeile 9 (gelesen); `git diff v0.15.0..HEAD --stat -- spec` trifft nur `spezifikation.md` | ein Satz wie in `[0.15.0]` |
| CI: Docker-Integrationstests grün | Lead, Closure | erst nach M2 auf dem neuen Push-Tip | **Kein CI-Satz im Lead**, solange M2 aussteht; die Status-Blöcke behalten ihren bestehenden Satz |
| Kein Sicherheitsfix, keine Einstufung als Schwachstelle | — | Entscheidung des Projektinhabers steht aus (§1) | **keine Aussage** im Lead zur Sicherheit; weder „security“ noch „vulnerability“ |

### Aussage des Release — Soll für den CHANGELOG-Lead

Stil von `[0.13.1]` (englisch, `Patch release.` und fette Kernaussage). Tabellenzeile und Status-Block beider READMEs tragen eine Kurzform ohne Kennungen; sie widersprechen dem Lead nicht und sagen nichts, was er nicht sagt.

- **Kernaussage:** zwei Prüfungen, die `u-boot.yaml` betreffen, greifen jetzt dort, wo sie zugesagt waren: `u-boot generate devcontainer` weist die Quelle eines aktivierten Features ab, die nicht freigegeben ist; `u-boot config set …features.<feature>.version` weist `version` ohne `source` bei `ollama`, `gh` und `glab` ab.
- **Verhaltensänderung für Skripte:** `generate devcontainer` endet bei einer nicht freigegebenen Quelle mit `10` statt `0` und schreibt nichts; die zwei Auswege (Quelle freigeben, Feature abschalten) nennt der Eintrag. `config set` endet bei `version` ohne `source` mit `10` statt `0`.
- **Bezug zu `v0.15.0`:** der Lead von `[0.15.0]` nannte die Prüfung in `config set` als Folge; `v0.15.1` liefert sie. Eine Datei, die schon einen solchen `version`-Eintrag trägt, bleibt abgelehnt (Abhilfe von Hand).
- **Gemessen / nicht gemessen:** nach der Beleg-Tabelle (Messort, Mutationen; Whitespace am Rand, `doctor` bei fehlendem `enabled`, `--json` der Ablehnung von `config set`, andere Hosts ungemessen).
- **Unverändert:** Lastenheft `0.22.0`; kein neuer Befehl, keine neue Option. Keine Aussage zur CI vor M2; keine Einstufung als Sicherheitsfix.

Tabellenzeile der READMEs (Muster `v0.13.1`, `v0.14.0`): Name „Two checks that were missing“ / „Zwei fehlende Prüfungen“ (Patch), die zwei Verhaltensänderungen in Worten, die Auswege kurz, Link auf den GitHub-Release `v0.15.1`. Status-Block: ein Satz zum Stand in der Gegenwartsform, kein „released“, kein Datum.

### Reihenfolge der Maintainer-Aktionen

Der Implementer bereitet vor (Commits lokal, Gates lokal) und **führt M1 und M5 nicht aus**. **Push (M1) und Tag (M5) führt nur der Hauptagent, je mit eigener ausdrücklicher Freigabe des Projektinhabers im Gespräch** (nicht vorwegnehmen).

> **Rückfrage-Punkt (Hauptagent, nicht Planner und nicht Implementer):** Vor **M1** (Push des Cuts) und vor **M5** (Tag) fragt der Hauptagent den Projektinhaber und wartet auf die ausdrückliche Freigabe im Gespräch. Eine Freigabe für M1 ist keine für M5. Der Slice bleibt bis zum Tag in `in-progress/`.

| # | Aktion | Wer | Prüfung / Erwartung |
|---|---|---|---|
| M0 | `date +%F` (Ortszeit); Datum in `CHANGELOG.md`, Tabellenzeile beider READMEs und Handbuch-`Stand` gegen das Tag-Datum halten | Implementer vor dem Handoff; Hauptagent vor M1 und vor M5 | Weicht es ab: ein weiterer Cut-Commit vor dem Tag und M2 auf dem neuen Push-Tip erneut ([`releasing.md`](../../../maintainer/releasing.md) §5); Cut, Push und Tag am selben Tag |
| M-T | Tap-Trockenlauf (`releasing.md` §4), im Docker, `--network none`, ohne Schreiben ins Tap; Muster Cut `v0.15.0` §3 | Implementer; Ergebnis in §3; der Verifier fährt ihn nach | `tap-formula-check.sh` Exit `0` mit `v0.15.1` und Exit `1` mit `v0.15.0`; `tap-nachzug.sh` mit Stellvertretern protokolliert mit `v0.15.1` `clone`, `add`, `commit -m u-boot 0.15.1` und `push` und bricht mit `v0.15.0` vor `clone` ab. **Beweist nicht:** `is_stable`, echter `curl`, echter `git push`, Token |
| M1 | `git push origin main` (alle Commits des Cuts in einem Push) | Hauptagent, **nach Freigabe des Projektinhabers** | Review- und Verification-Report liegen vor; `make docs-check` lokal grün; Arbeitsbaum sauber |
| M2 | `gh run list --limit 4`: `ci` und `integration-docker` auf dem Push-Tip; Logs über `gh api repos/pt9912/u-boot/actions/jobs/<id>/logs`, Warnungen über `gh api repos/pt9912/u-boot/check-runs/<id>/annotations` ([`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)) | Hauptagent | `ci` `success`: Jobs `gates`, `security-gates`, `image-scan`; keine Annotation „Unexpected input(s)“. `integration-docker` `success` (Konklusion **und** Zählung): **4 `--- SKIP`** (die vier bekannten Nested-Podman-Tests `TestE2E_SandboxDevcontainer_NestedPodman_Build`, `…_NestedPodman_BuildEgress`, `…_NestedPodman`, `…_StatusFiles_VFS`), **0 `--- FAIL`**, kein `panic: test timed out`; **Dauer** des Jobs gegen `timeout-minutes: 30` festhalten; Abschnitt „Operating System“ der Logs lesen. Abweichend: Release-Aussage vor dem Tag nachziehen |
| M3 | Vorab-Tag `v0.15.1-rc.1`: entfällt | — | Vorschlag des Planners; bestätigt der Projektinhaber beim Start (§7) |
| M4 | `gh workflow run tap-check` | Hauptagent, unmittelbar vor M5 | grün = das Token darf auf das Tap schreiben (`releasing.md` §2, §4) |
| M5 | `git switch main && git pull --ff-only`, `git tag v0.15.1` auf dem in M2 geprüften Commit, `git push origin v0.15.1` (nicht `--tags`) | **Hauptagent, nur nach ausdrücklicher Freigabe des Projektinhabers im Gespräch** | `HEAD` = der Commit aus M2; nach dem Push wird der Tag nicht verschoben (`releasing.md` §8); `v0.15.0` bleibt unberührt |
| M5b | Security Advisory: **entfällt** | — | Einstufung *Fixed* (§1); will der Projektinhaber *Security*, ändert sich dieser Plan vor M1 |
| M6 | `releasing.md` §6: `gh run watch`; `gh release view v0.15.1 --json assets,isPrerelease,isDraft`; `docker run --rm ghcr.io/pt9912/u-boot:0.15.1 --version`; Digest `:0.15.1` = `:latest`; Log und Annotationen des Jobs `tap`; Tap-Commit und Smoke-Workflow im Tap | Hauptagent, unmittelbar nach dem Push des Tags | DoD Punkt 2, Kontrollzeilen |
| M7 | Auf einem Mac: `brew tap pt9912/u-boot`, `brew trust pt9912/u-boot`, `brew install u-boot` (oder `brew upgrade u-boot`), `u-boot --version` = `0.15.1` | Projektinhaber | `releasing.md` §6; **offen**, nicht vom Cut entschieden |
| M8 | Nacharbeit (`releasing.md` §7): Roadmap-Zeile mit Tag-Hash und Datum, Ruhe-Marker; Datum gegen `git log -1 --format=%cs v0.15.1`; Zeiger „zuletzt“ in `releasing.md`; Closure-Notiz, Register, `git mv` nach `done/` | Planner (Closure), Implementer (Roadmap) | Closure erst nach M6 |

### Sensoren (vom Implementer vor dem Handoff, Exit-Codes der `make`-Prozesse in §3 *Umsetzung*)

`make gates`, `make ci`, `make fullbuild`, `make test-docker` (mit `-v`: Zahl von PASS, SKIP und FAIL aus dem Log; ein lokaler Lauf ersetzt den CI-Lauf nicht; am Paket-Zeitlimit von 30 Minuten endete der lokale Lauf in mehreren Cuts mit Exit 2, [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)), `make build VERSION=0.15.1` mit `docker run --rm u-boot:latest --version` = `u-boot version 0.15.1` (lokal, nicht gepusht), `make doc-immutable RANGE="v0.15.0..HEAD"`, `make doc-trace`, `make doc-commits RANGE="origin/main..HEAD"`, `make doc-planning`; M-T. `make verify-depguard` entfällt (keine Importregel berührt), mit Begründung zu nennen.

### Umsetzung (Implementer)

**Stand (Implementer, 2026-10-24; vor M1: nichts gepusht, kein Tag).** Commits lokal auf dem Plan-Commit `6335bf18`: `f246e155` reiner Move nach `next/` · `b1d94f5b` Verweis-Nachzug · `f053b1ec` reiner Move nach `in-progress/` · `617cd0d8` Verweis-Nachzug · `e29f2fa5` Start (Verantwortlich, Entscheidungen in §7, Roadmap-Zeiger, Zustandsmarker) · `521ab770` **Cut-Commit**. Der Tag gehört auf den letzten Commit des Cuts nach allen Fix-Runden (M5). Keine Commit-Message trägt eine `SPEC-`/`ARC-`Kennung.

**M0:** `date +%F` = `2026-10-24` (Ortszeit). Das Datum steht in `CHANGELOG.md` (`## [0.15.1] - 2026-10-24`), in der Tabellenzeile beider READMEs und im Handbuch-`Stand` (vorher `2026-10-23`). Weicht `date +%F` vor M1 oder M5 ab, wird neu geschnitten (`releasing.md` §5).

**Ist gegen Plan, Abweichungen und Entscheidungen im Zweifel:**

1. **Versionsstellen** wie im Plan (Suche neu gefahren vor der Änderung: gleiche Fundstellen): `cmd/uboot/main.go:43`, `repo.mk:51`, `:54`, `:55`, `Dockerfile:35`, `:37`, `:39`, `README.md` Status-Block und `VERSION=`, `README.de.md` ebenso, Handbuch `:4`, `:5`, `:82`, `:108`, `:135`, `:2998`. Nach dem Cut trifft `grep -rIn -e "0\.15\.0"` außerhalb der Ausschlüsse nur noch: die beiden Tabellenzeilen `v0.15.0` (bleiben, Geschichte), `docs/maintainer/releasing.md:57` (Zeiger „zuletzt“, Closure), `spec/lastenheft.md:1591` und `spec/spezifikation.md:1435` (Spec-Versionen). `0.15.0-dev` kommt nirgends mehr vor.
2. **Plattform im Lead: „one Linux host“, ohne `amd64`.** Die beiden Verify-Reports nennen weder Architektur noch Docker-Version der Messung (gelesen: Abschnitte *Ende-zu-Ende*, *Ungemessen*); `docker info` des Hosts, auf dem dieser Cut läuft, meldet `linux x86_64 29.9.2`, belegt aber nicht den Host der Verifier-Läufe. Der Slice `generate` nennt „Linux `amd64`“ für die Messung des Implementers (§7 jenes Slice), der Slice `config set` nennt keine Architektur; der Lead trägt deshalb die schwächere Form für beide. Die READMEs ebenso.
3. **Lead nach dem Soll aus §3:** Kernaussage, Verhaltensänderung für Skripte (zwei Auswege im Eintrag `Fixed`), Bezug zu `[0.15.0]`, Gemessen/Ungemessen (Whitespace am Rand einer Quelle, `doctor` bei fehlendem `enabled`, `--json` der Ablehnung bei `config set`, macOS, Podman, andere Hosts), „No new command and no new option. Lastenheft 0.22.0.“ Kein CI-Satz, keine Sicherheitsaussage (weder „security“ noch „vulnerability“; geprüft per Suche im Block `[0.15.1]`). Die Einträge unter `### Fixed` blieben wörtlich. Der Text unter `[0.15.0]` ist unberührt (`git diff` zeigt dort nur Kontext).
4. **READMEs:** Der Status-Block trägt einen Satz in Gegenwartsform (ohne Datum, ohne „released“, ohne Kennungen) hinter dem Katalogsatz zu `gh` und `glab`; die Tabellenzeile `v0.15.1` steht vor `v0.5.0` (nach `v0.15.0`, Reihenfolge der Tabelle). Zeilen `v0.15.0` und der bestehende CI-Satz im Status-Block unverändert.
5. **Handbuch:** nur Kopf und Versionsbeispiele; die Absätze zu beiden Fixes (`:1118`–`:1127`, `:1965`, `:2062`–`:2075`, `:2417`) wurden nicht angefasst (Plan: gelesen, tragen die Fixes schon).
6. **Rot gesehen, Grenze:** Kein Test pinnt die Nummer (`cmd/uboot/main_test.go`: `--version` nur „nicht leer“, `package-level version must not be empty`); eine falsche Versionszeile färbt kein Gate. Wächter sind die Prüfung `--version` = Tag im Workflow `publish` (nach dem Push) und die lokale Gegenprobe: `make build VERSION=0.15.1` meldet `0.15.1` (Label `0.15.1`), der Standardbau (`make fullbuild`) `0.15.1-dev` (Label `0.15.1-dev`). Die Prüfung auf den Tag in M-T ist rot gesehen (2b, 3b).
7. **`make ci`: ein Netzfehler, dann grün.** Der erste Lauf endete mit Exit 2 im Schritt `image-scan`: Trivy konnte die Schwachstellen-Datenbank nicht laden (`mirror.gcr.io/aquasec/trivy-db:2`, `404 Not Found`; auch zwei Wiederholungen von `make image-scan` im Abstand von einer Minute). Zur Diagnose lief derselbe Trivy-Aufruf mit `--db-repository ghcr.io/aquasecurity/trivy-db:2` (kein Gate, kein Ersatz): Exit 0, 0 Treffer. Der spätere Lauf von `make fullbuild` und danach `make ci` traf den Mirror wieder und lief grün. Keine Änderung an Makefile, Pins oder Gates; Ursache: Umgebung (Mirror), nicht der Cut.
8. **`make test-docker` nicht gelaufen**, Entscheidung des Auftrags: Paket-Zeitlimit (`BEO-E2E/e2e-paket-am-zeitlimit`) und Speicher des Hosts; der Diff berührt keinen Produktcode (`git diff v0.15.0..HEAD --stat -- cmd/uboot/main.go`: eine Zeile, `var version`). Die Docker-Tests laufen in M2 im CI; M2 zählt dort 4 `--- SKIP`, 0 `--- FAIL`.
9. **`make verify-depguard` nicht gelaufen:** keine Importregel und kein Go-Import berührt.
10. **Fix-Runde LOW-V1 (Messaussage getrennt).** Die Aufzählung `gh`, `glab`, `java` stand in `f16b89c2` für beide Fixes; belegt ist sie nur für `config set` (Verify `config set`, *Ende-zu-Ende*). Ursache: ein fehlerhafter Auftrag des Projektleiters an die vorige Fix-Runde, der den Umfang des `config set`-Verifys auf beide Fixes bezog; der `generate`-Report trägt diese Aufzählung nicht. Jetzt gelesen (Verify `generate`, *Ende-zu-Ende*: Eintrag `evil`, Standard- und Sandbox-Profil, Abhilfen, `node` ohne `source`) und gleichlautend an den vier Stellen getrennt: `CHANGELOG.md` `[0.15.1]` (*Measured*), Status-Block und Tabellenzeile `v0.15.1` in `README.md` und `README.de.md`; „one Linux host“ bleibt ohne `amd64`. Keine Sicherheits- und keine CI-Aussage, keine Kennung in den READMEs.

**Sensoren** (Docker, Host `linux/amd64`, Docker Engine 29.9.2; Exit-Codes der `make`-Prozesse; auf `521ab770`):

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0: `lint` 0 issues; `coverage-gate` 92,10 % gegen 90 %; `docs-check` 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung` 26 Fälle, 0 Fehlschläge |
| `make ci` | erster Lauf Exit 2 (Trivy-DB, Punkt 7); Wiederholung nach `make fullbuild` Exit 0: `govulncheck` „No vulnerabilities found“, `image-scan` `u-boot:latest (debian 12.15)` 0, `usr/local/bin/u-boot` 0 |
| `make fullbuild` | Exit 0; `docker run --rm u-boot:latest --version` = `u-boot version 0.15.1-dev`, OCI-Label `0.15.1-dev` |
| `make build VERSION=0.15.1` | Exit 0; `--version` = `u-boot version 0.15.1`, Label `0.15.1` (lokal, nicht gepusht) |
| `make doc-immutable RANGE="v0.15.0..HEAD"` | Exit 0, 802 Dateien, 0 Befunde |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0, 0 Befunde |
| `make doc-planning` | Exit 0, 0 Befunde |
| `make test-docker` | nicht gelaufen (Punkt 8) |
| `make verify-depguard` | nicht gelaufen (Punkt 9) |

**M-T ausgeführt (Trockenlauf des Tap-Pfads).** Befehl: `docker run -i --rm --network none -v /Development/u-boot:/src:ro u-boot:test bash < <Skript>` (Image aus `make gates`, Arbeitsbaum schreibgeschützt; Skript im Scratchpad, nicht im Repo; Fixture nach `sumsFull` in `cmd/uboot/homebrew_scripts_test.go`; Stellvertreter für `curl` und `git` im `PATH`). Ausgabe:

```text
(1)  homebrew-formula-fill.sh v0.15.1 <sums> scripts/homebrew-formula.rb.tmpl <Ausgabe>   exit=0   (vier url-Zeilen mit v0.15.1)
(2a) tap-formula-check.sh v0.15.1 <Ausgabe>   exit=0
(2b) tap-formula-check.sh v0.15.0 <Ausgabe>   exit=1
     tap-nachzug: the u-boot.rb asset of v0.15.0 does not carry four download URLs of v0.15.0 (found 4 url lines, 0 of v0.15.0).
(3)  tap-nachzug.sh, TAP_TOKEN=dummy-token, TAG=v0.15.1:   exit=0
     Protokoll: curl … /releases/download/v0.15.1/u-boot.rb -o <tmp> · git clone --depth 1 … · git diff --quiet -- Formula/u-boot.rb · git … add Formula/u-boot.rb · git … commit -m u-boot 0.15.1 · git push origin HEAD
     tap-nachzug: pt9912/homebrew-u-boot updated to v0.15.1.
(3b) tap-nachzug.sh mit TAG=v0.15.0 gegen dieselbe Formel: exit=1, Protokoll nur der curl-Aufruf, 0 clone-Aufrufe, 0 push-Aufrufe
```

Rot gesehen: (2b) und (3b) sind die rote Gegenprobe der Prüfung auf den Tag. **Grenze, wie in M-T benannt:** `is_stable == 'true'`, der Workflow-Schritt, der echte `curl` gegen das Release-Asset, der echte `git push` und das Token sind nicht geprobt; der Stellvertreter von `git` meldet immer eine Änderung. Befund: keiner; kein Skript geändert. `tap-check` (M4) ist Sache des Hauptagenten.

**Plan-Tabelle gegen `git diff --stat 6335bf18..HEAD`** (zehn Dateien): `CHANGELOG.md`, `cmd/uboot/main.go`, `repo.mk`, `Dockerfile`, `README.md`, `README.de.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md` (Start-Zeiger; Zeile `v0.15.1` unter *Nächste Wellen*: „Plan in `open/`“ → „in Arbeit“), `docs/plan/planning/open/README.md` (Zeile dieses Slice: „in Arbeit“), dieser Slice (Moves, Verantwortlich, §3, §7). Jeder Eintrag der Tabelle mit Diff steht dort; **kein Diff, gelesen:** `spec/*` (Lastenheft Zeile 9 `0.22.0` unverändert), `releasing.md` (Zeiger „zuletzt“ `:56`–`:57` nennt `v0.15.0`, Closure), `.github/`, `scripts/` (`tap-nachzug.sh`, `tap-formula-check.sh`, `homebrew-formula-fill.sh` gelesen), `packaging/`, `Makefile`, `go.mod`, `tools/`, `d-check.mk`, `Dockerfile` außer Versionsblock (Pins im Gate-Log: Go 1.27.3, golangci-lint v2.14.0, Trivy 0.70.0 unverändert), `internal/`. Die Zeilen der Reports und der Beobachtungen entstehen bei Review, Verifikation und Closure.

## 4. Trigger

**Start** (`next` → `in-progress`; der Slice liegt in `open/`, der Projektinhaber hat den Patch am 2026-10-24 entschieden): Freigabe zum Start durch den Projektinhaber; die beiden Fix-Slices liegen in `done/` (eingetreten, gelesen). Der Slice geht per `make slice-mv` über `next/` nach `in-progress/`; der Roadmap-Zeiger folgt im eigenen Commit.

**Rückführungen — vorab benennen:**

- `in-progress` → `open` (blockiert): ein roter `image-scan` durch einen neuen Advisory (Go-Pin als eigener Vorgang davor), ein roter CI-Lauf auf dem Push-Tip, der nicht am Cut liegt, oder die Entscheidung des Projektinhabers für *Security* (dann Plan ändern, M5b aufnehmen).
- `in-progress` → `next` (zu groß): trifft nicht zu (reiner Textcut).

## 5. Closure-Trigger

DoD vollständig, der Tag `v0.15.1` gesetzt und nach §6 der `releasing.md` kontrolliert (M6), Closure-Notiz mit Lerneintrag geschrieben (§7). Der `→ done`-Übergang ist ein reiner `git mv` nach dem Inhalt ([`AGENTS.md`](../../../../AGENTS.md) §3.3).

## 6. Risiken und offene Punkte

- **R1 — Der Lead sagt mehr, als die Reports tragen** (Plattform-Aussage `amd64`, „auch mit `--dry-run`“ für beide Fixes, Vorher-Verhalten am Release-Image). Gegenmaßnahme: Beleg-Tabelle, Formulierung nur mit gelesener Fundstelle ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)). — **Ausgang:** offen bis zur Closure.
- **R2 — Einstufung als Security ändert sich nach der Planung** (Projektinhaber). Gegenmaßnahme: der Plan nennt sie als Entscheidung des Projektinhabers und M5b als bedingten Schritt; ein Advisory ist unumkehrbar, deshalb vor M1 klären ([`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md)). — **Ausgang:** offen bis zur Closure.
- **R3 — Plan-Aussage über Code ohne Lesung** (Versionsstellen, Zeilen der Tabelle verschoben, „unverändert“ per Suche statt Lesung; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md)). Gegenmaßnahme: Abgleich gegen `git diff --stat` nach der letzten Fix-Runde (§3). — **Ausgang:** offen bis zur Closure.
- **R4 — CI auf dem Push-Tip rot oder am Zeitlimit** ([`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md)); lokales `make test-docker` kann am Paket-Zeitlimit mit Exit 2 enden. Gegenmaßnahme: M2 mit Dauer; lokaler Lauf nicht als CI-Ersatz. — **Ausgang:** offen bis zur Closure.
- **R5 — Datum driftet** (Cut und Tag an verschiedenen Tagen). Gegenmaßnahme: M0 vor M1 und M5. — **Ausgang:** offen bis zur Closure.
- **R6 — Projekte mit nicht freigegebener Quelle brechen nach dem Update** (Verhaltensänderung für Skripte). Gegenmaßnahme: Lead nennt es als solche mit den zwei Auswegen; ein Patch-Release trägt sie, weil die Zusage schon galt ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

**Entscheidungen beim Start** (2026-10-24, Implementer nach dem Auftrag des Hauptagenten, der sie als vom Projektleiter bestätigt benennt; ein Wortlaut des Projektinhabers liegt nicht vor): (1) Nummer `v0.15.1`; (2) Einstufung *Fixed*, kein *Security*, kein Advisory, M5b entfällt; (3) kein Vorab-Tag, M3 entfällt; (4) kein CI-Satz im Lead vor M2; (5) **Reichweite:** Push (M1) und Tag (M5) führt nur der Hauptagent, je Vorgang nach ausdrücklicher Freigabe des Projektinhabers; der Implementer pusht und taggt nie, auch nicht lokal; das volle `make test-docker` läuft lokal nicht (Paket-Zeitlimit, Speicher), die Docker-Tests laufen in M2 im CI.

**Ursache LOW-V1** (Fix-Runde nach `9db50093`): Die zu breite Messaussage in `f16b89c2` (Binary-Proben mit `gh`, `glab`, `java` für beide Fixes) stammte aus einem fehlerhaften Auftrag des Projektleiters, der den Umfang des `config set`-Verifys auf beide Fixes bezog; sie ist nicht dem Implementer allein zuzurechnen. Behoben durch Trennung je Fix (§3 *Umsetzung* Punkt 10).

*Wird bei der Closure (Planner, nach M6) gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen · Läufe und Abfragen von M0 bis M7 (Hauptagent).*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `cmd/uboot` (Wiring, eine Versionszeile; keine Sub-Area laut Modus-Deklaration, gelesen in `harness/conventions.md`) und `spec/`, `harness/`, `docs/plan/`, `docs/messungen/` (DOK, Greenfield; für Release-Text, Roadmap, Handbuch, README). Kein Produktcode in `internal/`, keine neue Sub-Area, keine zu grobe Fläche.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (`docs/plan/planning/observations/`, Stand 2026-10-24), Treffer je berührter Sub-Area (DOK) und Cut-Klasse:

- [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) — verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10, 50 Evidence-Dateien, Rückfälle zuletzt auch im Cut `v0.15.0`; hier relevant für Versionsstellen und „unverändert“ (R3).
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) — verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10, 24 Evidence-Dateien; hier relevant, weil der Lead Bedingungen („auch mit `--dry-run`“) nennt, die ein Test pinnt (Verify-Mutationen).
- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — offen, 20 Belege, über der Schwelle, Ausgang beim Projektinhaber; relevant für den Lead (R1).
- [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) — offen, 1×; relevant für M2.
- [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md) — offen, 1×; relevant für M1/M5 und die Einstufung (R2).
- [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (offen, 3×, Schwelle erreicht), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) und [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) — Lese-Schritte des Cuts (M2, lokales `make test-docker`).
- [`BEO-APP/pruefung-nur-beim-lesen`](../observations/BEO-APP/pruefung-nur-beim-lesen/observation.md) — offen, zwei Vorgänge (die zwei Inhalte dieses Release); der Cut ändert keinen Code und zählt keinen dritten Vorgang, die Closure liest nur.

Evidenz-/Diskrepanz-Risiko: R1 bis R3.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
