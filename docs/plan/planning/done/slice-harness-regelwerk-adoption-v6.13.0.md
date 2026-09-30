# Slice Harness: Baseline-Bump Regelwerk v3.5.2 → v6.13.0

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Phase:** Harness-/Baseline-Wartung (kein Produkt-Meilenstein, keine Welle).

**Bezug:** Baseline-Pin in
[`harness/conventions.md`](../../../../harness/conventions.md) Abschnitt
Baseline; `MR-004` Bump-Prozedur; Freshness-Audit ebenda. Beruehrt
`MR-001`..`MR-009` als Gegenproben, aendert sie nicht. Kein `LH`-Neubezug.

**Autor:** pt9912. **Datum:** 2026-09-29.

---

## 1. Ziel

Baselines-Bump des adoptierten Betriebsregelwerks (AI-Harness-Kurs,
`pt9912/ai-harness-course`) von `v3.5.2` auf `v6.13.0` als Einheit nach
`MR-004` (Pin, Vendor-Pfad, `AGENTS.md`, `harness/README.md`) samt
Delta-Lektuere und `MR-*`-Gegenprobe.

## 2. Verification-Evidence: Freshness-Audit

- **Sensor:** `tools/harness/fetch-baseline-cache.sh --verify` — **gruen**
  (`42` Dateien, Stand `v3.5.2`, vor dem Bump) und **gruen** nach dem Bump
  (`54` Dateien, Stand `v6.13.0`).
- **Sensor:** `tools/harness/fetch-baseline-cache.sh --check-freshness` —
  **Exit 3** (Review-Bump faellig), Ausfuehrung 2026-09-29: Pin `v3.5.2`,
  46 neuere Tags ab `v3.6.0` bis **`v6.13.0`** (neuester, 2026-09-28).
  Negativbefund wie gefordert in dieser Evidence festgehalten; kein
  Auto-Update, Bump als Slice.

## 3. Delta-Lektuere (v3.5.2 → v6.13.0)

Umfang (empirisch, Release-Bundle gegen den Vendortbestand):
~2.100 Diff-Zeilen im `regelwerk/`, ~1.600 in `templates/`; 43 → 54 Dateien.

- **Grundlagen-Schicht ausdifferenziert:** 3 → 7 Dateien (neu u. a.
  `grundlagen-begriffe.md`, `grundlagen-source-precedence.md`,
  `grundlagen-traceability.md`, `grundlagen-referenz-richtung.md`,
  `grundlagen-harness-dateien.md`, `grundlagen-bootstrap.md`).
- **Modul-Renamings:** `modul-03-lastenheft.md` → `modul-03-spec.md`,
  `modul-04-architektur-adrs.md` → `modul-04-adrs.md`.
- **Neue Planning-Templates:** `observation.template.md`
  (Beobachtungs-Register, `BEO-*`), `reconciliation.template.md`
  (bewusst nur fuer Brownfield-Bootstrap), `archiv-stub-{slice,welle}.template.md`,
  `welle-results.template.md`.
- **Templates/Harness:** `conventions.template.md` bleibt als Einzeldatei
  gueltig ("die Form (Einzeldatei vs. Verzeichnis, ADR-artig vs. Prosa)
  ist Wahl"); neu dazu `harness/conventions/MR-NNN-titel.template.md`
  (optionale ADR-artige Form) und `harness/sensors/gate.template.md`.

## 4. MR-Gegenproben

| MR | Befund |
|---|---|
| `MR-000`/`MR-001` | **valide.** Zwei-Straten-Source-Precedence bleibt kompatibel; die Regel-Stelle liegt jetzt verteilt in `grundlagen-source-precedence.md` / `grundlagen-traceability.md` (Prosa-Zitate in `MR-001` nennen den historischen Pfad — Fundstellen-Erwaehnung, kein Link). |
| `MR-002` | **valide.** Carveout-Templates existieren weiter; das Inventar-Muster bleibt die dokumentierte Abweichung. |
| `MR-003` | **valide.** `welle-results.template.md` ist neu, bleibt aber optionale Form — die Ablehnung der `welle-NN-results.md`-Ebene ist unveraendert repo-lokal begruendet. |
| `MR-004` | **Vollzug** durch diesen Slice (Pin, Vendor-Pfad, `AGENTS.md`, `harness/README.md`). |
| `MR-005` | **valide.** d-check-Template-Delta ist dokumentarisch (Zwei-Schritte-Hinweis, `link-policy`-Erlaeuterung); u-boots `.d-check.yml` traegt die aktive Config. |
| `MR-008` | **valide, kein CR.** ADR-Template behaelt die MADR-Form (Inline-Kopf-Felder, `Schaerft`-Feld) — vertraegskompatibel mit [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format); Delta sind Klaerungen (per-Sub-Area-Zahlraum optional, Kennungs-Vorrang im `Schaerft`-Feld, Baseline-Pfade statt Kurs-URLs). |
| `MR-009` | **valide.** Reviewer-Template-Delta praezisiert Dissens-vs.-Nicht-Determinismus bei mehreren Rolleninhabern und Template-Kommentar-Disziplin — Hinweis fuer die naechste Auffrischung der Skill-Datei, keine Pflicht. |

## 5. Definition of Done

- [x] Vendor-Lauf `v6.13.0` (`54` Dateien) mit `--verify` gruen.
- [x] Alten Bestand `.harness/baseline/v3.5.2/` entfernt (nur gepinnter
  Stand bleibt im Vendor-Pfad).
- [x] Pin-Einheit nach `MR-004`: `**Stand:**`, Integritaets-Pfad,
  Adoptierte-Quellen, Sync-Trigger in
  [`harness/conventions.md`](../../../../harness/conventions.md);
  `AGENTS.md` (3 Stellen); `harness/README.md` (1 Stelle).
- [x] Live-Verweise nachgezogen: `.harness/skills/reviewer.md`,
  `docs/reviews/README.md`, `docs/plan/planning/in-progress/roadmap.md`
  (je Baseline-Pfad). Historische Nennungen in `done/`, `docs/reviews/`
  (Beleg-Kontext) und `MR-008`-Nachtrag bleiben absichtlich stehen.
- [x] `MR`-Gegenprobe dokumentiert (§4).
- [x] `make doc-check` gruen.

## 6. Nicht-Ziele

- Keine Produktcode-Aenderung; keine Konvertierung des Ledgers in die
  neue `harness/conventions/`-Verzeichnisform (Form ist Wahl).
- Kein Rueckbau der `check_refs.py`-Referenz in historischen ADRs/Slices;
  das Upstream-Template entfernt sie stattdessen selbst.

## 7. Closure-Notiz

Vollzogen am 2026-09-29. Der Bump war wegen des Umfangs (46 Releases, drei
Major-Spruenge, Modul-Renamings) als Adoptions-Slice gefuehrt, nicht als
toter Punkt-Bump. Die Renamings (`modul-03-spec`, `modul-04-adrs`) und die
neuen Referenz-Formen (`observation`, `welle-results`) erzeugen keinen
Nachlauf-Pflichtbestand; erste Anwendung faellt bei der naechsten
Artefakt-Neuanlage an. Steering-Loop: d-check-Regel "Kennungen brauchen
Markdown-Links (auch in Inline-Code)" hat den Spec-Change am selben Tag
gepraegt (vgl.
[`slice-cr-sandbox-devcontainer`](../in-progress/slice-cr-sandbox-devcontainer.md)
Korrekturschleifen) — die Link-Pflicht ist also nicht nur aktiviert,
sondern auch belastet.
