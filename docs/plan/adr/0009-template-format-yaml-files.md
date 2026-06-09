# ADR 0009: Externe Template-Format-Konvention — YAML-Metadaten + `text/template`-Files

## Status

Accepted

## Datum

2026-05-31

## Kontext

`LH-OPEN-004` (Template-Format) ist in `spec/lastenheft.md` §14 offen:

> *„Das genaue Format für Templates ist noch festzulegen. Mögliche
> Optionen: YAML-Metadaten plus Dateivorlagen / Cookiecutter-kompatible
> Templates / eigenes Template-System / OCI-basierte Template-Pakete."*

`LH-FA-TPL-001..004` (V1) fordern Projektvorlagen mit Metadaten,
Listing und CLI-Auswahl (`u-boot init --template <name>`,
`u-boot template list`). `LH-FA-TPL-003` (Later) fordert lokale
User-Templates (`u-boot init --template ./my-template`).

**Vorhandenes Pattern (M3-T2):** u-boot nutzt intern bereits
`text/template + embed.FS` für die *applikativen* Init-Templates
(`u-boot.yaml`, `compose.yaml`, `README.md`, …, alle in
`internal/adapter/driven/templates/`). Das ist gut etablierter
Go-Stack ohne externe Toolchain-Abhängigkeit, geprüft durch
`make gates` + e2e-Tests aus M3.

Das **externe** Template-System (`LH-FA-TPL-*`) ist ein anderer
Lebenszyklus: es liefert keine u-boot-Codebase-Templates, sondern
*User-Projekt-Templates* (basic, micronaut, sveltekit, …). Trotzdem
gilt: jede zweite Engine im Projekt erhöht den Wartungsaufwand
und die Test-Surface.

Sicherheits-Rahmen: `LH-NFA-SEC-004` (MVP) verbietet die
verdeckte Ausführung externer Skripte ohne ausdrückliche
Nutzer-Zustimmung. Template-Engines mit Code-Eval-Fähigkeit
(Jinja2 `{% set %}` mit beliebigem Python-Ausdruck) müssen mit
Sandbox-Setup oder Allowlist arbeiten — Pflicht, kein Optional.

Vergleichbare Tools:

- `cookiecutter` (Python): Jinja2-Templates plus `cookiecutter.json`
  als Variable-Schema. Reife Ökosystem-Reichweite, aber bringt
  vollständige Python-Toolchain als Voraussetzung mit. Pre-/Post-
  Hooks erlauben Python-Skript-Ausführung — direkter Konflikt mit
  `LH-NFA-SEC-004`.
- `gh repo create --template` (GitHub): Repository-Templates,
  keine Variable-Substitution. Reduzierte Mächtigkeit, einfach,
  aber nicht für lokale Initialisierung gedacht.
- `helm`: Go-Templates plus `Chart.yaml` Metadaten plus
  `values.yaml` für Variablen. Pures Go-Ökosystem, gut etabliert.
- `kubectl create`-Generators: hartkodierte Go-Strukturen ohne
  Template-Files; flexibel nicht erweiterbar.

In u-boots Go-/Solo-Projekt-Kontext ist die `helm`-Lösung das
nächste Idiom: `text/template` plus YAML-Metadaten.

## Entscheidung

**YAML-Metadaten + `text/template`-Files** (Option 1 aus dem Slice-
Plan). Jedes externe Template ist ein Verzeichnis mit:

- `template.yaml` — Metadaten + Variable-Schema, von u-boot beim
  Listing (`u-boot template list`) und beim Init geparst.
- Datei-Templates (`*.tmpl` für gerenderte Files, einfache Files
  ohne Suffix werden 1:1 kopiert) unter dem Template-Wurzel-
  Verzeichnis. `text/template`-Syntax (`{{ .ProjectName }}`),
  exakt das M3-T2-Pattern, nur auf externe Templates ausgedehnt.

Konkrete Setzungen:

- **Engine: Go `text/template`** — identisch zur M3-T2-Implementierung.
  Eine Engine im Projekt, kein zweiter Stack zum Lernen.
- **Metadaten-Schema (`template.yaml` v1):**
  ```yaml
  apiVersion: github.com/pt9912/u-boot/template/v1
  name: micronaut
  description: "Micronaut starter project (Java, Gradle)."
  version: 1.0.0
  supportedAddOns: [postgres, keycloak]
  generatedFiles:
    - build.gradle
    - src/main/java/Application.java
  requiredTools:
    - jdk:>=21
  variables:
    - name: groupId
      description: "Maven group ID"
      default: "com.example"
      required: true
  ```
- **Built-in-Templates:** liegen unter
  `internal/adapter/driven/externaltemplates/` als `embed.FS`-
  Verzeichnisse (eines pro Template). Listing via
  `template list` enumeriert die `embed.FS`-Wurzel; `init
  --template <name>` löst über den Namen auf. Hinweis: dieses ADR
  hatte ursprünglich `external-templates/` mit Hyphen vorgesehen;
  [`slice-v1-template-list`](../planning/done/slice-v1-template-list.md)
  T1 hat auf `externaltemplates/` ohne Hyphen konsolidiert, weil
  alle bestehenden `driven/`-Adapter-Verzeichnisse (clock,
  confirm, docker, fs, git, logger, netprobe, progress, runtime,
  yaml) hyphen-frei sind und Go-Package-Namen ohnehin keine
  Hyphen erlauben.
- **Lokale User-Templates (`LH-FA-TPL-003`, Later):** `--template
  ./mein-template` löst gegen das Dateisystem statt `embed.FS`
  auf. Same Schema, same Engine.
- **Pfad-Eskalation verhindert:** beim Rendern werden absolute
  Pfade und `..`-Sequenzen in Template-Datei-Listings strikt
  abgewiesen. Der konkrete Validator wird im
  `slice-v1-template-init` als `domain.TemplatePath` (analog zum
  bestehenden `domain.ConfigPath`-Pattern aus M8) neu eingeführt,
  weil es heute keinen passenden Domain-Path-Validator gibt; die
  ADR-Setzung ist der Vertrag, dass der Slice ihn liefert.
- **Keine Pre-/Post-Hooks im Template** — kein Code-Eval-Pfad,
  also kein `LH-NFA-SEC-004`-Risiko. Falls später Setup-Skripte
  gewünscht werden, wären sie ein expliziter Slice mit eigenem
  Sandbox-Modell.
- **CLI-Subkommandos:**
  - `u-boot template list [--json]` — Listing aus dem Template-
    Katalog (`LH-FA-TPL-004`).
  - `u-boot init --template <name|path> [--var key=value …]` —
    Template-getriebener Init (`LH-FA-TPL-001`).
  - Beide werden als getrennte Slice-Pläne implementiert (siehe
    §Folgepunkte), dieses ADR liefert nur die Format-Entscheidung.

Cookiecutter, eigenes System und OCI-Pakete werden verworfen:

- **Cookiecutter:** Python-Toolchain-Abhängigkeit (verletzt
  `LH-NFA-PORT-002` „möglichst wenige Host-Deps"), Jinja2-Code-Eval
  mit `LH-NFA-SEC-004`-Risiko, doppelter Template-Stack.
- **Eigenes System:** maximal pflegeintensiv, kein erkennbarer
  Vorteil gegenüber `text/template`.
- **OCI-Pakete:** prospektive Architektur ohne Use-Case-Trigger.
  Wenn später ein externer Template-Author einen
  Distributions-Wunsch hat, ist das ein eigener Re-Eval-Slice
  (analog Plugin-System Folgepunkt 1).

## Konsequenzen

Positiv:

- **Eine Template-Engine im Projekt.** M3-T2 und das externe
  System teilen sich Code-Pfade und Test-Patterns; Bug-Fixes in
  `text/template`-Renderer-Helfern wirken überall.
- **Pure Go.** Keine Python-Toolchain im Image, kein
  zusätzlicher Docker-Stage, kein zusätzlicher CI-Job für
  Jinja2-Linting.
- **`LH-NFA-SEC-004` trivial erfüllt.** `text/template`-Engine ist
  per Design ohne Code-Eval; ohne Pre-/Post-Hooks gibt es keinen
  Pfad, über den ein Template fremden Code ausführt.
- **`LH-NFA-PORT-002` gewahrt.** Keine neuen Host-Voraussetzungen
  über Docker + Make hinaus.
- **Konsistent zum gewählten Stack:** ADR-0001 (Go), ADR-0002
  (hexagonal — Templates sind Driven-Adapter-Ressourcen),
  ADR-0008 (statische Add-ons, gleiche Verteilungs-Logik).

Negativ / Trade-offs:

- **Keine Cookiecutter-Reichweite.** Existierende Cookiecutter-
  Templates aus dem Python-Ökosystem müssen manuell auf das
  YAML+Go-Format portiert werden. Mit dem aktuell sehr begrenzten
  geplanten Template-Katalog (basic, micronaut, sveltekit,
  micronaut-sveltekit) ist die Portier-Last überschaubar.
- **Variable-Validierung selbst implementiert.** Cookiecutter
  hat eine entwickelte Variable-/Hook-Validierung; bei der
  YAML-Lösung muss u-boot Default-Werte, Required-Flags und
  Type-Checks selbst pflegen (über `template.yaml` Schema +
  Validator-Helfer).
- **Kein integrierter Variable-Prompt-UI.** Cookiecutter prompted
  interaktiv für Variablen; u-boot muss diesen Pfad über die
  bestehenden `LH-FA-CLI-005A`-Mode-Flags
  (`--yes`/`--no-interactive`) selbst aufsetzen.

Alternativen (verworfen):

- **Cookiecutter:** siehe oben — Toolchain, Security, doppelter
  Stack.
- **Eigenes Template-System:** unverhältnismäßiger Pflegeaufwand
  ohne Mehrwert.
- **OCI-Template-Pakete:** prospektive Architektur ohne Use-Case;
  Default ist statisch, analog ADR-0008.

## Folgepunkte (Implementierungs-Slices)

Dieses ADR liefert nur die Format-Entscheidung. Die Implementierung
braucht eigene Slices:

- ✅ [`slice-v1-template-list`](../planning/done/slice-v1-template-list.md)
  — `u-boot template list [--json]`-Subkommando + `embed.FS`-
  Katalog-Scan (`LH-FA-TPL-004`). Geliefert 2026-06-01 (T1
  `65795b5` Domain+Port+Adapter, T2 `a099d63` Use-Case+Service,
  T3 `23bd91b` CLI+Wiring, T4 Slice-Closure).
- ✅ [`slice-v1-template-init`](../planning/done/slice-v1-template-init.md)
  — `u-boot init --template <name>` mit Render-Loop
  (`LH-FA-TPL-001`, `LH-FA-TPL-002` Metadaten-Surface).
  Geliefert 2026-06-01 (T1 `9e81b02` Domain `TemplatePath` +
  Driven-Port `TemplateFiles` + Adapter `Open`, T2 `65a1ce8`
  Driving-Port + `TemplateInitService` mit Walk-Render-Skip-Loop,
  T3 `ed6d9a0` `basic`-Bootstrap-Content + Byte-Identity-Pin,
  T4 `daaaa9a` CLI-Flag + `InitProjectService`-Delegation + E2E,
  T5 Slice-Closure). Variable-Resolution + Prompt-Pfad bewusst
  out-of-scope — `basic` hat `variables: []`; eigene Folge-Slice
  bei erstem variable-bedürftigem Built-in.
- ✅ [`slice-later-local-templates`](../planning/done/slice-later-local-templates.md)
  — `--template ./pfad`-Auflösung gegen das Dateisystem
  (`LH-FA-TPL-003`). Geliefert (T1–T5): geteiltes `templateyaml`-Paket,
  reine `domain.ClassifyTemplateRef`-Klassifikation, `localtemplates`-
  FS-Resolver + Composite (zweite `TemplateFiles`-Impl), Render-Loop
  wiederverwendet mit Symlink-Guard + `ErrTemplateInvalid`-Exit-10-
  Pfad. `--var`-Variablen out-of-scope (eigener Folge-Slice).
- ✅ Built-in `basic` als Bootstrap-Stand wurde mit
  `slice-v1-template-list` T1 ausgeliefert
  (`templates/basic/template.yaml`); weitere Templates (Micronaut,
  SvelteKit, …) je nach konkretem Bedarf in eigenen Slices.

Re-Evaluation-Trigger (analog ADR-0008): wenn externer Bedarf an
Cookiecutter-Kompatibilität oder OCI-Verteilung konkret entsteht,
wird ein neuer Slice eröffnet, der dieses ADR ersetzt oder
ergänzt.

Der carveouts-Eintrag `LH-OPEN-004` wird mit der
[`slice-v1-template-format-entscheidung`](../planning/done/slice-v1-template-format-entscheidung.md)-
Closure entfernt.
