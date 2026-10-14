# Slice V1: Degradation und Doctor für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** Docker als verschachtelte Runtime (Lastenheft 0.5.0).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).
Abnahmekriterien 5 und 7 aus dem CR [`slice-cr-sandbox-docker-runtime`](slice-cr-sandbox-docker-runtime.md) (§1).

**Berührte Spec-Stellen:** — (die Spec-Änderung liefert der CR; dieser Slice
setzt sie um). Spezifikations- und Architektur-Folgeänderungen benennt die
Folge-ADR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein nicht betriebsbereiter Docker-Daemon führt bei `warn` zu Warnung ohne Daemon und bei `fail` zu Exit-Code `11`; `u-boot doctor` weist die Lockerungen aus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Generator-Ausgabe — [`slice-v1-sandbox-docker-generator`](slice-v1-sandbox-docker-generator.md) liefert sie und ist Voraussetzung.
- Podman als Ersatz — der Vertrag verbietet stillschweigenden Ersatz.
- Integrationstest — [`slice-v1-sandbox-docker-integrationstest`](slice-v1-sandbox-docker-integrationstest.md).

## 2. Definition of Done

- [x] `onUnavailable`: Daemon nicht betriebsbereit → `warn` ohne Daemon, `fail` Exit-Code `11` (Abnahmekriterium 5).
- [x] `u-boot doctor`: Warnung ohne Sandbox-Profil; Warnung zur unterlaufenen Egress-Leitplanke (Abnahmekriterium 7).
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
| Config-Validierung, Fehlerklassifikation (`internal/`) | update | Exit-Code `11` laut [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) |
| Doctor-Checks (`internal/`) | update | Abnahmekriterium 7 |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.5.0 und die Folge-ADR des CR [`slice-cr-sandbox-docker-runtime`](slice-cr-sandbox-docker-runtime.md) sind angenommen; [`slice-v1-sandbox-docker-generator`](slice-v1-sandbox-docker-generator.md) ist geliefert.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet
  3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): die Folge-ADR ändert den Zuschnitt oder
  ein Vorgänger-Slice verfehlt seine DoD.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Bereitschaftsprüfung mit Frist hängt an der Folge-ADR — **Ausgang:** entfallen: [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) legt 60 s und `docker info` fest.

## 7. Closure-Notiz

- **Was hat funktioniert:** Degradation (`warn` Exit 0 mit Warnung, `fail` Exit 11) und die beiden `doctor`-Warnungen sind in Unit- und Integrationstests gepinnt (`TestE2E_SandboxDevcontainer_NestedDocker_Degradation`, `TestDoctor_SandboxRuntime`/`_SandboxEgress`; Delivery `6d726e0`).
- **Was ging anders als geplant:** Das 60-Sekunden-Limit war nicht hart (F-5 im Review [`docs/reviews/2026-10-14-sandbox-docker-runtime.md`](../../../reviews/2026-10-14-sandbox-docker-runtime.md)); `docker info` läuft jetzt mit `timeout 5` (`9fa95ce`).
- **Steering-Loop-Eintrag:** Sensor `make test-docker` ergänzt: Degradationspfade des Daemon-Startskripts mit und ohne `--privileged`.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-integrationstest`](slice-v1-sandbox-docker-integrationstest.md) (im selben Lauf geschlossen).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../open/slice-v1-sandbox-docker-nachholmessung-desktop.md) in `open/` · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/driving/cli` und `internal/hexagon/application` (Greenfield laut `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Scope:
- Slice: `slice-v1-sandbox-docker-degradation-doctor`
- IDs: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
- Artefakte: `templates/devcontainer/docker-init.sh.tmpl`, `internal/hexagon/application/doctor_sandbox.go`, `doctor_sandbox_test.go`, `internal/e2e/sandbox_devcontainer_docker_test.go`

DoD-Abgleich (§2):
- [x] `onUnavailable`: `warn` ohne Daemon, `fail` Exit 11 (Kriterium 5) — Evidence: `TestE2E_SandboxDevcontainer_NestedDocker_Degradation` (Start ohne `--privileged`: Exit 0 bzw. 11, Meldung „no daemon available“, kein Podman); `..._SandboxDocker_FailPolicy` pinnt die Policy im Skript. Der 60-s-Timeout-Pfad (Daemon startet, wird aber nicht bereit) ist nicht getestet, nur im Skript gelesen.
- [x] `doctor`: Warnung ohne Sandbox-Profil; Warnung zur unterlaufenen Egress-Leitplanke (Kriterium 7) — Evidence: `doctor_sandbox_test.go` Fälle „docker without sandbox profile“ (Warn, „no effect“) und „on with nested docker“ (Warn, „bypass the restriction“); `docker-init.sh` ist in `credentialScanFiles` aufgenommen.
- [x] `make gates` grün — Evidence: lint/test/coverage-gate grün; `docs-check` siehe Sensors.
- [x] Review, Report liegt vor — Evidence: `docs/reviews/2026-10-14-sandbox-docker-runtime.md` (F-5 betraf dieses Slice; behoben in `9fa95ce` per `timeout 5` auf `docker info`, im Template gesichtet, durch keinen eigenen Test gepinnt).
- [x] Doku-Update, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang, Paarungen — Evidence: §6/§7.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | fail (nur bekannter Befund) | `lint` 0 issues, `test` ok, `coverage-gate` OK (90.90 % ≥ 90 %); `docs-check` endet mit 7 Befunden, alle durch die noch nicht erfolgte Verschiebung nach `done/` erwartet (6 × `target-missing` auf `done/slice-*sandbox-docker*`) plus 1 × `id-unlinked` in `docs/reviews/2026-10-14-sandbox-docker-runtime.md:30` (nicht erwartet, siehe Bericht) |
| `make docs-check` | fail (7 Befunde, s. o.) | `d-check: 192 Datei(en) geprüft, 7 Befund(e)`; keine Befunde in den vier Slice-Dateien selbst |
| `make test-docker` | pass | alle Pakete `ok`, `internal/e2e` 262 s, lief gegen die lokale Colima-Engine; Ausgabe ohne `-v`, Einzeltests nicht einzeln ausgewiesen, kein `SKIP` in der Ausgabe |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | E2E `..._NestedDocker_Degradation`, `TestDoctor_SandboxRuntime` |
| [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) | `TestDoctor_SandboxEgress` (docker + Egress) |
| [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | Exit 11 im Skript (`fail`) und im E2E-Test; ungültige Werte laufen über den bestehenden Exit-10-Pfad (laut Review, nicht eigens neu geprüft) |
| [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | 60 s / `docker info` im Template |

Carveouts:
- Neu: none
- Geloest: none
- Unveraendert: none berührt

Nicht ausgefuehrt:
- Test des 60-s-Timeout-Pfads — nicht vorhanden, nicht von mir ergänzt (Verifier ändert keinen Code).
- Neuer `make docs-check` nach dem finalen `git mv` — steht aus.

Commit / Artefakt:
- Implementation `6d726e0`, Review-Fixes `9fa95ce` und `e1c2c3a`, Go-Pin `2a888f7`, Golden-Tests `8bc3d15`; Delivery-Hash (letzte Code-Änderung): `e1c2c3a`.
- Review-Report: `docs/reviews/2026-10-14-sandbox-docker-runtime.md` (Gegenstand `fa6aa1b..6d726e0`) und Folgelauf [`2026-10-14-sandbox-docker-runtime-r2`](../../../reviews/2026-10-14-sandbox-docker-runtime-r2.md) (Gegenstand `6d726e0..8bc3d15`; dessen Befunde sind in `e1c2c3a` behoben).
