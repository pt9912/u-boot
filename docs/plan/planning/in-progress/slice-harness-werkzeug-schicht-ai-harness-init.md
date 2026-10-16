# Slice Harness: Werkzeug-Schicht aus ai-harness-init prüfen und übernehmen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Phase:** Harness-/Baseline-Wartung (kein Produkt-Meilenstein, keine Welle).

**Bezug:** [`slice-harness-regelwerk-adoption-v6.18.1`](../done/slice-harness-regelwerk-adoption-v6.18.1.md);
[`MR-004`](../../../../harness/conventions/MR-004-regelwerk-vendored.md);
Gates-Haltung [`MR-005`](../../../../harness/conventions/MR-005-gate-haltung-docs-check.md);
Adaption [`MR-011`](../../../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

**Ziel:** Entscheiden, welche Teile der Werkzeug-Schicht eines mit
`ai-harness-init` (Baseline `v6.18.1`) erzeugten Repos u-boot übernimmt, und
die gewählten Teile einhängen.

**Ausdrücklich NICHT in diesem Slice:**

- Produktcode und Produkt-Gates — die Docker-only-Regel (`AGENTS.md` §3.1) und
  `make gates` bleiben maßgeblich.
- Stilles Lockern bestehender Gates (`AGENTS.md` §3.6): die Gate-Liste von
  `make gates` ist unverändert, `record-gates` kommt als letzter Schritt dazu.
- Der d-check-Bump (`v0.51.1` → ≥ v0.86.0) und die Auffrischung von
  `.harness/skills/reviewer.md`: eigene Folge-Arbeit.

## 2. Entscheidung je Fragment

| Fragment des Referenzrepos | Entscheidung | Grund |
|---|---|---|
| `slice-mv.mk`, `tools/harness/slice-mv.sh` | übernommen | setzt `AGENTS.md` §3.3 mechanisch um; in einem Klon gegen den echten Bestand erprobt (reiner Move-Commit, kein Verweis zu ziehen) |
| `hooks-install.mk`, `.githooks/commit-msg`, `commit-msg-traceability.sh` | übernommen | setzt die Kennungs-Regel von `AGENTS.md` §5 durch; Aktivierung pro Klon per `make hooks-install`, nicht automatisch |
| `enforce.mk`, `record-gates.sh`, `working-tree-hash.sh`, Stop-Hook | übernommen | Gate-Nachweis vor Handoff (`AGENTS.md` §6 Schritt 6/8); `gates` hängt `record-gates` nach allen Checks an, die Gate-Liste bleibt |
| `.claude/hooks/pretooluse-command-guard.sh`, `extract-command.awk`, `blocked/go` | übernommen | Docker-only (`AGENTS.md` §3.1) mechanisch; Basisliste plus Go-Toolchain |
| `.claude/agents/` (`architect`, `implementer`, `reviewer`, `verifier`, `validator`) | übernommen, gekürzt | Rollen-Zuschnitt zu `harness/roles.md`; der Span-/Erfassungs-Absatz entfällt |
| `.claude/commands/implement-slice.md` | übernommen, auf u-boot neu gefasst | die Vorlage setzt das Beobachtungs-Register und die Welle-Archivierung voraus |
| `.claude/commands/plan-welle.md`, `close-welle.md`, Agent `planner` | abgelehnt | Wellen sind hier Roadmap-Zeilen ([`MR-003`](../../../../harness/conventions/MR-003-roadmap-wellen.md)); kein Register, keine Archivierung |
| `baseline-verify` | abgelehnt | `fetch-baseline-cache.sh --verify` leistet dasselbe |
| `doc-gate.mk`, `history-range-guard` | abgelehnt | `d-check.mk` ist eingebunden ([`MR-005`](../../../../harness/conventions/MR-005-gate-haltung-docs-check.md)) |
| `traeger-fetch`, `selbstpruefung`, `e2e-abdeckung`, `erfassung`, `archivierung`, Span-Hooks | abgelehnt | gehören zum Bootstrap-Werkzeug und seinem Träger-Binary |

## 3. Definition of Done

- [x] Entscheidung je Fragment dokumentiert (§2).
- [x] Übernommene Teile in `Makefile` eingebunden, `harness/README.md`
  §Sensors und `.d-check.yml` (`targets.makefiles`) nachgezogen, Adaption als
  [`MR-011`](../../../../harness/conventions/MR-011-werkzeug-fragmente-repo-gefuehrt.md).
- [x] `make gates` grün (lint, test, coverage-gate, docs-check; Stempel
  geschrieben, Stop-Hook gibt frei) und `make docs-check` grün.
- [x] Review durchgeführt: [`Report`](../../../reviews/2026-10-16-slice-harness-werkzeug-schicht-ai-harness-init.md)
  (0 HIGH, 1 MEDIUM, 3 LOW, 2 INFO); MEDIUM und die zwei Kopf-/Hash-LOW behoben.
- [ ] `make hooks-install` im eigenen Klon ausgeführt (bewusst nicht automatisch).

## 4. Risiken

- Hooks und Skripte verändern den Agenten-Lauf: der Guard ist quote-blind und
  blockt auch ein Heredoc mit einem Toolchain-Wort am Zeilenanfang —
  **Ausgang:** eingetreten, Abhilfe `git commit -F <datei>` bzw. Dateien per
  Editor-Tool schreiben (in `implement-slice` vermerkt).
- Der Guard erkennt nur das erste Wort je Segment; Toolchain-Aufrufe in
  `if`-/`for`-Rümpfen gehen durch (Review F-4) — **Ausgang:** weiter offen,
  der Guard ist ein Stolperdraht, keine Sandbox.
- `make slice-mv` zieht präfixlose Links auf `../in-progress/…` um (Review F-5)
  und `working-tree-hash.sh` braucht `sha256sum` auf dem Host (Review F-6) —
  **Ausgang:** entfallen, beides ist Verhalten der Werkzeuge, die Links lösen auf
  und ein fehlendes `sha256sum` bricht den Lauf laut.
- Der Commit-Hook weist Messages ohne Kennung ab (z. B. reine Querverweis-
  Commits) — **Ausgang:** weiter offen, Aktivierung ist pro Klon freiwillig.

## 5. Closure-Notiz

Offen.
