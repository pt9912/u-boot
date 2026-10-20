# Slice v1: Docker-Test — die Negativprobe `NestedPodman_PinMismatchBreaksBuild` hinterlässt die Zwischencontainer der fehlgeschlagenen Baue (`LH-FA-DEV-007`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung: nach dem Release-Cut `v0.12.0`. Der Slice hält keinen Cut auf und hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht aus dem zweiten Vorkommen hervor, das [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../done/slice-v1-e2e-pin-mismatch-zwischencontainer.md) (Umsetzung, Review INFO 1, Verification Abschnitt 8) meldete.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Option `nestedRuntime: podman`; der Test pinnt Prüfsumme und Signaturschlüssel der Lieferkette, [`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)), [`ADR-0023`](../../adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) (Negativprobe zu Prüfsumme und Signaturschlüssel).
Anlass: Der Vorgänger räumte `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` auf und fand beim Lesen denselben Bau fehlschlagender Images in `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`. Gelesen, nicht gemessen: Der Test baut in drei Teilfällen ohne `--force-rm` und ohne Cleanup für den Zwischencontainer.

**Berührte Spec-Stellen:** — (reine Test-Hygiene; kein Satz der Spezifikation oder der Sicht ist berührt).

**Verantwortlich:** pt9912

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` hinterlässt nach seinem Lauf keinen Container eines fehlgeschlagenen Baus auf dem Docker-Host, solange der Bau mit dem klassischen Builder läuft (Image `test-docker-tools`, ohne buildx; unter BuildKit wird der Test rot statt grün, siehe §6 Risiko 2); der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den jeweiligen Pin im Env tragen (zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen).

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Test:** `internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` (Zeilenangabe entfällt, sie wandert mit dem Code): drei Teilfälle in einer Map: „wrong checksum“ (`PODMAN_SHA256_AMD64` und `PODMAN_SHA256_ARM64` je 64 Nullen, erwartet `FAILED`), „signing subkey as the pin“ (`PODMAN_SIGNING_KEY` mit einem Fingerabdruck) und „unknown fingerprint as the pin“ (`PODMAN_SIGNING_KEY` mit 40 Nullen), beide erwarten `.tar.gz: OK`. Der Aufruf ist `docker build -q -t <tag> --build-arg … <Kontext>` ohne `--force-rm`; eine `t.Cleanup`-Funktion steht nur im Pfad, in dem der Bau wider Erwarten gelingt (`rmi -f` des Tags). Die Teilfälle scheitern an verschiedenen Schritten des Baus (Prüfsumme gegen Signaturprüfung).
- **Der Helfer des Vorgängers:** `containerIDs` und `newContainersWithCommand` (dieselbe Datei, nach `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`); `newContainersWithCommand` nimmt das Suchwort als Parameter und liest Path, Args und Env per `docker inspect`.
- **Bei der Anlage nicht gemessen, mit der Umsetzung gemessen (Host, Docker 29.9.1, `amd64`; Tabelle im Abschnitt Umsetzung):** die drei Teilfälle hinterlassen je einen Container des fehlgeschlagenen Schritts, der Wert des Pins steht in jedem Fall im Env. Weiter nicht gemessen: der Docker-Daemon des CI-Runners.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Templates, Lastenheft, Spezifikation** — Test-Hygiene; die Prüfung im erzeugten Dockerfile bleibt wie sie ist.
- **Das Ollama-Gegenstück** — der Vorgänger hat es geliefert und geschlossen; hier wird es nicht angefasst, es sei denn, ein gemeinsamer Helfer ändert seine Signatur (dann nur mechanisch, mit grünem Lauf des Ollama-Tests).
- **Eine allgemeine Aufräum-Hilfe für alle Docker-Tests** — erst, wenn die Messung hier und im Vorgänger denselben Weg trägt und ein drittes Vorkommen bekannt ist; sonst ein eigener Vorgang. Ob andere Tests in `internal/e2e/` Container nach fehlgeschlagenen Bauten hinterlassen, ist nicht gelesen.
- **Gate-Änderungen** (Test-Image mit buildx, `-v` im CI-Aufruf, Required Checks) — beim Projektinhaber ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)).

## 2. Definition of Done

- [x] **Mechanismus gemessen und Test-Eingriff:** Der Slice misst je Teilfall, ob und welcher Container zurückbleibt und welchen Pin er im Env (oder anderswo) trägt; dann gilt dieselbe Lösung wie im Vorgänger (`--force-rm` im Bau und Prüfung nach neuer ID und Pin, `newContainersWithCommand` mit dem Pin des Teilfalls als Suchwort) oder eine begründete andere. Der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den Pin des Teilfalls tragen; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen.
- [x] **Pin des Eingriffs:** Die Prüfung steht im selben Docker-Testlauf, deckt alle drei Teilfälle und wird je Teilfall (oder für den Test als Ganzes, mit Begründung) einmal rot gesehen (Eingriff entfernt) und ist danach grün.
- [x] `make gates` grün; der berührte Docker-Test gezielt grün (der volle `make test-docker` nur, wenn der Fremdtest `NestedPodman_BuildEgress` nicht mehr rot ist, sonst mit benannter Begründung); Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen. — Evidence: `make gates` Exit 0 (Verifier auf `3e5c2c5`, Closure-Lauf siehe §7); Docker-Test gezielt grün auf `3e5c2c5`-Stand (`--- PASS`, drei Teilfälle, §7); der volle `make test-docker` vom Implementer auf `034e124` gemeldet (Exit 0, `NestedPodman_BuildEgress` grün), vom Verifier und bei der Closure nicht wiederholt; Review-Report `acc095a`, Verification-Report `7574cb8`; Inhalt der Closure in `4b26028`.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` (`TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`; Helfer `newContainersWithCommand`) | update | Aufräumen der Zwischencontainer je Teilfall (`--force-rm`, Vorher-Liste, `t.Cleanup`, Suchwort je Teilfall); der Helfer bleibt, wie er ist, oder wird begründet verallgemeinert |
| `docs/messungen/<slice>/` und `tools/messung/<name>/` | entfällt | Der Vorgänger maß den Mechanismus im Wegwerf-Kontext und führte das Protokoll als Text im Slice; hier gilt dasselbe, es sei denn, die Messung wird wiederholbar gebraucht |
| `docs/maintainer/quality.md` | unverändert | gelesen (Suche nach Aufräumen/Container/cleanup): keine Konvention zum Aufräumen von Docker-Objekten in Tests; die Treffer betreffen den Lint-Container, den Test-Container mit `--network=host` und Trivy |
| diese Slice-Datei (nach der Closure in `done/`), `open/README.md`, Roadmap (`in-progress/roadmap.md`) bei Start und Closure | update | Lifecycle |
| Register `observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/` (`observation.md`, `state.md`, neue `evidence/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md`) | update; neu | zweites Vorkommen, Ollama-Wächterlücke, CI-Runner |
| `docs/reviews/` (Review-Report, Verification-Report) | neu | Review und Verifikation |
| `done/slice-v1-e2e-pin-mismatch-zwischencontainer.md` | update (Link) | Verweis-Nachzug von `make slice-mv` auf diese Datei |
| `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, `.github/workflows/`, `spec/` | unverändert (erwartet) | kein Produktcode, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../done/slice-v1-e2e-pin-mismatch-zwischencontainer.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung zeigt je Teilfall einen anderen Mechanismus oder eine andere Zuordnung; dann trennt der Architect die Teilfälle.
- `in-progress` → `open` (blockiert): Das Aufräumen ist ohne eine Gate-Änderung (Test-Image, Runner-Daemon) nicht zu erreichen; Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Aufräumen entfernt einen fremden Container auf einem geteilten Host.** — **Ausgang:** getragen mit Grenze; Zuordnung nur über den nach der Vorher-Liste neu entstandenen Container mit dem Pin des Teilfalls, nie über Status oder Zeitfenster allein; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen (der Vorgänger trägt dieselbe Begrenzung). Beleg: Eingriff in `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` und Doc-Kommentar von `newContainersWithCommand` (Review F-3); die Grenze gleichzeitiger Läufe gegen einen Daemon bleibt und ist nicht erzeugt worden.
- **Ein Wechsel des Test-Images auf einen Builder mit buildx (BuildKit) macht die Prüfung vakuum-grün.** — **Ausgang:** geteilt. **Nested Podman: getragen** durch den Wächter `failedStepRemoved` in `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`: jeder Teilfall verlangt in der Bau-Ausgabe „Removed intermediate container …“ unmittelbar vor „The command '…“; die Mutation „Muster ohne Treffer“ färbt alle drei Teilfälle rot. Unter BuildKit wird der Test damit rot statt vakuum-grün (aus der Mutation abgeleitet, auf dem Host gegengeprüft: 0 Zeilen „Removed intermediate“ bei `--progress=plain`; mit BuildKit nicht gesehen). **Ollama-Hälfte: weiter offen** → Register [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md): `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `-q --force-rm` und trägt keinen Wächter; unter BuildKit bliebe er grün (Review F-2). Das Ollama-Gegenstück ist nach §1 nicht Gegenstand dieses Slice; ob ein Folge-Slice entsteht, entscheidet der Projektinhaber (§7). Ein Wechsel des Images ist eine Gate-Änderung und liegt beim Projektinhaber.
- **Der Docker-Daemon des CI-Runners verhält sich anders als der lokale Host.** — **Ausgang:** weiter offen → Register [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Abschnitt „Benannt, nicht gezählt“); nicht gemessen; der Slice belegt den Fall auf dem Host und nennt den CI-Lauf als nicht gemessen.
- **Die drei Teilfälle scheitern an verschiedenen Schritten, der Pin steht womöglich nicht in jedem Fall im Env des Zwischencontainers.** — **Ausgang:** entfallen; die Messung je Teilfall (Tabelle im Abschnitt Umsetzung) zeigt den Pin in jedem der drei Teilfälle im Env des zurückbleibenden Containers, das Suchwort `pin` trägt in allen.

## Umsetzung (Implementer)

**Änderung:** `internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`: jeder Teilfall hält vorher die Container-IDs fest (`containerIDs`), baut mit `--force-rm` und ohne `-q` und prüft in einem `t.Cleanup` mit `newContainersWithCommand`, ob ein neu entstandener Container den Pin des Teilfalls trägt (Feld `pin` der Tabelle); einen solchen entfernt er mit `docker rm -f` und meldet den Fehler. Der Helfer blieb unverändert. Zusätzlich verlangt jeder Teilfall in der Bau-Ausgabe die Zeile „Removed intermediate container …“ unmittelbar vor „The command '…“ (`failedStepRemoved`): sie belegt, dass der klassische Builder einen Container des fehlgeschlagenen Schritts angelegt und `--force-rm` ihn entfernt hat.

**Mechanismus je Teilfall (gemessen auf dem Host, Docker 29.9.1, nativer `amd64`, 2026-10-20; Lauf im Image `u-boot:test-docker-tools`, Test-Binary mit der Änderung, aber ohne `--force-rm` und mit einer vorübergehenden Ausgabe von `docker inspect`; die Ausgabe ist wieder entfernt):**

| Teilfall | Pin im Env des zurückbleibenden Containers | Zustand des Containers | Zuordnung |
|---|---|---|---|
| wrong checksum | `PODMAN_SHA256_AMD64=` plus 64 Nullen (und ARM64 ebenso) | `exited`, Exit 1 (`sha256sum -c`: `FAILED`) | genau ein neuer Container, Pin im Env |
| signing subkey as the pin | `PODMAN_SIGNING_KEY=6E8E…0BA4` | `exited`, Exit 1 (`VALIDSIG`-Vergleich) | genau ein neuer Container, Pin im Env |
| unknown fingerprint as the pin | `PODMAN_SIGNING_KEY=` plus 40 Nullen | `exited`, Exit 22 (`curl -f` auf den Schlüssel: 404) | genau ein neuer Container, Pin im Env |

Alle drei Teilfälle scheitern im selben `RUN`-Schritt (Step 7/12), an verschiedenen Stellen darin; der Wert steht in jedem Fall im Env des Containers. Mit `--force-rm` bleibt keiner zurück. Vorher/nachher der Container-Liste des Hosts, gleicher Lauf, geliefertes Testprogramm: dieselbe Menge IDs (3 vorher, 3 nachher, Prüfsumme der sortierten ID-Liste gleich).

**Warum das „Removed …“-Muster mit Nachbarzeile:** Ein erster Entwurf prüfte nur das Vorkommen von „Removed intermediate container“. Bei kaltem Schicht-Cache (der erste Teilfall des Mutationslaufs) druckt der klassische Builder die Zeile auch für erfolgreiche Zwischenschritte; die Prüfung blieb dort grün, obwohl der Container des fehlgeschlagenen Schritts zurückblieb. Das Muster verlangt deshalb die unmittelbar folgende Zeile „The command '“.

**Rote Mutationen (Gegenbeispiele, alle gesehen):**

| Mutation | Ergebnis |
|---|---|
| `--force-rm` aus dem `docker build` entfernt | `FAIL` in allen drei Teilfällen: „the failed build left container … with the pin …; cleanup removed it“ und „does not show the failed step's container removed“ |
| Muster `failedStepRemoved` verfälscht (kein Treffer möglich), `--force-rm` bleibt | `FAIL` in allen drei Teilfällen (der Wächter gegen den vakuum-grünen Lauf greift) |
| geliefert | `PASS`, alle drei Teilfälle |

**Wechsel auf BuildKit:** nicht ausgeführt (das Image hat kein buildx; ein Wechsel ist eine Gate-Änderung beim Projektinhaber). Gedeckt durch die Bedingung: BuildKit druckt „Removed intermediate container“ nicht, der Test würde rot (Wirkung aus der Mutation „Muster ohne Treffer“ abgeleitet, nicht mit BuildKit gesehen).

**Sensoren (Implementer, Stand `034e124`; der Kommentar-Fix `3e5c2c5` danach ändert nur einen Kommentar, sein Lauf steht in §7):** `make gates` Exit 0. Der berührte Test gezielt im Image des Ziels `test-docker` (`-tags docker -timeout 30m -v -count=1 -run '^TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild$' ./internal/e2e/`): `PASS`, drei Teilfälle. Voller `make test-docker`: Exit 0 (`internal/e2e` 655,7 s; `NestedPodman_BuildEgress` ist darin grün); die Liste der Container-IDs des Hosts war danach gleich der vorher. `make verify-depguard` lief nicht (keine Importregel berührt).

**Abgleich Plan-Tabelle gegen `git diff --stat -M 1b4352c..HEAD` (Endstand nach Review-Fix, Verification-Report, Closure-Commits und Move):** Der Diff berührt den Test (Zeile 1), diese Slice-Datei (nach der Closure in `done/`), `open/README.md` und `in-progress/roadmap.md` (Lifecycle-Zeile), das Register (`observation.md`, `state.md`, neue Evidence-Datei), Review- und Verification-Report unter `docs/reviews/` sowie den Verweis-Nachzug von `make slice-mv` im Vorgänger-Slice `done/slice-v1-e2e-pin-mismatch-zwischencontainer.md` (Link von `in-progress/` nach `done/`; Plan-Zeile „Link“). Der Verweis in `observation.md` wird im Register-Commit auf `done/` gesetzt. `quality.md`, `docs/messungen/`, `tools/messung/`, Produktcode, `Makefile`, `repo.mk`, Workflows und `spec/` sind nicht im Diff.

**Nicht gemessen:** der Docker-Daemon des CI-Runners; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon (sie sähen jeweils den Container des anderen als neu, wie im Vorgänger).

## 7. Closure-Notiz

**Gegenstand und Commits.** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` hält je Teilfall vorher die Container-IDs fest, baut mit `--force-rm` und ohne `-q` und prüft in einem `t.Cleanup` mit `newContainersWithCommand`, ob ein neu entstandener Container den Pin des Teilfalls trägt; einen solchen entfernt es mit `docker rm -f` und meldet ihn als Fehler. Der Wächter `failedStepRemoved` verlangt in der Bau-Ausgabe „Removed intermediate container …“ unmittelbar vor „The command '…“ (`internal/e2e/sandbox_devcontainer_docker_test.go`). Mechanismus: der klassische Builder des Images `test-docker-tools` (ohne buildx) behält den Container des fehlgeschlagenen Schritts, in allen drei Teilfällen mit dem Pin im Env; `--force-rm` entfernt ihn. Die Zusage gilt unter der Bedingung des klassischen Builders (BuildKit: Test rot, nicht grün) und eines Daemons ohne gleichzeitigen zweiten Lauf. Commits: `1b4352c` (Roadmap-Arbeitszeile), `034e124` (Test und Umsetzungsabschnitt), `acc095a` (Review-Report), `3e5c2c5` (Kommentar zum Wächter im Indikativ, Review F-1), `7574cb8` (Verification-Report); Closure: `4b26028` (Inhalt), `2210d26` (reiner Move), `c33e8f3` (Verweis-Nachzug von `make slice-mv`), danach der Nachzug dieses Absatzes (Roadmap „Nichts in Arbeit.“, Hashes).

**Belege.** Review-Report [`2026-10-20-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../../../reviews/2026-10-20-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (`acc095a`: 0 HIGH, 0 MEDIUM, 1 LOW behoben durch `3e5c2c5`, 5 INFO) und Verification-Report [`2026-10-20-verification-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer`](../../../reviews/2026-10-20-verification-slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md) (`7574cb8`: GO mit Auflagen V-1 bis V-8). Abgearbeitet: V-1 (§6, alle vier Risiken mit Ausgang), V-2 (Ollama-Hälfte im Register, Folge-Punkt unten), V-3 (Register: Evidence-Datei, `state.md` 2×, `observation.md`), V-4 (Plan-Tabelle und Abgleich im Endstand), V-5 (Roadmap, `open/README.md`, Verweise), V-6 (Zeilenangabe und „nicht gemessen“ auf Ist-Stand, voller Lauf als gemeldet ausgewiesen), V-7 (gezielter Lauf unten), V-8 (Grenzen unten).

**Läufe und Exit-Codes.** Vom Verifier gefahren (Stand `3e5c2c5`): `make gates` Exit 0 (Coverage 91,40 % bei Schwelle 90 %, `ablage-pruefung-test`: 26 Fälle, 0 Fehlschläge). Bei der Closure im Image des Ziels `test-docker` (Aufruf wie in `repo.mk`: `--network=host`, Quellbaum, Docker-Socket und `/tmp` gemountet; `go test -tags docker -timeout 30m -v -count=1 -run 'NestedPodman_PinMismatch' ./internal/e2e/`), Exit 0, Stand mit dem Kommentar-Fix `3e5c2c5`: `--- PASS: TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild (30.64s)` mit den Teilfällen `wrong_checksum` (10,57 s), `signing_subkey_as_the_pin` (10,41 s) und `unknown_fingerprint_as_the_pin` (9,65 s), alle `PASS`. `make gates` und `make docs-check` der Closure stehen im Bericht des Planners.

**Nicht gelaufen und Grenzen.** (1) Der volle `make test-docker` lief nur beim Implementer auf `034e124` (Exit 0, `internal/e2e` 655,7 s, `NestedPodman_BuildEgress` grün); der Verifier und die Closure haben ihn nicht wiederholt, die Angabe ist übernommen, nicht belegt. Der DoD-Punkt verlangt den berührten Test gezielt; der volle Lauf war nicht Bedingung, weil `NestedPodman_BuildEgress` im Vorgänger-Stand noch rot war. (2) Host-Docker 29.9.1, nativer `amd64`: der Mechanismus ist nur dort gemessen. (3) BuildKit: abgeleitet aus der Mutation „Muster ohne Treffer“ und auf dem Host am Wortlaut gegengeprüft (Review), nicht mit einem umgestellten Test-Image gesehen. (4) Der Docker-Daemon des CI-Runners ist nicht gemessen (§6 Risiko 3, Register). (5) Zwei gleichzeitige Läufe gegen einen Daemon sind nicht erzeugt; der Eingriff sähe den Container des anderen Laufs als neu (Review F-3). (6) Wortlaut-Kopplung (Review F-5): der Wächter hängt am Wortlaut des klassischen Builders; verschwindet der Builder oder ändert sich der Wortlaut, wird der Test rot mit einer Meldung, die Builder und `--force-rm` nennt, und ein Wortlaut-Wechsel ist von „Mechanismus nicht mehr belegt“ am Ausgang nicht zu unterscheiden. (7) Doppelung (Review F-6): der Aufräumblock steht in Ollama- und Nested-Test fast wortgleich; eine gemeinsame Hilfe ist nach §1 erst bei einem dritten Vorkommen vorgesehen. (8) `make verify-depguard` lief nicht (keine Importregel berührt). (9) Das Zählen der Container vorher/nachher taugt auf einem geteilten Host nicht als Beleg; der Beleg ist die Prüfung nach Pin und neuer ID.

**Was hat funktioniert:** Die Messung je Teilfall mit einer vorübergehenden Ausgabe von `docker inspect` trennte Mechanismus und Zuordnung und zeigte, dass der Pin in jedem Teilfall im Env steht. Der Review wiederholte die Mutationen unabhängig und fand die Prüfung nicht vakuum-grün.

**Was ging anders als geplant:** Der erste Entwurf des Wächters prüfte nur das Vorkommen von „Removed intermediate container“; bei kaltem Schicht-Cache druckt der Builder die Zeile auch für erfolgreiche Zwischenschritte, die Prüfung blieb grün, obwohl der Container des fehlgeschlagenen Schritts zurückblieb. Erst die Nachbarzeile „The command '“ bindet sie an den fehlgeschlagenen Schritt (Mutation gesehen). Ein Kommentar im Wächter hatte die Form „Without … “ (Review F-1, behoben). Die Closure-Voraussetzungen (Zustandsmarker in §6, Register, Plan-Tabelle) standen nach der Umsetzung noch auf dem Stand der Anlage (Verifier V-1 bis V-5).

**Steering-Loop-Eintrag.** Lerneintrag als **gezählte Beobachtung, keine geschärfte Regel verkörpert, kein neuer Sensor:** (1) [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md): zweites Vorkommen belegt (drei Teilfälle, gleicher Mechanismus, neue Evidence-Datei `slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md`); Stand **2×**, unter der Schwelle von 3×, Ausgang nicht zugewiesen. (2) Der Kandidat aus dem Vorgänger („ein Test, der eine Prüfung gegen einen Mechanismus pinnt, sagt, welchen Mechanismus er voraussetzt“) ist hier als Test-Voraussetzung ausgesprochen: der Wächter `failedStepRemoved` macht den vakuum-grünen Lauf unter BuildKit zum roten. Er bleibt ein Einzelfall im Nested-Test; die Ollama-Hälfte trägt ihn nicht (Register). Es wird nichts in `AGENTS.md` verkörpert.

**Beobachtungs-Register (`../observations/`).** Fortgeschrieben: `BEO-E2E/docker-test-hinterlaesst-zwischencontainer` (`evidence/slice-v1-e2e-nested-podman-pin-mismatch-zwischencontainer.md`, `state.md` auf 2×, `observation.md` auf Ist-Stand mit der Ollama-Wächterlücke und dem CI-Runner). Gelesen, ohne neuen Beleg: `BEO-E2E/ci-ueberspringt-nested-podman-build` (offen; der Slice ändert den CI-Aufruf nicht).

**Folge-Punkt, zur Entscheidung vorgelegt (Projektinhaber), nicht entschieden:** Soll `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` den Wächter gegen den vakuum-grünen Lauf übernehmen (Bau ohne `-q`, Muster „Removed intermediate container“ vor „The command '“)? Ein gemeinsamer Helfer ist damit nicht vorgegeben: er ist nach §1 erst bei einem dritten Vorkommen vorgesehen. Ein Wechsel des Test-Images auf BuildKit bleibt Gate-Änderung beim Projektinhaber.

**Risiken aus §6:** jedes trägt einen Ausgang; getragen mit Grenze (fremder Container), geteilt (buildx: Nested getragen, Ollama-Hälfte offen → Register), weiter offen → Register (CI-Runner), entfallen (Pin nicht in jedem Env).

**Drei Paarungen:** Anker: kein Feld `liegt in` in diesem Slice, nichts zu paaren. Folge-Slice: keiner angelegt; die Ollama-Entscheidung liegt beim Projektinhaber und ist oben vorgelegt. Register: `BEO-E2E/docker-test-hinterlaesst-zwischencontainer` ist ein Verzeichnis im Register mit nicht leerem `evidence/` (zwei Dateien).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/e2e` (E2E, Test-Infrastruktur) und gegebenenfalls `docs/` laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); beide Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (1×, mit dem Vorgänger angelegt; dieser Slice liefert den zweiten Beleg oder den Ausgang), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; CI-Aufruf ohne `-v`, berührt dieselbe Testgruppe). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
