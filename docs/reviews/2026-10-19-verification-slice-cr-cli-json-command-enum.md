# Verification-Report: slice-cr-cli-json-command-enum — 2026-10-19 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-cr-cli-json-command-enum` (`docs/plan/planning/in-progress/`), Lastenheft 0.15.0 auf 0.16.0, Commit-Range `d20be6d..HEAD` (6 Commits, HEAD `2e9660d`, 12 Dateien). Davon zwei reine Move-/Nachzug-Commits (`62aab4e`, `8d837b7`) und vier Inhalts-Commits (`42fc809` Umsetzung, `84e2239` Review-Report, `903c73f` und `2e9660d` Review-Fix). Kein Produktcode, keine Spezifikation, keine Architektur-Sicht, keine ADR im Diff. <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-19-slice-cr-cli-json-command-enum.md` (F-1 MEDIUM, F-2 bis F-4 LOW, F-5, F-6 INFO), DoD-Haken des Slice (zwei gesetzt, übrige offen), Slice §1 (Wortlaut-Vorschlag), §3 (Plan-Tabelle). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Das Lastenheft steht wortgleich zum Vorschlag im Slice (bis auf die benannte Link-Form), der Lastenheft-Diff hat genau drei Hunks, `make gates` und `make doc-trace` sind grün, alle Review-Findings F-1 bis F-4 sind behoben, F-5 und F-6 sind übergeben beziehungsweise ausgewiesen. Es gibt keine DoD-Verletzung und keine Abweichung vom Plan-Wortlaut. Zwei kleine Nachträge zum Plan (V-2, V-3) und die Closure-Arbeit liegen beim Planner. Das Release-Gate (Bestätigung des Projektinhabers) ist offen und kein Mangel. Der Slice ist bis zur Closure-Notiz nicht schließbar. <!-- d-check:ignore (Kennungen als Zitat) -->

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | Exit-Code des `make`-Prozesses selbst in die Log-Datei geschrieben (`EXIT=0`); Ende des Logs gelesen: `coverage-gate: OK — coverage 91.40% meets threshold 90%`; Lint-Stufe (`golangci-lint run ./...`) ohne Abbruch; HEAD `2e9660d`; Arbeitsbaum vorher und nachher sauber (`git status --short` leer) |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ (keine neue Waise gegenüber dem Stand des Ollama-CR) |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | kein Produktcode, keine Importregel, kein Docker-Pfad des Produkts im Diff (`git diff --stat` ohne `internal/`, `cmd/`, Dockerfile, Makefile) |
| `make doc-immutable` | nicht gelaufen | keine Datei unter `docs/plan/adr/` im Diff; nicht aufgetragen |

## 2. Lastenheft: Wortlaut-Konformität (Prüfpunkt 1)

Gelesen: Slice §1 „Wortlaut-Vorschlag“ gegen `git diff d20be6d..HEAD -- spec/lastenheft.md`.

| Stelle | Ergebnis |
|---|---|
| 1: Enum [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) | Absatz „Für `--dry-run --json` …“: `… generate`, `config`, `egress` oder `template`.` wortgleich zum Vorschlag „Neu“; Rest des Absatzes unverändert (Diff ersetzt genau eine Zeile) | <!-- d-check:ignore (Kennungen als Zitat) -->
| 2: Versionsfeld | `0.15.0` auf `0.16.0`, Spaltenbreite unverändert, wortgleich |
| 2: Historie-Zeile | Zeile `0.16.0`, Datum `2026-10-19` (Datum des Lastenheft-Commits, wie der Vorschlag es verlangt), nach `0.15.0`; Text wortgleich zum Vorschlag; einzige Abweichung die Link-Form (Anker der eigenen Datei statt Relativpfad, wie im Lastenheft üblich); Klausel „Bestätigung des Projektinhabers vor dem Release ausstehend“ vorhanden; Quelle „Vereinbarung mit dem Projektinhaber“ |
| Sonst im Lastenheft | genau **drei Hunks** (Versionsfeld, Enum-Absatz, Historie), 3 Einfügungen, 2 Löschungen; `LH-NFA-USE-004`, `LH-SA-CLI-001` und die Aufzählung der `subcommand`-Pflicht unverändert (laut Slice bewusst nicht zu ändern) | <!-- d-check:ignore (Kennungen als Zitat) -->

Keine Slice-, Wellen-, ADR-Kennung in den neuen Zeilen (`AGENTS.md` §3.4). Der Vergleich ist von Hand geführt (Lesen und `git diff`), nicht mit einem Zeichen-Diff-Werkzeug. Eine Freigabe des Wortlauts gibt es nicht und der Slice verlangt sie nicht (abgeleitet, Gate vor dem Release).

## 3. DoD-Abgleich (Slice §2), Punkt für Punkt

| DoD-Zeile | Haken im Slice | Stand | Beleg |
|---|---|---|---|
| Lastenheft 0.16.0 (Stellen 1 und 2, Klausel) | gesetzt | **erfüllt, der Haken trägt** | Abschnitt 2 |
| Vorab-Suche und Nachzug (Vermerk §6.10, „zehn Subcommands“ in §6) | gesetzt | **erfüllt, mit Einschränkung** | Gegen den Baum gelesen: §6.10 trägt jetzt „**Vertrag:** … `command: "egress"` liegt damit im Vertrag … Bestätigung des Projektinhabers steht vor dem Release aus“, der Vermerk „Offene Vertragslücke“ ist entfernt; §6 sagt „listet elf Werte“ und „Zehn davon (alle außer `egress`, siehe §6.10)“. Die Vorab-Suche des Implementers ist in keinem Artefakt festgehalten (Slice §7 leer); der Plan-Text §3 vermerkt nur den Satz „Der Implementer hat sie … wiederholt (nur über das Lastenheft …)“. Auflage 3 |
| `make gates` grün; `make doc-trace` ohne neue Waise | offen | **erfüllt** (Abschnitt 1), Haken noch nicht gesetzt | Closure-Schritt |
| Review durchgeführt, Report liegt vor | offen | **erfüllt**, Haken noch nicht gesetzt | Report vorhanden, lief über den vorgelegten Wortlaut mit HEAD `42fc809`; die Fix-Runde (`903c73f`, `2e9660d`) ändert danach keinen Lastenheft-Text (nur Handbuch-Tabelle, Contract-Satz, Slice-/Plan-Text), ein Nachlauf-Review ist für sie nicht nötig, ihre Verifikation steht in Abschnitt 5 |
| Doku-Update: Benutzerhandbuch (`command`, `subcommand`) und `cli-json-contract.md`, kein CHANGELOG | offen | **erfüllt**, Haken noch nicht gesetzt | Handbuch Zeilen 2149 und 2150 tragen `egress` (`list`); `cli-json-contract.md` §6, §6.10 nachgezogen; `CHANGELOG.md` nicht im Diff |
| Closure-Notiz mit Steering-Loop-Lerneintrag | offen | offen, zu Recht | §7 trägt noch den Platzhalter (Auflage 3) |
| Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ | offen | offen, zu Recht | `docs/plan/planning/observations/` nicht im Diff (Auflage 3) |
| Jedes Risiko aus §6 trägt einen Ausgang | offen | offen, zu Recht | §6 steht auf „weiter offen“ beziehungsweise „offen“ (Auflage 3) |
| Drei Paarungen (Anker · Folge-Slice · Register) | offen | offen, zu Recht | Welle-Closure prüft sie |

Entscheidungs-Konformität: keine Accepted-ADR im Diff (`AGENTS.md` §3.5), kein Gate gelockert (§3.6), kein Kommentar im Code berührt (§3.7 nur für Prosa-Zeilen gelesen: Indikativ, keine „früher“-Muster in den geänderten Zeilen).

## 4. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff --stat d20be6d..HEAD`, `AGENTS.md` §3.10)

12 Dateien im Diff (Stand `2e9660d`). Mit diesem Report: 13.

| Diff-Datei | Plan-Zeile | Befund |
|---|---|---|
| `spec/lastenheft.md` | Zeile 1 | ok (Abschnitt 2) |
| `docs/maintainer/cli-json-contract.md` | Zeile 2 | ok; zwei Hunks (§6, §6.10), je Aussage gegen Code gelesen (`egress.go`: `newDataEnvelope("egress", "list", …)`, `ErrEgressSubcommandRequired`, `ErrEgressPreviewNotApplicable`) |
| `done/slice-cr-cli-json-command-enum.md` | Zeile 3 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/reviews/2026-10-19-slice-cr-cli-json-command-enum.md` (neu) | Zeile 4 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/user/benutzerhandbuch.md` | Zeile 5 (durch den Review-Fix ergänzt) | ok |
| `in-progress/roadmap.md` | Zeile 6 (durch den Review-Fix ergänzt) | ok: Arbeitszeile „In Arbeit: …“ (Zustand und Anker, keine Chronik) und Verweis in der Welle-Zeile |
| `CHANGELOG.md` | „unverändert“ | bestätigt: nicht im Diff |
| `done/slice-cr-devcontainer-ollama.md`, `done/slice-v1-cli-egress-anzeige.md`, `open/README.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-verification-slice-cr-devcontainer-ollama.md` | **keine Zeile** | Diffs gelesen: reine Pfad-Nachzüge `open/` nach `in-progress/` (Commit `8d837b7`, Muster `slice-mv`); keine Inhaltsänderung; **V-2** | <!-- d-check:ignore (Kennungen als Zitat) -->
| `open/slice-v2-spec-sandbox-egress-ollama.md` | **keine Zeile** | zwei Anteile: Pfad-Nachzug (`8d837b7`) und eine **Inhaltsänderung** in `903c73f` (Plan-Zeile `spec/architecture.md` führt die Aufzählung „Implementiert sind …“ der Driving-Adapter); **V-2** | <!-- d-check:ignore (Kennungen als Zitat) -->

Kein Eintrag der Plan-Tabelle ohne Diff. Zwei Diff-Gruppen ohne Plan-Zeile (V-2). Der Review-Fix hat die Tabelle um Handbuch und Roadmap ergänzt (Review F-2 behoben), aber nicht um die Verweis-Nachzüge und die Übergabe im Spec-Slice; der Abgleich nach der letzten Fix-Runde (§3.10) ist damit nicht vollständig.

## 5. Review-Findings F-1..F-6 einzeln (Prüfpunkt 2)

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 MEDIUM Handbuch ohne `egress` im Wertebereich | **behoben** | `903c73f`: Handbuch Zeile 2149 `command` führt `egress` an der Stelle zwischen `config` und `template` (gleiche Reihenfolge wie das Lastenheft); Zeile 2150 `subcommand` „nur bei `config` (`get`, `set`, `list`), `egress` (`list`) und `template` (`list`)“. Gegen den Code: `egress.go:158` `newDataEnvelope("egress", "list", …)`, `template.go:213` `"template", "list"`; die übrigen Befehle setzen `""`. Weitere Fundstellen des Enums im Baum gesucht (`generate`, `config` in README, `README.de.md`, `docs/user`, `docs/maintainer`, `spec`): README-Zeilen sind Release-Historie `v0.1.0`; sonst nur Lastenheft, Handbuch und `architecture.md` (V-5) |
| F-2 LOW Plan-Tabelle ohne Roadmap | **behoben, mit Rest** | Zeilen Handbuch und Roadmap stehen in §3; Rest siehe V-2 |
| F-3 LOW Zustandsmarker (a) bis (c) im Slice | **behoben** | (a) „Der CR … belegt `0.15.0` (noch nicht gelandet)“ ersetzt durch Ist-Stand der Historie; (b) Stelle 2 „Alt: `0.15.0`“ / „Neu: `0.16.0`“, Ziel und Risiko 2 ohne „erwartet“ (`903c73f`, `2e9660d`); (c) Vorab-Suche „hat sie … wiederholt (nur über das Lastenheft …)“. Suche im Slice nach „erwartet“, „noch nicht gelandet“, „wiederholt“ ohne Rest. Ein Marker außerhalb des Slice steht noch: V-3 |
| F-4 LOW „Zehn davon“ ohne benannte Differenz | **behoben** | `cli-json-contract.md` §6: „Zehn davon (alle außer `egress`, siehe §6.10) migrieren **inkrementell** über neun Folge-Slices“ |
| F-5 INFO Kopien des Enums in `SPEC-001`, im Zitat von §3 und in `architecture.md` | **übergeben, getragen** | `SPEC-001` ist im Spec-Slice als Übergabe eingetragen; die Aufzählung „Implementiert sind …“ (`spec/architecture.md` Zeile 232) steht seit `903c73f` in der Plan-Zeile `spec/architecture.md` von `slice-v2-spec-sandbox-egress-ollama` §3 (Zeile 77). Die DoD-Zeile „Sicht und Historie“ (dort Zeile 60) nennt die Aufzählung nicht ausdrücklich (V-5). Das Zitat in `cli-json-contract.md` §3 bleibt bis zum Spec-Slice wortgleich zu `SPEC-001`, der Vermerk §6.10 nennt die Nachführung | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-6 INFO DoD-Haken ohne festgehaltenen Beleg der Vorab-Suche | **ausgewiesen, an die Closure übergeben** | Der Plan-Text hält jetzt Umfang und Lücke fest; ein Beleg im Closure-Abschnitt §7 fehlt noch (Auflage 3) |

## 6. Doku-Konsistenz gegen den Code (`internal/adapter/driving/cli/egress.go`)

- `command: "egress"`, `subcommand: "list"`: `runEgressList` (`egress.go:158`, `newDataEnvelope("egress", "list", toEgressListData(resp), nil, 0)`), gelesen. Handbuch Zeilen 2149, 2150 und 2191 (Datenfelder `active`, `sandboxProfile`, `egressEnabled`, `cloneSource`, `hosts`) und `cli-json-contract.md` §6.10 stimmen damit überein.
- Bare `egress --json`: `ErrEgressSubcommandRequired` (`egress.go:21`, `:73`), kein Envelope; Handbuch Zeile 2199 nennt den Fall zusammen mit `template`. `--dry-run` und `--diff`: `ErrEgressPreviewNotApplicable`, Exit 2; Handbuch Zeile 2101 stimmt.
- `jsontestutil` prüft das Enum nicht (`checkSubcommand` nur für `template` und `config`): die Aussage in §6.10 steht im Review als gelesen; der Verifier hat `jsontestutil.go` nicht neu gelesen (Abschnitt 9).

## 7. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | INFO (Release-Gate, kein Mangel) | Bestätigung des Projektinhabers für den abgeleiteten Wortlaut steht aus (Slice §1, §6 Risiko 1; Welle-Datei §3). Die Historie-Zeile `0.16.0` trägt die Klausel; bei Änderung oder Ablehnung: CR-Nachlauf mit Review des Diffs (Register `BEO-DOK/vertragsaenderung-nach-review`). Nicht im Verifier-Zugriff. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-2 | LOW | Plan-Tabelle §3 deckt zwei Diff-Gruppen nicht: (a) die Verweis-Nachzüge von `slice-mv` in fünf Dateien (`done/slice-cr-devcontainer-ollama.md`, `done/slice-v1-cli-egress-anzeige.md`, `open/README.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-verification-slice-cr-devcontainer-ollama.md`) und `open/slice-v2-spec-sandbox-egress-ollama.md`; (b) die Inhaltsänderung in `open/slice-v2-spec-sandbox-egress-ollama.md` (Übergabe der Aufzählung in `architecture.md`). Nachzug: eine Zeile „Verweis-Nachzug `slice-mv` (reine Pfad-Änderung)“ und eine Zeile für den Spec-Slice; Abgleich gegen `git diff --stat` nach der Closure wiederholen (`AGENTS.md` §3.10). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-3 | LOW | Zustandsmarker mit Ablaufpunkt in `docs/plan/planning/welle-sandbox-egress-ollama.md` (Zeile 44, „Zweites Release-Gate“): „Zustand: noch nicht gelandet, Bestätigung offen.“ Der CR ist im Lastenheft gelandet (`0.16.0`); offen ist nur die Bestätigung. Auf den Ist-Stand setzen („im Lastenheft `0.16.0` gelandet, Bestätigung offen“). Die Datei steht im Diff (Pfad-Nachzug), die Zeile wurde in der Fix-Runde nicht gelesen (`AGENTS.md` §3.10, Zustandsmarker). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-4 | INFO | Roadmap-Arbeitszeile „In Arbeit: `slice-cr-cli-json-command-enum` …“ ist bei der Closure auf „Nichts in Arbeit.“ zurückzusetzen, in derselben Reihenfolge wie beim Ollama-CR. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-5 | INFO | Zwei Randpunkte ohne Bezug zum Mangel: (a) `slice-v2-spec-sandbox-egress-ollama` DoD „Sicht und Historie“ nennt die Aufzählung „Implementiert sind …“ nicht; die Plan-Zeile trägt sie. Bei Bedarf in die DoD-Zeile aufnehmen. (b) Handbuch Zeile 2150 führt für `config` die Subcommands `get`, `set`, `list`; `config.go:367` liefert zusätzlich `"show"`. Der Drift besteht vor diesem CR und liegt außerhalb seines Gegenstands (Beobachtung, kein Befund dieses Slice). | <!-- d-check:ignore (Kennungen als Zitat) -->

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, keine Abweichung vom Plan-Wortlaut.

## 8. Auflagen

**Release-Gate (Projektinhaber, kein Mangel):** Bestätigung des abgeleiteten Wortlauts (V-1); bis dahin Risiko 1 mit Ausgang „weiter offen → Release-Gate“, Welle-Datei §3 hält den Release-Cut an.

**Beim Planner vor `done/` (Closure-Notiz §7 und übrige Haken):**

1. **Closure-Notiz §7** füllen: Vorab-Suche (Umfang: Lastenheft vom Implementer, Handbuch und Sicht erst im Review; Beleg), Stand des Gates, Steering-Loop-Eintrag (Lernpunkt: Vorab-Suche über `spec/` allein übersah die Nutzer-Doku, Register `BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht` und `BEO-DOK/plan-aussage-ueber-code` zählen), Folge-Slices (Spec-Slice mit `SPEC-001` und `architecture.md`). <!-- d-check:ignore (Kennungen als Zitat) -->
2. **Beobachtungs-Register** fortschreiben oder „keine Beobachtung“ in §7 vermerken.
3. **Risiken §6** mit Ausgängen: Risiko 1 weiter offen (Gate); Risiko 2 (Versionsreihenfolge) eingetreten wie erwartet, `0.16.0` nach `0.15.0`; Risiko 3 (Spec führt nicht nach) weiter offen, Übergabe getragen; Risiko 4 (Review vor Bestätigung) weiter offen, Nachlauf-Review nach Bestätigung.
4. **DoD-Haken** setzen: `make gates`/`doc-trace` (Abschnitt 1), Review, Doku-Update, Closure, Register, Risiken, Paarungen.
5. **V-2 und V-3** nachziehen, Plan-Tabelle nach der letzten Änderung gegen `git diff --stat` erneut abgleichen.
6. **V-4:** Roadmap auf „Nichts in Arbeit.“ zurücksetzen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3), über `make slice-mv`; Welle-Datei und Roadmap-Verweise auf den neuen Pfad (der Move-Commit zieht die Verweise nach). <!-- d-check:ignore (Kennungen als Zitat) -->

## 9. Nicht verifiziert

Die Bestätigung des Wortlauts durch den Projektinhaber (Gate, offen); Wortgleichheit mit einem Zeichen-Diff-Werkzeug (Lesen statt Werkzeug); `jsontestutil.go` (`checkSubcommand`) nicht neu gelesen, die Aussage stützt sich auf den Review; `make test-docker`, `make verify-depguard`, `make doc-immutable` (kein Produktcode, keine ADR im Diff); die Laufzeit-Ausgabe von `egress list --json` (keine Tests ausgeführt, nur Code und Doku gelesen); `make docs-check` einzeln (läuft als Teil von `make gates`, Exit 0); die Vorab-Suche des Implementers selbst (nicht wiederholt, die Review-Lesung dient als Gegenprobe).
