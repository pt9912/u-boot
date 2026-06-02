# Slice Followup: Devcontainer-Features Drift-Doctor (`devcontainer.features.drift`)

> **Status:** ready — Trigger gefeuert. T1-T4 des Parent-Slice
> [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md)
> sind real ≈ 1009 LOC (siehe Parent-Plan §Tranchen-Tabelle), über
> der 800-LOC-Carveout-Schwelle. Dieser Folge-Slice nimmt
> Doctor-Teil-B (über Spec hinaus) aus dem Parent-T5 heraus.
> Implementierung startet nach Abschluss des Parent-T5 (Teil A,
> `devcontainer.features.allowlist`) — der allowlist-Check ist die
> strukturelle Grundlage für die Drift-Erkennung.

## Auslöser

Der Parent-Slice §T0-Outcomes hat eine zweigeteilte
Doctor-Integration definiert:

- **Teil A (Spec-mandatiert):**
  `devcontainer.features.allowlist` — `devcontainer.features.<name>.source`
  referenziert nur Quellen aus der Allowlist *oder* der Eintrag ist
  Katalog-aktivierbar. Bleibt im Parent-Slice T5.
- **Teil B (über Spec hinaus, ausgelagert):**
  `devcontainer.features.drift` — der gerenderte `devcontainer.json`
  enthält die aktivierten Features tatsächlich (Managed-Block-
  Disziplin analog M5/M7). Wandert hierher.

Der Spec-Pin
[`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md)
verlangt nur, dass Doctor keine `error`s gegen `devcontainer`-
Konfiguration oder Feature-Quellen ohne legitimen Anlass wirft
— die Drift-Erkennung selbst geht darüber hinaus und ist u-boot-
eigene Ergonomie-Disziplin (Plan-Erweiterung, analog M5 service-
block-drift und M7 managed-block-drift).

## Aufhebungsbedingung

`u-boot doctor` erkennt die folgenden Drift-Situationen und meldet
sie mit Check-ID `devcontainer.features.drift`. Die **Vergleichs-
Granularität** ist verbindlich:

- **cfg-Seite:** für jeden `cfg.Devcontainer.Features.<name>`-
  Eintrag wird der projizierte Render-Key
  `<source>:<version>` über `application.projectFeatureEntry`
  (T3-Helper) gebildet — die *gleiche* Projection, die der
  Generator verwendet. Damit sind cfg- und JSON-Keys
  byte-vergleichbar.
- **JSON-Seite:** die Schlüssel des `features:`-Objekts in
  `.devcontainer/devcontainer.json` (nach `stripJSONC` + JSON-
  Parse).

Drei Drift-Cases (alle Severity `warn`, weil
`u-boot generate devcontainer` als Repair-Hint genügt):

1. **Aktiviertes Feature fehlt im JSON.** Eingangs-Menge: nur
   Einträge mit `Enabled = &true`. Wenn `projectFeatureEntry`
   einen Render-Key liefert, der nicht in den JSON-Keys vorkommt
   → Drift. Sonderfall: wenn `.devcontainer/devcontainer.json`
   ganz fehlt, gilt **jeder** aktivierte Feature-Eintrag als
   Case 1 (siehe §AK „Datei-fehlt-Disziplin" unten).
2. **JSON-Key ohne projizierbares cfg-Pendant.** Eingangs-Menge:
   ALLE `cfg.Devcontainer.Features.<name>`-Einträge — auch
   `Enabled = &false` und `Enabled = nil` (Plan-Hinzufügung zu
   T0-(b)-Pointer-Semantik). Wenn ein JSON-Key zu keinem dieser
   projizierten Einträge passt → Drift mit Sub-Klassifizierung:
   - **Case 2a (Feature deaktiviert/unset):** Render-Key fände
     sich in der Vollprojection — d. h. der User hat
     `enabled: false` gesetzt, aber `generate devcontainer`
     nicht erneut gerufen. Hinweis: „Feature `<name>` wurde
     deaktiviert; entferne den JSON-Eintrag via
     `u-boot generate devcontainer`."
   - **Case 2b (JSON-Key komplett unbekannt):** kein cfg-Eintrag
     mit diesem Render-Key existiert (weder enabled, disabled
     noch unset). Hinweis: „JSON enthält den Feature-Eintrag
     `<key>`, der nicht via u-boot.yaml registriert ist —
     entweder Hand-Edit oder Drift aus früherem u-boot-Stand."

Case 2a und 2b unterscheiden sich nur in der Repair-Message, der
Drift wird in beiden gemeldet.

## Akzeptanzkriterien

- ✅ **Check-ID `devcontainer.features.drift`** im
  `doctorCheckID`-Enum (`internal/hexagon/application/doctor.go`)
  ergänzt. Punktnotation analog zu
  `devcontainer.features.allowlist` (Teil A aus Parent-T5).
- ✅ **Drift-Cases (1 / 2a / 2b, siehe Aufhebungsbedingung)
  erkannt** mit klaren `Repair`-Hints (Case 1 und 2a: „run
  `u-boot generate devcontainer`"; Case 2b zusätzlich Hand-Edit-
  Erkennung in der Message). Severity `warn`.
- ✅ **Vollprojection-Vertrag:** Der Check nutzt
  [`application.projectFeatureEntry`](../../../../internal/hexagon/application/devcontainer_features.go)
  als normative Projection, *nicht* `collectDevcontainerFeatures`.
  Begründung: `collectDevcontainerFeatures` filtert auf enabled —
  für Case 2a (Disabled-Drift) müssen wir auch die deaktivierten
  Einträge projizieren können, um den Render-Key gegen die
  JSON-Keys zu matchen. T1 dieser Followup-Slice darf
  `projectFeatureEntry` öffentlich machen (oder einen kleinen
  Wrapper schaffen), wenn nötig.
- ✅ **Skip-Disziplin (präzise):** Der Check skippt nur, wenn
  weder cfg-Features konfiguriert sind noch JSON-Features
  vorliegen. Konkret:
  - `cfg.Devcontainer == nil` ODER (`cfg.Devcontainer.Features
    == nil`) UND (kein `.devcontainer/devcontainer.json` ODER
    JSON enthält keine `features:`-Section) → skip mit
    Begründung „no devcontainer features configured anywhere".
  - `cfg.Devcontainer.Features == map[]{}` (explizit leere Map)
    UND JSON enthält `features:` → **kein Skip**; Case 2b kann
    feuern.
  - `cfg.Devcontainer.Features` enthält Einträge UND JSON fehlt
    → **kein Skip**; Case 1 feuert (siehe Datei-fehlt-Disziplin).
- ✅ **Datei-fehlt-Disziplin:** Wenn
  `.devcontainer/devcontainer.json` fehlt und mindestens ein
  Feature-Eintrag mit `Enabled = &true` existiert, meldet
  dieser Check `warn` mit Case-1-Hinweis. Die Lücke im
  bestehenden `checkDevcontainerJSON`
  ([`internal/hexagon/application/doctor.go:693`](../../../../internal/hexagon/application/doctor.go)
  — derzeit `SeverityOK` bei `!exists`, Test
  `TestDoctor_DevcontainerJSON_OKWhenAbsent` in
  `doctor_test.go:758` pinnt das) bleibt **bewusst unberührt**:
  jener Check ist auf die File-selbst-Existenz fokussiert
  (LH-FA-DEV-001 „darf fehlen, ist optional"), unser Check ist
  auf die *Konsistenz zwischen cfg und gerendertem JSON*
  fokussiert. T0-Decision: keine Severity-Eskalation in
  `checkDevcontainerJSON` in diesem Slice; falls dort ein
  größeres Refactoring sinnvoll wird, lebt das in einem eigenen
  Folge-Slice.
- ✅ **JSON-Parse-Fehler:** Wenn die JSON ungültig ist (User-
  Edit kaputt), gibt dieser Check `skip` mit Begründung „cannot
  classify drift against unparseable devcontainer.json; fix the
  file or run `u-boot generate devcontainer`". Der bestehende
  `checkDevcontainerJSON` ist für die Validity-Severity zuständig.
- ✅ **Tests** (in `doctor_test.go` oder eigene
  `doctor_features_drift_test.go`):
  - Case 1 Happy-Path (feature aktiv, JSON ohne Key → warn).
  - Case 1 mit fehlender JSON-Datei (feature aktiv, JSON fehlt
    → warn).
  - Case 2a (feature deaktiviert, JSON enthält noch → warn mit
    Case-2a-Message).
  - Case 2b (JSON enthält Key, der zu keinem cfg-Eintrag matcht
    → warn mit Case-2b-Message).
  - No-Drift-Negativtest (cfg + JSON konsistent → OK).
  - Skip-Pins: (a) nichts konfiguriert, (b) leere `features: {}`
    in cfg ohne JSON-Pendant → OK statt skip (Case 2b feuert
    nicht ohne JSON-Eintrag, aber Check ist „fühle mich
    zuständig"-Side OK).
  - JSON-Parse-Error → skip.
  - `nil` vs `features: {}` explizit unterscheiden (zwei
    separate Fixtures), damit die Skip-Bedingung nicht versehentlich
    Pointer-Identität vs Length-Null verwechselt.
- ✅ **README + `docs/user/devcontainer-features.md`** erwähnen
  den Check explizit (Case 1 / 2a / 2b je eine Zeile), damit
  User die Warn-Nachricht zuordnen können (passt zur Parent-T7-
  Doku-Closure; gemeinsamer Block).

## Tranchen (Skizze, wird beim Übergang nach `next/` verfeinert)

| T   | Inhalt | LOC (Schätzung) |
| --- | ------ | --------------- |
| T1  | **Check-ID + Drift-Detector.** Neue Check-ID `devcontainer.features.drift`. Reader für `devcontainer.json` (Production-Pfad-Helper, der `stripJSONC` + `json.Unmarshal` kapselt — der bestehende `StripJSONCForTest`-Bridge in `export_test.go` ist Test-only). **Projection-Schritt:** für jeden cfg-Eintrag (alle, nicht nur enabled) `application.projectFeatureEntry` aufrufen → `(renderKey, enabled)`-Tupel; renderKey = `<Source>:<Version>`. Zwei Mengen: `expectedKeys = {renderKey | enabled == true}`, `knownProjectableKeys = {renderKey | jeder cfg-Eintrag, der projizierbar ist}`. **Set-Differences:** Case 1 = `expectedKeys \ jsonKeys`. Case 2a = `(jsonKeys ∩ knownProjectableKeys) \ expectedKeys`. Case 2b = `jsonKeys \ knownProjectableKeys`. | ~120 |
| T2  | **Tests** (Case 1 happy, Case 1 file-missing, Case 2a, Case 2b, no-drift, skip-pins für nil/leer/parse-error, nil-vs-leer-Distinction). | ~80 |
| T3  | **Slice-Closure:** `open/` → `done/`, Carveouts-Tabelle aktualisieren, Parent-Slice §T5-Status-Update. | — (Plan-Arbeit) |

LOC-Schätzung **~200** in Summe — über der ursprünglichen
Parent-Plan-Vorhersage (~150). Grund: die korrekte Drift-Semantik
braucht eine zweite Projection-Menge (Case 2a vs 2b), und das
Datei-fehlt-Handling plus die nil-vs-leer-Distinction in der
Skip-Logik kosten je ein paar Tests extra. Risiko: wenn
`projectFeatureEntry` aus dem application-Paket-internen Sichtfeld
heraus exportiert werden muss (heute lowercase), kommt eine
zusätzliche export_test.go-Bridge dazu — oder ein Wrapper
`ProjectAllFeatureEntries(cfg) []DriftKey` als neuer
Production-Helper. T0-Entscheidung wird beim Übergang nach `next/`
getroffen.

## Out of Scope

- **Drift-Repair-Automatik** (`doctor --fix` für genau diesen
  Check): nice-to-have, eigener Folge-Slice mit eigenem Trigger.
  Bis dahin reicht der `repair`-Hint.
- **Drift-Detection für `featureSources.allow`:** Allowlist-Drift
  (Eintrag wird nie referenziert) ist ein Linting-Wunsch, kein
  Korrektheits-Bug. Out of scope; eigener Folge-Slice falls
  jemand danach fragt.
- **JSON-Schema-Validierung der `features:`-Section gegen den
  Devcontainers-Spec:** zu groß; wenn ein User den JSON-Block
  manuell editiert, fallen Devcontainer-Tools selbst beim
  Build-Versuch um.
- **Severity-Eskalation in `checkDevcontainerJSON` für den
  `!exists`-Branch:** dieser Slice umgeht die Lücke (Datei-fehlt-
  Disziplin oben), aber das *darunterliegende* Verhalten —
  `checkDevcontainerJSON` returnt `SeverityOK` bei `!exists`
  obwohl `devcontainer.enabled == true` Error verlangen würde —
  wird hier bewusst nicht angefasst. Wenn ein konkreter Trigger
  („Doctor schweigt zu fehlender Datei trotz enabled=true")
  feuert, lebt das in einem eigenen Slice (vermutlich M5-T7-
  Followup), nicht hier.

## Bezug

- Parent-Slice:
  [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md)
  §T0-Outcomes (Doctor-Integration Teil B) + Tranchen-Tabelle T5
  (LOC-Carveout-Trigger).
- Spec: `LH-FA-DEV-003`
  ([`spec/lastenheft.md:692`](../../../../spec/lastenheft.md))
  + Doctor-Pin
  [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md)
  (Negativ-Pin: keine `error`s, kein Spec-Wortlaut zu Drift —
  daher ist dieser Slice „über Spec hinaus").
- Vorbild-Slices:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  (Doctor-Service-Block-Drift-Pattern als Vorlage),
  [`slice-m7-generate`](../done/slice-m7-generate.md)
  (Managed-Block-Disziplin in `devcontainer.json`),
  [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md)
  Parent-T5 Teil A (`devcontainer.features.allowlist`) — die
  zwei Checks teilen sich Render-Helper.
- Code-Anker:
  `internal/hexagon/application/devcontainer_features.go`
  (`collectDevcontainerFeatures`, `projectFeatureEntry`) — die
  Drift-Detection vergleicht gegen genau diese Projection.
- Carveouts: dieser Slice ist selbst kein Carveout-Auflösungs-
  Slice; er ist der Carveout-Trigger-Slice für den 800-LOC-
  Trigger im Parent-Plan (siehe Parent-§Tranchen-Tabelle).
- Phase: V1-Folge, geplant nach Parent-Slice-Abschluss.
