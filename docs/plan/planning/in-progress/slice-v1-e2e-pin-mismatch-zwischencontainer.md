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

**Ziel:** Die Negativprobe `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` hinterlässt nach ihrem Lauf keinen Container des fehlgeschlagenen Baus auf dem Docker-Host, ohne einen fremden Container anzufassen.

### Gefundener Stand (gelesen am 2026-10-19)

- **Der Test:** `internal/e2e/sandbox_devcontainer_docker_test.go`, `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` (ab Zeile 1006): `generateSandboxProjectWith`, dann `dockerOutput(ctx, t, "build", "-q", "-t", tag, "--build-arg", "OLLAMA_SHA256_AMD64="+zeros, "--build-arg", "OLLAMA_SHA256_ARM64="+zeros, …)`; der Test erwartet einen Exit-Code ungleich 0 und die Ausgabe `FAILED`. Eine `t.Cleanup`-Funktion steht nur im Pfad, in dem der Bau wider Erwarten gelingt (`rmi -f` des Tags); für den Container des fehlgeschlagenen Baus steht keine.
- **Das Rohdatum:** [`rohdaten-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-09.txt) weist einen neuen Container nach dem Lauf aus (`status=exited exit=1`), dessen `Path`/`Args` den Pin `OLLAMA_SHA256_AMD64` tragen; das Aufräum-Skript des Messprotokolls entfernt ihn nicht, weil nichts ihn eindeutig der Messung zuordnet (`tools/messung/ollama-amd64/09-tests.sh` weist ihn aus, `11-aufraeumen.sh` löscht nur nach Präfix).
- **Nicht untersucht:** welcher Builder den Container hinterlässt (klassischer Builder oder BuildKit im Image `test-docker-tools`) und ob der Docker-Daemon des CI-Runners dasselbe tut; die Messung zeigt das Ergebnis auf einem Host, nicht den Mechanismus. Ob andere Docker-Tests in `internal/e2e/` Container nach fehlgeschlagenen Bauten hinterlassen, ist nicht gelesen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Templates, Lastenheft, Spezifikation** — der Befund ist Test-Hygiene; die Prüfsummen-Prüfung im erzeugten Dockerfile bleibt wie sie ist.
- **Das Aufräum-Skript des Messprotokolls** — der Container hat dort keine eindeutige Zuordnung; die Lösung gehört in den Test, nicht in ein Fremdlösch-Skript (Aufräum-Reichweite nur nach Präfix).
- **Gate-Änderungen** (`-v` im CI-Aufruf, Required Checks) — beim Projektinhaber ([`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md)).
- **Eine allgemeine Aufräum-Hilfe für alle Docker-Tests** — erst, wenn die Prüfung in der DoD ein zweites Vorkommen zeigt; sonst ein eigener Vorgang.

## 2. Definition of Done

- [x] **Mechanismus gemessen und Test-Eingriff:** Der Slice belegt, welcher Builder den Container hinterlässt, und räumt ihn im Test auf (Kandidaten, vom Architect zu messen und zu wählen: eine Build-Option, die Zwischencontainer auch bei einem Fehler entfernt, wo der Builder sie unterstützt; sonst die Differenz der Container-Listen vor und nach dem Bau, beschränkt auf den neu entstandenen Container mit dem Pin im `Path`/`Args`). Der Eingriff berührt keinen fremden Container.
- [x] **Pin des Eingriffs:** Ein Test oder eine Prüfung im selben Docker-Testlauf stellt sicher, dass nach `PinMismatchBreaksBuild` kein Container des fehlgeschlagenen Baus bleibt; der Fall wird einmal rot gesehen (Eingriff entfernt) und ist danach grün.
- [ ] `make gates` grün; `make test-docker` grün; Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`); Closure-Notiz mit Steering-Loop-Lerneintrag; Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist; jedes Risiko aus §6 trägt einen Ausgang; die drei Paarungen sind getragen.

## 3. Plan (vor Code)

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; der Implementer ergänzt sie in seinem ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab (`AGENTS.md` §3.10). Jede Zeile ist eine Erwartung, keine gelesene Aussage über unveränderte Dateien.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/e2e/sandbox_devcontainer_docker_test.go` (`TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild`) | update | Aufräumen des Zwischencontainers (`t.Cleanup`, Zuordnung nach Mechanismus); gegebenenfalls ein gemeinsamer Helfer neben `dockerOutput` |
| `docs/messungen/<slice>/` und `tools/messung/<name>/` | entfällt | Der Mechanismus ist im Wegwerf-Kontext gemessen; das Protokoll steht als Text in „Umsetzung (Implementer)“, kein Skript, keine Rohdaten-Datei |
| `docs/maintainer/quality.md` | unverändert | gelesen (Suche nach Aufräumen/Container/cleanup): die Datei führt keine Konvention zum Aufräumen von Docker-Objekten in Tests; die Treffer betreffen den Lint-Container (Z. 54), den Test-Container mit `--network=host` (Z. 184) und Trivy (Z. 244) |
| diese Slice-Datei, `open/README.md`, Register, Roadmap bei Start und Closure; `docs/reviews/` | update; neu | Lifecycle, Review und Verifikation |
| `internal/hexagon/**`, `internal/adapter/**`, `Makefile`, `repo.mk`, `.github/workflows/`, `spec/` | unverändert | kein Produktcode, keine Gate-Änderung |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): `Verantwortlich:` gesetzt; der Release-Cut `v0.12.0` ist geschnitten; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Das Aufräumen verlangt eine Änderung an mehreren Docker-Tests; dann trennt der Architect den Helfer vom Einzelfall.
- `in-progress` → `open` (blockiert): Das Aufräumen ist ohne eine Gate-Änderung (Test-Image, Runner-Daemon) nicht zu erreichen; Befund in §7, der Projektinhaber entscheidet.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Report liegen vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

- **Das Aufräumen entfernt einen fremden Container auf einem geteilten Host.** — **Ausgang:** weiter offen bis zum Architect-Schnitt; Zuordnung nur über den neu entstandenen Container mit dem Pin im Befehl, nie über Status oder Zeitfenster allein.
- **Der Mechanismus unterscheidet sich zwischen lokalem Host und CI-Runner.** — **Ausgang:** weiter offen; der Slice belegt den Fall auf dem Host, den CI-Lauf nennt er als nicht gemessen.

## Umsetzung (Implementer)

**Änderung:** `internal/e2e/sandbox_devcontainer_docker_test.go`: `TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild` baut mit `--force-rm`, hält vorher die Container-IDs fest (`containerIDs`) und prüft in einem `t.Cleanup`, ob ein neu entstandener Container den Null-Pin in Path, Args oder Env trägt (`newContainersWithCommand`); einen solchen entfernt es mit `docker rm -f` und meldet den Fehler. Der Pin der Prüfung ist die Zeichenkette `OLLAMA_SHA256_AMD64=` plus 64 Nullen; kein Status, kein Zeitfenster, kein Name zählt.

**Mechanismus (gelesen und gemessen auf dem Host, Docker 29.9.1, nativer `amd64`, 2026-10-19):**

- Das Image `test-docker-tools` (`Dockerfile`, Stage `test-docker-tools`) installiert nur `docker-ce-cli` und `docker-compose-plugin`; `docker buildx version` im Image antwortet `unknown command: docker buildx`, `DOCKER_BUILDKIT` ist leer. Der Bau des Tests läuft damit über den **klassischen Builder**, der den Container des fehlgeschlagenen Schritts behält.
- Wegwerf-Probe (Alpine, ein `RUN … && false`, im Image `test-docker-tools` gegen den Host-Daemon): ohne Option blieb ein Container zurück (Container-Zahl +1); mit `--force-rm` blieb keiner (Ausgabe „Removed intermediate container“, Zahl unverändert). Der Probe-Container wurde danach über seinen Befehl (`echo probe-`) identifiziert und entfernt.
- `--label` am `docker build` setzt das Label auf das Image, nicht auf den Zwischencontainer (`docker ps -a --filter label=…` fand ihn nicht). Eine Zuordnung über Label scheidet deshalb aus; die Zuordnung läuft über den Null-Pin im Env des Containers. Path/Args tragen nur den Namen der Variable im Skripttext, der Wert steht im Env (`Config.Env`).

**Wirkung vorher/nachher (Test gezielt, `go test -tags docker -timeout 30m -v -count=1 -run '^TestE2E_SandboxDevcontainer_Ollama_PinMismatchBreaksBuild$' ./internal/e2e/` im Image `u-boot:test-docker-tools`, Aufruf-Muster wie A8):**

| Lauf | Container vorher | Container nachher | Ergebnis des Tests |
|---|---|---|---|
| Gegenprobe: `--force-rm` entfernt, Prüfung aktiv | 3 | 3 | `FAIL`: „the failed build left 1 container(s) …“; der Cleanup hat den Container entfernt |
| Änderung wie geliefert | 3 | 4 | `PASS`; kein Container im Status `exited` auf dem Host |

Die Gesamtzahl taugt auf diesem geteilten Host nicht als Beleg: Ein fremder Lauf (Image `golang:1.27-alpine`) startete und beendete Container zwischen den Messungen, das erklärt die 4 nachher; der Beleg ist die Prüfung nach Pin und neuer ID, nicht die Zahl. Die Gegenprobe findet genau einen Container mit dem Null-Pin, den Zwischencontainer aus [`rohdaten-09.txt`](../../../messungen/slice-v1-devcontainer-ollama-amd64-messung/rohdaten-09.txt). Zwei Zwischenläufe mit einer Prüfung über Path/Args blieben grün, weil der Pin dort nicht steht; sie sind verworfen, die Prüfung liest das Env.

**Rote Mutation:** `--force-rm` aus dem `docker build`-Aufruf entfernen färbt den Test rot (Tabelle, erste Zeile).

**Zweites Vorkommen (nur gemeldet, nicht behoben):** `TestE2E_SandboxDevcontainer_NestedPodman_PinMismatchBreaksBuild` (dieselbe Datei) baut in drei Teilfällen absichtlich fehlschlagende Images ohne `--force-rm` und ohne Cleanup; gleicher Mechanismus, nicht gemessen. Das ist ein Folge-Vorgang (§1: allgemeine Aufräum-Hilfe erst bei zweitem Vorkommen, eigener Vorgang).

**Nicht gemessen:** der Docker-Daemon des CI-Runners und ein Host mit BuildKit-Standard (§6, zweites Risiko).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `internal/e2e` (E2E, Test-Infrastruktur) und gegebenenfalls `tools/`, `docs/` laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); alle Greenfield.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-E2E/test-gruen-nur-mit-lokalem-cache`](../observations/BEO-E2E/test-gruen-nur-mit-lokalem-cache/observation.md) (Tests im selben Paket, die vom Docker-Zustand des Hosts abhängen; der Fehler liegt im Test, nicht im Produkt), [`BEO-E2E/ci-ueberspringt-nested-podman-build`](../observations/BEO-E2E/ci-ueberspringt-nested-podman-build/observation.md) (offen; CI-Aufruf ohne `-v`). Evidenz-/Diskrepanz-Risiko: niedrig.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
