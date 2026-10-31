# Review-Report: [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/in-progress/slice-v1-konfigdatei-name-als-konstante.md) (Nachlauf zur Fix-Runde) — 2026-10-31

**Review-Art:** Code gegen Plan, Entscheidungen und Hard Rules (Wartbarkeit). Nachlauf im frischen Kontext zur Fix-Runde der Befunde F-1 bis F-3 des ersten Reports. Die DoD-Frage gehört dem Verifier.

**Gegenstand:** Commit `a289e9b0` gegen Basis `836dea1b` (4 Dateien: Slice-Text, `internal/hexagon/domain/config.go`, `internal/hexagon/domain/config_test.go` neu, `spec/architecture.md`).

**Skill:** `.harness/skills/reviewer.md` @ HEAD `a289e9b0` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-31

> **Zitier-Form.** Dieser Report friert ein: Kennung statt Adresse, `make <target>` statt Sensor-Link, Baseline-Stellen als Tag + Pfad in Inline-Code. Das Werkzeug zum Ausfüllen war `Write` (kein Edit-Werkzeug); die Datei ist per `cp` aus der Vorlage entstanden (`diff -q` vor dem Ausfüllen: gleich).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/in-progress/slice-v1-konfigdatei-name-als-konstante.md) und der erste Report [erste Report](2026-10-31-slice-v1-konfigdatei-name-als-konstante.md)
- [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration)
- `spec/architecture.md` (Aufzählungspunkt „Name der Konfigurationsdatei“ unter `hexagon/domain`)
- `AGENTS.md` (Hard Rules §3.2, §3.4, §3.6, §3.7, §3.9, §3.10), `harness/review.md`, `harness/replay.md`

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Slice-Text führt nach der Berichtigung auf 32 (§2 Befund) weiter die Zahl 33 im Gegenwartston: Titel („statt 33 Literale“), Ziel („Änderung an 33 Stellen“), Weg A („die 33 Stellen greifen darauf zu“), Zeile des Architect-Befunds („33 Literale … bestätigt“, „gehören zu den 33“), Weg C („Suche über 33 Stellen“), §3-Tabellenzeile zu den 15 Anwendungsdateien („33 Pfad-Literale durch die Konstante ersetzen“) und §4 („Eine der 33 Stellen“). Nur §2 und §5 tragen die Berichtigung; ein Leser, der Titel, Ziel oder die Plan-Tabelle liest, erfährt 33. | `AGENTS.md` §3.10 (Tabelle gegen Stand nach der letzten Fix-Runde abgleichen) | Slice-Datei · „33 Pfad-Literale durch die Konstante ersetzen“ (§3) und Titelzeile „statt 33 Literale“ | ja, `grep -n '33' <slice>` | Zählung nach Berichtigung nur teilweise nachgezogen |
| F-2 | LOW | Der Satz in `spec/architecture.md` behauptet „Jeder Pfad zur Konfigurationsdatei greift darauf zu“ (Konstante oder Funktion). Gepinnt ist nur die Negativseite (Test im Paket `application`: kein Pfad-Literal außerhalb `domain/config.go`); `ConfigFileCandidates()` hat im Nicht-Test-Code keinen Aufrufer (`grep -rn ConfigFileCandidates internal cmd` trifft nur Definition und `config_test.go`). Ein aus Teilstrings zusammengesetzter Name oder ein Pfad ohne Zugriff auf die Funktion fiele nicht auf. | `AGENTS.md` §3.10 (Bedingung „immer/jeder“ in der Spec braucht einen pinnenden Test oder wird gestrichen) | `spec/architecture.md` · „Jeder Pfad zur Konfigurationsdatei greift darauf zu“ | ja, `make test` (kein Test pinnt die Positivseite) | Spec-Allquantor nur negativ gepinnt |
| F-3 | INFO | In `TestConfigFileCandidates_ReturnsFreshSlicePerCall` ist die Zeile `_ = append(first[:0], "a", "b")` für das Ergebnis wirkungslos (Kapazität 1, die Zuweisung davor trägt den Test); schadet nicht. | Maintainability | `internal/hexagon/domain/config_test.go` · `_ = append(first[:0], "a", "b")` | nein, Review-only | Test-Zeile ohne Wirkung |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Spec-Satz „Name der Konfigurationsdatei“, jede Bedingung gegen Code und Test: „genau eine Stelle in `domain/config.go`“ (Test `TestConfigFileNameLiteralOnlyInDomain` im Paket `application_test`, Schutz-Test aus dem Erstlauf; am HEAD genau 2 Treffer, Konstante und Doc-Kommentar zu `Label`), „Konstante `ConfigFileName`“ (Wert gepinnt in `config_test.go`), „Liste der Dateinamen, die die Konfiguration tragen darf“ und „frische Liste je Aufruf“ (beide Tests), „Meldungstexte und Kommentare nennen den Namen weiter als Wort“ (Test-Doc-Kommentar grenzt das ab). Die gestrichene Wendung „Suchreihenfolge“ steht nirgends mehr in Spec und Doc-Kommentar. Ausnahme: F-2 | geprüft, mit F-2 |
| Spec-Satz ohne Wellen, Slices, ADR, Commit-Hashes (§3.4); ohne Slice-Kenntnis verständlich | geprüft, ohne Befund |
| Mutation (a): `ConfigFileCandidates()` gibt geteilte Paketvariable zurück; `make test` | rot (`TestConfigFileCandidates_ReturnsFreshSlicePerCall` und `…_HoldsOnlyConfigFileName` schlagen an); zurückgenommen |
| Mutation (b): zweiter Eintrag `.u-boot.yml`; `make test` | rot (beide Tests); zurückgenommen. Danach `git status --short` und `git diff --stat` leer, `make test` grün |
| Kommentare `config.go`, `config_test.go` (§3.7, §3.9 a bis c): Indikativ, geltende Zusage („The list holds [ConfigFileName] only“, „Each call returns a fresh slice“, vom Test gepinnt); kein Rückblick, keine Verweise auf geänderte Symbole, keine Helfer zwischen Doc-Kommentar und Funktion | geprüft, ohne Befund |
| Kommentare im Slice-Text der Fix-Runde: „berichtigt in der Fix-Runde nach dem Review“ in §2 und die Zeilen der Plan-Tabelle sind Planungsprosa mit auflösbarem Beleg (Stand, Befehl), keine Code-Kommentare; Zählung selbst stimmt: 6+6+3+4+2+2+9 = 32 Literale, 33 Treffer am Stand `076437c2` (3 Literale und 1 Doc-Kommentar in `doctor.go`), 2 Treffer am HEAD | geprüft, ohne Befund |
| Zählung 33 versus 32: Die verbleibenden „33“ sind teils zulässige Wiedergabe der Planung (Architect-Befund vom Planungsstand, Weg C, Trigger), täuschen aber Leser ohne Verweis auf §2, wo Titel, Ziel und Plan-Tabelle als geltende Aussage stehen; siehe F-1 | geprüft, mit F-1 |
| Plan-Tabelle §3 gegen `git show --stat a289e9b0`: Slice-Datei, `domain/config.go`, `config_test.go` (neue Zeile „neu (Fix-Runde)“), `spec/architecture.md` sind alle geführt; kein Diff-Eintrag ohne Zeile, keine Zeile ohne Diff. Zeile `harness/replay.md` „unverändert“: gelesen — §Scope (Zeilen 8 bis 19) bindet Replay an Pfade, die User-Artefakte erzeugen; §Verification Evidence (ab Zeile 86) verlangt den Replay-Block nur bei geändertem Generatorverhalten; beide Stellen tragen die Aussage | geprüft, ohne Befund |
| Suppression/Gate-Lockerung (§3.2, §3.6): kein `//nolint`, `.golangci.yml` und Schwellen unberührt | geprüft, ohne Befund |
| Architekturgrenze: `config_test.go` im Paket `domain_test` importiert nur Standardbibliothek und `domain` | geprüft, ohne Befund |
| Provenienz des Reports: `cp` aus Vorlage, `diff -q` gleich | geprüft, ohne Befund |
| `make docs-check` | 0 Befunde |
| `make lint`, `make gates`, `make test-docker`, `make verify-depguard` | nicht gelaufen: `lint` optional und der Diff ändert nur eine Kommentarzeile plus Testdatei; `gates` ist Sache des Verifiers; kein Docker- oder Importregel-Pfad berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zählung nach Berichtigung nur teilweise nachgezogen · Spec-Allquantor nur negativ gepinnt · Test-Zeile ohne Wirkung

## Verdikt

**Merge-blockierend:** nein — F-1 bis F-3 des Erstlaufs sind behoben (Test wirksam, Spec-Wendung gestrichen, Fundstelle zu `harness/replay.md` gelesen); die zwei LOW sind Nachzug-Punkte ohne Vertragsbruch.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation — DoD-Konformität und `make gates` prüft der Verifier.
