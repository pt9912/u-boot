# Slice Harness: Auffrischung des Reviewer-Skills (d-check-Bump übernommen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-harness-wartung`](../welle-harness-wartung.md) — geht dem Slice
[`slice-harness-internal-readmes-festlegungen-in-spec`](../in-progress/slice-harness-internal-readmes-festlegungen-in-spec.md) voraus.

**Bezug:** [`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md)
(Reviewer-Skill); Herkunft: Delta-Lektüre in
[`slice-harness-regelwerk-adoption-v6.18.1`](../done/slice-harness-regelwerk-adoption-v6.18.1.md)
§3/§4 und Closure-Notiz von
[`slice-harness-werkzeug-schicht-ai-harness-init`](../done/slice-harness-werkzeug-schicht-ai-harness-init.md).
Keine Anforderungsänderung — Skill-Pflege.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `.harness/skills/reviewer.md` an die Reviewer-Vorlage der Baseline
`v6.18.1` angleichen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Bump von d-check auf `v0.86.1`: den führt
  [`slice-harness-ai-harness-init-adoption`](../done/slice-harness-ai-harness-init-adoption.md)
  aus, weil das Werkzeug `d-check.mk` samt digest-gepinntem Image mitbringt. Der
  Dateiname dieses Slice behält die frühere Fassung seines Gegenstands.
- `closure-note-reviewer.md`: nicht Gegenstand des Delta (nur ein
  Pfad-Hinweis), wird nur mitgezogen, falls der Skill-Abgleich dort denselben
  Pfad-Anker verlangt.

## 2. Definition of Done

- [x] `.harness/skills/reviewer.md` trägt die Vorlagen-Änderungen:
  `pfad` als Datei plus wörtliches Kurzzitat (Zeile nur als Lesehilfe), LOW nur
  mit Konventions-Anker, „Kein Stil-Polizist“, kein HIGH/MEDIUM ohne
  Failure-Szenario, `klasse` im Output-Schema, dazu die HIGH-Anker der Vorlage,
  die das Repo schon trägt (Kommentar-Klassen und Zustandsfelder nach
  `AGENTS.md` §3.7, Norm nur im Template-Kommentar, Inline-`//nolint` nach
  `AGENTS.md` §3.2); `harness/review.md` ist im selben Zug nachgezogen, wo es
  dasselbe Schema führt.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Doku-Update für `MR-009`, falls deren Aussage den neuen Stand berührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/skills/reviewer.md` | update | Angleichung an `reviewer.template.md` der Baseline `v6.18.1` |
| `harness/review.md` | update | Output-Schema `pfad` und `klasse`, Kategorien-Regel, HIGH-Anker und Failure-Szenario-Regel, falls dort geführt |
| `docs/plan/planning/in-progress/roadmap.md` | update | Der Slice in `in-progress/` nimmt der Roadmap den Ruhe-Marker (`make docs-check`) |

Das Werkzeug `ai-harness-init` meldet bei jedem Lauf, dass der Skill von der
mitgelieferten Fassung abweicht und nennt die Vorlage zum Abgleich; es
überschreibt ihn nicht.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Verantwortlichen;
kein Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Abgleich zieht eine Neufassung von
  `harness/review.md` nach sich — diese als eigenen Slice abspalten.
- `in-progress` → `open` (blockiert): die Vorlage widerspricht `harness/review.md`
  und die Quelle ist zu klären.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
mit Steering-Loop-Eintrag geschrieben, `make gates` auf dem Stand vor dem
`git mv` nach `done/` grün.

## 6. Risiken und offene Punkte

- Der Reviewer-Skill ist durch MR-009 repo-gepflegt; die Vorlage ist kein
  Pflichtstand — übernommen wird nur, was `harness/review.md` nicht
  widerspricht — **Ausgang:** entfallen: geprüft, nichts Widersprüchliches
  übernommen; der eine Widerspruch (Failure-Szenario gegen die Kommentar-Anker,
  Review F-1) ist in Skill und `harness/review.md` aufgelöst.

## 7. Closure-Notiz

**Belege:** Review-Report [`2026-10-16-slice-harness-dcheck-bump-reviewer-skill`](../../../reviews/2026-10-16-slice-harness-dcheck-bump-reviewer-skill.md)
(0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO; F-1, F-2, F-4 behoben); Verification Evidence
[`2026-10-16-verification-slice-harness-dcheck-bump-reviewer-skill`](../../../reviews/2026-10-16-verification-slice-harness-dcheck-bump-reviewer-skill.md)
(bestanden); `make gates` grün (Exit 0).

- **Was hat funktioniert:** Der Abgleich ließ sich als Diff gegen `reviewer.template.md` führen; der Reviewer prüfte
  ihn mit dem Skill, wie er nach dem Commit stand.
- **Was ging anders als geplant:** Die Vorlage bringt drei HIGH-Anker mit, die der Plan nicht nannte (Kommentar,
  Zustandsfeld, Norm nur im Template-Kommentar); sie standen durch `AGENTS.md` §3.7 schon als Repo-Regel da und
  kamen mit hinein, die DoD wurde nachgezogen. Die neue Regel „kein HIGH/MEDIUM ohne Failure-Szenario" kollidierte
  mit diesen Ankern (F-1); aufgelöst ist das Szenario dort das Versagen des nächsten Änderers. Der Slice nahm außerdem
  der Roadmap den Ruhe-Marker, weil `in-progress/` nicht mehr leer war.
- **Grenzen:** Der Skill ist die größere Menge gegenüber `harness/review.md`: die älteren Skill-HIGHs (Dual-Classifier,
  Referenzmodell) und der Anker „Kommentar bricht mitten im Satz ab" stehen nur im Skill. Das Report-Gerüst der Baseline
  kennt die Spalte `Risiko` nicht, die Skill und `review.md` verlangen; die vendorte Vorlage ändert dieser Slice nicht.
- **Steering-Loop-Eintrag:** Reviewer-Skill und `harness/review.md` geschärft: `pfad` als Kurzzitat, `klasse` im Schema,
  Failure-Szenario-Pflicht für HIGH/MEDIUM, LOW nur mit Konventions-Anker — liegt in `.harness/skills/reviewer.md`.
  Gezählt, nicht als Regel mit Anker verkörpert, bleibt die Restlücke Skill ⊃ `review.md`.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt; die Restlücke trat zum ersten Mal auf.
- **Folge-Slices:** keine.
- **Risiken aus §6:** das eine Risiko trägt seinen Ausgang (entfallen, siehe §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Sub-Area
`spec/`, `harness/`, `docs/plan/` samt `.harness/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** u-boot führt bisher kein
Beobachtungs-Register; keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
