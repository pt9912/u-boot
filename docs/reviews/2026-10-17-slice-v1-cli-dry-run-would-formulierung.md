# Review-Report: slice-v1-cli-dry-run-would-formulierung — 2026-10-17 <!-- d-check:ignore (Slice-Kennung im Titel) -->

**Review-Art:** Code — gegen Plan, Lastenheft und Hard Rules (Maintainability).

**Gegenstand:** Commit-Range `874b6b8..HEAD` (`73f1f4b`, `1536c8e`), Slice `slice-v1-cli-dry-run-would-formulierung`, Lösung A. <!-- d-check:ignore (Slice-Kennung) -->

**Skill:** `.harness/skills/reviewer.md` @ `1536c8e` · <!-- d-check:ignore (Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-cli-dry-run-would-formulierung` (§1 Lösung A, §2 DoD, §3 Plan-Tabelle) <!-- d-check:ignore (Slice-Kennung) -->
- [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)
- `AGENTS.md` §3.1, §3.6, §3.7, §3.9, §3.10, §5; `harness/review.md`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der CHANGELOG-Satz „… `--json` output, errors and `--quiet` are unchanged“ ist für `--quiet --dry-run` missverständlich: Sichtbarkeit der Zusammenfassung bleibt, ihr Wortlaut ändert sich aber (die Tests pinnen „Would …“ auch mit `--quiet`). | `AGENTS.md` §5 Zeile 5 (Doku-Präzision) | `CHANGELOG.md` · „errors and `--quiet` are unchanged“ | Review-only | Doku-Aussage mehrdeutig zum Verbosity-Zusammenspiel |
| F-2 | INFO | Der `default`-Zweig von `printGenerateSummary` (künftige `GenerateAction`) kennt `dryRun` nicht; eine neue Aktion würde im Dry-Run wieder in Vergangenheitsform ausgeben, ohne dass ein Test es fängt. | Maintainability | `generate.go` · „Defensive: a future GenerateAction value“ | Review-only | Defensiver Zweig nicht dry-run-aware |
| F-3 | INFO | `config set` mit mehreren Paaren: gepinnt ist nur die Zeile des zweiten Paars; ein gemischter Aufruf (ein Paar unverändert → „already … no changes.“, eines geändert → „Would set“) ist nicht getestet. Der Code rendert je Paar einheitlich. | [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) | `dry_run_would_test.go` · „config set many“ | ja — `make test` (Mutationstest nötig) | Mehrpaar-Mischfall nicht gepinnt |
| F-4 | INFO | `TestDryRun_JSONAndErrorsUnchanged` prüft „kein Would“ per Teilstring auf der gesamten JSON-Ausgabe; ein künftiger Diagnostic-Text mit dem Wort „Would“ würde den Test ohne Vertragsbruch brechen. Eine sinnvolle Mutation der JSON-Seite gibt es nicht, weil der JSON-Zweig vor jedem Summary-Printer zurückkehrt. | Maintainability | `dry_run_would_test.go` · `strings.Contains(out, "Would")` | Review-only | Teilstring-Negativ-Assert zu breit |
| F-5 | INFO | Der Implementer meldet einen einmaligen Host-`python3`-Aufruf (Verstoß gegen §3.1). Ergebnis geprüft: Der Diff von `remove_acceptance_test.go` umfasst genau die Kommentar-Berichtigung und eine zusätzliche Headline-Assertion; keine Neben-Änderungen. Der Nachweis lief über `make test`, der Host-Lauf ist kein Handoff-Nachweis. Für den Verifier: Verstoß ohne Wirkung auf das Ergebnis, im Slice-Bericht zu benennen. | `AGENTS.md` §3.1 | `remove_acceptance_test.go` · „Dry-Run headline must read Would remove“ | nein — Host-Aufruf nicht nachprüfbar, Ergebnis per Diff | Host-Werkzeug außerhalb Docker (einmalig) |

## Prüfung der Schwerpunkte

1. **Vollständigkeit.** `--dry-run` ist nur in `init`, `add`, `remove`, `generate`, `config set` registriert; `config show/get/list` weisen es per `ErrDryRunNotApplicable` ab, `up`/`down`/`doctor`/`status`/`logs`/`template` haben es nicht. `init`/`add` sagen bereits „Would …“. Alle Summary-Pfade mit Vergangenheitsform sind abgedeckt, einschließlich `generate` aller drei ändernden Aktionen, `config set` einzeln und über `runConfigSetMany`. Kein weiterer Pfad.
2. **Kriterium.** `printRemoveSummary` nutzt `previewMode == PreviewDryRun`; `generate`/`config set` `flags.DryRun`. Beides deckt sich mit `previewModeFromFlags` (`--dry-run` gewinnt, `--diff` allein = `PreviewAndApply` bleibt Vergangenheitsform). `--json` und Fehler kehren vor den Printern zurück; No-op-Zeilen sind in beiden Modi identisch; `summaryVisible` unverändert; Exit-Codes unberührt.
3. **Wortlaut.** `Would remove service` + `Would change:` ist wortgleich zu `add` (`Would add service` / `Would change:`); `Would generate/update/repair` folgen dem `Would repair service` von `add`; `config: Would set <path> <old> → <new>.` je Paar. Stimmig. Das Wort „Would“ im Englischen ist sprachvertragskonform ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)).
4. **Mutationen (`make test`, je zurückgenommen, Baum danach sauber):** M1 `runConfigSetMany` übergibt `false` statt `flags.DryRun` — rot (`TestDryRun_SummarySaysWould/config_set_many`, Quiet-Matrix); M2 `Repaired`-Zweig ohne Dry-Run — rot (`TestDryRun_GenerateActions/repaired-manual`); M3 `remove` ohne `Would change:` — rot (`TestDryRun_SummarySaysWould/remove`, `TestRemove_PurgeDryRun_…`); M4 Kriterium `!= PreviewNone` (feuert bei `--diff` allein) — rot (`TestDryRun_ActualWriteKeepsPastForm/remove/--diff`); M5 fehlendes Leerzeichen in „Would set“ — rot; M6 `config set` einzeln immer Dry-Run — rot (mehrere Tests). Die Tests sind nicht tautologisch. JSON-Seite siehe F-4.
5. **Doku.** Handbuch-Absatz nennt Zustand (keine Kennungen, keine Chronik), stimmt mit Code und Tests überein. CHANGELOG behauptet nichts Ungeliefertes und weist Skripte auf `--json` hin (F-1 zur `--quiet`-Formulierung). Alter Wortlaut im Baum (README, Handbuch, Spec/Spezifikation, `harness/replay.md`, `docs/maintainer/cli-json-contract.md`, Tests): keine veraltete Dry-Run-Beispielausgabe; verbleibende Treffer in `remove_test.go`, `cli_test.go`, `quiet_modifying_test.go` sind Schreib-Läufe ohne `--dry-run` und korrekt. Der Kommentar „Pre-T7 zeigte …“ in `remove_acceptance_test.go` ist vorbestehend und nicht Teil des Diffs.
6. **Plan-Tabelle (§3.10).** Gegen `git diff --stat 874b6b8..HEAD`: alle neun Dateien außer dem Slice-Plan selbst haben einen Eintrag, kein Eintrag ohne Diff. §3.9 (a)/(b)/(c): keine Chronik-Kommentare, keine veralteten Symbolverweise (Suche nach altem Wortlaut im ganzen Baum), Helfer stehen nicht zwischen Doc-Kommentar und Funktion. Die neuen Doc-Kommentare stehen im Indikativ.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapter/driving/cli/remove.go`, `generate.go`, `config.go` | geprüft, ohne Befund (außer F-2) |
| `quiet_preview_test.go`, `remove_acceptance_test.go` | geprüft, ohne Befund |
| `docs/user/benutzerhandbuch.md` | geprüft, ohne Befund |
| Dual-Classifier / Exit-Codes / Sentinels | geprüft, unberührt, ohne Befund |
| Referenzmodell (Rang 1–3 verweisen nicht abwärts) | geprüft, Spec/Architektur unverändert, ohne Befund |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft, Entfernung „Nichts in Arbeit.“ entspricht Plan §3, ohne Befund |
| `make test-docker`, `make verify-depguard` | nicht ausgeführt: kein Docker-Pfad, keine Importregel berührt |
| `make gates` | nicht ausgeführt (Verifier-Aufgabe); `make test` lief grün |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Doku-Aussage mehrdeutig zum Verbosity-Zusammenspiel · Defensiver Zweig nicht dry-run-aware · Mehrpaar-Mischfall nicht gepinnt · Teilstring-Negativ-Assert zu breit · Host-Werkzeug außerhalb Docker (einmalig)

## Verdikt

**Merge-blockierend:** nein.

**Übergabe:** Findings gehen an den Implementer; F-5 zusätzlich als Hinweis an den Verifier. Die Finding-Klassen gehen in die Slice-Closure §7. DoD-Konformität prüft der Verifier.
