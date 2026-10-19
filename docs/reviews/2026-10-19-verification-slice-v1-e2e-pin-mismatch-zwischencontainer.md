# Verification-Report: slice-v1-e2e-pin-mismatch-zwischencontainer — 2026-10-19

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../plan/planning/in-progress/slice-v1-e2e-pin-mismatch-zwischencontainer.md) (`in-progress/`), Commit-Range `5eab741..HEAD` (HEAD `b37c715`). Implementierung: `1689c68` (Test), `be66a0a` (Slice-Stand). Review-Report `7422493`. Fix-Runde: `18f4edc` (Test), `b37c715` (Slice). Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md).

**Eingang:** Implementer-Zusage (DoD, Plan-Tabelle, Messtabelle), Review [`2026-10-19-slice-v1-e2e-pin-mismatch-zwischencontainer`](2026-10-19-slice-v1-e2e-pin-mismatch-zwischencontainer.md) (0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO).

**Gesamturteil: GO mit Auflagen.** Die Zusage trägt: Der Docker-Test ist von mir gefahren und grün, die Mutation färbt ihn rot, der Mechanismus ist mit eigener Probe bestätigt, nach dem Lauf bleibt kein Container mit dem Null-Pin. Der Diff berührt nur den Test, die Slice-Datei und den Review-Report. Eine DoD-Verletzung (Zusage ohne Beleg) liegt nicht vor. Auflagen sind ausschließlich Closure-Arbeit des Planners (Abschnitt 10).

## 1. Sensoren (selbst gefahren, nur `make` beziehungsweise das Image des `make test-docker`-Ziels)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | `coverage-gate: OK — coverage 91.40% meets threshold 90%`, `baseline-verify: v6.18.1 OK — 54 Dateien` |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n).“ |
| `make docs-check` | pass | 516 Dateien, 0 Befunde (vor Anlage dieses Reports) |
| `make doc-planning` | pass | 516 Dateien, 0 Befunde |
| `make ablage-pruefung` | pass | „26 Faelle, 0 Fehlschlaege“ |
| `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` (Build-Tag `docker`) | **PASS** | Abschnitt 2 |
| Mutation (`--force-rm` entfernt) | **FAIL wie erwartet** | Abschnitt 3 |
| `make test-docker` (voll) | nicht gelaufen | dauert lange und ist, wie im vorigen Verifier-Bericht dokumentiert, bekannt rot durch einen anderen, nicht berührten Test; der berührte Test lief gezielt im Image des Ziels |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (nur eine `_test.go`-Datei mit Build-Tag `docker`; `make lint` in `make gates` grün) |

`git status` sauber vor dem Lauf; `git diff --stat -M 5eab741..HEAD`: 3 Dateien (Test 49 Einfügungen und 2 Löschungen, Slice, Review-Report), 172 Einfügungen, 9 Löschungen. `git diff --stat 5eab741..HEAD` über `internal/hexagon`, `internal/adapter`, `spec`, `.github`, `repo.mk`, `Makefile`, `.d-check.yml`, `docs/plan/adr`, `.harness`: **leer**. Produktcode, Templates, Lastenheft, Spezifikation, ADRs, Workflows unverändert.

## 2. Docker-Test selbst gefahren

Aufruf wie das Ziel `test-docker` in `repo.mk`: erst der Bau-Schritt (`--target test-docker-tools`, Exit 0), dann im Image mit denselben Mounts (Host-Netz, Docker-Socket, `/tmp`) `go test -tags docker -timeout 30m -count=1 -v -run '^TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild$' ./internal/e2e/`.

Ergebnis: **`--- PASS: TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild (32.46s)`**, `ok … 32.464s`, Gesamtdauer des Aufrufs 45 s. Weder SKIP noch FAIL.

Container-Zählung (nur Status `exited` und Env-Treffer auf `OLLAMA_SHA256_AMD64=` plus 64 Nullen, per `docker inspect` über Path, Args und Env; keine Namen ausgegeben):

| Zeitpunkt | Container gesamt | Status `exited` | davon mit Null-Pin | Container mit Null-Pin in irgendeinem Status |
|---|---|---|---|---|
| vorher | 3 | 0 | 0 | 0 |
| nachher | 3 | 0 | 0 | 0 |

**Nach dem Lauf bleibt kein Container mit dem Null-Pin im Env zurück.** Es war nichts zu räumen; ich habe nichts entfernt.

## 3. Wirksamkeit unabhängig (Mutation) und Mechanismus

Wegwerf-Kopie des Repos (`git archive HEAD`) im Scratchpad, nicht im Arbeitsbaum; einzige Änderung gegenüber HEAD: `"--force-rm"` aus dem `docker build`-Aufruf entfernt (per `diff` gegen HEAD bestätigt, eine Zeile). Gezielter Lauf im selben Image:

```
sandbox_devcontainer_docker_test.go:1030: the failed build left container 69364ba755b7 with the zero pin; cleanup removed it
--- FAIL: TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild (31.08s)
```

Das ist die erwartete Rot-Sicht, und der Cleanup hat den Container selbst entfernt: Container-Zählung vorher 3, nachher 3, kein Container mit Null-Pin übrig. Ich musste nichts nachräumen. Damit ist belegt: Ohne `--force-rm` entsteht der Zwischencontainer (die Prüfung findet ihn genau einmal), mit `--force-rm` nicht (Abschnitt 2); die Prüfung ist nicht vakuum-grün.

Mechanismus, eigene Probe im Image `test-docker-tools`: `docker buildx version` antwortet `docker: unknown command: docker buildx`, Exit 1; `DOCKER_BUILDKIT` ist leer. Der Bau im Test läuft damit über den klassischen Builder; die Aussage des Slice ist bestätigt. Nicht gemessen bleibt (vom Slice so ausgewiesen, §6 Risiko 2): der Daemon des CI-Runners und ein Client mit buildx.

## 4. Zusage begrenzt, Cleanup, Fremdnamen

- **Zusage „nur neu entstandene Container mit Null-Pin; zwei gleichzeitige Läufe nicht ausgeschlossen“:** konsistent im Code-Kommentar von `newContainersWithCommand` („Two runs against the same daemon at the same time are not told apart“), in §1 (Ziel), in der DoD (erster Punkt) und in §6 (Risiko 1). LOW 1 des Reviews: **behoben** (Zusage auf den Code zugeschnitten).
- **Cleanup unterscheidet beide Ausgänge:** gelesen im Diff von `18f4edc`: scheitert `docker rm -f`, meldet der Test „removing it failed: …“ und fährt fort; gelingt es, „cleanup removed it“. Der Zweig „Entfernen scheitert“ ist von mir nicht ausgelöst (kein Test dafür, nur gelesen); der Zweig „Entfernen gelungen“ ist durch die Mutation gesehen (Abschnitt 3).
- **Fremdnamen:** Suche im Baum nach den vorgegebenen Fremdnamen (Präfix- und Image-Namen fremder Läufe): **keine Treffer** im Arbeitsbaum (ohne `.git`). In der Slice-Datei steht das generische Basis-Image der Wegwerf-Probe nur noch als Wort im Abschnitt Umsetzung (Mechanismus); der Test selbst benutzt dasselbe öffentliche Basis-Image an anderen Stellen. Ich werte das als kein fremdes Objekt (Review INFO 4); wer strikt sein will, kürzt das Wort zu „ein Basis-Image“.

## 5. Plan-vs-Diff (`AGENTS.md` §3.10), nach dem letzten Commit

| Plan-Zeile (§3) | Diff | Befund |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` (update) | berührt (`TestE2E_..._PinMismatchBreaksBuild`, Helfer `containerIDs`, `newContainersWithCommand` nach der Testfunktion) | stimmt |
| `docs/messungen/<slice>/`, `tools/messung/<name>/` (entfällt) | nicht im Diff | stimmt; Messung als Text im Slice |
| `docs/maintainer/quality.md` (unverändert) | nicht im Diff | stimmt (nicht neu gelesen; die Aussage des Plans steht mit Trefferzeilen) |
| diese Slice-Datei, `open/README.md`, Register, Roadmap, `docs/reviews/` | Slice-Datei und Review-Report im Diff; `open/README.md`, Register, Roadmap nicht | stimmt für den Stand vor der Closure; die Bewegung ist Closure-Arbeit |
| `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, `.github/workflows/`, `spec/` (unverändert) | nicht im Diff (Abschnitt 1) | stimmt |

Jede Diff-Zeile hat einen Plan-Eintrag, jeder Eintrag mit Diff eine Zeile. Mit diesem Report kommt eine vierte Datei hinzu (`docs/reviews/`, durch die Sammelzeile gedeckt). Die Aussage des Slice zum Abgleich („drei Dateien“) gilt für den Stand vor diesem Report und ist bei der Closure auf den Endstand zu setzen.

## 6. DoD, Punkt für Punkt

| DoD-Punkt | Befund |
|---|---|
| Mechanismus gemessen und Test-Eingriff (`[x]`) | **erfüllt, bestätigt.** Klassischer Builder im Image ohne buildx (eigene Probe), `--force-rm` im Aufruf, Prüfung nach neuer ID plus Null-Pin im Env; Zusage begrenzt formuliert. |
| Pin des Eingriffs (`[x]`) | **erfüllt, bestätigt.** Die Prüfung steht im selben Testlauf; einmal rot gesehen (Mutation, von mir wiederholt), danach grün (Abschnitt 2). |
| `make gates` grün, `make test-docker` grün | `make gates` von mir grün. `make test-docker` (voll) nicht von mir gefahren, der berührte Test gezielt grün; der Slice führt den Punkt offen (`[ ]`). Entscheidung bei der Closure (Auflage 1). |
| Review durchgeführt, Report liegt vor | erfüllt (`7422493`). |
| Closure-Notiz mit Steering-Loop-Lerneintrag; Register fortgeschrieben oder „keine Beobachtung“ notiert; Ausgang je Risiko in §6; die drei Paarungen | **offen**: §7 trägt „Wird bei der Closure gefüllt“, Risiko 1 und 2 in §6 tragen „weiter offen“. Planner-Aufgabe. |

## 7. Review-Findings einzeln

| Finding | Stand | Beleg |
|---|---|---|
| LOW 1: Zuordnung trennt gleichzeitige Läufe nicht | **behoben** (Zusage auf den Code zugeschnitten) | Code-Kommentar, §1, DoD, §6 nennen die Grenze; `18f4edc`, `b37c715` |
| LOW 2: Plan-Aussage Path/Args | **behoben** | In der Slice-Datei nennen nur noch Messaussagen Path/Args, richtig (Wert im Env, Variablenname im Skripttext); DoD und §6 tragen „Pin in Path/Args“ nicht mehr (Suche in der Datei) |
| INFO 1: zweites Vorkommen, Folge-Vorgang ohne Anker | **übergeben** an den Planner | Abschnitt 8, Auflage 2 |
| INFO 2: Mechanismus teils ungemessen | **begründet offen** | Klassischer Builder auf dem Host bestätigt (Abschnitt 3); CI-Runner-Daemon und Client mit buildx nicht gemessen, vom Slice ausgewiesen; bei einem Wechsel auf buildx würde die Prüfung vakuum-grün, ohne dass ein Test es sagt (Auflage 3) |
| INFO 3: Beleg-Tabelle, Fremdname | **behoben** bis auf das Wort für das Basis-Image (Abschnitt 4) | `b37c715` |
| INFO 4: Altbestand, Rückgabewert von `docker rm -f` | Rückgabewert **behoben**; Altbestand früherer Läufe **begründet offen** (gewollt) | `18f4edc`; Slice §1 grenzt das Aufräum-Skript aus |

## 8. Zweites Vorkommen (`TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`)

Gelesen: drei Teilfälle („wrong checksum“ mit `PODMAN_SHA256_*`, „signing subkey as the pin“ und „unknown fingerprint as the pin“ mit `PODMAN_SIGNING_KEY`), `docker build -q -t …` ohne `--force-rm`, Cleanup nur im Erfolgspfad. Derselbe Mechanismus ist anzunehmen, aber nicht gemessen; die Pins unterscheiden sich je Teilfall, und die Teilfälle scheitern an verschiedenen Schritten.

Bewertung gegen die DoD-Bedingung (§1, vierter NICHT-Punkt: „erst, wenn die Prüfung in der DoD ein zweites Vorkommen zeigt“): Die Prüfung der DoD ist auf den Ollama-Test zugeschnitten und zeigt kein zweites Vorkommen; es ergibt sich aus Lesen. Die Empfehlung des Reviews (Folge-Vorgang statt jetzt) ist damit **konform mit der DoD**: Die allgemeine Aufräum-Hilfe steht ausdrücklich außerhalb dieses Slice, §4 benennt die Rückführung für „mehrere Docker-Tests“, und eine Ausweitung hätte je Teilfall Messung und Pin-Wahl verlangt. Ich schließe mich der Empfehlung an. Bedingung: Der Folge-Vorgang bekommt einen auflösbaren Anker (Auflage 2); heute steht er nur im Fließtext des Slice.

## 9. Kommentare, Suppression, Build-Tag

- Build-Tag `//go:build docker` steht in Zeile 1 der Datei; kein `//nolint` in der Datei (Suche leer).
- Die drei neuen Kommentare sind im Indikativ und tragen Kopplung (Stage `test-docker-tools` ohne buildx), Zusage (Zuordnung, Grenze gleichzeitiger Läufe) oder Abgrenzung; keine Chronik, kein „früher“ (`AGENTS.md` §3.7, §3.9 (a)). Der Verweis `--force-rm` und der Stage-Name stimmen mit dem Baum überein (§3.9 (b)). Die Helfer stehen **nach** der Testfunktion (§3.9 (c)).
- Fehlermeldungen des Tests sind Englisch.

## 10. Auflagen (Planner, Closure)

1. **DoD-Rest und Closure-Notiz:** §7 füllen (Steering-Loop-Lerneintrag), Register fortschreiben oder „keine Beobachtung“ notieren, Ausgang für die zwei Risiken in §6 setzen (Risiko 1: Zusage auf den Code zugeschnitten, gleichzeitiger Lauf als Grenze benannt; Risiko 2: Host belegt, CI-Runner nicht gemessen), die drei Paarungen tragen; den Punkt „`make test-docker` grün“ entweder mit dem Beleg des vollen Laufs oder mit einer benannten Begründung (bekannt roter Fremdtest, berührter Test gezielt grün) schließen.
2. **Folge-Slice für das zweite Vorkommen** in `open/` anlegen (`TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`, drei Teilfälle, je Pin und Mechanismus zu messen; Verallgemeinerung von `newContainersWithCommand` als Option) und in §7 verlinken.
3. **Rest-Risiko führen** (im Folge-Slice oder Register): Ein Wechsel des Test-Images auf buildx macht die Prüfung vakuum-grün; ein Test oder Hinweis, der den klassischen Builder voraussetzt, ist offen (kein Befund dieses Slice, Review INFO 2).
4. **Abgleich der Plan-Tabelle** nach diesem Report und der Closure-Bewegung wiederholen (§3.10); die Aussage „drei Dateien“ im Abschnitt Umsetzung auf den Endstand setzen.
5. Optional: das Wort für das Basis-Image in der Mechanismus-Aussage der Slice-Datei zu „ein Basis-Image“ kürzen.
6. Reihenfolge der Bewegung: erst Inhalt (DoD-Häkchen, Closure-Notiz), dann reiner `git mv` nach `done/` (`AGENTS.md` §3.3); der Done-Slice trägt den Commit-Hash.
