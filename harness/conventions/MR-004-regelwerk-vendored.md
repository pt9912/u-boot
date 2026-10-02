# MR-004 — Regelwerk-Lese-Form committet vendored; Baseline-Pin v6.13.0; beide Baeume

- **Datum:** 2026-07-24
- **Geltungsbereich:** `.harness/baseline/`,
  `tools/harness/fetch-baseline-cache.sh`, [`AGENTS.md`](../../AGENTS.md) §1,
  [`harness/README.md`](../README.md) §Guides, [`.d-check.yml`](../../.d-check.yml)
  (`scan.ignore`), `.gitignore`, Abschnitt Baseline oben.
- **Ersetzt-Baseline-Regel:** [`modul-02-harness-bootstrap.md` §Greenfield-Bootstrap: Schritt-Sequenz (Modul 2)](../../.harness/baseline/v6.13.0/regelwerk/modul-02-harness-bootstrap.md#greenfield-bootstrap-schritt-sequenz-modul-2)
- **Adaption:** Die Lese-Form des adoptierten Regelwerks ist **committet
  vendored** (kein Remote-ZIP pro Lauf, kein Submodule):
  `.harness/baseline/v6.13.0/{regelwerk,templates}/` + `SHA256SUMS`, netzlos auf
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

### Sync-Trigger (T1/T2)

Begriffe aus dem Regelwerk (`modul-02-harness-bootstrap`): Ein Pointer auf die
vendored Baseline muss an **zwei** Stellen synchron gehalten werden.

- **T1** - Pointer in [`harness/README.md`](../README.md) §Guides (Verweis auf das
  vendored Regelwerk und diese Datei).
- **T2** - Pointer in der Source-Precedence-/Kopf-Sektion von
  [`AGENTS.md`](../../AGENTS.md) (Verweis auf die vendored Baseline + Lesemodell).

Beide zeigen auf `.harness/baseline/v6.13.0/regelwerk/README.md` (Index) und
werden bei einem Baseline-Bump gemeinsam nachgezogen (`MR-004` Bump-Prozedur).
Fundstelle: `.harness/baseline/v6.13.0/regelwerk/modul-02-harness-bootstrap.md`.

### Freshness-Audit

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
