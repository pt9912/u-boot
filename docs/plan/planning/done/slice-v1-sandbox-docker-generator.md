# Slice V1: Sandbox-Generator für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** Docker als verschachtelte Runtime (Lastenheft 0.5.0).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).
Abnahmekriterien 1 und 4 aus dem CR [`slice-cr-sandbox-docker-runtime`](slice-cr-sandbox-docker-runtime.md) (§1).

**Berührte Spec-Stellen:** — (die Spec-Änderung liefert der CR; dieser Slice
setzt sie um). Spezifikations- und Architektur-Folgeänderungen benennt die
Folge-ADR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** `u-boot generate devcontainer` erzeugt bei `nestedRuntime: docker` ein Sandbox-Devcontainer-Artefakt mit eigenem Docker-Daemon.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Degradation und `doctor`-Warnungen — [`slice-v1-sandbox-docker-degradation-doctor`](slice-v1-sandbox-docker-degradation-doctor.md) trägt sie; der Generator liefert zuerst das Artefakt.
- Integrationstest und macOS-Messung — [`slice-v1-sandbox-docker-integrationstest`](slice-v1-sandbox-docker-integrationstest.md) trägt sie; der Nachweis braucht das fertige Artefakt.
- Änderungen an `podman` und `none` — der CR schließt sie aus.

## 2. Definition of Done

- [x] `--privileged` in `runArgs` mit Warnung in der Befehlsausgabe; Volume `<projekt>-docker-${devcontainerId}` nach [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) (Abnahmekriterium 1).
- [x] Golden Cases: `none` und `podman` unverändert, nur `docker` erzeugt `--privileged` (Abnahmekriterium 4).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Pfad-Kandidaten; der Implementer-Agent schärft die Liste nach der Folge-ADR.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Generator-Templates und Config-Validierung (`internal/`) | update | Wert `docker`, `runArgs`, Volume, Daemon-Start nach Folge-ADR |
| Golden-Cases | neu / update | Abnahmekriterium 4 |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.5.0 und die Folge-ADR des CR [`slice-cr-sandbox-docker-runtime`](slice-cr-sandbox-docker-runtime.md) sind angenommen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet
  3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): die Folge-ADR ändert den Zuschnitt oder
  ein Vorgänger-Slice verfehlt seine DoD.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Daemon-Herkunft und -Start hängen an der Folge-ADR — **Ausgang:** entfallen: [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) ist angenommen (Debian-Pakete, `postStartCommand`, `--privileged`).

## 7. Closure-Notiz

- **Was hat funktioniert:** Generator, Templates (`docker-init.sh`, Dockerfile, `devcontainer.json`) und Golden Cases (`TestGenerateDevcontainer_SandboxDocker*`) liefen im ersten Anlauf grün; der Integrationstest baute das erzeugte Dockerfile unverändert (Delivery `6d726e0`).
- **Was ging anders als geplant:** Der Review ([`docs/reviews/2026-10-14-sandbox-docker-runtime.md`](../../../reviews/2026-10-14-sandbox-docker-runtime.md)) fand im Neustart-Pfad des Startskripts einen Fehler (F-1: PID-/Socket-Reste); behoben in `9fa95ce` samt Restart-Test.
- **Steering-Loop-Eintrag:** Sensor `make test-docker` deckt jetzt auch den Neustart des Daemon-Startskripts ab; Skripte, die Hintergrunddienste starten, brauchen einen Neustart-Fall im Integrationstest.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-degradation-doctor`](slice-v1-sandbox-docker-degradation-doctor.md), [`slice-v1-sandbox-docker-integrationstest`](slice-v1-sandbox-docker-integrationstest.md) (beide im selben Lauf geschlossen).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) in `open/` · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/**` und `internal/hexagon/**` (Greenfield laut `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Scope:
- Slice: `slice-v1-sandbox-docker-generator`
- IDs: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md), [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md), [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
- Artefakte: `internal/hexagon/application/templates/devcontainer/{docker-init.sh,Dockerfile,devcontainer.json}.tmpl`, `devcontainer_sandbox.go`, `devcontainer_egress.go`, `generate.go`, `templates.go`, `internal/hexagon/domain/{sandbox,configpath}.go`, `generate_sandbox_test.go`, `harness/replay.md`

DoD-Abgleich (§2):
- [x] `--privileged` in `runArgs` mit Warnung, Volume `<projekt>-docker-${devcontainerId}` (Kriterium 1) — Evidence: `TestGenerateDevcontainer_SandboxDocker`.
- [x] Golden Cases: `none`/`podman` unverändert, nur `docker` erzeugt `--privileged` (Kriterium 4) — Evidence: `TestGenerateDevcontainer_SandboxDocker_NotGenerated`; Podman-Goldens im Diff nicht verändert.
- [x] `make gates` grün — Evidence: lint/test/coverage-gate grün; `docs-check` nur mit den erwarteten `target-missing` sowie einem `id-unlinked` im Review-Report (siehe Sensors).
- [x] Review, Report liegt vor — Evidence: `docs/reviews/2026-10-14-sandbox-docker-runtime.md`.
- [x] Doku-Update — Evidence: Diff `CHANGELOG.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/cli-json-contract.md`, `harness/replay.md`, Spec.
- [x] Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang, Paarungen — Evidence: §6/§7.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | fail (nur bekannter Befund) | `lint` 0 issues, `test` ok, `coverage-gate` OK (90.90 % ≥ 90 %); `docs-check` endet mit 7 Befunden, alle durch die noch nicht erfolgte Verschiebung nach `done/` erwartet (6 × `target-missing` auf `done/slice-*sandbox-docker*`) plus 1 × `id-unlinked` in `docs/reviews/2026-10-14-sandbox-docker-runtime.md:30` (nicht erwartet, siehe Bericht) |
| `make docs-check` | fail (7 Befunde, s. o.) | `d-check: 192 Datei(en) geprüft, 7 Befund(e)`; keine Befunde in den vier Slice-Dateien selbst |
| `make test-docker` | pass | alle Pakete `ok`, `internal/e2e` 262 s, lief gegen die lokale Colima-Engine; Ausgabe ohne `-v`, Einzeltests nicht einzeln ausgewiesen, kein `SKIP` in der Ausgabe |

Replay / Golden:
- Affected generators: `generate devcontainer` / `init --devcontainer --sandbox` mit `devcontainer.sandbox.nestedRuntime: docker` (neu: `docker-init.sh`, Dockerfile-Zweig, `devcontainer.json`-Mounts-Zweig, `runArgs`, `postStartCommand`); `none`/`podman` nur durch Template-Verzweigung berührt.
- Golden cases added (In-Memory, `internal/hexagon/application/generate_sandbox_test.go`): `TestGenerateDevcontainer_SandboxDocker`, `_SandboxDocker_FailPolicy`, `_SandboxDocker_NotGenerated`, `_SandboxEgressWithDocker`; im Arbeitsbaum zusätzlich (uncommittet, im Lauf von `make gates`/`make test-docker` enthalten): `_SandboxDocker_ScriptConflict`, `_SandboxDocker_StaleBlockKeepsUserContent`.
- Golden cases updated: `harness/replay.md` (neue Tabellenzeile für `docker`); bestehende Goldens `config_sandbox_test.go`, `initproject_test.go`, `doctor_sandbox_test.go`, `sandbox_test.go` laut Diff nur um den Wert `docker` ergänzt.
- Golden cases replayed: gesamtes Paket `internal/hexagon/application` über `make test` und `make test-docker` (ok).
- Intentional output changes: nur bei `nestedRuntime: docker` (neue Datei `docker-init.sh`, `--privileged`, Mounts, Dockerfile-Block, `postStartCommand`); `none`/`podman` unverändert (nicht byte-identisch gegen die Vorversion verglichen, nur Abwesenheit der Docker-Artefakte und unveränderte Podman-Goldens belegt).
- Replay Rules für `docker-init.sh`: Fresh-State — `..._SandboxDocker` (Dateien und Inhalt; die Action `Created` wird nicht explizit assertet); Double-Run-NoOp — `..._SandboxDocker` (zweiter Lauf `NoOp`); stale block + User-Content — `..._StaleBlockKeepsUserContent`; kein Partial-Write — `..._ScriptConflict` (prüft nur, dass `devcontainer.json` nicht entsteht; Dockerfile und weitere Dateien werden nicht geprüft). Die beiden letzten Tests sind noch nicht committet.

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | `TestGenerateDevcontainer_SandboxDocker*`, `_SandboxEgressWithDocker` |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | Egress-/Podman-Goldens unverändert grün |
| [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) | Volume-Namen `<projekt>-docker-${devcontainerId}` im Golden |
| [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | Debian-Pakete, `postStartCommand`, `--privileged` im Golden; E2E `NestedDocker` |

Carveouts:
- Neu: none
- Geloest: none
- Unveraendert: none berührt

Nicht ausgefuehrt:
- `make ci`/`govulncheck`/`image-scan` — nicht Teil des Handoffs.
- Neuer `make docs-check` nach dem finalen `git mv` — steht aus.

Commit / Artefakt:
- Implementation `6d726e0`, Review-Fixes `9fa95ce` und `e1c2c3a`, Go-Pin `2a888f7`, Golden-Tests `8bc3d15`; Delivery-Hash (letzte Code-Änderung): `e1c2c3a`.
- Review-Report: `docs/reviews/2026-10-14-sandbox-docker-runtime.md` (Gegenstand `fa6aa1b..6d726e0`) und Folgelauf [`2026-10-14-sandbox-docker-runtime-r2`](../../../reviews/2026-10-14-sandbox-docker-runtime-r2.md) (Gegenstand `6d726e0..8bc3d15`; dessen Befunde sind in `e1c2c3a` behoben).
