# Slice M7: `u-boot generate <artifact>`-Flow

> **Status:** Done
> **DoD:** T1 ✅ `67fc181` / T2 ✅ `3c5de48` / T3 ✅ `037ab00` / T4 ✅ `19c4110` / T5 ✅ `294e492` / T6 ✅ `d32a733` / Review-Followup ✅ `27de9c5` (9 Findings S1..S4 + N1..N5 aus dem Post-Merge-Review adressiert; siehe `## Review-Followup` unten)

## Auslöser

Nach M3 (`u-boot init`), M4 (`u-boot doctor`), M5 (`u-boot add
postgres`) und M6 (`u-boot up`/`down`) fehlt das letzte MVP-Subkommando
aus §4.8: **`u-boot generate <artifact>`**. Erst damit schließt sich
`LH-AK-007` (Changelog-Generator) und der MVP-Akzeptanz-Pfad für
artefaktbasierte Wiederherstellung / -Aktualisierung.

Spec-Pflicht für M7 (alle MVP-Priorität, `spec/lastenheft.md` §4.8):

- **`LH-FA-GEN-001`** Befehlsstruktur `u-boot generate <artifact>`.
  Erlaubte Werte: `changelog`, `readme`, `env-example`, `devcontainer`.
  Unbekannter Wert ⇒ Exit-Code **2** mit explizit aufgelistetem Katalog
  (analog `domain.NewServiceName` für `add`, aber Exit-Code 2 weil
  CLI-Validierung — `add` mappt unbekannte Service-Namen auf Code 10
  als fachliche Inkonsistenz; `generate` macht das **nicht**, weil die
  Spec explizit Code 2 fordert).
- **`LH-FA-GEN-002`** `u-boot generate changelog` (LH-AK-007).
- **`LH-FA-GEN-003`** `u-boot generate readme`.
- **`LH-FA-GEN-004`** `u-boot generate env-example`.
- **`LH-FA-GEN-005`** Idempotenz: mehrfaches Ausführen erzeugt keine
  Duplikate; bestehende manuelle Inhalte bleiben erhalten; automatisch
  verwaltete Bereiche sind eindeutig markiert.

Plus aus angrenzenden Spec-Punkten, die `generate` einlöst:

- **`LH-FA-DEV-001`** Devcontainer-Mindestdateien
  (`.devcontainer/devcontainer.json` + `.devcontainer/Dockerfile`).
  M7 liefert die Templates und den `generate devcontainer`-Pfad; das
  parallele `init --devcontainer`-Flag (LH-AK-005) bleibt **MVP-
  Closure-Slice**, der die hier eingeführten Templates wiederverwendet.
- **`LH-FA-DEV-004`** Devcontainer mit nicht-root Default-User
  (User wird in der Template-Defaultsektion verankert).
- **`LH-FA-DEV-005`** Ports aus aktiven Services in
  `devcontainer.json`.`forwardPorts`. Quelle: u-boot.yaml-`services`-
  Tree + `compose.yaml`-Ports; MVP liest die Container-Seite der
  Compose-Port-Mappings aus `compose.yaml`-managed-Blocks (siehe T5).
- **`LH-FA-CLI-006`** Exit-Code-Mapping:
  - `2` — CLI-Validierung (unbekanntes Artefakt, fehlende
    positional args).
  - `10` — fachlicher Validierungsfehler
    (`ErrProjectNotInitialized` weil kein `u-boot.yaml`,
    `ErrGenerateManualConflict` wenn Datei ohne managed-Block
    existiert).
  - `14` — technischer Persistenz-/Dateisystemfehler beim Lesen oder
    Schreiben. Code 14 existiert bereits in `cli.ExitCode`
    (`isFilesystemError` in `cli.go:264-267` für
    `ErrBackupSuffixExhausted` + `ErrBackupSourceMissing`,
    `TestExitCode_BaseMappings` in `cli_test.go:682-684`); M7 ergänzt
    in dieser Liste den **neuen** Sentinel `ErrGenerateFileSystem`,
    der FS-Fehler aus `driven.FileSystem` wrappt. Falls die
    Driven-Layer keinen dedizierten Sentinel exportiert
    (`driven.ErrFileSystem*` — Scan ergab: existiert heute nicht),
    bleibt der Wrap; sonst zeigt `errors.Is` direkt auf den
    Driven-Sentinel.
- **`LH-SA-FILE-002`** Managed-Block-Konvention pro Datei-Format:
  - `.env.example` ⇒ `StyleHash` (`# BEGIN ...`).
  - `README.md`, `CHANGELOG.md` ⇒ `StyleHTMLComment` (`<!-- BEGIN ... -->`).
  - `.devcontainer/devcontainer.json` ⇒ `StyleDoubleSlash` (`// BEGIN ...`).
  - `.devcontainer/Dockerfile` ⇒ `StyleHash`.
  Die vier Datei-Mappings nutzen die **drei** im `managedblock`-Paket
  bereits implementierten Stile (`.env.example` und `Dockerfile` teilen
  sich `StyleHash`) — siehe
  `internal/hexagon/application/managedblock/managedblock_test.go` und
  die Style-Enum-Definition in `managedblock.go:31-39`.

Out of Scope (V1+):

- **`LH-FA-TPL-001..004`** Template-System mit benutzerdefinierten
  Projektvorlagen — eigener V1-Slice
  ([`slice-v1-template-format-entscheidung`](../done/slice-v1-template-format-entscheidung.md)).
- **`LH-FA-DOC-005`** Compose-Validierung — M7 generiert keinen
  Compose-Output, nur die vier oben genannten Artefakte.
- **`LH-FA-DEV-003`** Externe Devcontainer-Feature-Quellen
  (`--allow-external-feature-sources`) — V1 mit eigenem Flow,
  hier nur das Mindest-Template.
- **`LH-FA-GEN-001` mit anderen Artefakten** (z. B. `generate
  compose`, `generate dockerfile`) — Spec listet nur die vier oben;
  weitere kommen als eigene Slices.
- **`u-boot init --devcontainer`** (LH-AK-005) — gehört zu
  MVP-Closure, nutzt aber die T5-Templates aus M7.

## Vorbereitende Slices (Status)

- [`slice-m3-init-flow`](../done/slice-m3-init-flow.md) — Template-
  Embed, `renderTemplate`, `planFile`-Tabelle und
  `actionReplaceBlock`-Pfad existieren. M7 reused die gesamte Plan-and-
  Execute-Mechanik aus `initproject.go`.
- [`slice-m5-add-postgres`](../done/slice-m5-add-postgres.md) —
  `PatchScalar` + `PatchMappingEntryYAML` im `YAMLCodec`-Port;
  `renderManagedBlockOnly`-Helper; `managedblock.StyleHash`-Block-
  Konvention für `.env.example` etabliert.
- [`slice-m6-up-down`](../done/slice-m6-up-down.md) — grenzt die
  Port-Semantik ab: `upservice_portparse.go` normalisiert
  **Host-Ports** für TCP-Probes; T5 darf diesen Parser nicht direkt
  für `forwardPorts` wiederverwenden. Referenz für T5 ist stattdessen
  die bestehende Doctor-Logik `devcontainer.forwardPorts.consistency`,
  die **Container-Ports** aus `compose.yaml` ableitet.

## Architektur-Punkte

- **Neuer Driving-Port `GenerateUseCase`** in
  `internal/hexagon/port/driving/generate.go`:
  ```go
  type GenerateRequest struct {
      BaseDir  string
      Artifact domain.Artifact // Changelog | Readme | EnvExample | Devcontainer
  }
  type GenerateResponse struct {
      Artifact domain.Artifact
      Action   GenerateAction // Created | UpdatedBlock | NoOp | RepairedManual
      Changed  []string       // relative paths, deterministisch sortiert
  }
  type GenerateUseCase interface {
      Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)
  }
  ```
  Sentinels:
  - `ErrArtifactUnknown` — Mappt auf Exit-Code 2 (vor dem Use-Case-
    Aufruf in der CLI gefangen, parallel zu `ErrConflictingModeFlags`).
  - `ErrProjectNotInitialized` — Wiederverwendet aus M5; Code 10.
  - `ErrGenerateManualConflict` — Datei existiert, hat aber **keinen**
    managed-Block und damit keinen Reset-Anker (oder
    `managedblock.ErrBlockMalformed`-Fall mit BEGIN ohne END /
    duplicate BEGIN); Code 10 mit Repair-Hint („Benenne die Datei um
    und führe `u-boot generate <artifact>` erneut aus, oder füge den
    formatgerechten BEGIN/END-Marker aus `LH-SA-FILE-002` manuell
    ein"). MVP schreibt nichts, V1-Folge-Slice könnte `--replace`
    ergänzen — Runtime-Fehlermeldungen in M7 erwähnen diesen nicht
    implementierten Flag aber bewusst nicht; siehe „Out of Scope".
  - `ErrGenerateFileSystem` — Wrapt unerwartete IO-/Permissions-
    Fehler aus `driven.FileSystem.ReadFile`/`WriteFile`/`Stat`;
    Code 14. T6 reiht ihn in die existierende `isFilesystemError`-
    Liste (`cli.go:264-267`) ein. Falls die Driven-Layer einen
    passenden Sentinel exportiert (`driven.ErrFileSystem*` — der
    Pre-T1-Scan ergab: existiert heute nicht), könnte der Wrap
    entfallen und das `errors.Is`-Mapping zeigt direkt auf den
    Driven-Sentinel — Entscheidung fällt in T1 nach erneutem
    bestätigenden Scan der Driven-Pakete.

- **Neuer Domain-Type `domain.Artifact`** in
  `internal/hexagon/domain/artifact.go` als enum-String mit
  `NewArtifact(s string) (Artifact, error)`-Konstruktor. Vier Werte:
  `ArtifactChangelog`, `ArtifactReadme`, `ArtifactEnvExample`,
  `ArtifactDevcontainer`. `String()` liefert den CLI-Argumentwert
  (`"changelog"`, …) — synchron mit dem Spec-Katalog für die
  Fehlermeldung.

- **Neuer Application-Service `GenerateService`** in
  `internal/hexagon/application/generate.go`:
  - DI: `fs driven.FileSystem`, `yaml driven.YAMLCodec` (nur für
    Devcontainer-Port-Auslese aus `u-boot.yaml`/`compose.yaml`),
    `logger driven.Logger`.
  - Dispatch nach `Artifact`-Wert. Pro Artefakt ein privater
    Handler-Helper (`generateChangelog`, `generateReadme`,
    `generateEnvExample`, `generateDevcontainer`), damit Tests jeden
    Handler isoliert ansteuern können.
  - Top-Level `Generate(ctx, req)` macht zuerst die Projekt-State-
    Detection (analog `detectServiceState` aus M5): kein `u-boot.yaml`
    ⇒ `ErrProjectNotInitialized`; sonst Dispatch.

- **Wiederverwendung der `renderManagedBlockOnly` + `replaceBlock`-
  Mechanik** aus M5-T4a / M3-T4. M7 introduces keine neue
  Block-Splice-Logik; `generate` ist die zweite Quelle für
  `actionReplaceBlock`-Aufrufe (M3 hatte `init --force`, M7 hat
  `generate <artifact>`). Erweiterung des `planFile`-Helpers nicht
  nötig — `generate` baut seinen eigenen schmaleren Plan ohne
  `--backup`/`--force`-Verzweigung (siehe T1-Vertrag unten).

- **Wiederverwendung der Compose-Port-Detektion** aus
  `application/doctor.go`: T5 ruft die bereits vorhandenen
  package-internen Helper
  `activeServiceNames(cfg ubootYAMLConfig) []string`
  (`doctor.go:846`) und
  `collectActiveServicePorts(fs, yaml, baseDir, services) ([]int, error)`
  (`doctor.go:902`) direkt auf — beide leben im selben Paket
  (`application/`) wie der neue `GenerateService`, deshalb ist keine
  Extraktion in ein Sub-Package nötig. Damit teilen `generate
  devcontainer` und der Doctor-Check
  `devcontainer.forwardPorts.consistency` exakt dieselbe
  Ports-Quelle (sortiert, dedupliziert, Container-Seite, normiert für
  Map-/Scalar-/`host:cnt/proto`-Einträge). DoD-Pin in T5: ein Test
  ruft beide Pfade auf derselben Fixture auf und vergleicht die
  Listen byte-/element-identisch — Drift zwischen Doctor und
  Generator wird damit explizit verboten.

- **Block-Name in allen generierten Dateien: `init`.** Alle vier
  Artefakte verwenden denselben Block-Namen `init`, identisch zu den
  M3-Templates (`# BEGIN U-BOOT MANAGED BLOCK: init`). M7 führt
  bewusst **keinen** dedizierten `devcontainer`-Block-Namen ein,
  damit das spätere `init --devcontainer`-Flag (LH-AK-005, MVP-
  Closure) und das `generate devcontainer` denselben Block reaktivieren
  — sonst entstünden zwei konkurrierende Marker in derselben Datei.

- **`u-boot.yaml`-`generate:`-Tree (bewusst nicht eingeführt).**
  Spec §4.8 verlangt **keine** Konfigurierbarkeit der Generatoren in
  `u-boot.yaml` (im Gegensatz zu `services:` und `devcontainer:`).
  Generate ist heute parameterlos. Ein V1-Folge-Slice könnte
  `generate.<artifact>.enabled` einführen (um Generatoren je nach
  Projekt-Typ zu deaktivieren), aber das ist keine M7-Pflicht.

## Tranchen-Schnitt

Sechs Tranchen, in Reihenfolge implementierbar. T1–T2 sind
Vorarbeit (Port + simpler Erstgenerator), T3–T5 sind die drei
verbleibenden Artefakte mit jeweils eigenem Idempotenz-Profil, T6
ist CLI + Doku + Carveout-Beseitigung.

### T1 — Driving-Port `GenerateUseCase` + Skeleton

- `internal/hexagon/domain/artifact.go` mit `Artifact`-Enum-String,
  `NewArtifact`-Konstruktor und `Artifact.String()`. Validierungs-Test:
  unbekannter Wert ⇒ Fehler mit dem Spec-Katalog in der Message.
- `internal/hexagon/port/driving/generate.go` mit `GenerateRequest`,
  `GenerateResponse`, `GenerateAction`-Enum (`Created`,
  `UpdatedBlock`, `NoOp`, `RepairedManual`), `GenerateUseCase`-
  Interface, Sentinels (`ErrArtifactUnknown`,
  `ErrProjectNotInitialized` reuse, `ErrGenerateManualConflict`,
  `ErrGenerateFileSystem`).
- `internal/hexagon/application/generate.go` mit
  `GenerateService`-Skeleton:
  - DI-Constructor `NewGenerateService(fs, yaml, logger)`.
  - `Generate(ctx, req)` macht (a) Project-State-Check
    (`u-boot.yaml`-Exists wie M5), (b) Dispatch über
    `req.Artifact.String()`-Switch → vier Handler, die in T1 **alle**
    `errors.New("generate <artifact>: handler not implemented")`
    returnen. Bewusst **kein** Slice-Marker („M7-T2..T5") in der
    Runtime-Message, weil solche Refs in Prod-Fehlern rotten. Der
    Build-Fail-Switch ist stattdessen ein unexportierter Marker
    `errStubHandler` in `generate.go`, auf den ein Paket-interner
    Test in T1 pinnt; sobald ein Tranchen-Slice (T2..T5) seinen
    Handler ersetzt, fällt der Pin auf weniger Stubs und in T5 schließlich
    auf null (Test wird in T5 entfernt — siehe T5-DoD).
- Tests in `_test`-Package (`generate_test.go`):
  - Project-not-initialized: kein `u-boot.yaml` ⇒
    `ErrProjectNotInitialized` (errors.Is).
  - Stub-Pfade: jeder der vier Artefaktwerte triggert seinen Handler-
    Stub und gibt einen Fehler zurück, der `errors.Is(err,
    errStubHandler)` erfüllt (Paket-interner Sentinel, **nicht**
    Teil der Driving-Port-API — damit ein versehentliches Mergen
    ohne T2–T5 laut auffällt, aber keine Slice-Refs in der Public
    Error-Surface erscheinen).
  - `BaseDir == ""` ⇒ non-nil error (kein Sentinel; analog
    `AddServiceService`).

**DoD T1:**
- [ ] `domain.Artifact` + `NewArtifact` 100 % Coverage.
- [ ] `GenerateUseCase`-Interface in `driving/generate.go` exportiert,
  Sentinels (`ErrArtifactUnknown`, `ErrGenerateManualConflict`,
  `ErrGenerateFileSystem`) deklariert. Driven-Sentinel-Scan für
  `driven.ErrFileSystem*` durchgeführt; Entscheidung „Wrap vs.
  Direct-Is" dokumentiert in `generate.go`-Top-Kommentar.
- [ ] `GenerateService.Generate` dispatcht korrekt; alle vier Handler
  geben einen Fehler zurück, der `errors.Is(err, errStubHandler)`
  erfüllt (paket-interner Sentinel, **nicht** Teil der
  Driving-Port-API).
- [ ] Keine CLI-Verkabelung (das ist T6); Use-Case ist erreichbar nur
  über direkte Test-Aufrufe.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T1 ✅ <commit-hash>`.

### T2 — `generate env-example`

Einfachster Generator: ein einzelnes Template (`env.example.tmpl` aus
M3 existiert bereits) mit einem managed-Block in `StyleHash`-Form.

**State-Machine:**

| Zustand | `.env.example` | managed-Block `init` | Aktion |
| ------- | -------------- | -------------------- | ------ |
| **absent** | fehlt | — | `actionWrite`: komplettes Template rendern, neue Datei. `GenerateAction=Created`. |
| **present-with-block** | vorhanden | vorhanden | `actionReplaceBlock`: nur Block neu rendern + splicen; Inhalt außerhalb des Blocks (Service-Add-on-Blöcke, User-Vars) bleibt byte-identisch. `GenerateAction=UpdatedBlock`. Wenn der gerenderte Block identisch zum existierenden ist ⇒ `NoOp` (Idempotenz-Pin). |
| **present-no-block** | vorhanden | fehlt | `ErrGenerateManualConflict` mit Repair-Hint. Code 10. Kein Write. |
| **block-malformed** | vorhanden | BEGIN ohne END / duplicate BEGIN | `managedblock.ErrBlockMalformed` → wrap als `ErrGenerateManualConflict` mit anderer Detail-Message. Kein Write. |

**Tests:**
- Vier State-Fixtures plus ein Add-on-Erhaltungs-Test: präparieren
  einer `.env.example` mit dem `init`-Block plus einem hinzugefügten
  `service.postgres`-Block (analog M5-T4c-Output); `generate
  env-example` re-rendered den init-Block; assertion: `service.postgres`-
  Block byte-identisch, init-Block content semantisch identisch
  (Render des aktuellen `env.example.tmpl`).
- `NoOp`-Pin: zweimaliges `generate env-example` hintereinander →
  zweiter Lauf returnt `NoOp`, `Changed=nil`, **und** der
  `FileSystem`-Fake registriert **null** `WriteFile`-Aufrufe für den
  zweiten Lauf (Counting-Fake, analog M5-T4c). Action-Field allein
  reicht als Idempotenz-Beweis nicht; ohne den Schreib-Zähler könnte
  ein Handler `NoOp` zurückgeben **und** trotzdem schreiben.

**DoD T2:**
- [ ] `generateEnvExample`-Handler implementiert; Stub aus T1 ersetzt.
- [ ] 5 State-Tests grün, NoOp-Pin grün.
- [ ] Stub-Pin-Test in `generate_test.go` auf **3** verbliebene
  Stubs (`readme`, `changelog`, `devcontainer`) reduziert — damit
  ein Auslassen einer Folgetranche laut auffällt, ohne dass der
  Pin lautlos stale wird.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T2 ✅ <commit-hash>`.

### T3 — `generate readme`

Strukturell wie T2, aber mit `StyleHTMLComment` und einem
Markdown-Template (`readme.md.tmpl` aus M3 existiert).

**Idempotenz-Vertrag:** Identische State-Machine zu T2; User-Content
nach dem `<!-- END U-BOOT MANAGED BLOCK: init -->`-Marker
(üblicherweise: User-eigene Setup-Anleitung, Screenshots,
Lizenz-Sektion) bleibt byte-identisch. Der `init`-Block aus dem
Template wird re-gerendered.

**Tests** (gleicher Schnitt wie T2, plus):
- User-Content-nach-Block-Test: Fixture mit init-Block + frei
  bearbeitetem Markdown danach (`## Custom section\n\nUser text.\n`);
  `generate readme` darf den User-Bereich nicht verändern.

**DoD T3:**
- [ ] `generateReadme`-Handler implementiert.
- [ ] State + User-Content-Tests grün.
- [ ] Stub-Pin-Test auf **2** verbliebene Stubs (`changelog`,
  `devcontainer`) reduziert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T3 ✅ <commit-hash>`.

### T4 — `generate changelog`

Komplexer als T2/T3 wegen der Doppelrolle des `init`-Blocks: der
M3-Template-Block enthält den initialen Header und einen
`## [Unreleased]`-Stub. LH-AK-007 verlangt: „vorhandene Inhalte
werden nicht zerstört, neuer Abschnitt wird korrekt ergänzt oder
vorbereitet" — was im Widerspruch dazu steht, den `init`-Block
einfach zu re-rendern (denn User-Einträge unter `### Added`
würden überschrieben).

**Entscheidung für M7:** Zweistufiges Verhalten, abhängig vom
Datei-Zustand:

| Zustand | `CHANGELOG.md` | `init`-Block | `## [Unreleased]`-Sektion | Aktion |
| ------- | -------------- | ------------ | -------------------------- | ------ |
| **absent** | fehlt | — | — | `actionWrite`: Template rendern, neue Datei. `Created`. |
| **present-with-block-no-edits** | vorhanden | vorhanden | unverändert vs. Template | `actionReplaceBlock`: Block re-rendern. `UpdatedBlock` (wenn diff != ∅) oder `NoOp`. |
| **present-with-block-edited** | vorhanden | vorhanden | User hat Einträge ergänzt | **No-op-Pfad mit Diagnose**: kein Re-Render des init-Blocks; stattdessen prüfen, ob die `## [Unreleased]`-Sektion existiert und nicht-leer ist; falls **alle** Pflicht-Subsektionen (`### Added`, `### Changed`, `### Fixed`) fehlen, einen neuen Stub-`## [Unreleased]`-Header **außerhalb** des managed-Blocks vor der ersten Versions-Sektion einfügen (`RepairedManual`). Detection-Heuristik: Hash des Block-Bodys aus der existierenden Datei != Hash des Block-Bodys aus `renderTemplate("changelog.md.tmpl", {Name: cfg.Name})` ⇒ user-edited. Wichtig: gegen das *gerenderte* Template mit dem aktuellen Projektnamen vergleichen, nicht gegen die rohe Template-Quelle — sonst kippt jedes Projekt sofort in den user-edited-Pfad und der `NoOp`-Pin am Ende dieser Tranche schlägt fehl. |
| **present-no-block** | vorhanden | fehlt | — | `ErrGenerateManualConflict`. |

Die Wahl konservativ-no-op bei user-edited Blocks ist die
Idempotenz-sichere Variante. Die alternative Strategie (managed-Block
nur für *Struktur*, User-Einträge wandern in eine Sektion **außerhalb**
des Blocks) würde eine Template-Migration für bestehende Projekte
erzwingen und ist deshalb V1.

**Bekannte Fragilität der Hash-Heuristik:** Sobald
`changelog.md.tmpl` jemals geändert wird (Header-Wording, Datum,
Sektions-Reihenfolge), kippt **jedes** existierende Projekt-Changelog
schlagartig in den „user-edited"-Pfad und bekommt keine Block-
Aktualisierung mehr — auch wenn der User die Datei nie angefasst
hat. M7 akzeptiert das bewusst: das `init`-Template ist nach M3
eingefroren und Template-Änderungen sind ein **Breaking-Migration-
Event**, das einen eigenen Folge-Slice mit `--migrate`-Pfad oder
versionierten Marker (`<!-- BEGIN U-BOOT MANAGED BLOCK: init v2 -->`)
braucht. Siehe „Out of Scope".

**Tests:**
- Vier State-Fixtures.
- `RepairedManual`-Pfad: Fixture mit user-edited init-Block plus
  Versions-Sektion `## [0.1.0] - 2026-01-01` aber **ohne**
  `## [Unreleased]`; `generate changelog` fügt einen Unreleased-Stub
  vor `[0.1.0]` ein.
- Idempotenz-Pin: doppelter Lauf auf einer user-edited Datei mit
  bereits vorhandenem Unreleased-Stub ⇒ zweiter Lauf `NoOp`.

**DoD T4:**
- [ ] `generateChangelog`-Handler implementiert.
- [ ] LH-AK-007-Pin: ein End-to-end-Test, der genau dem Spec-Wortlaut
  folgt (`u-boot init && u-boot generate changelog`) — Datei existiert,
  Vor-Inhalt nicht zerstört, Sektion korrekt ergänzt.
- [ ] Stub-Pin-Test auf **1** verbliebenen Stub (`devcontainer`)
  reduziert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T4 ✅ <commit-hash>`.

### T5 — `generate devcontainer`

Net-neuer Generator: zwei neue Templates und ein neuer Marker-Stil
(`StyleDoubleSlash` für JSONC). Wichtigste Komplexität: `forwardPorts`
muss aus dem aktiven Service-State abgeleitet werden, nicht
hartcodiert.

**Neue Templates:**
- `internal/hexagon/application/templates/devcontainer/devcontainer.json.tmpl`
  — JSONC-Template mit `// BEGIN U-BOOT MANAGED BLOCK: init`-Marker.
  Pflichtfelder:
  - `name`: `"{{.Name}}"` (Projektname aus `u-boot.yaml`).
  - `build`: Object mit `dockerfile: "./Dockerfile"` und
    `context: "."`. (Spec verlangt mindestens eines aus `build` oder
    `image`, LH-AK-005; `build` ist konsistent mit dem mit-
    generierten Dockerfile.)
  - `forwardPorts`: Array; T5 lässt das Template das Feld **leer**
    (`[]`) und füllt es im Render-Schritt aus dem Service-State.
  - `remoteUser`: `"vscode"` (LH-FA-DEV-004 nicht-root).
- `internal/hexagon/application/templates/devcontainer/Dockerfile.tmpl`
  — Multi-Stage-Dockerfile-Template mit `# BEGIN U-BOOT MANAGED BLOCK: init`-
  Marker. Basisimage z. B. `mcr.microsoft.com/devcontainers/base:debian`,
  `USER vscode`-Sektion am Ende.
- `internal/hexagon/application/templates.go`-Embed erweitern:
  `templates/devcontainer/*.tmpl` in `//go:embed` aufnehmen, damit
  `renderTemplate("devcontainer/...")` die neuen Templates findet.

**`forwardPorts`-Quelle (entscheidend für LH-FA-DEV-005):**

Drei Kandidaten:
1. `u-boot.yaml`-`services.<name>.ports` (existiert heute nicht;
   müsste neu eingeführt werden).
2. `compose.yaml`-managed-Blocks ⇒ `ports:`-Sub-Schlüssel parsen.
3. Service-Add-on-Katalog (heute hartcodiert: `postgres` ⇒ 5432).

**Entscheidung M7:** Variante **(2)** — `compose.yaml`-Lesen über den
existierenden `YAMLCodec` (Read-only, kein neuer Port-Patch). Pro
`services.<name>:`-Block, der einen managed-Marker (`service.<name>`)
trägt **und** in `u-boot.yaml` als `enabled: true` markiert ist,
werden die Container-Ports gesammelt (letztes Segment eines
Compose-Port-Mappings; `8080:80` ⇒ `80`,
`127.0.0.1:8080:80/tcp` ⇒ `80`). Das entspricht der bestehenden
Doctor-Logik für `devcontainer.forwardPorts.consistency`. Sortiert,
dedupliziert. Bei `[]` ⇒ `forwardPorts` fehlt im
generierten JSON (LH-FA-DEV-005: „darf fehlen"). Kandidat (1) wäre
ein zusätzlicher Spec-Tree (V1-Folgeslice), Kandidat (3) würde Add-on-
Katalog-Pflege erzwingen.

**State-Machine** (pro Datei `.devcontainer/devcontainer.json` und
`.devcontainer/Dockerfile` separat ermittelt — beide werden in einem
**atomaren Plan-and-Execute-Schritt** behandelt: erst Plan für
**beide** Dateien aufstellen, alle Vorbedingungen prüfen, **dann**
schreiben. Bricht die Validierung einer Datei ab, wird **keine** der
beiden geschrieben — sonst entstehen halbe Schreibzustände, die der
nächste Lauf erneut als Konflikt sieht):

| Zustand | Datei A (`devcontainer.json`) | Datei B (`Dockerfile`) | Aktion |
| ------- | ----------------------------- | ---------------------- | ------ |
| **absent**       | fehlt              | fehlt              | Beide schreiben. `Created`. |
| **partial-clean**| fehlt **oder** vorhanden+Block | umgekehrt          | Fehlende neu schreiben, vorhandene per Block-Replace. Aggregat-Action: `UpdatedBlock` wenn mindestens eine geupdated wurde, sonst `Created`. Häufigster Re-Run-Pfad. |
| **both-present-with-block** | vorhanden+Block | vorhanden+Block | Beide per Block-Replace. `UpdatedBlock` oder `NoOp` (NoOp nur wenn **beide** rendern-identisch). |
| **block-missing-in-any**    | vorhanden, **kein** Block (oder malformed) | beliebig | `ErrGenerateManualConflict` mit Hinweis auf **alle** betroffenen Dateien (kann eine oder beide sein). **Kein** Write, auch nicht für die intakte Datei. |
| **fs-error**     | Read/Write-Fehler  | beliebig           | `ErrGenerateFileSystem`-Wrap; Aufruf endet mit Exit 14. Kein Teil-Write. |

**Tests:**
- Mindestfelder-Pin (LH-FA-DEV-001/004/005, **nicht** LH-AK-005 —
  letzteres verlangt explizit `u-boot init --devcontainer` und
  gehört in den MVP-Closure-Slice): `u-boot init && u-boot add
  postgres && u-boot generate devcontainer` ⇒ beide Dateien
  existieren, JSON ist syntaktisch gültig (JSONC mit Kommentaren ist
  OK; gegen `stripJSONC`+ `encoding/json.Valid` prüfen), Name
  korrekt, `build` vorhanden, `forwardPorts: [5432]` enthält
  Postgres-Port.
- Forward-Ports-Detection: drei Fixture-Variationen:
  (a) keine Services aktiv ⇒ `forwardPorts` fehlt;
  (b) Postgres aktiv ⇒ `[5432]`;
  (c) Postgres aktiv + zweiter Service mit `ports: ["8080:80"]` ⇒
  `[80, 5432]` (sortiert; Container-Ports).
- Idempotenz: doppelter Lauf ohne Service-Änderung ⇒ `NoOp`.
- `remoteUser: vscode`-Pin (LH-FA-DEV-004).
- `doctor`-Integration (Bestätigungstest, kein neuer Code) — Test
  benannt `TestDoctor_AfterGenerateDevcontainer_PinsWarnPath_WhenEnabledFalse`,
  damit ein späterer Maintainer die Intention nicht versehentlich
  kippt: nach `generate devcontainer` muss `u-boot doctor` (M4-Stand)
  bei `devcontainer.enabled=false` die Dateien als `warn` (nicht
  `error`) melden. M7 setzt `devcontainer.enabled=true` in
  `u-boot.yaml` **nicht** automatisch — `generate devcontainer` ist
  ein Datei-Schreiber, kein Konfig-Mutator; ein V1-Folge-Slice
  könnte das per Flag ergänzen (analog `add postgres` ⇒
  `services.postgres.enabled=true`).
- **Anti-Drift-Pin gegen `doctor.collectActiveServicePorts`:** ein
  Test legt eine Fixture mit `u-boot.yaml` (postgres + dummy-Service
  mit `ports: ["8080:80"]`) und `compose.yaml` an, ruft den von T5
  genutzten Pfad **und** den Doctor-Helper auf derselben Fixture auf
  und vergleicht die zurückgegebene Port-Liste byte-identisch
  (`reflect.DeepEqual`). Damit ist explizit verboten, dass `generate
  devcontainer` und `devcontainer.forwardPorts.consistency`
  jemals auseinanderdriften.

**DoD T5:**
- [ ] Zwei neue Templates eingecheckt (`devcontainer.json.tmpl` +
  `Dockerfile.tmpl`).
- [ ] `//go:embed` in `templates.go` deckt `templates/devcontainer/*.tmpl`
  ab; Template-Integrity-Test listet die beiden Dateien.
- [ ] `generateDevcontainer`-Handler implementiert, inkl. Port-Detection
  aus `compose.yaml` via Aufruf der bestehenden package-internen
  Helper `activeServiceNames` + `collectActiveServicePorts`.
- [ ] Atomarer Plan-and-Execute: bei Block-Konflikt in einer der beiden
  Dateien wird **keine** geschrieben (eigener Test pinnt das, indem
  `FileSystem.WriteFile`-Counter auf 0 prüft).
- [ ] Anti-Drift-Pin gegen `doctor.collectActiveServicePorts` grün.
- [ ] Mindestfelder-Pin grün (End-to-end mit Postgres,
  LH-FA-DEV-001/004/005).
- [ ] Stub-Pin-Test aus T1 (`errStubHandler`) wird hier entfernt — alle
  vier Handler sind ab T5 implementiert.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T5 ✅ <commit-hash>`.

### T6 — CLI-Subcommand + Doku + Carveouts

- `internal/adapter/driving/cli/generate.go` mit
  `newGenerateCommand(a)` analog `newAddCommand`:
  - `Use: "generate <artifact>"`.
  - `Args: cobra.ExactArgs(1)`.
  - Mapping: positional arg → `domain.NewArtifact`. Fehler ⇒ Wrap
    in `ErrArtifactUnknown`, der über `ExitCode` auf Code **2**
    mappt (anders als bei `add`, wo unbekannter Service-Name auf
    Code 10 mappt; in `ExitCode`-Switch eine eigene `if errors.Is(...,
    ErrArtifactUnknown)`-Klausel).
  - Use-Case-Aufruf, dann `printGenerateSummary` mit drei Formen:
    - `Created` → `"Generated <artifact> (<paths>)."`
    - `UpdatedBlock` → `"Updated <artifact> managed block (<paths>)."`
    - `NoOp` → `"<artifact> already up to date; no changes."`
    - `RepairedManual` → `"Repaired <artifact> structure (<paths>)."`
- `cli.go` `newRootCommand` ergänzen: `cmd.AddCommand(newGenerateCommand(a))`.
- `App`-Struct in `cli.go:35` bekommt
  `generateUseCase driving.GenerateUseCase`-Field. **Breaking-Change
  in `cli.New`-Signatur** (`cli.go:111`): der Konstruktor erhält
  einen neuen positionalen Parameter `genUC
  driving.GenerateUseCase`. Betroffen sind exakt:
  - `cmd/uboot/main.go:106` (einziger Produktiv-Aufruf von
    `cli.New`; neu: `cli.New(version, initSvc, doctorSvc, addSvc,
    upSvc, downSvc, generateSvc, cli.WithLogLevel(...))`).
  - `internal/adapter/driving/cli/cli_test.go` — die fünf
    Helper-Funktionen `newApp` / `newAppWithDoctor` /
    `newAppWithAdd` / `newAppWithUp` / `newAppWithDown` (Zeilen 99,
    105, 110, 115, 120) sind die einzigen Test-Aufrufer von
    `cli.New` und tragen jeweils alle anderen Use-Cases als
    Fake-Defaults. Hier ergänzt T6 einen neuen Default-Fake
    `fakeGenerateUseCase` (analog zu den fünf bestehenden
    `fake*UseCase`-Typen ab `cli_test.go:22`) plus einen sechsten
    Helper `newAppWithGenerate(genUC, opts...)`.
  - `verbosity_test.go` und `statusview_test.go` rufen `cli.New`
    **nicht** direkt auf — sie gehen ausschließlich über `newApp` —
    und brauchen deshalb keine eigene Migration; sie kompilieren
    weiter, sobald `newApp` den neuen Default-Fake mitführt.
  - Es existiert **keine** `fakes_test.go` in diesem Paket; Fakes
    leben in `cli_test.go`. M7 legt keine separate `fakes_test.go`
    an (wäre eine Strukturänderung, die zu einer Codestil-
    Entscheidung gehört, nicht zu diesem Slice).
  Alle direkten `cli.New`-Aufrufstellen (`main.go` + die fünf
  `cli_test.go`-Helper) müssen in **einem** Commit mitgezogen werden,
  sonst bricht `go build ./...`. — Wenn dieses Mitziehen disruptiv wirkt,
  wäre die Alternative ein Functional-Option-Constructor
  (`cli.WithGenerateUseCase(genUC)`), aber das wäre eine
  Abweichung vom etablierten Muster der anderen fünf Use-Cases und
  ist deshalb hier **nicht** der gewählte Weg.
- `docs/user/quality.md` §2 Tests: kein Eintrag nötig (Generate-Tests
  laufen im Standard-`make gates`-Pfad, kein neuer Build-Tag).
- `docs/user/cli.md` (oder analog, falls existent — sonst README.md):
  `generate`-Subkommando dokumentieren.
- **Carveouts:** `carveouts.md` enthält Stand Pre-T1 keinen
  offenen Eintrag zu `LH-FA-GEN-*` oder „Devcontainer-Generator"
  (gegen `docs/plan/planning/in-progress/carveouts.md` verifiziert);
  dieser Schritt ist deshalb erwartungsgemäß ein No-op. Sollte vor
  T6-Merge dennoch ein passender Eintrag entstanden sein, wird er
  mit M7-DoD-Hash als gelöst markiert.

**Tests:**
- `cli_test.go` (Cobra-Smoke): `u-boot generate unknown-artifact` ⇒
  Exit 2, Stderr enthält Katalogliste.
- `u-boot generate changelog` ohne `u-boot.yaml` ⇒ Exit 10,
  `ErrProjectNotInitialized`.
- `u-boot generate changelog` ohne positional arg ⇒ Exit 2,
  Cobra-Usage-Message.

**DoD T6:**
- [ ] CLI-Subkommando verfügbar; `u-boot generate --help` listet die vier
  Artefakte explizit.
- [ ] Smoke-Tests im `cli_test.go` grün, inkl. ein Test
  `TestExitCode_GenerateFileSystemError_MapsTo14`, der den neu
  hinzugefügten Sentinel `ErrGenerateFileSystem` in die
  bestehende `isFilesystemError`-Liste (`cli.go:264-267`,
  bereits getestet via `TestExitCode_BaseMappings` für
  `ErrBackupSuffixExhausted`/`ErrBackupSourceMissing`) einreiht
  und sein Code-14-Mapping pinnt.
- [ ] `cli.New`-Aufrufstellen migriert: `cmd/uboot/main.go:106`
  plus die fünf Helper in `cli_test.go` (Zeilen 99, 105, 110, 115,
  120). Neuer `fakeGenerateUseCase` in `cli_test.go` ergänzt;
  `verbosity_test.go` und `statusview_test.go` brauchen keine
  Änderung, weil sie nur über `newApp` gehen (keine
  `cli.New`-Direktaufrufe). `go build ./...` und `go test ./...`
  grün.
- [ ] `u-boot generate readme` produziert ein README, dessen Markdown-
  Links im `docs-check` (Markdown-Link-Validator-Slice schon Done)
  sauber durchlaufen.
- [ ] Eintrag in [`roadmap.md`](../in-progress/roadmap.md) auf
  „Done" gesetzt mit Slice-Link.
- [ ] `make gates` grün.
- [ ] DoD-Line: `T6 ✅ <commit-hash>`.

## Akzeptanzkriterien (Slice-übergreifend)

### Struktur

- Vier neue Handler-Funktionen im `GenerateService`, jeweils mit
  eigener State-Machine-Tabelle in `generate.go`-Top-Kommentaren.
- Kein neuer Driven-Port: M7 nutzt `FileSystem`, `YAMLCodec` (read-
  only für T5-Port-Detection) und `Logger`. **Keine** Erweiterung
  von `YAMLCodec` nötig — der `Unmarshal`-Pfad reicht für die
  Compose-Ports-Auslese.
- Drei Marker-Stile aus `managedblock` (`StyleHash`,
  `StyleHTMLComment`, `StyleDoubleSlash`) decken die vier Datei-
  Mappings ab; alle bereits implementiert und getestet.

### Verhalten

- **LH-FA-GEN-001**: `u-boot generate <artifact>` mit den vier
  erlaubten Werten; alle anderen Werte ⇒ Exit-Code 2 mit explizit
  aufgelistetem Katalog in Stderr.
- **LH-FA-GEN-002 / LH-AK-007**: `generate changelog` schafft oder
  aktualisiert `CHANGELOG.md`, zerstört keine User-Einträge.
- **LH-FA-GEN-003**: `generate readme` schafft oder aktualisiert
  `README.md`, zerstört keine User-Sektionen außerhalb des
  `init`-Blocks.
- **LH-FA-GEN-004**: `generate env-example` schafft oder aktualisiert
  `.env.example`, zerstört keine Service-Add-on-Blöcke noch User-
  Variablen außerhalb des Blocks.
- **LH-FA-GEN-005 / LH-AK-006-analog**: doppelter Aufruf ist ein
  `NoOp` (Action-Field expliziter Beleg; CLI-Output „already up to
  date").
- **LH-FA-DEV-001**: `generate devcontainer` produziert beide
  Pflichtdateien; JSON ist syntaktisch gültig (JSONC-Stripper +
  `encoding/json.Valid`). **Hinweis:** LH-AK-005 ist hiermit
  **nicht** geschlossen — die Akzeptanz verlangt explizit
  `u-boot init --devcontainer`; dieser Flag-Pfad gehört in den
  MVP-Closure-Slice (siehe Out of Scope), nutzt aber die hier
  eingeführten Templates wieder.
- **LH-FA-DEV-005**: `forwardPorts` enthält die Container-Ports aller
  aktiven Services (Postgres ⇒ 5432); leer ⇒ Feld fehlt.

### Negative

- Kein `u-boot.yaml` ⇒ Exit 10 mit `ErrProjectNotInitialized` und
  Repair-Hint („Run `u-boot init` first").
- Datei existiert **ohne** managed-Block ⇒ Exit 10 mit
  `ErrGenerateManualConflict`; M7 schreibt nichts.
- Malformed managed-Block (BEGIN ohne END, duplicate BEGIN) ⇒ Exit 10,
  gleiche Sentinel, andere Detail-Message.
- Unerwarteter IO-/Permissions-Fehler beim Lesen oder Schreiben ⇒
  Exit 14 (`ErrGenerateFileSystem`-Wrap, in T6 erstmals in
  `cli.ExitCode` verdrahtet).

## Out of Scope (M7-spezifisch)

- **`--replace`-Flag** zum erzwungenen Überschreiben einer Datei
  ohne managed-Block. Runtime-Fehlermeldungen in M7 erwähnen den
  Flag-Namen nicht als Repair-Hint, weil der Flag selbst erst mit V1
  ([`slice-v1-template-format-entscheidung`](../done/slice-v1-template-format-entscheidung.md)
  oder eigener Folge-Slice). Begründung: MVP ist konservativ
  no-write, der Recovery-Pfad ist „Datei umbenennen, generate
  erneut laufen lassen".
- **`init --devcontainer`-Flag** (LH-AK-005) — gehört zu
  MVP-Closure; M7 liefert nur die Templates. **Risikohinweis für
  den Closure-Slice:** weil M7 für die zwei Devcontainer-Dateien
  denselben Block-Namen `init` verwendet wie alle M3-Templates
  (siehe Architektur-Punkt „Block-Name in allen generierten
  Dateien: `init`"), muss `init --devcontainer` zwingend einen
  `actionReplaceBlock`-Pfad für eine bereits durch `generate
  devcontainer` angelegte Datei anbieten — sonst kollidiert die
  Flag-Variante mit dem Generator-Output und muss über
  `--force/--backup` umgeleitet werden. Diese Anforderung wird im
  Closure-Slice explizit als T0-Carveout-Bedingung übernommen.
- **Auto-Setzen von `devcontainer.enabled=true`** beim `generate
  devcontainer` — analog wäre es zu `add postgres` ⇒
  `services.postgres.enabled=true`, aber Spec verlangt das nicht für
  `generate`. Ein V1-Folge-Slice könnte ein `--enable`-Flag
  ergänzen.
- **`generate changelog --release <version>`**-Workflow (User-
  initiierter Release-Schnitt) — V1-Folge-Slice, ggf. gekoppelt an
  [`slice-v1-release-pipeline`](../done/slice-v1-release-pipeline.md).
- **`generate dockerfile`** für ein Anwendungs-Dockerfile außerhalb
  von `.devcontainer/` (LH-FA-DOC-002 V1) — eigener V1-Slice.
- **`--json`-Output** für Generate-Summary — analog zur M4/M6-
  Entscheidung V1.
- **Template-Migration für `init`-Blocks** (versionierte Marker
  `init v2` oder `--migrate`-Flag) — sobald
  `changelog.md.tmpl`/`readme.md.tmpl`/`env.example.tmpl` jemals
  inhaltlich geändert werden, kippen alle existierenden Projekte in
  den konservativen No-op-Pfad (siehe T4-Heuristik). M7 friert die
  M3-Templates ein und überlässt die Migration einem V1-Folge-Slice.

## Bezug

- Auslösende Spec: `spec/lastenheft.md` §4.8 (`LH-FA-GEN-001..005`),
  §4.3 (`LH-FA-DEV-001/004/005`), §LH-AK-005, §LH-AK-007.
- Hängt von: M3 (Template-Embed + `actionReplaceBlock`), M5
  (`YAMLCodec`/`managedblock`-Konventionen und Doctor-Referenz für
  `devcontainer.forwardPorts.consistency`), M6 (Abgrenzung:
  Host-Port-Probing aus `upservice_portparse.go` ist nicht die
  `forwardPorts`-Quelle).
- Phase: M7 (MVP-Abschluss vor MVP-Closure).
- Roadmap: ersetzt `Open` in `roadmap.md` durch
  `Open (Slice-Plan vorhanden)` mit Link auf diese Datei; nach T6
  wird der Eintrag auf `Done` gesetzt.

## Review-Followup (`27de9c5`)

Nach dem Merge von T6 lief ein Code-Review über T1..T6 (Agent
`code-documentation:code-reviewer`); ergab **keine Blocker** und
9 Findings, alle in einem Followup-Commit adressiert:

### Should-fix

- **S1 — fenced-code-block-Heuristik (echter Datenverlust-Pfad).**
  `firstReleaseSectionOffset` und `hasChangelogUnreleased` matchen
  jetzt nur `## [...]`-Header außerhalb von Backtick-Fences. Ohne
  diesen Fix hätte ein User-Changelog mit einem dokumentierten
  Keep-a-Changelog-Beispiel den RepairedManual-Splice mitten in den
  Fence gerouted und das Markdown korrumpiert. Helper
  `isOffsetInsideFencedBlock` (Parity-Count über die Prefix-Bytes).
  Tilde-Fences und indented Code-Blöcke sind bewusst nicht erkannt.
- **S2 — CRLF-Normalisierung im `bytes.Equal`-Vergleich.** Neuer
  `normaliseLF`-Helper, verdrahtet in allen drei Vergleichsstellen
  (`generateChangelog`, `generateManagedFile`, `planDevcontainerFile`).
  Ein auf Windows mit CRLF gespeichertes File registriert jetzt als
  fresh statt fälschlich in den user-edited-Pfad zu kippen. Der
  Splice selbst nutzt weiter Originalbytes; User-Line-Endings
  außerhalb des Blocks bleiben.
- **S3 — Doku-Hinweis zu Projekt-Rename.** Top-Kommentar von
  `generateChangelog` dokumentiert: bei Rename in
  `u-boot.yaml.project.name` flippt der Block in den user-edited-
  Pfad und bleibt auf dem alten Namen stuck (T2/T3 schreiben den
  Namen transparent um, T4 lässt ihn bewusst stehen).
- **S4 — `Artifact.String()` Out-of-Range.** Default-Branch zeigt
  jetzt `Artifact(N)` statt `unknown`, sodass Debug-Logs und
  Fehler-Messages den tatsächlichen Int sehen. Test aktualisiert.

### Nice-to-have

- **N1 — `executeDevcontainerPlans` default case.** Schaltet auf
  Programmer-Error um, wenn ein künftiges `devcontainerFileAction`-
  Enum-Value ohne Switch-Update landet.
- **N2 — `compose.yaml`-Parse-Error-Klassifikation dokumentiert.**
  Top-Kommentar von `collectDevcontainerForwardPorts` erklärt, dass
  ein YAML-Parse-Fehler als `ErrGenerateFileSystem` (Code 14)
  durchgereicht wird, weil der Doctor-Helper Read- und Parse-
  Failure nicht unterscheidet. Ein künftiger `driven.ErrYAMLParse`-
  Sentinel könnte den Pfad auf Code 10 verschieben — bekannte
  Klassifikations-Lücke, kein Bug.
- **N3 — Anti-Drift-Test-Sanity-Guard.** Vor der `reflect.DeepEqual`
  steht jetzt eine `len(doctorPorts) > 0`-Pflicht, damit ein
  künftiger doppelter Regression (Generator + Doctor liefern beide
  nil) nicht silent durchgeht.
- **N4 — `classifyExistingBlock`-Helper extrahiert.** Drei
  duplizierte Switch-Blöcke (~13 Zeilen je) zusammengezogen; alle
  drei Callsites routen jetzt `ErrBlockNotFound`/`ErrBlockMalformed`
  durch dieselbe `ErrGenerateManualConflict`-Message (format-
  agnostic, referenziert LH-SA-FILE-002).
- **N5 — `printGenerateSummary` default-Branch.** Zeigt jetzt
  `resp.Action` und `resp.Changed` statt silent auf
  `"Generated <name>"` zu truncaten.

### Neue Tests

`TestGenerateChangelog_UserEditedBlock_FencedReleaseOnly_NoOp` /
`_FencedReleaseBeforeReal_SpliceAtReal` /
`_FencedUnreleased_DoesNotCount` /
`TestGenerateChangelog_CRLFFreshBlock_NoOp` /
`TestGenerateEnvExample_CRLFFreshBlock_NoOp`.

Coverage steigt minimal von 90.10 % auf 90.20 %.
