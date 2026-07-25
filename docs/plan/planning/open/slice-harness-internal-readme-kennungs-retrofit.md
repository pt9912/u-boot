# Slice Harness: `internal/`-READMEs in den Referenz-Scan holen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** ohne Welle (Harness-Wartung) — Folge-Befund aus
[`slice-harness-sub-area-modus-audit`](slice-harness-sub-area-modus-audit.md)
§9.2.

**Bezug:** [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(Kennungs-/Link-Disziplin) und
[`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)
(die `scan.ignore`-Ausnahme ist ein Carveout und braucht einen Plan — dieser
hier). Berührt
[`ADR-0013`](../../adr/0013-dokumentationsreferenzmodell.md) lesend.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die Code-Paket-READMEs unter `internal/**` aus dem toten Winkel holen:
`LH-*`-Kennungen verlinken, den zeitlichen Ballast (Meilenstein-/Tranchen-Tags
in Status-Abschnitten) entfernen und den `scan.ignore`-Glob `internal/**` in
[`.d-check.yml`](../../../../.d-check.yml) auflösen. Damit bekommt die
Inventur-Linie „README-Aussage vs. Code-Bestand" erstmals einen Sensor.

Der Befund stammt aus dem Sub-Area-Audit: Diese READMEs sind die einzige
**Brownfield**-Sub-Area des Repos (Doku beschreibt Bestand nachträglich), und
ihre Ausnahme vom Referenz-Gate ist bisher nur im Konfigurations-Kommentar
begründet — ohne Inventar-Eintrag, ohne Plan.

## 2. Definition of Done

- [ ] **Kennungs-Inventar:** alle ungelinkten `LH-*`-/`ADR-*`-Kennungen in
  `internal/**/README.md` erfasst (Ist-Zahl im Slice festgehalten, nicht nur
  „einige").
- [ ] **Links gesetzt** — relative Pfade auf
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) bzw. den ADR-Index,
  Anker-Form wie im übrigen Repo.
- [ ] **Zeitliche Schicht entfernt:** Status-Abschnitte tragen keine
  Meilenstein-/Tranchen-Tags (`M3-T2`, `MVP-Closure-T1`, „Stand M8") mehr,
  sondern beschreiben den Ist-Zustand. Die zeitliche Schicht lebt in Roadmap
  und Closure-Notizen, nicht im Code-README.
- [ ] **`scan.ignore` verengt:** `internal/**` entfällt; falls einzelne
  Restfälle bleiben müssen, wird der Glob so eng wie möglich gefasst und der
  Rest als Carveout mit Plan geführt — kein pauschaler Verzeichnis-Ausschluss.
- [ ] **Carveout-Eintrag aufgelöst:** der Eintrag zum `internal/**`-Ausschluss
  im Inventar wird geschlossen (Audit-Trail bleibt).
- [ ] `make docs-check` grün mit dem erweiterten Scope (Datei-Zahl steigt
  sichtbar — das ist die Evidence).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/**/README.md` (7 Dateien) | update | Kennungen verlinken, Status entzeitlichen |
| [`.d-check.yml`](../../../../.d-check.yml) `scan.ignore` | update | `internal/**` entfernen bzw. minimal fassen |
| [`carveouts.md`](../in-progress/carveouts.md) | update | Eintrag auflösen mit Audit-Trail |

## 4. Trigger

Kein externer Trigger; direkt einplanbar. Sinnvoll gebündelt mit der nächsten
Doku-/Harness-Welle, weil der Diff breit (7 Dateien), aber flach ist.

## 5. Closure-Trigger

`make docs-check` grün ohne `internal/**`-Ausschluss, Carveout-Eintrag
aufgelöst, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Befund-Lawine:** Der erweiterte Scan kann mehr Befunde melden als erwartet
  (Anker-Drift in älteren READMEs). Gegenmaßnahme: erst Trockenlauf mit
  entferntem Ignore, Ist-Zahl notieren, dann entscheiden, ob der Slice in
  Tranchen zerfällt.
- **Entzeitlichen vs. Informationsverlust:** Die Meilenstein-Tags tragen
  historischen Kontext. Sie werden nicht ersatzlos gelöscht, sondern durch
  Ist-Aussagen ersetzt; die Historie steht ohnehin in den `done/`-Slices.
- **Carveout bleibt möglich:** Falls ein Restbestand nicht sinnvoll linkbar ist
  (z. B. Kennungen in Code-Blöcken), bleibt ein enger, begründeter Ausschluss —
  dann aber inventarisiert und mit Plan-Anker.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

### Sub-Area: `internal/**/README.md` (Code-Paket-READMEs)

- **Modus:** BF (Doku beschreibt den Code-Bestand nachträglich).
- **Konventionen-Dichte:** niedrig — keine Form-Regel für diese READMEs
  vorhanden; dieser Slice zieht die erste ein (Kennungs-Linkpflicht,
  keine zeitliche Schicht).
- **Phase-Reife:** 3 (partiell) → Ziel 4 (kohärent, mit Sensor).
- **Evidenz-/Diskrepanz-Risiko:** hoch — genau das macht dieser Slice
  sichtbar; die Ist-Zahl der Befunde ist das erwartete Diskrepanz-Signal.
- **Reconciliation-Aufwand:** ein Slice, ggf. zwei Tranchen (Kennungen /
  Entzeitlichung). Graduation nach GF, sobald `internal/**` im Scan liegt.
