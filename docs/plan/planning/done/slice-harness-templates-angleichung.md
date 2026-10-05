# Slice: Harness-Dateien an die Baseline-Vorlagen angleichen

> **Status:** **Done** (2026-10-02, **Delivery-Hash: `4c70f04`**).

**Welle:** ohne Welle.
**Bezug:** Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt und §harness/conventions.md als Konventionsspeicher, `modul-09-implementierung.md` §Ziel-Form: AGENTS.md; Vorlagen unter `.harness/baseline/v6.13.0/templates/`.
**Berührte Spec-Stellen:** —
**Autor:** pt9912. **Datum:** 2026-10-02.

## Auslöser

Der Projektinhaber verlangt, `harness/README.md`, `harness/conventions.md` und `AGENTS.md` entsprechend den Vorlagen anzupassen.

## Ziel und Abgrenzung

**Ziel:** Die drei Dateien folgen der Gliederung der Vorlagen: `AGENTS.md` mit den Abschnitten 1–6 (Hard Rules als 3.x, Gate-Index nur als Zeiger, Dokumentations-Regeln als Tabelle); `harness/README.md` mit Purpose, Source precedence, Guides, Sensors (Target · Vertrag · Bindung), Traceability rules, Safety and scope boundaries, Minimal agent workflow, Leseordnung; `harness/conventions.md` als Index mit je einer Datei pro Adaption unter `harness/conventions/` (aktiv und `done/`) und dem Pflichtfeld „Ersetzt-Baseline-Regel“ je Eintrag.

**Ausdrücklich NICHT in diesem Slice:**

- **Inhaltliche Änderung der Adaptionen oder Hard Rules** — Text wird übernommen und nur neu gegliedert.
- **Append-only-Sensor für `MR`-Dateien (`vcs`-Modul)** — die Konfiguration kennt einen Immutabilitäts-Block für ADRs; ein zweiter für `MR`-Dateien ist optional und wird nicht aktiviert.

## Definition of Done

- [x] Drei Dateien nach Vorlage gegliedert, zehn Adaptionen als Einzeldateien, Gate-Index zentral in `harness/README.md` (`.d-check.yml` `targets`), Verweise nachgezogen.
- [x] `make gates` grün.

## Tranchen

| T | Inhalt |
| - | ------ |
| T1 | `harness/conventions.md` zum Index, Adaptionen als Einzeldateien mit Pflichtfeld „Ersetzt-Baseline-Regel“, Freshness-Audit und Sync-Trigger in `MR-004`. |
| T2 | `harness/README.md` nach Vorlage, Gate-Index mit Bindung. |
| T3 | `AGENTS.md` nach Vorlage; `.d-check.yml` (`targets.authority`), Skript-Kommentar und Verweise. |

## Risiken

- Der Gate-Index wandert von `AGENTS.md` nach `harness/README.md`: Das `targets`-Modul prüft gegen die neue Autorität; `make docs-check` ist grün.

## Closure

- **Geliefert (`4c70f04`):** `AGENTS.md` mit den Abschnitten 1–6 der Vorlage (Hard Rules 3.1–3.7 nach Vorlage, die repo-spezifischen Rules stehen als Zeiger-Tabelle 3.8 (kanonische Quelle je Thema; die Rules selbst werden in `AGENTS.md` nicht wiederholt), Gate-Index nur als Zeiger, Dokumentations-Regeln als Tabelle); `harness/README.md` mit allen Vorlagen-Abschnitten und dem einzigen Gate-Index (Target · Vertrag · Bindung); `harness/conventions.md` als Index (Purpose, Baseline, Konventions-Quellen, `MR-000`, Tabellen „Aktive“ und „Aufgelöste“ Adaptionen, Zusatzklassen, Modus-Deklaration, Glossar) mit zehn Adaptionen als Einzeldateien unter `harness/conventions/`; alle zehn Adaptionen sind aktiv (`MR-006` und `MR-008` tragen geltende Reste und bleiben im Index; „Aufgelöste Adaptionen“ ist leer); Freshness-Audit und Sync-Trigger stehen jetzt in `MR-004`.
- **Anpassungen:** `.d-check.yml` `targets` (`doc-tables` und `authority` auf `harness/README.md`), Skript-Kommentar in `tools/harness/fetch-baseline-cache.sh` (`--verify` weiter ok), Verweise auf Abschnitte von `AGENTS.md` in den Adaptions-Dateien und der JSON-Vertragsdoku. Ersetzte-Baseline-Regel je Adaption verlinkt mit Anker in das vendorte Regelwerk (Linkprüfung grün).
- **Sensoren:** `make gates` grün (lint, test, coverage-gate, docs-check); `tools/harness/fetch-baseline-cache.sh --verify` ok. Nicht ausgeführt: `make ci`, `make test-docker`.
- **Review (unabhängige Rolle):** `docs/reviews/2026-10-02-harness-templates-angleichung.md` — 0 HIGH, 4 MEDIUM, 6 LOW, 3 INFO; behoben: Bindung von `doc-immutable`, Pflichtfelder und Selbstverweise in `MR-004`, Rangzählung (AGENTS.md jetzt zehn Ränge wie die README) samt `MR-001`-Erweiterung, `MR-006`/`MR-008` zurück in den aktiven Index, abgeschnittene Index-Zellen, Verweise auf `conventions.md` für einzelne Adaptionen, Bump-Kette, Traceability-Regel; aus der Rückfrage des Projektinhabers entfällt die Hard Rule „Hexagonale Architektur“ in `AGENTS.md` (steht in `spec/architecture.md`). Bewusst offen: Sensors-Tabelle führt nur `doc-immutable` aus den `d-check.mk`-Targets (INFO).
