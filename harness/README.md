# Harness

Dieser Harness verbindet Spezifikationen, ADRs, Slice-Plaene,
Quality-Gates und Betriebsdokumentation fuer `u-boot`. Er ist kein
Ersatz fuer `spec/` oder `docs/`, sondern der Einstiegspunkt fuer
Menschen und AI-Coding-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, gewinnt die
kanonische Quelle und diese Datei wird angepasst.

## Source Precedence

| Rang | Quelle | Charakter |
| --- | --- | --- |
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | Normative Anforderungen, Akzeptanzkriterien, Exit-Codes, Sprachvertrag |
| 2 | [`spec/architecture.md`](../spec/architecture.md) | Hexagonale Architektur, Schichten, Importregeln |
| 3 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 4 | [`docs/plan/planning/in-progress/`](../docs/plan/planning/in-progress/) und [`next/`](../docs/plan/planning/next/) | Aktuelle Slice-Arbeit |
| 5 | [`Makefile`](../Makefile), [`Dockerfile`](../Dockerfile), [`.golangci.yml`](../.golangci.yml), [`.github/workflows/`](../.github/workflows/) | Ausfuehrbare Vertraege |
| 6 | [`docs/user/`](../docs/user/) | Quality, Branch Protection, Nutzerpfade |
| 7 | [`README.md`](../README.md), [`README.de.md`](../README.de.md), [`CHANGELOG.md`](../CHANGELOG.md) | Produktueberblick und Release-Kommunikation |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing und Hard Rules |
| 9 | Diese Datei | Harness-Einstieg |

## Guides

Feedforward-Quellen, die Arbeit vor der Umsetzung lenken:

| Quelle | Inhalt |
| --- | --- |
| [`spec/lastenheft.md`](../spec/lastenheft.md) | `LH-*`-IDs, Prioritaeten, funktionale und nicht-funktionale Anforderungen |
| [`spec/architecture.md`](../spec/architecture.md) | Layering, Port-/Adapter-Regeln, depguard-Kontrakt |
| [`docs/plan/adr/README.md`](../docs/plan/adr/README.md) | ADR-Index und Entscheidungsueberblick |
| [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Release- und Slice-Status |
| [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md) | Temporaere und permanente Carveouts |
| [`docs/user/quality.md`](../docs/user/quality.md) | Quality-Gates, Linter-Profil, Coverage, Security |
| [`harness/roles.md`](roles.md) | Rollen, Uebergaben und Konfliktpfade |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Minimal Workflow |

## Sensors

Feedback-Gates, die reale Projektzustaende messen:

| Target | Charakter | Wann verwenden |
| --- | --- | --- |
| `make lint` | Computational feedback: statische Analyse, depguard, SOLID-nahe Linter | Nach Go-Code- und Architektur-Aenderungen |
| `make test` | Computational feedback: Unit- und Default-Tests im Docker-Test-Stage | Nach Codeaenderungen |
| `make test-docker` | Computational feedback: Docker-tag Integrationstests | Nach Docker-/Compose-/E2E-Aenderungen |
| `make coverage-gate` | Computational feedback: Coverage-Schwelle, Default 90 Prozent | Nach produktiven Codeaenderungen |
| `make docs-check` | Computational feedback: Markdown-Link- und Pfadpruefung | Nach Doku-, Spec-, ADR- oder Planning-Aenderungen |
| `make govulncheck` | Computational feedback: Go-Vulnerability-Scan | Vor CI-/Release-Handoff |
| `make image-scan` | Computational feedback: Trivy gegen Runtime-Image | Vor CI-/Release-Handoff |
| `make verify-depguard` | Computational feedback: depguard-Regeln feuern wirklich | Bei Aenderungen an Layern oder depguard-Konfig |
| `make gates` | Inner-loop Closure: lint + test + coverage-gate + docs-check | Normaler Abschluss fuer Codeaenderungen |
| `make ci` | CI-Spiegel: gates + govulncheck + image-scan | Vor groesseren Handoffs oder Releases |
| `make fullbuild` | Voller Buildabschluss: ci + Runtime-Image | Vor Release-Closure |

Wenn ein Sensor wegen Umgebung oder Sandbox nicht laeuft, den Grund im
Handoff nennen. Keine gruene Closure behaupten, wenn der passende Sensor
nicht ausgefuehrt wurde.

## Traceability

- Jede oeffentliche Verhaltensaenderung braucht einen `LH-*`-, `ADR-*`-
  oder Slice-Anker.
- Neue oder geaenderte Anforderungen brauchen einen Nachweis: Test,
  Gate, Demo, ADR oder dokumentierte Closure.
- Neue ADRs muessen den ADR-Index aktualisieren.
- Planning-Dokumente folgen `open/ -> next/ -> in-progress/ -> done/`.
- Temporaere Carveouts brauchen parallel Inventar-Eintrag und
  Plan-Anker.

## Role Separation

Rollen sind Kontextgrenzen, keine Personen. Die verbindliche
Rollenreferenz liegt in [`harness/roles.md`](roles.md).

Standardsequenz fuer Slice-Arbeit:

```text
Planner -> Architect -> Implementation -> Reviewer -> Verifier -> Validator -> Planner
```

Jeder Rollenwechsel braucht ein Uebergabe-Artefakt. Eine Rolle darf in
eine fruehere Rolle zurueckgeben, aber nicht deren Entscheidung
stillschweigend ersetzen.

## Scope Boundaries

- `u-boot` ist ein CLI zum Bootstrapping reproduzierbarer Docker-
  Entwicklungsumgebungen, kein allgemeiner Project-Generator ohne
  Docker-/Compose-Vertrag.
- Application-Code bleibt frei von konkreter externer I/O; I/O sitzt in
  Driven-Adaptern und wird ueber Ports erreicht.
- Generated files und User-Projektdateien sind sicherheitsrelevant:
  managed blocks, Backups, Two-Phase-Planung und Bestaetigungen sind
  Produktvertraege, keine Komfortdetails.
- CLI-Output und generierte Artefakte sind Englisch; normative Specs und
  Planning-Dokumente koennen Deutsch bleiben.
- Release- und Distributionsaenderungen muessen ADR-0004/ADR-0007,
  CI-Gates und README/CHANGELOG zusammen betrachten.

## Minimal Agent Workflow

1. Diese Datei und [`AGENTS.md`](../AGENTS.md) lesen.
2. Rolle aus [`harness/roles.md`](roles.md) bestimmen.
3. Relevante Spec, Architektur, ADR und aktiven Slice lesen.
4. Betroffene IDs und Produktvertraege benennen.
5. Kleinste sinnvolle Aenderung ausfuehren.
6. Engsten nuetzlichen Sensor laufen lassen.
7. Bei Codeaenderungen nach Moeglichkeit `make gates` ausfuehren.
8. Oeffentliche Doku, Planning-Artefakte und CHANGELOG aktualisieren,
   wenn ein oeffentlicher Vertrag beruehrt wurde.
9. Handoff mit Rolle, ausgefuehrten Sensors, offenen Sensors und Risiken.
