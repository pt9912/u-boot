# ADR 0009: Externe Template-Format-Konvention — YAML-Metadaten + `text/template`-Files

## Status

Accepted

## Datum

2026-05-31

## Kontext

[`LH-OPEN-004`](../../../spec/lastenheft.md#lh-open-004--template-format-entschieden) (Template-Format) ist in `spec/lastenheft.md` §14 offen:

> *„Das genaue Format für Templates ist noch festzulegen. Mögliche
> Optionen: YAML-Metadaten plus Dateivorlagen / Cookiecutter-kompatible
> Templates / eigenes Template-System / OCI-basierte Template-Pakete."*

[`LH-FA-TPL-001`](../../../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen)..[`LH-FA-TPL-004`](../../../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten) (V1) fordern Projektvorlagen mit Metadaten,
Listing und CLI-Auswahl (`u-boot init --template <name>`,
`u-boot template list`). [`LH-FA-TPL-003`](../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates) (Later) fordert lokale
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

Sicherheits-Rahmen: [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) (MVP) verbietet die
verdeckte Ausführung externer Skripte ohne ausdrückliche
Nutzer-Zustimmung. Template-Engines mit Code-Eval-Fähigkeit
(Jinja2 `{% set %}` mit beliebigem Python-Ausdruck) müssen mit
Sandbox-Setup oder Allowlist arbeiten — Pflicht, kein Optional.

Vergleichbare Tools:

- `cookiecutter` (Python): Jinja2-Templates plus `cookiecutter.json`
  als Variable-Schema. Reife Ökosystem-Reichweite, aber bringt
  vollständige Python-Toolchain als Voraussetzung mit. Pre-/Post-
  Hooks erlauben Python-Skript-Ausführung — direkter Konflikt mit
  [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte).
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

**YAML-Metadaten + `text/template`-Files** (Option 1 aus der
Optionsliste). Jedes externe Template ist ein Verzeichnis mit:

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
  hatte ursprünglich `external-templates/` mit Hyphen vorgesehen; die
  Umsetzung hat auf `externaltemplates/` ohne Hyphen konsolidiert, weil
  alle bestehenden `driven/`-Adapter-Verzeichnisse (clock,
  confirm, docker, fs, git, logger, netprobe, progress, runtime,
  yaml) hyphen-frei sind und Go-Package-Namen ohnehin keine
  Hyphen erlauben.
- **Lokale User-Templates ([`LH-FA-TPL-003`](../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates), Later):** `--template
  ./mein-template` löst gegen das Dateisystem statt `embed.FS`
  auf. Same Schema, same Engine.
- **Pfad-Eskalation verhindert:** beim Rendern werden absolute
  Pfade und `..`-Sequenzen in Template-Datei-Listings strikt
  abgewiesen. Der konkrete Validator nutzt ein Domain-Path-Pattern
  analog zum bestehenden `domain.ConfigPath`.
- **Keine Pre-/Post-Hooks im Template** — kein Code-Eval-Pfad,
  also kein [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)-Risiko. Falls später Setup-Skripte
  gewünscht werden, brauchen sie eine explizite Entscheidung mit eigenem
  Sandbox-Modell.
- **CLI-Subkommandos:**
  - `u-boot template list [--json]` — Listing aus dem Template-
    Katalog ([`LH-FA-TPL-004`](../../../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten)).
  - `u-boot init --template <name|path> [--var key=value …]` —
    Template-getriebener Init ([`LH-FA-TPL-001`](../../../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen)).
  - Beide werden als getrennte Inkremente implementiert; dieses ADR
    liefert nur die Format-Entscheidung.

Cookiecutter, eigenes System und OCI-Pakete werden verworfen:

- **Cookiecutter:** Python-Toolchain-Abhängigkeit (verletzt
  [`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) „möglichst wenige Host-Deps"), Jinja2-Code-Eval
  mit [`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte)-Risiko, doppelter Template-Stack.
- **Eigenes System:** maximal pflegeintensiv, kein erkennbarer
  Vorteil gegenüber `text/template`.
- **OCI-Pakete:** prospektive Architektur ohne Use-Case-Trigger.
  Wenn später ein externer Template-Author einen
  Distributions-Wunsch hat, ist das eine eigene Re-Eval-Entscheidung
  (analog Plugin-System Folgepunkt 1).

## Konsequenzen

Positiv:

- **Eine Template-Engine im Projekt.** M3-T2 und das externe
  System teilen sich Code-Pfade und Test-Patterns; Bug-Fixes in
  `text/template`-Renderer-Helfern wirken überall.
- **Pure Go.** Keine Python-Toolchain im Image, kein
  zusätzlicher Docker-Stage, kein zusätzlicher CI-Job für
  Jinja2-Linting.
- **[`LH-NFA-SEC-004`](../../../spec/lastenheft.md#lh-nfa-sec-004--keine-verdeckte-ausführung-fremder-skripte) trivial erfüllt.** `text/template`-Engine ist
  per Design ohne Code-Eval; ohne Pre-/Post-Hooks gibt es keinen
  Pfad, über den ein Template fremden Code ausführt.
- **[`LH-NFA-PORT-002`](../../../spec/lastenheft.md#lh-nfa-port-002--keine-unnötigen-systemabhängigkeiten) gewahrt.** Keine neuen Host-Voraussetzungen
  über Docker + Make hinaus.
- **Konsistent zum gewählten Stack:** [ADR-0001](0001-implementierungssprache-go.md) (Go), [ADR-0002](0002-hexagonale-architektur.md)
  (hexagonal — Templates sind Driven-Adapter-Ressourcen),
  [ADR-0008](0008-plugin-system-statisch.md) (statische Add-ons, gleiche Verteilungs-Logik).

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
  bestehenden [`LH-FA-CLI-005A`](../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung)-Mode-Flags
  (`--yes`/`--no-interactive`) selbst aufsetzen.

Alternativen (verworfen):

- **Cookiecutter:** siehe oben — Toolchain, Security, doppelter
  Stack.
- **Eigenes Template-System:** unverhältnismäßiger Pflegeaufwand
  ohne Mehrwert.
- **OCI-Template-Pakete:** prospektive Architektur ohne Use-Case;
  Default ist statisch, analog [ADR-0008](0008-plugin-system-statisch.md).

## Folgepunkte

Dieses ADR liefert nur die Format-Entscheidung. Die Implementierung
folgt den `LH-FA-TPL-*`-Anforderungen:

- `u-boot template list [--json]`-Subkommando + `embed.FS`-Katalog-Scan
  ([`LH-FA-TPL-004`](../../../spec/lastenheft.md#lh-fa-tpl-004--templates-auflisten)).
- `u-boot init --template <name>` mit Render-Loop ([`LH-FA-TPL-001`](../../../spec/lastenheft.md#lh-fa-tpl-001--projektvorlagen),
  [`LH-FA-TPL-002`](../../../spec/lastenheft.md#lh-fa-tpl-002--template-metadaten) Metadaten-Surface).
- `--template ./pfad`-Auflösung gegen das Dateisystem
  ([`LH-FA-TPL-003`](../../../spec/lastenheft.md#lh-fa-tpl-003--eigene-templates)).
- Variable-Resolution + Prompt-Pfad bleiben eigener Produktumfang bei
  erstem variable-bedürftigem Built-in.
- Weitere Templates (Micronaut, SvelteKit, …) folgen konkretem Bedarf.

Re-Evaluation-Trigger (analog [ADR-0008](0008-plugin-system-statisch.md)): wenn externer Bedarf an
Cookiecutter-Kompatibilität oder OCI-Verteilung konkret entsteht,
wird eine neue Entscheidung vorbereitet, die dieses ADR ersetzt oder
ergänzt.
