# AGENTS.md - Briefing fuer AI-Coding-Agenten

Diese Datei ist das verpflichtende Onboarding fuer jede AI-Session, die
in diesem Repo Code oder Dokumentation aendert. Sie traegt Hard Rules
und Pointer auf kanonische Quellen; sie ersetzt keine Spec, ADR oder
Slice-Doku.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gewinnt
die kanonische Quelle. Dann ist diese Datei anzupassen.

## Source Precedence

In dieser Reihenfolge lesen und aufloesen:

1. [`spec/lastenheft.md`](spec/lastenheft.md) - normative Anforderungen,
   Akzeptanzkriterien, Exit-Code- und Sprachvertraege.
2. [`spec/architecture.md`](spec/architecture.md) - Schichten,
   Komponenten, Importregeln, Podman-/Docker-Annahmen.
3. [`docs/plan/adr/`](docs/plan/adr/) - Architekturentscheidungen.
4. Aktiver Slice in [`docs/plan/planning/in-progress/`](docs/plan/planning/in-progress/)
   oder [`docs/plan/planning/next/`](docs/plan/planning/next/) - konkrete
   Arbeit, Tranchen, DoD und Closure-Bedingungen.
5. Ausfuehrbare Harness-Vertraege: [`Makefile`](Makefile),
   [`Dockerfile`](Dockerfile), [`.golangci.yml`](.golangci.yml) und
   [`.github/workflows/`](.github/workflows/).
6. Nutzer- und Quality-Doku unter [`docs/user/`](docs/user/), besonders
   [`docs/user/quality.md`](docs/user/quality.md).
7. [`README.md`](README.md), [`README.de.md`](README.de.md) und
   [`CHANGELOG.md`](CHANGELOG.md).
8. [`harness/README.md`](harness/README.md) und diese Datei.

## Hard Rules

### Role Separation

Rollen sind Kontextgrenzen. Nutze
[`harness/roles.md`](harness/roles.md) fuer Planner-, Architect-,
Implementation-, Reviewer-, Verifier- und Validator-Vertraege.

Wer geplant oder implementiert hat, reviewt oder verifiziert nicht mit
demselben Eingabe-Kontext. Jeder Rollenwechsel braucht ein
Uebergabe-Artefakt: Plan, ADR-Bezug, Diff, Findings,
Verification-Evidence, Validation-Evidence oder Closure-Notiz.

### Docker-only Workflow

`u-boot` hat keinen Host-Go-Toolchain-Vertrag. Build, Lint, Tests,
Coverage und Security-Gates laufen ueber `make`, das Docker nutzt.
Host-Voraussetzungen sind Docker und GNU `make`.

Verwende fuer Verifikation bevorzugt:

- `make lint`
- `make test`
- `make coverage-gate`
- `make docs-check`
- `make gates`
- `make ci`
- `make fullbuild`

Lokale Host-Toolchain-Befehle duerfen nicht als alleiniger Nachweis
fuer einen Handoff dienen.

### Hexagonale Architektur

Die Import- und Verantwortungsregeln aus
[`spec/architecture.md`](spec/architecture.md) sind verbindlich.
Insbesondere:

- `hexagon/application` kennt keine konkreten Adapter.
- Ports bleiben kreuz-blind (`driving` importiert nicht `driven` und
  umgekehrt).
- Konkrete Adapter werden nur im Wiring unter `cmd/uboot` verbunden.
- Docker-/Compose-Zugriffe laufen ueber Ports und Adapter, nicht direkt
  aus Application-Code.

### Spec-Traceability

Code-, Test- und Doku-Aenderungen muessen die betroffenen `LH-*`,
`ADR-*` oder Slice-IDs kennen. Neue oeffentliche CLI-Vertraege brauchen
mindestens einen Spec- oder ADR-Anker und einen Test- oder Gate-Nachweis.

CLI-Ausgaben, Fehlermeldungen und generierte Dateien bleiben Englisch
(`LH-LESE-002`), auch wenn Plan- und Spec-Dokumente deutsch sind.

### Exit-Code-Vertraege

Die `LH-FA-CLI-006`-Klassifikation ist ein Produktvertrag. Neue
Subcommands muessen ihre Fehlerpfade auf die bestehenden Exit-Code-
Kategorien abbilden und Tests fuer relevante Sentinels pinnen.

### Managed-Block- und Dateisicherheit

Generatoren und Re-Init-Pfade duerfen User-Dateien nicht opportunistisch
ueberschreiben. Nutze die vorhandenen managed-block-, Plan-and-Execute-,
Backup- und Two-Phase-Patterns. Destruktive Operationen brauchen die im
Spec/Slice verlangte Bestaetigungslogik.

### Suppression-Disziplin

Inline-`//nolint` ist nicht die Standardloesung. Dauerhafte Ausnahmen
leben zentral in [`.golangci.yml`](.golangci.yml) mit Begruendung und,
wenn temporaer, zusaetzlich im Carveout-Inventar.

### ADR-Disziplin

Accepted ADRs werden nicht inhaltlich umgeschrieben. Neue oder
geaenderte Entscheidungen entstehen als neue ADR oder als klar
dokumentierte Folgeentscheidung mit Verweis auf die alte ADR.

Gates, Coverage-Schwellen und Architekturregeln duerfen nicht still
gelockert werden. Eine Abschwaechung braucht einen expliziten Plan- oder
ADR-Anker; temporaere Ausnahmen brauchen einen Carveout.

### Planning-Lifecycle

Planning-Artefakte folgen:

```text
open/ -> next/ -> in-progress/ -> done/
```

Lifecycle-Bewegungen erfolgen per `git mv`, damit Historie erhalten
bleibt. Substanzielle Aenderungen an `done/`-Artefakten erzeugen einen
neuen Slice statt die alte Closure umzuschreiben.

Jeder neue temporaere Carveout bekommt parallel einen Eintrag in
[`docs/plan/planning/in-progress/carveouts.md`](docs/plan/planning/in-progress/carveouts.md)
und einen Plan-Anker.

## Quality Gates

Nur reale Make-Targets zaehlen als Harness-Sensoren:

| Target | Zweck |
| --- | --- |
| `make lint` | `golangci-lint` mit SOLID-nahem Profil und depguard |
| `make test` | `go test ./...` im Docker-Test-Stage |
| `make test-docker` | Docker-tag Integrationstests gegen echte Docker Engine |
| `make coverage-gate` | Coverage-Schwelle, Default 90 Prozent |
| `make docs-check` | relative Markdown-Links in `docs/`, `spec/`, README-Dateien |
| `make govulncheck` | Go-Vulnerability-Scan |
| `make image-scan` | Trivy HIGH/CRITICAL gegen Runtime-Image |
| `make verify-depguard` | On-demand Nachweis, dass depguard-Regeln feuern |
| `make gates` | `lint` + `test` + `coverage-gate` + `docs-check` |
| `make ci` | `gates` + `govulncheck` + `image-scan` |
| `make fullbuild` | `ci` + Runtime-Image-Build |

Vor Handoff mindestens den engsten sinnvollen Sensor ausfuehren. Fuer
Codeaenderungen ist `make gates` der normale Abschluss, sofern die
Umgebung Docker zulaesst.

## Minimal Agent Workflow

1. [`harness/README.md`](harness/README.md) lesen.
2. Rolle aus [`harness/roles.md`](harness/roles.md) bestimmen.
3. Source Precedence anwenden und die relevante Spec/ADR/Slice-Doku
   lesen.
4. Betroffene `LH-*`, `ADR-*` und Slice-IDs benennen.
5. Kleinste sinnvolle Aenderung umsetzen.
6. Engsten passenden Sensor laufen lassen; bei Codeaenderungen nach
   Moeglichkeit `make gates`.
7. Oeffentliche Vertraege in README, `docs/user/`, ADR-Index, Roadmap,
   Slice oder CHANGELOG nachziehen.
8. Im Handoff ausgefuehrte Sensoren, nicht ausgefuehrte Sensoren und
   verbleibende Risiken klar nennen.
