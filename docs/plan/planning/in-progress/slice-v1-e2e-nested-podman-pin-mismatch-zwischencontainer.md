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

**Ziel:** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` hinterlässt nach seinem Lauf keinen Container eines fehlgeschlagenen Baus auf dem Docker-Host; der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den jeweiligen Pin im Env tragen (zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen).

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Test:** `internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` (ab Zeile 572): drei Teilfälle in einer Map: „wrong checksum“ (`PODMAN_SHA256_AMD64` und `PODMAN_SHA256_ARM64` je 64 Nullen, erwartet `FAILED`), „signing subkey as the pin“ (`PODMAN_SIGNING_KEY` mit einem Fingerabdruck) und „unknown fingerprint as the pin“ (`PODMAN_SIGNING_KEY` mit 40 Nullen), beide erwarten `.tar.gz: OK`. Der Aufruf ist `docker build -q -t <tag> --build-arg … <Kontext>` ohne `--force-rm`; eine `t.Cleanup`-Funktion steht nur im Pfad, in dem der Bau wider Erwarten gelingt (`rmi -f` des Tags). Die Teilfälle scheitern an verschiedenen Schritten des Baus (Prüfsumme gegen Signaturprüfung).
- **Der Helfer des Vorgängers:** `containerIDs` und `newContainersWithCommand` (dieselbe Datei, nach `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`); `newContainersWithCommand` nimmt das Suchwort als Parameter und liest Path, Args und Env per `docker inspect`.
- **Nicht gemessen:** ob die drei Teilfälle Zwischencontainer hinterlassen (gleicher Mechanismus ist vermutet, nicht belegt: der klassische Builder des Images `test-docker-tools` behält den Container des fehlgeschlagenen Schritts, im Vorgänger am Ollama-Test gemessen), an welcher Stelle der Wert des Pins im Container steht (Env, wie beim Ollama-Pin, oder anderswo) und ob der Docker-Daemon des CI-Runners dasselbe tut.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Templates, Lastenheft, Spezifikation** — Test-Hygiene; die Prüfung im erzeugten Dockerfile bleibt wie sie ist.
- **Das Ollama-Gegenstück** — der Vorgänger hat es geliefert und geschlossen; hier wird es nicht angefasst, es sei denn, ein gemeinsamer Helfer ändert seine Signatur (dann nur mechanisch, mit grünem Lauf des Ollama-Tests).
- **Eine allgemeine Aufräum-Hilfe für alle Docker-Tests** — erst, wenn die Messung hier und im Vorgänger denselben Weg trägt und ein drittes Vorkommen bekannt ist; sonst ein eigener Vorgang. Ob andere Tests in `internal/e2e/` Container nach fehlgeschlagenen Bauten hinterlassen, ist nicht gelesen.
- **Gate-Änderungen** (Test-Image mit buildx, `-v` im CI-Aufruf, Required Checks) — beim Projektinhaber ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)).

## 2. Definition of Done

- [x] **Mechanismus gemessen und Test-Eingriff:** Der Slice misst je Teilfall, ob und welcher Container zurückbleibt und welchen Pin er im Env (oder anderswo) trägt; dann gilt dieselbe Lösung wie im Vorgänger (`--force-rm` im Bau und Prüfung nach neuer ID und Pin, `newContainersWithCommand` mit dem Pin des Teilfalls als Suchwort) oder eine begründete andere. Der Eingriff entfernt nur Container, die nach der Vorher-Liste neu entstanden sind und den Pin des Teilfalls tragen; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen.
- [x] **Pin des Eingriffs:** Die Prüfung steht im selben Docker-Testlauf, deckt alle drei Teilfälle und wird je Teilfall (oder für den Test als Ganzes, mit Begründung) einmal rot gesehen (Eingriff entfernt) und ist danach grün.
- [ ] `make gates` grün; der berührte Docker-Test gezielt grün (der volle `make test-docker` nur, wenn der Fremdtest `NestedPodman_BuildEgress` nicht mehr rot ist, sonst mit benannter Begründung); Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` (`TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild`; Helfer `newContainersWithCommand`) | update | Aufräumen der Zwischencontainer je Teilfall (`--force-rm`, Vorher-Liste, `t.Cleanup`, Suchwort je Teilfall); der Helfer bleibt, wie er ist, oder wird begründet verallgemeinert |
| `docs/messungen/<slice>/` und `tools/messung/<name>/` | entfällt | Der Vorgänger maß den Mechanismus im Wegwerf-Kontext und führte das Protokoll als Text im Slice; hier gilt dasselbe, es sei denn, die Messung wird wiederholbar gebraucht |
| `docs/maintainer/quality.md` | unverändert | gelesen (Suche nach Aufräumen/Container/cleanup): keine Konvention zum Aufräumen von Docker-Objekten in Tests; die Treffer betreffen den Lint-Container, den Test-Container mit `--network=host` und Trivy |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |
| `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, `.github/workflows/`, `spec/` | unverändert (erwartet) | kein Produktcode, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-e2e-pin-mismatch-zwischencontainer`](../done/slice-v1-e2e-pin-mismatch-zwischencontainer.md) liegt in `done/`; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung zeigt je Teilfall einen anderen Mechanismus oder eine andere Zuordnung; dann trennt der Architect die Teilfälle.
- `in-progress` → `open` (blockiert): Das Aufräumen ist ohne eine Gate-Änderung (Test-Image, Runner-Daemon) nicht zu erreichen; Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Aufräumen entfernt einen fremden Container auf einem geteilten Host.** — **Ausgang:** weiter offen bis zur Umsetzung; Zuordnung nur über den nach der Vorher-Liste neu entstandenen Container mit dem Pin des Teilfalls, nie über Status oder Zeitfenster allein; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon sind nicht ausgeschlossen (der Vorgänger trägt dieselbe Begrenzung).
- **Ein Wechsel des Test-Images auf einen Builder mit buildx (BuildKit) macht die Prüfung vakuum-grün.** — **Ausgang:** weiter offen bis zur Umsetzung; der Vorgänger maß, dass das Image `test-docker-tools` kein buildx trägt und der klassische Builder den Container behält; BuildKit legt nach dieser Annahme keinen solchen Container an, die Prüfung fände dann nichts und bliebe grün, ohne dass ein Test es sagt. Der Slice führt dazu entweder eine Prüfung, die den klassischen Builder voraussetzt (der Test wird rot, wenn der Bau keinen Container anlegt und der Mechanismus damit nicht mehr belegt ist), oder benennt die Lücke als Beobachtung im Register. Ein Wechsel des Images ist eine Gate-Änderung und liegt beim Projektinhaber.
- **Der Docker-Daemon des CI-Runners verhält sich anders als der lokale Host.** — **Ausgang:** weiter offen; nicht gemessen, weder im Vorgänger noch hier geplant; der Slice belegt den Fall auf dem Host und nennt den CI-Lauf als nicht gemessen.
- **Die drei Teilfälle scheitern an verschiedenen Schritten, der Pin steht womöglich nicht in jedem Fall im Env des Zwischencontainers.** — **Ausgang:** weiter offen; die Messung je Teilfall entscheidet, welches Suchwort trägt. Ein Teilfall ohne belegbare Zuordnung bleibt ohne Eingriff und ist in §7 ausgewiesen.

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

**Sensoren (Implementer, Stand nach der letzten Änderung am Test):** `make gates` Exit 0. Der berührte Test gezielt im Image des Ziels `test-docker` (`-tags docker -timeout 30m -v -count=1 -run '^TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild$' ./internal/e2e/`): `PASS`, drei Teilfälle. Voller `make test-docker`: Exit 0 (`internal/e2e` 655,7 s; `NestedPodman_BuildEgress` ist darin grün); die Liste der Container-IDs des Hosts war danach gleich der vorher. `make verify-depguard` lief nicht (keine Importregel berührt).

**Abgleich Plan-Tabelle gegen `git diff --stat` (Stand der Umsetzung, vor Review):** Der Diff berührt den Test (Zeile 1) und diese Slice-Datei (Sammelzeile); `quality.md`, `docs/messungen/`, `tools/messung/`, Produktcode, `Makefile`, `repo.mk`, Workflows und `spec/` sind nicht im Diff. Review-Report, Register und Roadmap folgen mit Review und Closure.

**Nicht gemessen:** der Docker-Daemon des CI-Runners; zwei gleichzeitige Läufe desselben Tests gegen denselben Daemon (sie sähen jeweils den Container des anderen als neu, wie im Vorgänger).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/e2e` (E2E, Test-Infrastruktur) und gegebenenfalls `docs/` laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); beide Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (1×, mit dem Vorgänger angelegt; dieser Slice liefert den zweiten Beleg oder den Ausgang), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; CI-Aufruf ohne `-v`, berührt dieselbe Testgruppe). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
