# Verification-Report: slice-cr-sandbox-egress-anzeige — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-cr-sandbox-egress-anzeige` (`docs/plan/planning/in-progress/`), Lastenheft 0.13.0 auf 0.14.0, Commit-Range `1b833ea..HEAD` (5 Commits, HEAD `36fb552`, 11 Dateien). Kein Produktcode, keine Spezifikation, keine Architektur-Sicht, keine ADR im Diff. <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-cr-sandbox-egress-anzeige.md` (F-1 MEDIUM, F-2..F-5 LOW, F-6..F-9 INFO; Findings-Nummern bis F-9), Fix-Runde `36fb552` (Message nennt F-1 bis F-7), DoD-Haken und Sensor-Behauptung aus dem Slice. <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflage.** Die Stellen 1 bis 3 stehen wortgleich zum freigegebenen Wortlaut im Lastenheft, Stelle 4 ist als abgeleiteter, nicht freigegebener Satz im Slice und in der Historie-Zeile ehrlich ausgewiesen, sonst ist im Lastenheft nichts verändert. Alle drei Sensoren sind grün. Auflage: **Stelle 4 (Zeile `u-boot egress list` in der Befehlstabelle und die Klausel der Historie-Zeile 0.14.0) wird vor dem Release vom Projektinhaber bestätigt**; die Klausel „Bestätigung ausstehend“ im Lastenheft wird danach entfernt oder auf den Ist-Stand gesetzt (V-1). Der Slice ist bis dahin nicht schließbar.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** (Exit-Code des `make`-Prozesses selbst erfasst: `gates exit=0`, nicht der von `tail`) | Lauf im Verifier-Kontext nach HEAD `36fb552`, Arbeitsbaum sauber |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass, Exit 0 | „d-check: 421 Datei(en) geprüft, 0 Befund(e)“ (vor Anlage dieses Reports) |
| `make test-docker`, `make verify-depguard`, `make doc-immutable` | nicht gelaufen | kein Produktcode, keine Importregeln, kein Docker-Pfad des Produkts, keine ADR im Diff (`git diff --name-status` zeigt nur Planung, Reviews und `spec/lastenheft.md`) |
| Test `TestE2E_SandboxDevcontainer_EgressSubdomains` (Beleg des Subdomain-Satzes) | nicht neu gelaufen | braucht Colima mit `NET_ADMIN` und Netz; Review F-8 hat Messung des Vorgänger-Slice gelesen |
| Host-Werkzeuge | siehe V-5 | `git`, `make`, Read-Tool; Ausnahme `tail` und `echo` in einem Shell-Aufruf (nur Lesen der Log-Dateien) |

## 2. Lastenheft: Wortlaut-Konformität (Read-Tool, Zeile für Zeile)

| Stelle | Slice-Wortlaut (§1) | Lastenheft | Ergebnis |
|---|---|---|---|
| 1: neuer Spiegelstrich in `LH-FA-DEV-008` | Slice-Datei Zeile 52 | `spec/lastenheft.md` Zeile 472, direkt nach „Erlaubte Ziele“ (Zeile 471) | **wortgleich** (vier Sätze, Backticks, „gälte“, Klammerinhalte); übrige Spiegelstriche unverändert | <!-- d-check:ignore (Kennungen als Zitat) -->
| 2: Aufzählung in `LH-FA-CLI-005` | Slice-Datei Zeile 64 | Aufzählung `config get`, `config list`, `egress list`, `logs`, `template list`; Rest des Absatzes laut `git diff` unverändert | **wortgleich** | <!-- d-check:ignore (Kennungen als Zitat) -->
| 3: Versionsfeld | `0.14.0` | `| Version          | 0.14.0 |`, Spaltenbreite unverändert | **wortgleich** |
| 3: Historie-Zeile 0.14.0 | Slice-Datei Zeile 74 | Zeile hat Datum 2026-10-18, beide Anker, Quelle „Vereinbarung mit dem Projektinhaber“ und **zusätzlich** die Klausel zu Stelle 4 | Wortlaut-Teil wortgleich; die Klausel ist Zusatz, im Slice (§1 „trägt zusätzlich die Klausel“, §7) ausgewiesen; siehe V-1 |
| 4: Zeile in der Befehlstabelle `LH-SA-CLI-002` | Slice-Datei Zeile 82 | `spec/lastenheft.md` Zeile 1121: `| u-boot egress list | Egress-Liste mit Herkunft anzeigen | V2 |`, zwischen `config migrate` und `template list` | **nicht freigegeben, abgeleitet**; Priorität V2 gleich `LH-FA-DEV-008` (Priorität V2, Zeile 468) | <!-- d-check:ignore (Kennungen als Zitat) -->
| Sonst im Lastenheft | — | `git diff 1b833ea..HEAD -- spec/lastenheft.md`: fünf Hunks (Versionsfeld, `LH-FA-CLI-005`, `LH-FA-DEV-008`, Tabellenzeile, Historie), 5 Einfügungen, 2 Löschungen | nichts sonst verändert | <!-- d-check:ignore (Kennungen als Zitat) -->

Weitere Lastenheft-Prüfungen: Link-Anker `#lh-fa-dev-008--egress-restriktion`, `#lh-fa-cli-005--verbosity-und-logging`, `#lh-sa-cli-002--vorgesehene-befehle` in den neuen Zeilen lösen auf (`make docs-check`, 0 Befunde); keine Slice-, Wellen-, ADR-Kennung in den geänderten Zeilen (Diff gelesen; nur Anker der eigenen Datei). Tabellenkopf „Priorität: MVP/V1 gemischt (siehe Spalte)“ unverändert; die Tabelle trug schon eine Zeile `Later`, jetzt zusätzlich `V2` (V-3). Der Vergleich ist von Hand geführt (Lesen und `git diff`), kein Zeichen-Diff-Werkzeug; der Slice-Wortlaut steht stellvertretend für die Vorlage im Gespräch, die Freigabe selbst ist nicht prüfbar.

## 3. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff 1b833ea..HEAD --name-status`)

11 Dateien im Diff: 1 neu (Review-Report), 1 Rename mit Änderung (`open/` nach `in-progress/`, R061), 9 geändert.

| Diff-Datei | Plan-Zeile | Befund |
|---|---|---|
| `spec/lastenheft.md` | Zeile 1 (Stellen 1 bis 4) | ok |
| `done/slice-cr-sandbox-egress-anzeige.md` (Move + Inhalt) | Zeile 2 | ok; Move in `8a7a827` rein, Inhalt danach | <!-- d-check:ignore (Kennungen als Zitat) -->
| `in-progress/roadmap.md` | Zeile 3 | ok (Arbeitszeile; Rückstellung ist Closure-Aufgabe) |
| `in-progress/slice-v1-cli-egress-anzeige.md`, `open/slice-v2-spec-sandbox-egress-ollama.md` | Zeile 4 | ok (Verweis-Nachzug und Übergabe aus dem Review, Diff gelesen) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `done/slice-v1-egress-wildcard-pruefung.md`, `done/slice-v1-sandbox-egress-neuladen-entwurf.md`, `open/README.md`, `welle-sandbox-egress-ollama.md`, `docs/reviews/2026-10-18-verification-slice-v1-sandbox-egress-neuladen-entwurf.md` | Zeile 5 | ok (je 2 bis 6 Zeilen Pfad-Nachzug) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/reviews/2026-10-18-slice-cr-sandbox-egress-anzeige.md` (neu) | Zeile 6 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| dieser Report | Zeile 6 („Verification Evidence“) | ok, Report ist die letzte Datei der Range |
| `CHANGELOG.md` | Zeile 7 „unverändert“ | bestätigt: nicht im Diff |

Kein Diff ohne Plan-Zeile, keine Plan-Zeile mit Datei-Eintrag ohne Diff. Der Abgleich gilt für den Stand nach der Fix-Runde (`36fb552`); mit diesem Report wächst die Range auf 12 Dateien (V-6). <!-- d-check:ignore (Kennungen als Zitat) -->

## 4. DoD-Abgleich (Slice §2)

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| Lastenheft 0.14.0 (Stellen 1 bis 3, Haken gesetzt) | **erfüllt, der Haken trägt** | Abschnitt 2; Freigabe nur im Gespräch (so in der DoD und §7 vermerkt, nicht prüfbar) |
| Vorab-Suche und Nachzug (Haken gesetzt) | **erfüllt, der Haken trägt** | §3 Vorab-Suche; Review hat sie wiederholt und die Lücke F-1 (Befehlstabelle) gefunden, mit Stelle 4 geschlossen; weitere Stellen (`LH-FA-CLI-007`, `LH-SA-CLI-001`, `LH-NFA-USE-004`, `LH-NFA-PERF-001`) als Beispiel-Aufzählungen begründet unverändert | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make gates` grün; `make doc-trace` ohne neue Waise | erfüllt (Abschnitt 1), Haken noch nicht gesetzt | Closure-Schritt |
| Review, Report liegt vor | erfüllt, Haken noch nicht gesetzt | Report vorhanden; die Fix-Runde `36fb552` (Lastenheft: Tabellenzeile und Klausel) ist nach dem Review entstanden und nicht erneut reviewt; der Verifier hat sie gelesen (V-2) |
| Doku-Update, `CHANGELOG.md` | erfüllt (nicht nötig) | kein Produktverhalten; Plan-Zeile „unverändert“ bestätigt |
| Closure-Notiz, Steering-Loop-Eintrag, Register, Risiken §6, drei Paarungen | offen, zu Recht | §7 trägt noch den Platzhalter „Wird bei der Closure gefüllt“; Risiken 1 und 3 stehen auf „offen“ (Risiko 3: der Review lief über die freigegebene Fassung; Ausgang bei der Closure zu setzen) |

Entscheidungs-Konformität: Befehlsname `u-boot egress list` und Priorität V2 im Lastenheft stimmen mit dem Slice (§1, §7) und den Folge-Slices überein; `slice-v1-cli-egress-anzeige` §1 führt Ziel, Abgrenzung und die Übergabe (Hinweis in der Ausgabe, Auslegung unter `--quiet`, Reichweite der Clone-Quelle, Stellen mit der Aufzählung) als Bedarf ohne Wortlaut; der freigegebene Lastenheft-Wortlaut bleibt unverändert. `slice-v2-spec-sandbox-egress-ollama` trägt die Spezifikations-Nachzüge. Keine Accepted-ADR berührt (AGENTS.md §3.5), kein Gate gelockert (§3.6), Lastenheft ohne abwärts gerichtete Kennungen (§3.4). <!-- d-check:ignore (Kennungen als Zitat) -->

## 5. Review-Findings F-1..F-9

Der Review-Report trägt keine Spalte „Behandlung“; die Zuordnung steht in der Commit-Message `36fb552` und im Slice (§1, §3, §7). Verifier-Prüfung am Baum:

| Finding | Stand | Verifier-Beleg |
|---|---|---|
| F-1 MEDIUM Befehlstabelle führt den Befehl nicht | **behoben, mit Auflage** | Zeile 1121 im Lastenheft (Abschnitt 2); als abgeleiteter, nicht freigegebener Satz ausgewiesen (Slice §1 Stelle 4, §7, Historie-Klausel). Der Review verlangte „neue Freigabe des Projektinhabers“; sie liegt nicht vor, die Bestätigung steht als Auflage 1 aus. |
| F-2 LOW „wirksame“ Allowlist gegen laufenden Container | **begründet, nicht im Lastenheft geändert** | Der freigegebene Wortlaut bleibt unverändert (Slice §7 „Abgrenzung der Anzeige“); Bedarf „Hinweis in der Ausgabe“ im Code-Slice §1 und „Grenze nennen“ im Spezifikations-Slice §1. Das Lastenheft nennt die Grenze weiter nicht (V-4). |
| F-3 LOW Plan-Tabelle ohne Verweis-Nachzug | **behoben** | Plan-Zeilen 4 und 5; Abschnitt 3 |
| F-4 LOW Zustandsmarker „Vorschlag“, „Arbeitsname“, Platzhalter | **behoben** | §1: „Befehlsname (freigegeben)“, Überschrift „Wortlaut (… freigegeben, im Lastenheft 0.14.0 eingetragen)“, Datum 2026-10-18 statt Platzhalter (Diff `21df190..36fb552` gelesen) |
| F-5 LOW Start-Bedingung `Verantwortlich:` | **ausgewiesen, nicht erfüllt** | §7 „Ist-Stand des Starts“: Feld bleibt `—`; Bedingung in §4 unverändert (Entscheidung beim Projektinhaber) |
| F-6 INFO zwei Auslegungsfragen | **übergeben** | Code-Slice §1 („Auslegung unter `--quiet`“, „Reichweite der Clone-Quelle“) mit Fundstellen; Spezifikations-Slice §1 |
| F-7 INFO weitere Stellen der Aufzählung | **übergeben** | `SPEC-043`, Handbuch-Zeile `--quiet`, Kommentar in `root.go`, Test `TestExecute_Quiet_QueryCommandsKeepTheirOutput` in den beiden Folge-Slices genannt; nicht geändert (Zuständigkeit dort) | <!-- d-check:ignore (Kennung als Zitat) -->
| F-8 INFO Subdomain-Satz belegt | **ok, kein Handlungsbedarf** | Messung E1 und Test; Messgrenzen im Handbuch; nicht neu gelaufen (Abschnitt 1) |
| F-9 INFO Host-Werkzeuge im Review-Lauf | **ausgewiesen** | im Report benannt; Register `BEO-WZG/host-werkzeug-ausserhalb-docker`; ob der Zähler fortgeschrieben wird, entscheidet die Closure (DoD „Register“) |

## 6. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | LOW | Die Historie-Zeile 0.14.0 im Lastenheft (Rang 1, vertraglich) trägt den Zustandsmarker „Bestätigung des Projektinhabers vor dem Release ausstehend“. Ein Marker mit Ablaufpunkt wird beim Ablauf entfernt oder auf den Ist-Stand gesetzt (AGENTS.md §3.10), doch die Historie ist Verlaufstext, der sonst nicht umgeschrieben wird. Die Bestätigung braucht daher eine ausdrückliche Änderung genau dieser Klausel (Zustand „bestätigt am <Datum>“ oder Streichen), im selben Zug wie die Bestätigung von Stelle 4. Sonst steht nach dem Release ein falscher Zustand im Vertrag. |
| V-2 | INFO | Die Fix-Runde `36fb552` änderte `spec/lastenheft.md` (Tabellenzeile, Historie-Klausel) nach dem Review; kein zweites Review über diese drei Zeilen. Der Verifier hat sie gegen den Slice-Text gelesen (ohne Befund). Wer die Rollentrennung eng liest, kann einen Review des Diffs der Fix-Runde verlangen. |
| V-3 | INFO | Tabellenkopf „Priorität: MVP/V1 gemischt (siehe Spalte)“ nennt weder `Later` (bestand schon) noch das neue `V2`; die Spalte ist maßgeblich („siehe Spalte“). Kein Widerspruch zum Wortlaut, beim Bestätigen von Stelle 4 mit entscheidbar. |
| V-4 | INFO | „Die wirksame Allowlist“ bleibt ohne Grenzenennung im Lastenheft (F-2); die Grenze lebt nur in den Plänen der Folge-Slices. Trägt der Code-Slice sie nicht (Hinweis in der Ausgabe mit Test), besteht die Lücke fort. Als Pflichtpunkt im Code-Slice geführt (`slice-v1-cli-egress-anzeige` §1). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-5 | INFO | Regelhinweis des Verifiers: ein Shell-Aufruf nutzte `tail` am Host (Lesen der Log-Dateien von `make gates`, `make doc-trace`, `make docs-check`) und `echo`; kein Schreiben, kein Nachweis aus diesem Werkzeug außer der Anzeige der Logzeilen. Der Nachweis ist der Exit-Code von `make`. Verstoß gegen „alles mit make und docker“ im Sinn von Lesen, sonst `git`, `make` und Read-Tool. |
| V-6 | INFO | Die Zahl der Diff-Dateien (11) ändert sich durch diesen Report und durch die Closure; die Plan-Tabelle ist nach der letzten Änderung erneut abzugleichen (AGENTS.md §3.10). |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung, keine Abweichung vom freigegebenen Wortlaut der Stellen 1 bis 3.

## 7. Auflagen für die Closure

1. **Stelle 4 vor dem Release vom Projektinhaber bestätigen** (Zeile in `LH-SA-CLI-002` und Klausel der Historie-Zeile 0.14.0); danach Klausel „Bestätigung ausstehend“ im Lastenheft auf den Ist-Stand setzen oder streichen (V-1), Vermerk mit Datum in §7. Lehnt er die Zeile ab oder ändert die Priorität, ist das ein Nachlauf (Diff der Änderung reviewen, nicht die ganze Datei). <!-- d-check:ignore (Kennungen als Zitat) -->
2. Entscheidung zu `Verantwortlich:` (bleibt `—`, in §7 festgehalten) und Ausgang der Risiken in §6 setzen (Risiko 3: Review lief über die freigegebene Fassung; Risiko 1: offen bis zur Architect-Entscheidung des Code-Slice).
3. §7 füllen: Vorab-Suche, Steering-Loop-Eintrag, Register (`BEO-WZG/host-werkzeug-ausserhalb-docker`: Beleg aus Review und Verifikation, falls zählbar), Folge-Slices, Risiken; DoD-Haken (`make gates`, Review, Doku-Update, Closure, Register, Risiken, Paarungen) setzen.
4. Plan-Tabelle nach der letzten Änderung gegen `git diff --stat` wiederholen (V-6).
5. Roadmap auf „Nichts in Arbeit.“ zurücksetzen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (AGENTS.md §3.3), über `make slice-mv`. <!-- d-check:ignore (Kennungen als Zitat) -->
6. Optional: Review der Fix-Runde (V-2), Tabellenkopf der Befehlstabelle (V-3).

## 8. Nicht verifiziert

Die Freigabe des Wortlauts im Gespräch am 2026-10-18 (kein Datei-Beleg, laut Slice gewollt); Wortgleichheit mit einem Zeichen-Diff-Werkzeug (Lesen statt Werkzeug, Docker-only-Regel); Subdomain-Messung auf anderen Messständen und der e2e-Test selbst; `make test-docker`, `make verify-depguard`, `make doc-immutable` (kein Produktcode, keine ADR); Wirkung des Hinweises zur Abgrenzung (kommt im Code-Slice).
