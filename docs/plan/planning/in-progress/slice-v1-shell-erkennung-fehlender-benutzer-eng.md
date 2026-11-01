# Slice v1: Erkennung eines fehlenden Benutzers bei `u-boot devcontainer shell` enger fassen, oder die Ausnahme im Lastenheft nennen (Grenze (e), Verification V-1)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** Release `v0.17.0` (geplant, Roadmap *Nächste Wellen*; die Welle-Datei entsteht bei der Eröffnung). Folge-Punkt von [`welle-devcontainer-shell-und-zweites-up`](../done/welle-devcontainer-shell-und-zweites-up.md) und von [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../done/slice-v2-spec-devcontainer-shell-und-zweites-up.md) (Verification V-1), eigener Slice; er hält die Closure der Welle und den Release `v0.14.0` nicht auf, weil die Spezifikation die Abweichung in der Spalte `0` ihrer Tabelle der Exit-Codes und im Unterabschnitt `shell` benennt (§1). Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden wäre.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („Hat die Shell begonnen, endet der Befehl mit `0`, unabhängig von ihrem Exit-Status“; „Scheitert vor dem Öffnen der Shell die Prüfung … endet er mit Exit-Code `12`“, Lastenheft `0.21.0`, Zeile 494, gelesen am 2026-10-22), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (`10`, `11`, `12` für fachliche Fehler); [`ADR-0034`](../../adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) als Herkunft der Sitzung.

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Unterabschnitt `shell`, Liste der erkannten Meldungen mit Grenze (e), Zeilen 842 bis 844 am 2026-10-22 gelesen; Spalte `0` der Tabelle der Exit-Codes, Zeile 815), je nach Weg das Lastenheft ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)). Der Verweis zeigt aufwärts: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Zusage „Exit-Code `0` nach Beginn der Shell“ des Lastenhefts gilt ohne die Ausnahme der Grenze (e) — oder das Lastenheft nennt die Ausnahme.

**Befund** (aus dem Verification-Report [V-1](../../../reviews/2026-10-22-slice-v2-spec-devcontainer-shell-und-zweites-up-verify.md); Code gelesen am 2026-10-22):

- **Die Grenze.** Eine Sitzung ohne Terminal auf einer Engine der Form 28.0.4, deren Standardausgabe genau die Zeile `unable to find user <Benutzer der Sitzung>: …` ist, bei leerer Fehlerausgabe und Status `126`, gilt als Fehler des Clients vor Beginn der Sitzung und endet mit `12` (`userMissingOnStdout`, `internal/adapter/driven/docker/devcontainer_shell.go`, Zeilen 223 bis 228). Dieselbe Zeile erzeugt ein `docker exec -u <Benutzer der Sitzung> <anderer Container>` in der Sitzung gegen eine Engine dieser Form als einzige Ausgabe; der Launcher trennt beides nicht, der Befehl endete dann mit `12` statt `0`. Der Erzeuger in der Sitzung ist aus der Form abgeleitet, **nicht gefahren**.
- **Die Spezifikation sagt es.** Die Liste der erkannten Meldungen in [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) endet mit dieser Aussage, und die Spalte `0` der Tabelle der Exit-Codes nimmt „die Grenzen der Erkennung“ aus. Das Lastenheft nennt keine Ausnahme; bei Konflikt gilt es (Rang 1). Der Verifier wertet die Spezifikation nicht als Erweiterung des Lastenhefts, weil sie die Abweichung benennt und nichts behauptet, was das Lastenheft verneint.
- **Die Grenzen (a) bis (d)** der Spezifikation sind kein Gegenstand: Dort gilt der Satz des Lastenhefts zu `12` und die Bindung von [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); der Code endet für die nicht erkannten Fälle mit `0`.
- **Ein Vorbild im Code.** `engineUnreachable` (dieselbe Datei, Zeile 237) fragt die Engine nach dem Ende des Clients noch einmal, bevor die Meldung `Cannot connect …` als Fehler des Clients gilt (Aufruf in `clientMessageOnStderr`, Zeile 199).

**Gewählt:** Weg A (Entscheidung des Projektinhabers am 2026-11-01) wurde als erneute Probe nach dem Ende des Clients umgesetzt (Kandidat „Rückfrage“, Muster `engineUnreachable`; Commits `000ae5ca` bis `75f89727`) und **nach der Messung zurückgenommen** (Entscheidung des Projektinhabers 2026-11-01: „ja beides zurücknehmen“, der dritte `docker exec` und die Folge-ADR). Grund: Das [Messprotokoll](../../../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md) (Engine und Client 29.9.2) zeigt die Zeile auf der Fehlerausgabe mit Status 1, nicht auf der Standardausgabe mit Status 126; der dritte Aufruf greift dort nie. Die Form der Engine 28.0.4 ist weder bestätigt noch widerlegt (keine Engine 28.0.4 verfügbar). Der Gegenstand des dritten Aufrufs entfällt; `ADR-0039` (Datei `docs/plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md`) ist `Deprecated`. Geblieben sind Messung, Messprotokoll und die Berichte; die Spezifikation führt den Fall „die Sitzung schreibt die Zeile selbst“ als Grenze der Erkennung (Lesung des Codes, kein gefahrener Lauf des Launchers), das Lastenheft bleibt unverändert.

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, Erkennung enger fassen** (Code-Slice; Spezifikation nach dem Code): die Zeile des Clients wird von einer in der Sitzung erzeugten getrennt. Kandidaten, **nicht gelesen und nicht gemessen**: eine Rückfrage nach dem Ende, ob der Benutzer im Container vorhanden ist (Muster `engineUnreachable`; sie wäre ein zweiter Aufruf im Container neben der Probe und ändert die Zusage „im Container ruft der Adapter nur die Probe und die Sitzung auf“, `TestLauncher_Shell_ProbesThenOpensTheSession`); der Zeitpunkt der ersten Ausgabe; ein Marker der Sitzung. Jeder Kandidat ist zu messen, bevor er gewählt wird.
- **B, Ausnahme im Lastenheft nennen** (Lastenheft-CR mit Wortlaut-Freigabe des Projektinhabers): ein Satz in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) zu der Meldung, die die Sitzung selbst erzeugt; die Spezifikation bleibt.
- **C, keine Änderung:** die Spezifikation trägt die Grenze weiter; Begründung in §7 (der Erzeuger ist abgeleitet und nie beobachtet).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Grenzen (a) bis (d) der Erkennung** — sie fallen unter die Bindung von [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) an `10`, `11` oder `12` und nennen keinen Widerspruch zu einer begonnenen Sitzung (Verification V-1).
- **Ein Test für `Cannot connect` bei nicht leerer Standardausgabe** (Verification V-3, Review I-3) — die Spezifikation behauptet das Gegenteil nicht; ein Test ist Code und eine eigene Entscheidung, nicht Teil dieser Frage.
- **Die Erkennung mit Terminal** — der Launcher liest den Terminal-Strom des Clients nicht; die Aussage steht in der Spezifikation (Grenze (a)) und bleibt.
- **Colima, andere Engines, andere Hosts** — nicht zugesagt ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Messung des Erzeugers:** `docker exec -u <Benutzer der Sitzung> <anderer Container>` in der Sitzung gegen die Engine des Hosts (29.9.2) ergibt die Zeile, die der Launcher als Fehler des Clients liest — oder nicht; mindestens zwei Läufe, datiert, mit Umgebung (Engine, Client, Architektur). Eine Engine der Form 28.0.4 stand nicht zur Verfügung; diese Hälfte der Messung ist **nicht gefahren** und bleibt in Risiko 1 und 2 offen. *Beleg:* [Messprotokoll](../../../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md) mit Rohdaten (zwei Läufe, je Zelle zwei Wiederholungen), Skript [`tools/messung/shell-benutzer-fehlt/messen.sh`](../../../../tools/messung/shell-benutzer-fehlt/messen.sh) ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)). Ergebnis: Auf 29.9.2 erzeugt der Erzeuger die Zeile auf der Fehlerausgabe mit Präfix und Status 1, nicht in der Form von `userMissingOnStdout`; die Grenze entfällt dadurch nicht (siehe Protokoll, Urteil zur Kernfrage).
- [ ] **Der gewählte Weg:** Weg A ist umgesetzt und nach der Messung zurückgenommen (Entscheidung des Projektinhabers 2026-11-01); der Gegenstand des dritten Aufrufs entfällt, der Code steht wie vor dem Slice. In [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) steht der Fall „die Sitzung schreibt die Zeile selbst“ (Exit-Code `12` trotz begonnener Shell) als benannte Grenze der Erkennung, ausgenommen in der Spalte `0` über „die Grenzen der Erkennung“ (gemessene Form der Engine 29.9.2; die Einstufung als Fehler des Clients ist eine Lesung des Codes, kein gefahrener Lauf); das Lastenheft bleibt unverändert. *Beleg:* der Revert-Commit und der Spec-Diff dieses Zweigs (`git log`), die Fälle von `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure`.
- [ ] `make gates` grün; `make docs-check` 0 Befunde; `make doc-trace` ohne Waise. `make test-docker` läuft nicht (Entscheidung des Projektinhabers 2026-11-01); der Docker-Test-Code steht wie vor dem Slice ([`AGENTS.md`](../../../../AGENTS.md) §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Ändert eine Fix-Runde Vertragstext, liest ein Nachlauf-Review im frischen Kontext genau diesen Diff, bevor der Verifier startet.
- [ ] Doku-Update: Spezifikation (Grenze der Erkennung, Änderungsprotokoll), Wartungsdoku `zusagen-und-messstaende.md`; Lastenheft unverändert; Handbuch 4.9 (*Ende*, `docs/user/benutzerhandbuch.md` Zeile 863 am 2026-11-01 gelesen) und `CHANGELOG.md` unverändert, weil der Code wie vor dem Slice steht und keine Zusage hinzukommt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt, mit der Stelle, an der der Planner die Aussage gelesen hat (2026-10-22). Sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen (`AGENTS.md` §3.10).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/shell-benutzer-fehlt/messen.sh`, `docs/messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/` (`README.md`, `rohdaten-lauf1.txt`, `rohdaten-lauf2.txt`), `docs/messungen/README.md` | neu / update (DoD-Punkt 1) | Skript und Protokoll der Messung des Erzeugers auf der Engine 29.9.2; Bestandszeile im Index; eine Engine 28.0.4 stand nicht zur Verfügung. Muster: `tools/messung/devcontainer-instanz/messen.sh`, gelesen |
| `docs/plan/planning/open/README.md`, `docs/plan/planning/in-progress/roadmap.md` | update | Zeile des Slice in `open/README.md` (Weg A, Zeiger auf den Stand) und Zeiger unter *Offene Wellen* der Roadmap, beim Start gelesen |
| `internal/adapter/driven/docker/devcontainer_shell.go`, `devcontainer_shell_test.go`, `devcontainer_shell_docker_test.go`, `devcontainer_launcher_test.go` | Weg A umgesetzt (`000ae5ca` bis `75f89727`) und zurückgenommen (Revert-Commit); Netto-Diff gegenüber dem Stand vor Weg A: leer | `userMissingNow`, `TestLauncher_Shell_MissingUserLineRepeatsTheProbeOnce`, der Docker-Guard-Fall und die Fake-Flags `user-gone-late` und `probe-late-not-running` entfallen; `userMissingOnStdout`, `clientErrorBeforeSession` und `TestLauncher_Shell_ClientErrorBeforeSessionIsRuntimeFailure` stehen wie vor dem Slice, nach dem Revert gelesen |
| `spec/spezifikation.md` | update | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes): Punkt *Sitzung* („zwei Aufrufe und sonst nichts“, ohne Ausnahme) steht wie vor dem Slice; neu der Absatz „Nicht zu trennen ist der Fall …“ unter *Grenzen der Erkennung* und eine Zeile im Änderungsprotokoll; Spalte `0` (Zeile 898) nennt „die Grenzen der Erkennung“ und bleibt, nach dem Revert gelesen |
| `spec/lastenheft.md` | update (Weg B), sonst **unverändert** (Weg A gewählt) | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile 503 gelesen; Wortlaut-Freigabe des Projektinhabers vor dem Start |
| `docs/maintainer/zusagen-und-messstaende.md` | update | Stand der unbelegten Beobachtung zur Zeile der Sitzung (Abschnitt zur Spezifikation von `shell`, Punkt „Grenzen der Erkennung“): Messung auf 29.9.2, Fall der Sitzung ohne Test |
| `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `spec/architecture.md` | **unverändert** (nach dem Revert) | Der Eintrag unter `[Unreleased]` entfällt; Handbuch Zeile 863 (*Ende*) steht wie vor dem Slice; die Sicht (Adapter-Absatz zu `docker/`, „Im Container ruft der Adapter nur die Probe und die Sitzung auf“) nennt wieder keinen dritten Aufruf, gelesen |
| `docs/plan/adr/0039-devcontainer-shell-dritter-aufruf-bei-fehlendem-benutzer.md`, `docs/plan/adr/README.md` | update | Folge-ADR auf `Deprecated` gesetzt (nie `Accepted`), mit Grund-Absatz; Index-Zeile; [`ADR-0034`](../../adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) unverändert |
| `docs/reviews/2026-11-01-slice-v1-shell-erkennung-fehlender-benutzer-eng.md`, `…-verify.md` | neu (Chronik, bleibt) | Review- und Verifikationsbericht zum zurückgenommenen Stand |
| diese Datei | update | Verantwortlicher gesetzt; Stand nach der Zurücknahme (§1, §2, §3, §5, §6) |
| `docs/plan/planning/observations/` | update, **offen** (kein Diff; folgt mit der Closure) | Evidence im Beobachtungs-Register |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; der Release `v0.14.0` ist geschnitten (Tag gesetzt, beobachtbar in `git tag`); der Projektinhaber hat den Weg (A, B oder C) entschieden und, bei B, den Wortlaut freigegeben; kein anderer Slice in `in-progress/` (WIP-Limit 1). **Früher**, wenn ein Nutzer oder der Validator der Welle Exit-Code `12` nach begonnener Sitzung meldet: dann ist der Erzeuger beobachtet, nicht abgeleitet, und der Slice hat Vorrang.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg B verlangt einen Lastenheft-CR vor dem Code; dann ist der CR ein eigener Slice mit Wortlaut-Freigabe.
- `in-progress` → `open` (blockiert): die Engine der Form 28.0.4 steht für die Messung nicht zur Verfügung, oder der Projektinhaber entscheidet den Weg nicht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` und `make docs-check` sind grün; die Spalte `0` der Tabelle der Exit-Codes in [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) nimmt die Grenzen der Erkennung aus und nennt dort den Fall „die Sitzung schreibt die Zeile selbst“ (Exit-Code `12` trotz begonnener Shell), und [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) bleibt unverändert (Entscheidung des Projektinhabers 2026-11-01; die Spezifikation sagt selbst, dass das Lastenheft keine Ausnahme nennt); Review-Report, ein Nachlauf-Review über jeden Diff von Vertragstext und der Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Die Ausgänge stehen bei der Closure; bis dahin nennt jede Zeile die Vorsorge.

1. **Der Erzeuger erzeugt die Zeile nicht** (abgeleitet, nie gefahren). — **Ausgang:** offen bis zur Closure. Stand der Messung 2026-11-01 ([Protokoll](../../../messungen/slice-v1-shell-erkennung-fehlender-benutzer-eng/README.md)): Auf der Engine 29.9.2 erzeugt der Erzeuger die Zeile, aber auf der Fehlerausgabe mit Präfix und Status 1, nicht in der Form von `userMissingOnStdout`; der dritte Aufruf griffe dort nie und ist zurückgenommen. Die Form der Engine 28.0.4 ist nicht gemessen.
2. **Eine engere Erkennung verliert einen Fall des Clients.** — **Ausgang:** offen bis zur Closure. Die engere Erkennung ist zurückgenommen; der Code steht wie vor dem Slice, die Fälle des Clients bleiben erkannt. Geblieben ist der Befund der Messung: Der Fall „die Sitzung schreibt die Zeile selbst, mit dem Namen des Benutzers der Sitzung“ hat auf 29.9.2 dieselbe Form wie die Zeile des Clients (Fehlerausgabe, Status 1); nach dem Lesen des Codes endet er mit `12` (kein gefahrener Lauf des Launchers). Er steht als Grenze in der Spezifikation.
3. **Eine Rückfrage im Container ändert die Zusage „nur die Probe und die Sitzung“.** — **Ausgang:** offen bis zur Closure; der Gegenstand entfällt: Es gibt keine Rückfrage, die Zusage aus [`ADR-0034`](../../adr/0034-devcontainer-shell-ende-des-clients-und-keine-login-shell.md) gilt unverändert, die Folge-ADR `ADR-0039` ist `Deprecated`.
4. **Die Spezifikation reicht weiter als der Test** (Engines, Clients, Hosts). — **Ausgang:** offen bis zur Closure ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)); Vorsorge: Messbedingungen im Satz (Engine und Client 29.9.2, Linux `amd64`, Datum); der Absatz zur Sitzung, die die Zeile selbst schreibt, kennzeichnet die Einstufung als Lesung des Codes.
5. **Restfall:** Schreibt die Sitzung die Zeile des fehlenden Benutzers selbst (Name des Benutzers der Sitzung, Status 1, leere Standardausgabe, kein Terminal), endet der Befehl nach dem Lesen des Codes mit `12` trotz begonnener Shell; das widerspricht [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in diesem einen Fall. [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) nennt ihn als Grenze (Grenzen der Erkennung, Spalte `0`). — **Ausgang:** weiter offen bis zur Closure; entschieden am 2026-11-01 (Projektinhaber): Die Grenze bleibt in der Spezifikation, das Lastenheft bleibt unverändert, eine Lastenheft-CR (Weg B) wird nicht gestellt.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/adapter/driven` (ADP), `tools/` (WZG, Messskript) sowie `spec/`, `docs/plan/` und `docs/messungen/` (DOK); alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-22; Zähler = Dateien in `evidence/`): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (16, offen, über der Schwelle) — Risiko 4; [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (offen, über der Schwelle) — Nachlauf-Review in der DoD; [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert) — Beleg-Tabelle. Evidenz-/Diskrepanz-Risiko: niedrig bis mittel (ein abgeleiteter Randfall).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
