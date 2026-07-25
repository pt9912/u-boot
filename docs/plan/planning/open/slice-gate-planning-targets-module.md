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

- [ ] **Trockenlauf zuerst, Ist-Zahl notiert** — je Modul einzeln, bevor es in
  die Modul-Liste wandert. Ein Modul, dessen Befundzahl man nicht kennt, wird
  nicht scharfgeschaltet.
- [ ] **`planning` konfiguriert und grün:** `roadmap:` zeigt auf
  [`roadmap.md`](../in-progress/roadmap.md). Zu prüfen ist, ob u-boots
  Ruhe-Formulierung („— (keine aktive Welle)") zum erwarteten `marker` passt
  oder ob `heading`/`marker` überschrieben werden müssen — u-boots
  `## Aktuelle Welle` trägt einen mehrzeiligen Block, kein einzelnes Stichwort.
- [ ] **`targets` konfiguriert und grün:** `makefiles: [Makefile]`,
  `doc-tables`/`authority` auf [`AGENTS.md`](../../../../AGENTS.md). Befunde
  sind **inhaltlich** zu klären, nicht per `exempt-targets` wegzudrücken:
  `gate-phantom` heißt, die Doku verspricht ein Target, das es nicht gibt;
  `gate-undocumented` heißt, ein Target ist undokumentiert. Beides sind echte
  Aussagen über u-boot, keine Werkzeug-Artefakte.
- [ ] **Ausnahmen nur begründet:** Falls `exempt-targets` nötig wird (etwa für
  reine Utility-Regeln wie `help` oder `clean`), steht je Eintrag ein Grund im
  Konfigurations-Kommentar.
- [ ] **Modul-Liste in `.d-check.yml` und `MR-005` synchron.**
- [ ] `make docs-check` grün mit erweiterter Modul-Liste.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.

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

<!-- Erst nach Abschluss füllen. -->

## 8. Sub-Area-Modus-Begründung

Berührte Sub-Area: *harness / Konventionen* — **GF**. Falls `targets`
Änderungen am [`Makefile`](../../../../Makefile) erzwingt, ist zusätzlich
*Harness-Tooling* berührt, ebenfalls **GF**.
