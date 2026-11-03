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

- [x] Entcheidungen 1 bis 3 je vom Projektinhaber bejaht (Wortlaut im CR)
- [x] `scripts/bau.sh`-Pfad ist der einzige Bau-Pfad: `repo.mk`-Zweige ohne `BUILD_MODE`, `modus-pruefung` je Entscheidung 3
- [x] Die CI-Workflows laufen die Gatter in der Devcontainer-Umgebung (je Entscheidung 1) je grünem Lauf `ci` und `integration-docker`? — Nachzug je Messbasis
- [x] Vertrags-Nachzug: [`SPEC-031`](../../../../spec/spezifikation.md)-Umformulierung (kein Zweimodus-Vertrag), Folge-ADR `Supersedes`-Verweis zu [`ADR-0037`](../../adr/0037-bauen-im-devcontainer-lokaler-modus.md), [`AGENTS.md`](../../../../AGENTS.md) §3.1 je Entscheidung 2
- [x] Doku je Träger (Handbuch `BUILD_MODE`, READMEs falls betroffen), `CHANGELOG` `[Unreleased]` Eintrag
- [x] Gates, Reports, Closure je Muster der Umsetzungs-Slices

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

### Closure (Planner, 2026-11-03, nach Umsetzung `ff4d8dde`/`e18bb158`/`a8f578a0`/`73e0fda6`, CR-Ortsform `97bebdc1`/`bc8056d3` und Review-/Verifier-Berichtigung `7c048849`/`e5146427`)

**Evidence:** Review-Report [`2026-11-03-slice-v1-bau-nur-devcontainer`](../../../reviews/2026-11-03-slice-v1-bau-nur-devcontainer.md) (1 HIGH/2 MEDIUM/3 LOW/2 INFO je berichtigt) · Verify-Report [`…-verify`](../../../reviews/2026-11-03-slice-v1-bau-nur-devcontainer-verify.md) — **DoD-Zeilen 1, 2, 4 belegt erfüllt; CI real gelaufen:** Run 41676843984 (gates 2 m 22 s inkl. Coldcache-Image-Bau, `tool-pruefung` OK in der CI) und Run 41676844151 (`integration-docker` success 18 m 12 s) — die M1-Netzbau-Grenze entlastet sich am Runner; M3 im eigenen Lauf reproduziert (33 s warm, Coverage 92,20 %); die zwei Bewejs des Laufzeit-Gates je eigener Nachfahrt (Container-Form grün: go 1.27.3/lint 2.14.0/make ≥ 4.3 samt Inventur; Host-Form rot: kein `go` am Host).

**Die neue Nachweis-Form:** der Lauf-Ort ist die Werkzeugkette, nicht der Container — `tool-pruefung` (die Vollständigkeits-Form der Inventur je Projektinhaber-Auftrag) bricht jeden Lauf, dessen Werkzeugkette die Inventur nicht trägt; der Stempel-Hash ist nicht mehr umgebungsgetragen ([`BEO-WZG/tree-hash-collation`](../observations/BEO-WZG/tree-hash-collation/observation.md): `LC_ALL=C`-Pin, Host/Container je byte-gleich — Evidence-Datei ist angelegt); die DIND-Startkante (das Socket-Rennen des ersten Container-Laufs) trägt die Warteschlang-Form (60 s) und als Grenze/Empfehlung im Protokoll — Retry-Form der Läufe.

**Risiken je Ausgang:** 1 CI-Migration — **belegt** (die CI-Läufe je Log oben; der Coldcache-Netzbau fiel am Runner) · 2 der Host-Vertrag — bleibt je [`ADR-0043`](../../adr/0043-nur-devcontainer-bau-pfad.md) (Entscheidung 2) samt der Werkzeugketten-Formulierung · 3 Plan-Aussagen über Code — die §3-Zeilen tragen je Berichten (F-6 halb, V-3 Rest je Fix; V-4 die Fix-Runden-Träger nachgetragen) — beide berichtigt vor der Closure.

**Was ging anders als geplant:** (a) der erste Container-Gates-Lauf endete am fehlenden verschachtelten Socket (Start-Rennen — der zweite Lauf trug in 35 s); (b) der Umsetzungs-Agent wurde mid-flight gestoppt (Arbeitsbaum sauber, vier Commits draußen; die Berichtigung kam vom Hauptagenten); (c) der Hash-Nachweis war umgebungsgetragen (uutils- vs. GNU-sort) — beobachtet und berichtigt je Pin.

**Steering-Loop-Eintrag** (gezählt, nicht verkörpert; kein Eintrag trägt das Feld `liegt in`): *Lehre 1:* die Werkzeugkette trägt, nicht der Lauf-Ort — die künftigen Werkzeuginsel der Inventur folgen der Vollständigkeits-Zeile (je Werkzeug je Form). *Lehre 2:* die Collation-Klasse ist beobachtet und getragen (Pin); `BEO-WZG/tree-hash-collation` trägt den Vorgang (1×, unter der Schwelle). *Neuer Sensor:* keiner. *Benannte Spec-Lücke:* keine.

**Beobachtungs-Register:** *Evidence neu:* [`BEO-WZG/tree-hash-collation`](../observations/BEO-WZG/tree-hash-collation/observation.md) ([`evidence/slice-v1-bau-nur-devcontainer.md`](../observations/BEO-WZG/tree-hash-collation/evidence/slice-v1-bau-nur-devcontainer.md): V-1; berichtigt im Pin). *Nicht gezählt:* [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) — die zustandsmarker-Instanz der Roadmap-Zeile (V-4/INFO) trägt die Berichtigung, keine Zählung; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) — kein Fund.

**Carveouts:** keiner — kein Gate gesenkt, keine Schwelle geändert (die `tool-pruefung`-Aufnahme ist eine Gate-ERWEITERUNG, keine Lockerung). **Folge-Slice:** der docker-cli-Weg (in `open/`) bleibt; keine Cuts versprochen. **Weiter offen:** nichts aus der DoD.

**Drei Paarungen:** *Anker* — kein Steering-Loop-Eintrag trägt ein Zielort-Feld; nichts verkörpert (die ADR-/Inventur-Form tragen die Regeln). *Folge-Slice* — keiner versprochen (docker-cli-Weg ist offene Nachbarkandidatin, Verzeichnis-Lesung). *Register* — die eine neue `BEO-WZG`-Kennung löst in ein Verzeichnis mit der Evidence-Datei dieser Closure (`ls`).

**Verifikation** (Container-Form-Läufe nach Closure-Inhalt, Move und Marker; Stempel je HEAD): `make gates` Exit 0 am 2026-11-03 (`record-gates`-Stempel = HEAD nach dem Nachzug). Der Lauf gilt für den Stand vor dieser Zeile; die Zeile ändert nur Prosa.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Berührt: `scripts/`/`tools/` (WZG), Code-Sub-Areas indirekt über die Gatter (keine Code-Diffs), `DOK`. Alle Greenfield. Evidenz-Risiko: mittel (CI-Läufe im verschachtelten Daemon nicht gemessen)
| tools/harness/tool-pruefung.sh (das Laufzeit-Gate samt Vollstaendigkeits-Form: Pins exakt, make-Floor, Anwesenheit von bash/curl/gh/git/docker-Client) und das Fragment tools/mk/tool-pruefung.mk (per Glob include tools/mk/*.mk gebunden) | neu | die Trager der Fix-Runde (F-1-Form je Freigabe); der Ablage-Waechter-Test bleibt bei tools/harness/ |
| die Fragmente tools/mk/ablage-pruefung.mk (der include-Glob) und repo.mk-Traeger | update | die include-Form der Fix-Runde (keine Doppelnennungen; make -n gates exit 0 je Probelauf) |.