# Harness-Konventionen - u-boot

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln von `u-boot`
gegenueber der adoptierten Harness-Konvention (Baseline). Sie ist der
Default-Ort fuer:

- **Adaptionen** gegenueber der Baseline (mit Begruendung und Aufloesungs-Trigger).
- **ID-Schema-Deklaration** - welches Praefix-Schema dieses Repo nutzt.
- **Zusatzklassen-Deklarationen** fuer repo-spezifische Bindung-Klassen der
  Sensors-Tabelle, die ueber die kanonischen hinausgehen.
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield / Hybrid).

Sie **dupliziert keinen Baseline-Text** - sie verweist und ergaenzt. Bei
Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die kanonische
Quelle (Source Precedence). Diese Datei ist konformitaets-bringend fuer
*Form*-Fragen, nicht autoritativ ueber Inhalt.

## Baseline

- **Konvention:** AI-Harness-Kurs (`pt9912/ai-harness-course`)
- **Stand:** v3.5.1 (Regelwerk-Bundle; Kurs-Welle 33)
- **Datum der Adoption:** 2026-07-24 (Erst-Adoption direkt auf `v3.5.1`).
- **Integritaets-Pin:** `.harness/baseline/v3.5.1/SHA256SUMS` ueber den vendorten
  Bestand (`regelwerk/` + `templates/`); offline pruefbar per
  `tools/harness/fetch-baseline-cache.sh --verify`.

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):** <https://github.com/pt9912/ai-harness-course/tree/v3.5.1>
- **Regelwerk (committet vendored, `MR-004`/`MR-007`):** die Lese-Form ist das
  nach Modulen und Grundlagen-Abschnitten aufgeteilte Bundle, entpackt und
  committet unter `.harness/baseline/v3.5.1/regelwerk/` (Index
  `regelwerk/README.md`), samt `.harness/baseline/v3.5.1/SHA256SUMS` - netzlos
  auf jedem Checkout, offline verifizierbar. Bundle-Quelle: Release-Asset
  `lab-regelwerk.zip`, Tag `v3.5.1`.
- **Templates (committet vendored, `MR-004`):** die Skelett-Vorlagen liegen
  vendored unter `.harness/baseline/v3.5.1/templates/` (aus demselben Bundle)
  und tragen zwei Rollen: **Referenz-Form**, auf die das Regelwerk mit
  `../templates/...` als "Ziel-Form" verweist (netzlos, weil parallel zu
  `regelwerk/` vendored), und **Kopiervorlage** - beim Anlegen neuer Artefakte
  (ADR, Slice, Welle, Carveout, Review-Report) das passende Template **kopieren
  und ausfuellen** statt frei zu formulieren.
- **In-Repo (verkoerperte Form):** die Gate-Baseline (`.d-check.yml`, `Makefile`,
  `Dockerfile`, `.golangci.yml`, `.github/workflows/`) und die
  autoren-gepflegte Harness-Prosa unter `harness/` (`README.md`, `roles.md`,
  `review.md`, `verification.md`, `replay.md`, diese Datei).

## Sync-Trigger (T1/T2)

Begriffe aus dem Regelwerk (`modul-02-harness-bootstrap`): Ein Pointer auf die
vendored Baseline muss an **zwei** Stellen synchron gehalten werden.

- **T1** - Pointer in [`harness/README.md`](README.md) Abschnitt Guides (Verweis auf das
  vendored Regelwerk und diese Datei).
- **T2** - Pointer in der Source-Precedence-/Kopf-Sektion von
  [`AGENTS.md`](../AGENTS.md) (Verweis auf die vendored Baseline + Lesemodell).

Beide zeigen auf `.harness/baseline/v3.5.1/regelwerk/README.md` (Index) und
werden bei einem Baseline-Bump gemeinsam nachgezogen (`MR-004` Bump-Prozedur).
Fundstelle: `.harness/baseline/v3.5.1/regelwerk/modul-02-harness-bootstrap.md`.

## Freshness-Audit

Der Integritaets-Pin (`SHA256SUMS`) beantwortet *"ist der vendorte Bestand
unversehrt?"*. Das Freshness-Audit beantwortet die andere Frage: *"ist der
gepinnte Stand noch der aktuelle?"*. Ohne den zweiten Sensor altert die
Baseline still.

- **Sensor:** `tools/harness/fetch-baseline-cache.sh --check-freshness` -
  liest die Release-**Liste** von `pt9912/ai-harness-course` und vergleicht sie
  mit dem `**Stand:**`-Pin oben. Read-only: kein Vendoring, kein Pin-Update,
  kein Schreibzugriff auf `.harness/baseline/`.
- **Exit-Codes:** `0` = Pin ist der neueste Tag; `3` = neuerer Tag vorhanden
  (Review-Bump faellig); `1` = Ausfuehrungsfehler (Netz, Werkzeug, Format,
  Pin nicht in der Liste). Ein Fehler ist nie ein stilles "alles aktuell" -
  Fail-loud ist Pflicht, weil ein Falsch-Negativ die Routine wertlos macht.
- **Kadenz:** (a) **ereignisgetrieben** - jeder Harness-/Baseline-Slice fuehrt
  den Check aus und haelt sein Ergebnis in der Verification Evidence fest,
  auch den Negativbefund; (b) **kalendarisch** - mindestens quartalsweise,
  falls in dem Zeitraum kein Harness-Slice lief. Die Ereignis-Kadenz ist die
  tragende; die kalendarische ist nur das Sicherheitsnetz gegen lange Pausen.
- **Zustaendig:** die Rolle, die den jeweiligen Harness-/Baseline-Slice fuehrt.
- **Auslöser bei Befund:** ein **Review-Bump** - ein eigener Slice, der die
  Aenderungen zwischen den Tags liest und den Bump als Einheit ausfuehrt
  (`MR-004` Bump-Prozedur: Pin, Vendor-Pfad, `AGENTS.md`, `harness/README.md`).
  **Kein Auto-Update**: Ein Regelwerks-Wechsel kann Adaptionen dieses Ledgers
  ungueltig machen; das ist eine Lese- und Entscheidungsleistung, nichts, was
  ein Skript still tut.
- **Nicht-Ziele:** kein Eintrag in `make gates` oder `make ci` - ein
  Kurs-Release darf u-boots Pipeline nicht rot faerben; kein automatischer
  Vendor-Lauf; keine Pin-Mutation durch den Check.

## Adaptions-Block

### MR-000 - Baseline-Aussage

- **Datum:** 2026-07-24
- **Geltungsbereich:** gesamtes Repo
- **Adaption:** *keine inhaltlichen Adaptionen gegenueber Baseline-Default fuer
  Verzeichniskonvention, Lifecycle-Regeln (`open` -> `next` -> `in-progress` ->
  `done`) und die etablierten ID-Schemata* (`ADR-<NNNN>`, `LH-*` samt
  Verifikations-Aliassen `PH-*`/`TC-*`, `slice-<phase>-<slug>`,
  `tranche-<nr>-<slug>`, Carveout-Familie `CO-*`). Konkrete Abweichungen sind als
  eigene `MR-<NNN>` unten dokumentiert.
- **Begruendung:** Initial-Setzung. u-boot war vor der Adoption bereits
  harness-geformt; dieser Block haelt den konformen Grundstand fest, spaetere
  Adaptionen folgen als `MR-<NNN>`.
- **Aufloesungs-Trigger:** permanent.

### MR-001 - Source Precedence mit zwei Spec-Straten, ohne Technik-Stratum

- **Datum:** 2026-07-24
- **Geltungsbereich:** [`AGENTS.md`](../AGENTS.md) Abschnitt Source Precedence,
  [`harness/README.md`](README.md) Abschnitt Source Precedence.
- **Adaption:** u-boot fuehrt eine 9-Rang-Source-Precedence mit **zwei**
  Spec-Straten: `contract_spec` ([`spec/lastenheft.md`](../spec/lastenheft.md),
  vertraglich) und `view_spec`
  ([`spec/architecture.md`](../spec/architecture.md), Sicht). Ein **Technik-
  Stratum** (`spec/spezifikation.md`) wird bewusst **nicht** gefuehrt.
  Repo-Klasse: **Tooling/Referenz**.
- **Begruendung:** Das Technik-Stratum ist laut Regelwerk (`grundlagen-
  konventionen` Abschnitt Spec-Straten) optional; nur Vertrag und Sicht sind
  obligatorisch. u-boots Zwei-Straten-Klassifikation ist dort selbst als
  Referenz-Implementierung benannt. Kein dritter Stratum -> keine Luecke.
- **Aufloesungs-Trigger:** permanent, solange u-boot ohne separates
  Technik-Stratum auskommt.

### MR-002 - Carveout-Inventar an fester Stelle statt `docs/plan/carveouts/`

- **Datum:** 2026-07-24
- **Geltungsbereich:** Carveout-Ablage;
  [`docs/plan/planning/in-progress/carveouts.md`](../docs/plan/planning/in-progress/carveouts.md),
  [`AGENTS.md`](../AGENTS.md) Abschnitt Planning-Lifecycle, `.d-check.yml` `matrix`.
- **Adaption:** Carveouts werden **inventarisiert** in der einen Datei
  `docs/plan/planning/in-progress/carveouts.md` (mit Plan-Anker je Eintrag),
  nicht als je eine Datei `docs/plan/carveouts/CO-<NNN>-<titel>.md`. Das
  `CO-*`-Schema bleibt fuer die ID-Vergabe gueltig.
- **Begruendung:** Etablierte u-boot-Struktur, bereits in `AGENTS.md` und der
  d-check-`matrix` (Klasse `carveout`) verankert; eine zusaetzliche
  Ein-Datei-pro-Carveout-Ebene braechte keinen Mehrwert und erzeugte Drift.
- **Aufloesungs-Trigger:** permanent, solange Carveouts zentral inventarisiert
  werden.

### MR-003 - Roadmap folgt Wellen-Template; Release-Versionen = Wellen; Ort in-progress/

- **Datum:** 2026-07-24
- **Geltungsbereich:**
  [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md),
  [`docs/plan/planning/README.md`](../docs/plan/planning/README.md).
- **Adaption:** Die Roadmap folgt der v3.5.1-`roadmap.template.md`-Struktur
  (Aktuelle Welle, Naechste Wellen, Meilensteine, Abhaengigkeitsgraph,
  Abgeschlossene Wellen, Historische Trigger-Verschiebungen). u-boots
  **Release-Versionen sind die Wellen**; Termine erscheinen nur als *Konsequenz*
  einer abgeschlossenen Welle (Release-Datum), nicht als Treiber (Template-
  Format-Regel "Wellen, keine Termine" damit gewahrt). Zwei Orts-/Form-
  Abweichungen: (a) die Roadmap liegt unter
  `docs/plan/planning/in-progress/roadmap.md` (nicht
  `docs/plan/planning/roadmap.md` - das Regelwerk nennt beide Pfade; u-boot loest
  zugunsten des Lifecycle-Verzeichnisses auf); (b) es gibt **keine**
  eigenstaendigen `welle-NN-results.md` und keine Wellen-Plan-Dateien - die
  Welle-Closure lebt im jeweiligen `done/`-Release-Cut-Slice (Detailquelle),
  Wellen sonst inline in der Roadmap.
- **Begruendung:** Die Template-Struktur macht Wellen-Reihenfolge, Trigger und
  Abhaengigkeiten explizit. u-boot liefert dated Releases, aber scope-getrieben
  (Datum = wann Scope fertig war, kein Deadline) - kompatibel mit der
  Template-Format-Regel. Die `welle-NN-results.md`-Ebene entfaellt, weil der
  `done/`-Release-Cut-Slice die Closure bereits vollstaendig traegt.
- **Aufloesungs-Trigger:** permanent, solange Release-Versionen die Wellen sind.

### MR-004 - Regelwerk-Lese-Form committet vendored; Baseline-Pin v3.5.1; beide Baeume

- **Datum:** 2026-07-24
- **Geltungsbereich:** `.harness/baseline/`,
  `tools/harness/fetch-baseline-cache.sh`, [`AGENTS.md`](../AGENTS.md) Abschnitt
  "Betriebsregelwerk (vendored Baseline)",
  [`harness/README.md`](README.md) Abschnitt Guides, [`.d-check.yml`](../.d-check.yml)
  (`scan.ignore`), `.gitignore`, Abschnitt Baseline oben.
- **Adaption:** Die Lese-Form des adoptierten Regelwerks ist **committet
  vendored** (kein Remote-ZIP pro Lauf, kein Submodule):
  `.harness/baseline/v3.5.1/{regelwerk,templates}/` + `SHA256SUMS`, netzlos auf
  jedem Checkout, offline verifizierbar. u-boot vendored **beide** Baeume
  (Upstream-Default), damit die `../templates/...`-Verweise der Module netzlos
  aufloesen und die Templates als Kopiervorlage bereitstehen - **kein**
  Consumer-Ausschluss der Templates.
- **Aufloesungs-Trigger / Bump-Prozedur:** Der `**Stand:**`-Pin ist nur die
  **Skript-Eingabe**, **nicht** vollumfaenglicher Single Source of Truth. Ein
  Versions-Bump ist als Einheit auszufuehren und fasst mindestens vier Stellen an:
  (1) `**Stand:**`-Pin oben, (2) Vendor-Pfad `.harness/baseline/<tag>/`
  (Skript-Lauf), (3) `AGENTS.md`-Pointer, (4) `harness/README.md`-Guides-Zeile.
  Ein neuer Kurs-Tag wird ueber die Release-**Liste** erkannt und loest einen
  Review-Bump aus, keinen Auto-Update. Sensor, Exit-Codes, Kadenz und
  Zustaendigkeit stehen im Abschnitt Freshness-Audit oben
  (`--check-freshness`, seit 2026-07-25 ausfuehrbar statt nur zugesagt).

### MR-005 - Gate-Haltung: `docs-check` via direktem Container-Lauf; `scan.ignore` erweitert

- **Datum:** 2026-07-24
- **Geltungsbereich:** [`.d-check.yml`](../.d-check.yml), [`Makefile`](../Makefile)
  (`docs-check`).
- **Adaption:** `make docs-check` laeuft via direktem
  `docker run ... $(D_CHECK_IMAGE)` gegen die repo-lokale `.d-check.yml` (kein
  tool-generiertes `d-check.mk`-`--print-mk`-Fragment). Aktive Module:
  `[links, anchors, ids, matrix]`. Bewusst **kein** `MR-<NNN>`-ID-Pattern - die
  Adaptions-IDs dieses Ledgers bleiben linkfrei. `.harness/baseline/**` liegt im
  `scan.ignore` (tag-agnostischer Glob `**`), damit die repo-relativen Links der
  vendorten Regelwerk-/Template-Dateien nicht gewertet werden.
- **Begruendung:** Der direkte Container-Lauf ist die etablierte u-boot-Form
  (digest-gepinntes `D_CHECK_IMAGE`); die Modul-Auswahl deckt den bestehenden
  Doku-Referenz-Vertrag. Der `scan.ignore`-Glob verengt **nicht** auf einen Tag,
  damit kuenftige vendored Staende automatisch erfasst sind; `.harness/skills/`
  (`MR-009`) bleibt ausserhalb des Baseline-Globs und damit pruefbar.
- **Aufloesungs-Trigger:** permanent; Modul-Auswahl bei d-check-Upgrade
  re-evaluieren.

### MR-006 - Modus-Deklaration pro Sub-Area

- **Datum:** 2026-07-24
- **Geltungsbereich:** Abschnitt Modus-Deklaration unten.
- **Adaption:** u-boot traegt Bestandscode (`hexagon/`, `cmd/`, `internal/`) neben
  den Doku-Sub-Areas. Der Abschnitt Modus-Deklaration unten ordnet jede Sub-Area
  als GF/BF/Hybrid ein; jede BF-/Hybrid-Markierung traegt eine
  Graduation-Bedingung. **Audit ausgefuehrt (2026-07-25):** Der Erst-Pass
  (drei grobe Sub-Areas, pauschal GF) ist durch eine auditierte Einordnung
  ersetzt - Drei-Achsen-Inklusion je Kandidat, dann vier Modus-Kriterien je
  qualifizierter Sub-Area. Ergebnis: acht Sub-Areas statt drei, davon eine
  Hybrid (`internal/adapter/driving/cli`) und eine Brownfield
  (`internal/**/README.md`); `cmd/uboot` faellt auf Sub-Area-Aspirantin zurueck.
- **Begruendung:** Das Regelwerk verlangt eine Modus-Aussage pro qualifizierter
  Sub-Area; eine BF-Sub-Area ohne Graduation-Plan waere "permanente Ausnahme als
  temporaer getarnt".
- **Aufloesungs-Trigger:** Audit erledigt; Delivery-Verweis im Slice
  [`slice-harness-sub-area-modus-audit`](../docs/plan/planning/open/slice-harness-sub-area-modus-audit.md)
  §9. Re-evaluieren bei jeder neuen Pfad-Familie im Produktivcode sowie beim
  Erfuellen einer der beiden Graduation-Bedingungen oben.

### MR-007 - Ortswahl `.harness/` (dot-prefixed, committet) neben `harness/`

- **Datum:** 2026-07-24
- **Geltungsbereich:** `.harness/` (vendored Baseline, kuenftig
  `.harness/skills/`), `harness/` (Autoren-Prosa), `.gitignore`.
- **Adaption:** Maschinen-materialisierte / vendorte Harness-Artefakte liegen im
  **dot-prefixed, aber getrackten** `.harness/` (Baseline unter
  `.harness/baseline/<tag>/`; Skills-Dateien kuenftig unter `.harness/skills/`).
  Die **handgeschriebene** Harness-Vertragsdoku bleibt im getrackten `harness/`
  (`README.md`, `roles.md`, `review.md`, `verification.md`, `replay.md`, diese
  Datei). `.harness/baseline/**` ist **committet** - bewusste Ausnahme zur
  "Dot-Prefix = ignorieren"-Lesart; nur ephemere Nebenprodukte
  (`.harness/cache/`) werden ignoriert.
- **Begruendung:** Der Dot-Prefix haelt den vendorten/generierten Bestand optisch
  vom Autoren-Bestand getrennt und folgt der Kurs-Oekosystem-Konvention
  (`.harness/baseline/`). Die Trennung ist hier explizit dokumentiert, damit der
  Zwei-Verzeichnis-Split (`.harness/` vs. `harness/`) kein
  Verwechslungs-Fallstrick ist.
- **Aufloesungs-Trigger:** permanent.

### MR-008 - ADR-Form: Bestand lean+grandfathered, neue ADRs MADR (per CR)

- **Datum:** 2026-07-24
- **Geltungsbereich:** ADR-Form-Politik; verweist auf den CR-Slice
  [`slice-cr-adr-format-madr`](../docs/plan/planning/done/slice-cr-adr-format-madr.md).
  Aendert `spec/lastenheft.md` NICHT von hier aus.
- **Adaption:** Das vendored ADR-Template (MADR-/Nygard-Stil) kollidiert mit dem
  heutigen [`LH-FA-PROJDOCS-002`](../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  (Status/Datum als Inline-Felder statt `##`-Ueberschriften; Titel-/Superseded-
  Format; zusaetzliche Pflicht-Sections). Die Angleichung aendert das
  **Vertrags-Stratum** und ist damit ein **Change Request**, nicht per
  conventions-MR moeglich. Die Aenderung traegt der CR-Slice, nicht dieser Block.
  **CR ausgefuehrt (2026-07-24):** `LH-FA-PROJDOCS-002` traegt jetzt die
  MADR-Form; die zum CR-Zeitpunkt Accepted ADRs (`0001`-`0010`, `0013`) bleiben
  lean + immutabel (grandfathered, Hard Rule); Proposed (`0011`, `0012`) und alle
  neuen ADRs sind MADR-konform.
- **Begruendung:** conventions.md ist form-bringend, nicht vertrags-aendernd; eine
  Vertrags-Anforderung darf hier nur referenziert werden. Das Template-Feld
  "Schaerft" deckt sich mit u-boots Referenzmodell
  ([`ADR-0013`](../docs/plan/adr/0013-dokumentationsreferenzmodell.md) /
  [`LH-FA-PROJDOCS-006`](../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)):
  die `Schaerft`-Aufwaerts-Deklaration ist die Aenderungskopplung Spec-ADR.
- **Aufloesungs-Trigger:** erledigt (CR ausgefuehrt 2026-07-24; Delivery-Hash im
  CR-Slice). Die status-basierte Grandfather-Grenze gilt permanent; neue
  Accepted-ADRs entstehen bereits MADR-konform, kein neuer Grandfather noetig.

### MR-009 - Skill-Dateien unter `.harness/skills/`, Review-Reports unter `docs/reviews/`

- **Datum:** 2026-07-25
- **Geltungsbereich:** `.harness/skills/reviewer.md`,
  `.harness/skills/closure-note-reviewer.md`,
  [`docs/reviews/`](../docs/reviews/README.md) (Ablage + README);
  Quellen-Rolle von [`review.md`](review.md) und
  [`verification.md`](verification.md).
- **Adaption:** Review-Wissen liegt in u-boot an **zwei** Orten mit klarer
  Rollentrennung: die kanonische **Prosa** bleibt im Autoren-Bestand `harness/`
  (`review.md`, `verification.md`), die vom Regelwerk geforderte **Skill-Form**
  (Kontext-Eingang, repo-spezifische Klassifikations-Anker, Nicht-Zustaendig-
  keiten, Output-Schema, Steering-Loop) liegt maschinen-materialisiert unter
  `.harness/skills/` (`MR-007`-Ortswahl). Die Skills **duplizieren keine
  Tabellen**, sie verweisen aufwaerts auf die Prosa; bei Abweichung gewinnt die
  Prosa. Die Report-Ablage ist `docs/reviews/` (ein Report pro Lauf, Folgelaeufe
  als neue Datei); die Report-**Vorlage** wird bewusst **nicht** ins Repo
  kopiert, sondern aus
  `.harness/baseline/<tag>/templates/docs/reviews/` genommen - eine zweite Kopie
  waere eine Drift-Quelle. Zwei Abweichungen vom Vorlagen-Wortlaut: (a) u-boot
  hat **kein** `agent-review`-Make-Target - Reviews sind agentisch, Beleg ist
  der Report; (b) der Closure-Note-Skill hat **kein** vorgelagertes
  computational Struktur-Gate (ein `check_closure_notes.py`-Aequivalent
  existiert hier nicht) und traegt deshalb Struktur- **und** Inhaltspruefung,
  mit den Pflichtfeldern aus [`verification.md`](verification.md) als
  Struktur-Quelle.
- **Begruendung:** Das Regelwerk fordert eine per-Repo Skill-Datei, nicht nur
  Prosa - ohne sie driftet der Reviewer zwischen Sessions. Die Doku-Mindest-
  struktur
  ([`LH-FA-PROJDOCS-001`](../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur))
  ist als **Mindest-**, nicht als Maximalstruktur formuliert; `docs/reviews/`
  ergaenzt sie additiv und erfuellt die dortige README-Pflicht je
  Unterverzeichnis. Kein Change Request am Vertrags-Stratum noetig.
- **Aufloesungs-Trigger:** permanent. Re-evaluieren, falls u-boot ein
  computational Closure-Note-Gate einfuehrt (dann wird der Closure-Note-Skill
  auf die semantische Schicht zurueckgeschnitten) oder falls ein
  `agent-review`-Target entsteht.

## Zusatzklassen-Deklaration fuer Sensors-Bindung

Ueber die kanonischen Bindung-Klassen (ADR, Carveout, Kalibrierung/Schwelle,
Reproduzierbarkeit) hinaus nutzt dieses Repo:

| Klasse | Form | Bedeutung | Beispiel |
|---|---|---|---|
| Anforderungs-Bindung | `LH-*` | Gate prueft eine bestimmte Lastenheft-Anforderung direkt | Exit-Code-Vertrag [`LH-FA-CLI-006`](../spec/lastenheft.md#lh-fa-cli-006--exit-codes) |
| Golden-/Replay-Bindung | Golden-Case-Satz | Replay-Gate haengt an einem fixierten Generator-Output | Fresh-State-/Idempotenz-Cases je CLI-Generator |

## Modus-Deklaration pro Sub-Area

Stand: auditiert am 2026-07-25 (Drei-Achsen-Inklusion + vier Modus-Kriterien je
Kandidat). Das **Audit-Protokoll** - inklusive der abgewiesenen Kandidaten
(Sub-Area-Aspirantinnen) - liegt im Slice
[`slice-harness-sub-area-modus-audit`](../docs/plan/planning/open/slice-harness-sub-area-modus-audit.md)
§9; hier steht nur das Ergebnis.

| Sub-Area (Pfad) | Modus | Begruendung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|
| `spec/`, `harness/`, `docs/plan/` (Spec, Architektur, ADR, Planung, Konventionen) | Greenfield | Doku-fuehrt-Sub-Areas; Spec/Architektur beschreiben vor dem Code. | n/a (GF) |
| `internal/hexagon/domain` | Greenfield | Reine Datentypen; I/O-Freiheit und Value-Object-Pflicht standen als Architektur-Regel vor dem Code. | n/a (GF) |
| `internal/hexagon/application` | Greenfield | Use-Case-Schnittstellen lagen als Ports in der Architektur-Sicht, bevor die Services entstanden. | n/a (GF). Eine code-seitig entstandene Konvention (nil-tolerante Ports via `noop*`-Defaults) wird in die Sicht-Spec nachgezogen. |
| `internal/hexagon/port` (`driving` + `driven`) | Greenfield | Kreuz-blinde Port-Trennung ist Architektur-Vorgabe und depguard-durchgesetzt. | n/a (GF) |
| `internal/adapter/driving/cli` | **Hybrid** | Exit-Code-/JSON-Vertrag sind spec-gefuehrt, aber zwei tragende Implementierungs-Konventionen (Sentinel-Schichtung, Dual-Classifier-Regel) sind aus Slice-/Review-Arbeit entstanden und leben nur im Code-Kommentar - dort laeuft die Richtung Code -> Doku. | Beide Regeln stehen in der Sicht-Spec (Abschnitt Fehlermodelle), geliefert von [`slice-harness-architecture-template-konformitaet`](../docs/plan/planning/open/slice-harness-architecture-template-konformitaet.md) -> dann GF. |
| `internal/adapter/driven` | Greenfield | Adapter implementieren vorher definierte Driven-Ports; der Port-Pin (`var _ driven.X`) macht Drift zum Build-Fehler. | n/a (GF) |
| `internal/e2e` (Test-Infrastruktur) | Greenfield | Build-Tag-Konvention (`//go:build docker`) und Fake-Clock-Pflicht sind in der Architektur-Sicht (Abschnitt Tests) verankert. | n/a (GF) |
| `tools/`, `scripts/` (Harness-Tooling) | Greenfield | Skripte materialisieren zuerst geschriebene Konventionen (`MR-004`, `MR-005`, Coverage-Bootstrap). | n/a (GF) |
| `internal/**/README.md` (Code-Paket-READMEs) | **Brownfield** | Beschreiben den Code-Bestand nachtraeglich, tragen Meilenstein-/Tranchen-Tags und ungelinkte `LH-*`-Kennungen; `internal/**` liegt im `scan.ignore` - die Inventur-Linie existiert, hat aber keinen Sensor. | Kennungen verlinkt, Status-Abschnitte entzeitlicht, `internal/**` aus dem `scan.ignore` entfernt: [`slice-harness-internal-readme-kennungs-retrofit`](../docs/plan/planning/open/slice-harness-internal-readme-kennungs-retrofit.md) -> dann GF. |
| `cmd/uboot` (Wiring/Entrypoint) | - (keine Sub-Area) | Erfuellt nur eine der drei Inklusions-Achsen (eigenes Verzeichnis, aber keine eigene Konvention und keine eigenstaendige Inventur-Linie): **Sub-Area-Aspirantin**, gefuehrt in der Hexagon-Schichtungs-Linie. | n/a - re-evaluieren, wenn das Wiring eigene Regeln traegt (z. B. DI-Container). |

> Zwei Nicht-GF-Aussagen und ihre Bedeutung: **Hybrid** heisst, in dieser
> Sub-Area laufen beide Richtungen nebeneinander (Vertrag fuehrt, Detail-
> Konvention folgt dem Code); **Brownfield** heisst, die Doku beschreibt den
> Bestand nachtraeglich. Beide tragen oben eine benannte Graduation-Bedingung -
> eine BF-/Hybrid-Markierung ohne Graduation-Plan waere eine permanente Ausnahme
> als temporaer getarnt.

## Glossar (optional)

Repo-spezifische Begriffe stehen im Lastenheft
([`spec/lastenheft.md`](../spec/lastenheft.md) Abschnitt Glossar/Begriffe); hier keine
Wiederholung.
