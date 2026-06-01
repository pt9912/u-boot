# Slice V1: `u-boot template list`

## Auslöser

[ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte
führt drei Implementierungs-Slices auf, die erst nach der Format-
Entscheidung sinnvoll waren: `slice-v1-template-list`,
`slice-v1-template-init`, `slice-later-local-templates`. Dieses
Slice ist der erste — kleinster Scope (read-only, kein Render-Pfad),
gute Aufwärm-Tranche für die externe Template-Engine und legt die
Port-/Adapter-Verdrahtung, gegen die `slice-v1-template-init` später
anbindet.

Spec-IDs: `LH-FA-TPL-004` (Templates auflisten, V1) sowie
implizit `LH-FA-TPL-002` (Template-Metadaten — mindestens Name +
Beschreibung + Version sind in der Listing-Ausgabe sichtbar).

## Aufhebungsbedingung

`u-boot template list` listet alle eingebauten Templates aus dem
`embed.FS`-Katalog mit den `LH-FA-TPL-004`-Pflichtfeldern (Name,
Beschreibung, Version). `--json` liefert dieselben Daten in
strukturierter Form für maschinelle Verwendung. Mindestens ein
Built-in-Template (`basic`) liegt im Katalog (Bootstrap-Stand,
analog zur Mindestauslieferung in ADR-0009 §Folgepunkte). Der
hexagonale Port `driven.TemplateCatalog` ist als
Lese-Schnittstelle gegen den `slice-v1-template-init`-Render-Pfad
wiederverwendbar.

## Akzeptanzkriterien

- ✅ `u-boot template list` (ohne Flag) druckt eine tabellarische
  Liste aller eingebauten Templates: Name + Beschreibung +
  Version, eine Zeile pro Template, deterministisch sortiert
  (Namen alphabetisch). Tabwriter-Layout mit padding=2.
- ✅ `u-boot template list --json` druckt ein JSON-Array
  `[{"name":"basic","description":"…","version":"…","supportedAddOns":[…],…}]`
  mit allen `LH-FA-TPL-002`-Metadatenfeldern. Reihenfolge identisch
  zur Tabellen-Form. Nil-Slices werden auf `[]` normalisiert
  (kein `null` im JSON).
- ✅ Hexagonale Verdrahtung: Driven-Port `TemplateCatalog` im
  `internal/hexagon/port/driven/`-Verzeichnis; Adapter
  `internal/adapter/driven/externaltemplates/` mit `embed.FS`
  (1 Built-in: `basic`); Application-Service
  `TemplateListService` mit Driving-Port `TemplateListUseCase`;
  CLI-Subkommando in `internal/adapter/driving/cli/template.go`.
  `depguard`-Schichten unverletzt.
- ✅ Mindestens ein Built-in-Template (`basic`) liegt im Katalog
  als Bootstrap-Stand: `template.yaml` mit den `LH-FA-TPL-002`-
  Pflichtfeldern. Datei-Templates für den `basic`-Init-Pfad
  fügt `slice-v1-template-init` hinzu — dieser Slice braucht sie
  noch nicht (kein Render-Pfad).
- ✅ Exit-Codes per `LH-FA-CLI-006`: 0 bei Erfolg, 14 bei
  technischen Adapter-Fehlern (`driving.ErrTemplateCatalog`
  in `isFilesystemError` aufgenommen), 2 bei CLI-Fehlern.

## Tranchen

| T | Commit | Inhalt |
| - | ------ | ------ |
| T1 | `65795b5` | Domain-Value-Struct `domain.TemplateMetadata` + Driven-Port `port/driven.TemplateCatalog` + Driven-Adapter `internal/adapter/driven/externaltemplates/` mit `embed.FS`-Scan (`templates/*/template.yaml`), deterministische Sortierung, YAML-Parse via `gopkg.in/yaml.v3` über privater `rawTemplateYAML`-Projektion. `domain.ErrInvalidTemplate`-Sentinel + `Validate()`-Methode (kebab-case-Name, Description, Version pflicht). Bootstrap-Built-in: `templates/basic/template.yaml` mit den `LH-FA-TPL-002`-Pflichtfeldern + `supportedAddOns: [postgres]` + `generatedFiles` analog dem heutigen u-boot init-Ausstoß. 8 Adapter-Tests + Domain-Validate-Coverage 100%. |
| T2 | `a099d63` | Driving-Port `port/driving.TemplateListUseCase` (`List(ctx, TemplateListRequest) (TemplateListResponse, error)`) + `ErrTemplateCatalog`-Sentinel (Exit-Code 14). Application-Service `application.TemplateListService` — thin Pass-through über den Catalog mit Multi-`%w`-Wrap, damit `errors.Is(err, domain.ErrInvalidTemplate)` durch die Service-Schicht überlebt. Nil-Logger-Fallback auf `noopLogger`. 6 Application-Unit-Tests (Delegate, Leerer Katalog, Sentinel-Wrap, Cause-Preservation, Context-Propagation, Nil-Logger-Construction). |
| T3 | `23bd91b` | CLI-Subkommando `u-boot template list [--json]` in `internal/adapter/driving/cli/template.go` (`newTemplateCommand` parent + `newTemplateListCommand` leaf analog `config`-Subbaum). Human-Output via `text/tabwriter` (NAME/DESCRIPTION/VERSION-Header, padding=2); JSON-Output via `encoding/json.MarshalIndent` auf eine CLI-lokale DTO (`templateJSON`), damit die Domain-Schicht presentation-agnostic bleibt (ADR-0002 / LH-FA-ARCH-002). Nil-Slices auf `[]` normalisiert. Wiring in `cmd/uboot/main.go` (`externaltemplates.New()` → `NewTemplateListService` → 9. Positional in `cli.New`). `cli.New`-Signatur um `tmplUC` erweitert; alle 7 bestehenden `newApp*`-Test-Helper + `fakeTemplateListUseCase` aktualisiert. `ErrTemplateCatalog` zu `isFilesystemError` und ExitCode-Doc-Comment hinzugefügt. 6 CLI-Tests (Human, Leerer Katalog, JSON-Roundtrip, JSON-Leer-Array, Error-→-Exit-14, Help-Listet-list). Smoketest gegen das gebaute Image: `u-boot template list` → tabellarisch, `u-boot template list --json` → JSON-Array mit den `basic`-Metadaten. |
| T4 | `a7e0d7b` | Slice-Plan nach `done/`; README.{md,de.md} bekommen einen neuen Bullet für `u-boot template list`; `CHANGELOG.md ## [Unreleased]` Added-Eintrag; `roadmap.md` §Nächste Schritte 3 mit T1-T3-Hashes aktualisiert + slice-v1-template-list ✅-Häkchen; ADR-0009 §Folgepunkte `slice-v1-template-list` ✅-Häkchen + §Entscheidung-Pfad-Korrektur (`external-templates/` → `externaltemplates/`, mit kurzer Begründung). `make docs-check` grün. |
| Review | dieser Commit | Code-Review-Followup: fünf Findings (N1..N5) direkt am Slice gepatcht — yaml `KnownFields(true)` gegen Tippfehler, Regex-Verschärfung gegen konsekutive Bindestriche, `default,omitempty` entfernt, apiVersion-Gate, embed.FS-Adapter ehrt jetzt ctx. Siehe Review-Followup-Tabelle unten. |

## Review-Followup (post-T4)

Code-Review (xhigh-effort, 5 Angles + Sweep) hat fünf
Inkonsistenzen zwischen Doku-Anspruch und Implementierung
gefunden. Im Followup-Commit gefixt (gleicher Branch, kein
eigener Slice — `make gates` grün nach den Patches):

| N | Datei | Befund | Fix |
| - | ----- | ------ | --- |
| N1 | `externaltemplates/catalog.go` | `yaml.Unmarshal` schluckte unbekannte Felder: `requiredTool:` (singular) wäre stillschweigend verschwunden. | `yaml.NewDecoder(...).KnownFields(true).Decode(...)` — Typos failen jetzt am Load. Neuer Test `TestCatalog_List_UnknownYAMLFieldRejected`. |
| N2 | `domain/template_metadata.go` | Regex `^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$` erlaubte konsekutive Bindestriche (`my--bad`). | Regex auf `^[a-z0-9]+(?:-[a-z0-9]+)*$` verschärft — Einzel-Bindestrich-Segmente erzwungen. Neuer Test-Case (`consecutive dashes rejected (review-followup N2)`) + Leading-dash-Test. |
| N3 | `cli/template.go` | `json:"default,omitempty"` zerstörte gerade die Missing-vs-Empty-Unterscheidung, die der Kommentar zu bewahren versprach. | `,omitempty` entfernt; Kommentar umgekehrt. Bestehende Tests passieren unverändert (Fixtures benutzen non-empty Defaults). |
| N4 | `externaltemplates/catalog.go` | `rawTemplateYAML.APIVersion` wurde geparst aber nie validiert; ein Template mit `apiVersion: …/v999` wäre durchgerutscht. | Konstante `supportedAPIVersion` + Gate in `readTemplate` mit `domain.ErrInvalidTemplate`-Wrap. Neuer Test `TestCatalog_List_UnsupportedAPIVersionRejected`. |
| N5 | `externaltemplates/catalog.go` + `port/driven/template_catalog.go` | Port-Doc versprach `ctx is honored`, der embed.FS-Adapter ignorierte ihn aber explizit (siehe „ignores it"-Sonderfall). | Adapter macht jetzt `ctx.Err()`-Check am Entry; Port-Doc gleicht den Sonderfall aus. Neuer Test `TestCatalog_List_HonorsCancelledContext`. |

Gates nach Followup grün: lint + test + coverage-gate 90.30%
(>= 90%) + docs-check. Smoke-Test gegen `docker run --rm u-boot
template list` und `... --json` zeigt unverändertes Output für
das `basic`-Bootstrap-Template (keine Wire-Shape-Regression).

## Out of Scope

- `u-boot init --template <name>` Render-Pfad — eigener Slice
  `slice-v1-template-init` (ADR-0009 §Folgepunkte 2).
- `u-boot init --template ./pfad` Filesystem-Auflösung — eigener
  Slice `slice-later-local-templates` (ADR-0009 §Folgepunkte 3).
- Weitere Built-in-Templates über `basic` hinaus (`micronaut`,
  `sveltekit`, …) — kommen mit konkretem Bedarf in eigenen Slices
  (ADR-0009 §Folgepunkte 4).
- Lokal-Template-Discovery (Filesystem-Walk außerhalb `embed.FS`)
  — gehört zum Later-Slice.
- Variable-Resolution / Render-Loop / `text/template`-Execution
  — gehört zum `slice-v1-template-init`.

## Bezug

- ADR: [ADR-0009](../../adr/0009-template-format-yaml-files.md)
  §Folgepunkte (Implementierungs-Slices, mit T4 ✅-Häkchen für
  diesen Slice), §Entscheidung (Metadaten-Schema + Pfad-Layout —
  Pfad mit T4 auf `externaltemplates/` ohne Hyphen konsolidiert).
- Spec: `LH-FA-TPL-002` (Template-Metadaten), `LH-FA-TPL-004`
  (Templates auflisten), beide V1 (`spec/lastenheft.md` §4.8).
- Roadmap: §Nächste Schritte Punkt 3 (V1-Templates-Implementation,
  drei ADR-0009-Slices). Mit T4 auf
  „slice-v1-template-list ✅ geliefert" aktualisiert.
- Architektur: hexagonale Schicht-Trennung (LH-FA-ARCH-001..003),
  neuer Driven-Port `TemplateCatalog` analog zum
  `RuntimeEnvironment`-Pattern aus v0.1.1.
- Phase: V1 (nach v0.1.1).
