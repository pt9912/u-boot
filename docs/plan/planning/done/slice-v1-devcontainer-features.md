# Slice V1: Devcontainer-Features / Toolchains ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features))

> **Status:** ✅ Done (v0.4.0-Material). T0 Discovery `7a8b8ad`,
> T1 Schema + `domain.FeatureName` `a97337a`,
> T2 Catalogue (8 Spec-Beispiele) `6de1464`,
> T3 Generator-Patch `e1646f1`,
> T4 CLI (`--allow-external-feature-sources` + 4 ConfigPath-Kinds)
> `a8d10bd` + Review-Followup R1..R6 `420e19f`,
> T5 Doctor Teil A (`devcontainer.features.allowlist`) `35a5bae`,
> T6 LHFADEV003 Acceptance-Pins `e9f7282`,
> T7 Doku-Closure `268141f`. Audit-Followup A1..A4
> (Schema-Validation-Wiring in 3 Production-Load-Pfade +
> Idempotenz-Fix bei Doppel-Source-Versionen + README-Update)
> `f69c14b` + README-Detail-Entlastung `d47ffa1`.
> Doctor Teil B (`devcontainer.features.drift`) ausgelagert in
> [`slice-followup-devcontainer-features-drift-doctor`](slice-followup-devcontainer-features-drift-doctor.md)
> (`37300c5` + S1/S2-Plan-Followup `0c34f0c` + T1+T2 `c2ff32f` +
> T3 Closure `2995524` + S1..S6 Code-Review-Followup `91f3fb2`)
> nach 800-LOC-Carveout-Trigger.

## Auslöser

Der heutige Add-on-Katalog kennt nur **Laufzeit-Services**
(Postgres, Keycloak, OpenTelemetry) — alles Compose-Services mit
Image-Pin, Healthcheck und Port. Eine reproduzierbare
Entwicklungsumgebung braucht aber auch **Toolchains** (Java,
Gradle, Go, Node.js, C++, …). Ohne sie ist `u-boot init` zwar
ein lauffähiger Compose-Stack, aber kein einsatzbereites
Dev-Setup — der User muss die Sprach-Toolchain selbst im
Devcontainer-`Dockerfile` oder per `features:`-Section
nachziehen.

Die Spec hat das vorgesehen, aber die Anforderung ist bisher
nicht implementiert:

- **[`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)** (Priorität V1) listet explizit Git, Docker
  CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools
  als Beispiele und legt das Security-Modell fest
  (`--allow-external-feature-sources`,
  `devcontainer.featureSources.allow`, `--yes` reicht nicht).
- **[ADR-0008](../../adr/0008-plugin-system-statisch.md)** §78 nennt die Allowlist ausdrücklich als einzigen
  ausnahme-pflichtigen externer-Code-Pfad ausserhalb des
  Compose-Stacks.
- **Code-Anker:** `internal/hexagon/application/initproject.go:66-69`
  markiert die Future-Field-Stelle im YAML-Schema
  (`Future fields:`-Kommentarblock mit
  `FeatureSources ubootYAMLFeatureSources`).
- **Schema-Skizze:** `spec/lastenheft.md:1340-1343` zeigt die
  YAML-Form (`devcontainer.featureSources.allow: [URL, …]`).

## Aufhebungsbedingung

`u-boot init --devcontainer` plus ein neuer Feature-Flow erzeugt
einen `.devcontainer/devcontainer.json`, der eine vom User
gewählte Toolchain als `features:`-Section enthält; die Allowlist
in `u-boot.yaml` ist gepflegt; `u-boot doctor` läuft ohne Errors;
`u-boot generate devcontainer` ist idempotent (managed-block-
Stabilität analog zu M7).

## Akzeptanzkriterien

- ✅ **Schema (`featureSources.allow`, Spec-mandatiert):**
  `ubootYAMLDevcontainer.FeatureSources` ist implementiert (heute
  als Future-Field-Kommentar in `initproject.go:66-69`).
  `devcontainer.featureSources.allow` ist Liste von Quell-Strings,
  dedupliziert (`spec/lastenheft.md:1352`). Konkrete
  Failure-Cases der Validation: leerer Quell-String, fehlendes
  URL-Scheme, fehlende Host-Komponente — siehe T1 für die
  vollständige Failure-Tabelle. Diagnose-Code: [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features),
  Exit-Code `10` (`spec/lastenheft.md:1353`).
- ✅ **Feature-Aktivierung in u-boot.yaml (Plan-eigene Schema-
  Erweiterung, *nicht* Spec-mandatiert):** Spec definiert nur
  `featureSources.allow` (Quellen-Whitelist, §719). Für die
  deterministische Regeneration des `devcontainer.json`-Blocks
  braucht der Plan zusätzlich eine *Aktivierungs-Map* analog zur
  `services.<name>.enabled`-Konvention aus M5
  ([`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)).
  Diese Erweiterung geht bewusst über den Spec-Wortlaut hinaus
  und ist als Plan-Hinzufügung explizit markiert.
  **T0-Entscheidung (§T0-Outcomes, (b)):** Map-analog
  `devcontainer.features.<name>.enabled` mit Pointer-`*bool`-
  Semantik wie `ubootYAMLService.Enabled` aus
  `initproject.go:54-56`, plus optionalem `Source string` und
  optionalem `Version string` pro Eintrag.
- ✅ **CLI:** **T0-Entscheidung (§T0-Outcomes, (a)):** Aktivierung
  ausschließlich über
  `u-boot config set devcontainer.features.<name>.enabled true`
  (M8-Path-Whitelist erweitern) + nachgelagertes
  `u-boot generate devcontainer`. Kein neues Top-Level-
  Kommando. Die drei Spec-Pfade aus `spec/lastenheft.md:714-717`
  für `--allow-external-feature-sources` (`init --devcontainer`,
  `generate devcontainer`,
  `config set devcontainer.featureSources.allow`) bleiben
  unverändert; das Flag wird strikt auf diese drei Pfade
  beschränkt und auf dem Aktivierungs-Pfad **nicht** akzeptiert
  (Default-Variante, Spec-treu).
- ✅ **Allowlist-Enforcement:** Ohne explizit erlaubte Quelle
  bricht der Versuch, eine externe Feature-Quelle zu aktivieren,
  mit code [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) / Exit-Code `10` ab. `--yes` reicht
  nicht ([`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)). `--allow-external-feature-sources`-
  Flag-Argumente werden in `u-boot.yaml`
  (`devcontainer.featureSources.allow`) persistiert — Flag ist
  Erweiterung der Liste, nicht Ein-Aufruf-Override
  (`spec/lastenheft.md:719`: „als explizit freigegebene Liste in
  der Projektkonfiguration gespeichert").
- ✅ **Devcontainer-Template:** **T0-Entscheidung
  (§T0-Outcomes, (d)):** Der bestehende `init`-Managed-Block in
  `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  wird um eine `"features": { … }`-Eigenschaft erweitert
  (single-block, keine Marker-Nesting). JSON-Stabilität
  (alphabetische Sortierung nach `Source`-URL, deterministische
  Ausgabe) ist gewahrt.
- ✅ **Doctor-Integration Teil A (Spec-mandatiert):** `u-boot
  doctor` validiert Allowlist-Konformität: `devcontainer.features`
  referenziert nur Quellen aus `devcontainer.featureSources.allow`
  *oder* aus dem eingebauten Katalog (Spec-Pin
  [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md):
  „`u-boot doctor` enthält keinen `error` zu `devcontainer`-
  Konfiguration oder Feature-Quellen"). Check-ID:
  `devcontainer.features.allowlist` (Punktnotation analog zum
  bestehenden `doctorCheckID`-Enum in `doctor.go:69 ff.`).
- 🟡 **Doctor-Integration Teil B (über Spec hinaus,
  ausgelagert):** Drift-Erkennung —
  `devcontainer.json` enthält die aktivierten Features
  tatsächlich (Managed-Block-Disziplin analog M5/M7). Check-ID:
  `devcontainer.features.drift`. **Carveout-Trigger gefeuert
  nach T4:** T1-T4-Real-LOC ≈ 1009 > 800-Schwelle (siehe LOC-
  Bilanz unter den Tranchen). Folge-Slice
  [`slice-followup-devcontainer-features-drift-doctor`](slice-followup-devcontainer-features-drift-doctor.md)
  in `open/` angelegt; Teil B implementiert dort, T5 in diesem
  Slice realisiert nur noch Teil A.
- ✅ **Statischer Katalog** ([ADR-0008](../../adr/0008-plugin-system-statisch.md)-konform): u-boot bringt
  einen kuratierten Default-Satz Features mit (mindestens die
  Spec-Beispiele: Git, Docker CLI, Node, Java/SDKMAN, Go, C++,
  K8s-Tools, PostgreSQL-Client), die ohne Allowlist-Eintrag
  aktivierbar sind. Externe Quellen jenseits davon brauchen die
  Allowlist. **Lesart des Spec-Begriffs „lokal hinterlegt"**
  (`spec/lastenheft.md:711`): in u-boot eingebauter statischer
  Go-Mapping-Katalog analog `serviceCatalogue()` in
  `addservice_execute.go:234 ff.`. Die Alternativ-Lesart
  „Features im Repo-Pfad" ist bewusst ausgeschlossen und Trigger
  für ein eigenes Folge-Slice (siehe Out of Scope).
  **T0-Entscheidung (§T0-Outcomes, (c)):** Katalog-Key ist die
  neue Domain-Type `domain.FeatureName` analog
  `domain.ServiceName` (name-only, Regex
  `^[a-z]([a-z0-9-]{0,30}[a-z0-9])?$`, 32-Char-Cap). Kein
  Version-Slot im Domain-Type — Version lebt in
  `featureCatalogueEntry.defaultVersion` plus optionalem
  `ubootYAMLDevcontainerFeature.Version`-Override.
- ✅ **Spec-Pin:** `internal/hexagon/application/acceptance_test.go`
  deckt [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) ab. Test-Naming-Konvention analog zum
  bestehenden `TestLHAK###_<Title>`-Muster (siehe
  `acceptance_test.go:39, 98`):
  `TestLHFADEV003_AllowlistEnforcement` (negativer Pfad — externe
  Quelle ohne Allowlist → Exit-Code `10`) +
  `TestLHFADEV003_CatalogueActivation` (positiver Pfad — Katalog-
  Feature aktivieren + `devcontainer.json` enthält Eintrag).
  Docker-e2e in `internal/e2e/` ist **nicht** erforderlich:
  Devcontainer-Features werden nicht zur Laufzeit gegen einen
  echten Docker-Stack geprüft, sondern strukturell im
  generierten JSON.
- ✅ **Doku:** README (EN + DE) erweitert um Feature-Beispiel;
  `docs/user/devcontainer-features.md` (oder vergleichbar) listet
  Katalog + Allowlist-Mechanik.

## Tranchen

| T   | Inhalt | LOC (Schätzung) |
| --- | ------ | --------------- |
| T0  | **Discovery / Design.** ✅ Vier Design-Entscheidungen + LOC-Tabelle festgezurrt — siehe [§T0-Outcomes](#t0-outcomes). | — (Plan-Arbeit) |
| T1  | **Schema-Erweiterung.** ✅ Done. `domain.FeatureName` (`featurename.go`, 56 LOC) + `ubootYAMLFeatureSources`/`ubootYAMLDevcontainerFeature` + Erweiterung von `ubootYAMLDevcontainer` (initproject.go-Diff ~40 LOC) + `validateFeatureSource`/`dedupeFeatureSources`/`normaliseFeatureSources`/`validateDevcontainerFeatures` (`devcontainer_features.go`, 174 LOC) + Tests (featurename_test.go 67 LOC + devcontainer_features_test.go 301 LOC). **Failure-Tabelle vollständig gepinnt** (allowlist: empty / no-scheme / unsupported-scheme / no-host / silent-dedupe; features.\<name\>: domain-regex). `Enabled fehlt` (T5-Doctor) und `Source-not-in-Allow` (T4-Enforcement) bewusst deferred. `make gates` grün (Coverage 90.20%). Depguard-Constraint `application-no-net` zwang manuellen URL-Parser statt `net/url` — Code-Größe entsprechend höher als Original-Schätzung. | ~150 geschätzt / **~230 real** (produktion ohne Tests) |
| T2  | **Feature-Katalog.** ✅ Done. `featureCatalogueEntry`-Struct (`source` als OCI-ref-Form ohne Scheme, `defaultVersion`, `shortDesc`) + `featureCatalogue()`-Map mit 8 Spec-Beispielen (`cpp`, `docker-cli`, `git`, `go`, `java`, `kubectl-helm`, `node`, `postgres-client`) + `featureFor`-Lookup. Tests pinnen Spec-Beispiele-Vollständigkeit, Per-Entry-Invarianten (FeatureName-Regex / source-Prefix / non-empty version+desc) und Lookup-Contract (known/unknown). | ~100 geschätzt / **~84 real** (unter Budget) |
| T3  | **Generator-Patch.** ✅ Done. `devcontainerFeatureData`-Projection + `collectDevcontainerFeatures`/`projectFeatureEntry`-Helper (devcontainer_features.go); `templateData.Features` (templates.go); `devcontainer.json.tmpl` um conditional `"features": {…}`-Block erweitert (innerhalb init-Block, kein Nesting); `generateDevcontainer` verdrahtet. Skip-Logik: enabled fehlt/false → skip; unknown catalogue ohne source-Override → skip (T4 enforced den Error). Sort by Source. Tests: 9 Projection-Sub-Cases + 3 End-to-End-Tests (no-features-key-Backwards-Compat, KeysSortedAndShape, Idempotent). | ~200 geschätzt / **~115 real** (deutlich unter Budget) |
| T4  | **CLI-Subkommando.** ✅ Done. M8-Path-Whitelist um 4 neue Kinds erweitert (`devcontainer.featureSources.allow` + `devcontainer.features.<name>.{enabled,source,version}`); `domain.ConfigPath.Feature`-Feld + Parser-Erweiterung. ConfigService: scalar-Pfade via PatchScalar; list-Pfad (`featureSources.allow`) via Marshal-Rewrite (`setFeatureSourcesAllow`); Source-Allowlist-Enforcement in `revalidateFeatureEntry` (Exit-Code 10, `--yes` reicht nicht). `--allow-external-feature-sources <quelle>[,<quelle>...]` per `StringSliceVar` (comma-split + multi-flag cumulate) auf den drei Spec-§714-717-Pfaden (init/generate/config-set) registriert; CLI-Layer rejected den Flag auf nicht-Spec-Pfaden. Init-Use-Case rejected Flag ohne `--devcontainer`. Generate-Use-Case appended Flag-URLs in u-boot.yaml NACH erfolgreichem Plan/Execute (Review-Followup R2). Komplexitäts-Refactoring: `validateInitPreconditions`, `revalidateFeatureEntry`, `extractDevcontainerFeatureValue` als Helper. **Review-Followup R1..R6** in `dc8…` (siehe Commit-Hash unten): R1 `ErrInvalidFeatureSource` nach `domain/` verschoben + in `isValidationError` registriert (Exit-Code 10 für init+generate-Flag-Pfade); R2 atomare Reihenfolge (validate früh, write spät); R3 Orphan-Feature-Info-Log nach config-set-enabled; R4 Allowlist-Membership-Error-Message + Doku-Klärung (trailing-slash/case); R5 Doc-Präzisierung second-normalise; R6 misleading `--yes`-Hinweis entfernt. Tests: 11 neue ConfigService-Tests + 3 Init-Tests + 3 Generate-Tests + 2 ConfigPath-Tests + 2 Exit-Code-CLI-Tests + 1 Atomic-Pin-Test + 3 Orphan-Warning-Tests; Coverage 90.0 %. | ~250 geschätzt / **~580 real** (T4 ~510 + Followup ~70; +132 %; primär durch 4 neue config-Pfade × Stages 1-5 + list-Pfad-Sonderfall + R1-R6-Härtung) |
| T5  | **Doctor-Checks (nur Teil A; Teil B ausgelagert).** ✅ Done. Check-ID `devcontainer.features.allowlist` mit drei Klassifikationen pro Feature-Eintrag: (1) `Source`-Override nicht in Allowlist → **Error** [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)/[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte); (2) Orphan-Activation (Source leer + Name nicht im Katalog) → **Warn**; (3) `Enabled == nil` → **Warn** (analog [`LH-FA-ADD-005`](../../../../spec/lastenheft.md#lh-fa-add-005--mehrfaches-hinzufügen-verhindern) §893). Worst-severity-wins-Aggregation. Skip-Pfade: u-boot.yaml fehlt/unparsable, kein devcontainer-subtree, leere features-Map. Spec §2394 (kein error ohne legitimen Anlass) durch `OKWhenNoFeatures`-Test gepinnt. `classifyFeatureEntries`-Helper extrahiert (gocognit). 8 neue Doctor-Tests + Anpassung der Total-Anzahl 11→12. Teil B (`devcontainer.features.drift`) lebt in [`slice-followup-devcontainer-features-drift-doctor`](slice-followup-devcontainer-features-drift-doctor.md). | ~120 geschätzt / **~163 real** (+36 %; primär durch dreifach-Klassifikation + Worst-severity-wins-Aggregation) |
| T6  | **E2E + Spec-Pin.** ✅ Done. `TestLHFADEV003_CatalogueActivation` (Init→config-set→generate→JSON-Key-Verify für catalogued `node`) + `TestLHFADEV003_AllowlistEnforcement` (negativ: `features.X.source` ohne Allowlist → `ErrConfigValueInvalid` mit [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features)/[`LH-NFA-SEC-004`](../../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)-Message; danach Allowlist-Seed via Spec-§717-Pfad + Wiederholung succeeds). Coverage durch die E2E-Pfade auf 90.4 % gestiegen (von 90.0 %). | ~100 geschätzt / **~144 real** (+44 %; verzweigte Happy+Negative+Recovery-Paths) |
| T7  | **Doku-Closure.** READMEs (EN + DE), `docs/user/devcontainer-features.md`, CHANGELOG `## [Unreleased]`-Eintrag (Conventional-Commit-Stil: feat(devcontainer): Devcontainer-Features-Allowlist und Katalog ([`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features))), Slice `open/` → `done/` **mit DoD-Hash-Line direkt** (kein `git log --grep`-Platzhalter, vgl. `feedback-done-slice-dod-hash`), Roadmap-Status-Update. | — (Doku, nicht im Go-LOC-Budget) |

**Go-LOC-Summe T1-T6 ≈ 920** (geschätzt nach Carveout: T1 150 + T2 100 + T3 200 + T4 250 + T5 120 + T6 100; Teil-B-Schätzung des Folge-Slice ≈ 200 LOC nach Plan-Followup-S1/S2, siehe dort). **Real-LOC nach T5:** T1 ≈ 230 (+53 %), T2 ≈ 84 (−16 %), T3 ≈ 115 (−43 %), T4 ≈ 510 (+104 %), T4-Review-Followup R1..R6 ≈ 70, T5 ≈ 163 (+36 %, dreifach-Klassifikation + Worst-severity-Aggregation). **T1-T5-Summe (in diesem Slice, ohne ausgelagerten Folge-Slice) ≈ 1172 LOC**, davon T1-T4 inkl. Followup ≈ 1009 LOC > 800-Schwelle (Carveout-Trigger gefeuert, Folge-Slice [`slice-followup-devcontainer-features-drift-doctor`](slice-followup-devcontainer-features-drift-doctor.md) (~200 LOC) angelegt). T6 (~100 LOC E2E-Pin) noch offen. Vergleichswerte aus
[`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) +
[`slice-v1-keycloak`](../done/slice-v1-keycloak.md).

## T0-Outcomes

Discovery-Ergebnis (zugleich Plan-Refinement vor `next/`). Vier
Design-Entscheidungen mit kurzer Begründung; abweichende Pfade
sind als Out of Scope notiert oder explizit dokumentiert.

### T0-(a) CLI-Pfad: `config set` + `generate devcontainer`

**Entscheidung:** Aktivierung über bestehendes
`u-boot config set devcontainer.features.<name>.enabled true`
(M8-Path-Whitelist erweitern) plus nachgelagertes
`u-boot generate devcontainer`. Kein neues Top-Level-Kommando
(`add feature` / `devcontainer feature add` verworfen).

**Begründung:** Reuse bestehender M7/M8-Pfade; kein neuer
CLI-Namespace; UX-Konsistenz mit
`config set devcontainer.enabled` aus M8 ([`LH-FA-CONF-005`](../../../../spec/lastenheft.md#lh-fa-conf-005--konfiguration-anzeigen-und-ändern)). Die
Spec-§714-717-Pfade für `--allow-external-feature-sources`
(`init --devcontainer`, `generate devcontainer`,
`config set devcontainer.featureSources.allow`) lassen sich
direkt verdrahten, weil sie bereits existieren — nur das Flag
selbst und ein neuer Whitelist-Pfad in M8 kommen dazu.

**Sub-Entscheidung (Spec-§714-717-Verträglichkeit):** Der
Aktivierungs-Pfad (`config set devcontainer.features.<name>.enabled`)
akzeptiert `--allow-external-feature-sources` **nicht** — das
Flag bleibt strikt auf die drei Spec-Pfade beschränkt
(Default-Variante, „Spec-treu"). Konsequenz: Wer eine externe
Quelle aktivieren will, muss zuerst einen der drei Spec-Pfade
benutzen, um die Allowlist zu pflegen, und kann erst danach
`config set devcontainer.features.<name>.source <url>` + `enabled true`
laufen lassen.

### T0-(b) Schema: Map-analog `services.<name>` mit `*bool`

**Entscheidung:** `devcontainer.features` ist eine Map
`map[string]ubootYAMLDevcontainerFeature` analog
`services: map[string]ubootYAMLService` aus
`initproject.go:54-56`. Per-Entry-Felder:

- `Enabled *bool` (Pointer-Semantik wie `ubootYAMLService.Enabled`:
  `nil` = unset → Doctor-Warn; `&false` = registriert + deaktiviert;
  `&true` = aktiviert).
- `Source string` (`omitempty`) — optionaler External-Source-
  Override; muss in `featureSources.allow` stehen. Fehlt das
  Feld, greift der Katalog-Lookup (Source aus
  `featureCatalogue()[name]`).
- `Version string` (`omitempty`) — optionaler Per-Feature-Pin-
  Override (z. B. `version: "21"` für Java-21). Fehlt das Feld,
  greift `featureCatalogueEntry.defaultVersion`.

**Begründung:** Pointer-Semantik ist im Codebase etabliert
(`ubootYAMLService.Enabled`, `ubootYAMLDevcontainer.Enabled`) und
deckt die Spec-Pflicht „services.<name>.enabled ist immer
explizit zu setzen" (`spec/lastenheft.md:1355`) — die gleiche
Disziplin gilt sinngemäß für Features. Die `Source`+`Version`-
Overrides sind Plan-Erweiterung über die Spec hinaus, aber
nötig, weil sonst externe Quellen nicht aktivierbar wären.

**Beispiel-YAML:**

```yaml
devcontainer:
  enabled: true
  featureSources:
    allow:
      - https://ghcr.io/orgX/features/custom-rust
  features:
    node:
      enabled: true                                # → catalogue lookup
    java:
      enabled: true
      version: "21"                                # → per-feature pin override
    custom-rust:
      enabled: true
      source: https://ghcr.io/orgX/features/custom-rust   # → external, allowlist-required
```

### T0-(c) Katalog + Domain-Type

**Entscheidung:** Statischer Go-Katalog `featureCatalogue() map[string]featureCatalogueEntry`
analog `serviceCatalogue()` aus
`addservice_execute.go:234 ff.`. Struct-Felder:

```go
type featureCatalogueEntry struct {
    source         string // canonical OCI ref, e.g. "ghcr.io/devcontainers/features/node"
    defaultVersion string // pinned version slug, e.g. "1" or "1.2.0"
    shortDesc      string // user-facing description (doctor hints, future `feature list`)
}
```

Erst-Befüllung (Spec-Beispiele aus §698-707):

| Catalogue-Key       | Source                                                   | Default-Version |
| ------------------- | -------------------------------------------------------- | --------------- |
| `git`               | `ghcr.io/devcontainers/features/git`                     | `1`             |
| `docker-cli`        | `ghcr.io/devcontainers/features/docker-outside-of-docker`| `1`             |
| `node`              | `ghcr.io/devcontainers/features/node`                    | `1`             |
| `java`              | `ghcr.io/devcontainers/features/java`                    | `1`             |
| `go`                | `ghcr.io/devcontainers/features/go`                      | `1`             |
| `cpp`               | `ghcr.io/devcontainers/features/cpp`                     | `1`             |
| `kubectl-helm`      | `ghcr.io/devcontainers/features/kubectl-helm-minikube`   | `1`             |
| `postgres-client`   | `ghcr.io/devcontainers/features/postgresql-client`       | `1`             |

(SDKMAN ist in `java` enthalten — die offizielle
`features/java`-Image-Beschreibung listet SDKMAN als
Installations-Backend.)

**Domain-Type `domain.FeatureName`:** Name-only, analog
`domain.ServiceName` aus `domain/servicename.go` (Regex
`^[a-z]([a-z0-9-]{0,30}[a-z0-9])?$`, 32-Char-Cap). Kein
Version-Slot im Domain-Type — Version ist ein Render-/Config-
Concern, nicht ein Identity-Concern. Versionspinning lebt in
`featureCatalogueEntry.defaultVersion` plus optionalem
`ubootYAMLDevcontainerFeature.Version`-Override.

**Lesart des Spec-Begriffs „lokal hinterlegt":** in u-boot
eingebauter statischer Katalog (Go-Mapping). Alternativ-Lesart
„Features im Repo-Pfad" ist bewusst ausgeschlossen und Trigger
für ein eigenes Folge-Slice (siehe Out of Scope).

### T0-(d) Managed-Block: bestehenden `init`-Block erweitern

**Entscheidung:** Der vorhandene `init`-Managed-Block in
`templates/devcontainer/devcontainer.json.tmpl` (heute 13
Zeilen — `name` / `build` / `forwardPorts` / `remoteUser`) wird
um eine `"features": { … }`-Eigenschaft erweitert. **Kein**
zweiter, eigenständiger Managed-Block. **Kein** Nesting von
Managed-Block-Markern innerhalb des JSON-Objekts.

**Begründung:**

- Ein zweiter Marker innerhalb des Objekt-Literals (zwischen
  Property-Paaren) würde zwar in JSONC syntaktisch erlaubt sein
  (Line-Comments dürfen überall stehen, wo Whitespace erlaubt
  ist), aber `managedblock.Find` ist ein Regex-Matcher ohne
  Nesting-Bewusstsein — der äußere `init`-Block-Text würde den
  inneren `features`-Block-Text einschließen, was bei
  Drift-Detection und Re-Render zu doppelten Patches führt.
- Der Single-Block-Ansatz macht jede Feature-Änderung zu einem
  Voll-Re-Render des `init`-Blocks. Block ist klein (≤ 30 Zeilen
  selbst mit acht aktivierten Features), Byte-Vergleich für
  No-Op-Detection ist bereits etabliert (`generate.go:763`).
- Konsistent mit dem M7-`generate`-Design: ein Block pro Datei
  ist die Default-Disziplin; Ausnahmen (mehrere Blöcke) gibt es
  bisher keine in u-boot-Templates.

**Template-Skizze (Pseudo-Pfad):**

```jsonc
// BEGIN U-BOOT MANAGED BLOCK: init
{
  "name": "{{.Name}}",
  "build": {
    "dockerfile": "./Dockerfile",
    "context": "."
  },
{{- if .ForwardPorts }}
  "forwardPorts": [{{range $i, $p := .ForwardPorts}}{{if $i}}, {{end}}{{$p}}{{end}}],
{{- end }}
{{- if .Features }}
  "features": {
{{- range $i, $f := .Features }}
{{- if $i }},{{ end }}
    "{{ $f.Source }}:{{ $f.Version }}": {}
{{- end }}
  },
{{- end }}
  "remoteUser": "vscode"
}
// END U-BOOT MANAGED BLOCK: init
```

Feature-Reihenfolge ist alphabetisch nach `Source`-URL (T3-
Sortierung im Renderer), damit der JSON-Output deterministisch
bleibt und Re-Generation No-Op-fähig ist.

### Offene Sub-Entscheidung (defer-fähig bis T2)

`featureCatalogueEntry.defaultVersion`: heute überall `"1"` —
sobald reale Version-Pins ausgewählt werden (T2), wird die
Tabelle oben aktualisiert. Quelle: Upstream-Tags der
`devcontainers/features`-Repository.

## Out of Scope

- **Eigene/lokale Features** (Custom-Feature im Repo-Pfad statt
  externer Quelle): geht über [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) hinaus; eigenes
  Folge-Slice mit eigenem Trigger.
- **Sprach-spezifische Build-Aktionen** (Gradle-Wrapper anlegen,
  Go-Module-Init, npm-init, …): das ist Template-Job
  (`LH-FA-TPL-*`), nicht Devcontainer-Feature-Job. Klare
  Trennung: u-boot pinnt die Toolchain, das Template (oder der
  User) ruft sie auf.
- **Feature-Version-Updates** (Renovate/Dependabot-Hook für
  Katalog-Pins): nice-to-have; eigener Folge-Slice nach
  Erstauslieferung.
- **Migration bestehender Devcontainer-Configs ohne managed-Block:**
  M3-Re-Init-Disziplin gilt; User mit selbst-gepflegtem
  `devcontainer.json` müssen `--force`/`--backup` benutzen oder
  `features:` manuell migrieren.

## Bezug

- Spec: [`LH-FA-DEV-003`](../../../../spec/lastenheft.md#lh-fa-dev-003--devcontainer-features) ([`spec/lastenheft.md:692`](../../../../spec/lastenheft.md))
  + Schema-Skizze [`spec/lastenheft.md:1340-1353`](../../../../spec/lastenheft.md)
  + Doctor-Pin [`spec/lastenheft.md:2394`](../../../../spec/lastenheft.md).
- ADR: [ADR-0008 §78](../../adr/0008-plugin-system-statisch.md)
  — Allowlist als einziger ausnahme-pflichtiger externer-Code-Pfad.
- Code-Anker: `internal/hexagon/application/initproject.go:66-69`
  (`Future fields:`-Kommentarblock mit `FeatureSources …`).
- Generator-Pfad: `internal/hexagon/application/generate.go:705`
  (devcontainer.json-Mapping) +
  `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  (heute nur `init`-Managed-Block).
- Vorbild-Slices:
  [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md)
  (Add-on-Pattern + `enabled`-State-Machine als Vorbild für die
  `devcontainer.features.<name>.enabled`-Disziplin gemäß
  T0-(b); abweichend: keine Managed-Block-Drift-Zustände, weil
  `devcontainer.json` einfacher zu klassifizieren ist als ein
  Compose-Service-Block),
  [`slice-v1-keycloak`](../done/slice-v1-keycloak.md)
  (Service-Catalogue-Erweiterung um `serviceCatalogueEntry`-Typ
  als Vorlage für `featureCatalogueEntry` gemäß T0-(c)),
  [`slice-m7-generate`](../done/slice-m7-generate.md)
  (Managed-Block-Disziplin in `devcontainer.json`; T0-(d) hat
  Single-Block-Erweiterung gewählt).
- Roadmap:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0.
- Phase: V1, geplant für v0.4.0-Bündelung mit weiteren
  V1-Generators (`u-boot logs`, `--json`/`--dry-run`); deren
  Slice-Pläne existieren noch nicht (Roadmap-Stichworte) und
  folgen separat.
