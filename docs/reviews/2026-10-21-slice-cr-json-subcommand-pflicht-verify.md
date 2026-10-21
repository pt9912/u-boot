# Verification-Report: [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md) — 2026-10-21

**Rolle:** Verifier (frischer Kontext, kein Self-Review, repariert nichts). Frage: Wurde der CR richtig gebaut, gemessen an Plan, DoD, den Entscheidungen des Projektinhabers und den Sensoren? Form nach [`harness/verification.md`](../../harness/verification.md). Außer dieser Datei ist nichts geändert, nichts committet, nichts gepusht. Jede Mutation ist zurückgesetzt, `git status` war danach jeweils leer.

**Gegenstand:** `git diff 95d43213..ce3c61a3` ohne `cecf6876`. Die Implementer-Commits sind `f68a5473`, `9ac26881`, `2f2b5b71`, `ee0f0c8f`, `e5be6620`, `757ade7b`, `9d87a200`, `6bfec81c`, `e4049fd0` und die Fix-Runde `ce3c61a3`. `HEAD` = `ce3c61a3`, der Arbeitsbaum war vor dieser Datei sauber. `cecf6876` (Hauptagent, [`slice-v1-release-cut-v0.13.0`](../plan/planning/done/slice-v1-release-cut-v0.13.0.md)) habe ich gelesen: Er fügt in dessen Slice-Datei unter `done/` genau eine Zeile hinzu (Vermerk M7) und berührt sonst nichts.

**Eingang (gelesen, aber nicht als Beweis genommen):**

- Slice §1 bis §8 samt *Umsetzung* und *Fix-Runde*.
- Review-Report [`2026-10-21-slice-cr-json-subcommand-pflicht`](2026-10-21-slice-cr-json-subcommand-pflicht.md): 0 HIGH, 0 MEDIUM, 2 LOW, 5 INFO, GO.
- Entscheidungen des Projektinhabers vom 2026-10-21 (gesetzt): Vertrag schließen; Wortlaut der Stellen 1 bis 3 und S1 bis S5 wie im Plan; Variante A; die Beispielsätze in [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) bleiben; Handbuch-Zeile `subcommand` mitnehmen; Register-Nachtrag bei der Closure.
- Quellen: [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) und §16 im Lastenheft; [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) und §7 in der Spezifikation; `jsontestutil/`; [`cli-json-contract.md`](../maintainer/cli-json-contract.md); die Handbuch-Zeile `subcommand`.

**Gesamturteil: GO für die Closure.** Ich habe keine DoD-Verletzung gefunden. Die sechs Liefer-Punkte der DoD sind erfüllt; an zwei davon hängt eine Ablage-Auflage für §7. Die übrigen fünf sind Closure-Punkte und bis dahin offen. Die vier Befunde V-1 bis V-4 sind INFO; keiner blockiert.

## Verification Evidence

Scope:

- Slice: [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md)
- IDs: [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe). Berührt sind nur die Befehle aus [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (`egress list`) und [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`devcontainer up`/`down`), deren Text unverändert bleibt. Keine ADR ist berührt.
- Artefakte: `spec/lastenheft.md`, `spec/spezifikation.md`, `jsontestutil.go`, `jsontestutil_test.go`, `jsonenvelope_test.go` (nur Kommentare), `docs/maintainer/cli-json-contract.md`, `docs/user/benutzerhandbuch.md`, Roadmap, `open/README.md`, Messstände-CR, Slice, Review-Report.

### DoD-Abgleich (§2 des Slice, Stand `ce3c61a3`)

- [x] **Lastenheft** — erfüllt. Den Wortlaut habe ich zeichengenau gegen den Plan verglichen (siehe *Einzelnachweise*). Version `0.19.0` steht in Zeile 9; die Historie-Zeile `0.19.0` (Zeile 1589) trägt das Datum `2026-10-21`, das ist das Commit-Datum von `e5be6620`, dazu „Vereinbarung mit dem Projektinhaber“. Es gibt keinen Slice-, ADR- oder Prozessverweis. `0.19.0` war beim Commit die nächste freie Version: Der Elter `e5be6620^` trägt `0.18.0` im Versionsfeld und als letzte Historie-Zeile, und `0.19.0` kommt vorher im Lastenheft nicht vor. **Für die Closure:** Die DoD verlangt die Fundstellen „in §7“, sie stehen heute in §3 *Umsetzung*; ein Verweis aus §7 genügt (V-2).
- [x] **Spezifikation S1 bis S5** — erfüllt, zeichengleich. Die Bedingung in S4 („erzeugen keine Vorschau-Ausgabe“) pinnen `TestExecute_EgressList_RejectsPreviewFlagsAndArguments` und `TestExecute_Devcontainer_RejectsPreviewFlagsAndArguments`. Beide habe ich gelesen: Sie prüfen `--dry-run` und `--diff` je mit `--json`, erwarten Exit 2 und laufen über `AssertMinimalEnvelope`, das `checkNoFullFields` aufruft.
- [x] **Wächter** — erfüllt. `checkSubcommand` prüft per `slices.Contains` die vier Befehle; die Fehlermeldung ist unverändert, und `checkSubcommand` läuft in beiden Helfern (Zeilen 142 und 172). Der Wächter-Test deckt je Befehl drei Fälle ab: `missing`, `empty` und `set` (ohne Befund). Dazu kommen die Gegenprobe `doctor/missing` und `TestAssertFullEnvelope_RequiresSubcommandForGroupedCommands`. Die Mutationen M1 bis M9 stehen mit Ergebnis in §3; M2, M4 vor/nach und M7 habe ich selbst nachgefahren, mit demselben Ergebnis (siehe *Mutationen*). Die drei Kommentare in `jsontestutil.go` und die zwei in `jsonenvelope_test.go` nennen die vier Befehle; die doppelte Kennung im Doc-Kommentar von `AssertMinimalEnvelope` ist berichtigt.
- [x] **Freigabe** — erfüllt. §7 vermerkt die fünf Entscheidungen F-1 bis F-5; sie decken sich mit den Entscheidungen im Auftrag. Das Gespräch selbst kann ich nicht nachprüfen; Beleg sind der Vermerk und der Auftrag des Hauptagenten.
- [x] **Vorab-Suche** — inhaltlich erfüllt. Die Suche lief vor dem Lastenheft-Commit (`ee0f0c8f`) und nach den Commits (`6bfec81c`); jeder Treffer steht mit Befund in §3 *Umsetzung*. Meine eigene Suche (siehe *Einzelnachweise*) fand keine weitere Kopie der Regel. **Für die Closure:** Die DoD sagt „in §7“ (V-2).
- [x] **`make gates`, `make docs-check`, `make doc-trace`** — selbst gefahren auf `ce3c61a3`, alle mit Exit 0 (siehe *Sensors*). Die Grenze von `doc-trace` ist richtig benannt (V-4). **Für die Closure:** Läufe und Exit-Codes „in §7“; heute stehen sie in §3, der Lauf nach den Closure-Commits fehlt noch.
- [x] **Review** — erfüllt; der Report liegt seit `ce3c61a3` im Repo. Ein Nachlauf-Review ist nicht fällig, denn die Fix-Runde berührt weder Lastenheft noch Spezifikation (`git show --stat ce3c61a3`: Slice, `open/README.md`, Review-Report).
- [x] **Doku-Update** — erfüllt. K1 bis K5 in `cli-json-contract.md` und die Handbuch-Zeile sind gegen Quelle und Code geprüft (siehe *Einzelnachweise*). `CHANGELOG.md` liegt nicht im Diff.
- [ ] **Closure-Notiz mit Steering-Loop-Lerneintrag** — offen (Closure).
- [ ] **Beobachtungs-Register** — offen (Closure). Für F-5 steht der Nachtrag in der Evidence-Datei des Vorgangs [`slice-v2-spec-devcontainer-start`](../plan/planning/done/slice-v2-spec-devcontainer-start.md) unter [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../plan/planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) aus. Die Datei existiert; `observations/` liegt nicht im Diff.
- [ ] **Jedes Risiko aus §6 mit Ausgang** — offen (Closure). Eines ist bereits „entfallen“, sechs stehen auf „weiter offen“; Eingaben dafür in der *Prüfliste*.
- [ ] **Drei Paarungen** — offen (Closure).

### Einzelnachweise (selbst gelesen oder gemessen, nicht aus dem Slice übernommen)

**Wortlaut, zeichengenau.** Ich habe die Zeilen der Abschnitte *Wortlaut-Vorschlag* im Plan per Skript mit den Zeilen in `spec/` verglichen. Das Skript lag im Scratchpad und war ein Text-Werkzeug, kein Gate. Vor dem Vergleich habe ich das Zitatzeichen entfernt, den Pfad `../../../../spec/<datei>#` auf `#` gesetzt und den Platzhalter `<Datum…>` durch `2026-10-21` ersetzt.

Zuerst die Vorfrage, ob der Vorschlag im Plan seit der Vorlage geändert wurde. Zwischen Plan-Commit `95d43213` und `ce3c61a3` unterscheiden sich die Abschnitte nur in drei Überschriften und zwei Vermerken (Freigabe, „nicht gewählt“); der Text der Vorschläge ist gleich.

| Stelle | Ergebnis |
|---|---|
| Stelle 1, Spiegelstrich `subcommand` | identisch mit Lastenheft Zeile 974 |
| Stelle 2 (Variante A) | Zeile 978 endet identisch mit dem Neu-Text. Die alte Zeile mit dem Ersatz „`; bei \`template\` und \`config\` ist \`subcommand\` verpflichtend, und`“ → „`, und`“ ist gleich der neuen Zeile; sonst ist der Absatz unverändert. Der alte Teilsatz kommt im Lastenheft nicht mehr vor. |
| Stelle 3, Versionsfeld | identisch mit Zeile 9, Spaltenbreite gleich |
| Stelle 3, Historie | identisch mit Zeile 1589; beide Halbsätze stehen da (Variante A) |
| Übriges Lastenheft | Der Diff hat genau drei Hunks. [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) Zeile 250 und [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) Zeile 1117 sind unverändert (Entscheidung F-2). |
| S1 | identisch mit Spezifikation Zeile 276 |
| S2 | identisch mit Zeile 822 |
| S3 | Der Block mit 11 Zeilen ist identisch ab Zeile 907, zwischen `config` und `devcontainer`. Der Schema-Block ist gültiges JSON (geparst); `allOf` hat genau die Werte `template`, `config`, `egress` und `devcontainer`. |
| S4 | identisch mit Zeile 938, direkt nach dem Absatz zum Wert `devcontainer` (Zeile 936) |
| S5 | identisch mit Zeile 1363, wenn [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) ohne Link steht. Der Plan sagt genau das zu: „führt die Kennung der Schema-Stelle wie ihre Nachbarzeilen ohne Link“. |

**Referenzrichtung** ([`AGENTS.md`](../../AGENTS.md) §3.4): Unter den hinzugefügten Zeilen in `spec/` trifft keine `slice-`, `welle-`, `ADR-`, `adr/`, `docs/plan`, `BEO-`, `roadmap`, `carveout` oder `tranche`. Die neuen Verweise zeigen nur auf Anker der eigenen Datei.

**Kopien in `cli-json-contract.md`:**

- K2: Die Schema-Kopie in §3 (117 Zeilen, ab Zeile 140) ist zeilengleich mit dem Schema in der Spezifikation (ab Zeile 815).
- K1: Beide Zitatblöcke in §2 habe ich inhaltlich gegen ihre Quelle verglichen, ohne Zeilenumbruch und ohne die Link-Form `[X]` gegen `` [`X`] ``. Ergebnis: Lastenheft-Zitat = Lastenheft Zeilen 970 bis 976; Spezifikations-Zitat = [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) Zeilen 270 bis 278. Beide Quellenangaben stimmen.
- K3 bis K5 habe ich gelesen: Sie nennen das Lastenheft als Quelle der Pflicht und `checkSubcommand` als Prüfung. Die Aussage „das Enum prüft `jsontestutil` nicht“ ist wahr (`checkCommand` vergleicht nur mit `WithCommand`).
- §6.9 („Pflicht-`subcommand`“, Zeile 1082) bleibt wahr.

**Handbuch-Zeile `subcommand` gegen den Code.** Ich habe jede Stelle gelesen, an der ein Umschlag entsteht: alle Aufrufe von `newMinimalEnvelope`, `newDataEnvelope`, `newFullEnvelope`, `reportErrorSub`, `jsonArgsValidator`, `configArgsValidator`, `writeConfigEntries` und `reportDevcontainerError` in den Nicht-Test-Dateien des Pakets `cli`.

| Befehl | Werte im Code | Pinnender Test |
|---|---|---|
| `config` | `show` (`newConfigCommand`: `configArgsValidator(a, "show", …)`, `runConfigShow`), `get`, `set`, `list` | `TestConfigJSON_ShowMinimalDataEnvelope` und `TestConfigJSON_SubcommandAlwaysSet` (show, get, set), `config_many_test.go` Zeile 140 (list) |
| `egress` | `list` | `WithSubcommand("list")` in `egress_test.go` |
| `devcontainer` | `up`, `down` | `WithSubcommand(sub)` in `devcontainer_test.go` |
| `template` | `list` | `template_test.go` Zeile 99 |

Alle übrigen Befehle übergeben `""`: `doctor`, `add`, `init`, `remove`, `generate`, `logs`, `up`, `down` sowie `reportError` mit `""` für diese sieben. Weil das Feld `omitempty` ist, fehlt es dort.

Der neue Wert `show` ist wahr und gepinnt. Die Zeile trägt keine Kennung ([`AGENTS.md`](../../AGENTS.md) §5 Regel 6). Dass `subcommand` „nur bei“ den vier Befehlen steht, pinnt kein Test je Befehl (V-1).

**Eigene Vorab-Suche** (Host-`git grep` als Text-Werkzeug). Gesucht habe ich in `spec`, `docs/user`, `docs/maintainer`, `docs/plan/adr`, `harness`, `AGENTS.md`, den READMEs, `CHANGELOG.md`, `internal`, `cmd` und `.harness/skills`.

- Muster `subcommand` zusammen mit Pflicht, verpflichtend, required, mandatory, gruppiert oder grouped: Jeder Treffer steht im Plan, als unverändert gelesen oder ist geändert. Dazu zählen die Chronik in `CHANGELOG.md` `[0.4.0]` (Zeilen 514 und 874), die Kommentare in `template.go` und `erroremission.go` sowie die Sentinels in `egress.go` und `devcontainer.go`.
- Die alten Aufzählungen „template/config“, „`template` oder/und `config`“ und „{template, config}“ kommen außerhalb des Slice nicht mehr vor.
- Der ersetzte Testname `TestAssertMinimalEnvelope_RejectsBadSubcommandForTemplate` steht nur noch im Slice ([`AGENTS.md`](../../AGENTS.md) §3.9 (b)).

**Kommentare** ([`AGENTS.md`](../../AGENTS.md) §3.7 und §3.9): Die neuen Kommentare in `jsontestutil.go`, `jsontestutil_test.go` und `jsonenvelope_test.go` sind Zusagen im Indikativ, ohne Chronik. Keine neue Funktion steht zwischen einem Doc-Kommentar und seiner Funktion. `slices` war schon vorher importiert (`95d43213`); Go `1.26.0` (`go.mod`), damit ist die Schleifenvariable in `t.Run` pro Durchlauf gebunden.

**Commit-Messages:**

| Commits | Slice-Kennung | [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) |
|---|---|---|
| Alle elf Slice-Commits (`95d43213` bis `ce3c61a3` ohne `cecf6876`) | ja | — |
| alle außer den Verschiebe-Commits `9ac26881` und `2f2b5b71` | — | ja |

Kein Subject und kein Body trägt `SPEC-<NNN>` oder `ARC-<NNN>`. `make doc-commits` meldet keinen Befund. Den Befund F-7 des Reviewers (Verfeinerungs-Kennung im Body von `6bfec81c`) teile ich als INFO.

**Lifecycle** ([`AGENTS.md`](../../AGENTS.md) §3.3): `9ac26881` ist ein reiner Rename (0 Zeilen geändert). Die Datei liegt nur in `in-progress/`; `in-progress/` enthält außer den drei Master-Dokumenten nur diesen Slice (WIP 1). Die Roadmap-Zeile „In Arbeit: …“ ist gesetzt.

### Mutationen des Verifiers (je `make test` im Docker-Lauf, danach `git checkout`, `git status` leer)

| # | Mutation | Ergebnis | Ausgabe (Auszug) |
|---|---|---|---|
| M2 | `"devcontainer"` aus der Menge in `checkSubcommand` | **rot, Exit 2**. Nur die Fälle `devcontainer` im Paket `jsontestutil`; `internal/adapter/driving/cli` ist `ok` | `--- FAIL: TestAssertMinimalEnvelope_RequiresSubcommandForGroupedCommands/devcontainer/missing` und `/devcontainer/empty` · `jsontestutil_test.go:153: reported "subcommand required when command == \"devcontainer\"" = false, want true` |
| M7 | Aufruf `checkSubcommand` aus `AssertFullEnvelope` entfernt (Zeile 172) | **rot, Exit 2**. Nur `TestAssertFullEnvelope_RequiresSubcommandForGroupedCommands`; `cli` ist `ok` | `jsontestutil_test.go:183: expected subcommand-required reject for command=config, errors=[]` |
| M4 vor | `jsontestutil.go` und `jsontestutil_test.go` im Stand `757ade7b`, dazu `newDataEnvelope("devcontainer", "", …)` in `runDevcontainerDown` und `WithSubcommand("down")` aus `TestExecute_DevcontainerDown` entfernt | **grün, Exit 0**: Das ist die Lücke, nichts fängt es | `ok … internal/adapter/driving/cli` |
| M4 nach | dieselbe Mutation auf `ce3c61a3` | **rot, Exit 2**. Fehlschläge: `TestExecute_DevcontainerDown/stopped`, `/already_stopped` und `/no_container` | `devcontainer_test.go:374: subcommand required when command == "devcontainer"`, dahinter die Kennung in Klammern |
| M10 (eigene) | `doctor` setzt `subcommand: "x"` (`doctor.go`: `newMinimalEnvelope("doctor", "x", …)`) | **grün, Exit 0** | Kein Test pinnt, dass `subcommand` „nur bei“ den vier Befehlen steht (V-1) |

Der Reviewer hat M1 und M3 nach gefahren, ich M2, M4 vor/nach, M7 und M10. Damit ist der Lückennachweis für beide neuen Befehle selbst gesehen: M3 nach (Reviewer) und M4 vor/nach (Verifier). Nicht selbst gefahren sind M3 vor, M5, M6, M8 und M9 (siehe *Nicht ausgeführt*).

### Plan-Tabelle gegen `git diff --stat 95d43213..ce3c61a3` ([`AGENTS.md`](../../AGENTS.md) §3.10, nach der letzten Fix-Runde)

Der Diff umfasst 13 Dateien. Eine davon kommt aus dem fremden Commit `cecf6876`: die Datei von [`slice-v1-release-cut-v0.13.0`](../plan/planning/done/slice-v1-release-cut-v0.13.0.md), eine Zeile. Die übrigen 12 haben je eine Zeile in der Tabelle:

- `cli-json-contract.md`
- `roadmap.md`
- der Slice (Rename `open/` → `in-progress/`)
- `open/README.md`
- die Messstände-CR [`slice-cr-lastenheft-messstaende`](../plan/planning/open/slice-cr-lastenheft-messstaende.md)
- der Review-Report (Zeile `docs/reviews/`)
- `benutzerhandbuch.md`
- `jsonenvelope_test.go`
- `jsontestutil.go`
- `jsontestutil_test.go`
- `lastenheft.md`
- `spezifikation.md`

Umgekehrt habe ich `git diff --stat 95d43213..ce3c61a3` über die als „unverändert“ geführten Pfade laufen lassen. Geprüft waren `spec/architecture.md`, `CHANGELOG.md`, beide READMEs, `docs/plan/adr`, `egress.go`, `devcontainer.go`, `template.go`, `config.go`, `jsonargs.go`, `erroremission.go`, `root.go`, `cli.go`, `rootjson_test.go`, `egress_test.go`, `devcontainer_test.go`, `config_acceptance_test.go`, `template_acceptance_test.go`, `drift_test.go`, `internal/e2e`, `benutzerhandbuch-standard.md` und `carveouts.md`. Die Ausgabe ist leer.

Ohne Diff sind drei Einträge, alle begründet: dieser Report, das Register (Closure) und die Zeile `docs/reviews/`, deren Text „der Verification-Report folgt“ bei der Closure auf den Ist-Stand geht. Kein Plan-Eintrag ist ohne Grund leer, keine Datei im Diff ist ohne Eintrag. Mit den Closure-Commits ist der Abgleich zu wiederholen.

### Fix-Runde `ce3c61a3`

- **F-1 (LOW) — behoben.** Die Zeile des Slice in `open/README.md` (Zeile 49) nennt den Ist-Stand: „geschrieben als Lastenheft `0.19.0`“, „Wortlaut vom Projektinhaber am 2026-10-21 freigegeben (Variante A)“, „in Arbeit (liegt in `in-progress/`)“. Das stimmt mit Zeile 33 überein („nach 0.19.0“). Kein offenes Slice in `open/` oder `next/` nennt sonst `0.18.0` als nächste freie Version; der Colima-Slice nennt keine Lastenheft-Version.
- **F-2 (LOW) — behoben.** Die Aussage in §3 *Umsetzung* ist wahr für die aufgezählten Commits (eigene Prüfung oben). Eine Randnotiz dazu steht in V-3.
- **Die Berichtigung über den Auftrag hinaus ist zulässig.** Gemeint ist die Änderung im Titel („Lastenheft, nächste freie Version, heute 0.19.0“ → „Lastenheft 0.19.0“) und in der DoD-Zeile *Lastenheft* („heute `0.19.0`“ → „`0.19.0`, belegt durch `e5be6620`“). Sie ist eine zulässige Berichtigung eines abgelaufenen Zustandsmarkers, keine Änderung des Kriteriums. Begründung:
  1. Das Kriterium lautet vorher und nachher „nächste freie Minor-Version zum Zeitpunkt des Commits“. Geändert ist nur die Klammer mit dem erwarteten Wert.
  2. „heute `0.19.0`“ ist ein Zustandsmarker. Er lief mit dem Lastenheft-Commit ab, und [`AGENTS.md`](../../AGENTS.md) §3.10 verlangt, ihn dann „zu entfernen oder auf den Ist-Stand zu setzen“.
  3. Den Ist-Stand habe ich selbst bestätigt: Der Elter von `e5be6620` trägt `0.18.0`, `0.19.0` kommt vorher nicht vor, und die Messstände-CR hat das Lastenheft nicht berührt.
  4. Die Änderung steht offen in der Fix-Runde („das Kriterium bleibt“). Kein Haken ist gesetzt, keine Bedingung gelockert.
  5. Lastenheft und Spezifikation sind unberührt, also ist kein Nachlauf-Review fällig.
  6. Die Form „Wert, belegt durch `<Hash>`“ entspricht der Praxis für DoD-Zeilen in `done/`.

  Formal gehört die DoD dem Planner. Er bestätigt die Klammer bei der Closure mit einem Satz in §7; eine Rückführung ist nicht nötig.
- F-3, F-4 und F-6 sind als Hinweise an den Projektinhaber in §3 *Fix-Runde* und §7 *Vorgemerkt* geführt. F-5 und F-7 sind „keine Änderung“; das teile ich.

### Konformität mit den Entscheidungen

| Entscheidung | Befund |
|---|---|
| Vertrag schließen, kein Produktcode | Im Diff unter `internal/` liegen nur `jsontestutil.go` (Test-Helfer), `jsontestutil_test.go` und `jsonenvelope_test.go` (Kommentare). Keine Produktdatei ist im Diff. |
| Wortlaut Stellen 1 bis 3 und S1 bis S5 | zeichengleich (oben) |
| Variante A | Die Pflicht steht nur im Spiegelstrich; der Teilsatz ist weg; die Historie trägt beide Halbsätze. |
| Beispielsätze bleiben | Zeilen 250 und 1117 sind unverändert. |
| Handbuch-Zeile mitnehmen | `show` ergänzt, nur mit Werten aus dem Code, ohne Kennung |
| Register-Nachtrag bei der Closure | `observations/` liegt nicht im Diff; offen für die Closure |

### Sensors

Alle Läufe habe ich selbst auf `ce3c61a3` gefahren. Umgebung: Docker Engine 29.9.2, Host `linux/amd64`. Angegeben ist der Exit-Code des `make`-Prozesses.

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` (sauberer Baum, vor den Mutationen) | pass, Exit 0, 46 s | `lint`: `0 issues.` · `coverage-gate: OK — coverage 91.70% meets threshold 90%` · `d-check: 605 Datei(en) geprüft, 0 Befund(e)` · `baseline-verify: v6.18.1 OK — 54 Dateien` · `ablage-pruefung: OK`, Wächter-Test `26 Faelle, 0 Fehlschlaege` · Stempel `gates-passed.head` = `ce3c61a3…` |
| `make docs-check` | pass, Exit 0 | `605 Datei(en) geprüft, 0 Befund(e)`; Lauf mit dieser Datei unter *Abschluss* |
| `make doc-commits RANGE="95d43213^..ce3c61a3"` | pass, Exit 0 | `12 Commit(s) — OK`, 0 Befunde (`cecf6876` eingeschlossen) |
| `make doc-immutable RANGE="v0.13.0..ce3c61a3"` | pass, Exit 0 | `16 Commit(s) — OK`, 0 Befunde (Accepted-ADRs unverändert) |
| `make doc-trace` | pass, Exit 0 | `84 Anforderung(en), 0 Waise(n).` Grenze: V-4 |
| `make doc-planning` | pass, Exit 0 | `605 Datei(en) geprüft, 0 Befund(e)` |
| `make test` mit M2 | fail (erwartet), Exit 2 | siehe *Mutationen* |
| `make test` mit M7 | fail (erwartet), Exit 2 | siehe *Mutationen* |
| `make test` mit M4 vor | pass (erwartet, die Lücke), Exit 0 | siehe *Mutationen* |
| `make test` mit M4 nach | fail (erwartet), Exit 2 | siehe *Mutationen* |
| `make test` mit M10 | pass, Exit 0 | Befund V-1 |
| `make docs-check` und `make gates` mit dieser Datei | pass, je Exit 0 | siehe *Abschluss* |

Traceability:

| ID | Beleg |
| --- | --- |
| [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Spiegelstrich `subcommand`) | Wächter `checkSubcommand` in beiden Helfern. Er wird belegt durch `TestAssertMinimalEnvelope_RequiresSubcommandForGroupedCommands` (vier Befehle × `missing`/`empty`/`set` und `doctor/missing`) und `TestAssertFullEnvelope_RequiresSubcommandForGroupedCommands`. Gemessen mit M1, M3 nach (Reviewer) sowie M2, M4 vor/nach und M7 (Verifier). Im Produkt pinnen die Einzel-Pins `WithSubcommand(…)` in `egress_test.go`, `devcontainer_test.go`, `config_acceptance_test.go` und `template_test.go`. `make gates` ist grün. |
| [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) | S1 bis S5 zeichengleich; Schema gültiges JSON; `allOf` hat die vier Werte; die Bedingung in S4 pinnen die zwei Reject-Tests. Für die Kopie in `cli-json-contract.md` gibt es kein Gate, nur meinen Zeilenvergleich (Lerneintrag, Plan §1). |
| [`LH-FA-DEV-008`](../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion), [`LH-FA-DEV-010`](../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) | Text unverändert. Die Werte `list` bzw. `up`/`down` habe ich an `egress.go` und `devcontainer.go` gelesen; sie sind durch die Einzel-Pins gepinnt. |
| [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) | unverändert (Entscheidung F-2), gelesen |

Carveouts:

- Neu: none
- Gelöst: none
- Unverändert: alle. `carveouts.md` liegt nicht im Diff. Kein Gate und keine Schwelle ist gelockert; der Wächter ist strenger.

Nicht ausgeführt:

- `make test-docker`: Der Diff berührt keinen Docker-Pfad und keinen Produktcode. Kein Test mit Build-Tag `docker` ruft `jsontestutil` (Implementer und Reviewer haben das gesucht). `internal/e2e/devcontainer_start_docker_test.go` liegt nicht im Diff.
- `make verify-depguard`: Keine Importregel und kein Import ist berührt; `slices` war schon importiert.
- `make ci`, `make govulncheck`, `make image-scan`: Keine Abhängigkeit, kein `Dockerfile`, kein Workflow und kein Release-Pfad ist berührt. Nach der Tabelle in `verification.md` sind sie für diese Änderungsklasse kein Closure-Sensor.
- Mutationen M3 vor, M5, M6, M8 und M9: nicht selbst gefahren. M3 vor ist symmetrisch zu M4 vor, die übrigen habe ich am Code auf Plausibilität geprüft. Dabei habe ich Folgendes am Code gelesen: M6 trifft den Zweig `sub == ""` und M9 den Zweig `!hasSub`. Für M8 sind die `doctor`-Tests über `AssertMinimalEnvelope` und die Gegenprobe betroffen.
- Validierung gegen den Bedarf der Konsumenten: Aufgabe des Validators.

Commit / Artefakt:

- Geprüft auf `ce3c61a3`; kein Release-Artefakt.

## Review-Findings einzeln

| Finding | Schwere | Stand | Beleg des Verifiers |
| --- | --- | --- | --- |
| F-1 Zeile in `open/README.md` | LOW | **behoben** | Zeile 49 nennt den Ist-Stand und stimmt mit Zeile 33 überein (oben). |
| F-2 Aussage über Commit-Messages | LOW | **behoben** | eigene Zählung je Commit (oben); Randnotiz V-3 |
| F-3 Kommentare zu bare-Rejects | INFO | **Hinweis an den Projektinhaber** (§7 *Vorgemerkt*) | Kein Symbol davon liegt im Diff (Liste der unveränderten Dateien oben). |
| F-4 doppelte Kennungen, Chronik in `AssertFullEnvelope` | INFO | **Hinweis an den Projektinhaber** | Am Code bestätigt: Der Doc-Kommentar von `AssertFullEnvelope` sagt „In diesem Slice (Doctor) … NICHT verwendet“; `checkNoFullFields` nennt die Kennung doppelt. |
| F-5 `show` auch im Fehler-Umschlag | INFO | **keine Änderung** | Die Zeile behauptet keine Ausschließlichkeit für `show`. |
| F-6 Argumentfehler ohne Umschlag | INFO | **Hinweis an den Projektinhaber** | Außerhalb des Diffs; ohne Umschlag gibt es auch keinen ohne `subcommand`. |
| F-7 Verfeinerungs-Kennung in Commit-Message | INFO | **keine Änderung** | `make doc-commits` ist grün; nach §5 Regel 1 nicht eindeutig. |

## Befunde und Diskrepanzen des Verifiers

| ID | Klasse | Befund | Wirkung |
|---|---|---|---|
| V-1 | INFO | Die Handbuch-Zeile sagt, `subcommand` erscheine „nur bei“ `config`, `egress`, `devcontainer` und `template`. Das ist am Code wahr: Alle anderen Aufrufer übergeben `""`, und das Feld ist `omitempty`, gepinnt durch `TestMinimalEnvelope_SubcommandOmittedWhenEmpty` über einen Test-Export. Je Befehl gepinnt ist es aber nicht: M10 (`doctor` setzt `"x"`) bleibt grün. Die Bedingung „nur bei“ steht so seit `e363b3ba` (2026-10-02), älter als die Pin-Regel in [`AGENTS.md`](../../AGENTS.md) §3.10 (`d99b5fff`, 2026-10-18). Dieser Slice hat nur den gepinnten Wert `show` ergänzt. | Kein DoD-Bruch: Die DoD verlangt die Zeile „nach F-4“, und der neue Wert ist gepinnt. Der Planner entscheidet bei der Closure: Hinweis an den Projektinhaber, oder ein Folge-Pin als Test im Paket `cli` über die acht übrigen Befehle, der `subcommand` als abwesend prüft. Der Pin gehört **nicht** in `checkSubcommand`, denn das Lastenheft lässt `subcommand` bei den übrigen Befehlen „optional“. |
| V-2 | INFO | Drei DoD-Punkte verlangen ihren Beleg „in §7“: die Fundstellen des Lastenhefts, die Treffer der Vorab-Suche sowie Läufe und Exit-Codes. Heute stehen alle drei in §3 *Umsetzung*; §7 trägt nur die Entscheidungen und die Vormerkung. | Closure: §7 verweist auf §3 *Umsetzung* und *Fix-Runde* (ein Satz) und trägt den Gate-Lauf nach den Closure-Commits. |
| V-3 | INFO | F-2 ist behoben: Der Satz „Jede Message trägt die Kennung“ (die Slice-Kennung) „…“ ist wahr für die aufgezählten Commits. Der genannte Lesebefehl `git log --format=%B 95d43213^..HEAD` umfasst beim Schreiben aber auch `cecf6876`, der keine der beiden Kennungen trägt. | Keine Änderungspflicht. Wird der Satz bei der Closure angefasst, die Ausnahme `cecf6876` nennen oder die Range auf die Slice-Commits begrenzen. |
| V-4 | INFO | „`make doc-trace` ohne neue Waise“ ist erfüllt (84/0), sagt aber nichts über [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe): Die Belegsicht führt nur `LH-FA-*` (80) und `LH-QA-*` (4) (Ausgabe gelesen). Der Implementer nennt die Grenze in §3 richtig. | Eingabe für den Lerneintrag: Die DoD-Zeile ist für Anforderungen der Familie `LH-NFA-*` kein Beleg. |

Ich habe keine DoD-Verletzung gefunden und keine Behauptung ohne Beleg in den geprüften Vertrags- und Nutzer-Aussagen.

## Prüfliste für die Closure (Planner / Implementer)

1. **Diese Datei committen.** Die Nachricht trägt die Slice-Kennung und [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe). Danach die Zeile `docs/reviews/` der Plan-Tabelle auf den Ist-Stand setzen (Verification-Report liegt vor).
2. **§7 Closure-Notiz:**
   - DoD-Haken setzen.
   - Auf §3 *Umsetzung* und *Fix-Runde* verweisen, dort stehen Fundstellen, Vorab-Suche und Sensoren (V-2).
   - Den Gate-Lauf nach den Closure-Commits mit Exit-Code eintragen.
   - Die DoD-Klammer „`0.19.0`, belegt durch `e5be6620`“ als Berichtigung eines Zustandsmarkers bestätigen (Fix-Runde, oben).
   - Hinweise an den Projektinhaber: Review F-3, F-4, F-6 und V-1. Die Rolle `docs/user/` als Kandidat der Modus-Deklaration steht in §8.
3. **Steering-Loop-Lerneintrag:**
   - Die Finding-Klassen des Review-Reports.
   - Die Lücke ohne Gate für die Kopien in `cli-json-contract.md` gegen [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) und [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) (Plan §1 *NICHT*).
   - Die Grenze der Belegsicht von `doc-trace` (V-4).
   - Die ungepinnte Bedingung „nur bei“ (V-1).
4. **Register (Entscheidung F-5):**
   - Nachtrag in der Evidence-Datei des Vorgangs [`slice-v2-spec-devcontainer-start`](../plan/planning/done/slice-v2-spec-devcontainer-start.md) unter [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../plan/planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) und ein Absatz in `state.md` dieses Eintrags, ohne neuen Zähler und ohne neue Datei. Inhalt: die Drift der Kopien in `cli-json-contract.md` §2 und §3 (Enum, `description`, `allOf`, Quelle des zweiten Zitatteils).
   - Entscheiden, ob Review F-1 und F-2 unter [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) zählen, und wenn ja, mit einer Evidence-Datei nach der Kennung dieses Slice.
   - Sonst in §7 notieren, dass keine weitere Beobachtung angefallen ist.
5. **Risiken §6, Ausgang je Zeile.** Vorschläge aus meiner Prüfung:

   | Risiko | Vorschlag | Grund |
   |---|---|---|
   | Versionskollision | entfallen | `0.19.0` ist belegt, die Messstände-CR nachgezogen, eine Colima-CR gibt es noch nicht |
   | Vorab-Suche übersieht Kopie | nicht eingetreten | Reviewer und Verifier fanden keine weitere Kopie |
   | Vertragstext ändert sich nach dem Review | nicht eingetreten | Die Fix-Runde berührt `spec/` nicht. |
   | Wächter ohne Wirkung | entfallen | M3 und M4 vor/nach |
   | Spec-Bedingung ohne Test | entfallen | Zwei Tests pinnen S4. |
   | Plan-Aussage stimmt nicht oder Datei fehlt | eingetreten und behoben | Messstände-Zeile in `open/README.md` (Abweichung 1), Review F-1 und F-2 |

6. **Drei Paarungen** (Anker · Folge-Slice · Register): Anker nur, wenn ein Eintrag `liegt in` gesetzt wird. Einen Folge-Slice gibt es nicht; das Enum, das Gate für die Kopien und V-1 sind ausdrücklich nicht angelegt. Register: Punkt 4.
7. **Roadmap:** Den Ruhe-Marker „Nichts in Arbeit.“ zurücksetzen, sobald der Slice `in-progress/` verlässt (sonst `planning-drift`).
8. **`open/README.md` Zeile 49:** „geliefert (liegt in `done/`)“.
9. **Reihenfolge** ([`AGENTS.md`](../../AGENTS.md) §3.3): erst der Inhalt (Closure-Notiz, Haken, Register), dann der reine Move per Make-Target zum Verschieben von Slices (Ziel `done`), dann der Verweis-Nachzug. Dieser zieht auch die Links dieses Reports und des Review-Reports auf `in-progress/` nach. Beide Commits gehen in denselben Push.
10. **Danach:** Plan-Tabelle erneut gegen `git diff --stat 95d43213..HEAD` abgleichen (ohne `cecf6876`). Dann `make docs-check` und `make gates` auf dem letzten Commit, Exit-Codes in §7.

## Minimal Handoff

```text
Role: Verifier
Input context: Slice §1–§8 (Umsetzung, Fix-Runde), Review-Report 2026-10-21, git diff 95d43213..ce3c61a3 (ohne cecf6876), spec/lastenheft.md, spec/spezifikation.md, jsontestutil, cli-json-contract.md, Handbuch-Zeile subcommand, alle Umschlag-Stellen im Paket cli
Changed artefacts: dieser Verification-Report (neu, nicht committet)
Evidence: make gates, docs-check, doc-commits, doc-immutable, doc-trace, doc-planning grün; Wortlaut zeichengleich; Schema-Kopie zeilengleich; Mutationen M2, M7, M4 vor/nach, M10 selbst gefahren
Open risks: V-1 (Handbuch "nur bei" je Befehl ungepinnt), Kopien ohne Gate, doc-trace-Sicht ohne LH-NFA-*
Next role: Hauptagent (Commit dieser Datei), dann Planner/Implementer (Closure nach Prüfliste)
```

## Abschluss

Nach dem Anlegen dieser Datei liefen beide Sensoren auf dem Arbeitsbaum (`HEAD` `ce3c61a3` plus diese Datei, ungestagt):

- `make docs-check`: Exit 0, `d-check: 606 Datei(en) geprüft, 0 Befund(e)`. Ein erster Lauf fand elf unverlinkte Kennungen in dieser Datei (`id-unlinked`); sie sind verlinkt oder umformuliert.
- `make gates`: Exit 0, 56 s; `lint` `0 issues.`, `coverage-gate` 91,70 %, `d-check` 606 Dateien und 0 Befunde, `baseline-verify` OK, `ablage-pruefung` OK (26 Fälle).

`git status` zeigt nur diese Datei als ungetrackt.
