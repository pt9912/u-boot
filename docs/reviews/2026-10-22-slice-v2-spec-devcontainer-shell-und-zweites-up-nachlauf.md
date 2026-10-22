# Review: slice-v2-spec-devcontainer-shell-und-zweites-up (Nachlauf zur Fix-Runde)

* Rolle: Reviewer (frischer Kontext, nicht Autor), Review-Art: Nachlauf-Review der Fix-Runde gegen den ersten Report
  [`2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up`](2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up.md) (4 LOW, GO)
* Gegenstand: Diff `3635dcc2..HEAD` (Fix-Commit `e600abc9` und zwei Nachträge): `spec/spezifikation.md` (vier Stellen) und Plan
  [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/in-progress/slice-v2-spec-devcontainer-shell-und-zweites-up.md) §3 (Sicht-Zeile, Streichungstabelle, Nachtrag *Fix-Runde nach dem Review*)
* Geprüft gegen: Plan, [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten),
  [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand*, Grenze (e) im Fundament-Slice
  [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md), Code
  `internal/adapter/driven/docker/devcontainer_shell.go` (`clientErrorBeforeSession`, `clientMessageOnStderr`, `userMissingOnStdout`), die Tests
  `TestLauncher_Shell_SessionEndIsNoError`, `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure`, `TestIsTerminal`, AGENTS.md §3.4, §3.7, §3.9, §3.10. Nicht gegen die DoD (Verifier).
* Sensoren: `make gates` Exit 0 auf dem Stand nach `e38456b9` (Lint `0 issues`, Coverage 92,00 %, `d-check: 713 Datei(en) geprüft, 0 Befund(e)`, Baseline-Verifikation, Ablage-Prüfung).
  Nicht gelaufen: `make test-docker` (kein Code und kein Test im Diff; Sache des Verifiers, siehe I-1 des ersten Reports) und `make verify-depguard` (keine Importregel berührt).

## Behebung der Findings des ersten Reports

| Finding | Ergebnis | Beleg |
|---|---|---|
| L-1 (`ioctl` ohne Test) | behoben | Satz in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) gestrichen; `grep -in ioctl spec/` ohne Treffer; im Plan nur in der Streichungszeile, der Fix-Tabelle und der Zeile des Auftrags (Zeile 50, Auftrag), nirgends als Zusage |
| L-2 (Grenze (e)) | behoben | siehe unten |
| L-3 (Aufrufer der Prüffunktionen) | behoben | Plan-Zeile nennt `devcontainer_launcher.go` mit `ValidateDevcontainerProjectFolder` (nur `Up`, Zeile 131) und `ValidateDevcontainerEnvName` (über `passedEnvNames`, Zeile 271), `devcontainer_shell.go` mit `ValidateDevcontainerShellUser` und `ValidateDevcontainerWorkspace` (Zeilen 56, 59), Application mit den ersten zwei (Zeilen 255, 372): stimmt mit `grep` über `internal` überein |
| L-4 (Wortlaut zu *Stand*) | behoben | Satz lautet jetzt „die Datei beschreibt den letzten Lauf ihres Scripts, nicht notwendig den letzten Start des Containers“, wörtlich wie [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand* |
| I-3 (Bedingung „Standardausgabe leer“) | als Grenze vertretbar | Der Code (`clientMessageOnStderr`) verlangt sie ohne Terminal für jede Meldung der Fehlerausgabe; die Spezifikation behauptet sie nur beim ersten Spiegelstrich und deckt die übrigen über „eine Meldung des Clients nach anderer Ausgabe der Sitzung gilt als Sitzung“. Die Aussage ist zu schwach, nicht zu stark; keine Änderung nötig (INFO, kein Finding) |

### Zu L-2 im Einzelnen

* Grenze (e) im Fundament-Slice (§3 *Ist gegen Plan* Punkt 3 und Nachlauf-F-2): ein `docker exec -u <Benutzer der Sitzung> <anderer Container>` in der Sitzung gegen eine Engine der Form 28.0.4 als einzige Ausgabe, fünf Bedingungen zugleich, aus der Matrix abgeleitet, nicht gefahren.
* Code `userMissingOnStdout`: `!opts.Terminal && exit == 126 && stderr.total == 0 && stdout.singleLine() && HasPrefix(line, missingUserLead(User))`. Das sind genau fünf Bedingungen; die Liste der erkannten Meldungen in der Spezifikation nennt sie (ohne Terminal, Engine 28.0.4, einzige Zeile der Standardausgabe, leere Fehlerausgabe, Status `126`, Doppelpunkt als Namensgrenze). Der neue Satz am Ende dieses Spiegelstrichs trifft diese Klasse.
* „der Befehl endete dann mit `12` statt `0`“: `Shell` gibt bei `clientErrorBeforeSession` `ErrDevcontainerRuntime` zurück, das ist Klasse `12`. Der positive Fall *the user does not exist after the probe (engine 28.0.4: status 126 on standard output)* in `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` pinnt die Erkennung der Form, die Wächterfälle in `TestLauncher_Shell_SessionEndIsNoError` (*status 1, the missing-user line as the only output*, *status 126, the missing-user line and text on standard error*, *other output before/after*, *terminal, status 126*, *vscode2*) pinnen die Ränder. Der Erzeuger in der Sitzung ist als „aus der Form abgeleitet, nicht gefahren“ gekennzeichnet und deckt sich mit dem Fundament-Slice.
* Spalte `0` der Tabelle der Exit-Codes: „(Ausnahme: die Grenzen der Erkennung im Unterabschnitt *`shell`*)“; der Spiegelstrich *Grenzen der Erkennung* existiert unter diesem Namen.
* Verhältnis zum Lastenheft: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt „Exit-Code `0` nach Beginn der Shell“ ohne Ausnahme. Die Ausnahme der Spezifikation engt diese Zusage für eine abgeleitete Randklasse ein, erweitert das Lastenheft nicht. Dieselbe Klasse (Fehler des Clients wird von einer begonnenen Sitzung nicht getrennt) trugen die Grenzen (c) und (d) schon im ersten Durchlauf; siehe I-1.

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

Keine.

### INFO

**I-1 · Ausnahme der Spalte `0` verengt die Zusage des Lastenhefts, ohne dass das Lastenheft sie nennt**
Die Spalte `0` und Grenze (e) sagen, dass nach begonnener Sitzung in einer Randklasse `12` entsteht; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt keine Ausnahme. Der Fundament-Slice benennt die Abweichung (Entscheidung zu Grenzen (d), (e)) und sie liegt im Bereich der Präzisierung durch den Rang 2; ob das Lastenheft einen Hinweis bekommt, ist eine Entscheidung des Architekten, kein Mangel dieses Diffs.

**I-2 · Grenze (e) steht im Plan in der Tabelle *Streichungen***
Die Zeile ist keine Streichung, sondern eine übernommene Grenze; sie steht dort, weil die Tabelle die Auftragsstücke aus dem Fundament-Slice bucht. Die Spalte *Wortlaut/Ausgang* macht es kenntlich. Keine Änderung nötig.

**I-3 · „Standardausgabe leer“ ohne Test außerhalb der `Error response from daemon`-Meldung**
Für `Cannot connect to the Docker daemon at` bei nicht leerer Standardausgabe fehlt ein Unit-Fall; die Spezifikation behauptet nur den Rand über den Absatz *Alles andere*. Unverändert zum ersten Report (I-3 dort), keine Änderung nötig.

**I-4 · Docker-Tests liefen nicht**
Wie im ersten Report (I-1): Sache des Verifiers.

## Negativbefunde

| Bereich / Linse | Ergebnis |
|---|---|
| §3.10: neue Aussagen über Verhalten im Diff (Grenze (e), Spalte `0`, *Stand*-Wortlaut, `ioctl`-Streichung) gegen Code und Tests gelesen | geprüft, ohne Befund |
| Test-Namen der neuen/geänderten Zeilen: `grep -rn "^func <Name>("` für `TestLauncher_Shell_SessionEndIsNoError` (`devcontainer_shell_test.go:229`), `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` (`:304`), `TestIsTerminal` (`cli/devcontainer_shell_test.go:333`): je genau eine Definition; Fälle gelesen, sie pinnen die Aussage | geprüft, ohne Befund |
| §3.4: Diff von `spec/` auf ADR-, Slice-, Welle-, Carveout-, Release-Verweise und Phase-Tags | geprüft, ohne Befund (null Treffer in den hinzugefügten Zeilen) |
| §3.7/§3.9: Kommentar- und Chronik-Muster im Diff („früher“, „no longer“, „ohne X wäre“): der neue Satz zu Grenze (e) beschreibt den Ist-Zustand im Konjunktiv nur als abgeleiteten Fall und ist als solcher gekennzeichnet | geprüft, ohne Befund |
| Bare Kennungen im Nachtrag des Plans: [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) ist verlinkt; `make docs-check` (Teil von `make gates`) 0 Befunde | geprüft, ohne Befund |
| Nachtrag *Fix-Runde nach dem Review* gegen den Diff: fünf Zeilen (L-1 bis L-4, I-3) stimmen mit den Änderungen in Spezifikation und Plan überein; die Aussage zur Sicht-Tabelle („die Sicht stimmte“) deckt sich mit `spec/architecture.md` (im Diff unverändert) | geprüft, ohne Befund |
| Plan-Tabelle gegen `git diff --stat 0d653453..HEAD`: berührt sind Spezifikation, Sicht, Planungsdateien (Move-Nachzug, Welle, Roadmap, `open/README.md`) und Reports; kein `internal/`, `cmd/`, Lastenheft, ADR, `docs/user/`, README, CHANGELOG | geprüft, ohne Befund |
| `make test-docker`, `make verify-depguard` | nicht geprüft: Doku-Diff ohne Code; Lauf ist Sache des Verifiers |
| Handbuch, CHANGELOG, Lastenheft | nicht geprüft: nicht im Diff |

## Summary

0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO. L-1 bis L-4 sind behoben, I-3 ist als Grenze vertretbar. Die Fix-Runde hat keine neue Aussage über Verhalten ohne Test und keinen Verweis nach unten eingeführt.

**GO** für den Verifier.
