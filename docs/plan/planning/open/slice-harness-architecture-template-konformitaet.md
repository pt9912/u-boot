# Slice Harness: `spec/architecture.md` auf Template-Konformität heben

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-harness-konformitaet-nachlauf` (Harness-Konformitäts-Nachlauf,
s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** Der in
[`slice-harness-regelwerk-adoption-v3.5.1`](../done/slice-harness-regelwerk-adoption-v3.5.1.md)
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

- [ ] **Kopf-Form angeglichen:** `**Status:** Aktiv.` + `**Letzte Änderung:**`
  statt der Entwurf-/Datum-Tabellenzeilen; die Hard-Rule-Aussage „keine Wellen,
  Slices, Commit-Hashes oder Closure-Daten in dieser Datei" steht sichtbar im
  Kopf. Der Bezug-Eintrag auf die `LH-FA-ARCH-*`-Kette bleibt erhalten.
- [ ] **Neuer Abschnitt Externe Abhängigkeiten** (Template §3): Tabelle mit
  Rolle und Substituierbarkeit je externem System/Library (Container-Runtime
  und Compose-Schnittstelle, CLI-Framework, YAML-Serialisierung). Die
  *Wahl-Begründung* bleibt draußen — die trägt die jeweilige
  Architekturentscheidung, nicht die Sicht.
- [ ] **Neuer Abschnitt Sequenz-Diagramme** (Template §4) mit mindestens drei
  kritischen Use-Cases entlang der Lastenheft-IDs — Projekt-Initialisierung,
  Add-on-Hinzufügen und Umgebung-Starten. Lanes sind die Architektur-Schichten
  (Driving-Adapter → Driving-Port → Application → Driven-Port →
  Driven-Adapter), nicht Funktionsnamen; jede Sequenz nennt ihre `LH-*`-ID.
- [ ] **Neuer Abschnitt Fehlermodelle und Resilienz** (Template §5): Tabelle
  Fehlerquelle → behandelnde Schicht → Logging/Sichtbarkeit, plus die zwei
  tragenden Regeln der Bestandsimplementierung, die bisher nur im Code stehen:
  (a) **Sentinel-Schichtung** — Driven-Sentinels werden vor
  Driving-/Application-Sentinels klassifiziert, weil die Schicht-Hierarchie die
  Klassifikations-Reihenfolge vorgibt; (b) **Dual-Classifier-Regel** — ein neuer
  oder aufgeteilter Driving-Sentinel muss sowohl in der Exit-Code-Klassifikation
  als auch in der Envelope-/Diagnostic-Abbildung des Driving-Adapters eingetragen
  werden, sonst driften Exit-Code und maschinenlesbare Ausgabe auseinander.
- [ ] **Sprach- und meilensteinfrei:** keine Phase-/Meilenstein-Tags, keine
  Slice-Referenzen, keine LOC-/Commit-Angaben, keine ADR-Verweise (die
  Änderungskopplung deklariert die ADR aufwärts, nicht die Sicht abwärts).
- [ ] **Referenzmodell gewahrt:** `make docs-check` `matrix` grün — die
  view-spec-Klasse verweist weiterhin nur auf das Vertrags-Stratum.
- [ ] **Bestand unverletzt:** Die neuen Abschnitte beschreiben den
  *implementierten* Stand, nicht ein Wunschbild. Jede Aussage ist am Code
  belegbar; abweichende Funde werden zu Findings, nicht stillschweigend
  „weggeschrieben".
- [ ] `make docs-check` grün; `make gates` als Closure-Sensor nicht nötig
  (kein Code-Delta) — Begründung in der Evidence.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Architektur-Sicht* — **GF** (Doku-führt; die
hexagonalen Regeln standen vor dem Code). Besonderheit: Die neuen Abschnitte
sind **nachträglich** aus dem Bestand geschrieben, also lokal
Brownfield-artig — deshalb die DoD-Auflage „am Code belegbar" und die
Finding-Pflicht bei Abweichung. Der Modus der Sub-Area bleibt GF, weil die
Änderungsrichtung Doku → Code erhalten bleibt; die Rückwärts-Belegung ist eine
einmalige Lücken-Schließung, keine Inventur-Linie.
