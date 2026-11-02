# ADR-0042: Dauerhafte Lint-Ausnahmen — `contextcheck` für die CLI-Schicht, `interfacebloat` für den `driven.FileSystem`-Port, `!**/*_test.go` in den depguard-Regelblöcken

**Status:** Accepted

**Datum:** 2026-11-02

**Autor:** pt9912 (Rolle Planner — Werkzeug-Wahl nach dem Trichter, Baseline-Regelwerk `modul-07-carveouts.md` §Werkzeug-Wahl bei Diskrepanz, mit der Projektinhaber-Entscheidung 2026-11-02 zur Carveout-Ablage)

**Bezug:** [`LH-QA-004`](../../../spec/lastenheft.md#lh-qa-004--linting-solid-nahes-lint-profil), [`ADR-0003`](0003-solid-nahes-lint-profil.md), [MR-002](../../../harness/conventions/done/MR-002-carveout-inventar.md)

**Schärft:** [`SPEC-035`](../../../spec/spezifikation.md#spec-035--pflichten-des-lint-profils) — die Pro-Pfad-Carveouts der `exclude-rules` tragen je eine benannte Entscheidung; [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) — die konfigurierte Deckungsgleichheit der depguard-Regeln kennt die production-only-Ausnahme (begründete Klammer) — die Kennung nennen, wo das Zielelement eine trägt.

**Verhältnis:** **Ergänzt**, ersetzt nichts: [ADR-0003](0003-solid-nahes-lint-profil.md) bleibt `Accepted` und unverändert; diese ADR trägt drei dauerhafte Ausnahmen des Profils, die der Bestand von ADR-0003 nicht nennt (Fundstellen: die `Why:`-Kommentare der `issues.exclude-rules` in `.golangci.yml`; die Überführung folgt der Werkzeug-Wahl des Trichters mit der Projektinhaber-Entscheidung 2026-11-02 zur Carveout-Ablage — die dauerhaften Zeilen des früheren Master-Inventars werden über den Trichter in eine ADR überführt statt in CO-Dateien, denn „[f]ehlt der Folge-Slice, ist der Carveout *de facto* permanent … dann gehört er nicht in `carveouts/`, sondern über den Trichter … in eine ADR“, Baseline-Regelwerk `modul-07-carveouts.md` §Ziel-Form: Carveout).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lint-Profil steht in [ADR-0003](0003-solid-nahes-lint-profil.md) (29 Linter,
`//nolint` verboten, Ausnahmen zentral in `.golangci.yml` mit `Why:`-Kommentar).
Drei Ausnahmen des Bestands tragen dort keinen Namen:

1. **`contextcheck` für `internal/adapter/driving/cli/`:** Cobras
   `RunE`-Signatur (`func(cmd *cobra.Command, args []string) error`) hat
   keinen Context-Parameter; der Closure muss `cmd.Context()` extrahieren
   und durchreichen. `contextcheck` sieht die Closure-Grenze nicht und
   flaggt das als fehlende Propagation; die echte Propagation steigt eine
   Ebene tiefer in `runInit` (Context als erster Parameter). Fundstelle des
   Begriffs: der `Why:`-Kommentar des Ausschluss-Bocks in `.golangci.yml`
   (gelesen 2026-11-02).
2. **`interfacebloat` für `internal/hexagon/port/driven/filesystem.go`:**
   `driven.FileSystem` ist die zentrale FS-Abstraktion (12 Methoden:
   Exists/Lstat/ReadFile/WriteFile/WriteFileExclusive/Mkdir/MkdirAll/
   Rename/ReadDir/RemoveAll/Copy/CopyExclusive); sie bedient die
   Backup-Mechanik von [`LH-FA-INIT-005`](../../../spec/lastenheft.md#lh-fa-init-005--überschreibschutz) in ihrer ganzen Breite und
   speist die Test-Fakes aller Adapter. Eine künstliche Aufspaltung
   verkompliziert die Fakes ohne semantischen Mehrwert; das Limit (10) ist
   bewusst für genau diese eine Schnittstelle aufgeweicht (sonst globale
   `interfacebloat.max=10`, [ADR-0003](0003-solid-nahes-lint-profil.md)).
3. **`!**/*_test.go` als erstes files-Pattern in jedem depguard-Regelblock:**
   die Schicht-Regeln ([`LH-FA-ARCH-003`](../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement) via
   [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)) gelten production-only — Tests
   müssen Fakes und Test-Bibliotheken (`testify`, …) frei importieren
   können. (Die Zusage von [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) trägt die
   Ausnahme bereits wörtlich; die Zeile der Konfiguration wird hier
   benannt.)

Das frühere Master-Inventar `docs/plan/planning/in-progress/carveouts.md`
(Ablage der Adaption [MR-002](../../../harness/conventions/done/MR-002-carveout-inventar.md),
im Migrationsslice aufgelöst) führte diese drei unter „Permanente
Carveouts (kein Plan, im Inventar dokumentiert)“ — Zeilen ohne
Folge-Slice. Nach dem Trichter (`modul-07-carveouts.md` §Ziel-Form:
Carveout) ist der de-facto-permanente Carveout kein Werkzeug der
`carveouts/`-Ablage (siehe das Verhältnis oben). Die dauerhaften
Lint-Ausnahmen sind Architekturentscheidungen des Profils (keine
Übergangs-Zustände), ihre Trigger sind ehrlich nicht erfüllbar — der
Trichter führt hier auf diese ADR, nicht auf eine temporäre CO-Datei.

## Entscheidung

Wir wählen **eine gemeinsame ADR für die drei dauerhaften Lint-Ausnahmen**
(des Profils, Folgeentscheidung zu [ADR-0003](0003-solid-nahes-lint-profil.md)):

1. **`contextcheck` bleibt für `internal/adapter/driving/cli/` aus —**
   dauerhaft (Trichter: Trigger nie erreichbar, Cobras `RunE`-Signatur ist
   Framework-Fest; die strenge Propagation lebt in `runInit`). Die
   Konfiguration (`.golangci.yml`) trägt den `Why:`-Kommentar und nennt
   diesen ADR.
2. **`interfacebloat` bleibt für `internal/hexagon/port/driven/filesystem.go` aus —**
   dauerhaft (die Breite des Ports ist eine Architektur-Zusage der
   Backup-Mechanik und der Fake-Stabilität, keine Übergangslage; ein
   Umbau der Schnittstelle wäre ein eigener Architektur-Vorgang).
3. **`!**/*_test.go` bleibt erstes files-Pattern jedes depguard-Regelblocks —**
   dauerhaft (die Schicht-Regeln gelten production-only,
   [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)).
   Die Zeile ist von [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln) getragen und wird hier
   explizit benannt.

Die drei gehören zu diesem einen ADR, weil sie dieselbe Natur sind:
dauerhafte Ausnahmen des Lint-Profils, die keine Gates lockern — Schwellen
bleiben ([`AGENTS.md`](../../../AGENTS.md) §3.6, die Schwelle selbst bleibt);
der Review-Kontext für jede dieser Ausnahmen ist die gleiche Lesart
(„bewusste Profil-Entscheidungen, keine Suppressions“ — die Lesart der Pflicht-Carveouts, docs/maintainer/quality.md §1.3).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — drei temporäre CO-Dateien mit Aufhebungs-Trigger | die Baseline-Mechanik (`modul-07-carveouts.md`) greift wörtlich; das Audit-Slice prüft | die Trigger sind ehrlich nie erreichbar (Framework-Grenze/Cobras-Signatur, Port-Breite als Zusage, production-only-Regel) — jeder CO-Block wäre ein permanenter Carveout (der lügt), genau das, was Modul 7 ausschließt; drei Wächter-Artefakte ohne Aufhebungs-Weg |
| B — je Ausnahme ein eigener ADR | je Entscheidung ihre Datei; feiner Verweis | drei ADR-Dateien mit je Kontext/Alternativen/Konsequenzen für dieselbe Klasse; kein eigener Inhalt — der gemeinsame Kontext (Lint-Profil) wäre dreifach gehalten; Drift-Risiko gegen ADR-0003 |
| **C — ein gemeinsamer ADR für die drei (gewählt)** | ein Ort für dieselbe Klasse; je Ausnahme eine nummerierte Entscheidung (§Entscheidung); Aufhebung einer einzigen über `Supersedes` in einer Folge-ADR | eine Aufhebung einer einzigen Ausnahme muss dann den ganzen ADR supersedieren (Muster `modul-04-adrs.md`) |
| D — Einträge in ADR-0003 nachtragen | ein Ort | [ADR-0003](0003-solid-nahes-lint-profil.md) ist `Accepted` und damit unveränderbar ([`AGENTS.md`](../../../AGENTS.md) §3.5); ein Nachfassen bricht die Unveränderlichkeit |
| E — nichts tun: die Ausnahmen bleiben Kommentar-getragen | kein neues Artefakt | die Ausschluss-Liste hat keinen benannten Entscheidungs-Anker; der `Why:`-Kommentar wäre die einzige Begründung — eine Ausnahme ohne ADR-Anker widerspricht [`AGENTS.md`](../../../AGENTS.md) §3.6 (Gates dürfen nicht ohne ADR/Plan-Anker gelockert werden); bei der neuen Carveout-Ablage (`docs/plan/carveouts/`) tragen die dauerhaften Einträge sonst gar keinen Anker |

## Konsequenzen

- **Positiv:** die drei Ausnahmen tragen einen Entscheidungs-Namen;
  ADR-0003 bleibt unangetastet; der README-Index der Carveout-Ablage
  (`docs/plan/carveouts/`) zeigt an einen Ort statt acht Zeilen zu duplizieren;
  jeder neue dauerhafte Lint-Ausschluss hat eine Form (hier nachtragen im
  Sinne dieses Musters oder Trichter).
- **Negativ / Trade-offs:** eine ADR-Datei führt drei Ausnahmen — eine
  Aufhebung einer einzigen (etwa ein Redesign des `FileSystem` Ports)
  berührt den gemeinsamen Kontext; die ADR bleibt `Accepted` und die
  Aufhebung braucht dann eine Folge-ADR mit `Supersedes` (Muster §3.5).

## Fitness Function

| Tooling | Regel | Make-Target |
|---|---|---|
| golangci-lint `issues.exclude-rules` | die drei Ausschlussregeln stehen konfiguriert, mit dem `Why:`-Kommentar, der diese ADR nennt; ohne die drei läuft das `lint`-Gate auf den ausgeschlossenen Funden rot | `make lint` |

## Re-Evaluierungs-Trigger

**permanent** — die drei Ausnahmen haben keinen erreichbaren
Auflösungs-Trigger (Cobras `RunE`-Signatur ist Framework-Fest; die
Port-Breite ist eine Architektur-Zusage; die production-only-Regel ist die
Zusage selbst, [`SPEC-034`](../../../spec/spezifikation.md#spec-034--enforcement-der-import-regeln)). Einer Neubewertung dient der
folgende Ereignis-Trigger: ein Redesign der
`driven.FileSystem`-Schnittstelle oder ein Frameworkwechsel der
CLI-Schicht (Cobra → ein anderes Driving-Modell) — dann eine Folge-ADR mit
`Supersedes ADR-0042`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-11-02 | Entwurf, `Proposed` (Planner), auf die Anweisung und Werkzeug-Wahl des Projektinhabers zur Carveout-Ablage vom selben Tag | Vereinbarung mit dem Projektinhaber |
| 2026-11-02 | `Accepted` durch den Projektinhaber mit derselben Werkzeug-Wahl (dauerhafte Ausnahmen über den Trichter in eine ADR) | Commit des Migrationsslices |