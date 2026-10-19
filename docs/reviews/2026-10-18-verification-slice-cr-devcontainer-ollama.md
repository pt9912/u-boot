# Verification-Report: slice-cr-devcontainer-ollama — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-cr-devcontainer-ollama` (`docs/plan/planning/in-progress/`), Lastenheft 0.14.0 auf 0.15.0, Commit-Range `19d8e2a..HEAD` (4 Commits, HEAD `bc3b313`, 15 Dateien). Kein Produktcode, keine Spezifikation, keine Architektur-Sicht, keine ADR im Diff. <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-cr-devcontainer-ollama.md` (F-1 MEDIUM, F-2 MEDIUM, F-3 LOW, F-4 LOW, F-5..F-8 INFO), DoD-Haken und Sensor-Behauptung des Slice (zwei Haken gesetzt, übrige offen), Slice §7 (Freigabe und Vorab-Suche). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Das Lastenheft steht wortgleich zum freigegebenen Text (bis auf die benannte Link-Form), der Lastenheft-Diff hat genau drei Hunks und sonst nichts, alle vier Sensoren sind grün, `ADR-0024` ist unverändert. Zwei Auflagen sind **Gates vor dem Release** und liegen beim Projektinhaber: **F-1** (Hosts des Modell-Pull, Stabilität des Blob-Hosts nicht gemessen) und **F-2** („weist aus“ gegen `--quiet`). Beide betreffen freigegebenen Text und bleiben unverändert offen; eine Änderung braucht eine neue Freigabe. Dazu drei LOW-Befunde zu fehlenden Test-Pins im Code-Slice (V-3 bis V-5) und Closure-Auflagen. Der Slice ist bis zur Closure-Notiz nicht schließbar. <!-- d-check:ignore (Kennungen als Zitat) -->

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** (Exit-Code des `make`-Prozesses selbst erfasst: `EXIT=0` direkt nach dem `make`-Aufruf in die Log-Datei geschrieben; Log nur mit begrenztem Read gelesen: `coverage-gate: OK`, Total 91,4 %) | Lauf im Verifier-Kontext, HEAD `bc3b313`, Arbeitsbaum davor und danach sauber |
| `make doc-trace` | pass, Exit 0 | „83 Anforderung(en), 0 Waise(n).“; `LH-FA-DEV-003` führt `slice-cr-devcontainer-ollama` und `ADR-0024`, Status ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make docs-check` | pass, Exit 0 | „d-check: 444 Datei(en) geprüft, 0 Befund(e)“ (vor Anlage dieses Reports) |
| `make doc-immutable RANGE=19d8e2a..HEAD` | pass, Exit 0 | „Range aufgeloest, 4 Commit(s) — OK“, 0 Befunde; `git diff --stat` zeigt keine Datei unter `docs/plan/adr/` |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | kein Produktcode, keine Importregel, kein Docker-Pfad des Produkts im Diff |

## 2. Lastenheft: Wortlaut-Konformität (Prüfpunkt 1)

Gelesen: Slice §1 „Wortlaut“ (Zeilen 53 bis 83) gegen `git diff 19d8e2a..HEAD -- spec/lastenheft.md`.

| Stelle | Ergebnis |
|---|---|
| 1: Spiegelstrich „Ollama (lokale Sprachmodelle, nur CPU)“ | wortgleich, letzter Strich der Beispiel-Liste |
| 1: Absatz (neun Sätze, „Ollama ist als eingebautes Feature `ollama` wählbar …“ bis „… durch `u-boot`.“) | wortgleich; einzige Abweichung die Link-Form `(#lh-fa-dev-008--egress-restriktion)` statt Relativpfad (Lastenheft verlinkt nur Anker der eigenen Datei, in Slice §7 benannt) |
| 3: Versionsfeld | `0.15.0`, Spaltenbreite unverändert, wortgleich |
| 3: Historie-Zeile 0.15.0 | wortgleich (Datum 2026-10-18, Quelle „Vereinbarung mit dem Projektinhaber“); trägt, anders als 0.14.0, keine Bestätigungs-Klausel |
| Sonst im Lastenheft | genau **drei Hunks** (Versionsfeld; Beispiel-Liste mit neuem Absatz; Historie), 5 Einfügungen, 1 Löschung; `LH-FA-DEV-006`, `LH-FA-CLI-005`, `LH-NFA-SEC-004` unverändert (Stellen 2 und 2a entfallen laut `ADR-0024` Entscheidung 5) | <!-- d-check:ignore (Kennungen als Zitat) -->

Keine Slice-, Wellen-, ADR-Kennung in den neuen Zeilen (`AGENTS.md` §3.4). Der Vergleich ist von Hand geführt (Lesen und `git diff`), nicht mit einem Zeichen-Diff-Werkzeug; die Freigabe selbst ist Gespräch und nicht prüfbar (Slice §7, laut DoD gewollt).

## 3. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff 19d8e2a..HEAD --stat`)

15 Dateien im Diff (Stand `bc3b313`): 1 Lastenheft, 1 Slice (Move rein in `e84d565`, Inhalt in `d531479`), 1 Roadmap, 11 Verweis-Nachzüge, 1 neuer Review-Report.

| Diff-Datei(en) | Plan-Zeile | Befund |
|---|---|---|
| `spec/lastenheft.md` | Zeile 1 | ok (Abschnitt 2) |
| `done/slice-cr-devcontainer-ollama.md` | Zeile 2 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `in-progress/roadmap.md` | Zeile 3 | ok: Arbeitszeile „In Arbeit: …“ ersetzt den Ruhe-Marker, Zustand und Anker ohne Chronik; Verweis der Welle-Zeile nachgezogen |
| `done/slice-v1-cli-egress-anzeige.md`, `done/slice-v1-devcontainer-ollama-messung.md`, `done/slice-v1-sandbox-egress-neuladen-entwurf.md`, `open/README.md`, `done/slice-cr-cli-json-command-enum.md`, `in-progress/slice-v1-devcontainer-feature-ollama.md`, `open/slice-v1-devcontainer-start-ohne-cli.md`, `open/slice-v2-spec-sandbox-egress-ollama.md`, `welle-sandbox-egress-ollama.md`, zwei Reports in `docs/reviews/` | Zeile 4 (Verweis-Nachzug) | ok; Diffs der Feature-, Spec-Slice-, Welle-, README- und Roadmap-Dateien gelesen: reine Pfad-Änderungen `open/` nach `in-progress/`, kein Inhalt | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/reviews/2026-10-18-slice-cr-devcontainer-ollama.md` (neu) | Zeile 5 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `CHANGELOG.md` | „unverändert“ | bestätigt: nicht im Diff |

Kein Diff ohne Plan-Zeile, keine Plan-Zeile mit Datei-Eintrag ohne Diff. Mit diesem Report wächst die Range auf 16 Dateien; die Zeile „docs/reviews (Review-Report, Verification Evidence)“ trägt ihn. Der Abgleich nach der Closure ist zu wiederholen (V-9). Review F-4 (a) bis (e): die Zustandsmarker sind unverändert und gehören in die Closure (Auflage 3).

## 4. DoD-Abgleich (Slice §2)

| DoD-Zeile | Stand | Beleg |
|---|---|---|
| Lastenheft 0.15.0 (Haken gesetzt) | **erfüllt, der Haken trägt** | Abschnitt 2; Freigabe nur im Gespräch (so vermerkt) |
| Vorab-Suche und Nachzug (Haken gesetzt) | **erfüllt, der Haken trägt** | Slice §7 listet die berührten Stellen (`LH-FA-DEV-006`, `-008`, `LH-NFA-SEC-004`, `LH-FA-CLI-005`, `LH-SA-CLI-002`, Glossar, `LH-FA-DEV-009`) mit Begründung „nicht geändert“; der Review hat die Stellen unabhängig gelesen (`LH-SA-CLI-002` Volltext, `LH-MOD-002`, Nicht-Kernumfang) ohne Widerspruch. Die Suche lief laut §7 als `grep -i` ohne Angabe des Ortes (Review F-8) | <!-- d-check:ignore (Kennungen als Zitat) -->
| `make gates` grün; `make doc-trace` ohne neue Waise | **erfüllt** (Abschnitt 1), Haken noch nicht gesetzt | Closure-Schritt |
| Review, Report liegt vor | **erfüllt**, Haken noch nicht gesetzt | Report vorhanden, über die freigegebene Fassung; keine Fix-Runde nach dem Review (kein Commit nach `bc3b313`) |
| Doku-Update (kein Handbuch, kein CHANGELOG) | erfüllt (nicht nötig) | kein Produktverhalten; `CHANGELOG.md` nicht im Diff |
| Closure-Notiz mit Steering-Loop-Eintrag, Register, Risiken §6, drei Paarungen | offen, zu Recht | §7 trägt Freigabe und Vorab-Suche, aber noch den Platzhalter „Closure durch den Planer“; §6-Risiken stehen auf „offen bis zum Start“ bzw. „offen bis O5“ (Auflage 3) |

Entscheidungs-Konformität gegen `ADR-0024` (`Accepted`): Absatz stimmt mit den Entscheidungen 1 (gepinntes Archiv, SHA-256), 3 (nicht signiert, Ausweisung), 4 (Server startet nicht von selbst), 5 (nur Loopback), 7 (Hosts in Default-Allowlist), 8 (benanntes Volume, Löschen des Containers statt Neubau) und 2 (nur Sandbox-Profil) überein; die Messstand-Zahlen (Docker 29.5.2, Colima 0.10.3, arm64, Ollama 0.40.6) stimmen mit `ADR-0024` §Messung überein. Entscheidungen 6 (GPU-Ausschluss) und 10 (`amd64`) stehen im Lastenheft als „nicht zugesagt: GPU, andere Architekturen“ ohne Vorbehalt über den Archivinhalt, vereinbar. Keine Accepted-ADR berührt (`AGENTS.md` §3.5), kein Gate gelockert (§3.6). <!-- d-check:ignore (Kennungen als Zitat) -->

## 5. Review-Findings F-1..F-8 und die zwei MEDIUM als Release-Gates (Prüfpunkt 2)

| Finding | Stand | Verifier-Beleg und Behandlung |
|---|---|---|
| F-1 MEDIUM Hosts des Modell-Pull, Stabilität des Blob-Hosts nicht gemessen | **offen, Entscheidung des Projektinhabers, Gate vor dem Release** | Freigegebener Text unverändert (Abschnitt 2). `ADR-0024` Entscheidung 7 begrenzt die Aussage („zwei Modelle, ein Zeitpunkt, ein Netzstandort; Stabilität nicht gemessen“); das Lastenheft sagt es nicht. Als **Risiko mit Ausgang „weiter offen“** in der Closure zu führen, Entscheidung: Bestätigung des Wortlauts oder Änderung (neue Freigabe, neuer Lastenheft-Commit, Historie), vor dem Release. Wo die Grenze bis dahin trägt: Feature-Slice §6 Risiko „Hosts des Modell-Pull wechseln“ (Rückgabe laut Welle §5) und `ADR-0024` Entscheidung 7 (Handbuch nennt Symptom und Weg); Spec-Slice §3 „Übergabe“ nennt „Hinweis auf die ungemessene Stabilität“ für `SPEC-019` | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-2 MEDIUM „weist aus“ gegen `--quiet` | **offen, Entscheidung des Projektinhabers, Gate vor dem Release** | Freigegebener Text unverändert. Das Lastenheft sagt zur Wirkung unter `--quiet` nichts; `ADR-0024` Entscheidung 3 (Warnung entfällt unter `--quiet`, ohne `SecurityRelaxation`) trägt die Lesart nur im Rang 4. Als Risiko „weiter offen“ in der Closure; Entscheidung: Bestätigung der ADR-Lesart oder Wortlaut-Änderung. Folge-Slice-Pin siehe V-3 | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-3 LOW „Gemessen ist das“ überdeckt Generator-Zusagen; Slice-Aussage „eine Aussage ohne Messwert steht nicht im Text“ nicht zutreffend | **offen, Planer in der Closure** | Slice-Text (nicht Lastenheft) berichtigen: Lastenheft ist Zielbild, Generator-Verhalten (Ausweisung, Allowlist, Wirkung im Profil) ist nicht gemessen (`ADR-0024`: „Nicht gemessen: Bau eines vom Generator erzeugten Dockerfiles“); Lastenheft-Text bleibt | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 LOW Zustandsmarker (a) bis (e) | **offen, Planer in der Closure** | Stellen bestätigt am Baum (Slice §1 Stelle 3 „Alt: 0.13.0“, „Variante, nur bei Bedarf“, Kopfzeile `Bezug`/`Berührte Spec-Stellen` mit „gegebenenfalls `LH-FA-DEV-006`“, §3 „der Implementer führt sie … aus“, §1 Ziel „wenn der Anzeige-CR 0.14.0 trägt“); zusätzlich Welle-Datei Tabellenzeile des Slice nennt `LH-FA-DEV-006` weiter als berührt, und `open/README.md` Zeile des Slice trägt „Wortlaut-Vorschlag … Freigabe des Projektinhabers nötig“ (V-8) | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-5 INFO Bedingungen Loopback (bei ungesetztem `OLLAMA_HOST`) und „wirkt im Sandbox-Profil“ nur in der ADR | **übergeben, Pin im Code-Slice nötig** | siehe V-4, V-5 |
| F-6 INFO Messzahlen stimmen | **bestätigt** | Verifier hat `ADR-0024` §Messung gegen den Absatz gelesen (Abschnitt 4); Rohdaten nicht neu gelesen | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-7 INFO Freigabe nur im Gespräch | **ausgewiesen, nicht prüfbar** | Abschnitt 8 |
| F-8 INFO Host-Werkzeuge im Review-Lauf | **ausgewiesen** | Register `BEO-WZG/host-werkzeug-ausserhalb-docker` (Ausweisung im Report; Zähler entscheidet die Closure); ein weiterer Beleg aus diesem Lauf siehe V-10 |

### Folge-Slice-Pins (gelesen: `in-progress/slice-v1-devcontainer-feature-ollama.md`, `open/slice-v2-spec-sandbox-egress-ollama.md`)

Gefragt war, ob die Pins als DoD oder Test getragen werden. Ergebnis je Pin:

| Pin | Feature-Slice (Code) | Spec-Slice | Befund |
|---|---|---|---|
| Hosts in der Default-Allowlist per Test | **getragen**: DoD Zeile 1 „`egressFeatureHosts()` trägt die zwei gemessenen Laufzeit-Hosts“, Golden Case „Egress-Liste die sortierte, duplikatfreie Vereinigung“ (auch mit `claude-code`), Plan-Zeile Tests „Hosts“; DoD Zeile 2 Docker-Test „ein Pull-Host erreichbar und ein fremder Host gesperrt“; Risiko „Hosts … wechseln“ | `SPEC-019` Zeile `ollama` mit zwei Hosts und Hinweis auf ungemessene Stabilität (§3 „Übergabe“, DoD Zeile 1) | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| Warnung unter `--quiet` entfällt (`ADR-0024` Entscheidung 3) | **nur teilweise**: Plan-Zeile `generate.go` nennt „entfällt unter `--quiet`“, `Bezug` nennt `LH-FA-CLI-005`; die **DoD** und die Test-Zeile (Katalogumfang, Render, Hosts, Volume, Start, Golden Case) nennen keinen `--quiet`-Test; „Hinweise genau wie die ADR vorschreibt“ im Golden Case prüft gerenderte Dateien, nicht Meldungen | `SPEC-043` Aufzählung der ungekennzeichneten Meldungen (§1 Berührte Spec-Stellen, §3) | **V-3 (LOW)** | <!-- d-check:ignore (Kennungen als Zitat) -->
| Loopback, `OLLAMA_HOST` ungesetzt (`ADR-0024` Entscheidung 5, „ein Test pinnt sie“) | **nur teilweise**: Docker-Test „zum Beispiel … Server erreichbar auf der gemessenen Adresse“ (Beispiel, kein Pflicht-Pin); §6 Risiko „Bedingung ohne Entsprechung“ verlangt Test oder Streichung generisch; weder „kein `OLLAMA_HOST` in `devcontainer.json`/Dockerfile“ noch „Aufruf über die Container-Adresse scheitert“ steht in der DoD | Spec-Slice nennt die Bindungsbedingung nicht (SPEC-023/-044-Nachzug nur Katalog, Warnung, Sandbox-Regel) | **V-4 (LOW)** | <!-- d-check:ignore (Kennungen als Zitat) -->
| Sandbox-Profil-Bedingung („wirkt im Sandbox-Profil“, ohne Profil keine Wirkung und Meldung) | **nur teilweise**: Plan-Zeile `devcontainer_sandbox.go` nennt „keine Wirkung ohne Sandbox-Profil“-Meldung; DoD/Golden Case nennen nur den Fall **mit** Sandbox-Profil, kein Negativfall ohne Profil (kein Block, Meldung) | `SPEC-044` Warnungstabelle „keine Wirkung ohne Sandbox-Profil“ (§3 Plan-Zeile) | **V-5 (LOW)** | <!-- d-check:ignore (Kennungen als Zitat) -->
| Spec-Nachzug (`slice-v2-spec-sandbox-egress-ollama`): SPEC-023 (Dockerfile-Block-Typ, `version:` ohne `source:` Exit-Code 10), SPEC-019 (zwei Hosts), SPEC-044 (Warnung, Sandbox-Regel begrenzt auf `sshd`/`claude-code`), SPEC-043 | n/a | **getragen**: §3 „Übergabe aus dem Entwurf `ADR-0024`“, DoD Zeile 1 mit Fundstellen-Tabelle, Risiko „Einträge zur Ollama-Warnung hängen an gebautem Code“; Start erst nach den Code- und CR-Slices | ok; Ergänzung V-6 (INFO) | <!-- d-check:ignore (Kennungen als Zitat) -->

## 6. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | MEDIUM (Gate vor dem Release) | **F-1 bleibt offen**: Der Lastenheft-Absatz (Rang 1) sagt „gehören die Hosts des Modell-Pull zur Default-Allowlist“, `ADR-0024` Entscheidung 7 belegt zwei Modelle, einen Zeitpunkt, einen Netzstandort und nennt die Stabilität „nicht gemessen“. Kein Verifier-Mittel schließt das. Bestätigung oder Änderung des Wortlauts durch den Projektinhaber vor dem Release; bei Änderung neuer Lastenheft-Commit samt Historie-Hinweis und Review des Diffs. Bis dahin Ausgang „weiter offen“. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-2 | MEDIUM (Gate vor dem Release) | **F-2 bleibt offen**: „weist das beim Erzeugen aus“ steht neben `LH-FA-DEV-006`/`LH-FA-CLI-005` Ausnahme 1, wo „ausweisen“ unter `--quiet` bleibt; `ADR-0024` Entscheidung 3 lässt die Warnung entfallen. Bestätigung der ADR-Lesart oder Änderung des Wortlauts durch den Projektinhaber vor dem Release. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-3 | LOW | Feature-Slice trägt den `--quiet`-Pin nicht als DoD oder Test (Abschnitt 5). Folge: Code-Slice kann die Warnung unter `--quiet` sichtbar bauen, und der Golden Case würde es nicht merken. Nachzug in `slice-v1-devcontainer-feature-ollama` §2: ein Test, dass die Lieferketten-Warnung und die „keine Wirkung“-Meldung unter `--quiet` entfallen (bzw. das Gegenteil, falls V-2 anders entschieden wird); Entscheidung V-2 vor dem Test festschreiben. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-4 | LOW | Loopback-Bedingung ohne ausdrücklichen Pin in der DoD (Abschnitt 5). Nachzug: Golden Case „kein `OLLAMA_HOST` in `devcontainer.json`, Dockerfile, Skript; kein `forwardPorts`“ und Docker-Test „Aufruf über die Container-Adresse scheitert, `127.0.0.1:11434` antwortet“ in die DoD (so steht es in der Fitness-Function-Tabelle der ADR, nicht im Slice); die Spezifikation führt die Bedingung „bei ungesetztem `OLLAMA_HOST`“ im Spec-Slice nach. |
| V-5 | LOW | Sandbox-Profil-Bedingung ohne Negativfall in DoD und Golden Case (Abschnitt 5). Nachzug: Test „ohne Sandbox-Profil kein Block im Dockerfile und die ungekennzeichnete Meldung“ (so in der ADR-Fitness-Function). |
| V-6 | INFO | Der Spec-Slice trägt die Bedingungen Loopback/`OLLAMA_HOST`, `amd64` „nicht verifiziert“ (`ADR-0024` Entscheidung 10, „in Spezifikation und Handbuch als nicht verifiziert“) und das Volume mit `id_ed25519` nicht ausdrücklich in seiner Übergabe-Liste; Lastenheft deckt `amd64` mit „andere Architekturen … nicht zugesagt“. Bei der Closure des CR in die Übergabe des Spec-Slice aufnehmen. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-7 | INFO | Review-Report ist ein Plan-/Spec-Review; er lief laut seinem Kopf mit HEAD `d531479`, der Report-Commit `bc3b313` ändert danach nur Dateien in `docs/reviews/`. Keine Fix-Runde, daher kein ungeprüfter Nachtrag (anders als beim Anzeige-CR). |
| V-8 | LOW | Zustandsmarker, die über Review F-4 hinausgehen: `open/README.md` führt den Slice weiter als „Lastenheft-CR mit Wortlaut-Vorschlag … Freigabe des Projektinhabers nötig“ (Freigabe liegt vor); Welle-Datei nennt `LH-FA-DEV-006` weiter als berührt (Stelle 2 entfallen). Bei der Closure auf den Ist-Stand setzen (`AGENTS.md` §3.10). | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-9 | INFO | Die Plan-Tabelle ist nach der Closure erneut gegen `git diff --stat` abzugleichen (diese Datei, Closure-Notiz, Roadmap-Rückstellung, `git mv` nach `done/`, Verweis-Nachzug). |
| V-10 | INFO | Regelhinweis des Verifiers (ehrlich): ein Shell-Aufruf nutzte `head` als Pipe-Stufe auf der Ausgabe von `make -n` (nur Anzeige der ersten Zeilen) und `cat` als Pipe-Stufe auf `git diff --stat`; ferner `echo "EXIT=$?"` zum Erfassen des Exit-Codes in der Log-Datei. Kein Nachweis stützt sich auf `head` oder `cat`; Nachweise sind die Exit-Codes von `make` und das Read-Tool. Sonst Read-Tool, `git`, `make`. Verstoß gegen „keine Tools außerhalb des Dockerfiles“ im Sinn des Lesens. |

Kein HIGH, keine DoD-Verletzung, keine Abweichung vom freigegebenen Wortlaut.

## 7. Auflagen für die Closure

1. **Gate vor dem Release (Projektinhaber):** F-1 und F-2 (V-1, V-2): Bestätigung oder Änderung des Wortlauts. Beide in §6 des Slice als Risiko mit Ausgang „weiter offen, Gate vor dem Release“ führen, in der Roadmap/Welle-Datei als Release-Gate sichtbar machen (Muster Enum-CR und Anzeige-CR). Bei Änderung: neue Freigabe, Lastenheft-Diff reviewen.
2. **Folge-Slice-Pins nachziehen (Planer, vor Start des Code-Slice):** `slice-v1-devcontainer-feature-ollama` §2 und §3 um V-3, V-4, V-5 ergänzen (Test `--quiet`, Loopback/`OLLAMA_HOST`, Negativfall ohne Sandbox-Profil); `slice-v2-spec-sandbox-egress-ollama` §3 um V-6. <!-- d-check:ignore (Kennungen als Zitat) -->
3. **Closure-Notiz (Planer):** F-3, F-4 und V-8 (Zustandsmarker auf Ist-Stand, `LH-FA-DEV-006` in Kopf/Bezug/Welle-Tabelle/`open/README.md`); §6-Risiken mit Ausgängen (Risiko 1 Messwert: entfallen, mit der Einschränkung F-3 und Gate F-1; Risiko 2 „Neubau“: entfallen, Lastenheft nennt ihn „nicht zugesagt“; Risiko 3 Version: eingetreten wie erwartet, `0.15.0`; Risiko 4 Review über die freigegebene Fassung: erfüllt); Steering-Loop-Eintrag; Register (`BEO-WZG/host-werkzeug-ausserhalb-docker` mit Belegen aus Review F-8 und V-10; `BEO-DOK/spec-bedingung-ohne-code-entsprechung`-Bezug); DoD-Haken (`make gates`, Review, Doku-Update, Closure, Register, Risiken, Paarungen) setzen. <!-- d-check:ignore (Kennungen als Zitat) -->
4. Plan-Tabelle nach der letzten Änderung gegen `git diff --stat` wiederholen (V-9).
5. Roadmap auf „Nichts in Arbeit.“ zurücksetzen; `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3), über `make slice-mv`. <!-- d-check:ignore (Kennungen als Zitat) -->

## 8. Nicht verifiziert

Die Freigabe des Wortlauts im Gespräch am 2026-10-18 (kein Datei-Beleg, laut Slice gewollt); Wortgleichheit mit einem Zeichen-Diff-Werkzeug (Lesen statt Werkzeug, Docker-only-Regel); die Rohdaten des Messprotokolls (nur `ADR-0024` gelesen); die Stabilität des Blob-Hosts (nicht messbar, F-1); `make test-docker`, `make verify-depguard` (kein Produktcode); Wirkung der Warnung unter `--quiet` (entsteht im Code-Slice); die Vorab-Suche des Implementers (`grep -i`) nicht wiederholt, die Review-Lesung der Lastenheft-Stellen dient als Gegenprobe. <!-- d-check:ignore (Kennungen als Zitat) -->
