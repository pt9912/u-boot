# Slice v1: `devcontainer up` urteilt nach dem Zustand des Containers — `failed` ergibt Exit-Code 11 auch nach einem gelungenen Lauf der CLI, die Erfolgszeile nennt einen schon laufenden Container (`alreadyRunning`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-shell-und-zweites-up`](../welle-devcontainer-shell-und-zweites-up.md) (Release `v0.14.0`). Zweiter von sechs Slices; nach [`slice-cr-devcontainer-shell-und-zweites-up`](../done/slice-cr-devcontainer-shell-und-zweites-up.md), vor [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md). Der Befehl `shell` setzt auf das Urteil nach dem Zustand auf ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1, *Startphase*).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Fehlerklassen in der Fassung des CR), [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Strenge `onUnavailable: fail`), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (Abbruch des Egress-Startscripts), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Klasse 11), [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) („Weitere Felder sind erlaubt“: `alreadyRunning`), [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Umschlag); [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) (Entscheidungen 1 bis 3, `Accepted`), [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (Entscheidung 1 bis auf die zwei abgelösten Stellen; Grenze (ii) gilt fort).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (*Fallbacks und Abbruch eines Startscripts*, *Ausgabe von `up`*, Tabelle *Exit-Codes*), [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung), [`SPEC-048`](../../../../spec/spezifikation.md#spec-048--startscript-der-egress-restriktion-upstream-dns-abbruch-und-wiederholungslauf). Der Slice **liest** sie und ändert keine; die Nachführung ist [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md). Der Verweis zeigt aufwärts.

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-22); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `u-boot devcontainer up` urteilt nach dem Zustand des Containers ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md)): Nach jedem Lauf der CLI, der nicht an der Vorab-Prüfung scheitert und nicht abgebrochen wird, ergibt `result=failed` in der Statusdatei einer konfigurierten Fähigkeit Exit-Code `11`, auch wenn die CLI `success` meldet (Entscheidung 1); die Erfolgszeile sagt `Devcontainer already running: …` und der JSON-Umschlag `alreadyRunning: true`, wenn der Container vor dem Aufruf lief und die CLI dieselbe Kennung meldet (Entscheidung 2); `u-boot` merkt sich nichts und heilt nichts selbst (Entscheidung 3). Dazu die Docker-Tests der Fitness Function mit Wächter und die Dokumente.

**Herkunft.** Lücke L-2 des [Validator-Berichts](../../../reviews/2026-10-21-validation-welle-devcontainer-start-ohne-editor.md) (nach Exit 11, Exit 12 oder Strg-C meldet ein zweites `up` Exit 0 und `Devcontainer started`); offene Entscheidung (a) der Closure von [`welle-devcontainer-start-ohne-editor`](../done/welle-devcontainer-start-ohne-editor.md) §7; entschieden in [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) (Option C, Feld `alreadyRunning`, Lastenheft-CR).

**Gelesener Stand (2026-10-21; die Zeilen von `devcontainer_start.go` am 2026-10-22 nachgelesen).** `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Up` (Zeilen 111 bis 149): Statusdateien liest der Dienst nur bei `driven.ErrDevcontainerRuntime` (`startScriptFailed`, Zeilen 205 bis 212); im Erfolgsfall nur `fallbackWarnings` (Zeile 147), das `failed` schweigend übergeht (`internal/hexagon/application/devcontainer_fallback.go`, `fallbackWarnings`, Zeilen 276 bis 298). Vor dem Aufruf des Launchers fragt `Up` den Inspector nicht; danach fragt ihn der `fallbackReader` einmal je Lauf, und nur, wenn die `u-boot.yaml` eine Fähigkeit konfiguriert (`devcontainer_fallback.go`, `fallbackReader.read`, Zeilen 149 bis 157). `Down` fragt ihn vor `Stop` (Zeile 219). Der Unit-Test `TestDevcontainerUp_FallbackWarningsAfterStart` pinnt das Schweigen (Fall „failed is silent after a success“, `devcontainer_start_test.go` Zeile 331). Die Erfolgszeile steht in `runDevcontainerUp` (`internal/adapter/driving/cli/devcontainer.go`, Zeilen 270 bis 273), das `data`-Objekt in `devcontainerUpData` (Zeilen 48 bis 52). Beide Klassifikatoren tragen `driving.ErrDevcontainerStartScriptFailed` schon (`isEnvironmentError` in `cli.go` Zeile 339, `mapDevcontainerEngineSentinel` in `devcontainer.go` Zeilen 340 bis 349).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Befehl `shell`** — [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md) und [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md); dieser Slice liefert das Urteil nach dem Zustand, auf das `shell` aufsetzt.
- **`spec/`** — das Lastenheft ändert der CR davor, Spezifikation und Sicht der Spezifikations-Slice danach ([`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md)); jede Bedingung im Handbuch braucht einen Test aus diesem Slice oder entfällt ([`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Selbstheilung, Zustandsspeicher und Frische-Prüfung** — verworfen ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Optionen D und E, Entscheidung 3); ein Abbruch ohne Statusdatei (etwa der Clone im `postCreateCommand`) bleibt für den Exit-Code unsichtbar, Erfolgszeile und Handbuch tragen ihn.
- **Ein neuer Sentinel oder Port** — die ADR nennt keinen: `ErrDevcontainerStartScriptFailed` und `DevcontainerInspector.FindContainer` bestehen.
- **Ein Port-Schritt `Preflight` vor der Abfrage** (offene Entscheidung (b) der Vorwelle, L-6) — die vorgezogene Abfrage verzögert bei nicht antwortender Engine den Exit-Code `11` um die Zeitgrenze des Inspectors; das ist ein akzeptiertes Negativ der ADR, kein Auftrag hier.
- **`devcontainer down`** — unverändert; es fragt den Inspector schon.
- **Korrekturhinweise nach [`LH-FA-CLI-004`](../../../../spec/lastenheft.md#lh-fa-cli-004--fehlerausgabe) für die übrigen Meldungen** (offene Entscheidung (i), L-5) — nur die neue Meldung nennt den Weg `u-boot devcontainer up --recreate`, weil die ADR es verlangt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Verhalten** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) in der Fassung des CR, [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Entscheidungen 1 bis 3): Nach einem gelungenen Lauf der CLI liest `Up` die Statusdateien der konfigurierten Fähigkeiten auf demselben Weg wie im Fehlerpfad; `result=failed` ergibt `driving.ErrDevcontainerStartScriptFailed` (Exit 11 in `ExitCode` und in `mapDevcontainerUpErrorToDiagnostic`, nicht zugleich `driven.ErrDevcontainerRuntime`); die Meldung nennt den Satz der Fähigkeit für `failed` (derselbe wie bei `doctor`), den Zeitpunkt `at`, dass der Container läuft und dieser Aufruf das Startscript nicht ausgeführt hat, und `u-boot devcontainer up --recreate`. Ohne `--recreate` fragt `Up` vor dem Launcher den Inspector einmal nach dem Container; lief er und meldet die CLI dieselbe Kennung, trägt die Antwort `AlreadyRunning`, die Textausgabe `Devcontainer already running: container <zwölf Zeichen>, remote user <Benutzer>, workspace <Pfad>` und `data` das Feld `alreadyRunning` (`true`, sonst `false`; ohne `omitempty`); ein Fehler der Abfrage ist kein Fehler des Befehls. Langtext von `up` (Zeile `11` der Exit-Codes) und die Kommentare von Sentinel und Port beschreiben den Ist-Stand ([`AGENTS.md`](../../../../AGENTS.md) §3.7, §3.9). *Beleg:* Unit-Tests nach der Fitness Function der ADR (§3).
- [ ] **Docker-Tests mit Wächter** (Fitness Function, Build-Tag `docker`): (1) Folge zu S8 — erzeugtes Projekt mit Egress-Restriktion, `onUnavailable: fail`, entzogener Capability und `u-boot.yaml`: erstes `up` `11`, zweites `up` `11` (mit `--json`: `exitCode` `11`); *Wächter:* ohne `u-boot.yaml` endet das zweite `up` mit `0`. (2) Zähler — ein handgeschriebenes `postStartCommand` schreibt je Lauf eine Zeile in eine Datei im Container: zweites `up` gegen den laufenden Container lässt die Zeilenzahl unverändert und meldet `alreadyRunning: true`; nach `down` und `up` steht eine Zeile mehr und `alreadyRunning: false`. Beide sichtbar `PASS` im CI-Lauf von `make test-docker` (`-v`), nicht übersprungen. *Beleg:* Lauf `integration-docker` mit Job-Log.
- [ ] **Doku:** Handbuch 4.9 (*Ausgabe*, Tabelle *Exit-Codes von `devcontainer up`*, *Erneutes `up` nach einem Fehler oder Abbruch*), §7 (*Exit-Codes*, *Sandbox-Container startet nicht* Schritt 3), Anhang *JSON-Ausgabe* (`data` von `devcontainer up`); `docs/maintainer/cli-json-contract.md` §6.11; `CHANGELOG.md` `[Unreleased]` unter *Changed* (ein Skript, das `up` nach einem Fehlschlag wiederholt, bekommt `11` statt `0`). Handbuch ohne Anforderungs-, ADR- und Spec-Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6); jede neue Bedingung im Handbuch mit Test aus diesem Slice oder gestrichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). *Beleg:* Beleg-Tabelle (Satz → Test) in der Übergabe.
- [ ] `make gates` grün; `make test-docker` grün als CI-Lauf auf dem Push-Tip (ein Host-Lauf ersetzt ihn nicht, [`AGENTS.md`](../../../../AGENTS.md) §3.1).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor; ändert eine Fix-Runde den Port-Vertrag oder einen Text, den das Handbuch zusagt, liest ein Nachlauf-Review genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md))
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei Start und zurück auf „Nichts in Arbeit.“ bei der Closure; ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; widerlegt ein Test die Entscheidung 2 (die CLI führt bei `up` gegen einen laufenden Container Lebenszyklus-Befehle aus), steht der Befund in §7 und der Slice geht an den Architect zurück (Re-Evaluierungs-Trigger der ADR).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei oder Stelle, die der Planner nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/application/devcontainer_start.go` | update | `Up` (gelesen, Zeilen 111 bis 149): nach einem gelungenen `launcher.Up` vor `fallbackWarnings` (Zeile 147) die Prüfung `startScriptFailed` (Zeilen 205 bis 212) mit eigener Meldung nach Entscheidung 1; ohne `req.Recreate` vor `s.launcher.Up` (Zeile 129) eine Abfrage `s.inspector.FindContainer(ctx, req.BaseDir)`; `AlreadyRunning`, wenn sie einen laufenden Container mit der Kennung ergab, die `res.ContainerID` meldet; ein Fehler der Abfrage ergibt `false`. Ob die Abfrage vor oder nach der Ausweisung über `req.OnWarning` (Zeilen 123 bis 127) steht, entscheidet der Implementer; die Ausweisung bleibt vor dem Start des Hilfs-Containers. Doc-Kommentar von `startScriptFailed` („The file describes the last run of its script, not necessarily this start“) auf den Ist-Stand |
| `internal/hexagon/application/devcontainer_fallback.go` | update (prüfen) | der Satz der Fähigkeit für `failed` steht je Probe im Feld `failed` (gelesen: `vfsProbe`, `daemonProbe`, `egressProbe`, Zeilen 46 bis 80), der Zeitpunkt über `rd.status.At()` (`reportedSuffix`, Zeile 213); ob die Meldung hier als Funktion neben `fallbackWarnings` (Zeilen 276 bis 298) oder in `devcontainer_start.go` entsteht, entscheidet der Implementer; der Doc-Kommentar von `fallbackWarnings` („Every other outcome (ok, failed, …) is silent“) bleibt wahr, solange `failed` vor dem Aufruf geprüft wird |
| `internal/hexagon/port/driving/devcontainer.go` | update | `DevcontainerUpResponse` (gelesen, Zeilen 33 bis 46) bekommt `AlreadyRunning bool`; der Doc-Kommentar von `ErrDevcontainerStartScriptFailed` („signals that the devcontainer CLI did not succeed and a start script … reported `result=failed`“, Zeilen 86 bis 92) und der Vertrag von `DevcontainerUseCase` („a failed start whose container reports `result=failed`“, Zeilen 104 bis 108) auf den Ist-Stand ([`AGENTS.md`](../../../../AGENTS.md) §3.9 (b)) |
| `internal/hexagon/application/devcontainer_start_test.go` | update | Fall „failed is silent after a success“ in `TestDevcontainerUp_FallbackWarningsAfterStart` (gelesen, Zeile 331) kehrt sich um; neue Fälle nach der Fitness Function: `success` mit `failed` einer konfigurierten Fähigkeit ergibt den Sentinel (nicht zugleich `ErrDevcontainerRuntime`), die Meldung nennt `at` und `--recreate`; `success` mit `fallback` bleibt Warnung; `failed` einer nicht konfigurierten Fähigkeit und ohne `u-boot.yaml` bleibt Erfolg; `AlreadyRunning` wahr bei laufendem Container und gleicher Kennung, falsch bei `--recreate`, vorher gestopptem Container, anderer Kennung und Fehler der Abfrage. Der Fake `fakeInspector` (`devcontainer_fallback_test.go`, gelesen Zeilen 30 bis 66) zählt die Aufrufe über `findFolders` und liefert einen festen Container; ob er je Aufruf andere Antworten braucht (vorher gestoppt, nachher laufend), prüft der Implementer |
| `internal/adapter/driving/cli/devcontainer.go` | update | `devcontainerUpData` (gelesen, Zeilen 48 bis 52) bekommt `AlreadyRunning bool` mit `json:"alreadyRunning"`; `runDevcontainerUp` (Zeilen 270 bis 273) schreibt je nach Antwort die Zeile `Devcontainer started: …` oder `Devcontainer already running: …`; Langtext von `up` (Zeilen 117 bis 131): Absatz zu den Fallbacks und Zeile `11` der Exit-Codes nennen das `failed` aus einem früheren Start |
| `internal/adapter/driving/cli/devcontainer_test.go` | update | `TestExecute_DevcontainerUp_HumanOutput` (Zeile 128), `…_Quiet` (Zeile 157), `…_JSON` (Zeile 193): Zeile und Feld in beiden Fällen; `TestExecute_Devcontainer_SentinelsPerSubcommand` (Zeile 276) deckt den Sentinel schon (Fall nach einem gelungenen Lauf ergänzen, prüfen) |
| `internal/adapter/driving/cli/cli.go` | **unverändert**, prüfen | `isEnvironmentError` (gelesen, Zeilen 333 bis 340) bildet `ErrDevcontainerStartScriptFailed` auf `11` ab; der Kommentar von `ExitCode` (Zeilen 266 bis 269, „a container start script that reported `result=failed`“) gilt ohne Bedingung weiter; der Doc-Kommentar von `isEnvironmentError` nennt als Rang-Zeiger nur [`ADR-0029`](../../adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (Zeile 335) — ob er [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) mitnennt, entscheidet der Implementer ([`AGENTS.md`](../../../../AGENTS.md) §3.7) |
| `internal/e2e/devcontainer_start_docker_test.go` | update | Folge zu `TestDevcontainerUp_Docker_S8_FailurePicture` (gelesen, Zeile 671; der Test ruft heute keinen zweiten Aufruf) und ein neuer Test für den Zähler; Hilfen `dcGenerateProject` (Zeile 191), `dcWithoutNetAdmin` (Zeile 229), `dcHandWritten` (Zeile 247), `dcRun` (Zeile 646, gebautes Binary), `dcRegisterCleanup` (Zeile 165); Budget `dcTestBudget` (Zeile 50) |
| `cmd/uboot/main.go` | **unverändert** | Wiring (gelesen, Zeilen 170 bis 174) bleibt: kein neuer Port; der Kommentar „`devcontainer up|down` … the inspector finds the container and its status files“ bleibt wahr |
| `internal/adapter/driven/docker/inspector.go` | **unverändert** | `inspectorTimeout = 10 * time.Second` (gelesen, Zeile 26) begrenzt die vorgezogene Abfrage; das ist die Verzögerung des akzeptierten Negativs der ADR |
| `docs/user/benutzerhandbuch.md` | update | gelesen: Abschnitt 4.9 *Ausgabe* (Zeilen 804 bis 806, Erfolgszeile), Tabelle *Exit-Codes von `devcontainer up`* (Zeilen 851 bis 859, Zeile `11`), *Erneutes `up` nach einem Fehler oder Abbruch* (Zeilen 874 bis 883: nach Exit 11 „nicht gemessen“), §7 *Exit-Codes* (Zeile 2179), *Sandbox-Container startet nicht* Schritt 3 (Zeilen 2363 bis 2368), Anhang *JSON-Ausgabe* (Zeile 2590, `data` von `devcontainer up`; Besonderheiten Zeilen 2603 bis 2606). Keine Kennungen |
| `docs/maintainer/cli-json-contract.md` | update | §6.11 (gelesen, Zeilen 1244 bis 1302): Beispiel-Umschlag mit `alreadyRunning`, Satz zu `ErrDevcontainerStartScriptFailed` auch nach einem gelungenen Lauf (Zeile 1291) |
| `CHANGELOG.md` | update | `[Unreleased]` (gelesen, Zeile 12, leer): Eintrag unter *Changed* (Verhaltensänderung für Skripte) und die neue Zeile samt Feld |
| `README.md`, `README.de.md` | **unverändert**, prüfen | Befehlstabelle (`README.md` Zeile 224, `README.de.md` Zeile 227) beschreibt `up` ohne Erfolgszeile; Status-Block beschreibt `v0.13.0` (Zeile 11 bzw. 12) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start, „Nichts in Arbeit.“ bei der Closure |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Zeile dieses Slice |
| diese Datei, die Welle-Datei, die übrigen Slices der Welle | update | Plan-Abgleich und Closure; Verweis-Nachzug durch `make slice-mv` |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | Evidence je angefallener Beobachtung |
| `spec/` <!-- d-check:ignore (Pfad-Zitat) --> | **unverändert** | Lastenheft im CR davor, Spezifikation und Sicht im Spezifikations-Slice danach |

### Entscheidungen beim Start (Implementer, 2026-10-22; Zweifelsfälle nach Plan, ADR und Lastenheft)

Der Plan überlässt dem Implementer mehrere Punkte (Reihenfolge der Abfrage, Ort der Meldung, Fake des Inspectors, Rang-Zeiger in `cli.go`); die Punkte 1 bis 3 sind hier entschieden, bevor Code entsteht (4 und 5 legen den Umgang mit Handbuch und Mutationen fest), der Rest steht in *Umsetzung*.

1. **Reihenfolge in `Up`.** Vorab-Prüfung (Konfiguration, Pfad, Namen) → Ausweisung über `OnWarning` → Abfrage `FindContainer` (nur ohne `--recreate`) → Launcher. Grund: Die Ausweisung des Sockets erscheint sofort, auch wenn die Abfrage bei nicht antwortender Engine bis zur Zeitgrenze des Inspectors dauert; die ADR beschreibt die Ausgabe von Exit 11 so, dass die Ausweisung „schon erschienen“ ist, bevor der Fehler steht ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Entscheidung 1, *Ausgabe*). Die Abfrage bleibt vor dem Start des Hilfs-Containers.
2. **Meldung nach einem gelungenen Lauf.** `Up` prüft die konfigurierten Fähigkeiten in der Reihenfolge `vfs`, `daemon`, `egress` und nennt die erste mit `result=failed` (wie `startScriptFailed` im Fehlerpfad); die Meldung trägt den Satz der Fähigkeit aus `fallbackProbe.failed` (derselbe wie bei `doctor`), den Zeitpunkt `at` der Datei, dass der Container läuft und dieser Aufruf das Startscript nicht ausgeführt hat, und `u-boot devcontainer up --recreate`. Sie entsteht in `devcontainer_start.go` neben der Prüfung; `devcontainer_fallback.go` bleibt unverändert (der Doc-Kommentar von `fallbackWarnings`, `failed` sei dort still, gilt weiter, weil `Up` `failed` vor diesem Aufruf prüft und `generate` und `init` es weiter still übergehen).
3. **Ein Reader je Aufruf, nach dem Launcher.** Die vorgezogene Abfrage ist eine eigene `FindContainer`-Abfrage; der `fallbackReader` bleibt wie heute (eine Abfrage nach dem Launcher, nur mit konfigurierter Fähigkeit), weil der Container nach dem Launcher ein anderer sein kann als davor (`--recreate`, Neuanlage nach dem Entfernen).
4. **Zusagen im Handbuch.** Das Handbuch sagt über die CLI gegen einen laufenden Container nur, was der Zähler-Test zeigt (`postStartCommand` läuft nicht erneut, läuft nach `down` und `up` wieder); die Beobachtungen des Validators zu `postCreateCommand` und Exit 12 bleiben als „beobachtet in einem Lauf“ stehen und werden nicht erweitert ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
5. **Mutationen.** Der Plan nennt keine; der Implementer legt sie in *Umsetzung* fest (je Zusage der DoD eine, die ihren Test rot färbt) und fährt sie.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`slice-cr-devcontainer-shell-und-zweites-up`](../done/slice-cr-devcontainer-shell-und-zweites-up.md) liegt in `done/` (das Lastenheft trägt Stelle 1, den Satz zur Klasse `11` beim erneuten `up`); kein anderer Slice in `in-progress/` (WIP-Limit 1). Erfüllt am 2026-10-22: `Verantwortlich:` gesetzt; der CR liegt in `done/` ([`slice-cr-devcontainer-shell-und-zweites-up`](../done/slice-cr-devcontainer-shell-und-zweites-up.md), gelesen) und `spec/lastenheft.md` trägt Version `0.21.0` mit dem Satz zur Klasse `11` beim erneuten `up` (Zeile 497, gelesen) und der Zusage-Menge „ohne diese Datei … ein erneutes `up` meldet einen früheren Abbruch nicht“ (Zeile 498); `in-progress/` trägt `README.md`, `carveouts.md` und `roadmap.md`, keinen Slice. Der Projektinhaber hat angeordnet, ohne Rückfrage weiterzumachen (2026-10-22): Zweifelsfälle entscheidet der Implementer nach Plan, [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) und Lastenheft; die Begründung steht in §3 *Umsetzung*.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Zähler-Test zeigt, dass die CLI bei `up` gegen einen laufenden Container Lebenszyklus-Befehle ausführt; dann fällt der Schluss im Kontext der ADR und Entscheidung 2 ([`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md), Re-Evaluierungs-Trigger), und der Slice geht an den Architect.
- `in-progress` → `open` (blockiert): der CR ist nicht angenommen (Re-Evaluierungs-Trigger der ADR: „der Lastenheft-CR zur Klasse `11` wird nicht angenommen“).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; die zwei Docker-Tests sind im CI-Lauf von `make test-docker` sichtbar bestanden und nicht übersprungen; Review- und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die CLI führt bei `up` gegen einen laufenden Container doch Lebenszyklus-Befehle aus.** Belegt ist das Gegenteil nur an einer Messvariante auf Colima und am Validator-Lauf (je ein Lauf, [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Kontext). — **Ausgang:** weiter offen; der Zähler-Test entscheidet, ein Widerspruch ist Rückführung (§4).
- **Das e2e-Paket überschreitet das Zeitlimit** von `make test-docker` (lokal beobachtet, in der CI rund zwölf Minuten); zwei Tests mehr, die je einen Devcontainer starten. — **Ausgang:** weiter offen; der CI-Lauf ist der Beleg, eine Änderung am Limit ist Gate-Gegenstand ([`AGENTS.md`](../../../../AGENTS.md) §3.6; [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md), 1×).
- **Ein Docker-Test hinterlässt Container oder Images.** — **Ausgang:** weiter offen; die Tests nutzen `dcProjectDir` und `dcRegisterCleanup` ([`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md), 2×).
- **Ein Test ist nur mit lokalem Image-Cache grün.** — **Ausgang:** weiter offen; Beleg ist der CI-Lauf ([`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md), 1×).
- **Verhaltensänderung für Skripte:** ein `up`, das nach einem Fehlschlag wiederholt wird, endet mit `11` statt `0`. — **Ausgang:** weiter offen; `CHANGELOG.md` *Changed* und Handbuch nennen sie und die Abhilfe `--recreate` (für `egress` und `daemon` auch `down` und `up`, [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Konsequenzen).
- **Akzeptierte Negative der ADR** (ein im Container gefälschtes `failed` hält `up` bis `--recreate` auf; die vorgezogene Abfrage verzögert den Exit-Code `11` bei nicht antwortender Engine um bis zu zehn Sekunden). — **Ausgang:** weiter offen; das Handbuch nennt nur, was ein Test trägt.
- **Ein Kommentar beschreibt, was entfernt wurde** („failed was silent“, „no longer“), oder ein Kommentar nennt das alte Verhalten nach dem Umbau. — **Ausgang:** weiter offen; Selbstcheck nach [`AGENTS.md`](../../../../AGENTS.md) §3.9 vor jedem Commit ([`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md), verkörpert).
- **Eine Bedingung im Handbuch ohne Test** („ein erneutes `up` … endet wieder mit Exit-Code 11“, „nach `down` und `up` läuft das `postStartCommand` wieder“). — **Ausgang:** weiter offen; Beleg-Tabelle in der Übergabe ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert).
- **Eine Aussage im Plan über den Code stimmt nicht** (die Zeilen „prüfen“). — **Ausgang:** weiter offen ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT, driving), `internal/adapter/driving/cli` (CLI), `internal/e2e` (E2E) und `docs/` (DOK), laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) alle Greenfield; `cmd/uboot` bleibt unverändert (keine Sub-Area). Die CLI-Sub-Area bleibt Greenfield, solange die neue Regel (Urteil nach dem Zustand, Erfolgszeile) in die Spezifikation kommt und nicht nur im Code steht — das trägt der Spezifikations-Slice.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-21; Zähler = Dateien in `evidence/`): [`BEO-E2E/e2e-paket-am-zeitlimit`](../observations/BEO-E2E/e2e-paket-am-zeitlimit/observation.md) (1), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (2), [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (1) — zwei neue Docker-Tests, Risiken in §6; [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) (5, verkörpert) — ein Test kehrt sich um, Kommentare beschreiben den Ist-Stand; [`BEO-CLI/kennung-in-cli-ausgabe`](../observations/BEO-CLI/kennung-in-cli-ausgabe/observation.md) (2) — der Langtext von `up` trägt die Zeile „Exit codes (…)“ mit Kennung (Bestand, offene Entscheidung (j) der Vorwelle; der Slice ändert die Form nicht); [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (35) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (14), beide verkörpert — Prüfpflicht; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (10, offen) — das Handbuch sagt zum Verhalten der CLI nur, was der Zähler-Test zeigt; [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md) (7, offen) — Nachweise nur über `make`. Evidenz-/Diskrepanz-Risiko: mittel (Verhalten der CLI gegen einen laufenden Container nur an zwei Läufen belegt).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
