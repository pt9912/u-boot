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

## Spec-Bezug

- `LH-FA-CONF-005` (Path-Whitelist) — Spec macht keine
  Aussage zur Hint-Form. Strukturierung ist Konsument-
  Komfort-Argument.
