# Slice V1: `u-boot init --template <name>`

## Auslöser

Zweiter der drei [ADR-0009](../../adr/0009-template-format-yaml-files.md)
§Folgepunkte-Slices. Der erste
([`slice-v1-template-list`](../done/slice-v1-template-list.md)) hat
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

- `u-boot init <name> --template basic` produziert dieselben Dateien
  wie `u-boot init <name>` ohne Flag (per Bytewise-Vergleich in einem
  E2E-Test gepinnt).
- Unbekannter Template-Name (`--template nonexistent`) failed mit
  `ErrTemplateNotFound` und LH-FA-CLI-006 Exit-Code 10 (fachlich,
  Nutzer-Aktion erforderlich).
- Pfad-Sicherheit: ein Template, dessen `template.yaml.generatedFiles`
  einen absoluten Pfad (`/etc/passwd`) oder `..`-Komponente enthält,
  failed beim Render mit `domain.ErrInvalidTemplatePath`. Pin via
  Fixture-Template im Adapter-Test.
- Render-Engine: `*.tmpl`-Dateien werden via `text/template` gegen
  `templateData{Name, ForwardPorts}` (analog M3-T2) gerendert; nicht-
  `.tmpl`-Dateien (z. B. `.gitignore`) werden 1:1 kopiert (ADR-0009
  §Entscheidung).
- Hexagonale Verdrahtung: neuer Driven-Port `TemplateRenderer`
  (oder `TemplateCatalog.OpenTemplate(name)`) liefert das per-
  Template-File-Tree als `iofs.FS`; Application-Service
  `TemplateInitService` orchestriert Lookup + Render-Loop; CLI flag
  `--template <name>` auf der bestehenden `init`-Subkommando.
- `domain.TemplatePath`-Validator ist eingeführt (analog M8
  `domain.ConfigPath`), ADR-0009 §Entscheidung verspricht das.
- Exit-Codes per LH-FA-CLI-006: 0 Erfolg; 10 Template-Not-Found
  (`ErrTemplateNotFound`) oder Path-Eskalation
  (`ErrInvalidTemplatePath`); 14 Render-/IO-Fehler
  (`ErrTemplateRender`).

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | Domain `domain.TemplatePath` mit Konstruktor (Path-Eskalation rejected: absolute Pfade + `..`-Komponenten + Leerstring). Driven-Port-Erweiterung: `TemplateCatalog.OpenTemplate(ctx, name) (iofs.FS, error)` (oder neuer Port `TemplateRenderer`, T1-Decision). Adapter liefert `fs.Sub(templatesFS, "templates/"+name)`. Sentinel `driven.ErrTemplateNotFound` (falls Subdir nicht existiert). Adapter-Tests (Open für `basic` → fs.FS mit `template.yaml` zugänglich; unbekannter Name → ErrTemplateNotFound; Path-Validator-Tests). |
| T2 | Application `TemplateInitService` mit Render-Loop: per `iofs.WalkDir` über die Template-FS, `.tmpl` via `text/template` rendern, sonst 1:1 kopieren; Pfade durch `domain.TemplatePath` validiert. Driving-Port `TemplateInitUseCase.Init(ctx, req)`. Sentinels: `ErrTemplateNotFound`, `ErrTemplateRender` (Exit 14), `ErrInvalidTemplatePath` (Exit 10). Application-Unit-Tests mit Fake-Catalog + Fake-FileSystem. |
| T3 | Bootstrap-Template-Inhalt: `externaltemplates/templates/basic/` bekommt `u-boot.yaml.tmpl`, `compose.yaml.tmpl`, `README.md.tmpl`, `CHANGELOG.md.tmpl`, `.env.example.tmpl`, `.gitignore` — Inhalt aus `internal/hexagon/application/templates/` portiert. Byte-Identity-Pin: Test rendert `basic` und vergleicht gegen die InitProjectService-Default-Outputs. |
| T4 | CLI: `init`-Subkommando bekommt `--template <name>`-Flag (mutex zu `--devcontainer`? T4-Decision). InitProjectRequest erhält optionales `Template` Feld; wenn gesetzt, delegiert InitProjectService an TemplateInitService statt an die Default-File-Plan-Pfad. Wiring in `cmd/uboot/main.go`. CLI-Tests + E2E-Test "Default-Init und --template basic produzieren bytewise identische Outputs". `isValidationError` für ErrTemplateNotFound, ErrInvalidTemplatePath; `isFilesystemError` für ErrTemplateRender. |
| T5 | READMEs + CHANGELOG `## [Unreleased]` + roadmap.md §Nächste Schritte 3 + ADR-0009 §Folgepunkte ✅-Häkchen für template-init + slice-Move `open/` → `done/` mit Tranchen-Commit-Spalte. |

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
- **Refactor: Default-Init auf externe Templates umstellen**: konnte
  man theoretisch — `internal/.../templates/` würde dann nach
  `externaltemplates/templates/basic/` wandern und InitProjectService
  würde standardmäßig `--template basic` laufen. Defer, weil
  Byte-Identity-Pin in T3 die Voraussetzung erst zementiert.

## Bezug

- ADR: [ADR-0009](../../adr/0009-template-format-yaml-files.md)
  §Entscheidung (`text/template`-Engine, Pfad-Eskalation-Validator,
  `embed.FS`-Layout) + §Folgepunkte (`slice-v1-template-init`).
- Voraussetzungs-Slice:
  [`slice-v1-template-list`](../done/slice-v1-template-list.md)
  liefert TemplateCatalog-Port + `basic`-Metadaten.
- Spec: `LH-FA-TPL-001` (V1, MVP), `LH-FA-TPL-002` (V1).
- Architektur: hexagonale Schichten unverletzt; `domain.TemplatePath`
  analog M8 `domain.ConfigPath`-Pattern.
- Phase: V1 (nach `slice-v1-template-list`).
