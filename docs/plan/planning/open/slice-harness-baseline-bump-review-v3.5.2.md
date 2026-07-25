# Slice Harness: Review-Bump der Baseline v3.5.1 → v3.5.2

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** ohne Welle (Harness-Wartung).

**Bezug:** Ausgelöst vom **ersten** Freshness-Audit-Lauf
(`tools/harness/fetch-baseline-cache.sh --check-freshness`, 2026-07-25,
Exit 3): Der Kurs `pt9912/ai-harness-course` trägt den Release-Tag `v3.5.2`,
der lokale Pin steht auf `v3.5.1`. Kein `LH`-/`ADR`-Bezug — reine
Baseline-Wartung nach der `MR-004`-Bump-Prozedur
([`harness/conventions.md`](../../../../harness/conventions.md)).

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Den Versions-Bump der vendorten Baseline als **eine Einheit** ausführen —
inklusive der Lese-Leistung, die ihn von einem Auto-Update unterscheidet:
Was hat sich zwischen `v3.5.1` und `v3.5.2` geändert, und macht eine dieser
Änderungen eine u-boot-Adaption (`MR-000`..`MR-009`) ungültig?

## 2. Definition of Done

- [ ] **Delta gelesen, nicht nur gezogen:** Die Änderungen zwischen den beiden
  Tags sind gesichtet und in einer Kurz-Liste festgehalten (welche Module,
  welche Templates, welche Regel-Änderungen).
- [ ] **Adaptions-Ledger gegengeprüft:** Für jede `MR-*`-Adaption ist
  festgehalten, ob sie unter `v3.5.2` unverändert gilt, angepasst werden muss
  oder entfällt. Ein Bump ohne diesen Abgleich wäre genau das Auto-Update, das
  die Routine verhindern soll.
- [ ] **Bump als Einheit ausgeführt** (`MR-004`): (1) `**Stand:**`-Pin,
  (2) Vendor-Pfad `.harness/baseline/v3.5.2/` per Skript-Lauf,
  (3) `AGENTS.md`-Pointer, (4) `harness/README.md`-Guides-Zeile.
- [ ] **Alt-Stand behandelt:** entweder `.harness/baseline/v3.5.1/` entfernt
  (Vendor-Pfad ist tag-versioniert, der Glob `.harness/baseline/**` bleibt
  tag-agnostisch) oder bewusst als Referenz behalten — mit Begründung.
- [ ] `--verify` gegen den neuen Stand grün (42+ Dateien, Manifest-Deckung).
- [ ] `--check-freshness` danach Exit 0.
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/v3.5.2/` | neu (Skript-Lauf) | vendorter Stand |
| `.harness/baseline/v3.5.1/` | entfernen oder begründet behalten | ein Stand pro Repo ist die Regel |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | Pin, Adoptionsdatum, `MR-*`-Abgleich |
| [`AGENTS.md`](../../../../AGENTS.md), [`harness/README.md`](../../../../harness/README.md) | update | T1-/T2-Pointer |

## 4. Trigger

Bereits gefeuert (Freshness-Audit 2026-07-25). Einplanung ist eine
Priorisierungs-Entscheidung, kein weiterer Trigger.

## 5. Closure-Trigger

Bump vollständig, `--verify` und `--check-freshness` grün, `MR-*`-Abgleich
dokumentiert, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Stiller Regel-Wechsel:** Eine Regeländerung in `v3.5.2` kann eine
  bestehende Adaption entwerten, ohne dass ein Gate ausschlägt. Genau dagegen
  steht der `MR-*`-Abgleich in der DoD.
- **Template-Drift:** Geänderte Templates wirken auf künftige Artefakte, nicht
  rückwirkend. Bestehende Artefakte werden **nicht** migriert; falls das
  Regelwerk das fordert, ist es ein eigener Slice (und bei Vertrags-Berührung
  ein Change Request, vgl. `MR-008`).
- **Bundle-Layout:** Vor dem Entpack-Lauf gilt weiterhin Schritt 0 —
  Asset-Name und innerer Tree gegen den echten Release prüfen.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen / Baseline* und
*Harness-Tooling* — beide **GF** nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
