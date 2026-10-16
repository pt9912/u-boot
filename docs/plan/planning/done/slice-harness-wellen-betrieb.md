# Slice Harness: Wellen-Betrieb nach Baseline (Welle-Dateien, Beobachtungs-Register)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Prozess-Umbau am Harness).

**Bezug:** [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)
(Roadmap-Form, Auflösungs-Trigger „solange Release-Versionen die Wellen sind“),
[`MR-002`](../../../../harness/conventions/MR-002-carveout-inventar.md),
[`MR-012`](../../../../harness/conventions.md#mr-012) (Punkt 4: planner, plan-welle,
close-welle, Register liegen ungenutzt). Voraussetzung:
[`slice-cr-lastenheft-produktanforderungen`](../done/slice-cr-lastenheft-produktanforderungen.md)
— die Planungs-Anforderungen verlassen das Lastenheft; ohne sie braucht dieser Slice
keinen Vertragswortlaut (heute schließt
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
Welle-Dateien als Dateinamen aus).

**Berührte Spec-Stellen:** — (nach dem CR keine)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** u-boot führt Wellen als flache Welle-Dateien (`welle-<kebab-slug>.md`
in `docs/plan/planning/`, bei Closure per `git mv` nach `done/`) und das
Beobachtungs-Register `docs/plan/planning/observations/`; die Roadmap trägt Zeiger
statt Wellen-Inhalt, und `plan-welle`, `close-welle` und `planner` laufen auf u-boot.

**Entscheidung (2026-10-16):** Welle-Dateien und Register ja, Archivierung nein.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Archivierung (`make archive-welle`, `done/<welle-id>/archiv.zip`, gekürzte Stubs):
  `done/`-Inhalte würden durch Stubs ersetzt und der Volltext verschwände aus dem
  Arbeitsbaum; eine eigene Entscheidung, falls je gewollt. Das Ziel `archive-welle`
  und `archivierung.mk` bleiben ungenutzt.
- Rückwirkende Welle-Dateien für abgeschlossene Releases (`v0.1.0`–`v0.9.0`): ihre
  Closure lebt in den `done/`-Release-Cut-Slices und in der Roadmap; nur neue Wellen
  bekommen eine Datei.
- Änderung am Lastenheft: führt
  [`slice-cr-lastenheft-produktanforderungen`](../done/slice-cr-lastenheft-produktanforderungen.md).

## 2. Definition of Done

- [x] **Roadmap und Doku auf die Welle-Datei-Form:** `MR-003` neu gefasst,
  `roadmap.md` mit Abschnitt „Offene Wellen“ samt Ruhe-Marker, `planning/README.md`,
  d-check-Modul `planning` (`heading`, `marker`) angepasst, `make doc-planning` grün.
- [x] **Werkzeuge auf u-boot gefasst:** `plan-welle`, `close-welle`, `planner`
  (ANPASSEN-Blöcke gefüllt, Hinweis „ungenutzt“ entfernt, Archivierungs-Schritt von
  `close-welle` gestrichen); das Beobachtungs-Register `observations/README.md` bleibt
  generisch, die Kürzel stehen in `harness/conventions.md` (Spalte Kürzel).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      ([`Report`](../../../reviews/2026-10-16-slice-harness-wellen-betrieb.md); 0 HIGH, 0 MEDIUM, 4 LOW, 3 INFO; F-1 bis F-4 behoben; `.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/conventions/MR-003-roadmap-wellen.md` | update | Release-Versionen bleiben die Wellen, aber als Welle-Dateien |
| `docs/plan/planning/in-progress/roadmap.md` | update | „Offene Wellen“ als Zeiger, Ruhe-Marker |
| `docs/plan/planning/README.md` | update | Slices vs. Wellen |
| `.d-check.yml` | update | Modul `planning`: Überschrift und Marker |
| `.claude/commands/plan-welle.md`, `close-welle.md`, `.claude/agents/planner.md` | update | ANPASSEN-Blöcke, Hinweis entfernen |
| `docs/plan/planning/observations/README.md` | unverändert | generisch, verweist auf die Kürzel-Spalte in `harness/conventions.md` |
| `harness/conventions.md` | update | ID-Schema `BEO-<KUERZEL>/<slug>`, Kürzel-Deklaration |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-cr-lastenheft-produktanforderungen` ist in
`done/`; die Planungs-Anforderungen stehen nicht mehr im Lastenheft.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Roadmap-Umbau und Werkzeug-Anpassung trennen.
- `in-progress` → `open` (blockiert): der CR ist nicht freigegeben (Stand 2026-10-16).

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Die Roadmap trägt die Closure heute schon; eine zweite Fassung als Welle-Datei
  dupliziert sie — **Ausgang:** entfallen, Welle-Dateien gelten nur für neue Wellen, die
  Detailquelle der Closure bleibt der Release-Cut-Slice
  ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md) (b)).
- `close-welle` ruft die Archivierung (Schritt 4) und verlangt sie als Vorbedingung —
  **Ausgang:** entfallen, der Schritt und die Results-Datei sind gestrichen, `close-welle`
  hat fünf Schritte.
- Die Zeiger-Hälfte von „Offene Wellen“ (ein Zeiger je Welle-Datei) hat keinen Sensor —
  **Ausgang:** weiter offen, benannt in `.d-check.yml`; kein Register, daher hier festgehalten.

## 7. Closure-Notiz

Vollzogen am 2026-10-16. u-boot führt Wellen als Welle-Dateien: Die Roadmap trägt
„Offene Wellen“ mit dem Ruhe-Marker (Modul `planning`, Marker-Hälfte geprüft), `MR-003` ist
neu gefasst, `plan-welle`, `close-welle` und `planner` sind auf u-boot gefasst (ohne
Archivierung und Results-Datei), und die Spalte Kürzel in `harness/conventions.md` bedient das
Beobachtungs-Register. Review: 0 HIGH, 0 MEDIUM, 4 LOW, 3 INFO, F-1 bis F-4 und F-7 behoben.
Verifier: DoD bestätigt, alle Sensor-Läufe grün.

- **Was hat funktioniert:** Der Lastenheft-CR zuerst (Planungs-Anforderungen als Grabsteine) hat
  den Vertragswortlaut für Welle-Dateien überflüssig gemacht. Der Review hat die Marker-Hälfte
  des Sensors in einem Wegwerf-Klon in beide Richtungen belegt.
- **Was ging anders als geplant:** Der Slice war zunächst als CR-Vehikel angelegt und wurde nach
  der Grundsatzentscheidung (Lastenheft nur mit Produktanforderungen) auf die reine Prozess-
  Umsetzung zurückgeschnitten. Mechanisch kopierte Vorlagentexte (Nummerierung in `close-welle`,
  Schreibweise `KUERZEL`) fand erst der Review.
- **Steering-Loop-Eintrag:** Nach dem Entfernen von Schritten aus einer übernommenen Anweisung
  werden Nummerierung und Querverweise der verbleibenden Schritte geprüft — liegt in
  `.claude/commands/close-welle.md` (Schritte 1–5).
- **Folge-Slices:** [`slice-harness-dcheck-bump-reviewer-skill`](../done/slice-harness-dcheck-bump-reviewer-skill.md),
  [`slice-harness-internal-readmes-festlegungen-in-spec`](../done/slice-harness-internal-readmes-festlegungen-in-spec.md).
- **Risiken aus §6:** je ein Ausgang, siehe §6.
- **Erste Welle-Datei:** entsteht mit der nächsten Eröffnung (`/plan-welle`); bis dahin trägt
  die Roadmap „Keine offene Welle-Datei“.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `harness/` und `docs/plan/`
samt `.claude/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.
