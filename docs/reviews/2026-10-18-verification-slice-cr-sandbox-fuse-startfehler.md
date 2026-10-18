# Verification Evidence: slice-cr-sandbox-fuse-startfehler — 2026-10-18 <!-- d-check:ignore (Slice-Kennung als Zitat) -->

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der CR richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach [`harness/verification.md`](../../harness/verification.md).

**Gegenstand:** CR-Slice `slice-cr-sandbox-fuse-startfehler` (`docs/plan/planning/in-progress/`), Commit-Range `17e4b90..1204c1e` (10 Commits: `c1ef60e` Lastenheft, `4c16901` Spezifikation, `a28a9f5` Code und Tests, `e4e5d0d` Handbuch/CHANGELOG, `b0919ad` Roadmap, `702335a`/`2ff6db8`/`f839d28`/`1204c1e` Nachzüge, `175d02d` Review-Report). Maßstab: Slice §1 (freigegebener Wortlaut, Entscheidungen 2026-10-17/18), §2 (DoD), §3, §6. Review-Report: [`2026-10-18-slice-cr-sandbox-fuse-startfehler.md`](2026-10-18-slice-cr-sandbox-fuse-startfehler.md). <!-- d-check:ignore (Slice-Kennung und Pfade als Zitat) -->

## Verification Evidence

Scope:
- Slice: `slice-cr-sandbox-fuse-startfehler` <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- IDs: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Änderung), [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten), [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md) (Messung M13/M14)
- Artefakte: `spec/lastenheft.md`, `spec/spezifikation.md`, `internal/hexagon/application/doctor_sandbox.go` (`checkSandboxRuntime`), `doctor_sandbox_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go`, `docs/user/benutzerhandbuch.md`, `CHANGELOG.md`, `README.md`, `README.de.md`, `roadmap.md`. <!-- d-check:ignore (Pfade als Zitat) -->

DoD-Abgleich (Slice §2):
- [x] Lastenheft 0.13.0 — Evidence: Versionsfeld `0.13.0`; Stelle 1 (Punkt „Bei `podman` startet …“, beide Devices) und Stelle 2 (Absatz Degradation) im Lastenheft; Stelle 2 per Textvergleich (Slice §1 „Neu (vollständig)“ gegen `spec/lastenheft.md` Zeile 464, Links normalisiert) **zeichengleich**, Stelle 1 und §16-Zeile (Datum 2026-10-18, Messstand, „Vereinbarung mit dem Projektinhaber“) stimmen mit dem Vorschlag überein; der Lastenheft-Diff umfasst sieben Zeilen, keine weitere Zusage berührt (`LH-FA-CLI-005` Ausnahme 2 unverändert, wie §1 begründet); kein Verweis nach unten. <!-- d-check:ignore (Pfade als Zitat) -->
- [x] Nachzug Spezifikation — Evidence: `.a` Tabellenzeile „`/dev/fuse` im Container nicht verfügbar (Container ohne `--device=/dev/fuse` gestartet)“ plus neuer Absatz (Host-Fall kein Degradationsschritt, `doctor` Fehler, Exit-Code `11`, unabhängig von `onUnavailable`, ohne Fallback-Kennzeichen); `.b` beide Devices, `doctor` weist `/dev/fuse` nur aus, soweit der Host es bestimmen kann; Code-Registry `devcontainer.sandbox.runtime`; §7-Zeile vom 2026-10-18. Die „lokale Engine“-Bedingung (Review F-1) ist entfernt (`f839d28`). Die übrigen Sätze von `.a` und `SPEC-043` unverändert, wie §3 vorsieht. <!-- d-check:ignore (Kennung als Zitat) -->
- [x] Nachzug Code — Evidence: `checkSandboxRuntime` liefert bei fehlendem `/dev/fuse` auf dem Linux-Host immer `SeverityError`, kein `DegradationFallback`, Text „… the container engine refuses to start the sandbox container with `--device=/dev/fuse` (LH-FA-DEV-007, exit 11).“, Hint `modprobe fuse` / `nestedRuntime` auf `none`; Nicht-Linux-Zweig OK mit Engine-Hinweis; der Zweig `onUnavailablePolicy` ist dort entfallen (die Funktion bleibt in `devcontainer_sandbox.go` genutzt, `make lint` grün). Funktionskommentar nennt den neuen Befund. Startscript und `runArgs` unverändert (Diff berührt kein Template; Golden Cases unberührt). <!-- d-check:ignore (Kennung als Zitat) -->
- [x] Nachzug Tests — Evidence: `TestDoctor_SandboxRuntime` (beide Policies `SeverityError`/„refuses to start“, neuer Nicht-Linux-Fall), neu `TestDoctor_SandboxRuntime_FuseMissingHint` (pinnt den Hint, Review F-3), `TestDoctor_SandboxFallbackFlagIsOnTheItem` (`fallback=false`); e2e `TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart` (Kontrolllauf mit echten `runArgs`, dann `--device=/dev/fuse-nonexistent`, erwartet Exit ≠ 0 und `error gathering device information`). Format (Review F-4) berichtigt (`2ff6db8`).
- [x] Doku-Update — Evidence: Handbuch (Absatz `/dev/net/tun` auf beide Devices; `vfs` nur für den ohne Device gestarteten Container; Tabelle „fehlt im Container“; `doctor`-Fehler, Exit-Code 11; Fehlerbild) ohne Kennungen und ohne Chronik; CHANGELOG `[Unreleased]`/`Changed` mit Messstand, **Regression für Skripte ohne `--strict`** benannt und „generated files do not change“ (stimmt mit dem Diff); READMEs EN/DE nennen beide Devices (Review F-2, `f839d28`).
- [x] `make gates` grün; `make test-docker` grün (Colima); `make doc-trace` ohne Waise — siehe Sensors (selbst gefahren).
- [x] Review durchgeführt — Evidence: `docs/reviews/2026-10-18-slice-cr-sandbox-fuse-startfehler.md` (0 HIGH, 0 MEDIUM, 3 LOW, 4 INFO); F-1 (Spec), F-2 (READMEs), F-3 (Hint-Test), F-4 (Format) am Code/Text bestätigt behoben; F-5 bewusst offen (Umsetzungs-Slice), F-6/F-7 ohne Handlung (F-7 durch diese Verification geschlossen). <!-- d-check:ignore (Pfad als Zitat) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen: Slice §7 trägt „Wird bei der Closure gefüllt.“ (Auflage 1).
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7 — offen, mit der Closure-Notiz (Auflage 1).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: alle sechs Risiken tragen noch „offen“ (Auflage 2).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — Folge-Slice getragen (`slice-v1-sandbox-fallback-host-ausweisung` liegt in `next/`); Anker und Register entstehen mit der Closure-Notiz (Auflage 1). <!-- d-check:ignore (Kennung als Zitat) -->

Plan-Tabelle §3 gegen `git diff --stat 17e4b90..HEAD` (`AGENTS.md` §3.10): zwölf Dateien im Diff, jede im Plan: `spec/lastenheft.md`, `spec/spezifikation.md`, `doctor_sandbox.go`, `doctor_sandbox_test.go`, e2e-Test, Handbuch, CHANGELOG, `README.md`, `README.de.md`, `roadmap.md` (Ruhe-Marker entfernt), Review-Report, der Slice selbst (die Plan-Tabelle). Jede „update“-Zeile hat einen Diff; jede „geprüft, unverändert“-Zeile (Template, `open/README.md`, `cli_test.go` & Co., `architecture.md`, `.golangci.yml`, `.d-check.yml`, Carveouts) hat keinen. Keine Diskrepanz. <!-- d-check:ignore (Dateinamen als Zitat) -->

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | selbst gefahren, Exit 0; `coverage-gate: OK — coverage 91.20% meets threshold 90%` (lint, test, docs-check, baseline-verify Teil von `gates`) |
| `make doc-trace` | pass | selbst gefahren, Exit 0; „83 Anforderung(en), 0 Waise(n).“ |
| `make doc-immutable RANGE="17e4b90..HEAD"` | pass | selbst gefahren, Exit 0; 10 Commits aufgelöst, „0 Befund(e)“ (kein ADR im Diff) |
| `make test-docker` | pass | selbst gefahren (Colima), Exit 0; alle Pakete `ok`, `internal/e2e` 488 s |
| gezielt `go test -tags docker -v -count=1 -run 'TestE2E_SandboxDevcontainer_(MissingHostFuseDeviceRefusesStart\|FailPolicyWithoutFuse\|NestedPodman$)' ./internal/e2e/` im Image `u-boot:test-docker-tools` mit Mount-Muster von `repo.mk` | pass, kein Skip | `--- PASS` für `NestedPodman` (26,8 s), `MissingHostFuseDeviceRefusesStart` (20,8 s), `FailPolicyWithoutFuse` (20,0 s); kein `SKIP` im Log |
| Praxislauf `doctor` (Docker-only-Binary `make build` → `u-boot:latest`, Wegwerf-Verzeichnis, danach gelöscht) | pass | siehe unten |
| `make docs-check` | pass | Teil von `make gates`; nach Anlegen dieser Datei erneut gefahren (siehe Abschluss) |

Praxislauf `doctor` (Profil `sandbox`, `nestedRuntime: podman`):
- Der Lauf im Runtime-Image überspringt die Host-Prüfungen (`docker.*`, `git.installed`: „skipped — u-boot is running inside a container“), **nicht** aber `devcontainer.sandbox.runtime`. Der Container ist ein Linux-System ohne `/dev/fuse`; damit ist der Fehlerzweig praktisch reproduzierbar (er prüft das Device des Containers, nicht des Docker-Hosts; das Verhalten auf einem echten Linux-Host belegt der Unit-Test).
- `onUnavailable: warn` und `fail`: beide `✗ devcontainer.sandbox.runtime  /dev/fuse is not available on this host; the container engine refuses to start the sandbox container with --device=/dev/fuse (LH-FA-DEV-007, exit 11).` mit Hint `modprobe fuse` / `nestedRuntime` `none`; Exit-Code des Befehls jeweils `11`. Der Befund ist unabhängig von der Policy, wie Lastenheft und Spezifikation es sagen. <!-- d-check:ignore (Kennung als Zitat) -->
- Gegenprobe mit `docker run --device /dev/fuse`: `✓ … /dev/fuse is available for nested Podman.`, Exit-Code `0`.
- Nicht-Linux-Zweig (macOS-Host, Item OK mit Hinweis auf die Engine-Ablehnung): **nicht praktisch geprüft** — es gibt kein macOS-Binary im Docker-only-Werkzeug; belegt durch `TestDoctor_SandboxRuntime` („non-linux host cannot tell“, „non-linux host without fuse names the engine refusal“).

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | Lastenheft 0.13.0 (Stelle 1/2); `TestDoctor_SandboxRuntime`, `TestDoctor_SandboxRuntime_FuseMissingHint`, `TestDoctor_SandboxFallbackFlagIsOnTheItem`; e2e `TestE2E_SandboxDevcontainer_MissingHostFuseDeviceRefusesStart`, `…FailPolicyWithoutFuse`, `…NestedPodman`; Praxislauf `doctor` |
| [`LH-FA-DEV-007.a`](../../spec/spezifikation.md#lh-fa-dev-007a--degradation-und-strenge-der-sandbox-fähigkeiten) / [`.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman) | Spezifikations-Diff (Tabelle, Absatz, Ausstattung, Registry, §7); `make docs-check` |
| [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Exit-Code `11` des Fehler-Items (Praxislauf, Unit-Test über die zentrale Klassifikation) |
| [`ADR-0022`](../plan/adr/0022-fallback-ausweisung-statusdateien-im-container.md) | M13 am Test `…MissingHostFuseDeviceRefusesStart` erneut gemessen (Exit ≠ 0, `error gathering device information`); `make doc-immutable` ohne Befund |

Carveouts:
- Neu: none
- Gelöst: none
- Unverändert: alle

Nicht ausgeführt:
- `make verify-depguard` — Importregeln nicht berührt (kein Import- oder `.golangci.yml`-Diff).
- `make ci` / `make image-scan` / `make govulncheck` — kein Security-/Release-Pfad berührt.
- Replay/Golden — kein Generator-Output geändert (`runArgs` und Templates unberührt; `generate_sandbox_test.go` pinnt weiter sechs Optionen).
- Praxislauf des Nicht-Linux-Zweigs auf dem macOS-Host — Grund siehe oben, durch Unit-Test belegt.
- Linux ohne VM, amd64, Docker Desktop, Podman als Host-Engine — nicht gemessen (Messstand Colima, Docker 29.5.2, arm64 steht im Lastenheft und im CHANGELOG).

## Diskrepanzen

Keine zwischen Plan, Lastenheft, Spezifikation, Code, Tests, Handbuch, CHANGELOG und READMEs. Plan-Aussagen über den Code (Zweig, Texte, Tests) am Code gelesen und bestätigt.

## Auflagen für die Closure

1. **Closure-Notiz §7** mit Steering-Loop-Lerneintrag und Register-Vermerk („keine Beobachtung“ oder Fortschreibung unter `observations/`). Kandidaten für den Lerneintrag aus dem Review: Plan-Tabelle gegen die Suche nach dem *Wort* (`fuse`) statt gegen die *Aussage* (Voraussetzung `/dev/net/tun`) geführt (README, F-2); ein Hint-Text ohne Test.
2. **Risiko-Ausgänge §6** setzen, alle sechs derzeit „offen“: (a) Einstufung `doctor` Warnung → Fehler: eingetreten und im CHANGELOG benannt (Regression ohne `--strict`); (b) Falsch-positiv bei entferntem `DOCKER_HOST`: weiter offen (Review F-1 bleibt als Grenze, nicht im Vertrag zugesagt) oder als eingetreten-akzeptiert; (c) Messung nur Colima: weiter offen; (d) `vfs`-Zweig als Fallback mit enger Ursache: weiter offen, getragen vom Umsetzungs-Slice; (e) Überschneidung mit dem Umsetzungs-Slice: entfallen (Reihenfolge CR → Umsetzung, `SPEC-043` und die übrigen Sätze von `.a` unberührt); (f) Spec ohne Verweise nach unten: entfallen (`make docs-check` 0 Befunde, Reviewer ohne Befund). <!-- d-check:ignore (Kennung als Zitat) -->
3. **Roadmap-Ruhe-Marker** nach dem Move nach `done/`: „Nichts in Arbeit.“ in `roadmap.md` wieder einsetzen (der Marker wurde in `b0919ad` entfernt, solange der Slice in `in-progress/` liegt; d-check `planning`). `make slice-mv` zieht die `done/`-Verweise nach, den Marker nicht. <!-- d-check:ignore (Pfad als Zitat) -->
4. **Reihenfolge nach AGENTS.md §3.3:** erst Inhalt (DoD-Häkchen, Closure-Notiz), dann der reine `git mv`; die Zeile im Planning-README (`open/README.md`) und die Verweise aus dem Umsetzungs-Slice danach auf Auflösbarkeit prüfen. Der Umsetzungs-Slice startet erst mit dem Slice in `done/`; er trägt den offenen Nachzug von `SPEC-043` Ausnahme 2 und der Ausnahme 2 in `LH-FA-CLI-005` (Review F-5, host-seitig ist `vfs` nicht mehr ermittelbar). <!-- d-check:ignore (Pfad und Kennung als Zitat) -->
5. **Register-Paarung:** keine neue Beobachtung angefallen, die ein Register-Eintrag verlangt; falls der Lerneintrag aus Auflage 1 die dritte Beobachtung zu `BEO-DOK/plan-aussage-ueber-code` ist, gilt die 3×-Schwelle erneut (Verkörperung prüfen). <!-- d-check:ignore (Beobachtungs-Kennung als Zitat) -->

## Verdikt

**DoD:** fünf von neun Punkten erfüllt (Lastenheft, Nachzug Spezifikation/Code/Tests, Doku, Gates/Sensors, Review); offen sind ausschließlich die Closure-Punkte (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen), die laut Slice §5 bei der Closure entstehen. Keine DoD-Verletzung im Gelieferten, keine Verifier-only-Diskrepanz.

**Übergabe:** Bericht an den Planner; die Closure ist nach Auflagen 1 und 2 möglich, Auflage 3 gehört in den Move-Schritt.
