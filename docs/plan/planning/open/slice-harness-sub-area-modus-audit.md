# Slice Harness: Sub-Area-Modus-Audit des Bestandscodes (FS-3)

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** Konformitätslücke der Regelwerk-Adoption, dokumentiert als FS-3 in
[`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md)
§7/§9 und als Auflösungs-Trigger in `MR-006`
([`harness/conventions.md`](../../../../harness/conventions.md)). Berührt
[`LH-FA-ARCH-002`](../../../../spec/lastenheft.md#lh-fa-arch-002--schichten-und-verzeichnislayout)
nur lesend (Verzeichnislayout als Ausgangsinventar), ändert keine Anforderung.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die pragmatische Erst-Einordnung des Bestandscodes („`hexagon/`, `cmd/`,
`internal/` sind GF") durch ein **vollständiges Audit** ersetzen: erst
Inklusions-Prüfung (drei Achsen, Schwelle ≥ 2 — welche Struktur *ist* überhaupt
eine Sub-Area?), dann pro qualifizierter Sub-Area die vier Modus-Pflichtkriterien
(Konventionen-Dichte · Phase-Reife · Evidenz-/Diskrepanz-Risiko ·
Reconciliation-Aufwand) mit GF/BF/Hybrid-Diagnose und — bei BF/Hybrid — einer
benannten Graduation-Bedingung.

Der Erst-Pass war bewusst grob und hat die drei Top-Level-Verzeichnisse als
Sub-Areas geführt. Das ist nach der Aggregations-/Inklusions-Logik des Regelwerks
zu grob (`internal/` bündelt Adapter, Hexagon-Kern und Test-Infrastruktur mit je
eigener Konventions-Härte). Dieser Slice differenziert aus und macht die
Modus-Tabelle in `harness/conventions.md` belastbar.

## 2. Definition of Done

- [ ] **Kandidaten-Inventar** aller Code-/Tooling-Strukturen aufgestellt
  (mindestens: `internal/hexagon/domain`, `internal/hexagon/application`,
  `internal/hexagon/port`, `internal/adapter/driving/cli`,
  `internal/adapter/driven/*`, `cmd/uboot`, `internal/e2e`,
  `tools/` + `scripts/`).
- [ ] **Inklusions-Prüfung je Kandidat** dokumentiert: Achse 1
  (Konventions-Härte — eigene `MR-NNN` plausibel?), Achse 2 (Inventur-Linie —
  eigene Diskrepanz-Zeile sinnvoll?), Achse 3 (struktureller Cluster — eigene
  Pfad-/Datei-Familie?), Schwelle ≥ 2. Nicht-qualifizierte Kandidaten werden
  als **Sub-Area-Aspirantin** benannt und der aufnehmenden Sub-Area zugeordnet
  — Aggregation ist ein Ergebnis, kein Versäumnis.
- [ ] **Vier Pflichtkriterien je qualifizierter Sub-Area** ausgefüllt, inklusive
  Phase-Reife (0–5) gegen die Phase-×-Modus-Matrix und expliziter
  GF/BF/Hybrid-Aussage.
- [ ] **Jede BF-/Hybrid-Markierung trägt eine Graduation-Bedingung** (Trigger-
  Klasse oder Folge-Slice); jede GF-Aussage trägt den Beleg, dass Doku vor Code
  lag. Keine Modus-Aussage ohne Beleg — „GF, weil neu" ist kein Beleg.
- [ ] **`harness/conventions.md` §Modus-Deklaration ersetzt**: die Vier-Zeilen-
  Erst-Pass-Tabelle weicht der auditierten Tabelle; der Erst-Pass-Vorbehalt
  („bewusst pragmatisch", „Folge-Slice") wird eingelöst, nicht bloß umformuliert.
  `MR-006` bekommt den erledigten Auflösungs-Trigger mit Delivery-Verweis.
- [ ] **Audit-Detailtabelle** (drei Achsen je Kandidat) lebt in **diesem Slice**,
  nicht in `conventions.md` — dort steht nur das Ergebnis. Begründung:
  `conventions.md` ist Regelwerks-Deklaration, nicht Audit-Protokoll.
- [ ] `make docs-check` grün.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| dieser Slice §Audit | neu | Audit-Protokoll (drei Achsen + vier Kriterien je Kandidat) |
| [`harness/conventions.md`](../../../../harness/conventions.md) §Modus-Deklaration | update | Ergebnis-Tabelle ersetzt den Erst-Pass |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-006` | update | Auflösungs-Trigger eingelöst, Verweis auf dieses Audit |

Kein Code-Delta: Das Audit beschreibt den Bestand, es ändert ihn nicht.

## 4. Trigger

Nach FS-1/FS-2
([`slice-harness-reviewer-skills-und-review-ablage`](slice-harness-reviewer-skills-und-review-ablage.md)),
weil der Reviewer-Skill die Modus-Aussagen als Prüf-Kontext nutzt. Rückführung
nach `open/`, falls das Audit eine echte BF-Sub-Area mit nennenswertem
Reconciliation-Aufwand aufdeckt — dann wird die Reconciliation ein eigener
Slice, nicht ein Anhang an dieses Audit.

## 5. Closure-Trigger

Alle Kandidaten qualifiziert oder als Aspirantin abgewiesen, Modus-Tabelle in
`conventions.md` ersetzt, `make docs-check` grün, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Bestätigungs-Audit statt Audit:** Das Risiko ist, dass das Audit den
  Erst-Pass nur nachbetet („alles GF"). Gegenmaßnahme: Jede GF-Aussage braucht
  einen *positiven* Beleg (Spec-/Architektur-Stelle, die dem Code vorausging),
  nicht nur das Fehlen von Gegenanzeichen.
- **Retrofit-verdächtige Bereiche:** Kandidaten mit nachträglich eingezogener
  Konvention (Exit-Code-Klassifikation, JSON-Envelope, Coverage-Bootstrap) sind
  die realistischen Hybrid-Anwärter. Sie sind explizit zu prüfen, nicht zu
  überspringen.
- **Über-Differenzierung:** Zu feine Sub-Areas erzeugen Inventur-Linien ohne
  eigene Diskrepanz (Anti-Refactoring). Die Aggregations-Regel (gleiche Trigger
  *und* gleiche Modus-Aussage → eine Sub-Area) ist gegenzuhalten.
- **Kein Carveout erwartet:** reine Doku-/Deklarations-Arbeit ohne Gate-Bezug.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Schreibend berührt wird nur die Sub-Area *harness / Konventionen* (**GF**,
Doku-führt). Die auditierten Code-Sub-Areas werden **lesend** inventarisiert —
ihre Modus-Aussage ist das *Ergebnis* dieses Slice und steht in §Audit bzw. in
[`harness/conventions.md`](../../../../harness/conventions.md).
