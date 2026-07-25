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
ihre Modus-Aussage ist das *Ergebnis* dieses Slice und steht in §9 bzw. in
[`harness/conventions.md`](../../../../harness/conventions.md).

## 9. Audit-Protokoll

Durchgeführt 2026-07-25 gegen den Stand `d5d896a`. Das Protokoll lebt hier,
nicht in `harness/conventions.md` — dort steht nur das Ergebnis.

### 9.1 Inklusions-Prüfung (drei Achsen, Schwelle ≥ 2)

Achse 1 = eigene `MR-NNN`-Adaption plausibel formulierbar; Achse 2 = eigene
Diskrepanz-/Inventur-Zeile sinnvoll, ohne Nachbar mitzuziehen; Achse 3 = eigene
Pfad-/Datei-Familie.

| Kandidat | A1 | A2 | A3 | Ergebnis |
|---|---|---|---|---|
| `internal/hexagon/domain` | ✅ I/O-Freiheit, Value-Object-Pflicht | ✅ „anämischer Domänentyp" ist eine eigene Diskrepanz | ✅ | **Sub-Area** (3/3) |
| `internal/hexagon/application` | ✅ ein Service pro Use-Case-Familie; nil-tolerante Ports via `noop*`-Defaults | ✅ Use-Case ohne Port-Abstraktion | ✅ | **Sub-Area** (3/3) |
| `internal/hexagon/port` (`driving` + `driven`) | ✅ Kreuz-Blindheit; Sentinel-Definition liegt im `driving`-Port | ✅ Port-/Adapter-Drift | ✅ | **Sub-Area** (3/3) |
| `internal/adapter/driving/cli` | ✅ Exit-Code-Klassifikation, Dual-Classifier-Regel, JSON-Envelope-Form | ✅ Sentinel ohne Klassifikations-Eintrag | ✅ | **Sub-Area** (3/3) |
| `internal/adapter/driven` | ✅ Port-Pin `var _ driven.X = (*Adapter)(nil)` im Produktivcode | ✅ Adapter ohne Pin / ohne Port | ✅ | **Sub-Area** (3/3) |
| `cmd/uboot` | ❌ „nur Wiring" ist eine Architektur-Regel, keine eigene `MR` | ❌ eine Verletzung ist ein Import-Regel-Befund und gehört zur Hexagon-Linie | ✅ eigenes Verzeichnis | **Aspirantin** (1/3) — 2 Dateien, 302 LOC; aggregiert in die Hexagon-Schichtungs-Linie |
| `internal/e2e` (Test-Infrastruktur) | ✅ Build-Tag-Konvention `//go:build docker`, „kein reales `time.Sleep` in Tests" | ✅ Acceptance-Test ohne `LH-*`-Anker | ✅ | **Sub-Area** (3/3) |
| `tools/` + `scripts/` (Harness-Tooling) | ✅ self-contained Bash, Docker-only-Harness, `MR-004`/`MR-005` betreffen sie direkt | ✅ Skript-Modus vs. dokumentierte Zusage | ✅ | **Sub-Area** (3/3) |
| `internal/**/README.md` (Code-Paket-READMEs) | ✅ eigene Form-Regel plausibel (Status-/Inventar-Abschnitt, Kennungs-Linkpflicht) | ✅ **eigene Linie:** README-Aussage vs. Code-Bestand — heute ohne Sensor, weil `internal/**` im `scan.ignore` liegt | ✅ eigenes Dateimuster | **Sub-Area** (3/3) — im Erst-Pass **nicht** sichtbar |
| `internal/hexagon/application/{templates,managedblock}` | ❌ keine von `application` getrennte Konvention | ❌ keine eigenständig abgleichbare Linie | ✅ | **Aspirantin** (1/3) — Teil von `application` |
| `internal/adapter/driven/<einzelner Adapter>` | ❌ | ❌ | ✅ | **Aspirantin** (1/3) — Aggregation nach `driven` (gleiche Trigger, gleiche Modus-Aussage) |

**Ergebnis der Inklusion:** acht qualifizierte Sub-Areas statt der drei
Erst-Pass-Einträge. Zwei Korrekturen gegenüber dem Erst-Pass: `cmd/uboot` fällt
auf Aspirantin zurück (Struktur ohne eigene Substanz), und die **Code-Paket-
READMEs** kommen als eigene Sub-Area hinzu — sie tragen eine Inventur-Linie, die
der Erst-Pass gar nicht gesehen hat.

### 9.2 Modus-Diagnose (vier Pflichtkriterien)

Nur die Sub-Areas mit nicht-trivialem Befund sind hier ausgeführt; die
GF-Diagnose der reinen Hexagon-Schichten ist in der Ergebnis-Tabelle
zusammengefasst.

#### `internal/adapter/driving/cli` — **Hybrid**

- **Konventionen-Dichte:** gemischt. Exit-Code-Vertrag und maschinenlesbare
  Ausgabe sind im Lastenheft verankert
  ([`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe)).
  Die zwei **tragenden Implementierungs-Konventionen** sind es nicht:
  (a) **Sentinel-Schichtung** — Driven-Sentinels werden vor Driving-Sentinels
  klassifiziert; (b) **Dual-Classifier-Regel** — ein Sentinel muss in der
  Exit-Code-Klassifikation *und* in der Envelope-/Diagnostic-Abbildung stehen.
  Beide leben ausschließlich als Code-Kommentar im Adapter, entstanden aus
  Slice-Arbeit und Review-Runden — also **Code → Doku**.
- **Phase-Reife:** 4 (kohärent). Der Bestand ist vollständig und konsistent,
  aber die Sicht-Spec kennt die beiden Regeln nicht — genau die Diskrepanz, die
  Phase 4 im BF-Lesart sichtbar macht.
- **Evidenz-/Diskrepanz-Risiko:** mittel. Ohne Doku-Anker ist die Regel nur so
  lange wirksam, wie jemand den Kommentar liest; ein neuer Sentinel-Split ist
  der realistische Fehlerfall.
- **Reconciliation-Aufwand / Graduation:** klein und **bereits eingeplant** —
  [`slice-harness-architecture-template-konformitaet`](slice-harness-architecture-template-konformitaet.md)
  nimmt beide Regeln in den neuen Abschnitt Fehlermodelle auf. Danach steht die
  Doku vor der nächsten Änderung → **Graduation nach GF**.

#### `internal/**/README.md` — **Brownfield**

- **Konventionen-Dichte:** niedrig. Es gibt keine Form-Regel für diese READMEs:
  sie tragen Status-Abschnitte mit Meilenstein-/Tranchen-Tags („Stand M8",
  „M3-T2", „MVP-Closure-T1") und `LH-*`-Kennungen **ohne Links**.
- **Phase-Reife:** 3 (partiell). Manche Pakete dokumentieren ihren Stand
  detailliert, andere knapp; die Aktualität hängt am jeweiligen Slice-Autor.
- **Evidenz-/Diskrepanz-Risiko:** **hoch, und heute unbeobachtet.** `internal/**`
  liegt im `scan.ignore` von [`.d-check.yml`](../../../../.d-check.yml) („historisch
  ungelinkte LH-Kennungen, dokumentierter Migrations-Restbestand"). Die
  Ausnahme ist im Konfigurations-Kommentar begründet, steht aber **weder im
  Carveout-Inventar noch trägt sie einen Plan-Anker** — nach
  [`LH-FA-PROJDOCS-005`](../../../../spec/lastenheft.md#lh-fa-projdocs-005--carveout-disziplin)
  ist genau das die Lücke. Zweitbefund: Die Meilenstein-Tags machen diese
  Dateien zu einer *zeitlichen Schicht* neben Roadmap und Closure-Notizen.
- **Reconciliation-Aufwand / Graduation:** ein eigener Slice
  ([`slice-harness-internal-readme-kennungs-retrofit`](slice-harness-internal-readme-kennungs-retrofit.md)):
  Kennungen verlinken, Status-Abschnitte entzeitlichen, `internal/**` aus dem
  `scan.ignore` nehmen. Danach existiert ein Sensor für die Inventur-Linie →
  **Graduation nach GF**.

#### `internal/hexagon/application` — GF mit benannter Rückwärts-Lücke

GF ist belegt: Die Use-Case-Schnittstellen standen als Ports in
[`spec/architecture.md`](../../../../spec/architecture.md) §2.2/§2.4, bevor die
Services entstanden. **Eine** Konvention ist dennoch code-seitig entstanden:
„alle Ports nil-tolerant via package-private `noop*`-Defaults" steht nur im
Paket-README. Kein Modus-Wechsel — aber die Aussage gehört in die Sicht-Spec
und wird im Architektur-Slice dieser Welle mitgenommen.

### 9.3 Ergebnis

| Sub-Area | Modus | Änderung ggü. Erst-Pass |
|---|---|---|
| `internal/hexagon/domain` | GF | keine (ausdifferenziert aus „`hexagon/`") |
| `internal/hexagon/application` | GF | keine; Rückwärts-Lücke benannt (§9.2) |
| `internal/hexagon/port` | GF | keine |
| `internal/adapter/driving/cli` | **Hybrid** | **neu** — Erst-Pass sagte GF |
| `internal/adapter/driven` | GF | keine |
| `internal/e2e` | GF | **neu qualifiziert** |
| `tools/` + `scripts/` | GF | **neu qualifiziert** |
| `internal/**/README.md` | **BF** | **neu qualifiziert**, im Erst-Pass unsichtbar |
| `cmd/uboot` | — | **Rückstufung auf Aspirantin** |

Zwei Folge-Wirkungen außerhalb der Modus-Tabelle: ein Carveout-Inventar-Eintrag
für den `scan.ignore`-Glob `internal/**` (mit Plan-Anker) und ein neuer
`open/`-Slice für den Retrofit.

> **Nachtrag zur Graduation (2026-07-25, noch innerhalb der Welle).** Die
> Hybrid-Aussage zu `internal/adapter/driving/cli` hat ihre Graduation-Bedingung
> im selben Durchlauf erfüllt: Sentinel-Schichtung und Dual-Classifier-Regel
> stehen seither in der Sicht-Spec. Die Modus-Tabelle in
> [`harness/conventions.md`](../../../../harness/conventions.md) führt die
> Sub-Area deshalb als **GF (graduiert)**; die Hybrid-Diagnose oben bleibt als
> Audit-Befund stehen — sie war zum Prüfzeitpunkt korrekt und ist der Grund,
> warum die Regeln überhaupt in die Spec gehoben wurden. Aufgedeckt hat den
> nicht nachgezogenen Ledger der Review-Lauf zu dieser Welle (F-1,
> [`docs/reviews/`](../../../reviews/README.md)).
