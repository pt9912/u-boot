# ADR 0038: Mehrere Devcontainer je Projekt — `--instance <name>` an `up`, `down` und `shell` über ein zweites Id-Label; der Container ohne Instanz bleibt unverändert, ohne Angabe wirkt der Befehl nur auf ihn

**Status:** Proposed

**Datum:** 2026-10-31

**Autor:** pt9912 (Rolle Architect)

**Bezug:** [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [ADR-0015](0015-sandbox-volumes-pro-instanz.md), [ADR-0028](0028-devcontainer-start-ohne-editor-hilfs-image.md), [ADR-0031](0031-devcontainer-shell-ohne-editor.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Absatz *Aufruf*, Voraussetzungen von `up` und `shell`, Hilfs-Container: Flag `--instance`, Namensregel, zweites Id-Label, Verhalten ohne Flag, Exit-Code `10`). Von [ADR-0015](0015-sandbox-volumes-pro-instanz.md) werden zwei Sätze geschärft: der tragende Satz der Entscheidung „Eine Instanz ist ein Ordner … Parallele Instanzen desselben Projekts laufen aus getrennten Ordnern (Klone, Git-Worktrees)“ und der Satz „zwei Instanzen aus **demselben** Ordner bleiben eine Instanz“ (Variante D, Spalte Contra). Beide gelten weiter für den Container ohne Instanz; mit `--instance` laufen mehrere Instanzen aus demselben Ordner, jede mit eigener `devcontainerId`. [ADR-0015](0015-sandbox-volumes-pro-instanz.md), [ADR-0028](0028-devcontainer-start-ohne-editor-hilfs-image.md) und [ADR-0031](0031-devcontainer-shell-ohne-editor.md) bleiben unverändert; diese ADR ergänzt sie und löst keine ihrer Entscheidungen ab.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Projektinhaber will mehrere Agenten getrennt arbeiten lassen und statt Git-Worktrees mehrere Devcontainer desselben Projekts starten (Gespräch am 2026-10-31). Heute ist ein Container an den Projektordner gebunden.

**Stand im Code (gelesen am 2026-10-31, mit Fundstelle).**

- `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`) ruft die Dev-Container-CLI als `devcontainer up --workspace-folder <Ordner>` und gibt kein `--id-label` mit. Die CLI leitet ihre Id-Labels dann aus dem Ordner ab.
- `FindContainer` (`internal/adapter/driven/docker/inspector.go`) sucht mit `docker ps -a --filter label=devcontainer.local_folder=<Ordner>` und liefert den ersten Treffer, der zuletzt angelegte (Port `DevcontainerInspector`, `internal/hexagon/port/driven/devcontainer_inspector.go`).
- `FindContainer` wird an den drei Stellen `runningContainerID` und `Down` (`internal/hexagon/application/devcontainer_start.go`) und `fallbackReader.read` (`internal/hexagon/application/devcontainer_fallback.go`) gerufen. Der `fallbackReader` dient `up` und `shell` (Urteil nach dem Start, `startOutcome`) sowie `doctor`, `generate devcontainer` und `init`.
- Die Namen der Volumes im Sandbox-Profil tragen `${devcontainerId}` ([ADR-0015](0015-sandbox-volumes-pro-instanz.md)); die CLI bildet die `devcontainerId` aus ihren Id-Labels.

**Messung.** Das Messprotokoll [`docs/messungen/slice-v1-devcontainer-instanz-messung/README.md`](../../messungen/slice-v1-devcontainer-instanz-messung/README.md) trägt die Werte; diese ADR wiederholt sie nicht. Aus ihm folgen, mit Stufen und Läufen dort, die fünf Punkte, auf denen die Entscheidung steht: (1) `devcontainer up --id-label devcontainer.local_folder=<Ordner> --id-label u-boot.instance=<Name>` legt je Name einen eigenen Container mit eigenen Volumes an (Abschnitte *M1*, *Urteil zur Kernfrage*); (2) der Container ohne Flag trägt `devcontainer.config_file`, die Instanz-Container nicht (*M1*); (3) `up` ohne Flag neben laufenden Instanzen übernimmt den zuletzt angelegten Instanz-Container, statt einen eigenen anzulegen (*M6*; der Mechanismus ist dort Vermutung); (4) `docker ps --filter label=devcontainer.local_folder=<Ordner>` trifft alle Container des Ordners, die Instanz zuerst, wenn sie die neuere ist (*M1*, *M6 umgekehrt*); (5) ein Namenswert mit `=` lässt `devcontainer up` nicht zurückkehren, die leere Zeichenkette lehnt die CLI ab (*M8*).

**Gemessen im Standardprofil** (Rohdaten `rohdaten-stufe{2,3,4}-lauf{1,2}.txt`, Abschnitt `M6 umgekehrt`, sechs Läufe, Bestandscontainer ohne Flag, danach Instanz `i1`): `up` ohne Flag trifft danach den Container ohne Instanz, nicht die neuere Instanz (Zeile „ohne Flag trifft“). Das Protokoll-README sagt im ersten M6-Absatz, das Standardprofil neben Instanzen sei nicht gemessen, und im Abschnitt *M6 umgekehrt* das Gegenteil; die Rohdaten gelten, weil sie die Zeile tragen. **Nicht gemessen** und deshalb hier nicht behauptet: wie `docker ps` ein fehlendes Label ausdrückt; `up` ohne Flag bei Container ohne Instanz und neuerer Instanz im Sandbox-Profil; `down` und `shell` in dieser Lage in beiden Profilen; wiederholtes `up --instance <name>` ohne `--recreate` (gemessen ist nur `--remove-existing-container`); `down --instance <name>` mit einem Namen, zu dem es keinen Container gibt; dass die Instanz-Container kein `devcontainer.config_file` tragen, ist gemessen (*M1*), dass Werkzeuge, die darüber suchen, sie deshalb nicht finden, nicht.

Gilt Punkt 1 nicht mehr (eine neue Version der Dev-Container-CLI wechselt die `devcontainerId` nicht mehr mit einem zusätzlichen Id-Label), entfällt die Entscheidung.

## Entscheidung

### Entscheidungen des Projektinhabers (2026-10-31)

Im Gespräch am 2026-10-31 entschieden; ein anderer Beleg als diese Zeilen liegt nicht vor.

- **Weg A.** Drei Wege standen zur Wahl: A, ein Flag `--instance <name>`, das `u-boot` in Id-Labels übersetzt; B, Git-Worktrees als Projektordner ohne Produktcode; C, nichts im Produkt. Gewählt ist A.
- **Begriff Instanz.** Eine Instanz ist ein vom Nutzer mit `--instance <name>` benannter weiterer Devcontainer desselben Projektordners; „der Container ohne Instanz“ ist der Container, den `up` ohne Flag anlegt.
- **Namensregel.** Kleinbuchstaben, Ziffern und Bindestriche, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32; die leere Zeichenkette und ein Name mit `=` sind ungültig; ein ungültiger Name endet mit Exit-Code `10`.
- **Fassung (i) für das Verhalten ohne Flag.** Zur Wahl standen (i), `up`, `down` und `shell` ohne `--instance` wirken nur auf den Container ohne Instanz, und eine Fassung, in der sie einen Container des Projekts übernehmen. Gewählt ist (i); gibt es nur Container von Instanzen, endet jeder der drei mit Exit-Code `10` und einem Hinweis auf `--instance`, geprüft von `u-boot` vor dem Aufruf der CLI.
- **Kein Label für den Container ohne Instanz** (Entscheidung 3).
- **Der Labelname `u-boot.instance` — wartet auf Bestätigung bei `Accepted`.** Der Name ist der Arbeitsname aus dem Messprotokoll, keine Entscheidung des Projektinhabers; mit `Accepted` wird er bestätigt oder ersetzt, ein späterer Wechsel ließe bestehende Instanzen verwaisen.

Die Mechanik der Entscheidungen unten schlägt der Architect vor; der Projektinhaber bestätigt sie mit `Accepted`.

1. **Flag.** `u-boot devcontainer up`, `down` und `shell` nehmen `--instance <name>`. Mit dem Flag wirkt der Befehl nur auf den Container dieser Instanz; mehrere Instanzen desselben Projektordners laufen gleichzeitig.
2. **Zweites Id-Label.** Mit `--instance` hängt `u-boot` dem Aufruf der CLI zwei Id-Labels an: `devcontainer.local_folder=<Ordner>` und `u-boot.instance=<name>`. Gemessen ist, dass die `devcontainerId` mit dem zweiten Label wechselt und damit die Volume-Namen aus [ADR-0015](0015-sandbox-volumes-pro-instanz.md); jede Instanz hat eigene Volumes. `FindContainer` findet die Instanz über `docker ps -a --filter label=devcontainer.local_folder=<Ordner> --filter label=u-boot.instance=<name>`. Der Labelname `u-boot.instance` ist der Arbeitsname aus dem Messprotokoll und wartet auf die Bestätigung bei `Accepted` (siehe oben).
3. **Der Container ohne Instanz bekommt kein Label.** Ohne `--instance` ruft `u-boot` die CLI wie heute, ohne `--id-label`. Ein Label (etwa `u-boot.instance=default`) änderte die `devcontainerId` und damit die Volume-Namen bestehender Projekte; ein Container ohne Instanz soll nach einem Upgrade von `u-boot` derselbe bleiben.
4. **Ohne Flag nur der Container ohne Instanz.** Die Suche ohne `--instance` liefert einen Container mit `devcontainer.local_folder=<Ordner>`, der **kein** `u-boot.instance` trägt. Ein Filter auf ein fehlendes Label ist im Repo nicht gelesen und nicht gemessen; `u-boot` liest darum das Label je Treffer aus der Ausgabe von `docker ps` und wählt selbst, der neueste Treffer ohne Instanz-Label zuerst. Dass `docker ps --format` das Label je Treffer ausgeben kann (das Muster `{{.Label "…"}}`), ist nicht im Repo gelesen und nicht gemessen; der Docker-Test des Slice belegt es, sonst ist die Wahl der Ausgabe `{{.Labels}}` zu prüfen. Das gilt für `up`, `down`, `shell` und für die Aufrufer, die heute ohne Instanz suchen (`doctor`, `generate devcontainer`, `init`): Sie bleiben beim Container ohne Instanz.
5. **Prüfung vor dem Aufruf der CLI.** Gibt es im Projekt Container, aber keinen ohne Instanz, endet `up`, `down` und `shell` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, bevor ein Hilfs-Container oder ein `docker stop` gestartet wird. `u-boot` prüft das selbst, weil die CLI sonst nach dem Messprotokoll (*M6*) den zuletzt angelegten Instanz-Container übernähme. Gibt es im Projekt gar keinen Container, startet `up` und `shell` wie heute den Container ohne Instanz, und `down` meldet wie heute, dass nichts zu stoppen ist (Exit-Code `0`).
6. **Namensregel.** Ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist 1 bis 32 Zeichen lang; die leere Zeichenkette und ein Name mit `=` sind ungültig (*M8*). Ein ungültiger Name endet mit Exit-Code `10`, bevor `u-boot` Docker anspricht. Die Regel steht einmal als Prüffunktion im Port-Paket `driven` (Muster `ValidateDevcontainerProjectFolder`); die Application ruft sie vor dem Start, der Adapter ruft sie vor dem Bau des Aufrufs und der Suche, weil der Wert mit `=` die CLI hängen lässt.
7. **Schnitt durch die Schichten.** Die Instanz reist als Wert (`""` heißt: ohne Instanz) von der Anfrage der CLI (`DevcontainerUpRequest`, `DevcontainerDownRequest`, `DevcontainerShellRequest`) durch die Application in `DevcontainerUpOptions` des Launchers (`runArgs` hängt die zwei `--id-label` an) und in die Suche des Ports `DevcontainerInspector`. Der Port unterscheidet „es gibt keinen Container“ von „es gibt nur Container von Instanzen“; die Form (Sentinel im Port oder Feld im Ergebnis) legt der Slice fest. Die Application bildet beides auf zwei Sentinels des Pakets `driving` ab, den ungültigen Namen und den Fall „nur Instanzen“; beide sind Klasse `10` ([`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)).
8. **Exit-Code-Abbildung nach der Dual-Classifier-Regel** (Sicht-Spec, Abschnitt Fehlermodelle): Jeder der zwei Sentinels steht in der zentralen Klassifikation (`isDevcontainerValidationError`, von `isValidationError` eingebunden) und in der Diagnostic-Abbildung von `up` und `down`. `down` führt bisher keinen Sentinel der Klasse `10` (`mapDevcontainerDownErrorToDiagnostic`) und bekommt hier den ersten. Die Abbildung von `shell` bleibt auf Klasse `2` beschränkt, weil `--json` dort vor dem Use-Case abgewiesen wird; der Exit-Code von `shell` läuft über die zentrale Klassifikation. Je Sentinel und Unterbefehl ein Test, der die Klasse pinnt.
9. **Urteil nach dem Start.** `up` und `shell` lesen die Statusdateien und fragen, ob der Container läuft, bei dem Container der gewählten Instanz; der Wert der Instanz geht in `runningContainerID` und `startOutcome` ein.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts im Produkt (Weg C) | keine Vertragsfläche, kein Code | das Problem bleibt: ein Container je Projektordner |
| B — Git-Worktrees als Projektordner (Weg B), nur Handbuch | kein Produktcode; jeder Worktree ist ein eigener Ordner und damit ein eigener Container (eigenes `devcontainer.local_folder`) | nicht gemessen; ein Worktree je Instanz zwingt zu einem Ordner je Agent auf dem Host, obwohl der Container aus `origin` klont und den Ordner nur für die Identität braucht; [ADR-0015](0015-sandbox-volumes-pro-instanz.md) nennt diesen Weg bereits (Klone, Git-Worktrees); der Projektinhaber hat sich am 2026-10-31 für Weg A entschieden |
| C — Instanzname im Ordnerpfad (Kopie oder Verknüpfung des Projektordners je Instanz) | nutzt den Mechanismus von [ADR-0015](0015-sandbox-volumes-pro-instanz.md) ohne neues Flag | verändert den Pfad, den der Hilfs-Container einhängt (`--mount`, `--workspace-folder`), und legt Ordner an, die `u-boot` verwalten müsste; `u-boot` führt nach [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) keine Liste und räumt nichts auf |
| D — Compose-Projektname je Instanz | Compose kennt die Trennung | der Compose-Stack gehört laut [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) nicht zum Befehl und nicht zu einer Instanz; ob ein Compose-Projektname die Id-Labels oder die Volume-Namen des Devcontainers ändert, ist nicht gelesen und nicht gemessen |
| E — ein Label auch für den Container ohne Instanz (etwa `u-boot.instance=default`) | eindeutige Suche: jeder Container trägt das Label, `docker ps --filter` genügt | ändert die `devcontainerId` und damit die Volume-Namen jedes bestehenden Projekts (die `devcontainerId` wechselt mit den Id-Labels, gemessen); die Suche der CLI nach einem bestehenden Container läuft über diese Labels, daher fände schon das nächste `up` ohne `--recreate` den Bestandscontainer nicht mehr und legte einen neuen mit neuen Volumes an; Container und Volumes des alten blieben verwaist (die Suche über die Id-Labels ist die Beschreibung der CLI in `devcontainer up --help`, ihr Ergebnis hier nicht gemessen) |
| F — ohne `--instance` den neuesten Container des Ordners nehmen | keine Prüfung, kein Exit-Code `10` | `down` stoppt eine fremde Instanz, `shell` öffnet sie, `up` startet sie; das ist das gemessene Verhalten der CLI (*M6*) und genau die Verwechslung, die das Flag vermeiden soll |
| **G — zweites Id-Label für Instanzen, kein Label für den Container ohne Instanz, Prüfung durch `u-boot` vor dem Aufruf (gewählt)** | gemessen: eigene Container und Volumes je Instanz; bestehende Projekte ändern sich nicht; ohne Flag ist das Verhalten das heutige, mit einer klaren Ablehnung statt einer Verwechslung | `u-boot` liest das Label selbst aus der Ausgabe von `docker ps`, weil ein Filter auf ein fehlendes Label nicht gemessen ist; Container von Instanzen bleiben für `doctor`, `generate devcontainer` und `init` unsichtbar |

## Konsequenzen

- Positiv: Mehrere Devcontainer desselben Projekts laufen mit getrennten Volumes (die Zusage im Lastenheft nennt die gemessene Menge); `down` und `shell` können eine Instanz gezielt ansprechen; bestehende Projekte behalten Container, Volumes und Aufruf.
- Negativ: Ein Projekt mit nur Instanz-Containern verlangt auch für `down` das Flag (Exit-Code `10`), und ein Projekt, das den Container ohne Instanz nie angelegt hat, kann ihn nicht über `--instance` ersetzen. `doctor`, `generate devcontainer` und `init` weisen die Fallbacks eines Instanz-Containers nicht aus; für sie ist ein Projekt mit nur Instanz-Containern „nicht ermittelbar“ ([`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) wird das nennen); das ist ein akzeptiertes Negativ, weil `doctor` ein Projekt prüft, nicht eine Instanz, und `u-boot` keine Liste von Instanzen führt. Der Race zwischen der Prüfung und dem Aufruf der CLI (ein Container wird dazwischen angelegt oder entfernt) ist nicht abgesichert: ein Nutzer, der denselben Ordner in zwei Terminals bedient, ist nicht Gegenstand.
- Offene Messung mit Rückführung: Dass `up` ohne Flag bei einem Container ohne Instanz und einer neueren Instanz den Container ohne Instanz trifft, ist im Standardprofil gemessen, im Sandbox-Profil nicht; für `down` und `shell` gilt die Suche von `u-boot` (Entscheidung 4), nicht die der CLI. Ein Docker-Test pinnt es; trifft die CLI dort die Instanz, trägt Entscheidung 4 für `up` nicht, und die Lösung (etwa beide Id-Labels der CLI, `devcontainer.local_folder` und `devcontainer.config_file`, auch ohne Instanz selbst mitzugeben) ist eine Folge-Entscheidung, die diese ADR ablöst. Die `devcontainerId` ändert sie nur, wenn die Labels nicht gleich sind; auch das ist dann zu messen.
- Nicht gemessen und im Test-Auftrag des Slice zu pinnen: wiederholtes `up --instance <name>` ohne `--recreate`, `down --instance <name>` mit unbekanntem Namen (erwartet: nichts zu stoppen, Exit-Code `0`, wie `down` ohne Container), die Suche nach einem Container ohne Instanz neben Instanzen im Sandbox-Profil und für `down` und `shell` (siehe Kontext).
- Folgepflicht: [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (Aufruf, Fehlerklassen, Verhalten ohne Flag), Handbuch Abschnitt 4.9, `CHANGELOG.md`; Docker-Test mit mindestens drei Instanzen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Tests (Unit, CLI) | Exit-Klasse je Sentinel und Unterbefehl (ungültiger Name, nur Instanzen); die Grenzfälle der Namensregel (Rand, Länge 32 und 33, `=`, leer, Großbuchstabe) | `make test` |
| Go-Test (Application) | bei nur Instanz-Containern ruft weder `up` noch `down` noch `shell` den Launcher | `make test` |
| Docker-Test | Container und Volumes von mindestens drei Instanzen paarweise verschieden; ohne Flag trifft der Befehl den Container ohne Instanz neben einer neueren Instanz | `make test-docker` |
| Go-Test (Port-Paket `driven`) | die Prüffunktion der Namensregel liegt im Port-Paket und lehnt die Grenzfälle ab; Application und Adapter rufen sie (je ein Test, der ohne gültigen Namen den Launcher beziehungsweise `docker` nicht erreicht) | `make test` |
| depguard | der Adapter `docker` importiert `domain` nicht (Importrichtung; der Ort der Funktion ist damit nicht geprüft) | `make lint` |

## Re-Evaluierungs-Trigger

- Eine neue Version der Dev-Container-CLI ändert, ob die `devcontainerId` mit einem zusätzlichen Id-Label wechselt, oder löst ein Id-Label anders auf (Punkt 1 des Kontexts).
- Der Docker-Test zur offenen Messung zeigt, dass `up` ohne Flag neben einer neueren Instanz die Instanz trifft.
- Der Projektinhaber will Instanzen auflisten, aufräumen oder im Standardprofil zusagen (Nicht-Ziele unten).

## Nicht-Ziele

`u-boot` führt keine Liste der Instanzen und räumt keine auf; es wählt keinen Zweig je Instanz, startet keinen Agenten und verteilt keine Arbeit; der Compose-Stack gehört nicht zu einer Instanz; mehr als die gemessene Zahl von Instanzen, andere Hosts als Linux `amd64` und die Isolation im Standardprofil sind nicht zugesagt (die Zusage steht im Lastenheft).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-31 | Proposed | — |
