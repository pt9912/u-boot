# Verifier-Report: slice-cr-lastenheft-messstaende — 2026-10-25

**Rolle:** Verifier (Modul 8), frischer Kontext. Frage: Bauen wir es richtig, gemessen an Plan, DoD und Entscheidungen? Kein Produktcode; der Slice ist ein Lastenheft-CR (`0.23.2` auf `0.24.0`) mit Abgleich nach unten.

**Gegenstand:** `git diff 0d50bfd7..HEAD` (11 Dateien), Commits `9ffbfda4` (Freigabe-Eintrag), `e2ea8398` (nur `spec/lastenheft.md`), `939c3739` (Abgleich nach unten), `9f04393d`/`b6b25c68` (Review-Report), `b17460e2` (Fix-Runde, ohne Lastenheft-Diff).

**Gelesen:** `harness/verification.md`; Slice (§1 Wortlaut Stellen 1 bis 5, Messlage M1 bis M5, Herkunftstabelle, Fragen; §2 DoD; §3 Plan-Tabelle; §6; §7); Review-Report `docs/reviews/2026-10-25-slice-cr-lastenheft-messstaende.md`; Messprotokoll `docs/messungen/slice-cr-lastenheft-messstaende/README.md`; `spec/lastenheft.md` (Diff und Resttreffer); [`ADR-0024`](../plan/adr/0024-ollama-im-sandbox-devcontainer.md) Entscheidung 10; [`ADR-0036`](../plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md) (Messung, Bedingung jeder Zeile: Linux `amd64`); Handbuch, beide READMEs, CHANGELOG; CI-Log selbst. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

## Ergebnis

Der Wortlaut im Lastenheft ist mit dem Entwurf in §1 und dem Freigabe-Eintrag in §7 Zeichen für Zeichen identisch, die Reihenfolge Freigabe vor Änderung belegt `git log`, und die Zusagen gehen nicht über das Gemessene hinaus. Kein HIGH, kein MEDIUM. Zwei LOW (eine Aussage im Handbuch ist breiter als das CI-Log; ein Zustandsmarker und eine Plan-Zeile sind nach Review und Fix-Runde nicht nachgezogen), vier INFO. Die DoD-Punkte 1 und 2 sind erfüllt; Punkt 3 ist bis auf die Carveout-Bindung erfüllt (offen für die Closure); Punkt 4 ist nicht abgeschlossen (Verifier-Report ist dieser Report, Closure-Notiz, Lerneintrag und Register fehlen noch, wie vorgesehen).

## Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | grün, Exit 0 | `lint`, `test`, `coverage-gate`, `docs-check`, `baseline-verify`, `ablage-pruefung` (26 Fälle, 0 Fehlschläge), `record-gates` |
| `make doc-trace` | grün | 84 Anforderungen, 0 Waisen |
| `make doc-commits RANGE="origin/main..HEAD"` | grün | 10 Commits aufgelöst, 0 Befunde |
| `make doc-immutable RANGE="d9dd171e..HEAD"` | grün | 82 Commits, 0 Befunde |
| ADR-Diff leer | ja | `git diff --stat d9dd171e..HEAD -- docs/plan/adr` ohne Ausgabe; ebenso `0d50bfd7..HEAD` |
| CI-Log nachgelesen | ja | `gh api repos/pt9912/u-boot/actions/jobs/121082471598/logs`, 11608 Zeilen |
| `make test-docker`, `make verify-depguard` | nicht gelaufen | Diff ohne Code, kein Docker-Pfad, keine Importregel berührt (Docs-only, wie beauftragt) |

Ausweisung: Ich habe an drei Stellen `sed -n` als reinen Lesebefehl auf Repo-Dateien benutzt (Verstoß gegen die Vorgabe „kein sed“); es wurde nichts geschrieben, kein Befund hängt davon ab, und die Zeilen sind mit dem Read-Werkzeug nicht nachgeprüft. Der Wortvergleich lief mit `grep -F` und `grep -F -x` auf Host-Shell mit bash-Ersetzung (reiner Text, kein Gate-Nachweis).

## Prüfung 1: Wortlaut gegen Freigabe

- Für Stellen 1, 2, 3, 4 (Zeilen 414 und 418) und die Historie-Zeile `0.24.0`: der Text aus §1 *Neu* (Datei-Links durch Anker, `JJJJ-MM-TT` durch `2026-10-25`) steht im Lastenheft: Stelle 2, 3 und Historie-Zeile mit `grep -F -x` je 1 Treffer, Stellen 1 und 4 als Teilsatz je 1 Treffer. Der *Alt*-Text der Stellen 1 bis 4 hat im Lastenheft 0 Treffer.
- Diff: nur Zeilen 9, 412, 414, 418, 472, 474 geändert und 1607 neu; Klasse B, C und Historie `0.7.0` bis `0.23.2` unberührt. Resttreffer von „gemessen“, „Stand der Messung“, „gemessenen Hosts“ nur in Historie-Zeilen (`0.7.0`, `0.12.0`, `0.13.0`, `0.15.0`, `0.23.1`, `0.24.0`) und in Zeile 416 (`gh`/`glab`, siehe V-4); „Gemessen ist das auf“, „Docker 29“, „Colima 0.10“ im Anforderungstext: 0 Treffer.
- Reihenfolge: `9ffbfda4` (2026-10-25 09:06:52, nur Roadmap und Slice) liegt vor `e2ea8398` (09:07:33, nur `spec/lastenheft.md`). In `9ffbfda4` steht kein Lastenheft-Diff.
- **Nicht prüfbar aus dem Repo:** dass die Vorlage an den Projektinhaber den Entwurf ungekürzt zeigte. §7 behauptet es („Vorlage ungekürzt mit Alt und Neu“, „Freigeben wie vorgelegt“); der Beleg ist das Gespräch. Der Abstand der beiden Commits (41 Sekunden) beweist weder das eine noch das andere. Ausgewiesen, nicht geraten.

## Prüfung 2: keine Zusage weiter als Messung

| Zusage im Lastenheft `0.24.0` | Messung | Befund |
|---|---|---|
| `ollama`-Installation, Linux `amd64`, lokale Docker-Engine | zwei Läufe (M1) und CI-Pass; im Log gelesen: `PASS TestE2E_SandboxDevcontainer_Ollama` und `…_Ollama_PinMismatchBreaksBuild`, Run 40164437737, Commit `f0d85e5b`, Runner-Image `ubuntu-24.04` 20261012.334.1 | gedeckt; ADR-0024 Entscheidung 10 nennt als Beleg „ein CI-Lauf mit sichtbarem Pass“ und `-v`: beides erfüllt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| `ollama`-Installation, Colima `arm64` | M2, ein von Hand nachgeschriebenes Probe-Image | bleibt (bestehende Zusage seit `0.12.0`); Grenze ausgewiesen (Slice §6, Satz 1.6, Handbuch Zeile 2174 bis 2175: Probe-Image, nicht das Erzeugnis). Kein neues Versprechen, keine Rücknahme |
| `podman build`, nur Images von Docker Hub (`docker.io`) | M3 und M4, beide `docker.io/library/…` | gedeckt; die Einschränkung steht im Wortlaut |
| Start von `podman`, `amd64` und Colima, jeweils mit `…=0` | M1 A6 und M4 (`amd64`, `…=0`), M3 (Colima) | gedeckt; **nicht** durch die CI: Log hat vier `--- SKIP` (`NestedPodman_Build`, `…_BuildEgress`, `NestedPodman`, `StatusFiles_VFS`), Meldung nennt `kernel.apparmor_restrict_unprivileged_userns=1`. Die Zusage hat für `amd64` kein CI-Gate; Handbuch, READMEs und CHANGELOG weisen das aus |
| Cloud-Modelle und Host-Anschluss | unverändert, nur `amd64` | unverändert; ADR-0036 *Bedingung jeder Zeile*: Linux `amd64`, lokale Engine | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

Die Bedingung `…=0` mit Exit-Code `11` bei `1` stützt sich laut Herkunftstabelle auf M3 und auf `internal/hexagon/application/doctor_sandbox.go` sowie `sandbox-init.sh.tmpl`; sie war schon vor diesem Slice Vertragstext und ist unverändert.

## Prüfung 3: Abgleich nach unten

- Das Messprotokoll stimmt mit dem Log überein: Run, Job, Commit, Start/Ende, Runner, 856,974 s, die elf `PASS`/`SKIP`-Zeilen des Auszugs wörtlich, vier `SKIP` mit derselben Meldung. F-2 und F-3 des Reviews sind behoben (Pass nur als Zeile `--- PASS`, Ableitung aus dem Testcode ausgewiesen; vier statt drei Skips). Die Marker `LOOPBACK_ANSWERS`, `LIBRARY_CPU`, `PULL_OK`, `MODEL_SURVIVES_CONTAINER_REMOVAL` und die Prüfliste stehen im Testcode (Zeilen 1229 bis 1244, 1307, 1324).
- CHANGELOG `[Unreleased]`: sagt nicht mehr „resolves“; er sagt, die zwei Mengen bestehen, stehen im Wortlaut und sind nicht vereinheitlicht. F-1 ist behoben. Der Satz stimmt mit dem Lastenheft überein (Installation `amd64` und Colima; Cloud und Anschluss nur `amd64`).
- README.md und README.de.md: Aussagen („promised … measured … in one CI run on `amd64`“, „the CI skips the Podman tests because its runner restricts nested user namespaces“) stimmen mit Lastenheft und Log. Keine Kennung im Diff von README und `docs/user/` (Suche über die hinzugefügten Zeilen nach `LH-`, `ADR-`, `SPEC-`, `slice-`, `welle-`: 0 Treffer; AGENTS.md §5 Regel 6 gehalten).
- Handbuch: die Zusage-Sätze (Zeilen 1383, 2159 bis 2161, 2464) stimmen mit dem Lastenheft. Die Bestände „nicht gemessen … `amd64`“ an Zeile 1723 und 1764 betreffen andere Messungen (Neuaufbau nach Egress-Änderung; Zeiten), keine Zusage-Menge, und kollidieren nicht. Eine Aussage ist breiter als das Log: V-1.

## Prüfung 4: Widersprüche in Lastenheft, Spec, Handbuch

- `spec/lastenheft.md`: kein Messstand mehr im Anforderungstext außer Zeile 416 (V-4). Zeile 414 trennt die Zusage der Cloud-Modelle ausdrücklich von der Installation; der Absatz `ollama` (Colima zugesagt) und Zeile 414 (Colima nicht zugesagt) widersprechen sich dadurch nicht.
- `spec/spezifikation.md`: Zeile 742 (Nicht zugesagt) verweist für „andere Architekturen und Container-Engines“ auf [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features); der Verweis trägt weiter (der Lastenheft-Absatz nennt die Menge jetzt im Wortlaut). Zeile 331: kein Messstand der Menge. Keine Treffer für „nicht verifiziert“ in `spec/`, Handbuch und READMEs; die in ADR-0024 Entscheidung 10 für `amd64` genannte Kennzeichnung steht dort nirgends mehr, die ADR bleibt unverändert (Accepted, Diff leer). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Handbuch: Zeilen 1545 und 1573 („die CI“ nicht gemessen) betreffen die DNS-Messung und stimmen mit den übersprungenen Podman-Tests.

## Prüfung 5: Plan-Tabelle, Zustandsmarker, Carveout

`git diff --stat 0d50bfd7..HEAD` hat 11 Dateien; alle stehen in der Plan-Tabelle oder in §7 (`roadmap.md` über den Startvermerk, `docs/messungen/README.md` als nachträglich genannt). Abweichungen siehe V-2. `spec/spezifikation.md`, `spec/architecture.md`, `docs/plan/adr/`, `harness/replay.md`, `internal/**`, `Makefile`, `repo.mk`, `.github/workflows/` unverändert (Diff leer).

## Findings

| ID | Kategorie | Befund | Pfad | Klasse |
|---|---|---|---|---|
| V-1 | LOW | Das Handbuch sagt, die CI habe „die Docker-Tests mit `nestedRuntime: podman`“ übersprungen (Zeilen 1384, 2173, 2474). Das Log zeigt: übersprungen sind `NestedPodman_Build`, `…_BuildEgress`, `NestedPodman` und `StatusFiles_VFS`; **bestanden** hat `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` mit drei Unterfällen (Pin-Prüfung im Podman-Pfad). Das Messprotokoll ist genau (drei Podman-Tests plus ein vierter), das Handbuch nennt die Menge pauschal. Eine Aussage über Verhalten, breiter als der gelesene Beleg (AGENTS.md §3.10). | `docs/user/benutzerhandbuch.md` Zeilen 1384, 2173, 2474: „die CI übersprang am 2026-10-25 die Docker-Tests mit `nestedRuntime: podman`“ | aussage-breiter-als-beleg |
| V-2 | LOW | Plan-Abgleich (§3.10) nach der letzten Fix-Runde nicht vollständig: (a) `docs/plan/planning/open/README.md` Zeile 33 trägt noch „Review, Verifier und Closure stehen aus“; der Review und die Fix-Runde liegen vor (Zustandsmarker mit Ablaufpunkt, nicht auf den Ist-Stand gesetzt). (b) Die Plan-Tabelle führt `open/README.md` mit „Zeile 33 und Zeile 73“; der Diff berührt nur Zeile 33. (c) `carveouts.md` steht als „update“ in der Tabelle, hat keinen Diff in `0d50bfd7..HEAD` (die Auflösung folgt bei der Closure; die Tabelle sagt es für den Planner-Teil, aber nicht, dass der Diff bis zur Closure leer ist). | `docs/plan/planning/open/README.md` Zeile 33; Slice §3, Zeilen „diese Slice-Datei, `open/README.md` (Zeile 33 und Zeile 73)“ und „`carveouts.md`“ | plan-tabelle-gegen-diff |
| V-3 | INFO | **Was die Closure entscheiden muss: Carveout „Zwei Zusage-Mengen“.** Lage: Das Lastenheft `0.24.0` führt für `ollama` weiter zwei Mengen (Installation: `amd64` und Colima `arm64`; Cloud-Weg und Host-Anschluss: nur `amd64`), jetzt im Wortlaut statt als Messstand, und Zeile 414 weist die Trennung aus; der CHANGELOG sagt ausdrücklich „keeps them apart, it does not unify them“. Die Carveout-Zeile beschreibt die Spannung als „Absatz zur Installation führt die Menge des Messstands (Colima, `arm64`)“; diese Beschreibung ist nach `0.24.0` nicht mehr wahr (kein Messstand mehr), die Mengendifferenz besteht fort. Frage 1d (Projektinhaber, 2026-10-25) entschied „bei der Closure aufgelöst“; §7 *Fix-Runde* sagt dagegen „entscheidet der Planner dort“ (Auflösung oder Weiterführung). Die Quelle widerspricht sich: Entscheidung 1d ist die höhere Bindung, §7 bezeichnet sie als offen. Die Closure muss eines von zweien tun und begründen: (1) auflösen nach Frage 1d (Grund: die Spannung bestand durch den Messstand als Mengenbestimmung; die zwei Mengen sind jetzt Vertrag und ausgewiesen, kein temporärer Zustand mehr; Zeile aus der Tabelle, Eintrag im Audit-Trail), oder (2) die Carveout-Zeile verlängern mit neuer Beschreibung und neuem Trigger, falls der Projektinhaber die Mengen doch vereinheitlichen will. Der Verifier entscheidet das nicht. Rückfrage an den Projektinhaber nötig, wenn Planner und 1d auseinandergehen. | Slice §2 Punkt 3, §7 *Fix-Runde*; `docs/plan/planning/in-progress/carveouts.md` Zeile 27 | closure-entscheidung |
| V-4 | INFO | Offener Messbefund im Anforderungstext, nicht Teil dieses Slice: `LH-FA-DEV-003` Absatz `gh`/`glab` Zeile 416 „hat die Messung keine Signatur gefunden“ (Frage 8); der Kommentar vor `TestE2E_SandboxDevcontainer_Ollama` (`internal/e2e/sandbox_devcontainer_docker_test.go` Zeile 1252, „the amd64 run is not a verification … until a CI run with a visible pass“) nennt einen erledigten Zustand. Beides ist im Slice benannt und für Folgeschritte gemeldet; die Closure trägt es als Ausgang. | `spec/lastenheft.md` Zeile 416; Testkommentar | folge-schritt | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| V-5 | INFO | Die Zusagen für `podman` auf `amd64` haben kein CI-Gate: der Runner hat `…=1`, alle Podman-Tests überspringen. Wer die Zusage prüfen will, braucht einen Runner mit `…=0` oder eine Nachmessung. Das ist kein Mangel dieses Slice (Messung gelesen, Zusage nach Entscheidung des Projektinhabers), aber die Closure sollte es als Risiko mit Ausgang tragen (Slice §6 führt es nicht als eigenes Risiko). | Messprotokoll *Lesehilfen*; Log `--- SKIP` | zusage-ohne-ci-beleg |
| V-6 | INFO | Die Freigabe-Vorlage ist aus dem Repo nicht prüfbar (siehe Prüfung 1). Die Beobachtung `freigabe-vorlage-verdichtet-den-entwurf` bleibt offen; ein Beleg „keine Verdichtung“ ist aus dem Repo nur als Aussage des Implementers lesbar. | Slice §7 *Freigabe je Stelle* | nicht-pruefbar |

## DoD-Abgleich (Slice §2)

- [x] **Freigabe vor der Änderung:** §7 *Freigabe* je Stelle mit Datum, eigener Commit `9ffbfda4` vor `e2ea8398`, Abweichung gegen den Entwurf „keine“. Vorbehalt: Inhalt der Vorlage nicht prüfbar (V-6).
- [x] **Wortlaut im Lastenheft:** wortgleich (Prüfung 1); Klasse A ohne Messumgebung, Prüfzustand, Messbilder; B und C unverändert; Version `0.24.0` und Historie-Zeile ohne Messangabe; kein Verweis auf ADR, Slice oder Roadmap im geänderten Text (Anker nur auf Lastenheft-Abschnitte); jede Bedingung steht in der Herkunftstabelle mit Messung und Pfad (gegen Log, Messprotokolle und Testcode gelesen).
- [~] **Abgleich nach unten und Carveout:** Handbuch, READMEs, CHANGELOG nachgezogen und gegen das Log gelesen (V-1 offen); Messprotokoll liegt im Repo; `make doc-trace` und `make docs-check` ohne Befund. Der Carveout ist **noch nicht** aufgelöst oder weitergeführt: Closure-Aufgabe (V-3).
- [~] **Gates und Review:** `make gates` grün (selbst gefahren); Review-Report liegt vor, Fix-Runde `b17460e2` ohne Lastenheft-Diff (`git diff --stat 939c3739 HEAD -- spec` leer), ein Nachlauf-Review ist deshalb nicht gefordert. Offen für die Closure: Closure-Notiz mit Steering-Loop-Lerneintrag, Beobachtungs-Register, Ausgang je Risiko aus §6 (mindestens V-5 ergänzen), die drei Paarungen.

## Sensoren nach `harness/verification.md`

Mindest-Sensor „Nur Markdown/Doku“ ist `make docs-check` (Teil von `make gates`, grün) plus gezielter Linkcheck (`make doc-trace`, `make doc-commits`, `make doc-immutable`, grün). `make test-docker` und `make verify-depguard`: nicht gelaufen, Grund: Diff ohne Code und ohne Importregel.

## Nicht geprüft

- Die ungekürzte Vorlage an den Projektinhaber und seine Antwort (Gespräch).
- Die Engine-Version des CI-Runners (steht nicht im Log).
- Eine Messung des Erzeugnisses unter Colima (`arm64`) und ein Lauf von Podman auf `amd64` in der CI (nicht Gegenstand).

## Empfehlung an den Planner

Freigabe zur Closure nach: V-1 (Handbuch-Formulierung auf die vier übersprungenen Tests eingrenzen oder den bestandenen Pin-Test erwähnen) und V-2 (Zeile 33 in `open/README.md` auf den Ist-Stand, Plan-Tabelle berichtigen) in einer Fix-Runde ohne Lastenheft-Diff; V-3 als Entscheidung der Closure mit dem Projektinhaber klären.
