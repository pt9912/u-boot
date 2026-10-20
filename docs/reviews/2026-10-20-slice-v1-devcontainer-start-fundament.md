# Review-Report: [`slice-v1-devcontainer-start-fundament`](../plan/planning/in-progress/slice-v1-devcontainer-start-fundament.md) — 2026-10-20

**Review-Art:** Code-Review gegen Plan, Entscheidungen und Hard Rules (Produktcode, Docker-Tests, Messung). Die DoD prüft der Verifier.

**Gegenstand:** `git diff eea0c92..3fb0c1b` (Code `d9780d3`, Messung `558b369`, Slice-Dokument `b0a4ba3`, Plan-Nachzug im Befehl-Slice `3fb0c1b`); 17 Dateien, 2188 Einfügungen.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `3fb0c1b` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-20

**Eingangs-Kontext:** [`slice-v1-devcontainer-start-fundament`](../plan/planning/in-progress/slice-v1-devcontainer-start-fundament.md) (§2, §3, §6, §7), [`slice-v1-cli-devcontainer-up-down`](../plan/planning/open/slice-v1-cli-devcontainer-up-down.md) (nur die Stellen, die der Fundament-Slice berührt), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md), [`ADR-0029`](../plan/adr/0029-devcontainer-start-exit-11-und-socket-haltung.md), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte), `spec/architecture.md` §2.4 und §2.6, `.golangci.yml` (`depguard`), [`AGENTS.md`](../../AGENTS.md) (Hard Rules, §3.1, §3.2, §3.6, §3.7, §3.9, §3.10).

**Selbst gelesen (nichts vom Implementer übernommen):** `internal/adapter/driven/docker/devcontainer_launcher.go`, `devcontainer_launcher_test.go`, `devcontainer_launcher_docker_test.go`, `export_test.go`, `inspector.go` (`requireLocalEngine`, `limitedBuffer`), `engine.go` (`progressSinkOrDiscard`), `devcontainercli/{Dockerfile,package.json,package-lock.json}`, `internal/hexagon/port/driven/devcontainer_launcher.go`, beide Messskripte, Messprotokoll mit Rohdaten, Suche nach `ErrDevcontainerRuntime` im ganzen Baum.

**Sensor-Läufe des Reviewers:**

- `make gates`: Exit 0.
- `make doc-trace`: `84 Anforderung(en), 0 Waise(n)`.
- `make doc-immutable RANGE="eea0c92..HEAD"`: 0 Befunde.
- Zehn Mutationen am Produktcode oder am eingebetteten Kontext, je einmal angewendet und mit `git checkout -- <Datei>` zurückgenommen (Läufe im Image `test-docker-tools` per `docker run`, `git status` danach sauber).

  | Mutation | Wächter | Ergebnis |
  |---|---|---|
  | `--rm` entfernt | `TestRunArgs` | rot |
  | `--rm` entfernt | `TestLauncher_Docker_S6_Lifetime/after_Up` (Docker, mit Engine) | rot (`helper containers after Up`); danach kein Hilfs-Container auf dem Daemon |
  | `docker rm -f` im Fehlerpfad entfernt | `TestLauncher_Up_CancelRemovesHelperContainer`, `TestLauncher_Up_RuntimeFailures` | rot |
  | `outcome`-Prüfung entfernt | `TestParseUpResult` | rot |
  | `requireLocalEngine` aus der Vorab-Prüfung entfernt | `TestLauncher_Preflight` | rot |
  | Spiegelpfad `target=/w`+Ordner | `TestRunArgs` | rot |
  | `-e NAME` → `-e NAME=x` | `TestLauncher_Up_PassesEnvByNameOnly` | rot |
  | Digest beim Node-`FROM` entfernt; `npm ci --ignore-scripts` → `npm ci` | `TestHelperContext_Pins` | je rot |
  | `--force-rm` entfernt | `TestLauncher_Up_BuildsHelperImageOnlyWhenTagIsMissing` | rot |

- Zusätzliche Probe (Dev-Container-CLI im Hilfs-Image, Projekt mit `containerEnv` und `remoteEnv` auf `${localEnv:…}`, Platzhalter-Wert): Befund F-1.
- `gh run list`: `ci` zu `3fb0c1b` grün; `integration-docker` zu `3fb0c1b` zum Zeitpunkt des Reports noch `in_progress` (nach 7 min 44 s). Information, kein Beweis.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Dev-Container-CLI schreibt bei jedem `up` ihre `docker run`-Zeile des Dev-Containers auf stderr; mit `"containerEnv": {"X": "${localEnv:NAME}"}` enthält sie `-e X=<Wert>` im Klartext (reproduziert: Platzhalter-Wert in stderr, nicht in stdout). Der Adapter leitet stderr an `opts.Progress` und hängt den Tail (4096 Byte) an die Fehlermeldung (`ErrDevcontainerRuntime`); der Wert steht damit auf der Konsole, in CI-Logs und bei einem Fehler im Fehlertext. Der Test `…S6_ValueIsNeverInImageLabelFileOrCommandLine` nutzt nur `remoteEnv` (taucht in dieser Zeile nicht auf) und liest den Fortschrittspuffer nicht. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 (Umgebung) | `internal/adapter/driven/docker/devcontainer_launcher.go` · `cmd.Stderr = io.MultiWriter(stderr, progressSinkOrDiscard(opts.Progress))` und `devcontainer up failed (%s: %s)`; `devcontainer_launcher_docker_test.go` · `"remoteEnv":         map[string]string{"UBL_SECRET"` | ja — Docker-Test mit `containerEnv`, Prüfung des Fortschrittspuffers | Zugangsdaten-Wert im Fortschrittsstrom |
| F-2 | LOW | Eingabefehler (`validateProjectFolder`, `passedEnvNames`, ungültige Container-ID) liefert der Adapter als `ErrDevcontainerRuntime`, der Test pinnt das (`RejectsUnsafeInput`). Die Port-Doku bindet den Sentinel an einen Fehler „after the environment checks passed“; die Prüfung läuft davor. Der Befehl-Slice erwartet für den Komma-Pfad Exit 10 und plant die Pfad-Prüfung in der Application: dieselbe Zeichenmenge steht dann an zwei Stellen mit zwei Klassen. | [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 7 (Fehlerklassen), Maintainability | `internal/hexagon/port/driven/devcontainer_launcher.go` · `failed after the environment checks passed`; `devcontainer_launcher.go` · `project folder %q is not an absolute path a --mount value can carry` | nein — Review-only | Fehlerklasse der Eingabeprüfung weicht von der Port-Doku ab |
| F-3 | LOW | `limitedBuffer` behält die ersten 1 MiB und verwirft den Rest; `parseUpResult` liest die letzte nicht leere Zeile. Übersteigt stdout 1 MiB, fehlt die Ergebniszeile, und die Meldung lautet „last stdout line is not a JSON result“ (oder „no result line“). Ob die CLI Hook-Ausgabe nach stdout schreibt, ist im Slice nicht belegt. | Maintainability ([`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 3: Ergebniszeile = letzte nicht leere Zeile) | `devcontainer_launcher.go` · `stdout := &limitedBuffer{limit: upStdoutLimit}` | nein — kein Test mit großer Ausgabe | Puffer behält Kopf, Parser liest Ende |
| F-4 | LOW | Plan-Tabelle und §7 sagen, der Allaussage-Test prüfe „die Kommandozeilen aller Prozesse während des Laufs (`/proc`)“. `make test-docker` startet den Test-Container ohne `--pid=host`: `/proc` zeigt nur dessen Prozesse (den `docker`-Client). Die Prozesse im Hilfs-Container (die CLI ruft dort `docker`-Befehle auf) sind nicht erfasst. Eine Aussage über Verhalten ist weiter formuliert als das Gelesene. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Aussage über Verhalten nennt nur, was gelesen oder gemessen wurde) | Slice [`slice-v1-devcontainer-start-fundament`](../plan/planning/in-progress/slice-v1-devcontainer-start-fundament.md) · `die **Kommandozeilen** aller Prozesse während des Laufs` ; `repo.mk` · `docker run --rm --network=host` | nein — Review-only | Prüfreichweite kleiner als die formulierte Zusage |
| F-5 | LOW | Der Slice ist nach den Entscheidungen des Projektinhabers (`3fb0c1b`) nicht nachgezogen. (a) §7 „Stand: 2026-10-20, HEAD `eea0c92`, nichts committet“ ist überholt. (b) §7 Punkt 3 „Entscheidung des Architect erbeten“ und „Offene Fragen an den Architect (a) bis (d)“ sind entschieden. (c) §6 führt weder `readonly`/schreibend, Abbruchdauer bis zur Frist noch Größe und Funde des Basis-Images (1716 MiB, 597 HIGH/CRITICAL, Messprotokoll) als Risiko mit Ausgang; die DoD-Zeile „Jedes Risiko aus §6 trägt einen Ausgang“ erfasst sie nicht. (d) §3 beschreibt den Umgebungs-Wächter noch als „Variante alle Variablen durchreichen“ und den Hilfs-Container-Inspect-Test; §7 Punkt 1 und 2 haben beides geändert. (e) Die Plan-Tabelle hat keine Zeile für [`slice-v1-cli-devcontainer-up-down`](../plan/planning/open/slice-v1-cli-devcontainer-up-down.md) (im Diff, `3fb0c1b`); die Commit-Nachricht nennt den Fundament-Slice, der Commit ändert ihn nicht. | [`AGENTS.md`](../../AGENTS.md) §3.10 (Zustandsmarker mit Ablaufpunkt, Abgleich Plan-Tabelle und `git diff --stat`) | [`slice-v1-devcontainer-start-fundament`](../plan/planning/in-progress/slice-v1-devcontainer-start-fundament.md) · `Stand: 2026-10-20, HEAD` und `Entscheidung des Architect erbeten` und `**Offene Fragen an den Architect.**` | nein — Review-only | Zustandsmarker abgelaufen, Risiko-Liste unvollständig |
| F-6 | LOW | DoD und Slice nennen „nur `make` und Docker am Host“; die Skripte rufen `bash`, `tar` mit GNU-Optionen (`--mtime`, `--owner`), `date +%s.%N`, `stat -c`, `awk`, `sed`, `grep`, `cut`, `mktemp`, `head`, `tail` auf. Der Kopf und das Messprotokoll nennen nur „Docker, `tar` und `date`“. `date +%N` und `stat -c` laufen auf macOS nicht, wo der Projektinhaber die Colima-Einzelmessung macht. Im Messprotokoll ist der Absatz „Jedes Skript wird im Verzeichnis … `P` ist der Pfad …“ nicht lesbar. | [`BEO-WZG/host-werkzeug-ausserhalb-docker`](../plan/planning/observations/BEO-WZG/host-werkzeug-ausserhalb-docker/observation.md), Slice §2 (Einmal-Messung) | `tools/messung/devcontainer-start-fundament/01-hilfs-image.sh` · `Voraussetzung: nur Docker und tar/date auf dem Host`; Messprotokoll [`README.md`](../messungen/slice-v1-devcontainer-start-fundament/README.md) · `Voraussetzung: nur Docker, `tar` und` | nein — Review-only | Host-Voraussetzung untertrieben, Linux-spezifische Aufrufe |
| F-7 | LOW | Die Docker-Tests räumen über das Eigene hinaus auf: `newProject` entfernt per `removeContainers("name=" + Präfix)` alle Hilfs-Container des Daemons, und `helperLeftovers` zählt alle. Ein parallel laufender echter Hilfs-Container auf demselben Daemon wird entfernt oder macht `after_Up` rot. In `…S6_ValueIsNeverIn…` bleibt der `scanDuring`-Goroutine bei einem `t.Fatalf` in `mustUp` aktiv (`close(stop)` fehlt im Fehlerpfad). | Slice §3 „Aufräumen: jeder Test entfernt, was er anlegt“, [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../plan/planning/observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) | `devcontainer_launcher_docker_test.go` · `removeContainers("name=" + docker.HelperContainerPrefixForTest)`; `hits := scanDuring(value, stop)` | nein — Review-only | Aufräumen reicht über den Test hinaus |
| F-8 | INFO | `removeHelper` verwirft das Ergebnis von `docker rm -f` (`_ =`); schlägt es fehl, bleibt der Hilfs-Container (und mit ihm der Socket) bestehen, `Up` meldet nur den Kontextfehler. Der Port sagt „after the helper container has been removed“. Den Fehlerfall pinnt kein Test. | [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 (Lebensdauer) | `devcontainer_launcher.go` · `_ = exec.CommandContext(rmCtx, l.binary, "rm", "-f", name).Run()` | nein | Aufräumen ohne Rückmeldung |
| F-9 | INFO | Der S2-Wächter („Messweg `/w/…`“) ist eine Probe am Messweg, keine Mutation am Produktcode; er wird nur rot, wenn das Label des Messwegs dem Host-Pfad gleicht. Der Haupttest S2 wird mit der Mutation `target=/w`+Ordner rot (beim Start), nicht über das Label — der Abbruch ist der ehrlich berichtete Ausgang (§7 Mutations-Tabelle). | [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9 (Wächter je Zusage) | `…docker_test.go` · `TestLauncher_Docker_S2_GuardMeasurementPathDeviates` | ja — `make test-docker` | Wächter ohne Produktcode-Mutation |
| F-10 | INFO | Der Befehl-Slice nimmt Volume-Test, `readonly`-Haltung, Abbruchdauer und Handbuch-Hinweise (Befunde 8 und 9) auf (`3fb0c1b`); Zuständigkeit Verifier: ob die Einzelmessung zum Basis-Image (Größe, Funde) als Handbuch-Zahl oder als Risiko zählt, entscheidet der Projektinhaber — im Fundament-Slice ist sie nur Befund 7. | Slice §2 / §6 | [`slice-v1-cli-devcontainer-up-down`](../plan/planning/open/slice-v1-cli-devcontainer-up-down.md) · `Der Projektordner ist im Hilfs-Container schreibend eingehängt` | nein | Hinweis an Verifier |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Argumentliste des Hilfs-Containers gegen [`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 4 (`runArgs`): `--rm`, `--name u-boot-devcontainer-up-<8 hex>`, `--mount type=bind` statt `-v`, `-e <NAME>` ohne Wert, Spiegelpfad (Quelle = Ziel = `--workspace-folder`), Socket nur im Hilfs-Container, kein `--privileged`/`--network`/Capabilities | geprüft, ohne Befund (Unit, Mutationen oben) |
| Komma, Anführungszeichen, Zeilenumbruch, NUL und relativer Pfad im Projektordner | geprüft, ohne Befund (abgewiesen, `RejectsUnsafeInput`; Klassifikation siehe F-2) |
| Zugangsdaten in Argumenten des Adapters, in [`docs/messungen/slice-v1-devcontainer-start-fundament/`](../messungen/slice-v1-devcontainer-start-fundament/README.md) und `tools/messung/devcontainer-start-fundament/` | geprüft, ohne Befund (Suche nach `token`, `secret`, `password`, `authorization`: 0 Treffer; Skripte setzen keine Umgebung durch; Stromweg: F-1) |
| Abbruch (SIGINT über `cmd.Cancel`, `WaitDelay`, `docker rm -f` mit `context.WithoutCancel`), Goroutinen und Prozesse im Produktcode | geprüft, ohne Befund (Docker-Test `after_a_cancelled_Up` grün gesehen, kein Hilfs-Container danach; Rückmeldung bei `rm`-Fehler: F-8) |
| Fehlerabbildung 11 gegen 12: Vorab-Prüfung (Binary, Daemon, lokal) vor Bau und Start; danach jeder Fehler `ErrDevcontainerRuntime`; Kontextabbruch ohne Sentinel | geprüft, ohne Befund |
| Dual-Classifier ([`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): `ErrDevcontainerRuntime` kommt im ganzen Baum nur im Port und im Adapter vor; `ExitCode`, `isValidationError` und die Diagnostic-Abbildungen sind nicht berührt; der Fundament-Slice führt keinen Befehl, die Einträge samt Pin-Test stehen in der Plan-Tabelle des Befehl-Slice | geprüft, ohne Befund |
| Ergebniszeile: letzte nicht leere Zeile, `containerId` gegen `[0-9a-f]{12,64}`, `outcome`, stdout und stderr getrennt | geprüft, ohne Befund (Fall große Ausgabe: F-3) |
| Lieferkette: beide `FROM` mit Index-Digest, exakte CLI-Version 0.89.0 in `package.json` und Lockdatei, `integrity` (sha512), `npm ci --ignore-scripts`, Tag aus Version und Kontext-Hash, Tar-Strom mit festen Metadaten und `docker build --force-rm -t <Tag> -`; Unit-Test der Pins (Mutationen oben) | geprüft, ohne Befund |
| Docker-Tests: `//go:build docker`, kein `t.Skip`, `requireDocker` als Hart-Fehler, Bau nach `docker rmi`, Aufräumen per `t.Cleanup` auch bei Fehlschlag | geprüft, ohne Befund (Reichweite des Aufräumens: F-7) |
| `t.Parallel` in `devcontainer_launcher_test.go` | geprüft, ohne Befund (keine Verwendung; der Kommentar an `newFakeDocker` nennt Grenze und Ursache im Indikativ, §3.7) |
| Kommentare nach §3.7 und §3.9 (a) bis (c) im Diff (Go, Dockerfile, Skripte) | geprüft, ohne Befund (Muster „früher“, „ohne X wäre“, „no longer“ nicht gefunden; Helfer stehen nach der Funktion mit Doc-Kommentar) |
| `inspector.go`: `requireLocalEngine` als freie Funktion, Rumpf unverändert, Methode ruft sie auf; `inspector_test.go` unverändert und grün | geprüft, ohne Befund |
| Port-Pin `var _ driven.DevcontainerLauncher = (*Launcher)(nil)`; Sentinel-Nutzung `ErrDockerUnavailable`/`ErrContainerEngineNotLocal`/`ErrDevcontainerRuntime` entspricht der Plan-Tabelle | geprüft, ohne Befund |
| Hard Rules: kein `//nolint`, `.golangci.yml`, `Makefile`, `repo.mk`, `Dockerfile`, `.github/` unverändert (keine Gate-Lockerung), Docker-only; `depguard`: Port importiert nur `context`/`errors`/`io`, Adapter nur Stdlib und den Port | geprüft, ohne Befund; `make verify-depguard` nicht gelaufen — keine Importregel berührt |
| Sprachvertrag ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)): Fehlertexte des Adapters | geprüft, ohne Befund (Englisch) |
| Plan-Tabelle gegen `git diff --stat eea0c92..3fb0c1b` | geprüft; alle Dateien des Diffs stehen in der Tabelle bis auf [`slice-v1-cli-devcontainer-up-down`](../plan/planning/open/slice-v1-cli-devcontainer-up-down.md) (F-5); `observations/…` steht als Closure-Eintrag, kein Diff erwartet |

## Nicht geprüft

| Bereich | Grund |
|---|---|
| `make test-docker` (volle Läufe) | Das Target führt das ganze Modul aus; ein gezielter Lauf ginge nur mit Änderung von `repo.mk`. Stattdessen drei gezielte Docker-Tests (`S6_Lifetime`, `StopStopsButKeepsTheContainer`, Mutation `--rm`) per `docker run` im Image `test-docker-tools`, ohne Dateiänderung. Der bekannte Umgebungsfehler (drei e2e-Tests, TLS zu `download.docker.com`) ist kein Befund. |
| Mutationen der Docker-only-Wächter S2, S5, S6-Umgebung, `readonly` | Laufzeit je mehrere Minuten; nur die Lebensdauer-Mutation nachvollzogen. Die übrigen Zeilen der Mutations-Tabelle in §7 sind nicht reproduziert. |
| `arm64`/Colima, Docker Desktop, Podman, rootless Docker | nicht zugesagt ([`ADR-0028`](../plan/adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) Entscheidung 9), kein Host verfügbar |
| `make verify-depguard` | keine Importregel berührt |
| `integration-docker` zu `3fb0c1b` | zum Zeitpunkt des Reports `in_progress`; das Ergebnis ist Information, kein Beweis |
| DoD, Closure-Notiz, Verification-Evidence | Aufgabe des Verifiers |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 6 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zugangsdaten-Wert im Fortschrittsstrom · Fehlerklasse der Eingabeprüfung weicht von der Port-Doku ab · Puffer behält Kopf, Parser liest Ende · Prüfreichweite kleiner als die formulierte Zusage · Zustandsmarker abgelaufen, Risiko-Liste unvollständig · Host-Voraussetzung untertrieben · Aufräumen reicht über den Test hinaus

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (MEDIUM). Der Wert einer über `${localEnv:NAME}` in `containerEnv` referenzierten Variable erreicht Konsole, CI-Log und Fehlertext; die Allaussage zu Zugangsdaten ([`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)) ist in Wortlaut und Handbuch des Befehl-Slice nur so weit zu tragen, wie ein Test sie pinnt. Für das vom Projekt erzeugte Konfigurations-Muster (`remoteEnv`) trat der Wert im Probelauf nicht in stderr auf. Die LOW-Findings blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Ein Self-Review der Behebung entfällt: Nachlauf als neue Datei (`-r2`).
