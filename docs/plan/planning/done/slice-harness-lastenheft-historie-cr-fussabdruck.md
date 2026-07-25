# Slice Harness: Lastenheft-Historie als CR-Fußabdruck nachziehen

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in
[`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Befund aus
[`slice-harness-baseline-bump-review-v3.5.2`](slice-harness-baseline-bump-review-v3.5.2.md).

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

- [x] **Abschnitt `## Historie`** am Ende von
  [`spec/lastenheft.md`](../../../../spec/lastenheft.md) angelegt, Spalten
  Version / Datum / Änderung / Verweis (Form aus der vendorten
  `lastenheft.template.md`).
- [x] **Bestandszeile für die MADR-Umstellung** eingetragen: Datum 2026-07-24,
  geänderte Anforderung `LH-FA-PROJDOCS-002`, Verweis auf den ausführenden
  `done/`-Slice als *Vehikel* — mit dem ausdrücklichen Vermerk, dass die
  Entscheidung außerhalb des Repos fiel (der Slice ist nicht die Autorität).
- [x] **Versions-Feld entschieden und gesetzt.** Zu klären: Bleibt die
  Kopf-Tabelle die Versionsquelle oder wird auf die Template-Form
  (`**Version:**`-Zeile) umgestellt? Und welcher Sprung ist richtig —
  `0.1.0` → `0.2.0` (Minor, weil eine bestehende Anforderung inhaltlich
  geändert wurde)? Die Entscheidung wird im Slice begründet, nicht still
  gesetzt.
- [x] **Status-Feld geprüft:** Das Lastenheft trägt `Entwurf`, obwohl seine IDs
  seit Monaten als verbindlich behandelt werden (Gates, Traceability-Matrix,
  Exit-Code-Vertrag). Entweder Status nachziehen oder die Abweichung als
  bewusste Aussage begründen.
- [x] **Rückwirkende Vollständigkeit geprüft:** Gab es vor dem 2026-07-24
  weitere Änderungen an bereits angenommenen `LH-*`? Falls ja, gehören sie
  ebenfalls in die Historie; falls nein, wird das als Negativbefund
  festgehalten.
- [x] `make docs-check` grün (Anker der neuen Überschrift, keine
  Referenzmodell-Verletzung — das Lastenheft verweist weiterhin **nicht**
  abwärts auf Slices; der Verweis in der Historie-Zeile ist deshalb sorgfältig
  zu formulieren bzw. als externer Vorgang zu benennen).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

### Die Wurzel lag eine Ebene tiefer

Der Slice trat an, eine fehlende `## Historie` nachzutragen. Beim ersten Blick
auf den Kopf zeigte sich, dass die Historie nur das *Symptom* war: Das
Lastenheft trug seit dem 2026-05-21 **Status `Entwurf`**, und im Modell der
Baseline steuert genau dieses Feld die Verbindlichkeit der IDs. Formal gab es
im Entwurfsstatus überhaupt keine Fußabdruck-Pflicht — praktisch behandelte
u-boot die IDs seit Monaten als bindend (Exit-Code-Vertrag, Coverage- und
Referenz-Gates, Traceability-Matrix), und die MADR-Umstellung lief ausdrücklich
als „Change Request am Vertrags-Stratum".

Die Deklaration widersprach also der gelebten Praxis, und das *unbemerkt*: Ein
`Entwurf`-Status macht keine Gate rot. Aufgelöst per Entscheidung des
Projektinhabers zugunsten der Praxis (Status → `Accepted`), nicht zugunsten der
Formalie.

### Verification Evidence

Scope:
- Slice: `slice-harness-lastenheft-historie-cr-fussabdruck`
- IDs: **keine Anforderung inhaltlich geändert.** Der Slice ändert die
  *Kopf-Metadaten* (Version, Status) und ergänzt einen Historie-Abschnitt —
  Buchführung über eine bereits vereinbarte Änderung, keine neue oder geänderte
  `LH-*`-Aussage.
  [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  wird in der Historie-Zeile nur *benannt*, nicht angefasst.
- Artefakte: [`spec/lastenheft.md`](../../../../spec/lastenheft.md) (Kopf +
  neuer §16 Historie), [`harness/conventions.md`](../../../../harness/conventions.md)
  (`MR-008`-Nachtrag aufgelöst).

DoD-Abgleich: alle Punkte erfüllt.

- **Versions-Sprung:** `0.1.0` → `0.2.0` (Minor). Begründung: Eine bestehende
  Anforderung wurde inhaltlich geändert, aber keine entfernt und keine
  Nummerierung gebrochen — Major bleibt dem Fall vorbehalten, dass Anforderungen
  entfallen oder IDs neu geschnitten werden.
- **Form des Kopfes:** Die bestehende u-boot-Tabelle bleibt, statt auf die
  `**Version:**`-Zeilenform des Templates zu wechseln. Die Template-**Substanz**
  (Version steuert Änderungs-Buchführung, Status steuert Verbindlichkeit) ist
  übernommen; die Zeilenform ist kosmetisch und hätte einen breiten Diff ohne
  Erkenntnisgewinn erzeugt — dieselbe Reconciliation-Logik wie bei der
  ADR-Titelform (`MR-008`).
- **Rückwirkende Vollständigkeit (Negativbefund):** 33 Commits haben
  `spec/lastenheft.md` je berührt; **32 davon** liegen vor dem Status-Wechsel
  und damit in der Entwurfsphase, in der die IDs keine Verbindlichkeit trugen.
  Sie bekommen bewusst **keine** Einzeleinträge — eine rekonstruierte
  Pseudo-Historie der Entwurfsphase wäre Genauigkeit vortäuschende Arbeit. Die
  Zeile `0.1.0` sagt das ausdrücklich.
- **Verweis-Form (die eigentliche offene Frage des Plans):** Die Spalte nennt
  den externen Vorgang in Prosa („Vereinbarung mit dem Projektinhaber,
  ausgelöst durch die Adoption des externen Betriebsregelwerks") und **keine
  Slice-Kennung**. Damit ist die Frage doppelt gelöst: inhaltlich, weil die
  Baseline genau den externen Vorgang meint, und mechanisch, weil ein
  `slice-*`-Token hier in eine Zwickmühle geführt hätte — die `ids`-Regel
  verlangt für jede Planning-Kennung einen Link, die `matrix`-Regel verbietet
  dem Vertrags-Stratum jeden Abwärts-Link. Kein Token, kein Konflikt.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make docs-check` | pass | 132 Dateien / 0 Befunde; `matrix` bestätigt, dass das Vertrags-Stratum weiterhin nicht abwärts verweist, `anchors` den neuen §16-Anker |
| Git-Historie als Inventar | pass | `git log -- spec/lastenheft.md` = 33 Commits; Abgrenzung Entwurfsphase/Vertragsphase am Status-Wechsel belegt |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| Baseline `v3.5.2` §Spec-Stratifizierung (CR-Fußabdruck) | §16 Historie + Versions-Feld + Status-Feld |
| [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell) | Verweis-Spalte ohne Abwärts-Kante; `matrix` grün |
| `MR-008` Nachtrag | aufgelöst, mit Datum und Verweis auf diesen Slice |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta. Anmerkung: Der Status-Wechsel auf
  `Accepted` ändert **keine** Gate-Semantik; die Gates prüften die IDs ohnehin
  schon als bindend. Genau das war der Widerspruch.

Independent Review: nicht durchgeführt. Der Diff ist klein und vollständig
deklarativ (drei Kopf-Zeilen, ein neuer Abschnitt); die eine
Ermessens-Entscheidung — Status und Versions-Sprung — ist vom Projektinhaber
getroffen, nicht vom Autor. Ein Frischkontext-Review würde hier den
Entscheidungsträger prüfen, nicht die Umsetzung.

Commit / Artefakt: `2cd9146` (Kopf-Felder, §16 Historie, `MR-008`-Nachtrag aufgeloest, Roadmap).

### Steering-Loop-Lerneintrag

- **Ein Statusfeld, das nichts durchsetzt, driftet unbemerkt.** `Entwurf` stand
  zwei Monate lang im Kopf, während jedes Gate die IDs als bindend behandelte.
  Kein Sensor konnte das melden, weil das Feld nirgends ausgewertet wird. Das
  ist der Normalfall bei Deklarations-Feldern — sie brauchen entweder einen
  Leser oder einen Anlass. Der Anlass war hier der Baseline-Bump.
- **Ein Fußabdruck-Befund führt auf die Frage, wer entscheiden darf.** Die
  Baseline-Regel („weder ADR noch Slice ändern `LH-*`") wirkt formalistisch, bis
  man den eigenen Bestand daran hält: Der `done/`-Slice, der die MADR-Umstellung
  ausgeführt hat, hätte man ohne diese Regel als Entscheidungs-Autorität lesen
  können. Der Fußabdruck macht sichtbar, dass die Entscheidung woanders fiel.
- **Nicht jede Lücke will vollständig gefüllt werden.** Die naheliegende
  Fleißarbeit wäre gewesen, aus 33 Commits eine Historie zu rekonstruieren. Sie
  hätte Präzision suggeriert, wo Entwurfsphase war. Eine Zeile, die die
  Abgrenzung *benennt*, ist ehrlicher als dreißig, die sie verwischen.
- **Folge-Slices:** keine.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *spec / Vertrags-Stratum* — **GF** (Doku-führt) nach
[`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Der Slice schreibt keine neue Anforderung, sondern die
Änderungs-Buchführung über eine bestehende.
