# Maschinen-lesbare CLI-Ausgabe (`--json`, `--dry-run`, `--diff`) — u-boot

| Dokument         | CLI-JSON-Vertrag                                              |
| ---------------- | ------------------------------------------------------------- |
| Projektname      | `u-boot`                                                      |
| Bezug            | `LH-NFA-USE-004`, `LH-FA-CLI-007`, `LH-FA-CLI-008`, `LH-FA-CLI-006` in [`spec/lastenheft.md`](../../spec/lastenheft.md) |
| ADR              | [`docs/plan/adr/0010-kein-http-driving-adapter.md`](../plan/adr/0010-kein-http-driving-adapter.md) |
| Slice-Anker      | [`docs/plan/planning/in-progress/slice-v1-cli-json-dry-run.md`](../plan/planning/in-progress/slice-v1-cli-json-dry-run.md) (Cluster) + [`docs/plan/planning/done/slice-v1-cli-json-dry-run-doctor.md`](../plan/planning/done/slice-v1-cli-json-dry-run-doctor.md) (Doctor-Folge-Slice, done) |
| Status           | Entwurf 0.4.0                                                 |
| Datum            | 2026-06-04                                                    |

## Zweck

Vertragsdokument für die maschinen-lesbare CLI-Ausgabe von u-boot.
Spec-Pflichtkontrakte aus dem Lastenheft sind hier verbatim
zitiert; die zugehörige Code-Lokation im Repo, die Spec-konformen
Diagnostic-Codes und die Migrations-Reihenfolge der 10 Spec-Enum-
Subcommands sind als verbindliche Quellen dokumentiert.

Die Pflichtaussagen leben im Lastenheft (§1809-1853 für `LH-NFA-
USE-004`, §302-447 für `LH-FA-CLI-007`, §451-489 für `LH-FA-CLI-008`).
Dieses Dokument ist die Detail-Doku für CLI-Konsumenten und für
den Test-Helper [`internal/adapter/driving/cli/jsontestutil/`](../../internal/adapter/driving/cli/).

ADR-0010 §Folgepunkte Re-Eval-Trigger 2 verankert: **JSON-CLI ist
die kanonische Maschinen-Schnittstelle von u-boot**; HTTP-/gRPC-
/WebSocket-Adapter sind ausdrücklich gegen genau dieses Surface
abgewogen und verworfen.

---

## 1. Zwei Kontraktstufen — wann gilt was?

`LH-NFA-USE-004` (§1841-1842) trennt zwei Vertragsstufen, die
**beide** im selben Wire-Format (`cliJSONEnvelope`, siehe §4)
gerendert werden:

| Aufruf-Modus | Pflicht-Vertrag | Spec |
| --- | --- | --- |
| `u-boot <cmd> --json` (read-only oder ohne Dry-Run/Diff) | **Minimalkontrakt** (§2) | `LH-NFA-USE-004` §1841 |
| `u-boot <cmd> --dry-run --json` | **Voll-Schema** (§3) | `LH-FA-CLI-007` §1842 |
| `u-boot <cmd> --diff --json` (mit oder ohne `--dry-run`) | **Voll-Schema** (§3) | `LH-FA-CLI-008` §468 |

Voll-Schema ist eine **Obermenge** des Minimalkontrakts: zusätzlich
zu den Minimal-Pflichtfeldern werden `dryRun`, `diff`, `plannedFiles`,
`changes` zu Pflichtfeldern. Im Minimalkontrakt sind diese vier
Felder **nicht** zulässig (sonst wäre es kein Minimalkontrakt mehr);
der Test-Helper `AssertMinimalEnvelope` rejected sie aktiv.

---

## 2. Minimalkontrakt (`LH-NFA-USE-004` §1823-1842)

Verbatim-Zitat des Lastenhefts:

> Für alle `--json`-Ausgaben gilt ergänzend ein gemeinsames
> Minimalkontrakt-Schema:
>
> - `status` (`ok`/`warn`/`error`)
> - `command` (Hauptbefehl als in `LH-FA-CLI-007` definiertes Enum)
> - optional `subcommand` (für gruppierte Befehle wie `template`
>   oder `config`)
> - `diagnostics` (Liste von Objekten mit mind. `level`, `code`,
>   `message`, optional `file`)
> - `exitCode` (vgl. `LH-FA-CLI-006`)
>
> Für `--json`-Antworten gilt zusätzlich:
>
> - `diagnostics`, wenn leer, darf als `[]` ausgegeben werden.
> - `diagnostics.level` darf nur `warn` oder `error` enthalten.
> - `diagnostics.code` folgt der Konvention: LH-Kennung der
>   verursachenden Anforderung (z. B. `LH-FA-DEV-003`). Tool-
>   interne Codes ohne LH-Bezug dürfen nur dann verwendet werden,
>   wenn ihre Bedeutung in der Dokumentation festgehalten ist
>   (Verweis: `LH-FA-CLI-007`).
> - `diagnostics.file` ist optional.
> - `status` ist an den höchsten in `diagnostics` enthaltenen
>   `level` gekoppelt: `error` → `status == "error"`; `warn` ohne
>   `error` → `status == "warn"`; sonst `status == "ok"`.
> - Bei `command == "template"` oder `command == "config"` ist
>   `subcommand` verpflichtend.
> - Die Felder `status`, `command`, `diagnostics` und `exitCode`
>   sind minimal verpflichtend und sollten mit anderen Feldern
>   ergänzt werden.

Spec-Beispiel für `u-boot doctor --json` im All-OK-Fall:

```json
{
  "status": "ok",
  "command": "doctor",
  "diagnostics": [],
  "exitCode": 0
}
```

### 2.1 Konsequenz: `level: "ok"` und `level: "info"` sind nicht zulässig

Das Lastenheft beschränkt `diagnostics[].level` strikt auf
`warn` oder `error`. Doctor-Checks mit interner `SeverityOK`-
oder `SeverityInfo`-Klassifikation werden im JSON-Modus
**übersprungen** — sie tauchen nicht als `diagnostics[]`-Eintrag
auf. Der All-OK-Fall serialisiert ein leeres Array (`diagnostics: []`).

### 2.2 Konsequenz: `--quiet` ist im JSON-Modus ein No-op

`--quiet` filtert im Plaintext-Output die OK-Items. Im JSON-Modus
sind OK-/Info-Items per §2.1 bereits ausgeschlossen — `--quiet`
hat keinen zusätzlichen Filter-Effekt. `u-boot doctor --quiet --json`
liefert einen **semantisch identischen** Envelope wie
`u-boot doctor --json` (gleiche `status`/`exitCode`, gleiche
`diagnostics`-Reihenfolge mit gleichen `code`/`level`-Paaren).

### 2.3 Exit-Codes (`LH-FA-CLI-006`)

`exitCode` im Envelope spiegelt den Prozess-Exit-Code:

- `0` — `status ∈ {ok, warn}` ohne `--strict`
- `2` — CLI-Validierungsfehler (unzulässige Argument-Kombination,
  noch nicht implementierte Subcommand-Forms, siehe §6)
- `11` — Doctor-Report enthält Errors, oder Warns unter `--strict`
  (`LH-FA-DIAG-003`)
- `10` — destruktive Operation ohne Freigabe, projekt-spezifisch
  (siehe `LH-FA-CLI-005A`)

---

## 3. Voll-Schema (`LH-FA-CLI-007` §322-417)

Voll-Schema gilt für `--dry-run --json` und `--diff --json`. Verbatim-
Zitat des Lastenhefts:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["status", "command", "dryRun", "diff", "plannedFiles", "changes", "diagnostics", "exitCode"],
  "properties": {
    "subcommand": {
      "type": "string",
      "description": "Unterkommando bei gruppierten Hauptkommandos wie `template` oder `config`"
    },
    "status": {
      "type": "string",
      "enum": ["ok", "warn", "error"]
    },
    "command": {
      "type": "string",
      "enum": ["init", "add", "remove", "up", "down", "doctor", "logs", "generate", "config", "template"]
    },
    "dryRun": {
      "type": "boolean"
    },
    "diff": {
      "type": "boolean"
    },
    "plannedFiles": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "action"],
        "properties": {
          "path": { "type": "string" },
          "action": {
            "type": "string",
            "enum": ["create", "modify", "delete"]
          }
        },
        "additionalProperties": true
      }
    },
    "changes": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "count"],
        "properties": {
          "path": { "type": "string" },
          "count": { "type": "integer", "minimum": 0 }
        },
        "additionalProperties": true
      }
    },
    "diagnostics": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["level", "code", "message"],
        "properties": {
          "level": { "type": "string", "enum": ["warn", "error"] },
          "code": { "type": "string" },
          "message": { "type": "string" },
          "file": { "type": "string" }
        },
        "additionalProperties": true
      }
    },
    "exitCode": {
      "type": "integer",
      "minimum": 0
    }
  },
  "allOf": [
    {
      "if": {
        "properties": {
          "command": { "const": "template" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    },
    {
      "if": {
        "properties": {
          "command": { "const": "config" }
        },
        "required": ["command"]
      },
      "then": {
        "required": ["subcommand"]
      }
    }
  ],
  "additionalProperties": true
}
```

Spec-Beispielinstanz für `u-boot add postgres --dry-run --json`:

```json
{
  "status": "warn",
  "command": "add",
  "dryRun": true,
  "diff": false,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "create" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 12 }
  ],
  "diagnostics": [
    { "level": "warn", "code": "LH-FA-CLI-007", "message": "Geplante Datei fehlt bereits" }
  ],
  "exitCode": 0
}
```

Spec-Beispielinstanz für `u-boot add postgres --diff --json` (Vorschau
mit anschließendem Schreiben):

```json
{
  "status": "ok",
  "command": "add",
  "dryRun": false,
  "diff": true,
  "plannedFiles": [
    { "path": "compose.yaml", "action": "modify" }
  ],
  "changes": [
    { "path": "compose.yaml", "count": 6 }
  ],
  "diagnostics": [],
  "exitCode": 0
}
```

---

## 4. Wire-Type und Lokation

Der Go-Wire-Type lebt im CLI-Adapter:

- **Envelope-Type**: `cliJSONEnvelope` in
  [`internal/adapter/driving/cli/jsonenvelope.go`](../../internal/adapter/driving/cli/)
  (kommt im Doctor-Folge-Slice, siehe Slice-Anker oben). Ein Typ
  mit beiden Modi: Minimal-Felder pflicht, Voll-Felder per
  `*bool`/nil-Slice `omitempty`.
- **Konstruktoren**: `newMinimalEnvelope(...)` und
  `newFullEnvelope(...)` — der Konstruktor pinnt die Spec-
  Stufe. Read-only-Pfade rufen ausschließlich den Minimal-
  Konstruktor, modifying-Pfade ausschließlich den Voll-
  Konstruktor.
- **Test-Helper**: `internal/adapter/driving/cli/jsontestutil/`
  mit zwei Modi `AssertMinimalEnvelope(t, raw, opts ...)` und
  `AssertFullEnvelope(t, raw, opts ...)`. Helper prüfen Schema-
  Konformität als Go-Code (kein embedded JSON-Schema, keine
  zusätzliche Dependency). Options-Pattern:
  `WithCommand(string)`, `WithSubcommand(string)`,
  `WithExpectedCodes(...string)` (Subset-Pin, **keine**
  Allowlist-Erweiterung), `WithExitCode(int)`.

Architektur-Grenze: Envelope-Type und Helper leben im **CLI-Adapter**,
**nicht** in Domain oder Application (`LH-FA-ARCH-002`/`-003`).

---

## 5. Code-Registry für `diagnostics[].code`

Spec §1835 (`LH-NFA-USE-004`) erlaubt für `diagnostics[].code`
zwei Quellen: **(a)** LH-Kennung der verursachenden Anforderung
(z. B. `LH-FA-DEV-003`, `LH-FA-CLI-007`), oder **(b)** tool-interne
Codes, falls ihre Bedeutung **in der Dokumentation festgehalten**
ist. u-boot verwendet aktuell tool-interne Codes mit Dotted-
Notation (`docker.installed`, `uboot.yaml.valid` etc.); diese
Registry ist die **kanonische Doku-Sektion** für Spec §1835.

Source-of-Truth: `internal/adapter/driving/cli/jsontestutil/coderegistry.go`
(`DefaultAllowedCodes`-Map). Diese Tabelle ist die spec-pflichtige
Doku-Form derselben Map; der Test-Helper `AssertMinimalEnvelope`
rejected `diagnostics[].code`-Werte außerhalb. Drift-Schutz
(zwei aktive Gates plus Acceptance-Helper):

- **Gate 1** (Code ↔ Map): Unit-Test prüft, dass jede
  `checkID*`-Konstante aus [`internal/hexagon/application/doctor.go`](../../internal/hexagon/application/doctor.go)
  einen Map-Eintrag hat.
- **Gate 2** (Map ↔ Doku): Unit-Test parst diese Tabelle und
  vergleicht sie symmetrisch gegen die Map. Bricht in beide Drift-
  Richtungen.

### 5.1 Doctor-Checks

Die folgenden Codes sind die kanonische Quelle für Gate 2 (Map ↔
Markdown). Die Markdown-Tabelle ist zwischen den HTML-Markern
eingegrenzt — Folge-Slices fügen weitere Tabellen für ihre eigenen
Subcommand-Codes an, ohne den Doctor-Block zu berühren.

<!-- code-registry:start -->

| Code | Bedeutung |
| --- | --- |
| `fs.write-permissions` | Schreib-Permission im Working Directory |
| `git.installed` | Git-Binary verfügbar |
| `docker.installed` | Docker-Binary verfügbar |
| `docker.reachable` | Docker-Daemon erreichbar |
| `docker.compose.installed` | Compose-Plugin verfügbar |
| `uboot.yaml.valid` | `u-boot.yaml` syntaktisch valide |
| `compose.yaml.valid` | `compose.yaml` syntaktisch valide |
| `devcontainer.json.valid` | `.devcontainer/devcontainer.json` syntaktisch valide |
| `devcontainer.dockerfile.valid` | `.devcontainer/Dockerfile` parsebar |
| `services.enabled-key` | `u-boot.yaml` `services`-Block konsistent |
| `devcontainer.forwardPorts.consistency` | `devcontainer.json` `forwardPorts` konsistent |
| `devcontainer.features.allowlist` | `devcontainer` Features auf Allowlist |
| `devcontainer.features.drift` | `devcontainer` Features ohne Drift |

<!-- code-registry:end -->

Weitere Subcommand-Sektionen kommen mit den jeweiligen Folge-Slices
(siehe §6).

---

## 6. Per-Command-Migrations-Reihenfolge

Spec-Enum (`LH-FA-CLI-007` §338) listet zehn Subcommands; alle
sollen `--json` tragen (`LH-NFA-USE-004` §1813). Migration läuft
**inkrementell** über neun Folge-Slices unterhalb des Cluster-
Slices [`slice-v1-cli-json-dry-run`](../plan/planning/in-progress/slice-v1-cli-json-dry-run.md);
Reihenfolge gemäß Cluster-T0-(e):

| # | Folge-Slice | Subcommand-Form(en) | Status |
| --- | --- | --- | --- |
| 1 | `slice-v1-cli-json-dry-run-doctor` | `doctor` + `template list`-Flag-Schnitt | done |
| 2 | `slice-v1-cli-json-dry-run-add` | `add` (modifying, etabliert Full-Modus + RecordingFileSystem) | done (DoD-Hash-Tabelle im Slice-File) |
| 3 | `slice-v1-cli-json-dry-run-init` | `init` | done (DoD-Hash-Tabelle im Slice-File) |
| 4 | `slice-v1-cli-json-dry-run-generate` | `generate` | done (DoD-Hash-Tabelle im Slice-File) |
| 5 | `slice-v1-cli-json-dry-run-remove` | `remove` | done (DoD-Hash-Tabelle im Slice-File) |
| 6 | `slice-v1-cli-json-dry-run-up-down` | `up`, `down` (gebündelt, read-only Compose-Status) | done (DoD-Hash-Tabelle im Slice-File) |
| 7 | `slice-v1-cli-json-dry-run-logs` | `logs` (Sub-Decision: JSON-Lines vs. Single-Envelope) | offen |
| 8 | `slice-v1-cli-json-dry-run-config` | `config`, `config get`, `config set` (gebündelt, drei Formen) | offen |
| 9 | `slice-v1-cli-json-dry-run-template` | `template` (bare) + `template list`-Envelope-Migration | offen |

### 6.1 Übergangs-Reject für nicht-migrierte Forms

Bis der jeweilige Folge-Slice landet, rejected `u-boot <sub> --json`
für die nicht-migrierten Subcommand-Formen mit Exit-Code `2`
(`LH-FA-CLI-006`-Klasse) und Fehlermeldung:

```
JSON-Ausgabe für 'u-boot <sub>' ist noch nicht implementiert
(siehe slice-v1-cli-json-dry-run-<sub>).
```

Mechanik: zentrale Allowlist-Map in
[`internal/adapter/driving/cli/root.go`](../../internal/adapter/driving/cli/)
plus `PersistentPreRunE` am Root-Command. Map-Key ist Cobra-
`cmd.CommandPath()` (`"u-boot doctor"`, `"u-boot config get"`,
…). Pro Folge-Slice-Merge wandern die Allowlist-Einträge der
migrierten Formen rein; Cluster-T_close entfernt Allowlist und
`PersistentPreRunE` komplett.

### 6.2 Sonderfall `template list --json`

`template list --json` existiert heute mit eigenem lokalem
`--json`-Flag und einem subcommand-spezifischen Array-Output
(`templateJSON`). Der Doctor-Folge-Slice führt **nur** den
Flag-Schnitt durch (lokales Flag → Root-Flag), das Output-Format
bleibt **unverändert** und ist deshalb **nicht** Minimalkontrakt-
konform. Carveouts-Eintrag in
[`docs/plan/planning/in-progress/carveouts.md`](../plan/planning/in-progress/carveouts.md)
§Temporäre Carveouts mit Re-Trigger
`slice-v1-cli-json-dry-run-template` (Platz 9 — Envelope-Migration).

### 6.3 `u-boot add --json` (slice-v1-cli-json-dry-run-add, done)

`u-boot add` ist der erste modifying-Subcommand mit JSON-
Envelope-Migration. Drei Flag-Kombinationen, drei Output-Formen
(`LH-FA-CLI-007/008`):

- **`--json` ohne `--dry-run`/`--diff`** → Minimalkontrakt
  (Spec §1841). Die Operation schreibt das FS um, das JSON-Output
  trägt aber **keine** Plan- oder Change-Information. Für die
  Liste der veränderten Files den Preview-Pfad nutzen:
  `--dry-run --json` (Vorschau ohne Schreiben) oder
  `--diff --json` (Vorschau plus Schreiben, mit Hunks).
- **`--dry-run --json`** → Voll-Schema mit
  `plannedFiles[]`/`changes[]`, `dryRun: true`, `diff: false`. Es
  wird **nichts** auf die Disk geschrieben — der
  `RecordingFileSystem` capturet alle geplanten Mutations.
- **`--diff --json`** → Voll-Schema mit
  `plannedFiles[].hunks[]`, `dryRun: false`, `diff: true`. Es wird
  geschrieben **und** capturet (Preview-and-Apply, Spec §465-470).
- **`--dry-run --diff --json`** → wie `--dry-run --json`, aber mit
  Hunks und `diff: true`. Kein Write.

**Mid-Write-Failure-UX** (`LH-NFA-REL-003`): wenn im
`--diff --json`-Preview-and-Apply-Pfad ein Write mid-stream failt,
trägt `plannedFiles[]` nur die Aufrufe bis zur Failure-Stelle. Der
Diagnostic-Eintrag (`diagnostics[].code: "LH-NFA-REL-003"`) trägt
`file:` mit der Failure-Stelle; `exitCode: 14`. Die später nicht
mehr geschriebenen Files erscheinen **nicht** in `plannedFiles[]`
— ein Roll-back-aware Capture ist V1-Out-of-Scope (Cluster-T0-(b)
Variante 3 ChangeSet-Pattern explizit verworfen).

**Sub-Decisions**: der Slice etabliert `RecordingFileSystem` als
driven-Adapter (`internal/adapter/driven/recordingfs/`), `Hunks`-
Schema (`plannedFiles[].hunks: [{oldStart, oldLines, newStart,
newLines, content}]`, T0-(l)), Pure-Go LCS-Diff-Renderer
(`internal/adapter/driving/cli/diff/`, T0-(d)), Composition-Root-
`fsFactory(driving.AddPreviewMode)`-Closure in `cmd/uboot/main.go`
(T0-(e)) und `changes[].count`-Semantik gemäß Spec §477
(`CountAdditions` über die `+`-Lines der Hunks, T0-(g)).
Diagnostic-Codes sind LH-Kennungen
(`LH-FA-ADD-{001,002,005,006}`/`LH-FA-INIT-{004,005,006}`/
`LH-NFA-REL-003`); keine erfundenen `add.*`-Codes (T0-(j)).

### 6.4 `u-boot init --json` (slice-v1-cli-json-dry-run-init, done)

`u-boot init <name>` ist der zweite modifying-Subcommand und der
wichtigste Onboarding-Use-Case. Vier Flag-Kombinationen mit
exakter Symmetrie zu `add` (T0-(a) Pattern-Erbe-Disziplin):

- **`--json` ohne `--dry-run`/`--diff`** → Minimalkontrakt-
  Envelope. Operation schreibt das FS um (Skeleton-Dirs +
  Skeleton-Files + `u-boot.yaml`), das JSON-Output trägt aber
  **keine** Plan- oder Change-Information. Preview-Information
  liefert der Preview-Pfad: `--dry-run --json` oder
  `--diff --json`.
- **`--dry-run --json`** → Voll-Schema mit `plannedFiles[]`/
  `changes[]`, `dryRun: true`, `diff: false`. Es wird **nichts**
  auf die Disk geschrieben — der `RecordingFileSystem` capturet
  alle geplanten Mutations. Auch der separate
  `driven.GitClient`-Port wird in diesem Modus **geskippt**
  (`PreviewMode == PreviewDryRun` umgeht `s.initGit`, weil git am
  Recorder vorbei auf die echte Disk schreiben würde — T0-(n)).
- **`--diff --json`** → Voll-Schema mit
  `plannedFiles[].hunks[]`, `dryRun: false`, `diff: true`. Es
  wird geschrieben **und** capturet (Preview-and-Apply); `git
  init` läuft regulär.
- **`--dry-run --diff --json`** → wie `--dry-run --json`, aber
  mit Hunks und `diff: true`. Kein Write.

**Mid-Write-Failure-UX** (`LH-NFA-REL-003`): identisch zu `add`
— Voll-Schema-Envelope trägt `plannedFiles[]` bis zur Failure-
Stelle, `diagnostics[].file` markiert die Position,
`exitCode: 14`. Backup-Sentinels
(`ErrBackupSuffixExhausted`/`ErrBackupSourceMissing`) und der neue
`ErrInitFileSystem`-Sentinel klassifizieren als `LH-NFA-REL-003`
(technische Persistenz, T0-(f)). **Switch-Order-Pflicht**: das
init-CLI prüft `ErrInitFileSystem` als ersten `errors.Is`-Case,
weil Multi-`%w`-Wraps (Go 1.20+) einen FS-Fehler sonst auf einen
Exit-10-Fachfehler downgraden könnten.

**Planning-Phase-Failures** (T0-(q)): Fehler vor dem
Recorder-Capture (Force-without-Backup, ungültige Service-Namen,
Template-Read-Failures) produzieren `plannedFiles: []` mit
fachlichem Diagnostic-Code (z. B. `LH-FA-INIT-005` / Exit 10) —
**nicht** Exit 14. Die Unterscheidung zu Mid-Write-Failure ist
load-bearing: Planning-Errors sind User-Action-Klasse,
Write-Errors sind FS-Klasse.

**Template-Modus-Mutex** (T0-(i)): `init --template <name>`
kombiniert mit `--dry-run` oder `--diff` wird mit
`ErrTemplateConflictsWithFlag` (Exit-Code 2, `LH-FA-CLI-006`)
rejected. V1 liefert keine Template-Preview; die Migration läuft
über einen eigenen Folge-Slice (Cluster-Roadmap).

**ProgressPort-Silencing-Hint** (T0-(o)): Modifying-Subcommands
mit stdout-bound ProgressPorts MÜSSEN den Port im JSON-Mode
silencen — der Recorder schützt nur die FS-Layer, nicht stdout.
`init` setzt `req.SilenceProgress = flags.JSON` und swappt im
Use-Case auf einen Noop-Adapter; sonst landen
`progress.AffectedFiles`-Events vor dem JSON-Envelope und brechen
JSON-Parser-Konsumenten. Add hat diesen Port nicht; bei künftigen
modifying-Subcommands mit Progress-Events ist das Silencing
verbindlich.

**Context-Cancellation-Carveout** (T0-(p)): Ctrl-C oder
`context.Canceled` während eines modifying-Subcommands fällt
heute auf die `LH-FA-CLI-006`-Default-Klausel und Exit 2. Eine
Interrupt-aware Exit-130-Convention bleibt eigener
Cross-Cutting-Folge-Slice für **alle** modifying-Subcommands —
init ändert den Status-quo nicht.

**Path-Anchor-Vertrag** (T0-(k)): `plannedFiles[].path` ist
**project-relativ** (anchor = `<name>`/), unabhängig davon, wie
der positional `<name>` geschrieben wurde (trailing-slash,
dot-slash, absoluter Pfad). Acceptance-Pins decken alle vier
Varianten.

**Concurrency**: `InitProjectService` trägt einen eigenen
`initMu sync.Mutex` — der `s.fs`-Swap pro Request ist nicht
race-frei ohne Lock. Zwei Goroutinen, die parallel `init` auf
denselben Service mit unterschiedlichen TempDirs aufrufen, werden
serialisiert (Acceptance-Pin).

### 6.5 `u-boot generate --json` (slice-v1-cli-json-dry-run-generate, done)

`u-boot generate <artifact>` ist der **erste** Subcommand, der
mehrere Artefakte (`changelog` / `readme` / `env-example` /
`devcontainer`) über einen einzigen Subcommand bedient. Vier
Flag-Kombinationen mit derselben Symmetrie wie `add`/`init`:

- **`--json` ohne `--dry-run`/`--diff`** → Minimal+Data-Envelope
  via `newDataEnvelope`. Operation schreibt das FS um; das
  JSON-Output trägt **keinen** Plan-Inhalt, aber `data.artifact`
  und `data.action` für Konsumenten-Klassifikation.
- **`--dry-run --json`** → Voll-Schema mit `plannedFiles[]`/
  `changes[]` plus `data.action`, `dryRun: true`, `diff: false`.
  Es wird **nichts** auf die Disk geschrieben.
- **`--diff --json`** → Voll-Schema mit `plannedFiles[].hunks[]`,
  `dryRun: false`, `diff: true`, plus `data.action`. Es wird
  geschrieben **und** capturet.
- **`--dry-run --diff --json`** → wie `--dry-run --json`, aber
  mit Hunks. Kein Write.

**Envelope-Shape** (T0-(m)): `command="generate"`, **kein
`subcommand`-Feld** (Cobra-`<artifact>` ist Positional-Arg, kein
Subcommand). Artefakt wird in `data.artifact:
"<changelog|readme|env-example|devcontainer>"` getragen.

**Action-Klassifikation via `data.action`** (T0-(f)):
`<created|updated-block|no-op|repaired-manual>`. **NoOp**
produziert zusätzlich `plannedFiles: []` UND `changes: []`
(beide leer). **UpdatedBlock** und **RepairedManual** sind
FS-semantisch identisch (`plannedFiles[i].action: "modify"`);
`data.action` ist der **einzige** Discriminator zwischen ihnen
(Acceptance-Pin in `generate_acceptance_test.go`).

**Error-Envelope** trägt `data.artifact` aber **kein
`data.action`** (T0-(q)): die Use-Case-Response ist Zero auf
dem Error-Pfad, eine Action existiert nicht. Bei unbekanntem
Artefakt (`ErrArtifactUnknown` aus
`domain.NewArtifact`-Validation) entfällt `data` komplett —
das Artefakt ist nicht klassifizierbar.

**Per-Artefakt LH-Code-Tabelle** (T0-(e)) für
`ErrGenerateManualConflict`:

| Artefakt | LH-Code | Exit-Code |
| --- | --- | --- |
| `changelog` | `LH-FA-GEN-002` | 10 |
| `readme` | `LH-FA-GEN-003` | 10 |
| `env-example` | `LH-FA-GEN-004` | 10 |
| `devcontainer` | `LH-FA-DEV-001` | 10 |

Plus weitere Diagnostic-Codes: `ErrGenerateFileSystem` →
`LH-NFA-REL-003` / Exit 14 (FS-Klasse, Switch-Order-First);
`ErrConfigValueInvalid` (ungültige
`--allow-external-feature-sources`-URL) → `LH-FA-DEV-003` /
Exit 10 (Spec §720); `ErrArtifactUnknown` → `LH-FA-CLI-006` /
Exit 2 (CLI-Validation, Spec §1157); `ErrProjectNotInitialized`
→ `LH-FA-INIT-001` / Exit 10.

**`--allow-external-feature-sources`-Mutex-Check**: der Flag ist
NUR für `generate devcontainer` gültig. Auf anderen Artefakten
rejected die CLI mit `ErrArtifactUnknown` / Exit 2 **vor** dem
Use-Case-Call (Acceptance-Pin).

**Devcontainer-Atomicity-Asymmetrie** (T0-(i)): Phase 1
(`planDevcontainerFiles`) ist Pre-Write-Validation-atomar (kein
WriteFile bei `ErrGenerateManualConflict`); Phase 2
(`executeDevcontainerPlans`) ist **nicht** Roll-back-atomar —
Mid-Write zweiter File hinterlässt halbgeschriebenen Zustand.
Carveout dokumentiert in
[`carveouts.md`](../plan/planning/in-progress/carveouts.md)
§Temporäre Carveouts; Rollback-Slice
[`slice-v2-generate-devcontainer-rollback-aware-write`](../plan/planning/open/slice-v2-generate-devcontainer-rollback-aware-write.md)
on hold pending trigger.

**Concurrency**: `GenerateService.generateMu sync.Mutex`
serialisiert konkurrierende `Generate()`-Calls auf demselben
Service (analog `InitProjectService.initMu`).

### 6.6 `u-boot remove --json` (slice-v1-cli-json-dry-run-remove, done)

`u-boot remove <service>` ist der fünfte modifying-Subcommand und
die **inverse Operation zu `add`**: strip managed-block aus
`compose.yaml` + `.env.example`, flip `services.<name>.enabled`
auf `false` in `u-boot.yaml`, optional Volume-Purge via
`--purge`-Gate (`LH-FA-CLI-005A` §254). Acht Flag-Kombinationen
plus die orthogonale `--purge`-Dimension (T0-(h)):

- **`--json` ohne `--dry-run`/`--diff`** → Minimal+Data-Envelope
  via `newDataEnvelope`. Operation schreibt das FS um; das
  JSON-Output trägt **keinen** Plan-Inhalt, aber das Success-
  Quartet `data: {service, priorState, state, volumesPurged}`
  für Konsumenten-Klassifikation.
- **`--dry-run --json`** → Voll-Schema mit `plannedFiles[]`/
  `changes[]` plus dem `data`-Quartet, `dryRun: true`,
  `diff: false`. Es wird **nichts** auf die Disk geschrieben —
  `RecordingFileSystem` capturet alle geplanten Mutations,
  inklusive der `RemoveAll`-Captures als
  `plannedFiles[].action: "delete"`.
- **`--diff --json`** → Voll-Schema mit `plannedFiles[].hunks[]`,
  `dryRun: false`, `diff: true`, plus `data`-Quartet. Es wird
  geschrieben **und** capturet (Preview-and-Apply).
- **`--dry-run --diff --json`** → wie `--dry-run --json`, aber
  mit Hunks. Kein Write.

**Envelope-Shape** (T0-(f)): `command="remove"`, **kein
`subcommand`-Feld** (Cobra-`<service>` ist Positional-Arg).
Success-`data` ist typed `removeEnvelopeData` mit
**Pointer-Wrapping**:

```json
{
  "service": "postgres",
  "priorState": "active",
  "state": "deactivated",
  "volumesPurged": false
}
```

`PriorState`/`State`/`VolumesPurged` sind Pointer (`*string`/
`*bool`), damit `omitempty` Key-**Abwesenheit** statt nur
Zero-Value-Drop pinnen kann (Spec §1841). `VolumesPurged` MUSS
`*bool` weil `false` ein valider Success-Wert ist (v0.3.0
deferred-Volumes). **Error-Envelope** trägt nur
`data: {"service": "<…>"}` ohne PriorState/State/VolumesPurged
(Zero-Response auf Error-Pfad, analog generate). Pre-Service-
Validation-Pfade (NoPositionalArg, ConflictingModeFlags,
InvalidServiceName) übergeben `data: null` — kein Service-
Kontext existiert.

**Idempotenz-NoOp-Semantik**: nur `PriorState=Deactivated`
qualifiziert als NoOp. Single+Repeat-Call gegen bereits
disabled Service liefern `plannedFiles: []` UND `changes: []`,
`data.priorState=data.state="deactivated"`, `status: ok`,
Exit 0. `EnabledUnset` und `InconsistentBlock` sind
state-transitioning (`Changed!=nil`, voll-`plannedFiles[]`)
— `EnabledUnset` wird auf `Deactivated` normalisiert mit
`data.priorState: "enabled-unset"` und
`data.state: "deactivated"`.

**`--purge`-Gate-Verhalten im JSON-Mode** (T0-(j)):
Confirmer-Prompt-Silencing via Service-Field-Mutation —
`req.SilenceConfirmer = flags.JSON` swappt `s.confirmer` auf
`noopConfirmer` (M4 Confirmer-Slice) innerhalb der
`removeMu`-Lock-Region; defer-Restore beim Wrapper-Return.
Pattern ist **nicht** aus init's `ProgressPort`-Silencing
geerbt (init swappt `s.progress`, nicht `s.confirmer`), sondern
ein remove-spezifisches Neu-Pattern. Bei
`--purge --no-interactive --json` OHNE `--yes` returnt der Gate
`ErrConfirmationRequired` → `LH-FA-INIT-005`-Envelope mit
Exit 10.

**WARN-Migration in `diagnostics[]`** (T0-(g)): heutige
`printRemoveSummary`-stderr-WARNING bei
`--purge && !VolumesPurged` (Volume-Removal ist in v0.3.0
deferred) wandert im JSON-Mode in `diagnostics[]`-Eintrag mit
`code: "LH-FA-ADD-007"`, `level: "warn"`, plus
`data.volumesPurged: false`. **`LH-FA-ADD-007` Multi-Use**:
derselbe Code identifiziert die Spec-Anforderung "Service
entfernen" (§924-947) UND markiert die deferred-Volumes-WARN.
Konsumenten disambiguieren ausschließlich über
`(code, level)`-Tupel: ERROR-Pfad
`ErrServiceUnregistered` liefert `code: "LH-FA-ADD-007"`,
`level: "error"`, Exit 10; WARN-Pfad liefert `level: "warn"`,
Exit 0. **Dry-Run-WARN-Suppression**: in `PreviewDryRun`
skippt die Use-Case den `runPurgeGate` (T0-(h)(a)) und führt
keine Mutation aus — die WARN-Prosa wäre semantisch falsch
("ist-deferred" statt "würde-deferred"); Human-Mode
`printRemoveSummary` unterdrückt entsprechend den WARN-Block.
`PreviewAndApply` behält die WARN.

**Custom-`Args`-Validator** (R11/R12/R13-Pins):
`validateRemoveArgs(a *App)` ist eine Cobra-PositionalArgs-
Closure mit `*App`-Capture, die `cobra.ExactArgs(1)` ersetzt.
Drei Cases:

- `len(args) == 1` → ok, weiter zu RunE.
- `len(args) == 0` → emit `LH-FA-CLI-006`-Envelope auf stdout
  bei `--json`, dann return `ErrServiceNameMissing`. Exit 2.
- `len(args) > 1` → emit `LH-FA-CLI-006`-Envelope auf stdout
  bei `--json` mit dem Cobra-Roh-Error `"accepts 1 arg(s),
  received N"`, dann return den Cobra-Error.
  Exit 2 via `isUsageError`-`"accepts "`-Prefix-Match.

In allen drei Fällen liest der Validator
`cmd.Flags().GetBool("dry-run"/"diff")` zur Validator-Zeit
(Cobra hat die Subcommand-Local-Flags zu diesem Zeitpunkt
bereits geparst) — Voll-Schema-Envelope bei `--dry-run` ODER
`--diff` (Spec §1842), sonst Minimal-Schema.

**`baseDirSanitizedError`-Wrapper**: FS-Wraps in der Use-Case
der Form `fmt.Errorf("remove write %s: %w: %w", absPath,
ErrRemoveFileSystem, raw)` tunneln den absoluten Filesystem-
Pfad in `diagnostic.message`. `runRemove` wrappt den UC-Error
vor dem `reportError`-Call mit `sanitizeBaseDir(removeErr,
cwd)`. Regeln:

- `<baseDir>/foo` → `foo` (project-relative)
- bare `<baseDir>` → `.` (project-root reference, an
  Word-Boundaries via `replaceBareBaseDir` — robust gegen
  Substring-Kollisionen wie `<baseDir>-cache/lock`)
- leerer baseDir → unverändert (defensive identity)

`Error()` ersetzt nur die Text-Form; `errors.Is`/`As` bleiben
intakt via Unwrap-Chain. Mapper-Switch-Order, Sentinel-
Identity, Exit-Code-Mapping — alles unverändert.

**Remove-LH-Code-Mapper-Tabelle** (FS-first Switch-Order,
Multi-`%w`-Defense):

| Sentinel | LH-Code | Exit |
| --- | --- | --- |
| `ErrRemoveFileSystem` | `LH-NFA-REL-003` | 14 |
| `ErrConfirmerUnavailable` | `LH-FA-CLI-005A` | 10 |
| `ErrConfirmationRequired` | `LH-FA-INIT-005` | 10 |
| `ErrServiceUnsupported` | `LH-FA-ADD-002` | 10 |
| `ErrServiceUnregistered` (ERROR) | `LH-FA-ADD-007` | 10 |
| `ErrServiceInconsistent` | `LH-FA-ADD-005` | 10 |
| `ErrProjectNotInitialized` | `LH-FA-ADD-001` | 10 |
| `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 |
| `ErrConflictingModeFlags` | `LH-FA-CLI-005A` | 2 |
| `cli.ErrServiceNameMissing` | `LH-FA-CLI-006` | 2 |
| Default (unknown) | `LH-FA-CLI-006` | 1 |

Plus die WARN-Multi-Use von `LH-FA-ADD-007` (siehe oben) — der
Code identifiziert die Spec-Anforderung, nicht die Sub-Semantik;
Disambiguation via `(code, level)`-Tupel.

**`delete`-Action im Recorder** (T0-(p)): `RemoveAll`-Captures
für `extraFiles` aus dem Catalog werden im Recorder als
`plannedFiles[].action: "delete"` gewired. `--diff` rendert
für `delete`-Aktionen Old-Content (der heute existierende
Datei-Inhalt) plus leeren New-Content; die Unified-Diff zeigt
nur `-`-Lines.

**Concurrency**: `RemoveServiceService.removeMu sync.Mutex`
serialisiert konkurrierende `Remove()`-Calls auf demselben
Service (analog `InitProjectService.initMu` /
`GenerateService.generateMu`). Confirmer-Swap UND fs-Swap
laufen **innerhalb** der Lock-Region — sonst könnten parallel
laufende Goroutinen ihre Swaps gegenseitig korrumpieren.

### 6.7 `u-boot up --json` / `u-boot down --json` (slice-v1-cli-json-dry-run-up-down, done)

`u-boot up` und `u-boot down` sind im sechsten Folge-Slice
gebündelt weil beide Subcommands den Compose-Status lesen und
das Confirmer-Swap-Pattern teilen. **Read-only-Klasse** auf
lokalem FS: weder `--dry-run` noch `--diff` (Cluster-Slice
Z. 464-467) — nur `--json` mit typisierten Data-Carriern.
Docker-Daemon-State ändern beide (Container starten / stoppen,
optional Volumes entfernen), aber das ist keine `--dry-run`-
fähige Surface.

- **`u-boot up --json`** → Minimal+Data-Envelope mit
  `data: {services: [{name, state, port, healthcheck}],
  timeoutFireAndForget?: bool}`. Stabilisierungs-Polling läuft
  wie im Human-Mode; nur die Compose-Progress-Stream-Ausgabe
  auf stderr wird unterdrückt
  (`req.SilenceProgress = flags.JSON` triggert
  Application-Layer-Branch auf `io.Discard`).
- **`u-boot down --json`** → Minimal+Data-Envelope mit
  `data: {removedVolumes: bool}`. Bool ohne `omitempty`:
  `false` ist der legitime Success-Wert "kein `--volumes`
  gesetzt".

**Envelope-Shape**: `command="up"` bzw. `command="down"`, KEIN
`subcommand`-Feld (beide sind Top-Level-Subcommands ohne Sub-
Form). KEIN `dryRun`/`diff`/`plannedFiles`/`changes`/`hunks`-
Feld (read-only-Klasse).

**`--quiet --json` semantisch identisch zu `--json`**:
Cluster-T0-(a) doctor-Pattern. `--quiet` darf den JSON-Output
NICHT unterdrücken — JSON ist die Maschinen-Schnittstelle.
Beide Flag-Reihenfolgen `--quiet --json` und `--json --quiet`
produzieren denselben Envelope.

**`--timeout=0` Fire-and-Forget im JSON-Mode**: Marker
`data.timeoutFireAndForget: true` erscheint ausschließlich im
`--timeout=0`-Pfad (`*bool`-Pointer mit `omitempty` —
Key-Absence-Disambiguation analog remove's `volumesPurged`).
`data.services: []` (NICHT `null`) auch bei Fire-and-Forget;
nil-Slice wird im CLI-Layer mit `[]serviceStatus{}`
initialisiert.

**`down --volumes` Confirmer-Branch im JSON-Mode** (T0-(d)
Option (b)): Request-time Gate-Branch in `runConfirmationGate`
Row 4 swappt den wired Confirmer auf einen lokalen
`noopConfirmer{}` (`(false, nil)`-Returns) wenn
`req.SilenceConfirmer == true`. **Refuse-by-Default-Semantik**:
bei `--volumes --json` OHNE `--yes` MUSS der User `--yes`
explizit setzen — Symmetrie zum `--no-interactive`-Pfad.
Direkter-Skip (proceed wie `AssumeYes`) wäre Security-by-
Default-Verletzung. Refuse-Pfad liefert
`ErrConfirmationRequired` → `LH-FA-INIT-005`/Exit 10 (geteilt
mit init/remove). Rows 1-3 (`!RemoveVolumes` / `AssumeYes` /
`NonInteractive`) behalten Vorrang vor dem JSON-Silencing-
Branch — explicit Flags > impliziter Mode-Default. Kein
State-Mutiert, kein neuer `downMu`-Mutex nötig, race-frei by
construction.

**Mapper-Tabelle mit verbindlicher Switch-Order** (T0-(e)
R3-HIGH-1) — Reihenfolge ist Switch-Sequenz im Mapper-Code:

| # | Sentinel | LH-Code | Exit | Mapper-Heim | Begründung |
| - | -------- | ------- | ---- | ----------- | ---------- |
| 1a | `driving.ErrUpFileSystem` | `LH-NFA-REL-003` | 14 | `mapUp` | FS-first damit Multi-`%w` mit FS+Docker auf FS-Klasse fällt |
| 1b | `driving.ErrDownFileSystem` | `LH-NFA-REL-003` | 14 | `mapDown` | analog 1a |
| 2 | `driven.ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | `helper` | Docker-Daemon vor Compose-Runtime |
| 3 | `driven.ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | `helper` | Compose-Runtime nach Daemon |
| 4 | `driving.ErrStabilizationTimeout` | `LH-FA-UP-001` | 12 | `mapUp` | Up-spezifische Runtime-Klasse |
| 5 | `driving.ErrConfirmationRequired` | `LH-FA-INIT-005` | 10 | `mapDown` | Confirmer-Refuse (geteilt mit init/remove) |
| 6 | `driving.ErrComposeFileMissing` | `LH-FA-UP-001` | 10 | `beide` | Fachliche Validierung (Datei-Schema) |
| 7 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | `beide` | Pattern-Erbe generate (Environment-Operation) |
| 8 | `cli.ErrInvalidTimeout` | `LH-FA-CLI-006` | 2 | `mapUp` | CLI-Form-Validierung |
| 9 | `cli.ErrConflictingModeFlags` | `LH-FA-CLI-005A` | 2 | `mapDown` | Mode-Mutex-Verträge |
| 10 | Default (unknown) | `LH-FA-CLI-006` | 1 | `beide` | Fallback |

**Cross-Slice-Klassen-Pin für `ErrProjectNotInitialized`**:
derselbe Sentinel mappt auf **`LH-FA-INIT-001`** bei
Environment-Subcommands (up/down/generate) UND auf
**`LH-FA-ADD-001`** bei Service-Subcommands (add/remove) —
bewusste Cluster-Konvention. Konsumenten dürfen NICHT erwarten
dass derselbe Sentinel cluster-weit auf denselben LH-Code
mappt; sie disambiguieren über `command` plus `code`.

**`(code, exitCode)`-Tupel-Disambiguation für Multi-`%w`-
Wraps**: Mapper-Tabelle ist FS-first (Row 1), aber der
ExitCode-Helper (`cli/cli.go:285-313`) checkt Driven-Sentinels
ZUERST. Das ergibt eine bewusste Zwei-Pfad-Klassifikation:

- Mapper → `diagnostics[0].code` (FS-Klasse-Signal)
- ExitCode-Helper → `exitCode` (Sub-Sentinel-Quelle)

Beispiele bei synthetisch konstruierten Multi-`%w`-Wraps:

| Multi-Wrap | `code` (Mapper) | `exitCode` (Helper) | Interpretation |
| --- | --- | --- | --- |
| `%w: %w` mit `ErrUpFileSystem` + `ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | FS-Klasse-Signal + Docker-Daemon-Sub |
| `%w: %w` mit `ErrUpFileSystem` + `ErrStabilizationTimeout` | `LH-NFA-REL-003` | 12 | FS-Klasse-Signal + Stabilization-Sub |
| `%w: %w` mit `ErrUpFileSystem` + `ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | FS-Klasse-Signal + Compose-Runtime-Sub |

Konsumenten MÜSSEN auf `(code, exitCode)`-Tupel filtern, NICHT
nur auf `code` allein. Pattern-Erbe remove's `LH-FA-ADD-007`-
Multi-Use ist die Disambiguation-Vorlage. Heute existiert kein
realer Code-Pfad der diese Sentinel-Paare chained
(`readComposeFile` failed VOR `ComposeUp`, `runConfirmationGate`
failed VOR `ComposeDown`) — die Disambiguation ist Defense-
only gegen künftige Multi-Wrap-Konstruktionen. Acceptance-
Tests pinnen ein Repräsentant pro Sub-Klasse (FS+Docker für
Exit 11, FS+StabilizationTimeout für Exit 12); FS+ConfirmRequired
/ FS+ComposeRuntime / FS+ProjectNotInitialized sind als
by-design-Carveouts dokumentiert.

**Path-Leak-Defense**: `runUp`/`runDown` wrappen UC-Errors mit
`sanitizeBaseDir(err, cwd)` vor `reportError` — der Sanitizer
extrahiert aus `cli/remove.go` (T7) in `cli/sanitize.go` (T5)
lebt jetzt package-intern. 11 FS-Read- und Compose-Runtime-
Wraps in upservice/downservice tunneln nach T3 keinen
absoluten Filesystem-Pfad mehr in `diagnostic.message`.

**Concurrency**: kein `upMu` / `downMu` — Pattern-Erbe init
T0-(d) (Service-Field-Swap mit Mutex) ist HIER explizit
**verworfen**. `SilenceProgress` und `SilenceConfirmer` sind
Bool-Felder im Request, die Branches sind lokale Variablen im
Use-Case-Method-Body. Kein Service-State mutiert → race-frei
by construction.

---

### 6.8 `u-boot logs --json` (slice-v1-cli-json-dry-run-logs, done)

`u-boot logs` ist der siebte Folge-Slice. **Read-only-Klasse** auf
lokalem FS (analog up-down): weder `--dry-run` noch `--diff` — nur
`--json` mit typisiertem Data-Carrier. Docker-Daemon-State unverändert
(`compose logs` ist read-only auf der Daemon-Seite).

- **`u-boot logs --json`** → Minimal+Data-Envelope mit
  `data: {lines: [string, ...]}`. **`lines` ohne `omitempty`**
  (Empty-Array-Pin: bei leerem Service-Set wird `"lines":[]`
  emittiert, NICHT `"lines":null` — Pattern-Erbe up-down's
  `services []serviceStatus`).

**T0-(a) Single-Envelope + `--follow --json` Reject** (Option (A)):
Spec-§1841-Konsens (Single-Envelope pro CLI-Call) wird honoriert.
`--follow` produziert konzeptionell einen Tail-Stream, nicht eine
beschränkte Antwort — und die NDJSON-Stream-Form ist Cluster-weit
nicht vorgesehen. Daher wird `--follow --json` in `runLogs` Stage-1
(VOR jedem UC-Call) mit `cli.ErrFollowJSONNotSupported` →
`LH-FA-CLI-006/Exit 2` rejected. **Bounded `--tail=N`** ist die
einzige `--json`-akquisitionsform.

**T0-(i) Validation-Reihenfolge**: bei `--follow --json --tail=-1`
ist `--follow --json` der dominante Reject-Pfad — die CLI-Stage-1-
Reihenfolge pinnt Reject-Sentinel VOR Tail-Validation (siehe
`TestLogsJSON_ValidationOrder_FollowJSONBeatsInvalidTail`). Konsumenten
können sich darauf verlassen, dass die kombinierte Fehler-Klasse
immer `LH-FA-CLI-006`/Exit 2 ist, ohne die Reject-Reihenfolge raten zu
müssen.

**Envelope-Shape**: `command="logs"`, KEIN `subcommand`-Feld. KEIN
`dryRun`/`diff`/`plannedFiles`/`changes`/`hunks`-Feld (read-only-Klasse).

**`--quiet --json` semantisch identisch zu `--json`**: Cluster-T0-(a)
doctor-Pattern. `--quiet` darf den JSON-Output NICHT unterdrücken.

**Mapper-Tabelle mit verbindlicher Switch-Order** (T0-(f)) —
Reihenfolge ist Switch-Sequenz im Mapper-Code (`mapLogsErrorToDiagnostic`):

| # | Sentinel | LH-Code | Exit | Begründung |
| - | -------- | ------- | ---- | ---------- |
| 1 | `driving.ErrLogsFileSystem` | `LH-NFA-REL-003` | 14 | FS-first damit Multi-`%w` mit FS+Docker auf FS-Klasse fällt |
| 2 | `driven.ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | helper `mapComposeRuntimeSentinel`, Docker-Daemon vor Compose-Runtime |
| 3 | `driven.ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | helper, Compose-Runtime nach Daemon |
| 4 | `driving.ErrComposeFileMissing` | `LH-FA-UP-001` | 10 | Fachliche Validierung (Datei-Schema) |
| 5 | `driving.ErrProjectNotInitialized` | `LH-FA-INIT-001` | 10 | Pattern-Erbe generate (Environment-Operation) |
| 6 | `domain.ErrInvalidServiceName` | `LH-FA-INIT-006` | 10 | Domain-level service-name Validierung |
| 7 | `cli.ErrFollowJSONNotSupported` | `LH-FA-CLI-006` | 2 | T0-(a) Reject-Pfad |
| 8 | `cli.ErrInvalidLogsTail` | `LH-FA-CLI-006` | 2 | CLI-Form-Validierung |
| 9 | Default (unknown) | `LH-FA-CLI-006` | 1 | Fallback |

**Cross-Slice-Klassen-Pin für `ErrProjectNotInitialized`**: derselbe
Sentinel mappt auf **`LH-FA-INIT-001`** bei Environment-Subcommands
(up/down/generate/logs) UND auf **`LH-FA-ADD-001`** bei Service-
Subcommands (add/remove) — bewusste Cluster-Konvention.

**`(code, exitCode)`-Tupel-Disambiguation für Multi-`%w`-Wraps**
(Pattern-Erbe up-down §6.7): Mapper ist FS-first (Row 1), aber der
ExitCode-Helper (`cli/cli.go`) checkt Driven-Sentinels ZUERST.
Beispiel-Tabelle:

| Multi-Wrap | `code` (Mapper) | `exitCode` (Helper) | Interpretation |
| --- | --- | --- | --- |
| `%w: %w` mit `ErrLogsFileSystem` + `ErrDockerUnavailable` | `LH-NFA-REL-003` | 11 | FS-Klasse-Signal + Docker-Daemon-Sub |
| `%w: %w` mit `ErrLogsFileSystem` + `ErrComposeRuntime` | `LH-NFA-REL-003` | 12 | FS-Klasse-Signal + Compose-Runtime-Sub |

Konsumenten MÜSSEN auf `(code, exitCode)`-Tupel filtern. Heute existiert
kein realer Code-Pfad der diese Sentinel-Paare chained (`Exists()`
failed VOR `ComposeLogs`) — Defense-only-Pin
(`TestLogsJSON_MultiWrap_FSAndDocker_SwitchOrderFSFirst_ByDesign`)
gegen künftige Multi-Wrap-Konstruktionen.

**Path-Leak-Defense**: `runLogs` wrappt UC-Errors mit
`sanitizeBaseDir(err, cwd)` vor `reportError` (Pattern-Erbe up-down).

**Bekannte Limitationen**:

- **`--follow --json` Reject**: konsequent inkompatibel, kein
  künftiger Streaming-Pfad geplant. Konsumenten die einen Tail-Stream
  brauchen müssen Human-Mode benutzen oder direkt `docker compose
  logs --follow` aufrufen.
- **`--tail=all`-Buffer-Caveat** (Pre-T6-Review MED-3): bei sehr großen
  Log-Volumina hält `--tail=all` die komplette Antwort vor dem JSON-
  Emit im Speicher. Konsumenten mit Mehr-GB-Logs sollten ein endliches
  `--tail=N` setzen. Eine Streaming-Form ist Cluster-weit nicht
  vorgesehen (siehe T0-(a) Reject-Entscheidung).
- **CRLF-Edge-Case** (Pre-T8-Bestätigungsrunde LOW-2): `splitLogLines`
  splittet nur auf `\n`. Bei einem hypothetischen CRLF-terminierten
  Compose-Output (Windows-Edge-Case) trügen alle `lines[i]` ein
  trailing `\r`. Heute kein bekannter Konsument-Pfad — Compose-CLI
  emittiert auf allen Plattformen `\n`. Falls je relevant: Folge-Slice
  ergänzt explizites `\r\n`-Handling.

**Concurrency**: kein `logsMu`. `JSON` und `Quiet` sind Bool-Felder im
`logsFlags`; der Reject-Pfad und der Acquisition-Pfad sind reine
Funktionsaufrufe ohne State-Mutation. Race-frei by construction.

---

## 7. Mutations-Matrix pro Subcommand

Drift-Anker für den Cluster-Folge-Slice-Block (T0-(f)): jeder
modifying-Subcommand listet, welche `driven.FileSystem`-Mutations-
Methoden sein Use-Case heute aufruft. Der `RecordingFileSystem`
implementiert alle 8 Mutations-Methoden (Drift-Schutz für künftige
Use-Cases); die Matrix dokumentiert die heutige Konsumenten-
Realität. Wenn ein zukünftiger Slice einen Use-Case erweitert
(z. B. neuer Add-on im Catalog ruft `Copy`), muss er die Matrix
in derselben PR ergänzen.

| Subcommand | WriteFile | WriteFileExclusive | Mkdir | MkdirAll | Rename | RemoveAll | Copy | CopyExclusive |
| --- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| `add` (slice-v1-cli-json-dry-run-add) | ✓ (3 Slots + N ExtraFiles) | — | — | ✓ (implizit, Recorder synthetisiert) | — | — | — | — |
| `init` (slice-v1-cli-json-dry-run-init) | ✓ (Skeleton-Files + `u-boot.yaml`) | — | ✓ (Backup-Verz.) | ✓ (Skeleton-Dirs direkt; Backup indirekt) | — | ✓ (Backup) | ✓ (Backup) | ✓ (Backup) |
| `generate` (slice-v1-cli-json-dry-run-generate) | ✓ (Artefakt-Files + `u-boot.yaml`-Allowlist-Mutation) | — | — | ✓ (devcontainer-Dir) | — | — | — | — |
| `remove` (slice-v1-cli-json-dry-run-remove) | ✓ (managed-block strip auf `compose.yaml` + `.env.example` + `enabled: false` auf `u-boot.yaml`) | — | — | — | — | ✓ (Catalog-`extraFiles`, optional) | — | — |
| `up` (slice-v1-cli-json-dry-run-up-down) | — | — | — | — | — | — | — | — |
| `down` (slice-v1-cli-json-dry-run-up-down) | — | — | — | — | — | — | — | — |
| `logs` (slice-v1-cli-json-dry-run-logs) | — | — | — | — | — | — | — | — |

Andere modifying-Subcommands (`config set`) ergänzen ihre Zeile
im jeweiligen Slice — heute `doctor`, `add`, `init`, `generate`,
`remove`, `up`, `down` und `logs` migriert. `up`/`down`/`logs` sind
read-only auf lokalem FS (nur `Exists`-Pre-Checks; bei `logs`
zusätzlich kein `ReadFile`, weil `compose logs` direkt gegen die
Docker-Daemon-Schicht läuft); die Docker-Daemon-Reads/-Mutationen
passieren durch den `DockerEngine`-Adapter ausserhalb dieser Matrix.

---

## 8. ADR-Bezug

[`ADR-0010 — kein HTTP-Driving-Adapter`](../plan/adr/0010-kein-http-driving-adapter.md)
§Folgepunkte Re-Eval-Trigger 2 macht JSON-CLI verbindlich zur
kanonischen Maschinen-Schnittstelle; HTTP-/gRPC-/WebSocket-
Adapter sind gegen genau dieses Surface abgewogen und verworfen.
Dieses Doku ist der Liefer-Anker für ADR-0010 §Trigger 2. ADR-0010
selbst bleibt **unverändert** (AGENTS.md §ADR-Disziplin: accepted
ADRs werden nicht umgeschrieben); ob bei Cluster-T_close eine
neue Folge-ADR „JSON-CLI ausgeliefert" angelegt wird, entscheidet
der Cluster-Slice.
