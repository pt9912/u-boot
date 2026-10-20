# Verification-Report: slice-v1-devcontainer-start-fundament — 2026-10-20

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD und Entscheidungen). Kein Review, keine Validierung, keine Reparatur: es wurde kein Produktionscode und kein Planungsdokument geändert, nichts committet, keine Closure ausgeführt.

**Gegenstand:** [`slice-v1-devcontainer-start-fundament`](../plan/planning/done/slice-v1-devcontainer-start-fundament.md) (`in-progress/`), `HEAD` = `origin/main` = `5d63e52`, Diff-Basis `eea0c92~1..HEAD` (alle Commits des Slice, `eea0c92` bis `5d63e52`, 23 Dateien, 3139 Einfügungen). Gelesen: Slice §1 bis §8 (Implementer-Übergabe, Review-Fix-Runde, r2), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) (beide `Accepted`), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) im Lastenheft, Port, Adapter, `secret_mask.go`, die drei eingebetteten Dateien, die Docker-Testdatei (Testnamen, Skip-Suche), das Messprotokoll.

**Eingang:** die Review-Reports [`2026-10-20-slice-v1-devcontainer-start-fundament`](2026-10-20-slice-v1-devcontainer-start-fundament.md) (F-1 bis F-10) und [`…-r2`](2026-10-20-slice-v1-devcontainer-start-fundament-r2.md) (R-1 bis R-10) als Eingang, nicht als Beweis. Die Behebung wurde am Ist-Stand des Codes nachgelesen (Abschnitt 3).

**Gesamturteil: GO (mit Auflagen für die Closure).** Alle drei Liefer-Punkte der DoD sind im Code und in den Belegen getragen; `make gates` ist von mir grün gefahren; der Nachweis für `make test-docker` (CI-Lauf für `5d63e52`) ist belegt. Kein HIGH. Ein MEDIUM (V-1: Zusage „nie auf einer Kommandozeile“ gegen den gemessenen Befund, förmliche Entscheidung offen) und mehrere LOW zu Plan-Aussagen und Zustandsmarkern, die der Planner vor oder bei der Closure berichtigen sollte.

## 1. Sensoren

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (selbst gefahren, sauberer Arbeitsbaum) | **pass, Exit 0** | `baseline-verify: v6.18.1 OK — 54 Dateien`; `d-check: 562 Datei(en) geprüft, 0 Befund(e)`; `lint`: `0 issues`; `coverage-gate: OK — coverage 91.60% meets threshold 90%`; `ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege` |
| `make test-docker` (CI, nicht lokal) | **pass** | GitHub-Lauf `integration-docker` 39409233721, `headSha` = `5d63e524c347c8c61a5a7ce6c8fae496b0fdfc7e`, Konklusion `success`, kein Schritt rot; Job-Log: `ok …/internal/adapter/driven/docker 152.031s`, `ok …/internal/e2e 593.785s`, `ok …/internal/hexagon/port/driven`; `ci` für denselben Commit (39409233725) `success` |
| Stichproben-Mutation (selbst, unabhängig vom Implementer) | pass (Wächter rot wie behauptet) | in einer Kopie von `HEAD` außerhalb des Repos: `--rm` aus `runArgs` entfernt und `-e NAME` → `-e NAME=x`: `TestRunArgs/defaults`, `TestRunArgs/env_by_name,_recreate,_read-only` und `TestLauncher_Up_PassesEnvByNameOnly` rot (Image `u-boot:test-docker-tools`, `go test -run`); das Repo blieb unberührt (`git status` leer) |
| `make verify-depguard` | **nicht gelaufen** | `.golangci.yml` im Diff unberührt; der Port importiert nur Stdlib (`context`, `errors`, `fmt`, `io`, `strings`), der Adapter liegt unter `adapter/driven` |
| `make test-docker` lokal | **nicht gelaufen** | Host-Netz (TLS zu `download.docker.com`), laut Auftrag und Slice §7 kein Nachweis; ersetzt durch den CI-Lauf oben |
| `make ci` / `make fullbuild` | **nicht gelaufen** | `ci` war als GitHub-Lauf für `5d63e52` grün; Security-/Release-Pfad nicht berührt, kein lokaler Lauf nötig |

Zur Reichweite des CI-Belegs: Der Lauf ruft `go test` ohne `-v` auf, ein Skip wäre im Log nicht pro Test zu trennen (der Slice kennt das, §5). Gegenbelege: die Docker-Testdatei hat **keinen** `t.Skip`-Pfad (Suche über beide Launcher-Testdateien: kein Treffer; `requireDocker` ruft `t.Fatalf`), und die Paketlaufzeit von 152 s passt zu den ausgeführten Docker-Tests (126,9 s lokal bei warmem Cache, dazu der Pull auf dem Runner). Das trägt die Zusage „nicht übersprungen“ hinreichend (siehe V-3).

## 2. DoD-Abgleich (Slice §2, Zeile für Zeile)

Die Kästchen im Slice sind alle leer; für einen Slice in `in-progress/` ist das der erwartete Zustand (Häkchen setzt der Planner bei der Closure). Bewertung nach Ist-Stand:

| # | DoD-Zeile | Urteil | Beleg |
|---|---|---|---|
| 1 | Port, Adapter, Sentinel, eingebetteter Kontext | **erfüllt** | `internal/hexagon/port/driven/devcontainer_launcher.go`: `DevcontainerLauncher` (`Up`, `Stop`), `ErrDevcontainerRuntime`, Optionen, Ergebnis, zwei Prüffunktionen. `internal/adapter/driven/docker/devcontainer_launcher.go`: eigene Vorab-Prüfung (`preflight`: Binary, `Probe.Info`, `requireLocalEngine`; nicht `Engine.preflight`), Bau nur bei fehlendem Tag (`EnsureImage`, `docker image inspect`, dann `docker build --force-rm -t <tag> -` mit Tar-Strom), `runArgs` mit `--rm`, `--name`, zwei `--mount`, `-e <NAME>` ohne Wert, Spiegelpfad, kein `--privileged`/`-v`/`--network`; Abbruch: `cmd.Cancel` = SIGINT, `WaitDelay` = Frist, `removeHelper` (`docker rm -f <Name>`); Port-Pin `var _ driven.DevcontainerLauncher = (*Launcher)(nil)`. Pins: `devcontainercli/Dockerfile` (`FROM node:22-bookworm@sha256:0e5f9065…`, `docker:29-cli@sha256:6fae2214…` als Stufe mit `COPY --from`, `RUN npm ci --ignore-scripts`), `package.json` (exakt `0.89.0`), `package-lock.json` (eine `integrity`). `--ignore-scripts` trägt (Docker-Tests S5 grün in CI), die Rückführung an den Architect entfällt. Einschränkung: V-2 (Digests im Netz nicht neu gelesen; durch den grünen Bau in CI indirekt belegt) |
| 2 | Docker-Tests mit Wächter, Unit-Tests | **erfüllt** | `devcontainer_launcher_docker_test.go` (`//go:build docker`): S2 (`…S2_LabelCarriesHostPath`, Wächter `…GuardMeasurementPathDeviates`), S6 Socket (`…SocketNotInStartedContainer`, Wächter `…GuardHandWrittenSocketMountIsDetected`, `…InitializeCommandRunsWithSocket`), S6 Lebensdauer (`…Lifetime`), S6 Umgebung (`…EnvReachesOnlyPassedNames`, `…ValueIsNeverInImageLabelFileOrCommandLine`, Wächter `…GuardCommandLineScanFindsValueInArguments`, `…ContainerEnvValueIsMasked…`, die zwei `…Limit…`-Tests), S6 `readonly`, S5 (`…BuildFromEmbeddedContext`, Wächter `…GuardFalsifiedIntegrityBreaksTheBuild`), S7 (`…S7_EngineOfTheRunner`), `…StopStopsButKeepsTheContainer`. Unit: `devcontainer_launcher_test.go` (Argumentliste, Vorab-Prüfung, Ergebniszeile, Pins, Maskierung), `port/driven/devcontainer_launcher_test.go`. Jede Zeile der Zusage-Tabelle in [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 für den Adapter hat einen Test samt Wächter; die in §3 genannten Test-Namen existieren im Code |
| 3 | Einmal-Messung (Größe, Bauzeit kalt/warm, erster `up`, Trivy) | **erfüllt** | [Messprotokoll](../messungen/slice-v1-devcontainer-start-fundament/README.md) (README, zwei Rohdaten), Skripte unter `tools/messung/devcontainer-start-fundament/`; Linux-/GNU-Voraussetzung im Protokoll und in den Skriptköpfen; `docs/messungen/README.md` führt die Zeile; Suche nach Zugangsdaten-Mustern in Protokoll und Skripten ohne Treffer; `make ablage-pruefung` grün. Nicht gemessen und benannt: `arm64`, kalter Pull der Basis-Images |
| 4 | `make gates` und `make test-docker` grün | **erfüllt** | Abschnitt 1 (lokal `gates`; `test-docker` über den CI-Lauf nach Entscheidung des Projektinhabers) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | beide Reports unter `docs/reviews/`; die Behebung von F-1 bis F-8 und R-1 bis R-6, R-9 steht in §7 und ist im Code nachgelesen (Abschnitt 3). Ein Nachlauf `-r3` über die r2-Behebung ist im Slice als Sache des Reviewers genannt und liegt nicht vor (V-4) |
| 6 | Doku-Update | **erfüllt, Rest bei Closure** | ADR-Index, Handbuch, `CHANGELOG.md`, `spec/` im Diff unberührt wie geplant; Roadmap-Zeile „In Arbeit“ gesetzt (V-5 zur Plan-Tabelle) |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen (Closure)** | §7 „Closure-Notiz“ trägt den Platzhalter; „Finding-Klassen für die Closure“ liegen als Eingabe vor |
| 8 | Beobachtungs-Register fortgeschrieben | **offen (Closure)** | `observations/BEO-APP/podman-static-pin-alterung/` nicht im Diff; §7 „Beobachtungs-Register (für die Closure)“ nennt den Inhalt des dritten Belegs |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **erfüllt** | 11 Risiken, alle mit „Ausgang:“ (entschieden oder weiter offen mit Anker) |
| 10 | Drei Paarungen | **offen (Closure)** | nur bei Closure prüfbar |

## 3. Review-Befunde: Behebung am Ist-Stand (Code gelesen)

| Befund | Ist-Stand | Urteil |
|---|---|---|
| F-1 Wert im Fortschritt/Fehlertext | `runUp`: stderr durch `newMaskingWriter` vor Tail-Puffer und `Progress`; `Flush` nach `cmd.Run`; Parse-Fehler durch `maskString`. `secret_mask.go`: Mindestlänge 4, längster Wert zuerst, Hold-back `maxLen-1` | behoben |
| F-2 Fehlerklasse der Eingabeprüfung | Zeichenmenge einmal im Port (`ValidateDevcontainerProjectFolder`, `…EnvName`), Adapter wendet sie vor jedem Docker-Aufruf an und wrappt mit `ErrDevcontainerRuntime`; Port-Kommentar sagt es | behoben |
| F-3 Ergebniszeile > 1 MiB | stdout ist `tailBuffer` (Ende bleibt), Parser liest das Ende | behoben |
| F-4 Prüfreichweite `/proc` | Aussage im Slice auf das Getestete begrenzt; `docker top`-Scan des Hilfs-Containers ergänzt | behoben |
| F-5 Slice nachgezogen | §3, §6, §7 fortgeschrieben; Plan-Tabelle gegen Diff (Abschnitt 4) | behoben, Reste V-5, V-6 |
| F-6 Host-Voraussetzung der Skripte | DoD-Zeile, Skriptköpfe, Protokoll nennen Linux/GNU | behoben |
| F-7 Aufräumen über den Test hinaus | `Launcher.namePrefix` (Voreinstellung `helperContainerPrefix`), Testpräfix je Test, `scanDuring` mit Kontext | behoben |
| F-8 Fehler von `docker rm -f` | `removeHelper` meldet außer bei „no such container“; `context.Canceled` bleibt erkennbar, kein `ErrDevcontainerRuntime` | behoben |
| R-1 Ergebniszeile durch Maskierung verfälscht | `cmd.Stdout = stdout` ohne Maskierer; Felder sind Daten des Aufrufers (Port-Kommentar zu `DevcontainerUpResult`) | behoben |
| R-2 Flush ohne Wächter | `TestLauncher_Up_ReleasesHeldStderrTail` vorhanden | behoben |
| R-3 Grenze `Config.Env` | Port-Kommentar zu `EnvNames`, §3, §6, Test `…LimitContainerEnvValueIsInDevContainerMetadata` | behoben |
| R-4 Port-Kommentare gegen Code | `ErrDevcontainerRuntime`-Kommentar und Prüfreihenfolge (`Up`: Ordner, Namen, dann Umgebung; `Stop`: ID, dann Umgebung) stimmen mit `Up`/`Stop` überein | behoben |
| R-5, R-6, R-9 | Zustandsmarker (bis auf V-6), Konjunktiv, Hold-back-Verzögerung im `Progress`-Kommentar | behoben (Rest V-6) |
| R-7 Ort der Prüffunktionen | Plan-Eintrag im Spec-Slice (nur Plan, kein Spec-Text); Entscheidung Architect | offen, als solche geführt, kein Blocker |
| R-8 Wortlaut „nie“ gegen `Config.Env` | siehe V-1 | offen |
| R-10 Abtastung des Hilfs-Container-Scans | Wächter verlangt einen Treffer (rot bei Verfehlen, kein falsches Grün); CI-Lauf für `5d63e52` grün, Flackerrisiko trifft Verfügbarkeit, nicht die Zusage | kein Befund; Flackern bei Wiederholung des CI-Laufs beobachten |

## 4. Plan-vs-Code (AGENTS.md §3.10)

Plan-Tabelle (§3) gegen `git diff --stat eea0c92~1 HEAD`:

- Jede Datei des Diffs steht in der Tabelle, **bis auf** `docs/plan/planning/in-progress/roadmap.md` (Zeile „In Arbeit“, Commit `eea0c92`; V-5). Die Tabelle des Slice selbst, die beiden Reports und die Pläne im `open/` stehen drin.
- Jede Tabellenzeile mit Änderungs-Art „neu“/„update“ hat einen Diff; die Zeilen „unverändert“ stimmen (`docker_test.go`, `.golangci.yml`, `port/driven/README.md`, `cmd/uboot/main.go`, `internal/adapter/driving/cli/`, `Dockerfile`, `repo.mk`, `.github/workflows/`, Handbuch, `CHANGELOG.md`, `spec/` sind im Diff nicht vorhanden). Die einzige Zeile ohne Diff ist `observations/BEO-APP/podman-static-pin-alterung/`: sie ist als Closure-Arbeit geplant (DoD 8), kein Abgleichsfehler.
- Stichproben gelesener Aussagen gegen den Code: `requireLocalEngine` ist frei stehende Funktion, `Inspector.requireLocalEngine` ruft sie (Rumpf unverändert, `git diff` auf `inspector.go`); `limitedBuffer` bleibt im Inspector und wird vom Launcher nicht genutzt (`tailBuffer` ist neu); `main.go` konstruiert den Adapter nicht (nicht im Diff); der Launcher wird nur in Tests aufgerufen.
- Eingebettete Dateien: Digest bei beiden `FROM`, exakte Version, `npm ci --ignore-scripts`, Lockdatei mit einer `integrity` wie in §7 Punkt „Pins“ beschrieben.

## 5. Entscheidungs-Konformität

- **[`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 1** (Kontext per `go:embed`, Tar-Strom, Tag aus Version und Hash, Bau nur bei fehlendem Tag): im Code gefunden (`helperContextFS`, `helperImageTag`, `helperContextTar`, `EnsureImage`). **Entscheidung 2** (Digest-Pins, exakte Version, Lockdatei, `npm ci`, `--ignore-scripts`): erfüllt; die Formulierung „`integrity` des ganzen Abhängigkeitsbaums“ ist mit einem Paket erfüllt (die CLI ist gebündelt, §7 benennt es). **Entscheidung 4** (Argumentliste, `--mount`, Spiegelpfad, Namen ohne Wert, SIGINT, Frist, `docker rm -f`): `runArgs` und `runUp` entsprechen; `readonly` ist als Option (Voreinstellung `false`) umgesetzt, die Spannung zu [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) Entscheidung 4 ist vom Projektinhaber entschieden und in §6 festgehalten. **Entscheidung 7** (Port, eigene Vorab-Prüfung zwei plus eine, `ErrDockerUnavailable`/`ErrContainerEngineNotLocal` wiederverwendet, neuer Sentinel): erfüllt; der Eintrag in Mapper und `ExitCode` ist ausdrücklich Sache des Befehl-Slice (`cli.go` unverändert). **Entscheidung 9** (Wächter je Zusage): siehe DoD 2.
- **[`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) Entscheidung 4** (Wächter `initializeCommand`): `…S6_InitializeCommandRunsWithSocket` vorhanden und in CI grün; der Teilsatz des Warntextes ist damit getragen (§7 Punkt 4). **Entscheidung 1** (Exit 11 über `result=failed`) ist Gegenstand des Befehl-Slice; der Port-Kommentar sagt korrekt, dass Fehler nach der Umgebungsprüfung `ErrDevcontainerRuntime` sind.
- **ADR-Immutabilität:** [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) sind im Diff nicht enthalten (Diff-Liste); die Abweichungen von ADR-Text stehen im Slice, nicht in der ADR (§7 Abweichungen 1 und 2, V-1).
- **Lastenheft/Spec:** [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) sagt „Zugangsdaten … stehen nie in einer Datei oder einem Image“; der Code verletzt das für `remoteEnv` nicht, für `containerEnv` steht der Wert in `Config.Env` des Dev-Containers (Metadaten, weder Datei noch Image; V-1). [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): der Socket ist nur im Hilfs-Container (S6-Tests). `spec/` im Diff unverändert.
- **Referenzmodell ([`ADR-0013`](../plan/adr/0013-dokumentationsreferenzmodell.md)):** `docs-check` 0 Befunde; Spec-Straten verweisen nicht nach unten.

## 6. Kommentarregeln (§3.7, §3.9)

Diff `eea0c92~1..HEAD` über `internal/` und `tools/` auf die Muster (a) Chronik/Konjunktiv („former“, „früher“, „originally“, „no longer“, „would“, „wäre“, „instead of“, „previously“), (b) Verweise auf Symbole, die der Diff umbenennt, (c) Helfer zwischen Doc-Kommentar und Funktion durchsucht:

- (a) Kein Treffer in Port, Adapter, `secret_mask.go`, Skripten; ein „instead of“ in `export_test.go` Zeile 44 (`start instead of os/exec`) beschreibt die gegenwärtige Brücke, keine entfernte Fassung (kein Verstoß). `minMaskedValueLen` steht im Indikativ.
- (b) Die Symbole `limitedBuffer`, `validateProjectFolder`, `envNamePattern` kommen nur in ihren eigenen, bestehenden Stellen vor (`limitedBuffer` im Inspector, unverändert); kein Kommentar verweist auf ein umbenanntes Symbol.
- (c) `Inspector.requireLocalEngine` (Methode) steht **nach** der frei stehenden Funktion, deren Doc-Kommentar davor steht; kein Helfer trennt Kommentar und Funktion. Die Methode trägt keinen Doc-Kommentar (zulässig).
- Die Kommentare der Testdateien führen ADR-Kennungen als Rang-Zeiger (zulässig). Prosa-Verstoß gegen §3.7 siehe V-6 (Zustandsmarker in der Implementer-Übergabe).

## 7. Findings

| ID | Schwere | Befund | Fundstelle | Auflage |
|---|---|---|---|---|
| V-1 | MEDIUM | **Die Zusage „der Wert steht nie in einer Kommandozeile“ trägt für Prozesse im Hilfs-Container nicht; die förmliche Entscheidung steht aus.** Gemessen und von Tests gepinnt: die Dev-Container-CLI schreibt den Wert für die Dauer eines Lebenszyklus-Befehls als `docker exec -e NAME=Wert` (`remoteEnv`) bzw. `docker run -e NAME=Wert` (`containerEnv`) in die Kommandozeile eines Prozesses *im Hilfs-Container* (`…S6_LimitRemoteEnvValueIsOnHelperCommandLine`); bei `containerEnv` steht er dazu im `Config.Env` des Dev-Containers (`…LimitContainerEnvValueIsInDevContainerMetadata`), und der laufende Hilfs-Container trägt ihn in `Config.Env` (§7 Abweichung 1). [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 sagt „Der Wert steht nie in der Kommandozeile des Prozesses“ und Entscheidung 9 prüft das nur für die Argumentliste des Launchers; die Slice-Rückführung „ein Docker-Test widerlegt eine Zeile der Entscheidung → Architect, Folge-ADR“ (§4) wurde nicht förmlich ausgeführt, sondern im Slice als Grenze geführt („Entscheidung des Projektinhabers über den Wortlaut steht aus“, §6). Der Wortlaut in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („nie in einer Datei oder einem Image“) wird nicht verletzt; die Allaussage „nie auf einer Kommandozeile“ steht dort nicht, wohl aber in der Linie von [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und dem Plan des Befehl-Slice (Handbuch). Die Grenze ist im Port-Kommentar, im Slice und mit zwei Tests offen benannt; keine Harness-Lüge, aber eine offene Entscheidung | Slice §3 „S6 Umgebung“, §6 letzter Punkt; `devcontainer_launcher.go` (Port) `EnvNames`; [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 und 9 | Architect/Projektinhaber vor der Closure: den Wortlaut entscheiden (Handbuch und Befehl-Slice) und festhalten, ob eine Ergänzung der ADR (Folge-ADR) nötig ist, die die Zusage auf „Prozesse außerhalb des Hilfs-Containers und Launcher-Argumentliste“ begrenzt; im Slice §6 klarstellen, *welcher* Wortlaut aussteht (das Lastenheft ist 0.18.0 bereits geschrieben) |
| V-2 | LOW | **Die Digests der Basis-Images wurden von mir nicht neu gegen die Registry gelesen** (kein Netz-Auftrag). Der Beleg ist der grüne Bau in CI (S5 baut aus dem Kontext nach `docker rmi`) und die Angabe der Plattformen im Slice §7 (beide Indizes listen `amd64` und `arm64/v8`). Die `arm64`-Aussage bleibt unbelegt (§6 „Digest und Plattform“, weiter offen, korrekt geführt) | Slice §7 Pins, §6 | keine; bei der Closure „weiter offen“ für `arm64` belassen |
| V-3 | LOW | **Closure-Trigger §5 verlangt „bis dahin zeigt der Nachweis ein lokaler `make test-docker`-Lauf mit Ausgabe“; lokal ist der Lauf rot (Host-Netz), der Nachweis ist der CI-Lauf ohne `-v`.** Die Abweichung ist vom Projektinhaber entschieden und in §7 genannt, aber der Wortlaut in §5 nicht angepasst. Ein Skip ist im CI-Log nicht pro Test sichtbar; getragen wird es durch das Fehlen jedes `t.Skip`-Pfads (geprüft) und die Paketlaufzeit | Slice §5; §7 Sensor-Tabelle | Closure-Notiz nennt den CI-Lauf `39409233721` und die Begründung (kein `t.Skip`-Pfad, `-v` kommt mit dem Befehl-Slice); §5 bleibt als Trigger lesbar |
| V-4 | LOW | **Die Behebung von r2 (R-1 bis R-6, R-9) ist nicht im frischen Kontext des Reviewers geprüft** (der Slice nennt `-r3`; die Datei liegt nicht vor). Der Verifier hat die Behebung im Code nachgelesen (Abschnitt 3) und eine Stichprobe der Wächter mutiert, ersetzt aber kein Review | Slice §7 „r2“ | Planner entscheidet: `-r3` nachholen oder in der Closure-Notiz benennen, dass die Verifier-Nachlesung an seine Stelle trat |
| V-5 | LOW | **Plan-Tabelle (§3) führt `roadmap.md` nicht, obwohl der Slice sie ändert** (Commit `eea0c92`: „Nichts in Arbeit.“ → „In Arbeit: …“). Die DoD-Zeile 6 sagt „Roadmap-Zeile ‚In Arbeit‘ bei der Closure“; die Zeile steht bereits seit dem Start. Abgleich nach der letzten Fix-Runde (AGENTS.md §3.10) hat sie nicht erfasst | Slice §3 Tabelle, §2 DoD 6; `docs/plan/planning/in-progress/roadmap.md` | Zeile `roadmap.md` in die Tabelle aufnehmen (Roadmap bei Closure auf „nichts in Arbeit“ bzw. Folge-Slice setzen), DoD-Text auf den Ist-Stand |
| V-6 | LOW | **Abgelaufene Zustandsmarker in der Implementer-Übergabe (§3.7, §3.10):** (a) „Ein voller `make test-docker`-Lauf nach dem Zusatz steht aus.“ — der CI-Lauf für `5d63e52` hat ihn geliefert; (b) „**Nicht gelaufen / nicht belegt.** … `make test-docker` grün (siehe oben); Review, Verifier.“ — Review ist gelaufen, die Zeile ist nach dem CI-Nachweis nicht mehr zutreffend; (c) die Sensor-Tabelle der Erst-Übergabe führt `make test-docker` „Exit 2“ ohne Hinweis auf den späteren CI-Beleg (die Nachträge in den Folge-Tabellen stehen) | Slice §7 „Implementer-Übergabe“ (Absatz „Reihenfolge“, Absatz „Nicht gelaufen / nicht belegt“) | Planner bei der Closure auf den Ist-Stand setzen oder entfernen |
| V-7 | INFO | **Port-Paket `port/driven` trägt erstmals Funktionen** (R-7 des Reviews); keine Regel in `.golangci.yml` verbietet das, `depguard` ist unberührt, Präzedenzfall `NewConfigHintError`. Der Spec-Slice trägt den Plan-Eintrag | Port; [`slice-v2-spec-devcontainer-start`](../plan/planning/open/slice-v2-spec-devcontainer-start.md) | Architect entscheidet den Ort (Sicht §2.4 oder `hexagon/domain`), kein Blocker |
| V-8 | INFO | **Bekannte Grenzen, korrekt als „weiter offen“ geführt:** Abbruch dauert bis zur Frist (~10,9 s, `--init` nicht untersucht), Hilfs-Image 1716 MiB mit 597 HIGH/CRITICAL im Debian-Anteil, Maskierung deckt nicht jede Form, `arm64` unbelegt, Laufzeit von `make test-docker` knapp am Limit (e2e-Paket 593,8 s in CI, Gesamtlauf liegt in CI bei ~11 Minuten gegen 30 min Limit) | Slice §6 | keine; Ausgänge bei der Closure übernehmen |

## 8. Traceability, Carveouts, Nicht ausgeführt

Traceability:

| ID | Beleg |
|---|---|
| [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidungen 1, 2, 4, 5, 7, 9 | Adapter und Port im Code gelesen; Docker-Tests S2, S5, S6, S7 samt Wächtern grün in CI (`integration-docker` für `5d63e52`); Unit-Tests in `make gates` |
| [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md) Entscheidung 4 | `…S6_InitializeCommandRunsWithSocket` (CI grün) |
| [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) | `TestLauncher_Up_PassesEnvByNameOnly` (von mir mutiert: rot), `…S6_ValueIsNeverInImageLabelFileOrCommandLine`; Grenze: V-1 |
| [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | `…S6_SocketNotInStartedContainer` mit Wächter |
| [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) | Pins in den eingebetteten Dateien (`TestHelperContext_Pins`), S5 mit verfälschter `integrity` |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Sentinel im Port, Preflight-Tests (11); Eintrag in beide Klassifikatoren im Befehl-Slice |

Carveouts: Neu: none. Gelöst: none. Unverändert: Inventar nicht berührt (`docs/plan/planning/in-progress/carveouts.md` nicht im Diff). Gates nicht gelockert (`.golangci.yml`, `repo.mk`, `Dockerfile`, `.github/workflows/` unverändert).

Nicht ausgeführt:

- `make test-docker` lokal — Host-Netz (`download.docker.com`); Nachweis ist der grüne CI-Lauf für `5d63e52`.
- `make verify-depguard` — keine Importregel berührt.
- `make ci` / `make fullbuild` lokal — kein Security-/Release-Pfad berührt; `ci` in GitHub grün.
- Neu-Lesen der Registry-Digests und `arm64`-Bau — außerhalb des Verifier-Auftrags, `arm64` im Slice als offen geführt.

Commit / Artefakt: `HEAD` = `origin/main` = `5d63e52`; dieser Report ist uncommitted (kein Commit durch den Verifier).
