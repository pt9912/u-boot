# Slice v1: `u-boot --help` verweist auf das Handbuch — Link auf den Tag und Link auf `main`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — kein Release-Bezug festgelegt; der Slice läuft **nach** dem Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) und nach den beiden Hygiene-Slices [`slice-v1-cli-dry-run-log-would`](../done/slice-v1-cli-dry-run-log-would.md) und [`slice-v1-harness-gofmt-und-depguard`](../done/slice-v1-harness-gofmt-und-depguard.md); er ist nicht Teil von `v0.11.0`. Herkunft: Wunsch des Projektinhabers 2026-10-18 („bei u-boot --help wird das readme referenziert. Können wir noch unser Handbuch in die Meldung mitaufnehmen“).

**Bezug:** [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) (Hilfeausgabe),
[`LH-FA-CLI-003`](../../../../spec/lastenheft.md#lh-fa-cli-003--versionsausgabe) (die Version, aus der der Tag-Link abgeleitet wird),
[`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache) (CLI-Ausgaben Englisch).

**Berührte Spec-Stellen:** keine. Gelesen: `spec/spezifikation.md` enthält keinen Wortlaut des Hilfetexts (Suche nach `readme` und `--help`: nur der Generator `readme` und der Beispielaufruf in [`SPEC-024`](../../../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen)); [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) verlangt Befehle, Kurzbeschreibung, Optionen und Beispiele, nicht die Dokumentations-Links. Das Lastenheft bleibt unverändert.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot --help` nennt neben dem README zwei Handbuch-Links: einen auf das Handbuch der Version des Binarys (Tag `v<version>`) und einen auf das aktuelle Handbuch (`main`); der Tag-Link wird aus der Version abgeleitet, pro Release-Cut entsteht kein Handarbeitsschritt.

### Entscheidung des Projektinhabers (2026-10-18)

Auf die Rückfrage „Link auf Tag oder `main`?“: „können wir beides machen. eins auf den tag und eins auf das aktuelle handbuch“. Der Text ist Englisch ([`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache)); das Handbuch ist Deutsch, der Text sagt es („(German)“).

### Gefundener Stand im Code (Ist-Stand, gelesen am 2026-10-18)

- **Wo der Text lebt.** `internal/adapter/driving/cli/root.go`, Funktion `buildRootCommand`: das Feld `Long` des Wurzelkommandos ist ein Raw-String-Literal und endet mit `Quickstart and user docs: https://github.com/pt9912/u-boot#readme`; `Version: a.version` steht im selben Literal darunter. Die Funktion bekommt `*App`; `App.version` (`cli.go`, Feld `version`, Kommentar „surfaced via `u-boot --version`. The wiring layer passes it in; the CLI package does not own version metadata“) setzt `New(version, …)`. `App.Execute` baut das Wurzelkommando bei **jedem** Aufruf neu (`cmd := buildRootCommand(a)`), die Version liegt beim Bau des `Long`-Textes vor. Es braucht keinen Port: Die Version ist bereits im Adapter, der Text ist reine Darstellung.
- **Wie die Version in den Code kommt.** `cmd/uboot/main.go`: `var version = "0.11.0-dev"` (Kommentar: überschreibbar per `-ldflags="-X main.version=…"`); `run` reicht sie an `cli.New(version, …)` weiter. `Dockerfile`: `ARG UBOOT_VERSION=0.11.0-dev`, die Stage `build` übergibt `-ldflags="-s -w -X main.version=${UBOOT_VERSION}"`. `repo.mk`: `--build-arg UBOOT_VERSION=$(VERSION)` und `-X main.version=$(VERSION)` für die Binaries. `.github/workflows/publish.yml`: `VERSION="${TAG_REF#v}"` — der Tag-Build setzt die Version **ohne** führendes `v`; ein Tag mit `-` (Vorabversion) ist „nicht stabil“, trägt aber seinen Tag (`v<version>`).
- **Was `--version` liefert.** Das Wurzelkommando setzt `Version: a.version` und keine eigene Versionsvorlage (Suche nach `SetVersionTemplate` in `internal`: ohne Treffer); Cobra gibt also seine Standardvorlage aus. `publish.yml` wertet die Ausgabe als `awk '{print $NF}'` und vergleicht sie mit `VERSION` (Prüfschritt „Pin u-boot --version output against the tag-derived VERSION“): das letzte Wort der Ausgabe ist die Version ohne `v`. **Release-Build:** `0.11.0` (kein `v`, kein Suffix). **Dev-Build** (lokales `make build`, Quellstand): `0.11.0-dev`. Der Tag `v0.11.0-dev` existiert nie; `docs/maintainer/releasing.md` §3 Schritt 2 setzt die `-dev`-Version nach jedem Cut auf `<neu>-dev`.
- **Tests zum Help-Text.** `cli_test.go`: `TestExecute_RootHelpPointsToUserDocs` prüft `u-boot <command> --help` und `https://github.com/pt9912/u-boot#readme`, verbietet die Strings `See spec/lastenheft.md …` und `LH-` sowie die Katalog-Begriffe `PostgreSQL`, `Keycloak`, `OpenTelemetry`; `TestExecute_HelpListsInit` prüft `init`. Alle Testfabriken (`newApp` u. a.) bauen den `App` mit der Version `"0.0.0-test"`.
- **Handbuch und Doku.** `docs/user/benutzerhandbuch.md` beschreibt `-v, --version` und `-h, --help` in der Optionstabelle des Anhangs (je eine Zeile) und im FAQ „Wo finde ich alle Optionen eines Befehls?“ (`u-boot <befehl> --help` …); der Hilfetext des Wurzelkommandos ist dort nicht beschrieben. Der Kopf des Handbuchs trägt `Software-Version: v0.11.0`. `docs/maintainer/releasing.md` §3 (Cut-Schritte) nennt die Versionsstrings an genau drei Stellen (`cmd/uboot/main.go`, `repo.mk`, `Dockerfile`); ein Help-Text steht dort nicht.
- **Repo-URL.** Die URL `https://github.com/pt9912/u-boot` steht im `Long`-Text hartkodiert; Modulpfad `github.com/pt9912/u-boot` (`go.mod`).

### Entwurf — die Fragen, die der Slice klärt, mit Antwort

1. **Wie wird der Tag-Link aus der Version abgeleitet?** `https://github.com/pt9912/u-boot/blob/v<version>/docs/user/benutzerhandbuch.md` mit `<version>` = `a.version`, das `v` wird vorangestellt (die Version trägt keines, s. o.).
2. **Dev-Stände (`…-dev`)?** Der Tag existiert nicht; ein Link wäre tot (404). Entscheidung: für eine Version mit Suffix `-dev` (und für eine leere Version) zeigt die Hilfe **nur** den `main`-Link; das Abschneiden von `-dev` wird verworfen, weil es auf den Tag des **letzten** Release (oder keinen) zeigte und Handbuch und Binary nicht mehr zusammenpassten, ohne dass der Benutzer es merkt. Festlegung des Leads (2026-10-18, strenger als der erste Entwurf): der Tag-Link erscheint **nur**, wenn die Version genau einer Release-Form entspricht, `X.Y.Z` oder `X.Y.Z-rc.N` (Zahlen ohne führende Null; `rc.N` ist die in `docs/maintainer/releasing.md` §1 und §5 genannte Vorab-Schreibweise `vX.Y.Z-rc.N`, `publish.yml` Schritt „Validate SemVer tag“ lässt weitere Vorab-Formen zu); jede andere Version (`-dev`, `0.0.0-test`, `0.12.0-rc1`, leer, selbst gebaut) bekommt nur den `main`-Link. Umgesetzt in `root.go`, `releaseVersionPattern` und `rootLongText`.
3. **Wo lebt der Text?** In `root.go`: `Long` wird aus einer kleinen reinen Funktion (Arbeitsname `rootLongText(version string) string`) zusammengesetzt, die ein Format-String-Muster mit `fmt.Sprintf` füllt; die Repo-URL steht als Konstante daneben. Kein Template-Paket, kein Port, kein neuer Typ.
4. **Wortlaut (Entwurf, Englisch).**

   ```text
   Quickstart and user docs: https://github.com/pt9912/u-boot#readme
   User manual (German), this version: https://github.com/pt9912/u-boot/blob/v0.11.0/docs/user/benutzerhandbuch.md
   User manual (German), latest: https://github.com/pt9912/u-boot/blob/main/docs/user/benutzerhandbuch.md
   ```

   Bei jeder Version außerhalb der Release-Formen (etwa `…-dev`) entfällt die mittlere Zeile; die letzte Zeile heißt dann ebenfalls „latest“. Der Wortlaut darf kein `LH-` enthalten (bestehender Test).
5. **Handarbeit pro Cut.** Keine: Der Tag-Link folgt `version`, die der Cut ohnehin an den drei Stellen setzt; `docs/maintainer/releasing.md` und die Vorlage der Cut-Slices bleiben unverändert (Beleg: §3 Schritt 2 setzt nur die Versionsstrings, der Link wird zur Laufzeit abgeleitet).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Link auf das Handbuch in den Hilfetexten der Unterbefehle:** Das Wurzelkommando ist der Einstieg, den der Wunsch nennt; je Befehl ein Link wäre eine eigene Zusage mit Anker-Pflege.
- **Eine Prüfung, dass der Tag oder die Seite im Netz existiert:** `--help` bleibt netzlos und unter 200 ms ([`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort): `u-boot --help` und `u-boot --version` unter 200 ms); die Ableitung aus der Version ist die Garantie.
- **Ein Konfigurationsschalter für die Repo-URL (Fork, Umbenennung):** Die URL stand schon im README-Link hartkodiert; ein Schalter wäre eine Zusage ohne Nachfrage. Das Risiko steht in §6.
- **Eine Änderung der Versionsausgabe (`--version`) oder der Versionsstrings:** unberührt.
- **Ein Gate, das die Links im Netz prüft:** Der Test pinnt den Text (§2); ein Link-Checker gegen GitHub wäre ein neues, netzgebundenes Gate mit eigener Begründung.

## 2. Definition of Done

- [x] **Code:** `u-boot --help` enthält neben dem README-Link den Tag-Link und den `main`-Link (Wortlaut §1, Entwurf 4), wenn die Version eine Release-Form ist (`X.Y.Z` oder `X.Y.Z-rc.N`); bei jeder anderen Version (`-dev`, `0.0.0-test`, `0.12.0-rc1`, leer) enthält sie nur den `main`-Link und **keinen** Tag-Link ([`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe), [`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache)).
- [x] **Test:** ein Test baut den `App` mit `0.11.0` und prüft beide Links vollständig (Tag `v0.11.0`, `main`, Pfad `docs/user/benutzerhandbuch.md`); ein Fall mit `0.11.0-dev` (und `0.0.0-test`, `0.12.0-rc1`, leer u. a.) prüft den `main`-Link und die Abwesenheit von `/blob/v`; ein Fall mit der Vorabversion `0.12.0-rc.1` prüft den Tag-Link `v0.12.0-rc.1`. Mutation: Der Tag-Link ohne `v`, die Regex gelockert und die Verzweigung entfernt — jede Mutation lässt mindestens einen Test rot werden, am Beispiel gesehen und zurückgenommen. Der bestehende Test `TestExecute_RootHelpPointsToUserDocs` bleibt grün (README-Link, kein `LH-`, keine Katalog-Begriffe).
- [x] Doku-Update ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 5): `CHANGELOG.md` `[Unreleased]` nennt die Änderung (Added: `u-boot --help` zeigt den Handbuch-Link der Version und den des aktuellen Handbuchs; Dev-Builds nur den aktuellen); das Handbuch erhält einen Satz im FAQ „Wo finde ich alle Optionen eines Befehls?“ (die Hilfe verweist auf das Handbuch Ihrer Version), oder §7 begründet, warum nicht.
- [x] `docs/maintainer/releasing.md` und die Cut-Vorlage unverändert; §7 hält fest, dass der Cut keinen neuen Handarbeitsschritt trägt (der Cut-Slice des nächsten Release ist der Beleg: kein Link wird von Hand gesetzt).
- [x] `make gates` grün; `make doc-trace` ohne neue Waise.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Jede Aussage über Code steht mit Fundstelle (Datei und Symbol, gelesen am 2026-10-18). Der Implementer gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10) und belegt jede Datei, die er als „unverändert“ führt, durch eine gelesene Aussage.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | `buildRootCommand`: `Long` ruft `rootLongText(a.version)`; neu: Konstanten `repoURL`, `manualPath`, `releaseVersionPattern` und die reine Funktion `rootLongText`, die den Tag-Link nur bei einer Release-Version (`X.Y.Z`, `X.Y.Z-rc.N`) setzt |
| `internal/adapter/driving/cli/cli_test.go` | update | neu `helpForVersion` und `TestExecute_RootHelpManualLinks` (Tabelle über Versionen, baut den `App` mit `cli.New(version, …)`) vor `TestExecute_UnknownCommand`; `TestExecute_RootHelpPointsToUserDocs` unverändert (gelesen: README-Link, kein `LH-`, keine Katalog-Begriffe) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Ruhe-Marker „Nichts in Arbeit.“ in „Offene Wellen“ durch die Arbeitszeile mit Slice-Link ersetzt (d-check `planning-drift`) |
| `internal/adapter/driving/cli/export_test.go` | prüfen | gelesen: `WalkRootCommandPathsForTest` ruft `buildRootCommand(a)` und liest nur die Pfade; ändert sich die Signatur von `buildRootCommand` nicht, bleibt die Datei unberührt |
| `CHANGELOG.md` | update | `## [Unreleased]` (derzeit nur `### Changed` des Dry-Run-Slice) bekommt `### Added`; Englisch, wie der Bestand |
| `docs/user/benutzerhandbuch.md` | update, ein Satz | FAQ „Wo finde ich alle Optionen eines Befehls?“ (Zeile mit `u-boot <befehl> --help`); die Optionstabelle im Anhang (`-h`, `--help`: „Hilfe zum Befehl“) bleibt, sofern der Satz im FAQ genügt |
| `README.md`, `README.de.md` | prüfen | gelesen: die Suche nach `help` in `README.md` findet nur `make help`; kein Text über den Hilfetext, keine Änderung geplant |
| `spec/spezifikation.md`, `spec/lastenheft.md` | unverändert | gelesen (siehe Kopf): kein Wortlaut des Hilfetexts; [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) nennt keine Dokumentations-Links |
| `docs/maintainer/releasing.md`, Cut-Vorlage | unverändert | gelesen: §3 Schritt 2 setzt die drei Versionsstrings; der Link wird aus `version` abgeleitet, es entsteht kein Handarbeitsschritt |
| `docs/reviews/<Datum>-slice-v1-cli-help-handbuch-link.md` | neu, durch den Reviewer | Review-Report |

**Hinweise aus dem Review** (`docs/reviews/2026-10-18-slice-v1-cli-help-handbuch-link.md`):

- **Gelöschte Vorab-Tags.** `docs/maintainer/releasing.md` §5 und §7 löschen Vorab-Tags nach dem Cut; ein danach noch laufendes `X.Y.Z-rc.N`-Binary zeigt einen Tag-Link, der ins Leere führt (die Tags `v0.10.0-rc.1` und `v0.11.0-rc.1` bestehen nicht mehr). Der `main`-Link bleibt gültig. Der CHANGELOG-Eintrag nennt die Grenze in einem Satz.
- **Laufzeit.** `rootLongText` kompiliert die Regex bei jedem Aufruf (`root.go`, `regexp.MustCompile`); der Mehraufwand ist ein einmaliger Regex-Abgleich je Prozess. Die Grenze [`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort) (unter 200 ms) wurde nicht gemessen; der Docker-Lauf (1,5 s einschließlich Container-Start) taugt dafür nicht.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) und die beiden Hygiene-Slices (`slice-v1-cli-dry-run-log-would`, `slice-v1-harness-gofmt-und-depguard`) liegen in `done/`; kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).

**Rollenfolge:** Planner → Implementer → Reviewer → Verifier; ein Architect ist nicht nötig (eine Funktion im Adapter, keine neue Schnittstelle, kein Port, keine Spec-Änderung).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Der Projektinhaber ändert den Wortlaut, die Repo-URL oder den Handbuch-Pfad (etwa eine Veröffentlichung des Handbuchs an anderer Stelle); dann wird der Entwurf in §1 neu gefasst.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Wird mit diesem Slice eine Regel verkörpert, trägt die Closure-Notiz das Feld `liegt in` und das Ziel den Anker `seit slice-v1-cli-help-handbuch-link`.

## 6. Risiken und offene Punkte

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten / entfallen / weiter offen).

- **Toter Link bei Dev-Build, Fork oder umbenanntem Repo:** Dev-Builds und jede Version außerhalb der Release-Formen (`0.0.0-test`, `0.12.0-rc1`, leer) bekommen nur den `main`-Link; ein Fork oder ein umbenanntes Repo zeigt weiter auf `pt9912/u-boot`, und ein selbst gebautes Binary mit einer Release-Form als Version bekäme einen Tag-Link, dessen Tag es dort nicht gibt. — **Ausgang:** offen; Mitigation: Test pinnt die Regel; die Grenze (nur `X.Y.Z` und `X.Y.Z-rc.N` tragen den Tag-Link) steht im Code-Kommentar als Zusage. **Bei der Closure:** weiter offen für Fork, umbenanntes Repo und ein selbst gebautes Binary mit Release-Form als Version (der Text kennt den Ort des Repos nicht); der Dev-Teil ist belegt (Praxislauf mit vier Versionen: `0.11.0-dev` und `0.12.0-rc1` zeigen README und `main`, keinen Tag-Link).
- **Link zum Handbuch einer Version ohne Tag:** Ein Release-Build trägt `version` in einer Release-Form; der Tag entsteht erst mit dem Cut. Ein vor dem Tag gebautes Binary mit der Release-Version zeigt kurz ins Leere. — **Ausgang:** offen; Mitigation: `publish.yml` baut erst aus dem Tag, `make build` setzt `-dev`. **Bei der Closure:** weiter offen. Für `0.11.0` ist der Tag-Link nicht tot (`gh api` auf `docs/user/benutzerhandbuch.md?ref=v0.11.0` Exit 0; Tag und Release `v0.11.0` existieren). Offen bleiben ein vor dem Tag gebautes Binary und, aus dem Review (F-3), gelöschte Vorab-Tags: `docs/maintainer/releasing.md` räumt `X.Y.Z-rc.N` nach dem Cut auf (die Tags `v0.10.0-rc.1` und `v0.11.0-rc.1` bestehen nicht mehr); ein danach noch laufendes Binary zeigt einen toten Tag-Link, der `main`-Link bleibt gültig. Die Grenze steht im CHANGELOG-Eintrag.
- **Hartkodierung der Repo-URL:** `https://github.com/pt9912/u-boot` in `root.go` (schon heute für den README-Link). — **Ausgang:** offen; Mitigation: eine Konstante, ein Test pinnt sie; ein Umzug des Repos ist ein eigener Vorgang. **Bei der Closure:** weiter offen (unverändert zur heutigen Lage; die Mutation von `repoURL` im Review ließ Tests rot werden).
- **Pfad- oder Anker-Änderung des Handbuchs:** `docs/user/benutzerhandbuch.md` wird verschoben oder umbenannt; die Links in alten Binaries bleiben auf dem Tag gültig (Tag friert den Stand ein), der `main`-Link bricht. — **Ausgang:** offen; Mitigation: der Test nennt den Pfad; eine Umbenennung des Handbuchs zieht den Test (und damit den Code) nach. **Bei der Closure:** entfallen für diesen Slice: Der Pfad besteht auf dem Tag `v0.11.0` und auf `main` (`gh api` Exit 0, beide Refs), die Mutation von `manualPath` ließ im Review Tests rot werden; eine künftige Verschiebung des Handbuchs bricht den Test, nicht still den Link.
- **Der Test pinnt den Wortlaut zu eng:** Das Mutationsbeispiel zeigt, welche Teile (Tag mit `v`, Pfad, Release-Form-Regex) tragend sind; der Satz „(German)“ ist es nicht. — **Ausgang:** entfallen: Der Review mutierte sieben Stellen, sechs ließen Tests rot werden, eine überlebte bewusst (die Beschriftung `latest:` ist nicht gepinnt); die 13 Fälle prüfen hartkodierte Zeichenketten, keine aus den Konstanten abgeleiteten.
- **Laufzeit-Grenze [`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort) (unter 200 ms) ungemessen** (aus dem Review F-4, hier nachgetragen): `rootLongText` kompiliert die Regex bei jedem Aufruf; reine String-Logik ohne Netz, der Docker-Lauf (1,5 s einschließlich Container-Start) taugt nicht als Messung. — **Ausgang:** weiter offen; nicht gemessen, nicht behauptet.

## 7. Closure-Notiz

**Entscheidung des Projektinhabers (2026-10-18), wörtlich:** „können wir beides machen. eins auf den tag und eins auf das aktuelle handbuch“. **Festlegung des Leads (2026-10-18):** Der Tag-Link erscheint nur bei der Release-Form `X.Y.Z` oder `X.Y.Z-rc.N` (Zahlen ohne führende Null, Regex `releaseVersionPattern`); jede andere Version bekommt nur den `main`-Link. Das Lastenheft bleibt unverändert (kein CR; [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) nennt keine Dokumentations-Links).

**Belege:** Review-Report [`2026-10-18-slice-v1-cli-help-handbuch-link`](../../../reviews/2026-10-18-slice-v1-cli-help-handbuch-link.md) (0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO). Verification Evidence [`2026-10-18-verification-slice-v1-cli-help-handbuch-link`](../../../reviews/2026-10-18-verification-slice-v1-cli-help-handbuch-link.md): keine Diskrepanz zwischen Plan und Code, Auflagen 1 bis 5 in dieser Closure erfüllt. Commits `ab931a2`, `1cacc58`, `0861e23`, `e8a02c9`, `3213066`, `c3576b8`, `f6ddb18`, `1966568`, `8297111`.

- **Praxislauf (Runtime-Image aus dem Arbeitsbaum, vier Versionen):** `0.11.0` und `0.12.0-rc.1` zeigen README-, Tag- und `main`-Link; `0.12.0-rc1` und der Default `0.11.0-dev` zeigen README und `main`, keinen Tag-Link. `--version` unverändert; `init --help` ohne Handbuch-Zeile. `gh api` auf `docs/user/benutzerhandbuch.md` mit `ref=v0.11.0` und `ref=main`: Exit 0, Pfad vorhanden; Tag und Release `v0.11.0` existieren.
- **Abweichung (bewusst, sichere Richtung):** `0.12.0-rc1` (ohne Punkt) bekommt nur den `main`-Link, obwohl `publish.yml` (Schritt „Validate SemVer tag“) die Form als Tag zulässt. Die Regex ist eine echte Teilmenge der Tag-Prüfung; der Preis ist, dass ein Vorab-Binary dieser Schreibweise das Handbuch seiner Version nicht verlinkt, der Gewinn ein nie toter Tag-Link für alle anderen Formen. Die Praxis-Schreibweise `-rc.N` steht in `docs/maintainer/releasing.md` §1 und §5.
- **Der Cut trägt keinen Handarbeitsschritt:** Der Tag-Link folgt `a.version`, die der Cut ohnehin an den drei Stellen setzt (`cmd/uboot/main.go` `version`, `repo.mk`, `Dockerfile` `UBOOT_VERSION`; `docs/maintainer/releasing.md` §3 Schritt 2); `publish.yml` leitet die Version aus dem Tag ab (`VERSION="${TAG_REF#v}"`), auch für `-rc.N`. `releasing.md` und die Cut-Vorlage sind unverändert (nicht im Diff). Der Beleg im engeren Sinn steht aus: der Cut-Slice des nächsten Release, in dem kein Link von Hand gesetzt wird.

**Fundstellen je Aussage (`AGENTS.md` §3.10):**

| Aussage | Fundstelle |
|---|---|
| Der Hilfetext wird aus der Version zusammengesetzt; der Tag-Link steht nur bei einer Release-Form, der `main`-Link immer | `rootLongText` in `internal/adapter/driving/cli/root.go`; Aufruf `Long: rootLongText(a.version)` in `buildRootCommand` |
| Die Release-Form ist `X.Y.Z` oder `X.Y.Z-rc.N` ohne führende Null, eine Teilmenge der Tag-Prüfung | Konstante `releaseVersionPattern` in `root.go`; Gegenstück Schritt „Validate SemVer tag“ in `.github/workflows/publish.yml` |
| Repo-URL und Handbuch-Pfad sind Konstanten | `repoURL` und `manualPath` in `root.go` |
| Der Test pinnt beide Links, die Release-Formen und zehn Nicht-Formen mit hartkodierten Zeichenketten | `TestExecute_RootHelpManualLinks` und `helpForVersion` in `internal/adapter/driving/cli/cli_test.go` (13 Fälle); `TestExecute_RootHelpPointsToUserDocs` unverändert |
| Die Version kommt aus dem Wiring und ist beim Bau des Textes bekannt | `New(version, …)` und `App.version` in `internal/adapter/driving/cli/cli.go`; `Execute` baut das Wurzelkommando je Aufruf neu |
| Der CHANGELOG nennt Änderung und Grenze gelöschter Vorab-Tags | `## [Unreleased]` / `### Added` in `CHANGELOG.md` |
| Das Handbuch sagt, dass `--help` auf die Fassung der Version verweist | FAQ „Wo finde ich alle Optionen eines Befehls?“ in `docs/user/benutzerhandbuch.md` |
| Vorab-Tags werden nach dem Cut gelöscht (Grund der Grenze) | `docs/maintainer/releasing.md` §5 und §7 (Zeile „Vorab-Tags aufräumen“) |

**Plan-Tabelle gegen `git diff --stat 2ab2ee6` (§3.10, nach der letzten Änderung):** Der Diff nennt 15 Einträge. Mit Plan-Zeile: `root.go`, `cli_test.go`, `CHANGELOG.md`, `benutzerhandbuch.md`, der Review-Report und dieser Slice (doppelt: Löschung unter `in-progress/`, Neuanlage unter `done/`, weil Move und Inhalt getrennte Commits tragen). `roadmap.md` steht mit Plan-Zeile, ist aber **nicht** im Diff: Die Arbeitszeile wurde gesetzt und mit der Closure auf „Nichts in Arbeit.“ zurückgesetzt, der Netto-Diff ist leer. Closure-Dateien ohne Plan-Zeile (sie sind nicht Gegenstand des Plans): der Verification-Report, `open/README.md`, die Evidence-Datei und `state.md` von `BEO-DOK/plan-aussage-ueber-code`, und die vier Link-Nachzüge von `make slice-mv` in `done/slice-v1-cli-dry-run-log-would.md`, `done/slice-v1-harness-gofmt-und-depguard.md`, dem Review-Report und der Verification des Dry-Run-Slice. „Unverändert“ per gelesener Aussage: `export_test.go` (`WalkRootCommandPathsForTest` ruft `buildRootCommand(a)`, Signatur unverändert), `docs/maintainer/releasing.md` und die Cut-Vorlage (kein Diff; §3 Schritt 2 setzt nur die drei Versionsstrings), `README.md` und `README.de.md` (keine Aussage über den Hilfetext), Lastenheft und Spezifikation (kein Wortlaut des Hilfetexts).

**Finding-Klassen des Reviews:**
- **F-1 (LOW, behoben):** Kommentar beschreibt eine Zusage zu weit (`AGENTS.md` §3.7): „exactly the versions the publish workflow tags“ und „because no tag exists for it“ in `root.go`; der Workflow lässt mehr Vorab-Formen zu, die Regex ist eine Teilmenge. Behoben im Indikativ (Commit `c3576b8`).
- **F-2 (LOW, behoben):** Plan-Aussage nach der Festlegung des Leads nicht nachgezogen (§3.10, Klasse `plan-aussage-ueber-code`): §2 Test-Punkt und §6 erstes Risiko standen auf dem alten Stand. Behoben (Commit `1966568`).
- **F-3 (INFO, benannt):** Die Grenze gelöschter Vorab-Tags (§6 zweites Risiko, weiter offen; CHANGELOG-Halbsatz).
- **F-4 (INFO, benannt):** Laufzeit nicht gemessen; Regex je Aufruf kompiliert (§6, weiter offen).

- **Was hat funktioniert:** Der Praxislauf mit vier Versionen am Binary und die hartkodierten Strings im Test; die Mutationen des Reviews trafen jede tragende Stelle; die strengere Form (Teilmenge der Tag-Prüfung) machte aus „Tag-Link kann tot sein“ eine begrenzte, benannte Ausnahme.
- **Was ging anders als geplant:** Der erste Entwurf zog nur `-dev` ab; die Festlegung des Leads verengte auf die Release-Form, der Plan-Text folgte nur in §1 (F-2). Kommentare behaupteten mehr, als die Quelle (`publish.yml`) trägt (F-1).
- **Steering-Loop-Eintrag:** **Gezählt, nicht verkörpert.** Kein `liegt in`, keine Entscheidung über Verkörperung. (1) F-2 gehört zu `BEO-DOK/plan-aussage-ueber-code` (verkörpert in `AGENTS.md` §3.10): vierzehnter Vorgang, Rückfall nach der Verkörperung, Evidence-Datei angelegt, Stand fortgeschrieben; eine weitere Schärfung ist nicht entschieden. Die Lücke, die dieser Vorgang zeigt: §3.10 verlangt den Abgleich der **Tabelle** gegen `git diff --stat`; die Aussagen im Fließtext von DoD und Risiko, die eine Festlegung nachträglich ändert, fängt er nicht. (2) F-1 ist **nicht** `BEO-CLI/kommentar-beschreibt-abwesendes`, siehe Register. **Lehre:** Wechselt eine Festlegung die Zusage, sucht der Autor den neuen Wert im ganzen Plan (§1, §2, §6), nicht nur in dem Abschnitt, in dem er die Festlegung einträgt.
- **Beobachtungs-Register (`../observations/`):** Sichtung von `BEO-CLI` (`dry-run-log-vergangenheitsform`, gestrichen; `kommentar-beschreibt-abwesendes`, verkörpert) und `BEO-DOK` (`plan-aussage-ueber-code`, `spec-bedingung-ohne-code-entsprechung`, beide verkörpert): keine Beobachtung zum Help-Text oder zu Links in Hilfetexten. [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md): eine weitere Evidence-Datei (F-2), Stand fortgeschrieben. [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md): **unverändert, keine Evidence-Datei für F-1.** Begründung: Die Beobachtung fasst die drei Muster von `AGENTS.md` §3.9 (a Vorzustand erzählt, b Verweis auf geänderte Symbole, c Helfer am falschen Funktionskopf): Der Kommentar war einmal richtig oder steht am falschen Ort, der Diff macht ihn falsch. F-1 ist ein anderer Mechanismus: Der Kommentar war von Anfang an zu weit („exactly“, „no tag exists“), weil die Quelle (`publish.yml`) nicht gelesen oder falsch gelesen war; kein Diff hat ihn veraltet, §3.9 (a) bis (c) hätten ihn nicht gefangen. Nah dran ist die Klasse von `plan-aussage-ueber-code` (Aussage über Code nicht an der Stelle gelesen), die §3.10 aber nur für Plan, Closure-Notiz und Spec nennt, nicht für Code-Kommentare. Einmalig, daher benannt und kein neuer Pfad (ein Zähler 1 ohne Wiederauftreten sagt nichts; tritt „Kommentar behauptet eine Zusage über eine andere Quelle, ohne sie zu lesen“ in einem zweiten Vorgang auf, legt dessen Closure den Pfad an). **Nicht entschieden hier:** ob F-1 doch unter `kommentar-beschreibt-abwesendes` zählt (Projektinhaber).
- **Folge-Slices:** keine.
- **Risiken aus §6:** alle sechs tragen einen Ausgang: weiter offen (Fork, Selbstbau; Link zu Version ohne Tag samt gelöschten Vorab-Tags; Hartkodierung der URL; Laufzeit-Grenze ungemessen), entfallen (Pfad-Änderung des Handbuchs; Test pinnt zu eng).
- **Roadmap:** Die Arbeitszeile steht nach dem Move wieder auf „Nichts in Arbeit.“ (`in-progress/` trägt keinen Slice).
- **`open/README.md`:** Die Zeile dieses Slice zeigt auf `done/` und trägt „geliefert“. Damit trägt `open/` außer dem README keinen Slice mehr; jede Tabellenzeile verweist auf `done/`.
- **Drei Paarungen:** Anker (kein `liegt in`, nichts verkörpert) · Folge-Slice (keiner) · Register (ein Pfad mit neuer Evidence-Datei, zitiert in dieser Notiz) — nach dem `git mv` nach `done/` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/adapter/driving/cli` (CLI, Greenfield) und `spec/`, `harness/`, `docs/plan/` (DOK) sowie `docs/user/` (Handbuch) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `cmd/uboot` bleibt unberührt (keine Sub-Area).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) ist bei der Anlage nur auf Verzeichnis-Ebene gesichtet (Kürzel `BEO-APP`, `BEO-CLI`, `BEO-DOK`, `BEO-E2E`, `BEO-WZG`); eine Beobachtung zum Help-Text oder zu Links in Hilfetexten wurde bei der Anlage nicht gesucht. **Gesichtet bei der Closure:** die Einträge von `BEO-CLI` (`dry-run-log-vergangenheitsform`, gestrichen; `kommentar-beschreibt-abwesendes`, verkörpert in `AGENTS.md` §3.9) und `BEO-DOK` (`plan-aussage-ueber-code`, `spec-bedingung-ohne-code-entsprechung`, beide verkörpert in `AGENTS.md` §3.10) gelesen: keine Beobachtung zum Help-Text oder zu Links in Hilfetexten; die beiden verkörperten Regeln galten für Kommentare und Plan-Aussagen dieses Slice (Ergebnis in §7). Die Sichtung ist vom Planer nachgeholt, nicht vom Implementer vor dem Start.

Alle berührten Sub-Areas GF.
