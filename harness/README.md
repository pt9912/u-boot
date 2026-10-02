# Harness

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area, Zusatzklassen für Sensors-Bindung) sowie Adaptionen ggü.
der adoptierten Baseline leben in [`conventions.md`](conventions.md).
Diese Datei dupliziert sie nicht.

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend (Anforderungen, Akzeptanzkriterien, Exit-Code- und Sprachvertrag) |
| 2 | [`spec/spezifikation.md`](../spec/spezifikation.md) | technisch fortschreibbar (Algorithmen, Schemata, Defaults, Fehler-Codes, externe Verträge) |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/`](../docs/plan/planning/in-progress/) und [`next/`](../docs/plan/planning/next/) | Wellen-Sequenz (`roadmap.md`) und aktive Slice-Arbeit |
| 6 | [`Makefile`](../Makefile), [`Dockerfile`](../Dockerfile), [`.golangci.yml`](../.golangci.yml), [`.github/workflows/`](../.github/workflows/) | ausführbare Verträge |
| 7 | [`docs/user/`](../docs/user/) und [`docs/maintainer/`](../docs/maintainer/) | Nutzerhandbuch; Quality, Branch Protection, Release |
| 8 | [`README.md`](../README.md), [`README.de.md`](../README.de.md), [`CHANGELOG.md`](../CHANGELOG.md) | Projekt-Überblick und Release-Kommunikation |
| 9 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing und Hard Rules |
| 10 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die drei Spec-Straten (Vertrag, Technik, Sicht). Die Ränge 5
> und 6 (Roadmap/aktive Slices und ausführbare Verträge) sind eine repo-lokale
> Ergänzung der Baseline-Tabelle; Begründung in
> [`MR-001`](conventions/MR-001-source-precedence-drei-straten.md).

## Guides (Feedforward-Quellen)

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen (`LH-*`), Akzeptanzkriterien |
| [`spec/spezifikation.md`](../spec/spezifikation.md) | technische Details (`SPEC-*`, Verfeinerungen), Defaults |
| [`spec/architecture.md`](../spec/architecture.md) | Komponenten (`ARC-*`), Schichten, Importregeln, depguard-Kontrakt |
| [`docs/plan/adr/`](../docs/plan/adr/README.md) | Architekturentscheidungen, ADR-Index |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md) | temporäre und permanente Carveouts |
| [`docs/maintainer/quality.md`](../docs/maintainer/quality.md) | Quality-Gates, Linter-Profil, Coverage, Security |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| [`roles.md`](roles.md) | Rollen, Übergaben und Konfliktpfade |
| [`review.md`](review.md) | Review-Kategorien, Prüflinsen und Output-Schema |
| [`replay.md`](replay.md) | Replay-/Golden-Set-Regeln für CLI-Generatoren |
| [`verification.md`](verification.md) | Verification-Evidence und Slice-Closure-Schema |
| `.harness/skills/reviewer.md` | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema (Modul 10) — nächste Rolle nach Schritt 8 des Minimal Agent Workflow, nicht Teil der Implementer-Eingabe |
| `.harness/baseline/v6.13.0/regelwerk/` (vendored; `README.md` = Index) | adoptiertes Betriebsregelwerk in Agenten-Kurzform — **präsente nachschlagbare Vertiefung**, pro Entscheidung abschnittsweise (siehe [`AGENTS.md`](../AGENTS.md) §1); derivativ, Stand/Tag siehe [`conventions.md`](conventions.md) §Baseline |
| `.harness/baseline/v6.13.0/templates/` (vendored, parallel) | Referenz-Form der Skelette, auf die das Regelwerk mit `../templates/…` als „Ziel-Form“ verweist (netzlos, weil parallel zu `regelwerk/`); Vorlagen zum Kopieren-und-Ausfüllen |

## Sensors (Feedback-Gates)

| Target | Vertrag | Bindung |
|---|---|---|
| `make docs-check` | Doku-Referenzen (d-check): Link-Pfade, Heading-Anker, ADR-/LH-/SPEC-/Planning-Kennungs-Links, Referenzmodell, Planning-Lifecycle, Gate-Index | [`ADR-0013`](../docs/plan/adr/0013-dokumentationsreferenzmodell.md) |
| `make lint` | statische Analyse, `depguard`, SOLID-nahe Linter; Verstöße sind PR-blockierend | [`ADR-0003`](../docs/plan/adr/0003-solid-nahes-lint-profil.md) |
| `make test` | Unit- und Default-Tests im Docker-Test-Stage | — |
| `make test-docker` | Integrationstests (Build-Tag `docker`) gegen eine echte Docker Engine | — |
| `make coverage-gate` | Coverage-Schwelle, bootstrap-aware | Schwelle 90 % |
| `make govulncheck` | Go-Vulnerability-Scan | [`ADR-0004`](../docs/plan/adr/0004-ci-system.md) |
| `make image-scan` | Trivy HIGH/CRITICAL gegen das Runtime-Image | [`ADR-0004`](../docs/plan/adr/0004-ci-system.md) |
| `make verify-depguard` | on-demand: die `depguard`-Regeln feuern wirklich | [`ADR-0003`](../docs/plan/adr/0003-solid-nahes-lint-profil.md) |
| `make doc-immutable RANGE="<base>..<head>"` | on-demand: Accepted-ADRs sind über eine Commit-Range unverändert (lokal `STAGED=1`); Aufruf mit Argument, daher nicht im Makefile-Gate-Index | [`MR-005`](conventions.md#mr-005) |
| `make gates` | alle inneren Gates: `lint` + `test` + `coverage-gate` + `docs-check` | — |
| `make ci` | `gates` + `govulncheck` + `image-scan` | [`ADR-0004`](../docs/plan/adr/0004-ci-system.md) |
| `make fullbuild` | volle Closure: `ci` + Runtime-Image-Build | — |

**Aktueller Lauf-Status:** CI-Badge bzw. lokal `make help` / `make gates`.
**Nicht behauptet** (geplant): — keine —.
**Rote Gates:** keine strukturell roten; Begründung eines etwaigen Carveouts im Master-Inventar [`carveouts.md`](../docs/plan/planning/in-progress/carveouts.md).
Läuft ein Sensor wegen Umgebung oder Sandbox nicht, wird der Grund im Handoff genannt; eine grüne Closure wird nicht behauptet, wenn der passende Sensor nicht lief.

## Traceability rules

- PRs/Commits **müssen** mindestens eine `LH-*`-, `ADR-*`- oder Slice-ID nennen.
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Dokument-Referenzen folgen dem Referenzmodell ([`ADR-0013`](../docs/plan/adr/0013-dokumentationsreferenzmodell.md)): Normative Kraft nur auf aufwärtsgerichteten Kanten; Slice-, Carveout- und Roadmap-Kanten sind Kontext.
- Slice-Closure braucht Verification-Evidence nach [`verification.md`](verification.md); Generator-Änderungen brauchen Replay-/Golden-Evidence nach [`replay.md`](replay.md).
- Neue ADRs müssen im ADR-Index ergänzt werden.
- Änderungen an Planning-Dokumenten müssen die Lifecycle-Regeln beachten (open → next → in-progress → done; reine `git mv`-Commits siehe [`AGENTS.md`](../AGENTS.md) §3.3).
- Temporäre Carveouts brauchen parallel einen Inventar-Eintrag und einen Plan-Anker.

## Safety and scope boundaries

- `u-boot` ist ein CLI zum Bootstrapping reproduzierbarer Docker-Entwicklungsumgebungen, kein allgemeiner Project-Generator ohne Docker-/Compose-Vertrag.
- Application-Code bleibt frei von konkreter externer I/O; I/O sitzt in Driven-Adaptern und wird über Ports erreicht.
- Generierte Dateien und User-Projektdateien sind sicherheitsrelevant: Managed Blocks, Backups, Two-Phase-Planung und Bestätigungen sind Produktverträge, keine Komfortdetails.
- CLI-Output und generierte Artefakte sind Englisch; normative Specs und Planning-Dokumente bleiben Deutsch.
- Release- und Distributionsänderungen betrachten [`ADR-0004`](../docs/plan/adr/0004-ci-system.md)/[`ADR-0007`](../docs/plan/adr/0007-distributionswege-ghcr.md), CI-Gates sowie README und CHANGELOG zusammen.

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten); die Rolle bestimmt [`roles.md`](roles.md).
3. Betroffene IDs identifizieren (`LH-*`, `ADR-*`, Slice).
4. Kleinste Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt; bei Generator-Änderungen Replay-Evidence nach [`replay.md`](replay.md), bei Slice-Closure Verification-Evidence nach [`verification.md`](verification.md).
8. Ausgeführte Sensors und verbleibende Risiken berichten.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe §Guides; Findings nach
[`review.md`](review.md)) → Verifier. Kein Self-Review — anderer Kontext findet
andere Findings, derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).

## Leseordnung

1. [`AGENTS.md`](../AGENTS.md) — Hard Rules und Workflow
2. [`spec/lastenheft.md`](../spec/lastenheft.md) — was das Produkt zusagt
3. [`conventions.md`](conventions.md) — repo-lokale Strukturregeln und Adaptionen, bei Bedarf
