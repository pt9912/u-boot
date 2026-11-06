# Verification-Report (zweiter Teil, Nachlauf): slice-v1-tap-smoke-runner-macos-15 — 2026-11-06

**Rolle:** Verifier (frischer Kontext; Baseline-Regelwerk `modul-11`), nur lesend

**Gegenstand:** [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md), die Liefer-Punkte **Tap** und **Smoke** (Plan §2). Der erste Teil ([`verification-2026-11-06-slice-v1-tap-smoke-runner-macos-15`](verification-2026-11-06-slice-v1-tap-smoke-runner-macos-15.md)) hat Quelle, Gates und Plan-Tabelle bis `3e1d81a9` geprüft; sein Abschnitt *Übergabe → An den zweiten, lesenden Teil* ist der Prüfauftrag dieses Berichts. Lokaler Tip beim Nachlauf: `d665e62b` (Report des ersten Teils); der Arbeitsbaum war vor und nach allen Lesungen sauber, `git diff d665e62b HEAD` leer.

**Umfang.** Gelesen wurde ausschließlich per `gh api` (GET) gegen `pt9912/homebrew-u-boot` und per lokalem `git`/`cmp`. Geschrieben habe ich in den Tap und auf GitHub nichts: kein Push, kein Dispatch, kein Cancel, kein Rerun. Plan und Quelle sind unverändert. Die Zeitangaben in den Metadaten sind UTC.

**Was der Hauptagent berichtet hat und was ich gelesen habe.** Berichtet: Tap-Commit `662c5fb5`, Lauf `42212837961`, Jobs `127875452461` und `127875452090`, Vorläufer `42196563230` als `queued`. Alles davon habe ich unten selbst gelesen; die Spalte *Beleg* nennt die Abfrage. Nichts wird aus dem Bericht übernommen.

**Eingangs-Kontext:**

- Slice-Plan §1 (Planentscheidung `macos-15`), §2 (DoD „Tap“, „Smoke“), §4 bis §7
- Bericht des ersten Teils (Findings V-1 bis V-5, Übergabe an den Nachlauf Punkte 1 bis 6)
- [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)
- [`harness/verification.md`](../../harness/verification.md), [`harness/roles.md`](../../harness/roles.md), [`AGENTS.md`](../../AGENTS.md)
- [`BEO-DOK/ci-log-falsch-gelesen`](../plan/planning/observations/BEO-DOK/ci-log-falsch-gelesen/observation.md), [`BEO-DOK/zusage-weiter-als-messung`](../plan/planning/observations/BEO-DOK/zusage-weiter-als-messung/observation.md)

---

## Verification Evidence

Scope:
- Slice: [`slice-v1-tap-smoke-runner-macos-15`](../plan/planning/in-progress/slice-v1-tap-smoke-runner-macos-15.md)
- IDs: [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung), [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md)
- Artefakte: `.github/workflows/smoke.yml` im Tap-Repo `pt9912/homebrew-u-boot` (Commit `662c5fb5`), Quelle `packaging/homebrew-tap/.github/workflows/smoke.yml`, Lauf `smoke` Run `42212837961`

### DoD-Abgleich (Plan §2, Punkte „Tap“ und „Smoke“)

Die Häkchen im Slice stehen weiter offen; sie gehören dem Planner.

| # | DoD-Punkt | Zustand | Evidence |
|---|---|---|---|
| 2 | **Tap:** Datei im Tap byte-gleich mit der Quelle (`cmp`, Exit 0), Tap-Commit in §7 | **erbracht, eigene Lesung** | siehe *Tap* unten: `cmp` Exit 0; Blob-SHA beider Seiten `db62a07fb710fb2752857b76bcb09422806a2800`; Tap-Commit `662c5fb5`. Dass der Tap-Commit in §7 steht, ist Sache der Closure (§7 trägt den Stand „Tap-Push … stehen aus“) |
| 3 | **Smoke:** Lauf `smoke` im Tap per `workflow_dispatch`, Lauf `success`, `install (macos-15)` und `install (macos-15-intel)` `success`; arm64-Beleg aus Job-Metadaten und „Set up job“; Formel-Version genannt; Run-ID und Job-IDs in §7 | **erbracht, eigene Lesung**; Run-ID und Job-IDs stehen noch nicht in §7 (Closure) | siehe *Smoke* unten |

Die übrigen Punkte (4 bis 7 aus dem ersten Teil) ändern sich durch den Nachlauf nicht; 8 bis 11 (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) bleiben **offen — Closure**. Dieser Bericht trägt keinen Befund über sie.

### Tap (Auftragspunkte 1 und Push-Ausgang)

| Aussage | Beleg (lesend, 2026-11-06) |
|---|---|
| Datei im Tap = Quelle | `gh api repos/pt9912/homebrew-u-boot/contents/.github/workflows/smoke.yml --jq .content \| base64 -d` gegen `packaging/homebrew-tap/.github/workflows/smoke.yml`: `cmp` **Exit 0** |
| Gleicher Blob | `git hash-object` der Quelle: `db62a07fb710fb2752857b76bcb09422806a2800`; `contents`-API `.sha` im Tap: derselbe Wert. Er ist auch der Blob, den der erste Teil für den Tip `3e1d81a9` festgehalten hat; die Quelle hat sich seither nicht geändert |
| Tap-Kopf | `commits/main` → `662c5fb5c57fc61751f5cff127f213be57573211`, Eltern `b9f23331a5…` (der Kopf, den der erste Teil vor dem Push gelesen hat), Zeitstempel 2026-11-06T09:38:43Z |
| Inhalt des Tap-Commits | `commits/662c5fb5` → genau eine Datei: `.github/workflows/smoke.yml`, `modified`, `+4 −1`. Der Hunk-Rumpf des Patches (`Accept: application/vnd.github.patch`) ist zeilengleich mit dem Hunk von `737dae18` in der Quelle (`diff` ohne Befund) |
| Push angenommen | der Commit liegt auf `main`; die Messstelle von V-3 (erster Teil: Recht für Workflow-Dateien) ist damit belegt: GitHub hat den Push einer Datei unter `.github/workflows/` angenommen. Mit welchem Token, ist nicht gelesen und für den Plan unerheblich |

### Smoke (Auftragspunkte 2 bis 6)

**Lauf** (`actions/runs/42212837961`):

| Feld | Wert |
|---|---|
| `name` / `path` | `smoke` / `.github/workflows/smoke.yml` |
| `event` | `workflow_dispatch` (der Push-Trigger beobachtet nur `Formula/u-boot.rb`; der Dispatch ist der im Plan vorgesehene Weg) |
| `head_sha` | `662c5fb5c57fc61751f5cff127f213be57573211` — der Tap-Kopf, also die Datei mit `macos-15` |
| `status` / `conclusion` | `completed` / **`success`**, `run_attempt` 1 (kein Rerun) |
| Zeiten | `created_at` 11:47:29Z, `updated_at` 11:49:41Z |

**Jobs** (`actions/runs/42212837961/jobs`):

| Job | ID | `labels` | `runner_name` | `conclusion` | Schritte |
|---|---|---|---|---|---|
| `install (macos-15)` | `127875452461` | `["macos-15"]` | `GitHub Actions 1000036052` (nicht leer) | **success** | alle 9 `success`, darunter „Formel pruefen (gatend)“ |
| `install (macos-15-intel)` | `127875452090` | `["macos-15-intel"]` | `GitHub Actions 1000036051` (nicht leer) | **success** | alle 9 `success`, darunter „Formel pruefen (gatend)“ |

Beide Jobs tragen `head_sha` `662c5fb5…`. Die Gate-Schritte sind nicht übersprungen und nicht lockerer geworden: der Schritt „Formel pruefen (gatend)“ (`brew audit --strict --online pt9912/u-boot/u-boot`, `exit 1` bei Befund) steht in der Datei unverändert und ist in beiden Läufen `success`.

**Job-Log, Abschnitt „Runner Image“** (`gh api --allow-escape-sequences repos/pt9912/homebrew-u-boot/actions/jobs/<id>/logs`, nicht `gh run view --log`; die ANSI-Folgen habe ich vor dem Lesen entfernt — der erste Abruf ohne das Flag wurde von `gh` abgelehnt und hat nichts geliefert):

| Job | OS laut Log | `Image:` | Image-Version | `brew`-Präfix im Log |
|---|---|---|---|---|
| `install (macos-15)` | macOS 15.8.2 (24H232) | `macos-15-arm64` | 20261028.0489.1 | `/opt/homebrew/…` (`Cloning into '/opt/homebrew/Library/Taps/…'`, `/opt/homebrew/Cellar/u-boot/0.24.0`); Gems `rubydex 0.4.1 (arm64-darwin)` |
| `install (macos-15-intel)` | macOS 15.7.9 (24G830) | `macos-15` | 20260824.0482.1 | `/usr/local/…` (`/usr/local/Cellar/u-boot/0.24.0`, `/usr/local/Homebrew/…`); Gems `rubydex 0.3.0 (x86_64-darwin)` |

Die Aussage **„arm64“ für den Job `install (macos-15)` trägt das Log**: Das Feld `Image:` nennt `macos-15-arm64`, das Homebrew-Präfix ist `/opt/homebrew`, die Gems sind `arm64-darwin`. Eine Zeile mit `uname -m` oder `file` auf das Binary steht im Log nicht (der Workflow führt sie nicht aus); die Aussage stützt sich auf die drei genannten Stellen, nicht auf ein Architektur-Kommando. Der Intel-Job ist entsprechend `x86_64`. Beide Jobs laufen auf macOS 15 (Patch-Stand verschieden, 15.8.2 gegen 15.7.9); das trägt die Plan-Begründung „beide Jobs laufen auf derselben Hauptversion“.

**Formel-Version (Schritt „Version“ und „Tap und Install“):**

| Job | `Formula u-boot (…)` | `Cellar/u-boot/…` | `u-boot --version` |
|---|---|---|---|
| `install (macos-15)` | `0.24.0` | `0.24.0` | `u-boot version 0.24.0` |
| `install (macos-15-intel)` | `0.24.0` | `0.24.0` | `u-boot version 0.24.0` |

`Formula/u-boot.rb` im Tap nennt auf `main` die URLs von `v0.24.0`, und beide Logs nennen `Formula u-boot (0.24.0)`; „0.24.0“ ist also die Version, die diese Läufe installiert haben (Plan §2 Punkt „Smoke“, letzter Satz). **Die Aussage „grün“ gilt für:** Formel `u-boot 0.24.0`, Tap-Kopf `662c5fb5`, die zwei Images aus der Tabelle oben, die neun Schritte der Datei (Tap und Install, Formel-`test`-Block, Version, `init`, `add postgres`, `doctor` mit Exit 0 oder 11, gatender `brew audit --strict --online`). Nicht gelesen und nicht behauptet: andere macOS-Versionen, andere Formel-Versionen, ein Lauf mit Docker auf dem Runner.

**Schritt `doctor`** (Plan-Schritt „Docker fehlt auf dem Runner, Exit 11 erwartet“): im Intel-Log gelesen: drei `✗` (`docker.installed`, `docker.compose.installed`, `docker.reachable`), `Summary: 3 error, 0 warn, 14 ok`; der Schritt ist `success`, weil die Datei Exit 0 oder 11 zulässt. Das ist der erwartete Zustand, kein Befund. Die entsprechende Zeile im arm64-Log habe ich nicht einzeln gelesen; der Schritt ist dort `success` (Job-Metadaten).

**Annotationen** (`check-runs/<id>/annotations`):

| Job | Annotationen |
|---|---|
| `install (macos-15)` | genau eine, Stufe `notice`: „Due to capacity constraints, jobs targeting macOS arm64 runners may experience longer queue times.“ Keine `warning`, kein `failure`; **keine** Meldung der Art „configuration … is not supported“ |
| `install (macos-15-intel)` | drei: zwei `warning` (AWS SAM CLI, Homebrew-Installer, 2023) und eine Stufe `failure`: „Failed to import: /usr/local/Homebrew/Library/Taps/hashicorp/homebrew-tap/Formula/vagrant.rb“. Der Job ist trotzdem `success`; die Meldung stammt aus dem Schritt „Tap und Install“ (im Log als `##[error]` am Beginn des Schritts), nicht aus einem Gate |

Die `failure`-Stufe auf einem `success`-Job ist der Hinweis „Failed to import … vagrant.rb“, den Plan §6 (Risiko „Ein neuer Befund unter `macos-15`“) als Angabe des Projektinhabers im Cut-Slice [`slice-v1-release-cut-v0.24.0`](../plan/planning/done/slice-v1-release-cut-v0.24.0.md) nennt und der Planner nicht selbst gelesen hat. **Jetzt gelesen**, am Intel-Job, dessen Label sich nicht geändert hat. Er betrifft einen anderen Tap im Homebrew-Verzeichnis des Runners, nicht die Formel dieses Slice (siehe N-3).

### Der frühere Dispatch-Lauf `42196563230` (Auftragspunkt: kein Beleg)

| Feld | Befund (`actions/runs/42196563230`, gelesen 2026-11-06 ~11:52Z) |
|---|---|
| `event` / `head_sha` | `workflow_dispatch` / `662c5fb5…` (derselbe Kopf) |
| `created_at` | 2026-11-06T09:39:00Z |
| `status` / `conclusion` | `queued` / leer; `…/jobs` liefert `total_count` 0 |

Der Lauf trägt keinen Job und hat keinen Runner bekommen; er ist **kein Beleg** für „Smoke“ und wird hier nicht gezählt. Die Zahl der Läufe auf demselben Kopf ist zwei; der grüne Beleg ist ausschließlich `42212837961`. Warum `42196563230` keine Jobs hat, ist nicht gelesen (siehe N-2).

### Entscheidungs-Konformität

| Quelle | Ergebnis | Beleg |
|---|---|---|
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 (Smoke im Tap) | **konform** | der Tap-Commit ändert nur das Label der Matrix und den Kommentar (Patch gelesen); Prüfungen, Trigger und Schritte der Datei unverändert |
| [`AGENTS.md`](../../AGENTS.md) §3.6 (keine Gate-Lockerung) | **bestätigt** | der gatende Audit-Schritt läuft in beiden Jobs und ist `success`; weder die Datei im Tap noch der Lauf zeigen `continue-on-error` oder übersprungene Schritte (alle 9 Schritte beider Jobs `success`) |
| Plan §1 (nichts außer Matrix-Zeile und Kommentar) | **bestätigt** | Tap-Commit `+4 −1`, eine Datei, Hunk-Rumpf gleich dem Quell-Hunk |
| Plan §2 „Smoke“ (Messweg über `gh api`) | **bestätigt** | Job-Metadaten und Log per `gh api` gelesen; ([`BEO-DOK/ci-log-falsch-gelesen`](../plan/planning/observations/BEO-DOK/ci-log-falsch-gelesen/observation.md)) |
| Plan §4 Rückführung `in-progress` → `open` (Push verweigert / Lauf startet nicht) | **nicht eingetreten** | Push angenommen, Lauf gelaufen |

### Sensors

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `cmp` Quelle gegen Tap-Datei | **pass** | Exit 0 (Auftragspunkt 1) |
| Lesung Lauf `42212837961` und beider Jobs per `gh api` | **pass** | siehe *Smoke* |
| Lesung der Job-Logs per `gh api …/logs` | **pass** | beide Logs lesbar (232 und 219 Zeilen nach Entfernen der ANSI-Folgen); „Runner Image“-Abschnitt und `u-boot --version` gelesen |
| Lesung der Annotationen | **pass** | siehe *Annotationen* |
| `make docs-check` für diesen Report | **pass** | `0 Befund(e)`, siehe *Sensor-Lauf* |
| `make gates` | not run | Der Nachlauf liest nur den Tap; im Hauptrepo ändert sich nichts außer diesem Report (Datei unter `docs/reviews/`). Der Gate-Lauf des ersten Teils (Tip `3e1d81a9`, grün) deckt den Quelldiff; der Report des ersten Teils (`d665e62b`) kam danach, ohne Code. Für diesen Report genügt der engste Sensor für Doku, `make docs-check` ([`harness/verification.md`](../../harness/verification.md) §Sensor-Auswahl) |
| `make test-docker`, `make verify-depguard`, `make ci`, `make fullbuild` | not run | keine Docker-, Importregel-, Produkt- oder Image-Änderung; Begründung wie im ersten Teil |

### Traceability

| ID | Beleg |
|---|---|
| [`LH-OPEN-002`](../../spec/lastenheft.md#lh-open-002--paketierung) | Der Distributionsweg (Homebrew, eigener Tap) hat im Tap-Lauf `42212837961` auf `macos-15` (arm64) und `macos-15-intel` Formel `0.24.0` installiert, `brew test`, `init`, `add postgres`, `doctor` und das gatende Audit durchlaufen. Das belegt die Funktion der Smoke-Prüfung, nicht mehr |
| [`ADR-0016`](../plan/adr/0016-homebrew-distribution-per-tap.md) | Entscheidung 5: `smoke.yml` im Tap, Prüfungen unverändert; belegt durch Patch und Lauf |

### Carveouts

- Neu: none
- Gelöst: none
- Unverändert: none berührt

### Nicht ausgeführt

- **Ankündigung von GitHub zur Außerbetriebnahme** — weiter nicht nachgelesen (V-2 des ersten Teils); der Lauf macht sie für die Ursachen-Frage entbehrlich (N-1), nicht gelesen bleibt sie trotzdem.
- **Ein Architektur-Kommando auf dem Runner (`uname -m`)** — der Workflow führt es nicht aus, und ein Eingriff in Datei oder Lauf war nicht Teil des Auftrags; „arm64“ stützt sich auf `Image:`, Homebrew-Präfix und Gem-Plattform (siehe *Smoke*).
- **Schreibende Schritte** (Cancel von `42196563230`, Rerun, neuer Dispatch) — nicht Teil dieser Rolle.
- **`make gates`, `make test-docker`, `make verify-depguard`, `make ci`, `make fullbuild`** — Begründung in der Sensor-Tabelle.

### Commit / Artefakt

- Tap-Commit `662c5fb5c57fc61751f5cff127f213be57573211` (Eltern `b9f23331a5…`), Blob `.github/workflows/smoke.yml` `db62a07fb710fb2752857b76bcb09422806a2800` = Blob der Quelle.
- Lauf `smoke` `42212837961`; Jobs `127875452461` (`install (macos-15)`) und `127875452090` (`install (macos-15-intel)`).
- Quellcommit im Hauptrepo `737dae18`; lokaler Tip beim Nachlauf `d665e62b`. Kein Tag, kein Release-Asset.

---

## Findings des Verifiers

| ID | Kategorie | Befund | Auflösung |
|---|---|---|---|
| N-1 | INFO | **Die Ursachen-Annahme des Plans ist durch den Lauf bestätigt.** Plan §6 (erstes Risiko) macht die Ursache („Außerbetriebnahme von `macos-14`“) erst mit dem grünen Lauf unter dem neuen Label zum Beweis. Der Lauf `42212837961` ist grün, `macos-15` hat einen Runner bekommen (`runner_name` nicht leer, `labels` `["macos-15"]`), und keine Annotation meldet eine nicht unterstützte Konfiguration. Das zeigt: Mit dem Label `macos-15` läuft die Datei, wo `macos-14` (Lauf `42191356009`, laut erstem Teil `runner_name` leer, Annotation „not supported“) keinen Runner bekam. Die Ankündigung von GitHub habe ich weiterhin nicht gelesen; die Aussage „Stilllegung zum 2026-11-02“ bleibt Aussage des Plans. | Planner setzt das Risiko bei der Closure auf *entfallen* mit Verweis auf Run und Jobs; die Ursachen-Formulierung bleibt auf das Gelesene beschränkt (Label bedient, Lauf grün), nicht auf die Ankündigung. |
| N-2 | INFO | **Der Dispatch-Lauf `42196563230` steht seit 09:39Z ohne Job auf `queued`** (bei meiner Lesung ~11:52Z, also rund 2 Std. 13 Min.), auf demselben Kopf `662c5fb5…`. Er hat null Jobs; das ist kein „Runner knapp“ (dann gäbe es queued Jobs), und warum er so steht, ist nicht gelesen. Er ist nicht Beleg. Er könnte später starten und dann einen zweiten Lauf mit denselben Jobs erzeugen; ein schreibender Eingriff (Cancel) war nicht meine Rolle. | Hauptagent: entscheidet, ob der Lauf abgebrochen wird; der Planner nennt in §7 beide Läufe getrennt (`42212837961` Beleg, `42196563230` nicht). Plan §6 „Knappe macOS-Runner“ passt nicht wörtlich (keine Jobs): beim Ausgang des Risikos die Formulierung auf das Gelesene begrenzen. |
| N-3 | INFO | **Der Hinweis „Failed to import … vagrant.rb“ ist jetzt gelesen.** Er steht in den Annotationen des Intel-Jobs als Stufe `failure` bei Job-Konklusion `success`, stammt aus dem Schritt „Tap und Install“ und betrifft den Tap `hashicorp/homebrew-tap` im Homebrew-Verzeichnis des Runners (Pfad `/usr/local/Homebrew/Library/Taps/hashicorp/…`), nicht `pt9912/u-boot`; ob er vorinstalliert ist, ist nicht gelesen; der arm64-Job zeigt ihn nicht. Plan §6 (Risiko „Ein neuer Befund unter `macos-15`“) führte ihn als nicht gelesene Angabe des Projektinhabers. Neuer Befund unter `macos-15` (arm64): keiner. | Planner übernimmt den Beleg bei der Closure in den Ausgang des Risikos (entfallen für `macos-15`; der Hinweis am Intel-Job bleibt unverändert und gatet nicht). Kein Gate betroffen. |
| N-4 | INFO | **Der arm64-Beleg hat keine Architektur-Zeile im Log.** „arm64“ trägt `Image: macos-15-arm64`, das Präfix `/opt/homebrew` und die Gem-Plattform `arm64-darwin`; ein Kommando, das die Architektur ausgibt, führt der Workflow nicht aus. Das genügt der DoD („Job-Metadaten und Image-Abschnitt von ‚Set up job‘“), ist aber eine Lesung des Image-Abschnitts, keine Messung auf dem Runner. | Keine Aktion; die Closure nennt die drei Stellen, nicht „gemessen auf dem Runner“. |
| N-5 | LOW | **Offen aus dem ersten Teil: V-1 (Dateiname des Verification-Reports in der Plan-Tabelle) ist nicht berichtigt.** Plan §3, Zeile 161, nennt weiter die Form `docs/reviews/<Datum>-verification-slice-…`; die Reports heißen `verification-<Datum>-<Slice-Kennung>.md` (erster Teil, [Link](verification-2026-11-06-slice-v1-tap-smoke-runner-macos-15.md)) und `…-nachlauf.md` (dieser). Mit diesem Report stehen zwei Dateien statt einer. Die Plan-Aussage über den Dateinamen stimmt nicht mit den Dateien überein (`AGENTS.md` §3.10); der Abgleich gegen `git diff --stat` nach der letzten Fix-Runde wird dadurch beim Closure erneut fällig. | Planner berichtigt die Zeile auf die Ist-Namen (beide Reports) und gleicht Tabelle und `git diff --stat` ab. V-4 (Zustandsmarker in §7) läuft mit diesem Bericht weiter ab: §7 sagt „Verification, Tap-Push und Smoke-Lauf stehen aus“; auch der Tap-Push und der Lauf stehen nicht mehr aus. |

Keine HIGH- und keine MEDIUM-Befunde. Eine DoD-Verletzung, die dem Review oder den Tests entgeht, liegt nicht vor.

## Verdikt

**Nachlauf bestanden — DoD „Tap“ und „Smoke“ sind erbracht und mit eigener Lesung belegt.** Die Datei im Tap ist byte-gleich mit der Quelle (`cmp` Exit 0, gleicher Blob `db62a07f…`), der Tap-Commit `662c5fb5` besteht aus der einen Datei mit `+4 −1` und demselben Hunk wie `737dae18`. Der Lauf `42212837961` (`workflow_dispatch`, Kopf `662c5fb5`) ist `success`; `install (macos-15)` (Job `127875452461`, `labels` `["macos-15"]`, `runner_name` nicht leer, Image `macos-15-arm64`, Präfix `/opt/homebrew`) und `install (macos-15-intel)` (Job `127875452090`) sind beide `success`; installiert wurde die Formel `u-boot 0.24.0`; das gatende Audit lief in beiden Jobs und ist `success`. Die Aussage „grün“ gilt für diese Formel-Version, diese zwei Images und die neun Schritte der Datei, nicht darüber hinaus.

**Mit diesem Bericht erfüllt der Verifier das Verification-Evidence-Minimum von [`harness/verification.md`](../../harness/verification.md) für den Slice:** erster Teil (Quelle, Gates, Plan-Tabelle) plus dieser Teil (Tap, Smoke). Der Slice ist damit **nicht** geschlossen; offen bleibt die Closure (Plan §2 Punkte 8 bis 11).

**Offen für die Closure (Planner):**

1. §7 auf den Ist-Stand setzen und die Belege trennen: gelesen von Planner und Verifier (Tap-Datei, Lauf, Jobs, Logs, Annotationen) und nur berichtet (Push-Ausgang, Dispatch-Auslösung). Tap-Commit `662c5fb5`, Run-ID `42212837961`, beide Job-IDs eintragen (DoD „Tap“ und „Smoke“); `42196563230` als Nicht-Beleg (N-2).
2. Risiko-Ausgänge aus §6: Ursache (N-1) *entfallen*; Schreibzugriff auf den Tap *entfallen* (Push angenommen); Datei-Drift: Vergleich nach der letzten Fix-Runde wiederholen (`cmp` zur Closure); neuer Befund unter `macos-15`: *entfallen* (N-3); Knappe Runner: auf das Gelesene begrenzen (N-2); weiterer Runner-Wandel und Stilllegung erst beim nächsten Formel-Push: *weiter offen* → Register und Bericht an den Projektinhaber.
3. Plan-Tabelle berichtigen (N-5 und V-1) und gegen `git diff --stat` abgleichen.
4. Register: Beobachtung „Runner-Label stillgelegt“ als zweites Vorkommen (nach `macos-13`), gezählt, nicht verkörpert; Kürzel nachschlagen.

**Übergabe an den Hauptagenten:** (a) Entscheidung über den Vorlauf `42196563230` (stehen lassen oder abbrechen; ich habe ihn nicht angefasst). (b) Vor der Closure den Vergleich Quelle gegen Tap erneut ziehen, falls zwischen den Lesungen etwas in einer der beiden Dateien geändert wurde; meine Lesung gilt für den Stand `d665e62b` und Tap-Kopf `662c5fb5`.

## Sensor-Lauf

`make docs-check` für diesen Report: `1204 Datei(en) geprüft, 0 Befund(e)` (der erste Lauf fand eine unverlinkte Slice-Kennung in N-5; berichtigt).
