# Review-Report: slice-v1-devcontainer-instanz-test-leicht — 2026-11-01

**Review-Art:** Code — Code-Review gegen Plan, Entscheidungen und Konventionen (Weg A).

**Gegenstand:** Commits `73b42e21` und `54c001db` samt den vier Move- und Nachzug-Commits davor (`622a626b`, `218b68ee`, `cf3e2ad1`, `456ce5fd`), Basis `c370280c`. Slice [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md). Sieben Dateien, 91 Einfügungen, 225 Löschungen (`git diff --stat c370280c..HEAD`).

**Skill:** `.harness/skills/reviewer.md` @ `54c001db`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-01

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan [`slice-v1-devcontainer-instanz-test-leicht`](../plan/planning/done/slice-v1-devcontainer-instanz-test-leicht.md) (§1 entschiedene Fragen, §2 DoD, §3 Plan, §6 Risiken); Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/next/slice-cr-devcontainer-instanz-zusage-getestet.md)
- [`ADR-0038`](../plan/adr/0038-devcontainer-instanz-zweites-id-label.md)
- [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`SPEC-049`](../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz Instanzen)
- Messung [`slice-v1-devcontainer-instanz-messung`](../messungen/slice-v1-devcontainer-instanz-messung/README.md)
- `AGENTS.md` (Hard Rules, besonders §3.6, §3.7, §3.9, §3.10)

**Sensors:** Kein Go-Toolchain im Prüfkontext (`go` nicht im `PATH`), daher ist die Übersetzbarkeit des Tags `docker` aus dem Quelltext gelesen, nicht gebaut. `make test-docker` lief nicht (Auftrag); `make verify-depguard` nicht nötig (keine Importregeln berührt). `make gates` habe ich nicht neu gefahren (Review ohne Produktänderung; Nachweis liegt beim Implementer/Verifier). `make docs-check` siehe Verdikt.

---

## Übersetzbarkeit des Tests (Referenz für Referenz)

Datei `internal/e2e/devcontainer_instance_docker_test.go`, Build-Tag `docker`, Paket `e2e_test`; alle Helfer in `internal/e2e/*.go` mit demselben Tag oder ohne.

| Verwendet (Zeile in der Test-Datei) | Fundstelle | Signatur passt |
|---|---|---|
| `requireDocker(t)` (118) | `sandbox_status_docker_test.go:41` | ja |
| `dcGenerateProject(ctx, t, true, nil)` (119) | `devcontainer_start_docker_test.go:195` `(ctx, t, bool, map[string]string) string` | ja (`nil`-Map zulässig, `range` über nil) |
| `dcService()` (120) | `devcontainer_start_docker_test.go:285` | ja; `svc.Up`/`svc.Down` wie im unveränderten Standardprofil-Test |
| `ubootBinary(t)` (121) | `devcontainer_start_docker_test.go:714` `string` | ja |
| `dcUpInstance` (124, 128, 183, 194, 227, 232) | in der Datei (37), `driving.DevcontainerUseCase` | ja |
| `dcRequireDistinct(t, what, []string)` (137, 163) | in der Datei (47) | ja |
| `dcInspect(t, id)` -> `dcInspectDoc` mit `.State.Running`, `.Config.Labels`, `.Mounts`, `volumeAt(string) string` (140–161, 236, 249) | `devcontainer_start_docker_test.go:75, 96, 142` | ja |
| `dcFolderLbl`, `dcConfigLbl`, `dcTestBudget` (116, 144, 148) | `devcontainer_start_docker_test.go:47–49` | ja |
| `dcInstanceLabel`, `dcWorkspaceDest` (147, 154, 161) | in der Datei (29, 32) | ja |
| `dcExec(ctx, t, id, user, script) (string, int)` (168, 174, 240) | `devcontainer_start_docker_test.go:319` | ja (zwei Rückgaben je Aufruf) |
| `runBinary(ctx, t, bin, dir, stdin, args...) (string, string, int)` (187, 190, 262) | `devcontainer_shell_docker_test.go:222` | ja |
| `dcRunning(t, id) bool`, `dcExists(id) bool` (208, 223, 233, 243) | in der Datei (74), `devcontainer_start_docker_test.go:106` | ja |
| `dcMustDocker(t, args...) string` (250, 67) | `devcontainer_start_docker_test.go:65` | ja |
| `dcDocker(args...) (string, error)` (253) | `devcontainer_start_docker_test.go:54` | ja (beide Rückgaben mit `_, _ =` verworfen) |
| `dcFolderContainers(t, dir) map[string]bool` (198, 256, 268) | in der Datei (64) | ja (`len`, `range`, Index) |
| `dcWatchNoHelper` (261, 296), `dcWatchHelpers`, `dcHelperNames` | in der Datei (87); `devcontainer_shell_docker_test.go:98, 68` | ja |
| `dcProjectDir` (292), `dcMarkerPath` (168, 174, 240) | `devcontainer_start_docker_test.go:158`; in der Datei (81) | ja |
| `driving.DevcontainerUpRequest{BaseDir, Instance, Recreate}`, `DevcontainerDownRequest{BaseDir, Instance}`, `ErrDevcontainerOnlyInstances` | `internal/hexagon/port/driving/devcontainer.go` (Zeilen 20–24, 70, 105, 140) | ja |

Importe: `context`, `errors` (277), `os` (342), `path/filepath` (342), `strings`, `testing`, `time` (290, 319, `time.Minute`), `driving` — alle benutzt; `fmt`, `strconv`, `syscall` sind entfernt und im Rest der Datei nicht mehr referenziert (Suche ohne Treffer). Ungenutzte Variablen: keine (`n`, `id`, `doc`, `markers`, `baseVolume`, `repeat`, `none`, `stopped`, `restarted`, `recreated` werden gelesen). Deklarationen: `down, err :=` (203) deklariert `err` neu im Funktionsscope, `none, err :=` und `stopped, err :=` (213, 218) führen je ein neues Symbol ein, kein Redeclare-Fehler; `doc` in der Schleife (140) und `doc :=` (236) liegen in verschiedenen Scopes. Doppelte Deklarationen in der Test-Datei: keine (je Name ein `func`/`const`).

Entfernte Symbole `dcInstanceCount`, `dcInstanceMinMemAvailable`, `dcInstanceMinSwapFree`, `dcInstanceMinFreeBytes`, `dcInstanceMemory`, `dcInstanceNanoCPU`, `dcDockerDest`, `dcOllamaDest`, `dcMeminfoBytes`, `dcRequireInstanceHost`, `dcLimitRunArgs`, `dcMarkerPaths`: im ganzen Baum (außer `.claude/worktrees`) in keiner Go-Datei mehr referenziert; Treffer nur in `docs/` (Chronik, siehe F-3). Ergebnis: Ich finde keinen Übersetzungsfehler; Gewissheit bringt nur ein Lauf mit Tag `docker` (F-6).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/next/slice-cr-devcontainer-instanz-zusage-getestet.md) trägt weiter, die Suite belege „Container ohne Instanz plus drei Instanzen“ (Abschnitt „Nach … belegt die Suite“, Variante 1 „zwei und drei gleichzeitig laufende Instanzen“, Risiko „Die Suite belegt ‚zwei‘ nicht“) und der Test werde „ein manuelles Messskript“ (Anlass, Ziel, „Ausdrücklich NICHT“). Entschieden und gebaut sind zwei Instanzen und kein Messskript; der Diff hat in dieser Datei nur die Links nachgezogen. Wer den Wortlaut für das Lastenheft aus Variante 1 freigibt, sagt drei Instanzen zu, die keine Suite fährt. | `AGENTS.md` §3.10 (Aussage im Plan über Code/Test nicht gelesen); [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Entwurf zur Zusage (Verzeichnis `open/`) · „Container ohne Instanz plus drei Instanzen, ein Workspace-Volume je Instanz“ | ja — `grep` nach „drei Instanzen“ und „Messskript“ in `open/` | Folge-Entwurf trägt überholte Test-Annahme |
| F-2 | LOW | Die Plan-Tabelle (§3) führt die Datei des Entwurfs zur Zusage nicht, obwohl der Diff sie ändert (vier Link-Zeilen); der Abgleich gegen `git diff --stat` ist für diese Zeile offen. Dieselbe Tabelle führt `zusagen-und-messstaende.md` und `CHANGELOG.md` zweimal. Der DoD-Punkt „Wächter und Doku-Verweise (b)“ verlangt weiter „den leichten Test und das Messskript“, die entschiedene Frage 3 hat kein Skript. | `AGENTS.md` §3.10 | Slice-Datei (Verzeichnis `in-progress/`) · „nennt den leichten Test und das Messskript statt des schweren Tests“ | Review-only | Plan-Tabelle gegen Diff nicht abgeglichen |
| F-3 | LOW | Das Zustandsfeld im Roadmap-Zeiger lautet „(in Arbeit, anderer Lauf als der vorige)“; „anderer Lauf als der vorige“ ist Lauf-Chronik und löst auf nichts auf. Zudem liegen zwei Slices in `in-progress/`, während §4 des Slice WIP-Limit 1 nennt („Start erst nach dessen Closure oder auf Anweisung des Projektinhabers“); die Anweisung ist weder in §7 noch im Roadmap-Zeiger als Beleg genannt. | `AGENTS.md` §3.7 (Zustandsfeld nennt Zustand und Beleg, nicht Chronik) | `docs/plan/planning/in-progress/roadmap.md` · „in Arbeit, anderer Lauf als der vorige“ | Review-only | Zustandsfeld mit Chronik; WIP-Anweisung ohne Beleg |
| F-4 | LOW | `dcInspectDoc.HostConfig` trägt weiter `Memory`, `MemorySwap`, `NanoCpus` samt Kommentar („are the limits of the container“); die einzigen Leser, der Grenzen-Rücklesetest, sind entfernt. Die Felder haben jetzt keinen Leser (Suche über `internal`: nur die Deklaration). | `AGENTS.md` §3.9 (b), Maintainability | `internal/e2e/devcontainer_start_docker_test.go` · „Memory, MemorySwap (bytes) and NanoCpus are the limits of the container.“ | ja — `make lint` bestätigt es nicht (Exporte, unmarshal-Felder) | Reste der entfernten Prüfung |
| F-5 | LOW | `docs/maintainer/zusagen-und-messstaende.md` nennt die Zusage („bis zu vier … `nestedRuntime: docker` und `ollama` …“) unverändert und danach, was der Test nicht läuft; der Satz „die Grenze von einer CPU und 3 GiB und vier Instanzen läuft die Suite nicht; sie stehen nur in der Messung“ passt, die Unit-Tests der Mount-Texte sind benannt. Gegenüber dem Lastenheft bleibt die Lücke (Zusage weiter als die Suite) offen und ist als Risiko 1 des Slice und dem Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/next/slice-cr-devcontainer-instanz-zusage-getestet.md) zugeordnet; der Bericht belegt nur, dass sie getragen, nicht dass sie geschlossen ist. Die Handbuch-Bedingung „soweit aktiviert“ (Docker-Speicher und Modelle) hat weiter nur den Text-Pin der Mounts, keinen Lauf. | `AGENTS.md` §3.10 (Bedingung in Spec/Handbuch ohne Lauf) | `docs/maintainer/zusagen-und-messstaende.md` · „läuft die Suite nicht; sie stehen nur in der Messung“ | Review-only | Zusage weiter als die Suite (bekannt, getragen) |
| F-6 | INFO | Der neue Test ist weder gelaufen noch gebaut. Zwei Annahmen sind erst durch den Lauf auf dem Runner belegt: (a) das Sandbox-Profil mit Standard-Einstellungen startet dort ohne die Einstellungen `egress.enabled`/`onUnavailable`, die jeder andere Sandbox-Test explizit setzt (`dcGenerateProject(ctx, t, true, nil)` ist sonst nirgends benutzt); (b) die Marker im Volume gehören `root` und sind per `-u root` lesbar. Zuständig: Verifier/CI (`integration-docker`). | Slice §6 Risiko 3 | `internal/e2e/devcontainer_instance_docker_test.go` · `dcGenerateProject(ctx, t, true, nil)` | ja — `make test-docker` (Auftrag: nicht jetzt) | Test ohne Lauf |
| F-7 | INFO | Läuft-auf-GitHub-Runner (Prüfpunkt 2) gelesen, ohne Befund: keine Hostressourcen-Prüfung mehr; drei Sandbox-Container ohne Grenze sind plausibel (Basis-Image `devcontainers/base:debian`, ohne `nestedRuntime`/`ollama`, nichts läuft darin außer dem Leerlaufprozess). Der Projektordner heißt `dcs<UnixNano>` (`dcSuffix`), die Volumes `…-workspace-${devcontainerId}` leiten sich aus den Labels (Ordner, Instanz) ab, die Images aus dem Ordnernamen; Containernamen vergibt die Dev-Container-CLI selbst. Eine Kollision mit anderen Tests oder parallelen Läufen (Namenskonflikt wie im früheren Messlauf) entsteht damit nicht. Aufräumen: `dcRegisterCleanup` entfernt Container (`rm -f -v`), Volumes und Images je Ordner; das Volume des vor dem Ende entfernten Containers ohne Instanz räumt der zusätzliche `t.Cleanup` (läuft vor dem Cleanup des Ordners, LIFO, ohne Konflikt). Zeitrahmen: Kontext `dcTestBudget` 20 Minuten gegen `go test -timeout 30m ./...` und `timeout-minutes: 30` des Jobs; das Budget des Pakets (alle e2e-Tests) liegt unverändert beim Job, nicht beim Test. | Maintainability | `.github/workflows/integration.yml` · `timeout-minutes: 30`; `internal/e2e/devcontainer_start_docker_test.go` · `func dcRegisterCleanup` | Review-only | — |
| F-8 | INFO | Wächter (Prüfpunkt 3), aus dem Test gelesen, tragfähig: (1) `up --instance` ohne Instanz-Label: `got != names[n-1]` schlägt an. (2) Gemeinsamer Workspace-Volume-Name: `dcRequireDistinct "workspace volume"` und die Paar-Probe `(code == 0) != (i == j)` schlagen an. (3) `up` ohne `--instance` wählt die neueste Instanz: `again.ContainerID != existing.ContainerID` und der Hostname der Shell. (4) `--recreate` entfernt das Volume: weil der Volume-Name aus den Labels folgt und gleich bleibt, fängt `volumeAt != workspaces[1]` das nicht; es trägt allein die Prüfung der Marker-Datei nach `--recreate` (`test -e`). Das Mutationsprotokoll des Implementers liegt nicht im Diff; der Verifier fährt die Mutationen selbst. Zuständig: Verifier. | Slice §2 DoD (a) | `internal/e2e/devcontainer_instance_docker_test.go` · „recreate of %s: the file is gone“ | ja — Verifier (Mutationen) | — |

## Zu den Prüfpunkten

1. **Trägt der Test die Aussagen?** Slice §1 und CHANGELOG (Zeile 16: „two instances and the container without an instance in the sandbox profile (default settings: no nested runtime, no feature `ollama`, no resource limit) and with three containers in the standard profile“) decken sich mit dem Test: Identität/Labels, Workspace-Volumes, Marker über alle geordneten Paare, Container-Wahl ohne `--instance`, `down`, `--recreate`, Exit 10 ohne Container ohne Instanz. Was die Suite nicht läuft, benennt `zusagen-und-messstaende.md` (F-5) und der Kopfkommentar des Tests; der CHANGELOG nennt es nicht, er sagt nur, was gepinnt ist (korrekt, kein Befund).
2. **Runner:** F-7.
3. **Wächter:** F-8.
4. **Kommentare §3.7/§3.9:** Kopfkommentar und Docstring der Test-Datei beschreiben den Ist-Stand in Indikativ; kein „früher“, „entfernt“, „no longer“ (Suche in Diff und Datei ohne Treffer). Die verneinenden Teile („no nested runtime, no feature `ollama`, no resource limit …“) sind Abgrenzung, keine Chronik. Offene Reste: F-4.
5. **Plan-Tabelle gegen `git diff --stat`:** sieben Dateien im Diff. Geführt und im Diff: Test-Datei, `zusagen-und-messstaende.md`, `CHANGELOG.md`, `roadmap.md`, `open/README.md`, der Slice selbst. Nicht geführt: der Entwurf zur Zusage in `open/` (F-2). Tabelle-Zeilen ohne Diff: `BEO-E2E/…` (Closure, noch offen, zulässig), die „unverändert“-Zeilen.
6. **Handbuch/Spec/Lastenheft:** `docs/user/benutzerhandbuch.md`, `spec/`, `README*.md` sind im Diff unverändert; Handbuch Abschnitt 4.9 nennt keinen Test (Zeilen 888 ff.). „Unverändert“ trifft zu (keine Zahl der Instanzen, kein Test genannt).
7. **Alte Verweise auf den Testnamen `…FourInstancesAreIsolated`:** außerhalb von Chronik nur im Entwurf [`slice-cr-devcontainer-instanz-zusage-getestet`](../plan/planning/next/slice-cr-devcontainer-instanz-zusage-getestet.md) (als Beschreibung des früheren Belegs, F-1) und im Slice selbst (Beschreibung des Gegenstands, zulässig). Chronik (unverändert gelassen, laut Plan zulässig): der abgeschlossene Slice [`slice-v1-devcontainer-instanz-label`](../plan/planning/done/slice-v1-devcontainer-instanz-label.md), die zwei Reports zu ihm in `docs/reviews/`, die Evidence unter `observations/BEO-APP/…`. In `internal/`, `spec/`, `docs/user/`, `.github/`, `Makefile`: kein Treffer.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/e2e/devcontainer_instance_docker_test.go` (Übersetzbarkeit, Importe, Variablen, Deklarationen) | geprüft, ohne Befund (gelesen, nicht gebaut) |
| Entfernte Hilfsfunktionen und Konstanten in anderen `docker`-getaggten Dateien | geprüft, ohne Befund (keine Referenz) |
| `CHANGELOG.md` Zeile 16 | geprüft, ohne Befund |
| `docs/maintainer/zusagen-und-messstaende.md` Fundstelle | geprüft, ohne Befund außer F-5 (bekannte Lücke) |
| `docs/user/benutzerhandbuch.md`, `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md` | geprüft, unverändert, ohne Befund |
| Instanz-Messung unter `docs/messungen/` | geprüft, unverändert (Diff ohne Datei dort) |
| `.github/workflows/integration.yml`, `Makefile`, `repo.mk` (kein Filter auf Testnamen) | geprüft, ohne Befund |
| Lifecycle-Bewegung des Slice (`open/` → `next/` → `in-progress/`, je reiner Move vor dem Nachzug-Commit) | geprüft, ohne Befund (`AGENTS.md` §3.3) |
| Gates (`AGENTS.md` §3.6) | geprüft: keine Schwelle gesenkt; ein Test wird durch einen leichteren ersetzt, im Slice und in der Fundstelle benannt |
| `make test-docker` / Lauf des neuen Tests | nicht geprüft: Auftrag (F-6) |
| Go-Übersetzung mit Tag `docker` | nicht geprüft: kein Toolchain im Prüfkontext |
| `.claude/worktrees/` | nicht angefasst (paralleler Lauf) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Folge-Entwurf trägt überholte Test-Annahme · Plan-Tabelle gegen Diff nicht abgeglichen · Zustandsfeld mit Chronik · Reste der entfernten Prüfung · Zusage weiter als die Suite · Test ohne Lauf

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): der Entwurf, der das Lastenheft anpassen soll, rechnet mit drei Instanzen und einem Messskript, die der entschiedene Stand nicht liefert. Der Code der Test-Datei ist nach dem Lesen übersetzbar; über seinen Lauf entscheidet erst `integration-docker` (F-6).

**Übergabe:** Findings gehen an den Implementer; F-1 bis F-5 sind Doku-/Plan-Nachzüge ohne Produktänderung. Die Finding-Klassen gehen in die Slice-Closure §7. DoD-/Spec-Konformität prüft der Verifier separat.
