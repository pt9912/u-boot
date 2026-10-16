# Slice Harness: ai-harness-init als Werkzeug des Repos (Aggregator, Fragmente, Sensor-Dateien)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Wartungsarbeit am Harness).

**Bezug:** [`MR-011`](../../../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md)
(Auflösungs-Trigger: u-boot führt `ai-harness-init` ein),
[`MR-004`](../../../../harness/conventions/MR-004-regelwerk-vendored.md),
[`MR-005`](../../../../harness/conventions/MR-005-gate-haltung-docs-check.md);
Herkunft:
[`slice-harness-werkzeug-schicht-ai-harness-init`](../done/slice-harness-werkzeug-schicht-ai-harness-init.md)
(Handübernahme einzelner Fragmente). Keine Anforderungsänderung.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** `ai-harness-init` (Baseline `v6.18.1`) läuft gegen u-boot und ist
fortan das Werkzeug, das `Makefile`, `d-check.mk`, `harness/mk/`,
`tools/harness/` (Werkzeug-Skripte) und `.claude/hooks/` schreibt. u-boots
eigene Targets leben in `repo.mk`, der Gate-Index wird zweiteilig
(`harness/README.md` plus `harness/mk/ai-harness-init.md`), und Gates mit mehr
als einem Satz Vertrag bekommen Sensor-Dateien unter `harness/sensors/`.

**Übernimmt:** den d-check-Bump aus
[`slice-harness-dcheck-bump-reviewer-skill`](../open/slice-harness-dcheck-bump-reviewer-skill.md)
(das Werkzeug bringt `d-check` `v0.86.1` digest-gepinnt mit); dort bleibt die
Auffrischung des Reviewer-Skills.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Wellen-Betrieb (flache Welle-Dateien, `archive-welle`, Beobachtungs-Register
  als Arbeitsform, Änderung von [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)):
  baut die Roadmap-Form und über hundert erledigte Slices um; eigener Slice
  [`slice-harness-wellen-betrieb`](../open/slice-harness-wellen-betrieb.md). Das Werkzeug
  legt `plan-welle`, `close-welle`, `planner` und `archivierung.mk` bei diesem
  Lauf mit an; sie bleiben bis dahin ungenutzt.
- Auffrischung von `.harness/skills/reviewer.md`: bleibt im Slice
  `slice-harness-dcheck-bump-reviewer-skill`, das Werkzeug fasst den Skill nicht an.
- Aktivierung der `reviews`-, `spans`- und `structure`-Module von `d-check`:
  eigene Entscheidungen mit Wirkung auf den Bestand, erst nach dem Bump zu messen.

## 2. Definition of Done

- [ ] u-boots Targets stehen in `repo.mk`; `make gates` fährt dieselben
  Checks wie vorher (`lint`, `test`, `coverage-gate`, `docs-check`) plus
  `baseline-verify` und endet mit `record-gates`; `ci` und `fullbuild` sind
  unverändert verfügbar.
- [ ] `ai-harness-init` ist gegen das Repo gelaufen; `harness/README.md`
  führt die Zeile „Targets der Werkzeug-Fragmente“, `.d-check.yml` misst
  gegen beide Index-Teile (`authority-disjoint`), `make docs-check` ist grün.
- [ ] Sensor-Dateien unter `harness/sensors/` für die Gates mit mehr als einem
  Satz Vertrag; die Target-Zellen verlinken sie.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Doku-Update: `MR-011` aufgelöst (nach `conventions/done/`), `MR-004`/`MR-005`
  auf den neuen Stand, neue Adaption für `repo.mk` und das behaltene
  `fetch-baseline-cache.sh` (Freshness-Audit).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `repo.mk` | neu (vor dem Lauf) | alle u-boot-Targets aus dem `Makefile`; das Werkzeug ersetzt das `Makefile` ohne Meldung und legt `repo.mk` nur an, wo sie fehlt |
| `Makefile`, `d-check.mk`, `harness/mk/*.mk`, `tools/harness/*` (Werkzeug-Skripte), `.claude/settings.json`, `.claude/hooks/*` | vom Werkzeug geschrieben | kanonische Teile, kein Handeingriff |
| `harness/mk/ai-harness-init.md` | vom Werkzeug geschrieben | werkzeug-eigener Index-Teil |
| `.d-check.yml` | update | `targets`: `makefiles`, `doc-tables`, `authority` als Listen, `authority-disjoint: true` |
| `harness/README.md` | update | Zeile der Werkzeug-Fragmente, Zeilen von `record-gates`/`slice-mv`/`hooks-install` entfallen (jetzt im Werkzeug-Teil, disjunkt) |
| `harness/sensors/*.md` | neu | aus `gate.template.md` für `coverage-gate`, `verify-depguard`, `test-docker` |
| `harness/conventions/MR-011-*`, `MR-004`, `MR-005`, `conventions.md` | update / `git mv` | Adaptionen nachziehen |

Ansatz: Probelauf im Wegwerf-Klon liegt vor (Umfang: `Makefile` ersetzt,
`d-check.mk` auf `v0.86.1`, Fragmente und Skripte kanonisiert, Werkzeug-Dateien
neu). Reihenfolge im Repo: erst `repo.mk` anlegen (eigener Commit), dann den
Lauf, dann `.d-check.yml`, Index und Sensor-Dateien.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung; kein anderer Slice in `in-progress/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Bump auf `d-check` `v0.86.1` macht so
  viele neue Befunde sichtbar, dass sie den Umbau verdrängen — den Bump als
  eigenen Slice abspalten.
- `in-progress` → `open` (blockiert): das Werkzeug läuft ohne Netz oder
  Docker nicht — Carveout nach Modul 7.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
mit Steering-Loop-Eintrag geschrieben, `make gates` auf dem Stand vor dem
`git mv` nach `done/` grün.

## 6. Risiken und offene Punkte

- Das Werkzeug ersetzt `Makefile` und `d-check.mk` ohne Meldung; ein vergessenes
  Target in `repo.mk` ginge verloren — **Ausgang:** offen; die Target-Liste vor
  und nach dem Lauf wird verglichen.
- `d-check` springt von `v0.51.1` auf `v0.86.1`; neue Befunde im Bestand —
  **Ausgang:** offen bis zur Closure.
- Das Werkzeug legt Wellen-Dateien und -Commands an, die `MR-003` widersprechen —
  **Ausgang:** weiter offen, übernimmt `slice-harness-wellen-betrieb`.
- Der Träger-Binary liegt nicht im Repo (`.harness/state/bin/`, ignoriert); ohne
  ihn erfasst `span-emit` nichts — **Ausgang:** offen; `make traeger-fetch` pro Klon.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Sub-Area
`spec/`, `harness/`, `docs/plan/` samt `.harness/`, den Gate-Dateien und
`.claude/` ([`conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Werkzeug legt
`docs/plan/planning/observations/` mit seiner `README.md` an; sie trägt keine
Einträge, keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
