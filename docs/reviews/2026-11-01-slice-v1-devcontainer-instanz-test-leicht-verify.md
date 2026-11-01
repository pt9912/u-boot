# Review-Report: slice-v1-devcontainer-instanz-test-leicht (Verifier) — 2026-11-01

**Review-Art:** Code — Verifikation gegen Plan, DoD, Entscheidungen und Sensors (Modul 11: „Bauen wir es richtig?“). Kein Review des Diffs und keine Validierung.

**Gegenstand:** Slice [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-test-leicht.md) im Hauptbaum (Weg A, Entscheidungen des Projektinhabers 2026-11-01), HEAD `172387b5` (Fix-Runde Review `32cecdf6` enthalten). Der Slice liegt in `in-progress/`. Diff für diesen Slice: `73b42e21`, `54c001db`, `32cecdf6` samt den Move-Commits seit Basis `c370280c`; der Diff dieser Spanne trägt auch Dateien des geschlossenen Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md), die hier nicht bewertet werden.

**Skill:** `.harness/skills/verifier.md` existiert nicht; die Rolle folgt `harness/roles.md` und `harness/verification.md`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-test-leicht.md) (§1 entschiedene Fragen, §2 DoD, §3 Plan, §4 Trigger, §5 Closure-Trigger, §6 Risiken); Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/open/slice-cr-devcontainer-instanz-zusage-getestet.md)
- [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz Instanzen)
- [Review-Report](2026-11-01-slice-v1-devcontainer-instanz-test-leicht.md) (F-1 bis F-8)
- `harness/verification.md`, `AGENTS.md` (Hard Rules)

**Sensors dieses Laufs:**

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` (Docker-Modus, HEAD `172387b5`) | pass | Exit 0; `baseline-verify` v6.18.1 OK (54 Dateien); `docs-check` 971 Dateien, 0 Befunde; Lint „0 issues“; Coverage 92,20 % gegen Schwelle 90 %; `ablage-pruefung` 26 Fälle, 0 Fehlschläge; Stempel `.harness/state/gates-passed.head` zeigt `172387b5`. Während des Laufs lag eine fremde ungetrackte Review-Datei im Baum (`docs/reviews/…-nachlauf-ruecknahme.md`); `docs-check` meldete nichts, eine Wiederholung war nicht nötig. |
| `make docs-check` mit diesem Report | im Abschlussbericht des Verifiers genannt (nach dem Schreiben gefahren) | — |
| Übersetzbarkeit Tag `docker` | **nicht gebaut**, aus dem Quelltext gelesen (unten) | `make lint` baut den Tag nicht; kein `make test-docker` (Entscheidung des Projektinhabers) |
| `make test-docker` | **nicht gefahren** | Entscheidung des Projektinhabers; der Lauf `integration-docker` auf `main` entscheidet (Anker [`slice-v1-test-docker-in-ci`](../plan/planning/open/slice-v1-test-docker-in-ci.md)) |
| `make verify-depguard` | nicht gefahren | keine Importregel berührt (kein Diff an `.golangci.yml`, Importe der Test-Datei nur Standardbibliothek und `driving`) |

---

## Übersetzbarkeit des Tests, erneut gelesen

`internal/e2e/devcontainer_instance_docker_test.go` (Tag `docker`, Paket `e2e_test`). Jede Referenz mit Fundstelle im Paket:

| Referenz | Fundstelle | Befund |
|---|---|---|
| `requireDocker` | `sandbox_status_docker_test.go:41` | ok |
| `dcGenerateProject(ctx, t, sandbox bool, map[string]string) string` | `devcontainer_start_docker_test.go:191` | ok, `nil` zulässig |
| `dcService() *DevcontainerStartService` | `devcontainer_start_docker_test.go:281`; `Up`/`Down` mit Signatur des Ports `driving.DevcontainerUseCase` (`application/devcontainer_start.go:193, 462`) | ok |
| `ubootBinary` | `devcontainer_start_docker_test.go:710` | ok |
| `dcInspect`, `dcInspectDoc.State.Running`, `.Config.Labels`, `.Mounts`, `volumeAt` | `devcontainer_start_docker_test.go:75, 92, 138` | ok |
| `dcFolderLbl`, `dcConfigLbl`, `dcTestBudget` | `devcontainer_start_docker_test.go:47–49` | ok |
| `dcExec` (zwei Rückgaben), `dcExists`, `dcMustDocker`, `dcDocker` | `devcontainer_start_docker_test.go:315, 102, 65, 54` | ok |
| `runBinary` (drei Rückgaben) | `devcontainer_shell_docker_test.go:222` | ok |
| `dcWatchHelpers` (`stop func() []string`), `dcHelperNames` | `devcontainer_shell_docker_test.go:98, 68` | ok, `started := stop()` passt |
| `dcProjectDir` | `devcontainer_start_docker_test.go:154` | ok |
| `driving.DevcontainerUpRequest{BaseDir, Instance, Recreate}`, `DevcontainerDownRequest{BaseDir, Instance}`, Felder `ContainerID`, `AlreadyRunning`, `Stopped`, `RemoteWorkspaceFolder`, `ErrDevcontainerOnlyInstances` | `hexagon/port/driving/devcontainer.go` (Zeilen 20–24, 45–60, 70, 105–118, 140) | ok |

Importe `context`, `errors` (Zeile 277), `os` (342), `path/filepath` (342), `strings`, `testing`, `time` (290, 319), `driving` sind benutzt; keine weiteren nötig. Variablen: `names`, `ids`, `markers`, `workspaces`, `existing`, `again`, `repeat`, `down`, `none`, `stopped`, `restarted`, `recreated`, `doc`, `baseVolume`, `before`, `after`, `bin` werden alle gelesen; `down, err :=` führt `err` ein, `none, err :=` und `stopped, err :=` je ein neues Symbol (kein Redeclare); `doc` (Zeile 140, Schleifen-Scope) und `doc :=` (236) liegen in verschiedenen Scopes. Doppelte Deklarationen der Helfer `dcRunning`, `dcMarkerPath`, `dcFolderContainers`, `dcWatchNoHelper`, `dcRequireDistinct`, `dcUpInstance`, `dcInstanceLabel`, `dcWorkspaceDest` im Paket: je genau eine (Suche über `internal/e2e/*.go`).

**`devcontainer_start_docker_test.go`, `HostConfig`:** die Struktur trägt nur noch `Binds` und `Mounts` (Zeilen 81–84). Leser über `internal/`: `hasSocket` (Zeilen 113, 118) im selben File; `envHolds(doc dcInspectDoc, …)` in `devcontainer_shell_docker_test.go:300` liest `Config.Env`, nicht `HostConfig`; `ollama_host_docker_test.go` liest `HostConfig.ExtraHosts` per `docker inspect --format`, nicht über die Struktur. Suche nach `HostConfig.Memory|MemorySwap|NanoCpus` und `.NanoCpus|.MemorySwap`: kein Treffer. Entfernte Symbole (`dcInstanceCount`, `dcInstanceMin*`, `dcInstanceMemory`, `dcInstanceNanoCPU`, `dcDockerDest`, `dcOllamaDest`, `dcMeminfoBytes`, `dcRequireInstanceHost`, `dcLimitRunArgs`, `dcMarkerPaths`, `…FourInstancesAreIsolated`): außerhalb von `.claude/worktrees` in keiner Go-Datei; Treffer nur in Chronik und im Slice selbst.

**Ergebnis:** Ich finde keinen Übersetzungsfehler. Das ist Lesen, nicht Bauen (V-3).

---

## DoD-Abgleich (Slice §2, Punkt für Punkt)

| # | DoD-Punkt | Stand | Beleg / Grund |
|---|---|---|---|
| 1 | Leichter Test liegt in der Datei, keine Vorprüfung, keine Grenzen, kein `nestedRuntime`/`ollama` (Weg A), **läuft grün** auf dem Runner | **Code erfüllt, Lauf offen** | Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (Zeile 115); `dcGenerateProject(ctx, t, true, nil)` (119), keine Einstellung, kein `runArgs`, kein Meminfo-Zugriff. Der grüne Lauf `integration-docker` liegt nicht vor (nicht gelaufen). Pflicht des Punktes nur zur Hälfte nachweisbar. |
| 2 | Schwerer Test ersatzlos entfernt; Konstanten und Helfer entfernt; `make lint` ohne ungenutzte Konstante | **erfüllt** | Suche nach allen genannten Namen und `FourInstances` im Baum (siehe oben); Lint „0 issues“ (der Lint baut den Tag `docker` nicht, daher belegt er nur die Nicht-`docker`-Dateien; die Abwesenheit von Resten in der Test-Datei ist gelesen). Kein Skript unter `tools/messung/` im Diff. |
| 3(a) | Mutationsprobe „im Bericht des Implementers“, jede Mutation färbt rot | **offen** | Es liegt kein Mutationsprotokoll des Implementers vor, und der Test ist nicht gelaufen; die vier Mutationen sind nur aus dem Quelltext bewertet (unten). Nach der Plan-Zusage (§6 Risiko „Wächter für `--recreate`“) bleibt die Probe „weiter offen bis zur Mutationsprobe des Verifiers“; ich kann sie nicht fahren. |
| 3(b) | `zusagen-und-messstaende.md` (Fundstelle Zeile 178) nennt den leichten Test | **erfüllt** | Zeilen 177–183: `…TwoInstancesAreIsolated`, „ohne `nestedRuntime`, ohne `ollama`, ohne Grenze“; „Volumes für Docker-Speicher und Modelle, die Grenze von einer CPU und 3 GiB und vier Instanzen läuft die Suite nicht“. Kein Messskript mehr genannt (F-2 des Reviews behoben). |
| 3(c) | Kopfkommentar (Zeilen 3–14) und Docstring beschreiben den Ist-Stand (§3.7) | **erfüllt** | Zeilen 3–11 und 96–114 gelesen: Indikativ über den Zustand; die Verneinungen („no nested runtime, no feature `ollama`, no resource limit“) sind Abgrenzung; keine Chronikwörter (Suche „former/früher/no longer/gone/entfernt“ in der Datei: kein Treffer). |
| 3(d) | Slice-/Register-Verweise auf den schweren Test nachgezogen | **erfüllt** | Verweise außerhalb der Chronik: nur im Slice selbst (Beschreibung des Gegenstands) und im Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/open/slice-cr-devcontainer-instanz-zusage-getestet.md), dort auf „Docker-Test mit vier Instanzen, ersatzlos entfernt“ gesetzt (Review F-1); Zeilen mit „drei Instanzen“ im Entwurf: keine mehr; „Messskript“ nur noch verneinend (Zeilen 15, 50). Chronik unverändert: abgeschlossener Slice [`slice-v1-devcontainer-instanz-label`](../plan/planning/done/slice-v1-devcontainer-instanz-label.md), Reports, Evidence unter `observations/BEO-APP/…` (Plan §3, Zeile „unverändert (Chronik)“). |
| 4 | `make gates` grün | **erfüllt** | selbst gefahren, siehe Sensors |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | [Review-Report](2026-11-01-slice-v1-devcontainer-instanz-test-leicht.md); F-1 bis F-5 in `32cecdf6` behoben (von mir gegengelesen, siehe „Review-Findings“); F-6 bis F-8 sind INFO und tragen Risiken in §6 |
| 6 | Doku-Update: `zusagen-und-messstaende.md`; `CHANGELOG.md` Zeile 16 auf den Ist-Stand; `docs/user/` unverändert | **erfüllt** | `CHANGELOG.md` Zeile 16 (`[Unreleased]`): „Docker tests with two instances and the container without an instance in the sandbox profile (default settings: no nested runtime, no feature `ollama`, no resource limit) and with three containers in the standard profile“. `docs/user/benutzerhandbuch.md`: der Diff der Spanne berührt nur Zeile 864–865 (Abschnitt `devcontainer shell`, Gegenstand des anderen Slice), nicht den Instanz-Abschnitt (Zeilen 896 ff.); keine Zahl der Instanzen, kein Test genannt |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen** | §7 lautet „Wird bei der Closure gefüllt.“ (Slice liegt in `in-progress/`) |
| 8 | Beobachtungs-Register fortgeschrieben | **offen** | Kein Diff unter `observations/BEO-E2E/…` für diesen Slice (Plan §3 führt die Zeile als Closure-Arbeit) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **offen** | Fünf Risiken tragen „weiter offen“ oder „offen bis zum ersten Lauf“; zwei sind „entfallen“ (Messweg, Testname in Chronik), eines „entfallen, wenn (c) erfüllt“ (erfüllt, noch nicht gesetzt) |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen** | hängt an 7 bis 9 |

Alle Kästchen in §2 stehen unmarkiert; für einen Slice in `in-progress/` ist das stimmig.

## Entscheidungs-Konformität (Projektinhaber 2026-11-01, ADR, Lastenheft)

| Entscheidung | Befund |
|---|---|
| Schwerer Vier-Instanzen-Test ersatzlos weg, kein Messskript | konform: Test, Konstanten, Helfer entfernt; kein Diff unter `tools/messung/`; Messprotokoll [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md) ohne Diff |
| Weg A: Sandbox-Profil mit Standard-Einstellungen | konform: `dcGenerateProject(ctx, t, true, nil)`; kein `nestedRuntime`, `ollama`, keine Grenzen, keine Speicher-Vorprüfung. Das Sandbox-Profil ohne diese trägt genau ein Volume `…-workspace-${devcontainerId}` am `workspaceFolder` (`devcontainer_sandbox.go` Zeilen 103–104, 139–141) = `dcWorkspaceDest` `/workspaces/dcstart` (gleicher Wert nutzt `…S1_S4_SandboxProfileVolumes`, `devcontainer_start_docker_test.go:516`) |
| Zwei Instanzen plus Container ohne Instanz | konform: `names := []string{"i1","i2"}`, Container ohne Instanz zuerst (Zeile 124) |
| Läuft auf dem GitHub-Runner | **nicht belegt**: keine Hostressourcen-Prüfung, drei Container ohne Grenze; `.github/workflows/integration.yml` ruft `make test-docker` ohne Filter (gelesen im Plan; unverändert im Diff). Ob die Annahme trägt, zeigt erst der Lauf (Risiken 3 und 4 des Slice) |
| [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) (zweites Id-Label) | konform: Labels `devcontainer.local_folder` und `u-boot.instance` an Instanzen, kein `devcontainer.config_file`; am Container ohne Instanz Config-Label, kein Instanz-Label (Zeilen 144–160); Test hat keine Produktänderung (Diff ohne Produktdatei dieses Slice) |
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) unverändert | konform: Lastenheft und Spec im Diff nur durch den anderen Slice berührt (`spezifikation.md`, `architecture.md`); das Lastenheft ist unberührt. Die Diskrepanz Zusage („bis zu vier“) gegen Suite trägt Risiko 1 → Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/open/slice-cr-devcontainer-instanz-zusage-getestet.md) |
| Keine Gate-Senkung (`AGENTS.md` §3.6) | konform: Schwelle, Lint, Importregeln unberührt; der Test wird ersetzt, nicht das Gate gesenkt. Zu beachten: die Suite belegt danach weniger (Docker-Speicher, Modelle, Grenze, drei und vier Instanzen); das ist in Slice §1, `zusagen-und-messstaende.md` und Risiko 1/6 offen benannt |

## Plan-vs-Code-Diff (Plan §3 gegen `git diff --stat`, nach der Fix-Runde)

Diff für diesen Slice (ohne die Dateien des anderen Slice):

| Datei im Diff | Plan §3 | Befund |
|---|---|---|
| `internal/e2e/devcontainer_instance_docker_test.go` | update | deckt sich (−235/+, schwerer Test entfernt, leichter ergänzt) |
| `internal/e2e/devcontainer_start_docker_test.go` | update (`HostConfig`) | deckt sich (−4 Zeilen: drei Felder und Kommentar) |
| `docs/maintainer/zusagen-und-messstaende.md` | update | deckt sich (Fundstelle Zeilen 177–183); der Diff trägt zusätzlich Änderungen des anderen Slice in Zeile ~1053 |
| `CHANGELOG.md` | update | deckt sich (Zeile 16); die Zeile `### Changed` zum `devcontainer shell` ist vom anderen Slice entfernt |
| `docs/plan/planning/open/README.md` | update | deckt sich (Zeile des Slice entfernt) |
| `docs/plan/planning/in-progress/roadmap.md` | update | deckt sich (Zeiger „Leichter Isolationstest … (in Arbeit)“; F-3 des Reviews: „anderer Lauf als der vorige“ ist entfernt) |
| [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/open/slice-cr-devcontainer-instanz-zusage-getestet.md) | update | deckt sich (Plan-Zeile ergänzt, F-2) |
| [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-test-leicht.md) | (der Slice selbst) | Änderungen der Fix-Runde |
| `docs/plan/planning/observations/BEO-DOK/{plan-aussage-ueber-code,spec-bedingung-ohne-code-entsprechung,zusage-weiter-als-messung}/state.md` | nicht geführt | Commits `1fe883f7` und `172387b5` tragen den Slice [`slice-v1-shell-erkennung-fehlender-benutzer-eng`](../plan/planning/done/slice-v1-shell-erkennung-fehlender-benutzer-eng.md), nicht diesen Slice; kein Befund |
| [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md) | unverändert | im Diff nicht enthalten, deckt sich |
| `spec/lastenheft.md`, `docs/user/` (Instanz-Abschnitt), `README*.md`, `.github/workflows/integration.yml`, `Makefile` | unverändert | im Diff nicht enthalten (für `docs/user/` siehe DoD 6), deckt sich |
| `docs/plan/planning/observations/BEO-E2E/…` | update, neu (Closure) | im Diff nicht enthalten; zulässig, die Closure liegt noch aus |

Eine Datei im Diff ohne Plan-Zeile: keine (für diesen Slice). Eine Plan-Zeile ohne Diff: nur die Closure-Zeile `BEO-E2E`.

## Review-Findings gegengelesen (F-1 bis F-8)

| Finding | Stand |
|---|---|
| F-1 (MEDIUM, Entwurf trägt drei Instanzen und Messskript) | behoben: Entwurf auf zwei Instanzen, kein Messskript, Variante 1 „zwei Instanzen“, Risiko neu gefasst (`Die Suite belegt „zwei“ erst nach einem grünen Lauf`) |
| F-2 (LOW, Plan-Tabelle/DoD (b) „Messskript“) | behoben: Plan-Zeile für den Entwurf ergänzt, DoD (b) nennt nur den leichten Test |
| F-3 (LOW, Roadmap „anderer Lauf als der vorige“, WIP) | behoben: Zeiger „(in Arbeit)“; §4 Start nennt „WIP-Limit 1 pro Lauf … Anweisung des Projektinhabers, 2026-11-01“; `Verantwortlich:` nennt den Lauf. Der andere Slice liegt inzwischen in `done/` |
| F-4 (LOW, `HostConfig`-Felder ohne Leser) | behoben (siehe oben) |
| F-5 (LOW, Zusage weiter als die Suite) | getragen: Risiko 1 und 6 → Entwurf; nicht geschlossen, so geplant |
| F-6 (INFO, Test nicht gelaufen, zwei Annahmen) | als Risiko 3/4 in §6 geführt, offen |
| F-7 (INFO, Runner) | als Risiko 3 geführt, offen |
| F-8 (INFO, `--recreate`-Wächter) | siehe Mutationen, V-2 |

## Mutationen, aus dem Quelltext bewertet (nicht gefahren)

Der Test ist ein Docker-Test; ich fahre sie nicht. Bewertung durch Lesen der Prüfungen, nicht durch Lauf:

| Mutation | Prüfung, die anschlägt | Bewertung |
|---|---|---|
| (1) `up --instance` ohne Instanz-Label oder mit dem Label einer anderen Instanz | `got != names[n-1]` (Zeile 154); bei fehlendem Label zusätzlich geteilte Volume-Namen (`dcRequireDistinct "workspace volume"`) | trägt (Fehlerzweig mit `t.Errorf`, der Test läuft weiter und färbt rot) |
| (2) Workspace-Volume gemeinsam (`-${devcontainerId}` entfällt) | `dcRequireDistinct(…"workspace volume"…)` (163) und die Paar-Probe `(code == 0) != (i == j)` (175); jede Marker-Datei ist dann in jedem Container sichtbar | trägt doppelt; die Paar-Probe fängt es auch, wenn die Namen verschieden, der Inhalt aber geteilt wäre |
| (3) `up` ohne `--instance` wählt die neueste Instanz | `again.ContainerID != existing.ContainerID || !again.AlreadyRunning` (184); `shell` ohne `--instance`: `hostname != existing.ContainerID[:12]` (187); Exit-10-Zweig (Container ohne Instanz entfernt) | trägt. Annahme der Shell-Prüfung: der `hostname` im Container ist die gekürzte Container-ID (Docker-Standard ohne `--hostname`; so liest es auch der Code der Shell-Tests); durch den Lauf zu belegen |
| (4) `--recreate` entfernt das Volume | allein die Marker-Prüfung `test -e …marker-i1` nach `--recreate` (240). `doc.volumeAt(…) != workspaces[1]` (237) fängt es nicht, weil der Volume-Name aus den Labels folgt und gleich bleibt (wie Review F-8) | trägt, aber nur mit einer Prüfung. Sie trägt, solange die Datei im Volume liegt und das neue Volume leer wäre; das Verhalten bei einem Volume mit gleichem Namen, das erhalten bleibt, ist der Positivfall |
| (5, zusätzlich gelesen) `down` ohne `--instance` stoppt auch Instanzen | `dcRunning` je Container nach `Down` (207–211) | trägt |
| (6, zusätzlich gelesen) `down --instance nope` meldet Fehler/Stop | Zeilen 213–216 | trägt |
| (7, zusätzlich gelesen) `--recreate` ersetzt auch die andere Instanz | `dcExists(ids[2]) && dcRunning(ids[2])` (243) | trägt |

Gerade (4) und die Hostname-Annahme von (3) sind erst durch einen Lauf belegt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | MEDIUM | DoD Punkt 1 („läuft im Lauf `integration-docker` auf dem GitHub-Runner grün“) und der Closure-Trigger (§5) sind nicht erfüllt: Der neue Test ist weder lokal noch in der CI gelaufen. Zwei Annahmen (Standard-Sandbox startet ohne `egress.enabled`/`onUnavailable`; Marker per `exec -u root` les- und schreibbar) und die Annahme „drei Sandbox-Container tragen auf dem Runner“ stehen im Plan als offene Risiken 3 und 4. | Slice §2 Punkt 1, §5, §6 | Slice · „läuft im Lauf `integration-docker` auf dem GitHub-Runner grün“ | ja — Lauf `integration-docker` auf `main` | Test ohne Lauf (Closure-Bedingung offen) |
| V-2 | LOW | DoD Punkt 3(a) verlangt eine Mutationsprobe im Bericht des Implementers; es liegt keines vor, und der Verifier kann sie nicht fahren. Die vier Mutationen sind aus dem Quelltext als wirksam bewertet; die Mutation (4) hängt an einer einzigen Prüfung (`test -e` auf die Marker-Datei). | Slice §2 Punkt 3(a), §6 Risiko „Wächter für `--recreate`“ | Slice · „Mutationsprobe am leichten Test, im Bericht des Implementers“ | ja — Lauf mit den vier Mutationen auf einem Docker-Host | Wächter ohne Mutationslauf |
| V-3 | LOW | Die Übersetzbarkeit des Tags `docker` ist nur gelesen, nicht gebaut; `make lint` und `make test` bauen den Tag nicht, ein Fehler fiele erst im CI-Lauf `integration-docker` auf. Ich habe keinen Fehler gefunden. | `harness/verification.md` §Sensor-Auswahl (Docker-Tests: `make test-docker`) | `internal/e2e/devcontainer_instance_docker_test.go` · Build-Tag `docker` | ja — `go vet -tags docker ./internal/e2e/` oder `make test-docker` | Test-Tag nicht von einem lokalen Gate gebaut |
| V-4 | INFO | Closure-Pflichten sind erwartungsgemäß offen (DoD 7 bis 10): Closure-Notiz, Register `BEO-E2E`, Risiko-Ausgänge, Paarungen. Der Slice liegt noch in `in-progress/`. | Slice §2, §7 | Slice · „Wird bei der Closure gefüllt.“ | ja — Prüfung bei der Closure | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Quelltext des neuen Tests (Referenzen, Importe, Variablen, Deklarationen) | geprüft, ohne Befund (gelesen, nicht gebaut) |
| `devcontainer_start_docker_test.go`, `dcInspectDoc.HostConfig`: Leser in anderen `docker`-getaggten Dateien | geprüft, ohne Befund (kein Leser von `Memory`, `MemorySwap`, `NanoCpus`) |
| Entfernte Konstanten und Helfer | geprüft, ohne Befund (keine Referenz im Baum außer Chronik) |
| Kommentare der Test-Datei (`AGENTS.md` §3.7, §3.9) | geprüft, ohne Befund |
| `CHANGELOG.md` Zeile 16, `zusagen-und-messstaende.md` Fundstelle | geprüft, ohne Befund |
| `docs/user/`, `spec/`, `README*.md`, Messprotokoll, Workflow, `Makefile` | geprüft, für diesen Slice unverändert |
| Plan §3 gegen `git diff --stat` nach der Fix-Runde | geprüft, ohne Befund |
| Lifecycle-Bewegungen des Slice (je reiner Move, Nachzug-Commit danach) | geprüft, ohne Befund (`AGENTS.md` §3.3) |
| Carveouts | unverändert: keiner neu, keiner gelöst |
| `make test-docker`, Lauf des neuen Tests | nicht geprüft: Entscheidung des Projektinhabers (V-1) |
| `make verify-depguard` | nicht geprüft: keine Importregel berührt |

## Verification Evidence (Schema `harness/verification.md`)

Scope:
- Slice: [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/in-progress/slice-v1-devcontainer-instanz-test-leicht.md)
- IDs: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md)
- Artefakte: `internal/e2e/devcontainer_instance_docker_test.go`, `internal/e2e/devcontainer_start_docker_test.go`, `docs/maintainer/zusagen-und-messstaende.md`, `CHANGELOG.md`, Entwurf zur Zusage, Roadmap-Zeiger, `open/README.md`

Traceability:

| ID | Beleg |
|---|---|
| [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Docker-Test `TestDevcontainerInstance_Docker_TwoInstancesAreIsolated` (Container- und Label-Identität, Workspace-Volume je Container, Marker-Isolation über alle geordneten Paare, Container-Wahl ohne `--instance`, `down`, `--recreate`, Exit 10 ohne Container ohne Instanz); **nicht gelaufen**; Standardprofil `…StandardProfileSharesTheProjectFolder` unverändert |
| [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md) | dieselbe Testlinie: Labels, `--instance`-Wahl; Unit-Tests des Produkts unverändert grün in `make gates` |

Carveouts: Neu: keiner · Gelöst: keiner · Unverändert: alle.

Nicht ausgeführt: `make test-docker` (Entscheidung des Projektinhabers; CI-Lauf auf `main` entscheidet, Anker [`slice-v1-test-docker-in-ci`](../plan/planning/open/slice-v1-test-docker-in-ci.md)); Docker-Test-Mutationen; Bau mit Tag `docker`; `make verify-depguard` (keine Importregel berührt).

Commit / Artefakt: HEAD `172387b5`; Fix-Runde `32cecdf6`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Test ohne Lauf (Closure-Bedingung offen) · Wächter ohne Mutationslauf · Test-Tag nicht von einem lokalen Gate gebaut

## Verdikt

**Merge-blockierend:** nein für den Stand der Umsetzung; der Slice ist **nicht schließbar**, solange V-1 offen ist. Code, Doku-Nachzüge und `make gates` stimmen mit Plan und DoD überein; die Fix-Runde hat F-1 bis F-5 des Reviews behoben. Offen sind der Lauf des Tests (DoD Punkt 1, Closure-Trigger §5), das Mutationsprotokoll (DoD 3(a)) und die Closure-Pflichten (DoD 7 bis 10).

**Bewertung zur Frage „mit ‚Test nicht gelaufen‘ schließen?“ (Entscheidung beim Planner/Projektinhaber, nicht hier):** Der Closure-Trigger §5 verlangt ausdrücklich einen grünen Lauf `integration-docker` oder einen als Carveout dokumentierten roten Lauf; „nicht gelaufen“ ist keiner von beiden. Ein Schluss ohne Lauf braucht daher eine Änderung des Triggers oder die Vorgabe, auf den CI-Lauf auf `main` zu warten; die Risiken 3 und 4 (Runner trägt drei Container, zwei Annahmen) und der Rücksprung nach `open/` (§4) sind genau für den roten Lauf vorgesehen. Aus dem Quelltext sehe ich keinen Grund für Rot (V-3).

**Übergabe:** Bericht an den Planner. Der Implementer hat nichts nachzubessern außer den Closure-Pflichten, sobald der Lauf auf `main` vorliegt; die Finding-Klassen gehen in die Slice-Closure §7.
