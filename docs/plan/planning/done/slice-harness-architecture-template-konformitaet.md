# Slice Harness: `spec/architecture.md` auf Template-Konformität heben

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** Der in
[`slice-harness-regelwerk-adoption-v3.5.1`](slice-harness-regelwerk-adoption-v3.5.1.md)
als „struktureller Umbau, auf Ansage" zurückgestellte Teil der
Template-Konformität. Betroffene Sicht-Spec:
[`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement)
(lesend, unverändert) sowie der Fehler-/Exit-Code-Vertrag
[`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes)
(lesend — die Architektur *beschreibt*, wo er durchgesetzt wird, sie *setzt*
ihn nicht neu). **Kein Change Request:** Die Struktur von `spec/architecture.md`
ist im Vertrags-Stratum nicht festgeschrieben; ergänzt werden Sichten, keine
Anforderungen.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

`spec/architecture.md` trägt heute Überblick, Schichten, Import-Regeln,
Enforcement, Tests, Anti-Patterns und Evolution — aber **keine
Sequenz-Sicht** und **kein Fehlermodell**. Beides fordert die vendorte
Referenz-Form (`.harness/baseline/v3.5.1/templates/spec/architecture.template.md`
§4/§5), und beides fehlt inhaltlich: Der Aufruf-Fluss durch die Schichten und
die Frage „wo wird ein Fehler abgefangen, wo propagiert, wo klassifiziert"
sind heute nur verstreut im Code-Kommentar dokumentiert.

Dieser Slice schließt die Lücke und gleicht zusätzlich die **Kopf-Form** an
(Status-/Änderungsdatum-Zeile plus die Hard-Rule-Aussage, dass die Datei keine
zeitliche Schicht trägt).

## 2. Definition of Done

- [x] **Kopf-Form angeglichen:** `**Status:** Aktiv.` + `**Letzte Änderung:**`
  statt der Entwurf-/Datum-Tabellenzeilen; die Hard-Rule-Aussage „keine Wellen,
  Slices, Commit-Hashes oder Closure-Daten in dieser Datei" steht sichtbar im
  Kopf. Der Bezug-Eintrag auf die `LH-FA-ARCH-*`-Kette bleibt erhalten.
- [x] **Neuer Abschnitt Externe Abhängigkeiten** (Template §3): Tabelle mit
  Rolle und Substituierbarkeit je externem System/Library (Container-Runtime
  und Compose-Schnittstelle, CLI-Framework, YAML-Serialisierung). Die
  *Wahl-Begründung* bleibt draußen — die trägt die jeweilige
  Architekturentscheidung, nicht die Sicht.
- [x] **Neuer Abschnitt Sequenz-Diagramme** (Template §4) mit mindestens drei
  kritischen Use-Cases entlang der Lastenheft-IDs — Projekt-Initialisierung,
  Add-on-Hinzufügen und Umgebung-Starten. Lanes sind die Architektur-Schichten
  (Driving-Adapter → Driving-Port → Application → Driven-Port →
  Driven-Adapter), nicht Funktionsnamen; jede Sequenz nennt ihre `LH-*`-ID.
- [x] **Neuer Abschnitt Fehlermodelle und Resilienz** (Template §5): Tabelle
  Fehlerquelle → behandelnde Schicht → Logging/Sichtbarkeit, plus die zwei
  tragenden Regeln der Bestandsimplementierung, die bisher nur im Code stehen:
  (a) **Sentinel-Schichtung** — Driven-Sentinels werden vor
  Driving-/Application-Sentinels klassifiziert, weil die Schicht-Hierarchie die
  Klassifikations-Reihenfolge vorgibt; (b) **Dual-Classifier-Regel** — ein neuer
  oder aufgeteilter Driving-Sentinel muss sowohl in der Exit-Code-Klassifikation
  als auch in der Envelope-/Diagnostic-Abbildung des Driving-Adapters eingetragen
  werden, sonst driften Exit-Code und maschinenlesbare Ausgabe auseinander.
- [x] **Sprach- und meilensteinfrei:** keine Phase-/Meilenstein-Tags, keine
  Slice-Referenzen, keine LOC-/Commit-Angaben, keine ADR-Verweise (die
  Änderungskopplung deklariert die ADR aufwärts, nicht die Sicht abwärts).
- [x] **Referenzmodell gewahrt:** `make docs-check` `matrix` grün — die
  view-spec-Klasse verweist weiterhin nur auf das Vertrags-Stratum.
- [x] **Bestand unverletzt:** Die neuen Abschnitte beschreiben den
  *implementierten* Stand, nicht ein Wunschbild. Jede Aussage ist am Code
  belegbar; abweichende Funde werden zu Findings, nicht stillschweigend
  „weggeschrieben".
- [x] `make docs-check` grün; `make gates` als Closure-Sensor nicht nötig
  (kein Code-Delta) — Begründung in der Evidence.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/architecture.md`](../../../../spec/architecture.md) Kopf | update | Template-Kopf-Form + Hard-Rule-Aussage |
| [`spec/architecture.md`](../../../../spec/architecture.md) §Externe Abhängigkeiten | neu | Template §3 fehlt bisher vollständig |
| [`spec/architecture.md`](../../../../spec/architecture.md) §Sequenz-Diagramme | neu | Template §4; Aufruf-Fluss bisher nur als ASCII-Überblick |
| [`spec/architecture.md`](../../../../spec/architecture.md) §Fehlermodelle und Resilienz | neu | Template §5; Sentinel-Schichtung und Dual-Classifier-Regel bisher nur Code-Kommentar |

Abschnitts-Nummerierung: Die neuen Abschnitte werden in die bestehende
Nummerierung eingefügt; nachfolgende Abschnitte verschieben sich. Interne
Anker-Verweise sind dabei mitzuziehen (`docs-check` `anchors` ist der Sensor).

## 4. Trigger

Direkt fortsetzbar; unabhängig von FS-1..FS-4, deshalb zuletzt in der Welle
(größter Diff, geringste Kopplung an die anderen Slices).

## 5. Closure-Trigger

Alle vier Abschnitte vorhanden und am Code belegt, `make docs-check` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Sicht driftet gegen Code:** Ein Sequenz-Diagramm ist schnell geschrieben und
  altert leise. Gegenmaßnahme: Die Sequenzen bleiben auf Schicht-Granularität
  (Ports/Adapter), nicht auf Funktions-Granularität — diese Ebene ist durch
  `depguard` maschinell abgesichert und driftet nicht ohne Gate-Ausschlag.
- **Anker-Drift durch Umnummerierung:** Neue Abschnitte verschieben bestehende
  Nummern; interne Verweise („siehe §4") und Anker aus anderen Dokumenten
  können brechen. `docs-check` `anchors` fängt Link-Anker; reine Prosa-Verweise
  („§4") sind manuell zu prüfen.
- **Vertrags-Kollision:** Falls sich beim Schreiben zeigt, dass eine
  Fehlerklasse im Code anders klassifiziert ist als im Lastenheft beschrieben,
  ist das ein **Finding mit CR-Bedarf**, kein stiller Angleich der Sicht.
- **Kein Carveout erwartet:** reine Doku-Ergänzung, kein Code-Delta.

## 7. Closure-Notiz (nach `done/`)

### Verification Evidence

Scope:
- Slice: `slice-harness-architecture-template-konformitaet`
- IDs: **keine** Anforderung geändert. Lesend belegt:
  [`LH-FA-ARCH-001`](../../../../spec/lastenheft.md#lh-fa-arch-001--hexagonales-pattern)..[`LH-FA-ARCH-003`](../../../../spec/lastenheft.md#lh-fa-arch-003--import-regeln-und-enforcement),
  [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes),
  [`LH-FA-CLI-005A`](../../../../spec/lastenheft.md#lh-fa-cli-005a--interaktivität-und-automatisierung),
  [`LH-NFA-USE-004`](../../../../spec/lastenheft.md#lh-nfa-use-004--maschinenlesbare-ausgabe),
  [`LH-AK-006`](../../../../spec/lastenheft.md#lh-ak-006--idempotenz).
- Artefakte: [`spec/architecture.md`](../../../../spec/architecture.md) (Kopf,
  §3 Externe Abhängigkeiten, §6 Sequenz-Diagramme, §7 Fehlermodelle und
  Resilienz, Umnummerierung),
  [`slice-harness-architecture-bestandsabgleich`](slice-harness-architecture-bestandsabgleich.md)
  (Folge-Plan), [`roadmap.md`](../in-progress/roadmap.md).

DoD-Abgleich: alle Punkte erfüllt — der Punkt „Bestand unverletzt: jede Aussage
ist am Code belegbar" allerdings **erst nach der Finding-Abarbeitung**. Der
unabhängige Review hat genau an diesem Punkt fünf MEDIUM gefunden (F-2..F-6);
sie sind in `4fed84f` behoben. Die DoD-Auflage hat damit gehalten, was sie
sollte — sie hat den eigenen Slice geprüft, nicht nur begleitet.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` | pass | 129 Dateien / 0 Befunde nach dem Umbau, 131 / 0 nach Finding-Abarbeitung; `anchors` bestätigt die mitgezogenen internen Verweise, `matrix` die unveränderte Aufwärts-Richtung der view-spec-Klasse |
| Code-Gegenprobe (lesend) | pass mit Findings | Exit-Klassen-Tabelle, Sentinel-Schichtung und Resilienz-Muster am Bestand belegt; drei Aussagen wichen ab und wurden korrigiert (F-3, F-4, F-5) |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| Template `spec/architecture.template.md` §3/§4/§5 | die drei neuen Abschnitte; Kopf-Form + Hard-Rule-Aussage übernommen |
| [`LH-FA-CLI-006`](../../../../spec/lastenheft.md#lh-fa-cli-006--exit-codes) | §7.1-Tabelle Zeile für Zeile gegen die Exit-Code-Klassifikation des Adapters geprüft (Review-Negativbefund) |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) | `matrix` grün; die Datei enthält keine Verweise auf ADR, Slice, Carveout oder Roadmap |
| Graduation `internal/adapter/driving/cli` | §7.2 trägt Sentinel-Schichtung und Dual-Classifier-Regel — die Bedingung aus dem FS-3-Audit |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make lint` / `make test` — kein Go-Delta (`git diff --stat`
  ohne `*.go`); `make docs-check` ist der Sensor der Zeile „Nur Markdown/Doku".
  Der Code wurde ausschließlich lesend als Beleg herangezogen.

Independent Review (Frischkontext, Rollentrennung): kein HIGH, aber der
Schwerpunkt der Welle-Findings liegt hier — fünf MEDIUM plus zwei LOW, alle
abgearbeitet in `4fed84f`:
- **F-2** Dual-Classifier-Regel als „zwei Stellen" formuliert, tatsächlich eine
  zentrale Klassifikation plus je Subkommando eine Abbildung (neun).
- **F-3** Init-Sequenz kehrte die im Code begründete Reihenfolge um
  (Soft-Existing-Bestätigung läuft **vor** dem Datei-Plan).
- **F-4** Up-Sequenz zeigte eine Runtime-Validierung im Use-Case; tatsächlich
  prüft er über Datei-Ports, die Runtime-Probe liegt adapter-intern.
- **F-5** Exit-Klassen `13`/`15` als „reserviert" bezeichnet — der Vertrag
  erlaubt sie bei dokumentierter Bedeutung.
- **F-6** §6 berief sich auf `depguard` als Schutz gegen Sequenz-Drift; das Gate
  prüft Import-Richtungen, nicht Schritt-Reihenfolgen.
- **F-8/F-9** (LOW) veraltetes Datum in §10; Bestätigungs-Gate ohne die
  `2`-Verzweigung des Vertrags.
Ein INFO (F-13) betrifft die unvollständige Sentinel-Liste in §2.3 — bereits als
Folge-Slice eingeplant.

Commit / Artefakt: `c35249d` (Kopf, §3/§6/§7, Umnummerierung); `4fed84f`
(Finding-Abarbeitung); Lifecycle-Move `open/` → `done/` im Folge-Commit.

### Steering-Loop-Lerneintrag

- **Die eigene Risiko-Gegenmaßnahme war der Fehler.** Dieser Plan hat in §6
  notiert, Schicht-Granularität sei „durch `depguard` maschinell abgesichert" —
  und genau dieser Satz wanderte als Behauptung in die Spec (F-6), während zwei
  Sequenzen daneben lagen (F-3, F-4). Lehre: Eine Gegenmaßnahme, die auf ein
  Gate verweist, muss geprüft werden, *was das Gate tatsächlich misst*.
  `depguard` misst Import-Richtungen, nicht Abläufe.
- **Rückwärts geschriebene Doku braucht einen zweiten Leser.** Alle drei
  Sequenz-/Modell-Fehler waren beim Schreiben plausibel und beim Nachlesen im
  Code falsch. Die DoD-Auflage „am Code belegbar" hat sie nicht verhindert —
  gefunden hat sie erst der Frischkontext-Review. Für den nächsten
  Bestandsabgleich heißt das: Beleg-Stichproben gehören *in* den Arbeitsgang,
  nicht in die Abnahme.
- **Ein Widerspruch zwischen Rang 1 und Rang 2 entsteht durch Zusammenfassen.**
  F-5 und F-9 sind beide daraus entstanden, dass eine Vertragspassage in eine
  Tabellenzelle verdichtet wurde. Wo der Vertrag verzweigt, verzweigt die Sicht
  auch — oder sie verweist, statt zu verdichten.
- **Folge-Slices:**
  [`slice-harness-architecture-bestandsabgleich`](slice-harness-architecture-bestandsabgleich.md)
  — §2 beschreibt weitgehend den Frühstand (elf Driving-Ports existieren, drei
  sind dokumentiert).

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Architektur-Sicht* — **GF** (Doku-führt; die
hexagonalen Regeln standen vor dem Code). Besonderheit: Die neuen Abschnitte
sind **nachträglich** aus dem Bestand geschrieben, also lokal
Brownfield-artig — deshalb die DoD-Auflage „am Code belegbar" und die
Finding-Pflicht bei Abweichung. Der Modus der Sub-Area bleibt GF, weil die
Änderungsrichtung Doku → Code erhalten bleibt; die Rückwärts-Belegung ist eine
einmalige Lücken-Schließung, keine Inventur-Linie.
