# Slice Harness: Wellen-Betrieb nach Baseline (Welle-Dateien, Archivierung, Beobachtungs-Register)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden ist (Prozess-Umbau am Harness).

**Bezug:** [`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)
(Roadmap-Form, Auflösungs-Trigger „solange Release-Versionen die Wellen sind“),
[`MR-002`](../../../../harness/conventions/MR-002-carveout-inventar.md);
Voraussetzung:
[`slice-harness-ai-harness-init-adoption`](../next/slice-harness-ai-harness-init-adoption.md)
(das Werkzeug legt `archivierung.mk`, `plan-welle`, `close-welle`, `planner` und
`observations/` an). Keine Anforderungsänderung — Prozessform, aber
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
ist zu prüfen.

**Berührte Spec-Stellen:** — (offen: Lastenheft-Aussagen zum Planning-Lifecycle
gegen die Welle-Datei-Form prüfen; bei Abweichung ein CR vor der Umsetzung)

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Entscheiden und, falls ja, umsetzen, ob u-boot den Wellen-Betrieb der
Baseline führt: Welle als flache Datei mit Closure per `/close-welle`,
Zeitdokumente geschlossener Wellen per `make archive-welle` archiviert,
Beobachtungs-Register `docs/plan/planning/observations/` als Arbeitsform.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Installation des Werkzeugs: liefert
  [`slice-harness-ai-harness-init-adoption`](../next/slice-harness-ai-harness-init-adoption.md).
- Rückwirkende Neufassung abgeschlossener Release-Cuts als Welle-Dateien: die
  Closure der bisherigen Wellen lebt in den `done/`-Release-Cut-Slices; ein
  Altbestand wird gesammelt archiviert (`archive-welle` mit `WELLE=altbestand`),
  nicht umgeschrieben.

## 2. Definition of Done

- [ ] Entscheidung dokumentiert (Adaption oder ADR): Welle-Dateien und Archivierung
  ja/nein; bei „ja“ `MR-003` aufgelöst oder neu gefasst und die Roadmap
  umgestellt, bei „nein“ die vom Werkzeug angelegten Commands und der Agent
  `planner` entfernt oder als nicht genutzt gekennzeichnet.
- [ ] Bei „ja“: erste Archivierung (`altbestand`) ausgeführt, Verweise
  aufgelöst, `make docs-check` grün.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/conventions/MR-003-roadmap-wellen.md` | update / auflösen | Roadmap-Form und Ort |
| `docs/plan/planning/in-progress/roadmap.md` | update | Offene Wellen als Zeiger auf flache Welle-Dateien |
| `docs/plan/planning/done/`, `docs/reviews/` | `git mv` / Archiv | Archivierung nach Wellen |
| `.claude/commands/`, `.claude/agents/planner.md` | update | auf u-boot gefasst |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-ai-harness-init-adoption` ist in `done/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Archivierung des Altbestands trennen von der
  Umstellung der Roadmap.
- `in-progress` → `open` (blockiert): das Lastenheft widerspricht der Welle-Datei-Form —
  CR zuerst.

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Die Archivierung bewegt Slice- und Review-Dateien und bricht Verweise darauf;
  das Werkzeug zieht sie nach, Immutabilität der Accepted-ADRs bleibt unberührt —
  **Ausgang:** offen bis zur Closure.
- Roadmap und Release-Cut-Slices tragen die Closure heute schon; eine zweite
  Fassung als Welle-Datei dupliziert sie — **Ausgang:** offen bis zur Entscheidung.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `docs/plan/` samt
`.claude/` ([`conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.
