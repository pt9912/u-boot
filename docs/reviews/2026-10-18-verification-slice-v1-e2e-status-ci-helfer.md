# Verification-Report: slice-v1-e2e-status-ci-helfer — 2026-10-18

**Rolle:** Verifier („Bauen wir es richtig?“, gegen Plan und DoD; kein Review, keine Reparatur).

**Gegenstand:** Slice [`slice-v1-e2e-status-ci-helfer`](../plan/planning/done/slice-v1-e2e-status-ci-helfer.md), Commit-Range `575281f..HEAD` (6 Commits: `814072a`, `db4eb13`, `b440a1c`, `fe4e5f5`, `342ac02`, `cd6b0a2`). Bezug: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md), [`ADR-0004`](../plan/adr/0004-ci-system.md). Review-Eingang: [`2026-10-18-slice-v1-e2e-status-ci-helfer`](2026-10-18-slice-v1-e2e-status-ci-helfer.md) (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO).

**Umgebung:** Colima, VM `kernel.apparmor_restrict_unprivileged_userns=0`, Docker-only (§3.1); Image `u-boot:test-docker-tools`, Rezeptur wie `repo.mk` Ziel `test-docker`.

## Verdikt

**DoD: die Liefer-Punkte sind erfüllt, die Closure-Punkte stehen aus (erwartungsgemäß, §7 ist noch leer).** Kein Befund gegen den Code. Der gemeldete Flake ist **nicht reproduzierbar** (siehe unten).

## Verification Evidence

Scope:

- Slice: `slice-v1-e2e-status-ci-helfer`; IDs siehe oben. <!-- d-check:ignore (Slice-Kennung als Name, Definition im Kopf verlinkt) -->
- Artefakte: `internal/e2e/sandbox_status_docker_test.go` (`startLabeled`, Vorab-Pull), `internal/e2e/container_id_test.go` (neu, untagged, `containerIDFromRunOutput` und Tests).

DoD-Abgleich (je Zeile gegen den Code gelesen):

- [x] **Helfer** — `startLabeled` startet `docker run -d -q --name <generiert> --label …`, erfasst stdout und stderr getrennt (`strings.Builder`, `cmd.Run`), wertet nur stdout über `containerIDFromRunOutput` aus (letzte nichtleere Zeile, `^[0-9a-f]{64}$`, sonst Fehler mit der ganzen Ausgabe) und endet bei Fehler mit `t.Fatalf` (stdout und stderr genannt). `docker pull -q docker.io/library/alpine:3` steht vor dem ersten `startLabeled` des Tests `…InspectorAgainstEngine`. Der zweite Aufruf (`second`) profitiert vom selben Pull.
- [x] **Tests des Helfers ohne Docker-Pflicht** — `TestContainerIDFromRunOutput`: 14 Fälle (nur ID, ohne Newline, Pull-Zeilen plus ID, CRLF plus Leerzeilen, Leerzeichen/Tab, nur CRLF, zwei IDs, leer, nur Blank, keine ID, Kurz-ID, ID plus Fehlerzeile, Großbuchstaben, Präfix); Datei ohne Build-Tag, `make test` fährt sie ohne Engine. Die Mutationsnachweise der DoD (ohne Hex-Prüfung, gemischte Ausgabe) stammen aus dem Review (drei erkannte Mutationen, zwei überlebende zu F-2) und gehören nach §7; ich habe die Mutationen **nicht** erneut gefahren. Die zwei überlebenden Mutanten sind mit `fe4e5f5`/`342ac02` adressiert (Fälle „spaces and tab“, „CRLF line ending only“, „two IDs“; das redundante CRLF-`ReplaceAll` ist entfernt, `TrimSpace` deckt `\r`).
- [~] **Lokaler Nachweis** — von mir erbracht und belegt (Tabelle unten): ungecachtes `alpine:3` grün, vollständiges `make test-docker` grün. Offen bleibt nur, dass die Ausgaben **in der Closure-Notiz (§7)** stehen; die Notiz ist noch nicht geschrieben.
- [x] **`make gates` grün, `make doc-trace` grün** — selbst gelaufen (Tabelle).
- [x] **Review durchgeführt, Report liegt vor** — `docs/reviews/2026-10-18-slice-v1-e2e-status-ci-helfer.md`. <!-- d-check:ignore (Slice-Kennung als Name) -->
- [ ] **Doku-Update** — nichts berührt; Aussage des Plans (§3), `harness/replay.md` und `docs/maintainer/quality.md` nennen den Helfer nicht, vom Review mit 0 Treffern bestätigt. Vermerk in §7 steht aus.
- [ ] **Closure-Notiz mit Steering-Loop-Lerneintrag** — offen (§7 leer).
- [ ] **Beobachtungs-Register** — offen (Entscheidung in §7 oder neuer Eintrag, siehe Auflagen).
- [ ] **Risiken aus §6 mit Ausgang** — offen; alle sechs Risiken tragen derzeit „offen“.
- [ ] **Drei Paarungen (Anker · Folge-Slice · Register)** — erst bei der Closure prüfbar.

Review-Findings am Code:

| Finding | Status | Beleg |
|---|---|---|
| F-1 (LOW, Plan-Tabelle) | behoben | §3 führt `roadmap.md`, `open/README.md`, `done/slice-v1-podman-aktuelle-version.md` (als anderer Vorgang) und den Review-Report; Abgleich siehe unten | <!-- d-check:ignore (Slice-Kennung als Name) -->
| F-2 (INFO, Mutanten) | behoben | drei neue Fälle in `container_id_test.go`; `ReplaceAll` entfernt |
| F-3 (INFO, Pull-Abhängigkeit) | als Risiko geführt | §6 „Pull-Fehler des Tests“, Ausgang offen |
| F-4 (INFO, Cleanup vor ID) | behoben | `t.Cleanup(rm -f -v <name>)` steht vor `docker run`; der Container trägt den generierten Namen `sbxe2e-lbl-…` |
| F-5 (INFO, Kommentar) | behoben | Doc-Kommentar nennt „container ID, read from stdout only“, keine Aussage mehr über stderr |

Plan-Tabelle §3 gegen `git diff --stat 575281f..HEAD` (§3.10):

| Datei im Diff | in §3 |
|---|---|
| `docs/plan/planning/in-progress/roadmap.md` | ja |
| `docs/plan/planning/done/slice-v1-e2e-status-ci-helfer.md` | ja | <!-- d-check:ignore (Slice-Kennung als Name) -->
| `docs/plan/planning/open/README.md`, `done/slice-v1-podman-aktuelle-version.md` | ja (anderer Vorgang, `db4eb13`) | <!-- d-check:ignore (Slice-Kennung als Name) -->
| `docs/reviews/2026-10-18-slice-v1-e2e-status-ci-helfer.md` | ja (durch den Reviewer) | <!-- d-check:ignore (Slice-Kennung als Name) -->
| `internal/e2e/container_id_test.go` | ja |
| `internal/e2e/sandbox_status_docker_test.go` | ja |

Jeder Eintrag der Tabelle mit „unverändert“ (`sandbox_devcontainer_docker_test.go`, `repo.mk`, `Dockerfile`, `integration.yml`, `replay.md`, `quality.md`) fehlt im Diff. Kein Diff ohne Eintrag, kein Eintrag ohne Diff (ausgenommen der Beobachtungs-Zeile „vom Implementer zu entscheiden“).

Sensors (selbst gelaufen):

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | Exit 0, 25,5 s; `baseline-verify` v6.18.1 OK (54 Dateien), `d-check` 356 Dateien 0 Befunde, `coverage-gate` 91,00 % (Schwelle 90 %) |
| `make doc-trace` | pass | 83 Anforderungen, 0 Waisen |
| `make doc-immutable RANGE="575281f..HEAD"` | pass | Range aufgelöst, 6 Commits, 0 Befunde |
| `make test-docker` (vollständig) | pass | `ok internal/e2e 635.727s`, Gesamtlauf 10 min 42 s (`-timeout 30m`), 0 `FAIL`; Log `/private/tmp/claude-501/-Users-dburkard-Dev-pt9912-u-boot/b510ef43-cf1e-44d6-a7c2-f018e826498b/scratchpad/v-testdocker.log` |
| `make docs-check` | pass | in `make gates` enthalten (0 Befunde); nach Anlegen dieser Datei erneut gelaufen, siehe Schluss |
| `make verify-depguard` | not run | keine Importregel berührt (nur Testcode) |

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `TestE2E_SandboxDevcontainer_StatusFiles_*` im vollständigen `make test-docker` grün |
| [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md) | `…StatusFiles_InspectorAgainstEngine` (Vertrauensgrenze gegen eine echte Engine), 28 grüne Einzelläufe (unten) |
| [`ADR-0004`](../plan/adr/0004-ci-system.md) | Gate nicht gelockert (`repo.mk`, `Dockerfile`, Workflow nicht im Diff); der CI-Beleg `integration-docker` steht aus (nach Push, braucht Freigabe) |

Carveouts: Neu: keine. Gelöst: keine. Unverändert: alle.

Nicht ausgeführt: `make verify-depguard` (kein Importregel-Bezug); Mutationsläufe (Review hat sie gefahren); CI-Lauf `integration-docker` (Push nicht freigegeben).

## Schwerpunkt Flake: `…StatusFiles_InspectorAgainstEngine` (30,05 s)

**Befund: nicht reproduzierbar.** In keinem meiner Läufe trat ein Fehlschlag auf. Alle Läufe stehen als volle Ausgabe unter `/private/tmp/claude-501/-Users-dburkard-Dev-pt9912-u-boot/b510ef43-cf1e-44d6-a7c2-f018e826498b/scratchpad/` (`sf-1.log` bis `sf-10.log`, `ia-1.log` bis `ia-10.log`, `par-1-1.log` bis `par-2-4.log`, `v-testdocker.log`).

| Serie | Läufe | `rmi alpine:3` vorher | Ergebnis | Zeit `InspectorAgainstEngine` |
|---|---|---|---|---|
| `-run StatusFiles -v -count=1`, nacheinander | 10 | Läufe 1 bis 6 | 10 grün | 5,9 bis 9,1 s (Gesamtlauf 81 bis 97 s, `StatusFiles_VFS` dominiert mit 60 bis 77 s) |
| nur `InspectorAgainstEngine`, nacheinander | 10 | Läufe 1 bis 3 | 10 grün | 6,5 bis 8,6 s |
| vier Prozesse gleichzeitig, CPU-Last (`--cpus 2`, Endlosschleifen), Runde 1 ungecacht, Runde 2 gecacht | 8 (2 mal 4) | Runde 1 | 8 grün | 6,7 bis 9,7 s |
| vollständiges `make test-docker` | 1 | nein (gecacht) | grün | im 635,7-s-Lauf |

Summe: 12 Läufe mit vorherigem `rmi` (davon 4 mit gleichzeitiger Last), 28 Läufe des Tests insgesamt, 0 Fehlschläge. Einschränkung der Lastserie: Beim Aufbau der Last blieben durch einen Namenskonflikt (`loadgen-*`) nicht alle sechs Lastcontainer in Runde 1 so auf, wie beabsichtigt; die Last war vorhanden, aber nicht exakt kontrolliert. Es ist kein weiterer Fehlerpfad dadurch sichtbar geworden.

**Code-Lesung zu den Hypothesen** (Zeitlimit 30 s):

- Der Inspector setzt `inspectorTimeout = 10 s` je Aufruf (`internal/adapter/driven/docker/inspector.go`, `FindContainer` und `ReadFile`); ein Timeout dort ergäbe ein Fehlerbild nach 10 s je Aufruf, nicht nach 30 s als Summe des Tests (der Test macht mehr als drei Inspector-Aufrufe, ein 30-s-Gesamtwert wäre nur bei drei aufeinanderfolgenden Timeouts erklärbar; unwahrscheinlich).
- Der Testkontext hat 10 Minuten; `dockerOutput` und `startLabeled` setzen kein eigenes Zeitlimit.
- Parallelität: `t.Parallel()` tragen nur die drei anderen `StatusFiles`-Tests (Zeilen 163, 232, 302 der Datei); `…InspectorAgainstEngine` ist seriell und läuft vor dem Freigeben der parallelen Tests. Keine Überlappung der `devcontainer.local_folder`-Labels: jeder Test nutzt ein eigenes `t.TempDir()`; die Container tragen generierte Namen `sbxe2e-lbl-<UnixNano base36>`. Dass `docker ps` „jüngster zuerst“ liefert und der Test das nutzt, betrifft nur Container desselben Ordners innerhalb **eines** Tests (`first`, `second`, dazwischen 1,5 s Pause).
- Alte Container desselben Labels aus Vorläufen: nach jedem meiner Läufe standen keine `sbxe2e-*`-Container oder -Volumes mehr da (geprüft); Ordner sind je Lauf neu.
- Auffällig bleibt, dass 30,05 s sehr nahe an den 30 s des Daemons liegt: Der Registry-Client des Docker-Daemons dialt mit 30 s Zeitlimit. Der Fehlschlag **nach dem `rmi`** fiel genau auf den Schritt, der Netz braucht (`docker pull -q alpine:3`, neu im Test). Das ist die plausibelste Erklärung (kurzzeitige Docker-Hub-Erreichbarkeit, Rate-Limit oder DNS), aber eine **Vermutung ohne Meldung**: Die Fehlermeldung wurde nicht gesehen, ich konnte nichts belegen.

**Empfehlung:** (1) Kein Code-Befund, kein Blocker für die Closure. (2) Beim nächsten Auftreten die vollständige Ausgabe sichern (`-v`, Datei); die Meldung unterscheidet sofort: `docker pull alpine:3: …` (Netz/Registry, Pull-Fehler, §6-Risiko F-3) gegen `seed files`/`FindContainer`/`ReadFile` (echter Fehler im Test oder Adapter). (3) Das Risiko „Pull-Fehler des Tests“ in §6 mit Ausgang **eingetreten, nicht belegt (ein Lauf von vier ohne Meldung)** oder **weiter offen** schließen und den Flake **nicht** als behoben melden. (4) Optional (Entscheidung des Projektinhabers, kein Teil dieses Slice): ein zweiter Pull-Versuch im Test oder `t.Skip` bei Registry-Fehler; das widerspricht der gewählten „ehrlichen“ Haltung aus Review F-3 und braucht einen eigenen Vorgang.

## Diskrepanzen

- Keine zwischen Plan und Code. Die `sandbox_status_docker_test.go`-Aussagen des Plans (Zeilenangaben 351/405 etc.) sind durch die Änderung um einige Zeilen verschoben; das ist Chronik, kein Befund.
- Der Slice verlangt in der DoD „Mutation … Nachweis in §7“; §7 ist leer, die Mutationsergebnisse liegen nur im Review-Report. Beim Schreiben von §7 nachzutragen (und die Behebung von F-2 als zusätzliche Fälle zu nennen).

## Auflagen für die Closure

1. **§7 schreiben:** lokaler Nachweis (diese Läufe: ungecacht grün, `make test-docker` 635,7 s, `make gates`, `make doc-trace`), Mutationsnachweis aus dem Review plus F-2-Nachtrag, Entscheidung „Tests des Helfers untagged“ (Befund: `make test` kompiliert und fährt `internal/e2e` ohne Engine, Coverage 91,0 %), Vermerk „Doku nicht berührt“, Steering-Loop-Lerneintrag (Finding-Klassen des Reviews).
2. **§6:** jedem der sechs Risiken einen Ausgang geben. „Ursache nur abgeleitet“ ist belegt (Review-Reproduktion Vorher-Stand rot, Nachher grün); „Weitere Tests mit gemischter Ausgabe“ entfallen (Review: keine); „Lokaler Nachweis ohne CI-Beleg“, „Unversionierte `docker-ce-cli`“, „Pull-Fehler des Tests“ weiter offen (der CI-Beleg ist Teil der DoD des Cuts); „Hygiene-Tests unter dem Tag“ entfallen (Tests untagged).
3. **Beobachtungs-Register:** Die Klasse „Test grün nur mit lokalem Cache“ ist ein Einzelfall dieses Slice (Zähler 1); als neue Beobachtung unter `docs/plan/planning/observations/BEO-E2E/` anlegen (Kennung nach Schema `BEO-E2E/<slug>`) oder in §7 begründet „nicht angelegt“ festhalten. Der Flake (ein Lauf von vier, Ursache unbelegt) gehört mit demselben Vorbehalt in die Beobachtung, nicht in die Behebungsaussage.
4. **Reihenfolge (§3.3):** erst Inhalt (DoD-Häkchen, §7), dann reiner `git mv` nach `done/` (`make slice-mv`). <!-- d-check:ignore (Make-Target, keine Slice-Kennung) -->
5. **Roadmap-Ruhe-Marker nach dem Move:** `roadmap.md` Sektion „Offene Wellen“ trägt derzeit den Verweis auf den laufenden Slice (Commit `b440a1c`). Nach dem Move nach `done/` muss dieser Verweis durch den Ruhe-Marker („Nichts in Arbeit.“) oder den nächsten laufenden Slice ersetzt werden, sonst meldet `planning-drift` in `make docs-check` einen Befund. Der Podman-Slice (`open/`) startet laut eigenem Plan nach diesem Slice.
6. **Offener Beleg für den Cut:** `integration-docker` grün auf dem Commit nach dem Push (Freigabe des Projektinhabers); kein Closure-Kriterium dieses Slice, aber DoD von `slice-v1-release-cut-v0.11.0`. <!-- d-check:ignore (Slice-Kennung als Name) -->

## Host-Zustand am Ende

`docker.io/library/alpine:3` gecacht (Image `294b683cb724`), keine `sbxe2e-*`-Container oder -Volumes, keine `loadgen-*`-Container; die acht fremden beendeten Container (`friendly_pasteur`, `cool_colden`, `heuristic_robinson`, `boring_wilbur`, `schema-guard-dev-pg`, `cdc-examples-*`) unangetastet. Arbeitsbaum außer dieser Datei unverändert; kein Push, kein Tag.
