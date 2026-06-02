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
> Per-Command-Folge-Slice-Serie. T0 (Discovery + Sub-Decisions)
> festzuhalten beim Übergang nach `next/`; danach Spawn der
> ersten Per-Command-Slices.

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
(2) das ADR-0010-Konsequenzen-Update darf **nicht** „JSON-CLI
als Maschinen-Schnittstelle ausgeliefert" sagen, sondern
muss den Carveout zitieren und Re-Trigger-Pfad nennen.

Default-Erwartung: keine Slips, alle 9 Folge-Slices schließen.
Die Hard-Rule ist Notfall-Pfad, nicht Standard.

## Akzeptanzkriterien (Cluster-Ebene)

- ✅ **Schema-Vertrag dokumentiert**: zentraler
  Reference-Block (vermutlich `docs/user/cli-json-output.md`
  oder Sektion im Architecture-Doc) zitiert das `LH-FA-CLI-
  007`-Schema verbatim und benennt die DTO-Lokalisation-
  Konvention (Per-Subcommand-DTO im Driving-Adapter,
  Common-Fields-Helper z. B. `cliJSONEnvelope`).
- ✅ **Per-Command-Folge-Slices angelegt** für alle zehn
  Subcommands, jeder mit eigenem T0-Discovery und LOC-
  Schätzung. Reihenfolge nach Use-Case-Druck festgezurrt
  (siehe §T0-Discovery (e) unten).
- ✅ **Erster Folge-Slice abgeschlossen** als belastbares
  Pattern-Vorbild für die restlichen (z. B. `add --dry-run
  --diff --json` als modifying-Pilot oder `doctor --json` als
  read-only-Pilot — T0-(e) entscheidet).
- ✅ **Schema-Konformitäts-Helper** im CLI-Adapter (oder als
  Test-Helper in `internal/adapter/driving/cli/jsontestutil/`):
  parst die `--json`-Ausgabe und prüft Pflichtfelder. Jeder
  Folge-Slice verwendet ihn in seinen Tests, damit Schema-
  Drift einheitlich kracht.
- ✅ **ADR-0010-Kreuzverweis** in §Konsequenzen aktualisiert:
  „JSON-CLI als kanonische Maschinen-Schnittstelle (sobald V1
  ausgeliefert)" wandert von „prospektiv" auf „ausgeliefert per
  `slice-v1-cli-json-dry-run`" sobald alle Folge-Slices in
  `done/` sind.
- ✅ **Roadmap-Status** zeigt den Cluster-Fortschritt: jede
  Folge-Slice-Closure aktualisiert die v0.4.0-Tabelle in
  `roadmap.md`.

## T0-Discovery (vor `next/`-Übergang festzulegen)

Sub-Decisions, die der Cluster-Slice klären muss, bevor die
ersten Per-Command-Folge-Slices spawnen. Jede Sub-Decision wird
beim `open/ → next/`-Lifecycle in ein `T0-Outcomes`-Layout
(analog [`slice-v1-logs`](../done/slice-v1-logs.md) §T0-Outcomes)
gegossen.

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
| T0 | **Discovery + Sub-Decisions.** Fünf T0-Fragen aus §T0-Discovery klären (Flag-Scope, Dry-Run-Architektur, DTO-Lokalisation, Diff-Renderer, Reihenfolge). Entscheidung pro Frage mit Begründung in einem `T0-Outcomes`-Block dokumentieren. ADR-0010-Kreuzverweis aktualisieren. | — (Plan-Arbeit) |
| T1 | **Schema-Vertrag-Doku.** `docs/user/cli-json-output.md` (neu) zitiert `LH-FA-CLI-007`-Schema verbatim, dokumentiert DTO-Konvention, listet Per-Command-Folge-Slice-Reihenfolge. README EN+DE bekommt Verweis-Zeile. | ~80 (reine Doku) |
| T2..Tn | **Spawn Folge-Slice-Stubs** in `open/` für alle 9 Per-Command-Slices. Pro Stub: Auslöser + grobe AKs + LOC-Schätzung + Verweis auf gemeinsamen Schema-Vertrag. Reihenfolge nach T0-(e). | ~30 LOC pro Stub × 9 = ~270 |
| T_close | **Cluster-Closure.** Sobald ein „kritisches Quorum" der Folge-Slices in `done/` ist (siehe §Out of Scope für Definition), Cluster-Slice mit DoD-Hash-Line aller Folge-Slices nach `done/` mit Roadmap-Update + ADR-0010-Konsequenzen-Update. | — (Doku) |

LOC-Schätzung Cluster-Slice: ~350 LOC, deutlich unter
800-LOC-Schwelle. Folge-Slice-LOC-Bandbreite: 200..600 je
Subcommand (T0-(b)-Architektur-Decision dominant).

## Out of Scope

- **Reihenfolge der Folge-Slice-Implementierung als
  „MVP-First-Strategie"**: in T0-(e) wird zwar eine
  Reihenfolge nach Use-Case-Druck festgezurrt, aber
  **kein** Read-only-only- oder „MVP-Quorum"-Closure-Pfad
  für den Cluster — die Closure-Hard-Rule in der
  Aufhebungsbedingung schließt das aus (alle 9
  Folge-Slices done/, sonst Carveout-Inventar-Pflicht
  vor Cluster-Move). T0-(e) entscheidet nur, **in welcher
  Reihenfolge** die Folge-Slices angefasst werden, nicht
  welche „erstmal reichen".
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
