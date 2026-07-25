# Slice Gate: Traceability-Matrix (`--trace`) für u-boot konfigurieren

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** noch keiner Welle zugeordnet (Wartungs-Kandidat in
[`roadmap.md`](../in-progress/roadmap.md) §Nächste Wellen) — Befund aus dem
Gate-Ausbau, bewusst **nach** `welle-gate-ausbau-v0.51` eingeplant.

**Bezug:** [`LH-FA-PROJDOCS-006`](../../../../spec/lastenheft.md#lh-fa-projdocs-006--dokumentationsreferenzmodell)
(Referenzmodell) und die handgepflegte Traceability-Matrix in
`spec/lastenheft.md` §13. **Achtung:** §13 liegt im **Vertrags-Stratum** — jede
Änderung daran ist ein Change Request mit Version-Bump und Historie-Zeile
([`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)-Muster).

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Die Frage „welche Anforderung ist durch kein Artefakt belegt?" maschinell
beantwortbar machen. Heute ist sie nur von Hand zu beantworten — beim
Durchsehen der offenen Punkte musste genau das gemacht werden.

`--trace` erzeugt die Matrix bereits; sie ist nur nie konfiguriert worden.
**Vorab gemessen (2026-07-25):** Zwei Zeilen Konfiguration genügen, um aus
79 „Waisen" **12** zu machen:

```yaml
trace:
  slices:
    dir: docs/plan/planning
    file-pattern: '^(slice-.+)\.md$'
```

Der Default erwartet Slice-Dateinamen der Form `slice-<NNN>-…`; u-boot nutzt
`slice-<phase>-<slug>`. Ein reiner Formatunterschied — die ADR-Spalte
funktionierte von Anfang an.

**Die verbleibenden 12 sind Traceability-Lücken, keine
Implementierungslücken:** [`LH-FA-CLI-001`](../../../../spec/lastenheft.md#lh-fa-cli-001--cli-aufruf) (CLI-Aufruf), `-002` (Hilfeausgabe),
`-003` (Versionsausgabe) sind demonstrierbar geliefert; sie stammen aus frühen
Skeleton-Slices, die ihre Kennungen nie benannt haben. Dasselbe Muster bei
[`LH-FA-DOC-001`](../../../../spec/lastenheft.md#lh-fa-doc-001--compose-datei-erzeugen)/`-003`/`-004`, [`LH-QA-001`](../../../../spec/lastenheft.md#lh-qa-001--automatisierte-tests)/`-002`,
[`LH-FA-PROJDOCS-004`](../../../../spec/lastenheft.md#lh-fa-projdocs-004--archivierung),
[`LH-FA-BUILD-003`](../../../../spec/lastenheft.md#lh-fa-build-003--build-args-und-pin-politik),
[`LH-FA-DEV-002`](../../../../spec/lastenheft.md#lh-fa-dev-002--vs-code-kompatibilität).

## 2. Definition of Done

- [ ] **Vollständiger Kennungs-Umfang entschieden und konfiguriert.** Die RTM
  sieht heute nur `LH-FA-*` und `LH-QA-*` — **79 von 139** Kennungen. Draußen
  bleiben `NFA`, `SA`, `DA`, `AK`, `ZB`, `PE`, `PÜ`, `LESE`, `ABG`, `RISK`,
  `MVP`, `OPEN`. Zu entscheiden ist, welche Familien überhaupt
  belegpflichtig sind: Eine Lesehinweis-Kennung ([`LH-LESE-001`](../../../../spec/lastenheft.md#lh-lese-001--modalverben) Modalverben)
  oder eine Abgrenzung (`LH-ABG-*`) wird nie einen Slice haben — sie als Waise
  zu führen wäre Rauschen. Die Entscheidung wird begründet, nicht durch ein
  weites Muster umgangen.
- [ ] **Umlaut-Fall gelöst:** `LH-PÜ-001`/`-002` matchen weder das neue
  `trace.requirements.id-pattern` noch das bestehende `ids`-Muster
  (`LH(?:-[A-Z0-9]+)+-\d{3}`) — beide kennen kein `Ü`. Der Nebenbefund gehört
  mitbehoben, sonst bleiben zwei Kennungen dauerhaft unsichtbar.
- [ ] **Advisory vor Gate.** `--require-complete` bleibt **aus**. Ein Gate, das
  ab Tag eins rot ist, wird abgeschaltet statt befolgt. `make doc-trace` ist
  das informative Target (das Fragment liefert es bereits); über ein Gate wird
  erst entschieden, wenn die Lückenliste abgearbeitet oder bewusst akzeptiert
  ist.
- [ ] **Die 12 Lücken einzeln bewertet** — je Kennung eine von drei Antworten:
  (a) ein `done/`-Slice benennt sie nachträglich (Querverweis-Korrektur, nach
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  zulässig), (b) sie ist tatsächlich unbelegt und bekommt einen Plan, oder
  (c) sie ist strukturell belegfrei und wird als solche deklariert.
  Keine Sammelantwort.
- [ ] **§13-Kollision entschieden, nicht übergangen.** Es gäbe dann zwei
  Matrizen: die handgepflegte im Vertrag und die generierte. Optionen:
  §13 bleibt (Prioritäts-Index, die generierte ist die Belegsicht) — oder §13
  wird per CR auf die generierte Sicht reduziert. **Empfehlung im Plan:**
  §13 unangetastet lassen; die Spalten „Pflichtenheft"/„Testfall" stehen dort
  ohnehin auf `-`, es ist faktisch ein Index mit Prioritäten, kein
  Coverage-Artefakt. Ein CR nur, wenn sich das als falsch erweist.
- [ ] **`MR-005` nachgezogen** (Trace-Konfiguration als Teil der Gate-Haltung).
- [ ] `make docs-check` grün; `make doc-trace` liefert eine Matrix, deren
  Waisenzahl im Slice dokumentiert ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) `trace` | neu | `slices`-Block + `requirements.id-pattern` |
| [`.d-check.yml`](../../../../.d-check.yml) `ids` | update | Umlaut-Fall im Kennungs-Muster |
| `docs/plan/planning/done/*.md` | update falls (a) | nachträgliche Kennungs-Nennung, Querverweis-Korrektur |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | Trace-Konfiguration dokumentieren |

## 4. Trigger

Gefeuert: Befund beim Durchsehen der offenen Lastenheft-Punkte (2026-07-25).
Bewusst **nach** `welle-gate-ausbau-v0.51` und nach dem anstehenden
Release-Cut eingeplant — der Sicherheits-Fix im Runtime-Image drängt stärker.

## 5. Closure-Trigger

Umfang entschieden, 12 Lücken einzeln bewertet, §13-Frage beantwortet,
`make doc-trace` mit dokumentierter Waisenzahl, `make docs-check` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Ein weites Muster erzeugt Rauschen statt Signal.** Alle 139 Kennungen in
  die RTM zu nehmen produziert Dutzende „Waisen", die keine sind
  (Lesehinweise, Abgrenzungen, Risiken). Die Familien-Auswahl ist die
  eigentliche inhaltliche Arbeit dieses Slice, nicht die Konfiguration.
- **Nachträgliche Kennungs-Nennung in `done/`-Slices ist eine Grenzfrage.**
  Einen fehlenden Link zu setzen ist klar eine Querverweis-Korrektur. Eine
  Kennung *neu einzufügen*, die dort nie stand, ist mehr als das — es
  behauptet rückwirkend einen Bezug. Je Fall zu prüfen, im Zweifel Variante
  (c) statt (a).
- **Zwei Matrizen sind eine Drift-Quelle.** Genau das Muster, das `targets`
  und `planning` bei anderen Doppelpflegen aufgedeckt haben. Deshalb ist die
  §13-Frage ein DoD-Punkt und keine Fußnote.
- **Kein Carveout erwartet.**

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *spec / docs* — beide **GF**
nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration. Sollte die §13-Frage zu einem CR führen, ist zusätzlich das
Vertrags-Stratum berührt — dann greift der Fußabdruck-Weg aus
[`slice-harness-lastenheft-historie-cr-fussabdruck`](../done/slice-harness-lastenheft-historie-cr-fussabdruck.md).
