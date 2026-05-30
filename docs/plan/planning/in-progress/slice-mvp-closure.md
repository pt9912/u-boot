# Slice MVP-Closure: `init --devcontainer` + LH-AK-Pin-Vervollständigung

> **Status:** In progress (T1 + T2 done)
> **DoD:** T1 ✅ `bfe6416` / T2 ✅ `8525c4c` (inkl. Doctor-Severity-Fix für `compose.yaml.valid` no-services: Error → Warn, LH-AK-001 §2299-Konformität) / T3 ⬜

## Auslöser

Nach M7 sind alle MVP-`LH-FA-*`-Anforderungen geliefert (`init`,
`doctor`, `add postgres`, `up`, `down`, `generate`). Die Roadmap-
Zeile `MVP-Closure | Open` zieht die letzten zwei Lücken zu
einem Slice zusammen:

1. **`LH-AK-005` Devcontainer-Flow** verlangt `u-boot init
   --devcontainer`. M7-T5 lieferte zwei Templates und den
   `generate devcontainer`-Pfad, ließ den Flag aber explizit als
   MVP-Closure-Scope offen (siehe slice-m7-generate.md §"Out of
   Scope" + §Architektur-Punkt "Block-Name in allen generierten
   Dateien").
2. **`LH-AK-001`** (`init && doctor`) und **`LH-AK-006`**
   (Doppel-`add postgres`) haben keine dedizierten e2e-Test-Pins
   mit dem Spec-Wortlaut. Die zugrundeliegende Semantik ist
   getestet (Init-Tests + `TestAdd_ActiveWithAllArtifactsIsNoOp`),
   aber kein einziger Test trägt einen `LHAK00X`-Namen für
   diese zwei IDs.

Damit ist die MVP-Akzeptanzkriterien-Matrix nach diesem Slice
geschlossen: alle fünf MVP-`LH-AK-*` (001, 002, 005, 006, 007)
haben einen benannten e2e-Pin (002 + 007 sind seit M6 / M7
gepinnt; 001/005/006 ergänzt dieser Slice).

Spec-Pflicht (alle MVP-Priorität):

- **`LH-AK-001`** Vorbedingung Docker-Engine erreichbar;
  `mkdir demo && cd demo && u-boot init && u-boot doctor`
  liefert „keinen `error`-Eintrag" und „vorhandene Dateien wurden
  nicht ungewollt überschrieben" (`spec/lastenheft.md` §2281).
- **`LH-AK-005`** `u-boot init --devcontainer` ⇒ beide
  Devcontainer-Dateien existieren, `devcontainer.json` enthält
  `name` + `build`/`image`, `forwardPorts` falls aktive Ports,
  `u-boot doctor` ohne `error` zu Devcontainer-Konfig (§2367).
- **`LH-AK-006`** `u-boot add postgres && u-boot add postgres`
  ⇒ Postgres ist genau einmal in der Konfig, verständliche
  Meldung (§2387).

Out of Scope (V1+):

- **`LH-AK-003`** Keycloak-Flow / **`LH-AK-004`** OTel-Flow —
  V1-Add-ons, nicht MVP (`LH-FA-ADD-003`/`-004`).
- **`LH-FA-DEV-003`** Externe Feature-Quellen (V1, eigener Flow
  mit `--allow-external-feature-sources`).
- **MVP-Release-Pipeline** — kein MVP-Scope; eigener V1-Slice
  ([`slice-v1-release-pipeline`](../open/slice-v1-release-pipeline.md)).
- **Volle MVP-`LH-FA-*`-Audit** — Spec listet 60+ MVP-IDs; ein
  Eintrag-für-Eintrag-Pin-Check ist eine eigene Disziplin-Übung,
  nicht hier. Dieser Slice fokussiert auf die fünf
  `LH-AK-*`-Acceptance-Flows; latente `LH-FA-*`-Gaps treten erst
  bei realer Verifikation auf.

## Vorbereitende Slices (Status)

- [`slice-m3-init-flow`](../done/slice-m3-init-flow.md) — Init-
  Service, `planFile`/`filePlan`-Pipeline, `actionWrite`/
  `actionReplaceBlock`, `--force`/`--backup`-Mechanik. T1 baut
  darauf auf.
- [`slice-m7-generate`](../done/slice-m7-generate.md) §T5 +
  §Out-of-Scope — Devcontainer-Templates + Risiko-Hinweis für
  T1: dieselben Marker-Namen (`init`) zwischen `init
  --devcontainer` und `generate devcontainer`, deshalb muss der
  Init-Pfad eine `actionReplaceBlock`-Kompatibilität mit einer
  bereits durch `generate devcontainer` angelegten Datei
  anbieten — sonst kollidiert er sofort.
- [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) —
  Add-Service-State-Machine mit Idempotenz-Pfad
  (`Active → Active, Changed=nil`). T2 reused den existierenden
  `TestAdd_ActiveWithAllArtifactsIsNoOp` als Backbone und
  ergänzt nur den `LHAK006`-benannten Wrapper-Test.

## Architektur-Punkte

- **`InitProjectRequest` bekommt `Devcontainer bool`-Field.**
  CLI-Adapter setzt den Wert aus dem neuen `--devcontainer`-
  Flag. Default `false` ⇒ M3-Verhalten unverändert.
- **`InitProjectService.Init`-Erweiterung (T1).** Wenn
  `req.Devcontainer == true`, werden die zwei
  Devcontainer-Templates (`devcontainer/devcontainer.json.tmpl`
  + `devcontainer/Dockerfile.tmpl`) in die `planFile`-Schleife
  aufgenommen. Reuse der bestehenden `--force`/`--backup`-
  Disziplin: existierende Dateien mit `init`-Marker bekommen
  `actionReplaceBlock`, ohne Marker `actionAbort` (oder
  `--force --backup` für Splice).
- **`u-boot.yaml.devcontainer.enabled` wird auf `true` gesetzt
  (T1).** Mirror zu `services.<name>.enabled` aus M5; setzt das
  Doctor-Severity-Eskalations-Gate aus M5-T7 scharf
  (`LH-FA-DIAG-002` §1073: `devcontainer.enabled == true` ⇒
  `.devcontainer/devcontainer.json` Syntax-Check ist `error`).
- **`forwardPorts` bei `init --devcontainer` ist leer.** Beim
  ersten Init existieren noch keine aktiven Services; die
  Render-Daten-Struktur (`templateData.ForwardPorts: nil`)
  greift dieselbe LH-FA-DEV-005-Optionalität wie `generate
  devcontainer` (Feld fehlt im JSON). Erst nach `u-boot add
  postgres` und einem späteren `u-boot generate devcontainer`
  taucht der Port auf — das ist konsistent mit der M7-T5-
  Detection-Logik und braucht keinen Sonderpfad.
- **Kollision mit `generate devcontainer`-Output.** Wenn der
  User zuerst `u-boot generate devcontainer` (ohne `init
  --devcontainer`) lief und dann `u-boot init --devcontainer`,
  finden die Init-Helper bereits Dateien mit dem `init`-Marker
  vor. Per M3-Konvention triggert das **ohne** `--force` einen
  `ErrFileExists` → Code 10. Mit `--force --backup` läuft der
  Block-Replace; Inhalt außerhalb des Blocks bleibt
  byte-identisch. Diese Reihenfolge bekommt einen eigenen Test
  in T1.

## Tranchen-Schnitt

Drei Tranchen, in Reihenfolge implementierbar. T1 ist die
substantielle Arbeit; T2 + T3 sind kurze Doku-/Test-Pin-
Stücke.

### T1 — `u-boot init --devcontainer` (LH-AK-005)

- `internal/hexagon/port/driving/initproject.go`:
  `InitProjectRequest.Devcontainer bool`-Field hinzufügen,
  Default `false`.
- `internal/hexagon/application/initproject.go`:
  - `fileTemplates()` ⇒ kein neuer Eintrag. Stattdessen ein
    neuer Helper `devcontainerFileTemplates()`, der die zwei
    Devcontainer-Templates returnt (Style + Path + Template-
    Name analog zur bestehenden `fileTemplate`-Struktur).
  - `Init()` ruft `devcontainerFileTemplates()` zusätzlich auf,
    wenn `req.Devcontainer == true`; die Pläne werden in der
    bestehenden `planFile`-Schleife mitverarbeitet. Damit gilt
    die volle `--force`/`--backup`-Disziplin **automatisch**
    auch für Devcontainer-Files; kein Sonderpfad.
  - `executeUBootYAML` ⇒ wenn `req.Devcontainer == true`, im
    geschriebenen YAML `devcontainer.enabled: true` setzen.
    Das geht durch das bestehende `ubootYAMLConfig`-Struct;
    `Devcontainer = &ubootYAMLDevcontainer{Enabled: &t}` mit
    `t := true`.
- `internal/adapter/driving/cli/init.go`:
  - Neuer Flag `--devcontainer` (Bool, Default `false`).
  - Wert wird in den `InitProjectRequest` propagiert.
- `internal/adapter/driving/cli/init.go` `Long`-Help:
  Beschreibung des Flags ergänzen (analog `--force`/`--backup`).
- Tests (`internal/hexagon/application/initproject_test.go`):
  - `TestInit_LHAK005_DevcontainerFlow_FreshProject`: fresh
    init mit `--devcontainer` ⇒ beide Dateien existieren, JSON
    ist valid via `stripJSONC + json.Valid`, enthält
    `name`/`build`/`remoteUser`. `u-boot.yaml` hat
    `devcontainer.enabled: true`.
  - `TestInit_LHAK005_DevcontainerFlow_AfterGenerate_ReinitWithForceBackup`:
    seed mit dem Output von `generate devcontainer` (zwei
    Files vorhanden mit `init`-Marker), dann `init
    --devcontainer --force --backup` ⇒ Block re-spliced,
    `.bak`-Files existieren, User-Content außerhalb des Blocks
    bleibt byte-identisch.
  - `TestInit_LHAK005_DevcontainerFlow_ConflictWithoutForce`:
    Seed mit existierenden Devcontainer-Files **ohne**
    `init`-Marker, `init --devcontainer` (ohne `--force`) ⇒
    `ErrFileExists`. Doku-Pflicht: User soll `--force --backup`
    oder File-Rename benutzen.
- CLI-Tests (`internal/adapter/driving/cli/cli_test.go`):
  - `TestExecute_Init_DevcontainerFlag_Propagates`: assert
    dass `--devcontainer` den Request-Field setzt.

**DoD T1:**
- [ ] `InitProjectRequest.Devcontainer` deklariert + dokumentiert.
- [ ] `devcontainerFileTemplates()`-Helper + Init-Verkabelung;
  bestehende Init-Tests grün (Backward-Compat).
- [ ] `u-boot.yaml.devcontainer.enabled = true` beim Init mit
  Flag (Pin via parse-and-check).
- [ ] LH-AK-005-e2e-Pin grün (init + doctor → kein Error für
  Devcontainer-Konfig; M5-T7-Severity-Gate fired correctly).
- [ ] `--devcontainer`-Flag in `u-boot init --help` sichtbar.
- [ ] Kollisions-Test mit Generate-Output dokumentiert (force-
  required-Pfad pinnt das Risiko aus dem M7-T5-OOO).
- [ ] `make gates` grün.
- [ ] DoD-Line: `T1 ✅ <commit-hash>`.

### T2 — LH-AK-001 + LH-AK-006 benannte e2e-Pins

Beide Pins folgen dem Muster von `TestGenerateChangelog_LHAK007_FlowEndToEnd`:
direkter Service-Aufruf im `application_test`-Package, kein
Docker, kein CLI-Layer.

- `TestLHAK001_InitFlow_DoctorClean`:
  - `tempDir := t.TempDir()` (mit explizitem Name analog
    M6-T4-fund).
  - `application.NewInitProjectService(...)` → Init.
  - `application.NewDoctorService(...)` → Check.
  - Assert: `report.HasErrors() == false`.
  - Per Spec: „vorhandene Dateien wurden nicht ungewollt
    überschrieben" — zweite Pflicht. Test seedet KEINE
    existierenden Files (frischer `t.TempDir()`); der negative
    Pfad (existierende Files schützen) ist bereits durch
    M3-T4a/`init --force`-Tests gepinnt. Hier nur der positive
    Spec-Wortlaut.
- `TestLHAK006_DoubleAddPostgres_NoDuplicate`:
  - Init + `add postgres` + `add postgres` (zweimal).
  - Assert: zweiter Aufruf returnt `PriorState=Active`,
    `State=Active`, `Changed=nil`.
  - Assert: `u-boot.yaml` enthält `postgres`-Key nur einmal
    (Service-Map-Größe = 1).
  - Assert: zweiter Aufruf produziert null `WriteFile`-Calls
    (analog T2-NoOp-Pin aus M7).
- Beide Tests landen in einer neuen Datei
  `internal/hexagon/application/acceptance_test.go` —
  dediziert für die LH-AK-Pins, damit ein künftiger LH-AK-003/
  -004-Pin (V1) eine offensichtliche Heimat hat.

**DoD T2:**
- [ ] Beide LH-AK-Test-Funktionen existieren, grün, und tragen
  den `LHAK00X`-Namen-Suffix.
- [ ] `acceptance_test.go` als neue Datei mit klarem
  Package-Header-Kommentar („Spec-pin-tests for the MVP
  acceptance criteria from `spec/lastenheft.md` §9").
- [ ] `make gates` grün.
- [ ] DoD-Line: `T2 ✅ <commit-hash>`.

### T3 — MVP-Closure-Dokumentation

- Slice-Plan auf `Done` flippen, Hash-Liste in DoD.
- `git mv` nach `done/`.
- Roadmap `MVP-Closure | Open` → `Done` mit Slice-Link.
- Roadmap "Nächste Schritte"-Block: nach MVP-Closure folgt die
  V1-Phase. Konkret offene V1-Slices: Plugin-System
  (`LH-OPEN-003`), Template-Format (`LH-OPEN-004`),
  Release-Pipeline (`LH-OPEN-002`), YAML-Parse-Error-Sentinel
  (Review-Followup N2). Diese vier hängen nicht voneinander
  ab; der Triage-Trigger ist „erster externer User-Report" oder
  „erste Release-Vorbereitung".
- Statt eines weiteren Carveout-Cleanups: ein knapper
  MVP-Bilanz-Block in der `## Nächste Schritte`-Sektion der
  Roadmap mit „MVP komplett — alle MVP-`LH-AK-*` gepinnt,
  alle MVP-`LH-FA-*` ausgeliefert".

**DoD T3:**
- [ ] Slice in `done/`, Roadmap-Eintrag Done.
- [ ] MVP-Bilanz-Block in Roadmap mit konkreter
  V1-Trigger-Beschreibung.
- [ ] DoD-Line: `T3 ✅ <commit-hash>`.

## Akzeptanzkriterien (Slice-übergreifend)

### Struktur

- Keine neuen Driven-Ports.
- Keine neuen Sentinels (existierende reuse: `ErrFileExists`
  + `ErrForceRequiresBackup` für die Generate-Kollisions-
  Routen).
- Block-Name `init` bleibt für die Devcontainer-Files (M7-T5
  Architektur-Punkt), damit `init --devcontainer` und
  `generate devcontainer` denselben Block bedienen.

### Verhalten

- **LH-AK-001**: `init + doctor` ohne Error.
- **LH-AK-002**: bereits gepinnt
  (`TestE2E_LHAK002_PostgresAcceptanceFlow`).
- **LH-AK-005**: `init --devcontainer` produziert beide
  Pflichtdateien, `devcontainer.enabled=true` in
  `u-boot.yaml`, JSON valid.
- **LH-AK-006**: Doppel-`add postgres` ⇒ keine Duplikate, NoOp-
  Aktion beim zweiten Aufruf.
- **LH-AK-007**: bereits gepinnt
  (`TestGenerateChangelog_LHAK007_FlowEndToEnd`).

### Negative

- `init --devcontainer` auf existierende Devcontainer-Files
  **ohne** `init`-Marker (z. B. User-handgeschrieben) ⇒
  `ErrFileExists` → Code 10. Test pinnt diesen Pfad.
- `init --devcontainer --force` ohne `--backup` auf Files mit
  init-Marker ⇒ Block-Replace (M3-Konvention); kein Backup.
- `init --devcontainer --backup` ohne `--force` auf Files mit
  init-Marker ⇒ Full-Overwrite mit Backup (M3-Konvention).
- `init --devcontainer --force --backup` auf Files mit init-
  Marker ⇒ Block-Replace + Backup.

## Bezug

- Auslösende Spec: `spec/lastenheft.md` §9 LH-AK-001/005/006,
  §4.3 LH-FA-DEV-001/004/005 (DEV-002 implizit über
  VS-Code-Kompatibilität von `name` + `build`).
- Hängt von: M3 (`actionReplaceBlock`/`--force`/`--backup`),
  M5 (`devcontainer.enabled`-Gate aus T7), M7-T5
  (Devcontainer-Templates + Block-Name-Konvention).
- Phase: MVP-Closure (letzter MVP-Slice vor V1).
- Roadmap: `MVP-Closure | Open` → `Done` nach T3.
