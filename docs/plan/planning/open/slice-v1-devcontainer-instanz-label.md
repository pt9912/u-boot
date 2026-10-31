# Slice V1: Mehrere Devcontainer je Projekt — ein Instanzname an `up`, `down` und `shell` über das Id-Label der Dev-Container-CLI

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Es gibt keine Closure-Bedingung, die von der DoD dieses Slice verschieden ist. Anlass: Der Projektinhaber will mehrere Agenten getrennt arbeiten lassen und statt Git-Worktrees mehrere Devcontainer desselben Projekts starten (Gespräch am 2026-10-31).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, `down` und `shell`), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Sandbox-Profil: Workspace im benannten Volume je Instanz), [ADR-0028](../../adr/0028-devcontainer-start-ohne-editor-hilfs-image.md) und [ADR-0031](../../adr/0031-devcontainer-shell-ohne-editor.md) (aktive ADRs zu `up` und `shell`; im Volltext nicht gelesen).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz *Aufruf*: `u-boot devcontainer up [--recreate]`, `down`, `shell`).

**Verantwortlich:** —

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Projekt kann mehrere Devcontainer gleichzeitig haben, jeden mit eigenem Workspace-Volume, und `u-boot devcontainer up`, `down` und `shell` wählen den Container über einen Instanznamen; ohne Angabe verhalten sich alle drei wie heute.

**Befund** (gelesen am 2026-10-31, Fundstellen je Punkt; [`AGENTS.md`](../../../../AGENTS.md) §3.10):

- **Ein Container je Projektordner.** `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`, Zeilen 284 bis 303) ruft die Dev-Container-CLI als `devcontainer up --workspace-folder <Ordner>`, bei `--recreate` mit `--remove-existing-container`, und gibt kein `--id-label` mit. `FindContainer` (`internal/adapter/driven/docker/inspector.go`, Filter in Zeile 86, Konstante `devcontainer.local_folder` in Zeile 21) sucht den Container allein über dieses Label und liefert bei mehreren Treffern den zuletzt angelegten (Port `DevcontainerInspector`, `internal/hexagon/port/driven/devcontainer_inspector.go`, Zeilen 44 bis 49).
- **Die CLI kann mehr.** Die Dev-Container-CLI 0.89.0 im Hilfs-Image von `u-boot` kennt `devcontainer up --id-label name=value`: Die Labels werden am Container gesetzt und dienen der Suche nach einem bestehenden Container; ohne Angabe leitet die CLI sie aus dem Workspace-Ordner ab (`devcontainer up --help`, gemessen am 2026-10-31). Ob die Namen der Volumes (`${devcontainerId}` in `.devcontainer/devcontainer.json`) mit dem Id-Label mitwechseln, ist aus der Spezifikation der CLI angenommen und **nicht gemessen**.
- **`u-boot` reicht nichts durch.** `u-boot devcontainer up` kennt nur `--recreate` neben den globalen Flags (`u-boot devcontainer up --help`; `internal/adapter/driving/cli/devcontainer.go`, Flag in Zeile 164, `down` ab Zeile 173, `shell` ab Zeile 209).
- **Die Suche ist an drei Stellen verdrahtet:** `internal/hexagon/application/devcontainer_start.go` (Zeilen 297 und 429) und `devcontainer_fallback.go` (Zeile 155) rufen `FindContainer(ctx, <Ordner>)`. Das Beobachtungs-Register führt den Mangel bereits: [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (ein Beleg, unter der Schwelle; Trigger: „Messung mit zwei Containern am selben Label“ und „eine Änderung am Port `DevcontainerInspector`“ — beides tritt mit diesem Slice ein).
- **Zusage heute:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) beschreibt `up`, `down` und `shell` für den Devcontainer „des Projekts“, ohne Instanz; zugesagt auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine.

**Wege** (Entscheidung des Projektinhabers vor dem Start; der Slice wählt nicht):

- **A, `--instance <name>` an `up`, `down` und `shell`:** `u-boot` übergibt der CLI zwei Id-Labels (`devcontainer.local_folder=<Ordner>` und ein eigenes Label mit dem Instanznamen, Arbeitsname `u-boot.instance`), und `FindContainer` filtert auf beide. Ohne Flag bleibt es beim heutigen Aufruf. Neue öffentliche Flags: vermutlich ein Lastenheft-CR ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt keine Instanz).
- **B, Git-Worktrees als Projektordner, kein Produktcode:** Jeder Worktree ist ein eigener Projektordner und damit ein eigener Container (eigenes `devcontainer.local_folder`); das Handbuch beschreibt den Weg. Der Container klont weiter aus `origin`, der Worktree liefert nur die Identität. Nicht gemessen.
- **C, nichts im Produkt.** Das Problem bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Welchen Zweig eine Instanz bearbeitet.** Der Clone folgt `devcontainer.sandbox.repository` oder `origin` und checkt den Standardzweig aus; die Wahl des Zweigs ist ein anderer Vorgang mit eigener Entscheidung (ein Feld je Instanz in `u-boot.yaml` oder ein Startparameter).
- **Agenten starten oder verteilen.** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nennt als Nicht-Ziel, dass `u-boot` keinen Agenten startet; ein Ablauf über mehrere Instanzen bleibt beim Aufrufer.
- **Der Compose-Stack je Instanz.** Der Befehl startet laut [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nicht den Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)).
- **Last und Grenzen bei vielen Instanzen.** Gemessen wird mit zwei Instanzen; ein Urteil über zehn wäre ein Schluss über die Messung hinaus ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)).
- **Andere Hosts als Linux `amd64`.** Sie sind für `up` und `shell` heute nicht zugesagt; der Slice ändert das nicht.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Zwei Instanzen desselben Projektordners laufen gleichzeitig (Linux `amd64`, lokale Engine, mindestens zwei Läufe, datiert, mit Umgebung): eigene Containerkennungen, eigene Volumes für Workspace, Docker-Speicher und Modelle, Image geteilt oder je Instanz gebaut, Dauer und Platzbedarf des zweiten Starts. *Beleg:* Messprotokoll unter `docs/messungen/<kennung>/` mit Rohdaten, Skript unter `tools/messung/<name>/` ([`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten)).
- [ ] **Der gewählte Weg** (A, B oder C): Bei A akzeptieren `up`, `down` und `shell` den Instanznamen, `FindContainer` findet den Container der Instanz, ohne Angabe ändert sich nichts (Bestandstest), ein ungültiger Name endet mit einem Exit-Code nach [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes); Tests je Bedingung und je Mutation, ein Docker-Test mit zwei Instanzen. Bei B oder C steht der Ablauf mit Datum und Messstand im Handbuch.
- [ ] **Spezifikation, Lastenheft (CR bei A), Handbuch und `CHANGELOG.md` nach dem Code** mit Beleg-Tabelle (Bedingung → Test oder Streichung); Handbuch und READMEs ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] `make gates` grün (im Modus der Umgebung).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die öffentliche Schnittstelle (CLI-Flag, Handbuch, README, `CHANGELOG.md`), soweit ein öffentlicher Vertrag berührt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Gelesen am 2026-10-31; die Zeilennummern verschieben sich, der Implementer liest sie am Ist-Stand nach ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat. Welche Zeilen der Weg berührt, entscheidet der Projektinhaber; die Tabelle nennt Weg A.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/messung/<name>/`, `docs/messungen/<kennung>/` | neu | Skript und Protokoll der Messung mit zwei Instanzen; Index in `docs/messungen/README.md` |
| `internal/adapter/driving/cli/devcontainer.go` | update | `up` (Flag `--recreate`, Zeile 164), `down` (ab Zeile 173) und `shell` (ab Zeile 209) erhalten `--instance <name>`; die Hilfetexte mit den Aufrufbeispielen (ab Zeile 84) folgen |
| `internal/hexagon/port/driven/devcontainer_launcher.go` | update | `DevcontainerUpOptions` (ab Zeile 113: `ProjectFolder`, `EnvNames`, `Recreate`, `ReadOnlyProject`, `Progress`) bekommt ein Feld für die Instanz |
| `internal/adapter/driven/docker/devcontainer_launcher.go` | update | `runArgs` (Zeilen 284 bis 303) hängt mit Instanz zwei `--id-label` an den Aufruf von `devcontainer up` |
| `internal/hexagon/port/driven/devcontainer_inspector.go` | update | `FindContainer(ctx, localFolder)` (Zeile 49) wird um die Instanz erweitert oder bekommt eine zweite Methode |
| `internal/adapter/driven/docker/inspector.go` | update | Filter `label=devcontainer.local_folder=<Ordner>` (Zeile 86) wird um das Instanz-Label ergänzt |
| `internal/hexagon/application/devcontainer_start.go`, `devcontainer_fallback.go` | update | Aufrufe von `FindContainer` (Zeilen 297, 429 und 155); `down` und `shell` suchen den Container der Instanz. Die Tests daneben (`devcontainer_start_test.go`, `devcontainer_shell_test.go`) folgen |
| `internal/e2e/devcontainer_start_docker_test.go` und Nachbarn | prüfen | Docker-Test mit zwei Instanzen; nicht gelesen |
| `spec/spezifikation.md` | update | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (ab Zeile 829, Absatz *Aufruf*): Flag, Fehlerklassen, Verhalten ohne Flag |
| `spec/lastenheft.md` | CR vor dem Code (Weg A) | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 507) nennt keine Instanz; die Zusage ändert sich, Wortlaut-Freigabe des Projektinhabers |
| `docs/user/benutzerhandbuch.md`, `README.md`, `README.de.md`, `CHANGELOG.md` | update (prüfen) | Abschnitt zu `up`, `down` und `shell` und die Release-Hinweise; nicht gelesen |
| `docs/plan/adr/` | neu | eine ADR für die Wahl des Weges und die Zusage-Menge (Vorbild: [ADR-0031](../../adr/0031-devcontainer-shell-ohne-editor.md)) |
| `harness/replay.md` | unverändert, prüfen | `up`, `down` und `shell` erzeugen keine Dateien; nicht gelesen |
| diese Datei, `docs/plan/planning/open/README.md` | update | Zeile dieses Slice; Moves per `make slice-mv` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; der Projektinhaber hat den Weg (A, B oder C) entschieden; bei A ist der Lastenheft-CR freigegeben; kein anderer Slice in `in-progress/` (WIP-Limit 1). Die Messung ist der erste Schritt und kann den Weg ändern.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Weg A berührt das Lastenheft — dann ist das ein Lastenheft-CR-Slice vor dem Code; oder die Messung zeigt, dass `devcontainerId` nicht mit dem Id-Label wechselt und die Volumes geteilt würden (dann ändert sich das Schema der Mounts, ein eigener Vorgang).
- `in-progress` → `open` (blockiert): Die Messung kann zwei parallele Instanzen nicht herstellen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün; die Zusage in Lastenheft, Spezifikation und Handbuch nennt nur das Gemessene (zwei Instanzen, Linux `amd64`); Review-Report und Verification-Report liegen vor; die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Mehrdeutige Suche:** Tragen zwei Container dasselbe Label `devcontainer.local_folder`, kann der Treffer der falsche sein ([`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md)); ein Instanz-Label löst das nur, wenn **alle** drei Aufrufer von `FindContainer` es mitgeben. — **Ausgang:** bei Closure.
- **Volumes und Image je Instanz** sind angenommen, nicht gemessen; wechselt die `devcontainerId` nicht mit dem Id-Label, teilen sich zwei Instanzen ein Workspace-Volume. — **Ausgang:** bei Closure (die Messung steht am Anfang).
- **Bestandscontainer ohne Instanz-Label:** Ein vor der Änderung angelegter Container trägt nur `devcontainer.local_folder`; ohne Flag muss er weiter gefunden werden. — **Ausgang:** bei Closure.
- **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): zwei Instanzen auf einem Host sind kein Beleg für viele oder für andere Hosts. — **Ausgang:** bei Closure.
- **Ressourcen:** Jede Instanz läuft mit `--privileged`, hat einen verschachtelten Docker mit eigenem Speicher und einen kalten Go-Cache im Workspace-Volume; Platz und Zeit sind nicht gemessen. — **Ausgang:** bei Closure.
- **Name der Instanz:** Zeichensatz und Länge (Label-Wert, Docker-Namen) sind nicht festgelegt. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Drei Paarungen.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/adapter/driving/cli` (CLI), `internal/hexagon/application` (APP), `internal/hexagon/port` (PORT), `internal/adapter/driven` (ADP), `internal/e2e` (E2E), `tools/` (WZG, Messskript) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK). Alle sind in der Modus-Deklaration von [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) als Greenfield auditiert; eine zu grobe Sub-Area ist nicht dabei.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen. Gelesen: [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (APP, ein Beleg, unter der Schwelle; sein Trigger tritt mit diesem Slice ein, ein zweiter Beleg wäre die Messung mit zwei Containern). Nur nach Titel gesichtet, Inhalt beim Start zu lesen: [`BEO-E2E/test-zaehlt-container-paralleler-pakete`](../observations/BEO-E2E/test-zaehlt-container-paralleler-pakete/observation.md) und [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Docker-Tests mit mehreren Containern), [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
