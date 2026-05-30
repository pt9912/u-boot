# Slice M8: `u-boot config`-Flow

> **Status:** Open
> **DoD:** T1 ⬜ / T2 ⬜ / T3 ⬜ / T4 ⬜ / T5 ⬜

## Auslöser

Nach MVP-Closure ist `u-boot config` der **letzte MVP-blockierende
Slice** (siehe Roadmap §Nächste Schritte / MVP-Bilanz). Spec-Pflicht
(alle MVP-Priorität, `spec/lastenheft.md` §4.10):

- **`LH-FA-CONF-001`** Projektkonfiguration über `u-boot config
  get project.name` / `u-boot config set project.name <wert>`
  pflegbar.
- **`LH-FA-CONF-003`** Konfiguration lesen + bei Befehlen
  berücksichtigen (read-Pfad ist seit M3 implementiert; das hier
  bestätigt nur den User-CLI-Zugriff).
- **`LH-FA-CONF-004`** Konfiguration aktualisieren bei Add-on-
  Änderungen — durch M5 (`u-boot add`) bereits abgedeckt; M8
  liefert die ergänzende User-direkte-Mutation via `set`.
- **`LH-FA-CONF-005`** `u-boot config` ohne Argumente zeigt die
  gesamte Konfiguration; `get` ein einzelner Wert; `set <pfad>
  <wert>` muss die Schema-Konformität prüfen.

Out of Scope (Later):

- **`LH-FA-CONF-006`** `u-boot config migrate` — Priorität
  „Later", eigener Slice mit Schema-Versionierung.

## Design-Entscheidungen (kanonisch, vor Implementierung
unterschrieben)

### D1 — Set-Pfad-Whitelist statt freie YAML-Pfade

Ein freier YAML-Pfad würde User-Schreiben überall hin erlauben
(`set foo.bar.baz xxx`), aber die LH-FA-CONF-002-Schema-Pflicht
verlangt, dass nur die dort definierten Felder existieren. **Set
akzeptiert ausschließlich Pfade aus einer Whitelist**:

| Pfad                              | Typ     | Domain-Validierung                                     |
| --------------------------------- | ------- | ------------------------------------------------------ |
| `project.name`                    | string  | `domain.NewProjectName` (LH-FA-INIT-006 regex)         |
| `services.<svc>.enabled`          | bool    | `domain.NewServiceName(svc)` + bool-Parse              |
| `devcontainer.enabled`            | bool    | bool-Parse                                             |

Unbekannte Pfade ⇒ `ErrConfigPathUnknown` → Exit-Code 10. Die
Liste ist eng absichtlich: V1-Felder (`services.keycloak.persistence`,
`devcontainer.featureSources.allow`) kommen erst mit den
entsprechenden Add-on-Slices in die Whitelist; `schemaVersion`
ist read-only (Migration ist Later).

`<svc>` ist ein variabler Wildcard-Segment; der Domain-Validator
([`domain.NewServiceName`](../../../../internal/hexagon/domain/servicename.go))
prüft Format und Catalogue-Mitgliedschaft.

### D2 — Set-Operation ist `PatchScalar`-only

M5-T3 hat den `YAMLCodec.PatchScalar`-Port eingeführt. M8 reused
ihn 1:1. Mehrwertige Sets (`set services.postgres '{enabled:
true, version: 16}'`) sind out-of-scope; jedes Setting ist
genau ein Skalar.

### D3 — Schema-Validierung post-write via Roundtrip

Nach jedem `set` wird die geschriebene `u-boot.yaml` re-gelesen
und in das `ubootYAMLConfig`-Struct unmarshalled. Bricht die
Unmarshal ab, gilt das als Schema-Verletzung ⇒ `ErrConfigSchemaInvalid`
→ Exit-Code 10 mit Hint („restore manually or run `u-boot
config get <path>` to inspect"). Der Patch selbst wird als
Plan-vor-Write klassifiziert: erst Plan-and-Validate auf einer
In-Memory-Kopie, dann commit. Bei Schema-Verletzung wird die
existierende Datei **nicht** überschrieben (atomar).

### D4 — `config get` returnt YAML-formatierten Skalar

`u-boot config get project.name` druckt nur den nackten Wert auf
stdout (`my-service\n`). Keine YAML-Anführungszeichen, keine
JSON-Quoting. Trailing Newline genau einer (analog `echo`).

### D5 — `config` (show) reicht den Datei-Inhalt durch

`u-boot config` ohne Subcommand reicht den `u-boot.yaml`-Inhalt
byte-identisch auf stdout. Das ist deutlich einfacher als ein
Reformat-Pfad und User-erwartbar (`cat u-boot.yaml`-Ersatz).
Keine Schema-Filterung; Kommentare bleiben sichtbar.

### D6 — Exit-Code-Mapping

- `0` — Erfolg.
- `2` — CLI-Validierung: unbekannter Subcommand, fehlende
  positional args (Cobra), `set` mit zu vielen/zu wenigen
  Argumenten.
- `10` — fachlicher Validierungsfehler:
  - `ErrProjectNotInitialized` (kein `u-boot.yaml`).
  - `ErrConfigPathUnknown` (nicht in der D1-Whitelist).
  - `ErrConfigValueInvalid` (domain-Validierung der Value
    schlug fehl, z. B. ungültiger Project-Name, ungültige
    bool-Repräsentation).
  - `ErrConfigSchemaInvalid` (Schema-Roundtrip nach Set
    bricht).
- `14` — technischer FS-Fehler beim Lesen/Schreiben
  (`ErrConfigFileSystem`, analog zu `ErrGenerateFileSystem`
  aus M7).

## Vorbereitende Slices (Status)

- [`slice-m3-init-flow`](../done/slice-m3-init-flow.md) —
  `u-boot.yaml`-Schema (`ubootYAMLConfig`) + initialer Schreib-
  Pfad. M8 liest und mutiert dasselbe Struct.
- [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) —
  `YAMLCodec.PatchScalar`-Port (M5-T3). M8 reused den Port
  ohne Erweiterung.
- [`slice-v1-yaml-parse-error-sentinel`](../done/slice-v1-yaml-parse-error-sentinel.md) —
  `driven.ErrYAMLParse`-Sentinel. M8 nutzt ihn für
  `ErrConfigSchemaInvalid`-Routing (Schema-Roundtrip-Parse-
  Fehler → Code 10).

## Architektur-Punkte

- **Neuer Driving-Port `ConfigUseCase`** in
  `internal/hexagon/port/driving/config.go`:
  ```go
  type ConfigGetRequest struct { BaseDir, Path string }
  type ConfigGetResponse struct { Path, Value string }

  type ConfigSetRequest struct { BaseDir, Path, Value string }
  type ConfigSetResponse struct { Path, OldValue, NewValue string }

  type ConfigShowRequest struct { BaseDir string }
  type ConfigShowResponse struct { Body []byte }

  type ConfigUseCase interface {
      Get(context.Context, ConfigGetRequest) (ConfigGetResponse, error)
      Set(context.Context, ConfigSetRequest) (ConfigSetResponse, error)
      Show(context.Context, ConfigShowRequest) (ConfigShowResponse, error)
  }
  ```
  Drei Methoden statt einer mit Action-Enum, weil die Request-
  Shapes sich strukturell unterscheiden (Show hat keinen Pfad)
  und die Cobra-Subkommandos sich sauberer 1:1 abbilden lassen.
  Sentinels: `ErrConfigPathUnknown`, `ErrConfigValueInvalid`,
  `ErrConfigSchemaInvalid`, `ErrConfigFileSystem`.
  `ErrProjectNotInitialized` reuse aus M5.

- **Neuer Application-Service `ConfigService`** in
  `internal/hexagon/application/config.go`:
  - DI: `fs driven.FileSystem`, `yaml driven.YAMLCodec`,
    `logger driven.Logger`.
  - Top-Level-Dispatcher pro Methode + zwei Helper:
    `validateConfigPath(path string) (configPath, error)` —
    Whitelist-Lookup + Wildcard-Substitution.
    `validateConfigValue(p configPath, raw string) (validated string, error)` —
    domain-spezifische Coercion (`NewProjectName`, bool-parse).
  - Set-Path-Tabelle als Code-Struktur (kein zur Laufzeit
    geparster Datenfile), damit die Whitelist bei Linter-/
    Test-Zeit prüfbar bleibt.

- **Wiederverwendung**: keine neuen Driven-Ports. `PatchScalar`
  ist der einzige Mutator; Get/Show nutzen `ReadFile +
  Unmarshal`.

## Tranchen-Schnitt

Fünf Tranchen, in Reihenfolge implementierbar.

### T1 — Domain `ConfigPath` + Whitelist-Parser

- `internal/hexagon/domain/configpath.go`:
  - `ConfigPath` als typed-string mit `Kind ConfigPathKind`
    (`ConfigProjectName` / `ConfigServiceEnabled` /
    `ConfigDevcontainerEnabled`) und `Service domain.ServiceName`
    (nur für `ConfigServiceEnabled` populated).
  - `NewConfigPath(raw string) (ConfigPath, error)`-Konstruktor
    mit der D1-Whitelist.
  - `ErrInvalidConfigPath` als domain-sentinel.
- Tests: jede erlaubte Pfad-Form roundtrip, alle bekannten
  Fehlpfade (`unknown.path`, `services..enabled`,
  `services.invalid-service-name.enabled`).

**DoD T1:**
- [ ] `domain.ConfigPath` + Whitelist + 100 % Coverage.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T1 ✅ <commit-hash>`.

### T2 — Driving-Port + Application-Skeleton

- `internal/hexagon/port/driving/config.go`: 3 Request-/3
  Response-Structs, `ConfigUseCase`-Interface, 4 Sentinels.
- `internal/hexagon/application/config.go`:
  - `ConfigService`-Struct + `NewConfigService`-Constructor.
  - Get/Set/Show-Methoden mit Stub-Body (alle returnen
    `errStubConfigHandler`-Pattern aus M7-T1, paket-intern,
    `ErrStubConfigHandlerForTest`-Export für Tests).
  - Project-State-Gate (kein `u-boot.yaml` ⇒
    `ErrProjectNotInitialized`) in einem Helper, von allen
    drei Methoden geteilt.
- Tests: Stubs feuern erwarteten Stub-Pin; Project-not-init-
  Pfad pinnt Sentinel für alle drei Methoden.

**DoD T2:**
- [ ] Port + Skeleton kompilieren; Stub-Pin in alle drei
  Methoden integriert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T2 ✅ <commit-hash>`.

### T3 — `Get` + `Show`

- `Get`: Read u-boot.yaml + Unmarshal in `ubootYAMLConfig` +
  Switch auf `ConfigPath.Kind` + Wert-Extraktion + String-
  Formatierung. Fehlende optionale Felder (z. B.
  `devcontainer.enabled` wenn `devcontainer:` fehlt) liefern
  einen klar formulierten Sentinel mit Hint („run `u-boot
  config set devcontainer.enabled false` first").
- `Show`: ReadFile + byte-identische Rückgabe in
  `ConfigShowResponse.Body`. Kein Re-Parse.
- Tests:
  - Get: drei Whitelist-Pfade roundtrip, ein Unknown-Path-Path.
  - Get: missing-optional-field-Pfad für `devcontainer.enabled`.
  - Show: Body matches the disk content byte-identisch
    (inkl. Kommentare).

**DoD T3:**
- [ ] Get/Show-Handler vollständig.
- [ ] Stub-Pin auf 1 (`Set`) reduziert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T3 ✅ <commit-hash>`.

### T4 — `Set` mit Schema-Roundtrip-Validierung

- Wert-Coercion: für `project.name` über `domain.NewProjectName`;
  für `*.enabled` über `strconv.ParseBool` (akzeptiert
  `true/false/0/1/T/F/...` gemäß Go-Standard) und Re-Marshal
  als `true`/`false`-String.
- Pre-Patch: Read → Patch in Memory → Re-Unmarshal in
  `ubootYAMLConfig` → bei Fehler `ErrConfigSchemaInvalid`,
  keine Datei-Mutation. Bei Erfolg WriteFile.
- `OldValue`/`NewValue` in der Response werden aus dem Pre-/
  Post-Snapshot extrahiert (für CLI-Summary).
- Schema-Validation-Implementation: das Unmarshal in
  `ubootYAMLConfig` ist die Schema-Pflicht. Wenn der V1-yaml-
  parse-Sentinel sich gewrappt wird (driven.ErrYAMLParse), wird
  er hier auf `ErrConfigSchemaInvalid` umgemappt — gleicher
  Exit-Code, klarere User-Message.
- Tests:
  - Set project.name happy + invalid name (z. B. "Demo-Project").
  - Set services.postgres.enabled = false (roundtrip).
  - Set unbekannter Pfad ⇒ ErrConfigPathUnknown.
  - Set value falsche Form (z. B. `bool` mit `"vielleicht"`) ⇒
    ErrConfigValueInvalid.
  - Set führt zu einem ungültigen YAML-Schema (über
    Patch-Test-Hook?) ⇒ ErrConfigSchemaInvalid + keine
    Datei-Mutation (writesBefore == writesAfter).
- LH-FA-CONF-005-Pin: das End-to-End-Test reproduziert exakt
  `u-boot config get project.name → u-boot config set
  project.name foo → u-boot config get project.name` und
  asserted dass der gesetzte Wert beim nächsten Get sichtbar
  ist.

**DoD T4:**
- [ ] Set-Handler vollständig + Schema-Roundtrip.
- [ ] Stub-Pin entfernt (alle drei Handler real).
- [ ] LH-FA-CONF-005-Pin im acceptance_test.go ergänzt.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T4 ✅ <commit-hash>`.

### T5 — CLI-Subkommando + ExitCode-Wiring + Doku-Update

- `internal/adapter/driving/cli/config.go`:
  - `newConfigCommand(a)` mit drei Sub-Subkommandos:
    - `u-boot config` (ohne Args) ⇒ Show.
    - `u-boot config get <path>` ⇒ Get.
    - `u-boot config set <path> <value>` ⇒ Set.
  - Cobra-Args-Validierung: `ExactArgs(1)` für `get`, `ExactArgs(2)` für `set`.
- `cli.go`-Erweiterung: `App.configUseCase` Field, `cli.New`-Signatur-Bruch
  (analog M7-T6). Alle Test-Helper in `cli_test.go` migrieren (heute
  6 Slots, neu 7).
- `cmd/uboot/main.go`: `application.NewConfigService` wiring.
- ExitCode-Mapping in `cli.go`:
  - `ErrConfigPathUnknown`, `ErrConfigValueInvalid`,
    `ErrConfigSchemaInvalid` in `isValidationError` → Code 10.
  - `ErrConfigFileSystem` in `isFilesystemError` → Code 14.
  - Cobra-Args-Fehler → Code 2 über bestehende `isUsageError`.
- Doku:
  - `u-boot config --help` listet die drei Subkommandos.
  - `docs/user/README.md` Zähler von „sechs verdrahteten
    Subcommands" auf „sieben" + `config` ergänzen.
- Tests (`cli_test.go`):
  - `TestExecute_Config_Show_PrintsFile` — Show-Output ist die
    Datei.
  - `TestExecute_Config_Get_ProjectName` — Get returnt Wert.
  - `TestExecute_Config_Set_ProjectName_Roundtrip` — Set +
    Get roundtrip.
  - `TestExecute_Config_UnknownPath_Code10`
  - `TestExecute_Config_SetSchemaInvalid_NoWrite_Code10`
  - `TestExitCode_ConfigSchemaInvalid_MapsTo10`,
    `TestExitCode_ConfigFileSystemError_MapsTo14`
    (analog zu M7-T6 named tests).
- Carveouts: kein offener Eintrag zu LH-FA-CONF-* heute (M8 war
  als „Open"-Phase in der Roadmap, kein Carveout); erwartungsgemäß
  no-op.

**DoD T5:**
- [ ] CLI-Subcommand verfügbar, alle drei Pfade gepinnt.
- [ ] `cli.New`-Signatur-Migration: `main.go` + 6 Test-Helper.
- [ ] ExitCode-Mapping + benannte Tests grün.
- [ ] Roadmap M8 → Done + **MVP komplett** Block ergänzen.
- [ ] Slice nach `done/`, DoD-Line `T5 ✅ <commit-hash>`.
- [ ] `make gates` grün.

## Akzeptanzkriterien (Slice-übergreifend)

### Struktur

- Keine neuen Driven-Ports. `PatchScalar` und `Unmarshal` decken
  alle Mutationen / Reads.
- Neue Sentinels (4): `ErrConfigPathUnknown`,
  `ErrConfigValueInvalid`, `ErrConfigSchemaInvalid`,
  `ErrConfigFileSystem` — alle drei `driving`-paket-resident,
  CLI-mapbar.
- Whitelist-Tabelle in einem Helper, nicht runtime-mutable.

### Verhalten

- **LH-FA-CONF-001/-005**: `u-boot config get/set/show`
  funktionieren wie spec'd.
- **LH-FA-CONF-002**: Schema-Konformität nach Set geprüft;
  Schema-Verletzung ⇒ Datei unverändert.
- **LH-FA-CONF-003**: read-Pfad reused den existierenden
  M3-Read-Pfad.
- **LH-FA-CONF-004**: Set ist die direkte User-Mutation;
  `add`/`remove` sind die fachlichen Mutationen.

### Negative

- Kein `u-boot.yaml` ⇒ Exit 10 (`ErrProjectNotInitialized`).
- Unbekannter Pfad ⇒ Exit 10 (`ErrConfigPathUnknown`).
- Invalider Wert (z. B. Project-Name mit ungültigem Zeichen)
  ⇒ Exit 10 (`ErrConfigValueInvalid`).
- Schema-Verletzung post-set ⇒ Exit 10
  (`ErrConfigSchemaInvalid`) + null WriteFile-Mutation
  (atomar).
- FS-Fehler ⇒ Exit 14 (`ErrConfigFileSystem`).

## Out of Scope (Later / V1)

- **`LH-FA-CONF-006`** Migration (`u-boot config migrate`) —
  Priorität „Later" in der Spec.
- **Nested set** (`set services '{postgres: {enabled: true}}'`)
  — kein Spec-Mandat; M5-Path-Granularität reicht.
- **Listen-Operationen** (`set
  devcontainer.featureSources.allow '[a, b]'`) — V1-Feld, wartet
  auf den entsprechenden V1-Slice.
- **JSON-Output** für `show` / `get` — analog M4/M6/M7-
  Entscheidung; Text-Output zuerst.
- **`config validate`** als eigener Subkommando — heute
  unnötig, weil `set` post-write validiert.

## Bezug

- Auslösende Spec: `spec/lastenheft.md` §4.10
  (`LH-FA-CONF-001..005`).
- Hängt von: M3 (`u-boot.yaml`-Schema), M5
  (`YAMLCodec.PatchScalar`), V1-yaml-parse-sentinel
  (`driven.ErrYAMLParse` für Schema-Roundtrip-Klassifikation).
- Phase: M8 — letzter MVP-blockierender Slice. Nach Done ist
  MVP vollständig (alle 5 LH-AK-* gepinnt + alle MVP-LH-FA-*
  ausgeliefert).
- Roadmap: ersetzt `M8 | Open` durch `M8 | Done` mit Slice-Link;
  der MVP-Bilanz-Block in „Nächste Schritte" wird auf „MVP
  vollständig" geflippt.
