# Slice V1: Maschinen-lesbare CLI — `--json`, `--dry-run`, `--diff`

> **Status:** geplant für v0.4.0+ — Spec ✅
> (`LH-FA-CLI-007` Dry-Run [`spec/lastenheft.md:302-447`](../../../../spec/lastenheft.md),
> `LH-FA-CLI-008` Diff [`spec/lastenheft.md:451-489`](../../../../spec/lastenheft.md),
> `LH-NFA-USE-004` Maschinen-lesbar
> [`spec/lastenheft.md:1809-1853`](../../../../spec/lastenheft.md)),
> ADR-Anker ✅ ([`ADR-0010`](../../adr/0010-kein-http-driving-adapter.md)
> — JSON-CLI ist *die* Maschinen-Schnittstelle, HTTP-Adapter
> verworfen). **Cluster-Slice, kein Code-Implementation-Slice:**
> definiert Reihenfolge + geteilte Konventionen für die
> Per-Command-Folge-Slice-Serie. T0 ✅ festgezurrt
> (§T0-Outcomes — 5 Sub-Decisions plus Mutations-Matrix-Pre-Scan);
> in `in-progress/`. **Cluster-Stand (2026-06-08): 9/9 Folge-Slices
> done — Cluster-Slice selbst noch `in-progress/` (T_close pending).**
>
> **Done (9/9 Folge-Slices)**:
> [`doctor`](../done/slice-v1-cli-json-dry-run-doctor.md) (1/9),
> [`add`](../done/slice-v1-cli-json-dry-run-add.md) (2/9),
> [`init`](../done/slice-v1-cli-json-dry-run-init.md) (3/9),
> [`generate`](../done/slice-v1-cli-json-dry-run-generate.md) (4/9),
> [`remove`](../done/slice-v1-cli-json-dry-run-remove.md) (5/9),
> [`up-down`](../done/slice-v1-cli-json-dry-run-up-down.md) (6/9),
> [`logs`](../done/slice-v1-cli-json-dry-run-logs.md) (7/9),
> [`config`](../done/slice-v1-cli-json-dry-run-config.md) (8/9 —
> T0–T8 + drei Review-Runden; erster Read-only+Modifying-Hybrid),
> [`template`](../done/slice-v1-cli-json-dry-run-template.md) (9/9 —
> T0→T2→T4, `template list --json` Array→Envelope; T3 (bare-Reject)
> nach T_close verschoben).
>
> **Cluster-T_close (offen — der letzte Schritt der Serie)**: Alle
> neun Folge-Slices sind in `done/` → die Closure-Hard-Rule ist
> erfüllt. T_close baut die Übergangs-Mechanik ab (Allowlist-Map +
> `applyJSONRejectGate` + `PersistentPreRunE`) UND führt dabei den
> bare-`template`-RunE-Reject (`cli.ErrTemplateSubcommandRequired`,
> envelope-LOS §1838) + Help-Leak-Pin ein (aus template-T3 hierher
> verschoben); verifiziert, dass alle neun Forms nach Mechanik-Abbau
> korrekt antworten; optional eine Folge-ADR „JSON-CLI ausgeliefert"
> (ADR-0010-Nachfolger). Danach Cluster-Slice selbst nach `done/`.

## Auslöser

Drei V1-Pflicht-Spec-IDs fordern eine maschinen-lesbare CLI:

- **`LH-FA-CLI-007`** (Dry-Run, V1): für dateiverändernde
  Befehle muss `--dry-run` zeigen, welche Dateien erzeugt /
  geändert / gelöscht würden, **ohne** Dateisystem-Schreiben.
  Pflicht-JSON-Schema definiert (`$schema` draft 2020-12,
  Pflichtfelder `status`, `command`, `dryRun`, `diff`,
  `plannedFiles`, `changes`, `diagnostics`, `exitCode`).
- **`LH-FA-CLI-008`** (Diff, V1): `--diff` zeigt Unterschiede
  zwischen aktuellem und geplantem Zustand der betroffenen
  Dateien. Kombinierbar mit `--dry-run`. Bei `--diff --json`
  gilt das LH-FA-CLI-007-Schema mit `diff: true`.
- **`LH-NFA-USE-004`** (Maschinen-lesbar, V1): `--json` für
  **alle** zehn Subcommands (Spec-Enum:
  `init, add, remove, up, down, doctor, logs, generate,
  config, template`). Minimalkontrakt für read-only
  Subcommands: `status`, `command`, `diagnostics`, `exitCode`.
  Bei gruppierten Befehlen (`template`, `config`)
  zusätzlich `subcommand`-Pflicht.

Heute existiert genau ein `--json`-Pfad im Repo:
[`template list --json`](../../../../internal/adapter/driving/cli/template.go)
(`renderTemplateListJSON` mit `templateJSON`-DTO + nil-Slice→`[]`-
Normalisierung). Der Pfad ist Vorbild für die DTO-Lokalisation
(Driving-Adapter besitzt das Wire-Format, Domain bleibt
präsentations-agnostisch — `LH-FA-ARCH-002` / ADR-0002), trägt
aber **noch nicht** das `LH-FA-CLI-007`-Pflicht-Schema mit
`status`/`command`/`diagnostics`/`exitCode` — auch der existierende
`template list --json`-Pfad muss im Zuge dieser Slice-Cluster-
Serie spec-konform werden (Folge-Slice
`slice-v1-cli-json-dry-run-template`).

ADR-0010 §Entscheidung verbindet diesen Slice mit der
Architektur: **JSON-CLI ist die kanonische Maschinen-
Schnittstelle**; HTTP-Adapter wurde mit ADR-0010 §Folgepunkte
Re-Eval-Trigger 2 ausdrücklich gegen das hier kommende
`LH-NFA-USE-004`-Surface abgewogen. Wenn dieser Slice slipt,
würde ADR-0010 selbst angreifbar (Folge-Trigger 2 aus §144 —
„Maschinen-Schnittstelle über LH-NFA-USE-004 hinaus"). Deshalb
ist `slice-v1-cli-json-dry-run` V1-pünktlich zu liefern, nicht
„V1+1" oder „nach Trigger".

Roadmap-Notiz aus
[`slice-v1-logs`](../done/slice-v1-logs.md) §Auslöser:
„`--json`-Mode kommt im Folge-Slice `slice-v1-cli-json-dry-run`
nachträglich auf `logs` drauf." Das ist die einzige
Code-Abhängigkeit: `logs` ist V0.4.0 ausgeliefert (✅ `e9a5392`
+ `357e40a`), `--json` für `logs` ist ein klassisches Read-only-
Subcommand-Beispiel und liegt im Cluster.

## Aufhebungsbedingung

Alle zehn Spec-Enum-Subcommands tragen einen `--json`-Pfad
gemäß `LH-NFA-USE-004`-Minimalkontrakt:

```bash
u-boot doctor --json                # alle read-only-Befehle
u-boot logs --json                  # (existierend: stream-orientiert,
                                    # Output-Format ist Sub-Decision)
u-boot template list --json         # bereits ✅, Schema-Audit nötig
u-boot config --json                # bare config — Listing/Default-View
u-boot config get <path> --json     # read-only Pfad-Lookup
u-boot up --json                    # Compose-Up-Status-Report
u-boot down --json                  # Compose-Down-Status-Report
```

Alle dateiverändernden Subcommands (`init`, `add`, `remove`,
`generate`, `config set`) tragen zusätzlich `--dry-run` +
`--diff`:

```bash
u-boot add postgres --dry-run --json
u-boot add postgres --diff
u-boot add postgres --dry-run --diff --json
u-boot config set <path> <value> --dry-run --diff --json
```

**`config`-Cluster-Pflicht (Review-Finding MEDIUM):**
`u-boot config`, `u-boot config get <path>`, `u-boot config
set <path> <value>` sind drei separate Subcommand-Formen
unter dem `config`-Hauptkommando. Alle drei brauchen
`--json` (`LH-NFA-USE-004` gilt für alle
Spec-Enum-Subcommands, nicht nur Schreibpfade). LH-FA-CLI-007
fordert für `command == "config"` zusätzlich
`subcommand`-Pflicht — der Wert für **bare** `u-boot config`
(ohne weiteren Pfad) ist Sub-Decision im Folge-Slice
`slice-v1-cli-json-dry-run-config` (Kandidaten: `"list"`,
`"show"`, oder explizit `""` falls Spec leeren Subcommand
erlaubt — Aufklärung gegen `LH-FA-CONF-005` im Folge-Slice).

Jede `--json`-Ausgabe validiert gegen das
`LH-FA-CLI-007`-Pflicht-Schema (oder den
`LH-NFA-USE-004`-Minimalkontrakt für read-only). Validierung
in Tests pinnt:

- `status` an höchstem `diagnostics.level` gekoppelt
  (`error → "error"`, `warn ohne error → "warn"`, sonst `"ok"`),
- `diagnostics[].code` LH-Kennung-konform (Convention §445),
- `subcommand` bei `command == "template"` / `"config"` gesetzt,
- `plannedFiles[].action` ∈ `{create, modify, delete}`,
- `changes[].count` ≥ 0,
- `exitCode` ≥ 0 und konsistent mit `LH-FA-CLI-006`.

Die einzelnen Subcommand-Implementierungen leben in
**Per-Command-Folge-Slices** (siehe §Per-Command-Folge-Slices).

**Closure-Hard-Rule (Review-Finding HIGH):** Dieser
Cluster-Slice schließt **ausschließlich**, sobald **alle**
Per-Command-Folge-Slices in `done/` sind. Es gibt **kein**
MVP-Quorum, **kein** Verteilungs-Audit-Bypass, **keinen**
Restweg-Carveout als Closure-Alternative — das wäre eine
direkte Aufweichung des V1-Pflicht-Surfaces aus `LH-NFA-USE-
004` und würde ADR-0010 §Folgepunkte Re-Eval-Trigger 2
unterminieren (HTTP-Adapter wurde mit dem Argument verworfen,
dass diese Spur V1-pünktlich kommt).

Wenn ein Folge-Slice tatsächlich slipt (z. B. weil die
Sub-Decision in `next/` blockiert): **vor** dem Cluster-Move
nach `done/` muss
(1) ein expliziter Carveout-Eintrag in
[`carveouts.md`](../in-progress/carveouts.md) §Temporäre
Carveouts erscheinen, der den fehlenden Subcommand benennt
und mit einem benannten Re-Trigger-Slice-Plan-Stub in `open/`
verlinkt (`LH-FA-PROJDOCS-005` Carveout-Plan-Anker-Pflicht);
(2) der **done/-Eintrag dieses Cluster-Slices** und der
**Roadmap-Liefer-Vermerk** dürfen **nicht** „JSON-CLI als
Maschinen-Schnittstelle ausgeliefert" sagen, sondern müssen
den Carveout zitieren und den Re-Trigger-Pfad nennen.
ADR-0010 selbst bleibt **unverändert** (AGENTS.md
§ADR-Disziplin: accepted ADRs werden nicht umgeschrieben);
eine neue Folge-ADR mit abgeschwächter Aussage ist möglich,
aber nicht erzwungen — Sub-Decision von T_close (siehe §AK
„ADR-0010-Liefer-Anker").

Default-Erwartung: keine Slips, alle 9 Folge-Slices schließen.
Der Slip-Pfad ist Notfall-Restlauf, **kein wählbarer**
Closure-Pfad.

## Akzeptanzkriterien (Cluster-Ebene)

- ✅ **Schema-Vertrag dokumentiert**: zentraler
  Reference-Block (vermutlich `docs/user/cli-json-output.md`
  oder Sektion im Architecture-Doc) zitiert das `LH-FA-CLI-
  007`-Schema verbatim und benennt die DTO-Lokalisation
  gemäß T0-(c) (**Common-Envelope `cliJSONEnvelope`** im
  CLI-Adapter mit Subcommand-spezifischem `Data`-Feld).
- ✅ **Per-Command-Folge-Slices angelegt** für alle zehn
  Spec-Enum-Subcommands — verteilt auf **9 Folge-Slices**
  (`up`/`down` gebündelt, weil beide read-only-JSON sind
  und denselben Compose-Status-Reader brauchen; siehe
  §Per-Command-Folge-Slices). Jeder Stub mit eigenem
  T0-Discovery und LOC-Schätzung. Reihenfolge gemäß T0-(e)
  festgezurrt.
- ✅ **Erster Folge-Slice abgeschlossen**
  (`slice-v1-cli-json-dry-run-doctor`, gemäß T0-(e)) als
  belastbares Pattern-Vorbild für read-only-Envelope +
  Schema-Helper. Der zweite Folge-Slice
  (`slice-v1-cli-json-dry-run-add`) trägt das modifying-
  Surface-Vorbild (Recorder + Dry-Run + Diff).
- ✅ **Schema-Konformitäts-Helper** im CLI-Adapter (oder als
  Test-Helper in `internal/adapter/driving/cli/jsontestutil/`):
  parst die `--json`-Ausgabe und prüft Pflichtfelder. Jeder
  Folge-Slice verwendet ihn in seinen Tests, damit Schema-
  Drift einheitlich kracht.
- ✅ **ADR-0010-Liefer-Anker** dokumentiert — **ohne**
  inhaltlichen Rewrite der accepted ADR (AGENTS.md
  §ADR-Disziplin). Auslieferungs-Anker ist primär der
  `done/`-Eintrag dieses Cluster-Slices (DoD-Hash-Line +
  Folge-Slice-Verweise); ADR-0010 selbst wird **nicht**
  in §Konsequenzen umgeschrieben. Sub-Decision für T_close:
  ob zusätzlich eine neue **Folge-ADR** angelegt wird, die
  ADR-0010 als Vorgänger referenziert und „JSON-CLI
  ausgeliefert" als eigenständigen Entscheid trägt, **oder**
  ob der Roadmap-Liefer-Vermerk plus done/-Slice
  ausreichen. ADR-Disziplin entscheidet T_close, nicht
  jetzt.
- ✅ **Roadmap-Status** zeigt den Cluster-Fortschritt: jede
  Folge-Slice-Closure aktualisiert die v0.4.0-Tabelle in
  `roadmap.md`.

## T0-Discovery (vor `next/`-Übergang festzulegen)

Sub-Decisions, die der Cluster-Slice klären muss, bevor die
ersten Per-Command-Folge-Slices spawnen. **Status:** alle fünf
Sub-Decisions plus Mutations-Matrix-Pre-Scan ✅ festgezurrt —
siehe [§T0-Outcomes](#t0-outcomes) für die Entscheidungen mit
Begründung. Layout analog [`slice-v1-logs`](../done/slice-v1-logs.md)
§T0-Outcomes.

### T0-(a) Globale Flag oder per Subcommand?

`--json` / `--dry-run` / `--diff` als persistente Cobra-Flags
am Root-Command (jeder Subcommand erbt sie) **oder** explizit
pro Subcommand registriert?

- Pro Root-Flag: weniger Wiring, konsistenter Help-Output,
  aber Subcommands wie `template list` müssen das Flag bewusst
  zurückweisen falls nicht unterstützt.
- Pro Subcommand: explizit, jeder Folge-Slice trägt eigenes
  Flag-Wiring, kein Risiko ungewollter Vererbung.

Vorbild: `--verbose` / `--debug` sind heute persistent (root-
level, slice-followup-verbosity-wiring §`7c6fbce`); `--json`
auf `template list` ist lokal definiert.

### T0-(b) Wo lebt die Dry-Run-Logik?

Drei Architektur-Varianten:

1. **Application-Layer-Flag im Request** (`Request.DryRun bool`):
   Use-Case entscheidet, ob mutierende FileSystem-Methoden
   aufgerufen werden. Pro: Adapter bleibt symmetrisch; Contra:
   jeder Use-Case trägt Dry-Run-If-Verzweigungen an jeder
   FS-Mutation.
2. **Recording-FileSystem-Wrapper** (driven-Adapter-Variante):
   `RecordingFileSystem` implementiert `driven.FileSystem`,
   capturet **alle 8 mutierenden Methoden** statt sie
   auszuführen (Review-Finding MEDIUM: vollständige Liste —
   `WriteFile`, `WriteFileExclusive`, `Mkdir`, `MkdirAll`,
   `Rename`, `RemoveAll`, `Copy`, `CopyExclusive`; vgl.
   [`driven.FileSystem`](../../../../internal/hexagon/port/driven/filesystem.go),
   permanent-Carveout interfacebloat). Use-Case weiß nichts
   vom Dry-Run-Modus. Pro: Use-Case sauber; Contra: ALLE acht
   Mutations-Methoden müssen geschlossen capturet werden,
   sonst lückt der Plan; `Rename`/`RemoveAll` müssen pro
   Folge-Slice in `plannedFiles[].action` als `delete`/`modify`
   gemappt werden (Spec-Enum nur `create|modify|delete`, nicht
   `rename`).
3. **ChangeSet-Pattern** (separates Apply-Step): Use-Case
   berechnet `ChangeSet`, ein separater `Apply`-Step führt
   alle 8 Mutationen aus. Pro: Dry-Run = Apply weglassen;
   Contra: alle Use-Cases müssen auf ChangeSet-Pattern
   refactoren — größter Eingriff.

T0-Decision sollte den Eingriffs-Radius pro Variante gegen die
Folge-Slice-LOC-Schätzungen abwägen.

**Mutations-Matrix-Pflicht (Review-Finding MEDIUM):**
unabhängig davon welche Variante T0-(b) wählt, muss der erste
Folge-Slice eine **vollständige Mutations-Matrix** liefern:
pro modifying Subcommand (`init`, `add`, `remove`, `generate`,
`config set`) wird aufgelistet, **welche FS-Mutations-Methoden
er heute aufruft**. Diese Matrix ist die Pin-Grundlage für
zwei Test-Disziplinen, die jeder Folge-Slice tragen muss:

- **Positive Pin:** für jeden modifying Subcommand existiert
  ein `--dry-run`-Test, der die laut Matrix erwartete
  `plannedFiles`-Action je Datei pinnt.
- **Negative Pin** („kein FS-Write"): für jeden modifying
  Subcommand existiert ein `--dry-run`-Test, der nach dem
  Run prüft, dass **null** der 8 FS-Mutations-Methoden
  ausgeführt wurde (Spy/Fake auf Production-FileSystem-Port,
  Call-Count == 0). Das schließt das Risiko aus, dass eine
  vergessene Methode am Dry-Run-Filter vorbeiläuft.

Die Matrix wandert in den Schema-Vertrag-Doc-Block aus T1,
damit Folge-Slices sie referenzieren und beim FS-Port-
Erweiterungen einen Drift-Trigger haben.

### T0-(c) DTO-Lokalisation: gemeinsam oder per Subcommand?

Spec-Pflichtfelder (`status`, `command`, `dryRun`, `diff`,
`plannedFiles`, `changes`, `diagnostics`, `exitCode`) sind
über alle Subcommands gleich. Optionen:

- **Common-Envelope-Type** (`cliJSONEnvelope` mit
  embeddable Subcommand-Spezifik): DRY, aber Schema-Drift
  über alle Folge-Slices auf einmal lösbar.
- **Pro-Subcommand-DTO** mit Field-Tag-Duplikation: weniger
  Helper-Code, aber Drift-Risiko zwischen Subcommand-DTOs.

Vorbild: `templateJSON` im `template list`-Pfad ist
subcommand-spezifisch (kein gemeinsamer Envelope).

### T0-(d) `--diff`-Renderer: Unified oder strukturiert?

`LH-FA-CLI-008` lässt das Format offen. Optionen:

- **Unified-Diff** (`go-diff`-Library oder eigene Impl): klassisch,
  human-lesbar, aber im JSON-Modus als String-Field eingebettet
  schwerer zu konsumieren.
- **Strukturiert** (`{path, hunks: [{oldStart, newStart, lines}]}`):
  maschinen-freundlich, im Human-Mode aber zusätzlich gerendert.

Ggf. beides: human → unified, JSON → strukturiert. Sub-Decision
beeinflusst die LOC-Schätzung der modifying-Folge-Slices stark.

### T0-(e) Reihenfolge der Per-Command-Folge-Slices

Use-Case-Druck (geschätzt):

1. `add --dry-run --diff --json` — höchster CI-Bedarf
   (Service-Add-Plan-Vorschau vor Commit).
2. `doctor --json` — niedrige Komplexität, schon
   read-only, gutes Schema-Pilot.
3. `init --dry-run --diff --json` — Onboarding-Use-Case.
4. `generate --dry-run --diff --json` — Build-Tooling.
5. `remove --dry-run --diff --json` — destructive, höchster
   Dry-Run-Nutzen.
6. `up --json` / `down --json` — read-only-Output von
   Compose-Zustand.
7. `logs --json` — stream-orientiert, Output-Modell-Frage
   (JSON-Lines pro Compose-Log-Zeile vs. Single-Envelope nach
   Stream-Ende). Diese Sub-Decision ist im Folge-Slice
   `slice-v1-cli-json-dry-run-logs` zu treffen — der
   ausgelieferte [`slice-v1-logs`](../done/slice-v1-logs.md)
   hat den `--json`-Pfad bewusst hierher ausgelagert und keine
   Vorab-Entscheidung getroffen (Review-Finding LOW: vorherige
   Stub-Version verwies fälschlich auf logs T0-(b), das aber
   Service-Name-Validation regelt).
8. `config` (alle drei Formen — bare `config`, `config get
   <path>`, `config set <path> <value>`) — `config` und
   `config get` sind read-only-`--json`, `config set` ist
   modifying-`--dry-run --diff --json`. Drei-Form-Bündel,
   damit `subcommand`-Pflicht aus `LH-FA-CLI-007` einmal
   geschlossen wird.
9. `template list --json` — Audit + Schema-Migration
   (existierender Pfad spec-konform machen).

T0-(e) festzurrt die Reihenfolge mit Begründung. Erster
Folge-Slice ist gleichzeitig Pattern-Vorbild — T0-(b)/(c)/(d)
sollten parallel mit dem ersten Folge-Slice festgezurrt werden.

## Per-Command-Folge-Slices

Per-Command-Inkrementell-Strategie (Cluster-Entscheidung):
**kein zentraler Pilot-Slice**, jeder Folge-Slice trägt eigenes
Wiring. Geteilt werden Schema-Vertrag, DTO-Konventionen, Test-
Helper. Gemeinsamer Code-Anker wandert mit dem ersten Folge-
Slice in den Repo (z. B. `cliJSONEnvelope` und
`jsontestutil.AssertSchemaConform`).

Folge-Slice-Plan-Namen (in `open/` zu erzeugen, ein Stub pro):

- `slice-v1-cli-json-dry-run-add`
- `slice-v1-cli-json-dry-run-doctor`
- `slice-v1-cli-json-dry-run-init`
- `slice-v1-cli-json-dry-run-generate`
- `slice-v1-cli-json-dry-run-remove`
- `slice-v1-cli-json-dry-run-up-down`
- `slice-v1-cli-json-dry-run-logs`
- `slice-v1-cli-json-dry-run-config`
- `slice-v1-cli-json-dry-run-template`

Bündelung von `up`/`down` in einem Slice ist sinnvoll, weil
beide read-only-JSON sind und denselben Compose-Status-Reader
brauchen. Alle anderen einzeln, weil jeder eigene Use-Case-
Logik trägt.

## Tranchen (vorgeschlagen — Cluster-Slice-eigene Arbeit)

| T | Inhalt | LOC (Schätzung) |
| - | ------ | --------------- |
| T0 | **Discovery + Sub-Decisions.** Fünf T0-Fragen aus §T0-Discovery klären (Flag-Scope, Dry-Run-Architektur, DTO-Lokalisation, Diff-Renderer, Reihenfolge). Entscheidung pro Frage mit Begründung in einem `T0-Outcomes`-Block dokumentieren. ADR-0010 bleibt unangetastet (§AK ADR-Disziplin). | — (Plan-Arbeit) |
| T1 | **Schema-Vertrag-Doku.** `docs/user/cli-json-output.md` (neu) zitiert `LH-FA-CLI-007`-Schema verbatim, dokumentiert DTO-Konvention, listet Per-Command-Folge-Slice-Reihenfolge. README EN+DE bekommt Verweis-Zeile. **Delegation: geliefert via `slice-v1-cli-json-dry-run-doctor`** ([`done/slice-v1-cli-json-dry-run-doctor.md`](../done/slice-v1-cli-json-dry-run-doctor.md) T1-Tranche, DoD-Hash `299e792`). Cluster-T_close-Reviewer sucht **keinen** separaten Cluster-T1-Commit. | ~80 (reine Doku, via Doctor-Slice) |
| T2..Tn | **Spawn Folge-Slice-Stubs** in `open/` für alle 9 Per-Command-Slices. Pro Stub: Auslöser + grobe AKs + LOC-Schätzung + Verweis auf gemeinsamen Schema-Vertrag. Reihenfolge nach T0-(e). | ~30 LOC pro Stub × 9 = ~270 |
| T_close | **Cluster-Closure.** Pflicht-Bedingung gemäß §Aufhebungsbedingung Closure-Hard-Rule (§124-148): **alle 9** Folge-Slices in `done/`. Punkt. Cluster-Slice mit DoD-Hash-Line aller Folge-Slices nach `done/` plus Roadmap-Update. ADR-0010 bleibt **unverändert**; optionale Folge-ADR ist Sub-Decision (siehe AK „ADR-0010-Liefer-Anker"). **Kein** „kritisches Quorum", **kein** MVP-Bypass, **kein** Restweg-Carveout als Closure-Alternative. Der **Notfall-Slip-Pfad** aus §Aufhebungsbedingung §137-148 (Carveout-Eintrag + abgeschwächte Liefer-Aussage) ist explizit **kein wählbarer** T_close-Pfad — er ist Restlauf-Disziplin für einen unvermeidbaren Slip und tritt nur in Kraft, **nachdem** ein Slip bereits passiert ist (Default-Erwartung bleibt „keine Slips"). | — (Doku) |

LOC-Schätzung Cluster-Slice: ~350 LOC, deutlich unter
800-LOC-Schwelle. Folge-Slice-LOC-Bandbreite: 200..600 je
Subcommand (T0-(b)-Architektur-Decision dominant).

## T0-Outcomes

Fünf Sub-Decisions vor `next/`-Übergang festgezurrt — die fünf
Fragen aus [§T0-Discovery](#t0-discovery-vor-next-übergang-festzulegen)
mit Begründung. Die Mutations-Matrix-Pflicht (Review-Finding
M2) wandert per Outcome (b) als Pre-Scan ins T0 vorgezogen,
die vollständige Pro-Subcommand-Matrix bleibt Lieferung des
ersten modifying Folge-Slices.

### T0-(a) Flag-Scope: `--json` Root-persistent, `--dry-run`/`--diff` per Subcommand

**Entscheidung:** `--json` lebt als PersistentFlag am Root-Cobra-
Command (`u-boot --json …`); jeder der 10 Spec-Enum-Subcommands
erbt es. `--dry-run` und `--diff` werden ausschließlich auf den
5 modifying Subcommands (`init`, `add`, `remove`, `generate`,
`config set`) per `cmd.Flags()` lokal registriert.

**Begründung:** `LH-NFA-USE-004` fordert `--json` für **alle**
zehn Subcommands — Root-Persistent ist die natürliche Heimat
(sonst 10× lokales Wiring + Drift-Risiko, vgl. heutiges Vorbild
`--verbose/--debug/--quiet/--yes/--no-interactive` aus
[`cli/root.go:38-54`](../../../../internal/adapter/driving/cli/root.go),
slice-followup-verbosity-wiring §`7c6fbce`). `--dry-run`/`--diff`
sind nur für die 5 modifying-Subcommands spec-pflichtig
(`LH-FA-CLI-007`/`LH-FA-CLI-008`); persistent würde die 5
read-only-Subcommands zwingen, sie aktiv abzulehnen — 5× extra
Reject-Wiring ohne Gegenwert.

### T0-(b) Dry-Run-Architektur: `RecordingFileSystem`-Wrapper mit Passthrough-Modus (Variante 2)

**Entscheidung:** Ein neuer driven-Adapter
`RecordingFileSystem` implementiert
[`driven.FileSystem`](../../../../internal/hexagon/port/driven/filesystem.go)
und delegiert die **4 Read-Methoden** (`Exists`, `ReadFile`,
`ReadDir`, `Lstat`) an die underlying Production-`fs.FileSystem`.
Für die **8 Mutations-Methoden** (`WriteFile`,
`WriteFileExclusive`, `Mkdir`, `MkdirAll`, `Rename`,
`RemoveAll`, `Copy`, `CopyExclusive`) gilt ein **Passthrough-
Schalter** (Konstruktor-Option / config-Field):

- `Passthrough=false` (Dry-Run-Modus, `--dry-run` aktiv):
  capturen als `plannedFiles`/`changes`-Einträge, **ohne**
  Production-FS aufzurufen.
- `Passthrough=true` (Preview-and-Apply, `--diff` ohne
  `--dry-run`): capturen **und** an Production-FS weiterreichen
  — die Aufruf-Reihenfolge ist „erst aufzeichnen, dann
  durchführen", damit der Preview-Renderer auch im
  Mid-Failure-Fall den geplanten Zustand sehen kann
  (Sub-Decision exakte Semantik im ersten modifying Folge-Slice).

Use-Case-Code bleibt unverändert. Das **Wiring lebt
ausschließlich im Composition-Root**
[`cmd/uboot/main.go`](../../../../cmd/uboot/main.go) — der
**einzige Ort**, an dem `RecordingFileSystem` und die
Production-FS koexistieren. Hard Rule Hexagonale Architektur
([`spec/architecture.md:154`](../../../../spec/architecture.md)
plus depguard `LH-FA-ARCH-002`/`LH-FA-ARCH-003`): der CLI-Adapter
unter `internal/adapter/driving/cli/` darf **keine**
`hexagon/port/driven`-Abhängigkeit tragen.

**Konkretes Wiring-Pattern (Sub-Decision exakter Field-Name
im ersten modifying Folge-Slice):** Für jeden der 5 modifying
Subcommands konstruiert das Composition-Root **zwei
driving-Port-Instanzen** statt einer:

- `addServiceUseCase` (Production-FS, Normal-Mode), und
- `addServicePreviewUseCase` (RecordingFS-wrapped, beide
  Passthrough-Schalter-Stellungen aus T0-(b)).

Beide Instanzen werden in die App-Struktur **als
`driving`-Port-Interface-Typen** injiziert (z. B.
`driving.AddServiceUseCase`). Die CLI-RunE-Funktion wählt
zwischen den **driving-Port-Instanzen** anhand der parsed
`--dry-run`/`--diff`-Flag-Kombination — der CLI-Adapter
**sieht den `driven.FileSystem`-Typ nirgends** und importiert
ihn auch nicht. depguard bleibt sauber.

Für die read-only-Pfade (`doctor`, `logs`, `up`, `down`,
`template list`, bare `config`, `config get`) gibt es nur
**eine** driving-Port-Instanz pro Use-Case; sie tragen weder
`--dry-run` noch `--diff`, also keine Variante.

**Begründung:** Variante (1) (`Request.DryRun bool` mit
If-Branches an jeder Mutation-Site) würde ~25 If-Branches
über 5 Use-Cases ziehen — eine vergessene → stiller Write im
Dry-Run, Negative-Pin-Test catcht erst post-hoc, der Write ist
bereits passiert. Variante (3) (ChangeSet-Pattern) refactort
alle 5 Use-Cases — YAGNI für V1. Variante (2) hält Use-Cases
sauber, das Negative-Pin-Test-Design wird trivial (Recorder
mit `Passthrough=false` zählt Mutations-Aufrufe — Production-FS
sieht null Calls), und der Recorder ist der einzige Ort, an dem
die FS-Mutations-Liste vollständig getragen werden muss → ein
Drift-Anker. Der **Passthrough-Modus** löst gleichzeitig
`LH-FA-CLI-008` `--diff` ohne `--dry-run` (Preview-and-Apply):
derselbe Adapter capturet den Preview-Plan **und** führt die
Writes aus — keine doppelte Use-Case-Ausführung nötig,
Preview-vs-Apply-Drift ausgeschlossen.

**Pre-Scan-Matrix (T0-Investigation, vor `next/`):**
Mutations-Aufrufe in der Production-Code-Pfaden — direkter
Aufruf plus indirekter via [`BackupPath`](../../../../internal/hexagon/application/backup.go):

| Use-Case | Direkt | Indirekt via `BackupPath` |
| --- | --- | --- |
| `init` ([`initproject.go`](../../../../internal/hexagon/application/initproject.go)) | `MkdirAll`, `WriteFile` | `CopyExclusive`, `Mkdir`, `MkdirAll`, `Copy`, `RemoveAll` |
| `add` ([`addservice_execute.go`](../../../../internal/hexagon/application/addservice_execute.go)) | `WriteFile` | — |
| `remove` ([`removeservice.go`](../../../../internal/hexagon/application/removeservice.go)) | `WriteFile`, `RemoveAll` | — |
| `generate` ([`generate.go`](../../../../internal/hexagon/application/generate.go)) | `MkdirAll`, `WriteFile` | — |
| `config set` ([`config.go`](../../../../internal/hexagon/application/config.go)) | `WriteFile` | — |

Heute direkt aufgerufen: `WriteFile`, `MkdirAll`, `RemoveAll`
(plus über Backup-Helper: `CopyExclusive`, `Mkdir`, `Copy`).
Heute **nicht** aus Use-Case-Pfaden: `WriteFileExclusive` (nur
[`doctor.go:183`](../../../../internal/hexagon/application/doctor.go)
für Sentinel-Write-Probe — kein Mutator-Artefakt) und
`Rename`. Der Recorder muss trotzdem alle 8 abdecken, damit
zukünftige Use-Cases keinen Lukentest am Dry-Run-Filter
vorbeischmuggeln.

**Reservation (Lieferpflicht des ersten modifying Folge-Slice
`slice-v1-cli-json-dry-run-add`):**

- `RecordingFileSystem` deckt **alle 8** `driven.FileSystem`-
  Mutations-Methoden ab; `WriteFileExclusive` und `Rename`
  sind im Recorder ebenfalls implementiert, auch wenn heute
  keine Use-Case sie aufruft (Drift-Schutz).
- **Passthrough-Modus-Verträge gepinnt:** Der erste modifying
  Folge-Slice fixiert die exakte Aufruf-Reihenfolge im
  Passthrough=`true`-Pfad (Vorschlag: erst Plan-Eintrag
  capturen, dann Production-Mutation aufrufen — bei
  Mutation-Fehler bleibt der Plan-Eintrag bestehen, das
  Diagnostic-Item trägt den Fehler) und beweist über
  Acceptance-Tests sowohl den `--dry-run`-Pfad (Passthrough=
  false, null FS-Mutationen) als auch den `--diff`-ohne-
  `--dry-run`-Pfad (Passthrough=true, Preview-Output + echte
  Writes in einem Lauf, `LH-FA-CLI-008`).
- **Wiring-Kontrolle:** Das Composition-Root in
  `cmd/uboot/main.go` ist die **einzige** Stelle, an der
  konkrete `driven.FileSystem`-Adapter (Production-FS und
  `RecordingFileSystem`) instanziiert und verkabelt werden.
  Der CLI-Adapter empfängt **driving-Port-Instanzen** (je
  modifying Subcommand zwei — Normal-Mode und Preview-Mode)
  über die App-Struktur und wählt zwischen ihnen anhand der
  Flag-Kombination. Der CLI-Adapter importiert weder
  `driven.FileSystem` noch den `RecordingFileSystem`-Konkret-
  Typ — depguard/Hexagonal-Architektur-Hard-Rule
  (`LH-FA-ARCH-002`/`003`) bleibt grün, prüfbar via
  `make lint`.
- **Read-after-Write-Audit pro Use-Case:** Der erste modifying
  Folge-Slice prüft explizit, ob ein Use-Case in derselben
  Sequenz erst `WriteFile(p, …)` und anschließend `Exists(p)`
  / `ReadFile(p)` / `Lstat(p)` auf demselben Pfad aufruft.
  Stichprobe T0-Investigation zeigt Read-then-Write-Muster
  (config liest u-boot.yaml → schreibt patched; add liest
  catalog → schreibt service-Files), kein Write-then-Read —
  Pflicht-Re-Validation im Folge-Slice mit dokumentiertem
  Ergebnis pro Use-Case.
- **Overlay-Fallback:** Wenn die Re-Validation ein Write-then-
  Read in einem Use-Case findet, erweitert der Recorder seine
  Read-Methoden um eine kleine In-Memory-Overlay-Map
  (geplante Writes überlagern die delegierte Read-Antwort).
  Geschätzter LOC-Aufwand: ~30 LOC; bewusst als Fallback
  und nicht als Default, damit der Recorder bei nicht-
  benötigtem Overlay sauber bleibt. Im Passthrough=true-Modus
  ist Overlay nicht nötig — der echte Write findet statt, der
  reale FS-Zustand stimmt.

### T0-(c) DTO-Lokalisation: Common-Envelope `cliJSONEnvelope` mit Subcommand-Payload-Feld

**Entscheidung:** Ein gemeinsamer Envelope-Type
`cliJSONEnvelope` im CLI-Adapter (Lokation: vermutlich
`internal/adapter/driving/cli/jsonenvelope.go` — Sub-Decision
im ersten Folge-Slice) trägt die `LH-FA-CLI-007`-Pflichtfelder
**einmal**: `Status`, `Command`, `Subcommand` (omitempty,
gesetzt für `template`/`config`), `DryRun`, `Diff`,
`PlannedFiles`, `Changes`, `Diagnostics`, `ExitCode`, plus
ein Subcommand-spezifisches `Data` (`json.RawMessage` oder
`any`) für Read-only-Payloads (z. B. `template list`-Array,
`doctor`-Report).

**Begründung:** Pro-Subcommand-DTOs würden die 8 Pflichtfeld-
Tags 10× duplizieren → Schema-Drift garantiert (eine Slice
schreibt `"exit_code"` statt `"exitCode"` und das Schema-
Conformance-Test der anderen 9 fängt es nicht). Common-
Envelope ist gleichzeitig die natürliche Verankerungsstelle
für den `jsontestutil.AssertSchemaConform`-Helper aus den
Cluster-AKs (ein Validator, der den Envelope parst und gegen
das `LH-FA-CLI-007`-Schema prüft — jeder Folge-Slice ruft
ihn auf). Vorbild `templateJSON` aus
[`cli/template.go:163-186`](../../../../internal/adapter/driving/cli/template.go)
ist subcommand-spezifisch, trägt aber **null** der
Spec-Pflichtfelder — wird im Folge-Slice
`slice-v1-cli-json-dry-run-template` ohnehin auf den Envelope
migriert (Schema-Audit).

### T0-(d) Diff-Renderer: Beides — Unified-String im Human-Mode, strukturierte Hunks per `plannedFiles[]` im JSON-Mode

**Entscheidung:** `LH-FA-CLI-008` lässt das Format offen, daher
zweigleisig.

- **Human-Mode** (`--diff` ohne `--json`): klassischer
  Unified-Diff als String an stdout (`+`/`-`-Prefix,
  Hunk-Header `@@ -oldStart,oldLines +newStart,newLines @@`).
- **JSON-Mode** (`--diff --json`): das Envelope-Feld `diff`
  bleibt **Boolean** und auf `true` gesetzt, wie
  `LH-FA-CLI-007` es als Pflichtfeld definiert. Die
  strukturierten Hunks landen **nicht** in `diff`, sondern
  als **omitempty-Hunk-Array per `plannedFiles[]`-Eintrag**
  (Vorschlag-Field-Name: `plannedFiles[].hunks` mit
  `[{oldStart, oldLines, newStart, newLines, content}]`).
  Begründung: `plannedFiles[]` trägt bereits Pfad und
  `action` pro betroffener Datei — die Hunks gehören
  semantisch dorthin, ein zusätzliches Top-Level-Array
  würde Pfad-Korrelation verdoppeln. Der **exakte
  Field-Name** und ob `hunks` direkt unter `plannedFiles[]`
  oder unter einem `plannedFiles[].diff`-Sub-Objekt sitzt,
  ist Sub-Decision des ersten modifying Folge-Slices —
  Verankerung gegen den `LH-FA-CLI-007`-Schema-Wortlaut
  passiert in T1 (Schema-Vertrag-Doku) und der erste
  Folge-Slice referenziert das.

Beide Modi berechnen denselben LCS-Hunk-Datentyp; der
Unified-String ist ein zweiter Renderer auf demselben Datum.

**Begründung:** Embedded Unified-String im JSON ist tooling-
unfreundlich (Konsumenten müssten den String selbst parsen).
Reine strukturierte Form wäre human nicht direkt lesbar.
Beide Renderer aus einer gemeinsamen Hunk-Repräsentation
spart LOC und schließt Format-Drift zwischen den Modi aus.
**Dep-Policy-Hinweis:** `go.mod` ist diszipliniert minimal
(4 Deps total). Pure-Go LCS+Unified-Diff intern (~150-200
LOC) hält die Linie; ob stattdessen `pmezard/go-difflib`
(0-Dep, tiny, MIT) zulässig ist, ist Sub-Decision des
ersten modifying Folge-Slice.

### T0-(e) Folge-Slice-Reihenfolge: `doctor` zuerst, `add` direkt danach

**Entscheidung:**

1. `slice-v1-cli-json-dry-run-doctor` — **Pattern-Vorbild für
   read-only-Envelope + Schema-Helper.**
2. `slice-v1-cli-json-dry-run-add` — **Pattern-Vorbild für
   modifying-Surface (Recorder, Dry-Run, Diff).**
3. `slice-v1-cli-json-dry-run-init`
4. `slice-v1-cli-json-dry-run-generate`
5. `slice-v1-cli-json-dry-run-remove`
6. `slice-v1-cli-json-dry-run-up-down`
7. `slice-v1-cli-json-dry-run-logs`
8. `slice-v1-cli-json-dry-run-config`
9. `slice-v1-cli-json-dry-run-template`

**Begründung:** Die Stub-Vorschlag-Reihenfolge (`add` zuerst
nach Use-Case-Druck) bündelt im ersten Folge-Slice
gleichzeitig **drei** neue Etablierungen (Envelope, Schema-
Helper, RecordingFileSystem) plus Diff-Renderer — zu viel
Risiko an einer Stelle. `doctor` zuerst etabliert Envelope +
Schema-Helper auf Read-only-Boden mit **null** Architektur-
Last; `add` direkt danach validiert dann die schwerere
RecordingFileSystem-Decision auf bereits stabilem Envelope.
Der Use-Case-Druck "add zuerst" löst sich, sobald der
Cluster mit doctor anläuft — `add` ist der unmittelbar
nächste Slice, kein Slip. Plätze 3..9 unverändert nach
Use-Case-Druck aus dem Stub.

**Übergangs-Disziplin Root---json vs. existierendes
`template list --json`-Flag (Review-Round-2-Finding M3):**
Der erste Folge-Slice (`doctor`) führt das Root-PersistentFlag
`--json` aus T0-(a) ein. Sobald das Root-Flag landet, würde
das **lokale** `--json`-Flag auf `template list`
([`cli/template.go:83`](../../../../internal/adapter/driving/cli/template.go))
mit dem Root-PersistentFlag kollidieren (Cobra: Duplicate-
Flag-Registration oder Shadow-State). `template list` ist
laut Reihenfolge aber erst Platz 9 — der Output-Migration
auf den Envelope passiert dort, nicht früher. **Pflicht-
Schnitt im doctor-Slice (Platz 1):** das lokale `--json`-Flag
in `template list` wird **gleichzeitig mit der Root-Flag-
Einführung entfernt**, und der `runTemplateList`-Pfad liest
den `--json`-Flag-State von der Root-Persistent-Stelle (über
die App-Struktur oder `cmd.Root().PersistentFlags().GetBool(…)`).
**Output-Format bleibt unverändert** (heutige `templateJSON`-
Array-Struktur), nur das Flag-Wiring wandert. CLI-Pin-Tests
für `u-boot template list --json` **und** `u-boot --json
template list` müssen beide grün bleiben — gleicher Output,
gleicher Exit-Code. Die spätere Envelope-Migration auf Platz 9
(`slice-v1-cli-json-dry-run-template`) ersetzt dann den
`templateJSON`-Array-Output durch die Envelope-Form.

---

**ADR-0010-Liefer-Anker (post-T0):** ADR-0010 selbst bleibt
inhaltlich **unangetastet** (AGENTS.md §ADR-Disziplin: accepted
ADRs werden nicht umgeschrieben). Der Liefer-Anker für die
JSON-CLI-Spur ist dieser Slice (T0-Outcomes hier, T1-Schema-
Vertrag-Doku als nächster Tranchen-Schritt, Folge-Slices in
`open/`/`next/`/`done/`). Vor Cluster-Closure entscheidet
T_close, ob zusätzlich eine **neue Folge-ADR** (Nummer:
nächste freie nach ADR-Index-Stand) angelegt wird, die
ADR-0010 als Vorgänger referenziert und „JSON-CLI
ausgeliefert" als eigenständigen Entscheid trägt — oder ob
der Roadmap-Liefer-Vermerk und der done/-Slice ausreichen.
Diese Entscheidung gehört nicht ins T0.

## Resume-Punkt-Morgen (Session-Ende 2026-06-08)

Cluster-Stand **8/9 done** (config T8-Closure abgeschlossen, nach
`done/` verschoben), 1/9 open (template).

**Config-Slice 8/9 vollständig erledigt (T0–T8 + drei Review-
Runden)** — `make gates` grün (lint + test, Coverage 91.30 % ≥ 90 %,
docs-check). DoD-Hash-Tabelle im
[`done/slice-v1-cli-json-dry-run-config.md`](../done/slice-v1-cli-json-dry-run-config.md).
Zusammenfassung der Tranchen:

- **T2 (Port + CLI-Scaffold)**: Port-Felder + zwei Sentinels
  `ErrConfigWriteRejected` + `ErrConfigPostPatchSanityFailed`
  (T0-(m)-Split), `cli.ErrDryRunNotApplicable`,
  `configSetFlags.JSON`/`Quiet` read-through, Pin-Tests.
- **T3 (Application-Layer)**: Multi-`%w`-Migration der 5 FS-Wrap-
  Sites; Sentinel-Split-Wiring (`writeRejectedError` →
  `ErrConfigWriteRejected`; Post-Patch-Sanity inkl. der
  `revalidateFeatureEntry`-Sites → `ErrConfigPostPatchSanityFailed`;
  Coercion + Allowlist-Enforcement bleiben `ErrConfigValueInvalid`);
  SilenceLogger-Branch; Orphan-WARN als freie Funktion mit
  Dual-Emission. Tests `config_t3_test.go` + zwei WriteRejected-
  Test-Updates.
- **T4 (PreviewMode-Cluster + Composition-Root)**:
  `ConfigService.fsFactory`-Feld + nil-safe `selectFS(mode)`;
  `Set` + `setFeatureSourcesAllow` schreiben über
  `s.selectFS(req.PreviewMode)` (Reads auf `s.fs`);
  `NewConfigServiceWithFactory`-Konstruktor; `cmd/uboot/main.go`
  `configFSFactory` (jetzt fünf Factories) + Konstruktor-Wechsel.
  Tests `config_factory_test.go` (DryRun-touch-nichts, PreviewNone-
  persistiert, Legacy-ignoriert-Mode).
- **T5 (CLI-RunE)**: vollständige `config.go`-Neufassung — drei
  `runConfig*` auf Cluster-Signatur + JSON-Pfade; 3 Data-Carrier
  (`configShowData`/`configGetData`/`configSetData`); Subcommand-
  Pflicht auf ALLEN Envelopes inkl. Error-Pfad via neue subcommand-
  bewusste `reportErrorSub`/`writeErrorEnvelopeSub` (additiv in
  `erroremission.go`, Single-Form-Caller unverändert); Allowlist
  3 Forms (Reject 4→1); `mapConfigErrorToDiagnostic` Switch-Order
  T0-(f) FS-first; konsolidierter `configArgsValidator`; Voll-
  Schema/Dry-Run/Diff aus `resp.PlannedFiles`; bare/get-Reject via
  `ErrDryRunNotApplicable` + `isUsageError`-Branch (Exit 2);
  `SilenceLogger=flags.JSON`; WARN→`diagnostics[]`; `sanitizeBaseDir`.
  Acceptance-Tests `config_acceptance_test.go`.
- **T6 (Acceptance-Vervollständigung)**: white-box Mapper-Rows 1-10
  + Switch-Order-`_ByDesign` + Pure-FS-Exit-14 (`config_internal_
  test.go`); `--quiet --json` (3 Forms); Cobra-unknown-sub-Pin
  (bestätigt Parent-Validator-Dispatch); Help-Edge-Case; CONF-005-
  Disambiguation; Sanitizer-Worst-Case; Subcommand-Pflicht-Pin;
  Mid-Stage-Shapes. Coverage 91.20→91.30 %.
- **T7 (Review-Fix-Rounds)**: dritte Review-Runde — unabhängiger
  Agent über den CLI-Layer (T5/T6) → **R-CLI-1 (MED)**:
  `configArgsValidator` leakte ein Voll-Schema-Error-Envelope auf
  den Read-only-Forms bei `--dry-run`. Fix: Read-only-Validatoren
  hartkodieren `false/false` + zwei Regression-Pins. LOW-1:
  Diff-Hunks-Test auf echte Hunk-Inhalt-Assertion verschärft.

**Bewusste Tranchen-Verschiebungen** (alle dokumentiert + erledigt):
- **T2→T5**: `configGetFlags`/`configShowFlags` + `DryRun`/`Diff` +
  Flag-Registrierung — in T5 gelandet.
- **T3→T4**: PreviewMode-Handling — in T4 gelandet.
- **Plan-Refinement T3**: Post-Patch-Sanity-Split auf alle
  `revalidateFeatureEntry`-Sites ausgeweitet (Mapper-Row-6-Konsistenz).

**T8-Closure erledigt**: CHANGELOG `### Added`; `cli-json-output.md`
neue §6.9-Sektion (drei Sub-Form-Envelopes + Set-Voll-Schema +
Subcommand-Pflicht + Reject-Doku + Mapper-Tabelle + CONF-005-
Disambiguation + YAML-Comment-Limitation) + §6-Tabelle Rows 7+8
→done (logs-Row-Drift mitgefixt) + §7-Mutations-Zeile `config set`;
carveouts.md vier Folge-Stub-Einträge (R3-MED-3); Slice nach `done/`
mit DoD-Hash-Tabelle. Cluster-Stand jetzt **9/9 Folge-Slices done**.

**Folge-Slice 9/9 template — vollständig done (T0→T2→T4)**
([`done/`](../done/slice-v1-cli-json-dry-run-template.md)): der
**letzte und kleinste** Cluster-Slice (~60 LOC). T0-Discovery +
R1+R2+R3 (Asymptote HIGH 1→0→0) + T2 (`template list --json`
Array→Minimalkontrakt-Envelope, Breaking-Change CHANGELOG
`### Changed`; `mapTemplateErrorToDiagnostic`) + T4-Closure. **T3
entfiel** (bare-`template`-Reject + `cli.ErrTemplateSubcommandRequired`
nach Cluster-T_close verschoben — wäre solange das Gate existiert
toter Code). **Nächster Schritt: der finale Cluster-T_close** —
alle neun Folge-Slices sind in `done/`, die Closure-Hard-Rule ist
erfüllt: Allowlist-Mechanik + `PersistentPreRunE`-Abbau **+ bare-
`template`-RunE-Reject + Help-Leak-Pin**, optional Folge-ADR, dann
Cluster-Slice selbst nach `done/`.

**Zweistufiger T2–T4-Review (2026-06-08) vor T5 — zwei HIGH-Findings**:
- **R-T4-1** (Selbst-Review): T4-Recorder-Verzicht hätte dem
  T5-`--diff`-Pfad die Byte-Quelle entzogen; Fix
  `ConfigSetResponse.PlannedFiles`-Feld + Recorder-Surfacing wie add.
- **R-IR-1** (unabhängiger Reviewer-Agent, frischer Kontext): die
  T3-Split-Sentinels `ErrConfigWriteRejected`/`ErrConfigPostPatch
  SanityFailed` fehlten in `cli.go isConfigValidationError` — ein von
  der geplanten T5-Mapper-Tabelle **unabhängiger**, bereits live
  wirksamer ExitCode-Klassifikator. Folge: heutige Exit-10→Exit-1-
  Regression bei `config set services.x.enabled` + jedem Post-Patch-
  Sanity-Failure. Vom Selbst-Review übersehen (Fokus lag auf dem
  neuen Mapper). Gefixt + `TestExitCode_ConfigValidationSentinels`-
  Tabellen-Pin. Lehre: bei Sentinel-Splits IMMER beide Klassifikator-
  Pfade (Mapper + ExitCode) prüfen.
- Übrige Punkte clear: Multi-`%w` (kein Fremd-Sentinel-Match),
  Read/Write-Trennung (Reads auf `s.fs`, kein dry-run-Leak),
  SilenceLogger (alle 5 Sites), kein Mutex nötig (config swappt
  `s.fs` nie).

Session-Commits 2026-06-08:
- config-T2: Port-Felder + 2 Sentinels + CLI-Scaffold +
  Pin-Tests + Lifecycle `next/`→`in-progress/`.
- config-T3: Multi-`%w` + Sentinel-Split + SilenceLogger +
  Orphan-WARN→Warnings + Tests.
- config-T4: PreviewMode-Cluster + Composition-Root + Factory-Tests.
- config-T4-Review-Followup R-T4-1: `ConfigSetResponse.PlannedFiles`
  + Recorder-Surfacing.
- config-Review-Followup R-IR-1: ExitCode-Regression der zwei
  Split-Sentinels gefixt + ExitCode-Pin.
- config-T5: CLI-RunE-Neufassung + subcommand-bewusste Error-Envelopes
  + Allowlist + Mapper + Validator + Voll-Schema/Dry-Run/Diff +
  Reject + WARN-Mapping + Acceptance-Tests.
- config-T6: Acceptance-Vervollständigung (white-box Mapper-Rows +
  Switch-Order-`_ByDesign` + Cobra-unknown-sub + Help-Edge-Case +
  CONF-005-Disambiguation + Sanitizer + Subcommand-Pflicht +
  `--quiet --json` + Mid-Stage-Shapes).
- config-T7: unabhängiger CLI-Review → R-CLI-1 (MED) Args-Validator-
  Voll-Schema-Leak gefixt + 2 Regression-Pins + Diff-Hunks-Assertion
  verschärft.
- config-T8-Closure: CHANGELOG + `cli-json-output.md` §6.9/§6/§7 +
  carveouts (4 Stubs) + Slice nach `done/` + DoD-Hash-Tabelle
  (dieser Commit; DoD-Hash-Followup trägt den Closure-Hash nach).

## T_close — Detaillierter Umsetzungsplan (review-bereit, 2026-06-08)

Der **finale Schritt der Serie**. Alle 9 Folge-Slices sind in
`done/` → Closure-Hard-Rule erfüllt. T_close baut die Übergangs-
Mechanik ab, führt den verschobenen bare-`template`-Reject ein und
schließt den Cluster-Slice.

### Kern-Einsicht: atomare Kopplung

Gate-Abbau und bare-`template`-RunE-Reject **müssen zusammen
landen**. Solange das Allowlist-Gate existiert, rejected es bare
`template --json` (`ErrJSONNotImplemented`) **vor** der RunE — ein
RunE-Reject wäre toter Code. Erst mit dem Gate-Abbau wird der
RunE-Reject erreichbar + testbar (genau der Grund, warum template-T3
hierher verschoben wurde).

### Pre-Scan (Code-Realität — was abzubauen ist)

- [`jsonallowlist.go`](../../../../internal/adapter/driving/cli/jsonallowlist.go):
  `jsonAllowlist()`, `jsonRejectError()`, `jsonSliceSuffix()`,
  `applyJSONRejectGate()`, `helpRequested()` — die **gesamte
  Datei** ist nur die Gate-Mechanik.
- [`root.go`](../../../../internal/adapter/driving/cli/root.go) `PersistentPreRunE`
  (Z. 76-83): Gate-Call `applyJSONRejectGate(cmd, a.json)` (der
  Logger-Level-Teil bleibt).
- [`cli.go`](../../../../internal/adapter/driving/cli/cli.go):
  `ErrJSONNotImplemented` (Def Z. 216 + `isUsageError`-Eintrag
  Z. 481). **Verifiziert: keine externe Nutzung** (`grep` über
  Non-cli-Package = 0) → sicher entfernbar.
- [`template.go`](../../../../internal/adapter/driving/cli/template.go)
  Z. 56-58: bare-RunE = `cmd.Help()`.
- Tests: `jsonallowlist_test.go` (7 Test-Funktionen + `appendStubArgs`),
  `erroremission_internal_test.go` (`helpRequested`- + `jsonSliceSuffix`-
  Cases), `export_test.go` (`JSONSliceSuffixForTest`,
  `JSONAllowlistPathsForTest`, `WalkRootCommandPathsForTest`).

### Änderungen (file-by-file)

1. **template.go** — bare-RunE: `if a.json { return ErrTemplate
   SubcommandRequired }; return cmd.Help()` (Human-Modus
   unverändert). Neuer Sentinel `ErrTemplateSubcommandRequired`
   (Heim template.go, Pattern-Erbe logs `ErrFollowJSONNotSupported`).
   Doc-Block (Z. 47-55) aktualisieren.
2. **cli.go** — `ErrJSONNotImplemented` + dessen `isUsageError`-
   Eintrag entfernen; `ErrTemplateSubcommandRequired` in
   `isUsageError` (Exit 2).
3. **root.go** — Gate-Call (Z. 81-83) entfernen; PersistentPreRunE
   behält nur Logger-Level; Kommentare (Z. 57-62, 77-80) aktualisieren.
4. **jsonallowlist.go** — `git rm` (ganze Datei).
5. **export_test.go** — `JSONSliceSuffixForTest` +
   `JSONAllowlistPathsForTest` entfernen; `WalkRootCommandPathsForTest`
   **behalten** (vom repurposed Tree-Walk genutzt).
6. **Tests**:
   - `jsonallowlist_test.go` (`git mv` → `rootjson_test.go`):
     **entfernen** `RejectsAllNonMigratedForms`, `AcceptsDoctor`,
     `AllowlistAndTreeMatch`, `JSONSliceSuffix_StableMapping`,
     `AcceptsHelpFlag` (Gate-spezifisch/obsolet). **Behalten**
     `AcceptsTemplateList_BothFlagPositions`,
     `AcceptsTemplateList_FlagBeforeSubcommand` (reales Verhalten).
     **Repurposen** `TreeWalkAllowlistCompleteness` → „alle Forms
     antworten post-T_close, keiner leakt rohen Output" (kein
     `ErrJSONNotImplemented`-Ref mehr). **Neu**: bare-`template
     --json` → `ErrTemplateSubcommandRequired`/Exit 2/kein
     Help-Leak; bare-`template` Human-Modus → `cmd.Help()` (kein
     Reject).
   - `erroremission_internal_test.go`: `helpRequested`- +
     `jsonSliceSuffix`-Test-Cases entfernen.
7. **Public-Contract-Doku** (Review-MEDIUM — sonst driftet die
   Vertragsdoku gegen das CLI-Verhalten):
   - `docs/user/cli-json-output.md` **§6-Tabelle Zeile 9**: „bare
     `template --json` bleibt Gate-Reject bis Cluster-T_close" →
     „bare `template --json` → RunE-Reject
     `ErrTemplateSubcommandRequired`/Exit 2 (envelope-LOS §1838)".
   - **§6.1 „Übergangs-Reject für nicht-migrierte Forms"**: die
     Sektion beschreibt die **abgebaute** Gate-Mechanik
     (Allowlist + `PersistentPreRunE` + `ErrJSONNotImplemented`).
     Umschreiben auf: „alle Spec-Enum-Forms sind migriert; die
     Übergangs-Mechanik wurde im Cluster-T_close **entfernt**; die
     einzige verbliebene Reject-Form ist bare `template --json`
     (RunE-getragen, §1838-Ausnahme)." `ErrJSONNotImplemented`-
     Erwähnung raus.
   - **§6.2 bare-`template`-Absatz**: „Bis Cluster-T_close trägt
     der Allowlist-Gate diesen Reject (`ErrJSONNotImplemented`);
     mit dem Gate-Abbau übernimmt ein RunE-Reject" → „trägt ein
     RunE-Reject (`ErrTemplateSubcommandRequired`, envelope-LOS)".
   - **§ADR-Bezug** (~Z. 1178): „ob bei Cluster-T_close eine neue
     Folge-ADR angelegt wird, entscheidet der Cluster-Slice" →
     „Cluster-T_close hat entschieden: **keine** neue Folge-ADR
     (SD-1 (b)); der done-Slice + Roadmap-Liefervermerk
     dokumentieren die Auslieferung; ADR-0010 bleibt unverändert."
   - **CHANGELOG.md** (§Unreleased ### Changed): den Mechanik-
     Wechsel spiegeln — bare `template --json` wird jetzt per
     `ErrTemplateSubcommandRequired` (RunE) statt Gate
     (`ErrJSONNotImplemented`) rejected; die Allowlist-/Reject-
     Gate-Übergangs-Mechanik ist mit dem Cluster-Abschluss
     entfernt. (Entweder Erweiterung des bestehenden template-
     Eintrags oder eigener Closure-Eintrag.)

### Verifikation (Aufhebungsbedingung T_close)

- `make gates` grün (inkl. `docs-check` — fängt Vertragsdoku-
  Drift via Link-/Anchor-Checks; die §6.1/§6.2-Prosa wird im
  Review der genannten Abschnitte verifiziert).
- **Alle registrierten JSON-relevanten Cobra-Forms** (Spec-Enum
  inkl. der gruppierten Subcommands — `config`/`config get`/`config
  set`, `template`/`template list` zählen separat; Tree-Walk über
  `WalkRootCommandPathsForTest`): `--json` antwortet korrekt, kein
  Help-Leak, kein roher Output. Bare `template --json` →
  `ErrTemplateSubcommandRequired`/Exit 2 (expliziter Pin).
- `ErrJSONNotImplemented` nirgends mehr referenziert (`grep` == 0,
  inkl. Tests).

### Sub-Decisions (User-Review 2026-06-08 — festgezurrt)

- **SD-1 — Folge-ADR? → (b) KEINE Folge-ADR.** Nur die bestehende
  User-Doku (cli-json-output.md §ADR-Bezug) finalisieren: „T_close
  hat entschieden, keine neue ADR; Auslieferung via done-Slice +
  Roadmap-Vermerk". ADR-0010 bleibt unverändert.
- **SD-2 — Test-Redistribution → `git mv jsonallowlist_test.go →
  rootjson_test.go`** (behält Git-Historie).
- **SD-3 — Sentinel-Message → `"u-boot template requires a
  subcommand (try u-boot template list)"`** (Englisch korrekt,
  Exit 2 via `isUsageError`).
- **SD-4 — Commit-Struktur → zwei Commits**: erst Code-T_close
  (Mechanik-Abbau + Sentinel + Tests + Public-Doku) grün, dann
  Cluster-Slice → `done/` + DoD-Tabelle (+ DoD-Hash-Followup).

### Risiken

- **Coverage**: Entfernen von Gate-Code **plus** seiner Tests ist
  netto-neutral; der neue Reject + repurposed Tree-Walk decken das
  Neue. `make gates` (90%-Gate) fängt Rest-Lücken.
- **Negativer LOC-Slice**: ~ **-120 LOC** Mechanik+Tests, **+30 LOC**
  Sentinel+Reject+Pins → netto **~-90 LOC**. Der einzige Slice der
  Serie, der Code *entfernt*.
- **Irreversibilität**: nach Mechanik-Abbau gibt es kein Reject-Gate
  mehr; ein künftiger nicht-migrierter Subcommand müsste seinen
  eigenen `--json`-Pfad mitbringen (das ist gewollt — alle Spec-Enum-
  Forms sind migriert).

### Out-of-Scope für T_close

- Keine neuen `--json`-Features (alle Forms sind migriert).
- Keine Änderung an den 9 done-Slices (nur Mechanik + bare-template).

## Out of Scope

- **Reihenfolge der Folge-Slice-Implementierung als
  „MVP-First-Strategie"**: in T0-(e) wird zwar eine
  Reihenfolge nach Use-Case-Druck festgezurrt, aber
  **kein** Read-only-only- oder „MVP-Quorum"-Closure-Pfad
  für den Cluster — die Closure-Hard-Rule in der
  Aufhebungsbedingung schließt das strict aus: Cluster-
  Closure verlangt **alle 9 Folge-Slices in `done/`**.
  Der Slip-Notfall-Pfad (Carveout-Inventar + abgeschwächter
  Liefer-Vermerk im done/-Slice und in der Roadmap) ist
  **kein wählbarer** Closure-Pfad, sondern Restlauf-Disziplin
  nach einem bereits eingetretenen Slip. T0-(e) entscheidet
  nur, **in welcher Reihenfolge** die Folge-Slices angefasst
  werden, nicht welche „erstmal reichen".
- **JSON-Output für nicht-Spec-Enum-Subcommands** (z. B.
  zukünftige `u-boot exec`-Spec-Erweiterung): außerhalb dieses
  Cluster-Slices. Wenn neue Subcommands dazukommen, bekommen
  sie ihren `--json`-Pfad im selben Cluster-Stil — aber nicht
  in diesem Slice.
- **Schema-Versionierung** (`schemaVersion: 1` o. ä. im
  JSON-Output): Spec fordert nicht, also YAGNI. Triggern bei
  erstem Spec-Breaking-Change.
- **GraphQL / gRPC / WebSocket-Schnittstellen**: ADR-0010
  schließt explizit aus (HTTP-Adapter Re-Eval-Trigger).
  Re-Triggern dort, nicht hier.
- **Stream-Output für `logs --json`** als JSON-Lines-Format
  vs. Single-Envelope: Sub-Decision im
  `slice-v1-cli-json-dry-run-logs`-Folge-Slice, nicht hier.

## Bezug

- Spec: `LH-FA-CLI-007` (Dry-Run), `LH-FA-CLI-008` (Diff),
  `LH-NFA-USE-004` (Maschinen-lesbar) — alle V1
  ([`spec/lastenheft.md:302-489`](../../../../spec/lastenheft.md),
  [`spec/lastenheft.md:1809-1853`](../../../../spec/lastenheft.md)).
- ADR: [`ADR-0010`](../../adr/0010-kein-http-driving-adapter.md)
  §Entscheidung + §Folgepunkte Re-Eval-Trigger 2 — dieser
  Slice ist die JSON-CLI-Spur, die ADR-0010 voraussetzt.
- Vorbild-Code: existierender `template list --json`-Pfad
  ([`internal/adapter/driving/cli/template.go:143-156`](../../../../internal/adapter/driving/cli/template.go),
  `templateJSON`-DTO mit `[]`-Normalisierung,
  `json.MarshalIndent`). Erste Schema-Audit-Aufgabe: prüfen,
  ob dieser Pfad das Pflichtschema bereits trägt — vermutlich
  nein, deshalb `slice-v1-cli-json-dry-run-template` für
  Migration.
- Vorbild-Slice für T0-Outcomes-Layout:
  [`slice-v1-logs`](../done/slice-v1-logs.md) §T0-Outcomes
  (vier Sub-Decisions mit Begründung).
- Roadmap-Anker:
  [`roadmap.md`](../in-progress/roadmap.md) §v0.4.0-
  Arbeitspakete-Tabelle.
- Phase: V1 — V1-pünktlich notwendig, weil ADR-0010
  Re-Eval-Trigger 2 auf diesen Slice referenziert.
