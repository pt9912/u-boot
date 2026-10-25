# Review-Report (Nachlauf 2): slice-cr-lastenheft-ollama-cloud-modelle — 2026-10-24

**Review-Art:** zweiter Nachlauf-Review über den Lastenheft-Diff der Fassung 3 (Lastenheft `0.23.1` auf `0.23.2`, Commit `29cb728b`). Gelesen wird genau dieser Diff (`git show 29cb728b -- spec/lastenheft.md` und der übrige Diff des Commits: Slice-Notiz, `open/README.md`, `in-progress/slice-cr-lastenheft-messstaende.md`, `welle-ollama-cloud-modelle.md`) gegen Plan (Slice §1 Fassung 3, Fragen H1 bis H4, §7 *Freigabe der Fassung 3*), ADR-0036 (Accepted; Z1 bis Z13, *Nicht gemessen, nicht zugesagt*), die Messprotokolle (nur Gemessenes) und den Vorreport `docs/reviews/2026-10-24-slice-cr-lastenheft-ollama-cloud-modelle-nachlauf.md` (N-1 bis N-7). Die DoD ist Sache des Verifiers. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

**Skill:** `.harness/skills/reviewer.md` @ HEAD `29cb728b` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24 · Arbeitsbaum zu Beginn sauber.

**Sensor-Läufe:** `make docs-check` (858 Dateien, 0 Befunde), `make doc-trace` (84 Anforderungen, 0 Waisen). Der Wortvergleich lief als `grep -c -x -F` mit Bash-Parameterersetzung auf dem Host (reiner Text, kein Nachweis im Sinn von AGENTS.md §3.1). `make gates`, `make test`, `make test-docker`, `make verify-depguard` nicht gelaufen: der Diff enthält keinen Code, keine Importregel, keinen Docker-Pfad.

**Ergebnis:** kein HIGH, kein MEDIUM. Ein LOW, vier INFO. Der Text im Lastenheft `0.23.2` ist mit dem freigegebenen Wortlaut der Fassung 3 Zeichen für Zeichen identisch. N-1 bis N-4 sind im Text behoben. Die Fassung 3 sagt nichts zu, was Z1 bis Z13 oder die Rohdaten nicht tragen, bis auf den Aufbau-Rest in L-1.

---

## Findings

```text
LOW spec/lastenheft.md · "in einem Container, den die Dev-Container-CLI aus der erzeugten Konfiguration des Sandbox-Profils startet (mit Workspace-Volume und Startbefehl; im Folgenden der CLI-Container)" - CLI-Container ohne Nennung der Egress-Restriktion definiert
Klasse: Aufbau eines Satzes weiter benannt als gemessen
Quelle: ADR-0036 Z6 (Protokoll ohne Konto, Abschnitt F: Modi `discover`, `restricted`, `allow`) und Abschnitt "Nicht gemessen, nicht zugesagt" ("das von `u-boot` erzeugte Sandbox-Profil ohne Egress-Restriktion (ohne `NET_ADMIN`, ohne protokollierenden Resolver, mit Workspace-Volume und `postStartCommand`)"); AGENTS.md §3.10
Befund: Der Geräteschlüssel-Satz gilt im "CLI-Container", definiert als Container der Dev-Container-CLI aus der erzeugten Konfiguration des Sandbox-Profils mit Workspace-Volume und Startbefehl, ohne Bezug zur Egress-Restriktion. Z6 ist nur in Projekten mit Egress-Skript als `postStartCommand` gemessen (`restricted`, `allow`; `discover` ist kein CLI-Container); die ADR führt den Container mit Workspace-Volume und Startbefehl ohne Egress-Restriktion ausdrücklich als nicht gemessen. Die Definition des Wortlauts deckt beide Lagen. Die Egress-Sätze desselben Absatzes sind durch "Hinter der Egress-Restriktion" auf die gemessene Lage begrenzt, der Schlüssel-Satz nicht.
Risiko: gering. Der nächste Änderer (Handbuch, Spezifikation) liest den Schlüssel-Satz als Zusage auch für den Container der CLI ohne Egress-Restriktion, den die ADR ungemessen führt. Die Schlüsselablage ist ein Verhalten von `ollama serve`, das dort kaum abweicht; belegt ist es nicht.
Verifizierbar: Review-only (Eingrenzung ist Entscheidung des Projektinhabers)
```

```text
INFO spec/lastenheft.md · "Hinter der Egress-Restriktion ... scheitert im CLI-Container `ollama run` eines Cloud-Modells ohne Anmeldung" gegenüber "der Ablauf mit Anmeldung in dem Container, den `u-boot devcontainer up` aus der erzeugten Konfiguration startet" - zwei Namen für nahezu denselben Container
Klasse: Start-Pfad als einziger Unterschied zweier Container-Begriffe
Quelle: Protokoll mit Konto, Tabelle der Handvarianten ("derselbe Mechanismus darunter; der Befehl von `u-boot` selbst wurde nicht gefahren"); Slice §1 Frage G3 und H3 (Standard "so")
Befund: Der CLI-Container (zugesagt für Schlüssel-Satz und Egress-Sätze) und der Container von `u-boot devcontainer up` (aus "Nicht zugesagt" für den Ablauf mit Anmeldung) unterscheiden sich im Text nur durch das startende Werkzeug. Die Egress-Sätze ohne Anmeldung sind für den Container von `u-boot devcontainer up` weder zugesagt noch ausgenommen.
Risiko: kein Vertragsbruch; mit H3 und G3 freigegeben. Hinweis für den Verifier und den Code-Slice, dass ein Nutzer den CLI-Container nur über `u-boot devcontainer up` erreicht und ein Test der Egress-Sätze deshalb am Befehl von `u-boot` hängt.
Verifizierbar: Review-only; Test des Code-Slice
```

```text
INFO spec/lastenheft.md · "Der Server legt beim ersten Start einen Geräteschlüssel im Volume der Modelle an, nicht im Workspace-Volume; das gilt im Aufbau des Ablaufs" - "nicht im Workspace-Volume" im Aufbau ohne Workspace-Volume
Klasse: Aussage im Aufbau ohne Gegenstand
Quelle: ADR-0036 Z6 (Bestätigung mit Konto im `discover`-Aufbau, Pfad Handvariante); Maintainability
Befund: Der Aufbau des Ablaufs hat laut Satz "ohne Workspace-Volume"; die Negation "nicht im Workspace-Volume" ist dort leer, die Aussage "im Volume der Modelle" trägt (`rohdaten-06-vorbereiten-1-114212.txt`, Zeilen 22 und 23). Für den CLI-Container trägt die Negation (Abschnitt F des Protokolls ohne Konto).
Risiko: keines; der Satz bleibt wahr.
Verifizierbar: Review-only
```

```text
INFO docs/plan/adr/0036-cloud-modelle-von-ollama-im-devcontainer.md · "Z4 | Neuer Container mit demselben Volume" - Zusage-Menge der ADR größer als der Vertrag
Klasse: Vertrag enger als die Zusage-Menge der ADR
Quelle: ADR-0036 Entscheidung 10 ("Der Lastenheft-Wortlaut enthält keinen Satz, der in dieser Menge keine Zeile hat" - die Gegenrichtung ist nicht verlangt); AGENTS.md §3.5 (Accepted-ADR bleibt unverändert)
Befund: Z4 steht in der Tabelle "Zugesagt"; das Lastenheft führt die Neuanlage jetzt als "nicht zugesagt". ADR-0036 und ADR-0024 sind gegenüber 3b135e0e unverändert (`git diff --stat` leer).
Risiko: kein; der Spezifikations-Slice und das Handbuch führen Z4 als Handvariante ohne Zusage, nicht als Zusage. Hinweis für Architect und Spezifikations-Slice.
Verifizierbar: Review-only
```

```text
INFO docs/plan/planning/done/slice-cr-lastenheft-ollama-cloud-modelle.md · "Zugesagt ist der Ablauf nur in dem Aufbau, in dem er geprüft ist" (Wortlaut, wie in N-6) - Prüfzustand im Zielbild; Freigabe-Beleg und Lastenheft-Commit im selben Commit (wie N-7); Z13 bis zum Code-Slice ohne Beleg
Klasse: bekannte INFO, unverändert
Quelle: Vorreport N-6, N-7; Vorreport F-11 des ersten Reviews
Befund: Alle drei Hinweise bestehen unverändert fort; die Fassung 3 berührt sie nicht.
Risiko: kein; Hinweis für den Verifier (Freigabe über §7 und Gespräch, Z13 über den Test des Code-Slice).
Verifizierbar: Verifier
```

---

## Prüfung der Vorreport-Befunde N-1 bis N-4 gegen den Text (Prüfpunkt 2)

| Befund | Stand im Text `0.23.2` | Urteil |
|---|---|---|
| N-1 Neuanlage-Zusage (MEDIUM) | Der Satz "Ein neu angelegter Container mit demselben Volume antwortet ohne neue Anmeldung; …" fehlt in der Zusage (`grep` auf "neu angelegt": Zeile 414 nur in "Nicht zugesagt", sowie Historie-Zeile 1606). "Nicht zugesagt" führt "ein Aufruf mit Anmeldung hinter der Egress-Restriktion, mit oder ohne den Eintrag von `ollama.com`, auch in einem neu angelegten Container mit demselben Volume". Z4 ist Handvariante (`restricted`-Aufbau der CLI, Eintrag von Hand, 1 Lauf); der Weg über `config set` mit Konto steht in der ADR als nicht gemessen. Der Wortlaut sagt dazu nichts zu. | behoben |
| N-2 "Aufruf" verallgemeinert (LOW) | Lookup-Fehler und Anmeldeaufforderung sind an `ollama run` gebunden ("scheitert im CLI-Container `ollama run` eines Cloud-Modells ohne Anmeldung …; mit dem Eintrag fordert `ollama run` ohne Anmeldung zur Anmeldung auf"). Z5 und das Protokoll ohne Konto, Abschnitt C (`ollama run`, Exit 1, `lookup ollama.com`; C2b Anmeldeaufforderung, Exit 0) decken das. `ollama launch claude` und `claude -p` ohne Anmeldung hinter der Egress-Restriktion stehen in "Nicht zugesagt". Gegen Abschnitt C gelesen: `launch claude` startet dort das installierte `claude` mit anderem Fehler, ist also nicht gleich zu `run`. | behoben |
| N-3 Zustandsmarker "heute 0.23.0" (LOW) | In `welle-ollama-cloud-modelle.md` nennt Frage 3 und die Zeile des CR in §4 "Lastenheft-Stand nach der Fassung 3 des Wortlauts: 0.23.2, gelesen in `spec/lastenheft.md` Zeile 9"; Zeile 9 trägt `0.23.2`. Kein "heute"-Wort mehr, die Vergabe `0.23.0` bleibt als Vergabe beschrieben. `grep` auf "heute 0.23.0" in `docs/plan/planning`: kein Treffer. | behoben |
| N-4 Aufbau von Z5/Z6 (LOW) | Schlüssel-Satz: "das gilt im Aufbau des Ablaufs, den dieser Absatz weiter unten nennt, und in einem Container, den die Dev-Container-CLI … startet (… der CLI-Container)"; Egress-Sätze: "im CLI-Container". Gegen Rohdaten gelesen: Der CLI-Container (Workspace-Volume als Docker-Volume `/workspaces/<projekt>`, Modell-Volume `/home/vscode/.ollama`, Egress-Skript als `postStartCommand`) ist der Aufbau von Z5 und Z6. Rest siehe L-1. | behoben bis auf L-1 |
| N-5 bis N-7 (INFO) | unverändert fortgeführt (vierter Eintrag oben) | weiter INFO |

Die übrigen Sätze der Stelle 3 sind zeichengleich zur Fassung 2 (Diff `9c8c16eb..29cb728b`: drei Hunks, nur Zeile 9, 414, 1606). Geräteschlüssel-Satz: im Aufbau des Ablaufs durch `rohdaten-06-vorbereiten-1-114212.txt` (Zeilen 22, 23: `id_ed25519` 387 Byte Rechte 600 und `.pub` 81 Byte, Einhängepunkt `/home/vscode/.ollama`) und im CLI-Container durch das Protokoll ohne Konto, Abschnitt F (Modi `restricted`, `allow`; Schlüssel "vor jedem Anmeldeversuch", auf dem Modell-Volume, "kein Pfad im Workspace-Volume") gedeckt. Der Lookup-Fehler-Satz ist durch Z5 gedeckt; Pfad "Auslieferung bis auf Resolver-Protokoll und `devcontainer up` der CLI", Eintrag über `config set`. Z5 ist ohne Konto gemessen; der Wortlaut sagt "ohne Anmeldung". Keine Zusage über Z1 bis Z13 hinaus, bis auf den Aufbau-Rest in L-1.

## Wortlaut-Vergleich (Prüfpunkt 1)

Methode: Zeile 173 (*Neu* Stelle 3) und 177 (Historie-Zeile) aus Slice §1 mit Bash-Parameterersetzung ohne `> `-Präfix, Pfadpräfix `../../../../spec/lastenheft.md#` durch `#` ersetzt, bei der Historie-Zeile `JJJJ-MM-TT` durch `2026-10-24`, dann `grep -c -x -F -f` gegen `spec/lastenheft.md`.

| Stelle | Ergebnis |
|---|---|
| 3 Fassung 3, *Neu* (3553 Zeichen) | ganze Zeile identisch, 1 Treffer, Zeile 414 |
| 7 Historie-Zeile `0.23.2` | identisch bis auf das Datum, 1 Treffer, Zeile 1606; Zeilen `0.23.0` (1604) und `0.23.1` (1605) unverändert |
| 7 Versionsfeld | `0.23.2` (Zeile 9, Tabellenbreite erhalten) |
| Stellen 1, 2, 4, 5, 6, 8 | unverändert: `git diff 9c8c16eb 29cb728b -- spec/lastenheft.md` hat genau drei Hunks (Zeile 9, Zeile 414, Zeile 1606); Stellen 1 (410), 2 (412), 4 (418), 5 (460), 6 (191) liegen außerhalb |

## Negativbefunde

Geprüft, ohne Befund:

- Prüfpunkt 1: Wortlaut Stelle 3, Historie-Zeile, Versionsfeld, Unverändertheit der übrigen Stellen (Tabelle oben).
- Prüfpunkt 3: Kennungen und Namen in Zeile 414: keine Treffer für ADR, Slice, Welle, SPEC, ARC, `host.docker.internal`, Versionsnummern, Modellnamen; Hostname nur `ollama.com`; Links ausschließlich Anker innerhalb der Datei (`#lh-fa-dev-008--egress-restriktion`); `Colima` nur in der Liste der nicht zugesagten Engines; Feature-/Variablen-/Capability-Namen (`claude-code`, `node`, `NET_ADMIN`, `ANTHROPIC_*`, `OLLAMA_API_KEY`) wie in der Fassung 2 freigegeben. Neu: "Dev-Container-CLI" als Werkzeugname (Frage H3, freigegeben). Referenzmodell eingehalten (kein Verweis abwärts, kein HIGH nach Skill-Anker "Referenzmodell-Verstoß").
- Prüfpunkt 4: Absatz `ollama` (Zeile 412, Installation) gegen Cloud-Absatz (414) und Host-Absatz (418): die Spannung der zwei Zusage-Mengen ist unverändert ausgewiesen und im Carveout-Inventar geführt; Host-Absatz widerspricht dem Cloud-Absatz nicht (beide nehmen das ohne Egress-Restriktion erzeugte Sandbox-Profil bzw. den Container von `u-boot devcontainer up` aus). `LH-FA-DEV-006` (Zeile 460) und `LH-FA-CLI-005` Ausnahme 1 (Zeile 191, "Vier Ausnahmen") unverändert; `LH-FA-DEV-009` (Schlüssel im Modell-Volume) ohne Widerspruch; Chronik `0.15.0` unverändert. Carveouts (`carveouts.md`) vom Diff nicht berührt; Inhalt nennt keinen Satz, den die Fassung 3 ändert. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Prüfpunkt 5: Plan-Tabelle §3 gegen `git diff --stat 9c8c16eb 29cb728b` (5 Dateien: `spec/lastenheft.md`, Slice, `open/README.md`, `in-progress/slice-cr-lastenheft-messstaende.md`, `welle-ollama-cloud-modelle.md`): jede Datei hat eine Zeile. Messstände-CR: Zeilenangaben 412, 472, 474, 1588, 1593, 1594, 1596, 1101, 433, 1509, 1584 gegen `spec/lastenheft.md` (1613 Zeilen) gelesen, stimmen; Version `0.23.2` in Kopf, Frage 2 und DoD. Zustandsmarker in Slice (§1 Kopf Fassung 3 "umgesetzt in Lastenheft `0.23.2`", §7), Welle (Frage 3, Zeile §4) und `open/README.md` auf Ist-Stand; `grep` auf `0.23.0`/`0.23.1` in `docs/plan/planning` außerhalb des Slice: nur `done/slice-v1-release-cut-v0.12.0.md` (andere Bedeutung) und Zeile 33 des Messstände-CR (Verlauf der Versionsfolge, als solcher benannt). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Prüfpunkt 6: `make doc-trace` 84 Anforderungen, 0 Waisen; `make docs-check` 0 Befunde; ADR-0036 und ADR-0024 unverändert (`git diff --stat 3b135e0e HEAD -- docs/plan/adr` leer; letzter Commit an ADR-0036 `d9dd171e`). <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->
- Kommentar-Selbstcheck AGENTS.md §3.9 (a) bis (c) auf die neuen Lastenheft-Zeilen: kein "früher"/"former"/"no longer"; die Historie-Zeile `0.23.2` nennt "nicht mehr zugesagt" und beschreibt damit die Änderung der Version (Klasse Chronik), nicht einen abwesenden Text im Normtext.
- Sprachvertrag: Lastenheft-Text deutsch, kein CLI-Output im Diff.

Nicht geprüft:

- Nutzer-Doku (`README.md`, `README.de.md`, `docs/user/`), `spec/spezifikation.md`, `spec/architecture.md`: vom Diff nicht berührt; die stehenbleibenden Widersprüche sind im Inventar geführt.
- `make gates`, `make test`, `make test-docker`, `make verify-depguard`: kein Code, keine Importregel, kein Docker-Pfad im Diff.
- Rohdaten: nur die für die Aussagen gelesenen (Protokoll ohne Konto, Abschnitte C, D, F, G, Handvarianten; `rohdaten-06-vorbereiten-1-114212.txt`); `rohdaten-06-restricted-1-114624.txt` (Z4) nicht einzeln geöffnet, weil die Fassung 3 die Neuanlage nicht mehr zusagt.
- Ob der Projektinhaber die Vorlage der Fassung 3 ungekürzt sah (`BEO-DOK/freigabe-vorlage-verdichtet-den-entwurf`): aus dem Repo nicht prüfbar; §7 trägt "Abweichung gegen den Entwurf in §1: keine", das Lastenheft stimmt mit §1 überein; Verifier. <!-- d-check:ignore (Kennung als Zitat im Lauf-Beleg) -->

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Aufbau eines Satzes weiter benannt als gemessen · Start-Pfad als einziger Unterschied zweier Container-Begriffe · Aussage im Aufbau ohne Gegenstand · Vertrag enger als die Zusage-Menge der ADR · bekannte INFO, unverändert

## Verdikt

**Merge-blockierend:** nein. Kein HIGH, kein MEDIUM. L-1 ändert, falls behoben, Lastenheft-Text und braucht eine Freigabe des Projektinhabers (Eingrenzung des CLI-Containers auf die Egress-Restriktion oder Annahme mit Begründung in §7); nimmt der Projektinhaber den Befund an, ohne den Text zu ändern, genügt der Eintrag in §7. Ein Widerspruch des Implementers stuft den Befund nicht herab. Ein weiteres Nachlauf-Review ist nur nötig, wenn der Wortlaut sich ändert.

**Übergabe:** Findings gehen an den Implementer; L-1 zusätzlich an den Projektinhaber (Entscheidung über den Wortlaut). Verifier: Satz-für-Satz-Vergleich gegen §7 und Fassung 3, Z13 bis zum Code-Slice, Freigabe-Reihenfolge. Architect: INFO zu Z4 für die Spezifikation.

## Kennzeichnung der Findings

L-1 = LOW CLI-Container ohne Egress-Qualifier; I-1 = INFO zwei Container-Begriffe; I-2 = INFO Negation im Aufbau ohne Workspace-Volume; I-3 = INFO Z4 in der ADR; I-4 = INFO fortgeführte N-5 bis N-7.
