# Review-Report (Nachlauf, klein, Code): CI-Fix von [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) — 2026-10-22

**Review-Art:** Nachlauf-Review gegen Plan, Entscheidungen und Hard Rules (Wartbarkeit), nur der CI-Fix. Die DoD-Frage gehört dem Verifier. Nichts vom Implementer ist ungeprüft übernommen: die Erkennung, ihre Bindungen, die Wächterfälle, die Port-Doc und die Engine-Matrix sind am Code, an den Tests und am Slice-Text gelesen; vier Mutationen und eine Sonde sind selbst gelaufen. Die Docker-Tests und die Messungen gegen Engine 28.0.4 sind gelesen, nicht wiederholt.

**Gegenstand:** `git diff 87420c45..8fa83a34 -- internal` (`devcontainer_shell.go`, `devcontainer_shell_test.go`, Port `devcontainer_launcher.go`; 57 Einfügungen, 9 Löschungen, wie der Slice sie zählt) und `f25b456b` (nur Slice §3 *CI-Fix*, §7). `HEAD` = `f25b456b`; der Arbeitsbaum war zu Beginn und ist am Ende außer diesem Report unverändert; nichts committet, nichts gepusht, nichts auf GitHub geschrieben.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `f25b456b` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-22

**Eingangs-Kontext:** Slice §3 *CI-Fix* (Befund, Engine-Matrix, Änderung, M36 bis M41) und §7; [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 3; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Spiegelstrich zu `shell`), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); die Reports [`2026-10-22-slice-v1-devcontainer-shell-fundament`](2026-10-22-slice-v1-devcontainer-shell-fundament.md) und [`…-verify`](2026-10-22-slice-v1-devcontainer-shell-fundament-verify.md); Register-Eintrag [`BEO-ADP/erkennung-am-meldungstext-engine-version`](../plan/planning/observations/BEO-ADP/erkennung-am-meldungstext-engine-version/observation.md); `AGENTS.md` §3.7, §3.9, §3.10.

## Sensor-Läufe (alle über `make`, Docker-only)

| Sensor | Ergebnis |
|---|---|
| `make test` (Arbeitsbaum) | Exit 0 |
| `make lint` (enthält `gofmt`) | Exit 0, `0 issues.` |
| `make docs-check` | Exit 0; 690 Dateien ohne diesen Report, 691 Dateien mit ihm, jeweils 0 Befunde |
| `make doc-commits RANGE="87420c45..f25b456b"`, `make doc-immutable RANGE="87420c45..f25b456b"` | Exit 0, 0 Befunde |
| Commit-Messages von `87420c45..f25b456b` (vier Commits) | `grep -E 'SPEC-\|ARC-'` leer; jede Message nennt den Slice, [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) |
| Mutationen und Sonde | in einer Kopie von `HEAD` im Scratchpad (`git archive`, `git apply`, `make test IMAGE=<eigener Tag>`); danach `make clean IMAGE=<eigener Tag>`; `git status` im Repo vor und nach den Läufen leer |

| Mutation (`devcontainer_shell.go`) | Ergebnis |
|---|---|
| M37: `userMissingOnStdout` ohne `exit == 126` | Exit 2, rot: `…SessionEndIsNoError/status_1,_the_missing-user_line_as_the_only_output` (wie im Slice) |
| M40: ohne den Benutzer (`HasPrefix(line, clientUserMissingText)`) | Exit 2, rot: `…/status_126,_the_missing-user_line_of_another_user` (wie im Slice) |
| M39: ohne `stdout.singleLine()` | Exit 2, rot: nur `…/status_126,_other_output_after_the_missing-user_line` (wie im Slice) |
| **M-A (neu): `singleLine` ohne `if b.total > len(b.data) { return false }`** | **Exit 0, grün, die Mutation überlebt** (F-3) |
| **Sonde (neu): Wächterfall „Zeile eines Benutzers, dessen Name mit dem der Sitzung beginnt“** (`unable to find user vscode2: …`, Status 126, Sitzung als `vscode`, erwartet: kein Fehler) | **Exit 2: die Erkennung wertet die Zeile als Fehler des Clients** (F-1) |

## Prüfpunkte (falsche Positive zuerst)

**1. Korrektheit der Erkennung.** Die fünf Bindungen stehen im Code, wie der Slice sie nennt (`!opts.Terminal`, `exit == 126`, `stderr.total == 0`, `stdout.singleLine()`, `HasPrefix` auf Text und Benutzer). Die vorhandene Erkennung ist unverändert: `clientUserMissingText` ersetzt nur das Wortstück im Vergleich der Fehlerausgabe (derselbe Text), die Kette in `clientErrorBeforeSession` ist Fehlerausgabe, dann OCI, dann neu.

- **Falsche Positive.** (i) `su nobody`, `sudo -u`, `runuser` und die Shells selbst schreiben „unable to find user“ nicht; den Text schreibt ein `docker`-Client. Der einzige natürliche Erzeuger in einer Sitzung ist ein `docker exec -u <Benutzer der Sitzung> <anderer Container> …` im Inneren der Sitzung, gegen eine Engine in der Form 28.0.4 (abgeleitet aus der Matrix, nach der die Form an der Engine hängt; nicht gefahren). Seine Zeile müsste die einzige Ausgabe der Sitzung sein, ohne Terminal, als letzter Befehl (der Status 126 eines nicht interaktiven `bash` ist der des letzten Befehls), Fehlerausgabe leer: fünf zugleich geltende Bedingungen, und ein Docker-Client mit Zugriff auf eine Engine im Inneren der Sitzung. Eng, aber nicht leer; derselbe Zielkonflikt wie die benannte Grenze (d). (ii) Die Bindung an den Benutzer ist ein Präfix-Vergleich ohne Namensgrenze und weitet das Fenster auf Namen mit gleichem Anfang (F-1). (iii) Die Klasse steht nicht in den benannten Grenzen (F-2).
- **Falsche Negative.** CRLF und LF: beides bedient (`TrimSpace` je Zeile; Positivfall mit `\r\n`, Wächter mit `\n`). `xubluser` für `ubluser`: nicht erkannt, richtig (das Präfix liegt am Zeilenanfang); der umgekehrte Fall ist F-1. Unicode im Namen: bytegleicher Vergleich ohne Normalisierung; mit Unicode-Namen nicht gemessen, kein Befund. Andere Engine-Versionen: nur 28.0.4 und 29.9.2 gemessen, die Port-Doc sagt es; die Stelle, an der die Form wechselt, ist unbekannt (F-9). `name:group` als Benutzer: nicht gemessen (F-8). Mit Terminal auf 28.0.4: nicht erkannt, Grenze (a) (Prüfpunkt 3).

**2. Wächterfälle und Mutationen.** Gepinnt sind Status (1 mit derselben Zeile), Fehlerausgabe, zweite nicht leere Zeile (davor, danach) und anderer Benutzer; „kein Terminal“ ist gleichwertig gedeckt (siehe M41). Ungepinnt bleiben der Namens-Präfix (F-1) und die Hälfte „nichts abgeschnitten“ von `singleLine` (F-3). M36 bis M40 sind begründet rot (M37, M39, M40 selbst bestätigt). **M41 ist eine gleichwertige Mutante, die Begründung stimmt:** mit Terminal ist `cmd.Stdout = opts.Stdout` (in `runSession`), der Puffer der Standardausgabe bleibt leer, `firstLine()` ist leer und `singleLine()` ergibt `false`; die Bedingung `!opts.Terminal` ist dort doppelt gedeckt und dokumentiert die Absicht wie bei der OCI-Erkennung.

**3. Port-Doc und Zusage.** Wahr für beide Engines ohne Terminal (Matrix, Unit-Fälle beider Formen, Docker-Läufe laut Slice); mit Terminal nur für 29.9.2 (Fehlerausgabe). Der Satz „on the engines measured (28.0.4 and 29.9.2, client 29.9.2), the user does not exist after the probe“ ist für sich gelesen für 28.0.4 mit Terminal zu weit; der letzte Satz der Doc („with Terminal … does not see lines on that stream“) und der Wächterfall *terminal, status 126, the missing-user line* tragen die Einschränkung (F-6, INFO). „Any other message of the client is part of the session“ gilt weiter. **Konsistenz mit dem Lastenheft:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) bindet Exit-Code `12` an das Scheitern der *Prüfung* vor dem Öffnen; die Prüfung ist die Probe (`probeShell`, `CombinedOutput` und Fehlerstatus, unabhängig von Strom, Text und Engine). Die Lücke (a) betrifft nur den Benutzer, der *nach* der Probe verschwindet, und nur mit Terminal auf einer Engine in der Form 28.0.4; dort sieht der Mensch die Meldung im Terminal, nur der Exit-Code ist `0`. Das ist gegen das Lastenheft **vertretbar**, solange die Spezifikation die Grenze nennt (sie steht in (a), F-7). Exit `12` für Fehler des Clients nach der Probe ist die Entscheidung des Projektinhabers (Slice §3 *CI-Fix*), nicht der Wortlaut des Lastenhefts; die Klasse bleibt `ErrDevcontainerRuntime`.

**4. Engine-Matrix (§3.10).** Der Abschnitt trägt Datum, Engine (Digest des `dind`-Images), Client 29.9.2, Plattform und Basis-Images einmal vor der Tabelle; „gleiche Zelle heißt: gleich gemessen“ ist die einzige Lesehilfe und reicht. Die 29.9.2-Spalte stimmt mit der früheren Messtabelle des Slice überein (Zeilen *Benutzer fehlt*, *Workspace fehlt*, *Container gestoppt*, *Engine nicht erreichbar*). Der Satz „Den Mechanismus … habe ich nicht gelesen; gemessen ist nur die Form in der Tabelle“ ist ehrlich; „die Form der Meldung hängt an der Engine“ davor ist ein Schluss aus „gleicher Client, andere Engine“ und wird durch den Runner gestützt (zweite Installation von 28.0.4, derselbe Text). Offen bleibt der Status 126 auf dem Runner (F-5).

**5. Kommentare und Hard Rules.** §3.7: die neuen Kommentare stehen im Indikativ und tragen Zusage, Abgrenzung oder Grenze. §3.9 (a): kein Treffer in den hinzugefügten Zeilen für „former“, „früher“, „originally“, „no longer“, „is gone“, „used to“, „previously“, „anymore“, „would“, „wäre“, „würde“, „instead of“, „ohne“. (b): `clientUserMissingText` wird in keinem Kommentar unter dem alten Wert genannt; die Verweise `[Launcher.clientMessageOnStderr]`, `[ociFailureOnStdout]`, `[userMissingOnStdout]` existieren; der Kommentarkopf der Konstanten („Docker 29.9.2“) bleibt wahr; **aber** die Plan-Aussage zur Engine-Abhängigkeit im Folgeplan ist durch den Fix veraltet (F-4). (c): `userMissingOnStdout` und `singleLine` tragen je ihren Doc-Kommentar und stehen nicht zwischen einem Doc-Kommentar und seiner Funktion. Keine Suppression (`nolint`: 0 Treffer im Diff), `gofmt` über `make lint` grün. Commit-Messages: siehe Sensor-Läufe.

**6. Robusterer Weg.** Siehe F-9 (nur bewertet, nicht umgebaut).

## Findings

0 HIGH, 0 MEDIUM, 4 LOW, 5 INFO. Kein HIGH und kein MEDIUM, weil die falschen Positive an fünf zugleich geltende Bedingungen gebunden sind (kein natürlicher Erzeuger außer einem `docker exec` im Inneren der Sitzung gegen eine Engine der Form 28.0.4) und denselben Zielkonflikt tragen wie die benannte Grenze (d), die der Erstreview als INFO geführt hat (F-10 dort).

**F-1 · LOW** `internal/adapter/driven/docker/devcontainer_shell.go` · „strings.HasPrefix(line, clientUserMissingText+opts.User)“ — Namensgrenze fehlt
- Klasse: Bindung an den Benutzer ohne Namensgrenze
- Quelle: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („Hat die Shell begonnen, endet der Befehl mit `0`“), [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 3 (Grenze)
- Befund: Slice, Commit-Message und Doc-Kommentar sagen „die Zeile nennt den Benutzer der Sitzung“; der Code prüft ein Präfix, also auch `unable to find user vscode2: …` für die Sitzung als `vscode` (Sonde in der Kopie: die Zeile wird als Fehler des Clients gewertet, Exit 2 des Tests). Der Wächterfall *…of another user* nutzt `other`, das kein Präfix des Sitzungsbenutzers ist.
- Risiko: Failure-Szenario: Sitzung als `dev` ohne Terminal, im Inneren `docker exec -u developer <anderer Container> …` gegen eine Engine der Form 28.0.4 als einzige Ausgabe und letzter Befehl; die Zeile `unable to find user developer: …` beginnt mit `unable to find user dev`, eine begonnene Sitzung endet mit Exit `12` statt `0`. Sehr eng, aber enger gefasst, als der Text zusagt. Der Vergleich der Fehlerausgabe (`strings.Contains`) hat dieselbe Eigenschaft und ist unverändert.
- Verifizierbar: `make test` mit dem Sondenfall (hier gelaufen, Exit 2)

**F-2 · LOW** [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) · „(d) eine Sitzung ohne Terminal, deren einzige Ausgabe eine Zeile“ — neue Klasse falscher Positive nicht benannt
- Klasse: Grenze der Erkennung am Meldungstext
- Quelle: `AGENTS.md` §3.10 (Bedingung in Spec ohne Test; die Grenzen gehen in den Spezifikations-Slice), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Entscheidung 4 (Zusage-Menge)
- Befund: Der CI-Fix erweitert (a) und (b), nicht (d); „eine Sitzung ohne Terminal, deren einzige Ausgabe die Zeile `unable to find user <Benutzer der Sitzung>…` ist und die mit Status 126 endet, gilt als Fehler des Clients“ steht in keiner benannten Grenze.
- Risiko: Der Spezifikations-Slice leitet seine Grenzen aus (a) bis (d) ab; die Klasse fehlt dort, die Spezifikation sagt dann „nur bei erkanntem Fehler“ ohne den Vorbehalt.
- Verifizierbar: Review-only

**F-3 · LOW** `internal/adapter/driven/docker/devcontainer_shell.go` · „if b.total > len(b.data) {“ — Bedingung ohne Wächterfall
- Klasse: Bedingung der Erkennung ohne Wächterfall (zweites Auftreten im Slice; erstes: R4 im Erstreview, behoben)
- Quelle: `AGENTS.md` §3.10 (Bedingung mit Test oder gestrichen; der Folgeplan führt „Jede Bedingung mit ihrem Test in der Beleg-Tabelle oder gestrichen“)
- Befund: Der Slice führt „nichts abgeschnitten“ als Bedingung von `singleLine`; M39 prüft nur die zweite Hälfte (zweite nicht leere Zeile). Die Mutation ohne die Abschneide-Prüfung läuft grün (hier gelaufen, Exit 0). Kein Fall liefert eine erste Zeile und danach mehr als 4096 Byte, von denen im Kopf nur Leerraum steht.
- Risiko: Der nächste Änderer von `singleLine` streicht die Prüfung, kein Test wird rot; die Zusage „genau diese eine Zeile“ gilt für lange Ausgaben nicht mehr.
- Verifizierbar: `make test` mit Mutation M-A (hier gelaufen, Exit 0)

**F-4 · LOW** [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../plan/planning/done/slice-v2-spec-devcontainer-shell-und-zweites-up.md) · „auf der Engine 28.0.4 des CI-Runners (Client 29.9.2) wird der Fall“ — Plan-Aussage durch den Fix veraltet
- Klasse: Plan-Aussage nach Änderung nicht nachgezogen
- Quelle: `AGENTS.md` §3.10 (Aussagen über Verhalten; Abgleich nach der letzten Änderung), §3.9 (b) sinngemäß
- Befund: Der Plan (aus `80c88dca`) sagt, auf der Engine 28.0.4 werde „Benutzer nach der Probe entfernt“ nicht erkannt und `Shell` liefere keinen Fehler. Seit `8fa83a34` ist das ohne Terminal falsch; `f25b456b` berührt die Datei nicht, der Abgleich des CI-Fix führt `open/` nicht auf.
- Risiko: Der Spezifikations-Slice schreibt die Grenze aus dem Satz ab; der Verweis „ist vor dem Schreiben zu lesen“ mildert, nennt aber den Ist-Stand nicht.
- Verifizierbar: Review-only (Zuständigkeit Planner)

**F-5 · INFO** [`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md) · „Der Exit-Status des Clients im fehlgeschlagenen Fall steht nicht im Log.“ — Status 126 nur am Docker-in-Docker gemessen
- Klasse: Messung an Ersatzaufbau
- Quelle: `AGENTS.md` §3.10
- Befund: Die Bindung `exit == 126` stützt sich auf die Messung gegen `docker:28.0.4-dind`; der Lauf des Runners (host-installierte 28.0.4, gleicher Text) zeigt den Status nicht. Zuständigkeit Verifier: der Lauf `integration-docker` auf dem neuen Push-Tip entscheidet, ob der Status auf dem Runner 126 ist.
- Risiko: lokal PASS, Runner weiter rot, falls der Status dort ein anderer ist.
- Verifizierbar: CI-Lauf `integration-docker`

**F-6 · INFO** `internal/hexagon/port/driven/devcontainer_launcher.go` · „measured (28.0.4 and 29.9.2, client 29.9.2), the user does not exist“ — Zusage dicht, Einschränkung nur im letzten Satz
- Klasse: Doc-Kommentar mit Einschränkung außerhalb der Zusage
- Quelle: `AGENTS.md` §3.7
- Befund: Die Aufzählung nennt „der Benutzer fehlt nach der Probe“ für beide Engines ohne Vorbehalt zum Terminal; mit Terminal auf 28.0.4 trägt erst der letzte Satz („the launcher does not see lines on that stream“). Ableitbar, nicht falsch; der Wächterfall pinnt es.
- Risiko: Der Befehl-Slice liest die Aufzählung allein.
- Verifizierbar: Review-only

**F-7 · INFO** `internal/adapter/driven/docker/devcontainer_shell_test.go` · „terminal, status 126, the missing-user line“ — Grenze (a) nur im Fake gepinnt
- Klasse: Grenze der Erkennung am Meldungstext
- Quelle: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), `AGENTS.md` §3.10
- Befund: „Befehl endet dann mit `0`“ (Terminal, 28.0.4) folgt aus dem Code und dem Fake-Fall; ein Lauf mit dem echten Client und Terminal auf 28.0.4 steht nicht in der Matrix als Test (die Familie `B18` läuft ohne Terminal). Bewertung zu Prüfpunkt 3: vertretbar; der Spezifikations-Slice nennt die Grenze mit dieser Herkunft.
- Risiko: gering
- Verifizierbar: Review-only (der Verifier kann mit Terminal gegen 28.0.4 nachmessen)

**F-8 · INFO** `internal/hexagon/port/driven/devcontainer_launcher.go` · „`name:group` pass;“ — `name:group` ungemessen
- Klasse: Messlücke bei zulässiger Eingabe
- Quelle: `AGENTS.md` §3.10
- Befund: Der Port lässt `name:group` zu; die Erkennung vergleicht den ganzen Wert. Ob die Meldung der Engine den ganzen Wert oder nur den Namen nennt, steht in keiner Matrix-Zeile; nennt sie nur den Namen, würde keine der beiden Formen erkannt. Der Benutzer kommt aus dem Ergebnis von `Up` (`remoteUser`), ein Paar ist dort nicht zu erwarten.
- Risiko: gering
- Verifizierbar: Review-only

**F-9 · INFO** [`BEO-ADP/erkennung-am-meldungstext-engine-version`](../plan/planning/observations/BEO-ADP/erkennung-am-meldungstext-engine-version/observation.md) · „Kein Gate prüft lokal gegen eine zweite Engine-Version.“ — Bewertung zu Prüfpunkt 6
- Klasse: Grenze der Erkennung am Meldungstext
- Quelle: Register `BEO-ADP/erkennung-am-meldungstext-engine-version`; [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Re-Evaluierungs-Trigger
- Befund (Bewertung für Architect und Folge, kein Auftrag an diesen Slice): Ein Weg, der nicht am Text hängt, ist die **Gegenprobe nach dem Ende des Clients**: Endet der Client mit 1, 126 oder 127, läuft die Probe noch einmal; scheitert sie jetzt, gilt die Sitzung als nicht begonnen. Vorteile: unabhängig von Strom, Text und Engine-Version, und sie schließt auch die Lücke (a), weil die Probe kein Terminal braucht. Kosten: Sie ersetzt die falschen Negative durch eine andere Klasse falscher Positive (eine Sitzung, die selbst Benutzer, Workspace oder Shell entfernt und dann mit 1, 126 oder 127 endet, ergäbe Exit `12` statt `0`), widerspricht der Grenze in [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 3 („Was in der Sitzung geschieht, ändert den Exit-Code nicht“) und braucht eine Entscheidung des Architects, keine des Implementers. Ein Positivsignal „Sitzung hat begonnen“ über die Shell ändert die Argumentliste ([`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Gegenoption D); die Engine-API (Exec-Inspect) verlangte eine neue Abhängigkeit am Adapter. Die Erkennung am Text bleibt in jedem Fall best effort; das Lastenheft ist durch die Probe (Prüfung vor dem Öffnen) von der Engine-Form unabhängig erfüllt. Empfehlung: Weg A so lassen, die Engine-Version der Tests oder eine zweite Version im CI nach dem Trigger des Registers entscheiden lassen (die Runner-Engine ändert sich ohne Zutun des Repos), die Gegenprobe erst bei einem zweiten Vorgang als ADR prüfen.
- Risiko: Jede weitere Engine-Form (28.1 bis 29.8 ungemessen) verlangt Code, Test und Matrix-Zeile.
- Verifizierbar: Review-only

## Geprüft ohne Befund

- `internal/adapter/driven/docker/devcontainer_shell.go`: Verkettung der Erkennung (Fehlerausgabe, OCI, neu); vorhandene Erkennung unverändert; CRLF und LF; `headBuffer`-Zähler.
- `internal/adapter/driven/docker/devcontainer_shell_test.go`: ein positiver Fall und sechs Wächterfälle entsprechen dem Slice (acht Fälle in `…ClientErrorBeforeSessionIsRuntimeFailure`, 26 in `…SessionEndIsNoError`, wie B8 und B9 sie zählen).
- Architekturgrenzen: Imports von `devcontainer_shell.go` und Port unverändert (Diff); kein neuer Sentinel, deshalb kein Dual-Classifier-Fall; `make verify-depguard` nicht nötig.
- Sprachvertrag: die neuen Texte im Code sind Englisch.
- Plan-Tabelle gegen `git diff 87420c45..HEAD --stat -- internal`: drei Dateien, 57 Einfügungen, 9 Löschungen, wie im Abgleich des Slice; die „unverändert“-Pfade des Abgleichs haben keinen Diff (F-4 betrifft `open/`, nicht diese Liste).
- Nutzertexte, Lastenheft, Spezifikation, ADRs, `.golangci.yml`, `Makefile`, `Dockerfile`: kein Diff in `87420c45..f25b456b`.

## Nicht geprüft

- Lauf der Docker-Tests (`make test-docker`, Familie `TestShell_Docker_` auf beiden Engines) und der CI-Lauf `integration-docker` auf `f25b456b`: Auftrag „lesen“; der CI-Lauf läuft parallel und gehört dem Verifier.
- Messungen gegen Engine 28.0.4 (Docker-in-Docker): nicht wiederholt (kein Docker-Aufruf außerhalb von `make`).
- `make gates` als Ganzes (`coverage-gate`, `baseline-verify`, `ablage-pruefung`): nicht gelaufen; der Diff berührt drei Go-Dateien und Slice-Text, `lint`, `test` und `docs-check` liefen einzeln; die Coverage 91,90 % ist die Meldung des Implementers.
- DoD-Erfüllung und Closure: Aufgabe des Verifiers.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 5 |

**Finding-Klassen:** Bindung an den Benutzer ohne Namensgrenze (F-1) · Grenze der Erkennung am Meldungstext (F-2, F-7, F-9; im Register `BEO-ADP/erkennung-am-meldungstext-engine-version`) · Bedingung der Erkennung ohne Wächterfall (F-3; zweites Auftreten im Slice) · Plan-Aussage nach Änderung nicht nachgezogen (F-4) · Messung an Ersatzaufbau (F-5) · Doc-Kommentar mit Einschränkung außerhalb der Zusage (F-6) · Messlücke bei zulässiger Eingabe (F-8).

## Verdikt

**Merge-blockierend:** nein. **Freigabe für die Verifier-Nachprüfung: GO.** Die Erkennung tut, was der Slice sagt; die Wächterfälle pinnen Status, Fehlerausgabe, zweite Zeile und anderen Benutzer; die Port-Doc ist für die gemessenen Engines wahr; die Matrix ist als Messung gekennzeichnet. F-1 und F-3 sind kleine Änderungen an Code und Unit-Tests (kein Docker-Lauf nötig); F-2 und F-4 sind Plan-Text für den Planner. Der Beleg für §5 bleibt der CI-Lauf auf dem Push-Tip (F-5).
