# Verification-Report: slice-v1-ollama-cloud-messung-mit-konto

* Rolle: Verifier (`harness/verification.md`), frischer Kontext
* Gegenstand: `git diff origin/main..HEAD` (HEAD `8283435a`, zehn Commits), Slice
  [`slice-v1-ollama-cloud-messung-mit-konto`](../plan/planning/in-progress/slice-v1-ollama-cloud-messung-mit-konto.md)
  (Welle [`welle-ollama-cloud-modelle`](../plan/planning/welle-ollama-cloud-modelle.md))
* Bezug: [`LH-FA-DEV-003`](../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features), [`LH-FA-DEV-006`](../../spec/lastenheft.md#lh-fa-dev-006--sandbox-profil), [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer)
* Review: [`2026-10-24-slice-v1-ollama-cloud-messung-mit-konto.md`](2026-10-24-slice-v1-ollama-cloud-messung-mit-konto.md) (F-1 bis F-10 in `5cb018e0` behoben)
* Protokoll: [`README.md`](../messungen/slice-v1-ollama-cloud-messung-mit-konto/README.md)
* Nicht wiederholt: keine Messung, kein Prompt, keine Schlüsseldatei gelesen, nichts gepusht, nichts am Konto geändert.

## Urteil

**Messung verifiziert, Closure noch nicht schließbar.** Sensoren grün, Zugangsregeln im Repo eingehalten, Zählung 40/40 belegt. Drei Bedingungen vor `done/`:

* **V-1 (MEDIUM)** Das Protokoll behauptet für Weg (b) mehr, als gemessen ist (Geräteschlüssel-Zuordnung, Unabhängigkeit des Weges, siehe unten): als Einschränkung und offener Punkt auszuweisen.
* **V-2 (LOW)** DoD-Wortlaut „zwei Läufe“ trägt ein `[x]`, das Protokoll belegt für Weg (b) und Teile von Weg (a) einen Lauf: die Abweichung ist als Planänderung im DoD zu kennzeichnen (F-11, siehe unten).
* Offene Closure-Pflichten, die nicht mir gehören: Häkchen Review/Doku/Gates, Closure-Notiz, Register, Risiken-Ausgänge, Paarungen, Bestätigung des Rückbaus im Konto.

## Sensoren (selbst gefahren)

| Sensor | Ergebnis |
|---|---|
| `make gates` | Exit 0 (`ablage-pruefung-test: 26 Faelle, 0 Fehlschlaege`); Baum danach unverändert |
| `make doc-commits RANGE="origin/main..HEAD"` | Exit 0, 10 Commits, 0 Befunde (`d-check: 825 Datei(en) geprüft, 0 Befund(e)`) |
| `09-suche.sh` (Vorgabe: Messverzeichnis + `tools/messung/ollama-cloud-konto/`) | Exit 0, keine Treffer |
| `09-suche.sh --alle` | Exit 0, 37 Treffer nur in benannten Pfaden der älteren Protokolle |
| eigene Suche (`ssh-ed25519`, `BEGIN … PRIVATE`, `connect?x=`, `ghp_`, `sk-ant`, `glpat-`, Mailadressen des Projektinhabers, Host-Pfade `/home/db`) über Messverzeichnis, Skripte, Slice, Review | Treffer nur in Mustern von `lib.sh`, `09-suche.sh`, `redact-probe.sh` und als Wortnennung im Review-Report; Autor-Mail der Commits ist die `noreply`-Adresse des Repos |
| `make test-docker`, `make verify-depguard` | nicht gelaufen: kein Produktcode, keine Importregel, kein Docker-Pfad des Produkts im Diff (`internal/`, `cmd/`, `spec/`, `docs/user/` nicht im Diff) |
| Messung wiederholt | nein (Kontingent 40/40; Auftrag) |

## DoD, Punkt für Punkt

| # | Zusage | Befund | Beleg |
|---|---|---|---|
| 1 | Sicherung vor dem ersten Lauf (i) bis (v) | **erfüllt, mit ausgewiesener Einschränkung** | (i)/(ii) Artefakt in Commit `5784e921`/`6bfc17e8` vor dem ersten Prompt (11:07:26 UTC); die Datei im Baum ist der zweite Vorlauf, das Protokoll sagt es (Zeile 40). Gegenproben (4) bis (8) entstanden nach dem ersten Prompt, der Abbruch bei 40 wurde nie rot gesehen: im Protokoll ehrlich ausgewiesen (Abschnitt *Gegenproben*, Zeitbezug). (iii) Suche leer, Ausnahmeliste benannt. (iv) zwei Stufen (`10-uebernehmen.sh --gelesen`). (v) kein Antwort-Inhalt, keine Anmelde-URL mit Query (von mir gesucht). |
| 2 | Weg (a): K1 bis K3, zwei Läufe, Engine-Abhängigkeit, Handvarianten, bedingungsgleich | **erfüllt mit Abweichung** | K1 3 Läufe bzw. 2/2, `std`/`discover` 2/2; `restricted`/`allow`-Prompt und K3c-Prompt nur **1 Lauf** (K1 dort 2/2 ohne Prompt). Handvarianten-Tabelle vollständig; Formel „mit dem Konto des Projektinhabers … Docker 29.9.2 …“ im Kopf. Abweichung von „zwei Läufe“: V-2. |
| 3 | Weg (b): K4 bis K8, Nutzerpfad, Liste „nicht gemessen“, Rückabwicklung ausgewiesen | **erfüllt mit Abweichung und Einschränkung** | Ein Lauf (Protokoll „1“). K4 bis K8 mit Rohdaten `rohdaten-06-*`. Rückabwicklung steht, der letzte Schritt (Schlüssel im Konto) ist offen, siehe Abschnitt *Weg (b)*. Einschränkung V-1. |
| 4 | `make gates`, `docs-check`, `ablage-pruefung`, `doc-commits` | **erfüllt, von mir gefahren** | Tabelle Sensoren |
| 5 bis 11 | Review-Häkchen, Doku, Closure-Notiz, Register, Risiken-Ausgänge, Paarungen | **offen, nicht meine Pflicht** | Review-Report liegt vor; §7 trägt den Platzhalter „Closure-Notiz folgt“ |

### F-11: „zwei Läufe“ gegen einen Lauf

Befund: Die DoD-Zeilen 2 und 3 nennen unverändert „zwei, mit frischem Volume und neuem Container“ und stehen auf `[x]`; das Protokoll belegt für Weg (b) und für die Prompt-Fälle in `restricted`/`allow`/`c` einen Lauf. Das ist **keine stille** Abschwächung (Protokoll: Spalte *Läufe*, Tabelle *Nicht gemessen* „ein zweiter Lauf“), aber der DoD-Wortlaut ist nicht angepasst.

Bewertung: **als Planänderung gelten lassen, wenn sie gekennzeichnet wird.** Gründe: (1) der Plan benennt den Fall vorab: §4 Rückführung `in-progress` → `open` (c) „die Obergrenze der Prompts ist vor Abschluss eines Weges erreicht (der Rest steht als ‚nicht gemessen‘)“; (2) die Obergrenze 40 ist Entscheidung des Projektinhabers (Frage 1), 40 von 40 ist belegt (15 `prompts-gestellt`-Zeilen in den Rohdaten, letzte `40 obergrenze=40`), und eine Wiederholung ist nicht möglich; (3) der Zweck „kein Einzellauf als Zusage“ ist durch die bedingungsgleiche Formulierung und die Spalte *Läufe* getragen. Wie bei der Vorgänger-Messung (DoD 3 dort) wirkt die Abschwächung durch die Implementer-Rolle erst mit dieser Entscheidung; der Planner trägt sie in DoD 2 und 3 ein (Wortlaut „ein Lauf, Kontingent“) und streicht das Kennzeichen beim Schließen. Die Verschwendung von 6 Prompts (Nr. 5, 7, 30 bis 33) ist ausgewiesen und vom Skriptentwurf verschuldet; sie hat die Läufe-Lücke mit verursacht und gehört in die Closure-Notiz (Lerneintrag), nicht in die Messung.

## Aussagen nicht weiter als die Messung

Stichproben gegen Rohdaten bestätigt: Prompt-Zähler 40 (Weg a 22, Weg b 18); Relay-Bindung `172.17.0.1`; `restricted`: `allowed` leer, `claude -p` Exit 124; K3c-B Antwort; `OLLAMA_NO_CLOUD` ohne Verkehr; `[1m]` vom Client abgelehnt. Die Fassung nach der Fix-Runde nimmt F-1 (kein „verlangte“ mehr), F-2 (zwei Antwort-Kriterien), F-5 (K6 auf den Ausgabenunterschied begrenzt), F-3/F-4 (Zeitbezug) auf. Einzige Rest-Überreichung: Weg (b), siehe V-1.

### Gegenproben-Zeitbezug

Wie im Protokoll ausgewiesen, von mir gegen die Commit-Historie geprüft: Artefakt des Vorlaufs `5784e921`/`6bfc17e8` vor dem ersten Prompt; `rohdaten-09-gegenproben.txt` ohne Zeitstempel und nach Weg (a) neu gefahren. Das Protokoll behauptet vor dem ersten Prompt nur (1) bis (3) und sagt, dass (4) bis (8) später entstanden: stimmig.

### Zugangsregeln

1. Keine Schlüsseldatei des Hosts: Skripte mounten nur den Messordner, `docker.sock`, das Skriptverzeichnis `:ro` und das Projekt; `id_ed25519` steht nur als `find`/`df`-Pfad im Container (`inner.sh`); der Host-Server wurde nur aus Containern angesprochen. Kein Befund.
2. Kein Schlüssel, keine Anmelde-URL mit Query, kein Token, kein Konto im Repo: Suchen leer, siehe Tabelle. Kein Befund.
3. Gegenprobe vor dem ersten Lauf mit Artefakt: ja (Historie), siehe DoD 1.

## Weg (b): Zuordnung des Geräteschlüssels und Unabhängigkeit von Weg (a)

Nachträge des Projektinhabers (2026-10-24): In den Kontoeinstellungen steht **ein** Geräteschlüssel (Name entspricht seinem Host-Rechner), und **keiner trägt etwas von heute**; kein Eintrag passt zur Anmeldung im Wegwerf-Container.

### Was das Protokoll behauptet

* Zeile 171: „Die Anmeldung hängt damit am Volume (am Geräteschlüssel `id_ed25519`), nicht am Container.“ Gemessen ist nur: Neuanlage des Containers mit demselben Volume, `ollama run` mit Exit 0 und Ausgabe ohne Meldungsmuster. **Dass die Anmeldung am Schlüssel hängt, ist nicht gemessen** (kein Test mit entferntem Schlüssel bei erhaltenem Rest; Schlüsselinhalt nie gelesen; die Anmeldung selbst verändert keine Datei, Zeile 170).
* Zeilen 173, 234, 280 sowie Slice §2/§6: „der Geräteschlüssel **dieses Wegwerf-Containers** in den Kontoeinstellungen“. Das setzt voraus, dass die Anmeldung im Container einen eigenen, im Konto sichtbaren Geräteeintrag erzeugte. Das ist **Plan-Annahme** (Slice §6, Frage 3), aus der Messung nicht belegt, und durch die Angabe des Projektinhabers jetzt unbestätigt: es gibt keinen solchen Eintrag.
* Zeile 170 enthält den richtigen Vorbehalt („ob die Anmeldung als Zustand des Servers im Konto liegt, ist nicht gemessen“), steht aber im Widerspruch zu den Zeilen 171, 173 und 205 („ob der Server den Schlüssel im Konto prüft“ ist ungemessen; „nach dem Entfernen des Geräteschlüssels im Konto“ unterstellt, dass es ihn gibt).
* Zeile 172 („Geräteschlüssel und [`LH-FA-DEV-009`](../../spec/lastenheft.md#lh-fa-dev-009--git-zugangsdaten-im-sandbox-devcontainer): Ort Modell-Volume …“) belegt den **Ort der Datei**, nicht ihre Rolle bei der Anmeldung. Die Folge-ADR darf daraus nicht ableiten, die Datei sei das Zugangsdatum der Cloud-Anmeldung.

### Erklärt das Protokoll, wie der Container ohne neu registrierten Schlüssel Cloud-Modelle ansprach?

**Nein.** Das Protokoll behauptet „Anmeldung im Container“ (K4/K5) und führt die Antworten darauf zurück, erklärt aber nicht, welche Kontobindung der Container nutzte. Folgende Erklärungen sind **mit den Rohdaten nicht unterscheidbar** und deshalb offen (keine ist gemessen): (a) die Anmeldung registrierte den Container unter einem Eintrag, den der Projektinhaber so nicht erkennt oder der später zusammenfiel; (b) die Anmeldung nutzte eine bestehende Kontobindung, ohne einen neuen Schlüssel anzulegen; (c) die Antworten kamen nicht über eine Anmeldung des Containers (Pfad über einen anderen Server, siehe unten); (d) die Anzeige der Kontoeinstellungen listet nicht alles. Das Protokoll darf keine davon vorwegnehmen.

### Ist Weg (b) belegt unabhängig von Weg (a)?

Belegt (Rohdaten `rohdaten-06-*`, Skript `06-lauf-container.sh`):

* Der Container `…-b` entsteht per `docker run -d --cap-add=NET_ADMIN -v <Volume>:/home/vscode/.ollama … sleep` im Standardnetz, **ohne** `--add-host`; die Zeile für `host.docker.internal` in `/etc/hosts` ist **leer** (`rohdaten-06-vorbereiten-1-114212.txt` Z. 7, `-restricted-` Z. 50).
* Die Umgebung des Containers trägt keine Variable `OLLAMA_HOST` und keine verbotenen Namen (`env-namen: keine`; `umgebung-container-verbotene-namen=0`); `claude -p` läuft gegen `http://127.0.0.1:11434`, `ollama run`/`launch` ohne `OLLAMA_HOST`.
* Der Server im Container meldet `{"version":"0.40.6"}` (`serve bereit`), der Host-Server `0.40.10`: der angesprochene Server ist der Container-Server, nicht der Host-Server.
* Der Resolver des Containers sah bei `ollama run` die Namen `ollama.com` (2) und `registry.ollama.ai` (4): der Container-Server hat die Cloud-Dienste selbst angefragt.
* Der Relay lief in Weg (b) nicht (kein Aufruf in `06-lauf-container.sh`); der Host-Server lauscht nur auf `127.0.0.1` (`ss`, `rohdaten-04-probe-113412.txt`), ohne Relay ist er aus dem Container nicht erreichbar (K1 gemessen). `08-aufraeumen.sh alles` belegt `relay-container-nachher=0` (nach Weg (b)).
* Das Meldungsmuster `MELDUNG` enthält `sign`; „You need to be signed in“ hätte als Meldungszeile erscheinen müssen (Vorgänger ohne Konto lieferte sie bei Exit 0).

**Nicht belegt:**

* In den Rohdaten von Weg (b) steht **keine** Prüfung „Relay entfernt“ und keine Container-Liste (Weg-(a)-Fall-Dateien haben die Zeile `relay: entfernt`, die Weg-(b)-Dateien nicht). Dass zur Zeit von Weg (b) kein Relay lief, folgt aus dem Skript und dem Aufräumen davor (`08-aufraeumen.sh weg-a`), nicht aus einer Messung in den Weg-(b)-Rohdaten. Der Relay hätte ohnehin nur über `host.docker.internal`/`172.17.0.1` gewirkt, und nichts im Container zeigt dorthin (kein `OLLAMA_HOST`, kein `--add-host`): die Abhängigkeit ist **unwahrscheinlich, aber nicht durch eine eigene Messung ausgeschlossen** (kein `ss`/`curl` auf das Gateway im Weg-(b)-Container).
* Der Weg, auf dem der Container-Server die Cloud-Modelle **authentifiziert** ansprach, ist nicht beobachtet: der Resolver zeigt Namen, nicht Verbindungen oder Kontobindung; die Anmeldung selbst erzeugte keine Namensanfrage und keine Dateiänderung (Zeile 169/170). Das Protokoll sagt es in *Nicht gemessen* („Verkehr des Servers im Container zum Cloud-Dienst“), zieht aber daraus nicht die Folge für die Aussage „angemeldet im Container“.
* Ob die Anmeldung wirklich durch den Container und nicht durch eine schon vorhandene Kontobindung wirkte, ist offen (Reihenfolge: Server-Start legte `id_ed25519` an, Anmeldung durch den Projektinhaber, danach Antworten; kein Aufruf **vor** der Anmeldung mit einer Cloud-Antwort-Gegenprobe im selben Container; der Vorlauf `vorbereiten` enthielt keinen Prompt).

### Beobachtung des Projektinhabers: Kontoseite „Recent requests“ (ohne Schluss)

Angabe des Projektinhabers (2026-10-24, Screenshot, nicht im Repo, von mir nicht gesehen): 24 sichtbare Zeilen (Liste unten abgeschnitten), je nur Modellname (`glm-5.3-flash`, `gpt-oss:120b`) und relative Zeit, keine Spalte für Gerät, Schlüssel oder Herkunft; Gruppen etwa 56 bis 53, 43 bis 42 und 20 bis 19 Minuten vor der Aufnahme.

Einordnung: Das stützt nur, dass Anfragen beider Modelle am Konto ankamen. Eine Zuordnung zu Weg (a) oder Weg (b) ist aus der Angabe nicht ablesbar, und sie klärt weder die fehlende Schlüsselzeile noch die Unabhängigkeit von Weg (b); beides bleibt oben aus den Rohdaten beurteilt. Die Zeitgruppen sind nicht gegen die Rohdaten abgeglichen (Aufnahmezeit unbekannt, die Liste ist abgeschnitten) und tragen keinen Schluss. Als Beobachtung in der Closure-Notiz festhalten, nicht in das Protokoll als Messwert.

### Offene Punkte / Einschränkungen (Auftrag an den Implementer, im Protokoll auszuweisen)

* **O-1** Zuordnung Konto/Gerät ungemessen: Aussagen „Geräteschlüssel dieses Wegwerf-Containers im Konto“ (Z. 173, 205, 234, 280; Slice §2, §6) als Annahme kennzeichnen. Der Rückbau-Schritt („Schlüssel im Konto entfernen“) ist **gegenstandslos oder anders zu fassen**, solange kein passender Eintrag existiert; den Befund des Projektinhabers (nur ein Schlüssel, kein heutiger Eintrag) als Beobachtung ins Protokoll, ohne Namen und ohne Kontodaten.
* **O-2** Zeile 171 („hängt am Geräteschlüssel“) auf das Gemessene begrenzen: „hängt am Volume; welche Datei oder welcher Kontozustand die Anmeldung trägt, ist nicht gemessen“.
* **O-3** In *Nicht gemessen* aufnehmen: Kontobindung der Anmeldung im Container (ob und wie der Container im Konto erscheint), Authentifizierungsweg der Cloud-Anfragen des Container-Servers, kein Relay-/Gateway-Nachweis in den Weg-(b)-Rohdaten.
* **O-4** Die Folge-ADR bekommt zur Frage „Geräteschlüssel als Zugangsdatum“ die Messung des **Ortes** (Modell-Volume, überlebt Neuanlage), nicht die der **Rolle**; die Wertung ist dort offen zu lassen.
* **O-5** Die Closure-Bedingung „Rückbau vom Projektinhaber bestätigt“ (Slice §6 Risiko Geräteschlüssel, §4) ist neu zu fassen: Bestätigt wird der Zustand des Kontos (kein Eintrag aus der Messung), nicht das Entfernen eines bestimmten Schlüssels.

Diese Einschränkungen berühren nicht die Zugangsregeln (kein Schlüssel gelesen) und nicht die Messwerte; sie begrenzen die **Lesart** von K4/K5.

## Findings

| ID | Kat. | Befund | Pfad | Empfehlung |
|---|---|---|---|---|
| V-1 | MEDIUM | Weg (b) behauptet Anmeldung im Container und Schlüssel-Bindung stärker als gemessen; eigener Geräteeintrag im Konto ist durch den Befund des Projektinhabers nicht bestätigt | Protokoll Z. 170 bis 173, 205, 234, 280; Slice §2, §4, §6 | O-1 bis O-5 |
| V-2 | LOW | DoD 2/3 „zwei Läufe“ bei `[x]` und ein Lauf im Protokoll; Abweichung nicht im DoD gekennzeichnet | Slice §2 | Planner kennzeichnet („ein Lauf, Kontingent 40/40, Rückführung §4 (c)“), siehe F-11 |
| V-3 | INFO | Weg-(b)-Rohdaten tragen kein „Relay entfernt“ und keine Container-Liste | `rohdaten-06-*` | in *Nicht gemessen* aufnehmen (O-3) |
| V-4 | INFO | Review-F-12 (Rückbau im Konto) ist durch den Befund des Projektinhabers neu zu lesen (O-1, O-5) | Review F-12 | Closure-Notiz |

## Nicht ausgeführt

* `make test-docker`, `make verify-depguard`: kein Produktcode, keine Importregel berührt.
* `make ci`, `make fullbuild`: nicht Teil der DoD.
* Wiederholung der Messung, Lesen von Schlüsseldateien, Aufruf des Host-Ollama: nach Auftrag ausgeschlossen.

## Commit / Artefakt

HEAD der Prüfung: `8283435a`.
