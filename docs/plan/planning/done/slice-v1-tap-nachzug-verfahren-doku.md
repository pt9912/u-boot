# Slice V1: Verfahren „Änderung der Tap-Dateien aus `packaging/homebrew-tap/` in das Tap-Repo“ in der Maintainer-Doku

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Dokumentations-Lücke aus der Closure von
[`slice-v1-tap-smoke-runner-macos-15`](../done/slice-v1-tap-smoke-runner-macos-15.md); es gibt keine
Closure-Bedingung, die von der DoD dieses Slice verschieden ist).

**Anlass.** Die Closure-Notiz des Vorgängers
[`slice-v1-tap-smoke-runner-macos-15`](../done/slice-v1-tap-smoke-runner-macos-15.md) (§6 „Datei-Drift“, §7
Rest- und Folgepunkt 1) stellt die Frage, wie eine geänderte Vorlage
`packaging/homebrew-tap/.github/workflows/smoke.yml` in das Tap-Repo `pt9912/homebrew-u-boot` kommt; der
Projektinhaber hat am 2026-11-06 entschieden, nur zu dokumentieren (Variante A: kein Skript, kein Sensor, kein
Test, keine Änderung der Token-Rechte, keine Änderung am Job `tap`; die Entscheidung steht als Angabe im
Auftrag, nicht als Artefakt im Repo). Gelesen am 2026-11-06 (Planner): Keine Maintainer-Doku sagt, wie die Datei in den
Tap kommt. `packaging/homebrew-tap/README.md` Zeilen 19 bis 21 nennt nur, dass beim Anlegen des Tap-Repos
README und `smoke.yml` hineinkopiert werden; `docs/maintainer/releasing.md` §2 (Zeilen 50 und 51) nennt „README
und Smoke-Workflow aus `packaging/homebrew-tap/`“ als Inhalt des Tap-Repos, §6 (Zeilen 162 bis 164) den Lauf
des Smoke nach dem Formel-Commit. `scripts/tap-nachzug.sh` klont den Tap, legt `Formula/` an, kopiert die Formel
nach `Formula/u-boot.rb` und fügt nur diese Datei hinzu (Zeilen 25 bis 27 und 34); `smoke.yml` und `README.md`
fasst kein Skript an. Der Smoke läuft bei einem Push auf `Formula/u-boot.rb` und per `workflow_dispatch`
(`smoke.yml` Zeilen 10 bis 14); ein stabiler Tag startet ihn also mit der Datei, die im Tap liegt, nicht mit
der Vorlage. Zustand der beiden Tap-Dateien am 2026-11-06 (Planner, `gh api`, nur GET): `cmp` der Tap-Datei
gegen die Quelle Exit `0` für `.github/workflows/smoke.yml` (Blob `db62a07f…`) und für `README.md` (Blob
`5b333c7e…`), Tap-Kopf `662c5fb5…`. Das Verfahren, das beim Vorgänger lief, ist nur berichtet (Hauptagent):
Tap-Commit `662c5fb5`, Push der Datei per `gh api` `PUT` mit dem Blob-sha der ersetzten Datei, danach `cmp`
gegen die Quelle, danach `gh workflow run smoke --repo pt9912/homebrew-u-boot` und Lesen des Laufs; der
Wortlaut des Push-Aufrufs steht in keinem Artefakt des Repos.

**Bezug:** [`LH-OPEN-002`](../../../../spec/lastenheft.md#lh-open-002--paketierung) (Paketierung: Homebrew
gewählt, eigener Tap), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) (Entscheidung 5: Smoke im
Tap).

**Berührte Spec-Stellen:** — (Maintainer-Doku, kein Vertrag: weder Lastenheft noch Spezifikation noch Sicht
nennen das Verfahren; [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 nennt `smoke.yml`
und seine Prüfungen, keinen Weg in den Tap).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner, im Auftrag des Hauptagenten). **Datum:** 2026-11-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Ziel:** Die Maintainer-Doku `docs/maintainer/releasing.md` beschreibt in einem eigenen Abschnitt (§9), wie
eine Änderung an `README.md` oder `.github/workflows/smoke.yml` aus `packaging/homebrew-tap/` in das Tap-Repo
`pt9912/homebrew-u-boot` kommt und wie der Maintainer prüft, dass beide Seiten gleich sind und der Smoke läuft;
die Stellen, an denen ein Maintainer das Verfahren sucht (§2, §4, §8 und der Index `docs/maintainer/README.md`),
zeigen auf diesen Abschnitt.

**Planentscheidung: Ort der Doku ist `docs/maintainer/releasing.md`, ein neuer §9.** Umgesetzt ist dieser Ort; die
Entscheidung traf der Hauptagent (Herkunft in §7).

- **Grund für `releasing.md`:** Das Dokument führt den Tap schon an drei Stellen — §2 (Tap-Repo und Token), §6
  (Smoke nach dem Formel-Commit), §8 (Fehlerbilder des Jobs `tap`) — und ein Maintainer, der den Tap nach einem
  roten Smoke sucht, liest es. Ein neuer Abschnitt braucht weder eine neue Datei noch einen neuen Index-Eintrag;
  `docs/maintainer/README.md` Zeile 10 nennt das Dokument mit dem „Homebrew-Tap“. Die Nummerierung der
  bestehenden Abschnitte ändert sich nicht (§9 steht hinter §8 samt dessen Schlussabsatz).
- **Gegenposition 1: `packaging/homebrew-tap/README.md`.** Verworfen: Die Datei ist selbst eine der Dateien, die
  ins Tap-Repo gehen (und dort die Nutzer-README des Tap); ein Satz dort ändert die Quelle und verlangt damit
  genau den Push, den der Satz beschreibt. Aus demselben Grund steht im Kommentar von `smoke.yml` kein Zeiger auf den
  neuen Abschnitt: Jede Änderung beider Dateien ist eine Änderung im Tap (Variante A: keine Änderung am Tap).
- **Gegenposition 2: eine eigene Datei `docs/maintainer/homebrew-tap.md`.** Möglich, wenn der Projektinhaber
  die Tap-Pflege vom Release-Ablauf trennen will; Kosten: eine Datei, eine Zeile im Index
  `docs/maintainer/README.md`, und die drei Stellen in `releasing.md` verweisen dorthin. Der Planner wählt sie nicht, weil
  die Änderung einer Tap-Datei in der Praxis am Release hängt: ein stabiler Tag startet den Smoke mit der Datei im Tap
  (§Anlass), und das Rot fällt erst dort auf.
- **Neubewertung, wenn:** `releasing.md` durch weitere Tap-Verfahren (zweiter Tap, Änderung der Tap-Rechte)
  über einen Abschnitt hinauswächst, oder der Projektinhaber Variante B oder C (Skript, Sensor) wählt; dann gehört
  das Verfahren in die Beschreibung dieses Werkzeugs, und der Planner legt das Ergebnis als eigenen Slice vor.
- **Keine ADR:** Das Verfahren ändert weder Entscheidung 5 von [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md)
  noch ein Gate; keine Accepted-ADR wird berührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Skript und kein Sensor, der `smoke.yml` oder `README.md` in den Tap schreibt oder Quelle und Tap
  vergleicht.** Entscheidung des Projektinhabers (Variante A, 2026-11-06). Die Drift-Gefahr bleibt: nach jedem
  Verfahrenslauf kann jede Seite allein geändert werden, und nur der Vergleich im Verfahren (und der Sichtschritt
  in `releasing.md` §4, Entscheidung 2 in §4) fängt es; §6 führt sie als offen.
- **Kein Test.** Kein Test liest `smoke.yml` (Vorgänger-Plan §3, Zeile „Tests“) oder pinnt `scripts/tap-nachzug.sh`
  auf „schreibt nur die Formel“ (`cmd/uboot/homebrew_scripts_test.go` hat `TestHomebrewFormulaFill_Fills`,
  `TestHomebrewFormulaDesc_EnglishASCII`, `TestHomebrewFormulaFill_FailsClosed`, `TestTapFormulaCheck`, gelesen
  die Funktionsnamen, Zeilen 38, 72, 103, 132). Ein Satz der Doku über das Skript ist deshalb ein Lesebefund mit
  Fundstelle, keine Zusage (§6).
- **Keine Änderung der Token-Rechte und keine Änderung am Job `tap`.** Entscheidung des Projektinhabers (Variante A).
  Der Job zieht weiter nur die Formel nach (`publish.yml` Zeilen 312 bis 334, `scripts/tap-nachzug.sh`); das Secret
  `HOMEBREW_TAP_GITHUB_TOKEN` ist laut `releasing.md` §2 ein PAT mit Contents Read & Write. Ob es `.github/workflows/`
  schreiben dürfte, ist nicht gemessen und nicht Gegenstand; die Doku sagt dazu nichts über dieses Token.
- **Kein Schreibzugriff auf den Tap und keine Änderung in `packaging/`.** Der Slice liest den Tap nur
  (`gh api`, GET); die Dateien `README.md` und `smoke.yml` bleiben Byte für Byte, in der Quelle wie im Tap.
- **Kein regelmäßiger Lauf (`schedule`) und kein Umgang mit dem hängenden Lauf `42196563230`.** Beides sind die Punkte 2
  und 3 der Closure des Vorgängers ([`slice-v1-tap-smoke-runner-macos-15`](../done/slice-v1-tap-smoke-runner-macos-15.md) §7)
  und stehen beim Projektinhaber.
- **Kein CHANGELOG-Eintrag, kein Release, kein Tag.** Maintainer-Doku ändert weder Binary noch Formel noch Pakete.
- **Kein Eintrag in `harness/README.md` und keine Änderung an einem Vertrag.** Es entsteht kein Target und kein Gate.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Verfahren:** `docs/maintainer/releasing.md` trägt einen neuen §9 „Dateien des Tap-Repos nachziehen
      (`packaging/homebrew-tap/`)“ mit diesem Inhalt (Wortlaut beim Implementer):
      (a) *Geltung:* `README.md` und `.github/workflows/smoke.yml`; die Formel `Formula/u-boot.rb` zieht der Job
      `tap` nach (§1, §6) und gehört nicht hierher.
      (b) *Rechte, als Hinweis:* GitHubs Seite „Create or update file contents“ (gelesen 2026-11-06) nennt für klassische
      PATs den Scope `repo` und zusätzlich `workflow`, um Dateien unter `.github/workflows` zu ändern, und für
      fine-grained PATs Contents (write) sowie Contents (write) mit Workflows (write); die Doku behauptet nichts über
      ein bestimmtes Token — welche Rechte das genutzte Token trägt, zeigt nur der Push.
      (c) *Schritte, in dieser Reihenfolge:* 1. Ausgangslage vergleichen (`git hash-object` der Quelle gegen
      `gh api repos/pt9912/homebrew-u-boot/contents/<Pfad> --jq .sha`; bei Abweichung zuerst klären, was nur im Tap
      steht); 2. Quelle ändern, `make gates` im Devcontainer-Container, auf `main` committen, Review und Verification
      der Quelle vor dem Tap-Push (die Quelle führt, der Tap folgt); 3. Push per `gh api --method PUT` auf
      `contents/<Pfad>` mit `message`, `content` (base64 der Quelle), `sha` (aus Schritt 1) und `branch`; der Push
      startet keinen Lauf (Trigger nur `Formula/u-boot.rb` und `workflow_dispatch`, `smoke.yml` Zeilen 10 bis 14);
      4. Nachvergleich (`gh api … --jq .content | base64 -d | cmp - <Quelle>` und der sha-Vergleich); 5. Lauf starten
      (`gh workflow run smoke --repo pt9912/homebrew-u-boot`) und lesen: den neuen Lauf an `created_at` und `head_sha`
      erkennen (ein älterer Dispatch-Lauf kann daneben in der Liste stehen), Jobs, `labels`, `runner_name` und
      `conclusion` mit `gh api …/actions/runs/<id>/jobs`, Annotationen mit `gh api …/check-runs/<job-id>/annotations`,
      Log mit `gh api --allow-escape-sequences …/actions/jobs/<job-id>/logs` (nicht allein `gh run view --log`); „grün“
      heißt die installierte Formel-Version aus dem Log und die Image-Zeile, nicht „macOS grün“; 6. Rückweg: ein
      weiterer Commit mit dem vorigen Inhalt, derselbe Aufruf mit dem dann aktuellen `sha`.
      (d) *Grenze, im Indikativ:* Kein Skript, kein Gate und kein Test hält Quelle und Tap gleich oder schreibt diese
      zwei Dateien in den Tap; der Vergleich der Schritte 1 und 4 ist die einzige Prüfung, dazwischen kann jede Seite
      allein geändert werden, und ein stabiler Tag startet den Smoke mit der Datei, die im Tap liegt. Für den
      Schreibschritt (3) gibt es keine Probe, die nichts schreibt.
      Der Abschnitt enthält keine Run-, Commit- oder Job-IDs als Beleg und keinen Zustandsmarker mit Ablaufpunkt
      (`AGENTS.md` §3.10); jede Aussage über `scripts/tap-nachzug.sh`, `scripts/tap-check.sh` oder `smoke.yml` nennt ihre
      Fundstelle (Datei und Symbol oder Zeilen) als Lesebefund. Jeder lesende Befehl des Abschnitts ist vom Implementer
      gegen den Tap ausgeführt (Ausgabe-Auszug in §7); der Schreibbefehl ist als ungeprobt ausgewiesen, sein Wortlaut
      stammt vom Aufruf des Hauptagenten beim Vorgänger oder, liegt er nicht vor, aus der API-Beschreibung (Entscheidung 4
      in §4).
- [x] **Auffindbarkeit:** Vier Stellen zeigen auf §9, jede mit höchstens einem Satz: `releasing.md` §2 (Aufzählung
      „Tap-Repo“, Zeilen 50 und 51), §4 (ein Satz: vor einem stabilen Tag Quelle und Tap der zwei Dateien vergleichen,
      Schritt 1 aus §9), §8 (eine Tabellenzeile: Smoke im Tap rot, Job ohne Runner, Annotation
      „The configuration '…' is not supported“ → Runner-Label der Vorlage ändern und nach §9 in den Tap bringen) und der Index
      `docs/maintainer/README.md` (Zeile 10 nennt das Nachziehen der Tap-Dateien). Der Satz in §4 ist ein Sichtschritt und kein
      Gate; er steht nur, wenn der Projektinhaber Entscheidung 2 in §4 so trifft.
- [x] **Nichts im Tap und in `packaging/` geändert:** `git diff --stat` über die Spanne des Slice nennt keine Datei
      unter `packaging/`, `scripts/` oder `.github/`; die Dateien `README.md` und `.github/workflows/smoke.yml` im Tap sind
      bei der Closure byte-gleich mit der Quelle (`cmp` Exit `0`, sha gleich), vom Verifier unabhängig gelesen.
- [x] `make gates` grün (im Devcontainer-Container, `AGENTS.md` §3.1).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Verification-Report ([`harness/verification.md`](../../../../harness/verification.md)) liegt vor; der Verifier
      führt die lesenden Befehle von §9 selbst aus und gleicht die Plan-Tabelle gegen `git diff --stat` ab.
- [x] Doku-Update: dieser Slice ist das Doku-Update (Maintainer-Doku); kein öffentlicher Vertrag berührt, daher keine
      Datei unter `docs/user/`, kein README, kein CHANGELOG; die Zeilen „unverändert“ in §3 sind mit gelesenen Stellen
      belegt und bei der Closure gegen `git diff --stat` abgeglichen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht. Jede Zeile „unverändert“ trägt die Stelle, die gelesen wurde
(`AGENTS.md` §3.10); Zeilennummern gelten für den Stand vor dem Diff.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/maintainer/releasing.md` | update | Neuer §9 hinter §8 samt dessen Schlussabsatz (Zeilen 191 und 192), Inhalt nach der DoD „Verfahren“. Dazu drei kurze Verweise: §2, Aufzählungspunkt „Tap-Repo“ (Zeilen 50 und 51); §4, Absatz „Zusätzlich:“ (Zeilen 118 bis 121) um den Sichtschritt „Tap-Dateien vergleichen“; §8, eine Tabellenzeile hinter „Formel fehlt im Tap“ (Zeile 187). Die Nummerierung §1 bis §8 bleibt. Tabellenzeile und §9 nennen kein Make-Target, das es nicht gibt (`docs-check`, Modul `targets`). |
| `docs/maintainer/README.md` | update | Zeile 10 beschreibt `releasing.md` als „Release-Ablauf (Tag, `publish`, GHCR, Binaries, Pakete, Homebrew-Tap, Fehlerbilder)“; sie nennt zusätzlich das Nachziehen der Tap-Dateien. |
| `docs/plan/planning/done/slice-v1-tap-nachzug-verfahren-doku.md` | neu | dieser Plan, im Endstand am Ruheort `done/`; die Zustandswechsel sind reine `git mv` (open → next → in-progress → done) |
| `docs/plan/planning/open/README.md` | update | Bestandszeile des Plans (Tabelle „Bestand“, erste Zeile); die Zeile folgt dem Slice bei jedem `git mv` |
| `docs/plan/planning/in-progress/roadmap.md` | unverändert bis zum Start; dann update, Nettodiff der Spanne leer | Gelesen: Zeile 45 „Nichts in Arbeit.“ (Ruhe-Marker); `.d-check.yml` (`planning`) verlangt ihn genau dann, wenn `in-progress/` keinen Slice trägt. In `open/` und `next/` bleibt die Datei, wie sie ist (die Roadmap führt die Slices in `open/` nicht einzeln; Zeile 28 nennt nur einen davon). Mit dem Start-Commit entfällt die Zeile, mit der Closure steht sie an derselben Stelle wieder; die Datei fehlt deshalb in `git diff --stat` der Spanne. |
| `docs/reviews/2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md` | neu | Review-Report (`.harness/skills/reviewer.md`), Gegenstand `11fc3b01` gegen `0afd3272`; Datum des Review-Laufs 2026-11-06, gleich dem Namen |
| `docs/reviews/verification-2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md` | neu | Verification-Report; Namensform `verification-<Datum>-<Kennung>.md` wie die neueren Berichte in `docs/reviews/`; eine Datei, mit zwei Nachträgen (Nachtrag, Nachtrag 2 zu V-1 an `f902d127`); ein Nachlauf-Teil zum Tap entsteht nicht, weil nichts im Tap geschrieben wird |
| `docs/plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/` (`evidence/slice-v1-tap-nachzug-verfahren-doku.md` neu, `state.md` update) | neu / update, bei der Closure | Review F-1 und Verification V-1 sind ein Auftreten der Klasse (Aussage über Verhalten ohne gelesene Stelle in der gelieferten Doku); Eintrag nur dort, wo ein Auftreten vorliegt, andere Beobachtungen tragen keine Datei (§7). |
| Tests | keine Änderung | Gelesen `cmd/uboot/homebrew_scripts_test.go`: die Funktionen `TestHomebrewFormulaFill_Fills`, `TestHomebrewFormulaDesc_EnglishASCII`, `TestHomebrewFormulaFill_FailsClosed` und `TestTapFormulaCheck` (Zeilen 38, 72, 103, 132) prüfen das Füllen der Formel und `scripts/tap-formula-check.sh`; kein Test ruft `scripts/tap-nachzug.sh`, `scripts/tap-check.sh` oder liest `smoke.yml` oder `releasing.md`. Eine Suche nach `releasing` über `*.go`, `*.sh`, `*.mk`, `*.yml` und `*.yaml` fand keine Datei, die den Text des Dokuments liest (Treffer nur in Dokumenten). |
| `packaging/homebrew-tap/.github/workflows/smoke.yml` | unverändert | Gelesen Zeilen 1 bis 70: Kopfkommentar Zeilen 3 bis 8 (Zweck, Docker fehlt auf dem Runner), Trigger Zeilen 10 bis 14 (`push` auf `Formula/u-boot.rb`, `workflow_dispatch`, kein `schedule`), Matrix-Kommentar Zeilen 24 bis 26 und `runner: [macos-15, macos-15-intel]` in Zeile 27. Kein Zeiger auf die Doku: Jede Änderung der Datei ist eine Änderung, die in den Tap muss (Variante A). |
| `packaging/homebrew-tap/README.md` | unverändert | Gelesen Zeilen 1 bis 21: Nutzer-README des Tap; Zeilen 19 bis 21 nennen das Anlegen des Tap-Repos. Kein Zeiger auf die Doku aus demselben Grund. |
| `scripts/tap-nachzug.sh` | unverändert | Gelesen: Klon des Tap Zeile 25, `mkdir -p "$work/tap/Formula"` Zeile 26, `cp … Formula/u-boot.rb` Zeile 27, `add Formula/u-boot.rb` Zeile 34, `git push origin HEAD` Zeile 37; Quelle der Doku-Aussage „kein Skript schreibt die zwei Dateien“ (Lesebefund, kein Test, §6). |
| `scripts/tap-check.sh` | unverändert | Gelesen Zeilen 32 bis 41: Klon und `git push --dry-run origin "HEAD:refs/heads/${branch}"` (Zeile 34) des geklonten Stands; der Befehl pusht keinen neuen Commit, belegt also kein Recht für `.github/workflows/` (Lesebefund). |
| `.github/workflows/publish.yml`, `.github/workflows/tap-check.yml` | unverändert | `publish.yml` Job `tap` Zeilen 312 bis 334: `if` auf `is_stable`, Schritt „Pull the formula into the tap“ ruft `scripts/tap-nachzug.sh`; `tap-check.yml` gelesen ganz: ein Job `check`, Aufruf von `scripts/tap-check.sh`. `integration.yml` gelesen (Job `integration-docker`, `ubuntu-latest`, `make test-docker`, kein Tap-Pfad); `ci.yml` gelesen in den Job-Köpfen (`gates` Zeile 43, `security-gates` Zeile 87, `image-scan` Zeile 102: Hauptrepo-Gates, kein Tap-Pfad). |
| `harness/README.md` | unverändert | Gelesen Abschnitte Sensors und Werkzeuge: keine Zeile zu Tap, Smoke oder `releasing.md`; es entsteht kein Target. |
| `CHANGELOG.md` | unverändert | Gelesen: `[Unreleased]` leer (Zeile 12, der nächste Abschnitt `[0.24.0]` beginnt in Zeile 14); Maintainer-Doku ändert das Produkt nicht. |
| `docs/user/benutzerhandbuch.md`, `docs/maintainer/branch-protection.md` | unverändert | Handbuch Zeilen 114 bis 124: nennt den Tap als Installationsweg (`brew tap`, `brew trust`, `brew install`), nichts zu `smoke.yml` oder zum Nachziehen; `branch-protection.md` Zeilen 80 bis 83: „Pushes in andere Repos mit dem persönlichen Token (Homebrew-Tap)“ sind von der Regel auf `main` des Hauptrepos nicht betroffen — der Satz gilt dem Hauptrepo und bleibt. |
| `docs/plan/adr/0016-homebrew-distribution-per-tap.md` | unverändert | Accepted; gelesen Zeile 29 (Entscheidung 5): nennt `smoke.yml` als „Vorlage für das Tap-Repo“ und seine Prüfungen, keinen Weg in den Tap. |

**Ansatz in der Reihenfolge der Rollen:** Implementer schreibt §9 und die vier Verweise, führt jeden lesenden Befehl des
Abschnitts gegen den Tap aus (Auszüge nach §7), setzt den Schreibbefehl als ungeprobt und läuft `make gates` → Reviewer →
Verifier (liest den Tap selbst: `cmp` und sha für beide Dateien; gleicht die Plan-Tabelle gegen `git diff --stat`) →
Closure. Es gibt keinen Push in den Tap und keinen Dispatch; die Rolle Hauptagent liefert nur den Wortlaut des Push-Aufrufs
vom Vorgänger (Entscheidung 4 in §4).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Freigabe zum Umsetzen liegt vor, wenn der Projektinhaber die vier
Entscheidungen unten getroffen oder die Vorschläge bestätigt hat. Der Hauptagent beginnt mit dem Implementer-Lauf,
sobald in `in-progress/` kein anderer Slice liegt (WIP-Limit 1; beim Schnitt liegt dort nur `roadmap.md`). Die Übergänge
`open` → `next` → `in-progress` sind je ein reiner `git mv`.

**Vor dem Start zu entscheiden** (Vorschlag des Planners je Punkt):

1. **Ort** der Doku: `docs/maintainer/releasing.md` §9 (Vorschlag) oder die eigene Datei `docs/maintainer/homebrew-tap.md`
   (§1, Gegenposition 2).
2. **Der Satz in `releasing.md` §4** (vor einem stabilen Tag Quelle und Tap der zwei Dateien vergleichen): aufnehmen
   (Vorschlag; er ist der einzige Ort, an dem die Prüfung vor dem Release steht) oder streichen. Er ist Doku, kein Gate und
   kein Schritt, den ein Sensor prüft.
3. **Der Zeiger im Ausgang der Drift-Gefahr** (§6, Risiko 1): Entscheidungs-Zeiger auf Variante A (Vorschlag, weil das
   Register Auftreten zählt und noch keines vorliegt) oder eine neue Beobachtung, die die Closure ohne Auftreten anlegt.
4. **Wortlaut des Push-Aufrufs:** Der Hauptagent liefert den Aufruf, der den Tap-Commit `662c5fb5` erzeugt hat (kein Artefakt
   im Repo trägt ihn); liegt er nicht vor, schreibt der Implementer den Aufruf nach der API-Beschreibung und weist ihn als
   ungeprobt aus. Eine Probe gegen den Tap (etwa mit absichtlich veraltetem `sha`) ist ein Schreibaufruf in ein fremdes Repo
   und nicht vorgeschlagen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Braucht der Abschnitt mehr als Dokumentation — ein Skript, einen
  Sensor, einen Test, eine Änderung der Token-Rechte oder des Jobs `tap` (Varianten, die der Projektinhaber am 2026-11-06
  nicht gewählt hat) —, ist der Zuschnitt falsch; der Mehrbedarf wird ein eigener Slice, mit dem Architect.
- `in-progress` → `open` (blockiert): Liest der Implementer den Tap nicht (kein Zugang mit `gh`), kann weder er noch der
  Verifier die lesenden Befehle ausführen, und der Slice kehrt zurück, bis der Zugang steht. Ein Carveout entsteht dabei
  nicht: kein gelockertes Gate, sondern ein fehlender Zugang.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (zwei Liefer-Punkte, die Zusicherung „nichts im Tap und in `packaging/` geändert“ plus Gate- und
Closure-Pflichten); Review und Verification; Closure-Notiz mit Lerneintrag; `git mv` nach `done/`. Die Closure liest den Tap
selbst (`gh api`, nur lesend) und wiederholt den Vergleich beider Dateien gegen die Quelle sowie den Abgleich der
Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde (`AGENTS.md` §3.10); was sie selbst gelesen hat und was nur
der Hauptagent berichtet, trennt §7. Lerneintrag: benannte Lücke — das Verfahren eines Schritts in einem Fremd-Repo stand nur
im Auftrag und in der Closure des Vorgängers, jetzt in `releasing.md`; die Drift-Gefahr bleibt als Grenze des Abschnitts
benannt (Variante A). Es wird nichts verkörpert, das einen Herkunfts-Anker braucht (die 3×-Schwelle ist nicht erreicht), also
trägt der Eintrag kein Feld `liegt in …`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst. Der Ausgang steht bei der Closure; vorab steht, woran er
sich entscheidet.

- **Drift zwischen Quelle und Tap.** Kein Skript, kein Gate und kein Test hält `README.md` und `smoke.yml` in Quelle und Tap
  gleich (Variante A); der Vergleich im Verfahren belegt den Stand zum Zeitpunkt des Vergleichs. Ein stabiler Tag startet den
  Smoke mit der Datei im Tap, eine Abweichung fällt also erst beim nächsten Formel-Push auf. Die Doku mindert das Risiko
  (Schritt 1 und 4, Sichtschritt in §4), hebt es nicht auf. — **Ausgang: weiter offen** → Entscheidung des Projektinhabers
  vom 2026-11-06 (Variante A, Angabe im Auftrag) und die Grenze (d), der Absatz „Grenze“ in `releasing.md` §9. Quelle und Tap
  sind bei der Closure gleich (§7, `cmp` Exit `0`, sha gleich); gehalten wird das von nichts. Eine Beobachtung im Register
  entsteht mit dem ersten Auftreten einer Abweichung, nicht vorher (Entscheidung 3 in §4, der Zeiger auf Variante A).
- **Der Schreibbefehl ist ungeprobt.** Der Push lässt sich im Tap nicht ohne Schreiben proben; der Wortlaut in §9 kann in einem
  Detail abweichen (Hilfsprogramme wie `base64` sind je Plattform nicht geprüft), und der Fehler fiele erst beim nächsten Push
  auf. — **Ausgang: weiter offen** → die Grenze (d), der Absatz „Grenze“ in `releasing.md` §9. Der Wortlaut des Aufrufs liegt in
  keinem Artefakt; der Absatz „Gelaufen und nicht gelaufen“ weist aus, dass die Zeile für `B64` im abgedruckten Block nicht
  als gelaufen gilt, der abgedruckte Block also nicht der gelaufene Aufruf ist (Verification V-2). Gelesen ist das Ergebnis
  (Tap-Commit, eine Datei, +4 −1), nicht der Aufruf. Entscheidung 4 ist damit nicht „entfallen“, sondern bleibt ein Bericht
  des Hauptagenten (§7).
- **Der Rechte-Hinweis geht über das Gelesene hinaus.** Die Doku zitiert eine Seite von GitHub (gelesen am 2026-11-06),
  nicht das Recht eines Tokens; welches Token einen Push von `.github/workflows/smoke.yml` trägt, ist nicht gemessen.
  — **Ausgang: entfallen.** §9 führt den Hinweis als Zitat mit Quelle und sagt über `HOMEBREW_TAP_GITHUB_TOKEN` nichts
  (gelesen, Abschnitt „Rechte, als Hinweis“); der Verifier hat die Seite von GitHub gelesen und das Zitat bestätigt, der
  Review nicht (Review-Negativbefund „nicht geprüft“). Die Rechte eines bestimmten Tokens bleiben ungemessen.
- **Ein Satz der Doku über ein Skript ohne Test.** „Kein Skript schreibt die zwei Dateien in den Tap“ liest
  `scripts/tap-nachzug.sh`; ändert jemand das Skript, wird der Satz falsch, und kein Test merkt es (die Klasse von
  `AGENTS.md` §3.10, dort für Spezifikation und Handbuch). — **Ausgang: weiter offen** → [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
  (Klasse nahe, Ort ist Maintainer-Doku; kein Auftreten gezählt, weil der Review keinen Befund zum Satz meldet). Die Sätze
  stehen als Lesebefund mit Fundstelle, tragen aber „Kein Skript, kein Gate und kein Test …“ und „die einzige Prüfung“ und
  damit das Plan-Kriterium „kein ‚nur‘, ‚nie‘ oder ‚immer‘“ dem Sinn nach nicht ganz; der Review hat den Satz per Suche
  gegengelesen (Stand des Laufs, ohne Sensor). Der Verifier empfahl „entfallen“; die Closure führt es als offen, weil kein
  Test den Satz pinnt.
- **Der falsche Lauf wird gelesen.** Die Lauf-Liste des Smoke trägt neben dem neuen Dispatch-Lauf einen älteren ohne Job
  (gelesen am 2026-11-06 mit `gh run list` und `gh api`: Lauf `42196563230`, Zustand `queued`, `total_count` der Jobs `0`;
  Rest- und Folgepunkt 3 des Vorgängers); wer den ersten Eintrag nimmt, liest den falschen. — **Ausgang: entfallen.** Schritt 5 in §9 erkennt den Lauf an `created_at`; `head_sha` ist der Commit aus Schritt 3
  und für beide Dispatch-Läufe gleich, so dass §9 hier genauer ist als der Plan (Verification V-4). Der Verifier hat den
  Schritt gegen die Lauf-Liste gelesen; die Closure hat die Liste am 2026-11-06 erneut gelesen (zwei Dispatch-Läufe auf
  `662c5fb5`, der ältere ohne Job).

## 7. Closure-Notiz

**Abschluss:** 2026-11-06. **Verantwortlich:** —. Die Pfade dieses Abschnitts lösen vom Ruheort `done/` auf
(Ruheort-Regel, Baseline-Regelwerk `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln). Die
Belege sind nach Herkunft getrennt: was der Planner selbst las, was Reviewer und Verifier lasen und der Planner nicht
wiederholte, was nur der Hauptagent berichtet, und was niemand las.

- **Geliefert (Fundstellen):** `docs/maintainer/releasing.md` trägt den neuen §9 „Dateien des Tap-Repos nachziehen
  (`packaging/homebrew-tap/`)“ (Geltung, „Rechte, als Hinweis“, sechs Schritte, „Gelaufen und nicht gelaufen“,
  „Grenze“) und drei Verweise: §2 (Klammer im Punkt „Tap-Repo“), §4 (Absatz „Sichtschritt vor einem stabilen Tag, kein
  Gate“), §8 (Tabellenzeile „Smoke im Tap rot, ein Job ohne Runner …“ hinter „Formel fehlt im Tap“);
  `docs/maintainer/README.md` Zeile 10 nennt „Nachziehen der Tap-Dateien“. Commits: Umsetzung `11fc3b01`
  (`11fc3b01c4b0f77565eedd19031e51f89bd186d3`, zwei Dateien), Review-Report `8558da8b`, Review-Fix `8cc97257` (F-1 bis
  F-3), Verification-Report `fb7530c4`, Verification-Fix `f902d127` (V-1), Verification-Nachtrag 2 `b4a09b4d`. Unter
  `packaging/`, `scripts/` und `.github/` ändert die Spanne nichts; im Tap ist nichts geschrieben.
- **Vom Planner gelesen** (`gh api`, nur GET, und lokales `git`/`cmp`, 2026-11-06 bis 12:34Z; Stand HEAD `b4a09b4d`, nach
  der letzten Fix-Runde `f902d127`): `git hash-object` der Quelle und `.sha` im Tap gleich für
  `.github/workflows/smoke.yml` (`db62a07fb710fb2752857b76bcb09422806a2800`) und `README.md`
  (`5b333c7e886ae1db79e763254aa55de6eeda883b`); `cmp` der Tap-Datei (`contents`-API, base64 dekodiert) gegen die Quelle
  Exit `0` für beide; `commits/main` im Tap `662c5fb5c57fc61751f5cff127f213be57573211`; `commits/662c5fb5`: Eltern
  `b9f23331a5…`, genau eine Datei (`.github/workflows/smoke.yml`, `modified`, +4 −1); Lauf-Liste von `smoke.yml`:
  `42212837961` (`workflow_dispatch`, `completed`/`success`, `head_sha` `662c5fb5…`, 11:47:29Z), `42196563230`
  (`workflow_dispatch`, `queued`, 09:39:00Z, gleicher `head_sha`, `…/jobs` `total_count` `0`), `42191356009` (`push`,
  `head_sha` `b9f23331…`, `cancelled`); `git diff --stat cf4c3cb8 HEAD -- packaging scripts .github` ohne Ausgabe;
  die Fundstellen der Zeilen „unverändert“ in §3: `smoke.yml` Zeilen 10 bis 14, `scripts/tap-nachzug.sh` Zeilen 25 bis
  27, 34 und 37, `scripts/tap-check.sh` Zeile 34, `packaging/homebrew-tap/README.md` Zeilen 19 bis 21, `CHANGELOG.md`
  `[Unreleased]` leer (Zeile 12), [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md) Entscheidung 5 (Zeile 29),
  `publish.yml` Job `tap` ab Zeile 312 (`if` auf `is_stable`), `branch-protection.md` Zeilen 80 bis 83; `releasing.md`
  §9 ganz, am HEAD. Nicht neu gelesen: `tap-check.yml`, `integration.yml`, `ci.yml`, `harness/README.md`,
  `docs/user/benutzerhandbuch.md`; sie stehen nicht in der Stat.
- **Vom Reviewer gelesen** ([Review-Report](../../../reviews/2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md), vom
  Planner nicht wiederholt): die lesenden Befehle der Schritte 1, 4 und 5 gegen den Tap; Annotationen des Laufs
  `42212837961`; die Zeilenangaben der Lesebefunde gegen `11fc3b01`; eine Suche nach `homebrew-tap` und `smoke.yml`
  außerhalb von `docs/` (einzige Fundstelle `packaging/homebrew-tap/README.md`). Nicht geprüft: das Zitat der Seite
  „Create or update file contents“.
- **Vom Verifier gelesen** ([Verification-Report](../../../reviews/verification-2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md),
  vom Planner nicht wiederholt): alle lesenden Befehle von §9 selbst, dazu die Zeile für `B64` lokal (3368 Zeichen, Rücklesen
  gleich); der Patch des Tap-Commits `662c5fb5` (Matrix `[macos-14, macos-15-intel]` → `[macos-15, macos-15-intel]` plus
  drei Kommentarzeilen) und dass auf diesem Commit kein `push`-Lauf liegt; die Seite „Create or update file contents“ der
  GitHub-REST-Dokumentation (`docs.github.com/en/rest/repos/contents`), deren Zitat in §9 stimmt; Jobs, `labels`,
  `runner_name`, Log und Annotationen von `42212837961`; `make gates` am Stand `8cc97257` und `f902d127`
  (Nachtrag 2), den Diff `8cc97257..f902d127`.
- **Auszug der lesenden Befehle (DoD „Verfahren“, letzter Absatz):** Die DoD verlangt den Lauf *durch den Implementer* mit
  Auszug in §7. Ein Artefakt des Implementer-Laufs liegt nicht vor (die Commit-Nachricht von `11fc3b01` nennt keinen Lauf);
  ausgeführt belegt sind die Befehle durch Reviewer, Verifier und Planner (oben). Der Haken ist mit dieser Abweichung vom
  Wortlaut gesetzt; die DoD-Zeile bleibt unverändert. Auszug (Planner): Schritt 1 und 4: je `db62a07f…`/`db62a07f…` und
  `5b333c7e…`/`5b333c7e…`, `cmp` Exit `0`; Schritt 5: die Lauf-Liste oben, Spalten Kennung, `created_at`, `head_sha`, Event,
  Status, Ergebnis. Jobs, Annotationen und Log von Schritt 5 hat der Planner nicht neu gelesen (Verifier-Tabelle, Zeilen
  „5“).
- **Nur berichtet (Hauptagent und Projektinhaber), nicht gelesen:** (1) der Wortlaut des Push-Aufrufs, der `662c5fb5`
  erzeugt hat: Herkunft Sitzungsprotokoll des Hauptagenten, in keinem Artefakt des Repos; die Commit-Kennung
  `662c5fb5` stammt aus dem Vorgänger [`slice-v1-tap-smoke-runner-macos-15`](../done/slice-v1-tap-smoke-runner-macos-15.md)
  und steht aus diesem Grund hier und nicht in `releasing.md`; (2) die Wahl der Variante A durch den Projektinhaber am
  2026-11-06 (Angabe im Auftrag, kein Artefakt); (3) die **vier Entscheidungen zu §4** (Ort §9, Sichtschritt in §4, Zeiger
  auf Variante A ohne neue Beobachtung, Push-Aufruf aus dem Sitzungsprotokoll) traf der Hauptagent auf Grundlage der
  Variante-A-Wahl; der Projektinhaber hat sie nicht einzeln freigegeben. Der Start-Trigger von §4 („die vier Entscheidungen
  getroffen oder die Vorschläge bestätigt“) ist damit dem Wortlaut nach nicht erfüllt; die Nachträglichkeit steht als
  Rest- und Folgepunkt 1.
- **Von niemandem gelesen:** die Rechte von `HOMEBREW_TAP_GITHUB_TOKEN` und jedes anderen Tokens für
  `.github/workflows/` (§9 behauptet sie nicht); der Schreibschritt (`PUT`) für `README.md`, der Rückweg (Schritt 6) und
  der Dispatch von Schritt 5 als Aufruf; das Verhalten von `base64 -w0` auf einer anderen Plattform als Linux; der Inhalt des
  Sitzungsprotokolls. Die Aussage „Quelle und Tap sind gleich“ gilt für die zwei Dateien, den Tap-Kopf `662c5fb5` und die
  Lesezeit oben.
- **Was hat funktioniert:** Die Messstellen standen vorab im Plan (Tap-Dateien und Lauf lesen, Rechte nicht behaupten, den
  Schreibschritt als ungeprobt ausweisen) und wurden von Reviewer und Verifier an den Befehlen selbst geprüft, nicht
  übernommen; der Verifier las die GitHub-Seite zu den Rechten nach, die der Review offen ließ. Die DoD-Zeile „keine Run-,
  Commit- oder Job-IDs als Beleg“ gab dem Verifier den Maßstab, die Abweichung V-1 zu finden. Der Tap-Kopf blieb von der
  Plananlage bis zur Closure `662c5fb5…`; ein Schreibzugriff fand nicht statt.
- **Was ging anders als geplant:** (1) Der Review-Fix zu F-1 setzte in den Absatz „Gelaufen und nicht gelaufen“ eine
  Commit-Kennung und die Aussage „Reviewer und Verifier haben ihn gelesen“. Beides widersprach der DoD-Zeile (keine
  Commit-IDs als Beleg), und die zweite war zur Zeit des Fix-Commits unbelegt. Ein Nachlauf-Review des Fix-Commits gab es
  nicht; der Verifier fand es (V-1, MEDIUM), berichtigt in `f902d127`, Nachprüfung im Nachtrag 2 des Verification-Reports.
  Der Hauptagent ordnet V-1 der Klasse
  [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) zu bzw. als Fehler
  seiner Vorgabe an den Fix (Angabe im Auftrag an den Planner, der Fix-Auftrag ist nicht gelesen). (2) Der Wortlaut der DoD
  (c) weicht in zwei Punkten vom Abschnitt ab: Schritt 3 nennt `branch`, §9 lässt ihn weg und weist es aus (der Push geht
  auf den Standard-Branch des Tap; Review F-4); Schritt 5 erkennt den Lauf an `created_at` und `head_sha`, §9 an
  `created_at`, weil `head_sha` bei zwei Dispatch-Läufen auf demselben Commit gleich ist (Verification V-4). Die DoD bleibt
  unverändert; maßgeblich ist der Abschnitt. (3) Nicht im Plan standen die Voraussetzung vor Schritt 3 (`git fetch`,
  `git status --short`, `git diff origin/main -- packaging/homebrew-tap/`, Review F-3) und der Hinweis auf `base64 -w0`
  als GNU-Option (F-2); beide stehen seit `8cc97257` in §9. Restgrenze der Voraussetzung (Verification): sie vergleicht nur
  `packaging/homebrew-tap/`, nicht den Stand von `HEAD` gegen `origin/main`. (4) Die Entscheidungen zu §4 lagen nur als
  Angabe des Hauptagenten vor (Review F-7, Verification V-3).
- **Steering-Loop-Eintrag und Lerneintrag:** *Benannte Lücke, jetzt gefüllt:* Das Verfahren eines Schritts in einem Fremd-Repo
  (`README.md` und `smoke.yml` aus `packaging/homebrew-tap/` in den Tap) stand nur im Auftrag und in der Closure des
  Vorgängers, **weil** kein Skript es trägt (`scripts/tap-nachzug.sh` klont den Tap und kopiert nur die Formel, Zeilen 25
  bis 27 und 34) und `packaging/homebrew-tap/README.md` nur das Anlegen des Tap-Repos nennt (Zeilen 19 bis 21); es steht
  jetzt in `docs/maintainer/releasing.md` §9. Die Drift-Gefahr bleibt als Grenze benannt und ist nicht behoben (Variante A,
  §6 Risiko 1). *Zweiter Befund, beobachtet:* Ein Fix, der auf einen Review-Befund antwortet, ändert den Text, an dem die DoD
  gemessen wird, **weil** Fix-Auftrag und DoD-Zeile nicht nebeneinander gelesen wurden (der Review verlangte einen Bezugspunkt,
  die DoD verbot Commit-Kennungen); ohne Nachlauf-Review über genau den Diff fing es erst der Verifier. Das ist eine
  Instanz der Klasse in
  [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Aussage über
  Verhalten und Zustand ohne gelesene Stelle, Zustandsmarker ohne Ablaufpunkt) und verwandt mit
  [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (Änderung
  nach dem Review; dort der Vertragstext, hier Maintainer-Doku; nicht gezählt). Es wurde nichts verkörpert, die
  3×-Schwelle ist für keine neue Regel erreicht; der Eintrag trägt kein Feld `liegt in …`. Ein neuer Sensor entsteht nicht:
  ein Gate, das eine DoD-Zeile gegen den Text eines Abschnitts liest, wäre eine Zusage über Prosa.
- **Beobachtungs-Register (`../observations/`; Zähler = Dateien in `evidence/`, gelesen am 2026-11-06):** neue Evidence-Datei
  `evidence/slice-v1-tap-nachzug-verfahren-doku.md` in
  [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (jetzt 91 Dateien;
  verkörpert in `AGENTS.md` §3.10, Rückfall; Passung eingeschränkt, die Aussage stand in der gelieferten Doku statt im Plan),
  `state.md` ergänzt. Kein neues Verzeichnis. Keine Datei und keine Zeile in `state.md`:
  [`BEO-WZG/runner-label-stillgelegt`](../observations/BEO-WZG/runner-label-stillgelegt/observation.md) (2×, offen; der Slice
  ist der Anlass, kein Auftreten),
  [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
  (nahe, nicht gezählt, §6 Risiko 4),
  [`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md)
  (Anwendung ohne Beleg: die Messstellen standen vorab im Plan; 3 Dateien, Ausgang beim Projektinhaber),
  [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) (Schritt 5 liest Jobs,
  Annotationen und Log mit `gh api`; der Verifier las so, ohne Fund),
  [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) („grün“ in §9 nennt
  Formel-Version und Image-Zeile),
  [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md)
  (Schritt 2 vor Schritt 3 in §9; Push ist per Schritt 6 umkehrbar) und
  [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (verwandt, nicht
  gleich). *Lese-Schritt:* kein Eintrag erreichte mit diesem Slice neu die Schwelle von drei; die Einträge über der Schwelle
  tragen ihren Ausgang in ihrer `state.md`, hier wird keiner neu zugewiesen.
- **Rest- und Folgepunkte** — Fragen an den Projektinhaber bzw. benannte Grenzen; keine Slice-Datei angelegt, ob einer als Slice
  entsteht, entscheidet er:
  1. **Die vier Entscheidungen zu §4 nachträglich.** Ort §9 (statt der eigenen Datei `docs/maintainer/homebrew-tap.md`),
     Sichtschritt in §4, Zeiger auf Variante A ohne neue Beobachtung, Push-Aufruf aus dem Sitzungsprotokoll: bestätigen oder
     anders entscheiden. Eine andere Antwort ändert `releasing.md` (§4-Satz) bzw. verlangt einen Folge-Slice (Ort).
     **Ziel: Projektinhaber.**
  2. **Die Punkte 2 und 3 des Vorgängers bleiben unverändert beim Projektinhaber:** ein regelmäßiger Lauf (`schedule`) und der
     hängende Lauf `42196563230` (bei der Closure weiter `queued`, kein Job). §9 Schritt 5 benennt den Fall (älterer
     Dispatch-Lauf ohne Job), nicht die Lauf-Kennung. **Ziel: Projektinhaber.**
  3. **Annotationen als Rot-Signal** (Review F-5, Verification V-5): ein grüner Lauf trägt eine Annotation der Stufe
     `failure` („Failed to import … vagrant.rb“); §9 nennt nur „not supported“ als Erkennungszeichen und sagt nicht, dass
     Annotationen allein kein Rot belegen. Nicht aufgenommen: eine Änderung an `releasing.md` nach der Verification ohne
     Nachlauf-Review wäre die Form von V-1. Der Hinweis gehört in die nächste Änderung von §9. **Ziel: keiner; benannt.**
  4. **Die Roadmap führt eine veraltete Zeile (Verification V-6):** Zeile 43 von `in-progress/roadmap.md` nennt
     [`slice-v1-bau-nur-devcontainer`](../done/slice-v1-bau-nur-devcontainer.md) als „liegt in `in-progress/`“, die Datei
     liegt in `done/`. Vorbestand, nicht im Diff dieses Slice und nicht berichtigt (Plan-Tabelle: Nettodiff der Roadmap
     leer). **Ziel: Planner beim nächsten Eingriff in die Roadmap.**
  5. **Grenzen, benannt:** die Modus-Deklaration in `harness/conventions.md` nennt `docs/maintainer/` nicht; die Zeile wird
     nicht erweitert (kein Modus ändert sich, `harness/conventions.md` steht nicht in der Plan-Tabelle); die Zeilenangaben der
     Lesebefunde in §9 (`smoke.yml`, `scripts/tap-nachzug.sh`, `scripts/tap-check.sh`) stehen ohne Sensor und verschieben sich
     mit einer Änderung dieser Dateien (Review F-6); die Sätze „Kein Skript, kein Gate und kein Test …“ stehen ohne Test
     (§6 Risiko 4). Die Neubewertung des Ortes steht in §1 („Neubewertung, wenn“). **Ziel: keiner; benannt.**
- **Folge-Slices:** keine angelegt. Die Punkte 1 bis 4 haben keine Slice-Datei, bis der Projektinhaber entscheidet.
- **Risiken aus §6:** je ein Ausgang, siehe §6 (eingetreten: keines; entfallen: 3, 5; weiter offen: 1, 2, 4).
- **Drei Paarungen:** (a) *Anker* — kein Herkunfts-Anker gesetzt, der Steering-Loop-Eintrag trägt kein Feld `liegt in …`, die
  Paarung hat keinen Gegenstand; (b) *Folge-Slice* — keiner genannt, die Paarung hat keinen Gegenstand; (c) *Register* — die
  genannten Kennungen `BEO-DOK/plan-aussage-ueber-code`, `BEO-DOK/vertragsaenderung-nach-review`,
  `BEO-WZG/runner-label-stillgelegt`, `BEO-DOK/spec-bedingung-ohne-code-entsprechung`,
  `BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`, `BEO-DOK/ci-log-falsch-gelesen`,
  `BEO-DOK/zusage-weiter-als-messung`, `BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle` und
  `BEO-WZG/kopplung-nicht-nachgezogen` existieren als Verzeichnis, jedes mit nicht leerem `evidence/` (geprüft nach dem Move).
- **Plan-Tabelle gegen `git diff --stat cf4c3cb8..HEAD`** (`AGENTS.md` §3.10; Stand nach der letzten Fix-Runde `f902d127`,
  Spitze nach dem Ruhe-Marker-Commit): 9 Dateien — `docs/maintainer/releasing.md`, `docs/maintainer/README.md`, diese Datei
  (neu, am Ruheort `done/`), `open/README.md`, der Review-Report, der Verification-Report (mit Nachtrag und Nachtrag 2),
  die neue Evidence-Datei und `state.md` von `BEO-DOK/plan-aussage-ueber-code/`. Die Plan-Tabelle führt jede davon; die
  Roadmap steht nicht in der Stat, weil der Start-Commit `0afd3272` den Ruhe-Marker entfernt und die Closure ihn an derselben
  Stelle wiederherstellt (Nettodiff leer, in der Tabelle so geführt). Plan-Einträge ohne Diff und ohne Behauptung des
  Gegenteils: die Zeilen „unverändert“ (Tests, `smoke.yml`, `packaging/homebrew-tap/README.md`, `scripts/tap-nachzug.sh`,
  `scripts/tap-check.sh`, `publish.yml`, `tap-check.yml`, `harness/README.md`, `CHANGELOG.md`,
  `docs/user/benutzerhandbuch.md`, `docs/maintainer/branch-protection.md`, [ADR-0016](../../adr/0016-homebrew-distribution-per-tap.md));
  ihre Stellen sind oben als gelesen vermerkt, soweit der Planner sie nachgelesen hat. Diff-Dateien ohne Plan-Eintrag: keine.
- **Carveouts:** keine neu, keine gelöst, keiner berührt (kein Eintrag unter `docs/plan/carveouts/` im Diff).
- **Verification und Gates:** [Review-Report](../../../reviews/2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md) (0 HIGH,
  1 MEDIUM, 2 LOW, 4 INFO; F-1 bis F-3 behoben in `8cc97257`, F-4 bis F-7 ohne Änderungspflicht);
  [Verification-Report](../../../reviews/verification-2026-11-06-slice-v1-tap-nachzug-verfahren-doku.md) (V-1 MEDIUM, behoben
  in `f902d127`, nachgeprüft im Nachtrag 2: Verdikt „DoD ‚Verfahren‘ erfüllt“; V-2 → §6 Risiko 2; V-3 → oben unter „Nur
  berichtet“; V-4 → „Was ging anders“; V-5 → Punkt 3; V-6 → Punkt 4). *Grenze:* zum Review-Fix `8cc97257` und zum
  Verification-Fix `f902d127` liegt kein eigener Review-Report vor; den Diff von `f902d127` hat der Verifier gelesen
  (Nachtrag 2), den von `8cc97257` hat er in der Spanne `8074e4cf^..8cc97257` mitgelesen. `make gates` im
  Devcontainer-Container durch den Verifier am Stand `f902d127`, Exit `0`: `d-check` 1212 Dateien 0 Befunde, `lint`
  0 issues, `coverage-gate` 92,40 % gegen 90 %, `baseline-verify` OK, Stempel `f902d127`. Nicht ausgeführt: `make test-docker`,
  `make verify-depguard`, `make ci`, `make fullbuild` (nur Markdown im Diff, kein Docker-, Importregel- oder Image-Pfad).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Der Slice berührt `docs/maintainer/` (Maintainer-Doku) und `docs/plan/planning/`;
beide liegen in der Sub-Area DOK (`docs/plan/`, `harness/`, `spec/`, `docs/messungen/` stehen in der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); `docs/maintainer/` steht dort
nicht wörtlich, es ist Dokumentation neben der Spec und wird der Zeile DOK zugerechnet; die Closure erweitert die Zeile nicht
(§7). Greenfield. `packaging/homebrew-tap/`, `scripts/` und `.github/workflows/` werden nur gelesen und
nicht geändert; keine eigene Sub-Area, die Tap-Dateien bleiben Aspirantin wie beim Vorgänger.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen am 2026-11-06 (alle Kürzel aufgelistet; die Einträge
unten gelesen, `state.md` und die Zahl der Dateien in `evidence/`).

- [`BEO-WZG/runner-label-stillgelegt`](../observations/BEO-WZG/runner-label-stillgelegt/observation.md) (2×, offen): trifft
  den Anlass, nicht den Gegenstand — die Stilllegung eines Labels erzwang den Push in den Tap, den das Verfahren beschreibt.
  Dieser Slice ist kein neues Auftreten.
- [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (90 Dateien,
  verkörpert in `AGENTS.md` §3.10): angewandt — jede Aussage in Anlass und §3 steht mit gelesener Stelle; der Abgleich gegen
  `git diff --stat` wird nach der letzten Fix-Runde wiederholt; die Namen der Reports in §3 tragen ein festes Datum (die
  Lücke aus V-1 des Vorgängers, `docs/reviews/verification-2026-11-06-slice-v1-tap-smoke-runner-macos-15.md`).
- [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)
  (46 Dateien, verkörpert in `AGENTS.md` §3.10): verwandt — die Regel gilt Spezifikation und Handbuch; die Maintainer-Doku
  trägt hier Sätze über Skripte (§6, Risiko 4). Kein Auftreten gezählt.
- [`BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle`](../observations/BEO-DOK/plan-annahme-ueber-fremdzustand-ohne-messstelle/observation.md)
  (3 Dateien, offen, Ausgang beim Projektinhaber): trifft — Rechte des Tokens und Verhalten der GitHub-API sind Fremdzustand.
  Die Messstellen stehen im Plan: Tap-Dateien und Lauf werden gelesen (Anlass, DoD), die Rechte des Tokens werden nicht
  behauptet (§6, Risiko 3), der Schreibschritt ist als ungeprobt ausgewiesen (§6, Risiko 2).
- [`BEO-DOK/ci-log-falsch-gelesen`](../observations/BEO-DOK/ci-log-falsch-gelesen/observation.md) (1 Datei, offen): trifft —
  Schritt 5 liest Jobs, Annotationen und Log mit `gh api`, nicht allein mit `gh run view --log`.
- [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (41 Dateien, über der
  Schwelle, offen): trifft — „grün“ nennt in §9 die installierte Formel-Version und die Image-Zeile, kein „macOS grün“.
- [`BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle`](../observations/BEO-DOK/unumkehrbarer-schritt-vor-der-kontrolle/observation.md)
  (1 Datei, offen): verwandt — Schritt 2 in §9 setzt Review und Verification der Quelle vor den Push in den Tap; der Push ist
  per Folge-Commit umkehrbar (Schritt 6), startet aber keinen Lauf, den startet erst Schritt 5. Keine Zählung. Der `state.md`
  der Beobachtung nennt `releasing.md` als Kandidaten-Ort für Verfahren, die nur im Plan stehen; dieser Slice schreibt dort
  keinen Advisory-Zyklus.
- [`BEO-WZG/kopplung-nicht-nachgezogen`](../observations/BEO-WZG/kopplung-nicht-nachgezogen/observation.md) (2 Dateien, offen):
  verwandt, nicht gleich — dort hängt ein Skript an der Form eines Erzeugers, hier steht dieselbe Datei an zwei Orten (§6,
  Risiko 1). Nicht gezählt.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF bzw. keine Sub-Area (Aspirantin).
