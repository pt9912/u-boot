# Slice implementieren (Harness)

Argument: $ARGUMENTS

Dieser Command führt die **Implementer**-Rolle (Baseline-Regelwerk `modul-09-implementierung.md`) für
*einen* Slice. Rollenfolge, Übergabe-Artefakte und Kontextgrenzen stehen in `harness/roles.md`;
Review und Verifikation laufen in **frischem Kontext**, nie im Kontext, der den Code schrieb.

Kanonisch ist `AGENTS.md` (Hard Rules, §6 Minimal Agent Workflow). Dieser Command wiederholt sie nicht,
er ergänzt, was die Werkzeuge dieses Repos erzwingen.

## Was die Werkzeuge erzwingen

- **Docker-only.** Der PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`) blockt
  Host-Toolchains (`go`, `golangci-lint`, Paketmanager) und prüft Sub-Shell-Strings. Nur
  `make`-Targets (`AGENTS.md` §3.1).
- **Gate-Nachweis und Stop-Hook.** `make gates` endet mit `record-gates`, das Content-Hash und
  HEAD-SHA stempelt. Der Stop-Hook verweigert das Turn-Ende, wenn HEAD ein anderer ist als beim
  letzten grünen Lauf **und** der Inhalt nicht gedeckt ist; streng (jedes Turn-Ende mit
  ungedecktem Inhalt) mit der Datei `.harness/stop-gate-streng` oder `STOP_GATE_STRENG=1`.
- **Commit-Kennung.** Der git-eigene Hook `.githooks/commit-msg` weist eine Message ohne Kennung ab
  (`ADR-NNNN`, `LH-XX-NN`, `MR-NNN`, `slice-<Kennung>`; Menge in `tools/harness/commit-msg-traceability.sh`).
  Aktiviert wird er pro Klon mit `make hooks-install`; `--no-verify` umgeht ihn. Geprüft wird die
  Anwesenheit einer Kennung, nicht ihre Wahrheit.
- **Commit via Message-Datei** (`git commit -F <datei>`): der Guard scannt den Command-String und
  ist quote-blind.
- **Lifecycle-Wechsel** per `make slice-mv SLICE=slice-<Kennung> TO=<open|next|in-progress|done>`: <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
  reiner `git mv`-Commit, Verweis-Nachzug als zweiter Commit (`AGENTS.md` §3.3). Bei `→ done`
  steht die Closure-Notiz **vor** dem Move. `make docs-check` danach zeigt, was stehen blieb.
- **Neue Artefakte** per `cp` aus `.harness/baseline/<tag>/templates/`, dann ausfüllen.

## Ablauf

1. `harness/README.md`, `AGENTS.md`, `harness/conventions.md` und die übergebene Slice-Datei lesen,
   dazu die referenzierten `LH-*`- und `ADR-*`-Quellen. Das Regelwerk nur abschnittsweise.
2. Berichten: Slice · `LH`-IDs · `ADR`-IDs · betroffene Komponenten · zu laufende Gates.
3. Der Slice liegt in `in-progress/`; sonst zuerst per `make slice-mv` dorthin. WIP-Limit 1. <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
4. **Plan vor Code.** Ist-Zustand gegen den Slice-Plan messen (`grep`/`diff`), Drift abgleichen,
   kleinste sinnvolle Änderung gegen die DoD planen.
5. Implementieren; engsten nützlichen Sensor, dann `make gates`. Ein roter Sensor führt zurück zum
   Plan, nicht zu Schritt 1.
6. Doku, ADR-Index, README, CHANGELOG nachziehen, wenn ein öffentlicher Vertrag berührt ist.
7. Pre-completion: DoD Punkt für Punkt **behaupten**, Sensor-Belege anhängen (`make gates` und die
   Nicht-Gate-Sensoren, die den Slice betreffen, z. B. `make test-docker`). Ein nicht gelaufener
   Sensor ist ein Befund und wird begründet. Zu jedem neuen oder geänderten Wächter die rot
   färbende Mutation benennen (`AGENTS.md` §3.6); neue oder geänderte Kommentare gegen §3.7 prüfen.
8. Übergabe: Diff + Plan-Verweis an den Reviewer (`.harness/skills/reviewer.md`, Report unter
   `docs/reviews/`), danach Verifier (`harness/verification.md`); Validator, falls der Slice
   Nutzerwert liefert. Kein Self-Review.

Gates nicht überspringen. Keine Erfolgsmeldung ohne Command-Ausgabe.
