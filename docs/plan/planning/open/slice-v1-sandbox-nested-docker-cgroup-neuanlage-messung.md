# Slice V1: Messung nach der Neuanlage des Devcontainers — Zustand der Wurzel-cgroup, `docker exec` und neue Terminals nach dem Umbau, Prozesse der Erweiterung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Folge-Slice von [`slice-v1-sandbox-nested-docker-cgroup`](../done/slice-v1-sandbox-nested-docker-cgroup.md), das mit eingeschränktem Gegenstand geschlossen wurde. Entscheidung des Projektinhabers 2026-11-04 (Wortlaut: „dazu muss ich diese session beenden -- deshalb machen wir das erst nach dem Release“): Die Neuanlage des Devcontainers und die Messung folgen **nach dem Release `v0.23.0`**. Eine Closure-Bedingung, die von der DoD dieses Slice verschieden ist, gibt es nicht.

**Bezug:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Container-Runtime im Sandbox-Devcontainer), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Dateien); Entscheidung: [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) (Abschnitt Re-Evaluierungs-Trigger, erster Punkt).

**Berührte Spec-Stellen:** [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) (Grenze, Beleg-Tabelle), [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) (Abschnitt *Nicht zugesagt*); je nach Befund eine Folge-ADR mit `Supersedes` zu [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-11-04.

---

## 1. Ziel und Abgrenzung

**Ziel:** Nach einer Neuanlage des Devcontainers dieses Repos (`make dev-up RECREATE=1`) sind die Risiken 2 bis 4 aus §6 von [`slice-v1-sandbox-nested-docker-cgroup`](../done/slice-v1-sandbox-nested-docker-cgroup.md) gemessen und im Protokoll belegt: (2) die Wurzel-cgroup ist nach der Neuanlage umgebaut, (3) `docker exec` und ein neues Terminal der Erweiterung nehmen nach dem Umbau Prozesse in der Wurzel auf, (4) die Erweiterung hängt sich nach dem Umbau an und ihre Prozesse laufen. Danach laufen die zwei Tests `TestE2E_Resources_DefaultProfileLimitsReachTheContainer` und `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` (`internal/e2e/resources_docker_test.go`) grün, oder ihr Befund steht im Protokoll.

**Herkunft.** Der Vorgänger lieferte den Umbau (ADR, [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag), Template, Tests, Docker-Test im Wegwerf-Container, eigenes Projekt regeneriert). Die Neuanlage des Entwicklungs-Devcontainers beendet die Sitzung, in der der Projektinhaber arbeitet; er hat sie auf die Zeit nach dem Release `v0.23.0` gelegt (siehe Kopf).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Produktcode ohne Befund:** Der Slice misst. Tritt eines der Risiken ein, ist die Folge eine Folge-ADR mit `Supersedes` und ein eigener Code-Slice; Code im selben Slice würde Messung und Änderung vermischen.
- **Keine Änderung von [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) in place:** Sie ist nach `Accepted` immutable ([`AGENTS.md`](../../../../AGENTS.md) §3.5).
- **Keine Zusage für andere Hosts und Engines:** gemessen wird ein Container; die Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) bleibt Gegenstand von [`slice-v1-devcontainer-ressourcenlimits-messung`](slice-v1-devcontainer-ressourcenlimits-messung.md).
- **Kein Test des Umbaus in der Zustandsform `domain threaded`:** der Docker-Test des Vorgängers stellt sie nicht her; eine Probe dafür wäre eine eigene Messung (Register: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)).

## 2. Definition of Done

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [ ] **Neuanlage-Messung.** Messprotokoll unter `docs/messungen/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung/` (Skript unter `tools/messung/sandbox-nested-docker-cgroup/`, Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)) zu den Risiken 2 bis 4: Zustand der Wurzel-cgroup nach der Neuanlage (`cgroup.type`, `subtree_control`), `docker exec` und neues Terminal der Erweiterung nach dem Umbau, Anhängen der Erweiterung und ihre Prozesse; je Zelle mindestens zwei Läufe.
- [ ] **Tests und Entscheidung.** `make test-docker` im neu angelegten Container: die zwei `TestE2E_Resources_*` grün, oder der Befund je Test im Protokoll. Tritt eines der Risiken ein (Re-Evaluierungs-Trigger von [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md)), liegt eine Folge-ADR mit `Supersedes` zu ihr vor (`Accepted`, Index nachgezogen); sonst steht im Protokoll, dass keiner eingetreten ist.
- [ ] **Doku.** Die Stellen, die die Messung als ausstehend nennen, tragen den Ist-Stand: [`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag) (Grenze), `docs/maintainer/devcontainer-arbeiten.md` (Zeile zur Neuanlage), `CHANGELOG.md` (Abschnitt *Not promised* des Eintrags zum Umbau); Nutzer-Doku ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6), nur das Gemessene.
- [ ] `make gates` grün im Devcontainer-Container. *(nicht gezählt — Gate-Lauf)*
- [ ] Review, Verifikation, Closure-Notiz, Register, Risiko-Ausgänge, Drei-Paarungen. *(nicht gezählt — Closure-Pflichten)*

## 3. Plan (vor Code)

Fundstellen am 2026-11-04 gelesen; der Implementer ergänzt die Tabelle im ersten Lauf und gleicht sie nach der letzten Fix-Runde gegen `git diff --stat` ab ([`AGENTS.md`](../../../../AGENTS.md) §3.10).

| Datei / Komponente | Änderungs-Art | Begründung / Fundstelle |
|---|---|---|
| `tools/messung/sandbox-nested-docker-cgroup/` | neu | Skript der Messung (Zustand der Wurzel, `docker exec`, Prozesse); Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) |
| `docs/messungen/slice-v1-sandbox-nested-docker-cgroup-neuanlage-messung/` | neu | Messprotokoll (`README.md`, Rohdaten) |
| `internal/e2e/resources_docker_test.go` | prüfen | `TestE2E_Resources_DefaultProfileLimitsReachTheContainer` (Zeile 41), `TestE2E_Resources_SandboxNestedDockerLimitsReachTheContainer` (Zeile 66); Lauf nach der Neuanlage, keine Änderung erwartet |
| `docs/plan/adr/` | bedingt neu | Folge-ADR mit `Supersedes` zu [`ADR-0045`](../../adr/0045-cgroup-umbau-im-startscript-vor-dockerd.md) samt Index, nur wenn ein Risiko eintritt |
| `spec/spezifikation.md` ([`SPEC-056`](../../../../spec/spezifikation.md#spec-056--cgroup-umbau-im-startscript-docker-initsh-bedingung-schritte-und-fehlschlag), [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs)), `docs/maintainer/devcontainer-arbeiten.md`, `CHANGELOG.md` | update | Stellen, die die Messung als ausstehend nennen (Stand: [`SPEC-053`](../../../../spec/spezifikation.md#spec-053--grenzen-für-arbeitsspeicher-und-cpu-des-devcontainers-schlüssel-werte-und-runargs) *Nicht zugesagt*, Maintainer-Doku Zeile zur Neuanlage) |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Release `v0.23.0` geschnitten **und** der Projektinhaber löst die Neuanlage des Devcontainers mit `make dev-up RECREATE=1` aus (sie beendet seine Sitzung); kein anderer Slice in `in-progress/` (WIP-Limit 1). Ohne beides bleibt der Slice in `open/`.

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Die Neuanlage gelingt nicht oder der Container startet nicht; die Messung ruht, bis der Projektinhaber einen neuen Container anlegt.
- `in-progress` → `next` (zu groß): Tritt ein Risiko ein, wandert die Änderung (Folge-ADR, Code) in einen eigenen Slice; dieser Slice schließt mit dem Befund.

## 5. Closure-Trigger

DoD vollständig; Messprotokoll mit den Risiken 2 bis 4 und je zwei Läufen; `make test-docker` im neu angelegten Container mit den zwei `TestE2E_Resources_*` (grün oder Befund im Protokoll); Entscheidung zur Folge-ADR dokumentiert; Review und Verifikation durch andere Rollen; Closure-Notiz; `git mv` nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

1. **Die Wurzel ist nach der Neuanlage nicht umgebaut** (etwas fügt vor `postStartCommand` Prozesse in die Wurzel ein, oder der Hook läuft nicht). Gelesen ist der Startort, gemessen ist nur der Wegwerf-Container des Docker-Tests. — **Ausgang:** offen (Risiko 2 des Vorgängers) → Messung dieses Slice.
2. **`docker exec` oder ein neues Terminal der Erweiterung scheitert nach dem Umbau mit `EBUSY`** (Wurzel mit gesetztem `subtree_control`). Der Docker-Test des Vorgängers prüft ein `docker exec` in einer Engine, nicht die Erweiterung. — **Ausgang:** offen (Risiko 3 des Vorgängers) → Messung dieses Slice.
3. **Die Erweiterung fällt aus, weil der Umbau ihre Prozesse verschiebt.** — **Ausgang:** offen (Risiko 4 des Vorgängers) → Messung dieses Slice.
4. **Die Messung ist ein Container, ein Host.** Ein Befund trägt nicht für andere Hosts oder Engines. — **Ausgang:** offen → die Zusage-Menge bleibt bei [`slice-v1-devcontainer-ressourcenlimits-messung`](slice-v1-devcontainer-ressourcenlimits-messung.md).

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf, Start-Trigger in §4).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `docs/plan/`, `docs/messungen/` (DOK), `tools/` (WZG) und `internal/e2e` (E2E, nur Lauf); alle sind in [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) deklariert, keine neue Kandidatin.

**Vorgelagert — offene Beobachtungen sichten:** Register gesichtet 2026-11-04; Treffer: [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10; der Zähler ist die Zahl der Evidence-Dateien) und [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (offen über der Schwelle, kein Ausgang zugewiesen). Beide betreffen den Wortlaut der Doku-Änderung dieses Slice: jede Aussage nennt nur das Gemessene und die Messgrenze.

**Modus-Begründungsblock:** entfällt — alle berührten Sub-Areas GF.
