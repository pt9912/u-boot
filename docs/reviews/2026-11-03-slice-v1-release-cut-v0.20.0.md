# Review-Report: slice-v1-release-cut-v0.20.0 — 2026-11-03

**Review-Art:** Code/Release (Cut-Diff gegen Plan, Entscheidungen und Hard Rules;
[`releasing.md`](../maintainer/releasing.md) §§3 bis 8 als Ablauf-Quelle; Lead des
Blocks `[0.20.0]` gegen Lastenheft `0.31.0` und [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)).
Nicht geprüft: DoD-Closure und Verification-Evidence im Ganzen — Verifier.

**Gegenstand:** [`slice-v1-release-cut-v0.20.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.20.0.md),
Range `3bc9757f..90f2a1d2` — Plan-Anlage `3bc9757f` (außerhalb des Diff-Umfangs,
als Kontext gelesen), reine Moves `a25484a8`/`70262862` (R100 100 %, Verweis-Nachzug
eigener Commit), Roadmap-Start `ad4e93e1`, Cut `90f2a1d2` (7 Dateien). Tag
`v0.20.0` auf `90f2a1d2` gesetzt und gepusht (Freigabe des Projektinhabers
2026-11-03 — Auftragskontext); Release gepubliziert (12 Assets, `publishedAt`
2026-11-03T05:26:34Z — eigener Lesezugriff über `gh release view`).

**Skill:** `.harness/skills/reviewer.md` @ `v0.20.0`
**Modell:** glm-5.3-flash · **Datum:** 2026-11-03 · **Rolle:** Reviewer (Modul 10)

**Eingangs-Kontext** (gelesen, gegen die geprüft wurde):

- Slice-Plan §1–§6 samt Risiken; Vorgänger-Plan
  [`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md)
  (Form-Gleichheit, Lehren F-5/M-F-5 und Marker-Klasse)
- [`releasing.md`](../maintainer/releasing.md) §§3–8; Vorgänger-Review
  [`2026-11-02-slice-v1-release-cut-v0.19.0.md`](2026-11-02-slice-v1-release-cut-v0.19.0.md)
  (Finding-Klassen als Suchbild)
- Lastenheft `0.31.0`: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
  Absatz `java-toolchain` (Zeile 427), Historie-Zeile `0.31.0` (Zeile 1631);
  [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
  samt Beleg-Tabelle und „Nicht zugesagt“
- `CHANGELOG.md` Blöcke `[0.20.0]` und `[0.19.0]` (Lead-Stil),
  `README.md`/`README.de.md`, `docs/user/benutzerhandbuch.md`
- [`AGENTS.md`](../../AGENTS.md) §§3.3, 3.7, 3.9, 3.10, 5 Regel 6;
  [`roadmap.md`](../plan/planning/in-progress/roadmap.md) und
  [`open/README.md`](../plan/planning/open/README.md) in der Fassung am Tag-Commit

---

## Findings

**F-1 — MEDIUM** [`docs/plan/planning/in-progress/slice-v1-release-cut-v0.20.0.md §7`](../plan/planning/in-progress/slice-v1-release-cut-v0.20.0.md) ·
"Noch nicht geschlossen; der Slice liegt in `open/` (Plan)." — Abgelaufener
Zustandsmarker im Slice (identischer Wortlaut, dritter Lauf)
Klasse: Zustandsmarker nach Ablauf nicht abgebaut
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt),
[`harness/review.md`](../../harness/review.md) MEDIUM (Wiederholung → Drift-Signal)
Befund: Der Slice liegt seit `a25484a8` in `in-progress/`; §7 trägt unverändert
den Plan-Moment `open/`. Denselben Satz trug der `v0.19.0`-Cut (dort F-4,
MEDIUM) und der `v0.18.0`-Cut (dort F-4, LOW) — der Marker-Klasse-Wortlaut
steht jetzt bei drei Läufen und wurde in keiner Berichtigungsrunde geändert.
Soll–Ist: der Marker ist beim Move abgelaufen; die §7-Fassung ist erst mit der
Closure-Notiz ohnehin neu zu setzen.
Risiko: Der nächste Änderer (Verifier, Closure-Lauf) liest „liegt in `open/`“
als Ort-Aussage gegen die Verzeichnis-Wahrheit; kein Gate prüft
Statusmarker gegen das Lifecycle-Verzeichnis.
Verifizierbar: Review-only (d-check prüft Links, nicht Zustandsmarker)

**F-2 — MEDIUM** `docs/plan/planning/in-progress/roadmap.md` (Tabelle *Nächste
Wellen*, Zeile 34) · "`v0.20.0` — Feature-Katalog: Quellen, Optionen und
Basis-Image" — Vergebene Vorschlags-Version nach dem Tag nicht auf den Ist-Stand
Klasse: Vorschlags-Version vergeben (3. Befund) / Zustandsmarker nach dem Tag
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker beim Ablauf auf den
Ist-Stand), [`releasing.md`](../maintainer/releasing.md) §3 Punkt 4
(§7-Roadmap-Nacharbeit), Eigenes Plan-Risiko 5 (dieser Slice)
Befund: Die Tabelle *Nächste Wellen* trägt die Restthemen (Quellen, Optionen,
Basis-Image, [`slice-v1-docker-cli-bau-auf-standard-basis-image`](../plan/planning/open/slice-v1-docker-cli-bau-auf-standard-basis-image.md) und
[`slice-v1-katalog-features-quellen-und-freigabe`](../plan/planning/open/slice-v1-katalog-features-quellen-und-freigabe.md), beide weiter in `open/`) unter
dem Label `v0.20.0` — das Label war die Zuordnung des Nachzugs `b32e9c38` nach
dem `v0.19.0`-F-7-Bericht. Der Tag `v0.20.0` ist jetzt gesetzt und trägt allein
den Katalog-Schlüssel `java-toolchain`; die Restthemen stehen ohne Nummer. Der
Plan selbst (§3, Roadmap-Zeile) nennt genau diese Berichtigung als Auftragspunkt
(„Zustandsmarker („Vorschlag Release"-Zeile der `Nächste Wellen`-Tabelle) nach
dem Tag auf den Ist-Stand (Muster-Marker-Klasse, vier Auftritte)“) — im Range
bis zum Tag-Commit trägt kein Diff sie. Soll–Ist: das Label ist vergeben, der
Eintrag braucht die Folgerversion (oder die Entzeitlichung).
Risiko: Der nächste Planner/Cut liest `v0.20.0` als Zielnummer der Restthemen
und plant einen Schnitt auf eine bereits veröffentlichte Version; dieselbe
Klasse löste in `v0.18.0` (F-8/F-9) und `v0.19.0` (F-7) je eine Nachzugs-Berichtigungsrunde aus. Kein Gate prüft Roadmap-Labels gegen gesetzte Tags.
Verifizierbar: Review-only

**F-3 — LOW** `docs/plan/planning/open/README.md` (Bestand-Zeile, Zeile 102) ·
"Push und Tag nur nach ausdrücklicher Freigabe des Projektinhabers; nur Plan" —
Abgelaufener Zustandsmarker „nur Plan“ nach dem Cut
Klasse: Zustandsmarker nach dem Tag nicht auf den Ist-Stand
Quelle: [`AGENTS.md`](../../AGENTS.md) §3.10, Plan §6 Risiko 5 (dieser Slice)
Befund: Die Bestand-Zeile (angelegt im Plan-Commit `3bc9757f` — die M-F-5-Lektion
des `v0.19.0`-Cuts ist damit ehrlich getragen) trägt am Ende „nur Plan“; am
Tag-Commit ist der Schnitt geschnitten, getaggt und gepubliziert. Das Muster der
`v0.19.0`-Zeile („gesetzt …, `publish` und Job `tap` grün — geliefert …“) zeigt
die erwartete Ist-Fassung; sie gehört in die Berichtigung vor dem Abhaken des
DoD-Punkts „Bestand-Zeilen“.
Risiko: Der DoD-Punkt „Bestand-Zeilen“ wird gegen einen abgelaufenen Stand
geprüft; der nächste Abgleich (AGENTS §3.10, nach der letzten Fix-Runde) findet
einen Eintrag, der den Ist-Stand nicht trägt.
Verifizierbar: Review-only

**F-4 — LOW** `CHANGELOG.md` (`[0.20.0]`-Lead) · "An unreadable or
not-covering checksum source, and the combination of the external `java`
feature with `java-toolchain`" — Lead-Abweisungsliste enger als die
Vertrags-Abweisung
Klasse: Lead-Abweisungsliste enger als die Vertrags-Abweisung (2. Befund nach
`v0.19.0` F-3); Schwester der Lead-Grenzlisten-Klasse (Schwelle, siehe
Steering-Loop)
Quelle: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Zeile 427,
dritte Abweisung des Typs), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv)
(Abschnitt *Abweisungen*)
Befund: Der Vertrag kennt für den Eintragstyp drei Exit-10-Abweisungen: Quelle
nicht lesbar/nicht überdeckend, Eintrag des Typs **ohne Quellen-Angabe**, und
externes `java` plus Schlüssel. Der Lead nennt die erste und die dritte; der
dritte Fall („ Eintrag des Typs ohne Angabe seiner offiziellen Prüfsummen-Quelle
weist eine `version:`-Angabe ab “) fehlt. Der Träger-Eintrag trägt alle drei
(„an entry of the type without a source“) — nur der Lead kürzt gegenüber der
Vertragsliste und gegenüber dem `[0.19.0]`-Lead, der alle drei nannte
(„an entry without a source“).
Risiko: Ein Nutzer, der einem Eintrag des Typs ohne Quellen-Angabe eine
`version:` setzt, liest aus dem Lead die Annahme; die Abweisung steht nur im
Eintrag. Identisches Muster wie `v0.19.0` F-3 (dort fehlte „nicht überdeckend“).
Verifizierbar: Review-only — kein Gate liest Lead gegen Vertragsliste

**F-5 — INFO** `docs/plan/planning/in-progress/roadmap.md` (Zeile 26) ·
"liegt in `in-progress/` — Minor-Release" — Post-Tag-Roadmap-Aufträge offen
(Verweis Verifier)
Klasse: Zustandsmarker mit Ablaufpunkt (noch nicht abgelaufen; Fälligkeit
Closure)
Quelle: [`releasing.md`](../maintainer/releasing.md) §3 Punkt 4 und §7;
Plan §2 DoD „Doku nachgezogen“
Befund: Die Offene-Wellen-Zeile ist am Tag-Commit zutreffend (der Slice liegt da;
Zustand und Beleg, keine Chronik — Form korrekt). Mit dem Tag und der Closure
fallen nach `releasing.md` §3 Punkt 4/§7 die *Abgeschlossene-Wellen*-Zeile mit
Tag-Hash (`90f2a1d2`) und Datum (2026-11-03) sowie der Umzug des Slices nach
`done/` an; beide sind im Plan (§3) getragen und offen. Der Abgleich
Plan-Tabelle gegen `git diff` ist nach der letzten Fix-Runde zu wiederholen
(AGENTS §3.10).
Risiko: — (kein Verstoß am Prüfstand; Fälligkeit benannt)
Verifizierbar: Verifier (DoD), nicht dieser Lauf

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Lead des Blocks `[0.20.0]` Satz für Satz gegen Lastenheft `0.31.0` (Zeilen 416, 427, 1631) und [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv): dritter Schlüssel, Temurin LTS 21, Prüfsummen je Architektur vor dem Entpacken, Abbruch vor dem Extrahieren, Übersteuerung samt Lesung **bei der Erzeugung** (`the generation then reads … the build runs against them` — die `v0.19.0`-F-2-Lesung/Verifikation-Verwechslung ist berichtigt), Exit `10` ohne Schreiben, externes `java` allein wählbar, `JAVA_HOME`/`PATH` im Block ohne `containerEnv`, Homebrew-Befund ohne zweite Installationsform, Promised/Not-promised gegen die Vertrags-Liste | geprüft, ohne Befund abgesehen von F-4 (Abweisungs-Enumeration in der Lead-Summe) |
| `[Unreleased]`-Wanderung: Block-Inhalt (`### Added` + java-toolchain-Eintrag, 1 Abschnitt, kein *Changed*) gegen `git show 3bc9757f:CHANGELOG.md` zeichenweise verglichen — nach Header, Lead und Compare-Links byte-identisch; `[0.19.0]`-Block und ältere Blöcke unverändert; nur Plan-Aussage von §3 Punkt 1 erfüllt | geprüft, ohne Befund |
| Compare-Links und leerer Anker: `## [Unreleased]` ohne Inhalt, `[Unreleased]: …/compare/v0.20.0...HEAD`, `[0.20.0]: …/compare/v0.19.0...v0.20.0`, `[0.19.0]` unverändert; Block-Datum 2026-11-03 = Tag-Datum (M0) = Handbuch-`Stand` | geprüft, ohne Befund |
| Versionszeilen `0.20.0-dev` an den drei Stellen samt Kommentar-Nennungen — am Tag-Commit gelesen: `cmd/uboot/main.go` Zeile 45 (`var version = "0.20.0-dev"`), `repo.mk` Zeile 63 (Kommentar) und 66 (Kommentar) und 67 (`VERSION ?=`), `Dockerfile` Zeilen 47/49/51; `grep -rIn '0\.19\.0-dev'` über den Baum außer Historie/Reviews: 0 Treffer (nur die Plan-Zeile [`slice-v1-release-cut-v0.20.0.md`](../plan/planning/in-progress/slice-v1-release-cut-v0.20.0.md) §3 nennt den Übergang) | geprüft, ohne Befund |
| §3.7/§3.9 über den Cut-Diff: die berührten Kommentare beschreiben den Ist-Stand (Fallback-/Default-Formulierung), kein „früher/originally/no longer“, kein Verweis auf den alten Wert, kein Helfer zwischen Doc-Kommentar und Funktion | geprüft, ohne Befund |
| README-Regeln (§3 Punkt 3, AGENTS §5 Regel 6): Status-Zeile je README auf `v0.20.0`, `VERSION=0.20.0`-Beispiele, `grep` nach `LH-`/`SPEC-`/`ADR-`/`slice-`/„released“/Release-Tabelle: 0 Treffer | geprüft, ohne Befund |
| Handbuch: `Software-Version: v0.20.0`, `Stand: 2026-11-03` = Tag-Datum, die drei Versionsbeispiele an den Fundstellen des Plans (Zeilen 84, 110, 136) und `Gültigkeitsbereich` (Zeile 3240) auf `v0.20.0`; `grep '0\.19\.0'` im Nutzertext: 0 Treffer; Kennungen: 0 | geprüft, ohne Befund |
| „Twentieth release.“: Zählung folgt dem etablierten Minor-Ordinal der bisherigen Leads (`Nineteenth release.` bei `0.19.0`, dort negativfestgestellte fortlaufende Zählung) | geprüft, ohne Befund |
| Move-Disziplin (AGENTS §3.3 Fall 1): `a25484a8` reiner Move (R100), Verweis-Nachzug `70262862` eigener Commit; 1 eingehender Verweis nachgezogen; Plan-Kopf-Pfade lösen am Ruheort | geprüft, ohne Befund |
| Plan §3 gegen `git diff --name-status 3bc9757f..90f2a1d2` (10 Einträge): `CHANGELOG.md`, drei Versionsstellen, beide READMEs, Handbuch, Roadmap, `open/README.md`, Slice-Move — jede berührte Datei hat eine Zeile; die „kein Diff“-Zeile (`spec/`, `docs/plan/adr/`, `templates/`, `.github/`, `scripts/`, `packaging/`, `Makefile`, `go.mod`) bestätigt sich; offen am Prüfstand nur die nach dem Tag fälligen Plan-Aufträge (F-2, F-5) | geprüft (F-2, F-5) |
| Plan-Form gegen [`slice-v1-release-cut-v0.19.0`](../plan/planning/done/slice-v1-release-cut-v0.19.0.md): DoD-Formulierungen, Plan-Zeilen-Reihenfolge und M-Folge M0/M-T/M1/M2/M4/M5/M6/M7 übereinstimmend; erweiternd (kein Bruch): Lead-Grenz-Aufzählung in der CHANGELOG-Zeile, `templates/` neu in der „kein Diff“-Zeile, `open/README`-Zeile mit der M-F-5-Lektion, Risiko 5 (Marker-Klasse) und Fundstellen-Konkretisierung | geprüft, ohne Befund |
| Commit-Messages im Range: jede trägt [Slice-Kennung](../plan/planning/in-progress/slice-v1-release-cut-v0.20.0.md) samt `LH-*`ID (`90f2a1d2`: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); `ad4e93e1`: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)); die Message von `3bc9757f` („Bestand-Zeile im open/README mit dem Plan-Commit“, M-F-5-Lehre) entspricht ihrem Diff | geprüft, ohne Befund |
| Roadmap-Start-Zeile (`ad4e93e1`): Zustand wahr zum Commit-Moment (Move lag davor), Beleg als Verweis, kein Drift-Log-Eintrag für eine Schließung (§3.7) | geprüft, ohne Befund |
| Release-Kontrolle §6 (eigenständig gelesen, nicht übernommen): `gh release view v0.20.0` — 12 Assets (6 Binaries, 2 `.deb`, 2 `.rpm`, `SHA256SUMS`, `u-boot.rb`), kein Draft, kein Prerelease, Ziel-Commit `90f2a1d2`; CI-Läufe des Push-Tips `ci` 41581784136 / `integration-docker` 41581784179 / `tap-check` 41581828207 je success (eigener `gh`-Zugriff); `publish` 41583702257 success mit vier success-Jobs, Job-`tap`-Log trägt `tap-nachzug: pt9912/homebrew-u-boot updated to v0.20.0.` (kein Warn-Zweig); `docker run --rm ghcr.io/pt9912/u-boot:0.20.0 --version` = `u-boot version 0.20.0` (Probe von mir geführt); `:latest`-Digest `sha256:370bdfbe1b993bafebe5994c163df2826a9b323e170012f221b660ffbef058b5` = `:0.20.0` (`buildx imagetools inspect`, beide Tags) | geprüft, ohne Befund |
| Vorab `make gates`/`make ci`/`make fullbuild` (Exit 0, Stand `90f2a1d2`), Tap-Trockenlauf in `u-boot:test` mit `--network none` (fill+check Exit 0, Gegenprobe `v0.19.0` Exit 1, Nachzug-Protokoll und Gegenprobe), `make build VERSION=0.20.0`-Probe, M7-Mac-Bestätigung des Projektinhabers | Auftragskontext — als übernommene Evidence nicht eigenständig verifiziert (Reviewer-Gratwanderung: nur Lesezugriffe; Nachweis beim Verifier, der sie laut Plan §2 eigentlich führt) |
| Reviewer-HIGH-Liste (Dual-Classifier, Referenzmodell-Rang 1–3, Sprachvertrag, Architekturgrenze, Gate-Lockerung, Kommentar-Klassen, Template-Normen, Dateisicherheit, Carveout, Immutabilität, Lifecycle, Baseline-Bump): der Cut ändert keinen Produktcode, kein Gate, keinen ADR, keine Spec; `CHANGELOG`-Links auf Lastenheft-Anker laufen aufwärts; Lifecycle: ein Slice in genau einem Verzeichnis | geprüft, ohne Befund |

## Sensors

| Sensor | Ergebnis |
|---|---|
| `make docs-check` (Docker-Modus, eigener Lauf am Prüfstand) | Exit 0 — 1053 Dateien geprüft, 0 Befunde |
| `gh run list` / `gh run view` (Push-Tip `90f2a1d2`: `ci` 41581784136, `integration-docker` 41581784179, `tap-check` 41581828207; `publish` 41583702257 samt Job-`tap`-Log) | je success, log-Zeile gelesen (siehe Negativbefunde) |
| Release- und GHCR-Proben (`gh release view`, `docker run … :0.20.0 --version`, `buildx imagetools inspect` beider Tags) | wie angegeben |
| `make gates`, `make ci`, `make fullbuild`, `make build VERSION=0.20.0`, Tap-Trockenlauf, `make test-docker` | nicht durch diesen Lauf — Auftragskontext (Vorab-Nachweis des Implementers); `make test-docker`: im Plan als Nicht-Ziel geführt (`integration-docker` am Push-Tip ist der Nachweis); Wiederholung beim Verifier |
| `make verify-depguard` | nicht ausgeführt — keine Importregel berührt (kein Produktcode im Diff) |
| `make ablage-pruefung`, `make baseline-verify` | in `make gates` enthalten (Implementer-Lauf); hier nicht wiederholt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (F-1, F-2) |
| LOW | 2 (F-3, F-4) |
| INFO | 1 (F-5) |

**Finding-Klassen dieses Laufs:** Zustandsmarker nach Ablauf nicht abgebaut
(dritter Lauf, identischer §7-Wortlaut — F-1) · Vorschlags-Version vergeben
(dritter Befund, `v0.18.0` F-8/F-9, `v0.19.0` F-7 — F-2) ·
Zustandsmarker nach dem Tag nicht auf den Ist-Stand (F-3, die Klasse des
Plan-Risikos 5) · Lead-Abweisungsliste enger als die Vertrags-Abweisung
(2. Befund nach `v0.19.0` F-3 — F-4).

**Steering-Loop-Vermerk:** „Zustandsmarker nach Ablauf nicht abgebaut“ steht mit
F-1 bei drei Läufen mit identischem §7-Wortlaut (`v0.18.0` F-4, `v0.19.0` F-4,
dieser Lauf F-1) — die Schwelle ist auch mit dem `v0.19.0`-Vermerk rechnerisch
überschritten; Kandidaten: ein computationaler Prüflauf (d-check-Heuristik oder
Closure-Checklisten-Schritt „Statusmarker gegen Lifecycle-Verzeichnis“) und die
Klassen-Zusammenlegung mit „Vorschlags-Version vergeben“ (F-2) unter der
Marker-Linse in `harness/review.md` und dem Skill. „Lead-Grenzlisten-Klasse“
(Umbrella, seit `v0.19.0` auf der Schwelle) tritt mit F-4 erneut auf — der
Architect-Auftrag aus der `v0.19.0`-Closure (Pflichtschritt „Lead gegen die
Vertragsliste lesen“ vor dem Tag) bleibt offen und wird zur Schwelle geführt.
Beide Setzungen liegen beim Projektinhaber; dieser Lauf verkörpert nichts.

## Verdikt

**Merge-blockierend: nein** — der Tag `v0.20.0` ist gesetzt, gepubliziert und
wird nicht zurückgenommen ([`releasing.md`](../maintainer/releasing.md) *Zweck*).
F-1 bis F-3 betreffen Roadmap-, Bestand- und Slice-Marker und gehen als
Berichtigungen der Fix-Runde vor der Closure an den Implementer; F-4 (Lead-Abweisung
in `[0.20.0]`) ist ein Nachzugs-Kandidat auf `main` (Muster `v0.19.0` F-1–F-3:
Entscheidung des Projektinhabers, Nachzug vor der Closure oder mit dem nächsten
Patch). Der Slice geht nicht nach `done/`, bevor F-1–F-3 berichtigt sind; der
DoD-Abgleich (AGENTS §3.10) nach der letzten Fix-Runde deckt die Zeilen
„Bestand-Zeilen“ (F-3) und „Roadmap“ (F-2, F-5) gegen den Ist-Stand. HIGH gibt
es nicht — ein Konfliktpfad über den Architect entfällt; die MEDIUM-Findings
stehen in der Schwelle-Klasse des Steering-Loop (Lern-Eintrag in §7 der
Closure).

**DoD-Vermerk:** Nicht Bestandteil dieses Laufs (Verifier). Festgehalten: der
DoD-Punkt „Vorab-Prüfungen“ trägt die CI- und Release-Evidence (oben eigenständig
nachgelesen); die offenen Punkte (Review-Report — dieser Lauf —,
Verification-Report, Closure-Notiz mit Lerneintrag, Roadmap-Abgeschlossenen-Zeile
mit Tag-Hash, Marker-Berichtigungen F-1–F-3, `git mv` nach `done/`) sind im
Slice ehrlich offen.