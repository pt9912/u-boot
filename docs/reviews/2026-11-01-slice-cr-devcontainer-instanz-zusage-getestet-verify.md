# Review-Report: slice-cr-devcontainer-instanz-zusage-getestet (Verifier) — 2026-11-01

**Review-Art:** Code — Verifikation gegen Plan, DoD, Entscheidungen und Sensors (Modul 11: „Bauen wir es richtig?“). Kein Review des Diffs und keine Validierung.

**Gegenstand:** Slice [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md) im Hauptbaum, HEAD `1bf451c4` (Fix-Runde des Reviews enthalten). Der Slice liegt in `in-progress/`. Diff dieses Slice: `3ec5cda5` (Freigabe), `d38c0908` (Lastenheft 0.26.0), `1bf451c4` (Fix-Runde) samt den Move-Commits `a8f0ece1`, `c077a6a0`, `8296197c`, `44d14987`. Mitgeprüft: Änderung A, Commit `9fec8752` ([`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Unterabschnitt `shell`, auf die Engine 29.9.2 beschränkt; Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md)).

**Skill:** `.harness/skills/verifier.md` existiert nicht; die Rolle folgt `harness/roles.md` und `harness/verification.md`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md) (§1, §2 DoD, §3 Plan, §6 Risiken, §7 Freigabe 2026-11-01); Muster [`slice-cr-devcontainer-instanz`](../plan/planning/done/slice-cr-devcontainer-instanz.md); [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md)
- [Review-Report](2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md) (0 HIGH, 0 MEDIUM, F-1 bis F-8)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Lastenheft 0.26.0), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)
- [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md), [`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)
- `harness/verification.md`, `AGENTS.md` (Hard Rules, besonders §3.4, §3.10)

**Sensors dieses Laufs:**

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (Docker-Modus, HEAD `1bf451c4`) | pass | Exit 0, selbst gefahren; Coverage 92,20 % gegen Schwelle 90 %; `baseline-verify` v6.18.1 OK (54 Dateien); `ablage-pruefung-test` 26 Fälle, 0 Fehlschläge; `docs-check` Teil von `gates` |
| `make doc-trace` | pass | Exit 0; „84 Anforderung(en), 0 Waise(n).“; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) führt [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md) in der Slice-Spalte, Status ok |
| `make docs-check` mit diesem Report | siehe Abschlussbericht (nach dem Schreiben gefahren) | — |
| CI-Beleg `integration-docker` | gelesen | Lauf 41285105546, Kopf `a6264402`, `push` auf `main`, Konklusion `success`; Job-Log 124784712075 enthält `--- PASS: TestDevcontainerInstance_Docker_TwoInstancesAreIsolated (16.32s)` und `ok  github.com/pt9912/u-boot/internal/e2e 849.686s`; `git diff a6264402 HEAD --stat -- internal/e2e CHANGELOG.md` leer, der Beleg trägt den heutigen Teststand |
| `make test-docker` | **nicht gefahren** | Entscheidung des Projektinhabers; Ersatz ist der gelesene CI-Lauf (oben) |
| `make verify-depguard` | nicht gefahren | keine Importregel berührt (Diff trägt nur Doku und einen Doc-Kommentar) |

---

## DoD-Abgleich (Punkt für Punkt, gegen gelesene Fundstellen)

| # | DoD-Punkt | Stand | Beleg |
|---|---|---|---|
| 1 | Wortlaut ohne offene Entscheidung | erfüllt | §7 vermerkt die Freigabe 2026-11-01 (Variante 1, Fragen 1 bis 4) vor dem Lastenheft-Commit: `3ec5cda5` liegt vor `d38c0908`. |
| 2 | Lastenheft `0.26.0` | erfüllt | Zeile 9 `0.26.0`; Historie-Zeile `0.26.0` (Zeile 1610) additiv hinter `0.25.0`, Verweis „Vereinbarung mit dem Projektinhaber“; Absatz Zeile 504: Satz 1 „Zwei“ statt „Bis zu vier“ (Entscheidung 2), Satz 2 gegen §7 maschinell verglichen (Freigabe-Text, `(…)` aufgefüllt durch „etwa bei einem Go-Build oder einem Modell-Pull“, ist Teilzeichenkette des Lastenhefts: wahr). Die Zeile `0.25.0` ist im Diff unberührt. |
| 3 | Folgen nachgezogen | erfüllt | `docs/maintainer/zusagen-und-messstaende.md` Zusage-Teil auf den Wortlaut gesetzt; Handbuch 4.9 („wenn aktiviert“, ohne Isolationsaussage zu Docker-Speicher und Modellen, keine Kennung, keine Zahl); [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeilen 837 ff. und `CHANGELOG.md` Zeile 16 in §7 mit gelesener Fundstelle als unverändert geführt (Suche in `spec/` nach Zahl und Grenze ohne Treffer; `CHANGELOG.md` und `internal/e2e` im Diff seit `a6264402` leer). Die Text-Pins `dockerMount` und `ollamaMount` (`generate_ollama_test.go`) trägt der Review (Prüfpunkt 3). |
| 4 | Vorab-Suche und Beleg-Tabelle | erfüllt | §7: `grep -i`-Suche mit Treffern nur in Zeile 504 und Historie `0.25.0`; Tabelle „Bedingung → Beleg“ mit fünf Zeilen. Zeile „Linux `amd64`, lokale Docker-Engine“ nennt ehrlich „kein Test pinnt die Architektur“ (Zusage-Menge der Anforderung und Runner). |
| 5 | `make gates` grün, `make doc-trace` ohne neue Waise | erfüllt (dieser Lauf) | oben; die Läufe und Exit-Codes sind in §7 des Slice noch nicht eingetragen (V-1). |
| 6 | Review durchgeführt, Report liegt vor | erfüllt | [Review-Report](2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md), 0 HIGH und 0 MEDIUM; F-1 bis F-4 und F-7 in `1bf451c4` behoben (gelesen: `devcontainer_shell.go` Kommentar, Wortlaut von [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), §7-Beleg, Plan-Tabellen, §6-Risiko). |
| 7 | Doku-Update `CHANGELOG.md` | erfüllt (nicht nötig) | Zeile 16 nennt Docker-Tests mit zwei Instanzen, keine Zusage über vier; unberührt, kein Produktverhalten. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | §7 trägt Freigabe, Vorab-Suche, Folgen und Beleg-Tabelle; die Felder „Was hat funktioniert“, „Was ging anders“, „Steering-Loop-Eintrag“, „Folge-Slices“, „Risiken aus §6“, „Drei Paarungen“ sind leer; die Zeile „Gegenstand“ (nur bei Variante 2) steht noch und ist zu löschen. |
| 9 | Beobachtungs-Register fortgeschrieben | **offen** | Kein Beleg für diesen Slice: `BEO-DOK/zusage-weiter-als-messung/evidence/` enthält keine Datei zu [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md), `state.md` nennt den Slice nicht. Die Finding-Klassen des Reviews sind nicht eingetragen. |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **offen** | Von sieben Risiken tragen zwei einen Ausgang (Lastenheft ist Vertrag: entfallen; Test-Lücken Review F-7: weiter offen mit Anker). Fünf stehen noch auf „offen bis …“, obwohl der Bedingungspunkt erreicht ist (Zuordnung unten). |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen** | Hängt an Punkt 8 und 9; Haken fehlen. |

Alle Kästchen in §2 stehen auf `[ ]`, auch die erfüllten; der Haken wird mit der Closure gesetzt (V-1).

## Zuordnung der Risiken aus §6 (Vorschlag für die Closure, nicht entschieden)

| Risiko | Befund der Verifikation | Vorschlag für den Ausgang |
|---|---|---|
| Wortlaut sagt mehr zu, als Messung und Suite belegen | Satz 2 gegen den Test und die Anforderung gelesen: jedes Merkmal belegt oder aus der Zusage-Menge übernommen; „Isolation der Workspace-Volumes“ aus Satz 1; eingetreten ist es nicht. Satz 1 trägt weiter „soweit aktiviert … eigenes Volume“ für Docker-Speicher und Modelle (Freigabe, Review F-5). | entfallen |
| Suite belegt „zwei“ erst nach einem grünen Lauf | Lauf 41285105546 grün, Test im Job-Log `PASS`, `internal/e2e` seither unverändert. | entfallen |
| Handbuch-Satz zu Docker-Speicher und Modellen | umgeschrieben und gepinnt (Text-Pins); Isolation nicht mehr behauptet. | entfallen |
| Version `0.26.0` vergeben | bei Start gelesen, frei, geschrieben. | entfallen |
| Test-Lücken der Eingrenzung von `shell` (Review F-7) | (a) und (b) bleiben, (c) bleibt (Job-Log nennt nur den Client 29.9.2, nicht die Engine des Runners); siehe V-3. | weiter offen (zulässig, Anker steht) |
| Lastenheft wird enger, nicht weiter | CR mit Freigabe. | entfallen (steht schon) |
| Accepted-ADR trägt den Satz | §7 liest [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) Zeilen 73 und 97; die Aussage „die Zusage nennt eine Menge“ bleibt wahr; die ADR ist nicht im Diff. | entfallen |

## Plan-vs-Code (`git diff --stat` gegen §3)

Diff `3ec5cda5..HEAD` (12 Dateien) und `d38c0908`/`3ec5cda5` gegen die Plan-Tabelle:

| Datei im Diff | Plan-Zeile | Befund |
|---|---|---|
| `spec/lastenheft.md` | ja | passt |
| `docs/maintainer/zusagen-und-messstaende.md` | ja | passt |
| `docs/user/benutzerhandbuch.md` | ja („update, falls die Lesung es verlangt“) | passt, Lesung in §7 |
| die Slice-Datei selbst, `open/README.md`, `roadmap.md` | ja | passt |
| [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md), zwei Reports `docs/reviews/…test-leicht*.md` | ja (Werkzeug-Nachzug, F-4 behoben) | passt |
| der [Review-Report](2026-11-01-slice-cr-devcontainer-instanz-zusage-getestet.md) | ja | passt |
| `spec/spezifikation.md`, `internal/adapter/driven/docker/devcontainer_shell.go` (in `1bf451c4`), Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) | Plan-Zeilen „spec/spezifikation.md unverändert“, „internal/** unverändert“ | Für den CR gilt „unverändert“ (kein Diff aus `3ec5cda5` oder `d38c0908`). Die Änderung stammt aus Änderung A beziehungsweise F-1 und gehört dem Slice der Engine-Eingrenzung; die Plan-Zeile „unverändert“ des CR gilt damit nur für den CR-Inhalt, nicht für den Stand von `HEAD` (V-2). |
| `docs/plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/` (Plan: update) | ja | Eintrag im Diff fehlt (V-1, DoD 9) |
| `CHANGELOG.md`, `README.md`, `README.de.md`, `spec/architecture.md` | unverändert | kein Diff, passt |

Decken-Regel: `git diff a6264402 HEAD -- spec` (nur Zeilen mit `+`) enthält keinen Verweis auf Slice, ADR, Roadmap, Carveout oder Planung; die Historie-Zeilen von Lastenheft und Spezifikation nennen nur „Vereinbarung mit dem Projektinhaber“ beziehungsweise „Entscheidung des Projektinhabers“. Erfüllt ([`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md), `AGENTS.md` §3.4).

## Änderung A (`9fec8752`, `SPEC-049` auf Engine 29.9.2)

Gelesen: Der Listenpunkt ohne Terminal und der Absatz *Engine-Abhängigkeit* sagen übereinstimmend: zugesagt Engine 29.9.2 mit Client 29.9.2, Linux `amd64`; Engine 28.0.4 erkannt, nicht zugesagt, nicht bestätigt (nach `1bf451c4`); der Launcher prüft die Form unabhängig von der Version (Review-Prüfpunkt 6). Das Änderungsprotokoll der Spezifikation trägt die Zeile vom 2026-11-01. Der Port-Kommentar `devcontainer_launcher.go` steht im Indikativ ohne Chronik. Das Lastenheft nennt für `shell` weiterhin „lokale Engine“ ohne Version; die Spezifikation ist enger als der Vertrag (Review F-6, Zuständigkeit Projektinhaber, kein Konflikt).

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | INFO | Die Closure-Arbeit steht aus: DoD-Kästchen alle `[ ]`, Closure-Notiz ohne Lerneintrag, Register ohne Beleg für diesen Slice, §7 ohne die Läufe von `make gates` und `make doc-trace`, fünf Risiken ohne Ausgang. Kein Fehler des Gebauten, sondern Stand vor der Closure. | `harness/verification.md` §Pflicht bei Slice-Closure; DoD §2 Punkte 5, 8 bis 11 | [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/in-progress/slice-cr-devcontainer-instanz-zusage-getestet.md) · „*Wird bei der Closure gefüllt.*“ | ja — `make gates` (docs-check) bestätigt die Verweise, nicht die Vollständigkeit | DoD-Rest vor Closure |
| V-2 | INFO | Die Plan-Tabelle des CR führt `spec/spezifikation.md` und `internal/**` als „unverändert“; `HEAD` zeigt Änderungen aus Änderung A und aus der Fix-Runde `1bf451c4` (`devcontainer_shell.go`, Spezifikation). Für den CR-Inhalt stimmt die Zeile, für den Baum nicht; der Abgleich nach der letzten Fix-Runde (`AGENTS.md` §3.10) gilt für beide Slices zugleich. | `AGENTS.md` §3.10 (Abgleich nach der letzten Fix-Runde) | Slice §3 · „`spec/spezifikation.md` … unverändert“ | nein — Verifier-only (`git show --stat 1bf451c4`) | Plan-Zeile „unverändert“ gilt nur für den CR-Anteil |
| V-3 | INFO | Die Zusage von `shell` auf der Engine 29.9.2 hat keinen Lauf auf einer solchen Engine im Verifier-Beleg: Das Job-Log nennt `docker-ce-cli 29.9.2` (Client), nicht die Engine des Runners; der Fall „Sitzung schreibt die Zeile selbst“ und die Form 28.0.4 stehen nur als Lesung und Fake-Test. Das ist das Risiko F-7 des Reviews, im Slice als „weiter offen“ geführt; es betrifft Änderung A, nicht den Lastenheft-Wortlaut des CR. | `AGENTS.md` §3.10 (Aussage über Verhalten nennt, was gelesen oder gemessen wurde) | `spec/spezifikation.md` · „Das ist eine Lesung des Codes gegen die gemessene Form, kein gefahrener Lauf des Launchers“ | nein — Verifier-only (`make test-docker` nicht gefahren, Entscheidung des Projektinhabers) | Zusage-Menge ohne Lauf auf der zugesagten Engine |
| V-4 | INFO | Die Beleg-Tabelle nennt für „Linux `amd64`, lokale Docker-Engine“ keinen Test; die Bedingung stützt sich auf die Zusage-Menge der Anforderung. §7 sagt das offen. Die Bedingung steht seit `0.25.0` im Lastenheft und wird nicht neu behauptet. | `AGENTS.md` §3.10 (Bedingung in Spec bekommt Test oder wird gestrichen) | Slice §7 · „kein Test pinnt die Architektur“ | nein — Verifier-only | Bedingung ohne pinnenden Test (offen benannt) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Lastenheft-Wortlaut Satz 2 gegen die Freigabe in §7 | geprüft, ohne Befund (maschineller Teilstring-Vergleich wahr) |
| Version (Zeile 9) und Historie | geprüft, ohne Befund |
| Decken-Regel in `spec/` | geprüft, ohne Befund |
| Handbuch 4.9, `zusagen-und-messstaende.md` | geprüft, ohne Befund (Chronik-Satz „Bis Lastenheft 0.25.0“ ist Review F-8, belassen) |
| [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) | geprüft, im Diff unberührt |
| `CHANGELOG.md`, `internal/e2e` seit CI-Kopf `a6264402` | geprüft, unverändert |
| Gate-Läufe `make gates`, `make doc-trace` | gelaufen, grün, 0 Waisen |
| `make test-docker`, `make verify-depguard` | nicht gefahren, begründet (oben) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** DoD-Rest vor Closure · Plan-Zeile „unverändert“ gilt nur für den CR-Anteil · Zusage-Menge ohne Lauf auf der zugesagten Engine · Bedingung ohne pinnenden Test (offen benannt)

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM, kein LOW. Der Lastenheft-Wortlaut entspricht der Freigabe, die Decken-Regel ist gewahrt, `make gates` und `make doc-trace` sind grün, der CI-Beleg für „zwei Instanzen“ ist im Job-Log gelesen.

**Schließt der Slice?** Noch nicht, aber nichts Inhaltliches steht im Weg. Offen sind allein die Closure-Pflichten: DoD-Haken, §7 (Läufe, Lerneintrag, Risiko-Ausgänge, drei Paarungen), Register-Beleg für `BEO-DOK/zusage-weiter-als-messung` und die Verschiebung nach `done/` (erst Inhalt, dann reiner `git mv`, `AGENTS.md` §3.3). Die Ausgänge der Risiken trägt der Planner; ein Risiko darf „weiter offen“ bleiben (F-7 des Reviews, Anker vorhanden), die übrigen sind nach obiger Zuordnung entfallen.

**Übergabe:** an den Planner zur Closure. Die Finding-Klassen gehen in die Closure §7 und den Zähler. Der Verifier hat nichts repariert und nichts geändert außer dieser Datei.
