# Review-Report: slice-cr-devcontainer-java-toolchain — 2026-11-02

**Review-Art:** Code-Review des CR-Diffs — geprüft wird die
Wortlaut-Lieferung (Lastenheft 0.31.0, TECH-Spec, Planning-Nachzug)
gegen den freigegebenen
[CR-Vorschlag](../plan/planning/done/cr-vorschlag-gepinnte-werkzeugkette.md)
der Werkzeugkette bzw. den freigegebenen
[Java-CR-Vorschlag](../plan/planning/done/cr-vorschlag-java-toolchain.md),
die Hard Rules
([`AGENTS.md`](../../AGENTS.md) §3) und die Rollenentscheidungen F1 bis
F8. Keine DoD-Closure — die prüft der Verifier
([`harness/roles.md`](../../harness/roles.md) §Kontext-Trennung).

**Gegenstand:** die drei Contract-Commits `1ddada3e` (Lastenheft 0.31.0),
`c1f14f30` (TECH-Spec samt Anker-Nachzug), `5da10dba` (Planning-Nachzug)
gegen den Slice
[`slice-cr-devcontainer-java-toolchain`](../plan/planning/done/slice-cr-devcontainer-java-toolchain.md);
Grundlage des §3.10-Abgleichs: `6731bd52..5da10dba` = `HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ Commit `5da10dba` ·
**Modell:** `glm-5.3-flash:cloud[1m]` · **Datum:** 2026-11-02

**Eingangs-Kontext:**

- [Slice-Plan](../plan/planning/done/slice-cr-devcontainer-java-toolchain.md) samt §7-Abgleich des Implementers
- [CR-Vorschlag](../plan/planning/done/cr-vorschlag-java-toolchain.md) (Freigabe-Block, §§2 Stellen 1 bis 12, F1 bis F8)
- [`spec/lastenheft.md`](../../spec/lastenheft.md) ([`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [`LH-FA-DEV-011`](../../spec/lastenheft.md#lh-fa-dev-011--pin-politik-der-gepinnten-werkzeugkette), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)), [`spec/spezifikation.md`](../../spec/spezifikation.md) ([`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion), [`SPEC-023`](../../spec/spezifikation.md#spec-023--devcontainer-basis-image-und-feature-quellen), [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config), [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-050`](../../spec/spezifikation.md#spec-050--katalog-features-im-sandbox-profil-abweisung-von-docker-cli-lockerungen-von-go-features-aus-externen-quellen), [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv))
- [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) (`Accepted`, nur Referenzebene)
- Beobachtungen: [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md), [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules §3, §5 Regel 6), [`harness/conventions.md`](../../harness/conventions.md), [`harness/review.md`](../../harness/review.md)
- Umsetzungs-Slice [`slice-v1-katalog-feature-jdk-gepinnt`](../plan/planning/in-progress/slice-v1-katalog-feature-jdk-gepinnt.md) (Konsument der offenen Punkte)

---

## Findings

Schema nach `.harness/skills/reviewer.md` §Output-Schema
([`harness/review.md`](../../harness/review.md) ist die kanonische Prosa).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der §3.10-Abgleich (Slice §7, berichtigt auf den Stand nach den Commits) behauptet die „Sicht-Zeilen“ der §3-Tabelle „[bleiben] im Diff unberührt“, während die Sicht vom Anker-Nachzug berührt ist (zwei reine Anker-Zeilen in [`spec/architecture.md`](../../spec/architecture.md)); die §3-Tabelle führt die Datei unter „ungeöffnet / unverändert“. Der Abgleich ist der beweisende Sensor dieser Regel — seine Aussage trägt nicht. | [`AGENTS.md`](../../AGENTS.md) §3.10; [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) | [`spec/architecture.md`](../../spec/architecture.md) · „ungeöffnet / unverändert“ in der §3-Tabelle (Plan-Zeile 7); „bleiben im Diff unberührt“ in Slice §7 | Review-only (`git diff-tree -r c1f14f30`: die zwei Zeilen sind Anker-Berichtigungen) | Plan-Abgleich lässt berührte Datei unter den Tisch fallen |
| F-2 | MEDIUM | Für `harness/replay.md` (Anker-Zeile des Anker-Nachzugs) gibt es weder eine §3-Plan-Zeile noch eine §7-Erwähnung; der Abgleich klassifiziert die 17 Nachzug-Dateien als „planning- und report-Dateien“, was weder [`spec/architecture.md`](../../spec/architecture.md) (Sicht-Stratum) noch `harness/replay.md` (Replay-Vertrag, [`harness/replay.md`](../../harness/replay.md)) trägt — der eigene Commit-Text listet beide korrekt. | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`harness/replay.md`](../../harness/replay.md) · „planning- und report-Dateien … (17 Dateien, reine Querverweis-Berichtigungen)“ in Slice §7 | Review-only (`git diff-tree -r c1f14f30`) | Plan-Abgleich lässt berührte Datei unter den Tisch fallen |
| F-3 | LOW | Der CR-Vorschlag steht in keiner §3-Plan-Zeile: er ist im Implementer-Range neu angelegt (Planner-Commit) und zweimal geändert (Anker, Freigabe-Block) — der Vorgänger-CR führte ihn als Zeile der Tabelle. | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`cr-vorschlag-java-toolchain.md`](../plan/planning/done/cr-vorschlag-java-toolchain.md) · fehlt in der §3-Tabelle (Vorgänger: Plan-Zeile mit „CR-Vorschlag“) | Review-only | Plan-Abgleich lässt berührte Datei unter den Tisch fallen |
| F-4 | LOW | Der versprochene Nachzug in den Umsetzungs-Slice (Plan-Zeile „Umsetzungs-Slice … (Berührte Spec-Stellen und §3-Tabelle nachziehen)“) ist über den Anker-Nachzug nicht hinausgekommen: dessen „Berührte Spec-Stellen“ nennen [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion) weiter nur „Hosts-Zeile `feature:java`“, und seine §3-Zeile plant „die Schlüssel-Namen in [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) (Zeilen 522 und 544)“, eine [`SPEC-019`](../../spec/spezifikation.md#spec-019--default-allowlist-der-egress-restriktion)-Hosts-Zeile und eine eigene Stelle neben [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) — alles seit dem CR-Diff geliefert oder entschieden (F6 erweitern, F8 im CR-Commit). | [`AGENTS.md`](../../AGENTS.md) §3.10; Plan-Zeile des CR-Slice | [`slice-v1-katalog-feature-jdk-gepinnt.md`](../plan/planning/in-progress/slice-v1-katalog-feature-jdk-gepinnt.md) · „(Hosts-Zeile `feature:java`, gelesen)“ | ja — `make docs-check` trägt die Anker, der Inhalt ist Plan-Text (Berichtigung beim nächsten Implementer-Lauf) | Plan des Folgeslices zum Ist-Stand nicht nachgezogen |
| F-5 | LOW | Der Slice ist per `git mv` in `in-progress/` gestartet, `Verantwortlich:` steht aber noch auf „—“ — der eigene §4-Start-Trigger setzt ein gesetztes Feld voraus, und beide Vorgänger-CR füllten es beim Start. | Slice §4 (Start-Trigger); [`harness/roles.md`](../../harness/roles.md) | [`slice-cr-devcontainer-java-toolchain.md`](../plan/planning/done/slice-cr-devcontainer-java-toolchain.md) · „Verantwortlich: — (Beim Start gesetzt.)“ | Review-only | Start-Trigger-Zusage nicht erfüllt |
| F-6 | INFO | Lastenheft 0.31.0 und [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) tragen den Schlüssel als Zielbild — Code und Tests tragen ihn noch nicht (Katalog-Code ohne Java-Konstante, die neun Java-Zeilen der Beleg-Tabelle sind Aufträge, die Konflikt-Meldung liegt in [`devcontainer_toolchain.go`](../../internal/hexagon/application/devcontainer_toolchain.go) nur als Go-Muster vor). Die Reihenfolge ist die abgesprochene des Musters, aber der Release-Cut darf von diesem HEAD den Schlüssel nicht als implementiert schneiden. | Maintainability (Zielbild-Muster des Vorgänger-CR) | [`spec/lastenheft.md`](../../spec/lastenheft.md) java-Absatz · „eine abweichende Prüfsumme bricht den Bau vor dem Entpacken ab“ | nein — DoD-Zeile 4 und Umsetzungs-Slice (Verifier) | Vertragstext vor Umsetzung |
| F-7 | INFO | Die Historie-Zeile 0.31.0 sagt „die Zahlen stehen im `CHANGELOG` und in dieser Historie als Chronik“ — auf diesem Stand trägt keine der beiden Stellen die Temurin-Revision (`CHANGELOG` `[Unreleased]` ist leer, die Zeile trägt keine Zahl); der §7-Offene-Punkt schiebt sie an den Umsetzungs-Slice, während das Vorgänger-Muster F4 die Zahlen direkt in die Zeile legte. | [`AGENTS.md`](../../AGENTS.md) §3.7 (Zustandsaussage; Ablaufpunkt) | [`spec/lastenheft.md`](../../spec/lastenheft.md) Historie 0.31.0 · „die Zahlen stehen im `CHANGELOG` und in dieser Historie als Chronik“ | Review-only | Zustands-Aussage über Ablauf-Zustand |
| F-8 | INFO | Die gelesene Vertriebs-Lage hinter der F8-Hosts-Zeile (Fundstelle im Slice §7: Release-Listing `adoptium/temurin21-binaries`, Redirect eines Prüfsummen-Assets, Datum) kann der Review nicht wiederholen (netzfremder Lauf) — die Fundstelle ist mit Datum, Listing und Asset-Host dokumentiert, und der Docker-Test des Umsetzungs-Slices pinnt die Form; eine abweichende Messung berichtigt die Zeile (auch in der Spec so gesetzt). | [`AGENTS.md`](../../AGENTS.md) §3.10 | [`spec/spezifikation.md`](../../spec/spezifikation.md) Abschnitt „Zeile `java-toolchain`“ | ja — `make test-docker` am Umsetzungs-Diff | Netz-Lese des Reviewers nicht wiederholbar |

### Failure-Szenarien

F-1/F-2: der nächste Änderer (Fix-Runde, Verifier, Baseline-Bump) liest die
§3-Tabelle als Inventur — [`spec/architecture.md`](../../spec/architecture.md)
steht dort unter „unverändert“, `harness/replay.md` fehlt gänzlich; der
§7-Abgleich, der die Vollständigkeit behauptet, ist genau der Beleg, den die
Closure der DoD entnimmt — dieselbe Blende, die die Verkörperung des §3.10
in zehn Vorgängen trug ([`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md),
zwei Auftritte nach der Verkörperung; der Vorgänger-CR-Review trug dieselbe
Klasse als Plan-Berichtigung nach). F-3: der Verifier des Closures prüft die
§3-Tabelle gegen den Range und liest eine Dateileiste, die eine real
berührte Datei fehlt — die Vorlage des Vorgänger-CR trug den Eintrag. F-4:
der Start des Umsetzungs-Slices liest eine §3-Zeile, die das Hinzufügen der
Schlüssel-Namen-Zeilen und der Hosts-Zeile plant, und legt sie erneut an oder
bemängelt den Zustand — die Arbeit des CR-Diff erscheint doppelt oder der
Slice-Start bricht am Widerspruch zum Lastenheft 0.31.0. F-5: die Closure
findet ein Übergabe-Feld leer und trägt die Verantwortlichkeit rückwirkend
ein — dieselbe Buchführungslücke, die [`harness/roles.md`](../../harness/roles.md)
bei einem Handoff ohne beschriftbare Übergabe sieht.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Wortlaut-Treue java-Absatz ([`spec/lastenheft.md`](../../spec/lastenheft.md) gegen den Block-Wortlaut des [CR-Vorschlags](../plan/planning/done/cr-vorschlag-java-toolchain.md), Stelle 2) | geprüft, ohne Befund — Programmatic-Abgleich (Marker „(J1)“ bis „(J11)“, „(Stelle 4)“ und die Datei-Anker-Form entfernt): Zeichen für Zeichen gleich, je 2768 Zeichen |
| Listen-Zeile (Position zwischen den zwei Toolchain-Zeilen und der `ollama-host`-Zeile; Wortlaut der Stelle 1) | geprüft, ohne Befund (Zeilen 414 bis 417) |
| Abweichungs-Versprechen (genau drei deklarierte Abweichungen, keine weitere) | geprüft, ohne Befund |
| Decken-Regel im Lastenheft (Suchlauf nach `ADR-*`, `SPEC-*`, `slice-*`, `BEO-*`, Planning-Kennungen über die drei neuen Stellen und die Historie-Zeile) | geprüft, ohne Befund — nur [`LH-*`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)-Anker innerhalb der Datei |
| F7 — kein `sdkman`-Fall in den Abweisungen (grep über beide Spec-Dateien) | geprüft, ohne Befund — `sdkman` nur in den bestandseigenen Zeilen ([`SPEC-019`](../../spec/spezifikation.md), [`SPEC-023`](../../spec/spezifikation.md)) |
| F8 — Hosts-Zeile ohne erfundenen Host; `api.adoptium.net` in der Tabelle nicht und im Absatz ausdrücklich negiert; Same-Origin begründet | geprüft, ohne Befund (Netz-Lese des Reviewers nicht wiederholbar — siehe F-8) |
| Ist-Zustand-Regel — `docs/user/benutzerhandbuch.md`, beide README-Dateien, `CHANGELOG.md` ohne Diff im Range | geprüft, ohne Befund (`git diff --stat 6731bd52..HEAD` führt keine dieser Dateien); Plan-Zeile der Stellen 11/12 auf Ist-Zustand berichtigt |
| [`SPEC-040`](../../spec/spezifikation.md) — die zwei Zeilen (522, 544), nur Klammer trägt den dritten Namen, Zahlen unverändert | geprüft, ohne Befund |
| [`SPEC-019`](../../spec/spezifikation.md) — Tabellen-Zeile und Absatz „Zeile `java-toolchain`“ | geprüft, ohne Befund |
| [`SPEC-054`](../../spec/spezifikation.md#spec-054--katalog-eintragstyp-gepinnte-werkzeugkette-go-toolchain-golangci-lint-und-java-toolchain-prüfsummen-quellen-version-übersteuerung-und-cache-werte-unter-containerenv) — Titel, Absätzen, Quellen-Tabelle (Same-Origin, `api.adoptium.net` negiert), Abweisungen, Cache-Abgrenzung, Meldungen, Beleg-Tabelle (Auftrags-Form), Historie | geprüft, ohne Befund |
| [`SPEC-050`](../../spec/spezifikation.md) und [`SPEC-023`](../../spec/spezifikation.md) — keine Berührung, wie im Commit gemeldet | geprüft, ohne Befund (die Zeilen sind im Diff nicht berührt) |
| Anker-Nachzug — die 17 Dateien je reine Querverweis-Berichtigung (alle berührten Zeilen je gegen den Anker-Muster geprüft, keine inhaltliche Änderung versteckt) | geprüft, ohne Befund |
| Fundstellen der Vorschlags-Belege gegen Code (`devcontainer_features.go` 173 bis 182, 206 bis 221, 273 bis 277, 338 bis 347; `devcontainer_toolchain.go` 159 bis 166, 168 bis 179, 181 bis 192) | geprüft, ohne Befund — die gelesenen Aussagen tragen am Fundort |
| Fundstellen des §7-Abgleichs (Lastenheft 9, 398 bis 437, 425/427/429, 523 bis 527, 529, 1627 bis 1631; Spec 973 ff. mit 991/993/995/997, 1464/1475/1483/1495, 1561, 634, 1609; [`SPEC-040`](../../spec/spezifikation.md#spec-040--konfigurationspfade-und-schreibpipeline-von-u-boot-config) 522/544, [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) 600) | geprüft, ohne Befund — alle Fundstellen zeigen auf das Gemeldete |
| Referenzrichtung in der Spec (nur aufwärts: [`LH-*`](../../spec/lastenheft.md) und Anker innerhalb der Datei; die Beleg-Tabelle nennt Tests, wie seit 0.30.0) | geprüft, ohne Befund |
| [[`LH-FA-DEV-001`](../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen)] — Cache-Zusage benennt die zwei Go-Schlüssel (kein Widerspruch zum neuen „kein Cache-Verzeichnis“) | geprüft, ohne Befund |
| §3.9-Muster über die Diffs (früher/entfernt-Sätze, umbenannte Symbol-Verweise, eingefügte Helfer vor Doc-Kommentar) | geprüft, ohne Befund |
| F-Tabelle F1 bis F8 im CR-Vorschlag auf entschieden Stand; Status-Zeile auf freigegeben | geprüft, ohne Befund |
| Historie-Form (Lastenheft §16 nach der 0.30.0-Zeile; Spec-Historie-Zeile 2026-11-02; Versionsfeld 0.31.0 konsistent) | geprüft, ohne Befund (die Chronik-Sätze selbst: siehe F-7) |
| [`ADR-0041`](../plan/adr/0041-gepinnte-werkzeugkette-hybrid-versionierung.md) — unverändert (`Accepted`) | geprüft, ohne Befund |

**Nicht geprüft:** `make gates` und die übrigen Sensors — Review-only
(aufgetragener Umfang, kein Gate-Lauf; der Implementer meldete
`make docs-check` 0 Befunde je Commit). Die letzte grüne CI-Läufe
(`ci`, `integration-docker`) stehen auf `51614375`; für die drei
Contract-Commits ist kein CI-Lauf sichtbar. Die `make gates`-Strecke ist
DoD-Zeile 4 und Verifier-Aufgabe.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Plan-Abgleich lässt berührte Datei
unter den Tisch fallen · Plan des Folgeslices zum Ist-Stand nicht
nachgezogen · Start-Trigger-Zusage nicht erfüllt · Vertragstext vor
Umsetzung · Zustands-Aussage über Ablauf-Zustand · Netz-Lese des
Reviewers nicht wiederholbar

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (der Abgleich des Slices behauptet
eine Vollständigkeit, die sein eigener Diff nicht trägt) blockieren die
DoD-Zeile 4; eine Berichtigungs-Runde über Planning-Text allein (kein
Lastenheft-Text) löst F-1 bis F-5, ein Nachlauf-Review über Vertragstext
entsteht dadurch nicht (die Bedingung des
[`BEO-DOK/vertragsaenderung-nach-review`](../plan/planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)
greift nicht).

**Übergabe:** Findings an den Implementer (Fix-Runde ohne
Lastenheft-Diff); die Finding-Klassen gehen in die Slice-Closure §7 und
von dort in den Beobachtungs-Zähler — F-1 bis F-3 zählen dieselbe Familie
wie [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md)
(dritter Auftritt dieser Klasse im CR-Strang). Dieser Report selbst ist
ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier separat
([`harness/verification.md`](../../harness/verification.md)).