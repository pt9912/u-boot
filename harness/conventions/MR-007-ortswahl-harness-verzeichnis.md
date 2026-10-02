# MR-007 — Ortswahl `.harness/` (dot-prefixed, committet) neben `harness/`

- **Datum:** 2026-07-24
- **Geltungsbereich:** `.harness/` (vendored Baseline, kuenftig
  `.harness/skills/`), `harness/` (Autoren-Prosa), `.gitignore`.
- **Ersetzt-Baseline-Regel:** [`grundlagen-harness-dateien.md` §Verzeichniskonvention](../../.harness/baseline/v6.13.0/regelwerk/grundlagen-harness-dateien.md#verzeichniskonvention)
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
