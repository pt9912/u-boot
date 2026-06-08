# Slice V1: `template list --json` — Envelope-Migration

> **Status:** `next/` — **T0-Discovery + R1+R2+R3 gefahren,
> Asymptote erreicht, Lifecycle `open/`→`next/` (2026-06-08)**.
> R1 (1 HIGH + 2 MED) drehte T0-(a) am Spec-Text um (→ Reject).
> R2 (0 HIGH + 2 MED + 3 LOW) ergänzte den Error-Envelope-Pfad
> (T0-(f)) + Envelope-Asymmetrie. R3 (0 HIGH + 1 MED + 1 LOW)
> härtete T0-(a) gegen das stärkste Gegenargument (Cluster-„alle
> Enum-Subcommands tragen `--json`" — aufgelöst über Daten-Kommando
> bare-`config` vs. Help-Parent bare-`template`) + Akzeptanz-
> kriterien-Politur. **Asymptote** (HIGH 1→0→0). Heute-Stand-Pre-
> Scan + Sub-Decisions (a)-(f) festgezurrt. Bereit für T2-Start.
> Letzter
> Folge-Slice (9/9) des Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 9). Closure-Pflicht-Slice für den
> Cluster-T_close-Lauf, weil er den bewussten Übergangs-Carveout
> aus dem Doctor-Slice schließt. **Der einfachste Slice des
> Clusters** (Read-only, kein FS-Mutation, kein neuer Sentinel,
> `Data`-Konstruktor seit generate T1 vorhanden) — der Substanz-
> Kern liegt in **zwei Sub-Decisions** (bare-`template`-Verhalten
> §1813-Spannung + Allowlist-Abbau-Grenze), nicht im Code.

## Auslöser

Cluster-Slice §T0-(e)-Reihenfolge stellt diesen Slice auf
**Platz 9** (letzter). Vorgänger
[`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
T3+T4 hat das lokale `--json`-Flag auf `template list` entfernt
und das Output-Format **bewusst unverändert** gelassen: heutige
`templateJSON`-Array-Struktur ohne Spec-§1841-Minimalkontrakt-
Felder. Carveouts-Eintrag in
[`carveouts.md`](../in-progress/carveouts.md) §Temporäre Carveouts
verweist auf diesen Slice als Re-Trigger.

Code-Realität heute:

- [`internal/adapter/driving/cli/template.go`](../../../../internal/adapter/driving/cli/template.go)
  `runTemplateList` ruft `renderTemplateListJSON` bei `a.json`,
  serialisiert `[]templateJSON`-Array via `json.MarshalIndent`.
- Helper `jsontestutil.AssertMinimalEnvelope` rejected die
  heutige Array-Form, weil `status`/`command`/`diagnostics`/
  `exitCode`-Felder fehlen.

## Heute-Stand-Pre-Scan (T0-Discovery 2026-06-08)

[`cli/template.go`](../../../../internal/adapter/driving/cli/template.go),
232 LOC. Zwei Cobra-Commands:

| Aspekt | bare `template` | `template list` |
| --- | --- | --- |
| Cobra-`Use` | `template` (Help-Parent) | `list` |
| Args | `cobra.NoArgs` | `cobra.NoArgs` |
| RunE heute | `cmd.Help()` (Hilfetext) | `runTemplateList` |
| `--json` heute | **rejected** (nicht in Allowlist; Reject-Gate feuert VOR `cmd.Help()`) | **akzeptiert** — `renderTemplateListJSON` → **rohes `[]templateJSON`-Array** (kein Envelope) |
| Allowlist-Stand | — (rejected) | `jsonallowlist.go:37` `"u-boot template list": true` (Doctor-Slice-Carveout) |
| Use-Case | — | `TemplateListUseCase.List` (read-only, `embed.FS`-Katalog, kein lokaler FS-Read) |
| Sentinels | — | keine fachlichen; nur Katalog-Adapter-IO (Exit 14, in CI nie erreicht — `embed.FS` validiert load-time) |
| DTO | — | `templateJSON{name, description, version, supportedAddOns, generatedFiles, requiredTools, variables[]}` (7 Felder; nil-Slices → `[]`-normalisiert; `templateVariableJSON{name, description, default, required}`) |

**Schlüssel-Befunde des Pre-Scans:**

1. **`template list --json` ist bereits in der Allowlist** (seit
   Doctor-Slice T3/T4) — der Reject-Gate lässt es durch, aber das
   Output-Format ist das **rohe Array**, NICHT der Minimalkontrakt-
   Envelope. Das ist der zu schließende Carveout. Der Slice
   migriert das Format, NICHT den Allowlist-Status.
2. **bare `template` ist ein Help-Parent ohne eigenes Datum**
   (RunE = `cmd.Help()`) — anders als bare `config` (= `show`,
   trägt `u-boot.yaml`-Body). Im `--json`-Modus feuert der
   Reject-Gate VOR `cmd.Help()` → heute Exit 2. **Das ist die
   zentrale Sub-Decision (a)** und steht in Spannung zu
   LH-NFA-USE-004 §1813 (alle zehn Spec-Enum-Subcommands tragen
   `--json`) — siehe Sub-Decisions.
3. **`newDataEnvelope` + `cliJSONEnvelope.Data` existieren bereits**
   (aus generate 4/9 T1, `bd3de20`) — T1 entfällt, der Slice
   konsumiert nur (`newDataEnvelope("template", "list", dtos, nil,
   0)`).
4. **Breaking-Change**: `template list --json` ist ein **bereits
   ausgeliefertes** JSON-Surface (seit Doctor-Slice). Die
   Array→`{…, "data": [...]}`-Migration verschiebt die Top-Level-
   Form → jeder existierende Konsument, der das Top-Level-Array
   liest, bricht. Sub-Decision (b).

## Sub-Decisions (T0-Discovery — füllt sich in R-Runden)

- **T0-(a) bare `template --json`-Verhalten** (HIGH — zentrale
  Decision): heute Reject/Exit 2 (Help-Parent, kein Datum).
  Optionen:
  (i) **Reject Exit 2 beibehalten** (Status quo, Stub-
      Vorschlag). **Adversariale Spannung**: `template` IST ein
      Spec-Enum-Subcommand (`LH-FA-CLI-007` §338 listet zehn,
      inkl. `template`); `LH-NFA-USE-004` §1813 fordert `--json`
      für **alle zehn**. Ein dauerhafter Reject von bare
      `template --json` wäre damit eine **Spec-Lücke** — die
      Cluster-Closure-Hard-Rule fordert „alle Subcommand-Formen
      tragen `--json`". Zu klären in R1: Ist bare `template`
      eine eigene „Subcommand-Form" im Sinne von §1813, oder
      deckt `template list` die `template`-Enum-Pflicht ab?
  (ii) **Default-Subcommand `"list"`**: bare `u-boot config`
       migrierte zu `subcommand: "show"` mit echtem Datum;
       analog könnte bare `template --json` ≡ `template list
       --json` sein (Default-View = Listing). **Vorteil**:
       schließt die §1813-Lücke; **Nachteil**: Doppeldeutigkeit
       Help-vs-List im Human-Mode (bare `template` ohne `--json`
       druckt Help, mit `--json` listet → asymmetrisch).
  (iii) **Minimal-Envelope mit `subcommand: ""` + Hinweis-
        Diagnostic** „use `template list`": trägt einen Envelope
        (§1813 erfüllt) ohne Datum, aber `subcommand: ""`
        kollidiert mit der `command="template"`-§322-Subcommand-
        Pflicht. Wahrscheinlich verworfen.
  **R1-festgezurrt: (i) Reject Exit 2** (HIGH-1 — Discovery-
  Empfehlung (ii) am Spec-Text widerlegt). Begründung:
  - **§1838/§420**: bei `command == "template"` ist `subcommand`
    **verpflichtend** — bare `template` hat keinen natürlichen
    Subcommand, kann also gar keinen spec-validen `command:
    "template"`-Envelope erzeugen. Option (ii) müsste künstlich
    `subcommand: "list"` setzen.
  - **Cluster-Aufhebungsbedingung** (Cluster-Slice §Aufhebungs-
    bedingung bash-Block) verlangt nur `u-boot template list
    --json`, **nicht** bare `u-boot template --json` — bare
    template ist NICHT in der Cluster-Pflicht-Formenliste.
  - **§1813**: `--json` ist „**optional** maschinenlesbare
    Ausgabe" — kein Zwang, dass ein Help-Parent ohne Datum JSON
    emittiert.
  - **Asymmetrie-Wart**: Option (ii) machte `u-boot template`
    (Human → Help) und `u-boot template --json` (→ Liste) zu
    zwei verschiedenen Operationen je nach Flag. Reject ist
    konsistent: „Subcommand fehlt" in beiden Modi.
  **R1-Verfeinerung (HIGH-1b)**: der Reject muss in die **bare-
  `template`-RunE** wandern (nicht nur am Allowlist-Gate hängen).
  Heute feuert der Reject-Gate VOR `cmd.Help()`; sobald der
  Cluster-T_close das Gate + `PersistentPreRunE` abbaut, würde
  bare `template --json` sonst auf `cmd.Help()` fallen und
  **Hilfetext statt Reject** leaken. Fix: bare-`template`-RunE
  prüft `a.json` selbst und returnt `cli.ErrTemplateSubcommand
  Required`/Exit-2 (T0-(f); Pattern-Erbe config's RunE-Reject),
  damit das Verhalten T_close-stabil ist. Im **Human-Modus**
  (`!a.json`) bleibt `cmd.Help()` unverändert — der Reject gilt
  nur im JSON-Pfad. Das ist die **einzige Code-Berührung an bare
  `template`** in diesem Slice (~20 LOC RunE + Sentinel + Pin).

  **R3-Härtung gegen das stärkste Gegenargument** (R3-MED-1): die
  Cluster-Aufhebungsbedingung (`config`-Cluster-Pflicht-Callout)
  formuliert „`LH-NFA-USE-004` gilt für **alle** Spec-Enum-
  Subcommands" — und `template` IST im Enum (§338). Liest man das
  wörtlich, müsste bare `template --json` einen Envelope tragen.
  **Auflösung**: das Argument trägt nicht, weil es **bare `config`
  vs. bare `template`** verwechselt:
  - bare `config` ist ein **Daten-Kommando** (`runConfigShow`
    emittiert den `u-boot.yaml`-Body) → bekam zu Recht einen
    Envelope (`subcommand: "show"`).
  - bare `template` ist ein **reiner Help-Parent** (`RunE =
    cmd.Help()`, kein eigenes Datum) — das `template`-Enum wird
    durch seine **einzige daten-produzierende Form**
    `template list --json` erfüllt; bare trägt nichts, was ein
    Envelope serialisieren könnte, und §1838 verbietet einen
    `subcommand`-losen `command:"template"`-Envelope unabhängig.
  Damit ist Reject die einzige spec-kohärente Wahl — **nicht**
  trotz, sondern **wegen** „alle Enum-Subcommands tragen `--json`"
  (die Form, die das Enum trägt, ist `list`).

- **T0-(b) Breaking-Change-Politik** (MED): Array→Envelope ist
  ein Breaking-Change am ausgelieferten `template list --json`-
  Surface. Optionen:
  (i) **Migrieren + CHANGELOG-Breaking-Note**: Spec-§1841-Pflicht
      schlägt Konsumenten-Verträglichkeit; JSON-Surface ist
      pre-1.0 (v0.4.0), Breaking-Changes sind hier legitim.
  (ii) Array-Form als **permanenten Carveout** behalten —
       verstößt gegen §1841 + Cluster-Closure-Hard-Rule
       (verworfen).
  Plan-Empfehlung: **(i) Migrieren**, Breaking-Change im
  CHANGELOG `### Changed` explizit markiert. Kein Major-Bump
  nötig (pre-1.0).

- **T0-(c) Allowlist-/Reject-Gate-Abbau-Grenze** (MED): nach
  diesem Slice sind **alle 9** Folge-Slices migriert → die
  Allowlist + `applyJSONRejectGate` + `PersistentPreRunE`
  haben keinen Reject-Fall mehr. Wer baut sie ab — dieser
  Slice oder der **Cluster-T_close**? Optionen:
  (i) **Cluster-T_close** baut die Mechanik ab (Cluster-Slice
      §T0-(g) sagt „Cluster-T_close entfernt Allowlist und
      `PersistentPreRunE` komplett"). Template-Slice lässt die
      Mechanik stehen (mit leerer Reject-Menge) und schließt
      nur seinen Carveout.
  (ii) Template-Slice baut die Mechanik gleich mit ab (spart
       einen T_close-Schritt). **Risiko**: vermischt
       Slice-Scope mit Cluster-Scope; der Cluster-Slice-Move
       nach `done/` ist ein separater Closure-Akt.
  Plan-Empfehlung: **(i)** — Boundary sauber halten. Template-
  Slice migriert `template list` + klärt bare-`template`; der
  Allowlist-Mechanik-Abbau gehört zum Cluster-T_close (eigener
  Schritt nach template-Slice-done). Der Stub-T4 („Allowlist-
  Mechanik komplett abbauen") ist damit **falsch zugeordnet** —
  R-Runde korrigiert die Tranchen-Tabelle.

- **T0-(d) `subcommand: "list"`-Pflicht** (LOW): der migrierte
  Envelope MUSS `subcommand: "list"` setzen (`LH-FA-CLI-007`
  §322 Subcommand-Pflicht bei `command="template"`). Trivial,
  aber T-Pin gegen Empty-Subcommand-Drift (analog config).

- **T0-(e) Code-Registry** (LOW): `template list` emittiert
  **keine** `diagnostics[]`-Codes (read-only Happy-Path, leeres
  `diagnostics: []`). `docs/user/cli-json-output.md` §5 braucht
  also **keinen** neuen Code-Registry-Eintrag — nur ggf. einen
  Hinweis. Der Katalog-Adapter-IO-Fehler (Exit 14) ist in CI
  nie erreichbar (`embed.FS` load-time-validiert); falls ein
  Envelope-Error-Pfad doch gebaut wird, nutzt er `LH-NFA-REL-003`
  (bestehender Code, kein neuer Registry-Eintrag).

- **T0-(f) Error-Envelope-Pfad + Envelope-Asymmetrie**
  (R2-MED-1 + R2-MED-2): zwei getrennte Reject/Error-Formen, die
  **nicht** gleich behandelt werden dürfen:
  - **`template list` Error → Envelope-VOLL**: Cluster-Symmetrie
    (logs nutzt `mapLogsErrorToDiagnostic` + `reportError`; §1841
    fordert Envelope auch im Error-Fall). Auch wenn der einzige
    reale Fehler (Katalog-Adapter-IO) in CI unerreichbar ist
    (`embed.FS` load-time-validiert), trägt `template list` einen
    minimalen **`mapTemplateErrorToDiagnostic`** (2 Rows: Katalog-
    IO `LH-NFA-REL-003`/Exit 14; Default `LH-FA-CLI-006`/Exit 1)
    + `reportErrorSub(out, err, …, "template", "list", mapErr,
    nil)`. `subcommand: "list"` ist auch im Error-Envelope
    gesetzt (§322). Bare-`return err` wäre Cluster-Inkonsistenz.
  - **bare `template --json` Reject → Envelope-LOS**: anders als
    `list` KANN bare `template` keinen spec-validen Envelope
    erzeugen, weil §1838/§420 `subcommand` für `command=
    "template"` verpflichtend macht und bare keinen hat. Der
    Reject ist daher ein **plain Usage-Error** (kein Envelope) →
    Exit 2, stderr-Hinweis „use `u-boot template list`". Das
    entspricht 1:1 dem heutigen Gate-Verhalten (`ErrJSON…` ohne
    Envelope) und ist nach T_close-Gate-Abbau via RunE
    T_close-stabil. **Cluster-Ausnahme** (bewusst dokumentiert):
    alle anderen Rejects (config `--dry-run`, logs `--follow
    --json`) tragen einen Envelope; bare-`template` ist die
    EINE envelope-lose Reject-Form, weil §1838 sie erzwingt.
  - **Reject-Sentinel** (Sub-Decision): der bare-RunE-Reject
    braucht einen Exit-2-Sentinel. `ErrJSONNotImplemented`
    (heutiger Gate-Wert) ist semantisch falsch („nicht
    implementiert" — es IST entschieden, nämlich Reject). R2-
    Vorschlag: neuer `cli.ErrTemplateSubcommandRequired`
    (Message „u-boot template requires a subcommand (try `u-boot
    template list`)", Exit 2 via `isUsageError`). R3 bestätigt
    den Namen.

## Aufhebungsbedingung

`u-boot --json template list` (und das Synonym `u-boot template
list --json`) liefert einen Spec-§1841-konformen Minimalkontrakt-
Envelope:

```json
{
  "status": "ok",
  "command": "template",
  "subcommand": "list",
  "diagnostics": [],
  "exitCode": 0,
  "data": [/* heutige templateJSON-Array-Struktur */]
}
```

Carveouts-Eintrag aus `carveouts.md` ist entfernt. Acceptance-
Test pinnt die Envelope-Form via
`jsontestutil.AssertMinimalEnvelope` mit
`WithCommand("template")` + `WithSubcommand("list")`.

`u-boot --json template` (bare) wird **per R1-festgezurrt mit
Exit-Code 2 rejected** (T0-(a) (i)) — und zwar T_close-stabil in
der bare-`template`-RunE selbst (nicht nur am Allowlist-Gate),
damit der spätere Cluster-T_close-Gate-Abbau keinen Hilfetext
leakt (R1-HIGH-1b). Begründung im Detail siehe Sub-Decision
T0-(a): §1838 macht `subcommand` für `command="template"`
verpflichtend (bare hat keinen), die Cluster-Aufhebungsbedingung
verlangt nur `template list`, §1813 macht `--json` optional, und
Default-`list` erzeugte eine Human-vs-JSON-Asymmetrie.

## Akzeptanzkriterien

- ✅ **Envelope-Migration**: `cliJSONEnvelope` mit
  `command: "template"`, `subcommand: "list"`, `Data` als
  `[]templateJSON`. Konstruktor `newDataEnvelope(command,
  subcommand, data any, diags, exitCode)` ist **bereits im
  generate-Slice 4/9 T1 eingeführt** (Ownership aus T0-(p)
  vorgezogen, Commit `bd3de20`); Template-Slice **erbt** nur
  und ruft mit `subcommand="list"` + `data=[]templateJSON`.
  Marshal-Pin-Tests `TestDataEnvelope_DataPresent` /
  `TestDataEnvelope_DataNilOmitted` in
  `jsonenvelope_test.go` decken das Feld bereits ab.
- ✅ **`Data`-Feld im `cliJSONEnvelope`**: ist bereits
  vorhanden (generate T1, Commit `bd3de20`). Type `any` mit
  `omitempty`-Tag, in der Struct-Definition Z. 45 in
  `jsonenvelope.go`. Template-Slice fügt KEIN Feld mehr hinzu,
  sondern verbraucht es nur — die ursprüngliche Sub-Decision
  Plan-Vorgabe ist damit erfüllt; T1-LOC-Schätzung sinkt
  entsprechend (siehe Tranchen-Tabelle).
- ✅ **Error-Envelope-Pfad** (R2-MED-1 / T0-(f)): `template list`
  trägt einen minimalen `mapTemplateErrorToDiagnostic` (Katalog-IO
  → `LH-NFA-REL-003`/Exit 14; Default → `LH-FA-CLI-006`/Exit 1) +
  `reportErrorSub(…, "template", "list", mapErr, nil)`. **Keine
  neue §5-Code-Registry-Sektion** (R3): es werden nur bestehende
  LH-Codes genutzt, keine tool-internen Codes; `template list`
  emittiert auf dem Happy-Path `diagnostics: []`.
- ✅ **Carveouts-Eintrag entfernt**: Zeile aus
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts gestrichen (T4).
- ✅ **bare-`template`-Verhalten festgezurrt** (T0-(a) (i),
  R1+R3): `u-boot template --json` → **Reject Exit 2** via
  `cli.ErrTemplateSubcommandRequired`, RunE-getragen (T_close-
  stabil), **envelope-LOS** (§1838-Ausnahme). Human-Modus
  unverändert `cmd.Help()`. Pin-Test gegen Help-Leak.
- ✅ **Allowlist unverändert**:
  [`jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  behält `"u-boot template list": true`; bare `template` wird
  **NICHT** eingetragen (bleibt rejected, jetzt RunE-getragen).
  Der Allowlist-Mechanik-**Abbau** gehört zum Cluster-T_close
  (R1-MED-2), nicht zu diesem Slice.
- ✅ **CLI-Pin-Tests**: bestehende `TestRootJSON_AcceptsTemplate
  List_BothFlagPositions`-Logik bleibt grün (beide Flag-
  Positionen), aber der schwache „JSON array"-Assert
  (jsonallowlist_test.go:90) wird auf die Envelope-Form
  verschärft (`AssertMinimalEnvelope` + `WithCommand("template")`
  + `WithSubcommand("list")` + `data`-Inhalt). Neuer
  bare-`template --json`-Reject-Pin (Exit 2, kein Envelope,
  Human-Help-unberührt). Empty-Catalog-`data: []`-Pin (R2-LOW-2).
- ✅ **README EN+DE Verweis-Block** auf
  `docs/user/cli-json-output.md` bleibt unverändert (kein neuer
  Doku-Pfad).

## Tranchen (vorgeschlagen)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + R-Runden**: Pre-Scan + Sub-Decisions (a)-(e); T0-(a) per R1 auf Reject festgezurrt. `Data`-Konstruktor seit generate T1 (`bd3de20`) etabliert. | — (Plan-Arbeit) |
| T1 | **Entfällt** — `cliJSONEnvelope.Data` + `newDataEnvelope(command, subcommand, data, diags, exitCode)` seit generate-Slice 4/9 T1 (`bd3de20`) vorhanden inkl. Marshal-Pin-Tests. Template-Slice nutzt sie nur (T2). | — (entfällt) |
| T2 | **`runTemplateList`-Envelope-Migration**: `renderTemplateListJSON` ersetzt das rohe `[]templateJSON`-Array durch `newDataEnvelope("template", "list", dtos, nil, 0)` → `writeEnvelope` (`subcommand: "list"`-Pflicht, T0-(d)). **Error-Pfad** (T0-(f)/R2-MED-1): minimaler `mapTemplateErrorToDiagnostic` (2 Rows) + `reportErrorSub(out, err, …, "template", "list", mapErr, nil)`. **Format-Change** (R2-LOW-1): `MarshalIndent` (indent) → `writeEnvelope`/`json.Marshal` (compact, single-line) — Teil des Breaking-Change. **Test-Update** (R2-LOW-3): `TestRootJSON_AcceptsTemplateList_BothFlagPositions` (jsonallowlist_test.go:90) asserted heute „JSON array" — Pin auf Envelope-Form umstellen (`AssertMinimalEnvelope` + `WithCommand`/`WithSubcommand("list")` + `data`-Inhalt); Both-Flag-Positionen bleiben. **Empty-Catalog-Pin** (R2-LOW-2): leerer Katalog → `data: []` (nicht `null`; `make([]templateJSON,0,…)` normalisiert bereits). | ~60 |
| T3 | **bare `template --json`-Reject in der RunE** (R1-HIGH-1b): bare-`template`-RunE prüft `a.json` und returnt **`cli.ErrTemplateSubcommandRequired`** (R2-MED-2, neuer Exit-2-Sentinel via `isUsageError`; Message „u-boot template requires a subcommand (try `u-boot template list`)") statt `cmd.Help()` — **envelope-LOS** (§1838-Ausnahme, T0-(f)), T_close-stabil. `template list` bleibt in der Allowlist; bare `template` bleibt rejected (jetzt RunE-getragen, nicht gate-abhängig). + Pin-Test. | ~30 |
| T4 | **Closure** (NICHT Allowlist-Mechanik-Abbau — R1-MED-2: das ist Cluster-T_close-Scope, eigener Schritt nach diesem Slice, Cluster-Slice §T0-(g)): carveouts.md `template list`-Eintrag entfernen; CHANGELOG **`### Changed`** (Breaking: `template list --json` Array→Envelope + indent→compact, T0-(b)); `cli-json-output.md` **§6.2** (bestehende „Sonderfall template list --json"-Carveout-Sektion auf Envelope-Form aktualisieren — R2-LOW: KEINE separate §6.10, der template-Inhalt lebt schon in §6.2) + §6-Tabelle (template→done); roadmap; Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) |

LOC-Schätzung: **~90 LOC** (T2 ~60 inkl. Error-Pfad + T3 ~30 inkl.
Sentinel; T1 + Allowlist-Abbau entfallen) — der **kleinste Slice
des Clusters**. Nach template-Slice-
done greift die **Cluster-Closure-Hard-Rule** → der Cluster-Slice
selbst geht via **T_close** nach `done/` (Allowlist + `applyJSONReject
Gate` + `PersistentPreRunE` Abbau, optional Folge-ADR — alles
Cluster-Scope, NICHT dieser Slice).

**Cluster-T_close Forward-Concern** (R1-Notiz): nach dem Gate-Abbau
darf bare `template --json` nicht auf `cmd.Help()` zurückfallen —
deshalb trägt T3 den Reject in der RunE (T_close-stabil). T_close
muss verifizieren, dass nach Mechanik-Abbau alle 9 migrierten Forms
weiterhin korrekt antworten (kein read-only-Form leakt rohen Output).

## Out of Scope

- **HTTP- oder gRPC-Schnittstellen**: ADR-0010 schließt
  explizit aus.
- **Schema-Versionierung** (`schemaVersion: 1`): siehe
  Cluster-Slice §Out of Scope.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  §T0-(e) Platz 9.
- Vorgänger-Slice:
  [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  T3+T4 (Flag-Schnitt + Carveouts-Eintrag).
- Carveouts:
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts §`template list --json`.
- Code-Realität: `internal/adapter/driving/cli/template.go`,
  `internal/adapter/driving/cli/jsonenvelope.go`.
- Spec: `LH-NFA-USE-004` Minimalkontrakt
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md) §1841),
  `LH-FA-TPL-004` Template-Listing.
- ADR: [`ADR-0010`](../../adr/0010-kein-http-driving-adapter.md).
- Phase: V1 (Cluster-Closure-Pflicht).
