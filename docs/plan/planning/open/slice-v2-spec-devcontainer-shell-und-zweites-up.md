# Slice v2: Spezifikation und Sicht zu `u-boot devcontainer shell` und zum erneuten `up` nach dem Zustand des Containers (nach dem Code, mit Beleg-Tabelle)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-shell-und-zweites-up`](../welle-devcontainer-shell-und-zweites-up.md) (Release `v0.14.0`). Fünfter von sechs Slices; nach allen Code- und CR-Slices der Welle (zuletzt [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md)), vor dem Release-Cut. Muster [`slice-v2-spec-devcontainer-start`](../done/slice-v2-spec-devcontainer-start.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (jeweils in der Fassung des CR dieser Welle); [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) (*Schärft:* [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf)), [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (*Schärft:* [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet), `architecture.md` §2.2, §2.4, §2.6).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet); gelesen und voraussichtlich unverändert [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json); `architecture.md` §2.2, §2.3, §2.4, §2.5, §2.6 (§2.3 und §2.5 sind Funde der Vorab-Suche, nicht im *Schärft* der ADRs). Der Verweis zeigt aufwärts: Spezifikation und Sicht nennen diesen Slice nie.

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Spezifikation und Sicht beschreiben, was Lastenheft und Code nach dieser Welle tragen: das Urteil von `up` nach dem Zustand des Containers mit der Erfolgszeile für einen schon laufenden Container und dem Feld `alreadyRunning`, und den Befehl `u-boot devcontainer shell` (Startphase, Sitzung, Umgebung, Ausgabe auf der Fehlerausgabe, `--quiet`, Exit-Codes, Grenzen). Jede Bedingung steht mit dem Test, der sie pinnt, in einer Beleg-Tabelle oder ist gestrichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Geschrieben wird **nach** dem Code: Die Commit-Reihenfolge belegt es.

**Fundstellen aus dem Gelesenen (2026-10-21; Zeilen der Fassung bei Planung, der Implementer liest sie am Ist-Stand nach).**

| Stelle | Heutiger Satz | Auftrag |
|---|---|---|
| [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Titel und *Aufruf* (Zeilen 739 und 743) | „`u-boot devcontainer up` und `down`“; Exit-Code `2` für beide Unterbefehle | `shell` als **neuer Unterabschnitt** ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) *Schärft*); der Titel bleibt, weil sein Anker in der Sicht (§2.2, §2.4), in der Spezifikation und in zwei `Accepted`-ADRs steht (Suche nach dem Anker: [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md), [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md), `spec/architecture.md`, `spec/spezifikation.md`) und eine Link-Korrektur in einer `Accepted`-ADR eine Entscheidung des Projektinhabers braucht; Abweisung von `--json` bei `shell` im *Aufruf* |
| dieselbe, *Fallbacks und Abbruch eines Startscripts* (Zeilen 775 und 777) | „Jeder andere Zustand bleibt still“; „Zugesagt ist die Abbildung für ein `failed`, das das Startscript dieses Starts schreibt … welcher Exit-Code für ein `failed` aus einem früheren Lauf gilt, ist nicht zugesagt“ | Urteil nach dem Zustand in beiden Pfaden ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Entscheidung 1), Meldung nach einem gelungenen Lauf, Grenze (ii) |
| dieselbe, *Ausweisung des Sockets* (Zeile 779) | fünf Teilaussagen, (2) „für die Dauer des Aufrufs“ | Variante „für die Dauer des Starts“ bei `shell` ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2) |
| dieselbe, *Ausgabe von `up`* (Zeilen 781 bis 785) | nur `Devcontainer started: …`; `data` mit drei Feldern | Zeile `Devcontainer already running: …`, Feld `alreadyRunning`, die vorgezogene Abfrage und ihre Grenze ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Entscheidung 2) |
| dieselbe, Tabelle *Exit-Codes* (Zeilen 789 bis 799) | Spalten `up` und `down`; Zeile `11` „Abbruch eines Startscripts mit `failed` in der Statusdatei“ | Zeile `11` nach dem Zustand; Spalte `shell` (`0` nach Beginn der Sitzung, `12` bei gescheiterter Probe) |
| dieselbe, *Reichweite* (Zeile 801) | „beide Unterbefehle“ | drei Unterbefehle |
| [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), *Ermittlung durch den Host* (Zeile 653) | „`devcontainer up` kennt die Container-Erkennung nicht und fragt ebenfalls immer über die Engine“ | gilt auch für `shell` (prüfen am Code) |
| dieselbe, *Abbildung in `generate devcontainer` und `init`* (Zeile 671) | „`devcontainer up` fragt … nach einem gelungenen Start für `result=fallback` (Warnung) und nach einem fehlgeschlagenen für `result=failed` (Exit-Code `11` …)“ | `failed` in beiden Pfaden; `shell` ebenso |
| [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf), *Abbruch* (Zeile 735) | „`u-boot devcontainer up` bildet es ebenfalls auf Exit-Code `11` ab, wenn …“ | auch nach einem gelungenen Lauf der CLI und bei `shell` |
| [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) (Zeilen 584, 588, 591 und 595) | Liste der ändernden Subkommandos mit `devcontainer up` und `devcontainer down`; Ausnahme 1; `DegradationFallback` auch bei `devcontainer up`; „`devcontainer up --quiet` hält den Fortschritt zurück …; `devcontainer down --quiet` gibt nichts aus“ | `shell` nach dem freigegebenen Wortlaut von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Stellen 6a und 6b des CR); die Spezifikation erweitert das Lastenheft nicht |
| [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile 830, Wertebereich von `command`) | `devcontainer` enthalten; `subcommand` ohne Wertebereich | voraussichtlich unverändert (prüfen am `allOf`) |
| `architecture.md` §2.2 (Zeilen 128 bis 140) | `DevcontainerStartService` startet und stoppt; „entscheidet nach einem Fehler des Starts … ob ein Startscript abgebrochen hat“ | Urteil nach dem Zustand; `shell` |
| `architecture.md` §2.3 (Zeile 176) | „`DevcontainerUseCase` (`Up`/`Down`)“ | Methode für `shell` nach dem Code (Fund der Vorab-Suche) |
| `architecture.md` §2.4 (Zeile 232) | „`DevcontainerLauncher` (`Up`/`Stop`)“ | `Shell` |
| `architecture.md` §2.5 (Zeilen 264 bis 270) | „`devcontainer` (mit `up`/`down`)“; Abweisung von `--dry-run`/`--diff` durch `up` und `down` | `shell` (Fund der Vorab-Suche) |
| `architecture.md` §2.6 (Zeile 301) | `DevcontainerLauncher` „für Start und Stopp des Devcontainers“ | Sitzung über `docker exec` des Hosts |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code, Tests, Handbuch, READMEs, `CHANGELOG.md`** — die Code-Slices davor; dieser Slice schreibt nur `spec/spezifikation.md` und `spec/architecture.md`. Findet er eine Bedingung ohne Test, wird sie gestrichen oder als „nicht zugesagt“ geführt, nicht durch einen Test hier nachgeholt (ein Test ist Code; Folge-Slice, wenn der Projektinhaber die Zusage will).
- **Das Lastenheft** — der CR dieser Welle; die Spezifikation präzisiert und erweitert nicht ([`AGENTS.md`](../../../../AGENTS.md) §2).
- **Aussagen zu Colima, `updateRemoteUserUID` und zum Port-Schritt `Preflight`** — offene Entscheidungen (b), (c) und der Folge-Slice [`slice-v1-devcontainer-start-colima-messung`](slice-v1-devcontainer-start-colima-messung.md); dieser Slice ändert die heutigen Sätze dazu nicht.
- **ADR-Bezüge in Spezifikation und Sicht** — die Straten verweisen nicht abwärts ([`AGENTS.md`](../../../../AGENTS.md) §3.4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Spezifikation** nach den Aufträgen der Fundstellen-Tabelle (§1): [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) mit neuem Unterabschnitt für `shell` und dem Urteil nach dem Zustand, [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf) und [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet); nichts darin weiter als das Lastenheft in der Fassung des CR.
- [ ] **Sicht** `spec/architecture.md` §2.2 bis §2.6 nach dem Code; sprach- und meilensteinfrei, ohne ADR-, Slice- und Roadmap-Bezug ([`AGENTS.md`](../../../../AGENTS.md) §3.4).
- [ ] **Beleg-Tabelle** (Bedingung → Test oder Streichung) für jede Bedingung, die der Diff schreibt oder deren Zeile er berührt (Kandidat der Regel aus [`slice-cr-json-subcommand-pflicht`](../done/slice-cr-json-subcommand-pflicht.md) §7: die Pin-Pflicht gilt für die ganze berührte Zeile, sonst nennt der Plan sie als ungepinnten Bestand); jeder genannte Testname existiert im Baum (Suche je Name, V-2 der Vorwelle fand drei falsche Namen). *Beleg:* Tabelle in der Übergabe, Verification gegen den Baum.
- [ ] `make gates` grün, `make docs-check` 0 Befunde, `make doc-trace` ohne neue Waise.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei Start und zurück auf „Nichts in Arbeit.“ bei der Closure; Handbuch, READMEs und `CHANGELOG.md` unverändert (kein Verhalten); ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Vor dem Schreiben sucht der Implementer in beiden Dateien nach `devcontainer up`, `devcontainer down`, `up`/`down`, `Up`/`Down`, `Up`/`Stop`, `Hilfs-Container`, `Dauer des Aufrufs` und den Werten von `subcommand` und trägt jeden Treffer mit Befund in die Fundstellen-Tabelle ein, auch den unveränderten ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Fundstellen-Tabelle in §1 (gelesen am 2026-10-21, Zeilen 584 bis 801) |
| `spec/architecture.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | §2.2, §2.3, §2.4, §2.5, §2.6 (gelesen, Zeilen 128 bis 140, 176, 232, 264 bis 270, 301) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start, „Nichts in Arbeit.“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Zeile dieses Slice |
| diese Datei, die Welle-Datei, die übrigen Slices der Welle | update | Beleg-Tabelle, Plan-Abgleich und Closure; Verweis-Nachzug durch `make slice-mv` |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence je angefallener Beobachtung |
| `spec/lastenheft.md`, `docs/plan/adr/`, `internal/`, `cmd/`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | Lastenheft im CR, ADRs `Accepted` und unveränderlich ([`AGENTS.md`](../../../../AGENTS.md) §3.5), Code und Nutzer-Dokumente in den Code-Slices |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-cr-devcontainer-shell-und-zweites-up`](slice-cr-devcontainer-shell-und-zweites-up.md), [`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md), [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md) und [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md) liegen in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Vorab-Suche findet Stellen außerhalb der genannten Abschnitte, die eine eigene Entscheidung brauchen (etwa eine Umbenennung mit Anker-Wechsel in einer `Accepted`-ADR).
- `in-progress` → `open` (blockiert): eine Bedingung, die Lastenheft oder Handbuch zusagen, hat keinen Test; dann zuerst ein Test-Slice oder ein Folge-CR, nicht ein Satz ohne Beleg.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; die Spezifikations-Commits liegen nach dem Move des Befehl-Slice nach `done/` (`git log` belegt die Reihenfolge, `git log <Move>..HEAD -- internal cmd` ist leer); Review- und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine Bedingung der Spezifikation ohne Test** (etwa zum Terminal-Pfad, zu Strg-C ohne Terminal, zur vorgezogenen Abfrage bei nicht antwortender Engine). — **Ausgang:** weiter offen; Beleg-Tabelle; was kein Test trägt, ist gestrichen oder „nicht zugesagt“ ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Die Vorab-Suche übersieht eine Aufzählung** (in der Vorwelle übersah sie die Pflicht von `subcommand` und einen Satz in [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung)). — **Ausgang:** weiter offen; Suche nach §3 vor und nach der Fix-Runde ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), 3×).
- **Die Fix-Runde ändert Spezifikationstext ohne Nachlauf-Review.** — **Ausgang:** weiter offen; Nachlauf-Review über genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), 7×).
- **Die Spezifikation reicht weiter als der Test** (eine Mengenaussage über Hosts, Images oder Läufe). — **Ausgang:** weiter offen ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 10×).
- **Ein geänderter Titel bricht Anker in `Accepted`-ADRs.** — **Ausgang:** weiter offen; der Titel von [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) bleibt (§1).
- **Eine Aussage im Plan über den Text stimmt nicht.** — **Ausgang:** weiter offen ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert).

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

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6. Die drei Paarungen prüft die Closure der Welle.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (Spezifikation und Sicht) und `docs/plan/`, beide in der Sub-Area mit dem Kürzel DOK laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) (Greenfield). Die CLI-Sub-Area bleibt Greenfield, weil ihre neuen Konventionen (Urteil nach dem Zustand, Ausgabe von `shell` auf der Fehlerausgabe, Terminal-Erkennung) hier in die Spezifikation kommen.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-21; Zähler = Dateien in `evidence/`): [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (14, verkörpert) und [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (35, verkörpert) — Beleg-Tabelle und Fundstellen; [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (3, offen), [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (7, offen) und [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (10, offen) — Risiken in §6; [`BEO-DOK/doc-immutable-greift-nicht-bei-inline-status`](../observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md) (1) — die ADRs 0028 bis 0031 tragen den Status inline, `make doc-immutable` sieht sie nicht; dieser Slice ändert keine ADR, die Welle-Closure prüft sie zusätzlich mit `git diff`. Evidenz-/Diskrepanz-Risiko: mittel.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
