# MR-003 — Roadmap folgt Wellen-Template; Release-Cuts wellenlos; Ort in-progress/

- **Datum:** 2026-07-24
- **Geltungsbereich:**
  [`docs/plan/planning/in-progress/roadmap.md`](../../docs/plan/planning/in-progress/roadmap.md),
  [`docs/plan/planning/README.md`](../../docs/plan/planning/README.md),
  `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md`,
  `.claude/agents/planner.md` (Kurzform).
- **Ersetzt-Baseline-Regel:** [`modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Modul 6)](../../.harness/baseline/v6.18.1/regelwerk/modul-06-roadmap.md#roadmap-struktur-fünf-abschnitte-modul-6)
- **Adaption (neu gefasst 2026-11-02):** Die Roadmap folgt der `roadmap.template.md`-Struktur
  der vendorten Baseline (Offene Wellen, Naechste Wellen, Meilensteine, Abhaengigkeitsgraph,
  Abgeschlossene Wellen, Historische Trigger-Verschiebungen). u-boots **Release-Cuts sind
  keine Wellen** (Beschluss des Projektinhabers 2026-11-02): Ein Release-Cut ist ein
  welle-loses Release-Cut-Slice, das in der Roadmap unter *Offene
  Wellen* geführt wird, nach dem Tag unter *Abgeschlossene Wellen* mit Tag-Hash und Datum
  steht und in §7 seiner `done/`-Fassung schließt (Tag-Hash, Lauf-IDs, Lerneintrag).
  **Eine Welle-Datei bekommt der Cut nicht**; die Closure eines Release lebt vollständig
  im `done/`-Cut-Slice und der Roadmap-Zeile — eine Welle-Datei daneben würde dieselbe
  Beziehung zweimal führen. Grundlage ist die dokumentierte Praxis: die vier welle-losen
  Cuts `v0.13.1`, `v0.15.1`, `v0.16.1` und `v0.17.0` sowie die Releases `v0.1.0` bis
  `v0.9.0` ohne Welle-Datei; bis zur Neufassung trug die Praxis keinen MR-Anker
  (Drift-Audit-Kandidat, Frage im Cut-Slice [`slice-v1-release-cut-v0.18.0`](../../docs/plan/planning/open/slice-v1-release-cut-v0.18.0.md) §6).
  Wellen-Dateien (`welle-<kebab-slug>.md`, Vorlage `welle.template.md`) bleiben der Arbeit
  vorbehalten, die nicht als Release-Cut läuft (etwa eine mehrschichtige Welle mit eigenem
  Validation-Beleg); eine solche Welle bekommt bei der Eroeffnung die flache Datei unter
  `docs/plan/planning/`, bei Closure wandert sie per `git mv` nach `done/`; geplante Wellen
  stehen nur unter *Naechste Wellen*. Termine erscheinen nur als *Konsequenz* abgeschlossener
  Arbeit (Release-Datum nach dem Tag), nicht als Treiber. Drei Abweichungen bleiben:
  (a) die Roadmap liegt unter
  `docs/plan/planning/in-progress/roadmap.md` (das Regelwerk nennt beide Pfade; u-boot loest
  zugunsten des Lifecycle-Verzeichnisses auf); (b) es gibt **keine** `welle-<Kennung>-results.md` -
  die Closure-Notiz steht in §7 der Welle-Datei bzw. — für Release-Cuts — in §7 des
  `done/`-Cut-Slices; (c) die Zeitdokumente geschlossener Wellen werden **nicht** archiviert
  (`make archive-welle` bleibt ungenutzt, `done/` behaelt den Volltext). Rueckwirkende
  Welle-Dateien fuer die Releases `v0.1.0` bis `v0.9.0` gibt es nicht; ihre Closure lebt in den
  `done/`-Release-Cut-Slices und in der Roadmap.
- **Begruendung:** Die Template-Struktur macht Wellen-Reihenfolge, Trigger und
  Abhaengigkeiten explizit; die Welle-Datei traegt Ziel, Trigger, Closure-Kriterien und
  Out-of-Scope vor der Arbeit. Für Release-Cuts lebt die Closure bereits vollstaendig im
  Cut-Slice, eine zweite Ablage (Welle-Datei oder Results-Datei) dupliziert sie — die Praxis
  seit `v0.13.1` hat das viermal gezeigt und der Beschluss setzt sie als Regel. Die
  Archivierung ersetzte `done/`-Inhalte durch Stubs und liefe gegen die Sichtbarkeit des
  Volltexts im Arbeitsbaum.
- **Aufloesungs-Trigger:** permanent, solange Release-Cuts wellenlos laufen.