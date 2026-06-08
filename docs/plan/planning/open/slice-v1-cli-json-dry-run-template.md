# Slice V1: `template list --json` — Envelope-Migration

> **Status:** `open/` — **T0-Discovery + R1 gefahren (2026-06-08)**.
> R1 (3 Findings: 1 HIGH + 2 MED) hat die zentrale Sub-Decision
> T0-(a) am echten Spec-Text **umgedreht** (Default-`list` →
> **Reject**) und die Tranchen-Tabelle korrigiert (Allowlist-Abbau
> ist Cluster-T_close, nicht dieser Slice). R2/R3 **noch offen**.
> Heute-Stand-Pre-Scan + Sub-Decisions (a)-(e). Letzter
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
  prüft `a.json` selbst und returnt einen `ErrJSON…`/Exit-2-
  Reject (Pattern-Erbe config's `ErrDryRunNotApplicable`-RunE-
  Reject), damit das Verhalten T_close-stabil ist. Das ist die
  **einzige Code-Berührung an bare `template`** in diesem Slice
  (~15 LOC RunE + Pin-Test).

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
- ✅ **Code-Registry-Sektion** in
  [`docs/user/cli-json-output.md`](../../../user/cli-json-output.md)
  §5 erweitert um eine `template`-Sektion (sofern eigene Codes
  emittiert werden — heute null, also evtl. nur Hinweis
  „template list emittiert keine diagnostics-Codes").
- ✅ **Carveouts-Eintrag entfernt**: Zeile aus
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
  Carveouts gestrichen.
- ✅ **bare-`template`-Sub-Decision**: `u-boot template --json`
  Verhalten festgezurrt (Reject oder Default-Subcommand,
  Vorschlag siehe Aufhebungsbedingung).
- ✅ **Allowlist-Erweiterung**:
  [`jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  `jsonAllowlist`-Map enthält `"u-boot template list"` (heute
  schon) und entweder den bare-`template`-Pfad (Default-
  Sub-Decision) oder bleibt bei Reject.
- ✅ **CLI-Pin-Tests**: bestehende `TestRootJSON_AcceptsTemplate
  List_BothFlagPositions`-Logik bleibt grün; **neuer**
  Envelope-Acceptance-Test prüft `command: "template"`/
  `subcommand: "list"`/`data`-Inhalt.
- ✅ **README EN+DE Verweis-Block** auf
  `docs/user/cli-json-output.md` bleibt unverändert (kein neuer
  Doku-Pfad).

## Tranchen (vorgeschlagen)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + R-Runden**: Pre-Scan + Sub-Decisions (a)-(e); T0-(a) per R1 auf Reject festgezurrt. `Data`-Konstruktor seit generate T1 (`bd3de20`) etabliert. | — (Plan-Arbeit) |
| T1 | **Entfällt** — `cliJSONEnvelope.Data` + `newDataEnvelope(command, subcommand, data, diags, exitCode)` seit generate-Slice 4/9 T1 (`bd3de20`) vorhanden inkl. Marshal-Pin-Tests. Template-Slice nutzt sie nur (T2). | — (entfällt) |
| T2 | **`runTemplateList`-Envelope-Migration**: `renderTemplateListJSON` ersetzt das rohe `[]templateJSON`-Array durch `newDataEnvelope("template", "list", dtos, nil, 0)` → `writeEnvelope`. `subcommand: "list"`-Pflicht (T0-(d)). Bestehende `TestRootJSON_AcceptsTemplateList_BothFlagPositions`-Logik bleibt grün (beide Flag-Positionen), aber der Output-Pin wechselt von Array auf Envelope (Breaking-Change am ausgelieferten Surface, T0-(b)). | ~60 |
| T3 | **bare `template --json`-Reject in der RunE** (R1-HIGH-1b): bare-`template`-RunE prüft `a.json` und returnt einen Exit-2-Reject (Pattern-Erbe config `ErrDryRunNotApplicable`-RunE-Reject) statt `cmd.Help()` — T_close-stabil. `template list` bleibt in der Allowlist; bare `template` bleibt rejected (jetzt RunE-getragen, nicht gate-abhängig). + Pin-Test. | ~25 |
| T4 | **Closure** (NICHT Allowlist-Mechanik-Abbau — R1-MED-2: das ist Cluster-T_close-Scope, eigener Schritt nach diesem Slice, Cluster-Slice §T0-(g)): carveouts.md `template list`-Eintrag entfernen; CHANGELOG **`### Changed`** (Breaking: `template list --json` Array→Envelope, T0-(b)); `cli-json-output.md` §6.2/§6-Tabelle (template→done) + §6.10-Sektion; roadmap; Slice nach `done/` mit DoD-Hash-Tabelle. | — (Doku) |

LOC-Schätzung: **~85 LOC** (T2 ~60 + T3 ~25; T1 + Allowlist-Abbau
entfallen) — der **kleinste Slice des Clusters**. Nach template-Slice-
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
