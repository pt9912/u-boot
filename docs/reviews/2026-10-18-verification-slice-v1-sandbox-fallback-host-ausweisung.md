# Verification-Report: slice-v1-sandbox-fallback-host-ausweisung — 2026-10-18 <!-- d-check:ignore (Slice-Kennung im Titel, Zitat) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan, DoD, ADR und Spec). Kein Review, keine Validierung, keine Reparatur.

**Gegenstand:** Slice `slice-v1-sandbox-fallback-host-ausweisung` (`docs/plan/planning/in-progress/`), Commit-Range `3b740aa..HEAD` (20 Commits, HEAD `3d71b88`), [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md) (`Accepted`, Entscheidungen 10 und 11), Lastenheft 0.13.0 ([`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 2). <!-- d-check:ignore (Slice-Kennung als Zitat) -->

**Eingang:** Review `docs/reviews/2026-10-18-slice-v1-sandbox-fallback-host-ausweisung.md` (F-1 HIGH, F-2/F-3 MEDIUM, F-4..F-6 LOW, F-7..F-11 INFO) und die Zusage des Implementers, alle seien behoben oder vermerkt. Geprüft wurden die Belege am Code und Text, nicht die Prosa; die Sensoren wurden selbst gefahren (Host: Colima, Docker 29.x, arm64; kein Host-go). <!-- d-check:ignore (Kennungen als Zitat) -->

**Gesamturteil:** Alle sachlichen DoD-Zeilen sind erfüllt; offen sind nur die Closure-Zeilen, die dieser Rolle vorgelagert sind (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) — Auflagen unten. Keine Diskrepanz, die die Closure blockiert; zwei Restbefunde (LOW) und eine Annahme ohne ausdrückliche Bestätigung sind in der Closure-Notiz zu führen.

## 1. Sensoren (selbst gefahren)

| Sensor | Ergebnis | Beleg |
|---|---|---|
| `make gates` | pass (Exit 0, zweimal gefahren, zuletzt am sauberen Baum) | lint, test, `coverage-gate`: 91,00 % bei Schwelle 90 %, `docs-check`: 340 Dateien, 0 Befunde, `baseline-verify` v6.18.1 OK |
| `make doc-trace` | pass | „83 Anforderung(en), 0 Waise(n)“ |
| `make doc-immutable RANGE="3b740aa..HEAD"` | pass | Range aufgelöst, 20 Commits, 340 Dateien, 0 Befunde (Accepted-ADRs unverändert) |
| `make verify-depguard` | 8 von 9 PASS, wie erwartet; `domain-isoliert` FAIL = Bestand | Log: `import cycle not allowed … port/driven → domain → port/driven (typecheck)`, nicht „domain must not depend on port“. Am Basiscommit lesend bestätigt: `scripts/verify-depguard.sh` im Bereich `3b740aa..HEAD` ohne Diff, die Zeile `domain-isoliert\|…` steht dort schon (Zeile 63), und `port/driven/template_catalog.go` importiert `domain` bereits am Basiscommit. Kein Gegenstand des Slice (Entscheidung des Projektinhabers). |
| `make test-docker` (Colima) | pass | alle Pakete `ok`, e2e-Paket 584,7 s, gesamt 590,8 s (`real`), kein `FAIL`/`panic`/`SKIP`; danach keine Container, Images oder Volumes `sbxe2e-*` auf dem Host |
| gezielt `-run StatusFiles -v -count=1` im Image `u-boot:test-docker-tools` (Mount-Muster aus `repo.mk`) | pass | `StatusFiles_InspectorAgainstEngine` 3,98 s, `_Egress` 39,2 s, `_VFS` 47,0 s, `_Daemon` 55,4 s; alle vier PASS, kein SKIP; Paket 59,3 s |

Zur Laufzeit von `make test-docker`: 590,8 s liegen unter dem Standardlimit von zehn Minuten (600 s), aber nur um 9 s; zusammen mit der Messung des Implementers (477 bis 635 s) und dem Review-Lauf (über 600 s) belegt das die Streuung. Das Limit `-timeout 30m` in `repo.mk` ist daher fachlich begründet und nach Entscheidung des Projektinhabers (2026-10-18, „-timeout 30m in repo.mk“) gedeckt.

## 2. DoD-Abgleich (jede Zeile)

| # | DoD-Zeile | Stand | Evidence (am Code geprüft) |
|---|---|---|---|
| 1 | Ermittlung und Ausweisung (Code), Tests je Fallback und je Befehl, Replay/Golden | erfüllt | Wertobjekt `internal/hexagon/domain/sandbox_status.go` (geschlossene Menge, 512 Byte, strikter UTC-Zeitstempel, Fehlertexte ohne Dateiinhalt); Port `internal/hexagon/port/driven/devcontainer_inspector.go` (`FindContainer`, `ReadFile`, drei Sentinels); Adapter `internal/adapter/driven/docker/inspector.go` (`ps`, `cp … -`, Tar-Kopf, Kennungs-Muster, 10-s-Limit, entfernte Engine = nicht ermittelbar); Application `devcontainer_fallback.go` (drei Sonden, Abbildung, Warnungen), `doctor_sandbox.go` (`checkSandboxFallbacks`), `generate.go` (`existingContainerFallbackWarnings` nach `writeDevcontainerPhase`, nicht bei `PreviewDryRun`), `initproject.go`, Wiring `cmd/uboot/main.go` an `doctorSvc`, `generateSvc`, `initSvc`; Tests: `TestDoctor_SandboxFallback_*` (Happy, nicht ermittelbar, ungültig, gestoppt, Engine-Fehler, im Devcontainer), `TestGenerateDevcontainer_ExistingContainer*`, `…DryRunDoesNotAskTheEngine`, `TestInit_*` (Warnung, still bei ok/failed/ungültig/Engine-Fehler, kein Aufruf bei `--dry-run`), CLI `doctor_fallback_test.go` (`--quiet`, `--strict`, `--json`, Volltext), Adapter-Tests gegen Fake-Prozess, Golden/Replay `TestGenerateDevcontainer_SandboxStartScriptsWriteStatusFiles` und `…StatusDirectoryOnlyWithStartScripts`, Zeilen in `harness/replay.md`. Templates `sandbox-init.sh.tmpl`, `docker-init.sh.tmpl`, `egress-init.sh.tmpl`: `write_status` in jedem Zweig (`ok`, `fallback`, `failed` direkt vor `exit 11`), Temp-Datei und `mv -f`, Rückgabe immer 0, Exit-Codes und Stderr-Zeilen unverändert; `Dockerfile.tmpl` legt das Verzeichnis nur bei `NestedPodman`/`NestedDocker`/`Egress` an. Einschränkung: der `init`-Pfad ist in Produktion nicht erreichbar (siehe D-2). |
| 2 | Spezifikation und Sicht nachgezogen | erfüllt | `LH-FA-DEV-007.a`, `LH-FA-DEV-008.a`, `SPEC-043` Ausnahme 2 ohne „soweit vom Host aus ermittelbar“ (der Rest-Wortlaut steht nur noch im Lastenheft-Absatz zu `u-boot doctor` und `/dev/fuse`, vom CR gesetzt); `SPEC-045` (Statusformat, Ermittlung, Stand, Abbildung, Grenze); drei Kennungen in der `SPEC-017`-Registry und in `DefaultAllowedCodes` (Drift-Test grün); `spec/architecture.md` (Port, Adapter, `ARC-008`) ohne ADR-, Slice- oder Wellen-Verweis (§3.4); Historie-Zeile in §7. | <!-- d-check:ignore (Kennungen als Zitat) -->
| 3 | Doku-Update (Handbuch, CHANGELOG) | erfüllt | Handbuch: Abschnitt „Fallbacks sehen“, `--quiet`-Zeile, doctor-Abschnitt, Daemon-, Egress- und Fehlerbehebungs-Absatz; CHANGELOG `[Unreleased]` → Added. Aussagen zu `init` und Startkette sind ehrlich (F-2, F-6); Praxislauf bestätigt die Handbuch-Tabelle (siehe §4). |
| 4 | `make gates` grün; `make test-docker` grün (Colima); `make doc-trace` ohne neue Waise | erfüllt | §1 |
| 5 | Folge-ADR `Accepted`, im Index | erfüllt (bereits abgehakt) | `5bb01f7` vor dem Basiscommit; `doc-immutable` ohne Befund |
| 6 | Review durchgeführt, Report liegt vor | erfüllt (abzuhaken) | `docs/reviews/2026-10-18-slice-v1-sandbox-fallback-host-ausweisung.md`; Findings siehe §3 | <!-- d-check:ignore (Kennungen als Zitat) -->
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | offen | §7 des Slice steht auf „Wird bei der Closure gefüllt.“ — Auflage A-1 |
| 8 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 | offen | nicht geschrieben — Auflage A-3 |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | offen | nur das Risiko `/dev/fuse` trägt „eingetreten“, die übrigen stehen auf „offen“ — Auflage A-2 |
| 10 | Die drei Paarungen (Anker, Folge-Slice, Register) | offen | Closure-Schritt — Auflage A-1/A-3 |

## 3. Review-Findings am Code und Text

| Finding | Stand | Prüfung |
|---|---|---|
| F-1 HIGH Zeitlimit | behoben | `repo.mk`: `go test -tags docker -timeout 30m ./...` mit Kommentar (Zeitbudget, kein Test gestrichen, keine Prüfung gelockert); im eigenen Lauf grün bei 590,8 s. Die vorübergehende Änderung `40bb659` wurde mit `8340f51` zurückgenommen, `13fd933` setzt sie nach der Entscheidung des Projektinhabers neu. Plan-Tabelle nennt `repo.mk`. |
| F-2 MEDIUM `init`-Wortlaut | behoben | `LH-FA-DEV-007.a`, `SPEC-045`, CHANGELOG und Handbuch sagen übereinstimmend, dass `init` fragt, aber nichts abzufragen hat. Praxislauf: `init --devcontainer` im Ordner mit beschriftetem Container (alle drei Dateien `fallback`) erzeugt `u-boot.yaml` ohne `nestedRuntime` und ohne Egress und gibt keine Warnung aus. Der Pfad ist nur mit Fake-Codec testbar (D-2). | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-3 MEDIUM Devcontainer-Betrieb | im Wortlaut behoben, ein Rest | `SPEC-045` sagt jetzt, `generate`/`init` fragten immer über die Engine und im Devcontainer bleibe die Abfrage still. Für `doctor` (Dateisystem-Pfad) pinnt `TestDoctor_SandboxFallback_InsideContainer`. Für `generate`/`init` im Devcontainer pinnt kein Test etwas; „erreicht keine Engine“ gilt nicht allgemein (Rest D-1). | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-4 LOW Kommentar | behoben | `devcontainer_sandbox.go`: Kommentar nennt `doctor` und `generate devcontainer` über `fallbackWarnings`; Verweis auf `init` fehlt dort sachlich zu Recht (kein Erzeuger). |
| F-5 LOW Beleg-Zeiger | behoben | `harness/replay.md` nennt `internal/e2e/sandbox_status_docker_test.go`; die Tests stehen dort. |
| F-6 LOW Stand-Aussage | behoben | `SPEC-045` Absatz „Stand“ und „Nicht ermittelbar“ nennen den Abbruch der Kette `egress && docker` und blockierte User-Namespaces; Handbuch ebenso; entspricht `devcontainerTemplateData` (Egress vor Daemon). | <!-- d-check:ignore (Kennungen als Zitat) -->
| F-7 INFO `domain-isoliert` | Bestand, nicht Gegenstand | selbst bestätigt (§1); Harness-Pflege, eigener Vorgang. |
| F-8 INFO Vertrauensgrenze | akzeptiert | Praxislauf bestätigt: Freitext mit Steuerzeichen, Wert außerhalb der Menge und Symlink auf `/etc/passwd` ergeben nur „invalid“, nie Dateiinhalt in der Ausgabe. |
| F-9 INFO Commit-Message `816f6ec` | vermerkt | beschreibt zurückgenommenen Code; der Baum trägt die Aussage nicht; Hinweis in die Closure-Notiz (Auflage A-1). |
| F-10 INFO Diagnosegrenzen | akzeptiert | ADR „Nicht gemessen“; Hinweis in die Closure-Notiz. |
| F-11 INFO Aufräumen | behoben | `3d71b88`: Volume-Cleanup vor dem Container registriert (LIFO); nach `make test-docker` und dem gezielten Lauf keine Reste auf dem Host. |

## 4. Praxislauf vs. Test-Beleg

**Praktisch geprüft** (Docker-only, Wegwerf-Verzeichnis, danach vollständig entfernt, ebenso Container und Volumes `sbxe2e-*`; Baum unverändert, `git status` leer):

- `doctor` aus dem Runtime-Image `u-boot:latest` (Container-Erkennung aktiv) mit eingehängtem Statusverzeichnis `/var/lib/u-boot/status` und Projekt mit `nestedRuntime: docker` plus Egress-Restriktion. `ok`: zwei OK-Zeilen, Exit 0. `fallback`: je Eintrag `⚠` mit Zeitpunkt und Hinweis, Exit 0, mit `--strict` Exit 11. `failed`: zwei `✗`, Exit 11. Ungültig (Freitext mit Steuerzeichen, `result=weird`, Symlink): `⚠` „invalid“, ohne Inhalt; mit `--strict` Exit 11. Fehlend und Verzeichnis nicht vorhanden: `?` „not determinable“ mit Grund, Exit 0, auch mit `--strict`.
- `--quiet`: nur `fallback` und `failed` sichtbar, „nicht ermittelbar“ und `ok` stumm. `--json`: `diagnostics[]` trägt die Fallback-Warnungen, kein `info`-Eintrag, `exitCode` 0.
- `nestedRuntime: podman` im Runtime-Image: der Eintrag zu `/dev/fuse` (`devcontainer.sandbox.runtime`) ist dort ein Fehler, weil der Container selbst kein `/dev/fuse` hat; das ist das Bestandsverhalten aus dem vorgelagerten CR, kein Befund dieses Slice, und der Grund, warum der `vfs`-Eintrag nur in der „nicht ermittelbar“-Variante gesehen wurde.
- `generate devcontainer`, mit dem Docker-gebauten Binary (Image `u-boot:test-docker-tools`, Socket gemountet) und einem Container mit Label `devcontainer.local_folder=<Projektordner>`: ohne Container still; mit `fallback` zwei Warnungen („state of the existing devcontainer container at its last start (…)“, Neubau-Hinweis), auch unter `--quiet`, in `--json` im Warnungsfeld; gestoppter Container wird gelesen; `ok`, `failed` und ungültig still; `--dry-run` fragt nicht.
- `init --devcontainer` im Ordner mit beschriftetem Container: keine Abfrage-Warnung (wie in Spec und Handbuch).

**Belegt durch Tests, nicht praktisch geprüft:** die echten Startscripte mit Statusdateien und das Lesen durch Host-Adapter und `doctor` (`StatusFiles_*`, §1); Inspector gegen Symlink, FIFO, Verzeichnis, Übergröße, mehrere Container, entfernte Engine (e2e und Adaptertests); `init` mit Fähigkeit (Fake-Codec, D-2); Devcontainer-Betrieb von `doctor` mit nested Daemon. Nicht verfügbar: Devcontainer-CLI (Annahme der ADR), Linux ohne VM, amd64, VS Code.

## 5. Plan-Tabelle §3 gegen `git diff --stat 3b740aa..HEAD` (§3.10)

Alle 33 berührten Dateien stehen in der Tabelle; „unverändert“-Zeilen (`cli-json-contract.md`, READMEs, `open/README.md`, `carveouts.md`, `.golangci.yml`, `.d-check.yml`, `spec/lastenheft.md`) haben keinen Diff; die beiden Zeilen „in der Planung nicht aufgeführt“ (`devcontainer_fallback.go`, `coderegistry.go`) sind vermerkt, `repo.mk` ist aufgenommen. Eine Zeile ohne Diff in der Range: `docs/plan/adr/` samt Index (Arbeit in `5bb01f7` und `fae000a`, vor dem Basiscommit; in der Zeile als „erledigt“ vermerkt) — kein Befund.

## 6. Diskrepanzen und Annahmen

- **D-1 (LOW, Spec-Wortlaut ungepinnt):** `SPEC-045` sagt, im Devcontainer erreiche der Engine-Weg von `generate`/`init` keine Engine. `GenerateService` und `InitProjectService` kennen `RuntimeEnvironment` nicht; mit einem verschachtelten Docker-Daemon erreicht der Aufruf den verschachtelten Daemon, der keinen Container mit dem Label des Arbeitsordners kennt (still), mit dem Podman-Alias scheitert er (still). Wirkung (still) stimmt, die Begründung („keine Engine“) und die fehlende Pinnung sind nicht belegt. Vorschlag für den Planner: Satz auf „bleibt dort still“ kürzen oder mit einem Test pinnen. <!-- d-check:ignore (Kennungen als Zitat) -->
- **D-2 (LOW, bekannt, F-2):** Die `init`-Abfrage ist in Produktion ein toter Pfad (das geschriebene `u-boot.yaml` konfiguriert keine Fähigkeit); die DoD-Anforderung „Test je Befehl `init`“ ist mit Fake-Codec (`yamlWithNestedRuntime`) erfüllt, der Produktionspfad erreicht sie nicht. Spec, Handbuch und CHANGELOG sagen es offen. In die Closure-Notiz als „Grenze, bewusst“ aufnehmen.
- **A-ADR (Annahme ohne ausdrückliche Bestätigung):** „keine Abfrage bei `--dry-run`“ steht in `ADR-0022` Entscheidung 10 (angenommen, `Accepted`), ist im Slice (§1), Spec, Handbuch und `harness/replay.md` übernommen; der Projektinhaber hat sie nicht ausdrücklich bestätigt. Code und Tests entsprechen der ADR; eine Änderung wäre eine Verhaltensänderung an `generate`/`init`. In §7 ausweisen. <!-- d-check:ignore (Kennungen als Zitat) -->
- **D-3 (INFO):** Das Register führt „Spec-Bedingung ohne Code-Entsprechung“ noch nicht; im Verzeichnis `docs/plan/planning/observations/` gibt es keinen Eintrag zu dieser Klasse. Der Review nennt „einmal im Vorlauf (CR, F-1), hier zweimal (F-2, F-3)“. Nach der Zählregel des Registers („ein Vorgang zählt einmal“) sind das zwei Vorgänge (CR und dieser Slice), nicht drei. Die Verkörperung beruht auf der Entscheidung des Projektinhabers (2026-10-18); die Closure-Notiz sollte die Zählbasis ehrlich nennen (Beobachtungen: 3 Findings, 2 Vorgänge) und nicht „3×“ behaupten, sonst weicht sie von der Register-Regel ab.

## 7. Auflagen für die Closure

- **A-1 Closure-Notiz (§7) mit Steering-Loop-Lerneintrag:** Finding-Klassen des Reviews (Zeitbudget des Integrations-Targets ohne Reserve; Spec-Bedingung ohne Code-Entsprechung 2×; Kommentar nennt Erzeuger unvollständig; Beleg-Zeiger auf falsche Datei; Stand-Aussage ohne Abbruchpfad; Sensor-Fall nicht durchführbar; Vertrauensgrenze akzeptiert; Commit-Message beschreibt zurückgenommenen Code; Diagnosegrenze ohne Log; Abbruch ohne Aufräumen); Hinweise F-7 (Bestand `domain-isoliert`, eigener Harness-Vorgang), F-9 (Message von `816f6ec` nennt Zurückgenommenes; die Message bleibt, der Baum ist sauber), F-10 (Grenzen ohne Log, nicht gemessene Umgebungen), D-1, D-2 und die Annahme zu `--dry-run`; „bei `init` nichts abzufragen“ als Grenze. DoD-Zeilen 6 bis 10 erst nach Schreiben abhaken.
- **A-2 Risiko-Ausgänge §6** (jeweils genau einer): Messung auf Colima trägt nicht → entfallen für Colima/arm64 (gemessen M1–M14, e2e grün), weiter offen für Linux ohne VM, amd64, VS Code; Floating Container-Zustand → eingetreten als Eigenschaft (Zeitstempel und „container stopped“ im Eintrag, Kette `egress && docker`), mitigiert; Container-Identifikation → weiter offen (nur Devcontainer-CLI-Label, mehrere Container: zuletzt erzeugter); `docker`-Zugriff in `doctor` → entfallen (Fehler ergeben `info`, keine Exit-Code-Verschiebung, belegt); Vertrauen in das Dokument → entfallen als Risiko der Ausgabe (Praxislauf, Parser-Tests), als Grenze weiter dokumentiert; `--quiet`-Sichtbarkeit → entfallen (Handbuch und Spec benennen es); Docker-only-Tests → entfallen (Lauf mit Engine grün, 590,8 s, ohne Skip); Änderung der erzeugten Ausgabe → eingetreten (Startscripte und Dockerfile ändern sich nach `generate`), im CHANGELOG genannt; Spec-Nachzug ohne Verweise nach unten → entfallen (`docs-check` grün); `/dev/fuse` → eingetreten (bereits vermerkt); Abfrage in `generate`/`init` → entfallen (Zeitlimit 10 s je Aufruf, alle Fehler still, Tests); Warnung beschreibt alten Container → entfallen (Wortlaut nennt „last start“ und Neubau); Reihenfolge CR → Slice → entfallen (CR in `done/`); Umfang → entfallen (Slice nicht zerlegt, e2e-Zeit als Folgekosten notieren: `test-docker` rund zehn Minuten).
- **A-3 Beobachtungs-Register:** Eintrag der Klasse „Spec-Bedingung ohne Code-Entsprechung“ als `BEO-DOK/<slug>` mit `evidence/` je Vorgang (CR `slice-cr-sandbox-fuse-startfehler`, dieser Slice) und Ausgang nach Entscheidung des Projektinhabers (Verkörperung in `AGENTS.md` §3.10, mit Herkunfts-Anker `seit slice-v1-sandbox-fallback-host-ausweisung` bzw. der Welle, wenn sie den Slice aufnimmt, und `liegt in`-Feld in der Closure-Notiz, Anker-Paarung); Zählbasis ehrlich benennen (D-3). Für `BEO-DOK/plan-aussage-ueber-code`: ein Rückfall (Review F-4, F-5 und F-6 betreffen Aussagen über Code im Plan und Beleg-Zeiger) führt nach Entscheidung zur zweiten Schärfung von `AGENTS.md` §3.10; Retirement-Check beachten (der Anker bleibt, die Regel wird nur geschärft, nicht gelockert). Bei der Schärfung ist die Prüfung „Wortlaut gegen Code“ für Aussagen der Spec zu Randbedingungen (D-1) mitzunehmen. <!-- d-check:ignore (Kennungen als Zitat) -->
- **A-4 Roadmap-Ruhe-Marker nach dem Move:** `git mv` nach `done/` als reiner Move-Commit (erst Inhalt, dann Move, `AGENTS.md` §3.3). Mit der Closure liegt nichts mehr in `in-progress/`: der Marker „Nichts in Arbeit.“ gehört zurück in `roadmap.md` unter „Offene Wellen“ (Commit `2374d64` hat ihn entfernt, solange der Slice dort liegt; d-check `planning`), `next/README.md` und `done/README.md` nachziehen; Verweise vom Ruheort aus schreiben (Ruheort-Regel), Pfade in der Slice-Datei prüfen (`../done/…`, `../observations/…`).
- **A-5 Öffentliche Verträge:** `CHANGELOG.md` ist nachgezogen; beim Release-Schnitt bleibt der Eintrag unter `[Unreleased]`. Keine weitere Doku offen.
- **A-6 Hygiene-Befunde:** `verify-depguard` `domain-isoliert` und gofmt-Bestand sind nicht Gegenstand dieses Slice (Entscheidung des Projektinhabers); als eigene Vorgänge führen, nicht in der Closure-Notiz als erledigt behaupten.

## 8. Verification Evidence (Block nach `harness/verification.md`)

Scope:
- Slice: `slice-v1-sandbox-fallback-host-ausweisung` <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- IDs: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-005`](../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md)
- Artefakte: `sandbox_status.go`, `devcontainer_inspector.go`, `docker/inspector.go`, `devcontainer_fallback.go`, `doctor_sandbox.go`, `generate.go`, `initproject.go`, `cmd/uboot/main.go`, vier Templates, `repo.mk`, e2e `sandbox_status_docker_test.go`, Spec und Handbuch

Sensors:

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | Exit 0, 91,00 % Coverage, `docs-check` 0 Befunde |
| `make test-docker` | pass | 590,8 s, kein Skip |
| `make doc-trace` | pass | 0 Waisen |
| `make doc-immutable RANGE="3b740aa..HEAD"` | pass | 0 Befunde |
| `make verify-depguard` | 8 von 9 PASS | `domain-isoliert` Bestand (Importzyklus, am Basiscommit lesend bestätigt) |

Traceability:

| ID | Beleg |
|---|---|
| `LH-FA-DEV-007` | `StatusFiles_*` (e2e), `TestDoctor_SandboxFallback_*`, `TestGenerateDevcontainer_ExistingContainer*`, `SPEC-045`, `SPEC-017`-Registry | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-DEV-008` | `StatusFiles_Egress` (e2e), Egress-Sonde und Test je Zustand | <!-- d-check:ignore (Kennungen als Zitat) -->
| `LH-FA-CLI-005` Ausnahme 2 | `doctor_fallback_test.go` (`--quiet`), `SPEC-043`, Praxislauf | <!-- d-check:ignore (Kennungen als Zitat) -->
| `ADR-0022` | Port, Adapter, Wertobjekt, `doc-immutable`, Praxislauf (Entscheidungen 6 bis 10); Entscheidung 10 `--dry-run` ungeprüft durch Projektinhaber | <!-- d-check:ignore (Kennungen als Zitat) -->

Carveouts:
- Neu: none · Gelöst: none · Unverändert: bestehende Einträge

Nicht ausgeführt:
- `make ci` (`govulncheck`, `image-scan`) — Sicherheits-/Release-Pfad nicht berührt, kein Abhängigkeits- oder Image-Wechsel im Diff
- Devcontainer-CLI-Lauf, Linux ohne VM, amd64, VS Code — nicht verfügbar; in der ADR als nicht gemessen geführt

Commit / Artefakt:
- HEAD `3d71b88` (Range `3b740aa..HEAD`)
