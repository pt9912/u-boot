# Slice V2: `devcontainer remove` — Spezifikation: Aufruf, Suche, Abfrage, Ausgabe und Exit-Codes von `remove`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-devcontainer-remove`](../welle-devcontainer-remove.md), Release-Zuordnung offen. Dritter von vier Slices; ihm geht [`slice-cr-devcontainer-remove`](../done/slice-cr-devcontainer-remove.md) voraus. Reihe (vier Slices, je einzeln lieferbar): [`slice-v1-devcontainer-remove-messung`](../done/slice-v1-devcontainer-remove-messung.md) → [`slice-cr-devcontainer-remove`](../done/slice-cr-devcontainer-remove.md) → [`slice-v2-spec-devcontainer-remove`](slice-v2-spec-devcontainer-remove.md) → [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Scope), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes), [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Aufruf, Suche, Ausgabe und Exit-Codes von `up` und `down`; `remove` erweitert sie) und der Abschnitt zur `--json`-Ausgabe, in dem `devcontainer` als Wert von `command` steht; ob eine neue `SPEC-<NNN>` entsteht oder [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) wächst, entscheidet der Architect. Außerdem die Kopien der Befehlsliste in `spec/spezifikation.md` (`--quiet`-Liste, Aussage zu Instanzen, Wert `devcontainer` des Feldes `command`, Verweisliste zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)), in `spec/architecture.md` (Abschnitt zum CLI-Adapter und die Tabelle der Fehlerklassen) und in `docs/maintainer/cli-json-contract.md` (Abschnitt 6.11, Matrix der Dateisystem-Mutationen); Träger je Stelle in §3.

**Verantwortlich:** —

**Autor:** Claude (Planner-Entwurf auf Wunsch des Projektinhabers). **Datum:** 2026-11-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Spezifikation legt fest, wie `u-boot devcontainer remove` den Container sucht (wie `down`: Label `devcontainer.local_folder`, Instanz-Label), welche Volumes `--volumes` erfasst, wie die Abfrage und `--yes` wirken, was die Ausgabe und die `--json`-Ausgabe tragen und welche Exit-Codes gelten. Sie legt dazu die Grenzen fest, die der Lastenheft-Wortlaut bewusst nicht trägt (Review des CR F-8, Entscheidung F5 des CR): was `remove --volumes` mit einem Volume tut, das ein anderer Container hält (die Messung belegt nur, dass das Entfernen scheitert, solange ein Container es hält); wie anonyme Volumes behandelt werden (die Messung: `docker rm` lässt sie stehen, `docker rm -v` entfernt sie; der Lastenheft-Satz „Die Volumes des Containers bleiben dabei erhalten“ trägt eine der beiden Wahlen); die Reihenfolge Abfrage, Stoppen, Entfernen (das Lastenheft verlangt die Abfrage vor jedem Eingriff); die Zuordnung der Exit-Codes `10` (Freigabe fehlt oder abgelehnt, nur Instanzen), `11` (Engine nicht erreichbar oder nicht lokal) und `12` (Stoppen oder Entfernen scheitert, die Engine verweigert); Plattform- und Messgrenzen (die Messung lief auf der verschachtelten Engine).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft** — liegt in [`slice-cr-devcontainer-remove`](../done/slice-cr-devcontainer-remove.md); die Spezifikation präzisiert es, erweitert es nie.
- **Code** — Folge-Slice [`slice-v1-devcontainer-remove`](../open/slice-v1-devcontainer-remove.md).

## 2. Definition of Done

- [ ] **Spezifikation:** Aufruf, Suche, Volume-Auswahl, Abfrage, Ausgabe (Text und `--json`), Exit-Codes und Grenzen stehen in `spec/spezifikation.md`; jede Bedingung, die den Code beschreibt („nur wenn“, „immer“), bekommt im Code-Slice einen pinnenden Test oder entfällt (`AGENTS.md` §3.10). Die Kopien der Befehlsliste (Plan-Tabelle §3) nennen `remove`, und keine Aussage der Spezifikation widerspricht dem Lastenheft (insbesondere „räumt keine auf“ in `spec/spezifikation.md`; das Lastenheft sagt „räumt keine von sich aus auf“).

Pflichten der Closure (zählen nicht mit):

- [ ] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Verification-Report unter `docs/reviews/`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) bzw. eine neue `SPEC-<NNN>` (gelesen am Stand `7f5110cc`): Zeile 866 (*Aufruf*: nennt `up`, `down`, `shell`), Zeile 870 („`u-boot` führt keine Liste der Instanzen und räumt keine auf“ — widerspricht dem Lastenheft „räumt keine von sich aus auf“), Absatz *Ohne `--instance`, nur Instanzen* und *Mit `--instance`* (Zeilen 875 und 876: nennen `up`, `down`, `shell`), die Tabelle der Exit-Codes ab Zeile 938 (Spalten `up`, `down`, `shell`; eine Spalte oder Zeile für `remove`); Zeile 948 (*Reichweite*, gelesen: „Ohne `docker` im `PATH` enden alle drei Unterbefehle mit Exit-Code `11`“ und der Verweis auf die Anforderung für Hosts, Engines und die Grenze im ausgelieferten Container-Image: `remove` kommt hinzu; der Image-Satz gilt im Lastenheft `0.37.2` ausdrücklich auch für `remove`, die Spezifikation führt dafür einen Test des Umsetzungs-Slice oder weist ihn als Nicht-Zusage aus, `AGENTS.md` §3.10); die Reihenfolge von Containersuche und Freigabe bei `remove --volumes` (Lastenheft `0.37.2`: ohne Container zu entfernen keine Abfrage, Exit-Code `0`; bei `u-boot down --volumes` steht die Freigabe vor dem Engine-Zugriff, `internal/hexagon/application/downservice.go`, `Down`); die Verweisliste Zeile 1602 (Eintrag zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): `u-boot devcontainer up`, `up --recreate`, `down`, `shell`) |
| `spec/spezifikation.md` (`--quiet` und `--json`) | update | Zeile 619 (Liste der ändernden Subkommandos unter `--quiet`: endet auf `devcontainer shell`), Zeile 631 (`devcontainer down --quiet` gibt nichts aus; `remove --quiet` fehlt), Zeile 1354 (Wert `devcontainer` im Feld `command`: „`up` beziehungsweise `down`“; ob `remove` eine `--json`-Ausgabe hat, legt dieser Slice fest) |
| `spec/architecture.md` | update | Zeilen 299 bis 314 (Abschnitt zum CLI-Adapter: „`devcontainer` (mit `up`/`down`/`shell`)“, „jeder der drei Unterbefehle trägt eine eigene Diagnostic-Abbildung“, Abweisung von `--dry-run` und `--diff`); Zeile 381 (`ARC-008`: Lifecycle-Operationen, Sitzung von `devcontainer shell`); Zeile 656 (Fehlerklasse „Laufzeitfehler der Runtime“, Exit-Code `12`: Stoppen und Entfernen von `remove`); Zeile 651 (Bestätigungs-Gate, Exit-Code `10`) trägt keinen Befehlsnamen und bleibt unverändert (gelesen). Die Architektur bleibt meilenstein- und ADR-frei (`AGENTS.md` §3.4) |
| `docs/maintainer/cli-json-contract.md` | update | Abschnitt 6.11 (Überschrift Zeile 1244: `devcontainer up --json` / `devcontainer down --json`, `devcontainer shell`; ein Umschlag oder die Abweisung für `remove`) und die Matrix der Dateisystem-Mutationen (Zeilen 1356 bis 1358: `devcontainer up`, `down`, `shell`; eine Zeile `remove`, dessen Wirkung die Container-Engine, nicht eine Projektdatei ist) |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeilen des Slice, Zustandsmarker |

## 4. Trigger

**Start** (`next` → `in-progress`): [`slice-cr-devcontainer-remove`](../done/slice-cr-devcontainer-remove.md) liegt in `done/`; `in-progress/` ist leer.

**Rückführungen:** `in-progress` → `next`, wenn der CR-Wortlaut eine Spec-Stelle nicht trägt.

**Bezug zum CR:** Der Start verlangt das Lastenheft `0.37.3` (Folgeänderungen nach dem Review und den Nachlauf-Reviews des CR: Exit-Codes `10`, `0` und `12`, abgelehnte Abfrage, Abfrage nur bei einem Container zu entfernen, auch in der Regel zu destruktiven Operationen ([`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)), Zusage-Menge samt Image-Satz) freigegeben und in `done/` geschlossen; die Spezifikation präzisiert diesen Wortlaut.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Reports liegen vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Spezifikation beschreibt Verhalten, das der Code noch nicht hat (Zielbild). — **Ausgang:** offen; der Code-Slice pinnt jede Aussage mit einem Test.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `spec/` und `docs/plan/` (DOK); alle stehen in `harness/conventions.md` §Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** vor dem Übergang nach `next/` auszuführen (Register `docs/plan/planning/observations/`); der Entwurf trägt noch kein Ergebnis.

**Modus-Begründungsblock:** alle berührten Sub-Areas sind Greenfield. Kein BF-/Hybrid-Block nötig.
