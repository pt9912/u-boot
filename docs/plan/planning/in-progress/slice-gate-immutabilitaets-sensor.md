# Slice Gate: Immutabilitäts-Sensor für Accepted-ADRs und `done/`-Slices

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-gate-ausbau-v0.51` (s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
(Grandfathering: Accepted-ADRs sind **unveränderlich**) und
[`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
(`done/`-Inhalte nur korrigierend änderbar). Beide Hard Rules stehen zusätzlich
in [`AGENTS.md`](../../../../AGENTS.md).

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

u-boot trägt **zwei Immutabilitäts-Hard-Rules** — Accepted-ADRs werden nicht
inhaltlich umgeschrieben, `done/`-Artefakte nur korrigierend geändert — und
für beide gibt es heute **keinen Sensor**. Sie gelten durch Aufmerksamkeit.
Das Modul `vcs` prüft genau das: den normalisierten **Core** einer geschützten
Datei-Klasse über eine Commit-Range.

Entscheidend für den Aufwand: `vcs` arbeitet **musterbasiert** (`paths`-Glob
plus `immutable-when`-Regex) — es braucht **keine** Pin-Marker in den rund 40
`done/`-Slices und 13 ADRs. Das war die ursprüngliche Sorge und ist gegenstandslos.

## 2. Definition of Done

- [ ] **Regex an u-boots Form angepasst und verifiziert:** Die
  Handbuch-Beispiele setzen `**Status:** Accepted` als Inline-Feld voraus;
  u-boots ADRs tragen `## Status` als Überschrift mit `Accepted` in der
  Folgezeile — **auch die MADR-konformen**, weil die Umstellung erst beim
  nächsten inhaltlichen Anfassen greift. `immutable-when`, `status-line` und
  `head-allow` sind entsprechend zu setzen und an einem echten ADR zu prüfen.
- [ ] **Zulässige Übergänge erlaubt:** `Accepted` → `Superseded by <NNNN>-<slug>`
  ist nach u-boot-Konvention legitim (nicht die Template-Form
  `Superseded by ADR-NNNN`) und darf nicht als Drift melden.
- [ ] **`exclude-sections` bestimmt:** Bei ADRs gehört die `## Geschichte`
  (bzw. ihr u-boot-Äquivalent) nicht zum Core — sie wächst bei jedem
  Status-Ereignis. Bei `done/`-Slices ist zu klären, ob Querverweis-Korrekturen
  (ausdrücklich erlaubt) den Core berühren; falls ja, braucht die Regel eine
  Antwort, sonst meldet jeder Lifecycle-Nachzug Drift.
- [ ] **Zweite Datei-Klasse `done/`:** Ob `vcs` beide Klassen in einem
  Konfigurations-Block trägt oder ob `done/` eine eigene Behandlung braucht
  (dort gibt es keine Status-Zeile, die Immutabilität folgt aus dem
  **Verzeichnis**), ist zu entscheiden und zu begründen.
- [ ] **Integrationspunkt gewählt:** `vcs` braucht eine Commit-Range
  (`--range`/`--staged`) und gehört damit **nicht** in `make docs-check`,
  sondern an CI (Range über den PR) oder einen lokalen pre-commit-Lauf. Die
  Wahl ist zu begründen und in [`AGENTS.md`](../../../../AGENTS.md)
  §Quality Gates einzutragen.
- [ ] **Gegenprobe am realen Verlauf:** Der Sensor läuft über eine Range, die
  eine bekannte legitime Änderung enthält (Lifecycle-Move, Hash-Nachtrag) und
  über eine konstruierte illegitime — er muss die eine durchlassen und die
  andere melden. Ein Sensor, der nie ausgeschlagen hat, ist nicht verifiziert.
- [ ] **`MR-005`/`MR-008` nachgezogen** (Gate-Haltung bzw. ADR-Politik).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) | update | `vcs`-Block (`paths`, `immutable-when`, `exclude-sections`, `head-allow`) |
| CI-Workflow bzw. `Makefile` | update | Integrationspunkt mit Range |
| [`AGENTS.md`](../../../../AGENTS.md) §Quality Gates | update | neues Target dokumentieren |
| [`harness/conventions.md`](../../../../harness/conventions.md) | update | `MR-005`, ggf. `MR-008` |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump. Letzter Slice
der Welle — er hat die meisten offenen Detailfragen.

## 5. Closure-Trigger

Sensor konfiguriert, an legitimer **und** illegitimer Änderung gegengeprüft,
Integrationspunkt dokumentiert, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **Falsch-positive lähmen mehr als sie schützen.** Ein Sensor, der bei jedem
  Lifecycle-Move oder Hash-Nachtrag anschlägt, wird abgeschaltet. Die
  `exclude-sections`- und Core-Definition entscheidet über den Nutzen —
  deshalb die Gegenprobe an einer realen Range in der DoD.
- **Grandfathering-Grenze:** Die lean-ADRs (`0001`–`0010`, `0013`) und die
  MADR-Form unterscheiden sich im Kopf. Ein Regex, der nur eine Form trifft,
  schützt die andere Hälfte nicht — und das fiele nicht auf, weil „keine
  Befunde" wie Erfolg aussieht.
- **Range-Abhängigkeit:** `vcs` ist fail-closed ohne `.git`. In einem
  Container-Lauf mit `:ro`-Mount muss `.git` mitgemountet sein; das ist beim
  Integrationspunkt zu prüfen.
- **Kein Carveout erwartet.**

## 7. Closure-Notiz (nach `done/`)

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
