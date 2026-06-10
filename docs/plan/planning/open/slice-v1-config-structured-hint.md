# Slice V1: `config set/get` Strukturiertes `data.hint{action, argument}`-Field

> **Status:** `open/`, on hold pending trigger. Cleanup-/Feature-
> Slice zum WriteAllowed-Reverse-Mapping-Carveout aus
> [`slice-v1-cli-json-dry-run-config`](../done/slice-v1-cli-json-dry-run-config.md)
> §Out of Scope. Carveout-Plan-Anker
> ([[feedback_carveouts_need_plans]]); verlinkt aus
> [`docs/plan/planning/in-progress/carveouts.md`](../in-progress/carveouts.md)
> §Temporäre Carveouts (T8-Closure des config-Slice trägt den
> Eintrag nach).

## Auslöser

`config set` bei WriteAllowed-Reject (Pfad ist nicht
writable, z. B. `services.<svc>.enabled`) liefert heute
einen Reject mit eingebettetem Hint-String "u-boot add
<svc>" in der `diagnostic.message`. `config get` bei
`ErrConfigValueNotSet` liefert analog Hint "u-boot init
--devcontainer". Konsumenten müssen die Hint-Strings per
Substring-Match parsen, was brittle ist.

Ein strukturiertes `data.hint{action: "add", argument: "<svc>"}`
oder `data.hint{action: "init", flag: "--devcontainer"}`
wäre maschinenlesbar und Pipeline-tauglich (Konsument
kann direkt den Hint-Command zusammenstellen ohne Parsing).

V1-Trade-off: Hint-im-Message-String ist heute Pattern-Erbe
(remove `mapWarningsToDiagnostics`). Strukturierte Form
wandert in diesen Folge-Slice.

## Trigger

Plan-Stub bleibt `on hold` bis einer der folgenden Trigger feuert:

- **Real-World-Druck nach Hint-Strukturierung**: CI-Use-Case
  mit Auto-Fix-Pipeline beschwert sich über Substring-Match-
  Pflicht.
- **`remove`/`add` Hint-Strukturierung** als Cluster-übergrei-
  fender Folge-Slice: dann config einbinden.

## Lösungs-Skizze (vorläufig)

`ConfigSetResponse.Hint *ConfigHint` (nullable;
`omitempty`) mit `{Action string, Argument string, Flag
string}`-Form. Application-Layer-Setter bei WriteAllowed-
Reject + ValueNotSet. CLI-Layer mappt auf
`envelope.data.hint`. JSON-Konsument prüft `hint != null`
und baut den Folge-Command zusammen.

## Heute-Stand-Pre-Scan + T0-Discovery (2026-06-08)

> **Kein Trigger gefeuert** — reine Konsument-Komfort-Strukturierung,
> proaktiv geplant (User-Wunsch „beide planen, dann entscheiden").
> Bleibt `open/` bis Umsetzungs-Entscheid.

**Code-Realität (grounding):**

- Der Hint ist **Teil der Error-Message**, NICHT ein Response-Feld.
  Gebaut via `fmt.Errorf("%w: … use \`u-boot add %s\` …", sentinel,
  path.Service.String())`:
  - `writeRejectedError(path)` — `application/config.go:357-367`
    (`ErrConfigWriteRejected`, Kind `ConfigServiceEnabled` →
    „`u-boot add <svc>`").
  - `ErrConfigValueNotSet`-Hints — `application/config.go:743-810`,
    sechs Kind-Varianten (`u-boot init --devcontainer`,
    `u-boot add <svc>`, `u-boot config set <path> <value>`, …).
- Am Bau-Punkt verfügbar: `path domain.ConfigPath` mit `.Kind`,
  `.Service`, `.Feature`, `.String()`. Der **„action"-Teil** (`add`/
  `init`/`config set`) ist heute **hartkodiert pro Kind** im Format-
  String; das **„argument"** ist `path.Service.String()` bzw.
  `path.String()`.
- Response-Structs (`port/driving/config.go`):
  `ConfigSetResponse{Path, OldValue, NewValue, Warnings,
  PlannedFiles}`, `ConfigGetResponse{Path, Value}` — **kein
  Hint-Feld**. Der Hint existiert nur auf dem **Error-Pfad** (die
  Funktion gibt einen `error` zurück, keine Response).
- CLI-Daten-Carrier: `configGetData{path,value}`,
  `configSetData{path,oldValue,newValue,noOp,appendedSources}`
  (`cli/config.go:27-45`).

**Plan-kritische Korrektur gegenüber dem Stub:** Der Stub schlug
`ConfigSetResponse.Hint *ConfigHint` vor — **das passt nicht**, weil
der Hint auf dem **Error-Pfad** lebt (kein Response). Der CLI-Mapper
(`mapConfigErrorToDiagnostic`) sieht nur den sentinel-gewrappten
Fehler; um die `path.Service`-Info zu bekommen, müsste er die
Message **parsen** — genau die Brittleness, die wir abbauen. Richtig
ist ein **typed-error-Carrier**, der die strukturierten Teile trägt.

## Sub-Decisions (T0, zum Review)

- **SD-A1 — Typed-Error-Carrier (statt Response-Feld).** Neuer
  `driving.ConfigHintError struct { inner error; Action, Argument,
  Flag string }` mit `Error() string` (= heutige Message, unverändert
  fürs Human-Auge) + `Unwrap() error` (→ Sentinel, damit
  `errors.Is(err, ErrConfigWriteRejected)`/`ErrConfigValueNotSet` +
  `ExitCode` intakt bleiben — Pattern-Erbe `baseDirSanitizedError`).
  `writeRejectedError`/die ValueNotSet-Pfade wrappen ihr Ergebnis
  damit. Der CLI-Error-Pfad (`reportErrorSub`) extrahiert via
  `errors.As(err, *ConfigHintError)` die Felder ins
  `data.hint`-Envelope-Feld.
- **SD-A2 — Hint-Schema.** Die Hints haben drei Formen:
  `u-boot add <svc>` (action+argument), `u-boot init --devcontainer`
  (action+flag), `u-boot config set <path> <value>` (action+
  argument). Vorschlag: `data.hint{command: "<ready-to-run string>",
  action: "add|init|config-set", argument?: "<svc|path>", flag?:
  "<--flag>"}` — der `command`-String ist die kanonische, direkt
  ausführbare Form (Konsument braucht kein Re-Assembly), die
  strukturierten Felder erlauben Programmatik. **Review-Frage:**
  reicht `{action, argument, flag}` ohne den `command`-String, oder
  ist der ready-to-run-String der Hauptwert?
- **SD-A3 — Scope.** Beide Hint-Quellen abdecken
  (`ErrConfigWriteRejected` set-Pfad + `ErrConfigValueNotSet`
  get-Pfad), da beide brittle Hint-Strings tragen. `data.hint` ist
  ein Error-Envelope-Feld auf beiden Formen.

## Tranchen (vorläufig)

| Tranche | Inhalt | LOC |
| --- | --- | --- |
| T0 | Discovery (dieser Pre-Scan) + Sub-Decisions | — |
| T1 | `ConfigHintError`-Carrier (port) + `Error`/`Unwrap`; `writeRejectedError` + ValueNotSet-Pfade wrappen Action/Argument/Flag (aus `.Kind`/`.Service`/`.String()`); `errors.Is`/`ExitCode`-Pins bleiben grün | ~80 |
| T2 | CLI: `errors.As` im Error-Pfad → `data.hint` (neuer `configHintData`-Carrier); Envelope-Schema-Pin | ~60 |
| T3 | Acceptance-Tests: set-WriteRejected + get-ValueNotSet → `data.hint{...}` strukturiert; Human-Message unverändert; `errors.Is`-Chain intakt | ~80 |
| T4 | Closure: carveouts-Eintrag entfernen, CHANGELOG `### Added`, `cli-json-output.md` §6.9-Ergänzung, `done/`-Move + DoD | — |

LOC **~220** (mittel). Berührt Port-Contract (neuer Error-Typ),
aber **kein neuer Sentinel** (wrappt bestehende), kein neues
Subcommand.

## Spec-Bezug

- [`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern) (Path-Whitelist) — Spec macht keine
  Aussage zur Hint-Form. Strukturierung ist Konsument-
  Komfort-Argument (V1+-Ergonomik, kein Spec-Pflicht-Surface).
