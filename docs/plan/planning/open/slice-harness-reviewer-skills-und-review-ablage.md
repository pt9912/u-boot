# Slice Harness: Reviewer-Skill-Dateien und Review-Report-Ablage (FS-1/FS-2)

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** Konformitätslücke der Regelwerk-Adoption, dokumentiert als FS-1/FS-2
in [`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md)
§7. Kein `LH`-/`ADR`-Neubezug auf Anforderungsebene: reine Harness-Konformität.
Die Doku-Mindeststruktur
([`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur))
wird **additiv** ergänzt, nicht geändert — sie ist als Mindest-, nicht als
Maximalstruktur formuliert; deshalb kein Change Request am Vertrags-Stratum.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die beiden Regelwerk-Pflichtartefakte nachziehen, die die Kern-Adoption bewusst
ausgelassen hat: **maschinenlesbare Reviewer-Skill-Dateien** unter
`.harness/skills/` (statt nur Review-Prosa unter `harness/`) und die
**Report-Ablage** `docs/reviews/` mit ihrer Ein-Report-pro-Lauf-Konvention.

Die Skills **duplizieren** die bestehende Prosa nicht, sie **destillieren** sie
in die vom Regelwerk geforderte Skill-Form (Kontext-Eingang, Klassifikation mit
repo-spezifischen Ankern, Nicht-Zuständigkeiten, Output-Schema, Steering-Loop).
[`harness/review.md`](../../../../harness/review.md) und
[`harness/verification.md`](../../../../harness/verification.md) bleiben die
kanonischen Quellen; die Skills verweisen aufwärts auf sie.

## 2. Definition of Done

- [ ] **`.harness/skills/reviewer.md`** angelegt (aus der vendorten Vorlage
  `.harness/baseline/v3.5.1/templates/.harness/skills/reviewer.template.md`
  kopiert und ausgefüllt, Template-Hinweis-Block gelöscht) mit **mindestens zwei
  repo-spezifischen HIGH-Regeln**, die ein generischer Reviewer-Skill nicht
  abdeckt.
- [ ] **`.harness/skills/closure-note-reviewer.md`** angelegt, angepasst an
  u-boot: es gibt **kein** computational Closure-Note-Gate
  (`check_closure_notes.py` existiert hier nicht) — die Abgrenzung „semantische
  Schicht über einem Struktur-Gate" wird durch die u-boot-Realität ersetzt
  (Struktur-Pflichtfelder aus `harness/verification.md`, inferentielle Prüfung
  von Substanz vs. Floskel). Der Skill benennt diese Abweichung explizit.
- [ ] **`docs/reviews/README.md`** angelegt: Zweck, Dateiname-Konvention
  `<YYYY-MM-DD>-<slice-oder-diff-ref>.md`, Regel „ein Report pro Lauf,
  Folgeläufe als neue Datei statt Überschreibung", Verweis auf das vendorte
  Report-Gerüst als Kopiervorlage (keine Kopie im Repo → keine Drift-Quelle).
- [ ] **Adaptions-Ledger-Eintrag** `MR-009` in
  [`harness/conventions.md`](../../../../harness/conventions.md): Ortswahl
  `.harness/skills/` + `docs/reviews/`, Abgrenzung zur Prosa unter `harness/`,
  Begründung, warum `LH-FA-PROJDOCS-001` dadurch nicht verletzt ist.
- [ ] **Erster formaler Lauf** des Reviewer-Skills auf dem Diff dieser Welle,
  abgelegt als Report unter `docs/reviews/` — damit ist `docs/reviews/` nicht
  nur ein leeres Verzeichnis mit README (FS-2 fordert genau diesen
  Entstehungs-Zeitpunkt). Rollentrennung: der Lauf erfolgt in Frischkontext
  gegen den Diff, nicht aus dem Autoren-Kontext heraus.
- [ ] `make docs-check` grün — insbesondere: die neuen Dateien unter
  `.harness/skills/` liegen **außerhalb** des `scan.ignore`-Globs
  `.harness/baseline/**` (so in `MR-005` zugesagt) und werden damit tatsächlich
  auf Link-/Anker-/ID-Pflicht geprüft.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/skills/reviewer.md` | neu | Regelwerk fordert per-Repo Skill-Datei, nicht nur Prosa (FS-1) |
| `.harness/skills/closure-note-reviewer.md` | neu | Schwester-Skill für die engere Closure-Prüfung (FS-1) |
| `docs/reviews/README.md` | neu | Ablage + Negativbefund-Disziplin (FS-2) |
| `docs/reviews/<datum>-<ref>.md` | neu | erster formaler Lauf, macht die Ablage real (FS-2) |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | `MR-009` Ortswahl/Abgrenzung; FS-1/FS-2-Verweise in `MR-005` einlösen |

## 4. Trigger

Direkt fortsetzbar, kein externer Trigger — erster Slice der Welle
`welle-harness-konformitaet-nachlauf`. Rückführung nach `open/`, falls die
Skill-Schärfung eine Vertrags-Kollision aufdeckt (dann CR wie bei `MR-008`).

## 5. Closure-Trigger

DoD vollständig, `make docs-check` grün, erster Review-Report abgelegt,
Closure-Notiz mit Verification Evidence geschrieben.

## 6. Risiken und offene Punkte

- **Skill-Drift gegen Prosa:** Zwei Orte für Review-Regeln (`harness/review.md`
  als Prosa, `.harness/skills/reviewer.md` als Skill) können auseinanderlaufen.
  Gegenmaßnahme: Der Skill nennt die Prosa als kanonische Quelle und
  dupliziert Tabellen nicht, sondern verweist; die Anker-Listen sind bewusst
  die *scharfe Kurzform*, nicht die Vollkopie.
- **d-check-Kennungspflicht in neuen Verzeichnissen:** `.harness/skills/` und
  `docs/reviews/` werden erstmals gescannt; jede `LH-*`-, `ADR-NNNN`- und
  `slice-*`-Kennung dort braucht einen Link (oder ein bewusstes
  `d-check:ignore`). Erwartet, kein Carveout — Diagnose vor Carveout.
- **Rollentrennung beim ersten Lauf:** Autor und Reviewer dürfen nicht denselben
  Eingabe-Kontext teilen (`AGENTS.md` §Role Separation). Der Report muss den
  verwendeten Kontext offenlegen, sonst ist er nicht reproduzierbar.
- **Kein Carveout erwartet:** additive Doku-/Harness-Arbeit, kein Code-Delta,
  keine Gate-Lockerung.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührt schreibend nur die Doku-/Harness-Sub-Areas; alle berührten Sub-Areas
sind **GF** (Doku-führt), Einordnung nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Bestandscode wird nicht angefasst.
