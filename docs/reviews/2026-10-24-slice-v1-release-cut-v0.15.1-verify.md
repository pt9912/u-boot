# Verification-Report: slice-v1-release-cut-v0.15.1

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff 6335bf18..HEAD` (HEAD `f16b89c2`), Slice
  [`slice-v1-release-cut-v0.15.1`](../plan/planning/in-progress/slice-v1-release-cut-v0.15.1.md)
  (Patch `v0.15.1`, Einstufung *Fixed*); Muster [`slice-v1-release-cut-v0.15.0`](../plan/planning/done/slice-v1-release-cut-v0.15.0.md)
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-CONF-005`](../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)
* Review: [`2026-10-24-slice-v1-release-cut-v0.15.1.md`](2026-10-24-slice-v1-release-cut-v0.15.1.md) (LOW-1 und INFO-1 in `f16b89c2` behoben)
* Belege der Inhalte: Verify-Reports [`config set`](2026-10-24-slice-v1-config-set-validiert-feature-version-verify.md) und [`generate`](2026-10-24-slice-v1-generate-prueft-freigabeliste-verify.md)

## Urteil

**Push (M1) freigabereif: ja, mit einer Empfehlung vor dem Push (LOW-V1).** Alle Sensoren, die ich fuhr, sind
grün; Tap-Pfad (M-T) mit roter Gegenprobe reproduziert; keine Tag-/Push-Aktion, keine vergessene Versionsstelle.
Die Fix-Runde `f16b89c2` hat LOW-1 des Reviews nur zur Hälfte geschlossen: Die Messaussage nennt `gh`, `glab`
und `java` für **beide** Fixes am Binary, belegt ist das nur für `config set`. Das ist eine Textkorrektur
(CHANGELOG-Lead, Status-Block und Tabellenzeile beider READMEs) und gehört **vor** den Push, weil der Lead
nach dem Tag nicht mehr geändert wird (`releasing.md` §8). Wird sie nicht gemacht, bleibt ein bekanntes LOW im
Release-Text.

## Sensoren (selbst gefahren, Exit-Codes sichtbar)

| Sensor | Ergebnis |
|---|---|
| `make gates` | `EXIT=0` (Log bis `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`) |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen |
| `make doc-planning` | Exit 0, 803 Dateien, 0 Befunde |
| `make doc-commits RANGE=d98f5755..HEAD` | Exit 0, 10 Commits, 0 Befunde |
| `make doc-immutable RANGE=v0.15.0..HEAD` | Exit 0, 47 Commits, 0 Befunde |
| `make build VERSION=0.15.1` | Exit 0; `docker run --rm u-boot:latest --version` = `u-boot version 0.15.1`; Label `org.opencontainers.image.version` = `0.15.1` |
| `make test-docker` | nicht gelaufen: Auftrag (Speicher); der Produktcode-Diff gegenüber `v0.15.0` ist die Zeile `var version` |
| `make verify-depguard` | nicht gelaufen: keine Importregel berührt |
| `make ci`, `make fullbuild` | nicht gefahren; Darstellung des Implementers (Trivy-DB 404, Wiederholung grün) von mir nicht reproduziert |

## M-T: Tap-Trockenlauf, selbst reproduziert

`docker run -i --rm --network none -v /Development/u-boot:/src:ro u-boot:test bash < <Skript>`; Skript im
Scratchpad (nicht im Repo), eigene Stellvertreter für `curl` und `git` im `PATH`, Fixture-Summen wie in
`cmd/uboot/homebrew_scripts_test.go`, Tap-Token `dummy`. Kein Netz, kein Schreiben ins Tap.

| Schritt | Ergebnis |
|---|---|
| `homebrew-formula-fill.sh v0.15.1` | Exit 0, vier `url`-Zeilen mit `v0.15.1` |
| `tap-formula-check.sh v0.15.1` | Exit 0 |
| `tap-formula-check.sh v0.15.0` (rot) | Exit 1, „found 4 url lines, 0 of v0.15.0“ |
| eigene Gegenprobe: eine URL mit `v0.15.0`, Prüfung auf `v0.15.1` | Exit 1, „3 of v0.15.1“ |
| `tap-nachzug.sh`, `TAG=v0.15.1` | Exit 0; Protokoll `curl`, `git clone --depth 1`, `git diff --quiet`, `add`, `commit -m u-boot 0.15.1`, `git push origin HEAD` |
| `tap-nachzug.sh`, `TAG=v0.15.0` (rot) | Exit 1; Protokoll nur der `curl`-Aufruf, kein `clone`, kein `push` |

**Beweist nicht:** `is_stable`, echter `curl` gegen das Asset, echter `git push`, Token; der `git`-Stellvertreter
meldet immer eine Änderung. `tap-check` (M4) bleibt Sache des Hauptagenten.

## Nachlauf der Fix-Runde `f16b89c2`: Messaussage

Aussage im CHANGELOG-Lead, Status-Block und Tabellenzeile beider READMEs: „both fixes at the built binary on one
Linux host … for `gh`, `glab` and `java`; `ollama` pinned by tests only“.

| Teilaussage | Beleg | Befund |
|---|---|---|
| `config set`: Binary-Proben mit `gh`, `glab`, `java` | Verify `config set`, *Ende-zu-Ende*: `gh.version` Exit 10 (vier Aufrufformen), `java.version` Exit 0, `glab.version` vor `glab.source` Exit 10, `gh`-Paare Exit 0 | belegt |
| `ollama` nur durch Tests | Verify `config set`: `ollama`/`gh`/`glab` je Fall in beiden Testfunktionen, am Binary keine `ollama`-Probe | belegt |
| `generate`: Binary-Proben „for `gh`, `glab` and `java`“ | Verify `generate`, *Ende-zu-Ende*: ausschließlich Eintrag `evil` (`source: https://evil.example.org/feat`), zusätzlich Katalogschlüssel `node` ohne `source` → Exit 0; `gh`, `glab`, `java` kommen im Report nicht vor | **nicht belegt, siehe LOW-V1** |
| „Linux host“ | Slice `generate` §7 nennt „Linux `amd64`“; Verify `generate` und Verify `config set` nennen weder Betriebssystem noch Architektur der Verifier-Läufe | schwach: Linux ist für `generate` im Slice belegt, für `config set` nur aus dem Docker-Aufruf zu schließen; der Lead nennt kein `amd64` (richtig) |
| Mutationstests der Verifier | beide Reports tragen Mutationstabellen; der Verifier `generate` fuhr einen Teil nicht selbst nach („Nicht von mir nachgefahren: M1b, M2, M5, M6, M8a bis M8g, M9, M11, M12“) | belegt, Formulierung „of the verifiers“ deckt es |
| Ungemessen-Liste (Whitespace am Rand, `doctor` ohne `enabled`, `--json` der Ablehnung bei `config set`, macOS, Podman, andere Hosts) | beide *Ungemessen*-Abschnitte | deckungsgleich |

### LOW-V1

* `kategorie`: LOW
* `quelle`: `AGENTS.md` §3.10 (eine Aussage nennt nur, was gemessen wurde), `BEO-DOK/zusage-weiter-als-messung`
* `pfad`: `CHANGELOG.md` `[0.15.1]` (Absatz *Measured*); `README.md` und `README.de.md` Status-Block (Satz „both are measured … for `gh`, `glab` and `java`“ / „beides ist … für `gh`, `glab` und `java` gemessen“); Tabellenzeile `v0.15.1` beider READMEs
* `befund`: Die Fix-Runde hat den Messort für `config set` richtig auf `gh`, `glab`, `java` und `ollama` nur durch Tests gestellt, die Aufzählung aber auf **beide** Fixes bezogen. Der Fix `generate` ist am Binary mit einem benutzerdefinierten Feature `evil` und dem Katalogschlüssel `node` gemessen; kein Binary-Lauf von `generate` mit `gh`, `glab` oder `java` liegt in einem Report. Die Freigabeprüfung hängt an der Quelle des Eintrags, nicht am Feature-Namen; der Code liefert also keinen Anlass zur Sorge, aber die Zusage geht weiter als die Messung (dieselbe Klasse wie LOW-1 des Reviews).
* `risiko`: Leser nimmt `generate` mit `gh` oder `glab` als am Binary gemessen. Gering, weil die Regel namensunabhängig ist und service-seitig gepinnt.
* `empfehlung`: Aussage trennen: `config set` am Binary mit `gh`, `glab`, `java` (`ollama` nur durch Tests); `generate` am Binary mit einem Eintrag einer nicht freigegebenen Quelle (Standard- und Sandbox-Profil). Gleichlautend an den vier Stellen. Danach `make gates` und ein kurzer Nachlauf-Blick des Reviewers auf den geänderten Wortlaut; Datum prüfen (M0).
* `verifizierbar`: Verifier (Reports gelesen)

## DoD des Cuts gegen `releasing.md` (M0, M-T, M1 bis M8)

Stand: vor dem Push. „erfüllt“ nur, was ich selbst gesehen habe.

| Punkt | Stand vor dem Push | Beleg |
|---|---|---|
| M0 Datum | **erfüllt**: `date +%F` = `2026-10-24` = CHANGELOG `[0.15.1] - 2026-10-24`, Tabellenzeile beider READMEs, Handbuch-`Stand` | eigener Lauf; gilt nur, solange Push und Tag am selben Tag liegen |
| Versionsstrings | **erfüllt**: `0.15.1-dev` in `cmd/uboot/main.go:43`, `repo.mk`, `Dockerfile`; `0.15.0-dev` nirgends mehr | eigene Suche; Build mit `VERSION=0.15.1` meldet `0.15.1` |
| Vergessene `0.15.0`-Stellen | **keine**: `grep -rIn "0\.15\.0"` außerhalb `.git`, `.harness`, `planning`, `reviews`, `messungen`, `CHANGELOG.md` trifft nur die zwei Tabellenzeilen `v0.15.0` (Geschichte), `releasing.md:57` (Zeiger „zuletzt“, Closure M8) und die Spec-Versionen `lastenheft.md:1591`, `spezifikation.md:1435` | eigene Suche |
| CHANGELOG | **erfüllt**: `[Unreleased]` leerer Anker; Vergleichslinks `[Unreleased]: v0.15.1...HEAD`, `[0.15.1]: v0.15.0...v0.15.1`; Einträge *Fixed* und Text unter `[0.15.0]` unverändert; im Block `[0.15.1]` kein „security“, „vulnerab“, kein CI-Satz | Diff gelesen, Suche |
| README / Handbuch | **erfüllt**: kennungsfrei (Diff auf `LH-`, `ADR-`, `SPEC-`, Slice, Roadmap gelesen: keine Treffer); Zeilen `v0.15.0` unberührt; Handbuch nur Kopf und Versionsbeispiele | Diff gelesen |
| M-T Tap-Trockenlauf | **erfüllt und nachgefahren**, rote Gegenprobe gesehen | Abschnitt oben |
| `make gates` | **erfüllt** (Exit 0, selbst gefahren) | oben |
| `make ci` / `make fullbuild` | **ungeprüft von mir**; Implementer: Exit 0 nach einer Trivy-DB-404-Störung | Review: nicht reproduziert |
| `make test-docker` lokal | **nicht gelaufen**, Planänderung dokumentiert in DoD und §3 Punkt 8 (INFO-1 des Reviews damit geschlossen; der Haken bleibt offen, weil Punkt 2 der DoD auf M2 wartet) | Slice §2, §3 |
| Review-Report | **erfüllt** (liegt vor, LOW-1/INFO-1 behoben; LOW-V1 neu) | |
| Verification-Report | dieser Bericht (vor M1) | |
| Zustandsmarker | Roadmap-Zeiger („in Arbeit“, Push und Tag stehen aus), `open/README.md` „in Arbeit“: richtig für den Stand vor dem Push; Ablaufpunkt Closure | INFO-2 des Reviews bleibt Closure-Pflicht |
| M1 Push | **steht aus**: Hauptagent, nach ausdrücklicher Freigabe | `origin/main` = `d98f5755`, 10 Commits lokal davor, Arbeitsbaum sauber |
| M2 CI auf dem Push-Tip | **steht aus, ungemessen**: `ci` (`gates`, `security-gates`, `image-scan`) und `integration-docker` (4 `--- SKIP`, 0 `--- FAIL`, Dauer gegen 30 Minuten, Job-Log über REST-API und Annotationen, nicht per Textsuche) | keine CI-Aussage im Release-Text, richtig |
| M3 Vorab-Tag | entfällt (Entscheidung beim Start, Slice §7) | |
| M4 `tap-check` | steht aus, unmittelbar vor M5; M-T beweist das Token nicht | |
| M5 Tag | **steht aus**: nur nach eigener ausdrücklicher Freigabe; Tag auf den in M2 geprüften Commit | `git tag -l 'v0.15*'` = nur `v0.15.0` auf `a895ffd2` |
| M5b Advisory | entfällt (Einstufung *Fixed*, Projektinhaber bestätigt laut Slice §7; ein Wortlaut des Projektinhabers liegt dort nicht vor, der Hauptagent benennt sie als bestätigt) | |
| M6 Kontrolle `publish`, GHCR, zwölf Assets, Job `tap` | steht aus, nach dem Tag | |
| M7 Mac | **ungemessen, offen** (Projektinhaber) | |
| M8 Nacharbeit | steht aus: Roadmap-Zeile mit Tag-Hash und Datum, Zeiger „zuletzt“ in `releasing.md:57`, Closure-Notiz, Register, `git mv` nach `done/` (Closure erst nach M6) | |

## Außenwirkungs- und Historien-Prüfung

* **Keine Tag-/Push-Aktion:** `git tag -l 'v0.15*'` = nur `v0.15.0`, `git rev-parse 'v0.15.0^{commit}'` = `a895ffd2`; `origin/main` = `d98f5755`, `HEAD` = `f16b89c2`.
* **Zwischenstände nicht Spitze:** die roten Zwischenstände der reinen Move- und Verweis-Nachzug-Commits (`f246e155`, `b1d94f5b`, `f053b1ec`, `617cd0d8`) liegen unter der Spitze; sie gehören in **einen** Push (M1: „alle Commits des Cuts in einem Push“). Der Tip `f16b89c2` (oder ein Nachfolger der LOW-V1-Korrektur) ist der Stand, den M2 liest.
* **Commits:** jede Message nennt Slice und `LH-*`; keine `SPEC-`/`ARC-`-Kennung; `doc-commits` 0 Befunde.
* **Harte Regeln:** kein `//nolint`, kein Gate gelockert, keine Spec-, ADR- oder Workflow-Änderung im Diff.

## Was dem Push vorausgehen muss

1. LOW-V1 berichtigen (ein Commit; Wortlaut an vier Stellen), `make gates` erneut, `make docs-check`; Reviewer-Blick auf den neuen Wortlaut.
2. M0 am Tag des Pushs wiederholen (`date +%F` = `2026-10-24`); weicht es ab, neuer Cut-Commit.
3. Ausdrückliche Freigabe des Projektinhabers im Gespräch für M1. Danach M2 mit gelesenem Job-Log; `tap-check` (M4) und eine eigene Freigabe für M5.

## Nicht geprüft / ungemessen

* CI auf dem Push-Tip (M2), `integration-docker` und die vier übersprungenen Nested-Podman-Tests.
* `make test-docker` lokal, `make verify-depguard`, `make ci`, `make fullbuild` (nicht reproduziert).
* `is_stable`, echter `curl`/`git push`, Token (M-T beweist sie nicht); `tap-check`.
* Mac, `brew install` (M7).
* Plattform und Docker-Version der Verifier-Läufe der beiden Fix-Slices (in den Reports nicht genannt).
