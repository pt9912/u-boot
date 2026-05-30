# Slice M7: `u-boot generate <artifact>`-Flow

> **Status:** Open
> **DoD:** T1 ⬜ / T2 ⬜ / T3 ⬜ / T4 ⬜ / T5 ⬜ / T6 ⬜

## Auslöser

Nach M3 (`u-boot init`), M4 (`u-boot doctor`), M5 (`u-boot add
postgres`) und M6 (`u-boot up`/`down`) fehlt das letzte MVP-Subkommando
aus §4.8: **`u-boot generate <artifact>`**. Erst damit schließt sich
`LH-AK-007` (Changelog-Generator) und der MVP-Akzeptanz-Pfad für
artifaktive Wiederherstellung / -Aktualisierung.

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
  Tree + (optional) `services.<name>.ports`-Override; MVP liest aus
  `compose.yaml`-managed-Blocks (siehe T5).
- **`LH-FA-CLI-006`** Exit-Code-Mapping:
  - `2` — CLI-Validierung (unbekanntes Artefakt, fehlende
    positional args).
  - `10` — fachlicher Validierungsfehler
    (`ErrProjectNotInitialized` weil kein `u-boot.yaml`).
  - `13` — Datei-/IO-Fehler beim Schreiben.
- **`LH-SA-FILE-002`** Managed-Block-Konvention pro Datei-Format:
  - `.env.example` ⇒ `StyleHash` (`# BEGIN ...`).
  - `README.md`, `CHANGELOG.md` ⇒ `StyleHTMLComment` (`<!-- BEGIN ... -->`).
  - `.devcontainer/devcontainer.json` ⇒ `StyleDoubleSlash` (`// BEGIN ...`).
  - `.devcontainer/Dockerfile` ⇒ `StyleHash`.
  Alle vier Stile sind im bestehenden `managedblock`-Paket implementiert
  (siehe `internal/hexagon/application/managedblock/managedblock_test.go`).

Out of Scope (V1+):

- **`LH-FA-TPL-001..004`** Template-System mit benutzerdefinierten
  Projektvorlagen — eigener V1-Slice
  ([`slice-v1-template-format-entscheidung`](slice-v1-template-format-entscheidung.md)).
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
- [`slice-m6-up-down`](../done/slice-m6-up-down.md) — Port-Probe-
  Logik in `upservice_portparse.go`, die T5 (Devcontainer-
  `forwardPorts`) als Datenquelle wiederverwenden kann (Compose-Block
  → Port-Liste).

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
    managed-Block und damit keinen Reset-Anker; Code 10 mit
    Repair-Hint („Run `u-boot generate <artifact> --replace` oder
    füge `# BEGIN U-BOOT MANAGED BLOCK: init`-Marker manuell ein"). MVP
    schreibt nichts, V1-Folge-Slice könnte `--replace` ergänzen — der
    Flag-Name steht in der Fehlermeldung, wird in M7 aber nicht
    implementiert; siehe „Out of Scope".

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
  Interface, drei Sentinels (`ErrArtifactUnknown`,
  `ErrProjectNotInitialized` reuse, `ErrGenerateManualConflict`).
- `internal/hexagon/application/generate.go` mit
  `GenerateService`-Skeleton:
  - DI-Constructor `NewGenerateService(fs, yaml, logger)`.
  - `Generate(ctx, req)` macht (a) Project-State-Check
    (`u-boot.yaml`-Exists wie M5), (b) Dispatch über
    `req.Artifact.String()`-Switch → vier Handler, die in T1 **alle**
    `errors.New("generate <artifact>: not yet implemented (M7-T2..T5)")`
    returnen.
- Tests in `_test`-Package (`generate_test.go`):
  - Project-not-initialized: kein `u-boot.yaml` ⇒
    `ErrProjectNotInitialized` (errors.Is).
  - Stub-Pfade: jeder der vier Artefaktwerte triggert seinen Handler-
    Stub und gibt den „not yet implemented"-Fehler zurück (kein
    Sentinel, damit ein versehentliches Mergen ohne T2–T5
    laut auffällt).
  - `BaseDir == ""` ⇒ non-nil error (kein Sentinel; analog
    `AddServiceService`).

**DoD T1:**
- `domain.Artifact` + `NewArtifact` 100 % Coverage.
- `GenerateUseCase`-Interface in `driving/generate.go` exportiert.
- `GenerateService.Generate` dispatcht korrekt; Tests grün.
- Keine CLI-Verkabelung (das ist T6); Use-Case ist erreichbar nur
  über direkte Test-Aufrufe.
- `make gates` grün.
- DoD-Line: `T1 ✅ <commit-hash>`.

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
  zweiter Lauf returnt `NoOp`, `Changed=nil`.

**DoD T2:**
- `generateEnvExample`-Handler implementiert; Stub aus T1 ersetzt.
- 5 State-Tests grün, NoOp-Pin grün.
- `make gates` grün.
- DoD-Line: `T2 ✅ <commit-hash>`.

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
- `generateReadme`-Handler implementiert.
- State + User-Content-Tests grün.
- `make gates` grün.
- DoD-Line: `T3 ✅ <commit-hash>`.

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
| **present-with-block-edited** | vorhanden | vorhanden | User hat Einträge ergänzt | **No-op-Pfad mit Diagnose**: kein Re-Render des init-Blocks; stattdessen prüfen, ob die `## [Unreleased]`-Sektion existiert und nicht-leer ist; falls **alle** Pflicht-Subsektionen (`### Added`, `### Changed`, `### Fixed`) fehlen, einen neuen Stub-`## [Unreleased]`-Header **außerhalb** des managed-Blocks vor der ersten Versions-Sektion einfügen (`RepairedManual`). Detection-Heuristik: Block-Body-Hash != Template-Body-Hash ⇒ user-edited; konservativ. |
| **present-no-block** | vorhanden | fehlt | — | `ErrGenerateManualConflict`. |

Die Wahl konservativ-no-op bei user-edited Blocks ist die
Idempotenz-sichere Variante. Die alternative Strategie (managed-Block
nur für *Struktur*, User-Einträge wandern in eine Sektion **außerhalb**
des Blocks) würde eine Template-Migration für bestehende Projekte
erzwingen und ist deshalb V1.

**Tests:**
- Vier State-Fixtures.
- `RepairedManual`-Pfad: Fixture mit user-edited init-Block plus
  Versions-Sektion `## [0.1.0] - 2026-01-01` aber **ohne**
  `## [Unreleased]`; `generate changelog` fügt einen Unreleased-Stub
  vor `[0.1.0]` ein.
- Idempotenz-Pin: doppelter Lauf auf einer user-edited Datei mit
  bereits vorhandenem Unreleased-Stub ⇒ zweiter Lauf `NoOp`.

**DoD T4:**
- `generateChangelog`-Handler implementiert.
- LH-AK-007-Pin: ein End-to-end-Test, der genau dem Spec-Wortlaut
  folgt (`u-boot init && u-boot generate changelog`) — Datei existiert,
  Vor-Inhalt nicht zerstört, Sektion korrekt ergänzt.
- `make gates` grün.
- DoD-Line: `T4 ✅ <commit-hash>`.

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
werden die Host-Ports (`HOST:CONTAINER`-Split, Host-Seite vor `:`)
gesammelt. Sortiert, dedupliziert. Bei `[]` ⇒ `forwardPorts` fehlt im
generierten JSON (LH-FA-DEV-005: „darf fehlen"). Kandidat (1) wäre
ein zusätzlicher Spec-Tree (V1-Folgeslice), Kandidat (3) würde Add-on-
Katalog-Pflege erzwingen.

**State-Machine** (pro Datei `.devcontainer/devcontainer.json` und
`.devcontainer/Dockerfile` separat — beide werden in einem Aufruf
gemeinsam erzeugt):

| Zustand | Beide Dateien | Aktion |
| ------- | ------------- | ------ |
| **absent** | beide fehlen | Beide neu schreiben. `Created`. |
| **partial** | eine fehlt | Fehlende neu schreiben, vorhandene per Block-Replace. `Created` für die fehlende, `UpdatedBlock` für die vorhandene; `Action=UpdatedBlock` als Aggregat (häufigster Pfad bei Re-Run nach Service-Add). |
| **both-present-with-block** | beide haben init-Block | Beide per Block-Replace. `UpdatedBlock` oder `NoOp`. |
| **block-missing-in-one** | eine Datei ohne init-Block | `ErrGenerateManualConflict` mit Hinweis auf welche Datei. |

**Tests:**
- LH-AK-005-Pin: `u-boot init && u-boot add postgres && u-boot generate
  devcontainer` ⇒ beide Dateien existieren, JSON ist syntaktisch
  gültig (JSONC mit Kommentaren ist OK; gegen `stripJSONC`+
  `encoding/json.Valid` prüfen), Name korrekt, `build` vorhanden,
  `forwardPorts: [5432]` enthält Postgres-Port.
- Forward-Ports-Detection: drei Fixture-Variationen:
  (a) keine Services aktiv ⇒ `forwardPorts` fehlt;
  (b) Postgres aktiv ⇒ `[5432]`;
  (c) Postgres aktiv + zweiter Service mit `ports: ["8080:80"]` ⇒
  `[5432, 8080]` (sortiert).
- Idempotenz: doppelter Lauf ohne Service-Änderung ⇒ `NoOp`.
- `remoteUser: vscode`-Pin (LH-FA-DEV-004).
- `doctor`-Integration (Bestätigungstest, kein neuer Code): nach
  `generate devcontainer` muss `u-boot doctor` (M4-Stand,
  `devcontainer.enabled=true` falls gesetzt) keinen Error für die
  Devcontainer-Datei-Existenz mehr liefern. M7 setzt
  `devcontainer.enabled=true` in `u-boot.yaml` **nicht** automatisch
  — `generate devcontainer` ist ein Datei-Schreiber, kein Konfig-
  Mutator; ein V1-Folge-Slice könnte das per Flag ergänzen
  (analog `add postgres` ⇒ `services.postgres.enabled=true`). Doctor
  prüft mit `devcontainer.enabled=false` die Dateien als `warn` —
  Test pinnt diesen warn-Pfad statt eines erzwungenen `error`-Pfads.

**DoD T5:**
- Zwei neue Templates eingecheckt (`devcontainer.json.tmpl` +
  `Dockerfile.tmpl`).
- `generateDevcontainer`-Handler implementiert, inkl. Port-Detection
  aus `compose.yaml`.
- LH-AK-005-Pin grün (End-to-end mit Postgres).
- `make gates` grün.
- DoD-Line: `T5 ✅ <commit-hash>`.

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
- `App`-Struct in `cli.go` bekommt `generateUseCase
  driving.GenerateUseCase`-Field; `main.go`-Wireup ergänzen.
- `docs/user/quality.md` §2 Tests: kein Eintrag nötig (Generate-Tests
  laufen im Standard-`make gates`-Pfad, kein neuer Build-Tag).
- `docs/user/cli.md` (oder analog, falls existent — sonst README.md):
  `generate`-Subkommando dokumentieren.
- **Carveouts:** Prüfen, ob in `carveouts.md` ein offener Eintrag zu
  `LH-FA-GEN-*` oder „Devcontainer-Generator" existiert; falls ja,
  mit M7-DoD-Hash als gelöst markieren.

**Tests:**
- `cli_test.go` (Cobra-Smoke): `u-boot generate unknown-artifact` ⇒
  Exit 2, Stderr enthält Katalogliste.
- `u-boot generate changelog` ohne `u-boot.yaml` ⇒ Exit 10,
  `ErrProjectNotInitialized`.
- `u-boot generate changelog` ohne positional arg ⇒ Exit 2,
  Cobra-Usage-Message.

**DoD T6:**
- CLI-Subkommando verfügbar; `u-boot generate --help` listet die vier
  Artefakte explizit.
- Smoke-Tests im `cli_test.go` grün.
- `u-boot generate readme` produziert ein README, das `markdownlint`
  (Markdown-Link-Validator-Slice schon Done) sauber durchläuft.
- Eintrag in [`roadmap.md`](../in-progress/roadmap.md) auf
  „Done" gesetzt mit Slice-Link.
- `make gates` grün.
- DoD-Line: `T6 ✅ <commit-hash>`.

## Akzeptanzkriterien (Slice-übergreifend)

### Struktur

- Vier neue Handler-Funktionen im `GenerateService`, jeweils mit
  eigener State-Machine-Tabelle in `generate.go`-Top-Kommentaren.
- Kein neuer Driven-Port: M7 nutzt `FileSystem`, `YAMLCodec` (read-
  only für T5-Port-Detection) und `Logger`. **Keine** Erweiterung
  von `YAMLCodec` nötig — der `Unmarshal`-Pfad reicht für die
  Compose-Ports-Auslese.
- Vier Marker-Stile aus `managedblock` werden genutzt; alle bereits
  implementiert und getestet.

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
- **LH-FA-DEV-001 / LH-AK-005**: `generate devcontainer` produziert
  beide Pflichtdateien; JSON ist syntaktisch gültig (JSONC-Stripper
  + `encoding/json.Valid`).
- **LH-FA-DEV-005**: `forwardPorts` enthält die Host-Ports aller
  aktiven Services (Postgres ⇒ 5432); leer ⇒ Feld fehlt.

### Negative

- Kein `u-boot.yaml` ⇒ Exit 10 mit `ErrProjectNotInitialized` und
  Repair-Hint („Run `u-boot init` first").
- Datei existiert **ohne** managed-Block ⇒ Exit 10 mit
  `ErrGenerateManualConflict`; M7 schreibt nichts.
- Malformed managed-Block (BEGIN ohne END, duplicate BEGIN) ⇒ Exit 10,
  gleiche Sentinel, andere Detail-Message.

## Out of Scope (M7-spezifisch)

- **`--replace`-Flag** zum erzwungenen Überschreiben einer Datei
  ohne managed-Block. Die Fehlermeldung erwähnt den Flag-Namen als
  Repair-Hint, aber der Flag selbst kommt erst mit V1
  ([`slice-v1-template-format-entscheidung`](slice-v1-template-format-entscheidung.md)
  oder eigener Folge-Slice). Begründung: MVP ist konservativ
  no-write, der Recovery-Pfad ist „Datei umbenennen, generate
  erneut laufen lassen".
- **`init --devcontainer`-Flag** (LH-AK-005) — gehört zu
  MVP-Closure; M7 liefert nur die Templates.
- **Auto-Setzen von `devcontainer.enabled=true`** beim `generate
  devcontainer` — analog wäre es zu `add postgres` ⇒
  `services.postgres.enabled=true`, aber Spec verlangt das nicht für
  `generate`. Ein V1-Folge-Slice könnte ein `--enable`-Flag
  ergänzen.
- **`generate changelog --release <version>`**-Workflow (User-
  initiierter Release-Schnitt) — V1-Folge-Slice, ggf. gekoppelt an
  [`slice-v1-release-pipeline`](slice-v1-release-pipeline.md).
- **`generate dockerfile`** für ein Anwendungs-Dockerfile außerhalb
  von `.devcontainer/` (LH-FA-DOC-002 V1) — eigener V1-Slice.
- **`--json`-Output** für Generate-Summary — analog zur M4/M6-
  Entscheidung V1.

## Bezug

- Auslösende Spec: `spec/lastenheft.md` §4.8 (`LH-FA-GEN-001..005`),
  §4.3 (`LH-FA-DEV-001/004/005`), §LH-AK-005, §LH-AK-007.
- Hängt von: M3 (Template-Embed + `actionReplaceBlock`), M5
  (`PatchScalar` reuse für u-boot.yaml-Reads bei T5), M6 (Port-
  Detection-Heuristik aus `upservice_portparse.go`).
- Phase: M7 (MVP-Abschluss vor MVP-Closure).
- Roadmap: ersetzt `Open` in `roadmap.md` durch
  `Open (Slice-Plan vorhanden)` mit Link auf diese Datei; nach T6
  wird der Eintrag auf `Done` gesetzt.
