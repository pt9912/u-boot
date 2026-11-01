# Slice V1: Devcontainer mit sprechendem, eindeutigem Namen — der Container, den `u-boot devcontainer up` startet, trägt einen Namen mit Präfix statt des Zufallsnamens der Engine; Wege zur Namensform und zum Setzen zur Entscheidung des Projektinhabers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle, Release-Zuordnung offen (Frage 8 in §6). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden wäre. Anlass: Der Projektinhaber hat am 2026-11-01 eine Instanz gestartet (`u-boot devcontainer up --instance slice-v1-hilfetexte`); `docker ps` zeigte den Container als `mystifying_poincare` (Zufallsname der Engine), neben dem Label `u-boot.instance=slice-v1-hilfetexte`. Sein Urteil: „Dein container name ist kryptisch“, „der instance-name sollte eindeutig sein (mit prefix)“.

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`up`, `down`, `shell`, `--instance`), [`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) (erzeugte Devcontainer-Konfiguration; nur beim Weg über `runArgs`), [`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache) (Ausgaben Englisch), [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (Instanz über ein zweites Id-Label; bleibt unverändert, `Accepted`, [`AGENTS.md`](../../../../AGENTS.md) §3.5).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Aufruf von `up`, Hilfs-Container und Ablauf; nur bei Weg über den Aufruf der CLI), je nach Weg [`SPEC-005`](../../../../spec/spezifikation.md#spec-005--schema-der-projektkonfiguration-u-bootyaml) und die Beschreibung der erzeugten `runArgs` (nicht gelesen, prüfen).

**Verantwortlich:** pt9912 (Implementer)

**Autor:** pt9912 (Rolle Implementer, Entwurf im Auftrag des Projektinhabers). **Datum:** 2026-11-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Entwicklungscontainer, den `u-boot devcontainer up` startet, bekommt einen sprechenden, über Instanzen und Projekte eindeutigen Namen mit Präfix (Beispiel der Form: `u-boot-<projekt>-<instanz>`), statt des Zufallsnamens der Engine. Die Identität des Containers (Labels, `devcontainerId`, Volume-Namen) bleibt unverändert; ein Name darf sie nicht beeinflussen.

### Gefundener Stand (gelesen am 2026-11-01; Zeilen vor dem Diff)

- **Der Entwicklungscontainer bekommt von `u-boot` keinen Namen.** `runArgs` in `internal/adapter/driven/docker/devcontainer_launcher.go` (Zeilen 292 bis 316) baut die Argumente des **Hilfs-Containers**: `run --rm --name <name> --mount … --mount … [-e …] <image> devcontainer up --workspace-folder <Ordner> [--remove-existing-container] [--id-label … --id-label …]`. Das `--name` in Zeile 299 gehört dem Hilfs-Container (Präfix `helperContainerPrefix = "u-boot-devcontainer-up-"`, Zeile 29, plus Zufallsteil, Zeile 150); es ist **nicht** der Name des Entwicklungscontainers. Der Entwicklungscontainer wird von der Dev-Container-CLI angelegt; `u-boot` übergibt ihr keinen Namen. Dass die CLI dann einen Zufallsnamen der Engine stehen lässt, zeigt die Messung (`bold_lehmann`, `elated_faraday` in Abschnitt M9 des Messprotokolls [`slice-v1-devcontainer-instanz-messung`](../../../messungen/slice-v1-devcontainer-instanz-messung/README.md)) und der Anlass.
- **Ob die Dev-Container-CLI einen Namen annimmt, ist nicht gelesen und nicht gemessen.** Die Aufrufform `devcontainer up` im Aufruf oben kennt die Optionen `--workspace-folder`, `--remove-existing-container` und `--id-label` (gelesen in `runArgs`); ob es eine Option für den Containernamen gibt, ist die erste Messung (Messung 1).
- **`runArgs` der erzeugten `devcontainer.json`.** `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl` rendert `runArgs` aus `.RunArgs` (Zeile 15 im Sandbox-Zweig, Zeile 40 im Standardzweig); `.RunArgs` füllt `devcontainerTemplateData` (`internal/hexagon/application/devcontainer_sandbox.go`, Zeile 89 ff.) nur aus festen Quellen. Die Datei ist **eine je Projektordner**: Alle Instanzen desselben Ordners starten aus derselben `devcontainer.json` (der Hilfs-Container hängt den Ordner ein, `--workspace-folder <Ordner>`; die Instanz kommt nur über die zwei `--id-label`). Ein festes `--name=…` in `runArgs` gälte damit für **jede** Instanz und ließe die zweite Instanz an einem Namenskonflikt scheitern; ein Weg über `runArgs` allein trägt die Instanz nicht (siehe Weg B). Dass die Dev-Container-Spezifikation `${devcontainerId}` in `runArgs` auflöst, ist nicht gelesen und nicht gemessen (Messung 2).
- **Die Suche nutzt keinen Namen.** `Inspector.FindContainer` (`internal/adapter/driven/docker/inspector.go`, Zeile 85 ff.) sucht mit `docker ps -a --no-trunc --filter label=devcontainer.local_folder=<Ordner>` und, mit Instanz, `--filter label=u-boot.instance=<Name>`; die Ausgabe ist `{{.ID}} {{.State}} {{.Label "u-boot.instance"}}`. Weder `down` noch `shell` noch die Fallback-Auswertung brauchen den Containernamen ([ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md), Kontext: drei Aufrufstellen von `FindContainer`: `runningContainerID`, `Down`, `fallbackReader.read`). Ein Name ändert die Suche nicht.
- **Namensregel der Instanz.** `ValidateDevcontainerInstance` (`internal/hexagon/port/driven/devcontainer_launcher.go`, Zeile 47): 1 bis 32 Zeichen `a-z`, `0-9`, `-`, Rand nur Buchstabe oder Ziffer; `DevcontainerInstanceMaxLength = 32` (Zeile 37). Die Regel gilt dem Label-Wert.
- **Die Instanz-ADR und die Messung zu Identität und Volumes.** Die `devcontainerId` wechselt mit den Id-Labels (Messung M2: je Instanz ein eigener 52-stelliger Wert, alle paarweise verschieden), die Volume-Namen im Sandbox-Profil tragen `${devcontainerId}` ([ADR-0015](../../adr/0015-sandbox-volumes-pro-instanz.md)). Entscheidung 3 von [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md): Der Container ohne Instanz bekommt **kein** Label, weil ein Label die `devcontainerId` und damit die Volume-Namen bestehender Projekte ändert. Folgerung: Ein Weg, der Labels, Id-Labels oder Namen mit Einfluss auf die `devcontainerId` ändert, ist ausgeschlossen. Ein Name als solcher beeinflusst die `devcontainerId` nach dem Stand der Messung nicht (sie hängt an den Id-Labels, M2); gemessen mit einem gesetzten Namen ist es nicht (Messung 3).
- **Die Messung sagt zu Containernamen nur:** Die Container der CLI tragen Zufallsnamen (M9, Zeile mit `bold_lehmann`, `elated_faraday`); der Label-Wert nimmt viele Zeichen an (M8, 23 von 25 Namen, darunter Länge 1 bis 128, Großbuchstaben, `.`, `_`). Über `--name` oder Namensregeln der Engine sagt sie nichts.
- **Wortlaut zum Namen im Lastenheft und in der Spezifikation:** Suche nach `--name`, `container name`, `Containername`, `mystifying` in `spec/lastenheft.md`, `spec/spezifikation.md`, `docs/user/benutzerhandbuch.md`: einziger Treffer `spec/spezifikation.md` Zeile 853 (`docker run --rm --name <Name>`, der **Hilfs-Container**). Zum Namen des Entwicklungscontainers besteht keine Zusage.

**Deutung der Aussage „der instance-name sollte eindeutig sein (mit prefix)“** (Frage 1 in §6): (i) Gemeint ist der **Containername** (Docker `--name`), das ist der Gegenstand des Slice; (ii) gemeint wäre der Label-Wert `u-boot.instance` selbst mit einem Präfix. (ii) ändert die `devcontainerId` und damit die Volume-Namen jeder Instanz und ist durch [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (Label-Name und Mechanik `Accepted`, Entscheidung 2 und 3) ausgeschlossen; der Slice nimmt (i), solange der Projektinhaber nicht anders entscheidet.

**Wege zur Namensform** (der Slice wählt nicht):

- **A, `u-boot-<projekt>-<instanz>`, ohne Instanz `u-boot-<projekt>`.** *Nutzen:* sprechend (Projekt und Instanz sofort zu lesen), eindeutig, solange der Projektname im Docker-Daemon eindeutig ist. *Haken:* Zwei Klone mit gleichem Projektnamen (`project.name` oder Ordnername) und gleicher Instanz kollidieren (Docker-Fehler „name already in use“, nicht gemessen); die Länge ist Projektname plus Präfix plus bis zu 32 Zeichen Instanz. **Empfehlung des Entwurfs** zusammen mit einer Kollisionsregel (Weg D), weil der Name lesbar bleibt.
- **B, Weg A plus kurzer Hash des Projektordners** (`u-boot-<projekt>-<hash6>-<instanz>`). *Nutzen:* eindeutig über gleichnamige Klone. *Haken:* weniger sprechend; der Hash des absoluten Ordnerpfads ist host-spezifisch (kein Problem, der Container ist es auch), muss aber zwischen Läufen stabil sein.
- **C, nur `u-boot-<instanz>`.** *Nutzen:* kurz. *Haken:* zwei Projekte mit derselben Instanz (`slice-a`) kollidieren; für den Container ohne Instanz fehlt ein Name.
- **D, Kollisionsregel zu A:** bei vorhandenem Container gleichen Namens mit anderem `devcontainer.local_folder` Abbruch mit klarer Meldung (Exit-Code offen, Frage 5) statt eines Docker-Fehlers, oder Anhängen von `-2`. Ergänzung zu A, kein eigener Weg.

**Wege zum Setzen des Namens** (der Slice wählt nicht):

- **S1, `u-boot` hängt dem Aufruf der Dev-Container-CLI einen Namen an** (Option der CLI, falls sie existiert; Messung 1). *Nutzen:* gilt je Instanz, die `devcontainer.json` bleibt eine Datei; ändert die erzeugte Datei nicht (kein Golden Case). *Haken:* hängt davon ab, dass die gepinnte Version der CLI die Option kennt; sonst entfällt S1.
- **S2, `--name=…` in `runArgs` der erzeugten `devcontainer.json`.** *Nutzen:* wirkt auch, wenn der Nutzer die CLI außerhalb von `u-boot` benutzt. *Haken:* eine Datei je Ordner, also **ein** Name für alle Instanzen; ohne eine tragende Variable (`${devcontainerId}`, Messung 2) nur für den Container ohne Instanz brauchbar; ändert die erzeugte Datei (Golden Cases in [`harness/replay.md`](../../../../harness/replay.md), Bestand ohne Namen müsste byte-identisch bleiben oder wird bewusst geändert); eine Hand-Ergänzung überlebt `generate devcontainer` nicht (Messung des Entwurfs `slice-v1-devcontainer-ressourcenlimits`).
- **S3, beides:** S1 für Instanzen, S2 für den Container ohne Instanz. *Haken:* zwei Mechanismen für dieselbe Zusage.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Name des Hilfs-Containers** (`u-boot-devcontainer-up-<zufall>`, `devcontainer_launcher.go` Zeile 29 und 150): kurzlebig, bleibt (Spezifikation Zeile 853).
- **Label-Name, Label-Wert und `devcontainerId`** ([ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidungen 2 und 3): unverändert; bestehende Instanzen verwaisen sonst.
- **Umbenennen bestehender Container:** Sie behalten ihren Zufallsnamen bis zum Neuanlegen (`up --recreate`); `u-boot` führt kein `docker rename` aus (Frage 5).
- **Namen der Volumes, Netzwerke und des Compose-Stacks:** Volume-Namen folgen aus `devcontainerId`; der Stack gehört nicht zur Instanz ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Nicht-Ziele).
- **Suche nach Namen:** `FindContainer` bleibt auf Labels; ein Name ist Anzeige, kein Schlüssel.

## 2. Definition of Done

Bis zu drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten zählen nicht mit.

- [ ] **Messung und Freigabe** ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)): Die fünf Messungen aus §6 sind gefahren (Linux `amd64`, lokale Engine, mindestens zwei Läufe je Zelle, datiert, mit Umgebung; Protokoll unter `docs/messungen/slice-v1-devcontainer-container-name/`, Skripte unter `tools/messung/`, keine Zugangsdaten); der Projektinhaber hat Namensform, Weg zum Setzen und den Wortlaut des Lastenheft-CR freigegeben (Alt/Neu-Vorschlag in §6).
- [ ] **Der Name:** `u-boot devcontainer up` (mit und ohne `--instance`, `--recreate`) startet einen Container mit dem Namen der gewählten Form; `down` und `shell` finden ihn weiter über die Labels (Test); `devcontainerId` und Volume-Namen sind gegenüber dem Lauf ohne Namen **gleich** (Test oder Messung, die [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) Entscheidung 3 pinnt); ein ungültiger oder zu langer Bestandteil endet nach der gewählten Regel (Exit-Code `10` bei einem Namen, den die Engine ablehnen würde); Kollision nach Frage 5. Tests je Bedingung und je Mutation; Docker-Test hinter dem Tag `docker` für das Ankommen des Namens im gestarteten Container (Nachweis: CI-Lauf, nicht lokal, [`slice-v1-test-docker-in-ci`](../open/slice-v1-test-docker-in-ci.md)).
- [ ] **Doku und Zusage-Menge:** Spezifikation, Lastenheft (CR), Handbuch und `CHANGELOG.md` nennen nur das Gemessene (Linux-Host `amd64`, lokale Engine); jede Bedingung in Spec oder Handbuch hat einen pinnenden Test oder entfällt ([`AGENTS.md`](../../../../AGENTS.md) §3.10, letzter Punkt); Handbuch ohne Kennungen ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6).
- [ ] `make gates` grün (im Modus der Umgebung); Review (`.harness/skills/reviewer.md`), Verification, Closure-Notiz mit Lerneintrag, Beobachtungs-Register fortgeschrieben, jedes Risiko aus §6 mit Ausgang, drei Paarungen getragen.

## 3. Plan (vor Code)

Gelesen am 2026-11-01 ([`AGENTS.md`](../../../../AGENTS.md) §3.10); der Implementer gleicht die Tabelle nach der letzten Fix-Runde gegen `git diff --stat` ab. „prüfen“ nennt eine Datei, die der Planner nicht an der berührten Stelle gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung / Fundstelle |
|---|---|---|
| `tools/messung/devcontainer-container-name/messen.sh`, `docs/messungen/slice-v1-devcontainer-container-name/`, `docs/messungen/README.md` | neu / update | Messungen 1 bis 5; Ablage nach [`harness/conventions.md`](../../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) |
| `internal/adapter/driven/docker/devcontainer_launcher.go`, `devcontainer_launcher_test.go` | update (S1') | `runArgs` (Zeile 292): wenn ein Name gesetzt ist, rendert der Adapter die je Aufruf gültige Konfigurationskopie (Inhalt der erzeugten `devcontainer.json` plus `--name=<Name>` in `runArgs`), reicht sie per zusätzlicher Bind-Mount in den Hilfs-Container und hängt `--override-config <Pfad im Hilfs-Container>` an; ohne Name unverändert; Unit-Test der Argumentliste mit und ohne Namen |
| `internal/hexagon/port/driven/devcontainer_launcher.go` | update | `DevcontainerUpOptions` (Zeile 139) bekommt den Namen (oder seine Bestandteile) als Feld |
| `internal/hexagon/application/devcontainer_start.go` | update | bildet den Namen aus Projektname und Instanz und reicht ihn an den Launcher (Projektname-Quelle zu lesen: `u-boot.yaml` `project.name` oder Ordnername, nicht gelesen) |
| `internal/hexagon/domain/` (neue Datei, Name beim Bau) | neu | Wertobjekt `ContainerName` mit der Bildungsregel und der Prüfung nach den Engine-Regeln (Messung 4) |
| `internal/hexagon/application/devcontainer_sandbox.go`, `templates/devcontainer/devcontainer.json.tmpl`, `generate_*_test.go`, `harness/replay.md` | **unverändert** (S1') | erzeugte Datei bleibt byte-identisch; der Name wird nicht in die Vorlage gerendert (Messung 6 belegt den Weg am Overlay, Golden Cases laufen unverändert) |
| `internal/adapter/driven/docker/inspector.go` | **unverändert** | gelesen: `FindContainer` (Zeile 85) sucht über Labels; ein Name berührt die Suche nicht |
| `internal/hexagon/port/driven/devcontainer_launcher.go` (`ValidateDevcontainerInstance`) | **unverändert** | gelesen: Regel des Label-Werts; der Containername ist eine andere Regel |
| `internal/e2e/` (Datei mit Build-Tag `docker`, Name beim Bau) | neu | Name des gestarteten Containers, `devcontainerId` und Volumes unverändert; Vorbild `devcontainer_start_docker_test.go` (nicht gelesen) |
| `internal/adapter/driving/cli/devcontainer.go` | update | Ausgabe von `up`/`shell` („Devcontainer started: container <zwölf Zeichen> …“, Spezifikation Zeile 887): Name neben der ID ausgeben (Frage 6 entschieden, 2026-11-01) |
| `internal/hexagon/application/doctor*.go` | prüfen | ob `doctor` Container sucht (nach Name oder Label); nicht gelesen |
| `harness/replay.md` | update (S2/S3) | Golden Case zu `generate devcontainer` mit `runArgs` |
| `spec/lastenheft.md` | **CR vor dem Code** (Wortlaut-Freigabe) | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten); Alt/Neu in §6 Frage 7; nicht in diesem Entwurf geschrieben |
| `spec/spezifikation.md` | update (nach dem Code) | [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), Block „Hilfs-Container“ ab Zeile 850 (gelesen: Aufrufform ohne Containernamen) |
| `docs/user/benutzerhandbuch.md`, `CHANGELOG.md` | update | Abschnitt zu mehreren Instanzen; ohne Kennungen |
| `docs/plan/adr/` und `docs/plan/adr/README.md` | neu, nur wenn Frage 9 „ja“ | Folge-ADR zu [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) (`Supersedes` entfällt, die ADR bleibt gültig; `Schärft:`) mit Index-Eintrag |
| diese Datei, `docs/plan/planning/open/README.md` | neu / update | Zeile im Bestand; Moves per `make slice-mv` |

**Wächter und rot färbende Mutationen** (Implementer führt sie aus):

| Wächter | Mutation, die ihn rot färbt |
|---|---|
| Unit-Test der Argumentliste: Name steckt in der je Aufruf gültigen Konfiguration und `--override-config` steht im Aufruf (S1') | das Overlay-Datei-Rendern oder die Option beim Zusammenbau von `runArgs` weglassen |
| Test „Bildungsregel des Namens“ (mit/ohne Instanz, Länge, Zeichen) | Präfix oder Trenner ändern; Instanz nicht anhängen |
| Test „Name verändert `devcontainerId` und Volume-Namen nicht“ (Docker-Test oder Messskript) | den Namen als `--id-label` an die CLI geben |
| Test „`down`/`shell` finden den Container weiter über Labels“ | `FindContainer` auf den Namen umstellen |
| Docker-Test „Name im gestarteten Container“ | Namen im Launcher auslassen; `docker inspect` zeigt den Zufallsnamen |
| Kollisionstest (Frage 5) | Prüfung entfernen, der Docker-Fehler „name already in use“ erscheint unübersetzt |

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): Start in einer eigenen Devcontainer-Instanz, Projektinhaber 2026-11-01; kein Termin. Dazu: Namensform und Weg zum Setzen entschieden (Fragen 2 und 3), `Verantwortlich:` gesetzt, Wortlaut des CR freigegeben (Frage 7), kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `open` (blockiert): Der Projektinhaber entscheidet eine neue Frage nicht (etwa ein Wortlaut-Nachtrag des CR nach dem Code). Die ursprüngliche Rückführung „zu groß“ entfällt: Messung 1 und die Zusatzprobe 6 zeigen, dass S1' trägt und die Zusage je Instanz ohne neues Konstrukt gesetzt werden kann.

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Docker-Test im CI-Lauf grün; Review- und Verification-Report liegen vor; Closure-Notiz mit Steering-Loop-Lerneintrag; `git mv` nach `done/` als eigener Commit nach dem Inhalt ([`AGENTS.md`](../../../../AGENTS.md) §3.3 Fall 2).

## 6. Risiken und offene Punkte

**Fragen an den Projektinhaber** (Entwurf; entschieden am 2026-11-01: Fragen 1 bis 5, 7 bis 9 — 1 durch Nicht-Widerspruch zur Deutung des Entwurfs, 3 und 7 nach den Messungen; Frage 6 auf „Name mit ausgeben“):

1. **Deutung:** Containername (Docker `--name`), wie der Entwurf annimmt — oder der Label-Wert der Instanz mit Präfix ([ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) schließt das aus, Volumes wechseln).
2. **Namensform:** A, B oder C aus §1, jeweils auch für den Container ohne Instanz (`u-boot-<projekt>`). **Entschieden (Projektinhaber, 2026-11-01): A**, `u-boot-<projekt>-<instanz>` beziehungsweise `u-boot-<projekt>` ohne Instanz.
3. **Weg zum Setzen:** **S1'** (Projektinhaber, 2026-11-01, nach der Zusatzprobe 6): der Hilfs-Aufruf übergibt `--override-config` mit einer je Aufruf gerenderten Kopie der erzeugten Konfiguration samt Namen in `runArgs`; die erzeugte `devcontainer.json` bleibt unberührt. S1 entfällt (CLI kennt keine Namens-Option), S2 verworfen (erzeugte Datei, Golden Cases, ein Name für alle Instanzen).
4. **Eindeutigkeit über Projekte:** **A mit Abbruch bei Kollision** (Projektinhaber, 2026-11-01) — gleichnamige Klone kollidieren mit klarer Meldung, statt eines Hash-Teils (B).
5. **Kollision und Bestand:** **Abbruch mit klarer Meldung** (Projektinhaber, 2026-11-01) bei einem Container gleichen Namens mit anderem `devcontainer.local_folder`; die Messungen 4 und 5 belegen: die Engine weist laufende und gestoppte Container mit exit=125 zurück, die CLI endet exit=1 und `--remove-existing-container` berührt den fremden Container nicht; bestehende Container behalten ihren Namen bis zum Neuanlegen.
6. **Ausgabe:** **Name mit ausgeben** (Projektinhaber, 2026-11-01) — `up` und `shell` nennen den Containernamen neben der Container-ID (Spezifikation Zeile 887 wird erweitert).
7. **Lastenheft-CR: Wortlaut und Gefäß** — Ort: Absatz zu `up` in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), im selben Slice. *Neu (Entwurf):* „Der Container, den `u-boot devcontainer up` startet, trägt den Namen `u-boot-<Projekt>` beziehungsweise, mit `--instance <name>`, `u-boot-<Projekt>-<name>`. Die Suche nach dem Container, die Labels und die Namen der Volumes bleiben unverändert. Zugesagt ist der Name nur auf einem Linux-Host `amd64` mit lokaler Engine und nur, soweit ein Test mit Build-Tag `docker` ihn belegt; ein Container, der vor dieser Fassung angelegt wurde, behält seinen Namen bis zum Neuanlegen.“ **Freigegeben (Projektinhaber, 2026-11-01)** wie in §6 Frage 7 vorgelegt; „vor dieser Fassung“ ist Chronik und wird beim Schreiben des CR in den Ist-Zustand übertragen.
8. **Release-Zuordnung:** **Minor, *Changed*** (Projektinhaber, 2026-11-01) — nächstes Minor nach `v0.17.0`; das Verhalten von `up` ändert sich für bestehende Aufrufe.
9. **ADR-Bedarf:** **kein Folge-ADR** (Projektinhaber hat der Empfehlung nicht widersprochen, 2026-11-01): [ADR-0038](../../adr/0038-devcontainer-instanz-zweites-id-label.md) bleibt unberührt; der Name ist eine Produktzusage im CR, der Weg S1' ist Umsetzungsdetail des Hilfs-Aufrufs.

**Messungen vor dem Code** (gemessen am 2026-11-01, Protokoll [`docs/messungen/slice-v1-devcontainer-container-name`](../../../messungen/slice-v1-devcontainer-container-name/README.md), je Messung zwei Läufe vollständig, Skript `tools/messung/devcontainer-container-name/messen.sh`):

1. **Nimmt die CLI einen Namen an?** Die gepinnte CLI 0.89.0 kennt **keine** Option für den Containernamen an `devcontainer up`; die Eigenschaft `containerName` der `devcontainer.json` wird **nicht angewendet**; **`runArgs` mit `--name=<Name>` wird angewendet**.
2. **Variablen in `runArgs`:** `${devcontainerId}` wird **aufgelöst** (Name am Container enthielt die 52-stellige ID) — eindeutig je Instanz, aber nicht sprechend.
3. **Identität mit gesetztem Namen:** In allen vier Zellen (beide Profile, mit und ohne Instanz) bleiben `devcontainerId`, Volume-Namen und Label-Wert **gleich** zum Lauf ohne Namen (Risiko 2 gemessen widerlegt).
4. **Namensregeln der Engine (29.9.2):** erlaubt sind die Zeichen nach `[a-zA-Z0-9][a-zA-Z0-9_.-]` mit Erstzeichen Buchstabe/Ziffer; Länge wurde bis 257 Zeichen ohne Zurückweisung gemessen (Grenze oberhalb gesucht nicht); exit=125 mit Wortlaut `invalid container name (<Name>): only [a-zA-Z0-9][a-zA-Z0-9_.-] are allowed` für führende `-`/`.`/`_`, Leer- und Sonderzeichen; **Kollision**: laufender *und* gestoppter Container gleichen Namens → exit=125 `Conflict. The container name … is already in use`, der haltende Container bleibt.
5. **Kollision und Bestand über die CLI:** gleicher Ordner + gleicher Name → Übernahme desselben Containers (exit=0); **fremder Container gleichen Namens** (anderes `devcontainer.local_folder`, laufend oder gestoppt) → CLI endet exit=1, `--remove-existing-container` endet ebenfalls exit=1 und berührt den fremden Container nicht; der Docker-Wortlaut steht hinter der 260-Zeichen-Kürzung der Rohdaten (Exit-Weg gemessen, Wortlaut von Messung 4).
6. **Zusatzprobe (für Frage 3):** `--override-config` mit einer je Aufruf gerenderten Kopie der erzeugten `devcontainer.json`, Name in `runArgs` — Name wird gesetzt (die sprechende Form ist damit **je Instanz** über den Hilfs-Aufruf erreichbar), `devcontainerId`/Volumes/Labels gleich, erzeugte Datei unberührt.

**Risiken:**

1. **Die Dev-Container-CLI kennt keine Namensoption (Messung 1) und `runArgs` trägt keine Instanz:** Dann bleibt nur `docker rename` nach dem Start (ein zweiter Docker-Aufruf, Fensterzeit mit Zufallsnamen) oder der Wechsel der Entscheidung. — **Ausgang:** bei Closure.
2. **Der Name ändert die Identität:** wäre `devcontainerId` namensabhängig, wechselten Volume-Namen bestehender Projekte (Datenverlust-Eindruck). — **Ausgang:** Messung 3 und Test, der es pinnt; bei Closure.
3. **Kollision zweier Projekte oder Klone:** Docker-Fehler statt einer klaren Meldung. — **Ausgang:** bei Closure (Frage 5).
4. **Zusage weiter als die Messung** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)): Linux-Host `amd64`, lokale Engine; andere Hosts nicht zugesagt. — **Ausgang:** bei Closure.
5. **Bedingung in Spec oder Handbuch ohne Test** ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)): „Name“, „Labels unverändert“ brauchen je einen pinnenden Test. — **Ausgang:** bei Closure.
6. **Docker-Test läuft nicht lokal:** Nachweis ist der CI-Lauf; der Docker-Test hinterlässt keine Container ([`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md)). — **Ausgang:** bei Closure.

## 7. Closure-Notiz

Noch nicht geschlossen; der Slice liegt in `open/` (Entwurf).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden je nach Weg `internal/adapter/driven/docker` (ADP), `internal/hexagon/application` (APP), `internal/hexagon/domain` (DOM), `internal/hexagon/port/driven` (PORT), `internal/adapter/driving/cli` (CLI, nur Ausgabe), `internal/e2e` (E2E), `tools/` (WZG) sowie `spec/`, `harness/`, `docs/plan/` und `docs/messungen/` (DOK); alle Greenfield laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area), die Schwelle ≥ 2 von 3 Achsen ist dort entschieden.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen, nicht verändert): [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (Risiken 4 und 5), [`BEO-DOK/messvariante-belegt-nicht-die-auslieferung`](../observations/BEO-DOK/messvariante-belegt-nicht-die-auslieferung/observation.md) (Messung an der Umsetzung, nicht an einer Hand-Konstruktion), [`BEO-E2E/docker-test-hinterlaesst-zwischencontainer`](../observations/BEO-E2E/docker-test-hinterlaesst-zwischencontainer/observation.md) (Risiko 6), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (Prüflinse dieser Tabelle; „nicht gelesen“ ist an den Zeilen vermerkt). Zähler nicht neu gezählt.

Evidenz-/Diskrepanz-Risiko: mittel (Messung 1 entscheidet, ob S1 trägt). Alle berührten Sub-Areas Greenfield (kein Modus-Begründungsblock nötig).
