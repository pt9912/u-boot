# Review-Report: slice-v1-cli-dry-run-log-would — 2026-10-18

**Review-Art:** Code — Code-Review gegen Plan, Hard Rules und Konventionen (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** Commits `1f42dfd`, `a79bf59`, `64e853a`, `878e453` (Range `fa019d0..HEAD`, ohne den fremden Planner-Commit `044e0c9` und dessen Dateien `docs/plan/planning/open/README.md` und [`slice-v1-cli-help-handbuch-link`](../plan/planning/done/slice-v1-cli-help-handbuch-link.md)).

**Skill:** `.harness/skills/reviewer.md` @ Stand HEAD · <!-- d-check:ignore (Skill-Pfad als Zitat) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- [`slice-v1-cli-dry-run-log-would`](../plan/planning/done/slice-v1-cli-dry-run-log-would.md) (§1–§6)
- Vorgänger [`slice-v1-cli-dry-run-would-formulierung`](../plan/planning/done/slice-v1-cli-dry-run-would-formulierung.md) (Zusammenfassung auf stdout)
- [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)
- `AGENTS.md` (Hard Rules §3.2, §3.6, §3.7, §3.9, §3.10; §5)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Felder der Dry-Run-Zeile (`path`, `old`, `new`, `changed`, `items`) sind nicht gepinnt: Drei Mutationen (`l.Info(would)` ohne `args`, `"changed", 0` in `generateDevcontainer`, `"old", ""` im Egress-Pfad) ließen `make test` grün, obwohl Plan §3 „dieselben Felder“ zusagt. | Maintainability (Slice-Plan §3: „trägt dieselben Felder“) | `internal/hexagon/application/dryrun_log_test.go` · `infoMessages` (liest nur `e.Msg`) | ja — `make test` mit den drei Mutationen | Test pinnt Zusage „Felder unverändert“ nicht |
| F-2 | LOW | Der Doc-Kommentar von `devcontainerWriteMessage` nennt die Nicht-Dry-Run-Zeile „im bisherigen Wortlaut“; das beschreibt einen Vorzustand statt der geltenden Zusage. | `AGENTS.md` §3.7, §3.9 (a) | `internal/hexagon/application/dryrun_log.go` · „sonst `<Aktion>` im bisherigen Wortlaut“ | Review-only | Kommentar beschreibt Vorzustand |
| F-3 | INFO | Der Plan nennt `TestConfigSet_SilenceLoggerFalse_EmitsInfo (:72)`; das Symbol steht bei `:52`, `:72` ist die `t.Errorf`-Zeile darin. Symbol und Aussage stimmen. | `AGENTS.md` §3.10 | [`slice-v1-cli-dry-run-log-would`](../plan/planning/done/slice-v1-cli-dry-run-log-would.md) §3 · „`TestConfigSet_SilenceLoggerFalse_EmitsInfo` (`:72`)“ | Review-only | Zeilenhinweis im Plan ungenau |
| F-4 | INFO | Der Doc-Kommentar von `applyAllowExternalFeatureSources` erwähnt den neuen Parameter `mode` nicht (kein Verhaltensfehler; `logWrite` trägt ihn). | Maintainability | `internal/hexagon/application/generate.go` · `func (s *GenerateService) applyAllowExternalFeatureSources(` | Review-only | Doc-Kommentar kennt neuen Parameter nicht |
| F-5 | INFO | Der Zweig `would change` in `devcontainerWriteMessage` ist nicht erreichbar (`devcontainerAggregateAction` liefert nur `UpdatedBlock`, `Created`, `NoOp`) und ungetestet; `TestDryRunLog_GenerateDevcontainerUpdatedBlock` deckt `PreviewAndApply` nicht ab (andere Fälle in `TestDryRunLog_WriteLines` tun es). | Maintainability | `internal/hexagon/application/dryrun_log.go` · `return prefix + "would change"` | ja — `make test` (Zweig/Mutation) | Toter Rückfallzweig ohne Test |

## Prüfung nach Auftrag

**(1) Vollständigkeit.** Suche nach `.Info(`, `.Warn(`, `.Log(`, `.Error(` in `internal/` ohne Tests, jede Stelle gelesen. Application-Layer: alle Schreibmeldungen laufen über `logWrite` bzw. `devcontainerWriteMessage` (`config.go` 2x, `config_many.go`, `devcontainer_egress.go`, `generate.go` 6x). Übrig: `config.go` `maybeWarnOrphanFeatureActivation` (Zustandshinweis, im Dry-Run wahr), `generate_rollback.go` „rolled back“ (`rollbackOnError` mit `nil`-Journal gibt im Dry-Run `err` ohne Logzeile zurück: `s.rollbackFS` ist dort `nil`, `newGenerateJournal(nil)` = `nil`), `doctor.go` „checks complete“. `add`, `remove`, `init`, `up`, `down`, `logs`, `template` haben keine Info-/Warn-Zeile über Schreibergebnisse (`init --devcontainer` überspringt `generate` im Dry-Run: `initproject.go` `req.PreviewMode == driving.PreviewDryRun`). Adapter: `slogLogger` reicht nur durch; keine Schreibmeldung auf Logger-Pfad. Kein übersehener Rest.

**(2) Kriterium.** `mode == driving.PreviewDryRun` genau; `--dry-run --diff` ergibt `PreviewDryRun` (`previewModeFromFlags`), `--diff` allein `PreviewAndApply` und bleibt Vergangenheitsform (Test, Mutation m4). Legacy-Konstruktoren: `selectFS` ignoriert den Modus ohne Factory, `cmd/uboot/main.go` verdrahtet beide Services mit `…WithFactory`; das Argument des Plans stimmt. `--quiet` = Pegel `Error` (`slogLevel`), Lauf `generate devcontainer --dry-run --quiet` ohne Zeile. Pegel `Info`, Felder in allen sechs Aufrufen unverändert (Lesen des Diffs; Pinnung siehe F-1). No-op und Fehler: Test `TestDryRunLog_DryRunWritesNothingAndNoOpStaysQuiet`.

**(3) Mutationen** (`make test`, je danach `git checkout -- <Datei>`, Arbeitsbaum sauber):

| Mutation | Ergebnis |
|---|---|
| `config_many.go`: Dry-Run-Text auf `updated` | gefangen (`config_set_many`) |
| `dryrun_log.go`: `would update block` → `would change` | gefangen (`GenerateDevcontainerUpdatedBlock`) |
| `dryrun_log.go`: Kriterium `!= PreviewNone` (`--diff` allein trifft) | gefangen (Boundary) |
| `generate.go`: Allowlist-Zeile mit `PreviewNone` | gefangen |
| `dryrun_log.go`: `l.Info(would)` ohne `args` | überlebt (F-1) |
| `generate.go`: `"changed", 0` | überlebt (F-1) |
| `devcontainer_egress.go`: `"old", ""` | überlebt (F-1) |

**(4) Praxis.** Runtime-Images aus `git archive fa019d0` und `HEAD` gebaut, Läufe in gemounteten Verzeichnissen unter `~`; Images und Verzeichnisse danach entfernt. Vorher/Nachher stderr: `config set project.name renamed --dry-run` → `config set: updated` / `config set: would update` (`path`, `old`, `new` vorhanden); `generate devcontainer --dry-run` → `generate devcontainer: created` / `generate devcontainer: would create` (Felder `project`, `forwardPorts`, `features`, `changed` gleich). Ohne `--dry-run` und bei `--diff` allein in beiden Ständen `config set: updated` und `generate devcontainer: created`; `--dry-run --quiet` ohne Zeile; `u-boot.yaml` nach den Dry-Runs unverändert (`renamed2` erst nach dem echten Lauf). Hinweis: `init` brach im Runtime-Image ohne `git` ab (nach Anlage von `u-boot.yaml`); beide Stände gleich, ohne Einfluss auf die geprüften Aufrufe.

**(5) Doku.** Handbuch-Absatz in „Vorschau statt Überraschung“: keine Kennungen, Stand statt Chronik, wahr gegen Code (Beispielzeilen, `--diff` allein, `--quiet`). CHANGELOG `[Unreleased]` · Changed: nennt alle acht Zeilen, den Skript-Hinweis und `--json` als Maschinenvertrag; ehrlich. Alter Wortlaut im Baum (Suche nach `config set: updated`, `generate: created`, `allowlist updated`, `repaired`, `updated block`): Treffer nur in `config_t3_test.go` (gilt weiter, kein Dry-Run), CHANGELOG-Eintrag, Plan und Aktions-Namen (`updated-block`, `repaired-manual` als JSON-Werte, unverändert); `harness/replay.md` und Spec tragen den Logwortlaut nicht (`docs/user/cli-json-contract.md` existiert nicht).

**(6) Plan-Tabelle §3 gegen `git diff --stat` der vier Slice-Commits:** 6 Code-Dateien, CHANGELOG, Handbuch, Roadmap, Slice-Datei — jede Zeile hat einen Plan-Eintrag, jeder Eintrag (außer „unverändert“ und Closure-Zeile) einen Diff. „Unverändert“-Aussagen (`generate_rollback.go`, `doctor.go`, `config_t3_test.go`, `logger.go`) per gelesener Aussage bestätigt. `044e0c9` ist fremd (anderer Slice) und nicht Gegenstand.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/application` (alle Logger-Aufrufe) | geprüft, ohne Befund (Vollständigkeit) |
| `internal/adapter/**` (Logger, CLI-Summaries) | geprüft, ohne Befund (Summaries sind Vorgänger-Scope) |
| Kriterium `PreviewDryRun`, Legacy-Konstruktoren, `--quiet` | geprüft, ohne Befund |
| Architekturgrenze, Sprachvertrag, Suppression (§3.2), Gates (§3.6) | geprüft, ohne Befund (keine `//nolint`, keine Gate-Änderung, Zeilen englisch) |
| Handbuch, CHANGELOG | geprüft, ohne Befund |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Docker-Pfad und keine Importregel berührt (`dryrun_log.go` importiert nur Ports) |
| `make gates` | nicht im Reviewer-Lauf wiederholt: Verifier-Frage; `make test` lief je Mutation |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Test pinnt Zusage „Felder unverändert“ nicht · Kommentar beschreibt Vorzustand · Zeilenhinweis im Plan ungenau · Doc-Kommentar kennt neuen Parameter nicht · Toter Rückfallzweig ohne Test

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 und F-2 sind LOW mit Konventions-Anker.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Prüfung (u. a. `make gates`, `make doc-trace`) bleibt beim Verifier.
