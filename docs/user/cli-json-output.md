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
| 4 | `slice-v1-cli-json-dry-run-generate` | `generate` | T0-Discovery + R1-R5 adressiert, `next/` |
| 5 | `slice-v1-cli-json-dry-run-remove` | `remove` | offen |
| 6 | `slice-v1-cli-json-dry-run-up-down` | `up`, `down` (gebündelt, read-only Compose-Status) | offen |
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

Andere modifying-Subcommands (`generate`, `remove`,
`config set`) ergänzen ihre Zeile im jeweiligen Slice — heute
`doctor`, `add` und `init` migriert.

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
