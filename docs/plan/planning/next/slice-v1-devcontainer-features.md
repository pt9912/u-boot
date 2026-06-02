# Slice V1: Devcontainer-Features / Toolchains (`LH-FA-DEV-003`)

> **Status:** scheduled für v0.4.0 — Spec ✅, ADR-Hook ✅,
> Code-Anker ✅, T0-Discovery ✅ (§T0-Outcomes), Implementation
> ausstehend (T1..T7). Plan-Refinement für `next/` kann jetzt
> ausschließlich auf den T0-Outcomes aufsetzen.

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

- **`LH-FA-DEV-003`** (Priorität V1) listet explizit Git, Docker
  CLI, Node.js, Java, SDKMAN, PostgreSQL Client, Kubernetes Tools
  als Beispiele und legt das Security-Modell fest
  (`--allow-external-feature-sources`,
  `devcontainer.featureSources.allow`, `--yes` reicht nicht).
- **ADR-0008** §78 nennt die Allowlist ausdrücklich als einzigen
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
  vollständige Failure-Tabelle. Diagnose-Code: `LH-FA-DEV-003`,
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
  mit `code LH-FA-DEV-003` / Exit-Code `10` ab. `--yes` reicht
  nicht (`LH-NFA-SEC-004`). `--allow-external-feature-sources`-
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
- ✅ **Doctor-Integration Teil B (über Spec hinaus,
  konditional → im Slice):** Drift-Erkennung —
  `devcontainer.json` enthält die aktivierten Features
  tatsächlich (Managed-Block-Disziplin analog M5/M7). Check-ID:
  `devcontainer.features.drift`. **Carveout-Trigger (T0-LOC-
  Schätzung):** T1-T4 ≈ 700 LOC < 800-Schwelle → Teil B bleibt
  in T5. Re-Check vor T4-Start (siehe LOC-Hinweis unter den
  Tranchen): wenn die Schätzung dann auf > 800 LOC läuft, wird
  Folge-Slice `slice-followup-devcontainer-features-drift-doctor`
  angelegt und Teil B ausgelagert.
- ✅ **Statischer Katalog** (ADR-0008-konform): u-boot bringt
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
  deckt `LH-FA-DEV-003` ab. Test-Naming-Konvention analog zum
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
| T1  | **Schema-Erweiterung.** `ubootYAMLDevcontainer.FeatureSources` aus dem `Future fields`-Kommentar in `initproject.go:66-69` ziehen; zusätzlich `ubootYAMLDevcontainerFeature` (Map-Value gemäß T0-(b)) mit `Enabled *bool`, optional `Source string` (external-source-Override) und `Version string` (per-Feature-Pin-Override). YAML-Codec + Schema-Validierung. **Failure-Tabelle (`featureSources.allow`-Einträge):** leerer Quell-String, fehlendes URL-Scheme (`http://`/`https://`/`oci://`), fehlende Host-Komponente, doppelter Eintrag (silent-dedupe gemäß `spec/lastenheft.md:1352`). **Failure-Tabelle (`features.<name>`):** `name` verletzt `domain.FeatureName`-Regex (analog `ServiceName` aus `domain/servicename.go`), `Enabled` fehlt (Doctor-Warn, kein Error), `Source` gesetzt aber nicht in Allowlist (Error, Exit-Code `10`, erst in T4 enforced). Tests in `application/*_test.go`. Inkludiert in T1: Domain-Type `domain.FeatureName` (analog `ServiceName`, name-only). | ~150 |
| T2  | **Feature-Katalog.** Statischer Go-Katalog `featureCatalogue() map[string]featureCatalogueEntry` (analog `addservice_execute.go:234 ff. serviceCatalogue()`) mit Spec-Beispielen (Git, Docker CLI, Node, Java/SDKMAN, Go, C++, K8s-Tools, PostgreSQL-Client). Quell-URLs gepinnt (`ghcr.io/devcontainers/features/<name>`); `defaultVersion` pro Eintrag. Struct-Felder: `source string`, `defaultVersion string`, `shortDesc string`. Tests: Katalog-Lookup, `domain.FeatureName`-Roundtrip pro Eintrag. | ~100 |
| T3  | **Generator-Patch.** `devcontainer.json.tmpl` — bestehender `init`-Managed-Block (T0-(d)) wird um `"features": { … }`-Block erweitert (innerhalb des init-Blocks, keine Nesting-Marker). Render-Side: Katalog-Lookup pro aktiviertem Feature (oder `source:`-Override-Lookup gegen Allowlist), deterministische Reihenfolge (alphabetisch nach Source-URL), JSON-Determinismus. `templateData`-Struct in `generate.go` um `Features []devcontainerFeature` erweitert. `generate devcontainer` idempotent (M7-Managed-Block-Disziplin gilt). | ~200 |
| T4  | **CLI-Subkommando.** Kein neues Top-Level-Kommando — Aktivierung läuft ausschließlich über bestehendes `u-boot config set devcontainer.features.<name>.enabled true` (M8-Path-Whitelist erweitern) + nachgelagertes `u-boot generate devcontainer`. Plus M8-Whitelist um `devcontainer.featureSources.allow` und `devcontainer.features.<name>.source`/`.version` erweitern. **`--allow-external-feature-sources <quelle>[,<quelle>...]`-Parser** (gültig für die drei Spec-Pfade aus §714-717: `init --devcontainer`, `generate devcontainer`, `config set devcontainer.featureSources.allow`): Komma-Trennung, Whitespace-Toleranz (trim pro Element), Multi-Flag-Vorkommen kumulieren (kein last-wins, analog `--with-deps`-Pattern aus M5), Duplikate gegen bestehende Allowlist silent-dedupen (`spec/lastenheft.md:1352`). **Allowlist-Enforcement:** Aktivierung eines Features mit `source:`-Override, dessen Wert nicht in `featureSources.allow` steht, bricht mit `code LH-FA-DEV-003`/Exit-Code `10` ab. `--yes` reicht nicht (`LH-NFA-SEC-004`). | ~250 |
| T5  | **Doctor-Checks (Teil A + B im Slice; Carveout-Trigger entfällt).** T1-T4-Summe ≈ 700 LOC < 800-Schwelle → Teil B bleibt drin. Check-ID `devcontainer.features.allowlist` (Spec-mandatiert): `devcontainer.features.<name>.source` referenziert nur Quellen aus `featureSources.allow` *oder* `features.<name>` ist Katalog-Eintrag (kein `source:`-Override). Check-ID `devcontainer.features.drift` (über Spec hinaus): `devcontainer.json` enthält die aktivierten Features tatsächlich (Managed-Block-Disziplin analog M5/M7). Severity-Klassifikation analog M5. | ~270 (A ~120 + B ~150) |
| T6  | **E2E + Spec-Pin.** `internal/hexagon/application/acceptance_test.go`: `TestLHFADEV003_CatalogueActivation` (positiver Pfad — `config set devcontainer.features.node.enabled true` + `generate devcontainer` + JSON enthält `features."ghcr.io/devcontainers/features/node": {…}`) + `TestLHFADEV003_AllowlistEnforcement` (negativer Pfad — `features.<name>.source: <url>` ohne Allowlist-Eintrag → Exit-Code `10`). Docker-e2e in `internal/e2e/` nicht erforderlich (kein Laufzeit-Stack zu prüfen). | ~100 |
| T7  | **Doku-Closure.** READMEs (EN + DE), `docs/user/devcontainer-features.md`, CHANGELOG `## [Unreleased]`-Eintrag (Conventional-Commit-Stil: `feat(devcontainer): Devcontainer-Features-Allowlist und Katalog (LH-FA-DEV-003)`), Slice `open/` → `done/` **mit DoD-Hash-Line direkt** (kein `git log --grep`-Platzhalter, vgl. `feedback-done-slice-dod-hash`), Roadmap-Status-Update. | — (Doku, nicht im Go-LOC-Budget) |

**Go-LOC-Summe T1-T6 ≈ 1070** (T1 150 + T2 100 + T3 200 + T4 250 + T5 270 + T6 100), davon **T1-T4 ≈ 700 LOC < 800-Schwelle** → Doctor-Teil-B-Carveout-Trigger feuert nicht; Teil B bleibt in T5. Vergleichswerte aus
[`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) +
[`slice-v1-keycloak`](../done/slice-v1-keycloak.md). Schätzungen ±25 %; bei T1-Abschluss wird die Schwellen-Frage neu geprüft (Re-Check vor T4-Start: wenn T1+T2+T3 realer LOC > 600, dann T4-Schätzung neu kalibrieren — bei dann projizierter T1-T4 > 800 LOC Folge-Slice `slice-followup-devcontainer-features-drift-doctor` in `open/` anlegen *vor* T5-Start).

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
`config set devcontainer.enabled` aus M8 (LH-FA-CONF-005). Die
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
  externer Quelle): geht über `LH-FA-DEV-003` hinaus; eigenes
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

- Spec: `LH-FA-DEV-003` ([`spec/lastenheft.md:692`](../../../../spec/lastenheft.md))
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
