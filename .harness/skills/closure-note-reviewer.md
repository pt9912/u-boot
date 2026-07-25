# Closure-Note-Reviewer-Skill — u-boot

* Status: Accepted
* Bezug: [`harness/verification.md`](../../harness/verification.md)
  (Pflichtfelder der Verification Evidence),
  [`docs/plan/planning/README.md`](../../docs/plan/planning/README.md) §Lifecycle,
  [`LH-FA-PROJDOCS-003`](../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
* Gilt für: die Closure-Abschnitte der Slices in
  [`docs/plan/planning/done/`](../../docs/plan/planning/done/)

**Adaption gegenüber der Baseline-Vorlage.** Die vendorte Vorlage beschreibt
diesen Skill als *semantische* Schicht über einem *computational* Struktur-Gate
(`make verify-closure-notes` / `tools/check_closure_notes.py`). **u-boot hat
dieses Gate nicht** — `make docs-check` prüft Referenzen, nicht
Closure-Substanz. Konsequenz: Dieser Skill trägt **beide** Ebenen, Struktur und
Inhalt, und darf sich nicht auf ein vorgelagertes Gate berufen. Die
Struktur-Pflichtfelder stammen aus
[`harness/verification.md`](../../harness/verification.md) §Pflicht bei
Slice-Closure, nicht aus einem Skript.

## Kontext-Eingang (Pflicht)

Was der Reviewer *immer* mitbringt, bevor er urteilt:

- die Closure-/Evidence-Abschnitte der zu prüfenden Slices in
  [`docs/plan/planning/done/`](../../docs/plan/planning/done/)
- [`harness/verification.md`](../../harness/verification.md) — die sieben
  Pflichtfelder (Scope, DoD-Abgleich, Sensors, Traceability, Carveouts, Nicht
  ausgeführt, Commit/Artefakt) und die harten Regeln
- der zugehörige Slice-Plan (DoD-Liste — wogegen der Abgleich läuft)
- das Carveout-Inventar
  ([`carveouts.md`](../../docs/plan/planning/in-progress/carveouts.md)) für die
  Gegenprobe der Carveout-Zeile
- die tatsächliche Commit-Historie des Slice (belegt die Hash-Angabe)

Ohne diese Liste prüft der Reviewer Text, nicht *die Pflicht-Inhalte, die eine
Closure tragen muss*.

## Prüf-Auftrag

Lies die Closure-Notiz und die Verification Evidence jedes Slice in `done/`.
Markiere jede, die *keinen* der folgenden Inhalte trägt: (a) ein konkretes
Lernsignal mit Ursache („Test rot, **weil** X"), (b) ein konkretes Folge-Slice
oder die belegte Aussage, dass keines entsteht, (c) eine konkrete Architektur-
oder Harness-Beobachtung. Floskeln ohne Inhalt sind ein HIGH-Finding.

Der Inhalt-Teil ist inferentiell — „Substanz vs. Floskel" ist semantisch und
nicht regex-fähig. Der Struktur-Teil (Feld vorhanden? Sensor benannt? Hash
statt Platzhalter?) ist beobachtbar und wird deshalb *zuerst* geprüft, damit
nicht Substanz diskutiert wird, wo schon das Feld fehlt.

## Klassifikation

**HIGH** — eines der folgenden:

- **Floskel ohne Substanz.** Die Closure-Notiz existiert, trägt aber *keinen*
  der drei Pflicht-Inhalte („wie geplant umgesetzt", „lief gut", „fertig").
- **Pflichtfeld fehlt ganz.** Eines der sieben Felder aus
  [`harness/verification.md`](../../harness/verification.md) ist nicht
  vorhanden — insbesondere ein fehlendes „Nicht ausgeführt" (dann ist nicht
  unterscheidbar, ob ein Sensor bewusst ausgelassen oder vergessen wurde).
- **Hash-Platzhalter statt Delivery-Hash (repo-spezifisch #1).** Ein
  `done/`-Slice trägt im Commit-/Artefakt-Feld einen Sucher-Platzhalter
  (`git log --grep …`, „siehe Historie", „TBD") statt des konkreten
  Commit-Hashes. Ein Slice in `done/` ohne auflösbaren Hash ist nicht
  auditierbar; der Hash steht direkt in der Zeile.
- **Sensor behauptet, nicht belegt (repo-spezifisch #2).** Die Sensors-Tabelle
  meldet `pass` ohne kurzen Beleg (Zahl, Ausgabe-Zeile, Schwellwert) — z. B.
  `make coverage-gate: pass` ohne Prozentwert oder `make docs-check: pass` ohne
  Datei-/Befundzahl. Ein grüner Sensor ohne Beleg ist eine Behauptung, und
  [`harness/verification.md`](../../harness/verification.md) verlangt Evidence
  je Sensor.
- **Carveout-Zeile fehlt oder widerspricht dem Inventar.** Neuer Carveout ohne
  Eintrag in
  [`carveouts.md`](../../docs/plan/planning/in-progress/carveouts.md)
  ([`LH-FA-PROJDOCS-005`](../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)).

**MEDIUM** — genau *einer* der drei Pflicht-Inhalte fehlt oder ist unkonkret:

- Lernsignal ohne das „weil X" (Behauptung statt Ursache).
- Folge-Slice benannt, aber ohne zugehörigen `open/`-Eintrag.
- Architektur-/Harness-Beobachtung als Etikett statt als beobachtbare Aussage.
- DoD-Abgleich pauschal („alle Punkte erfüllt") bei einer DoD-Liste, deren
  Punkte nicht offensichtlich gemeinsam belegt sind.

**LOW** — alle Inhalte vorhanden, aber schwer nachvollziehbar formuliert
(Substanz da, Klarheit fehlt); inkonsistente Feld-Reihenfolge gegenüber dem
Evidence-Block.

**INFO** — Hinweis ohne erwartete Aktion, z. B. „verweist auf ein Folge-Slice,
das noch nicht in `open/` liegt — Nachtrag durch die Planning-Rolle".

## Was dieser Skill NICHT macht

- Keine Bewertung, ob der Slice *fachlich* korrekt abgeschlossen wurde — das
  prüft der Verifier gegen Spec und Sensors, nicht dieser Skill gegen Text.
- Keine Umschreibung der Closure-Notiz — der Autor formuliert nach, der
  Reviewer kategorisiert.
- Keine Prüfung von Slices außerhalb `done/` — `open/`, `next/` und
  `in-progress/` tragen noch keine Closure-Pflicht.
- Kein Nach-Ausführen von Sensors. Fehlt ein Beleg, ist das ein Finding, keine
  Aufforderung an den Reviewer, den Sensor selbst laufen zu lassen.

## Output-Schema

Jedes Finding:

- `kategorie`: HIGH | MEDIUM | LOW | INFO
- `quelle`: `Verification-Pflichtfeld <Name>` | `Closure-Inhaltspflicht (a/b/c)`
  | `LH-*`-ID
- `pfad`: `docs/plan/planning/done/<slice>.md`:<Zeile>
- `befund`: *welches* Feld bzw. *welcher* Pflicht-Inhalt fehlt, 1–2
  beobachtbare Sätze, ohne Formulierungs-Vorschlag
- `verifizierbar`: Struktur-Findings ja (am Dokument nachlesbar),
  Substanz-Findings nein (inferentiell)

Zusätzlich am Ende eine Zeile „geprüft, ohne Befund" pro betrachteter
Slice-Charge — sonst ist die Abdeckung unsichtbar.

## Pflege (Steering-Loop)

Bei dreimaligem HIGH derselben Art:

- Muster in [`harness/verification.md`](../../harness/verification.md) als
  benanntes Anti-Pattern aufnehmen.
- Prüfen, ob das Muster *strukturell* fangbar wäre (Feld-Präsenz,
  Hash-Format, Floskel-Liste) — dann als computational Gate bauen statt
  weiter inferentiell nachzulesen. Ein solches Gate wäre für u-boot neu und
  gehört als eigener Slice geplant, nicht nebenbei eingezogen.
- Slice-Template bzw. Evidence-Block schärfen, falls das Feld die
  Pflicht-Inhalte nicht klar genug abfragt.

Diese Datei wird nicht überschrieben, sondern versioniert fortgeschrieben.
