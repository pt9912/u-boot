# Slice V1: Umstellung der Entwicklung auf nur Devcontainer (`BUILD_MODE=local` als einziger Modus)

**Lifecycle:** Der Zustand dieser Datei ist das Verzeichnis; Wechsel nur durch `git mv` (Modul 5 §Lifecycle).

**Welle:** ohne Welle-Datei (Vorschlag des Projektinhabers, 2026-11-03: „[ein] weiter[er] Slice für die Umstellung auf nur devcontainer (siehe `bau.sh` und `BUILD_MODE=local`) — also nur noch `build_mode=local`").

**Bezug:** [`SPEC-031`](../../../../spec/spezifikation.md) (Regeln des Docker-only-Workflows; der lokale Modus ist dort der Devcontainer-Modus), [`AADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) (`Accepted` — der Umzug braucht eine Folge-ADR `Supersedes`), `scripts/bau.sh` (die lokalen Rezepte), `repo.mk` `BUILD_MODE`-Zweig, `modus-pruefung` (die zwei Modi je `make -n`), `AGENTS.md` §3.1 (Docker-only: der Host braucht nur Docker und GNU make — der Host-Vertrag würde neu), `.github/workflows/` (CI läuft heute Docker-Modus). **Verantwortlich:** — (Setzung beim Start). **Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

## 1. Ziel und Abgrenzung

**Ziel:** `BUILD_MODE` entfällt als Betriebsentscheidung — die Rezepte (`bau.sh`) sind der einzige Weg, das Repo zu bauen/linten/testen; der Container (Devcontainer) ist die Bau-Umgebung, und der Modus-Gewirr (`modus-pruefung` mit seinen Kreuzproben) entfällt zugunsten einer Lage.

**Die Entscheidungen des Projektinhabers (2026-11-03, alle drei bejaht):**

1. **CI:** bejaht — die Gatter laufen in der CI in einem Devcontainer-Container (`docker run` auf dem Devcontainer-Image statt der Docker-Stage-Zweige); die Migration ist der eigentliche Träger-Zug.
2. **Host-Vertrag:** bleibt — der Host braucht Docker und GNU make (AGENTS §3.1); die Nutzungs-Sicht verschiebt sich auf den Devcontainer, der Vertrag selbst wird nicht gesenkt.
3. **`modus-pruefung`/`BUILD_MODE`-Zweig:** bejaht entfallen — der Zweig entfällt, `modus-pruefung` entfällt mit ihm (keine Abweisungs-Form).

**Ausdrücklich NICHT in diesem Slice:** Pin-/Werkzeug-Änderungen (die Pin-Inventur läuft separat), das Publish-/Tap-Verhalten der Release-Pipeline (läuft in Docker-Containern, unabhängig vom Entwicklungsmodus; falls die CI-Migration sie berührt, je Entscheidung 1), Doku-Nutzertext-Nachzug gehört zur Umsetzung.

## 2. Definition of Done

- [ ] Entcheidungen 1 bis 3 je vom Projektinhaber bejaht (Wortlaut im CR)
- [ ] `scripts/bau.sh`-Pfad ist der einzige Bau-Pfad: `repo.mk`-Zweige ohne `BUILD_MODE`, `modus-pruefung` je Entscheidung 3
- [ ] Die CI-Workflows laufen die Gatter in der Devcontainer-Umgebung (je Entscheidung 1) je grünem Lauf `ci` und `integration-docker`? — Nachzug je Messbasis
- [ ] Vertrags-Nachzug: [`SPEC-031`](../../../../spec/spezifikation.md)-Umformulierung (kein Zweimodus-Vertrag), Folge-ADR `Supersedes`-Verweis zu [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md), [`AGENTS.md`](../../../../AGENTS.md) §3.1 je Entscheidung 2
- [ ] Doku je Träger (Handbuch `BUILD_MODE`, READMEs falls betroffen), `CHANGELOG` `[Unreleased]` Eintrag
- [ ] Gates, Reports, Closure je Muster der Umsetzungs-Slices

## 3. Plan (vor Code)

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `scripts/bau.sh` | update | bleibt (der einzige Bau-Pfad); Aufrufe je Entscheidung |
| `repo.mk`, `Makefile` | update | `BUILD_MODE`-Zweige entfallen je Entscheidung 3; Targets rufen `bau.sh` direkt |
| `.github/workflows/*.yml` | update | CI-Läufe über das Devcontainer-Image je Entscheidung 1 |
| `spec/spezifikation.md` | update | [`SPEC-031`](../../../../spec/spezifikation.md) Umformulierung, Beleg-Tabelle |
| `docs/plan/adr/` | neu | Folge-ADR mit `Supersedes` zu [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md) |
| `AGENTS.md`, `docs/user/benutzerhandbuch.md`, READMEs | update | Host-Vertrag je Entscheidung 2, Modus-Doku entfällt |
| `.d-check.yml` | update | der `modus-pruefung`-Index-Zeile je Entscheidung 3 |

## 4. Trigger

**Start** (`open` → `next`): die drei Entscheidungen sind bejaht (2026-11-03, oben) und der CR-Schnitt der Pin-Inventur [`slice-harness-dcheck-pinsync`](../in-progress/slice-harness-dcheck-pinsync.md) ist geschlossen; `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`. **Rückführungen:** `in-progress` → `open` (die CI-Messung widerspricht — etwa der Devcontainer-Container schlägt in der CI an dem Punkt, der die Docker-Stage trägt).

## 5. Closure-Trigger

DoD vollständig; `make gates` und die CI-Läufe je am neuen Weg grün; Reports; Closure-Notiz mit Lerneintrag.

## 6. Risiken

1. Die CI-Migration verschiebt die Last auf GitHub-Runner mit anderen Umgebungeigenschaften (Nested-Docker, Cache) — die e2e-Test-Basis (docker-getaggt) braucht im Devcontainer-Modus den verschachtelten Daemon, dessen CI-Fähigkeit ist nicht gemessen. Vorsorge: Messung vor dem Wortlaut ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) — **Ausgang:** bei Closure.
2. Der Host-Vertrag ist öffentlich (README/Handbuch/AGENTS) — die Umstellung ändert eine Zusage mit Nutzer-Sicht — **Ausgang:** bei Closure.
3. Plan-Aussagen über Code (Workflows, `repo.mk`-Zweige) — §3.10-Abgleich je Runde — **Ausgang:** bei Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Berührt: `scripts/`/`tools/` (WZG), Code-Sub-Areas indirekt über die Gatter (keine Code-Diffs), `DOK`. Alle Greenfield. Evidenz-Risiko: mittel (CI-Läufe im verschachtelten Daemon nicht gemessen).