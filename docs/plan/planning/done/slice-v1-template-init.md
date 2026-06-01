# Slice V1: `u-boot init --template <name>`

## Auslöser

Zweiter der drei [ADR-0009](../../adr/0009-template-format-yaml-files.md)
§Folgepunkte-Slices. Der erste
([`slice-v1-template-list`](slice-v1-template-list.md)) hat
Katalog + Driven-Port + `basic`-Bootstrap-Metadaten geliefert; dieser
Slice baut darauf den Render-Pfad und verdrahtet die externe Vorlage
in `u-boot init`.

Spec-IDs: `LH-FA-TPL-001` (Projektvorlagen, V1 — Beispiele
`u-boot init --template basic|micronaut|sveltekit|micronaut-sveltekit`),
plus die volle Surface von `LH-FA-TPL-002` (Template-Metadaten — Listing
hatte nur Name/Description/Version; hier kommen GeneratedFiles und
ggf. später Variables zum Tragen).

## Aufhebungsbedingung

`u-boot init --template basic` erzeugt ein Projekt, das **byte-
identisch** zum Default-`u-boot init` ist. Damit ist der Render-Pfad
nachweisbar korrekt — die default-Init-Templates aus
`internal/hexagon/application/templates/` und die externen
`externaltemplates/templates/basic/`-Files produzieren denselben
Output für identische Eingaben. Weitere Built-in-Templates (z. B.
`micronaut`) kommen in eigenen Slices, dann mit Variable-Resolution.

## Akzeptanzkriterien

- ✅ `u-boot init <name> --template basic` produziert dieselben
  Dateien wie `u-boot init <name>` ohne Flag (per Bytewise-Vergleich
  in einem E2E-Test gepinnt; `diff -r` der zwei Outputs ist leer).
- ✅ Unbekannter Template-Name (`--template nonexistent`) failed mit
  `ErrTemplateNotFound` und LH-FA-CLI-006 Exit-Code 10 (fachlich,
  Nutzer-Aktion erforderlich).
- ✅ Pfad-Sicherheit: `domain.TemplatePath` rejected absolute Pfade,
  `..`-Segmente (auch wenn `path.Clean` sie wegnormalisieren würde),
  Windows-Drive-Letters und Leerstring. 14 Test-Cases mit 100%-
  Coverage; Service-Wrap auf `driving.ErrInvalidTemplatePath` über
  Multi-%w.
- ✅ Render-Engine: `*.tmpl`-Dateien werden via `text/template` gegen
  `templateData{Name}` (analog M3-T2) gerendert; nicht-`.tmpl`-Dateien
  (z. B. `.gitignore`) werden 1:1 kopiert. Engine-parallel zur M3-T2-
  `renderTemplate`-Helper.
- ✅ Hexagonale Verdrahtung: neuer Driven-Port `TemplateFiles.Open()`
  (separat von `TemplateCatalog` — saubere SRP, bestehender
  `fakeCatalog` unbehelligt); Application-Service `TemplateInitService`
  orchestriert Lookup + Walk + Render + Write; CLI-Flag
  `--template <name>` auf der bestehenden `init`-Subkommando;
  InitProjectService delegiert via `WithTemplateInit`-Option.
- ✅ `domain.TemplatePath`-Validator eingeführt (analog M8
  `domain.ConfigPath`), ADR-0009 §Entscheidung verspricht das.
- ✅ Exit-Codes per LH-FA-CLI-006: 0 Erfolg; 10 Template-Not-Found
  oder Path-Eskalation (`isTemplateInitValidationError`-Helper); 14
  Render-/IO-Fehler (`isFilesystemError`); 2 Mutex-Verletzung
  (`ErrTemplateConflictsWithFlag` in `isUsageError`).

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `9e81b02` | `domain.TemplatePath` mit `NewTemplatePath`-Konstruktor analog `domain.ConfigPath`-Pattern. Reject-Liste: empty, absolute (Unix + Windows-Backslash), Windows-Drive-Letter, jede `..`-Sequenz im rohen Input (vor `path.Clean`). 14 Test-Cases (6 accept + 8 reject) + Round-Trip-Pin, 100% Funktions-Coverage. Driven-Port `driven.TemplateFiles.Open(ctx, name) (iofs.FS, error)` als SEPARATER Port (statt Erweiterung von `TemplateCatalog`) — bestehender `fakeCatalog` bleibt unbehelligt, Single-Responsibility. `driven.ErrTemplateNotFound`-Sentinel. Adapter `externaltemplates.Catalog.Open()` via `iofs.ReadDir` (Existenz-Check) + `iofs.Sub`, `ctx.Err()`-Entry-Check analog T1-`List`. 5 neue Adapter-Tests. |
| T2 | `65a1ce8` | Driving-Port `port/driving.TemplateInitUseCase` + Request/Response (`{BaseDir, ProjectName, TemplateName}` → `{Created []string}`). Drei Sentinels: `ErrTemplateNotFound` (10), `ErrInvalidTemplatePath` (10), `ErrTemplateRender` (14). Application-Service `TemplateInitService` mit Walk-Render-Skip-Loop: `.tmpl` via `text/template` rendert gegen `templateData{Name}`, sonst byte-identische Copy, `template.yaml` wird übersprungen, Parent-Dirs via `MkdirAll` on-the-fly. `renderOne`-Helper validiert jeden Pfad durch `domain.NewTemplatePath`. 7 Application-Unit-Tests (Happy-Path, Nested-Dirs, UnknownTemplate-Multi-%w-Chain, RenderFailure ohne Partial-Writes, EmptyBaseDir, NilLogger). InvalidTemplatePath-Boundary bewusst nicht als Integration-Test (Domain-Tests covern den Reject; fstest.MapFS rejected `..`-Pfade selbst via `fs.ValidPath`, Custom-FS wäre disproportional). |
| T3 | `ed6d9a0` | Bootstrap-Content für `externaltemplates/templates/basic/`: sechs Source-Files. Fünf sind byte-identische Kopien aus `internal/hexagon/application/templates/` (`compose.yaml.tmpl`, `README.md.tmpl`, `CHANGELOG.md.tmpl`, `.env.example.tmpl`, `.gitignore.tmpl`); `u-boot.yaml.tmpl` ist neu und mirror auf den `yaml.v3.Marshal`-Output (4-Space-Indent: `schemaVersion: 1` / `project:` / `    name: {{.Name}}`). embed-Pattern auf `all:templates/*` umgestellt — Go-`embed`-Default schließt führende-Punkt-Dateien aus, aber `.env.example.tmpl` und `.gitignore.tmpl` brauchen sie. Byte-Identity-Pin: `TestTemplateInitService_BasicByteIdenticalToDefaultInit` verifiziert alle sechs Outputs einzeln gegen captured Strings aus `docker run u-boot init demo --no-git`. |
| T4 | `daaaa9a` | CLI `--template <name>`-Flag (StringVar mit Help-Text + LH-FA-TPL-001-Referenz). `InitProjectRequest.Template`-Feld. `ErrTemplateConflictsWithFlag`-Sentinel (Exit 2 via `isUsageError`). `InitProjectService` bekommt `WithTemplateInit(uc)`-Functional-Option — additive `opts ...InitProjectOption` an `NewInitProjectService` (non-breaking für 7 Test-Callsites). Init() hat einen frühen Branch: wenn `req.Template != ""` → `initFromTemplate()`, der Soft-Existing-Detection / Project-Structure-Dirs / git init bewahrt und nur File-Rendering an `TemplateInitUseCase` delegiert. Mutex-Reject für `--template` + `--devcontainer`/`--force`/`--backup` (v1 fresh-init only). `isFilesystemError` um `ErrTemplateRender` ergänzt; neuer `isTemplateInitValidationError`-Helper (gocyclo-Carve-Out) für NotFound + InvalidPath. Wiring in `cmd/uboot/main.go`: TemplateInitService vor InitProjectService konstruiert, via Option durchgereicht. 6 Integration-Tests inkl. **E2E-Byte-Identity-Pin** (gegen die echte `externaltemplates.New()`-Adapter): default-Pfad und Template-Pfad produzieren bytewise identische Outputs für alle 6 Dateien. Smoke-Test gegen das gebaute Image: `diff -r` zwischen den Ausgaben ist leer. |
| T5 | dieser Commit | Slice-Plan nach `done/`; README.{md,de.md} `init`-Bullet erwähnt `--template <name>`; `CHANGELOG.md ## [Unreleased]` Added-Eintrag; `roadmap.md` §Nächste Schritte 3 mit T1-T4-Hashes und Markierung der zweiten ADR-0009-Folge als ✅; ADR-0009 §Folgepunkte template-init ✅-Häkchen + Verweis auf den done-Slice. `make docs-check` grün. |

## Out of Scope

- **Variable-Resolution + Prompt-Pfad**: das `basic`-Template hat
  `variables: []`. Erst wenn ein Built-in (z. B. `micronaut`) tatsächlich
  Variablen einführt, lohnt sich der Prompt/`--var key=value`-Pfad.
  Eigener Folge-Slice (`slice-v1-template-init-variables` oder
  zusammen mit dem ersten variable-bedürftigen Template-Slice).
- **`--template ./pfad`-Filesystem-Auflösung**: gehört zum
  Later-Slice `slice-later-local-templates` (ADR-0009 §Folgepunkte 3).
- **Weitere Built-in-Templates** über `basic` hinaus (`micronaut`,
  `sveltekit`, …): konkreter Bedarf-getriebene eigene Slices.
- **`--template` mit `--devcontainer`/`--force`/`--backup`**: T4
  rejected die Kombination mit `ErrTemplateConflictsWithFlag`.
  Re-Init-Pfad (managed-block-Semantics auf Template-Files)
  und Devcontainer-Layer beim Template-Render sind eigene
  Folge-Slices.
- **Refactor: Default-Init auf externe Templates umstellen**: konnte
  man theoretisch — `internal/.../templates/` würde dann nach
  `externaltemplates/templates/basic/` wandern und InitProjectService
  würde standardmäßig `--template basic` laufen. Defer, weil
  Byte-Identity-Pin in T3+T4 die Voraussetzung erst zementiert; ein
  Folge-Refactor-Slice kann den Default ohne Rendering-Risk auf
  externe Templates schwenken.

## Bezug

- ADR: [ADR-0009](../../adr/0009-template-format-yaml-files.md)
  §Entscheidung (`text/template`-Engine, Pfad-Eskalation-Validator,
  `embed.FS`-Layout) + §Folgepunkte (`slice-v1-template-init` mit
  T5 ✅).
- Voraussetzungs-Slice:
  [`slice-v1-template-list`](slice-v1-template-list.md)
  liefert TemplateCatalog-Port + `basic`-Metadaten.
- Spec: `LH-FA-TPL-001` (V1) — komplett geliefert für `basic`;
  `LH-FA-TPL-002` (V1) — Metadaten-Surface komplett, Variable-
  Resolution defer-pflichtig (Out-of-Scope).
- Architektur: hexagonale Schichten unverletzt; `domain.TemplatePath`
  analog M8 `domain.ConfigPath`-Pattern; `TemplateFiles`-Port
  parallel zu `TemplateCatalog` (zwei Rollen, eine Adapter-Instanz).
- Phase: V1 (nach `slice-v1-template-list`).
