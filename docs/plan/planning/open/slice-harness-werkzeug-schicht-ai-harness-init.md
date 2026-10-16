# Slice Harness: Werkzeug-Schicht aus ai-harness-init prüfen und übernehmen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Phase:** Harness-/Baseline-Wartung (kein Produkt-Meilenstein, keine Welle).

**Bezug:** [`slice-harness-regelwerk-adoption-v6.18.1`](../done/slice-harness-regelwerk-adoption-v6.18.1.md);
[`MR-004`](../../../../harness/conventions/MR-004-regelwerk-vendored.md);
Gates-Haltung [`MR-005`](../../../../harness/conventions/MR-005-gate-haltung-docs-check.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Entscheiden, welche Teile der Werkzeug-Schicht eines mit
`ai-harness-init` (Baseline `v6.18.1`) erzeugten Repos u-boot übernimmt, und
die gewählten Teile einhängen. Der Referenzbestand umfasst: `harness/mk/*.mk`
mit `harness/mk/ai-harness-init.md`, `tools/harness/*.sh` (u. a.
`record-gates.sh`, `history-range-guard.sh`, `commit-msg-traceability.sh`,
`slice-mv.sh`), `.githooks/commit-msg`, `.claude/` (Agenten, Kommandos, Hooks),
`.harness/skills/closure-note-reviewer.md`.

**Ausdrücklich NICHT in diesem Slice:**

- Produktcode und Produkt-Gates — die Docker-only-Regel (`AGENTS.md` §3.1) und
  `make gates` bleiben maßgeblich.
- Stilles Lockern bestehender Gates (`AGENTS.md` §3.6).

## 2. Offene Fragen

- Überschneidung mit vorhandenen Werkzeugen (`tools/harness/`, `d-check.mk`,
  `.github/workflows/`) je Fragment klären; Übernahme nur, wo sie ersetzt oder
  ergänzt, nicht doppelt.
- d-check-Bump (`v0.51.1` → ≥ v0.86.0) für `reviews`-Schlüssel und die
  Vereinigung mehrerer Gate-Index-Teile.
- Auffrischung `.harness/skills/reviewer.md` (Finding-Anker, Kategorien).
- Spec-Stelle für Werkzeug-Festlegungen (neuer Spezifikations-Abschnitt der
  Baseline) und Gate-ADR(s) für übernommene Gates.

## 3. Definition of Done

- [ ] Entscheidung je Fragment dokumentiert (übernommen / abgelehnt, Grund).
- [ ] Übernommene Teile in `Makefile`/CI eingebunden, `harness/README.md`
  §Sensors nachgezogen.
- [ ] `make gates` und `make doc-check` grün.
- [ ] Review durchgeführt.

## 4. Risiken

- Hooks und Skripte verändern den Agenten-Lauf; Übernahme ohne Gegenprobe kann
  bestehende Workflows brechen — **Ausgang:** weiter offen.

## 5. Closure-Notiz

Offen.
