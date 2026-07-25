# Slice Harness: Freshness-Audit-Routine der Baseline (FS-4)

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** FS-4 aus
[`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md)
§7; löst die in `MR-004`
([`harness/conventions.md`](../../../../harness/conventions.md)) offen
formulierte Zusage ein („Ein neuer Kurs-Tag wird über die Release-**Liste**
erkannt (Freshness-Audit) und löst einen Review-Bump aus, keinen Auto-Update").
Kein `LH`-/`ADR`-Neubezug.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die Freshness-Prüfung der vendorten Baseline von einer *Absichtserklärung* in
eine **ausführbare, terminierte Routine** überführen: eine dokumentierte
Prüf-Kadenz plus ein Skript-Modus, der die Release-**Liste** des Kurs-Repos
gegen den lokalen `**Stand:**`-Pin hält und einen **Review-Bump meldet** —
ohne je selbst zu aktualisieren.

Die Trennung ist der Kern: Der Integritäts-Pin (`SHA256SUMS`) beantwortet „ist
der vendorte Bestand unversehrt?", das Freshness-Audit beantwortet „ist der
gepinnte Stand noch der aktuelle?". Der zweite Sensor fehlt bisher; ohne ihn
altert die Baseline still.

## 2. Definition of Done

- [ ] **`tools/harness/fetch-baseline-cache.sh --check-freshness`**: liest die
  Release-Liste von `pt9912/ai-harness-course`, vergleicht sie mit dem
  `**Stand:**`-Pin und meldet neuere Tags. **Read-only** — kein Schreibzugriff
  auf `.harness/baseline/`, kein Vendoring, kein Pin-Update. Der Modus ist
  bewusst kein Gate in `make gates`: ein neuer Kurs-Release darf u-boots CI
  nicht rot färben.
- [ ] **Exit-Code-Semantik dokumentiert und implementiert:** `0` = Pin ist der
  neueste Tag, `3` = neuerer Tag vorhanden (Review-Bump fällig), `1` =
  Ausführungsfehler (Netz, fehlendes Werkzeug, unlesbarer Pin). Der
  Unterschied zwischen „aktuell" und „veraltet" muss maschinenlesbar sein,
  sonst ist die Routine nicht automatisierbar.
- [ ] **Kadenz dokumentiert** in [`harness/conventions.md`](../../../../harness/conventions.md)
  (eigener Abschnitt Freshness-Audit): *wann* geprüft wird (bei jedem
  Harness-/Baseline-Slice sowie mindestens quartalsweise), *wer* prüft
  (die Rolle, die den Harness-Slice führt), *was* ein positiver Befund auslöst
  (Review-Bump-Slice nach `MR-004`-Bump-Prozedur, kein Auto-Update) und *was
  nicht* (kein automatischer Vendor-Lauf, kein CI-Gate).
- [ ] **Abgrenzung zu `--verify` explizit** in Skript-Kopf und Doku: Integrität
  ≠ Aktualität; `--verify` bleibt offline, `--check-freshness` braucht Netz.
- [ ] **`MR-004` eingelöst**: Der offene Satz zur Release-Liste verweist auf die
  jetzt existierende Routine statt sie nur anzukündigen.
- [ ] Erster Lauf ausgeführt und sein Ergebnis in der Closure-Evidence
  festgehalten (auch ein Negativbefund „Pin ist aktuell" ist Evidence).
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/fetch-baseline-cache.sh` | update | dritter Modus `--check-freshness` (read-only, Netz) |
| [`harness/conventions.md`](../../../../harness/conventions.md) §Freshness-Audit | neu | Kadenz, Zuständigkeit, Auslöser, Nicht-Ziele |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-004` | update | Zusage eingelöst, Verweis auf die Routine |
| [`harness/README.md`](../../../../harness/README.md) | update falls nötig | Skript-Modi-Zeile konsistent halten |

## 4. Trigger

Nach FS-3
([`slice-harness-sub-area-modus-audit`](slice-harness-sub-area-modus-audit.md)),
weil beide `harness/conventions.md` anfassen und eine serielle Bearbeitung
Merge-Reibung spart. Fachlich unabhängig — die Reihenfolge ist Ökonomie, keine
Abhängigkeit.

## 5. Closure-Trigger

Skript-Modus lauffähig und einmal gelaufen, Kadenz dokumentiert, `MR-004`
eingelöst, `make docs-check` grün, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Netz-Abhängigkeit:** Der Modus braucht GitHub-Erreichbarkeit. Er darf
  deshalb nirgends in einem Pflicht-Gate hängen; ein Fehlschlag ist ein
  Audit-Befund, kein Build-Fehler.
- **API-Form-Drift:** Die Release-Liste wird über die GitHub-API bzw. deren
  öffentliche Auflistung gelesen. Ändert sich die Form, meldet der Modus einen
  Ausführungsfehler (Exit 1) statt still „alles aktuell" — Fail-loud ist hier
  Pflicht, weil ein stiller Falsch-Negativ die Routine wertlos macht.
- **Kadenz ohne Erinnerung:** Eine dokumentierte Kadenz ohne Auslöser wird
  vergessen. Gegenmaßnahme: Die Kadenz hängt primär am *Ereignis*
  (jeder Harness-/Baseline-Slice prüft), nicht nur am Kalender.
- **Kein Carveout erwartet:** additive Tooling-/Doku-Arbeit ohne Gate-Lockerung.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* (`tools/`) —
beide **GF** (Doku-/Konvention-führt; das Skript materialisiert eine
Konvention, die zuerst geschrieben wurde). Einordnung nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration in der von
[`slice-harness-sub-area-modus-audit`](slice-harness-sub-area-modus-audit.md)
auditierten Fassung.
