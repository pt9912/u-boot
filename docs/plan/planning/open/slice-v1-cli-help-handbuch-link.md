# Slice v1: `u-boot --help` verweist auf das Handbuch — Link auf den Tag und Link auf `main`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — kein Release-Bezug festgelegt; der Slice läuft **nach** dem Cut [`slice-v1-release-cut-v0.11.0`](../done/slice-v1-release-cut-v0.11.0.md) und nach den beiden Hygiene-Slices [`slice-v1-cli-dry-run-log-would`](../done/slice-v1-cli-dry-run-log-would.md) und [`slice-v1-harness-gofmt-und-depguard`](../in-progress/slice-v1-harness-gofmt-und-depguard.md); er ist nicht Teil von `v0.11.0`. Herkunft: Wunsch des Projektinhabers 2026-10-18 („bei u-boot --help wird das readme referenziert. Können wir noch unser Handbuch in die Meldung mitaufnehmen“).

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
2. **Dev-Stände (`…-dev`)?** Der Tag existiert nicht; ein Link wäre tot (404). Entscheidung: für eine Version mit Suffix `-dev` (und für eine leere Version) zeigt die Hilfe **nur** den `main`-Link; das Abschneiden von `-dev` wird verworfen, weil es auf den Tag des **letzten** Release (oder keinen) zeigte und Handbuch und Binary nicht mehr zusammenpassten, ohne dass der Benutzer es merkt. Vorabversionen (`0.12.0-rc1`) und alle übrigen Versionen tragen ihren Tag (`publish.yml` legt sie unter `v<version>` an) und bekommen beide Links.
3. **Wo lebt der Text?** In `root.go`: `Long` wird aus einer kleinen reinen Funktion (Arbeitsname `rootLongText(version string) string`) zusammengesetzt, die ein Format-String-Muster mit `fmt.Sprintf` füllt; die Repo-URL steht als Konstante daneben. Kein Template-Paket, kein Port, kein neuer Typ.
4. **Wortlaut (Entwurf, Englisch).**

   ```text
   Quickstart and user docs: https://github.com/pt9912/u-boot#readme
   User manual (German), this version: https://github.com/pt9912/u-boot/blob/v0.11.0/docs/user/benutzerhandbuch.md
   User manual (German), latest: https://github.com/pt9912/u-boot/blob/main/docs/user/benutzerhandbuch.md
   ```

   Bei `…-dev` entfällt die mittlere Zeile; die letzte Zeile heißt dann ebenfalls „latest“. Der Wortlaut darf kein `LH-` enthalten (bestehender Test).
5. **Handarbeit pro Cut.** Keine: Der Tag-Link folgt `version`, die der Cut ohnehin an den drei Stellen setzt; `docs/maintainer/releasing.md` und die Vorlage der Cut-Slices bleiben unverändert (Beleg: §3 Schritt 2 setzt nur die Versionsstrings, der Link wird zur Laufzeit abgeleitet).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Link auf das Handbuch in den Hilfetexten der Unterbefehle:** Das Wurzelkommando ist der Einstieg, den der Wunsch nennt; je Befehl ein Link wäre eine eigene Zusage mit Anker-Pflege.
- **Eine Prüfung, dass der Tag oder die Seite im Netz existiert:** `--help` bleibt netzlos und unter 200 ms ([`LH-NFA-PERF-001`](../../../../spec/lastenheft.md#lh-nfa-perf-001--schnelle-cli-antwort): `u-boot --help` und `u-boot --version` unter 200 ms); die Ableitung aus der Version ist die Garantie.
- **Ein Konfigurationsschalter für die Repo-URL (Fork, Umbenennung):** Die URL stand schon im README-Link hartkodiert; ein Schalter wäre eine Zusage ohne Nachfrage. Das Risiko steht in §6.
- **Eine Änderung der Versionsausgabe (`--version`) oder der Versionsstrings:** unberührt.
- **Ein Gate, das die Links im Netz prüft:** Der Test pinnt den Text (§2); ein Link-Checker gegen GitHub wäre ein neues, netzgebundenes Gate mit eigener Begründung.

## 2. Definition of Done

- [ ] **Code:** `u-boot --help` enthält neben dem README-Link den Tag-Link und den `main`-Link (Wortlaut §1, Entwurf 4); bei einer Version mit Suffix `-dev` (und bei leerer Version) enthält sie nur den `main`-Link und **keinen** Tag-Link ([`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe), [`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache)).
- [ ] **Test:** ein Test baut den `App` mit `0.11.0` und prüft beide Links vollständig (Tag `v0.11.0`, `main`, Pfad `docs/user/benutzerhandbuch.md`); ein zweiter mit `0.11.0-dev` prüft den `main`-Link und die Abwesenheit von `/blob/v`; ein dritter mit einer Vorabversion (`0.12.0-rc1`) prüft den Tag-Link `v0.12.0-rc1`. Mutation: Der Tag-Link ohne `v`, die Verzweigung auf `-dev` entfernt und der `main`-Link durch den Tag-Link ersetzt — jede Mutation lässt mindestens einen Test rot werden, am Beispiel gesehen und zurückgenommen. Der bestehende Test `TestExecute_RootHelpPointsToUserDocs` bleibt grün (README-Link, kein `LH-`, keine Katalog-Begriffe).
- [ ] Doku-Update ([`AGENTS.md`](../../../../AGENTS.md) §5 Zeile 5): `CHANGELOG.md` `[Unreleased]` nennt die Änderung (Added: `u-boot --help` zeigt den Handbuch-Link der Version und den des aktuellen Handbuchs; Dev-Builds nur den aktuellen); das Handbuch erhält einen Satz im FAQ „Wo finde ich alle Optionen eines Befehls?“ (die Hilfe verweist auf das Handbuch Ihrer Version), oder §7 begründet, warum nicht.
- [ ] `docs/maintainer/releasing.md` und die Cut-Vorlage unverändert; §7 hält fest, dass der Cut keinen neuen Handarbeitsschritt trägt (der Cut-Slice des nächsten Release ist der Beleg: kein Link wird von Hand gesetzt).
- [ ] `make gates` grün; `make doc-trace` ohne neue Waise.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Jede Aussage über Code steht mit Fundstelle (Datei und Symbol, gelesen am 2026-10-18). Der Implementer gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10) und belegt jede Datei, die er als „unverändert“ führt, durch eine gelesene Aussage.

| Datei / Komponente | Änderungs-Art | Begründung (Fundstelle) |
|---|---|---|
| `internal/adapter/driving/cli/root.go` | update | `buildRootCommand`: `Long` (Raw-String, endet mit dem README-Link) wird aus einer reinen Funktion der Version (`a.version`) zusammengesetzt; neue Konstante für die Repo-URL und die Handbuch-Pfade; Verzweigung für Suffix `-dev` und leere Version |
| `internal/adapter/driving/cli/cli_test.go` | update | neue Tests neben `TestExecute_RootHelpPointsToUserDocs`; die Fabrik `newApp` baut mit `"0.0.0-test"`, die neuen Tests bauen den `App` mit eigener Version über `cli.New(version, …)` (die Argumentliste der Fabriken zeigt die Reihenfolge); der bestehende Test bleibt unverändert |
| `internal/adapter/driving/cli/export_test.go` | prüfen | gelesen: `WalkRootCommandPathsForTest` ruft `buildRootCommand(a)` und liest nur die Pfade; ändert sich die Signatur von `buildRootCommand` nicht, bleibt die Datei unberührt |
| `CHANGELOG.md` | update | `## [Unreleased]` (derzeit nur `### Changed` des Dry-Run-Slice) bekommt `### Added`; Englisch, wie der Bestand |
| `docs/user/benutzerhandbuch.md` | update, ein Satz | FAQ „Wo finde ich alle Optionen eines Befehls?“ (Zeile mit `u-boot <befehl> --help`); die Optionstabelle im Anhang (`-h`, `--help`: „Hilfe zum Befehl“) bleibt, sofern der Satz im FAQ genügt |
| `README.md`, `README.de.md` | prüfen | gelesen: die Suche nach `help` in `README.md` findet nur `make help`; kein Text über den Hilfetext, keine Änderung geplant |
| `spec/spezifikation.md`, `spec/lastenheft.md` | unverändert | gelesen (siehe Kopf): kein Wortlaut des Hilfetexts; [`LH-FA-CLI-002`](../../../../spec/lastenheft.md#lh-fa-cli-002--hilfeausgabe) nennt keine Dokumentations-Links |
| `docs/maintainer/releasing.md`, Cut-Vorlage | unverändert | gelesen: §3 Schritt 2 setzt die drei Versionsstrings; der Link wird aus `version` abgeleitet, es entsteht kein Handarbeitsschritt |
| `docs/reviews/<Datum>-slice-v1-cli-help-handbuch-link.md` | neu, durch den Reviewer | Review-Report |

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

- **Toter Link bei Dev-Build, Fork oder umbenanntem Repo:** Dev-Builds sind durch die `-dev`-Verzweigung abgedeckt; ein Fork oder ein umbenanntes Repo zeigt weiter auf `pt9912/u-boot`. Ein selbst gebautes Binary mit einer beliebigen Version (etwa `0.0.0-test`) bekäme einen Tag-Link, dessen Tag nicht existiert. — **Ausgang:** offen; Mitigation: Test pinnt die Regel; die Grenze (nur `-dev` ist ausgenommen) steht im Code-Kommentar als Zusage.
- **Link zum Handbuch einer Version ohne Tag:** Ein Release-Build trägt `version` ohne `-dev`; der Tag entsteht erst mit dem Cut. Ein vor dem Tag gebautes Binary mit der Release-Version zeigt kurz ins Leere. — **Ausgang:** offen; Mitigation: `publish.yml` baut erst aus dem Tag, `make build` setzt `-dev`.
- **Hartkodierung der Repo-URL:** `https://github.com/pt9912/u-boot` in `root.go` (schon heute für den README-Link). — **Ausgang:** offen; Mitigation: eine Konstante, ein Test pinnt sie; ein Umzug des Repos ist ein eigener Vorgang.
- **Pfad- oder Anker-Änderung des Handbuchs:** `docs/user/benutzerhandbuch.md` wird verschoben oder umbenannt; die Links in alten Binaries bleiben auf dem Tag gültig (Tag friert den Stand ein), der `main`-Link bricht. — **Ausgang:** offen; Mitigation: der Test nennt den Pfad; eine Umbenennung des Handbuchs zieht den Test (und damit den Code) nach.
- **Der Test pinnt den Wortlaut zu eng:** Das Mutationsbeispiel zeigt, welche Teile (Tag mit `v`, Pfad, `-dev`-Ausnahme) tragend sind; der Satz „(German)“ ist es nicht. — **Ausgang:** offen.

## 7. Closure-Notiz

Wird bei der Closure gefüllt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `internal/adapter/driving/cli` (CLI, Greenfield) und `spec/`, `harness/`, `docs/plan/` (DOK) sowie `docs/user/` (Handbuch) — alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#mr-006), jede erfüllt die Schwelle ≥ 2 von 3 Achsen; keine Ausdifferenzierung nötig. `cmd/uboot` bleibt unberührt (keine Sub-Area).

**Vorgelagert — offene Beobachtungen sichten:** Beobachtungs-Register ([`../observations/`](../observations/README.md)) ist bei der Anlage nur auf Verzeichnis-Ebene gesichtet (Kürzel `BEO-APP`, `BEO-CLI`, `BEO-DOK`, `BEO-E2E`, `BEO-WZG`); eine Beobachtung zum Help-Text oder zu Links in Hilfetexten wurde nicht gesucht. Der Implementer sichtet die Einträge von `BEO-CLI` und `BEO-DOK` vor dem Start und trägt das Ergebnis hier ein. `BEO-DOK/plan-aussage-ueber-code` (verkörpert, `AGENTS.md` §3.10) gilt für die Plan-Tabelle.

Alle berührten Sub-Areas GF.
