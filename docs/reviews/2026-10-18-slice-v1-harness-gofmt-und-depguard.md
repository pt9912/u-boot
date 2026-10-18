# Review-Report: slice-v1-harness-gofmt-und-depguard — 2026-10-18 <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Review-Art:** Code (Plan-Konformität des Diffs gegen Slice, ADR, Lastenheft und `AGENTS.md`; Schwerpunkt Gate-Verschärfung ohne Lockerung, Beweis-Skript, Format-Commit nur Formatierung)

**Gegenstand:** Commit-Range `38f17a5..HEAD` (6 Commits, `d94da25` bis `f535c63`) des Slice `slice-v1-harness-gofmt-und-depguard` <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `f535c63` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

**Eingangs-Kontext:**

- Slice `slice-v1-harness-gofmt-und-depguard` (§1 Ziel und Abgrenzung, §3 Plan-Tabelle und gewählte Wege, §6 Risiken, Anhang A; §2 DoD nur als Scope-Quelle) <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md), [`LH-QA-004`](../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil), [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
- `AGENTS.md` §3.2, §3.3, §3.6, §3.7, §3.9, §3.10; `.golangci.yml`, `scripts/verify-depguard.sh`, `harness/sensors/verify-depguard.md`, `docs/maintainer/quality.md`, `harness/README.md`
- Sensor-Läufe des Reviewers (Docker-only, Baum nach jeder Probe zurückgesetzt, `git status` leer): `make lint` (grün) · `make test` (grün) · Vet mit Build-Tag `docker` im Image `u-boot:test-docker-tools` (`VET_OK`) · `golangci-lint fmt --diff ./...` im Image `u-boot:lint` (Exit 0, keine Änderung) · `make verify-depguard` (9 von 9 PASS, Exit 0) · Rot-Probe Formatter und Gegenprobe `domain-isoliert` (siehe Negativbefunde). Nicht gelaufen: `make test-docker` (der Diff ändert in den Docker-Tests nur den Umbruch eines Closure-Literals; Vet mit Tag `docker` belegt, dass das Paket baut), `make gates` als Ganzes (Teilziele einzeln gelaufen).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Format-Commit hat die Direktive `//nolint:unparam` ans Ende des Kommentarblocks gesetzt; ihre Begründung bricht dort mit Strichpunkt ab (`(postgres-only catalogue);`), der Satzrest („slice-v1-keycloak adds the first non-nil row …“) steht jetzt über der Direktive im Doc-Kommentar. Die Wirkung der Direktive ist unverändert (`make lint` grün, `unparam` meldet nicht); der Plan nennt die Verschiebung und „im Wortlaut unverändert“, nicht den abgebrochenen Satz. | `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist) und §3.9 | `internal/hexagon/application/addservice.go` · `//nolint:unparam // result is always nil today (postgres-only catalogue);` | nein — Review-only | Kommentar-Satz durch mechanische Umordnung getrennt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-2 | INFO | Bestandsbefund, nicht dieses Slice (Zeilen liegen vor `38f17a5`): Im Quellbaum stehen Inline-`nolint`-Kommentare, obwohl `AGENTS.md` §3.2 und `docs/maintainer/quality.md` §1 sie ausschließen: wirksame Direktiven in `addservice.go` (`unparam`), `upservice_portprobe_docker_test.go` (`gosec`), `template_list_test.go` (`containedctx`), `up_acceptance_test.go`, `logs_acceptance_test.go`, `down_acceptance_test.go` (je `dupl`); in `up.go`, `down.go`, `logs.go` stehen `// nolint:dupl` mit Leerzeichen, die keine Direktive sind. `nolintlint` ist in `.golangci.yml` nicht aktiviert, kein Gate fängt Drift in diesen Zeilen. Zuständig: Architect/Projektinhaber (Beobachtung, gegebenenfalls Folge-Slice mit Entscheidung, ob §3.2 oder der Bestand weicht). | `AGENTS.md` §3.2 | `grep -rn "nolint" internal cmd` (9 Treffer) | ja — `grep`; `make lint` bleibt grün | Inline-`nolint` im Bestand gegen Suppression-Verbot | <!-- d-check:ignore (Kennungen als Zitat im Lauf-Beleg) -->
| F-3 | INFO | Das Wegwerf-Paket trägt selbst einen `revive`-Befund (`package-comments`, `probe.go:1:1`); der Lauf von `make lint` endet im Fall `domain-isoliert` also auch ohne den `depguard`-Treffer rot. Die Entscheidung trägt allein die Beschreibungsprüfung (`grep -F` auf `domain must not depend on port`); die Gegenprobe zeigt, dass sie bei entfernter Regel FAIL meldet. Die Prüfung des Exit-Codes ist in diesem Fall nicht trennscharf. | Maintainability | `scripts/verify-depguard.sh` · `printf 'package %s\n' "$(basename "$PROBE_DIR")" > "$PROBE_DIR/probe.go"` | ja — `make verify-depguard` (Log von `domain-isoliert`: 2 issues, `depguard: 1`, `revive: 1`) | Begleitbefund im Probe-Paket maskiert Exit-Code-Prüfung | <!-- d-check:ignore (Kennungen als Zitat im Lauf-Beleg) -->
| F-4 | INFO | `docs/maintainer/quality.md` §1 nennt den Aufruf `golangci-lint fmt ./...` im Container ohne Nutzerangabe (Datei-Eigentümer nach dem Lauf auf Linux-Hosts nicht gemessen) und nicht, dass `gofmt` mit dem Default `simplify` läuft (der Plan nennt es); künftige Dateien werden also auch auf `simplify`-Muster geprüft. Kein Widerspruch zur Spec. | `AGENTS.md` §3.1 (Docker-only), Maintainability | `docs/maintainer/quality.md` · `golangci-lint fmt ./...` | nein — Review-only | Formatter-Aufruf ohne Nutzer-/Simplify-Hinweis |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Format-Commit `d94da25`: `git diff -w d94da25~1 d94da25` ganz gelesen | geprüft, ohne Befund bis auf F-1: 6 Einfügungen/13 Entfernungen ohne Whitespace sind Leerzeilen, ein `//`-Trenner in `erroremission.go`, das umbrochene `t.Cleanup`-Literal (gleiche Anweisung, gleiche Argumente) und die versetzte Direktive; kein Verhalten geändert |
| 29 Dateien im Diff-Stat gegen Anhang A | geprüft, ohne Befund (29 von 29, 170/177 wie im Plan); Gesamt-Diff 37 Dateien = 29 + 8 (`.git-blame-ignore-revs`, `.golangci.yml`, `quality.md`, `roadmap.md`, Slice, `harness/README.md`, `verify-depguard.md`, `verify-depguard.sh`), deckt sich mit Plan-Tabelle §3 |
| Build, Vet, Tests | geprüft, ohne Befund: `make test` grün, Vet mit Tag `docker` grün, `golangci-lint fmt --diff` ohne Änderung |
| Commit-Trennung §3.3 | geprüft, ohne Befund: Format-Commit allein, Gate-Änderung (`2265909`), Skript, Doku, Blame-Datei je eigene Commits; Slice-Commits tragen Slice-ID und `LH-QA-004` | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| `.git-blame-ignore-revs` | geprüft, ohne Befund: voller Hash `d94da2561469a7bdbec2d0f8b05fdee3315d37af` ist der Format-Commit, Kommentarform korrekt, `git blame --ignore-revs-file` überspringt ihn (Zeile mit Direktive blieb sinnvoll zugeordnet; die eingefügte Leerzeile erscheint als `d94da2561`, weil sie dort eingefügt wurde) |
| `.golangci.yml` | geprüft, ohne Befund: `git diff 38f17a5..HEAD` zeigt nur Kopfkommentar und `formatters.enable: [gofmt]`; keine `depguard`-Regel, kein Linter, keine Schwelle berührt; Schema `version: "2"` gültig für `v2.14.0` (Version aus `Dockerfile` und Lauf bestätigt); Verschärfung, keine Lockerung (§3.6) |
| Rot/grün-Probe Formatter | geprüft, ohne Befund: eine absichtlich falsch eingerückte, sonst saubere Datei in `internal/hexagon/domain` ließ `make lint` mit `File is not properly formatted (gofmt)` und Exit 2 enden; Datei gelöscht, Baum sauber, danach `make lint` grün |
| CI-Wirkung, Target-Index | geprüft, ohne Befund: `ci.yml` ruft `make gates` (enthält `lint`); kein neues Target, `harness/README.md` nennt keines hinzu, nur die Zeile `make lint` um den Formatter |
| `scripts/verify-depguard.sh` Logik | geprüft, ohne Befund: `trap` auf `EXIT INT TERM` räumt Verletzungsdatei und Probe-Verzeichnis; `rm -rf "$PROBE_DIR"` nur auf den gesetzten, in Anführungszeichen stehenden Pfad aus `CASES` (nie leer, nie ungequotet, nie ein vorhandenes fremdes Verzeichnis, weil ein vorhandenes mit Exit 2 abbricht, bevor `PROBE_DIR` gesetzt ist); Exit-Codes 0/1/2 wie dokumentiert; nach dem Lauf Baum sauber |
| Injektion trifft die Regel | geprüft, ohne Befund: Log von `domain-isoliert` zeigt `is not allowed from list 'domain-isoliert': domain must not depend on port (LH-FA-ARCH-003) (depguard)`, kein `import cycle` (zur zweiten Meldung siehe F-3) | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| Gegenprobe kaputtes Regelwerk | geprüft, ohne Befund: Klon aus `git archive HEAD` unter `~/uboot-review-depguard` (eigenes `git init`), Regelblock `domain-isoliert` aus `.golangci.yml` entfernt: `FAIL domain-isoliert: expected desc not found`, übrige 8 PASS, Exit 1; Klon und Logs gelöscht |
| `make verify-depguard` im Repo | geprüft, ohne Befund: 9 von 9 PASS, Exit 0 |
| Doku `quality.md`, `verify-depguard.md`, `harness/README.md` | geprüft, ohne Befund: „29 Linter“ bleibt wahr (5 + 24, der Formatter ist ausdrücklich kein Linter), `spec/spezifikation.md` Zeile „29 Linter insgesamt“ unverändert wahr; Vertrag (fünftes Feld, Exit 2 bei vorhandenem Probe-Verzeichnis, neun Regeln) stimmt mit dem Skript; CHANGELOG nicht nötig (kein Nutzer-Effekt) |
| Plan-Tabelle §3 gegen Code, „unverändert“-Zeilen | geprüft, ohne Befund: `repo.mk` und `Dockerfile` ohne Diff, Aussage zum `lint`-Target stimmt; §3.9 (a)/(b)/(c) außer F-1 ohne Treffer; Sprachvertrag: keine CLI-Ausgabe berührt, Skriptmeldungen und Skript-Kommentare Englisch, Slice- und Doku-Prosa Deutsch |
| Nicht geprüft | `make test-docker`, `make coverage-gate`, `make docs-check` über den Review-Stand (Letzteres läuft nach dem Schreiben dieses Reports); DoD-Häkchen und Closure-Notiz (Verifier-Aufgabe) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentar-Satz durch mechanische Umordnung getrennt · Inline-`nolint` im Bestand gegen Suppression-Verbot · Begleitbefund im Probe-Paket maskiert Exit-Code-Prüfung · Formatter-Aufruf ohne Nutzer-/Simplify-Hinweis

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; F-1 ist LOW mit Konventions-Anker und unverändert wirksamer Direktive.

**Übergabe:** F-1 und F-4 an den Implementer. F-2 ist ein Bestandsbefund und geht als Beobachtung ans Register (Entscheidung über Folge-Slice beim Architect bzw. Projektinhaber); er blockiert diesen Slice nicht. Die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Prüfung bleibt beim Verifier.
