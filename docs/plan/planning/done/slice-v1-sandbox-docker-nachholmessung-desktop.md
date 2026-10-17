# Slice V1: Docker-Desktop-Messung für `nestedRuntime: docker`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle. Das Lastenheft 0.7.0 sagt Docker Desktop für `nestedRuntime: docker` nicht mehr zu
([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)); der Slice trägt keine Closure-Bedingung
einer Welle. Er ist ohne Lieferung geschlossen (§7 `Gegenstand:`); eine Wiederaufnahme der Zusage verlangt einen CR am Lastenheft.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer),
[ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md).
Abnahmekriterium 6 (Nachholmessung) aus dem CR
[`slice-cr-sandbox-docker-runtime`](../done/slice-cr-sandbox-docker-runtime.md) §1, Teil Docker Desktop.

**Berührte Spec-Stellen:** — (Messung und Doku; ändert sich das Ergebnis, folgt ein CR).

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Messreihe aus [ADR-0018](../../../../docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) liegt für Docker Desktop (macOS) vor, und die Doku sagt die Lauffähigkeit erst danach zu.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Podman als Host-Engine und VS Code/Codespaces — [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) sagt Podman als Host-Engine für `docker` nicht zu; Codespaces hat keinen Anlass im CR.
- Änderungen am Lockerungsprofil — ein Ergebnis, das `--privileged` nicht trägt, führt über den Re-Evaluierungs-Trigger von [ADR-0018](../../../../docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md) zu einem neuen CR.

## 2. Definition of Done

- [ ] Integrationstest `TestE2E_SandboxDevcontainer_NestedDocker` (`//go:build docker`) läuft gegen Docker Desktop; Ergebnis (Engine-Version, VM-Typ, Kernel, cgroup) steht in der Geschichte von [ADR-0018](../../../../docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md).
- [ ] Benutzerhandbuch §Grenzen nennt Docker Desktop als geprüft oder mit der gemessenen Einschränkung.
- [ ] Trägt die Messung, nimmt ein CR die Zusage für Docker Desktop in [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) wieder auf (Lastenheft 0.7.0 hat sie gestrichen).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0018-nested-docker-daemon-sandbox-devcontainer.md` (Geschichte) | update | Messergebnis; Status bleibt `Accepted` (Geschichte-Zeile, keine inhaltliche Umschreibung) |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt „Grenzen" |

## 4. Trigger

**Start:** entfällt; der Slice ist ohne Lieferung geschlossen (§7 `Gegenstand:`). Eine Messung unter Docker Desktop beginnt als neuer Slice, sobald ein Nutzerbericht oder ein Host den Re-Evaluierungs-Trigger von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) auslöst.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die dritte Ebene (`docker run --privileged docker:dind`) ist unter Docker Desktop nicht erreichbar — dann CR auf die Ebenen 1–2 schneiden (CR §4).
- `in-progress` → `open` (blockiert): kein Docker-Desktop-Host verfügbar.

## 5. Closure-Trigger

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Docker Desktop startet den Daemon im Container nicht (andere cgroup-/iptables-Lage in der VM) — **Ausgang:** entfallen: Docker Desktop ist seit Lastenheft 0.7.0 nicht zugesagt und der Slice wird nicht gearbeitet; ein Nutzerbericht fällt unter den Re-Evaluierungs-Trigger von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** — (nichts geliefert; der Slice wurde nicht gearbeitet.)
- **Was ging anders als geplant:** Das Lastenheft 0.7.0 nahm Docker Desktop für `nestedRuntime: docker` aus der Zusage; danach blieb der Slice ohne Welle liegen.
- **Gegenstand:** entfallen: Die Zusage für Docker Desktop ist seit Lastenheft 0.7.0 zurückgenommen ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)), und es steht kein Docker-Desktop-Host zur Verfügung. Der Re-Evaluierungs-Trigger von [ADR-0018](../../adr/0018-nested-docker-daemon-sandbox-devcontainer.md) deckt einen Nutzerbericht oder einen späteren Host ab; eine Wiederaufnahme der Zusage verlangt einen CR am Lastenheft (§2).
- **Steering-Loop-Eintrag:** Ein Slice, dessen Voraussetzung dauerhaft fehlt, wird als entfallen geschlossen statt in `open/` gehalten. Gezählt, nicht verkörpert; erstes Auftreten.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung angelegt.
- **Folge-Slices:** keine.
- **Risiken aus §6:** jedes trägt einen Ausgang — siehe §6 (entfallen: Grund).
- **Drei Paarungen:** Anker · Folge-Slice · Register nach dem `git mv` zu prüfen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** `docs/plan/`, `docs/user/` (Greenfield laut `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** kein Beobachtungs-Register
mit Treffern für die berührten Sub-Areas.

Alle berührten Sub-Areas GF.
