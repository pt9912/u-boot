# MR-003 — Roadmap folgt Wellen-Template; Release-Versionen = Wellen; Ort in-progress/

- **Datum:** 2026-07-24
- **Geltungsbereich:**
  [`docs/plan/planning/in-progress/roadmap.md`](../../docs/plan/planning/in-progress/roadmap.md),
  [`docs/plan/planning/README.md`](../../docs/plan/planning/README.md).
- **Ersetzt-Baseline-Regel:** [`modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Modul 6)](../../.harness/baseline/v6.13.0/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6)
- **Adaption:** Die Roadmap folgt der `roadmap.template.md`-Struktur der
  vendorten Baseline
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
