# Verification Evidence: slice-cr-lastenheft-produktanforderungen — 2026-10-16

**Rolle:** Verifier (frischer Kontext, kein Self-Review). Frage: Wurde der Slice richtig gebaut, gemessen an Plan, DoD und Sensors? Form nach `harness/verification.md`. <!-- d-check:ignore (Pfad als Zitat) -->

**Gegenstand:** Diff `938df12..HEAD` (letzter Commit `7d032b3`), 28 Dateien. Alle Sensors wurden in diesem Lauf selbst gefahren (Docker-only), nicht aus der Implementer-Behauptung übernommen.

## Scope

- Slice: `slice-cr-lastenheft-produktanforderungen` <!-- d-check:ignore (Slice-Kennung als Zitat) -->
- IDs: Grabsteine `LH-FA-BUILD-001`..`009`, `LH-FA-PROJDOCS-001`..`006`, `LH-FA-ARCH-001`..`003`, `LH-QA-001`..`004`; neu `LH-NFA-PORT-004`, `LH-LESE-004`; Verfeinerungen `SPEC-025`..`SPEC-035`. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
- Artefakte: `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, `harness/`, Skills, `docs/`-Verweise.

## DoD-Abgleich

- [x] Freigabe der Einzelentscheidungen — Evidence: Vermerk im Slice (2026-10-16); Gegenbeleg liegt im Zielkontext des Aufrufers, nicht im Repo prüfbar (Aussage des Projektinhabers).
- [x] Lastenheft 0.6.0 — bestätigt:
  - 22 Grabsteine: genau 22 Treffer `Entfallen (0.6.0)`; je ID genau eine Überschrift und der Satz direkt darunter (Schleife über alle 22: head=1, tomb=1).
  - Überschriften unverändert: Überschriften-Diff `938df12` gegen HEAD zeigt ausschließlich zwei Zugänge (`LH-LESE-004`, `LH-NFA-PORT-004`), keine Änderung oder Entfernung. Damit sind auch die Anker unverändert. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - Grabstein-Texte enthalten keinen Verweis auf ADR, Slice, Welle oder `docs/plan` (grep leer).
  - §12 bereinigt (Diff: sechs Prozess-Punkte ersetzt durch den Verweis auf das Image), Kopf-Version `0.6.0`, Zeile 0.6.0 in §16 vorhanden, Einleitungen §4.11, §4.12, §4.13 angepasst.
- [x] Spezifikation — bestätigt:
  - 13 Verfeinerungen: 11 neu als `SPEC-025`..`SPEC-035` (Anhang-Tabelle stimmt mit den Überschriften überein), 2 als `LH-NFA-PORT-004.a`/`.b`. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - SPEC-Nummern: Liste 001..010, 012..035; keine Doppelvergabe (`uniq -c` ohne Treffer >1). Die Lücke bei `SPEC-011` besteht schon in `938df12` (Alt-Datei enthält sie nicht) und gehört nicht zum CR. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - Alte Anker: alle 13 als `<a id=...>` direkt unter der neuen Überschrift; `make docs-check` löst die Alt-Links aus `done/` und `docs/reviews` auf (0 Befunde).
- [x] Lebende Verweise — **mit Einschränkung**, siehe Abweichung A-1: die DoD-Zeile ist für die dort genannten Stellen erfüllt; eine nicht genannte Klasse (Paket-READMEs unter `internal/`) zitiert die Grabsteine weiter als geltende Zusage.
- [x] `trace` und `doc-complete`: 83 Anforderungen, 0 Waisen — Evidence: Sensor-Lauf unten. Aussagekraft für die 22 Grabsteine begrenzt (Review F-7, INFO).
- [x] `make gates` grün — Evidence: Sensor-Lauf unten (Exit 0, Coverage 90,90 % bei Schwelle 90 %).
- [x] Review durchgeführt — Report `docs/reviews/2026-10-16-slice-cr-lastenheft-produktanforderungen.md` liegt vor; Zahlen stimmen (0 HIGH, 0 MEDIUM, 6 LOW, 2 INFO). <!-- d-check:ignore (Kennungen im Zeitdokument) -->
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag — offen, bewusst (Slice §7 „Offen.“); vor `done/` zu schreiben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang — offen: die drei Risiken stehen noch auf „offen“ (siehe A-1, A-3).

## Prüfung der Review-Befunde

| Befund | Soll | Ist-Beleg |
| --- | --- | --- |
| F-1 SPEC-026 Doppelverweis | behoben | SPEC-026 nennt `SPEC-008` unverändert, aber mit je eigenem Bezug (MVP-Pflicht „via diesem SPEC-026“); kein doppelter „SPEC-008 (V1)“-Wortlaut mehr. Bestätigt. | <!-- d-check:ignore (Kennungen im Zeitdokument) -->
| F-2 architecture.md | behoben | `Bezug` nennt genau `SPEC-012`, `SPEC-013`, `SPEC-034`; `Letzte Änderung: 2026-10-16`. Bestätigt. | <!-- d-check:ignore (Kennungen im Zeitdokument) -->
| F-3 Decken-Regel in §16 | behoben | §16-Einleitung belegt die Regel mit „Decken-Regel des adoptierten Baseline-Regelwerks“ statt mit dem Grabstein. Bestätigt. |
| F-5 Archiv/Roadmap | behoben | `docs/archive/README.md` ohne Link auf den Grabstein; `roadmap.md` zitiert `modul-07-carveouts.md` und `modul-05-planning-harness.md` §Lifecycle statt Grabsteinen. Bestätigt (Inhaltsgleichheit der Archiv-Konvention mit dem Baseline-Regelwerk nicht gesondert geprüft, Reviewer-Urteil übernommen). |
| F-4, F-6, F-7, F-8 | bewusst offen / INFO | unverändert offen, so im Slice und Report benannt. F-6 siehe A-1. |

## Sensors

| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make gates` | pass | Exit 0; `coverage-gate: OK — coverage 90.90% meets threshold 90%`; alle Pakete `ok` |
| `make docs-check` | pass | `d-check: 235 Datei(en) geprüft, 0 Befund(e)` (inkl. Alt-Anker-Auflösung aus `done/` und Reviews) |
| `make doc-complete` | pass | `83 Anforderung(en), 0 Waise(n)` |
| `make doc-planning` | pass | 0 Befunde |
| `make doc-targets` | pass | 0 Befunde |
| `make baseline-verify` | pass | `v6.18.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `git diff 938df12..HEAD -- docs/plan/adr/` | pass | nur `docs/plan/adr/README.md` (1 Zeile, Index); keine Accepted ADR verändert |

## Traceability

| ID | Beleg |
| --- | --- |
| `LH-LESE-003` (Dokumentenordnung) | Grabsteine und Decken-Regel im Diff; `make docs-check` | <!-- d-check:ignore (Kennungen im Zeitdokument) -->
| `LH-NFA-PORT-004`, `LH-LESE-004` | Überschriften-Diff gegen `938df12`; Matrix `doc-complete` führt `LH-LESE-004` über den Slice; `LH-NFA-PORT-004` laut Review F-7 nicht in der Matrix (Familien-Muster) — hier nicht nachgeprüft <!-- d-check:ignore (Kennungen im Zeitdokument) -->
| `SPEC-025`..`SPEC-035` | Überschriften und Anker in `spec/spezifikation.md`; Anhang-Tabelle im Slice | <!-- d-check:ignore (Kennungen im Zeitdokument) -->
| Accepted ADRs | `git diff` auf `docs/plan/adr/`: nur Index |

## Abweichungen

- **A-1 (LOW, bestätigt F-6 und erweitert es):** Lebende Stellen zitieren die entfallenen Anforderungen weiter als Geltungsgrund.
  - Slice-DoD und Plan nehmen Code- und Konfig-Kommentare ausdrücklich aus; das ist gedeckt: `Dockerfile`, `repo.mk`, `.golangci.yml`, `.dockerignore`, `scripts/`, `.github/workflows/`, `cmd/` und `internal/**/*.go`.
  - **Nicht gedeckt ist** die Klasse Markdown unter `internal/`: `internal/README.md` (Zeile 4: Architektur-Pattern „nach“ `LH-FA-ARCH-001..003`; Zeile 88/89: Coverage-Scope „(`LH-FA-BUILD-008`, `-009`)“; Zeile 94: „Verbindliche Schicht-Regeln in `LH-FA-ARCH-003`“), `internal/hexagon/domain/README.md` (Zeile 68), `internal/adapter/driven/README.md`, `internal/adapter/driving/README.md`, `internal/hexagon/application/README.md`, `internal/hexagon/domain/README.md`, `internal/hexagon/port/driven/README.md`, `internal/hexagon/port/driving/README.md`. Sie sind keine Planungs-Historie und stehen nicht in der Ausnahmeliste der DoD („READMEs“ dort meint die Wurzel-READMEs laut Review); der Link zeigt auf einen Grabstein. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - `spec/architecture.md` §4 trägt in den `depguard`-`desc`-Zeilen weiter `LH-FA-ARCH-003` (11 Stellen); das Dokument ist Sicht und gibt die `.golangci.yml` wieder, die gleichlautende Zeichenketten trägt. Eine Änderung dort erfordert Gleichlauf mit der Konfiguration (Kopplung), daher offen entschieden. <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - `CHANGELOG.md`: historische Einträge; als Chronik tragbar.
  - Der Risiko-Ausgang „Suche nach den 22 IDs vor der Closure“ ist damit nicht abgeschlossen: Ergebnis der Suche liegt hier vor, die Entscheidung über die `internal/`-READMEs und die `desc`-Zeilen steht beim Planner.
- **A-2 (INFO):** `make baseline-verify` meldet Baseline `v6.18.1`, `AGENTS.md` nennt `v6.13.0` als vendored Stand. Nicht Gegenstand dieses CR (`AGENTS.md` im Diff nur Zeiger-Tabelle); als Beobachtung für den Planner, vor Folgearbeit zu klären.
- **A-3 (INFO):** Grabstein-Texte nennen „technisches Dokument“ und „Architektur-Sicht“ ohne Link (Review F-4); unverändert bewusst offen.

## Plan-vs-Code-Diff (Planänderungen über den CR hinaus)

- **Review F-8 bestätigt:** `docs/plan/planning/in-progress/slice-harness-wellen-betrieb.md` ist im Diff mit 95 Zeilen Änderung. Es ist eine Planänderung eines **anderen** Slice, nicht Teil des CR-Umfangs „Grenze Produkt/Prozess“: <!-- d-check:ignore (Kennungen im Zeitdokument) -->
  - Titel und Ziel umgeschrieben (Archivierung gestrichen, Entscheidung „Welle-Dateien und Register ja, Archivierung nein“ vom 2026-10-16),
  - zwei DoD-Punkte ersetzt (Entscheidungs-Haken und „erste Archivierung“ entfallen; neue Punkte „Roadmap und Doku auf Welle-Datei-Form“, „Werkzeuge auf u-boot gefasst“),
  - Voraussetzung auf den CR-Slice umgestellt, `Verantwortlich` gesetzt, Nicht-Ziele erweitert.
  - Der CR-Slice nennt diese Änderung nur als Auflösung der Blockade („Löst die Blockade von …“), nicht als Teil des Umfangs. Entscheidung ist im Slice selbst als Notiz belegt; ein Plan-Anker außerhalb (Roadmap-Drift-Log, ADR) wurde nicht gefunden. Empfehlung: Planner bestätigt die Umplanung ausdrücklich oder legt sie ins Drift-Log der Roadmap.
- Sonst keine Planabweichung: Plan-Tabelle (§3) gegen Diff — `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, `harness/`, `AGENTS.md`, `docs/`, `README*` berührt; `.d-check.yml` geprüft und nur im Kommentar geändert (keine Schwelle, kein Muster).
- Gates nicht gelockert: Coverage-Schwelle 90 % unverändert (Gate-Ausgabe), keine Konfigurationsänderung an `.golangci.yml`.

## Carveouts

- Neu: none. Gelöst: none. Unverändert: Inventar `docs/plan/planning/in-progress/carveouts.md` (nur Verweis-Nachzug im Diff).

## Nicht ausgeführt

- `make test-docker` — kein Docker-Laufzeitpfad berührt (Diff enthält keinen Code, nur Markdown/Kommentar-freie Doku und `.d-check.yml`-Kommentar).
- `make verify-depguard` — Importregeln nicht berührt; die `desc`-Zeilen in `spec/architecture.md` wurden nicht verändert (Wortlaut gleich der Konfiguration).
- Prüfung der Freigabe durch den Projektinhaber — nicht aus dem Repo beweisbar.

## Gesamturteil

DoD-Punkte 1, 2, 3 (mit A-1), 5, 6, 7 bestätigt; Punkte 8 und 9 bewusst offen. Keine Blocker für die Closure aus den Sensors. Vor `done/`: A-1 entscheiden (internal-READMEs nachziehen oder als Ausnahme in den Slice schreiben), F-8 bestätigen lassen, Closure-Notiz und Risiko-Ausgänge schreiben.

Commit / Artefakt: Stand `7d032b3`, dieser Bericht.
