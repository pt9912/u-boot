# Slice Gate: `ids` auf `link-policy: always` — Kennungen auch in Code-Spans

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-gate-ausbau-v0.51` (s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(Kennungs-/Link-Disziplin) und
[`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md). `MR-005` führt
die Gate-Haltung.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Den Blindfleck schließen, der beim `internal/`-README-Retrofit aufgefallen
ist: `ids` prüft im Default (`prose`) nur *nackte* Kennungen. Steht eine ID in
Backticks, gilt sie als Code-Beispiel und bleibt ungeprüft — und genau so
liegt der Großteil der Kennungen in u-boots Doku. `link-policy: always` prüft
Inline-Code-Spans zusätzlich.

**Gemessen: 100 Befunde.** Die Verteilung ist der eigentliche Inhalt dieses
Slice — 78 in den `internal/`-READMEs (mechanisch behebbar), 13 im
Review-Report und 15 in `done/`-Slices (**nicht** mechanisch behebbar, s. §6).

## 2. Definition of Done

- [x] **Politik-Entscheidung dokumentiert:** Für jede Artefaktklasse ist
  festgehalten, ob sie verlinkt oder per `exempt-paths` ausgenommen wird — mit
  Begründung, nicht als Liste. Mindestens zu entscheiden:
  `docs/reviews/**` (Audit-Artefakte, ein Report pro Lauf, keine
  Überschreibung), `docs/plan/planning/done/**` (nur korrigierend änderbar),
  `.harness/baseline/**` (bereits per `scan.ignore` draußen).
- [x] **Muster-granular statt global:** `link-policy` wird je Muster gesetzt.
  Zu prüfen ist, ob alle vier Muster (`ADR-*`, `LH-*`, `PH|TC|CO-*`,
  `slice|tranche-*`) dieselbe Politik brauchen — das `slice`-Muster hat einen
  Sonderfall (s. §6).
- [x] **Verbleibende Befunde auf null**, ohne dass eine Ausnahme eine
  behebbare Stelle verdeckt: Ein `exempt-paths`-Eintrag ist nur für
  Artefaktklassen zulässig, deren Unveränderlichkeit eine **Regel** ist, nicht
  für Bequemlichkeit.
- [x] **`repair`-Weg geprüft:** `make doc-repair` (bzw. `--repair`) erzeugt
  einen git-apply-fähigen Patch. Ob er die 78 README-Fälle trägt, ist zu
  messen — von Hand verlinken war beim letzten Mal die teuerste Position.
- [x] **`MR-005` nachgezogen** (Gate-Haltung inkl. Politik).
- [x] `make docs-check` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) | update | `link-policy` je Muster, `exempt-paths` je Politik |
| `internal/**/README.md`, `harness/*.md`, Spec-Dateien | update | verbleibende Kennungen verlinken |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Politik dokumentieren |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump. Reihenfolge
innerhalb der Welle: **nach**
[`slice-gate-planning-targets-module`](slice-gate-planning-targets-module.md),
weil dieser Slice den breitesten Diff hat.

## 5. Closure-Trigger

Politik entschieden und dokumentiert, Befunde auf null, `MR-005` nachgezogen,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Der Kern-Konflikt:** `docs/reviews/` und `done/`-Slices sind per
  Konvention unveränderlich. Ihre Kennungen nachträglich zu verlinken
  widerspräche der eigenen Regel; sie ungeprüft zu lassen ist eine
  Abdeckungslücke. `exempt-paths` löst das mechanisch — die Frage ist, ob es
  das *richtig* löst. Alternative: Politik gilt ab jetzt, Bestand bleibt
  ausgenommen, und die Ausnahme trägt ein Ablaufdatum statt „permanent".
- **Versionsnummern in Slice-Namen:** Das `slice`-Muster bricht am Punkt —
  `slice-harness-baseline-bump-review-v3.5.2` erzeugt einen Treffer auf
  `slice-harness-baseline-bump-review-v3`, der nicht sauber verlinkbar ist.
  Entweder Muster schärfen, `exempt`-Behandlung, oder künftige Slice-Namen
  ohne Versionspunkte (Namenskonvention). Drei Befunde hängen daran.
- **Ausnahme als Gewohnheit:** Jeder `exempt-paths`-Eintrag schwächt das Gate.
  Wer beim ersten Widerstand ausnimmt, hat am Ende ein grünes Gate ohne
  Aussage — dieselbe Falle wie beim `internal/**`-Glob, den wir gerade erst
  aufgelöst haben.
- **Kein Carveout erwartet** — aber wenn doch, dann inventarisiert und mit
  Plan-Anker.

## 7. Closure-Notiz (nach `done/`)

### Die Plan-Prämisse war falsch — und das war der Ertrag

Der Plan ging davon aus, die 13 Report- und 15 `done/`-Befunde seien „**nicht**
mechanisch behebbar", und stellte die Frage, ob `exempt-paths` sie *richtig*
löst. Beim Nachlesen der eigenen Regeln stellte sich heraus: Die Prämisse hält
nicht.

- [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  nennt **„Querverweise"** ausdrücklich als zulässige nachträgliche Korrektur in
  `done/`. Eine fehlende Kennung zu verlinken *ist* eine — die Regel verbietet
  substanzielle Änderungen, nicht Referenz-Pflege.
- Die `docs/reviews/`-Regel lautet „ein Report pro Lauf, Folgeläufe als **neue
  Datei** statt Überschreibung". Sie verbietet, einen Report durch einen neuen
  Lauf zu ersetzen — nicht, einen Link zu setzen.

Damit war die Ausnahme, um die sich der ganze Plan drehte, gar nicht nötig.

### Drei Klassen, drei Mechanismen — kein einziger Glob

Die 100 Befunde waren nicht ein Problem, sondern drei:

| Klasse | Fälle | Antwort |
|---|---|---|
| Echte Referenz | 91 | verlinkt |
| Beleg, kein Verweis | 6 | `d-check:ignore` **je Zeile**, Begründung an Ort und Stelle |
| Muster-Artefakt | 3 | Regex geschärft statt ausgenommen |

**Kein `exempt-paths`, nirgends.** Genau das war die Sorge, die den Slice
ausgelöst hat: Ein Verzeichnis-Glob hätte den Bestand *und* jedes künftige
Dokument ausgenommen — der nächste Review-Report und jeder neue `done/`-Slice
wären ungeprüft entstanden, ohne dass es jemandem aufgefallen wäre. Dieselbe
Falle wie beim `internal/**`-Glob, den wir am selben Tag aufgelöst haben.

Die sechs Zeilenmarker sitzen dort, wo eine Kennung **Daten** sind statt einer
Referenz: ein `code: "LH-FA-DEV-003"`-JSON-Payload in einem `done/`-Slice und  <!-- d-check:ignore (zitiertes Payload-Beispiel, keine Referenz) -->
fünf `pfad:`-Felder des Review-Reports, die den Fundort *zum Prüfzeitpunkt*
festhalten. Diese auf `done/` umzubiegen hätte den Befund verfälscht — F-7 des
Reports lautete gerade, dass die Slices damals in `open/` lagen.

Die Bereichs-Schreibweisen (`LH-FA-INIT-001..007`, 30 Stück) sind zu verlinkten  <!-- d-check:ignore (Notations-Beispiel, keine Referenz) -->
Paaren aufgelöst — die Form, die
[`spec/architecture.md`](../../../../spec/architecture.md) längst verwendet.

### Der Report wurde angefasst — und das ist vertretbar

Entscheidung des Projektinhabers, hier begründet: Acht echte Referenzen sind
verlinkt, fünf `pfad:`-Zeilen tragen einen Marker. **Kein Befund, keine
Kategorie, kein Pfad-Wert und kein Verdikt hat sich geändert** — die
Auditierbarkeit hängt an diesen, nicht an der Markdown-Syntax. Die Alternative
(namentlicher Datei-Ausschluss) hätte den Report unberührt gelassen, aber auch
seine acht echten Referenzen dauerhaft ungeprüft.

### Verification Evidence

Scope:
- Slice: `slice-gate-ids-link-policy-always`
- IDs: **keine** Anforderung geändert. Tragend:
  [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
  (Kennungs-Disziplin),
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  (erlaubt die Querverweis-Korrektur),
  [`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md) (lesend).
- Artefakte: [`.d-check.yml`](../../../../.d-check.yml) (`link-policy: always`
  je Muster, geschärftes `slice`-Muster), 18 Dateien mit gesetzten Links bzw.
  Zeilenmarkern, [`harness/conventions.md`](../../../../harness/conventions.md)
  (`MR-005` Linkpolitik).

DoD-Abgleich: alle Punkte erfüllt — mit einer Abweichung nach oben: Der Punkt
„Politik-Entscheidung dokumentiert" verlangte eine Entscheidung *je
Artefaktklasse*; das Ergebnis ist für alle drei Klassen dieselbe Antwort
(**keine Ausnahme**), was der Plan nicht erwartet hatte. Der `repair`-Weg wurde
**nicht** genutzt: Die Kennungs-Ziele stammen aus einer selbst gebauten
Anker-Karte (aus `spec/lastenheft.md`-Überschriften, dem ADR-Verzeichnis und den
Lifecycle-Verzeichnissen erzeugt), weil sie Bereichsauflösung und
Beleg-Klassifikation beherrschen musste — beides kann ein generischer
Reparatur-Patch nicht wissen.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make doc-check` (nach Umstellung, vor Behebung) | 100 Befunde | Ausgangs-Ist, nach Verzeichnis und Klasse aufgeschlüsselt |
| Anker-Karte gegen Bestand | pass | 137 `LH-*`, 13 ADRs, 92 Slices; Stichprobe gegen bekannte Anker (`#lh-fa-cli-006--exit-codes`, `#lh-fa-init-005--überschreibschutz`) deckungsgleich |
| `make docs-check` (final) | pass | **144 Dateien / 0 Befunde** mit `link-policy: always` |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) | Kennungen sind jetzt auch in Code-Spans linkpflichtig — die Abdeckungslücke aus dem README-Retrofit ist geschlossen |
| [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) | Die `done/`-Änderungen sind Querverweis-Korrekturen, ausdrücklich zulässig |
| `MR-005` | Linkpolitik samt der drei Mechanismen dokumentiert |

Carveouts: Neu: none. Gelöst: none. Unverändert: none. **Ausdrücklich:** Die
sechs Zeilenmarker sind keine Carveouts — sie erklären, dass an dieser Stelle
keine Referenz vorliegt, sie setzen keine Prüfung aus.

Nicht ausgeführt:
- `make gates` / `make test` — kein Go-Delta.

Independent Review: nicht durchgeführt. Die Behauptung „alle Kennungen
verlinkt" ist maschinell geprüft; die inhaltliche Frage — ist eine Kennung
Referenz oder Beleg? — ist für die sechs Grenzfälle einzeln im Diff
nachvollziehbar begründet.

Commit / Artefakt: `<wird beim Closure-Commit eingetragen>`

### Steering-Loop-Lerneintrag

- **Bevor man eine Ausnahme baut, liest man die Regel, gegen die sie schützen
  soll.** Der ganze Slice war um die Frage herum geplant, wie man `done/` und
  `docs/reviews/` ausnimmt — und die Antwort stand in
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle): Querverweis-Korrekturen sind dort erlaubt. Der Plan
  hatte eine Unveränderlichkeit angenommen, die die Regel so nie behauptet hat.
- **„Nicht mechanisch behebbar" war eine Vermutung, keine Messung.** Sie stand
  im Plan als Tatsache. Erst das Auflisten der Fundstellen mit Kontext hat
  gezeigt, dass 91 von 100 schlichte Referenzen sind. Aufwandsschätzungen
  gehören gemessen, bevor sie eine Architektur-Entscheidung tragen.
- **Ein Glob ist die bequemste Art, ein Gate zu entwerten.** Er wirkt sofort,
  sieht sauber aus und nimmt lautlos auch alles Künftige mit. Die drei
  zeilengenauen Mechanismen waren teurer und sind die einzigen, die den
  nächsten Report und den nächsten `done/`-Slice noch prüfen.
- **Folge-Slices:** keine. Offen bleibt eine **Namenskonvention**: Slice-Namen
  mit Versionspunkten (`…-v3.5.2`) haben das Muster gebrochen und mussten es
  erweitern. Künftige Namen kommen besser ohne aus — das ist eine Notiz für den
  nächsten Planning-Slice, kein eigener.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen*, *spec / docs* und
`internal/**/README.md` — alle **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration (die READMEs sind seit dem Retrofit graduiert).
