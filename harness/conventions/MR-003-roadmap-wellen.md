# MR-003 — Roadmap folgt Wellen-Template; Release-Versionen = Wellen; Ort in-progress/

- **Datum:** 2026-07-24
- **Geltungsbereich:**
  [`docs/plan/planning/in-progress/roadmap.md`](../../docs/plan/planning/in-progress/roadmap.md),
  [`docs/plan/planning/README.md`](../../docs/plan/planning/README.md).
- **Ersetzt-Baseline-Regel:** [`modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Modul 6)](../../.harness/baseline/v6.18.1/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6)
- **Adaption (neu gefasst 2026-10-16):** Die Roadmap folgt der `roadmap.template.md`-Struktur
  der vendorten Baseline (Offene Wellen, Naechste Wellen, Meilensteine, Abhaengigkeitsgraph,
  Abgeschlossene Wellen, Historische Trigger-Verschiebungen). u-boots **Release-Versionen sind die
  Wellen**; Termine erscheinen nur als *Konsequenz* einer abgeschlossenen Welle (Release-Datum),
  nicht als Treiber. Eine Welle bekommt bei der Eroeffnung eine flache Datei
  `welle-<kebab-slug>.md` unter `docs/plan/planning/` (Vorlage `welle.template.md`), bei Closure
  wandert sie per `git mv` nach `done/`; geplante Wellen stehen nur unter *Naechste Wellen*. Drei
  Abweichungen bleiben: (a) die Roadmap liegt unter
  `docs/plan/planning/in-progress/roadmap.md` (das Regelwerk nennt beide Pfade; u-boot loest
  zugunsten des Lifecycle-Verzeichnisses auf); (b) es gibt **keine** `welle-<Kennung>-results.md` -
  die Closure-Notiz steht in §7 der Welle-Datei und verweist auf den `done/`-Release-Cut-Slice als
  Detailquelle; (c) die Zeitdokumente geschlossener Wellen werden **nicht** archiviert
  (`make archive-welle` bleibt ungenutzt, `done/` behaelt den Volltext). Rueckwirkende
  Welle-Dateien fuer die Releases `v0.1.0` bis `v0.9.0` gibt es nicht; ihre Closure lebt in den
  `done/`-Release-Cut-Slices und in der Roadmap.
- **Begruendung:** Die Template-Struktur macht Wellen-Reihenfolge, Trigger und Abhaengigkeiten
  explizit; die Welle-Datei traegt Ziel, Trigger, Closure-Kriterien und Out-of-Scope vor der Arbeit.
  Der `done/`-Release-Cut-Slice traegt die Closure bereits vollstaendig, eine zweite Results-Datei
  dupliziert sie. Die Archivierung ersetzte `done/`-Inhalte durch Stubs und liefe gegen die
  Sichtbarkeit des Volltexts im Arbeitsbaum.
- **Aufloesungs-Trigger:** permanent, solange Release-Versionen die Wellen sind.
