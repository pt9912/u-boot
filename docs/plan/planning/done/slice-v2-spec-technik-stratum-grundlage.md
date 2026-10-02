# Slice V2: Technik-Stratum einführen (Fundament der Spec-Stratifizierung)

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `fdf4c92`**) — Technik-Stratum eingeführt; die Schnitt-Slices der Welle können starten.

**Welle:** `welle-spec-technik-stratum` (siehe [`roadmap.md`](../in-progress/roadmap.md)).
**Bezug:** [`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur)..[`-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) (Doku-Struktur und Referenzmodell, Scope).
**Berührte Spec-Stellen:** `harness/conventions.md` `MR-001`, `AGENTS.md` §Source Precedence, `harness/README.md`, `.d-check.yml` (Referenzmatrix), neues Dokument `spec/spezifikation.md`.
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Das Lastenheft ([`spec/lastenheft.md`](../../../../spec/lastenheft.md), 2984 Zeilen) vermischt Vertrag und Technik: JSON-Schema,
YAML-Schema von `u-boot.yaml`, Defaults, Dockerfile-/CI-Details, Verzeichnislayout. Die Baseline (v6.13.0,
`regelwerk/modul-03-spec.md`) sieht drei obligatorische Straten vor: Lastenheft (Vertrag) › Spezifikation (Technik) ›
Architektur (Sicht). `harness/conventions.md` `MR-001` führt bisher bewusst nur zwei. Der Projektinhaber hat die
Vollform beschlossen (2026-10-02): Technik-Stratum einführen und das Lastenheft auf den Vertrag zurückschneiden.

## Ziel und Abgrenzung

**Ziel:** `spec/spezifikation.md` existiert nach Vorlage (`.harness/baseline/v6.13.0/templates/spec/spezifikation.template.md`), ist in Source Precedence und Referenzmatrix verankert und wird von `make docs-check` als eigenes Stratum geprüft; `MR-001` ist abgelöst.

**Ausdrücklich NICHT in diesem Slice:**

- **Inhalt in die Spezifikation verschieben** — das tun die vier Schnitt-Slices; hier entsteht nur das leere, geprüfte Gefäß.
- **Lastenheft-Version anheben** — der Vertrag ändert sich erst mit dem ersten Schnitt; der Bump (0.4.0) und die Historie-Zeile gehören
  in den Abschluss-Slice, damit es genau einen Vertragsstand gibt.
- **Architektur-Sicht umbauen** — eigener Slice (`slice-v2-spec-architektur-angleichung`); die Sicht ist von der Technik unabhängig prüfbar.

**Schnittregeln** (gelten für alle Schnitt-Slices der Welle; die Schnitt-Slices verweisen auf diesen Abschnitt):

| Bleibt im Lastenheft (Vertrag, „Was") | Zieht in die Spezifikation (Technik, „Wie genau") |
|---|---|
| Anforderungstext mit Modalverb und Priorität | Algorithmen und Datenflüsse als Verfeinerung `LH-<…>-<NNN>.<a>` der jeweiligen Anforderung |
| Akzeptanzkriterien (Happy · Boundary · Negative) und Out-of-Scope | Datenstrukturen und Schemata (YAML-, JSON-, CLI-Tabellen) mit `SPEC-<NNN>` |
| Randbedingungen (`-RB-`) und Qualitätsanforderungen mit Messmethode | Defaults, Konstanten, Grenzwerte mit `SPEC-<NNN>` |
| Produktverträge: Exit-Code-Klassen, Sprachvertrag, Sicherheitszusagen | Fehler-, Diagnose- und Logging-Codes, Metrik- und Tracing-Felder |
| Alle `LH-*`-Kennungen und Überschriften **unverändert** (Anker bleiben) | Externe Verträge (Docker, Compose, Devcontainer, Versionsannahmen) |

- **Übernahme wörtlich, keine inhaltliche Änderung.** Der Schnitt verschiebt, er entscheidet nichts neu. Zweifelsfälle
  bleiben im Lastenheft und stehen im Slice unter „Offen".
- **Decken-Regel:** Das Lastenheft verweist nicht auf die Spezifikation (kein Rückzeiger); die Spezifikation verweist
  nie auf ADRs oder Slices.
- **Verbleib-Tabelle:** Jede verschobene Stelle steht im Slice mit Quelle (Lastenheft-§) und Ziel (Spezifikations-ID);
  nichts geht verloren, ein Gegenlesen zeigt, dass jede Quelle genau ein Ziel hat.
- **Zeilenverweise** auf `lastenheft.md` (Test- und Code-Kommentare der Form `spec/lastenheft.md:692-721`) werden im selben
  Slice für die berührten Stellen auf Kennung oder Abschnitt umgestellt.
## Definition of Done

- [x] `MR-001` abgelöst (Konventionen: drei Straten; Begründung und Auflösungs-Trigger), `AGENTS.md` und `harness/README.md` führen die neue Source Precedence.
- [x] `spec/spezifikation.md` angelegt (Kopf, sieben Abschnitte, Historie; Platzhalter-frei, noch ohne übernommene Inhalte) und in `.d-check.yml` als Stratum **Technik** mit Referenzregeln und Kennungs-Mustern (`SPEC-<NNN>`, Verfeinerung `LH-<…>-<NNN>.<a>`) geführt.
- [x] `make gates` grün (inklusive `make docs-check`), Exit-Code separat geprüft.
- [x] Review durch eine andere Rolle als die Umsetzung (`harness/review.md`), Report unter `docs/reviews/`.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | **Konvention.** `MR-001` ersetzen (neuer Text: drei Straten, Begründung „Vollform nach Baseline v6.13.0", Auflösungs-Trigger permanent); ID-Schemata (`SPEC-<NNN>`, Verfeinerung `.a`) in `harness/conventions.md` aufnehmen; `AGENTS.md` §Source Precedence und `harness/README.md` nachziehen. |
| T2 | **Gefäß.** `spec/spezifikation.md` aus der Vorlage anlegen (Kopf mit Bezug zum Lastenheft, §1–§7, Historie ohne ADR-/Slice-Verweise). |
| T3 | **Sensor.** `.d-check.yml`: Klasse `tech-spec` mit denselben Regeln wie `view-spec` (Decken-Regel), `ids.patterns` für `SPEC-<NNN>` und Verfeinerungen, Abgleich mit `.harness/baseline/v6.13.0/templates/.d-check.yml`. `make docs-check` und `make gates` grün. |

## Risiken

- Verfeinerungs-Kennungen mit Buchstabensuffix (Anforderungs-Kennung, Punkt, Buchstabe): Das bestehende Kennungs-Muster für `LH-…` matcht nur den Stamm; ohne Anpassung entstehen Falschbefunde oder Lücken. T3 prüft das mit einem Gegenbeispiel.
- Eine leere Spezifikation darf `docs-check` nicht rot färben (Platzhalter, Pflichtabschnitte): T2/T3 pinnen das mit dem ersten Lauf.

## Closure

- **Geliefert (`fdf4c92`):** `MR-001` in `harness/conventions.md` auf Fassung 2 (drei Straten, ID-Schemata `SPEC-<NNN>` und Verfeinerung mit Buchstabensuffix); `AGENTS.md` und `harness/README.md` führen die 10-Rang-Source-Precedence; `spec/spezifikation.md` als Gefäß nach Vorlage (Kopf, §1–§7, Historie); `.d-check.yml` mit Klasse `tech-spec`, Decken-Regel-Regeln und Kennungs-Muster `SPEC-<NNN>`.
- **Befund zu Verfeinerungs-Kennungen (Risiko ausgeräumt):** Eine nackte oder in Backticks gesetzte Verfeinerung (`<Anforderung>.a`) meldet d-check als ungelinkten Stamm der Anforderungs-Kennung; ein Link (auch auf `spezifikation.md`) genügt, Überschriften sind ausgenommen. Es braucht deshalb **kein** eigenes Muster; die Konvention lautet: Verfeinerungen werden im Fließtext verlinkt.
- **Leere Spezifikation:** hält `make docs-check` grün (166 Dateien, 0 Befunde).
- **Sensoren:** `make gates` grün (lint, test, coverage-gate, docs-check). Nicht ausgeführt: `make ci`, `make test-docker` (kein Code berührt).
- **Review:** unabhängiges Review über die gesamte Welle am Abschluss (`slice-v2-spec-schnitt-qualitaet-akzeptanz-abschluss`); dieser Slice ist reine Konfiguration und Gefäß.
- **Lerneintrag:** Die Matrix ließ sich mit vier Regelzeilen erweitern, weil `contract-spec`/`view-spec` bereits als Klassen existierten; ein generischer `spec`-Block mit `order` (Vorlage) wäre ein Umbau ohne Zusatznutzen gewesen.
