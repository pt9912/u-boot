# Slice V1: Sandbox-Generator für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** Docker als verschachtelte Runtime (Lastenheft 0.5.0).

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil),
[ADR-0014](../../adr/0014-nested-podman-sandbox-devcontainer.md),
[ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md).
Abnahmekriterien 1 und 4 aus dem CR [`slice-cr-sandbox-docker-runtime`](../next/slice-cr-sandbox-docker-runtime.md) (§1).

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

- [ ] `--privileged` in `runArgs` mit Warnung in der Befehlsausgabe; Volume `<projekt>-docker-${devcontainerId}` nach [ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md) (Abnahmekriterium 1).
- [ ] Golden Cases: `none` und `podman` unverändert, nur `docker` erzeugt `--privileged` (Abnahmekriterium 4).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review.
- [ ] Doku-Update, falls ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Pfad-Kandidaten; der Implementer-Agent schärft die Liste nach der Folge-ADR.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Generator-Templates und Config-Validierung (`internal/`) | update | Wert `docker`, `runArgs`, Volume, Daemon-Start nach Folge-ADR |
| Golden-Cases | neu / update | Abnahmekriterium 4 |

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.5.0 und die Folge-ADR des CR [`slice-cr-sandbox-docker-runtime`](../next/slice-cr-sandbox-docker-runtime.md) sind angenommen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice überschreitet
  3 Liefer-Punkte.
- `in-progress` → `open` (blockiert): die Folge-ADR ändert den Zuschnitt oder
  ein Vorgänger-Slice verfehlt seine DoD.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Daemon-Herkunft und -Start hängen an der Folge-ADR — **Ausgang:** weiter offen: Folge-ADR des CR.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** siehe §6
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `internal/adapter/**` und `internal/hexagon/**` (Greenfield laut `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.
