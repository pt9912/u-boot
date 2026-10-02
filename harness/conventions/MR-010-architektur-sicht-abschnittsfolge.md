# MR-010 — Architektur-Sicht: Vorlagen-Inhalt, erweiterte Abschnittsfolge

- **Datum:** 2026-10-02
- **Geltungsbereich:** [`spec/architecture.md`](../../spec/architecture.md) (Sicht-Stratum).
- **Ersetzt-Baseline-Regel:** [`modul-03-spec.md` §Ziel-Form: Architektur-Sicht](../../.harness/baseline/v6.13.0/regelwerk/modul-03-spec.md#ziel-form-architektur-sicht)
- **Adaption:** Die Sicht folgt der Vorlage `architecture.template.md` im Inhalt
  (Kopf mit Rolle und Hard Rule, `ARC-<NNN>` fuer Komponenten und externe
  Beruehrungspunkte, Schichten mit Constraints, externe Abhaengigkeiten,
  Sequenzen, Fehlermodelle), behaelt aber ihre **erweiterte Abschnittsfolge**
  (Ueberblick, Schichten je Komponente, Import-Regeln, Enforcement,
  Sequenz-Diagramme, Fehlermodelle, Tests, Anti-Patterns, Evolution). Die
  Import-Regel-Tabelle fuehrt die Spezifikation ([`SPEC-013`](../../spec/spezifikation.md#spec-013--import-regel-tabelle-der-schichten)); die Sicht
  visualisiert sie und verweist aufwaerts.
- **Begruendung:** Die Abschnittsnummern sind Anker fuer Verweise aus ADRs
  (`Schaerft:`); eine Umnummerierung nach Vorlage waere ein Anker-Bruch ohne
  Nutzen. Die Zusatzabschnitte sind reine Sicht (kein Anforderungsinhalt).
- **Aufloesungs-Trigger:** permanent, solange die Abschnittsnummern als Anker
  genutzt werden.
