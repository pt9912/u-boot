# Slice Followup: Devcontainer-Features Drift-Doctor (`devcontainer.features.drift`)

> **Status:** ✅ Done. Carveout-Trigger des Parent-Slice
> [`slice-v1-devcontainer-features`](../done/slice-v1-devcontainer-features.md)
> aufgelöst (T1-T4 dort ≈ 1009 LOC > 800-Schwelle). Lifecycle:
> Plan-Anlage `37300c5`, Plan-Followup S1/S2 (Drift-Semantik
> präzisiert) `0c34f0c`, open→in-progress `18ac5a3`, T1+T2
> Implementierung + Tests + User-Doc + CHANGELOG `c2ff32f`, T3
> Closure (slice in-progress→done, carveout aufgelöst)
> `2995524`, Code-Review-Followup S1..S6 (renderKeyOf-Refactor,
> 4 Cross-Check-Tests, formatDriftMessage-Pin, stringSet-
> Konsolidierung) `91f3fb2`, Audit-Followup A3 (Skip-Disziplin
> für `features: {}` präzisiert: explizit-leere-Map + JSON-leer
> → „in sync" statt skip; Test-Assertion-Verschärfung)
> `f69c14b`.

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
  ([`LH-FA-DEV-001`](../../../../spec/lastenheft.md#lh-fa-dev-001--devcontainer-erzeugen) „darf fehlen, ist optional"), unser Check ist
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

## Tranchen

| T   | Inhalt | LOC (Schätzung / Real) |
| --- | ------ | ---------------------- |
| T1  | **Check-ID + Drift-Detector.** ✅ Done in `c2ff32f`. Check-ID `devcontainer.features.drift` + Dispatcher-Slot (Doctor-Total 12→13). `driftJSONFeatureKeys`-Helper (production-Pfad mit `stripJSONC` + `json.Unmarshal`; `ErrDevcontainerJSONUnparsable`-Sentinel). `projectAllFeatureEntries`-Helper baut zwei Mengen (`expectedKeys` enabled-only / `knownProjectableKeys` alle) via Re-use von `projectFeatureEntry` (kein Export nötig — gleicher Package). Drei Set-Differences in `classifyDriftCase1`/`classifyDriftCase2`. `checkDevcontainerFeaturesDrift` aggregiert Warn-Message mit Repair-Hint + Case-2b-Hand-Edit-Hinweis. | ~120 / **~135 real** (+13 %; durch den separaten `classifyDriftCase{1,2}`-Split etwas größer, aber gocognit-konform). |
| T2  | **Tests.** ✅ Done in `c2ff32f`. `doctor_features_drift_test.go` mit 7 Test-Funktionen (3 davon mit Sub-Cases): OKWhenNothingConfigured, OKWhenInSync, Case1_FeatureMissingInJSON, Case1_FileMissing, Case2a_DisabledStillInJSON, Case2b_HandEditUnknownKey, SkipOnParseError, NilVsEmptyFeaturesMap (3 Sub-Cases). | ~80 / **~270 real** (+237 %; Plan-Schätzung hat die separate JSON-Fixture pro Test unterschätzt). |
| T3  | **Slice-Closure.** ✅ Done (dieser Commit). `in-progress/` → `done/`, `carveouts.md`-Zeile auf "aufgelöst durch …", roadmap-Note aktualisiert. User-Doc + CHANGELOG bereits in T1+T2-Commit. | — (Plan-Arbeit) |

**LOC-Bilanz:** ~135 produktion (Schätzung 120, +13 %), ~270 Tests
(Schätzung 80, +237 %). Test-LOC-Inflation kommt vom JSON-Fixture-
Bedarf pro Sub-Case. Plan-Schätzung **~200 total** war realistisch
für Produktion, hat aber die Test-Last unterschätzt. Reuse von
`projectFeatureEntry` (im selben Package, kein Export-Bridge
nötig) hat T1 unter Plan-Schätzung gehalten.

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
- Spec: [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)
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
