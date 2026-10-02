# MR-009 — Skill-Dateien unter `.harness/skills/`, Review-Reports unter `docs/reviews/`

- **Datum:** 2026-07-25
- **Geltungsbereich:** `.harness/skills/reviewer.md`,
  `.harness/skills/closure-note-reviewer.md`,
  [`docs/reviews/`](../../docs/reviews/README.md) (Ablage + README);
  Quellen-Rolle von [`review.md`](../review.md) und
  [`verification.md`](../verification.md).
- **Ersetzt-Baseline-Regel:** [`modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill](../../.harness/baseline/v6.13.0/regelwerk/modul-10-review-harness.md#ziel-form-reviewer-skill)
- **Adaption:** Review-Wissen liegt in u-boot an **zwei** Orten mit klarer
  Rollentrennung: die kanonische **Prosa** bleibt im Autoren-Bestand `harness/`
  (`review.md`, `verification.md`), die vom Regelwerk geforderte **Skill-Form**
  (Kontext-Eingang, repo-spezifische Klassifikations-Anker, Nicht-Zustaendig-
  keiten, Output-Schema, Steering-Loop) liegt maschinen-materialisiert unter
  `.harness/skills/` (`MR-007`-Ortswahl). Die Skills **duplizieren keine
  Tabellen**, sie verweisen aufwaerts auf die Prosa; bei Abweichung gewinnt die
  Prosa. Die Report-Ablage ist `docs/reviews/` (ein Report pro Lauf, Folgelaeufe
  als neue Datei); die Report-**Vorlage** wird bewusst **nicht** ins Repo
  kopiert, sondern aus
  `.harness/baseline/<tag>/templates/docs/reviews/` genommen - eine zweite Kopie
  waere eine Drift-Quelle. Zwei Abweichungen vom Vorlagen-Wortlaut: (a) u-boot
  hat **kein** `agent-review`-Make-Target - Reviews sind agentisch, Beleg ist
  der Report; (b) der Closure-Note-Skill hat **kein** vorgelagertes
  computational Struktur-Gate (ein `check_closure_notes.py`-Aequivalent
  existiert hier nicht) und traegt deshalb Struktur- **und** Inhaltspruefung,
  mit den Pflichtfeldern aus [`verification.md`](../verification.md) als
  Struktur-Quelle.
- **Begruendung:** Das Regelwerk fordert eine per-Repo Skill-Datei, nicht nur
  Prosa - ohne sie driftet der Reviewer zwischen Sessions. Die Doku-Mindest-
  struktur
  ([`LH-FA-PROJDOCS-001`](../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur))
  ist als **Mindest-**, nicht als Maximalstruktur formuliert; `docs/reviews/`
  ergaenzt sie additiv und erfuellt die dortige README-Pflicht je
  Unterverzeichnis. Kein Change Request am Vertrags-Stratum noetig.
- **Aufloesungs-Trigger:** permanent. Re-evaluieren, falls u-boot ein
  computational Closure-Note-Gate einfuehrt (dann wird der Closure-Note-Skill
  auf die semantische Schicht zurueckgeschnitten) oder falls ein
  `agent-review`-Target entsteht.
