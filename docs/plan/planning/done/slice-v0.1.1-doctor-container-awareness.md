# Slice v0.1.1: Doctor-Container-Awareness

## Auslöser

**Real-world-Befund 2026-05-31, kurz nach v0.1.0-Release:** ein
Nutzer hat `docker run -it --rm ghcr.io/pt9912/u-boot:latest doctor`
gegen ein gesundes Host-System ausgeführt. Resultat: vier `error`-
Diagnostiken, weil `doctor` die Tools `docker`, `docker compose`
und `git` im Container-PATH sucht. Das distroless-Image
(`gcr.io/distroless/static-debian12:nonroot`, ADR-0007 / Dockerfile
L142) enthält bewusst keine dieser Binaries — nur das u-boot-Binary.

Konkrete Fehlausgabe (Auszug):

```
✗  docker.installed    docker binary not available: docker version failed: exec: "docker": executable file not found in $PATH
✗  docker.compose.installed  ...
✗  docker.reachable    ...
✗  git.installed       git binary not available: ...
```

Ein ansonsten gesunder Host-Stand wird als „4 errors" gemeldet,
weil `doctor` keinen Begriff von „läuft im Container" hat. Das
ist ein UX-Bug, der aus der Distributions-Entscheidung (ADR-0007:
GHCR-Image als primärer Distributionsweg) hervorgeht, aber dort
nicht antizipiert wurde.

## Aufhebungsbedingung

`doctor` läuft im Container nicht mehr in den False-Positive-
Pfad. Mindestens eine der folgenden Strategien ist umgesetzt:

1. **Container-Detection + Skip mit Hinweis:** `doctor` erkennt
   per `/.dockerenv` / `/run/.containerenv` / `cgroup`-Heuristik,
   dass es im Container läuft. Die Host-Voraussetzungs-Checks
   (`docker.*`, `git.installed`) werden als
   `severity: info, status: skipped` markiert mit einer
   Repair-Hint-Zeile, die erklärt: „Diese Checks gelten dem Host,
   nicht dem Container. Installiere u-boot lokal (siehe
   `slice-v2-binary-distribution`) oder bewerte den Host
   gesondert."
2. **Doc + Help-Text:** `doctor --help` und die README-Quickstart
   erläutern, dass `doctor` für die Host-Installation gedacht ist;
   im Container-Modus liefert es nur die Datei-Checks
   (`uboot.yaml`, `compose.yaml`, `.devcontainer/*`) sinnvoll.
3. **Binary-Distribution-Trigger:** dieser Befund zieht den ersten
   ADR-0007-Re-Evaluation-Trigger („erste konkrete Nachfrage nach
   Cross-Plattform-Distribution") in einen aktiven Slice-Plan
   `slice-v2-binary-distribution`.

## Akzeptanzkriterien

- `doctor` im Container-Modus (Detection grünt z. B. via
  `/.dockerenv`) liefert KEINE `error`-Diagnostik für
  `docker.installed`, `docker.reachable`, `docker.compose.installed`,
  `git.installed`. Stattdessen `skipped`-Status mit Repair-Hint.
- Exit-Code im Container-Modus bei sonst gesundem Projekt: 0
  (statt 11 wie heute), weil keine Errors mehr.
- `doctor --help` und README dokumentieren das Container-Verhalten.
- Slice-Plan `slice-v2-binary-distribution.md` ist in `open/`
  angelegt (Trigger ADR-0007 §Folgepunkte 1 ist materialisiert).

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `9a99bbf` | Port `driven.RuntimeEnvironment` (`InContainer() bool`) + Adapter `internal/adapter/driven/runtime/runtime.go` via `/.dockerenv` und `/run/.containerenv`. Tests: sieben Tabellen-Cases plus Production-Smoketest. Architektur-Korrektur: nicht in `domain/` (kein I/O dort), sondern als Driven-Port-Pattern analog `Clock`/`Git`. |
| T2 | `c35360f` | `internal/hexagon/application/doctor`: vier Host-Prerequisite-Checks (`git.installed`, `docker.installed`, `docker.reachable`, `docker.compose.installed`) bekommen frühen `inContainer()`-Gate → `SeverityInfo` + gemeinsamer `hostHintSkippedInContainer`-Hint. Probes werden NICHT aufgerufen (short-circuit). nil-runtime preserves pre-v0.1.1 behaviour (12 bestehende Test-Aufrufe via perl-bulk auf 6-arg-Signatur gehoben). Zwei dedizierte neue Pin-Tests (`HostChecksSkipped_WhenInContainer`, `HostChecksRunNormally_WhenRuntimeNil`). cmd/uboot/main.go: `runtime.New()`-Adapter im Wiring. |
| T3 | `111e725` | `CHANGELOG.md ## [0.1.1] - TBD` mit Added/Changed/Notes (Datum wird beim Tag-Push gesetzt). README.{md,de.md} "doctor and the container caveat"-Block von Future-tense auf "Starting with v0.1.1 / Ab v0.1.1, doctor detects … emits SeverityInfo … exit code 0 (not 11)" umformuliert. |
| T4 | dieser Commit | Slice-Plan nach done/; roadmap-Sync (Carveout-Auflösungs-Tabelle + §Nächste Schritte Punkt 1). Release-Tag `v0.1.1` bleibt Nutzer-Aktion (analog v0.1.0-T4): (a) CHANGELOG-Datum auf Push-Datum aktualisieren, (b) push, (c) ersten CI-Lauf abwarten, (d) `git tag v0.1.1 && git push origin v0.1.1` → publish.yml triggert GHCR-Push. |

## Out of Scope

- Volle Binary-Distribution (eigener Slice
  `slice-v2-binary-distribution`, in `open/` als Trigger-Ziel).
- Devcontainer-Mount-Variante (z. B. Host-PATH ins Container
  durchreichen via `-v`); zu hässlich für die README, lässt sich
  aber nach Bedarf in der Doku als Workaround erwähnen.

## Bezug

- Auslöser: v0.1.0-Container-Run am 2026-05-31, real-world-feedback.
- Hängt von:
  [ADR-0007](../../adr/0007-distributionswege-ghcr.md) §Folgepunkte 1
  („erste konkrete Nachfrage nach Cross-Plattform-Distribution")
  als Trigger für `slice-v2-binary-distribution`.
- Plan-Anker: dieser Slice plus
  [`slice-v2-binary-distribution`](slice-v2-binary-distribution.md).
- Phase: V1-Followup (post-v0.1.0, Pre-v0.1.1).
