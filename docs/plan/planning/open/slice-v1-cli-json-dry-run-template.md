# Slice V1: `template list --json` — Envelope-Migration

> **Status:** geplant für v0.4.0+ — letzter Folge-Slice des
> Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Platz 9). Closure-Pflicht-Slice für den
> Cluster-T_close-Lauf, weil er den bewussten Übergangs-Carveout
> aus dem Doctor-Slice schließt.

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

`u-boot --json template` (bare) klärt die Spec-§1838-
`subcommand`-Pflicht für den Help-Parent — Sub-Decision in
diesem Slice: Reject mit Exit-Code 2 (status quo) oder
Default-`subcommand: "list"`? Vorschlag: Reject (Help-Parent
trägt kein eigenes Datum, Default-Subcommand würde Doppeldeutig-
keit zwischen Help und List schaffen).

## Akzeptanzkriterien

- ✅ **Envelope-Migration**: `cliJSONEnvelope` mit
  `command: "template"`, `subcommand: "list"`, `Data` als
  `[]templateJSON`. **Neuer Konstruktor `newDataEnvelope`**
  (oder analoger Helper) etabliert die "Data-im-Minimal-Mode"-
  Disziplin mit eigenem Marshal-Pin-Test (Review-Finding H1 aus
  `slice-v1-cli-json-dry-run-doctor` adressiert).
- ✅ **`Data`-Feld im `cliJSONEnvelope`**: wird in diesem Slice
  wieder eingeführt (war in T2 entfernt, weil kein Konsument
  existierte). Type bleibt `any` mit `omitempty`-Tag. Sub-
  Decision: Pattern-Disziplin per `newDataEnvelope(command,
  subcommand, data any, ...)`-Konstruktor (Plan-Vorgabe aus
  Doctor-Slice T0-(c)).
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
| T0 | **Discovery + Sub-Decisions** für bare-`template`-Verhalten (Reject vs. Default-Subcommand), `Data`-Konstruktor-Form, Code-Registry-Bedarf. | — (Plan-Arbeit) |
| T1 | **`cliJSONEnvelope` erweitern**: `Data any \`json:"data,omitempty"\`` wieder ergänzen; `newDataEnvelope(command, subcommand string, data any, diags, exitCode)`-Konstruktor + Marshal-Pin-Test. | ~80 |
| T2 | **`runTemplateList`-Migration**: Array-Output durch Envelope ersetzen, `templateJSON`-Slice als `Data`. Bestehender Pin-Test `TestRootJSON_AcceptsTemplateList_BothFlagPositions` muss überarbeitet werden — Format-Wechsel ist Breaking-Change im JSON-Surface (rechtfertigt v0.5.0-Bump oder Carveouts.md-permanent-Eintrag falls Konsument-Verträglichkeit erforderlich). | ~80 |
| T3 | **bare `template` Sub-Decision umsetzen**: Default-Subcommand oder expliziter Reject. | ~30 |
| T4 | **Cluster-T_close-Vorbereitung**: Carveouts-Eintrag entfernen, Allowlist-Mechanik komplett abbauen (siehe Cluster-Slice §T0-(g) Cluster-T_close-Pflicht-Check). | ~40 |
| T5 | **Closure.** CHANGELOG, Roadmap, Cluster-Slice nach `done/` (zusammen mit diesem Slice), DoD-Hash-Tabelle. | — (Doku) |

LOC-Schätzung: ~230 LOC, im niedrigen Bereich der Cluster-
LOC-Bandbreite.

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
