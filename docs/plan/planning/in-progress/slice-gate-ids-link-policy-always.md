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

- [ ] **Politik-Entscheidung dokumentiert:** Für jede Artefaktklasse ist
  festgehalten, ob sie verlinkt oder per `exempt-paths` ausgenommen wird — mit
  Begründung, nicht als Liste. Mindestens zu entscheiden:
  `docs/reviews/**` (Audit-Artefakte, ein Report pro Lauf, keine
  Überschreibung), `docs/plan/planning/done/**` (nur korrigierend änderbar),
  `.harness/baseline/**` (bereits per `scan.ignore` draußen).
- [ ] **Muster-granular statt global:** `link-policy` wird je Muster gesetzt.
  Zu prüfen ist, ob alle vier Muster (`ADR-*`, `LH-*`, `PH|TC|CO-*`,
  `slice|tranche-*`) dieselbe Politik brauchen — das `slice`-Muster hat einen
  Sonderfall (s. §6).
- [ ] **Verbleibende Befunde auf null**, ohne dass eine Ausnahme eine
  behebbare Stelle verdeckt: Ein `exempt-paths`-Eintrag ist nur für
  Artefaktklassen zulässig, deren Unveränderlichkeit eine **Regel** ist, nicht
  für Bequemlichkeit.
- [ ] **`repair`-Weg geprüft:** `make doc-repair` (bzw. `--repair`) erzeugt
  einen git-apply-fähigen Patch. Ob er die 78 README-Fälle trägt, ist zu
  messen — von Hand verlinken war beim letzten Mal die teuerste Position.
- [ ] **`MR-005` nachgezogen** (Gate-Haltung inkl. Politik).
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) | update | `link-policy` je Muster, `exempt-paths` je Politik |
| `internal/**/README.md`, `harness/*.md`, Spec-Dateien | update | verbleibende Kennungen verlinken |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Politik dokumentieren |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump. Reihenfolge
innerhalb der Welle: **nach**
[`slice-gate-planning-targets-module`](../done/slice-gate-planning-targets-module.md),
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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen*, *spec / docs* und
`internal/**/README.md` — alle **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration (die READMEs sind seit dem Retrofit graduiert).
