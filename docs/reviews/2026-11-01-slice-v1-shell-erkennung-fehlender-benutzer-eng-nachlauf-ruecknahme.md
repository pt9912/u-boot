# Review-Report: slice-v1-shell-erkennung-fehlender-benutzer-eng, Nachlauf zur Rücknahme — 2026-11-01

**Review-Art:** Code — Code-Review des Rücknahme-Diffs gegen Plan, Spec, ADR und Konventionen (kein Self-Review, frischer Kontext).

**Gegenstand:** `git diff 350cdc72 HEAD -- internal/ spec/ docs/user docs/maintainer docs/plan/adr CHANGELOG.md` sowie die Commits `34357bbf` (Revert), `ad0e824f` ([ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) Deprecated), `8db00a3e` ([SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Wartungsdoku, Slice), `63257933`, `8a7ce7c0` und die Merge-Commits. Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (§7 Closure).

**Skill:** `.harness/skills/reviewer.md` @ `172387b5`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (§3 Plan-Tabelle, §6 Risiken, §7 Closure)
- Messprotokoll [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md) mit Rohdaten
- [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md), [`MR-008`](../../harness/conventions/MR-008-adr-form.md)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Tabelle der Exit-Codes, Unterabschnitt `shell`, Grenzen der Erkennung, Änderungsprotokoll)
- `AGENTS.md` (Hard Rules, besonders §3.5, §3.7, §3.9, §3.10)

**Sensors:** `make docs-check` lief: 971 Dateien, 0 Befunde. `git diff 350cdc72 HEAD -- internal/adapter/driven/docker/` ist leer (0 Zeilen). `make gates` und `make test-docker` liefen nicht: Der Auftrag schließt `make test-docker` aus, und der Code-Diff gegen die Basis ist im Paket leer, so dass `make gates` hier nichts prüfen würde, was nicht schon am Basisstand grün war (Verifier-Aufgabe, falls gewünscht).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Das Messprotokoll führt drei Aussagen im Stand vor der Rücknahme: Zeile „Ob Code und die Folge-Entscheidung ([ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md)) daraufhin angepasst werden, entscheidet der Auftraggeber; diese Messung ändert beides nicht“ (entschieden, ADR ist `Deprecated`); der Pfad sei in der Spezifikation „ohne die Ausnahme“ beschrieben (die Spezifikation trägt sie jetzt); der „dritte Aufruf (`userMissingNow`)“ als gegenwärtig („auf der Engine 29.9.2 ist der Pfad inaktiv“, der Code ist entfernt). Der nächste Leser des Protokolls liest einen offenen Zustand, der abgelaufen ist. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) und §3.9 (b) | [Messprotokoll](../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md) · „entscheidet der Auftraggeber; diese Messung ändert beides nicht“ | nein — Review-only (`make docs-check` sieht keinen Zustandsmarker) | Zustandsmarker im Messprotokoll nach Ablauf nicht nachgeführt |
| F-2 | LOW | Ein Kommentar im Messskript verweist auf das Symbol `userMissingNow`, das im Baum nicht mehr existiert; die beiden Rohdateien tragen dieselbe Überschrift als Ausgabe des Skripts. Eine Suche nach dem alten Namen im ganzen Baum (§3.9 b) trifft das Skript. | `AGENTS.md` §3.9 (b) | `tools/messung/shell-benutzer-fehlt/messen.sh` · „## G. Die erneute Probe nach der Sitzung (userMissingNow)“ | nein — Review-only (`grep -rn userMissingNow`) | Verweis auf entfernten Namen in Skript |
| F-3 | LOW | Die zweite Zeile der Tabelle „Geschichte“ in [ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) hat vier Zellen in einer Tabelle mit drei Spalten (`Datum \| Ereignis \| Verweis`); der Verweis „Slice zur Erkennung eines fehlenden Benutzers bei `shell`“ steht doppelt. | `MR-008`; ADR-Vorlage (Spalten `Datum`, `Ereignis`, `Verweis`) | `docs/plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md` · „Projektinhaber nimmt den dritten Aufruf und diese ADR zurück \| Messprotokoll, Slice zur Erkennung“ | nein — Review-only | Tabellenzeile mit überzähliger Zelle |
| F-4 | LOW | Die Tabelle der Exit-Codes nimmt in Spalte `0` „die Grenzen der Erkennung“ aus, der Punkt *Exit-Code `12`* im selben Unterabschnitt sagt unbedingt „Der Exit-Status der Shell nach dem Beginn der Sitzung ist kein Fehler … der Befehl endet mit `0`“ und verweist für `12` auf *Grenzen der Erkennung*, nennt aber die Ausnahme zu `0` nicht. Die Stellen sind nicht widersprüchlich im Code, aber untereinander nicht deckungsgleich. | `AGENTS.md` §3.10 (Bedingung in der Spec) | `spec/spezifikation.md` · „Der Exit-Status der Shell nach dem Beginn der Sitzung ist kein Fehler: Er erreicht den Aufrufer nicht“ | nein — Review-only | Ausnahme an einer von mehreren „unbedingt 0“-Stellen nachgetragen |

## Beantwortung der Prüfpunkte

1. **Code auf dem Stand vor der Slice-Arbeit.** `git diff 350cdc72 HEAD -- internal/adapter/driven/docker/` ist leer. `userMissingNow`, `user-gone-late` und `probe-late-not-running` kommen im ganzen Baum nur noch vor in: Messprotokoll (F-1), Messskript und Rohdaten (F-2), [ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) (eingefrorener Text), dem Slice in `done/` (als Aufzählung des Entfallenen) und den Berichten vom Stand mit Weg A. Weder in `internal/` noch in Tests noch in Kommentaren des Codes. Die Differenz in `internal/` gegenüber `350cdc72` liegt allein in `internal/e2e/` (siehe F-7, anderer Slice).
2. **Absatz der Spezifikation.** Die Lesung des Codes stimmt. `runSession` setzt bei `ExitError` `exit` auf den Status; `clientErrorBeforeSession` gibt bei `exit == 0` `false` zurück (deckt „Endet die Shell danach mit Status `0` … gilt die Sitzung“); `clientMessageOnStderr` verlangt ohne Terminal `stdout.total == 0`, eine erste Zeile der Fehlerausgabe mit `clientDaemonErrorPrefix` und `missingUserLead(opts.User)` (Name mit Doppelpunkt) oder die Kennung und liefert `named` als Urteil; das Ergebnis wird zu `ErrDevcontainerRuntime` und damit zu `12`. Die Form (Präfix, Fehlerausgabe, Status 1, Name des Benutzers) trägt die Zelle E5 der Messung (zwei Läufe, je zwei Wiederholungen), E6 trägt den Fall `exit 0`. „Gemessen am 2026-11-01, Linux `amd64`“ und „Lesung des Codes, kein gefahrener Lauf des Launchers“ stehen im Text, das Protokoll sagt dasselbe. Der Satz „Ist dieser Befehl der letzte der Sitzung …“ nennt einen hinreichenden Fall; der Code verlangt nur Status ungleich `0`, leere Standardausgabe und diese Zeile als erste der Fehlerausgabe, der Absatz verengt das nicht falsch. Spalte `0`, Punkt *Sitzung* („zwei Aufrufe und sonst nichts“, jetzt ohne Ausnahme, das stimmt wieder mit `Shell`/`probeShell`/`runSession`: ein Probe-Versuch je Shell, eine Sitzung) und das Änderungsprotokoll (Zeile `2026-11-01`) sind untereinander und zum Absatz konsistent; Rest siehe F-4. Das Lastenheft ist unverändert, der Widerspruch zu ihm steht im Absatz („Das Lastenheft nennt diese Ausnahme zu `0` nicht; bei Konflikt gilt es“).
3. **[ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) `Deprecated`.** Formal konform: `**Status:** Deprecated` ist ein Wert der Vorlage und des ADR-Index; `Grund`-Absatz ist wahr (Messung, Zeile auf der Fehlerausgabe mit Status 1; 28.0.4 weder bestätigt noch widerlegt; [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) gilt unverändert; „nie in Kraft getreten“). Index-Zeile in `docs/plan/adr/README.md` entspricht dem Grund. [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md) und [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) sind im Diff gegen die Basis unverändert (leere Diff-Stat). Der Wortlaut von `deprecated` im Regelwerk („Entscheidung entfällt ersatzlos, der bedarf besteht nicht mehr“) passt nur teilweise: Der Bedarf besteht als Grenze fort, die Spezifikation führt ihn; der Grund-Absatz sagt das (INFO F-5). Dass kein Slice auf die `Deprecated`-ADR verlinkt, ist in Ordnung: `make docs-check` meldet 0 Befunde, und die ADR-Kennung steht in `done/` nur als Text oder Pfad in Backticks. Der Messprotokoll-Link auf die ADR (F-1) ist der einzige verbliebene Link außerhalb der Berichte.
4. **Reste im Baum.** „Ausnahme dritter Aufruf“ und „[ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md)“ außerhalb von `done/`, Berichten, der Evidence-Datei, [ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) selbst und dem Index: nur das Messprotokoll (F-1). `spec/architecture.md`, `docs/user/`, `CHANGELOG.md` (zu `shell`) und `spec/lastenheft.md` nennen weder die Kennung noch einen dritten Aufruf (Diff gegen die Basis leer bzw. nur Fremdänderung, siehe F-7). Die Zeile in `docs/plan/planning/open/README.md` und der Roadmap-Zeiger benennen den Stand ohne Chronik-Prosa über das nötige Maß.
5. **`docs/maintainer/zusagen-und-messstaende.md`.** Aussagen belegt: Engine 28.0.4 „nicht gefahren“, Engine 29.9.2 gefahren (Messprotokoll, Präfix, Fehlerausgabe, Status 1); „die Einstufung dieser Zeile pinnt der Fall ‚the user does not exist after the probe‘“ trifft zu (`devcontainer_shell_test.go`, `msgUserGone` mit Präfix und Name `vscode`, Fehlerausgabe, Status 1); „der Fall der Sitzung, die die Zeile selbst schreibt, ist durch keinen Test gepinnt“ trifft zu. Das Wort „gepinnt“ steht nur dort, wo ein Test existiert.
6. **Offener Punkt Engine 28.0.4.** Im Slice §7 ehrlich benannt (siehe F-6): Er steht als offener Punkt 1 mit den richtigen Fundstellen (`spec/spezifikation.md` Zeilen 925 und 929 stimmen mit dem Stand überein), mit Herkunft (frühere Messung ohne Protokoll, Abschnitt 3.15), mit der Aussage „nicht bestätigt und nicht widerlegt“, dem Hinweis, dass dieser Slice die Zeile nicht ändert, und den zwei Wegen der Entscheidung. Er ist zusätzlich als Risiko 4 (Anker `BEO-DOK/zusage-weiter-als-messung`) geführt. Die Form für 28.0.4 steht im Code nur gegen den Fake (`devcontainer_shell_test.go`, Fälle „engine 28.0.4“).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapter/driven/docker/` (Diff gegen `350cdc72`) | geprüft, leer; ohne Befund |
| `devcontainer_shell.go`, `devcontainer_shell_test.go` (Lesung gegen den Spec-Absatz) | geprüft, ohne Befund |
| `spec/spezifikation.md` Absatz „Nicht zu trennen ist der Fall …“, Änderungsprotokoll, Spalte `0` | geprüft, ohne Befund (Rest F-4) |
| `docs/plan/adr/0031-…`, `0034-…` | geprüft, unverändert gegen die Basis, ohne Befund |
| `docs/plan/adr/0039-…` und Index | geprüft; F-3, sonst ohne Befund |
| `docs/maintainer/zusagen-und-messstaende.md` (Abschnitt Spezifikation von `shell`) | geprüft, ohne Befund |
| Slice in `done/` §7, §6, Plan-Tabelle | geprüft, ohne Befund (Plan-Tabelle gegen Diff abgeglichen; die Tabelle führt den Revert, die Spec-, Wartungsdoku-, ADR- und Index-Zeilen) |
| `spec/architecture.md`, `docs/user/`, `spec/lastenheft.md` | geprüft, unverändert, ohne Befund (§3.4: keine ADR-, Slice- oder Wellenbezüge) |
| Referenzmodell (Rang 1–3 verweisen nicht abwärts) | geprüft: der Spec-Absatz nennt weder Slice noch ADR; ohne Befund |
| Messskript und Messprotokoll | F-1, F-2 |
| `make test-docker`, `make gates` | nicht gelaufen, Gründe siehe Sensors |
| `internal/e2e/` und der CHANGELOG-Satz zu `--instance`, `zusagen-und-messstaende.md` Abschnitt zu Instanzen | nicht geprüft: Änderungen eines anderen Slice (F-7) |

## Weitere Beobachtungen (INFO)

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-5 | INFO | [ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) führt nach dem Grund-Absatz den Stand beim Verwerfen unverändert fort (Schärft [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), „Supersedes (teilweise)“ [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) mit „ist diese ADR die geltende Fassung“, „Nachgezogen (im Slice) … Zeile 325“). Der Grund-Absatz stellt klar, dass nichts davon in Kraft trat. Eine Aussage, die ein Leser ohne den Absatz als gültig läse, bleibt im Text; die Lesereihenfolge schützt. | `AGENTS.md` §3.5 | `docs/plan/adr/0039-…` · „Der folgende Text ist der Stand beim Verwerfen.“ | nein | Verworfene ADR trägt Stand-Text neben dem Grund |
| F-6 | INFO | Offener Punkt zur Engine 28.0.4 ist im Slice §7 benannt; die Entscheidung (Messung nachholen oder Zusage auf 29.9.2 beschränken) liegt beim Projektinhaber. Der Spec-Satz „Zugesagt sind die Engines 29.9.2 und 28.0.4“ bleibt bis dahin eine Zusage weiter als die Messung. Zuständigkeit Projektinhaber/Validator. | `AGENTS.md` §3.10 | `spec/spezifikation.md` · „Zugesagt sind die Engines 29.9.2 und 28.0.4“ | nein | Zusage weiter als Messung (Register: `BEO-DOK/zusage-weiter-als-messung`) |
| F-7 | INFO | Der Diff `350cdc72..HEAD` enthält Änderungen eines anderen Slice (`internal/e2e/devcontainer_instance_docker_test.go`, `devcontainer_start_docker_test.go`, CHANGELOG-Satz zu `--instance`, Abschnitt zu Instanzen in `zusagen-und-messstaende.md`, [slice-v1-devcontainer-instanz-test-leicht](../plan/planning/in-progress/slice-v1-devcontainer-instanz-test-leicht.md)). Sie sind hier nicht bewertet; der Slice hat einen eigenen Review-Bericht. | Maintainability | `internal/e2e/devcontainer_instance_docker_test.go` | nein | Fremdänderung im Diff-Bereich |
| F-8 | INFO | Das Benutzerhandbuch 4.9 sagt weiter „was immer der Exit-Status der Shell war“, die Spezifikation nennt jetzt die Ausnahme zu `0`. Das Handbuch folgt dem Lastenheft (Rang 1), das die Ausnahme nicht kennt; es wurde im Diff nicht geändert. Wer die Grenze für Nutzer sichtbar machen will, braucht eine Entscheidung zu Lastenheft oder Handbuch. | [LH-FA-DEV-010](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | `docs/user/benutzerhandbuch.md` · „was immer der Exit-Status der Shell war“ | nein | Handbuch bleibt hinter Spec-Grenze zurück |
| F-9 | INFO | §7 des Slice führt als offenen Punkt 2 den Nachlauf-Review („in keinem frischen Kontext gelesen“); mit diesem Bericht ist der Zustandsmarker abgelaufen. Verifier/Implementer führen ihn nach. | `AGENTS.md` §3.10 | Slice §7 · „Nachlauf-Review des Rücknahme-Diffs“ | nein | Zustandsmarker läuft mit diesem Lauf ab |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Zustandsmarker im Messprotokoll nach Ablauf nicht nachgeführt · Verweis auf entfernten Namen in Skript · Tabellenzeile mit überzähliger Zelle · Ausnahme an einer von mehreren „unbedingt 0“-Stellen nachgetragen

Transiente Befunde: keine; `make docs-check` lief mit 0 Befunden.

## Verdikt

**Merge-blockierend:** nein. Der Code ist exakt auf dem Stand vor der Slice-Arbeit, die Spec-Lesung stimmt mit dem Code und der Messzelle E5 überein, [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md) und [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) sind unverändert, der offene Punkt zur Engine 28.0.4 ist im Slice ehrlich benannt. Die vier LOW sind Nachführungen ohne Vertragsbruch.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ersetzt keine Verifikation.
