# Slice V1: Umstellung der Entwicklung auf nur Devcontainer (`BUILD_MODE=local` als einziger Modus)

**Lifecycle:** Der Zustand dieser Datei ist das Verzeichnis; Wechsel nur durch `git mv` (Modul 5 §Lifecycle).

**Welle:** ohne Welle-Datei (Vorschlag des Projektinhabers, 2026-11-03: „[ein] weiter[er] Slice für die Umstellung auf nur devcontainer (siehe `bau.sh` und `BUILD_MODE=local`) — also nur noch `build_mode=local`").

**Bezug:** [`SPEC-031`](../../../../spec/spezifikation.md) (Regeln des Docker-only-Workflows; der lokale Modus ist dort der Devcontainer-Modus), [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md) (die Folge-ADR mit `Supersedes` zu [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md), `Accepted` — Stand bei Eröffnung des Slices), `scripts/bau.sh` (die lokalen Rezepte), `repo.mk` `BUILD_MODE`-Zweig, `modus-pruefung` (die zwei Modi je `make -n`), `AGENTS.md` §3.1 (Docker-only: der Host braucht nur Docker und GNU make — der Host-Vertrag würde neu), `.github/workflows/` (die CI trug den Docker-Modus bis zum Übergang). **Verantwortlich:** pt9912 (Implementer; Messung CI-im-Devcontainer zuerst, CR je danach; Entscheidungen bejaht 2026-11-03). **Autor:** pt9912 (Rolle Planner, im Auftrag). **Datum:** 2026-11-03.

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
| `scripts/bau.sh` | update | bleibt (der einzige Bau-Pfad); Kopfkommentar nennt die neue Aufruf-Ebene je Entscheidung |
| `repo.mk`, `Makefile` | update | `BUILD_MODE`-Zweige entfallen je Entscheidung 3; Targets rufen `bau.sh` direkt (`Makefile` ohne Diff — es trägt keine Modus-Logik) |
| `Dockerfile`, `scripts/bau.sh` (Kopfkommentare) | update | die Kommentare nennen den Docker-Modus als Default je §3.7 berichtigt (Abgleich-Nachtrag, nicht im ersten Plan-Zug) |
| `.github/workflows/*.yml` | update | CI-Läufe über das Devcontainer-Image je Entscheidung 1 (`ci.yml::gates`); `integration.yml` nur Kommentar je Sockelform |
| `spec/spezifikation.md` | update | [`SPEC-031`](../../../../spec/spezifikation.md) Umformulierung, Beleg-Tabelle; [`SPEC-008`](../../../../spec/spezifikation.md) Folgepflicht (Zweck-Spalte ohne Zweitmodus) |
| `docs/plan/adr/` | neu / update | Folge-ADR 0043 mit Supersedes-Kante zu [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md); `Geschichte`-Zeile und Index des Vorgängers |
| `AGENTS.md`, READMEs | update | Host-Vertrag je Entscheidung 2, Modus-Doku entfällt; `docs/user/benutzerhandbuch.md` ohne Diff (geprüft: die `BUILD_MODE`-Nennung(en) des Handbuchs tragen die Zusage zu generierten Projekten, nicht zum Bau-Pfad des Repos) |
| `harness/README.md`, `docs/maintainer/quality.md`, `CHANGELOG.md` | update | Die `modus-pruefung`-Sensors-Zeile entfällt und der Lauf-Ebene-Absatz folgt je Entscheidung 3; Handbuch-Abschnitte je Ist-Stand; `[Unreleased]`-Eintrag (Nachtrag des ersten Plan-Zugs) |
| `.d-check.yml` | — (kein Diff) | Plan-Korrektur: `modus-pruefung` ist dort nicht eingetragen; die Entfernung trägt die Sensors-Zeile von `harness/README.md` allein |
| `tools/messung/devcontainer-ci/` (Skripte), `docs/messungen/slice-v1-bau-nur-devcontainer/` (Protokoll + Rohdaten), `docs/messungen/README.md` (Bestands-Zeile) | neu / update | Messung M1 (Plan-Risiko 1), vor dem Vertrags-Wortlaut; nur Lesen und Protokollieren, kein Code- oder Vertrags-Diff |
| `docs/plan/planning/in-progress/cr-vorschlag-spec-031-devcontainer-bau.md` | update | Link-Berichtigungen je d-check-Befund (Messbasis-Href, `Supersedes`-Link im Backtick) |

## 4. Trigger

**Start** (`open` → `next`): die drei Entscheidungen sind bejaht (2026-11-03, oben) und der CR-Schnitt der Pin-Inventur [`slice-harness-dcheck-pinsync`](../done/slice-harness-dcheck-pinsync.md) ist geschlossen; `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/`. **Rückführungen:** `in-progress` → `open` (die CI-Messung widerspricht — etwa der Devcontainer-Container schlägt in der CI an dem Punkt, der die Docker-Stage trägt).

## 5. Closure-Trigger

DoD vollständig; `make gates` und die CI-Läufe je am neuen Weg grün; Reports; Closure-Notiz mit Lerneintrag.

## 6. Risiken

1. Die CI-Migration verschiebt die Last auf GitHub-Runner mit anderen Umgebungeigenschaften (Nested-Docker, Cache) — die e2e-Test-Basis (docker-getaggt) braucht im Devcontainer-Modus den verschachtelten Daemon, dessen CI-Fähigkeit ist nicht gemessen. Vorsorge: Messung vor dem Wortlaut ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)) — **Ausgang:** bei Closure.
2. Der Host-Vertrag ist öffentlich (README/Handbuch/AGENTS) — die Umstellung ändert eine Zusage mit Nutzer-Sicht — **Ausgang:** bei Closure.
3. Plan-Aussagen über Code (Workflows, `repo.mk`-Zweige) — §3.10-Abgleich je Runde — **Ausgang:** bei Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Berührt: `scripts/`/`tools/` (WZG), Code-Sub-Areas indirekt über die Gatter (keine Code-Diffs), `DOK`. Alle Greenfield. Evidenz-Risiko: mittel (CI-Läufe im verschachtelten Daemon nicht gemessen)
| tools/harness/tool-pruefung.sh (das Laufzeit-Gate samt Vollstaendigkeits-Form: Pins exakt, make-Floor, Anwesenheit von bash/curl/gh/git/docker-Client) und das Fragment tools/mk/tool-pruefung.mk (per Glob include tools/mk/*.mk gebunden) | neu | die Trager der Fix-Runde (F-1-Form je Freigabe); der Ablage-Waechter-Test bleibt bei tools/harness/ |
| die Fragmente tools/mk/ablage-pruefung.mk (der include-Glob) und repo.mk-Traeger | update | die include-Form der Fix-Runde (keine Doppelnennungen; make -n gates exit 0 je Probelauf) |.