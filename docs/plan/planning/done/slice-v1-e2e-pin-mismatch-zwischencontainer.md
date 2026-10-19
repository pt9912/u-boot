# Slice v1: Docker-Test — die Negativprobe `PinMismatchBreaksBuild` hinterlässt den Zwischencontainer des fehlgeschlagenen Baus (`LH-FA-DEV-003`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nach dem Release-Cut `v0.12.0` (Empfehlung des Verifiers, Abschnitt 5 seines Berichts; die Entscheidung über die Reihenfolge liegt beim Projektinhaber). Der Slice hält den Cut nicht auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht aus dem Befund B-4 von [`slice-v1-devcontainer-ollama-amd64-messung`](../done/slice-v1-devcontainer-ollama-amd64-messung.md) hervor.

**Bezug:** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) (Katalog-Feature `ollama`; der Test pinnt die Prüfsumme der Lieferkette, [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)), [`ADR-0024`](../../adr/0024-ollama-im-sandbox-devcontainer.md) (Negativprobe zur Prüfsumme).
Anlass: Befund B-4 der `amd64`-Messung. Belegt (Messfrage A8, zwei Läufe): Nach dem Lauf von `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` steht ein Container im Status `exited`, Exit 1, mit von Docker vergebenem Namen und ohne Label; die Anzahl der Container nachher liegt um einen über der vorher ([`rohdaten-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-09.txt), [`rohdaten-lauf2-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-lauf2-09.txt), [`rohdaten-11-vorher.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-11-vorher.txt) und [`rohdaten-11-nachher.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-11-nachher.txt)). Der Befund betrifft jeden Lauf von `make test-docker` auf einem geteilten Host, nicht den Nutzer des Produkts.

**Berührte Spec-Stellen:** — (reine Test-Hygiene; kein Satz der Spezifikation oder der Sicht ist berührt).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Negativprobe `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` hinterlässt nach ihrem Lauf keinen Container des fehlgeschlagenen Baus auf dem Docker-Host; der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den Null-Pin im Env tragen (zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen).

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Test:** `internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` (ab Zeile 1006): `generateSandboxProjectWith`, dann `dockerOutput(ctx, t, "build", "-q", "-t", tag, "--build-arg", "OLLAMA_SHA256_AMD64="+zeros, "--build-arg", "OLLAMA_SHA256_ARM64="+zeros, …)`; der Test erwartet einen Exit-Code ungleich 0 und die Ausgabe `FAILED`. Eine `t.Cleanup`-Funktion steht nur im Pfad, in dem der Bau wider Erwarten gelingt (`rmi -f` des Tags); für den Container des fehlgeschlagenen Baus steht keine.
- **Das Rohdatum:** [`rohdaten-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-09.txt) weist einen neuen Container nach dem Lauf aus (`status=exited exit=1`), ohne weitere Angabe zu Path, Args oder Env (der Env-Befund ist eine eigene Messung dieses Slice, siehe „Umsetzung (Implementer)“); das Aufräum-Skript des Messprotokolls entfernt ihn nicht, weil nichts ihn eindeutig der Messung zuordnet (`tools/messung/ollama-amd64/09-tests.sh` weist ihn aus, `11-aufraeumen.sh` löscht nur nach Präfix).
- **Nicht untersucht:** welcher Builder den Container hinterlässt (klassischer Builder oder BuildKit im Image `test-docker-tools`) und ob der Docker-Daemon des CI-Runners dasselbe tut; die Messung zeigt das Ergebnis auf einem Host, nicht den Mechanismus. Ob andere Docker-Tests in `internal/e2e/` Container nach fehlgeschlagenen Bauten hinterlassen, ist nicht gelesen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Templates, Lastenheft, Spezifikation** — der Befund ist Test-Hygiene; die Prüfsummen-Prüfung im erzeugten Dockerfile bleibt wie sie ist.
- **Das Aufräum-Skript des Messprotokolls** — der Container hat dort keine eindeutige Zuordnung; die Lösung gehört in den Test, nicht in ein Fremdlösch-Skript (Aufräum-Reichweite nur nach Präfix).
- **Gate-Änderungen** (`-v` im CI-Aufruf, Required Checks) — beim Projektinhaber ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)).
- **Eine allgemeine Aufräum-Hilfe für alle Docker-Tests** — erst, wenn die Prüfung in der DoD ein zweites Vorkommen zeigt; sonst ein eigener Vorgang.

## 2. Definition of Done

- [x] (`1689c68`) **Mechanismus gemessen und Test-Eingriff:** Der Slice belegt, welcher Builder den Container hinterlässt, und räumt ihn im Test auf (Kandidaten, vom Architect zu messen und zu wählen: eine Build-Option, die Zwischencontainer auch bei einem Fehler entfernt, wo der Builder sie unterstützt; sonst die Differenz der Container-Listen vor und nach dem Bau, beschränkt auf den neu entstandenen Container mit dem Null-Pin im Env). Der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den Null-Pin im Env tragen; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen.
- [x] (`1689c68`, Mutation vom Verifier wiederholt, `1251e58`) **Pin des Eingriffs:** Ein Test oder eine Prüfung im selben Docker-Testlauf stellt sicher, dass nach `PinMismatchBreaksBuild` kein Container des fehlgeschlagenen Baus bleibt; der Fall wird einmal rot gesehen (Eingriff entfernt) und ist danach grün.
- [x] `make gates` grün (Verifier, `1251e58`; Closure-Lauf siehe §7); `make test-docker` **nicht als voller Lauf**, begründet: der volle Lauf ist bekannt rot durch den nicht berührten Fremdtest `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (Folge: [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md)), der berührte Test lief gezielt grün (`--- PASS`, Verifier und Implementer); Review durchgeführt, Report unter `docs/reviews/` liegt vor (`7422493`, Fix-Runde `18f4edc`, `b37c715`; `.harness/skills/reviewer.md`); Verification-Report (`1251e58`); Closure-Notiz mit Steering-Loop-Lerneintrag (§7); Beobachtungs-Register (`../observations/`) fortgeschrieben (§7); jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen (§7).

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` (`TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`) | update | Aufräumen des Zwischencontainers (`t.Cleanup`, Zuordnung nach Mechanismus); gegebenenfalls ein gemeinsamer Helfer neben `dockerOutput` |
| `docs/messungen/<slice>/` und `tools/messung/<name>/` | entfällt | Der Mechanismus ist im Wegwerf-Kontext gemessen; das Protokoll steht als Text in „Umsetzung (Implementer)“, kein Skript, keine Rohdaten-Datei |
| `docs/maintainer/quality.md` | unverändert | gelesen (Suche nach Aufräumen/Container/cleanup): die Datei führt keine Konvention zum Aufräumen von Docker-Objekten in Tests; die Treffer betreffen den Lint-Container (Z. 54), den Test-Container mit `--network=host` (Z. 184) und Trivy (Z. 244) |
| diese Slice-Datei, `open/README.md`, Register (`observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/`), Roadmap bei Start und Closure; Folge-Slice `open/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md`; `docs/reviews/` (Review-Report, Verification-Report) | update; neu | Lifecycle, Review und Verifikation; Folge-Slice und Register-Eintrag für das zweite Vorkommen |
| `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, `.github/workflows/`, `spec/` | unverändert | kein Produktcode, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Das Aufräumen verlangt eine Änderung an mehreren Docker-Tests; dann trennt der Architect den Helfer vom Einzelfall.
- `in-progress` → `open` (blockiert): Das Aufräumen ist ohne eine Gate-Änderung (Test-Image, Runner-Daemon) nicht zu erreichen; Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Aufräumen entfernt einen fremden Container auf einem geteilten Host.** — **Ausgang:** getragen, mit Grenze: Zuordnung nur über den nach der Vorher-Liste neu entstandenen Container mit dem Null-Pin im Env, nie über Status oder Zeitfenster allein (Code-Kommentar von `newContainersWithCommand`, Verification Abschnitt 4); zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen und stehen so in §1, DoD und Code. Der Verifier fand beim Lauf keinen fremden Container entfernt (Zahl vorher/nachher je 3).
- **Der Mechanismus unterscheidet sich zwischen lokalem Host und CI-Runner.** — **Ausgang:** weiter offen → [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Host belegt, klassischer Builder; Daemon des CI-Runners und ein Client mit buildx nicht gemessen; im Folge-Slice [`slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../open/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) als Risiko geführt, dort zusammen mit „Wechsel auf buildx macht die Prüfung vakuum-grün“).

## Umsetzung (Implementer)

**Änderung:** `internal/e2e/sandbox_devcontainer_docker_test.go`: `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `--force-rm`, hält vorher die Container-IDs fest (`containerIDs`) und prüft in einem `t.Cleanup`, ob ein neu entstandener Container den Null-Pin in Path, Args oder Env trägt (`newContainersWithCommand`); einen solchen entfernt es mit `docker rm -f` und meldet den Fehler. Der Pin der Prüfung ist die Zeichenkette `OLLAMA_SHA256_AMD64=` plus 64 Nullen; kein Status, kein Zeitfenster, kein Name zählt.

**Mechanismus (gelesen und gemessen auf dem Host, Docker 29.9.1, nativer `amd64`, 2026-10-19):**

- Das Image `test-docker-tools` (`Dockerfile`, Stage `test-docker-tools`) installiert nur `docker-ce-cli` und `docker-compose-plugin`; `docker buildx version` im Image antwortet `unknown command: docker buildx`, `DOCKER_BUILDKIT` ist leer. Der Bau des Tests läuft damit über den **klassischen Builder**, der den Container des fehlgeschlagenen Schritts behält.
- Wegwerf-Probe (ein Basis-Image, ein `RUN … && false`, im Image `test-docker-tools` gegen den Host-Daemon): ohne Option blieb ein Container zurück (Container-Zahl +1); mit `--force-rm` blieb keiner (Ausgabe „Removed intermediate container“, Zahl unverändert). Der Probe-Container wurde danach über seinen Befehl (`echo probe-`) identifiziert und entfernt.
- `--label` am `docker build` setzt das Label auf das Image, nicht auf den Zwischencontainer (`docker ps -a --filter label=…` fand ihn nicht). Eine Zuordnung über Label scheidet deshalb aus; die Zuordnung läuft über den Null-Pin im Env des Containers. Path/Args tragen nur den Namen der Variable im Skripttext, der Wert steht im Env (`Config.Env`).

**Wirkung vorher/nachher (Test gezielt, `go test -tags docker -timeout 30m -v -count=1 -run '^TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild$' ./internal/e2e/` im Image `u-boot:test-docker-tools`, Aufruf-Muster wie A8):**

| Lauf | Container vorher | Container nachher | Ergebnis des Tests |
|---|---|---|---|
| Gegenprobe: `--force-rm` entfernt, Prüfung aktiv | 3 | 3 | `FAIL`: „the failed build left 1 container(s) …“; der Cleanup hat den Container entfernt |
| Änderung wie geliefert | 3 | 4 | `PASS`; kein Container im Status `exited` auf dem Host |

Die Gesamtzahl taugt auf diesem geteilten Host nicht als Beleg: Läufe anderer Arbeiten auf dem Host starteten und beendeten Container zwischen den Messungen, das erklärt die 4 nachher; der Beleg ist die Prüfung nach Pin und neuer ID, nicht die Zahl. Die Gegenprobe findet genau einen Container mit dem Null-Pin, den Zwischencontainer, den [`rohdaten-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-09.txt) als `exited`, Exit 1 ausweist. Hinweis: Zwei verworfene Zwischenläufe prüften Path/Args; sie blieben wirkungslos grün, weil der Pin nur im Env steht (Path/Args tragen den Variablennamen im Skripttext, nicht den Wert). Die Prüfung liest das Env.

**Rote Mutation:** `--force-rm` aus dem `docker build`-Aufruf entfernen färbt den Test rot (Tabelle, erste Zeile).

**Zweites Vorkommen (nur gemeldet, nicht behoben):** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` (dieselbe Datei) baut in drei Teilfällen absichtlich fehlschlagende Images ohne `--force-rm` und ohne Cleanup; gleicher Mechanismus, nicht gemessen. **Vermerk für die Closure:** Folge-Vorgang für diesen Test (drei Teilfälle, gleiche Klasse, nicht gemessen); die Anlage des Folge-Slice macht der Planner bei der Closure (§1: allgemeine Aufräum-Hilfe erst bei zweitem Vorkommen, eigener Vorgang).

**Abgleich Plan-Tabelle gegen `git diff --stat -M 5eab741..HEAD` (Endstand nach Review-Fix-Runde, Verification-Report, Closure-Commits und Move):** Der Diff berührt den Test (Zeile 1), diese Slice-Datei (Sammelzeile; nach der Closure in `done/`), den Review-Report und den Verification-Report unter `docs/reviews/` (Sammelzeile), den Folge-Slice in `open/`, `open/README.md`, das Register (`BEO-E2E/docker-test-hinterlaesst-zwischencontainer/` mit drei Dateien; `BEO-DOK/plan-aussage-ueber-code/` mit einer neuen Evidence-Datei und der `state.md`) die Roadmap (Sammelzeile „Register, Roadmap“) sowie die Verweis-Nachzüge von `make slice-mv` in `done/slice-v1-devcontainer-ollama-amd64-messung.md` und `done/slice-v1-release-cut-v0.12.0.md` (Link auf diese Datei von `in-progress/` nach `done/`; keine inhaltliche Änderung, in der Plan-Tabelle durch die Sammelzeile „Lifecycle“ nicht gedeckt und hier benannt); jede Zeile hat einen Plan-Eintrag, jeder Eintrag mit Diff eine Zeile. Das Register-Verzeichnis `BEO-DOK/plan-aussage-ueber-code/` steht in der Plan-Tabelle unter „Register“ ohne Pfad und ist hier benannt.

**Nicht gemessen:** der Docker-Daemon des CI-Runners und ein Host mit BuildKit-Standard (§6, zweites Risiko).

## 7. Closure-Notiz

**Gegenstand und Commits.** `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `--force-rm`, hält vorher die Container-IDs fest und prüft in einem `t.Cleanup`, ob ein neu entstandener Container den Null-Pin im Env trägt; einen solchen entfernt es mit `docker rm -f` und meldet das Ergebnis (`internal/e2e/sandbox_devcontainer_docker_test.go`, Helfer `containerIDs` und `newContainersWithCommand`). Mechanismus: der klassische Builder des Images `test-docker-tools` (ohne buildx) behält den Container des fehlgeschlagenen Schritts, `--force-rm` entfernt ihn. Commits: `1689c68` (Test), `be66a0a` (Umsetzung und Plan-Tabelle), `7422493` (Review-Report), Fix-Runde `18f4edc` (Test: Cleanup unterscheidet Entfernen gelungen und fehlgeschlagen, Grenze im Kommentar) und `b37c715` (Slice: Env-Befund, begrenzte Zusage, Folge-Vermerk), `1251e58` (Verification-Report), `8c30937` (Folge-Slice und Register-Eintrag); Closure: `b24324a` (Inhalt), `80ab9fb` (reiner Move), `18173ff` (Verweis-Nachzug von `make slice-mv`), danach der Nachzug dieses Absatzes.

**Belege.** Review-Report [`2026-10-19-slice-v1-e2e-pin-mismatch-zwischencontainer`](../../../reviews/2026-10-19-slice-v1-e2e-pin-mismatch-zwischencontainer.md) (`7422493`: 0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO; beide LOW behoben durch `18f4edc` und `b37c715`) und Verification-Report [`2026-10-19-verification-slice-v1-e2e-pin-mismatch-zwischencontainer`](../../../reviews/2026-10-19-verification-slice-v1-e2e-pin-mismatch-zwischencontainer.md) (`1251e58`: GO mit Auflagen). Die Auflagen sind mit dieser Closure abgearbeitet: 1 (diese Notiz, Register, Ausgang der Risiken, Begründung zu `make test-docker`, Paarungen), 2 (Folge-Slice in `open/`), 3 (Rest-Risiko „Wechsel auf buildx“ im Folge-Slice §6 geführt), 4 (Plan-Tabelle und Abgleich), 5 (Wort für das Basis-Image gekürzt), 6 (erst Inhalt, dann reiner `git mv`).

**Läufe und Exit-Codes (Exit-Code des `make`-Prozesses selbst).** Vom Verifier gefahren (Stand `b37c715`): `make gates` Exit 0 (Coverage 91,40 % bei Schwelle 90 %), `make doc-trace` (83 Anforderungen, 0 Waisen), `make docs-check` (516 Dateien, 0 Befunde), `make doc-planning` (0 Befunde), `make ablage-pruefung` („26 Faelle, 0 Fehlschlaege“). Der Docker-Test (Build-Tag `docker`), vom Verifier im Image des Ziels `test-docker` gefahren: `--- PASS: TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild (32.46s)`; Container vorher 3, nachher 3, keiner mit Null-Pin im Env. Mutation (`--force-rm` entfernt, Wegwerf-Kopie): `--- FAIL` mit „the failed build left container … with the zero pin; cleanup removed it“. Mechanismus-Probe: `docker buildx version` im Image `test-docker-tools` antwortet „unknown command“. Die Läufe der Closure (nach Inhalt, Move und Nachzug) stehen im Bericht des Planners.

**Nicht gelaufen und Grenzen.** (1) Der volle `make test-docker` lief nicht: bekannt rot durch den nicht berührten Test `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` (Folge: [`slice-v1-egress-dns-upstream-schliessen`](../done/slice-v1-egress-dns-upstream-schliessen.md)); der berührte Test lief gezielt grün. Der DoD-Punkt ist mit dieser Begründung geschlossen, nicht mit einem vollen grünen Lauf. (2) `make verify-depguard` lief nicht (keine Importregel berührt; `make lint` in `make gates` grün). (3) Der Zweig „`docker rm -f` scheitert“ im Cleanup ist gelesen, nicht ausgelöst; der Zweig „Entfernen gelungen“ ist durch die Mutation gesehen. (4) Der Docker-Daemon des CI-Runners und ein Client mit buildx sind nicht gemessen. (5) Zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen. (6) Die Gesamtzahl der Container taugt auf einem geteilten Host nicht als Beleg; der Beleg ist die Prüfung nach Pin und neuer ID.

**Was hat funktioniert:** Die Messung im Wegwerf-Kontext trennte den Mechanismus (klassischer Builder, `--force-rm`) von der Zuordnung (Label am Image, nicht am Container; Pin im Env). Der Verifier wiederholte die Mutation unabhängig und fand die Prüfung nicht vakuum-grün.

**Was ging anders als geplant:** Zwei verworfene Zwischenläufe prüften Path und Args auf den Null-Pin und blieben wirkungslos grün, weil der Pin nur im Env steht; der Plan hatte „Pin in Path/Args“ als Aussage geführt, ohne sie am Container gelesen zu haben (Review LOW 2). Die erste Zusage („entfernt nur den Container dieses Laufs“) war weiter als der Code (Review LOW 1) und ist auf „neu seit Vorher-Liste und Pin im Env“ mit der Grenze gleichzeitiger Läufe zugeschnitten. Das zweite Vorkommen (`NestedPodman_PinMismatchBreaksBuild`) zeigte sich erst beim Lesen und nicht in der DoD-Prüfung; es ist ein Folge-Slice.

**Steering-Loop-Eintrag.** Lerneintrag als **gezählte Beobachtung, keine geschärfte Regel verkörpert, kein neuer Sensor:** (1) Neue Beobachtung [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md): ein Docker-Test mit absichtlich fehlschlagendem Bau hinterlässt Objekte auf dem geteilten Host, ohne dass der Test rot wird; **1×**, Ausgang offen (nicht zugewiesen); `verkörpert` oder `geplant` wären unter 3× nicht zulässig. Das zweite Vorkommen ist gelesen, nicht gemessen (Folge-Slice). (2) Die Aussage „Pin in Path/Args“ im Plan und die zu weite Zusage zur Zuordnung sind Rückfälle der Klasse [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in `AGENTS.md` §3.10; Fund im Review, nicht im Plan; die Regel bleibt unverändert, hier ein gezählter Beleg: neue Evidence-Datei `slice-v1-e2e-pin-mismatch-zwischencontainer.md`. (3) Der Kandidat „ein Test, der eine Prüfung gegen den Mechanismus pinnt, sagt, welchen Mechanismus er voraussetzt“ (Wechsel auf buildx macht die Prüfung vakuum-grün) wird im Folge-Slice als Risiko mit Ausgang geführt, **nicht** eigenmächtig in `AGENTS.md` verkörpert.

**Beobachtungs-Register (`../observations/`).** Neu angelegt: [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) mit `evidence/slice-v1-e2e-pin-mismatch-zwischencontainer.md` (1×). Fortgeschrieben: [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (eine weitere Evidence-Datei, verkörpert in `AGENTS.md` §3.10, keine weitere Schärfung). Gelesen, ohne neuen Beleg: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (anderer Mechanismus, kein Auftreten) und [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; der berührte Test ist ein anderer, kein neuer Beleg). Die Evidence-Dateien tragen nur die Beobachtungsform. **Kein Eintrag über der Schwelle von 3× außer den schon verkörperten.**

**Folge-Slices:** [`slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../open/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (zweites Vorkommen, drei Teilfälle, Messung zuerst; führt das Risiko „Wechsel des Test-Images auf buildx macht die Prüfung vakuum-grün“ und den Daemon des CI-Runners).

**Zur Entscheidung vorgelegt (Projektinhaber), nicht entschieden:** (1) **Gate-Lücke `doc-immutable`** bei Inline-Status-ADRs (Beobachtung `doc-immutable-greift-nicht-bei-inline-status`, bereits vorgelegt). (2) **Kopf von `.github/workflows/ci.yml`** (Gate-Kette; Workflow-Änderung, nicht dieser Slice). (3) **Archiv-Hinweis** in `done/`-Slices. Es wird nichts eigenmächtig in `AGENTS.md` verkörpert.

**Risiken aus §6:** jedes trägt einen Ausgang; getragen (fremder Container, mit Grenze gleichzeitiger Läufe), weiter offen → Register (Host belegt, CI-Runner und buildx nicht gemessen).

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: `slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer` existiert als Datei in `open/`. Register: `BEO-E2E/docker-test-hinterlaesst-zwischencontainer` und `BEO-DOK/plan-aussage-ueber-code` sind Verzeichnisse im Register, jedes mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/e2e` (E2E, Test-Infrastruktur) und gegebenenfalls `tools/`, `docs/` laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (Tests im selben Paket, die vom Docker-Zustand des Hosts abhängen; der Fehler liegt im Test, nicht im Produkt), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; CI-Aufruf ohne `-v`). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
