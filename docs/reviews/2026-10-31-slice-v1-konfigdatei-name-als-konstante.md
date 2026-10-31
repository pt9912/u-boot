# Review-Report: [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/in-progress/slice-v1-konfigdatei-name-als-konstante.md) — 2026-10-31

**Review-Art:** Code gegen Plan, Entscheidungen und Hard Rules (Wartbarkeit). Die DoD-Frage gehört dem Verifier. Nichts vom Implementer ist übernommen: Zählung, Ersetzungen, Schutz-Test, Mutationen und Spec-Aussage sind am Code und an eigenen Läufen geprüft.

**Gegenstand:** Commit `836dea1b` gegen Basis `076437c2` (19 Dateien, 155 Einfügungen, 38 Löschungen; Arbeitsbaum am Ende außer diesem Report und den fremden ungetrackten Dateien eines anderen Vorgangs unverändert, keine Produktdatei verändert).

**Skill:** `.harness/skills/reviewer.md` @ HEAD `836dea1b` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-31

> **Zitier-Form.** Dieser Report friert ein: Kennung statt Adresse, `make <target>` statt Sensor-Link, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/in-progress/slice-v1-konfigdatei-name-als-konstante.md) (§1 Befund, Wege und Architect-Prüfung, §2 DoD, §3 Plan-Tabelle, §6 Risiken)
- [`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md) (Lint-Profil unverändert; Schutz als Test statt Lint-Regel)
- [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)
- `spec/architecture.md` (Abschnitt `hexagon/domain`), [`SPEC-034`](../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)
- `AGENTS.md` §3.2, §3.4, §3.6, §3.7, §3.9, §3.10; `harness/review.md`

**Sensor-Läufe:** `make test` rot mit Mutation 1 (`filepath.Join(baseDir, "u-boot.yaml")` in neuer Datei unter `internal/hexagon/application/`) und rot mit Mutation 2 (String-Literal in `cmd/uboot/` und Raw-String in `internal/adapter/driving/cli/`), danach grün ohne Mutation; die Mutationsdateien sind gelöscht, `git status` zeigt keine Produktänderung. `make verify-depguard` grün (alle 9 Regeln). `make docs-check` meldet 2 Befunde, beide in der fremden ungetrackten Datei eines anderen Vorgangs (Devcontainer-Ressourcenlimits); die Dateien des Commits und dieser Report tragen keinen Befund. Nicht gelaufen: `make gates` (Nachweis des Implementers; ein Teilbefund von `docs-check` liegt außerhalb des Diffs), `make test-docker` (kein Docker-Pfad berührt), `make image-scan` und `make govulncheck` (Image und Abhängigkeiten unverändert).

---

## Findings

Jedes Finding folgt dem §Output-Schema des Reviewer-Skills; die Spalten spiegeln es (`v6.18.1` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der neue Aufzählungspunkt in der Sicht-Spec sagt von `ConfigFileCandidates()`, sie liefere „bei jedem Aufruf eine frische Liste“, und nennt „die Kandidaten in Suchreihenfolge“. Kein Test ruft die Funktion auf (kein Treffer in `*_test.go`, kein Aufrufer im Nicht-Test-Code); der Schutz-Test pinnt nur das Literal. Eine Änderung zu einer geteilten Slice-Rückgabe bleibt grün, ein Aufrufer, der das Ergebnis verändert, verfälscht dann die Liste für alle. | `AGENTS.md` §3.10 (Bedingung in Spec ohne pinnenden Test); [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) | `spec/architecture.md` · „bei jedem Aufruf eine frische Liste“; `internal/hexagon/domain/config.go` · „Each call returns a fresh slice“ | ja — `make test` bleibt grün, wenn die Funktion eine geteilte Slice liefert | Spec-Bedingung ohne pinnenden Test |
| F-2 | LOW | Die Zählung „33 Pfad-Literale“ (Slice §1, §2, Commit-Nachricht) ist die Trefferzahl von `grep -rn '"u-boot\.yaml"'`; ein Treffer ist ein Kommentar (`doctor.go` Doc-Kommentar zu `Label`, Beispiel `"u-boot.yaml"`). Tatsächlich ersetzt sind 32 Literale; `doctor.go` trägt 3 statt 4 (Diff: `RelPath`, `Label`, `ReadFile`). Am HEAD liefert dasselbe grep 2 Treffer (Konstante und dieser Kommentar), nicht „eine Stelle“; der Test zählt Literale und ist grün. | `AGENTS.md` §3.10 (Aussage nennt nur, was gelesen wurde) | Plan §1 · „`doctor.go` (4)“; Plan §2 · „Zählung: eine Stelle statt 33“; `internal/hexagon/application/doctor.go` · „probe/read/parse-error messages (e.g. `\"u-boot.yaml\"`)“ | ja — `git grep -n '"u-boot.yaml"' 076437c2 -- internal cmd` gegen HEAD | Zählung zählt Kommentar als Literal |
| F-3 | LOW | Die Plan-Tabelle führt `harness/replay.md` als „unverändert, prüfen“ mit „nicht gelesen“; `AGENTS.md` §3.10 verlangt für „unverändert“ eine gelesene Aussage mit Fundstelle. Der Abgleich nach der Umsetzung hat die Zeile nicht berichtigt (die Datei ist im Diff unberührt). | `AGENTS.md` §3.10 (Datei „unverändert“ durch gelesene Aussage belegen) | Plan §3 · „Der Slice ändert kein Generator-Verhalten; nicht gelesen“ | nein — Review-only | Plan-Zeile „unverändert“ ohne gelesene Aussage |
| F-4 | INFO | Grenzen des Schutz-Tests (alle ohne Fehlalarm und ohne Vertragsbruch heute): Er erfasst nur ein `BasicLit` genau gleich `u-boot.yaml`; `"u-boot" + ".yaml"`, `fmt.Sprintf` und Teilpfade wie `".u-boot.yaml"` fallen durch. Der Guard `scanned == 0` zählt beide Wurzeln zusammen, eine leere Wurzel `cmd` bliebe unbemerkt (`internal` trägt genug Dateien). Der Test liegt im Paket `application_test`, scannt aber auch Adapter und `cmd/`; er läuft nur im Gesamtlauf `make test`, ein Teilfilter (`RUN`) überspringt ihn. Das entspricht der Zusage des Plans (Weg A, nur eigenständiges Literal). | Maintainability; Plan §6 („erfasst nur das Pfad-Literal als eigenständigen String“) | `internal/hexagon/application/config_filename_literal_test.go` · „val == configFileNameLiteral“ | ja — Mutationen von Hand | Schutz erfasst nur das eigenständige Literal |
| F-5 | INFO | Zuständigkeit Verifier: `make gates` grün und Golden-/Replay-Cases „unverändert grün“ sind Aussagen der DoD, die dieser Lauf nicht beantwortet. Zuständig: Verifier. | `harness/verification.md` | Plan §2 · „`make gates` grün (im Modus der Umgebung)“ | nein — Verifier-Aufgabe | DoD-Beleg nicht Gegenstand des Reviews |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Ersetzungen in den 15 Dateien unter `internal/hexagon/application/` (jede Hunk im Diff gelesen: Join-Pfade, `fileWrite.Path`, `fileTemplate.Path`, `RelPath`, Marker-Listen in `initproject.go` und `doctor_sandbox.go`, `changed`-Liste in `generate.go`) | geprüft, ohne Befund; jede Stelle ersetzt das Literal 1:1 durch `domain.ConfigFileName` (Wert `u-boot.yaml`), Verhalten und Ausgabewortlaut (Label, `changed`-Liste, Meldungstexte, Fehlerstrings, JSON) unverändert; Meldungstexte wurden nicht angefasst |
| Zählung am HEAD und Formen, die das Muster verfehlt (Raw-Strings, zusammengesetzte Pfade, `"u-boot"`/`".yaml"`) | geprüft, ohne Befund bis auf F-2; Adapter und `cmd/` tragen kein Pfad-Literal, nur Hilfetexte und Meldungen (Produkttext, im Plan ausgeschlossen) |
| `internal/hexagon/domain/config.go` (Konstante, `ConfigFileCandidates()` als Funktion statt Paket-Variable, `gochecknoglobals`) | geprüft, ohne Befund; Kommentare im Indikativ ohne Chronik und ohne Konjunktiv über Verworfenes; §3.9 (a) bis (c): Grep des Diffs auf „former/früher/originally/no longer/is gone/ohne X wäre“ ohne Treffer, kein Kommentar verweist auf ein umbenanntes Symbol, kein Helfer zwischen Doc-Kommentar und Funktion |
| Schutz-Test `config_filename_literal_test.go`: Wurzel-Ermittlung (`runtime.Caller` bis `go.mod`), Ausnahme nur für `config.go`, Testdateien und Kommentare nicht erfasst, `scanned == 0`-Guard | geprüft, ohne Befund bis auf F-4; Mutation 1 und 2 färben `make test` rot mit Fundstelle, Rücknahme grün |
| `spec/architecture.md`: Platz (Abschnitt `hexagon/domain`, nach „Weiterer Bestand“, vor „Erlaubte Imports“), Übereinstimmung mit Code (Datei, Konstante, Funktion, „den `init` anlegt“ gegen `devcontainerFileTemplates`/`initFromTemplate`) | geprüft, ohne Befund bis auf F-1; keine Welle, kein Slice, keine ADR-Kennung (§3.4) |
| Plan-Tabelle §3 gegen `git diff --stat` des Commits | geprüft: alle 15 Application-Dateien, `config.go`, Test, Spec und Slice stehen im Diff und in der Tabelle; Abweichungen nur F-2 und F-3 |
| Schichtregel: Imports `application` → `domain`; `downservice.go`, `egress_list.go`, `logsservice.go` erhalten den Domain-Import neu | geprüft, ohne Befund; `make verify-depguard` grün, `domain` importiert nichts Neues |
| Hard Rules §3.2 (kein `//nolint`), §3.6 (`.golangci.yml` unverändert, keine Lockerung), Sprachvertrag (Test- und Doc-Texte Englisch) | geprüft, ohne Befund |
| `make test-docker`, `make gates`, `make image-scan`, `make govulncheck` | nicht gelaufen (Gründe oben unter Sensor-Läufe) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Spec-Bedingung ohne pinnenden Test · Zählung zählt Kommentar als Literal · Plan-Zeile „unverändert“ ohne gelesene Aussage · Schutz erfasst nur das eigenständige Literal · DoD-Beleg nicht Gegenstand des Reviews

## Verdikt

**Merge-blockierend:** nein für HIGH (keines); F-1 (MEDIUM) ist vor dem Merge zu klären, entweder durch einen Test, der die frische Liste pinnt, oder durch Streichen der Bedingung aus der Spec. F-2 und F-3 sind Plan-Berichtigungen ohne Verhaltensauswirkung.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation; DoD-Konformität prüft der Verifier separat.
