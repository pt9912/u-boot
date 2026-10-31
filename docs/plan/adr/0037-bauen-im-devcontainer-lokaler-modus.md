# ADR 0037: Bauen im Devcontainer — lokaler Modus für die Pflicht-Targets neben dem Docker-Modus, Schalter ausdrücklich, Pins im Devcontainer-Dockerfile, CI unverändert

**Status:** Accepted

**Datum:** 2026-10-31

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten), [`LH-FA-DEV-001`](../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen), [ADR-0004](0004-ci-system.md)

**Schärft:** [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) (Regeln des Docker-only-Workflows: neuer lokaler Modus im Devcontainer) und [`SPEC-008`](../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) (Spalte *Zweck* der Pflicht-Targets nennt beide Modi)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Pflicht-Targets des Makefiles bauen ausschließlich mit `docker build` und `docker run` ([`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows), erster Punkt); am Host sind nur Docker und `make` vorausgesetzt, eine Go-Toolchain nicht. Das Ziel dieser ADR ist ein schnelleres Bauen im Devcontainer, den `u-boot` für dieses Repo erzeugt, ohne diese Zusage für Host und CI aufzugeben.

**Stand im Code (gelesen am 2026-10-31, mit Fundstelle).**

- `repo.mk` Zeilen 81 bis 84 und 93 bis 103, 141 bis 145: Jedes Target ruft `$(DOCKER_BUILD) --target <stage>`. Für `lint`, `test` und `coverage-gate` steht zusätzlich `--no-cache-filter` (Zeilen 77 bis 79); der Kommentar dort nennt den Grund: ein veralteter Layer soll keinen Fehler verdecken.
- `Dockerfile`: sechs `COPY . .` (jede Stage, die den Quelltext braucht), kein `--mount=type=cache` (Zählung: null). Die Stage `test` führt `go test ./...` (Zeile 82), die Stage `coverage` führt dieselbe Suite mit Coverage-Profil noch einmal (Zeilen 138 bis 143).
- `.dockerignore`: schließt `docs/`, `spec/`, `harness/` und `.harness/` nicht aus; der Build-Kontext hat etwa 24 MB.
- `.claude/hooks/pretooluse-command-guard.sh` (Kopfkommentar): blockt Host-Toolchains in der Befehlsposition, auch `go`; `make`-Aufrufe sind davon nicht betroffen.
- `u-boot.yaml` und `.devcontainer/`: Sandbox-Profil mit verschachteltem Docker. In `u-boot` 0.16.1 setzt `devcontainer.features.<feature>.version` den Tag des Features, nicht die Version der Toolchain; `go:1.27.3` gibt es in der Registry nicht (`docker manifest inspect`, 2026-10-31), `go:1` und `go:1.4.0` schon. Das Feature `go` hat die Option `version` (Standard `latest`), die `u-boot` nicht schreibt.

**Messung (2026-10-31).** Linux `amd64`, Docker 29.9.2, ein Lauf je Zelle (Schwankung nicht bestimmt):

| Gate | Cache warm, Quellen unverändert | nach einer Zeile in einer Markdown-Datei |
|---|---|---|
| `make lint` | 30,7 s | 20,4 s |
| `make test` | 20,4 s | 18,5 s |
| `make coverage-gate` | 19,9 s | 19,2 s |
| `make docs-check` | 6,3 s | — |
| `make ablage-pruefung` | 14,4 s | — |
| `make baseline-verify` | 0,1 s | — |

Die Zeiten sind mit und ohne Änderung gleich: `--no-cache-filter` lässt die drei Stages immer neu laufen, und weil das `Dockerfile` keinen Cache-Mount hat, bleibt zwischen zwei Läufen kein Go-Build- und kein Test-Cache erhalten (aus dem Dockerfile gelesen, nicht gesondert gemessen).

Gilt die Annahme „der Docker-Pfad ist für Entwickler schnell genug“ nicht mehr, kippt die Entscheidung nicht, sie entfällt an einem Punkt: der lokale Modus wird dann nicht mehr gebraucht.

## Entscheidung

Vom Projektinhaber am 2026-10-31 entschieden: **die Erweiterung von [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows)** um das Bauen im Devcontainer ohne `Dockerfile`-Stages, **erst der Devcontainer-Modus, dieser aber konsequent**, **ein Nachweis statt zwei**, **keine Prüfskripte für Pins in diesem Vorhaben** (das Gate dafür liefert d-check, Issue 14 in `pt9912/d-check`) und die Annahme, dass **die CI zunächst unverändert** bleibt. Die Mechanik unten schlug der Autor vor; der Projektinhaber hat sie mit der Annahme dieser ADR am 2026-10-31 bestätigt.

1. **Zwei Modi.** Der *Docker-Modus* ist der heutige und bleibt der Standard. Der *lokale Modus* führt die Befehle der Stages mit der im Devcontainer installierten, gepinnten Toolchain aus, ohne `docker build` der Stage. Am Host und in der CI gilt der Docker-Modus; eine Go-Toolchain am Host bleibt für Standard-Targets nicht vorausgesetzt.
2. **Der Schalter ist ausdrücklich.** Die Variable `BUILD_MODE` (Standard `docker`) wählt den Modus. Der Devcontainer setzt `BUILD_MODE=local` per `ENV` im handgeschriebenen Teil von `.devcontainer/Dockerfile` (hinter dem verwalteten Block; `u-boot generate devcontainer` lässt Zeilen dort stehen, gemessen am 2026-10-31). Eine Erkennung von `go` im `PATH` gibt es nicht: ein Host mit Go würde sonst still umschalten.
3. **Umfang.** Im lokalen Modus laufen `deps`, `compile`, `lint`, `test`, `coverage-gate`, `test-docker`, `govulncheck` und `build-binaries` ohne `docker build`. Im Docker-Modus bleiben `build`, `run`, `image-scan` und `packages`: sie arbeiten am Runtime-Image oder in dessen Release-Umfeld.
4. **Eine Quelle für die Stage-Befehle.** Die Befehle von `compile`, `lint`, `test` und `coverage` stehen in einem Skript unter `scripts/`, das das `Dockerfile` und der lokale Modus gleich aufrufen. Das Ergebnis des Docker-Modus ändert sich nicht.
5. **Pins im Devcontainer.** Go und `golangci-lint` kommen als Release-Archive mit SHA-256 in den handgeschriebenen Teil von `.devcontainer/Dockerfile`, mit den Versionen `GO_VERSION` und `GOLANGCI_LINT_VERSION` aus `repo.mk`. Jeder dieser Pins steht damit an drei Stellen (`Dockerfile`, `repo.mk`, `.devcontainer/Dockerfile`); ein Abgleich ist hier nicht vorgesehen und kommt mit d-check.
6. **Nachweis.** `make gates` im Modus der Umgebung ist der Nachweis für den Handoff und die Closure. Ein zweiter Lauf im anderen Modus gehört nicht zum Prozess; die CI läuft im Docker-Modus bei jedem Pull Request.
7. **Test-Cache.** Im lokalen Modus darf der Ergebnis-Cache von Go Testergebnisse wiederverwenden (Schlüssel: Inhalt der Eingaben). Der Docker-Modus behält `--no-cache-filter`. Ob `go test` mit Coverage-Profil Ergebnisse cached, ist nicht gemessen; der Slice misst es zuerst und legt danach fest, ob Test und Coverage ein Lauf werden.
8. **Go-Cache im Workspace-Volume.** `GOCACHE` zeigt auf einen Pfad im Workspace-Volume des Devcontainers, damit `u-boot devcontainer up --recreate` den Cache behält.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun | keine Regeländerung | `lint`, `test` und `coverage-gate` laufen jedes Mal vollständig (nacheinander gemessen etwa 71 s; die Gesamtzeit von `make gates`, dessen Checks parallel laufen, ist nicht gemessen) |
| B — Cache-Mounts und engeres `.dockerignore` im `Dockerfile` | Docker-only bleibt unberührt, kleiner Eingriff | `--no-cache-filter` lässt die Stages trotzdem neu laufen; Test und Coverage laufen doppelt; Cache bleibt an den Docker-Daemon gebunden |
| C — Feature `go` im Devcontainer | keine eigene Installation | `u-boot` setzt die Go-Version nicht (nur `latest`), bringt `SYS_PTRACE` und `seccomp=unconfined` mit |
| D — Lokaler Modus, automatisch bei `go` im `PATH` | kein Schalter | schaltet am Host mit Go still um |
| E — Gemeinsames Toolchain-Image, Targets per `docker run` mit Bind-Mount und Cache-Volumes | ein Modus für Host, Devcontainer und CI; Pins an einer Stelle | größerer Umbau von [`SPEC-008`](../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) und der CI; Entscheidung des Projektinhabers: später |
| **F — Lokaler Modus mit ausdrücklichem Schalter, Docker-Modus als Standard** | Zeitgewinn im Devcontainer, Host und CI unberührt, kleinster Schritt Richtung E | zwei Modi können auseinanderlaufen, bis die CI es meldet; Pins an drei Stellen |

## Konsequenzen

- Positiv: Go-Build- und Test-Cache bleiben zwischen Läufen erhalten, es gibt keinen Kontext-Transfer, ein einzelnes Paket oder ein einzelner Test soll gezielt startbar sein (vorgesehen: `make test PKG=… RUN=…`).
- Negativ: Zwei Modi; eine Abweichung fällt erst in der CI auf. Go und `golangci-lint` stehen an drei Stellen. `golangci-lint`, Go und Docker im Devcontainer sind Teil der Lieferkette des Devcontainers, nicht mehr nur des Build-Images.
- Folgepflicht: [`SPEC-031`](../../../spec/spezifikation.md#spec-031--regeln-des-docker-only-workflows) und [`SPEC-008`](../../../spec/spezifikation.md#spec-008--eigenschaften-und-pflicht-targets-des-makefiles) werden angepasst; `AGENTS.md` §3.1 und `docs/maintainer/quality.md` folgen. Der Pin-Abgleich über drei Dateien kommt mit dem d-check-Modul (Issue 14 in `pt9912/d-check`).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `make -n` | Im lokalen Modus enthält keiner der in Entscheidung 3 genannten Docker-fremden Targets ein `docker build`; im Docker-Modus (Standard) enthalten sie es | `make modus-pruefung` (neu, kein Teil der Gate-Kette) |

## Re-Evaluierungs-Trigger

Wenn die CI auf das Toolchain-Image umgestellt wird (Option E), ersetzt diese Entscheidung ihr Ergebnis. Wenn d-check das Modul für gleiche Pins in mehreren Dateien liefert, wird der Abgleich der drei Stellen eingebunden. Wenn eine Abweichung zwischen den Modi in der CI auftritt, wird über einen zweiten Nachweis neu entschieden.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-31 | Proposed | — |
| 2026-10-31 | Accepted | Entscheidung des Projektinhabers im Gespräch am 2026-10-31 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
