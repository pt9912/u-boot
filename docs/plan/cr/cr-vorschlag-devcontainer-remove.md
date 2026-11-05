# CR-Vorschlag: `u-boot devcontainer remove` — Lastenheft 0.36.1 → 0.37.0, Alt/Neu je Stelle

**Status:** Entwurf des Planners (2026-11-05); der Wortlaut wird dem Projektinhaber vorgelegt und wartet auf die Freigabe. Das Lastenheft ist nicht geändert.

**Gegenstand:** die Lieferung zur DoD-Zeile 1 des CR-Slices [`slice-cr-devcontainer-remove`](../planning/in-progress/slice-cr-devcontainer-remove.md) (Lastenheft-Wortlaut); Folge-Slices: Spezifikation [`slice-v2-spec-devcontainer-remove`](../planning/open/slice-v2-spec-devcontainer-remove.md), Umsetzung [`slice-v1-devcontainer-remove`](../planning/open/slice-v1-devcontainer-remove.md). Diese Datei ist kein Slice und trägt keinen Lebenszyklus; sie ändert nichts an den Verträgen. Ihr **Ablage-Ort** ist dieser Ordner (`docs/plan/cr/`) je [`harness/conventions.md`](../../../harness/conventions.md#ablage-von-messprotokollen-und-skripten) (Ablage-Tabelle).

**Grundlage:** die Entscheidungen E1 bis E5 und die Antwort zur Frage aus der Messung (M5) des Projektinhabers vom 2026-11-05 (Antwort im Sitzungsdialog; im Slice unter *Entscheidungen* mit Ausgang geführt) und das Messprotokoll [`slice-v1-devcontainer-remove-messung`](../../messungen/slice-v1-devcontainer-remove-messung/README.md) (gelesen 2026-11-05; Verification: [`verification-2026-11-05-slice-v1-devcontainer-remove-messung`](../../reviews/verification-2026-11-05-slice-v1-devcontainer-remove-messung.md)). Das Lastenheft nimmt nur die Anforderung auf: den Befehl, sein Verhalten aus Sicht des Benutzers, die Abfrage und die Anwendung der Exit-Codes ([`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) wird angewandt, nicht geändert). Suchweg, Schlüsselnamen, Volume-Erkennung, Konfliktregeln, Plattform- und Messgrenzen gehören in die Spezifikation und in die Tests der Folge-Slices (Liste unter „Nicht geändert“).

**Die Entscheidungen des Projektinhabers, auf die der Wortlaut zurückgeht (2026-11-05):**

| | Entscheidung | Stelle |
|---|---|---|
| E1 | Name `u-boot devcontainer remove` | (c), (e), (f) |
| E2 | Standard: Container entfernen, Volumes bleiben; `--volumes` entfernt zusätzlich die Volumes dieses Containers, mit Abfrage, `--yes` und Liste wie `u-boot down --volumes` | (c) |
| E3 | laufender Container: stoppen, dann entfernen | (c) |
| E4 | `--volumes` entfernt auch das Volume der Ollama-Modelle; die Abfrage nennt es | (c) |
| E5 / M5 | ohne `--instance` wie `down`: nur der Container ohne Instanz; gibt es nur Container von Instanzen, Exit-Code `10` mit Hinweis auf `--instance` | (c), (d) |

---

## Stelle (a) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Einleitungssatz (Zeile 512)

**Alt:** „*Priorität: V2.* Das Produkt soll den erzeugten Devcontainer eines Projekts starten, darin eine Shell öffnen und ihn stoppen können, ohne dass auf dem Host ein Devcontainer-fähiger Editor, Node.js oder die Dev-Container-CLI installiert ist; auf dem Host genügen Docker (Client und lokale Engine).“

**Neu (Vorschlag):**
> *Priorität: V2.* Das Produkt soll den erzeugten Devcontainer eines Projekts starten, darin eine Shell öffnen, ihn stoppen und seinen Container entfernen können, ohne dass auf dem Host ein Devcontainer-fähiger Editor, Node.js oder die Dev-Container-CLI installiert ist; auf dem Host genügen Docker (Client und lokale Engine).

Die Überschrift („Devcontainer ohne Editor starten“) bleibt, damit der Anker der Anforderung bestehen bleibt.

## Stelle (b) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Spiegelstrich zu `down` (Zeile 515), letzter Satz

**Alt:** „… gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `down` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`.“ (der Spiegelstrich endet hier; sein erster Satz „`u-boot devcontainer down` stoppt den Container; er entfernt weder den Container noch seine Volumes.“ bleibt unverändert)

**Neu (Vorschlag):** der Spiegelstrich bekommt einen letzten Satz:
> … endet `down` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`. Den Container entfernt `u-boot devcontainer remove`.

`down` behält seine Zusage („stoppen, nichts entfernen“); der Satz verweist nur auf den neuen Befehl.

## Stelle (c) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), neuer Spiegelstrich nach dem Spiegelstrich zu `shell` (nach Zeile 516)

**Alt:** Es gibt keinen Spiegelstrich zu einem Befehl, der den Container entfernt; nach dem Spiegelstrich zu `u-boot devcontainer shell` folgt der Spiegelstrich zu den Volumes zweier Instanzen („Zwei gleichzeitig laufende Instanzen …“).

**Neu (Vorschlag):** ein neuer Spiegelstrich:
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (Exit-Code `10` oder abgelehnte Abfrage) bleiben Container und Volumes bestehen. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst aus der Konfiguration.

Zur Lesart („Satz für Satz“, kein verdichteter Vortrag): Satz 1 trägt E3, Satz 2 und 3 E2 und E4, Satz 4 bis 6 die Abfrage von E2 und die Anwendung von [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), Satz 7 und 8 E5, Satz 9 die Auslegung zu F2, Satz 10 grenzt gegen `u-boot remove <service>` ab (Risiko des Slice).

## Stelle (d) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Fehlerklassen (Zeile 521) und Nicht-Ziele (Zeile 523)

**(d1) Alt (Fehlerklassen, Auszug):** „… und wenn `up`, `down` oder `shell` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt; `11`, …“

**(d1) Neu (Vorschlag):**
> … und wenn `up`, `down`, `shell` oder `remove` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt, und wenn `remove --volumes` im nicht-interaktiven Modus ohne `--yes` aufgerufen wird; `11`, …

Der Rest des Spiegelstrichs (Exit-Codes `11` und `12`) bleibt wörtlich; er nennt Docker und die Engine ohne Befehl als Subjekt und gilt damit auch für `remove`.

**(d2) Alt (Nicht-Ziele, Auszug):** „… Zu Instanzen: `u-boot` führt keine Liste der Instanzen, räumt keine auf, wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zu einer Instanz.“

**(d2) Neu (Vorschlag):**
> … Zu Instanzen: `u-boot` führt keine Liste der Instanzen, räumt keine von sich aus auf (eine Instanz entfernt nur der ausdrückliche Aufruf `u-boot devcontainer remove --instance <name>`), wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zu einer Instanz. Volumes eines Containers, der nicht mehr vorhanden ist (etwa nach `u-boot devcontainer remove` ohne `--volumes` oder nach einem Entfernen des Containers von Hand), entfernt `u-boot` nicht.

Grund für (d2): „räumt keine auf“ liest sich neben `remove --instance` als Widerspruch; der zweite Satz benennt die Grenze, die der Benutzer sonst erst beim Ausprobieren findet (siehe F3).

## Stelle (e) — [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Befehlsübersicht (Zeilen 1186 bis 1188)

**Alt:** die Tabelle endet im Block der Devcontainer-Befehle mit
```text
| `u-boot devcontainer up`     | Devcontainer ohne Editor starten   | V2        |
| `u-boot devcontainer down`   | Devcontainer stoppen               | V2        |
| `u-boot devcontainer shell`  | Shell im Devcontainer öffnen       | V2        |
| `u-boot template list`       | Templates anzeigen                 | V1        |
```

**Neu (Vorschlag):** eine Zeile nach `devcontainer shell`:
```text
| `u-boot devcontainer remove` | Devcontainer-Container entfernen   | V2        |
```

## Stelle (f) — Verweis-Stellen auf die Befehle: [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging) (Zeile 189) und [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Zeilen 210 und 214)

**(f1) Alt (CLI-005, Auszug):** „Unter `--quiet` geben die Subkommandos, die etwas ändern (`init`, `generate`, `add`, `remove`, `config set`, `up`, `down`, `devcontainer up`, `devcontainer down`, `devcontainer shell`), nur Fehler aus; …“

**(f1) Neu (Vorschlag):**
> … `devcontainer up`, `devcontainer down`, `devcontainer remove`, `devcontainer shell`), nur Fehler aus; …

**(f2) Alt (CLI-005A, Zeile 210):** „- Die Optionen sind auf Befehle anzuwenden, die Bestätigungsentscheidungen benötigen (insb. `u-boot init`, `u-boot add`, `u-boot remove`, `u-boot config set`, `u-boot down --volumes`).“

**(f2) Neu (Vorschlag):**
> - Die Optionen sind auf Befehle anzuwenden, die Bestätigungsentscheidungen benötigen (insb. `u-boot init`, `u-boot add`, `u-boot remove`, `u-boot config set`, `u-boot down --volumes`, `u-boot devcontainer remove --volumes`).

**(f3) Alt (CLI-005A, Zeile 214, Anfang):** „Bei destruktiven Operationen (insb. `u-boot down --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad …“

**(f3) Neu (Vorschlag):**
> Bei destruktiven Operationen (insb. `u-boot down --volumes`, `u-boot devcontainer remove --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad …

Der Rest von Zeile 214 („Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.“) trägt den Exit-Code `10` in Stelle (c); die Zeile bleibt sonst wörtlich.

## Stelle (g) — Versionsfeld (Zeile 9) und §16 Historie (nach Zeile 1659)

**Alt:** Version `0.36.1`; §16 Historie endet mit der Zeile `0.36.1` (2026-11-05).

**Neu (Vorschlag):** Version `0.37.0` und eine Historie-Zeile mit dem Datum der Übernahme:
> | 0.37.0 | *(Datum der Übernahme)* | [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): neuer Befehl `u-boot devcontainer remove` (entfernt den Container, ein laufender wird zuerst gestoppt; `--volumes` entfernt zusätzlich die Volumes des Containers samt dem Volume der Modelle des Features `ollama`, mit Nennung der Volumes vor der Abfrage, `--yes` und Nennung der entfernten Volumes; `--instance <name>`; ohne `--instance` nur der Container ohne Instanz, Exit-Code `10` bei nur Instanzen; ohne Container Exit-Code `0`); `devcontainer down` und `up --recreate` behalten ihre Zusage; die Nicht-Ziele zu Instanzen und die Fehlerklassen nennen den Befehl. [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): der Befehl in den Aufzählungen und der Befehlsübersicht. | Entscheidung des Projektinhabers 2026-11-05 *(Zusatz bei der Übernahme: „Wortlaut freigegeben wie vorgelegt“ oder die Abweichung je Stelle)* |

Die Historie-Zeile verweist auf nichts unterhalb des Lastenhefts (Decken-Regel). Die Zahl `0.37.0` ist der nächste freie Schritt der Lastenheft-Version, nicht die Release-Version; der Slice nannte bis zu diesem Entwurf noch `0.36.0` (vergeben an die Anforderung zu den Homebrew-Formeln, [`cr-vorschlag-brew-pakete`](cr-vorschlag-brew-pakete.md)) — der Slice ist berichtigt.

## Nicht geändert

- [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): wird angewandt (`0`, `10`), nicht geändert.
- [`LH-FA-UP-004`](../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen): `u-boot down --volumes` deckt den Compose-Stack ab; der Slice grenzt sich ab (Abgrenzung „`remove` für Compose-Dienste“).
- [`LH-FA-CLI-007`](../../../spec/lastenheft.md#lh-fa-cli-007--dry-run): der Wert `devcontainer` des Feldes `command` steht schon in der Aufzählung; ob und wie `remove` eine Vorschau oder eine `--json`-Ausgabe kennt, ist Sache der Spezifikation.
- `devcontainer down`, `up --recreate`, `devcontainer shell` und `u-boot remove <service>`: unverändert.
- Der Spiegelstrich zur Zusage-Menge („Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`) …“, Zeile 522) bleibt wörtlich (siehe F4).
- **Gehört in die Spezifikation und in die Tests der Folge-Slices, nicht ins Lastenheft** (Anweisung des Projektinhabers aus dem Vorgängerfall): die Suche des Containers (Labels, Instanz-Label, die Suche auch nach gestoppten Containern); die Bestimmung der Volumes eines Containers vor dem Entfernen (Messung: nur am vorhandenen Container möglich) und der Umgang mit anonymen Volumes; die Reihenfolge Stoppen, Abfrage, Entfernen; der Fall eines Volumes, das ein anderer Container hält; das Verhalten mit `--json`; die Texte der Ausgabe und der Meldungen; die Zuordnung der Exit-Codes `11` und `12`; Plattform-, Engine- und Messgrenzen; die Kopien der Befehle in `docs/maintainer/cli-json-contract.md` (Abschnitt 6.11 und die Matrix der Dateisystem-Mutationen je Subcommand) — sie ziehen im Spezifikations-Slice nach, wie es [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) verlangt.

## Prüfung des Wortlauts gegen die Messgrenzen

Das Protokoll führt je Frage unter „Darf nicht tragen“ Sätze, die der CR nicht tragen darf. Der Wortlaut enthält keinen davon:

| Messgrenze (Protokoll) | Wie der Wortlaut sie hält |
|---|---|
| M1: `u-boot` ermittelt die Volumes eines schon entfernten Containers nicht | (d2) sagt, dass `u-boot` solche Volumes nicht entfernt; (c) verlangt die Nennung der Volumes nur vor dem Entfernen des Containers |
| M1: Namen aus `${devcontainerId}` berechnen; Labels mit Verweis; anonyme Volumes ausgeschlossen; Runtime „funktioniert“; andere Hosts und Engines | kein Satz dazu: kein Name, kein Label, kein Schlüssel, keine Plattform im Wortlaut |
| M2: `docker volume rm -f` entfernt ein gehaltenes Volume; gestoppter Container gibt sein Volume frei; `docker rm -v` entfernt benannte Volumes | kein Satz zu den Docker-Befehlen; der Fall eines von einem anderen Container gehaltenen Volumes bleibt der Spezifikation (F5) |
| M3: Präfix bestimmt die Volumes einer Instanz; Isolation von Docker-Speicher und Modellen zwischen Instanzen; drei und mehr Instanzen; Podman-Instanzen; Standardprofil mit Katalog-Features | (c) benennt kein Präfix, keine Isolation und kein Profil; „lässt die übrigen Container … und ihre Volumes unberührt“ ist die Anforderung an den Befehl, die der Umsetzungs-Slice mit einem Test pinnt, keine Aussage über die Messung |
| M4: ein Modell übersteht das Entfernen; die Modelle überstehen den Neubau des Images; Volume-Namen stabil bei Umzug; Podman-Speicher; die Entscheidung über `--volumes` (Sache des CR) | (c) sagt nichts zu Modellen oder Neubau; das Volume der Modelle wird als Volume des Features `ollama` genannt (E4), nicht als Inhalt |
| M5: ein `remove` ohne `--instance` trifft dieselbe Auswahl wie `down` | die Auswahl ist die Entscheidung E5/M5 des Projektinhabers (Anforderung), nicht eine Aussage der Messung; `shell`, Standardprofil, eine Instanz, `paused`/`restarting`/`created` stehen nicht im Wortlaut |

Die Bedingungen im Wortlaut („nur auf den Container ohne Instanz“, „im nicht-interaktiven Modus ohne `--yes`“, „ohne Container Exit-Code `0`“) bekommen im Umsetzungs-Slice je einen pinnenden Test oder entfallen ([`AGENTS.md`](../../../AGENTS.md) §3.10, Bedingung in der Spezifikation ohne Test).

---

**Abweichungen der Übernahme je Stelle** *(nur nach der Freigabe tragen — jede Abweichung des Vertrags-Diffs gegenüber dem Vorschlags-Wortlaut je Stelle ausweisen; der §7-Abgleich des CR-Slices trägt dieselbe Liste).*

**Offen zur Freigabe** (je Frage der Vorschlag und die Gegenposition; F1 bis F5 sind Auslegungen, die E1 bis E5 offen ließen, F6 und F7 Form):

- **F1 — Folge einer abgelehnten oder fehlenden Freigabe.** Vorschlag (Stelle (c), Satz 5): bleibt `--volumes` ohne Freigabe, bleiben Container und Volumes bestehen (wie bei `u-boot down --volumes`: Abbruch, nichts geschieht). Gegenposition: Der Container wird auch bei abgelehnter Abfrage entfernt, nur die Volumes bleiben.
- **F2 — Aufruf ohne Container.** Vorschlag (Stelle (c), Satz 9): Exit-Code `0`, nichts wird entfernt (wie `down`, gemessen für `down`; für `remove` eine Anforderung). Gegenposition: Exit-Code `10` mit Hinweis, weil eine zerstörerische Aktion ohne Ziel ein Eingabefehler sein kann. Hintergrund: ein zweiter Aufruf nach einem abgebrochenen oder wiederholten `remove` soll nicht scheitern.
- **F3 — Volumes nach einem `remove` ohne `--volumes`.** Die Messung trägt nur die Bestimmung der Volumes am vorhandenen Container; nach dem Entfernen bleibt nur der Name (M1, M4). Folge für den Benutzer: Wer `remove` ohne `--volumes` aufgerufen hat, entfernt die Volumes danach nicht mehr mit `u-boot` (Stelle (d2), zweiter Satz). Vorschlag: so festschreiben und im Handbuch die Reihenfolge nennen („Volumes zugleich mit dem Container entfernen“). Gegenposition: Satz streichen und die Frage der Spezifikation überlassen; ein namensbasiertes Entfernen verwaister Volumes bräuchte eine neue Messung (das Verfahren der Namensbildung ist nach dem Protokoll weder gelesen noch gemessen).
- **F4 — Zusage-Menge.** Der Spiegelstrich „Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`) …“ nennt `remove` nicht. Vorschlag: unverändert lassen und Plattform- und Messgrenzen von `remove` in der Spezifikation führen (Anweisung des Projektinhabers: keine Plattform-/Messgrenzen-Sätze im Lastenheft). Gegenposition: Ein Halbsatz nimmt `remove` in die Zusage-Menge auf; das Protokoll ist ein Lauf auf der verschachtelten Engine des Entwicklungs-Devcontainers, keine native Engine (Risiko 1 der Messung).
- **F5 — Volume, das ein anderer Container hält.** Vorschlag: nicht im Lastenheft; die Spezifikation legt fest, was `remove --volumes` dann tut (die Messung belegt nur, dass das Entfernen des Volumes scheitert, solange ein Container es hält, M2). Der Wortlaut in Stelle (c) sagt „die Volumes dieses Containers“ ohne Ausnahme. Gegenposition: ein Halbsatz „soweit kein anderer Container sie hält“.
- **F6 — Priorität `V2`** (wie die übrigen Devcontainer-Befehle in [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)). Gegenposition: `V1`.
- **F7 — Version des Lastenhefts.** `0.37.0` (neuer Befehl); die Zahl vergibt die Übernahme. Gegenposition: `0.36.2`, falls der Befehl als Präzisierung gelten soll.

**Regeln dieser Form:** die Kennungen sind verlinkt; der Neu-Wortlaut nennt Anforderungen, keine Plattform-, Bau- oder Messgrenzen, keine Labelnamen, Pfade und Flag-Interna; die Messbasis steht je Link mit Datum und Grenzen; die Historie-Zeile entsteht bei der Übernahme mit dem Datum des Commits und folgt dem Muster der Vorgänger-CRs.
