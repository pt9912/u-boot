# Verification-Report: slice-v1-devcontainer-feature-ollama — 2026-10-19 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur; kein Code geändert.

**Gegenstand:** Slice `slice-v1-devcontainer-feature-ollama` (`docs/plan/planning/in-progress/`), `ADR-0024` (Accepted, Entscheidungen 1 bis 10), Lastenheft 0.15.0 (`LH-FA-DEV-003`, `LH-FA-DEV-008`, `LH-NFA-SEC-004`). Commit-Range `1e57e82..HEAD` (7 Commits, HEAD `38d3f2a`, 20 Dateien): Implementierung `08cc5f6`, `be4feec`, `919512f`, `a5b7c2f`; Review-Report `2579e19`; Fix-Runde `38d3f2a`; fremd `054091b` (Roadmap-Zeile der Folge-Welle). <!-- d-check:ignore (Kennungen als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-19-slice-v1-devcontainer-feature-ollama.md` (F-1..F-4 LOW, F-5..F-9 INFO), DoD-Haken des Slice (fünf von zehn gesetzt), Umsetzungsstand und Sensor-Behauptung in Slice §3. <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil: GO mit Auflagen.** Alle vier Liefer-Haken der DoD (Katalog und Generator, Lauf belegt, Pins V-3 bis V-5, Doku) tragen; `make gates` ist grün, die zwei Ollama-Docker-Tests sind grün, der Rotstand von `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` ist am Ausgangsstand `1e57e82` mit identischem Fehlerbild gemessen und damit nicht vom Slice verursacht. Kein HIGH, keine DoD-Verletzung, keine Abweichung von `ADR-0024` oder vom Lastenheft. Offen sind Closure-Arbeit des Planners und drei Zustandsmarker im Slice-Text, die nach der Fix-Runde nicht mit dem Ist-Stand übereinstimmen (V-1). Der Slice ist bis zur Closure-Notiz nicht schließbar. <!-- d-check:ignore (Kennungen als Zitat) -->

## 1. Sensoren (selbst gefahren)

Alle Läufe über `make` bzw. `docker`, keine Host-Go-Werkzeuge. Docker 29.9.1 (nativ, Linux `amd64`), Netz verfügbar. Die Exit-Codes sind die der Prozesse selbst (in Log-Dateien nachgestellt: `EXIT=`-Zeile direkt nach dem Aufruf).

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass, **Exit 0** | `coverage-gate: OK — coverage 91.40% meets threshold 90%`; HEAD `38d3f2a`, Arbeitsbaum davor und danach sauber |
| `make doc-trace` | pass, **Exit 0** | „83 Anforderung(en), 0 Waise(n).“ |
| `make doc-immutable RANGE="1e57e82..HEAD"` | pass, Exit 0 | „Range aufgeloest, 7 Commit(s) — OK“; `git diff --stat` zeigt keine Datei unter `spec/`, `docs/plan/adr/`, `.github/`, kein `Makefile`, `repo.mk`, `.golangci.yml` im Diff (kein Gate gelockert, `AGENTS.md` §3.6) |
| Ollama-Docker-Tests (Aufbau siehe Abschnitt 2) | pass, Exit 0 | `TestE2E_SandboxDevcontainer_Ollama` 54,6 s, `..._Ollama_PinMismatchBreaksBuild` 28,7 s, Paket 83,3 s |
| Voller Docker-Lauf des e2e-Pakets (Aufbau wie `make test-docker`, mit `-v`) | **1 Test rot**, Exit 1 | nur `TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress` rot (Abschnitt 3); Paket `internal/e2e` 503,9 s, Gesamtlauf 512 s |
| `make docs-check` | pass, Exit 0 | „d-check: 456 Datei(en)… 0 Befund(e)“ (nach Anlage dieses Reports, siehe Abschnitt 9) |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (Diff nur im Paket `application`, `e2e`, Doku; `.golangci.yml` unverändert) |

## 2. Laufzeit der Ollama-Docker-Tests (Review F-3)

Aufruf (nicht `make test-docker` selbst, sondern dessen zwei Schritte mit einem engeren `-run`, damit nur die Ollama-Tests laufen; Build-Argumente aus dem Lauf von `make`: `GO_VERSION=1.27.2`, `GOLANGCI_LINT_VERSION=v2.14.0`, `UBOOT_VERSION=0.11.0-dev`):

```
docker build <build-args> --target test-docker-tools -t u-boot:test-docker-tools .
docker run --rm --network=host -v /Development/u-boot:/src -w /src \
  -v /var/run/docker.sock:/var/run/docker.sock -v /tmp:/tmp u-boot:test-docker-tools \
  go test -tags docker -timeout 30m -v -run '^TestE2E_SandboxDevcontainer_Ollama' ./internal/e2e/
```

| Messung | Wert | Einordnung |
|---|---|---|
| `..._Ollama` (Bau, Image-Prüfung, Server hinter Egress, Pull, Lauf, zweiter Container) | 54,6 s und 55,4 s (zwei Läufe) | **Layer-Cache warm** (das erzeugte Dockerfile ist byte-gleich zu den Läufen des Implementers; ein Cache-Treffer des Bau-Schritts ist anzunehmen, im Testprotokoll nicht ausgewiesen) |
| `..._Ollama_PinMismatchBreaksBuild` | 28,7 s und 29,0 s | geänderte `--build-arg`-Werte machen den `RUN`-Schritt zum **Cache-Fehlschlag**: Das ist die Zeit eines kalten Schritts (apt, Download 1,44 GB, Prüfsumme, Abbruch) auf diesem Host |
| Paket `internal/e2e` insgesamt | 503,9 s (8,4 min), davon die zwei Ollama-Tests 84 s | Makefile-Kommentar nennt 7 bis 11 Minuten für das Paket (Colima, arm64); `-timeout 30m` hält mit großem Abstand |
| Kombination Ollama plus `nestedRuntime` podman bzw. docker (Review F-7), Bau und Versionsprobe | 45,9 s und 52,9 s | siehe Abschnitt 6 |

Aussage und Grenze: Auf diesem Host (schneller Anschluss, Cache warm) bleibt das 30-Minuten-Limit sicher gewahrt. **Nicht gemessen** ist ein kalter Lauf auf einem CI-Runner ohne Cache; der Download beider Pfade (Positivtest und Negativtest) kann dort jeweils minutenlang dauern. Die Archivgrößen stimmen mit der Zusage: `content-length` der Release-Assets `ollama-linux-amd64.tar.zst` = 1 438 364 593 Byte (1,44 GB), `ollama-linux-arm64.tar.zst` = 1 560 414 226 Byte (1,56 GB), gelesen per HTTP-HEAD am 2026-10-19.

## 3. Rotstand `NestedPodman_BuildEgress` (Review F-6) und CI-Einordnung

Messung in einem getrennten `git worktree` (Scratchpad, danach mit `git worktree remove --force` und `prune` entfernt; `git worktree list` und `git status` nach dem Lauf sauber), Aufruf wie in Abschnitt 2 mit `-run '^TestE2E_SandboxDevcontainer_NestedPodman_BuildEgress$'`:

| Stand | Ergebnis | Fehlerbild |
|---|---|---|
| `1e57e82` (Ausgangsstand) | **FAIL**, 21,4 s | `STEP 3/3: RUN ! getent hosts example.com …` löst `example.com` auf (zwei IPv6-Adressen `2606:4700:10::…`), `Error: building at STEP …: exit status 1` |
| HEAD `38d3f2a` (im vollen Lauf) | **FAIL**, 22,9 s | dasselbe Fehlerbild, dieselben Adressen, derselbe Schritt |

Der Test ist am Ausgangsstand und auf HEAD gleichermaßen rot, mit gleicher Ursache (`getent hosts` für einen nicht freigegebenen Namen löst im verschachtelten Podman-Bau auf, hier mit IPv6-Antworten). Die Behauptung in Slice §3 („auf `1e57e82` ebenfalls rot“) ist damit **belegt**; die Änderung des Slice ist nicht die Ursache. Die Ursache selbst ist nicht untersucht (Host-Abhängigkeit, etwa IPv6 oder DNS-Weg, ist denkbar, nicht gemessen). Die Fehlerbild-Gleichheit schließt aus, dass die `mounts`-Umstellung den Test berührt.

**CI-Konfiguration** (`.github/workflows/integration.yml`, gelesen): Der Job `integration-docker` führt `make test-docker` aus (`timeout-minutes: 30`, Trigger `pull_request`, `push` auf `main`, `workflow_dispatch`). Er trägt **kein `continue-on-error`**: `grep` über `.github/` findet das Wort nur in den Kommentarzeilen 13 und 22 des Workflow-Kopfes; das Flag wurde mit Commit `8865ca13` („remove integration-docker continue-on-error flag“) entfernt. Der Kopfkommentar beschreibt es noch als geltende befristete Maßnahme. Der Review-Befund F-3 („steht laut Workflow-Kopf noch auf `continue-on-error`“) ist deshalb in diesem Punkt **unzutreffend**; der Kopf ist veraltet, der Job ist ein regulärer Job und bei dieser Rotlage rot. **Ob er ein Pflicht-Check für Merges ist**, steht nicht im Repo: `docs/maintainer/branch-protection.md` und `docs/maintainer/quality.md` führen als Required Status Checks nur `gates`, `security-gates` und `image-scan`; der Job `integration-docker` ist dort nicht genannt (die GitHub-Einstellung selbst ist nicht prüfbar). Der CI-Aufruf läuft ohne `-v`, ein Skip ist von einem Pass nicht zu trennen (so auch `ADR-0024` Entscheidung 10). Aus beidem folgt für den Slice: Der Docker-Test des Slice ist durch kein sichtbares, grünes CI-Ergebnis belegt (siehe V-2). <!-- d-check:ignore (Kennungen als Zitat) -->

## 4. Plan-vs-Diff (Plan-Tabelle §3 gegen `git diff --stat 1e57e82..HEAD`, `AGENTS.md` §3.10)

20 Dateien im Diff bei HEAD `38d3f2a` (954 Einfügungen aus dem Review-Stand plus die Fix-Runde und den Review-Report; mit diesem Report wächst die Range auf 21 Dateien).

| Diff-Datei(en) | Plan-Zeile | Befund |
|---|---|---|
| `devcontainer_features.go` | Zeile 1 | ok; Marker `dockerfileBlock`, `projectFeatureSourceVersion` mit `ok=false`, `validateDevcontainerFeatures` mit Exit 10, `thirdPartyFeatureWarnings` unberührt: am Diff gelesen |
| `devcontainer_egress.go` | Zeile 2 | ok; genau die zwei Namen |
| `templates.go` | Zeile 3 | ok, mit benannter Abweichung (Schalter `Ollama` und `Mounts` statt Feldern; die drei Volume-Felder entfallen): Diff bestätigt, kein Verweis auf die entfernten Felder im Baum |
| `devcontainer_sandbox.go` | Zeile 4 | ok, mit benannter Abweichung: `sandboxRelaxationWarnings` unverändert (nicht im Diff), die Meldungen stehen in `ollamaWarnings` (`devcontainer_features.go`); `catalogueFeatureActive` wird nicht benutzt |
| `Dockerfile.tmpl` | Zeile 5 | ok (Abschnitt 5, Entscheidung 1 bis 3, 6, 8) |
| `devcontainer.json.tmpl` | Zeile 6 | ok; zusammensetzbare `mounts`, Fälle Standard/Podman/Docker als Literale gepinnt |
| `generate.go` | Zeile 7 | ok; eine Zeile, `ollamaWarnings` neben `thirdPartyFeatureWarnings` und `agentNodeWarnings` |
| `devcontainer_features_test.go`, `export_test.go` | Zeile 8 | ok |
| `generate_ollama_test.go` (neu) | Zeile 8 (Abweichung benannt) | ok; die in Zeile 8 genannten `generate_features_test.go` und `generate_sandbox_test.go` sind **nicht** im Diff (Plan-Eintrag ohne Diff, in Slice §3 „Neu gegenüber der Tabelle“ benannt) |
| `sandbox_devcontainer_docker_test.go` | Zeile 9 | ok |
| `benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/replay.md` | Zeile 10 | ok |
| `cmd/uboot/main_test.go` | Zeile 14 (seit Fix-Runde) | ok |
| `roadmap.md` | Zeile 15, „fremd“ (seit Fix-Runde) | ok (`054091b`) |
| `slice-v1-devcontainer-feature-ollama.md` | Zeile 13 | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/reviews/2026-10-19-slice-v1-devcontainer-feature-ollama.md` | Zeile 12 (`docs/reviews/`) | ok | <!-- d-check:ignore (Kennungen als Zitat) -->
| `docs/plan/planning/observations/BEO-APP/…` | Zeile „entfällt“ | bestätigt: nicht im Diff |

Kein Diff ohne Plan-Zeile. Beide Fix-Runden-Änderungen des Reviews (F-4: `main_test.go` und `roadmap.md` in der Tabelle) sind eingetragen. **Aber:** der Abgleich in Slice §3 „Umsetzungsstand“ trägt weiter den Stempel „(Stand nach der Umsetzung, vor Review)“ und führt `cmd/uboot/main_test.go` unter „Neu gegenüber der Tabelle“, obwohl die Tabelle die Zeile jetzt hat; ein Abgleich „nach der letzten Fix-Runde“ (so verlangt die Zeile direkt unter der Tabelle und `AGENTS.md` §3.10) ist nicht vermerkt (V-1).

Fix-Runde `38d3f2a` gegen die Review-Befunde: F-1 behoben (Doc-Kommentar nennt `TestRun_Generate_OllamaNotices`; der alte Name `TestRun_Generate_Quiet_OllamaNoticesDropped` kommt im Baum nur noch im Review-Report als Zitat vor), F-2 behoben (die Sammelaussage „drop a new entry … done“ und „slug is `1` across the board“ sind durch die geltende Zusage ersetzt; `dockerfileBlockFeatureActive` existiert als Verweisziel), F-4 behoben (Tabelle). Kein Kommentar der Fix-Runde trägt eines der Muster aus `AGENTS.md` §3.9 (a) bis (c).

## 5. DoD-Abgleich (Slice §2) und Entscheidungs-Konformität

| DoD-Zeile | Stand | Beleg (vom Verifier gelesen oder gefahren) |
|---|---|---|
| **Katalog und Generator** (Haken gesetzt) | **erfüllt, der Haken trägt** | `ollama` als Katalogschlüssel mit Marker (`devcontainer_features.go`); Pin als `ARG` an einer Stelle, `sha256sum -c` vor `tar --zstd`, in derselben `&&`-Kette (`Dockerfile.tmpl`); `egressFeatureHosts()` mit den zwei Namen; Volume-Mount, kein `postStartCommand`-Anteil; Golden Case: `TestGenerateDevcontainer_Ollama_Sandbox`, `_EgressHostsAndTogetherWithClaudeCode`, `_MountsComposition` (bestehende Konfigurationen byte-identisch), `_NotActive`, `_VersionOverrideWithoutSourceRejected`, `TestDoctor_FeaturesDrift_Ollama`; `make gates` grün |
| **Lauf belegt** (Haken gesetzt) | **erfüllt, der Haken trägt** | `TestE2E_SandboxDevcontainer_Ollama` baut das erzeugte Dockerfile, prüft Version (`0.40.6`), Bibliothekspfad, keine GPU-Verzeichnisse, kein `zstd`, Besitzer des Modellspeichers, kein `OLLAMA*` in der Umgebung, startet den Server hinter Egress, prüft `LISTEN:127.0.0.1:11434`, Container-Adresse verweigert, `library=cpu`, Pull von `smollm2:135m`, `ollama run`, fremder Host gesperrt, Modell im zweiten Container; selbst gefahren, grün (Abschnitt 2) |
| **Pins V-3 bis V-5** (Haken gesetzt) | **erfüllt, der Haken trägt** | (a) `--quiet`: `TestRun_Generate_OllamaNotices` (Nutzerpfad, echte Adapter, prüft Text ohne `--quiet` und Abwesenheit des Wortes `ollama` unter `--quiet`) und `TestGenerateDevcontainer_Ollama_NoticesAreNotRelaxations`; (b) kein `OLLAMA_HOST`, `forwardPorts`, `ENV`, `EXPOSE`, `CMD`, `ENTRYPOINT`, `ollama serve` in den erzeugten Dateien (Golden Case) und Docker-Test „Container-Adresse verweigert, `127.0.0.1:11434` antwortet“; (c) `TestGenerateDevcontainer_Ollama_DefaultProfileGeneratesNothing` und der Fall `default profile` in `TestRun_Generate_OllamaNotices` |
| **Doku** (Haken gesetzt) | **erfüllt, der Haken trägt** | Handbuch (Katalogtabelle, Egress-Tabelle, Abschnitt *Ollama (optional)*, Volume-Namen-Absatz), `README.md` und `README.de.md` („14 Features“ stimmt zu 13 plus `ollama`), `CHANGELOG.md` `[Unreleased]`, Golden-Zeile in `harness/replay.md`; Nutzer-Doku ohne Kennungen (`AGENTS.md` §5 Regel 6) |
| `make gates` grün (Haken gesetzt) | **erfüllt** | Abschnitt 1 |
| Review, Report liegt vor | **erfüllt**, Haken noch nicht gesetzt | Report vom 2026-10-19 liegt vor; Fix-Runde `38d3f2a` danach (Closure-Schritt des Planners) |
| Doku-Update öffentlicher Vertrag (Haken gesetzt) | erfüllt | s. o. |
| Closure-Notiz mit Steering-Loop-Eintrag, Register, §6-Risiken mit Ausgang, drei Paarungen | **offen, zu Recht** | §7 trägt den Platzhalter; §6-Risiken stehen auf „offen bis zur Closure“ (Auflage 2) |

### Entscheidungs-Konformität `ADR-0024` (Accepted, unverändert im Diff)

| Nr. | Entscheidung | Stand | Beleg |
|---|---|---|---|
| 1 | gepinntes Archiv, SHA-256 vor dem Entpacken, `ARG`-Pin, `zstd` entfernt, kein Skript | konform | `Dockerfile.tmpl`: `ARG OLLAMA_VERSION=0.40.6` und beide Prüfsummen identisch mit der ADR; `sha256sum -c` vor `tar`; `apt-get purge -y --auto-remove zstd`; Negativtest grün (Bau bricht mit Prüfsumme null ab); Test `NO_ZSTD` |
| 2 | Katalogeintrag ohne OCI-Quelle (Marker), nur Sandbox-Profil, `version:` ohne `source:` Exit 10, `source:` macht ihn extern | konform | Diff `devcontainer_features.go`; Tests `_NotActive`, `_VersionOverrideWithoutSourceRejected`, `TestRun_Generate_OllamaNotices` (Fall `version without source`, Exit 10) |
| 3 | eigene Meldung, `LH-NFA-SEC-004`, ohne `SecurityRelaxation`, entfällt unter `--quiet`; ohne Sandbox die ungekennzeichnete „keine Wirkung“-Meldung, kein `doctor`-Eintrag | konform, mit **benannter Abweichung im Wortlaut** | `ollamaWarnings` (Level `warn`, Subjekt `.devcontainer/devcontainer.json`); der Text nennt die Version nicht selbst, sondern `OLLAMA_VERSION` im Dockerfile (ADR: „Quelle und Version“; Slice §3 nennt die Abweichung mit Begründung, Pin an einer Stelle) | <!-- d-check:ignore (Kennungen als Zitat) -->
| 4 | kein Auto-Start | konform | kein `postStartCommand`-Anteil, kein `CMD`/`ENTRYPOINT` (Golden Case) |
| 5 | nur Loopback, `OLLAMA_HOST` nie gesetzt, kein `forwardPorts` | konform | Golden Case und Docker-Test (`ss`, `curl` gegen Container-Adresse) |
| 6 | GPU-Verzeichnisse ausgelassen | konform | `--exclude` für `cuda_*` und `vulkan`; Docker-Test `NO_GPU_LIBRARIES` und `LIBRARY_CPU` **grün auf `amd64`** (die ADR führt Muster auf `amd64` als nicht gemessen; der Lauf des Verifiers und des Implementers deckt das nativ ab, ohne dass Handbuch oder CHANGELOG mehr zusagen) |
| 7 | zwei feste Hosts | konform | `egressFeatureHosts()`; Handbuch nennt Symptom und Weg bei Wechsel |
| 8 | benanntes Volume, `mkdir`/`chown`, Zusage nur „Löschen des Containers“ | konform | Mount `…-ollama-${devcontainerId}`; `STORE_OWNED_BY_VSCODE`; `MODEL_SURVIVES_CONTAINER_REMOVAL` (zweiter Container, gleiche Volumes); Neubau nicht zugesagt |
| 9 | Pin-Pflege per Slice | konform | Pin an einer Stelle, `strings.Count`-Wächter im Test; Register unberührt, Plan-Zeile „entfällt“ (Beleg liegt aus dem Mess-Slice) |
| 10 | `amd64` gepinnt, ungelaufen, nicht verifiziert | konform | Handbuch, CHANGELOG, Lastenheft führen `amd64` als nicht verifiziert; beide Docker-Tests liefen für diesen Report nativ auf `amd64` grün, das ist **kein CI-Beleg** im Sinn der Entscheidung (CI ohne `-v`, siehe Abschnitt 3) |

### Lastenheft 0.15.0

`LH-FA-DEV-003` (Absatz Ollama): jede Zusage trägt einen Test oder eine Messung: „wählbar“ und „wirkt im Sandbox-Profil“ (Generator-Tests), „nur CPU“ (`LIBRARY_CPU`), „gepinntes Archiv, geprüfte SHA-256“ (Golden Case, Negativtest), „nicht signiert, weist aus“ (Meldungstest), „startet nicht von selbst“ (Golden Case), „nur Loopback“ (Docker-Test), „Volume, überstehen das Löschen des Containers“ (Docker-Test), „nicht zugesagt“: GPU, Host-Ollama, andere Architekturen, Neubau, Modellverwaltung (nicht im Code und nicht in Handbuch oder CHANGELOG als Zusage). `LH-FA-DEV-008`: die zwei Hosts in der Default-Allowlist (Golden Case und Docker-Test mit gesperrtem fremden Host); die Stabilität des zweiten Namens bleibt als Gate F-1 des CR-Reports beim Projektinhaber (nur Hinweis). `LH-NFA-SEC-004`: kein Fremdskript im Pfad (Dockerfile-Test über alle Nicht-Kommentar-Zeilen). Das Lastenheft ist im Diff unverändert. <!-- d-check:ignore (Kennungen als Zitat) -->

## 6. Stichprobe der Zusagen aus Handbuch, CHANGELOG und READMEs gegen Test oder Messung

| Zusage | Beleg | Befund |
|---|---|---|
| zwei Hosts | `egressFeatureHosts()` (Diff), Golden Case, Docker-Test `PULL_OK` mit `FOREIGN_BLOCKED` | belegt |
| „nur Loopback“ | Docker-Test: `LISTEN:127.0.0.1:11434 ` (genau diese Adresse), `LOOPBACK_ANSWERS`, `CONTAINER_ADDRESS_REFUSED`; kein `OLLAMA*` in der Image-Umgebung (`NO_OLLAMA_ENV`) | belegt; Grenze (Review F-8): `CONTAINER_ADDRESS_REFUSED` steht für jeden Fehlschlag von `curl -m 3`, auch Verwerfen; die Bedingung trägt `LISTEN:127.0.0.1`, das Handbuch sagt „abgelehnt“: unschädlich |
| „Modelle überstehen das Löschen des Containers“ | Docker-Test: Container entfernt, zweiter Container mit denselben Volumes, `MODEL_SURVIVES_CONTAINER_REMOVAL` | belegt |
| „Überstehen des Image-Neubaus nicht gemessen“ | Handbuch, CHANGELOG nennen es als Nicht-Zusage; kein Test behauptet das Gegenteil | konform |
| „amd64 nicht verifiziert“ | Handbuch, CHANGELOG, Lastenheft; der Docker-Test lief für diesen Report auf `amd64` grün, kein CI-Beleg | konform; die Aussage bleibt richtig, bis ein CI-Lauf mit sichtbarem Pass vorliegt |
| Archivgröße 1,56 GB (`arm64`) und 1,44 GB (`amd64`) | HTTP-HEAD der Release-Assets: 1 560 414 226 und 1 438 364 593 Byte | belegt (gerundet auf zwei Stellen) |
| „bereit nach höchstens zwei Sekunden“, „66 Sekunden `max retries exceeded` ohne Blob-Host“, „`ollama.com` gesperrt: zwei Warnzeilen“ | stehen im Handbuch als **gemessen**; die Quelle ist die Messung des Mess-Slice (`ADR-0024`), nicht der Docker-Test (dieser wartet bis zu 30 s auf den Server und sperrt `ollama.com` nicht gesondert) | nicht neu gemessen (Rohdaten des Mess-Slice nicht erneut gelesen); als Messwert gekennzeichnet, kein Test-Versprechen | <!-- d-check:ignore (Kennungen als Zitat) -->
| „ohne Sandbox-Profil: kein Block, kein Volume, Meldung“ | `TestGenerateDevcontainer_Ollama_DefaultProfileGeneratesNothing`, `TestRun_Generate_OllamaNotices` | belegt |
| „Meldungen entfallen unter `--quiet`“ | `TestRun_Generate_OllamaNotices` (zweiter Lauf mit `--quiet`) | belegt (Lesart der ADR; Gate F-2 des CR-Reports laut Auftrag vom Projektinhaber entschieden, nur Hinweis) |
| `version` ohne `source` Exit-Code `10` | `TestRun_Generate_OllamaNotices` (Fall `version without source`, Exit 10) | belegt |
| „14 Features“ in beiden READMEs | `featureCatalogue()` hat 13 plus `ollama` (Review: nachgezählt) | belegt |
| Kombination mit `nestedRuntime` podman oder docker (Review F-7, im Handbuch nicht zugesagt) | **vom Verifier gemessen**: temporärer Test in einem getrennten Worktree (danach entfernt, nicht committet): Bau aus dem erzeugten Dockerfile mit `ollama` plus `nestedRuntime: podman` bzw. `docker`; `ollama --version` meldet `0.40.6`, `podman version 6.1.4` bzw. `Docker version 29.9.1` (Client und `dockerd`), `fuse-overlayfs`, `newuidmap` bzw. `iptables` vorhanden | Bau und Versionsprobe grün. **Nicht gemessen:** Lauf von Podman oder Docker in der Kombination (kein Start, kein Build im Container). `apt-get purge --auto-remove zstd` hat in diesen Images die geprüften Pakete nicht entfernt; `slirp4netns` ist im Template nicht vorgesehen (nur `passt`), also kein Verlust |

## 7. Review-Findings F-1..F-9 einzeln

| Finding | Stand | Verifier-Beleg und Behandlung |
|---|---|---|
| F-1 LOW Verweis auf nicht existierenden Test | **behoben** | `38d3f2a`; Suche über den Baum findet den alten Namen nur im Review-Report als Zitat |
| F-2 LOW Katalog-Kommentare | **behoben** | `38d3f2a`, Kommentar am Diff gelesen |
| F-3 LOW Docker-Test ohne Laufzeitmessung | **gemessen, offen als Risiko** | Abschnitt 2: Paket 503,9 s warm, zwei Ollama-Tests 84 s, kalter `RUN`-Schritt 29 s; kalter Lauf auf Runner ohne Cache nicht gemessen. **Berichtigung der Befund-Begründung:** der Job steht nicht auf `continue-on-error` (Abschnitt 3). Ausgang für Slice §6 Risiko „Docker-Test langsam“: weiter offen (kein Workflow-Eingriff, `AGENTS.md` §3.6) |
| F-4 LOW Plan-Abgleich unvollständig | **behoben** in der Tabelle (Zeilen `main_test.go`, `roadmap.md`); der Stand-Vermerk im Umsetzungsstand ist veraltet (V-1) |
| F-5 INFO Gate F-2 (`--quiet`) | **übergeben** | Entscheidung des Projektinhabers laut Auftrag am 2026-10-19; Slice-Text führt sie noch als offen (V-1) |
| F-6 INFO Rotstand am Ausgangsstand ungeprüft | **belegt** | Abschnitt 3 |
| F-7 INFO Kombination ohne Bautest | **teilweise gemessen** | Abschnitt 6: Bau und Versionsprobe grün; Lauf nicht gemessen. Als Beobachtung/Risikoausgang für die Closure, kein Pflicht-Test des Slice |
| F-8 INFO Probe unterscheidet Ablehnung nicht | **begründet offen** | kein Änderungszwang; `LISTEN:127.0.0.1:11434` trägt die Zusage |
| F-9 INFO `-tags docker` nicht in Lint und Unit-Test | **begründet offen, Repo-Struktur** | der Docker-Test wurde für diesen Report übersetzt und gefahren (grün); Lint und `make test` übersetzen ihn nicht |

## 8. Findings des Verifiers

| ID | Klasse | Befund |
|---|---|---|
| V-1 | LOW | **Zustandsmarker im Slice stimmen nach der Fix-Runde nicht mit dem Ist-Stand** (`AGENTS.md` §3.10, Marker mit Ablaufpunkt): (a) Slice §3 Umsetzungsstand, erster Punkt: „Abgleich … (Stand nach der Umsetzung, vor Review)“ und „Der Abgleich ist nach jeder Fix-Runde zu wiederholen“ ohne vermerkten Abgleich nach `38d3f2a`; (b) zweiter Punkt „Neu gegenüber der Tabelle: … `cmd/uboot/main_test.go`“, obwohl die Tabelle die Zeile seit `38d3f2a` trägt; (c) DoD Zeile „Pins aus der Verification“ und Umsetzungsstand „Offen: Gate F-2 (Projektinhaber)“ führen das Gate als offen, die Entscheidung liegt nach Auftrag seit dem 2026-10-19 vor (Lesart der ADR unter `--quiet` bestätigt); der Absatz „Umsetzungsstand“ spricht zudem von „beide Tests grün“, ohne den Lauf des Verifiers zu kennen. Auflage 1. |
| V-2 | LOW | **Kein sichtbarer grüner CI-Beleg für den Docker-Test, und der Kopfkommentar des Workflows ist veraltet:** `integration.yml` beschreibt `continue-on-error: true` als geltende Maßnahme, das Flag ist seit `8865ca13` entfernt; `integration-docker` steht nicht in den dokumentierten Required Status Checks; der Aufruf läuft ohne `-v`. Gleichzeitig ist `NestedPodman_BuildEgress` auf dem gemessenen Host (und am Ausgangsstand) rot, damit wäre das Job-Ergebnis dort rot, unabhängig vom Slice. Konsequenz: Für `amd64` bleibt der Lauf „nicht verifiziert“ (richtig so, `ADR-0024` Entscheidung 10); der Slice darf `amd64` nicht als verifiziert nachführen. Der Workflow-Kommentar und das Fehlen von `-v` sind Gate-Gegenstand und liegen beim Projektinhaber (`AGENTS.md` §3.6), nicht in diesem Slice. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-3 | INFO | Laufzeitmessung ist eine Messung mit warmem Cache (Abschnitt 2); eine Aussage über CI-Zeiten folgt daraus nicht. |
| V-4 | INFO | Wortlaut der Lieferketten-Meldung nennt die Version nicht, `ADR-0024` Entscheidung 3 nennt „Quelle und Version“; Abweichung im Slice benannt und begründet (Pin an einer Stelle). Kein Handlungsbedarf, wenn die Closure sie als bewusste Abweichung führt. | <!-- d-check:ignore (Kennungen als Zitat) -->
| V-5 | INFO | Ursache des Rotstands von `NestedPodman_BuildEgress` nicht untersucht (Fehlerbild: IPv6-Auflösung von `example.com` im verschachtelten Bau). Ein Register-Eintrag existiert für das CI-Überspringen (`BEO-E2E/ci-ueberspringt-nested-podman-build`); ob dieser Beleg dort anschließt, entscheidet der Planner. |

Kein HIGH, keine DoD-Verletzung, keine Abweichung von `ADR-0024`. <!-- d-check:ignore (Kennungen als Zitat) -->

## 9. Auflagen (alle beim Planner, vor `done/`)

1. **Zustandsmarker auf den Ist-Stand (V-1):** Abgleich der Plan-Tabelle gegen `git diff --stat` nach der letzten Änderung wiederholen und vermerken (diesen Report und die Closure-Notiz eingeschlossen); die Aussage „Neu gegenüber der Tabelle: `cmd/uboot/main_test.go`“ und den Stempel „vor Review“ berichtigen; Gate F-2 in DoD-Zeile und Umsetzungsstand auf „entschieden am 2026-10-19 (Lesart der ADR unter `--quiet`)“ setzen, Gate F-1 des CR-Reports entsprechend im Register der Closure.
2. **Closure-Notiz §7 mit Steering-Loop-Lerneintrag, Register, §6-Risiken mit Ausgang** (Vorschlag aus dem Messstand: Marker-Typ: entfallen/getragen; `mounts`-Umbau: entfallen, byte-identisch gepinnt; Docker-Test langsam: weiter offen, kalter Lauf nicht gemessen; Pin altert: weiter offen; Bedingung ohne Entsprechung: entfallen, jede Zusage hat Test oder Messung, Ausnahmen als „gemessen“ gekennzeichnet; Hosts wechseln: weiter offen, Gate F-1 laut Auftrag entschieden), DoD-Haken (Review, Closure, Register, Risiken, Paarungen) setzen. Beobachtungen für das Register: der veraltete Kopfkommentar von `integration.yml` (V-2) und die Kombination `ollama` plus `nestedRuntime` mit Bau, aber ohne Lauf (F-7) als offene Punkte.
3. **Validator-Hinweis aus Slice §5:** der Nutzerpfad (`config set … enabled true`, `generate devcontainer`, Container starten, Modell holen) ist nachzuspielen; der Verifier hat den Teil „Bau und Lauf im Container“ durch den Docker-Test, nicht den Nutzerpfad über die Dev-Container-CLI geprüft.
4. **Hinweis an den Projektinhaber (nicht Slice-Gegenstand):** Kopfkommentar von `.github/workflows/integration.yml` und das Fehlen von `-v` im CI-Aufruf (V-2) sind Gate-Gegenstand; ohne sie bleibt `amd64` nicht verifiziert.
5. `git mv` nach `done/` als reiner Move nach der Closure-Notiz (`AGENTS.md` §3.3), Roadmap-Arbeitszeile zurücksetzen.

## 10. Nicht verifiziert

Kalter Lauf der Ollama-Tests auf einem Runner ohne Layer-Cache und mit CI-Netz; ein CI-Lauf mit sichtbarem Pass auf `amd64`; Docker Desktop, Linux ohne VM, Colima/`arm64` (der Messstand des Lastenhefts; der Verifier lief auf nativem Docker 29.9.1, `amd64`); die Rohdaten des Mess-Slice (Zwei-Sekunden-Start, 66 Sekunden, `ollama.com`-Warnzeilen); Lauf von Podman oder Docker in der Kombination mit `ollama`; Start des Servers über die Dev-Container-CLI; die Ursache des Rotstands von `NestedPodman_BuildEgress`; die GitHub-Einstellung der Required Status Checks; `make verify-depguard` (keine Importregel berührt); die Delegation des Projektinhabers zu Gate F-1 und F-2 (Gespräch, nicht prüfbar, laut Auftrag nur Hinweis).

## Sensor-Belege (Verification-Evidence-Block, `harness/verification.md`)

Scope: Slice `slice-v1-devcontainer-feature-ollama`; IDs `LH-FA-DEV-003`, `LH-FA-DEV-008`, `LH-NFA-SEC-004`, `ADR-0024`. <!-- d-check:ignore (Kennungen als Zitat) -->

Carveouts: neu keine, gelöst keine, unverändert alle (nicht im Diff).

Commit / Artefakt: HEAD `38d3f2a`; dieser Report.
