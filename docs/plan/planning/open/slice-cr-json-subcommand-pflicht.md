# Slice CR: `subcommand` verpflichtend für die vier gruppierten Befehle der `--json`-Ausgabe — [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Lastenheft, nächste freie Version, heute 0.19.0), Spezifikation und Wächter in `jsontestutil`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht (Modul 6)). Herkunft: offene Entscheidung (e) der Closure von [`welle-devcontainer-start-ohne-editor`](../done/welle-devcontainer-start-ohne-editor.md) §7 („`egress` fehlt in der Pflicht von `subcommand` und im `allOf` von [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json)“). Release-Zuordnung: das nächste Release nach `v0.13.0`; die Nummer vergibt der Release-Cut. Der Slice ändert kein Produktverhalten und hält kein Release auf.

**Bezug:** [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (Scope: Minimalkontrakt aller `--json`-Ausgaben, Pflicht von `subcommand`). Gelesen und unverändert vorgeschlagen (Frage F-2): [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Enum von `command`; Beispielsatz „Bei gruppierten Befehlen wie …“), [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) (Beispiel „z. B. `template`, `config`“). Die Befehle der Lücke: [`LH-FA-DEV-008`](../../../../spec/lastenheft.md#lh-fa-dev-008--egress-restriktion) (`u-boot egress list`), [`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten) (`u-boot devcontainer up` und `down`). Keine ADR berührt.
Anlass: Entscheidung des Projektinhabers vom 2026-10-21 im Gespräch zu (e): „Vertrag schließen“.

**Berührte Spec-Stellen:** `lastenheft.md` [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), §16 (Historie), Versionsfeld; `spezifikation.md` [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), §7 (Historie). `architecture.md`: — (gelesen, Zeile in §3). Der Verweis zeigt aufwärts: Lastenheft und Spezifikation nennen diesen Slice nie.

**Verantwortlich:** pt9912 (Rolle Implementer, Start 2026-10-21); Freigaben: Projektinhaber.

**Autor:** pt9912 (Rolle Planner). **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Lastenheft verlangt `subcommand` in jeder `--json`-Ausgabe der gruppierten Befehle `template`, `config`, `egress` und `devcontainer`, die Spezifikation führt die Pflicht für alle vier, und `jsontestutil` prüft sie für alle vier als Wächter; die Kopien in [`cli-json-contract.md`](../../../maintainer/cli-json-contract.md) folgen. Der Code setzt `subcommand` heute in jedem Umschlag der vier Befehle (Fundstellen unten); der CR schließt den Vertrag um den Bestand und ändert kein Verhalten.

**Entscheidung des Projektinhabers** (2026-10-21, im Gespräch, zu (e)): „Vertrag schließen“ — kleiner CR: in [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) wird `subcommand` für `template`, `config`, `egress` und `devcontainer` verpflichtend (Lastenheft `0.18.0` → `0.19.0`); [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) bekommt `egress` in den `allOf`-Block; `checkSubcommand` prüft alle vier Befehle; Schema-Kopie und Lastenheft-Zitat in `cli-json-contract.md` werden angeglichen; kein Produktcode. **Den Wortlaut legt dieser Plan als Vorschlag vor** (unten); der Projektinhaber hat ihn am 2026-10-21 beim Start freigegeben, unverändert, Variante A (§4, §7).

### Gefundener Stand (gelesen am 2026-10-21 auf `8ae2431e`; Fundstellen nach Datei, Abschnitt und Symbol)

**Vertrag.**

- `spec/lastenheft.md`, [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), Spiegelstrich des Minimalkontrakts: „optional `subcommand` (für gruppierte Befehle wie `template` oder `config`)“; Absatz „Für `--json`-Antworten gilt zusätzlich: … bei `template` und `config` ist `subcommand` verpflichtend, …“ — eine geschlossene Aufzählung ohne „wie“; `egress` und `devcontainer` stehen in keiner der beiden Stellen. Der Absatz „Für `--dry-run`- oder `--diff`-Ausgaben mit `--json` gilt zusätzlich das vollständige Schema …“ macht den Minimalkontrakt auch für die Vorschau bindend.
- `spec/lastenheft.md`, [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run): Enum von `command` mit `egress` und `devcontainer`; Absatz „Bei gruppierten Befehlen wie `command == "template"` oder `command == "config"` muss das Feld `subcommand` gesetzt sein (z. B. `list`, `get`, `set`).“ — Beispiel mit „wie“, in der Anforderung zur Vorschau.
- `spec/lastenheft.md`, [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur): „`subcommand` ist für kommandospezifische Unterbefehle reserviert (z. B. `template`, `config`).“ — Beispiel.
- `spec/lastenheft.md`, Versionsfeld `0.18.0`; §16 letzte Zeile `0.18.0` ([`LH-FA-DEV-010`](../../../../spec/lastenheft.md#lh-fa-dev-010--devcontainer-ohne-editor-starten)).
- `spec/spezifikation.md`, [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), Spiegelstrich: „Bei `command == "template"`, `command == "config"` oder `command == "devcontainer"` ist `subcommand` verpflichtend.“ **Abweichung vom Befund des Auftrags:** Die Spezifikation führt `devcontainer` in der Pflicht schon (eingeführt in `54324651`, Fix-Runde von [`slice-v2-spec-devcontainer-start`](../done/slice-v2-spec-devcontainer-start.md)), `egress` nicht. Damit geht sie heute über [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) hinaus, gestützt nur auf das „wie“ in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), das die Vorschau betrifft; nach dem CR präzisiert sie wieder.
- `spec/spezifikation.md`, [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json): `properties.subcommand.description` „Unterkommando bei gruppierten Hauptkommandos wie `template`, `config` oder `devcontainer`“; `enum` von `command` mit `egress` und `devcontainer`; `allOf` mit je einem Block für `template`, `config`, `devcontainer`, keinem für `egress`. Die zwei Absätze unter dem Schema: `egress list` setzt `subcommand: "list"` und weist `--dry-run` und `--diff` ab ([`SPEC-046`](../../../../spec/spezifikation.md#spec-046--egress-anzeige-u-boot-egress-list-berechnung-ausgabe-und-abgrenzung)); `devcontainer up`/`down` setzt `up` beziehungsweise `down` und weist beide ab ([`SPEC-049`](../../../../spec/spezifikation.md#spec-049--devcontainer-ohne-editor-starten-und-stoppen-u-boot-devcontainer-up-und-down-aufruf-vorab-prüfung-umgebung-ausgabe-und-exit-codes)). Das Schema der Vorschau greift für beide Befehle also nie.

**Code (Produkt, bleibt unverändert).**

- `internal/adapter/driving/cli/egress.go`: jeder Umschlag von `egress list` trägt `"list"` — `jsonArgsValidator(a, "egress", "list", …)` in `newEgressListCommand`, `reportErrorSub(… "egress", "list" …)` in der Funktion `fail` von `runEgressList` (auch für `ErrEgressPreviewNotApplicable`), `newDataEnvelope("egress", "list", …)` im Erfolgspfad. Bares `egress --json` gibt `ErrEgressSubcommandRequired` zurück (`newEgressCommand`), Exit 2, kein Umschlag (Test `TestExecute_Egress_Bare`, prüft `out == ""`).
- `internal/adapter/driving/cli/devcontainer.go`: `jsonArgsValidator(a, "devcontainer", "up"|"down", …)`; `reportDevcontainerError` schreibt `newMinimalEnvelope("devcontainer", subcommand, …)`, gerufen mit `"up"` aus `runDevcontainerUp` und `"down"` aus `runDevcontainerDown`; Erfolg `newDataEnvelope("devcontainer", "up"|"down", …)`. Bares `devcontainer` gibt mit und ohne `--json` `ErrDevcontainerSubcommandRequired` zurück, kein Umschlag (`TestExecute_Devcontainer_BareIsUsageError`, prüft `out == ""`).
- `internal/adapter/driving/cli/template.go`: `reportErrorSub(… "template", "list" …)` und `newDataEnvelope("template", "list", …)`; bares `template --json` → `ErrTemplateSubcommandRequired` ohne Umschlag. `template list` trägt `cobra.NoArgs` ohne Umschlag-Haken: ein Argumentfehler schreibt keinen Umschlag, also auch keinen ohne `subcommand`.
- `internal/adapter/driving/cli/config.go`: bares `config --json` läuft über `runConfigShow` mit `subcommand: "show"`; `get`, `set` und `list` setzen ihren Namen in jedem Pfad (`configArgsValidator`, `reportErrorSub`, `newDataEnvelope`, `newFullEnvelope("config", "set", …)`, `writeConfigEntries`).
- Ein Umschlag entsteht nur in den Befehlsdateien, in `jsonargs.go` (`jsonArgsValidator`, mit `command` und `subcommand` aus der Befehlsdefinition) und in `erroremission.go` (`writeErrorEnvelopeSub`); `root.go` und `cli.go` schreiben keinen. Die neue Pflicht verletzt damit keinen heutigen Umschlag.

**Wächter.**

- `internal/adapter/driving/cli/jsontestutil/jsontestutil.go`, `checkSubcommand`: Pflicht nur bei `cmd == "template" || cmd == "config"`; dieselbe Aufzählung tragen der Doc-Kommentar von `WithSubcommand` („`command ∈ {template, config}`“), der Doc-Kommentar von `AssertMinimalEnvelope` („subcommand Pflicht bei template/config, sonst optional“) und der Zeilenkommentar in `checkSubcommand`. `checkSubcommand` läuft in `AssertMinimalEnvelope` und in `AssertFullEnvelope`.
- `jsontestutil_test.go`: ein Test pinnt die Regel, nur für `template` (`TestAssertMinimalEnvelope_RejectsBadSubcommandForTemplate`); für `config`, `egress` und `devcontainer` gibt es keinen.
- `jsontestutil/drift_test.go`: gleicht nur die Code-Registry ab (`DefaultAllowedCodes` gegen `application.DoctorCheckIDs()` und gegen den markierten Abschnitt der Spezifikation); weder das Enum noch die Schema-Kopie in `cli-json-contract.md`. **Kein Gate gleicht die Kopie gegen [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) ab.**
- **Abweichung vom Befund des Auftrags (Mutation):** Jeder Test, der einen Umschlag von `egress` oder `devcontainer` prüft (`egress_test.go`: `TestExecute_EgressList_JSON`, `TestExecute_EgressList_Errors`, `TestExecute_EgressList_RejectsPreviewFlagsAndArguments`; `devcontainer_test.go`: alle Aufrufe von `AssertMinimalEnvelope`, darunter `TestExecute_DevcontainerDown`), übergibt schon `WithSubcommand(…)`. Ein geleertes `subcommand` im Umschlag von `egress list` färbt `make test` deshalb heute schon rot — über den Einzel-Pin, nicht über eine generische Prüfung. Die Mutation, die die neue Prüfung belegt, steht unter *Wächter und Mutation*.
- `internal/adapter/driving/cli/jsonenvelope_test.go`: die Doc-Kommentare von `TestMinimalEnvelope_SubcommandOmittedWhenEmpty` („mandatory for template/config“) und `TestMinimalEnvelope_SubcommandPresentWhenSet` („the template/config gruppierten Befehle“) tragen die Aufzählung (Kopplung, [`AGENTS.md`](../../../../AGENTS.md) §3.9 (b)).

**Kopien.**

- [`cli-json-contract.md`](../../../maintainer/cli-json-contract.md) §2: ein Zitatblock, ausgewiesen als „Verbatim-Zitat des Lastenhefts“. Sein erster Teil (Liste des Minimalkontrakts) ist das Lastenheft; sein zweiter Teil („Für `--json`-Antworten gilt zusätzlich:“ als Liste) ist die Spezifikation [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), die die Regeln seit [`slice-v2-spec-lastenheft-verschlankung`](../done/slice-v2-spec-lastenheft-verschlankung.md) trägt (dort Fundstellen-Tabelle), im Stand vor `54324651` (Spiegelstrich „Bei `command == "template"` oder `command == "config"` …“). **Weitere Abweichung vom Befund:** Nicht nur §3, auch der zweite Teil von §2 liegt hinter der Spezifikation, und seine Quellenangabe stimmt nicht.
- §3 „Zitat der Spezifikation ([`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json))“: `description` ohne `devcontainer`, `enum` ohne `devcontainer`, `allOf` nur `template` und `config` — Stand vor `871ed781` (Enum) und `54324651` (`allOf`, `description`). Der Spezifikations-Slice hat die Datei nicht angefasst; seit `871ed781` hat kein Commit sie geändert. Die Kopie war damit auch im Release `v0.13.0` hinter [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json).
- §6.1: drei Reject-Formen (bares `template`, `egress`, `devcontainer`); der Absatz zu `template` begründet den Reject mit der Pflicht, der Satz danach: „`egress` ist ein Help-Parent ohne eigenes Datum und weist `--json` gleichartig ab.“
- §6.9: „Pflicht-`subcommand` ([`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run))“ für die Formen von `config` — bleibt wahr.
- §6.10 und §6.11, Punkt **Vertrag:** „`checkSubcommand` verlangt `subcommand` nur für `template` und `config`“ (§6.10) beziehungsweise „verlangt `subcommand` nur für `template` und `config`“ (§6.11).
- `docs/user/benutzerhandbuch.md`, Abschnitt *JSON-Ausgabe*, Tabelle „Immer vorhanden“, Zeile `subcommand`: „nur bei `config` (`get`, `set`, `list`), `egress` (`list`), `devcontainer` (`up`, `down`) und `template` (`list`)“ — nennt alle vier Befehle; bei `config` fehlt `show` (bares `config --json`, `runConfigShow`), bekannt als V-5b aus [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) §7 und dort nicht behoben (Frage F-4).

### Wortlaut-Vorschlag Lastenheft (alt/neu; vom Projektinhaber freigegeben am 2026-10-21, Variante A, §7)

Das Lastenheft nennt keine Slice- oder Wellen-Kennung. Die Verweise stehen hier als Pfad vom Slice aus; im Lastenheft sind es Anker `#…` der Datei.

**Stelle 1 — [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe), Spiegelstrich `subcommand` des Minimalkontrakts.**

Alt:

> - optional `subcommand` (für gruppierte Befehle wie `template` oder `config`)

Neu:

> - `subcommand` (Unterbefehl; verpflichtend bei den gruppierten Befehlen `template`, `config`, `egress` und `devcontainer`, bei den übrigen Befehlen optional)

**Stelle 2 — dieselbe Anforderung, Absatz „Für `--json`-Antworten gilt zusätzlich“** (nur der Teilsatz zu `subcommand`; sonst bleibt der Absatz unverändert).

Alt:

> … `diagnostics.file` ist optional; bei `template` und `config` ist `subcommand` verpflichtend, und `diagnostics.code` trägt die Kennung der verursachenden Anforderung oder einen dokumentierten tool-internen Code.

Neu:

> … `diagnostics.file` ist optional, und `diagnostics.code` trägt die Kennung der verursachenden Anforderung oder einen dokumentierten tool-internen Code.

Begründung der Form (Variante A): Die Pflicht steht danach an **einer** Stelle der Anforderung, im Feld selbst; heute stehen „optional … wie `template` oder `config`“ und „bei `template` und `config` verpflichtend“ an zwei Stellen, und eine Ergänzung müsste beide treffen — die Lücke dieses CR ist genau eine solche zweite Kopie. „Bei den übrigen Befehlen optional“ hält den Vertrag der übrigen Befehle unverändert (heute „optional“; der Code lässt das Feld bei ihnen weg, `TestMinimalEnvelope_SubcommandOmittedWhenEmpty`).
**Variante B** (kleinerer Eingriff, Frage F-3; vom Projektinhaber nicht gewählt, §7): Stelle 1 wie oben, Stelle 2 statt der Streichung „bei `template`, `config`, `egress` und `devcontainer` ist `subcommand` verpflichtend, und …“ — dann steht die Liste zweimal.

**Stelle 3 — Versionsfeld und Historie (§16).** Die Version ist die nächste freie Minor-Version zum Zeitpunkt des Lastenheft-Commits; beim heutigen Stand `0.19.0`.

Alt: `| Version          | 0.18.0                                                             |`

Neu: `| Version          | 0.19.0                                                             |`

Neue Historie-Zeile (Datum: das des Lastenheft-Commits):

> | 0.19.0 | `<Datum des Lastenheft-Commits>` | [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe): `subcommand` ist auch bei `egress` und `devcontainer` verpflichtend; die Pflicht für die vier gruppierten Befehle (`template`, `config`, `egress`, `devcontainer`) steht im Feld `subcommand` des Minimalkontrakts. | Vereinbarung mit dem Projektinhaber |

(Bei Variante B entfiele der zweite Halbsatz; gewählt ist Variante A, die Zeile trägt beide Halbsätze.)

**Gelesen, unverändert vorgeschlagen** (Frage F-2; der Projektinhaber hat am 2026-10-21 entschieden, dass die Beispiele bleiben, §7):

- [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Absatz „Bei gruppierten Befehlen wie `command == "template"` oder `command == "config"` muss das Feld `subcommand` gesetzt sein (z. B. `list`, `get`, `set`).“ — ein Beispiel („wie“) in der Anforderung zur Vorschau; `egress` und `devcontainer` erzeugen keine Vorschau-Ausgabe, und für jede `--json`-Ausgabe, auch die Vorschau, gilt die Pflicht aus [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe).
- [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur), „(z. B. `template`, `config`)“ — Beispiel für die Befehlsstruktur, keine Liste.
- [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run), Enum von `command` — nennt `egress` und `devcontainer` bereits.

### Wortlaut-Vorschlag Spezifikation (nach dem Lastenheft-Commit, deckungsgleich mit ihm; vom Projektinhaber freigegeben am 2026-10-21, §7)

Die Spezifikation verweist nur aufwärts; keine Slice-, ADR- oder Wellen-Kennung. Die Verweise stehen hier als Pfad vom Slice aus; in der Spezifikation sind es Anker `#…` der Datei beziehungsweise `lastenheft.md#…`, und die Historie-Zeile führt die Kennung der Schema-Stelle wie ihre Nachbarzeilen ohne Link.

**S1 — [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten), Spiegelstrich zu `subcommand`.**

Alt:

> - Bei `command == "template"`, `command == "config"` oder `command == "devcontainer"` ist `subcommand` verpflichtend.

Neu:

> - Bei `command == "template"`, `command == "config"`, `command == "egress"` oder `command == "devcontainer"` ist `subcommand` verpflichtend.

**S2 — [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), `properties.subcommand.description`.**

Alt:

```text
"description": "Unterkommando bei gruppierten Hauptkommandos wie `template`, `config` oder `devcontainer`"
```

Neu:

```text
"description": "Unterkommando der gruppierten Hauptkommandos `template`, `config`, `egress` und `devcontainer`"
```

**S3 — [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), `allOf`:** ein Block für `egress` zwischen den Blöcken für `config` und `devcontainer`, in der Form der übrigen:

```json
    {
      "if": {
        "properties": {
          "command": { "const": "egress" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    },
```

**S4 — [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json), neuer Absatz nach dem Absatz zum Wert `devcontainer`:**

> Die Pflicht von `subcommand` für `template`, `config`, `egress` und `devcontainer` gilt für jede `--json`-Ausgabe ([`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten)); `allOf` gibt sie für das Schema der Vorschau wieder, auch für `egress` und `devcontainer`, die keine Vorschau-Ausgabe erzeugen.

Die Bedingung „erzeugen keine Vorschau-Ausgabe“ pinnen `TestExecute_EgressList_RejectsPreviewFlagsAndArguments` und `TestExecute_Devcontainer_RejectsPreviewFlagsAndArguments`: beide prüfen unter `--dry-run`/`--diff` mit `--json` den Minimalkontrakt mit Exit 2, und `AssertMinimalEnvelope` weist `dryRun`, `diff`, `plannedFiles` und `changes` ab (`checkNoFullFields`) ([`AGENTS.md`](../../../../AGENTS.md) §3.10, letzter Punkt).

**S5 — §7 Historie, neue Zeile** (Datum: das des Spezifikations-Commits):

> | `<Datum>` | [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) und [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) folgen dem Lastenheft 0.19.0: `subcommand` ist auch bei `egress` verpflichtend (Regel, Beschreibung des Feldes, Block in `allOf`); ein Absatz nennt die Herkunft der Pflicht aus dem Minimalkontrakt. |

### Wächter und Mutation

**Wächter** (`internal/adapter/driving/cli/jsontestutil/`): `checkSubcommand` verlangt `subcommand` (vorhanden und nicht leer), wenn `command` einer der vier Werte `template`, `config`, `egress`, `devcontainer` ist — eine Menge an einer Stelle des Helfers, statt der zwei Vergleiche; die Fehlermeldung bleibt unverändert (`subcommand required when command == %q`, gefolgt von der Kennung der Anforderung in Klammern). Die drei Kommentare (Doc von `WithSubcommand`, Doc von `AssertMinimalEnvelope`, Zeilenkommentar) nennen die vier Befehle; der Doc-Kommentar von `AssertMinimalEnvelope` trägt in seinen ersten zwei Zeilen die Kennung der Anforderung doppelt (Rest eines Zeilenumbruchs); das wird im selben Block berichtigt.

**Wächter des Wächters** (`jsontestutil_test.go`): ein Tabellen-Test ersetzt `TestAssertMinimalEnvelope_RejectsBadSubcommandForTemplate` (oder ergänzt ihn) — je Befehl der vier ein Umschlag ohne `subcommand` und einer mit `"subcommand": ""`, beide müssen „subcommand required“ melden; Gegenprobe: `doctor` ohne `subcommand` meldet nichts (die Pflicht greift nur bei den vier); ein Fall über `AssertFullEnvelope`, weil `checkSubcommand` auch dort läuft.

**Mutationen** (vom Implementer per `make test` im Docker-Lauf, danach zurückgesetzt, `git diff` danach leer; Ergebnis und Ausgabe-Ausschnitt in der Evidence):

| # | Mutation | Erwartung |
|---|---|---|
| M1 | `"egress"` aus der Menge in `checkSubcommand` entfernen | rot: Tabellenfall `egress` des Wächter-Tests (nur er) |
| M2 | `"devcontainer"` aus der Menge entfernen | rot: Tabellenfall `devcontainer` (nur er) |
| M3 | Lückennachweis: in `runEgressList` `newDataEnvelope("egress", "list", …)` → `newDataEnvelope("egress", "", …)` **und** in `TestExecute_EgressList_JSON` `WithSubcommand("list")` entfernen | **vor** der Änderung von `checkSubcommand` grün (die Lücke: nichts fängt es), **nach** ihr rot mit „subcommand required when command == "egress"“ |
| M4 | wie M3 für `devcontainer down` (`newDataEnvelope("devcontainer", "down", …)` in `runDevcontainerDown`, Pin `WithSubcommand("down")` in `TestExecute_DevcontainerDown`) | vor der Änderung grün, nach ihr rot |
| M5 | nur Produkt: `subcommand` im Erfolgsumschlag von `egress list` leeren, Tests unverändert | rot schon heute, über `WithSubcommand("list")` — belegt die vorhandenen Einzel-Pins, nicht die neue Prüfung; wird so in der Evidence ausgewiesen |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktcode und Verhalten** — jeder Umschlag der vier Befehle setzt `subcommand` schon (Fundstellen oben); der CR schließt den Vertrag um den Bestand. Keine Änderung an den Abweisungen (`--dry-run`/`--diff` bei `egress list` und `devcontainer`, bare Gruppenbefehle), kein neues Feld, kein neuer Wert von `command`.
- **Eine Prüfung des Enums von `command` in `jsontestutil`** — ein neuer Sensor; der Projektinhaber hat die Pflicht von `subcommand` gewählt, nicht das Enum. Kein Folge-Slice angelegt.
- **Ein Gate, das die Kopien in `cli-json-contract.md` gegen [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) und [`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten) abgleicht** — ein neuer Sensor (Muster: `drift_test.go` für die Code-Registry), nicht Teil der Entscheidung; die Lücke wird bei der Closure als Lerneintrag benannt.
- **Die Kommentare zu den bare-Rejects in `root.go`, `cli.go` und `rootjson_test.go`** („Two `--json` rejects remain: bare `u-boot template` … and bare `u-boot egress`“; `devcontainer` fehlt) — sie betreffen die Abweisungen, nicht die Pflicht, der Diff berührt ihre Symbole nicht (§3.9 (b) greift nicht), und `root.go`/`cli.go` sind Produktcode; als Hinweis an den Projektinhaber im Handback.
- **[`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) (Beispielsatz) und [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur)** — Beispiele mit „wie“ beziehungsweise „z. B.“, keine Listen (Frage F-2); ein Angleichen wäre eine eigene Stelle im selben CR.
- **`CHANGELOG.md`** — kein Produktverhalten (Muster [`slice-cr-cli-json-command-enum`](../done/slice-cr-cli-json-command-enum.md) und [`slice-cr-devcontainer-start-ohne-editor`](../done/slice-cr-devcontainer-start-ohne-editor.md)); die Lastenheft-Version nennt der nächste Release-Cut in seiner Release-Notiz (Muster: „Lastenheft 0.18.0.“ in `[0.13.0]`).
- **`spec/architecture.md`** — nennt weder die gruppierten Befehle noch die Pflicht (Zeile in §3).

## 2. Definition of Done

- [ ] **Lastenheft** (nächste freie Minor-Version zum Zeitpunkt des Commits, heute `0.19.0`): Stellen 1 bis 3 im freigegebenen Wortlaut (Variante nach F-3); Versionsfeld und eine Historie-Zeile mit Datum und „Vereinbarung mit dem Projektinhaber“; kein Slice-, ADR- oder Prozessverweis. *Beleg:* Fundstellen in §7, Vergleich Satz für Satz mit dem freigegebenen Vorschlag (Verification).
- [ ] **Spezifikation** deckungsgleich mit dem Lastenheft: S1 bis S5; die Bedingung in S4 ist durch die zwei genannten Tests gepinnt.
- [ ] **Wächter:** `checkSubcommand` verlangt `subcommand` bei allen vier Befehlen; der Wächter-Test deckt jeden der vier (fehlend und leer), die Gegenprobe `doctor` und einen Fall über `AssertFullEnvelope`; Mutationen M1 bis M5 mit Ergebnis in der Evidence; die Kommentare in `jsontestutil.go` (drei) und `jsonenvelope_test.go` (zwei) nennen die vier Befehle.
- [ ] **Freigabe** des Projektinhabers (Wortlaut Stellen 1 bis 3 und S1 bis S5, Fragen F-1 bis F-5) im Gespräch belegt und in §7 vermerkt.
- [ ] **Vorab-Suche** vor dem Lastenheft-Commit wiederholt (Muster in §3) und jeder Treffer einer Aufzählung mit Befund in §7, auch der unveränderte.
- [ ] `make gates` grün; `make docs-check` ohne Befund; `make doc-trace` ohne neue Waise. Läufe und Exit-Codes in §7.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8); ändert eine Fix-Runde Lastenheft- oder Spezifikationstext, liest ein Nachlauf-Review im frischen Kontext genau diesen Diff, bevor der Verifier startet.
- [ ] Doku-Update: `cli-json-contract.md` (K1 bis K5 in §3) und die Zeile `subcommand` im Benutzerhandbuch (nach F-4); `CHANGELOG.md` unverändert.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (nach F-5), oder in §7 notiert, dass keine Beobachtung angefallen ist.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure geprüft (ohne Welle).

## 3. Plan (vor Code)

Die Tabelle nennt jede Datei, die der Diff voraussichtlich berührt, und die gelesenen Dateien ohne Diff mit ihrer Fundstelle; sie wird nach der letzten Fix-Runde gegen `git diff --stat` abgeglichen ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Reihenfolge der Commits: Lastenheft (nach der Freigabe) → Spezifikation → Wächter (M3 und M4 laufen vor der Änderung von `checkSubcommand`) → Kopien in der Doku → Closure.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | Stellen 1 bis 3 ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe): Spiegelstrich `subcommand`, Teilsatz im Absatz „zusätzlich“; Versionsfeld; §16) |
| `spec/spezifikation.md` | update | S1 bis S5 ([`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten); [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json): `description`, Block `egress` in `allOf`, Absatz; §7 Historie); erst nach dem Lastenheft-Commit |
| `internal/adapter/driving/cli/jsontestutil/jsontestutil.go` | update | `checkSubcommand`: Pflicht für die vier Befehle (eine Menge); Kommentare von `WithSubcommand`, `AssertMinimalEnvelope` (auch die doppelte Kennung) und in `checkSubcommand` |
| `internal/adapter/driving/cli/jsontestutil/jsontestutil_test.go` | update | Wächter-Test nach [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe): je Befehl fehlend und leer (Negativ), Gegenprobe `doctor` (Positiv), ein Fall `AssertFullEnvelope` |
| `internal/adapter/driving/cli/jsonenvelope_test.go` | update (nur Kommentare) | Doc-Kommentare von `TestMinimalEnvelope_SubcommandOmittedWhenEmpty` und `TestMinimalEnvelope_SubcommandPresentWhenSet` nennen die vier Befehle (§3.9 (b)); Testcode unverändert |
| `docs/maintainer/cli-json-contract.md` | update | **K1** §2: zwei Zitatblöcke mit je eigener Quelle — „Zitat des Lastenhefts ([`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe))“ mit Stelle 1 (der Block zitiert den Absatz „zusätzlich“ des Lastenhefts nicht, Stelle 2 berührt ihn deshalb nicht), „Zitat der Spezifikation ([`LH-NFA-USE-004.a`](../../../../spec/spezifikation.md#lh-nfa-use-004a--regeln-für---json-antworten))“ mit S1; **K2** §3: Schema wie [`SPEC-001`](../../../../spec/spezifikation.md#spec-001--json-schema-der-vorschau-ausgabe---dry-run---json) nach S2 und S3 (`description`, `enum` mit `devcontainer`, `allOf` mit vier Blöcken); **K3** §6.1: der Satz zu `egress` nennt die Pflicht für `egress` und `devcontainer`; **K4** §6.10 und **K5** §6.11, Punkt *Vertrag*: [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) macht `subcommand` für den Befehl verpflichtend, `checkSubcommand` prüft es, das Enum prüft `jsontestutil` nicht. §6.9 („Pflicht-`subcommand`“, wahr), §6-Einleitung („zwölf Werte“, wahr) und der Kopf (`Status`/`Datum`, seit dem Umzug `e363b3ba` nicht fortgeschrieben, auch nicht vom Befehl-Slice `06299bbd`) bleiben |
| `docs/user/benutzerhandbuch.md` | update (nach F-4) | Abschnitt *JSON-Ausgabe*, Zeile `subcommand`: bei `config` zusätzlich `show` (`runConfigShow`); keine Kennung (§5 Regel 6). Sonst gelesen und unverändert: Abschnitt zu `devcontainer up`/`down` („`subcommand` ist `up` oder `down`“) |
| `docs/plan/planning/open/README.md` | update | Bestand: Zeile dieses Slice (Plan-Commit); bei der Closure Zustand der Zeile |
| `docs/plan/planning/<Lifecycle>/slice-cr-json-subcommand-pflicht.md` (diese Datei) | neu, update | Plan; beim Start `Verantwortlich:` und Freigabe-Vermerk; Closure (§7); Moves per `make slice-mv` |
| `docs/plan/planning/in-progress/roadmap.md` | update | Zeile unter *Offene Wellen* („Nichts in Arbeit.“) beim Start und bei der Closure; sonst nichts (Slices ohne Welle stehen nicht in der Roadmap, gelesen: `open/slice-cr-lastenheft-messstaende.md` und `open/slice-v1-devcontainer-start-colima-messung.md` fehlen dort) |
| Dateien mit Verweis auf diesen Slice (heute nur `open/README.md`) | update (nur Verweise) | `make slice-mv` zieht den Pfad nach |
| `docs/reviews/` | neu | Review-Report und Verification-Report dieses Slice |
| `docs/plan/planning/observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/` | update (Closure, nach F-5) | Drift der Kopien in `cli-json-contract.md` §2 und §3: Vorschlag Nachtrag in `evidence/slice-v2-spec-devcontainer-start.md` und Absatz in `state.md`, kein neuer Zähler (ein Vorgang zählt einmal) |
| `docs/plan/planning/observations/BEO-DOK/…` (weitere) | update (Closure, falls angefallen) | z. B. `plan-aussage-ueber-code`, `vertragsaenderung-nach-review`; jeder Fund mit Evidence-Datei `slice-cr-json-subcommand-pflicht.md` |
| `docs/plan/planning/open/slice-cr-lastenheft-messstaende.md` | update (falls dieser CR `0.19.0` belegt) | „nächste freie Version nach 0.18.0“ an drei Stellen (Titel, Frage 2 der Historie, DoD „zuletzt `0.18.0`“) auf `0.19.0`; Muster [`slice-cr-devcontainer-start-ohne-editor`](../done/slice-cr-devcontainer-start-ohne-editor.md) §3 |
| `spec/architecture.md` | **unverändert** | gelesen: kein Treffer für `subcommand` oder „gruppiert“; die drei Verweise auf [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe) (§2.5 `adapter/driving`, Punkt *Maschinenlesbare Ausgabe*; §7.2, Dual-Classifier-Regel; §10, *Nicht geplant*: HTTP-Driving-Adapter) nennen keine Pflicht von `subcommand` |
| `internal/adapter/driving/cli/egress.go`, `devcontainer.go`, `template.go`, `config.go`, `jsonargs.go`, `erroremission.go` | **unverändert** | Fundstellen unter *Gefundener Stand*; die Kommentare in `template.go` (`ErrTemplateSubcommandRequired`, `newTemplateCommand`: „make `subcommand` mandatory for command="template"“) und `erroremission.go` (`reportErrorSub`, `writeErrorEnvelopeSub`: Pflicht für die Formen von `config`) bleiben wahr; M3/M4/M5 ändern `egress.go` und `devcontainer.go` nur vorübergehend |
| `internal/adapter/driving/cli/root.go`, `cli.go`, `rootjson_test.go` | **unverändert** | Kommentare zu den bare-Rejects (`buildRootCommand`, Feld `json` in `App`, `TestRootJSON_AllFormsRespondPostTClose`) nennen `template` und `egress`, nicht die Pflicht; siehe §1 *NICHT* |
| `internal/adapter/driving/cli/egress_test.go`, `devcontainer_test.go`, `config_acceptance_test.go`, `template_acceptance_test.go`, `rootjson_test.go` | **unverändert** | pinnen `subcommand` je Aufruf mit `WithSubcommand(…)`; `egress_test.go` und `devcontainer_test.go` werden nur für M3/M4 vorübergehend geändert |
| `internal/e2e/devcontainer_start_docker_test.go` | **unverändert** | pinnt im Fehler-Umschlag von `devcontainer up --json` `Subcommand == "up"` (Build-Tag `docker`, läuft nicht in `make test`) |
| `internal/adapter/driving/cli/jsontestutil/drift_test.go` | **unverändert** | gleicht nur die Code-Registry ab (`TestDriftGate1_MapVsDoctorCheckIDs`, `TestDriftGate2_MapVsMarkdownDoc`) |
| `README.md`, `README.de.md`, `docs/user/benutzerhandbuch-standard.md`, `CHANGELOG.md` | **unverändert** | gelesen: `README.md` nennt „subcommands“ nur allgemein (Exit-Codes, Befehlstabelle, Release-Zeilen `v0.1.0` und `v0.4.0`), `README.de.md` kein Treffer, das Standard-Handbuch kein Treffer für `subcommand`; `CHANGELOG.md` trägt die Pflicht nur in zwei Einträgen von `[0.4.0]` (`config` mit „Pflicht-`subcommand`“, `template list`), die Chronik des Release bleiben |
| `docs/plan/adr/` | **unverändert** | keine ADR nennt die Pflicht von `subcommand` (Suche ohne Treffer außer einer unverwandten Zeile in `0012`) |

**Vorab-Suche** (Pflicht vor jeder Vertragsänderung; nach dem Kandidaten der Regel aus [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md): außer den Namen der Befehle die Felder und Werte, die sie setzen, und jede Aufzählung). Planner, 2026-10-21, über `spec/`, `docs/user/`, `docs/maintainer/`, `docs/plan/adr/`, `internal/`, `cmd/`, `README.md`, `README.de.md`, `CHANGELOG.md`, `harness/`, `AGENTS.md` nach `subcommand` in Verbindung mit „Pflicht“, „verpflichtend“, „required“, „mandatory“, nach „gruppiert“ und „grouped command“, nach den Aufzählungen „`template` oder `config`“, „`template` und `config`“, „template/config“, „template, config“, „{template, config}“; die Treffer-Absätze gelesen (Host-`grep` als Text-Werkzeug, kein Gate). Jeder Treffer steht oben unter *Gefundener Stand* oder in der Tabelle. Der Implementer wiederholt die Suche nach dem Lastenheft-Commit mit denselben Mustern und zusätzlich nach den neuen Werten (`egress`, `devcontainer`) in Verbindung mit `subcommand`.

## 4. Trigger

**Start** (`open` → `in-progress`, ggf. über `next`): `Verantwortlich:` gesetzt; kein anderer Slice in `in-progress/` (WIP-Limit 1); der Projektinhaber hat entschieden (erfüllt am 2026-10-21, Entscheidungen in §7):

- **F-1 — Wortlaut:** Stellen 1 bis 3 (Lastenheft) und S1 bis S5 (Spezifikation) freigegeben, wie vorgelegt oder geändert. — *Entschieden:* freigegeben wie vorgelegt.
- **F-2 — Beispiele:** [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) („Bei gruppierten Befehlen wie …“) und [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) („z. B. `template`, `config`“) bleiben (Vorschlag) oder werden angeglichen (dann je eine weitere Stelle im selben CR). — *Entschieden:* sie bleiben unverändert.
- **F-3 — Form:** Variante A (Pflicht einmal, im Feld `subcommand`; Teilsatz im Absatz entfällt; Vorschlag) oder Variante B (Liste an beiden Stellen). — *Entschieden:* Variante A.
- **F-4 — Handbuch:** `show` bei `config` in der Zeile `subcommand` mitnehmen (Vorschlag: ja, ein Wort in der Zeile, die die jetzt bindende Pflicht für Nutzer beschreibt) oder streichen; das Handbuch stand nicht in der Aufzählung der Entscheidung. — *Entschieden:* mitnehmen, nur mit Werten, die der Code setzt.
- **F-5 — Register:** Die Drift der Kopien in `cli-json-contract.md` §2 und §3 entstand im Vorgang [`slice-v2-spec-devcontainer-start`](../done/slice-v2-spec-devcontainer-start.md), der in [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) schon einen Beleg trägt. Vorschlag: Nachtrag in dessen Evidence-Datei, kein neuer Zähler (Register-README: ein Vorgang zählt einmal); die Alternative, eine eigene Evidence-Datei, bräuchte einen anderen abgeschlossenen Vorgang als Kennung. — *Entschieden:* Nachtrag, kein neuer Zähler, keine neue Datei.

Unabhängig von [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) und dem Folge-CR für Colima aus [`slice-v1-devcontainer-start-colima-messung`](slice-v1-devcontainer-start-colima-messung.md): wer zuerst ins Lastenheft schreibt, bekommt die kleinere Versionsnummer.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Projektinhaber erweitert den Gegenstand, etwa um eine offene Regel „jeder Befehl mit Unterbefehlen“ statt der geschlossenen Liste oder um eine Prüfung des Enums oder der Kopien in `jsontestutil` (neuer Sensor).
- `in-progress` → `open` (blockiert): die Freigabe steht aus oder der Projektinhaber lehnt den Wortlaut ab; das Lastenheft bleibt unverändert.

## 5. Closure-Trigger

DoD vollständig, Review-Report und Verification-Evidence liegen vor, Freigabe des Projektinhabers vermerkt, Closure-Notiz mit Steering-Loop-Lerneintrag geschrieben (§7); Abgleich der Plan-Tabelle gegen `git diff --stat` nach der letzten Fix-Runde.

## 6. Risiken und offene Punkte

- **Der Wortlaut wird nicht freigegeben oder geändert.** — **Ausgang:** entfallen; der Projektinhaber hat ihn am 2026-10-21 unverändert freigegeben (§7).
- **Die Versionsnummer kollidiert** mit [`slice-cr-lastenheft-messstaende`](slice-cr-lastenheft-messstaende.md) oder dem Folge-CR für Colima. — **Ausgang:** weiter offen; Version ist die nächste freie zum Zeitpunkt des Commits, Versionsangaben im jeweils anderen offenen Slice werden nachgezogen (Zeile in §3).
- **Die Vorab-Suche übersieht eine weitere Kopie** ([`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), Schwelle erreicht). — **Ausgang:** weiter offen; Suche mit Feldern, Werten und Aufzählungen (§3), Wiederholung durch den Implementer, Review prüft die Fundstellen.
- **Vertragstext ändert sich nach dem Review** ([`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md)). — **Ausgang:** weiter offen; Nachlauf-Review über den Diff (DoD).
- **Der Wächter bleibt ohne Wirkung belegt:** alle heutigen Tests pinnen `subcommand` schon einzeln, die generische Prüfung greift erst bei einem Test ohne Pin. — **Ausgang:** weiter offen; Wächter-Test im Paket `jsontestutil` und Mutationen M1 bis M4, M5 ausdrücklich als Beleg der Einzel-Pins.
- **Eine Bedingung der Spezifikation ohne Test** (S4, „erzeugen keine Vorschau-Ausgabe“; [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md)). — **Ausgang:** weiter offen; die zwei pinnenden Tests stehen in §1, der Review liest sie gegen den Satz.
- **Eine Aussage dieses Plans über Code stimmt nicht oder eine berührte Datei fehlt in §3** ([`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md), verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10). — **Ausgang:** weiter offen; Abgleich gegen `git diff --stat` nach der letzten Fix-Runde.

## 7. Closure-Notiz

*Die Closure-Notiz wird bei der Closure geschrieben. Beim Start vermerkt sind nur die Entscheidungen des Projektinhabers.*

**Entscheidungen des Projektinhabers (2026-10-21, im Gespräch; an den Implementer übermittelt durch den Hauptagenten, beim Start vermerkt).**

1. **F-1 — Wortlaut freigegeben:** Lastenheft Stellen 1 bis 3 und Spezifikation S1 bis S5, wie in §1 vorgeschlagen.
2. **F-3 — Variante A:** Die Pflicht steht nur im Spiegelstrich `subcommand` des Minimalkontrakts; der Teilsatz „bei `template` und `config` ist `subcommand` verpflichtend“ im Absatz „Für `--json`-Antworten gilt zusätzlich“ entfällt; die Historie-Zeile `0.19.0` trägt beide Halbsätze.
3. **F-2 — Beispiele bleiben:** Die Beispielsätze in [`LH-FA-CLI-007`](../../../../spec/lastenheft.md#lh-fa-cli-007--dry-run) und [`LH-SA-CLI-001`](../../../../spec/lastenheft.md#lh-sa-cli-001--befehlsstruktur) bleiben unverändert.
4. **F-4 — Handbuch mitnehmen:** In der Tabellenzeile `subcommand` wird bei `config` der Wert `show` ergänzt; die Zeile nennt nur Werte, die der Code setzt. Gelesen in `internal/adapter/driving/cli/config.go`: `newConfigCommand` (`configArgsValidator(a, "show", …)`) und `runConfigShow` setzen `show`, `newConfigGetCommand`/`runConfigGet` `get`, `newConfigSetCommand`/`runConfigSet` `set`, `newConfigListCommand` (`configArgsValidator(a, "list", …)`) und `runConfigList` (`writeConfigEntries(out, "list", …)`) `list` — `list` in der heutigen Zeile stimmt.
5. **F-5 — Register:** Nachtrag in der bestehenden Evidence-Datei `evidence/slice-v2-spec-devcontainer-start.md` von [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md), kein neuer Zähler und keine neue Datei; geschrieben bei der Closure (Plan-Tabelle §3).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt werden `spec/` und `docs/plan/` (DOK, Greenfield) und `internal/adapter/driving/cli` (CLI, Greenfield; `jsontestutil/` und `jsonenvelope_test.go` liegen darunter) laut [`harness/conventions.md`](../../../../harness/conventions.md#modus-deklaration-pro-sub-area); die Schwelle aus zwei von drei Achsen ist dort erfüllt. `docs/maintainer/` und `docs/user/` stehen weder in der Modus-Deklaration noch unter den abgewiesenen Kandidaten des Audits ([`slice-harness-sub-area-modus-audit`](../done/slice-harness-sub-area-modus-audit.md), gelesen: kein Treffer für beide Pfade); `docs/user/` hat ein eigenes Verzeichnis und eine eigene Konvention ([`AGENTS.md`](../../../../AGENTS.md) §5 Regel 6), also zwei Achsen — ein Kandidat für die Modus-Deklaration, nicht Gegenstand dieses Slice (Hinweis an den Projektinhaber).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer: [`BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht`](../observations/BEO-DOK/vorab-suche-uebersieht-befehlsuebersicht/observation.md) (drei Belege, Schwelle erreicht, Ausgang beim Projektinhaber; dieser Plan wendet den Kandidaten der Regel an, und die Drift in `cli-json-contract.md` gehört zum dritten Beleg, F-5), [`BEO-DOK/vertragsaenderung-nach-review`](../observations/BEO-DOK/vertragsaenderung-nach-review/observation.md) (sieben Belege, Ausgang beim Projektinhaber; Nachlauf-Review in der DoD), [`BEO-DOK/plan-aussage-ueber-code`](../observations/BEO-DOK/plan-aussage-ueber-code/observation.md) und [`BEO-DOK/spec-bedingung-ohne-code-entsprechung`](../observations/BEO-DOK/spec-bedingung-ohne-code-entsprechung/observation.md) (verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.10: jede Aussage über Code an Datei und Symbol gelesen, die Bedingung in S4 mit Test), [`BEO-CLI/kommentar-beschreibt-abwesendes`](../observations/BEO-CLI/kommentar-beschreibt-abwesendes/observation.md) (verkörpert in [`AGENTS.md`](../../../../AGENTS.md) §3.9: die neuen Kommentare nennen die geltende Menge, keine Chronik „bisher nur `template`/`config`“). Evidenz-/Diskrepanz-Risiko: mittel (drei Orte tragen Kopien der Regel).

Alle berührten Sub-Areas sind Greenfield (kein Modus-Begründungsblock nötig).
