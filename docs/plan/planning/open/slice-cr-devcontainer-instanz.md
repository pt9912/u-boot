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

**Bezug:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Änderung), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Codes: nur Anwendung, keine Änderung), [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Sandbox-Profil; nur DEV-007 sagt „je Instanz“, für den Docker-Speicher, Zeile 473; beide unverändert, siehe Entscheidung 2).
Anlass: Projektinhaber, 2026-10-31 (Weg A des Entwurfs [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md); Anforderung: Instanzen teilen kein Volume).

**Berührte Spec-Stellen:** `lastenheft.md`, Abschnitt zu [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), §16 (Historie), Versionsfeld. `spezifikation.md` und `architecture.md` berührt der CR nicht; [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) folgt im Code-Slice.

**Verantwortlich:** pt9912 (gesetzt beim Start, 2026-10-31)

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-31.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Zusage „ein Projekt kann mehrere benannte Devcontainer (Instanzen) gleichzeitig haben; `up`, `down` und `shell` wählen den Container mit `--instance <name>`; im Sandbox-Profil hat jede Instanz eigene Volumes“ steht im Lastenheft, mit den Messwerten des Mess-Slice und in der vom Projektinhaber freigegebenen Fassung. Die Versionsnummer (voraussichtlich `0.25.0`; gelesen am 2026-10-31: Kopf des Lastenhefts `0.24.0`, letzte Historie-Zeile `0.24.0`, Zeilen 9 und 1607) wird erst bei der Annahme aus der dann gültigen Historie gelesen. Der Code-Slice beginnt erst danach.

**Dieser Wortlaut ist ein ENTWURF zur Freigabe durch den Projektinhaber.** Er steht nur in dieser Datei; `spec/lastenheft.md` ist nicht geändert und wird erst nach der Freigabe geschrieben. Alle Messwert-Platzhalter sind aus dem Protokoll [`docs/messungen/slice-v1-devcontainer-instanz-messung`](../../../messungen/slice-v1-devcontainer-instanz-messung/README.md), der Closure-Notiz des Mess-Slice (§7) und dem Verification-Report [`2026-10-31-slice-v1-devcontainer-instanz-messung-verify`](../../../reviews/2026-10-31-slice-v1-devcontainer-instanz-messung-verify.md) (Abschnitt 5, Platzhalter des CR) ersetzt. **Offen bleiben nur die Versionsnummer und das Datum der Historie-Zeile** (Abschnitt „Entscheidungen des Projektinhabers“). Die Spec-Dateien sind Zielbild: Aussagen über noch ungemessene Teile, etwa `down --instance`, dürfen als Zielbild im Wortlaut stehen. Das Lastenheft nennt keine Slice- oder Wellen-Kennung; im Lastenheft stehen Verweise als Anker `#…` der Datei, hier als Pfad von `open/` aus.

### Entscheidungen des Projektinhabers (2026-10-31)

1. **Namensregel der Instanz:** Zeichen `a` bis `z`, `0` bis `9` und `-`; Rand (erstes und letztes Zeichen) nur Buchstabe oder Ziffer; Länge 1 bis 32; ein ungültiger Name endet mit Exit-Code `10`. Die leere Zeichenkette und ein Name mit `=` sind ausgeschlossen (gemessen: die Dev-Container-CLI lehnt die leere Zeichenkette als Label-Wert mit Exit-Code `1` ab; bei `=` im Wert kehrt `devcontainer up` nicht zurück, Ursache nicht geklärt; Protokoll M8 und Abschnitt Hänger). Die Obergrenze 32 ist eine Entscheidung, kein Messwert: gemessen sind Label-Werte bis Länge 128 ohne Befund, eine Grenze darunter hat die Messung nicht gefunden.
2. **Begriff „Instanz“:** [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) bleiben unverändert (dort heißt „je Instanz“ der einzelne Devcontainer). Der CR definiert „Instanz“ einmal, in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Stelle 1). Das Glossar (§15, Zeile 1542 ff.) hat keinen Eintrag zu „Instanz“ (gelesen: Zeilen 1542 bis 1556; `grep -i instanz` trifft das Lastenheft nur in den Zeilen 416, 473 und der Historie), also kein Glossar-Eintrag.
3. **Flag:** `--instance <name>` an `up`, `down` und `shell`.
4. **`up`, `down` und `shell` ohne `--instance` (Messfrage M6): Fassung (i).** Ohne `--instance` wirken die Befehle nur auf den Container ohne Instanz (Bestandscontainer und Standard). Gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `up` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ebenso `down` und `shell` (Entscheidung desselben Tages, symmetrisch zu `up`); `u-boot` prüft das vor dem Aufruf der Dev-Container-CLI, weil die CLI sonst eine Instanz übernähme (gemessen für `up`; ob `down` und `shell` die CLI vor der Prüfung aufrufen, ist nicht gelesen). Die Fassung, nach der `up` wie die CLI den zuletzt angelegten Container übernimmt, ist verworfen.
5. **Spec-Dateien sind Zielbild:** Aussagen über noch ungemessene Teile (etwa `down --instance`) dürfen als Zielbild im Wortlaut stehen. Was gemessen ist und was nicht, steht in diesem Slice (Begleittexte, Tabelle „Nicht-Zusagen“), nicht im Vertragstext.

**Offen bleiben nur:** die Versionsnummer (aus der Historie lesen; vorläufig `0.25.0`, Kopf des Lastenhefts heute `0.24.0`, letzte Historie-Zeile `0.24.0`, Zeilen 9 und 1607 — am Ist-Stand gelesen, keine Lastenheft-Änderung seit `0.24.0`) und das Datum der Historie-Zeile (Tag der Annahme).

**Warum eine Lastenheft-Änderung:** [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) beschreibt `up`, `down` und `shell` für den Devcontainer „des Projekts“ ohne Instanz; ein neues öffentliches Flag und eine Zusage zur Isolation (Volumes) sind Produktverhalten. Der Entwurf [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) sah einen CR bereits vor.

**Begründung zu Entscheidung 4 (Fassung (i)).** Gemessen ist (Protokoll M2 und M6): Die Dev-Container-CLI hängt die Volume-Namen an `${devcontainerId}`, und der Wert hängt von den Id-Labels ab — ohne `--id-label` ergibt einen anderen Wert als mit nur `--id-label devcontainer.local_folder=<Ordner>` (alle sechs Läufe „Gleich: nein“). Gäbe `u-boot` dem Standard-Container ein eigenes Id-Label, änderte sich seine `devcontainerId` und damit die Namen seiner Volumes: Bestehende Projekte fänden ihren Workspace nicht wieder. Deshalb bleibt der Container ohne Instanz label-frei, und die Trennung läuft über das Instanz-Label. Gemessen ist außerdem: Ist ein Container ohne Flag da, trifft `up` ohne Flag ihn weiter, auch wenn danach eine Instanz angelegt wurde; gibt es nur Instanzen, **übernimmt** die CLI mit `up` ohne Flag den zuletzt angelegten Instanz-Container (Sandbox-Profil, 6 von 6 Läufen, Vermutung zum Mechanismus nicht gemessen) — in dem Fall würde ohne eigene Prüfung eine Instanz ohne Namen angesprochen. Hinzu kommt: Die Suche nach dem ersten Treffer von `docker ps -a` trifft bei einer neueren Instanz neben einem Bestandscontainer die Instanz (Protokoll M6, umgekehrt), während `up` der CLI den Bestandscontainer bevorzugt — `down` und `shell` ohne Flag und `up` würden also auf verschiedene Container wirken, solange die Suche nicht auf den Container ohne Instanz eingeschränkt ist; das ist ein Auftrag des Code-Slice.

**Prüfung von (i) am Code** (gelesen am 2026-10-31):

- `FindContainer` (`internal/adapter/driven/docker/inspector.go`, Zeilen 75 bis 100) filtert allein mit `docker ps -a --no-trunc --filter label=devcontainer.local_folder=<Ordner>` und nimmt den ersten Treffer. Für (i) müsste die Suche Container mit dem Instanz-Label ausschließen; ob `docker ps` ein fehlendes Label per `--filter` ausdrücken kann, ist nicht gelesen und nicht gemessen (eine Auswertung des Label-Werts über `--format` in `u-boot` ist denkbar, nicht gemessen).
- `runArgs` (`internal/adapter/driven/docker/devcontainer_launcher.go`, Zeilen 284 bis 303) ruft `devcontainer up --workspace-folder <Ordner>` ohne `--id-label` und reicht nur `--remove-existing-container` an die CLI; die Prüfung „nur Instanz-Container vorhanden → Exit-Code `10`“ muss deshalb **vor** diesem Aufruf in `u-boot` liegen. Wo genau, hängt an den Aufrufstellen von `FindContainer` (`internal/hexagon/application/devcontainer_start.go` und `devcontainer_fallback.go`, laut Entwurf des Code-Slice); diese Dateien sind **nicht gelesen**.
- [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) (gelesen: Kontext, Entscheidung, Konsequenzen): Volume-Namen enthalten `-${devcontainerId}`; „eine Instanz ist ein Ordner“; das Schema ist die Grundlage der Begründung oben. Dass die ADR „Instanz“ für „Ordner“ verwendet, ist ein dritter Gebrauch des Wortes; die Folge-ADR des Code-Slice schärft sie (Entwurf des Code-Slice, Bezug).

### Gefundener Stand (gelesen am 2026-10-31)

- [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 508): Überschrift, Einleitungsabsatz (499), acht Spiegelstriche — `up` (501, mit `--recreate`), `down` (502), `shell` (503), Hilfs-Container und Socket (504), Zugangsdaten (505), Fehlerklassen (506), Zusage-Menge (507) und Nicht-Ziele (508). Keiner nennt eine Instanz.
- [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Zeile 473) sagt „Der Daemon-Speicher liegt in einem benannten Volume je Instanz“; gemeint ist der einzelne Container (`${devcontainerId}`), nicht ein vom Nutzer vergebener Name. Unverändert (Entscheidung 2); der Wortlaut unten sagt „benannte Instanz“ und definiert den Begriff in [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).
- [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle) (Zeilen 1147 bis 1149): die Zeilen `devcontainer up`, `devcontainer down`, `devcontainer shell` nennen Zweck und Priorität, keine Flags (nicht geändert).
- [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeile 189) zählt `devcontainer up`, `down` und `shell` unter den ändernden Befehlen unter `--quiet`; das Flag ändert daran nichts (nicht geändert).
- [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Zeilen 227 und 236): `10` ist der fachliche Validierungsfehler, „typische Zuordnung: Struktur-, Namens- und Konfigurationsvalidierung“; ein ungültiger Instanzname fällt darunter. Der Fall „nur Instanz-Container vorhanden“ (Fassung (i)) ist kein Namensfehler, sondern ein Zustand des Projekts; `10` ist hier die Entscheidung des Projektinhabers (2026-10-31), nach dem Wortlaut der Tabelle zu `11` (Umgebungs-/Prüfungsfehler) wäre auch eine andere Zuordnung denkbar — nicht Gegenstand.
- Das `--json`-Schema (`command: devcontainer`, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)) ist **nicht gelesen**: ob die Ausgabe von `up` ein Feld für die Instanz bekommt, ist eine Frage der Spezifikation im Code-Slice; der Wortlaut unten sagt dazu nichts.
- Die Egress-, Ressourcen- und Zweig-Fragen berührt der CR nicht (siehe unten).

### Wortlaut (alt/neu; ENTWURF, nicht freigegeben)

Die Alt-Stellen sind aus `spec/lastenheft.md` wörtlich übernommen; die Zeile steht dabei. Im Lastenheft sind Verweise Anker `(#…)` der Datei, hier Pfade von `open/` aus; der Text ist sonst unverändert. „Neu“ ist der vollständige Spiegelstrich, die geänderten Stellen sind **fett**. Das Lastenheft ist Zielbild des Befehls: Belegt sind die gemessenen Aufbauten und das Verhalten der Dev-Container-CLI und der Engine; dass `u-boot devcontainer up --instance` so arbeitet, pinnt erst der Docker-Test des Code-Slice ([`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md)).

**Stelle 1 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 501 (`up`): `--instance` und die Definition der Instanz.**

Alt (Zeile 501):

> - `u-boot devcontainer up` startet den Devcontainer des Projekts im Basisverzeichnis mit der Dev-Container-CLI, die `u-boot` in einem kurzlebigen Hilfs-Container ausführt. Das Hilfs-Image baut `u-boot` lokal aus gepinnten Quellen; der erste Aufruf braucht Netzzugriff. `--recreate` entfernt einen vorhandenen Container und legt ihn neu an; benannte Volumes bleiben erhalten. Nach dem Start weist der Befehl die Degradations-Fallbacks des Containers aus ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)).

Neu:

> - `u-boot devcontainer up` startet den Devcontainer des Projekts im Basisverzeichnis mit der Dev-Container-CLI, die `u-boot` in einem kurzlebigen Hilfs-Container ausführt. Das Hilfs-Image baut `u-boot` lokal aus gepinnten Quellen; der erste Aufruf braucht Netzzugriff. `--recreate` entfernt einen vorhandenen Container und legt ihn neu an; benannte Volumes bleiben erhalten. **In diesem Abschnitt meint „der Devcontainer des Projekts“ den Container ohne Instanz; eine Instanz ist ein vom Nutzer mit `--instance <name>` benannter weiterer Devcontainer desselben Projekts. Mit `--instance <name>` startet, sucht und ersetzt (mit `--recreate`) der Befehl nur den Container dieser Instanz; mehrere Instanzen desselben Projektordners können gleichzeitig laufen. Ohne `--instance` wirkt der Befehl nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `up` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, statt den Container einer Instanz zu übernehmen.** Nach dem Start weist der Befehl die Degradations-Fallbacks des Containers aus ([`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)).

**Stelle 2 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeilen 502 (`down`) und 503 (`shell`).**

Alt (Zeile 502):

> - `u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes.

Neu:

> - `u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes. **Mit `--instance <name>` stoppt er nur den Container dieser Instanz und lässt die übrigen Container des Projekts laufen. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `down` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`.**

*Beleg (nicht Teil des Vertragstexts):* Gemessen ist das Verhalten der Engine und der CLI (`docker stop` eines Containers lässt die übrigen N − 1 laufen; `devcontainer up --remove-existing-container` mit den Labels einer Instanz entfernt nur deren Container, `devcontainerId` und alle Volumes bleiben; alle sechs Läufe, Protokoll M7). `u-boot devcontainer down --instance` ist nicht gemessen; der Wortlaut steht als Zielbild (Entscheidung 5), der Docker-Test des Code-Slice pinnt ihn.

Alt (Zeile 503, erster Satz; die weiteren Sätze der Zeile bleiben unverändert):

> - `u-boot devcontainer shell` startet den Devcontainer wie `u-boot devcontainer up` ohne `--recreate` und öffnet danach eine interaktive Shell im Container, mit dem Benutzer und dem Arbeitsverzeichnis, die der Start meldet.

Neu (erster Satz):

> - `u-boot devcontainer shell` startet den Devcontainer wie `u-boot devcontainer up` ohne `--recreate` **(mit `--instance <name>` den Container dieser Instanz; ohne `--instance` wirkt er nur auf den Container ohne Instanz und endet bei einem Projekt mit nur Containern von Instanzen mit Exit-Code `10` und einem Hinweis auf `--instance`)** und öffnet danach eine interaktive Shell im Container, mit dem Benutzer und dem Arbeitsverzeichnis, die der Start meldet.

**Stelle 3 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), neuer Spiegelstrich nach Zeile 503: Isolation je Instanz (Zusage nur für das Gemessene).**

Alt: — (neu).

Neu:

> - **Bis zu vier gleichzeitig laufende Instanzen desselben Projekts teilen im Sandbox-Profil ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)) kein Volume: Jede Instanz hat ein eigenes Workspace-Volume und, soweit aktiviert, ein eigenes Volume für den Docker-Speicher (`nestedRuntime: docker`, [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer)) und für die Modelle des Features `ollama`; eine Datei im Workspace-Volume einer Instanz ist in keiner anderen Instanz sichtbar. Gemessen ist das mit zwei, drei und vier gleichzeitig laufenden Instanzen, jede mit höchstens einer CPU und 3 GiB Arbeitsspeicher (Grenze samt Swap), im Sandbox-Profil mit `nestedRuntime: docker` und dem Feature `ollama`, im Zustand direkt nach dem Start, auf einem Linux-Host (`amd64`) mit lokaler Docker-Engine. Nicht zugesagt sind mehr als vier Instanzen, der Speicher von Podman (`nestedRuntime: podman`), Instanzen mit mehr oder ohne Begrenzung von CPU und Arbeitsspeicher, das Verhalten unter Last (etwa bei einem Go-Build oder einem Modell-Pull), mit kaltem Cache oder bei kaltem Bau des Images, Dauer und Platzbedarf des Starts, andere Hosts und Engines als gemessen und die Isolation im Standardprofil: Dort teilen alle Instanzen den Projektordner (Bind-Mount).**

*Belege der Einzelaussagen* (Protokoll, Rohdaten `rohdaten-stufe{2,3,4}-lauf{1,2}.txt`): Volume-Namen 3 × N (6, 9, 12), je Familie paarweise verschieden, in allen sechs Läufen (M2/M3); Inhaltsprobe über alle Paare (Datei im Workspace, Datei im Modell-Volume, Image im verschachtelten Docker): 0 Paare mit sichtbarem Fremdinhalt (M3); Grenze zurückgelesen an den Instanzen der Sandbox-Phase (`HostConfig.Memory=3221225472`, `MemorySwap=3221225472`, `NanoCpus=1000000000`; M5); Grenze in keinem Lauf erreicht, `OOMKilled=false`, rund 55 bis 62 MiB je Instanz direkt nach dem Start (M5); Standardprofil: alle N Instanzen und der Host sehen dieselbe Datei im Projektordner (M3); Podman nicht messbar auf dem Messhost (`kernel.apparmor_restrict_unprivileged_userns = 1`, M0, M3). „Bis zu vier“ ist die höchste vollständig gemessene Zahl; Stufe 4 lief mit gelockerter Vorprüfung auf Entscheidung des Projektinhabers (Closure-Notiz des Mess-Slice §7), ohne Swap-Druck; das gehört in die Closure, nicht in das Lastenheft. Nicht im Wortlaut, weil nicht gemessen oder nur Deutung: die Host-Größen (Speicher, CPU), Dev-Container-CLI- und Docker-Version (Vorschlag des Verifiers: nicht nennen; sie stehen im Protokoll), Zahlen zu Dauer und Platz (Handbuch, nicht Lastenheft), der Mechanismus der Übernahme bei `up` ohne `--instance` (Vermutung), die Ursache des Hängers bei `=`. Die Annahme „1 CPU und 3 GiB je Instanz, nicht für alle zusammen“ ist die Deutung des Planners der Vorgabe des Projektinhabers und nicht zurückgefragt.

**Stelle 4 — [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 506 (Fehlerklassen).**

Alt (Zeile 506):

> - Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält; `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `11` auch, wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`; `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind.

Neu:

> - Fehlerklassen ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält, **wenn der Name der Instanz ungültig ist (ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist 1 bis 32 Zeichen lang; die leere Zeichenkette und ein Name mit `=` sind ungültig) und wenn `up`, `down` oder `shell` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt**; `11`, wenn Docker nicht erreichbar … (der Rest der Zeile bleibt wörtlich wie in Alt).

*Begründung der Namensregel* (Entscheidung 1): Der Name wird Label-Wert am Container und Suchwert im Filter von `docker ps`; ein kleiner, geschlossener Zeichensatz schließt die Zeichen aus, die ein Filter, ein Shell-Aufruf oder ein Name anders lesen; Kleinbuchstaben vermeiden zwei Instanzen, die sich nur in der Groß- und Kleinschreibung unterscheiden; der Rand folgt den Namensregeln von Docker (Annahme, nicht gelesen); `=` ist durch den Zeichensatz ausgeschlossen, die leere Zeichenkette durch die Länge 1 (gemessen: Exit-Code `1` der CLI, `Unmatched argument format`; `=`: `up` hängt). Ein kleinerer Kreis ist leichter zu erweitern als zu verengen.

**Stelle 5 — Nicht-Ziele, [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Zeile 508.**

Alt (Zeile 508):

> - Nicht-Ziele: `u-boot` startet keinen Agenten ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), und der Befehl startet nicht den Compose-Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)).

Neu:

> - Nicht-Ziele: `u-boot` startet keinen Agenten ([`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil)), und der Befehl startet nicht den Compose-Stack ([`LH-FA-UP-001`](../../../../spec/lastenheft.md#lh-fa-up-001--umgebung-starten)). **Zu Instanzen: `u-boot` führt keine Liste der Instanzen, räumt keine auf, wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zu einer Instanz.**

**Stelle 6 — Versionsfeld und Historie (§16).**

Alt: `| Version          | 0.24.0 |` (Zeile 9; Stand 2026-10-31).

Neu: `| Version          | 0.25.0 |` — erst bei Annahme; die Nummer ist die dann nächste freie Minor-Version aus der Historie (letzte Zeile heute `0.24.0`, Zeile 1607).

Neue Historie-Zeile:

> | 0.25.0 | ⟦Datum der Annahme⟧ | [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): benannte Instanz an `u-boot devcontainer up`, `down` und `shell` (`--instance <name>`; mehrere Devcontainer desselben Projekts gleichzeitig; ohne `--instance` wirkt der Befehl nur auf den Container ohne Instanz); Namensregel (Kleinbuchstaben, Ziffern, Bindestriche, 1 bis 32 Zeichen), Exit-Code `10` bei ungültigem Namen und bei `up`, `down` und `shell` ohne `--instance`, wenn es nur Container von Instanzen gibt; Isolation der Volumes für bis zu vier gleichzeitige Instanzen im Sandbox-Profil mit `nestedRuntime: docker` und `ollama`, jede mit höchstens einer CPU und 3 GiB, im Zustand direkt nach dem Start, auf einem Linux-Host (`amd64`) mit lokaler Engine; nicht zugesagt: mehr als vier Instanzen, Podman-Speicher, andere Grenzen, Last, Dauer und Platz, andere Hosts, Isolation im Standardprofil, Zweig je Instanz. | Vereinbarung mit dem Projektinhaber |

*(Der Verweis auf die Entscheidung steht im Lastenheft als Anker `#…` der Datei; das Datum ist die einzige Stelle im Wortlaut, die erst bei der Annahme feststeht.)*

### Nicht-Zusagen aus dem Verification-Report: wo sie im Wortlaut stehen

| Nicht in einen Lastenheft-Wortlaut als Zusage | Stelle |
|---|---|
| mehr als vier Instanzen | Stelle 3 (Nicht zugesagt), Stelle 6 |
| andere Hosts und Engines | Stelle 3, Stelle 6 |
| Isolation durch Podman | Stelle 3 (Nicht zugesagt), Stelle 6 |
| Dauer- und Platzzahlen | Stelle 3 (Nicht zugesagt: Dauer und Platzbedarf), keine Zahl im Wortlaut |
| Verbrauch unter Go-Build, Modell-Pull, kaltem Cache oder kaltem Bau; Last | Stelle 3 (Nicht zugesagt), „im Zustand direkt nach dem Start“ |
| Isolation im Standardprofil | Stelle 3 (Bind-Mount geteilt, nicht zugesagt), Stelle 6 |
| Mechanismus der Übernahme bei `up` ohne Flag (Vermutung) | nicht im Wortlaut; Begleittext zu Stelle 3 und Begründung (i) nennen ihn als nicht gemessen |
| Ursache des Hängers bei `=` | nicht im Wortlaut; Stelle 4 schließt `=` über den Zeichensatz aus, ohne die Ursache zu behaupten |
| dass `u-boot devcontainer up --instance` so arbeitet | Entscheidung 5 (Spec-Dateien sind Zielbild); im Slice-Text, nicht im Vertragstext: Beleg-Absatz zu Stelle 2; der Docker-Test des Code-Slice pinnt den Befehl |
| Messbedingung „1 CPU und 3 GiB je Instanz“ als Versprechen über den Zustand direkt nach dem Start hinaus | Stelle 3: „jede mit höchstens …, im Zustand direkt nach dem Start“, „Nicht zugesagt: andere Grenzen, Last“ |

**Keine Änderung** (gelesen, bewusst nicht angefasst): [`LH-FA-DEV-006`](../../../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil) (auch die Regeln zum Hilfs-Container: die Instanz ändert den Socket-Mount nicht) und [`LH-FA-DEV-007`](../../../../spec/lastenheft.md#lh-fa-dev-007--container-runtime-im-sandbox-devcontainer) (Entscheidung 2: „je Instanz“ bleibt, der Begriff ist dort der einzelne Devcontainer), die Zusage-Menge-Zeile 507 von [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) („Linux-Host `amd64`, lokale Engine“; die neue Zusage bleibt innerhalb), [`LH-FA-CLI-005`](../../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) und [`LH-SA-CLI-002`](../../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), das Glossar (kein Eintrag zu „Instanz“).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Messung** — [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md); der CR trägt nur, was dort gemessen ist.
- **Code, Tests, Spezifikation, Handbuch, `CHANGELOG.md`** — [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md); das Lastenheft ist Zielbild, der Code folgt.
- **Welchen Zweig eine Instanz bearbeitet, Agenten verteilen, Compose-Stack je Instanz, Last bei vielen Instanzen** — im Entwurf des Code-Slice ausgeschlossen (§1 dort); der CR führt sie als Nicht-Ziele (Stelle 5) oder als Nicht-Zusage (Stelle 3).
- **Entscheidung über Release oder Welle** — beim Projektinhaber (Roadmap nicht geändert).
- **Änderung an `spec/lastenheft.md`** — erst nach der Freigabe dieses Wortlauts.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Wortlaut ohne offene Entscheidung:** Der Wortlaut ist dem Projektinhaber Stelle für Stelle vorgelegt (Alt/Neu); die Entscheidungen des Projektinhabers (Abschnitt in §1) sind eingearbeitet; offen sind nur Versionsnummer und Datum der Historie-Zeile; die Freigabe ist in §7 vermerkt.
- [ ] **Lastenheft `0.25.0`** (oder die dann nächste freie Version) geschrieben: Stellen 1 bis 6 in der freigegebenen Fassung, Versionsfeld, Historie-Zeile mit Datum (Verweis „Vereinbarung mit dem Projektinhaber“).
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
| `spec/lastenheft.md` | update (erst nach der Freigabe) | Stellen 1 bis 6; gelesen: Kopf (Zeile 9), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 497 bis 508), [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Zeilen 218 bis 236), Glossar (Zeilen 1542 bis 1556), Fundstellen per `grep` (Zeilen 189, 191, 457, 473, 1147 bis 1149), Historie bis `0.24.0` (Zeile 1607) |
| `docs/plan/planning/open/slice-cr-devcontainer-instanz.md` (diese Datei) | update | Wortlaut-Entwurf, Freigabe und Vorab-Suche in §7; Zustandsmarker auf Ist-Stand |
| `docs/plan/planning/open/slice-v1-devcontainer-instanz-label.md` | update | Namensregel (Entscheidung 1) und Test-Auftrag zu Entscheidung 4 (Exit-Code `10`, Prüfung vor dem Aufruf der CLI); Zustandsmarker „CR freigegeben“ im Start-Trigger, Verweis-Nachzug der Moves |
| `docs/plan/planning/open/README.md` | update | Zeile dieses Slice auf „geliefert“; Verweis-Nachzug der Moves (`make slice-mv`) |
| `docs/reviews/` (Review-Report, Verification Evidence) | neu | Review und Verifikation nach `harness/review.md` und `harness/verification.md` |
| `docs/plan/planning/observations/` | update, neu | Register bei der Closure (§7) |
| `spec/spezifikation.md`, `spec/architecture.md`, `docs/user/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `internal/**` | unverändert | der CR ändert kein Produktverhalten; der Code-Slice zieht nach |

Vorab-Suche (Pflicht vor jeder Vertragsänderung, Muster [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md)): `grep -i` über `spec/lastenheft.md` nach den Begriffen der DoD; ausgeführt vor dem Lastenheft-Commit, Treffer in §7. Vorab (diesen Entwurf): `grep -i instanz` — Treffer nur in den Zeilen 416, 473 und der Historie (0.4.0, 0.5.0, 0.22.0), kein Glossar-Eintrag.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt; [`slice-v1-devcontainer-instanz-messung`](../done/slice-v1-devcontainer-instanz-messung.md) liegt in `done/` (eingetreten); der Projektinhaber hat den Wortlaut freigegeben (Gespräch) und die Entscheidungen der Tabelle getroffen; kein anderer Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Der Wortlaut verlangt mehr als eine Anforderung oder eine Entscheidung, die der Mess-Slice nicht getroffen hat.
- `in-progress` → `open` (blockiert): keine Freigabe des Wortlauts; das Lastenheft bleibt unverändert, kein Code-Slice startet.

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

- **Ein Satz im Lastenheft stützt sich auf einen Messwert, den der Mess-Slice nicht liefert** ([`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md)); vor der Freigabe wird die Liste „nicht gemessen“ Satz für Satz gegen den Wortlaut gelesen (geschärfte Regel aus [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §7; die Tabelle „Nicht-Zusagen“ in §1 ist dieser Abgleich). — **Ausgang:** bei Closure.
- **Eine Accepted-ADR trägt einen Satz, den dieser CR überholt:** [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) (Variante D) sagt „zwei Instanzen aus **demselben** Ordner bleiben eine Instanz“. Die Folge-ADR des Code-Slice löst diesen Satz per `Supersedes (teilweise)` oder `Schärft` ab; [`ADR-0015`](../../adr/0015-sandbox-volumes-pro-instanz.md) wird nicht umgeschrieben (`AGENTS.md` §3.5). Der Wortlaut dieses CR nennt keine ADR. — **Ausgang:** weiter offen → [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) (Folge-ADR dort in §3).
- **Das Standardprofil teilt den Projektordner:** Eine Zusage „Instanzen sind isoliert“ ohne die Einschränkung auf das Sandbox-Profil wäre falsch. — **Ausgang:** bei Closure (Stelle 3 nennt die Grenze).
- **Die Annahme „1 CPU und 3 GiB je Instanz“ ist eine Deutung der Vorgabe** (nicht für alle Instanzen zusammen) und nicht zurückgefragt; sie steht als Messbedingung in Stelle 3. — **Ausgang:** bei Closure (Rückfrage an den Projektinhaber bei der Vorlage).
- **Entscheidung 4 verlangt eine Prüfung in `u-boot` vor dem Aufruf der CLI und eine Suche, die Container mit Instanz-Label ausschließt;** beides ist nicht gelesen oder gemessen (Begründung zu Entscheidung 4). — **Ausgang:** weiter offen → [`slice-v1-devcontainer-instanz-label`](slice-v1-devcontainer-instanz-label.md) (Test-Auftrag dort).
- **Die Version `0.25.0` ist vergeben, wenn der Slice startet:** andere Entwürfe in `open/` mit eigenem Lastenheft-CR können sie früher nehmen. — **Ausgang:** bei Closure (Version bei der Annahme aus der Historie lesen).

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt: Freigabe des Wortlauts (Datum) · Vorab-Suche mit Treffern · Review und Verification · Steering-Loop-Eintrag · Beobachtungs-Register · Folge-Slices · Risiken aus §6 · Grenzen · Drei Paarungen · Abgleich der Plan-Tabelle gegen `git diff --stat`.)*

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

**Vorgelagert — offene Beobachtungen sichten:** Das Register ([`observations/`](../observations/)) ist nach den Titeln durchgesehen. Gelesen: [`BEO-DOK/zusage-weiter-als-messung`](../observations/BEO-DOK/zusage-weiter-als-messung/observation.md) als Verweis aus [`slice-cr-devcontainer-ollama`](../done/slice-cr-devcontainer-ollama.md) §7 (Inhalt hier nicht gelesen; die Regel daraus steht in §6) und [`BEO-APP/container-identitaet-nur-ueber-label`](../observations/BEO-APP/container-identitaet-nur-ueber-label/observation.md) (Kopf und Trigger; der Messbeleg des Mess-Slice liegt als Evidence-Datei dort). [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) gilt (`AGENTS.md` §3.10): Jede Bedingung im Wortlaut („nur im Sandbox-Profil“, „kein geteiltes Volume“, „Exit-Code `10` bei ungültigem Namen“, „Exit-Code `10` bei nur Container von Instanzen“) braucht im Code-Slice einen pinnenden Test oder entfällt.

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
