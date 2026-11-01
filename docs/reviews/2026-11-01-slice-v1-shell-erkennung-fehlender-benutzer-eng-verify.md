# Review-Report: slice-v1-shell-erkennung-fehlender-benutzer-eng (Verifier) — 2026-11-01

**Review-Art:** Code — Verifikation gegen Plan, DoD, Entscheidungen und Sensors (Modul 11: „Bauen wir es richtig?“). Kein Review des Diffs und keine Validierung.

**Gegenstand:** Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/open/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (Weg A, Entscheidung des Projektinhabers 2026-11-01), HEAD `59fbd8a2` (Merge des Zweigs `worktree-agent-afb663c9f1176132c`), Diff `350cdc72..HEAD`: 13 Dateien, 324 Einfügungen, 28 Löschungen. Der Slice liegt in `open/`.

**Skill:** `.harness/skills/verifier.md` existiert nicht; die Rolle folgt `harness/roles.md` und `harness/verification.md`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/open/slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (§1 Gewählt, §2 DoD, §3 Plan, §5 Closure-Trigger, §6 Risiken)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
- [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Punkt *Sitzung*, Tabelle der Exit-Codes Spalte `0`, Punkt *Exit-Code `12`*, Grenzen der Erkennung, Änderungsprotokoll)
- [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md), [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md), [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) (Proposed)
- `spec/architecture.md` (Adapter `docker/`), [Review-Report](2026-11-01-slice-v1-shell-erkennung-fehlender-benutzer-eng.md) samt Nachlauf (F-N1 bis F-N4)
- `harness/verification.md`, `AGENTS.md` (Hard Rules)

**Sensors dieses Laufs:**

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (Docker-Modus, HEAD `59fbd8a2`) | pass | Exit 0; Lint „0 issues“; Coverage 92,20 % gegen Schwelle 90 %; Stempel `.harness/state/gates-passed.head` zeigt `59fbd8a2`; Arbeitsbaum danach sauber |
| `make doc-trace` | pass | 84 Anforderungen, 0 Waisen; [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) und [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) führen den Slice und [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) |
| `make docs-check` | pass (Teil von `make gates`, vor diesem Report); mit diesem Report erneut gefahren | Ergebnis im Abschlussbericht |
| Mutationen in einer Kopie im Scratchpad (`make test`, Docker-Modus; Produktbaum unberührt) | siehe zweite Tabelle | — |
| `make test-docker` | **nicht gefahren** | Entscheidung des Projektinhabers (nicht lokal). Der Implementer berichtet den Docker-Test des Slice (`TestShell_Docker_B18_SessionDoesNotBegin`, Fall „guard: the session itself writes the missing-user line …“) grün auf Engine 29.9.2; **so berichtet, nicht nachgeprüft**. Eine Engine 28.0.4 stand nicht zur Verfügung. |
| `make verify-depguard` | nicht gefahren | keine Importregel berührt (keine Änderung an `.golangci.yml` oder Importen im Diff) |

Mutationen (am Stand `59fbd8a2`, Kopie unter dem Scratchpad):

| Mutation | Ergebnis |
|---|---|
| `userMissingNow` liefert immer wahr | rot: `TestLauncher_Shell_SessionEndIsNoError` (3 Fälle mit erneuter Probe), `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce/the_user_exists` |
| erneute Probe entfällt (`return line, found`) | rot: dieselben Tests plus `.../the_user_is_missing_now` |
| `missingUserLead(opts.User)` durch `unable to find user` ersetzt (Benutzerbindung) | rot: `SessionEndIsNoError/…a_probe_that_reports_another_user` |
| `err != nil` in `userMissingNow` entfernt | nicht gefahren; der Review-Nachlauf liest sie als überlebend (die Probe nach der Sitzung schreibt die Zeile nur bei Fehler); kein realistischer Fehler |

---

## DoD-Abgleich (Slice §2, Punkt für Punkt)

| # | DoD-Punkt | Stand | Beleg / Grund |
|---|---|---|---|
| 1 | Messung des Erzeugers (Skript, Protokoll, zwei datierte Läufe je Engine) | **offen** | Kein Diff unter `tools/messung/` oder `docs/messungen/` (`git diff --stat`: keine Treffer); der Plan führt die Zeile ehrlich als „neu, offen“. Der Slice-Text (§1 *Gewählt*) sagt „ohne vorherige Messung des Erzeugers“. Das Änderungsprotokoll der Spec und `zusagen-und-messstaende.md` nennen: „ein `docker exec` gegen eine Engine 28.0.4 ist nicht gefahren“. |
| 2 | Der gewählte Weg | erfüllt für Weg A, mit einer benannten Abweichung vom Wortlaut | Code: `userMissingOnStdout` (Zeile 226) ruft bei erfüllter Zeilen-Bedingung `userMissingNow` (Zeile 236), das `probeArgs(opts, shell)` wiederholt; `clientErrorBeforeSession` (Zeile 165) reicht `shell` durch. Tests je Bedingung: `TestLauncher_Shell_SessionEndIsNoError` (Zähler 2 oder 3 je Fall; Fälle: Probe gelingt, scheitert anders, meldet `vscode2`), `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` (`goneLate`), `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce`; Mutationen oben rot. Spec nach dem Code mit Beleg-Tabelle im Änderungsprotokoll (Zeile `2026-11-01`). **Abweichung:** DoD verlangt „die Grenze (e) entfällt in [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) … und die Ausnahme in der Spalte `0`“. Sie entfällt nicht ganz: Sie bleibt als benannter Restfall (Grenze am Punkt zum fehlenden Benutzer, Spalte `0` „Ausnahme: der Restfall“, Punkt *Exit-Code `12`* „Ausnahme: der Restfall …“). Die DoD-Formulierung ist mit Weg A in dieser Form nicht erreichbar; sie ist nicht angepasst. |
| 3 | `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise; bei A `make test-docker` als CI-Lauf `integration-docker` | teilweise | `make gates` grün (selbst gefahren), `make doc-trace` 0 Waisen (selbst gefahren), `make docs-check` Teil von `make gates` (grün). `integration-docker` als CI-Lauf mit sichtbar bestandenen Tests liegt nicht vor: `make test-docker` ist laut Projektinhaber nicht lokal zu fahren, ein CI-Lauf ist für den Merge-Stand `59fbd8a2` nicht belegt. |
| 4 | Review durchgeführt, Report, Nachlauf-Review | erfüllt | [Review-Report samt Nachlauf](2026-11-01-slice-v1-shell-erkennung-fehlender-benutzer-eng.md) (frischer Kontext, Diff der Fix-Runden); F-1 bis F-8 und F-N1 bis F-N4 in den Fix-Commits `2f09b5dd`, `f234fe63`, `75f89727` behoben (von mir gegengelesen: F-N1 — `Folgepflicht` und „nicht angepasst“ kommen in [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) nicht mehr vor; F-N2 — Punkt *Exit-Code `12`* trägt die Ausnahme; F-N3 — Handbuch nennt den Restfall; F-N4 — Plan-Tabelle führt `observations/` als „offen, kein Diff“). |
| 5 | Doku-Update (Spec, Handbuch 4.9, CHANGELOG) | erfüllt | Spec: Spalte `0` (Zeile 898), Punkt *Sitzung* (Zeile 910), Punkt *Exit-Code `12`*, Grenzen (Zeilen 920 bis 927), Änderungsprotokoll; Handbuch 4.9 *Ende* (Zeile 863 ff.) nennt den Restfall in Worten ohne Kennungen (`AGENTS.md` §5 Regel 6); `CHANGELOG.md` unter `[Unreleased]`/`Changed`; `spec/architecture.md` Zeile 325 ohne ADR-, Slice- oder Wellenbezug (§3.4); `docs/maintainer/zusagen-und-messstaende.md` nachgezogen; Lastenheft nicht berührt (Weg A). |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | §7 lautet „(Wird bei der Closure gefüllt.)“. |
| 7 | Beobachtungs-Register fortgeschrieben | **offen** | Kein Diff unter `docs/plan/planning/observations/`. |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **offen** | Risiken 1 bis 4 tragen „offen bis zur Closure“; Risiko 5 trägt „weiter offen“ (Entscheidung beim Projektinhaber). |
| 9 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen** | Hängt an den Punkten 6 bis 8. |

Alle Kästchen in §2 stehen unmarkiert; das ist für einen Slice in `open/` stimmig.

## Entscheidungs-Konformität

| Entscheidung | Befund |
|---|---|
| Weg A (Projektinhaber 2026-11-01) | Umgesetzt als erneute Probe nach dem Ende des Clients (Kandidat „Rückfrage“, Muster `engineUnreachable`); §1 *Gewählt* nennt Weg und Entscheidungsdatum. Weg B (Lastenheft) und Weg C nicht gewählt. |
| Accepted-ADRs unverändert | bestätigt: `git diff --name-status 350cdc72..HEAD` führt keine Änderung an [`ADR-0031`](../plan/adr/0031-devcontainer-shell-ohne-editor.md) und an [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md); neu sind nur [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) und die Index-Zeile. |
| [`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md) Proposed | bestätigt: `Status: Proposed`; `Supersedes (teilweise)` nennt genau das erste Glied von Entscheidung 4 aus [`ADR-0034`](../plan/adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md); Index-Zeile in `docs/plan/adr/README.md`. Der Code ruft den dritten Aufruf bereits, bevor die ADR `Accepted` ist; das ist offen zu nennen (Finding V-3). |
| Lastenheft unverändert | bestätigt: `spec/lastenheft.md` nicht im Diff. Die Spec sagt selbst „das Lastenheft nennt dafür keine Ausnahme“: Der Restfall (`12` trotz begonnener Shell) steht dem Satz von [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gegenüber offen. |
| Hard Rules | §3.2: kein `//nolint` im Diff (Lint „0 issues“). §3.4: `spec/architecture.md` Zeile 325 ohne ADR-/Slice-Bezug. §3.5: keine Accepted-ADR geändert. §3.6: keine Gate-Lockerung (Coverage 92,20 % ≥ 90 %). |

## Plan-vs-Code (Plan §3 gegen `git diff --stat 350cdc72..HEAD`)

| Diff-Datei | Plan-Zeile |
|---|---|
| `internal/adapter/driven/docker/devcontainer_shell.go` | ja (Zeilen 165, 226, 236 gelesen, treffen zu) |
| `…/devcontainer_shell_test.go`, `…/devcontainer_shell_docker_test.go`, `…/devcontainer_launcher_test.go` | ja |
| `spec/spezifikation.md`, `spec/architecture.md` | ja |
| `CHANGELOG.md`, `docs/maintainer/zusagen-und-messstaende.md`, `docs/user/benutzerhandbuch.md` | ja |
| `docs/plan/adr/0039-…md`, `docs/plan/adr/README.md` | ja |
| Review-Report des Slice | ja |
| der Slice selbst (`Verantwortlich:`, §1, §3, §6) | ja (Zeile „diese Datei“) |

Plan-Einträge ohne Diff: `tools/messung/<name>/`, `docs/messungen/<kennung>/` (offen, DoD 1) und `docs/plan/planning/observations/` (offen, DoD 7), beide im Plan als „kein Diff“ gekennzeichnet; `spec/lastenheft.md` als „unverändert“ geführt, im Diff unberührt. Diff-Dateien ohne Plan-Zeile: keine. Der Report dieses Laufs ist neu und kommt hinzu.

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | HIGH | DoD-Punkt 1 (Messung des Erzeugers) ist offen. Die Notwendigkeit des dritten `docker exec` hängt an einer abgeleiteten, nie beobachteten Zeile; der Slice sagt selbst, „zeigt die Messung, dass der Erzeuger die Zeile nicht erzeugt, entfällt die Grenze, und der Slice geht … nach `done/`“ (Risiko 1 bleibt ohne Ausgang). Der Slice kann **so nicht** nach `done/` schließen, ohne dass die DoD geändert oder die Messung gefahren wird. | Slice §2 Punkt 1, §6 Risiko 1 | Slice · „Messung des Erzeugers“ | ja — Skript und Protokoll liegen nicht vor | Messung vor Wahl offen |
| V-2 | MEDIUM | Der Restfall ist Risiko 5 „weiter offen“ und widerspricht dem Lastenheft ([`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)) in einem Fall. Der Closure-Trigger (Slice §5) verlangt, dass die Spalte `0` und das Lastenheft dasselbe über „`0` nach Beginn der Shell“ sagen; das ist mit dem Restfall nicht erfüllt, solange das Lastenheft unverändert ist. | Slice §5, §6 Risiko 5 | Spec · „das Lastenheft nennt dafür keine Ausnahme“ | nein — Entscheidung des Projektinhabers | Zusage und Vertrag nicht deckungsgleich |
| V-3 | MEDIUM | Der Code mit dem dritten `docker exec` ist gemergt, die Entscheidung dazu ([`ADR-0039`](../plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md)) ist `Proposed`; die Sicht und die Spec führen die Ausnahme von „nur die Probe und die Sitzung“ bereits als geltend. | `AGENTS.md` §3.5; Slice §6 Risiko 3 | `docs/plan/adr/0039-…md` · „**Status:** Proposed“ | nein — Accept durch den Projektinhaber | Folge-ADR vor Annahme im Code wirksam |
| V-4 | MEDIUM | `integration-docker` als CI-Lauf mit sichtbar bestandenen Tests (DoD 3, bei A verlangt) liegt nicht vor; der Docker-Test des Slice lief nur im Lauf des Implementers auf Engine 29.9.2. Die Form der Engine 28.0.4 ist nur durch Einheitstests mit einem Skript an Stelle von `docker` belegt (Risiko 2: Vorsorge „Test auf beiden Engines“ nicht eingelöst). | Slice §2 Punkt 3, §6 Risiko 2 | Spec · Änderungsprotokoll `2026-11-01` · „nicht auf der Engine 28.0.4“ | ja — CI-Lauf `integration-docker` | Docker-Beleg aussteht |
| V-5 | LOW | DoD-Punkt 2 verlangt, dass „die Grenze (e) entfällt“; sie bleibt als Restfall bestehen. Die DoD-Formulierung ist nicht an Weg A angepasst. | Slice §2 Punkt 2 | Slice · „die Grenze (e) entfällt in …“ | Review-only | DoD-Wortlaut hinter dem Weg zurück |
| V-6 | LOW | `open/README.md` führt den Slice noch mit „Messung, dann Weg … nur Plan“ und ohne Weg A; der Plan nennt den Nachzug „beim Zusammenführen“. Die Zeile ist zu Weg A veraltet. | `AGENTS.md` §3.10 (Zustandsmarker) | `docs/plan/planning/open/README.md` · „Messung, dann Weg (…) als Entscheidung des Projektinhabers; nur Plan“ | Review-only | Zustandsmarker nach Ablauf |
| V-7 | INFO | Die Mutation `err != nil` in `userMissingNow` überlebt (vom Review-Nachlauf bereits gelesen; hier nicht neu gefahren). Kein realistischer Fehlerfall. | Maintainability | `devcontainer_shell.go` · `return err != nil && strings.Contains(…)` | ja — `make test` mit der Mutation | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Code `devcontainer_shell.go` gegen Spec (Bedingungen: kein Terminal, Status 126, leere Fehlerausgabe, eine Zeile, Benutzerpräfix, erneute Probe mit Fehler und derselben Zeile) | geprüft, ohne Befund |
| Zähler „zwei oder drei `docker exec`“ gegen Punkt *Sitzung* | geprüft, ohne Befund |
| Spec-Spalte `0`, Punkt *Exit-Code `12`*, Grenze, Änderungsprotokoll: untereinander widerspruchsfrei | geprüft, ohne Befund |
| Handbuch 4.9, CHANGELOG, `spec/architecture.md`: decken sich mit Code und Spec; keine Kennungen im Handbuch | geprüft, ohne Befund |
| Die zwei Accepted-ADRs und das Lastenheft unverändert | geprüft, ohne Befund |
| Plan §3 gegen Diff | geprüft, ohne Befund (siehe Tabelle) |
| Gates, Coverage, Suppression, Importregeln | geprüft, ohne Befund |
| Docker-Test auf einer Engine 28.0.4 und `make test-docker` | nicht geprüft (nicht gefahren, siehe Sensors) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Messung vor Wahl offen · Zusage und Vertrag nicht deckungsgleich · Folge-ADR vor Annahme im Code wirksam · Docker-Beleg aussteht · Zustandsmarker nach Ablauf

## Bewertung: kann der Slice schließen? (keine Entscheidung)

Nein, nicht in der heutigen Fassung. Offen sind DoD 1 (Messung), 6 bis 9 (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen), der CI-Lauf `integration-docker` aus DoD 3 und der Closure-Trigger aus §5 (Spalte `0` und Lastenheft deckungsgleich, mit Restfall nicht erfüllt). Die Entscheidung liegt beim Projektinhaber und Planner; Wege, die sich aus dem Slice-Text ergeben (hier nur genannt, nicht gewählt):

1. Die Messung fahren (Engine 28.0.4 und die Engine des Hosts, je zwei Läufe); zeigt sie, dass der Erzeuger die Zeile nicht schreibt, entfällt der dritte Aufruf (Slice nach `done/` mit `Gegenstand:` und Grund, Code und Folge-ADR zurückzunehmen).
2. DoD 1 ändern oder in einen Folge-Slice schieben (Planner), wenn die Engine 28.0.4 nicht verfügbar ist; dann trägt der Slice die Aussage „Erzeuger abgeleitet“ ausdrücklich in der Closure-Notiz, mit Risiko 1 und 2 als „weiter offen“ (zulässig nur mit einem Ausgang je Risiko).
3. Den Restfall entscheiden: Grenze belassen (der Closure-Trigger §5 ist dann anzupassen, weil Lastenheft und Spalte `0` nicht dasselbe sagen) oder Lastenheft-CR (Weg B) als eigener Slice.
4. Die Folge-ADR (Proposed) annehmen oder zurückziehen.

## Verdikt

**Merge-blockierend:** n. a. (der Zweig ist gemergt). **Closure-blockierend:** ja — wegen V-1 (HIGH) und V-2 bis V-4 (MEDIUM). Was gebaut ist, entspricht Plan §3, der Spec und Weg A; `make gates` und `make doc-trace` sind grün, die Mutationen färben die Tests rot. Eine Verifikation des Verhaltens auf Engine 28.0.4 und im CI-Lauf `integration-docker` steht aus.

**Übergabe:** Bericht an den Planner (Rückkante Verifier → Planner). Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Docker-Test des Slice ist vom Implementer berichtet, nicht von mir nachgeprüft.
