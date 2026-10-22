# Review: slice-v2-spec-devcontainer-shell-und-zweites-up (Spezifikation und Sicht nach dem Code)

* Rolle: Reviewer (frischer Kontext, nicht der Autor), Review-Art: Design-Review der Doku gegen Plan, Lastenheft, ADRs und Code
* Gegenstand: Diff `0d653453..HEAD` (10 Commits), Plan
  [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/in-progress/slice-v2-spec-devcontainer-shell-und-zweites-up.md);
  Änderungen an `spec/spezifikation.md` ([`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) und die Nachzüge in [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen), der Absatz zu [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), Änderungshistorie), an `spec/architecture.md` (§2.2 bis §2.6, ARC-008, §7.1, §7.2) und am Plan §3 (Beleg-Tabelle B1 bis B29, Streichungen)
* Geprüft gegen: Plan, Lastenheft 0.21.0
  ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten),
  [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging),
  [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`LH-SA-CLI-002`](../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)),
  die ADRs [`ADR-0030`](../plan/adr/0030-devcontainer-up-failed-nach-zustand.md),
  [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md),
  [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md) und
  [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) (als Quelle, die Spezifikation nennt sie nicht),
  den Code (`internal/adapter/driven/docker/devcontainer_shell.go`, `internal/hexagon/application/devcontainer_start.go`,
  `internal/adapter/driving/cli/devcontainer.go`, `cli.go`, `terminal.go`) und die Tests der Beleg-Tabelle,
  AGENTS.md §3.4, §3.7, §3.9, §3.10. Nicht geprüft gegen die DoD (Verifier).
* Sensoren: `make gates` Exit 0 auf `dba15a4a` (docs-check, Lint, Test, Coverage-Gate, Baseline-Verifikation, Ablage-Prüfung);
  `make doc-trace`: 84 Anforderungen, 0 Waisen. Nicht gelaufen: `make test-docker` (der Diff ändert keinen Code und keinen Test;
  die Docker-Tests der Beleg-Tabelle sind Bestand der Slices davor, ihr Lauf auf diesem Stand ist Sache des Verifiers, siehe I-1) und
  `make verify-depguard` (keine Importregel berührt).

## Findings

### HIGH

Keine.

### MEDIUM

Keine.

### LOW

**L-1 · Spezifikation sagt „kein `ioctl`“ ohne Test**
* quelle: AGENTS.md §3.10 (Bedingung in der Spezifikation, die eine Eigenschaft des Codes behauptet, braucht einen Test oder wird gestrichen)
* pfad: `spec/spezifikation.md` · „die Erkennung benutzt kein `ioctl`“ (Spiegelstrich *Terminal* im Unterabschnitt `shell` von [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes))
* klasse: code-eigenschaft-ohne-pin
* befund: Der Satz beschreibt die Implementierung von `isTerminal` (`internal/adapter/driving/cli/terminal.go`: `Stat` und Dateimodus). Kein Test der Beleg-Tabelle (B18: `TestIsTerminal`, `TestExecute_DevcontainerShell_Terminal`) kann rot werden, wenn die Erkennung später ein `ioctl` benutzt; gepinnt ist nur das Ergebnis je Dateiart. Die Streichungen führen den Satz nicht.
* risiko: Der nächste Änderer folgt einem Satz, den kein Gate stützt; wechselt die Erkennung auf `ioctl` (`TCGETS`), ist die Spezifikation falsch, ohne dass etwas rot wird. Der Zusatz trägt für den Leser nichts, was das Ergebnis nicht schon sagt.
* verifizierbar: Review-only

**L-2 · Grenze (e) des Fundament-Slice steht weder in der Spezifikation noch in den Streichungen**
* quelle: AGENTS.md §3.10; Plan §1 (Auftrag: Grenzen (a) bis (e) übernehmen oder streichen); [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („Hat die Shell begonnen, endet der Befehl mit `0`“)
* pfad: `spec/spezifikation.md` · „Alles andere gilt als Sitzung“ (Unterabschnitt `shell`, *Grenzen der Erkennung*); Plan §3 · Tabelle *Streichungen*
* klasse: benannte-grenze-nicht-uebernommen
* befund: Die Grenzen (a), (c), (d) stehen als Bedingungen im Text, (b) (`name:group`) steht in den Streichungen. Grenze (e), die Meldung zu einem fehlenden Benutzer, die eine Sitzung selbst erzeugt (ein `docker exec -u <Benutzer> <anderer Container>` auf der Engine der Form 28.0.4 als einziger Ausgabe), kommt in keiner der beiden Listen vor. Die Spalte `0` der Tabelle der Exit-Codes sagt „unabhängig vom Exit-Status der Shell“; die Liste der erkannten Meldungen (Status 126, einzige Zeile, leere Fehlerausgabe) trifft diese Sitzung und macht daraus `12` (`userMissingOnStdout`; Unit-Fall *the user does not exist after the probe (engine 28.0.4: status 126 on standard output)* in `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` erwartet dafür `12`, ohne dass der Code die Zeile der Sitzung von der des Clients trennen kann).
* risiko: Wer die Spezifikation liest, sieht die Möglichkeit eines `12` nach begonnener Sitzung nicht und hält die Zusage „nach Beginn der Sitzung immer `0`“ für ausnahmslos; die Lücke ist in der Übergabe benannt und verschwindet hier still.
* verifizierbar: Review-only

**L-3 · Plan-Aussage über den Aufrufer der vier Prüffunktionen ist falsch**
* quelle: AGENTS.md §3.10 (Aussage über Code im Plan an der Stelle gelesen)
* pfad: Plan §3 · „(Aufrufer der vier Prüffunktionen: Suche im Baum, die Application ruft die ersten zwei, `devcontainer_shell.go` alle vier)“
* klasse: plan-aussage-ueber-code
* befund: `devcontainer_shell.go` ruft selbst nur `ValidateDevcontainerShellUser` und `ValidateDevcontainerWorkspace`; `ValidateDevcontainerEnvName` steht in `devcontainer_launcher.go` (über `passedEnvNames`, auch für `Shell`) und `ValidateDevcontainerProjectFolder` nur im `Up`-Pfad desselben Pakets. Die Sicht trägt die richtige, gröbere Aussage („der Adapter ruft alle vier“); nur die Fundstelle im Plan ist falsch.
* risiko: Kleine Fundstellen-Drift im Plan, den der Verifier als Beleg liest. Die Spezifikation und die Sicht sind nicht betroffen.
* verifizierbar: Review-only

**L-4 · [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) begründet mit [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand* in anderen Worten, als [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) sie sagt**
* quelle: Maintainability (Konsistenz innerhalb der Spezifikation); AGENTS.md §3.10
* pfad: `spec/spezifikation.md` · „denn sie beschreibt den Zustand des Containers“ (Absatz „`up` urteilt nach dem Zustand des Containers“) gegen [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) · „die Datei beschreibt den letzten Lauf ihres Scripts, nicht notwendig den letzten Start des Containers“
* klasse: verweis-traegt-aussage-nicht
* befund: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) verweist für „Datei aus einem früheren Start“ auf *Stand* und sagt, sie beschreibe „den Zustand des Containers“; *Stand* sagt das Gegenteil an Schärfe: Sie beschreibt den letzten Lauf *ihres Scripts*. Gemeint ist dasselbe (ein `failed` bleibt stehen, bis das Script neu läuft), die Worte tragen es nicht.
* risiko: Ein Änderer von [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Stand* sieht in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) keine Kopplung und kann die Aussage zum früheren Start brechen.
* verifizierbar: Review-only

### INFO

**I-1 · Docker-Tests der Beleg-Tabelle liefen nicht in diesem Slice (Verifier)**
Die Zeilen B5, B8, B13 bis B15, B19, B20, B23, B25 und B29 stützen sich auf Tests mit Build-Tag `docker` (`internal/e2e`, `internal/adapter/driven/docker`). Ich habe ihren Wortlaut gegen die Spezifikation gelesen, nicht gefahren. Ob sie auf `dba15a4a` im CI grün sind (`make test-docker`), prüft der Verifier.

**I-2 · Engine-Versionen im Text sind Messbedingung, kein Pin**
`29.9.2` und `28.0.4` stehen in der Spezifikation (Absatz zu den Lebenszyklus-Befehlen, *Messbedingungen*, *Engine-Abhängigkeit*). Kein Test fragt die Version ab; die Docker-Tests laufen gegen die Engine des Runners. Gedeckt ist die Aussage durch die Formulierung „gemessen“ und „ein Lauf je Fall“ und durch die Test-Kommentare (*Observed with engine and client 29.9.2*); die Form der Engine 28.0.4 pinnt nur der Unit-Test mit dem Fake (`\r\n`, Status 126). Das entspricht [`ADR-0033`](../plan/adr/0033-devcontainer-up-laufender-container-lebenszyklus-befehle.md) und [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md); „die gepinnte Dev-Container-CLI“ ist ohne Zahl geschrieben, wie der Plan es festlegt.

**I-3 · Bedingung „Standardausgabe leer“ gilt im Code für alle Meldungen der Fehlerausgabe**
`clientMessageOnStderr` verwirft jede Meldung der Fehlerausgabe, wenn ohne Terminal die Standardausgabe nicht leer ist; die Spezifikation nennt die Bedingung nur beim ersten Spiegelstrich (`Error response from daemon: `). Für „Cannot connect…“ und „stdin is not a terminal“ deckt „eine Meldung des Clients nach anderer Ausgabe der Sitzung“ im Absatz *Alles andere* die Aussage ab; ein eigener Test für diese zwei Meldungen mit nicht leerer Standardausgabe fehlt, die Spezifikation behauptet ihn nicht. Keine Änderung nötig.

**I-4 · `TestIsTerminal`: der positive Fall hängt an `/dev/ptmx`**
Der Fall „Zeichengerät, kein Nullgerät, ergibt Terminal“ wird nur angelegt, wenn `/dev/ptmx` sich öffnen lässt; sonst bleibt er stumm aus. Auf dem Test-Image des Gates ist das erfüllt (`make gates` grün), die Zusage „ein Zeichengerät, das nicht das Nullgerät ist“ ist damit bedingt gepinnt.

## Negativbefunde

| Bereich / Linse | Ergebnis |
|---|---|
| Strata-Regeln (AGENTS.md §3.4): Suche im Diff von `spec/` nach ADR-, Slice-, Welle-, Carveout-, Roadmap-Kennungen, Release-Versionen (`v0.`), Phase-Tags, Verifier-/Mutations-/Messprotokoll-Begriffen, Prozesszuständen, Zeilenzahlen | geprüft, ohne Befund (einzige Zahlen: Engine-Versionen als Messbedingung, `0.21.0` in der Änderungshistorie wie bei den Zeilen davor) |
| Sicht sprach- und meilensteinfrei, keine ADR-Bezüge; Code-Pfade und Symbole nur als Module | geprüft, ohne Befund |
| Beleg-Tabelle: alle 60 genannten Testnamen mit `grep -rn "^func <Name>(" internal` gesucht | geprüft, ohne Befund (jeder Name genau eine Definition) |
| Beleg-Tabelle auf Test-Ebene gelesen (17 von 29 Zeilen: B1, B3, B4, B5, B7 bis B9, B14 bis B21, B23, B25; die übrigen nur auf Namensebene) und gegen Spezifikation gehalten: Engine-Texte, Exit-Codes, Terminal, `alreadyRunning`, Zähler der Lebenszyklus-Befehle, `sessionEnvNames`, OCI-Zeile Status 127 (Status 1 ist Sitzung), `u-boot.yaml`-Grenze, `vfs` vor `egress`, kein Hilfs-Container während der Sitzung | geprüft, ohne Befund außer L-1 |
| Lastenheft 0.21.0: Spalte `shell` der Exit-Code-Tabelle, `alreadyRunning`, `--json` bei `shell` (Exit `2`, Umschlag mit `subcommand: "shell"`, Test `jsontestutil` prüft Command, Subcommand, Exit-Code, Code [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)), Socket „für die Dauer des Starts“, `--quiet`-Sätze zu `shell`, Befehlsübersicht | geprüft, ohne Befund (die Spezifikation präzisiert, erweitert nicht) |
| Konsistenz: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) gegen [`SPEC-043`](../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), [`SPEC-024`](../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) und die Sicht; „zwei Prüffunktionen“ ist durch „vier“ ersetzt; Suche nach „beide Unterbefehle“, „Up/Down“, „nach einem gelungenen Start“ ohne Rest | geprüft, ohne Befund außer L-4 |
| Sicht §7.1/§7.2 und §2.3 (Code 11) gegen `ExitCode` (`isEnvironmentError` vor `isRuntimeError`), `mapDevcontainerEngineSentinel`, `mapDevcontainerUpErrorToDiagnostic`, `mapDevcontainerShellErrorToDiagnostic` (nur Klasse `2`, weil `--json` vor dem Use-Case abgewiesen wird); Dual-Classifier-Regel gewahrt, Pin `TestExecute_Devcontainer_SentinelsPerSubcommand` und `TestExitCode_DevcontainerShellUsageSentinels` | geprüft, ohne Befund |
| Streichungen: keine gestrichene Bedingung wird von einem Test getragen (OCI-Zeile `126`, Terminal-Fälle „Container entfernt“ und „Engine nicht erreichbar“, `name:group`, Ende der OCI-Zeile, Reihenfolge `daemon`, Mechanismus der untätigen Shell); gehaltene Bedingungen mit Test (Wächterfälle in `TestLauncher_Shell_SessionEndIsNoError`) | geprüft, ohne Befund außer L-2 |
| Änderungshistorie: Datum 2026-10-22, Lastenheft 0.21.0 (Spezifikation führt keine eigene Versionsnummer), Aufzählung der berührten Stellen und der nicht zugesagten Punkte | geprüft, ohne Befund |
| d-check/Anker: `make docs-check` im `make gates` grün, `make doc-trace` ohne Waise | geprüft, ohne Befund |
| Plan-Tabelle §3 gegen `git diff --stat 0d653453..HEAD`: nur Spezifikation, Sicht, Planung (Move-Nachzug in vier `done/`-Slices, Welle, `open/README.md`, Roadmap) und elf Reports; kein `internal/`, `cmd/`, Lastenheft, ADR, Handbuch | geprüft, ohne Befund |
| Kommentar-Muster (AGENTS.md §3.9) im Diff: keine Beschreibung von Entferntem; „aus einem früheren Start“ beschreibt den Zustand | geprüft, ohne Befund |
| `make test-docker`, Docker-Tests, Mutationen | nicht geprüft: Doku-Slice ohne Code; Lauf ist Sache des Verifiers (I-1) |
| Handbuch `docs/user/`, CHANGELOG, Lastenheft | nicht geprüft: nicht im Diff, kein Gegenstand des Slice |

## Summary

0 HIGH, 0 MEDIUM, 4 LOW, 4 INFO. Die Spezifikation und die Sicht stimmen an den gelesenen Stellen mit dem Code und mit den Tests überein; die Aussagen der Beleg-Tabelle tragen. Die LOW-Findings sind Aufräumpunkte (ein ungepinnter Satz, eine nicht übernommene Grenze, eine falsche Plan-Fundstelle, eine Formulierungsdrift).

**GO** für den Verifier.
