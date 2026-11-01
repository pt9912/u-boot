# Verifikations-Report: slice-v1-hilfetexte-ohne-kennungen — 2026-11-01

**Rolle:** Verifier (eigener Kontext, Frage „Bauen wir es richtig?“), kein Produktcode geändert.
**Stand:** `761e8bd6` auf `main`, Diff-Basis `1805e54a`. Review-Report: [`2026-11-01-slice-v1-hilfetexte-ohne-kennungen.md`](2026-11-01-slice-v1-hilfetexte-ohne-kennungen.md).
**Entscheidungen des Projektinhabers:** Frage 1 Weg A, 3 im selben Slice, 4 mit; Frage 2 (Release) offen.

## Sensor-Belege

| Sensor | Beleg | Befund |
|---|---|---|
| `make gates` | Behauptung des Implementers (grün, Coverage 92,20 %); vom Verifier nicht neu gefahren | Gate-Beleg nicht selbst reproduziert; DoD-Punkt 4 bleibt offen |
| Wächter A, selbst gefahren | `make test PKG=./internal/adapter/driving/cli RUN='TestHelpTree\|TestOptionUsage\|TestUserTextFindings'` | grün |
| Wächter B, selbst gefahren | `make test PKG=./internal/hexagon/application RUN='TestMessage'` | grün |
| Mutation A | `(LH-FA-ADD-002)` in `Short` von `add` eingesetzt | rot: `help_identifiers_test.go:130: u-boot add [Short]: "LH-FA-ADD-002"`; zurückgesetzt (`git status` sauber) | <!-- d-check:ignore (wörtliches Zitat der Mutation bzw. des Fundes) -->
| Mutation B | `(LH-FA-INIT-006)` in `projectname.go:45` eingesetzt | rot: `message_identifiers_test.go:214 … projectname.go:45:25: LH-FA-INIT-006`; zurückgesetzt | <!-- d-check:ignore (wörtliches Zitat der Mutation bzw. des Fundes) -->
| `make test-docker` | nicht gelaufen (CI) | begründet; der Wortlaut der Warnzeile steht in `internal/e2e/devcontainer_shell_docker_test.go` (im Diff, 1 Zeile) |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (kein Eintrag in `.golangci.yml`, `Makefile`, `repo.mk` im Diff) |

Nicht selbst gefahren: die zwei Mutationen des Plans zu Ankündigung/Chronik (`remove.go:155` „v0.3.0 does NOT …“, `template` „Coming in later slices:“); gedeckt nur durch den Positivfall-Test `TestUserTextFindings_Patterns`, der diese Sätze als rot pinnt (grün).

## DoD §2, Punkt für Punkt

| Punkt | Urteil | Beleg im Code |
|---|---|---|
| **Hilfe** | erfüllt | Wächter A läuft über den Cobra-Baum (`HelpTextsForTest` in `export_test.go`; Mindestzahl 50 Texte als Schutz gegen leeren Lauf), Muster für Kennungen, Deutsch, Ankündigung, Chronik; Rot-Fälle und Grün-Fälle (`no longer needed`, `planned changes`, `blob/v`) gepinnt. `add`, `remove`, `template` gelesen (Diff): Katalog „postgres, keycloak or otel“ stimmt mit der Plan-Fundstelle `servicename.go` überein; `--purge`-Text stimmt mit `applyVolumePurge` überein; `template` ohne „Coming“-Abschnitt. „Exit codes:“ ohne Kennung; Sentinel-Namen durch Worte ersetzt. Nicht gelesen im Verifier-Kontext: Wortlaut je Befehl Alt gegen Neu außerhalb von `add`/`remove`/`template` (das hat der Reviewer geleistet). |
| **Meldungen** | erfüllt, mit Grenze | Wächter B (`go/parser`, Callees, `Message/Hint/MissingHint/Detail`, `warn/relax` zweites Argument, namensbasierte Konstanten) über fünf Wurzeln; `Code:`-Felder unberührt. Breite Suche des Verifiers über alle Nicht-Test-Literale (`LH-…`, `SPEC-`, `ADR-`, `slice-v`, `§`): übrig bleiben nur Codes (`return "LH-FA-GEN-002"` u. a.), `jsontestutil` (laut Plan ausgenommen) und kein Meldungstext. Suche nach deutschen Wörtern, Ankündigung, Chronik in Literalen: kein Treffer. Warnzeile Weg A: `up.go:214`, `devcontainer.go:555`, `generate.go:232`, `init.go:441` schreiben `Warning: %s`; Suche nach `Warning (` in `spec/`, `docs/user`, `docs/maintainer`, READMEs, `internal`, `cmd`: einziger Treffer ist der Negativ-Pin `sandbox_test.go:92`. | <!-- d-check:ignore (wörtliches Zitat der Mutation bzw. des Fundes) -->
| **Doku und Eintrag** | erfüllt | `CHANGELOG.md` `[Unreleased]` `### Changed` mit zwei Punkten (abweichend vom Plan-Wortlaut „ein Satz“: inhaltlich vollständiger, nennt Warnzeile und Backticks). `spec/spezifikation.md` Zeilen 887, 916 auf `Warning: <Text>` (Diff). README/`docs/user` unverändert; Suche nach `Warning (` ohne Treffer. |
| **Vorlagen** (Frage 3) | erfüllt | Suche nach `LH-`, `SPEC-`, `ADR-`, `slice-`, `§` in `internal/hexagon/application/templates`, `internal/adapter/driven/externaltemplates/templates` und `.devcontainer`: kein Treffer. Golden/Replay deckt `template_init_test.go` (im Diff, grün unter `make gates` laut Implementer, nicht selbst reproduziert). |
| **Backticks** (Frage 4) | erfüllt | `TestOptionUsageCarriesNoBacktickWord` grün; `root.go` Option `--quiet` ohne Backtick. |
| **`make gates` grün** | offen | nur Behauptung; Verifier hat es nicht gefahren. Der Haken fehlt zu Recht, bis der Lauf belegt ist. |
| **Review, Verification, Closure, Register** | offen | Review vorhanden (alle HIGH/MEDIUM behoben); dieser Report ist die Verification; Closure-Notiz §7 („Noch nicht geschlossen“) und Register `BEO-CLI/kennung-in-cli-ausgabe` stehen aus (Plan: bei der Closure). |

## Plan-Tabelle §3 gegen `git diff --stat 1805e54a..HEAD` (57 Dateien)

- Jede Code-, Test-, Vorlagen-, Doku- und Spec-Datei des Diffs hat einen Plan-Eintrag (Zeilen 107–123), auch `.devcontainer/*`, `export_test.go`, `message_identifiers_test.go`, `help_identifiers_test.go`, `cmd/uboot/main_test.go`, `removeservice.go`.
- Der Review-Report des Slice ist in der Plan-Tabelle gelistet (nachgezogen nach diesem Lauf).
- Plan-Eintrag ohne Diff: `observations/BEO-CLI/kennung-in-cli-ausgabe/` (Closure, erwartet); `cli-json-contract.md`, README, `docs/user` (als „unverändert“ geführt, korrekt).
- Die Plan-Zeile zum Lifecycle nennt den Pfad in `in-progress/` (nachgezogen nach diesem Lauf).

## Restfunde

| Nr. | Stufe | Befund | Fundstelle |
|---|---|---|---|
| V1 | LOW | Kommentar beschreibt abwesenden Text („the legacy „deferred“ prose is moot“, „no „deferred“ warning“), [`AGENTS.md`](../../AGENTS.md) §3.7, §3.9 (a); außerhalb der Wächter-Reichweite (Kommentare) | `internal/adapter/driving/cli/remove.go:479`, `internal/hexagon/application/removeservice.go:314` |
| V2 | LOW | Kommentar mit deutschem Rest und unverständlicher Klammer „(ohne entfernte Volumes)“ nach dem Streichen von „v0.3.0 deferred“ | `internal/adapter/driving/cli/remove.go:60` |
| V3 | LOW | Plan-Zeile 124 trägt den Pfad `open/…` statt `in-progress/…` | Slice §3 |
| V4 | INFO | Hilfetext von `remove --purge` sagt „the named volume of the service“; für Dienste ohne Volume unscharf (wie Reviewer I2) | `remove.go`, Option `--purge` |
| V5 | INFO | Grenzen der Wächter wie im Review I1: `Use`/`Aliases`, `:=`-Zuweisungen, `%s`-Einsetzung, `cmd/uboot`, `.tmpl` ungedeckt; „currently“, „for now“, „future“ nicht in der Liste. Die Suche des Verifiers fand dort keinen Treffer in Meldungs-Literalen | — |

Kein Fund in Hilfetexten oder Meldungs-Literalen (Kennung, Deutsch, Ankündigung, Chronik) außerhalb der Wächter.

## Urteil

Die Liefer-Punkte der DoD (Hilfe, Meldungen, Doku/CHANGELOG, Warnzeile, Vorlagen, Backticks) sind im Code erfüllt und von beiden Wächtern mit rot färbender Mutation belegt. Offen sind: der Gate-Lauf als reproduzierter Beleg (Verifier fuhr `make gates` nicht), `make test-docker` (CI), die Closure-Notiz samt Lerneintrag, das Register und Frage 2 (Release). Die drei LOW-Funde (V1 bis V3) blockieren die Closure nicht, sollten vorher berichtigt werden (V3 ist Pflicht nach §3.10).

## Nicht geprüft

`make gates`, `make test-docker`, `make verify-depguard`; die Alt/Neu-Lesung der Hilfen von `config`, `doctor`, `down`, `generate`, `init`, `logs`, `up`, `devcontainer`, `egress` Zeile für Zeile (Verifier las die Befehle `add`, `remove`, `template`, `root` im Diff und die Wächter über Suchen).
