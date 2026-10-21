# Review-Report: [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md) — 2026-10-21

**Review-Art:** Plan- und Code-Review gegen Plan, Entscheidungen des Projektinhabers und Hard Rules (Wartbarkeit). Gegenstand ist der Lastenheft-CR zu [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (`0.18.0` → `0.19.0`), dazu die Spezifikation, der Wächter in `jsontestutil` und die Kopien in der Maintainer- und Nutzer-Doku. Die DoD-Frage gehört dem Verifier und wird hier nicht beantwortet. Keine Einschätzung des Implementers ist ungeprüft übernommen: Jede Aussage unten habe ich an Datei und Symbol gelesen oder per Sensor nachgefahren.

**Gegenstand:** `git diff 95d43213..e4049fd0` (neun Commits von `f68a5473` bis `e4049fd0`, dazu die Message des Plan-Commits `95d43213`; 11 Dateien, 215 Einfügungen, 62 Löschungen). Nicht Gegenstand ist `cecf6876`: Laut `git diff --stat e4049fd0..cecf6876` ändert er eine Zeile in [`slice-v1-release-cut-v0.13.0`](../plan/planning/done/slice-v1-release-cut-v0.13.0.md). Der Arbeitsbaum war zu Beginn sauber, `HEAD` = `cecf6876`, nichts ist gepusht.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `cecf6876` · **Modell:** claude-opus-5-5 · **Datum:** 2026-10-21

**Eingangs-Kontext:**

- Slice [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md): §1 (Ziel, Gefundener Stand, Wortlaut-Vorschläge Stelle 1 bis 3 und S1 bis S5, Wächter und Mutation, Nicht-Ziele), §2 (nur als Scope-Quelle), §3 (Plan-Tabelle und *Umsetzung*), §4 (F-1 bis F-5), §6, §7 (Entscheidungen), §8.
- Die Entscheidungen des Projektinhabers vom 2026-10-21 gelten als gesetzt und werden nicht neu verhandelt: Vertrag schließen; Wortlaut der Stellen 1 bis 3 und S1 bis S5 wie im Plan; Variante A; die Beispielsätze in [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) bleiben; die Handbuch-Zeile bekommt `show`; Register-Nachtrag bei der Closure.
- Lastenheft: [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run), [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur), [`LH-FA-CLI-006`](../../spec/lastenheft.md#lh-fa-cli-006--exit-codes), §16. Spezifikation: [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), §7. `spec/architecture.md`: Suche nach `subcommand` und „gruppiert“ ergibt 0 Treffer. Keine ADR berührt.
- `AGENTS.md` §3.3, §3.4, §3.7, §3.9, §3.10, §5 Regeln 1 und 6; `harness/conventions.md` (Modus: DOK und CLI Greenfield); `harness/review.md`.
- Vorherige Reports am Bereich: [`2026-10-19-slice-cr-cli-json-command-enum`](2026-10-19-slice-cr-cli-json-command-enum.md) (F-3 Zustandsmarker, F-5 Kopie des Enums) und [`2026-10-21-slice-v2-spec-devcontainer-start`](2026-10-21-slice-v2-spec-devcontainer-start.md) (F-1 `devcontainer` in der Pflicht).

**Selbst gelesen:**

- Der vollständige Diff.
- Im Code jede Stelle, an der ein Umschlag der vier Befehle entsteht:
  - `egress.go`: `newEgressCommand`, `newEgressListCommand`, `runEgressList`
  - `devcontainer.go`: `newDevcontainerCommand`, `reportDevcontainerError`, `runDevcontainerUp`, `runDevcontainerDown`
  - `template.go`: `Args` der Gruppe und von `list`
  - `config.go`: `configArgsValidator` und jeder Aufruf von `reportErrorSub`, `newDataEnvelope`, `newFullEnvelope` und `writeConfigEntries`
  - `jsonargs.go` (`jsonArgsValidator`) und `erroremission.go` (`writeErrorEnvelopeSub`)
- Suche über alle Nicht-Test-Dateien des Pakets `cli` nach `writeEnvelope`, `json.Marshal`, `cliJSONEnvelope{`, `FlagErrorFunc` und `Args:`.
- `jsontestutil.go` und `jsontestutil_test.go` im Stand `e4049fd0`.
- Die Tests `TestExecute_EgressList_RejectsPreviewFlagsAndArguments`, `TestExecute_Devcontainer_RejectsPreviewFlagsAndArguments`, `TestExecute_Egress_Bare`, `TestExecute_Devcontainer_BareIsUsageError`, `TestConfigJSON_ShowMinimalDataEnvelope`, `TestConfigJSON_UnknownSubcommandEmitsEnvelope`, `TestConfigJSON_SubcommandAlwaysSet` und `TestRootJSON_AllFormsRespondPostTClose`.
- Die Kommentare in `root.go` (`buildRootCommand`) und `cli.go` (Feld `json` in `App`).

**Sensor-Läufe des Reviewers** (Docker; angegeben ist der Exit-Code des `make`-Prozesses):

- `make gates` auf `HEAD` `cecf6876`, sauberer Baum: **Exit 0**. `lint` 0 issues; `test` grün; `coverage-gate` 91,70 % bei einer Schwelle von 90 %; `docs-check` 604 Dateien, 0 Befunde; `baseline-verify` v6.18.1 OK; `ablage-pruefung` OK (26 Fälle).
- **Mutation M3 nach** (eigener Lauf, Stand `e4049fd0`):
  - Änderung: in `runEgressList` `newDataEnvelope("egress", "", …)`, in `TestExecute_EgressList_JSON` `WithSubcommand("list")` entfernt.
  - `make test`: **Exit 2**. Rot ist genau `TestExecute_EgressList_JSON` mit `egress_test.go:197: subcommand required when command == "egress"`, dahinter die Kennung [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) in Klammern.
  - Zurückgesetzt per `git checkout`; `git status` danach leer.
- **Mutation M1** (eigener Lauf):
  - Änderung: `"egress"` aus der Menge in `checkSubcommand` entfernt.
  - `make test`: **Exit 2**. Rot sind genau `TestAssertMinimalEnvelope_RequiresSubcommandForGroupedCommands/egress/missing` und `/egress/empty` (Paket `jsontestutil`); das Paket `cli` ist `ok`.
  - Zurückgesetzt; `git status` danach leer.
- `make doc-commits RANGE="95d43213^..e4049fd0"`: **Exit 0**, 10 Commits (der Plan-Commit eingeschlossen), 0 Befunde.
- `make docs-check` mit diesem Report im Baum: **Exit 0**, 605 Dateien, 0 Befunde.
- Nicht gelaufen sind `make test-docker` und `make verify-depguard`; die Gründe stehen unter *Nicht geprüft*.

---

## Findings

Es gibt kein HIGH und kein MEDIUM. Der Wortlaut von Lastenheft und Spezifikation stimmt Zeichen für Zeichen mit dem freigegebenen Vorschlag überein. Die Schema-Kopie ist zeichengleich mit [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json). Jede neue Aussage über Code trägt einen Test, der sie pinnt, und der Wächter greift, wie der Plan es sagt (M1 und M3 habe ich selbst nachgefahren). Die zwei LOW betreffen Plan-Buchhaltung. Die INFO-Findings beschreiben Befunde außerhalb des Diffs und ordnen sie zu.

| ID | Kategorie | Befund | Risiko | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|---|
| F-1 | LOW | Die Zeile dieses Slice im Bestand von `open/README.md` lautet weiterhin „Lastenheft-CR auf die nächste freie Version (heute `0.19.0`), nur Plan: … der Wortlaut braucht die Freigabe des Projektinhabers“. Der Wortlaut ist seit `f68a5473` freigegeben, und `e5be6620` belegt `0.19.0`. Dieselbe Tabelle nennt weiter oben (Zeile 33, in `e5be6620` nachgezogen) die nächste freie Version „nach 0.19.0“. | Zwei Zeilen derselben Tabelle widersprechen sich zur nächsten freien Version. Wer die Messstände-CR oder den Colima-CR startet und diese Zeile liest, hält `0.19.0` für frei und den Wortlaut für offen. Der Plan verschiebt die Zeile auf die Closure (§3, Zeile `open/README.md`), die Regel bindet die Änderung aber an den Ablaufpunkt. | `AGENTS.md` §3.10 (ein Zustandsmarker mit Ablaufpunkt wird beim Ablauf entfernt oder auf den Ist-Stand gesetzt) | `docs/plan/planning/open/README.md` · „(heute `0.19.0`)“ (Zeile 49) | Review-only; `make docs-check` prüft keine Zustandsaussagen | Zustandsmarker mit Ablaufpunkt nach der Umsetzung stehen geblieben |
| F-2 | LOW | §3 *Umsetzung* sagt „Jede Message trägt [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md) und [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)“. Die zwei Verschiebe-Commits `9ac26881` und `2f2b5b71` tragen nur die Slice-Kennung (`git log --format=%s 95d43213^..e4049fd0`). | Der Verifier übernimmt die Aussage in die Evidence, obwohl sie ein ungelesenes Artefakt beschreibt; das ist dasselbe Muster wie im Register-Eintrag `BEO-DOK/plan-aussage-ueber-code`. Die Traceability selbst ist erfüllt: Die Slice-Kennung genügt, und `make doc-commits` meldet 0 Befunde. | `AGENTS.md` §3.10 (eine Aussage steht erst im Text, wenn sie gelesen ist) | [`slice-cr-json-subcommand-pflicht`](../plan/planning/done/slice-cr-json-subcommand-pflicht.md) §3 · „Jede Message trägt“ | Review-only (`git log`) | Plan-Aussage über Artefakt nicht gelesen |
| F-3 | INFO | Kommentare außerhalb des Diffs nennen nur zwei bare-Rejects, `devcontainer` fehlt. Betroffen sind `root.go` `buildRootCommand` („Two `--json` rejects remain“ und „remaining `--json` rejects are bare `u-boot template` and“), `cli.go` am Feld `json` von `App` („Bare `template --json` and bare `egress --json` are the RunE-borne rejects“, über einen Zeilenumbruch) und `rootjson_test.go` am Doc von `TestRootJSON_AllFormsRespondPostTClose`, dazu dessen Map `bareRejects` ohne `u-boot devcontainer`. Der Diff berührt weder diese Symbole noch die Abweisungen. `cli-json-contract.md` §6.1, auf den die Kommentare zeigen, führte die dritte Form schon vor dem Diff. Die Drift entstand mit der dritten Abweisung in `07f219b2` ([`slice-v1-cli-devcontainer-up-down`](../plan/planning/done/slice-v1-cli-devcontainer-up-down.md)). | Kein Risiko für diesen Slice. Die Herausnahme im Plan (§1 *NICHT*) ist begründet: §3.9 (b) greift nicht, weil der Diff zwar die Datei ändert, auf die die Kommentare zeigen (§6.1, K3), deren Liste der Abweisungen aber nicht. Wer den Kommentaren folgt, sieht zwei statt drei Formen. Der Befund gehört nicht in diesen Slice, sondern zum Projektinhaber, wie im Handback vorgesehen. | `AGENTS.md` §3.7, §3.9 (b) | `internal/adapter/driving/cli/root.go` · „Two `--json` rejects remain“ | Review-only | Kommentar-Kopplung bei neuer Reject-Form nicht nachgezogen |
| F-4 | INFO | In `jsontestutil.go`, außerhalb der geänderten Symbole, steht die Kennung doppelt. `checkNoFullFields` nennt in seiner Meldung [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) zweimal, durch „/“ getrennt, und der Kommentar in `checkPlannedFiles` trägt [`LH-FA-CLI-008`](../../spec/lastenheft.md#lh-fa-cli-008--diff-ausgabe) zweimal über einen Zeilenumbruch hinweg. Dasselbe Muster berichtigt der Diff im Doc-Kommentar von `AssertMinimalEnvelope`. Im ganzen `internal/` sind es 14 Fundstellen (gleiche Kennung zweimal in Folge); sie stammen aus `8119daa2`, der Umstellung der Zeilenverweise auf Kennungen. Außerdem sagt der Doc-Kommentar von `AssertFullEnvelope` „In diesem Slice (Doctor) wird der Voll-Helper als Stub angelegt, aber NICHT verwendet“, obwohl sieben Testdateien des Pakets `cli` den Helfer aufrufen. | Kein Risiko für den Vertrag. Der nächste Änderer liest eine Chronik statt einer Zusage. Der Plan nennt ausdrücklich nur die eine Stelle, und der Diff berührt diese Symbole nicht. Das gehört nicht in diesen Slice, sondern ist als eigener Befund außerhalb zu führen. | `AGENTS.md` §3.7 | `internal/adapter/driving/cli/jsontestutil/jsontestutil.go` · `minimal envelope must not contain %q` | Review-only | Doppelte Kennung aus Kennungs-Umstellung |
| F-5 | INFO | Die Handbuch-Zeile `subcommand` sagt „`show` für `u-boot config` ohne Unterbefehl“; das ist wahr und durch `TestConfigJSON_ShowMinimalDataEnvelope` gepinnt. `show` steht aber auch im Fehler-Umschlag von `u-boot config <unbekannt> --json` mit Exit 2 (`configArgsValidator(a, "show", cobra.NoArgs)`, gepinnt durch `TestConfigJSON_UnknownSubcommandEmitsEnvelope`). | Die Zeile behauptet nicht, dass `show` nur dort vorkommt, und `status`/`exitCode` unterscheiden die beiden Fälle. Kein Vertragsbruch und keine Änderungspflicht. Ob Konsumenten das Wort missverstehen, kann der Validator prüfen. | `AGENTS.md` §3.10 (Handbuch-Aussage gegen Code) | `docs/user/benutzerhandbuch.md` · „`show` für `u-boot config` ohne Unterbefehl“ | `make test` (beide Tests vorhanden) | Handbuch-Wert ohne Nennung des Fehlerpfads |
| F-6 | INFO | Ein Bestand außerhalb des Diffs: `doctor`, `up`, `down` und `template list` validieren ihre Argumente mit bloßem `cobra.NoArgs` (`Args:` in `doctor.go`, `up.go`, `down.go`, `template.go`). Ein Argumentfehler oder ein unbekanntes Flag unter `--json` schreibt dort keinen Umschlag: Es gibt kein `FlagErrorFunc`, und `root.go`/`cli.go` schreiben keinen Umschlag. Der Doc-Kommentar von `jsonArgsValidator` nennt genau diesen Fall „a [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) violation“. Das Handbuch sagt „Mit `--json` gibt jeder Befehl genau ein JSON-Objekt auf der Standardausgabe aus“ und nennt als Ausnahmen nur die bare Gruppen und `logs --follow --json`. | Für diesen Slice ohne Folgen: Ohne Umschlag gibt es auch keinen Umschlag ohne `subcommand`, und Plan §1 liest das zu `template.go` richtig. Es bleibt eine mögliche Lücke im Vertrag und im Handbuch. Zuständig sind Projektinhaber und Planner (eigener Vorgang). Einen pinnenden Test habe ich nicht gefunden (Suche nach Aufrufen mit Zusatzargument unter `--json`). | [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | `internal/adapter/driving/cli/jsonargs.go` · „JSON on stdout — a“ | Review-only | Argumentfehler unter --json ohne Umschlag |
| F-7 | INFO | Der Body von Commit `6bfec81c` nennt die Verfeinerungs-Kennung [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten). `AGENTS.md` §5 Regel 1 schließt nur `SPEC-<NNN>` und `ARC-<NNN>` aus der Commit-Message aus. Die Spezifikation führt `.a`-Kennungen als Adressen ihrer Verfeinerungen (Abschnitt *Kennungen*), also nahe an der Klasse der Struktur-IDs. Vier frühere Commit-Bodies (unter den letzten 400 Commits) verwenden dieselbe Form, und `make doc-commits` meldet nichts. Die Subjects aller zehn Commits tragen weder `SPEC-` noch `ARC-`. | Kein Risiko. Die Regel ist für diese Form nicht eindeutig; ein Kandidat für eine Klarstellung, keine Änderungspflicht. | `AGENTS.md` §5 Regel 1 | Commit `6bfec81c` · „Spezifikation [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten)“ | `make doc-commits` (grün) | Verfeinerungs-Kennung in Commit-Message |

### Vorschläge (Entscheidung beim Implementer)

- **F-1:** Die Zeile in `open/README.md` auf den Ist-Stand setzen: Lastenheft `0.19.0`, umgesetzt, Wortlaut freigegeben. Im Plan dann den Eintrag „bei der Closure Zustand der Zeile“ entsprechend anpassen. Die Fix-Runde berührt so weder Lastenheft noch Spezifikation, und das Nachlauf-Review aus der DoD wird nicht ausgelöst.
- **F-2:** Die Aussage in §3 auf den gelesenen Stand bringen, zum Beispiel: „Jede Message trägt die Slice-Kennung; die Commits mit Inhalt tragen zusätzlich …“.
- **F-3, F-4, F-6:** In diesem Slice nicht beheben. Sie gehen als Hinweise an den Projektinhaber im Handback und in die Closure-Notiz. Ob daraus ein Register-Eintrag oder ein eigener Slice wird, entscheidet die Closure bzw. der Projektinhaber.
- **F-5, F-7:** Keine Änderung nötig.

## Wortlaut-Vergleich (Schwerpunkt 1)

Ich habe programmatisch per Zeichenvergleich verglichen: Zeilen aus dem Plan (§1, Abschnitte *Wortlaut-Vorschlag*), bei denen das Zitatzeichen und die Pfad-Präfixe `../../../../spec/…` entfernt sind, gegen die Zeilen in `spec/`.

| Stelle | Ergebnis |
|---|---|
| Stelle 1, Spiegelstrich `subcommand` in [`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) | zeichengleich |
| Stelle 2, Teilsatz im Absatz „Für `--json`-Antworten gilt zusätzlich“ (Variante A: der Teilsatz entfällt) | zeichengleich; sonst ist der Absatz unverändert |
| Stelle 3, Versionsfeld `0.19.0` (Spaltenbreite unverändert) und Historie-Zeile `0.19.0` mit `2026-10-21`, beiden Halbsätzen und „Vereinbarung mit dem Projektinhaber“ | zeichengleich, bis auf das eingesetzte Datum |
| Übriges Lastenheft | Der Diff hat genau drei Hunks. Der Beispielsatz in [`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Zeile 250) und das Beispiel in [`LH-SA-CLI-001`](../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) (Zeile 1117) sind unverändert (F-2 der Entscheidungen). |
| S1 [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) | zeichengleich |
| S2 `description` | zeichengleich |
| S3 `allOf`-Block für `egress` zwischen `config` und `devcontainer` | Form und Einrückung wie die übrigen drei Blöcke; vier Blöcke |
| S4, Absatz nach dem Absatz zum Wert `devcontainer` | zeichengleich |
| S5, Historie §7 | zeichengleich. Die Kennung der Schema-Stelle steht ohne Link wie in den Nachbarzeilen; der Plan nennt das so. |

## Wahrheit gegen Code (Schwerpunkt 2)

| Aussage | Fundstelle im Code | Pinnender Test |
|---|---|---|
| Jeder Umschlag von `egress list` trägt `list` | `newEgressListCommand` (`jsonArgsValidator(a, "egress", "list", …)`), die Funktion `fail` in `runEgressList` (`reportErrorSub(… "egress", "list" …)`) und der Erfolgspfad `newDataEnvelope("egress", "list", …)` | `TestExecute_EgressList_JSON`, `…_Errors`, `…_RejectsPreviewFlagsAndArguments` (je `WithSubcommand("list")`) |
| Jeder Umschlag von `devcontainer up`/`down` trägt `up` bzw. `down` | `jsonArgsValidator(a, "devcontainer", "up"\|"down", …)`, `reportDevcontainerError(…, "up"\|"down", …)`, `newDataEnvelope("devcontainer", "up"\|"down", …)` | `devcontainer_test.go` (`WithSubcommand(sub)`); im Docker-Lauf zusätzlich `internal/e2e/devcontainer_start_docker_test.go` (`Subcommand != "up"`) |
| Bares `egress --json`, bares `devcontainer` (mit und ohne `--json`) und bares `template --json` schreiben keinen Umschlag | `Args: cobra.NoArgs` plus RunE-Sentinel in `newEgressCommand`, `newDevcontainerCommand` und `newTemplateCommand`; kein `FlagErrorFunc`, kein Schreiber in `root.go`/`cli.go` | `TestExecute_Egress_Bare` (`out != ""`), `TestExecute_Devcontainer_BareIsUsageError` (`out != ""`), `TestRootJSON_AllFormsRespondPostTClose` |
| `config` setzt `show`/`get`/`set`/`list` in jedem Pfad | `configArgsValidator(a, "show"\|"get"\|"list"\|"set", …)`, alle `reportErrorSub(… "config", …)`, `writeConfigEntries(out, "get"\|"list", …)`, `newFullEnvelope("config", "set", …)` | `TestConfigJSON_SubcommandAlwaysSet`, `TestConfigJSON_ShowMinimalDataEnvelope`, `TestConfigJSON_UnknownSubcommandEmitsEnvelope` |
| S4: „`egress` und `devcontainer` erzeugen keine Vorschau-Ausgabe“ | `flags.DryRun \|\| flags.Diff` ist die erste Prüfung in `runEgressList`, `runDevcontainerUp` und `runDevcontainerDown` und führt auf `reportErrorSub`/`reportDevcontainerError` mit Minimal-Umschlag. Ein Argumentfehler geht mit `previewFlags=false` ebenfalls auf den Minimal-Umschlag. | `TestExecute_EgressList_RejectsPreviewFlagsAndArguments` und `TestExecute_Devcontainer_RejectsPreviewFlagsAndArguments` prüfen `--dry-run` und `--diff` je einzeln, Exit 2, über `AssertMinimalEnvelope` → `checkNoFullFields`. Die Kombination beider Flags hat keinen eigenen Test; der Code nimmt für sie denselben Zweig (`\|\|`). |
| Handbuch: `subcommand` „nur bei `config` (`show` …, `get`, `set`, `list`), `egress` (`list`), `devcontainer` (`up`, `down`) und `template` (`list`)“ | Alle übrigen Befehle übergeben `""` (`add`, `init`, `generate`, `remove`, `up`, `down`, `logs`, `doctor`), und `omitempty` lässt das Feld weg | `TestMinimalEnvelope_SubcommandOmittedWhenEmpty`; die Werte siehe oben; zu `show` siehe F-5 |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Referenzrichtung in Lastenheft und Spezifikation (`AGENTS.md` §3.4, Decken-Regel; Skill-HIGH „Referenzmodell-Verstoß“) | geprüft, ohne Befund. Die neuen Zeilen verweisen nur auf Anker der eigenen Datei oder aufwärts; keine Slice-, ADR-, Wellen- oder Roadmap-Kennung. `spec/architecture.md` ist unberührt. |
| Die Spezifikation präzisiert das Lastenheft, ohne es zu erweitern | geprüft, ohne Befund. Die Menge in S1 und in `allOf` ist gleich der Menge der Stelle 1. Vor dem Diff ging die Spezifikation mit `devcontainer` über das Lastenheft hinaus; jetzt ist sie deckungsgleich. |
| Wächter `checkSubcommand` und Lint-Profil | geprüft, ohne Befund. Die Menge ist ein Literal in der Funktion (kein Paket-Global, `gochecknoglobals`), kein `//nolint`; `checkSubcommand` läuft in beiden Helfern. Ein `subcommand` ohne String-Typ (etwa `null`) wird über `hasSub == false` gemeldet. Die Fehlermeldung ist unverändert. `make lint`: 0 issues. |
| Wächter-Test und Mutationsangaben in §3 | geprüft, ohne Befund. M1 und M3 nach habe ich selbst nachgefahren; beide sind rot wie angegeben, mit genau den angegebenen Tests. Die übrigen Angaben sind am Code plausibel: M2 verhält sich symmetrisch zu M1; M6 trifft nur die Fälle `empty`; M7 trifft nur `TestAssertFullEnvelope_…`, weil kein anderer Test einen Voll-Umschlag ohne `subcommand` prüft; M8 trifft die Gegenprobe und die `doctor`-Tests; M9 trifft die Fälle `set` und jeden Test, der einen der vier Befehle pinnt. M4 und M5 bis M9 habe ich nicht selbst gefahren. |
| Kommentare im Diff (`AGENTS.md` §3.7, §3.9 a bis c) | geprüft, ohne Befund. Die neuen Kommentare sind Zusagen bzw. Rang-Zeiger ([`LH-NFA-USE-004`](../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)), ohne Chronik. Der ersetzte Testname steht nur noch im Slice. Keine neue Funktion sitzt zwischen einem Doc-Kommentar und seiner Funktion. Eigene Suche nach den alten Aufzählungen („template/config“, „`template` oder/und `config`“, „{template, config}“, „template, config“) über `spec/`, `docs/user`, `docs/maintainer`, `docs/plan/adr`, `internal/`, `cmd/`, `harness/`, `AGENTS.md`, READMEs, `CHANGELOG.md`: Treffer gibt es nur noch in den zwei Beispielen nach F-2 der Entscheidungen. Treffer außerhalb des Diffs siehe F-3 und F-4. |
| `cli-json-contract.md` K1 bis K5 (Schwerpunkt 5) | geprüft, ohne Befund. Die Schema-Kopie in §3 ist zeichengleich mit [`SPEC-001`](../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json): 117 Zeilen innerhalb der Code-Zäune, `diff` ohne Ausgabe. §2 hat zwei Zitatblöcke mit richtiger Quelle: das Lastenheft mit dem Minimalkontrakt (wortgleich mit Stelle 1) und die Spezifikation mit [`LH-NFA-USE-004.a`](../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) (wortgleich mit S1, „Tool-interne“ ohne Trennung). K3 in §6.1 stimmt mit Code und Tests überein. K4 und K5 nennen das Lastenheft als Quelle der Pflicht und `checkSubcommand` als Prüfung; die Aussage, dass `jsontestutil` das Enum nicht prüft, ist wahr (`checkCommand` vergleicht nur mit `WithCommand`). §6.9 „Pflicht-`subcommand` ([`LH-FA-CLI-007`](../../spec/lastenheft.md#lh-fa-cli-007--dry-run))“ bleibt wahr, weil der Beispielsatz `config` nennt. |
| Nutzer-Doku (`AGENTS.md` §5 Regel 6) | geprüft, ohne Befund. Die geänderte Handbuch-Zeile trägt keine Kennung und keinen Verweis auf Planungs- oder Entwicklerdokumente. |
| Plan-Tabelle gegen `git diff --stat 95d43213..e4049fd0` (Schwerpunkt 6) | geprüft, ohne Befund. Jede der 11 Dateien hat einen Eintrag. Einträge ohne Diff sind `docs/reviews/` (dieser Report) und das Register (Closure); beides ist begründet. Keine Datei, die als „unverändert“ geführt ist, liegt im Diff: `git diff --name-only` über `spec/architecture.md`, die sechs Befehlsdateien, `root.go`, `cli.go`, `CHANGELOG.md`, die READMEs und `docs/plan/adr` liefert 0 Dateien. |
| Commit-Messages (Schwerpunkt 7) | geprüft. Alle zehn Commits tragen die Slice-Kennung; kein Subject und kein Body trägt `SPEC-` oder `ARC-`. Die Reihenfolge der Commits ist die geplante: Lastenheft → Spezifikation → Wächter → Kopien → §3. `make doc-commits` Exit 0. Siehe F-2 und F-7. |
| Lifecycle und `git mv` (`AGENTS.md` §3.3) | geprüft, ohne Befund. `9ac26881` ist ein reiner Move von `open/` nach `in-progress/`, `2f2b5b71` zieht die Verweise nach; die Datei liegt nur in einem Lifecycle-Verzeichnis. Der rote `docs-check` auf `2f2b5b71` ist in §3 benannt und mit `ee0f0c8f` behoben. Bedingung: Beide Commits gehen im selben Push hinaus. |
| Roadmap-Zeile (`AGENTS.md` §3.7) | geprüft, ohne Befund. Die Zeile nennt Zustand, Anker und Herkunft; das Drift-Log ist unberührt. |
| Versionsverweise der Messstände-CR | geprüft, ohne Befund. Titel, Frage 2, DoD und die Zeile in `open/README.md` sagen „nach 0.19.0“ bzw. „zuletzt `0.19.0`“. Sonst nennt kein offenes Slice die nächste freie Version. |
| Dual-Classifier, Sprachvertrag, Architekturgrenze, Dateisicherheit, Carveouts, ADR-Immutabilität, Baseline-Bump, Replay/Golden | nicht berührt: kein Produktcode, kein Sentinel, kein CLI-String, kein Import, keine ADR, kein Generator |
| Modus-Deklaration | geprüft, ohne Befund. `spec/` und `docs/plan/` (DOK) sowie `internal/adapter/driving/cli` (CLI) sind Greenfield, und der Diff verschiebt keinen Modus. Der Hinweis in Slice §8 zu `docs/user/` als Kandidat ist nicht Gegenstand. |

## Nicht geprüft

- `make test-docker`: nicht gelaufen. Der Diff berührt keinen Docker-Pfad und keinen Produktcode. Kein Test mit `//go:build docker` ruft `jsontestutil` (eigene Suche: 0 Dateien). `internal/e2e/devcontainer_start_docker_test.go` pinnt `Subcommand == "up"` über eine eigene Struktur und ist unverändert.
- `make verify-depguard`: nicht gelaufen. Der Diff ändert keine Importregel und keinen Import; `jsontestutil.go` importierte `slices` schon vorher.
- Die Mutationen M2, M4 und M5 bis M9 habe ich nicht selbst gefahren; ihre Plausibilität habe ich am Code geprüft (siehe Negativbefunde).
- DoD-Abgleich, Verification-Evidence und Register-Nachtrag (F-5 der Entscheidungen): Aufgabe von Verifier und Closure.
- Validierung gegen den Bedarf der Konsumenten: Aufgabe des Validators (siehe F-5 und F-6).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Zustandsmarker mit Ablaufpunkt nach der Umsetzung stehen geblieben · Plan-Aussage über Artefakt nicht gelesen · Kommentar-Kopplung bei neuer Reject-Form nicht nachgezogen · Doppelte Kennung aus Kennungs-Umstellung · Handbuch-Wert ohne Nennung des Fehlerpfads · Argumentfehler unter --json ohne Umschlag · Verfeinerungs-Kennung in Commit-Message

Die Klasse von F-1 trägt denselben Namen wie F-3 im Report [`2026-10-19-slice-cr-cli-json-command-enum`](2026-10-19-slice-cr-cli-json-command-enum.md). Das Register zählt Zustandsmarker unter [`BEO-DOK/plan-aussage-ueber-code`](../plan/planning/observations/BEO-DOK/plan-aussage-ueber-code/observation.md) (verkörpert in `AGENTS.md` §3.10); dort gehören auch F-1 und F-2 hin. Die Zuordnung entscheidet die Closure.

## Verdikt

**Merge-blockierend:** nein. Es gibt kein HIGH und kein MEDIUM. Die zwei LOW sind Plan-Buchhaltung ohne Vertragsbezug; ihre Behebung berührt weder Lastenheft noch Spezifikation, also wird kein Nachlauf-Review nach DoD fällig.

**GO für den Verifier.** Die Fix-Runde zu F-1 und F-2 kann vor oder mit der Closure laufen.

**Übergabe:** Die Findings gehen an den Implementer. F-3, F-4 und F-6 gehen als Hinweise an den Projektinhaber, als Befunde außerhalb dieses Slice. Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
