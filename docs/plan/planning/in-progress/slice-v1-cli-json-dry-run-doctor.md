# Slice V1: `doctor --json` — Pattern-Vorbild für read-only-Envelope

> **Status:** geplant für v0.4.0 — erster Folge-Slice des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 1). Liefert die gemeinsame Infrastruktur für die
> 9er-Folge-Slice-Serie: Root-PersistentFlag `--json` (T0-(a)),
> Common-Envelope `cliJSONEnvelope` (T0-(c)) und Schema-Helper
> `jsontestutil.AssertMinimalEnvelope` /
> `jsontestutil.AssertFullEnvelope`. Trägt zusätzlich den
> Übergangs-Schnitt für das existierende `template list --json`
> (Review-Round-2-Finding M3). T0 ✅ festgezurrt (§T0-Outcomes
> — acht Sub-Decisions plus Review-Findings H1-L2). In
> `in-progress/`, T1 in Arbeit (`docs/user/cli-json-output.md`
> — Cluster-T1 via L1-Delegation).

## Auslöser

Cluster-Slice `slice-v1-cli-json-dry-run` §T0-Outcomes
festgezurrt: `doctor` ist **Platz 1** der Folge-Slice-Reihenfolge
(T0-(e)) — niedrige Komplexität, schon read-only, „gutes Schema-
Pilot" ohne Architektur-Last. Der schwerere `RecordingFileSystem`-
Wiring-Komplex aus T0-(b) wird im **zweiten** Folge-Slice
`slice-v1-cli-json-dry-run-add` validiert; `doctor` etabliert
vorher das Envelope + den Schema-Helper auf Read-only-Boden.

Spec-Bezug für `doctor --json` (read-only-Pfad):

- **`LH-NFA-USE-004`** (Maschinen-lesbar, V1): Pflicht-`--json`
  für alle 10 Spec-Enum-Subcommands. **Minimalkontrakt** (bindend
  für `--json` ohne `--dry-run`/`--diff`, Lastenheft §1841):
  `status` ∈ `{ok, warn, error}`, `command` ∈ Spec-Enum,
  optional `subcommand` (Pflicht bei `command ∈ {template,
  config}`), `diagnostics` (leer = `[]`), `exitCode` (vgl.
  `LH-FA-CLI-006`)
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
  Zusätzliche Pflichtregeln: `diagnostics[].level` **nur**
  `warn` oder `error` (Lastenheft §1834 — der All-OK-Fall
  serialisiert ein leeres `diagnostics`-Array, **kein**
  `level: "ok"`-Eintrag); `diagnostics[].code` LH-Kennung-
  konform (Lastenheft §1835 + §445); `status`-Kopplung an
  höchstes `level` (`error → "error"`, `warn ohne error →
  "warn"`, sonst `"ok"`).
- **`LH-FA-CLI-007`** (Dry-Run, V1, Voll-Schema): gilt
  ausschließlich für `--dry-run --json` und `--diff --json`
  (Lastenheft §1842, §468). Pflichtfelder
  `dryRun`/`diff`/`plannedFiles`/`changes` plus die Minimal-
  kontrakt-Felder. **`doctor` ist read-only** — der `--json`-
  Pfad trägt den Minimalkontrakt, **nicht** das Voll-Schema.
  Dieser Slice etabliert daher beide Helper-Modi
  (`AssertMinimalEnvelope` und `AssertFullEnvelope`), damit
  der zweite Folge-Slice (`add`) den Full-Mode auf bereits
  stabilem Minimal-Helper aufsetzt.
- **`LH-FA-DIAG-003`** (Severity-Klassifikation): bereits
  implementiert in
  [`doctor.go`](../../../../internal/adapter/driving/cli/doctor.go),
  Exit-Code 11 für Error (und für Warn unter `--strict`). Der
  `--json`-Pfad muss `exitCode` konsistent zu
  [`ErrDoctorFailures`](../../../../internal/adapter/driving/cli/cli.go)
  serialisieren (`LH-FA-CLI-006`).

Heute existierender Pfad:
[`doctor.go:71-78`](../../../../internal/adapter/driving/cli/doctor.go)
trägt **keinen** `--json`-Flag — der Bericht geht ausschließlich
über `writeDoctorReport` als Human-Plaintext an stdout. Der
domain-Typ
[`domain.Diagnostic`](../../../../internal/hexagon/domain/diagnostic.go)
liefert `ID`, `Severity`, `Message`, `Hint` — Kandidaten für
`diagnostics[].code` / `level` / `message` / `hint` im Envelope
(genaue Feldnamen-Mapping = Sub-Decision dieses Slices).

**Übergangs-Schnitt `template list --json` (Cluster T0-(e)
§Review-Round-2-Finding M3):** das lokale `--json`-Flag auf
[`template list`](../../../../internal/adapter/driving/cli/template.go)
kollidiert mit dem neu eingeführten Root-PersistentFlag aus
T0-(a). Dieser Slice **muss** den Schnitt durchführen — lokales
Flag entfernen, `runTemplateList` liest Root-Flag-State,
**Output-Format bleibt unverändert** (die spätere Envelope-
Migration für `template list` ist Platz 9
`slice-v1-cli-json-dry-run-template`).

## Aufhebungsbedingung

`u-boot doctor --json` liefert einen **Minimalkontrakt-konformen**
Envelope (`LH-NFA-USE-004` §1841) und `make test` + `make lint` +
`make docs-check` bleiben grün. Im selben Commit-Set ist das
lokale `template list --json`-Flag entfernt und beide CLI-Pin-
Tests (`u-boot template list --json` **und** `u-boot --json
template list`) sind grün — gleicher Output, gleicher Exit-Code
wie heute.

Subcommand-Form-Inventar (Spec-Enum §338 zählt 10 Subcommands,
einige sind gruppiert):

- **Migrate in diesem Slice (2 Formen):** `doctor`, `template list`
  (Flag-Schnitt ohne Output-Migration — Carveouts-Eintrag).
- **Reject in diesem Slice (11 Formen):** `init`, `add`, `remove`,
  `up`, `down`, `logs`, `generate`, `config` (bare), `config get`,
  `config set`, `template` (bare). `config`/`template` als
  gruppierte Befehle (Spec §338 + §420) tragen jeweils mehrere
  Forms; bare `config --json` und bare `template --json` sind
  Help-Parents (`Args: cobra.NoArgs`, kein eigenes JSON-vorge-
  sehenes Output) und rejecten bis zum jeweiligen Folge-Slice
  (`slice-v1-cli-json-dry-run-config`/`-template`), der für
  bare die `subcommand`-Pflicht aus `LH-FA-CLI-007` §420 löst.

Reject-Form heute (alle 11): Exit-Code `2`
(`LH-FA-CLI-006`-Klasse) und Verweis auf den jeweiligen
Folge-Slice — kein Subcommand „akzeptiert `--json` still und
liefert Human-Output" (sonst untergrabener V1-Maschinenvertrag).
Cluster-T_close-Pflicht-Check: alle 13 Formen (2 Migrate +
11 Reject) sind in der Allowlist **oder** die Allowlist-
Mechanik komplett entfernt.

Konkrete Output-Pins für `doctor --json`:

```bash
# All-OK-Fall — diagnostics ist leer (level "ok" ist Spec-invalid)
u-boot doctor --json
# {
#   "status": "ok",
#   "command": "doctor",
#   "diagnostics": [],
#   "exitCode": 0
# }

# Warn-Fall — nur Warn-/Error-Items erscheinen in diagnostics
u-boot doctor --json
# {
#   "status": "warn",
#   "command": "doctor",
#   "diagnostics": [
#     {"level":"warn","code":"devcontainer.json.valid","message":"…",
#      "file":".devcontainer/devcontainer.json"}
#   ],
#   "exitCode": 0
# }

# Error-Fall — status koppelt an höchsten level, exitCode == 11
u-boot doctor --json
# {
#   "status": "error",
#   "command": "doctor",
#   "diagnostics": [
#     {"level":"error","code":"docker.installed","message":"…"}
#   ],
#   "exitCode": 11
# }
```

Pin-Werte für `diagnostics[].code` spiegeln den **T0-(h)-
Vorschlag Option (3)** (Tool-interne Codes aus
`domain.Diagnostic.ID`, dokumentiert in der Code-Registry).
Wird T0-(h) auf Option (1) LH-ID-Mapping festgezurrt, werden
die Codes durch `LH-FA-DIAG-002` / `LH-FA-DIAG-003` ersetzt
— die Pin-Tests müssen dann gegen die finale T0-(h)-Wahl
geschnitten werden.

`status` ist an das höchste `diagnostics[].level` gekoppelt
(`error → "error"`, `warn ohne error → "warn"`, sonst `"ok"`,
Lastenheft §447), analog zu der heutigen `HasErrors()`/
`HasWarnings()`-Logik in
[`doctor.go`](../../../../internal/adapter/driving/cli/doctor.go).
`SeverityOK`- **und** `SeverityInfo`-Items des heutigen Plaintext-
Reports werden im JSON-Modus **nicht** als `diagnostics[]`-
Eintrag serialisiert (Spec §1834 lässt `level` nur `warn|error`
zu — `domain.Severity` trägt seit M6-T1 vier Stufen
[`diagnostic.go:14-36`](../../../../internal/hexagon/domain/diagnostic.go),
heute emittiert kein Doctor-Check `SeverityInfo`, der Filter
ist Drift-Schutz für zukünftige Checks). `--quiet` ist im
JSON-Modus darum semantisch ein No-op.

## Akzeptanzkriterien

- ✅ **Root-PersistentFlag `--json` (Cluster T0-(a))**: am
  Cobra-Root registriert, persistent für alle 10 Subcommands,
  Wiring analog
  [`--verbose`/`--debug`/`--quiet`](../../../../internal/adapter/driving/cli/root.go).
- ✅ **Reject-Pfad für nicht-migrierte Subcommand-Formen**:
  alle **11** noch nicht migrierten Subcommand-Formen (`init`,
  `add`, `remove`, `up`, `down`, `logs`, `generate`, `config`,
  `config get`, `config set`, `template` (bare) — `config`
  zählt als 3 Formen, `template` als 1 bare-Form plus
  Sonderregel `template list` unten) rejecten `--json` mit
  Exit-Code `2` (`LH-FA-CLI-006`-Klasse) und Fehlermeldung
  `JSON-Ausgabe für 'u-boot <sub>' ist noch nicht implementiert
  (siehe slice-v1-cli-json-dry-run-<sub>).` Mechanik T0-(g).
  Pflicht-Pin-Test je nicht-migrierter Subcommand-Form (11
  Pin-Tests insgesamt). Pro Folge-Slice-Merge fallen **genau**
  die Reject-Pfade der migrierten Subcommand-Formen
  (`config`-Slice schließt drei auf einmal, `template`-Slice
  schließt bare-`template` und führt die Envelope-Migration
  für `template list`); **Cluster-T_close-Pflicht-Check:**
  alle 13 Formen (2 in diesem Slice migriert + 11 sukzessive
  via Folge-Slices) sind in der Allowlist enthalten **oder**
  Allowlist-Mechanik komplett entfernt — null offene
  Reject-Pfade.
- ✅ **Common-Envelope `cliJSONEnvelope` (Cluster T0-(c))**: neuer
  Typ im CLI-Adapter (Lokation-Sub-Decision T0-(a) dieses Slices,
  Vorschlag `internal/adapter/driving/cli/jsonenvelope.go`),
  trägt die Minimalkontrakt-Felder **immer** (`Status`,
  `Command`, `Subcommand omitempty`, `Diagnostics`, `ExitCode`,
  `Data omitempty`) und die Voll-Schema-Felder **omitempty**
  (`DryRun`, `Diff`, `PlannedFiles`, `Changes` — befüllt nur
  beim `add`-/modifying-Pfad). Read-only-Output trägt die
  Voll-Schema-Felder gar nicht im JSON (Lastenheft §1841,
  Minimalkontrakt-bindend); modifying-Output muss sie alle
  vier explizit setzen (`LH-FA-CLI-007` §326).
- ✅ **Schema-Helper-Modus-Split**: zwei Helper-Modi im
  Sub-Package `internal/adapter/driving/cli/jsontestutil/`:
  - `AssertMinimalEnvelope(t, raw)` prüft den Lastenheft-§1841-
    Minimalkontrakt — Pflicht-Set, `status` ∈ `{ok, warn, error}`,
    `diagnostics[].level` ∈ `{warn, error}` (rejected **sowohl**
    `"ok"` **als auch** `"info"` — `domain.Severity` trägt vier
    Stufen, Spec §1834 nur zwei zulässige `level`-Werte),
    `diagnostics[].code` LH-Kennung-konform **oder** in der
    Code-Registry (siehe AK Code-Registry), `exitCode` ≥ 0,
    `subcommand` nur bei `command ∈ {template, config}`.
    **Verbietet** `dryRun`/`diff`/`plannedFiles`/`changes` im
    Envelope (sonst ist es kein Minimalkontrakt mehr).
  - `AssertFullEnvelope(t, raw)` prüft zusätzlich die
    `LH-FA-CLI-007`-Voll-Felder. Wird in diesem Slice
    angelegt, aber **nicht** verwendet — Pattern-Anker für
    den nachfolgenden `add`-Slice.

  Wird ab diesem Slice von **jedem** Folge-Slice in seinen
  Schema-Tests aufgerufen (read-only → Minimal, modifying → Full).
- ✅ **`doctor --json` Envelope**: read-only-Pfad in
  [`runDoctor`](../../../../internal/adapter/driving/cli/doctor.go)
  oder einem geschwisterlichen `runDoctorJSON`. `diagnostics[]`
  wird aus
  [`domain.DiagnosticReport`](../../../../internal/hexagon/domain/diagnostic.go)
  gemappt, **wobei `SeverityOK`- *und* `SeverityInfo`-Items
  übersprungen werden** (Spec §1834 lässt `level` nur
  `warn|error` zu — `domain.Severity` trägt seit M6-T1 vier
  Stufen, der Filter pinnt sich gegen Drift wenn ein
  zukünftiger Doctor-Check `SeverityInfo` emittiert). Sort-
  Order analog
  [`SortedByIssuesFirst`](../../../../internal/hexagon/domain/diagnostic.go),
  aber gefiltert auf Warn/Error. `--quiet` und `--strict`
  interagieren mit `--json` definiert (Sub-Decision T0-(e)
  dieses Slices).
- ✅ **Code-Registry-Doku (Spec §1835 + §445)**: in
  `docs/user/cli-json-output.md` eine Sektion „Code-Registry"
  pflichtbestandteil. Wenn T0-(h) `Diagnostic.ID`-basierte
  Tool-interne Codes wählt (z. B. `docker.available`), trägt
  diese Sektion jeden Code mit seiner Bedeutung; wenn T0-(h)
  LH-ID-Mapping wählt (z. B. alle Doctor-Checks → `LH-FA-DIAG-002`
  oder `LH-FA-DIAG-003`), dokumentiert die Sektion das Mapping.
  Der Minimal-Helper konsultiert diese Sektion über eine
  embedded Allowlist (Sub-Decision T0-(b)/(h)).
- ✅ **`template list --json`-Schnitt (Cluster T0-(e) §M3)**:
  lokales `cmd.Flags().Bool("json", …)` auf
  [`template list`](../../../../internal/adapter/driving/cli/template.go)
  entfernt. `runTemplateList` liest Root-Flag über die
  App-Struktur oder `cmd.Root().PersistentFlags().GetBool("json")`.
  **Output-Format unverändert** (heutige `templateJSON`-Array-
  Struktur, `[]`-Normalisierung) — das ist Cluster-T0-(e)-
  Vorgabe, **bewusst nicht** der Minimalkontrakt. Bestehender
  Pin-Test bleibt grün, zusätzlicher Pin-Test für
  `u-boot --json template list` beweist identisches Verhalten.
  **Bekanntes Übergangs-Spec-Loch** (Cluster-Verantwortung):
  Der heutige Array-Output ist **nicht** Minimalkontrakt-konform.
  Der Reject-Pfad oben gilt deshalb **nicht** für `template
  list --json` (würde sonst bestehenden Output brechen). Die
  Minimal-konforme Form folgt mit Cluster-Platz 9
  `slice-v1-cli-json-dry-run-template`. **Carveouts-Eintrag-
  Pflicht**: dieser Slice trägt einen Carveout-Eintrag in
  [`carveouts.md`](carveouts.md) §Temporäre
  Carveouts mit dem Re-Trigger-Verweis auf
  `slice-v1-cli-json-dry-run-template`.
- ✅ **Schema-Konformität via Helper**: drei Acceptance-Tests
  pinnen `doctor --json` (All-OK-Fall mit `diagnostics: []`,
  Warn-Fall, Error-Fall) via `jsontestutil.AssertMinimalEnvelope`.
  Exit-Code konsistent mit `LH-FA-CLI-006` (`0` für ok / warn,
  `11` für error / strict-warn — gleicher Pfad wie heute,
  gemeinsamer `ErrDoctorFailures`-Anker).
- ✅ **Architektur-Grenzen sauber**: `cliJSONEnvelope` und
  Helper leben im CLI-Adapter, **nicht** im Domain- oder
  Application-Layer (`LH-FA-ARCH-002`/`003`). `make lint`
  (depguard) grün; CLI-Adapter importiert keine neuen
  `hexagon/port/driven`-Typen (RecordingFileSystem kommt erst
  im `add`-Folge-Slice).
- ✅ **Schema-Vertrag-Doku** (Cluster T1, in diesem Slice
  parallel mitliefern oder eigener Cluster-T1-Schritt
  vorab — Sub-Decision T0-(f) dieses Slices): zentraler
  Reference-Block (Kandidat `docs/user/cli-json-output.md`)
  zitiert sowohl den `LH-NFA-USE-004`-Minimalkontrakt
  (§1823-1842) als auch das `LH-FA-CLI-007`-Voll-Schema
  (§322-417) verbatim und benennt Envelope-Lokation plus
  Code-Registry; README EN+DE bekommen Verweis-Zeile.

## T0-Discovery (vor `next/`-Übergang)

Sub-Decisions, die dieser Folge-Slice klären muss, bevor er in
`next/` wandert. Layout analog zum Cluster-Slice §T0-Outcomes.

### T0-(a) Lokation des Envelope-Types

`internal/adapter/driving/cli/jsonenvelope.go` oder ein
Sub-Package `internal/adapter/driving/cli/jsonenvelope/`?
Erstere Variante ist leichter, letztere isoliert besser für
Test-Helper-Reuse. Sub-Decision mit LOC-Auswirkung.

Vorschlag: **Single-File** `internal/adapter/driving/cli/jsonenvelope.go`.
Der Envelope-Type ist kompakt (~10 Felder + JSON-Tags); ein
Sub-Package würde leeren `package`-Wrapper kosten ohne Symbol-
Reuse-Gewinn. Test-Helper-Reuse wandert in das parallele
Sub-Package `jsontestutil/` aus T0-(b) — die beiden Belange
(Wire-Type vs. Test-Helper) leben architektur-sauber getrennt,
ohne dass der Envelope-Type selbst ein eigenes Package braucht.

### T0-(b) Lokation und Form des Test-Helpers

`internal/adapter/driving/cli/jsontestutil/` mit zwei Modi
`AssertMinimalEnvelope(t, raw []byte, opts ...)` und
`AssertFullEnvelope(t, raw []byte, opts ...)` als Public-API.
Sub-Decision: ob der Helper die Schema-Regeln als embedded
JSON-Schema (Lastenheft §322-417 verbatim) oder als Go-Code-
Regeln trägt. Embed wäre langfristig sauberer gegen Drift,
kostet aber `embed`-Wiring und einen JSON-Schema-Validator —
Dep-Policy-Hinweis: `go.mod` 4-Dep-Disziplin, Sub-Decision
dieses Slices ob `github.com/santhosh-tekuri/jsonschema/v6`
o. ä. zulässig. Vorschlag: Go-Code-Regeln für V1 (näher am
heutigen Repo-Stil, kein zusätzlicher Dep). Allowlist für
`diagnostics[].code` (LH-IDs vs. dokumentierte Tool-Codes)
hängt von T0-(h) ab.

### T0-(c) `Data`-Inhalt für `doctor --json`

`Data` trägt im Envelope den subcommand-spezifischen Payload.
Für `doctor` Optionen:

1. Komplettes Mirroring von `domain.DiagnosticReport` (alle
   Items mit ID/Severity/Message/Hint).
2. Normalisierte Mini-Form (z. B. `{ "byID": { "docker.available":
   "ok", … } }` für schnellen Lookup).
3. **Kein** `Data`-Field für `doctor` — alle Informationen
   stecken im Pflicht-`diagnostics[]`-Array, `data` wäre Duplikation.

Vorschlag: Option 3 (Diagnostics-Pflichtfeld trägt schon alles).
Sub-Decision verbindlich machen.

### T0-(d) Minimalkontrakt-vs.-Voll-Schema-Disziplin im Envelope

Lastenheft §1841 macht den Minimalkontrakt für `--json` ohne
`--dry-run`/`--diff` **bindend**, das Voll-Schema gilt nur für
`--dry-run --json` / `--diff --json` (§1842, §468). Konsequenz
für `cliJSONEnvelope`: die Voll-Schema-Felder
`dryRun`/`diff`/`plannedFiles`/`changes` **müssen** im
read-only-Output **fehlen** (sonst läuft `doctor --json`
falsche Output-Form, und der Minimal-Helper kann sie nicht
sauber detektieren).

Optionen für die Go-Struktur:

1. Ein einziger `cliJSONEnvelope`-Typ mit
   `omitempty`-Tags auf den Voll-Feldern; read-only-Pfad lässt
   sie unbefüllt → werden weggelassen. Modifying-Pfad setzt
   alle vier explizit (auch `false`/`[]`).
2. Zwei separate Typen `minimalEnvelope` + `fullEnvelope`,
   beide tragen die Minimal-Felder, der zweite zusätzlich die
   Voll-Felder.

Vorschlag: Option 1 (ein Typ, `omitempty` auf den vier
Voll-Feldern; der Helper-Split aus T0-(b) garantiert, dass im
read-only-Pfad geprüft wird, dass die Felder im JSON-Output
fehlen, und im modifying-Pfad, dass sie alle vier gesetzt sind).

### T0-(e) Interaktion `--json` × `--quiet` × `--strict`

`--quiet` heute unterdrückt OK-Items im Plaintext-Output.
Im `--json`-Modus: `--quiet` ignorieren (JSON ist
Maschinen-Output, Konsumenten filtern selbst) oder
`--quiet`-Filter auf `diagnostics[]` anwenden?
`--strict` bleibt Exit-Code-Modifier (warn → 11), das ist
JSON-unabhängig. Sub-Decision: `--quiet --json` ignoriert
`--quiet` (Vorschlag).

### T0-(f) Schema-Vertrag-Doku-Zeitpunkt (Cluster T1)

Der Cluster-Slice T1 (`docs/user/cli-json-output.md`) ist als
„reine Doku" gelistet. Dieser Slice braucht den Schema-Anker zum
Verbatim-Zitieren. Sub-Decision: Cluster-T1 wird in diesem
Folge-Slice mitgeliefert (ein Stub-Doc, das später um weitere
Subcommand-Sektionen wächst) oder als eigener Vor-Schritt
parallel zum `next/`-Übergang.

Vorschlag: in diesem Folge-Slice mitliefern — sonst kein
Verbatim-Anker für den `jsontestutil`-Helper.

### T0-(g) Reject-Mechanik für nicht-migrierte Subcommands

Wenn das Root-PersistentFlag `--json` ab T3 für alle 10
Spec-Enum-Subcommands existiert, müssen die 10 noch nicht
migrierten Subcommand-Formen ihn explizit rejecten
(`config` zählt als 3 Formen — `config`, `config get`,
`config set`; `template list` ist Sonderfall mit Flag-Schnitt
ohne Reject; vgl. Aufhebungsbedingung + AK Reject-Pfad). Drei
Mechanik-Optionen:

1. **PreRunE pro Subcommand**: jeder noch nicht migrierte
   Subcommand bekommt einen `PreRunE`-Hook, der den Root-Flag
   abfragt und bei `--json` mit Exit-Code `2` abbricht.
   Vorteil: lokal, leicht zu entfernen pro Folge-Slice-Merge.
   Nachteil: 10 separate Wirings, einfach zu vergessen.
2. **Zentrale Subcommand-Allowlist** in `cli/root.go`: eine
   Map `migratedJSONSubcommands` listet die migrierten
   Subcommand-Formen. Cobra-`PersistentPreRunE` am Root prüft,
   ob `--json` gesetzt und Subcommand nicht in Allowlist →
   Reject. Pro Folge-Slice-Merge wandern **die Einträge der
   migrierten Formen** in die Allowlist (`config`-Slice
   schließt drei auf einmal). Vorteil: zentraler Drift-Anker;
   Nachteil: zentrales File ändert sich oft. Cluster-T_close-
   Check trivial (Allowlist muss alle 11 Formen — die 10
   gerejekteten plus `template list` — enthalten; bzw.
   Allowlist und Reject-Mechanik können dann ganz entfernt
   werden).
3. **Reject im Envelope-Builder**: wenn ein Subcommand keinen
   Envelope-Builder verlinkt hat, fällt der Build durch ein
   `nil`-Check und gibt Exit-Code 2. Vorteil: implicit;
   Nachteil: schwer zu testen, leise Fehler-Klasse.

Vorschlag: Option (2) — zentrale Allowlist in `cli/root.go`
mit `PersistentPreRunE`. Die Allowlist wandert mit jedem
Folge-Slice ein Stück nach vorne; Cluster-T_close entfernt
Allowlist + Reject-Pfad komplett.

### T0-(h) `diagnostics[].code`-Quelle vs. Spec-Code-Konvention

Heute trägt
[`domain.Diagnostic.ID`](../../../../internal/hexagon/domain/diagnostic.go)
einen stabilen Check-Identifier wie `"docker.available"`. Spec
§1835 + §445 erlaubt für `diagnostics[].code` **zwei** Formen:
(a) LH-Kennung der verursachenden Anforderung (z. B.
`LH-FA-DIAG-002`, `LH-FA-DEV-003`) oder (b) Tool-interne Codes,
**falls** ihre Bedeutung in der Doku festgehalten ist.
`docker.available` ist tool-intern → Option (b) braucht eine
Code-Registry in `docs/user/cli-json-output.md`.

Drei Sub-Optionen:

1. **LH-ID-Mapping in der Adapter-Schicht**: CLI-Adapter
   übersetzt `Diagnostic.ID → LH-ID` (z. B. alle Doctor-Checks
   gemeinsam → `LH-FA-DIAG-002` oder `LH-FA-DIAG-003`). Vorteil:
   keine Domain-Änderung; Nachteil: Mapping-Logik im Adapter,
   Mehrere Checks teilen eine LH-ID (verliert Stabilität für
   CI-Scrapings).
2. **Domain-Erweiterung `Diagnostic.LHCode string`**:
   `Diagnostic`-Struct um optionales Feld erweitern; jeder
   Check-Konstruktor setzt eine LH-ID. Vorteil: pro-Check-
   Spezifik; Nachteil: rückwirkend alle Doctor-Checks
   anpassen — größerer Eingriff, Domain-Layer-Berührung
   (`LH-FA-ARCH-002`-OK, aber spec-getrieben).
3. **Tool-interne Codes plus Registry-Doku**: heutige
   `Diagnostic.ID` direkt durchreichen, Code-Registry-Sektion
   in `docs/user/cli-json-output.md` dokumentiert die Bedeutung
   jeder Tool-ID. Vorteil: minimaler Eingriff; Nachteil: jede
   neue Check-ID muss gleichzeitig im Doc auftauchen — Drift-
   Risiko (Pflicht-Check im Schema-Helper).

Vorschlag: Option (3) für V1 — pragmatisch, kein Domain-
Refactor, der Drift-Schutz wandert in den Schema-Helper
(`AssertMinimalEnvelope` lädt die Registry und rejected
Codes außerhalb).

Spec-konsistente Folgepflicht: jeder Folge-Slice trägt **seine**
Checks in der Registry nach, der Helper rejected sonst.

## T0-Outcomes

Acht Sub-Decisions vor `next/`-Übergang festgezurrt — die acht
Fragen aus [§T0-Discovery](#t0-discovery-vor-next-übergang) mit
Begründung. Layout analog Cluster-Slice §T0-Outcomes.

### T0-(a) Envelope-Lokation: Single-File `cli/jsonenvelope.go`

**Entscheidung:** Neuer Wire-Type `cliJSONEnvelope` lebt als
**Single-File** in
`internal/adapter/driving/cli/jsonenvelope.go`. Kein eigenes
Sub-Package.

**Begründung:** Der Envelope-Type ist kompakt (Minimal-Felder
`Status`/`Command`/`Subcommand`/`Diagnostics`/`ExitCode` +
Voll-Felder `DryRun`/`Diff`/`PlannedFiles`/`Changes` +
`Data` — ~10 Felder mit `json`-Tags). Ein Sub-Package würde
einen leeren `package`-Wrapper kosten ohne Symbol-Reuse-Gewinn.
Test-Helper-Reuse lebt architektur-sauber im **parallelen**
Sub-Package `jsontestutil/` aus T0-(b) — Wire-Type-Heimat und
Test-Helper-Heimat sind bewusst getrennt, ohne dass der
Envelope-Type selbst ein eigenes Package braucht. Repo-Vorbild:
[`cli/template.go`](../../../../internal/adapter/driving/cli/template.go)
trägt `templateJSON` lokal — analoge Disziplin.

### T0-(b) Test-Helper-Form: Sub-Package `jsontestutil/` mit Go-Code-Regeln und Options-Pattern

**Entscheidung:** Neues Sub-Package
`internal/adapter/driving/cli/jsontestutil/` mit zwei Public-
API-Funktionen:

```go
func AssertMinimalEnvelope(t testing.TB, raw []byte, opts ...AssertOption)
func AssertFullEnvelope(t testing.TB, raw []byte, opts ...AssertOption)
```

`AssertOption`-Funktionen:

- `WithCommand(string)` — pinnt den erwarteten `command`-Wert
  (z. B. `"doctor"`); rejected `command`-Mismatch.
- `WithSubcommand(string)` — pinnt `subcommand`-Pflichtwert für
  `command ∈ {template, config}`.
- `WithExpectedCodes(...string)` — pinnt, **welche** Codes der
  konkrete Test-Output enthalten muss (Subset-Pin gegen den
  Helper-Default-Set aus T0-(h) `DefaultAllowedCodes`). **Keine**
  Allowlist-Erweiterung — die globale Code-Allowlist bleibt
  exklusiv `DefaultAllowedCodes`, und neue Subcommand-Codes
  müssen dort **plus** in der Markdown-Doku eingetragen werden
  (siehe T0-(h) Single-Source-of-Truth-Disziplin). Diese Option
  ist nur Pin-Hilfe für konkrete Test-Fälle, kein Schema-Bypass.
- `WithExitCode(int)` — pinnt erwarteten `exitCode`-Wert
  (häufiger Pin-Wunsch).

Schema-Regeln werden als **Go-Code** geprüft, **kein** embedded
JSON-Schema, **kein** zusätzlicher Dep.

**Begründung:** Die heutige `go.mod`-Disziplin (4 Deps total)
verbietet leichtfertige Dep-Additionen. Ein JSON-Schema-Validator
(`github.com/santhosh-tekuri/jsonschema/v6` o. ä.) würde die
Dep-Liste auf 5 anheben — der Drift-Schutz, den der Validator
bringt, wandert pragmatisch in den Helper-Code (`AssertMinimalEnvelope`
prüft die Lastenheft-§1841-Regeln direkt mit `t.Errorf`-Calls).
Falls die Code-Regeln später wirklich driften, ist der Wechsel
auf embedded Schema ein lokaler Refactor — keine Architektur-
Berührung in den Folge-Slices.

Options-Pattern statt Konstruktor-Vielfalt: jeder Folge-Slice
trägt unterschiedliche Pins (`add` braucht `WithCommand("add")`,
`config get` braucht `WithSubcommand("get")`), die Variadik
hält die Helper-API klein. **Bewusster Verzicht** auf eine
`WithAllowedCodes`-Option, die das globale `DefaultAllowedCodes`
erweitern würde — die globale Map ist exklusive Source-of-Truth
gemäß T0-(h), Folge-Slice-Erweiterungen passieren nur dort
(Single-Source-Disziplin).

### T0-(c) `Data`-Inhalt für `doctor --json`: kein `Data`-Field (Option 3)

**Entscheidung:** `doctor --json` setzt **kein** `Data`-Field
im Envelope. Alle Doctor-Informationen liegen im Pflicht-
`diagnostics[]`-Array.

**Begründung:** `diagnostics[]` trägt nach dem Mapping aus T5
bereits `level`/`code`/`message`/`file`/`hint` pro Check —
zusätzliche `Data`-Felder wären Duplikation. Das Envelope-
Struct trägt `Data any` mit `omitempty`-Tag (Voll-Schema-Field
für andere read-only-Subcommands wie `up`/`down`-Compose-Status,
`template list`-Array nach Platz-9-Migration). Bei `doctor`
bleibt es nil → nicht im JSON.

**Pattern-Disziplin für die 8 weiteren Folge-Slices (Review-Finding L2):**
`Data any` als Wire-Type-Field verliert jede Compile-Time-
Garantie. Pattern-Vorbild-Last dieses Slices: jeder Folge-Slice,
der `Data` setzt, definiert dafür einen **subcommand-spezifischen
Wire-Type** mit dedizierten `json`-Tags (z. B. `type
upStatusData struct { Services []serviceStatus \`json:"services"\` }`
für `slice-v1-cli-json-dry-run-up-down`) und übergibt eine
**typisierte Instanz** (`Data: upStatusData{...}`) — **kein**
`map[string]any`, **keine** Inline-Anonym-Structs. Acceptance-
Helper kann den Sub-Type via `json.RawMessage`-Re-Marshal pinnen,
falls strikter Schema-Pin gewünscht.

### T0-(d) Envelope-Struktur: ein Typ mit Pointer-Wrapping auf Voll-Boolean-Feldern (Option 1)

**Entscheidung:** Ein Typ `cliJSONEnvelope`, Modus-Unterscheidung
über Pointer/Slice-Nullbarkeit plus zwei Konstruktoren:

```go
type cliJSONEnvelope struct {
    Status       string           `json:"status"`
    Command      string           `json:"command"`
    Subcommand   string           `json:"subcommand,omitempty"`
    DryRun       *bool            `json:"dryRun,omitempty"`
    Diff         *bool            `json:"diff,omitempty"`
    PlannedFiles []plannedFile    `json:"plannedFiles,omitempty"`
    Changes      []changeEntry    `json:"changes,omitempty"`
    Diagnostics  []diagnosticItem `json:"diagnostics"`
    ExitCode     int              `json:"exitCode"`
    Data         any              `json:"data,omitempty"`
}

func newMinimalEnvelope(command, subcommand string, diags []diagnosticItem, exitCode int) cliJSONEnvelope
func newFullEnvelope(command, subcommand string, dryRun, diff bool, planned []plannedFile, changes []changeEntry, diags []diagnosticItem, exitCode int) cliJSONEnvelope
```

`Diagnostics` trägt **kein** `omitempty` — leeres Array muss
laut Spec §1833 als `[]` serialisiert werden, nicht weggelassen.
Im Read-only-Pfad wird `newMinimalEnvelope` verwendet → `DryRun`/
`Diff`-Pointer bleiben nil, `PlannedFiles`/`Changes`-Slices
bleiben nil → alle vier Felder fallen via `omitempty` aus dem
JSON. Im Modifying-Pfad setzt `newFullEnvelope` `*bool`-Pointer
auf `&dryRun`/`&diff` und Slices auf `[]plannedFile{...}` /
`[]changeEntry{...}` (auch leere als `make([]X, 0)`, damit der
JSON-Marshaller `[]` schreibt statt das Feld wegzulassen).

**Begründung:** Option 2 (zwei separate Wire-Typen) würde Symbol-
und Test-Helper-Duplikation erzeugen, ohne semantischen Gewinn
— der Spec-Schnitt aus §1841/§1842 ist klar genug, dass ein
Typ mit zwei Konstruktoren die Disziplin trägt. Die Pointer-Wahl
auf `*bool` ist Go-typisch (gleicher Trick wie z. B. `json.RawMessage`
für optionale Boolean-Felder) und macht die `omitempty`-Semantik
trivial. Der Helper-Split aus T0-(b) pinnt: `AssertMinimalEnvelope`
rejected `dryRun`/`diff`/`plannedFiles`/`changes` im JSON-Output
(durch reines String-Suchen oder `gjson`-freies Parsing als
`map[string]any` und Key-Check), `AssertFullEnvelope` verlangt
alle vier.

**Anti-Drift-Pin gegen `*bool → bool`-Refactor (Review-Finding M1):**
Ein späterer Folge-Slice könnte versucht sein, `*bool` durch
`bool` mit `omitempty` zu ersetzen ("sieht harmloser aus") —
das würde `dryRun:false`/`diff:false` aus dem JSON werfen und
Spec §326 Required-Set verletzen (`["status", "command",
"dryRun", "diff", "plannedFiles", "changes", "diagnostics",
"exitCode"]` — alle acht pflicht im modifying-Pfad, auch wenn
Wert `false`). Drei Schutz-Maßnahmen in T2:

1. **Struct-Kommentar am Envelope-Type:** `// IMPORTANT: *bool
   (not bool) — Spec §326 requires dryRun/diff in modifying-mode
   even when value is false. Plain bool + omitempty would drop
   false from JSON.`
2. **Positive-Marshal-Pin:** Unit-Test pinnt, dass
   `newFullEnvelope(..., dryRun=false, diff=false, ...)`
   die JSON-Repräsentation `"dryRun":false,"diff":false`
   produziert (nicht weggelassen).
3. **`AssertFullEnvelope`-Required-Check:** der Voll-Helper
   prüft beim ersten Acceptance-Lauf (im `add`-Folge-Slice)
   das Spec §326 Required-Set komplett — bricht, wenn eines
   der acht Felder im JSON fehlt.

### T0-(e) `--json` × `--quiet` × `--strict`: `--quiet` ist No-op, `--strict` bleibt Exit-Code-Modifier

**Entscheidung:** Im JSON-Modus ignoriert die Logik den
`--quiet`-Flag-State semantisch. `--strict` bleibt unverändert
ein Exit-Code-Modifier (Warn-Fund unter `--strict` → Exit-Code 11),
JSON-Body identisch zu `--json` ohne `--strict`.

**Begründung:** `--quiet` filtert heute OK-Items aus dem
Plaintext-Output. Im JSON-Modus sind OK-/Info-Items bereits
durch den Spec-§1834-Filter (`level ∈ {warn, error}`)
ausgeschlossen — `--quiet` hätte keinen zusätzlichen Filter-
Effekt mehr. Die Pin-Test-Disziplin in T5 erzwingt das:
`u-boot doctor --quiet --json` produziert einen **semantisch
identischen** Envelope wie `u-boot doctor --json` —
gleicher `status`, gleiche `diagnostics`-Reihenfolge mit
gleichen `code`/`level`-Paaren, gleicher `exitCode`. Die
Pin-Form ist bewusst **nicht** byte-identisch (Review-Finding
M4): sobald ein zukünftiger Doctor-Check zeitabhängige
Information in `message` einbettet (Daemon-Version, Plugin-
Version), wäre der byte-identische Vergleich brüchig und
würde zur Aufweichung verleiten. Die semantische Form trägt
den intendierten Drift-Schutz robuster und bleibt scharf
gegen versehentliche Quiet-Filter-Einführung.

`--strict` bleibt orthogonal: die Severity-→-Exit-Code-Logik in
[`doctor.go`](../../../../internal/adapter/driving/cli/doctor.go)
bleibt unverändert, T5 setzt nur das Pin: `u-boot doctor
--strict --json` mit Warn-Fund liefert `status: "warn"`,
`exitCode: 11` (statt `0` im non-strict-Pfad).

### T0-(f) Schema-Vertrag-Doku in T1 dieses Slices mitliefern

**Entscheidung:** Cluster-T1 (`docs/user/cli-json-output.md`)
wird als T1-Tranche dieses Slices mitgeliefert, nicht als
eigener Cluster-Vor-Schritt.

**Begründung:** Der Test-Helper aus T2 muss gegen den verbatim
zitierten Schema-Wortlaut entwickelt werden — ein vorgelagerter
eigener Cluster-T1-Schritt würde nur die Kosten einer separaten
Tranche-Plan-Pflege bringen, ohne Reihenfolge-Vorteil. T1 in
diesem Slice baut Minimal-Schema (§1823-1842) und Voll-Schema
(§322-417) als zwei klar getrennte Sektionen auf; spätere
Folge-Slices erweitern die Code-Registry-Sektion, lassen die
Schema-Sektionen unverändert.

**Cluster-T1-Doppelzählung vermeiden (Review-Finding L1):**
Der Cluster-Slice führt T1 noch in seiner eigenen Lieferpflicht-
Tabelle (`docs/plan/planning/in-progress/slice-v1-cli-json-dry-run.md`
§Tranchen T1). Pflicht-Nachzug **in diesem Slice-Closure (T6)**:
Cluster-Slice-T1-Zelle muss eine „geliefert via Doctor-Slice
(`slice-v1-cli-json-dry-run-doctor`)"-Notiz tragen, damit
Cluster-T_close-Reviewer den T1-Schritt nicht doppelt zählt
oder vergeblich nach einem separaten Cluster-T1-Commit sucht.

### T0-(g) Reject-Mechanik: zentrale Allowlist mit `PersistentPreRunE` (Option 2)

**Entscheidung:** Reject-Pfad lebt in `cli/root.go` als
`PersistentPreRunE` am Root-Command. Eine Allowlist-Map listet
die migrierten Subcommand-Formen per Cobra-`cmd.CommandPath()`:

```go
var jsonAllowlist = map[string]bool{
    "u-boot doctor":             true,  // dieser Slice
    "u-boot template list":      true,  // dieser Slice (Flag-Schnitt, Carveout)
    // weitere Einträge pro Folge-Slice-Merge
}
```

`PersistentPreRunE` prüft: wenn Root-Flag `--json` gesetzt und
`cmd.CommandPath()` nicht in `jsonAllowlist` → Reject mit
Exit-Code 2 und Fehlermeldung
`JSON-Ausgabe für '<cmd.CommandPath()>' ist noch nicht
implementiert (siehe slice-v1-cli-json-dry-run-<sub>).`

**Subcommand-Form-Inventar (13 Formen):**

| Spec-Enum (§338) | Forms | Heute |
| --- | --- | --- |
| `init` | 1 | Reject |
| `add` | 1 | Reject |
| `remove` | 1 | Reject |
| `up` | 1 | Reject |
| `down` | 1 | Reject |
| `doctor` | 1 | **Migrate** (dieser Slice) |
| `logs` | 1 | Reject |
| `generate` | 1 | Reject |
| `config` | 3 (bare, `get`, `set`) | Reject (alle 3) |
| `template` | 2 (bare, `list`) | bare: Reject; `list`: **Migrate** (Flag-Schnitt, Carveout) |

Heute (nach diesem Slice): 2 Migrate + 11 Reject = 13 Formen
in der Allowlist. **Cluster-T_close-Pflicht-Check:** Allowlist
enthält **alle 13** Formen (alle Reject-Pfade durch Folge-Slices
abgebaut), **dann** wird die Allowlist-Mechanik komplett
entfernt (`PersistentPreRunE` raus, `jsonAllowlist`-Map raus).
Pflicht-Pin im Cluster-Close-Slice.

**Bare-`template`-/bare-`config`-Klärung:** beide sind in
Cobra heute `Args: cobra.NoArgs` ohne `RunE` — sie drucken
Help, wenn ohne Subcommand aufgerufen. `u-boot config --json`
und `u-boot template --json` haben **kein** Spec-vorgesehenes
Output-Format (Spec §420 fordert `subcommand`-Pflicht bei
`command ∈ {template, config}`); sie rejecten daher bis zum
jeweiligen Folge-Slice (`slice-v1-cli-json-dry-run-config`/
`-template`), der die `subcommand`-Pflicht-Auflösung für
bare-Form trifft (Kandidaten gemäß Cluster-Plan §96-107:
`"list"`, `"show"`, explizit `""`).

**Anti-Drift-Pin gegen `cmd.Use`-Umbenennungen (Review M2):**
T2 trägt einen Cobra-Tree-Walk-Test, der über `rootCmd.Commands()`
rekursiv läuft und für jeden Leaf-Command (plus Help-Parent-
Commands wie bare `template`/`config`) den `cmd.CommandPath()`
gegen den erwarteten Map-Key-String matched. Bricht, wenn ein
späterer Slice `cmd.Use` ändert (z. B. Synonym oder Spec-Drift),
ohne den Map-Key mitzuziehen — der Test zeigt sofort, welcher
Pfad-String aus der Allowlist gelaufen ist.

**Begründung:** Option 1 (PreRunE pro Subcommand) würde 11
separate Wirings bedeuten — vergessens-anfällig, und der
Reviewer für jeden Folge-Slice müsste die korrekte Entfernung
des PreRunE-Hooks prüfen. Option 3 (Reject im Envelope-Builder)
ist implizit/leise und schwer testbar. Option 2 ist der zentrale
Drift-Anker: ein einziger Map-Eintrag pro Folge-Slice, ein
einziger Test-Code-Pfad. Map-Key = `cmd.CommandPath()` ist
Cobra-idiomatisch (vgl.
[`cmd.CommandPath()`](https://pkg.go.dev/github.com/spf13/cobra#Command.CommandPath))
und liefert für `config get` den vollen Pfad `u-boot config get`
— die drei `config`-Formen sind damit individuell adressierbar
und werden vom `config`-Folge-Slice in einem Merge zu drei
Allowlist-Einträgen. Der Cobra-Tree-Walk-Anti-Drift-Pin schließt
das Drift-Restrisiko gegen `Use`-Renames.

### T0-(h) `diagnostics[].code`-Quelle: Tool-interne Codes + Registry (Option 3)

**Entscheidung:** `diagnostics[].code` trägt für Doctor die
heutigen tool-internen `Diagnostic.ID`-Werte (z. B.
`"docker.installed"`, `"docker.reachable"`,
`"uboot.yaml.valid"`) unverändert durch. Eine Code-Registry
in `docs/user/cli-json-output.md` §Code-Registry dokumentiert
die Bedeutung jedes Codes. Eine Go-Map in
`internal/adapter/driving/cli/jsontestutil/coderegistry.go` ist
die Source-of-Truth für den Helper:

```go
var DefaultAllowedCodes = map[string]string{
    "fs.write-permissions":            "doctor: Schreib-Permission im Working Directory",
    "git.installed":                   "doctor: Git-Binary verfügbar",
    "docker.installed":                "doctor: Docker-Binary verfügbar",
    "docker.reachable":                "doctor: Docker-Daemon erreichbar",
    "docker.compose.installed":        "doctor: Compose-Plugin verfügbar",
    "uboot.yaml.valid":                "doctor: u-boot.yaml syntaktisch valide",
    "compose.yaml.valid":              "doctor: compose.yaml syntaktisch valide",
    "devcontainer.json.valid":         "doctor: devcontainer.json syntaktisch valide",
    "devcontainer.dockerfile.valid":   "doctor: devcontainer/Dockerfile parsebar",
    "services.enabled-key":            "doctor: u-boot.yaml services-Block konsistent",
    "devcontainer.forwardPorts.consistency": "doctor: devcontainer.json forwardPorts konsistent",
    "devcontainer.features.allowlist": "doctor: devcontainer features auf Allowlist",
    "devcontainer.features.drift":     "doctor: devcontainer features ohne Drift",
}
```

`AssertMinimalEnvelope` lädt diese Map plus optionale
`WithAllowedCodes(...)`-Erweiterungen aus T0-(b) und rejected
`diagnostics[].code`-Werte außerhalb.

**Drift-Disziplin (drei aktive Drift-Gates, Review-Finding H2):**
Die Go-Map ist die Source-of-Truth; die Markdown-Sektion ist
spec-pflichtige Doku (§1835). Drei automatische Gates in T2:

1. **Map ↔ Code-Realität-Pin:** Unit-Test prüft, dass jede
   `checkID*`-Konstante aus
   [`doctor.go:74-114`](../../../../internal/hexagon/application/doctor.go)
   einen Map-Eintrag hat. Bricht, sobald ein neuer Doctor-Check
   im Code ohne Map-Eintrag landet.

2. **Map ↔ Markdown-Roundtrip-Pin:** Unit-Test parst die
   Code-Registry-Sektion in `docs/user/cli-json-output.md` mit
   einem schlanken Markdown-Tabellen-Regex
   (`^\|\s*\` ``…`` `\s*\|\s*(.+?)\s*\|$` o. ä., ~20 LOC, keine
   neue Dep) und vergleicht die Code-Set-Symmetrie gegen die
   Go-Map: jeder Map-Key hat eine Markdown-Zeile mit identischer
   Beschreibung; jede Markdown-Zeile hat einen Map-Eintrag.
   Bricht in **beide** Drift-Richtungen.

3. **Helper-Reject im Acceptance-Pfad:** `AssertMinimalEnvelope`
   rejected jeden `diagnostics[].code`, der nicht in
   `DefaultAllowedCodes` steht (Folge-Slice-Acceptance-Test
   schlägt sofort fehl, **nicht** erst im Cluster-T_close).

Gate 1+2 laufen in `make test`; Gate 3 läuft per Subcommand-
Acceptance-Test. **Kein** `//go:generate`-Pfad — der Roundtrip-
Test ist kürzer als ein Generator (~20 LOC vs. ~30 LOC plus
`make lint`-Wiring) und gibt sofort einen Failure mit Zeile statt
einem `git diff`-Stub im CI.

Spec-konsistente Folgepflicht für jeden Folge-Slice: seine
Check-IDs landen in der Map **und** in der Doku, im selben
Slice-Closure-Commit — sonst bricht Gate 2 (Markdown-Roundtrip)
oder Gate 3 (Helper-Reject) im selben PR-Lauf.

## Tranchen (vorgeschlagen)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + Sub-Decisions** aus §T0-Discovery klären (acht Sub-Decisions, inkl. T0-(g) Reject-Mechanik und T0-(h) `diagnostics[].code`-Quelle). Entscheidungen mit Begründung in einem `T0-Outcomes`-Block dokumentieren. | — (Plan-Arbeit) |
| T1 | **Schema-Vertrag-Doku.** `docs/user/cli-json-output.md` anlegen: Minimalkontrakt (`LH-NFA-USE-004` §1823-1842) und Voll-Schema (`LH-FA-CLI-007` §322-417) verbatim getrennt zitiert; **Code-Registry-Sektion** für Doctor-Checks (gemäß T0-(h)); Envelope-Lokation benennen; Minimal-vs.-Voll-Diff klargestellt. README EN+DE bekommt einen Verweis-Eintrag. (Cluster-T1; per T0-(f) hier mitgeliefert.) | ~120 |
| T2 | **Envelope + Helper-Split + Drift-Gates.** `cliJSONEnvelope`-Typ (Minimal-Felder Pflicht, Voll-Felder `omitempty` via `*bool`/nil-Slices — gemäß T0-(d), inkl. Anti-Drift-Struct-Kommentar gegen `*bool → bool`-Refactor) plus zwei Helper `AssertMinimalEnvelope` und `AssertFullEnvelope` mit Options (`WithCommand`, `WithSubcommand`, `WithExpectedCodes`, `WithExitCode`). **Pin-Tests:** (i) Helper positive/negative (fehlende Pflichtfelder, ungültiges `status`, `level: "ok"`/`level: "info"`-Reject, Voll-Schema-Feld im Minimal-Pfad-Reject, undokumentierter Code-Reject); (ii) Marshal-Pin `newFullEnvelope(..., dryRun=false, diff=false, ...)` produziert `"dryRun":false,"diff":false` (nicht weggelassen — M1-Schutz); (iii) **Drift-Gate 1** Map ↔ `doctor.go:74-114`-Vollständigkeit; (iv) **Drift-Gate 2** Map ↔ Markdown-Roundtrip-Parser. | ~220 |
| T3 | **Root-PersistentFlag `--json` + Reject-Allowlist + Cobra-Tree-Walk-Pin** am Cobra-Root registrieren plus `PersistentPreRunE`-Reject für nicht-migrierte Subcommand-Formen (Mechanik gemäß T0-(g)). **11 Reject-Pin-Tests**, je einer pro nicht-migrierter Form (`init`, `add`, `remove`, `up`, `down`, `logs`, `generate`, `config`, `config get`, `config set`, `template` (bare); `template list` siehe T4 — kein Reject). **Anti-Drift-Pin** Cobra-Tree-Walk: rekursiver `rootCmd.Commands()`-Traversal vergleicht jeden Leaf-Path mit erwartetem Map-Key (M2-Schutz gegen `cmd.Use`-Renames). App-Struktur-Field für den Flag-State. | ~110 |
| T4 | **`template list`-Schnitt + Carveouts-Eintrag.** Lokales Flag entfernen, `runTemplateList` liest Root-Flag. Zwei Pin-Tests grün (`u-boot template list --json` + `u-boot --json template list` → gleicher Output). Carveouts-Eintrag in `carveouts.md` §Temporäre Carveouts: `template list --json` heute nicht Minimalkontrakt-konform, Re-Trigger `slice-v1-cli-json-dry-run-template`. | ~40 |
| T5 | **`doctor --json`-Pfad.** Envelope-Befüllung aus `domain.DiagnosticReport` mit `SeverityOK`- + `SeverityInfo`-Filter (kein `level: "ok"`/`level: "info"`-Eintrag), `status`-Mapping (höchstes vorhandenes Severity-Level), `exitCode` konsistent mit `ErrDoctorFailures`. Drei Acceptance-Tests (All-OK mit `diagnostics: []`, Warn-Fall, Error-Fall) via `jsontestutil.AssertMinimalEnvelope`. `--quiet`/`--strict`-Interaktion gemäß T0-(e), **inkl. Pin-Test `u-boot doctor --quiet --json` liefert semantisch identischen Envelope wie `u-boot doctor --json`** (gleiche `status`/`exitCode`, gleiche `diagnostics`-Reihenfolge mit gleichen `code`/`level`-Paaren — bewusst nicht byte-identisch, M4-Schutz gegen brüchige Pins bei zeitabhängigen Messages). Doctor-Check-Codes in Code-Registry aus T1 ergänzt. | ~140 |
| T6 | **Closure.** CHANGELOG `## [Unreleased]` Added-Eintrag (Envelope + Helper-Split + Drift-Gates + Root-Flag + Reject-Allowlist + Cobra-Tree-Walk-Pin + `doctor --json` + `template list`-Flag-Schnitt + Carveout). **roadmap.md v0.4.0-Tabelle:** Cluster-Slice-Zelle aktualisiert (Doctor done, nächster Schritt `add`). **Cluster-Slice-T1-Zelle** in `docs/plan/planning/in-progress/slice-v1-cli-json-dry-run.md` mit „geliefert via `slice-v1-cli-json-dry-run-doctor`"-Notiz versehen (L1-Konsistenz). Slice-File `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. | — (Doku) |

LOC-Schätzung Folge-Slice: ~630 LOC — leicht über der vom
Cluster-Slice gesetzten 200..600-Bandbreite. Treiber des Anstiegs
gegenüber dem ursprünglichen Stub (~560): die nach dem T0-Review
ergänzten Drift-Gates (Map↔Markdown-Roundtrip-Parser, Cobra-Tree-
Walk-Pin, Marshal-Pin) plus 11 statt 10 Reject-Pin-Tests. Die
Bandbreiten-Überschreitung ist begründet (Pattern-Vorbild-Last für
8 weitere Folge-Slices); Cluster-Plan §LOC-Bandbreite ist deshalb
**keine** Hard-Rule.

## Out of Scope

- **`--dry-run` / `--diff` für `doctor`**: `doctor` ist
  read-only, schreibt nichts. `LH-FA-CLI-007`/`-008` Voll-Schema
  gelten nicht — `doctor --json` erfüllt den Minimalkontrakt
  aus `LH-NFA-USE-004` §1841 (T0-(d)).
- **`AssertFullEnvelope`-Nutzung in diesem Slice**: der
  Voll-Helper wird hier nur als Stub angelegt, **nicht** verwendet
  (keine modifying-Tests). Erstnutzung im Folge-Slice
  `slice-v1-cli-json-dry-run-add`.
- **`template list`-Envelope-Migration**: dieser Slice macht
  nur den Flag-Schnitt (lokal → Root), Output bleibt heutiges
  Array. Die Envelope-Migration mit Minimalkontrakt landet
  im Cluster-Platz-9 `slice-v1-cli-json-dry-run-template`
  (Carveouts-Eintrag aus T4 ist der Re-Trigger).
- **`RecordingFileSystem` und alles aus Cluster-T0-(b)**:
  kommt erst im Folge-Slice `slice-v1-cli-json-dry-run-add`.
  Dieser Slice **importiert keinen** neuen `driven.FileSystem`-
  Typ und etabliert kein Composition-Root-Wiring für
  Preview-Adapter. Konkret:
  [`cmd/uboot/main.go`](../../../../cmd/uboot/main.go)
  wird **nicht angefasst** — keine Doppel-Wiring der driving-
  Port-Instanzen (Normal-Mode + Preview-Mode), keine
  `RecordingFileSystem`-Instanziierung, keine App-Struktur-
  Erweiterung um Preview-Use-Case-Felder. Der einzige
  Composition-Root-Touch dieses Slices ist die Root-Persistent-
  Flag `--json` in
  [`cli/root.go`](../../../../internal/adapter/driving/cli/root.go)
  (analog zum heutigen `--verbose`/`--debug`/`--quiet`-Wiring).
  Der `add`-Folge-Slice erbt das Composition-Root-Doppel-Wiring
  als geschlossenen Cluster-T0-(b)-Outcome-Block.
- **Envelope-Migration für `template list`-Output**: nur der
  Flag-Schnitt wandert hierher (T0-(e)-Pflicht aus dem Cluster).
  Die Migration der `templateJSON`-Array-Form auf den
  Envelope-`Data`-Field passiert in Platz 9
  `slice-v1-cli-json-dry-run-template`.
- **JSON-Output für alle anderen 8 Subcommands**: jeder eigener
  Folge-Slice gemäß Cluster-T0-(e)-Reihenfolge.
- **Neue Folge-ADR** „JSON-CLI ausgeliefert": entscheidet
  Cluster-T_close, nicht dieser Slice. ADR-0010 bleibt
  unverändert (AGENTS.md §ADR-Disziplin).
- **Schema-Versionierung** (`schemaVersion: 1` im Envelope):
  Cluster §Out-of-Scope hat das als YAGNI für V1 eingestuft.

## Bezug

- Cluster-Slice:
  [`slice-v1-cli-json-dry-run`](slice-v1-cli-json-dry-run.md)
  — §T0-Outcomes (a, c, e) sind die Vorgaben dieses Slices,
  §Aufhebungsbedingung Closure-Hard-Rule ist verbindlich für die
  Cluster-Schließung.
- Spec: `LH-NFA-USE-004` (Minimalkontrakt read-only),
  `LH-FA-CLI-007` (Pflicht-Schema), `LH-FA-DIAG-003` (Severity),
  `LH-FA-CLI-006` (Exit-Codes)
  ([`spec/lastenheft.md`](../../../../spec/lastenheft.md)).
- ADR: [`ADR-0010`](../../adr/0010-kein-http-driving-adapter.md)
  §Folgepunkte Re-Eval-Trigger 2 — dieser Folge-Slice ist Teil
  der JSON-CLI-Spur, die ADR-0010 voraussetzt.
- Code-Anker heute:
  [`cli/doctor.go`](../../../../internal/adapter/driving/cli/doctor.go),
  [`cli/template.go`](../../../../internal/adapter/driving/cli/template.go)
  (`templateJSON`, lokales `--json`-Flag),
  [`cli/root.go`](../../../../internal/adapter/driving/cli/root.go)
  (Vorbild PersistentFlag-Wiring),
  [`domain/diagnostic.go`](../../../../internal/hexagon/domain/diagnostic.go)
  (`Diagnostic.ID/Severity/Message/Hint` als Mapping-Quelle),
  [`cli/cli.go`](../../../../internal/adapter/driving/cli/cli.go)
  (`ErrDoctorFailures` für Exit-Code-Pfad).
- Vorbild-Slice für T0-Outcomes-Layout:
  [`slice-v1-logs`](../done/slice-v1-logs.md) §T0-Outcomes.
- Phase: V1 (Teil des V1-pünktlichen
  `slice-v1-cli-json-dry-run`-Clusters).
