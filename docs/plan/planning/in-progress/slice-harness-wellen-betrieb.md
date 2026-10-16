# Slice Harness: Wellen-Betrieb nach Baseline (Welle-Dateien, Beobachtungs-Register) — zugleich Change Request

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
close-welle, Register liegen ungenutzt);
Change Request auf
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
(Dateinamen in `planning/`, Master-Dokumente). Der Slice ist das CR-Vehikel, wie
[`slice-cr-adr-format-madr`](../done/slice-cr-adr-format-madr.md).

**Berührte Spec-Stellen:** [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) (Lastenheft, Vertrags-Stratum) —
Änderung nur nach Vereinbarung mit dem Projektinhaber (Lastenheft §16).

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** u-boot führt Wellen als flache Welle-Dateien (`welle-<kebab-slug>.md`
in `docs/plan/planning/`, bei Closure per `git mv` nach `done/`) und das
Beobachtungs-Register `docs/plan/planning/observations/`; die Roadmap trägt Zeiger
statt Wellen-Inhalt, und `plan-welle`, `close-welle` und `planner` laufen auf u-boot.

**Entscheidung (2026-10-16):** Welle-Dateien und Register ja, Archivierung nein.

**Warum ein CR:** [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
kennt in `planning/` nur `slice-…` und `tranche-…` als Dateinamen sowie
`roadmap.md` und `carveouts.md` als Master-Dokumente; Welle-Dateien und
`observations/` sind weder das eine noch das andere. Die Anforderung selbst ändert
sich, ein `harness/conventions.md`-Eintrag reicht nicht (Lastenheft §16).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Archivierung (`make archive-welle`, `done/<welle-id>/archiv.zip`, gekürzte Stubs):
  widerspricht [`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung)
  (Verschieben per `git mv` nach `docs/archive/`, Inhalt bleibt sichtbar) und der
  Regel „Inhalte in `done/` nur korrigierend änderbar“; ein eigener CR, falls je gewollt.
  Das Ziel `archive-welle` und `archivierung.mk` bleiben ungenutzt.
- Rückwirkende Welle-Dateien für abgeschlossene Releases (`v0.1.0`–`v0.9.0`): ihre
  Closure lebt in den `done/`-Release-Cut-Slices und in der Roadmap; nur neue Wellen
  bekommen eine Datei.
- Änderung der Lastenheft-Aussagen zu Slices, Tranchen und Carveouts.

## 2. Definition of Done

- [ ] **Lastenheft-Änderung** nach Freigabe des Projektinhabers:
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  nennt `welle-<kebab-slug>.md` für Welle-Pläne und `observations/` als Master-Ablage;
  Version `0.6.0`, Zeile in §16 Historie, Verweis auf die Vereinbarung.
- [ ] **Roadmap und Doku auf die Welle-Datei-Form:** `MR-003` neu gefasst,
  `roadmap.md` mit Abschnitt „Offene Wellen“ samt Ruhe-Marker, `planning/README.md`,
  d-check-Modul `planning` (`heading`, `marker`) angepasst, `make doc-planning` grün.
- [ ] **Werkzeuge auf u-boot gefasst:** `plan-welle`, `close-welle`, `planner`
  (ANPASSEN-Blöcke gefüllt, Hinweis „ungenutzt“ entfernt), Beobachtungs-Register
  `observations/README.md` an u-boots Kürzel angepasst.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

**Vorgeschlagener Lastenheft-Wortlaut** ([`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle), Absatz Dateinamen;
Freigabe steht aus):

> Dateinamen in `planning/`: `slice-<phase>-<kebab-slug>.md` für Slice-Pläne,
> `tranche-<nr>-<kebab-slug>.md` für Tranchen-Pläne und `welle-<kebab-slug>.md` für
> Welle-Pläne; eine offene Welle liegt flach in `planning/`, bei Closure wandert sie
> per `git mv` nach `done/`. Übergreifende Master-Dokumente (`roadmap.md`,
> `carveouts.md`) liegen dauerhaft in `in-progress/`, das Beobachtungs-Register
> `observations/` dauerhaft unter `planning/`; sie folgen keinem der drei Formate.

§16-Zeile (Version `0.6.0`, 2026-10-16): [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle): Welle-Pläne (`welle-<kebab-slug>.md`) und das Beobachtungs-Register als Planning-Artefakte zugelassen; Archivierung unberührt. Verweis: Vereinbarung mit dem Projektinhaber.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle), Version, §16 (nach Freigabe) |
| `harness/conventions/MR-003-roadmap-wellen.md` | update | Release-Versionen bleiben die Wellen, aber als Welle-Dateien |
| `docs/plan/planning/in-progress/roadmap.md` | update | „Offene Wellen“ als Zeiger, Ruhe-Marker |
| `docs/plan/planning/README.md` | update | Slices vs. Wellen |
| `.d-check.yml` | update | Modul `planning`: Überschrift und Marker |
| `.claude/commands/plan-welle.md`, `close-welle.md`, `.claude/agents/planner.md` | update | ANPASSEN-Blöcke, Hinweis entfernen; Archivierungs-Schritt 4 von `close-welle` entfällt |
| `docs/plan/planning/observations/README.md` | update | Kürzel nach `harness/conventions.md` |
| `harness/conventions.md` | update | ID-Schema `BEO-<KUERZEL>/<slug>`, Kürzel-Deklaration |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-ai-harness-init-adoption` ist in
`done/` (erfüllt, 2026-10-16).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Lastenheft-Änderung und Roadmap-Umbau in zwei Slices
  trennen.
- `in-progress` → `open` (blockiert): der Projektinhaber gibt die Lastenheft-Änderung
  nicht frei — dann Variante „Wellen bleiben Roadmap-Zeilen“ (Werkzeug-Dateien entfernen).

## 5. Closure-Trigger

DoD vollständig, Review- und Verification-Evidence liegen vor, Closure-Notiz
geschrieben, `make gates` grün.

## 6. Risiken und offene Punkte

- Die Lastenheft-Änderung braucht die Freigabe des Projektinhabers; ohne sie bleibt
  der Slice blockiert — **Ausgang:** offen bis zur Freigabe.
- Die Roadmap trägt die Closure heute schon; eine zweite Fassung als Welle-Datei
  dupliziert sie — **Ausgang:** offen; Welle-Dateien gelten nur für neue Wellen,
  die Detailquelle der Closure bleibt der Release-Cut-Slice.
- `close-welle` ruft die Archivierung (Schritt 4) und verlangt sie als Vorbedingung —
  **Ausgang:** offen; der Schritt wird für u-boot gestrichen und dort vermerkt.

## 7. Closure-Notiz

Offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
`docs/plan/` samt `.claude/`
([`conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register ist leer.

**Modus:** alle berührten Sub-Areas GF.
