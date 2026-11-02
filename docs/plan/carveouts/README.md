# Carveouts — u-boot

Aktive Carveouts mit Auflösungs-Trigger. Aufgelöste Carveouts wandern nach
`done/` (reiner `git mv`); der Zustand ist die Verzeichnis-Position.

## Aktive Carveouts

| ID | Titel | Gate | Trigger | Folge-Slice |
|---|---|---|---|---|
| [CO-001](CO-001-crlf-logs.md) | CRLF-Edge-Case in `u-boot logs --json` (`splitLogLines` splittet nur auf `\n`) | `make test` (Deckung der CLI-Vertrags-Tests; keine Schwelle — der Carveout nimmt nichts aus, er dokumentiert die Grenze) | erster bestätigter Befund über ein trailing `\r` in der `--json`-Antwort (Konsumer-Meldung, Issue oder Review-Finding) | [`slice-v1-logs-crlf`](../planning/open/slice-v1-logs-crlf.md) |

## Aufgelöste Carveouts

(noch keine)

## Konventionen

- Jeder aktive Carveout braucht: Auflösungs-Trigger, Folge-Slice, letzten
  Prüf-Termin; die Dateikonvention ist `docs/plan/carveouts/CO-<NNN>-<kurztitel>.md`,
  Vorlage `.harness/baseline/v6.18.1/templates/docs/plan/carveouts/carveout.template.md`.
- Bei Welle-Closure: Carveout-Audit zwingend — welche gültig, welche aufgelöst?
- **Dauerhafte Begründungen leben nicht hier** (Trichter,
  Baseline-Regelwerk `modul-07-carveouts.md` §Werkzeug-Wahl mit Projektinhaber-Entscheidung
  2026-11-02): die dauerhaften Ausschlüsse des früheren Master-Inventars tragen
  [ADR-0003](../adr/0003-solid-nahes-lint-profil.md) (vier Lint-Ausnahmen, Zeilen 62 bis 72),
  [ADR-0042](../adr/0042-dauerhafte-lint-ausnahmen.md) (`contextcheck`, `interfacebloat` —
  im Migrationsslice angelegt), [SPEC-031](../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)
  (GNU `make` am Host) und [SPEC-034](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
  (depguard-Regeln production-only); die Werkzeug-Wahl-Spur der
  Überführung trägt [`MR-002`](../../../harness/conventions/done/MR-002-carveout-inventar.md) (resolved).
- Siehe Baseline-Regelwerk `modul-07-carveouts.md` (§Ziel-Form: Carveout,
  §Carveout-Audit-Slice). Diese Datei ist der Index; ihre Zeilen sind Pointer,
  keine zweite Quelle — Quelle der Wahrheit sind die CO-Dateien und die Anker.