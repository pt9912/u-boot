# Verification-Report: [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/done/slice-v1-konfigdatei-name-als-konstante.md) — 2026-10-31

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und den ausgeführten Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Die Vorlage `review-report.template.md` der vendored Baseline wurde per `cp` angelegt (`diff -q` gegen die Vorlage: identisch) und dann mit `Write` ersetzt, weil kein `Edit` zur Verfügung stand; Form nach den vorhandenen `-verify.md`-Reports. Außer dieser Datei ist nichts geändert, nichts committet; `git status` war nach jedem Mutationslauf leer.

**Gegenstand:** `git diff 076437c2..HEAD` (`HEAD` = `e2a012d0`; Implementierung `836dea1b`, Fix-Runde `a289e9b0`, Reviews `121e975c`, `e2a012d0`). Arbeitsbaum zu Beginn sauber. Fremde Vorgänge im Diff ([`slice-v1-devcontainer-ressourcenlimits`](../plan/planning/open/slice-v1-devcontainer-ressourcenlimits.md), `docs/messungen/`, `tools/messung/`, `docs/reviews/`) sind ausgenommen. Umgebung: `BUILD_MODE` nicht gesetzt, also Docker-Modus (Standard), `docker` vorhanden, kein Host-`go`.

**Verdikt:** DoD-Liefer-Punkte 1 bis 3 und `make gates`: **erfüllt**. Keine HIGH, kein MEDIUM, zwei LOW, zwei INFO. Die Closure-Punkte (Review-Häkchen, Notiz, Register, Risiken, Paarungen) liegen beim Planner. **Bereit für die Closure: ja**, nach Nachzug von F-1.

## Verification Evidence

Scope:
- Slice: [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/done/slice-v1-konfigdatei-name-als-konstante.md)
- IDs: [`LH-FA-CONF-002`](../../spec/lastenheft.md#lh-fa-conf-002--inhalt-der-konfiguration) (Name der Datei unverändert); keine ADR berührt ([`ADR-0003`](../plan/adr/0003-solid-nahes-lint-profil.md): Profil unverändert, `.golangci.yml` nicht im Diff)
- Artefakte: `internal/hexagon/domain/config.go`, `config_test.go`, 15 Dateien unter `internal/hexagon/application/`, `config_filename_literal_test.go`, `spec/architecture.md`

### 1. DoD-Abgleich

- [x] **Konstante und Ersetzung** — selbst gezählt: `grep -rn '"u-boot\.yaml"' --include=*.go internal cmd | grep -v _test.go` liefert genau 2 Treffer: `internal/hexagon/domain/config.go:7` (`const ConfigFileName = "u-boot.yaml"`) und `internal/hexagon/application/doctor.go:495` (Doc-Kommentar zu `Label`, kein Literal). `git diff 076437c2..HEAD -- internal/hexagon/application | grep -c '^-.*"u-boot\.yaml"'` = 32 entfernte Zeilen mit Literal, 15 geänderte Nicht-Test-Dateien (`doctor.go` hatte vorher 4 Treffer, 3 Literale + Kommentar). Verhalten: Diff Zeile für Zeile gelesen; jede Änderung ersetzt `"u-boot.yaml"` durch `domain.ConfigFileName` (gleicher Wert), sonst nur die Importzeile `domain` in `downservice.go`, `egress_list.go`, `logsservice.go`. Die Sonderstellen: `Label: domain.ConfigFileName` (`doctor.go`, Anzeigetext unverändert), `changed = append(changed, domain.ConfigFileName)` (`generate.go`, Ausgabeliste gleicher Wortlaut), Marker-Listen `[]string{domain.ConfigFileName, "compose.yaml", ".env.example"}` in `initproject.go` und `doctor_sandbox.go` (Reihenfolge unverändert). Keine Meldungstexte, Testdateien, Templates oder Golden-Dateien im Diff (Diff-Suche unter `testdata`/`golden*`: leer). Alle Tests grün (`make gates`, siehe unten).
- [x] **Schutz** — `internal/hexagon/application/config_filename_literal_test.go` (`TestConfigFileNameLiteralOnlyInDomain`, `go/parser`, `BasicLit` mit Wert `u-boot.yaml` in Nicht-Test-Dateien unter `internal` und `cmd`, nur `domain/config.go` ausgenommen; Wächter `scanned == 0` schlägt an). **Mutation selbst gefahren:** neue Datei `internal/hexagon/application/zz_mutation.go` mit `filepath.Join(x, "u-boot.yaml")` → `make test` Exit 2, `--- FAIL: TestConfigFileNameLiteralOnlyInDomain … string literal "u-boot.yaml" at /src/internal/hexagon/application/zz_mutation.go:5:60: use domain.ConfigFileName`. Datei entfernt, `git status --short` leer, danach `make gates` Exit 0. Weder `//nolint` noch Lockerung von Gate oder Schwelle (`.golangci.yml` unverändert), also AGENTS.md §3.2 und §3.6 nicht berührt.
- [x] **Konvention in der Sicht-Spec** — `spec/architecture.md`, §2.1 `hexagon/domain`, Aufzählungspunkt „Name der Konfigurationsdatei“ nach „Weiterer Bestand“ (ein Zeilen-Zusatz im Diff). Nennt Datei, `ConfigFileName`, `ConfigFileCandidates()` (frische Liste je Aufruf: gepinnt durch `TestConfigFileCandidates_ReturnsFreshSlicePerCall`; genau ein Eintrag: `TestConfigFileCandidates_HoldsOnlyConfigFileName`), den Test im Paket `application` (gepinnt, Mutation oben). Keine Welle, kein Slice, keine ADR, kein Commit im Satz (AGENTS.md §3.4).
- [x] **`make gates` grün** — selbst gefahren im Docker-Modus: Exit 0; `baseline-verify: v6.18.1 OK — 54 Dateien`, `coverage-gate: OK — coverage 92.10% meets threshold 90%`, `d-check: 930 Datei(en) geprüft, 0 Befund(e)`, `ablage-pruefung: OK`, `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`.
- [ ] Review durchgeführt — Reports liegen vor (Erstlauf 0 HIGH, 1 MEDIUM, 2 LOW, 2 INFO; Nachlauf mit zwei LOW, einem INFO, „Merge-blockierend: nein“). Das Häkchen setzt der Planner.
- [ ] Closure-Notiz, Register, Risiko-Ausgänge, Paarungen — offen, siehe Abschnitt 5.

### 2. Plan gegen Code

`git diff --stat 076437c2..HEAD` (nur dieser Slice) gegen die Plan-Tabelle §3:

| Diff-Datei | Plan-Zeile | Ergebnis |
|---|---|---|
| `internal/hexagon/domain/config.go` | neu | stimmt; trägt `ConfigFileName` und `ConfigFileCandidates()` |
| `internal/hexagon/domain/config_test.go` | neu (Fix-Runde) | stimmt |
| 15 Dateien `internal/hexagon/application/*.go` (Liste identisch mit der Tabelle) | update | stimmt, 15 von 15 |
| `internal/hexagon/application/config_filename_literal_test.go` | neu | stimmt |
| `spec/architecture.md` | update | stimmt |
| der Slice selbst | update | stimmt |
| `harness/replay.md`, `internal/hexagon/domain/configpath.go`, `.golangci.yml` | unverändert | nicht im Diff; stimmt |

Keine Diff-Zeile ohne Plan-Eintrag, kein Plan-Eintrag ohne Diff. Zahlen im Slice gegen Code: 32 Literale und 15 Dateien stimmen (s. o.); `doctor.go` 3 Literale + 1 Kommentar stimmt; „am Stand `836dea1b` 2 Treffer“ stimmt am HEAD ebenfalls.

### 3. Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make test` mit Mutation | fail (erwartet) | Exit 2, Test benennt die Fundstelle |
| `make gates` | pass | Exit 0, Zeilen oben |
| `make verify-depguard` | pass | Exit 0, „all 9 depguard rules verified“, Arbeitsbaum danach sauber |
| `make docs-check` | pass | in `make gates`: 0 Befunde; nach Anlage dieses Reports erneut gefahren (siehe Schluss) |

Nicht ausgeführt:
- `make test-docker` — kein Docker-Pfad berührt: nur Konstante gleichen Werts, keine Runtime-, Compose- oder E2E-Änderung; Begründung trägt die Sensor-Auswahl in `harness/verification.md` (Go-Code ohne Docker-Runtime: `make gates`).
- `make ci`, `make fullbuild` — kein Security-/Release-Pfad berührt (kein Abhängigkeits-, Image- oder Workflow-Diff).
- `make coverage-gate` einzeln — Teil von `make gates`.

### 4. Befunde

| ID | Schwere | Befund |
|---|---|---|
| F-1 | LOW | `docs/plan/planning/in-progress/roadmap.md` Zeile 20 führt weiter „statt 33 Pfad-Literale“; Slice-Titel und §2 sagen 32 (Zählung berichtigt in der Fix-Runde). Die Roadmap ist nicht im Diff und wurde beim Zählungs-Nachzug nicht erfasst. |
| F-2 | LOW | Slice §2 Punkt 1 hält „Golden-/Replay-Cases unverändert grün“ als Beleg, aber der Slice nennt keinen eigenen Replay-Lauf; belegt ist es nur durch `make test`/`make gates` (Golden-Dateien nicht im Diff). Genügt als Nachweis, sollte in §7 so benannt werden. |
| I-1 | INFO | Spec-Satz „(der Name, den `init` anlegt)“: Kein Test pinnt ausdrücklich, dass `init` den Namen aus der Konstante nimmt; mittelbar gedeckt durch die Bestandstests auf `u-boot.yaml` und den Literal-Test (kein anderer Weg zum Namen im Nicht-Test-Code). |
| I-2 | INFO | Der Schutz erfasst nur Literale mit exakt dem Wert `u-boot.yaml`; Zusammensetzungen oder Raw-Strings mit Schrägstrich fallen durch (Risiko 4 im Slice). Gemessen: Raw-Strings mit dem Namen kommen nur in Kommentaren und CLI-Hilfetexten (`adapter/driving/cli/config.go`) vor; ein Raw-String `` `u-boot.yaml` `` würde `strconv.Unquote` dagegen als Wert erfassen. |

### 5. Offene Schuld für die Closure (Planner schreibt; hier nichts geschrieben)

- DoD-Häkchen: Review durchgeführt (Reports Erstlauf und Nachlauf liegen vor); Closure-Notiz; Register; Risiko-Ausgänge; Paarungen. `make gates` ist abgehakt, Verifikation liegt mit diesem Report vor (§5 Closure-Trigger nennt ihn).
- §7 ist Platzhalter: Was funktioniert hat / anders lief (Zählung 33 → 32 nach Review, Test zählt Literale, Kommentar-Treffer in `doctor.go`); Steering-Loop-Lerneintrag (z. B. Zählung per Befehl und Kategorie Literal/Kommentar trennen; Spec-Bedingung ohne Test strich der Nachlauf); Folge-Slice-Zeiger auf [`slice-v1-konfigdatei-alternativer-name`](../plan/planning/open/slice-v1-konfigdatei-alternativer-name.md).
- Beobachtungs-Register: §8 kündigt „kann als neue Beobachtung angelegt werden“ an; entweder `BEO-APP/<slug>` mit `evidence/` anlegen oder „keine Beobachtung angefallen“ in §7 notieren.
- Risiken §6, alle vier mit „Ausgang: bei Closure“, keines hat einen Ausgang: (1) Verhaltensänderung — entfallen (Diff nur Wert-gleiche Ersetzung, Tests grün); (2) Schutz schlägt falsch an — entfallen (Test grün auf 506/251-Zeilen-Bestand, nur exaktes Literal); (3) Test ohne Mutation — entfallen (Mutation gemessen, oben); (4) Zählung unvollständig — weiter offen oder entfallen laut Wahl: Raw-String ist abgedeckt, Zusammensetzung mit Schrägstrich nicht (I-2). Die Ausgänge setzt der Planner.
- Paarungen (Anker · Folge-Slice · Register): Anker nur, wenn der Eintrag `liegt in` trägt (die Konstante/der Test wäre ein Zielort mit `seit` plus der Kennung dieses Slice, [`slice-v1-konfigdatei-name-als-konstante`](../plan/planning/done/slice-v1-konfigdatei-name-als-konstante.md)); Folge-Slice [`slice-v1-konfigdatei-alternativer-name`](../plan/planning/open/slice-v1-konfigdatei-alternativer-name.md) liegt in `open/`.
- Nachzug F-1 (Roadmap-Zeile 20 auf 32). Moves nach AGENTS.md §3.3 (Inhalt zuerst, dann reiner `git mv` nach `done/`).
- Carveouts: keine neu, keine gelöst, keine berührt.

Commit / Artefakt: `836dea1b` (Implementierung), `a289e9b0` (Fix-Runde), `e2a012d0` (Nachlauf, `HEAD`).
