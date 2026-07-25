# Slice Harness: Freshness-Audit-Routine der Baseline (FS-4)

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** FS-4 aus
[`slice-harness-regelwerk-adoption-v3.5.1`](slice-harness-regelwerk-adoption-v3.5.1.md)
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

- [x] **`tools/harness/fetch-baseline-cache.sh --check-freshness`**: liest die
  Release-Liste von `pt9912/ai-harness-course`, vergleicht sie mit dem
  `**Stand:**`-Pin und meldet neuere Tags. **Read-only** — kein Schreibzugriff
  auf `.harness/baseline/`, kein Vendoring, kein Pin-Update. Der Modus ist
  bewusst kein Gate in `make gates`: ein neuer Kurs-Release darf u-boots CI
  nicht rot färben.
- [x] **Exit-Code-Semantik dokumentiert und implementiert:** `0` = Pin ist der
  neueste Tag, `3` = neuerer Tag vorhanden (Review-Bump fällig), `1` =
  Ausführungsfehler (Netz, fehlendes Werkzeug, unlesbarer Pin). Der
  Unterschied zwischen „aktuell" und „veraltet" muss maschinenlesbar sein,
  sonst ist die Routine nicht automatisierbar.
- [x] **Kadenz dokumentiert** in [`harness/conventions.md`](../../../../harness/conventions.md)
  (eigener Abschnitt Freshness-Audit): *wann* geprüft wird (bei jedem
  Harness-/Baseline-Slice sowie mindestens quartalsweise), *wer* prüft
  (die Rolle, die den Harness-Slice führt), *was* ein positiver Befund auslöst
  (Review-Bump-Slice nach `MR-004`-Bump-Prozedur, kein Auto-Update) und *was
  nicht* (kein automatischer Vendor-Lauf, kein CI-Gate).
- [x] **Abgrenzung zu `--verify` explizit** in Skript-Kopf und Doku: Integrität
  ≠ Aktualität; `--verify` bleibt offline, `--check-freshness` braucht Netz.
- [x] **`MR-004` eingelöst**: Der offene Satz zur Release-Liste verweist auf die
  jetzt existierende Routine statt sie nur anzukündigen.
- [x] Erster Lauf ausgeführt und sein Ergebnis in der Closure-Evidence
  festgehalten (auch ein Negativbefund „Pin ist aktuell" ist Evidence).
- [x] `make docs-check` grün.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

### Verification Evidence

Scope:
- Slice: `slice-harness-baseline-freshness-audit` (FS-4)
- IDs: **keine** Anforderung geändert — Harness-Tooling ohne Produktvertrag.
- Artefakte: `tools/harness/fetch-baseline-cache.sh` (dritter Modus
  `--check-freshness` + Kopf-Doku),
  [`harness/conventions.md`](../../../../harness/conventions.md)
  (§Freshness-Audit neu, `MR-004` eingelöst),
  [`AGENTS.md`](../../../../AGENTS.md) (Pointer auf die Ereignis-Kadenz),
  [`slice-harness-baseline-bump-review-v3.5.2`](../open/slice-harness-baseline-bump-review-v3.5.2.md)
  (ausgelöster Folge-Plan), [`roadmap.md`](../in-progress/roadmap.md).

DoD-Abgleich: alle Punkte erfüllt. Der Plan-Eintrag
„[`harness/README.md`](../../../../harness/README.md) update **falls nötig**"
löst sich zu „nicht nötig" auf: Die Datei trägt keine Skript-Modi-Zeile, nur den
T1-Pointer — vom Review gegengeprüft.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `bash -n tools/harness/fetch-baseline-cache.sh` | pass | Syntax ok |
| `fetch-baseline-cache.sh --check-freshness` | **Exit 3** | erster Lauf meldet `v3.5.2` gegen Pin `v3.5.1`; Arbeitsbaum danach unverändert (`git status` leer) — Read-only bestätigt |
| `fetch-baseline-cache.sh --verify` | pass | 42 Dateien, vollständig, offline — Regression des bestehenden Modus |
| `make docs-check` | pass | 129 Dateien / 0 Befunde |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| `MR-004` Zusage „neuer Tag über die Release-**Liste**, Review-Bump statt Auto-Update" | ausführbarer Modus + §Freshness-Audit (Sensor, Exit-Codes, Kadenz, Zuständigkeit, Auslöser, Nicht-Ziele) |
| Regelwerk `modul-02` (Baseline-Aktualität) | Trennung Integrität (`--verify`, offline) vs. Aktualität (`--check-freshness`, Netz) in Skript-Kopf, `AGENTS.md` und `conventions.md` |
| Regelkonforme Reaktion auf den Erst-Befund | kein Pin-Delta im Diff; stattdessen `open/`-Plan mit Delta-Lektüre und `MR-*`-Gegenprobe **vor** dem Bump |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta; das Skript ist Harness-Tooling
  außerhalb der Build-Pipeline.
- Inhaltlicher Abgleich `v3.5.1` gegen `v3.5.2` — ausdrücklich **nicht** Teil
  dieses Slice; das ist der Auftrag des ausgelösten Folge-Plans. Ein Bump im
  selben Zug wäre genau das Auto-Update, das die Routine verhindern soll.

Independent Review (Frischkontext, Rollentrennung): kein Befund an diesem Slice.
Read-only-Zusage statisch und dynamisch gegengeprüft, Fail-loud-Semantik
inklusive der Sonderfälle (leere Antwort, Pin nicht in der Liste, fehlendes
`curl`) bestätigt, Exit-Code-Propagation über `check_freshness || rc=$?`
verifiziert, Gate-Politik als bewusstes Nicht-Ziel anerkannt. Zwei INFO ohne
Änderungsbedarf: Exit `3` kollidiert zahlenmäßig mit dem für das *Produkt*
reservierten Bereich (bindet Harness-Skripte nicht), und die Skript-Ausgaben
sind deutsch wie die der bestehenden Modi (
[`LH-LESE-002`](../../../../spec/lastenheft.md#lh-lese-002--sprache) bindet
Produkt-Ausgaben, nicht Harness-Tooling).

Commit / Artefakt: `c7fe437`; Lifecycle-Move `open/` → `done/` im Folge-Commit.

### Steering-Loop-Lerneintrag

- **Ein Sensor, der beim ersten Lauf anschlägt, war überfällig.** Die Routine
  fand sofort einen neueren Tag — die Baseline war also bereits veraltet, ohne
  dass es jemandem aufgefallen wäre. Das ist das stärkste Argument gegen
  „Kadenz später": Der Sensor hat sich mit seinem ersten Lauf bezahlt gemacht.
- **Fail-loud ist bei Freshness kein Luxus.** Ein Netz- oder Formatfehler, der
  als „alles aktuell" durchginge, wäre schlimmer als kein Sensor — er erzeugte
  belegte Sicherheit statt bekannter Unsicherheit. Deshalb `exit 1` statt
  stillem `return 0`, auch im Sonderfall „Pin nicht in der Liste".
- **Die Ereignis-Kadenz trägt, die Kalender-Kadenz sichert nur ab.** Eine reine
  Quartals-Regel wird vergessen; die Kopplung an „jeder Harness-Slice prüft"
  hat einen natürlichen Auslöser und wurde deshalb ins Agenten-Briefing gehoben.
- **Folge-Slices:**
  [`slice-harness-baseline-bump-review-v3.5.2`](../open/slice-harness-baseline-bump-review-v3.5.2.md)
  — vom Sensor selbst ausgelöst, bewusst nicht in dieser Welle ausgeführt.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* (`tools/`) —
beide **GF** (Doku-/Konvention-führt; das Skript materialisiert eine
Konvention, die zuerst geschrieben wurde). Einordnung nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration in der von
[`slice-harness-sub-area-modus-audit`](slice-harness-sub-area-modus-audit.md)
auditierten Fassung.
