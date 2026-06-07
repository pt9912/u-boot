# Slice V1: `up --json` / `down --json` — read-only Compose-Status-Envelope

> **Status:** `open/`. Sechster Folge-Slice (6/9) des Cluster-Slice
> [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
> (T0-(e) Reihenfolge 6/9). **Read-only-Klasse**: weder `--dry-run`
> noch `--diff` (Cluster-Slice Z. 464-467), nur `--json` mit
> typisiertem Data-Carrier. **Bündelt up+down** in einem Slice
> (Cluster-T0-(e) Z. 369-372): beide read-only-JSON, gemeinsamer
> Compose-Status-Reader, Confirmer-Pattern bei `down --volumes`
> (Spec §1015) 1:1 erbbar aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> T0-(j).
>
> Erbt das `Data any`-Wire-Field aus
> [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
> T0-(c)/(d): explizit *"für `slice-v1-cli-json-dry-run-up-down`"*
> vorgesehen mit Vorbild
> `type upStatusData struct { Services []serviceStatus
> ` + "`json:\"services\"`" + ` }`. Erbt
> `driving.WarningEntry`-Type (inkl. `Subject`-Feld) aus
> [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
> T2 (R7-MED-F2 + R9-MED-F2 + R12-LOW-F4): bewusst generisch für
> up/down *"recreate-Warnings"* (Multi-Service-WARN "container
> 'postgres' will be replaced") plus config-set (8/9) value-warnings.

## Auslöser

Cluster-Slice §T0-Outcomes (a) macht jeden read-only-Subcommand
für `--json` verbindlich (`LH-NFA-USE-004` §1813). `up` und `down`
sind nach `doctor`/`add`/`init`/`generate`/`remove` die nächsten
und einzig verbleibenden modifying-CLI-Subcommands mit
Compose-Side-Effects: sie schreiben **nichts** auf das lokale
Filesystem (nur `ReadFile(compose.yaml)`), aber sie ändern den
Docker-Daemon-State (Container starten/stoppen, optional Volumes
entfernen).

Spec-Bezug:

- `LH-FA-UP-001` (Umgebung starten, §955-§978) — `u-boot up` mit
  `--timeout`-Stabilisierung
- `LH-FA-UP-002` (Docker Compose verwenden, §980-§986)
- `LH-FA-UP-003` (Startstatus anzeigen, §988-§1000) — Service-
  Name, Containerstatus, Port, Healthcheck-Status als
  Mindestangaben → direkter Carrier-Field-Anker für
  `upStatusData`
- `LH-FA-UP-004` (Umgebung stoppen, §1003-§1019) — `u-boot down`
  mit `--volumes`-destructive-Opt-in
- `LH-FA-CLI-005A` §234/§246/§254 — Confirmation-Gate für
  `down --volumes` (gemeinsam mit `remove --purge`, `init
  --force`, `config set` destructive paths)
- `LH-NFA-USE-004` §1813 / §1841 — Minimalkontrakt-Pflicht
- `LH-FA-CLI-007` §322-417 — Voll-Schema-Vertrag (NICHT für
  up/down weil keine `--dry-run`-Variante, aber das
  `Data any`-Field gehört zur `cliJSONEnvelope`-Struktur)

Heute-Stand-Pre-Scan
(`internal/adapter/driving/cli/up.go`, 110 LOC;
`internal/adapter/driving/cli/down.go`, 117 LOC;
`internal/hexagon/application/upservice.go`, 521 LOC;
`internal/hexagon/application/downservice.go`, 131 LOC):

| Aspekt | up | down |
| --- | --- | --- |
| Positional-Args | `cobra.NoArgs` (`up.go:63`) | `cobra.NoArgs` (`down.go:73`) |
| Lokale Flags | `--timeout` (default 60, 0=fire-and-forget) | `--volumes` (destructive opt-in) |
| Persistente Flags read-through | `--quiet` | `--yes`, `--no-interactive`, `--quiet` |
| FS-Mutation | KEINE | KEINE |
| FS-Read | `ReadFile(compose.yaml)` einmal in `readComposeFile` | `ReadFile(compose.yaml)` analog |
| Docker-Mutation | `engine.ComposeUp` + `engine.ComposePs`-Polling | `engine.ComposeDown` (+ optional `ComposeRm --volumes`) |
| WARN-Emission heute | `renderUpDiagnostics(stdout, resp.Result.Diagnostics, quiet)` — auf **stdout** (M6 §T6) | `renderDownSuccess` — minimal output, keine Diagnostics |
| Confirmer-Pfad | KEIN | `down --volumes` ruft Confirmer analog `remove --purge` (Spec §1015 + §254) |
| ProgressSink | `stderr` für Compose-pull/create/start/healthcheck (LH-NFA-PERF-002) | `stderr` für compose-down-phases (LH-NFA-PERF-002) |
| Exit-Codes | 10 (no-yaml), 11 (Docker unreachable), 12 (Compose-Runtime) | 10 (no-yaml, confirmation refused), 11, 12 |
| Allowlist heute | `jsonallowlist.go:74-75` Reject mit Follow-up `up-down` | analog |

Use-Case-Deps `UpService`: `driven.FileSystem`, `driven.YAMLCodec`,
`driven.DockerEngine`, `driven.NetProbe`, `driven.Clock`,
`driven.Logger`. Use-Case-Deps `DownService`: `driven.FileSystem`,
`driven.YAMLCodec`, `driven.DockerEngine`, `driven.Confirmer`,
`driven.Logger` (+ optional). **KEIN** RecordingFileSystem-Bedarf
(read-only FS), **KEIN** PreviewMode-Bedarf, **KEIN** Diff-
Renderer-Bedarf.

## Aufhebungsbedingung

Zwei Flag-Kombinationen pro Subcommand liefern spec-konforme
Outputs:

```bash
u-boot up                                # Human-Mode (heutiges Verhalten)
u-boot up --json                         # Minimal+Data-Envelope
u-boot up --json --timeout=0             # analog fire-and-forget
u-boot up --json --quiet                 # semantisch identisch zu --json (Cluster-T0-(a))

u-boot down                              # Human-Mode (heutiges Verhalten)
u-boot down --json                       # Minimal+Data-Envelope, KEIN Confirmer-Prompt
u-boot down --volumes --yes --json       # Confirmer geswapped, Volumes removed, Envelope
u-boot down --volumes --no-interactive --json   # ErrConfirmationRequired-Envelope, exit 10
u-boot down --volumes --json             # default interactive prompt — analog remove?
                                          # (T0-Sub-Decision: Prompt im JSON-Mode unterdrücken oder erlauben?)
```

`make gates` grün (lint + test + coverage-gate ≥ 90 % + docs-check).

## Akzeptanzkriterien (vorläufig — T0-Review präzisiert)

- ✅ **`--json`-Allowlist-Migration** (R2-LOW-1 Form-
  Klarstellung): `jsonAllowlist()` ist `map[string]bool` mit
  Cobra-`CommandPath` als Key (`cli/jsonallowlist.go`). Migration
  fügt **zwei separate Einträge** hinzu: `"u-boot up": true` UND
  `"u-boot down": true`. KEIN gemeinsamer `"u-boot up-down"`-Key
  (das ist nur der Folge-Slice-Name, nicht der Cobra-Pfad). Der
  Reject-Mapping in `jsonallowlist.go:74-75` (Follow-up `"up-down"`)
  bleibt für die Übergangsphase relevant solange einer der zwei
  noch nicht migriert ist; beide gleichzeitig migrieren ist
  einfachste Form (eine PR-Pflicht aus T0-(a) Bündelung).
- ✅ **Envelope-Shape**: `command="up"` bzw. `command="down"`,
  KEIN `subcommand`-Feld (beide sind Top-Level-Subcommands ohne
  Sub-Form). KEIN `dryRun`/`diff`/`plannedFiles`/`changes`/`hunks`-
  Feld (read-only-Klasse). Pflicht-Felder pro Spec §1841:
  `status`/`command`/`diagnostics`/`exitCode` plus typed `data`-
  Carrier.
- ✅ **`upStatusData`-Carrier-Form** (doctor T0-(c)/(d) Vorlauf,
  T0-(g) Review-Finding MED-1 Field-Korrektur + R2-MED-3
  Pointer-Konsistenz-Klärung): pro Service `{name, state, port,
  healthcheck}` als `serviceStatus`-Sub-Struct mit `json:"…"`-
  Tags. **Single `port string`** (NICHT `ports []string`) —
  matched heutigen `domain.ServiceStatus.Port`-Display-String-
  Vertrag (`cli/statusview.go:11ff`, Spec LH-FA-UP-003
  Mindestangabe "Port" Singular). Multi-Port-Form wäre eigener
  Sub-Decision-Pfad (Use-Case-Layer-Anpassung nötig — Folge-
  Slice falls Real-World-Bedarf). **Pointer-Wrap-Disziplin
  (R2-MED-3 festgezurrt)**: alle vier Felder als plain
  Go-Strings mit `omitempty`-Tag, **KEIN** Pointer-Wrap:
  ```go
  type serviceStatus struct {
      Name        string `json:"name"`
      State       string `json:"state"`
      Port        string `json:"port,omitempty"`
      Healthcheck string `json:"healthcheck,omitempty"`
  }
  ```
  Begründung: (a) `domain.ServiceStatus.Port` ist `string` plain
  (`domain/serviceup.go:155-181`) — Pointer-Wrap im Wire-Type
  bräuchte Conversion-Layer (`if s.Port != "" { p := s.Port;
  carrier.Port = &p }`), unnötiger Aufwand. (b) Remove's
  `*bool`-Pattern (`removeEnvelopeData.VolumesPurged`)
  rechtfertigt sich durch Three-State-Disambiguation
  (Success-false vs. Error-zero); Port/Healthcheck haben keine
  Three-State-Sit — Empty-String `""` ist semantisch identisch
  zu "kein Port/kein Healthcheck", omitempty droppt beide ohne
  Verlust. (c) Pattern-Konsistenz mit doctor's
  `serviceStatus`-Vorbild (T0-(c) Z.565-575) das ebenfalls
  plain-Strings nutzt. **Name** und **State** sind Pflicht-
  Felder ohne omitempty (Spec LH-FA-UP-003 Mindestangabe).
- ✅ **`downStatusData`-Carrier-Form** (T0-(h) revidiert nach
  Review-Finding HIGH-1 + Followup): matched den heutigen
  Port-Vertrag `DownResponse{RemovedVolumes bool}`
  (`port/driving/down.go:80`). Der Port-Kommentar verbietet
  explizit Counts/Namen: *"No stop / removed counters — docker
  compose down emits a human-readable progress stream rather
  than a structured count, and inventing an 'unknown' sentinel
  value would force every caller to special-case it. If a
  future slice needs precise counts (e.g. for --json output,
  LH-NFA-USE-004 V1), it would add a ComposePs diff
  before/after the call rather than parse the stderr stream."*
  Carrier-Form: `{removedVolumes bool}` — Spec-konform-minimal,
  matched Port. **Kein `omitempty`** auf dem Feld
  (Followup-Pin): `false` ist der legitime Success-Wert
  "nichts entfernt" und MUSS im Erfolgs-Envelope explizit
  erscheinen, sonst kann der Konsument Key-Abwesenheit nicht
  von "kein --volumes gesetzt" disambiguieren. Pattern-Erbe
  remove's `*bool` ist hier nicht nötig weil `down` keine
  Three-State-Disambiguation hat (Error-Pfad trägt `data=nil`
  laut T0-(i), Success-Pfad immer `bool` mit klarem Wert).
  **Feldname-Anmerkung**: `removedVolumes` reflektiert den
  Port-Vertrag-Namen 1:1. `volumesRemoved` wäre als JSON-Name
  natürlicher (Substantiv + Past-Participle), aber Spiegelung
  Port-↔-Wire ist Pattern-konsistent mit remove's
  `volumesPurged`/`VolumesPurged`. ComposePs-Diff-Form für
  `[]string`-Namen ist expliziter Folge-Slice-Pfad (Out-of-Scope,
  siehe T0-(h) unten).
- ✅ **Idempotenz-Pin**: `down` gegen bereits-gestoppte Umgebung
  liefert `removedVolumes: false`, `status: ok`, Exit 0 (analog
  remove NoOp-Semantik — `false` ist der valide
  "nichts-zu-removen"-Wert).
- ✅ **Empty-Array-Pin** (T0-(j) Review-Finding MED + R2-LOW-3
  Pre-Scan-Korrektur): leere Service-Listen MÜSSEN als `[]`
  serialisieren, NICHT `null`. Pre-Scan-Realität (R2-LOW-3):
  `upservice.go:84` returnt heute `domain.UpResult{Stabilized:
  false, Diagnostics: [...]}` — `Services`-Feld ist NICHT
  initialisiert, also **nil-Slice** (NICHT `[]`). Naive
  `json.Marshal(nil-slice)` würde `null` produzieren. T5-CLI-
  Layer MUSS bei nil-Slice mit `[]serviceStatus{}`
  initialisieren (oder Conversion-Layer mit explicit Pre-
  Allokation). Besonders relevant bei `up --timeout=0` (fire-
  and-forget, heute `Services: nil`) und bei Mid-Failure-
  Pfaden. T6-Pin prüft `json.Unmarshal → []serviceStatus{}`,
  nicht `[]serviceStatus(nil)`. Pattern-Erbe doctor
  `diagnostics: []` (Spec §1846-1852 Beispiel).
- ✅ **`--quiet --json` semantisch identisch zu `--json`**
  (Cluster-T0-(a) Pattern aus doctor T6-Pin): `--quiet` darf den
  JSON-Output NICHT unterdrücken — JSON ist die Maschinen-
  Schnittstelle, `--quiet` ist die Human-Schnittstelle-Unterdrückung
  und kollidiert semantisch.
- ✅ **ProgressSink-Silencing im JSON-Mode**: Compose-Phase-
  Streaming auf stderr (LH-NFA-PERF-002 "pull/create/start/
  healthcheck phases stream to stderr live") MUSS in `--json`
  unterdrückt werden, sonst polluten Live-Phasen-Logs den stderr
  für JSON-Konsumenten. Pattern-Erbe von init T0-(o)
  `ProgressPort`-Silencing — aber up/down nutzen `ProgressSink
  io.Writer` (kein Port), also via Request-Field
  `req.ProgressSink = io.Discard` oder Service-Field-Swap.
  T0-Sub-Decision: welche der zwei Formen.
- ✅ **`down --volumes`-Confirmer-Branch** (R2-HIGH-1 ↔ T0-(d)-
  Synchronisierung): `req.SilenceConfirmer = flags.JSON` triggert
  einen **Request-time Gate-Branch** im `DownService.Down()`-
  Code-Pfad (T0-(d) Option (b), nicht Service-Field-Mutation —
  pre-R1-Wortlaut "analog remove `RemoveServiceService.Remove()`-
  Wrapper" war Drift). Use-Case verzweigt im `runConfirmation
  Gate(ctx, req)`-Aufruf auf `req.SilenceConfirmer` und ruft
  entweder `s.confirmer` oder einen lokalen `noopConfirmer{}` —
  kein State-Mutiert, kein neuer `downMu` nötig. **Branch-
  Semantik festgezurrt** (R2-MED-2): bei `--volumes --json`
  OHNE `--yes` → `noopConfirmer.ConfirmRemoveVolumes` returnt
  `(false, nil)` → `ErrConfirmationRequired`-Envelope mit
  `LH-FA-INIT-005`/Exit 10 (Symmetrie zum `--no-interactive`-
  Pfad). Konsistenz-Vertrag: JSON-Mode-Konsumenten erleben das
  **Confirmer-Silencing als Refuse-by-Default**, NICHT als
  Implicit-Auto-Confirm. User MUSS `--yes` explizit setzen für
  destructive `--volumes` im JSON-Mode.
- ✅ **WARN-Migration in `diagnostics[]`**: heutige
  `renderUpDiagnostics`-Calls auf stdout (M6 §T6) wandern im
  JSON-Mode in `diagnostics[]` mit `level: "warn"` und passenden
  LH-Codes. Multi-Service-WARNs (z. B. *"container 'postgres'
  will be replaced"*) nutzen das proaktiv eingeführte
  `Subject`-Feld auf `driving.WarningEntry` (R12-LOW-F4-Vorlauf
  aus remove T2).
- ✅ **Mapper-Tabelle** (analog `mapRemoveErrorToDiagnostic`,
  R2-HIGH-2 Sentinel-Liste erweitert): per-Subcommand-Mapper
  `mapUpErrorToDiagnostic` und `mapDownErrorToDiagnostic` mit
  geteiltem internem Helper `mapComposeRuntimeSentinel` in
  `cli/composesentinel.go` (T0-(e) R2-LOW-2). **Switch-Order
  FS-first** (T0-(f) R2-HIGH-2): die zwei neuen FS-Sentinels
  matchen vor Docker-Klasse. Sentinels:
  `driving.ErrUpFileSystem`/`driving.ErrDownFileSystem` (NEU
  in T2) → `LH-NFA-REL-003`/Exit 14;
  `driven.ErrDockerUnavailable` → `LH-NFA-REL-003`/Exit 11
  (T0-(f) Konsolidierung mit Doku/Test-Pin);
  `driven.ErrComposeRuntime` → `LH-NFA-REL-003`/Exit 12 (dito);
  `ErrConfirmationRequired` → `LH-FA-INIT-005`/Exit 10
  (geteilt mit init/remove); `ErrConflictingModeFlags` →
  `LH-FA-CLI-005A`/Exit 2; `ErrInvalidTimeout` →
  `LH-FA-CLI-006`/Exit 2; `ErrStabilizationTimeout` →
  `LH-FA-UP-001`/Exit 12 (Compose-Runtime-Klasse, Up-spezifisch);
  `ErrComposeFileMissing` → `LH-FA-UP-001`/Exit 10 (fachliche
  Klasse); `ErrProjectNotInitialized` → `LH-FA-INIT-001`/Exit
  10 (Pattern-Erbe init/add/remove — NICHT `LH-FA-UP-001`,
  weil ProjectNotInitialized eine cross-cutting Klasse ist).
- ✅ **Mid-Operation-Failure-UX** (T0-(i) revidiert nach
  Review-Finding HIGH-2 + Followup): heute liefert `UpService`
  bei `ComposeUp`-Fehlern (`upservice.go:76-80`) UND bei
  terminalen Poll-Failures (`upservice.go:200-202`) eine
  **Zero-Response** zurück, keinen Snapshot. Plus:
  `domain.ContainerState`-Enum (`domain/serviceup.go:20ff`)
  kennt nur `unknown|starting|running|restarting|dead`, KEIN
  `failed`. Plan-Empfehlung: Failure-Pfad trägt **nur**
  `diagnostics[]`-Eintrag mit Failure-Service-Name + Failure-
  State + Exit-Code 12 (für `ErrComposeRuntime`) / 11 (für
  `ErrDockerUnavailable`) / 10 (für `ErrProjectNotInitialized`).
  `data` ist `nil` auf Error-Pfad (Zero-Response analog
  generate Error-Envelope). **Mapper-als-Single-Source-of-Truth-
  Pin** (Followup): wenn `data=nil`, MUSS der Mapper
  (`mapUpErrorToDiagnostic` / `mapDownErrorToDiagnostic`) ALLE
  relevanten Failure-Felder in `diagnostics[0]` liefern —
  `code`, `level: "error"`, `message` (mit Failure-Service-
  Name + Terminal-State falls anwendbar), plus `exitCode` am
  Envelope-Top-Level. Kein `data.lastObservedService` o.ä.
  Backup-Channel. T6-Pin verifiziert das End-to-End:
  konstruierter `ErrComposeRuntime`-Failure mit
  "postgres reached terminal state dead" liefert
  `diagnostics[0].message` als einzige Service-Name-Quelle.
  Partial-Snapshot-Form (Snapshot der teilweise gestarteten
  Services bis zur Failure-Stelle) wäre eigener Application-
  Port-Contract mit T0-Sub-Decision-Pfad — siehe T0-(i) unten.
  **Call-Site-Pin** (R2-MED-4): `runUp`/`runDown` rufen
  `reportError(out, sanitizeBaseDir(err, cwd), nil, ..., "up"/
  "down", mapErr, **nil**)` mit `data` als **interface{} nil**,
  NICHT als Zero-Value-Struct `upStatusData{}` oder
  `downStatusData{}`. Zero-Value-Struct würde `services: null`
  serialisieren (genau die Empty-Array-Pin-T0-(j)-Verletzung).
  Pattern-Erbe remove `remove.go:299` (data nil bei Pre-
  Service-Validation-Pfaden) — up/down haben keinen
  Service-Kontext auf Error-Pfad weil kein positional Arg
  existiert (`cobra.NoArgs`).
- ✅ **CLI-Pin-Tests**: ~10-14 Acceptance-Tests in
  `up_acceptance_test.go` + `down_acceptance_test.go` (oder
  einer gebündelter Form `updown_acceptance_test.go`).
- ✅ **`cli-json-output.md`-Update**: §6-Tabelle (up-down→done),
  neue §6.7-Sektion mit Pattern-Vorgabe aus doctor T0-(d), §7
  Mutations-Matrix-Zeilen ("`up`: nur ReadFile", "`down`: nur
  ReadFile" — beide read-only).
- ✅ **CHANGELOG `### Added`-Eintrag** analog
  remove/generate/init.

## Sub-Decisions (TODO — füllt sich in Review-Runden)

- **T0-(a) Bündelung up+down in einem Slice — wie viel Code-
  Sharing?** Cluster-T0-(e) Z. 369-372 sagt "denselben Compose-
  Status-Reader brauchen". Konkret: gibt's ein gemeinsames
  `composeStatusEnvelope`-Helper-Pattern, oder bleibt jeder
  Subcommand selbst-tragend mit kopierter Envelope-Logik?
  Sub-Decision-Optionen:
  (a) Beide Subcommands rufen einen geteilten
      `writeComposeStatusJSON(out, data, warnings)`-Helper,
      lebt in `cli/composestatus.go` (neu).
  (b) Jeder Subcommand hat eigenes `writeUpJSON` /
      `writeDownJSON`, kein geteilter Helper — Pattern wie
      remove `writeRemoveJSON`. Mehr Code-Duplikation, weniger
      Abstraktions-Last.
  Plan-Empfehlung: (b) — die Carrier-Types unterscheiden sich
  (Services-Array vs. RemovedVolumes-Array), das geteilte
  Pattern wäre nur `newDataEnvelope`-Call + Allowlist-Eintrag
  (beides existierende Helper). Helper-Extraction-Druck reift
  erst nach 8/9 (consolidation-Slice).
- **T0-(b) `--quiet --json` Pattern festzurren** (Cluster-T0-(a)
  doctor T6 Vorbild): `--quiet` UNterdrückt im Human-Mode
  status-table + diagnostics-Section auf stdout; im JSON-Mode
  MUSS der Envelope erscheinen. Code-Pfad-Pin: `if flags.JSON
  { … return writeUpJSON(...) } if flags.Quiet { return nil }
  … renderUpStatus(...)`. T6-Pin verifiziert beide Reihenfolgen
  `--quiet --json` und `--json --quiet`.
- **T0-(c) ProgressSink-Silencing-Form**: drei Optionen:
  (a) `req.ProgressSink = io.Discard` aus dem CLI bei
      `flags.JSON` — request-time Substitution, Use-Case sieht
      das Discard-Writer.
  (b) Service-Field-Swap analog init T0-(o) ProgressPort —
      braucht Service-Mutation und defer-Restore innerhalb der
      Lock-Region (gibt's heute keinen `upMu` / `downMu`).
  (c) CLI-Layer-Wrap: ein `discardOnJSONWriter`-Decorator,
      injiziert in den ProgressSink-Field.
  Plan-Empfehlung: (a) — `io.Discard` ist die schlankste Form,
  Use-Case-Signatur unverändert, kein Mutex-Erfordernis.
- **T0-(d) `down --volumes` Confirmer-Pattern** (Review-
  Finding MED-3 Form-Korrektur): drei Optionen:
  (a) **Service-Field-Mutation mit defer-Restore** PLUS neuer
      `downMu sync.Mutex` — vollständig analog remove T0-(j)
      (`removeservice.go:159-178`). Race-Sicherheit fordert
      Mutex; ohne ihn wäre Field-Swap nicht race-frei.
      `DownService` hat heute KEINEN Mutex.
  (b) **Request-time Gate-Branch** ohne Field-Mutation: der
      Use-Case verzweigt im Code-Pfad selbst auf
      `req.SilenceConfirmer` und benutzt entweder
      `s.confirmer` oder einen lokalen `noopConfirmer{}`.
      Kein Service-State mutiert → kein Mutex nötig → race-frei
      by construction.
  (c) **Request-time Confirmer-Field** in `DownRequest`:
      `req.Confirmer driven.Confirmer` (optional, default
      Service-Field). CLI injiziert `noopConfirmer{}` bei
      `flags.JSON`. Schlankste Form, aber bricht heutigen
      `DownRequest`-Vertrag (`SilenceConfirmer bool` wäre
      Pattern-Erbe-konsistenter).
  Plan-Empfehlung **WECHSELT auf (b)** Request-time Gate-
  Branch: kein Service-State mutiert (race-frei), kein
  neuer Mutex nötig (kleinere Application-Layer-Erweiterung),
  Pattern-Erbe-Konsistenz mit remove bleibt nur **konzeptuell**
  (gleiches Ergebnis im JSON-Mode) — nicht 1:1 strukturell.
  remove brauchte Field-Mutation weil `runPurgeGate`-Aufruf
  außerhalb der Verzweigung lag und nur über `s.confirmer`
  drauf zugreifen konnte; in `down` ist die Confirmer-Nutzung
  lokaler und kann via Branch direkt ausgewählt werden.

  **Branch-Semantik festgezurrt** (R2-MED-2 Fix): bei
  `--volumes --json` OHNE `--yes` MUSS der Branch
  `noopConfirmer{}.ConfirmRemoveVolumes` aufrufen (returnt
  `(false, nil)`) → fällt durch in den existierenden
  `downservice.go:128`-Pfad `--volumes declined by user: %w
  ErrConfirmationRequired` → Exit 10. **Konsistenz-Vertrag**:
  JSON-Mode-Confirmer-Silencing ist **Refuse-by-Default**,
  NICHT Implicit-Auto-Confirm. Symmetrie zum
  `--no-interactive`-Pfad (`downservice.go:120` returnt
  ebenfalls `ErrConfirmationRequired`/Exit 10).
  Direkter-Skip-Path (proceed wie `AssumeYes`) ist explizit
  **verworfen** weil er destructive Operation ohne expliziten
  User-Consent im JSON-Mode triggern würde —
  Security-by-Default-Verletzung. T6-Pin:
  `TestDown_VolumesJSONWithoutYes_EmitsErrConfirmationRequired`
  verifiziert Exit 10 + `LH-FA-INIT-005`-Diagnostic.
- **T0-(e) Mapper-Tabelle Layer-Heim** (R2-LOW-2 Heim-
  Festzurrung): zwei separate Mapper
  `mapUpErrorToDiagnostic`/`mapDownErrorToDiagnostic` ODER ein
  gemeinsamer `mapComposeErrorToDiagnostic(err, command
  string)`. Heutige Sentinels überlappen stark
  (`ErrDockerUnavailable`, `ErrComposeRuntime`, `ErrProject
  NotInitialized` sind in beiden), aber `down`-spezifisch sind
  `ErrConfirmationRequired`/`ErrConflictingModeFlags`/
  `volumes`-Sentinels und `up`-spezifisch `ErrInvalidTimeout`/
  `ErrStabilizationTimeout`. Plan-Empfehlung: **separate
  Mapper** analog `mapRemoveErrorToDiagnostic`/
  `mapAddErrorToDiagnostic` (Pattern-Erbe), aber geteilter
  Helper für die Docker-/Compose-Runtime-Sentinels als
  interner `mapComposeRuntimeSentinel(err) (code string,
  exitCode int, matched bool)`-Helper. **Helper-File-Heim
  festgezurrt** (R2-LOW-2): neuer File
  `cli/composesentinel.go` (neben `cli/jsonenvelope.go`).
  Begründung: (a) `jsonenvelope.go` hat heute keine Mapper-
  Logic — Mischung würde Layer-Grenze auflösen. (b) Eigener
  File macht Pattern-Wiederverwendung für künftige Compose-
  Subcommands (z. B. ein hypothetisches `restart`-Subcommand)
  sauberer. (c) File-Größe wird ~30-50 LOC, klein aber
  fokussiert. T5-Cell ergänzt um das File-Heim.
- **T0-(f) LH-Code-Klassifikation für Docker-Runtime-Klasse**
  (Review-Finding MED-2 Risiko-Klarstellung):
  `ErrDockerUnavailable` → Exit 11 ist gesetzt, aber welcher
  LH-Code? Spec hat `LH-NFA-REL-003` für FS-Failure und
  `LH-FA-CLI-006` als Default. Für Docker-Daemon-
  Unverfügbarkeit gibt's keinen dedizierten LH-Code.
  Sub-Decision: (i) neuer `LH-NFA-REL-005`-Code ODER (ii)
  Konsolidierung auf existierenden `LH-NFA-REL-003` mit
  Sub-Semantik-Dehnung (analog remove's Triple-Use von
  `LH-FA-ADD-005`).
  Plan-Vorschlag: **(ii) Konsolidierung mit explizitem
  Doku-/Test-Pin-Block**. Risiko (Review-MED-2): `LH-NFA-REL-003`
  ist heute im Repo stark mit technischen Persistenz-/FS-
  Fehlern UND Exit 14 assoziiert. Bei Konsolidierung MUSS der
  Slice drei Pins liefern:
  (1) **Doku-Pin** in `cli-json-output.md` §6.7: dass derselbe
      `LH-NFA-REL-003`-Code mit Exit 11 (Docker-Daemon) oder
      Exit 12 (Compose-Runtime) erscheinen kann, NICHT nur 14
      (FS) — Disambiguation via `(code, exitCode)`-Tupel
      analog remove's `LH-FA-ADD-007` Multi-Use (ERROR + WARN
      via `(code, level)`).
  (2) **Test-Pin** `TestUp_DockerUnavailable_DiagnosticCodeIs
      RELN003_ExitCode11` verifiziert die Kombination explizit.
  (3) **Mapper-Switch-Order-Pin** verifiziert dass FS-Klasse
      VOR Docker-Klasse matched bei Multi-`%w`-Wrap.
      **R2-HIGH-2 Sentinel-Bedarf-Klärung**: heute existiert
      WEDER `ErrUpFileSystem` NOCH `ErrDownFileSystem` NOCH
      ein cluster-weiter `driven.ErrFileSystem`. Die FS-Read-
      Wraps in `upservice.go:105/138/148` und
      `downservice.go:81/97` nutzen rohes `%w` ohne typed
      Sentinel — sie fallen heute auf Default-Mapper
      `LH-FA-CLI-006`/Exit 1. **Konsequenz**: ohne neuen
      Sentinel ist der Switch-Order-Pin ein Phantom-Test
      (es gibt keinen FS-Sentinel der matchen würde —
      konstruierter Multi-Wrap wäre nur `ErrDockerUnavailable
      + ErrComposeRuntime`, beide Docker-Klasse). **Fix**:
      T2-Cell um zwei neue Port-Sentinels erweitern:
      `driving.ErrUpFileSystem` in `port/driving/up.go` und
      `driving.ErrDownFileSystem` in `port/driving/down.go`
      (Pattern-Erbe `driving.ErrRemoveFileSystem`). T3
      migriert die fünf FS-Wrap-Stellen
      (`upservice.go:105/138/148`, `downservice.go:81/97`)
      auf Multi-`%w` mit dem neuen Sentinel. Mapper-Tabelle
      T0-(e) ergänzt um `ErrUpFileSystem`/
      `ErrDownFileSystem` → `LH-NFA-REL-003`/Exit 14
      (kanonische FS-Klasse). Switch-Order-Pin ist dann
      **real, nicht Phantom**: konstruierter
      `fmt.Errorf("%w: %w", ErrUpFileSystem,
      ErrDockerUnavailable)` MUSS `LH-NFA-REL-003`/Exit 14
      liefern, nicht `LH-NFA-REL-003`/Exit 11. Pattern-Erbe
      remove `mapRemoveErrorToDiagnostic` Switch-Order T0-(e).
  Alternative-Wechsel auf neuen `LH-NFA-REL-005`: zieht
  Spec-Erweiterung und Lastenheft-Edit. Plan bleibt bei (ii)
  weil Spec-Footprint-Stabilität V1-prioritär.
- **T0-(g) `upStatusData`-Field-Granularität** (Review-Finding
  MED-2 Port-Vertrag-Korrektur): doctor T0-(c) zitiert das
  Vorbild als `Services []serviceStatus`, aber `serviceStatus`-
  Felder sind nicht festgenagelt. Spec LH-FA-UP-003 fordert:
  Name, Containerstatus, Port (Singular), Healthcheck
  (optional). Heutiger `domain.ServiceStatus.Port`
  (`cli/statusview.go:11ff`) ist **single Display-String**, NICHT
  `[]string`. Sub-Decision-Optionen:
  (i) **`port string`** matched heutigen Port-Vertrag
      Spec-konform-minimal.
  (ii) **`ports []string`** braucht Use-Case-Layer-Anpassung
       (`domain.ServiceStatus.Port` zu `Ports []string` umbauen
       plus alle Aufrufstellen). Bricht Pattern.
  (iii) Beide Felder mit `port` deprecated → `ports`: zwei
        JSON-Keys parallel für eine Übergangszeit. Doppelarbeit,
        wenig Nutzen.
  Plan-Empfehlung **(i) single `port string`** — matched Port-
  Vertrag, kein Domain-Refactor, JSON-Konsument kann via
  `strings.Split(port, ", ")` parsen falls Mehrfach-Ports
  drinstehen. Multi-Port-Form als eigener Folge-Slice falls
  Real-World-Bedarf (Domain-Erweiterung notwendig).
- **T0-(h) `downStatusData`-Field-Definition** (Review-Finding
  HIGH-1 Port-Vertrag-Korrektur): heutiger Port liefert
  `DownResponse{RemovedVolumes bool}` (`port/driving/down.go:80`).
  Der Port-Kommentar verbietet explizit Counts/Namen:
  *"No stop / removed counters — docker compose down emits a
  human-readable progress stream rather than a structured
  count, and inventing an 'unknown' sentinel value would force
  every caller to special-case it. If a future slice needs
  precise counts (e.g. for --json output, LH-NFA-USE-004 V1),
  it would add a ComposePs diff before/after the call rather
  than parse the stderr stream."* Drei Sub-Decision-Optionen:
  (i) **`{removedVolumes bool}`** — 1:1-Echo von
      `DownResponse.RemovedVolumes`. Spec-konform-minimal,
      kein Port-Refactor. JSON-Konsument bekommt einen
      Boolean-Status statt einer namensbasierten Liste.
  (ii) **`{removedVolumes []string}` mit ComposePs-Diff**:
       expliziter Application-Port-Vertrag, der vor und nach
       `ComposeDown` ein `ComposePs --filter "label=…
       project=<n>" --format json` aufruft und die Differenz
       als Volume-Namen-Liste trägt. Großer Architektur-
       Eingriff: neuer `DockerEngine.ListVolumes`-Port-Method,
       zusätzlicher Compose-Daemon-Roundtrip, Volume-vs-
       Container-Naming-Disambiguation, Roll-back-Semantik bei
       Mid-Failure. Nicht V1-würdig (Port-Kommentar verweist
       explizit auf Folge-Slice).
  (iii) Hybrid — `{removedVolumesEcho bool, removedVolumeNames
        []string}` mit Names als optional (omitempty). Doppel-
        Field zieht Klassifikations-Verwirrung.
  Plan-Empfehlung **(i) `bool`** — matched heutigen Port,
  Spec-konform, kein Architektur-Eingriff. Option (ii) als
  **Out-of-Scope** Carveout mit eigenem Folge-Slice
  `slice-v1-down-volumes-named-list` (Trigger: Real-World-
  Konsumenten-Bedarf nach Namen-Liste).
- **T0-(i) Mid-`ComposeUp`-Failure-Capture-Vertrag** (Review-
  Finding HIGH-2 Port/Enum-Vertrag-Korrektur):
  Pre-Plan-Empfehlung-Vorschlag *"`state: "failed"` analog
  `domain.ContainerState`-Enum + Snapshot bis zur Failure-
  Stelle"* ist nicht umsetzbar:
  (a) `domain.ContainerState` kennt nur `unknown|starting|
      running|restarting|dead`, KEIN `failed`
      (`domain/serviceup.go:20-40`).
  (b) `UpService` returnt bei `ComposeUp`-Fehlern
      (`upservice.go:76-80`) UND bei terminalen Poll-Failures
      (`upservice.go:200-202`) eine **Zero-Response** — kein
      Snapshot. Architektur-Vertrag heute: Error-Pfad ohne
      Data.
  Drei Sub-Decision-Optionen:
  (i) **Failure nur in `diagnostics[]`** — `data` ist `nil`
      auf Error-Pfad (analog generate Error-Envelope-Form aus
      generate T0-(q)). Diagnostic-Eintrag trägt Failure-
      Service-Name + Failure-State + LH-Code + Exit 11/12/10.
      Pattern-konsistent mit heutigem Port-Vertrag, kein
      Architektur-Eingriff.
  (ii) **Partial-Snapshots-Application-Port-Contract**: neuer
       `UpResponse.PartialServices []domain.ServiceStatus`-Feld
       PLUS Use-Case-Refactor (Z. 76/200) um vor dem
       Error-Return einen `ComposePs`-Snapshot zu ziehen.
       Großer Eingriff (drei Aufrufstellen in `upservice.go`
       müssen Snapshot statt Zero-Response liefern; Enum-
       Erweiterung um `StateFailed` mit Migrationspflicht für
       alle bestehenden Switch-Statements). Nicht V1-würdig.
  (iii) Hybrid mit `data.lastObservedServices []serviceStatus`
        nur im JSON-Pfad: CLI-Layer ruft `ComposePs` nochmal
        bei Error. Bricht Layer-Trennung (CLI macht Docker-
        Side-Effects).
  Plan-Empfehlung **(i) Failure nur in `diagnostics[]`** —
  matched heutigen Port-Vertrag, Spec-konform-minimal, kein
  Application-Refactor. Option (ii) als **Out-of-Scope**
  Carveout mit eigenem Folge-Slice
  `slice-v1-up-partial-snapshot-on-failure` (Trigger: Real-
  World-Bedarf nach "was lief schon" bei Mid-Failure-Debugging
  + Domain-Enum-Erweiterung).
- **T0-(j) `--timeout=0` Fire-and-Forget im JSON-Mode**
  (Review-Finding MED-2 Empty-Array-Klarstellung): heute
  bedeutet `--timeout=0` "no polling, no probes, status table
  omitted, info diagnostic shown" (`up.go:53-54`). Im JSON-Mode
  bleibt das info-diagnostic problematisch (`level: "info"` ist
  Spec §1834 verboten — siehe doctor-Pattern aus §97 doctor-
  Slice). Sub-Decision-Optionen:
  (i) `level: "warn"` Upgrade — semantisch unscharf weil
      Fire-and-Forget ein User-Wunsch ist, kein Problem.
  (ii) **Field-Drop des info-diagnostic + Marker**
       `data.timeoutFireAndForget: true` im Carrier.
       Konsumenten-Klassifikation via Marker-Field, kein
       Severity-Niveau-Stretching.
  (iii) Field-Drop ohne Marker — Konsument kann nur indirekt
        ableiten (`services: []` UND Exit 0).
  Plan-Empfehlung **(ii) Field-Drop + Marker** (Followup-
  Präzisierung): `timeoutFireAndForget` MUSS in genau **einem**
  Modus erscheinen — `true` ausschließlich bei `--timeout=0`-
  Pfad; in jedem anderen Up-Pfad **fehlt** das Feld komplett
  (Key-Abwesenheit via Pointer-Wrapping `*bool` mit
  omitempty), NICHT als explizit `false` getragen. Sub-
  Decision-Form: `TimeoutFireAndForget *bool
  ` + "`json:\"timeoutFireAndForget,omitempty\"`" + ` ` — analog
  remove's `*bool`-Pattern für `volumesPurged`-Key-Presence-vs-
  Absence-Disambiguation. **Empty-Array-Pin** für BEIDE
  Carrier-Felder explizit (Followup): `data.services: []` UND
  `diagnostics: []` (NICHT `null`) — beide Felder OHNE
  omitempty serialisieren, bei nil-Slice mit `[]…{}`
  initialisieren. Pattern-Analog doctor `diagnostics: []`
  (Spec §1846-1852 Beispiel). T6-Pins:
  `TestUp_TimeoutZero_JSON_ServicesIsEmptyArrayNotNull` und
  `TestUp_AllStable_JSON_DiagnosticsIsEmptyArrayNotNull` —
  beide mit `json.RawMessage`-Re-Marshal-Check (verifiziert
  Byte-Sequenz `"services":[]` statt `"services":null`).
- **T0-(k) Recreate-Warnings-Semantik** (R12-LOW-F4 aus remove
  T2 setzt den Type-Vorlauf): wann emittiert `up` eine WARN
  *"container 'postgres' will be replaced"*? Heute existiert
  in `upservice.go` keine recreate-Detection. Sub-Decision:
  (i) Recreate-Detection als V1-Scope ODER (ii) als Carveout
  für Folge-Slice (analog Volume-Auto-Removal aus remove).
  Plan-Empfehlung: (ii) — Recreate-Detection braucht
  Compose-Plan-Pre-Walk (`docker compose config`-Parse +
  Container-Hash-Vergleich), nicht-trivial. WARN-Carrier-Type
  ist proaktiv da, aber konkrete Detection kommt später.

## Tranchen (vorgeschlagen — präzisiert in T0-Outcomes)

| T | Inhalt | LOC (Schätzung) | Voraussetzung |
| - | --- | --- | --- |
| T0 | Discovery + Sub-Decisions (a)-(k) klären; Review-Runden | — (Plan) | — |
| T1 | **Entfällt** (analog remove T1): `noopConfirmer` lebt bereits in `application/noop.go:17-33`, `io.Discard` ist Go-stdlib — beide Helper für ProgressSink-Silencing und Confirmer-Swap existieren | — (entfällt) | T0 |
| T2 | Port-Types: `UpRequest.SilenceProgress` request-time ProgressSink-Discard; `DownRequest.SilenceConfirmer`-Feld; `UpResponse.Warnings []driving.WarningEntry` (Type schon da aus remove T2). **Zwei neue Port-Sentinels** (R2-HIGH-2 Fix): `driving.ErrUpFileSystem` (`port/driving/up.go`) und `driving.ErrDownFileSystem` (`port/driving/down.go`) analog `driving.ErrRemoveFileSystem` — für die fünf FS-Read-Wraps (`upservice.go:105/138/148`, `downservice.go:81/97`) als Multi-`%w`-Wrap-Target | ~120 | T0 |
| T3 | Application-Layer: ProgressSink-Discard-Wiring im JSON-Mode (request-time analog T0-(c)); `DownService.Down()`-Request-time Gate-Branch ohne Field-Mutation (T0-(d) Option (b) — kein `downMu` nötig, race-frei) — KEIN remove-1:1-Service-Field-Swap mit Mutex. **Multi-`%w`-Wrap-Migration** (R2-HIGH-2 Fix) der fünf FS-Read-Stellen auf `ErrUpFileSystem`/`ErrDownFileSystem`. Mapper-Helper `mapComposeRuntimeSentinel(err)` in `cli/composesentinel.go` für die geteilten Docker/Compose-Sentinels (T0-(e) R2-LOW-2) | ~120 | T2 |
| T4 | Composition-Root: heute existiert KEINE up/down-fsFactory (kein Recorder) — T4 prüft ob `cmd/uboot/main.go` Wiring-Updates braucht (vermutlich nur Confirmer-Bezug für down) | ~20 | T3 |
| T5 | CLI-RunE: `runUp`/`runDown` ruft generische Helper mit `command="up"`/`"down"`, `mapErr=mapUpErrorToDiagnostic`/`mapDownErrorToDiagnostic`; Envelope-Pfade; Allowlist-Migration mit zwei separaten Einträgen (R2-LOW-1); Mapper neu; `data`-Structs (`upStatusData`/`downStatusData`); WARN-Migration. **`baseDirSanitizedError`-Helper-Extraktion** (R2-MED-5 Fix): geteilter Helper aus `cli/remove.go:465-491` nach neuem `cli/sanitize.go` extrahieren; `runUp`/`runDown` rufen `sanitizeBaseDir(err, cwd)` vor `reportError` analog `remove.go:299` | ~250 | T2 |
| T6 | Acceptance-Tests: ~12-16 Tests (Envelope-Pin both Subcommands, Idempotenz-Pin für down, `--quiet --json`-Pin, ProgressSink-Silencing-Pin, Confirmer-Branch-Pin für `down --volumes --json` ohne `--yes` (R2-MED-2), ConflictingModeFlags-Pin, Service-Sentinels-Pins, Multi-`%w`-Switch-Order-Pin für FS+Docker (R2-HIGH-2), Path-Leak-Sanitizer-Pin (R2-MED-5), Empty-Array-Pins für services+diagnostics (R2-LOW-3)) | ~500-600 | T5 |
| T7 | Review-Fix-Rounds (~1-2 Runden bei Pattern-Erbe) | ~50 | T6 |
| T8 | Closure: CHANGELOG, cli-json-output.md §6/§6.7/§7 (zwei §7-Zeilen "nur ReadFile"), roadmap done-Zähler 5→6, **carveouts.md** drei neue Einträge (Recreate-Warnings, Volume-Named-Liste, Partial-Snapshot — Multi-Port als vierter falls Real-World-Trigger), **vier open/-Stubs** schaffen (R2-MED-1 Fix Memory-`carveouts_need_plans`): `slice-v1-recreate-detection`, `slice-v1-down-volumes-named-list`, `slice-v1-up-partial-snapshot-on-failure`, ggf. `slice-v1-multi-port-services`. Slice nach `done/` mit DoD-Hash-Tabelle | — (Doku) | T7 |

LOC-Bilanz vorläufig: **~1060-1160** (R2-Korrektur nach oben:
T2+40 für zwei neue FS-Sentinels, T3+40 für Multi-`%w`-Migration
der 5 FS-Wrap-Stellen, T5+50 für `sanitizeBaseDir`-Helper-
Extraktion + Aufrufe, T6+100 für drei neue Pin-Klassen
(Sanitizer, Switch-Order, Confirmer-Branch). Pattern-Erbe von
remove (Confirmer-Konzept + FS-Sentinel-Pattern + Sanitizer-
Pattern), init (ProgressPort-Silencing-Vorbild) und doctor
(typed Data-Carrier).

## Out of Scope

- **Recreate-Detection** (T0-(k) Sub-Decision): Compose-Plan-Pre-
  Walk + Container-Hash-Vergleich für *"container 'postgres' will
  be replaced"*-WARN ist V1-Out-of-Scope. WARN-Carrier-Type ist
  via `driving.WarningEntry`-Vorlauf (remove T2 R12-LOW-F4)
  proaktiv vorhanden, konkrete Detection wandert in Folge-Slice
  (Trigger: User-Feedback über fehlende Replace-Warnings oder
  Cluster-T_close-Audit).
- **`down --volumes` Named-Volume-Liste** (T0-(h) Option (ii)):
  `removedVolumes` ist heute `bool` auf dem Port-Vertrag
  (`port/driving/down.go:80`). Named-Liste braucht einen neuen
  `DockerEngine.ListVolumes`-Port-Method plus ComposePs-Diff-
  Pattern vor/nach `ComposeDown` (so der heutige Port-Kommentar
  selbst). Folge-Slice
  `slice-v1-down-volumes-named-list` (Trigger: Real-World-
  Konsumenten-Bedarf nach Namen-Liste z. B. für Audit-Logs oder
  CI-Cleanup-Scripts; aktueller `removedVolumes: bool` ist
  Spec-konform-minimal).
- **Partial-Snapshot bei Mid-`ComposeUp`-Failure** (T0-(i)
  Option (ii)): heutige Zero-Response auf Error-Pfad
  (`upservice.go:76/200`) reflektiert den Port-Vertrag.
  Partial-Snapshot brauchte (a) `UpResponse.PartialServices
  []domain.ServiceStatus`-Feld, (b) Use-Case-Refactor um vor
  Error-Return einen `ComposePs`-Snapshot zu ziehen,
  (c) `domain.ContainerState`-Enum-Erweiterung um `StateFailed`
  mit Migrations-Pflicht für alle Switch-Statements. Großer
  Architektur-Eingriff. Folge-Slice
  `slice-v1-up-partial-snapshot-on-failure` (Trigger: Real-
  World-Bedarf nach "was lief schon"-Mid-Failure-Debugging,
  z. B. interaktive CI-Diagnose).
- **`up --service <name>`-Selective-Form**: heute liefert `up`
  alle Compose-Services (`cobra.NoArgs`); ein zukünftiger
  Sub-Form für Single-Service-Start wäre eigener Slice mit Args-
  Validator-Pattern (würde dann auch das envelope-consolidation-
  Pattern erben).
- **Args-Validator-Pattern-Drift** mit add/init/generate
  (R15-Cross-Slice-1 aus remove): up/down landen VOR der
  Konsolidierung
  ([`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md))
  und übernehmen den heutigen Drift-Status mit — **soweit
  Args-Validator betrifft**: `cobra.NoArgs` (`up.go:63`,
  `down.go:73`), kein positional-Arg-Pfad, kein
  `validateRemoveArgs`-Erbe nötig. Bleibt out-of-scope.

  **R2-MED-5 Korrektur — `baseDirSanitizedError`-Bedarf IST
  vorhanden**: Pre-R2-Wortlaut "kein BaseDirSanitizedError-
  Bedarf für Compose-Runtime-Errors" war faktisch falsch.
  Tatsächlicher Path-Leak-Inventar (R2-MED-5 Audit):

  | File | Zeile | Pfad-Source | Klasse |
  | --- | --- | --- | --- |
  | `upservice.go` | 80 | `req.BaseDir` (abs) | Compose-Runtime |
  | `upservice.go` | 105 | `filepath.Join(baseDir, "u-boot.yaml")` | FS-Read |
  | `upservice.go` | 108 | dito | FS-fachlich |
  | `upservice.go` | 138 | `filepath.Join(baseDir, "compose.yaml")` | FS-Read |
  | `upservice.go` | 141 | dito | FS-fachlich |
  | `upservice.go` | 146/148 | dito | FS-Read |
  | `downservice.go` | 69 | `req.BaseDir` (abs) | Compose-Runtime |
  | `downservice.go` | 81 | `filepath.Join(baseDir, "u-boot.yaml")` | FS-Read |
  | `downservice.go` | 84 | dito | FS-fachlich |
  | `downservice.go` | 97 | `filepath.Join(baseDir, "compose.yaml")` | FS-Read |
  | `downservice.go` | 100 | dito | FS-fachlich |

  Im JSON-Mode wandern diese absoluten Pfade via `err.Error()`
  in `diagnostics[0].message` — exakt der R15-Cross-Slice-1-
  Defekt-Klasse. **T5 MUSS** den
  `sanitizeBaseDir`-Wrapper-Helper aus `cli/remove.go:465-491`
  in einen geteilten Helper extrahieren (analog T0-(e) Mapper-
  Helper-Heim) und in `runUp`/`runDown` vor `reportError`
  anwenden. Helper-Heim festgelegt unter R15-LOW-1-Korrektur:
  `cli/sanitize.go` (oder analog `cli/composesentinel.go`-
  Pattern). Pattern-Erbe `cli/remove.go:299` mit
  `sanitizeBaseDir(removeErr, cwd)`-Aufruf.

  **Sub-Sequenz-Folge** für envelope-consolidation-Stub
  Update: der existierende Stub
  ([`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md))
  beschreibt Pattern-Drift nur für add/init/generate, NICHT
  für up/down — der existierende Stub-Trigger "Cluster-Stand
  8/9" wäre semantisch eh in Cluster-T_close-Nähe, also
  passt das. Wenn up/down den Sanitizer-Wrap selbst nutzen
  (geteilter Helper), trägt der up/down-Slice den Pattern-Fix
  bereits — der envelope-consolidation-Slice braucht nur noch
  add/init/generate retroaktiv zu refactorn.
- **Multi-Port-Form** für Services mit mehreren exposed Ports
  (T0-(g) Option (ii)): heutiger `domain.ServiceStatus.Port`
  ist single Display-String. Umstellung auf `Ports []string`
  bricht das Domain-Pattern und ist eigener Slice falls
  Real-World-Druck (z. B. Multi-Port-Service-Health-Reporting).
- **Docker-Daemon-Version-Reporting im Envelope** (z. B.
  `data.dockerVersion`): nicht in Spec gefordert, eigener Slice
  falls Konsumenten-Bedarf.

## Bezug

- Cluster:
  [`slice-v1-cli-json-dry-run`](../in-progress/slice-v1-cli-json-dry-run.md)
  (Folge-Slice 6/9).
- Pattern-Vorbilder:
  [`slice-v1-cli-json-dry-run-doctor`](../done/slice-v1-cli-json-dry-run-doctor.md)
  (Data-Carrier `upStatusData`-Vorbild + Read-Only-Envelope-Form),
  [`slice-v1-cli-json-dry-run-init`](../done/slice-v1-cli-json-dry-run-init.md)
  (ProgressPort-Silencing-Vorbild im JSON-Mode),
  [`slice-v1-cli-json-dry-run-remove`](../done/slice-v1-cli-json-dry-run-remove.md)
  (Confirmer-Swap-Pattern T0-(j) + `driving.WarningEntry`-Type).
- Code-Anker:
  [`cli/up.go`](../../../../internal/adapter/driving/cli/up.go),
  [`cli/down.go`](../../../../internal/adapter/driving/cli/down.go),
  [`application/upservice.go`](../../../../internal/hexagon/application/upservice.go),
  [`application/downservice.go`](../../../../internal/hexagon/application/downservice.go),
  [`cli/jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go)
  Z. 29/74-75.
- Folge-Slices:
  [`slice-v1-cli-json-envelope-consolidation`](slice-v1-cli-json-envelope-consolidation.md)
  (R15-Cross-Slice-1 — retroaktive Helper-Extraction, Trigger
  Cluster-Stand 8/9).
- Phase: V1 (Teil des V1-pünktlichen Cluster-Slices).
