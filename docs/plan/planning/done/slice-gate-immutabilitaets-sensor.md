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

- [x] **Regex an u-boots Form angepasst und verifiziert:** Die
  Handbuch-Beispiele setzen `**Status:** Accepted` als Inline-Feld voraus;
  u-boots ADRs tragen `## Status` als Überschrift mit `Accepted` in der
  Folgezeile — **auch die MADR-konformen**, weil die Umstellung erst beim
  nächsten inhaltlichen Anfassen greift. `immutable-when`, `status-line` und
  `head-allow` sind entsprechend zu setzen und an einem echten ADR zu prüfen.
- [x] **Zulässige Übergänge erlaubt:** `Accepted` → `Superseded by <NNNN>-<slug>`
  ist nach u-boot-Konvention legitim (nicht die Template-Form
  `Superseded by ADR-NNNN`) und darf nicht als Drift melden.
- [x] **`exclude-sections` bestimmt:** Bei ADRs gehört die `## Geschichte`
  (bzw. ihr u-boot-Äquivalent) nicht zum Core — sie wächst bei jedem
  Status-Ereignis. Bei `done/`-Slices ist zu klären, ob Querverweis-Korrekturen
  (ausdrücklich erlaubt) den Core berühren; falls ja, braucht die Regel eine
  Antwort, sonst meldet jeder Lifecycle-Nachzug Drift.
- [x] **Zweite Datei-Klasse `done/`:** Ob `vcs` beide Klassen in einem
  Konfigurations-Block trägt oder ob `done/` eine eigene Behandlung braucht
  (dort gibt es keine Status-Zeile, die Immutabilität folgt aus dem
  **Verzeichnis**), ist zu entscheiden und zu begründen.
- [x] **Integrationspunkt gewählt:** `vcs` braucht eine Commit-Range
  (`--range`/`--staged`) und gehört damit **nicht** in `make docs-check`,
  sondern an CI (Range über den PR) oder einen lokalen pre-commit-Lauf. Die
  Wahl ist zu begründen und in [`AGENTS.md`](../../../../AGENTS.md)
  §Quality Gates einzutragen.
- [x] **Gegenprobe am realen Verlauf:** Der Sensor läuft über eine Range, die
  eine bekannte legitime Änderung enthält (Lifecycle-Move, Hash-Nachtrag) und
  über eine konstruierte illegitime — er muss die eine durchlassen und die
  andere melden. Ein Sensor, der nie ausgeschlagen hat, ist nicht verifiziert.
- [x] **`MR-005`/`MR-008` nachgezogen** (Gate-Haltung bzw. ADR-Politik).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

### Vier Messungen, eine davon war der eigentliche Befund

| Probe | Erwartung | Ergebnis |
|---|---|---|
| Range ohne ADR-Änderung | grün | grün |
| Inhaltliche Änderung an einem `Accepted`-ADR | Befund | `core-drift-vcs` |
| Änderung an einem `Proposed`-ADR | grün | grün |
| Zulässiger Übergang `Accepted` → `Superseded by …` | grün | **Befund — falsch-positiv** |

Die vierte Probe war der Grund, warum die Gegenprobe in der DoD stand. Ursache:
`status-line` strippt eine **Kopf-Zeile** (`**Status:** Accepted`), wie die
d-check-Beispiele sie voraussetzen. u-boots Status lebt in einem **Abschnitt**
(`## Status`, Wert in der Folgezeile) — die Kopf-Regel greift dort ins Leere,
und der Statuswechsel zählte zum Core.

Gelöst über `exclude-sections: [Status, Geschichte]`: der ganze Abschnitt fällt
aus dem Core. Danach sind Übergang **und** Inhaltsschutz gleichzeitig richtig —
beides nachgemessen, nicht angenommen.

Der Rest der Konfiguration ist unauffällig: `immutable-when: '^Accepted$'`
trennt sauber zwischen `Accepted` (geschützt) und `Proposed` (änderbar), wie
Probe 2 und 3 zeigen. Die im Plan befürchtete Grandfathering-Grenze ist keine:
Lean- und MADR-Form tragen beide `## Status` als Überschrift, weil die
MADR-Umstellung erst beim nächsten inhaltlichen Anfassen greift.

### `done/`-Slices: bewusst nicht erfasst

Die zweite Datei-Klasse aus dem Plan bekommt **keinen** Sensor — mit Grund,
nicht aus Vergesslichkeit.

Ihre Regel lautet „nur korrigierend änderbar"; [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) nennt
ausdrücklich Tippfehler, **Querverweise** und Archiv-Hinweise. Das ist eine
**semantische** Bedingung. `vcs` kann nur „verändert oder nicht" sagen — es
würde jede erlaubte Querverweis-Korrektur als Drift melden. Genau solche
Korrekturen haben in dieser Session dutzendfach stattgefunden (Lifecycle-Moves,
Kennungs-Verlinkung, Hash-Nachträge); ein Sensor, der sie alle anschlägt, wird
nach zwei Tagen abgeschaltet.

Ein Sensor, der die eigene Regel bricht, ist schlechter als keiner. Die
Immutabilität von `done/` bleibt damit eine Disziplin ohne Automat — das ist
ein bewusst getragener Rest, kein Versäumnis.

### Integrationspunkt: on-demand, nicht CI

`vcs` braucht eine Commit-Range und gehört deshalb nicht in `make docs-check`.
Gewählt: **On-demand-Sensor** nach dem Vorbild von `make verify-depguard`, in
[`AGENTS.md`](../../../../AGENTS.md) §Quality Gates eingetragen mit der
Range-Syntax.

**Nicht** in die CI eingebaut, und das ist eine Entscheidung, keine
Auslassung: Der Workflow läuft auf `pull_request` und `push`; u-boot arbeitet
heute direkt auf `main` (diese Session: über zwanzig Commits ohne PR). Ein
Schritt, der nie ausgeführt wird, ist kein Sensor, sondern Dekoration — und ein
Range-Ausdruck, den ich lokal nicht verifizieren kann, wäre blind geliefert.
Die CI-Einbindung folgt, wenn ein PR-basierter Fluss tatsächlich genutzt wird.

### Verification Evidence

Scope:
- Slice: `slice-gate-immutabilitaets-sensor`
- IDs: **keine** Anforderung geändert. Tragend:
  [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format)
  (Accepted-ADRs unveränderlich — die Regel, die jetzt einen Sensor hat),
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
  (die Regel, die **gegen** einen `done/`-Sensor spricht).
- Artefakte: [`.d-check.yml`](../../../../.d-check.yml) (`vcs`-Block mit
  Begründungs-Kommentar), [`AGENTS.md`](../../../../AGENTS.md) §Quality Gates,
  [`harness/conventions.md`](../../../../harness/conventions.md) (`MR-005`).

DoD-Abgleich: alle Punkte erfüllt. Zwei mit abweichendem Ergebnis:
- „Zweite Datei-Klasse `done/`" ist **entschieden und begründet abgelehnt**,
  nicht konfiguriert — der Plan hatte beide Ausgänge offengelassen.
- „Integrationspunkt in `AGENTS.md` eintragen" ist erfüllt; die im Plan
  genannte CI-Variante ist begründet vertagt.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make doc-immutable RANGE="HEAD~5..HEAD"` | pass | 0 Befunde (Negativprobe) |
| `make doc-immutable STAGED=1`, Inhaltsänderung an `0013` | **1 Befund** | `core-drift-vcs` — der Sensor schlägt an |
| `make doc-immutable STAGED=1`, `Proposed`-ADR `0011` geändert | pass | 0 Befunde — Abgrenzung greift |
| `make doc-immutable STAGED=1`, Status → `Superseded` | pass (nach Fix) | vor `exclude-sections`: falsch-positiv |
| `make docs-check` | pass | 147 Dateien / 0 Befunde |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| [`LH-FA-PROJDOCS-002`](../../../../spec/lastenheft.md#lh-fa-projdocs-002--adr-format) | Accepted-ADRs sind erstmals maschinell gegen inhaltliche Änderung geschützt |
| [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) | begründet, warum `done/` **keinen** Sensor bekommt |
| `MR-005` | Sensor, Konfigurations-Begründung und die `done/`-Absage im Ledger |

Carveouts: Neu: none. Gelöst: none. Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make test` — kein Go-Delta.
- CI-Lauf des neuen Sensors — siehe Integrationspunkt oben.

Independent Review: nicht durchgeführt. Der Sensor ist durch vier Proben in
**beide** Richtungen belegt (schlägt an, wo er soll; schweigt, wo er soll) —
das ist stärker als eine zweite Meinung über die Konfiguration.

Commit / Artefakt: `7d9c4c2` (`.d-check.yml` `vcs`-Block, `AGENTS.md` Gate-Tabelle, `MR-005`).

### Steering-Loop-Lerneintrag

- **Die Gegenprobe war der ganze Wert des Slice.** Nach den ersten beiden
  Messungen sah die Konfiguration fertig aus: Der Sensor schlug an, wo er
  sollte. Erst die dritte Frage — „lässt er auch durch, was durch soll?" —
  deckte auf, dass jeder legitime Statuswechsel gemeldet worden wäre. Ein
  Sensor ist erst verifiziert, wenn er in beide Richtungen gezeigt hat.
- **Eine Werkzeug-Option kann formal passen und semantisch danebenliegen.**
  `status-line` war offensichtlich das richtige Feld — und war es nicht, weil
  u-boots Status kein Kopf-Feld ist. Die Lösung lag ein Konzept weiter
  (`exclude-sections`). Wo ein Werkzeug ein Dokumentmodell annimmt, gehört das
  eigene Modell dagegengehalten, bevor man Felder ausfüllt.
- **Nicht jede Regel verdient einen Automaten.** Für `done/` wäre ein Sensor
  technisch möglich und praktisch schädlich gewesen. Die Entscheidung, ihn
  nicht zu bauen, ist Teil des Ergebnisses — dokumentiert, damit sie beim
  nächsten Anlauf nicht als Lücke gelesen wird.
- **Folge-Slices:** keine. Offen als benannter Punkt: CI-Einbindung, sobald ein
  PR-basierter Fluss existiert.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Areas: *harness / Konventionen* und *Harness-Tooling* — beide
**GF** nach [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration.
