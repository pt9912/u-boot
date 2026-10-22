# Review: slice-v1-cli-devcontainer-shell (`u-boot devcontainer shell`)

* Rolle: Reviewer (frischer Kontext, nicht der Autor)
* Gegenstand: Diff `41e0b262..HEAD` (10 Commits), Plan
  [`slice-v1-cli-devcontainer-shell`](../plan/planning/done/slice-v1-cli-devcontainer-shell.md)
* Geprüft gegen: Plan, [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging),
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), AGENTS.md §3.1 bis §3.10, §5 Regel 1 und 6.
  Nicht geprüft gegen die DoD (Verifier).
* Sensoren: `make gates` Exit 0 (lint, test, coverage-gate, docs-check, baseline-verify,
  ablage-pruefung). `make test-docker` nicht gefahren: Der Reviewer stützt sich auf das Lesen
  der Docker-Tests; der Bericht des Implementers über den lokalen Lauf ist nicht
  nachgeprüft. `make verify-depguard` nicht gefahren: kein neuer Import über eine
  Schichtgrenze (Application importiert nur `io`, `strings`).

## Findings

### MEDIUM

**M-1 · Handbuch `--quiet`-Zeile: Ausnahmen (1) und (2) nennen `devcontainer shell` nicht**
* quelle: Doku driftet gegen implementiertes Verhalten (Reviewer-Skill MEDIUM); AGENTS.md §3.10
* pfad: `docs/user/benutzerhandbuch.md` · „erzeugen nur `generate devcontainer` und `devcontainer up` (die Einbindung des Docker-Sockets in den Hilfs-Container)“ und „in `doctor` und in der Warnungsausgabe von `generate devcontainer`, `init` und `devcontainer up`“
* klasse: doku-aufzaehlung-nicht-nachgezogen
* befund: Dieselbe Tabellenzelle führt `devcontainer shell` jetzt in der Liste der Befehle mit
  entfallender Erfolgsmeldung und beim zurückgehaltenen Fortschritt, nennt es aber in den
  zwei Aufzählungen nicht, welche Befehle unter `--quiet` die Socket-Ausweisung und die
  Degradations-Fallbacks sichtbar halten. `shell` tut beides
  (`TestExecute_DevcontainerShell_Quiet`, `TestDevcontainerShell_FallbackWarningsComeBeforeTheSession`,
  `visibleWarnings` in `runDevcontainerShell`); das Wort „nur“ macht die erste Aufzählung falsch.
* risiko: Ein Nutzer, der `--quiet` setzt, liest, dass `shell` die Sicherheitslockerung nicht
  ausgibt; die Ausweisung ist aber ein Produktvertrag ([`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2). Die Vorab-Suche
  des Plans nennt die Zelle als „geändert“, die zwei Aufzählungen sind übersehen.
* verifizierbar: Review-only (`make docs-check` fängt das nicht)

### LOW

**L-1 · Hilfetext: Satz zu `failed` ist unbedingt, die Messung und das Handbuch sind bedingt**
* quelle: AGENTS.md §3.10 (Aussage über Verhalten, Bedingung)
* pfad: `internal/adapter/driving/cli/devcontainer.go` · „ends the command with exit 11 and opens no shell, also when the container was already running“
* klasse: aussage-ohne-bedingung
* befund: Der Docker-Wächter (`TestDevcontainerShell_Docker_FailedStatusOpensNoSession`) zeigt:
  ohne `u-boot.yaml` öffnet `shell` im selben Container die Shell. Das Handbuch sagt die
  Bedingung, die Zeile `11` der Exit-Liste im Hilfetext auch, der erklärende Absatz nicht;
  `TestExecute_DevcontainerShell_Help` pinnt den unbedingten Wortlaut.
* risiko: Hilfe und Handbuch widersprechen sich in einem Randfall; der Test zementiert die
  kürzere Fassung.
* verifizierbar: `make test` (Help-Test), Docker-Wächter

**L-2 · Handbuch nennt `remoteEnv` als Quelle, der Code liest die Form überall**
* quelle: AGENTS.md §3.10 (Bedingung in Handbuch ohne pinnenden Test)
* pfad: `docs/user/benutzerhandbuch.md` · „Die Variablen, auf die `remoteEnv` des erzeugten `devcontainer.json` mit“ gegen `internal/hexagon/application/devcontainer_start.go` · „finds the form wherever it stands, not only below `remoteEnv`“
* klasse: handbuch-bedingung-ohne-test
* befund: `sessionEnvNames` gibt den Eintrag `"X": "${localEnv:X}"` auch aus `containerEnv`
  oder jedem anderen Objekt weiter; `TestSessionEnvNames` hat keinen Fall, der den Ort pinnt
  (der Fall „sorted and distinct“ verschachtelt nur zufällig). Handbuch und Code sagen
  Verschiedenes.
* risiko: Wer das Handbuch liest, erwartet eine Einschränkung auf `remoteEnv`; ein späterer
  Änderer, der sie einführt, bricht keinen Test.
* verifizierbar: `make test` (fehlender Fall)

**L-3 · Port-Kommentar `Terminal` sagt mehr zu, als der Adapter entscheidet**
* quelle: AGENTS.md §3.7 (Zusage), Maintainability
* pfad: `internal/hexagon/port/driving/devcontainer.go` · „The adapter sets it when its standard input is a terminal; Stdin and Stdout are then the files of that terminal.“
* klasse: kommentar-zusage-ohne-garantie
* befund: `isTerminal` prüft nur die Standardeingabe. Bei `u-boot devcontainer shell | tee log`
  ist `Terminal` wahr, `Stdout` aber eine Pipe. Der Launcher reicht dann `-t` plus Pipe an den
  Client; das ist nicht gemessen (die Docker-Tests decken Terminal auf Ein- und Ausgabe sowie
  Pipe auf beiden ab, nicht gemischt).
* risiko: Der nächste Änderer verlässt sich auf die Zusage „Stdout ist das Terminal“.
* verifizierbar: Review-only

**L-4 · Testkommentar zählt falsch**
* quelle: AGENTS.md §3.7
* pfad: `internal/adapter/driving/cli/devcontainer_shell_test.go` · „// Both sentinels of exit class 2 are in `ExitCode`.“
* klasse: kommentar-stimmt-nicht-mit-code
* befund: Die Tabelle prüft drei Sentinels und eine Umhüllung.
* risiko: gering; irreführend beim Erweitern.
* verifizierbar: Review-only

### INFO

**I-1 · Terminal-Erkennung ohne `ioctl` (Schwerpunkt 3)** · `internal/adapter/driving/cli/terminal.go`.
Andere Zeichengeräte (`/dev/zero`, `/dev/urandom`) zählen als Terminal; `docker exec -t` weist
sie ab, der Launcher macht Exit 12. Das Risiko ist gering (ungewöhnliche Eingabeumleitung),
die Grenze steht im Doc-Kommentar und in Plan Entscheidung 6, `TestIsTerminal` pinnt Nullgerät,
Pipe, Datei, geschlossene Datei, nil. Der positive Fall hängt an `/dev/ptmx` und wird sonst
vom Docker-Test unter `script` gedeckt. Kein Handbuch-Satz behauptet mehr. Entscheidung des
Projektinhabers bleibt (Plan §3).

**I-2 · `sessionEnvNames` ignoriert Kommentare, `localEnvNames` nicht (Schwerpunkt 4)**.
Konsistent mit Least-Privilege und mit der Fitness Function von [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md); die Teilmengen-
Beziehung gilt (jede Session-Variable steht auch in `localEnvNames`). Der Test
`…OpensTheSessionWithTheResultOfTheStart` (Mutation M5) pinnt den Unterschied. Beobachtung:
Der Hilfs-Container, der das Engine-Socket hält, bekommt weiter auch die Namen aus
auskommentierten Verweisen (und ein auskommentierter ungültiger Name ergibt Exit 10). Das ist
Verhalten von `up` (nicht Teil des Diffs) und im Handbuch nicht zugesagt.

**I-3 · Wächter-Qualität (Schwerpunkt 6)**. Die Prüfung `Config.Env` des Devcontainers hat keine
eigene Mutation: D6 scheitert laut Plan vor der Prüfung, und der Wächter-Container belegt nur,
dass `envHolds` einen per `docker run -e` gesetzten Wert findet. Der Plan sagt das ehrlich. Kein
Nutzersatz hängt neu daran (die bestehende Aussage zu `Config.Env` im Handbuch ist nicht Teil
des Diffs). Die übrigen Wächter sind stark: Helfer-Beobachtung mit Gegenprobe für den Präfix,
Gegenprobe ohne `-e`, Pipe-Gegenprobe zu `script`, Gegenprobe ohne `u-boot.yaml`. Die Mutationstabelle
ist plausibel, aber nicht nachgefahren; die Zeile „drei neue Docker-Tests“ im Plan zählt Teiltests
(zwei Top-Level-Funktionen).

**I-4 · Verifier-Hinweis**. Docker-Läufe (`make test-docker`, `integration-docker` auf dem
Push-Tip) sind Sache des Verifiers.

## Negativbefund (geprüft, ohne Befund)

* Dual-Classifier: `ErrDevcontainerShellJSONNotApplicable` steht in `isDevcontainerUsageError`
  (damit in `isUsageError`, `ExitCode` Klasse 2) und in der Diagnostic-Abbildung von `shell`
  (`mapDevcontainerShellErrorToDiagnostic`, nur Klasse 2; Begründung im Kommentar stimmt, weil
  `--json` vor dem Use-Case abgewiesen wird). Pins: `TestExitCode_DevcontainerShellUsageSentinels`,
  `TestExecute_DevcontainerShell_SentinelsAndExitCodes` (alle Klassen 10/11/12/14/1, Vorrang 11 vor 12),
  `TestExecute_DevcontainerShell_RejectsFlagsAndArguments`, Wiring-Test in `cmd/uboot/main_test.go`.
* Schichten/Importregeln: Application und Port ohne neue externe I/O; Terminal-Erkennung im
  Driving-Adapter; Wiring unverändert außer Kommentaren.
* Sprachvertrag: neue Fehlertexte, Hilfe und Warnungen Englisch.
* Standardströme: Ausweisung, Fortschritt, Warnungen, Fehler auf stderr; stdout nur Sitzung;
  `--quiet` hält Fortschritt zurück und schreibt ihn bei Fehler (C5), Ausweisung bleibt sichtbar.
* Kennungen: keine `LH-`/`ADR-`/`SPEC-`/`ARC-`/`slice-`/`welle-` in den hinzugefügten Zeilen von
  `docs/user/`, `README.md`, `README.de.md`; Hilfe ohne Kennung (Test).
* Handbuch/Hilfe/CHANGELOG ohne Aussagen zu Lebenszyklus-Befehlen, Login-Shell, Dauer oder
  Abbruchcode (die `kill -9`-Aussage folgt [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md); Test im Fundament).
* §3.9 (a) bis (c): keine Chronik-Kommentare im Diff (Wortsuche und Lesen), kein Helfer zwischen
  Doc-Kommentar und Funktion; umbenannte Symbole (`socketRelaxationText`, `referencedEnvNames`)
  ohne Altnamen in Kommentaren.
* Plan-Tabellen-Abgleich gegen `git diff --stat 41e0b262..HEAD`: Jede Go- und Doku-Datei hat
  Plan-Zeile oder Eintrag unter „Ist gegen Plan“; Dateien in `done/`, `observations/`, `open/`
  stammen vom Verweis-Nachzug. Keine Auslassung gefunden.
* `terminal.go`, `WaitDelay`/Standardströme: Der Launcher gehört dem Fundament (nicht im Diff);
  der Befehl reicht `os.Stdin`/`os.Stdout`/`os.Stderr` (verschiedene Dateien) unverändert durch.
* Gates: keine Lockerung, kein `//nolint`, `.golangci.yml` unverändert.

## Nicht geprüft

* `make test-docker`, `make verify-depguard`, Engine 28.0.4 des CI-Runners: nicht gefahren
  (siehe oben); die Mutationstabelle ist nicht reproduziert.
* DoD-Erfüllung, Closure-Notiz, Evidence: Verifier.

## Empfehlung

**GO** mit Auflage: M-1 vor dem Merge klären (Handbuch-Zelle `--quiet` nachziehen); L-1 bis L-4
nach Ermessen des Implementers im selben Zug. Kein HIGH.
