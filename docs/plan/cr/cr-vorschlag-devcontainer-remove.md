# CR-Vorschlag: `u-boot devcontainer remove` — Lastenheft 0.36.1 → 0.37.0 und Folgeänderungen 0.37.1, 0.37.2 und 0.37.3, Alt/Neu je Stelle

**Status:** freigegeben wie vorgelegt (Antwort des Projektinhabers im Sitzungsdialog, 2026-11-05: „freigegeben wie vorgelegt“); Wortlaut, Lastenheft-Version `0.37.0`, Priorität `V2` und F1 bis F7 wie vorgeschlagen. Übernommen ins Lastenheft `0.37.0`.

**Folgeänderung `0.37.1`** (Review [`2026-11-05-slice-cr-devcontainer-remove`](../../reviews/2026-11-05-slice-cr-devcontainer-remove.md) F-1, F-3, F-4, F-5; Entscheidungen des Projektinhabers vom 2026-11-05, Sitzungsdialog nach dem Review): der Wortlaut steht im Abschnitt „Folgeänderung 0.37.1“ unten (Stellen (h) bis (k)). Der Projektinhaber hat ihn im Sitzungsdialog am 2026-11-05 freigegeben wie vorgelegt; übernommen ist er im Lastenheft `0.37.1` (Commit `a7914c99`), der Nachlauf-Review liegt vor ([`2026-11-05-slice-cr-devcontainer-remove-nachlauf`](../../reviews/2026-11-05-slice-cr-devcontainer-remove-nachlauf.md)).

**Folgeänderung `0.37.2`** (Nachlauf-Review, MEDIUM „Freigabe und Containersuche“ und LOW „Image-Satz“; Entscheidungen des Projektinhabers im Sitzungsdialog nach dem Nachlauf-Review): der Wortlaut steht im Abschnitt „Folgeänderung 0.37.2“ unten (Stellen (l) bis (o)) und in `spec/lastenheft.md` Version `0.37.2`. Der Projektinhaber hat ihn im Sitzungsdialog am 2026-11-06 freigegeben wie vorgelegt, einschließlich der Auslegung von `--instance` ohne Container der Instanz (Exit-Code `0`, keine Abfrage); der Nachlauf-Review über den Diff von `0.37.1` zu `0.37.2` liegt vor ([`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](../../reviews/2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md)); sein MEDIUM behandelt die Folgeänderung `0.37.3`.

**Folgeänderung `0.37.3`** (zweiter Nachlauf-Review, MEDIUM „Ausnahme der Version 0.37.2 widerspricht der Regel in der Anforderung zu Interaktivität und Automatisierung“; Entscheidung des Projektinhabers im Sitzungsdialog, 2026-11-06): der Wortlaut steht im Abschnitt „Folgeänderung 0.37.3“ unten (Stellen (p) und (q)) und in `spec/lastenheft.md` Version `0.37.3`. Der Projektinhaber hat ihn im Sitzungsdialog am 2026-11-06 freigegeben wie vorgelegt; der Nachlauf-Review über den Diff von `0.37.2` zu `0.37.3` steht aus (DoD-Zeile 7 des CR-Slice).

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

## Stelle (c) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), neuer Spiegelstrich nach dem Spiegelstrich zu `shell` (nach Zeile 516; im Ist-Stand Zeile 517, Wortlaut der Übernahme `0.37.0`, geändert in `0.37.1`: Stelle (h))

**Alt:** Es gibt keinen Spiegelstrich zu einem Befehl, der den Container entfernt; nach dem Spiegelstrich zu `u-boot devcontainer shell` folgt der Spiegelstrich zu den Volumes zweier Instanzen („Zwei gleichzeitig laufende Instanzen …“).

**Neu (Vorschlag):** ein neuer Spiegelstrich:
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (Exit-Code `10` oder abgelehnte Abfrage) bleiben Container und Volumes bestehen. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst aus der Konfiguration.

Zur Lesart („Satz für Satz“, kein verdichteter Vortrag): Satz 1 trägt E3, Satz 2 und 3 E2 und E4, Satz 4 bis 6 die Abfrage von E2 und die Anwendung von [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) und [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), Satz 7 und 8 E5, Satz 9 die Auslegung zu F2, Satz 10 grenzt gegen `u-boot remove <service>` ab (Risiko des Slice).

## Stelle (d) — [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten), Fehlerklassen (vor der Übernahme Zeile 521, im Ist-Stand 522; geändert in `0.37.1`: Stelle (i)) und Nicht-Ziele (vor der Übernahme Zeile 523, im Ist-Stand 524)

**(d1) Alt (Fehlerklassen, Auszug):** „… und wenn `up`, `down` oder `shell` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt; `11`, …“

**(d1) Neu (Vorschlag, Wortlaut der Übernahme `0.37.0`; die Folgeänderung `0.37.1` ändert den Spiegelstrich, Stelle (i)):**
> … und wenn `up`, `down`, `shell` oder `remove` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt, und wenn `remove --volumes` im nicht-interaktiven Modus ohne `--yes` aufgerufen wird; `11`, …

Der Rest des Spiegelstrichs (Exit-Codes `11` und `12`) blieb in `0.37.0` wörtlich; er nennt Docker und die Engine ohne Befehl als Subjekt. Das Review F-1 liest ihn als unvollständig (Zuordnung der abgelehnten Abfrage, des Fehlers beim Entfernen und der Klauseln ohne Befehls-Subjekt); `0.37.1` ändert ihn (Stelle (i)).

**(d2) Alt (Nicht-Ziele, Auszug):** „… Zu Instanzen: `u-boot` führt keine Liste der Instanzen, räumt keine auf, wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zu einer Instanz.“

**(d2) Neu (Vorschlag):**
> … Zu Instanzen: `u-boot` führt keine Liste der Instanzen, räumt keine von sich aus auf (eine Instanz entfernt nur der ausdrückliche Aufruf `u-boot devcontainer remove --instance <name>`), wählt keinen Zweig je Instanz und verteilt keine Arbeit auf Instanzen; der Compose-Stack gehört nicht zu einer Instanz. Volumes eines Containers, der nicht mehr vorhanden ist (etwa nach `u-boot devcontainer remove` ohne `--volumes` oder nach einem Entfernen des Containers von Hand), entfernt `u-boot` nicht.

Grund für (d2): „räumt keine auf“ liest sich neben `remove --instance` als Widerspruch; der zweite Satz benennt die Grenze, die der Benutzer sonst erst beim Ausprobieren findet (siehe F3).

## Stelle (e) — [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle), Befehlsübersicht (vor der Übernahme Zeilen 1186 bis 1188, der Devcontainer-Block endete mit `shell` auf 1188; im Ist-Stand die Zeilen 1187 bis 1190 mit `remove` auf 1190)

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

Der Rest von Zeile 214 („Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.“) trägt den Exit-Code `10` in Stelle (c); die Zeile bleibt sonst wörtlich. (In `0.37.3` ergänzt Stelle (p) diesen Satz um die Ausnahme für `remove --volumes` ohne Container.)

## Stelle (g) — Versionsfeld (Zeile 9) und §16 Historie (vor der Übernahme nach Zeile 1659; im Ist-Stand die Zeile `0.37.0` auf 1662 und die Zeile `0.37.1` auf 1663)

**Alt:** Version `0.36.1`; §16 Historie endet mit der Zeile `0.36.1` (2026-11-05). (Stand vor der Übernahme `0.37.0`; die Folgeänderung `0.37.1` ist Stelle (k).)

**Neu (Vorschlag):** Version `0.37.0` und eine Historie-Zeile mit dem Datum der Übernahme (`2026-11-05`):
> | 0.37.0 | 2026-11-05 | [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): neuer Befehl `u-boot devcontainer remove` (entfernt den Container, ein laufender wird zuerst gestoppt; `--volumes` entfernt zusätzlich die Volumes des Containers samt dem Volume der Modelle des Features `ollama`, mit Nennung der Volumes vor der Abfrage, `--yes` und Nennung der entfernten Volumes; `--instance <name>`; ohne `--instance` nur der Container ohne Instanz, Exit-Code `10` bei nur Instanzen; ohne Container Exit-Code `0`); `devcontainer down` und `up --recreate` behalten ihre Zusage; die Nicht-Ziele zu Instanzen und die Fehlerklassen nennen den Befehl. [`LH-FA-CLI-005`](../../../spec/lastenheft.md#lh-fa-cli-005--verbosity-und-logging), [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle): der Befehl in den Aufzählungen und der Befehlsübersicht. | Entscheidung des Projektinhabers 2026-11-05; Wortlaut freigegeben wie vorgelegt |

Die Historie-Zeile verweist auf nichts unterhalb des Lastenhefts (Decken-Regel). Die Zahl `0.37.0` ist der nächste freie Schritt der Lastenheft-Version, nicht die Release-Version; der Slice nannte bis zu diesem Entwurf noch `0.36.0` (vergeben an die Anforderung zu den Homebrew-Formeln, [`cr-vorschlag-brew-pakete`](cr-vorschlag-brew-pakete.md)) — der Slice ist berichtigt.

## Nicht geändert

- [`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes): wird angewandt (`0`, `10`), nicht geändert.
- [`LH-FA-UP-004`](../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen): `u-boot down --volumes` deckt den Compose-Stack ab; der Slice grenzt sich ab (Abgrenzung „`remove` für Compose-Dienste“).
- [`LH-FA-CLI-007`](../../../spec/lastenheft.md#lh-fa-cli-007--dry-run): der Wert `devcontainer` des Feldes `command` steht schon in der Aufzählung; ob und wie `remove` eine Vorschau oder eine `--json`-Ausgabe kennt, ist Sache der Spezifikation.
- `devcontainer down`, `up --recreate`, `devcontainer shell` und `u-boot remove <service>`: unverändert.
- Der Spiegelstrich zur Zusage-Menge („Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`) …“, vor der Übernahme Zeile 522, im Ist-Stand 523) blieb in `0.37.0` wörtlich (siehe F4); `0.37.1` nimmt `remove` auf (Stelle (j)).
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

## Folgeänderung 0.37.1 (nach dem Review): Stellen (h) bis (k), Alt/Neu je Stelle

Anlass: [`2026-11-05-slice-cr-devcontainer-remove`](../../reviews/2026-11-05-slice-cr-devcontainer-remove.md) F-1, F-3, F-4 und F-5 (merge-blockierend: F-1 und F-3). Der Wortlaut änderte sich nach dem Review des CR; der Projektinhaber hat ihn am 2026-11-05 freigegeben wie vorgelegt, der Nachlauf-Review liegt vor ([`BEO-DOK/vertragsaenderung-nach-review`](../planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). Alt ist der Wortlaut von `0.37.0` (`git show 602d8fa4:spec/lastenheft.md`), Neu der Wortlaut von `0.37.1` (`git show a7914c99:spec/lastenheft.md`).

**Entscheidungen des Projektinhabers, auf die die Folgeänderung zurückgeht (2026-11-05, Sitzungsdialog nach dem Review; der Wortlaut ist am 2026-11-05 im Sitzungsdialog freigegeben wie vorgelegt):**

| | Entscheidung | Stelle |
|---|---|---|
| F-3a | Abgelehnte Abfrage: nichts geschieht; die Abfrage kommt vor jedem Eingriff; bei Ablehnung läuft ein laufender Container weiter, nichts wird gestoppt oder entfernt, Exit-Code `10` (wie `u-boot down --volumes`) | (h) |
| F-3b | `remove --instance a` ohne Container der Instanz, auch wenn andere Container des Projekts existieren: Exit-Code `0`, nichts zu entfernen (wie `down`) | (h) |
| F-1 | Abgelehnte Abfrage: Exit-Code `10`; Fehler beim Stoppen oder Entfernen selbst (die Engine verweigert): Exit-Code `12`, der Fehlerklassen-Wortlaut nennt das Entfernen; die Klauseln der gemeinsamen Fehlerklassen-Liste tragen ein Befehls-Subjekt | (i) |
| F-4 | Die Zusage-Menge (Linux-Host mit lokaler Engine, `amd64`) gilt ausdrücklich auch für `remove`; ein Halbsatz; Messgrenzen (verschachtelte Engine) stehen nicht im Lastenheft | (j) |
| F-5 | Abgrenzungssatz zu `u-boot remove <service>` aus [`LH-FA-ADD-007`](../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) gelesen: Dienst des Compose-Stacks gegenüber Container des Devcontainers; die Gemeinsamkeit (Volumes nur auf ausdrückliche Anforderung, `--purge`) genannt | (h) |

### Stelle (h) — Spiegelstrich zu `remove` (Ist-Stand Zeile 517): abgelehnte Abfrage, Instanz ohne Container, Abgrenzung

**Alt (`0.37.0`):**
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (Exit-Code `10` oder abgelehnte Abfrage) bleiben Container und Volumes bestehen. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst aus der Konfiguration.

**Neu (Vorschlag):**
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach, bevor er etwas stoppt oder entfernt; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage) endet der Befehl mit Exit-Code `10` und es geschieht nichts: Ein laufender Container läuft weiter, und kein Container und kein Volume wird gestoppt oder entfernt. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt; gibt es keinen Container dieser Instanz, entfernt er nichts und endet mit Exit-Code `0`, auch wenn es andere Container des Projekts gibt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst des Compose-Stacks aus der Konfiguration, jener entfernt den Container des Devcontainers; beide löschen Volumes nur auf ausdrückliche Anforderung (`--purge` bei `u-boot remove`, `--volumes` bei `u-boot devcontainer remove`).

Lesart der Änderung (Sätze wie in der „Lesart“ zu Stelle (c); die Zahl der Sätze bleibt zehn): Satz 4 sagt, dass die Abfrage vor jedem Eingriff steht („bevor er etwas stoppt oder entfernt“); Satz 5 trägt F-3a und den Teil von F-1(a) zur abgelehnten Abfrage („endet der Befehl mit Exit-Code `10` und es geschieht nichts: Ein laufender Container läuft weiter, und kein Container und kein Volume wird gestoppt oder entfernt“) — der Satz ersetzt „bleiben Container und Volumes bestehen“, das einen laufenden Container offen ließ; Satz 7 trägt F-3b (die Klausel nach dem Semikolon); Satz 10 trägt F-5. Gelesen für Satz 10: [`LH-FA-ADD-007`](../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) (`u-boot remove <service>` nimmt einen Service aus der Konfiguration; Mindestumfang: Eintrag im Compose-Stack, verwaltete Blöcke, `enabled: false`; „Volumes nur auf explizite Anforderung (`--purge`) löschen“, `--purge` ist eine destruktive Operation). Der Satz benennt weder die Datei noch die Schlüssel.

### Stelle (i) — Fehlerklassen (Ist-Stand Zeile 522): Zuordnung von `remove`, Befehls-Subjekte

**Alt (`0.37.0`):**
> - Fehlerklassen ([`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis keine Devcontainer-Konfiguration enthält, wenn der Name der Instanz ungültig ist (ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist 1 bis 32 Zeichen lang; die leere Zeichenkette und ein Name mit `=` sind ungültig), wenn ein Bestandteil des Containernamens die Engine-Regel für `--name` verletzt, wenn ein Container eines anderen Projektordners den Containernamen hält, und wenn `up`, `down`, `shell` oder `remove` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt, und wenn `remove --volumes` im nicht-interaktiven Modus ohne `--yes` aufgerufen wird; `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `11` auch, wenn die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`; `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind.

**Neu (Vorschlag):**
> - Fehlerklassen ([`LH-FA-CLI-006`](../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)): Exit-Code `10`, wenn das Basisverzeichnis für `up` oder `shell` keine Devcontainer-Konfiguration enthält, wenn bei `up`, `down`, `shell` oder `remove` der Name der Instanz ungültig ist (ein Name besteht aus Kleinbuchstaben, Ziffern und Bindestrichen, beginnt und endet mit einem Buchstaben oder einer Ziffer und ist 1 bis 32 Zeichen lang; die leere Zeichenkette und ein Name mit `=` sind ungültig), wenn bei `up` oder `shell` ein Bestandteil des Containernamens die Engine-Regel für `--name` verletzt, wenn bei `up` oder `shell` ein Container eines anderen Projektordners den Containernamen hält, wenn `up`, `down`, `shell` oder `remove` ohne `--instance` auf ein Projekt trifft, in dem es nur Container von Instanzen und keinen Container ohne Instanz gibt, und wenn `remove --volumes` die Freigabe nicht erhält (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage); `11`, wenn Docker nicht erreichbar oder die Engine nicht lokal ist oder wenn bei `up` oder `shell` ein Startscript des Containers wegen eines Umgebungsproblems abbricht (wie bei `doctor`); `11` auch, wenn bei `up` oder `shell` die Statusdatei den Abbruch aus einem früheren Start desselben Containers meldet; ein erneutes `up` ohne `--recreate` endet dann wieder mit `11`; `12` bei Fehlern beim Bau oder Start, die kein Umgebungsproblem sind, und bei `remove`, wenn das Stoppen oder das Entfernen des Containers oder seiner Volumes scheitert, etwa weil die Engine es verweigert.

Lesart: (1) Die Klausel zur Devcontainer-Konfiguration trägt das Subjekt `up` oder `shell`, denn `down` liest keine Projektdatei und braucht keine `devcontainer.json` (gelesen: [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes), „Voraussetzungen von `up` und `shell`“; `internal/hexagon/application/devcontainer_start.go`, `DevcontainerStartService.Down`: `wired`, `validInstance`, `FindContainer`, `Stop`, keine Konfiguration); für `remove` gilt dasselbe wie für `down` (Anforderung, für `remove` ohne Code). (2) Die Namensklauseln des Containers (Engine-Regel für `--name`, Namenskollision) tragen `up` oder `shell`; nach [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (*Namenskollision vor dem Start*) gibt es diese Prüfung bei `down` nicht, weil `down` stoppt und keinen Namen vergibt; `remove` vergibt keinen Namen. (3) Die Klausel zum Namen der Instanz trägt `up`, `down`, `shell` oder `remove`: `--instance <name>` hat bei `remove` dieselbe Namensregel wie bei `down` (Anwendung der Entscheidung E5; für `remove` ohne Code, eine Anforderung). (4) Die zwei Startscript-Klauseln von `11` tragen `up` oder `shell`, denn nur der Start führt ein Startscript aus. (5) Die abgelehnte Abfrage von `remove --volumes` steht mit dem nicht-interaktiven Modus ohne `--yes` unter `10` (F-1(a); bei `u-boot down --volumes` ist die abgelehnte Abfrage `ErrConfirmationRequired` mit Code 10, gelesen: `internal/hexagon/application/downservice.go`). (6) `12` nennt das Stoppen und das Entfernen von Container und Volumes von `remove` (F-1(b)); „etwa weil die Engine es verweigert“ ist ein Beispiel, keine Aussage über die Messung (M2: `docker volume rm` scheitert, solange ein Container das Volume hält). Eine Frage an den Projektinhaber bei der Freigabe: Punkt (2) bis (4) ändern keine Entscheidung von `up`, `down` oder `shell`, sondern benennen ihr Subjekt; liest der Projektinhaber eine der Klauseln weiter, geht sie wieder ohne Subjekt in den Text.

### Stelle (j) — Zusage-Menge (Ist-Stand Zeile 523)

**Alt (`0.37.0`):**
> - Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`); alle übrigen Hosts und Engines, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt. Die Fallback-Ausweisung nach dem Start setzt eine `u-boot.yaml` im Basisverzeichnis voraus. Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus, die der Befehl über die `u-boot.yaml` im Basisverzeichnis findet: ohne diese Datei endet der Befehl bei einem Startscript-Abbruch mit `12`, und ein erneutes `up` meldet einen früheren Abbruch nicht. Im ausgelieferten Container-Image ([`LH-NFA-PORT-004`](../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)) läuft er nicht.

**Neu (Vorschlag):**
> - Der Befehl startet, und `u-boot devcontainer remove` entfernt, unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`); alle übrigen Hosts und Engines, darunter Colima, Docker Desktop, rootless Docker, Podman als Host-Engine, eine entfernte Engine, Linux-`arm64` und Windows, sind nicht zugesagt. Die Fallback-Ausweisung nach dem Start setzt eine `u-boot.yaml` im Basisverzeichnis voraus. Der Exit-Code `11` für den Startscript-Abbruch setzt die Statusdatei des Containers voraus, die der Befehl über die `u-boot.yaml` im Basisverzeichnis findet: ohne diese Datei endet der Befehl bei einem Startscript-Abbruch mit `12`, und ein erneutes `up` meldet einen früheren Abbruch nicht. Im ausgelieferten Container-Image ([`LH-NFA-PORT-004`](../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)) läuft er nicht.

Der Halbsatz „und `u-boot devcontainer remove` entfernt,“ nimmt `remove` in die Zusage-Menge auf (F-4); der Rest des Spiegelstrichs ist wörtlich. Eine Messgrenze (die Messung lief auf der verschachtelten Engine) steht nicht im Lastenheft; sie bleibt Sache der Spezifikation und der Tests. Der Satz „Im ausgelieferten Container-Image … läuft er nicht“ bleibt wörtlich.

### Stelle (k) — Versionsfeld (Zeile 9) und §16 Historie (Ist-Stand Zeile 1663)

**Alt:** Version `0.37.0`; §16 Historie endet mit der Zeile `0.37.0`.

**Neu (Vorschlag):** Version `0.37.1` und eine Historie-Zeile (Vorbild: die Folgeänderung `0.36.1` zu [`LH-FA-DEV-013`](../../../spec/lastenheft.md#lh-fa-dev-013--homebrew-formeln-im-devcontainer) nach dem Review des CR):
> | 0.37.1 | 2026-11-05 | [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Folgeänderung zu `u-boot devcontainer remove` nach dem Review des CR: die Abfrage von `remove --volumes` kommt vor jedem Eingriff, bei fehlender oder abgelehnter Freigabe endet der Befehl mit Exit-Code `10` und es geschieht nichts (ein laufender Container läuft weiter); ohne Container der angegebenen Instanz Exit-Code `0`, auch wenn es andere Container des Projekts gibt; Exit-Code `12`, wenn das Stoppen oder das Entfernen scheitert; die Fehlerklassen nennen für die Klauseln zu Konfiguration, Containername und Startscript den Befehl, für den sie gelten; die Zusage-Menge (Linux-Host `amd64`, lokale Engine) gilt ausdrücklich auch für `remove`; die Abgrenzung zu `u-boot remove <service>` nennt Dienst des Compose-Stacks gegenüber Container des Devcontainers und die Gemeinsamkeit, dass Volumes nur auf ausdrückliche Anforderung gelöscht werden. | Entscheidung des Projektinhabers 2026-11-05 nach dem Review des CR |

Die Historie-Zeile verweist auf nichts unterhalb des Lastenhefts (Decken-Regel).

### Prüfung der Folgeänderung gegen die Messgrenzen

| Satz | Deckung |
|---|---|
| (h) Instanz ohne Container: Exit-Code `0`, auch wenn es andere Container gibt | Anforderung (F-3b); für `down` gemessen: [`slice-v1-devcontainer-remove-messung`](../../messungen/slice-v1-devcontainer-remove-messung/README.md) M3 (`down --instance a` nach dem Entfernen von `a`, `b` und der Container ohne Instanz laufen: Exit-Code 0, „nothing to stop“) und [`SPEC-049`](../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes) (*Mit `--instance`*); für `remove` keine Messaussage |
| (h) abgelehnte Abfrage: Exit-Code `10`, nichts geschieht | Anforderung (F-3a); Muster `u-boot down --volumes` (Code, kein Messwert) |
| (i) `12` bei scheiterndem Stoppen oder Entfernen | Anforderung (F-1(b)); die Messung belegt nur, dass `docker volume rm` scheitert, solange ein Container das Volume hält (M2); der Wortlaut nennt keine Docker-Befehle |
| (j) Zusage-Menge auch für `remove` | Anforderung (F-4); das Protokoll ist ein Lauf auf der verschachtelten Engine (Risiko 1 der Messung): Der Satz sagt zu, was ein Docker-Test auf dem Linux-Host mit lokaler Engine tragen muss ([`slice-v1-devcontainer-remove`](../planning/open/slice-v1-devcontainer-remove.md) §2), keine Messaussage |
| (h) Abgrenzung | gelesen: [`LH-FA-ADD-007`](../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen); keine Messaussage |

Die Bedingungen „gibt es keinen Container dieser Instanz … Exit-Code `0`“, „abgelehnte Abfrage … läuft weiter“ und „`12`, wenn das Stoppen oder das Entfernen scheitert“ bekommen im Umsetzungs-Slice je einen pinnenden Test oder entfallen ([`AGENTS.md`](../../../AGENTS.md) §3.10).

## Folgeänderung 0.37.2 (nach dem Nachlauf-Review): Stellen (l) bis (o), Alt/Neu je Stelle

Anlass: [`2026-11-05-slice-cr-devcontainer-remove-nachlauf`](../../reviews/2026-11-05-slice-cr-devcontainer-remove-nachlauf.md), MEDIUM „Reihenfolge von Freigabe und Containersuche“ und LOW „Reichweite des Image-Satzes“. Alt ist der Wortlaut von `0.37.1` (`git show a7914c99:spec/lastenheft.md`), Neu der Wortlaut von `0.37.2` in `spec/lastenheft.md`. Er ist vom Projektinhaber freigegeben (Sitzungsdialog, 2026-11-06, wie vorgelegt) und geht an einen Nachlauf-Review über den Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). Der Review-Befund zum MEDIUM 1 (Zustandsmarker) ist keine Wortlaut-Änderung und steht hier nicht.

**Entscheidungen des Projektinhabers, auf die die Folgeänderung zurückgeht (Sitzungsdialog nach dem Nachlauf-Review; der Wortlaut ist am 2026-11-06 freigegeben wie vorgelegt):**

| | Entscheidung | Stelle |
|---|---|---|
| G-1 | `remove --volumes` ohne `--yes` im nicht-interaktiven Modus, ohne vorhandenen Container: Exit-Code `0` (nichts zu entfernen). Die Abfrage und der Exit-Code `10` gelten nur, wenn es etwas zu entfernen gäbe. | (l), (m) |
| G-2 | Der Satz zum ausgelieferten Container-Image gilt ausdrücklich auch für `remove`. | (n) |

### Stelle (l) — Spiegelstrich zu `remove` (Zeile 517): Abfrage und Exit-Code `10` nur bei einem Container zu entfernen

**Alt (`0.37.1`):**
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach, bevor er etwas stoppt oder entfernt; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage) endet der Befehl mit Exit-Code `10` und es geschieht nichts: Ein laufender Container läuft weiter, und kein Container und kein Volume wird gestoppt oder entfernt. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt; gibt es keinen Container dieser Instanz, entfernt er nichts und endet mit Exit-Code `0`, auch wenn es andere Container des Projekts gibt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst des Compose-Stacks aus der Konfiguration, jener entfernt den Container des Devcontainers; beide löschen Volumes nur auf ausdrückliche Anforderung (`--purge` bei `u-boot remove`, `--volumes` bei `u-boot devcontainer remove`).

**Neu (Vorschlag):**
> - `u-boot devcontainer remove` entfernt den Container; ein laufender Container wird dafür zuerst gestoppt. Die Volumes des Containers bleiben dabei erhalten. Mit `--volumes` entfernt der Befehl zusätzlich die Volumes dieses Containers, auch das Volume für die Modelle des Features `ollama` ([`LH-FA-DEV-003`](../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)). `--volumes` ist eine destruktive Operation nach [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Der Befehl nennt vor der Abfrage die Volumes, die er entfernen würde, und fragt nach, bevor er etwas stoppt oder entfernt; `--yes` ersetzt die Abfrage; im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`. Ohne Freigabe (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage) endet der Befehl mit Exit-Code `10` und es geschieht nichts: Ein laufender Container läuft weiter, und kein Container und kein Volume wird gestoppt oder entfernt. Nach dem Entfernen nennt der Befehl die entfernten Volumes. Mit `--instance <name>` wirkt er nur auf den Container dieser Instanz und lässt die übrigen Container des Projekts und ihre Volumes unberührt; gibt es keinen Container dieser Instanz, entfernt er nichts und endet mit Exit-Code `0`, auch wenn es andere Container des Projekts gibt. Ohne `--instance` wirkt er nur auf den Container ohne Instanz; gibt es im Projekt nur Container von Instanzen und keinen Container ohne Instanz, endet `remove` ohne `--instance` mit Exit-Code `10` und einem Hinweis auf `--instance`, ohne einen Container zu entfernen. Gibt es im Projekt keinen Container, entfernt der Befehl nichts und endet mit Exit-Code `0`. **Die Abfrage und der Exit-Code `10` ohne Freigabe gelten nur, wenn der Befehl einen Container entfernen würde; gibt es keinen Container zu entfernen (im Projekt oder, mit `--instance`, in der Instanz), fragt der Befehl nicht nach und endet auch mit `--volumes` im nicht-interaktiven Modus ohne `--yes` mit Exit-Code `0`.** `u-boot devcontainer remove` ist nicht `u-boot remove <service>`: Dieser nimmt einen Dienst des Compose-Stacks aus der Konfiguration, jener entfernt den Container des Devcontainers; beide löschen Volumes nur auf ausdrückliche Anforderung (`--purge` bei `u-boot remove`, `--volumes` bei `u-boot devcontainer remove`).

(Der fett gesetzte Satz ist die einzige Änderung der Stelle; im Lastenheft steht er ohne Fettdruck.) Lesart: Der Satz ist ein neuer Satz 10 hinter Satz 9; die Abgrenzung wird Satz 11. Er schränkt Satz 4 („im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`“) und Satz 5 („Ohne Freigabe … `10`“) ausdrücklich auf den Fall ein, dass es etwas zu entfernen gibt, und macht Satz 7 (Instanz ohne Container: `0`) und Satz 9 (kein Container im Projekt: `0`) auch für `--volumes` ohne Freigabe maßgeblich. Der Fall „nur Container von Instanzen, ohne `--instance`“ (Satz 8, Exit-Code `10` mit Hinweis auf `--instance`) bleibt: Der Exit-Code `10` kommt dort aus der fehlenden Auswahl, nicht aus der fehlenden Freigabe, und der Satz sagt „ohne Freigabe“. Der Befehl muss dafür den Container suchen, bevor er die Freigabe prüft; wie, legt die Spezifikation fest.

### Stelle (m) — Fehlerklassen (Zeile 522): `10` bei fehlender Freigabe nur mit einem Container zu entfernen

**Alt (`0.37.1`, Auszug):** „… und wenn `remove --volumes` die Freigabe nicht erhält (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage); `11`, …“

**Neu (Vorschlag, Auszug):**
> … und wenn `remove --volumes` bei einem Container, den der Befehl entfernen würde, die Freigabe nicht erhält (nicht-interaktiver Modus ohne `--yes` oder abgelehnte Abfrage); `11`, …

Der Rest des Spiegelstrichs ist wörtlich. Die Klausel folgt Stelle (l): ohne Container zu entfernen gibt es kein `10` aus fehlender Freigabe.

### Stelle (n) — Zusage-Menge (Zeile 523): der Image-Satz gilt auch für `remove`

**Alt (`0.37.1`, letzter Satz):** „Im ausgelieferten Container-Image ([`LH-NFA-PORT-004`](../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)) läuft er nicht.“

**Neu (Vorschlag):**
> Im ausgelieferten Container-Image ([`LH-NFA-PORT-004`](../../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image)) läuft er nicht, und `u-boot devcontainer remove` läuft dort ebenfalls nicht.

Der Rest des Spiegelstrichs ist wörtlich. „er“ bleibt der Befehl, der startet (Satz 1 des Spiegelstrichs); der Halbsatz nennt `remove` ausdrücklich, damit die Zusage-Menge für beide Befehle dieselbe Grenze trägt.

### Stelle (o) — Versionsfeld (Zeile 9) und §16 Historie (Zeile 1664)

**Alt:** Version `0.37.1`; §16 Historie endet mit der Zeile `0.37.1`.

**Neu (Vorschlag):** Version `0.37.2` und eine Historie-Zeile:
> | 0.37.2 | 2026-11-06 | [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten): Folgeänderung zu `u-boot devcontainer remove` nach dem Nachlauf-Review: die Abfrage von `remove --volumes` und der Exit-Code `10` ohne Freigabe gelten nur, wenn es einen Container zu entfernen gibt; ohne Container (im Projekt oder in der angegebenen Instanz) fragt der Befehl nicht nach und endet auch ohne `--yes` im nicht-interaktiven Modus mit Exit-Code `0`; der Satz zum ausgelieferten Container-Image gilt ausdrücklich auch für `remove`. | Entscheidung des Projektinhabers nach dem Nachlauf-Review des CR |

Die Historie-Zeile verweist auf nichts unterhalb des Lastenhefts (Decken-Regel). Das Datum ist das der Änderung im Lastenheft; die Entscheidungen des Projektinhabers stammen aus dem Sitzungsdialog nach dem Nachlauf-Review.

### Prüfung der Folgeänderung 0.37.2 gegen die Messgrenzen

| Satz | Deckung |
|---|---|
| (l), (m) ohne Container zu entfernen: keine Abfrage, Exit-Code `0` auch mit `--volumes` ohne `--yes` | Anforderung (G-1); für `remove` keine Messaussage. Für `down` gemessen ist nur der Exit-Code `0` ohne Container ([`slice-v1-devcontainer-remove-messung`](../../messungen/slice-v1-devcontainer-remove-messung/README.md) M3); die Reihenfolge von Freigabe und Containersuche ist Code des Umsetzungs-Slice. Gelesen: `u-boot down --volumes` (`internal/hexagon/application/downservice.go`, `DownService.Down`) prüft die Freigabe über `runConfirmationGate` vor dem ersten Engine-Zugriff (`ComposeDown`); `remove` übernimmt diese Reihenfolge nicht, weil der Wortlaut die Suche des Containers vor der Freigabe verlangt |
| (n) Image-Satz auch für `remove` | Anforderung (G-2); keine Messaussage; ob `remove` im ausgelieferten Image tatsächlich nicht läuft, ist nicht gelesen und nicht gemessen |

Die Bedingungen „ohne Container zu entfernen … keine Abfrage, Exit-Code `0`“ und „Im ausgelieferten Container-Image … läuft `remove` nicht“ bekommen im Umsetzungs-Slice ([`slice-v1-devcontainer-remove`](../planning/open/slice-v1-devcontainer-remove.md) §2) je einen pinnenden Test oder entfallen bzw. stehen in der Spezifikation als Nicht-Zusage ([`AGENTS.md`](../../../AGENTS.md) §3.10). Für den Image-Satz zu `up` ist kein pinnender Test gefunden (Suche nach `PORT-004` in `*.go`, `Makefile`, `*.mk`, `*.sh` ohne Treffer; Suche nach `Runtime-Image` in `internal/e2e/*.go` ohne Treffer); die Entscheidung dazu trifft der Spezifikations-Slice.

## Folgeänderung 0.37.3 (nach dem zweiten Nachlauf-Review): Stellen (p) und (q), Alt/Neu je Stelle

Anlass: [`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](../../reviews/2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md), MEDIUM „die Ausnahme der Version 0.37.2 widerspricht der Regel in der Anforderung zu Interaktivität und Automatisierung, die den Befehl namentlich nennt“. Alt ist der Wortlaut von `0.37.2` (`git show 7f568768:spec/lastenheft.md`), Neu der Wortlaut von `0.37.3` in `spec/lastenheft.md`. Er ist vom Projektinhaber freigegeben (Sitzungsdialog, 2026-11-06, wie vorgelegt) und geht an einen Nachlauf-Review über den Diff ([`BEO-DOK/vertragsaenderung-nach-review`](../planning/observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)).

**Entscheidung des Projektinhabers, auf die die Folgeänderung zurückgeht (Sitzungsdialog, 2026-11-06):**

| | Entscheidung | Stelle |
|---|---|---|
| H-1 | Der Widerspruch zwischen [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Zeile 214: nicht-interaktiv ohne `--yes` bei destruktiven Operationen immer Exit-Code `10`, genannt u. a. `u-boot devcontainer remove --volumes`) und [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (ohne Container Exit-Code `0`, keine Abfrage) wird durch eine Ausnahme in [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) aufgelöst: Für `u-boot devcontainer remove --volumes` gilt der Exit-Code `10` nur, wenn es einen Container zu entfernen gibt. `u-boot down --volumes` und `u-boot remove --purge` bleiben unverändert. | (p), (q) |

### Stelle (p) — [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), Absatz zu destruktiven Operationen (Zeile 214, letzter Satz): Ausnahme für `remove --volumes` ohne Container

**Alt (`0.37.2`, ganzer Absatz):**
> Bei destruktiven Operationen (insb. `u-boot down --volumes`, `u-boot devcontainer remove --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad (`--yes` oder aktiv bestätigten interaktiven Pfad) erfolgen. Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.

**Neu (Vorschlag):**
> Bei destruktiven Operationen (insb. `u-boot down --volumes`, `u-boot devcontainer remove --volumes` und `u-boot remove --purge`) darf eine Löschung nur über den expliziten Freigabepfad (`--yes` oder aktiv bestätigten interaktiven Pfad) erfolgen. Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen**; bei `u-boot devcontainer remove --volumes` gilt das nur, wenn es einen Container zu entfernen gibt ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).**

(Der fett gesetzte Teil ist die einzige Änderung der Stelle: Das Semikolon ersetzt den Punkt am Satzende, der Halbsatz kommt hinzu; im Lastenheft steht er ohne Fettdruck, der Link in der Form `#lh-…` ohne Dateipfad.) Lesart: Die Ausnahme steht an der Regel, die den Befehl namentlich nennt, und verweist auf [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) für den Fall ohne Container (Exit-Code `0`, keine Abfrage). Der erste Satz des Absatzes (Löschung nur über den Freigabepfad) bleibt unbedingt: Ohne Container gibt es keine Löschung. `u-boot down --volumes` und `u-boot remove --purge` fallen nicht unter den Halbsatz.

### Stelle (q) — Versionsfeld (Zeile 9) und §16 Historie (Zeile 1665)

**Alt:** Version `0.37.2`; §16 Historie endet mit der Zeile `0.37.2`.

**Neu (Vorschlag):** Version `0.37.3` und eine Historie-Zeile:
> | 0.37.3 | 2026-11-06 | [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung): Folgeänderung zu `u-boot devcontainer remove` nach dem zweiten Nachlauf-Review: der Exit-Code `10` im nicht-interaktiven Modus ohne `--yes` gilt bei `u-boot devcontainer remove --volumes` nur, wenn es einen Container zu entfernen gibt ([`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)); `u-boot down --volumes` und `u-boot remove --purge` bleiben unverändert. | Entscheidung des Projektinhabers nach dem zweiten Nachlauf-Review des CR |

Die Historie-Zeile verweist auf nichts unterhalb des Lastenhefts (Decken-Regel); im Lastenheft stehen die Links in der Form `#lh-…`.

### Prüfung weiterer Stellen im Lastenheft auf denselben unbedingten Satz für `remove --volumes` (gelesen in `spec/lastenheft.md`, Stand `0.37.2`)

| Stelle | Aussage | Ergebnis |
|---|---|---|
| [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeile 214 | „Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen“, `remove --volumes` namentlich genannt | geändert, Stelle (p) |
| [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeile 210 | „Die Optionen sind auf Befehle anzuwenden, die Bestätigungsentscheidungen benötigen (insb. … `u-boot devcontainer remove --volumes`)“ | unverändert: nennt keinen Exit-Code und keine Bedingung; ohne Container braucht der Befehl keine Entscheidung, die Optionen sind „anzuwenden“, wo eine gebraucht wird |
| [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeile 212 | „`--no-interactive` bricht bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab“ | unverändert: ohne Container gibt es keine offene Bestätigungsfrage; die vorbestehende Spannung zwischen Exit-Code `2` (Zeile 212) und `10` (Zeile 214) bleibt, sie gilt gleich für `u-boot down --volumes` und `u-boot remove --purge` (INFO des zweiten Nachlauf-Reviews; offen, §6 des CR-Slice) |
| [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile 517, Satz 4 („im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`“) und Satz 5 („Ohne Freigabe … `10`“) | unbedingte Fassung im Spiegelstrich zu `remove` | unverändert: der Satz 10 desselben Spiegelstrichs schränkt beide auf den Fall ein, dass es einen Container zu entfernen gibt (Lesart unter Stelle (l)); die Ausnahme in Zeile 214 sagt dasselbe |
| [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) Zeile 522 (Fehlerklassen) | `10` bei `remove --volumes` nur „bei einem Container, den der Befehl entfernen würde“ | schon bedingt (Stelle (m)) |
| [`LH-FA-ADD-007`](../../../spec/lastenheft.md#lh-fa-add-007--service-entfernen) Zeile 668 | `--purge` von `u-boot remove`: „Im nicht-interaktiven Modus ohne `--yes` muss der Aufruf mit Exit-Code `10` abgebrochen werden“ | unverändert: anderer Befehl, nicht von der Entscheidung H-1 erfasst |
| [`LH-FA-UP-004`](../../../spec/lastenheft.md#lh-fa-up-004--umgebung-stoppen) Zeile 710 | `u-boot down --volumes` | unverändert: nicht von der Entscheidung H-1 erfasst |
| Zeilen 189 und 1190 | `devcontainer remove` in der `--quiet`-Liste und in der Befehlsübersicht | keine Aussage über Exit-Codes |

### Prüfung der Folgeänderung 0.37.3 gegen die Messgrenzen

| Satz | Deckung |
|---|---|
| (p) `remove --volumes` ohne Container zu entfernen: kein Exit-Code `10` aus [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) | Anforderung (H-1); keine Messaussage und keine neue Zusage gegenüber `0.37.2` (Satz 10 und Fehlerklassen sagen es schon); sie beseitigt den Widerspruch im Rang-1-Dokument. Der Umsetzungs-Slice pinnt den Fall (d) („ohne Container, ohne `--yes`, nicht-interaktiv: `0`“) weiter, jetzt ohne Gegenstelle in Zeile 214 |
| (p) `u-boot down --volumes` und `u-boot remove --purge` unverändert | Der Halbsatz nennt nur `u-boot devcontainer remove --volumes`; gelesen: `DownService.Down` (`internal/hexagon/application/downservice.go`) prüft die Freigabe vor dem ersten Engine-Zugriff, also auch ohne Container mit Exit-Code `10`; für `u-boot remove --purge` ist die Reihenfolge nicht Gegenstand dieser Folgeänderung |

Die Bedingung „bei `u-boot devcontainer remove --volumes` nur, wenn es einen Container zu entfernen gibt“ bekommt im Umsetzungs-Slice ([`slice-v1-devcontainer-remove`](../planning/open/slice-v1-devcontainer-remove.md) §2, Fall (d)) einen pinnenden Test ([`AGENTS.md`](../../../AGENTS.md) §3.10); sie ist dieselbe Bedingung wie in Satz 10 von [`LH-FA-DEV-010`](../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten).

---

**Abweichungen der Übernahme je Stelle:** keine. Der Wortlaut der Stellen (a) bis (f) ist wortgleich übernommen; die Kennungen der Anforderungen sind in der Form des Lastenhefts verlinkt (Anker ohne Dateipfad, `#lh-…`, statt `../../../spec/lastenheft.md#lh-…`). In Stelle (g) sind die Platzhalter des Vorschlags aufgelöst (Datum der Übernahme `2026-11-05`, Quelle-Zelle „Entscheidung des Projektinhabers 2026-11-05; Wortlaut freigegeben wie vorgelegt“), der Vorschlags-Text dort trägt die aufgelösten Werte; die Historie-Zeile trägt keinen Verweis auf Dateien unterhalb des Lastenhefts.

**Abweichung der Folgeänderung `0.37.1` von der Übernahme `0.37.0`:** nach dem Review ändern die Stellen (h) bis (k) den übernommenen Wortlaut (Version `0.37.1`, Vorbild: `0.36.1` im CR [`cr-vorschlag-brew-pakete`](cr-vorschlag-brew-pakete.md)): abgelehnte Abfrage mit Exit-Code `10` und ohne Eingriff, Instanz ohne Container mit Exit-Code `0`, Exit-Code `12` für das Entfernen, Befehls-Subjekte in der Fehlerklassen-Liste, `remove` in der Zusage-Menge, Abgrenzung zu `u-boot remove <service>`. Der Wortlaut der Stellen (a) bis (g) bleibt, soweit die Stellen (h) bis (k) ihn nicht ändern, wie übernommen. Der Wortlaut ist freigegeben wie vorgelegt (Sitzungsdialog, 2026-11-05); die Folgeänderung hat einen Nachlauf-Review (siehe oben), dessen Befunde die Folgeänderung `0.37.2` behandelt.

**Abweichung der Folgeänderung `0.37.2` von `0.37.1`:** nach dem Nachlauf-Review ändern die Stellen (l) bis (o) den Wortlaut von `0.37.1` (Version `0.37.2`): ein Satz im Spiegelstrich zu `remove` und eine Klausel in den Fehlerklassen nehmen den Fall ohne Container zu entfernen von der Abfrage und vom Exit-Code `10` aus (Exit-Code `0`), und der Image-Satz der Zusage-Menge nennt `remove`. Der Wortlaut der Stellen (a) bis (k) bleibt, soweit die Stellen (l) bis (o) ihn nicht ändern, wie in `0.37.1`. Der Wortlaut ist freigegeben wie vorgelegt (Sitzungsdialog, 2026-11-06); die Folgeänderung hat einen Nachlauf-Review ([`2026-11-06-slice-cr-devcontainer-remove-nachlauf-2`](../../reviews/2026-11-06-slice-cr-devcontainer-remove-nachlauf-2.md)), dessen MEDIUM die Folgeänderung `0.37.3` behandelt.

**Abweichung der Folgeänderung `0.37.3` von `0.37.2`:** nach dem zweiten Nachlauf-Review ändern die Stellen (p) und (q) den Wortlaut von `0.37.2` (Version `0.37.3`): ein Halbsatz in [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Zeile 214) nimmt `u-boot devcontainer remove --volumes` ohne Container zu entfernen von der Regel „Exit-Code `10`“ aus, und das Versionsfeld sowie die Historie folgen. Der Wortlaut der Stellen (a) bis (o) bleibt wie in `0.37.2`. Der Wortlaut ist freigegeben wie vorgelegt (Sitzungsdialog, 2026-11-06); die Folgeänderung bekommt einen Nachlauf-Review, der aussteht.

**Bei der Freigabe entschieden** (je Frage der Vorschlag und die verworfene Gegenposition; F1 bis F5 sind Auslegungen, die E1 bis E5 offen ließen, F6 und F7 Form):

- **F1 — Folge einer abgelehnten oder fehlenden Freigabe.** Vorschlag (Stelle (c), Satz 5): bleibt `--volumes` ohne Freigabe, bleiben Container und Volumes bestehen (wie bei `u-boot down --volumes`: Abbruch, nichts geschieht). Gegenposition: Der Container wird auch bei abgelehnter Abfrage entfernt, nur die Volumes bleiben. — **Ausgang:** Vorschlag freigegeben; in `0.37.1` geschärft (F-3a: die Abfrage kommt vor jedem Eingriff, ein laufender Container läuft weiter, Exit-Code `10`; Stelle (h)).
- **F2 — Aufruf ohne Container.** Vorschlag (Stelle (c), Satz 9): Exit-Code `0`, nichts wird entfernt (wie `down`, gemessen für `down`; für `remove` eine Anforderung). Gegenposition: Exit-Code `10` mit Hinweis, weil eine zerstörerische Aktion ohne Ziel ein Eingabefehler sein kann. Hintergrund: ein zweiter Aufruf nach einem abgebrochenen oder wiederholten `remove` soll nicht scheitern. — **Ausgang:** Vorschlag freigegeben; in `0.37.1` auf `--instance <name>` ohne Container dieser Instanz bei vorhandenen anderen Containern ausgedehnt (F-3b, wie `down`; Stelle (h)).
- **F3 — Volumes nach einem `remove` ohne `--volumes`.** Die Messung trägt nur die Bestimmung der Volumes am vorhandenen Container; nach dem Entfernen bleibt nur der Name (M1, M4). Folge für den Benutzer: Wer `remove` ohne `--volumes` aufgerufen hat, entfernt die Volumes danach nicht mehr mit `u-boot` (Stelle (d2), zweiter Satz). Vorschlag: so festschreiben und im Handbuch die Reihenfolge nennen („Volumes zugleich mit dem Container entfernen“). Gegenposition: Satz streichen und die Frage der Spezifikation überlassen; ein namensbasiertes Entfernen verwaister Volumes bräuchte eine neue Messung (das Verfahren der Namensbildung ist nach dem Protokoll weder gelesen noch gemessen). — **Ausgang:** Vorschlag freigegeben.
- **F4 — Zusage-Menge.** Der Spiegelstrich „Der Befehl startet unter Docker auf einem Linux-Host mit lokaler Engine (`amd64`) …“ nennt `remove` nicht. Vorschlag: unverändert lassen und Plattform- und Messgrenzen von `remove` in der Spezifikation führen (Anweisung des Projektinhabers: keine Plattform-/Messgrenzen-Sätze im Lastenheft). Gegenposition: Ein Halbsatz nimmt `remove` in die Zusage-Menge auf; das Protokoll ist ein Lauf auf der verschachtelten Engine des Entwicklungs-Devcontainers, keine native Engine (Risiko 1 der Messung). — **Ausgang:** Vorschlag freigegeben; in `0.37.1` ersetzt durch die Gegenposition (F-4: die Zusage-Menge gilt ausdrücklich auch für `remove`, Entscheidung des Projektinhabers 2026-11-05; Messgrenzen wie die verschachtelte Engine stehen weiter nicht im Lastenheft; Stelle (j)).
- **F5 — Volume, das ein anderer Container hält.** Vorschlag: nicht im Lastenheft; die Spezifikation legt fest, was `remove --volumes` dann tut (die Messung belegt nur, dass das Entfernen des Volumes scheitert, solange ein Container es hält, M2). Der Wortlaut in Stelle (c) sagt „die Volumes dieses Containers“ ohne Ausnahme. Gegenposition: ein Halbsatz „soweit kein anderer Container sie hält“. — **Ausgang:** Vorschlag freigegeben.
- **F6 — Priorität `V2`** (wie die übrigen Devcontainer-Befehle in [`LH-SA-CLI-002`](../../../spec/lastenheft.md#lh-sa-cli-002--vorgesehene-befehle)). Gegenposition: `V1`. — **Ausgang:** `V2`.
- **F7 — Version des Lastenhefts.** `0.37.0` (neuer Befehl); die Zahl vergibt die Übernahme. Gegenposition: `0.36.2`, falls der Befehl als Präzisierung gelten soll. — **Ausgang:** `0.37.0`.

**Regeln dieser Form:** die Kennungen sind verlinkt; der Neu-Wortlaut nennt Anforderungen, keine Plattform-, Bau- oder Messgrenzen, keine Labelnamen, Pfade und Flag-Interna; die Messbasis steht je Link mit Datum und Grenzen; die Historie-Zeile entsteht bei der Übernahme mit dem Datum des Commits und folgt dem Muster der Vorgänger-CRs.
