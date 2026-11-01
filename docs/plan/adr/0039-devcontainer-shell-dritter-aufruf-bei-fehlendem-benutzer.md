# ADR 0039: `u-boot devcontainer shell` — ein dritter `docker exec` im Container, nur wenn die einzige Ausgabe einer Sitzung ohne Terminal die Zeile `unable to find user <Benutzer>` bei Status 126 ist; die Zusage „im Container nur die Probe und die Sitzung“ aus ADR-0034 gilt für alle anderen Fälle

**Status:** Deprecated

**Grund:** Der Projektinhaber hat diese ADR am 2026-11-01 vor `Accepted` verworfen; der dritte Aufruf ist aus dem Code genommen. Die Messung ([Messprotokoll](../../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md), Engine und Client 29.9.2) zeigt: Die Zeile `unable to find user …` entsteht dort auf der Fehlerausgabe mit Status 1, nicht auf der Standardausgabe mit Status 126; der dritte Aufruf greift auf 29.9.2 nie. Die Form der Engine 28.0.4 ist weder bestätigt noch widerlegt (keine Engine 28.0.4 verfügbar). Die Zusage „im Container nur die Probe und die Sitzung“ aus ADR-0034 gilt unverändert; die Teil-Ablösung unten ist nie in Kraft getreten. Der folgende Text ist der Stand beim Verwerfen.

**Datum:** 2026-11-01

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [ADR-0031](0031-devcontainer-shell-ohne-editor.md), [ADR-0034](0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Unterabschnitt `shell`: Punkt *Sitzung* und Punkt *Grenzen der Erkennung*, soweit sie die Zahl der Aufrufe im Container und den Restfall nennen; die Spezifikation führt sie nach dem Code, der Stand der Commits `000ae5ca` und `2f09b5dd` ist dort gelesen). Das Lastenheft ändert diese ADR nicht; der Restfall steht ihm gegenüber offen (Entscheidung 4).

**Supersedes (teilweise):** eine Stelle von [ADR-0034](0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), und nur diese: in Entscheidung 4, Aufzählung *Zugesagt* (Test), das erste Glied bis zum Semikolon („`u-boot` ruft im Container nur die Probe und die Sitzung auf (`TestLauncher_Shell_ProbesThenOpensTheSession`);“). Für diese Stelle ist diese ADR die geltende Fassung (Wortlaut in Entscheidung 1; dieselbe Form wie [ADR-0034](0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) gegenüber ADR-0031). Alles Übrige von ADR-0034 und ADR-0031 gilt unverändert; beide bleiben `Accepted` und werden nicht geändert. Die Zeile der Fitness Function von ADR-0034 zu `TestLauncher_Shell_ProbesThenOpensTheSession` („zwei `docker exec`“) bleibt wahr: Der Test fährt den Normalfall ohne die Zeile.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt: Hat die Shell begonnen, endet der Befehl mit `0`, unabhängig von ihrem Exit-Status; scheitert vor dem Öffnen der Shell die Prüfung, endet er mit `12`. Der Launcher erkennt einen Fehler des Clients vor Beginn der Sitzung an der Meldung des Clients, weil der Exit-Status ihn nicht von dem einer Shell trennt (`clientErrorBeforeSession`, `internal/adapter/driven/docker/devcontainer_shell.go`, gelesen am 2026-11-01). Auf der Engine 28.0.4 steht die Meldung zum fehlenden Benutzer ohne Präfix auf der Standardausgabe, mit Status 126 (`userMissingOnStdout`, ebd.). Dieselbe Zeile kann ein `docker exec -u <Benutzer der Sitzung> <anderer Container>` **in** der Sitzung als einzige Ausgabe schreiben; der Launcher trennte beides nicht, der Befehl endete dann mit `12` statt `0` (Befund V-1 des Verification-Reports zur Spezifikation von `shell`; Gegenstand des Slice zur Erkennung eines fehlenden Benutzers bei `shell`).

Der Projektinhaber hat am 2026-11-01 Weg A gewählt (die Zeile des Clients von der der Sitzung trennen, Code statt Lastenheft-Ausnahme). Der Slice setzt ihn um (Commits `000ae5ca` und `2f09b5dd`). Der Weg verlangt, nach dem Ende des Clients einmal im Container nachzusehen, ob der Benutzer fehlt; das ist ein dritter `docker exec` und berührt die Zusage von [ADR-0034](0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Entscheidung 4: „`u-boot` ruft im Container nur die Probe und die Sitzung auf“. Eine `Accepted`-ADR wird nicht nachgebessert, daher diese Folge-ADR.

**Gelesen am 2026-11-01** (Aussagen über Code, `AGENTS.md` §3.10):

- `devcontainer_shell.go`, `userMissingOnStdout`: `found` ist wahr nur bei `!opts.Terminal && exit == 126 && stderr.total == 0 && stdout.singleLine()` und einer ersten Zeile, die mit `missingUserLead(opts.User)` beginnt (`unable to find user <Benutzer>:`); erst dann ruft sie `userMissingNow`. Mit Terminal sieht der Launcher den Strom nicht (Kommentar am Symbol; `Shell` reicht mit Terminal `cmd.Stdout = opts.Stdout` durch).
- `userMissingNow`: wiederholt `probeArgs(opts, shell)` (dieselbe Probe wie vor der Sitzung, mit der Shell der Sitzung) unter `context.WithTimeout(ctx, inspectorTimeout)` und liefert wahr genau dann, wenn die Probe mit einem Fehler endet und ihre Ausgabe `missingUserLead(opts.User)` enthält; eine gelingende oder anders scheiternde Probe überlässt die Zeile der Sitzung.
- `devcontainer_shell_test.go`: `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce` (zwei Fälle: Benutzer da, Benutzer weg; erwartet jeweils drei `docker exec` und die dritte Argumentliste gleich der ersten); `TestLauncher_Shell_ProbesThenOpensTheSession` (zwei `docker exec` im Normalfall).
- [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Punkt *Sitzung* und *Grenzen der Erkennung* (Stand der Fix-Runde): nennen die Wiederholung genau in diesem Fall und den Restfall mit Exit-Code `12`. `spec/architecture.md` Zeile 325 nennt den dritten Aufruf als Ausnahme von „Im Container ruft der Adapter nur die Probe und die Sitzung auf“ (angepasst).

**Messstand.** Der Erzeuger (ein `docker exec -u <Benutzer> <anderer Container>` in der Sitzung gegen eine Engine der Form 28.0.4) ist aus der Form abgeleitet und **nie gefahren**; die Engine 28.0.4 steht nur als Fake (`fakeDockerScript`) zur Verfügung, gelaufen sind Docker-Tests auf 29.9.2 (`TestShell_Docker_B18_SessionDoesNotBegin`). Die Messung des Erzeugers (DoD-Punkt 1 des Slice) ist offen. Die Entscheidung unten gilt unabhängig davon; findet die Messung die Zeile nicht, entfällt die Grenze und mit ihr diese ADR (Re-Evaluierungs-Trigger).

## Entscheidung

Wir wählen **den dritten Aufruf nur im einen Fall**.

1. **Geltende Fassung der einen Stelle.** Das erste Glied der Aufzählung *Zugesagt* in ADR-0034 Entscheidung 4 lautet künftig: „`u-boot` ruft im Container die Probe und die Sitzung auf und sonst nichts, mit einer Ausnahme: Endete die Sitzung ohne Terminal mit Status `126` und ist die Zeile `unable to find user <Benutzer der Sitzung>: …` die einzige Ausgabe (bei leerer Fehlerausgabe), wiederholt der Launcher nach dem Ende des Clients die Probe der Shell der Sitzung genau einmal (`TestLauncher_Shell_ProbesThenOpensTheSession` für den Normalfall mit zwei Aufrufen, `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce` für die Ausnahme mit drei);“.
2. **Wann der dritte Aufruf läuft.** Nur wenn alle Bedingungen von `userMissingOnStdout` zugleich gelten (ohne Terminal, Status 126, leere Fehlerausgabe, genau eine Zeile auf der Standardausgabe, Zeile beginnt mit `unable to find user <Benutzer>:`). In jedem anderen Fall, auch mit Terminal und auf der Engine 29.9.2 (dort mit Präfix auf der Fehlerausgabe, ohne Wiederholung), bleibt es bei zwei Aufrufen. Der Aufruf verändert nichts im Container: Er ist die Probe, die schon vor der Sitzung lief (`bash -c true`, sonst `sh -c true`, ohne `-i`, `-t` und Variablen).
3. **Urteil.** Meldet die wiederholte Probe den Benutzer als fehlend (Status ungleich `0`, Ausgabe mit derselben Zeile für den Benutzer der Sitzung), gilt die Zeile als Fehler des Clients vor der Sitzung: Exit-Code `12`. Gelingt die Probe oder scheitert sie anders, gilt die Zeile als Ausgabe der Sitzung: Exit-Code `0` ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Der Doppelpunkt ist die Namensgrenze (`vscode2` ist nicht `vscode`).
4. **Restfall, offen beim Projektinhaber.** Schreibt die Sitzung die Zeile selbst und verschwindet der Benutzer im Container danach (etwa `userdel`), meldet die wiederholte Probe ihn als fehlend, und der Befehl endet mit `12`, obwohl die Shell begonnen hat; das widerspricht der Zusage des Lastenhefts in diesem einen Fall. Die Spezifikation nennt ihn als Grenze (Spalte `0`, *Grenzen der Erkennung*). Diese ADR **akzeptiert den Restfall als Negativ**. Grund: Er verlangt zwei seltene Ereignisse zugleich, die Zeile der Sitzung als einzige Ausgabe **und** das Verschwinden genau dieses Benutzers; gemessen ist keines der beiden in der Sitzung; ihn zu schließen wäre Weg B. Ob die Grenze bleibt oder eine Lastenheft-CR (Weg B) folgt, entscheidet der Projektinhaber (Risiko 5 des Slice); beides ändert diese ADR nicht, weil eine CR das Lastenheft ändert und diese ADR den Code beschreibt.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Weg C des Slice): die Spezifikation trägt die Grenze weiter, die Zeile der Sitzung und die des Clients werden nicht getrennt | kein Code, keine Änderung der Zusage von ADR-0034 | der Befehl endet mit `12` nach begonnener Shell, wenn die Sitzung die Zeile als einzige Ausgabe schreibt; gegen die Zusage des Lastenhefts, ohne dass der Projektinhaber sie gelockert hätte; er hat A (Erkennung enger) gewählt |
| B — Ausnahme im Lastenheft (Weg B): ein Satz in [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | nur Text, die Zusage „zwei Aufrufe“ bliebe | Lastenheft-CR mit Wortlaut-Freigabe; senkt die Zusage des Vertrags für einen Fall, den der Code trennen kann; vom Projektinhaber nicht gewählt (bleibt Ausgang für den Restfall) |
| C — die Zeile anders erkennen, ohne zweiten Aufruf (Zeitpunkt der ersten Ausgabe; Marker der Sitzung) | die Zusage von ADR-0034 bliebe wörtlich | nicht gelesen und nicht gemessen: der Zeitpunkt hängt an Last und Engine, ein Marker verlangte eine Änderung der Argumentliste oder der Shell (gegen ADR-0031 Entscheidung 1 und ADR-0034 Entscheidung 3); eine Schätzung als Urteil über einen Exit-Code |
| **D — dritter Aufruf nur im einen Fall (gewählt)** | trennt die Zeile mit einer Beobachtung statt einer Schätzung; Vorbild im Code (`engineUnreachable` fragt die Engine ebenso nach dem Ende des Clients); im Normalfall und mit Terminal unverändert zwei Aufrufe; die Ausnahme ist eng, getestet und in der Spezifikation benannt | die Zusage „nur die Probe und die Sitzung“ gilt nicht mehr ohne Ausnahme; die Rückfrage beurteilt den Benutzer nur danach (Restfall, Entscheidung 4); ein Aufruf mehr im Container |
| E — dritter Aufruf in jedem Fall (auch bei Erfolg) | einfacherer Code | ein Aufruf mehr im Container, den Zusage und Handbuch-Weg nicht kennen, ohne Anlass in fast allen Sitzungen |

## Konsequenzen

- Positiv: Der Befehl endet nach begonnener Shell mit `0`, auch wenn die Sitzung die Zeile des fehlenden Benutzers als einzige Ausgabe schreibt, solange der Benutzer im Container noch da ist; der Code folgt der Zusage des Lastenhefts in dem Fall, den er trennen kann.
- Negativ: Eine Ausnahme von der Zusage „nur die Probe und die Sitzung“; der Restfall (Entscheidung 4) endet mit `12` trotz begonnener Shell; die Engine-Form 28.0.4 ist nur gegen den Fake gelaufen, der Erzeuger nie.
- **Akzeptierte Negative, mit Grund:** (i) *Restfall:* siehe Entscheidung 4. (ii) *Erzeuger ungemessen:* Die Entscheidung stützt sich auf eine abgeleitete Form; fällt die Messung negativ aus, entfällt die Grenze, und der Code bleibt harmlos (die Zeile träte nie auf). (iii) *Mit Terminal* bleibt die Erkennung auf die Fehlerausgabe beschränkt; die Zeile der Engine 28.0.4 liegt im Strom des Clients, den der Launcher nicht liest, der Befehl endet dann mit `0` (in der Spezifikation bereits benannt).
- Nachgezogen (im Slice): die Sicht [`spec/architecture.md`](../../../spec/architecture.md) Zeile 325 nennt die Ausnahme (die Sicht nennt keine ADR, `AGENTS.md` §3.4), die Spezifikation trägt sie. Das Handbuch 4.9 nennt den Restfall in Worten, ohne Aussage zur Zahl der Aufrufe. Findet die Messung des Erzeugers die Zeile nicht, löst eine Folge-ADR diese ab und streicht den Aufruf.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Unit | `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce`: genau dann drei `docker exec`, wenn die Zeile als einzige Ausgabe bei Status 126 ohne Terminal steht; die dritte Argumentliste gleicht der Probe; Fehler nur bei fehlendem Benutzer | `make test` |
| Unit | `TestLauncher_Shell_ProbesThenOpensTheSession`: im Normalfall zwei `docker exec`; `TestLauncher_Shell_SessionEndIsNoError` zählt die Aufrufe je Fall | `make test` |
| Integrationstest (`//go:build docker`) | `TestShell_Docker_B18_SessionDoesNotBegin`: die Sitzung schreibt die Zeile selbst (Guard); gelaufen auf 29.9.2, auf 28.0.4 offen | `make test-docker` |

Die Worte der Spezifikation zur Ausnahme sind nicht maschinell prüfbar; `AGENTS.md` §3.10 und das Review decken sie.

## Re-Evaluierungs-Trigger

Die Messung des Erzeugers findet die Zeile in der Sitzung auf keiner Engine: Grenze und diese ADR entfallen (Folge-ADR). Ein Nutzerbericht, dass `u-boot devcontainer shell` mit `12` endet, obwohl die Shell begonnen hat (Restfall beobachtet): Weg B neu bewerten. Eine Engine, die die Zeile in anderer Form schreibt, oder ein Client, der den fehlenden Benutzer auf der Fehlerausgabe meldet (so die Engine 29.9.2): Erkennung und Aufruf neu lesen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-01 | Proposed durch den Architect; Weg A vom Projektinhaber gewählt, Umsetzung in den Commits `000ae5ca` und `2f09b5dd`; `Accepted` setzt der Projektinhaber |
| 2026-11-01 | Deprecated: Messung auf Engine 29.9.2 (Zeile auf der Fehlerausgabe, Status 1), Projektinhaber nimmt den dritten Aufruf und diese ADR zurück | Messprotokoll, Slice zur Erkennung eines fehlenden Benutzers bei `shell` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0039`.
