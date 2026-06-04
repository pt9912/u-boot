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
> in `in-progress/`, Folge-Slice 1/9
> [`slice-v1-cli-json-dry-run-doctor`](../open/slice-v1-cli-json-dry-run-doctor.md)
> in `open/`. Nächster Schritt: T0-Discovery des Doctor-Slices
> (8 Sub-Decisions), danach `next/`-Übergang dort.

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
| T1 | **Schema-Vertrag-Doku.** `docs/user/cli-json-output.md` (neu) zitiert `LH-FA-CLI-007`-Schema verbatim, dokumentiert DTO-Konvention, listet Per-Command-Folge-Slice-Reihenfolge. README EN+DE bekommt Verweis-Zeile. **Delegation:** Doctor-Slice T0-(f) übernimmt diese Tranche (Lieferung als T1 in [`slice-v1-cli-json-dry-run-doctor`](../open/slice-v1-cli-json-dry-run-doctor.md)) — Cluster-T_close-Reviewer sucht **keinen** separaten Cluster-T1-Commit, sondern verweist auf den Doctor-Slice-T1-Commit. | ~80 (reine Doku, via Doctor-Slice) |
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
