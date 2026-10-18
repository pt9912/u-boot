# Verification-Report: slice-v1-podman-aktuelle-version — 2026-10-18 <!-- d-check:ignore (Kennung im Titel) -->

**Rolle:** Verifier (Frage: „Bauen wir es richtig?“, gegen Plan und DoD; kein Review, keine Validierung).

## Scope

- Slice: [`slice-v1-podman-aktuelle-version`](../plan/planning/done/slice-v1-podman-aktuelle-version.md) (in-progress), Commit-Range `314b1e9..HEAD` (`d6bbe4d`, 12 Commits, 12 Dateien).
- IDs: [`LH-FA-DEV-007`](../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-NFA-SEC-004`](../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); Spezifikation [`LH-FA-DEV-007.b`](../../spec/spezifikation.md#lh-fa-dev-007b--ausstattung-bei-nestedruntime-podman); [`ADR-0023`](../plan/adr/0023-podman-aus-podman-static-im-sandbox-devcontainer.md) (Accepted).
- Eingang: DoD-Bestätigung des Implementers, Review-Report [`2026-10-18-slice-v1-podman-aktuelle-version.md`](2026-10-18-slice-v1-podman-aktuelle-version.md) (0 HIGH/MEDIUM, 3 LOW, 5 INFO).
- Host: Colima, Docker-Server 29.5.2, arm64, VM mit `kernel.apparmor_restrict_unprivileged_userns=0`; kein Host-Go.

## Sensors (selbst gefahren, Belege gesehen)

| Sensor | Ergebnis | Evidence |
|---|---|---|
| `make gates` | pass | lint, test, `coverage-gate: OK — coverage 91.00% meets threshold 90%`, `d-check: 362 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.1 OK`; Stempel `.harness/state/gates-passed.head` = `d6bbe4d` |
| `make doc-trace` | pass | `83 Anforderung(en), 0 Waise(n).` |
| `make doc-immutable RANGE="314b1e9..HEAD"` | pass | `12 Commit(s) — OK`, `0 Befund(e)` |
| `make test-docker` (vollständig) | pass | Exit 0, **789 s** (13 min 9 s) Gesamtlauf, `internal/e2e` 780,5 s bei `-timeout 30m`; kein Test-Fehler |
| Einzeltests, `-run … -v -count=1`, Image `u-boot:test-docker-tools`, Rezeptur `repo.mk` (Ziel `test-docker`), 338 s | pass, 0 Skip (`grep -ci skip` = 0) | `…_NestedPodman` 43,8 s; `…_NestedPodman_Build` 40,1 s; `…_NestedPodman_BuildEgress` 82,6 s; `…_NestedPodman_PinMismatchBreaksBuild` 11,2 s (`wrong_checksum`, `signing_subkey_as_the_pin`, `unknown_fingerprint_as_the_pin` je PASS); `…_MissingHostFuseDeviceRefusesStart` 42,0 s; `…_FailPolicyWithoutFuse` 33,7 s; `StatusFiles_` (`InspectorAgainstEngine`, `Egress`, `VFS`, `Daemon`) je PASS |
| `make docs-check` | pass | siehe unten (nach Ablage dieses Reports) |
| `make verify-depguard` | nicht gelaufen | keine Importregel berührt (Diff: Templates, Tests, Doku) |
| `make govulncheck`, `make image-scan` | nicht gelaufen | kein Security-/Release-Pfad berührt; kein Go-Modul-/Image-Wechsel |

Laufzeit-Einordnung: Die Messung P7 nennt „nicht gemessen“; jetzt gemessen sind 789 s, innerhalb der 30 Minuten, aber über dem Kommentar in `repo.mk` („zwischen sieben und elf Minuten“). Im Scratchpad liegt ein älterer Lauf desselben Pakets mit 536 s (`ok`, Herkunft und Stand unklar, **kein belastbarer Vergleich**); ob die längere Laufzeit auf den neuen Test (11 s) und Download-Zeiten oder auf die Host-Last zurückgeht, ist nicht gemessen. Der Kommentar in `repo.mk` ist zu prüfen (Hinweis, kein Befund dieses Slice).

## Praxislauf des Produkts (Docker-only, `make build` → Image `u-boot:latest`)

Wegwerf-Verzeichnis `~/uboot-verify-podman6` (zwei Projekte `plain`, `egress`), danach gelöscht; Probe-Images `vp6-*` entfernt.

| Prüfung | Ergebnis |
|---|---|
| `init --devcontainer --sandbox`, `config set … nestedRuntime podman`, `generate devcontainer` (Egress an/aus) | läuft; sechs Warnungen `LH-FA-DEV-007` je Lockerung, bei Egress zusätzlich `LH-FA-DEV-008` (`NET_ADMIN`); `runArgs` = die sechs Optionen (+ `--cap-add=NET_ADMIN`) | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| Erzeugtes Dockerfile | `ARG PODMAN_VERSION=6.1.4`, beide SHA-256 und Fingerprint `0CCF102C…B503` je einmal und gleich den Werten in `ADR-0023` und in `Dockerfile.tmpl`; `registries.conf` mit `['docker.io']`; Symlink `docker`; Debian nur `uidmap ca-certificates`; `nftables` nur im Egress-Projekt | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| Image bauen (arm64) | `docker build` plain 7 s (Cache), Egress-Image gebaut |
| Container mit sechs Lockerungen | `u-boot-sandbox-init` Exit 0; `podman version 6.1.4`; `crun` 1.30.1; `netavark` 2.1.0; `command -v podman` = `/usr/local/bin/podman`, `docker` löst auf `podman`; kein Debian-Paket `podman`; `/README.md` fehlt |
| `podman build` und `docker build` (`FROM docker.io/library/alpine:3`, `COPY`, `RUN`) | beide laufen, zweiter nutzt den Cache, `podman run` des Ergebnisses liefert `ran`/`hello`; Kurzname `FROM alpine` löst über `docker.io` auf |
| Egress (`u-boot-egress-init`, Default-Allowlist) | `FROM docker.io/library/debian:stable-slim` mit `apt-get update` läuft (`APT_OK`), ein nicht freigegebener Host löst nicht auf (`UNLISTED_BLOCKED`), `EGRESS_BUILD_OK` |
| `doctor` im Runtime-Image | unverändert sinnvoll im Sinn des Container-Modus: Host-Prüfungen `?` übersprungen (u-boot läuft im Container), Statusverzeichnis „nicht bestimmbar“; die Zeile `devcontainer.sandbox.runtime` meldet `/dev/fuse` als fehlend (Exit 11), weil der Prüfpfad im Runtime-Image-Container kein `/dev/fuse` sieht — Verhalten des unveränderten `doctor`-Codes, nicht Gegenstand des Slice |
| `amd64` (Emulation, `--platform linux/amd64`) | Dockerfile baut bis Prüfsumme, Signatur, Extraktion; `podman version 6.1.4`, `netavark 2.1.0`, `fuse-overlayfs 1.18`, `conmon 2.2.1` laufen. **Nicht prüfbar** und nicht behauptet: ein verschachteltes `run`/`build` auf `amd64` (Emulation; CI-Lauf nach dem Push), Docker Desktop, cgroups v1 |

## DoD-Abgleich (Zeile für Zeile)

- [x] **Liefer-Punkt 1 — Messung und Folge-ADR:** `ADR-0023` Accepted (2026-10-18, Projektinhaber; Geschichte-Zeile), MADR, `Schärft:` `LH-FA-DEV-007`/`-008`/`LH-NFA-SEC-004`, Teilablösung von Entscheidung 1 `ADR-0014` (und Aussage in `ADR-0021` Entscheidung 2); Messprotokoll mit Host, Engine, Architektur, Kernel, Sysctl, Digest, Quelle/Version/Prüfsumme/Signatur, Komponenten, Größe/Bauzeit, P1 bis P8; ADR-Index `Accepted` (`docs/plan/adr/README.md`). Verglichene Alternativen A bis H vorhanden. Erfüllt. Hinweis: Entscheidung 2 nennt das `VALIDSIG`-Feld „Feld 11“, Vorlage und Test prüfen `$12` (awk-Zählung mit `[GNUPG:]` als Feld 1; die eigene Rechnung: `$12` ist der Primär-Fingerprint, wie der Review gemessen hat). Eine Wortlaut-Ungenauigkeit in einer immutablen ADR, ohne Wirkung auf Code; kein Eingriff, in den Closure-Hinweisen zu nennen. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- [x] **Liefer-Punkt 2 — Umsetzung und Beleg:** `Dockerfile.tmpl` Block `NestedPodman`: Prüfsumme (`sha256sum -c -`) → Schlüsselabruf → `gpg --import` → `--verify` → `VALIDSIG`-Vergleich → `tar -xzf` (Reihenfolge im Golden `assertPodmanStaticDockerfile` gepinnt, Pins je genau einmal); Zweig für unbekannte Architektur (`exit 1`) im Golden gepinnt (F-2); `devcontainer_egress.go` und `SPEC-019` unverändert, belegt durch `TestGenerateDevcontainer_SandboxEgressWithPodman` (Hosts der Quelle nicht auf der Allowlist); e2e `…_NestedPodman`/`_Build`/`_BuildEgress` laufen mit 6.1.4, Versionsprüfung auf Hauptversion 6 (`podman version 6\.`); `harness/replay.md` nachgezogen; Negativprobe als Test `…_PinMismatchBreaksBuild` (dauerhaft, PASS) und vom Review per Mutation belegt. Erfüllt. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- [x] **Liefer-Punkt 3 — Doku:** Spezifikation `LH-FA-DEV-007.b` (Herkunft, Pin, Prüfung, Pakete, Archivinhalt, Build-Hosts, `registries.conf`, `vfs`-Drop-in) und Historiezeile; Handbuch (Herkunft und Version, Kurznamen, `vfs`, Gemessen, Grenzen); `CHANGELOG.md` `[Unreleased]` ohne Chronik, mit Hinweis auf Neugenerierung und Neubau bestehender Projekte, Messgrenzen genannt. Allowlist-Tabelle unverändert (keine Host-Änderung). Lastenheft, `ADR-0014`, `ADR-0019`, `ADR-0021` unverändert (`git diff --stat 314b1e9..HEAD` für diese Pfade leer; `make doc-immutable` OK). Erfüllt. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- [x] `make gates` grün; `make test-docker` grün (789 s, kein Skip). Erfüllt.
- [x] Review durchgeführt, Report liegt unter `docs/reviews/` (Reviewer-Rolle, anderer Kontext). Erfüllt.
- [x] Doku-Update bei öffentlichem Vertrag: Liefer-Punkt 3. Erfüllt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — **offen**: §7 des Slice ist „Leer — füllt der Planner bei der Closure“.
- [ ] Beobachtungs-Register fortgeschrieben (oder in §7 notiert, dass keine Beobachtung anfiel) — **offen**, Planner.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — **offen**: §6 trägt noch den Stand „offen bis …“ (siehe Auflagen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — **offen**: von der nächsten Welle-Closure getragen.
- Die DoD-Häkchen im Slice sind noch nicht gesetzt (Planner-Aufgabe bei der Closure); `Verantwortlich:` steht auf „—“.

## Review-Findings am Code/Text

| Finding | Bearbeitung | Befund des Verifiers |
|---|---|---|
| F-1 (LOW) Plan-Tabelle gegen Diff | `open/README.md` jetzt „gelesen, unverändert“ (Diff berührt die Datei nicht), README-Zeile auf 224–225 | erledigt; alle 12 Diff-Dateien stehen in §3 (ADR-Index in Zeile 1) |
| F-2 (LOW) Architektur-Zweig | im Golden als vollständige `case`-Zeile mit `exit 1` gepinnt | erledigt (am Test gelesen; Mutation nicht erneut gefahren) |
| F-3 (LOW) Build-Hosts | `github.com`, `objects.githubusercontent.com`, `keys.openpgp.org` in `LH-FA-DEV-007.b` und im Handbuch | erledigt | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| F-4 (INFO) Schlüsselablauf 2027-03-08 | nur im Slice §3 „Hinweise“ und in der ADR | als Wiedervorlage in §6/§7 und Register zu tragen (Auflage 3) |
| F-5 (INFO) Pins als `ARG` überschreibbar | §3-Hinweis „kein Handbuch-Hinweis“ | akzeptiert; Handbuch „Version ist festgelegt“ gilt für das erzeugte Dockerfile |
| F-6 (INFO) Archivinhalt | `rm -rf "$work" /README.md`, Golden und e2e `NO_ARCHIVE_README` (im Lauf PASS), Spezifikation nennt `/etc/containers/*` | erledigt |
| F-7 (INFO) `vfs` ohne Rückbau | §3-Hinweis | in §7/Register zu tragen |
| F-8 (INFO) Negativprobe-Wirkstufe | §3-Hinweis | in §7 zu tragen |

## Plan-vs-Code-Diff (§3, `git diff --stat 314b1e9..HEAD`, AGENTS.md §3.10)

Alle 12 geänderten Dateien haben einen Plan-Eintrag; jeder „update (umgesetzt)“-Eintrag hat einen Diff; die „gelesen, unverändert“-Zeilen (`devcontainer_egress.go`, `sandbox_status_docker_test.go`, `spec/lastenheft.md`, `README.md`, `README.de.md`, `releasing.md`, `open/README.md`) haben keinen Diff. Die gelesenen Aussagen stimmen mit dem Baum überein (README-Zeilen 12–13, 199, 224–225; `README.de.md` 13–14, 227–228; `releasing.md` ohne Treffer; `egressRegistryHosts` Zeile 63, `egressDockerAptHost` Zeile 70).

**Diskrepanzen (Fundstellen-Drift, LOW, kein inhaltlicher Befund):** drei Zeilenangaben sind nach den Review-Fixes (`a300817`, `d7c9821`) nicht nachgezogen: `TestGenerateDevcontainer_SandboxEgressWithPodman` steht in Zeile **629** (Plan: 624), `…_PinMismatchBreaksBuild` in Zeile **510** (Plan: 509), Handbuch „Grenzen“ in Zeile **1459** (Plan: 1451). Der Plan nennt die Tabelle als „nach der letzten Fix-Runde“ abgeglichen; diese Zeilen sind es nicht (Closure-Nachzieher).

## Traceability

| ID | Beleg |
|---|---|
| `LH-FA-DEV-007` / `.b` | Golden `TestGenerateDevcontainer_SandboxPodman*`, e2e `…_NestedPodman*`, `StatusFiles_VFS`; Spezifikation `.b`; Praxislauf | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `LH-FA-DEV-008` | `TestGenerateDevcontainer_SandboxEgressWithPodman` (Allowlist unverändert), e2e `…_BuildEgress`, Praxislauf Egress | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `LH-NFA-SEC-004` | `ADR-0023` Entscheidung 2 und 7 (Auslegung vom Projektinhaber mit `Accepted` bestätigt); Pin und Prüfung im Golden und `…_PinMismatchBreaksBuild` | <!-- d-check:ignore (Kennungen als Zitat im Report) -->
| `ADR-0023` | Golden, e2e, `make doc-immutable`, ADR-Index | <!-- d-check:ignore (Kennungen als Zitat im Report) -->

## Carveouts

Neu: keine. Gelöst: keine. Unverändert: alle (`docs/plan/planning/in-progress/carveouts.md` nicht im Diff).

## Nicht ausgeführt / nicht geprüft

- `make verify-depguard`, `make govulncheck`, `make image-scan` (Gründe oben).
- Verschachtelter Lauf auf `amd64`, Docker Desktop, cgroups v1, Hosts mit `kernel.apparmor_restrict_unprivileged_userns=1` (ADR-0023 „Nicht verifiziert“); CI-Lauf `integration.yml` nach dem Push steht aus. <!-- d-check:ignore (Kennungen als Zitat im Report) -->
- Mutation für F-2 (Architektur-Zweig) nicht erneut gefahren; am Test gelesen.

## Verdikt und Auflagen für die Closure

**DoD (Liefer-Punkte 1 bis 3, Gates, Review, Doku): erfüllt. Offen sind nur die Planner-Aufgaben der Closure.** Keine Verifier-only-DoD-Verletzung gefunden.

1. **§6 Risiken — Ausgänge:** Versionssprung 5.4 → 6.x: *eingetreten und getragen* (Messung P2 bis P5, Linie 6.1.x trägt; `vfs`-Befund P6 war das eingetretene Teilrisiko, behoben). Lieferkette Einzelperson: *weiter offen* → Register. Kollision der Konfiguration/`registries.conf`: *eingetreten* (vier Registries), gelöst durch `docker.io`-Setzung. amd64 ungemessen: *weiter offen* bis CI-Lauf nach dem Push. Image-Größe/Bauzeit: *entfallen* (P7 gemessen: kleiner und schneller). `make test-docker`-Laufzeit: *gemessen 789 s*, innerhalb des Limits; Kommentar in `repo.mk` (7 bis 11 Minuten) prüfen. Doppelbesetzung Debian/Archiv: *entfallen* (keine Debian-Pakete). „Neueste veraltet“ (f1): *weiter offen* → Register (Pin-Hebung je Podman-Release, Re-Evaluierungs-Trigger der ADR). Warn-/Sicherheitsfolgen: *entfallen* (Profil unverändert).
2. **§7 Closure-Notiz** mit Steering-Loop-Lerneintrag und den Hinweisen **F-4** (Schlüsselablauf 2027-03-08: Wiedervorlage vor dem Datum, Pin-Hebung mit Schlüsselprüfung), **F-5** (Pins per `--build-arg` überschreibbar), **F-7** (`vfs` ohne Rückbau, `containers.conf` wird ersetzt), **F-8** (Negativprobe: nur `signing_subkey_as_the_pin` prüft den `VALIDSIG`-Vergleich; `unknown_fingerprint` scheitert am Schlüsselabruf). Dazu die Wortlaut-Notiz „Feld 11“ gegen `$12` und die drei Fundstellen-Drifts (Zeilen 629, 510, 1459) in §3 berichtigen.
3. **Register:** Beobachtungen aus `BEO-DOK/plan-aussage-ueber-code` (Fundstellen-Drift nach der letzten Fix-Runde: 1× in diesem Slice), Lieferkette Einzelperson, Pin-Alterung, Schlüsselablauf; Zählerstand vor dem Start sichten (§8).
4. **Roadmap:** Zeile der Welle `v0.11.0` führt den Slice als weiteren Blocker (Eintrag durch den Lead, laut §3 nicht Teil des Plans); nach dem Move nach `done/` Marker und Stub-Zeile nachziehen, `open/README.md` und `in-progress/README.md` entsprechend (Reihenfolge nach `AGENTS.md` §3.3: erst Inhalt, dann reiner `git mv`, `make slice-mv`). <!-- d-check:ignore (Kennungen als Zitat im Report) -->
5. **Cut-Start-Trigger:** [`slice-v1-release-cut-v0.11.0`](../plan/planning/in-progress/slice-v1-release-cut-v0.11.0.md) darf erst nach der Closure starten (Slice §5: DoD vollständig, Review, Verification, Closure-Notiz, ADR Accepted — letzteres erfüllt). Vor dem Cut: Push und CI-Lauf `integration.yml` für die `amd64`-Aussage (bleibt sonst „nicht gemessen“ im Handbuch).

## Host-Zustand am Ende

Arbeitsbaum sauber (nur dieser Report hinzugefügt); keine `sbxe2e-*`-Container oder -Volumes, Wegwerf-Verzeichnis und `vp6-*`-Images entfernt, `alpine:3` gecacht; fremde Container und Images (`probe-cs`, `probe-java`, Compose-Projekte) unberührt; kein `docker image prune`, kein Push, kein Tag.
