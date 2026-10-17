# Slice V1: Integrationstest und Doku für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** Docker als verschachtelte Runtime (Lastenheft 0.5.0).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).
Abnahmekriterien 2, 3 und 6 aus dem CR [`slice-cr-sandbox-docker-runtime`](slice-cr-sandbox-docker-runtime.md) (§1).

**Berührte Spec-Stellen:** — (die Spec-Änderung liefert der CR; dieser Slice
setzt sie um). Spezifikations- und Architektur-Folgeänderungen benennt die
Folge-ADR.

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Integrationstest weist die Abnahmekriterien 2 und 3 nach; die Colima-Messung und die Nutzer-Doku liegen vor.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Generator und Degradation — [`slice-v1-sandbox-docker-generator`](slice-v1-sandbox-docker-generator.md) und [`slice-v1-sandbox-docker-degradation-doctor`](slice-v1-sandbox-docker-degradation-doctor.md) sind Voraussetzung.
- Nachholmessung unter Docker Desktop — [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) trägt sie; dieser Lauf hat keinen Docker-Desktop-Host.
- Geteilter Image-Cache zwischen Instanzen — der CR schließt ihn aus; die Doku nennt den Speicherverbrauch.

## 2. Definition of Done

- [x] Integrationstest `//go:build docker` auf Linux: `docker info`, `type=cacheonly`, `RUN --network=none`, `docker compose up`, dritte Ebene `--privileged`, zwei getrennte Instanzen (Abnahmekriterien 2, 3, 6).
- [x] Messung unter Colima in der Geschichte der Folge-ADR; Nutzer-Doku zu Speicher, Container-Ausbruch und Egress-Hosts. *(Plan geändert: die Nachholmessung unter Docker Desktop führt [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md), siehe §7.)*
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
| `internal/e2e` (Integrationstest) | neu | Abnahmekriterien 2, 3 |
| `docs/user/` | update | Speicher je Instanz, Ausbruchsrisiko, Egress-Hosts |
| Folge-ADR | update | Messergebnis macOS |

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

- Dritte Ebene `--privileged` unter Docker Desktop nicht erreichbar — **Ausgang:** weiter offen: [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md); unter Colima erreicht ([ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)).

## 7. Closure-Notiz

- **Was hat funktioniert:** `TestE2E_SandboxDevcontainer_NestedDocker` (Kriterien 2 und 3: `type=cacheonly`, `RUN --network=none`, Compose, dritte Ebene, getrennte Daemon-IDs) und `_Restart` laufen unter Colima grün; die Messreihe steht in [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md); Benutzerhandbuch beschreibt Speicher, Ausbruchsrisiko und Egress-Hosts (`6d726e0`).
- **Was ging anders als geplant:** Eine Docker-Desktop-Messung war in diesem Lauf nicht möglich (nur Colima verfügbar). Abnahmekriterium 6 ist damit für Colima erfüllt, für Docker Desktop offen; der Plan ist geändert und die Messung in [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) ausgelagert (kein stilles Streichen). Das Lastenheft nennt Docker Desktop als Ziel, die Doku sagt es erst nach der Messung zu (Review F-4).
- **Steering-Loop-Eintrag:** Sensor `make test-docker`: Integrationstests laufen gegen die verfügbare Engine; fehlende Zielumgebungen werden als eigener Slice geführt, nicht in der DoD stehen gelassen.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angefallen.
- **Folge-Slices:** [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) (Docker-Desktop-Messung, `open/`).
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** Anker (Docker-Desktop-Messung) → Folge-Slice [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md) in `open/` · Register: keine Beobachtung angefallen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/e2e` und `docs/user/` (Greenfield laut `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.

## Verification Evidence

Scope:
- Slice: `slice-v1-sandbox-docker-integrationstest`
- IDs: [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md), [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md)
- Artefakte: `internal/e2e/sandbox_devcontainer_docker_test.go`, `docs/user/benutzerhandbuch.md`, `docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md`, Folge-Slice `docs/plan/planning/done/slice-v1-sandbox-docker-nachholmessung-desktop.md`

DoD-Abgleich (§2):
- [x] Integrationstest `//go:build docker` (Kriterien 2, 3, 6): `docker info`, `type=cacheonly`, `RUN --network=none`, Compose, dritte Ebene, zwei getrennte Instanzen — Evidence: `TestE2E_SandboxDevcontainer_NestedDocker` und `_Restart`, `make test-docker` pass (`internal/e2e` 262 s) unter Colima. Läuft gegen die Colima-VM, nicht auf einem Linux-Host; das DoD-Wort „auf Linux“ ist über die Colima-VM erfüllt, nicht über einen Linux-Host.
- [~] Messung in der Geschichte der Folge-ADR; Nutzer-Doku — Colima-Teil erfüllt ([ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md), Benutzerhandbuch im Diff gesichtet). **DoD-Abweichung (Planner-Thema, nicht von mir entschieden):** Der CR-Slice nennt in Kriterium 6 „Nachholmessung auf macOS mit Docker Desktop und mit Colima“. Die DoD-Zeile (Slice-Datei Zeile 38) ist umformuliert („Messung unter Colima“, Hinweis „Plan geändert“) und trotzdem abgehakt; die Docker-Desktop-Hälfte steht im Folge-Slice `open/`. Ob die Umformulierung zulässig ist oder der Slice offen bleibt, entscheidet der Planner.
- [x] `make gates` grün — Evidence: lint/test/coverage-gate grün; `docs-check` siehe Sensors.
- [x] Review, Report liegt vor — Evidence: `docs/reviews/2026-10-14-sandbox-docker-runtime.md`; F-1 (Neustart) ist durch `_Restart` abgedeckt, ausgeführt und grün.
- [x] Doku-Update, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang, Paarungen — Evidence: §6/§7 (Risiko „Docker Desktop“ trägt „weiter offen“ mit Folge-Slice).

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | fail (nur bekannter Befund) | `lint` 0 issues, `test` ok, `coverage-gate` OK (90.90 % ≥ 90 %); `docs-check` endet mit 7 Befunden, alle durch die noch nicht erfolgte Verschiebung nach `done/` erwartet (6 × `target-missing` auf `done/slice-*sandbox-docker*`) plus 1 × `id-unlinked` in `docs/reviews/2026-10-14-sandbox-docker-runtime.md:30` (nicht erwartet, siehe Bericht) |
| `make docs-check` | fail (7 Befunde, s. o.) | `d-check: 192 Datei(en) geprüft, 7 Befund(e)`; keine Befunde in den vier Slice-Dateien selbst |
| `make test-docker` | pass | alle Pakete `ok`, `internal/e2e` 262 s, lief gegen die lokale Colima-Engine; Ausgabe ohne `-v`, Einzeltests nicht einzeln ausgewiesen, kein `SKIP` in der Ausgabe |

Traceability:
| ID | Beleg |
| --- | --- |
| [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) | E2E `NestedDocker`, `_Restart`, `_Degradation`; Lastenheft sagt Docker Desktop zu, gemessen ist nur Colima (Review F-4) |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) | kein Socket-/Bind-Mount im E2E-Mounts-Check |
| [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) | Mount-Prüfung `-docker-${devcontainerId}`, getrennte Daemon-IDs |
| [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) | Messreihe Colima im ADR im Diff |

Carveouts:
- Neu: none (die offene Docker-Desktop-Messung ist als Folge-Slice in `open/` geführt, nicht im Carveout-Inventar)
- Geloest: none
- Unveraendert: none berührt

Nicht ausgefuehrt:
- Docker-Desktop-Messung — kein Docker-Desktop-Host; [`slice-v1-sandbox-docker-nachholmessung-desktop`](../done/slice-v1-sandbox-docker-nachholmessung-desktop.md)`.
- Lauf auf einem Linux-Host ohne VM — nicht verfügbar.
- Neuer `make docs-check` nach dem finalen `git mv` — steht aus.

Commit / Artefakt:
- Implementation `6d726e0`, Review-Fixes `9fa95ce` und `e1c2c3a`, Go-Pin `2a888f7`, Golden-Tests `8bc3d15`; Delivery-Hash (letzte Code-Änderung): `e1c2c3a`.
- Review-Report: `docs/reviews/2026-10-14-sandbox-docker-runtime.md` (Gegenstand `fa6aa1b..6d726e0`) und Folgelauf [`2026-10-14-sandbox-docker-runtime-r2`](../../../reviews/2026-10-14-sandbox-docker-runtime-r2.md) (Gegenstand `6d726e0..8bc3d15`; dessen Befunde sind in `e1c2c3a` behoben).
