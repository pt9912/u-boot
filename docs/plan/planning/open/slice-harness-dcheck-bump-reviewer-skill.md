# Slice Harness: Auffrischung des Reviewer-Skills (d-check-Bump übernommen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Wartungsarbeit am Harness).

**Bezug:** [`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md)
(Reviewer-Skill); Herkunft: Delta-Lektüre in
[`slice-harness-regelwerk-adoption-v6.18.1`](../done/slice-harness-regelwerk-adoption-v6.18.1.md)
§3/§4 und Closure-Notiz von
[`slice-harness-werkzeug-schicht-ai-harness-init`](../done/slice-harness-werkzeug-schicht-ai-harness-init.md).
Keine Anforderungsänderung — Skill-Pflege.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `.harness/skills/reviewer.md` an die Reviewer-Vorlage der Baseline
`v6.18.1` angleichen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Bump von d-check auf `v0.86.1`: den führt
  [`slice-harness-ai-harness-init-adoption`](../in-progress/slice-harness-ai-harness-init-adoption.md)
  aus, weil das Werkzeug `d-check.mk` samt digest-gepinntem Image mitbringt. Der
  Dateiname dieses Slice behält die frühere Fassung seines Gegenstands.
- `closure-note-reviewer.md`: nicht Gegenstand des Delta (nur ein
  Pfad-Hinweis), wird nur mitgezogen, falls der Skill-Abgleich dort denselben
  Pfad-Anker verlangt.

## 2. Definition of Done

- [ ] `.harness/skills/reviewer.md` trägt die Vorlagen-Änderungen:
  `pfad` als Datei plus wörtliches Kurzzitat (Zeile nur als Lesehilfe), LOW nur
  mit Konventions-Anker, „Kein Stil-Polizist“, kein HIGH/MEDIUM ohne
  Failure-Szenario; `harness/review.md` ist im selben Zug nachgezogen, wo es
  dasselbe Schema führt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Doku-Update für `MR-009`, falls deren Aussage den neuen Stand berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/skills/reviewer.md` | update | Angleichung an `reviewer.template.md` der Baseline `v6.18.1` |
| `harness/review.md` | update | Output-Schema `pfad` und Kategorien-Regel, falls dort geführt |

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
  widerspricht — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Sub-Area
`spec/`, `harness/`, `docs/plan/` samt `.harness/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** u-boot führt bisher kein
Beobachtungs-Register; keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
