# Slice Harness: Lastenheft-Historie als CR-Fußabdruck nachziehen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in
[`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Befund aus
[`slice-harness-baseline-bump-review-v3.5.2`](../done/slice-harness-baseline-bump-review-v3.5.2.md).

**Bezug:** [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
(die am 2026-07-24 geänderte Anforderung, deren Fußabdruck fehlt) und die
Doku-Mindeststruktur
[`LH-FA-PROJDOCS-001`](../../../../spec/lastenheft.md#lh-fa-projdocs-001--mindeststruktur).
Der Slice **ändert keine `LH-*`-Anforderung** — er trägt Buchführung über eine
bereits vereinbarte Änderung nach.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

`spec/lastenheft.md` bekommt den Fußabdruck, den die Baseline `v3.5.2` für
angenommene Change Requests verlangt: einen **Versions-Stand**, der sich bei
Vertragsänderungen erhöht, und einen Abschnitt **Historie**, der jede solche
Änderung mit Datum und Verweis auf den externen Vorgang festhält.

Heute steht das Lastenheft unverändert auf Version `0.1.0` / Status `Entwurf` /
Datum `2026-05-21` und hat keinen Historie-Abschnitt — obwohl
`LH-FA-PROJDOCS-002` am 2026-07-24 auf die MADR-ADR-Form umgestellt wurde. Die
Änderung ist also passiert, aber am Vertrag selbst nicht ablesbar; sie ist nur
über einen `done/`-Slice rekonstruierbar.

## 2. Definition of Done

- [ ] **Abschnitt `## Historie`** am Ende von
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) angelegt, Spalten
  Version / Datum / Änderung / Verweis (Form aus der vendorten
  `lastenheft.template.md`).
- [ ] **Bestandszeile für die MADR-Umstellung** eingetragen: Datum 2026-07-24,
  geänderte Anforderung `LH-FA-PROJDOCS-002`, Verweis auf den ausführenden
  `done/`-Slice als *Vehikel* — mit dem ausdrücklichen Vermerk, dass die
  Entscheidung außerhalb des Repos fiel (der Slice ist nicht die Autorität).
- [ ] **Versions-Feld entschieden und gesetzt.** Zu klären: Bleibt die
  Kopf-Tabelle die Versionsquelle oder wird auf die Template-Form
  (`**Version:**`-Zeile) umgestellt? Und welcher Sprung ist richtig —
  `0.1.0` → `0.2.0` (Minor, weil eine bestehende Anforderung inhaltlich
  geändert wurde)? Die Entscheidung wird im Slice begründet, nicht still
  gesetzt.
- [ ] **Status-Feld geprüft:** Das Lastenheft trägt `Entwurf`, obwohl seine IDs
  seit Monaten als verbindlich behandelt werden (Gates, Traceability-Matrix,
  Exit-Code-Vertrag). Entweder Status nachziehen oder die Abweichung als
  bewusste Aussage begründen.
- [ ] **Rückwirkende Vollständigkeit geprüft:** Gab es vor dem 2026-07-24
  weitere Änderungen an bereits angenommenen `LH-*`? Falls ja, gehören sie
  ebenfalls in die Historie; falls nein, wird das als Negativbefund
  festgehalten.
- [ ] `make docs-check` grün (Anker der neuen Überschrift, keine
  Referenzmodell-Verletzung — das Lastenheft verweist weiterhin **nicht**
  abwärts auf Slices; der Verweis in der Historie-Zeile ist deshalb sorgfältig
  zu formulieren bzw. als externer Vorgang zu benennen).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) Kopf | update | Version/Status als belastbare Felder |
| [`spec/lastenheft.md`](../../../../spec/lastenheft.md) §Historie | neu | CR-Fußabdruck laut Baseline `v3.5.2` |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-008` | update | Nachtrag auflösen, sobald der Fußabdruck existiert |

## 4. Trigger

Kein externer Trigger. Sinnvoll **vor** der nächsten Vertragsänderung — danach
müsste man zwei Einträge rückwirkend rekonstruieren statt einen.

## 5. Closure-Trigger

Historie-Abschnitt vorhanden und für den Bestand vollständig, Version/Status
begründet gesetzt, `make docs-check` grün, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Referenzmodell-Falle:** Das Lastenheft ist Rang 1 und darf nicht abwärts auf
  Slices verweisen. Die Historie-Spalte „Verweis" zeigt laut Baseline auf den
  **externen** Vorgang (Ticket, Vertragsanhang). u-boot hat keinen externen
  Ticket-Raum — die Form dieses Verweises ist die eigentliche offene Frage des
  Slice und muss die `matrix`-Regel überleben.
- **Rückwirkende Rekonstruktion:** Je länger gewartet wird, desto mehr Historie
  ist aus Commits zu rekonstruieren. Heute ist es genau ein Eintrag.
- **Versions-Semantik ohne Vorbild:** Es gibt bisher keine Regel, wann Major,
  Minor oder Patch steigt. Der Slice setzt sie implizit — deshalb gehört die
  Begründung in die Closure, nicht nur die Zahl.
- **Kein Carveout erwartet:** reine Buchführung, keine Anforderungsänderung.

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Vertrags-Stratum* — **GF** (Doku-führt) nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Der Slice schreibt keine neue Anforderung, sondern die
Änderungs-Buchführung über eine bestehende.
