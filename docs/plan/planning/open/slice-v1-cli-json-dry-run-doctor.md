# Slice V1: `doctor --json` — Pattern-Vorbild für read-only-Envelope

> **Status:** geplant für v0.4.0 — erster Folge-Slice des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../next/slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 1). Liefert die gemeinsame Infrastruktur für die
> 9er-Folge-Slice-Serie: Root-PersistentFlag `--json` (T0-(a)),
> Common-Envelope `cliJSONEnvelope` (T0-(c)) und Schema-Helper
> `jsontestutil.AssertSchemaConform`. Trägt zusätzlich den
> Übergangs-Schnitt für das existierende `template list --json`
> (Review-Round-2-Finding M3).

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
wie heute. Alle 8 noch nicht migrierten Subcommands (`init`,
`add`, `remove`, `up`, `down`, `logs`, `generate`, `config` in
seinen drei Formen) rejecten `--json` mit Exit-Code `2`
(`LH-FA-CLI-006`-Klasse) und Verweis auf den jeweiligen
Folge-Slice — kein Subcommand „akzeptiert `--json` still und
liefert Human-Output" (sonst untergrabener V1-Maschinenvertrag).

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
#     {"level":"warn","code":"LH-FA-DIAG-002","message":"…",
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
#     {"level":"error","code":"LH-FA-DIAG-002","message":"…"}
#   ],
#   "exitCode": 11
# }
```

`status` ist an das höchste `diagnostics[].level` gekoppelt
(`error → "error"`, `warn ohne error → "warn"`, sonst `"ok"`,
Lastenheft §447), analog zu der heutigen `HasErrors()`/
`HasWarnings()`-Logik in
[`doctor.go`](../../../../internal/adapter/driving/cli/doctor.go).
`SeverityOK`-Items des heutigen Plaintext-Reports werden im
JSON-Modus **nicht** als `diagnostics[]`-Eintrag serialisiert
(Spec verbietet `level: "ok"`; `--quiet` ist im JSON-Modus
darum semantisch ein No-op).

## Akzeptanzkriterien

- ✅ **Root-PersistentFlag `--json` (Cluster T0-(a))**: am
  Cobra-Root registriert, persistent für alle 10 Subcommands,
  Wiring analog
  [`--verbose`/`--debug`/`--quiet`](../../../../internal/adapter/driving/cli/root.go).
- ✅ **Reject-Pfad für nicht-migrierte Subcommands**: alle
  9 noch nicht migrierten Spec-Enum-Subcommands (`init`, `add`,
  `remove`, `up`, `down`, `logs`, `generate`, `config` ×3,
  `template list` siehe Sonderregel unten) rejecten `--json`
  mit Exit-Code `2` (`LH-FA-CLI-006`-Klasse) und Fehlermeldung
  `JSON-Ausgabe für 'u-boot <sub>' ist noch nicht implementiert
  (siehe slice-v1-cli-json-dry-run-<sub>).` Mechanik T0-(h).
  Pflicht-Pin-Test je nicht-migrierter Subcommand. Pro
  Folge-Slice-Merge fällt **genau** der Reject-Pfad des
  migrierten Subcommands; **Cluster-T_close-Pflicht-Check:**
  null offene Reject-Pfade.
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
    `diagnostics[].level` ∈ `{warn, error}` (kein `ok`),
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
  gemappt, **wobei `SeverityOK`-Items übersprungen werden**
  (Spec §1834 verbietet `level: "ok"`). Sort-Order analog
  [`SortedByIssuesFirst`](../../../../internal/hexagon/domain/diagnostic.go),
  aber gefiltert auf Warn/Error. `--quiet` und `--strict`
  interagieren mit `--json` definiert (Sub-Decision T0-(e)
  dieses Slices).
- ✅ **Code-Registry-Doku (Spec §1835 + §445)**: in
  `docs/user/cli-json-output.md` eine Sektion „Code-Registry"
  pflichtbestandteil. Wenn T0-(g) `Diagnostic.ID`-basierte
  Tool-interne Codes wählt (z. B. `docker.available`), trägt
  diese Sektion jeden Code mit seiner Bedeutung; wenn T0-(g)
  LH-ID-Mapping wählt (z. B. alle Doctor-Checks → `LH-FA-DIAG-002`
  oder `LH-FA-DIAG-003`), dokumentiert die Sektion das Mapping.
  Der Minimal-Helper konsultiert diese Sektion über eine
  embedded Allowlist (Sub-Decision T0-(b)/(g)).
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
  [`carveouts.md`](../in-progress/carveouts.md) §Temporäre
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
hängt von T0-(g) ab.

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

### T0-(h) Reject-Mechanik für nicht-migrierte Subcommands

Wenn das Root-PersistentFlag `--json` ab T3 für alle 10
Subcommands existiert, müssen die 9 noch nicht migrierten
Subcommands ihn explizit rejecten (vgl. Aufhebungsbedingung +
AK Reject-Pfad). Drei Mechanik-Optionen:

1. **PreRunE pro Subcommand**: jeder noch nicht migrierte
   Subcommand bekommt einen `PreRunE`-Hook, der den Root-Flag
   abfragt und bei `--json` mit Exit-Code `2` abbricht.
   Vorteil: lokal, leicht zu entfernen pro Folge-Slice-Merge.
   Nachteil: 9 separate Wirings, einfach zu vergessen.
2. **Zentrale Subcommand-Allowlist** in `cli/root.go`: eine
   Map `migratedJSONSubcommands` listet die migrierten
   Subcommands. Cobra-`PersistentPreRunE` am Root prüft, ob
   `--json` gesetzt und Subcommand nicht in Allowlist → Reject.
   Pro Folge-Slice-Merge wird **ein Eintrag** in die Allowlist
   gesetzt. Vorteil: zentraler Drift-Anker; Nachteil: zentrales
   File ändert sich oft. Cluster-T_close-Check trivial (Allow-
   list muss alle 10 enthalten — bzw. Allowlist und Reject-
   Mechanik können dann ganz entfernt werden).
3. **Reject im Envelope-Builder**: wenn ein Subcommand keinen
   Envelope-Builder verlinkt hat, fällt der Build durch ein
   `nil`-Check und gibt Exit-Code 2. Vorteil: implicit;
   Nachteil: schwer zu testen, leise Fehler-Klasse.

Vorschlag: Option (2) — zentrale Allowlist in `cli/root.go`
mit `PersistentPreRunE`. Die Allowlist wandert mit jedem
Folge-Slice ein Stück nach vorne; Cluster-T_close entfernt
Allowlist + Reject-Pfad komplett.

### T0-(g) `diagnostics[].code`-Quelle vs. Spec-Code-Konvention

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

## Tranchen (vorgeschlagen)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + Sub-Decisions** aus §T0-Discovery klären (acht Sub-Decisions, inkl. T0-(h) Reject-Mechanik). Entscheidungen mit Begründung in einem `T0-Outcomes`-Block dokumentieren. | — (Plan-Arbeit) |
| T1 | **Schema-Vertrag-Doku.** `docs/user/cli-json-output.md` anlegen: Minimalkontrakt (`LH-NFA-USE-004` §1823-1842) und Voll-Schema (`LH-FA-CLI-007` §322-417) verbatim getrennt zitiert; **Code-Registry-Sektion** für Doctor-Checks (gemäß T0-(g)); Envelope-Lokation benennen; Minimal-vs.-Voll-Diff klargestellt. README EN+DE bekommt einen Verweis-Eintrag. (Cluster-T1; per T0-(f) hier mitgeliefert.) | ~120 |
| T2 | **Envelope + Helper-Split.** `cliJSONEnvelope`-Typ (Minimal-Felder Pflicht, Voll-Felder `omitempty` — gemäß T0-(d)) plus zwei Helper `AssertMinimalEnvelope` und `AssertFullEnvelope`. Unit-Tests beider Helper (positive + negative Cases — fehlende Pflichtfelder, ungültiges `status`, `level: "ok"`-Reject, Voll-Schema-Feld im Minimal-Pfad-Reject, undokumentierter Code-Reject). | ~180 |
| T3 | **Root-PersistentFlag `--json` + Reject-Allowlist** am Cobra-Root registrieren plus `PersistentPreRunE`-Reject für nicht-migrierte Subcommands (Mechanik gemäß T0-(h)). 9 Reject-Pin-Tests, je einer pro nicht-migrierter Subcommand-Form (`init`, `add`, `remove`, `up`, `down`, `logs`, `generate`, `config`, `config get`, `config set`; `template list` siehe T4). App-Struktur-Field für den Flag-State. | ~80 |
| T4 | **`template list`-Schnitt + Carveouts-Eintrag.** Lokales Flag entfernen, `runTemplateList` liest Root-Flag. Zwei Pin-Tests grün (`u-boot template list --json` + `u-boot --json template list` → gleicher Output). Carveouts-Eintrag in `carveouts.md` §Temporäre Carveouts: `template list --json` heute nicht Minimalkontrakt-konform, Re-Trigger `slice-v1-cli-json-dry-run-template`. | ~40 |
| T5 | **`doctor --json`-Pfad.** Envelope-Befüllung aus `domain.DiagnosticReport` mit `SeverityOK`-Filter (kein `level: "ok"`-Eintrag), `status`-Mapping (höchstes vorhandenes Severity-Level), `exitCode` konsistent mit `ErrDoctorFailures`. Drei Acceptance-Tests (All-OK mit `diagnostics: []`, Warn-Fall, Error-Fall) via `jsontestutil.AssertMinimalEnvelope`. `--quiet`/`--strict`-Interaktion gemäß T0-(e). Doctor-Check-Codes in Code-Registry aus T1 ergänzt. | ~140 |
| T6 | **Closure.** CHANGELOG `## [Unreleased]` Added-Eintrag (Envelope + Helper-Split + Root-Flag + Reject-Allowlist + `doctor --json` + `template list`-Flag-Schnitt + Carveout). roadmap.md v0.4.0-Tabelle: Cluster-Slice-Zelle aktualisiert (Doctor done, nächster Schritt `add`). Slice-File `open/` → `done/` mit Tranchen+Commit-Tabelle. `make docs-check` grün. | — (Doku) |

LOC-Schätzung Folge-Slice: ~560 LOC — am oberen Ende der vom
Cluster-Slice gesetzten 200..600-Bandbreite, weil dieser Slice
neben dem Doctor-Pfad gleichzeitig den Minimal-Helper, den
Voll-Helper-Stub, die Reject-Allowlist und die Code-Registry-
Doku etabliert (Pattern-Vorbild-Last).

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
  Preview-Adapter.
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
  [`slice-v1-cli-json-dry-run`](../next/slice-v1-cli-json-dry-run.md)
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
