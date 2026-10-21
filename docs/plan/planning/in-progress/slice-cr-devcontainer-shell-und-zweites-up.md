# Slice CR: `u-boot devcontainer shell` und erneutes `up` nach dem Zustand des Containers — Lastenheft-CR (Fehlerklasse `11`, neuer Befehl, Befehlsübersicht, Abgrenzung des Hilfs-Containers, `--quiet`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [`welle-devcontainer-shell-und-zweites-up`](../welle-devcontainer-shell-und-zweites-up.md) (Release `v0.14.0`). Erster von sechs Slices (Reihenfolge CR → Verhalten von `up` → Fundament von `shell` → Befehl `shell` → Spezifikation → Release-Cut; Muster [`slice-cr-devcontainer-start-ohne-editor`](../done/slice-cr-devcontainer-start-ohne-editor.md)).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Fehlerklassen, neuer Befehl, Socket-Dauer), [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Befehlsübersicht), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (Abgrenzung des Hilfs-Containers), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (`--quiet`, Ausnahme 1; Fund der Vorab-Suche, §1); [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) (Entscheidung des Projektinhabers c) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) (Entscheidung des Projektinhabers e), beide `Accepted` seit `ec8fdf5a`.

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), §4.1 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §6.1 ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (das tut [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md) nach dem Code). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-21); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Lastenheft sagt zu, was [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) entschieden haben: ein `failed` aus einem früheren Start desselben Containers ergibt beim erneuten `up` wieder Exit-Code `11`, und der neue Befehl `u-boot devcontainer shell` (Start wie `up`, danach eine interaktive Shell im Container ohne Hilfs-Container und ohne Socket) steht in der Anforderung, in der Befehlsübersicht und in der Abgrenzung des Sandbox-Profils. Es ist das **Zielbild**: kein Messstand, kein Prozesszustand, keine Slice- oder ADR-Kennung. Ob der Code es trägt, entscheiden die Tests der Code-Slices, nicht dieser CR.

**Herkunft.** Entscheidungen des Projektinhabers vom 2026-10-21 im Gespräch: [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) c („ein Lastenheft-CR mit einem Satz in den Fehlerklassen“) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) e („Absatz zu `shell`, Socket ‚für die Dauer des Starts‘, Fehlerklassen, Zeile in der Befehlsübersicht, Klammer zum Hilfs-Container in der Abgrenzung des Sandbox-Profils, im selben CR“). Lücken L-1 und L-2 des [Validator-Berichts](../../../reviews/2026-10-21-validation-welle-devcontainer-start-ohne-editor.md).

### Lastenheft-Wortlaut-Vorschlag (alt/neu; vom Projektinhaber freigegeben am 2026-10-21, alle Stellen einschließlich 6b und 7, §7)

Gelesen am 2026-10-21 in `spec/lastenheft.md` (Version `0.19.0`, Zeile 9). Die Verweise stehen im Lastenheft als Anker-Links (`#lh-…`); hier sind sie aus diesem Verzeichnis gesetzt.

**Stelle 1 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich *Fehlerklassen* (Zeile 496; Entscheidung c von [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md)).**

- *Alt:* „Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält; `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind.“
- *Neu:* „Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält; `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `11` auch, wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`; `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind.“

**Stelle 2 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Einleitungssatz (Zeile 490) und neuer Spiegelstrich nach dem zu `down` (Zeile 493; Entscheidung e von [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md)).**

- *Alt (Einleitung):* „*Priorität: V2.* Das Produkt soll den erzeugten Devcontainer eines Projekts starten und stoppen können, ohne dass auf dem Host ein Devcontainer-fähiger Editor, Node.js oder die Dev-Container-CLI installiert ist; auf dem Host genügen Docker (Client und lokale Engine).“
- *Neu (Einleitung):* „*Priorität: V2.* Das Produkt soll den erzeugten Devcontainer eines Projekts starten, darin eine Shell öffnen und ihn stoppen können, ohne dass auf dem Host ein Devcontainer-fähiger Editor, Node.js oder die Dev-Container-CLI installiert ist; auf dem Host genügen Docker (Client und lokale Engine).“
- *Neu (Spiegelstrich nach `down`):* „`u-boot devcontainer shell` startet den Devcontainer wie `u-boot devcontainer up` ohne `--recreate` und öffnet danach eine interaktive Shell im Container, mit dem Benutzer und dem Arbeitsverzeichnis, die der Start meldet. Die Shell läuft ohne Hilfs-Container und ohne den Container-Runtime-Socket. Die Umgebungsvariablen, auf die die erzeugte Konfiguration mit `${localEnv:<NAME>}` verweist, sind in der Shell gesetzt, soweit sie auf dem Host gesetzt sind. Für den Start gelten die Fehlerklassen von `up`; meldet eine Statusdatei den Abbruch eines Startscripts, öffnet der Befehl keine Shell. Lässt sich die Shell im Container nicht starten, endet er mit Exit-Code `12`. Hat die Shell begonnen, endet der Befehl mit `0`, unabhängig von ihrem Exit-Status; einen Befehl im Container mit Rückgabe seines Exit-Codes führt er nicht aus.“

**Stelle 3 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zum Socket (Zeile 494; Entscheidung e von [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md): „Socket ‚für die Dauer des Starts‘“).**

- *Alt:* „Der Hilfs-Container bindet für die Dauer des Aufrufs den Container-Runtime-Socket des Hosts ein. Der Befehl weist diese Sicherheitslockerung in der Befehlsausgabe aus; …“
- *Neu:* „Der Hilfs-Container bindet den Container-Runtime-Socket des Hosts ein, bei `up` für die Dauer des Aufrufs, bei `shell` für die Dauer des Starts. Der Befehl weist diese Sicherheitslockerung in der Befehlsausgabe aus; …“ (Rest des Spiegelstrichs unverändert.)
- *Begründung:* Ohne die Änderung sagte der Spiegelstrich für `shell` „für die Dauer des Aufrufs“, und der Aufruf umfasst die Sitzung, in der kein Socket eingebunden ist ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1, *Kein Socket während der Sitzung*). `down` startet keinen Hilfs-Container (gelesen: `internal/hexagon/port/driven/devcontainer_launcher.go`, Methode `Stop`: „starts no helper container“).

**Stelle 4 — [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Befehlsübersicht (Zeilen 1137 und 1138).**

- *Alt:* die Zeilen `u-boot devcontainer up` (V2) und `u-boot devcontainer down` (V2), danach `u-boot template list`.
- *Neu:* eine Zeile nach `u-boot devcontainer down`: „| `u-boot devcontainer shell`  | Shell im Devcontainer öffnen       | V2        |“.

**Stelle 5 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), Spiegelstrich zum Socket (Zeile 448).**

- *Alt:* „keinen Container-Runtime-Socket des Hosts einbinden (der kurzlebige Hilfs-Container von `u-boot devcontainer up` gehört nicht zum erzeugten Ergebnis, siehe [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten));“
- *Neu:* „keinen Container-Runtime-Socket des Hosts einbinden (der kurzlebige Hilfs-Container von `u-boot devcontainer up` und `u-boot devcontainer shell` gehört nicht zum erzeugten Ergebnis, siehe [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten));“

**Stelle 6 — [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeilen 189, 191 und 196). Fund der Vorab-Suche, nicht in der Aufzählung der Entscheidung e; Freigabe einzeln (am 2026-10-21 freigegeben: 6a und 6b, 6b mit beiden Teilen, §7).**

- *6a (notwendig), Ausnahme 1, alt:* „… zusätzliche Devices, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up`) …“ — *neu:* „… zusätzliche Devices, der Einbindung des Container-Runtime-Sockets in den Hilfs-Container von `u-boot devcontainer up` und `u-boot devcontainer shell`) …“. *Begründung:* [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2 verlangt die Ausweisung von `shell` „auch unter `--quiet` ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Ausnahme 1)“; Ausnahme 1 nennt heute nur `up`. Ohne 6a müsste die Spezifikation die Ausnahme erweitern; sie präzisiert das Lastenheft und erweitert es nie ([`AGENTS.md`](../../../../AGENTS.md) §2).
- *6b (empfohlen), Liste der ändernden Subkommandos (Zeile 189), alt:* „(`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up`, `devcontainer down`)“ — *neu:* „(`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up`, `devcontainer down`, `devcontainer shell`)“; dazu nach dem Satz zu den Abfragebefehlen (Zeile 196) der neue Satz: „Die Shell, die `u-boot devcontainer shell` öffnet, ist nicht Ausgabe von `u-boot`; `--quiet` wirkt nur auf die Ausgabe des Starts.“ *Begründung:* `shell` startet den Container wie `up`; [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2: „`--quiet` hält den Fortschritt zurück wie bei `up`“. Ohne 6b sagt das Lastenheft nicht, was `--quiet` bei `shell` bewirkt. *Gegenoption:* 6b entfällt, und die Spezifikation beschreibt `--quiet` bei `shell` nur so weit, wie die Regel „Alle anderen Hinweise und Warnungen entfallen unter `--quiet`“ (Zeile 196) es trägt.

**Stelle 7 — optional, Vorschlag des Planners über Entscheidung c hinaus: [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zur Zusage-Menge (Zeile 497; am 2026-10-21 freigegeben, §7).**

- *Alt:* „… Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus, die der Befehl über die `u-boot.yaml` im Basisverzeichnis findet: ohne diese Datei endet der Befehl bei einem Startscript-Abbruch mit `12`. …“
- *Neu:* „… ohne diese Datei endet der Befehl bei einem Startscript-Abbruch mit `12`, und ein erneutes `up` meldet einen früheren Abbruch nicht. …“
- *Begründung:* [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) Entscheidung 1, Grenze (ii), und deren Wächter in der Fitness Function („ohne `u-boot.yaml` endet das zweite `up` mit 0“). Ohne Stelle 7 schweigt das Lastenheft zum zweiten Aufruf ohne `u-boot.yaml` ([`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md)). Entfällt die Stelle, trägt die Spezifikation die Grenze als Präzisierung der vorhandenen Voraussetzung.

**Stelle 8 — Versionsfeld (Zeile 9) und §16 Historie (nach Zeile 1589).** Versionsfeld auf die nächste freie Minor-Version (Vorschlag `0.20.0`; läuft [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) vorher, die darauf folgende; am 2026-10-21 freigegeben: `0.20.0`, dieser CR geht vor jenem nach `done/`, §7). Historie-Zeile, Vorschlag: „| 0.20.0 | ‹Datum› | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Exit-Code `11` auch, wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; neuer Befehl `u-boot devcontainer shell` (Start wie `up`, danach eine interaktive Shell im Container ohne Hilfs-Container und ohne Socket, mit den Umgebungsvariablen aus `${localEnv:<NAME>}`, Exit-Code `0` nach Beginn der Shell); Socket bei `shell` nur für die Dauer des Starts; [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): eine Zeile; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Abgrenzung des Hilfs-Containers auch für `shell`; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Socket-Ausweisung von `shell` als Teil von Ausnahme 1 (und, mit 6b, `devcontainer shell` unter den ändernden Befehlen). | Vereinbarung mit dem Projektinhaber |“ — die Zeile nennt nur die freigegebenen Stellen.

### Vorab-Suche (Fundstellen mit Befund, auch die unveränderten)

Gesucht am 2026-10-21 in `spec/lastenheft.md` nach `devcontainer up`, `devcontainer down`, `Hilfs-Container`, `subcommand`, `devcontainer` (Wert von `command`), `shell`, `docker exec`, `remoteEnv` und `Dauer des Aufrufs` ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md)):

| Fundstelle | Befund |
|---|---|
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeile 189 (ändernde Subkommandos) | Stelle 6b |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeile 191 (Ausnahme 1, Hilfs-Container von `up`) | Stelle 6a |
| [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeile 196 (Abfragebefehle, Vorschau) | Stelle 6b (neuer Satz danach); der Satz selbst unverändert |
| [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) Zeilen 230 und 234 (reservierte Codes) | unverändert; `shell` gibt den Exit-Status der Shell nicht weiter ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 3) |
| [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) Zeile 248 (Wertebereich von `command`, enthält `devcontainer`) | unverändert; `shell` ist ein Unterbefehl, kein neuer Wert von `command` |
| [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) Zeile 250 (`subcommand` bei gruppierten Befehlen, Beispiele `list`, `get`, `set`) | unverändert; Beispiele, keine Aufzählung |
| [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Zeile 448 | Stelle 5 |
| [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile 490 (Einleitung) | Stelle 2 |
| dieselbe, Zeilen 492 (`up`) und 493 (`down`) | unverändert; Stelle 2 setzt den neuen Spiegelstrich nach Zeile 493 |
| dieselbe, Zeile 494 (Socket, „für die Dauer des Aufrufs“) | Stelle 3 |
| dieselbe, Zeile 495 (Zugangsdaten nur als Umgebungsvariablen aus `${localEnv:<NAME>}`) | unverändert; gilt für `shell` wortgleich (`docker exec -e <NAME>` ohne Wert, [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 1, *Umgebung*) |
| dieselbe, Zeile 496 (Fehlerklassen) | Stelle 1 |
| dieselbe, Zeile 497 (Zusage-Menge, `u-boot.yaml`, Runtime-Image) | unverändert, gilt für `shell` („der Befehl“); optional Stelle 7 |
| dieselbe, Zeile 498 (Nicht-Ziele) | unverändert; eine Shell ist kein Agent |
| [`LH-FA-DEV-009`](../../../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer) Zeilen 482 bis 486 | unverändert; Übergabe zur Laufzeit per Umgebungsvariable gilt auch für die Shell |
| [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) Zeile 974 (`subcommand` verpflichtend bei `devcontainer`) | unverändert; der Umschlag des Exit-Codes `2` von `shell --json` trägt `subcommand: shell` ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2) |
| [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) Zeilen 1116 und 1117 | unverändert |
| [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) Zeilen 1137 und 1138 | Stelle 4 |
| §16 Historie, Zeile 1588 (`0.18.0`) | unverändert (Fußabdruck einer angenommenen Änderung) |
| `shell`, `docker exec`, `remoteEnv` | kein Treffer |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode, Port, Befehl, Tests** — der CR schreibt den Vertrag, der Code folgt ihm ([`slice-v1-devcontainer-up-failed-nach-zustand`](slice-v1-devcontainer-up-failed-nach-zustand.md), [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md), [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md)).
- **`spec/spezifikation.md` und `spec/architecture.md`** — nach dem Code in [`slice-v2-spec-devcontainer-shell-und-zweites-up`](slice-v2-spec-devcontainer-shell-und-zweites-up.md); [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) sagt bis dahin „nicht zugesagt“ zum `failed` aus einem früheren Lauf (Zeile 777).
- **Handbuch, READMEs und `CHANGELOG.md`** — es ändert sich kein Verhalten; die Code-Slices bringen die Dokumente mit.
- **`shell -- <Befehl>` mit durchgereichtem Exit-Code** — vom Projektinhaber verworfen ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md), Entscheidung des Projektinhabers c); er verlangte eine Ausnahme von den reservierten Codes in [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) und eine Folge-ADR.
- **Colima in der Zusage-Menge** — Folge-CR nach der Einzelmessung ([`slice-v1-devcontainer-start-colima-messung`](slice-v1-devcontainer-start-colima-messung.md)); Zeile 497 bleibt.
- **Die Messstände in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)** — eigener CR ([`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md)); wer zuerst läuft, bekommt die kleinere Versionsnummer.
- **Die Release-Nummer** — Lastenheft-Version ist nicht Release-Version; `v0.14.0` schneidet der Release-Cut der Welle.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Wortlaut im Lastenheft** nach den freigegebenen Stellen 1 bis 7 aus §1 (Stelle 6b und Stelle 7 nur, wenn einzeln freigegeben), wortgleich mit der Freigabe; Zielbild ohne Prozesszustand, ohne Messstand und ohne Verweis auf ADR, Slice oder Roadmap (Decken-Regel, [`AGENTS.md`](../../../../AGENTS.md) §3.4). *Beleg:* Vergleich Satz für Satz im Verification-Report, Fundstellen in §7.
- [ ] **Versionsfeld und Historie:** die nächste freie Minor-Version (Vorschlag `0.20.0`) im Versionsfeld und eine Historie-Zeile, die nur die freigegebenen Stellen nennt; die Priorität `V2` steht in der Zeile von [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gleich. *Beleg:* Versionsfeld, Historie-Zeile.
- [ ] **Freigabe des Projektinhabers** im Gespräch belegt und in §7 vermerkt, je Stelle (Stellen 6b und 7 einzeln); kein Wortlaut nach dem Review ohne Nachlauf-Review über genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).
- [ ] `make gates` grün, `make docs-check` grün. `make doc-trace` ist für [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) ein Beleg (Familie `LH-FA-*`), für [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) nicht: die Belegsicht führt nur `LH-FA-*` und `LH-QA-*` (`.d-check.yml`, Abschnitt `trace`; [`BEO-DOK/dod-sensor-sieht-anforderung-nicht`](../observations/BEO-DOK/dod-sensor-sieht-anforderung-nicht/observation.md)); dort trägt der Vergleich im Verification-Report.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Roadmap-Zeile „In Arbeit“ bei Start und zurück auf „Nichts in Arbeit.“ bei der Closure; Handbuch, READMEs und `CHANGELOG.md` unverändert (kein Verhalten); ADR-Index unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Die Tabelle nennt die Dateien, die der Diff voraussichtlich berührt; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Eine Zeile „prüfen“ nennt eine Datei, die der Planner nicht gelesen hat.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Stellen 1 bis 8 aus §1 im freigegebenen Wortlaut. Gelesen am 2026-10-21: Versionsfeld `0.19.0` (Zeile 9); [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeilen 189, 191, 196; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Zeile 448; [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeilen 488 bis 498; [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) Zeilen 1120 bis 1139 (letzte Zeile `u-boot template list`, V1); §16, letzte Historie-Zeile `0.19.0` (Zeile 1589) |
| `docs/plan/planning/in-progress/roadmap.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile „In Arbeit“ bei Start, „Nichts in Arbeit.“ bei der Closure (`.d-check.yml`, Abschnitt `planning`; `make doc-planning`) |
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile dieses Slice (Version des CR, Zustand); Zeile von [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) („nach 0.19.0“, gelesen Zeile 33) auf die Version nach diesem CR |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Versionsangaben (gelesen: Zeile 52 „nach 0.19.0“, Zeile 65 „zuletzt `0.19.0`“) auf die Version dieses CR, weil dieser Slice die nächste freie belegt; entfällt, wenn jener vorher läuft |
| diese Datei | update | `Verantwortlich:`, Freigabe-Vermerk, Plan-Abgleich, Closure; Moves per `make slice-mv` |
| `docs/plan/planning/welle-devcontainer-shell-und-zweites-up.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Verweis-Nachzug durch `make slice-mv`; die vergebene Lastenheft-Version in §5 |
| die übrigen Slice-Dateien der Welle in `open/` | update (nur Verweise) | `make slice-mv` zieht den Pfad dieses Slice nach |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report dieses Slice |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | je eine Evidence-Datei für jede angefallene Beobachtung (§7); keine, wenn keine anfällt |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/`, `internal/`, `cmd/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten. Gelesen als Stand, den die Folge-Slices ändern: [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeile 777 („welcher Exit-Code für ein `failed` aus einem früheren Lauf gilt, ist nicht zugesagt“); `docs/user/benutzerhandbuch.md` Zeilen 874 bis 883 (*Erneutes `up` nach einem Fehler oder Abbruch*) und 818 bis 833 (*Im Container arbeiten*); `internal/adapter/driving/cli/devcontainer.go`, `newDevcontainerCommand` (Unterbefehle `up`, `down`) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; [`ADR-0030`](../../adr/0030-devcontainer-up-failed-nach-zustand.md) und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) sind `Accepted` (eingetreten: `ec8fdf5a`); der Projektinhaber hat den Wortlaut der Stellen 1 bis 8 freigegeben, die Stellen 6b und 7 je einzeln mit Ja oder Nein; kein anderer Slice in `in-progress/` (WIP-Limit 1). Erfüllt am 2026-10-21: Stellen 1 bis 5, 6a und 8 freigegeben wie vorgeschlagen, 6b Ja, 7 Ja (§7); `in-progress/` ohne Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Projektinhaber verlangt eine Zusage, die eine neue Entscheidung braucht (etwa `shell -- <Befehl>` oder Colima); dann zuerst eine Folge-ADR.
- `in-progress` → `open` (blockiert): die Freigabe des Wortlauts steht aus; §7 nennt die offene Stelle.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; der Wortlaut im Lastenheft ist mit der Freigabe Satz für Satz verglichen (Verification-Report), der Review-Report liegt vor, und die Closure-Notiz mit Steering-Loop-Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Vorab-Suche übersieht eine Aufzählung** (Befehlsübersicht, Wertebereich, Liste der ändernden Befehle). Die Stellen 6a und 6b sind solche Funde außerhalb der Aufzählung der Entscheidung e. — **Ausgang:** weiter offen; die Fundstellen-Tabelle in §1 nennt jeden Treffer mit Befund, der Reviewer prüft sie gegen eine eigene Suche ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), 3×, Ausgang beim Projektinhaber).
- **Der Wortlaut ändert sich in der Fix-Runde nach dem Review.** — **Ausgang:** weiter offen; ändert eine Fix-Runde Lastenheft-Text, liest ein Nachlauf-Review im frischen Kontext genau diesen Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md), 7×, Ausgang beim Projektinhaber).
- **Die Zusage reicht weiter als die Tests.** Stelle 2 sagt eine *interaktive* Shell zu; [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4 nennt den Terminal-Pfad ohne eigenen Test „nicht zugesagt“ und als Prüfauftrag. — **Ausgang:** weiter offen; [`slice-v1-devcontainer-shell-fundament`](slice-v1-devcontainer-shell-fundament.md) und [`slice-v1-cli-devcontainer-shell`](slice-v1-cli-devcontainer-shell.md) liefern einen Test des Terminal-Pfads oder benennen seine Grenze; trägt kein Test das Wort „interaktiv“, folgt ein Folge-CR oder die Freigabe nimmt das Wort heraus ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 10×, Ausgang beim Projektinhaber). Frage an den Projektinhaber (Welle §5, Frage 2), entschieden am 2026-10-21: „interaktiv“ bleibt im Wortlaut; den Terminal-Pfad pinnt ein Test mit `script` aus dem Test-Image, ohne neue Abhängigkeit (Auftrag an Fundament- bzw. Befehl-Slice, §7).
- **Eine Bedingung des Wortlauts ohne Test** („ohne Hilfs-Container und ohne den Container-Runtime-Socket“, „soweit sie auf dem Host gesetzt sind“, „endet mit `0`, unabhängig von ihrem Exit-Status“, „ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`“). — **Ausgang:** weiter offen; die Fitness Functions beider ADRs nennen je Bedingung einen Test, die Code-Slices liefern ihn, die Beleg-Tabelle des Spezifikations-Slice führt ihn ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Stelle 6a wird nicht freigegeben.** Dann trägt das Lastenheft die Ausweisung von `shell` unter `--quiet` nicht, und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2 stünde ohne Anforderung. — **Ausgang:** entfallen; der Projektinhaber hat Stelle 6a am 2026-10-21 freigegeben (§7).
- **Versionsfolge mit [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md).** — **Ausgang:** weiter offen; wer zuerst nach `done/` geht, bekommt die kleinere Version, der andere zieht seine Angaben nach (§3). Entschieden am 2026-10-21: dieser CR geht zuerst nach `done/` und belegt `0.20.0`, der Messstände-CR rechnet danach mit der nächsten freien Version (§7).
- **Eine Aussage über den Code im Plan stimmt nicht.** — **Ausgang:** weiter offen; jede Aussage in §1 und §3 nennt Datei und Zeile oder Symbol ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(Wird bei der Closure gefüllt: Was hat funktioniert · Was ging anders als geplant · Freigabe-Vermerk je Stelle · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6. Die drei Paarungen prüft die Closure der Welle.)* Beim Start vermerkt sind nur die Entscheidungen des Projektinhabers.

**Entscheidungen des Projektinhabers (2026-10-21, im Gespräch; an den Implementer übermittelt durch den Hauptagenten, beim Start vermerkt).** Freigabe-Vermerk je Stelle:

1. **Stellen 1, 2, 3, 4, 5 und 6a — freigegeben** wie in §1 vorgeschlagen (*Neu*-Wortlaut unverändert).
2. **Stelle 6b — Ja**, beide Teile: `devcontainer shell` in der Liste der ändernden Subkommandos von [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und nach dem Satz zu den Abfragebefehlen der Satz „Die Shell, die `u-boot devcontainer shell` öffnet, ist nicht Ausgabe von `u-boot`; `--quiet` wirkt nur auf die Ausgabe des Starts.“
3. **Stelle 7 — Ja:** Spiegelstrich zur Zusage-Menge von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) „… mit `12`, und ein erneutes `up` meldet einen früheren Abbruch nicht.“
4. **Stelle 8 — freigegeben:** Version `0.19.0` → `0.20.0`; die Historie-Zeile nennt nur die freigegebenen Stellen, also mit 6b (`devcontainer shell` unter den ändernden Befehlen) und mit 7.
5. **Reihenfolge mit [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md):** dieser CR geht zuerst nach `done/` und belegt `0.20.0`; der Messstände-CR rechnet danach mit der nächsten freien Version (Versionsangaben dort und in seiner Zeile in `open/README.md` werden mit dem Lastenheft-Commit nachgezogen).
6. **Für die Folge-Slices** (in der Welle-Datei §5 vermerkt, hier nur der Zeiger): Terminal-Pfad per Test mit `script` aus dem Test-Image, „interaktiv“ bleibt (Frage 2); die Hilfe von `shell` ohne Anforderungskennung, `up` und `down` unverändert bis (j) als Ganzes entschieden ist (Frage 3); Frage 4 ((g)) bleibt offen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird `spec/` (Lastenheft) und `docs/plan/` (Planung), beide in der Sub-Area mit dem Kürzel DOK laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area) (Greenfield; auditiert, drei Achsen erfüllt). Keine Code-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-21; Zähler = Dateien in `evidence/`): [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (3, offen, Ausgang beim Projektinhaber) — trifft genau diesen CR, die Fundstellen-Tabelle in §1 ist die Antwort; [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (7, offen) — Risiko in §6; [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) (10, offen) — Risiko in §6 (Terminal-Pfad); [`BEO-DOK/vertrag-schweigt-zum-nutzerpfad`](../observations/BEO-DOK/vertrag-schweigt-zum-nutzerpfad/observation.md) (2) — die Welle greift sie im Closure-Kriterium auf, Stelle 7 ist ein Teil davon; [`BEO-DOK/dod-sensor-sieht-anforderung-nicht`](../observations/BEO-DOK/dod-sensor-sieht-anforderung-nicht/observation.md) (1) — DoD-Zeile zu `make doc-trace` in §2; [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (35, verkörpert) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (14, verkörpert) — Prüfpflicht, Risiken in §6. Evidenz-/Diskrepanz-Risiko: mittel (zwei Funde außerhalb der Aufzählung der Entscheidung).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
