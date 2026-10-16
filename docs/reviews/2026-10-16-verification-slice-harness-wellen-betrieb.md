# Verification Evidence: slice-harness-wellen-betrieb — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Diff `67aa428..HEAD` (letzter Commit `99ac081`), 22 Dateien. Alle Sensors wurden in diesem Lauf selbst gefahren (Docker-only), nicht aus der Implementer-Behauptung übernommen.

## Scope

- Slice: [`slice-harness-wellen-betrieb`](../plan/planning/in-progress/slice-harness-wellen-betrieb.md)
- IDs: [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md), [`MR-012`](../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md) Punkt 4, [`MR-002`](../../harness/conventions/MR-002-carveout-inventar.md); keine `LH-*`/`ADR-*` berührt (Slice: „Berührte Spec-Stellen: keine“).
- Artefakte: `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/README.md`, `.d-check.yml`, `harness/conventions.md`, `harness/conventions/MR-003-roadmap-wellen.md`, `MR-012`, `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md`, `.claude/agents/*.md`.

## DoD-Abgleich

- [x] **Roadmap und Doku auf die Welle-Datei-Form** — bestätigt:
  - `roadmap.md` trägt `## Offene Wellen` (Zeiger-Absatz, benennt beide Slices in `in-progress/`), den Ruhe-Marker `Nichts in Arbeit.` nur in `.d-check.yml`, nicht in der Roadmap (richtig: es liegen zwei Slices in `in-progress/`).
  - `MR-003` neu gefasst: Adaption (a) Ort, (b) keine Results-Datei, (c) keine Archivierung, keine rückwirkenden Welle-Dateien; Aufnahme-Zeile im Index von `harness/conventions.md` unverändert gültig.
  - `docs/plan/planning/README.md` §Slices vs. Wellen: flache `welle-<kebab-slug>.md`, `git mv` nach `done/` bei Closure, Abweichungen (keine Results, keine Archivierung).
  - `.d-check.yml` Modul `planning`: `heading: "## Offene Wellen"`, `marker: "Nichts in Arbeit."`; Kommentar benennt die ungewächterte Zeiger-Hälfte.
  - `make doc-planning`: 0 Befunde (Sensor-Lauf unten).
- [x] **Werkzeuge auf u-boot gefasst** — bestätigt:
  - `grep -n 'ANPASSEN'` in `plan-welle.md`, `close-welle.md`, `planner.md`: kein Treffer.
  - `grep 'ungenutzt'` in `.claude/`: nur die beabsichtigten Negativ-Aussagen „`make archive-welle` bleibt ungenutzt“ in `plan-welle`/`close-welle`; der Hinweis „Anweisungssätze sind ungenutzt“ im `planner` ist entfernt (Diff `planner.md`).
  - `close-welle`: Archivierungs- und Results-Schritt gestrichen; „fünf Schritte“ = Schritt 1–5; Querverweise (Schritt 3, 4, 5) stimmen untereinander.
  - `planner.md`: Repo-Bezug-Absatz an die Stelle des ANPASSEN-Blocks.
  - `observations/README.md`: im Diff `67aa428..HEAD` unverändert (Plan-Tabelle: „unverändert“); die Kürzel stehen in der Spalte *Kürzel* in `harness/conventions.md`: neun Sub-Areas belegt (DOK, DOM, APP, PORT, CLI, ADP, E2E, WZG, PKG), `cmd/uboot` mit Strich (keine Sub-Area); Schreibweise `KUERZEL` in README, `close-welle`, `plan-welle` und `conventions.md` einheitlich.
- [x] **`make gates` grün** — Evidence: Sensor-Lauf unten (Exit 0, Coverage 90,90 % bei Schwelle 90 %).
- [x] **Review durchgeführt** — Report [`docs/reviews/2026-10-16-slice-harness-wellen-betrieb.md`](2026-10-16-slice-harness-wellen-betrieb.md) liegt vor; Zahlen stimmen (Tabelle: 0 HIGH, 0 MEDIUM, 4 LOW, 3 INFO; sieben Findings F-1..F-7). Behebung gegen den Ist-Stand geprüft:
  - F-1: Roadmap „Offene Wellen“ nennt beide Slices in `in-progress/` und kein „ein Slice ohne Welle“ mehr — behoben.
  - F-2: README-Zeile „ohne Wellen-Betrieb löst die jeweilige Slice-Closure …“ steht weiter (bedingt formuliert, generisch; Slice-Plan hält das Register-README bewusst unverändert); die Schreibweise ist nachgezogen — Schreibweise behoben, Zustandsaussage bewusst belassen, siehe Abweichung A-1.
  - F-3: Closure-Form-Notiz trennt „bisherige Wellen“ (Release-Cut-Slice) und „neue Wellen“ (§7 der Welle-Datei, Verweis auf den Slice) — behoben.
  - F-4: `close-welle` Listennummern 2–7 und „Schritt 1–5“ sind eine feste Versatz-Folge (Liste setzt bei Punkt 1 ein), kein Sprung mehr — behoben.
  - F-7: in der Roadmap-Tabelle „Abgeschlossene Wellen“ keine Leerzeile vor `welle-spec-technik-stratum` — behoben.
  - F-5, F-6 (INFO) bewusst offen: `archivierung.mk` ohne Sperre, Zeiger-Hälfte ohne Sensor — beides in `MR-012` Punkt 4 bzw. `.d-check.yml` benannt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen, bewusst (Slice §7 „Offen.“); vor `done/` zu schreiben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: beide Risiken stehen auf „offen“; vor `done/` aufzulösen.

## Zusatzprüfungen

- **Lastenheft-Konsistenz:** `LH-FA-PROJDOCS-003`..`006` sind seit 0.6.0 Grabsteine (`spec/lastenheft.md`, Zeilen 886–899); der Slice braucht keinen Vertragswortlaut. Der Slice-Kopf (Bezug) beschreibt den Vor-CR-Zustand („heute schließt … aus“) — Zeitdokument im Slice, keine Behauptung des Repos über den Ist-Stand. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
- **Residuen:** `grep -rIE 'Keine eigenständigen Welle|welle-NN-results|Aktuelle Welle|keine aktive Welle'` außerhalb `done/`, `docs/reviews/`, `.harness/`: kein Treffer. `grep -i 'keine welle-dateien|Wellen sonst inline'`: kein Treffer.
- **MR-012 Punkt 4 vs. Ist:** Text („`planner`, `plan-welle`, `close-welle` und Register … auf u-boot gefasst; `archivierung.mk`/`make archive-welle` bleiben ungenutzt“) deckt sich mit den Dateien; Auflösungs-Trigger (Archivierung einführen) deckt sich mit `MR-003` (c). Verweis auf den alten `open/`-Pfad des Slice entfernt.
- **Keine Lockerung (AGENTS.md §3.6):** Diff gegen `.golangci.yml`, `Makefile`, `harness/mk`, `spec/`, `.github` leer; `.d-check.yml` ändert nur Überschrift und Marker des Moduls `planning`, das Modul bleibt aktiv, `modules` unverändert.
- **Rollen-Schreibweise:** Erfassungs-Satz in den sechs `.claude/agents/*.md` ergänzt (Diff-Statistik je +1 Zeile) — wortgleich laut Review.

## Abweichungen

- **A-1 (INFO):** `docs/plan/planning/observations/README.md` §Wer liest führt weiter „ohne Wellen-Betrieb löst die jeweilige Slice-Closure diesen Schritt selbst aus“. Die Aussage ist generisch und bedingt formuliert und deckt sich mit dem DoD-Wortlaut („bleibt generisch“). Folge: kein Befund gegen die DoD; wer das Register-README trotzdem nachzieht, ändert die generische Vorlage.
- **A-2 (INFO):** Die Roadmap-Kopfzeile trägt „Letzte Änderung: 2026-10-14“, obwohl der Slice den Abschnitt „Offene Wellen“ am 2026-10-16 geändert hat. Reine Metadaten, kein DoD-Punkt.

Keine Abweichung gegen einen angehakten DoD-Punkt.

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0; `coverage-gate: OK — coverage 90.90% meets threshold 90%` |
| `make docs-check` | pass | `d-check: 238 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Berichts; nach Anlage erneut gefahren, siehe Schluss-Zeile) |
| `make doc-planning` | pass | 0 Befunde |
| `make doc-targets` | pass | 0 Befunde |
| `make doc-complete` | pass | `83 Anforderung(en), 0 Waise(n).` |
| `make baseline-verify` | pass | `baseline-verify: v6.18.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |

## Traceability

| ID | Beleg |
| --- | --- |
| [`MR-003`](../../harness/conventions/MR-003-roadmap-wellen.md) | `make doc-planning` (Marker-Hälfte), Diff `roadmap.md`/`planning/README.md`, Review-Wegwerf-Klon in beide Richtungen |
| [`MR-012`](../../harness/conventions/MR-012-ai-harness-init-werkzeug-des-repos.md) | Diff Punkt 4 gegen `plan-welle`/`close-welle`/`planner` und Archivierungs-grep |
| [`MR-002`](../../harness/conventions/MR-002-carveout-inventar.md) | unverändert, nicht berührt (Verweis im Slice-Kopf) |

## Carveouts

- Neu: none
- Gelöst: none
- Unverändert: alle bestehenden Einträge in `carveouts.md` (nicht berührt)

## Nicht ausgeführt

- `make test`, `make test-docker`, `make verify-depguard` — Diff berührt keinen Go-Code, keinen Docker-Pfad, keine Importregeln; `make gates` hat den Test- und Lint-Teil bereits abgedeckt.
- Wegwerf-Klon-Nachweis der Marker-Hälfte (beide Richtungen) — nicht wiederholt; Evidence des Reviews, `make doc-planning` bestätigt den Ist-Zustand.

## Commit / Artefakt

- Geprüfter Stand: `99ac081` (HEAD); dieser Bericht ist der einzige neue Inhalt des Verifier-Laufs.
