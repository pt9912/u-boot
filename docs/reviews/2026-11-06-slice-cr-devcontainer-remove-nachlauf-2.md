# Review-Report (Nachlauf 2): slice-cr-devcontainer-remove, Folgeänderung Lastenheft 0.37.2 — 2026-11-06

**Review-Art:** Nachlauf-Review zu einer Änderung des Wortlauts nach dem Nachlauf-Review (Beobachtung `BEO-DOK/vertragsaenderung-nach-review`, DoD-Zeile 5 des Slice). Vorlauf: [`2026-11-05-slice-cr-devcontainer-remove-nachlauf`](2026-11-05-slice-cr-devcontainer-remove-nachlauf.md).
**Gegenstand:** `git diff db8fa8af..7f568768` (HEAD), 7 Dateien, 111 Einfügungen, 38 Löschungen. Betroffen sind `spec/lastenheft.md`, der CR-Vorschlag [`cr-vorschlag-devcontainer-remove.md`](../plan/cr/cr-vorschlag-devcontainer-remove.md), der Slice [`slice-cr-devcontainer-remove`](../plan/planning/done/slice-cr-devcontainer-remove.md), [`open/README.md`](../plan/planning/open/README.md), die zwei Folge-Slices [`slice-v1-devcontainer-remove`](../plan/planning/done/slice-v1-devcontainer-remove.md) und [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) und die Welle [`welle-devcontainer-remove`](../plan/planning/done/welle-devcontainer-remove.md).
**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5.

## Eingangs-Kontext und Sensors

**Gelesen:**

- Slice (Diff in Wortform: Kopf, §1 Ausgänge G-1 und G-2, DoD, §3 Tabelle und Abgleich, §4, §6), CR-Vorschlag (Abschnitt „Folgeänderung 0.37.2“, Stellen (l) bis (o)), Vorlauf-Report.
- Im Lastenheft: [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (Zeilen 513 bis 524), [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) (Zeilen 205 bis 214), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), [`LH-NFA-PORT-004`](../../spec/lastenheft.md#lh-nfa-port-004--ausgeliefertes-container-image).
- Code und Messung als Belege der Aussagen im CR-Vorschlag: `internal/hexagon/application/downservice.go` (`Down`, `runConfirmationGate`), `Dockerfile` (Stage `runtime`, Zeilen 173 bis 187), `internal/adapter/driving/cli/devcontainer_shell_test.go` (Zeilen 309 bis 317), Messprotokoll [`slice-v1-devcontainer-remove-messung`](../messungen/slice-v1-devcontainer-remove-messung/README.md) M3, `spec/spezifikation.md` Zeile 948.

**Sensors:**

- `make docs-check`: gefahren auf HEAD `7f568768`, Arbeitsbaum sauber, `1168 Datei(en) geprüft, 0 Befund(e)`.
- `make gates`, `make test-docker`, `make verify-depguard`: nicht gefahren. Der Diff enthält nur Doku, keinen Code, keinen Docker-Pfad und keine Importregel.
- Eigene Proben, nur lesend:
  - Wortvergleich der Stellen (l), (m), (n) und der Historie-Zeile (o) des CR-Vorschlags gegen `spec/lastenheft.md` (Link-Form und Fettdruck normiert): 4 von 4 wortgleich.
  - `git diff -U0 db8fa8af HEAD -- spec/lastenheft.md`: Versionsfeld, Zeile 517 (ein Satz hinzu), Zeile 522 (eine Klausel), Zeile 523 (ein Halbsatz), Historie-Zeile 1664; sonst nichts.
  - Suche der Stellen nach festgehaltenen Zeilennummern des Slice (9, 517, 522, 523, 1664) gegen HEAD: stimmen.
  - Satzzählung des Spiegelstrichs Zeile 517: Der neue Satz ist Satz 10, die Abgrenzung wird Satz 11, wie die „Lesart“ des CR-Vorschlags sagt.

## Findings

### MEDIUM `spec/lastenheft.md` · „Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.“ — die Ausnahme der Version 0.37.2 widerspricht der Regel in LH-FA-CLI-005A, die den Befehl namentlich nennt

**Klasse:** Anforderung widerspricht einer unveränderten Schwesterstelle (Rang-1-interner Widerspruch)
**Quelle:** [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) (Exit-Code-Zuordnung bindend); Skill MEDIUM „Fehlerbehandlung … keiner Exit-Code-Kategorie eindeutig zugeordnet“
**Befund:**

- [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) Zeile 214 zählt `u-boot devcontainer remove --volumes` namentlich unter den destruktiven Operationen auf und sagt ohne Bedingung: „Im nicht-interaktiven Modus ohne `--yes` ist der Befehl mit Exit-Code `10` abzubrechen.“ Die Zeile steht im Diff nicht.
- Der neue Satz 10 in [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) lässt denselben Fall ohne Container mit Exit-Code `0` enden. Der Satz 4 des Spiegelstrichs wiederholt die unbedingte Fassung („im nicht-interaktiven Modus ohne `--yes` endet der Befehl mit Exit-Code `10`“), und Satz 5 ebenso; die Einschränkung steht erst in Satz 10, nach den Sätzen zu Instanz und Projekt.
- Der CR-Vorschlag prüft diesen Widerspruch nicht. Stelle (f3) hält fest, der Rest von Zeile 214 trage den Exit-Code `10` und die Zeile bleibe „sonst wörtlich“; der Abschnitt „Prüfung der Folgeänderung 0.37.2 gegen die Messgrenzen“ nennt [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) nicht. Der Slice-Plan führt die Zeile 214 in der Tabelle von §3 nur als gelesenen Stand der Übernahme `0.37.0`.
- Das Schwesterkommando in derselben Aufzählung verhält sich anders: `u-boot down --volumes` prüft die Freigabe vor jedem Engine-Zugriff (`downservice.go`, `Down`, Aufruf von `runConfirmationGate`; Zeile 63) und endet auch ohne Container im nicht-interaktiven Modus ohne `--yes` mit Exit-Code `10`. `remove --volumes` endet dort mit `0`. Der CR-Vorschlag nennt den Unterschied in der Prüftabelle; die Anforderung, die beide Befehle in einem Satz bindet, nennt ihn nicht.

**Risiko:** Der Welle-Plan verlangt eine Zusage-zu-Beleg-Tabelle, in der jeder Satz der Aufzählungen in [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) einen Test hat ([`welle-devcontainer-remove`](../plan/planning/done/welle-devcontainer-remove.md) Kriterium „Eine Aussage über `remove`“). Der Implementer pinnt Zeile 214 für `remove --volumes` als „ohne `--yes` nicht-interaktiv: `10`“ und widerspricht damit Fall (d) der DoD desselben Slice (ohne Container: `0`). Oder er liest Zeile 214 als maßgeblich und baut die Freigabe vor die Containersuche, wie `down` es tut; dann ist der Nachlauf-Review-Befund aus dem Vorlauf (Reihenfolge) wieder offen. Beide Lesarten sind mit dem Wortlaut vereinbar, der Vertrag nennt keine Vorrangregel. Welche Stelle gilt, entscheidet der Projektinhaber über den Architect, nicht der Reviewer. Jede Wortlaut-Änderung läuft danach über Freigabe und einen weiteren Nachlauf-Review.
**Verifizierbar:** Review-only (der Pin-Test des Umsetzungs-Slice hat bis zur Entscheidung kein eindeutiges Soll)

### LOW `docs/plan/planning/open/README.md` · „liegt in `in-progress/`, in Umsetzung, Wortlaut freigegeben und im Lastenheft `0.37.0` übernommen“ — die Zeile trägt die Versionskette als Chronik

**Klasse:** Zustandsfeld trägt Chronik (Ursache der wiederkehrenden Zustandsmarker-Drift in den Läufen davor)
**Quelle:** `AGENTS.md` §3.7 (Zustandsfelder), Skill LOW mit Konventions-Anker
**Befund:** Die Zeile zum Slice nennt für `0.37.0`, `0.37.1` und `0.37.2` je Freigabe, Übernahme und Nachlauf-Review („freigegeben und übernommen (Nachlauf-Review liegt vor)“, „weiterer Nachlauf-Review steht aus“). Das ist der Weg zum Zustand, nicht der Zustand. Dieselbe Kette steht in der Welle (Welle-Ziel, Reihenfolge, Freigaben) und im Kopf des Slice; jede Folgeänderung schreibt sie an vier Stellen fort.
**Risiko:** Die Stellen driften, wenn eine vergessen wird. Die Klasse „Zustandsmarker mit Ablaufpunkt nicht entfernt“ traf in den drei Läufen davor genau hier. In diesem Diff sind die Marker in allen sieben Dateien auf den Ist-Stand gebracht (siehe „Geprüft ohne Befund“); die nächste Folgeänderung (`0.37.3`) hat dieselbe Fortschreibung an vier Stellen zu leisten, und jeder bleibende Marker „Nachlauf-Review steht aus“ ist mit diesem Report überholt.
**Verifizierbar:** Review-only

### INFO `spec/lastenheft.md` · „und `u-boot devcontainer remove` läuft dort ebenfalls nicht“ — Wortlaut ohne Beleg im Test, Plan trägt den Auftrag

**Klasse:** Bedingung in der Spec ohne pinnenden Test (offen im Umsetzungs-Slice)
**Quelle:** `AGENTS.md` §3.10 (Bedingung in Spec oder Handbuch bekommt einen Test oder wird gestrichen)
**Befund:** Die Stufe `runtime` des `Dockerfile` (Zeilen 173 bis 187) ist `distroless/static` und kopiert nur das Binary; `docker` liegt dort nicht im `PATH`. Die Aussage ist damit plausibel, nach `spec/spezifikation.md` Zeile 948 endet der Befehl dann mit Exit-Code `11`. Ein Test, der sie pinnt, existiert nicht (der CR-Vorschlag sagt es selbst: nicht gelesen, nicht gemessen; für den Satz zu `up` ist ebenfalls kein Test gefunden). „läuft nicht“ nennt keinen Exit-Code. Den Auftrag tragen die DoD von [`slice-v1-devcontainer-remove`](../plan/planning/done/slice-v1-devcontainer-remove.md) (Zeile „Bedingungen des Lastenhefts `0.37.2`“, Punkt 1) und die Plan-Tabelle von [`slice-v2-spec-devcontainer-remove`](../plan/planning/done/slice-v2-spec-devcontainer-remove.md) (Zeile 948).
**Risiko:** keines im Diff; wird der Auftrag nicht erfüllt, steht im Lastenheft eine Bedingung ohne Beleg. Die Pflicht liegt beim Verifier des Umsetzungs-Slice.
**Verifizierbar:** Review-only (Test im Umsetzungs-Slice)

### INFO `spec/lastenheft.md` · „`--no-interactive` bricht bei jeder offenen Bestätigungsfrage mit Exit-Code `2` ab“ — Spannung zu Exit-Code `10` besteht weiter

**Klasse:** Vorbestehende Spannung, keine Änderung durch den Diff (F-7 des Vorlaufs des CR-Reviews)
**Quelle:** [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
**Befund:** Unverändert: Die Zeile mit Exit-Code `2` für `--no-interactive` bei offener Bestätigungsfrage steht neben der Zeile 214 mit Exit-Code `10`. Der Slice führt die Frage in §6 an den Architect. Sie hängt mit dem MEDIUM zusammen: Beide betreffen, welche Zeile von [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) für `remove --volumes` gilt.
**Risiko:** Keines aus dem Diff.
**Verifizierbar:** Review-only

## Auflösung der Findings des Vorlaufs

| Vorlauf | Stand im Diff |
|---|---|
| MEDIUM 1 (Zustandsmarker „Arbeitsbaum“, „nicht committet“, Version `0.37.0` in der Welle, `git show HEAD:`-Anweisung) | Aufgelöst. Die Anweisung nennt jetzt die Commits `602d8fa4` und `a7914c99`; beide tragen die genannte Version (Zeile 9 geprüft: `0.37.0`, `0.37.1`). Kein Treffer mehr für „Arbeitsbaum“ oder „nicht committet“ in den sieben Dateien. Die Marker „steht aus“ sind wahr und mit diesem Report überholt (siehe LOW). |
| MEDIUM 2 (Reihenfolge von Freigabe und Containersuche) | Zum Teil aufgelöst. Der Satz 10 und die Klausel in den Fehlerklassen entscheiden den Fall ohne Container; der Widerspruch zu Zeile 214 von [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) entsteht dadurch (MEDIUM oben). |
| LOW (Reichweite des Image-Satzes) | Aufgelöst im Wortlaut (Halbsatz „und `u-boot devcontainer remove` läuft dort ebenfalls nicht“). Der Test steht aus (INFO oben). |
| INFO (Spannung Exit-Code `2` und `10`) | offen, unverändert (INFO oben) |
| INFO (Machbarkeit des Falls (e) im Docker-Test) | Behandelt. Der Auftrag in der DoD des Umsetzungs-Slice (Punkt 2) nennt die Lücke und die gelesene Stelle; die Fundstelle `devcontainer_shell_test.go` Zeilen 309 bis 317 stimmt. Ob der dort genannte Docker-Test unter `script` existiert, bleibt nicht gelesen, und der Plan sagt das. |

## Geprüft ohne Befund

- **Wortgleichheit CR-Vorschlag und Lastenheft:** Stellen (l), (m), (n) und (o) 4 von 4 wortgleich. Der Lastenheft-Diff trägt nichts außerhalb dieser Stellen.
- **Satz 10 gegen Satz 7, 8 und 9 des Spiegelstrichs:** Die Klammer „(im Projekt oder, mit `--instance`, in der Instanz)“ trennt den Fall ohne Container vom Fall „ohne `--instance`, nur Container von Instanzen“ (Satz 8, Exit-Code `10` mit Hinweis). Die „Lesart“ im CR-Vorschlag benennt das; der Fehlerklassen-Satz („bei einem Container, den der Befehl entfernen würde“) stimmt dazu.
- **Fehlerklassen gegen Spiegelstrich:** Exit-Code `10` bei fehlender Freigabe nur mit Container, `12` für Stoppen und Entfernen, `11` für unerreichbare Engine: Die Containersuche vor der Freigabe verlangt die Engine, die Reihenfolge `11` vor `10` ist damit gegeben und widerspricht nichts.
- **Volumes eines nicht mehr vorhandenen Containers:** Der Nicht-Ziele-Spiegelstrich (Zeile 524) sagt, `u-boot` entfernt sie nicht; Exit-Code `0` ohne Container lässt das offene Nutzerpfad-Thema (`remove` ohne `--volumes`, dann `remove --volumes`) nicht neu entstehen.
- **Referenzrichtung (`AGENTS.md` §3.4, Skill HIGH #2):** Die hinzugefügten Lastenheft-Zeilen tragen keinen Verweis auf ADR, Slice, Carveout, Roadmap oder Dateien. Die Historie-Zeile nennt „Nachlauf-Review“ als Anlass in Worten, ohne Link.
- **Anforderungsebene des Wortlauts:** keine Labelnamen, Pfade, Docker-Befehle oder Flag-Interna außer der Bedienoberfläche (`--volumes`, `--yes`, `--instance`).
- **Aussagen des CR-Vorschlags über Code und Messung:** `Down` ruft `runConfirmationGate` vor `ComposeDown` (und vor jedem Engine-Aufruf in den Zeilen 1 bis 3 der Wahrheitstabelle) auf; die Messung M3 belegt `down --instance a` ohne Container mit Exit-Code `0`; der CR-Vorschlag behauptet davon nicht mehr.
- **Slice §3 Plan-Tabelle gegen `git diff --stat db8fa8af..HEAD`:** dieselben 7 Dateien, jede mit Zeile; die Zeilenangaben zum Lastenheft stimmen mit HEAD. Der Abgleichssatz nennt den Stand `HEAD` `db8fa8af` (Stand vor dem Commit der Folgeänderung) und ist damit ein benannter Stand, kein falscher Marker.
- **Zustandsmarker nach dem Commit:** Slice-Kopf, CR-Status, `open/README.md`, Welle (§1, §3, §4, §5) und DoD-Zeilen 3 bis 5 nennen für `0.37.2` den Ist-Stand (freigegeben, im Lastenheft; Nachlauf-Review offen).
- **Lifecycle (`AGENTS.md` §3.3):** keine Renames, der Slice bleibt in `in-progress/`, die Folge-Slices bleiben in `open/`.
- **Kommentar-Selbstcheck (`AGENTS.md` §3.9):** kein Code- oder Skript-Kommentar im Diff.
- **Accepted-ADR, Nutzer-Doku, CHANGELOG, Carveouts:** nicht berührt (`AGENTS.md` §3.5 und §5 Regeln 5 und 6).
- **Dual-Classifier, Sprachvertrag, Dateisicherheit, Importregeln:** nicht berührt; der Diff enthält keinen Code.

## Nicht geprüft

- **Code, Tests, Docker-Pfad:** im Diff nicht berührt. Die Sentinel- und Diagnostic-Abdeckung von `remove` ist Aufgabe des Review von [`slice-v1-devcontainer-remove`](../plan/planning/done/slice-v1-devcontainer-remove.md).
- **Ob der Projektinhaber den Wortlaut `0.37.2` gesehen und die Auslegung zu `--instance` bestätigt hat:** liegt im Sitzungsdialog und ist nicht belegbar; die Planungsstellen führen es als erfolgt.
- **Ob der Docker-Test unter `script` zu `devcontainer_shell_test.go` existiert und ob `internal/e2e` ein Terminal-Testmittel trägt:** nicht gelesen; der Plan führt es als offen.
- **`make gates`:** nicht gefahren (Review-Rolle, Doku-Diff).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Anforderung widerspricht einer unveränderten Schwesterstelle (Rang-1-interner Widerspruch, Zeile 214 von [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)) · Zustandsfeld trägt Chronik · Bedingung in der Spec ohne pinnenden Test (Image-Satz) · Vorbestehende Spannung (Exit-Code `2` und `10`).

## Verdikt

**Merge-blockierend:** ja, für das MEDIUM. Es ist vor der Closure des Slice zu beheben oder durch den Architect oder Projektinhaber mit Begründung zu entscheiden. Ein Widerspruch des Implementers stuft es nicht herab.

**Übergabe:**

- **Architect und Projektinhaber:** Das MEDIUM und das INFO zu Exit-Code `2` und `10` betreffen beide, welche Zeile von [`LH-FA-CLI-005A`](../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung) für `remove --volumes` ohne Container gilt. Jede Wortlaut-Änderung läuft über die Freigabe des Projektinhabers und einen weiteren Nachlauf-Review, nicht über eine stille Korrektur. Dieser Report ist der Nachlauf-Review zu `0.37.2` (DoD-Zeile 5 des Slice).
- **Implementer:** Nach diesem Report sind die Marker „Nachlauf-Review steht aus“ in CR-Vorschlag, Slice (Kopf, §1, §3, DoD-Zeile 5, §6), `open/README.md` und Welle auf den Ist-Stand zu setzen; der Abgleich der Plan-Tabelle gegen `git diff --stat` ist danach zu wiederholen (`AGENTS.md` §3.10).
- **Steering-Loop:** Die Klasse „Zustandsmarker mit Ablaufpunkt nicht entfernt“ hat in diesem Lauf nicht getroffen; der Schwellen-Übertritt des Vorlaufs (Schärfung von Skill und `harness/review.md`) ist im Diff nicht angegangen und bleibt beim Architect.
- **Verifier:** Der Image-Satz (INFO) und die Pin-Tests der Fälle (d) und (e) gehören in den Abgleich des Umsetzungs-Slice. Dieser Report belegt nur den Nachlauf, nicht die Closure.
