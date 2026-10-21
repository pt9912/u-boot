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

**Berührte Spec-Stellen:** `lastenheft.md` §4.3 ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), §4.1 ([`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging)), §6.1 ([`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)), §16 (Historie), Versionsfeld. `SPEC-*` und `ARC-*` berührt der CR nicht (das tut [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../open/slice-v2-spec-devcontainer-shell-und-zweites-up.md) nach dem Code). Der Verweis zeigt aufwärts: Das Lastenheft nennt diesen Slice nie.

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

**Stelle 8 — Versionsfeld (Zeile 9) und §16 Historie (nach Zeile 1589).** Versionsfeld auf die nächste freie Minor-Version (Vorschlag `0.20.0`; läuft [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) vorher, die darauf folgende; am 2026-10-21 freigegeben: `0.20.0`, dieser CR geht vor jenem nach `done/`, §7). Historie-Zeile, Vorschlag: „| 0.20.0 | ‹Datum› | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Exit-Code `11` auch, wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; neuer Befehl `u-boot devcontainer shell` (Start wie `up`, danach eine interaktive Shell im Container ohne Hilfs-Container und ohne Socket, mit den Umgebungsvariablen aus `${localEnv:<NAME>}`, Exit-Code `0` nach Beginn der Shell); Socket bei `shell` nur für die Dauer des Starts; [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): eine Zeile; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil): Abgrenzung des Hilfs-Containers auch für `shell`; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging): Socket-Ausweisung von `shell` als Teil von Ausnahme 1 (und, mit 6b, `devcontainer shell` unter den ändernden Befehlen). | Vereinbarung mit dem Projektinhaber |“ — die Zeile nennt nur die freigegebenen Stellen.

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

- **Produktcode, Port, Befehl, Tests** — der CR schreibt den Vertrag, der Code folgt ihm ([`slice-v1-devcontainer-up-failed-nach-zustand`](../open/slice-v1-devcontainer-up-failed-nach-zustand.md), [`slice-v1-devcontainer-shell-fundament`](../open/slice-v1-devcontainer-shell-fundament.md), [`slice-v1-cli-devcontainer-shell`](../open/slice-v1-cli-devcontainer-shell.md)).
- **`spec/spezifikation.md` und `spec/architecture.md`** — nach dem Code in [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../open/slice-v2-spec-devcontainer-shell-und-zweites-up.md); [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) sagt bis dahin „nicht zugesagt“ zum `failed` aus einem früheren Lauf (Zeile 777).
- **Handbuch, READMEs und `CHANGELOG.md`** — es ändert sich kein Verhalten; die Code-Slices bringen die Dokumente mit.
- **`shell -- <Befehl>` mit durchgereichtem Exit-Code** — vom Projektinhaber verworfen ([`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md), Entscheidung des Projektinhabers c); er verlangte eine Ausnahme von den reservierten Codes in [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) und eine Folge-ADR.
- **Colima in der Zusage-Menge** — Folge-CR nach der Einzelmessung ([`slice-v1-devcontainer-start-colima-messung`](../open/slice-v1-devcontainer-start-colima-messung.md)); Zeile 497 bleibt.
- **Die Messstände in [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)** — eigener CR ([`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md)); wer zuerst läuft, bekommt die kleinere Versionsnummer.
- **Die Release-Nummer** — Lastenheft-Version ist nicht Release-Version; `v0.14.0` schneidet der Release-Cut der Welle.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Drei Liefer-Punkte; Gate-Läufe und Closure-Pflichten darunter zählen nicht mit.

- [ ] **Wortlaut im Lastenheft** nach den freigegebenen Stellen 1 bis 7 aus §1 (Stelle 6b und Stelle 7 nur, wenn einzeln freigegeben), wortgleich mit der Freigabe; Zielbild ohne Prozesszustand, ohne Messstand und ohne Verweis auf ADR, Slice oder Roadmap (Decken-Regel, [`AGENTS.md`](../../../../AGENTS.md) §3.4). *Beleg:* Vergleich Satz für Satz im Verification-Report, Fundstellen in §7.
- [ ] **Versionsfeld und Historie:** die nächste freie Minor-Version (`0.20.0`, belegt durch `d01fc191`) im Versionsfeld und eine Historie-Zeile, die nur die freigegebenen Stellen nennt; die Priorität `V2` steht in der Zeile von [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) gleich. *Beleg:* Versionsfeld, Historie-Zeile.
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
| `docs/plan/planning/open/README.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Zeile dieses Slice (Version des CR, Zustand); Zeile von [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md) („nach 0.19.0“, gelesen Zeile 33) auf die Version nach diesem CR |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Versionsangaben (gelesen: Zeile 52 „nach 0.19.0“, Zeile 65 „zuletzt `0.19.0`“) auf die Version dieses CR, weil dieser Slice die nächste freie belegt; entfällt, wenn jener vorher läuft |
| diese Datei | update | `Verantwortlich:`, Freigabe-Vermerk, Plan-Abgleich, Closure; Moves per `make slice-mv` |
| `docs/plan/planning/welle-devcontainer-shell-und-zweites-up.md` <!-- d-check:ignore (Pfad-Zitat) --> | update | Verweis-Nachzug durch `make slice-mv`; die vergebene Lastenheft-Version in §5 |
| die übrigen Slice-Dateien der Welle in `open/` | update (nur Verweise) | `make slice-mv` zieht den Pfad dieses Slice nach |
| `docs/reviews/` <!-- d-check:ignore (Pfad-Zitat) --> | neu | Review- und Verification-Report dieses Slice |
| `docs/plan/planning/observations/` <!-- d-check:ignore (Pfad-Zitat) --> | update (Closure) | je eine Evidence-Datei für jede angefallene Beobachtung (§7); keine, wenn keine anfällt |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `docs/maintainer/`, `internal/`, `cmd/` <!-- d-check:ignore (Pfad-Zitate) --> | **unverändert** | kein Code, kein Verhalten. Gelesen als Stand, den die Folge-Slices ändern: [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeile 777 („welcher Exit-Code für ein `failed` aus einem früheren Lauf gilt, ist nicht zugesagt“); `docs/user/benutzerhandbuch.md` Zeilen 874 bis 883 (*Erneutes `up` nach einem Fehler oder Abbruch*) und 818 bis 833 (*Im Container arbeiten*); `internal/adapter/driving/cli/devcontainer.go`, `newDevcontainerCommand` (Unterbefehle `up`, `down`) |

### Umsetzung (Implementer, 2026-10-21; nichts gepusht, kein Tag)

Commits (lokal, über `origin/main` = `282f04b6`): `770113e8` (Plan, Planner, mit der Eröffnung der Welle), `5e54f672` (Start: Entscheidungen des Projektinhabers in §7 und in der Welle-Datei §5, `Verantwortlich:`), `961e22cd` und `a5ede8a3` (`make slice-mv` nach `in-progress/`), `9e6c47d4` (Roadmap „In Arbeit“), `d01fc191` (Lastenheft `0.20.0`, Versionsverweise der Messstände-CR, Zeile dieses Slice in `open/README.md`), danach der Commit dieses Abschnitts. Jede Message seit `5e54f672` trägt die Kennung `slice-cr-devcontainer-shell-und-zweites-up`; alle außer den zwei Commits von `make slice-mv` tragen zusätzlich [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), `d01fc191` dazu [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (gelesen per `git log --format=%B 770113e8..HEAD`). Der `commit-msg`-Träger ist im Klon nicht aktiviert (`core.hooksPath` leer); jede Message-Datei lief vor dem Commit durch `tools/harness/commit-msg-traceability.sh` (Exit 0), die Range zusätzlich durch `make doc-commits` (unten).

**Plan vor Code — Ist gegen Plan** (gelesen auf `9e6c47d4`, vor dem Lastenheft-Commit). Die Fundstellen stimmen mit §1 und §3 überein: Version `0.19.0` Zeile 9; [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) Zeilen 189, 191, 196; [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Zeile 448; [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeilen 488 bis 498 (490 Einleitung, 492 `up`, 493 `down`, 494 Socket, 495 Zugangsdaten, 496 Fehlerklassen, 497 Zusage-Menge, 498 Nicht-Ziele); [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) Zeilen 1137 bis 1139; §16, letzte Zeile `0.19.0` (Zeile 1589). Jeder *Alt*-Wortlaut aus §1 stand zeichengleich im Lastenheft; ersetzt wurde je Stelle nur bei genau einem Treffer. Die Aussagen der Zeile „unverändert“ in der Plan-Tabelle gelesen und zutreffend: [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) Zeile 777 trägt den Satz „nicht zugesagt“; Handbuch Zeile 818 *Im Container arbeiten*, Zeile 874 *Erneutes `up` nach einem Fehler oder Abbruch*; `newDevcontainerCommand` hängt `newDevcontainerUpCommand` und `newDevcontainerDownCommand` an (`devcontainer.go`, `cmd.AddCommand`). Zwei Punkte, an denen der Plan dem Implementer eine Lesart ließ:

1. **Ort des Satzes von Stelle 6b.** §1 sagt „nach dem Satz zu den Abfragebefehlen (Zeile 196)“, die Fundstellen-Tabelle zu Zeile 196 (Abfragebefehle, Vorschau) „neuer Satz danach“. Unmittelbar hinter dem Satz zu den Abfragebefehlen stünde der neue Satz vor „Ebenso ist die Ausgabe der Vorschau-Optionen …“, und „Ebenso“ bezöge sich auf die Shell statt auf die Abfragebefehle. Der Satz steht deshalb als letzter Satz desselben Absatzes, hinter dem Satz zu den Abfragebefehlen und den zwei Sätzen zur Vorschau; sein Wortlaut ist zeichengleich mit der Freigabe. **Offen:** Ob dieser Ort der Freigabe entspricht, bestätigt der Projektinhaber; Reviewer und Verifier prüfen ihn.
2. **Historie-Zeile `0.20.0`.** Der Vorschlag in §1 (Stelle 8) trägt für 6b die Klammer „(und, mit 6b, `devcontainer shell` unter den ändernden Befehlen)“ und für Stelle 7 keinen Teilsatz. Geschrieben ist 6b als „… als Teil von Ausnahme 1 und `devcontainer shell` unter den ändernden Befehlen“; Stelle 7 als Klammer im Teil zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): „(ohne `u-boot.yaml` meldet ein erneutes `up` einen früheren Abbruch nicht)“, aus dem freigegebenen Wortlaut von Stelle 7 genommen, weil die Zeile die freigegebenen Stellen nennt. Den Satz von 6b zur Shell nennt die Zeile nicht eigens (der Vorschlag sah ihn nicht vor). **Offen:** Die Klammer zu Stelle 7 steht nicht wörtlich im Vorschlag; Reviewer und Projektinhaber prüfen sie.

**Wortlaut-Vergleich** (Python auf dem Host als Text-Werkzeug, kein Gate): Jeder *Neu*-Wortlaut aus §1 (Stellen 1 bis 7; Verweise `../../../../spec/lastenheft.md#…` als `#…`, „…“ als Trennstelle) ist auf `d01fc191` als Teilzeichenkette im Lastenheft enthalten, dazu die Zeile der Befehlsübersicht, die Liste der ändernden Subkommandos und der Satz von 6b in der Form, die der Hauptagent übermittelt hat. `V2` steht in der neuen Zeile von [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) und in der Einleitung von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten). Kein Verweis nach unten, keine `SPEC-` oder `ARC-`Kennung im neuen Text (Mutationen M2 und M3 unten belegen, dass `docs-check` beides fängt).

**Vorab-Suche** (Host-`grep` als Text-Werkzeug, kein Gate; Muster aus §1 *Vorab-Suche*, dazu `interaktiv`, `Dauer des Starts`, `Subkommando`, `Unterbefehl`, `` `down` ``, `ändernden`, `Abfragebefehl`, `Befehlsausgabe`, `Exit-Status`, `Standardausgabe`, `stdout`, `stderr`):

- *Lastenheft vor dem Commit* (`9e6c47d4`): Jeder Treffer steht in der Fundstellen-Tabelle in §1, außer diesen; gelesen, keine Aufzählung der Devcontainer-Befehle, unverändert: Zeile 135 (Modul-Tabelle, Devcontainer-Generator); [`LH-FA-GEN-001`](../../../../spec/lastenheft.md#lh-fa-gen-001--generate-befehl) Zeile 739 (Artefakt `devcontainer` von `generate`); [`LH-AK-005`](../../../../spec/lastenheft.md#lh-ak-005--devcontainer-flow) Zeilen 1340 bis 1349 (`init --devcontainer`, Inhalt von `devcontainer.json`); [`LH-MVP-001`](../../../../spec/lastenheft.md#lh-mvp-001--muss-im-mvp-enthalten-sein) Zeilen 1443 bis 1456 (Umfang des MVP, ohne die V2-Befehle); [`LH-FA-CLI-008`](../../../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) Zeile 258 („dateiverändernde Befehle“; `shell` schreibt keine Datei); [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeile 210 (Befehle mit Bestätigung; `shell` fragt nichts). `shell`, `docker exec`, `remoteEnv`, `Standardausgabe`, `stdout`, `stderr`: kein Treffer.
- *Lastenheft nach dem Commit* (`d01fc191`): Der Befehlsname `shell` steht in den Zeilen 189, 191, 196, 448, 494, 495, 1140 und 1592, jede eine freigegebene Stelle; „Dauer des Aufrufs“ nur noch in Zeile 495, neben „Dauer des Starts“.
- *Spezifikation, Sicht, Handbuch, READMEs, `docs/maintainer/cli-json-contract.md`, `CHANGELOG.md`* (Muster `devcontainer up`, `devcontainer down`, Wertebereich von `subcommand` mit `up` und `down`, „für die Dauer“): in diesem Slice unverändert (kein Verhalten). Die meisten Treffer führen die Pläne der Folge-Slices schon in ihren Fundstellen-Tabellen ([`slice-v2-spec-devcontainer-shell-und-zweites-up`](../open/slice-v2-spec-devcontainer-shell-und-zweites-up.md) §1 und §3: [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), [`SPEC-043`](../../../../spec/spezifikation.md#spec-043--verbosity-pegelabbildung-rangfolge-und-wirkung-von---quiet) Zeilen 584 bis 595, `architecture.md` §2.2 bis §2.6; [`slice-v1-cli-devcontainer-shell`](../open/slice-v1-cli-devcontainer-shell.md) §3: Handbuch, beide READMEs, `CHANGELOG.md`, `cli-json-contract.md` §6.11; [`slice-v1-devcontainer-up-failed-nach-zustand`](../open/slice-v1-devcontainer-up-failed-nach-zustand.md) §3: Handbuch 4.9 und §7). Status-Block der READMEs (Zeile 11 bzw. 12) und die Einträge zu `v0.13.0` im `CHANGELOG.md` beschreiben das Release und bleiben bis zum Release-Cut. **Treffer, die kein Folge-Plan nennt — Aufträge:**
  - an [`slice-v2-spec-devcontainer-shell-und-zweites-up`](../open/slice-v2-spec-devcontainer-shell-und-zweites-up.md): [`SPEC-045`](../../../../spec/spezifikation.md#spec-045--statusdateien-der-startscripte-und-ihre-ausweisung) *Ermittlung durch den Host* (`spezifikation.md` Zeile 648: „`doctor`, `generate devcontainer`, `init` und `devcontainer up` lesen die Dateien über die Container-Engine“; die Startphase von `shell` liest sie ebenso); der Absatz unter dem Schema von [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) (Zeile 936: „… von `u-boot devcontainer up` und `u-boot devcontainer down`; der Befehl setzt `subcommand` auf `up` beziehungsweise `down`“; der Plan führt nur Zeile 830); [`SPEC-024`](../../../../spec/spezifikation.md#spec-024--beispielaufrufe-zu-den-anforderungen) *Beispielaufrufe* (Zeile 1181, Zeile zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).
  - an [`slice-v1-cli-devcontainer-shell`](../open/slice-v1-cli-devcontainer-shell.md): Handbuch 4.9 *Sicherheit: das Docker-Socket* (Zeilen 770 bis 773: „bindet für die Dauer des Aufrufs“, „`u-boot devcontainer up` sagt das, bevor …“; Gegenstück von Stelle 3); Handbuch 4.9, Spiegelstrich `--json` (Zeile 811: „`subcommand` ist `up` oder `down`“; der Plan führt nur Zeile 2548); `cli-json-contract.md` §7 *Mutations-Matrix pro Subcommand* (Zeilen 1328 und 1329, je eine Zeile für `devcontainer up` und `down`).
  - Aufzählungen im Code (Hilfe, Fehlertexte, Kommentare unter `internal/`) lagen außerhalb dieser Suche; der Befehl-Slice sucht sie selbst (dort §6, Risiko „Eine Aufzählung der Befehle bleibt ohne `shell`“).

**Mutationen** (je ein `make docs-check`, danach `git checkout -- spec/lastenheft.md`; `git status` danach sauber):

| # | Mutation | Ergebnis | Ausgabe (Auszug) |
|---|---|---|---|
| M0 | Stand `a5ede8a3` (Slice in `in-progress/`, Roadmap noch mit dem Ruhe-Marker) | **rot**, Exit 2 | `roadmap.md:14 … planning-drift … trägt den Ruhe-Marker „Nichts in Arbeit.“` — behoben durch `9e6c47d4` |
| M1 | Anker der neuen Historie-Zeile auf [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) verfälscht (`…-befehlx`) | **rot**, Exit 2 | `spec/lastenheft.md:1592 … anchor-missing` |
| M2 | im neuen Spiegelstrich zu `shell` ein Verweis auf die ADR-Datei von [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) | **rot**, Exit 2 | `spec/lastenheft.md:494 … matrix-forbidden … contract-spec → adr` |
| M3 | im Spiegelstrich zum Socket ein Verweis auf [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) | **rot**, Exit 2 | `spec/lastenheft.md:495 … matrix-forbidden … contract-spec → tech-spec` |

**Grenze der Sensoren:** Kein Gate vergleicht den Wortlaut mit der Freigabe (das tun der Vergleich oben und der Verification-Report), keines prüft den Bezug von „Ebenso“ (Punkt 1 oben), und `make doc-trace` führt [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) nicht ([`BEO-DOK/dod-sensor-sieht-anforderung-nicht`](../observations/BEO-DOK/dod-sensor-sieht-anforderung-nicht/observation.md)).

**Sensoren** (Docker, Host `linux/amd64`; Exit-Codes der `make`-Prozesse):

| Sensor | Ergebnis |
|---|---|
| `make gates` (auf `d01fc191`) | Exit 0, 67 s: `lint` 0 issues; `test` grün (24 Pakete `ok`); `coverage-gate` 91,70 % gegen 90 %; `docs-check` 619 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK, 54 Dateien; `ablage-pruefung` OK, Wächter-Test 26 Fälle, 0 Fehlschläge; `record-gates` stempelte `d01fc191` |
| `make docs-check` | vor den Commits `5e54f672`, `9e6c47d4`, `d01fc191` und dieses Abschnitts je 0 Befunde; rot auf `a5ede8a3` (M0), wie im Ablauf erwartet |
| `make doc-trace` | Exit 0, 84 Anforderungen, 0 Waisen; [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) je `ok`. **Grenze:** [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) steht nicht in der Belegsicht (nur `LH-FA-*` und `LH-QA-*`, `.d-check.yml` `trace`); dort trägt der Vergleich im Verification-Report |
| `make doc-immutable RANGE="282f04b6..HEAD"` | Exit 0, 8 Commits bis `d01fc191`, 0 Befunde; dazu `git diff ec8fdf5a..HEAD -- docs/plan/adr/` leer (beide ADRs führen den Status inline, [`BEO-DOK/doc-immutable-greift-nicht-bei-inline-status`](../observations/BEO-DOK/doc-immutable-greift-nicht-bei-inline-status/observation.md)) |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0, 8 Commits bis `d01fc191`, 0 Befunde |
| `make test-docker` | **nicht gelaufen:** kein Docker-Pfad, kein Code, kein Test berührt (`git diff --stat 770113e8..HEAD -- internal cmd` leer) |
| `make verify-depguard` | **nicht gelaufen:** keine Importregel und kein Import berührt |

**Kommentar- und Marker-Selbstcheck** ([`AGENTS.md`](../../../../AGENTS.md) §3.7, §3.9, §3.10) über `git diff 770113e8..HEAD`: Der Diff berührt nur Markdown, keinen Code-Kommentar. (a) „früher“ steht in den hinzugefügten Zeilen nur im Vertragstext („Abbruch aus einem früheren Start desselben Containers“, „meldet einen früheren Abbruch nicht“) und in seinen Zitaten; es beschreibt den Zustand des Containers, nicht abwesenden Text. (b) Kein Symbol und keine Datei umbenannt; die Versionszahl `0.19.0` als „nächste freie“ steht nur noch in Lese-Aussagen mit Datum (§1, §3), in der Messstände-CR ist sie nachgezogen. (c) entfällt. Abgelaufene Zustandsmarker auf den Ist-Stand gesetzt: DoD *Versionsfeld* („Vorschlag `0.20.0`“ → „`0.20.0`, belegt durch `d01fc191`“), Zeile dieses Slice in `open/README.md`, Tabellenzeile und Freigabe (2) in der Welle-Datei, Risiko zum Terminal-Pfad („Offene Frage“ → entschieden). Neu mit Ablaufpunkt: „in Arbeit (liegt in `in-progress/`)“ in `open/README.md` und die Roadmap-Zeile — beide setzt die Closure zurück.

**Abweichungen vom Plan.** (1) Ort des Satzes von 6b und (2) Klammer zu Stelle 7 in der Historie-Zeile, beide oben unter *Plan vor Code* begründet und als offen markiert. (3) Die Welle-Datei bekam beim Start mehr als den Verweis-Nachzug und die Version in §5: Vermerke zu den Fragen 1 bis 4, die erteilte Freigabe (2) und die Tabellenzeile des CR ohne „Vorschlag“ (Auftrag des Hauptagenten; Ursprungstext der Fragen unverändert). (4) Von den „übrigen Slice-Dateien der Welle“ zog `make slice-mv` drei nach (Fundament, Urteil nach dem Zustand, Spezifikation); [`slice-v1-cli-devcontainer-shell`](../open/slice-v1-cli-devcontainer-shell.md) verweist nicht auf diesen Slice. Sonst nach Plan; kein Produktcode, Spezifikation, Sicht, Handbuch, READMEs und `CHANGELOG.md` unverändert.

**Plan-Tabelle gegen `git diff --stat 770113e8..HEAD`** (Stand `d01fc191`, vor diesem Abschnitt): 9 Dateien — `docs/plan/planning/in-progress/roadmap.md`, dieser Slice (Rename aus `open/`), `docs/plan/planning/open/README.md`, `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md`, die drei Slice-Dateien der Welle aus Punkt (4), die Welle-Datei, `spec/lastenheft.md`; jede hat einen Eintrag. Einträge ohne Diff: `docs/reviews/` (zu diesem Stand ohne Report), das Beobachtungs-Register (Closure); die als „unverändert“ geführten Pfade liegen nicht im Diff (`git diff --stat 770113e8..HEAD` über `spec/spezifikation.md`, `spec/architecture.md`, `docs/user`, beide READMEs, `CHANGELOG.md`, `docs/maintainer`, `internal`, `cmd`, `docs/plan/adr` leer). Wiederholt wird der Abgleich nach der Fix-Runde und bei der Closure.

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
- **Die Zusage reicht weiter als die Tests.** Stelle 2 sagt eine *interaktive* Shell zu; [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 4 nennt den Terminal-Pfad ohne eigenen Test „nicht zugesagt“ und als Prüfauftrag. — **Ausgang:** weiter offen; [`slice-v1-devcontainer-shell-fundament`](../open/slice-v1-devcontainer-shell-fundament.md) und [`slice-v1-cli-devcontainer-shell`](../open/slice-v1-cli-devcontainer-shell.md) liefern einen Test des Terminal-Pfads oder benennen seine Grenze; trägt kein Test das Wort „interaktiv“, folgt ein Folge-CR oder die Freigabe nimmt das Wort heraus ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md), 10×, Ausgang beim Projektinhaber). Frage an den Projektinhaber (Welle §5, Frage 2), entschieden am 2026-10-21: „interaktiv“ bleibt im Wortlaut; den Terminal-Pfad pinnt ein Test mit `script` aus dem Test-Image, ohne neue Abhängigkeit (Auftrag an Fundament- bzw. Befehl-Slice, §7).
- **Eine Bedingung des Wortlauts ohne Test** („ohne Hilfs-Container und ohne den Container-Runtime-Socket“, „soweit sie auf dem Host gesetzt sind“, „endet mit `0`, unabhängig von ihrem Exit-Status“, „ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`“). — **Ausgang:** weiter offen; die Fitness Functions beider ADRs nennen je Bedingung einen Test, die Code-Slices liefern ihn, die Beleg-Tabelle des Spezifikations-Slice führt ihn ([`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10).
- **Stelle 6a wird nicht freigegeben.** Dann trägt das Lastenheft die Ausweisung von `shell` unter `--quiet` nicht, und [`ADR-0031`](../../adr/0031-devcontainer-shell-ohne-editor.md) Entscheidung 2 stünde ohne Anforderung. — **Ausgang:** entfallen; der Projektinhaber hat Stelle 6a am 2026-10-21 freigegeben (§7).
- **Versionsfolge mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md).** — **Ausgang:** weiter offen; wer zuerst nach `done/` geht, bekommt die kleinere Version, der andere zieht seine Angaben nach (§3). Entschieden am 2026-10-21: dieser CR geht zuerst nach `done/` und belegt `0.20.0`, der Messstände-CR rechnet danach mit der nächsten freien Version (§7).
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
5. **Reihenfolge mit [`slice-cr-lastenheft-messstaende`](../open/slice-cr-lastenheft-messstaende.md):** dieser CR geht zuerst nach `done/` und belegt `0.20.0`; der Messstände-CR rechnet danach mit der nächsten freien Version (Versionsangaben dort und in seiner Zeile in `open/README.md` werden mit dem Lastenheft-Commit nachgezogen).
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
