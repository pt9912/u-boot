# MR-011 — Werkzeug-Fragmente unter `harness/mk/` repo-geführt; Gate-Nachweis in `gates`

- **Datum:** 2026-10-16
- **Geltungsbereich:** [`Makefile`](../../Makefile) (`GATE_CHECKS`, `include harness/mk/*.mk`),
  `harness/mk/`, `tools/harness/` (außer `fetch-baseline-cache.sh`), `.githooks/`,
  `.claude/` (Agenten, Command, Hooks, `settings.json`), [`.gitignore`](../../.gitignore),
  [`.d-check.yml`](../../.d-check.yml) (`targets.makefiles`), [`harness/README.md`](../README.md) §Sensors.
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt](../../.harness/baseline/v6.18.1/regelwerk/grundlagen-harness-dateien.md#harnessreadmemd-als-einstiegspunkt)
- **Adaption:** Die Fragmente `slice-mv.mk`, `hooks-install.mk` und `enforce.mk` und die <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
  Skripte darunter stammen aus der Werkzeug-Schicht von `ai-harness-init` (Baseline `v6.18.1`),
  werden hier aber **vom Repo geführt**: kein Werkzeug schreibt sie bei einem Lauf neu, und es
  gibt keinen werkzeug-eigenen Index-Teil `harness/mk/<werkzeug>.md`. Ihre Targets stehen in
  [`harness/README.md`](../README.md) §Sensors, dem einen Gate-Index; `targets.makefiles`
  nennt die drei Dateien einzeln (die d-check-Pinversion kennt keinen Glob). `make gates` hängt `record-gates` als letzten Schritt an
  (`GATE_CHECKS` plus Ordnungskante); die Gate-Liste selbst ist unverändert. Nicht
  übernommen sind die Teile, die das Werkzeug selbst oder seinen Träger-Binary voraussetzen
  (`traeger-fetch`, `selbstpruefung`, `e2e-abdeckung`, `erfassung`/Span-Hooks, `archivierung`),
  `baseline-verify` (ersetzt durch `tools/harness/fetch-baseline-cache.sh --verify`,
  [`MR-004`](MR-004-regelwerk-vendored.md)), `doc-gate.mk` (u-boot bindet `d-check.mk` selbst
  ein, [`MR-005`](MR-005-gate-haltung-docs-check.md)) sowie die Commands `plan-welle` und
  `close-welle` und der Agent `planner` (sie setzen flache Welle-Dateien, das
  Beobachtungs-Register und die Archivierung voraus, [`MR-003`](MR-003-roadmap-wellen.md)).
  `implement-slice` ist auf u-boot gekürzt. Der PreToolUse-Guard sperrt zusätzlich zur
  Basisliste die Host-Toolchain der Sprache Go
  (`tools/harness/blocked/go`, [`AGENTS.md`](../../AGENTS.md) §3.1).
- **Begründung:** Die Teile erzwingen Regeln, die `AGENTS.md` schon trägt (Docker-only,
  Kennungen in Commits, reine `git mv`-Commits, Gate vor Handoff), mechanisch statt per
  Disziplin. Ein Index-Teil, den kein Werkzeug schreibt, wäre ein zweiter Index ohne Eigentümer.
- **Auflösungs-Trigger:** u-boot führt `ai-harness-init` ein und lässt es `harness/mk/` neu
  schreiben; dann entsteht `harness/mk/ai-harness-init.md` und diese Adaption wird aufgelöst.
