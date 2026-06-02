# Slice Followup: Devcontainer-Features Drift-Doctor (`devcontainer.features.drift`)

> **Status:** ready — Trigger gefeuert. T1-T4 des Parent-Slice
> [`slice-v1-devcontainer-features`](../in-progress/slice-v1-devcontainer-features.md)
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

`u-boot doctor` erkennt jede der drei Drift-Situationen und meldet
sie mit Check-ID `devcontainer.features.drift`:

1. **Feature aktiv, aber im JSON fehlend:** `cfg.Devcontainer.Features[<name>].Enabled = &true`,
   aber `devcontainer.json.features` enthält den per
   `<source>:<version>` zusammengesetzten Key nicht.
2. **Feature deaktiviert, aber im JSON noch present:** der Key
   steht im JSON, obwohl `cfg.Devcontainer.Features[<name>].Enabled`
   `nil` oder `&false` ist.
3. **Feature im JSON aber nicht in `cfg`:** ein Key in
   `devcontainer.json.features` hat keinen entsprechenden
   `cfg.Devcontainer.Features.<name>`-Eintrag.

Alle drei Situationen klassifizieren als Severity `warn` (nicht
`error`) — der User kann den Drift selbst beheben mit
`u-boot generate devcontainer`, und Doctor liefert genau diesen
Repair-Hint.

## Akzeptanzkriterien

- ✅ **Check-ID `devcontainer.features.drift`** im
  `doctorCheckID`-Enum (`internal/hexagon/application/doctor.go`)
  ergänzt. Punktnotation analog zu
  `devcontainer.features.allowlist` (Teil A aus Parent-T5).
- ✅ **Drei Drift-Cases (siehe Aufhebungsbedingung) erkannt** mit
  klaren `Repair`-Hints („run `u-boot generate devcontainer`").
  Severity `warn`.
- ✅ **Kein Trigger ohne Features:** Projekte ohne
  `cfg.Devcontainer.Features` (Pre-LH-FA-DEV-003-Stand) bekommen
  einen `skip`-Result mit Begründung — der Check fühlt sich nicht
  zuständig.
- ✅ **Kein Trigger ohne `.devcontainer/devcontainer.json`:** wenn
  das File fehlt, ist M7-`devcontainer.json.valid` (oder die
  M5-Severity-Eskalation gegen `devcontainer.enabled = true`)
  zuständig, nicht dieser Check.
- ✅ **Tests** (in `doctor_test.go` oder eigene
  `doctor_features_drift_test.go`): pro Drift-Case je ein
  Happy-Path + ein No-Drift-Negativtest. Plus „kein Features
  konfiguriert → skip"-Pin und „kein devcontainer.json → skip".
- ✅ **README + `docs/user/devcontainer-features.md`** erwähnen
  den Check explizit, damit User die Warn-Nachricht zuordnen
  können (passt zur Parent-T7-Doku-Closure; die Doku kann in
  einem gemeinsamen Block landen).

## Tranchen (Skizze, wird beim Übergang nach `next/` verfeinert)

| T   | Inhalt | LOC (Schätzung) |
| --- | ------ | --------------- |
| T1  | **Check-ID + Drift-Detector.** Neue Check-ID, Reader für `devcontainer.json` (kann den bestehenden M7-`StripJSONCForTest`-Helper nutzen, sofern der Production-Pfad einen analog-positionierten Helper hat — sonst Mini-Erweiterung), Set-Difference zwischen `cfg.Devcontainer.Features` (enabled) und `json.features`-Keys. | ~100 |
| T2  | **Tests** (3 Drift-Cases + 2 Skip-Cases) + Doku-Eintrag. | ~50 |
| T3  | **Slice-Closure:** `open/` → `done/`, Carveouts-Tabelle aktualisieren, Parent-Slice §T5-Status-Update. | — (Plan-Arbeit) |

LOC-Schätzung **~150** in Summe — exakt der Wert aus der Parent-
Plan-Vorhersage. Risiko: wenn die Set-Difference-Logik im
`projectFeatureEntry`-Aufruf-Pfad re-used werden kann (siehe
`internal/hexagon/application/devcontainer_features.go:projectFeatureEntry`),
fällt T1 kleiner aus. Re-Check bei T1-Abschluss.

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

## Bezug

- Parent-Slice:
  [`slice-v1-devcontainer-features`](../in-progress/slice-v1-devcontainer-features.md)
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
  [`slice-v1-devcontainer-features`](../in-progress/slice-v1-devcontainer-features.md)
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
