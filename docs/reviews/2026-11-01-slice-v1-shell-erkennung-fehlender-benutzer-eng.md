# Review-Report: slice-v1-shell-erkennung-fehlender-benutzer-eng — 2026-11-01

**Review-Art:** Code — Code-Review gegen Plan, Spec und Konventionen (Weg A).

**Gegenstand:** Commit `000ae5ca` gegenüber `f7f6c710` (Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md), Weg A). Acht Dateien, 91 Einfügungen, 15 Löschungen.

**Skill:** `.harness/skills/reviewer.md` @ `000ae5ca`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (§1 Ziel und Wege, §2 DoD, §3 Plan, §6 Risiken)
- [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Unterabschnitt `shell`, Tabelle der Exit-Codes, Änderungsprotokoll)
- `AGENTS.md` (Hard Rules, besonders §3.5, §3.7, §3.9, §3.10)

**Sensors:** `make gates` lief im Worktree auf `000ae5ca` und endete grün (Stempel `.harness/state/gates-passed.head` zeigt `000ae5ca`); der Beleg des Implementers lag nur für `f7f6c710` vor. `make test-docker` lief nicht (Entscheidung des Projektinhabers); der neue Docker-Test ist ungelaufen gelesen. `make verify-depguard` nicht nötig (keine Importregeln berührt). Die zwei „rot gesehenen" Mutationen habe ich nicht neu gefahren (keine Produktdatei-Änderung erlaubt); ihre Tragfähigkeit ist aus den Tests gelesen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Restfall „die Sitzung schreibt die Zeile selbst, und der Benutzer fehlt bei der erneuten Probe" (etwa nach `userdel` in der Sitzung, dann `docker exec -u <Benutzer> <anderer Container>` mit Status 126) endet mit `12`, obwohl die Shell begonnen hat. Weder die Spalte `0` („unabhängig vom Exit-Status der Shell", Ausnahme gestrichen) noch das Änderungsprotokoll („eine Sitzung, die dieselbe Zeile selbst erzeugt, endet mit `0`") noch die Liste der erkannten Meldungen nennen diese Grenze; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt `0` nach begonnener Shell. Der Slice-Plan nennt sie als Weg-A-Folge nirgends. | `AGENTS.md` §3.10 (Bedingung in der Spec ohne Entsprechung im Code); [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | `spec/spezifikation.md` · „gilt die Zeile als Sitzung (ein `docker exec -u <Benutzer der Sitzung> <anderer Container>` in der Sitzung kann dieselbe Zeile als einzige Ausgabe erzeugen)" und Zeile `\| 2026-11-01 \|` | nein — Review-only (Spec-Lesung gegen `userMissingNow`) | Spec-Zusage ohne benannte Restgrenze |
| F-2 | MEDIUM | Die Bedingung „genau dann einmal … sonst bleibt es bei den zwei Aufrufen" ist nur an einer Seite gepinnt: `TestLauncher_Shell_ProbesThenOpensTheSession` zählt zwei `exec` bei Status 0, `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce` drei bei der vollen Bedingung. Die Fälle von `TestLauncher_Shell_SessionEndIsNoError` (Terminal, andere Ausgabe davor oder danach, nicht leere Fehlerausgabe, Status ungleich 126) prüfen nur „kein Fehler", nicht die Zahl der `exec`; eine Mutation, die bei jedem Status ungleich 0 oder ohne die Zeilen-Bedingung erneut probt, bliebe grün, weil die Probe dort gelingt. | `AGENTS.md` §3.10 (Spec-Bedingung braucht einen pinnenden Test) | `internal/adapter/driven/docker/devcontainer_shell_test.go` · `func TestLauncher_Shell_SessionEndIsNoError`; `spec/spezifikation.md` · „sonst bleibt es bei den zwei Aufrufen" | ja — `make test` mit der Mutation „erneute Probe vor der Zeilen-Bedingung" | Bedingung der Spec nur einseitig gepinnt |
| F-3 | LOW | Die erneute Probe gilt nur für „denselben Benutzer"; kein Fall lässt die Probe die Zeile eines anderen Benutzers melden (`vscode2` bei Sitzung `vscode`). Eine Mutation von `missingUserLead(opts.User)` auf den allgemeinen Präfix `unable to find user` in `userMissingNow` bleibt grün (`goneLate` trägt in beiden Tests den Benutzer der Sitzung). | `AGENTS.md` §3.10 | `devcontainer_shell.go` · `strings.Contains(string(out), missingUserLead(opts.User))` | ja — `make test` mit der Mutation | Benutzer-Bindung der Probe ungepinnt |
| F-4 | LOW | Die Aussage „`u-boot` ruft im Container nur die Probe und die Sitzung auf" steht in einer Accepted-ADR als „Zugesagt (Test)" und gilt mit dem dritten `docker exec` nicht mehr uneingeschränkt; die Spec führt die Ausnahme, eine Folgeentscheidung mit Verweis auf die ADR (oder ein Vermerk dort) fehlt. Slice-Risiko 3 sah das Nachziehen „mit dem Code" vor. | `AGENTS.md` §3.5 (Folgeentscheidung mit Verweis) | `docs/plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md` · „ruft im Container nur die Probe und die Sitzung auf" | Review-only | Accepted-ADR-Aussage überholt, keine Folgeentscheidung |
| F-5 | LOW | Plan-Tabelle gegen `git diff --stat`: alle acht Dateien des Diffs sind geführt, die Zeile `tools/messung/<name>/`, `docs/messungen/<kennung>/` (Art „neu") hat keinen Diff-Gegenpart. Mehrere Fundstellen sind nach dem Nachziehen auf `f7f6c710` veraltet: [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) „Zeilen 842 bis 844 und Zeile 815" (jetzt um 895 und 920 bis 925), `userMissingOnStdout` „Zeilen 223 bis 228", Handbuch „Zeile 860" in der DoD gegenüber „Zeile 863" in der Tabelle (gelesen: 863 trifft zu). | `AGENTS.md` §3.10 (Tabelle nach der letzten Runde abgleichen) | [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) · „Zeilen 842 bis 844 und Zeile 815 (Spalte `0`)" | Review-only | Plan-Fundstellen nach Rebase veraltet |
| F-6 | INFO | DoD-Punkt 1 (Messung des Erzeugers) ist offen und im Diff ehrlich unberührt; das Änderungsprotokoll benennt, dass die Form der Engine 28.0.4 nur durch das Skript an Stelle von `docker` belegt ist und der Docker-Test auf 29.9.2 lief. Slice §1 verlangt „Jeder Kandidat ist zu messen, bevor er gewählt wird"; die Wahl von Weg A steht im Slice-Text nicht, und die Risiken 1 bis 3 tragen noch „offen bis zur Closure" (Risiko 3 ist mit dem Diff eingetreten, Risiko 2 hat die Vorsorge „Test auf 28.0.4" nicht eingelöst). Das ist ein Ausgang für die Closure, kein Fehler des Diffs; zeigt die Messung, dass der Erzeuger die Zeile nicht erzeugt, bleibt der dritte `exec` ohne Nutzen. Zuständig: Implementer (Ausgänge), Verifier (Beleg). | Slice §2, §6 | [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) · „Messung des Erzeugers" | Review-only | Messung vor Wahl offen |
| F-7 | INFO | Die Änderungsprotokoll-Zeile nennt keine Lastenheft-Version („deckungsgleich mit Lastenheft 0.xx" fehlt, alle Nachbarzeilen führen sie). Das ist hier vertretbar, solange F-1 offen ist, weil eine Deckungsgleichheit mit dem Lastenheft sonst behauptet würde; mit der Benennung des Restfalls ist die Zeile erneut zu lesen. | Maintainability | `spec/spezifikation.md` · „„[`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Unterabschnitt shell, deckungsgleich mit dem Code" | Review-only | — |
| F-8 | INFO | Die zwei Mutationen sind aus dem Diff nicht belegt (kein Mutationsprotokoll im Repo); gelesen sind sie tragfähig: „erneute Probe entfällt" färbt `MissingUserLineRepeatsTheProbeOnce`/`userExists` und den Docker-Guard rot, „Probe ignoriert die Ausgabe" färbt die Fälle mit `probe-late-not-running` rot. Die Flags `user-gone-late` und `probe-late-not-running` greifen erst, wenn `-i` in `argv.log` steht, trennen also Probe vor und nach der Sitzung sauber. Verifier prüft die Mutationen selbst. | Maintainability | `internal/adapter/driven/docker/devcontainer_launcher_test.go` · `user-gone-late` | ja — Verifier | — |

## Zu den Prüfpunkten

1. **Trennung der Zeilen.** `userMissingNow` wiederholt die Probe (`probeArgs`, Shell der Sitzung) nach dem Ende des Clients; sie trennt die Zeile der Sitzung von der des Clients, solange der Benutzer zur Zeit der erneuten Probe vorhanden ist. Der Restfall (Benutzer verschwindet nach der Zeile, Fehler `12`) ist nicht benannt: F-1. Ein echter Fehler des Clients, bei dem die erneute Probe aus anderem Grund scheitert (Container weg), gilt als Sitzung (`0`); das nennt die Spec ausdrücklich („scheitert sie anders, gilt die Zeile als Sitzung").
2. **Tests.** Je Bedingung vorhanden: Probe gelingt, Probe scheitert anders, Probe meldet den Benutzer (zwei Fälle `vscode`, `vscode2`), genau drei `exec`. Lücken: F-2, F-3.
3. **Spec-Bedingungen.** Die Bedingung im Satz „Sitzung" ist teils ungepinnt (F-2); die Änderungsprotokoll-Zeile ohne Lastenheft-Version ist vertretbar (F-7).
4. **Handbuch.** Gelesen: `docs/user/benutzerhandbuch.md` Zeile 863 („Hat die Shell begonnen, endet der Befehl mit Exit-Code 0, was immer der Exit-Status der Shell war") nennt keine Ausnahme; „unverändert" trifft zu. Das Handbuch nennt den dritten Aufruf nicht und braucht es nicht.
5. **Kommentare.** Ohne Befund gegen §3.7 und §3.9: Der Doc-Kommentar von `userMissingOnStdout` steht im Indikativ, nennt Zusage und Kopplung (`[Launcher.userMissingNow]`); `userMissingNow` steht nach der Funktion, deren Kommentar davor steht; der Querverweis in `clientErrorBeforeSession` ist auf `[Launcher.userMissingOnStdout]` nachgezogen. Der alte Name `userMissingOnStdout` kommt im Baum nur noch im Plan eines abgeschlossenen Slice ([`slice-v1-devcontainer-shell-fundament`](../plan/planning/done/slice-v1-devcontainer-shell-fundament.md)) als historische Fundstelle vor, nicht in Code oder Kommentaren.
6. **Plan-Tabelle.** F-5.
7. **Basis.** Der Diff gegenüber `f7f6c710` berührt keine Instanz-Logik (`Label`, Instanz-Code nicht berührt); ohne Befund.
8. **DoD-Punkt 1.** Ehrlich offen, Ausgang folgt bei der Closure: F-6.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapter/driven/docker/devcontainer_shell.go` (Argumente, Kontext, Zeitgrenze `inspectorTimeout`, Fehlerabbildung) | geprüft, ohne Befund (Abbruch des Kontexts macht die Probe zur Sitzung, entspricht der Spec) |
| Dual-Classifier / Exit-Code-Abbildung (`ErrDevcontainerRuntime`) | geprüft, ohne Befund (kein neuer Sentinel) |
| Sprachvertrag (Fehlertexte) | geprüft, ohne Befund (keine neuen Texte) |
| `internal/adapter/driven/docker/devcontainer_shell_docker_test.go` | geprüft, ohne Befund (Guard „die Sitzung schreibt die Zeile selbst" pinnt den Weg auf der Engine des Hosts; ungelaufen, `make test-docker` ausgesetzt) |
| `CHANGELOG.md`, `docs/maintainer/zusagen-und-messstaende.md` | geprüft, ohne Befund |
| `docs/user/benutzerhandbuch.md` | geprüft, ohne Befund (Zeile 863 gelesen) |
| Gates, Suppression, Importregeln | geprüft, ohne Befund (`make gates` grün; keine `//nolint`) |
| `spec/lastenheft.md` | nicht geprüft für Änderungen: unverändert im Diff; Befund F-1 betrifft die Spec gegen das Lastenheft |
| Mutationen, Docker-Tests auf 28.0.4 | nicht gefahren (Rolle Verifier; `make test-docker` ausgesetzt) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Spec-Zusage ohne benannte Restgrenze · Bedingung der Spec nur einseitig gepinnt · Plan-Fundstellen nach Rebase veraltet

## Verdikt

**Merge-blockierend:** ja — wegen der beiden MEDIUM-Findings (F-1, F-2). Kein HIGH; die Umsetzung trennt die Zeile der Sitzung von der des Clients wie geplant.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ersetzt keine Verifikation — DoD-Konformität (Messung, `make test-docker` als `integration-docker`, Mutationen) prüft der Verifier.

## Nachlauf (Fix-Runden `2f09b5dd` und `f234fe63`)

Geprüft im frischen Kontext: `git diff 000ae5ca..HEAD`, [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) und [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md) gegen die Basis. Kein Produktcode geändert (`devcontainer_shell.go` seit `000ae5ca` unverändert); keine Messung, kein `make test-docker`. Die Einheitstests wurden nicht ausgeführt (`PKG`/`RUN` gibt es nur im lokalen Modus, im Worktree-Kontext steht keine Toolchain bereit); die Mutationen unten sind aus dem Quelltext gelesen.

### Stand der Vorgänger-Findings

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (Restfall) | behoben; die Entscheidung des Projektinhabers bleibt offen | [`spec/spezifikation.md`](../../spec/spezifikation.md): Spalte `0` („Ausnahme: der Restfall“), **Grenze:** am Punkt zum fehlenden Benutzer, Änderungsprotokoll; [Lastenheft](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) unverändert, die Spec sagt selbst „das Lastenheft nennt dafür keine Ausnahme“; Slice §6 Risiko 5 „weiter offen“ |
| F-2 (Zähler einseitig) | behoben | `TestLauncher_Shell_SessionEndIsNoError`: Feld `reprobes`, Erwartung 3 nur bei den drei Fällen mit der Zeile des Sitzungsbenutzers, sonst 2 |
| F-3 (`vscode2`) | behoben | Fall `goneLate: "vscode2"` bei Sitzung als `vscode` (Exit 0, drei Aufrufe) und in `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` die Sitzung als `vscode2` (Exit 12) |
| F-4 ([ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md)) | behoben, Rest siehe F-N1 | [[ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md)](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md), Proposed |
| F-5, F-6 | behoben | Plan-Tabelle gegen die Dateiliste des Diffs, Fundstellen nachgezogen, „Gewählt: Weg A“ im Slice |

### Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-N1 | LOW | Die ADR führt die Sicht als „nicht angepasst“ (Kontext: „`spec/architecture.md` Zeile 325 sagt dagegen noch …“) und als offene Pflicht (Konsequenzen: „Folgepflicht (Nachzug im Slice …): die Sicht … wird um die Ausnahme ergänzt“). `f234fe63` hat Zeile 325 angepasst; der Satz der Sicht nennt den dritten Aufruf bereits. Die ADR ist Proposed und noch änderbar; nach `Accepted` bliebe der veraltete Satz stehen, und ein Leser hielte die Sicht für nachzuziehen. | `AGENTS.md` §3.10 (Zustandsmarker mit Ablaufpunkt) | `docs/plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md` · „(nicht angepasst; Folgepflicht)“ und „Folgepflicht (Nachzug im Slice“ | Review-only | Zustandsmarker nach Ablauf |
| F-N2 | LOW | Der Punkt *Exit-Code `12`* sagt: „Der Exit-Status der Shell nach dem Beginn der Sitzung ist kein Fehler: … der Befehl endet mit `0`.“ ohne den Restfall; Spalte `0` und die Grenze zwei Absätze weiter nennen ihn. Der nächste Änderer der Spec liest den Satz als unbedingt. | `AGENTS.md` §3.9 (b) | `spec/spezifikation.md` · „der Befehl endet mit `0`.“ (Punkt *Exit-Code `12`*) | Review-only | Spec-Zusage ohne benannte Restgrenze |
| F-N3 | LOW | Das Handbuch 4.9 *Ende* sagt „Hat die Shell begonnen, endet der Befehl mit Exit-Code 0, was immer der Exit-Status der Shell war.“ ohne Ausnahme; die Spec nennt jetzt den Restfall. Der Plan führt die Datei bewusst als „unverändert“ (Zeile 863 gelesen); die Grenzen (a) bis (d) der Spec fehlen im Handbuch ebenso, der Restfall ist die nächste. Ein Nutzer erfährt ihn nur über das CHANGELOG. | Maintainability (Doku gegen Spec) | `docs/user/benutzerhandbuch.md` · „endet der Befehl mit Exit-Code 0, was immer der Exit-Status der Shell war“ | Review-only; ob das Handbuch Grenzen nennt, entscheidet laut Slice der Slice | Spec-Zusage ohne benannte Restgrenze |
| F-N4 | INFO | Plan-Tabelle: Die Zeile „`docs/reviews/`, `docs/plan/planning/observations/`, Art „neu / update“ führt `observations/`, der Diff berührt es nicht (folgt mit der Closure); nach §3.10 wird ein Eintrag ohne Diff berichtigt oder als offen gekennzeichnet. In §6 steht eine Leerzeile zwischen Risiko 4 und 5. | `AGENTS.md` §3.10 | Slice §3 · „observations/“ | Review-only | Plan-Eintrag ohne Diff |

### Zu den Prüfpunkten

1. **Spec-Grenze und Widerspruchsfreiheit.** Spalte `0`, Punkt *Sitzung* („sonst bleibt es bei den zwei Aufrufen“), Grenze und Änderungsprotokoll sind untereinander widerspruchsfrei und decken sich mit dem Code (`userMissingOnStdout`: kein Terminal, Status 126, leere Fehlerausgabe, eine Zeile, Präfix `missingUserLead`; `userMissingNow`: Fehler und Ausgabe enthält `missingUserLead`). Das Lastenheft bleibt unverändert, der Widerspruch zu ihm ist offen benannt statt verdeckt. Weitere „unbedingt 0“-Stellen: Spec Punkt *Exit-Code `12`* (F-N2), Handbuch 4.9 (F-N3). Das CHANGELOG trägt den Satz „Limit: …“ und deckt sich mit der Spec; `docs/maintainer/zusagen-und-messstaende.md` und `spec/architecture.md` sagen kein „unbedingt 0“.
2. **Tests (aus dem Quelltext gelesen).** (a) Wiederholung entfernt oder Zeile immer als Fehler des Clients: der Fall „user exists“ und die drei `reprobes`-Fälle schlagen fehl. (b) Bedingung von `userMissingOnStdout` gelockert (Terminal, `singleLine`, Status, Fehlerausgabe, Präfix): die Fälle ohne `reprobes` erwarten zwei Aufrufe und fielen auf drei. (c) `missingUserLead` ohne Doppelpunkt oder ohne Benutzer: der Fall `goneLate: "vscode2"` ergäbe Exit 12, der Fall „a user whose name begins with the session user“ drei statt zwei Aufrufe. (d) `err != nil` in `userMissingNow` entfernt: überlebt, weil eine gelingende Probe die Zeile nicht schreibt; kein realistischer Fehler. Der Fake trennt die Probe nach der Sitzung über `-i` in `argv.log`. Der Restfall hat mit dem Fall „the user is missing now“ in `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce` (Fehler erwartet) einen pinnenden Test (§3.10).
3. **[ADR-0039](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md).** Form MADR vollständig (Kontext, Entscheidung, Alternativen, Konsequenzen, Fitness Function, Re-Evaluierungs-Trigger, Geschichte). `Supersedes (teilweise)` zitiert das erste Glied von [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) Entscheidung 4 (*Zugesagt*) wörtlich richtig; die Aussage zur Fitness-Zeile „zwei `docker exec`“ stimmt mit [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) überein; `Schärft` nennt [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes). `git diff --stat` gegen die Basis: [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md) und [ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) unverändert. Index-Zeile vorhanden (Proposed). Die Aussagen über Code tragen Fundstelle und Lesedatum und decken sich mit dem Quelltext; Messstand ehrlich (Erzeuger nie gefahren, Engine 28.0.4 nur als Fake). Rest: F-N1.
4. **`spec/architecture.md` Zeile 325.** Keine ADR-, Slice- oder Wellenbezüge (§3.4, im Diff 0 Treffer); der Satz (dritter Aufruf nur bei Sitzung ohne Terminal, Status `126`, Zeile über fehlenden Benutzer als einziger Ausgabe) deckt sich mit `userMissingOnStdout`/`userMissingNow` und [SPEC-049](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes); die leere Fehlerausgabe fehlt dort, folgt aber aus „einziger Ausgabe“.
5. **Plan-Tabelle** gegen `git diff --stat` der Basis bis HEAD (12 Dateien): jede Datei hat eine Zeile; `tools/messung/…` ist als „offen, kein Diff“ ehrlich gekennzeichnet; Rest F-N4.
6. **Kommentare §3.7/§3.9.** Neu sind die Feldkommentare (`goneLate`, `reprobes`), die Fehlermeldung des Zählers und die Doc-Kommentare der zwei Tests: Indikativ, Zusage oder Kopplung, kein Chronik- oder Konjunktivmuster, kein Treffer für „früher/former/no longer“; kein Helfer zwischen Doc-Kommentar und Funktion eingefügt; kein Verweis auf einen umbenannten Namen.

### Negativbefund

Geprüft, ohne Befund: `internal/adapter/driven/docker/devcontainer_shell_test.go` (Zähler, `vscode2`), `devcontainer_launcher_test.go` (Fake-Flags), `devcontainer_shell_docker_test.go` (Guard), `CHANGELOG.md`, `docs/plan/adr/README.md`, [ADR-0031](../plan/adr/0031-devcontainer-shell-ohne-editor.md)/[ADR-0034](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) unverändert, `spec/architecture.md`, `docs/maintainer/zusagen-und-messstaende.md`. Nicht geprüft: Ausführung der Einheitstests und `make gates` (Verifier), Docker-Tests (Vorgabe), die Messung des Erzeugers (DoD-Punkt 1, offen).

### Verdikt Nachlauf

0 HIGH, 0 MEDIUM, 3 LOW, 1 INFO. **Merge-blockierend:** nein. F-N1 ist vor `Accepted` der ADR zu berichtigen; F-N2 und F-N3 sind kleine Nachzüge oder eine Begründung im Slice. Die Entscheidung über den Restfall (Grenze behalten oder Lastenheft-CR, Weg B) bleibt beim Projektinhaber.
