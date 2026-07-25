# Slice Gate: Module `planning` und `targets` scharfschalten

**Status:** open → next → in-progress → done (Datei wird durch die
Verzeichnisse bewegt).

**Welle:** `welle-gate-ausbau-v0.51` (s. [`roadmap.md`](../in-progress/roadmap.md)).

**Bezug:** [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle)
(Lifecycle-Disziplin, die `planning` prüft) und
[`LH-FA-BUILD-005`](../../../../spec/lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets)
(Standard-Targets, deren Doku-Konsistenz `targets` prüft). `MR-005` führt die
aktive Modul-Liste.

**Autor:** pt9912. **Datum:** 2026-07-25.

---

## 1. Ziel

Zwei Disziplinen mechanisieren, die u-boot heute an Aufmerksamkeit hängt:

- **`planning`** — der Ruhe-Marker steht im `## Aktuelle Welle`-Block genau
  dann, wenn kein `slice-*` im Lifecycle-Verzeichnis liegt. Genau diese
  Buchführung wurde in der Konformitäts-Welle mehrfach von Hand nachgezogen;
  ein vergessener Nachzug ist heute unsichtbar.
- **`targets`** — jedes in einer Doku-Tabelle behauptete `make X` ist eine
  echte Makefile-Regel, und jede Regel steht in der Autoritäts-Doku. Die
  Gate-Tabelle in [`AGENTS.md`](../../../../AGENTS.md) §Quality Gates ist genau
  so eine Tabelle; sie wurde bisher nie gegen das
  [`Makefile`](../../../../Makefile) geprüft.

Beide sind hermetisch (kein git, kein Netz) und billig zu konfigurieren.

## 2. Definition of Done

- [x] **Trockenlauf zuerst, Ist-Zahl notiert** — je Modul einzeln, bevor es in
  die Modul-Liste wandert. Ein Modul, dessen Befundzahl man nicht kennt, wird
  nicht scharfgeschaltet.
- [x] **`planning` konfiguriert und grün:** `roadmap:` zeigt auf
  [`roadmap.md`](../in-progress/roadmap.md). Zu prüfen ist, ob u-boots
  Ruhe-Formulierung („— (keine aktive Welle)") zum erwarteten `marker` passt
  oder ob `heading`/`marker` überschrieben werden müssen — u-boots
  `## Aktuelle Welle` trägt einen mehrzeiligen Block, kein einzelnes Stichwort.
- [x] **`targets` konfiguriert und grün:** `makefiles: [Makefile]`,
  `doc-tables`/`authority` auf [`AGENTS.md`](../../../../AGENTS.md). Befunde
  sind **inhaltlich** zu klären, nicht per `exempt-targets` wegzudrücken:
  `gate-phantom` heißt, die Doku verspricht ein Target, das es nicht gibt;
  `gate-undocumented` heißt, ein Target ist undokumentiert. Beides sind echte
  Aussagen über u-boot, keine Werkzeug-Artefakte.
- [x] **Ausnahmen nur begründet:** Falls `exempt-targets` nötig wird (etwa für
  reine Utility-Regeln wie `help` oder `clean`), steht je Eintrag ein Grund im
  Konfigurations-Kommentar.
- [x] **Modul-Liste in `.d-check.yml` und `MR-005` synchron.**
- [x] `make docs-check` grün mit erweiterter Modul-Liste.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`.d-check.yml`](../../../../.d-check.yml) | update | `modules` + `planning`/`targets`-Blöcke |
| [`AGENTS.md`](../../../../AGENTS.md) bzw. [`Makefile`](../../../../Makefile) | update falls Befund | Phantom-/undokumentierte Targets auflösen |
| [`harness/conventions.md`](../../../../harness/conventions.md) `MR-005` | update | aktive Modul-Liste |

## 4. Trigger

Gefeuert: Entscheidung des Projektinhabers nach dem Image-Bump. Sinnvoll
**nach** der `--print-mk`-Einbindung
([`slice-gate-print-mk-einbindung`](slice-gate-print-mk-einbindung.md)), weil
das Fragment für beide Module fertige Targets mitbringt.

## 5. Closure-Trigger

Beide Module aktiv und grün, Befunde inhaltlich aufgelöst, `MR-005`
nachgezogen, Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`planning` erwartet ein Layout.** u-boots „Aktuelle Welle" ist Prosa mit
  mehreren Absätzen, nicht ein Stichwort-Marker. Passt die Default-Erwartung
  nicht, ist die ehrliche Antwort, den `marker` zu konfigurieren — **nicht**,
  die Roadmap dem Werkzeug anzupassen.
- **`targets` kann unbequeme Wahrheiten liefern.** Undokumentierte Regeln sind
  ein reales Doku-Defizit. Diagnose vor Carveout: erst verstehen, dann
  entweder dokumentieren oder begründet ausnehmen.
- **Fail-closed:** Beide Module brechen bei fehlender/mehrdeutiger Struktur
  ab (Exit 2). Das ist gewollt, macht aber die erste Konfiguration
  fehleranfällig — deshalb der Trockenlauf je Modul einzeln.
- **Kein Carveout erwartet.**

## 7. Closure-Notiz (nach `done/`)

### `planning`: die Invariante ist schärfer als erwartet

Trockenlauf in zwei Zuständen, weil die Semantik erst dadurch sichtbar wurde:

| Zustand | Ergebnis |
|---|---|
| Welle deklariert, **kein** `slice-*` in `in-progress/` | `planning-drift` |
| Welle deklariert, ein `slice-*` in `in-progress/` | grün |

Das Modul koppelt „die Roadmap benennt eine aktive Welle" an „ein Slice liegt
im Roadmap-Verzeichnis". u-boots Modell trennt beides bisher: Die Welle lebt in
der Roadmap, die Slices ziehen einzeln durch `open/ → in-progress/ → done/`.
Zwischen zwei Slices — also genau im Moment eines Closure-Commits — ist die
Welle deklariert und nichts in Arbeit.

**Bewusst als Invariante übernommen, nicht wegkonfiguriert.** Eine deklarierte
Welle ohne Slice in Arbeit ist kein harmloser Zwischenzustand, sondern eine
Roadmap, die Aktivität behauptet, die nicht stattfindet — genau die Drift, die
gestern noch von Hand gesucht werden musste. Die Konsequenz ist eine
Arbeitsfluss-Regel: **Wer einen Slice schließt, zieht im selben Commit den
nächsten nach oder schließt die Welle.** Sie steht als Kommentar im
`planning`-Block und wurde von diesem Closure-Commit als erstem befolgt.

Angepasst wurde nur der `marker`: u-boot schreibt „keine aktive Welle"
kleingeschrieben, der Default lautet „Keine aktive Welle" (literaler
Teilstring-Vergleich). Das Werkzeug an die Roadmap angepasst, nicht umgekehrt.

### `targets`: acht undokumentierte Regeln, kein Phantom

**Kein `gate-phantom`** — jedes in [`AGENTS.md`](../../../../AGENTS.md)
versprochene `make X` existiert. Acht `gate-undocumented`: `help`, `clean`,
`deps`, `compile`, `coverage`, `build`, `run`, `build-binaries`.

Inhaltlich geklärt statt weggedrückt: Die Tabelle in `AGENTS.md` §Quality Gates
führt ausdrücklich **Harness-Sensoren** („Nur reale Make-Targets zaehlen als
Harness-Sensoren"), nicht jede Regel. Die acht sind keine Sensoren — sie sind
Bedienung (`help`, `clean`), Zwischenstufen des Docker-Build-Flusses (`deps`,
`compile`), ein Alias auf einen dokumentierten Sensor (`coverage` →
`coverage-gate`) und Release-/Smoke-Pfade, die `fullbuild` abdeckt (`build`,
`run`, `build-binaries`). Sie in die Gate-Tabelle aufzunehmen hätte deren
Aussage verwässert.

Deshalb `exempt-targets` — aber **einzeln benannt mit Begründung je Gruppe**,
nicht als Sammelposten. Das Modul kennt ohnehin keine Globs; damit fällt ein
künftiges echtes Gate nicht versehentlich unter die Ausnahme.

Zweite Entscheidung: `makefiles: [Makefile]` listet **nur** das eigene
Makefile. Die `doc-*`-Targets stammen aus dem generierten `d-check.mk` und
werden vom Werkzeug selbst dokumentiert (`make doc-help`); sie sind
verfügbare Werkzeuge, keine u-boot-Gates.

### Verification Evidence

Scope:
- Slice: `slice-gate-planning-targets-module`
- IDs: **keine** Anforderung geändert. Lesend belegt:
  [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle),
  [`LH-FA-BUILD-005`](../../../../spec/lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets).
- Artefakte: [`.d-check.yml`](../../../../.d-check.yml) (`modules` erweitert,
  `planning`- und `targets`-Block),
  [`harness/conventions.md`](../../../../harness/conventions.md) (`MR-005`
  Modul-Liste + Arbeitsfluss-Regel).

DoD-Abgleich: alle Punkte erfüllt. Der Trockenlauf lief je Modul einzeln über
`make doc-planning` bzw. `make doc-targets` — beide Targets kamen mit dem
Fragment aus dem Vorgänger-Slice, hätten also von Hand gebaut werden müssen.

Sensors:
| Sensor | Ergebnis | Evidence |
| --- | --- | --- |
| `make doc-planning` (Zustand ohne Slice) | 1 Befund | `planning-drift` — die Invariante, dokumentiert |
| `make doc-planning` (Zustand mit Slice) | pass | 0 Befunde |
| `make doc-targets` (vor Ausnahmen) | 8 Befunde | alle `gate-undocumented`, kein `gate-phantom` |
| `make docs-check` (6 Module aktiv) | pass | 144 Dateien / 0 Befunde |

Traceability:
| ID / Pflicht | Beleg |
| --- | --- |
| [`LH-FA-PROJDOCS-003`](../../../../spec/lastenheft.md#lh-fa-projdocs-003--planning-lifecycle) | Lifecycle-Invariante erstmals maschinell geprüft statt per Aufmerksamkeit |
| [`LH-FA-BUILD-005`](../../../../spec/lastenheft.md#lh-fa-build-005--makefile-mit-standard-targets) | Gate-Tabelle und Makefile sind jetzt gegeneinander verankert |
| `MR-005` | aktive Modul-Liste (4 → 6) und Arbeitsfluss-Regel nachgezogen |

Carveouts: Neu: none — die `exempt-targets` sind **keine** Carveouts, sondern
eine Bereichsabgrenzung: Die ausgenommenen Regeln sind per Definition keine
Sensoren, es wird also keine Prüfung temporär ausgesetzt. Gelöst: none.
Unverändert: none.

Nicht ausgeführt:
- `make gates` / `make ci` — kein Go-Delta; `docs-check` ist der berührte
  Sensor und lief.

Independent Review: nicht durchgeführt. Beide Module sind durch ihre eigenen
Trockenläufe belegt, inklusive der Gegenprobe in zwei Zuständen für
`planning` — ein Sensor, der in beide Richtungen gezeigt hat, ist verifiziert.

Commit / Artefakt: `83efe46` (`.d-check.yml` Module + Bloecke, `MR-005` Modul-Liste und Arbeitsfluss-Regel).

### Steering-Loop-Lerneintrag

- **Ein Sensor zwingt zu einer Antwort, die vorher offenblieb.** Dass eine
  deklarierte Welle ohne Slice in Arbeit „eigentlich" nicht vorkommen soll,
  wusste niemand — es stand nirgends, und deshalb kam es vor. Das Modul hat
  keine Regel gefunden, es hat eine erzwungen. Genau das ist der Wert: Die
  Alternative wäre gewesen, den Marker so zu konfigurieren, dass die Frage
  nicht mehr auftaucht.
- **Der erste Anwendungsfall war der eigene Closure-Commit.** Die neue Regel
  („Slice schließen und nächsten nachziehen im selben Commit") wurde von genau
  dem Commit befolgt, der sie eingeführt hat. Ein Sensor, dessen erste
  Anwendung der eigene Slice ist, ist billig zu verifizieren.
- **„Nicht wegdrücken" braucht eine Alternative, sonst wird es Rhetorik.** Die
  acht `gate-undocumented` *waren* eine Ausnahme wert — aber nur, weil die
  Gate-Tabelle einen definierten Zweck hat (Sensoren) und die acht ihn
  nachweislich verfehlen. Ohne dieses Kriterium wäre `exempt-targets` bloß
  Bequemlichkeit gewesen.
- **Folge-Slices:** keine; die zwei verbleibenden der Welle standen bereits.

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *harness / Konventionen* — **GF**. Falls `targets`
Änderungen am [`Makefile`](../../../../Makefile) erzwingt, ist zusätzlich
*Harness-Tooling* berührt, ebenfalls **GF**.
