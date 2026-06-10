# ADR 0003: SOLID-nahes Lint-Profil

## Status

Accepted

## Datum

2026-05-21

## Kontext

[`LH-QA-004`](../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil) verlangt Linting für Quellcode und Konfigurationsdateien.
Die Bootstrap-`.golangci.yml` aus M1 aktiviert nur die 5 Default-Linter
(`govet`, `errcheck`, `staticcheck`, `unused`, `ineffassign`) plus
`depguard` für die Schicht-Regeln aus [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement). Das ist die
Mindestmenge gegen offensichtliche Bug-Patterns, deckt aber
Designsignale (Komplexität, Funktionslänge, Interface-Pollution,
globale Kopplung, Context-Disziplin) nicht ab.

Beide Vorlagen-Projekte fahren ein erweitertes Profil:

- `k-deskflight` (ADR 0012 §2.2): 5 Defaults + 24 SOLID-nahe Linter,
  `//nolint`-Pragmas verboten, Carveouts zentral in `.golangci.yml`.
- `m-trace` (`spec/lastenheft.md §10.1`, `docs/user/quality.md §1.2`):
  identische 5 + 24 Aufteilung; Schwellen produktiv eingespielt
  (`cyclop=15`, `funlen=100/60`, `gocognit=20`, `gocyclo=15`,
  `interfacebloat=10`, `maintidx=20`, `nestif=5`).

Lastenheft-Bezug:

- [`LH-QA-004`](../../../spec/lastenheft.md#lh-qa-004-linting-solid-nahes-lint-profil) – Linting (in diesem Commit auf MVP-Pflicht gehoben,
  Spezifikation auf SOLID-nahes Profil verschärft).
- [`LH-FA-BUILD-001`](../../../spec/lastenheft.md#lh-fa-build-001-multi-stage-dockerfile-u-boot-repo) – `lint`-Stage des Dockerfile.
- [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) – `depguard` für Schicht-Regeln (bereits aktiv).
- [`LH-NFA-MAINT-001`](../../../spec/lastenheft.md#lh-nfa-maint-001-modulare-architektur)..[`LH-NFA-MAINT-003`](../../../spec/lastenheft.md#lh-nfa-maint-003-testbarkeit) – modulare Architektur, Erweiterbarkeit,
  Testbarkeit.

## Entscheidung

u-boot übernimmt das **SOLID-nahe Lint-Profil aus m-trace**
unverändert für den Go-Teil: 5 Default-Linter plus 24 SOLID-nahe
Linter (`depguard` ist als Schicht-Regel-Linter Teil dieser 24, siehe
[`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement)). Damit sind 29 Linter aktiv.

Konkrete Setzungen:

- Linter-Liste 1:1 wie in `docs/user/quality.md` §1.2 — der dortige
  Doku-Block ist Single Source of Truth.
- Schwellen wie m-trace `apps/api/.golangci.yml`:
  `cyclop.max-complexity=15`, `funlen.lines=100/statements=60`,
  `gocognit.min-complexity=20`, `gocyclo.min-complexity=15`,
  `interfacebloat.max=10`, `maintidx.under=20`,
  `nestif.min-complexity=5`, `dupl.threshold=150`.
- `forbidigo` verbietet `fmt.Print*` (Logging gehört in `log/slog`
  bzw. zukünftig in einen dedizierten Logging-Port).
- `gomodguard_v2` aktiviert mit leerem Block-/Allow-Set; konkrete
  Regeln folgen, sobald externe Modul-Dependencies dazukommen. Der
  v1-`gomodguard` ist seit `golangci-lint` v2.12.0 deprecated und wird
  in u-boot nicht mehr verwendet.
- `//nolint`-Pragmas sind verboten (analog [`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003-import-regeln-und-enforcement) für
  `depguard`). Pro-Pfad-Carveouts in `issues.exclude-rules` mit
  `Why:`-Kommentar.
- Aktive MVP-Carveouts (Detail in `docs/user/quality.md` §1.3):
  - `*_test.go`: `cyclop`, `gocognit`, `gocyclo`, `nestif`, `funlen`,
    `noctx`, `unparam`, `revive.unused-parameter` deaktiviert
    (tabellengetriebene Tests, Fakes).
  - `cmd/uboot/**`: `testpackage` und `gochecknoglobals` deaktiviert
    (Wiring-Schicht braucht `main`-Package; `var version` wird per
    `-ldflags` überschrieben).
- `errcheck` mit Ausnahme für `fmt.Fprintln`/`Fprintf`/`Fprint`
  (CLI-Writes auf stdout/stderr; bereits aus M1 übernommen).

## Konsequenzen

Positiv:

- **Designdrift sichtbar** ab Tag 1: zu lange Funktionen, hohe
  Komplexität, globale Variablen, fehlende Context-Disziplin werden
  PR-blockierend gemeldet, statt erst im Review aufzufallen.
- **Konsistenz** mit `k-deskflight` und `m-trace` — Reviewer und
  Beitragende kennen das Profil bereits, Schwellen sind kalibriert.
- **Interface-Disziplin** (`iface`, `interfacebloat`, `ireturn`) passt
  direkt zu [`LH-FA-ARCH-002`](../../../spec/lastenheft.md#lh-fa-arch-002-schichten-und-verzeichnislayout): kleine, fokussierte Ports.
- **Carveouts zentral** in `.golangci.yml`; kein `//nolint`-Streuen
  über den Code.

Negativ / Trade-offs:

- **Lint-Stage wird langsamer:** 5 → 30 Linter, plus depguard.
  Erfahrung aus m-trace und k-deskflight: bleibt unter ~30 s für die
  bootstrap-große Codebase, akzeptabel.
- **Einarbeitung:** Beitragende ohne SOLID-Profil-Erfahrung brauchen
  einmal die Tabelle aus `docs/user/quality.md` §1.2 plus die
  Carveout-Liste in §1.3. Lint-Fehlermeldungen sind selbsterklärend.
- **Schwellen-Diskussion:** wenn ein produktiver Use-Case echt 110
  Zeilen braucht, muss entweder die Funktion zerlegt werden (was
  meistens die richtige Antwort ist) oder die Schwelle gehoben werden
  (Routine, Commit-Body begründet).
- **`gochecknoglobals` für `var version`** würde brechen, deshalb der
  `cmd/uboot/`-Carveout. Alternative wäre eine Konstante mit Compiler-
  Override via `-X main.Version`, aber das ist genau das, was das
  `var version` schon tut.

Alternativen (verworfen):

- **Nur Defaults belassen:** widerspricht [`LH-NFA-MAINT-001`](../../../spec/lastenheft.md#lh-nfa-maint-001-modulare-architektur)..[`LH-NFA-MAINT-003`](../../../spec/lastenheft.md#lh-nfa-maint-003-testbarkeit),
  überlässt Designdrift dem Review.
- **Engerer Schwellen-Satz als m-trace** (z. B. `funlen.lines=80`):
  ohne Daten-Basis spekulativ; m-trace-Werte sind in einem produktiven
  Codebase kalibriert.
- **Volles `revive`-Custom-Profil** (analog k-deskflight ADR 0012):
  über den MVP hinaus; kann in einem späteren ADR ergänzt werden.

## Folgepunkte

- Mit dem ersten produktiven Inkrement (M3): Lint-Findings aus der
  Erweiterung systematisch abarbeiten und Schwellen ggf. nachziehen.
- `gomodguard`-Regeln definieren, sobald externe Modul-Dependencies
  (z. B. `spf13/cobra`, `gopkg.in/yaml.v3`) konkret werden.
- Erweiterung um `revive`-Custom-Rules in einem Folge-ADR, falls die
  default-Konfiguration zu schwach wird.
