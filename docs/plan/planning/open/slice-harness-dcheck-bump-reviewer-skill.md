# Slice Harness: d-check-Bump auf v0.86.1 und Auffrischung des Reviewer-Skills

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Wartungsarbeit am Harness).

**Bezug:** [`MR-005`](../../../../harness/conventions/MR-005-gate-haltung-docs-check.md)
(Gate-Haltung `docs-check`),
[`MR-009`](../../../../harness/conventions/MR-009-skills-und-review-ablage.md)
(Reviewer-Skill),
[`MR-011`](../../../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md);
Herkunft: Delta-Lektüre in
[`slice-harness-regelwerk-adoption-v6.18.1`](../done/slice-harness-regelwerk-adoption-v6.18.1.md)
§3/§4 und Closure-Notiz von
[`slice-harness-werkzeug-schicht-ai-harness-init`](../done/slice-harness-werkzeug-schicht-ai-harness-init.md).
Keine Anforderungsänderung — Werkzeug-Pin und Skill-Pflege.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das digest-gepinnte Gate-Image von d-check `v0.51.1` auf `v0.86.1`
heben (Mindeststand der Baseline `v6.18.1` ist `v0.86.0`) und den
Reviewer-Skill `.harness/skills/reviewer.md` an die Reviewer-Vorlage der
Baseline `v6.18.1` angleichen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Aktivierung neuer d-check-Module oder Schlüssel (`reviews` mit
  `require-promises`, Vereinigung mehrerer Gate-Index-Teile, `authority-disjoint`):
  ein Pin-Bump ist kein Gate-Umbau (dieselbe Abgrenzung wie im früheren
  Bump-Slice [`slice-harness-dcheck-image-bump`](../done/slice-harness-dcheck-image-bump.md)).
  Neue Module sind eigene Entscheidungen mit eigenem Slice, sobald der Bump
  ihre Wirkung auf den Bestand sichtbar gemacht hat.
- Rückbau der Einzelauflistung in `targets.makefiles` zugunsten eines Globs:
  hängt an der Aktivierung der Vereinigungs-Schlüssel, nicht am Pin.
- `closure-note-reviewer.md`: nicht Gegenstand des Delta (nur ein
  Pfad-Hinweis), wird nur mitgezogen, falls der Skill-Abgleich dort
  denselben Pfad-Anker verlangt.

## 2. Definition of Done

- [ ] `DCHECK_DIGEST` im `Makefile` zeigt auf `v0.86.1`, `d-check.mk` ist mit
  `d-check --print-mk` dieses Stands neu erzeugt; `make docs-check` ist grün,
  neue Befunde sind behoben oder als Ausnahme in `.d-check.yml` begründet.
- [ ] `.harness/skills/reviewer.md` trägt die Vorlagen-Änderungen:
  `pfad` als Datei plus wörtliches Kurzzitat (Zeile nur als Lesehilfe), LOW nur
  mit Konventions-Anker, „Kein Stil-Polizist“, kein HIGH/MEDIUM ohne
  Failure-Szenario; `harness/review.md` ist im selben Zug nachgezogen, wo es
  dasselbe Schema führt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Doku-Update für `MR-005`/`MR-009`, falls deren Aussagen den neuen Stand
  berühren (Pin-Nennung, Aktiv-Module).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` (`DCHECK_DIGEST`) | update | Digest des Images `v0.86.1` (Image per Tag auflösen, Digest eintragen) |
| `d-check.mk` | neu erzeugt | `d-check --print-mk` des neuen Stands; Pin-Kommentar `DCHECK_IMAGE` |
| `.d-check.yml` | update, falls nötig | Schlüssel, die der neue Stand ablehnt oder umbenennt; Befunde, die der Bump sichtbar macht |
| `.harness/skills/reviewer.md` | update | Angleichung an `reviewer.template.md` der Baseline `v6.18.1` (§1/§2) |
| `harness/review.md` | update | Output-Schema `pfad` und Kategorien-Regel, falls dort geführt |
| `harness/conventions/MR-005-gate-haltung-docs-check.md` | update, falls nötig | Pin-/Modulaussage |

Ansatz: erst den Bump allein (Digest, `d-check.mk`, ggf. `.d-check.yml`) mit
grünem `make docs-check` als eigener Commit, danach der Skill-Abgleich — so
bleibt ein Gate-Befund dem Bump zurechenbar und nicht der Skill-Änderung.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Verantwortlichen;
kein Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): der Bump macht so viele neue Befunde
  sichtbar, dass ihre Behebung den Skill-Abgleich verdrängt — dann den Bump
  als eigenen Slice abspalten.
- `in-progress` → `open` (blockiert): das Image `v0.86.1` ist für
  `linux/arm64` oder offline-hermetisch (`--network none`) nicht lauffähig —
  Carveout nach Modul 7 oder Wartezeit auf ein Folge-Release.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
mit Steering-Loop-Eintrag geschrieben, `make gates` auf dem Stand vor dem
`git mv` nach `done/` grün.

## 6. Risiken und offene Punkte

- Der Bump von 0.51.1 auf 0.86.1 überspringt 35 Minor-Versionen; Schlüssel können
  umbenannt oder verschärft sein und Bestandsdateien neu beanstanden —
  **Ausgang:** offen bis zur Closure.
- Der Reviewer-Skill ist durch MR-009 repo-gepflegt; die Vorlage ist kein
  Pflichtstand — übernommen wird nur, was `harness/review.md` nicht
  widerspricht — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Sub-Area
`spec/`, `harness/`, `docs/plan/` samt `.harness/` und den Gate-Dateien
(`Makefile`, `d-check.mk`, `.d-check.yml`); sie ist in
[`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration
geführt, keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** u-boot führt kein
Beobachtungs-Register; keine Treffer.

**Modus:** alle berührten Sub-Areas GF.
