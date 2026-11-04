# Review-Report: slice-v1-test-host-umgebung-isolieren — 2026-11-04

**Review-Art:** Code — Diff gegen Plan + Konventionen.

**Gegenstand:** Commit `4995f65e` (HEAD), Slice [`slice-v1-test-host-umgebung-isolieren`](../plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md).

**Skill:** `.harness/skills/reviewer.md` @ HEAD
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-11-04

**Eingangs-Kontext:**

- Slice-Plan [`slice-v1-test-host-umgebung-isolieren`](../plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md) (§3, §6)
- [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)
- `AGENTS.md` §3.7, §3.9, §3.10
- `harness/review.md`, `harness/conventions.md` (Sub-Area ADP, E2E: Greenfield)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | `isolateHostEnv` pinnt `GIT_TOKEN`, `GIT_CONFIG_GLOBAL` und `GIT_CONFIG_NOSYSTEM`, nicht aber die übrigen Variablen, die `git config` und `git init` lesen (`GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_<n>`/`GIT_CONFIG_VALUE_<n>`, `GIT_CONFIG_PARAMETERS`). Läuft die Suite aus einem Git-Hook (der setzt `GIT_DIR`) oder mit injizierter Konfiguration, kommt die Antwort von `ConfigValue` wieder vom Host. | Maintainability; Ziel des Slice (Antworten nie vom Host) | `internal/adapter/driven/git/hostsource_test.go` · `t.Setenv("GIT_CONFIG_NOSYSTEM", "1")` | nein — Review-only (kein Sensor stellt einen Hook-Kontext nach) | Isolations-Helfer pinnt nicht alle gelesenen Git-Variablen |
| F-2 | LOW | Der Doc-Kommentar `// scriptBinary returns a script that prints its lines …` steht vor `TestHostSource_GitToken_HostVariableWins`. Das Symbol `scriptBinary` gibt es im Baum nicht (die Funktion heißt `writeScript`); der Kommentar liegt direkt unter dem neu eingefügten Helfer. | `AGENTS.md` §3.7, §3.9 (b)/(c) | `hostsource_test.go` · `// scriptBinary returns a script that prints its lines and` | ja — `make gates` fängt es nicht; Review-only | Kommentar verweist auf nicht vorhandenes Symbol |
| F-3 | LOW | Plan-Tabelle §3 führt `docs/plan/planning/observations/BEO-ADP/test-liest-host-umgebung/` als „neu". `git diff --stat HEAD~1 HEAD` zeigt nur die beiden Testdateien; das Verzeichnis kam mit `09d765f0`, die Evidence-Änderung mit `HEAD~1`. Der Slice-Plan selbst (Statusänderungen) und die Evidence fehlen in der Tabelle, soweit sie der Slice berührt. | `AGENTS.md` §3.10 (Abgleich Tabelle gegen `git diff --stat`) | [`docs/plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md`](../plan/planning/in-progress/slice-v1-test-host-umgebung-isolieren.md) · `\| docs/plan/planning/observations/BEO-ADP/test-liest-host-umgebung/ \| neu \|` | nein — Review-only | Plan-Tabelle weicht vom Diff ab |
| F-4 | INFO | Dasselbe Muster (Meldung gibt eine Ausgabe aus, die den Wert von `GIT_TOKEN` enthalten kann) steht noch in `TestLauncher_Up_…`: `t.Errorf("the environment of the run client = %q, want %q", got, "GIT_TOKEN="+resolved)`. `got` ist die `grep`-Zeile des Prozess-Environments; der Launcher hängt Werte an `os.Environ()` an, ein Host-`GIT_TOKEN` würde im Fehlerfall ausgegeben. Placeholder-Werte der Tests sind kein Leck. Slice §6 nennt „weitere Tests" als nicht gesucht. | [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) (Wert nicht ausgeben) | `internal/adapter/driven/docker/devcontainer_launcher_test.go` · `the environment of the run client = %q` | Review-only | Fehlermeldung gibt Token-Wert aus |
| F-5 | INFO | Der geänderte Docker-Test (Build-Tag `docker`) lief nicht; im Slice §6 benannt. Die Meldungsänderung ist kompilierbar gelesen (`out != ""`, `err` bleibt `%v`), aber nicht belegt. Zuständigkeit: Verifier (`make test-docker`). | Slice §6 | `internal/e2e/devcontainer_shell_docker_test.go` · `output set=%t` | `make test-docker` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Isolation `GIT_TOKEN` (Adapter liest nur `os.Getenv("GIT_TOKEN")`; `""` gilt als ungesetzt, passt zu `t.Setenv(…, "")`) | geprüft, ohne Befund |
| `GH_TOKEN`, `HOME`, `XDG_CONFIG_HOME` | geprüft, ohne Befund: `gh` ist in den Tests ein Skript; mit gesetztem `GIT_CONFIG_GLOBAL` liest git weder `~/.gitconfig` noch die XDG-Datei (Git 2.56 im Container; die Variable verlangt Git ≥ 2.32) |
| Aufrufe von `isolateHostEnv` (6 in Plan genannt, 6 im Diff; `HostVariableWins` setzt `GIT_TOKEN` selbst und ruft `gh` nie) | geprüft, ohne Befund |
| Kommentare im Diff (§3.9 a/b/c) | geprüft; Doc-Kommentar des Helfers im Indikativ, ohne Chronik; einziger Befund F-2 (vorbestehender Kommentar) |
| Sprachvertrag, Importregeln, `//nolint` | geprüft, ohne Befund (nur Testcode, englische Meldungen) |
| Andere Tests in `internal/adapter/driven/**` | geprüft: `git_test.go` ohne Env-/Konfig-Zugriff; `localtemplates_test.go` pinnt `HOME` selbst; Launcher-/Shell-Tests übergeben Werte per Map (Befund F-4 als INFO) |
| Tests ausgeführt | nicht geprüft: `make gates` und `make test-docker` nicht von diesem Lauf ausgeführt (Verifier-Aufgabe, keine Code-Änderung im Review) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Isolations-Helfer pinnt nicht alle gelesenen Git-Variablen · Kommentar verweist auf nicht vorhandenes Symbol · Plan-Tabelle weicht vom Diff ab · Fehlermeldung gibt Token-Wert aus

## Verdikt

**Merge-blockierend:** nein — keine HIGH/MEDIUM-Findings.

**Übergabe:** Findings F-1 bis F-3 an den Implementer; F-4 als Eingang für das Register `BEO-ADP/test-liest-host-umgebung`; F-5 an den Verifier. Finding-Klassen gehen in die Slice-Closure §7.
