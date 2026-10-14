# Review-Report: [`slice-cr-sandbox-docker-runtime`](../../docs/plan/planning/done/slice-cr-sandbox-docker-runtime.md) (Folgelauf r2) — 2026-10-14

**Review-Art:** Code — Diff gegen die Befunde des Erstlaufs (`2026-10-14-sandbox-docker-runtime`), Spec ([`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md), Replay-Regeln (`harness/replay.md` §Replay Rules) und Pin-Politik ([`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)).

**Gegenstand:** Diff-Range `6d726e0..HEAD` (Commits `9fa95ce`, `2a888f7`, `8bc3d15`).

**Skill:** `.harness/skills/reviewer.md` @ 8bc3d15 · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-14

**Eingangs-Kontext:**

- Erstreport `docs/reviews/2026-10-14-sandbox-docker-runtime.md` (F-1..F-8)
- [`slice-cr-sandbox-docker-runtime`](../../docs/plan/planning/done/slice-cr-sandbox-docker-runtime.md), [`slice-v1-sandbox-docker-generator`](../../docs/plan/planning/done/slice-v1-sandbox-docker-generator.md), [`slice-v1-sandbox-docker-degradation-doctor`](../../docs/plan/planning/done/slice-v1-sandbox-docker-degradation-doctor.md), [`slice-v1-sandbox-docker-integrationstest`](../../docs/plan/planning/done/slice-v1-sandbox-docker-integrationstest.md)
- [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md), [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik)
- `AGENTS.md` (Hard Rules), `harness/review.md`, `harness/roles.md`, `harness/replay.md`

Zusätzlich lokal ohne `--privileged` geprüft: Im Basis-Image `mcr.microsoft.com/devcontainers/base:debian` liegt `/usr/bin/timeout` (GNU coreutils 9.7); `/var/run` ist Symlink auf `/run`.

## Stand der Erstbefunde

| Erstbefund | Stand |
|---|---|
| F-1 (Stale PID/Socket) | Skript räumt auf (`docker-init.sh.tmpl:25`); der neue Neustart-Test belegt die Behebung nicht diskriminierend — siehe F-1 unten |
| F-2 (Replay-Zeile) | behoben (`harness/replay.md:52`); Inhalt der Zeile unvollständig — siehe F-4 |
| F-3 (Spec-Satz `--cap-add=NET_ADMIN`) | behoben (`spec/spezifikation.md:326` nennt NET_ADMIN nur mit Egress) |
| F-4 (Docker Desktop) | unverändert offen, nicht im Diff (Verifier, AK 6) |
| F-5 (60-s-Grenze) | teilweise behoben — siehe F-2 unten |
| F-6 (Nachziehstellen) | behoben (`cli-json-contract.md:616`, `benutzerhandbuch.md:1088`) |
| F-7 (DoD-Zustand) | nicht im Diff (Verifier schließt) |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `TestE2E_SandboxDevcontainer_NestedDocker_Restart` scheitert auch ohne die `rm -f`-Zeile nicht: Nach `docker restart` läuft in der frischen PID-Namespace nur `sleep` als PID 1; die alte `dockerd`-PID ist praktisch nie belegt, und `dockerd` ersetzt einen toten Socket selbst. Der Test löst die Fehlerbedingung des Erstbefunds („PID N belegt“) nicht aus und pinnt die neue Zeile nicht. Zudem prüft er nur `overlay2`, nicht, dass der Storage des ersten Daemons (Volume) wieder erreichbar ist, wie der Kommentar zusagt. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Replay Rules (Safety-Pfad) | `internal/e2e/sandbox_devcontainer_docker_test.go:479`, `:505-514` | nein — nur ein privilegierter Lauf mit entferntem `rm -f` zeigt, ob der Test rot wird (`make test-docker`) | Regressionstest pinnt Fehlerbedingung nicht |
| F-2 | LOW | Die 60-s-Zusage ist weiter nicht hart: Die Schleife zählt Runden (`i < 60`), jede Runde kostet bis zu 5 s (`timeout 5 docker info`) plus 1 s `sleep`; bei hängendem `docker info` vergehen bis zu 360 s bis zur Meldung. Ein einzelner Aufruf ist begrenzt, die Gesamtdauer nicht. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Spezifikation: „nach spätestens 60 Sekunden“) | `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl:18`, `:32-41` | nein — Review-only (Test nur mit hängendem Daemon-Socket) | Zeitzusage „spätestens“ nicht durchgesetzt |
| F-3 | LOW | Die Spezifikation führt als Default-Pin weiter „Stand 2026-10-02: `1.27.1`“; der Pin-Bump auf `1.27.2` zieht diese beiden Stellen nicht nach (Dockerfile, Makefile und CHANGELOG sind konsistent auf `1.27.2`). | [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik), `AGENTS.md` §5 Regel 5 | `spec/spezifikation.md:785`, `:800` | ja — Gegenlesen (`grep 1.27.1`); Wiederholung des Musters aus dem Report `2026-10-02-welle-spec-technik-stratum` F-10 | Veralteter Pin-Wert in Defaults |
| F-4 | LOW | Die Replay-Zeile für `nestedRuntime: docker` nennt weder „stale block → Block ersetzt, User-Content bleibt“ noch „User-Skript ohne Marker → Konflikt ohne Write“, obwohl die neuen Tests `..._StaleBlockKeepsUserContent` und `..._ScriptConflict` genau das pinnen (die Podman-Zeile führt den Konfliktfall). `..._StaleBlockKeepsUserContent` prüft nur `Action != NoOp`, nicht `UpdatedBlock`, wie die Replay-Konvention „stale managed block → UpdatedBlock“ es für Stale-Cases festhält. | Hard Rule §3.8 (Replay), `harness/replay.md` §Replay Rules | `harness/replay.md:52`, `internal/hexagon/application/generate_sandbox_test.go:827-829` | ja — Gegenlesen; Test bleibt auch bei `Created` grün | Replay-Tabelle trägt nicht alle gepinnten Cases |
| F-5 | LOW | Die Aufräumzeile läuft immer, wenn der erste `ready`-Aufruf fehlschlägt, auch wenn ein lebender, nur langsamer `dockerd` (Aufruf läuft in das 5-s-Limit) noch existiert: PID-Datei und Socket werden unter ihm entfernt, danach startet ein zweiter `dockerd`. Auslöser wäre ein Re-Run von `u-boot-docker-init` (manuell oder nach einem `warn`-Timeout) gegen einen noch startenden Daemon. | [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl:20-25` | nein — Review-only (abgeleitet, kein Lauf) | Aufräumen ohne Lebendprüfung |
| F-6 | INFO | Die Shell-Variable `timeout=60` und der Befehl `timeout 5 docker info` teilen den Namen; das funktioniert in `sh` (getrennte Namensräume), erschwert aber das Lesen. Der Commit `2a888f7` trägt die Begründung des Pin-Bumps nur in der Betreffzeile, kein Body (Dockerfile-Kommentar: „rationale goes into the commit body“). | Maintainability, [`LH-FA-BUILD-003`](../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik) | `internal/hexagon/application/templates/devcontainer/docker-init.sh.tmpl:8`, `:18`; `Dockerfile:22-24` | ja — Gegenlesen | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docker-init.sh.tmpl` Exit-Codes | `exit 0` (bereit/degradiert) und `exit 11` (`fail`) an genau einer Stelle (`unavailable`); keine neuen Pfade; geprüft, ohne Befund |
| `set`-Verhalten | kein `set -e`/`-u`; der Hintergrundstart, `kill -0` und die Schleife sind darauf nicht angewiesen; ein fehlschlagendes `timeout`/`docker` ergibt Rückgabewert ungleich 0 statt Abbruch; geprüft, ohne Befund |
| `timeout` im Basis-Image | `/usr/bin/timeout` (coreutils 9.7) im Debian-13-Basis-Image vorhanden; Skript läuft unter `sudo sh` mit diesem PATH; geprüft, ohne Befund |
| rm-Pfade | `/var/run` → `/run`; PID-Datei von `dockerd` (`/var/run/docker.pid`), Socket und `containerd.pid` des von `dockerd` verwalteten containerd sind die richtigen Reste; `rm -f` ist bei fehlenden Dateien unkritisch; geprüft, ohne Befund (Randfall: F-5) |
| Fehlerpfade Early-Exit / Timeout | unverändert gegenüber `6d726e0`; Meldungen Englisch ([`LH-LESE-002`](../../spec/lastenheft.md#lh-lese-002--sprache)); geprüft, ohne Befund |
| Golden-Tests gegen Replay Rules | Fresh-State (`Created`), Stale-Block mit User-Content, Konflikt „no partial write“ (`devcontainer.json`, `Dockerfile` unberührt), Skript-Zusagen als Substring-Pins vorhanden; Double-Run-NoOp aus `6d726e0`; geprüft (Lücken: F-4) |
| Spec-Satz zu `runArgs` (Erst-F-3) | `spec/spezifikation.md:326` deckt `NET_ADMIN` nur mit Egress ab, konsistent zu `TestGenerateDevcontainer_SandboxEgressWithDocker`; Lastenheft unverändert und mit der Spezifikation verträglich; geprüft, ohne Befund |
| Decken-Regel | Diff in `spec/` ohne ADR-/Slice-/Roadmap-Bezug; geprüft, ohne Befund |
| Nachziehstellen (Erst-F-6) | `docker-init.sh` in Rollback-Beschreibung und „Grenzen“ ergänzt; weitere `sandbox-init.sh`-Fundstellen in `docs/user`/`docs/maintainer`/READMEs enthalten keine Lücke; geprüft, ohne Befund |
| Pin-Bump 1.27.1 → 1.27.2 | `Dockerfile:29`, `Makefile:15` identisch; `Makefile` reicht `GO_VERSION` per `--build-arg` und an `golang:$(GO_VERSION)` durch; kein weiterer Pin in `go.mod`/Workflows; `CHANGELOG` unter `[Unreleased]`-nachfolgendem, noch ungetaggtem `[0.8.0]` (`### Security`, Keep-a-Changelog-Kategorie) konsistent; geprüft, ohne Befund (Rest: F-3, F-6) |
| Gate-Lockerung | keine Schwelle, Lint-Regel oder Gate-Definition im Diff gesenkt; geprüft, ohne Befund |
| Kommentar-Regel §3.7 | neue Kommentare (`docker-init.sh.tmpl:24`, Testkommentare) im Indikativ, Zusage/Kopplung; geprüft, ohne Befund |
| Architektur / Importregeln, Exit-Code-Klassifikation | nur Templates, Tests und Doku berührt; keine neuen Sentinels; geprüft, ohne Befund |

Nicht geprüft:

- Echter Lauf von `make test-docker` / `TestE2E_SandboxDevcontainer_NestedDocker_Restart` — Auftrag untersagt `--privileged`-Container; F-1 und F-5 sind aus Code und `dockerd`-Verhalten abgeleitet.
- Wahrheitsgehalt der genannten CVE-Kennungen und deren Behebung in Go 1.27.2 — kein Netz-/Advisory-Abgleich; Beleg wäre `make image-scan`.
- Docker Desktop, Podman als Host-Engine, Codespaces — laut [ADR-0018](../plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) ungemessen.
- `make test`/`make lint`/`make gates` — laut Auftrag grün, nicht erneut ausgeführt.
- `make docs-check` — nicht ausgeführt; Links auf `done/` lösen bis zum Closure `target-missing` auf (bekannt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Regressionstest pinnt Fehlerbedingung nicht · Zeitzusage „spätestens“ nicht durchgesetzt · Veralteter Pin-Wert in Defaults · Replay-Tabelle trägt nicht alle gepinnten Cases · Aufräumen ohne Lebendprüfung

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings. Die Erstbefunde F-1, F-2, F-3, F-6 sind im Code beziehungsweise in der Doku behoben, F-5 teilweise; F-1 (dieser Lauf) ist vor Closure zu klären, weil der Neustart-Test den Defekt nicht auslöst und die Behebung damit unbelegt bleibt.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. „Zeitzusage ‚spätestens‘ nicht durchgesetzt“ und „Veralteter Pin-Wert in Defaults“ treten wiederholt auf (Erstlauf F-5, Report `2026-10-02-welle-spec-technik-stratum` F-10). Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
