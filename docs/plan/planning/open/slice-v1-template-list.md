# Slice V1: `u-boot template list`

## Auslöser

[ADR-0009](../../adr/0009-template-format-yaml-files.md) §Folgepunkte
führt drei Implementierungs-Slices auf, die erst nach der Format-
Entscheidung sinnvoll waren: `slice-v1-template-list`,
`slice-v1-template-init`, `slice-later-local-templates`. Dieses
Slice ist der erste — kleinster Scope (read-only, kein Render-Pfad),
gute Aufwärm-Tranche für die externe Template-Engine und legt die
Port-/Adapter-Verdrahtung, gegen die T2 (`slice-v1-template-init`)
später anbindet.

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

- `u-boot template list` (ohne Flag) druckt eine tabellarische
  Liste aller eingebauten Templates: Name + Beschreibung +
  Version, eine Zeile pro Template, deterministisch sortiert
  (Namen alphabetisch).
- `u-boot template list --json` druckt ein JSON-Array
  `[{"name":"basic","description":"…","version":"…","supportedAddOns":[…],…},…]`
  mit allen `LH-FA-TPL-002`-Metadatenfeldern (Name, Beschreibung,
  Version, unterstützte Add-Ons, erzeugte Dateien, benötigte
  Tools). Reihenfolge identisch zur Tabellen-Form.
- Hexagonale Verdrahtung: Driven-Port `TemplateCatalog` im
  `internal/hexagon/port/driven/`-Verzeichnis; Adapter
  `internal/adapter/driven/external-templates/` mit `embed.FS`
  (1 Built-in: `basic`); Application-Service
  `TemplateListService` mit Driving-Port `TemplateListUseCase`;
  CLI-Subkommando in `internal/adapter/driving/cli/template.go`.
  `depguard`-Schichten unverletzt.
- Mindestens ein Built-in-Template (`basic`) liegt im Katalog
  als Bootstrap-Stand: `template.yaml` mit den `LH-FA-TPL-002`-
  Pflichtfeldern. Datei-Templates für den `basic`-Init-Pfad
  fügt `slice-v1-template-init` hinzu — dieser Slice braucht sie
  noch nicht (kein Render-Pfad).
- Exit-Codes per `LH-FA-CLI-006`: 0 bei Erfolg, 14 bei
  technischen Adapter-Fehlern (kaputter `template.yaml`-Parse,
  `embed.FS`-IO; sollte in der Test-Suite nie auftreten — wir
  validieren das `embed.FS`-Layout via Integrity-Smoketest beim
  Adapter-Build), 2 bei CLI-Fehlern (unbekanntes Subkommando).

## Tranchen (vorgeschlagen)

| T | Inhalt |
| - | ------ |
| T1 | Domain-Value-Struct `domain.TemplateMetadata` + Driven-Port `port/driven.TemplateCatalog` (`List(ctx) ([]TemplateMetadata, error)`) + Driven-Adapter `internal/adapter/driven/external-templates/` mit `embed.FS`-Scan für `*/template.yaml`, deterministische Sortierung, `template.yaml`-YAML-Parse via `gopkg.in/yaml.v3` (analog zum bestehenden `driven.yaml`-Codec). Bootstrap-Built-in: `basic/template.yaml` mit Pflichtfeldern. Adapter-Tests mit Fixture (mindestens: `basic` wird gefunden + Felder geparst; falsches `template.yaml` produziert klassifizierbaren Fehler). |
| T2 | Driving-Port `port/driving.TemplateListUseCase` (`List(ctx, TemplateListRequest) (TemplateListResponse, error)`) + Sentinel-Sammlung (`ErrTemplateCatalog`); Application-Service `application.TemplateListService` (delegiert an Catalog.List; mappt Adapter-Errors auf Sentinels). Application-Unit-Tests mit Catalog-Fake (mindestens: leerer Katalog → leere Response, Fake-Error → Sentinel-Wrap). |
| T3 | CLI-Subkommando `u-boot template list [--json]` in `internal/adapter/driving/cli/template.go` (`newTemplateCommand` + `newTemplateListCommand` analog `config`-Subbaum). Human-readable: tabellarische Ausgabe mit `text/tabwriter`-Layout (Spalten NAME / DESCRIPTION / VERSION); `--json`: `encoding/json`-Marshal mit `json.Indent`-Pretty-Print. Wiring in `cmd/uboot/main.go` (`external-templates`-Adapter konstruieren, an `TemplateListService` injizieren, an `cli.New` via Option erweitern). CLI-Smoketest + ein E2E-Test pro Format. |
| T4 | README.{md,de.md} bekommen einen kurzen „Templates" -Abschnitt im Quickstart-Kontext (oder unter „Status" — Variante in T4 zu entscheiden); `CHANGELOG.md ## [Unreleased]` Added-Eintrag; `roadmap.md` §Nächste Schritte 3 aktualisiert (template-list ✅, weitere zwei offen); ADR-0009 §Folgepunkte template-list ✅-Marke; Slice-Move `open/` → `done/` mit Tranchen-Tabelle inkl. Commit-Spalte (T4 als „dieser Commit"). |

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
  §Folgepunkte (Implementierungs-Slices), §Entscheidung
  (Metadaten-Schema + Pfad-Layout).
- Spec: `LH-FA-TPL-002` (Template-Metadaten), `LH-FA-TPL-004`
  (Templates auflisten), beide V1 (`spec/lastenheft.md` §4.8).
- Roadmap: §Nächste Schritte Punkt 3 (V1-Templates-Implementation,
  drei ADR-0009-Slices).
- Architektur: hexagonale Schicht-Trennung (LH-FA-ARCH-001..003),
  neuer Driven-Port `TemplateCatalog` analog zum
  `RuntimeEnvironment`-Pattern aus v0.1.1.
- Phase: V1 (nach v0.1.1).
