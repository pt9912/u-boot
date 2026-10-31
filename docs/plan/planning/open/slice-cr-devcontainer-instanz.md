# Slice CR: Mehrere Devcontainer je Projekt — benannte Instanz an `up`, `down` und `shell`, Isolation je Instanz als Zusage nur für Gemessenes (`LH-FA-DEV-010`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (wie der Code-Slice; Release-Zuordnung beim Projektinhaber, die Roadmap ist nicht geändert). Zweiter von drei Slices: ihm geht [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) voraus (Messwerte), ihm folgt [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) (Code, Tests, Docker-Test, Doku). Vorbild des Ablaufs: [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md).

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes: nur Anwendung, keine Änderung), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (nennen „je Instanz“ für Workspace und Docker-Speicher; Änderung nur, wenn der Projektinhaber Stelle 5 annimmt).
Anlass: Projektinhaber, 2026-10-31 (Weg A des Entwurfs [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md); Anforderung: Instanzen teilen kein Volume).

**Berührte Spec-Stellen:** `lastenheft.md`, Abschnitt zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) folgt im Code-Slice.

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

**Ziel:** Die Zusage „ein Projekt kann mehrere benannte Devcontainer gleichzeitig haben; `up`, `down` und `shell` wählen den Container mit `--instance <name>`; jede Instanz hat eigene Volumes“ steht im Lastenheft, auf die **nächste freie Minor-Version** (`0.25.0`; gelesen am 2026-10-31: Kopf des Lastenhefts `0.24.0`, letzte Historie-Zeile `0.24.0` — **Vorbehalt:** vergeben wird beim Start des Slice, nach einem erneuten Blick auf die Historie, denn andere Entwürfe in `open/` können einen CR vor diesem schreiben), **mit den Messwerten des Mess-Slice** und in der vom Projektinhaber freigegebenen Fassung. Der Code-Slice beginnt erst danach.

**Dieser Wortlaut ist ein ENTWURF zur Freigabe durch den Projektinhaber.** Stellen, deren Inhalt vom Messwert abhängt, tragen einen Platzhalter `⟦M<Nr>: …⟧` (die Nummer ist die Messfrage in [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) §1). Ein Platzhalter wird vor der Vorlage zur Freigabe durch den Messwert ersetzt oder die Aussage gestrichen; ein Wortlaut mit Platzhalter geht nicht zur Freigabe. Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen Verweise als Anker `#…` der Datei.

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) beschreibt `up`, `down` und `shell` für den Devcontainer „des Projekts“ ohne Instanz; ein neues öffentliches Flag und eine Zusage zur Isolation (Volumes) sind Produktverhalten. Der Entwurf [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) sah einen CR bereits vor.

### Gefundener Stand (gelesen am 2026-10-31)

- [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 507): acht Spiegelstriche — `up` (mit `--recreate`), `down`, `shell`, Hilfs-Container und Socket, Zugangsdaten, Fehlerklassen, Zusage-Menge und Nicht-Ziele. Keiner nennt eine Instanz.
- **„Instanz“ ist im Lastenheft schon besetzt:** [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Zeile 473) sagt „Der Daemon-Speicher liegt in einem benannten Volume je Instanz“; gemeint ist der einzelne Container (`${devcontainerId}`), nicht ein vom Nutzer vergebener Name. Der Wortlaut unten führt deshalb die **benannte Instanz** und lässt „je Instanz“ in den Nachbarn unberührt, außer der Projektinhaber nimmt Stelle 5 an.
- [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Zeilen 1147 bis 1149): die Zeilen `devcontainer up`, `devcontainer down`, `devcontainer shell` nennen Zweck und Priorität, keine Flags (nicht geändert, Vorschlag; die Vorab-Suche bestätigt es am Start).
- [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeile 189) zählt `devcontainer up`, `down` und `shell` unter den ändernden Befehlen unter `--quiet`; das Flag ändert daran nichts (Vorschlag, nicht geändert).
- Das `--json`-Schema (`command: devcontainer`, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) ist **nicht gelesen**: ob die Ausgabe von `up` ein Feld für die Instanz bekommt, ist eine Frage der Spezifikation im Code-Slice; der Wortlaut unten sagt dazu nichts.
- Die Egress-, Ressourcen- und Zweig-Fragen berührt der CR nicht (siehe unten).

### Wortlaut (alt/neu; ENTWURF, nicht freigegeben)

**Stelle 1 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), erster Spiegelstrich (`up`), Satz zu `--recreate`.**

Alt:

> `--recreate` entfernt einen vorhandenen Container und legt ihn neu an; benannte Volumes bleiben erhalten.

Neu:

> `--recreate` entfernt einen vorhandenen Container und legt ihn neu an; benannte Volumes bleiben erhalten. Mit `--instance <name>` wählt der Befehl eine benannte Instanz des Devcontainers dieses Projekts: Er startet, sucht und ersetzt (mit `--recreate`) nur den Container dieser Instanz, und mehrere Instanzen desselben Projektordners können gleichzeitig laufen. Ohne `--instance` wirkt der Befehl ⟦M6: Entscheidung des Projektinhabers offen. Gemessen (Sandbox-Profil, Dev-Container-CLI 0.89.0, sechs Läufe) übernimmt die CLI mit nur dem Label `devcontainer.local_folder` neben laufenden Instanzen den zuletzt angelegten Container einer Instanz und legt keinen eigenen an; „nur auf den Container ohne Instanz“ ist **nicht** das Verhalten der CLI. Soll der Wortlaut „nur auf den Container ohne Instanz“ lauten, muss der Code-Slice das gegen die CLI herstellen (und `FindContainer` auf den Container ohne Instanz-Label einschränken); soll er „auf den zuletzt angelegten Container des Projekts“ lauten, steht das hier. Für das Standardprofil ist das Verhalten neben laufenden Instanzen nicht gemessen — kein Satz dazu⟧.

**Stelle 2 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), zweiter und dritter Spiegelstrich (`down`, `shell`).**

Alt (zweiter Spiegelstrich, ganz):

> `u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes.

Neu:

> `u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes. Mit `--instance <name>` stoppt er nur den Container dieser Instanz ⟦M7: die übrigen Container des Projekts laufen weiter⟧.

Alt (dritter Spiegelstrich, erster Satz):

> `u-boot devcontainer shell` startet den Devcontainer wie `u-boot devcontainer up` ohne `--recreate` und öffnet danach eine interaktive Shell im Container, mit dem Benutzer und dem Arbeitsverzeichnis, die der Start meldet.

Neu (nur der erste Satz; die folgenden Sätze bleiben):

> `u-boot devcontainer shell` startet den Devcontainer wie `u-boot devcontainer up` ohne `--recreate` (auch mit `--instance <name>`, dann den Container dieser Instanz) und öffnet danach eine interaktive Shell im Container, mit dem Benutzer und dem Arbeitsverzeichnis, die der Start meldet.

**Stelle 3 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), neuer Spiegelstrich nach dem zu `shell`: Isolation je Instanz (Zusage nur für das Gemessene).**

Alt: — (neu).

Neu:

> Bis zu ⟦M9: Zahl der gemessenen Instanzen, höchstens vier; Fassung A: „vier“, wenn die Messung vier Stufen vollständig trägt, Fassung B: die höchste vollständig gemessene Zahl, wenn der Host eine Stufe nicht trägt (Abbruch mit Befund im Mess-Slice)⟧ gleichzeitig laufende Instanzen desselben Projekts teilen im Sandbox-Profil untereinander kein Volume (alle Volumes paarweise verschieden): Jede hat ein eigenes Workspace-Volume ⟦M3: und, soweit aktiviert, eigene Volumes für den Speicher von Docker, den Speicher von Podman und die Modelle von Ollama — Familien nur aufzählen, die die Messung getrennt gezeigt hat⟧; eine Datei im Workspace der einen Instanz ist in der anderen nicht sichtbar. Gemessen ist das mit ⟦M9: Zahl der gemessenen Instanzen⟧ gleichzeitig laufenden Instanzen, ⟦M10: jede mit höchstens 1 CPU und 3 GB Arbeitsspeicher (Messbedingung je Instanz, Vorgabe des Projektinhabers; nur aufnehmen, wenn die zurückgelesene Grenze am Container im Protokoll steht)⟧, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine ⟦M0: Dev-Container-CLI-Version, Docker-Version, Speicher und CPU des Hosts nur, wenn die Zusage-Menge sie nennen soll — Vorschlag: nicht nennen, sie stehen im Handbuch⟧. Im Standardprofil teilen die Instanzen den Projektordner (Bind-Mount) ⟦M3: gemessen⟧; dort ist keine Isolation zugesagt. Nicht zugesagt sind mehr als ⟦M9: dieselbe Zahl⟧ gleichzeitige Instanzen, Instanzen mit mehr oder ohne Begrenzung von CPU und Arbeitsspeicher ⟦M10: nur mit M10 oben⟧, das Verhalten unter Last, die Größe und die Dauer des Starts weiterer Instanzen ⟦M4, M5: sonst Zahlen der Messung nennen oder streichen — Vorschlag: nicht nennen⟧, andere Hosts und Engines als gemessen sowie der Einsatz einer Instanz auf einem anderen Zweig. Kann `u-boot` die Isolation nicht herstellen ⟦M2: nur aufnehmen, wenn `${devcontainerId}` nicht mit dem zweiten Id-Label wechselt — dann Fassung (2) der Messung; wechselt es, entfällt dieser Satz⟧, endet der Befehl mit Exit-Code `10`.

**Stelle 4 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich „Fehlerklassen“.**

Alt (Anfang): „Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält; …“

Neu (zusätzlicher Halbsatz hinter dem ersten Fall; der Rest bleibt):

> Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält, **und wenn der Name der Instanz ungültig ist** (ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist ⟦M8: höchstens 32⟧ Zeichen lang; die leere Angabe ist ungültig); …

*Vorschlag mit Begründung (der Projektinhaber entscheidet):* **Zeichensatz `[a-z0-9-]`, Rand nur Buchstabe oder Ziffer, Länge 1 bis 32.** Der Name wird Label-Wert am Container und Suchwert im Filter von `docker ps`; ein geschlossener, kleiner Zeichensatz schließt Zeichen aus, die ein Filter, ein Shell-Aufruf oder ein Volume-Name anders lesen (Doppelpunkt, Komma, Gleichheitszeichen, Leerzeichen, Anführungszeichen); Kleinbuchstaben vermeiden zwei Instanzen, die sich nur durch Groß- und Kleinschreibung unterscheiden und auf Dateisystemen ohne Unterscheidung kollidieren würden; der Bindestrich-Rand folgt den Namensregeln von Docker; die Obergrenze 32 lässt Platz unter den Längengrenzen der Docker-Namen, wenn der Name in einen Namen eingeht ⟦M8: Obergrenze gegen die gemessene Grenze von Label-Wert und Docker-Namen prüfen und notfalls senken, nie anheben⟧. Gegenposition: weiter Zeichensatz (Unterstrich, Punkt, Großbuchstaben) für Namen wie `feature_x.1` — kleinerer Kreis ist leichter zu erweitern als zu verengen.

**Stelle 5 — [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Begriff „je Instanz“ (optional, Entscheidung des Projektinhabers).**

Alt (Beispiel [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer), Zeile 473): „Der Daemon-Speicher liegt in einem benannten Volume je Instanz.“

Neu (Vorschlag, wenn angenommen): „Der Daemon-Speicher liegt in einem benannten Volume je Devcontainer (je benannter Instanz, wenn `--instance` gesetzt ist).“ — die Fassung ersetzt „Instanz“ in den Nachbarstellen durch „Devcontainer“, damit der Begriff nur noch bei [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) die benannte Instanz meint. *Alternative: unverändert lassen* (der Wortlaut bleibt dann so, wie er ist, und „benannte Instanz“ grenzt ab); die Vorab-Suche zählt die Fundstellen von „Instanz“ im Lastenheft, bevor der Projektinhaber entscheidet.

**Stelle 6 — Nicht-Ziele, [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), letzter Spiegelstrich.**

Alt:

> Nicht-Ziele: `u-boot` startet keinen Agenten ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), und der Befehl startet nicht den Compose-Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)).

(Im Lastenheft stehen dort Anker-Links der Datei `#…`; hier der Pfad von `open/` aus.)

Neu (Anhang):

> … Nicht-Ziele zu Instanzen: `u-boot` verwaltet keine Liste der Instanzen und räumt keine auf, wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zur Instanz.

**Stelle 7 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.24.0 |` (Stand 2026-10-31).

Neu: `| Version          | 0.25.0 |` ⟦Vorbehalt: nächste freie Minor-Version beim Start, aus der Historie gelesen⟧.

Neue Historie-Zeile:

> | 0.25.0 | ⟦Datum der Freigabe⟧ | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (im Lastenheft als Anker `#…` der Datei): benannte Instanz an `u-boot devcontainer up`, `down` und `shell` (`--instance <name>`, mehrere Devcontainer desselben Projekts gleichzeitig; Isolation der Volumes im Sandbox-Profil für die gemessenen Aufbauten, Exit-Code `10` bei ungültigem Namen ⟦M2: und bei nicht herstellbarer Isolation⟧); nicht zugesagt: mehr als ⟦M9: Zahl der gemessenen Instanzen⟧ Instanzen, Last, andere Hosts als gemessen, Instanzen im Standardprofil isoliert, Zweig je Instanz. | Vereinbarung mit dem Projektinhaber |

### Überblick: wovon der Wortlaut vom Messwert abhängt

| Platzhalter | Stelle | Messfrage | Fällt der Messwert anders aus … |
|---|---|---|---|
| ⟦M0⟧ | 3 | Dev-Container-CLI- und Docker-Version | Zusage-Menge steht mit exakt diesen Versionen |
| ⟦M9⟧ | 3, 7 | Zahl der Stufen, die der Mess-Slice vollständig gemessen hat (Anforderung des Projektinhabers: Isolation mit drei und vier Instanzen) | trägt der Host vier, steht „vier“; sonst die höchste vollständig gemessene Zahl (Abbruch mit Befund, nicht Weglassen); die Zahl steht an allen drei Stellen gleich |
| ⟦M10⟧ | 3 | Messbedingung je Instanz `--cpus=1 --memory=3g` (Vorgabe des Projektinhabers; Annahme „je Instanz“, im Mess-Slice §1 ausgewiesen) | nennt das Protokoll die Grenze nicht zurückgelesen oder hat der Projektinhaber „für alle zusammen“ gemeint, entfällt oder ändert sich der Halbsatz; kein Versprechen für größere oder keine Limits |
| ⟦M2⟧ | 3, 4, 7 | wechselt `${devcontainerId}` mit dem zweiten Id-Label | wechselt es nicht: Stellen 3 und 7 nennen den Abbruch mit Exit-Code `10` oder der CR wartet auf einen Slice, der die Lücke schließt (Mess-Slice §4) |
| ⟦M3⟧ | 3 | getrennte Volumes je Familie, Standardprofil geteilt | nur die gemessenen Familien stehen im Satz |
| ⟦M4, M5⟧ | 3 | Image geteilt, Dauer, Platz | Vorschlag: keine Zahl im Lastenheft, nur im Handbuch |
| ⟦M6⟧ | 1 (und 2) | `up` ohne Flag neben einer Instanz | **gemessen:** die CLI übernimmt den zuletzt angelegten Instanz-Container (Sandbox-Profil, 6 von 6 Läufen); der Satz ist an das Verhalten angepasst oder an eine Zusage des Code-Slice, Entscheidung des Projektinhabers |
| ⟦M7⟧ | 2 | `down` stoppt nur die Instanz | sonst Satz streichen |
| ⟦M8⟧ | 4 | Grenze von Label-Wert und Docker-Namen | Obergrenze nur senken, nie anheben |

**Keine Änderung** (gelesen, bewusst nicht angefasst): [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) Regeln zum Hilfs-Container (die Instanz ändert den Socket-Mount nicht), die Zusage-Menge-Zeile von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („Linux-Host `amd64`, lokale Engine“; die neue Zusage bleibt innerhalb), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (siehe oben, Vorschlag).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung und Wahl der Lückenlösung** — [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md); der CR trägt nur, was dort gemessen ist.
- **Code, Tests, Spezifikation, Handbuch, `CHANGELOG.md`** — [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md); das Lastenheft ist Zielbild, der Code folgt.
- **Welchen Zweig eine Instanz bearbeitet, Agenten verteilen, Compose-Stack je Instanz, Last bei vielen Instanzen** — im Entwurf des Code-Slice ausgeschlossen (§1 dort); der CR führt sie als Nicht-Ziele oder gar nicht.
- **Entscheidung über Release oder Welle** — beim Projektinhaber (Roadmap nicht geändert).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Wortlaut ohne Platzhalter:** Alle `⟦…⟧` des Entwurfs sind durch Messwerte ersetzt oder die Aussagen gestrichen; der Wortlaut ist dem Projektinhaber Stelle für Stelle vorgelegt (Alt/Neu), der Vorschlag zu Zeichensatz und Länge (Stelle 4) und die Fragen zu Stelle 5 sind entschieden, die Freigabe ist in §7 vermerkt.
- [ ] **Lastenheft `0.25.0`** (oder die dann nächste freie Version) geschrieben: Stellen 1 bis 4, 6, 7 (und 5, wenn angenommen) in der freigegebenen Fassung, Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“).
- [ ] **Vorab-Suche und Nachzug:** `grep -i` über `spec/lastenheft.md` nach `Instanz`, `devcontainer up`, `devcontainer down`, `devcontainer shell`, `recreate`, `Volume`, `quiet` und `json` ist in §7 festgehalten; mitzuändernde Stellen: keine über die freigegebenen hinaus.
- [ ] `make gates` grün; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); der Review läuft über die **freigegebene** Fassung.
- [ ] Doku-Update: kein Handbuch, kein `CHANGELOG.md` (kein Produktverhalten; der Code-Slice zieht nach).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 bis 4, 6, 7 (und 5); gelesen: Kopf (Zeilen 1 bis 10), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 507), Fundstellen per `grep` (Zeilen 189, 191, 457, 473, 1147 bis 1149), Historie bis `0.24.0` (Zeile 1607) |
| `docs/plan/planning/open/slice-cr-devcontainer-instanz.md` (diese Datei) | update | Wortlaut ohne Platzhalter, Freigabe und Vorab-Suche in §7; Zustandsmarker auf Ist-Stand |
| `docs/plan/planning/open/slice-v1-devcontainer-instanz-label.md` | update | Zustandsmarker „CR freigegeben“ im Start-Trigger, Verweis-Nachzug der Moves |
| `docs/plan/planning/open/README.md` | update | Zeile dieses Slice auf „geliefert“; Verweis-Nachzug der Moves (`make slice-mv`) |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/` | update, neu | Register bei der Closure (§7) |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `internal/**` | unverändert | der CR ändert kein Produktverhalten; der Code-Slice zieht nach |

Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)): `grep -i` über `spec/lastenheft.md` nach den Begriffen der DoD; ausgeführt vor dem Lastenheft-Commit, Treffer in §7.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) liegt in `done/` (Messwerte und Urteil zur Kernfrage liegen vor); der Projektinhaber hat den Wortlaut freigegeben (Gespräch); kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Der Wortlaut verlangt mehr als eine Anforderung oder eine Entscheidung, die der Mess-Slice nicht getroffen hat (etwa die Lückenlösung bei Fassung (2) des Urteils).
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts, oder die Messung hat die Kernfrage nicht entschieden (Fassung (3)); das Lastenheft bleibt unverändert, kein Code-Slice startet.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Verification-Evidence und Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7). Der Lastenheft-Commit trägt die Freigabe nur, wenn §7 sie vermerkt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein Satz im Lastenheft stützt sich auf einen Messwert, den der Mess-Slice nicht liefert** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)); vor der Freigabe wird die Liste „nicht gemessen“ Satz für Satz gegen den Wortlaut gelesen (geschärfte Regel aus [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §7). — **Ausgang:** bei Closure.
- **Der Begriff „Instanz“ kollidiert mit „je Instanz“ in [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer).** — **Ausgang:** bei Closure (Stelle 5: Entscheidung des Projektinhabers).
- **Das Standardprofil teilt den Projektordner:** Eine Zusage „Instanzen sind isoliert“ ohne die Einschränkung auf das Sandbox-Profil wäre falsch. — **Ausgang:** bei Closure (Stelle 3 nennt die Grenze).
- **Die Version `0.25.0` ist vergeben, wenn der Slice startet:** andere Entwürfe in `open/` mit eigenem Lastenheft-CR können sie früher nehmen. — **Ausgang:** bei Closure (Version beim Start aus der Historie lesen).
- **Die Zusage nennt mehr Instanzen, als der Host des Mess-Slice getragen hat** (Anforderung: bis zu vier). — **Ausgang:** bei Closure (⟦M9⟧ wird gegen die Messwert-Tabelle des Mess-Slice §7 gelesen; nicht gemessene Stufen stehen nicht im Wortlaut).
- **Die Kernfrage fällt ungünstig aus** (Volumes geteilt): Der Wortlaut ändert sich oder der CR wartet auf einen Vorgang vor ihm. — **Ausgang:** bei Closure.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Freigabe des Wortlauts (Datum, Fassung) · Vorab-Suche mit Treffern · Review und Verification · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Grenzen · Drei Paarungen · Abgleich der Plan-Tabelle gegen `git diff --stat`.)*

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird allein `spec/` (DOK, Greenfield) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle (zwei von drei Achsen) ist laut Modus-Deklaration erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen. Gelesen: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) als Verweis aus [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §7 (Inhalt hier nicht gelesen; die Regel daraus steht in §6). [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) gilt (`AGENTS.md` §3.10): Jede Bedingung im Wortlaut („nur im Sandbox-Profil“, „kein geteiltes Volume“, „Exit-Code `10` bei ungültigem Namen“) braucht im Code-Slice einen pinnenden Test oder entfällt.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
