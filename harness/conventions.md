# Harness-Konventionen — u-boot

## Purpose

Diese Datei deklariert die *repo-lokalen* Strukturregeln von `u-boot`
gegenüber der adoptierten Harnesskonvention (Baseline). Sie ist der
Default-Ort für:

- **Adaptionen** ggü. der Baseline (mit Begründung und Auflösungs-Trigger).
- **ID-Schema-Deklaration** — welches Präfix-Schema dieses Repo nutzt.
  Der Baseline-Default wird als Teil der `MR-000`-Aussage festgehalten;
  ein abweichendes Präfix oder Schema ist ein eigener `MR`-Eintrag.
- **Zusatzklassen-Deklarationen** für repo-spezifische
  Bindung-Klassen in der Sensors-Tabelle, die über die vier kanonischen
  hinausgehen (ADR, Carveout, Schwelle, Reproduzierbarkeit).
- **Modus-Deklarationen** pro Sub-Area (Greenfield / Brownfield /
  Hybrid) inklusive Konvergenz-Auftrag bei BF.

Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt die
kanonische Quelle (Source Precedence). Diese Datei ist konformitäts-
bringend für *Form*-Fragen, nicht autoritativ über Inhalt.

## Baseline

- **Konvention:** AI-Harness-Kurs (`pt9912/ai-harness-course`)
- **Stand:** v6.13.0 (Regelwerk-Bundle)
- **Datum der Adoption:** 2026-07-24 (Erst-Adoption direkt auf `v3.5.1`; Review-Bumps auf `v3.5.2` am 2026-07-25 und auf `v6.13.0` am 2026-09-29, Protokoll in den Bump-Slices, Verfahren in [`MR-004`](conventions/MR-004-regelwerk-vendored.md))
- **Integritäts-Pin:** `.harness/baseline/v6.13.0/SHA256SUMS` über den vendorten
  Bestand (`regelwerk/` + `templates/`); offline prüfbar per
  `tools/harness/fetch-baseline-cache.sh --verify`.

## Adoptierte Konventions-Quellen

- **Extern (Lehrmaterial):** <https://github.com/pt9912/ai-harness-course/tree/v6.13.0>
- **Vendored Baseline (Regelwerk + Templates):** aus dem self-contained
  Release-Asset `lab-regelwerk.zip` (Tag `v6.13.0`) nach
  `.harness/baseline/v6.13.0/{regelwerk,templates}/` entpackt (netzlos,
  `SHA256SUMS`); Index `regelwerk/README.md`. Die Skelett-Vorlagen unter
  `templates/` sind Referenz-Form („Ziel-Form“ des Regelwerks) und
  Kopiervorlage für neue Artefakte (ADR, Slice, Welle, Carveout, Review-Report).
  Bump-Prozedur, Sync-Trigger und Freshness-Audit stehen in
  [`MR-004`](conventions/MR-004-regelwerk-vendored.md).
- **In-Repo (verkörperte Form):** die Gate-Baseline (`.d-check.yml`, `Makefile`,
  `Dockerfile`, `.golangci.yml`, `.github/workflows/`) und die
  autoren-gepflegte Harness-Prosa unter `harness/` (`README.md`, `roles.md`,
  `review.md`, `verification.md`, `replay.md`, diese Datei samt
  `conventions/`) sowie `AGENTS.md`.

## Adaptions-Block

Regeln dieser Sektion: Diese Datei trägt den **Index**, nicht die Einträge.
Jede Adaption ist eine eigene Datei unter `harness/conventions/`, kopiert aus
`harness/conventions/MR-NNN-titel.template.md` der vendored Baseline;
ist ihr Auflösungs-Trigger eingetreten, wandert sie per `git mv` nach
`conventions/done/`. Der Zustand ist die Verzeichnis-Position, kein
Status-Feld. Der Grund für den Schnitt: Was hier steht, liest **jeder**
Agentenlauf — aufgelöste Adaptionen gehören nicht in diesen Pfad
(Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/conventions.md als Konventionsspeicher).

### MR-000 — Baseline-Aussage

Bleibt hier: Sie ist keine Adaption, sondern die Adoptions-Erklärung, und
sie gilt für jeden Lauf.

- **Datum:** 2026-07-24
- **Geltungsbereich:** gesamtes Repo
- **Ersetzt-Baseline-Regel:** — *(keine; dieser Eintrag ist die
  Adoptions-Erklärung, keine Adaption)*
- **Adaption:** *keine inhaltlichen Adaptionen ggü. Baseline-Default
  für Verzeichniskonvention, Lifecycle-Regeln (`open` → `next` → `in-progress` →
  `done`), Carveout-Disziplin und ID-Schema:* Vertrags-Präfix `LH`
  (`LH-FA-*`, `LH-QA-*`, `LH-NFA-*` u. a. Familien des Lastenhefts), `SPEC-<NNN>` und
  `ARC-<NNN>` (fest, kodieren das Stratum), `ADR-<NNNN>`, `CO-<NNN>`,
  `slice-<phase>-<slug>`, `tranche-<nr>-<slug>`, `MR-<NNN>`. Ein Bereichssegment
  wird nicht geführt (ein schreibender Entwickler-Kontext, Zählraum repo-weit).
  Konkrete Abweichungen sind als eigene `MR-<NNN>` im Index unten dokumentiert.
- **Begründung:** Initial-Setzung. u-boot war vor der Adoption bereits
  harness-geformt; dieser Block hält den konformen Grundstand fest, spätere
  Adaptionen folgen als `MR-<NNN>`.
- **Auflösungs-Trigger:** permanent.

### Aktive Adaptionen

| MR | Titel | Geltungsbereich | Ersetzt-Baseline-Regel |
|---|---|---|---|
| [001](conventions/MR-001-source-precedence-drei-straten.md) <a id="mr-001"></a> | Source Precedence mit drei Spec-Straten (Vollform) | `AGENTS.md` Abschnitt Source Precedence, `harness/README.md` Abschnitt Source Precedence, `.d-check.yml` (Referenzmatrix). | [`grundlagen-source-precedence.md` §Source Precedence](../.harness/baseline/v6.13.0/regelwerk/grundlagen-source-precedence.md#source-precedence) |
| [002](conventions/MR-002-carveout-inventar.md) <a id="mr-002"></a> | Carveout-Inventar an fester Stelle statt `docs/plan/carveouts/` | Carveout-Ablage; `docs/plan/planning/in-progress/carveouts.md`, `AGENTS.md` Abschnitt Planning-Lifecycle, `.d-check.yml` `matrix`. | [`modul-07-carveouts.md` §Ziel-Form: Carveout](../.harness/baseline/v6.13.0/regelwerk/modul-07-carveouts.md#ziel-form-carveout) |
| [003](conventions/MR-003-roadmap-wellen.md) <a id="mr-003"></a> | Roadmap folgt Wellen-Template; Release-Versionen = Wellen; Ort in-progress/ | `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/README.md`. | [`modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Modul 6)](../.harness/baseline/v6.13.0/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6) |
| [004](conventions/MR-004-regelwerk-vendored.md) <a id="mr-004"></a> | Regelwerk-Lese-Form committet vendored; Baseline-Pin v6.13.0; beide Baeume | `.harness/baseline/`, `tools/harness/fetch-baseline-cache.sh`, `AGENTS.md` Abschnitt "Betriebsregelwerk (vendored Baseline)", `harness/READM | [`modul-02-harness-bootstrap.md` §Greenfield-Bootstrap: Schritt-Sequenz (Modul 2)](../.harness/baseline/v6.13.0/regelwerk/modul-02-harness-bootstrap.md#greenfield-bootstrap-schritt-sequenz-modul-2) |
| [005](conventions/MR-005-gate-haltung-docs-check.md) <a id="mr-005"></a> | Gate-Haltung: `docs-check` via direktem Container-Lauf; `scan.ignore` erweitert | `.d-check.yml`, `Makefile` (`docs-check`). | [`modul-02-harness-bootstrap.md` §Gate-Fragment `d-check.mk` (Schritt 2)](../.harness/baseline/v6.13.0/regelwerk/modul-02-harness-bootstrap.md#gate-fragment-d-checkmk-schritt-2) |
| [007](conventions/MR-007-ortswahl-harness-verzeichnis.md) <a id="mr-007"></a> | Ortswahl `.harness/` (dot-prefixed, committet) neben `harness/` | `.harness/` (vendored Baseline, kuenftig `.harness/skills/`), `harness/` (Autoren-Prosa), `.gitignore`. | [`grundlagen-harness-dateien.md` §Verzeichniskonvention](../.harness/baseline/v6.13.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention) |
| [009](conventions/MR-009-skills-und-review-ablage.md) <a id="mr-009"></a> | Skill-Dateien unter `.harness/skills/`, Review-Reports unter `docs/reviews/` | `.harness/skills/reviewer.md`, `.harness/skills/closure-note-reviewer.md`, `docs/reviews/` (Ablage + README); Quellen-Rolle von `review.md`  | [`modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill](../.harness/baseline/v6.13.0/regelwerk/modul-10-review-harness.md#ziel-form-reviewer-skill) |
| [010](conventions/MR-010-architektur-sicht-abschnittsfolge.md) <a id="mr-010"></a> | Architektur-Sicht: Vorlagen-Inhalt, erweiterte Abschnittsfolge | `spec/architecture.md` (Sicht-Stratum). | [`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../.harness/baseline/v6.13.0/regelwerk/modul-03-spec.md#ziel-form-architektur-sicht) |

### Aufgelöste Adaptionen

| MR | aufgelöst durch |
|---|---|
| [006](conventions/done/MR-006-modus-deklaration.md) <a id="mr-006"></a> | — (Audit abgeschlossen; das Ergebnis steht im Abschnitt Modus-Deklaration unten) |
| [008](conventions/done/MR-008-adr-form.md) <a id="mr-008"></a> | — (Change Request ausgeführt: MADR-Form im Lastenheft, Version 0.2.0) |

## Zusatzklassen-Deklaration für Sensors-Bindung

Über die kanonischen Bindung-Klassen (ADR, Carveout, Kalibrierung/Schwelle,
Reproduzierbarkeit) hinaus nutzt dieses Repo:

| Klasse | Form | Bedeutung | Beispiel |
|---|---|---|---|
| Anforderungs-Bindung | `LH-*` | Gate prueft eine bestimmte Lastenheft-Anforderung direkt | Exit-Code-Vertrag [`LH-FA-CLI-006`](../spec/lastenheft.md#lh-fa-cli-006--exit-codes) |
| Golden-/Replay-Bindung | Golden-Case-Satz | Replay-Gate haengt an einem fixierten Generator-Output | Fresh-State-/Idempotenz-Cases je CLI-Generator |

## Modus-Deklaration pro Sub-Area

Das Kürzel-Segment (`ADR-<KÜRZEL>-NNNN`) wird in u-boot nicht geführt; die Tabelle trägt keine Kürzel-Spalte.

Stand: auditiert am 2026-07-25 (Drei-Achsen-Inklusion + vier Modus-Kriterien je
Kandidat). Das **Audit-Protokoll** - inklusive der abgewiesenen Kandidaten
(Sub-Area-Aspirantinnen) - liegt im Slice
[`slice-harness-sub-area-modus-audit`](../docs/plan/planning/done/slice-harness-sub-area-modus-audit.md)
§9; hier steht nur das Ergebnis.

| Sub-Area (Pfad) | Modus | Begruendung | Graduation-Bedingung / Folge-Slice |
|---|---|---|---|
| `spec/`, `harness/`, `docs/plan/` (Spec, Architektur, ADR, Planung, Konventionen) | Greenfield | Doku-fuehrt-Sub-Areas; Spec/Architektur beschreiben vor dem Code. | n/a (GF) |
| `internal/hexagon/domain` | Greenfield | Reine Datentypen; I/O-Freiheit und Value-Object-Pflicht standen als Architektur-Regel vor dem Code. | n/a (GF) |
| `internal/hexagon/application` | Greenfield | Use-Case-Schnittstellen lagen als Ports in der Architektur-Sicht, bevor die Services entstanden. | n/a (GF). Eine code-seitig entstandene Konvention (nil-tolerante Ports via `noop*`-Defaults) wird in die Sicht-Spec nachgezogen. |
| `internal/hexagon/port` (`driving` + `driven`) | Greenfield | Kreuz-blinde Port-Trennung ist Architektur-Vorgabe und depguard-durchgesetzt. | n/a (GF) |
| `internal/adapter/driving/cli` | Greenfield (**graduiert 2026-07-25**, vorher Hybrid) | Exit-Code- und JSON-Vertrag waren immer spec-gefuehrt; die beiden Implementierungs-Konventionen (Sentinel-Schichtung, Dual-Classifier-Regel) lebten nur im Code-Kommentar - Richtung Code -> Doku. Diese Luecke ist geschlossen: beide stehen jetzt in der Sicht-Spec (Abschnitt Fehlermodelle), die Dual-Classifier-Regel dabei praeziser als der Code-Kommentar (zentrale Klassifikation **plus** je Subkommando eine Diagnostic-Abbildung). | n/a (GF, Graduation vollzogen). Rueckstufung, falls wieder eine tragende Adapter-Konvention nur im Code entsteht. |
| `internal/adapter/driven` | Greenfield | Adapter implementieren vorher definierte Driven-Ports; der Port-Pin (`var _ driven.X`) macht Drift zum Build-Fehler. | n/a (GF) |
| `internal/e2e` (Test-Infrastruktur) | Greenfield | Build-Tag-Konvention (`//go:build docker`) und Fake-Clock-Pflicht sind in der Architektur-Sicht (Abschnitt Tests) verankert. | n/a (GF) |
| `tools/`, `scripts/` (Harness-Tooling) | Greenfield | Skripte materialisieren zuerst geschriebene Konventionen (`MR-004`, `MR-005`, Coverage-Bootstrap). | n/a (GF) |
| `internal/**/README.md` (Code-Paket-READMEs) | Greenfield (**graduiert 2026-07-25**, vorher Brownfield) | Sie beschrieben den Code-Bestand nachtraeglich, trugen Meilenstein-/Tranchen-Tags und ungelinkte `LH-*`-Kennungen, und `internal/**` lag im `scan.ignore` - eine Inventur-Linie ohne Sensor. Alle drei Punkte sind aufgeloest: Kennungen verlinkt, Status-Abschnitte entzeitlicht, Glob entfernt ([`slice-harness-internal-readme-kennungs-retrofit`](../docs/plan/planning/done/slice-harness-internal-readme-kennungs-retrofit.md)). | n/a (GF, Graduation vollzogen). **Bekannte Sensor-Grenze:** `ids` wertet nur *bare* Kennungen; in Backticks gesetzte IDs bleiben ungeprueft - die Abdeckung ist real, aber nicht vollstaendig. |
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
