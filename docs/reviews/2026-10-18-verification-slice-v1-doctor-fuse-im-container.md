# Verification-Report: slice-v1-doctor-fuse-im-container — 2026-10-18 <!-- d-check:ignore (Kennung im Titel) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan und DoD; kein Review, keine Validierung).

## Scope

- Slice: [`slice-v1-doctor-fuse-im-container`](../plan/planning/in-progress/slice-v1-doctor-fuse-im-container.md) (in-progress), Commit-Range `e3d328e..HEAD` (`f544d13`, 8 Commits, 10 Dateien). E1 = Lösung (b), E2 = nur der `/dev/fuse`-Zweig.
- IDs: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Lastenheft 0.13.0, im Diff unverändert), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); Spezifikation [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), [`SPEC-045`](../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) (gelesen, unverändert); [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md).
- Eingang: DoD-Bestätigung des Implementers, Review-Report [`2026-10-18-slice-v1-doctor-fuse-im-container.md`](2026-10-18-slice-v1-doctor-fuse-im-container.md) (0 HIGH/MEDIUM, 3 LOW, 5 INFO laut Tabelle: F-1..F-3 LOW, F-4..F-7 INFO).
- Host: Colima, Docker-Server arm64, VM-Kernel mit `kernel.apparmor_restrict_unprivileged_userns=0`; kein Host-Go.

## Sensors (selbst gefahren)

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass (Exit 0) | `baseline-verify: v6.18.1 OK — 54 Dateien`, `d-check: 369 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK — coverage 91.00% meets threshold 90%`; lint und test grün |
| `make doc-trace` | pass | `83 Anforderung(en), 0 Waise(n).` |
| `make doc-immutable RANGE="e3d328e..HEAD"` | pass | `8 Commit(s) — OK`, `0 Befund(e)` |
| `make docs-check` | pass | `369 Datei(en) geprüft, 0 Befund(e)` (vor Ablage dieses Reports; Nachlauf siehe unten) |
| `make test-docker` | nicht gelaufen, begründet | `git diff --stat e3d328e..HEAD`: kein Template, kein e2e-Test, kein Adapter im Diff; berührt sind `doctor.go` (Kommentar), `doctor_sandbox.go`, ein Unit-Test, Spec, Handbuch, CHANGELOG, Planung, Review. Die Container-Seite ist stattdessen im Praxislauf unten gemessen |
| `make verify-depguard` | nicht gelaufen, begründet | keine Importregel geändert (der Test importiert `application` und `port/driven` im Paket `application_test`; im Lint abgedeckt) |
| `make govulncheck`, `make image-scan` | nicht gelaufen | kein Security-/Release-Pfad berührt |

## DoD-Abgleich (Slice §2, Zeile für Zeile)

| # | DoD-Zeile | Befund | Evidence |
|---|---|---|---|
| 1 | Code und Tests (Lösung (b)) | **erfüllt** | `checkSandboxRuntime` in `internal/hexagon/application/doctor_sandbox.go`: `if s.inContainer()` unmittelbar vor `s.fs.Exists("/dev/fuse")`, nach Config-Zweigen und AppArmor-Zweig; Rückgabe `SeverityInfo`, kein `DegradationFallback`; Text nennt „inside a container“ und `devcontainer.sandbox.fallback.vfs`, Hint nennt den Container-Host. `inContainer()` ist bei `runtime == nil` falsch (Host-Verhalten bleibt). Tests: `TestDoctor_SandboxRuntime_InContainer` (Container × `warn`/`fail` × Gerät da/fehlt, Konfig-Warnung, `nestedRuntime: docker`, AppArmor, `nil`-Runtime, `&fakeRuntimeEnv{}`), `TestDoctor_SandboxRuntime` (neue Host-Zeile „fuse available, fail policy“), `TestDoctor_SandboxRuntime_FuseMissingHint` und `TestDoctor_SandboxFallbackFlagIsOnTheItem` unverändert grün (Host, `nil`-Runtime). Testmatrix des Plans vollständig abgedeckt (acht Zeilen) |
| 2 | Spec-Nachzug | **erfüllt** | `LH-FA-DEV-007.a` (Absatz nach der Degradationstabelle), `.b` (Klammer), Registry-Zeile `devcontainer.sandbox.runtime` („nicht im Container“), §7-Zeile vom 2026-10-18; Bedingung „soweit Statusdatei lesbar“ steht in `.a`, ist im Info-Text und im Test (`status file is readable`) gepinnt; „Konfigurations- und AppArmor-Zweig unverändert“ (§7) pinnen die Testzeilen „still warns“/„still reports privileged“/„AppArmor sysctl still an error“. Kein Verweis abwärts (Slice/ADR) in der Spec; `lastenheft.md`, `architecture.md`, `SPEC-045` nicht im Diff | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| 3 | Nutzer-Doku | **erfüllt** | Handbuch an vier Stellen (Runtime-Image-Aufzählung, `doctor`-Aufzählung, `/dev/fuse`-Absatz im Sandbox-Abschnitt, Absatz nach der Degradationstabelle bzw. Fehlerbehebung), jeweils mit „soweit seine Statusdatei lesbar ist (sonst „nicht ermittelbar“)“, ohne `LH-*`/`ADR-*`/Slice-Kennung. CHANGELOG `[Unreleased]`: Vergleich mit `v0.10.0` nennt `warn` (Warnung, Exit 0) **und** `onUnavailable: fail` (Fehler, Exit 11) und das Ergebnis „info, exit code 0 under both values“; Regression-Hinweis auf „Linux host outside a container“ begrenzt |
| 4 | Gates | **erfüllt, mit benannter Auslassung** | `make gates`, `make doc-trace` (0 Waisen), `make docs-check` grün; `make test-docker` nicht gelaufen (Grund oben) |
| 5 | Mutationen (Reviewer) | **erfüllt** | Review-Report: neun Mutationen, alle getötet (Gate entfernt/invertiert/vor Config-/AppArmor-Zweig, `Info` → `OK`, `nil`-Runtime als Container u. a.). Vom Verifier **nicht** erneut gefahren (Review-Beleg gesehen, Testfälle gegen die Mutationsliste gelesen: jede Mutation hat eine zuständige Zeile) |
| 6 | Review durchgeführt | **erfüllt** | `docs/reviews/2026-10-18-slice-v1-doctor-fuse-im-container.md`; getrennte Rolle | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| 7 | Verification Evidence mit Praxislauf | **erfüllt, Gegenprobe Host offen** | siehe Praxislauf; Linux-Host ohne Container: **nicht gemessen** |
| 8 | Closure-Notiz, Register | **offen (Closure, Planner)** | §7 des Slice steht auf „Wird bei der Closure gefüllt.“; keine Evidence-Dateien unter `observations/BEO-DOK/…/evidence/` für diesen Slice (Plan §3 nennt sie als „neu (Closure)“) |
| 9 | Risiken aus §6 tragen einen Ausgang | **offen (Closure)** | alle sechs Risiken stehen noch auf „offen bis …“; Vorschläge unten |
| 10 | Drei Paarungen | **offen (Closure)** | Anker/Folge-Slice/Register erst mit der Closure prüfbar |
| 11 | Release-Cut trägt den Slice als erfüllten Blocker | **offen (Closure, Planner)** | Cut-Slice §2/§4 nennen den Slice bereits im Futur-Wortlaut „liegt in `done/`“, obwohl er in `in-progress/` liegt; der Cut darf nicht starten, bevor die Closure das herstellt; Lead-Zeile gegen CHANGELOG erst dann abgleichbar |

## Praxislauf des Produkts (Docker-only)

Wegwerf-Images `ubverify-before` (aus `git archive e3d328e`, `make build IMAGE=ubverify-before`) und `ubverify-after` (HEAD, `make build IMAGE=ubverify-after`), zwei Wegwerf-Volumes mit `u-boot.yaml` (`profile: sandbox`, `nestedRuntime: podman`, `onUnavailable` je `warn`/`fail`), Aufruf `docker run --rm --user 0 -v <vol>:/work -w /work [--device /dev/fuse] <image> doctor`. Danach Images und Volumes entfernt (`docker rmi`, `docker volume rm`; Nachprüfung: kein `ubverify*` mehr in Images, Volumes, Containern; fremde Container unberührt).

| Fall | Vorher (`e3d328e`) | Nachher (HEAD) |
|---|---|---|
| `warn`, ohne `--device /dev/fuse` | `✗ devcontainer.sandbox.runtime … refuses to start …`, Summary `1 error`, **Exit 11** | `? devcontainer.sandbox.runtime /dev/fuse not checked: … inside a container …`, Hint „Run `u-boot doctor` on the container host to check /dev/fuse.“, Summary `0 error, 1 warn, 10 ok`, **Exit 0** |
| `warn`, mit `--device /dev/fuse` | `✓ … /dev/fuse is available`, Exit 0 | dieselbe Info-Zeile, **Exit 0** |
| `fail`, ohne Gerät | Fehler, **Exit 11** | Info, **Exit 0** |
| `fail`, mit Gerät | nicht gefahren | Info, **Exit 0** |
| `--quiet` (`fail`, ohne Gerät) | — | Eintrag nicht ausgegeben, Exit 0 |
| `--json` (`fail`, ohne Gerät) | — | Eintrag nicht im JSON (Info-Einträge erscheinen dort nicht), `"exitCode":0` |
| `--quiet --strict` | — | Exit 11 allein wegen der Warnung `compose.yaml.valid` (im Wegwerf-Projekt fehlt `compose.yaml`), nicht wegen des Info-Eintrags |

Im selben Lauf steht `devcontainer.sandbox.fallback.vfs` direkt darunter als „not determinable: … runs in a container without the status directory“; das bestätigt F-1 des Reviews und den nachgezogenen Wortlaut „soweit die Statusdatei lesbar ist“. Die Zeile `?` des neuen Eintrags und die Nicht-Zählung in der Summary entsprechen dem Muster der vier Host-Prüfungen.

**Nicht gemessen, ehrlich benannt:**

- **Gegenprobe Host (Linux ohne Container, ohne Modul `fuse`: Fehler, Exit 11):** nicht gefahren. Die Engine ist ein Colima-VM-Host; dort greift der Nicht-Linux-Zweig („cannot be determined from this host“). Die Host-Seite ist nur mit dem Fake gepinnt (`TestDoctor_SandboxRuntime`, `TestDoctor_SandboxFallbackFlagIsOnTheItem`, Test „not in container keeps the host behavior“ und „runtime nil …“).
- AppArmor-Wert `1` im echten Container (hätte den Sysctl der gemeinsamen VM verändert); nur mit dem Fake gepinnt.
- Podman-/Kubernetes-Marker der Container-Erkennung, Kata, gVisor: nicht gemessen; nur `/.dockerenv` (Docker).
- Devcontainer **mit** Statusdatei, `vfs`-Ausweisung im Container: nicht neu gemessen (unverändert aus dem Vorgänger-Slice; Pinning im Test `TestDoctor_SandboxFallback_InsideContainer`, nicht im Diff).

## Review-Findings am Code/Text (Bearbeitungsstand)

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (LOW, Ausweisung ohne Statusdatei-Bedingung) | behoben | Info-Text „as far as its status file is readable“ (Test gepinnt); `LH-FA-DEV-007.a`; Handbuch (zwei Stellen); CHANGELOG | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-2 (LOW, CHANGELOG nur `warn`) | behoben | CHANGELOG nennt `onUnavailable: fail` (Fehler, Exit 11 in `v0.10.0`) |
| F-3 (LOW, Kommentar `inContainer`) | behoben, **Restlücke (INFO)** | Kommentar nennt `checkSandboxRuntime`; ein weiterer Aufrufer, `checkSandboxFallbacks` (über `fallbackReader{inContainer: s.inContainer()}`, `doctor_sandbox.go`), steht weiter nicht im Kommentar. Der Kommentar nennt damit zwei von drei Nutzern; kein Fehler, aber dieselbe Klasse wie F-3. Mehr: der Satz „the four affected checks“ im Kommentar zu `hostHintSkippedInContainer` ist von der Aufgabe nicht berührt und stimmt weiter |
| F-4..F-6 (INFO) | in Slice §3 „Hinweise aus dem Review“ aufgenommen | vier Spiegelstriche (Marker-Erkennung, Hint auf Colima, AppArmor nicht gemessen, Bedingung Statusdatei) |
| F-7 (INFO, `open/README.md`) | behoben | Zelle nennt „Lösung (b) entschieden“ |

## Plan-Tabelle §3 gegen `git diff --stat e3d328e..HEAD` (AGENTS.md §3.10)

| Diff-Datei | Plan-Zeile | Abgleich |
|---|---|---|
| `internal/hexagon/application/doctor_sandbox.go` | update, `checkSandboxRuntime`, Kommentar `checkIDSandboxRuntime` | stimmt (Gate, Text, Hint, beide Kommentare) |
| `internal/hexagon/application/doctor_sandbox_test.go` | update, `TestDoctor_SandboxRuntime_InContainer`, Zeile „fuse available, fail policy“ | stimmt |
| **`internal/hexagon/application/doctor.go`** | „**prüfen, unverändert**“ | **Diskrepanz:** die Datei steht im Diff (Doc-Kommentar von `inContainer`, Review F-3); die Zeile trägt den alten Vermerk und nennt die Änderung nicht. Der Plan nennt die Datei als gelesen, ohne die Berichtigung — Rückfall der Klasse `plan-aussage-ueber-code` (Dateiliste nicht gegen den Diff nach der letzten Fix-Runde abgeglichen) |
| `spec/spezifikation.md` | update | stimmt (`.a`, `.b`, Registry, §7) |
| `docs/user/benutzerhandbuch.md` | update | stimmt (vier Stellen; Plan nennt (1)–(4), Diff trägt vier Hunks) |
| `CHANGELOG.md` | update | stimmt |
| `docs/plan/planning/in-progress/roadmap.md` | update (Planner), Ruhe-Marker | stimmt |
| `docs/plan/planning/open/README.md` | update (Planner) | stimmt (im Implementer-Diff trotz „Planner“, harmlos) |
| `docs/plan/planning/in-progress/slice-…md` | — | der Plan selbst, keine Zeile nötig |
| `docs/reviews/2026-10-18-slice-v1-doctor-fuse-im-container.md` | neu (nicht im Implementer-Diff) | liegt im Diff (Review-Commit `f544d13`), der Plan nennt den Report mit `<Datum>`; ohne Befund | <!-- d-check:ignore (Kennungen als Zitat im Report) -->

Zeilen „unverändert“ ohne Diff: `fakes_test.go`, `doctor_test.go`, `runtime.go`, CLI-Tests/Fixtures, `sandbox-init.sh.tmpl`, `devcontainer_sandbox.go`, `lastenheft.md`, `architecture.md`, READMEs — nicht im Diff, stimmen.

## Traceability

| ID | Beleg |
|---|---|
| `LH-FA-DEV-007` | `TestDoctor_SandboxRuntime`, `…_InContainer`, `…_FuseMissingHint`, `TestDoctor_SandboxFallbackFlagIsOnTheItem`; Praxislauf Vorher/Nachher | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `LH-FA-CLI-006` | Exit 11 → 0 im Container (Praxislauf); Host-Zweig Exit 11 mit Fake gepinnt | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `LH-FA-DEV-007.a`/`.b` | Spec-Diff, Info-Text-Pin (`status file is readable`) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `SPEC-045` | gelesen, unverändert wahr (Praxislauf: „not determinable“ im Runtime-Image) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| Commit-Hook | alle acht Commits nennen Slice-Kennung und `LH-FA-DEV-007` | <!-- d-check:ignore (Kennungen als Zitat im Report) -->

## Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle; kein Gate gelockert (Coverage 91,00 % über Schwelle 90 %).

## Auflagen für die Closure (§6/§7, Planner)

1. **§7 Closure-Notiz** mit Steering-Loop-Eintrag; Finding-Klassen des Reviews eintragen. Zu entscheiden: ob die Doku-Diskrepanz `doctor.go` („prüfen, unverändert“ trotz Kommentar-Änderung) und die Restlücke des `inContainer`-Kommentars (dritter Nutzer `checkSandboxFallbacks`) als weiterer Vorgang zu `plan-aussage-ueber-code` zählen.
2. **Plan-Tabelle berichtigen:** Zeile `doctor.go` auf „update (Kommentar von `inContainer`)“.
3. **Register:** je eine Evidence-Datei unter `observations/BEO-DOK/plan-aussage-ueber-code/evidence/` (Container-Bedingung des Plans nicht gelesen; zusätzlich der `doctor.go`-Abgleich) und `observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/evidence/` (Spec-Satz „soweit der Host es bestimmen kann“ ohne Test „Host gegen Container“); siehe [`plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`spec-bedingung-ohne-code-entsprechung`](../plan/planning/observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md); `state.md` nur anfassen, wenn sich der Stand ändert.
4. **Risiken §6, Ausgänge:** Sichtbarkeit ohne Statusdatei — eingetreten als benannter Zustand („nicht ermittelbar“, im Praxislauf gesehen), Abhilfe `generate devcontainer`; Info verdeckt Host-Fehler — weiter offen (Grenze); AppArmor im Container — E2 entschieden, Grenze: Wert `1` nicht gemessen; Marker-Erkennung — weiter offen; Test nur mit Fake — **entfallen für Docker** (Praxislauf maß Runtime-Image), offen für Podman/Kubernetes/Linux-Host; Verzögerung des Cuts — eingetreten (Slice war Blocker); Wiederholung des Plan-Fehlers — eingetreten in kleinem Maß (`doctor.go`-Zeile, Kommentar-Nutzer), sonst entfallen (Review las alle Prüfungen in `Check`).
5. **Roadmap-Marker nach dem Move:** der Verweis in `## Offene Wellen` zeigt auf `slice-v1-doctor-fuse-im-container.md` im Ruheort `in-progress/`; nach `make slice-mv … TO=done` Marker auf den Ruhe-Zustand zurücksetzen und die Zeile der Welle `v0.11.0` prüfen (Move-Commit rein, Marker eigener Commit, im selben Push; Inhalt vor dem Move, AGENTS.md §3.3). <!-- d-check:ignore (Kennungen als Zitat im Report) -->
6. **Cut-Start-Trigger:** `slice-v1-release-cut-v0.11.0` §4 und §2 sagen „liegt in `done/`“ im Futur; erst nach der Closure wahr. Die Lead-Zeile zu `doctor` und `/dev/fuse` gegen den CHANGELOG-Wortlaut halten (Linux-Host außerhalb eines Containers ist Fehler; `warn` und `fail` gegenüber `v0.10.0`). <!-- d-check:ignore (Kennungen als Zitat im Report) -->
7. **`open/README.md`:** Zeile des Slice zeigt auf `in-progress/`; nach dem Move auf `done/` und Beschreibung auf „geliefert“ umstellen (Zustand und Anker, keine Chronik).
8. **Nicht gemessen** (§7 benennen): Linux-Host ohne Container (Gegenprobe), Wert `1` im Container, Podman-/Kubernetes-Marker, Kata/gVisor.

## Verdikt

Der Code folgt dem Plan (E1 (b), E2 nur `/dev/fuse`), alle Gates grün, `doc-trace` 0 Waisen, `doc-immutable` OK, Praxislauf im Runtime-Image bestätigt Exit 11 → Exit 0 (mit und ohne Gerät, `warn` und `fail`, `--quiet`/`--json`). DoD-Zeilen 1–7 erfüllt (Zeile 7 mit benannter Auslassung der Host-Gegenprobe); die Zeilen 8–11 sind Closure-Arbeit und offen. Eine Diskrepanz (Plan-Zeile `doctor.go`), eine Kommentar-Restlücke (INFO), keine Merge-blockierende Abweichung.
