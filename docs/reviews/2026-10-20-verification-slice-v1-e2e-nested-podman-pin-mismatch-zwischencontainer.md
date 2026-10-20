# Verification-Report: slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer — 2026-10-20

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../plan/planning/in-progress/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (`in-progress/`), Commit-Range `1b4352c..3e5c2c5` (drei Commits: `034e124` Umsetzung, `acc095a` Review-Report, `3e5c2c5` Kommentar-Fix zu Review F-1). Vorgänger: [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../plan/planning/done/slice-v1-e2e-pin-mismatch-zwischencontainer.md). Anforderung: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer); Entscheidung: [`ADR-0023`](../plan/adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) (Negativprobe zu Prüfsumme und Signaturschlüssel, durch den Slice nicht berührt).

**Eingang:** Implementer-Zusage (Abschnitt Umsetzung im Slice), Review [`2026-10-20-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](2026-10-20-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (F-1 bis F-6, kein HIGH/MEDIUM), Fix-Commit `3e5c2c5`.

**Gesamturteil: GO mit Auflagen.** Der Test-Eingriff stimmt mit Plan, DoD-Punkten 1 und 2 und dem gelesenen Code überein; die Mechanismus-Tabelle ist an Dockerfile-Vorlage und Test nachgerechnet; der Wächter gegen den vakuum-grünen Lauf ist am Code gelesen und deckt die Annahme, dass der Builder der klassische ist. Offen ist nichts am Produkt oder am Testcode, sondern die Closure: vier Risiken in §6 stehen noch auf „weiter offen bis zur Umsetzung“, die Ollama-Hälfte des buildx-Risikos (Review F-2) hat keinen Ausgang, das Register trägt den zweiten Beleg nicht, §7 und der dritte DoD-Punkt sind leer, Roadmap und `open/README.md` tragen überholte Aussagen.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass | eigener Lauf auf `3e5c2c5`, sauberer Arbeitsbaum, Exit 0; `coverage-gate: OK — coverage 91.40% meets threshold 90%`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege` |
| `make docs-check` | in `make gates` enthalten; für diesen Report nach dem Schreiben erneut (siehe Abschnitt 8) | — |
| `make test-docker` (voll) | **nicht gefahren**, auf Auftrag | Der Implementer meldet Exit 0 (`internal/e2e` 655,7 s, `NestedPodman_BuildEgress` grün, Container-ID-Liste vorher/nachher gleich). Diese Angabe ist hier **nicht belegt**, nur übernommen. |
| Test gezielt (`-run NestedPodman_PinMismatch`) | **nicht gefahren**, auf Auftrag | Beleg dafür sind der Implementer-Lauf und die Mutationsläufe des Reviews (Review-Report, Schwerpunkt 1: geliefert PASS in 30,9 s; vier Mutationen); nicht selbst wiederholt. |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (`git diff --stat 1b4352c..HEAD`: Test, Slice, Review-Report) |

Der Build-Tag `docker` liegt auf dem Test; `make gates` (Lint, Unit-Tests) übersetzt die Datei nicht. Der Kommentar-Fix `3e5c2c5` (nur Kommentar, 9 Zeilen, selbst gelesen) ist daher von keinem Lauf abgedeckt, den ich gefahren habe; der volle Implementer-Lauf liegt vor diesem Commit (Reihenfolge der Commits: Umsetzung, Review, Fix).

## 2. Eingriff am Code gelesen (Plan-vs-Code)

Test selbst gelesen (`internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`, Helfer `containerIDs`, `newContainersWithCommand`):

- Je Teilfall: Vorher-Liste (`containerIDs`), Bau mit `--force-rm` und ohne `-q`, `t.Cleanup` am Subtest-`t` mit `newContainersWithCommand(…, before, tc.pin)`; ein Treffer wird mit `rm -f` entfernt und als `t.Errorf` gemeldet (der Test wird rot, auch wenn das Entfernen gelingt). Der Helfer ist laut Diff unverändert.
- Zuordnung nur über „nicht in der Vorher-Liste“ **und** Pin-String in `Path`/`Args`/`Env`; kein Status, kein Zeitfenster, kein Name. Die Pins sind der volle `NAME=Wert`-String (Nullen bzw. der Subkey-Fingerabdruck), entstehen im Paket nur in diesem Test. Plan-Zusage (§1, §6 Risiko 1) und Code stimmen überein, einschließlich der Grenze gleichzeitiger Läufe (Doc-Kommentar des Helfers).
- Wächter `failedStepRemoved` (`Removed intermediate container <hex>`, Zeilenumbruch, `The command '`): verlangt in jedem Teilfall die Entfernung des Containers des **fehlgeschlagenen** Schritts, nicht nur irgendeines Zwischenschritts (Nachbarzeile). Ohne Treffer ist der Teilfall rot mit einer Meldung, die „kein klassischer Builder“ und „`--force-rm` fehlt“ nennt. Das schließt den vakuum-grünen Lauf für dieses Gegenstück, wie §6 Risiko 2 es für den Test verlangt.
- Kommentare nach `AGENTS.md` §3.7/§3.9 gelesen: Der Kommentar nach `3e5c2c5` steht im Indikativ („The leftover check above finds a container only if the classic builder created one.“), kein „Ohne X“, keine Chronik; `failedStepRemoved` steht vor dem Doc-Kommentar der Testfunktion, nicht dazwischen. Review F-1 ist behoben.

## 3. Gelesene Aussagen im Slice nachgerechnet

| Aussage im Slice | Nachrechnung | Urteil |
|---|---|---|
| Alle drei Teilfälle scheitern im selben `RUN`-Schritt, „Step 7/12“ | `Dockerfile.tmpl` (Podman-Zweig): `FROM` (1), `USER root` (2), vier `ARG` (3 bis 6), das `RUN apt-get …` (7). Für ein Projekt ohne Egress, mit Sandbox: danach `COPY sandbox-init.sh` (8), `RUN mkdir … status` (9), `RUN mkdir … workspace` (10), `RUN git config` (11), `USER vscode` (12). Gesamtzahl 12 passt, der Schritt 7 ist der `RUN` mit Prüfsumme und Signatur. Mit `Egress` wären es 14; `generateSandboxProject(…, "podman", "warn")` erzeugt demnach keinen Egress-Zweig — die 12 hängen an dieser Annahme (aus der Zählung geschlossen, nicht am erzeugten Dockerfile gelesen). | stimmt, Annahme benannt |
| Exit 1 (wrong checksum), Exit 1 (Subkey), Exit 22 (unbekannter Fingerabdruck) | `sha256sum -c -` bei falscher Prüfsumme: 1; Subkey als Pin: der Schlüssel mit diesem Fingerabdruck ist abrufbar, der `VALIDSIG`-Vergleich (`test … = "${PODMAN_SIGNING_KEY}"`) schlägt fehl: 1; Nullen-Fingerabdruck: `curl -fsSL` auf `…/by-fingerprint/000…` liefert 404, `curl -f` endet mit 22. Der Review hat die Ausgabe `returned a non-zero code: 1` (2×) und `22` mit `curl: (22) … 404` gesehen. Die Reihenfolge der Befehle im `RUN` ist mit der Tabelle konsistent (Prüfsumme vor Schlüssel). | stimmt aus Code und Review-Beleg; nicht selbst gemessen |
| Pin steht im Env des zurückbleibenden Containers | `ARG` wird im klassischen Builder in die Umgebung des `RUN`-Containers übernommen (so lief schon der Vorgänger mit `OLLAMA_SHA256_AMD64=` im Env); die Mutation ohne `--force-rm` fand im Review in allen drei Teilfällen den Container über genau diesen String. | durch Review belegt |
| Messvariante (Test-Binary mit vorübergehender Ausgabe von `docker inspect`, ohne `--force-rm`) ≠ Auslieferung | Die Messvariante ist wieder entfernt und nirgends committet; der Diff enthält die Ausgabe nicht. Die Auslieferung wurde vom Review in Kopien mutiert (`--force-rm` entfernt, `--no-cache`, Muster verfälscht). Dass der **committete** Test läuft: Implementer-Lauf nach `034e124`, Review-Lauf auf `034e124`; `3e5c2c5` ändert nur Kommentar. Kein Code, der nur in der Messvariante stand, steht in der Auslieferung. | getrennt und sauber; der Lauf auf `3e5c2c5` fehlt (Abschnitt 1) |
| Plan-Tabelle gegen Diff | `git diff --stat 1b4352c..HEAD`: Test, Slice-Datei, Review-Report. Jede Datei steht in §3 (Test: Zeile 1; Slice: Sammelzeile; `docs/reviews/`: Sammelzeile). `quality.md`, `docs/messungen/`, `tools/messung/`, `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, Workflows, `spec/` ohne Diff, wie ausgewiesen. Der Abgleich im Slice trägt „Stand der Umsetzung, vor Review“; er muss nach der letzten Fix-Runde und mit Verification-Report und Closure-Dateien wiederholt werden (`AGENTS.md` §3.10). | stimmt für den Stand; Wiederholung offen (V-4) |

## 4. Zugesagt nicht weiter als gemessen

| Zusage | Beleg | Urteil |
|---|---|---|
| Klassischer Builder legt den Container an, `--force-rm` entfernt ihn (Host, Docker 29.9.1, `amd64`) | gemessen (Implementer), vom Review in drei Mutationsläufen wiederholt | gedeckt, auf diesen Host begrenzt |
| BuildKit: Test wird rot, nicht vakuum-grün | Slice sagt „abgeleitet, nicht mit BuildKit gesehen“; Review hat auf dem Host den Wortlaut gegengeprüft (0 Zeilen „Removed intermediate“ bei `--progress=plain`). Das Test-Image wurde nie auf BuildKit umgestellt. Von mir **nicht** nachgefahren (Auftrag und keine fremden Docker-Objekte). | ehrlich ausgewiesen; Bedingung der Aussage ist der Wortlaut des Builders (Review F-5) |
| CI-Runner | ungemessen; Slice §6 Risiko 3 und „Nicht gemessen“ sagen es | gedeckt, Ausgang offen |
| Zwei gleichzeitige Läufe gegen einen Daemon | nicht erzeugt; Slice §1 (Ziel), §6 Risiko 1, Doc-Kommentar des Helfers, „Nicht gemessen“ nennen die Grenze | gedeckt |
| „Hinterlässt nach seinem Lauf keinen Container“ (Ziel, §1) | gilt unter der Bedingung des klassischen Builders und eines Daemons ohne fremde, gleichzeitige Läufe; unter BuildKit wird der Test rot statt grün. Die Formulierung steht ohne die Builder-Bedingung im Ziel, die Bedingungen stehen in §1/§6. | tragbar; in der Closure-Notiz mit der Bedingung wiederholen |

Die Vorher/nachher-Aussage des Implementers (3 vorher, 3 nachher, gleiche Prüfsumme der ID-Liste) ist auf einem geteilten Host kein Beleg für Wirkung (der Vorgänger hielt das schon fest); der Beleg ist die Prüfung nach Pin und neuer ID. Ich habe sie nicht wiederholt.

## 5. Review-Finding F-2: Ollama-Test ohne Wächter

**Befund am Code:** `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `-q --force-rm`, prüft nur `FAILED` in der Ausgabe und die Rest-Container; es gibt keinen Wächter. Unter BuildKit bliebe er grün (Rest-Prüfung findet nichts, `FAILED` kann in der Fehlerausgabe stehen), der Nested-Test wird rot. Das ist die Ollama-Hälfte desselben Risikos, das der Slice in §6 führt.

**Verlangt der Slice einen Ausgang?** Ja, einen Ausgang, aber keinen Code-Eingriff:

- §6 Risiko 2 verlangt für den Slice „entweder eine Prüfung … oder die Lücke als Beobachtung im Register“. Die Prüfung ist für den Test dieses Slice geliefert. Für Ollama sieht §1 ausdrücklich keine Änderung vor („Das Ollama-Gegenstück — … hier nicht angefasst“). Der Vorgänger hat das Risiko „Wechsel auf buildx macht die Prüfung vakuum-grün“ in seiner Closure ausdrücklich an diesen Slice übergeben (Vorgänger §7, Steering-Loop-Eintrag (3), und §6). Damit gehört der Ausgang für **beide** Prüfungen hierher.
- Es ist deshalb **kein Befund am gelieferten Code** und kein Verstoß gegen §1, aber ein **offener Ausgang des Slice**, ohne den das DoD-Wort „jedes Risiko aus §6 trägt einen Ausgang“ nicht erfüllt ist.

**Was bei der Closure stehen muss:**

1. §6 Risiko 2: Ausgang geteilt und belegt. Nested-Podman: getragen durch `failedStepRemoved` (Fundstelle Datei und Symbol, Mutation „Muster ohne Treffer“ rot in allen drei Teilfällen). Ollama-Test: weiter offen → Register `docker-test-hinterlaesst-zwischencontainer` (Wächter fehlt, `-q` im Bau; unter BuildKit bliebe er grün), mit Folge-Punkt, nicht eigenmächtig geschlossen.
2. Register: der Ollama-Befund als Teil der Beobachtung benennen (Abschnitt „Benannt, nicht gezählt“ in `observation.md`) und entscheiden lassen, ob ein Folge-Slice entsteht (Ollama baut ohne `-q`, übernimmt den Wächter; gemeinsamer Helfer erst bei drittem Vorkommen laut §1 Nicht-Ziele). Die Entscheidung liegt beim Planner/Projektinhaber; der Verifier schlägt keinen Test vor, der den Slice-Umfang ändert.
3. Dasselbe Risiko gehört in die Closure-Notiz als „Zur Entscheidung vorgelegt“, nicht als Verkörperung (Zählung 2×, unter der Schwelle von 3×).

## 6. DoD-Abgleich

- [x] **Mechanismus gemessen und Test-Eingriff:** erfüllt. Messung je Teilfall im Abschnitt Umsetzung (Tabelle), nachgerechnet in Abschnitt 3; Eingriff mit Vorher-Liste, Pin je Teilfall und `--force-rm`, entfernt nur neu entstandene Container mit dem Pin; Grenze gleichzeitiger Läufe genannt. Vorbehalt: Messvariante nicht committet, Messrohdaten nicht ablegbar (Plan §3 verzichtet bewusst auf `docs/messungen/`).
- [x] **Pin des Eingriffs:** erfüllt. Die Prüfung steht im selben Docker-Testlauf, deckt alle drei Teilfälle, wurde rot gesehen (Implementer-Mutationen; vom Review wiederholt). Vom Verifier **nicht** erneut rot gesehen (Abschnitt 1, Auftrag); das Rot beruht auf dem Review.
- [ ] **Dritter Punkt** (`make gates`, Review, Report, Closure-Notiz, Register, §6-Ausgänge, Paarungen): `make gates` grün (eigener Lauf, Abschnitt 1); Review-Report liegt vor; Verification-Report mit diesem Dokument; **offen:** der Haken selbst, die Closure-Notiz (siehe V-6 zum vollen Lauf), Register, §6-Ausgänge, Paarungen (Abschnitt 7).

## 7. Befunde

| ID | Schwere | Befund | Fundstelle |
|---|---|---|---|
| V-1 | MEDIUM | §6 trägt für alle vier Risiken den Ausgang „weiter offen bis zur Umsetzung“ (drei Mal) bzw. „weiter offen“; die Umsetzung ist geschehen, die Marker sind überholt (`AGENTS.md` §3.10, Zustandsmarker mit Ablaufpunkt). Ausgänge fehlen, die DoD verlangt sie. Sachlich: Risiko 1 (fremder Container) getragen mit Grenze; Risiko 2 (buildx) geteilt, siehe Abschnitt 5; Risiko 3 (CI-Runner) weiter offen → Register; Risiko 4 (Pin nicht im Env jedes Teilfalls) **entfallen**, Messung in der Tabelle der Umsetzung. | Slice §6 |
| V-2 | MEDIUM | Ausgang der Ollama-Hälfte (Review F-2) fehlt, siehe Abschnitt 5: Register-Eintrag und Folge-Punkt sind Auflage der Closure, kein Eingriff am Ollama-Test in diesem Slice. | Slice §6 Risiko 2, Register |
| V-3 | MEDIUM | Register nicht fortgeschrieben: `state.md` führt „offen — 1×“, `observation.md` führt das Nested-Vorkommen als „gelesen, nicht gemessen“, und `evidence/` hat nur die Datei des Vorgängers. Mit der Messung ist das zweite Vorkommen belegt (drei Teilfälle, Mechanismus gleich); die Evidence-Datei für diesen Slice (Beobachtungsform) und die Zählung 2× sind zu schreiben, die „Benannt, nicht gezählt“-Passage und der Trigger der Beobachtung sind auf den Ist-Stand zu setzen. Die 3×-Schwelle ist nicht erreicht; nichts wird verkörpert. | `docs/plan/planning/observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/` |
| V-4 | LOW | Plan-Tabelle gegen `git diff --stat`: Der Abgleich im Slice nennt „Stand der Umsetzung, vor Review“. Nach der letzten Fix-Runde (`3e5c2c5`), mit Review-Report, Verification-Report, Register, Roadmap, `open/README.md` und den Verweis-Nachzügen des Moves ist er zu wiederholen; die Plan-Zeilen „Register“ und „Roadmap“ brauchen die konkreten Pfade, wie es der Vorgänger in seinem Endstand tat. Auch der Satz „Sensoren (Implementer, Stand nach der letzten Änderung am Test)“ trifft nach `3e5c2c5` nicht mehr zu: Der Test wurde danach (nur Kommentar) geändert. | Slice, Abschnitt Umsetzung |
| V-5 | LOW | Überholte Aussagen in Planungs-Indizes: (a) `roadmap.md` Z. 20 „In Arbeit: …“ (Zustandsmarker, bei Closure entfernen); (b) `open/README.md` Z. 32 sagt „erst messen, dann dieselbe Lösung oder eine begründete andere; … Start nach dem Vorgänger“ — der Slice ist gestartet, gemessen und liegt in `in-progress/`; die Zeile wird bei der Closure auf „geliefert (liegt in `done/`)“ gesetzt wie Z. 31 des Vorgängers, deren Halbsatz „das zweite Vorkommen führt [Slice] (../in-progress/…)“ nach dem Move auf `../done/` zeigt (der Move-Werkzeug-Nachzug zieht Verweise nach, danach prüfen); (c) Vorgänger-Slice §6 und §7 und `observation.md` verlinken den Slice unter `in-progress/`, dieselbe Nachziehpflicht. | `docs/plan/planning/in-progress/roadmap.md` Z. 20; `docs/plan/planning/open/README.md` Z. 31, 32 |
| V-6 | LOW | Slice-interne Zustandsmarker: §1 „Nicht gemessen: ob die drei Teilfälle Zwischencontainer hinterlassen … (vermutet, nicht belegt)“ ist durch die Umsetzung gemessen; als datierter Fundstand („gelesen am 2026-10-19“) zulässig, in der Closure-Notiz mit dem Ist-Stand zu ergänzen. Das „ab Zeile 572“ stimmt nur für den Stand vom 2026-10-19 (Review F-4). DoD-Punkt 3 verlangt „der berührte Docker-Test gezielt grün (der volle `make test-docker` nur, wenn …)“; der Implementer lieferte den vollen Lauf grün — die Closure-Notiz soll sagen, dass der volle Lauf **gemeldet**, vom Verifier nicht wiederholt ist. | Slice §1, §2 |
| V-7 | INFO | `make gates` übersetzt die Datei mit Build-Tag `docker` nicht; der Kommentar-Fix `3e5c2c5` ist von keinem Lauf des Verifiers getragen. Kommentar-Diff selbst gelesen, kein Code geändert. | `internal/e2e/sandbox_devcontainer_docker_test.go` |
| V-8 | INFO | Review F-5 (Wortlaut-Kopplung des Wächters an den klassischen Builder, „DEPRECATED“-Hinweis) und F-6 (Aufräumblock zweimal fast wortgleich, Hilfe erst beim dritten Vorkommen laut §1): beide als Grenze erklärt; Closure-Notiz nennt sie als bekannte Folge, nicht als Mangel. | Review F-5, F-6 |

## 8. Auflagen vor Closure

1. §6: für jedes der vier Risiken einen Ausgang mit Beleg eintragen (V-1): fremder Container — getragen mit Grenze (Doc-Kommentar des Helfers, Slice §1); buildx — geteilt getragen/offen (V-2, Abschnitt 5); CI-Runner — weiter offen → Register; Pin nicht im Env — entfallen (Tabelle der Umsetzung, drei Teilfälle mit Pin im Env).
2. Register fortschreiben (V-3, V-2): eine Evidence-Datei unter `evidence/` mit der Slice-Kennung als Dateinamen, wie beim Vorgänger (Beobachtungsform), `state.md` 2× mit Stand „offen“ und dem Trigger für das dritte Vorkommen, `observation.md` mit der Ollama-Lücke und dem CI-Runner. Das Nested-Vorkommen verkörpert nichts (2× unter 3×); die Ollama-Entscheidung (Folge-Slice oder nicht) gehört in „Zur Entscheidung vorgelegt“, nicht eigenmächtig in `AGENTS.md`.
3. Plan-Tabelle und Abgleich gegen `git diff --stat` für den Endstand wiederholen (V-4), mit Pfaden für Register und Roadmap.
4. §7 schreiben: Gegenstand und Commits, Belege (Review: 0 HIGH, 0 MEDIUM, 1 LOW behoben durch `3e5c2c5`; dieser Report), Läufe mit Exit-Codes, **Nicht gelaufen und Grenzen** (voller `make test-docker` vom Implementer gemeldet, vom Verifier nicht wiederholt; Test gezielt vom Verifier nicht gefahren; CI-Runner und BuildKit-Fall nur abgeleitet beziehungsweise auf dem Host gegengeprüft; gleichzeitige Läufe; `make verify-depguard` nicht berührt), Steering-Loop-Eintrag (zwei Vorkommen, Mechanismus-Wächter als Test-Voraussetzung ausgesprochen, kein neuer Sensor), Folge-Punkte, drei Paarungen (Anker: kein `liegt in`-Feld; Folge-Slice: nur falls entschieden, sonst „keiner“; Register: Verzeichnis mit nicht leerem `evidence/`).
5. Dritter DoD-Punkt: Haken erst nach 1 bis 4, mit Begründung für den vollen `make test-docker` (gemeldet, nicht wiederholt) und dem Commit-Hash des Inhalts.
6. Roadmap-Zeile „In Arbeit“ entfernen, `open/README.md` Z. 32 (und Z. 31) auf den Ist-Stand setzen (V-5). Closure-Reihenfolge nach `AGENTS.md` §3.3 Fall 2: erst Inhalt (DoD-Haken, Closure-Notiz mit Commit-Hash), dann der reine `git mv` nach `done/` (das Move-Werkzeug), danach die Verweis-Nachzüge prüfen, `make gates`.

## 9. Nicht geprüft

- Voller `make test-docker` und der Test `-run NestedPodman_PinMismatch`: auf Auftrag nicht gefahren; Ergebnisse stammen vom Implementer und vom Review.
- Die roten Mutationen nicht selbst wiederholt (Quelle: Review-Report, Schwerpunkt 1).
- BuildKit-Verhalten, Docker Desktop, der Daemon des CI-Runners, zwei gleichzeitige Läufe gegen einen Daemon: nicht erzeugt, im Slice als nicht gemessen ausgewiesen.
- Der erzeugte Dockerfile-Text des Testprojekts (Step-Zählung 7/12 aus der Vorlage gerechnet, nicht am erzeugten Projekt gelesen).
- Ob die Host-Container-Liste vor und nach dem Implementer-Lauf gleich war (Angabe des Implementers; kein Beleg am geteilten Host).
- Handbuch, Spezifikation und Lastenheft: nicht berührt (`git diff --stat`), nicht gelesen.

## 10. Carveouts und Artefakt

Carveouts: neu none, gelöst none, unverändert (Inventar im Diff leer; Test-Hygiene ohne Gate-Änderung, `AGENTS.md` §3.6 nicht berührt; kein `//nolint`, keine Docker-Befehle auf dem Host als Nachweis außer den `make`-Läufen).

Verifizierter Stand: `3e5c2c5`. Dieser Report liegt unter `docs/reviews/` (Verification-Report des Slice).
