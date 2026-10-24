# Review-Report (Nachlauf): slice-cr-lastenheft-ollama-cloud-modelle — 2026-10-24

**Review-Art:** Nachlauf-Review über den Lastenheft-Diff der Fassung 2 (Lastenheft `0.23.0` auf `0.23.1`, Commit `9c8c16eb`). Gelesen wird genau dieser Diff (`git show 9c8c16eb -- spec/lastenheft.md` und der übrige Diff des Commits: Slice-Notiz, `open/README.md`, `open/slice-cr-lastenheft-messstaende.md`) gegen Plan (Slice §1 Fassung 2, freigegebener Wortlaut der Stellen 3, 4, 7; §6; §7), ADR-0036 (Accepted; Z1 bis Z13, *Nicht gemessen, nicht zugesagt*), die Messprotokolle (nur Gemessenes) und den Vorreport `docs/reviews/2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle.md` (F-1 bis F-12). Die DoD ist Sache des Verifiers. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `9c8c16eb` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24 · Arbeitsbaum zu Beginn sauber.

**Sensor-Läufe:** `make docs-check` (857 Dateien, 0 Befunde), `make doc-trace` (84 Anforderungen, 0 Waisen). Der Wortvergleich lief als `grep -c -x -F` auf dem Host (reiner Text, kein Nachweis im Sinn von AGENTS.md §3.1). `make gates`, `make test`, `make test-docker`, `make verify-depguard` nicht gelaufen: der Diff enthält keinen Code, keine Importregel, keinen Docker-Pfad.

**Ergebnis:** kein HIGH. Ein MEDIUM, drei LOW, drei INFO. Der Text im Lastenheft `0.23.1` ist mit dem freigegebenen Wortlaut der Fassung 2 Zeichen für Zeichen identisch. F-1 bis F-3 und F-5, F-6 sind im Text behoben; F-4 ist bis auf einen Rest behoben (Befund N-1).

---

## Findings

```text
MEDIUM spec/lastenheft.md · "Ein neu angelegter Container mit demselben Volume antwortet ohne neue Anmeldung; das ist nur für ein Sandbox-Profil mit Egress-Restriktion und dem Eintrag von `ollama.com` zugesagt." - Neuanlage-Zusage trägt den Weg des Eintrags und den Start-Pfad nicht
Klasse: Zusage über die Handvariante hinaus auf die Auslieferung
Quelle: ADR-0036 Z4 und Abschnitt "Nicht gemessen, nicht zugesagt" ("`restricted` mit Konto und Eintrag über `config set` (nur Handvariante, 1 Lauf)"); Protokoll mit Konto, K8-Zeile und Tabelle der Handvarianten ("`devcontainer up` der Dev-Container-CLI statt `u-boot devcontainer up`"); AGENTS.md §3.10; Vorreport F-4
Befund: Z4 ist ein Lauf: Container der Dev-Container-CLI im `restricted`-Aufbau mit Workspace-Volume und Startbefehl, der Eintrag `ollama.com` als "dieselben zwei Zeilen wie im Egress-Skript" von Hand gesetzt. Der Satz nennt weder "von Hand" noch den Start-Pfad; der Absatz weist den Nutzer zugleich auf `devcontainer.sandbox.egress.allow` (also `config set`) als Ort des Eintrags, und diesen Weg mit Konto führt die ADR als nicht gemessen. Die erste Nicht-zugesagt-Position nimmt den Container aus, den `u-boot devcontainer up` mit Workspace-Volume und Startbefehl startet; der Neuanlage-Container in Z4 ist von dieser Art. Ob die Neuanlage dort zugesagt ist oder ausgenommen, folgt aus dem Text nicht.
Risiko: Nutzer, Code-Slice und Verifier lesen die Neuanlage-Zusage als Zusage für den Ablauf `config set` + `u-boot devcontainer up` + Neuanlage; gemessen ist er so nicht. Das ist der Rest von F-4 (Eintrag über `config set` mit Konto fehlt in der Liste) und die Klasse, die das Register mit 24 Belegen führt. Die Spannung zwischen den beiden Sätzen kann der Verifier nicht per Satz-Vergleich auflösen, weil beide wortgleich zur Freigabe sind.
Verifizierbar: Review-only (Messung oder Streichung des Satzes ist Entscheidung des Projektinhabers; G3 hat die Frage zum Start-Pfad gestellt und den Standard "so" freigegeben, die Frage nach dem Weg des Eintrags nicht)
```

```text
LOW spec/lastenheft.md · "scheitert der Aufruf eines Cloud-Modells ohne Anmeldung an der Namensauflösung von `ollama.com`" - "Aufruf" verallgemeinert den gemessenen `ollama run`
Klasse: Verhaltensaussage weiter als die Zusage-Zeile
Quelle: ADR-0036 Z5 (`ollama run` endet mit Exit 1 und `lookup ollama.com`); Protokoll ohne Konto Abschnitt G ("`run` mit Lookup-Fehler"); AGENTS.md §3.10
Befund: Z5 misst den Lookup-Fehler für `ollama run`. Der Satz sagt "Aufruf eines Cloud-Modells", und der Absatz nennt als Aufrufe `ollama run`, `ollama launch claude` und `claude -p`; für die beiden anderen hinter der Egress-Restriktion ohne Konto belegt die Zeile den Lookup-Fehler nicht.
Risiko: der nächste Änderer pinnt oder dokumentiert den Lookup-Fehler für alle drei Aufrufe (Handbuch, Spezifikation) und folgt einer Aussage, die nur für einen Aufruf gemessen ist.
Verifizierbar: Review-only
```

```text
LOW docs/plan/planning/welle-ollama-cloud-modelle.md · "die Version des CR dieser Welle wird beim Start des CR-Slice nach dem Stand des Lastenhefts vergeben (heute 0.23.0)" - Zustandsmarker "heute" stimmt nicht mit dem Ist-Stand
Klasse: Zustandsmarker nach Ablauf nicht auf Ist-Stand gesetzt
Quelle: AGENTS.md §3.10 (Zustandsmarker mit Ablaufpunkt), §3.7
Befund: Die Welle-Datei (Frage 3, Vorschlag des Planners) nennt "heute 0.23.0"; das Lastenheft trägt `0.23.1`, und zum Zeitpunkt des Vorschlags war es `0.22.0`. Die Zeile in §4 ("Version wurde beim Start des CR-Slice vergeben (0.23.0 …)") nennt die Fassung 2 und `0.23.1` nicht. Die Slice-Notiz §7 führt beide Nennungen als "gelesen, nicht geändert" und bewertet sie als Vergabe-Beschreibung; "heute" ist ein Zustandswort.
Risiko: Wer die Welle zur Closure liest ("die vergebene Lastenheft-Version bei der Closure", Plan-Tabelle §3), übernimmt `0.23.0` als Stand des Lastenhefts.
Verifizierbar: ja, `grep -n -F "heute 0.23.0" docs/plan/planning/welle-ollama-cloud-modelle.md`
```

```text
LOW spec/lastenheft.md · "Der Server legt beim ersten Start einen Geräteschlüssel im Volume der Modelle an, nicht im Workspace-Volume." - Satz steht im Absatz mit dem Zugesagt-Aufbau "ohne Workspace-Volume"
Klasse: Aufbau eines Satzes im Absatz nicht benannt
Quelle: ADR-0036 Z6 (Pfad: Auslieferung bis auf `devcontainer up` der CLI; Bestätigung mit Konto: Handvariante); Maintainability
Befund: Der Aufbau-Satz des Absatzes ("ohne Workspace-Volume und mit einem Volume für die Modelle") und die Egress-Sätze ("Hinter der Egress-Restriktion …", Z5, Z6) gehören zu verschiedenen gemessenen Aufbauten (Z5 und Z6: Container der CLI mit Workspace-Volume und Egress-Skript). Der Absatz sagt "Zugesagt ist der Ablauf nur in dem Aufbau …" und bindet damit nur den Ablauf, nicht die Verhaltenssätze davor; der Leser erfährt nicht, in welchem Aufbau Z5 und Z6 gelten.
Risiko: gering; der Verifier muss Z5 und Z6 gegen einen anderen Aufbau halten als Z1 bis Z3 und kann sich auf den Aufbau-Satz nicht stützen.
Verifizierbar: Review-only
```

```text
INFO spec/lastenheft.md · "warnt beim Erzeugen, ohne abzubrechen" - Sichtbarkeit der Zusatzwarnung unter `--quiet` bleibt dem Vertragstext nach offen
Klasse: Sichtbarkeit einer Zusatzwarnung unter --quiet offen
Quelle: Vorreport F-10; Slice §1 Frage G8 (Standard: Stelle 6 unverändert); `LH-FA-CLI-005` ("Alle anderen Hinweise und Warnungen entfallen unter `--quiet`")
Befund: Mit G8 entschieden und unverändert; der Test des Code-Slice pinnt die Lesart. Der Vertragstext allein löst sie nicht auf.
Risiko: kein; Hinweis für den Code-Slice und den Verifier.
Verifizierbar: Test des Code-Slice
```

```text
INFO spec/lastenheft.md · "Zugesagt ist der Ablauf nur in dem Aufbau, in dem er geprüft ist" - Zusage trägt einen Prüfstand im Wortlaut
Klasse: Zielbild trägt Prüfzustand
Quelle: AGENTS.md §3.4 (Lastenheft ohne Prozesszustand), Vorgabe 2026-10-19 (Messstand nicht im Lastenheft)
Befund: "in dem er geprüft ist" und "in dem es geprüft ist" nennen einen Prüfzustand, keinen Messwert; der Aufbau steht direkt dahinter im Wortlaut. Kein Messstand, keine Version, kein Host im Sinn der Vorgabe.
Risiko: kein; Hinweis für den Fall, dass der Projektinhaber den Wortlaut ohne Prüfzustand will.
Verifizierbar: Review-only
```

```text
INFO docs/plan/planning/in-progress/slice-cr-lastenheft-ollama-cloud-modelle.md · "kein Commit von `spec/lastenheft.md` vor diesem Eintrag" - Freigabe-Eintrag und Lastenheft-Commit im selben Commit
Klasse: Freigabe-Beleg und Änderung im selben Commit
Quelle: Slice §2 Liefer-Punkt 1; Verifier-Aufgabe
Befund: §7 *Freigabe der Fassung 2* und der Lastenheft-Diff liegen im Commit `9c8c16eb`; die Reihenfolge "Eintrag vor Commit" ist aus `git log` nicht ablesbar. Bei Fassung 1 verlief es gleich; der Vorreport hat es nicht beanstandet.
Risiko: kein Vertragsbruch; der Verifier belegt die Freigabe über §7 und den Gesprächsstand.
Verifizierbar: Verifier
```

---

## Prüfung der Vorreport-Befunde F-1 bis F-6 gegen den Text (Prüfpunkt 2)

| Befund | Stand im Text `0.23.1` | Urteil |
|---|---|---|
| F-1 Reichweite über die Handvariante | Satz 1 "im Container des Sandbox-Profils … in dem Aufbau, den dieser Absatz als zugesagt benennt"; Aufbau-Satz (Image mit `ollama`, `node`, `claude-code`, vom Nutzer gestartet, `NET_ADMIN`, ohne Egress-Skript, ohne Workspace-Volume, mit Volume für Modelle, ohne Terminal, Linux-`amd64`, lokale Engine); Nicht zugesagt nennt den Container von `u-boot devcontainer up` und das Sandbox-Profil ohne Egress-Restriktion. Aufbau gegen Protokoll mit Konto gelesen: Weg (b) `discover` mit Modell-Volume (`rohdaten-06-vorbereiten-1-114212.txt`), ohne Workspace-Volume. "ohne Egress-Restriktion und damit ohne `NET_ADMIN`" gegen Code gelesen (`internal/hexagon/application/devcontainer_sandbox.go`, `egressRelaxation`). | behoben, Rest siehe N-1 |
| F-2 Bedingungen von `claude -p` (Stelle 3) | Basis-URL, Platzhalter-Token, kein `ANTHROPIC_API_KEY` im Wortlaut; Features `claude-code`, `node` im Aufbau-Satz. Rohdatei `rohdaten-06-discover-1-114507.txt` nennt `ANTHROPIC_BASE_URL` und `ANTHROPIC_AUTH_TOKEN` als Umgebungsnamen des Aufrufs. | behoben |
| F-3 Bedingungen im Anschluss-Absatz | Basis-URL, Platzhalter-Token, kein `ANTHROPIC_API_KEY`; Features `claude-code`, `node`; Aufbau Standardprofil oder Container des Sandbox-Profils; "eine bestimmte Version von Claude Code und des Dienstes des Hosts" nicht zugesagt (Z9/Z10: Server 0.40.10, abweichend von der gepinnten). | behoben |
| F-4 Nicht-zugesagt-Listen | Cloud-Absatz: Hosts des Servers zum Cloud-Dienst, Container von `u-boot devcontainer up`, Anmeldung hinter Egress über die Neuanlage hinaus, zweiter Lauf, Rolle des Geräteschlüssels, Kontobindung. Anschluss-Absatz: Kontingent und Kosten, Modell, interaktiv, zweiter Lauf, Wirkung des erzeugten Eintrags. Gegen die Liste der ADR gelesen: alle berührten Positionen stehen, ausgenommen der Rest in N-1. | behoben bis auf N-1 |
| F-5 zwei Zusage-Mengen | Satz "Diese Zusage gilt für den Weg zu Cloud-Modellen; sie ist von der Zusage zur Installation von `ollama` im Absatz davor getrennt und ändert deren Grenzen nicht." Der Messstand-Halbsatz im Absatz `ollama` steht unverändert (Stelle 2, Frage F3). Die Spannung ist im Text ausgewiesen und im Inventar geführt (zwei temporäre Zeilen in `carveouts.md`, Plan-Anker `slice-cr-lastenheft-messstaende`, Trigger "liegt in `done/`", Bindung vor dem Release-Cut bei G6/Projektinhaber). | behoben (ausgewiesen, nicht aufgelöst; Auflösung beim Folge-Slice) | <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
| F-6 "geänderte Bindung" | "für einen Dienst mit anderer Bindung gibt es keine Zusage, auch nicht für eine Bindung, die der Nutzer ändert"; Z11 trägt den Loopback-Satz. | behoben |

Keine Zusage über Z1 bis Z13 hinaus, abgesehen von N-1 und dem LOW zu "Aufruf". Z13 bleibt ein Produktversprechen ohne Beleg bis zum Code-Slice (INFO F-11 des Vorreports, weiter beim Code-Slice und Verifier).

## Wortlaut-Vergleich (Prüfpunkt 1)

Methode: Zeilen aus Slice §1 (Blockzitate Fassung 2) mit dem Pfadpräfix `../../../../spec/lastenheft.md#` durch `#` ersetzt, dann `grep -c -x -F` gegen die Datei.

| Stelle | Ergebnis |
|---|---|
| 3 Fassung 2 (Zeile 414) | ganze Zeile identisch, 1 Treffer (3351 Zeichen) |
| 4 Fassung 2 (Zeile 418) | ganze Zeile identisch, 1 Treffer (3017 Zeichen) |
| *Alt* Stelle 3 und 4 | je 1 Treffer in `git show 3b135e0e:spec/lastenheft.md` (Lastenheft `0.23.0`) |
| 7 Historie-Zeile `0.23.1` (Zeile 1605) | identisch mit dem Entwurf, `JJJJ-MM-TT` durch `2026-10-24` ersetzt; 1 Treffer |
| 7 Versionsfeld (Zeile 9) | `0.23.1` |
| Historie-Zeile `0.23.0` (Zeile 1604) | unverändert (Chronik) |
| Stellen 1, 2, 5, 6, 8 | unverändert gegenüber `0.23.0`: `git diff 3b135e0e 9c8c16eb -- spec/lastenheft.md` hat drei Hunks (Zeile 9, Zeilen 414 bis 418, Zeile 1605), 4 Einfügungen, 3 Löschungen; Stellen 1 (410), 2 (412), 5 (460), 6 (191) liegen außerhalb |

## Negativbefunde

Geprüft, ohne Befund:

- Prüfpunkt 1: Wortlaut Stellen 3, 4, 7, Versionsfeld, Unverändertheit der Stellen 1, 2, 5, 6, 8 (Tabelle oben).
- Prüfpunkt 3: Kennungen in den drei geänderten Zeilen: keine Treffer für `ADR-`, `Slice`, `slice-`, `Welle`, `SPEC-`, `ARC-`. Hostnamen: nur `ollama.com`; kein `host.docker.internal`. Keine Versionsnummern, keine Engine-Versionen (kein `Docker 29`, kein `0.40`), keine Modellnamen. Verweise: ausschließlich Anker innerhalb der Datei (`#lh-fa-dev-008--egress-restriktion`, `#lh-fa-dev-006--sandbox-profil`, `#lh-fa-dev-003--devcontainer-features`); Referenzmodell (Decken-Regel, AGENTS.md §3.4) eingehalten, kein HIGH nach Skill-Anker "Referenzmodell-Verstoß".
- Prüfpunkt 4: `LH-FA-DEV-003` Absatz `ollama` (Installation) gegen Cloud-Absatz und Anschluss-Absatz: die Spannung der zwei Mengen ist im Text und im Inventar ausgewiesen (F-5); sonst kein Widerspruch. `LH-FA-DEV-006` (Zeile 460, Netzanschluss in der Klammer) und `LH-FA-CLI-005` Ausnahme 1 (Zeile 191, "Vier Ausnahmen" unverändert): konsistent mit dem Anschluss-Absatz ("ausgewiesene Sicherheitslockerung" im Sandbox-Profil). `LH-FA-DEV-009` (Geräteschlüssel im Modell-Volume, nicht im Workspace-Volume): kein Widerspruch zum Cloud-Absatz. Chronik `0.15.0` (Zeile 1596) unverändert. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Carveouts: zwei temporäre Zeilen in `docs/plan/planning/in-progress/carveouts.md` (Widerspruch der tieferen Straten und des Codes zum Lastenheft; zwei Zusage-Mengen), je mit Plan-Anker und Trigger; Inhalt gegen den Lastenheft-Text gelesen, konsistent.
- Prüfpunkt 5: Plan-Tabelle §3 gegen `git diff --stat 9c8c16eb~1 9c8c16eb` (4 Dateien: `spec/lastenheft.md`, Slice, `open/README.md`, `open/slice-cr-lastenheft-messstaende.md`): jede Datei hat eine Zeile. Gegen `git diff --stat 3b135e0e 9c8c16eb` (7 Dateien, einschließlich `carveouts.md`, `welle-ollama-cloud-modelle.md`, Vorreport): jede Datei hat eine Zeile. Zustandsmarker im Slice (§1, §2, §6, §7): auf Ist-Stand `0.23.1`. Messstände-CR: Version `0.23.1` und Zeilenangaben 412, 472, 474, 1588, 1593, 1594, 1596 sowie 1101, 433, 1509, 1584 gegen `spec/lastenheft.md` gelesen, stimmen; "gleich dem Stand 0.23.0" stimmt (Datei 1612 Zeilen, Fassung 2 ändert Zeilenlänge, nicht Zeilenzahl, vor Zeile 1604). `open/README.md`: Zeilen beider Slices auf `0.23.1` bzw. "vergeben beim Start: 0.23.0, in der Fix-Runde 0.23.1". <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Prüfpunkt 6: `make doc-trace` 84 Anforderungen, 0 Waisen; `make docs-check` 0 Befunde. ADR-0036 und ADR-0024 unverändert (`git diff --stat 3b135e0e 9c8c16eb -- docs/plan/adr` leer; keine ADR-Änderung seit `d9dd171e`). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Kommentar-Selbstcheck AGENTS.md §3.9 (a) bis (c) auf die neuen Lastenheft-Zeilen: kein "früher"/"no longer"/"former"; die Historie-Zeile `0.23.1` beschreibt die Änderung der Version (Klasse Chronik).
- Sprachvertrag: Lastenheft-Text deutsch, kein CLI-Output im Diff.

Nicht geprüft:

- Nutzer-Doku (`README.md`, `README.de.md`, `docs/user/`), `spec/spezifikation.md`, `spec/architecture.md`: vom Diff nicht berührt; die bewusst stehenbleibenden Widersprüche (SPEC-047 Zeile 739, Handbuch Zeilen 2023 und 2043, `README.de.md` Zeile 14) sind im Inventar geführt. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- `make gates`, `make test`, `make test-docker`, `make verify-depguard`: kein Code, keine Importregel, kein Docker-Pfad im Diff.
- Rohdaten: nur die für die Aussagen gelesenen (`rohdaten-06-discover-1-114507.txt`, `rohdaten-06-vorbereiten-1-114212.txt`); nicht alle Dateien einzeln geöffnet.
- Ob der Projektinhaber die Vorlage der Fassung 2 ungekürzt sah (`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`): aus dem Repo nicht prüfbar; §7 trägt "Abweichung gegen den Entwurf in §1: keine", das Lastenheft stimmt mit §1 überein; Verifier. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage über die Handvariante hinaus auf die Auslieferung · Verhaltensaussage weiter als die Zusage-Zeile · Zustandsmarker nach Ablauf nicht auf Ist-Stand gesetzt · Aufbau eines Satzes im Absatz nicht benannt · Sichtbarkeit einer Zusatzwarnung unter --quiet offen · Zielbild trägt Prüfzustand · Freigabe-Beleg und Änderung im selben Commit

## Verdikt

**Merge-blockierend:** nein für HIGH (keiner). Das MEDIUM zur Neuanlage-Zusage ändert, falls behoben, Lastenheft-Text und braucht eine Freigabe des Projektinhabers (Streichen, Eingrenzen auf den gemessenen Eintrag und Start-Pfad oder als Zusage annehmen mit Begründung in §7); dann folgt ein weiteres Nachlauf-Review über genau diesen Diff vor dem Verifier. Das LOW zur Welle-Datei ist eine Korrektur ohne Vertragswirkung. Ein Widerspruch des Implementers stuft das MEDIUM nicht herab; ab HIGH mit Rollen-Widerspruch oder beim dritten gleichen Konflikttyp läuft der Konfliktpfad über den Architect.

**Übergabe:** Findings gehen an den Implementer; N-1 (MEDIUM) zusätzlich an Architect und Projektinhaber (Rückkante Review → Plan: der Wortlaut folgt dem Entwurf des Planners). Verifier: Satz-für-Satz-Vergleich gegen §7 und Fassung 2, INFO zur Freigabe-Reihenfolge, Z13 bis zum Code-Slice.

## Kennzeichnung der Findings

N-1 = MEDIUM Neuanlage-Zusage; N-2 = LOW "Aufruf" (Z5); N-3 = LOW Welle-Datei "heute 0.23.0"; N-4 = LOW Aufbau von Z5/Z6; N-5 bis N-7 = INFO in der Reihenfolge oben (Warnung unter `--quiet`, Prüfstand im Wortlaut, Freigabe-Reihenfolge).
